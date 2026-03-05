package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore[
	// Account Address type.
	ADDR chains.Hashable,
] interface {
	CheckEnabled(ctx context.Context, address ADDR) error
	EnabledAddresses(ctx context.Context) ([]ADDR, error)
}
