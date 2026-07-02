package ccipcodec

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
	PackageName = "ccip-codec"
	PackageID   = "0e57bb05bbfe395a24b351f66e66976ae5e1c96a5b840a0b2a4b05fe4bed01a1"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	WaitForFinalityFlag      = types.TEXT("00000000")
	MinBlockDepth            = types.INT64(1)
	MaxBlockDepth            = types.INT64(65535)
	FinalityConfigByteLength = types.INT64(4)
	MaxNumeric0IntegerText   = types.TEXT("99999999999999999999999999999999999999")
)

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

// DecodedFinality is a Record type
type DecodedFinality struct {
	Raw       types.TEXT     `json:"raw"`
	Requested FinalityConfig `json:"requested"`
}

// ToMap converts DecodedFinality to a map for DAML arguments
func (t DecodedFinality) ToMap() map[string]any {
	m := make(map[string]any)

	m["raw"] = string(t.Raw)

	m["requested"] = model.NestedToDAMLValue(t.Requested)

	return m
}

func (t DecodedFinality) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DecodedFinality) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DecodedFinality to hex string (Canton MCMS format)
func (t DecodedFinality) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DecodedFinality from hex string (Canton MCMS format)
func (t *DecodedFinality) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalityConfig is a variant/union type
type FinalityConfig struct {
	WaitForFinality *types.UNIT  `json:"WaitForFinality,omitempty"`
	WaitForSafe     *types.UNIT  `json:"WaitForSafe,omitempty"`
	BlockDepth      *types.INT64 `json:"BlockDepth,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for FinalityConfig
func (v FinalityConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for FinalityConfig
func (v *FinalityConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes FinalityConfig to hex string (Canton MCMS format)
func (v FinalityConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes FinalityConfig from hex string (Canton MCMS format)
func (v *FinalityConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v FinalityConfig) GetVariantTag() string {

	if v.WaitForFinality != nil {
		return "WaitForFinality"
	}

	if v.WaitForSafe != nil {
		return "WaitForSafe"
	}

	if v.BlockDepth != nil {
		return "BlockDepth"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v FinalityConfig) GetVariantValue() any {

	if v.WaitForFinality != nil {
		return v.WaitForFinality
	}

	if v.WaitForSafe != nil {
		return v.WaitForSafe
	}

	if v.BlockDepth != nil {
		return v.BlockDepth
	}

	return nil
}

var _ types.VARIANT = (*FinalityConfig)(nil)

// GetVariantTagByte implements types.VariantWithTagByte interface for MCMS numeric tag encoding
func (v FinalityConfig) GetVariantTagByte() byte {

	if v.WaitForFinality != nil {
		return 0
	}

	if v.WaitForSafe != nil {
		return 1
	}

	if v.BlockDepth != nil {
		return 2
	}

	return 0xFF // Invalid/unknown variant
}

var _ types.VariantWithTagByte = (*FinalityConfig)(nil)

// MessageV1 is a Record type
type MessageV1 struct {
	SourceChainSelector types.NUMERIC    `json:"sourceChainSelector"`
	DestChainSelector   types.NUMERIC    `json:"destChainSelector"`
	SequenceNumber      types.NUMERIC    `json:"sequenceNumber"`
	ExecutionGasLimit   types.INT64      `json:"executionGasLimit"`
	CcipReceiveGasLimit types.INT64      `json:"ccipReceiveGasLimit"`
	Finality            DecodedFinality  `json:"finality"`
	CcvAndExecutorHash  types.TEXT       `json:"ccvAndExecutorHash"`
	OnRampAddress       types.TEXT       `json:"onRampAddress"`
	OffRampAddress      types.TEXT       `json:"offRampAddress" hex:"bytes"`
	Sender              types.TEXT       `json:"sender"`
	Receiver            types.TEXT       `json:"receiver"`
	DestBlob            types.TEXT       `json:"destBlob"`
	TokenTransfer       *TokenTransferV1 `json:"tokenTransfer" hex:"optional"`
	MessageData         types.TEXT       `json:"messageData"`
}

// ToMap converts MessageV1 to a map for DAML arguments
func (t MessageV1) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["destChainSelector"] = t.DestChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["executionGasLimit"] = int64(t.ExecutionGasLimit)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	m["onRampAddress"] = string(t.OnRampAddress)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["sender"] = string(t.Sender)

	m["receiver"] = string(t.Receiver)

	m["destBlob"] = string(t.DestBlob)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransfer),
		}
	} else {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["messageData"] = string(t.MessageData)

	return m
}

func (t MessageV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MessageV1 to hex string (Canton MCMS format)
func (t MessageV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MessageV1 from hex string (Canton MCMS format)
func (t *MessageV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferV1 is a Record type
type TokenTransferV1 struct {
	Amount             types.TEXT `json:"amount"`
	SourcePoolAddress  types.TEXT `json:"sourcePoolAddress"`
	SourceTokenAddress types.TEXT `json:"sourceTokenAddress"`
	DestTokenAddress   types.TEXT `json:"destTokenAddress"`
	TokenReceiver      types.TEXT `json:"tokenReceiver"`
	ExtraData          types.TEXT `json:"extraData"`
}

// ToMap converts TokenTransferV1 to a map for DAML arguments
func (t TokenTransferV1) ToMap() map[string]any {
	m := make(map[string]any)

	m["amount"] = string(t.Amount)

	m["sourcePoolAddress"] = string(t.SourcePoolAddress)

	m["sourceTokenAddress"] = string(t.SourceTokenAddress)

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["tokenReceiver"] = string(t.TokenReceiver)

	m["extraData"] = string(t.ExtraData)

	return m
}

func (t TokenTransferV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferV1 to hex string (Canton MCMS format)
func (t TokenTransferV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferV1 from hex string (Canton MCMS format)
func (t *TokenTransferV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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
