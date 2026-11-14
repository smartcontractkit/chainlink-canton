package bindings_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	registryratelimiter "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
)

// The factory decodes DeployRegistryRateLimiterParams with decodeDeployRegistryRateLimiterParams in
// contracts/ccip/factory/daml/CCIP/FactoryV2.daml. Direction and mode are read by
// decodeRegistryRateLimitDirectionAt and decodeRegistryRateLimitModeAt, which only accept the
// single tag bytes 00 and 01.
func TestDeployRegistryRateLimiterParams_TagBytes(t *testing.T) {
	const (
		instanceID     = "rl-1"
		poolInstanceID = "pool-1"
		selector       = "16015286601757825753"
	)
	// Text and Numeric fields are encoded as [uint8 utf8 length][utf8 bytes].
	lenPrefixed := func(s string) string {
		return hex.EncodeToString(append([]byte{byte(len(s))}, s...))
	}
	prefix := lenPrefixed(instanceID) + lenPrefixed(poolInstanceID) + lenPrefixed(selector)
	// isEnabled=true, capacity "1000", rate "10", empty observers list.
	tail := "01" + lenPrefixed("1000") + lenPrefixed("10") + "00"

	tests := []struct {
		name      string
		direction registryratelimiter.RateLimitDirection
		mode      registryratelimiter.RateLimitMode
		wantTags  string
	}{
		{
			name:      "Outbound_DefaultFinality",
			direction: "RateLimitDirection_Outbound",
			mode:      "RateLimitMode_DefaultFinality",
			wantTags:  "0000",
		},
		{
			name:      "Outbound_CustomFinality",
			direction: "RateLimitDirection_Outbound",
			mode:      "RateLimitMode_CustomFinality",
			wantTags:  "0001",
		},
		{
			name:      "Inbound_DefaultFinality",
			direction: "RateLimitDirection_Inbound",
			mode:      "RateLimitMode_DefaultFinality",
			wantTags:  "0100",
		},
		{
			name:      "Inbound_CustomFinality",
			direction: "RateLimitDirection_Inbound",
			mode:      "RateLimitMode_CustomFinality",
			wantTags:  "0101",
		},
		{
			// An empty direction encodes as 0xFF, which the Daml decoder rejects.
			name:      "EmptyDirection",
			direction: "",
			mode:      "RateLimitMode_DefaultFinality",
			wantTags:  "ff00",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := factory.DeployRegistryRateLimiterParams{
				InstanceId:          instanceID,
				PoolInstanceId:      poolInstanceID,
				RemoteChainSelector: selector,
				Direction:           tc.direction,
				Mode:                tc.mode,
				IsEnabled:           true,
				Capacity:            "1000",
				Rate:                "10",
				Observers:           []types.PARTY{},
			}

			raw, err := params.MarshalHex()
			require.NoError(t, err)
			// MarshalHex returns raw bytes; the MCMS transport hex-encodes them.
			encoded := hex.EncodeToString([]byte(raw))

			require.Greaterf(t, len(encoded), len(prefix)+6, "encoding too short: %q", encoded)
			require.Equalf(t, prefix, encoded[:len(prefix)], "fields before direction changed: %q", encoded)
			// Direction tag, mode tag, then the isEnabled bool byte.
			require.Equal(t, tc.wantTags+"01", encoded[len(prefix):len(prefix)+6],
				"direction and mode must be single tag bytes followed by the bool byte")
			require.Equal(t, prefix+tc.wantTags+tail, encoded)
		})
	}
}
