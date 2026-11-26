package accessors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

func TestCantonAccessorCCIPAPIParser(t *testing.T) {
	t.Parallel()

	a := &accessor{}
	parser, err := a.CCIPAPIParser()
	require.NoError(t, err)
	require.Same(t, a, parser)

	raw, err := contracts.RawInstanceAddressFromString("token@owner")
	require.NoError(t, err)
	got, err := parser.DecodeAddress(raw.String())
	require.NoError(t, err)
	require.Equal(t, raw.InstanceAddress().Bytes(), []byte(got))

	hashed := contracts.HexToInstanceAddress("0x01")
	got, err = parser.DecodeAddress(hashed.String())
	require.NoError(t, err)
	require.Equal(t, hashed.Bytes(), []byte(got))

	for _, address := range []string{"not-an-address", "0x01", "0x" + strings.Repeat("zz", 32)} {
		_, err = parser.DecodeAddress(address)
		require.Error(t, err, address)
	}
}
