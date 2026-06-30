package token_admin_registry

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	splice "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
)

func TestTokenAdminRegistryMCMSEncodingUsesParamsTypes(t *testing.T) {
	t.Parallel()

	inst := splice.InstrumentId{
		Admin: types.PARTY("party::alice"),
		Id:    types.TEXT("TOKEN"),
	}
	admin := types.PARTY("party::bob")
	pool := &core.PoolRegistration2{
		PoolOwner:      types.PARTY("party::pool"),
		PoolInstanceId: types.TEXT("pool-1"),
	}

	proposeParams := core.ProposeAdminParams{InstrumentId: inst, NewAdmin: admin}
	proposeEncoded, err := tarEncoder.ProposeAdministrator(proposeParams)
	require.NoError(t, err)
	proposeWire, err := proposeParams.MarshalHex()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString([]byte(proposeWire)), proposeEncoded.OperationData)
	require.Equal(t, "ProposeAdministrator", proposeEncoded.Choice)

	acceptParams := core.AcceptAdminParams{InstrumentId: inst}
	acceptEncoded, err := tarEncoder.AcceptAdminRole(acceptParams)
	require.NoError(t, err)
	acceptWire, err := acceptParams.MarshalHex()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString([]byte(acceptWire)), acceptEncoded.OperationData)
	require.Equal(t, "AcceptAdminRole", acceptEncoded.Choice)

	setPoolParams := core.SetPoolParams{InstrumentId: inst, TokenPool: pool}
	setPoolEncoded, err := tarEncoder.SetPoolParams(setPoolParams)
	require.NoError(t, err)
	setPoolWire, err := setPoolParams.MarshalHex()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString([]byte(setPoolWire)), setPoolEncoded.OperationData)
	require.Equal(t, "SetPool", setPoolEncoded.Choice)
}
