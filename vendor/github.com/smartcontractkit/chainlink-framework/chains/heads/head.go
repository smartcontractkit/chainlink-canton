package heads

import (
	"github.com/smartcontractkit/chainlink-framework/chains"
)

type Head[BLOCK_HASH chains.Hashable, CHAIN_ID chains.ID] interface {
	chains.Head[BLOCK_HASH]
	// ChainID returns the chain ID of the head.
	ChainID() CHAIN_ID
	// HasChainID returns true if the head has a chain ID.
	HasChainID() bool
	// IsValid returns true if the head is valid.
	IsValid() bool
}
