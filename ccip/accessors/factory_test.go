package accessors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccessorUnsupportedCapabilities(t *testing.T) {
	t.Parallel()
	accessor := newAccessor(nil)

	feeQuoterReader, err := accessor.FeeQuoterReader()
	require.Nil(t, feeQuoterReader)
	require.Error(t, err)

	feeQuoterWriter, err := accessor.FeeQuoterWriter()
	require.Nil(t, feeQuoterWriter)
	require.Error(t, err)

	gasPriceReader, err := accessor.GasPriceReader()
	require.Nil(t, gasPriceReader)
	require.Error(t, err)
}
