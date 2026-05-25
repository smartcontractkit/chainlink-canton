package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	ccvsbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	offrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	onrampBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployChainContractsFromFactory is the factory-backed Canton deploy path.
// It assumes CCIPFactory has already been bootstrapped and targets that
// existing factory for all core CCIP contract deployments.
var DeployChainContractsFromFactory = operations.NewSequence(
	"canton/ccip/deploy_chain_contracts_from_factory",
	semver.MustParse("2.0.0"),
	"Deploys CCIP contracts on Canton through CCIPFactory choices",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		var addresses []datastore.AddressRef
		var proposalOutputs []contract.ExerciseOutput

		ownerParty, ccipOwnerParty, err := requireOwnerParties(input)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}
		factoryRawInstanceAddress, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		rmnRemoteRawInstanceAddress, rmnAddressRef, rmnProposalOutputs, err := resolveOrDeployRMNRemote(
			b, deps, input, ccipOwnerParty,
		)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}
		if rmnAddressRef != nil {
			addresses = append(addresses, *rmnAddressRef)
		}
		proposalOutputs = append(proposalOutputs, rmnProposalOutputs...)

		globalConfigInstanceID, err := ensureInstanceID(input.GlobalConfig.Template.InstanceId, "globalconfig")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure GlobalConfig instance ID: %w", err)
		}
		globalConfigTemplate := common.GlobalConfig{
			InstanceId:         types.TEXT(globalConfigInstanceID),
			CcipOwner:          ccipOwnerParty,
			ChainSelector:      input.GlobalConfig.Template.ChainSelector,
			DestChainConfigs:   nil,
			SourceChainConfigs: nil,
		}
		deployGlobalConfigReport, err := operations.ExecuteOperation(b, factoryops.DeployGlobalConfig, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployGlobalConfig{Contract: globalConfigTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy GlobalConfig from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployGlobalConfigReport.Output, input.ProposalDriven)
		globalConfigRawInstanceAddress := globalConfigInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), globalConfigRawInstanceAddress, global_config.ContractType, global_config.Version, ""))

		tokenAdminRegistryInstanceID, err := ensureInstanceID("", "tokenadminregistry")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure TokenAdminRegistry instance ID: %w", err)
		}
		tokenAdminRegistryTemplate := tokenadminregistry.TokenAdminRegistry{
			InstanceId: types.TEXT(tokenAdminRegistryInstanceID),
			Owner:      ccipOwnerParty,
			EntryCount: 0,
		}
		deployTokenAdminRegistryReport, err := operations.ExecuteOperation(b, factoryops.DeployTokenAdminRegistry, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployTokenAdminRegistry{Contract: tokenAdminRegistryTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy TokenAdminRegistry from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployTokenAdminRegistryReport.Output, input.ProposalDriven)
		tokenAdminRegistryRawInstanceAddress := tokenAdminRegistryInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), tokenAdminRegistryRawInstanceAddress, token_admin_registry.ContractType, token_admin_registry.Version, ""))

		feeQuoterInstanceID, err := ensureInstanceID(input.FeeQuoterConfig.Template.InstanceId, "feequoter")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure FeeQuoter instance ID: %w", err)
		}
		// TODO: what is this code doing? "link-token" string hardcoded
		linkTokenID := input.FeeQuoterConfig.Template.LinkTokenInstrumentId
		if linkTokenID.Admin == "" {
			linkTokenID = splice_api_token_holding_v1.InstrumentId{
				Admin: ownerParty,
				Id:    types.TEXT("link-token"),
			}
		}
		feeQuoterTemplate := feequoter.FeeQuoter{
			InstanceId:                       types.TEXT(feeQuoterInstanceID),
			Owner:                            ownerParty,
			FeeTokens:                        types.SET{},
			DestChainConfigs:                 nil,
			TokenTransferFeeConfigs:          nil,
			UsdPerUnitGasByDestChainSelector: nil,
			UsdPerToken:                      nil,
			LinkTokenInstrumentId:            linkTokenID,
			PriceUpdaters:                    input.FeeQuoterConfig.Template.PriceUpdaters,
		}
		deployFeeQuoterReport, err := operations.ExecuteOperation(b, factoryops.DeployFeeQuoter, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployFeeQuoter{Contract: feeQuoterTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy FeeQuoter from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployFeeQuoterReport.Output, input.ProposalDriven)
		feeQuoterRawInstanceAddress := feeQuoterInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), feeQuoterRawInstanceAddress, fee_quoter.ContractType, fee_quoter.Version, ""))

		var firstCommitteeVerifierBinding mcms.RawInstanceAddress
		if len(input.CommitteeVerifiers) == 0 {
			if input.CcvRegistryBinding.Unpack == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("CcvRegistryBinding is required when CommitteeVerifiers is empty")
			}
			firstCommitteeVerifierBinding = input.CcvRegistryBinding
		}
		var ccvOwnerParty types.PARTY
		if len(input.CommitteeVerifiers) > 0 {
			ccvOwnerParty, err = requireCCVOwnerParty(input)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
		}
		for i, committeeVerifierParams := range input.CommitteeVerifiers {
			qualifier := committeeVerifierParams.Qualifier
			committeeVerifierOwner := ccvOwnerParty
			committeeVerifierInstanceID, err := ensureInstanceID(committeeVerifierParams.Template.InstanceId, "committeeverifier")
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure CommitteeVerifier #%d instance ID: %w", i, err)
			}
			committeeVerifierTemplate := ccvsbindings.CommitteeVerifier{
				InstanceId:                   types.TEXT(committeeVerifierInstanceID),
				Owner:                        committeeVerifierOwner,
				CcipOwner:                    ccipOwnerParty,
				VersionTag:                   committeeVerifierParams.Template.VersionTag,
				AllowListAdmin:               committeeVerifierParams.Template.AllowListAdmin,
				MessageSentObservers:         committeeVerifierParams.Template.MessageSentObservers,
				StorageLocations:             committeeVerifierParams.Template.StorageLocations,
				StorageLocationsAdmin:        committeeVerifierParams.Template.StorageLocationsAdmin,
				PendingStorageLocationsAdmin: committeeVerifierParams.Template.PendingStorageLocationsAdmin,
				RemoteChainConfigs:           committeeVerifierParams.Template.RemoteChainConfigs,
				SignerConfigs:                committeeVerifierParams.Template.SignerConfigs,
				Deps: ccvsbindings.CommitteeVerifierDeps{
					RmnRemote: rmnRemoteRawInstanceAddress.Binding(),
				},
			}
			deployCommitteeVerifierReport, err := operations.ExecuteOperation(b, factoryops.DeployCommitteeVerifier, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployCommitteeVerifier{Contract: committeeVerifierTemplate}, input.ProposalDriven))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier #%d from factory: %w", i, err)
			}
			proposalOutputs = appendExerciseOutput(proposalOutputs, deployCommitteeVerifierReport.Output, input.ProposalDriven)
			committeeVerifierRawInstanceAddress := committeeVerifierInstanceID.RawInstanceAddress(committeeVerifierOwner)
			if i == 0 {
				firstCommitteeVerifierBinding = committeeVerifierRawInstanceAddress.Binding()
			}
			addresses = append(addresses, newAddressRef(deps.ChainSelector(), committeeVerifierRawInstanceAddress, committee_verifier.ContractType, committee_verifier.Version, qualifier))
		}

		offRampInstanceID, err := ensureInstanceID("", "offramp")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure OffRamp instance ID: %w", err)
		}
		offRampTemplate := offrampBinding.OffRamp{
			InstanceId: types.TEXT(offRampInstanceID),
			CcipOwner:  ccipOwnerParty,
			Deps: offrampBinding.OffRampDeps{
				GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
				RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
				TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
			},
		}
		deployOffRampReport, err := operations.ExecuteOperation(b, factoryops.DeployOffRamp, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployOffRamp{Contract: offRampTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OffRamp from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployOffRampReport.Output, input.ProposalDriven)
		offRampRawInstanceAddress := offRampInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), offRampRawInstanceAddress, offramp.ContractType, offramp.Version, ""))

		onRampInstanceID, err := ensureInstanceID("", "onramp")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure OnRamp instance ID: %w", err)
		}
		onRampTemplate := onrampBinding.OnRamp{
			InstanceId:        types.TEXT(onRampInstanceID),
			CcipOwner:         ccipOwnerParty,
			MaxUSDCentsPerMsg: types.NUMERIC("100000000"),
			Deps: onrampBinding.OnRampDeps{
				GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
				RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
				TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
				FeeQuoter:          feeQuoterRawInstanceAddress.Binding(),
				CcvRegistry:        firstCommitteeVerifierBinding,
			},
		}
		deployOnRampReport, err := operations.ExecuteOperation(b, factoryops.DeployOnRamp, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployOnRamp{Contract: onRampTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy OnRamp from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployOnRampReport.Output, input.ProposalDriven)
		onRampRawInstanceAddress := onRampInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), onRampRawInstanceAddress, onramp.ContractType, onramp.Version, ""))

		perPartyRouterFactoryInstanceID, err := ensureInstanceID("", "perpartyrouterfactory")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure PerPartyRouterFactory instance ID: %w", err)
		}
		perPartyRouterFactoryTemplate := perpartyrouter.PerPartyRouterFactory{
			InstanceId: types.TEXT(perPartyRouterFactoryInstanceID),
			CcipOwner:  ccipOwnerParty,
			Deps: perpartyrouter.PerPartyRouterDeps{
				OnRamp:             onRampRawInstanceAddress.Binding(),
				OffRamp:            offRampRawInstanceAddress.Binding(),
				GlobalConfig:       globalConfigRawInstanceAddress.Binding(),
				TokenAdminRegistry: tokenAdminRegistryRawInstanceAddress.Binding(),
				FeeQuoter:          feeQuoterRawInstanceAddress.Binding(),
				RmnRemote:          rmnRemoteRawInstanceAddress.Binding(),
			},
			RegisteredRouters: nil,
		}
		deployPerPartyRouterFactoryReport, err := operations.ExecuteOperation(b, factoryops.DeployPerPartyRouterFactory, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployPerPartyRouterFactory{Contract: perPartyRouterFactoryTemplate}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy PerPartyRouterFactory from factory: %w", err)
		}
		proposalOutputs = appendExerciseOutput(proposalOutputs, deployPerPartyRouterFactoryReport.Output, input.ProposalDriven)
		perPartyRouterFactoryRawInstanceAddress := perPartyRouterFactoryInstanceID.RawInstanceAddress(ownerParty)
		addresses = append(addresses, newAddressRef(deps.ChainSelector(), perPartyRouterFactoryRawInstanceAddress, per_party_router_factory.ContractType, per_party_router_factory.Version, ""))

		for i, params := range input.Executors {
			executorInstanceID, err := ensureInstanceID(params.Template.InstanceId, "executor")
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure Executor #%d instance ID: %w", i, err)
			}
			executorTemplate := params.Template
			executorTemplate.InstanceId = types.TEXT(executorInstanceID)

			deployExecutorReport, err := operations.ExecuteOperation(b, factoryops.DeployExecutor, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployExecutor{
				Contract: executorTemplate,
			}, input.ProposalDriven))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy Executor #%d from factory: %w", i, err)
			}
			proposalOutputs = appendExerciseOutput(proposalOutputs, deployExecutorReport.Output, input.ProposalDriven)
			executorRawInstanceAddress := executorInstanceID.RawInstanceAddress(params.Template.Owner)
			addresses = append(addresses, newAddressRef(deps.ChainSelector(), executorRawInstanceAddress, executor.ContractType, executor.Version, params.Qualifier))
		}

		out := sequences.OnChainOutput{Addresses: addresses}
		if input.ProposalDriven {
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to build proposal batch for factory deploys: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				out.BatchOps = []mcms_types.BatchOperation{batchOp}
			}
		}

		return out, nil
	},
)

// DeployRMNFromFactory deploys RMNRemote via the rmn-qualified CCIPFactory.
var DeployRMNFromFactory = operations.NewSequence(
	"canton/ccip/deploy_rmn_from_factory",
	semver.MustParse("0.1.0"),
	"Deploys RMNRemote on Canton through the ccip CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		_, ccipOwnerParty, err := requireOwnerParties(input)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		factoryRawInstanceAddress, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		_, rmnAddressRef, proposalOutputs, err := deployRMNRemoteFromFactory(
			b, deps, input, factoryRawInstanceAddress, ccipOwnerParty,
		)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		out := sequences.OnChainOutput{Addresses: []datastore.AddressRef{*rmnAddressRef}}
		if input.ProposalDriven {
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to build RMN proposal batch: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				out.BatchOps = []mcms_types.BatchOperation{batchOp}
			}
		}

		return out, nil
	},
)

// DeployCCIPChainContractsFromFactory deploys CCIP contracts (no RMNRemote, no CommitteeVerifiers) via the ccip-qualified factory.
var DeployCCIPChainContractsFromFactory = operations.NewSequence(
	"canton/ccip/deploy_ccip_chain_contracts_from_factory",
	semver.MustParse("2.0.0"),
	"Deploys core CCIP contracts on Canton through CCIPFactory (excludes RMNRemote and CommitteeVerifier)",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		if input.RmnRemoteRawInstanceAddress == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("RmnRemoteRawInstanceAddress is required for core CCIP factory deploy")
		}

		coreInput := input
		coreInput.CommitteeVerifiers = nil

		report, err := operations.ExecuteSequence(b, DeployChainContractsFromFactory, deps, coreInput)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		return report.Output, nil
	},
)

// DeployCCVFromFactory deploys CommitteeVerifier contracts via a CCV-qualified CCIPFactory.
var DeployCCVFromFactory = operations.NewSequence(
	"canton/ccip/deploy_ccv_from_factory",
	semver.MustParse("0.1.0"),
	"Deploys CommitteeVerifier contracts on Canton through CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployChainContractsParams) (sequences.OnChainOutput, error) {
		if len(input.CommitteeVerifiers) == 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("at least one committee verifier is required")
		}
		if input.RmnRemoteRawInstanceAddress == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("RmnRemoteRawInstanceAddress is required for CCV factory deploy")
		}

		_, ccipOwnerParty, err := requireOwnerParties(input)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}
		ccvOwnerParty, err := requireCCVOwnerParty(input)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		factoryRawInstanceAddress, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		var addresses []datastore.AddressRef
		var proposalOutputs []contract.ExerciseOutput

		for i, committeeVerifierParams := range input.CommitteeVerifiers {
			qualifier := committeeVerifierParams.Qualifier
			committeeVerifierOwner := ccvOwnerParty
			committeeVerifierInstanceID, err := ensureInstanceID(committeeVerifierParams.Template.InstanceId, "committeeverifier")
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to ensure CommitteeVerifier #%d instance ID: %w", i, err)
			}
			committeeVerifierTemplate := ccvsbindings.CommitteeVerifier{
				InstanceId:                   types.TEXT(committeeVerifierInstanceID),
				Owner:                        committeeVerifierOwner,
				CcipOwner:                    ccipOwnerParty,
				VersionTag:                   committeeVerifierParams.Template.VersionTag,
				AllowListAdmin:               committeeVerifierParams.Template.AllowListAdmin,
				MessageSentObservers:         committeeVerifierParams.Template.MessageSentObservers,
				StorageLocations:             committeeVerifierParams.Template.StorageLocations,
				StorageLocationsAdmin:        committeeVerifierParams.Template.StorageLocationsAdmin,
				PendingStorageLocationsAdmin: committeeVerifierParams.Template.PendingStorageLocationsAdmin,
				RemoteChainConfigs:           committeeVerifierParams.Template.RemoteChainConfigs,
				SignerConfigs:                committeeVerifierParams.Template.SignerConfigs,
				Deps: ccvsbindings.CommitteeVerifierDeps{
					RmnRemote: input.RmnRemoteRawInstanceAddress.Binding(),
				},
			}
			deployReport, err := operations.ExecuteOperation(b, factoryops.DeployCommitteeVerifier, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployCommitteeVerifier{Contract: committeeVerifierTemplate}, input.ProposalDriven))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier #%d from factory: %w", i, err)
			}
			proposalOutputs = appendExerciseOutput(proposalOutputs, deployReport.Output, input.ProposalDriven)
			committeeVerifierRawInstanceAddress := committeeVerifierInstanceID.RawInstanceAddress(committeeVerifierOwner)
			addresses = append(addresses, newAddressRef(deps.ChainSelector(), committeeVerifierRawInstanceAddress, committee_verifier.ContractType, committee_verifier.Version, qualifier))
		}

		out := sequences.OnChainOutput{Addresses: addresses}
		if input.ProposalDriven {
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to build CCV proposal batch: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				out.BatchOps = []mcms_types.BatchOperation{batchOp}
			}
		}

		return out, nil
	},
)

func requireOwnerParties(input DeployChainContractsParams) (types.PARTY, types.PARTY, error) {
	if input.OwnerParty == "" {
		return "", "", fmt.Errorf("OwnerParty is required")
	}
	if input.CCIPOwnerParty == "" {
		return "", "", fmt.Errorf("CCIPOwnerParty is required")
	}

	return types.PARTY(input.OwnerParty), types.PARTY(input.CCIPOwnerParty), nil
}

func requireCCVOwnerParty(input DeployChainContractsParams) (types.PARTY, error) {
	if input.CCVOwnerParty == "" {
		return "", fmt.Errorf("CCVOwnerParty is required")
	}

	return types.PARTY(input.CCVOwnerParty), nil
}

func resolveOrDeployRMNRemote(
	b operations.Bundle,
	deps canton.Chain,
	input DeployChainContractsParams,
	ccipOwnerParty types.PARTY,
) (contracts.RawInstanceAddress, *datastore.AddressRef, []contract.ExerciseOutput, error) {
	if input.DevenvBundledDeploy {
		if input.RmnRemoteRawInstanceAddress != "" {
			return "", nil, nil, fmt.Errorf("RmnRemoteRawInstanceAddress must not be set when DevenvBundledDeploy is true")
		}
		rmnFactoryRaw, err := rawInstanceAddressFromAddressRef(input.RMNFactoryAddressRef)
		if err != nil {
			return "", nil, nil, fmt.Errorf("RMNFactoryAddressRef is required for DevenvBundledDeploy: %w", err)
		}

		return deployRMNRemoteFromFactory(b, deps, input, rmnFactoryRaw, ccipOwnerParty)
	}
	if input.RmnRemoteRawInstanceAddress == "" {
		return "", nil, nil, fmt.Errorf("RmnRemoteRawInstanceAddress is required")
	}

	return input.RmnRemoteRawInstanceAddress, nil, nil, nil
}

func deployRMNRemoteFromFactory(
	b operations.Bundle,
	deps canton.Chain,
	input DeployChainContractsParams,
	factoryRawInstanceAddress contracts.RawInstanceAddress,
	ccipOwnerParty types.PARTY,
) (contracts.RawInstanceAddress, *datastore.AddressRef, []contract.ExerciseOutput, error) {
	rmnInstanceID, err := ensureInstanceID(input.RMNRemote.Template.InstanceId, "rmn_remote")
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to ensure RMNRemote instance ID: %w", err)
	}
	rmnTemplate := rmn.RMNRemote{
		InstanceId:      types.TEXT(rmnInstanceID),
		RmnOwner:        ccipOwnerParty,
		CcipOwner:       ccipOwnerParty,
		CustomObservers: input.RMNRemote.Template.CustomObservers,
		CursedSubjects:  input.RMNRemote.Template.CursedSubjects,
	}
	deployRMNRemoteReport, err := operations.ExecuteOperation(b, factoryops.DeployRMNRemote, deps, newChoiceInput(factoryRawInstanceAddress, factorybindings.DeployRMNRemote{Contract: rmnTemplate}, input.ProposalDriven))
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to deploy RMNRemote from factory: %w", err)
	}

	rmnRemoteRawInstanceAddress := rmnInstanceID.RawInstanceAddress(ccipOwnerParty)
	addressRef := newAddressRef(deps.ChainSelector(), rmnRemoteRawInstanceAddress, rmn_remote.ContractType, rmn_remote.Version, "")

	var proposalOutputs []contract.ExerciseOutput
	if input.ProposalDriven {
		proposalOutputs = append(proposalOutputs, deployRMNRemoteReport.Output)
	}

	return rmnRemoteRawInstanceAddress, &addressRef, proposalOutputs, nil
}

func appendExerciseOutput(outputs []contract.ExerciseOutput, output contract.ExerciseOutput, proposalDriven bool) []contract.ExerciseOutput {
	if !proposalDriven {
		return outputs
	}

	return append(outputs, output)
}

func newChoiceInput[ARGS any](
	rawInstanceAddress contracts.RawInstanceAddress,
	args ARGS,
	proposalDriven bool,
) contract.ChoiceInput[ARGS] {
	return contract.ChoiceInput[ARGS]{
		InstanceAddress:    rawInstanceAddress.InstanceAddress(),
		RawInstanceAddress: rawInstanceAddress.String(),
		Args:               args,
		MCMSEnabled:        proposalDriven,
	}
}

func ensureInstanceID(current types.TEXT, prefix string) (contracts.InstanceID, error) {
	if current != "" {
		return contracts.InstanceID(current), nil
	}

	return contracts.NewInstanceID(prefix)
}

func newAddressRef(
	chainSelector uint64,
	rawInstanceAddress contracts.RawInstanceAddress,
	contractType deployment.ContractType,
	version *semver.Version,
	qualifier string,
) datastore.AddressRef {
	return datastore.AddressRef{
		Address:       rawInstanceAddress.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(rawInstanceAddress.String()),
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(contractType),
		Version:       version,
		Qualifier:     qualifier,
	}
}

func rawInstanceAddressFromAddressRef(ref datastore.AddressRef) (contracts.RawInstanceAddress, error) {
	labels := ref.Labels.List()
	if len(labels) == 0 {
		return "", fmt.Errorf("address ref for %s is missing raw instance address label", ref.Type)
	}

	rawInstanceAddress, err := contracts.RawInstanceAddressFromString(labels[0])
	if err != nil {
		return "", fmt.Errorf("parse raw instance address label %q: %w", labels[0], err)
	}

	return rawInstanceAddress, nil
}
