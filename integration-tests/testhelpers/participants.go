package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
)

type Participant struct {
	Name string
	// GetToken returns a JWT token to authenticate API calls against this participant
	GetToken func(ctx context.Context) (string, error)
	// GetConfig returns the participant's configuration
	GetConfig func() ParticipantConfig
	UserName  string
	Party     string

	// API Clients

	// Admin API
	PackageServiceClient participantv30.PackageServiceClient

	// Ledger API
	PartyManagementServiceClient admin.PartyManagementServiceClient
	UserManagementServiceClient  admin.UserManagementServiceClient

	StateServiceClient   apiv2.StateServiceClient
	CommandServiceClient apiv2.CommandServiceClient
	UpdateServiceClient  apiv2.UpdateServiceClient
	VersionServiceClient apiv2.VersionServiceClient

	// Validator API
	ScanProxyClient scanProxy.ClientWithResponsesInterface
}

func GetCurrentOffset(ctx context.Context, participant Participant) (int64, error) {
	ledgerEndResp, err := participant.StateServiceClient.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger end: %w", err)
	}

	return ledgerEndResp.GetOffset(), nil
}

func GetDisclosedContractById(ctx context.Context, participant Participant, contractId string) (*apiv2.DisclosedContract, error) {
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.Party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
								WildcardFilter: &apiv2.WildcardFilter{
									IncludeCreatedEventBlob: true,
								},
							},
						},
					},
				},
			},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			if c.ActiveContract.GetCreatedEvent().ContractId == contractId {
				return &apiv2.DisclosedContract{
					TemplateId:       c.ActiveContract.GetCreatedEvent().GetTemplateId(),
					ContractId:       c.ActiveContract.GetCreatedEvent().GetContractId(),
					CreatedEventBlob: c.ActiveContract.GetCreatedEvent().GetCreatedEventBlob(),
					SynchronizerId:   c.ActiveContract.GetSynchronizerId(),
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to find active contract with id %s", contractId)
}

func GetDisclosedContractByTemplateId(ctx context.Context, participant Participant, templateId *apiv2.Identifier) (*apiv2.DisclosedContract, error) {
	activeContracts, err := ListActiveContractsByTemplateId(ctx, participant, templateId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active contracts with templateId %v found on participant %vfor party %s", templateId, participant.Name, participant.Party)
	}

	contract := activeContracts[len(activeContracts)-1]

	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}

func ListActiveContractsByTemplateId(ctx context.Context, participant Participant, templateId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.Party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
								TemplateId:              templateId,
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts for party %v: %w", participant.Party, err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			activeContracts = append(activeContracts, c.ActiveContract)
		}
	}
	slices.SortFunc(activeContracts, func(a, b *apiv2.ActiveContract) int {
		return a.GetCreatedEvent().GetCreatedAt().AsTime().Compare(b.GetCreatedEvent().GetCreatedAt().AsTime())
	})

	return activeContracts, nil
}

func ListActiveContractsByInterfaceId(ctx context.Context, participant Participant, interfaceId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.StateServiceClient.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.Party: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
								InterfaceId:             interfaceId,
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
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer activeContractsResponse.CloseSend()
	for {
		activeContract, err := activeContractsResponse.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract); ok {
			activeContracts = append(activeContracts, c.ActiveContract)
		}
	}
	slices.SortFunc(activeContracts, func(a, b *apiv2.ActiveContract) int {
		return a.GetCreatedEvent().GetCreatedAt().AsTime().Compare(b.GetCreatedEvent().GetCreatedAt().AsTime())
	})

	return activeContracts, nil
}
