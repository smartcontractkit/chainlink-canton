package store

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-canton/eds/monitoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fakeUpdateService fakes apiv2.UpdateServiceClient; only GetUpdates is wired.
type fakeUpdateService struct {
	getUpdates func(ctx context.Context, in *apiv2.GetUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error)
}

func (f *fakeUpdateService) GetUpdates(ctx context.Context, in *apiv2.GetUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
	return f.getUpdates(ctx, in, opts...)
}
func (f *fakeUpdateService) GetUpdateByOffset(ctx context.Context, in *apiv2.GetUpdateByOffsetRequest, opts ...grpc.CallOption) (*apiv2.GetUpdateResponse, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeUpdateService) GetUpdateById(ctx context.Context, in *apiv2.GetUpdateByIdRequest, opts ...grpc.CallOption) (*apiv2.GetUpdateResponse, error) {
	return nil, errors.New("not implemented")
}

// fakeStream fakes grpc.ServerStreamingClient with a canned Recv.
type fakeStream struct {
	grpc.ClientStream // embedded to satisfy the interface; unused by the stream code
	recv              func() (*apiv2.GetUpdatesResponse, error)
}

func (f *fakeStream) Recv() (*apiv2.GetUpdatesResponse, error) { return f.recv() }
func (f *fakeStream) CloseSend() error                         { return nil }

// A stream that immediately fails like the production 404 (Unimplemented on Recv)
// must fail closed after MaxEmptyReconnects incarnations instead of looping forever.
func TestStreamFailsClosedAfterEmptyReconnects(t *testing.T) {
	var calls atomic.Int32
	svc := &fakeUpdateService{getUpdates: func(ctx context.Context, in *apiv2.GetUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
		calls.Add(1)
		return &fakeStream{recv: func() (*apiv2.GetUpdatesResponse, error) {
			return nil, status.Error(codes.Unimplemented, `unexpected HTTP status code received from server: 404 (Not Found)`)
		}}, nil
	}}
	s := &BackfilledStream{logger: zerolog.Nop(), metrics: monitoring.NoopEDSMetricLabeler{}, updateService: svc}

	cfg := DefaultStreamConfig()
	cfg.MaxEmptyReconnects = 3
	cfg.MaxDisconnected = 0 // isolate the empty-reconnect counter

	err := s.stream(context.Background(), 0, FiltersByParty{}, cfg, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "without receiving any update")
	assert.Equal(t, int32(3), calls.Load(), "must stop after exactly MaxEmptyReconnects reconnects")
}

// Receiving updates resets the budget: ten incarnations that each deliver one
// update must never trip it, even though MaxEmptyReconnects is 3.
func TestStreamBudgetResetsOnUpdates(t *testing.T) {
	var served atomic.Int32
	svc := &fakeUpdateService{getUpdates: func(ctx context.Context, in *apiv2.GetUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
		if served.Load() >= 10 {
			// Participant stops accepting streams: the connect retry loop then
			// runs until the test context expires, so the served count is exact.
			return nil, status.Error(codes.Unavailable, "participant down")
		}
		served.Add(1)
		first := true
		return &fakeStream{recv: func() (*apiv2.GetUpdatesResponse, error) {
			if first {
				first = false
				return &apiv2.GetUpdatesResponse{Update: &apiv2.GetUpdatesResponse_Transaction{Transaction: &apiv2.Transaction{Offset: 1}}}, nil
			}
			return nil, io.EOF
		}}, nil
	}}
	s := &BackfilledStream{logger: zerolog.Nop(), metrics: monitoring.NoopEDSMetricLabeler{}, updateService: svc}

	cfg := DefaultStreamConfig()
	cfg.MaxEmptyReconnects = 3
	cfg.MaxDisconnected = 0 // keep the connect retry loop running until the context expires

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := s.stream(ctx, 0, FiltersByParty{}, cfg, nil, nil, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded) // stopped by the test, not by the budget
	assert.Equal(t, int32(10), served.Load(), "budget must reset on every update")
}

// A participant that refuses all connections must fail closed after
// MaxDisconnected of retrying, not retry forever.
func TestStreamFailsClosedAfterDisconnected(t *testing.T) {
	svc := &fakeUpdateService{getUpdates: func(ctx context.Context, in *apiv2.GetUpdatesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
		return nil, status.Error(codes.Unavailable, "connection refused")
	}}
	s := &BackfilledStream{logger: zerolog.Nop(), metrics: monitoring.NoopEDSMetricLabeler{}, updateService: svc}

	cfg := DefaultStreamConfig()
	cfg.MaxEmptyReconnects = 0 // isolate the disconnect timer
	cfg.MaxDisconnected = 50 * time.Millisecond
	cfg.BackoffMin = time.Millisecond
	cfg.BackoffMax = 5 * time.Millisecond

	start := time.Now()
	err := s.stream(context.Background(), 0, FiltersByParty{}, cfg, nil, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disconnected for")
	assert.LessOrEqual(t, time.Since(start), 5*time.Second, "must not retry forever")
}
