package sequences

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	executor2 "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// cantonFeeQuoterUSDPerUnitGas formats V2Params.USDPerUnitGas for Canton FeeQuoter UpdatePrices.
// DAML stores usdPerUnitGas as Decimal (e.g. 0.0000000038). chainlink-ccip models it as *big.Int
// for cross-family tooling; on Canton that integer is scaled by 1e10 USD per gas unit (integration
// parity: 38 -> 0.0000000038, matching historical ApplyFeeTokenUpdates+UpdatePrices tests).
func cantonFeeQuoterUSDPerUnitGas(v *big.Int) types.NUMERIC {
	if v == nil || v.Sign() == 0 {
		return types.NUMERIC("0")
	}
	const scale int64 = 10_000_000_000 // 1e10
	r := new(big.Rat).SetFrac(new(big.Int).Set(v), big.NewInt(scale))
	s := strings.TrimRight(strings.TrimRight(r.FloatString(20), "0"), ".")
	if s == "" || s == "-" {
		return types.NUMERIC("0")
	}

	return types.NUMERIC(s)
}

// cantonFeeQuoterUsdPerToken formats lane TokenPrices for Canton FeeQuoter UpdatePrices.
// DAML stores usdPerToken as Decimal USD per whole token (FeeQuoter.daml tests use 20.0 for $20/LINK).
// Lane params carry USD*1e8; divide by 1e8 before encoding (e.g. 1_000_000_000 -> "10").
func cantonFeeQuoterUsdPerToken(v *big.Int) types.NUMERIC {
	if v == nil || v.Sign() <= 0 {
		return types.NUMERIC("0")
	}
	const scale int64 = 100_000_000 // 1e8
	r := new(big.Rat).SetFrac(new(big.Int).Set(v), big.NewInt(scale))
	s := strings.TrimRight(strings.TrimRight(r.FloatString(8), "0"), ".")
	if s == "" || s == "-" {
		return types.NUMERIC("0")
	}

	return types.NUMERIC(s)
}

// ConfigureLaneLegInput carries the lane update plus a datastore for fee quoter resolution.
type ConfigureLaneLegInput struct {
	Lane      lanes.UpdateLanesInput
	DataStore datastore.DataStore
}

var ConfigureLaneLegAsSourceWithInput = operations.NewSequence(
	"CantonConfigureLaneLegAsSourceWithInput",
	semver.MustParse("2.0.0"),
	"Configures a lane leg as source on CCIP 2.0.0",
	configureLaneLegAsSource,
)

var ConfigureLaneLegAsSource = operations.NewSequence(
	"CantonConfigureLaneLegAsSource",
	semver.MustParse("2.0.0"),
	"Configures a lane leg as source on CCIP 2.0.0",
	func(b operations.Bundle, deps chain.BlockChains, input lanes.UpdateLanesInput) (output sequences.OnChainOutput, err error) {
		return configureLaneLegAsSource(b, deps, ConfigureLaneLegInput{Lane: input})
	},
)

func configureLaneLegAsSource(b operations.Bundle, deps chain.BlockChains, input ConfigureLaneLegInput) (output sequences.OnChainOutput, err error) {
	b.Logger.Infof("Canton Configuring lane leg as source. src: %+v, dest: %+v", input.Lane.Source, input.Lane.Dest)

	chain, ok := deps.CantonChains()[input.Lane.Source.Selector]
	if !ok {
		return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.Lane.Source.Selector)
	}

	sourceChain := input.Lane.Source
	destChain := input.Lane.Dest
	participant := chain.Participants[0]
	mcmsEnabled := len(participant.ReadAsPartyIDs) > 0
	var proposalOutputs []contract.ExerciseOutput

	if sourceChain.CantonLaneConfig == nil {
		return sequences.OnChainOutput{}, fmt.Errorf("canton lane config is required on source chain")
	}

	globalConfigRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(sourceChain.CantonLaneConfig.GlobalConfig)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("global config raw instance address: %w", err)
	}
	feeQuoterRef, err := resolveFeeQuoterRef(input, sourceChain.Selector)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("resolve fee quoter: %w", err)
	}
	feeQuoterRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(feeQuoterRef)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("fee quoter raw instance address: %w", err)
	}

	// GlobalConfig - Dest Chain Config
	isEnabled := len(destChain.Router) > 0
	defaultExecutor, err := dsutils.GetRawInstanceAddressFromAddressRef(sourceChain.DefaultExecutor)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("getting default executor: %w", err)
	}
	laneMandatedOutboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(sourceChain.LaneMandatedOutboundCCVs))
	for _, ccv := range sourceChain.LaneMandatedOutboundCCVs {
		outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated outbound CCV: %w", err)
		}
		laneMandatedOutboundCCVs = append(laneMandatedOutboundCCVs, outboundCCV.Binding())
	}
	defaultOutboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(sourceChain.DefaultOutboundCCVs))
	for _, ccv := range sourceChain.DefaultOutboundCCVs {
		outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated outbound CCV: %w", err)
		}
		defaultOutboundCCVs = append(defaultOutboundCCVs, outboundCCV.Binding())
	}
	// Default to false if not specified
	tokenReceiverAllowed := false
	if destChain.TokenReceiverAllowed != nil {
		tokenReceiverAllowed = *destChain.TokenReceiverAllowed
	}
	destChainConfigReport, err := operations.ExecuteOperation(b, global_config.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[core.ApplyDestChainConfigUpdates]{
		InstanceAddress:    globalConfigRaw.InstanceAddress(),
		RawInstanceAddress: globalConfigRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplyDestChainConfigUpdates{
			DestChainConfigUpdates: []core.DestChainConfigArgs{
				{
					DestChainSelector:         types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
					IsEnabled:                 types.BOOL(isEnabled),
					AddressBytesLength:        types.INT64(destChain.AddressBytesLength),
					TokenReceiverAllowed:      types.BOOL(tokenReceiverAllowed),
					BaseExecutionGasCost:      types.INT64(destChain.BaseExecutionGasCost),
					OffRampAddress:            types.TEXT(hex.EncodeToString(destChain.OffRamp)),
					DefaultExecutor:           new(defaultExecutor.Binding()),
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
	if mcmsEnabled && !destChainConfigReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, destChainConfigReport.Output)
	}

	// Executor - Dest Chain Config
	executorRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(sourceChain.DefaultExecutor)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("getting default executor raw instance address: %w", err)
	}
	executorReport, err := operations.ExecuteOperation(b, executor2.ApplyDestChainUpdates, chain, contract.ChoiceInput[executor.ApplyDestChainUpdates]{
		InstanceAddress:    executorRaw.InstanceAddress(),
		RawInstanceAddress: executorRaw.String(),
		MCMSEnabled:        mcmsEnabled,
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
	if mcmsEnabled && !executorReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, executorReport.Output)
	}

	// FeeQuoter - Dest Chain Config
	feeQuoterDestConfigReport, err := operations.ExecuteOperation(b, feequoterop.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[core.ApplyFeeQuoterDestChainConfigUpdates]{
		InstanceAddress:    feeQuoterRaw.InstanceAddress(),
		RawInstanceAddress: feeQuoterRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplyFeeQuoterDestChainConfigUpdates{
			DestChainConfigArgs: []core.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
					DestChainConfig: core.FeeQuoterDestChainConfig{
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
	if mcmsEnabled && !feeQuoterDestConfigReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, feeQuoterDestConfigReport.Output)
	}

	ccipOwnerParty, err := resolveCcipOwnerPartyFromFeeQuoterRef(feeQuoterRef)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("resolve ccipOwner for fee quoter price updates: %w", err)
	}

	priceUpdatersReport, err := operations.ExecuteOperation(b, feequoterop.ApplyPriceUpdatersUpdate, chain, contract.ChoiceInput[core.ApplyPriceUpdatersUpdate]{
		InstanceAddress:    feeQuoterRaw.InstanceAddress(),
		RawInstanceAddress: feeQuoterRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplyPriceUpdatersUpdate{
			AddedPriceUpdaters:   []types.PARTY{types.PARTY(ccipOwnerParty)},
			RemovedPriceUpdaters: nil,
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("ensuring price updater on fee quoter: %w", err)
	}
	if mcmsEnabled && !priceUpdatersReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, priceUpdatersReport.Output)
	}

	tokenPriceUpdates, err := tokenPriceUpdatesFromParams(destChain.TokenPrices)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("building token price updates from lane params: %w", err)
	}

	// FeeQuoter - Update prices.
	updatePricesReport, err := operations.ExecuteOperation(b, feequoterop.UpdatePrices, chain, contract.ChoiceInput[core.UpdatePrices]{
		InstanceAddress:    feeQuoterRaw.InstanceAddress(),
		RawInstanceAddress: feeQuoterRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.UpdatePrices{
			PriceUpdates: core.PriceUpdates{
				TokenPriceUpdates: tokenPriceUpdates,
				GasPriceUpdates: []core.GasPriceUpdate{
					{
						DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
						UsdPerUnitGas: cantonFeeQuoterUSDPerUnitGas(
							destChain.FeeQuoterDestChainConfig.V2Params.USDPerUnitGas),
					},
				},
			},
			Caller: types.PARTY(ccipOwnerParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("updating prices in fee quoter: %w", err)
	}
	if mcmsEnabled && !updatePricesReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, updatePricesReport.Output)
	}

	var cvBatchOps []mcms_types.BatchOperation
	// CommitteeVerifier configure is emitted in a separate mcms-ccv proposal (Run 2).
	if !mcmsEnabled {
		for _, verifierConfig := range input.Lane.Source.CommitteeVerifiers {
			cvReport, err := operations.ExecuteSequence(b, ConfigureCommitteeVerifierAsSource, deps, ConfigureCommitteeVerifierAsSourceInput{
				ChainSelector:           chain.Selector,
				MCMSEnabled:             mcmsEnabled,
				CommitteeVerifierConfig: verifierConfig,
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("configuring committee verifier as source: %w", err)
			}
			cvBatchOps = append(cvBatchOps, cvReport.Output.BatchOps...)
		}
	}

	if !mcmsEnabled {
		return sequences.OnChainOutput{}, nil
	}
	batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for lane configuration: %w", err)
	}
	batchOps := cvBatchOps
	if len(batchOp.Transactions) > 0 {
		batchOps = append(batchOps, batchOp)
	}
	if len(batchOps) == 0 {
		return sequences.OnChainOutput{}, nil
	}

	return sequences.OnChainOutput{BatchOps: batchOps}, nil
}

// CantonOutboundDestFeeInput carries resolved lane legs for Canton→remote outbound dest fee hardening.
type CantonOutboundDestFeeInput struct {
	Source *lanes.ChainDefinition
	Dest   *lanes.ChainDefinition
}

// appendCantonOutboundDestFeeProposalOutputs emits GlobalConfig and FeeQuoter dest updates only
// (no Executor, price updater, or UpdatePrices ops).
func appendCantonOutboundDestFeeProposalOutputs(
	b operations.Bundle,
	chain canton.Chain,
	mcmsEnabled bool,
	globalConfigRaw, feeQuoterRaw contracts.RawInstanceAddress,
	input CantonOutboundDestFeeInput,
) ([]contract.ExerciseOutput, error) {
	if input.Source == nil || input.Dest == nil {
		return nil, fmt.Errorf("source and dest chain definitions are required")
	}
	sourceChain := input.Source
	destChain := input.Dest

	var proposalOutputs []contract.ExerciseOutput

	isEnabled := len(destChain.Router) > 0
	defaultExecutor, err := dsutils.GetRawInstanceAddressFromAddressRef(sourceChain.DefaultExecutor)
	if err != nil {
		return nil, fmt.Errorf("getting default executor: %w", err)
	}
	laneMandatedOutboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(sourceChain.LaneMandatedOutboundCCVs))
	for _, ccv := range sourceChain.LaneMandatedOutboundCCVs {
		outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
		if err != nil {
			return nil, fmt.Errorf("getting lane mandated outbound CCV: %w", err)
		}
		laneMandatedOutboundCCVs = append(laneMandatedOutboundCCVs, outboundCCV.Binding())
	}
	defaultOutboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(sourceChain.DefaultOutboundCCVs))
	for _, ccv := range sourceChain.DefaultOutboundCCVs {
		outboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
		if err != nil {
			return nil, fmt.Errorf("getting default outbound CCV: %w", err)
		}
		defaultOutboundCCVs = append(defaultOutboundCCVs, outboundCCV.Binding())
	}
	tokenReceiverAllowed := false
	if destChain.TokenReceiverAllowed != nil {
		tokenReceiverAllowed = *destChain.TokenReceiverAllowed
	}
	destChainConfigReport, err := operations.ExecuteOperation(b, global_config.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[core.ApplyDestChainConfigUpdates]{
		InstanceAddress:    globalConfigRaw.InstanceAddress(),
		RawInstanceAddress: globalConfigRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplyDestChainConfigUpdates{
			DestChainConfigUpdates: []core.DestChainConfigArgs{
				{
					DestChainSelector:         types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
					IsEnabled:                 types.BOOL(isEnabled),
					AddressBytesLength:        types.INT64(destChain.AddressBytesLength),
					TokenReceiverAllowed:      types.BOOL(tokenReceiverAllowed),
					BaseExecutionGasCost:      types.INT64(destChain.BaseExecutionGasCost),
					OffRampAddress:            types.TEXT(hex.EncodeToString(destChain.OffRamp)),
					DefaultExecutor:           new(defaultExecutor.Binding()),
					LaneMandatedCCVs:          laneMandatedOutboundCCVs,
					DefaultCCVs:               defaultOutboundCCVs,
					MessageNetworkFeeUSDCents: types.NUMERIC(strconv.FormatUint(uint64(destChain.MessageNetworkFeeUSDCents), 10)),
					TokenNetworkFeeUSDCents:   types.NUMERIC(strconv.FormatUint(uint64(destChain.TokenNetworkFeeUSDCents), 10)),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("applying dest chain config updates to global config: %w", err)
	}
	if mcmsEnabled && !destChainConfigReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, destChainConfigReport.Output)
	}

	feeQuoterDestConfigReport, err := operations.ExecuteOperation(b, feequoterop.ApplyDestChainConfigUpdates, chain, contract.ChoiceInput[core.ApplyFeeQuoterDestChainConfigUpdates]{
		InstanceAddress:    feeQuoterRaw.InstanceAddress(),
		RawInstanceAddress: feeQuoterRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplyFeeQuoterDestChainConfigUpdates{
			DestChainConfigArgs: []core.FeeQuoterDestChainConfigArgs{
				{
					DestChainSelector: types.NUMERIC(strconv.FormatUint(destChain.Selector, 10)),
					DestChainConfig: core.FeeQuoterDestChainConfig{
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
		return nil, fmt.Errorf("applying dest chain config updates to fee quoter: %w", err)
	}
	if mcmsEnabled && !feeQuoterDestConfigReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, feeQuoterDestConfigReport.Output)
	}

	return proposalOutputs, nil
}

func resolveFeeQuoterRef(input ConfigureLaneLegInput, chainSelector uint64) (datastore.AddressRef, error) {
	if input.DataStore == nil {
		return datastore.AddressRef{}, fmt.Errorf("datastore is required on ConfigureLaneLegInput")
	}

	return input.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	))
}

var ConfigureLaneLegAsDest = operations.NewSequence(
	"CantonConfigureLaneLegAsDest",
	semver.MustParse("2.0.0"),
	"Configures a lane lad as dest on CCIP 2.0.0",
	func(b operations.Bundle, deps chain.BlockChains, input lanes.UpdateLanesInput) (output sequences.OnChainOutput, err error) {
		b.Logger.Infof("Canton Configuring lane leg as dest. src: %+v, dest: %+v", input.Source, input.Dest)

		chain, ok := deps.CantonChains()[input.Dest.Selector]
		if !ok {
			return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found", input.Source.Selector)
		}

		sourceChain := input.Source
		destChain := input.Dest
		mcmsEnabled := len(chain.Participants[0].ReadAsPartyIDs) > 0
		var proposalOutputs []contract.ExerciseOutput

		globalConfigRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(destChain.CantonLaneConfig.GlobalConfig)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("global config raw instance address: %w", err)
		}

		// GlobalConfig - Source Chain Config
		laneMandatedInboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(destChain.LaneMandatedInboundCCVs))
		for _, ccv := range destChain.LaneMandatedInboundCCVs {
			inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated inbound CCV: %w", err)
			}
			laneMandatedInboundCCVs = append(laneMandatedInboundCCVs, inboundCCV.Binding())
		}
		defaultInboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(destChain.DefaultInboundCCVs))
		for _, ccv := range destChain.DefaultInboundCCVs {
			inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("getting lane mandated inbound CCV: %w", err)
			}
			defaultInboundCCVs = append(defaultInboundCCVs, inboundCCV.Binding())
		}
		inboundOnRampAddresses := []types.TEXT{
			types.TEXT(hex.EncodeToString(gethcommon.LeftPadBytes(sourceChain.OnRamp, 32))),
		}
		sourceChainConfigReport, err := operations.ExecuteOperation(b, global_config.ApplySourceChainConfigUpdates, chain, contract.ChoiceInput[core.ApplySourceChainConfigUpdates]{
			InstanceAddress:    globalConfigRaw.InstanceAddress(),
			RawInstanceAddress: globalConfigRaw.String(),
			MCMSEnabled:        mcmsEnabled,
			Args: core.ApplySourceChainConfigUpdates{
				SourceChainConfigUpdates: []core.SourceChainConfigArgs{
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
		if mcmsEnabled && !sourceChainConfigReport.Output.Executed() {
			proposalOutputs = append(proposalOutputs, sourceChainConfigReport.Output)
		}

		var cvBatchOps []mcms_types.BatchOperation
		// CommitteeVerifier configure is emitted in a separate mcms-ccv proposal (Run 2).
		if !mcmsEnabled {
			for _, verifierConfig := range input.Dest.CommitteeVerifiers {
				cvReport, err := operations.ExecuteSequence(b, ConfigureCommitteeVerifierAsDest, deps, ConfigureCommitteeVerifierAsDestInput{
					ChainSelector:           chain.Selector,
					MCMSEnabled:             mcmsEnabled,
					CommitteeVerifierConfig: verifierConfig,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("configuring committee verifier as dest: %w", err)
				}
				cvBatchOps = append(cvBatchOps, cvReport.Output.BatchOps...)
			}
		}

		if !mcmsEnabled {
			return sequences.OnChainOutput{}, nil
		}
		batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for lane configuration: %w", err)
		}
		batchOps := cvBatchOps
		if len(batchOp.Transactions) > 0 {
			batchOps = append(batchOps, batchOp)
		}
		if len(batchOps) == 0 {
			return sequences.OnChainOutput{}, nil
		}

		return sequences.OnChainOutput{BatchOps: batchOps}, nil
	},
)

func resolveCcipOwnerPartyFromFeeQuoterRef(feeQuoterRef datastore.AddressRef) (string, error) {
	for _, label := range feeQuoterRef.Labels.List() {
		at := strings.LastIndex(label, "@")
		if at < 0 || at+1 >= len(label) {
			continue
		}
		party := label[at+1:]
		if strings.Contains(party, "::") {
			return party, nil
		}
	}

	return "", fmt.Errorf("ccipOwner party not found in FeeQuoter labels")
}

func parseInstrumentPriceKey(instrument string) (admin, id string, err error) {
	instrument = strings.TrimSpace(instrument)
	lastColon := strings.LastIndex(instrument, ":")
	if lastColon <= 0 || lastColon+1 >= len(instrument) {
		return "", "", fmt.Errorf("invalid token price instrument key %q, expected format <admin>:<id>", instrument)
	}
	admin = strings.TrimSpace(instrument[:lastColon])
	id = strings.TrimSpace(instrument[lastColon+1:])
	if admin == "" || id == "" {
		return "", "", fmt.Errorf("invalid token price instrument key %q, expected format <admin>:<id>", instrument)
	}

	return admin, id, nil
}

func tokenPriceUpdatesFromParams(tokenPrices map[string]*big.Int) ([]core.TokenPriceUpdate, error) {
	if len(tokenPrices) == 0 {
		return nil, nil
	}
	updates := make([]core.TokenPriceUpdate, 0, len(tokenPrices))
	for instrument, price := range tokenPrices {
		if price == nil {
			continue
		}
		admin, id, err := parseInstrumentPriceKey(instrument)
		if err != nil {
			return nil, err
		}
		updates = append(updates, core.TokenPriceUpdate{
			InstrumentId: splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(admin),
				Id:    types.TEXT(id),
			},
			UsdPerToken: cantonFeeQuoterUsdPerToken(price),
		})
	}

	return updates, nil
}
