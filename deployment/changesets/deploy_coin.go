package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	coinBinding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/coin"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"

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

	party := chain.Participants[config.Participant].PartyID
	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, coin.Deploy, chain, contract.DeployInput[coinBinding.CoinRegistry]{
		Qualifier:        new(config.Config.Symbol),
		ParticipantIndex: config.Participant,
		Template: coinBinding.CoinRegistry{
			Issuer: types.PARTY(party),
			InstrumentId: splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(party),
				Id:    types.TEXT(config.Config.Symbol),
			},
		},
		OwnerParty: types.PARTY(party),
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
