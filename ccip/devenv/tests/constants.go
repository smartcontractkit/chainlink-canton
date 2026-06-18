package tests

// Shared devenv test constants (message and token paths).

const (
	// CantonToEVMFeeAmount is the Canton CCIP send fee in Amulet units.
	CantonToEVMFeeAmount int64 = 50

	// EVMDecimalsScale converts Canton token amounts to EVM 18-decimal balance units.
	EVMDecimalsScale int64 = 1_000_000_000_000_000_000

	// CantonToEVMTokenSequentialSends is how many token transfers the Canton→EVM e2e subtest sends.
	CantonToEVMTokenSequentialSends = 2
)
