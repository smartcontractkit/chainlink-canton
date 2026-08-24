// Package finality parses CCIP finality settings for the demo CLI.
package finality

import (
	"fmt"
	"strconv"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
)

// Parsed holds the EVM extraArgs finality and matching Canton receiver config.
type Parsed struct {
	EVM      protocol.Finality
	Receiver ccipcodec.FinalityConfig
	Label    string
}

// Parse parses --finality values:
//   - finality — wait for full source-chain finality (slowest, default)
//   - safe     — wait for Ethereum safe head
//   - N        — wait for N block confirmations (1–65535; 1 is fastest)
func Parse(s string) (Parsed, error) {
	switch s {
	case "finality":
		return Parsed{
			EVM:      protocol.FinalityWaitForFinality,
			Receiver: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
			Label:    "WaitForFinality",
		}, nil
	case "safe":
		return Parsed{
			EVM:      protocol.FinalityWaitForSafe,
			Receiver: ccipcodec.FinalityConfig{WaitForSafe: &types.UNIT{}},
			Label:    "WaitForSafe",
		}, nil
	default:
		n, err := strconv.ParseUint(s, 10, 16)
		if err != nil || n < 1 || n > 65535 {
			return Parsed{}, fmt.Errorf("invalid finality %q (finality|safe|1-65535)", s)
		}
		depth := uint16(n)

		return Parsed{
			EVM:      protocol.NewFinality().WithBlockDepth(depth),
			Receiver: ccipcodec.FinalityConfig{BlockDepth: new(types.INT64(int64(depth)))},
			Label:    fmt.Sprintf("BlockDepth(%d)", depth),
		}, nil
	}
}
