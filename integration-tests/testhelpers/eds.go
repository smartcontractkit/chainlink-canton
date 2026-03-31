package testhelpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts"
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

type ExecuteDisclosures struct {
	DisclosedContracts           []*apiv2.DisclosedContract
	ChoiceContext                *apiv2.Value
	PoolExtraContext             *apiv2.Value
	CCVContractIDs               []*apiv2.Value_ContractId
	TokenPoolContractID          *apiv2.Value_ContractId
	TokenPoolHoldingsContractIDs []*apiv2.Value
}

func GetCCIPExecuteDisclosures(
	ctx context.Context,
	encodedMessageHex string,
	edsClient *edsv1.ClientWithResponses,
	ccvs []contracts.InstanceAddress,
) (
	*ExecuteDisclosures,
	error,
) {
	// Add the requested CCVs to the API request
	request := edsv1.CCIPExecuteRequest{
		Ccvs:           make([]string, len(ccvs)),
		EncodedMessage: encodedMessageHex,
	}
	for i, ccv := range ccvs {
		request.Ccvs[i] = ccv.String()
	}

	resp, err := edsClient.CcipExecuteWithResponse(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error calling CCIPExecute: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// The top-level choice context has the following disclosures:
	// - OffRamp
	// - GlobalConfig
	// - TokenAdminRegistry
	// - RMNRemote
	// And, if a token transfer is present:
	// - TokenPool
	// - TokenPoolHolding
	// - InboundRateLimiter
	// - InboundCustomBlockConfirmationsRateLimiter
	// - OutboundRateLimiter
	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		disclosedContract, err := disclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}

		disclosedContracts = append(disclosedContracts, disclosedContract)
	}
	choiceContext, err := ChoiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	// Optionally populate poolExtraContext.
	// TODO: is there a cleaner approach to this?
	var poolExtraContext *apiv2.Value
	if len(resp.JSON200.PoolExtraContext) > 0 {
		poolExtraContext, err = ChoiceContextFromData(resp.JSON200.PoolExtraContext)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pool extra context: %w", err)
		}
	}

	// Optionally get token pool and holdings disclosures.
	var tokenPoolContractID *apiv2.Value_ContractId
	var tokenPoolHoldingsContractIDs []*apiv2.Value
	if resp.JSON200.TokenPool != nil {
		tokenPoolContractID = &apiv2.Value_ContractId{
			ContractId: resp.JSON200.TokenPool.ContractId,
		}
	}
	if resp.JSON200.TokenPoolHoldings != nil {
		tokenPoolHoldingsContractIDs = []*apiv2.Value{
			{Sum: &apiv2.Value_ContractId{ContractId: resp.JSON200.TokenPoolHoldings.ContractId}},
		}
	}

	// Handle CCVs
	ccvContractIDs := make([]*apiv2.Value_ContractId, len(ccvs))
	for i, ccv := range ccvs {
		// Check if the API returned an explicit disclosure for the requested CCV
		disclosure, ok := resp.JSON200.Ccvs[ccv.String()]
		if !ok || disclosure.DisclosedContract == nil {
			return nil, fmt.Errorf("failed to get disclosure for ccv: %s", ccv.String())
		}
		ccvContractIDs[i] = &apiv2.Value_ContractId{
			ContractId: disclosure.DisclosedContract.ContractId,
		}
		// Add the CCV's explicit disclosure to disclosedContracts
		disclosedContract, err := disclosedContractToProto(*disclosure.DisclosedContract)
		if err != nil {
			return nil, fmt.Errorf("failed to convert disclosed contract for ccv %s: %w", ccv.String(), err)
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	return &ExecuteDisclosures{
		DisclosedContracts:           disclosedContracts,
		ChoiceContext:                choiceContext,
		PoolExtraContext:             poolExtraContext,
		CCVContractIDs:               ccvContractIDs,
		TokenPoolContractID:          tokenPoolContractID,
		TokenPoolHoldingsContractIDs: tokenPoolHoldingsContractIDs,
	}, nil
}

func disclosedContractToProto(contract edsv1.DisclosedContract) (*apiv2.DisclosedContract, error) {
	id, err := TemplateIdFromString(contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       id,
		ContractId:       contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   contract.SynchronizerId,
	}, nil
}
