package exercisechoice

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// QueryState returns the output of the most recent run of
// ExerciseChoiceSequence for the given input, or false if none was recorded.
func QueryState(reporter operations.Reporter, input ExerciseChoiceInput) (ExerciseChoiceOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return ExerciseChoiceOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != ExerciseChoiceSequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[ExerciseChoiceInput](r.Input)
		if !ok || inp != input {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return ExerciseChoiceOutput{}, false
	}

	return ceremony.Rehydrate[ExerciseChoiceOutput](latest.Output)
}
