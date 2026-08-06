package registry

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func parseDecimal(label, s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, fmt.Errorf("%s: parse %q: %w", label, s, err)
	}

	return d, nil
}
