package sequences

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInstrumentPriceKey_CantonPartyIDs(t *testing.T) {
	t.Parallel()

	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"
	dso := "DSO::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337"

	admin, id, err := parseInstrumentPriceKey(ccipOwner + ":link-token")
	require.NoError(t, err)
	require.Equal(t, ccipOwner, admin)
	require.Equal(t, "link-token", id)

	admin, id, err = parseInstrumentPriceKey(dso + ":Amulet")
	require.NoError(t, err)
	require.Equal(t, dso, admin)
	require.Equal(t, "Amulet", id)
}
