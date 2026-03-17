package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

var ErrHoldingDisclosureNotFound = errors.New("holding disclosure not found")

// InstrumentHoldingStore provides access to the latest instrument holdings for a given instrument ID.
type InstrumentHoldingStore interface {
	// GetInstrumentHolding returns the latest instrument holding for the given instrument ID.
	// Returns nil if no instrument holding is found.
	GetInstrumentHolding(ctx context.Context, instrumentID splice_api_token_holding_v1.InstrumentId) (*apiv2.DisclosedContract, error)
}

var _ InstrumentHoldingStore = &InstrumentHoldingStoreService{}

type InstrumentHoldingStoreService struct {
	logger zerolog.Logger

	// owner is the party that we're tracking the holdings of.
	owner types.PARTY

	// holdingDisclosures is a map of instrument IDs to the owner's latest holding disclosure for that instrument.
	holdingDisclosures map[splice_api_token_holding_v1.InstrumentId]*apiv2.DisclosedContract

	// stateService is the service used to get the latest active contracts via interface id,
	// specifically the #splice-api-token-holding-v1 interface.
	stateService apiv2.StateServiceClient

	// MaxRetries is the maximum number of retries when creating the active contracts stream (0 = unlimited).
	MaxRetries int

	// mux is to protect the holdings map.
	mux sync.RWMutex
}

func NewInstrumentHoldingStore(
	owner types.PARTY,
	stateService apiv2.StateServiceClient,
	logger zerolog.Logger,
) *InstrumentHoldingStoreService {
	return &InstrumentHoldingStoreService{
		logger: logger.With().Str("component", "InstrumentHoldingStoreService").Logger(),

		owner:              owner,
		holdingDisclosures: make(map[splice_api_token_holding_v1.InstrumentId]*apiv2.DisclosedContract),

		stateService: stateService,

		mux: sync.RWMutex{},
	}
}

// activeContractsStreamFactory returns a StreamFactory that creates GetActiveContracts streams for this store's owner and Holding interface.
func (i *InstrumentHoldingStoreService) activeContractsStreamFactory() StreamFactory[apiv2.GetActiveContractsResponse] {
	return func(ctx context.Context, offset int64) (grpc.ServerStreamingClient[apiv2.GetActiveContractsResponse], error) {
		return i.stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
			ActiveAtOffset: offset,
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
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
				},
				Verbose: true,
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

	stream, err := GetStreamWithRetry(ctx, offset, i.activeContractsStreamFactory(), DefaultReliableStreamConfig(i.logger, i.MaxRetries))
	if err != nil {
		return fmt.Errorf("failed to get active contracts: %w", err)
	}
	defer stream.CloseSend()

	respChan, errChan := ReceiveFromStream(ctx, stream)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			if err != nil {
				return fmt.Errorf("failed to receive active contract: %w", err)
			}
		case resp, ok := <-respChan:
			if !ok {
				// respChan closed; error was already sent on errChan
				err := <-errChan
				if err != nil {
					return fmt.Errorf("failed to receive active contract: %w", err)
				}
				return nil
			}
			if ac, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
				i.logger.Debug().
					Str("contractID", ac.ActiveContract.GetCreatedEvent().GetContractId()).
					Str("packageName", ac.ActiveContract.GetCreatedEvent().GetPackageName()).
					Str("packageID", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetPackageId()).
					Str("moduleName", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetModuleName()).
					Str("entityName", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName()).
					Msg("Received active contract")

				holdingDisclosure, holdingView, err := getHoldingDisclosureAndView(i.owner, ac.ActiveContract)
				if err != nil {
					i.logger.Debug().Err(err).Str("contractID", ac.ActiveContract.GetCreatedEvent().GetContractId()).Msg("Failed to get holding disclosure for active contract, skipping")
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
