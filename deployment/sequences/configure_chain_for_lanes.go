package sequences

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	executor2 "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ConfigureLaneLegAsSource = operations.NewSequence(
	"CantonConfigureLaneLegAsSource",
	semver.MustParse("2.0.0"),
	"Configures a lane leg as source on CCIP 2.0.0",
	func(b operations.Bundle, deps chain.BlockChains, input lanes.UpdateLanesInput) (output sequences.OnChainOutput, err error) {
		b.Logger.Infof("Canton Configuring lane leg as source. src: %+v, dest: %+v", input.Source, input.Dest)

		chain, ok := deps.CantonChains()[input.Source.Selector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.Source.Selector)
		}

		sourceChain := input.Source
		destChain := input.Dest

		// GlobalConfig - Dest Chain Config
		globalConfigAddress := contracts.HexToInstanceAddress(sourceChain.CantonLaneConfig.GlobalConfig.Address)
		isEnabled := len(destChain.Router) > 0
		defaultExecutor, err := dsutils.GetRawInstanceAddressFromAddressRef(sourceChain.DefaultExecutor)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("getting default executor: %w", err)
		}
		laneMandatedOutboundCCVs := make([]mcms.RawInstanceAddress, 0, len(sourceChain.LaneMandatedOutboundCCVs))
		for _, ccv := range sourceChain.LaneMandatedOutboundCCVs {
			outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated outbound CCV: %w", err)
			}
			laneMandatedOutboundCCVs = append(laneMandatedOutboundCCVs, outboundCCV.Binding())
		}
		defaultOutboundCCVs := make([]mcms.RawInstanceAddress, 0, len(sourceChain.LaneMandatedOutboundCCVs))
		for _, ccv := range sourceChain.LaneMandatedOutboundCCVs {
			outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated outbound CCV: %w", err)
			}
			defaultOutboundCCVs = append(defaultOutboundCCVs, outboundCCV.Binding())
		}
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("getting global config dest chain config args: %w", err)
		}
		_, err = operations.ExecuteOperation(b, global_config.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[common.ApplyDestChainConfigUpdates]{
			InstanceAddress: globalConfigAddress,
			Args: common.ApplyDestChainConfigUpdates{
				DestChainConfigUpdates: []common.DestChainConfigArgs{
					common.DestChainConfigArgs{
						DestChainSelector:         types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
						IsEnabled:                 types.BOOL(isEnabled),
						AddressBytesLength:        types.INT64(destChain.AddressBytesLength),
						TokenReceiverAllowed:      true, // TODO: missing from input
						BaseExecutionGasCost:      types.INT64(destChain.BaseExecutionGasCost),
						OffRampAddress:            types.TEXT(hex.EncodeToString(destChain.OffRamp)),
						DefaultExecutor:           defaultExecutor.Binding(),
						LaneMandatedCCVs:          laneMandatedOutboundCCVs,
						DefaultCCVs:               defaultOutboundCCVs,
						MessageNetworkFeeUSDCents: types.NUMERIC(strconv.FormatUint(uint64(destChain.MessageNetworkFeeUSDCents), 10)),
						TokenNetworkFeeUSDCents:   types.NUMERIC(strconv.FormatUint(uint64(destChain.TokenNetworkFeeUSDCents), 10)),
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("applying dest chain config updates to global config: %w", err)
		}

		// Executor - Dest Chain Config
		executorAddress := contracts.HexToInstanceAddress(sourceChain.DefaultExecutor.Address)
		_, err = operations.ExecuteOperation(b, executor2.ApplyDestChainUpdates, chain, contract.ChoiceInput[executor.ApplyDestChainUpdates]{
			InstanceAddress: executorAddress,
			Args: executor.ApplyDestChainUpdates{
				DestChainSelectorsToRemove: nil,
				DestChainSelectorsToAdd: []executor.RemoteChainConfigArgs{
					{
						DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
						Config: executor.RemoteChainConfig{
							FeeUSDCents: types.NUMERIC(strconv.FormatUint(uint64(destChain.ExecutorDestChainConfig.USDCentsFee), 10)),
							Enabled:     types.BOOL(destChain.ExecutorDestChainConfig.Enabled),
						},
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("applying dest chain config updates to executor: %w", err)
		}

		// FeeQuoter - Dest Chain Config
		feeQuoterAddress := contracts.BytesToInstanceAddress(sourceChain.FeeQuoter)
		_, err = operations.ExecuteOperation(b, feequoterop.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[feequoter.ApplyDestChainConfigUpdates2]{
			InstanceAddress: feeQuoterAddress,
			Args: feequoter.ApplyDestChainConfigUpdates2{
				DestChainConfigArgs: []feequoter.DestChainConfigArgs2{
					feequoter.DestChainConfigArgs2{
						DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
						DestChainConfig: feequoter.DestChainConfig2{
							IsEnabled:                   types.BOOL(destChain.FeeQuoterDestChainConfig.IsEnabled),
							MaxDataBytes:                types.INT64(destChain.FeeQuoterDestChainConfig.MaxDataBytes),
							MaxPerMsgGasLimit:           types.INT64(destChain.FeeQuoterDestChainConfig.MaxPerMsgGasLimit),
							DestGasOverhead:             types.INT64(destChain.FeeQuoterDestChainConfig.DestGasOverhead),
							DestGasPerPayloadByteBase:   types.INT64(destChain.FeeQuoterDestChainConfig.DestGasPerPayloadByteBase),
							DefaultTxGasLimit:           types.INT64(destChain.FeeQuoterDestChainConfig.DefaultTxGasLimit),
							LinkFeeMultiplierPercent:    types.NUMERIC(strconv.FormatUint(uint64(destChain.FeeQuoterDestChainConfig.V2Params.LinkFeeMultiplierPercent), 10)),
							DefaultTokenFeeUSD:          types.NUMERIC(strconv.FormatUint(uint64(destChain.FeeQuoterDestChainConfig.DefaultTokenFeeUSDCents), 10)),
							DefaultTokenDestGasOverhead: types.INT64(destChain.FeeQuoterDestChainConfig.DefaultTokenDestGasOverhead),
						},
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("applying dest chain config updates to fee quoter: %w", err)
		}

		// FeeQuoter - Update prices (gas prices only, as these are per dest chain)
		_, err = operations.ExecuteOperation(b, feequoterop.UpdatePrices, chain, contract.ChoiceInput[feequoter.UpdatePrices]{
			InstanceAddress: feeQuoterAddress,
			Args: feequoter.UpdatePrices{
				PriceUpdates: feequoter.PriceUpdates{
					TokenPriceUpdates: nil,
					GasPriceUpdates: []feequoter.GasPriceUpdate{
						{
							DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
							UsdPerUnitGas:     types.NUMERIC(destChain.FeeQuoterDestChainConfig.V2Params.USDPerUnitGas.String()),
						},
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("updating gas prices in fee quoter: %w", err)
		}

		// CommitteeVerifier - Dest Chain Config
		for _, verifierConfig := range input.Source.CommitteeVerifiers {
			_, err = operations.ExecuteSequence(b, ConfigureCommitteeVerifierAsSource, deps, ConfigureCommitteeVerifierAsSourceInput{
				ChainSelector:           chain.Selector,
				CommitteeVerifierConfig: verifierConfig,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("configuring committee verifier as source: %w", err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)

var ConfigureLaneLegAsDest = operations.NewSequence(
	"CantonConfigureLaneLegAsDest",
	semver.MustParse("2.0.0"),
	"Configures a lane lad as dest on CCIP 2.0.0",
	func(b operations.Bundle, deps chain.BlockChains, input lanes.UpdateLanesInput) (output sequences.OnChainOutput, err error) {
		b.Logger.Infof("Canton Configuring lane leg as source. src: %+v, dest: %+v", input.Source, input.Dest)

		chain, ok := deps.CantonChains()[input.Dest.Selector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.Source.Selector)
		}

		sourceChain := input.Source
		destChain := input.Dest

		// GlobalConfig - Source Chain Config
		globalConfigAddress := contracts.HexToInstanceAddress(destChain.CantonLaneConfig.GlobalConfig.Address)
		laneMandatedInboundCCVs := make([]mcms.RawInstanceAddress, 0, len(destChain.LaneMandatedInboundCCVs))
		for _, ccv := range destChain.LaneMandatedInboundCCVs {
			inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated inbound CCV: %w", err)
			}
			laneMandatedInboundCCVs = append(laneMandatedInboundCCVs, inboundCCV.Binding())
		}
		defaultInboundCCVs := make([]mcms.RawInstanceAddress, 0, len(destChain.LaneMandatedInboundCCVs))
		for _, ccv := range destChain.LaneMandatedInboundCCVs {
			inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated inbound CCV: %w", err)
			}
			defaultInboundCCVs = append(defaultInboundCCVs, inboundCCV.Binding())
		}
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("getting global config dest chain config args: %w", err)
		}
		// TODO input doesn't support multiple OnRamps
		// TODO Must pad EVM addresses to 32 bytes
		inboundOnRampAddresses := []types.TEXT{
			types.TEXT(hex.EncodeToString(gethcommon.LeftPadBytes(sourceChain.OnRamp, 32))),
		}
		_, err = operations.ExecuteOperation(b, global_config.ApplySourceChainConfigUpdates, chain, contract.ChoiceInput[common.ApplySourceChainConfigUpdates]{
			InstanceAddress: globalConfigAddress,
			Args: common.ApplySourceChainConfigUpdates{
				SourceChainConfigUpdates: []common.SourceChainConfigArgs{
					{
						SourceChainSelector: types.NUMERIC(strconv.FormatUint(sourceChain.Selector, 10)),
						IsEnabled:           types.BOOL(!input.IsDisabled),
						OnRampAddresses:     inboundOnRampAddresses,
						DefaultCCVs:         defaultInboundCCVs,
						LaneMandatedCCVs:    laneMandatedInboundCCVs,
					},
				},
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("applying source chain config updates to global config: %w", err)
		}

		// CommitteeVerifier - Source Chain Config
		for _, verifierConfig := range input.Dest.CommitteeVerifiers {
			_, err = operations.ExecuteSequence(b, ConfigureCommitteeVerifierAsDest, deps, ConfigureCommitteeVerifierAsDestInput{
				ChainSelector:           chain.Selector,
				CommitteeVerifierConfig: verifierConfig,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("configuring committee verifier as dest: %w", err)
			}
		}

		return sequences.OnChainOutput{}, nil
	},
)
