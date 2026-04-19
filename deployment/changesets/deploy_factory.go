package changesets

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type DeployCCIPFactoryParams struct {
	OwnerParty string `json:"ownerParty" yaml:"ownerParty"`
	MCMSParty  string `json:"mcmsParty,omitempty" yaml:"mcmsParty,omitempty"`
	InstanceID string `json:"instanceID,omitempty" yaml:"instanceID,omitempty"`
	Qualifier  string `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
}

type DeployCCIPFactoryConfig struct {
	Params DeployCCIPFactoryParams `json:"params" yaml:"params"`
}

type DeployCCIPFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCCIPFactoryConfig]] = DeployCCIPFactory{}

func (d DeployCCIPFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCCIPFactoryConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}
	if config.Config.Params.OwnerParty == "" {
		return fmt.Errorf("owner party is required")
	}
	return nil
}

func (d DeployCCIPFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCCIPFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	out, err := operations.ExecuteSequence(e.OperationsBundle, deployCCIPFactorySequence, chain, config.Config.Params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployCCIPFactory sequence: %w", err)
	}
	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}
	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}

var deployCCIPFactorySequence = operations.NewSequence(
	"canton/ccip_factory/deploy",
	semver.MustParse("0.1.0"),
	"Deploys a CCIPFactory contract on Canton",
	func(b operations.Bundle, deps canton.Chain, input DeployCCIPFactoryParams) (ccipsequences.OnChainOutput, error) {
		if input.OwnerParty == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("owner party is required")
		}

		ownerParty := types.PARTY(input.OwnerParty)
		mcmsParty := ownerParty
		if input.MCMSParty != "" {
			mcmsParty = types.PARTY(input.MCMSParty)
		}

		var qualifier *string
		if input.Qualifier != "" {
			qualifier = &input.Qualifier
		}

		deployReport, err := operations.ExecuteOperation(b, factoryops.Deploy, deps, opcontract.DeployInput[factorybindings.CCIPFactory]{
			Qualifier: qualifier,
			Template: factorybindings.CCIPFactory{
				InstanceId:                    types.TEXT(input.InstanceID),
				Owner:                         ownerParty,
				McmsParty:                     mcmsParty,
				UsedInstanceIds:               types.GENMAP{},
				DeployedContracts:             types.GENMAP{},
				PerPartyRouterFactoryDeployed: types.BOOL(false),
			},
			OwnerParty: ownerParty,
		})
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("deploy CCIPFactory: %w", err)
		}

		return ccipsequences.OnChainOutput{Addresses: []datastore.AddressRef{deployReport.Output}}, nil
	},
)
