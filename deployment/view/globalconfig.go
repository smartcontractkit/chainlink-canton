package view

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/service/ledger"
)

// DestChainConfigView represents a destination chain configuration
type DestChainConfigView struct {
	IsEnabled        bool     `json:"isEnabled"`
	DefaultExecutor  string   `json:"defaultExecutor"`
	OffRampAddress   string   `json:"offRampAddress"`
	LaneMandatedCCVs []string `json:"laneMandatedCCVs"`
	DefaultCCVs      []string `json:"defaultCCVs"`
}

// SourceChainConfigView represents a source chain configuration
type SourceChainConfigView struct {
	IsEnabled        bool     `json:"isEnabled"`
	OnRampAddress    string   `json:"onRampAddress"`
	LaneMandatedCCVs []string `json:"laneMandatedCCVs"`
	DefaultCCVs      []string `json:"defaultCCVs"`
}

// GlobalConfigView represents the view of a GlobalConfig contract
type GlobalConfigView struct {
	ContractMetaData

	CcipOwner          string                           `json:"ccipOwner"`
	InstanceId         string                           `json:"instanceId"`
	ChainSelector      string                           `json:"chainSelector"` // Numeric as string
	OnRampAddress      string                           `json:"onRampAddress"`
	DestChainConfigs   map[string]DestChainConfigView   `json:"destChainConfigs"`   // Chain selector (as string) -> config
	SourceChainConfigs map[string]SourceChainConfigView `json:"sourceChainConfigs"` // Chain selector (as string) -> config
}

// GenerateGlobalConfigView generates a GlobalConfig view by querying the on-chain state
// stateService is the StateService from bindingClient.StateService
func GenerateGlobalConfigView(
	ctx context.Context,
	stateService ledger.StateService,
	globalConfigContractID string,
	ccipCommonPkgID string,
	party string, // Party that can read the contract
) (GlobalConfigView, error) {
	if globalConfigContractID == "" {
		return GlobalConfigView{}, fmt.Errorf("globalConfigContractID cannot be empty")
	}
	if party == "" {
		return GlobalConfigView{}, fmt.Errorf("party cannot be empty")
	}

	// Get current offset
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return GlobalConfigView{}, fmt.Errorf("failed to get ledger end: %w", err)
	}

	// Query active contracts for the party using wildcard filter
	// Following the pattern from GetActiveContractsForPartyWithFilters
	incl := &model.InclusiveFilters{
		TemplateFilters: []*model.TemplateFilter{
			{
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipCommonPkgID, "CCIP.GlobalConfig", "GlobalConfig"),
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
			if createdEvent.ContractID == globalConfigContractID {
				contractEvent = createdEvent
				templateID = createdEvent.TemplateID
				goto done
			}
			fmt.Printf("found CreatedEvent contractId=%s templateId=%s party=%s hasCreateArgs=%t\n",
				createdEvent.ContractID, createdEvent.TemplateID, party, createdEvent.CreateArguments != nil,
			)

		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			if err != nil {
				return GlobalConfigView{}, fmt.Errorf("failed to receive active contract: %w", err)
			}
		case <-ctx.Done():
			return GlobalConfigView{}, ctx.Err()
		}
	}
done:

	rec := contractEvent.CreateArguments.(*apiv2.Record)

	ccipOwner := PartyVal(FieldVal(rec, "ccipOwner"))
	instanceId := TextVal(FieldVal(rec, "instanceId"))
	chainSelectorStr := NormalizeNumeric0Val(NumericVal(FieldVal(rec, "chainSelector")))

	onRampAddress := TextVal(FieldVal(rec, "onRampAddress"))

	dest := map[string]DestChainConfigView{}
	for _, e := range GenMapEntriesVal(FieldVal(rec, "destChainConfigs")) {
		key := NormalizeNumeric0Val(NumericVal(e.GetKey())) // "2222222222."
		r := RecordVal(e.GetValue())
		dest[key] = DestChainConfigView{
			IsEnabled:        BoolVal(FieldVal(r, "isEnabled")),
			DefaultExecutor:  TextVal(FieldVal(r, "defaultExecutor")),
			OffRampAddress:   TextVal(FieldVal(r, "offRampAddress")),
			LaneMandatedCCVs: TextListVal(FieldVal(r, "laneMandatedCCVs")),
			DefaultCCVs:      TextListVal(FieldVal(r, "defaultCCVs")),
		}
	}

	src := map[string]SourceChainConfigView{}
	for _, e := range GenMapEntriesVal(FieldVal(rec, "sourceChainConfigs")) {
		key := NormalizeNumeric0Val(NumericVal(e.GetKey())) // "3333333333."
		r := RecordVal(e.GetValue())
		src[key] = SourceChainConfigView{
			IsEnabled:        BoolVal(FieldVal(r, "isEnabled")),
			OnRampAddress:    TextVal(FieldVal(r, "onRampAddress")),
			LaneMandatedCCVs: TextListVal(FieldVal(r, "laneMandatedCCVs")),
			DefaultCCVs:      TextListVal(FieldVal(r, "defaultCCVs")),
		}
	}

	return GlobalConfigView{
		ContractMetaData: ContractMetaData{
			Address:    globalConfigContractID,
			Owner:      ccipOwner,
			TemplateID: templateID,
			InstanceID: instanceId,
		},
		CcipOwner:          ccipOwner,
		InstanceId:         instanceId,
		ChainSelector:      chainSelectorStr,
		OnRampAddress:      onRampAddress,
		DestChainConfigs:   dest,
		SourceChainConfigs: src,
	}, nil

}
