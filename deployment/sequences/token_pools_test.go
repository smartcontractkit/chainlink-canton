package sequences

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
)

func TestResolveCantonTokenPoolType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    datastore.ContractType
		wantErr bool
	}{
		{name: "empty defaults to lock release", input: "", want: lockReleasePoolType},
		{name: "lock release", input: "LockReleaseTokenPool", want: lockReleasePoolType},
		{name: "burn mint", input: "BurnMintTokenPool", want: burnMintPoolType},
		{name: "invalid", input: "UnsupportedPool", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resolveCantonTokenPoolType(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseInstrumentIDFromTokenRefLabels(t *testing.T) {
	t.Parallel()

	t.Run("parses instrument labels", func(t *testing.T) {
		t.Parallel()

		instrumentID, instrumentAdmin, err := parseInstrumentIDFromTokenRefLabels(datastore.AddressRef{
			Labels: datastore.NewLabelSet("instrument-admin:link-admin", "instrument-id:LINK"),
		})
		require.NoError(t, err)
		require.Equal(t, "link-admin", instrumentAdmin)
		require.Equal(t, splice_api_token_holding_v1.InstrumentId{
			Admin: "link-admin",
			Id:    "LINK",
		}, instrumentID)
	})

	t.Run("fails when labels are missing", func(t *testing.T) {
		t.Parallel()

		_, _, err := parseInstrumentIDFromTokenRefLabels(datastore.AddressRef{})
		require.Error(t, err)
	})
}
