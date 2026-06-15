package factory

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
)

func TestEncodeDeployRateLimiterParamsWire_EnumTagBytes(t *testing.T) {
	t.Parallel()

	opData, err := encodeDeployRateLimiterParamsWire(factorybindings.DeployRateLimiterParams{
		InstanceId:          types.TEXT("0xabc-outbound-16015286601757825753"),
		PoolInstanceId:      types.TEXT("burnminttokenpool-LINK"),
		PoolOwner:           types.PARTY("ccipOwner::1220deadbeef"),
		RemoteChainSelector: types.NUMERIC("16015286601757825753"),
		Direction:           core.RateLimitDirectionRateLimitDirection_Outbound,
		Mode:                core.RateLimitModeRateLimitMode_DefaultFinality,
		IsEnabled:           types.BOOL(true),
		Capacity:            types.NUMERIC("0"),
		Rate:                types.NUMERIC("0"),
	})
	if err != nil {
		t.Fatalf("encodeDeployRateLimiterParamsWire: %v", err)
	}

	wire, err := hex.DecodeString(opData)
	if err != nil {
		t.Fatalf("decode operation data: %v", err)
	}

	remoteEnd := findSubslice(wire, []byte("16015286601757825753"))
	if remoteEnd < 0 {
		t.Fatalf("remote chain selector not found in wire: %x", wire)
	}

	directionOffset := remoteEnd + len("16015286601757825753") + 1 // length prefix byte
	if directionOffset >= len(wire) {
		t.Fatalf("wire too short after remote selector: %x", wire)
	}
	if wire[directionOffset] != 0x00 {
		t.Fatalf("direction byte = 0x%02x, want 0x00 (Outbound)", wire[directionOffset])
	}
	if wire[directionOffset+1] != 0x00 {
		t.Fatalf("mode byte = 0x%02x, want 0x00 (DefaultFinality)", wire[directionOffset+1])
	}
}

func findSubslice(haystack, needle []byte) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		match := true
		for j := range needle {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}

	return -1
}
