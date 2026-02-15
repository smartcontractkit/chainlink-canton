package codec

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
	"time"

	. "github.com/smartcontractkit/chainlink-canton/bindings/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonCodec_Marshal_BasicTypes(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "PARTY as string",
			input:    PARTY("alice"),
			expected: `"alice"`,
		},
		{
			name:     "TEXT as string",
			input:    TEXT("hello world"),
			expected: `"hello world"`,
		},
		{
			name:     "INT64 as string (default)",
			input:    INT64(42),
			expected: `"42"`,
		},
		{
			name:     "BOOL as boolean",
			input:    BOOL(true),
			expected: `true`,
		},
		{
			name:     "NUMERIC as string (default)",
			input:    NUMERIC("123456789"),
			expected: `"123456789"`,
		},
		{
			name:     "DECIMAL as string (default)",
			input:    DECIMAL(big.NewInt(987654321)),
			expected: `"987654321"`,
		},
		{
			name:     "TIMESTAMP as ISO string",
			input:    TIMESTAMP(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
			expected: `"2023-01-01T12:00:00.000000Z"`,
		},
		{
			name:     "DATE as ISO date string",
			input:    DATE(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
			expected: `"2023-01-01"`,
		},
		{
			name:     "UNIT as empty object",
			input:    UNIT{},
			expected: `{}`,
		},
		{
			name:     "CONTRACT_ID as string",
			input:    CONTRACT_ID("contract-123"),
			expected: `"contract-123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestJsonCodec_Marshal_NumericAsNumber(t *testing.T) {
	codec := NewJsonCodecWithOptions(false, false, false)

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "INT64 as number",
			input:    INT64(42),
			expected: `42`,
		},
		{
			name:     "NUMERIC as number",
			input:    NUMERIC("123"),
			expected: `123`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestJsonCodec_Marshal_Collections(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "GENMAP",
			input: GENMAP{
				"key1": "value1",
				"key2": "value2",
			},
			expected: `{"key1":"value1","key2":"value2"}`,
		},
		{
			name:     "LIST of strings",
			input:    LIST{"item1", "item2", "item3"},
			expected: `["item1","item2","item3"]`,
		},
		{
			name:     "slice of INT64",
			input:    []INT64{1, 2, 3},
			expected: `["1","2","3"]`, // INT64 encoded as strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

type TestRecord struct {
	Name     TEXT    `json:"name"`
	Age      INT64   `json:"age"`
	Active   BOOL    `json:"active"`
	Balance  NUMERIC `json:"balance"`
	Optional *TEXT   `json:"optional"`
}

func TestJsonCodec_Marshal_Records(t *testing.T) {
	codec := NewJsonCodec()

	optionalValue := TEXT("present")
	record := TestRecord{
		Name:     TEXT("Alice"),
		Age:      INT64(30),
		Active:   BOOL(true),
		Balance:  NUMERIC("1000"),
		Optional: &optionalValue,
	}

	result, err := codec.Marshal(record)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "Alice", parsed["name"])
	assert.Equal(t, "30", parsed["age"]) // INT64 as string
	assert.Equal(t, true, parsed["active"])
	assert.Equal(t, "1000", parsed["balance"]) // NUMERIC as string
	assert.Equal(t, "present", parsed["optional"])
}

func TestJsonCodec_Marshal_RecordWithNilOptional(t *testing.T) {
	codec := NewJsonCodec()

	record := TestRecord{
		Name:     TEXT("Bob"),
		Age:      INT64(25),
		Active:   BOOL(false),
		Balance:  NUMERIC("500"),
		Optional: nil, // nil optional
	}

	result, err := codec.Marshal(record)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "Bob", parsed["name"])
	assert.Equal(t, "25", parsed["age"])
	assert.Equal(t, false, parsed["active"])
	assert.Equal(t, "500", parsed["balance"])
	assert.Nil(t, parsed["optional"]) // nil optional included
}

func TestJsonCodec_Marshal_ExcludeNullValues(t *testing.T) {
	codec := NewJsonCodecWithOptions(true, true, true) // exclude null values

	record := TestRecord{
		Name:     TEXT("Charlie"),
		Age:      INT64(35),
		Active:   BOOL(true),
		Balance:  NUMERIC("750"),
		Optional: nil, // nil optional
	}

	result, err := codec.Marshal(record)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "Charlie", parsed["name"])
	assert.Equal(t, "35", parsed["age"])
	assert.Equal(t, true, parsed["active"])
	assert.Equal(t, "750", parsed["balance"])

	_, exists := parsed["optional"]
	assert.False(t, exists)
}

type TestColor string

const (
	TestColorRed   TestColor = "Red"
	TestColorGreen TestColor = "Green"
	TestColorBlue  TestColor = "Blue"
)

func (e TestColor) GetEnumConstructor() string {
	return string(e)
}

func (e TestColor) GetEnumTypeID() string {
	return "TestColor"
}

var _ ENUM = TestColor("")

func TestJsonCodec_Marshal_Enum(t *testing.T) {
	codec := NewJsonCodec()

	result, err := codec.Marshal(TestColorRed)
	require.NoError(t, err)

	assert.JSONEq(t, `"Red"`, string(result))
}

type ComplexRecord struct {
	Owner    PARTY                  `json:"owner"`
	Metadata map[string]interface{} `json:"metadata"`
	Tags     []string               `json:"tags"`
	Config   TestRecord             `json:"config"`
}

func TestJsonCodec_Marshal_ComplexNested(t *testing.T) {
	codec := NewJsonCodec()

	optionalText := TEXT("nested")
	complex := ComplexRecord{
		Owner: PARTY("alice"),
		Metadata: map[string]interface{}{
			"version": "1.0",
			"type":    "test",
		},
		Tags: []string{"important", "test"},
		Config: TestRecord{
			Name:     TEXT("Nested"),
			Age:      INT64(40),
			Active:   BOOL(true),
			Balance:  NUMERIC("2000"),
			Optional: &optionalText,
		},
	}

	result, err := codec.Marshal(complex)
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "alice", parsed["owner"])
	assert.Equal(t, []interface{}{"important", "test"}, parsed["tags"])

	metadata := parsed["metadata"].(map[string]interface{})
	assert.Equal(t, "1.0", metadata["version"])
	assert.Equal(t, "test", metadata["type"])

	config := parsed["config"].(map[string]interface{})
	assert.Equal(t, "Nested", config["name"])
	assert.Equal(t, "40", config["age"])
	assert.Equal(t, true, config["active"])
	assert.Equal(t, "2000", config["balance"])
	assert.Equal(t, "nested", config["optional"])
}

// ========== UNMARSHAL TESTS ==========

func TestJsonCodec_Unmarshal_BasicTypes(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		json     string
		target   interface{}
		expected interface{}
	}{
		{
			name:     "PARTY from string",
			json:     `"alice"`,
			target:   new(PARTY),
			expected: PARTY("alice"),
		},
		{
			name:     "TEXT from string",
			json:     `"hello world"`,
			target:   new(TEXT),
			expected: TEXT("hello world"),
		},
		{
			name:     "INT64 from string",
			json:     `"42"`,
			target:   new(INT64),
			expected: INT64(42),
		},
		{
			name:     "INT64 from number",
			json:     `42`,
			target:   new(INT64),
			expected: INT64(42),
		},
		{
			name:     "BOOL from boolean",
			json:     `true`,
			target:   new(BOOL),
			expected: BOOL(true),
		},
		{
			name:     "NUMERIC from string",
			json:     `"123456789"`,
			target:   new(NUMERIC),
			expected: NUMERIC("123456789"),
		},
		{
			name:     "NUMERIC from number",
			json:     `123`,
			target:   new(NUMERIC),
			expected: NUMERIC("123"),
		},
		{
			name:     "TIMESTAMP from ISO string",
			json:     `"2023-01-01T12:00:00.000000Z"`,
			target:   new(TIMESTAMP),
			expected: TIMESTAMP(time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)),
		},
		{
			name:     "DATE from ISO date string",
			json:     `"2023-01-01"`,
			target:   new(DATE),
			expected: DATE(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "UNIT from empty object",
			json:     `{}`,
			target:   new(UNIT),
			expected: UNIT{},
		},
		{
			name:     "CONTRACT_ID from string",
			json:     `"contract-123"`,
			target:   new(CONTRACT_ID),
			expected: CONTRACT_ID("contract-123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := codec.Unmarshal([]byte(tt.json), tt.target)
			require.NoError(t, err)

			rv := reflect.ValueOf(tt.target).Elem().Interface()
			assert.Equal(t, tt.expected, rv)
		})
	}
}

func TestJsonCodec_Unmarshal_Collections(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		json     string
		target   interface{}
		expected interface{}
	}{
		{
			name:   "GENMAP",
			json:   `{"key1":"value1","key2":"value2"}`,
			target: new(GENMAP),
			expected: GENMAP{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "LIST of strings",
			json:     `["item1","item2","item3"]`,
			target:   new(LIST),
			expected: LIST{"item1", "item2", "item3"},
		},
		{
			name:     "slice of INT64 from strings",
			json:     `["1","2","3"]`,
			target:   new([]INT64),
			expected: []INT64{1, 2, 3},
		},
		{
			name:     "slice of INT64 from numbers",
			json:     `[1,2,3]`,
			target:   new([]INT64),
			expected: []INT64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := codec.Unmarshal([]byte(tt.json), tt.target)
			require.NoError(t, err)

			rv := reflect.ValueOf(tt.target).Elem().Interface()
			assert.Equal(t, tt.expected, rv)
		})
	}
}

func TestJsonCodec_Unmarshal_Records(t *testing.T) {
	codec := NewJsonCodec()

	jsonData := `{
		"name": "Alice",
		"age": "30",
		"active": true,
		"balance": "1000",
		"optional": "present"
	}`

	var result TestRecord
	err := codec.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Equal(t, TEXT("Alice"), result.Name)
	assert.Equal(t, INT64(30), result.Age)
	assert.Equal(t, BOOL(true), result.Active)
	assert.Equal(t, NUMERIC("1000"), result.Balance)
	require.NotNil(t, result.Optional)
	assert.Equal(t, TEXT("present"), *result.Optional)
}

func TestJsonCodec_Unmarshal_RecordWithNilOptional(t *testing.T) {
	codec := NewJsonCodec()

	jsonData := `{
		"name": "Bob",
		"age": "25",
		"active": false,
		"balance": "500"
	}`

	var result TestRecord
	err := codec.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Equal(t, TEXT("Bob"), result.Name)
	assert.Equal(t, INT64(25), result.Age)
	assert.Equal(t, BOOL(false), result.Active)
	assert.Equal(t, NUMERIC("500"), result.Balance)
	assert.Nil(t, result.Optional)
}

func TestJsonCodec_RoundTrip_Marshal_Unmarshal(t *testing.T) {
	codec := NewJsonCodec()

	optionalValue := TEXT("present")
	original := TestRecord{
		Name:     TEXT("Alice"),
		Age:      INT64(30),
		Active:   BOOL(true),
		Balance:  NUMERIC("1000"),
		Optional: &optionalValue,
	}

	jsonBytes, err := codec.Marshal(original)
	require.NoError(t, err)

	var result TestRecord
	err = codec.Unmarshal(jsonBytes, &result)
	require.NoError(t, err)
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Age, result.Age)
	assert.Equal(t, original.Active, result.Active)

	assert.Equal(t, string(original.Balance), string(result.Balance))
	require.NotNil(t, result.Optional)
	assert.Equal(t, *original.Optional, *result.Optional)
}

func TestJsonCodec_RoundTrip_WithNumericAsNumber(t *testing.T) {
	codec := NewJsonCodecWithOptions(false, false, false)

	original := TestRecord{
		Name:     TEXT("Bob"),
		Age:      INT64(25),
		Active:   BOOL(false),
		Balance:  NUMERIC("500"),
		Optional: nil,
	}

	jsonBytes, err := codec.Marshal(original)
	require.NoError(t, err)

	var result TestRecord
	err = codec.Unmarshal(jsonBytes, &result)
	require.NoError(t, err)

	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Age, result.Age)
	assert.Equal(t, original.Active, result.Active)
	assert.Equal(t, string(original.Balance), string(result.Balance))
	assert.Nil(t, result.Optional)
}

func TestJsonCodec_Marshal_RELTIME(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		input    RELTIME
		expected string
	}{
		{
			name:     "1 second as microseconds",
			input:    RELTIME(1 * time.Second),
			expected: `"1000000"`,
		},
		{
			name:     "100 microseconds",
			input:    RELTIME(100 * time.Microsecond),
			expected: `"100"`,
		},
		{
			name:     "1 hour as microseconds",
			input:    RELTIME(1 * time.Hour),
			expected: `"3600000000"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestJsonCodec_Unmarshal_RELTIME(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		json     string
		expected RELTIME
	}{
		{
			name:     "microseconds from string",
			json:     `"1000000"`,
			expected: RELTIME(1 * time.Second),
		},
		{
			name:     "microseconds from number",
			json:     `100`,
			expected: RELTIME(100 * time.Microsecond),
		},
		{
			name:     "large value",
			json:     `"3600000000"`,
			expected: RELTIME(1 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result RELTIME
			err := codec.Unmarshal([]byte(tt.json), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonCodec_Unmarshal_RELTIME_FromMap(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		json     string
		expected RELTIME
	}{
		{
			name:     "microseconds from map with number",
			json:     `{"microseconds": 1000000}`,
			expected: RELTIME(1 * time.Second),
		},
		{
			name:     "microseconds from map with string",
			json:     `{"microseconds": "5000000"}`,
			expected: RELTIME(5 * time.Second),
		},
		{
			name:     "microseconds from map with large value",
			json:     `{"microseconds": 3600000000}`,
			expected: RELTIME(1 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result RELTIME
			err := codec.Unmarshal([]byte(tt.json), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonCodec_Unmarshal_RELTIME_FromMap_InvalidFormat(t *testing.T) {
	codec := NewJsonCodec()

	// Map without microseconds field should fail
	json := `{"other_field": 1000000}`
	var result RELTIME
	err := codec.Unmarshal([]byte(json), &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid map format for RELTIME")
}

func TestJsonCodec_Marshal_SET(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		input    SET
		expected string
	}{
		{
			name:     "SET of strings",
			input:    SET{"item1", "item2", "item3"},
			expected: `["item1","item2","item3"]`,
		},
		{
			name:     "SET of numbers",
			input:    SET{1, 2, 3},
			expected: `["1","2","3"]`,
		},
		{
			name:     "empty SET",
			input:    SET{},
			expected: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := codec.Marshal(tt.input)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestJsonCodec_Unmarshal_SET(t *testing.T) {
	codec := NewJsonCodec()

	tests := []struct {
		name     string
		json     string
		expected SET
	}{
		{
			name:     "SET from array",
			json:     `["item1","item2","item3"]`,
			expected: SET{"item1", "item2", "item3"},
		},
		{
			name:     "empty SET",
			json:     `[]`,
			expected: SET{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result SET
			err := codec.Unmarshal([]byte(tt.json), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJsonCodec_RoundTrip_RELTIME_SET(t *testing.T) {
	codec := NewJsonCodec()

	originalReltime := RELTIME(5 * time.Minute)
	jsonBytes, err := codec.Marshal(originalReltime)
	require.NoError(t, err)

	var resultReltime RELTIME
	err = codec.Unmarshal(jsonBytes, &resultReltime)
	require.NoError(t, err)
	assert.Equal(t, originalReltime, resultReltime)

	originalSet := SET{"a", "b", "c"}
	jsonBytes, err = codec.Marshal(originalSet)
	require.NoError(t, err)

	var resultSet SET
	err = codec.Unmarshal(jsonBytes, &resultSet)
	require.NoError(t, err)
	assert.Equal(t, originalSet, resultSet)
}

// func TestJsonCodec_Marshal_TUPLE2(t *testing.T) {
// 	codec := NewJsonCodec()

// 	t.Run("TUPLE2 of strings", func(t *testing.T) {
// 		input := TUPLE2[string, string]{First: "hello", Second: "world"}
// 		result, err := codec.Marshal(input)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, `{"_1":"hello","_2":"world"}`, string(result))
// 	})

// 	t.Run("TUPLE2 of mixed types", func(t *testing.T) {
// 		input := TUPLE2[INT64, TEXT]{First: INT64(42), Second: TEXT("test")}
// 		result, err := codec.Marshal(input)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, `{"_1":"42","_2":"test"}`, string(result))
// 	})

// 	t.Run("TUPLE2 of numbers", func(t *testing.T) {
// 		input := TUPLE2[int, int]{First: 1, Second: 2}
// 		result, err := codec.Marshal(input)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, `{"_1":"1","_2":"2"}`, string(result))
// 	})
// }

// func TestJsonCodec_Unmarshal_TUPLE2(t *testing.T) {
// 	codec := NewJsonCodec()

// 	t.Run("TUPLE2 from object", func(t *testing.T) {
// 		var result TUPLE2[string, string]
// 		err := codec.Unmarshal([]byte(`{"_1":"hello","_2":"world"}`), &result)
// 		require.NoError(t, err)
// 		assert.Equal(t, TUPLE2[string, string]{First: "hello", Second: "world"}, result)
// 	})

// 	t.Run("TUPLE2 with numbers", func(t *testing.T) {
// 		var result TUPLE2[float64, float64]
// 		err := codec.Unmarshal([]byte(`{"_1":42,"_2":100}`), &result)
// 		require.NoError(t, err)
// 		assert.Equal(t, TUPLE2[float64, float64]{First: float64(42), Second: float64(100)}, result)
// 	})
// }

// func TestJsonCodec_RoundTrip_TUPLE2(t *testing.T) {
// 	codec := NewJsonCodec()

// 	original := TUPLE2[string, string]{
// 		First:  "first value",
// 		Second: "second value",
// 	}
// 	jsonBytes, err := codec.Marshal(original)
// 	require.NoError(t, err)

// 	var result TUPLE2[string, string]
// 	err = codec.Unmarshal(jsonBytes, &result)
// 	require.NoError(t, err)
// 	assert.Equal(t, original, result)
// }
