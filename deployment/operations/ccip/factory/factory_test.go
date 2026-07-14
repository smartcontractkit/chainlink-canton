package factory

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	executorbindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ratelimiter"
)

// Byte-for-byte the DeployExecutor operationData from the prod_testnet proposal
// 1781126292119391259-...-from-factory_mcms_timelock_proposal_0.json, which was produced
// by the previous hand-rolled encoder:
// instanceId "executor-wwbhq" + owner + int64(10) + finality tag 0x00 + bool false.
const deployExecutorProposalHex = "0e6578656375746f722d77776268714f636369704f776e65723a3a3132323065333832663465353762303831356536626537333730303665333831653662376465343438653036626430333365636536646634393830313738373966353531000000000000000a0000"

// Wire format reference: encodeRequestedFinality/decodeRequestedFinalityAt in CCIP/Codec.daml
// (finality variant = uint8 tag 0x00/0x01, or 0x02 + int64 block depth).
func TestEncodeDeployExecutor(t *testing.T) {
	t.Parallel()

	blockDepth := types.INT64(12)
	prefix := deployExecutorProposalHex[:len(deployExecutorProposalHex)-4]

	tests := []struct {
		name     string
		finality ccipcodec.FinalityConfig
		enabled  types.BOOL
		wantHex  string
	}{
		{
			name:     "wait for finality matches previously emitted proposal bytes",
			finality: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
			enabled:  false,
			wantHex:  deployExecutorProposalHex,
		},
		{
			name:     "wait for safe",
			finality: ccipcodec.FinalityConfig{WaitForSafe: &types.UNIT{}},
			enabled:  true,
			wantHex:  prefix + "0101",
		},
		{
			name:     "block depth carries int64 payload",
			finality: ccipcodec.FinalityConfig{BlockDepth: &blockDepth},
			enabled:  false,
			wantHex:  prefix + "02000000000000000c00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := encodeDeployExecutor(factorybindings.DeployExecutor{
				Contract: executorbindings.Executor{
					InstanceId:    "executor-wwbhq",
					Owner:         "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551",
					MaxCCVsPerMsg: 10,
					DynamicConfig: executorbindings.DynamicConfig{
						AllowedFinalityConfig: tt.finality,
						CcvAllowlistEnabled:   tt.enabled,
					},
				},
			})
			require.NoError(t, err)
			assert.Equal(t, "DeployExecutor", encoded.Choice)
			assert.Equal(t, tt.wantHex, encoded.OperationData)
		})
	}
}

// The generated params struct must round-trip its own encoding so MCMS proposal
// decoding (mcms/sdk/canton) can verify candidates by re-encoding.
func TestDeployExecutorParamsHexRoundTrip(t *testing.T) {
	t.Parallel()

	original := factorybindings.DeployExecutorParams{
		InstanceId:            "executor-wwbhq",
		Owner:                 "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551",
		MaxCCVsPerMsg:         10,
		AllowedFinalityConfig: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
		CcvAllowlistEnabled:   false,
	}

	encoded, err := encodeDeployExecutor(factorybindings.DeployExecutor{
		Contract: executorbindings.Executor{
			InstanceId:    original.InstanceId,
			Owner:         original.Owner,
			MaxCCVsPerMsg: original.MaxCCVsPerMsg,
			DynamicConfig: executorbindings.DynamicConfig{
				AllowedFinalityConfig: original.AllowedFinalityConfig,
				CcvAllowlistEnabled:   original.CcvAllowlistEnabled,
			},
		},
	})
	require.NoError(t, err)

	var decoded factorybindings.DeployExecutorParams
	require.NoError(t, decoded.UnmarshalHex(encoded.OperationData))
	assert.Equal(t, original, decoded)
}

// TestDeployRateLimiterParams_EnumEncoding verifies that RateLimitDirection and RateLimitMode
// are encoded as single ordinal bytes (0x00/0x01), matching the Daml MCMS codec wire format
// (decodeRateLimitDirectionAt / decodeRateLimitModeAt in contracts/ccip/factory/daml/CCIP/Factory.daml).
// Before the fix, go-daml encoded enums as length-prefixed constructor-name strings (e.g. 28 bytes
// for "RateLimitDirection_Outbound"), which the ledger's decodeUint8 step rejects.
func TestDeployRateLimiterParams_EnumEncoding(t *testing.T) {
	t.Parallel()

	params := factorybindings.DeployRateLimiterParams{
		InstanceId:          "rl-out-default",
		PoolInstanceId:      "pool-1",
		PoolOwner:           "owner::abc123",
		RemoteChainSelector: "16015286601757825753",
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Outbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_DefaultFinality,
		IsEnabled:           true,
		Capacity:            "1000",
		Rate:                "10",
	}

	encoded, err := params.MarshalHex()
	require.NoError(t, err)

	raw := []byte(encoded)
	// Walk the encoded bytes to find Direction and Mode positions:
	// TEXT instanceId:     1 + 14 = 15
	// TEXT poolInstanceId: 1 + 6  = 7  → offset 22
	// PARTY poolOwner:     1 + 13 = 14 → offset 36
	// NUMERIC selector:    1 + 20 = 21 → offset 57
	// Direction: offset 57 must be 0x00 (Outbound)
	// Mode:      offset 58 must be 0x00 (DefaultFinality)
	// Byte layout: 1+14 instanceId | 1+6 poolInstanceId | 1+13 poolOwner | 1+20 chainSelector | dir | mode | ...
	// Direction at offset 57 (15+7+14+21), Mode at offset 58.
	require.Greater(t, len(raw), 59, "encoded too short")
	assert.Equal(t, byte(0x00), raw[57], "Direction Outbound must be byte 0x00")
	assert.Equal(t, byte(0x00), raw[58], "Mode DefaultFinality must be byte 0x00")

	// Inbound / CustomFinality → 0x01
	params.Direction = ratelimiter.RateLimitDirectionRateLimitDirection_Inbound
	params.Mode = ratelimiter.RateLimitModeRateLimitMode_CustomFinality
	encoded, err = params.MarshalHex()
	require.NoError(t, err)
	raw = []byte(encoded)
	assert.Equal(t, byte(0x01), raw[57], "Direction Inbound must be byte 0x01")
	assert.Equal(t, byte(0x01), raw[58], "Mode CustomFinality must be byte 0x01")
}

func TestDeployRateLimiterParams_HexRoundTrip(t *testing.T) {
	t.Parallel()

	original := factorybindings.DeployRateLimiterParams{
		InstanceId:          "rl-inbound",
		PoolInstanceId:      "pool-2",
		PoolOwner:           "owner::def456",
		RemoteChainSelector: "16015286601757825753",
		Direction:           ratelimiter.RateLimitDirectionRateLimitDirection_Inbound,
		Mode:                ratelimiter.RateLimitModeRateLimitMode_CustomFinality,
		IsEnabled:           false,
		Capacity:            "500",
		Rate:                "5",
	}

	// MarshalHex returns raw bytes (not a hex string); hex-encode for UnmarshalHex (same as executor test).
	encoded, err := original.MarshalHex()
	require.NoError(t, err)

	var decoded factorybindings.DeployRateLimiterParams
	require.NoError(t, decoded.UnmarshalHex(hex.EncodeToString([]byte(encoded))))
	assert.Equal(t, original, decoded)
}
