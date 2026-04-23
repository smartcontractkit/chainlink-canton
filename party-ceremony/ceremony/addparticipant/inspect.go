package addparticipant

import (
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// AddParticipantSequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input AddParticipantInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full AddParticipantOutput from the most recent run of
// AddParticipantSequence for the given input, or false if no run has been recorded.
func QueryState(reporter operations.Reporter, input AddParticipantInput) (AddParticipantOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return AddParticipantOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != AddParticipantSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[AddParticipantInput](r.Input)
		if !ok || inp.DecentralizedPartyID != input.DecentralizedPartyID || inp.NewParticipantID != input.NewParticipantID {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return AddParticipantOutput{}, false
	}

	out, ok := ceremony.Rehydrate[AddParticipantOutput](latest.Output)
	if !ok {
		return AddParticipantOutput{}, false
	}

	return out, true
}
