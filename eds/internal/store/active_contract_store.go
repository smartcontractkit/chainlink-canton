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

type ActiveContractStore Store[contracts.InstanceAddress, *apiv2.ActiveContract]

// RegisteredTemplate defines a template that the ActiveContractStore will keep track of.
// It defines a TemplateID, as well as a Party - the latter of which must be a stakeholder (signatory or observer)
// on the contract in order for the ActiveContractStore to pick up the contract.
type RegisteredTemplate struct {
	TemplateID contracts.TemplateID
	PartyID    string
}

type ActiveContractStoreConfig struct {
	Logger        zerolog.Logger
	UpdateService apiv2.UpdateServiceClient
	StateService  apiv2.StateServiceClient
	MaxRetries    int
}

// NewActiveContractStore returns a Store implementation that keeps track of active contracts by subscribing
// to incremental ledger updates.
// The ActiveContractStore is configured with a variable list of registeredTemplates which are the only templates is will
// keep track of. All templates must contain an 'InstanceId' field in order to calculate their InstanceAddress.
// For example, if configured with:
//
//	store.RegisteredTemplate{
//	   TemplateID: store.TemplateID{
//	       PackageID:  "#ccip-committeeverifier",
//	       ModuleName: "CCIP.CommitteeVerifier",
//	       EntityName: "CommitteeVerifier",
//	   },
//	   PartyID: "ccvOwnerParty::0x123567890",
//	}
//
// The ActiveContractStore will list and subscribe all CommitteeVerifier contracts that the 'ccvOwnerParty' can see.
// It will then index them by their calculated InstanceAddress using the combination of signatory + instanceId field.
//
// NewActiveContractStore itself will perform no RPC calls, it will immediately return.
// In order for the ActiveContractStore to initialize and subscribe to updates, (s *ContractStore) Run() needs to be run.
func NewActiveContractStore(
	_ context.Context,
	config ActiveContractStoreConfig,
	metrics Metrics,
	registeredTemplates ...RegisteredTemplate,
) (*ContractStore[contracts.InstanceAddress, *apiv2.ActiveContract], error) {
	filtersByParty := make(map[string]*apiv2.Filters) // Assemble filters
	for _, template := range registeredTemplates {
		existingFilterForParty, ok := filtersByParty[template.PartyID]
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
		filtersByParty[template.PartyID] = existingFilterForParty
	}

	return &ContractStore[contracts.InstanceAddress, *apiv2.ActiveContract]{
		logger:         config.Logger.With().Str("component", "ActiveContractStore").Logger(),
		updateService:  config.UpdateService,
		stateService:   config.StateService,
		metrics:        metrics,
		maxRetries:     config.MaxRetries,
		filtersByParty: filtersByParty,
		mux:            sync.RWMutex{},
		ledgerEnd:      0,
		contracts:      make(map[contracts.InstanceAddress]*apiv2.ActiveContract),
		handleActiveContract: func(ctx context.Context, store *ContractStore[contracts.InstanceAddress, *apiv2.ActiveContract], activeContract *apiv2.ActiveContract) (updates []ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract], err error) {
			instanceAddresses, err := getInstanceAddresses(activeContract.GetCreatedEvent())
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", activeContract.GetCreatedEvent().GetContractId()).Msg("Failed to get instance addresses for active contract, skipping")
				return nil, nil
			}

			updates = make([]ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract], len(instanceAddresses))
			for i, address := range instanceAddresses {
				updates[i] = ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract]{
					Key:   address,
					Value: activeContract,
				}
			}

			return updates, nil
		},
		handleCreatedEvent: func(ctx context.Context, store *ContractStore[contracts.InstanceAddress, *apiv2.ActiveContract], transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) (updates []ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract], err error) {
			instanceAddresses, err := getInstanceAddresses(createdEvent)
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", createdEvent.GetContractId()).Msg("Failed to get instance addresses for created contract, skipping")
				return nil, nil
			}

			updates = make([]ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract], len(instanceAddresses))
			for i, address := range instanceAddresses {
				updates[i] = ContractUpdate[contracts.InstanceAddress, *apiv2.ActiveContract]{
					Key: address,
					Value: &apiv2.ActiveContract{
						CreatedEvent:        createdEvent,
						SynchronizerId:      transaction.GetSynchronizerId(),
						ReassignmentCounter: 0,
					},
				}
			}

			return updates, nil
		},
	}, nil
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
