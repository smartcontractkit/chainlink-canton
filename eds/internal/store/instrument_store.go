package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
)

var holdingInterfaceId = &apiv2.Identifier{
	// PackageId seems to be the hashed package ID, not the human-readable one.
	// We probably shouldn't check against a specific hash, so we ignore that part of the interface ID.
	// PackageId:  fmt.Sprintf("#%s", splice_api_token_holding_v1.PackageName),
	ModuleName: "Splice.Api.Token.HoldingV1",
	EntityName: "Holding",
}

type InstrumentHoldingStore Store[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract]

// defaultReconnectBackoff is the delay before reconnecting after the update stream closes (e.g. server closed or error).
const defaultReconnectBackoff = 5 * time.Second

type InstrumentHoldingStoreConfig struct {
	Logger           zerolog.Logger
	Owner            types.PARTY
	StateService     apiv2.StateServiceClient
	UpdateService    apiv2.UpdateServiceClient
	MaxRetries       int
	ReconnectBackoff time.Duration // delay before reconnecting after stream close; 0 uses DefaultReconnectBackoff
}

func NewInstrumentHoldingStore(
	config InstrumentHoldingStoreConfig,
	metrics Metrics,
) *ContractStore[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract] {
	backoff := config.ReconnectBackoff
	if backoff == 0 {
		backoff = defaultReconnectBackoff
	}
	owner := config.Owner

	return &ContractStore[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract]{
		logger:           config.Logger.With().Str("component", "InstrumentHoldingStore").Logger(),
		updateService:    config.UpdateService,
		stateService:     config.StateService,
		metrics:          metrics,
		maxRetries:       config.MaxRetries,
		reconnectBackoff: backoff,
		filtersByParty: map[string]*apiv2.Filters{
			string(config.Owner): {
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
		mux:       sync.RWMutex{},
		ledgerEnd: 0,
		contracts: make(map[splice_api_token_holding_v1.InstrumentId]*apiv2.DisclosedContract),
		handleActiveContract: func(ctx context.Context, store *ContractStore[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract], activeContract *apiv2.ActiveContract) (updates []ContractUpdate[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract], err error) {
			relevantView, err := getRelevantInterfaceViewValue(activeContract.GetCreatedEvent().GetInterfaceViews(), holdingInterfaceId)
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", activeContract.GetCreatedEvent().GetContractId()).Msg("Failed to get relevant interface view value, skipping")
				return nil, nil
			}
			holdingView, err := getHoldingView(owner, relevantView)
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", activeContract.GetCreatedEvent().GetContractId()).Msg("Failed to get holding disclosure for active contract, skipping")
				return nil, nil
			}

			// create the disclosure
			holdingDisclosure := &apiv2.DisclosedContract{
				TemplateId:       activeContract.GetCreatedEvent().GetTemplateId(),
				ContractId:       activeContract.GetCreatedEvent().GetContractId(),
				CreatedEventBlob: activeContract.GetCreatedEvent().GetCreatedEventBlob(),
				SynchronizerId:   activeContract.GetSynchronizerId(),
			}
			store.logger.Info().
				Any("holdingView", holdingView).
				Any("holdingDisclosure", holdingDisclosure).
				Msg("Recording holding disclosure")

			return []ContractUpdate[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract]{{
				Key:   holdingView.InstrumentId,
				Value: holdingDisclosure,
			}}, nil
		},
		handleCreatedEvent: func(ctx context.Context, store *ContractStore[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract], transaction *apiv2.Transaction, createdEvent *apiv2.CreatedEvent) (updates []ContractUpdate[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract], err error) {
			relevantView, err := getRelevantInterfaceViewValue(createdEvent.GetInterfaceViews(), holdingInterfaceId)
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", createdEvent.GetContractId()).Msg("Failed to get relevant interface view value, skipping")
				return nil, nil
			}

			holdingView, err := getHoldingView(owner, relevantView)
			if err != nil {
				store.logger.Debug().Err(err).Str("contractID", createdEvent.GetContractId()).Msg("Failed to get holding disclosure for created event, skipping")
				return nil, nil
			}

			// create the disclosure
			holdingDisclosure := &apiv2.DisclosedContract{
				TemplateId:       createdEvent.GetTemplateId(),
				ContractId:       createdEvent.GetContractId(),
				CreatedEventBlob: createdEvent.GetCreatedEventBlob(),
				SynchronizerId:   transaction.GetSynchronizerId(),
			}
			store.logger.Info().
				Any("holdingView", holdingView).
				Any("holdingDisclosure", holdingDisclosure).
				Msg("Recording holding disclosure")

			return []ContractUpdate[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract]{{
				Key:   holdingView.InstrumentId,
				Value: holdingDisclosure,
			}}, nil
		},
	}
}

func getHoldingView(expectedOwner types.PARTY, record *apiv2.Record) (*splice_api_token_holding_v1.HoldingView, error) {
	var holdingView splice_api_token_holding_v1.HoldingView
	err := ledger.RecordToStruct(record, &holdingView)
	if err != nil {
		return nil, fmt.Errorf("failed to convert record to holding view: %w", err)
	}

	// At this point we have an active holding contract, we should form the disclosure for it.
	if holdingView.Owner != expectedOwner {
		return nil, fmt.Errorf("holding owner %v does not match expected owner %v", holdingView.Owner, expectedOwner)
	}

	return &holdingView, nil
}

func getRelevantInterfaceViewValue(interfaceViews []*apiv2.InterfaceView, expectedInterfaceId *apiv2.Identifier) (*apiv2.Record, error) {
	for _, interfaceView := range interfaceViews {
		// PackageId seems to be the hashed package ID, not the human-readable one.
		// We probably shouldn't check against a specific hash, so we ignore that part of the interface ID.
		if interfaceView.GetInterfaceId().GetModuleName() == expectedInterfaceId.GetModuleName() &&
			interfaceView.GetInterfaceId().GetEntityName() == expectedInterfaceId.GetEntityName() {
			return interfaceView.GetViewValue(), nil
		}
	}

	return nil, fmt.Errorf("no interface view found for interface id: %s", expectedInterfaceId.String())
}
