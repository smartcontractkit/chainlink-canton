package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

func normalizeConfiguredOnRamp(onRamp []byte) (string, error) {
	raw := onRamp
	if len(onRamp) > 0 {
		maybeHex := strings.TrimPrefix(string(onRamp), "0x")
		if len(maybeHex)%2 == 0 && maybeHex != "" && strings.IndexFunc(maybeHex, func(r rune) bool {
			return !strings.ContainsRune("0123456789abcdefABCDEF", r)
		}) == -1 {
			decoded, err := hex.DecodeString(maybeHex)
			if err != nil {
				return "", fmt.Errorf("failed to decode onramp hex bytes %q: %w", string(onRamp), err)
			}
			raw = decoded
		}
	}
	if len(raw) > 32 {
		return "", fmt.Errorf("onramp address exceeds 32 bytes: %d", len(raw))
	}

	return hex.EncodeToString(gethcommon.LeftPadBytes(raw, 32)), nil
}

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
		globalConfigSourceChainConfigArgs := make([]common.SourceChainConfigArgs, 0, len(input.RemoteChains))
		globalConfigDestChainConfigArgs := make([]common.DestChainConfigArgs, 0, len(input.RemoteChains))
		feeQuoterDestChainConfigArgs := make([]feequoter.DestChainConfigArgs2, 0, len(input.RemoteChains))

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
				// EVM messages encode source addresses as 32-byte left-padded values, so we must
				// normalize configured remote onramps to that same shape or OffRamp validation will reject them.
				normalizedOnRamp, err := normalizeConfiguredOnRamp(onRamp)
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to normalize onramp for remote chain %d: %w", remoteSelector, err)
				}
				onRamps = append(onRamps, types.TEXT(normalizedOnRamp))
			}
			globalConfigSourceChainConfigArgs = append(globalConfigSourceChainConfigArgs, common.SourceChainConfigArgs{
				SourceChainSelector: types.NUMERIC(remoteSelectorStr),
				IsEnabled:           types.BOOL(remoteConfig.AllowTrafficFrom),
				OnRampAddresses:     onRamps,
				LaneMandatedCCVs:    laneMandatedInboundCCVs,
				DefaultCCVs:         defaultInboundCCVs,
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
			globalConfigDestChainConfigArgs = append(globalConfigDestChainConfigArgs, common.DestChainConfigArgs{
				DestChainSelector:         types.NUMERIC(remoteSelectorStr),
				IsEnabled:                 types.BOOL(remoteConfig.AllowTrafficFrom),
				AddressBytesLength:        types.INT64(remoteConfig.AddressBytesLength),
				OffRampAddress:            types.TEXT(hex.EncodeToString(remoteConfig.OffRamp)), // Remote chain off-ramp for outbound execution
				LaneMandatedCCVs:          laneMandatedOutboundCCVs,
				DefaultCCVs:               defaultOutboundCCVs,
				MessageNetworkFeeUSDCents: types.NUMERIC(strconv.FormatInt(int64(remoteConfig.FeeQuoterDestChainConfig.NetworkFeeUSDCents), 10)),
				TokenNetworkFeeUSDCents:   types.NUMERIC(strconv.FormatInt(int64(remoteConfig.FeeQuoterDestChainConfig.DefaultTokenFeeUSDCents), 10)), // TODO: check if this is accurate
			})

			fqConfig := remoteConfig.FeeQuoterDestChainConfig
			feeQuoterDestChainConfigArgs = append(feeQuoterDestChainConfigArgs, feequoter.DestChainConfigArgs2{
				DestChainSelector: types.NUMERIC(remoteSelectorStr),
				DestChainConfig: feequoter.DestChainConfig2{
					IsEnabled:                   types.BOOL(fqConfig.IsEnabled),
					MaxDataBytes:                types.INT64(fqConfig.MaxDataBytes),
					MaxPerMsgGasLimit:           types.INT64(fqConfig.MaxPerMsgGasLimit),
					DestGasOverhead:             types.INT64(fqConfig.DestGasOverhead),
					DestGasPerPayloadByteBase:   types.INT64(fqConfig.DestGasPerPayloadByteBase),
					DefaultTxGasLimit:           types.INT64(fqConfig.DefaultTxGasLimit),
					DefaultTokenFeeUSD:          types.NUMERIC(strconv.FormatUint(uint64(fqConfig.DefaultTokenFeeUSDCents), 10)),
					DefaultTokenDestGasOverhead: types.INT64(fqConfig.DefaultTokenDestGasOverhead),
				},
			})
		}

		// Apply SourceChainConfigs to GlobalConfig
		_, err := operations.ExecuteOperation(b, global_config.ApplySourceChainConfigUpdates, deps, contract.ChoiceInput[common.ApplySourceChainConfigUpdates]{
			ChainSelector:   deps.Chain.Selector,
			InstanceAddress: input.GlobalConfig,
			ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
			Args: common.ApplySourceChainConfigUpdates{
				SourceChainConfigUpdates: globalConfigSourceChainConfigArgs,
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to apply source chain config updates: %w", err)
		}

		// Apply signature configs to CommitteeVerifiers
		for _, verifierConfig := range input.CommitteeVerifiers {
			_, err := operations.ExecuteSequence(b, ConfigureCommitteeVerifierForLanes, deps, ConfigureCommitteeVerifierForLanesInput{
				ChainSelector:           input.ChainSelector,
				CommitteeVerifierConfig: verifierConfig,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to configure committee verifier for lanes: %w", err)
			}
		}

		// Apply DestChainConfigs to GlobalConfig
		_, err = operations.ExecuteOperation(b, global_config.ApplyDestChainConfigUpdates, deps, contract.ChoiceInput[common.ApplyDestChainConfigUpdates]{
			ChainSelector:   deps.Chain.Selector,
			InstanceAddress: input.GlobalConfig,
			ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
			Args: common.ApplyDestChainConfigUpdates{
				DestChainConfigUpdates: globalConfigDestChainConfigArgs,
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to apply dest chain config updates: %w", err)
		}

		// Apply DestChainConfigs to FeeQuoter
		_, err = operations.ExecuteOperation(b, feequoterop.ApplyDestChainConfigUpdates, deps, contract.ChoiceInput[feequoter.ApplyDestChainConfigUpdates2]{
			ChainSelector:   deps.Chain.Selector,
			InstanceAddress: input.FeeQuoter,
			ActAs:           []string{deps.Chain.Participants[deps.Participant].PartyID},
			Args: feequoter.ApplyDestChainConfigUpdates2{
				DestChainConfigArgs: feeQuoterDestChainConfigArgs,
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to apply fee quoter dest chain config updates: %w", err)
		}

		return sequences.OnChainOutput{}, nil
	},
)
