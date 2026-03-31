package store

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// fakeActiveContractsStream implements grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse] for tests.
type fakeActiveContractsStream struct {
	ctx        context.Context //nolint:containedctx
	responses  []*apiv2.GetActiveContractsResponse
	err        error
	idx        int
	blockOnEOF bool // when true, Recv blocks on ctx.Done() instead of returning EOF (for context cancellation tests)
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

// fakeUpdatesStream implements grpc.ServerStreamingClient[apiv2.GetUpdatesResponse] for tests.
type fakeUpdatesStream struct {
	ctx        context.Context //nolint:containedctx
	responses  []*apiv2.GetUpdatesResponse
	err        error
	idx        int
	blockOnEOF bool
}

func (s *fakeUpdatesStream) Recv() (*apiv2.GetUpdatesResponse, error) {
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

func (s *fakeUpdatesStream) Header() (metadata.MD, error) { return metadata.MD{}, nil }
func (s *fakeUpdatesStream) Trailer() metadata.MD         { return metadata.MD{} }
func (s *fakeUpdatesStream) CloseSend() error             { return nil }
func (s *fakeUpdatesStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}

	return context.Background()
}
func (s *fakeUpdatesStream) SendMsg(any) error { return nil }
func (s *fakeUpdatesStream) RecvMsg(any) error { return nil }

var _ grpc.ServerStreamingClient[apiv2.GetUpdatesResponse] = (*fakeUpdatesStream)(nil)

// stubEmptyActiveContractsBackfill makes GetActiveContracts return an empty stream so backfill completes
// without blocking. Required for any Run test that passes GetLedgerEnd: unexpected GetActiveContracts
// from a non-test goroutine triggers mock fail → t.FailNow on wrong goroutine → deadlock on <-done.
func stubEmptyActiveContractsBackfill(stateClient *mocks.MockStateServiceClient) {
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything).
		Return(&fakeActiveContractsStream{}, nil).Once()
}

// holdingCreatedEvent returns a CreatedEvent that UnmarshalCreatedEvent[HoldingView] can parse.
// owner and instrumentId must match the store's owner and the instrument ID you want to look up.
func holdingCreatedEvent(contractID string, owner types.PARTY, instrumentID splice_api_token_holding_v1.InstrumentId) *apiv2.CreatedEvent {
	templateID := &apiv2.Identifier{
		PackageId:  "#" + splice_api_token_holding_v1.PackageName,
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}

	view := &apiv2.Record{
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
	}

	return &apiv2.CreatedEvent{
		ContractId:      contractID,
		TemplateId:      templateID,
		CreateArguments: view,
		InterfaceViews: []*apiv2.InterfaceView{
			{
				// getRelevantInterfaceViewValue matches module/entity only (ignores package id).
				InterfaceId: &apiv2.Identifier{ModuleName: templateID.ModuleName, EntityName: templateID.EntityName},
				ViewValue:   view,
			},
		},
	}
}

func TestInstrumentHoldingStoreService_Run(t *testing.T) {
	t.Parallel()
	const testOwner = "test-owner"
	owner := types.PARTY(testOwner)
	logger := zerolog.Nop()

	makeConfig := func(stateClient *mocks.MockStateServiceClient, updateClient *mocks.MockUpdateServiceClient, maxRetries int) InstrumentHoldingStoreConfig {
		return InstrumentHoldingStoreConfig{
			Logger:           logger,
			Owner:            owner,
			StateService:     stateClient,
			UpdateService:    updateClient,
			MaxRetries:       maxRetries,
			ReconnectBackoff: time.Millisecond, // short backoff in tests so reconnect tests don't wait
		}
	}

	t.Run("records holding when stream yields transaction with Created then reconnects on EOF", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		instrumentID := splice_api_token_holding_v1.InstrumentId{Admin: "admin-party", Id: "instrument-1"}
		contractID := "holding-contract-1"
		createdEvent := holdingCreatedEvent(contractID, owner, instrumentID)

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 5}, nil)
		stubEmptyActiveContractsBackfill(stateClient)

		tx := &apiv2.Transaction{
			UpdateId: "test-update-id",
			Offset:   6,
			Events: []*apiv2.Event{
				{Event: &apiv2.Event_Created{Created: createdEvent}},
			},
		}
		updateClient := mocks.NewMockUpdateServiceClient(t)
		// First call: one transaction then EOF; Run will reconnect. Second call: block until ctx done.
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.MatchedBy(func(req *apiv2.GetUpdatesRequest) bool {
			return req != nil && req.BeginExclusive == 5 &&
				req.UpdateFormat != nil &&
				req.UpdateFormat.IncludeTransactions != nil &&
				req.UpdateFormat.IncludeTransactions.EventFormat != nil &&
				req.UpdateFormat.IncludeTransactions.EventFormat.FiltersByParty[testOwner] != nil
		}), mock.Anything).
			Return(&fakeUpdatesStream{
				ctx: ctx,
				responses: []*apiv2.GetUpdatesResponse{
					{Update: &apiv2.GetUpdatesResponse_Transaction{Transaction: tx}},
				},
			}, nil)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.MatchedBy(func(req *apiv2.GetUpdatesRequest) bool {
			return req != nil && req.BeginExclusive == 6
		}), mock.Anything).
			Return(&fakeUpdatesStream{ctx: ctx, blockOnEOF: true}, nil)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0))
		done := make(chan error, 1)
		go func() { done <- svc.Run(ctx) }()

		// Wait for disclosure to be recorded (from first stream)
		require.Eventually(t, func() bool {
			_, err := svc.GetInstrumentHolding(ctx, instrumentID)
			return err == nil
		}, 2*time.Second, 10*time.Millisecond)

		disclosure, err := svc.GetInstrumentHolding(ctx, instrumentID)
		require.NoError(t, err)
		require.NotNil(t, disclosure)
		require.Equal(t, contractID, disclosure.ContractId)
		require.Equal(t, createdEvent.TemplateId, disclosure.TemplateId)

		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
	})

	t.Run("surfaces GetLedgerEnd error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		expectedErr := errors.New("ledger end failed")
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return((*apiv2.GetLedgerEndResponse)(nil), expectedErr)
		updateClient := mocks.NewMockUpdateServiceClient(t)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0))
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get ledger end")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("surfaces GetUpdates error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		expectedErr := errors.New("get updates failed")
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		stubEmptyActiveContractsBackfill(stateClient)
		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return((grpc.ServerStreamingClient[apiv2.GetUpdatesResponse])(nil), expectedErr).Times(2)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 1))
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to create update stream")
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("reconnects on stream Recv error", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		stubEmptyActiveContractsBackfill(stateClient)
		updateClient := mocks.NewMockUpdateServiceClient(t)
		// First stream returns recvErr; Run reconnects. Second stream blocks until ctx done.
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeUpdatesStream{ctx: ctx, err: errors.New("recv failed")}, nil)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeUpdatesStream{ctx: ctx, blockOnEOF: true}, nil)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0))
		done := make(chan error, 1)
		go func() { done <- svc.Run(ctx) }()
		// Let Run pass backfill and the first (failing) stream before cancelling, otherwise ctx may
		// still be canceled during backfill and GetUpdates is never called.
		time.Sleep(50 * time.Millisecond)
		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
	})

	t.Run("returns context error when context is cancelled", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&apiv2.GetLedgerEndResponse{Offset: 0}, nil)
		stubEmptyActiveContractsBackfill(stateClient)
		// Cancel before Run: backfill observes ctx.Done() and returns before GetUpdates.
		updateClient := mocks.NewMockUpdateServiceClient(t)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 1))
		cancel()
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
	})
}
