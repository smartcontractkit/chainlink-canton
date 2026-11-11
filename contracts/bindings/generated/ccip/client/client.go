package client

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipapi"
	ccipcodec "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	clientapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	extensionapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "ccip-client-v2"
	PackageID   = "60b61aade24f95d14567a45ed5bfdf676109da5db8cb9221a1b287468f85c3a7"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IICCIPReceiver is a DAML interface
type IICCIPReceiver interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// ICCIPReceiverGetRequiredCCVs executes the ICCIPReceiver_GetRequiredCCVs choice
	ICCIPReceiverGetRequiredCCVs(contractID string, args ICCIPReceiverGetRequiredCCVs) *model.ExerciseCommand

	// ICCIPReceiverExecute executes the ICCIPReceiver_Execute choice
	ICCIPReceiverExecute(contractID string, args ICCIPReceiverExecute) *model.ExerciseCommand
}

// IICCIPSender is a DAML interface
type IICCIPSender interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// ICCIPSenderGetRequiredCCVs executes the ICCIPSender_GetRequiredCCVs choice
	ICCIPSenderGetRequiredCCVs(contractID string, args ICCIPSenderGetRequiredCCVs) *model.ExerciseCommand

	// ICCIPSenderGetFee executes the ICCIPSender_GetFee choice
	ICCIPSenderGetFee(contractID string, args ICCIPSenderGetFee) *model.ExerciseCommand

	// ICCIPSenderSend executes the ICCIPSender_Send choice
	ICCIPSenderSend(contractID string, args ICCIPSenderSend) *model.ExerciseCommand
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

// ICCIPReceiverView is a Record type
type ICCIPReceiverView struct {
	CcipOwner              types.PARTY                       `json:"ccipOwner"`
	RequiredCCVs           []chainlinkapi.RawInstanceAddress `json:"requiredCCVs"`
	ReceiverFinalityConfig ccipcodec.FinalityConfig          `json:"receiverFinalityConfig"`
}

// ToMap converts ICCIPReceiverView to a map for DAML arguments
func (t ICCIPReceiverView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	return m
}

func (t ICCIPReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPReceiverView to hex string (Canton MCMS format)
func (t ICCIPReceiverView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPReceiverView from hex string (Canton MCMS format)
func (t *ICCIPReceiverView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPReceiverExecute is a Record type
type ICCIPReceiverExecute struct {
	ReceiverParty  types.PARTY                                `json:"receiverParty"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenTransfer  *ReceiverTokenTransferInput                `json:"tokenTransfer" hex:"optional"`
	CcvInputs      []ReceiverCCVInput                         `json:"ccvInputs"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ICCIPReceiverExecute to a map for DAML arguments
func (t ICCIPReceiverExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

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

	m["ccvInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvInputs))
		for _, e := range t.CcvInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ICCIPReceiverExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPReceiverExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPReceiverExecute to hex string (Canton MCMS format)
func (t ICCIPReceiverExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPReceiverExecute from hex string (Canton MCMS format)
func (t *ICCIPReceiverExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPReceiverGetRequiredCCVs is a Record type
type ICCIPReceiverGetRequiredCCVs struct {
	ReceiverParty  types.PARTY                                `json:"receiverParty"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenPoolCid   *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ICCIPReceiverGetRequiredCCVs to a map for DAML arguments
func (t ICCIPReceiverGetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

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

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ICCIPReceiverGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPReceiverGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPReceiverGetRequiredCCVs to hex string (Canton MCMS format)
func (t ICCIPReceiverGetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPReceiverGetRequiredCCVs from hex string (Canton MCMS format)
func (t *ICCIPReceiverGetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPSenderView is a Record type
type ICCIPSenderView struct {
	CcipOwner types.PARTY `json:"ccipOwner"`
}

// ToMap converts ICCIPSenderView to a map for DAML arguments
func (t ICCIPSenderView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	return m
}

func (t ICCIPSenderView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPSenderView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPSenderView to hex string (Canton MCMS format)
func (t ICCIPSenderView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPSenderView from hex string (Canton MCMS format)
func (t *ICCIPSenderView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPSenderGetFee is a Record type
type ICCIPSenderGetFee struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	CcvSendInputs            []SenderCCVSendInput                       `json:"ccvSendInputs"`
	TokenTransferInput       *SenderTokenTransferInput                  `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *SenderExecutorInput                       `json:"executorInput" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ICCIPSenderGetFee to a map for DAML arguments
func (t ICCIPSenderGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

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

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ICCIPSenderGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPSenderGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPSenderGetFee to hex string (Canton MCMS format)
func (t ICCIPSenderGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPSenderGetFee from hex string (Canton MCMS format)
func (t *ICCIPSenderGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPSenderGetRequiredCCVs is a Record type
type ICCIPSenderGetRequiredCCVs struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	TokenPoolCid             *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ICCIPSenderGetRequiredCCVs to a map for DAML arguments
func (t ICCIPSenderGetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

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

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ICCIPSenderGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPSenderGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPSenderGetRequiredCCVs to hex string (Canton MCMS format)
func (t ICCIPSenderGetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPSenderGetRequiredCCVs from hex string (Canton MCMS format)
func (t *ICCIPSenderGetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ICCIPSenderSend is a Record type
type ICCIPSenderSend struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	FeeTokenInput            SenderFeeTokenInput                        `json:"feeTokenInput"`
	CcvSendInputs            []SenderCCVSendInput                       `json:"ccvSendInputs"`
	TokenTransferInput       *SenderTokenTransferInput                  `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *SenderExecutorInput                       `json:"executorInput" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ICCIPSenderSend to a map for DAML arguments
func (t ICCIPSenderSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

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

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ICCIPSenderSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ICCIPSenderSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ICCIPSenderSend to hex string (Canton MCMS format)
func (t ICCIPSenderSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ICCIPSenderSend from hex string (Canton MCMS format)
func (t *ICCIPSenderSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReceiverCCVInput is a Record type
type ReceiverCCVInput struct {
	CcvCid          types.CONTRACT_ID                          `json:"ccvCid"`
	VerifierResults types.TEXT                                 `json:"verifierResults"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ReceiverCCVInput to a map for DAML arguments
func (t ReceiverCCVInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvCid"] = model.NestedToDAMLValue(t.CcvCid)

	m["verifierResults"] = string(t.VerifierResults)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ReceiverCCVInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReceiverCCVInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReceiverCCVInput to hex string (Canton MCMS format)
func (t ReceiverCCVInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReceiverCCVInput from hex string (Canton MCMS format)
func (t *ReceiverCCVInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReceiverExecuteResult is a Record type
type ReceiverExecuteResult struct {
	Router             types.CONTRACT_ID                          `json:"router"`
	MessageId          types.TEXT                                 `json:"messageId"`
	Message            ccipcodec.MessageV1                        `json:"message"`
	State              ccipapi.MessageExecutionState              `json:"state"`
	TokenReleaseResult *extensionapi.ReleaseOrMintResult          `json:"tokenReleaseResult" hex:"optional"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ReceiverExecuteResult to a map for DAML arguments
func (t ReceiverExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["messageId"] = string(t.MessageId)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["state"] = model.NestedToDAMLValue(t.State)

	if t.TokenReleaseResult != nil {
		m["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenReleaseResult),
		}
	} else {
		m["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ReceiverExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReceiverExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReceiverExecuteResult to hex string (Canton MCMS format)
func (t ReceiverExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReceiverExecuteResult from hex string (Canton MCMS format)
func (t *ReceiverExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReceiverTokenTransferInput is a Record type
type ReceiverTokenTransferInput struct {
	TokenPoolCid       types.CONTRACT_ID                          `json:"tokenPoolCid"`
	TokenReceiverParty types.PARTY                                `json:"tokenReceiverParty"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ReceiverTokenTransferInput to a map for DAML arguments
func (t ReceiverTokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenPoolCid"] = model.NestedToDAMLValue(t.TokenPoolCid)

	m["tokenReceiverParty"] = t.TokenReceiverParty.ToMap()

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ReceiverTokenTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReceiverTokenTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReceiverTokenTransferInput to hex string (Canton MCMS format)
func (t ReceiverTokenTransferInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReceiverTokenTransferInput from hex string (Canton MCMS format)
func (t *ReceiverTokenTransferInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SenderCCVSendInput is a Record type
type SenderCCVSendInput struct {
	CcvAddress chainlinkapi.RawInstanceAddress            `json:"ccvAddress"`
	CcvCid     types.CONTRACT_ID                          `json:"ccvCid"`
	Context    splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts SenderCCVSendInput to a map for DAML arguments
func (t SenderCCVSendInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvAddress"] = model.NestedToDAMLValue(t.CcvAddress)

	m["ccvCid"] = model.NestedToDAMLValue(t.CcvCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t SenderCCVSendInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SenderCCVSendInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SenderCCVSendInput to hex string (Canton MCMS format)
func (t SenderCCVSendInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SenderCCVSendInput from hex string (Canton MCMS format)
func (t *SenderCCVSendInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SenderExecutorInput is a Record type
type SenderExecutorInput struct {
	ExecutorCid types.CONTRACT_ID                          `json:"executorCid"`
	Context     splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts SenderExecutorInput to a map for DAML arguments
func (t SenderExecutorInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorCid"] = model.NestedToDAMLValue(t.ExecutorCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t SenderExecutorInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SenderExecutorInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SenderExecutorInput to hex string (Canton MCMS format)
func (t SenderExecutorInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SenderExecutorInput from hex string (Canton MCMS format)
func (t *SenderExecutorInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SenderFeeTokenInput is a Record type
type SenderFeeTokenInput struct {
	SenderInputCids         []types.CONTRACT_ID                    `json:"senderInputCids"`
	FeeTokenConfigCid       types.CONTRACT_ID                      `json:"feeTokenConfigCid"`
	FeeTokenTransferFactory types.CONTRACT_ID                      `json:"feeTokenTransferFactory"`
	FeeTokenExtraArgs       splice_api_token_metadata_v1.ExtraArgs `json:"feeTokenExtraArgs"`
}

// ToMap converts SenderFeeTokenInput to a map for DAML arguments
func (t SenderFeeTokenInput) ToMap() map[string]any {
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

func (t SenderFeeTokenInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SenderFeeTokenInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SenderFeeTokenInput to hex string (Canton MCMS format)
func (t SenderFeeTokenInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SenderFeeTokenInput from hex string (Canton MCMS format)
func (t *SenderFeeTokenInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SenderGetFeeResult is a Record type
type SenderGetFeeResult struct {
	FeeResult clientapi.GetFeeResult          `json:"feeResult"`
	PoolQuote *extensionapi.TokenPoolFeeQuote `json:"poolQuote" hex:"optional"`
}

// ToMap converts SenderGetFeeResult to a map for DAML arguments
func (t SenderGetFeeResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeResult"] = model.NestedToDAMLValue(t.FeeResult)

	if t.PoolQuote != nil {
		m["poolQuote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.PoolQuote),
		}
	} else {
		m["poolQuote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SenderGetFeeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SenderGetFeeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SenderGetFeeResult to hex string (Canton MCMS format)
func (t SenderGetFeeResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SenderGetFeeResult from hex string (Canton MCMS format)
func (t *SenderGetFeeResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SenderTokenTransferInput is a Record type
type SenderTokenTransferInput struct {
	SenderInputCids []types.CONTRACT_ID                        `json:"senderInputCids"`
	TokenPoolCid    types.CONTRACT_ID                          `json:"tokenPoolCid"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts SenderTokenTransferInput to a map for DAML arguments
func (t SenderTokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["tokenPoolCid"] = model.NestedToDAMLValue(t.TokenPoolCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t SenderTokenTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SenderTokenTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SenderTokenTransferInput to hex string (Canton MCMS format)
func (t SenderTokenTransferInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SenderTokenTransferInput from hex string (Canton MCMS format)
func (t *SenderTokenTransferInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IICCIPReceiverInterfaceID returns the interface ID for the IICCIPReceiver interface using the package name
func IICCIPReceiverInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Client.Receiver", "ICCIPReceiver")
}

// IICCIPReceiverInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IICCIPReceiverInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Client.Receiver", "ICCIPReceiver")
}

// IICCIPSenderInterfaceID returns the interface ID for the IICCIPSender interface using the package name
func IICCIPSenderInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Client.Sender", "ICCIPSender")
}

// IICCIPSenderInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IICCIPSenderInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Client.Sender", "ICCIPSender")
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
