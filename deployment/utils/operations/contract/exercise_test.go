package contract

import (
	"context"
	"errors"
	"io"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

func TestLedgerQueryParties(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		p        canton.Participant
		want     []string
	}{
		{
			name: "actas only",
			p: canton.Participant{
				PartyID: "operator",
			},
			want: []string{"operator"},
		},
		{
			name: "actas plus readas",
			p: canton.Participant{
				PartyID:         "operator",
				ReadAsPartyIDs:  []string{"ccip-owner", "pool-owner"},
			},
			want: []string{"operator", "ccip-owner", "pool-owner"},
		},
		{
			name: "dedupes actas when also in readas",
			p: canton.Participant{
				PartyID:        "operator",
				ReadAsPartyIDs: []string{"operator", "ccip-owner"},
			},
			want: []string{"operator", "ccip-owner"},
		},
		{
			name: "skips empty strings",
			p: canton.Participant{
				PartyID:        "operator",
				ReadAsPartyIDs: []string{"", "ccip-owner", ""},
			},
			want: []string{"operator", "ccip-owner"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, LedgerQueryParties(tt.p))
		})
	}
}

func TestFiltersByPartyForTemplate(t *testing.T) {
	t.Parallel()

	parties := []string{"party-a", "party-b", "party-c"}
	got := filtersByPartyForTemplate(parties, "pkg-id", "CCIP.GlobalConfig", "GlobalConfig")

	require.Len(t, got, 3)
	for _, party := range parties {
		f, ok := got[party]
		require.True(t, ok, "missing filter for %s", party)
		require.Len(t, f.Cumulative, 1)
		tplFilter, ok := f.Cumulative[0].IdentifierFilter.(*apiv2.CumulativeFilter_TemplateFilter)
		require.True(t, ok)
		tf := tplFilter.TemplateFilter
		require.Equal(t, "pkg-id", tf.TemplateId.PackageId)
		require.Equal(t, "CCIP.GlobalConfig", tf.TemplateId.ModuleName)
		require.Equal(t, "GlobalConfig", tf.TemplateId.EntityName)
		require.True(t, tf.IncludeCreatedEventBlob)
	}
}

func TestFindActiveContractByInstanceAddress_requiresParties(t *testing.T) {
	t.Parallel()

	_, err := FindActiveContractByInstanceAddress(
		t.Context(),
		mocks.NewMockStateServiceClient(t),
		nil,
		common.GlobalConfig{}.GetTemplateID(),
		contracts.InstanceAddress{},
	)
	require.ErrorContains(t, err, "at least one query party is required")
}

func TestFindActiveContractByInstanceAddress_queriesAllParties(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	const (
		instanceID = "globalconfig-abc"
		signatory  = "ccip-owner"
		contractID = "cid-001"
	)
	target := contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(signatory)).InstanceAddress()
	templateID := common.GlobalConfig{}.GetTemplateID()

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 42}, nil)
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.MatchedBy(func(req *apiv2.GetActiveContractsRequest) bool {
		if req.ActiveAtOffset != 42 || req.EventFormat == nil {
			return false
		}
		fbp := req.EventFormat.FiltersByParty
		return len(fbp) == 3 &&
			fbp["operator"] != nil &&
			fbp["ccip-owner"] != nil &&
			fbp["pool-owner"] != nil
	}), mock.Anything).Return(newFakeActiveContractsStream(ctx, []*apiv2.GetActiveContractsResponse{
		makeActiveContractACSResponse(instanceID, signatory, contractID),
	}), nil)

	got, err := FindActiveContractByInstanceAddress(
		ctx,
		stateClient,
		[]string{"operator", "ccip-owner", "pool-owner"},
		templateID,
		target,
	)
	require.NoError(t, err)
	require.Equal(t, contractID, got.GetCreatedEvent().GetContractId())
}

func TestFindActiveContractByInstanceAddress_dedupesSameContractAcrossParties(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	const (
		instanceID = "globalconfig-dedup"
		signatory  = "ccip-owner"
		contractID = "cid-dedup"
	)
	target := contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(signatory)).InstanceAddress()
	templateID := common.GlobalConfig{}.GetTemplateID()
	dup := makeActiveContractACSResponse(instanceID, signatory, contractID)

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil)
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		Return(newFakeActiveContractsStream(ctx, []*apiv2.GetActiveContractsResponse{dup, dup}), nil)

	got, err := FindActiveContractByInstanceAddress(ctx, stateClient, []string{"a", "b"}, templateID, target)
	require.NoError(t, err)
	require.Equal(t, contractID, got.GetCreatedEvent().GetContractId())
}

func TestFindActiveContractByInstanceAddress_notFound(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	signatory := types.PARTY("ccip-owner")
	target := contracts.InstanceID("missing").RawInstanceAddress(signatory).InstanceAddress()

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil)
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		Return(newFakeActiveContractsStream(ctx, nil), nil)

	_, err := FindActiveContractByInstanceAddress(
		ctx,
		stateClient,
		[]string{"operator"},
		common.GlobalConfig{}.GetTemplateID(),
		target,
	)
	require.ErrorContains(t, err, "no active contract found")
}

func TestFindActiveContractByInstanceAddress_multipleMatches(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	const signatory = "ccip-owner"
	target := contracts.InstanceID("gc-multi").RawInstanceAddress(types.PARTY(signatory)).InstanceAddress()
	templateID := common.GlobalConfig{}.GetTemplateID()

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 1}, nil)
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		Return(newFakeActiveContractsStream(ctx, []*apiv2.GetActiveContractsResponse{
			makeActiveContractACSResponse("gc-multi", signatory, "cid-1"),
			makeActiveContractACSResponse("gc-multi", signatory, "cid-2"),
		}), nil)

	_, err := FindActiveContractByInstanceAddress(ctx, stateClient, []string{"operator"}, templateID, target)
	require.ErrorContains(t, err, "multiple active contracts found")
}

func TestFindActiveContractIDByInstanceAddress(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	const (
		instanceID = "globalconfig-id"
		signatory  = "ccip-owner"
		contractID = "cid-wrap"
	)
	target := contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(signatory)).InstanceAddress()

	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(&apiv2.GetLedgerEndResponse{Offset: 7}, nil)
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
		Return(newFakeActiveContractsStream(ctx, []*apiv2.GetActiveContractsResponse{
			makeActiveContractACSResponse(instanceID, signatory, contractID),
		}), nil)

	got, err := FindActiveContractIDByInstanceAddress(
		ctx,
		stateClient,
		LedgerQueryParties(canton.Participant{PartyID: "operator", ReadAsPartyIDs: []string{signatory}}),
		common.GlobalConfig{}.GetTemplateID(),
		target,
	)
	require.NoError(t, err)
	require.Equal(t, contractID, got)
}

func makeActiveContractACSResponse(instanceID, signatory, contractID string) *apiv2.GetActiveContractsResponse {
	return &apiv2.GetActiveContractsResponse{
		ContractEntry: &apiv2.GetActiveContractsResponse_ActiveContract{
			ActiveContract: &apiv2.ActiveContract{
				CreatedEvent: &apiv2.CreatedEvent{
					ContractId: contractID,
					CreateArguments: &apiv2.Record{
						Fields: []*apiv2.RecordField{
							{
								Label: "instanceId",
								Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}},
							},
						},
					},
					Signatories: []string{signatory},
				},
			},
		},
	}
}

type fakeActiveContractsStream struct {
	ctx       context.Context
	responses []*apiv2.GetActiveContractsResponse
	err       error
	idx       int
}

func newFakeActiveContractsStream(ctx context.Context, responses []*apiv2.GetActiveContractsResponse) *fakeActiveContractsStream {
	return &fakeActiveContractsStream{ctx: ctx, responses: responses}
}

func (s *fakeActiveContractsStream) Recv() (*apiv2.GetActiveContractsResponse, error) {
	if s.idx < len(s.responses) {
		resp := s.responses[s.idx]
		s.idx++
		return resp, nil
	}
	if s.err != nil {
		return nil, s.err
	}
	return nil, io.EOF
}

func (s *fakeActiveContractsStream) Header() (metadata.MD, error)  { return metadata.MD{}, nil }
func (s *fakeActiveContractsStream) Trailer() metadata.MD           { return metadata.MD{} }
func (s *fakeActiveContractsStream) CloseSend() error               { return nil }
func (s *fakeActiveContractsStream) Context() context.Context       { return s.ctx }
func (s *fakeActiveContractsStream) SendMsg(any) error              { return nil }
func (s *fakeActiveContractsStream) RecvMsg(any) error              { return nil }

var _ grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse] = (*fakeActiveContractsStream)(nil)

func TestFindActiveContractByInstanceAddress_surfacesLedgerErrors(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	stateClient := mocks.NewMockStateServiceClient(t)
	stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
		Return(nil, errors.New("ledger unavailable"))

	_, err := FindActiveContractByInstanceAddress(
		ctx,
		stateClient,
		[]string{"operator"},
		common.GlobalConfig{}.GetTemplateID(),
		contracts.InstanceAddress{},
	)
	require.ErrorContains(t, err, "failed to get ledger end")
}
