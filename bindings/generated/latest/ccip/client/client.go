package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipapi"
	ccipcodec "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
	extensionapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
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
	PackageID   = "e5e7d0df77722323b656350e5d5a3a31ce1ba9f876d6af615a513131d1aade5e"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IPerPartyRouter is a DAML interface
type IPerPartyRouter interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// PerPartyRouterGetSequenceNumber executes the PerPartyRouter_GetSequenceNumber choice
	PerPartyRouterGetSequenceNumber(contractID string, args PerPartyRouterGetSequenceNumber) *model.ExerciseCommand

	// PerPartyRouterGetRequiredCCVsForSend executes the PerPartyRouter_GetRequiredCCVsForSend choice
	PerPartyRouterGetRequiredCCVsForSend(contractID string, args PerPartyRouterGetRequiredCCVsForSend) *model.ExerciseCommand

	// PerPartyRouterGetFee executes the PerPartyRouter_GetFee choice
	PerPartyRouterGetFee(contractID string, args PerPartyRouterGetFee) *model.ExerciseCommand

	// PerPartyRouterPrepareSend executes the PerPartyRouter_PrepareSend choice
	PerPartyRouterPrepareSend(contractID string, args PerPartyRouterPrepareSend) *model.ExerciseCommand

	// PerPartyRouterFinalizeFee executes the PerPartyRouter_FinalizeFee choice
	PerPartyRouterFinalizeFee(contractID string, args PerPartyRouterFinalizeFee) *model.ExerciseCommand

	// PerPartyRouterCCIPSend executes the PerPartyRouter_CCIPSend choice
	PerPartyRouterCCIPSend(contractID string, args PerPartyRouterCCIPSend) *model.ExerciseCommand

	// PerPartyRouterGetExecutionState executes the PerPartyRouter_GetExecutionState choice
	PerPartyRouterGetExecutionState(contractID string, args PerPartyRouterGetExecutionState) *model.ExerciseCommand

	// PerPartyRouterGetRequiredCCVsForExecute executes the PerPartyRouter_GetRequiredCCVsForExecute choice
	PerPartyRouterGetRequiredCCVsForExecute(contractID string, args PerPartyRouterGetRequiredCCVsForExecute) *model.ExerciseCommand

	// PerPartyRouterPrepareExecute executes the PerPartyRouter_PrepareExecute choice
	PerPartyRouterPrepareExecute(contractID string, args PerPartyRouterPrepareExecute) *model.ExerciseCommand

	// PerPartyRouterExecute executes the PerPartyRouter_Execute choice
	PerPartyRouterExecute(contractID string, args PerPartyRouterExecute) *model.ExerciseCommand
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

// CCIPSendResult is a Record type
type CCIPSendResult struct {
	Router                 types.CONTRACT_ID   `json:"router"`
	CcipMessageSent        types.CONTRACT_ID   `json:"ccipMessageSent"`
	MessageId              types.TEXT          `json:"messageId"`
	FeeChangeCids          []types.CONTRACT_ID `json:"feeChangeCids"`
	PendingFeeInstructions []types.CONTRACT_ID `json:"pendingFeeInstructions"`
}

// ToMap converts CCIPSendResult to a map for DAML arguments
func (t CCIPSendResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["ccipMessageSent"] = model.NestedToDAMLValue(t.CcipMessageSent)

	m["messageId"] = string(t.MessageId)

	m["feeChangeCids"] = func() []any {
		res := make([]any, 0, len(t.FeeChangeCids))
		for _, e := range t.FeeChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["pendingFeeInstructions"] = func() []any {
		res := make([]any, 0, len(t.PendingFeeInstructions))
		for _, e := range t.PendingFeeInstructions {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSendResult to hex string (Canton MCMS format)
func (t CCIPSendResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSendResult from hex string (Canton MCMS format)
func (t *CCIPSendResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCVExtraArg is a Record type
type CCVExtraArg struct {
	CcvAddress chainlinkapi.RawInstanceAddress `json:"ccvAddress"`
	CcvArgs    types.TEXT                      `json:"ccvArgs"`
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

// ExecuteResult is a Record type
type ExecuteResult struct {
	Router                types.CONTRACT_ID             `json:"router"`
	TokenReceiveTicket    *types.CONTRACT_ID            `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID             `json:"executionStateChanged"`
	MessageId             types.TEXT                    `json:"messageId"`
	Message               ccipcodec.MessageV1           `json:"message"`
	State                 ccipapi.MessageExecutionState `json:"state"`
}

// ToMap converts ExecuteResult to a map for DAML arguments
func (t ExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = model.NestedToDAMLValue(t.Router)

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenReceiveTicket),
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["executionStateChanged"] = model.NestedToDAMLValue(t.ExecutionStateChanged)

	m["messageId"] = string(t.MessageId)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["state"] = model.NestedToDAMLValue(t.State)

	return m
}

func (t ExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteResult to hex string (Canton MCMS format)
func (t ExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteResult from hex string (Canton MCMS format)
func (t *ExecuteResult) UnmarshalHex(data string) error {
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
	ExecutorAddress chainlinkapi.RawInstanceAddress `json:"executorAddress"`
	ExecutorArgs    types.TEXT                      `json:"executorArgs"`
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

// GetFeeResult is a Record type
type GetFeeResult struct {
	FeeTokenAmount types.NUMERIC `json:"feeTokenAmount"`
}

// ToMap converts GetFeeResult to a map for DAML arguments
func (t GetFeeResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenAmount"] = t.FeeTokenAmount

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

// PerPartyRouterView is a Record type
type PerPartyRouterView struct {
	CcipOwner types.PARTY `json:"ccipOwner"`
	Owner     types.PARTY `json:"owner"`
}

// ToMap converts PerPartyRouterView to a map for DAML arguments
func (t PerPartyRouterView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["owner"] = t.Owner.ToMap()

	return m
}

func (t PerPartyRouterView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterView to hex string (Canton MCMS format)
func (t PerPartyRouterView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterView from hex string (Canton MCMS format)
func (t *PerPartyRouterView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterCCIPSend is a Record type
type PerPartyRouterCCIPSend struct {
	Context                 splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	SendingMessageCid       types.CONTRACT_ID                          `json:"sendingMessageCid"`
	FeeTokenHoldingCids     []types.CONTRACT_ID                        `json:"feeTokenHoldingCids"`
	FeeTokenConfigCid       types.CONTRACT_ID                          `json:"feeTokenConfigCid"`
	FeeTokenTransferFactory types.CONTRACT_ID                          `json:"feeTokenTransferFactory"`
	FeeTokenExtraArgs       splice_api_token_metadata_v1.ExtraArgs     `json:"feeTokenExtraArgs"`
}

// ToMap converts PerPartyRouterCCIPSend to a map for DAML arguments
func (t PerPartyRouterCCIPSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeTokenHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["feeTokenConfigCid"] = model.NestedToDAMLValue(t.FeeTokenConfigCid)

	m["feeTokenTransferFactory"] = model.NestedToDAMLValue(t.FeeTokenTransferFactory)

	m["feeTokenExtraArgs"] = model.NestedToDAMLValue(t.FeeTokenExtraArgs)

	return m
}

func (t PerPartyRouterCCIPSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterCCIPSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterCCIPSend to hex string (Canton MCMS format)
func (t PerPartyRouterCCIPSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterCCIPSend from hex string (Canton MCMS format)
func (t *PerPartyRouterCCIPSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterExecute is a Record type
type PerPartyRouterExecute struct {
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
}

// ToMap converts PerPartyRouterExecute to a map for DAML arguments
func (t PerPartyRouterExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	return m
}

func (t PerPartyRouterExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterExecute to hex string (Canton MCMS format)
func (t PerPartyRouterExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterExecute from hex string (Canton MCMS format)
func (t *PerPartyRouterExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterFinalizeFee is a Record type
type PerPartyRouterFinalizeFee struct {
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts PerPartyRouterFinalizeFee to a map for DAML arguments
func (t PerPartyRouterFinalizeFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t PerPartyRouterFinalizeFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterFinalizeFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterFinalizeFee to hex string (Canton MCMS format)
func (t PerPartyRouterFinalizeFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterFinalizeFee from hex string (Canton MCMS format)
func (t *PerPartyRouterFinalizeFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterGetExecutionState is a Record type
type PerPartyRouterGetExecutionState struct {
	MessageHash types.TEXT                                 `json:"messageHash"`
	Context     splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller      types.PARTY                                `json:"caller"`
}

// ToMap converts PerPartyRouterGetExecutionState to a map for DAML arguments
func (t PerPartyRouterGetExecutionState) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageHash"] = string(t.MessageHash)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PerPartyRouterGetExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterGetExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterGetExecutionState to hex string (Canton MCMS format)
func (t PerPartyRouterGetExecutionState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterGetExecutionState from hex string (Canton MCMS format)
func (t *PerPartyRouterGetExecutionState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterGetFee is a Record type
type PerPartyRouterGetFee struct {
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Message           Canton2AnyMessage                          `json:"message"`
	CcvFeeQuotes      []extensionapi.CrossChainVerifierFeeQuote  `json:"ccvFeeQuotes"`
	TokenPoolFeeQuote *extensionapi.TokenPoolFeeQuote            `json:"tokenPoolFeeQuote" hex:"optional"`
	ExecutorFeeQuote  *extensionapi.ExecutorFeeQuote             `json:"executorFeeQuote" hex:"optional"`
}

// ToMap converts PerPartyRouterGetFee to a map for DAML arguments
func (t PerPartyRouterGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["ccvFeeQuotes"] = func() []any {
		res := make([]any, 0, len(t.CcvFeeQuotes))
		for _, e := range t.CcvFeeQuotes {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenPoolFeeQuote != nil {
		m["tokenPoolFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolFeeQuote),
		}
	} else {
		m["tokenPoolFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorFeeQuote != nil {
		m["executorFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorFeeQuote),
		}
	} else {
		m["executorFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t PerPartyRouterGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterGetFee to hex string (Canton MCMS format)
func (t PerPartyRouterGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterGetFee from hex string (Canton MCMS format)
func (t *PerPartyRouterGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterGetRequiredCCVsForExecute is a Record type
type PerPartyRouterGetRequiredCCVsForExecute struct {
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Message                   ccipcodec.MessageV1                        `json:"message"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	TokenPoolRequiredCCVs     []chainlinkapi.RawInstanceAddress          `json:"tokenPoolRequiredCCVs"`
}

// ToMap converts PerPartyRouterGetRequiredCCVsForExecute to a map for DAML arguments
func (t PerPartyRouterGetRequiredCCVsForExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverOptionalCCVs))
		for _, e := range t.ReceiverOptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalThreshold"] = int64(t.ReceiverOptionalThreshold)

	m["tokenPoolRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.TokenPoolRequiredCCVs))
		for _, e := range t.TokenPoolRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t PerPartyRouterGetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterGetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterGetRequiredCCVsForExecute to hex string (Canton MCMS format)
func (t PerPartyRouterGetRequiredCCVsForExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterGetRequiredCCVsForExecute from hex string (Canton MCMS format)
func (t *PerPartyRouterGetRequiredCCVsForExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterGetRequiredCCVsForSend is a Record type
type PerPartyRouterGetRequiredCCVsForSend struct {
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Message           Canton2AnyMessage                          `json:"message"`
	PoolReportedCCVs  []chainlinkapi.RawInstanceAddress          `json:"poolReportedCCVs"`
}

// ToMap converts PerPartyRouterGetRequiredCCVsForSend to a map for DAML arguments
func (t PerPartyRouterGetRequiredCCVsForSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["poolReportedCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolReportedCCVs))
		for _, e := range t.PoolReportedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t PerPartyRouterGetRequiredCCVsForSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterGetRequiredCCVsForSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterGetRequiredCCVsForSend to hex string (Canton MCMS format)
func (t PerPartyRouterGetRequiredCCVsForSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterGetRequiredCCVsForSend from hex string (Canton MCMS format)
func (t *PerPartyRouterGetRequiredCCVsForSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterGetSequenceNumber is a Record type
type PerPartyRouterGetSequenceNumber struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts PerPartyRouterGetSequenceNumber to a map for DAML arguments
func (t PerPartyRouterGetSequenceNumber) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t PerPartyRouterGetSequenceNumber) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterGetSequenceNumber) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterGetSequenceNumber to hex string (Canton MCMS format)
func (t PerPartyRouterGetSequenceNumber) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterGetSequenceNumber from hex string (Canton MCMS format)
func (t *PerPartyRouterGetSequenceNumber) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterPrepareExecute is a Record type
type PerPartyRouterPrepareExecute struct {
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	EncodedMessage            types.TEXT                                 `json:"encodedMessage"`
	ReceiverParty             types.PARTY                                `json:"receiverParty"`
	TokenReceiverParty        *types.PARTY                               `json:"tokenReceiverParty" hex:"optional"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	ReceiverFinalityConfig    ccipcodec.FinalityConfig                   `json:"receiverFinalityConfig"`
	Caller                    types.PARTY                                `json:"caller"`
}

// ToMap converts PerPartyRouterPrepareExecute to a map for DAML arguments
func (t PerPartyRouterPrepareExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	if t.TokenReceiverParty != nil {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
			"value": (*t.TokenReceiverParty).ToMap(),
		}
	} else {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverOptionalCCVs))
		for _, e := range t.ReceiverOptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalThreshold"] = int64(t.ReceiverOptionalThreshold)

	m["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PerPartyRouterPrepareExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterPrepareExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterPrepareExecute to hex string (Canton MCMS format)
func (t PerPartyRouterPrepareExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterPrepareExecute from hex string (Canton MCMS format)
func (t *PerPartyRouterPrepareExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterPrepareSend is a Record type
type PerPartyRouterPrepareSend struct {
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  Canton2AnyMessage                          `json:"message"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts PerPartyRouterPrepareSend to a map for DAML arguments
func (t PerPartyRouterPrepareSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t PerPartyRouterPrepareSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterPrepareSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterPrepareSend to hex string (Canton MCMS format)
func (t PerPartyRouterPrepareSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterPrepareSend from hex string (Canton MCMS format)
func (t *PerPartyRouterPrepareSend) UnmarshalHex(data string) error {
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

// IPerPartyRouterInterfaceID returns the interface ID for the IPerPartyRouter interface using the package name
func IPerPartyRouterInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ClientV1", "PerPartyRouter")
}

// IPerPartyRouterInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IPerPartyRouterInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.ClientV1", "PerPartyRouter")
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
