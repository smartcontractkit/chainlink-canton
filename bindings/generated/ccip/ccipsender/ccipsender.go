package ccipsender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
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
	PackageID   = "0bf87bcaba9af8a980de7811823ecc62e4ceaa58529a3c87de97d3d1a3aed3c6"
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

// Send is a Record type
type Send struct {
	DestinationChainSelector types.NUMERIC            `json:"destinationChainSelector"`
	Message                  common.Canton2AnyMessage `json:"message"`
	Context                  common.CCIPContext       `json:"context"`
	RouterCid                types.CONTRACT_ID        `json:"routerCid"`
	FeeTokenInput            interfaces.TokenInput    `json:"feeTokenInput"`
	FeeTokenHoldingCids      []types.CONTRACT_ID      `json:"feeTokenHoldingCids"`
	TokenTransfer            *TokenTransferInput      `json:"tokenTransfer" hex:"optional"`
	CcvSendInputs            types.GENMAP             `json:"ccvSendInputs"`
	ExecutorCid              *types.CONTRACT_ID       `json:"executorCid" hex:"optional"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

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

	m["ccvSendInputs"] = func() any {
		if t.CcvSendInputs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.CcvSendInputs}
	}()

	if t.ExecutorCid != nil {
		m["executorCid"] = map[string]any{
			"_type": "optional",
			"value": *t.ExecutorCid,
		}
	} else {
		m["executorCid"] = map[string]any{
			"_type": "optional",
		}
	}

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
	TokenPoolCid     types.CONTRACT_ID     `json:"tokenPoolCid"`
	TokenInput       interfaces.TokenInput `json:"tokenInput"`
	PoolExtraContext common.CCIPContext    `json:"poolExtraContext"`
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
