package store

import (
	"context"
	"fmt"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type ActiveContractStore struct {
	logger zerolog.Logger

	stream  *BackfilledStream
	filters FiltersByParty

	mux       sync.RWMutex
	contracts map[contracts.InstanceAddress]*apiv2.ActiveContract
}

func NewActiveContractStore(
	logger zerolog.Logger,
	updateService apiv2.UpdateServiceClient,
	stateService apiv2.StateServiceClient,
	metrics Metrics,
) *ActiveContractStore {
	logger = logger.With().Str("component", "ActiveContractStore").Logger()
	s := &ActiveContractStore{
		logger: logger,
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

type RegisteredTemplate struct {
	TemplateID contracts.TemplateID
	PartyID    string
}

// RegisterTemplates registers the given templates with the ActiveContractStore.
// RegisterTemplates must be called before any call to Run(), calling while the store is already running will lead to
// undefined behavior.
func (s *ActiveContractStore) RegisterTemplates(templates ...RegisteredTemplate) {
	for _, template := range templates {
		existingFilterForParty, ok := s.filters[template.PartyID]
		if !ok {
			existingFilterForParty = &apiv2.Filters{}
		}
		existingFilterForParty.Cumulative = append(existingFilterForParty.Cumulative, &apiv2.CumulativeFilter{
			IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
				TemplateId: &apiv2.Identifier{
					PackageId:  template.TemplateID.PackageID,
					ModuleName: template.TemplateID.ModuleName,
					EntityName: template.TemplateID.EntityName,
				},
				IncludeCreatedEventBlob: true,
			}},
		})
		s.filters[template.PartyID] = existingFilterForParty
	}
}

func (s *ActiveContractStore) Run(ctx context.Context, streamConfig StreamConfig) error {
	s.contracts = make(map[contracts.InstanceAddress]*apiv2.ActiveContract)

	return s.stream.Run(ctx, s.filters, streamConfig, s.onActiveContract, s.onCreatedEvent, nil, nil)
}

func (s *ActiveContractStore) Get(address contracts.InstanceAddress) (*apiv2.ActiveContract, bool) {
	s.mux.RLock()
	value, ok := s.contracts[address]
	s.mux.RUnlock()

	return value, ok
}

func (s *ActiveContractStore) onActiveContract(ctx context.Context, activeContract *apiv2.ActiveContract) error {
	instanceAddresses, err := getInstanceAddresses(activeContract.GetCreatedEvent())
	if err != nil {
		s.logger.Debug().Err(err).Str("contractID", activeContract.GetCreatedEvent().GetContractId()).Msg("Failed to get instance addresses for active contract, skipping")
		return nil
	}
	s.mux.Lock()
	for _, address := range instanceAddresses {
		s.contracts[address] = activeContract
	}
	s.mux.Unlock()

	return nil
}

func (s *ActiveContractStore) onCreatedEvent(ctx context.Context, transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) error {
	instanceAddresses, err := getInstanceAddresses(createdEvent)
	if err != nil {
		s.logger.Debug().Err(err).Str("contractID", createdEvent.GetContractId()).Msg("Failed to get instance addresses for created event, skipping")
		return nil
	}
	s.mux.Lock()
	for _, address := range instanceAddresses {
		s.contracts[address] = &apiv2.ActiveContract{
			CreatedEvent:        createdEvent,
			SynchronizerId:      transaction.GetSynchronizerId(),
			ReassignmentCounter: 0,
		}
	}
	s.mux.Unlock()

	return nil
}

// getInstanceAddresses returns all possible InstanceAddresses for a given active contract.
// The contract must contain a create argument 'instanceID' of type text.
// Since a contract with multiple signatories can have multiple InstanceAddresses, this returns all possible values.
func getInstanceAddresses(createdEvent *apiv2.CreatedEvent) ([]contracts.InstanceAddress, error) {
	var instanceID string
	for _, field := range createdEvent.GetCreateArguments().GetFields() {
		if field.GetLabel() == "instanceId" {
			instanceID = field.GetValue().GetText()
		}
	}
	if instanceID == "" {
		return nil, fmt.Errorf("no instanceId found in active contract")
	}

	instanceAddresses := make([]contracts.InstanceAddress, len(createdEvent.GetSignatories()))
	for i, s := range createdEvent.GetSignatories() {
		instanceAddresses[i] = contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(s)).InstanceAddress()
	}

	return instanceAddresses, nil
}
