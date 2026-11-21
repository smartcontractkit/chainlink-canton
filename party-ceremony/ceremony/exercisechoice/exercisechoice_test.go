package exercisechoice_test

import (
	"context"
	"encoding/base64"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/exercisechoice"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// fakeAdmin implements only the Admin API calls the ceremony makes; any other
// call panics on the nil embedded interface.
type fakeAdmin struct {
	client.CantonClient
	participantID string
}

func (f *fakeAdmin) GetParticipantUID(context.Context) (string, error) { return f.participantID, nil }

func (f *fakeAdmin) GetP2P(context.Context, string, string) (*client.P2PState, error) {
	return &client.P2PState{
		Participants:     []client.P2PParticipantInfo{{ParticipantUID: "p1"}, {ParticipantUID: "p2"}},
		PartySigningKeys: &client.P2PSigningKeysInfo{Keys: []string{"k1", "k2"}, Threshold: 2},
	}, nil
}

// fakeLedger records submitted commands and executions.
type fakeLedger struct {
	client.LedgerClient
	prepared      [][]*apiv2.Command
	disclosed     [][]*apiv2.DisclosedContract
	actAs, readAs []string
	executions    int
	tamper        func(*interactive.PreparedTransaction)
}

func (f *fakeLedger) PartyExists(context.Context, string) (bool, error) { return true, nil }

func (f *fakeLedger) PrepareSubmission(_ context.Context, commands []*apiv2.Command, _, _ []string, _ string, disclosed []*apiv2.DisclosedContract) (*interactive.PrepareSubmissionResponse, error) {
	f.prepared = append(f.prepared, commands)
	f.disclosed = append(f.disclosed, disclosed)
	tx := exerciseTx(commands[0].GetExercise())
	if f.tamper != nil {
		f.tamper(tx)
	}

	return &interactive.PrepareSubmissionResponse{PreparedTransactionHash: []byte{0xab, 0xcd}, PreparedTransaction: tx}, nil
}

func (f *fakeLedger) GrantActAs(_ context.Context, userID, _ string) error {
	f.actAs = append(f.actAs, userID)

	return nil
}

func (f *fakeLedger) GrantReadAs(_ context.Context, userID, _ string) error {
	f.readAs = append(f.readAs, userID)

	return nil
}

func exerciseTx(cmd *apiv2.ExerciseCommand) *interactive.PreparedTransaction {
	chosen := proto.Clone(cmd.GetChoiceArgument()).(*apiv2.Value)
	for _, f := range chosen.GetRecord().GetFields() {
		f.Label = ""
	}

	return &interactive.PreparedTransaction{Transaction: &interactive.DamlTransaction{
		Roots: []string{"0"},
		Nodes: []*interactive.DamlTransaction_Node{{
			NodeId: "0",
			VersionedNode: &interactive.DamlTransaction_Node_V1{V1: &v1.Node{NodeType: &v1.Node_Exercise{Exercise: &v1.Exercise{
				TemplateId:  cmd.GetTemplateId(),
				ContractId:  cmd.GetContractId(),
				ChoiceId:    cmd.GetChoice(),
				ChosenValue: chosen,
			}}}},
		}},
	}}
}

func (f *fakeLedger) ExecuteSubmission(context.Context, *interactive.PreparedTransaction, *interactive.PartySignatures, interactive.HashingSchemeVersion) (client.SubmissionResult, error) {
	f.executions++

	return client.SubmissionResult{
		UpdateID: "update-1",
		Created: []client.CreatedContract{
			{TemplateID: "pkg:Utility.Registry.App.V0.Service.Registrar:RegistrarService", ContractID: "registrar-cid"},
			{TemplateID: "pkg:Utility.Registry.App.V0.Service.AllocationFactory:AllocationFactory", ContractID: "factory-cid"},
		},
	}, nil
}

type fakeSigner struct{ signed *int }

func (s fakeSigner) Sign(_ context.Context, hash []byte) (*apiv2.Signature, error) {
	if s.signed != nil {
		*s.signed++
	}

	return &apiv2.Signature{Signature: hash, SignedBy: "fp"}, nil
}

func deps(participantID string, lc *fakeLedger) ledger.ContractDeployDeps {
	return ledger.ContractDeployDeps{
		AdminClient:  &fakeAdmin{participantID: participantID},
		LedgerClient: lc,
		Signer:       fakeSigner{},
		Logger:       logger.Nop(),
	}
}

func input() exercisechoice.ExerciseChoiceInput {
	return exercisechoice.ExerciseChoiceInput{
		DecentralizedPartyID: "ccipOwner::1220abcd",
		SynchronizerID:       "global",
		PackageID:            "#utility-registry-app-v0",
		TemplateModule:       "Utility.Registry.App.V0.Service.Provider",
		TemplateEntity:       "ProviderService",
		ContractID:           "provider-service-cid",
		Choice:               "ProviderService_AcceptRegistrarServiceRequest",
		ChoiceArgument:       `{"record":{"fields":[{"label":"cid","value":{"contractId":"request-cid"}}]}}`,
	}
}

func TestExerciseChoiceSequence_CollectsEverySignatureThenSubmitsTheExercise(t *testing.T) {
	t.Parallel()

	lc := &fakeLedger{}
	reporter := operations.NewMemoryReporter()
	run := func(pid string) (operations.SequenceReport[exercisechoice.ExerciseChoiceInput, exercisechoice.ExerciseChoiceOutput], error) {
		return operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), reporter), exercisechoice.ExerciseChoiceSequence, deps(pid, lc), input())
	}

	_, err := run("p1")
	require.ErrorContains(t, err, exercisechoice.ErrThresholdNotMet.Error(), "p1 signs, p2 has not signed yet")
	assert.Zero(t, lc.executions, "nothing is submitted before every participant signs")

	sr, err := run("p2")
	require.NoError(t, err)
	assert.Equal(t, 1, lc.executions)
	assert.Equal(t, "update-1", sr.Output.UpdateID)
	assert.Equal(t, []string{"registrar-cid", "factory-cid"}, []string{sr.Output.CreatedContracts[0].ContractID, sr.Output.CreatedContracts[1].ContractID})
	assert.Equal(t, exercisechoice.PhaseCompleted, sr.Output.State.Phase)
	assert.ElementsMatch(t, []string{"p1", "p2"}, sr.Output.State.Signed)

	require.Len(t, lc.prepared, 1, "the exercise is prepared once and the cached hash is reused")
	require.Len(t, lc.prepared[0], 1)
	exercise := lc.prepared[0][0].GetExercise()
	require.NotNil(t, exercise, "the ceremony submits an exercise command, not a create")
	assert.Equal(t, "#utility-registry-app-v0", exercise.GetTemplateId().GetPackageId())
	assert.Equal(t, "Utility.Registry.App.V0.Service.Provider", exercise.GetTemplateId().GetModuleName())
	assert.Equal(t, "ProviderService", exercise.GetTemplateId().GetEntityName())
	assert.Equal(t, "provider-service-cid", exercise.GetContractId())
	assert.Equal(t, "ProviderService_AcceptRegistrarServiceRequest", exercise.GetChoice())
	wantArg := &apiv2.Value{}
	require.NoError(t, protojson.Unmarshal([]byte(input().ChoiceArgument), wantArg))
	assert.True(t, proto.Equal(wantArg, exercise.GetChoiceArgument()), "choice argument round-trips from proto JSON")

	cached, err := run("p1")
	require.NoError(t, err, "a completed ceremony re-runs from cache")
	assert.Equal(t, 1, lc.executions, "a re-run does not submit again")
	assert.Equal(t, sr.Output.UpdateID, cached.Output.UpdateID)

	inspected, found := exercisechoice.QueryState(reporter, input())
	require.True(t, found, "inspect finds the completed run")
	assert.Equal(t, "update-1", inspected.UpdateID)
	assert.Len(t, inspected.CreatedContracts, 2)

	other := input()
	other.Choice = "ProviderService_CreateProviderConfiguration"
	_, found = exercisechoice.QueryState(reporter, other)
	assert.False(t, found, "inspect matches the full input, not just the party")
}

func TestExerciseChoiceSequence_RejectsMalformedChoiceArgument(t *testing.T) {
	t.Parallel()

	in := input()
	in.ChoiceArgument = `{"fields":[]}`
	lc := &fakeLedger{}
	reporter := operations.NewMemoryReporter()
	run := func(pid string) error {
		_, err := operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), reporter), exercisechoice.ExerciseChoiceSequence, deps(pid, lc), in)
		return err
	}

	require.ErrorContains(t, run("p1"), "parsing choice argument")
	assert.Empty(t, lc.prepared, "nothing reaches the ledger")
}

func TestExerciseChoiceSequence_SendsDisclosedContractsWithThePrepare(t *testing.T) {
	t.Parallel()

	in := input()
	in.DisclosedContracts = []ledger.DisclosedContract{{
		TemplateID:       "#utility-registry-v0:Utility.Registry.V0.Configuration.Registry:RegistryConfiguration",
		ContractID:       "registry-config-cid",
		CreatedEventBlob: base64.StdEncoding.EncodeToString([]byte("registry-config-blob")),
		SynchronizerID:   "global-domain::1220",
	}}
	lc := &fakeLedger{}
	_, err := operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter()),
		exercisechoice.ExerciseChoiceSequence, deps("p1", lc), in)
	require.ErrorContains(t, err, exercisechoice.ErrThresholdNotMet.Error())

	want := []*apiv2.DisclosedContract{{
		TemplateId: &apiv2.Identifier{
			PackageId:  "#utility-registry-v0",
			ModuleName: "Utility.Registry.V0.Configuration.Registry",
			EntityName: "RegistryConfiguration",
		},
		ContractId:       "registry-config-cid",
		CreatedEventBlob: []byte("registry-config-blob"),
		SynchronizerId:   "global-domain::1220",
	}}
	require.Len(t, lc.disclosed, 1)
	require.Len(t, lc.disclosed[0], 1)
	assert.True(t, proto.Equal(want[0], lc.disclosed[0][0]), "got %v", lc.disclosed[0][0])
}

func TestExerciseChoiceSequence_GrantsActAsOnlyToTheRunningParticipant(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	ledgers := map[string]*fakeLedger{"p1": {}, "p2": {}}
	run := func(pid string) error {
		d := deps(pid, ledgers[pid])
		d.UserID = pid + "-user"
		_, err := operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), reporter), exercisechoice.ExerciseChoiceSequence, d, input())

		return err
	}

	require.ErrorContains(t, run("p1"), "1/2 participants have signed", "p1 prepares and signs without waiting for p2's grant")
	assert.Len(t, ledgers["p1"].prepared, 1)
	assert.Equal(t, []string{"p1-user"}, ledgers["p1"].actAs)
	assert.Empty(t, ledgers["p2"].actAs)

	require.NoError(t, run("p2"))
	assert.Equal(t, 1, ledgers["p2"].executions, "p2 resumes last and submits")
	assert.Equal(t, []string{"p1-user"}, ledgers["p1"].actAs)
	assert.Equal(t, []string{"p2-user"}, ledgers["p2"].actAs)
	assert.Equal(t, []string{"p1-user"}, ledgers["p1"].readAs)
	assert.Equal(t, []string{"p2-user"}, ledgers["p2"].readAs)
}

func TestExerciseChoiceSequence_RefusesToSignAPreparedTransactionThatDoesNotMatch(t *testing.T) {
	t.Parallel()

	exercise := func(tx *interactive.PreparedTransaction) *v1.Exercise {
		return tx.GetTransaction().GetNodes()[0].GetV1().GetExercise()
	}
	for name, tamper := range map[string]func(*interactive.PreparedTransaction){
		"other choice":   func(tx *interactive.PreparedTransaction) { exercise(tx).ChoiceId = "Archive" },
		"other contract": func(tx *interactive.PreparedTransaction) { exercise(tx).ContractId = "other-cid" },
		"other template": func(tx *interactive.PreparedTransaction) { exercise(tx).TemplateId.EntityName = "RegistrarService" },
		"other argument": func(tx *interactive.PreparedTransaction) {
			exercise(tx).ChosenValue.GetRecord().Fields[0].Value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: "other-request-cid"}}
		},
		"two roots": func(tx *interactive.PreparedTransaction) {
			tx.Transaction.Roots = append(tx.Transaction.Roots, "1")
		},
		"root is a create": func(tx *interactive.PreparedTransaction) {
			tx.Transaction.Nodes[0].VersionedNode = &interactive.DamlTransaction_Node_V1{V1: &v1.Node{NodeType: &v1.Node_Create{Create: &v1.Create{}}}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			signed := 0
			d := deps("p1", &fakeLedger{tamper: tamper})
			d.Signer = fakeSigner{signed: &signed}
			_, err := operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter()),
				exercisechoice.ExerciseChoiceSequence, d, input())

			require.ErrorContains(t, err, ledger.ErrPreparedTxMismatch.Error())
			assert.NotContains(t, err.Error(), exercisechoice.ErrThresholdNotMet.Error(), "a refusal is not reported as pending signers")
			assert.Zero(t, signed, "the signer is never called")
		})
	}
}

type recordingConfirmer struct {
	ceremony.Confirmer
	got []ceremony.DAMLSignDetail
}

func (r *recordingConfirmer) ConfirmDAMLSign(_ context.Context, d ceremony.DAMLSignDetail) error {
	r.got = append(r.got, d)

	return nil
}

func TestExerciseChoiceSequence_ShowsTheChoiceWhenConfirmingTheSignature(t *testing.T) {
	t.Parallel()

	confirmer := &recordingConfirmer{}
	d := deps("p1", &fakeLedger{})
	d.Confirmer = confirmer
	_, err := operations.ExecuteSequence(operations.NewBundle(t.Context, logger.Nop(), operations.NewMemoryReporter()),
		exercisechoice.ExerciseChoiceSequence, d, input())
	require.ErrorContains(t, err, exercisechoice.ErrThresholdNotMet.Error())

	require.Len(t, confirmer.got, 1)
	assert.Equal(t, "Utility.Registry.App.V0.Service.Provider:ProviderService", confirmer.got[0].Template)
	assert.Equal(t, "provider-service-cid", confirmer.got[0].ContractID)
	assert.Equal(t, "ProviderService_AcceptRegistrarServiceRequest", confirmer.got[0].Choice)
	assert.Equal(t, input().ChoiceArgument, confirmer.got[0].ChoiceArgument)
}
