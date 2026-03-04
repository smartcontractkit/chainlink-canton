package changesets

import (
	"testing"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

func TestToCCIPConfigureChainForLanesInput_MapsCoreFields(t *testing.T) {
	in := sequences.ConfigureChainForLanesInput{
		ChainSelector: 77,
		GlobalConfig:  contracts.HexToInstanceAddress("0x01"),
		OnRamp:        contracts.HexToInstanceAddress("0x02"),
		FeeQuoter:     contracts.HexToInstanceAddress("0x03"),
		OffRamp:       contracts.HexToInstanceAddress("0x04"),
		RemoteChains: map[uint64]ccipadapters.RemoteChainConfig[[]byte, contracts.RawInstanceAddress]{
			88: {
				AllowTrafficFrom: true,
				OnRamps:          [][]byte{{0xaa}},
				OffRamp:          []byte{0xbb},
				DefaultExecutor:  contracts.RawInstanceAddress("0x05"),
			},
		},
	}

	out := sequences.ToCCIPConfigureChainForLanesInput(in)
	require.Equal(t, in.ChainSelector, out.ChainSelector)
	require.Equal(t, in.GlobalConfig.Hex(), out.Router)
	require.Equal(t, in.OnRamp.Hex(), out.OnRamp)
	require.Equal(t, in.FeeQuoter.Hex(), out.FeeQuoter)
	require.Equal(t, in.OffRamp.Hex(), out.OffRamp)
	require.Contains(t, out.RemoteChains, uint64(88))
	require.Equal(t, []byte{0xbb}, out.RemoteChains[88].OffRamp)
	require.Equal(t, "0x05", out.RemoteChains[88].DefaultExecutor)
}

