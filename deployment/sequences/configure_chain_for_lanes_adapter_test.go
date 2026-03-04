package sequences

import (
	"testing"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestToCantonConfigureChainForLanesInput_Valid(t *testing.T) {
	in := ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: 1,
		Router:        "0x01",
		OnRamp:        "0x02",
		FeeQuoter:     "0x03",
		OffRamp:       "0x04",
		RemoteChains: map[uint64]ccipadapters.RemoteChainConfig[[]byte, string]{
			2: {
				AllowTrafficFrom: true,
				OnRamps:          [][]byte{{0xaa}},
				OffRamp:          []byte{0xbb},
				DefaultExecutor:  "0x05",
			},
		},
	}

	out, err := toCantonConfigureChainForLanesInput(in)
	require.NoError(t, err)
	require.Equal(t, in.ChainSelector, out.ChainSelector)
	require.Equal(t, contracts.HexToInstanceAddress(in.Router).Hex(), out.GlobalConfig.Hex())
	require.Equal(t, contracts.HexToInstanceAddress(in.OnRamp).Hex(), out.OnRamp.Hex())
	require.Equal(t, contracts.HexToInstanceAddress(in.FeeQuoter).Hex(), out.FeeQuoter.Hex())
	require.Equal(t, contracts.HexToInstanceAddress(in.OffRamp).Hex(), out.OffRamp.Hex())
	require.Contains(t, out.RemoteChains, uint64(2))
	require.Equal(t, []byte{0xbb}, out.RemoteChains[2].OffRamp)
}

func TestToCantonConfigureChainForLanesInput_RequiresCoreFields(t *testing.T) {
	_, err := toCantonConfigureChainForLanesInput(ccipadapters.ConfigureChainForLanesInput{
		ChainSelector: 1,
		OnRamp:        "0x02",
		FeeQuoter:     "0x03",
		OffRamp:       "0x04",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Router")
}

