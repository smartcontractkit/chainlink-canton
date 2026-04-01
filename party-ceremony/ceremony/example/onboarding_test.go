package example_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	ceremony "github.com/chainlink/canton-party-ceremony/ceremony/example"
)

// newDeps wraps the mock client in a CantonDeps value.
func newDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{Client: &ceremony.MockCantonClient{ParticipantID: participantID}, Logger: logger.Nop()}
}

// TestOnboardingSequence_HappyPath verifies that the full ceremony completes
// when all participants cooperate.  It checks the output fields and that the
// SequenceReport contains one child report per operation executed.
func TestOnboardingSequence_HappyPath(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)

	input := ceremony.OnboardingInput{
		NamespaceName:  "test-happy-path",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "test-party",
		SynchronizerID: "global",
		Threshold:      1,
	}

	sr, err := operations.ExecuteSequence(b, ceremony.OnboardingSequence, newDeps("p1"), input)
	require.NoError(t, err)

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
	assert.Contains(t, sr.Output.PartyID, "test-party::")

	// InitMember×3 + CreateProposal×1 + SignProposal×3 (2 error out) + SubmitProposal×1 + OnboardingSequence×1 = 9
	// The sequence's own report is included alongside its child operation reports.
	assert.Len(t, sr.ExecutionReports, 9)
}

// TestOnboardingSequence_Idempotent verifies that calling the sequence twice
// with the same input returns the cached report without re-executing any
// operation (report ID must be identical on both calls).
func TestOnboardingSequence_Idempotent(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	deps := newDeps("p1")
	input := ceremony.OnboardingInput{
		NamespaceName:  "test-happy-path",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "test-party",
		SynchronizerID: "global",
		Threshold:      1,
	}

	// First call succeeds as threshold=1 is met.  Second call must return the same report from cache.
	sr1, err := operations.ExecuteSequence(b, ceremony.OnboardingSequence, deps, input)
	require.NoError(t, err)

	// Same input and same deps → should hit the cache and return the same report without re-running any operations.
	sr2, err := operations.ExecuteSequence(b, ceremony.OnboardingSequence, deps, input)
	require.NoError(t, err)

	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report, not a new execution")
	assert.Equal(t, sr1.Output, sr2.Output)
}

// TestOnboardingSequence_AsyncMultiActorFlow simulates the "resume" pattern:
// multiple operators run the tool at different times and share a single
// persistent Reporter (the stand-in for the git-backed ceremony directory).
//
// Run 1 — coordinator p1 executes everything (threshold=1 so it completes).
// Run 2 — a second actor reuses the shared reporter and gets the cached result
// instantly without re-running any side-effects.
func TestOnboardingSequence_AsyncMultiActorFlow(t *testing.T) {
	t.Parallel()

	// sharedReporter is the durable state store shared across tool invocations.
	sharedReporter := operations.NewMemoryReporter()

	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// threshold=1 so the first actor completes the ceremony on its own.
	input := ceremony.OnboardingInput{
		NamespaceName:  "async-test-2026",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "async-party",
		SynchronizerID: "global",
		Threshold:      1,
	}

	// Run 1: coordinator p1 — all operations execute, ceremony completes.
	sr1, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p1"), input)
	require.NoError(t, err)
	assert.NotEmpty(t, sr1.Output.PartyID)

	// Run 2: any other actor — whole sequence is cached, no side-effects repeat.
	sr2, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p2"), input)
	require.NoError(t, err)
	assert.Equal(t, sr1.ID, sr2.ID, "second invocation must return the cached sequence report")
}

// TestOnboardingSequence_ThresholdNotMet asserts that [ceremony.ErrThresholdNotMet]
// is returned when fewer signers than required have run the sign step.
func TestOnboardingSequence_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	// Only p1 is allowed to sign.
	partialDeps := newDeps("p1")

	input := ceremony.OnboardingInput{
		NamespaceName:  "threshold-test-2026",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "threshold-party",
		SynchronizerID: "global",
		Threshold:      2, // requires 2 signatures; only 1 will succeed
	}

	b := optest.NewBundle(t)
	_, err := operations.ExecuteSequence(b, ceremony.OnboardingSequence, partialDeps, input)

	// ExecuteSequence stores the error message in a ReportError, which breaks the
	// error chain.  Use ErrorContains to check the sentinel message is preserved.
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error())
}

// TestOnboardingSequence_ResumeAfterPartialSigning simulates the realistic
// async scenario where signers come online in separate runs:
//
//  1. First run: only p1 has signed → ErrThresholdNotMet.
//  2. Second run: p2 also signs (threshold=2 met) → ceremony completes.
//
// Both runs share the same reporter so cached work is reused.
func TestOnboardingSequence_ResumeAfterPartialSigning(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()

	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	input := ceremony.OnboardingInput{
		NamespaceName:  "resume-test-2026",
		Participants:   []string{"p1", "p2", "p3"},
		PartyName:      "resume-party",
		SynchronizerID: "global",
		Threshold:      3, // require all 3 signatures for this test
	}

	// Run 1: only p1 can sign — threshold not met.
	onlyP1 := newDeps("p1")
	_, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, onlyP1, input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error(), "run 1 must fail: only 1 of 2 required signatures present")

	// Run 2: only p2 can sign — threshold not met.
	onlyP2 := newDeps("p2")
	_, err = operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, onlyP2, input)
	require.ErrorContains(t, err, ceremony.ErrThresholdNotMet.Error(), "run 2 must fail: only 1 of 2 required signatures present")

	// // Run 3: only p3 can sign — threshold met.
	sr, err := operations.ExecuteSequence(newBundle(), ceremony.OnboardingSequence, newDeps("p3"), input)
	require.NoError(t, err, "run 3 must succeed: all required signatures present")

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
}
