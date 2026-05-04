package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
)

func GetCurrentOffset(ctx context.Context, stateService apiv2.StateServiceClient) (int64, error) {
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get ledger end: %w", err)
	}

	return ledgerEndResp.GetOffset(), nil
}

func GetDisclosedContractById(ctx context.Context, participant canton.Participant, contractId string) (*apiv2.DisclosedContract, error) {
	offset, err := GetCurrentOffset(ctx, participant.LedgerServices.State)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
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

func GetDisclosedContractByTemplateId(ctx context.Context, participant canton.Participant, templateId *apiv2.Identifier) (*apiv2.DisclosedContract, error) {
	activeContracts, err := ListActiveContractsByTemplateId(ctx, participant, templateId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active contracts with templateId %v found on participant %v for party %s", templateId, participant.Name, participant.PartyID)
	}

	contract := activeContracts[len(activeContracts)-1]

	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}

func ListActiveContractsByTemplateId(ctx context.Context, participant canton.Participant, templateId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant.LedgerServices.State)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
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
		return nil, fmt.Errorf("failed to get active contracts for party %v: %w", participant.PartyID, err)
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

func ListActiveContractsByInterfaceId(ctx context.Context, participant canton.Participant, interfaceId *apiv2.Identifier) ([]*apiv2.ActiveContract, error) {
	var activeContracts []*apiv2.ActiveContract
	offset, err := GetCurrentOffset(ctx, participant.LedgerServices.State)
	if err != nil {
		return nil, err
	}
	activeContractsResponse, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
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

// ContractCleanup performs a best-effort attempt at removing all CCIP and MCMS contracts.
// It can be run as a test cleanup function in order to allow multiple tests to run on a participant in sequence, without
// having to re-create the participant(s).
func ContractCleanup(t *testing.T, ctx context.Context, participants []canton.Participant) {
	for _, participant := range participants {
		ArchiveAllContracts(t, ctx, participant)
	}
}

func ArchiveAllContracts(t *testing.T, ctx context.Context, participant canton.Participant) {
	t.Logf("Archiving all active contracts for participant %v and party %s", participant.Name, participant.PartyID)
	offset, err := GetCurrentOffset(ctx, participant.LedgerServices.State)
	require.NoError(t, err)
	resp, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offset,
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: &apiv2.Filters{
					Cumulative: []*apiv2.CumulativeFilter{
						{IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{WildcardFilter: &apiv2.WildcardFilter{}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(ccipreceiver.CCIPReceiver{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(ccipreceiver.CCIPMessageReceived{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(ccipsender.CCIPSender{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(ccvs.CommitteeVerifier{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(onramp.OnRamp{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(offramp.OffRamp{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.RateLimiter{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.GlobalConfig{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.ExecutionStateChanged{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.CCIPMessageSent{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.TokenReceiveTicket{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(common.TokenReceiveTicketClaimed{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(perpartyrouter.PerPartyRouter{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(perpartyrouter.PerPartyRouterFactory{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(tokenadminregistry.TokenAdminRegistry{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(tokenadminregistry.TokenConfig{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(rmn.RMNRemote{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(feequoter.FeeQuoter{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(executor.Executor{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}).ToLedgerIdentifier()}}},
						// {IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{TemplateId: contracts.TemplateIDFromBinding(mcms.MCMS{}).ToLedgerIdentifier()}}},
					},
				},
			},
			FiltersForAnyParty: nil,
			Verbose:            false,
		},
	})
	require.NoError(t, err)
	wg := sync.WaitGroup{}
	for {
		activeContract, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		createdEvent := activeContract.GetActiveContract().GetCreatedEvent()

		if strings.HasPrefix(createdEvent.GetTemplateId().GetModuleName(), "CCIP") ||
			strings.HasPrefix(createdEvent.GetTemplateId().GetModuleName(), "MCMS") ||
			strings.HasPrefix(createdEvent.GetTemplateId().GetModuleName(), "Link") {
			wg.Go(func() {
				t.Logf("Archiving contract %q: %s:%s:%s", createdEvent.GetContractId(), createdEvent.GetTemplateId().GetPackageId(), createdEvent.GetTemplateId().GetModuleName(), createdEvent.GetTemplateId().GetEntityName())

				archiveResp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
					Commands: &apiv2.Commands{
						CommandId: uuid.NewString(),
						Commands: []*apiv2.Command{{
							Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
								TemplateId: createdEvent.GetTemplateId(),
								ContractId: createdEvent.GetContractId(),
								Choice:     "Archive",
								ChoiceArgument: &apiv2.Value{
									Sum: &apiv2.Value_Record{Record: &apiv2.Record{
										RecordId: nil,
										Fields:   nil,
									}},
								},
							}},
						}},
						ActAs: []string{participant.PartyID},
					},
					TransactionFormat: nil,
				})
				if err != nil {
					t.Logf("Failed to archive contract %q: %v", createdEvent.GetContractId(), err)
					return
				}
				t.Logf("Archived Contract %q in update: %v", createdEvent.GetContractId(), archiveResp.GetTransaction().GetUpdateId())
			})
		}
	}
	wg.Wait()
}
