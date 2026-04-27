package eds

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
)

type ExecutorDisclosure struct {
	ContractId         string
	Address            contracts.RawInstanceAddress
	DisclosedContracts []*apiv2.DisclosedContract
	ChoiceContext      common.CCIPContext
}

func GetExecutorSendDisclosure(
	ctx context.Context,
	executorAPIClient oapiExecutor.ClientWithResponsesInterface,
	message oapiCommon.Message,
	executorAddress contracts.InstanceAddress,
	ccvAddresses []string,
) (*ExecutorDisclosure, error) {
	ccvs := make([]oapiCommon.RawOrHashedAddress, len(ccvAddresses))
	for i, address := range ccvAddresses {
		_ = ccvs[i].FromInstanceAddress(address)
	}
	resp, err := executorAPIClient.PostExecutorSendWithResponse(ctx, executorAddress.Hex(), oapiExecutor.ExecutorSendRequest{
		Ccvs:    ccvs,
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling Executor Send: %w", err)
	}
	if resp.StatusCode() != 200 {
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

	return &ExecutorDisclosure{
		ContractId:         resp.JSON200.ContractId,
		Address:            address,
		DisclosedContracts: disclosedContracts,
		ChoiceContext:      choiceContext,
	}, nil
}
