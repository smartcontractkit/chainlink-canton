package changesets

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

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

	// TODO: change bindings to allow passing a JWT provider directly
	token, err := participant.JWTProvider.Token(e.GetContext())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get JWT: %w", err)
	}

	bindingClient, err := client.NewDamlClient(token, participant.Endpoints.GRPCLedgerAPIURL).Build(e.GetContext())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create Daml binding client: %w", err)
	}

	deps := dependencies.CantonDeps{
		Chain:         chain,
		BindingClient: bindingClient,
		Party:         config.Party,
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
