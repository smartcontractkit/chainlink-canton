package converters

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

// RawOrHashedAddress convenience helpers

func ResolveRawOrHashedAddress(address oapiCommon.RawOrHashedAddress) (contracts.InstanceAddress, error) {
	// Try to parse as RawInstanceAddress
	if rawInstanceAddressString, err := address.AsRawInstanceAddress(); err == nil {
		rawInstanceAddress, err := contracts.RawInstanceAddressFromString(rawInstanceAddressString)
		if err == nil {
			return rawInstanceAddress.InstanceAddress(), nil
		}
	}

	// Try to parse as InstanceAddress
	if instanceAddressString, err := address.AsInstanceAddress(); err == nil {
		instanceAddress := contracts.HexToInstanceAddress(instanceAddressString)
		if (instanceAddress != contracts.InstanceAddress{}) {
			return instanceAddress, nil
		}
	}

	return contracts.InstanceAddress{}, fmt.Errorf("invalid RawOrHashedAddress")
}

func ResolveAddress(address string) (contracts.InstanceAddress, error) {
	// Try to parse as RawInstanceAddress
	rawInstanceAddress, err := contracts.RawInstanceAddressFromString(address)
	if err == nil {
		return rawInstanceAddress.InstanceAddress(), nil
	}

	// Try to parse as InstanceAddress
	instanceAddress := contracts.HexToInstanceAddress(address)
	if (instanceAddress != contracts.InstanceAddress{}) {
		return instanceAddress, nil
	}

	return contracts.InstanceAddress{}, fmt.Errorf("invalid RawOrHashedAddress")
}

func RawOrHashedAddressAsString(address oapiCommon.RawOrHashedAddress) string {
	rawInstanceAddressString, _ := address.AsRawInstanceAddress()
	return rawInstanceAddressString
}

func InstanceAddressAsRawOrHashedAddress(instanceAddress contracts.InstanceAddress) oapiCommon.RawOrHashedAddress {
	var address oapiCommon.RawOrHashedAddress
	_ = address.FromInstanceAddress(instanceAddress.String())

	return address
}

func RawInstanceAddressAsRawOrHashedAddress(rawInstanceAddress contracts.RawInstanceAddress) oapiCommon.RawOrHashedAddress {
	var address oapiCommon.RawOrHashedAddress
	_ = address.FromRawInstanceAddress(rawInstanceAddress.String())

	return address
}

func ActiveContractToDisclosedContract(activeContract *apiv2.ActiveContract) oapiCommon.DisclosedContract {
	return oapiCommon.DisclosedContract{
		TemplateId:       fmt.Sprintf("%s:%s:%s", activeContract.GetCreatedEvent().GetTemplateId().GetPackageId(), activeContract.GetCreatedEvent().GetTemplateId().GetModuleName(), activeContract.GetCreatedEvent().GetTemplateId().GetEntityName()),
		ContractId:       activeContract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: base64.StdEncoding.EncodeToString(activeContract.GetCreatedEvent().GetCreatedEventBlob()),
		SynchronizerId:   activeContract.GetSynchronizerId(),
	}
}

func SerializeChoiceContext(context splice_api_token_metadata_v1.ChoiceContext) (map[string]any, error) {
	jsonBytes, err := context.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ChoiceContext: %w", err)
	}
	var data map[string]any
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ChoiceContext JSON: %w", err)
	}

	return data, nil
}
