package tokenpool

import (
	"context"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
)

func TestNewStaticFactoryResolver(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when no factories configured", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, newStaticFactoryResolver(nil, nil))
	})

	t.Run("returns configured transfer and burn mint factories", func(t *testing.T) {
		t.Parallel()

		transferFactoryID := "transfer-factory-cid"
		burnMintFactoryID := "burnmint-factory-cid"
		resolver := newStaticFactoryResolver(&transferFactoryID, &burnMintFactoryID)
		require.NotNil(t, resolver)

		factories, err := resolver(context.Background(), splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY("admin"),
			Id:    "LINK",
		})
		require.NoError(t, err)
		require.Equal(t, transferFactoryID, factories.TransferFactory)
		require.NotNil(t, factories.BurnMintFactory)
		require.Equal(t, burnMintFactoryID, *factories.BurnMintFactory)
		require.Equal(t, map[string]any{"values": map[string]any{}}, factories.ChoiceContext)
		require.Nil(t, factories.DisclosedContracts)
	})

	t.Run("allows burn mint only configuration", func(t *testing.T) {
		t.Parallel()

		burnMintFactoryID := "burnmint-factory-cid"
		resolver := newStaticFactoryResolver(nil, &burnMintFactoryID)
		require.NotNil(t, resolver)

		factories, err := resolver(context.Background(), splice_api_token_holding_v1.InstrumentId{})
		require.NoError(t, err)
		require.Empty(t, factories.TransferFactory)
		require.NotNil(t, factories.BurnMintFactory)
		require.Equal(t, burnMintFactoryID, *factories.BurnMintFactory)
	})
}
