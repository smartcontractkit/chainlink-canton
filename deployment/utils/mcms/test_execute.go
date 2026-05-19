package mcms

import (
	"encoding/json"
	"fmt"

	mcms_types "github.com/smartcontractkit/mcms/types"

	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// ExecuteOutputsForTest builds a batch operation from exercise outputs for test execution.
// This is a helper to collect MCMS transactions from multiple operation calls.
//
// Example usage in tests:
//
//	var outputs []contract.ExerciseOutput
//	out1, _ := operations.ExecuteOperation(b, global_config.ApplyDestChainConfigUpdates, chain, input)
//	outputs = append(outputs, out1.Output)
//	out2, _ := operations.ExecuteOperation(b, fee_quoter.UpdatePrices, chain, input)
//	outputs = append(outputs, out2.Output)
//
//	batchOp, _ := mcms.ExecuteOutputsForTest(outputs)
//	// Then execute batchOp via Bypasser pattern from integration-tests/mcms/mcms_timelock_test.go
func ExecuteOutputsForTest(
	outputs []opcontract.ExerciseOutput,
) (mcms_types.BatchOperation, error) {
	return opcontract.NewBatchOperationFromExercises(outputs)
}

// ExtractTransactions returns the MCMS transactions from exercise outputs.
// Useful for inspection or custom test execution.
func ExtractTransactions(outputs []opcontract.ExerciseOutput) []mcms_types.Transaction {
	var txs []mcms_types.Transaction
	for _, out := range outputs {
		if !out.Executed() {
			txs = append(txs, out.Tx)
		}
	}
	return txs
}

// DecodeCantonAdditionalFields decodes the MCMS transaction additional fields for Canton.
// This extracts the target contract address, choice name, and operation data.
func DecodeCantonAdditionalFields(raw []byte) (CantonAdditionalFields, error) {
	var fields CantonAdditionalFields
	if err := json.Unmarshal(raw, &fields); err != nil {
		return CantonAdditionalFields{}, fmt.Errorf("failed to unmarshal additional fields: %w", err)
	}
	return fields, nil
}

// CantonAdditionalFields holds the decoded MCMS transaction additional fields.
type CantonAdditionalFields struct {
	TargetInstanceAddress string `json:"targetInstanceAddress"`
	FunctionName          string `json:"functionName"`
	OperationData         string `json:"operationData"`
	TargetTemplateID      string `json:"targetTemplateID"`
}

// TestExecutePattern documents the pattern for executing MCMS proposals in tests.
// This is a reference implementation - actual tests should use the helper functions
// from integration-tests/mcms/mcms_timelock_test.go.
//
// For Bypasser execution (immediate, no timelock):
//
//   1. Generate operations with DisableMCMS=false (default)
//   2. Collect outputs: outputs = append(outputs, out.Output)
//   3. Build batch: batchOp, _ := ExecuteOutputsForTest(outputs)
//   4. Get MCMS contract CID and instance address
//   5. Get op count: queryBypasserOpCountForParty(t, participant, party, mcmsCid)
//   6. Build proposal: NewMCMSProposal(chainID, multisigID, opCount, false).AddOperation(...).Build()
//   7. Sign: proposal.Sign(validUntil, signers[:threshold])
//   8. SetRoot: setRootWithRoleAndDisclosureParty(t, participant, owner, party, mcmsCid, "Bypasser", proposal, validUntil, signatures)
//   9. Execute: bypasserExecuteBatchWithDisclosureParty(t, participant, owner, party, mcmsCid, targetCids, proposal.Operations[0], opProof)
//
// For Proposer execution (with timelock delay):
//
//   1-7. Same as above
//   8. SetRoot: setRootWithRoleAndDisclosureParty(..., "Proposer", ...)
//   9. Schedule: scheduleExecuteOp(t, participant, owner, party, mcmsCid, targetCids, proposal.Operations[0], opProof)
//   10. Wait for timelock delay
//   11. ExecuteScheduled: executeScheduledBatchWithDisclosure(t, participant, owner, party, mcmsCid, targetCids, proposal.Operations[0], opProof, validUntil)
const TestExecutePattern = "See integration-tests/mcms/mcms_timelock_test.go for helper functions"

// TODO: AutoExecuteProposal helper for local tests
//
// For production, MCMS proposals are executed manually via the MCMS UI or CLI.
// For local integration tests, we need an auto-execution helper that:
//
//   1. Takes MCMSTimelockProposals from changeset output
//   2. Signs with test signers
//   3. Executes via Bypasser role (immediate, no timelock)
//   4. Returns execution result
//
// This helper should be added to integration-tests/mcms package and would look like:
//
// func AutoExecuteTimelockProposals(
//     t testing.TB,
//     participant canton.Participant,
//     proposals []mcms.TimelockProposal,
//     mcmsContractInfo MCMSContractInfo,
//     signers []MCMSSigner,
// )
//
// The TestCCIPSend and other integration tests can then use this helper to
// automatically execute any generated MCMS proposals.
