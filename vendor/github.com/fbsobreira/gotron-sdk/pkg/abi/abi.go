package abi

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	eABI "github.com/ethereum/go-ethereum/accounts/abi"
	eCommon "github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"golang.org/x/crypto/sha3"
)

// Param list
type Param map[string]interface{}

// LoadFromJSON string into ABI data
func LoadFromJSON(jString string) ([]Param, error) {
	if len(jString) == 0 {
		return nil, nil
	}
	data := []Param{}
	err := json.Unmarshal([]byte(jString), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Selector of a function from its signature, e.g. transfer(address,uint256)
func Selector(functionSignature string) []byte {
	// hash functionSignature
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(functionSignature))
	b := hasher.Sum(nil)
	return b[:4]
}

func convertToAddress(v interface{}) (eCommon.Address, error) {
	switch v.(type) {
	case string:
		addr, err := address.StringToAddress(v.(string))
		if err != nil {
			return eCommon.Address{}, fmt.Errorf("invalid address %s: %+v", v.(string), err)
		}
		return eCommon.BytesToAddress(addr.Bytes()[len(addr.Bytes())-20:]), nil
	case address.Address:
		addr := v.(address.Address)
		return addr.EthAddress(), nil
	}
	return eCommon.Address{}, fmt.Errorf("invalid address %v", v)
}

func convertToInt(ty eABI.Type, v interface{}) interface{} {
	if ty.T == eABI.IntTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseInt(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = int8(tmp)
		case 16:
			v = int16(tmp)
		case 32:
			v = int32(tmp)
		case 64:
			v = int64(tmp)
		}
	} else if ty.T == eABI.UintTy && ty.Size <= 64 {
		tmp, _ := strconv.ParseUint(v.(string), 10, ty.Size)
		switch ty.Size {
		case 8:
			v = uint8(tmp)
		case 16:
			v = uint16(tmp)
		case 32:
			v = uint32(tmp)
		case 64:
			v = uint64(tmp)
		}
	} else {
		s := v.(string)
		// check for hex char
		if strings.HasPrefix(s, "0x") {
			v, _ = new(big.Int).SetString(s[2:], 16)
		} else {
			v, _ = new(big.Int).SetString(s, 10)
		}
	}
	return v
}

// GetPaddedParam from struct
func GetPaddedParam(params []any) ([]byte, error) {
	// check params len is even
	if len(params)%2 == 1 {
		return nil, fmt.Errorf("expected even number of params, got %d", len(params))
	}

	values := make([]interface{}, 0)
	arguments := eABI.Arguments{}

	for i := 0; i < len(params); i += 2 {
		k := params[i]
		v := params[i+1]

		kStr, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("invalid non-string type %+v", kStr)
		}

		ty, err := eABI.NewType(kStr, "", nil)
		if err != nil {
			return nil, fmt.Errorf("could not parse type %s: %w", kStr, err)
		}

		arguments = append(arguments,
			eABI.Argument{
				Name:    "",
				Type:    ty,
				Indexed: false,
			},
		)

		if ty.T == eABI.SliceTy || ty.T == eABI.ArrayTy {
			if ty.Elem.T == eABI.AddressTy {
				addressStrSlice, ok := v.([]string)
				if ok {
					// Handle base58 addresses
					updatedAddresses := []eCommon.Address{}
					for _, addressStr := range addressStrSlice {
						a, err := address.StringToAddress(addressStr)
						if err != nil {
							break
						}
						updatedAddresses = append(updatedAddresses, a.EthAddress())
					}
					if len(updatedAddresses) > 0 {
						if len(updatedAddresses) != len(addressStrSlice) {
							return nil, errors.New("failed to convert all base58 addresses")
						}
						v = updatedAddresses
					}
				}
				sdkAddressSlice, ok := v.([]address.Address)
				if ok {
					updatedAddresses := []eCommon.Address{}
					for _, sdkAddress := range sdkAddressSlice {
						updatedAddresses = append(updatedAddresses, sdkAddress.EthAddress())
					}
					v = updatedAddresses
				}
			} else if (ty.Elem.T == eABI.IntTy || ty.Elem.T == eABI.UintTy) &&
				ty.Elem.Size > 64 {
				tmpSlice, ok := v.([]string)
				if ok {
					tmp := make([]*big.Int, 0)
					for i := range tmpSlice {
						var value *big.Int
						// check for hex char
						if strings.HasPrefix(tmpSlice[i], "0x") {
							value, _ = new(big.Int).SetString(tmpSlice[i][2:], 16)
						} else {
							value, _ = new(big.Int).SetString(tmpSlice[i], 10)
						}
						tmp = append(tmp, value)
					}
					v = tmp
				}
			}
		} else if ty.T == eABI.AddressTy {
			if v, err = convertToAddress(v); err != nil {
				return nil, err
			}
		} else if (ty.T == eABI.IntTy || ty.T == eABI.UintTy) && reflect.TypeOf(v).Kind() == reflect.String {
			v = convertToInt(ty, v)
		} else if ty.T == eABI.BytesTy || ty.T == eABI.FixedBytesTy {
			if data, ok := v.(string); ok {
				converted, err := convertStringToBytes(ty, data)
				if err != nil {
					v = converted
				}
			}
		}

		values = append(values, v)
	}

	// convert params to bytes
	return arguments.PackValues(values)
}

func processJSONArray(input string) ([][]byte, error) {
	var jsonArray []string
	err := json.Unmarshal([]byte(input), &jsonArray)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON array from string %+v: %w", input, err)
	}
	bytesSlice := make([][]byte, len(jsonArray))
	for i, elem := range jsonArray {
		bytesSlice[i], err = hex.DecodeString(strings.TrimPrefix(elem, "0x"))
		if err != nil {
			return nil, fmt.Errorf("error decoding byte string from element [%d] %+v: %w", i, elem, err)
		}
	}
	return bytesSlice, nil
}

func convertStringToBytes(ty eABI.Type, data string) (interface{}, error) {
	// if string
	// convert from hex string
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		// try with base64
		dataBytes, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
	}
	// if array and size == 0
	if ty.T == eABI.BytesTy || ty.Size == 0 {
		return dataBytes, nil
	}
	if len(dataBytes) != ty.Size {
		return nil, fmt.Errorf("invalid size: %d/%d", ty.Size, len(dataBytes))
	}
	switch ty.Size {
	case 1:
		value := [1]byte{}
		copy(value[:], dataBytes[:1])
		return value, nil
	case 2:
		value := [2]byte{}
		copy(value[:], dataBytes[:2])
		return value, nil
	case 8:
		value := [8]byte{}
		copy(value[:], dataBytes[:8])
		return value, nil
	case 16:
		value := [16]byte{}
		copy(value[:], dataBytes[:16])
		return value, nil
	case 32:
		value := [32]byte{}
		copy(value[:], dataBytes[:32])
		return value, nil
	}
	return nil, errors.New("failed to convert")
}

// Pack data into bytes
func Pack(method string, params []any) ([]byte, error) {
	if params == nil {
		params = []any{}
	}
	selector := Selector(method)
	pBytes, err := GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	selector = append(selector, pBytes...)
	return selector, nil
}
