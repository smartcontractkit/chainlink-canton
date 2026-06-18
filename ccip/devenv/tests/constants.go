package tests

// Shared devenv test constants (message and token paths).

const (
	// CantonToEVMFeeAmount is the per-message CCIP fee budget in Amulet units for
	// message-only sends (200k gas). Kept low for prod-testnet compatibility.
	CantonToEVMFeeAmount int64 = 50

	// CantonToEVMTokenTransferFeeAmount is the per-message fee budget for token transfers
	// (500k execution gas), which quote ~127 Amulet per send in devenv. Used only by
	// token e2e/load tests; must cover one send and leave enough change for sequential sends.
	CantonToEVMTokenTransferFeeAmount int64 = 130

	// EVMDecimalsScale converts Canton token amounts to EVM 18-decimal balance units.
	EVMDecimalsScale int64 = 1_000_000_000_000_000_000

	// CantonToEVMTokenSequentialSends is how many token transfers the Canton→EVM e2e subtest sends.
	CantonToEVMTokenSequentialSends = 2
)
