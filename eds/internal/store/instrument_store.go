package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
)

var ErrHoldingDisclosureNotFound = errors.New("holding disclosure not found")

// InstrumentHoldingStore provides access to the latest instrument holdings for a given instrument ID.
type InstrumentHoldingStore interface {
	// GetInstrumentHolding returns the latest instrument holding for the given instrument ID.
	// Returns nil if no instrument holding is found.
	GetInstrumentHolding(ctx context.Context, instrumentID splice_api_token_holding_v1.InstrumentId) (*apiv2.DisclosedContract, error)
}

var _ InstrumentHoldingStore = &InstrumentHoldingStoreService{}

type InstrumentHoldingStoreConfig struct {
	Logger        zerolog.Logger
	Owner         types.PARTY
	StateService  apiv2.StateServiceClient
	UpdateService apiv2.UpdateServiceClient
	MaxRetries    int
}

type InstrumentHoldingStoreService struct {
	logger zerolog.Logger

	// owner is the party that we're tracking the holdings of.
	owner types.PARTY

	// holdingDisclosures is a map of instrument IDs to the owner's latest holding disclosure for that instrument.
	holdingDisclosures map[splice_api_token_holding_v1.InstrumentId]*apiv2.DisclosedContract

	// stateService is the service used to get the latest active contracts via interface id,
	// specifically the #splice-api-token-holding-v1 interface.
	stateService apiv2.StateServiceClient
	// updateService is the service used to subscribe to updates from the Canton participant,
	// specifically fetching _new_ holdings.
	updateService apiv2.UpdateServiceClient

	// maxRetries is the maximum number of retries when creating the active contracts stream (0 = unlimited).
	maxRetries int

	// mux is to protect the holdings map.
	mux sync.RWMutex
}

func NewInstrumentHoldingStore(
	config InstrumentHoldingStoreConfig,
) *InstrumentHoldingStoreService {
	return &InstrumentHoldingStoreService{
		logger: config.Logger.With().Str("component", "InstrumentHoldingStoreService").Logger(),

		owner:              config.Owner,
		holdingDisclosures: make(map[splice_api_token_holding_v1.InstrumentId]*apiv2.DisclosedContract),

		stateService:  config.StateService,
		updateService: config.UpdateService,

		maxRetries: config.MaxRetries,
	}
}

// updateStreamFactory returns a StreamFactory that creates GetUpdates streams for this store's owner and Holding interface.
// It mirrors UpdateStore's updateStreamFactory but filters for the Holding interface only.
func (i *InstrumentHoldingStoreService) updateStreamFactory() StreamFactory[apiv2.GetUpdatesResponse] {
	filtersByParty := map[string]*apiv2.Filters{
		string(i.owner): {
			Cumulative: []*apiv2.CumulativeFilter{
				{
					IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
						InterfaceId: &apiv2.Identifier{
							PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
							ModuleName: "Splice.Api.Token.HoldingV1",
							EntityName: "Holding",
						},
						IncludeInterfaceView:    true,
						IncludeCreatedEventBlob: true,
					}},
				},
			},
		},
	}

	return func(ctx context.Context, offset int64) (grpc.ServerStreamingClient[apiv2.GetUpdatesResponse], error) {
		return i.updateService.GetUpdates(ctx, &apiv2.GetUpdatesRequest{
			BeginExclusive: offset,
			EndInclusive:   nil, // not set, stream will not terminate
			UpdateFormat: &apiv2.UpdateFormat{
				IncludeTransactions: &apiv2.TransactionFormat{
					EventFormat: &apiv2.EventFormat{
						FiltersByParty: filtersByParty,
						Verbose:        true,
					},
					TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				},
			},
		})
	}
}

func (i *InstrumentHoldingStoreService) Run(ctx context.Context) error {
	i.logger.Debug().Msg("Starting InstrumentHoldingStoreService")
	ledgerEndResponse, err := i.stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return fmt.Errorf("failed to get ledger end: %w", err)
	}
	offset := ledgerEndResponse.Offset

	// Subscribe to updates; reconnect on stream error (including EOF)
	for {
		i.logger.Debug().Int64("offset", offset).Msg("Subscribing to update stream")
		stream, err := GetStreamWithRetry(ctx, offset, i.updateStreamFactory(), DefaultReliableStreamConfig(i.logger, i.maxRetries))
		if err != nil {
			return fmt.Errorf("failed to create update stream: %w", err)
		}

		i.logger.Debug().Int64("offset", offset).Msg("Update stream created, listening for updates")
		respChan, errChan := ReceiveFromStream(ctx, stream)
	reconnect:
		for {
			select {
			case err := <-errChan:
				if errors.Is(err, io.EOF) {
					i.logger.Debug().Msg("Update stream closed by server, reconnecting")
					break reconnect
				}
				if err != nil {
					i.logger.Warn().Err(err).Msg("Update stream closed by server, reconnecting")
					break reconnect
				}
			case <-ctx.Done():
				i.logger.Debug().Msg("Context cancelled, stopping InstrumentHoldingStoreService")
				_ = stream.CloseSend()

				return ctx.Err()
			case resp, ok := <-respChan:
				if !ok {
					// respChan closed; err was already sent on errChan, inner loop will handle it
					break reconnect
				}
				if tx, ok := resp.GetUpdate().(*apiv2.GetUpdatesResponse_Transaction); ok {
					for _, event := range tx.Transaction.GetEvents() {
						if createdEvent, ok := event.GetEvent().(*apiv2.Event_Created); ok {
							ac := &apiv2.ActiveContract{
								CreatedEvent:   createdEvent.Created,
								SynchronizerId: tx.Transaction.GetSynchronizerId(),
							}
							holdingDisclosure, holdingView, err := getHoldingDisclosureAndView(i.owner, ac)
							if err != nil {
								i.logger.Debug().Err(err).Str("contractID", createdEvent.Created.GetContractId()).Msg("Failed to get holding disclosure for created event, skipping")
								continue
							}
							i.logger.Info().
								Any("holdingView", holdingView).
								Any("holdingDisclosure", holdingDisclosure).
								Msg("Recording holding disclosure")
							i.mux.Lock()
							i.holdingDisclosures[holdingView.InstrumentId] = holdingDisclosure
							i.mux.Unlock()
						}
					}
					offset = tx.Transaction.GetOffset()
				}
			}
		}
	}
}

func getHoldingDisclosureAndView(expectedOwner types.PARTY, ac *apiv2.ActiveContract) (*apiv2.DisclosedContract, *splice_api_token_holding_v1.HoldingView, error) {
	// At this point we have an active holding contract, we should form the disclosure for it.
	holdingView, err := bindings.UnmarshalCreatedEvent[splice_api_token_holding_v1.HoldingView](ac.GetCreatedEvent())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal holding view for active contract: %w", err)
	}

	if holdingView.Owner != expectedOwner {
		return nil, nil, fmt.Errorf("holding owner does not match expected owner: %w", err)
	}

	holdingDisclosure := &apiv2.DisclosedContract{
		TemplateId:       ac.GetCreatedEvent().GetTemplateId(),
		ContractId:       ac.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: ac.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   ac.GetSynchronizerId(),
	}

	return holdingDisclosure, holdingView, nil
}

// GetInstrumentHolding implements [InstrumentHoldingStore].
func (i *InstrumentHoldingStoreService) GetInstrumentHolding(ctx context.Context, instrumentID splice_api_token_holding_v1.InstrumentId) (*apiv2.DisclosedContract, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	holdingDisclosure, ok := i.holdingDisclosures[instrumentID]
	if !ok {
		return nil, ErrHoldingDisclosureNotFound
	}

	return holdingDisclosure, nil
}
