package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jpillora/backoff"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// StreamFactory is a function that creates a new stream starting at the given offset.
// Both GetUpdates (BeginExclusive) and GetActiveContracts (ActiveAtOffset) use an offset.
type StreamFactory[T any] func(ctx context.Context, offset int64) (grpc.ServerStreamingClient[T], error)

// ReliableStreamConfig configures retry and backoff when creating a stream.
type ReliableStreamConfig struct {
	Logger        zerolog.Logger
	MaxRetries    int           // 0 means unlimited retries
	BackoffMin    time.Duration // default 100ms
	BackoffMax    time.Duration // default 3s
	BackoffFactor float64       // default 2
}

// DefaultReliableStreamConfig returns a config with standard backoff and no retry limit.
func DefaultReliableStreamConfig(logger zerolog.Logger, maxRetries int) ReliableStreamConfig {
	return ReliableStreamConfig{
		Logger:        logger,
		MaxRetries:    maxRetries,
		BackoffMin:    100 * time.Millisecond,
		BackoffMax:    3 * time.Second,
		BackoffFactor: 2,
	}
}

// GetStreamWithRetry creates a stream by calling createStream, retrying with backoff on failure.
// The stream starts at the given offset (meaning is defined by the factory, e.g. BeginExclusive or ActiveAtOffset).
func GetStreamWithRetry[T any](ctx context.Context, offset int64, createStream StreamFactory[T], config ReliableStreamConfig) (grpc.ServerStreamingClient[T], error) {
	min := config.BackoffMin
	if min == 0 {
		min = 100 * time.Millisecond
	}
	max := config.BackoffMax
	if max == 0 {
		max = 3 * time.Second
	}
	factor := config.BackoffFactor
	if factor == 0 {
		factor = 2
	}
	b := &backoff.Backoff{Min: min, Max: max, Factor: factor}
	b.Reset()

	var lastErr error
	for attempt := 0; ; attempt++ {
		if config.MaxRetries > 0 && attempt > config.MaxRetries {
			return nil, fmt.Errorf("max retries reached: %w", lastErr)
		}

		stream, err := createStream(ctx, offset)
		if err == nil {
			return stream, nil
		}
		lastErr = err

		wait := b.Duration()
		config.Logger.Warn().Err(err).Str("wait", wait.String()).Int("attempt", attempt).Msg("Failed to create stream, retrying")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
			// continue
		}
	}
}

// ReceiveFromStream runs a goroutine that Recv()s from the stream and sends each response to the returned channel.
// The stream's Recv returns (*T, error) for proto messages; the response channel carries *T.
// When Recv returns an error (including io.EOF), that error is sent on the error channel and both channels are closed.
// The caller must not close the stream before consuming the error channel; the goroutine does not close the stream.
func ReceiveFromStream[T any](ctx context.Context, stream grpc.ServerStreamingClient[T]) (<-chan *T, <-chan error) {
	respChan := make(chan *T)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		for {
			resp, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case respChan <- resp:
				// continue
			}
		}
	}()

	return respChan, errChan
}
