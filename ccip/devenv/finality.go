package devenv

import (
	"fmt"

	"github.com/smartcontractkit/go-daml/pkg/types"
)

const (
	receiverWaitForFinalityConfig types.TEXT = "00000000"
	receiverWaitForSafeConfig     types.TEXT = "00010000"
)

func encodeReceiverFinalityConfig(finality int64) (types.TEXT, error) {
	switch {
	case finality < 0:
		return "", fmt.Errorf("invalid finality %d: must be non-negative", finality)
	case finality == 0:
		return receiverWaitForFinalityConfig, nil
	case finality == 0x00010000:
		return receiverWaitForSafeConfig, nil
	case finality > 0xFFFF:
		return "", fmt.Errorf("invalid finality %d: max supported block depth is 65535", finality)
	default:
		// Receiver-side config uses bytes4 text with flags in the upper 16 bits.
		return types.TEXT(fmt.Sprintf("%08x", uint32(finality))), nil
	}
}
