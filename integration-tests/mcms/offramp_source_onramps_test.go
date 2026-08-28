package tests

import (
	"encoding/hex"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	// Registers the EVM chain family adapter. The Canton reader asks the registry
	// for the source family's address width.
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/committeeverifier"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	cantonadapters "github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// sepoliaChainSelector is the EVM source chain of the lane under test. A real EVM
// selector is required: the reader resolves the address width from the source
// chain's family.
const sepoliaChainSelector = uint64(16015286601757825753)

// baseSepoliaChainSelector is a second inbound lane that the update must leave
// alone.
const baseSepoliaChainSelector = uint64(10344971235874465080)

const (
	legacyOnRampHex = "0x1111111111111111111111111111111111111111"
	newOnRampHex    = "0x2222222222222222222222222222222222222222"
	// bystanderOnRampHex belongs to the lane that must not change.
	bystanderOnRampHex = "0x3333333333333333333333333333333333333333"
)

// TestOffRampSourceOnRamps_DualAllowlist covers the Canton side of an OnRamp
// redeployment on an EVM source chain, where Canton is the destination.
//
// The allowlist moves through the three states the phase changesets need:
//
//	[legacy]        -> the lane before the redeployment
//	[legacy, new]   -> both OnRamps accepted while the switch is in progress
//	[new]           -> after cleanup
//
// It also asserts that the setter leaves IsEnabled and DefaultCCVs untouched,
// and that it reports a skip when the allowlist already matches.
func TestOffRampSourceOnRamps_DualAllowlist(t *testing.T) {
	t.Parallel()

	sharedEnv := GetSharedCCIPMCMSEnvironment(t)
	participant := sharedEnv.Participant
	party := sharedEnv.CcipOwner

	// The adapter builds an MCMS batch operation instead of submitting when the
	// participant can only read as the owning party. This test drives the direct
	// path, so it needs a participant that can act as the ccipOwner.
	require.Empty(t, participant.ReadAsPartyIDs,
		"this test exercises the direct submission path; a read-as participant needs the MCMS flow")

	chainSelector := chainsel.CANTON_LOCALNET.Selector
	chainSelectorNumeric := types.NUMERIC(strconv.FormatUint(chainSelector, 10))

	cantonChain := &canton.Chain{
		ChainMetadata: canton.ChainMetadata{Selector: chainSelector},
		Participants:  []canton.Participant{participant},
	}
	bundle := cld_ops.NewBundle(t.Context, logger.Test(t), cld_ops.NewMemoryReporter())

	// ------------------------------------------------------------------
	// Deploy the CCIP contracts and record the GlobalConfig in a datastore
	// ------------------------------------------------------------------

	deployOut, err := cld_ops.ExecuteSequence(bundle, sequences.DeployChainContracts, *cantonChain, sequences.DeployChainContractsParams{
		CCIPOwnerParty: party,
		RMNOwnerParty:  party,
		CommitteeVerifiers: []sequences.CommitteeVerifierParams{{
			Template: committeeverifier.CommitteeVerifier{
				Owner:                        types.PARTY(party),
				CcipOwner:                    types.PARTY(party),
				VersionTag:                   "test-v1",
				StorageLocations:             []types.TEXT{"ipfs://test"},
				StorageLocationsAdmin:        types.PARTY(party),
				PendingStorageLocationsAdmin: types.PARTY(party),
			},
		}},
		GlobalConfig: sequences.GlobalConfigParams{
			Template: core.GlobalConfig{
				ChainSelector: chainSelectorNumeric,
			},
		},
		RMNRemote: sequences.RMNRemoteParams{
			Template: core.RMNRemote{
				RmnOwner: types.PARTY(party),
			},
		},
		NativeInstrumentId: splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(party),
			Id:    types.TEXT("LINK-" + uuid.New().String()[:8]),
		},
	})
	require.NoError(t, err, "deploy chain contracts")

	var (
		gcRaw  contracts.RawInstanceAddress
		ccvRaw chainlinkapi.RawInstanceAddress
	)
	for _, addr := range deployOut.Output.Addresses {
		raw, err := contracts.RawInstanceAddressFromString(addr.Labels.List()[0])
		require.NoError(t, err)
		switch string(addr.Type) {
		case string(global_config.ContractType):
			gcRaw = raw
		case "CommitteeVerifier":
			ccvRaw = chainlinkapi.RawInstanceAddress{Unpack: types.TEXT(raw.String())}
		}
	}
	require.NotEmpty(t, gcRaw.String(), "GlobalConfig address not found in deploy output")
	require.NotEmpty(t, ccvRaw.Unpack, "CommitteeVerifier address not found in deploy output")

	// The adapter resolves the GlobalConfig from the datastore. The raw instance
	// address lives in the first label, which is what
	// dsutils.GetRawInstanceAddressFromAddressRef reads.
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(global_config.ContractType),
		Version:       global_config.Version,
		Address:       gcRaw.InstanceAddress().Hex(),
		Labels:        datastore.NewLabelSet(gcRaw.String()),
	}))

	env := cldf.Environment{
		Logger:           logger.Test(t),
		GetContext:       t.Context,
		DataStore:        ds.Seal(),
		BlockChains:      chain.NewBlockChainsFromSlice([]chain.BlockChain{cantonChain}),
		OperationsBundle: bundle,
	}

	adapter := &cantonadapters.CantonChainFamilyAdapter{}

	// ------------------------------------------------------------------
	// An unconfigured lane reads as empty, and cannot be updated
	// ------------------------------------------------------------------

	onRamps, err := adapter.GetOffRampSourceOnRamps(env, chainSelector, sepoliaChainSelector)
	require.NoError(t, err, "read an unconfigured source chain")
	require.Empty(t, onRamps)

	_, _, err = adapter.SetOffRampSourceOnRamps(env, ccipadapters.OffRampSetSourceOnRampsEntry{
		LocalChainSelector:  chainSelector,
		SourceChainSelector: sepoliaChainSelector,
		OnRamps:             []string{legacyOnRampHex},
	})
	require.Error(t, err, "the setter must not create a lane")
	require.Contains(t, err.Error(), "configure the lane first")

	// ------------------------------------------------------------------
	// Seed the lane with the legacy OnRamp only
	//
	// A second inbound lane and an outbound lane are seeded as well, so the
	// blast-radius assertion below has something to protect.
	// ------------------------------------------------------------------

	applySourceChainConfig(t, bundle, cantonChain, gcRaw, core.SourceChainConfigArgs{
		SourceChainSelector: types.NUMERIC(strconv.FormatUint(sepoliaChainSelector, 10)),
		IsEnabled:           types.BOOL(true),
		OnRampAddresses:     []types.TEXT{paddedOnRamp(t, legacyOnRampHex)},
		DefaultCCVs:         []chainlinkapi.RawInstanceAddress{ccvRaw},
		LaneMandatedCCVs:    nil,
	})

	// The bystander lane differs from the target in its OnRamp and in IsEnabled,
	// so a write that leaked into it is visible. A CCV must not repeat across
	// DefaultCCVs and LaneMandatedCCVs; the contract rejects that with
	// DuplicateCCVNotAllowed.
	applySourceChainConfig(t, bundle, cantonChain, gcRaw, core.SourceChainConfigArgs{
		SourceChainSelector: types.NUMERIC(strconv.FormatUint(baseSepoliaChainSelector, 10)),
		IsEnabled:           types.BOOL(false),
		OnRampAddresses:     []types.TEXT{paddedOnRamp(t, bystanderOnRampHex)},
		DefaultCCVs:         []chainlinkapi.RawInstanceAddress{ccvRaw},
		LaneMandatedCCVs:    nil,
	})

	applyDestChainConfig(t, bundle, cantonChain, gcRaw, core.DestChainConfigArgs{
		DestChainSelector:         types.NUMERIC(strconv.FormatUint(sepoliaChainSelector, 10)),
		IsEnabled:                 true,
		AddressBytesLength:        20,
		TokenReceiverAllowed:      true,
		BaseExecutionGasCost:      21000,
		OffRampAddress:            types.TEXT("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		DefaultCCVs:               []chainlinkapi.RawInstanceAddress{ccvRaw},
		MessageNetworkFeeUSDCents: "100",
		TokenNetworkFeeUSDCents:   "50",
	})

	onRamps, err = adapter.GetOffRampSourceOnRamps(env, chainSelector, sepoliaChainSelector)
	require.NoError(t, err)
	require.Equal(t, []string{legacyOnRampHex}, hexOnRamps(onRamps),
		"the reader must return the EVM address width, not the padded ledger width")

	// ------------------------------------------------------------------
	// Phase 1: allowlist both OnRamps
	// ------------------------------------------------------------------

	before := readGlobalConfig(t, env, cantonChain, gcRaw)

	batchOp, skipped, err := adapter.SetOffRampSourceOnRamps(env, ccipadapters.OffRampSetSourceOnRampsEntry{
		LocalChainSelector:  chainSelector,
		SourceChainSelector: sepoliaChainSelector,
		OnRamps:             []string{legacyOnRampHex, newOnRampHex},
	})
	require.NoError(t, err)
	require.False(t, skipped)
	require.Nil(t, batchOp, "the exercise runs directly when MCMS is off")

	onRamps, err = adapter.GetOffRampSourceOnRamps(env, chainSelector, sepoliaChainSelector)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{legacyOnRampHex, newOnRampHex}, hexOnRamps(onRamps))

	// Nothing but the target lane's OnRamp list may have moved.
	after := readGlobalConfig(t, env, cantonChain, gcRaw)
	requireOnlyOnRampsChanged(t, before, after, sepoliaChainSelector)

	// ------------------------------------------------------------------
	// A repeat of the same update is a no-op
	// ------------------------------------------------------------------

	batchOp, skipped, err = adapter.SetOffRampSourceOnRamps(env, ccipadapters.OffRampSetSourceOnRampsEntry{
		LocalChainSelector:  chainSelector,
		SourceChainSelector: sepoliaChainSelector,
		// Reversed order, to prove the comparison ignores order.
		OnRamps: []string{newOnRampHex, legacyOnRampHex},
	})
	require.NoError(t, err)
	require.True(t, skipped, "an allowlist that already matches must be skipped")
	require.Nil(t, batchOp)

	// ------------------------------------------------------------------
	// Cleanup: drop the legacy OnRamp
	// ------------------------------------------------------------------

	before = readGlobalConfig(t, env, cantonChain, gcRaw)

	_, skipped, err = adapter.SetOffRampSourceOnRamps(env, ccipadapters.OffRampSetSourceOnRampsEntry{
		LocalChainSelector:  chainSelector,
		SourceChainSelector: sepoliaChainSelector,
		OnRamps:             []string{newOnRampHex},
	})
	require.NoError(t, err)
	require.False(t, skipped)

	onRamps, err = adapter.GetOffRampSourceOnRamps(env, chainSelector, sepoliaChainSelector)
	require.NoError(t, err)
	require.Equal(t, []string{newOnRampHex}, hexOnRamps(onRamps))

	// Shrinking the allowlist must be as narrow as growing it.
	after = readGlobalConfig(t, env, cantonChain, gcRaw)
	requireOnlyOnRampsChanged(t, before, after, sepoliaChainSelector)
}

// applySourceChainConfig writes a source chain config straight to the
// GlobalConfig, without going through the lane configure sequence.
func applySourceChainConfig(
	t *testing.T,
	bundle cld_ops.Bundle,
	cantonChain *canton.Chain,
	gcRaw contracts.RawInstanceAddress,
	args core.SourceChainConfigArgs,
) {
	t.Helper()

	_, err := cld_ops.ExecuteOperation(bundle, global_config.ApplySourceChainConfigUpdates, *cantonChain, opcontract.ChoiceInput[core.ApplySourceChainConfigUpdates]{
		InstanceAddress:    gcRaw.InstanceAddress(),
		RawInstanceAddress: gcRaw.String(),
		MCMSEnabled:        false,
		Args: core.ApplySourceChainConfigUpdates{
			SourceChainConfigUpdates: []core.SourceChainConfigArgs{args},
		},
	})
	require.NoError(t, err, "apply source chain config")
}

// readGlobalConfig reads the whole GlobalConfig off the ledger, so the test can
// compare the state before and after an update.
func readGlobalConfig(
	t *testing.T,
	env cldf.Environment,
	cantonChain *canton.Chain,
	gcRaw contracts.RawInstanceAddress,
) *core.GlobalConfig {
	t.Helper()

	active, err := opcontract.FindActiveContractByInstanceAddress(
		env.GetContext(),
		cantonChain.Participants[0].LedgerServices.State,
		opcontract.LedgerQueryParties(cantonChain.Participants[0]),
		core.GlobalConfig{}.GetTemplateID(),
		gcRaw.InstanceAddress(),
	)
	require.NoError(t, err)
	require.NotNil(t, active)

	created, err := bindings.UnmarshalCreatedEvent[core.GlobalConfig](active.GetCreatedEvent())
	require.NoError(t, err)

	return created
}

// requireOnlyOnRampsChanged asserts that an update touched nothing but
// OnRampAddresses on the target source chain config.
//
// The Canton choice replaces a whole source chain config entry, so the adapter
// has to copy every other field forward. This is the assertion that catches a
// dropped field.
func requireOnlyOnRampsChanged(
	t *testing.T,
	before *core.GlobalConfig,
	after *core.GlobalConfig,
	targetSourceChainSelector uint64,
) {
	t.Helper()

	require.Equal(t, before.InstanceId, after.InstanceId)
	require.Equal(t, before.CcipOwner, after.CcipOwner)
	require.Equal(t, before.ChainSelector, after.ChainSelector)
	require.Equal(t, before.DestChainConfigs, after.DestChainConfigs,
		"an inbound lane update must not touch the outbound lanes")

	require.Len(t, before.SourceChainConfigs, len(after.SourceChainConfigs),
		"an update must not add or remove a source chain config")

	targetKeys := 0
	for key, beforeConfig := range before.SourceChainConfigs {
		afterConfig, ok := after.SourceChainConfigs[key]
		require.Truef(t, ok, "source chain config %q disappeared", key)

		if sourceChainSelectorFromKey(t, key) != targetSourceChainSelector {
			require.Equalf(t, beforeConfig, afterConfig,
				"source chain config %q must not change", key)

			continue
		}

		targetKeys++

		// Rewind the one field the update is allowed to change. Everything else
		// must then match byte for byte.
		rewound := afterConfig
		rewound.OnRampAddresses = beforeConfig.OnRampAddresses
		require.Equal(t, beforeConfig, rewound,
			"the update must change OnRampAddresses and nothing else")
	}
	require.Equal(t, 1, targetKeys, "expected exactly one config for the target source chain")
}

// sourceChainSelectorFromKey parses a GlobalConfig source chain config map key.
// The ledger returns the NUMERIC key with a trailing dot.
func sourceChainSelectorFromKey(t *testing.T, key types.NUMERIC) uint64 {
	t.Helper()

	parsed, err := strconv.ParseUint(strings.TrimSuffix(string(key), "."), 10, 64)
	require.NoErrorf(t, err, "parse source chain selector key %q", key)

	return parsed
}

// applyDestChainConfig writes an outbound lane config, so the blast-radius
// assertion has outbound state to protect.
func applyDestChainConfig(
	t *testing.T,
	bundle cld_ops.Bundle,
	cantonChain *canton.Chain,
	gcRaw contracts.RawInstanceAddress,
	args core.DestChainConfigArgs,
) {
	t.Helper()

	_, err := cld_ops.ExecuteOperation(bundle, global_config.ApplyDestChainConfigUpdates, *cantonChain, opcontract.ChoiceInput[core.ApplyDestChainConfigUpdates]{
		InstanceAddress:    gcRaw.InstanceAddress(),
		RawInstanceAddress: gcRaw.String(),
		MCMSEnabled:        false,
		Args: core.ApplyDestChainConfigUpdates{
			DestChainConfigUpdates: []core.DestChainConfigArgs{args},
		},
	})
	require.NoError(t, err, "apply dest chain config")
}

// paddedOnRamp converts an operator-facing OnRamp address into the padded form
// that Canton stores on the ledger.
func paddedOnRamp(t *testing.T, addr string) types.TEXT {
	t.Helper()

	raw, err := hex.DecodeString(strings.TrimPrefix(addr, "0x"))
	require.NoError(t, err)
	padded := make([]byte, 32)
	copy(padded[32-len(raw):], raw)

	return types.TEXT(hex.EncodeToString(padded))
}

// hexOnRamps renders the reader's result in the operator-facing form, so a
// failure message names the address instead of a byte slice.
func hexOnRamps(onRamps [][]byte) []string {
	out := make([]string, 0, len(onRamps))
	for _, onRamp := range onRamps {
		out = append(out, "0x"+hex.EncodeToString(onRamp))
	}

	return out
}
