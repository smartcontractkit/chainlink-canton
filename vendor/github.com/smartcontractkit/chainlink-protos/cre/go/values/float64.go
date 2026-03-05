package values

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

type Float64 struct {
	Underlying float64
}

func NewFloat64(f float64) *Float64 {
	return &Float64{Underlying: f}
}

func (f *Float64) proto() *pb.Value {
	return pb.NewFloat64(f.Underlying)
}

func (f *Float64) Unwrap() (any, error) {
	var to float64
	return to, f.UnwrapTo(&to)
}

var decimalType = reflect.TypeOf(&decimal.Decimal{})
var float32Type = reflect.TypeOf((*float32)(nil))

func (f *Float64) UnwrapTo(to any) error {
	if f == nil {
		return errors.New("cannot unwrap nil values.Float64")
	}

	switch t := to.(type) {
	case *decimal.Decimal:
		if t == nil {
			return errors.New("cannot unwrap to nil pointer")
		}
		*t = decimal.NewFromFloat(f.Underlying)
		return nil
	case *float32:
		val, err := safeFloat64ToFloat32(f.Underlying)
		if err != nil {
			return err
		}
		*t = val
	default:
		toVal := reflect.ValueOf(to)
		toType := toVal.Type()
		if toType.ConvertibleTo(decimalType) {
			return f.UnwrapTo(toVal.Convert(decimalType).Interface())
		} else if toType.ConvertibleTo(float32Type) {
			return f.UnwrapTo(toVal.Convert(float32Type).Interface())
		}
	}

	return unwrapTo(f.Underlying, to)
}

func (f *Float64) copy() Value {
	if f == nil {
		return f
	}
	return &Float64{Underlying: f.Underlying}
}

func safeFloat64ToFloat32(f float64) (float32, error) {
	// Check for too large
	if f > math.MaxFloat32 {
		return 0, fmt.Errorf("%v is too large for a float32", f)
	}
	// Check for too small (excluding zero)
	if f != 0 && math.Abs(f) < math.SmallestNonzeroFloat32 {
		return 0, fmt.Errorf("%v is too small for a float32", f)
	}
	return float32(f), nil
}
