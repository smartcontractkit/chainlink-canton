package adapters

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
)

const (
	testEVMChainSelector    = uint64(16015286601757825753) // Sepolia
	testCantonChainSelector = uint64(10109143320554840099) // Canton testnet
)

// evmOnRampWire is an EVM OnRamp address in the form the source chain writes it
// into its messages: abi.encode(address), 32 bytes, left zero-padded. This is
// what the caller supplies and what the ledger stores.
const evmOnRampWire = "0x0000000000000000000000001111111111111111111111111111111111111111"

// evmOnRampLedger is evmOnRampWire in the ledger's own form: the same bytes, hex
// encoded without a 0x prefix.
const evmOnRampLedger = "0000000000000000000000001111111111111111111111111111111111111111"

// evmOnRampContract is the 20-byte contract address. The adapter must reject it:
// it is not what the source chain puts on the wire.
const evmOnRampContract = "0x1111111111111111111111111111111111111111"

func TestParseCantonOffRampSourceOnRamps(t *testing.T) {
	t.Parallel()

	second := "0x0000000000000000000000002222222222222222222222222222222222222222"

	got, err := parseCantonOffRampSourceOnRamps([]string{
		evmOnRampWire,
		"  " + evmOnRampWire + "  ", // duplicate after trimming whitespace
		second,
	})
	require.NoError(t, err)
	require.Equal(t, []types.TEXT{
		types.TEXT(evmOnRampLedger),
		types.TEXT("0000000000000000000000002222222222222222222222222222222222222222"),
	}, got, "the wire bytes are stored verbatim, and duplicates are dropped")
}

func TestParseCantonOffRampSourceOnRamps_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		addrs []string
		msg   string
	}{
		{name: "empty input", addrs: nil},
		{
			// The reader and the changesets both work in the wire encoding, so a
			// bare contract address must not be silently padded.
			name:  "20-byte contract address is rejected",
			addrs: []string{evmOnRampContract},
			msg:   "must be 32 bytes, got 20",
		},
		{name: "missing 0x prefix", addrs: []string{evmOnRampLedger}},
		{name: "not hex", addrs: []string{"0xzz"}},
		{name: "empty address", addrs: []string{"0x"}},
		{
			name:  "too wide",
			addrs: []string{"0x" + strings.Repeat("11", cantonOnRampAddressBytes+1)},
			msg:   "must be 32 bytes, got 33",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseCantonOffRampSourceOnRamps(test.addrs)
			require.Error(t, err)
			if test.msg != "" {
				require.Contains(t, err.Error(), test.msg)
			}
		})
	}
}

func TestDecodeCantonOnRampAddresses(t *testing.T) {
	t.Parallel()

	got, err := decodeCantonOnRampAddresses([]types.TEXT{types.TEXT(evmOnRampLedger)})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Len(t, got[0], cantonOnRampAddressBytes,
		"the reader must preserve the wire width, not reduce it to the 20-byte contract address")
	require.Equal(t, evmOnRampWire, "0x"+hex.EncodeToString(got[0]))
}

// TestDecodeCantonOnRampAddresses_RoundTrip pins the property the changesets
// depend on: what the caller supplies is what the reader gives back.
func TestDecodeCantonOnRampAddresses_RoundTrip(t *testing.T) {
	t.Parallel()

	entries, err := parseCantonOffRampSourceOnRamps([]string{evmOnRampWire})
	require.NoError(t, err)

	got, err := decodeCantonOnRampAddresses(entries)
	require.NoError(t, err)

	want, err := hex.DecodeString(strings.TrimPrefix(evmOnRampWire, "0x"))
	require.NoError(t, err)
	require.Equal(t, [][]byte{want}, got)
}

func TestDecodeCantonOnRampAddresses_InvalidHex(t *testing.T) {
	t.Parallel()

	_, err := decodeCantonOnRampAddresses([]types.TEXT{"not-hex"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "onRampAddresses[0]")
}

func TestCantonOnRampSetsEqual(t *testing.T) {
	t.Parallel()

	legacy := types.TEXT(evmOnRampLedger)
	current := types.TEXT("0000000000000000000000002222222222222222222222222222222222222222")

	tests := []struct {
		name    string
		current []types.TEXT
		desired []types.TEXT
		want    bool
	}{
		{
			name:    "equal ignoring order",
			current: []types.TEXT{legacy, current},
			desired: []types.TEXT{current, legacy},
			want:    true,
		},
		{
			name:    "equal ignoring hex case and 0x prefix",
			current: []types.TEXT{types.TEXT("0x" + strings.ToUpper(evmOnRampLedger))},
			desired: []types.TEXT{legacy},
			want:    true,
		},
		{
			// Width is significant on the wire, so an unpadded address is a
			// different onramp, not the same one.
			name:    "padding is significant",
			current: []types.TEXT{types.TEXT("1111111111111111111111111111111111111111")},
			desired: []types.TEXT{legacy},
			want:    false,
		},
		{
			name:    "different length",
			current: []types.TEXT{legacy},
			desired: []types.TEXT{legacy, current},
			want:    false,
		},
		{
			name:    "different members",
			current: []types.TEXT{legacy},
			desired: []types.TEXT{current},
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.want, cantonOnRampSetsEqual(test.current, test.desired))
		})
	}
}

func TestFindCantonSourceChainConfig(t *testing.T) {
	t.Parallel()

	want := core.SourceChainConfig2{
		IsEnabled:       types.BOOL(true),
		OnRampAddresses: []types.TEXT{types.TEXT(evmOnRampLedger)},
	}
	// The ledger returns the NUMERIC map key with a trailing dot. The write path
	// formats it without one, so the key is parsed rather than built.
	configs := map[types.NUMERIC]core.SourceChainConfig2{
		types.NUMERIC("16015286601757825753."): want,
	}

	got, found, err := findCantonSourceChainConfig(configs, testEVMChainSelector)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, want, got)

	_, found, err = findCantonSourceChainConfig(configs, testCantonChainSelector)
	require.NoError(t, err)
	require.False(t, found)

	_, _, err = findCantonSourceChainConfig(map[types.NUMERIC]core.SourceChainConfig2{
		types.NUMERIC("not-a-number"): want,
	}, testEVMChainSelector)
	require.Error(t, err)
}
