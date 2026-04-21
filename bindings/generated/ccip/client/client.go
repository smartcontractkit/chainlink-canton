package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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
	PackageName = "ccip-client"
	PackageID   = "199a76896211a22c486133f5c3adb5325e4779e8249e913d71ae909c37e5ed34"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
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

// CCVExtraArg is a Record type
type CCVExtraArg struct {
	CcvAddress mcms.RawInstanceAddress `json:"ccvAddress"`
	CcvArgs    types.TEXT              `json:"ccvArgs"`
}

// ToMap converts CCVExtraArg to a map for DAML arguments
func (t CCVExtraArg) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvAddress"] = model.NestedToDAMLValue(t.CcvAddress)

	m["ccvArgs"] = string(t.CcvArgs)

	return m
}

func (t CCVExtraArg) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVExtraArg) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCVExtraArg to hex string (Canton MCMS format)
func (t CCVExtraArg) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVExtraArg from hex string (Canton MCMS format)
func (t *CCVExtraArg) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Canton2AnyMessage is a Record type
type Canton2AnyMessage struct {
	Receiver      types.TEXT                               `json:"receiver"`
	Payload       types.TEXT                               `json:"payload"`
	TokenTransfer *TokenTransfer                           `json:"tokenTransfer" hex:"optional"`
	FeeToken      splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	ExtraArgs     ExtraArgs                                `json:"extraArgs"`
}

// ToMap converts Canton2AnyMessage to a map for DAML arguments
func (t Canton2AnyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

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

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t Canton2AnyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Canton2AnyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Canton2AnyMessage to hex string (Canton MCMS format)
func (t Canton2AnyMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Canton2AnyMessage from hex string (Canton MCMS format)
func (t *Canton2AnyMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorExtraArg is a variant/union type
type ExecutorExtraArg struct {
	ExecutorNoExecutor  *types.UNIT          `json:"Executor_NoExecutor,omitempty"`
	ExecutorUseDefault  *ExecutorUseDefault  `json:"Executor_UseDefault,omitempty"`
	ExecutorWithAddress *ExecutorWithAddress `json:"Executor_WithAddress,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ExecutorExtraArg
func (v ExecutorExtraArg) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for ExecutorExtraArg
func (v *ExecutorExtraArg) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes ExecutorExtraArg to hex string (Canton MCMS format)
func (v ExecutorExtraArg) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes ExecutorExtraArg from hex string (Canton MCMS format)
func (v *ExecutorExtraArg) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ExecutorExtraArg) GetVariantTag() string {

	if v.ExecutorNoExecutor != nil {
		return "Executor_NoExecutor"
	}

	if v.ExecutorUseDefault != nil {
		return "Executor_UseDefault"
	}

	if v.ExecutorWithAddress != nil {
		return "Executor_WithAddress"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ExecutorExtraArg) GetVariantValue() any {

	if v.ExecutorNoExecutor != nil {
		return v.ExecutorNoExecutor
	}

	if v.ExecutorUseDefault != nil {
		return v.ExecutorUseDefault
	}

	if v.ExecutorWithAddress != nil {
		return v.ExecutorWithAddress
	}

	return nil
}

var _ types.VARIANT = (*ExecutorExtraArg)(nil)

// ExecutorUseDefault is a Record type
type ExecutorUseDefault struct {
	ExecutorArgs types.TEXT `json:"executorArgs"`
}

// ToMap converts ExecutorUseDefault to a map for DAML arguments
func (t ExecutorUseDefault) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorArgs"] = string(t.ExecutorArgs)

	return m
}

func (t ExecutorUseDefault) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorUseDefault) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorUseDefault to hex string (Canton MCMS format)
func (t ExecutorUseDefault) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorUseDefault from hex string (Canton MCMS format)
func (t *ExecutorUseDefault) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorWithAddress is a Record type
type ExecutorWithAddress struct {
	ExecutorAddress mcms.RawInstanceAddress `json:"executorAddress"`
	ExecutorArgs    types.TEXT              `json:"executorArgs"`
}

// ToMap converts ExecutorWithAddress to a map for DAML arguments
func (t ExecutorWithAddress) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorAddress"] = model.NestedToDAMLValue(t.ExecutorAddress)

	m["executorArgs"] = string(t.ExecutorArgs)

	return m
}

func (t ExecutorWithAddress) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorWithAddress) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorWithAddress to hex string (Canton MCMS format)
func (t ExecutorWithAddress) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorWithAddress from hex string (Canton MCMS format)
func (t *ExecutorWithAddress) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExtraArgs is a variant/union type
type ExtraArgs struct {
	V3 *GenericExtraArgsV3 `json:"V3,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ExtraArgs
func (v ExtraArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for ExtraArgs
func (v *ExtraArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes ExtraArgs to hex string (Canton MCMS format)
func (v ExtraArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes ExtraArgs from hex string (Canton MCMS format)
func (v *ExtraArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ExtraArgs) GetVariantTag() string {

	if v.V3 != nil {
		return "V3"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ExtraArgs) GetVariantValue() any {

	if v.V3 != nil {
		return v.V3
	}

	return nil
}

var _ types.VARIANT = (*ExtraArgs)(nil)

// GenericExtraArgsV3 is a Record type
type GenericExtraArgsV3 struct {
	GasLimit      types.INT64      `json:"gasLimit"`
	Ccvs          []CCVExtraArg    `json:"ccvs"`
	Executor      ExecutorExtraArg `json:"executor"`
	TokenReceiver types.TEXT       `json:"tokenReceiver"`
	TokenArgs     types.TEXT       `json:"tokenArgs"`
}

// ToMap converts GenericExtraArgsV3 to a map for DAML arguments
func (t GenericExtraArgsV3) ToMap() map[string]any {
	m := make(map[string]any)

	m["gasLimit"] = int64(t.GasLimit)

	m["ccvs"] = func() []any {
		res := make([]any, 0, len(t.Ccvs))
		for _, e := range t.Ccvs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["executor"] = model.NestedToDAMLValue(t.Executor)

	m["tokenReceiver"] = string(t.TokenReceiver)

	m["tokenArgs"] = string(t.TokenArgs)

	return m
}

func (t GenericExtraArgsV3) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GenericExtraArgsV3) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GenericExtraArgsV3 to hex string (Canton MCMS format)
func (t GenericExtraArgsV3) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GenericExtraArgsV3 from hex string (Canton MCMS format)
func (t *GenericExtraArgsV3) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransfer is a Record type
type TokenTransfer struct {
	Token  splice_api_token_holding_v1.InstrumentId `json:"token"`
	Amount types.NUMERIC                            `json:"amount"`
}

// ToMap converts TokenTransfer to a map for DAML arguments
func (t TokenTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["token"] = model.NestedToDAMLValue(t.Token)

	m["amount"] = t.Amount

	return m
}

func (t TokenTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransfer to hex string (Canton MCMS format)
func (t TokenTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransfer from hex string (Canton MCMS format)
func (t *TokenTransfer) UnmarshalHex(data string) error {
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
