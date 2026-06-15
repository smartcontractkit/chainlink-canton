package tests

// Shared devenv test constants (message and token paths).

const (
	// OwnerParticipantIndex is Participants[0]: deploy, EDS, verifiers, fee aggregator.
	OwnerParticipantIndex = 0
	// ClientParticipantIndex is Participants[1]: send, receive, execute, token holdings.
	ClientParticipantIndex = 1

	// CantonToEVMFeeAmount is the Canton CCIP send fee in Amulet units.
	CantonToEVMFeeAmount int64 = 2_000

	// EVMDecimalsScale converts Canton token amounts to EVM 18-decimal balance units.
	EVMDecimalsScale int64 = 1_000_000_000_000_000_000

	// CantonToEVMTokenSequentialSends is how many token transfers the Canton→EVM e2e subtest sends.
	CantonToEVMTokenSequentialSends = 2
)
