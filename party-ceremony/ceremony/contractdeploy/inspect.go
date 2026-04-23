package contractdeploy

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

// LatestSequenceState returns the CeremonyState from the most recent run of
// ContractDeploySequence for the given input, or false if no run has been recorded.
func LatestSequenceState(reporter operations.Reporter, input ContractDeployInput) (CeremonyState, bool) {
	out, ok := QueryState(reporter, input)
	if !ok {
		return CeremonyState{}, false
	}

	return out.State, true
}

// QueryState returns the full ContractDeployOutput from the most recent run of
// ContractDeploySequence for the given input, or false if no run has been recorded.
func QueryState(reporter operations.Reporter, input ContractDeployInput) (ContractDeployOutput, bool) {
	reports, err := reporter.GetReports()
	if err != nil {
		return ContractDeployOutput{}, false
	}

	var latest *operations.Report[any, any]
	for i := range reports {
		r := &reports[i]
		if r.Def.ID != ContractDeploySequence.ID() {
			continue
		}
		inp, ok := ceremony.Rehydrate[ContractDeployInput](r.Input)
		if !ok || inp.DecentralizedPartyID != input.DecentralizedPartyID ||
			inp.TemplateModule != input.TemplateModule || inp.TemplateEntity != input.TemplateEntity {
			continue
		}
		if latest == nil || r.Timestamp.After(*latest.Timestamp) {
			latest = r
		}
	}
	if latest == nil {
		return ContractDeployOutput{}, false
	}

	out, ok := ceremony.Rehydrate[ContractDeployOutput](latest.Output)
	if !ok {
		return ContractDeployOutput{}, false
	}

	return out, true
}
