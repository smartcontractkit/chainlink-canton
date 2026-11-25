package accessors

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/protocol"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
)

var _ chainaccess.CCIPAPIParser = (*accessor)(nil)

// CCIPAPIParser returns the accessor itself.
func (a *accessor) CCIPAPIParser() (chainaccess.CCIPAPIParser, error) {
	return a, nil
}

// DecodeAddress accepts a raw id@party address or its 32-byte hashed hex form.
func (a *accessor) DecodeAddress(address string) (protocol.UnknownAddress, error) {
	if raw, err := contracts.RawInstanceAddressFromString(address); err == nil {
		return protocol.UnknownAddress(raw.InstanceAddress().Bytes()), nil
	}

	var hashed contracts.InstanceAddress
	if err := hashed.UnmarshalText([]byte(address)); err != nil {
		return nil, fmt.Errorf("invalid Canton address: %w", err)
	}

	return protocol.UnknownAddress(hashed.Bytes()), nil
}
