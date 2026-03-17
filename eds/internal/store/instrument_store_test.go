package store

import (
	"context"
	"errors"
	"io"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// fakeActiveContractsStream implements grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse] for tests.
type fakeActiveContractsStream struct {
	ctx         context.Context //nolint:containedctx
	responses   []*apiv2.GetActiveContractsResponse
	err         error
	idx         int
	blockOnEOF  bool // when true, Recv blocks on ctx.Done() instead of returning EOF (for context cancellation tests)
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
	if s.blockOnEOF && s.ctx != nil {
		<-s.ctx.Done()
		return nil, s.ctx.Err()
	}
	return nil, io.EOF
}

func (s *fakeActiveContractsStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s *fakeActiveContractsStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s *fakeActiveContractsStream) CloseSend() error {
	return nil
}

func (s *fakeActiveContractsStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}

func (s *fakeActiveContractsStream) SendMsg(any) error {
	return nil
}

func (s *fakeActiveContractsStream) RecvMsg(any) error {
	return nil
}

var _ grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse] = (*fakeActiveContractsStream)(nil)

// holdingCreatedEvent returns a CreatedEvent that UnmarshalCreatedEvent[HoldingView] can parse.
// owner and instrumentId must match the store's owner and the instrument ID you want to look up.
func holdingCreatedEvent(contractID string, owner types.PARTY, instrumentID splice_api_token_holding_v1.InstrumentId) *apiv2.CreatedEvent {
	templateID := &apiv2.Identifier{
		PackageId:  "#" + splice_api_token_holding_v1.PackageName,
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}
	return &apiv2.CreatedEvent{
		ContractId: contractID,
		TemplateId: templateID,
		CreateArguments: &apiv2.Record{
			Fields: []*apiv2.RecordField{
				{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(owner)}}},
				{Label: "instrumentId", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: string(instrumentID.Admin)}}},
					{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: string(instrumentID.Id)}}},
				}}}}},
				{Label: "amount", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
				{Label: "lock", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{}}}},
				{Label: "meta", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
				}}}}},
			},
		},
	}
}

func TestInstrumentHoldingStoreService_Run(t *testing.T) {
	t.Parallel()
	const testOwner = "test-owner"
	owner := types.PARTY(testOwner)
	logger := zerolog.Nop()

	t.Run("records holding when stream yields ActiveContract then EOF", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		instrumentID := splice_api_token_holding_v1.InstrumentId{Admin: "admin-party", Id: "instrument-1"}
		contractID := "holding-contract-1"
		createdEvent := holdingCreatedEvent(contractID, owner, instrumentID)

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.MatchedBy(func(req *apiv2.GetActiveContractsRequest) bool {
			return req != nil && req.ActiveAtOffset == 5 &&
				req.EventFormat != nil &&
				req.EventFormat.FiltersByParty[testOwner] != nil
		}), mock.Anything).
			Return(&fakeActiveContractsStream{
				ctx: ctx,
				responses: []*apiv2.GetActiveContractsResponse{
					{
						ContractEntry: &apiv2.GetActiveContractsResponse_ActiveContract{
							ActiveContract: &apiv2.ActiveContract{
								CreatedEvent:   createdEvent,
								SynchronizerId: "sync-1",
							},
						},
					},
				},
			}, nil)

		svc := NewInstrumentHoldingStore(owner, stateClient, logger)
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, io.EOF)
		require.Contains(t, err.Error(), "failed to receive active contract")

		// Run processes one contract then exits on EOF; disclosure should be recorded
		disclosure, err := svc.GetInstrumentHolding(ctx, instrumentID)
		require.NoError(t, err)
		require.NotNil(t, disclosure)
		require.Equal(t, contractID, disclosure.ContractId)
		require.Equal(t, createdEvent.TemplateId, disclosure.TemplateId)
	})

	t.Run("surfaces GetLedgerEnd error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		expectedErr := errors.New("ledger end failed")
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return((*apiv2.GetLedgerEndResponse)(nil), expectedErr)

		svc := NewInstrumentHoldingStore(owner, stateClient, logger)
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get ledger end")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("surfaces GetActiveContracts error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		expectedErr := errors.New("get active contracts failed")
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		// GetStreamWithRetry retries; with MaxRetries=1 we get 2 attempts total
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return((grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse])(nil), expectedErr).Times(2)

		svc := NewInstrumentHoldingStore(owner, stateClient, logger)
		svc.MaxRetries = 1
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get active contracts")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("surfaces stream Recv error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		recvErr := errors.New("recv failed")
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx, err: recvErr}, nil)

		svc := NewInstrumentHoldingStore(owner, stateClient, logger)
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to receive active contract")
		require.ErrorIs(t, err, recvErr)
	})

	t.Run("returns context error when context is cancelled", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		// Stream blocks on Recv until ctx is done, then returns context.Canceled
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx, blockOnEOF: true}, nil)

		svc := NewInstrumentHoldingStore(owner, stateClient, logger)
		svc.MaxRetries = 1 // no retries so we connect once and then block on Recv
		cancel()
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
	})
}
