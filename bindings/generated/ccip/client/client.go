package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
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
	PackageID   = "1ddcf1b906e309facec95c7ab2a2342358022fc47f365dc020c31124148e2e53"
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

// CCVSendInput is a Record type
type CCVSendInput struct {
	CcvCid          types.CONTRACT_ID  `json:"ccvCid"`
	VerifierArgs    types.TEXT         `json:"verifierArgs"`
	CcvExtraContext common.CCIPContext `json:"ccvExtraContext"`
}

// ToMap converts CCVSendInput to a map for DAML arguments
func (t CCVSendInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvCid).(mapper); ok {
			return m.toMap()
		}
		return t.CcvCid
	}()

	m["verifierArgs"] = string(t.VerifierArgs)

	m["ccvExtraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.CcvExtraContext
	}()

	return m
}

func (t CCVSendInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVSendInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCVSendInput to hex string (Canton MCMS format)
func (t CCVSendInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVSendInput from hex string (Canton MCMS format)
func (t *CCVSendInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Canton2AnyMessage is a Record type
type Canton2AnyMessage struct {
	Receiver      types.TEXT          `json:"receiver"`
	Payload       types.TEXT          `json:"payload"`
	TokenTransfer *TokenTransferInput `json:"tokenTransfer" hex:"optional"`
	FeeToken      FeeTokenInput       `json:"feeToken"`
	ExtraArgs     ExtraArgs           `json:"extraArgs"`
}

// ToMap converts Canton2AnyMessage to a map for DAML arguments
func (t Canton2AnyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenTransfer,
		}
	} else {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
		}
	}

	m["feeToken"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	m["extraArgs"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

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

// ExecutorInput is a Record type
type ExecutorInput struct {
	ExecutorCid          types.CONTRACT_ID  `json:"executorCid"`
	ExecutorArgs         types.TEXT         `json:"executorArgs"`
	ExecutorExtraContext common.CCIPContext `json:"executorExtraContext"`
}

// ToMap converts ExecutorInput to a map for DAML arguments
func (t ExecutorInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutorCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutorCid
	}()

	m["executorArgs"] = string(t.ExecutorArgs)

	m["executorExtraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutorExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutorExtraContext
	}()

	return m
}

func (t ExecutorInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorInput to hex string (Canton MCMS format)
func (t ExecutorInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorInput from hex string (Canton MCMS format)
func (t *ExecutorInput) UnmarshalHex(data string) error {
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

// FeeTokenInput is a Record type
type FeeTokenInput struct {
	Token           splice_api_token_holding_v1.InstrumentId `json:"token"`
	TokenInput      interfaces.TokenInput                    `json:"tokenInput"`
	SenderInputCids []types.CONTRACT_ID                      `json:"senderInputCids"`
}

// ToMap converts FeeTokenInput to a map for DAML arguments
func (t FeeTokenInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["token"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Token).(mapper); ok {
			return m.toMap()
		}
		return t.Token
	}()

	m["tokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t FeeTokenInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeTokenInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeTokenInput to hex string (Canton MCMS format)
func (t FeeTokenInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeTokenInput from hex string (Canton MCMS format)
func (t *FeeTokenInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GenericExtraArgsV3 is a Record type
type GenericExtraArgsV3 struct {
	GasLimit           types.INT64               `json:"gasLimit"`
	BlockConfirmations types.INT64               `json:"blockConfirmations"`
	Ccvs               []mcms.RawInstanceAddress `json:"ccvs"`
	Executor           *ExecutorInput            `json:"executor" hex:"optional"`
	TokenReceiver      types.TEXT                `json:"tokenReceiver"`
	TokenArgs          types.TEXT                `json:"tokenArgs"`
}

// ToMap converts GenericExtraArgsV3 to a map for DAML arguments
func (t GenericExtraArgsV3) ToMap() map[string]any {
	m := make(map[string]any)

	m["gasLimit"] = int64(t.GasLimit)

	m["blockConfirmations"] = int64(t.BlockConfirmations)

	m["ccvs"] = func() []any {
		res := make([]any, 0, len(t.Ccvs))
		for _, e := range t.Ccvs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.Executor != nil {
		m["executor"] = map[string]any{
			"_type": "optional",
			"value": *t.Executor,
		}
	} else {
		m["executor"] = map[string]any{
			"_type": "optional",
		}
	}

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

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	Token            splice_api_token_holding_v1.InstrumentId `json:"token"`
	Amount           types.NUMERIC                            `json:"amount"`
	SenderInputCids  []types.CONTRACT_ID                      `json:"senderInputCids"`
	TokenPoolCid     types.CONTRACT_ID                        `json:"tokenPoolCid"`
	TokenInput       interfaces.TokenInput                    `json:"tokenInput"`
	PoolExtraContext common.CCIPContext                       `json:"poolExtraContext"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["token"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Token).(mapper); ok {
			return m.toMap()
		}
		return t.Token
	}()

	m["amount"] = t.Amount

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["tokenPoolCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenPoolCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenPoolCid
	}()

	m["tokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["poolExtraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PoolExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolExtraContext
	}()

	return m
}

func (t TokenTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferInput to hex string (Canton MCMS format)
func (t TokenTransferInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferInput from hex string (Canton MCMS format)
func (t *TokenTransferInput) UnmarshalHex(data string) error {
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
