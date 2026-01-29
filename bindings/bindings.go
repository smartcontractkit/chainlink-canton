package bindings

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/noders-team/go-daml/pkg/codec"
	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

// UnmarshalActiveContract unmarshals a Canton ActiveContract response into a typed DAML binding struct.
// This provides type-safe access to contract fields instead of manual string-based field parsing.
//
// Example usage:
//
//	mcmsContract, err := UnmarshalActiveContract[mcms.MCMS](activeContract)
//	if err != nil {
//	    return err
//	}
//	// Now use type-safe fields: mcmsContract.McmsId, mcmsContract.Config.Signers, etc.
func UnmarshalActiveContract[T any](ac *apiv2.GetActiveContractsResponse_ActiveContract) (*T, error) {
	if ac == nil {
		return nil, errors.New("active contract is nil")
	}

	createdEvent := ac.ActiveContract.GetCreatedEvent()
	if createdEvent == nil {
		return nil, errors.New("no created event in active contract")
	}

	return UnmarshalCreatedEvent[T](createdEvent)
}

// UnmarshalCreatedEvent unmarshals a CreatedEvent into a typed DAML binding struct.
// This is useful when you have a CreatedEvent directly (e.g., from transaction events).
func UnmarshalCreatedEvent[T any](event *apiv2.CreatedEvent) (*T, error) {
	if event == nil {
		return nil, errors.New("created event is nil")
	}

	createArgs := event.GetCreateArguments()
	if createArgs == nil {
		return nil, errors.New("no create arguments in created event")
	}

	// Convert protobuf Record to map (same as go-daml's valueFromProto)
	recordMap := valueFromRecord(createArgs)

	// Marshal to JSON
	jsonData, err := json.Marshal(recordMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	// Use go-daml codec to unmarshal into typed struct
	var result T
	jsonCodec := codec.NewJsonCodec()
	err = jsonCodec.Unmarshall(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to type %T: %w", result, err)
	}

	return &result, nil
}

// valueFromRecord converts protobuf Record to map (replicates go-daml's logic)
func valueFromRecord(record *apiv2.Record) map[string]interface{} {
	if record == nil {
		return nil
	}

	result := make(map[string]interface{})
	for _, field := range record.Fields {
		result[field.Label] = valueFromProto(field.Value)
	}
	return result
}

// valueFromProto converts protobuf Value to interface{} (replicates go-daml's logic)
func valueFromProto(pb *apiv2.Value) interface{} {
	if pb == nil {
		return nil
	}

	switch v := pb.Sum.(type) {
	case *apiv2.Value_Unit:
		return map[string]interface{}{"_type": "unit"}
	case *apiv2.Value_Bool:
		return v.Bool
	case *apiv2.Value_Int64:
		return v.Int64
	case *apiv2.Value_Text:
		return v.Text
	case *apiv2.Value_Numeric:
		return v.Numeric
	case *apiv2.Value_Party:
		return v.Party
	case *apiv2.Value_ContractId:
		return v.ContractId
	case *apiv2.Value_Date:
		return v.Date
	case *apiv2.Value_Timestamp:
		return v.Timestamp
	case *apiv2.Value_Optional:
		if v.Optional.Value != nil {
			return valueFromProto(v.Optional.Value)
		}
		return nil
	case *apiv2.Value_List:
		result := make([]interface{}, len(v.List.Elements))
		for i, elem := range v.List.Elements {
			result[i] = valueFromProto(elem)
		}
		return result
	case *apiv2.Value_Record:
		return valueFromRecord(v.Record)
	case *apiv2.Value_TextMap:
		result := make(map[string]interface{})
		for _, entry := range v.TextMap.Entries {
			result[entry.Key] = valueFromProto(entry.Value)
		}
		return result
	case *apiv2.Value_GenMap:
		// GenMap as map for DAML codec
		result := make(map[string]interface{})
		for _, entry := range v.GenMap.Entries {
			// Convert key to string for map key
			keyStr := fmt.Sprintf("%v", valueFromProto(entry.Key))
			result[keyStr] = valueFromProto(entry.Value)
		}
		return result
	case *apiv2.Value_Enum:
		if v.Enum != nil {
			return v.Enum.Constructor
		}
		return nil
	case *apiv2.Value_Variant:
		if v.Variant != nil {
			return map[string]interface{}{
				"tag":   v.Variant.Constructor,
				"value": valueFromProto(v.Variant.Value),
			}
		}
		return nil
	default:
		return nil
	}
}
