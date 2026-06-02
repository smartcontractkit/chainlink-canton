package eds

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

type TokenPoolExecuteDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      splice_api_token_metadata_v1.ChoiceContext
	RequiredCCVs       []string
}

func GetTokenPoolExecuteDisclosure(
	ctx context.Context,
	tokenPoolAPIClient oapiTokenPool.ClientWithResponsesInterface,
	encodedMessageHex string,
	tokenPoolAddress contracts.InstanceAddress,
) (*TokenPoolExecuteDisclosure, error) {
	resp, err := tokenPoolAPIClient.PostTokenPoolExecuteWithResponse(ctx, tokenPoolAddress.String(), oapiTokenPool.TokenPoolExecuteRequest{
		EncodedMessage: encodedMessageHex,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling Token Pool Execute: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	address, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Token Pool address: %w", err)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := DisclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	requiredCCVs := make([]string, len(resp.JSON200.RequiredCCVs))
	for i, ccv := range resp.JSON200.RequiredCCVs {
		requiredCCVs[i], _ = ccv.AsRawInstanceAddress()
	}

	return &TokenPoolExecuteDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       requiredCCVs,
	}, nil
}

type TokenPoolSendDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      splice_api_token_metadata_v1.ChoiceContext
	RequiredCCVs       []string
}

func GetTokenPoolSendDisclosure(
	ctx context.Context,
	tokenPoolAPIClient oapiTokenPool.ClientWithResponsesInterface,
	message oapiCommon.Message,
	tokenPoolAddress contracts.InstanceAddress,
) (*TokenPoolSendDisclosure, error) {
	resp, err := tokenPoolAPIClient.PostTokenPoolSendWithResponse(ctx, tokenPoolAddress.String(), oapiTokenPool.TokenPoolSendRequest{
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling Token Pool Send: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	address, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Token Pool address: %w", err)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := DisclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	choiceContext, err := contracts.ChoiceContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	requiredCCVs := make([]string, len(resp.JSON200.RequiredCCVs))
	for i, ccv := range resp.JSON200.RequiredCCVs {
		requiredCCVs[i], _ = ccv.AsRawInstanceAddress()
	}

	return &TokenPoolSendDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       requiredCCVs,
	}, nil
}
