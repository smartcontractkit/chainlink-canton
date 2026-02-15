package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/feequoter"
	offrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/offramp"
	onrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
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
	// Qualifier distinguishes between multiple deployments of the committee verifier on the same chain.
	Qualifier string
	Template  ccvs.CommitteeVerifier
}

type RMNRemoteParams struct {
	Template rmn.RMNRemote
}

type GlobalConfigParams struct {
	Template common.GlobalConfig
}

type DeployChainContractsParams struct {
	CCIPOwnerParty     string
	CommitteeVerifiers []CommitteeVerifierParams
	GlobalConfig       GlobalConfigParams
	RMNRemote          RMNRemoteParams
}

var DeployChainContracts = operations.NewSequence(
	"canton/ccip/deploy_chain_contracts",
	semver.MustParse("1.7.0"),
	"Deploys all required contracts for CCIP 1.7.0 to a Canton chain",
	func(b operations.Bundle, deps dependencies.CantonDeps, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		var addresses []datastore.AddressRef

		// Deploy RMNRemote
		deployRMNRemoteReport, err := operations.ExecuteOperation(b, rmn_remote.Deploy, deps, contract.DeployInput[rmn.RMNRemote]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: rmn.RMNRemote{
				RmnOwner:       input.RMNRemote.Template.RmnOwner,
				CcipOwner:      types.PARTY(input.CCIPOwnerParty),
				CursedSubjects: input.RMNRemote.Template.CursedSubjects,
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy RMNRemote: %w", err)
		}
		addresses = append(addresses, deployRMNRemoteReport.Output)
		rmnRemoteRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployRMNRemoteReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse RMNRemote raw instance address: %w", err)
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
		addresses = append(addresses, deployGlobalConfigReport.Output)
		globalConfigRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployGlobalConfigReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse GlobalConfig raw instance address: %w", err)
		}

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
		addresses = append(addresses, deployTokenAdminRegistryReport.Output)
		tokenAdminRegistryRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployTokenAdminRegistryReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse TokenAdminRegistry raw instance address: %w", err)
		}

		// Deploy FeeQuoter
		deployFeeQuoterReport, err := operations.ExecuteOperation(b, fee_quoter.Deploy, deps, contract.DeployInput[feequoter.FeeQuoter]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: feequoter.FeeQuoter{
				Owner:                            types.PARTY(input.CCIPOwnerParty),
				FeeTokens:                        nil,
				DestChainConfigs:                 nil,
				TokenTransferFeeConfigs:          nil,
				UsdPerUnitGasByDestChainSelector: nil,
				UsdPerToken:                      nil,
				PriceUpdaters:                    nil,
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy FeeQuoter: %w", err)
		}
		addresses = append(addresses, deployFeeQuoterReport.Output)

		// Deploy OffRamp
		deployOffRampReport, err := operations.ExecuteOperation(b, offramp.Deploy, deps, contract.DeployInput[offrampBinding.OffRamp]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: offrampBinding.OffRamp{
				CcipOwner:                         types.PARTY(input.CCIPOwnerParty),
				GlobalConfigInstanceAddress:       globalConfigRawInstanceAddress.Binding(),
				RmnRemoteInstanceAddress:          rmnRemoteRawInstanceAddress.Binding(),
				TokenAdminRegistryInstanceAddress: tokenAdminRegistryRawInstanceAddress.Binding(),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OffRamp: %w", err)
		}
		addresses = append(addresses, deployOffRampReport.Output)

		// Deploy OnRamp
		deployOnRampReport, err := operations.ExecuteOperation(b, onramp.Deploy, deps, contract.DeployInput[onrampBinding.OnRamp]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: onrampBinding.OnRamp{
				CcipOwner:                         types.PARTY(input.CCIPOwnerParty),
				GlobalConfigInstanceAddress:       globalConfigRawInstanceAddress.Binding(),
				RmnRemoteInstanceAddress:          rmnRemoteRawInstanceAddress.Binding(),
				TokenAdminRegistryInstanceAddress: tokenAdminRegistryRawInstanceAddress.Binding(),
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OnRamp: %w", err)
		}
		addresses = append(addresses, deployOnRampReport.Output)

		// Deploy PerPartyRouterFactory
		deployPerPartyRouterFactoryReport, err := operations.ExecuteOperation(b, per_party_router_factory.Deploy, deps, contract.DeployInput[perpartyrouter.PerPartyRouterFactory]{
			ChainSelector: deps.Chain.Selector,
			ActAs:         []string{deps.Party},
			Template: perpartyrouter.PerPartyRouterFactory{
				CcipOwner:         types.PARTY(input.CCIPOwnerParty),
				RegisteredRouters: nil,
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy PerPartyRouterFactory: %w", err)
		}
		addresses = append(addresses, deployPerPartyRouterFactoryReport.Output)

		// Deploy Committee Verifiers
		for i, committeeVerifier := range input.CommitteeVerifiers {
			committeeVerifier.Template.CcipOwner = types.PARTY(input.CCIPOwnerParty)
			deployCommitteeVerifierReport, err := operations.ExecuteOperation(b, committee_verifier.Deploy, deps, contract.DeployInput[ccvs.CommitteeVerifier]{
				ChainSelector: deps.Chain.Selector,
				ActAs:         []string{deps.Party},
				Template: ccvs.CommitteeVerifier{
					Owner:                    committeeVerifier.Template.Owner,
					CcipOwner:                types.PARTY(input.CCIPOwnerParty),
					VersionTag:               committeeVerifier.Template.VersionTag,
					MessageSentObserver:      committeeVerifier.Template.MessageSentObserver,
					StorageLocation:          committeeVerifier.Template.StorageLocation,
					Threshold:                committeeVerifier.Template.Threshold,
					Signers:                  committeeVerifier.Template.Signers,
					RmnRemoteInstanceAddress: rmnRemoteRawInstanceAddress.Binding(),
					RemoteChainFeeConfigs:    committeeVerifier.Template.RemoteChainFeeConfigs,
				},
				OwnerParty: committeeVerifier.Template.Owner,
				Qualifier:  &committeeVerifier.Qualifier,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier #%v: %w", i, err)
			}
			addresses = append(addresses, deployCommitteeVerifierReport.Output)
		}

		return sequences.OnChainOutput{
			Addresses: addresses,
		}, nil
	},
)
