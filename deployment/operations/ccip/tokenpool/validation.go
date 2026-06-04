package tokenpool

import (
	"errors"
	"fmt"
)

// Bounds match CCIP.Math minTokenDecimals / maxTokenDecimals on-ledger.
const (
	MinTokenDecimals int64 = 0
	MaxTokenDecimals int64 = 37
)

// ValidateTokenDecimals checks pool deploy-time decimal precision.
func ValidateTokenDecimals(decimals int64) error {
	if decimals < MinTokenDecimals {
		return errors.New("decimals cannot be negative")
	}
	if decimals > MaxTokenDecimals {
		return fmt.Errorf("decimals %d exceeds maximum %d", decimals, MaxTokenDecimals)
	}

	return nil
}
