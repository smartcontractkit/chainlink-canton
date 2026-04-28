package contracts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
)

func CCIPContextFromData(contextData map[string]any) (common.CCIPContext, error) {
	valuesIn, ok := contextData["values"].(map[string]any)
	if !ok {
		return common.CCIPContext{}, fmt.Errorf("no values found in context")
	}

	values := map[string]common.AnyValue{}
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
			value.AVText = new(types.TEXT(valueString))
		case "AV_Int":
			// Int64s are encoded as JSON numbers or strings, depending on the encoder settings
			switch val := rawValue.(type) {
			case string:
				valueInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value.AVInt = new(types.INT64(valueInt))
			case json.Number:
				valueInt, err := strconv.ParseInt(val.String(), 10, 64)
				if err != nil {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value.AVInt = new(types.INT64(valueInt))
			case float64:
				// Some encoders may encode int64s as JSON numbers, which are float64s in Go. This can cause precision loss for large int64s, but we can still parse them if they fit within uint64.
				if val < 0 || val > float64(^uint64(0)) {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is out of range for uint64: %f", val)
				}
				value.AVInt = new(types.INT64(int64(val)))
			default:
				return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a string or number: %T", rawValue)
			}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value.AVDecimal = new(types.NUMERIC(valueString))
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value.AVBool = new(types.BOOL(valueBool))
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.DateOnly, valueString)
			if err != nil {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a DateOnly time: %s", valueString)
			}
			value.AVDate = new(types.DATE(t))
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return common.CCIPContext{}, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value.AVTime = new(types.TIMESTAMP(t))
		case "AV_RelTime":
			// Int64s are encoded as JSON numbers or strings, depending on the encoder settings
			switch val := rawValue.(type) {
			case string:
				valueInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value.AVRelTime = new(types.RELTIME(time.Duration(valueInt) * time.Microsecond))
			case json.Number:
				valueInt, err := strconv.ParseInt(val.String(), 10, 64)
				if err != nil {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a valid uint64 string: %s", val)
				}
				value.AVRelTime = new(types.RELTIME(time.Duration(valueInt) * time.Microsecond))
			case float64:
				// Some encoders may encode int64s as JSON numbers, which are float64s in Go. This can cause precision loss for large int64s, but we can still parse them if they fit within uint64.
				if val < 0 || val > float64(^uint64(0)) {
					return common.CCIPContext{}, fmt.Errorf("AV_Int value is out of range for uint64: %f", val)
				}
				value.AVRelTime = new(types.RELTIME(time.Duration(int64(val)) * time.Microsecond))
			default:
				return common.CCIPContext{}, fmt.Errorf("AV_Int value is not a string or number: %T", rawValue)
			}
		case "AV_Party":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_Party value is not a string: %T", rawValue)
			}
			value.AVParty = new(types.PARTY(valueString))
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return common.CCIPContext{}, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value.AVContractId = new(types.CONTRACT_ID(valueString))
		default:
			// TODO Add lists and maps
			return common.CCIPContext{}, fmt.Errorf("unimplemented tag: %v", tag)
		}

		values[k] = value
	}

	return common.CCIPContext{
		Values: values,
	}, nil
}

// Token metadata bindings generate a separate ChoiceContext/AnyValue type from
// the CCIP bindings, so we still need an explicit adapter until codegen emits a
// shared canonical type.
func CCIPContextToChoiceContext(ctx common.CCIPContext) (splice_api_token_metadata_v1.ChoiceContext, error) {
	var result splice_api_token_metadata_v1.ChoiceContext

	payload, err := json.Marshal(ctx)
	if err != nil {
		return splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("marshal CCIPContext: %w", err)
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("unmarshal ChoiceContext: %w", err)
	}

	return result, nil
}
