package extensionapi

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	core "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
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
	PackageName = "ccip-extension-api"
	PackageID   = "f2556f0df8057ffc942153e8d0b4efa1f420b9fb05d49e4401df7689eb04ca9b"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IICrossChainVerifier is a DAML interface
type IICrossChainVerifier interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// CrossChainVerifierVerifyMessage executes the CrossChainVerifier_VerifyMessage choice
	CrossChainVerifierVerifyMessage(contractID string, args CrossChainVerifierVerifyMessage) *model.ExerciseCommand

	// CrossChainVerifierCalculateFee executes the CrossChainVerifier_CalculateFee choice
	CrossChainVerifierCalculateFee(contractID string, args CrossChainVerifierCalculateFee) *model.ExerciseCommand

	// CrossChainVerifierGetFee executes the CrossChainVerifier_GetFee choice
	CrossChainVerifierGetFee(contractID string, args CrossChainVerifierGetFee) *model.ExerciseCommand

	// CrossChainVerifierForwardToVerifier executes the CrossChainVerifier_ForwardToVerifier choice
	CrossChainVerifierForwardToVerifier(contractID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand
}

// IIExecutor is a DAML interface
type IIExecutor interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// ExecutorCalculateFee executes the Executor_CalculateFee choice
	ExecutorCalculateFee(contractID string, args ExecutorCalculateFee) *model.ExerciseCommand

	// ExecutorGetFee executes the Executor_GetFee choice
	ExecutorGetFee(contractID string, args ExecutorGetFee) *model.ExerciseCommand
}

// IITokenPool is a DAML interface
type IITokenPool interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenPoolGetRequiredCCVs executes the TokenPool_GetRequiredCCVs choice
	TokenPoolGetRequiredCCVs(contractID string, args TokenPoolGetRequiredCCVs) *model.ExerciseCommand

	// TokenPoolVerifyInboundMessage executes the TokenPool_VerifyInboundMessage choice
	TokenPoolVerifyInboundMessage(contractID string, args TokenPoolVerifyInboundMessage) *model.ExerciseCommand

	// TokenPoolVerifyOutboundCCVs executes the TokenPool_VerifyOutboundCCVs choice
	TokenPoolVerifyOutboundCCVs(contractID string, args TokenPoolVerifyOutboundCCVs) *model.ExerciseCommand

	// TokenPoolReleaseFromTicket executes the TokenPool_ReleaseFromTicket choice
	TokenPoolReleaseFromTicket(contractID string, args TokenPoolReleaseFromTicket) *model.ExerciseCommand

	// TokenPoolLockOrBurn executes the TokenPool_LockOrBurn choice
	TokenPoolLockOrBurn(contractID string, args TokenPoolLockOrBurn) *model.ExerciseCommand

	// TokenPoolCalculateFee executes the TokenPool_CalculateFee choice
	TokenPoolCalculateFee(contractID string, args TokenPoolCalculateFee) *model.ExerciseCommand

	// TokenPoolGetFee executes the TokenPool_GetFee choice
	TokenPoolGetFee(contractID string, args TokenPoolGetFee) *model.ExerciseCommand
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

// CrossChainVerifierFeeQuote is a Record type
type CrossChainVerifierFeeQuote struct {
	CcvInstanceId      types.TEXT    `json:"ccvInstanceId"`
	CcvOwner           types.PARTY   `json:"ccvOwner"`
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
	GasForVerification types.INT64   `json:"gasForVerification"`
	PayloadSizeBytes   types.INT64   `json:"payloadSizeBytes"`
}

// ToMap converts CrossChainVerifierFeeQuote to a map for DAML arguments
func (t CrossChainVerifierFeeQuote) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["gasForVerification"] = int64(t.GasForVerification)

	m["payloadSizeBytes"] = int64(t.PayloadSizeBytes)

	return m
}

func (t CrossChainVerifierFeeQuote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierFeeQuote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierFeeQuote to hex string (Canton MCMS format)
func (t CrossChainVerifierFeeQuote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierFeeQuote from hex string (Canton MCMS format)
func (t *CrossChainVerifierFeeQuote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	InstanceId       types.TEXT   `json:"instanceId"`
	Owner            types.PARTY  `json:"owner"`
	CcipOwner        types.PARTY  `json:"ccipOwner"`
	StorageLocations []types.TEXT `json:"storageLocations"`
}

// ToMap converts CrossChainVerifierView to a map for DAML arguments
func (t CrossChainVerifierView) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CrossChainVerifierView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierView to hex string (Canton MCMS format)
func (t CrossChainVerifierView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierView from hex string (Canton MCMS format)
func (t *CrossChainVerifierView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierCalculateFee is a Record type
type CrossChainVerifierCalculateFee struct {
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts CrossChainVerifierCalculateFee to a map for DAML arguments
func (t CrossChainVerifierCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CrossChainVerifierCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierCalculateFee to hex string (Canton MCMS format)
func (t CrossChainVerifierCalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierCalculateFee from hex string (Canton MCMS format)
func (t *CrossChainVerifierCalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierForwardToVerifier is a Record type
type CrossChainVerifierForwardToVerifier struct {
	RmnRemoteCid      types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT                                 `json:"verifierArgs"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts CrossChainVerifierForwardToVerifier to a map for DAML arguments
func (t CrossChainVerifierForwardToVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["verifierArgs"] = string(t.VerifierArgs)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CrossChainVerifierForwardToVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierForwardToVerifier to hex string (Canton MCMS format)
func (t CrossChainVerifierForwardToVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierForwardToVerifier from hex string (Canton MCMS format)
func (t *CrossChainVerifierForwardToVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierGetFee is a Record type
type CrossChainVerifierGetFee struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts CrossChainVerifierGetFee to a map for DAML arguments
func (t CrossChainVerifierGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CrossChainVerifierGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierGetFee to hex string (Canton MCMS format)
func (t CrossChainVerifierGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierGetFee from hex string (Canton MCMS format)
func (t *CrossChainVerifierGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierVerifyMessage is a Record type
type CrossChainVerifierVerifyMessage struct {
	RmnRemoteCid        types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext        splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	VerifierResults     types.TEXT                                 `json:"verifierResults"`
	Caller              types.PARTY                                `json:"caller"`
}

// ToMap converts CrossChainVerifierVerifyMessage to a map for DAML arguments
func (t CrossChainVerifierVerifyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["verifierResults"] = string(t.VerifierResults)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CrossChainVerifierVerifyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CrossChainVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CrossChainVerifierVerifyMessage to hex string (Canton MCMS format)
func (t CrossChainVerifierVerifyMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CrossChainVerifierVerifyMessage from hex string (Canton MCMS format)
func (t *CrossChainVerifierVerifyMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorFeeQuote is a Record type
type ExecutorFeeQuote struct {
	ExecutorInstanceId types.TEXT    `json:"executorInstanceId"`
	ExecutorOwner      types.PARTY   `json:"executorOwner"`
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
}

// ToMap converts ExecutorFeeQuote to a map for DAML arguments
func (t ExecutorFeeQuote) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorInstanceId"] = string(t.ExecutorInstanceId)

	m["executorOwner"] = t.ExecutorOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	return m
}

func (t ExecutorFeeQuote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorFeeQuote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorFeeQuote to hex string (Canton MCMS format)
func (t ExecutorFeeQuote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorFeeQuote from hex string (Canton MCMS format)
func (t *ExecutorFeeQuote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorView is a Record type
type ExecutorView struct {
	InstanceId types.TEXT  `json:"instanceId"`
	Owner      types.PARTY `json:"owner"`
}

// ToMap converts ExecutorView to a map for DAML arguments
func (t ExecutorView) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	return m
}

func (t ExecutorView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorView to hex string (Canton MCMS format)
func (t ExecutorView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorView from hex string (Canton MCMS format)
func (t *ExecutorView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorCalculateFee is a Record type
type ExecutorCalculateFee struct {
	ExpectedExecutor  chainlinkapi.RawInstanceAddress            `json:"expectedExecutor"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	ExecutorArgs      types.TEXT                                 `json:"executorArgs"`
	ExtraContext      splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts ExecutorCalculateFee to a map for DAML arguments
func (t ExecutorCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedExecutor"] = model.NestedToDAMLValue(t.ExpectedExecutor)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorCalculateFee to hex string (Canton MCMS format)
func (t ExecutorCalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorCalculateFee from hex string (Canton MCMS format)
func (t *ExecutorCalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetFee is a Record type
type ExecutorGetFee struct {
	ExpectedExecutor  chainlinkapi.RawInstanceAddress   `json:"expectedExecutor"`
	DestChainSelector types.NUMERIC                     `json:"destChainSelector"`
	RequiredCCVs      []chainlinkapi.RawInstanceAddress `json:"requiredCCVs"`
	Caller            types.PARTY                       `json:"caller"`
}

// ToMap converts ExecutorGetFee to a map for DAML arguments
func (t ExecutorGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedExecutor"] = model.NestedToDAMLValue(t.ExpectedExecutor)

	m["destChainSelector"] = t.DestChainSelector

	m["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetFee to hex string (Canton MCMS format)
func (t ExecutorGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetFee from hex string (Canton MCMS format)
func (t *ExecutorGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOrBurnResult is a Record type
type LockOrBurnResult struct {
	PoolChangeCids    []types.CONTRACT_ID `json:"poolChangeCids"`
	SenderChangeCids  []types.CONTRACT_ID `json:"senderChangeCids"`
	SendingMessageCid types.CONTRACT_ID   `json:"sendingMessageCid"`
}

// ToMap converts LockOrBurnResult to a map for DAML arguments
func (t LockOrBurnResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolChangeCids"] = func() []any {
		res := make([]any, 0, len(t.PoolChangeCids))
		for _, e := range t.PoolChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["senderChangeCids"] = func() []any {
		res := make([]any, 0, len(t.SenderChangeCids))
		for _, e := range t.SenderChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	return m
}

func (t LockOrBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOrBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOrBurnResult to hex string (Canton MCMS format)
func (t LockOrBurnResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOrBurnResult from hex string (Canton MCMS format)
func (t *LockOrBurnResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReleaseOrMintResult is a Record type
type ReleaseOrMintResult struct {
	Output          ReleaseOrMintResultOutput `json:"output"`
	PoolChangeCids  []types.CONTRACT_ID       `json:"poolChangeCids"`
	ClaimedEventCid types.CONTRACT_ID         `json:"claimedEventCid"`
}

// ToMap converts ReleaseOrMintResult to a map for DAML arguments
func (t ReleaseOrMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["output"] = model.NestedToDAMLValue(t.Output)

	m["poolChangeCids"] = func() []any {
		res := make([]any, 0, len(t.PoolChangeCids))
		for _, e := range t.PoolChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["claimedEventCid"] = model.NestedToDAMLValue(t.ClaimedEventCid)

	return m
}

func (t ReleaseOrMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReleaseOrMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReleaseOrMintResult to hex string (Canton MCMS format)
func (t ReleaseOrMintResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReleaseOrMintResult from hex string (Canton MCMS format)
func (t *ReleaseOrMintResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReleaseOrMintResultCompleted is a Record type
type ReleaseOrMintResultCompleted struct {
	ReceiverHoldingCids []types.CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts ReleaseOrMintResultCompleted to a map for DAML arguments
func (t ReleaseOrMintResultCompleted) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverHoldingCids))
		for _, e := range t.ReceiverHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t ReleaseOrMintResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReleaseOrMintResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReleaseOrMintResultCompleted to hex string (Canton MCMS format)
func (t ReleaseOrMintResultCompleted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReleaseOrMintResultCompleted from hex string (Canton MCMS format)
func (t *ReleaseOrMintResultCompleted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ReleaseOrMintResultOutput is a variant/union type
type ReleaseOrMintResultOutput struct {
	ReleaseOrMintResultPending   *ReleaseOrMintResultPending   `json:"ReleaseOrMintResult_Pending,omitempty"`
	ReleaseOrMintResultCompleted *ReleaseOrMintResultCompleted `json:"ReleaseOrMintResult_Completed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ReleaseOrMintResultOutput
func (v ReleaseOrMintResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for ReleaseOrMintResultOutput
func (v *ReleaseOrMintResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes ReleaseOrMintResultOutput to hex string (Canton MCMS format)
func (v ReleaseOrMintResultOutput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes ReleaseOrMintResultOutput from hex string (Canton MCMS format)
func (v *ReleaseOrMintResultOutput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ReleaseOrMintResultOutput) GetVariantTag() string {

	if v.ReleaseOrMintResultPending != nil {
		return "ReleaseOrMintResult_Pending"
	}

	if v.ReleaseOrMintResultCompleted != nil {
		return "ReleaseOrMintResult_Completed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ReleaseOrMintResultOutput) GetVariantValue() any {

	if v.ReleaseOrMintResultPending != nil {
		return v.ReleaseOrMintResultPending
	}

	if v.ReleaseOrMintResultCompleted != nil {
		return v.ReleaseOrMintResultCompleted
	}

	return nil
}

var _ types.VARIANT = (*ReleaseOrMintResultOutput)(nil)

// ReleaseOrMintResultPending is a Record type
type ReleaseOrMintResultPending struct {
	TransferInstructionCid types.CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts ReleaseOrMintResultPending to a map for DAML arguments
func (t ReleaseOrMintResultPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferInstructionCid"] = model.NestedToDAMLValue(t.TransferInstructionCid)

	return m
}

func (t ReleaseOrMintResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ReleaseOrMintResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ReleaseOrMintResultPending to hex string (Canton MCMS format)
func (t ReleaseOrMintResultPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ReleaseOrMintResultPending from hex string (Canton MCMS format)
func (t *ReleaseOrMintResultPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolFeeQuote is a Record type
type TokenPoolFeeQuote struct {
	PoolInstanceId    types.TEXT    `json:"poolInstanceId"`
	PoolOwner         types.PARTY   `json:"poolOwner"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	TokenFeeBps       types.NUMERIC `json:"tokenFeeBps"`
	IsEnabled         types.BOOL    `json:"isEnabled"`
}

// ToMap converts TokenPoolFeeQuote to a map for DAML arguments
func (t TokenPoolFeeQuote) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["tokenFeeBps"] = t.TokenFeeBps

	m["isEnabled"] = bool(t.IsEnabled)

	return m
}

func (t TokenPoolFeeQuote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolFeeQuote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolFeeQuote to hex string (Canton MCMS format)
func (t TokenPoolFeeQuote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolFeeQuote from hex string (Canton MCMS format)
func (t *TokenPoolFeeQuote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolView is a Record type
type TokenPoolView struct {
	Owner        types.PARTY                              `json:"owner"`
	CcipOwner    types.PARTY                              `json:"ccipOwner"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// ToMap converts TokenPoolView to a map for DAML arguments
func (t TokenPoolView) ToMap() map[string]any {
	m := make(map[string]any)

	m["owner"] = t.Owner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	return m
}

func (t TokenPoolView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolView to hex string (Canton MCMS format)
func (t TokenPoolView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolView from hex string (Canton MCMS format)
func (t *TokenPoolView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolCalculateFee is a Record type
type TokenPoolCalculateFee struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	FeeQuoterCid          types.CONTRACT_ID                          `json:"feeQuoterCid"`
	TokenInstrumentId     splice_api_token_holding_v1.InstrumentId   `json:"tokenInstrumentId"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenPoolCalculateFee to a map for DAML arguments
func (t TokenPoolCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolCalculateFee to hex string (Canton MCMS format)
func (t TokenPoolCalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolCalculateFee from hex string (Canton MCMS format)
func (t *TokenPoolCalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolGetFee is a Record type
type TokenPoolGetFee struct {
	FeeQuoterCid      types.CONTRACT_ID                        `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	TokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts TokenPoolGetFee to a map for DAML arguments
func (t TokenPoolGetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["destChainSelector"] = t.DestChainSelector

	m["tokenInstrumentId"] = model.NestedToDAMLValue(t.TokenInstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolGetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolGetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolGetFee to hex string (Canton MCMS format)
func (t TokenPoolGetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolGetFee from hex string (Canton MCMS format)
func (t *TokenPoolGetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolGetRequiredCCVs is a Record type
type TokenPoolGetRequiredCCVs struct {
	RemoteChainSelector types.NUMERIC       `json:"remoteChainSelector"`
	SourceAmount        types.TEXT          `json:"sourceAmount"`
	Finality            core.FinalityConfig `json:"finality"`
	ExtraData           types.TEXT          `json:"extraData"`
	Direction           TransferDirection   `json:"direction"`
	Caller              types.PARTY         `json:"caller"`
}

// ToMap converts TokenPoolGetRequiredCCVs to a map for DAML arguments
func (t TokenPoolGetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["sourceAmount"] = string(t.SourceAmount)

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = model.NestedToDAMLValue(t.Direction)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolGetRequiredCCVs to hex string (Canton MCMS format)
func (t TokenPoolGetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolGetRequiredCCVs from hex string (Canton MCMS format)
func (t *TokenPoolGetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolLockOrBurn is a Record type
type TokenPoolLockOrBurn struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	SenderInputCids       []types.CONTRACT_ID                        `json:"senderInputCids"`
	Amount                types.NUMERIC                              `json:"amount"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenPoolLockOrBurn to a map for DAML arguments
func (t TokenPoolLockOrBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["senderInputCids"] = func() []any {
		res := make([]any, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolLockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolLockOrBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolLockOrBurn to hex string (Canton MCMS format)
func (t TokenPoolLockOrBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolLockOrBurn from hex string (Canton MCMS format)
func (t *TokenPoolLockOrBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolReleaseFromTicket is a Record type
type TokenPoolReleaseFromTicket struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	RmnRemoteCid          types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	TokenReceiveTicketCid types.CONTRACT_ID                          `json:"tokenReceiveTicketCid"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenPoolReleaseFromTicket to a map for DAML arguments
func (t TokenPoolReleaseFromTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolReleaseFromTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolReleaseFromTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolReleaseFromTicket to hex string (Canton MCMS format)
func (t TokenPoolReleaseFromTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolReleaseFromTicket from hex string (Canton MCMS format)
func (t *TokenPoolReleaseFromTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolVerifyInboundMessage is a Record type
type TokenPoolVerifyInboundMessage struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	ExecutingMessageCid   types.CONTRACT_ID                          `json:"executingMessageCid"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenPoolVerifyInboundMessage to a map for DAML arguments
func (t TokenPoolVerifyInboundMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolVerifyInboundMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolVerifyInboundMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolVerifyInboundMessage to hex string (Canton MCMS format)
func (t TokenPoolVerifyInboundMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolVerifyInboundMessage from hex string (Canton MCMS format)
func (t *TokenPoolVerifyInboundMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenPoolVerifyOutboundCCVs is a Record type
type TokenPoolVerifyOutboundCCVs struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext `json:"extraContext"`
	SendingMessageCid     types.CONTRACT_ID                          `json:"sendingMessageCid"`
	Amount                types.NUMERIC                              `json:"amount"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenPoolVerifyOutboundCCVs to a map for DAML arguments
func (t TokenPoolVerifyOutboundCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["extraContext"] = model.NestedToDAMLValue(t.ExtraContext)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["amount"] = t.Amount

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolVerifyOutboundCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPoolVerifyOutboundCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPoolVerifyOutboundCCVs to hex string (Canton MCMS format)
func (t TokenPoolVerifyOutboundCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPoolVerifyOutboundCCVs from hex string (Canton MCMS format)
func (t *TokenPoolVerifyOutboundCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferDirection is an enum type
type TransferDirection string

const (
	TransferDirectionOutbound TransferDirection = "Outbound"

	TransferDirectionInbound TransferDirection = "Inbound"
)

func (e TransferDirection) GetEnumConstructor() string { return string(e) }

func (e TransferDirection) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.TokenPool", "TransferDirection")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e TransferDirection) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Interfaces.TokenPool", "TransferDirection")
}

func (e TransferDirection) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *TransferDirection) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes TransferDirection to hex string (Canton MCMS format)
func (e TransferDirection) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes TransferDirection from hex string (Canton MCMS format)
func (e *TransferDirection) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = TransferDirection("")

// IICrossChainVerifierInterfaceID returns the interface ID for the IICrossChainVerifier interface using the package name
func IICrossChainVerifierInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.CrossChainVerifier", "ICrossChainVerifier")
}

// IICrossChainVerifierInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IICrossChainVerifierInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Interfaces.CrossChainVerifier", "ICrossChainVerifier")
}

// IIExecutorInterfaceID returns the interface ID for the IIExecutor interface using the package name
func IIExecutorInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.Executor", "IExecutor")
}

// IIExecutorInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IIExecutorInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Interfaces.Executor", "IExecutor")
}

// IITokenPoolInterfaceID returns the interface ID for the IITokenPool interface using the package name
func IITokenPoolInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.TokenPool", "ITokenPool")
}

// IITokenPoolInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IITokenPoolInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Interfaces.TokenPool", "ITokenPool")
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
