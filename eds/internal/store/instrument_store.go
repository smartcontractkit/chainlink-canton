package store

import (
	"context"
	"fmt"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_holding_v1"
)

type InstrumentHoldingStore struct {
	logger zerolog.Logger

	stream  *BackfilledStream
	filters FiltersByParty

	mux sync.RWMutex
	// Party -> InstrumentId -> ContractId -> Holding
	holdings        map[types.PARTY]map[splice_api_token_holding_v1.InstrumentId]map[types.CONTRACT_ID]*apiv2.ActiveContract
	activeContracts map[types.CONTRACT_ID]splice_api_token_holding_v1.HoldingView
}

func NewInstrumentHoldingStore(
	logger zerolog.Logger,
	updateService apiv2.UpdateServiceClient,
	stateService apiv2.StateServiceClient,
	metrics Metrics,
) *InstrumentHoldingStore {
	logger = logger.With().Str("component", "InstrumentHoldingStore").Logger()
	s := &InstrumentHoldingStore{
		logger: zerolog.Logger{},
		stream: &BackfilledStream{
			logger:        logger,
			metrics:       metrics,
			updateService: updateService,
			stateService:  stateService,
		},
		filters: make(FiltersByParty),
		mux:     sync.RWMutex{},
	}

	return s
}

func (s *InstrumentHoldingStore) RegisterParty(parties ...string) {
	for _, party := range parties {
		s.filters[party] = &apiv2.Filters{
			Cumulative: []*apiv2.CumulativeFilter{
				{
					IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
						InterfaceId: &apiv2.Identifier{
							PackageId:  "#splice-api-token-holding-v1",
							ModuleName: "Splice.Api.Token.HoldingV1",
							EntityName: "Holding",
						},
						IncludeInterfaceView:    true,
						IncludeCreatedEventBlob: true,
					}},
				},
			},
		}
	}
}

func (s *InstrumentHoldingStore) Run(ctx context.Context, streamConfig StreamConfig) error {
	s.holdings = make(map[types.PARTY]map[splice_api_token_holding_v1.InstrumentId]map[types.CONTRACT_ID]*apiv2.ActiveContract)
	s.activeContracts = make(map[types.CONTRACT_ID]splice_api_token_holding_v1.HoldingView)

	return s.stream.Run(ctx, s.filters, streamConfig, s.onActiveContract, s.onCreatedEvent, s.onArchivedEvent, nil)
}

func (s *InstrumentHoldingStore) GetHolding(party types.PARTY, instrumentId splice_api_token_holding_v1.InstrumentId) ([]*apiv2.ActiveContract, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	partyHoldings, ok := s.holdings[party]
	if !ok {
		return nil, false
	}
	instrumentHoldings, ok := partyHoldings[instrumentId]
	if !ok {
		return nil, false
	}
	holdings := maps.Values(instrumentHoldings)

	// Clone all holdings before returning them
	out := make([]*apiv2.ActiveContract, len(holdings))
	for i, holding := range holdings {
		out[i] = proto.CloneOf(holding)
	}

	return out, true
}

func (s *InstrumentHoldingStore) onActiveContract(ctx context.Context, activeContract *apiv2.ActiveContract) error {
	holdingViews, err := getHoldingViews(activeContract.GetCreatedEvent())
	if err != nil {
		s.logger.Debug().Err(err).Str("contractID", activeContract.GetCreatedEvent().GetContractId()).Msg("Failed to get HoldingViews for active contract, skipping")
		return nil
	}

	s.mux.Lock()
	for _, view := range holdingViews {
		partyHoldings := s.holdings[view.Owner]
		if partyHoldings == nil {
			partyHoldings = make(map[splice_api_token_holding_v1.InstrumentId]map[types.CONTRACT_ID]*apiv2.ActiveContract)
		}

		instrumentHoldings := partyHoldings[view.InstrumentId]
		if instrumentHoldings == nil {
			instrumentHoldings = make(map[types.CONTRACT_ID]*apiv2.ActiveContract)
		}

		s.activeContracts[types.CONTRACT_ID(activeContract.GetCreatedEvent().GetContractId())] = view
		instrumentHoldings[types.CONTRACT_ID(activeContract.GetCreatedEvent().GetContractId())] = proto.CloneOf(activeContract) // Create a copy of the ActiveContract
		partyHoldings[view.InstrumentId] = instrumentHoldings
		s.holdings[view.Owner] = partyHoldings
	}
	s.mux.Unlock()

	return nil
}

func (s *InstrumentHoldingStore) onCreatedEvent(ctx context.Context, transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) error {
	holdingViews, err := getHoldingViews(createdEvent)
	if err != nil {
		s.logger.Debug().Err(err).Str("contractID", createdEvent.GetContractId()).Msg("Failed to get HoldingViews for active contract, skipping")
		return nil
	}

	s.mux.Lock()
	for _, view := range holdingViews {
		partyHoldings := s.holdings[view.Owner]
		if partyHoldings == nil {
			partyHoldings = make(map[splice_api_token_holding_v1.InstrumentId]map[types.CONTRACT_ID]*apiv2.ActiveContract)
		}

		instrumentHoldings := partyHoldings[view.InstrumentId]
		if instrumentHoldings == nil {
			instrumentHoldings = make(map[types.CONTRACT_ID]*apiv2.ActiveContract)
		}

		s.activeContracts[types.CONTRACT_ID(createdEvent.GetContractId())] = view
		instrumentHoldings[types.CONTRACT_ID(createdEvent.GetContractId())] = &apiv2.ActiveContract{
			CreatedEvent:        proto.CloneOf(createdEvent), // Create a copy of the CreatedEvent
			SynchronizerId:      transaction.GetSynchronizerId(),
			ReassignmentCounter: 0,
		}
		partyHoldings[view.InstrumentId] = instrumentHoldings
		s.holdings[view.Owner] = partyHoldings
	}
	s.mux.Unlock()

	return nil
}

func (s *InstrumentHoldingStore) onArchivedEvent(ctx context.Context, _ *apiv2.Transaction, archivedEvent *apiv2.ArchivedEvent) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	holdingView, ok := s.activeContracts[types.CONTRACT_ID(archivedEvent.GetContractId())]
	if !ok {
		return fmt.Errorf("archived event for contract %v with no active contract record", archivedEvent.GetContractId())
	}
	partyHoldings, ok := s.holdings[holdingView.Owner]
	if !ok {
		return fmt.Errorf("archived event for contract %v with no holdings for owner %v", archivedEvent.GetContractId(), holdingView.Owner)
	}
	instrumentHoldings, ok := partyHoldings[holdingView.InstrumentId]
	if !ok {
		return fmt.Errorf("archived event for contract %v with no holdings for instrument %v for owner %v", archivedEvent.GetContractId(), holdingView.InstrumentId, holdingView.Owner)
	}
	delete(s.activeContracts, types.CONTRACT_ID(archivedEvent.GetContractId()))
	delete(instrumentHoldings, types.CONTRACT_ID(archivedEvent.GetContractId()))
	partyHoldings[holdingView.InstrumentId] = instrumentHoldings
	s.holdings[holdingView.Owner] = partyHoldings

	return nil
}

func getHoldingViews(createdEvent *apiv2.CreatedEvent) ([]splice_api_token_holding_v1.HoldingView, error) {
	var holdingViews []splice_api_token_holding_v1.HoldingView
	for _, interfaceView := range createdEvent.GetInterfaceViews() {
		// PackageId seems to be the hashed package ID, not the human-readable one.
		// We probably shouldn't check against a specific hash, so we ignore that part of the interface ID.
		if interfaceView.GetInterfaceId().GetModuleName() == "Splice.Api.Token.HoldingV1" &&
			interfaceView.GetInterfaceId().GetEntityName() == "Holding" {
			var holdingView splice_api_token_holding_v1.HoldingView
			if err := ledger.RecordToStruct(interfaceView.GetViewValue(), &holdingView); err != nil {
				return nil, fmt.Errorf("failed to parse HoldingView: %w", err)
			}
			holdingViews = append(holdingViews, holdingView)
		}
	}

	return holdingViews, nil
}
