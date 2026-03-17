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

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

// ContractStore provides access to currently active contracts on the ledger, indexed by InstanceAddress.
type ContractStore interface {
	// GetContract returns the currently active contract for a given InstanceAddress.
	// Returns nil if no active contract is found.
	GetContract(ctx context.Context, instanceAddress contracts.InstanceAddress) *apiv2.ActiveContract
}

var _ ContractStore = &UpdateStore{}

type UpdateStore struct {
	logger        zerolog.Logger
	updateService apiv2.UpdateServiceClient
	stateService  apiv2.StateServiceClient

	metrics Metrics

	// The filters to apply when backfilling/subscribing to updates
	filtersByParty map[string]*apiv2.Filters

	mux       sync.RWMutex
	ledgerEnd int64
	contracts map[contracts.InstanceAddress]*apiv2.ActiveContract

	maxRetries int
}

// RegisteredTemplate defines a template that the UpdateStore will keep track of.
// It defines a TemplateID, as well as a Party - the latter of which must be a stakeholder (signatory or observer)
// on the contract in order for the UpdateStore to pick up the contract.
type RegisteredTemplate struct {
	TemplateID contracts.TemplateID
	PartyID    string
}

type UpdateStoreConfig struct {
	Logger        zerolog.Logger
	UpdateService apiv2.UpdateServiceClient
	StateService  apiv2.StateServiceClient
	MaxRetries    int
}

// NewUpdateStore returns a ContractStore implementation that keeps track of active contracts by subscribing
// to incremental ledger updates.
// The UpdateStore is configured with a variable list of registeredTemplates which are the only templates is will
// keep track of. All templates must contain an 'InstanceId' field in order to calculate their InstanceAddress.
// For example, if configured with:
//
//	store.RegisteredTemplate{
//	   TemplateID: store.TemplateID{
//	       PackageID:  "#ccip-committeeverifier",
//	       ModuleName: "CCIP.CommitteeVerifier",
//	       EntityName: "CommitteeVerifier",
//	   },
//	   PartyID: "ccvOwerParty::0x123567890",
//	}
//
// The UpdateStore will list and subscribe all CommitteeVerifier contracts that the 'ccvOwnerParty' can see.
// It will then index them by their calculated InstanceAddress using the combination of signatory + instanceId field.
//
// NewUpdateStore itself will perform not RPC calls, it will immediately return.
// In order for the UpdateStore to initialize and subscribe to updates, (s *UpdateStore)Run() needs to be run.
func NewUpdateStore(
	ctx context.Context,
	config UpdateStoreConfig,
	metrics Metrics,
	registeredTemplates ...RegisteredTemplate,
) (*UpdateStore, error) {
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

	return &UpdateStore{
		logger:         config.Logger.With().Str("component", "UpdateStore").Logger(),
		updateService:  config.UpdateService,
		stateService:   config.StateService,
		metrics:        metrics,
		mux:            sync.RWMutex{},
		ledgerEnd:      0,
		filtersByParty: filtersByParty,
		contracts:      make(map[contracts.InstanceAddress]*apiv2.ActiveContract),
		maxRetries:     config.MaxRetries,
	}, nil
}

// Run runs the UpdateStore. It will start by initializing the UpdateStore with a backfill of all existing active contracts
// and subscribe to incremental updates afterward.
// Run is a long-running function that needs to keep running int the background in order for the UpdateStore to keep
// up-tp-date.
// To terminate Run, cancel the context.
func (s *UpdateStore) Run(ctx context.Context) error {
	s.logger.Debug().Msg("Starting UpdateStore")
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
	// Checking if the metric is increasing can be used to verify that the UpdateStore is still subscribed and ready to process updates
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
				s.logger.Debug().Msg("Context cancelled, stopping UpdateStore")
				_ = stream.CloseSend()

				return ctx.Err()
			case <-uptimeTicker.C:
				// Increment uptime metric on every tick
				s.metrics.IncrementStoreSubscriptionUptime(ctx)
			case resp := <-updateChan:
				if transaction, ok := resp.GetUpdate().(*apiv2.GetUpdatesResponse_Transaction); ok {
					s.logger.Trace().Str("updateID", transaction.Transaction.GetUpdateId()).Int("events", len(transaction.Transaction.GetEvents())).Msg("Update received")
					for _, event := range transaction.Transaction.GetEvents() {
						if createdEvent, ok := event.GetEvent().(*apiv2.Event_Created); ok {
							s.logger.Trace().
								Str("contractID", createdEvent.Created.GetContractId()).
								Str("packageName", createdEvent.Created.GetPackageName()).
								Str("packageID", createdEvent.Created.GetTemplateId().GetPackageId()).
								Str("moduleName", createdEvent.Created.GetTemplateId().GetModuleName()).
								Str("entityName", createdEvent.Created.GetTemplateId().GetEntityName()).
								Msg("Received CreatedEvent")
							instanceAddresses, err := getInstanceAddresses(createdEvent.Created)
							if err != nil {
								s.logger.Debug().Err(err).Str("contractID", createdEvent.Created.GetContractId()).Msg("Failed to get instance addresses for created contract, skipping")
								continue
							}
							// Increment updates metric
							s.metrics.IncrementStoreUpdatesCounter(ctx)
							for _, address := range instanceAddresses {
								s.logger.Trace().Stringer("instanceAddress", address).Msg("Updating active contract")
								s.mux.Lock()
								s.contracts[address] = &apiv2.ActiveContract{
									CreatedEvent:        createdEvent.Created,
									SynchronizerId:      transaction.Transaction.GetSynchronizerId(),
									ReassignmentCounter: 0,
								}
								// Update ledgerEnd
								s.ledgerEnd = transaction.Transaction.GetOffset()
								s.mux.Unlock()
								// Record latest ledger end
								s.metrics.RecordStoreLedgerEndGauge(ctx, transaction.Transaction.GetOffset())
							}
						}
					}
				}
			}
		}
	}
}

// backfill returns all currently-active contracts for s.filtersByParty at a given offset, indexed by InstanceAddress.
func (s *UpdateStore) backfill(ctx context.Context, offset int64) (map[contracts.InstanceAddress]*apiv2.ActiveContract, error) {
	activeContracts := make(map[contracts.InstanceAddress]*apiv2.ActiveContract)

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

			if ac, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
				s.logger.Trace().
					Str("contractID", ac.ActiveContract.GetCreatedEvent().GetContractId()).
					Str("packageName", ac.ActiveContract.GetCreatedEvent().GetPackageName()).
					Str("packageID", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetPackageId()).
					Str("moduleName", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetModuleName()).
					Str("entityName", ac.ActiveContract.GetCreatedEvent().GetTemplateId().GetEntityName()).
					Msg("Backfilling active contract")
				instanceAddresses, err := getInstanceAddresses(ac.ActiveContract.GetCreatedEvent())
				if err != nil {
					s.logger.Debug().Err(err).Str("contractID", ac.ActiveContract.GetCreatedEvent().GetContractId()).Msg("Failed to get instance addresses for active contract, skipping")
					continue
				}
				for _, address := range instanceAddresses {
					s.logger.Trace().Stringer("instanceAddress", address).Msg("Updating active contract")
					activeContracts[address] = ac.ActiveContract
				}
			}
		}
	}
}

// updateStreamFactory returns a StreamFactory that creates GetUpdates streams for this store's filters.
func (s *UpdateStore) updateStreamFactory() StreamFactory[apiv2.GetUpdatesResponse] {
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

// GetContract returns the active contract for the given InstanceAddress. If no active contract is found, it returns nil.
func (s *UpdateStore) GetContract(ctx context.Context, instanceAddress contracts.InstanceAddress) *apiv2.ActiveContract {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.contracts[instanceAddress]
}
