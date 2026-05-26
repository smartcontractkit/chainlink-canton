package eds

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
)

func DisclosedContractToProto(contract oapiCommon.DisclosedContract) (*apiv2.DisclosedContract, error) {
	id, err := contracts.TemplateIDFromString(contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       id.ToLedgerIdentifier(),
		ContractId:       contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   contract.SynchronizerId,
	}, nil
}

func GetGlobalDisclosureBatch(ctx context.Context, globalAPIClient oapiGlobal.ClientWithResponsesInterface, addresses []contracts.InstanceAddress) ([]*apiv2.DisclosedContract, error) {
	requestedAddresses := make([]oapiCommon.RawOrHashedAddress, len(addresses))
	for i, address := range addresses {
		_ = requestedAddresses[i].FromInstanceAddress(address.Hex())
	}
	resp, err := globalAPIClient.PostGetExplicitDisclosureBatchWithResponse(ctx, oapiGlobal.GetExplicitDisclosureBatchRequest{
		Addresses: requestedAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("error calling GetExplicitDisclosureBatch: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d; response: %s", resp.StatusCode(), string(resp.Body))
	}

	disclosedContracts := make([]*apiv2.DisclosedContract, len(resp.JSON200.Disclosures))
	for i, disclosure := range resp.JSON200.Disclosures {
		disclosedContracts[i], err = DisclosedContractToProto(disclosure)
		if err != nil {
			return nil, fmt.Errorf("failed to convert disclosed contract: %w", err)
		}
	}

	return disclosedContracts, nil
}
