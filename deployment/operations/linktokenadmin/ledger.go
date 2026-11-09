package linktokenadmin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	spliceHolding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

func disclosureOf(ev *apiv2.CreatedEvent, synchronizerID string) *apiv2.DisclosedContract {
	return &apiv2.DisclosedContract{
		TemplateId:       ev.GetTemplateId(),
		ContractId:       ev.GetContractId(),
		CreatedEventBlob: ev.GetCreatedEventBlob(),
		SynchronizerId:   synchronizerID,
	}
}

// findByInstance resolves an active contract signed by owner with the given instanceId.
func findByInstance(ctx context.Context, p canton.Participant, templateID, instanceID, owner string) (*apiv2.ActiveContract, error) {
	addr := contracts.NewRawInstanceAddress(contracts.InstanceID(instanceID), types.PARTY(owner))

	return contract.FindActiveContractByInstanceAddress(ctx, p.LedgerServices.State, contract.LedgerQueryParties(p), templateID, addr.InstanceAddress())
}

func findAuthorization(ctx context.Context, p canton.Participant, spec kindSpec, addr contracts.RawInstanceAddress) (authorization, *apiv2.DisclosedContract, error) {
	ac, err := findByInstance(ctx, p, spec.templateID, addr.InstanceID(), addr.Owner())
	if err != nil {
		return authorization{}, nil, fmt.Errorf("find authorization %s: %w", addr, err)
	}
	auth, err := spec.decode(ac.GetCreatedEvent())
	if err != nil {
		return authorization{}, nil, fmt.Errorf("decode authorization %s: %w", addr, err)
	}

	return auth, disclosureOf(ac.GetCreatedEvent(), ac.GetSynchronizerId()), nil
}

func findAdmin(ctx context.Context, p canton.Participant, auth authorization) (string, *apiv2.DisclosedContract, error) {
	ac, err := findByInstance(ctx, p, tokenadmin.LinkTokenAdmin{}.GetTemplateID(), auth.adminInstanceID, auth.ccipOwner)
	if err != nil {
		return "", nil, fmt.Errorf("find LinkTokenAdmin %s@%s: %w", auth.adminInstanceID, auth.ccipOwner, err)
	}

	return ac.GetCreatedEvent().GetContractId(), disclosureOf(ac.GetCreatedEvent(), ac.GetSynchronizerId()), nil
}

func wildcardFilters(parties []string) map[string]*apiv2.Filters {
	filters := make(map[string]*apiv2.Filters, len(parties))
	for _, party := range parties {
		filters[party] = &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{{
			IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
				WildcardFilter: &apiv2.WildcardFilter{IncludeCreatedEventBlob: true},
			},
		}}}
	}

	return filters
}

// discloseByID fetches the created event blob of an active contract visible to the participant.
func discloseByID(ctx context.Context, p canton.Participant, contractID string) (*apiv2.DisclosedContract, error) {
	resp, err := p.LedgerServices.EventQuery.GetEventsByContractId(ctx, &apiv2.GetEventsByContractIdRequest{
		ContractId:  contractID,
		EventFormat: &apiv2.EventFormat{FiltersByParty: wildcardFilters(contract.LedgerQueryParties(p))},
	})
	if err != nil {
		return nil, fmt.Errorf("get events for contract %s: %w", contractID, err)
	}
	if resp.GetArchived() != nil {
		return nil, fmt.Errorf("contract %s is archived", contractID)
	}
	created := resp.GetCreated()
	if created == nil || len(created.GetCreatedEvent().GetCreatedEventBlob()) == 0 {
		return nil, fmt.Errorf("no created event blob for contract %s", contractID)
	}

	return disclosureOf(created.GetCreatedEvent(), created.GetSynchronizerId()), nil
}

// listHoldings returns the unlocked holdings of owner for instrument, with their disclosures.
func listHoldings(ctx context.Context, p canton.Participant, owner string, instrument spliceHolding.InstrumentId) ([]holding, map[string]*apiv2.DisclosedContract, error) {
	packageID, moduleName, entityName, err := contracts.ParseTemplateIDFromString(spliceHolding.IHoldingInterfaceID())
	if err != nil {
		return nil, nil, err
	}
	interfaceFilter := &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{{
		IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
			InterfaceId:             &apiv2.Identifier{PackageId: packageID, ModuleName: moduleName, EntityName: entityName},
			IncludeInterfaceView:    true,
			IncludeCreatedEventBlob: true,
		}},
	}}}
	parties := contract.LedgerQueryParties(p)
	filters := make(map[string]*apiv2.Filters, len(parties))
	for _, party := range parties {
		filters[party] = interfaceFilter
	}

	end, err := p.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("get ledger end: %w", err)
	}
	stream, err := p.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: end.GetOffset(),
		EventFormat:    &apiv2.EventFormat{FiltersByParty: filters, Verbose: true},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("get active holdings: %w", err)
	}
	defer stream.CloseSend()

	var holdings []holding
	disclosed := map[string]*apiv2.DisclosedContract{}
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("receive active holdings: %w", err)
		}
		ac := resp.GetActiveContract()
		if ac == nil {
			continue
		}
		ev := ac.GetCreatedEvent()
		if _, seen := disclosed[ev.GetContractId()]; seen || len(ev.GetInterfaceViews()) == 0 {
			continue
		}
		// Another token's holding whose view failed to compute must not block selection.
		iv := ev.GetInterfaceViews()[0]
		if iv.GetViewStatus().GetCode() != 0 || iv.GetViewValue() == nil {
			continue
		}
		var view spliceHolding.HoldingView
		if err := ledger.RecordToStruct(iv.GetViewValue(), &view); err != nil {
			return nil, nil, fmt.Errorf("decode holding view %s: %w", ev.GetContractId(), err)
		}
		if string(view.Owner) != owner || view.InstrumentId != instrument || view.Lock != nil {
			continue
		}
		amount, ok := new(big.Rat).SetString(string(view.Amount))
		if !ok {
			return nil, nil, fmt.Errorf("holding %s has invalid amount %q", ev.GetContractId(), view.Amount)
		}
		holdings = append(holdings, holding{contractID: ev.GetContractId(), amount: amount})
		disclosed[ev.GetContractId()] = disclosureOf(ev, ac.GetSynchronizerId())
	}

	return holdings, disclosed, nil
}

// submit exercises cmd as actAs and returns the update id and the choice's result.
func submit(ctx context.Context, p canton.Participant, actAs []string, cmd *model.ExerciseCommand, disclosed []*apiv2.DisclosedContract) (string, *apiv2.Value, error) {
	packageID, moduleName, entityName, err := contracts.ParseTemplateIDFromString(cmd.TemplateID)
	if err != nil {
		return "", nil, fmt.Errorf("parse template id %s: %w", cmd.TemplateID, err)
	}
	filters := make(map[string]*apiv2.Filters, len(actAs))
	for _, party := range actAs {
		filters[party] = &apiv2.Filters{}
	}

	resp, err := p.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			ActAs:     actAs,
			Commands: []*apiv2.Command{{Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
				TemplateId:     &apiv2.Identifier{PackageId: packageID, ModuleName: moduleName, EntityName: entityName},
				ContractId:     cmd.ContractID,
				Choice:         cmd.Choice,
				ChoiceArgument: ledger.MapToValue(cmd.Arguments),
			}}}},
			DisclosedContracts: disclosed,
		},
		// Ledger effects carry the root exercise result, which names created contracts the
		// submitter is not a stakeholder of.
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat:      &apiv2.EventFormat{FiltersByParty: filters},
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		},
	})
	if err != nil {
		return "", nil, fmt.Errorf("submit %s: %w", cmd.Choice, err)
	}
	tx := resp.GetTransaction()
	for _, ev := range tx.GetEvents() {
		if ex := ev.GetExercised(); ex != nil && ex.GetChoice() == cmd.Choice && ex.GetContractId() == cmd.ContractID {
			return tx.GetUpdateId(), ex.GetExerciseResult(), nil
		}
	}

	return tx.GetUpdateId(), nil, nil
}

// dedupDisclosures drops nil entries and repeated contract ids, keeping the first.
func dedupDisclosures(in ...*apiv2.DisclosedContract) []*apiv2.DisclosedContract {
	seen := make(map[string]struct{}, len(in))
	out := make([]*apiv2.DisclosedContract, 0, len(in))
	for _, d := range in {
		if d == nil {
			continue
		}
		if _, ok := seen[d.GetContractId()]; ok {
			continue
		}
		seen[d.GetContractId()] = struct{}{}
		out = append(out, d)
	}

	return out
}
