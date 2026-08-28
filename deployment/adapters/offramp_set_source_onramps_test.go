package adapters

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"

	// Registers the EVM chain family adapter, which sourceOnRampAddressBytes reads
	// to resolve the EVM address width.
	_ "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/adapters"
)

const (
	testEVMChainSelector    = uint64(16015286601757825753) // Sepolia
	testCantonChainSelector = uint64(10109143320554840099) // Canton testnet
)

// evmOnRampHex is a 20-byte EVM address in the form the changesets pass in.
const evmOnRampHex = "0x1111111111111111111111111111111111111111"

// evmOnRampPadded is the same address in the form Canton stores on the ledger.
const evmOnRampPadded = "0000000000000000000000001111111111111111111111111111111111111111"

func TestPadAndTrimCantonOnRampAddress_EVMRoundTrip(t *testing.T) {
	t.Parallel()

	raw, err := hex.DecodeString(strings.TrimPrefix(evmOnRampHex, "0x"))
	require.NoError(t, err)

	padded := padCantonOnRampAddress(raw)
	require.Equal(t, types.TEXT(evmOnRampPadded), padded)

	decoded, err := hex.DecodeString(string(padded))
	require.NoError(t, err)
	require.Equal(t, raw, trimCantonOnRampAddress(decoded, 20))
}

func TestPadAndTrimCantonOnRampAddress_CantonRoundTrip(t *testing.T) {
	t.Parallel()

	raw := make([]byte, cantonOnRampAddressBytes)
	for i := range raw {
		raw[i] = byte(i + 1)
	}

	padded := padCantonOnRampAddress(raw)
	decoded, err := hex.DecodeString(string(padded))
	require.NoError(t, err)
	require.Equal(t, raw, decoded)

	// A 32-byte address must survive a trim request for a shorter width.
	require.Equal(t, raw, trimCantonOnRampAddress(decoded, 20))
}

func TestTrimCantonOnRampAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     []byte
		wantLen int
		want    []byte
	}{
		{
			name:    "trims zero padding",
			raw:     []byte{0, 0, 1, 2},
			wantLen: 2,
			want:    []byte{1, 2},
		},
		{
			name:    "keeps value when the prefix is not zero",
			raw:     []byte{0, 9, 1, 2},
			wantLen: 2,
			want:    []byte{0, 9, 1, 2},
		},
		{
			name:    "keeps value when it is already short",
			raw:     []byte{1, 2},
			wantLen: 4,
			want:    []byte{1, 2},
		},
		{
			name:    "keeps value when no width is wanted",
			raw:     []byte{0, 0, 1, 2},
			wantLen: 0,
			want:    []byte{0, 0, 1, 2},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.want, trimCantonOnRampAddress(test.raw, test.wantLen))
		})
	}
}

func TestSourceOnRampAddressBytes(t *testing.T) {
	t.Parallel()

	got, err := sourceOnRampAddressBytes(testEVMChainSelector)
	require.NoError(t, err)
	require.Equal(t, 20, got)

	got, err = sourceOnRampAddressBytes(testCantonChainSelector)
	require.NoError(t, err)
	require.Equal(t, int(new(CantonChainFamilyAdapter).GetAddressBytesLength()), got)

	_, err = sourceOnRampAddressBytes(1)
	require.Error(t, err)
}

func TestDecodeCantonOnRampAddresses_EVMSource(t *testing.T) {
	t.Parallel()

	addressBytes, err := sourceOnRampAddressBytes(testEVMChainSelector)
	require.NoError(t, err)

	got, err := decodeCantonOnRampAddresses([]types.TEXT{types.TEXT(evmOnRampPadded)}, addressBytes)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Len(t, got[0], 20)
	require.Equal(t, evmOnRampHex, "0x"+hex.EncodeToString(got[0]))
}

func TestDecodeCantonOnRampAddresses_InvalidHex(t *testing.T) {
	t.Parallel()

	_, err := decodeCantonOnRampAddresses([]types.TEXT{"not-hex"}, 20)
	require.Error(t, err)
	require.Contains(t, err.Error(), "onRampAddresses[0]")
}

func TestParseCantonOffRampSourceOnRamps(t *testing.T) {
	t.Parallel()

	got, err := parseCantonOffRampSourceOnRamps([]string{
		evmOnRampHex,
		"  " + evmOnRampHex + "  ", // duplicate after trim
		"0x2222222222222222222222222222222222222222",
	})
	require.NoError(t, err)
	require.Equal(t, []types.TEXT{
		types.TEXT(evmOnRampPadded),
		types.TEXT("0000000000000000000000002222222222222222222222222222222222222222"),
	}, got)
}

func TestParseCantonOffRampSourceOnRamps_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		addrs []string
	}{
		{name: "empty input", addrs: nil},
		{name: "missing 0x prefix", addrs: []string{evmOnRampPadded}},
		{name: "not hex", addrs: []string{"0xzz"}},
		{name: "empty address", addrs: []string{"0x"}},
		{
			name:  "too wide",
			addrs: []string{"0x" + strings.Repeat("11", cantonOnRampAddressBytes+1)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseCantonOffRampSourceOnRamps(test.addrs)
			require.Error(t, err)
		})
	}
}

func TestCantonOnRampSetsEqual(t *testing.T) {
	t.Parallel()

	legacy := types.TEXT(evmOnRampPadded)
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
			name:    "equal ignoring hex case and prefix",
			current: []types.TEXT{types.TEXT("0x" + strings.ToUpper(evmOnRampPadded))},
			desired: []types.TEXT{legacy},
			want:    true,
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
		OnRampAddresses: []types.TEXT{types.TEXT(evmOnRampPadded)},
	}
	// The ledger returns the NUMERIC map key with a trailing dot. The write path
	// formats it without one, so both forms must resolve.
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
