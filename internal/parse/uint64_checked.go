package parse

import (
	"fmt"
	"math/big"
)

// Uint64Checked parses a uint64 from a string that may have been serialized from a
// Daml NUMERIC.
//
// Daml NUMERIC values are decimal strings with arbitrary precision, and may be
// represented with a trailing decimal point (for example "12345."). In this
// codebase we only accept integer-valued NUMERICs when converting into uint64.
//
// The parsing is strict:
// - the value must parse as a rational number
// - it must be an integer (no fractional component)
// - it must fit into an unsigned 64-bit integer
func Uint64Checked(s string) (uint64, error) {
	// Since uint64s are represented by Numeric 0, which will be serialized as "1.0",
	// cannot use *big.Int directly - instead, parse as big.Rat and then convert to int.
	val, ok := new(big.Rat).SetString(s)
	if !ok {
		return 0, fmt.Errorf("failed to parse uint64: %s", s)
	}
	if !val.IsInt() {
		return 0, fmt.Errorf("is not an int: %s", s)
	}

	num := val.Num()
	if !num.IsUint64() {
		return 0, fmt.Errorf("is not uint64: %s", s)
	}

	return num.Uint64(), nil
}
