package types

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type (
	PARTY       string
	TEXT        string
	INT64       int64
	BOOL        bool
	DECIMAL     *big.Int
	NUMERIC     string
	DATE        time.Time
	TIMESTAMP   time.Time
	UNIT        struct{}
	LIST        []string
	MAP         map[string]interface{}
	OPTIONAL    *interface{}
	GENMAP      map[string]interface{}
	TEXTMAP     map[string]interface{}
	CONTRACT_ID string
	RELTIME     time.Duration
	SET         []interface{}
	TUPLE2      struct {
		First  interface{}
		Second interface{}
	}
)

func NewNumericFromDecimal(d decimal.Decimal) NUMERIC {
	// Convert decimal to string with 10 decimal places precision
	scaled := d.Shift(10)
	return NUMERIC(scaled.String())
}

// VARIANT represents a DAML variant/union type
type VARIANT interface {
	GetVariantTag() string
	GetVariantValue() interface{}
}

// ENUM represents a DAML enum type
type ENUM interface {
	GetEnumConstructor() string
	GetEnumTypeID() string
}

func (p PARTY) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"_type": "party",
		"value": string(p),
	}
}
