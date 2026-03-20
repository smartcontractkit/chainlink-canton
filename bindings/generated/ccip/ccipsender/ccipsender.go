package ccipsender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
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
	PackageName = "ccip-sender"
	PackageID   = "31e52b4079670904a32788d7765fa4948eae990eb423b1220ec1cd80fe3e1128"
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

// CCIPSender is a Template type
type CCIPSender struct {
	InstanceId types.TEXT  `json:"instanceId"`
	Owner      types.PARTY `json:"owner"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPSender) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPSender) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPSender) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPSender) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPSender) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSender) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSender to hex string (Canton MCMS format)
func (t CCIPSender) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSender from hex string (Canton MCMS format)
func (t *CCIPSender) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPSender

// Archive exercises the Archive choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPSender) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Send exercises the Send choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) Send(contractID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// SendWithPackageID exercises the Send choice using the provided package ID instead of package name
func (t CCIPSender) SendWithPackageID(contractID string, packageID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
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

// CantonExtraArgsV1 is a Record type
type CantonExtraArgsV1 struct {
	GasLimit           types.INT64                 `json:"gasLimit"`
	SenderRequiredCCVs []common.RawInstanceAddress `json:"senderRequiredCCVs"`
	ExecutorCid        types.CONTRACT_ID           `json:"executorCid"`
	ExecutorArgs       *types.TEXT                 `json:"executorArgs" hex:"optional"`
	TokenReceiver      *types.TEXT                 `json:"tokenReceiver" hex:"optional"`
	TokenArgs          types.TEXT                  `json:"tokenArgs"`
}

// ToMap converts CantonExtraArgsV1 to a map for DAML arguments
func (t CantonExtraArgsV1) ToMap() map[string]any {
	m := make(map[string]any)

	m["gasLimit"] = int64(t.GasLimit)

	m["senderRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.SenderRequiredCCVs))
		for _, e := range t.SenderRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["executorCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutorCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutorCid
	}()

	if t.ExecutorArgs != nil {
		m["executorArgs"] = map[string]any{
			"_type": "optional",
			"value": string(*t.ExecutorArgs),
		}
	} else {
		m["executorArgs"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.TokenReceiver != nil {
		m["tokenReceiver"] = map[string]any{
			"_type": "optional",
			"value": string(*t.TokenReceiver),
		}
	} else {
		m["tokenReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	m["tokenArgs"] = string(t.TokenArgs)

	return m
}

func (t CantonExtraArgsV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CantonExtraArgsV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CantonExtraArgsV1 to hex string (Canton MCMS format)
func (t CantonExtraArgsV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CantonExtraArgsV1 from hex string (Canton MCMS format)
func (t *CantonExtraArgsV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Send is a Record type
type Send struct {
	Context             common.CCIPContext                       `json:"context"`
	RouterCid           types.CONTRACT_ID                        `json:"routerCid"`
	DestChainSelector   types.NUMERIC                            `json:"destChainSelector"`
	Receiver            types.TEXT                               `json:"receiver"`
	Payload             types.TEXT                               `json:"payload"`
	FeeToken            splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	FeeTokenInput       interfaces.TokenInput                    `json:"feeTokenInput"`
	FeeTokenHoldingCids []types.CONTRACT_ID                      `json:"feeTokenHoldingCids"`
	TokenTransfer       *TokenTransferInput                      `json:"tokenTransfer" hex:"optional"`
	ExtraArgs           CantonExtraArgsV1                        `json:"extraArgs"`
	CcvSendInputs       []CCVSendInput                           `json:"ccvSendInputs"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["routerCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	m["destChainSelector"] = t.DestChainSelector

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["feeToken"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	m["feeTokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeTokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.FeeTokenInput
	}()

	m["feeTokenHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

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

	m["extraArgs"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

	m["ccvSendInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvSendInputs))
		for _, e := range t.CcvSendInputs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t Send) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Send) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Send to hex string (Canton MCMS format)
func (t Send) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Send from hex string (Canton MCMS format)
func (t *Send) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	TokenPoolCid      types.CONTRACT_ID                        `json:"tokenPoolCid"`
	TokenInput        interfaces.TokenInput                    `json:"tokenInput"`
	SenderInputCids   []types.CONTRACT_ID                      `json:"senderInputCids"`
	Amount            types.NUMERIC                            `json:"amount"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	PoolExtraContext  common.CCIPContext                       `json:"poolExtraContext"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

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

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["tokenInstrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInstrumentId
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
	Send(args Send) (*bind.EncodedChoice, error)
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

// Send encodes parameters for the Send choice.
func (e *encoder) Send(args Send) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Send", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
