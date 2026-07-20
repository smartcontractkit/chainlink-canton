package devenv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

// Shared devenv test constants (message and token paths).

const (
	// EVMToCantonFinalityConfig is the minimum block-depth FTF (1 confirmation).
	EVMToCantonFinalityConfig = protocol.Finality(1)

	// CantonToEVMFeeAmount is the per-message CCIP fee budget in Amulet units for
	// message-only sends (200k gas). Kept low for prod-testnet compatibility.
	CantonToEVMFeeAmount int64 = 5

	// CantonToEVMTokenTransferFeeAmount is the per-message fee budget for token transfers
	// (500k execution gas), which quote ~127 Amulet per send in devenv. Used only by
	// token e2e/load tests; must cover one send and leave enough change for sequential sends.
	CantonToEVMTokenTransferFeeAmount int64 = 13

	// Canton token amounts use 10-decimal fixed point (e.g. 1_000_000_000 = 0.1).
	CantonFixedPointScale      int64 = 10_000_000_000
	CantonFixedPointToEVMScale int64 = 100_000_000 // fixedPoint * this = EVM wei

	// CantonToEVMTokenSequentialSends is how many token transfers the Canton→EVM e2e subtest sends.
	CantonToEVMTokenSequentialSends = 2
)

// EVMToCantonNoExecutionExecutor is Client.NO_EXECUTION_ADDRESS on EVM source chains.
// Canton dest execution is manual (ManuallyExecuteMessage); no EVM Executor supports Canton.
var EVMToCantonNoExecutionExecutor = protocol.UnknownAddress(
	common.HexToAddress("0xEBa517d200000000000000000000000000000000").Bytes(),
)
