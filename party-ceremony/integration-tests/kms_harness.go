package integrationtests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	awskmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ctfcanton "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain/canton"
	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

const (
	defaultLocalstackImage = "localstack/localstack:3.8.1"
	localstackImageEnv     = "PARTY_CEREMONY_LOCALSTACK_IMAGE"

	localstackKMSAlias    = "localstack-kms"
	localstackEdgePort    = "4566/tcp"
	kmsRegion             = "us-east-1"
	localstackAccessKey   = "test"
	localstackSecretKey   = "test"
	localstackSessionName = ""
)

type CTFKMSChain struct {
	Chain *canton.Chain
	KMS   *KMSRegistry
}

type KMSMaterial struct {
	NamespaceKeyID string
	ProtocolKeyID  string
}

type KMSRegistry struct {
	client *awskms.Client

	mu        sync.Mutex
	materials map[string]map[int]KMSMaterial
	keyIDs    []string
}

// LoadChainWithCTFKMS starts the normal live CTF Canton topology with
// participant crypto configured for AWS KMS, and points both Canton and the Go
// contract-deploy signer at the same LocalStack KMS backend. SDK mocks are not
// sufficient here because Canton topology signing happens inside the Canton
// node, while Ledger API signatures must verify against the protocol key that
// Canton published in P2P topology. The harness pins a community LocalStack
// image and disables Pro activation so the suite does not require a
// LOCALSTACK_AUTH_TOKEN.
func LoadChainWithCTFKMS(t *testing.T, numberOfValidators int) (*CTFKMSChain, error) {
	t.Helper()
	ctx := t.Context()

	if numberOfValidators <= 0 {
		return nil, fmt.Errorf("number of validators must be greater than zero")
	}
	if err := framework.DefaultNetwork(defaultNetworkOnce); err != nil {
		return nil, err
	}

	localstack, hostEndpoint, internalEndpoint, err := startLocalStackKMS(ctx, t)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = localstack.Terminate(context.Background()) })

	postgresReq := ctfcanton.PostgresContainerRequest(numberOfValidators)
	postgres, err := startContainer(ctx, t, postgresReq)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = postgres.Terminate(context.Background()) })

	cantonReq := ctfcanton.ContainerRequest(numberOfValidators, "", postgresReq.Name)
	cantonReq, err = withParticipantKMSConfig(cantonReq, numberOfValidators, internalEndpoint)
	if err != nil {
		return nil, err
	}
	if cantonReq.Env == nil {
		cantonReq.Env = make(map[string]string)
	}
	cantonReq.Env["AWS_ACCESS_KEY_ID"] = localstackAccessKey
	cantonReq.Env["AWS_SECRET_ACCESS_KEY"] = localstackSecretKey
	cantonReq.Env["AWS_REGION"] = kmsRegion
	cantonReq.Env["AWS_DEFAULT_REGION"] = kmsRegion
	cantonReq.Env["AWS_EC2_METADATA_DISABLED"] = "true"
	cantonContainer, err := startContainer(ctx, t, cantonReq)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = cantonContainer.Terminate(context.Background()) })

	spliceReq := ctfcanton.SpliceContainerRequest(numberOfValidators, "", postgresReq.Name, cantonReq.Name)
	spliceContainer, err := startContainer(ctx, t, spliceReq)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() { _ = spliceContainer.Terminate(context.Background()) })

	port := freeport.GetOne(t)
	nginxReq, nginxName := ctfcanton.NginxContainerRequest(numberOfValidators, strconv.Itoa(port), cantonReq.Name, spliceReq.Name)
	nginxContainer, err := startContainer(ctx, t, nginxReq)
	if err != nil {
		freeport.Return([]int{port})
		return nil, err
	}
	t.Cleanup(func() { _ = nginxContainer.Terminate(context.Background()) })

	chain, err := buildCantonChain(ctx, t, numberOfValidators, port, nginxName, nginxContainer)
	if err != nil {
		return nil, err
	}

	return &CTFKMSChain{
		Chain: chain,
		KMS: &KMSRegistry{
			client:    newLocalStackKMSClient(hostEndpoint),
			materials: make(map[string]map[int]KMSMaterial),
		},
	}, nil
}

func (r *KMSRegistry) Config(t testing.TB, scope string, actorIndex int) client.KMSConfig {
	t.Helper()
	material := r.Material(t, scope, actorIndex)
	return client.KMSConfig{
		NamespaceKeyID: material.NamespaceKeyID,
		ProtocolKeyID:  material.ProtocolKeyID,
	}
}

func (r *KMSRegistry) ProtocolOnlyConfig(t testing.TB, scope string, actorIndex int) client.KMSConfig {
	t.Helper()
	material := r.Material(t, scope, actorIndex)
	return client.KMSConfig{ProtocolKeyID: material.ProtocolKeyID}
}

func (r *KMSRegistry) Material(t testing.TB, scope string, actorIndex int) KMSMaterial {
	t.Helper()
	if r == nil {
		return KMSMaterial{}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	scope = sanitizeKMSScope(scope)
	if _, ok := r.materials[scope]; !ok {
		r.materials[scope] = make(map[int]KMSMaterial)
	}
	if material, ok := r.materials[scope][actorIndex]; ok {
		return material
	}

	label := fmt.Sprintf("%s-actor-%d", scope, actorIndex+1)
	material := KMSMaterial{
		NamespaceKeyID: r.createSigningKey(t, label+"-namespace"),
		ProtocolKeyID:  r.createSigningKey(t, label+"-protocol"),
	}
	r.materials[scope][actorIndex] = material

	return material
}

func (r *KMSRegistry) AWSClient() *awskms.Client {
	if r == nil {
		return nil
	}

	return r.client
}

func (r *KMSRegistry) KeyIDs() []string {
	if r == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]string, len(r.keyIDs))
	copy(out, r.keyIDs)

	return out
}

func (r *KMSRegistry) createSigningKey(t testing.TB, label string) string {
	t.Helper()
	out, err := r.client.CreateKey(context.Background(), &awskms.CreateKeyInput{
		Description: aws.String(label),
		KeySpec:     awskmstypes.KeySpecEccNistP256,
		KeyUsage:    awskmstypes.KeyUsageTypeSignVerify,
	})
	require.NoError(t, err, "create LocalStack KMS signing key %q", label)
	require.NotNil(t, out.KeyMetadata, "missing key metadata for %q", label)

	keyID := aws.ToString(out.KeyMetadata.Arn)
	if keyID == "" {
		keyID = aws.ToString(out.KeyMetadata.KeyId)
	}
	require.NotEmpty(t, keyID, "missing key id for %q", label)

	r.keyIDs = append(r.keyIDs, keyID)

	return keyID
}

func startLocalStackKMS(ctx context.Context, t *testing.T) (testcontainers.Container, string, string, error) {
	t.Helper()

	name := framework.DefaultTCName(localstackKMSAlias)
	req := testcontainers.ContainerRequest{
		Image:    localstackImage(),
		Name:     name,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {name, localstackKMSAlias},
		},
		ExposedPorts: []string{localstackEdgePort},
		Env: map[string]string{
			"SERVICES":              "kms",
			"ACTIVATE_PRO":          "0",
			"LOG_LICENSE_ISSUES":    "0",
			"LOCALSTACK_AUTH_TOKEN": "",
			"LOCALSTACK_API_KEY":    "",
			"AWS_ACCESS_KEY_ID":     localstackAccessKey,
			"AWS_SECRET_ACCESS_KEY": localstackSecretKey,
			"AWS_DEFAULT_REGION":    kmsRegion,
		},
		WaitingFor: wait.ForHTTP("/_localstack/health").
			WithPort(localstackEdgePort).
			WithStartupTimeout(2 * time.Minute),
		Labels: framework.DefaultTCLabels(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", "", err
	}
	mappedPort, err := container.MappedPort(ctx, localstackEdgePort)
	if err != nil {
		return nil, "", "", err
	}

	return container,
		fmt.Sprintf("http://%s:%s", host, mappedPort.Port()),
		fmt.Sprintf("http://%s:4566", localstackKMSAlias),
		nil
}

func startContainer(ctx context.Context, _ *testing.T, req testcontainers.ContainerRequest) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func withParticipantKMSConfig(req testcontainers.ContainerRequest, numberOfValidators int, endpoint string) (testcontainers.ContainerRequest, error) {
	for i := range req.Files {
		if req.Files[i].ContainerFilePath != "/app/app.conf" {
			continue
		}
		if req.Files[i].Reader == nil {
			return req, fmt.Errorf("CTF Canton app.conf file has no reader")
		}

		raw, err := io.ReadAll(req.Files[i].Reader)
		if err != nil {
			return req, fmt.Errorf("reading CTF Canton app.conf: %w", err)
		}
		req.Files[i].Reader = strings.NewReader(string(raw) + participantKMSConfig(numberOfValidators, endpoint))

		return req, nil
	}

	return req, fmt.Errorf("CTF Canton app.conf file not found")
}

func localstackImage() string {
	if image := strings.TrimSpace(os.Getenv(localstackImageEnv)); image != "" {
		return image
	}

	return defaultLocalstackImage
}

func participantKMSConfig(numberOfValidators int, endpoint string) string {
	var b strings.Builder
	b.WriteString("\n# KMS-backed participants for party ceremony integration tests.\n")
	for i := 1; i <= numberOfValidators; i++ {
		fmt.Fprintf(&b, `
canton.participants.participant%[1]d.crypto.provider = kms
canton.participants.participant%[1]d.crypto.kms {
  type = aws
  region = "%[2]s"
  multi-region-key = false
  audit-logging = true
  endpoint-override = "%[3]s"
}
`, i, kmsRegion, endpoint)
	}

	return b.String()
}

func buildCantonChain(ctx context.Context, t *testing.T, numberOfValidators int, port int, nginxName string, nginxContainer testcontainers.Container) (*canton.Chain, error) {
	t.Helper()

	host, err := nginxContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	participants := make([]canton.Participant, 0, numberOfValidators)
	for i := 1; i <= numberOfValidators; i++ {
		endpoints, err := participantEndpoints(host, strconv.Itoa(port), i)
		if err != nil {
			return nil, err
		}
		assertParticipantHealthURL(t, ctx, endpoints.HTTPHealthCheckURL)
		internalEndpoints, err := participantEndpoints(nginxName, strconv.Itoa(ctfcanton.DefaultNginxInternalPort), i)
		if err != nil {
			return nil, err
		}

		authProvider := authentication.NewInsecureStaticProvider(endpoints.JWT)
		tokenSource := authProvider.TokenSource()
		transportCredentials := authProvider.TransportCredentials()
		perRPCCredentials := authProvider.PerRPCCredentials()

		ledgerConn, err := grpc.NewClient(
			endpoints.GRPCLedgerAPIURL,
			grpc.WithTransportCredentials(transportCredentials),
			grpc.WithPerRPCCredentials(perRPCCredentials),
		)
		if err != nil {
			return nil, fmt.Errorf("create Ledger API gRPC client for participant %d(%s): %w", i, endpoints.GRPCLedgerAPIURL, err)
		}
		ledgerServices := canton.CreateLedgerServiceClients(ledgerConn)

		adminConn, err := grpc.NewClient(
			endpoints.AdminAPIURL,
			grpc.WithTransportCredentials(transportCredentials),
			grpc.WithPerRPCCredentials(perRPCCredentials),
		)
		if err != nil {
			return nil, fmt.Errorf("create Admin API gRPC client for participant %d(%s): %w", i, endpoints.AdminAPIURL, err)
		}
		adminServices := canton.CreateAdminServiceClients(adminConn)

		resp, err := ledgerServices.Admin.UserManagement.GetUser(ctx, &adminv2.GetUserRequest{UserId: endpoints.UserID})
		if err != nil {
			return nil, fmt.Errorf("get user %q: %w", endpoints.UserID, err)
		}
		if resp.User.PrimaryParty == "" {
			return nil, fmt.Errorf("user %q has no primary party", endpoints.UserID)
		}

		participants = append(participants, canton.Participant{
			Name: fmt.Sprintf("Participant %d", i),
			Endpoints: canton.ParticipantEndpoints{
				JSONLedgerAPIURL: endpoints.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: endpoints.GRPCLedgerAPIURL,
				AdminAPIURL:      endpoints.AdminAPIURL,
				ValidatorAPIURL:  endpoints.ValidatorAPIURL,
			},
			InternalEndpoints: &canton.ParticipantEndpoints{
				JSONLedgerAPIURL: internalEndpoints.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: internalEndpoints.GRPCLedgerAPIURL,
				AdminAPIURL:      internalEndpoints.AdminAPIURL,
				ValidatorAPIURL:  internalEndpoints.ValidatorAPIURL,
			},
			LedgerServices: ledgerServices,
			AdminServices:  &adminServices,
			TokenSource:    tokenSource,
			UserID:         endpoints.UserID,
			PartyID:        resp.User.PrimaryParty,
		})
	}

	return &canton.Chain{
		ChainMetadata: canton.ChainMetadata{Selector: chainsel.CANTON_LOCALNET.Selector},
		Participants:  participants,
	}, nil
}

func participantEndpoints(host, port string, participantNumber int) (blockchain.CantonParticipantEndpoints, error) {
	userID := fmt.Sprintf("user-participant%d", participantNumber)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   userID,
		Audience:  []string{ctfcanton.AuthProviderAudience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(blockchain.TokenExpiry)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}).SignedString([]byte(ctfcanton.AuthProviderSecret))
	if err != nil {
		return blockchain.CantonParticipantEndpoints{}, fmt.Errorf("create token for participant%d: %w", participantNumber, err)
	}

	return blockchain.CantonParticipantEndpoints{
		JSONLedgerAPIURL:   fmt.Sprintf("http://participant%d.json-ledger-api.%s:%s", participantNumber, host, port),
		GRPCLedgerAPIURL:   fmt.Sprintf("participant%d.grpc-ledger-api.%s:%s", participantNumber, host, port),
		AdminAPIURL:        fmt.Sprintf("participant%d.admin-api.%s:%s", participantNumber, host, port),
		ValidatorAPIURL:    fmt.Sprintf("http://participant%d.validator-api.%s:%s/api/validator", participantNumber, host, port),
		HTTPHealthCheckURL: fmt.Sprintf("http://participant%d.http-health-check.%s:%s", participantNumber, host, port),
		GRPCHealthCheckURL: fmt.Sprintf("participant%d.grpc-health-check.%s:%s", participantNumber, host, port),
		UserID:             userID,
		JWT:                token,
	}, nil
}

func newLocalStackKMSClient(endpoint string) *awskms.Client {
	return awskms.New(awskms.Options{
		Region:       kmsRegion,
		Credentials:  credentials.NewStaticCredentialsProvider(localstackAccessKey, localstackSecretKey, localstackSessionName),
		BaseEndpoint: aws.String(endpoint),
	})
}

func sanitizeKMSScope(scope string) string {
	scope = strings.ToLower(scope)
	var b strings.Builder
	for _, r := range scope {
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
		return "kms"
	}

	return out
}

func assertParticipantHealthURL(t *testing.T, ctx context.Context, healthURL string) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL+"/health", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	_ = resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
