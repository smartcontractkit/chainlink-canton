package view

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
)

// ContractMetaData contains metadata about a DAML contract
type ContractMetaData struct {
	Address    string `json:"address"`    // Contract ID
	Owner      string `json:"owner"`      // Party that owns the contract
	TemplateID string `json:"templateId"` // Full template ID
	InstanceID string `json:"instanceId"` // Instance ID from contract
}

// OnRampView represents the view of an OnRamp contract
type OnRampView struct {
	ContractMetaData

	CcipOwner  string `json:"ccipOwner"`
	InstanceId string `json:"instanceId"`
}

// GenerateOnRampView generates an OnRamp view by querying the on-chain state
// stateService is the StateService from bindingClient.StateService (uses channels with model types)
func GenerateOnRampView(
	ctx context.Context,
	stateService ledger.StateService,
	onRampContractID string,
	ccipOnRampPkgID string,
	party string, // Party that can read the contract
) (OnRampView, error) {
	if onRampContractID == "" {
		return OnRampView{}, fmt.Errorf("onRampContractID cannot be empty")
	}
	if party == "" {
		return OnRampView{}, fmt.Errorf("party cannot be empty")
	}

	// Get current offset
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get ledger end: %w", err)
	}

	// Query active contracts for the party using wildcard filter
	// Following the pattern from GetActiveContractsForPartyWithFilters
	incl := &model.InclusiveFilters{
		TemplateFilters: []*model.TemplateFilter{
			{
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipOnRampPkgID, "CCIP.OnRamp", "OnRamp"),
			},
		}, // Empty = wildcard
	}

	req := &model.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.Offset,
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {Inclusive: incl},
			},
			Verbose: true,
		},
	}

	responseChan, errorChan := stateService.GetActiveContracts(ctx, req)

	// Find the contract with matching contract ID
	var contractEvent *model.CreatedEvent
	var templateID string
	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				goto done
			}
			if resp == nil || resp.ContractEntry == nil {
				continue
			}
			entry, ok := resp.ContractEntry.(*model.ActiveContractEntry)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.CreatedEvent == nil {
				continue
			}
			createdEvent := entry.ActiveContract.CreatedEvent
			if createdEvent.ContractID == onRampContractID {
				contractEvent = createdEvent
				templateID = createdEvent.TemplateID
				goto done
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			if err != nil {
				return OnRampView{}, fmt.Errorf("failed to receive active contract: %w", err)
			}
		case <-ctx.Done():
			return OnRampView{}, ctx.Err()
		}
	}
done:

	if contractEvent == nil {
		return OnRampView{}, fmt.Errorf("contract with ID %s not found for party %s", onRampContractID, party)
	}

	// CreateArguments is a Record with json dump like below
	// {
	// 	"fields": [
	// 	  {"label":"ccipOwner","value":{"Sum":{"Party":"..."}}},
	// 	  {"label":"instanceId","value":{"Sum":{"Text":"..."}}}
	// 	]
	//   }

	// Treat the interface{} as the Ledger API Record type so I can traverse its fields.
	rec, ok := contractEvent.CreateArguments.(*apiv2.Record)
	if !ok || rec == nil {
		return OnRampView{}, fmt.Errorf("CreateArguments is %T, expected *apiv2.Record", contractEvent.CreateArguments)
	}

	// lookup function inside the fields
	ccipOwner := PartyVal(FieldVal(rec, "ccipOwner"))
	instanceId := TextVal(FieldVal(rec, "instanceId"))

	return OnRampView{
		ContractMetaData: ContractMetaData{
			Address:    onRampContractID,
			Owner:      ccipOwner,
			TemplateID: templateID,
			InstanceID: instanceId,
		},
		CcipOwner:  ccipOwner,
		InstanceId: instanceId,
	}, nil
}
