package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
)

const ParticipantInputEnvVar = "PARTICIPANT_INPUT"

type TestEnvironment struct {
	Logger zerolog.Logger
	Chain  canton.Chain
}

// ParticipantInput provides the input configuration for participants, when not using CTF
type ParticipantInput struct {
	Selector     uint64              `toml:"chain_selector"`
	Participants []ParticipantConfig `toml:"participants"`
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

func LoadChainFromFile(t *testing.T, path string) (*canton.Chain, error) {
	content, err := os.ReadFile(filepath.Join(".", path))
	if err != nil {
		return nil, fmt.Errorf("error reading participant input file: %w", err)
	}
	var participantInput ParticipantInput
	if err := toml.Unmarshal(content, &participantInput); err != nil {
		return nil, fmt.Errorf("error unmarshalling participant input TOML: %w", err)
	}

	providerConfig := cantonProvider.RPCChainProviderConfig{
		Participants: make([]cantonProvider.ParticipantConfig, len(participantInput.Participants)),
	}

	for i, config := range participantInput.Participants {
		party := config.Party

		if party == "" {
			// If no party is specified, use UserManagement service to get the primary party for the user
			authProvider := authentication.NewInsecureStaticProvider(config.JWT)
			ledgerApiConn, err := grpc.NewClient(
				config.GRPCLedgerAPIURL,
				grpc.WithTransportCredentials(authProvider.TransportCredentials()),
				grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
			)
			require.NoError(t, err, "Failed to create gRPC connection to Ledger API for participant %d", i+1)
			userMgt := adminv2.NewUserManagementServiceClient(ledgerApiConn)
			user, err := userMgt.GetUser(t.Context(), &adminv2.GetUserRequest{UserId: config.UserName})
			require.NoError(t, err, "Failed to get user info for user %s for participant %d", config.UserName, i+1)
			if len(user.GetUser().GetPrimaryParty()) == 0 {
				return nil, fmt.Errorf("no primary party found for user %s for participant %d", config.UserName, i+1)
			}
			party = user.GetUser().GetPrimaryParty()
			log.Info().Str("user", config.UserName).Str("party", party).Msgf("No party specified for participant %d, using primary party of the user", i+1)
			_ = ledgerApiConn.Close()
		}

		authProvider := authentication.NewInsecureStaticProvider(config.JWT)
		providerConfig.Participants[i] = cantonProvider.ParticipantConfig{
			Endpoints: cantonProvider.Endpoints{
				JSONLedgerAPIURL: config.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: config.GRPCLedgerAPIURL,
				AdminAPIURL:      config.AdminAPIURL,
				ValidatorAPIURL:  config.ValidatorAPIURL,
			},
			UserID:       config.UserName,
			PartyID:      party,
			AuthProvider: authProvider,
		}
	}

	bc, err := cantonProvider.NewRPCChainProvider(participantInput.Selector, providerConfig).Initialize(t.Context())
	require.NoError(t, err, "Failed to initialize RPC chain provider")

	return bc.(*canton.Chain), nil
}

var defaultNetworkOnce = &sync.Once{}

func LoadChainWithCTF(t *testing.T, numberOfValidators int) (*canton.Chain, error) {
	const (
		maxAttempts       = 3
		retryDelay        = 5 * time.Second
		perAttemptTimeout = 6 * time.Minute
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		once := defaultNetworkOnce
		if attempt > 1 {
			// First attempt may leave a partially initialized/default network behind.
			// Use a fresh sync.Once on retries so CTF can rebuild network state.
			once = &sync.Once{}
		}
		attemptCtx, cancel := context.WithTimeout(t.Context(), perAttemptTimeout)
		bc, err := cantonProvider.NewCTFChainProvider(t, chainsel.CANTON_LOCALNET.Selector, cantonProvider.CTFChainProviderConfig{
			NumberOfValidators: numberOfValidators,
			Once:               once,
		}).Initialize(attemptCtx)
		cancel()
		if err == nil {
			return bc.(*canton.Chain), nil
		}

		lastErr = err
		t.Logf("CTF chain provider initialization failed (attempt %d/%d): %v", attempt, maxAttempts, err)
		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to initialize CTF chain provider after %d attempts: %w", maxAttempts, lastErr)
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
		cantonChain *canton.Chain
		err         error
	)

	// Load participant input from environment variable, if set.
	// If not set, use CLDF to spin up a network
	participantInputPath := os.Getenv(ParticipantInputEnvVar)
	if participantInputPath != "" {
		cantonChain, err = LoadChainFromFile(t, participantInputPath)
		require.NoError(t, err, "Failed to load chain from file")
		if len(cantonChain.Participants) < testConfig.NumberOfParticipants {
			t.Fatalf("Number of participants in chain (%d) does satisfy test's requirement (%d)", len(cantonChain.Participants), testConfig.NumberOfParticipants)
		}
	} else {
		cantonChain, err = LoadChainWithCTF(t, testConfig.NumberOfParticipants)
		require.NoError(t, err, "Failed to load chain with CTF")
	}

	env := TestEnvironment{
		Logger: log.
			Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.InfoLevel).
			With().
			Fields(map[string]any{"component": "TestEnvironment"}).
			Logger(),
		Chain: *cantonChain,
	}

	return env
}
