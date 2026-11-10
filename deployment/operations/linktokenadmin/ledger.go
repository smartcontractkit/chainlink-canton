package linktokenadmin

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
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
	ac, err := findByInstance(ctx, p, spec.template.GetTemplateID(), addr.InstanceID(), addr.Owner())
	if err != nil {
		return authorization{}, nil, fmt.Errorf("find authorization %s: %w", addr, err)
	}
	if spec.decode == nil {
		return authorization{}, nil, fmt.Errorf("authorization %s: kind has no decoder", addr)
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

func discloseAll(ctx context.Context, p canton.Participant, contractIDs []string) ([]*apiv2.DisclosedContract, error) {
	disclosed := make([]*apiv2.DisclosedContract, 0, len(contractIDs))
	for _, cid := range contractIDs {
		d, err := discloseByID(ctx, p, cid)
		if err != nil {
			return nil, err
		}
		disclosed = append(disclosed, d)
	}

	return disclosed, nil
}

var executeMintChoice = tokenadmin.LinkMintAuthorization{}.ExecuteMint("", tokenadmin.ExecuteMint{}).Choice

// mintedHoldingCids reads the ExecuteMint result from the update. Ledger effects carry the
// root exercise result, which names holdings the submitter is not a stakeholder of.
func mintedHoldingCids(ctx context.Context, p canton.Participant, updateID, authCid string) ([]string, error) {
	resp, err := p.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: updateID,
		UpdateFormat: &apiv2.UpdateFormat{IncludeTransactions: &apiv2.TransactionFormat{
			EventFormat:      &apiv2.EventFormat{FiltersByParty: map[string]*apiv2.Filters{p.PartyID: {}}},
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		}},
	})
	if err != nil {
		return nil, fmt.Errorf("get update %s: %w", updateID, err)
	}
	for _, ev := range resp.GetTransaction().GetEvents() {
		ex := ev.GetExercised()
		if ex == nil || ex.GetChoice() != executeMintChoice || ex.GetContractId() != authCid {
			continue
		}
		var cids []string
		for _, v := range ex.GetExerciseResult().GetList().GetElements() {
			cids = append(cids, v.GetContractId())
		}

		return cids, nil
	}

	return nil, fmt.Errorf("update %s has no ExecuteMint on %s", updateID, authCid)
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
