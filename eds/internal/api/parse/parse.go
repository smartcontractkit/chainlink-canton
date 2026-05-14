package parse

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

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
