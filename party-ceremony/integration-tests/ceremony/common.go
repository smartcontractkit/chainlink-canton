package tests

import (
	"context"
	"fmt"
	"testing"

	integrationtests "github.com/chainlink/canton-party-ceremony/integration-tests"
	"github.com/chainlink/canton-party-ceremony/internal/client"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Actor struct {
	client client.CantonClient
	uid    string // Canton participant UID
}

type CeremonyTestSuite struct {
	suite.Suite

	SynchronizerID string
	ParticipantIDs []string
	Actors         []Actor
}

func (s *CeremonyTestSuite) SetupSuite() {
	t := s.T()

	// Uncomment to use local Canton environment for faster test runs (requires local setup).
	// chain, err := s.NewLocalEnv()
	// require.NoError(t, err, "failed to create local environment")

	chain, err := integrationtests.LoadChainWithCTF(t, 3)
	require.NoError(t, err, "failed to load chain with CTF")
	require.Len(t, chain.Participants, 3, "expected 3 participants")

	// Build CantonClients for each participant.
	actors := make([]Actor, 3)
	for i := range 3 {
		c, conn := s.NewCantonClient(chain.Participants[i])
		t.Cleanup(func() { _ = conn.Close() })
		actors[i] = Actor{
			client: c,
			uid:    getParticipantUID(t, c),
		}
		t.Logf("Participant %d: %s", i+1, actors[i].uid)
	}

	synchronizerID := s.DiscoverSynchronizerID(chain.Participants[0])
	participantIDs := []string{actors[0].uid, actors[1].uid, actors[2].uid}

	s.SynchronizerID = synchronizerID
	s.ParticipantIDs = participantIDs
	s.Actors = actors
}

// Creates a GRPCCantonClient from a canton.Participant by
// dialling the participant's Admin API endpoint. The caller must close the
// returned connection when done.
func (s *CeremonyTestSuite) NewCantonClient(p canton.Participant) (client.CantonClient, *grpc.ClientConn) {
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

	return client.NewGRPCClient(conn), conn
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
	adminports := []int{1201, 1202, 1203}
	participants := make([]canton.Participant, len(adminports))
	for i, port := range adminports {
		adminURL := fmt.Sprintf("localhost:%d", port)
		conn, err := grpc.NewClient(adminURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "failed to dial admin API for participant %d", i+1)
		t.Cleanup(func() { _ = conn.Close() })
		adminServices := canton.CreateAdminServiceClients(conn)
		participants[i] = canton.Participant{
			Name: fmt.Sprintf("participant%c", 'A'+i),
			Endpoints: canton.ParticipantEndpoints{
				AdminAPIURL: adminURL,
			},
			AdminServices: &adminServices,
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
