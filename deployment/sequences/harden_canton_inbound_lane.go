package sequences

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Masterminds/semver/v3"
	gethcommon "github.com/ethereum/go-ethereum/common"

	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// HardenCantonInboundLaneInput carries pre-resolved contract refs and lane hardening params.
// Emits exactly two MCMS operations: FeeQuoter::UpdatePrices then GlobalConfig::ApplySourceChainConfigUpdates.
type HardenCantonInboundLaneInput struct {
	CantonChainSelector       uint64
	RemoteSourceChainSelector uint64
	GlobalConfigRef           datastore.AddressRef
	FeeQuoterRef              datastore.AddressRef
	DefaultInboundCCVs        []datastore.AddressRef
	RemoteOnRampAddress       gethcommon.Address
	TokenPrices               map[string]*big.Int
	USDPerUnitGas             *big.Int
}

var HardenCantonInboundLane = operations.NewSequence(
	"HardenCantonInboundLane",
	semver.MustParse("2.0.0"),
	"Applies Canton inbound lane hardening (fee quoter prices + invalid default inbound CCV)",
	hardenCantonInboundLane,
)

func hardenCantonInboundLane(
	b operations.Bundle,
	deps chain.BlockChains,
	input HardenCantonInboundLaneInput,
) (output ccipseq.OnChainOutput, err error) {
	chain, ok := deps.CantonChains()[input.CantonChainSelector]
	if !ok {
		return ccipseq.OnChainOutput{}, fmt.Errorf("canton chain %d not found", input.CantonChainSelector)
	}
	if len(chain.Participants) == 0 {
		return ccipseq.OnChainOutput{}, fmt.Errorf("canton chain %d has no participants", input.CantonChainSelector)
	}

	mcmsEnabled := len(chain.Participants[0].ReadAsPartyIDs) > 0
	var proposalOutputs []contract.ExerciseOutput

	globalConfigRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(input.GlobalConfigRef)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("global config raw instance address: %w", err)
	}
	feeQuoterRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(input.FeeQuoterRef)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("fee quoter raw instance address: %w", err)
	}

	ccipOwnerParty, err := resolveCcipOwnerPartyFromFeeQuoterRef(input.FeeQuoterRef)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve ccipOwner for fee quoter price updates: %w", err)
	}

	tokenPriceUpdates, err := tokenPriceUpdatesFromParams(input.TokenPrices)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("building token price updates: %w", err)
	}

	updatePricesReport, err := operations.ExecuteOperation(b, feequoterop.UpdatePrices, chain, contract.ChoiceInput[core.UpdatePrices]{
		InstanceAddress:    feeQuoterRaw.InstanceAddress(),
		RawInstanceAddress: feeQuoterRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.UpdatePrices{
			PriceUpdates: core.PriceUpdates{
				TokenPriceUpdates: tokenPriceUpdates,
				GasPriceUpdates: []core.GasPriceUpdate{
					{
						DestChainSelector: types.NUMERIC(strconv.FormatUint(input.RemoteSourceChainSelector, 10)),
						UsdPerUnitGas:     cantonFeeQuoterUSDPerUnitGas(input.USDPerUnitGas),
					},
				},
			},
			Caller: types.PARTY(ccipOwnerParty),
		},
	})
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("updating prices in fee quoter: %w", err)
	}
	if mcmsEnabled && !updatePricesReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, updatePricesReport.Output)
	}

	defaultInboundCCVs := make([]chainlinkapi.RawInstanceAddress, 0, len(input.DefaultInboundCCVs))
	for _, ccv := range input.DefaultInboundCCVs {
		inboundCCV, err := dsutils.GetRawInstanceAddressFromAddressRef(ccv)
		if err != nil {
			return ccipseq.OnChainOutput{}, fmt.Errorf("resolve default inbound CCV: %w", err)
		}
		defaultInboundCCVs = append(defaultInboundCCVs, inboundCCV.Binding())
	}

	sourceChainConfigReport, err := operations.ExecuteOperation(b, global_config.ApplySourceChainConfigUpdates, chain, contract.ChoiceInput[core.ApplySourceChainConfigUpdates]{
		InstanceAddress:    globalConfigRaw.InstanceAddress(),
		RawInstanceAddress: globalConfigRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplySourceChainConfigUpdates{
			SourceChainConfigUpdates: []core.SourceChainConfigArgs{
				{
					SourceChainSelector: types.NUMERIC(strconv.FormatUint(input.RemoteSourceChainSelector, 10)),
					IsEnabled:           types.BOOL(true),
					OnRampAddresses: []types.TEXT{
						types.TEXT(hex.EncodeToString(gethcommon.LeftPadBytes(input.RemoteOnRampAddress.Bytes(), 32))),
					},
					DefaultCCVs:      defaultInboundCCVs,
					LaneMandatedCCVs: nil,
				},
			},
		},
	})
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("applying source chain config updates to global config: %w", err)
	}
	if mcmsEnabled && !sourceChainConfigReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, sourceChainConfigReport.Output)
	}

	if !mcmsEnabled {
		return ccipseq.OnChainOutput{}, nil
	}

	batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("build MCMS batch for inbound hardening: %w", err)
	}
	if len(batchOp.Transactions) == 0 {
		return ccipseq.OnChainOutput{}, nil
	}

	return ccipseq.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
}
