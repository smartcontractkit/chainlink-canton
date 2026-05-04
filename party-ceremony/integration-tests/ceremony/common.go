package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	interactivepb "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	integrationtests "github.com/smartcontractkit/chainlink-canton/party-ceremony/integration-tests"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

type Actor struct {
	client client.CantonClient
	uid    string // Canton participant UID
}

type CeremonyTestSuite struct {
	suite.Suite

	chain          *canton.Chain
	KMS            *integrationtests.KMSRegistry
	KMSRunName     string
	SynchronizerID string
	ParticipantIDs []string
	Actors         []Actor
}

func (s *CeremonyTestSuite) SetupSuite() {
	t := s.T()

	if s.chain == nil {
		env, err := integrationtests.LoadChainWithCTFKMS(t, 3)
		require.NoError(t, err, "failed to load KMS-backed chain with CTF")
		require.Len(t, env.Chain.Participants, 3, "expected 3 participants")
		s.chain = env.Chain
		s.KMS = env.KMS
	}

	// Build CantonClients for each participant.
	actors := make([]Actor, 3)
	for i := range 3 {
		c, conn := s.NewCantonClient(s.chain.Participants[i])
		t.Cleanup(func() { _ = conn.Close() })
		actors[i] = Actor{
			client: c,
			uid:    getParticipantUID(t, c),
		}
		t.Logf("Participant %d: %s", i+1, actors[i].uid)
	}

	synchronizerID := s.DiscoverSynchronizerID(s.chain.Participants[0])
	participantIDs := []string{actors[0].uid, actors[1].uid, actors[2].uid}

	s.SynchronizerID = synchronizerID
	s.ParticipantIDs = participantIDs
	s.Actors = actors
}

// Creates a GRPCCantonClient from a canton.Participant by
// dialling the participant's Admin API endpoint. The caller must close the
// returned connection when done.
func (s *CeremonyTestSuite) NewCantonClient(p canton.Participant) (client.CantonClient, *grpc.ClientConn) {
	conn := s.NewAdminConn(p)

	return client.NewGRPCClient(conn), conn
}

func (s *CeremonyTestSuite) NewAdminConn(p canton.Participant) *grpc.ClientConn {
	t := s.T()
	require.NotNil(t, p.AdminServices, "participant %q has no admin API configured", p.Name)

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if p.TokenSource != nil {
		tok, err := p.TokenSource.Token()
		require.NoError(t, err, "failed to get JWT for participant %s", p.Name)
		if tok.AccessToken != "" {
			dialOpts = append(dialOpts,
				grpc.WithUnaryInterceptor(jwtUnary(tok.AccessToken)),
				grpc.WithStreamInterceptor(jwtStream(tok.AccessToken)),
			)
		}
	}

	conn, err := grpc.NewClient(p.Endpoints.AdminAPIURL, dialOpts...)
	require.NoError(t, err, "failed to dial admin API for participant %s", p.Name)

	return conn
}

// discoverSynchronizerID queries the first connected synchronizer from the
// participant's admin API and returns its synchronizer_id.
func (s *CeremonyTestSuite) DiscoverSynchronizerID(p canton.Participant) string {
	t := s.T()

	require.NotNil(t, p.AdminServices, "participant %q has no admin API configured", p.Name)

	resp, err := p.AdminServices.SynchronizerConnectivity.ListConnectedSynchronizers(
		t.Context(), &participantv30.ListConnectedSynchronizersRequest{},
	)
	require.NoError(t, err, "ListConnectedSynchronizers failed")
	require.NotEmpty(t, resp.GetConnectedSynchronizers(), "no connected synchronizers found")

	syncID := resp.GetConnectedSynchronizers()[0].GetSynchronizerId()
	t.Logf("Discovered synchronizer ID: %s", syncID)

	return syncID
}

func (s *CeremonyTestSuite) NewLocalEnv() (*canton.Chain, error) {
	t := s.T()
	t.Helper()
	// Participant ports follow the localnet compose convention:
	// admin: 1201, 1202, 1203 — ledger: 1301, 1302, 1303
	type localPorts struct{ admin, ledger int }
	ports := []localPorts{{1201, 1301}, {1202, 1302}, {1203, 1303}}
	participants := make([]canton.Participant, len(ports))
	for i, p := range ports {
		adminURL := fmt.Sprintf("localhost:%d", p.admin)
		ledgerURL := fmt.Sprintf("localhost:%d", p.ledger)
		conn, err := grpc.NewClient(adminURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "failed to dial admin API for participant %d", i+1)
		t.Cleanup(func() { _ = conn.Close() })
		ledgerConn, err := grpc.NewClient(ledgerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "failed to dial ledger API for participant %d", i+1)
		t.Cleanup(func() { _ = ledgerConn.Close() })
		ledgerServices := canton.CreateLedgerServiceClients(ledgerConn)
		adminServices := canton.CreateAdminServiceClients(conn)
		// UserID matches the additional-admin-user-id in simple-topology.conf.
		userID := fmt.Sprintf("user-participant%d", i+1)
		participants[i] = canton.Participant{
			Name: fmt.Sprintf("participant%c", 'A'+i),
			Endpoints: canton.ParticipantEndpoints{
				AdminAPIURL:      adminURL,
				GRPCLedgerAPIURL: ledgerURL,
			},
			AdminServices:  &adminServices,
			LedgerServices: ledgerServices,
			UserID:         userID,
		}
	}

	return &canton.Chain{
		Participants: participants,
	}, nil
}

// getParticipantID returns the Canton-assigned participant identifier by
// calling GetParticipantID on a CantonClient.
func getParticipantUID(t *testing.T, c client.CantonClient) string {
	t.Helper()
	uid, err := c.GetParticipantUID(t.Context())
	require.NoError(t, err, "GetParticipantUID failed")

	return uid
}

// jwtUnary injects a bearer token into outgoing gRPC unary requests.
func jwtUnary(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), method, req, reply, cc, opts...)
	}
}

// jwtStream injects a bearer token into outgoing gRPC streaming requests.
func jwtStream(token string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), desc, cc, method, opts...)
	}
}

// userIDInterceptor injects a user_id into PrepareSubmissionRequests and
// ExecuteSubmissionRequests when the Canton node runs without authentication
// (no JWT subject claim).
func userIDInterceptor(userID string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if r, ok := req.(*interactivepb.PrepareSubmissionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}
		if r, ok := req.(*interactivepb.ExecuteSubmissionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}
		if r, ok := req.(*interactivepb.ExecuteSubmissionAndWaitForTransactionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewLedgerClient creates a GRPCLedgerClient from a canton.Participant by
// dialling the participant's gRPC Ledger API endpoint. The caller must close
// the returned connection when done.
func (s *CeremonyTestSuite) NewLedgerClient(p canton.Participant) (client.LedgerClient, *grpc.ClientConn) {
	t := s.T()
	require.NotEmpty(t, p.Endpoints.GRPCLedgerAPIURL,
		"participant %q has no gRPC Ledger API URL", p.Name)

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if p.TokenSource != nil {
		tok, err := p.TokenSource.Token()
		require.NoError(t, err, "failed to get JWT for participant %s", p.Name)
		if tok.AccessToken != "" {
			dialOpts = append(dialOpts,
				grpc.WithUnaryInterceptor(jwtUnary(tok.AccessToken)),
				grpc.WithStreamInterceptor(jwtStream(tok.AccessToken)),
			)
		}
	}
	// When the participant has a UserID set and no JWT carries a subject claim
	// (local no-auth Canton), inject the user_id field into PrepareSubmission
	// requests so Canton can resolve the acting user.
	if p.UserID != "" && p.TokenSource == nil {
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(userIDInterceptor(p.UserID)))
	}

	conn, err := grpc.NewClient(p.Endpoints.GRPCLedgerAPIURL, dialOpts...)
	require.NoError(t, err, "failed to dial ledger API for participant %s", p.Name)

	return client.NewGRPCLedgerClient(conn), conn
}

func (s *CeremonyTestSuite) testScope(phase string) string {
	base := s.KMSRunName
	if base == "" {
		base = s.T().Name()
	}

	return sanitizeTestName(base + "-" + phase)
}

func (s *CeremonyTestSuite) uniqueName(phase string) string {
	return "kms-" + s.testScope(phase)
}

func (s *CeremonyTestSuite) kmsConfigFor(actorIndex int, phase string) client.KMSConfig {
	if s.KMS == nil {
		return client.KMSConfig{}
	}

	return s.KMS.Config(s.T(), s.testScope(phase), actorIndex)
}

func (s *CeremonyTestSuite) depsFor(actorIndex int, kmsCfg client.KMSConfig) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: s.Actors[actorIndex].client,
		KMS:    kmsCfg,
		Logger: logger.Test(s.T()),
	}
}

func (s *CeremonyTestSuite) kmsDepsFor(actorIndex int) ceremony.CantonDeps {
	return s.depsFor(actorIndex, s.kmsConfigFor(actorIndex, "onboarding"))
}

func (s *CeremonyTestSuite) assertKMSKeysRegistered(actorIndex int, kmsCfg client.KMSConfig) {
	t := s.T()
	t.Helper()
	if kmsCfg.NamespaceKeyID == "" && kmsCfg.ProtocolKeyID == "" {
		return
	}

	conn := s.NewAdminConn(s.chain.Participants[actorIndex])
	defer conn.Close()

	vault := cryptoadminv30.NewVaultServiceClient(conn)
	resp, err := vault.ListMyKeys(t.Context(), &cryptoadminv30.ListMyKeysRequest{})
	require.NoError(t, err, "ListMyKeys for participant %d", actorIndex+1)

	registered := make(map[string]struct{})
	for _, metadata := range resp.GetPrivateKeysMetadata() {
		if kmsKeyID := metadata.GetKmsKeyId(); kmsKeyID != "" {
			registered[kmsKeyID] = struct{}{}
		}
	}

	if kmsCfg.NamespaceKeyID != "" {
		require.Contains(t, registered, kmsCfg.NamespaceKeyID, "namespace KMS key should be registered for participant %d", actorIndex+1)
	}
	if kmsCfg.ProtocolKeyID != "" {
		require.Contains(t, registered, kmsCfg.ProtocolKeyID, "protocol KMS key should be registered for participant %d", actorIndex+1)
	}
}

func (s *CeremonyTestSuite) assertReportsDoNotContainKMS(reporter operations.Reporter) {
	t := s.T()
	t.Helper()
	if s.KMS == nil {
		return
	}

	reports, err := reporter.GetReports()
	require.NoError(t, err, "read reports for KMS leakage check")
	raw, err := json.Marshal(reports)
	require.NoError(t, err, "marshal reports for KMS leakage check")
	payload := string(raw)

	require.NotContains(t, payload, "arn:aws:kms", "workflow reports must not contain KMS ARNs")
	for _, keyID := range s.KMS.KeyIDs() {
		require.NotContains(t, payload, keyID, "workflow reports must not contain local KMS key IDs")
	}
}

func sanitizeTestName(name string) string {
	name = strings.ToLower(name)
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "ceremony"
	}
	if len(out) > 48 {
		return out[:48]
	}

	return out
}

type recordingKMSAPI struct {
	base client.AWSKMSAPI

	mu                sync.Mutex
	getPublicKeyCalls int
	signCalls         int
}

func newRecordingKMSAPI(base client.AWSKMSAPI) *recordingKMSAPI {
	return &recordingKMSAPI{base: base}
}

func (r *recordingKMSAPI) GetPublicKey(ctx context.Context, in *awskms.GetPublicKeyInput, optFns ...func(*awskms.Options)) (*awskms.GetPublicKeyOutput, error) {
	r.mu.Lock()
	r.getPublicKeyCalls++
	r.mu.Unlock()

	return r.base.GetPublicKey(ctx, in, optFns...)
}

func (r *recordingKMSAPI) Sign(ctx context.Context, in *awskms.SignInput, optFns ...func(*awskms.Options)) (*awskms.SignOutput, error) {
	r.mu.Lock()
	r.signCalls++
	r.mu.Unlock()

	return r.base.Sign(ctx, in, optFns...)
}

func (r *recordingKMSAPI) assertUsed(t *testing.T) {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()

	require.Positive(t, r.getPublicKeyCalls, "KMS GetPublicKey should be called")
	require.Positive(t, r.signCalls, "KMS Sign should be called")
}
