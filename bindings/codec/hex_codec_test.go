package codec

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-canton/bindings/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHexCodec_EncodeUint8(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(uint8(10))
	require.NoError(t, err)
	assert.Equal(t, "0a", result)
}

func TestHexCodec_EncodeUint8_Zero(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(uint8(0))
	require.NoError(t, err)
	assert.Equal(t, "00", result)
}

func TestHexCodec_EncodeUint8_Max(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(uint8(255))
	require.NoError(t, err)
	assert.Equal(t, "ff", result)
}

func TestHexCodec_EncodeUint32(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(uint32(10))
	require.NoError(t, err)
	assert.Equal(t, "0000000a", result)
}

func TestHexCodec_EncodeUint32_Large(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(uint32(0x12345678))
	require.NoError(t, err)
	assert.Equal(t, "12345678", result)
}

func TestHexCodec_EncodeInt(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(int(10))
	require.NoError(t, err)
	assert.Equal(t, "0000000a", result) // int encodes as int32
}

func TestHexCodec_EncodeInt64(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(int64(10))
	require.NoError(t, err)
	assert.Equal(t, "000000000000000a", result)
}

func TestHexCodec_EncodeText(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal("foo")
	require.NoError(t, err)
	// len=3 + "foo" in hex
	assert.Equal(t, "03666f6f", result)
}

func TestHexCodec_EncodeText_Empty(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal("")
	require.NoError(t, err)
	assert.Equal(t, "00", result) // len=0
}

func TestHexCodec_EncodeText_HelloWorld(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal("hello")
	require.NoError(t, err)
	// len=5 + "hello" in hex (68656c6c6f)
	assert.Equal(t, "0568656c6c6f", result)
}

func TestHexCodec_EncodeBool_True(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(true)
	require.NoError(t, err)
	assert.Equal(t, "01", result)
}

func TestHexCodec_EncodeBool_False(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(false)
	require.NoError(t, err)
	assert.Equal(t, "00", result)
}

func TestHexCodec_EncodeDAMLText(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.TEXT("bar"))
	require.NoError(t, err)
	assert.Equal(t, "03626172", result) // len=3 + "bar"
}

func TestHexCodec_EncodeDAMLInt64(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.INT64(256))
	require.NoError(t, err)
	assert.Equal(t, "0000000000000100", result)
}

func TestHexCodec_EncodeDAMLBool(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.BOOL(true))
	require.NoError(t, err)
	assert.Equal(t, "01", result)
}

func TestHexCodec_EncodeDAMLParty(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.PARTY("alice"))
	require.NoError(t, err)
	assert.Equal(t, "05616c696365", result) // len=5 + "alice"
}

func TestHexCodec_EncodeSlice(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal([]uint32{1, 2, 3})
	require.NoError(t, err)
	// len=3 + 3 * uint32 (4 bytes each)
	assert.Equal(t, "03000000010000000200000003", result)
}

func TestHexCodec_EncodeSlice_Empty(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal([]uint32{})
	require.NoError(t, err)
	assert.Equal(t, "00", result) // len=0
}

func TestHexCodec_EncodeStruct(t *testing.T) {
	type SimpleStruct struct {
		Name  string
		Value int
	}

	c := NewHexCodec()
	result, err := c.Marshal(SimpleStruct{
		Name:  "test",
		Value: 42,
	})
	require.NoError(t, err)
	// Name: len=4 + "test" (04 74657374)
	// Value: int32 = 42 (0000002a)
	assert.Equal(t, "0474657374"+"0000002a", result)
}

func TestHexCodec_EncodeNestedStruct(t *testing.T) {
	type Inner struct {
		X int
		Y int
	}
	type Outer struct {
		Name  string
		Inner Inner
	}

	c := NewHexCodec()
	result, err := c.Marshal(Outer{
		Name: "a",
		Inner: Inner{
			X: 1,
			Y: 2,
		},
	})
	require.NoError(t, err)
	// Name: len=1 + "a" (0161)
	// Inner.X: int32 = 1 (00000001)
	// Inner.Y: int32 = 2 (00000002)
	assert.Equal(t, "0161"+"00000001"+"00000002", result)
}

func TestHexCodec_EncodeStructWithSlice(t *testing.T) {
	type Config struct {
		Name   string
		Values []int
	}

	c := NewHexCodec()
	result, err := c.Marshal(Config{
		Name:   "cfg",
		Values: []int{10, 20},
	})
	require.NoError(t, err)
	// Name: len=3 + "cfg" (03 636667)
	// Values: len=2 + 10 + 20 (02 0000000a 00000014)
	assert.Equal(t, "03636667"+"02"+"0000000a"+"00000014", result)
}

// Decode tests

func TestHexCodec_DecodeUint8(t *testing.T) {
	c := NewHexCodec()
	var result uint8
	err := c.Unmarshal("0a", &result)
	require.NoError(t, err)
	assert.Equal(t, uint8(10), result)
}

func TestHexCodec_DecodeUint32(t *testing.T) {
	c := NewHexCodec()
	var result uint32
	err := c.Unmarshal("0000000a", &result)
	require.NoError(t, err)
	assert.Equal(t, uint32(10), result)
}

func TestHexCodec_DecodeInt64(t *testing.T) {
	c := NewHexCodec()
	var result int64
	err := c.Unmarshal("000000000000000a", &result)
	require.NoError(t, err)
	assert.Equal(t, int64(10), result)
}

func TestHexCodec_DecodeText(t *testing.T) {
	c := NewHexCodec()
	var result string
	err := c.Unmarshal("03666f6f", &result)
	require.NoError(t, err)
	assert.Equal(t, "foo", result)
}

func TestHexCodec_DecodeBool_True(t *testing.T) {
	c := NewHexCodec()
	var result bool
	err := c.Unmarshal("01", &result)
	require.NoError(t, err)
	assert.True(t, result)
}

func TestHexCodec_DecodeBool_False(t *testing.T) {
	c := NewHexCodec()
	var result bool
	err := c.Unmarshal("00", &result)
	require.NoError(t, err)
	assert.False(t, result)
}

func TestHexCodec_DecodeDAMLText(t *testing.T) {
	c := NewHexCodec()
	var result types.TEXT
	err := c.Unmarshal("03626172", &result)
	require.NoError(t, err)
	assert.Equal(t, types.TEXT("bar"), result)
}

func TestHexCodec_DecodeDAMLInt64(t *testing.T) {
	c := NewHexCodec()
	var result types.INT64
	err := c.Unmarshal("0000000000000100", &result)
	require.NoError(t, err)
	assert.Equal(t, types.INT64(256), result)
}

func TestHexCodec_DecodeDAMLBool(t *testing.T) {
	c := NewHexCodec()
	var result types.BOOL
	err := c.Unmarshal("01", &result)
	require.NoError(t, err)
	assert.Equal(t, types.BOOL(true), result)
}

func TestHexCodec_DecodeSlice(t *testing.T) {
	c := NewHexCodec()
	var result []uint32
	err := c.Unmarshal("03000000010000000200000003", &result)
	require.NoError(t, err)
	assert.Equal(t, []uint32{1, 2, 3}, result)
}

func TestHexCodec_DecodeStruct(t *testing.T) {
	type SimpleStruct struct {
		Name  string
		Value int
	}

	c := NewHexCodec()
	var result SimpleStruct
	err := c.Unmarshal("04746573740000002a", &result)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

func TestHexCodec_DecodeNestedStruct(t *testing.T) {
	type Inner struct {
		X int
		Y int
	}
	type Outer struct {
		Name  string
		Inner Inner
	}

	c := NewHexCodec()
	var result Outer
	err := c.Unmarshal("01610000000100000002", &result)
	require.NoError(t, err)
	assert.Equal(t, "a", result.Name)
	assert.Equal(t, 1, result.Inner.X)
	assert.Equal(t, 2, result.Inner.Y)
}

func TestHexCodec_RoundTrip_Struct(t *testing.T) {
	type Config struct {
		Name      string
		Count     int
		Enabled   bool
		Threshold int64
	}

	original := Config{
		Name:      "myconfig",
		Count:     100,
		Enabled:   true,
		Threshold: 1000000,
	}

	c := NewHexCodec()
	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded Config
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestHexCodec_RoundTrip_SliceOfStructs(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}
	type Container struct {
		Items []Item
	}

	original := Container{
		Items: []Item{
			{ID: 1, Name: "first"},
			{ID: 2, Name: "second"},
		},
	}

	c := NewHexCodec()
	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded Container
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

// Test matching Canton format (based on mcms_crypto.go EncodeSetConfigParams)
func TestHexCodec_SignerInfo_Format(t *testing.T) {
	// SignerInfo in Canton format:
	// - signerAddress: text (len + bytes)
	// - signerIndex: int32 (4 bytes)
	// - signerGroup: int32 (4 bytes)
	type SignerInfo struct {
		SignerAddress string
		SignerIndex   int
		SignerGroup   int
	}

	c := NewHexCodec()
	signer := SignerInfo{
		SignerAddress: "abc",
		SignerIndex:   0,
		SignerGroup:   0,
	}

	result, err := c.Marshal(signer)
	require.NoError(t, err)

	// Expected: 03 616263 00000000 00000000
	// - 03: length of "abc"
	// - 616263: "abc" in hex
	// - 00000000: signerIndex = 0 as int32
	// - 00000000: signerGroup = 0 as int32
	assert.Equal(t, "036162630000000000000000", result)
}

func TestHexCodec_SetConfigParams_Format(t *testing.T) {
	// SetConfigParams in Canton format (simplified):
	// - signers: list of SignerInfo
	// - groupQuorums: list of int
	// - groupParents: list of int
	// - clearRoot: bool
	type SignerInfo struct {
		SignerAddress string
		SignerIndex   int
		SignerGroup   int
	}
	type SetConfigParams struct {
		Signers      []SignerInfo
		GroupQuorums []int
		GroupParents []int
		ClearRoot    bool
	}

	c := NewHexCodec()
	params := SetConfigParams{
		Signers: []SignerInfo{
			{SignerAddress: "ab", SignerIndex: 0, SignerGroup: 0},
		},
		GroupQuorums: []int{1, 0},
		GroupParents: []int{0, 0},
		ClearRoot:    false,
	}

	result, err := c.Marshal(params)
	require.NoError(t, err)

	// Expected breakdown:
	// Signers: 01 (count=1) + 02 6162 00000000 00000000 (signerInfo)
	// GroupQuorums: 02 (count=2) + 00000001 00000000
	// GroupParents: 02 (count=2) + 00000000 00000000
	// ClearRoot: 00
	expected := "01" + // signers count
		"02" + "6162" + "00000000" + "00000000" + // signer[0]
		"02" + "00000001" + "00000000" + // group quorums
		"02" + "00000000" + "00000000" + // group parents
		"00" // clearRoot
	assert.Equal(t, expected, result)
}

func TestHexCodec_InvalidHex(t *testing.T) {
	c := NewHexCodec()
	var result uint8
	err := c.Unmarshal("zz", &result)
	assert.Error(t, err)
}

func TestHexCodec_NotEnoughData(t *testing.T) {
	c := NewHexCodec()
	var result uint32
	err := c.Unmarshal("00", &result) // Need 4 bytes, only have 1
	assert.Error(t, err)
}

func TestHexCodec_NonPointerTarget(t *testing.T) {
	c := NewHexCodec()
	var result uint8
	err := c.Unmarshal("0a", result) // Should be &result
	assert.Error(t, err)
}

// Test hex:"bytes" struct tag
func TestHexCodec_HexBytesTag_Encode(t *testing.T) {
	type SignerInfo struct {
		SignerAddress string `hex:"bytes"` // hex-encoded address
		SignerIndex   int
		SignerGroup   int
	}

	c := NewHexCodec()
	signer := SignerInfo{
		SignerAddress: "abcd", // hex string representing 2 bytes
		SignerIndex:   0,
		SignerGroup:   0,
	}

	result, err := c.Marshal(signer)
	require.NoError(t, err)

	// Expected: 02 (len=2 bytes) + abcd (raw bytes) + 00000000 (index) + 00000000 (group)
	assert.Equal(t, "02abcd0000000000000000", result)
}

func TestHexCodec_HexBytesTag_Decode(t *testing.T) {
	type SignerInfo struct {
		SignerAddress string `hex:"bytes"`
		SignerIndex   int
		SignerGroup   int
	}

	c := NewHexCodec()
	var result SignerInfo
	// Input: 02 (len=2) + abcd (bytes) + 00000000 (index) + 00000000 (group)
	err := c.Unmarshal("02abcd0000000000000000", &result)
	require.NoError(t, err)

	assert.Equal(t, "abcd", result.SignerAddress)
	assert.Equal(t, 0, result.SignerIndex)
	assert.Equal(t, 0, result.SignerGroup)
}

func TestHexCodec_HexBytesTag_RoundTrip(t *testing.T) {
	type SignerInfo struct {
		SignerAddress string `hex:"bytes"`
		SignerIndex   int
		SignerGroup   int
	}

	c := NewHexCodec()
	original := SignerInfo{
		SignerAddress: "1375dc8a4c1476e6628b03216545e5cdcbff3f84",
		SignerIndex:   5,
		SignerGroup:   2,
	}

	encoded, err := c.Marshal(original)
	require.NoError(t, err)
	t.Logf("Encoded: %s", encoded)

	var decoded SignerInfo
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestHexCodec_HexBytesTag_SetConfigParams(t *testing.T) {
	// This matches the Canton SetConfigParams structure
	type SignerInfo struct {
		SignerAddress string `hex:"bytes"`
		SignerIndex   int
		SignerGroup   int
	}
	type SetConfigParams struct {
		Signers      []SignerInfo
		GroupQuorums []int
		GroupParents []int
		ClearRoot    bool
	}

	c := NewHexCodec()
	params := SetConfigParams{
		Signers: []SignerInfo{
			{SignerAddress: "abcd", SignerIndex: 0, SignerGroup: 0},
		},
		GroupQuorums: []int{1, 0},
		GroupParents: []int{0, 0},
		ClearRoot:    false,
	}

	result, err := c.Marshal(params)
	require.NoError(t, err)

	// Expected (matching manual encoder):
	// 01 = 1 signer
	// 02 = addr len (2 bytes for "abcd" decoded)
	// abcd = addr bytes
	// 00000000 = signerIndex
	// 00000000 = signerGroup
	// 02 = 2 quorums
	// 00000001 00000000 = quorum values
	// 02 = 2 parents
	// 00000000 00000000 = parent values
	// 00 = clearRoot
	expected := "01" + "02" + "abcd" + "00000000" + "00000000" +
		"02" + "00000001" + "00000000" +
		"02" + "00000000" + "00000000" +
		"00"
	assert.Equal(t, expected, result)
}

func TestHexCodec_HexBytesTag_InvalidHex(t *testing.T) {
	type BadStruct struct {
		Address string `hex:"bytes"`
	}

	c := NewHexCodec()
	bad := BadStruct{Address: "zzzz"} // invalid hex

	_, err := c.Marshal(bad)
	assert.Error(t, err)
}

// Test DAML enum/variant type encoding
// DAML enums implement GetEnumConstructor() and should be encoded as text

// TestRole mimics the generated Role enum type from DAML
type TestRole string

const (
	TestRoleProposer  TestRole = "Proposer"
	TestRoleCanceller TestRole = "Canceller"
	TestRoleBypasser  TestRole = "Bypasser"
	TestRoleExecutor  TestRole = "Executor"
)

func (e TestRole) GetEnumConstructor() string { return string(e) }

func TestHexCodec_EncodeEnum(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(TestRoleProposer)
	require.NoError(t, err)
	// "Proposer" = 8 chars, so: 08 + "Proposer" in hex (50726f706f736572)
	assert.Equal(t, "0850726f706f736572", result)
}

func TestHexCodec_EncodeEnum_AllValues(t *testing.T) {
	c := NewHexCodec()

	tests := []struct {
		role     TestRole
		expected string
	}{
		{TestRoleProposer, "0850726f706f736572"},    // 08 + "Proposer"
		{TestRoleCanceller, "0943616e63656c6c6572"}, // 09 + "Canceller"
		{TestRoleBypasser, "084279706173736572"},    // 08 + "Bypasser"
		{TestRoleExecutor, "084578656375746f72"},    // 08 + "Executor"
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			result, err := c.Marshal(tt.role)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHexCodec_DecodeEnum(t *testing.T) {
	c := NewHexCodec()
	var result TestRole
	// Decode "Proposer" (08 + 50726f706f736572)
	err := c.Unmarshal("0850726f706f736572", &result)
	require.NoError(t, err)
	assert.Equal(t, TestRoleProposer, result)
}

func TestHexCodec_RoundTrip_Enum(t *testing.T) {
	c := NewHexCodec()

	roles := []TestRole{TestRoleProposer, TestRoleCanceller, TestRoleBypasser, TestRoleExecutor}
	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			encoded, err := c.Marshal(role)
			require.NoError(t, err)

			var decoded TestRole
			err = c.Unmarshal(encoded, &decoded)
			require.NoError(t, err)

			assert.Equal(t, role, decoded)
		})
	}
}

func TestHexCodec_EncodeStructWithEnum(t *testing.T) {
	type SetConfig struct {
		TargetRole TestRole
		ClearRoot  bool
	}

	c := NewHexCodec()
	config := SetConfig{
		TargetRole: TestRoleProposer,
		ClearRoot:  false,
	}

	result, err := c.Marshal(config)
	require.NoError(t, err)
	// TargetRole: 08 + "Proposer" (50726f706f736572)
	// ClearRoot: 00
	assert.Equal(t, "0850726f706f73657200", result)
}

func TestHexCodec_RoundTrip_StructWithEnum(t *testing.T) {
	type SetConfig struct {
		TargetRole TestRole
		ClearRoot  bool
	}

	c := NewHexCodec()
	original := SetConfig{
		TargetRole: TestRoleCanceller,
		ClearRoot:  true,
	}

	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded SetConfig
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

// Tests for new DAML types

func TestHexCodec_EncodeRELTIME(t *testing.T) {
	c := NewHexCodec()
	// 1 second = 1000000 microseconds
	result, err := c.Marshal(types.RELTIME(1000000 * 1000)) // 1 second as nanoseconds
	require.NoError(t, err)
	// 1000 microseconds = 0x03E8 (1000000/1000 since RELTIME is Duration in nanoseconds, we convert to microseconds)
	// Actually 1000000*1000 nanoseconds / 1000 = 1000000 microseconds
	assert.Equal(t, "00000000000f4240", result) // 1000000 in hex = 0xf4240
}

func TestHexCodec_DecodeRELTIME(t *testing.T) {
	c := NewHexCodec()
	var result types.RELTIME
	err := c.Unmarshal("00000000000f4240", &result)
	require.NoError(t, err)
	// 1000000 microseconds * time.Microsecond = 1 second
	assert.Equal(t, types.RELTIME(1000000*1000), result) // 1 second in nanoseconds
}

func TestHexCodec_RoundTrip_RELTIME(t *testing.T) {
	c := NewHexCodec()
	original := types.RELTIME(5 * 1000 * 1000 * 1000) // 5 seconds

	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded types.RELTIME
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestHexCodec_EncodeUNIT(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.UNIT{})
	require.NoError(t, err)
	assert.Equal(t, "", result) // UNIT encodes to empty
}

func TestHexCodec_DecodeUNIT(t *testing.T) {
	c := NewHexCodec()
	var result types.UNIT
	err := c.Unmarshal("", &result)
	require.NoError(t, err)
	assert.Equal(t, types.UNIT{}, result)
}

func TestHexCodec_EncodeTIMESTAMP(t *testing.T) {
	c := NewHexCodec()
	// Use a specific timestamp
	ts := types.TIMESTAMP(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	result, err := c.Marshal(ts)
	require.NoError(t, err)
	// Should encode as microseconds since epoch
	assert.Len(t, result, 16) // 8 bytes in hex
}

func TestHexCodec_RoundTrip_TIMESTAMP(t *testing.T) {
	c := NewHexCodec()
	// Note: TIMESTAMP only preserves microsecond precision
	original := types.TIMESTAMP(time.Date(2024, 6, 15, 12, 30, 45, 123000000, time.UTC))

	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded types.TIMESTAMP
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	// Compare truncated to microseconds
	assert.Equal(t, time.Time(original).UnixMicro(), time.Time(decoded).UnixMicro())
}

func TestHexCodec_EncodeDATE(t *testing.T) {
	c := NewHexCodec()
	// Unix epoch is day 0
	dt := types.DATE(time.Unix(86400, 0).UTC()) // Day 1
	result, err := c.Marshal(dt)
	require.NoError(t, err)
	assert.Equal(t, "00000001", result) // 1 day
}

func TestHexCodec_RoundTrip_DATE(t *testing.T) {
	c := NewHexCodec()
	// Jan 1, 2024 = some number of days since epoch
	original := types.DATE(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded types.DATE
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	// DATE only preserves day precision
	assert.Equal(t, time.Time(original).Unix()/86400, time.Time(decoded).Unix()/86400)
}

func TestHexCodec_EncodeDECIMAL(t *testing.T) {
	c := NewHexCodec()
	bigInt := big.NewInt(12345)
	result, err := c.Marshal(types.DECIMAL(bigInt))
	require.NoError(t, err)
	// "12345" = 5 chars, so: 05 + "12345" in hex (3132333435)
	assert.Equal(t, "053132333435", result)
}

func TestHexCodec_EncodeDECIMAL_Nil(t *testing.T) {
	c := NewHexCodec()
	result, err := c.Marshal(types.DECIMAL(nil))
	require.NoError(t, err)
	assert.Equal(t, "00", result) // empty string
}

func TestHexCodec_RoundTrip_DECIMAL(t *testing.T) {
	c := NewHexCodec()
	bigInt := big.NewInt(9876543210)
	original := types.DECIMAL(bigInt)

	encoded, err := c.Marshal(original)
	require.NoError(t, err)

	var decoded types.DECIMAL
	err = c.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, (*big.Int)(original).String(), (*big.Int)(decoded).String())
}

func TestHexCodec_EncodeSET(t *testing.T) {
	c := NewHexCodec()
	set := types.SET{uint32(1), uint32(2), uint32(3)}
	result, err := c.Marshal(set)
	require.NoError(t, err)
	// len=3 + 3 uint32s
	assert.Equal(t, "03000000010000000200000003", result)
}

func TestHexCodec_EncodeLIST(t *testing.T) {
	c := NewHexCodec()
	list := types.LIST{"foo", "bar"}
	result, err := c.Marshal(list)
	require.NoError(t, err)
	// len=2 + "foo" (03666f6f) + "bar" (03626172)
	assert.Equal(t, "0203666f6f03626172", result)
}

func TestHexCodec_EncodeTEXTMAP(t *testing.T) {
	c := NewHexCodec()
	// Note: map iteration order is not deterministic, so use single entry
	m := types.TEXTMAP{"key": "value"}
	result, err := c.Marshal(m)
	require.NoError(t, err)
	// len=1 + "key" (036b6579) + "value" (0576616c7565)
	assert.Equal(t, "01036b65790576616c7565", result)
}

func TestHexCodec_EncodeGENMAP(t *testing.T) {
	c := NewHexCodec()
	m := types.GENMAP{"num": int64(42)}
	result, err := c.Marshal(m)
	require.NoError(t, err)
	// len=1 + "num" (036e756d) + int64(42) (000000000000002a)
	assert.Equal(t, "01036e756d000000000000002a", result)
}

func TestHexCodec_EncodeMAP(t *testing.T) {
	c := NewHexCodec()
	m := types.MAP{"x": int64(1)}
	result, err := c.Marshal(m)
	require.NoError(t, err)
	// len=1 + "x" (0178) + int64(1)
	assert.Equal(t, "0101780000000000000001", result)
}

func TestHexCodec_EncodeTUPLE2(t *testing.T) {
	c := NewHexCodec()
	tuple := types.TUPLE2{First: int64(10), Second: int64(20)}
	result, err := c.Marshal(tuple)
	require.NoError(t, err)
	// int64(10) + int64(20)
	assert.Equal(t, "000000000000000a0000000000000014", result)
}

func TestHexCodec_EncodeTUPLE2_Mixed(t *testing.T) {
	c := NewHexCodec()
	tuple := types.TUPLE2{First: "hello", Second: int64(5)}
	result, err := c.Marshal(tuple)
	require.NoError(t, err)
	// "hello" (0568656c6c6f) + int64(5) (0000000000000005)
	assert.Equal(t, "0568656c6c6f0000000000000005", result)
}

// Test VARIANT interface encoding
type TestVariant struct {
	tag   string
	value interface{}
}

func (v TestVariant) GetVariantTag() string        { return v.tag }
func (v TestVariant) GetVariantValue() interface{} { return v.value }

func TestHexCodec_EncodeVARIANT(t *testing.T) {
	c := NewHexCodec()
	variant := TestVariant{tag: "Some", value: int64(42)}
	result, err := c.Marshal(variant)
	require.NoError(t, err)
	// "Some" (04536f6d65) + int64(42) (000000000000002a)
	assert.Equal(t, "04536f6d65000000000000002a", result)
}

func TestHexCodec_EncodeVARIANT_NilValue(t *testing.T) {
	c := NewHexCodec()
	variant := TestVariant{tag: "None", value: nil}
	result, err := c.Marshal(variant)
	require.NoError(t, err)
	// "None" (044e6f6e65) only
	assert.Equal(t, "044e6f6e65", result)
}
