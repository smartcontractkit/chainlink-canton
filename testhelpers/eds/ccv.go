package eds

import (
	"context"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

type CCVExecuteDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      common.CCIPContext
}

func GetCCVExecuteDisclosure(
	ctx context.Context,
	ccvAPIClient oapiCCV.ClientWithResponsesInterface,
	encodedMessageHex string,
	ccvAddress contracts.InstanceAddress,
) (*CCVExecuteDisclosure, error) {
	resp, err := ccvAPIClient.PostCCVExecuteWithResponse(ctx, ccvAddress.String(), oapiCCV.CCVExecuteRequest{
		EncodedMessage: encodedMessageHex,
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

	choiceContext, err := CCIPContextFromData(resp.JSON200.ContextData)
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
	ChoiceContext      common.CCIPContext
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

	choiceContext, err := CCIPContextFromData(resp.JSON200.ContextData)
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
