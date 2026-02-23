package testhelpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

func GetPerPartyRouterFactoryDisclosures(ctx context.Context, edsClient *edsv1.ClientWithResponses) (string, []*apiv2.DisclosedContract, error) {
	resp, err := edsClient.PerPartyRouterFactoryWithResponse(ctx, edsv1.CCIPPerPartyRouterFactoryRequest{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to get per party router factory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decode created event blob: %w", err)
		}
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}

	factoryCid := resp.JSON200.PerPartyRouterFactoryId

	return factoryCid, disclosedContracts, nil
}

func GetCCIPExecuteDisclosures(ctx context.Context, edsClient *edsv1.ClientWithResponses) ([]*apiv2.DisclosedContract, *apiv2.Value, error) {
	resp, err := edsClient.CcipExecuteWithResponse(ctx, edsv1.CCIPExecuteRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("error calling CCIPExecute: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		id, err := TemplateIdFromString(contract.TemplateId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse template id: %w", err)
		}
		createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode created event blob: %w", err)
		}
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}
	choiceContext, err := ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	return disclosedContracts, choiceContext, nil
}
