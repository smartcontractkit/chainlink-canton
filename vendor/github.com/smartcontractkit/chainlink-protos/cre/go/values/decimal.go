package values

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

type Decimal struct {
	Underlying decimal.Decimal
}

func NewDecimal(d decimal.Decimal) *Decimal {
	return &Decimal{Underlying: d}
}

func (d *Decimal) proto() *pb.Value {
	return pb.NewDecimalValue(d.Underlying)
}

func (d *Decimal) Unwrap() (any, error) {
	var dec decimal.Decimal
	return dec, d.UnwrapTo(&dec)
}

var float64Type = reflect.TypeOf((*float64)(nil))
var maxFloat64 = decimal.NewFromFloat(math.MaxFloat64)
var minFloat64 = decimal.NewFromFloat(math.SmallestNonzeroFloat64)
var maxFloat32 = decimal.NewFromFloat(math.MaxFloat32)
var minFloat32 = decimal.NewFromFloat(math.SmallestNonzeroFloat32)

func (d *Decimal) UnwrapTo(to any) error {
	if d == nil {
		return errors.New("could not unwrap nil values.Decimal")
	}

	underlying := d.Underlying
	switch to.(type) {
	case *float64:
		f, err := safeRange(minFloat64, maxFloat64, underlying, "float64")
		if err != nil {
			return err
		}
		*to.(*float64) = f
		return nil
	case *float32:
		f, err := safeRange(minFloat32, maxFloat32, underlying, "float32")
		if err != nil {
			return err
		}
		*to.(*float32) = float32(f)
		return nil
	}

	vto := reflect.ValueOf(to)
	if vto.CanConvert(float64Type) {
		return d.UnwrapTo(vto.Convert(float64Type).Interface())
	} else if vto.CanConvert(float32Type) {
		return d.UnwrapTo(vto.Convert(float32Type).Interface())
	}

	return unwrapTo(d.Underlying, to)
}

func (d *Decimal) copy() Value {
	if d == nil {
		return nil
	}
	return &Decimal{Underlying: d.Underlying.Copy()}
}

func safeRange(min, max, value decimal.Decimal, name string) (float64, error) {
	if value.GreaterThan(max) {
		return 0, fmt.Errorf("decimal value exceeds maximum %s", name)
	} else if value.LessThan(min) {
		return 0, fmt.Errorf("decimal value is below minimum %s", name)
	}

	return value.InexactFloat64(), nil
}
