package splice_api_token_metadata_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

const (
	PackageName = "splice-api-token-metadata-v1"
	PackageID   = "4ded6b668cb3b64f7a88a30874cd41c75829f5e064b3fbbadf41ec7e8363354f"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IAnyContract is a DAML interface
type IAnyContract interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

func argsToMap(args any) map[string]any {
	if args == nil {
		return map[string]any{}
	}

	if m, ok := args.(map[string]any); ok {
		return m
	}

	type mapper interface {
		ToMap() map[string]any
	}
	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]any{"args": args}
}

// AnyContractView is a Record type
type AnyContractView struct {
}

// ToMap converts AnyContractView to a map for DAML arguments
func (t AnyContractView) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t AnyContractView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AnyContractView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AnyContractView to hex string (Canton MCMS format)
func (t AnyContractView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AnyContractView from hex string (Canton MCMS format)
func (t *AnyContractView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AnyValue is a variant/union type
type AnyValue struct {
	AVText       *types.TEXT        `json:"AV_Text,omitempty"`
	AVInt        *types.INT64       `json:"AV_Int,omitempty"`
	AVDecimal    *types.NUMERIC     `json:"AV_Decimal,omitempty"`
	AVBool       *types.BOOL        `json:"AV_Bool,omitempty"`
	AVDate       *types.DATE        `json:"AV_Date,omitempty"`
	AVTime       *types.TIMESTAMP   `json:"AV_Time,omitempty"`
	AVRelTime    *types.RELTIME     `json:"AV_RelTime,omitempty"`
	AVParty      *types.PARTY       `json:"AV_Party,omitempty"`
	AVContractId *types.CONTRACT_ID `json:"AV_ContractId,omitempty"`
	AVList       *[]AnyValue        `json:"AV_List,omitempty"`
	AVMap        *types.TEXTMAP     `json:"AV_Map,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AnyValue
func (v AnyValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for AnyValue
func (v *AnyValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes AnyValue to hex string (Canton MCMS format)
func (v AnyValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes AnyValue from hex string (Canton MCMS format)
func (v *AnyValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
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
func (v AnyValue) GetVariantValue() any {

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

var _ types.VARIANT = (*AnyValue)(nil)

// ChoiceContext is a Record type
type ChoiceContext struct {
	Values types.TEXTMAP `json:"values"`
}

// ToMap converts ChoiceContext to a map for DAML arguments
func (t ChoiceContext) ToMap() map[string]any {
	m := make(map[string]any)

	m["values"] = model.NestedToDAMLValue(t.Values)

	return m
}

func (t ChoiceContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ChoiceContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ChoiceContext to hex string (Canton MCMS format)
func (t ChoiceContext) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ChoiceContext from hex string (Canton MCMS format)
func (t *ChoiceContext) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ChoiceExecutionMetadata is a Record type
type ChoiceExecutionMetadata struct {
	Meta Metadata `json:"meta"`
}

// ToMap converts ChoiceExecutionMetadata to a map for DAML arguments
func (t ChoiceExecutionMetadata) ToMap() map[string]any {
	m := make(map[string]any)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ChoiceExecutionMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ChoiceExecutionMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ChoiceExecutionMetadata to hex string (Canton MCMS format)
func (t ChoiceExecutionMetadata) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ChoiceExecutionMetadata from hex string (Canton MCMS format)
func (t *ChoiceExecutionMetadata) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExtraArgs is a Record type
type ExtraArgs struct {
	Context ChoiceContext `json:"context"`
	Meta    Metadata      `json:"meta"`
}

// ToMap converts ExtraArgs to a map for DAML arguments
func (t ExtraArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ExtraArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExtraArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExtraArgs to hex string (Canton MCMS format)
func (t ExtraArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExtraArgs from hex string (Canton MCMS format)
func (t *ExtraArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Metadata is a Record type
type Metadata struct {
	Values types.TEXTMAP `json:"values"`
}

// ToMap converts Metadata to a map for DAML arguments
func (t Metadata) ToMap() map[string]any {
	m := make(map[string]any)

	m["values"] = model.NestedToDAMLValue(t.Values)

	return m
}

func (t Metadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Metadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Metadata to hex string (Canton MCMS format)
func (t Metadata) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Metadata from hex string (Canton MCMS format)
func (t *Metadata) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IAnyContractInterfaceID returns the interface ID for the IAnyContract interface using the package name
func IAnyContractInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.MetadataV1", "AnyContract")
}

// IAnyContractInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IAnyContractInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.MetadataV1", "AnyContract")
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
}

// encoder provides typed encoding methods for choice parameters (unexported).
// It wraps bind.BoundTemplate to encode parameters to hex-encoded operation data.
type encoder struct {
	*bind.BoundTemplate
}

// Contract wraps template operations with Sui-style API access.
// Use NewContract to create instances, then call Encoder() for encoding methods.
type Contract struct {
	enc *encoder
}

// NewContract creates a Contract with encoder for the given template.
// This provides Sui-style API: contract.Encoder().Method(args)
func NewContract(packageID, moduleName, templateName string) *Contract {
	return &Contract{
		enc: &encoder{
			BoundTemplate: bind.NewBoundTemplate(packageID, moduleName, templateName),
		},
	}
}

// Encoder returns the encoder for Sui-style contract.Encoder().Method() usage.
func (c *Contract) Encoder() MCMSEncoder {
	return c.enc
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
