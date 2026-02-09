package changesets

import (
	"fmt"

	"https://github.com/smartcontractkit/go-daml/pkg/auth"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

type ConfigureChainForLanesConfig struct {
	Input sequences.ConfigureChainForLanesInput
}

var _ cldf.ChangeSetV2[CantonCSDeps[ConfigureChainForLanesConfig]] = ConfigureChainForLanes{}

type ConfigureChainForLanes struct{}

func (c ConfigureChainForLanes) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (c ConfigureChainForLanes) Apply(e cldf.Environment, config CantonCSDeps[ConfigureChainForLanesConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	token, err := participant.JWTProvider.Token(e.GetContext())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get JWT: %w", err)
	}

	insecureCreds := grpc.WithTransportCredentials(insecure.NewCredentials())
	ledgerApiClient, err := grpc.NewClient(participant.Endpoints.GRPCLedgerAPIURL, insecureCreds, grpc.WithPerRPCCredentials(auth.NewBearerToken(token)))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create gRPC ledger API client: %w", err)
	}

	deps := dependencies.CantonDeps{
		Chain:                chain,
		CommandServiceClient: apiv2.NewCommandServiceClient(ledgerApiClient),
		StateServiceClient:   apiv2.NewStateServiceClient(ledgerApiClient),
		Party:                config.Party,
	}

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.ConfigureChainForLanes, deps, config.Config.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployChainContracts sequence: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{},
	}, nil
}
