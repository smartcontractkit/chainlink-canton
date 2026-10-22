package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/jpillora/backoff"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type FiltersByParty map[string]*apiv2.Filters

type StreamConfig struct {
	MaxRetries    int           // 0 means unlimited retries
	BackoffMin    time.Duration // default 100ms
	BackoffMax    time.Duration // default 3s
	BackoffFactor float64       // default 2
	// MaxEmptyReconnects fails closed after this many consecutive stream
	// reconnects that produced zero updates (0 means never fail on this).
	MaxEmptyReconnects int
	// MaxDisconnected fails closed after a single failure-to-connect gap
	// longer than this (0 means never fail on this).
	MaxDisconnected time.Duration
}

func DefaultStreamConfig() StreamConfig {
	return StreamConfig{
		MaxRetries:         0,
		BackoffMin:         time.Millisecond * 100,
		BackoffMax:         time.Second * 3,
		BackoffFactor:      2,
		MaxEmptyReconnects: 20,
		MaxDisconnected:    time.Minute * 15,
	}
}

type BackfilledStream struct {
	logger  zerolog.Logger
	metrics Metrics

	// updateService is the service used to subscribe to updates from the Canton participant
	updateService apiv2.UpdateServiceClient
	// stateService is the service used to get the latest active contracts from the Canton participant
	stateService apiv2.StateServiceClient
}

func (s *BackfilledStream) Run(
	ctx context.Context,
	filters FiltersByParty,
	streamConfig StreamConfig,
	onActiveContract func(ctx context.Context, activeContract *apiv2.ActiveContract) error,
	onCreatedEvent func(ctx context.Context, transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) error,
	onArchivedEvent func(ctx context.Context, transaction *apiv2.Transaction, archivedEvent *apiv2.ArchivedEvent) error,
	onExercisedEvent func(ctx context.Context, transaction *apiv2.Transaction, exercisedEvent *apiv2.ExercisedEvent) error,
	onFinishedBackfill func(),
) error {
	s.logger.Debug().Msg("Starting BackfilledStream")
	ledgerEndResponse, err := s.stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return fmt.Errorf("failed to get ledger end: %w", err)
	}
	offset := ledgerEndResponse.Offset

	// Backfill
	err = s.backfill(ctx, offset, filters, onActiveContract)
	if err != nil {
		return fmt.Errorf("failed to backfill: %w", err)
	}

	s.logger.Debug().Msg("Backfill complete, starting stream...")
	if onFinishedBackfill != nil {
		onFinishedBackfill()
	}

	// Run stream
	return s.stream(ctx, offset, filters, streamConfig, onCreatedEvent, onArchivedEvent, onExercisedEvent)
}

func (s *BackfilledStream) backfill(
	ctx context.Context,
	offset int64,
	filters FiltersByParty,
	onActiveContract func(ctx context.Context, activeContract *apiv2.ActiveContract) error,
) error {
	activeContractsResponse, err := s.stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: filters,
			Verbose:        true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get active contracts: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			_ = activeContractsResponse.CloseSend()
			return ctx.Err()
		default:
			resp, err := activeContractsResponse.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to receive active contract: %w", err)
			}

			switch c := resp.GetContractEntry().(type) {
			case *apiv2.GetActiveContractsResponse_ActiveContract:
				s.logger.Trace().
					Str("contractID", c.ActiveContract.GetCreatedEvent().GetContractId()).
					Str("packageName", c.ActiveContract.GetCreatedEvent().GetPackageName()).
					Str("packageID", c.ActiveContract.GetCreatedEvent().GetTemplateId().GetPackageId()).
					Str("moduleName", c.ActiveContract.GetCreatedEvent().GetTemplateId().GetModuleName()).
					Str("entityName", c.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName()).
					Msg("Backfilling active contract")

				err = onActiveContract(ctx, c.ActiveContract)
				if err != nil {
					return fmt.Errorf("failed to handle active contract: %w", err)
				}
			default:
				continue
			}
		}
	}
}

func (s *BackfilledStream) stream(
	ctx context.Context,
	offset int64,
	filters FiltersByParty,
	streamConfig StreamConfig,
	onCreatedEvent func(ctx context.Context, transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) error,
	onArchivedEvent func(ctx context.Context, transaction *apiv2.Transaction, archivedEvent *apiv2.ArchivedEvent) error,
	onExercisedEvent func(ctx context.Context, transaction *apiv2.Transaction, exercisedEvent *apiv2.ExercisedEvent) error,
) error {
	// Create uptime ticker
	// Ticking once a second to increase the uptime counter metric, but only
	// while the current stream incarnation has delivered at least one update:
	// a quiet or broken stream no longer fakes subscription health.
	uptimeTicker := time.NewTicker(time.Second)
	defer uptimeTicker.Stop()

	// Fail-closed policy: reconnecting is pointless if we never receive updates
	// or cannot re-establish the stream at all. Run's caller shuts the service down.
	budget := &reconnectBudget{
		maxEmptyReconnects: streamConfig.MaxEmptyReconnects,
		maxDisconnected:    streamConfig.MaxDisconnected,
	}
	// Backoff between consecutive empty reconnects, so a burst of instant
	// stream failures cannot spend the whole reconnect budget in a second.
	reconnectBackoff := &backoff.Backoff{Min: streamConfig.BackoffMin, Max: streamConfig.BackoffMax, Factor: streamConfig.BackoffFactor}

	for {
		// Start streaming from offset (exclusive)
		s.logger.Debug().Int64("offset", offset).Msg("Subscribing to update stream")
		stream, err := s.getStreamWithRetry(ctx, offset, filters, streamConfig, budget)
		if err != nil {
			return fmt.Errorf("failed to create update stream: %w", err)
		}

		// Start receiving from stream
		s.logger.Debug().Int64("offset", offset).Msg("Update stream created, listening for updates")
		updateChan, errChan := ReceiveFromStream(ctx, stream)
		updated := false
	reconnect:
		for {
			select {
			case err := <-errChan:
				if errors.Is(err, io.EOF) {
					s.logger.Debug().Msg("Update stream closed by server, reconnecting")
					break reconnect
				}
				if err != nil {
					s.logger.Warn().Err(err).Msg("Update stream closed by server, reconnecting")
					break reconnect
				}
			case <-ctx.Done():
				s.logger.Debug().Msg("Context cancelled, stopping ContractStore")
				_ = stream.CloseSend()

				return ctx.Err()
			case <-uptimeTicker.C:
				// Increment uptime metric on every tick, but only while this
				// stream incarnation is delivering updates
				if updated {
					s.metrics.IncrementStoreSubscriptionUptime(ctx)
				}
			case resp, ok := <-updateChan:
				if !ok {
					// updateChan closed; err was already sent on errChan, inner loop will handle it
					break reconnect
				}
				if transaction, ok := resp.GetUpdate().(*apiv2.GetUpdatesResponse_Transaction); ok {
					// A real update is proof of health: reset the reconnect budget and backoff.
					budget.update()
					reconnectBackoff.Reset()
					updated = true
					s.logger.Trace().Str("updateID", transaction.Transaction.GetUpdateId()).Int("events", len(transaction.Transaction.GetEvents())).Msg("Update received")
					for _, event := range transaction.Transaction.GetEvents() {
						switch event := event.GetEvent().(type) {
						case *apiv2.Event_Created:
							if onCreatedEvent == nil {
								continue
							}
							s.logger.Trace().
								Str("contractID", event.Created.GetContractId()).
								Str("packageName", event.Created.GetPackageName()).
								Str("packageID", event.Created.GetTemplateId().GetPackageId()).
								Str("moduleName", event.Created.GetTemplateId().GetModuleName()).
								Str("entityName", event.Created.GetTemplateId().GetEntityName()).
								Msg("Received CreatedEvent")
							err = onCreatedEvent(ctx, transaction.Transaction, event.Created)
						case *apiv2.Event_Archived:
							if onArchivedEvent == nil {
								continue
							}
							s.logger.Trace().
								Str("contractID", event.Archived.GetContractId()).
								Str("packageName", event.Archived.GetPackageName()).
								Str("packageID", event.Archived.GetTemplateId().GetPackageId()).
								Str("moduleName", event.Archived.GetTemplateId().GetModuleName()).
								Str("entityName", event.Archived.GetTemplateId().GetEntityName()).
								Msg("Received ArchivedEvent")
							err = onArchivedEvent(ctx, transaction.Transaction, event.Archived)
						case *apiv2.Event_Exercised:
							if onExercisedEvent == nil {
								continue
							}
							s.logger.Trace().
								Str("contractID", event.Exercised.GetContractId()).
								Str("packageName", event.Exercised.GetPackageName()).
								Str("packageID", event.Exercised.GetTemplateId().GetPackageId()).
								Str("moduleName", event.Exercised.GetTemplateId().GetModuleName()).
								Str("entityName", event.Exercised.GetTemplateId().GetEntityName()).
								Msg("Received ExercisedEvent")
							err = onExercisedEvent(ctx, transaction.Transaction, event.Exercised)
						default:
							return fmt.Errorf("unknown event type: %T", event)
						}
						if err != nil {
							return fmt.Errorf("failed to handle event: %w", err)
						}

						// Update offset
						offset = transaction.Transaction.GetOffset()

						// Update metrics
						s.metrics.IncrementStoreUpdatesCounter(ctx)
						s.metrics.RecordStoreLedgerEndGauge(ctx, transaction.Transaction.GetOffset())
					}
				}
			}
		}

		// Stream ended without a context cancellation. If this incarnation never
		// delivered a single update, drain the reconnect budget and fail closed
		// when it is spent: reconnecting is evidently not making progress.
		if !updated {
			if err := ctx.Err(); err != nil {
				return err
			}
			if budget.emptyReconnect() {
				return fmt.Errorf("update stream reconnected %d times without receiving any update (budget %d), failing closed", budget.emptyReconnects, budget.maxEmptyReconnects)
			}

			// Wait before resubscribing: a burst of instant stream failures must
			// not spend the whole budget (or hot-spin) in under a second.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(reconnectBackoff.Duration()):
			}
		}
	}
}

// reconnectBudget implements the fail-closed policy for the update stream:
// give up after maxEmptyReconnects consecutive stream incarnations that broke
// without delivering a single update, or after one failure-to-connect gap
// longer than maxDisconnected. A zero value disables each check.
// shortcut: liveness is only proven by received updates, so a quiet ledger
// whose streams also keep closing can trip MaxEmptyReconnects; upgrade by
// counting successful stream establishments as liveness evidence if that bites.
type reconnectBudget struct {
	maxEmptyReconnects int
	maxDisconnected    time.Duration

	emptyReconnects   int
	disconnectedSince time.Time
}

// update resets the empty-reconnect streak: the stream delivered an update.
func (b *reconnectBudget) update() { b.emptyReconnects = 0 }

// emptyReconnect records a stream incarnation that broke with zero updates and
// reports whether the empty-reconnect budget is spent.
func (b *reconnectBudget) emptyReconnect() bool {
	b.emptyReconnects++
	return b.maxEmptyReconnects > 0 && b.emptyReconnects >= b.maxEmptyReconnects
}

// connectFailed records a failed connect attempt and reports the time spent so
// far without an established stream.
func (b *reconnectBudget) connectFailed(now time.Time) time.Duration {
	if b.disconnectedSince.IsZero() {
		b.disconnectedSince = now
	}
	return now.Sub(b.disconnectedSince)
}

// connected clears the disconnected timer: a stream is established again.
func (b *reconnectBudget) connected() { b.disconnectedSince = time.Time{} }

// GetStreamWithRetry creates a stream by calling createStream, retrying with backoff on failure.
// The stream starts at the given offset (meaning is defined by the factory, e.g. BeginExclusive or ActiveAtOffset).
func (s *BackfilledStream) getStreamWithRetry(ctx context.Context, offset int64, filters FiltersByParty, streamConfig StreamConfig, budget *reconnectBudget) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
	b := &backoff.Backoff{Min: streamConfig.BackoffMin, Max: streamConfig.BackoffMax, Factor: streamConfig.BackoffFactor}
	b.Reset()

	for attempt := 1; ; attempt++ {
		stream, err := s.updateService.GetUpdates(ctx, &apiv2.GetUpdatesRequest{
			BeginExclusive: offset,
			EndInclusive:   nil, // not set, stream will not terminate
			UpdateFormat: &apiv2.UpdateFormat{
				IncludeTransactions: &apiv2.TransactionFormat{
					EventFormat: &apiv2.EventFormat{
						FiltersByParty: filters,
						Verbose:        true,
					},
					TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				},
			},
		})
		if err != nil {
			// Fail closed if we have been unable to establish a stream for too long.
			if gap := budget.connectFailed(time.Now()); budget.maxDisconnected > 0 && gap > budget.maxDisconnected {
				return nil, fmt.Errorf("disconnected for %s (budget %s): %w", gap.Round(time.Second), budget.maxDisconnected, err)
			}

			// If MaxRetries is enabled/non-zero and the maximum number of retries has been reached, return an error.
			if streamConfig.MaxRetries > 0 && attempt >= streamConfig.MaxRetries {
				return nil, fmt.Errorf("max retries (%v) reached: %w", streamConfig.MaxRetries, err)
			}

			wait := b.Duration()
			s.logger.Warn().Err(err).Str("wait", wait.String()).Int("attempt", attempt).Msg("Failed to create stream, retrying")
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
				continue
			}
		}

		budget.connected()
		return stream, nil
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
			// Set up separate goroutine for reading from stream.
			// In case the stream hangs, this allows for ReceiveFromStream to still return if the context is cancelled.
			type recvResult struct {
				msg *T
				err error
			}
			ch := make(chan recvResult, 1)
			go func() {
				msg, err := stream.Recv()
				ch <- recvResult{msg: msg, err: err}
			}()

			// Read incoming messages until context is cancelled or Recv returns an error
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case result := <-ch:
				if result.err != nil {
					errChan <- result.err
					return
				}
				respChan <- result.msg
			}
		}
	}()

	return respChan, errChan
}
