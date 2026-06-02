package addparticipantwithacs

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// AddParticipantWithAcsSequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input AddParticipantWithAcsInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full AddParticipantWithAcsOutput from the most recent run of
// AddParticipantWithAcsSequence for the given input, or false if no run has been recorded.
func QueryState(reporter operations.Reporter, input AddParticipantWithAcsInput) (AddParticipantWithAcsOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return AddParticipantWithAcsOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != AddParticipantWithAcsSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[AddParticipantWithAcsInput](r.Input)
		if !ok || inp.DecentralizedPartyID != input.DecentralizedPartyID || inp.NewParticipantID != input.NewParticipantID {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return AddParticipantWithAcsOutput{}, false
	}

	out, ok := ceremony.Rehydrate[AddParticipantWithAcsOutput](latest.Output)
	if !ok {
		return AddParticipantWithAcsOutput{}, false
	}

	return out, true
}
