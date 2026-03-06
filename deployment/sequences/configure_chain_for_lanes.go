package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// TODO should align this with the EVM changesets if possible? Currently, these field are hardcoded
type ConfigureChainForLanesInput struct {
	// The selector of the chain being configured.
	ChainSelector uint64
	// The GlobalConfig address on the chain being configured.
	GlobalConfig contracts.InstanceAddress
	// The FeeQuoter address on the chain being configured.
	FeeQuoter contracts.InstanceAddress
	// The OnRamp address on the chain being configured.
	// Similarly, we assume that all connections will use the same OnRamp.
	OnRamp contracts.InstanceAddress
	// The OffRamp address on the chain being configured
	OffRamp contracts.InstanceAddress

	// The CommitteeVerifiers on the chain being configured.
	// There can be multiple committee verifiers on a chain, each controlled by a different entity.
	CommitteeVerifiers []adapters.CommitteeVerifierConfig[contracts.InstanceAddress]
	// The configuration for each remote chain that we want to connect to.
	RemoteChains map[uint64]adapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]
}

var ConfigureChainForLanes = operations.NewSequence(
	"canton/ccip/configure_chain_for_lanes",
	semver.MustParse("1.7.0"),
	"Configures a Canton chain as a source & destination for multiple remote chains",
	// TODO change deps to cldf_chain.BlockChains once clients are added
	func(b operations.Bundle, deps dependencies.CantonDeps, input ConfigureChainForLanesInput) (sequences.OnChainOutput, error) {

		// Create inputs for each operation
		globalConfigSourceChainConfigArgs := make([]common.UpdateSourceChainConfig, 0, len(input.RemoteChains))
		globalConfigDestChainConfigArgs := make([]common.UpdateDestChainConfig, 0, len(input.RemoteChains))

		for remoteSelector, remoteConfig := range input.RemoteChains {
			remoteSelectorStr := strconv.FormatUint(remoteSelector, 10)

			// Inbound / OffRamp
			defaultInboundCCVs := make([]common.RawInstanceAddress, 0, len(remoteConfig.DefaultInboundCCVs))
			for _, ccv := range remoteConfig.DefaultInboundCCVs {
				defaultInboundCCVs = append(defaultInboundCCVs, common.RawInstanceAddress{Unpack: types.TEXT(ccv)})
			}
			laneMandatedInboundCCVs := make([]common.RawInstanceAddress, 0, len(remoteConfig.LaneMandatedInboundCCVs))
			for _, ccv := range remoteConfig.LaneMandatedInboundCCVs {
				laneMandatedInboundCCVs = append(laneMandatedInboundCCVs, common.RawInstanceAddress{Unpack: types.TEXT(ccv)})
			}
			onRamps := make([]types.TEXT, 0, len(remoteConfig.OnRamps))
			for _, onRamp := range remoteConfig.OnRamps {
				onRamps = append(onRamps, types.TEXT(hex.EncodeToString(onRamp)))
			}
			globalConfigSourceChainConfigArgs = append(globalConfigSourceChainConfigArgs, common.UpdateSourceChainConfig{
				SourceChainSelector: types.NUMERIC(remoteSelectorStr),
				Config: common.SourceChainConfig{
					IsEnabled:        types.BOOL(remoteConfig.AllowTrafficFrom),
					OnRampAddress:    onRamps[0], // TODO: currently only supports one onRamp
					LaneMandatedCCVs: laneMandatedInboundCCVs,
					DefaultCCVs:      defaultInboundCCVs,
				},
			})

			defaultOutboundCCVs := make([]common.RawInstanceAddress, 0, len(remoteConfig.DefaultOutboundCCVs))
			for _, ccv := range remoteConfig.DefaultOutboundCCVs {
				defaultOutboundCCVs = append(defaultOutboundCCVs, common.RawInstanceAddress{Unpack: types.TEXT(ccv)})
			}
			laneMandatedOutboundCCVs := make([]common.RawInstanceAddress, 0, len(remoteConfig.LaneMandatedOutboundCCVs))
			for _, ccv := range remoteConfig.LaneMandatedOutboundCCVs {
				laneMandatedOutboundCCVs = append(laneMandatedOutboundCCVs, common.RawInstanceAddress{Unpack: types.TEXT(ccv)})
			}

			// Outbound / OnRamp
			globalConfigDestChainConfigArgs = append(globalConfigDestChainConfigArgs, common.UpdateDestChainConfig{
				DestChainSelector: types.NUMERIC(remoteSelectorStr),
				Config: common.DestChainConfig{
					IsEnabled:                 types.BOOL(remoteConfig.AllowTrafficFrom),
					OffRampAddress:            types.TEXT(hex.EncodeToString(remoteConfig.OffRamp)), // Remote chain off-ramp for outbound execution
					LaneMandatedCCVs:          laneMandatedOutboundCCVs,
					DefaultCCVs:               defaultOutboundCCVs,
					MessageNetworkFeeUSDCents: types.NUMERIC(strconv.FormatInt(int64(remoteConfig.FeeQuoterDestChainConfig.NetworkFeeUSDCents), 10)),
					TokenNetworkFeeUSDCents:   types.NUMERIC(strconv.FormatInt(int64(remoteConfig.FeeQuoterDestChainConfig.DefaultTokenFeeUSDCents), 10)), // TODO: check if this is accurate
				},
			})
			// TODO: Other configs once the contracts are ready
		}

		// Apply SourceChainConfigs to GlobalConfig
		for i, arg := range globalConfigSourceChainConfigArgs {
			_, err := operations.ExecuteOperation(b, global_config.UpdateSourceChainConfig, deps, contract.ChoiceInput[common.UpdateSourceChainConfig]{
				ChainSelector:   deps.Chain.Selector,
				InstanceAddress: input.GlobalConfig,
				ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
				Args:            arg,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply source chain config %d for remote chain %s: %w", i, string(arg.SourceChainSelector), err)
			}
		}

		// Apply signature configs to CommitteeVerifiers
		for _, verifierConfig := range input.CommitteeVerifiers {
			signatureConfigs := make([]ccvs.SignatureConfig, 0, len(verifierConfig.RemoteChains))
			for remoteSelector, remoteConfig := range verifierConfig.RemoteChains {
				signerKeys := make([]types.TEXT, len(remoteConfig.SignatureConfig.Signers))
				for i, signer := range remoteConfig.SignatureConfig.Signers {
					// Decode and encode signer pubkeys to ensure they're in the correct format
					signerBytes, err := hex.DecodeString(strings.TrimPrefix(signer, "0x"))
					if err != nil {
						return sequences.OnChainOutput{}, fmt.Errorf("failed to decode signer key %d for remote chain %d: %w", i, remoteSelector, err)
					}
					signerKeys[i] = types.TEXT(hex.EncodeToString(signerBytes))
				}

				signatureConfigs = append(signatureConfigs, ccvs.SignatureConfig{
					SourceChainSelector: types.NUMERIC(strconv.FormatUint(remoteSelector, 10)),
					Threshold:           types.INT64(remoteConfig.SignatureConfig.Threshold),
					SignerKeys:          signerKeys,
				})
			}
			for _, address := range verifierConfig.CommitteeVerifier {
				_, err := operations.ExecuteOperation(b, committee_verifier.ApplySignatureConfigs, deps, contract.ChoiceInput[ccvs.CommitteeVerifierApplySignatureConfigs]{
					ChainSelector:   deps.Chain.Selector,
					InstanceAddress: address,
					ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
					Args: ccvs.CommitteeVerifierApplySignatureConfigs{
						SourceChainSelectorsToRemove: nil, // This doesn't support removing chains
						SignatureConfigs:             signatureConfigs,
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to apply signature configs to CommitteeVerifier at address %s: %w", address.Hex(), err)
				}
			}
		}

		// Apply DestChainConfigs to GlobalConfig
		for i, arg := range globalConfigDestChainConfigArgs {
			_, err := operations.ExecuteOperation(b, global_config.UpdateDestChainConfig, deps, contract.ChoiceInput[common.UpdateDestChainConfig]{
				ChainSelector:   deps.Chain.Selector,
				InstanceAddress: input.GlobalConfig,
				ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
				Args:            arg,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to apply source chain config %d for remote chain %s: %w", i, string(arg.DestChainSelector), err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)
