package keyrotation

import (
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// KeyRotationSequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input KeyRotationInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full KeyRotationOutput from the most recent run of
// KeyRotationSequence for the given input, or false if no run has been recorded.
func QueryState(reporter operations.Reporter, input KeyRotationInput) (KeyRotationOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return KeyRotationOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != KeyRotationSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[KeyRotationInput](r.Input)
		if !ok ||
			inp.DecentralizedPartyID != input.DecentralizedPartyID ||
			inp.TargetParticipantID != input.TargetParticipantID {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return KeyRotationOutput{}, false
	}

	out, ok := ceremony.Rehydrate[KeyRotationOutput](latest.Output)
	if !ok {
		return KeyRotationOutput{}, false
	}

	return out, true
}
