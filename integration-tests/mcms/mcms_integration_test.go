package tests

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

func TestMCMS_Execute(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))

	participant := env.Participant(1)
	randomUserParticipant := env.Participant(2)

	// ========================
	// |   Setup: Upload DAR  |
	// ========================

	mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
	require.NoError(t, err)

	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant, randomUserParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded MCMS DAR, package IDs: %v", packageIDs)
	require.NotEmpty(t, packageIDs)
	mcmsPkgID := packageIDs[0]

	// ========================
	// |   Setup: Parties     |
	// ========================

	ccipOwner := participant.Party
	randomUser := randomUserParticipant.Party
	t.Logf("Using CCIP party: %s", ccipOwner)
	t.Logf("Using user party: %s", randomUser)

	// ========================
	// |   Setup: Signers     |
	// ========================

	t.Log("Creating signers...")

	// Create 3 signers for 2-of-3 config
	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	config := New2of3Config(signers)

	t.Logf("Signer 1: %s", signer1.Address)
	t.Logf("Signer 2: %s", signer2.Address)
	t.Logf("Signer 3: %s", signer3.Address)

	// Sort signers by address for consistent ordering
	sortedSigners := SortSignersByAddress(signers)

	// ========================
	// |   Contract Constants |
	// ========================

	chainId := int64(1)

	// Run tests
	t.Run("ExecuteOp Flow", func(t *testing.T) {
		t.Parallel()
		testExecuteOpFlow(t, mcmsPkgID, config, chainId, sortedSigners, participant, ccipOwner)
	})
	t.Run("Signature Verification Failure", func(t *testing.T) {
		t.Parallel()
		testSignatureVerificationFails(t, mcmsPkgID, config, chainId, participant, ccipOwner)
	})
	t.Run("Replay Protection", func(t *testing.T) {
		t.Parallel()
		testReplayProtection(t, mcmsPkgID, config, chainId, sortedSigners, participant, ccipOwner)
	})
	t.Run("MCMS Op", func(t *testing.T) {
		t.Parallel()
		testExecuteMCMSOp(t, mcmsPkgID, config, chainId, sortedSigners, participant, randomUserParticipant, ccipOwner, randomUser)
	})
	t.Run("Signatory Check", func(t *testing.T) {
		t.Parallel()
		testSignatoryCheck(t, mcmsPkgID, config, chainId, sortedSigners, participant, randomUserParticipant, ccipOwner, randomUser)
	})
}

// testExecuteOpFlow tests the complete MCMS execute flow with direct invocation:
// 3. Create proposal with "increment" operation
// 4. Sign with 2 signers
// 5. SetRoot with real signatures
// 6. ExecuteOp - direct call to Counter via MCMSReceiver interface
// 7. Verify counter value incremented
func testExecuteOpFlow(
	t *testing.T,
	mcmsPkgID string,
	config MCMSConfig,
	chainId int64,
	sortedSigners []*MCMSSigner,
	participant testhelpers.Participant,
	ccipOwnerParty string,
) {
	// ========================
	// |   Contract Constants |
	// ========================

	baseMcmsId := "mcms-integration-test-" + uuid.New().String()[:8]
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)
	multisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleBypasser)                                     // Use Bypasser for external calls
	counterInstanceId := fmt.Sprintf("counter-env-%s@%s", uuid.New().String()[:8], ccipOwnerParty) // Stable instance ID for Counter
	t.Logf("MCMS ID: %s", multisigId)
	t.Logf("Counter Instance ID: %s", counterInstanceId)

	// ========================
	// |   1. Create MCMS |
	// ========================

	t.Log("Creating MCMS contract...")

	// Build signer info values
	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	// Build group quorums (32 ints)
	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
	}

	// Build group parents (32 ints)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	// Create empty seen hashes map
	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}

	// Create epoch time for empty expiring root
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}

	// Empty expiring root
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}

	// Empty root metadata
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	// Multisig config record
	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}

	// RoleState record (config + replay + root state)
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}

	// Timelock fields (default empty)
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyTextList := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	// Create MCMS contract (simplified - no ticket tracking)
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "timelockTimestamps", Value: emptyMap},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyTextList},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// ========================
	// |   2. Create Counter  |
	// ========================

	t.Log("Creating Counter contract with MCMSReceiver interface...")

	counterCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Counter",
								EntityName: "Counter",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: counterInstanceId}}},
								{Label: "value", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	counterCid := counterCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created Counter contract: %s", counterCid)

	// ========================
	// |   3. Build Proposal  |
	// ========================

	t.Log("Building proposal for bypasser_execute_batch...")

	// Build bypasser_execute_batch proposal
	// The proposal targets MCMS itself with "bypasser_execute_batch" function
	// The actual Counter.increment call is encoded in operationData
	bypasserParams := BypasserExecuteParams{
		Calls: []TimelockCall{
			{
				TargetInstanceId: counterInstanceId,
				FunctionName:     "increment",
				OperationData:    "",
			},
		},
	}
	encodedBypasserParams := EncodeBypasserExecuteParams(bypasserParams)

	proposal := NewMCMSProposal(int(chainId), multisigId, 0, false)
	proposal.AddOperation(mcmsInstanceId, "bypasser_execute_batch", encodedBypasserParams)
	proposal.Build()

	root := proposal.GetRoot()
	t.Logf("Proposal root: %s", root)

	// Get proofs
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	t.Logf("Metadata proof: %v", metadataProof)

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)
	t.Logf("Op proof: %v", opProof)

	// ========================
	// |   4. Sign Proposal   |
	// ========================

	t.Log("Signing proposal with 2 signers...")

	// Valid for 1 hour
	validUntil := time.Now().Add(time.Hour)
	validUntilMicros := validUntil.UnixMicro()

	// Sign with first 2 sorted signers (to meet 2-of-3 quorum)
	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)
	require.Len(t, signaturesRaw, 2)

	t.Logf("Signature 1 from %s: r=%s..., s=%s...", sortedSigners[0].Address, signaturesRaw[0].R[:16], signaturesRaw[0].S[:16])
	t.Logf("Signature 2 from %s: r=%s..., s=%s...", sortedSigners[1].Address, signaturesRaw[1].R[:16], signaturesRaw[1].S[:16])

	// Build signature values for Canton
	signatureValues := make([]*apiv2.Value, len(signaturesRaw))
	for i, sig := range signaturesRaw {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	// Build metadata proof values
	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	// Build metadata value
	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	// ========================
	// |   5. SetRoot         |
	// ========================

	t.Log("Calling SetRoot for Bypasser role...")

	setRootRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Bypasser"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)

	// Get new MCMS contract ID from exercise result
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			mcmsCid = created.GetContractId()
			break
		}
	}
	t.Logf("SetRoot succeeded, new MCMS CID: %s", mcmsCid)

	// ========================
	// |   6. ExecuteOp       |
	// ========================

	t.Log("Calling ExecuteOp (direct invocation via MCMSReceiver interface)...")

	// Build op value
	op := proposal.Operations[0]
	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}

	// Build op proof values
	opProofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		opProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	executeOpRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "ExecuteOp",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Bypasser"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "op", Value: opValue},
									{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: opProofValues}}}},
									{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: counterCid}},
									}}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)

	// Get counter value, contract ID, and MCMSEntrypointEvent from the created events
	var counterValue int64 = -1
	var newCounterCid string
	var eventFound bool
	var eventInstanceId string
	var eventFunctionName string
	var eventOperationData string
	var eventContractIdsAsText []string

	for _, event := range executeOpRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil {
			entityName := created.GetTemplateId().GetEntityName()

			// Check for Counter contract
			if entityName == "Counter" {
				newCounterCid = created.GetContractId()
				counter, err := bindings.UnmarshalCreatedEvent[mcms.Counter](created)
				if err != nil {
					t.Logf("Failed to unmarshal Counter: %v", err)
				} else {
					counterValue = int64(counter.Value)
				}
			}

			// Check for MCMSEntrypointEvent (emitted by Counter.mcmsEntrypoint)
			if entityName == "MCMSEntrypointEvent" {
				eventFound = true
				event, err := bindings.UnmarshalCreatedEvent[mcms.MCMSEntrypointEvent](created)
				if err != nil {
					t.Logf("Failed to unmarshal MCMSEntrypointEvent: %v", err)
				} else {
					eventInstanceId = string(event.InstanceId)
					eventFunctionName = string(event.FunctionName)
					eventOperationData = string(event.OperationData)
					for _, cid := range event.ContractIdsAsText {
						eventContractIdsAsText = append(eventContractIdsAsText, string(cid))
					}
				}
			}
		}
	}
	t.Logf("ExecuteOp succeeded, counter value from event: %d", counterValue)
	require.NotEmpty(t, newCounterCid, "Should have new counter contract ID")

	// Verify MCMSEntrypointEvent was emitted with correct data
	require.True(t, eventFound, "MCMSEntrypointEvent should be emitted by Counter.mcmsEntrypoint")
	assert.Equal(t, counterInstanceId, eventInstanceId, "Event instanceId should match Counter instanceId")
	assert.Equal(t, "increment", eventFunctionName, "Event functionName should be 'increment'")
	assert.Empty(t, eventOperationData, "Event operationData should be empty for increment")

	// Note: contractIds is now empty in the current contract implementation
	t.Logf("MCMSEntrypointEvent verified: instanceId=%s, functionName=%s, contractIdsAsText=%v",
		eventInstanceId, eventFunctionName, eventContractIdsAsText)

	// ========================
	// |   7. Verify Counter  |
	// ========================

	// Also query the ACS to verify the counter value
	t.Log("Querying counter via ACS...")
	counterContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{
		PackageId:  mcmsPkgID,
		ModuleName: "MCMS.Counter",
		EntityName: "Counter",
	})
	require.NoError(t, err)

	var queriedValue int64 = -1
	for _, contract := range counterContracts {
		if contract.GetCreatedEvent().GetContractId() == newCounterCid {
			counter, err := bindings.UnmarshalCreatedEvent[mcms.Counter](contract.GetCreatedEvent())
			if err != nil {
				t.Logf("Failed to unmarshal Counter from ACS: %v", err)
			} else {
				queriedValue = int64(counter.Value)
			}
		}
	}
	t.Logf("Counter value from ACS query: %d", queriedValue)

	// Verify both values match
	require.Equal(t, int64(1), counterValue, "Counter from event should be 1")
	require.Equal(t, int64(1), queriedValue, "Counter from ACS query should be 1")
	require.Equal(t, counterValue, queriedValue, "Event value and ACS query value should match")

	t.Log("✓ Full MCMS ExecuteOp flow test completed successfully!")
	t.Log("Summary:")
	t.Log("  1. Created MCMS with 2-of-3 config")
	t.Log("  2. Created Counter contract (implements MCMSReceiver)")
	t.Log("  3. Built proposal with 'increment' operation targeting instanceId")
	t.Log("  4. Signed with 2 signers (real ECDSA signatures)")
	t.Log("  5. SetRoot with on-chain verification")
	t.Log("  6. ExecuteOp - MCMS directly calls Counter.MCMSReceiver_Entrypoint")
	t.Log("  7. Counter value = 1 ✓")
}

// testSignatureVerificationFails tests that invalid signatures are rejected
func testSignatureVerificationFails(
	t *testing.T,
	mcmsPkgID string,
	config MCMSConfig,
	chainId int64,
	participant testhelpers.Participant,
	ccipOwnerParty string,
) {
	// Build config values
	baseMcmsId := "mcms-sig-fail-test-" + uuid.New().String()[:8]
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)
	multisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleProposer)

	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyTextList := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	// Create MCMS contract
	t.Log("Creating MCMS contract...")
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "timelockTimestamps", Value: emptyMap},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyTextList},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// Build a valid proposal
	proposal := NewMCMSProposal(int(chainId), multisigId, 0, false)
	proposal.AddOperation("counter", "increment", "")
	proposal.Build()

	root := proposal.GetRoot()
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)

	validUntil := time.Now().Add(time.Hour)
	validUntilMicros := validUntil.UnixMicro()

	// Create INVALID signatures (using random data instead of actual signing)
	t.Log("Creating invalid signatures...")
	invalidSignatures := []*apiv2.Value{
		{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "04" + strings.Repeat("ab", 64)}}}, // fake pub key
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: strings.Repeat("12", 32)}}},                // fake r
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: strings.Repeat("34", 32)}}},                // fake s
			},
		}}},
		{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "04" + strings.Repeat("cd", 64)}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: strings.Repeat("56", 32)}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: strings.Repeat("78", 32)}}},
			},
		}}},
	}

	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	// Attempt SetRoot with invalid signatures - should fail
	t.Log("Attempting SetRoot with invalid signatures...")
	_, err = participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: invalidSignatures}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})

	require.Error(t, err)
	// Either E_NO_VALID_SIGNATURES (application-level) or CRYPTO_ERROR (Canton native crypto) is acceptable
	isExpectedError := strings.Contains(err.Error(), "E_NO_VALID_SIGNATURES") ||
		strings.Contains(err.Error(), "CRYPTO_ERROR") ||
		strings.Contains(err.Error(), "can not parse")
	require.True(t, isExpectedError,
		"Expected signature verification error, got: %v", err)

	t.Log("✓ SetRoot correctly rejected invalid signatures")
}

// testReplayProtection tests that the same root cannot be set twice
func testReplayProtection(
	t *testing.T,
	mcmsPkgID string,
	config MCMSConfig,
	chainId int64,
	sortedSigners []*MCMSSigner,
	participant testhelpers.Participant,
	ccipOwnerParty string,
) {
	// Build config values
	baseMcmsId := "mcms-replay-test-" + uuid.New().String()[:8]
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)
	multisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleProposer)

	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyTextList := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	// Create MCMS contract
	t.Log("Creating MCMS contract...")
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "timelockTimestamps", Value: emptyMap},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyTextList},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// Build proposal
	proposal := NewMCMSProposal(int(chainId), multisigId, 0, false)
	proposal.AddOperation("counter", "increment", "")
	proposal.Build()

	root := proposal.GetRoot()
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)

	validUntil := time.Now().Add(time.Hour)
	validUntilMicros := validUntil.UnixMicro()

	// Sign with 2 signers
	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	signatureValues := make([]*apiv2.Value, len(signaturesRaw))
	for i, sig := range signaturesRaw {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	// First SetRoot - should succeed
	t.Log("First SetRoot call (should succeed)...")
	setRootRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	t.Log("First SetRoot succeeded")

	// Get new MCMS contract ID
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			mcmsCid = created.GetContractId()
			break
		}
	}

	// Second SetRoot with SAME signatures - should fail with E_ALREADY_SEEN_HASH
	t.Log("Second SetRoot call with same signatures (should fail)...")
	_, err = participant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "E_ALREADY_SEEN_HASH",
		"Expected E_ALREADY_SEEN_HASH error, got: %v", err)

	t.Log("✓ Second SetRoot correctly rejected with E_ALREADY_SEEN_HASH (replay protection working)")
}

// testExecuteMCMSOp tests self-dispatch MCMS operations (Aptos pattern)
// This demonstrates changing MCMS config via a signed proposal
func testExecuteMCMSOp(
	t *testing.T,
	mcmsPkgID string,
	config MCMSConfig,
	chainId int64,
	sortedSigners []*MCMSSigner,
	ccipParticipant, userParticipant testhelpers.Participant,
	ccipOwnerParty, userParty string,
) {
	// ========================
	// |   Contract Constants |
	// ========================

	baseMcmsId := "mcms-op-test-" + uuid.New().String()[:8]
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)
	multisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleProposer)
	t.Logf("MCMS ID: %s", multisigId)

	// ========================
	// |   1. Create MCMS |
	// ========================

	t.Log("Creating MCMS contract...")

	// Build signer info values
	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyTextList := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	mcmsCreateRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: mcmsInstanceId}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "timelockTimestamps", Value: emptyMap},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyTextList},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// ========================
	// |   2. Build Proposal  |
	// ========================

	t.Log("Building MCMS proposal (set_config targeting 'self')...")

	// Prepare new config params (change from 2-of-3 to 1-of-3)
	newQuorums := make([]int, NumGroups)
	newQuorums[0] = 1 // Changed from 2 to 1
	newParents := make([]int, NumGroups)

	// Encode SetConfigParams into operationData (like Aptos BCS)
	setConfigParams := SetConfigParams{
		Signers:      config.Signers,
		GroupQuorums: newQuorums,
		GroupParents: newParents,
		ClearRoot:    false,
	}
	encodedParams := EncodeSetConfigParams(setConfigParams)
	t.Logf("Encoded set_config params: %s... (%d bytes)", encodedParams[:min(40, len(encodedParams))], len(encodedParams)/2)

	// Build proposal with MCMS operation targeting this MCMS instanceId
	// operationData contains the encoded params (like Aptos)
	proposal := NewMCMSProposal(int(chainId), multisigId, 0, false)
	proposal.AddOperation(mcmsInstanceId, "set_config", encodedParams) // Params encoded in operationData
	proposal.Build()

	root := proposal.GetRoot()
	t.Logf("MCMS proposal root: %s", root)

	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)

	// ========================
	// |   3. Sign Proposal   |
	// ========================

	t.Log("Signing MCMS proposal with 2 signers...")

	validUntil := time.Now().Add(time.Hour)
	validUntilMicros := validUntil.UnixMicro()

	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	signatureValues := make([]*apiv2.Value, len(signaturesRaw))
	for i, sig := range signaturesRaw {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	// ========================
	// |   4. SetRoot         |
	// ========================

	t.Log("Calling SetRoot with MCMS proposal...")

	setRootRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)

	// Get new contract ID from SetRoot result
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			mcmsCid = created.GetContractId()
			break
		}
	}
	t.Logf("SetRoot succeeded, new MCMS CID: %s", mcmsCid)

	// Query ACS to get disclosed contract with CreatedEventBlob (for bob to use)
	// Transaction events don't include the blob by default, so we query the ACS
	disclosedMcms, err := testhelpers.GetDisclosedContractById(t.Context(), ccipParticipant, mcmsCid)
	require.NoError(t, err, "Should get disclosed contract from ACS")
	t.Logf("Got disclosed contract with blob size: %d bytes", len(disclosedMcms.CreatedEventBlob))

	// ========================
	// |   5. ExecuteMcmsOp  |
	// ========================
	//
	// ANYONE CAN EXECUTE OPS IN MCMS
	// ==============================
	// Security does NOT come from who submits the transaction.
	// Security comes from:
	//   1. SetRoot validated the signatures from required signers (2-of-3)
	//   2. SetRoot stored the merkle root on-chain
	//   3. ExecuteMcmsOp validates the merkle proof against the stored root
	//   4. The operation data is cryptographically bound to the merkle leaf
	//
	// This means bob (randomUser) can execute the operation even though:
	//   - bob is NOT a signer
	//   - bob did NOT create the contract
	//   - bob only has visibility via contract disclosure
	//
	// The merkle proof verification ensures only pre-approved operations can execute.

	t.Log("Calling ExecuteMcmsOp as randomUser (bob) with disclosed contract...")

	// Build op value
	op := proposal.Operations[0]
	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}

	opProofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		opProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	// No separate params - params are encoded in op.operationData (like Aptos BCS)
	// randomUser (bob) submits via userParticipant with disclosed contract
	_, err = userParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "ExecuteOp",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: userParty}}},
									{Label: "op", Value: opValue},
									{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: opProofValues}}}},
									{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs:              []string{userParty},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedMcms}, // Grant bob visibility
		},
	})
	require.NoError(t, err)
	t.Log("ExecuteMcmsOp succeeded (bob executed with disclosed contract)")
	time.Sleep(5 * time.Second) // Wait for ACS to update on ccip participant

	// Store old contract ID to verify it changed
	oldMcmsCid := mcmsCid

	// Bob can't see the created event (not an observer), so query from alice's participant
	// The old contract should be archived, so find the active one
	t.Log("Querying ACS from alice's participant to find the new contract...")

	// ========================
	// |   6. Verify Config   |
	// ========================

	// Query via GetActiveContracts with verbose to get the actual config
	// Find the MCMS contract with our instanceId (the old one is archived, new one is active)
	var newNumSigners int64 = -1
	var newQuorum int64 = -1
	var newMcmsCid string
	offset, _ := testhelpers.GetCurrentOffset(t.Context(), ccipParticipant)
	acsRes, err := ccipParticipant.StateServiceClient.GetActiveContracts(t.Context(), &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				ccipOwnerParty: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
								TemplateFilter: &apiv2.TemplateFilter{
									TemplateId: &apiv2.Identifier{
										PackageId:  mcmsPkgID,
										ModuleName: "MCMS.Main",
										EntityName: "MCMS",
									},
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	require.NoError(t, err)
	defer acsRes.CloseSend()

	foundContract := false
	for {
		ac, err := acsRes.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		c, ok := ac.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}

		if c.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName() != bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			continue
		}

		// Parse fields manually to avoid RELTIME unmarshaling issues with bindings
		// The MCMS contract has a minDelay field of type RelTime which Canton returns
		// as a record {"microseconds": N} but go-daml expects a simple number
		fields := c.ActiveContract.GetCreatedEvent().GetCreateArguments().GetFields()
		var contractInstanceId string
		var proposerConfig *apiv2.Record
		for _, field := range fields {
			switch field.GetLabel() {
			case "instanceId":
				contractInstanceId = field.GetValue().GetText()
			case "proposer":
				// proposer is a RoleState record with a config field
				for _, pf := range field.GetValue().GetRecord().GetFields() {
					if pf.GetLabel() == "config" {
						proposerConfig = pf.GetValue().GetRecord()
					}
				}
			}
		}

		// Check if this is our MCMS by matching instanceId (new multi-role structure)
		if contractInstanceId != mcmsInstanceId {
			continue
		}

		newMcmsCid = c.ActiveContract.GetCreatedEvent().ContractId

		// Extract signers and groupQuorums from proposer config
		if proposerConfig != nil {
			for _, cf := range proposerConfig.GetFields() {
				switch cf.GetLabel() {
				case "signers":
					newNumSigners = int64(len(cf.GetValue().GetList().GetElements()))
				case "groupQuorums":
					quorums := cf.GetValue().GetList().GetElements()
					if len(quorums) > 0 {
						newQuorum = quorums[0].GetInt64()
					}
				}
			}
		}
		foundContract = true

		break
	}
	require.True(t, foundContract, "Should find MCMS contract with instanceId=%s in ACS", mcmsInstanceId)
	require.NotEmpty(t, newMcmsCid, "Should have new MCMS contract ID")
	require.NotEqual(t, oldMcmsCid, newMcmsCid, "MCMS contract ID should change after ExecuteMcmsOp (old contract archived, new created)")
	t.Logf("Found new MCMS contract: %s (changed from %s)", newMcmsCid, oldMcmsCid)
	t.Logf("Verified config from ACS: numSigners=%d, quorum=%d", newNumSigners, newQuorum)

	require.Equal(t, int64(3), newNumSigners, "Should still have 3 signers")
	require.Equal(t, int64(1), newQuorum, "Quorum should be changed to 1")

	t.Log("✓ ExecuteMcmsOp test completed successfully!")
	t.Log("Summary:")
	t.Log("  1. Created MCMS with 2-of-3 config")
	t.Log("  2. Built MCMS proposal with set_config targeting 'self'")
	t.Log("  3. Signed with 2 signers")
	t.Log("  4. SetRoot with on-chain verification")
	t.Log("  5. ExecuteMcmsOp - self-dispatch to change config")
	t.Log("  6. Config changed from 2-of-3 to 1-of-3 ✓")
}

// testSignatoryCheck verifies that MCMS can only execute operations
// on contracts where the MCMS owner is a signatory of the target.
//
// Note: Due to Canton's participant-party model, we cannot create a contract
// on participant1 owned by participant2's party. This test verifies that
// when both MCMS and Counter have the same owner, ExecuteOp succeeds.
//
// The signatory check rejection case is covered by Daml unit tests in
// contracts/mcms/test/daml/MCMS/FlowTest.daml:TestSignatoryCheckRejectsWrongOwner
// which creates MCMS and Counter with different owners in the same participant.
func testSignatoryCheck(
	t *testing.T,
	mcmsPkgID string,
	config MCMSConfig,
	chainId int64,
	sortedSigners []*MCMSSigner,
	ccipParticipant, _ testhelpers.Participant,
	ccipOwnerParty, _ string,
) {
	// ========================
	// |   Contract Constants |
	// ========================

	baseMcmsId := "mcms-signatory-test-" + uuid.New().String()[:8]
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)
	multisigId := MakeMcmsId(mcmsInstanceId, MCMSRoleBypasser) // Use Bypasser for external calls
	counterInstanceId := fmt.Sprintf("counter-signatory-%s@%s", uuid.New().String()[:8], ccipOwnerParty)
	t.Logf("MCMS ID: %s", multisigId)
	t.Logf("Counter Instance ID: %s", counterInstanceId)

	// ========================
	// |   1. Create MCMS (owned by ccipOwner) |
	// ========================

	t.Log("Creating MCMS contract owned by ccipOwner...")

	// Build signer info values
	signerInfoValues := make([]*apiv2.Value, len(config.Signers))
	for i, si := range config.Signers {
		signerInfoValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "signerAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: si.SignerAddress}}},
				{Label: "signerIndex", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerIndex)}}},
				{Label: "signerGroup", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(si.SignerGroup)}}},
			},
		}}}
	}

	groupQuorumValues := make([]*apiv2.Value, NumGroups)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := range NumGroups {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
		groupParentValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupParents[i])}}
	}

	emptyMap := &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{}}}}
	epochTime := &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: 0}}
	emptyExpiringRoot := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "root", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "validUntil", Value: epochTime},
			{Label: "opCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyRootMetadata := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
		},
	}}}

	configValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
			{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
			{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
		},
	}}}
	roleStateValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "config", Value: configValue},
			{Label: "seenHashes", Value: emptyMap},
			{Label: "expiringRoot", Value: emptyExpiringRoot},
			{Label: "rootMetadata", Value: emptyRootMetadata},
		},
	}}}
	minDelayValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "microseconds", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
		},
	}}}
	emptyTextList := &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}}

	mcmsCreateRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: fmt.Sprintf("%s@%s", baseMcmsId, ccipOwnerParty)}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "proposer", Value: roleStateValue},
								{Label: "canceller", Value: roleStateValue},
								{Label: "bypasser", Value: roleStateValue},
								{Label: "timelockTimestamps", Value: emptyMap},
								{Label: "minDelay", Value: minDelayValue},
								{Label: "blockedFunctions", Value: emptyTextList},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// ========================
	// |   2. Create Counter (also owned by ccipOwner, SAME as MCMS owner) |
	// ========================
	//
	// This tests the SUCCESS case of the signatory check:
	// - MCMS owner = ccipOwnerParty
	// - Counter owner = ccipOwnerParty (same!)
	// - Signatory check: ccipOwnerParty `elem` signatory Counter ✓ PASSES

	t.Log("Creating Counter contract also owned by ccipOwner (same as MCMS owner)...")

	counterCreateRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Counter",
								EntityName: "Counter",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: counterInstanceId}}},
								{Label: "value", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err)
	counterCid := counterCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created Counter contract: %s", counterCid)

	// ========================
	// |   3. Build Proposal  |
	// ========================

	t.Log("Building bypasser_execute_batch proposal targeting the counter...")

	// Build bypasser_execute_batch proposal with "increment" call
	bypasserParams := BypasserExecuteParams{
		Calls: []TimelockCall{
			{
				TargetInstanceId: counterInstanceId,
				FunctionName:     "increment",
				OperationData:    "",
			},
		},
	}
	encodedBypasserParams := EncodeBypasserExecuteParams(bypasserParams)

	proposal := NewMCMSProposal(int(chainId), multisigId, 0, false)
	proposal.AddOperation(mcmsInstanceId, "bypasser_execute_batch", encodedBypasserParams)
	proposal.Build()

	root := proposal.GetRoot()
	t.Logf("Proposal root: %s", root)

	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)

	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)

	// ========================
	// |   4. Sign Proposal   |
	// ========================

	t.Log("Signing proposal with 2 signers...")

	validUntil := time.Now().Add(time.Hour)
	validUntilMicros := validUntil.UnixMicro()

	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)
	require.Len(t, signaturesRaw, 2)

	signatureValues := make([]*apiv2.Value, len(signaturesRaw))
	for i, sig := range signaturesRaw {
		signatureValues[i] = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "publicKey", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.PublicKey}}},
				{Label: "r", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.R}}},
				{Label: "s", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: sig.S}}},
			},
		}}}
	}

	metadataProofValues := make([]*apiv2.Value, len(metadataProof))
	for i, p := range metadataProof {
		metadataProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	metadataValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: proposal.Metadata.MultisigId}}},
			{Label: "preOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PreOpCount)}}},
			{Label: "postOpCount", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(proposal.Metadata.PostOpCount)}}},
			{Label: "overridePreviousRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: proposal.Metadata.OverridePreviousRoot}}},
		},
	}}}

	// ========================
	// |   5. SetRoot         |
	// ========================

	t.Log("Calling SetRoot (should succeed - root doesn't check ownership)...")

	setRootRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Bypasser"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "newRoot", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: root}}},
									{Label: "validUntil", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: validUntilMicros}}},
									{Label: "metadata", Value: metadataValue},
									{Label: "metadataProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: metadataProofValues}}}},
									{Label: "signatures", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signatureValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err, "SetRoot should succeed for Bypasser role")

	// Get new MCMS contract ID
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.MCMS{}.GetTemplateID()) {
			mcmsCid = created.GetContractId()
			break
		}
	}
	t.Logf("SetRoot succeeded, new MCMS CID: %s", mcmsCid)

	// ========================
	// |   6. ExecuteOp (should SUCCEED - same owner) |
	// ========================
	//
	// Both MCMS and Counter are owned by ccipOwnerParty (same owner).
	// The signatory check: ccipOwnerParty `elem` signatory Counter ✓ PASSES
	// This verifies the signatory check logic is working correctly.

	t.Log("Calling ExecuteOp (should succeed - same owner, signatory check passes)...")

	// Build op value
	op := proposal.Operations[0]
	opValue := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
		Fields: []*apiv2.RecordField{
			{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.ChainId)}}},
			{Label: "multisigId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.MultisigId}}},
			{Label: "nonce", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(op.Nonce)}}},
			{Label: "targetInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.TargetInstanceId}}},
			{Label: "functionName", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.FunctionName}}},
			{Label: "operationData", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: op.OperationData}}},
		},
	}}}

	opProofValues := make([]*apiv2.Value, len(opProof))
	for i, p := range opProof {
		opProofValues[i] = &apiv2.Value{Sum: &apiv2.Value_Text{Text: p}}
	}

	// ExecuteOp - should SUCCEED because both MCMS and Counter have the same owner
	executeOpRes, err := ccipParticipant.CommandServiceClient.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  mcmsPkgID,
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "ExecuteOp",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "targetRole", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Bypasser"}}}},
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwnerParty}}},
									{Label: "op", Value: opValue},
									{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: opProofValues}}}},
									{Label: "targetCids", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: counterCid}},
									}}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwnerParty},
		},
	})
	require.NoError(t, err, "ExecuteOp should succeed when MCMS owner is a signatory of target")

	// Verify counter was incremented
	var counterValue int64 = -1
	for _, event := range executeOpRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == bindings.GetEntityName(mcms.Counter{}.GetTemplateID()) {
			for _, field := range created.GetCreateArguments().GetFields() {
				if field.GetLabel() == "value" {
					counterValue = field.GetValue().GetInt64()
					break
				}
			}
		}
	}
	require.Equal(t, int64(1), counterValue, "Counter should be incremented to 1")

	t.Log("✓ Signatory check test completed successfully!")
	t.Log("Summary:")
	t.Log("  1. Created MCMS owned by ccipOwner")
	t.Log("  2. Created Counter also owned by ccipOwner (same owner)")
	t.Log("  3. Built and signed proposal with valid signatures")
	t.Log("  4. SetRoot succeeded (root doesn't check ownership)")
	t.Log("  5. ExecuteOp succeeded ✓")
	t.Log("     Signatory check passed: ccipOwner is a signatory of Counter")
	t.Log("")
	t.Log("Note: The rejection case (different owners) is covered by Daml unit tests:")
	t.Log("  contracts/mcms/test/daml/MCMS/FlowTest.daml:TestSignatoryCheckRejectsWrongOwner")
}

// TestMCMS_GenerateDamlTestValues generates all cryptographic values needed for Daml unit tests.
// Run this test and copy the output to contracts/mcms/test/daml/MCMS/FlowTest.daml
// This uses FIXED values (not random) so the output is deterministic.
func TestMCMS_GenerateDamlTestValues(t *testing.T) {
	t.Parallel()

	t.Log("=======================================================================")
	t.Log("GENERATING DAML TEST VALUES")
	t.Log("Copy these values to contracts/mcms/test/daml/MCMS/FlowTest.daml")
	t.Log("=======================================================================")

	// Use fixed seeds for deterministic output across runs
	// Create 3 signers with deterministic keys
	signer1, err := NewMCMSSignerFromSeed(42001)
	require.NoError(t, err)
	signer2, err := NewMCMSSignerFromSeed(42002)
	require.NoError(t, err)
	signer3, err := NewMCMSSignerFromSeed(42003)
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	config := New2of3Config(signers)
	sortedSigners := SortSignersByAddress(signers)

	// Fixed test values
	chainId := 1
	baseMcmsId := "mcms-daml-test"
	mcmsInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, "ccip_owner-9cefe94d")
	mcmsId := MakeMcmsId(mcmsInstanceId, MCMSRoleProposer)

	t.Log("")
	t.Log("-- ===========================================================================")
	t.Log("-- SIGNER CONFIG (2-of-3)")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("testSigners : [SignerInfo]")
	t.Log("testSigners =")
	for i, si := range config.Signers {
		comma := ","
		if i == len(config.Signers)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s SignerInfo \"%s\" %d %d%s", bracket, si.SignerAddress, si.SignerIndex, si.SignerGroup, comma)
	}
	t.Log("  ]")
	t.Log("")

	// Build proposal with one operation
	// targetInstanceId must match Counter view (instanceId includes @partyId)
	// In Daml sandbox (daml test), allocateParty "ccip_owner" → partyToText = "ccip_owner-9cefe94d"
	// The suffix is deterministic in the sandbox based on the hint string.
	proposal := NewMCMSProposal(chainId, mcmsId, 0, false)
	proposal.AddOperation("counter@ccip_owner-9cefe94d", "increment", "")
	proposal.Build()

	root := proposal.GetRoot()
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)

	// Use a fixed validUntil far in the future (year 2030)
	validUntil := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	// Sign with first 2 sorted signers
	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- PROPOSAL VALUES")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Logf("testChainId : Int")
	t.Logf("testChainId = %d", chainId)
	t.Log("")
	t.Logf("testMcmsInstanceId : Text")
	t.Logf("testMcmsInstanceId = \"%s\"", mcmsInstanceId)
	t.Log("")
	t.Log("testOp : Op")
	t.Log("testOp = Op")
	t.Logf("  { chainId = %d", proposal.Operations[0].ChainId)
	t.Logf("  , multisigId = \"%s\"", proposal.Operations[0].MultisigId)
	t.Logf("  , nonce = %d", proposal.Operations[0].Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", proposal.Operations[0].TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", proposal.Operations[0].FunctionName)
	t.Logf("  , operationData = \"%s\"", proposal.Operations[0].OperationData)
	t.Log("  }")
	t.Log("")
	t.Log("testMetadata : RootMetadata")
	t.Log("testMetadata = RootMetadata")
	t.Logf("  { chainId = %d", proposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", proposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", proposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", proposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", proposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")
	t.Logf("testRoot : Text")
	t.Logf("testRoot = \"%s\"", root)
	t.Log("")

	t.Log("-- ===========================================================================")
	t.Log("-- MERKLE PROOFS")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("testMetadataProof : [Text]")
	if len(metadataProof) == 0 {
		t.Log("testMetadataProof = []")
	} else {
		t.Log("testMetadataProof =")
		for i, p := range metadataProof {
			comma := ","
			if i == len(metadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("testOpProof : [Text]")
	if len(opProof) == 0 {
		t.Log("testOpProof = []")
	} else {
		t.Log("testOpProof =")
		for i, p := range opProof {
			comma := ","
			if i == len(opProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")

	t.Log("-- ===========================================================================")
	t.Log("-- SIGNATURES (Real ECDSA secp256k1)")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Logf("-- ValidUntil: %s (Unix: %d)", validUntil.Format(time.RFC3339), validUntil.Unix())
	t.Log("-- Signatures are sorted by signer address (required by MCMS)")
	t.Log("")
	t.Log("testSignatures : [RawSignature]")
	t.Log("testSignatures =")
	for i, sig := range signaturesRaw {
		comma := ","
		if i == len(signaturesRaw)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")

	t.Log("-- ===========================================================================")
	t.Log("-- SIGNER PUBLIC KEYS (for config setup)")
	t.Log("-- ===========================================================================")
	t.Log("")
	for i, s := range sortedSigners {
		t.Logf("-- Signer %d: address=%s", i+1, s.Address)
		t.Logf("--           pubkey=%s", s.PublicKey)
	}
	t.Log("")

	t.Log("-- ===========================================================================")
	t.Log("-- COMPUTED VALUES (for verification)")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Logf("-- Op leaf hash: %s", HashOpLeaf(proposal.Operations[0]))
	t.Logf("-- Metadata leaf hash: %s", HashMetadataLeaf(proposal.Metadata))
	signedHash := ComputeSignedHash(root, validUntil)
	t.Logf("-- Signed hash: %s", signedHash)
	t.Log("")

	// ======================================================================
	// Additional vectors for Daml tests: schedule_batch + bypasser_execute_batch
	// ======================================================================

	baseMcmsId = "mcms-daml-test"
	timelockInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, "ccip_owner-9cefe94d")
	proposerMultisigId := MakeMcmsId(timelockInstanceId, MCMSRoleProposer)
	bypasserMultisigId := MakeMcmsId(timelockInstanceId, MCMSRoleBypasser)

	// Create inner calls for schedule_batch - these are the operations to be timelocked
	scheduleInnerCalls := []TimelockCall{
		{
			TargetInstanceId: timelockInstanceId,
			FunctionName:     "update_min_delay",
			OperationData:    "120", // new min delay of 120 seconds
		},
	}
	scheduleParams := ScheduleBatchParams{
		Calls:       scheduleInnerCalls,
		Predecessor: ZeroHash, // no predecessor
		Salt:        "test-schedule-salt-1",
		DelaySecs:   0, // use minimum delay
	}
	encodedScheduleParams := EncodeScheduleBatchParams(scheduleParams)

	// Proposer root authorizing schedule_batch (self)
	scheduleProposal := NewMCMSProposal(chainId, proposerMultisigId, 0, false).
		AddOperation(timelockInstanceId, "schedule_batch", encodedScheduleParams).
		Build()
	scheduleRoot := scheduleProposal.GetRoot()
	scheduleMetadataProof, err := scheduleProposal.GetMetadataProof()
	require.NoError(t, err)
	scheduleOpProof, err := scheduleProposal.GetOpProof(0)
	require.NoError(t, err)
	scheduleSignatures, err := scheduleProposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- TIMELOCK SELF-DISPATCH VECTORS (Proposer -> schedule_batch)")
	t.Log("-- Copy/paste into FlowTest.daml")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("realTestScheduleMetadata : RootMetadata")
	t.Log("realTestScheduleMetadata = RootMetadata")
	t.Logf("  { chainId = %d", scheduleProposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", scheduleProposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", scheduleProposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", scheduleProposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", scheduleProposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")
	t.Logf("realTestScheduleRoot : Text")
	t.Logf("realTestScheduleRoot = \"%s\"", scheduleRoot)
	t.Log("")
	t.Log("realTestScheduleMetadataProof : [Text]")
	if len(scheduleMetadataProof) == 0 {
		t.Log("realTestScheduleMetadataProof = []")
	} else {
		t.Log("realTestScheduleMetadataProof =")
		for i, p := range scheduleMetadataProof {
			comma := ","
			if i == len(scheduleMetadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("realTestScheduleOpProof : [Text]")
	if len(scheduleOpProof) == 0 {
		t.Log("realTestScheduleOpProof = []")
	} else {
		t.Log("realTestScheduleOpProof =")
		for i, p := range scheduleOpProof {
			comma := ","
			if i == len(scheduleOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("realTestScheduleSignatures : [RawSignature]")
	t.Log("realTestScheduleSignatures =")
	for i, sig := range scheduleSignatures {
		comma := ","
		if i == len(scheduleSignatures)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")
	t.Log("realTestScheduleOp : Op")
	t.Log("realTestScheduleOp = Op")
	t.Logf("  { chainId = %d", scheduleProposal.Operations[0].ChainId)
	t.Logf("  , multisigId = \"%s\"", scheduleProposal.Operations[0].MultisigId)
	t.Logf("  , nonce = %d", scheduleProposal.Operations[0].Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", scheduleProposal.Operations[0].TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", scheduleProposal.Operations[0].FunctionName)
	t.Logf("  , operationData = \"%s\"", scheduleProposal.Operations[0].OperationData)
	t.Log("  }")
	t.Log("")
	t.Log("realTestScheduleOpProofPath : [Text]")
	if len(scheduleOpProof) == 0 {
		t.Log("realTestScheduleOpProofPath = []")
	} else {
		t.Log("realTestScheduleOpProofPath =")
		for i, p := range scheduleOpProof {
			comma := ","
			if i == len(scheduleOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")

	// Create inner calls for bypasser_execute_batch - operations to execute immediately
	bypasserInnerCalls := []TimelockCall{
		{
			TargetInstanceId: timelockInstanceId,
			FunctionName:     "update_min_delay",
			OperationData:    "300", // new min delay of 300 seconds
		},
	}
	bypasserParams := BypasserExecuteParams{
		Calls: bypasserInnerCalls,
	}
	encodedBypasserParams := EncodeBypasserExecuteParams(bypasserParams)

	// Bypasser root authorizing bypasser_execute_batch (self)
	bypasserProposal := NewMCMSProposal(chainId, bypasserMultisigId, 0, false).
		AddOperation(timelockInstanceId, "bypasser_execute_batch", encodedBypasserParams).
		Build()
	bypasserRoot := bypasserProposal.GetRoot()
	bypasserMetadataProof, err := bypasserProposal.GetMetadataProof()
	require.NoError(t, err)
	bypasserOpProof, err := bypasserProposal.GetOpProof(0)
	require.NoError(t, err)
	bypasserSignatures, err := bypasserProposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- BYPASSER SELF-DISPATCH VECTORS (Bypasser -> bypasser_execute_batch)")
	t.Log("-- Copy/paste into FlowTest.daml")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("realTestBypasserMetadata : RootMetadata")
	t.Log("realTestBypasserMetadata = RootMetadata")
	t.Logf("  { chainId = %d", bypasserProposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", bypasserProposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", bypasserProposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", bypasserProposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", bypasserProposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")
	t.Logf("realTestBypasserRoot : Text")
	t.Logf("realTestBypasserRoot = \"%s\"", bypasserRoot)
	t.Log("")
	t.Log("realTestBypasserMetadataProof : [Text]")
	if len(bypasserMetadataProof) == 0 {
		t.Log("realTestBypasserMetadataProof = []")
	} else {
		t.Log("realTestBypasserMetadataProof =")
		for i, p := range bypasserMetadataProof {
			comma := ","
			if i == len(bypasserMetadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("realTestBypasserOpProof : [Text]")
	if len(bypasserOpProof) == 0 {
		t.Log("realTestBypasserOpProof = []")
	} else {
		t.Log("realTestBypasserOpProof =")
		for i, p := range bypasserOpProof {
			comma := ","
			if i == len(bypasserOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("realTestBypasserSignatures : [RawSignature]")
	t.Log("realTestBypasserSignatures =")
	for i, sig := range bypasserSignatures {
		comma := ","
		if i == len(bypasserSignatures)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")
	t.Log("realTestBypasserOp : Op")
	t.Log("realTestBypasserOp = Op")
	t.Logf("  { chainId = %d", bypasserProposal.Operations[0].ChainId)
	t.Logf("  , multisigId = \"%s\"", bypasserProposal.Operations[0].MultisigId)
	t.Logf("  , nonce = %d", bypasserProposal.Operations[0].Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", bypasserProposal.Operations[0].TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", bypasserProposal.Operations[0].FunctionName)
	t.Logf("  , operationData = \"%s\"", bypasserProposal.Operations[0].OperationData)
	t.Log("  }")
	t.Log("")
	t.Log("realTestBypasserOpProofPath : [Text]")
	if len(bypasserOpProof) == 0 {
		t.Log("realTestBypasserOpProofPath = []")
	} else {
		t.Log("realTestBypasserOpProofPath =")
		for i, p := range bypasserOpProof {
			comma := ","
			if i == len(bypasserOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")

	// ======================================================================
	// Additional vectors for External Call Tests (Bypasser calls Counter)
	// ======================================================================

	counterInstanceId := "counter@ccip_owner-9cefe94d"

	// Bypasser op with external call to Counter (increment)
	externalBypasserCalls := []TimelockCall{
		{
			TargetInstanceId: counterInstanceId,
			FunctionName:     "increment",
			OperationData:    "",
		},
	}
	externalBypasserParams := BypasserExecuteParams{
		Calls: externalBypasserCalls,
	}
	encodedExternalBypasserParams := EncodeBypasserExecuteParams(externalBypasserParams)

	externalBypasserProposal := NewMCMSProposal(chainId, bypasserMultisigId, 0, false).
		AddOperation(timelockInstanceId, "bypasser_execute_batch", encodedExternalBypasserParams).
		Build()
	externalBypasserRoot := externalBypasserProposal.GetRoot()
	externalBypasserMetadataProof, err := externalBypasserProposal.GetMetadataProof()
	require.NoError(t, err)
	externalBypasserOpProof, err := externalBypasserProposal.GetOpProof(0)
	require.NoError(t, err)
	externalBypasserSignatures, err := externalBypasserProposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- EXTERNAL CALL VECTORS (Bypasser -> bypasser_execute_batch -> Counter)")
	t.Log("-- Copy/paste into ExternalTargetTest.daml")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("externalBypasserOp : Op")
	t.Log("externalBypasserOp = Op")
	t.Logf("  { chainId = %d", externalBypasserProposal.Operations[0].ChainId)
	t.Logf("  , multisigId = \"%s\"", externalBypasserProposal.Operations[0].MultisigId)
	t.Logf("  , nonce = %d", externalBypasserProposal.Operations[0].Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", externalBypasserProposal.Operations[0].TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", externalBypasserProposal.Operations[0].FunctionName)
	t.Logf("  , operationData = \"%s\"", externalBypasserProposal.Operations[0].OperationData)
	t.Log("  }")
	t.Log("")
	t.Log("externalBypasserMetadata : RootMetadata")
	t.Log("externalBypasserMetadata = RootMetadata")
	t.Logf("  { chainId = %d", externalBypasserProposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", externalBypasserProposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", externalBypasserProposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", externalBypasserProposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", externalBypasserProposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")
	t.Logf("externalBypasserRoot : Text")
	t.Logf("externalBypasserRoot = \"%s\"", externalBypasserRoot)
	t.Log("")
	t.Log("externalBypasserMetadataProof : [Text]")
	if len(externalBypasserMetadataProof) == 0 {
		t.Log("externalBypasserMetadataProof = []")
	} else {
		t.Log("externalBypasserMetadataProof =")
		for i, p := range externalBypasserMetadataProof {
			comma := ","
			if i == len(externalBypasserMetadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("externalBypasserOpProof : [Text]")
	if len(externalBypasserOpProof) == 0 {
		t.Log("externalBypasserOpProof = []")
	} else {
		t.Log("externalBypasserOpProof =")
		for i, p := range externalBypasserOpProof {
			comma := ","
			if i == len(externalBypasserOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("externalBypasserSignatures : [RawSignature]")
	t.Log("externalBypasserSignatures =")
	for i, sig := range externalBypasserSignatures {
		comma := ","
		if i == len(externalBypasserSignatures)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")

	// Scheduled external call vectors (Proposer -> schedule_batch -> Counter)
	externalScheduleCalls := []TimelockCall{
		{
			TargetInstanceId: counterInstanceId,
			FunctionName:     "increment",
			OperationData:    "",
		},
	}
	externalScheduleParams := ScheduleBatchParams{
		Calls:       externalScheduleCalls,
		Predecessor: ZeroHash,
		Salt:        "test-external-schedule-salt-1",
		DelaySecs:   0,
	}
	encodedExternalScheduleParams := EncodeScheduleBatchParams(externalScheduleParams)

	externalScheduleProposal := NewMCMSProposal(chainId, proposerMultisigId, 0, false).
		AddOperation(timelockInstanceId, "schedule_batch", encodedExternalScheduleParams).
		Build()
	externalScheduleRoot := externalScheduleProposal.GetRoot()
	externalScheduleMetadataProof, err := externalScheduleProposal.GetMetadataProof()
	require.NoError(t, err)
	externalScheduleOpProof, err := externalScheduleProposal.GetOpProof(0)
	require.NoError(t, err)
	externalScheduleSignatures, err := externalScheduleProposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- SCHEDULED EXTERNAL CALL VECTORS (Proposer -> schedule_batch -> Counter)")
	t.Log("-- Copy/paste into ExternalTargetTest.daml")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Log("externalScheduleOp : Op")
	t.Log("externalScheduleOp = Op")
	t.Logf("  { chainId = %d", externalScheduleProposal.Operations[0].ChainId)
	t.Logf("  , multisigId = \"%s\"", externalScheduleProposal.Operations[0].MultisigId)
	t.Logf("  , nonce = %d", externalScheduleProposal.Operations[0].Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", externalScheduleProposal.Operations[0].TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", externalScheduleProposal.Operations[0].FunctionName)
	t.Logf("  , operationData = \"%s\"", externalScheduleProposal.Operations[0].OperationData)
	t.Log("  }")
	t.Log("")
	t.Log("externalScheduleMetadata : RootMetadata")
	t.Log("externalScheduleMetadata = RootMetadata")
	t.Logf("  { chainId = %d", externalScheduleProposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", externalScheduleProposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", externalScheduleProposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", externalScheduleProposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", externalScheduleProposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")
	t.Logf("externalScheduleRoot : Text")
	t.Logf("externalScheduleRoot = \"%s\"", externalScheduleRoot)
	t.Log("")
	t.Log("externalScheduleMetadataProof : [Text]")
	if len(externalScheduleMetadataProof) == 0 {
		t.Log("externalScheduleMetadataProof = []")
	} else {
		t.Log("externalScheduleMetadataProof =")
		for i, p := range externalScheduleMetadataProof {
			comma := ","
			if i == len(externalScheduleMetadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("externalScheduleOpProof : [Text]")
	if len(externalScheduleOpProof) == 0 {
		t.Log("externalScheduleOpProof = []")
	} else {
		t.Log("externalScheduleOpProof =")
		for i, p := range externalScheduleOpProof {
			comma := ","
			if i == len(externalScheduleOpProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")
	t.Log("externalScheduleSignatures : [RawSignature]")
	t.Log("externalScheduleSignatures =")
	for i, sig := range externalScheduleSignatures {
		comma := ","
		if i == len(externalScheduleSignatures)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")
	t.Logf("externalScheduleSalt : Text")
	t.Logf("externalScheduleSalt = \"%s\"", externalScheduleParams.Salt)
	t.Log("")

	t.Log("=======================================================================")
	t.Log("END OF DAML TEST VALUES")
	t.Log("=======================================================================")
}

// TestMCMS_GenerateMcmsOpTestValues generates and prints test values for MCMS self-dispatch operations
// These values can be used in Daml unit tests for cross-verification
func TestMCMS_GenerateMcmsOpTestValues(t *testing.T) {
	t.Parallel()

	t.Log("=======================================================================")
	t.Log("GENERATING DAML TEST VALUES FOR MCMS OP (SELF-DISPATCH)")
	t.Log("=======================================================================")
	t.Log("")

	// Create signers (deterministic for reproducible tests)
	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	sortedSigners := SortSignersByAddress(signers)

	chainId := 1
	baseMcmsId := "mcms-daml-test"
	mcmsOpInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, "ccip_owner-9cefe94d")
	mcmsId := MakeMcmsId(mcmsOpInstanceId, MCMSRoleProposer)

	// Build signer config for the test
	config := New2of3Config(signers)

	// Prepare new config params (change from 2-of-3 to 1-of-3)
	newQuorums := make([]int, NumGroups)
	newQuorums[0] = 1 // Changed from 2 to 1
	newParents := make([]int, NumGroups)

	// Encode SetConfigParams into operationData (like Aptos BCS)
	setConfigParams := SetConfigParams{
		Signers:      config.Signers,
		GroupQuorums: newQuorums,
		GroupParents: newParents,
		ClearRoot:    false,
	}
	encodedParams := EncodeSetConfigParams(setConfigParams)

	// Build proposal with MCMS operation targeting instanceId
	// operationData contains the encoded params
	proposal := NewMCMSProposal(chainId, mcmsId, 0, false)
	proposal.AddOperation(mcmsOpInstanceId, "set_config", encodedParams)
	proposal.Build()

	root := proposal.GetRoot()
	metadataProof, err := proposal.GetMetadataProof()
	require.NoError(t, err)
	opProof, err := proposal.GetOpProof(0)
	require.NoError(t, err)

	// Sign with 2030 validUntil for long-lasting test values
	validUntil := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	signaturesRaw, err := proposal.Sign(validUntil, sortedSigners[:2])
	require.NoError(t, err)

	t.Log("-- ===========================================================================")
	t.Log("-- MCMS OP TEST VALUES (Copy to FlowTest.daml)")
	t.Log("-- ===========================================================================")
	t.Log("")

	t.Log("-- Signer config (2-of-3 quorum)")
	t.Log("mcmsOpTestSigners : [SignerInfo]")
	t.Log("mcmsOpTestSigners =")
	for i, s := range sortedSigners {
		comma := ","
		if i == len(sortedSigners)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s SignerInfo \"%s\" %d 0%s", bracket, s.Address, i, comma)
	}
	t.Log("  ]")
	t.Log("")

	t.Log("mcmsOpTestConfig : MultisigConfig")
	t.Log("mcmsOpTestConfig = MultisigConfig")
	t.Log("  { signers = mcmsOpTestSigners")
	t.Log("  , groupQuorums = [2] <> replicate 31 0  -- 2-of-3 quorum in group 0")
	t.Log("  , groupParents = replicate 32 0")
	t.Log("  }")
	t.Log("")

	t.Logf("mcmsOpTestChainId : Int")
	t.Logf("mcmsOpTestChainId = %d", chainId)
	t.Log("")

	t.Logf("mcmsOpTestMcmsId : Text")
	t.Logf("mcmsOpTestMcmsId = \"%s\"", mcmsId)
	t.Log("")

	op := proposal.Operations[0]
	t.Log("mcmsOpTestOp : Op")
	t.Log("mcmsOpTestOp = Op")
	t.Logf("  { chainId = %d", op.ChainId)
	t.Logf("  , multisigId = \"%s\"", op.MultisigId)
	t.Logf("  , nonce = %d", op.Nonce)
	t.Logf("  , targetInstanceId = \"%s\"", op.TargetInstanceId)
	t.Logf("  , functionName = \"%s\"", op.FunctionName)
	t.Logf("  , operationData = \"%s\"", op.OperationData)
	t.Log("  }")
	t.Log("")

	t.Log("mcmsOpTestMetadata : RootMetadata")
	t.Log("mcmsOpTestMetadata = RootMetadata")
	t.Logf("  { chainId = %d", proposal.Metadata.ChainId)
	t.Logf("  , multisigId = \"%s\"", proposal.Metadata.MultisigId)
	t.Logf("  , preOpCount = %d", proposal.Metadata.PreOpCount)
	t.Logf("  , postOpCount = %d", proposal.Metadata.PostOpCount)
	t.Logf("  , overridePreviousRoot = %v", proposal.Metadata.OverridePreviousRoot)
	t.Log("  }")
	t.Log("")

	t.Logf("mcmsOpTestRoot : Text")
	t.Logf("mcmsOpTestRoot = \"%s\"", root)
	t.Log("")

	t.Log("mcmsOpTestMetadataProof : [Text]")
	t.Log("mcmsOpTestMetadataProof =")
	if len(metadataProof) == 0 {
		t.Log("  []")
	} else {
		for i, p := range metadataProof {
			comma := ","
			if i == len(metadataProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")

	t.Log("mcmsOpTestOpProof : [Text]")
	t.Log("mcmsOpTestOpProof =")
	if len(opProof) == 0 {
		t.Log("  []")
	} else {
		for i, p := range opProof {
			comma := ","
			if i == len(opProof)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s \"%s\"%s", bracket, p, comma)
		}
		t.Log("  ]")
	}
	t.Log("")

	t.Log("-- Real ECDSA secp256k1 signatures (ValidUntil: 2030-01-01T00:00:00Z)")
	t.Log("mcmsOpTestSignatures : [RawSignature]")
	t.Log("mcmsOpTestSignatures =")
	for i, sig := range signaturesRaw {
		comma := ","
		if i == len(signaturesRaw)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s RawSignature", bracket)
		t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
		t.Logf("      , r = \"%s\"", sig.R)
		t.Logf("      , s = \"%s\"", sig.S)
		t.Logf("      }%s", comma)
	}
	t.Log("  ]")
	t.Log("")

	t.Log("mcmsOpTestValidUntil : Time")
	t.Log("mcmsOpTestValidUntil = datetime 2030 Jan 1 0 0 0")
	t.Log("")

	t.Log("-- ===========================================================================")
	t.Log("-- COMPUTED VALUES (for verification)")
	t.Log("-- ===========================================================================")
	t.Log("")
	t.Logf("-- Op leaf hash: %s", HashOpLeaf(proposal.Operations[0]))
	t.Logf("-- Metadata leaf hash: %s", HashMetadataLeaf(proposal.Metadata))
	signedHash := ComputeSignedHash(root, validUntil)
	t.Logf("-- Signed hash: %s", signedHash)
	t.Log("")

	t.Log("=======================================================================")
	t.Log("END OF MCMS OP DAML TEST VALUES")
	t.Log("=======================================================================")
}

// TestMCMS_GenerateTimelockTestValues generates test vectors for ScheduleBatch/BypasserExecuteBatch
// self-dispatch flows. Run this and copy output into FlowTest.daml.
func TestMCMS_GenerateTimelockTestValues(t *testing.T) {
	t.Parallel()

	t.Log("=======================================================================")
	t.Log("GENERATING DAML TEST VALUES FOR TIMELOCK SELF-DISPATCH")
	t.Log("=======================================================================")

	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	sortedSigners := SortSignersByAddress(signers)

	chainId := 1
	baseMcmsId := "mcms-daml-test"
	timelockInstanceId := fmt.Sprintf("%s@%s", baseMcmsId, "ccip_owner-9cefe94d")
	validUntil := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	// Print signer config
	t.Log("")
	t.Log("-- Signer config (2-of-3 quorum)")
	t.Log("timelockTestSigners : [SignerInfo]")
	t.Log("timelockTestSigners =")
	for i, s := range sortedSigners {
		comma := ","
		if i == len(sortedSigners)-1 {
			comma = ""
		}
		bracket := " "
		if i == 0 {
			bracket = "["
		}
		t.Logf("  %s SignerInfo \"%s\" %d 0%s", bracket, s.Address, i, comma)
	}
	t.Log("  ]")
	t.Log("")
	t.Log("timelockTestConfig : MultisigConfig")
	t.Log("timelockTestConfig = MultisigConfig")
	t.Log("  { signers = timelockTestSigners")
	t.Log("  , groupQuorums = [2] <> replicate 31 0")
	t.Log("  , groupParents = replicate 32 0")
	t.Log("  }")
	t.Log("")

	printProposal := func(prefix string, proposal *MCMSProposal) {
		root := proposal.GetRoot()
		metadataProof, err := proposal.GetMetadataProof()
		require.NoError(t, err)
		opProof, err := proposal.GetOpProof(0)
		require.NoError(t, err)
		sigs, err := proposal.Sign(validUntil, sortedSigners[:2])
		require.NoError(t, err)

		op := proposal.Operations[0]
		t.Logf("%sOp : Op", prefix)
		t.Logf("%sOp = Op", prefix)
		t.Logf("  { chainId = %d", op.ChainId)
		t.Logf("  , multisigId = \"%s\"", op.MultisigId)
		t.Logf("  , nonce = %d", op.Nonce)
		t.Logf("  , targetInstanceId = \"%s\"", op.TargetInstanceId)
		t.Logf("  , functionName = \"%s\"", op.FunctionName)
		t.Logf("  , operationData = \"%s\"", op.OperationData)
		t.Log("  }")
		t.Log("")

		t.Logf("%sMetadata : RootMetadata", prefix)
		t.Logf("%sMetadata = RootMetadata", prefix)
		t.Logf("  { chainId = %d", proposal.Metadata.ChainId)
		t.Logf("  , multisigId = \"%s\"", proposal.Metadata.MultisigId)
		t.Logf("  , preOpCount = %d", proposal.Metadata.PreOpCount)
		t.Logf("  , postOpCount = %d", proposal.Metadata.PostOpCount)
		t.Logf("  , overridePreviousRoot = %v", proposal.Metadata.OverridePreviousRoot)
		t.Log("  }")
		t.Log("")

		t.Logf("%sRoot : Text", prefix)
		t.Logf("%sRoot = \"%s\"", prefix, root)
		t.Log("")

		printList := func(name string, items []string) {
			t.Logf("%s%s : [Text]", prefix, name)
			if len(items) == 0 {
				t.Logf("%s%s = []", prefix, name)
			} else {
				t.Logf("%s%s =", prefix, name)
				for i, p := range items {
					comma := ","
					if i == len(items)-1 {
						comma = ""
					}
					bracket := " "
					if i == 0 {
						bracket = "["
					}
					t.Logf("  %s \"%s\"%s", bracket, p, comma)
				}
				t.Log("  ]")
			}
			t.Log("")
		}

		printList("MetadataProof", metadataProof)
		printList("OpProof", opProof)

		t.Logf("%sSignatures : [RawSignature]", prefix)
		t.Logf("%sSignatures =", prefix)
		for i, sig := range sigs {
			comma := ","
			if i == len(sigs)-1 {
				comma = ""
			}
			bracket := " "
			if i == 0 {
				bracket = "["
			}
			t.Logf("  %s RawSignature", bracket)
			t.Logf("      { publicKey = \"%s\"", sig.PublicKey)
			t.Logf("      , r = \"%s\"", sig.R)
			t.Logf("      , s = \"%s\"", sig.S)
			t.Logf("      }%s", comma)
		}
		t.Log("  ]")
		t.Log("")
	}

	// Proposer schedule_batch proposal
	t.Log("-- ===========================================================================")
	t.Log("-- PROPOSER SCHEDULE_BATCH VECTORS")
	t.Log("-- ===========================================================================")
	t.Log("")
	proposerMsId := MakeMcmsId(timelockInstanceId, MCMSRoleProposer)
	scheduleProposal := NewMCMSProposal(chainId, proposerMsId, 0, false).
		AddOperation(timelockInstanceId, "schedule_batch", "").
		Build()
	printProposal("timelockSchedule", scheduleProposal)

	// Bypasser proposal
	t.Log("-- ===========================================================================")
	t.Log("-- BYPASSER EXECUTE_BATCH VECTORS")
	t.Log("-- ===========================================================================")
	t.Log("")
	bypasserMsId := MakeMcmsId(timelockInstanceId, MCMSRoleBypasser)
	bypasserProposal := NewMCMSProposal(chainId, bypasserMsId, 0, false).
		AddOperation(timelockInstanceId, "bypasser_execute_batch", "").
		Build()
	printProposal("timelockBypasser", bypasserProposal)

	t.Log("timelockTestValidUntil : Time")
	t.Log("timelockTestValidUntil = datetime 2030 Jan 1 0 0 0")
	t.Log("")

	t.Log("=======================================================================")
	t.Log("END OF TIMELOCK DAML TEST VALUES")
	t.Log("=======================================================================")
}
