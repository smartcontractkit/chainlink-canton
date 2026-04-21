package onboarding

import (
	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// OnboardingSequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input OnboardingInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full OnboardingOutput from the most recent run of
// OnboardingSequence for the given input, or false if no run has been recorded.
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
		if !ok || inp.NamespaceName != input.NamespaceName || inp.PartyPrefix != input.PartyPrefix {
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
