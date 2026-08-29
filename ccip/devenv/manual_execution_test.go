package devenv

import (
	"testing"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

func TestHashInstanceAddress_rawInstanceAddress(t *testing.T) {
	t.Parallel()

	raw := "test-verifier@ccvOwner::1220abcd"
	addr, err := hashInstanceAddress(protocol.UnknownAddress(raw))
	require.NoError(t, err)

	expectedRaw, err := contracts.RawInstanceAddressFromString(raw)
	require.NoError(t, err)
	require.Equal(t, protocol.UnknownAddress(expectedRaw.InstanceAddress().Bytes()), addr)
}

func TestHashInstanceAddress_alreadyHashed(t *testing.T) {
	t.Parallel()

	hashed := protocol.UnknownAddress(contracts.HexToInstanceAddress("0xec1e288bcf8bbf034ac2d31b67f9b15a3f1f828d086c5b9d8fc2866129cd02fe").Bytes())
	addr, err := hashInstanceAddress(hashed)
	require.NoError(t, err)
	require.Equal(t, hashed, addr)
}

func TestHashInstanceAddress_empty(t *testing.T) {
	t.Parallel()

	_, err := hashInstanceAddress(nil)
	require.Error(t, err)
}
