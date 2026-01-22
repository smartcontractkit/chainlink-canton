package tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

// TestMCMS_ExecuteOpFlow tests the complete MCMS execute flow with direct invocation:
// 1. Deploy MCMS and Counter contracts
// 2. Issue a target auth ticket
// 3. Configure signers (2-of-3)
// 4. Create proposal with "increment" operation
// 5. Sign with 2 signers
// 6. SetRoot with real signatures
// 7. ExecuteOp - validates ticket, then calls Counter
// 8. Verify counter value incremented
func TestMCMS_ExecuteOpFlow(t *testing.T) {
	// Setup context with JWT auth
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Connect to participant
	participant, err := NewParticipant(
		"participant1.admin-api.localhost:8080",
		"participant1.grpc-ledger-api.localhost:8080",
	)
	require.NoError(t, err)

	// ========================
	// |   Setup: Upload DAR  |
	// ========================

	t.Log("Uploading MCMS DAR...")

	mcmsDar, err := os.ReadFile("../../contracts/mcms/.daml/dist/mcms-1.0.0.dar")
	require.NoError(t, err)

	packageIDs, err := UploadDARstoMultipleParticipants(ctx, [][]byte{mcmsDar}, participant)
	require.NoError(t, err)
	t.Logf("Uploaded MCMS DAR, package IDs: %v", packageIDs)

	// ========================
	// |   Setup: Parties     |
	// ========================

	parties, err := EnsurePartyOnMultipleParticipants(ctx, participant)
	require.NoError(t, err)
	ccipOwner := parties[0]
	t.Logf("Using party: %s", ccipOwner)

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
	baseMcmsId := "mcms-integration-test-" + uuid.New().String()[:8]
	mcmsId := MakeMcmsId(baseMcmsId, MCMSRoleProposer)
	instanceId := "counter-env-" + uuid.New().String()[:8] // Stable instance ID for Counter
	t.Logf("MCMS ID: %s", mcmsId)
	t.Logf("Counter Instance ID: %s", instanceId)

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
	for i := 0; i < NumGroups; i++ {
		groupQuorumValues[i] = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(config.GroupQuorums[i])}}
	}

	// Build group parents (32 ints)
	groupParentValues := make([]*apiv2.Value, NumGroups)
	for i := 0; i < NumGroups; i++ {
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

	// Create MCMS contract (simplified - no ticket tracking)
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								{Label: "role", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "mcmsId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: mcmsId}}},
								{Label: "config", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
									Fields: []*apiv2.RecordField{
										{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
										{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
										{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
									},
								}}}},
								{Label: "seenHashes", Value: emptyMap},
								{Label: "expiringRoot", Value: emptyExpiringRoot},
								{Label: "rootMetadata", Value: emptyRootMetadata},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// ========================
	// |   2. Create Counter  |
	// ========================

	t.Log("Creating Counter contract with MCMSReceiver interface...")

	counterCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Counter",
								EntityName: "Counter",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceId}}},
								{Label: "value", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)
	counterCid := counterCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created Counter contract: %s", counterCid)

	// ====================================
	// |   2.5. Issue Auth Ticket         |
	// ====================================

	t.Log("Issuing MCMS target auth ticket...")

	issueTicketRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.MCMSReceiver",
								EntityName: "MCMSReceiver",
							},
							ContractId: counterCid,
							Choice:     "MCMSReceiver_IssueAuthTicket",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)

	var authTicketCid string
	for _, event := range issueTicketRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "TargetAuthTicket" {
			authTicketCid = created.GetContractId()
			break
		}
	}
	require.NotEmpty(t, authTicketCid, "Should have created a TargetAuthTicket")
	t.Logf("Issued auth ticket: %s", authTicketCid)

	// ========================
	// |   3. Build Proposal  |
	// ========================

	t.Log("Building proposal...")

	// Build proposal with one "increment" operation
	// Note: targetInstanceId is the instanceId, not the contract ID
	proposal := NewMCMSProposal(int(chainId), mcmsId, 0, false)
	proposal.AddOperation(instanceId, "increment", "")
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

	t.Log("Calling SetRoot...")

	setRootRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
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
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)

	// Get new MCMS contract ID from exercise result
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
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

	executeOpRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "ExecuteOp",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
									{Label: "targetCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: counterCid}}},
									{Label: "authTicketCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: authTicketCid}}},
									{Label: "op", Value: opValue},
									{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: opProofValues}}}},
									{Label: "contractIds", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{Sum: &apiv2.Value_ContractId{ContractId: counterCid}}, // Pass Counter CID for target to use
									}}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
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
				for _, field := range created.GetCreateArguments().GetFields() {
					if field.GetLabel() == "value" {
						counterValue = field.GetValue().GetInt64()
						break
					}
				}
			}

			// Check for MCMSEntrypointEvent (emitted by Counter.mcmsEntrypoint)
			if entityName == "MCMSEntrypointEvent" {
				eventFound = true
				for _, field := range created.GetCreateArguments().GetFields() {
					switch field.GetLabel() {
					case "instanceId":
						eventInstanceId = field.GetValue().GetText()
					case "functionName":
						eventFunctionName = field.GetValue().GetText()
					case "operationData":
						eventOperationData = field.GetValue().GetText()
					case "contractIdsAsText":
						if listValue := field.GetValue().GetList(); listValue != nil {
							for _, elem := range listValue.GetElements() {
								eventContractIdsAsText = append(eventContractIdsAsText, elem.GetText())
							}
						}
					}
				}
			}
		}
	}
	t.Logf("ExecuteOp succeeded, counter value from event: %d", counterValue)
	require.NotEmpty(t, newCounterCid, "Should have new counter contract ID")

	// Verify MCMSEntrypointEvent was emitted with correct data
	require.True(t, eventFound, "MCMSEntrypointEvent should be emitted by Counter.mcmsEntrypoint")
	assert.Equal(t, instanceId, eventInstanceId, "Event instanceId should match Counter instanceId")
	assert.Equal(t, "increment", eventFunctionName, "Event functionName should be 'increment'")
	assert.Equal(t, "", eventOperationData, "Event operationData should be empty for increment")

	// Verify contractIds were passed through the chain
	require.Len(t, eventContractIdsAsText, 1, "Should have 1 contractId passed through")
	t.Logf("MCMSEntrypointEvent verified: instanceId=%s, functionName=%s, contractIdsAsText=%v",
		eventInstanceId, eventFunctionName, eventContractIdsAsText)
	// Note: The show representation in Daml may wrap the CID, so we just check it's not empty
	assert.NotEmpty(t, eventContractIdsAsText[0], "contractIdsAsText[0] should contain the serialized Counter CID")

	// ========================
	// |   7. Verify Counter  |
	// ========================

	// Also query the ACS to verify the counter value
	t.Log("Querying counter via ACS...")
	counterContracts, err := GetActiveContractsForPartyTemplateId(ctx, participant, ccipOwner, &apiv2.Identifier{
		PackageId:  "#mcms",
		ModuleName: "MCMS.Counter",
		EntityName: "Counter",
	})
	require.NoError(t, err)

	var queriedValue int64 = -1
	for _, contract := range counterContracts {
		if contract.GetCreatedEvent().GetContractId() == newCounterCid {
			for _, field := range contract.GetCreatedEvent().GetCreateArguments().GetFields() {
				if field.GetLabel() == "value" {
					queriedValue = field.GetValue().GetInt64()
					break
				}
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
	t.Log("  3. Issued target auth ticket")
	t.Log("  4. Built proposal with 'increment' operation targeting instanceId")
	t.Log("  5. Signed with 2 signers (real ECDSA signatures)")
	t.Log("  6. SetRoot with on-chain verification")
	t.Log("  7. ExecuteOp - validates ticket and calls Counter.MCMSReceiver_Entrypoint")
	t.Log("  8. Counter value = 1 ✓")
}

// TestMCMS_SignatureVerificationFails tests that invalid signatures are rejected
func TestMCMS_SignatureVerificationFails(t *testing.T) {
	// Setup context with JWT auth
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Connect to participant
	participant, err := NewParticipant(
		"participant1.admin-api.localhost:8080",
		"participant1.grpc-ledger-api.localhost:8080",
	)
	require.NoError(t, err)

	// Upload DAR
	mcmsDar, err := os.ReadFile("../../contracts/mcms/.daml/dist/mcms-1.0.0.dar")
	require.NoError(t, err)
	_, err = UploadDARstoMultipleParticipants(ctx, [][]byte{mcmsDar}, participant)
	require.NoError(t, err)

	// Setup party
	parties, err := EnsurePartyOnMultipleParticipants(ctx, participant)
	require.NoError(t, err)
	ccipOwner := parties[0]

	// Create signers
	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	config := New2of3Config(signers)

	chainId := int64(1)
	baseMcmsId := "mcms-sig-fail-test-" + uuid.New().String()[:8]
	mcmsId := MakeMcmsId(baseMcmsId, MCMSRoleProposer)

	// Build config values
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
	for i := 0; i < NumGroups; i++ {
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

	// Create MCMS contract
	t.Log("Creating MCMS contract...")
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								{Label: "role", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "mcmsId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: mcmsId}}},
								{Label: "config", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
									Fields: []*apiv2.RecordField{
										{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
										{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
										{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
									},
								}}}},
								{Label: "seenHashes", Value: emptyMap},
								{Label: "expiringRoot", Value: emptyExpiringRoot},
								{Label: "rootMetadata", Value: emptyRootMetadata},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// Build a valid proposal
	proposal := NewMCMSProposal(int(chainId), mcmsId, 0, false)
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
	_, err = participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
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
			ActAs: []string{ccipOwner},
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

// TestMCMS_ReplayProtection tests that the same root cannot be set twice
func TestMCMS_ReplayProtection(t *testing.T) {
	// Setup context with JWT auth
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Connect to participant
	participant, err := NewParticipant(
		"participant1.admin-api.localhost:8080",
		"participant1.grpc-ledger-api.localhost:8080",
	)
	require.NoError(t, err)

	// Upload DAR
	mcmsDar, err := os.ReadFile("../../contracts/mcms/.daml/dist/mcms-1.0.0.dar")
	require.NoError(t, err)
	_, err = UploadDARstoMultipleParticipants(ctx, [][]byte{mcmsDar}, participant)
	require.NoError(t, err)

	// Setup party
	parties, err := EnsurePartyOnMultipleParticipants(ctx, participant)
	require.NoError(t, err)
	ccipOwner := parties[0]

	// Create signers
	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	config := New2of3Config(signers)
	sortedSigners := SortSignersByAddress(signers)

	chainId := int64(1)
	baseMcmsId := "mcms-replay-test-" + uuid.New().String()[:8]
	mcmsId := MakeMcmsId(baseMcmsId, MCMSRoleProposer)

	// Build config values
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
	for i := 0; i < NumGroups; i++ {
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

	// Create MCMS contract
	t.Log("Creating MCMS contract...")
	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								{Label: "role", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "mcmsId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: mcmsId}}},
								{Label: "config", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
									Fields: []*apiv2.RecordField{
										{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
										{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
										{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
									},
								}}}},
								{Label: "seenHashes", Value: emptyMap},
								{Label: "expiringRoot", Value: emptyExpiringRoot},
								{Label: "rootMetadata", Value: emptyRootMetadata},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)
	mcmsCid := mcmsCreateRes.GetTransaction().GetEvents()[0].GetCreated().GetContractId()
	t.Logf("Created MCMS contract: %s", mcmsCid)

	// Build proposal
	proposal := NewMCMSProposal(int(chainId), mcmsId, 0, false)
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
	setRootRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
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
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)
	t.Log("First SetRoot succeeded")

	// Get new MCMS contract ID
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			mcmsCid = created.GetContractId()
			break
		}
	}

	// Second SetRoot with SAME signatures - should fail with E_ALREADY_SEEN_HASH
	t.Log("Second SetRoot call with same signatures (should fail)...")
	_, err = participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
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
			ActAs: []string{ccipOwner},
		},
	})

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "E_ALREADY_SEEN_HASH"),
		"Expected E_ALREADY_SEEN_HASH error, got: %v", err)

	t.Log("✓ Second SetRoot correctly rejected with E_ALREADY_SEEN_HASH (replay protection working)")
}

// TestMCMS_GenerateDamlTestValues generates all cryptographic values needed for Daml unit tests.
// Run this test and copy the output to contracts/mcms/test/daml/MCMS/FlowTest.daml
// This uses FIXED values (not random) so the output is deterministic.
func TestMCMS_GenerateDamlTestValues(t *testing.T) {
	t.Log("=======================================================================")
	t.Log("GENERATING DAML TEST VALUES")
	t.Log("Copy these values to contracts/mcms/test/daml/MCMS/FlowTest.daml")
	t.Log("=======================================================================")

	// Use fixed seed for deterministic output
	// Create 3 signers
	signer1, err := NewMCMSSigner()
	require.NoError(t, err)
	signer2, err := NewMCMSSigner()
	require.NoError(t, err)
	signer3, err := NewMCMSSigner()
	require.NoError(t, err)

	signers := []*MCMSSigner{signer1, signer2, signer3}
	config := New2of3Config(signers)
	sortedSigners := SortSignersByAddress(signers)

	// Fixed test values
	chainId := 1
	mcmsId := "mcms-daml-test-proposer"

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
	proposal := NewMCMSProposal(chainId, mcmsId, 0, false)
	proposal.AddOperation("counter", "increment", "")
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
	t.Logf("testMcmsId : Text")
	t.Logf("testMcmsId = \"%s\"", mcmsId)
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

	t.Log("=======================================================================")
	t.Log("END OF DAML TEST VALUES")
	t.Log("=======================================================================")
}

// TestMCMS_GenerateMcmsOpTestValues generates and prints test values for MCMS self-dispatch operations
// These values can be used in Daml unit tests for cross-verification
func TestMCMS_GenerateMcmsOpTestValues(t *testing.T) {
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
	mcmsId := "mcms-daml-mcmsop-test-proposer"

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

	// Build proposal with MCMS operation targeting "self"
	// operationData contains the encoded params
	proposal := NewMCMSProposal(chainId, mcmsId, 0, false)
	proposal.AddOperation("self", "set_config", encodedParams)
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

// TestMCMS_ExecuteMcmsOp tests self-dispatch MCMS operations (Aptos pattern)
// This demonstrates changing MCMS config via a signed proposal
func TestMCMS_ExecuteMcmsOp(t *testing.T) {
	// Setup context with JWT auth
	jwToken, err := getJWT()
	require.NoError(t, err)
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Connect to participant
	participant, err := NewParticipant(
		"participant1.admin-api.localhost:8080",
		"participant1.grpc-ledger-api.localhost:8080",
	)
	require.NoError(t, err)

	// Connect to random user participant
	randomUserParticipant, err := NewParticipant(
		"participant2.admin-api.localhost:8080",
		"participant2.grpc-ledger-api.localhost:8080",
	)
	require.NoError(t, err)

	// Upload DAR
	t.Log("Uploading MCMS DAR...")
	mcmsDar, err := os.ReadFile("../../contracts/mcms/.daml/dist/mcms-1.0.0.dar")
	require.NoError(t, err)
	_, err = UploadDARstoMultipleParticipants(ctx, [][]byte{mcmsDar}, participant)
	require.NoError(t, err)

	_, err = UploadDARstoMultipleParticipants(ctx, [][]byte{mcmsDar}, randomUserParticipant)
	require.NoError(t, err)

	// Setup party
	parties, err := EnsurePartyOnMultipleParticipants(ctx, participant, randomUserParticipant)
	require.NoError(t, err)
	ccipOwner := parties[0]
	randomUser := parties[1]
	t.Logf("Using party: %s", ccipOwner)
	t.Logf("Using party: %s", randomUser)

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
	sortedSigners := SortSignersByAddress(signers)

	// ========================
	// |   Contract Constants |
	// ========================

	chainId := int64(1)
	baseMcmsId := "mcms-op-test-" + uuid.New().String()[:8]
	mcmsId := MakeMcmsId(baseMcmsId, MCMSRoleProposer)
	t.Logf("MCMS ID: %s", mcmsId)

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
	for i := 0; i < NumGroups; i++ {
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

	mcmsCreateRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
								{Label: "role", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "Proposer"}}}},
								{Label: "chainId", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: chainId}}},
								{Label: "mcmsId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: mcmsId}}},
								{Label: "config", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
									Fields: []*apiv2.RecordField{
										{Label: "signers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: signerInfoValues}}}},
										{Label: "groupQuorums", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupQuorumValues}}}},
										{Label: "groupParents", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: groupParentValues}}}},
									},
								}}}},
								{Label: "seenHashes", Value: emptyMap},
								{Label: "expiringRoot", Value: emptyExpiringRoot},
								{Label: "rootMetadata", Value: emptyRootMetadata},
							}},
						},
					},
				},
			},
			ActAs: []string{ccipOwner},
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

	// Build proposal with MCMS operation targeting "self"
	// operationData contains the encoded params (like Aptos)
	proposal := NewMCMSProposal(int(chainId), mcmsId, 0, false)
	proposal.AddOperation("self", "set_config", encodedParams) // Params encoded in operationData
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

	setRootRes, err := participant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "SetRoot",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ccipOwner}}},
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
			ActAs: []string{ccipOwner},
		},
	})
	require.NoError(t, err)

	// Get new contract ID from SetRoot result
	for _, event := range setRootRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil && created.GetTemplateId().GetEntityName() == "MCMS" {
			mcmsCid = created.GetContractId()
			break
		}
	}
	t.Logf("SetRoot succeeded, new MCMS CID: %s", mcmsCid)

	// Query ACS to get disclosed contract with CreatedEventBlob (for bob to use)
	// Transaction events don't include the blob by default, so we query the ACS
	disclosedMcms, err := QueryDisclosedContract(ctx, mcmsCid, participant)
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
	// randomUser (bob) submits via randomUserParticipant with disclosed contract
	_, err = randomUserParticipant.CommandServiceClient.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#mcms",
								ModuleName: "MCMS.Main",
								EntityName: "MCMS",
							},
							ContractId: mcmsCid,
							Choice:     "ExecuteMcmsOp",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{
								Fields: []*apiv2.RecordField{
									{Label: "submitter", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: randomUser}}},
									{Label: "op", Value: opValue},
									{Label: "opProof", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: opProofValues}}}},
								},
							}}},
						},
					},
				},
			},
			ActAs:              []string{randomUser},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedMcms}, // Grant bob visibility
		},
	})
	require.NoError(t, err)
	t.Log("ExecuteMcmsOp succeeded (bob executed with disclosed contract)")

	// Bob can't see the created event (not an observer), so query from alice's participant
	// The old contract should be archived, so find the active one
	t.Log("Querying ACS from alice's participant to find the new contract...")

	// ========================
	// |   6. Verify Config   |
	// ========================

	// Query via GetActiveContracts with verbose to get the actual config
	// Find the MCMS contract with our mcmsId (the old one is archived, new one is active)
	var newNumSigners int64 = -1
	var newQuorum int64 = -1
	var newMcmsCid string
	offset, _ := GetCurrentOffset(ctx, participant)
	acsRes, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersForAnyParty: &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{
				{
					IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
						WildcardFilter: &apiv2.WildcardFilter{
							IncludeCreatedEventBlob: true,
						},
					},
				},
			}},
			Verbose: true,
		},
	})
	require.NoError(t, err)
	defer acsRes.CloseSend()

	for {
		ac, err := acsRes.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if c, ok := ac.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			if c.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName() == "MCMS" {
				// Check if this is our MCMS by matching mcmsId
				for _, field := range c.ActiveContract.GetCreatedEvent().GetCreateArguments().GetFields() {
					if field.GetLabel() == "mcmsId" && field.GetValue().GetText() == mcmsId {
						newMcmsCid = c.ActiveContract.GetCreatedEvent().ContractId
						for _, f := range c.ActiveContract.GetCreatedEvent().GetCreateArguments().GetFields() {
							if f.GetLabel() == "config" {
								configRecord := f.GetValue().GetRecord()
								for _, configField := range configRecord.GetFields() {
									if configField.GetLabel() == "signers" {
										newNumSigners = int64(len(configField.GetValue().GetList().GetElements()))
									}
									if configField.GetLabel() == "groupQuorums" {
										quorums := configField.GetValue().GetList().GetElements()
										if len(quorums) > 0 {
											newQuorum = quorums[0].GetInt64()
										}
									}
								}
							}
						}
						break
					}
				}
			}
		}
		if newMcmsCid != "" {
			break
		}
	}
	require.NotEmpty(t, newMcmsCid, "Should find new MCMS contract in ACS")
	t.Logf("Found new MCMS contract: %s", newMcmsCid)
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
