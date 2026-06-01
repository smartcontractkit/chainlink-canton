package sender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	core "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/core"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/chainlink/chainlinkapi"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageID   = "3f62fff74aad2b68498e486990773cd4367b3960fac8cbdf48783a8c4ae5fda5"
	SDKVersion  = "3.4.11"
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

// GetFee exercises the GetFee choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) GetFee(contractID string, args GetFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t CCIPSender) GetFeeWithPackageID(contractID string, packageID string, args GetFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

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

// GetRequiredCCVs exercises the GetRequiredCCVs choice on this CCIPSender contract
// This method uses the package name in the template ID
func (t CCIPSender) GetRequiredCCVs(contractID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsWithPackageID exercises the GetRequiredCCVs choice using the provided package ID instead of package name
func (t CCIPSender) GetRequiredCCVsWithPackageID(contractID string, packageID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPSender", "CCIPSender"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// CCVSendInput is a Record type
type CCVSendInput struct {
	CcvAddress      chainlinkapi.RawInstanceAddress            `json:"ccvAddress"`
	CcvCid          types.CONTRACT_ID                          `json:"ccvCid"`
	CcvExtraContext splice_api_token_metadata_v1.ChoiceContext `json:"ccvExtraContext"`
}

// ToMap converts CCVSendInput to a map for DAML arguments
func (t CCVSendInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvAddress"] = model.NestedToDAMLValue(t.CcvAddress)

	m["ccvCid"] = model.NestedToDAMLValue(t.CcvCid)

	m["ccvExtraContext"] = model.NestedToDAMLValue(t.CcvExtraContext)

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

// ExecutorInput is a Record type
type ExecutorInput struct {
	ExecutorCid          types.CONTRACT_ID                          `json:"executorCid"`
	ExecutorExtraContext splice_api_token_metadata_v1.ChoiceContext `json:"executorExtraContext"`
}

// ToMap converts ExecutorInput to a map for DAML arguments
func (t ExecutorInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorCid"] = model.NestedToDAMLValue(t.ExecutorCid)

	m["executorExtraContext"] = model.NestedToDAMLValue(t.ExecutorExtraContext)

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

// FeeTokenInput is a Record type
type FeeTokenInput struct {
	SenderInputCids         []types.CONTRACT_ID                    `json:"senderInputCids"`
	FeeTokenConfigCid       types.CONTRACT_ID                      `json:"feeTokenConfigCid"`
	FeeTokenTransferFactory types.CONTRACT_ID                      `json:"feeTokenTransferFactory"`
	FeeTokenExtraArgs       splice_api_token_metadata_v1.ExtraArgs `json:"feeTokenExtraArgs"`
}

// ToMap converts FeeTokenInput to a map for DAML arguments
func (t FeeTokenInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["feeTokenConfigCid"] = model.NestedToDAMLValue(t.FeeTokenConfigCid)

	m["feeTokenTransferFactory"] = model.NestedToDAMLValue(t.FeeTokenTransferFactory)

	m["feeTokenExtraArgs"] = model.NestedToDAMLValue(t.FeeTokenExtraArgs)

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

// GetFee2 is a Record type
type GetFee2 struct {
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  core.Canton2AnyMessage                     `json:"message"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	CcvSendInputs            []CCVSendInput                             `json:"ccvSendInputs"`
	TokenTransferInput       *TokenTransferInput                        `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *ExecutorInput                             `json:"executorInput" hex:"optional"`
}

// ToMap converts GetFee2 to a map for DAML arguments
func (t GetFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["ccvSendInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvSendInputs))
		for _, e := range t.CcvSendInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenTransferInput != nil {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransferInput),
		}
	} else {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorInput != nil {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorInput),
		}
	} else {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t GetFee2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee2 to hex string (Canton MCMS format)
func (t GetFee2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee2 from hex string (Canton MCMS format)
func (t *GetFee2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeResult is a Record type
type GetFeeResult struct {
	FeeTokenAmount     types.NUMERIC `json:"feeTokenAmount"`
	PoolFeeTokenAmount types.NUMERIC `json:"poolFeeTokenAmount"`
}

// ToMap converts GetFeeResult to a map for DAML arguments
func (t GetFeeResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenAmount"] = t.FeeTokenAmount

	m["poolFeeTokenAmount"] = t.PoolFeeTokenAmount

	return m
}

func (t GetFeeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeResult to hex string (Canton MCMS format)
func (t GetFeeResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeResult from hex string (Canton MCMS format)
func (t *GetFeeResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVs is a Record type
type GetRequiredCCVs struct {
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  core.Canton2AnyMessage                     `json:"message"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	TokenPoolCid             *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	if t.TokenPoolCid != nil {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolCid),
		}
	} else {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t GetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVs to hex string (Canton MCMS format)
func (t GetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVs from hex string (Canton MCMS format)
func (t *GetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Send is a Record type
type Send struct {
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  core.Canton2AnyMessage                     `json:"message"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	FeeTokenInput            FeeTokenInput                              `json:"feeTokenInput"`
	CcvSendInputs            []CCVSendInput                             `json:"ccvSendInputs"`
	TokenTransferInput       *TokenTransferInput                        `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *ExecutorInput                             `json:"executorInput" hex:"optional"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["feeTokenInput"] = model.NestedToDAMLValue(t.FeeTokenInput)

	m["ccvSendInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvSendInputs))
		for _, e := range t.CcvSendInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenTransferInput != nil {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransferInput),
		}
	} else {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorInput != nil {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorInput),
		}
	} else {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
	SenderInputCids  []types.CONTRACT_ID                        `json:"senderInputCids"`
	TokenPoolCid     types.CONTRACT_ID                          `json:"tokenPoolCid"`
	PoolExtraContext splice_api_token_metadata_v1.ChoiceContext `json:"poolExtraContext"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["tokenPoolCid"] = model.NestedToDAMLValue(t.TokenPoolCid)

	m["poolExtraContext"] = model.NestedToDAMLValue(t.PoolExtraContext)

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
	GetFee(args GetFee2) (*bind.EncodedChoice, error)
	GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error)
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

// GetFee encodes parameters for the GetFee choice.
func (e *encoder) GetFee(args GetFee2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetRequiredCCVs encodes parameters for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// Send encodes parameters for the Send choice.
func (e *encoder) Send(args Send) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Send", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
