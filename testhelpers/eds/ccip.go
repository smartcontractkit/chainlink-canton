package eds

import (
	"context"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type PerPartyRouterFactoryDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
}

func GetPerPartyRouterFactoryDisclosure(
	ctx context.Context,
	ccipAPIClient oapiCCIP.ClientWithResponsesInterface,
	partyId string,
) (*PerPartyRouterFactoryDisclosure, error) {
	resp, err := ccipAPIClient.PostPerPartyRouterFactoryWithResponse(ctx, oapiCCIP.CCIPPerPartyRouterFactoryRequest{
		PartyID: partyId, // Unused for now
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get per party router factory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := DisclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	address, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PerPartyRouterFactory address: %w", err)
	}

	return &PerPartyRouterFactoryDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
	}, nil
}

func GetTokenPoolForToken(ctx context.Context, ccipAPIClient oapiCCIP.ClientWithResponsesInterface, token contracts.EncodedInstrumentID) (contracts.RawInstanceAddress, error) {
	resp, err := ccipAPIClient.GetTokenAdminRegistryTokenWithResponse(ctx, token.String())
	if err != nil {
		return contracts.RawInstanceAddress(""), fmt.Errorf("failed to get token pool for token: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return contracts.RawInstanceAddress(""), fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	tokenPoolAddress, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return contracts.RawInstanceAddress(""), fmt.Errorf("failed to parse token pool address: %w", err)
	}

	return tokenPoolAddress, nil
}

type CCIPExecuteDisclosure struct {
	ChoiceContext      common.CCIPContext
	DisclosedContracts []*apiv2.DisclosedContract
	TokenPool          *contracts.RawInstanceAddress
}

func GetCCIPExecuteDisclosure(
	ctx context.Context,
	ccipAPIClient oapiCCIP.ClientWithResponsesInterface,
	encodedMessageHex string,
) (*CCIPExecuteDisclosure, error) {
	resp, err := ccipAPIClient.PostCCIPExecuteWithResponse(ctx, oapiCCIP.CCIPExecuteRequest{
		EncodedMessage: encodedMessageHex,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling CCIPExecute: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := DisclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	choiceContext, err := contracts.CCIPContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	var tokenPool *contracts.RawInstanceAddress
	if resp.JSON200.TokenPool != nil {
		tp, err := contracts.RawInstanceAddressFromString(*resp.JSON200.TokenPool)
		if err != nil {
			return nil, fmt.Errorf("received invalid token pool raw instance address: %w", err)
		}
		tokenPool = &tp
	}

	return &CCIPExecuteDisclosure{
		ChoiceContext:      choiceContext,
		DisclosedContracts: disclosedContracts,
		TokenPool:          tokenPool,
	}, nil
}

type CCIPSendDisclosure struct {
	ChoiceContext      common.CCIPContext
	DisclosedContracts []*apiv2.DisclosedContract
	CCVs               []string
	Executor           *string
}

func GetCCIPSendDisclosure(
	ctx context.Context,
	ccipAPIClient oapiCCIP.ClientWithResponsesInterface,
	message oapiCommon.Message,
	senderRequiredCCVs []string,
	tokenPoolRequiredCCVs []string,
) (*CCIPSendDisclosure, error) {
	senderCCVs := make([]oapiCommon.RawOrHashedAddress, len(senderRequiredCCVs))
	for i, v := range senderRequiredCCVs {
		_ = senderCCVs[i].FromRawInstanceAddress(v)
	}
	tokenPoolCCVs := make([]oapiCommon.RawOrHashedAddress, len(tokenPoolRequiredCCVs))
	for i, v := range tokenPoolRequiredCCVs {
		_ = tokenPoolCCVs[i].FromRawInstanceAddress(v)
	}

	resp, err := ccipAPIClient.PostCCIPSendWithResponse(ctx, oapiCCIP.CCIPSendRequest{
		Message:               message,
		SenderRequiredCCVs:    &senderCCVs,
		TokenPoolRequiredCCVs: &tokenPoolCCVs,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling CCIPSend: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := DisclosedContractToProto(contract)
		if err != nil {
			return nil, err
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	choiceContext, err := contracts.CCIPContextFromData(resp.JSON200.ContextData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert choice context: %w", err)
	}

	ccvs := make([]string, len(resp.JSON200.Ccvs))
	for i, ccv := range resp.JSON200.Ccvs {
		ccvs[i], _ = ccv.AsRawInstanceAddress()
	}

	var executor *string
	if resp.JSON200.Executor != nil {
		executorAddr, _ := resp.JSON200.Executor.AsRawInstanceAddress()
		executor = &executorAddr
	}

	return &CCIPSendDisclosure{
		ChoiceContext:      choiceContext,
		DisclosedContracts: disclosedContracts,
		CCVs:               ccvs,
		Executor:           executor,
	}, nil
}
