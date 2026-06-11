package sequences

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type CommitteeVerifierParams struct {
	// Qualifier distinguishes between multiple deployments of the committee verifier on the same chain.
	Qualifier string
	Template  committeeverifier.CommitteeVerifier
}

type ExecutorParams struct {
	Qualifier string
	Template  executorBinding.Executor
}

// RMNRemoteParams configures RMNRemote deploy. Set Template.CustomObservers (typically ccvOwner)
// so NOP verifiers can poll curse state via sourcereader.
type RMNRemoteParams struct {
	Template core.RMNRemote
}

func rmnRemoteDeployTemplate(rmnOwner, ccipOwner types.PARTY, params RMNRemoteParams) core.RMNRemote {
	return core.RMNRemote{
		RmnOwner:        rmnOwner,
		CcipOwner:       ccipOwner,
		CursedSubjects:  params.Template.CursedSubjects,
		CustomObservers: params.Template.CustomObservers,
	}
}

type GlobalConfigParams struct {
	Template core.GlobalConfig
}

type FeeQuoterParams struct {
	Template core.FeeQuoter
	// The price of the native token to be set on the FeeQuoter.
	// If not-nil, native will be added as a fee token and the price will be set.
	USDPerNative *big.Int
}

type DeployChainContractsParams struct {
	// OwnerParty is the Canton instance owner (decentralized party); used for instanceId@owner addresses.
	OwnerParty string
	// CCIPOwnerParty is the operational ccipOwner on product templates (lanes, admin roles, etc.).
	CCIPOwnerParty string
	// CCVOwnerParty is the CommitteeVerifier signatory owner (ccvOwner); distinct from CCIPOwnerParty in dual-MCMS deploys.
	CCVOwnerParty string
	// RMNOwnerParty is the RMNRemote signatory owner (rmnOwner); distinct from CCIPOwnerParty in triple-MCMS deploys.
	RMNOwnerParty      string
	CommitteeVerifiers []CommitteeVerifierParams
	Executors          []ExecutorParams
	GlobalConfig       GlobalConfigParams
	RMNRemote          RMNRemoteParams
	FeeQuoterConfig    FeeQuoterParams
	// FactoryAddressRef is the ccip-qualified factory for core CCIP deploys.
	FactoryAddressRef datastore.AddressRef
	// RMNFactoryAddressRef is the rmn-qualified factory; required when DevenvBundledDeploy is true.
	RMNFactoryAddressRef datastore.AddressRef
	// ProposalDriven enables MCMS proposal generation for factory-backed deploys.
	ProposalDriven bool
	// CcvRegistryBinding is required for OnRamp deps when CommitteeVerifiers is empty (CCV deployed separately).
	CcvRegistryBinding chainlinkapi.RawInstanceAddress
	// RmnRemoteRawInstanceAddress is required for production split deploy paths.
	RmnRemoteRawInstanceAddress contracts.RawInstanceAddress
	// DevenvBundledDeploy runs RMN+CV+core in one sequence (devenv adapter only). Mutually exclusive with RmnRemoteRawInstanceAddress.
	DevenvBundledDeploy bool
	// The InstrumentId of the native token
	NativeInstrumentId splice_api_token_holding_v1.InstrumentId
}

var DeployChainContracts = operations.NewSequence(
	"canton/ccip/deploy_chain_contracts",
	semver.MustParse("2.0.0"),
	"Deploys all required contracts for CCIP 2.0.0 to a Canton chain",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		var addresses []datastore.AddressRef

		rmnOwnerParty, err := requireRMNOwnerParty(input)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		// Deploy RMNRemote
		deployRMNRemoteReport, err := operations.ExecuteOperation(b, rmn_remote.Deploy, deps, contract.DeployInput[core.RMNRemote]{
			Template:   rmnRemoteDeployTemplate(rmnOwnerParty, types.PARTY(input.CCIPOwnerParty), input.RMNRemote),
			OwnerParty: rmnOwnerParty,
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
		deployGlobalConfigReport, err := operations.ExecuteOperation(b, global_config.Deploy, deps, contract.DeployInput[core.GlobalConfig]{
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
		deployTokenAdminRegistryReport, err := operations.ExecuteOperation(b, token_admin_registry.Deploy, deps, contract.DeployInput[core.TokenAdminRegistry]{
			Template: core.TokenAdminRegistry{
				Owner:      types.PARTY(input.CCIPOwnerParty),
				EntryCount: 0,
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
		// TODO: what is this code doing? "link-token" string hardcoded
		linkTokenId := input.FeeQuoterConfig.Template.LinkTokenInstrumentId
		if linkTokenId.Admin == "" {
			linkTokenId = splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(input.CCIPOwnerParty),
				Id:    types.TEXT("link-token"),
			}
		}
		deployFeeQuoterReport, err := operations.ExecuteOperation(b, fee_quoter.Deploy, deps, contract.DeployInput[core.FeeQuoter]{
			Template: core.FeeQuoter{
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

		// Any token with a price is treated as a fee token, so pushing the native
		// token price is sufficient to register it as usable for fees.
		if input.FeeQuoterConfig.USDPerNative != nil {
			_, err = operations.ExecuteOperation(b, fee_quoter.UpdatePrices, deps, contract.ChoiceInput[core.UpdatePrices]{
				InstanceAddress: feeQuoterRawInstanceAddress.InstanceAddress(),
				Args: core.UpdatePrices{
					PriceUpdates: core.PriceUpdates{
						TokenPriceUpdates: []core.TokenPriceUpdate{
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

			err = ensureNativeFeeTokenConfig(
				b,
				deps,
				tokenAdminRegistryRawInstanceAddress.InstanceAddress(),
				input.CCIPOwnerParty,
				input.NativeInstrumentId,
			)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure native fee token config: %w", err)
			}
		}

		// Deploy OffRamp
		deployOffRampReport, err := operations.ExecuteOperation(b, offramp.Deploy, deps, contract.DeployInput[ccipruntime.OffRamp]{
			Template: ccipruntime.OffRamp{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
				Deps: ccipruntime.OffRampDeps{
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
		deployOnRampReport, err := operations.ExecuteOperation(b, onramp.Deploy, deps, contract.DeployInput[ccipruntime.OnRamp]{
			Template: ccipruntime.OnRamp{
				CcipOwner:         types.PARTY(input.CCIPOwnerParty),
				MaxUSDCentsPerMsg: types.NUMERIC("100000000"),
				Deps: ccipruntime.OnRampDeps{
					GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
					RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
					TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
					FeeQuoter:          feeQuoterRawInstanceAddress.Binding(),
					CcvRegistry:        chainlinkapi.RawInstanceAddress{},
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
		deployPerPartyRouterFactoryReport, err := operations.ExecuteOperation(b, per_party_router_factory.Deploy, deps, contract.DeployInput[ccipruntime.PerPartyRouterFactory]{
			Template: ccipruntime.PerPartyRouterFactory{
				CcipOwner: types.PARTY(input.CCIPOwnerParty),
				Deps: ccipruntime.PerPartyRouterDeps{
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
			deployCommitteeVerifierReport, err := operations.ExecuteOperation(b, committee_verifier.Deploy, deps, contract.DeployInput[committeeverifier.CommitteeVerifier]{
				Template: committeeverifier.CommitteeVerifier{
					Owner:                        committeeVerifier.Template.Owner,
					CcipOwner:                    types.PARTY(input.CCIPOwnerParty),
					VersionTag:                   committeeVerifier.Template.VersionTag,
					MessageSentObservers:         committeeVerifier.Template.MessageSentObservers,
					StorageLocations:             committeeVerifier.Template.StorageLocations,
					StorageLocationsAdmin:        committeeVerifier.Template.StorageLocationsAdmin,
					PendingStorageLocationsAdmin: committeeVerifier.Template.PendingStorageLocationsAdmin,
					RemoteChainConfigs:           committeeVerifier.Template.RemoteChainConfigs,
					SignerConfigs:                committeeVerifier.Template.SignerConfigs,
					Deps: committeeverifier.CommitteeVerifierDeps{
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

func ensureNativeFeeTokenConfig(
	b operations.Bundle,
	deps canton.Chain,
	tokenAdminRegistryAddress contracts.InstanceAddress,
	ccipOwnerParty string,
	instrumentId splice_api_token_holding_v1.InstrumentId,
) error {
	if instrumentId.Admin == "" || instrumentId.Id == "" {
		return nil
	}

	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).
		RawInstanceAddress(types.PARTY(ccipOwnerParty)).
		InstanceAddress()

	if _, found, err := findTokenConfigCid(b, deps, 0, tokenConfigAddress); err != nil {
		return fmt.Errorf("failed to lookup native fee token config: %w", err)
	} else if found {
		return nil
	}

	_, err := operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[core.ProposeAdministrator]{
		InstanceAddress: tokenAdminRegistryAddress,
		Args: core.ProposeAdministrator{
			TokenConfigCid: nil,
			InstrumentId:   instrumentId,
			NewAdmin:       types.PARTY(ccipOwnerParty),
			Caller:         types.PARTY(ccipOwnerParty),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to propose native fee token admin: %w", err)
	}

	tokenConfigCid, found, err := findTokenConfigCid(b, deps, 0, tokenConfigAddress)
	if err != nil {
		return fmt.Errorf("failed to lookup native fee token config after propose: %w", err)
	}
	if !found {
		return fmt.Errorf("native fee token config not found after propose")
	}

	_, err = operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[core.AcceptAdminRole]{
		InstanceAddress: tokenAdminRegistryAddress,
		Args: core.AcceptAdminRole{
			TokenConfigCid: tokenConfigCid,
			InstrumentId:   instrumentId,
			Caller:         types.PARTY(ccipOwnerParty),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to accept native fee token admin: %w", err)
	}

	return nil
}
