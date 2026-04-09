package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Store[K comparable, V any] interface {
	// Get returns the currently stored value for the given key.
	// Returns false if no value is found.
	Get(key K) (value V, ok bool)
}

type ContractUpdate[K comparable, V any] struct {
	Key   K
	Value V
}

type ContractStore[K comparable, V any] struct {
	logger zerolog.Logger
	// updateService is the service used to subscribe to updates from the Canton participant
	updateService apiv2.UpdateServiceClient
	// stateService is the service used to get the latest active contracts from the Canton participant
	stateService apiv2.StateServiceClient

	metrics Metrics

	// maxRetries is the maximum number of retries when creating the active contracts stream (0 = unlimited).
	maxRetries int
	// reconnectBackoff is the delay before reconnecting after the stream closes (avoids log thrashing on misconfiguration).
	reconnectBackoff time.Duration

	// The filters to apply when backfilling/subscribing to updates
	filtersByParty map[string]*apiv2.Filters

	// mutex to protect the contracts map and ledgerEnd
	mux sync.RWMutex
	// the last-processed ledgerEnd
	ledgerEnd int64
	// the currently stored values
	contracts map[K]V

	// handler functions

	// handleActiveContract is used to convert an ActiveContract to store updates.
	// It is being called during backfill for each currently active contract returned by the configured filters.
	handleActiveContract func(ctx context.Context, store *ContractStore[K, V], activeContract *apiv2.ActiveContract) (updates []ContractUpdate[K, V], err error)
	// handleCreatedEvent is used to convert a CreatedEvent to store updates.
	// It is being called for each CreatedEvent received in the update stream that matches the configured filters.
	handleCreatedEvent func(ctx context.Context, store *ContractStore[K, V], transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) (updates []ContractUpdate[K, V], err error)
}

// Run runs the ContractStore. It will start by initializing the ContractStore with a backfill of all existing active contracts
// and subscribe to incremental updates afterward.
// Run is a long-running function that needs to keep running in the background in order for the ContractStore to keep
// up-to-date.
// To terminate Run, cancel the context.
func (s *ContractStore[K, V]) Run(ctx context.Context) error {
	s.logger.Debug().Msg("Starting ContractStore")
	ledgerEndResponse, err := s.stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return fmt.Errorf("failed to get ledger end: %w", err)
	}
	offset := ledgerEndResponse.Offset

	// backfill
	s.logger.Debug().Int64("offset", offset).Msg("Starting backfill")
	activeContracts, err := s.backfill(ctx, offset)
	if err != nil {
		return fmt.Errorf("backfill failed: %w", err)
	}
	s.mux.Lock()
	s.ledgerEnd = offset
	s.contracts = activeContracts
	s.mux.Unlock()
	s.logger.Debug().Int("activeContracts", len(activeContracts)).Msg("Backfill complete")

	// Create uptime ticker
	// Ticking once a second to increase the uptime counter metric
	// Checking if the metric is increasing can be used to verify that the ContractStore is still subscribed and ready to process updates
	uptimeTicker := time.NewTicker(time.Second)
	defer uptimeTicker.Stop()

	// Subscribe to updates
	for {
		// Start streaming from s.ledgerEnd (exclusive)
		s.logger.Debug().Int64("offset", s.ledgerEnd).Msg("Subscribing to update stream")
		stream, err := GetStreamWithRetry(ctx, s.ledgerEnd, s.updateStreamFactory(), DefaultReliableStreamConfig(s.logger, s.maxRetries))
		if err != nil {
			return fmt.Errorf("failed to create update stream: %w", err)
		}

		// Start receiving from stream
		s.logger.Debug().Int64("offset", s.ledgerEnd).Msg("Update stream created, listening for updates")
		updateChan, errChan := ReceiveFromStream(ctx, stream)
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
				// Increment uptime metric on every tick
				s.metrics.IncrementStoreSubscriptionUptime(ctx)
			case resp, ok := <-updateChan:
				if !ok {
					// updateChan closed; err was already sent on errChan, inner loop will handle it
					break reconnect
				}
				if transaction, ok := resp.GetUpdate().(*apiv2.GetUpdatesResponse_Transaction); ok {
					s.logger.Trace().Str("updateID", transaction.Transaction.GetUpdateId()).Int("events", len(transaction.Transaction.GetEvents())).Msg("Update received")
					for _, event := range transaction.Transaction.GetEvents() {
						switch event := event.GetEvent().(type) {
						case *apiv2.Event_Created:
							s.logger.Trace().
								Str("contractID", event.Created.GetContractId()).
								Str("packageName", event.Created.GetPackageName()).
								Str("packageID", event.Created.GetTemplateId().GetPackageId()).
								Str("moduleName", event.Created.GetTemplateId().GetModuleName()).
								Str("entityName", event.Created.GetTemplateId().GetEntityName()).
								Msg("Received CreatedEvent")

							updates, err := s.handleCreatedEvent(ctx, s, transaction.Transaction, event.Created)
							if err != nil {
								return fmt.Errorf("failed to handle created event: %w", err)
							}

							s.mux.Lock()
							// Save updates
							for _, update := range updates {
								s.logger.Trace().Any("key", update.Key).Msg("Updating contract")
								s.contracts[update.Key] = update.Value
							}
							// Update ledgerEnd
							s.ledgerEnd = transaction.Transaction.GetOffset()
							s.mux.Unlock()

							// Update metrics
							s.metrics.IncrementStoreUpdatesCounter(ctx)
							s.metrics.RecordStoreLedgerEndGauge(ctx, transaction.Transaction.GetOffset())
						default:
							continue
						}
					}
				}
			}
		}

		// Backoff before reconnecting to avoid thrashing logs (e.g. on misconfiguration).
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.reconnectBackoff):
			s.logger.Debug().Str("backoff", s.reconnectBackoff.String()).Msg("Reconnect backoff elapsed, reconnecting")
		}
	}
}

// updateStreamFactory returns a StreamFactory that creates GetUpdates streams for this store's filters.
func (s *ContractStore[K, V]) updateStreamFactory() StreamFactory[apiv2.GetUpdatesResponse] {
	return func(ctx context.Context, offset int64) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
		return s.updateService.GetUpdates(ctx, &apiv2.GetUpdatesRequest{
			BeginExclusive: offset,
			EndInclusive:   nil, // not set, stream will not terminate
			UpdateFormat: &apiv2.UpdateFormat{
				IncludeTransactions: &apiv2.TransactionFormat{
					EventFormat: &apiv2.EventFormat{
						FiltersByParty: s.filtersByParty,
						Verbose:        true,
					},
					TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				},
			},
		})
	}
}

// backfill returns all currently-active contracts for s.filtersByParty at a given offset.
func (s *ContractStore[K, V]) backfill(ctx context.Context, offset int64) (map[K]V, error) {
	activeContracts := make(map[K]V)

	activeContractsResponse, err := s.stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: s.filtersByParty,
			Verbose:        true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			_ = activeContractsResponse.CloseSend()
			return activeContracts, ctx.Err()
		default:
			resp, err := activeContractsResponse.Recv()
			if errors.Is(err, io.EOF) {
				return activeContracts, nil
			}
			if err != nil {
				return nil, fmt.Errorf("failed to receive active contract: %w", err)
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

				updates, err := s.handleActiveContract(ctx, s, c.ActiveContract)
				if err != nil {
					return nil, fmt.Errorf("failed to handle active contract: %w", err)
				}

				for _, update := range updates {
					s.logger.Trace().Any("key", update.Key).Msg("Backfilling contract")
					activeContracts[update.Key] = update.Value
				}
			default:
				continue
			}
		}
	}
}

func (s *ContractStore[K, V]) Get(key K) (V, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	value, ok := s.contracts[key]

	return value, ok
}
