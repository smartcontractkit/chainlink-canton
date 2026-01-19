package perpartyrouter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

// IAnyContract is a DAML interface
type IAnyContract interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// AnyContractView is a Record type
type AnyContractView struct {
}

// toMap converts AnyContractView to a map for DAML arguments
func (t AnyContractView) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for AnyContractView using JsonCodec
func (t AnyContractView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for AnyContractView using JsonCodec
func (t *AnyContractView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AnyValue is a variant/union type
type AnyValue struct {
	AVText       *TEXT        `json:"AV_Text,omitempty"`
	AVInt        *INT64       `json:"AV_Int,omitempty"`
	AVDecimal    *NUMERIC     `json:"AV_Decimal,omitempty"`
	AVBool       *BOOL        `json:"AV_Bool,omitempty"`
	AVDate       *DATE        `json:"AV_Date,omitempty"`
	AVTime       *TIMESTAMP   `json:"AV_Time,omitempty"`
	AVRelTime    *RELTIME     `json:"AV_RelTime,omitempty"`
	AVParty      *PARTY       `json:"AV_Party,omitempty"`
	AVContractId *CONTRACT_ID `json:"AV_ContractId,omitempty"`
	AVList       *[]AnyValue  `json:"AV_List,omitempty"`
	AVMap        *TEXTMAP     `json:"AV_Map,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AnyValue
func (v AnyValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for AnyValue
func (v *AnyValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v AnyValue) GetVariantTag() string {

	if v.AVText != nil {
		return "AV_Text"
	}

	if v.AVInt != nil {
		return "AV_Int"
	}

	if v.AVDecimal != nil {
		return "AV_Decimal"
	}

	if v.AVBool != nil {
		return "AV_Bool"
	}

	if v.AVDate != nil {
		return "AV_Date"
	}

	if v.AVTime != nil {
		return "AV_Time"
	}

	if v.AVRelTime != nil {
		return "AV_RelTime"
	}

	if v.AVParty != nil {
		return "AV_Party"
	}

	if v.AVContractId != nil {
		return "AV_ContractId"
	}

	if v.AVList != nil {
		return "AV_List"
	}

	if v.AVMap != nil {
		return "AV_Map"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v AnyValue) GetVariantValue() interface{} {

	if v.AVText != nil {
		return v.AVText
	}

	if v.AVInt != nil {
		return v.AVInt
	}

	if v.AVDecimal != nil {
		return v.AVDecimal
	}

	if v.AVBool != nil {
		return v.AVBool
	}

	if v.AVDate != nil {
		return v.AVDate
	}

	if v.AVTime != nil {
		return v.AVTime
	}

	if v.AVRelTime != nil {
		return v.AVRelTime
	}

	if v.AVParty != nil {
		return v.AVParty
	}

	if v.AVContractId != nil {
		return v.AVContractId
	}

	if v.AVList != nil {
		return v.AVList
	}

	if v.AVMap != nil {
		return v.AVMap
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*AnyValue)(nil)

// ChoiceContext is a Record type
type ChoiceContext struct {
	Values TEXTMAP `json:"values"`
}

// toMap converts ChoiceContext to a map for DAML arguments
func (t ChoiceContext) toMap() map[string]interface{} {
	return map[string]interface{}{

		"values": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Values).(mapper); ok {
				return m.toMap()
			}
			return t.Values
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ChoiceContext using JsonCodec
func (t ChoiceContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ChoiceContext using JsonCodec
func (t *ChoiceContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ChoiceExecutionMetadata is a Record type
type ChoiceExecutionMetadata struct {
	Meta Metadata `json:"meta"`
}

// toMap converts ChoiceExecutionMetadata to a map for DAML arguments
func (t ChoiceExecutionMetadata) toMap() map[string]interface{} {
	return map[string]interface{}{

		"meta": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Meta).(mapper); ok {
				return m.toMap()
			}
			return t.Meta
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ChoiceExecutionMetadata using JsonCodec
func (t ChoiceExecutionMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ChoiceExecutionMetadata using JsonCodec
func (t *ChoiceExecutionMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExtraArgs is a Record type
type ExtraArgs struct {
	Context ChoiceContext `json:"context"`
	Meta    Metadata      `json:"meta"`
}

// toMap converts ExtraArgs to a map for DAML arguments
func (t ExtraArgs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"context": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Context).(mapper); ok {
				return m.toMap()
			}
			return t.Context
		}(),
		"meta": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Meta).(mapper); ok {
				return m.toMap()
			}
			return t.Meta
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ExtraArgs using JsonCodec
func (t ExtraArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExtraArgs using JsonCodec
func (t *ExtraArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Metadata is a Record type
type Metadata struct {
	Values TEXTMAP `json:"values"`
}

// toMap converts Metadata to a map for DAML arguments
func (t Metadata) toMap() map[string]interface{} {
	return map[string]interface{}{

		"values": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Values).(mapper); ok {
				return m.toMap()
			}
			return t.Values
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for Metadata using JsonCodec
func (t Metadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Metadata using JsonCodec
func (t *Metadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IAnyContractInterfaceID returns the interface ID for the IAnyContract interface
func IAnyContractInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "Splice.Api.Token.MetadataV1", "AnyContract")
}
