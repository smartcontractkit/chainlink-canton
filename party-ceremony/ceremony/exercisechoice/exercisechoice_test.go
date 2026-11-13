package exercisechoice_test

import (
	"context"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

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
	prepared   [][]*apiv2.Command
	executions int
}

func (f *fakeLedger) PartyExists(context.Context, string) (bool, error) { return true, nil }

func (f *fakeLedger) PrepareSubmission(_ context.Context, commands []*apiv2.Command, _, _ []string, _ string) (*interactive.PrepareSubmissionResponse, error) {
	f.prepared = append(f.prepared, commands)

	return &interactive.PrepareSubmissionResponse{PreparedTransactionHash: []byte{0xab, 0xcd}}, nil
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

type fakeSigner struct{}

func (fakeSigner) Sign(_ context.Context, hash []byte) (*apiv2.Signature, error) {
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
