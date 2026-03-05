package values

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

type Uint64 struct {
	Underlying uint64
}

func NewUint64(i uint64) *Uint64 {
	return &Uint64{Underlying: i}
}

func (i *Uint64) proto() *pb.Value {
	return pb.NewUInt64Value(i.Underlying)
}

func (i *Uint64) Unwrap() (any, error) {
	var u uint64
	return u, i.UnwrapTo(&u)
}

func (i *Uint64) copy() Value {
	if i == nil {
		return nil
	}
	return &Uint64{Underlying: i.Underlying}
}

func (i *Uint64) UnwrapTo(to any) error {
	if i == nil {
		return errors.New("cannot unwrap nil values.Uint64")
	}

	if to == nil {
		return fmt.Errorf("cannot unwrap to nil pointer: %+v", to)
	}

	switch tv := to.(type) {
	case *int64:
		if err := verifyUnsignedBounds(math.MaxInt64, i.Underlying, "uint32"); err != nil {
			return err
		}

		*tv = int64(i.Underlying)
		return nil
	case *int32:
		if err := verifyUnsignedBounds(math.MaxInt32, i.Underlying, "int32"); err != nil {
			return err
		}

		*tv = int32(i.Underlying)
		return nil
	case *int16:
		if err := verifyUnsignedBounds(math.MaxInt16, i.Underlying, "int16"); err != nil {
			return err
		}

		*tv = int16(i.Underlying)
		return nil
	case *int8:
		if err := verifyUnsignedBounds(math.MaxInt8, i.Underlying, "int8"); err != nil {
			return err
		}

		*tv = int8(i.Underlying)
		return nil
	case *int:
		if err := verifyUnsignedBounds(math.MaxInt, i.Underlying, "int"); err != nil {
			return err
		}

		*tv = int(i.Underlying)
		return nil
	case *uint64:
		*tv = i.Underlying
		return nil
	case *uint32:
		if err := verifyUnsignedBounds(math.MaxUint32, i.Underlying, "uint32"); err != nil {
			return err
		}

		*tv = uint32(i.Underlying)
		return nil
	case *uint16:
		if err := verifyUnsignedBounds(math.MaxUint16, i.Underlying, "uint16"); err != nil {
			return err
		}

		*tv = uint16(i.Underlying)
		return nil
	case *uint8:
		if err := verifyUnsignedBounds(math.MaxUint8, i.Underlying, "uint8"); err != nil {
			return err
		}

		*tv = uint8(i.Underlying)
		return nil
	case *uint:
		if verifyUnsignedBounds(math.MaxUint, i.Underlying, "uint") != nil {
			return fmt.Errorf("value %d is too large for uint", i.Underlying)
		}

		*tv = uint(i.Underlying)
		return nil
	case *any:
		*tv = i.Underlying
		return nil
	}

	rv := reflect.ValueOf(to)
	if rv.Kind() == reflect.Ptr {
		switch rv.Elem().Kind() {
		case reflect.Int64:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(int64(0)))).Interface())
		case reflect.Int32:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(int32(0)))).Interface())
		case reflect.Int16:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(int16(0)))).Interface())
		case reflect.Int8:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(int8(0)))).Interface())
		case reflect.Int:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(0))).Interface())
		case reflect.Uint64:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(uint64(0)))).Interface())
		case reflect.Uint32:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(uint32(0)))).Interface())
		case reflect.Uint16:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(uint16(0)))).Interface())
		case reflect.Uint8:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(uint8(0)))).Interface())
		case reflect.Uint:
			return i.UnwrapTo(rv.Convert(reflect.PointerTo(reflect.TypeOf(uint(0)))).Interface())
		default:
			// fall through to the error, default is required by lint
		}
	}

	return fmt.Errorf("cannot unwrap to type %T", to)
}

func verifyUnsignedBounds(max, value uint64, tpe string) error {
	if value > max {
		return fmt.Errorf("value %d is too large for %s", value, tpe)
	}
	return nil
}
