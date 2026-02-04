package contracts

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const InstanceAddressLength = 32

type InstanceAddress [InstanceAddressLength]byte

func BytesToInstanceAddress(b []byte) InstanceAddress {
	var addr InstanceAddress
	addr.SetBytes(b)

	return addr
}

// HexToInstanceAddress converts a hex string to an InstanceAddress.
// s may be prefixed with "0x".
func HexToInstanceAddress(s string) InstanceAddress { return BytesToInstanceAddress(common.FromHex(s)) }

func (a InstanceAddress) Cmp(other InstanceAddress) int { return bytes.Compare(a[:], other[:]) }

func (a InstanceAddress) Bytes() []byte { return a[:] }

func (a InstanceAddress) Hex() string { return hexutil.Encode(a[:]) }

func (a InstanceAddress) String() string { return a.Hex() }

func (a *InstanceAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-InstanceAddressLength:]
	}

	copy(a[InstanceAddressLength-len(b):], b)
}
