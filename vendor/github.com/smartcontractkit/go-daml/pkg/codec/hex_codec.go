package codec

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/go-daml/pkg/types"
)

// HexCodec encodes/decodes values to/from hex strings using the Canton MCMS format.
// This format is used for encoding operation parameters in Canton smart contracts.
//
// Encoding format:
//   - uint8: 1 byte, hex-encoded
//   - uint32/int: 4 bytes, big-endian, hex-encoded
//   - int64: 8 bytes, big-endian, hex-encoded
//   - bool: 1 byte (00 or 01)
//   - text/string: uint8(len) + utf8 bytes
//   - list/slice: uint8(count) + items (encoded sequentially)
//   - struct: concatenate fields in order
type HexCodec struct{}

// NewHexCodec creates a new HexCodec instance
func NewHexCodec() *HexCodec {
	return &HexCodec{}
}

// Marshal encodes value to hex string
func (c *HexCodec) Marshal(value interface{}) (string, error) {
	bytes, err := c.encode(value)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Unmarshal decodes hex string to value
func (c *HexCodec) Unmarshal(data string, target interface{}) error {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("failed to decode hex string: %w", err)
	}
	_, err = c.decode(bytes, 0, target)
	return err
}

// encode dispatches to type-specific encoders
func (c *HexCodec) encode(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	// Handle DAML types first
	switch v := value.(type) {
	case types.TEXT:
		return c.encodeText(string(v)), nil
	case types.INT64:
		return c.encodeInt64(int64(v)), nil
	case types.BOOL:
		return c.encodeBool(bool(v)), nil
	case types.PARTY:
		return c.encodeText(string(v)), nil
	case types.CONTRACT_ID:
		return c.encodeText(string(v)), nil
	case types.NUMERIC:
		return c.encodeText(string(v)), nil
	case types.RELTIME:
		// RELTIME stored as time.Duration, encode microseconds as int64
		microseconds := int64(time.Duration(v) / time.Microsecond)
		return c.encodeInt64(microseconds), nil
	case types.UNIT:
		// UNIT is empty - encode as 0 bytes
		return []byte{}, nil
	case types.DECIMAL:
		// DECIMAL is *big.Int, encode as length-prefixed string representation
		if v == nil {
			return c.encodeText(""), nil
		}
		return c.encodeText((*big.Int)(v).String()), nil
	case types.TIMESTAMP:
		// TIMESTAMP as microseconds since epoch (int64)
		microseconds := time.Time(v).UnixMicro()
		return c.encodeInt64(microseconds), nil
	case types.DATE:
		// DATE as days since epoch (int32)
		days := int32(time.Time(v).Unix() / 86400)
		return c.encodeInt32(days), nil
	case types.SET:
		// SET is []interface{}, encode like slice
		return c.encodeSlice(reflect.ValueOf([]interface{}(v)))
	case types.LIST:
		// LIST is []string, encode like slice
		return c.encodeSlice(reflect.ValueOf([]string(v)))
	case types.TEXTMAP:
		return c.encodeTextMap(v)
	case types.GENMAP:
		return c.encodeGenMap(v)
	case types.MAP:
		return c.encodeGenMap(map[string]interface{}(v))
	case types.TUPLE2:
		return c.encodeTuple2(v)
	}

	// Handle DAML enum/variant types (types with GetEnumConstructor method)
	// These are string-based enums that should be encoded as their constructor name
	type enumConstructorGetter interface {
		GetEnumConstructor() string
	}
	if e, ok := value.(enumConstructorGetter); ok {
		return c.encodeText(e.GetEnumConstructor()), nil
	}

	// Handle VARIANT types (implements types.VARIANT interface)
	if variant, ok := value.(types.VARIANT); ok {
		return c.encodeVariant(variant)
	}

	// Handle Go primitive types
	switch v := value.(type) {
	case string:
		return c.encodeText(v), nil
	case bool:
		return c.encodeBool(v), nil
	case int:
		return c.encodeInt32(int32(v)), nil
	case int8:
		return c.encodeUint8(uint8(v)), nil
	case int16:
		return c.encodeInt16(int16(v)), nil
	case int32:
		return c.encodeInt32(v), nil
	case int64:
		return c.encodeInt64(v), nil
	case uint:
		return c.encodeUint32(uint32(v)), nil
	case uint8:
		return c.encodeUint8(v), nil
	case uint16:
		return c.encodeUint16(v), nil
	case uint32:
		return c.encodeUint32(v), nil
	case uint64:
		return c.encodeUint64(v), nil
	case []byte:
		// Encode as text (length-prefixed)
		return c.encodeBytes(v), nil
	}

	// Handle via reflection
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		return c.encode(rv.Elem().Interface())
	}

	if rv.Kind() == reflect.Struct {
		return c.encodeStruct(rv)
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		return c.encodeSlice(rv)
	}

	return nil, fmt.Errorf("unsupported type for hex encoding: %T", value)
}

// Primitive encoders matching Canton Codec.daml

func (c *HexCodec) encodeUint8(v uint8) []byte {
	return []byte{v}
}

func (c *HexCodec) encodeUint16(v uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}

func (c *HexCodec) encodeUint32(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func (c *HexCodec) encodeUint64(v uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf
}

func (c *HexCodec) encodeInt16(v int16) []byte {
	return c.encodeUint16(uint16(v))
}

func (c *HexCodec) encodeInt32(v int32) []byte {
	return c.encodeUint32(uint32(v))
}

func (c *HexCodec) encodeInt64(v int64) []byte {
	return c.encodeUint64(uint64(v))
}

func (c *HexCodec) encodeText(s string) []byte {
	b := []byte(s)
	result := make([]byte, 1+len(b))
	result[0] = byte(len(b))
	copy(result[1:], b)
	return result
}

func (c *HexCodec) encodeBytes(b []byte) []byte {
	result := make([]byte, 1+len(b))
	result[0] = byte(len(b))
	copy(result[1:], b)
	return result
}

func (c *HexCodec) encodeBool(v bool) []byte {
	if v {
		return []byte{1}
	}
	return []byte{0}
}

func (c *HexCodec) encodeStruct(rv reflect.Value) ([]byte, error) {
	var result []byte
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanInterface() {
			continue // Skip unexported fields
		}

		fieldType := rt.Field(i)
		var encoded []byte
		var err error

		// Check for hex tag - supports various encoding overrides
		hexTag := fieldType.Tag.Get("hex")
		switch hexTag {
		case "bytes":
			// hex:"bytes" - string contains hex-encoded bytes, decode and encode as length-prefixed bytes
			if field.Kind() == reflect.String {
				hexStr := field.String()
				decoded, decErr := hex.DecodeString(hexStr)
				if decErr != nil {
					return nil, fmt.Errorf("failed to decode hex field %s: %w", fieldType.Name, decErr)
				}
				encoded = c.encodeBytes(decoded)
			} else {
				return nil, fmt.Errorf("hex:\"bytes\" tag only valid on string fields, got %v", field.Kind())
			}
		case "uint32":
			// hex:"uint32" - encode INT64 or int64 as 4-byte uint32 (for Canton compatibility)
			var val int64
			switch v := field.Interface().(type) {
			case types.INT64:
				val = int64(v)
			case int64:
				val = v
			case int:
				val = int64(v)
			default:
				return nil, fmt.Errorf("hex:\"uint32\" tag only valid on INT64/int64/int fields, got %T", field.Interface())
			}
			encoded = c.encodeUint32(uint32(val))
		case "uint16":
			// hex:"uint16" - encode as 2-byte uint16
			var val int64
			switch v := field.Interface().(type) {
			case types.INT64:
				val = int64(v)
			case int64:
				val = v
			case int:
				val = int64(v)
			default:
				return nil, fmt.Errorf("hex:\"uint16\" tag only valid on INT64/int64/int fields, got %T", field.Interface())
			}
			encoded = c.encodeUint16(uint16(val))
		case "uint8":
			// hex:"uint8" - encode as 1-byte uint8
			var val int64
			switch v := field.Interface().(type) {
			case types.INT64:
				val = int64(v)
			case int64:
				val = v
			case int:
				val = int64(v)
			default:
				return nil, fmt.Errorf("hex:\"uint8\" tag only valid on INT64/int64/int fields, got %T", field.Interface())
			}
			encoded = c.encodeUint8(uint8(val))
		case "[]uint32":
			// hex:"[]uint32" - encode slice of INT64/int64/int as length + uint32 elements
			if field.Kind() != reflect.Slice {
				return nil, fmt.Errorf("hex:\"[]uint32\" tag only valid on slice fields, got %v", field.Kind())
			}
			length := field.Len()
			if length > 255 {
				return nil, fmt.Errorf("slice length %d exceeds max 255 for hex:\"[]uint32\"", length)
			}
			encoded = []byte{byte(length)}
			for j := 0; j < length; j++ {
				elem := field.Index(j)
				var val int64
				switch v := elem.Interface().(type) {
				case types.INT64:
					val = int64(v)
				case int64:
					val = v
				case int:
					val = int64(v)
				default:
					return nil, fmt.Errorf("hex:\"[]uint32\" element %d: expected INT64/int64/int, got %T", j, elem.Interface())
				}
				encoded = append(encoded, c.encodeUint32(uint32(val))...)
			}
		case "[]uint64":
			// hex:"[]uint64" - encode slice of INT64/int64/int as length + uint64 elements (8 bytes each)
			if field.Kind() != reflect.Slice {
				return nil, fmt.Errorf("hex:\"[]uint64\" tag only valid on slice fields, got %v", field.Kind())
			}
			length := field.Len()
			if length > 255 {
				return nil, fmt.Errorf("slice length %d exceeds max 255 for hex:\"[]uint64\"", length)
			}
			encoded = []byte{byte(length)}
			for j := 0; j < length; j++ {
				elem := field.Index(j)
				var val int64
				switch v := elem.Interface().(type) {
				case types.INT64:
					val = int64(v)
				case int64:
					val = v
				case int:
					val = int64(v)
				default:
					return nil, fmt.Errorf("hex:\"[]uint64\" element %d: expected INT64/int64/int, got %T", j, elem.Interface())
				}
				encoded = append(encoded, c.encodeInt64(val)...)
			}
		default:
			// No tag or unknown tag - use default encoding
			encoded, err = c.encode(field.Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to encode field %s: %w", fieldType.Name, err)
			}
		}
		result = append(result, encoded...)
	}
	return result, nil
}

func (c *HexCodec) encodeSlice(rv reflect.Value) ([]byte, error) {
	length := rv.Len()
	if length > 255 {
		return nil, fmt.Errorf("slice length %d exceeds maximum of 255 for uint8 length prefix", length)
	}
	result := []byte{byte(length)}
	for i := 0; i < length; i++ {
		encoded, err := c.encode(rv.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to encode slice element %d: %w", i, err)
		}
		result = append(result, encoded...)
	}
	return result, nil
}

func (c *HexCodec) encodeTextMap(m types.TEXTMAP) ([]byte, error) {
	length := len(m)
	if length > 255 {
		return nil, fmt.Errorf("map length %d exceeds max 255", length)
	}
	result := []byte{byte(length)}
	for k, v := range m {
		result = append(result, c.encodeText(k)...)
		encoded, err := c.encode(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode map value for key %s: %w", k, err)
		}
		result = append(result, encoded...)
	}
	return result, nil
}

func (c *HexCodec) encodeGenMap(m map[string]interface{}) ([]byte, error) {
	length := len(m)
	if length > 255 {
		return nil, fmt.Errorf("map length %d exceeds max 255", length)
	}
	result := []byte{byte(length)}
	for k, v := range m {
		result = append(result, c.encodeText(k)...)
		encoded, err := c.encode(v)
		if err != nil {
			return nil, fmt.Errorf("failed to encode map value for key %s: %w", k, err)
		}
		result = append(result, encoded...)
	}
	return result, nil
}

func (c *HexCodec) encodeVariant(v types.VARIANT) ([]byte, error) {
	tag := v.GetVariantTag()
	value := v.GetVariantValue()
	result := c.encodeText(tag)
	if value != nil {
		encoded, err := c.encode(value)
		if err != nil {
			return nil, fmt.Errorf("failed to encode variant value: %w", err)
		}
		result = append(result, encoded...)
	}
	return result, nil
}

func (c *HexCodec) encodeTuple2(t types.TUPLE2) ([]byte, error) {
	var result []byte
	encoded, err := c.encode(t.First)
	if err != nil {
		return nil, fmt.Errorf("failed to encode tuple2 first field: %w", err)
	}
	result = append(result, encoded...)
	encoded, err = c.encode(t.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to encode tuple2 second field: %w", err)
	}
	result = append(result, encoded...)
	return result, nil
}

// Decode methods

func (c *HexCodec) decode(data []byte, offset int, target interface{}) (int, error) {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return offset, fmt.Errorf("target must be a pointer, got %v", rv.Kind())
	}
	if rv.IsNil() {
		return offset, fmt.Errorf("target pointer is nil")
	}

	elem := rv.Elem()
	return c.decodeValue(data, offset, elem)
}

func (c *HexCodec) decodeValue(data []byte, offset int, target reflect.Value) (int, error) {
	targetType := target.Type()

	// Handle DAML types
	switch targetType {
	case reflect.TypeOf(types.TEXT("")):
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		target.Set(reflect.ValueOf(types.TEXT(s)))
		return newOffset, nil

	case reflect.TypeOf(types.INT64(0)):
		if offset+8 > len(data) {
			return offset, fmt.Errorf("not enough data for INT64 at offset %d", offset)
		}
		v := int64(binary.BigEndian.Uint64(data[offset : offset+8]))
		target.Set(reflect.ValueOf(types.INT64(v)))
		return offset + 8, nil

	case reflect.TypeOf(types.BOOL(false)):
		if offset >= len(data) {
			return offset, fmt.Errorf("not enough data for BOOL at offset %d", offset)
		}
		target.Set(reflect.ValueOf(types.BOOL(data[offset] != 0)))
		return offset + 1, nil

	case reflect.TypeOf(types.PARTY("")):
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		target.Set(reflect.ValueOf(types.PARTY(s)))
		return newOffset, nil

	case reflect.TypeOf(types.CONTRACT_ID("")):
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		target.Set(reflect.ValueOf(types.CONTRACT_ID(s)))
		return newOffset, nil

	case reflect.TypeOf(types.NUMERIC("")):
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		target.Set(reflect.ValueOf(types.NUMERIC(s)))
		return newOffset, nil

	case reflect.TypeOf(types.RELTIME(0)):
		if offset+8 > len(data) {
			return offset, fmt.Errorf("not enough data for RELTIME at offset %d", offset)
		}
		microseconds := int64(binary.BigEndian.Uint64(data[offset : offset+8]))
		target.Set(reflect.ValueOf(types.RELTIME(time.Duration(microseconds) * time.Microsecond)))
		return offset + 8, nil

	case reflect.TypeOf(types.UNIT{}):
		target.Set(reflect.ValueOf(types.UNIT{}))
		return offset, nil

	case reflect.TypeOf(types.TIMESTAMP{}):
		if offset+8 > len(data) {
			return offset, fmt.Errorf("not enough data for TIMESTAMP at offset %d", offset)
		}
		microseconds := int64(binary.BigEndian.Uint64(data[offset : offset+8]))
		t := time.UnixMicro(microseconds)
		target.Set(reflect.ValueOf(types.TIMESTAMP(t)))
		return offset + 8, nil

	case reflect.TypeOf(types.DATE{}):
		if offset+4 > len(data) {
			return offset, fmt.Errorf("not enough data for DATE at offset %d", offset)
		}
		days := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
		t := time.Unix(int64(days)*86400, 0).UTC()
		target.Set(reflect.ValueOf(types.DATE(t)))
		return offset + 4, nil
	}

	// Handle DECIMAL (needs special handling as it's a type alias for *big.Int)
	if targetType == reflect.TypeOf(types.DECIMAL(nil)) {
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		if s == "" {
			target.Set(reflect.ValueOf(types.DECIMAL(nil)))
			return newOffset, nil
		}
		bigInt := new(big.Int)
		bigInt, ok := bigInt.SetString(s, 10)
		if !ok {
			return offset, fmt.Errorf("failed to parse DECIMAL value: %s", s)
		}
		target.Set(reflect.ValueOf(types.DECIMAL(bigInt)))
		return newOffset, nil
	}

	// Handle Go types based on kind
	switch target.Kind() {
	case reflect.String:
		s, newOffset, err := c.decodeText(data, offset)
		if err != nil {
			return offset, err
		}
		target.SetString(s)
		return newOffset, nil

	case reflect.Bool:
		if offset >= len(data) {
			return offset, fmt.Errorf("not enough data for bool at offset %d", offset)
		}
		target.SetBool(data[offset] != 0)
		return offset + 1, nil

	case reflect.Int, reflect.Int32:
		if offset+4 > len(data) {
			return offset, fmt.Errorf("not enough data for int32 at offset %d", offset)
		}
		v := int32(binary.BigEndian.Uint32(data[offset : offset+4]))
		target.SetInt(int64(v))
		return offset + 4, nil

	case reflect.Int8:
		if offset >= len(data) {
			return offset, fmt.Errorf("not enough data for int8 at offset %d", offset)
		}
		target.SetInt(int64(int8(data[offset])))
		return offset + 1, nil

	case reflect.Int16:
		if offset+2 > len(data) {
			return offset, fmt.Errorf("not enough data for int16 at offset %d", offset)
		}
		v := int16(binary.BigEndian.Uint16(data[offset : offset+2]))
		target.SetInt(int64(v))
		return offset + 2, nil

	case reflect.Int64:
		if offset+8 > len(data) {
			return offset, fmt.Errorf("not enough data for int64 at offset %d", offset)
		}
		v := int64(binary.BigEndian.Uint64(data[offset : offset+8]))
		target.SetInt(v)
		return offset + 8, nil

	case reflect.Uint, reflect.Uint32:
		if offset+4 > len(data) {
			return offset, fmt.Errorf("not enough data for uint32 at offset %d", offset)
		}
		v := binary.BigEndian.Uint32(data[offset : offset+4])
		target.SetUint(uint64(v))
		return offset + 4, nil

	case reflect.Uint8:
		if offset >= len(data) {
			return offset, fmt.Errorf("not enough data for uint8 at offset %d", offset)
		}
		target.SetUint(uint64(data[offset]))
		return offset + 1, nil

	case reflect.Uint16:
		if offset+2 > len(data) {
			return offset, fmt.Errorf("not enough data for uint16 at offset %d", offset)
		}
		v := binary.BigEndian.Uint16(data[offset : offset+2])
		target.SetUint(uint64(v))
		return offset + 2, nil

	case reflect.Uint64:
		if offset+8 > len(data) {
			return offset, fmt.Errorf("not enough data for uint64 at offset %d", offset)
		}
		v := binary.BigEndian.Uint64(data[offset : offset+8])
		target.SetUint(v)
		return offset + 8, nil

	case reflect.Ptr:
		// Allocate new value and decode into it
		newVal := reflect.New(target.Type().Elem())
		newOffset, err := c.decodeValue(data, offset, newVal.Elem())
		if err != nil {
			return offset, err
		}
		target.Set(newVal)
		return newOffset, nil

	case reflect.Slice:
		return c.decodeSlice(data, offset, target)

	case reflect.Struct:
		return c.decodeStruct(data, offset, target)

	default:
		return offset, fmt.Errorf("unsupported target type for hex decoding: %v", target.Type())
	}
}

func (c *HexCodec) decodeText(data []byte, offset int) (string, int, error) {
	if offset >= len(data) {
		return "", offset, fmt.Errorf("not enough data for text length at offset %d", offset)
	}
	length := int(data[offset])
	offset++

	if offset+length > len(data) {
		return "", offset, fmt.Errorf("not enough data for text of length %d at offset %d", length, offset)
	}
	s := string(data[offset : offset+length])
	return s, offset + length, nil
}

func (c *HexCodec) decodeSlice(data []byte, offset int, target reflect.Value) (int, error) {
	if offset >= len(data) {
		return offset, fmt.Errorf("not enough data for slice length at offset %d", offset)
	}
	length := int(data[offset])
	offset++

	slice := reflect.MakeSlice(target.Type(), length, length)
	for i := 0; i < length; i++ {
		var err error
		offset, err = c.decodeValue(data, offset, slice.Index(i))
		if err != nil {
			return offset, fmt.Errorf("failed to decode slice element %d: %w", i, err)
		}
	}
	target.Set(slice)
	return offset, nil
}

func (c *HexCodec) decodeStruct(data []byte, offset int, target reflect.Value) (int, error) {
	targetType := target.Type()
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		if !field.CanSet() {
			continue // Skip unexported fields
		}

		fieldType := targetType.Field(i)
		var err error

		// Check for hex tag - supports various encoding overrides
		hexTag := fieldType.Tag.Get("hex")
		switch hexTag {
		case "bytes":
			// hex:"bytes" - decode as length-prefixed bytes, store as hex string
			if field.Kind() == reflect.String {
				if offset >= len(data) {
					return offset, fmt.Errorf("not enough data for hex bytes length at offset %d", offset)
				}
				length := int(data[offset])
				offset++
				if offset+length > len(data) {
					return offset, fmt.Errorf("not enough data for hex bytes of length %d at offset %d", length, offset)
				}
				rawBytes := data[offset : offset+length]
				field.SetString(hex.EncodeToString(rawBytes))
				offset += length
			} else {
				return offset, fmt.Errorf("hex:\"bytes\" tag only valid on string fields, got %v", field.Kind())
			}
		case "uint32":
			// hex:"uint32" - decode 4-byte uint32 into INT64/int64/int
			if offset+4 > len(data) {
				return offset, fmt.Errorf("not enough data for uint32 at offset %d", offset)
			}
			val := binary.BigEndian.Uint32(data[offset : offset+4])
			offset += 4
			switch field.Type() {
			case reflect.TypeOf(types.INT64(0)):
				field.Set(reflect.ValueOf(types.INT64(val)))
			default:
				field.SetInt(int64(val))
			}
		case "uint16":
			// hex:"uint16" - decode 2-byte uint16 into INT64/int64/int
			if offset+2 > len(data) {
				return offset, fmt.Errorf("not enough data for uint16 at offset %d", offset)
			}
			val := binary.BigEndian.Uint16(data[offset : offset+2])
			offset += 2
			switch field.Type() {
			case reflect.TypeOf(types.INT64(0)):
				field.Set(reflect.ValueOf(types.INT64(val)))
			default:
				field.SetInt(int64(val))
			}
		case "uint8":
			// hex:"uint8" - decode 1-byte uint8 into INT64/int64/int
			if offset >= len(data) {
				return offset, fmt.Errorf("not enough data for uint8 at offset %d", offset)
			}
			val := data[offset]
			offset++
			switch field.Type() {
			case reflect.TypeOf(types.INT64(0)):
				field.Set(reflect.ValueOf(types.INT64(val)))
			default:
				field.SetInt(int64(val))
			}
		case "[]uint32":
			// hex:"[]uint32" - decode length + uint32 elements into slice of INT64/int64/int
			if offset >= len(data) {
				return offset, fmt.Errorf("not enough data for []uint32 length at offset %d", offset)
			}
			length := int(data[offset])
			offset++
			if offset+length*4 > len(data) {
				return offset, fmt.Errorf("not enough data for %d uint32 elements at offset %d", length, offset)
			}
			slice := reflect.MakeSlice(field.Type(), length, length)
			for j := 0; j < length; j++ {
				val := binary.BigEndian.Uint32(data[offset : offset+4])
				offset += 4
				elem := slice.Index(j)
				switch elem.Type() {
				case reflect.TypeOf(types.INT64(0)):
					elem.Set(reflect.ValueOf(types.INT64(val)))
				default:
					elem.SetInt(int64(val))
				}
			}
			field.Set(slice)
		case "[]uint64":
			// hex:"[]uint64" - decode length + uint64 elements into slice of INT64/int64/int
			if offset >= len(data) {
				return offset, fmt.Errorf("not enough data for []uint64 length at offset %d", offset)
			}
			length := int(data[offset])
			offset++
			if offset+length*8 > len(data) {
				return offset, fmt.Errorf("not enough data for %d uint64 elements at offset %d", length, offset)
			}
			slice := reflect.MakeSlice(field.Type(), length, length)
			for j := 0; j < length; j++ {
				val := binary.BigEndian.Uint64(data[offset : offset+8])
				offset += 8
				elem := slice.Index(j)
				switch elem.Type() {
				case reflect.TypeOf(types.INT64(0)):
					elem.Set(reflect.ValueOf(types.INT64(val)))
				default:
					elem.SetInt(int64(val))
				}
			}
			field.Set(slice)
		default:
			// No tag or unknown tag - use default decoding
			offset, err = c.decodeValue(data, offset, field)
			if err != nil {
				return offset, fmt.Errorf("failed to decode field %s: %w", fieldType.Name, err)
			}
		}
	}
	return offset, nil
}
