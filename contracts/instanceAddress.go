package contracts

import (
	"bytes"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"golang.org/x/crypto/sha3"

	ccip_common "github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
)

// InstanceAddressLength is the length in bytes of an InstanceAddress.
const InstanceAddressLength = 32

type RawInstanceAddress string

func RawInstanceAddressFromString(s string) RawInstanceAddress { return RawInstanceAddress(s) }

func (r RawInstanceAddress) String() string { return string(r) }

func (r RawInstanceAddress) InstanceAddress() InstanceAddress {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(r))

	return InstanceAddress(h.Sum(nil))
}

func (r RawInstanceAddress) Binding() ccip_common.RawInstanceAddress {
	return ccip_common.RawInstanceAddress{
		Unpack: types.TEXT(r),
	}
}

// InstanceAddress represents the 32 byte Keccak256 hash of an instance ID + owner party.
type InstanceAddress [InstanceAddressLength]byte

// BytesToInstanceAddress converts a byte slice to an InstanceAddress.
// If b is larger than InstanceAddressLength, b is cropped from the left.
func BytesToInstanceAddress(b []byte) InstanceAddress {
	var addr InstanceAddress
	addr.SetBytes(b)

	return addr
}

// HexToInstanceAddress converts a hex string to an InstanceAddress.
// s may be prefixed with "0x".
func HexToInstanceAddress(s string) InstanceAddress { return BytesToInstanceAddress(common.FromHex(s)) }

// Cmp compares two InstanceAddresses.
func (a InstanceAddress) Cmp(other InstanceAddress) int { return bytes.Compare(a[:], other[:]) }

// Bytes returns the byte slice representation of the InstanceAddress.
func (a InstanceAddress) Bytes() []byte { return a[:] }

// Hex returns the hex string representation of the InstanceAddress, prefixed with "0x".
func (a InstanceAddress) Hex() string { return hexutil.Encode(a[:]) }

// String returns the hex string representation of the InstanceAddress, prefixed with "0x".
func (a InstanceAddress) String() string { return a.Hex() }

var (
	instanceAddressT = reflect.TypeFor[InstanceAddress]()
)

// UnmarshalText parses an InstanceAddress in hex syntax.
func (a *InstanceAddress) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("InstanceID", text, a[:])
}

// UnmarshalJSON parses an InstanceAddress in hex syntax.
func (a *InstanceAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(instanceAddressT, input, a[:])
}

// MarshalText returns the hex representation of a.
func (a InstanceAddress) MarshalText() ([]byte, error) { return hexutil.Bytes(a[:]).MarshalText() }

// SetBytes sets the InstanceAddress to the value of b.
// If b is larger than InstanceAddressLength, b is cropped from the left.
func (a *InstanceAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-InstanceAddressLength:]
	}

	copy(a[InstanceAddressLength-len(b):], b)
}
