package cantonops

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestFinalityConfigEqual(t *testing.T) {
	t.Parallel()

	waitForFinality := core.FinalityConfig{WaitForFinality: &types.UNIT{}}
	waitForSafe := core.FinalityConfig{WaitForSafe: &types.UNIT{}}
	blockDepth1 := core.FinalityConfig{BlockDepth: new(types.INT64(1))}
	blockDepth5 := core.FinalityConfig{BlockDepth: new(types.INT64(5))}

	require.True(t, finalityConfigEqual(waitForFinality, waitForFinality))
	require.True(t, finalityConfigEqual(blockDepth1, blockDepth1))
	require.False(t, finalityConfigEqual(waitForFinality, waitForSafe))
	require.False(t, finalityConfigEqual(waitForFinality, blockDepth1))
	require.False(t, finalityConfigEqual(blockDepth1, blockDepth5))
}

func TestReceiverInstanceID(t *testing.T) {
	t.Parallel()

	require.Equal(t, "ccipreceiver-WaitForFinality", receiverInstanceID(core.FinalityConfig{
		WaitForFinality: &types.UNIT{},
	}))
	require.Equal(t, "ccipreceiver-WaitForSafe", receiverInstanceID(core.FinalityConfig{
		WaitForSafe: &types.UNIT{},
	}))
	require.Equal(t, "ccipreceiver-BlockDepth-1", receiverInstanceID(core.FinalityConfig{
		BlockDepth: new(types.INT64(1)),
	}))
}

func TestReceiverFinalityLabel(t *testing.T) {
	t.Parallel()

	require.Equal(t, "WaitForFinality", receiverFinalityLabel(core.FinalityConfig{
		WaitForFinality: &types.UNIT{},
	}))
	require.Equal(t, "BlockDepth(1)", receiverFinalityLabel(core.FinalityConfig{
		BlockDepth: new(types.INT64(1)),
	}))
}

func TestReceiverRequiredCCVConfigured(t *testing.T) {
	t.Parallel()

	ccv := contracts.RawInstanceAddress("committeeverifier-tqkny@ccvOwner::1220abc")
	recv := &receiver.CCIPReceiver{
		RequiredCCVs: []chainlinkapi.RawInstanceAddress{
			{Unpack: types.TEXT(ccv.String())},
		},
	}

	require.True(t, receiverRequiredCCVConfigured(recv, ccv))
	require.False(t, receiverRequiredCCVConfigured(recv, contracts.RawInstanceAddress("other@ccvOwner::1220abc")))
	require.False(t, receiverRequiredCCVConfigured(&receiver.CCIPReceiver{}, ccv))
}
