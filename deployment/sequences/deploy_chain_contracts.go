package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/noders-team/go-daml/pkg/types"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/feequoter"
	offrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/offramp"
	onrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/ccv_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type CommitteeVerifierParams struct {
	Template ccvs.CommitteeVerifier
}

type GlobalConfigParams struct {
	Template common.GlobalConfig
}

type DeployChainContractsParams struct {
	CCIPOwnerParty     string
	CommitteeVerifiers []CommitteeVerifierParams
	GlobalConfig       GlobalConfigParams
}

var DeployChainContracts = operations.NewSequence(
	"canton/ccip/deploy_chain_contracts",
	semver.MustParse("1.7.0"),
	"Deploys all required contracts for CCIP 1.7.0 to a Canton chain",
	func(b operations.Bundle, deps dependencies.CantonDeps, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		var results []datastore.AddressRef

		// Deploy CCVRegistry
		deployCCVRegistryReport, err := operations.ExecuteOperation(b, ccv_registry.Deploy, deps, contract.DeployInput[common.CCVRegistry]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: common.CCVRegistry{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CCVRegistry: %w", err)
		}
		results = append(results, deployCCVRegistryReport.Output)

		// Deploy FeeQuoter
		deployFeeQuoterReport, err := operations.ExecuteOperation(b, fee_quoter.Deploy, deps, contract.DeployInput[feequoter.FeeQuoter]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: feequoter.FeeQuoter{
				Owner: types.PARTY(input.CCIPOwnerParty),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy FeeQuoter: %w", err)
		}
		results = append(results, deployFeeQuoterReport.Output)

		// Deploy Token Admin Registry
		deployTokenAdminRegistryReport, err := operations.ExecuteOperation(b, token_admin_registry.Deploy, deps, contract.DeployInput[tokenadminregistry.TokenAdminRegistry]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: tokenadminregistry.TokenAdminRegistry{
				Owner:        types.PARTY(input.CCIPOwnerParty),
				TokenConfigs: nil,
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy TokenAdminRegistry: %w", err)
		}
		results = append(results, deployTokenAdminRegistryReport.Output)

		// Deploy OffRamp
		deployOffRampReport, err := operations.ExecuteOperation(b, offramp.Deploy, deps, contract.DeployInput[offrampBinding.OffRamp]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: offrampBinding.OffRamp{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OffRamp: %w", err)
		}
		results = append(results, deployOffRampReport.Output)

		// Deploy OnRamp
		deployOnRampReport, err := operations.ExecuteOperation(b, onramp.Deploy, deps, contract.DeployInput[onrampBinding.OnRamp]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: onrampBinding.OnRamp{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OnRamp: %w", err)
		}
		results = append(results, deployOnRampReport.Output)

		// Deploy PerPartyRouterFactory
		deployPerPartyRouterFactoryReport, err := operations.ExecuteOperation(b, per_party_router_factory.Deploy, deps, contract.DeployInput[perpartyrouter.PerPartyRouterFactory]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: perpartyrouter.PerPartyRouterFactory{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy PerPartyRouterFactory: %w", err)
		}
		results = append(results, deployPerPartyRouterFactoryReport.Output)

		// Deploy Committee Verifiers
		for i, committeeVerifier := range input.CommitteeVerifiers {
			committeeVerifier.Template.CcipOwner = types.PARTY(input.CCIPOwnerParty)
			deployCommitteeVerifierReport, err := operations.ExecuteOperation(b, committee_verifier.Deploy, deps, contract.DeployInput[ccvs.CommitteeVerifier]{
				ChainSelector: deps.Chain.Selector,
				ActAs:         []string{deps.Party},
				Template:      committeeVerifier.Template,
				OwnerParty:    committeeVerifier.Template.Owner,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier #%v: %w", i, err)
			}
			results = append(results, deployCommitteeVerifierReport.Output)
		}

		// Deploy Global Config
		input.GlobalConfig.Template.CcipOwner = types.PARTY(input.CCIPOwnerParty)
		deployGlobalConfigReport, err := operations.ExecuteOperation(b, global_config.Deploy, deps, contract.DeployInput[common.GlobalConfig]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template:      input.GlobalConfig.Template,
			OwnerParty:    types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy GlobalConfig: %w", err)
		}
		results = append(results, deployGlobalConfigReport.Output)

		return sequences.OnChainOutput{
			Addresses: results,
		}, nil
	},
)
