package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	seqcore "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

var _ ccipadapters.DeployChainContractsAdapter = (*CantonDeployChainContractsAdapter)(nil)

type CantonDeployChainContractsAdapter struct{}

func (c *CantonDeployChainContractsAdapter) GetDefaultDeployContractParams(_ uint64) ccipadapters.DeployContractParams {
	return ccipadapters.DeployContractParams{
		Executors: []ccipadapters.ExecutorDeployParams{defaultExecutorDeployParams()},
	}
}

func (c *CantonDeployChainContractsAdapter) ResolveDeployAddresses(
	_ deployment.Environment,
	_ uint64,
) (ccipadapters.DeployChainResolvedAddresses, error) {
	// Canton uses the participant party as deployer; no prerequisite contracts to resolve.
	return ccipadapters.DeployChainResolvedAddresses{}, nil
}

func (c *CantonDeployChainContractsAdapter) BuildDeployContractParams(
	input ccipadapters.BuildDeployContractParamsInput,
) (ccipadapters.DeployContractParams, error) {
	if len(input.CommitteeVerifiers) == 0 {
		return ccipadapters.DeployContractParams{}, fmt.Errorf("chain %d: at least one committee verifier is required", input.ChainSelector)
	}

	params := input.Defaults
	params.CommitteeVerifiers = input.CommitteeVerifiers

	if len(params.Executors) == 0 {
		exec := defaultExecutorDeployParams()
		exec.DynamicConfig.FeeAggregator = input.CommitteeVerifiers[0].FeeAggregator
		params.Executors = []ccipadapters.ExecutorDeployParams{exec}
	}

	return ccipadapters.ApplyDeployContractParamsOverrides(params, input.Overrides), nil
}

func (c *CantonDeployChainContractsAdapter) SetContractParamsFromImportedConfig() *cldf_ops.Sequence[ccipadapters.DeployChainConfigCreatorInput, ccipadapters.DeployContractParams, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton/set-contract-params-from-imported-config",
		semver.MustParse("2.0.0"),
		"Returns the user-provided contract params for Canton",
		func(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, input ccipadapters.DeployChainConfigCreatorInput) (ccipadapters.DeployContractParams, error) {
			return input.UserProvidedConfig, nil
		},
	)
}

func (c *CantonDeployChainContractsAdapter) DeployChainContracts() *cldf_ops.Sequence[ccipadapters.DeployChainContractsInput, ccipadapters.DeployChainContractsOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton/deploy-chain-contracts",
		semver.MustParse("2.0.0"),
		"Deploys CCIP contracts on a Canton chain",
		func(bundle cldf_ops.Bundle, chains cldf_chain.BlockChains, input ccipadapters.DeployChainContractsInput) (ccipadapters.DeployChainContractsOutput, error) {
			return DeployCantonChainContracts(bundle.GetContext(), bundle, chains, input)
		},
	)
}

func DeployCantonChainContracts(ctx context.Context, bundle cldf_ops.Bundle, chains cldf_chain.BlockChains, input ccipadapters.DeployChainContractsInput) (ccipadapters.DeployChainContractsOutput, error) {
	chain, ok := chains.CantonChains()[input.ChainSelector]
	if !ok || len(chain.Participants) == 0 {
		return ccipadapters.DeployChainContractsOutput{}, fmt.Errorf("canton chain %d not found or has no participants", input.ChainSelector)
	}

	participant := chain.Participants[0]
	ownerParty := deployerPartyID(input.DeployerContract, participant)

	nativeInstrumentID, err := lookupNativeInstrumentID(ctx, participant, nil, input.ChainSelector)
	if err != nil {
		return ccipadapters.DeployChainContractsOutput{}, err
	}

	factoryAddressRef, err := dsutils.FactoryAddressRefFromRefs(input.ChainSelector, dsutils.QualifierCCIP, input.ExistingAddresses)
	if err != nil {
		return ccipadapters.DeployChainContractsOutput{}, err
	}
	rmnFactoryAddressRef, err := dsutils.FactoryAddressRefFromRefs(input.ChainSelector, dsutils.QualifierRMN, input.ExistingAddresses)
	if err != nil {
		return ccipadapters.DeployChainContractsOutput{}, err
	}

	proposalDriven := shouldUseMCMSProposalDeployment(input, chain)

	out, err := cldf_ops.ExecuteSequence(bundle, sequences.DeployChainContractsFromFactory, chain, sequences.DeployChainContractsParams{
		OwnerParty:           ownerParty,
		CCIPOwnerParty:       ownerParty,
		CCVOwnerParty:        ownerParty,
		RMNOwnerParty:        ownerParty,
		DevenvBundledDeploy:  true,
		FactoryAddressRef:    factoryAddressRef,
		RMNFactoryAddressRef: rmnFactoryAddressRef,
		CommitteeVerifiers:   committeeVerifierParams(ownerParty, input.ContractParams.CommitteeVerifiers),
		GlobalConfig: sequences.GlobalConfigParams{
			Template: core.GlobalConfig{
				CcipOwner:     "",
				ChainSelector: types.NUMERIC(strconv.FormatUint(input.ChainSelector, 10)),
			},
		},
		ProposalDriven:     proposalDriven,
		NativeInstrumentId: nativeInstrumentID,
		FeeQuoterConfig: sequences.FeeQuoterParams{
			Template: core.FeeQuoter{
				PriceUpdaters: []types.PARTY{types.PARTY(ownerParty)},
			},
		},
		RMNRemote: sequences.RMNRemoteParams{
			Template: core.RMNRemote{
				CursedSubjects: nil,
			},
		},
		Executors: executorParams(ownerParty, input.ContractParams.Executors),
	})
	if err != nil {
		return ccipadapters.DeployChainContractsOutput{}, fmt.Errorf("failed to deploy canton chain contracts for selector %d: %w", input.ChainSelector, err)
	}

	return ccipadapters.DeployChainContractsOutput{
		OnChainOutput: seqcore.OnChainOutput{
			Addresses: out.Output.Addresses,
			BatchOps:  out.Output.BatchOps,
		},
		RefsToTransferOwnership: nil,
	}, nil
}

func shouldUseMCMSProposalDeployment(input ccipadapters.DeployChainContractsInput, chain canton.Chain) bool {
	if input.DeployerKeyOwned {
		return false
	}
	if len(chain.Participants) == 0 {
		return false
	}

	return len(chain.Participants[0].ReadAsPartyIDs) > 0
}

func deployerPartyID(deployerContract string, participant canton.Participant) string {
	if trimmed := strings.TrimPrefix(deployerContract, "canton:"); trimmed != "" && trimmed != deployerContract {
		return trimmed
	}
	if deployerContract != "" {
		return deployerContract
	}

	return participant.PartyID
}

func committeeVerifierParams(ownerParty string, verifiers []ccipadapters.CommitteeVerifierDeployParams) []sequences.CommitteeVerifierParams {
	params := make([]sequences.CommitteeVerifierParams, 0, len(verifiers))
	for _, verifier := range verifiers {
		storageLocations := make([]types.TEXT, len(verifier.StorageLocations))
		for i, location := range verifier.StorageLocations {
			storageLocations[i] = types.TEXT(location)
		}
		qualifier := verifier.Qualifier
		if qualifier == "" {
			qualifier = devenvcommon.DefaultCommitteeVerifierQualifier
		}
		params = append(params, sequences.CommitteeVerifierParams{
			Qualifier: qualifier,
			Template: committeeverifier.CommitteeVerifier{
				CcipOwner:                    types.PARTY(ownerParty),
				VersionTag:                   types.TEXT("e9a05a20"),
				MessageSentObservers:         nil,
				StorageLocations:             storageLocations,
				StorageLocationsAdmin:        types.PARTY(ownerParty),
				PendingStorageLocationsAdmin: types.PARTY(ownerParty),
				Deps:                         committeeverifier.CommitteeVerifierDeps{},
			},
		})
	}

	return params
}

func executorParams(
	ownerParty string,
	executors []ccipadapters.ExecutorDeployParams,
) []sequences.ExecutorParams {
	if len(executors) == 0 {
		return []sequences.ExecutorParams{defaultExecutorParams(ownerParty)}
	}

	params := make([]sequences.ExecutorParams, 0, len(executors))
	for _, exec := range executors {
		qualifier := exec.Qualifier
		if qualifier == "" {
			qualifier = devenvcommon.DefaultExecutorQualifier
		}
		maxCCVs := exec.MaxCCVsPerMsg
		if maxCCVs == 0 {
			maxCCVs = 10
		}
		params = append(params, sequences.ExecutorParams{
			Qualifier: qualifier,
			Template: executor.Executor{
				Owner:         types.PARTY(ownerParty),
				MaxCCVsPerMsg: types.INT64(maxCCVs),
				DynamicConfig: executor.DynamicConfig{
					FeeAggregator:         nil,
					AllowedFinalityConfig: requestedFinality(exec.DynamicConfig.AllowedFinalityConfig),
					CcvAllowlistEnabled:   types.BOOL(exec.DynamicConfig.CcvAllowlistEnabled),
				},
				AllowedCCVs: nil,
			},
		})
	}

	return params
}

func defaultExecutorDeployParams() ccipadapters.ExecutorDeployParams {
	return ccipadapters.ExecutorDeployParams{
		Qualifier:     devenvcommon.DefaultExecutorQualifier,
		MaxCCVsPerMsg: 10,
		DynamicConfig: ccipadapters.ExecutorDynamicDeployConfig{
			AllowedFinalityConfig: finality.Config{WaitForFinality: true},
			CcvAllowlistEnabled:   false,
		},
	}
}

func defaultExecutorParams(ownerParty string) sequences.ExecutorParams {
	return sequences.ExecutorParams{
		Qualifier: devenvcommon.DefaultExecutorQualifier,
		Template: executor.Executor{
			Owner:         types.PARTY(ownerParty),
			MaxCCVsPerMsg: types.INT64(10),
			DynamicConfig: executor.DynamicConfig{
				FeeAggregator:         nil,
				AllowedFinalityConfig: waitForFinalityRequested(),
				CcvAllowlistEnabled:   types.BOOL(false),
			},
			AllowedCCVs: nil,
		},
	}
}

func requestedFinality(cfg finality.Config) ccipcodec.FinalityConfig {
	if cfg.IsZero() || cfg.WaitForFinality {
		return waitForFinalityRequested()
	}
	if cfg.WaitForSafe && cfg.BlockDepth > 0 {
		panic(fmt.Sprintf("unsupported combined finality config: waitForSafe=true blockDepth=%d", cfg.BlockDepth))
	}
	if cfg.WaitForSafe {
		return waitForSafeRequested()
	}
	if cfg.BlockDepth > 0 {
		return ccipcodec.FinalityConfig{BlockDepth: new(types.INT64(cfg.BlockDepth))}
	}

	return waitForFinalityRequested()
}

func waitForFinalityRequested() ccipcodec.FinalityConfig {
	return ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}}
}

func waitForSafeRequested() ccipcodec.FinalityConfig {
	return ccipcodec.FinalityConfig{WaitForSafe: &types.UNIT{}}
}
