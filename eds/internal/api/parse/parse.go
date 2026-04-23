package parse

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func Uint64Checked(s string) (uint64, error) {
	// Since uint64s are represented by Numeric 0, which will be serialized as "1.0",
	// cannot use *big.Int directly - instead, parse as big.Rat and then convert to int.
	val, ok := new(big.Rat).SetString(s)
	if !ok {
		return 0, fmt.Errorf("failed to parse uint64: %s", s)
	}
	if !val.IsInt() {
		return 0, fmt.Errorf("is not an int: %s", s)
	}

	num := val.Num()
	if !num.IsUint64() {
		return 0, fmt.Errorf("is not uint64: %s", s)
	}

	return num.Uint64(), nil
}

func RawInstanceAddress(a mcms.RawInstanceAddress) (contracts.RawInstanceAddress, error) {
	address, err := contracts.RawInstanceAddressFromString(string(a.Unpack))
	if err != nil {
		return contracts.RawInstanceAddress(""), fmt.Errorf("failed to parse raw instance address: %w", err)
	}

	return address, nil
}

func RawInstanceAddressList(addrs []mcms.RawInstanceAddress) ([]contracts.RawInstanceAddress, error) {
	addresses := make([]contracts.RawInstanceAddress, len(addrs))
	for i, addr := range addrs {
		parsedAddr, err := RawInstanceAddress(addr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse raw instance address at index %d: %w", i, err)
		}
		addresses[i] = parsedAddr
	}

	return addresses, nil
}
