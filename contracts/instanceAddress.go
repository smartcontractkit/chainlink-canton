package contracts

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
)

type RawInstanceAddress string

func NewRawInstanceAddress(instanceId InstanceID, owner types.PARTY) RawInstanceAddress {
	return RawInstanceAddress(fmt.Sprintf("%s@%s", instanceId, owner))
}

func RawInstanceAddressFromString(s string) (RawInstanceAddress, error) {
	// Validate that s is a valid raw instance address by checking if it contains exactly one '@'
	split := strings.Split(s, "@")
	if len(split) != 2 {
		return "", errors.New("invalid raw instance address: must contain exactly one '@'")
	}
	for _, s := range split {
		if len(s) == 0 {
			return "", errors.New("invalid raw instance address: parts cannot be empty")
		}
	}

	return RawInstanceAddress(s), nil
}

func (r RawInstanceAddress) String() string { return string(r) }

func (r RawInstanceAddress) InstanceID() string { return strings.Split(string(r), "@")[0] }

func (r RawInstanceAddress) Owner() string { return strings.Split(string(r), "@")[1] }

func (r RawInstanceAddress) InstanceAddress() InstanceAddress {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(r))

	return InstanceAddress(h.Sum(nil))
}

func (r RawInstanceAddress) Binding() chainlinkapi.RawInstanceAddress {
	return chainlinkapi.RawInstanceAddress{
		Unpack: types.TEXT(r),
	}
}

// InstanceAddressLength is the length in bytes of an InstanceAddress.
const InstanceAddressLength = 32

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

var instanceAddressT = reflect.TypeFor[InstanceAddress]()

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

const HashedPartyLength = 32

// HashedParty represents the 32 byte Keccak256 hash of a party string.
type HashedParty [HashedPartyLength]byte

func BytesToHashedParty(b []byte) HashedParty {
	var party HashedParty
	party.SetBytes(b)

	return party
}

// HashedPartyFromString computes the Keccak256 hash of the given party string and returns it as a HashedParty.
func HashedPartyFromString(party string) HashedParty {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(party))

	return HashedParty(h.Sum(nil))
}

// HexToHashedParty converts a hex string to a HashedParty.
// s may be prefixed with "0x".
func HexToHashedParty(s string) HashedParty { return BytesToHashedParty(common.FromHex(s)) }

// Cmp compares two HashedParties.
func (p HashedParty) Cmp(other HashedParty) int { return bytes.Compare(p[:], other[:]) }

// Bytes returns the byte slice representation of the HashedParty.
func (p HashedParty) Bytes() []byte { return p[:] }

// Hex returns the hex string representation of the HashedParty, prefixed with "0x".
func (p HashedParty) Hex() string { return hexutil.Encode(p[:]) }

// String returns the hex string representation of the HashedParty, prefixed with "0x".
func (p HashedParty) String() string { return p.Hex() }

var hashedPartyT = reflect.TypeFor[HashedParty]()

// UnmarshalText parses a HashedParty in hex syntax.
func (p *HashedParty) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("HashedParty", text, p[:])
}

// UnmarshalJSON parses a HashedParty in hex syntax.
func (p *HashedParty) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashedPartyT, input, p[:])
}

// MarshalText returns the hex representation of p.
func (p HashedParty) MarshalText() ([]byte, error) { return hexutil.Bytes(p[:]).MarshalText() }

// SetBytes sets the HashedParty to the value of b.
// If b is larger than HashedPartyLength, b is cropped from the left.
func (p *HashedParty) SetBytes(b []byte) {
	if len(b) > len(p) {
		b = b[len(b)-HashedPartyLength:]
	}

	copy(p[HashedPartyLength-len(b):], b)
}

// EncodedInstrumentIDLength is the length in bytes of an EncodedInstrumentID
const EncodedInstrumentIDLength = 32

// EncodedInstrumentID represents the 32 byte Keccak256 hash of an Instrument ID.
type EncodedInstrumentID [EncodedInstrumentIDLength]byte

func BytesToEncodedInstrumentID(b []byte) EncodedInstrumentID {
	var instrumentID EncodedInstrumentID
	instrumentID.SetBytes(b)

	return instrumentID
}

// HexToEncodedInstrumentID converts a hex string to an EncodedInstrumentID.
// s may be prefixed with "0x"
func HexToEncodedInstrumentID(s string) EncodedInstrumentID {
	return BytesToEncodedInstrumentID(common.FromHex(s))
}

// EncodeInstrumentID takes the given InstrumentId and returns its EncodedInstrumentID.
func EncodeInstrumentID(id splice_api_token_holding_v1.InstrumentId) EncodedInstrumentID {
	h := sha3.NewLegacyKeccak256()
	fmt.Fprintf(h, "%s@%s", id.Id, id.Admin)

	return EncodedInstrumentID(h.Sum(nil))
}

func (i EncodedInstrumentID) Bytes() []byte { return i[:] }

func (i EncodedInstrumentID) Hex() string { return hexutil.Encode(i[:]) }

func (i EncodedInstrumentID) String() string { return i.Hex() }

func (i *EncodedInstrumentID) SetBytes(b []byte) {
	if len(b) > len(i) {
		b = b[len(b)-EncodedInstrumentIDLength:]
	}

	copy(i[EncodedInstrumentIDLength-len(b):], b)
}
