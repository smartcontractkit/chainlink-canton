package eds

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
)

type TokenPoolExecuteDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      common.CCIPContext
	RequiredCCVs       []string
	TokenInput         interfaces.TokenInput
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

	choiceContext, err := CCIPContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	requiredCCVs := make([]string, len(resp.JSON200.RequiredCCVs))
	for i, ccv := range resp.JSON200.RequiredCCVs {
		requiredCCVs[i], _ = ccv.AsRawInstanceAddress()
	}

	tokenInputContext, err := CCIPContextFromData(resp.JSON200.TokenInput.ExtraArgs.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to convert token input choice context: %w", err)
	}

	tokenPoolHoldings := make([]types.CONTRACT_ID, len(resp.JSON200.TokenInput.TokenPoolHoldings))
	for i, id := range resp.JSON200.TokenInput.TokenPoolHoldings {
		tokenPoolHoldings[i] = types.CONTRACT_ID(id)
	}

	tokenInput := interfaces.TokenInput{
		TransferFactory: types.CONTRACT_ID(resp.JSON200.TokenInput.TransferFactory),
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: CCIPContextToChoiceContext(tokenInputContext),
			Meta: splice_api_token_metadata_v1.Metadata{},
		},
		TokenPoolHoldings: tokenPoolHoldings,
	}

	return &TokenPoolExecuteDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       requiredCCVs,
		TokenInput:         tokenInput,
	}, nil
}

type TokenPoolSendDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      common.CCIPContext
	RequiredCCVs       []string
	TokenInput         interfaces.TokenInput
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

	choiceContext, err := CCIPContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	requiredCCVs := make([]string, len(resp.JSON200.RequiredCCVs))
	for i, ccv := range resp.JSON200.RequiredCCVs {
		requiredCCVs[i], _ = ccv.AsRawInstanceAddress()
	}

	tokenInputContext, err := CCIPContextFromData(resp.JSON200.TokenInput.ExtraArgs.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to convert token input choice context: %w", err)
	}

	tokenInput := interfaces.TokenInput{
		TransferFactory: types.CONTRACT_ID(resp.JSON200.TokenInput.TransferFactory),
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: CCIPContextToChoiceContext(tokenInputContext),
			Meta: splice_api_token_metadata_v1.Metadata{},
		},
	}

	return &TokenPoolSendDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
		RequiredCCVs:       requiredCCVs,
		TokenInput:         tokenInput,
	}, nil
}
