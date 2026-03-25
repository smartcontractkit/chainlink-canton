package sequences

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	offrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	onrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"

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

type ExecutorParams struct {
	Qualifier string
	Template  executorBinding.Executor
}

type RMNRemoteParams struct {
	Template rmn.RMNRemote
}

type GlobalConfigParams struct {
	Template common.GlobalConfig
}

type FeeQuoterParams struct {
	Template feequoter.FeeQuoter
	// The price of the native token to be set on the FeeQuoter.
	// If not-nil, native will be added as a fee token and the price will be set.
	USDPerNative *big.Int
}

type DeployChainContractsParams struct {
	CCIPOwnerParty     string
	CommitteeVerifiers []CommitteeVerifierParams
	Executors          []ExecutorParams
	GlobalConfig       GlobalConfigParams
	RMNRemote          RMNRemoteParams
	FeeQuoterConfig    FeeQuoterParams
	// The InstrumentId of the native token
	NativeInstrumentId splice_api_token_holding_v1.InstrumentId
}

var DeployChainContracts = operations.NewSequence(
	"canton/ccip/deploy_chain_contracts",
	semver.MustParse("1.7.0"),
	"Deploys all required contracts for CCIP 1.7.0 to a Canton chain",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		var addresses []datastore.AddressRef

		// Deploy RMNRemote
		deployRMNRemoteReport, err := operations.ExecuteOperation(b, rmn_remote.Deploy, deps, contract.DeployInput[rmn.RMNRemote]{
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
			Template:   input.GlobalConfig.Template,
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
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
		linkTokenId := input.FeeQuoterConfig.Template.LinkTokenInstrumentId
		if linkTokenId.Admin == "" {
			linkTokenId = splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(input.CCIPOwnerParty),
				Id:    types.TEXT("link-token"),
			}
		}
		deployFeeQuoterReport, err := operations.ExecuteOperation(b, fee_quoter.Deploy, deps, contract.DeployInput[feequoter.FeeQuoter]{
			Template: feequoter.FeeQuoter{
				Owner:                            types.PARTY(input.CCIPOwnerParty),
				FeeTokens:                        types.SET{},
				DestChainConfigs:                 nil,
				TokenTransferFeeConfigs:          nil,
				UsdPerUnitGasByDestChainSelector: nil,
				UsdPerToken:                      nil,
				LinkTokenInstrumentId:            linkTokenId,
				PriceUpdaters:                    input.FeeQuoterConfig.Template.PriceUpdaters,
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy FeeQuoter: %w", err)
		}
		addresses = append(addresses, deployFeeQuoterReport.Output)
		feeQuoterRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployFeeQuoterReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse FeeQuoter raw instance address: %w", err)
		}

		// Add native as a fee token on the FeeQuoter
		_, err = operations.ExecuteOperation(b, fee_quoter.ApplyFeeTokenUpdates, deps, contract.ChoiceInput[feequoter.ApplyFeeTokenUpdates]{
			InstanceAddress: feeQuoterRawInstanceAddress.InstanceAddress(),
			Args: feequoter.ApplyFeeTokenUpdates{
				FeeTokensToRemove: nil,
				FeeTokensToAdd: []feequoter.FeeTokenArgs{
					{
						InstrumentId: input.NativeInstrumentId,
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to add native as a fee token on FeeQuoter: %w", err)
		}

		// Update the native token price if specified
		if input.FeeQuoterConfig.USDPerNative != nil {
			_, err = operations.ExecuteOperation(b, fee_quoter.UpdatePrices, deps, contract.ChoiceInput[feequoter.UpdatePrices]{
				InstanceAddress: feeQuoterRawInstanceAddress.InstanceAddress(),
				Args: feequoter.UpdatePrices{
					PriceUpdates: feequoter.PriceUpdates{
						TokenPriceUpdates: []feequoter.TokenPriceUpdate{
							{
								InstrumentId: input.NativeInstrumentId,
								UsdPerToken:  types.NUMERIC(input.FeeQuoterConfig.USDPerNative.String()),
							},
						},
					},
				},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to update native token price on FeeQuoter: %w", err)
			}
		}

		// Deploy OffRamp
		deployOffRampReport, err := operations.ExecuteOperation(b, offramp.Deploy, deps, contract.DeployInput[offrampBinding.OffRamp]{
			Template: offrampBinding.OffRamp{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
				Deps: offrampBinding.OffRampDeps{
					GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
					RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
					TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
				},
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OffRamp: %w", err)
		}
		addresses = append(addresses, deployOffRampReport.Output)
		offRampRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployOffRampReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse OffRamp raw instance address: %w", err)
		}

		// Deploy OnRamp
		deployOnRampReport, err := operations.ExecuteOperation(b, onramp.Deploy, deps, contract.DeployInput[onrampBinding.OnRamp]{
			Template: onrampBinding.OnRamp{
				CcipOwner:         types.PARTY(input.CCIPOwnerParty),
				MaxUSDCentsPerMsg: types.NUMERIC("100000000"),
				Deps: onrampBinding.OnRampDeps{
					GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
					RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
					TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
					FeeQuoter:          feeQuoterRawInstanceAddress.Binding(),
					CcvRegistry:        mcms.RawInstanceAddress{},
				},
			},
			OwnerParty: types.PARTY(input.CCIPOwnerParty),
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OnRamp: %w", err)
		}
		addresses = append(addresses, deployOnRampReport.Output)
		onRampRawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployOnRampReport.Output.Labels.List()[0])
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to parse OnRamp raw instance address: %w", err)
		}

		// Deploy PerPartyRouterFactory
		deployPerPartyRouterFactoryReport, err := operations.ExecuteOperation(b, per_party_router_factory.Deploy, deps, contract.DeployInput[perpartyrouter.PerPartyRouterFactory]{
			Template: perpartyrouter.PerPartyRouterFactory{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
				Deps: perpartyrouter.PerPartyRouterDeps{
					OnRamp:             onRampRawInstanceAddress.Binding(),
					OffRamp:            offRampRawInstanceAddress.Binding(),
					GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
					TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
					FeeQuoter:          feeQuoterRawInstanceAddress.Binding(),
					RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
				},
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
				Template: ccvs.CommitteeVerifier{
					Owner:                        committeeVerifier.Template.Owner,
					CcipOwner:                    types.PARTY(input.CCIPOwnerParty),
					VersionTag:                   committeeVerifier.Template.VersionTag,
					MessageSentObservers:         committeeVerifier.Template.MessageSentObservers,
					StorageLocations:             committeeVerifier.Template.StorageLocations,
					StorageLocationsAdmin:        committeeVerifier.Template.StorageLocationsAdmin,
					PendingStorageLocationsAdmin: committeeVerifier.Template.PendingStorageLocationsAdmin,
					RemoteChainConfigs:           committeeVerifier.Template.RemoteChainConfigs,
					SignerConfigs:                committeeVerifier.Template.SignerConfigs,
					Deps: ccvs.CommitteeVerifierDeps{
						RmnRemote: rmnRemoteRawInstanceAddress.Binding()},
				},
				OwnerParty: committeeVerifier.Template.Owner,
				Qualifier:  &committeeVerifier.Qualifier,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier #%v: %w", i, err)
			}
			addresses = append(addresses, deployCommitteeVerifierReport.Output)
		}

		// Deploy Executors
		for i, params := range input.Executors {
			deployExecutorReport, err := operations.ExecuteOperation(b, executor.Deploy, deps, contract.DeployInput[executorBinding.Executor]{
				Template:   params.Template,
				OwnerParty: params.Template.Owner,
				Qualifier:  &params.Qualifier,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy Executor #%v: %w", i, err)
			}
			addresses = append(addresses, deployExecutorReport.Output)
		}

		return sequences.OnChainOutput{
			Addresses: addresses,
		}, nil
	},
)
