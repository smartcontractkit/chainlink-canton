package types

import (
	"context"

	"github.com/smartcontractkit/chainlink-framework/chains"
	"github.com/smartcontractkit/chainlink-framework/chains/fees"
)

// StuckTxDetector is used by the Confirmer to determine if any unconfirmed transactions are terminally stuck
type StuckTxDetector[CID chains.ID, ADDR chains.Hashable, THASH, BHASH chains.Hashable, SEQ chains.Sequence, FEE fees.Fee] interface {
	// Uses either a chain specific API or heuristic to determine if any unconfirmed transactions are terminally stuck. Returns only one transaction per enabled address.
	DetectStuckTransactions(ctx context.Context, enabledAddresses []ADDR, blockNum int64) ([]Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], error)
	// Loads the internal map that tracks the last block num a transaction was purged at using the DB state
	LoadPurgeBlockNumMap(ctx context.Context, addresses []ADDR) error
	// Sets the last purged block num after a transaction has been successfully purged with receipt
	SetPurgeBlockNum(fromAddress ADDR, blockNum int64)
	// Returns the error message to set in the transaction error field to mark it as terminally stuck
	StuckTxFatalError() string
}
