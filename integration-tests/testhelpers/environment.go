package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"https://github.com/smartcontractkit/go-daml/pkg/auth"

	"github.com/BurntSushi/toml"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

const ParticipantInputEnvVar = "PARTICIPANT_INPUT"

type TestEnvironment struct {
	Logger zerolog.Logger

	Selector     uint64
	Participants []Participant
	Splice       Splice

	Chain chain.BlockChain
}

type Splice struct {
	TokenMetadataClient       tokenMetadataV1.ClientWithResponsesInterface
	TransferInstructionClient transferInstructionV1.ClientWithResponsesInterface
}

// Participant returns the n-th participant (1-indexed)
func (e TestEnvironment) Participant(n int) Participant {
	return e.Participants[n-1]
}

// Optional input configuration for participants, when not using CTF

type ParticipantInput struct {
	Selector       uint64              `toml:"chain_selector"`
	Participants   []ParticipantConfig `toml:"participants"`
	ScanApiURL     string              `toml:"scan_api_url"`
	RegistryApiURL string              `toml:"registry_api_url"`
}

type ParticipantConfig struct {
	Name             string `toml:"name,omitempty"`
	JWT              string `toml:"jwt"`
	UserName         string `toml:"username"`
	Party            string `toml:"party,omitempty"`
	JSONLedgerAPIURL string `toml:"json_ledger_api_url"`
	GRPCLedgerAPIURL string `toml:"grpc_ledger_api_url"`
	AdminAPIURL      string `toml:"admin_api_url"`
	ValidatorAPIURL  string `toml:"validator_api_url"`
}

func LoadParticipantInputFromFile(path string) (ParticipantInput, error) {
	content, err := os.ReadFile(filepath.Join(".", path))
	if err != nil {
		return ParticipantInput{}, fmt.Errorf("error reading participant input file: %w", err)
	}
	var participantInput ParticipantInput
	if err := toml.Unmarshal(content, &participantInput); err != nil {
		return ParticipantInput{}, fmt.Errorf("error unmarshaling participant input TOML: %w", err)
	}

	return participantInput, nil
}

var defaultNetworkOnce = &sync.Once{}

func LoadParticipantsWithCLDF(t *testing.T, numberOfValidators int) (ParticipantInput, error) {
	bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
		NumberOfValidators: numberOfValidators,
		Once:               defaultNetworkOnce,
	}).Initialize(t.Context())
	require.NoError(t, err, "Failed to initialize CTF chain provider")

	chain := bc.(*canton.Chain)

	var participantInput ParticipantInput
	participantInput.Selector = chainsel.CANTON_LOCALNET.Selector
	for i, participant := range chain.Participants {
		// This doesn't handle dynamic token providers yet, add if needed
		jwt, err := participant.JWTProvider.Token(t.Context())
		require.NoErrorf(t, err, "Failed to get JWT token for participant %q", participant.Name)
		participantConfig := ParticipantConfig{
			Name:             participant.Name,
			JWT:              jwt,
			UserName:         fmt.Sprintf("user-participant%v", i+1),
			Party:            "", // TODO populate from CLDF
			JSONLedgerAPIURL: participant.Endpoints.JSONLedgerAPIURL,
			GRPCLedgerAPIURL: participant.Endpoints.GRPCLedgerAPIURL,
			AdminAPIURL:      participant.Endpoints.AdminAPIURL,
			ValidatorAPIURL:  participant.Endpoints.ValidatorAPIURL,
		}
		participantInput.Participants = append(participantInput.Participants, participantConfig)
	}

	// TODO properly expose these endpoints via CLDF
	// Also, these are CC specific
	port := strings.Split(participantInput.Participants[0].GRPCLedgerAPIURL, ":")[1]
	registryUrl := fmt.Sprintf("http://scan.localhost:%v", port)
	participantInput.RegistryApiURL = registryUrl
	participantInput.ScanApiURL = registryUrl

	return participantInput, nil
}

type TestConfig struct {
	NumberOfParticipants int
}

func DefaultTestConfig() TestConfig {
	return TestConfig{
		NumberOfParticipants: 1,
	}
}

type TestOption func(env *TestConfig)

func WithNumberOfParticipants(n int) TestOption {
	return func(env *TestConfig) {
		env.NumberOfParticipants = n
	}
}

func NewTestEnvironment(t *testing.T, options ...TestOption) TestEnvironment {
	testConfig := DefaultTestConfig()
	for _, option := range options {
		option(&testConfig)
	}

	var (
		participantInput ParticipantInput
		err              error
	)

	// Load participant input from environment variable, if set.
	// If not set, use CLDF to spin up a network
	participantInputPath := os.Getenv(ParticipantInputEnvVar)
	if participantInputPath != "" {
		participantInput, err = LoadParticipantInputFromFile(participantInputPath)
		require.NoError(t, err, "Failed to load participant input from file")
		if len(participantInput.Participants) < testConfig.NumberOfParticipants {
			t.Fatalf("Number of participants in input (%d) does satisfy test's requirement (%d)", len(participantInput.Participants), testConfig.NumberOfParticipants)
		}
	} else {
		participantInput, err = LoadParticipantsWithCLDF(t, testConfig.NumberOfParticipants)
		require.NoError(t, err, "Failed to load participants with CLDF")
	}

	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(participantInput.RegistryApiURL)
	require.NoError(t, err, "Failed to create token metadata client")
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(participantInput.RegistryApiURL)
	require.NoError(t, err, "Failed to create transfer instruction client")

	env := TestEnvironment{
		Logger: log.
			Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.InfoLevel).
			With().
			Fields(map[string]any{"component": "TestEnvironment"}).
			Logger(),
		Selector:     participantInput.Selector,
		Participants: make([]Participant, len(participantInput.Participants)),
		Splice: Splice{
			TokenMetadataClient:       tokenMetadataClient,
			TransferInstructionClient: transferInstructionClient,
		},
	}
	for i, participantConfig := range participantInput.Participants {
		env.Participants[i] = dialParticipant(t, participantConfig)
	}

	// Create CLDF chain instance
	cldfConfig := cantonProvider.RPCChainProviderConfig{
		Endpoints:    make([]canton.ParticipantEndpoints, len(participantInput.Participants)),
		JWTProviders: make([]canton.JWTProvider, len(participantInput.Participants)),
	}
	for i, participantConfig := range participantInput.Participants {
		cldfConfig.Endpoints[i] = canton.ParticipantEndpoints{
			AdminAPIURL:      participantConfig.AdminAPIURL,
			GRPCLedgerAPIURL: participantConfig.GRPCLedgerAPIURL,
			JSONLedgerAPIURL: participantConfig.JSONLedgerAPIURL,
			ValidatorAPIURL:  participantConfig.ValidatorAPIURL,
		}
		cldfConfig.JWTProviders[i] = canton.NewStaticJWTProvider(participantConfig.JWT)
	}
	chainProvider := cantonProvider.NewRPCChainProvider(env.Selector, cldfConfig)
	chainInstance, err := chainProvider.Initialize(t.Context())
	require.NoError(t, err, "Failed to initialize CLDF chain provider")
	env.Chain = chainInstance

	return env
}

func dialParticipant(t *testing.T, config ParticipantConfig) Participant {
	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	adminApiClient, err := grpc.NewClient(config.AdminAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(config.JWT)))
	require.NoError(t, err, "Failed to dial gRPC admin API")
	ledgerApiClient, err := grpc.NewClient(config.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(config.JWT)))
	require.NoError(t, err, "Failed to dial gRPC ledger API")

	authProvider, err := securityprovider.NewSecurityProviderBearerToken(config.JWT)
	require.NoError(t, err, "Failed to create security provider")
	scanProxyClient, err := scanProxy.NewClientWithResponses(config.ValidatorAPIURL, scanProxy.WithRequestEditorFn(authProvider.Intercept))
	require.NoError(t, err, "Failed to create scan proxy client")

	p := Participant{
		Name: config.Name,
		GetToken: func(ctx context.Context) (string, error) {
			return config.JWT, nil
		},
		GetConfig: func() ParticipantConfig {
			return config
		},
		UserName:                     config.UserName,
		Party:                        config.Party,
		PackageServiceClient:         participantv30.NewPackageServiceClient(adminApiClient),
		PartyManagementServiceClient: admin.NewPartyManagementServiceClient(ledgerApiClient),
		UserManagementServiceClient:  admin.NewUserManagementServiceClient(ledgerApiClient),
		StateServiceClient:           apiv2.NewStateServiceClient(ledgerApiClient),
		CommandServiceClient:         apiv2.NewCommandServiceClient(ledgerApiClient),
		UpdateServiceClient:          apiv2.NewUpdateServiceClient(ledgerApiClient),
		VersionServiceClient:         apiv2.NewVersionServiceClient(ledgerApiClient),
		ScanProxyClient:              scanProxyClient,
	}

	// Populate Party if empty
	if p.Party == "" {
		resp, err := p.UserManagementServiceClient.GetUser(t.Context(), &admin.GetUserRequest{
			UserId: p.UserName,
		})
		require.NoError(t, err, "Failed to get user from %q to determine primary party", p.Name)
		p.Party = resp.GetUser().PrimaryParty
	}

	return p
}
