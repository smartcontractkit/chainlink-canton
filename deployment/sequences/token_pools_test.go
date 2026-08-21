package sequences

import (
	"testing"

	tokenadaptersfinality "github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/ccipcodec"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
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

func TestToCantonFinalityConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  tokenadaptersfinality.Config
		want func(t *testing.T, got ccipcodec.FinalityConfig)
	}{
		{
			name: "block depth takes precedence over wait for finality",
			cfg: tokenadaptersfinality.Config{
				WaitForFinality: true,
				BlockDepth:      1,
			},
			want: func(t *testing.T, got ccipcodec.FinalityConfig) {
				t.Helper()
				require.NotNil(t, got.BlockDepth)
				require.Equal(t, types.INT64(1), *got.BlockDepth)
				require.Nil(t, got.WaitForFinality)
			},
		},
		{
			name: "wait for safe when no block depth",
			cfg: tokenadaptersfinality.Config{
				WaitForSafe: true,
			},
			want: func(t *testing.T, got ccipcodec.FinalityConfig) {
				t.Helper()
				require.NotNil(t, got.WaitForSafe)
			},
		},
		{
			name: "wait for finality when no block depth or safe",
			cfg: tokenadaptersfinality.Config{
				WaitForFinality: true,
			},
			want: func(t *testing.T, got ccipcodec.FinalityConfig) {
				t.Helper()
				require.NotNil(t, got.WaitForFinality)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toCantonFinalityConfig(tt.cfg)
			tt.want(t, got)
		})
	}
}
