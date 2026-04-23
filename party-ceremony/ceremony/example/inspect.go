package example

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// OnboardingSequence for the given input, or false if no run has been recorded.
//
// This is the observer-side equivalent of a Temporal query: the handler embeds
// the live state into its output on every path (success or error), the
// framework persists it in the sequence report, and this function retrieves it.
// The direct caller of ExecuteSequence can simply read sr.Output.State instead.
func LatestSequenceState(reporter operations.Reporter, input OnboardingInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full OnboardingOutput from the most recent run of
// OnboardingSequence for the given input, or false if no run has been recorded.
// Unlike LatestSequenceState it also provides access to the final result fields
// (PartyID, DNSConfirmed, P2PConfirmed) populated when the ceremony completes.
func QueryState(reporter operations.Reporter, input OnboardingInput) (OnboardingOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return OnboardingOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != OnboardingSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[OnboardingInput](r.Input)
		if !ok || inp.NamespaceName != input.NamespaceName || inp.PartyName != input.PartyName {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return OnboardingOutput{}, false
	}

	out, ok := ceremony.Rehydrate[OnboardingOutput](latest.Output)
	if !ok {
		return OnboardingOutput{}, false
	}

	return out, true
}
