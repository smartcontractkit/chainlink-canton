package eds

import (
	"context"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type CCVExecuteDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      splice_api_token_metadata_v1.ChoiceContext
}

func GetCCVExecuteDisclosure(
	ctx context.Context,
	ccvAPIClient oapiCCV.ClientWithResponsesInterface,
	encodedMessageHex string,
	ccvAddress contracts.InstanceAddress,
	receiver types.PARTY,
) (*CCVExecuteDisclosure, error) {
	resp, err := ccvAPIClient.PostCCVExecuteWithResponse(ctx, ccvAddress.String(), oapiCCV.CCVExecuteRequest{
		EncodedMessage: encodedMessageHex,
		Receiver:       string(receiver),
	})
	if err != nil {
		return nil, fmt.Errorf("error calling CCVExecute: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	address, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CCV address: %w", err)
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

	return &CCVExecuteDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
	}, nil
}

type CCVSendDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      splice_api_token_metadata_v1.ChoiceContext
}

func GetCCVSendDisclosure(
	ctx context.Context,
	ccvAPIClient oapiCCV.ClientWithResponsesInterface,
	message oapiCommon.Message,
	ccvAddress contracts.InstanceAddress,
) (*CCVSendDisclosure, error) {
	resp, err := ccvAPIClient.PostCCVSendWithResponse(ctx, ccvAddress.String(), oapiCCV.CCVSendRequest{
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling CCVExecute: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	address, err := contracts.RawInstanceAddressFromString(resp.JSON200.RawInstanceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CCV address: %w", err)
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

	return &CCVSendDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
	}, nil
}
