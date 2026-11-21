package client

import (
	"context"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type fakeInteractive struct {
	interactive.InteractiveSubmissionServiceClient
	events []*apiv2.Event
}

func (f fakeInteractive) ExecuteSubmissionAndWaitForTransaction(
	context.Context, *interactive.ExecuteSubmissionAndWaitForTransactionRequest, ...grpc.CallOption,
) (*interactive.ExecuteSubmissionAndWaitForTransactionResponse, error) {
	return &interactive.ExecuteSubmissionAndWaitForTransactionResponse{
		Transaction: &apiv2.Transaction{UpdateId: "update-1", Events: f.events},
	}, nil
}

func TestExecuteSubmission_LeavesOutContractsArchivedInTheSameTransaction(t *testing.T) {
	t.Parallel()

	tid := &apiv2.Identifier{PackageId: "pkg", ModuleName: "M", EntityName: "E"}
	c := &GRPCLedgerClient{interactive: fakeInteractive{events: []*apiv2.Event{
		{Event: &apiv2.Event_Created{Created: &apiv2.CreatedEvent{ContractId: "transient-cid", TemplateId: tid}}},
		{Event: &apiv2.Event_Archived{Archived: &apiv2.ArchivedEvent{ContractId: "transient-cid", TemplateId: tid}}},
		{Event: &apiv2.Event_Created{Created: &apiv2.CreatedEvent{ContractId: "active-cid", TemplateId: tid}}},
	}}}

	result, err := c.ExecuteSubmission(t.Context(), &interactive.PreparedTransaction{}, &interactive.PartySignatures{}, 0)

	require.NoError(t, err)
	require.Equal(t, []CreatedContract{{TemplateID: "pkg:M:E", ContractID: "active-cid"}}, result.Created)
}
