package store

import (
	"context"
	"fmt"
	"sync"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type ActiveContractStore struct {
	logger zerolog.Logger

	stream  *BackfilledStream
	filters FiltersByParty

	mux                        sync.RWMutex
	contractsByInstanceAddress map[contracts.InstanceAddress]*apiv2.ActiveContract
	contractsByTemplateId      map[types.PARTY]map[contracts.TemplateID]*apiv2.ActiveContract
	contractsByContractId      map[types.CONTRACT_ID]*apiv2.ActiveContract
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

func (s *ActiveContractStore) Run(ctx context.Context, streamConfig StreamConfig, opts ...RunOption) error {
	options := &runOptions{}
	for _, opt := range opts {
		opt(options)
	}

	s.mux.Lock()
	s.contractsByInstanceAddress = make(map[contracts.InstanceAddress]*apiv2.ActiveContract)
	s.contractsByTemplateId = make(map[types.PARTY]map[contracts.TemplateID]*apiv2.ActiveContract)
	s.contractsByContractId = make(map[types.CONTRACT_ID]*apiv2.ActiveContract)
	s.mux.Unlock()

	return s.stream.Run(ctx, s.filters, streamConfig, s.onActiveContract, s.onCreatedEvent, s.onArchivedEvent, nil, options.onBackfillCompleted)
}

func (s *ActiveContractStore) Get(address contracts.InstanceAddress) (*apiv2.ActiveContract, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	value, ok := s.contractsByInstanceAddress[address]

	// Return a copy of the ActiveContract
	return proto.CloneOf(value), ok
}

// GetByTemplateId returns the active contract for the given TemplateID, if it exists.
// It accepts TemplateIds in both the #packageName and PackageId syntax.
func (s *ActiveContractStore) GetByTemplateId(party types.PARTY, templateId contracts.TemplateID) (*apiv2.ActiveContract, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	contractsForParty, ok := s.contractsByTemplateId[party]
	if !ok {
		return nil, false
	}
	value, ok := contractsForParty[templateId]
	if !ok {
		return nil, false
	}

	// Return a copy of the ActiveContract
	return proto.CloneOf(value), true
}

func (s *ActiveContractStore) GetByContractId(contractId types.CONTRACT_ID) (*apiv2.ActiveContract, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	value, ok := s.contractsByContractId[contractId]

	// Return a copy of the ActiveContract
	return proto.CloneOf(value), ok
}

func (s *ActiveContractStore) onActiveContract(ctx context.Context, activeContract *apiv2.ActiveContract) error {
	// Create a copy of the ActiveContract, since it will be stored in state
	activeContract = proto.CloneOf(activeContract)

	instanceAddresses, _ := getInstanceAddresses(activeContract.GetCreatedEvent())

	s.mux.Lock()
	// Add by ContractId
	s.contractsByContractId[types.CONTRACT_ID(activeContract.GetCreatedEvent().GetContractId())] = activeContract
	// Add by InstanceAddress
	for _, address := range instanceAddresses {
		s.contractsByInstanceAddress[address] = activeContract
	}
	for _, witnessParty := range activeContract.GetCreatedEvent().GetWitnessParties() {
		party := types.PARTY(witnessParty)
		contractsForParty, ok := s.contractsByTemplateId[party]
		if !ok {
			contractsForParty = make(map[contracts.TemplateID]*apiv2.ActiveContract)
			s.contractsByTemplateId[party] = contractsForParty
		}
		// Add by PackageId
		contractsForParty[contracts.TemplateID{
			PackageID:  activeContract.GetCreatedEvent().GetTemplateId().GetPackageId(),
			ModuleName: activeContract.GetCreatedEvent().GetTemplateId().GetModuleName(),
			EntityName: activeContract.GetCreatedEvent().GetTemplateId().GetEntityName(),
		}] = activeContract
		// Add by PackageName
		contractsForParty[contracts.TemplateID{
			PackageID:  fmt.Sprintf("#%s", activeContract.GetCreatedEvent().GetPackageName()),
			ModuleName: activeContract.GetCreatedEvent().GetTemplateId().GetModuleName(),
			EntityName: activeContract.GetCreatedEvent().GetTemplateId().GetEntityName(),
		}] = activeContract
	}
	s.mux.Unlock()

	return nil
}

func (s *ActiveContractStore) onCreatedEvent(ctx context.Context, transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) error {
	instanceAddresses, _ := getInstanceAddresses(createdEvent)
	activeContract := &apiv2.ActiveContract{
		CreatedEvent:        proto.CloneOf(createdEvent), // Copy the CreatedEvent
		SynchronizerId:      transaction.GetSynchronizerId(),
		ReassignmentCounter: 0,
	}

	s.mux.Lock()
	// Add by ContractId
	s.contractsByContractId[types.CONTRACT_ID(createdEvent.GetContractId())] = activeContract
	// Add by InstanceAddress
	for _, address := range instanceAddresses {
		s.contractsByInstanceAddress[address] = activeContract
	}
	for _, witnessParty := range createdEvent.GetWitnessParties() {
		party := types.PARTY(witnessParty)
		contractsForParty, ok := s.contractsByTemplateId[party]
		if !ok {
			contractsForParty = make(map[contracts.TemplateID]*apiv2.ActiveContract)
			s.contractsByTemplateId[party] = contractsForParty
		}
		// Add by PackageId
		contractsForParty[contracts.TemplateID{
			PackageID:  createdEvent.GetTemplateId().GetPackageId(),
			ModuleName: createdEvent.GetTemplateId().GetModuleName(),
			EntityName: createdEvent.GetTemplateId().GetEntityName(),
		}] = activeContract
		// Add by PackageName
		contractsForParty[contracts.TemplateID{
			PackageID:  fmt.Sprintf("#%s", createdEvent.GetPackageName()),
			ModuleName: createdEvent.GetTemplateId().GetModuleName(),
			EntityName: createdEvent.GetTemplateId().GetEntityName(),
		}] = activeContract
	}
	s.mux.Unlock()

	return nil
}

func (s *ActiveContractStore) onArchivedEvent(ctx context.Context, _ *apiv2.Transaction, archivedEvent *apiv2.ArchivedEvent) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	activeContract, ok := s.contractsByContractId[types.CONTRACT_ID(archivedEvent.GetContractId())]
	if !ok {
		return fmt.Errorf("archived event for contract %v with no active contract record", archivedEvent.GetContractId())
	}

	// Delete by ContractId
	delete(s.contractsByContractId, types.CONTRACT_ID(archivedEvent.GetContractId()))
	// Delete by InstanceAddress
	instanceAddresses, _ := getInstanceAddresses(activeContract.GetCreatedEvent())
	for _, address := range instanceAddresses {
		delete(s.contractsByInstanceAddress, address)
	}
	for _, witnessParty := range activeContract.GetCreatedEvent().GetWitnessParties() {
		party := types.PARTY(witnessParty)
		contractsForParty, ok := s.contractsByTemplateId[party]
		if !ok {
			continue
		}
		// Delete by PackageId
		delete(contractsForParty, contracts.TemplateID{
			PackageID:  archivedEvent.GetTemplateId().GetPackageId(),
			ModuleName: archivedEvent.GetTemplateId().GetModuleName(),
			EntityName: archivedEvent.GetTemplateId().GetEntityName(),
		})
		// Delete by PackageName
		delete(contractsForParty, contracts.TemplateID{
			PackageID:  fmt.Sprintf("#%s", archivedEvent.GetPackageName()),
			ModuleName: archivedEvent.GetTemplateId().GetModuleName(),
			EntityName: archivedEvent.GetTemplateId().GetEntityName(),
		})
	}

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
