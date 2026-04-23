package kick

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// KickSequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input KickInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full KickOutput from the most recent run of
// KickSequence for the given input, or false if no run has been recorded.
func QueryState(reporter operations.Reporter, input KickInput) (KickOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return KickOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != KickSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[KickInput](r.Input)
		if !ok || inp.DecentralizedPartyID != input.DecentralizedPartyID {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return KickOutput{}, false
	}

	out, ok := ceremony.Rehydrate[KickOutput](latest.Output)
	if !ok {
		return KickOutput{}, false
	}

	return out, true
}
