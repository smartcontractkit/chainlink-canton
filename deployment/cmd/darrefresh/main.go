package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/deployment/changesets"
)

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

type participantInput struct {
	Selector     uint64              `toml:"chain_selector"`
	Participants []participantConfig `toml:"participants"`
}

type participantConfig struct {
	Name             string `toml:"name,omitempty"`
	JWT              string `toml:"jwt"`
	UserName         string `toml:"username"`
	Party            string `toml:"party,omitempty"`
	JSONLedgerAPIURL string `toml:"json_ledger_api_url"`
	GRPCLedgerAPIURL string `toml:"grpc_ledger_api_url"`
	AdminAPIURL      string `toml:"admin_api_url"`
	ValidatorAPIURL  string `toml:"validator_api_url"`
}

func main() {
	var (
		participantInputPath string
		participantIndex     int
		syncVetting          bool
		synchronizerID       string
		darPaths             multiFlag
		unvetMainPkgIDs      multiFlag
	)

	flag.StringVar(&participantInputPath, "participant-input", "", "Path to participant TOML config (same format as PARTICIPANT_INPUT)")
	flag.IntVar(&participantIndex, "participant-index", 0, "Participant index to run against (0-based)")
	flag.BoolVar(&syncVetting, "sync-vetting", true, "Synchronize package vetting while uploading DARs")
	flag.StringVar(&synchronizerID, "synchronizer-id", "", "Optional synchronizer ID for vet/unvet requests")
	flag.Var(&darPaths, "dar", "DAR file path to upload; repeat -dar for multiple files")
	flag.Var(&unvetMainPkgIDs, "unvet-main-package-id", "Optional main package ID to unvet; repeat flag for multiple IDs (default: auto-list all DARs)")
	flag.Parse()

	if participantInputPath == "" {
		exitWithError("missing required -participant-input")
	}
	if len(darPaths) == 0 {
		exitWithError("missing required -dar (provide one or more DAR file paths)")
	}

	ctx := context.Background()
	cantonChain, selector, err := loadChainFromFile(ctx, participantInputPath)
	if err != nil {
		exitWithError("failed to load chain from participant input: %v", err)
	}
	if participantIndex < 0 || participantIndex >= len(cantonChain.Participants) {
		exitWithError("participant index %d out of range, chain has %d participants", participantIndex, len(cantonChain.Participants))
	}

	dars := make([][]byte, 0, len(darPaths))
	for _, darPath := range darPaths {
		absPath, err := filepath.Abs(darPath)
		if err != nil {
			exitWithError("failed to resolve DAR path %q: %v", darPath, err)
		}
		data, err := os.ReadFile(absPath)
		if err != nil {
			exitWithError("failed to read DAR %q: %v", absPath, err)
		}
		dars = append(dars, data)
	}

	lggr, err := logger.New()
	if err != nil {
		exitWithError("failed to create logger: %v", err)
	}
	bundle := ops.NewBundle(func() context.Context { return ctx }, lggr, ops.NewMemoryReporter())
	env := cldf.Environment{
		Logger:           lggr,
		GetContext:       func() context.Context { return ctx },
		DataStore:        datastore.NewMemoryDataStore().Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{cantonChain}),
		OperationsBundle: bundle,
	}

	var syncPtr *string
	if synchronizerID != "" {
		syncPtr = &synchronizerID
	}

	cs := changesets.UnvetAndReuploadDARs{}
	cfg := changesets.CantonCSDeps[changesets.UnvetAndReuploadDARsConfig]{
		ChainSelector: selector,
		Participant:   participantIndex,
		Config: changesets.UnvetAndReuploadDARsConfig{
			DARs:                  dars,
			MainPackageIDsToUnvet: unvetMainPkgIDs,
			SynchronizeVetting:    syncVetting,
			SynchronizerID:        syncPtr,
		},
	}

	if err := cs.VerifyPreconditions(env, cfg); err != nil {
		exitWithError("precondition check failed: %v", err)
	}
	if _, err := cs.Apply(env, cfg); err != nil {
		exitWithError("changeset execution failed: %v", err)
	}

	fmt.Printf("Successfully unvetted and reuploaded %d DAR(s) on participant index %d\n", len(dars), participantIndex)
}

func loadChainFromFile(ctx context.Context, path string) (*canton.Chain, uint64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read participant input file: %w", err)
	}

	var input participantInput
	if _, err := toml.Decode(string(content), &input); err != nil {
		return nil, 0, fmt.Errorf("unmarshal participant input TOML: %w", err)
	}

	if input.Selector == 0 {
		input.Selector = chainsel.CANTON_LOCALNET.Selector
	}
	if len(input.Participants) == 0 {
		return nil, 0, fmt.Errorf("participant input has no participants")
	}

	providerConfig := cantonProvider.RPCChainProviderConfig{
		Participants: make([]cantonProvider.ParticipantConfig, len(input.Participants)),
	}

	for i, p := range input.Participants {
		if p.JWT == "" || p.UserName == "" || p.GRPCLedgerAPIURL == "" || p.AdminAPIURL == "" {
			return nil, 0, fmt.Errorf("participant #%d is missing required fields", i)
		}

		party := p.Party
		authProvider := authentication.NewInsecureStaticProvider(p.JWT)

		if party == "" {
			ledgerAPIConn, err := grpc.NewClient(
				p.GRPCLedgerAPIURL,
				grpc.WithTransportCredentials(authProvider.TransportCredentials()),
				grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
			)
			if err != nil {
				return nil, 0, fmt.Errorf("create ledger API gRPC client for participant #%d: %w", i, err)
			}

			userMgt := adminv2.NewUserManagementServiceClient(ledgerAPIConn)
			userResp, err := userMgt.GetUser(ctx, &adminv2.GetUserRequest{UserId: p.UserName})
			_ = ledgerAPIConn.Close()
			if err != nil {
				return nil, 0, fmt.Errorf("get user for participant #%d: %w", i, err)
			}
			party = userResp.GetUser().GetPrimaryParty()
			if party == "" {
				return nil, 0, fmt.Errorf("participant #%d has no party and user %q has no primary party", i, p.UserName)
			}
		}

		providerConfig.Participants[i] = cantonProvider.ParticipantConfig{
			JSONLedgerAPIURL: p.JSONLedgerAPIURL,
			GRPCLedgerAPIURL: p.GRPCLedgerAPIURL,
			AdminAPIURL:      p.AdminAPIURL,
			ValidatorAPIURL:  p.ValidatorAPIURL,
			UserID:           p.UserName,
			PartyID:          party,
			AuthProvider:     authProvider,
		}
	}

	bc, err := cantonProvider.NewRPCChainProvider(input.Selector, providerConfig).Initialize(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("initialize RPC chain provider: %w", err)
	}

	return bc.(*canton.Chain), input.Selector, nil
}

func exitWithError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
