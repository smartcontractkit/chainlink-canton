package changesets

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/auth"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	coinBinding "github.com/smartcontractkit/chainlink-canton/bindings/coin"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/coin"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type DeployCoinConfig struct {
	Symbol string
}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCoinConfig]] = DeployCoin{}

type DeployCoin struct{}

func (d DeployCoin) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCoinConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d DeployCoin) Apply(e cldf.Environment, config CantonCSDeps[DeployCoinConfig]) (cldf.ChangesetOutput, error) {
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

	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, coin.Deploy, deps, contract.DeployInput[coinBinding.CoinRegistry]{
		ChainSelector: config.ChainSelector,
		Qualifier:     ptr.String(config.Config.Symbol),
		ActAs:         []string{config.Party},
		Template: coinBinding.CoinRegistry{
			Issuer: types.PARTY(config.Party),
			InstrumentId: coinBinding.InstrumentId{
				Admin: types.PARTY(config.Party),
				Id:    types.TEXT(config.Config.Symbol),
			},
		},
		OwnerParty: types.PARTY(config.Party),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply DeployCoin operation: %w", err)
	}

	if err = ds.AddressRefStore.Add(out.Output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save deployed Coin contract address: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []cld_ops.Report[any, any]{},
	}, nil
}
