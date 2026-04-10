package store

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/eds/monitoring"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// stubEmptyActiveContractsBackfill makes GetActiveContracts return an empty stream so backfill completes
// without blocking. Required for any Run test that passes GetLedgerEnd: unexpected GetActiveContracts
// from a non-test goroutine triggers mock fail → t.FailNow on wrong goroutine → deadlock on <-done.
func stubEmptyActiveContractsBackfill(t *testing.T, stateClient *mocks.MockStateServiceClient) {
	activeContractStream := mocks.NewMockServerStreamingClient[apiv2.GetActiveContractsResponse](t)
	activeContractStream.EXPECT().Recv().Return(nil, io.EOF).Maybe()
	stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything).
		Return(activeContractStream, nil).Once()
	activeContractStream.EXPECT().CloseSend().Return(nil).Maybe()
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

// afterCancel waits for the context to be Done and then send the current time on the returned channel.
func afterCancel(ctx context.Context) <-chan time.Time {
	ch := make(chan time.Time, 1)
	go func() {
		<-ctx.Done()
		ch <- time.Now()
	}()

	return ch
}

func TestInstrumentHoldingStoreService_Run(t *testing.T) {
	t.Parallel()
	const testOwner = "test-owner"
	owner := types.PARTY(testOwner)
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.TraceLevel)

	makeConfig := func(stateClient *mocks.MockStateServiceClient, updateClient *mocks.MockUpdateServiceClient, maxRetries int) InstrumentHoldingStoreConfig {
		return InstrumentHoldingStoreConfig{
			Logger:        logger,
			Owner:         owner,
			StateService:  stateClient,
			UpdateService: updateClient,
			StreamConfig: ReliableStreamConfig{
				MaxRetries: maxRetries,
			},
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

		stubEmptyActiveContractsBackfill(t, stateClient)

		updatesStream := mocks.NewMockServerStreamingClient[apiv2.GetUpdatesResponse](t)
		updateClient := mocks.NewMockUpdateServiceClient(t)

		tx := &apiv2.Transaction{
			UpdateId: "test-update-id",
			Offset:   6,
			Events: []*apiv2.Event{
				{Event: &apiv2.Event_Created{Created: createdEvent}},
			},
		}
		// First subscription: one transaction then EOF; Run will reconnect
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.MatchedBy(func(req *apiv2.GetUpdatesRequest) bool {
			return req != nil && req.BeginExclusive == 5 && // Begin must match offset returned by GetLedgerEnd
				req.UpdateFormat != nil &&
				req.UpdateFormat.IncludeTransactions != nil &&
				req.UpdateFormat.IncludeTransactions.EventFormat != nil &&
				req.UpdateFormat.IncludeTransactions.EventFormat.FiltersByParty[testOwner] != nil
		}), mock.Anything).
			Return(updatesStream, nil).Once()
		updatesStream.EXPECT().Recv().Return(
			&apiv2.GetUpdatesResponse{
				Update: &apiv2.GetUpdatesResponse_Transaction{Transaction: tx},
			},
			nil,
		).Once()
		updatesStream.EXPECT().Recv().Return(nil, io.EOF).Once() // Send EOF, store will reconnect

		// Second subscription: block until ctx done; Run will call CloseSend during teardown
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.MatchedBy(func(req *apiv2.GetUpdatesRequest) bool {
			return req != nil && req.BeginExclusive == 6
		}), mock.Anything).
			Return(updatesStream, nil).Once()
		updatesStream.EXPECT().Recv().WaitUntil(afterCancel(ctx)).Return(nil, nil).Once()
		updatesStream.EXPECT().CloseSend().Return(nil).Once()

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0), monitoring.NoopEDSMetricLabeler{})
		done := make(chan error, 1)
		go func() { done <- svc.Run(ctx) }()

		// Wait for disclosure to be recorded (from first stream)
		require.Eventually(t, func() bool {
			_, ok := svc.Get(instrumentID)
			return ok
		}, 2*time.Second, 10*time.Millisecond)

		disclosure, ok := svc.Get(instrumentID)
		require.True(t, ok)
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

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0), monitoring.NoopEDSMetricLabeler{})
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
		stubEmptyActiveContractsBackfill(t, stateClient)
		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return((grpc.ServerStreamingClient[apiv2.GetUpdatesResponse])(nil), expectedErr).Times(2)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 1), monitoring.NoopEDSMetricLabeler{})
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
		stubEmptyActiveContractsBackfill(t, stateClient)
		updatesStream := mocks.NewMockServerStreamingClient[apiv2.GetUpdatesResponse](t)
		updateClient := mocks.NewMockUpdateServiceClient(t)
		// First stream returns recvErr; Run reconnects. Second stream blocks until ctx done.
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return(updatesStream, nil)
		updatesStream.EXPECT().Recv().Return(nil, errors.New("recv failed")).Once()
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything, mock.Anything).
			Return(updatesStream, nil)
		updatesStream.EXPECT().Recv().WaitUntil(afterCancel(ctx)).Return(nil, nil).Once()
		updatesStream.EXPECT().CloseSend().Return(nil).Once()

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 0), monitoring.NoopEDSMetricLabeler{})
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
		stubEmptyActiveContractsBackfill(t, stateClient)
		// Cancel before Run: backfill observes ctx.Done() and returns before GetUpdates.
		updateClient := mocks.NewMockUpdateServiceClient(t)

		svc := NewInstrumentHoldingStore(makeConfig(stateClient, updateClient, 1), monitoring.NoopEDSMetricLabeler{})
		cancel()
		err := svc.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, context.Canceled)
	})
}
