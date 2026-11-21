package client

import (
	"context"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	apiv2admin "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type fakeInteractive struct {
	interactive.InteractiveSubmissionServiceClient
	events   []*apiv2.Event
	prepared *interactive.PrepareSubmissionRequest
}

func (f *fakeInteractive) PrepareSubmission(
	_ context.Context, req *interactive.PrepareSubmissionRequest, _ ...grpc.CallOption,
) (*interactive.PrepareSubmissionResponse, error) {
	f.prepared = req
	return &interactive.PrepareSubmissionResponse{}, nil
}

type fakeUserManagement struct {
	apiv2admin.UserManagementServiceClient
	granted *apiv2admin.GrantUserRightsRequest
}

func (f *fakeUserManagement) GrantUserRights(
	_ context.Context, req *apiv2admin.GrantUserRightsRequest, _ ...grpc.CallOption,
) (*apiv2admin.GrantUserRightsResponse, error) {
	f.granted = req
	return &apiv2admin.GrantUserRightsResponse{}, nil
}

func (f *fakeInteractive) ExecuteSubmissionAndWaitForTransaction(
	context.Context, *interactive.ExecuteSubmissionAndWaitForTransactionRequest, ...grpc.CallOption,
) (*interactive.ExecuteSubmissionAndWaitForTransactionResponse, error) {
	return &interactive.ExecuteSubmissionAndWaitForTransactionResponse{
		Transaction: &apiv2.Transaction{UpdateId: "update-1", Events: f.events},
	}, nil
}

func TestExecuteSubmission_LeavesOutContractsArchivedInTheSameTransaction(t *testing.T) {
	t.Parallel()

	tid := &apiv2.Identifier{PackageId: "pkg", ModuleName: "M", EntityName: "E"}
	c := &GRPCLedgerClient{interactive: &fakeInteractive{events: []*apiv2.Event{
		{Event: &apiv2.Event_Created{Created: &apiv2.CreatedEvent{ContractId: "transient-cid", TemplateId: tid}}},
		{Event: &apiv2.Event_Archived{Archived: &apiv2.ArchivedEvent{ContractId: "transient-cid", TemplateId: tid}}},
		{Event: &apiv2.Event_Created{Created: &apiv2.CreatedEvent{ContractId: "active-cid", TemplateId: tid}}},
	}}}

	result, err := c.ExecuteSubmission(t.Context(), &interactive.PreparedTransaction{}, &interactive.PartySignatures{}, 0)

	require.NoError(t, err)
	require.Equal(t, []CreatedContract{{TemplateID: "pkg:M:E", ContractID: "active-cid"}}, result.Created)
}

func TestPrepareSubmission_SendsTheDisclosedContracts(t *testing.T) {
	t.Parallel()

	disclosed := []*apiv2.DisclosedContract{{
		TemplateId:       &apiv2.Identifier{PackageId: "pkg", ModuleName: "M", EntityName: "E"},
		ContractId:       "disclosed-cid",
		CreatedEventBlob: []byte("blob"),
		SynchronizerId:   "global",
	}}
	fi := &fakeInteractive{}
	c := &GRPCLedgerClient{interactive: fi}

	_, err := c.PrepareSubmission(t.Context(), nil, []string{"party"}, nil, "global", disclosed)

	require.NoError(t, err)
	require.Len(t, fi.prepared.GetDisclosedContracts(), 1)
	require.True(t, proto.Equal(disclosed[0], fi.prepared.GetDisclosedContracts()[0]))
}

func TestGrantRights_SendsOnlyTheRequestedRight(t *testing.T) {
	t.Parallel()

	actAs := &apiv2admin.Right{Kind: &apiv2admin.Right_CanActAs_{CanActAs: &apiv2admin.Right_CanActAs{Party: "party"}}}
	readAs := &apiv2admin.Right{Kind: &apiv2admin.Right_CanReadAs_{CanReadAs: &apiv2admin.Right_CanReadAs{Party: "party"}}}
	tests := []struct {
		name  string
		grant func(*GRPCLedgerClient) error
		want  []*apiv2admin.Right
	}{
		{"actAs", func(c *GRPCLedgerClient) error { return c.GrantActAs(t.Context(), "user", "party") }, []*apiv2admin.Right{actAs}},
		{"readAs", func(c *GRPCLedgerClient) error { return c.GrantReadAs(t.Context(), "user", "party") }, []*apiv2admin.Right{readAs}},
		{"both", func(c *GRPCLedgerClient) error { return c.GrantPartyRights(t.Context(), "user", "party") }, []*apiv2admin.Right{actAs, readAs}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			um := &fakeUserManagement{}

			require.NoError(t, tc.grant(&GRPCLedgerClient{userManagement: um}))

			require.Equal(t, "user", um.granted.GetUserId())
			require.Len(t, um.granted.GetRights(), len(tc.want))
			for i, want := range tc.want {
				require.True(t, proto.Equal(want, um.granted.GetRights()[i]), "right %d: got %v", i, um.granted.GetRights()[i])
			}
		})
	}
}
