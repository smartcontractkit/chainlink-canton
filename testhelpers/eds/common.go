package eds

import (
	"encoding/base64"
	"fmt"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
)

func CCIPContextFromData(contextData map[string]any) (common.CCIPContext, error) {
	valuesIn, ok := contextData["values"].(map[string]any)
	if !ok {
		return common.CCIPContext{}, fmt.Errorf("no values found in context")
	}

	values := types.TEXTMAP{}
	for k, v := range valuesIn {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value common.AnyValue
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			vText := types.TEXT(valueString)
			value.AVText = &vText
		case "AV_Int":
			// JSON numbers come as float64
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			// TODO
			vInt := types.INT64(int64(valueFloat))
			value.AVInt = &vInt
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			vNumeric := types.NUMERIC(valueString)
			value.AVDecimal = &vNumeric
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			vBool := types.BOOL(valueBool)
			value.AVBool = &vBool
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			vDate := types.DATE(t)
			value.AVDate = &vDate
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			vTime := types.TIMESTAMP(t)
			value.AVTime = &vTime
		case "AV_RelTime":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_RelTime value is not a number: %T", rawValue)
			}
			vRelTime := types.RELTIME(time.Duration(int64(valueFloat)) * time.Microsecond)
			value.AVRelTime = &vRelTime
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			vContractId := types.CONTRACT_ID(valueString)
			value.AVContractId = &vContractId
		default:
			// Add lists and maps
			return common.CCIPContext{}, fmt.Errorf("unimplemented tag: %v", tag)
		}

		values[k] = value
	}

	return common.CCIPContext{
		Values: values,
	}, nil
}

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
