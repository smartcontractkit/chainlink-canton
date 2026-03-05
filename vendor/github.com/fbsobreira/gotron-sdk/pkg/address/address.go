package address

import (
	"bytes"
	"crypto/ecdsa"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"math/big"

	eCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected byte length of the address
	AddressLength = 21
	// AddressLengthBase58 is the expected length of the address in base58format
	AddressLengthBase58 = 34
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

var (
	// Tron zero address - https://developers.tron.network/docs/faq#3-what-is-the-destruction-address-of-tron
	ZeroAddress = Address(append([]byte{TronBytePrefix}, make([]byte, AddressLength-1)...))
)

// Address represents the 21 byte address of an Tron account.
type Address []byte

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex gets Tron hex representation of the address (21 byte, 41-prefixed)
func (a Address) Hex() string {
	return common.BytesToHexString(a[:])[2:]
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address {
	id := b.Bytes()
	base := bytes.Repeat([]byte{0}, AddressLength-len(id))
	return append(base, id...)
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) (Address, error) {
	addr, err := common.FromHex(s)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// Base58ToAddress returns Address with byte values of s.
func Base58ToAddress(s string) (Address, error) {
	addr, err := common.DecodeCheck(s)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// Base64ToAddress returns Address with byte values of s.
func Base64ToAddress(s string) (Address, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return Address(decoded), nil
}

func EVMAddressToAddress(a eCommon.Address) Address {
	return Address(append([]byte{TronBytePrefix}, a[:]...))
}

func StringToAddress(s string) (Address, error) {
	// if evm address format (0x prefixed hex)
	if eCommon.IsHexAddress(s) {
		return EVMAddressToAddress(eCommon.HexToAddress(s)), nil
	}

	// if hex address format (41 prefixed hex, 21 bytes)
	if len(s) == AddressLength*2 && s[:2] == "41" {
		return HexToAddress(s)
	}

	// if base58 address format (T prefixed)
	if len(s) == AddressLengthBase58 && s[0] == 'T' {
		return Base58ToAddress(s)
	}

	return nil, fmt.Errorf("invalid address format: %s", s)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}

	// This is actually an invalid address, since all TRON addresses should start with `TronBytePrefix`.
	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}

	return common.EncodeCheck(a.Bytes())
}

// PubkeyToAddress returns address from ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	address := crypto.PubkeyToAddress(p)

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return addressTron
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	*a = Address(srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return []byte(a), nil
}

func (a Address) EthAddress() eCommon.Address {
	return eCommon.BytesToAddress(a.Bytes()[1:])
}

// MarshalJSON implements the json.Marshaler interface.
// This marshals the address into Base58.
func (a Address) MarshalJSON() ([]byte, error) {
	if len(a) == 0 {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, a.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *Address) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid address format")
	}

	str := string(data[1 : len(data)-1])
	if str == "" {
		*a = nil
		return nil
	}

	// If string starts with 'T', treat as base58
	if str[0] == 'T' {
		addr, err := Base58ToAddress(str)
		if err != nil {
			return fmt.Errorf("invalid base58 address: %w", err)
		}
		*a = addr
		return nil
	}

	// Otherwise treat as hex
	addr, err := HexToAddress(str)
	if err != nil {
		return fmt.Errorf("invalid Tron hex address: %w", err)
	}
	*a = addr
	return nil
}
