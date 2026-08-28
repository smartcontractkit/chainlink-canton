package adapters

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
)

var (
	_ ccipadapters.OffRampSourceOnRampSetter = (*CantonChainFamilyAdapter)(nil)
	_ ccipadapters.OffRampSourceOnRampReader = (*CantonChainFamilyAdapter)(nil)
)

// cantonOnRampAddressBytes is the width of an OnRamp address in the Canton
// GlobalConfig source chain config.
//
// The width comes from the source chain, not from Canton: an EVM source emits
// its message with the OnRamp address abi-encoded, so Canton must store the same
// 32 bytes for the equality check on an incoming message to succeed. Every other
// supported source family also uses 32 bytes. See configure_chain_for_lanes.go
// and harden_canton_inbound_lane.go, which store the same form.
const cantonOnRampAddressBytes = 32

// GetOffRampSourceOnRamps implements [ccipadapters.OffRampSourceOnRampReader].
//
// The result is the stored address exactly as the source chain writes it into
// its messages, which for an EVM source is abi.encode(address): 32 bytes, left
// zero-padded. The width is not reduced to the source family's contract address
// width, because the OffRamp matches an incoming message by hashing the onramp
// bytes the message carries.
func (a *CantonChainFamilyAdapter) GetOffRampSourceOnRamps(
	e cldf.Environment,
	localChainSelector uint64,
	sourceChainSelector uint64,
) ([][]byte, error) {
	config, found, err := a.readSourceChainConfig(e, localChainSelector, sourceChainSelector)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return decodeCantonOnRampAddresses(config.OnRampAddresses)
}

// SetOffRampSourceOnRamps implements [ccipadapters.OffRampSourceOnRampSetter].
//
// Only the OnRamp allowlist changes. IsEnabled, DefaultCCVs and
// LaneMandatedCCVs are copied from the current source chain config, because the
// Canton choice replaces the whole entry.
func (a *CantonChainFamilyAdapter) SetOffRampSourceOnRamps(
	e cldf.Environment,
	update ccipadapters.OffRampSetSourceOnRampsEntry,
) (*mcms_types.BatchOperation, bool, error) {
	chain, ok := e.BlockChains.CantonChains()[update.LocalChainSelector]
	if !ok || len(chain.Participants) == 0 {
		return nil, false, fmt.Errorf("canton chain %d not found or has no participants", update.LocalChainSelector)
	}

	desired, err := parseCantonOffRampSourceOnRamps(update.OnRamps)
	if err != nil {
		return nil, false, err
	}

	globalConfigRaw, err := a.globalConfigRawAddress(e.DataStore, update.LocalChainSelector)
	if err != nil {
		return nil, false, err
	}

	current, found, err := a.readSourceChainConfig(e, update.LocalChainSelector, update.SourceChainSelector)
	if err != nil {
		return nil, false, err
	}
	if !found {
		// This method updates an existing lane. Creating one needs the fee quoter
		// and the inbound CCVs, which only the lane configure changeset resolves.
		return nil, false, fmt.Errorf(
			"canton chain %d has no source chain config for source chain %d; configure the lane first",
			update.LocalChainSelector, update.SourceChainSelector)
	}

	if cantonOnRampSetsEqual(current.OnRampAddresses, desired) {
		e.Logger.Infow("Canton source OnRamp allowlist already matches the desired state, skipping",
			"localChain", update.LocalChainSelector,
			"sourceChain", update.SourceChainSelector,
			"globalConfig", globalConfigRaw.InstanceAddress().String(),
			"onRampCount", len(desired),
		)

		return nil, true, nil
	}

	mcmsEnabled := len(chain.Participants[0].ReadAsPartyIDs) > 0

	report, err := cldf_ops.ExecuteOperation(e.OperationsBundle, global_config.ApplySourceChainConfigUpdates, chain, contract.ChoiceInput[core.ApplySourceChainConfigUpdates]{
		InstanceAddress:    globalConfigRaw.InstanceAddress(),
		RawInstanceAddress: globalConfigRaw.String(),
		MCMSEnabled:        mcmsEnabled,
		Args: core.ApplySourceChainConfigUpdates{
			SourceChainConfigUpdates: []core.SourceChainConfigArgs{
				{
					SourceChainSelector: types.NUMERIC(strconv.FormatUint(update.SourceChainSelector, 10)),
					IsEnabled:           current.IsEnabled,
					OnRampAddresses:     desired,
					DefaultCCVs:         current.DefaultCCVs,
					LaneMandatedCCVs:    current.LaneMandatedCCVs,
				},
			},
		},
	})
	if err != nil {
		return nil, false, fmt.Errorf("apply source chain config on canton chain %d: %w", update.LocalChainSelector, err)
	}

	if !mcmsEnabled || report.Output.Executed() {
		// The exercise ran directly with the operator's party. There is nothing to
		// put in a proposal.
		return nil, false, nil
	}

	batchOp, err := contract.NewBatchOperationFromExercises([]contract.ExerciseOutput{report.Output})
	if err != nil {
		return nil, false, fmt.Errorf("build MCMS batch for source chain config update: %w", err)
	}
	if len(batchOp.Transactions) == 0 {
		return nil, false, nil
	}

	return &batchOp, false, nil
}

// readSourceChainConfig returns the GlobalConfig source chain config for
// sourceChainSelector. The bool is false when the lane is not configured.
func (a *CantonChainFamilyAdapter) readSourceChainConfig(
	e cldf.Environment,
	localChainSelector uint64,
	sourceChainSelector uint64,
) (core.SourceChainConfig2, bool, error) {
	chain, ok := e.BlockChains.CantonChains()[localChainSelector]
	if !ok || len(chain.Participants) == 0 {
		return core.SourceChainConfig2{}, false, fmt.Errorf(
			"canton chain %d not found or has no participants", localChainSelector)
	}

	globalConfigRaw, err := a.globalConfigRawAddress(e.DataStore, localChainSelector)
	if err != nil {
		return core.SourceChainConfig2{}, false, err
	}

	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}

	globalConfig, err := readCantonGlobalConfig(ctx, chain, globalConfigRaw.InstanceAddress())
	if err != nil {
		return core.SourceChainConfig2{}, false, err
	}

	return findCantonSourceChainConfig(globalConfig.SourceChainConfigs, sourceChainSelector)
}

func (a *CantonChainFamilyAdapter) globalConfigRawAddress(
	ds datastore.DataStore,
	chainSelector uint64,
) (contracts.RawInstanceAddress, error) {
	ref, err := findContractRef(
		ds,
		chainSelector,
		datastore.ContractType(global_config.ContractType),
		global_config.Version,
		"",
	)
	if err != nil {
		return "", fmt.Errorf("resolve global config on chain %d: %w", chainSelector, err)
	}

	raw, err := dsutils.GetRawInstanceAddressFromAddressRef(ref)
	if err != nil {
		return "", fmt.Errorf("global config raw instance address: %w", err)
	}

	return raw, nil
}

// readCantonGlobalConfig reads the active GlobalConfig contract from the ledger.
func readCantonGlobalConfig(
	ctx context.Context,
	chain canton.Chain,
	globalConfig contracts.InstanceAddress,
) (*core.GlobalConfig, error) {
	participant := chain.Participants[0]

	active, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		core.GlobalConfig{}.GetTemplateID(),
		globalConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("find active contract for global config %s: %w", globalConfig.String(), err)
	}
	if active == nil {
		return nil, fmt.Errorf("no active contract found for global config %s", globalConfig.String())
	}

	created, err := bindings.UnmarshalCreatedEvent[core.GlobalConfig](active.GetCreatedEvent())
	if err != nil {
		return nil, fmt.Errorf("unmarshal global config %s: %w", globalConfig.String(), err)
	}

	return created, nil
}

// findCantonSourceChainConfig looks up sourceChainSelector in a GlobalConfig
// source chain config map.
//
// The map key is a DAML NUMERIC. The ledger returns it with a trailing dot
// (e.g. "16015286601757825753."), but the write path formats it without one, so
// the key is parsed instead of built.
func findCantonSourceChainConfig(
	configs map[types.NUMERIC]core.SourceChainConfig2,
	sourceChainSelector uint64,
) (core.SourceChainConfig2, bool, error) {
	for key, config := range configs {
		parsed, err := internalparse.Uint64Checked(string(key))
		if err != nil {
			return core.SourceChainConfig2{}, false, fmt.Errorf(
				"parse source chain selector from GlobalConfig SourceChainConfigs map key %s: %w", key, err)
		}
		if parsed == sourceChainSelector {
			return config, true, nil
		}
	}

	return core.SourceChainConfig2{}, false, nil
}

// parseCantonOffRampSourceOnRamps decodes the caller-supplied OnRamps into the
// ledger form. Duplicates are removed.
//
// Each entry must already be the address exactly as the source chain writes it
// into its messages, which is what [ccipadapters.OffRampSetSourceOnRampsEntry]
// documents. No padding happens here: padding the caller's bytes would hide a
// caller that supplied the wrong encoding, and the resulting lane would reject
// every message.
func parseCantonOffRampSourceOnRamps(addrs []string) ([]types.TEXT, error) {
	out := make([]types.TEXT, 0, len(addrs))
	seen := make(map[types.TEXT]struct{}, len(addrs))
	for i, addr := range addrs {
		raw, err := hexutil.Decode(strings.TrimSpace(addr))
		if err != nil {
			return nil, fmt.Errorf("onRamps[%d]: invalid hex address %q: %w", i, addr, err)
		}
		if len(raw) != cantonOnRampAddressBytes {
			return nil, fmt.Errorf(
				"onRamps[%d]: address %q must be %d bytes, got %d; pass the address as the source chain writes it into its messages (for an EVM source, abi.encode(address))",
				i, addr, cantonOnRampAddressBytes, len(raw))
		}

		entry := types.TEXT(hex.EncodeToString(raw))
		if _, ok := seen[entry]; ok {
			continue
		}
		seen[entry] = struct{}{}
		out = append(out, entry)
	}
	if len(out) == 0 {
		return nil, errors.New("no valid onRamp addresses after parsing")
	}

	return out, nil
}

// decodeCantonOnRampAddresses decodes the ledger form into raw bytes. The width
// is preserved, so the result is the wire encoding the source chain uses.
func decodeCantonOnRampAddresses(entries []types.TEXT) ([][]byte, error) {
	out := make([][]byte, 0, len(entries))
	for i, entry := range entries {
		raw, err := hex.DecodeString(strings.TrimPrefix(string(entry), "0x"))
		if err != nil {
			return nil, fmt.Errorf("decode onRampAddresses[%d] %q: %w", i, entry, err)
		}
		out = append(out, raw)
	}

	return out, nil
}

// cantonOnRampSetsEqual compares two ledger-form OnRamp lists without regard to
// order.
func cantonOnRampSetsEqual(current []types.TEXT, desired []types.TEXT) bool {
	if len(current) != len(desired) {
		return false
	}
	counts := make(map[string]int, len(current))
	for _, entry := range current {
		counts[normalizeCantonOnRampEntry(entry)]++
	}
	for _, entry := range desired {
		key := normalizeCantonOnRampEntry(entry)
		counts[key]--
		if counts[key] < 0 {
			return false
		}
	}

	return true
}

// normalizeCantonOnRampEntry makes two ledger entries comparable regardless of
// hex case or a 0x prefix. The width is significant, so no padding is stripped:
// two addresses that differ only in width are different onramps on the wire.
func normalizeCantonOnRampEntry(entry types.TEXT) string {
	raw, err := hex.DecodeString(strings.TrimPrefix(string(entry), "0x"))
	if err != nil {
		return strings.ToLower(string(entry))
	}

	return string(raw)
}
