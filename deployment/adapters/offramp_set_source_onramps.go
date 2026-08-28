package adapters

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"
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
// GlobalConfig source chain config. Canton stores every source family's OnRamp
// address left-padded to this width. See configure_chain_for_lanes.go and
// harden_canton_inbound_lane.go, which write the same form.
const cantonOnRampAddressBytes = 32

// GetOffRampSourceOnRamps implements [ccipadapters.OffRampSourceOnRampReader].
//
// The result uses the source family's own address width, not the padded width
// stored on the ledger. Callers compare it against addresses that the source
// chain's family adapter produced, so the two encodings must agree.
func (a *CantonChainFamilyAdapter) GetOffRampSourceOnRamps(
	e cldf.Environment,
	localChainSelector uint64,
	sourceChainSelector uint64,
) ([][]byte, error) {
	addressBytes, err := sourceOnRampAddressBytes(sourceChainSelector)
	if err != nil {
		return nil, err
	}

	config, found, err := a.readSourceChainConfig(e, localChainSelector, sourceChainSelector)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return decodeCantonOnRampAddresses(config.OnRampAddresses, addressBytes)
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

// parseCantonOffRampSourceOnRamps decodes the operator-supplied OnRamps and
// left-pads each one to the Canton address width. Duplicates are removed.
func parseCantonOffRampSourceOnRamps(addrs []string) ([]types.TEXT, error) {
	out := make([]types.TEXT, 0, len(addrs))
	seen := make(map[types.TEXT]struct{}, len(addrs))
	for i, addr := range addrs {
		raw, err := hexutil.Decode(strings.TrimSpace(addr))
		if err != nil {
			return nil, fmt.Errorf("onRamps[%d]: invalid hex address %q: %w", i, addr, err)
		}
		if len(raw) == 0 || len(raw) > cantonOnRampAddressBytes {
			return nil, fmt.Errorf("onRamps[%d]: address %q must be 1-%d bytes, got %d",
				i, addr, cantonOnRampAddressBytes, len(raw))
		}

		padded := padCantonOnRampAddress(raw)
		if _, ok := seen[padded]; ok {
			continue
		}
		seen[padded] = struct{}{}
		out = append(out, padded)
	}
	if len(out) == 0 {
		return nil, errors.New("no valid onRamp addresses after parsing")
	}

	return out, nil
}

// padCantonOnRampAddress returns the ledger form of an OnRamp address: hex
// without a 0x prefix, left-padded to the Canton address width.
func padCantonOnRampAddress(raw []byte) types.TEXT {
	return types.TEXT(hex.EncodeToString(gethcommon.LeftPadBytes(raw, cantonOnRampAddressBytes)))
}

// decodeCantonOnRampAddresses converts the ledger form back to the source
// family's own address width. An entry shorter than wantLen, or one whose
// leading bytes are not all zero, is returned unchanged.
func decodeCantonOnRampAddresses(entries []types.TEXT, wantLen int) ([][]byte, error) {
	out := make([][]byte, 0, len(entries))
	for i, entry := range entries {
		raw, err := hex.DecodeString(strings.TrimPrefix(string(entry), "0x"))
		if err != nil {
			return nil, fmt.Errorf("decode onRampAddresses[%d] %q: %w", i, entry, err)
		}
		out = append(out, trimCantonOnRampAddress(raw, wantLen))
	}

	return out, nil
}

// trimCantonOnRampAddress removes left zero padding down to wantLen bytes.
func trimCantonOnRampAddress(raw []byte, wantLen int) []byte {
	if wantLen <= 0 || len(raw) <= wantLen {
		return raw
	}
	prefix := raw[:len(raw)-wantLen]
	for _, b := range prefix {
		if b != 0 {
			return raw
		}
	}

	return raw[len(raw)-wantLen:]
}

// sourceOnRampAddressBytes returns the address width of the source chain's
// family.
//
// An unknown selector, an unregistered family or a zero width is an error. The
// width decides the encoding of the returned addresses, so a wrong value gives
// a caller bytes that never match the source family's own encoding. Fail
// closed instead.
func sourceOnRampAddressBytes(sourceChainSelector uint64) (int, error) {
	family, err := chainsel.GetSelectorFamily(sourceChainSelector)
	if err != nil {
		return 0, fmt.Errorf("get chain family for source chain %d: %w", sourceChainSelector, err)
	}
	adapter, ok := ccipadapters.GetChainFamilyRegistry().GetChainFamily(family)
	if !ok {
		return 0, fmt.Errorf(
			"no chain family adapter registered for source family %q; the binary must import that family's adapter package",
			family)
	}
	length := adapter.GetAddressBytesLength()
	if length == 0 {
		return 0, fmt.Errorf("chain family adapter for %q reports a zero address width", family)
	}

	return int(length), nil
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
// hex case or a 0x prefix.
func normalizeCantonOnRampEntry(entry types.TEXT) string {
	raw, err := hex.DecodeString(strings.TrimPrefix(string(entry), "0x"))
	if err != nil {
		return strings.ToLower(string(entry))
	}

	return string(bytes.TrimLeft(raw, "\x00"))
}
