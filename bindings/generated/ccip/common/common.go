package common

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
	PackageName = "ccip-common"
	PackageID   = "1f57163d10f6a58caa575b04557d405d0e3e8d0b11ae0ac6b0231b0b2a50799b"
	SDKVersion  = "3.4.10"
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

	// CrossChainVerifierForwardToVerifier executes the CrossChainVerifier_ForwardToVerifier choice
	CrossChainVerifierForwardToVerifier(contractID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand
}

// IIExecutor is a DAML interface
type IIExecutor interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// ExecutorCalculateFee executes the Executor_CalculateFee choice
	ExecutorCalculateFee(contractID string, args ExecutorCalculateFee) *model.ExerciseCommand
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

// AddCCVFee is a Record type
type AddCCVFee struct {
	CcvInstanceId     types.TEXT    `json:"ccvInstanceId"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts AddCCVFee to a map for DAML arguments
func (t AddCCVFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddCCVFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCCVFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCCVFee to hex string (Canton MCMS format)
func (t AddCCVFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCCVFee from hex string (Canton MCMS format)
func (t *AddCCVFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCCVFeeMCMSParams is AddCCVFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddCCVFeeMCMSParams struct {
	CcvInstanceId     types.TEXT    `json:"ccvInstanceId"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// MarshalHex encodes AddCCVFeeMCMSParams to hex string for MCMS operationData.
func (t AddCCVFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCCVFeeMCMSParams from hex string.
func (t *AddCCVFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCCVVerification is a Record type
type AddCCVVerification struct {
	CcvInstanceId types.TEXT  `json:"ccvInstanceId"`
	VersionTag    types.TEXT  `json:"versionTag"`
	Caller        types.PARTY `json:"caller"`
}

// ToMap converts AddCCVVerification to a map for DAML arguments
func (t AddCCVVerification) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["versionTag"] = string(t.VersionTag)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddCCVVerification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCCVVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCCVVerification to hex string (Canton MCMS format)
func (t AddCCVVerification) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCCVVerification from hex string (Canton MCMS format)
func (t *AddCCVVerification) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCCVVerificationMCMSParams is AddCCVVerification without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddCCVVerificationMCMSParams struct {
	CcvInstanceId types.TEXT `json:"ccvInstanceId"`
	VersionTag    types.TEXT `json:"versionTag"`
}

// MarshalHex encodes AddCCVVerificationMCMSParams to hex string for MCMS operationData.
func (t AddCCVVerificationMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCCVVerificationMCMSParams from hex string.
func (t *AddCCVVerificationMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddExecutorFee is a Record type
type AddExecutorFee struct {
	ExecutorInstanceId types.TEXT    `json:"executorInstanceId"`
	ExecutorArgs       types.TEXT    `json:"executorArgs"`
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
	Caller             types.PARTY   `json:"caller"`
}

// ToMap converts AddExecutorFee to a map for DAML arguments
func (t AddExecutorFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorInstanceId"] = string(t.ExecutorInstanceId)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["feeUSDCents"] = t.FeeUSDCents

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddExecutorFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddExecutorFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddExecutorFee to hex string (Canton MCMS format)
func (t AddExecutorFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddExecutorFee from hex string (Canton MCMS format)
func (t *AddExecutorFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddExecutorFeeMCMSParams is AddExecutorFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddExecutorFeeMCMSParams struct {
	ExecutorInstanceId types.TEXT    `json:"executorInstanceId"`
	ExecutorArgs       types.TEXT    `json:"executorArgs"`
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
}

// MarshalHex encodes AddExecutorFeeMCMSParams to hex string for MCMS operationData.
func (t AddExecutorFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddExecutorFeeMCMSParams from hex string.
func (t *AddExecutorFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSend is a Record type
type AddTokenSend struct {
	PoolInstanceId   types.TEXT                               `json:"poolInstanceId"`
	PoolOwner        types.PARTY                              `json:"poolOwner"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount           types.NUMERIC                            `json:"amount"`
	DestTokenAddress types.TEXT                               `json:"destTokenAddress"`
	ExtraData        types.TEXT                               `json:"extraData"`
}

// ToMap converts AddTokenSend to a map for DAML arguments
func (t AddTokenSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = t.Amount

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["extraData"] = string(t.ExtraData)

	return m
}

func (t AddTokenSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddTokenSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddTokenSend to hex string (Canton MCMS format)
func (t AddTokenSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSend from hex string (Canton MCMS format)
func (t *AddTokenSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSendFee is a Record type
type AddTokenSendFee struct {
	PoolInstanceId    types.TEXT    `json:"poolInstanceId"`
	PoolOwner         types.PARTY   `json:"poolOwner"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts AddTokenSendFee to a map for DAML arguments
func (t AddTokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t AddTokenSendFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddTokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddTokenSendFee to hex string (Canton MCMS format)
func (t AddTokenSendFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddTokenSendFee from hex string (Canton MCMS format)
func (t *AddTokenSendFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddVerifierData is a Record type
type AddVerifierData struct {
	CcvInstanceId        types.TEXT    `json:"ccvInstanceId"`
	VersionTag           types.TEXT    `json:"versionTag"`
	VerifierBlob         types.TEXT    `json:"verifierBlob"`
	MessageSentObservers []types.PARTY `json:"messageSentObservers"`
	Caller               types.PARTY   `json:"caller"`
}

// ToMap converts AddVerifierData to a map for DAML arguments
func (t AddVerifierData) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["versionTag"] = string(t.VersionTag)

	m["verifierBlob"] = string(t.VerifierBlob)

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AddVerifierData) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddVerifierData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddVerifierData to hex string (Canton MCMS format)
func (t AddVerifierData) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddVerifierData from hex string (Canton MCMS format)
func (t *AddVerifierData) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddVerifierDataMCMSParams is AddVerifierData without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AddVerifierDataMCMSParams struct {
	CcvInstanceId        types.TEXT    `json:"ccvInstanceId"`
	VersionTag           types.TEXT    `json:"versionTag"`
	VerifierBlob         types.TEXT    `json:"verifierBlob"`
	MessageSentObservers []types.PARTY `json:"messageSentObservers"`
}

// MarshalHex encodes AddVerifierDataMCMSParams to hex string for MCMS operationData.
func (t AddVerifierDataMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddVerifierDataMCMSParams from hex string.
func (t *AddVerifierDataMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Any2CantonMessage is a Record type
type Any2CantonMessage struct {
	MessageId           types.TEXT    `json:"messageId"`
	SourceChainSelector types.NUMERIC `json:"sourceChainSelector"`
	Sender              types.TEXT    `json:"sender"`
	Payload             types.TEXT    `json:"payload"`
	DestTokenAmount     *TokenAmount  `json:"destTokenAmount" hex:"optional"`
}

// ToMap converts Any2CantonMessage to a map for DAML arguments
func (t Any2CantonMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageId"] = string(t.MessageId)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sender"] = string(t.Sender)

	m["payload"] = string(t.Payload)

	if t.DestTokenAmount != nil {
		m["destTokenAmount"] = map[string]any{
			"_type": "optional",
			"value": *t.DestTokenAmount,
		}
	} else {
		m["destTokenAmount"] = map[string]any{
			"_type": "optional",
		}
	}

	return m
}

func (t Any2CantonMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Any2CantonMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Any2CantonMessage to hex string (Canton MCMS format)
func (t Any2CantonMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Any2CantonMessage from hex string (Canton MCMS format)
func (t *Any2CantonMessage) UnmarshalHex(data string) error {
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

// ApplyDestChainConfigUpdates is a Record type
type ApplyDestChainConfigUpdates struct {
	DestChainConfigUpdates []DestChainConfigArgs `json:"destChainConfigUpdates"`
}

// ToMap converts ApplyDestChainConfigUpdates to a map for DAML arguments
func (t ApplyDestChainConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigUpdates"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigUpdates))
		for _, e := range t.DestChainConfigUpdates {
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

func (t ApplyDestChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainConfigUpdates to hex string (Canton MCMS format)
func (t ApplyDestChainConfigUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainConfigUpdates from hex string (Canton MCMS format)
func (t *ApplyDestChainConfigUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplySourceChainConfigUpdates is a Record type
type ApplySourceChainConfigUpdates struct {
	SourceChainConfigUpdates []SourceChainConfigArgs `json:"sourceChainConfigUpdates"`
}

// ToMap converts ApplySourceChainConfigUpdates to a map for DAML arguments
func (t ApplySourceChainConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainConfigUpdates"] = func() []any {
		res := make([]any, 0, len(t.SourceChainConfigUpdates))
		for _, e := range t.SourceChainConfigUpdates {
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

func (t ApplySourceChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplySourceChainConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplySourceChainConfigUpdates to hex string (Canton MCMS format)
func (t ApplySourceChainConfigUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplySourceChainConfigUpdates from hex string (Canton MCMS format)
func (t *ApplySourceChainConfigUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BuildMessage is a Record type
type BuildMessage struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts BuildMessage to a map for DAML arguments
func (t BuildMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t BuildMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BuildMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BuildMessage to hex string (Canton MCMS format)
func (t BuildMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BuildMessage from hex string (Canton MCMS format)
func (t *BuildMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BuildMessageMCMSParams is BuildMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type BuildMessageMCMSParams struct {
}

// MarshalHex encodes BuildMessageMCMSParams to hex string for MCMS operationData.
func (t BuildMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BuildMessageMCMSParams from hex string.
func (t *BuildMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCIPContext is a Record type
type CCIPContext struct {
	Values types.TEXTMAP `json:"values"`
}

// ToMap converts CCIPContext to a map for DAML arguments
func (t CCIPContext) ToMap() map[string]any {
	m := make(map[string]any)

	m["values"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Values).(mapper); ok {
			return m.toMap()
		}
		return t.Values
	}()

	return m
}

func (t CCIPContext) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPContext) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPContext to hex string (Canton MCMS format)
func (t CCIPContext) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPContext from hex string (Canton MCMS format)
func (t *CCIPContext) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCIPMessageSent is a Template type
type CCIPMessageSent struct {
	CcipOwner types.PARTY          `json:"ccipOwner"`
	CcvOwners []types.PARTY        `json:"ccvOwners"`
	Sender    types.PARTY          `json:"sender"`
	Observers []types.PARTY        `json:"observers"`
	Event     CCIPMessageSentEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPMessageSent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "CCIPMessageSent")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPMessageSent) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Events", "CCIPMessageSent")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPMessageSent) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPMessageSent) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPMessageSent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPMessageSent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPMessageSent to hex string (Canton MCMS format)
func (t CCIPMessageSent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPMessageSent from hex string (Canton MCMS format)
func (t *CCIPMessageSent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPMessageSent

// Archive exercises the Archive choice on this CCIPMessageSent contract
// This method uses the package name in the template ID
func (t CCIPMessageSent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPMessageSent) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Events", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CCIPMessageSentEvent is a Record type
type CCIPMessageSentEvent struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	SequenceNumber    types.NUMERIC `json:"sequenceNumber"`
	MessageId         types.TEXT    `json:"messageId"`
	EncodedMessage    types.TEXT    `json:"encodedMessage"`
	VerifierBlobs     []types.TEXT  `json:"verifierBlobs"`
	Receipts          []Receipt     `json:"receipts"`
}

// ToMap converts CCIPMessageSentEvent to a map for DAML arguments
func (t CCIPMessageSentEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["messageId"] = string(t.MessageId)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["verifierBlobs"] = func() []any {
		res := make([]any, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["receipts"] = func() []any {
		res := make([]any, 0, len(t.Receipts))
		for _, e := range t.Receipts {
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

func (t CCIPMessageSentEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPMessageSentEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPMessageSentEvent to hex string (Canton MCMS format)
func (t CCIPMessageSentEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPMessageSentEvent from hex string (Canton MCMS format)
func (t *CCIPMessageSentEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCVFee is a Record type
type CCVFee struct {
	CcvInstanceId     types.TEXT    `json:"ccvInstanceId"`
	CcvOwner          types.PARTY   `json:"ccvOwner"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts CCVFee to a map for DAML arguments
func (t CCVFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t CCVFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCVFee to hex string (Canton MCMS format)
func (t CCVFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVFee from hex string (Canton MCMS format)
func (t *CCVFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCVVerification is a Record type
type CCVVerification struct {
	CcvInstanceId types.TEXT  `json:"ccvInstanceId"`
	CcvOwner      types.PARTY `json:"ccvOwner"`
	VersionTag    types.TEXT  `json:"versionTag"`
}

// ToMap converts CCVVerification to a map for DAML arguments
func (t CCVVerification) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["versionTag"] = string(t.VersionTag)

	return m
}

func (t CCVVerification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCVVerification to hex string (Canton MCMS format)
func (t CCVVerification) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVVerification from hex string (Canton MCMS format)
func (t *CCVVerification) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelExecute is a Record type
type CancelExecute struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts CancelExecute to a map for DAML arguments
func (t CancelExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CancelExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelExecute to hex string (Canton MCMS format)
func (t CancelExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelExecute from hex string (Canton MCMS format)
func (t *CancelExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelExecuteMCMSParams is CancelExecute without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CancelExecuteMCMSParams struct {
}

// MarshalHex encodes CancelExecuteMCMSParams to hex string for MCMS operationData.
func (t CancelExecuteMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelExecuteMCMSParams from hex string.
func (t *CancelExecuteMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Canton2AnyMessage is a Record type
type Canton2AnyMessage struct {
	Receiver    types.TEXT                               `json:"receiver"`
	Payload     types.TEXT                               `json:"payload"`
	TokenAmount *TokenAmount                             `json:"tokenAmount" hex:"optional"`
	FeeToken    splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	ExtraArgs   ExtraArgs                                `json:"extraArgs"`
}

// ToMap converts Canton2AnyMessage to a map for DAML arguments
func (t Canton2AnyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	if t.TokenAmount != nil {
		m["tokenAmount"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenAmount,
		}
	} else {
		m["tokenAmount"] = map[string]any{
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

// ConsumeCapacity is a Record type
type ConsumeCapacity struct {
	Requested types.NUMERIC `json:"requested"`
}

// ToMap converts ConsumeCapacity to a map for DAML arguments
func (t ConsumeCapacity) ToMap() map[string]any {
	m := make(map[string]any)

	m["requested"] = t.Requested

	return m
}

func (t ConsumeCapacity) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeCapacity) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeCapacity to hex string (Canton MCMS format)
func (t ConsumeCapacity) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeCapacity from hex string (Canton MCMS format)
func (t *ConsumeCapacity) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeCapacityResult is a Record type
type ConsumeCapacityResult struct {
	RateLimiterCid         types.CONTRACT_ID `json:"rateLimiterCid"`
	AvailableBeforeConsume types.NUMERIC     `json:"availableBeforeConsume"`
	Consumed               types.NUMERIC     `json:"consumed"`
}

// ToMap converts ConsumeCapacityResult to a map for DAML arguments
func (t ConsumeCapacityResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimiterCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RateLimiterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RateLimiterCid
	}()

	m["availableBeforeConsume"] = t.AvailableBeforeConsume

	m["consumed"] = t.Consumed

	return m
}

func (t ConsumeCapacityResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeCapacityResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeCapacityResult to hex string (Canton MCMS format)
func (t ConsumeCapacityResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeCapacityResult from hex string (Canton MCMS format)
func (t *ConsumeCapacityResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	CcipOwner        types.PARTY  `json:"ccipOwner"`
	StorageLocations []types.TEXT `json:"storageLocations"`
}

// ToMap converts CrossChainVerifierView to a map for DAML arguments
func (t CrossChainVerifierView) ToMap() map[string]any {
	m := make(map[string]any)

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
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts CrossChainVerifierCalculateFee to a map for DAML arguments
func (t CrossChainVerifierCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

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
	RmnRemoteCid      types.CONTRACT_ID `json:"rmnRemoteCid"`
	ExtraContext      CCIPContext       `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT        `json:"verifierArgs"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts CrossChainVerifierForwardToVerifier to a map for DAML arguments
func (t CrossChainVerifierForwardToVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

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

// CrossChainVerifierVerifyMessage is a Record type
type CrossChainVerifierVerifyMessage struct {
	RmnRemoteCid        types.CONTRACT_ID `json:"rmnRemoteCid"`
	ExtraContext        CCIPContext       `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	VerifierResults     types.TEXT        `json:"verifierResults"`
	Caller              types.PARTY       `json:"caller"`
}

// ToMap converts CrossChainVerifierVerifyMessage to a map for DAML arguments
func (t CrossChainVerifierVerifyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
	}()

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

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

// DestChainConfig is a Record type
type DestChainConfig struct {
	IsEnabled                 types.BOOL                `json:"isEnabled"`
	AddressBytesLength        types.INT64               `json:"addressBytesLength"`
	TokenReceiverAllowed      types.BOOL                `json:"tokenReceiverAllowed"`
	BaseExecutionGasCost      types.INT64               `json:"baseExecutionGasCost"`
	OffRampAddress            types.TEXT                `json:"offRampAddress"`
	DefaultExecutor           mcms.RawInstanceAddress   `json:"defaultExecutor"`
	LaneMandatedCCVs          []mcms.RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []mcms.RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC             `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC             `json:"tokenNetworkFeeUSDCents"`
}

// ToMap converts DestChainConfig to a map for DAML arguments
func (t DestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["addressBytesLength"] = int64(t.AddressBytesLength)

	m["tokenReceiverAllowed"] = bool(t.TokenReceiverAllowed)

	m["baseExecutionGasCost"] = int64(t.BaseExecutionGasCost)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["defaultExecutor"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DefaultExecutor).(mapper); ok {
			return m.toMap()
		}
		return t.DefaultExecutor
	}()

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["defaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["messageNetworkFeeUSDCents"] = t.MessageNetworkFeeUSDCents

	m["tokenNetworkFeeUSDCents"] = t.TokenNetworkFeeUSDCents

	return m
}

func (t DestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DestChainConfig to hex string (Canton MCMS format)
func (t DestChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DestChainConfig from hex string (Canton MCMS format)
func (t *DestChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DestChainConfigArgs is a Record type
type DestChainConfigArgs struct {
	DestChainSelector         types.NUMERIC             `json:"destChainSelector"`
	IsEnabled                 types.BOOL                `json:"isEnabled"`
	AddressBytesLength        types.INT64               `json:"addressBytesLength"`
	TokenReceiverAllowed      types.BOOL                `json:"tokenReceiverAllowed"`
	BaseExecutionGasCost      types.INT64               `json:"baseExecutionGasCost"`
	OffRampAddress            types.TEXT                `json:"offRampAddress"`
	DefaultExecutor           mcms.RawInstanceAddress   `json:"defaultExecutor"`
	LaneMandatedCCVs          []mcms.RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []mcms.RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC             `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC             `json:"tokenNetworkFeeUSDCents"`
}

// ToMap converts DestChainConfigArgs to a map for DAML arguments
func (t DestChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["isEnabled"] = bool(t.IsEnabled)

	m["addressBytesLength"] = int64(t.AddressBytesLength)

	m["tokenReceiverAllowed"] = bool(t.TokenReceiverAllowed)

	m["baseExecutionGasCost"] = int64(t.BaseExecutionGasCost)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["defaultExecutor"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DefaultExecutor).(mapper); ok {
			return m.toMap()
		}
		return t.DefaultExecutor
	}()

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["defaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["messageNetworkFeeUSDCents"] = t.MessageNetworkFeeUSDCents

	m["tokenNetworkFeeUSDCents"] = t.TokenNetworkFeeUSDCents

	return m
}

func (t DestChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DestChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DestChainConfigArgs to hex string (Canton MCMS format)
func (t DestChainConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DestChainConfigArgs from hex string (Canton MCMS format)
func (t *DestChainConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutingMessageDeps is a Record type
type ExecutingMessageDeps struct {
	OffRamp            mcms.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts ExecutingMessageDeps to a map for DAML arguments
func (t ExecutingMessageDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["offRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OffRamp
	}()

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["tokenAdminRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistry
	}()

	return m
}

func (t ExecutingMessageDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageDeps to hex string (Canton MCMS format)
func (t ExecutingMessageDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageDeps from hex string (Canton MCMS format)
func (t *ExecutingMessageDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutingMessageState is an enum type
type ExecutingMessageState string

const (
	ExecutingMessageStateExecutingMessageState_RequirePoolCCVs ExecutingMessageState = "ExecutingMessageState_RequirePoolCCVs"

	ExecutingMessageStateExecutingMessageState_Prepared ExecutingMessageState = "ExecutingMessageState_Prepared"
)

func (e ExecutingMessageState) GetEnumConstructor() string { return string(e) }

func (e ExecutingMessageState) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e ExecutingMessageState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageState")
}

func (e ExecutingMessageState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *ExecutingMessageState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes ExecutingMessageState to hex string (Canton MCMS format)
func (e ExecutingMessageState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes ExecutingMessageState from hex string (Canton MCMS format)
func (e *ExecutingMessageState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = ExecutingMessageState("")

// ExecutingMessageV1 is a Template type
type ExecutingMessageV1 struct {
	CcipOwner                     types.PARTY                `json:"ccipOwner"`
	Message                       MessageV1                  `json:"message"`
	MessageId                     types.TEXT                 `json:"messageId"`
	Receiver                      types.PARTY                `json:"receiver"`
	TokenReceiver                 *types.PARTY               `json:"tokenReceiver" hex:"optional"`
	Executor                      types.PARTY                `json:"executor"`
	ObservingParties              []types.PARTY              `json:"observingParties"`
	CcvVerifications              []CCVVerification          `json:"ccvVerifications"`
	CcvOwners                     []types.PARTY              `json:"ccvOwners"`
	RequiredCCVs                  []mcms.RawInstanceAddress  `json:"requiredCCVs"`
	OptionalCCVs                  []mcms.RawInstanceAddress  `json:"optionalCCVs"`
	OptionalCCVThreshold          types.INT64                `json:"optionalCCVThreshold"`
	ReceiverMinBlockConfirmations types.INT64                `json:"receiverMinBlockConfirmations"`
	SourceDefaultCCVs             []mcms.RawInstanceAddress  `json:"sourceDefaultCCVs"`
	InboundPoolCCVs               *[]mcms.RawInstanceAddress `json:"inboundPoolCCVs" hex:"optional"`
	Deps                          ExecutingMessageDeps       `json:"deps"`
	State                         ExecutingMessageState      `json:"state"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutingMessageV1) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutingMessageV1) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutingMessageV1) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.TokenReceiver != nil {
		args["tokenReceiver"] = map[string]any{
			"_type": "optional",
			"value": (*t.TokenReceiver).ToMap(),
		}
	} else {
		args["tokenReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executor"] = t.Executor.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observingParties"] = func() []any {
		res := make([]any, 0, len(t.ObservingParties))
		for _, e := range t.ObservingParties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvVerifications"] = func() []any {
		res := make([]any, 0, len(t.CcvVerifications))
		for _, e := range t.CcvVerifications {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVThreshold"] = int64(t.OptionalCCVThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverMinBlockConfirmations"] = int64(t.ReceiverMinBlockConfirmations)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.SourceDefaultCCVs))
		for _, e := range t.SourceDefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.InboundPoolCCVs != nil {
		args["inboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": *t.InboundPoolCCVs,
		}
	} else {
		args["inboundPoolCCVs"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	if t.State != "" {
		args["state"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutingMessageV1) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.TokenReceiver != nil {
		args["tokenReceiver"] = map[string]any{
			"_type": "optional",
			"value": (*t.TokenReceiver).ToMap(),
		}
	} else {
		args["tokenReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executor"] = t.Executor.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observingParties"] = func() []any {
		res := make([]any, 0, len(t.ObservingParties))
		for _, e := range t.ObservingParties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvVerifications"] = func() []any {
		res := make([]any, 0, len(t.CcvVerifications))
		for _, e := range t.CcvVerifications {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVThreshold"] = int64(t.OptionalCCVThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverMinBlockConfirmations"] = int64(t.ReceiverMinBlockConfirmations)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.SourceDefaultCCVs))
		for _, e := range t.SourceDefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.InboundPoolCCVs != nil {
		args["inboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": *t.InboundPoolCCVs,
		}
	} else {
		args["inboundPoolCCVs"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	if t.State != "" {
		args["state"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutingMessageV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageV1 to hex string (Canton MCMS format)
func (t ExecutingMessageV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageV1 from hex string (Canton MCMS format)
func (t *ExecutingMessageV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutingMessageV1

// CancelExecute exercises the CancelExecute choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) CancelExecute(contractID string, args CancelExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "CancelExecute",
		Arguments:  argsToMap(args),
	}
}

// CancelExecuteWithPackageID exercises the CancelExecute choice using the provided package ID instead of package name
func (t ExecutingMessageV1) CancelExecuteWithPackageID(contractID string, packageID string, args CancelExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "CancelExecute",
		Arguments:  argsToMap(args),
	}
}

// AddCCVVerification exercises the AddCCVVerification choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) AddCCVVerification(contractID string, args AddCCVVerification) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "AddCCVVerification",
		Arguments:  argsToMap(args),
	}
}

// AddCCVVerificationWithPackageID exercises the AddCCVVerification choice using the provided package ID instead of package name
func (t ExecutingMessageV1) AddCCVVerificationWithPackageID(contractID string, packageID string, args AddCCVVerification) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "AddCCVVerification",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVs exercises the SetInboundPoolCCVs choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) SetInboundPoolCCVs(contractID string, args SetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVsWithPackageID exercises the SetInboundPoolCCVs choice using the provided package ID instead of package name
func (t ExecutingMessageV1) SetInboundPoolCCVsWithPackageID(contractID string, packageID string, args SetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecute exercises the FinalizeExecute choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) FinalizeExecute(contractID string, args FinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecuteWithPackageID exercises the FinalizeExecute choice using the provided package ID instead of package name
func (t ExecutingMessageV1) FinalizeExecuteWithPackageID(contractID string, packageID string, args FinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutingMessageV1) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutionMode is an enum type
type ExecutionMode string

const (
	ExecutionModeExecutionMode_Executor ExecutionMode = "ExecutionMode_Executor"

	ExecutionModeExecutionMode_NoExecutor ExecutionMode = "ExecutionMode_NoExecutor"
)

func (e ExecutionMode) GetEnumConstructor() string { return string(e) }

func (e ExecutionMode) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "ExecutionMode")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e ExecutionMode) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "ExecutionMode")
}

func (e ExecutionMode) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *ExecutionMode) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes ExecutionMode to hex string (Canton MCMS format)
func (e ExecutionMode) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes ExecutionMode from hex string (Canton MCMS format)
func (e *ExecutionMode) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = ExecutionMode("")

// ExecutionStateChanged is a Template type
type ExecutionStateChanged struct {
	CcipOwner types.PARTY                `json:"ccipOwner"`
	CcvOwners []types.PARTY              `json:"ccvOwners"`
	Receiver  types.PARTY                `json:"receiver"`
	Event     ExecutionStateChangedEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutionStateChanged) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "ExecutionStateChanged")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutionStateChanged) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Events", "ExecutionStateChanged")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutionStateChanged) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutionStateChanged) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutionStateChanged) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutionStateChanged) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutionStateChanged to hex string (Canton MCMS format)
func (t ExecutionStateChanged) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutionStateChanged from hex string (Canton MCMS format)
func (t *ExecutionStateChanged) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutionStateChanged

// Archive exercises the Archive choice on this ExecutionStateChanged contract
// This method uses the package name in the template ID
func (t ExecutionStateChanged) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutionStateChanged) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Events", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutionStateChangedEvent is a Record type
type ExecutionStateChangedEvent struct {
	SourceChainSelector types.NUMERIC         `json:"sourceChainSelector"`
	SequenceNumber      types.NUMERIC         `json:"sequenceNumber"`
	MessageId           types.TEXT            `json:"messageId"`
	State               MessageExecutionState `json:"state"`
	ReturnData          types.TEXT            `json:"returnData"`
}

// ToMap converts ExecutionStateChangedEvent to a map for DAML arguments
func (t ExecutionStateChangedEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	m["messageId"] = string(t.MessageId)

	m["state"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.State).(mapper); ok {
			return m.toMap()
		}
		return t.State
	}()

	m["returnData"] = string(t.ReturnData)

	return m
}

func (t ExecutionStateChangedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutionStateChangedEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutionStateChangedEvent to hex string (Canton MCMS format)
func (t ExecutionStateChangedEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutionStateChangedEvent from hex string (Canton MCMS format)
func (t *ExecutionStateChangedEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorFee is a Record type
type ExecutorFee struct {
	ExecutorInstanceId types.TEXT    `json:"executorInstanceId"`
	ExecutorOwner      types.PARTY   `json:"executorOwner"`
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
}

// ToMap converts ExecutorFee to a map for DAML arguments
func (t ExecutorFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorInstanceId"] = string(t.ExecutorInstanceId)

	m["executorOwner"] = t.ExecutorOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	return m
}

func (t ExecutorFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorFee to hex string (Canton MCMS format)
func (t ExecutorFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorFee from hex string (Canton MCMS format)
func (t *ExecutorFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorView is a Record type
type ExecutorView struct {
	Owner types.PARTY `json:"owner"`
}

// ToMap converts ExecutorView to a map for DAML arguments
func (t ExecutorView) ToMap() map[string]any {
	m := make(map[string]any)

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
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	ExecutorArgs      types.TEXT        `json:"executorArgs"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts ExecutorCalculateFee to a map for DAML arguments
func (t ExecutorCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["executorArgs"] = string(t.ExecutorArgs)

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

// FeeTokenAmount is a Record type
type FeeTokenAmount struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts FeeTokenAmount to a map for DAML arguments
func (t FeeTokenAmount) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeTokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeTokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeTokenAmount to hex string (Canton MCMS format)
func (t FeeTokenAmount) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeTokenAmount from hex string (Canton MCMS format)
func (t *FeeTokenAmount) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeTokenAmountMCMSParams is FeeTokenAmount without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type FeeTokenAmountMCMSParams struct {
}

// MarshalHex encodes FeeTokenAmountMCMSParams to hex string for MCMS operationData.
func (t FeeTokenAmountMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeTokenAmountMCMSParams from hex string.
func (t *FeeTokenAmountMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeExecute is a Record type
type FinalizeExecute struct {
	TokenAdminRegistryInstanceId types.TEXT                                `json:"tokenAdminRegistryInstanceId"`
	MaybePoolOwner               *types.PARTY                              `json:"maybePoolOwner" hex:"optional"`
	MaybeTicketReceiver          *types.PARTY                              `json:"maybeTicketReceiver" hex:"optional"`
	MaybeTokenReceiver           *types.PARTY                              `json:"maybeTokenReceiver" hex:"optional"`
	MaybeInstrumentId            *splice_api_token_holding_v1.InstrumentId `json:"maybeInstrumentId" hex:"optional"`
	MaybeAmount                  *types.NUMERIC                            `json:"maybeAmount" hex:"optional"`
	ReturnData                   types.TEXT                                `json:"returnData"`
}

// ToMap converts FinalizeExecute to a map for DAML arguments
func (t FinalizeExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	if t.MaybePoolOwner != nil {
		m["maybePoolOwner"] = map[string]any{
			"_type": "optional",
			"value": (*t.MaybePoolOwner).ToMap(),
		}
	} else {
		m["maybePoolOwner"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.MaybeTicketReceiver != nil {
		m["maybeTicketReceiver"] = map[string]any{
			"_type": "optional",
			"value": (*t.MaybeTicketReceiver).ToMap(),
		}
	} else {
		m["maybeTicketReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.MaybeTokenReceiver != nil {
		m["maybeTokenReceiver"] = map[string]any{
			"_type": "optional",
			"value": (*t.MaybeTokenReceiver).ToMap(),
		}
	} else {
		m["maybeTokenReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.MaybeInstrumentId != nil {
		m["maybeInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": *t.MaybeInstrumentId,
		}
	} else {
		m["maybeInstrumentId"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.MaybeAmount != nil {
		m["maybeAmount"] = map[string]any{
			"_type": "optional",
			"value": *t.MaybeAmount,
		}
	} else {
		m["maybeAmount"] = map[string]any{
			"_type": "optional",
		}
	}

	m["returnData"] = string(t.ReturnData)

	return m
}

func (t FinalizeExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeExecute to hex string (Canton MCMS format)
func (t FinalizeExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeExecute from hex string (Canton MCMS format)
func (t *FinalizeExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeExecuteResult is a Record type
type FinalizeExecuteResult struct {
	TokenReceiveTicket    *types.CONTRACT_ID `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID  `json:"executionStateChanged"`
}

// ToMap converts FinalizeExecuteResult to a map for DAML arguments
func (t FinalizeExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenReceiveTicket,
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
		}
	}

	m["executionStateChanged"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutionStateChanged).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutionStateChanged
	}()

	return m
}

func (t FinalizeExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeExecuteResult to hex string (Canton MCMS format)
func (t FinalizeExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeExecuteResult from hex string (Canton MCMS format)
func (t *FinalizeExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeFee is a Record type
type FinalizeFee struct {
	FeeTokenPrice             types.NUMERIC `json:"feeTokenPrice"`
	PremiumMultiplier         types.NUMERIC `json:"premiumMultiplier"`
	TotalExecutionGasLimit    types.INT64   `json:"totalExecutionGasLimit"`
	ExecutorDestGasLimit      types.INT64   `json:"executorDestGasLimit"`
	ExecutorDestBytesOverhead types.INT64   `json:"executorDestBytesOverhead"`
	ExecutionCostUSDCents     types.NUMERIC `json:"executionCostUSDCents"`
}

// ToMap converts FinalizeFee to a map for DAML arguments
func (t FinalizeFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenPrice"] = t.FeeTokenPrice

	m["premiumMultiplier"] = t.PremiumMultiplier

	m["totalExecutionGasLimit"] = int64(t.TotalExecutionGasLimit)

	m["executorDestGasLimit"] = int64(t.ExecutorDestGasLimit)

	m["executorDestBytesOverhead"] = int64(t.ExecutorDestBytesOverhead)

	m["executionCostUSDCents"] = t.ExecutionCostUSDCents

	return m
}

func (t FinalizeFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeFee to hex string (Canton MCMS format)
func (t FinalizeFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeFee from hex string (Canton MCMS format)
func (t *FinalizeFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeSend is a Record type
type FinalizeSend struct {
	MessageSender        types.PARTY   `json:"messageSender"`
	MessageSentObservers []types.PARTY `json:"messageSentObservers"`
	VerifierBlobs        []types.TEXT  `json:"verifierBlobs"`
	Receipts             []Receipt     `json:"receipts"`
}

// ToMap converts FinalizeSend to a map for DAML arguments
func (t FinalizeSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageSender"] = t.MessageSender.ToMap()

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["verifierBlobs"] = func() []any {
		res := make([]any, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["receipts"] = func() []any {
		res := make([]any, 0, len(t.Receipts))
		for _, e := range t.Receipts {
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

func (t FinalizeSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeSend to hex string (Canton MCMS format)
func (t FinalizeSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeSend from hex string (Canton MCMS format)
func (t *FinalizeSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeSendResult is a Record type
type FinalizeSendResult struct {
	CcipMessageSent types.CONTRACT_ID `json:"ccipMessageSent"`
}

// ToMap converts FinalizeSendResult to a map for DAML arguments
func (t FinalizeSendResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipMessageSent"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcipMessageSent).(mapper); ok {
			return m.toMap()
		}
		return t.CcipMessageSent
	}()

	return m
}

func (t FinalizeSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeSendResult to hex string (Canton MCMS format)
func (t FinalizeSendResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeSendResult from hex string (Canton MCMS format)
func (t *FinalizeSendResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GenericExtraArgsV3 is a Record type
type GenericExtraArgsV3 struct {
	GasLimit           types.INT64  `json:"gasLimit"`
	BlockConfirmations types.INT64  `json:"blockConfirmations"`
	Ccvs               []types.TEXT `json:"ccvs"`
	CcvArgs            []types.TEXT `json:"ccvArgs"`
	Executor           types.TEXT   `json:"executor"`
	ExecutorArgs       types.TEXT   `json:"executorArgs"`
	TokenReceiver      types.TEXT   `json:"tokenReceiver"`
	TokenArgs          types.TEXT   `json:"tokenArgs"`
}

// ToMap converts GenericExtraArgsV3 to a map for DAML arguments
func (t GenericExtraArgsV3) ToMap() map[string]any {
	m := make(map[string]any)

	m["gasLimit"] = int64(t.GasLimit)

	m["blockConfirmations"] = int64(t.BlockConfirmations)

	m["ccvs"] = func() []any {
		res := make([]any, 0, len(t.Ccvs))
		for _, e := range t.Ccvs {
			res = append(res, string(e))
		}
		return res
	}()

	m["ccvArgs"] = func() []any {
		res := make([]any, 0, len(t.CcvArgs))
		for _, e := range t.CcvArgs {
			res = append(res, string(e))
		}
		return res
	}()

	m["executor"] = string(t.Executor)

	m["executorArgs"] = string(t.ExecutorArgs)

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

// GetDestChainConfig is a Record type
type GetDestChainConfig struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts GetDestChainConfig to a map for DAML arguments
func (t GetDestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDestChainConfig to hex string (Canton MCMS format)
func (t GetDestChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChainConfig from hex string (Canton MCMS format)
func (t *GetDestChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestChainConfigMCMSParams is GetDestChainConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDestChainConfigMCMSParams struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// MarshalHex encodes GetDestChainConfigMCMSParams to hex string for MCMS operationData.
func (t GetDestChainConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChainConfigMCMSParams from hex string.
func (t *GetDestChainConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetSourceChainConfig is a Record type
type GetSourceChainConfig struct {
	SourceChainSelector types.NUMERIC `json:"sourceChainSelector"`
	Caller              types.PARTY   `json:"caller"`
}

// ToMap converts GetSourceChainConfig to a map for DAML arguments
func (t GetSourceChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetSourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetSourceChainConfig to hex string (Canton MCMS format)
func (t GetSourceChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetSourceChainConfig from hex string (Canton MCMS format)
func (t *GetSourceChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetSourceChainConfigMCMSParams is GetSourceChainConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetSourceChainConfigMCMSParams struct {
	SourceChainSelector types.NUMERIC `json:"sourceChainSelector"`
}

// MarshalHex encodes GetSourceChainConfigMCMSParams to hex string for MCMS operationData.
func (t GetSourceChainConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetSourceChainConfigMCMSParams from hex string.
func (t *GetSourceChainConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfig is a Template type
type GlobalConfig struct {
	InstanceId         types.TEXT    `json:"instanceId"`
	CcipOwner          types.PARTY   `json:"ccipOwner"`
	ChainSelector      types.NUMERIC `json:"chainSelector"`
	DestChainConfigs   types.GENMAP  `json:"destChainConfigs"`
	SourceChainConfigs types.GENMAP  `json:"sourceChainConfigs"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t GlobalConfig) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t GlobalConfig) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t GlobalConfig) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	if t.ChainSelector != "" {
		args["chainSelector"] = t.ChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceChainConfigs"] = func() any {
		if t.SourceChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SourceChainConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t GlobalConfig) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	if t.ChainSelector != "" {
		args["chainSelector"] = t.ChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceChainConfigs"] = func() any {
		if t.SourceChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SourceChainConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t GlobalConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalConfig to hex string (Canton MCMS format)
func (t GlobalConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalConfig from hex string (Canton MCMS format)
func (t *GlobalConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for GlobalConfig

// ApplyDestChainConfigUpdates exercises the ApplyDestChainConfigUpdates choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) ApplyDestChainConfigUpdates(contractID string, args ApplyDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainConfigUpdatesWithPackageID exercises the ApplyDestChainConfigUpdates choice using the provided package ID instead of package name
func (t GlobalConfig) ApplyDestChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "ApplyDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplySourceChainConfigUpdates exercises the ApplySourceChainConfigUpdates choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) ApplySourceChainConfigUpdates(contractID string, args ApplySourceChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "ApplySourceChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplySourceChainConfigUpdatesWithPackageID exercises the ApplySourceChainConfigUpdates choice using the provided package ID instead of package name
func (t GlobalConfig) ApplySourceChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplySourceChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "ApplySourceChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this GlobalConfig contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t GlobalConfig) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t GlobalConfig) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetDestChainConfig exercises the GetDestChainConfig choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) GetDestChainConfig(contractID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfigWithPackageID exercises the GetDestChainConfig choice using the provided package ID instead of package name
func (t GlobalConfig) GetDestChainConfigWithPackageID(contractID string, packageID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetSourceChainConfig exercises the GetSourceChainConfig choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) GetSourceChainConfig(contractID string, args GetSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetSourceChainConfigWithPackageID exercises the GetSourceChainConfig choice using the provided package ID instead of package name
func (t GlobalConfig) GetSourceChainConfigWithPackageID(contractID string, packageID string, args GetSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this GlobalConfig contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t GlobalConfig) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t GlobalConfig) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for GlobalConfig

var _ mcms.IMCMSReceiver = (*GlobalConfig)(nil)

// IssuerType is an enum type
type IssuerType string

const (
	IssuerTypeIssuerType_CCV IssuerType = "IssuerType_CCV"

	IssuerTypeIssuerType_Pool IssuerType = "IssuerType_Pool"

	IssuerTypeIssuerType_Executor IssuerType = "IssuerType_Executor"

	IssuerTypeIssuerType_Network IssuerType = "IssuerType_Network"
)

func (e IssuerType) GetEnumConstructor() string { return string(e) }

func (e IssuerType) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Tickets", "IssuerType")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e IssuerType) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Tickets", "IssuerType")
}

func (e IssuerType) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *IssuerType) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes IssuerType to hex string (Canton MCMS format)
func (e IssuerType) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes IssuerType from hex string (Canton MCMS format)
func (e *IssuerType) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = IssuerType("")

// MessageExecutionState is an enum type
type MessageExecutionState string

const (
	MessageExecutionStateUNTOUCHED MessageExecutionState = "UNTOUCHED"

	MessageExecutionStateIN_PROGRESS MessageExecutionState = "IN_PROGRESS"

	MessageExecutionStateSUCCESS MessageExecutionState = "SUCCESS"

	MessageExecutionStateFAILURE MessageExecutionState = "FAILURE"
)

func (e MessageExecutionState) GetEnumConstructor() string { return string(e) }

func (e MessageExecutionState) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Internal", "MessageExecutionState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e MessageExecutionState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Internal", "MessageExecutionState")
}

func (e MessageExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *MessageExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes MessageExecutionState to hex string (Canton MCMS format)
func (e MessageExecutionState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes MessageExecutionState from hex string (Canton MCMS format)
func (e *MessageExecutionState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = MessageExecutionState("")

// MessageV1 is a Record type
type MessageV1 struct {
	SourceChainSelector types.NUMERIC    `json:"sourceChainSelector"`
	DestChainSelector   types.NUMERIC    `json:"destChainSelector"`
	SequenceNumber      types.NUMERIC    `json:"sequenceNumber"`
	ExecutionGasLimit   types.INT64      `json:"executionGasLimit"`
	CcipReceiveGasLimit types.INT64      `json:"ccipReceiveGasLimit"`
	Finality            types.INT64      `json:"finality"`
	CcvAndExecutorHash  types.TEXT       `json:"ccvAndExecutorHash"`
	OnRampAddress       types.TEXT       `json:"onRampAddress"`
	OffRampAddress      types.TEXT       `json:"offRampAddress"`
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

	m["finality"] = int64(t.Finality)

	m["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	m["onRampAddress"] = string(t.OnRampAddress)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["sender"] = string(t.Sender)

	m["receiver"] = string(t.Receiver)

	m["destBlob"] = string(t.DestBlob)

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

// RateLimitDirection is an enum type
type RateLimitDirection string

const (
	RateLimitDirectionRateLimitDirection_Outbound RateLimitDirection = "RateLimitDirection_Outbound"

	RateLimitDirectionRateLimitDirection_Inbound RateLimitDirection = "RateLimitDirection_Inbound"
)

func (e RateLimitDirection) GetEnumConstructor() string { return string(e) }

func (e RateLimitDirection) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimitDirection")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e RateLimitDirection) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimitDirection")
}

func (e RateLimitDirection) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *RateLimitDirection) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes RateLimitDirection to hex string (Canton MCMS format)
func (e RateLimitDirection) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes RateLimitDirection from hex string (Canton MCMS format)
func (e *RateLimitDirection) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = RateLimitDirection("")

// RateLimitMode is an enum type
type RateLimitMode string

const (
	RateLimitModeRateLimitMode_DefaultFinality RateLimitMode = "RateLimitMode_DefaultFinality"

	RateLimitModeRateLimitMode_CustomFinality RateLimitMode = "RateLimitMode_CustomFinality"
)

func (e RateLimitMode) GetEnumConstructor() string { return string(e) }

func (e RateLimitMode) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimitMode")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e RateLimitMode) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimitMode")
}

func (e RateLimitMode) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *RateLimitMode) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes RateLimitMode to hex string (Canton MCMS format)
func (e RateLimitMode) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes RateLimitMode from hex string (Canton MCMS format)
func (e *RateLimitMode) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = RateLimitMode("")

// RateLimiter is a Template type
type RateLimiter struct {
	InstanceId          types.TEXT         `json:"instanceId"`
	PoolInstanceId      types.TEXT         `json:"poolInstanceId"`
	PoolOwner           types.PARTY        `json:"poolOwner"`
	RemoteChainSelector types.NUMERIC      `json:"remoteChainSelector"`
	Direction           RateLimitDirection `json:"direction"`
	Mode                RateLimitMode      `json:"mode"`
	IsEnabled           types.BOOL         `json:"isEnabled"`
	Capacity            types.NUMERIC      `json:"capacity"`
	Rate                types.NUMERIC      `json:"rate"`
	Tokens              types.NUMERIC      `json:"tokens"`
	LastUpdated         types.TIMESTAMP    `json:"lastUpdated"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RateLimiter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimiter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RateLimiter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimiter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RateLimiter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolInstanceId"] = string(t.PoolInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	if t.RemoteChainSelector != "" {
		args["remoteChainSelector"] = t.RemoteChainSelector
	}

	if t.Direction != "" {
		args["direction"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.Direction).(mapper); ok {
				return m.toMap()
			}
			return t.Direction
		}()
	}

	if t.Mode != "" {
		args["mode"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.Mode).(mapper); ok {
				return m.toMap()
			}
			return t.Mode
		}()
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isEnabled"] = bool(t.IsEnabled)

	if t.Capacity != "" {
		args["capacity"] = t.Capacity
	}

	if t.Rate != "" {
		args["rate"] = t.Rate
	}

	if t.Tokens != "" {
		args["tokens"] = t.Tokens
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastUpdated"] = t.LastUpdated

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RateLimiter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolInstanceId"] = string(t.PoolInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	if t.RemoteChainSelector != "" {
		args["remoteChainSelector"] = t.RemoteChainSelector
	}

	if t.Direction != "" {
		args["direction"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.Direction).(mapper); ok {
				return m.toMap()
			}
			return t.Direction
		}()
	}

	if t.Mode != "" {
		args["mode"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.Mode).(mapper); ok {
				return m.toMap()
			}
			return t.Mode
		}()
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isEnabled"] = bool(t.IsEnabled)

	if t.Capacity != "" {
		args["capacity"] = t.Capacity
	}

	if t.Rate != "" {
		args["rate"] = t.Rate
	}

	if t.Tokens != "" {
		args["tokens"] = t.Tokens
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastUpdated"] = t.LastUpdated

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RateLimiter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RateLimiter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RateLimiter to hex string (Canton MCMS format)
func (t RateLimiter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RateLimiter from hex string (Canton MCMS format)
func (t *RateLimiter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RateLimiter

// ConsumeCapacity exercises the ConsumeCapacity choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) ConsumeCapacity(contractID string, args ConsumeCapacity) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "ConsumeCapacity",
		Arguments:  argsToMap(args),
	}
}

// ConsumeCapacityWithPackageID exercises the ConsumeCapacity choice using the provided package ID instead of package name
func (t RateLimiter) ConsumeCapacityWithPackageID(contractID string, packageID string, args ConsumeCapacity) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "ConsumeCapacity",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RateLimiter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetConfig exercises the SetConfig choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) SetConfig(contractID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// SetConfigWithPackageID exercises the SetConfig choice using the provided package ID instead of package name
func (t RateLimiter) SetConfigWithPackageID(contractID string, packageID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// Receipt is a Record type
type Receipt struct {
	IssuerType        IssuerType    `json:"issuerType"`
	IssuerAddress     types.TEXT    `json:"issuerAddress"`
	VersionTag        *types.TEXT   `json:"versionTag" hex:"optional"`
	DestGasLimit      types.INT64   `json:"destGasLimit"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
	FeeTokenAmount    types.NUMERIC `json:"feeTokenAmount"`
	ExtraArgs         types.TEXT    `json:"extraArgs"`
}

// ToMap converts Receipt to a map for DAML arguments
func (t Receipt) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuerType"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.IssuerType).(mapper); ok {
			return m.toMap()
		}
		return t.IssuerType
	}()

	m["issuerAddress"] = string(t.IssuerAddress)

	if t.VersionTag != nil {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": string(*t.VersionTag),
		}
	} else {
		m["versionTag"] = map[string]any{
			"_type": "optional",
		}
	}

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeTokenAmount"] = t.FeeTokenAmount

	m["extraArgs"] = string(t.ExtraArgs)

	return m
}

func (t Receipt) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Receipt) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Receipt to hex string (Canton MCMS format)
func (t Receipt) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Receipt from hex string (Canton MCMS format)
func (t *Receipt) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageDeps is a Record type
type SendingMessageDeps struct {
	Router             mcms.RawInstanceAddress `json:"router"`
	OnRamp             mcms.RawInstanceAddress `json:"onRamp"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          mcms.RawInstanceAddress `json:"feeQuoter"`
}

// ToMap converts SendingMessageDeps to a map for DAML arguments
func (t SendingMessageDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	m["onRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OnRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OnRamp
	}()

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["tokenAdminRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistry
	}()

	m["feeQuoter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoter).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoter
	}()

	return m
}

func (t SendingMessageDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageDeps to hex string (Canton MCMS format)
func (t SendingMessageDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageDeps from hex string (Canton MCMS format)
func (t *SendingMessageDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageState is an enum type
type SendingMessageState string

const (
	SendingMessageStateSendingMessageState_RequirePoolCCVs SendingMessageState = "SendingMessageState_RequirePoolCCVs"

	SendingMessageStateSendingMessageState_Prepared SendingMessageState = "SendingMessageState_Prepared"

	SendingMessageStateSendingMessageState_TokenLocked SendingMessageState = "SendingMessageState_TokenLocked"

	SendingMessageStateSendingMessageState_ExecutorFinalized SendingMessageState = "SendingMessageState_ExecutorFinalized"

	SendingMessageStateSendingMessageState_FeeFinalized SendingMessageState = "SendingMessageState_FeeFinalized"
)

func (e SendingMessageState) GetEnumConstructor() string { return string(e) }

func (e SendingMessageState) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e SendingMessageState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageState")
}

func (e SendingMessageState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *SendingMessageState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes SendingMessageState to hex string (Canton MCMS format)
func (e SendingMessageState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes SendingMessageState from hex string (Canton MCMS format)
func (e *SendingMessageState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = SendingMessageState("")

// SendingMessageV1 is a Template type
type SendingMessageV1 struct {
	Deps                      SendingMessageDeps                        `json:"deps"`
	CcipOwner                 types.PARTY                               `json:"ccipOwner"`
	Sender                    types.PARTY                               `json:"sender"`
	DestChainSelector         types.NUMERIC                             `json:"destChainSelector"`
	SequenceNumber            types.NUMERIC                             `json:"sequenceNumber"`
	RequiredCCVs              []types.TEXT                              `json:"requiredCCVs"`
	ExecutorAddress           types.TEXT                                `json:"executorAddress"`
	ExecutionMode             *ExecutionMode                            `json:"executionMode" hex:"optional"`
	SourceChainSelector       types.NUMERIC                             `json:"sourceChainSelector"`
	SenderAddress             types.TEXT                                `json:"senderAddress"`
	Receiver                  types.TEXT                                `json:"receiver"`
	Payload                   types.TEXT                                `json:"payload"`
	ExecutionGasLimit         types.INT64                               `json:"executionGasLimit"`
	CcipReceiveGasLimit       types.INT64                               `json:"ccipReceiveGasLimit"`
	CcvAndExecutorHash        types.TEXT                                `json:"ccvAndExecutorHash"`
	OnRampAddress             types.TEXT                                `json:"onRampAddress"`
	OffRampAddress            types.TEXT                                `json:"offRampAddress"`
	TokenReceiver             types.TEXT                                `json:"tokenReceiver"`
	TokenArgs                 types.TEXT                                `json:"tokenArgs"`
	FeeToken                  splice_api_token_holding_v1.InstrumentId  `json:"feeToken"`
	NetworkFeeUSDCents        types.NUMERIC                             `json:"networkFeeUSDCents"`
	ExpectedTokenInstrumentId *splice_api_token_holding_v1.InstrumentId `json:"expectedTokenInstrumentId" hex:"optional"`
	OutboundPoolCCVs          *[]types.TEXT                             `json:"outboundPoolCCVs" hex:"optional"`
	ExecutorArgs              types.TEXT                                `json:"executorArgs"`
	ExecutorFee               *ExecutorFee                              `json:"executorFee" hex:"optional"`
	ExecutorDestGasLimit      types.INT64                               `json:"executorDestGasLimit"`
	ExecutorDestBytesOverhead types.INT64                               `json:"executorDestBytesOverhead"`
	ExecutorFeeTokenAmount    types.NUMERIC                             `json:"executorFeeTokenAmount"`
	ObservingParties          []types.PARTY                             `json:"observingParties"`
	CcvFees                   []CCVFee                                  `json:"ccvFees"`
	TokenSendFee              *TokenSendFee                             `json:"tokenSendFee" hex:"optional"`
	CcvFeeTokenAmounts        []types.NUMERIC                           `json:"ccvFeeTokenAmounts"`
	TokenSendFeeTokenAmount   types.NUMERIC                             `json:"tokenSendFeeTokenAmount"`
	NetworkFeeTokenAmount     types.NUMERIC                             `json:"networkFeeTokenAmount"`
	TokenSendData             *TokenSendData                            `json:"tokenSendData" hex:"optional"`
	VerifierData              []VerifierData                            `json:"verifierData"`
	CcvOwners                 []types.PARTY                             `json:"ccvOwners"`
	Message                   *MessageV1                                `json:"message" hex:"optional"`
	EncodedMessage            types.TEXT                                `json:"encodedMessage"`
	MessageId                 types.TEXT                                `json:"messageId"`
	State                     SendingMessageState                       `json:"state"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t SendingMessageV1) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t SendingMessageV1) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t SendingMessageV1) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	if t.DestChainSelector != "" {
		args["destChainSelector"] = t.DestChainSelector
	}

	if t.SequenceNumber != "" {
		args["sequenceNumber"] = t.SequenceNumber
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorAddress"] = string(t.ExecutorAddress)

	if t.ExecutionMode != nil {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": *t.ExecutionMode,
		}
	} else {
		args["executionMode"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderAddress"] = string(t.SenderAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = string(t.Receiver)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["payload"] = string(t.Payload)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionGasLimit"] = int64(t.ExecutionGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampAddress"] = string(t.OnRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["offRampAddress"] = string(t.OffRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = string(t.TokenReceiver)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenArgs"] = string(t.TokenArgs)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeToken"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	if t.NetworkFeeUSDCents != "" {
		args["networkFeeUSDCents"] = t.NetworkFeeUSDCents
	}

	if t.ExpectedTokenInstrumentId != nil {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": *t.ExpectedTokenInstrumentId,
		}
	} else {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.OutboundPoolCCVs != nil {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": *t.OutboundPoolCCVs,
		}
	} else {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorArgs"] = string(t.ExecutorArgs)

	if t.ExecutorFee != nil {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": *t.ExecutorFee,
		}
	} else {
		args["executorFee"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorDestGasLimit"] = int64(t.ExecutorDestGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorDestBytesOverhead"] = int64(t.ExecutorDestBytesOverhead)

	if t.ExecutorFeeTokenAmount != "" {
		args["executorFeeTokenAmount"] = t.ExecutorFeeTokenAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observingParties"] = func() []any {
		res := make([]any, 0, len(t.ObservingParties))
		for _, e := range t.ObservingParties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvFees"] = func() []any {
		res := make([]any, 0, len(t.CcvFees))
		for _, e := range t.CcvFees {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.TokenSendFee != nil {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenSendFee,
		}
	} else {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvFeeTokenAmounts"] = func() []any {
		res := make([]any, 0, len(t.CcvFeeTokenAmounts))
		for _, e := range t.CcvFeeTokenAmounts {
			res = append(res, e)
		}
		return res
	}()

	if t.TokenSendFeeTokenAmount != "" {
		args["tokenSendFeeTokenAmount"] = t.TokenSendFeeTokenAmount
	}

	if t.NetworkFeeTokenAmount != "" {
		args["networkFeeTokenAmount"] = t.NetworkFeeTokenAmount
	}

	if t.TokenSendData != nil {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenSendData,
		}
	} else {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifierData"] = func() []any {
		res := make([]any, 0, len(t.VerifierData))
		for _, e := range t.VerifierData {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	if t.Message != nil {
		args["message"] = map[string]any{
			"_type": "optional",
			"value": *t.Message,
		}
	} else {
		args["message"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["encodedMessage"] = string(t.EncodedMessage)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.State != "" {
		args["state"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t SendingMessageV1) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	if t.DestChainSelector != "" {
		args["destChainSelector"] = t.DestChainSelector
	}

	if t.SequenceNumber != "" {
		args["sequenceNumber"] = t.SequenceNumber
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorAddress"] = string(t.ExecutorAddress)

	if t.ExecutionMode != nil {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": *t.ExecutionMode,
		}
	} else {
		args["executionMode"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderAddress"] = string(t.SenderAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = string(t.Receiver)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["payload"] = string(t.Payload)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionGasLimit"] = int64(t.ExecutionGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampAddress"] = string(t.OnRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["offRampAddress"] = string(t.OffRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = string(t.TokenReceiver)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenArgs"] = string(t.TokenArgs)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeToken"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	if t.NetworkFeeUSDCents != "" {
		args["networkFeeUSDCents"] = t.NetworkFeeUSDCents
	}

	if t.ExpectedTokenInstrumentId != nil {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": *t.ExpectedTokenInstrumentId,
		}
	} else {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.OutboundPoolCCVs != nil {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": *t.OutboundPoolCCVs,
		}
	} else {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorArgs"] = string(t.ExecutorArgs)

	if t.ExecutorFee != nil {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": *t.ExecutorFee,
		}
	} else {
		args["executorFee"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorDestGasLimit"] = int64(t.ExecutorDestGasLimit)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorDestBytesOverhead"] = int64(t.ExecutorDestBytesOverhead)

	if t.ExecutorFeeTokenAmount != "" {
		args["executorFeeTokenAmount"] = t.ExecutorFeeTokenAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observingParties"] = func() []any {
		res := make([]any, 0, len(t.ObservingParties))
		for _, e := range t.ObservingParties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvFees"] = func() []any {
		res := make([]any, 0, len(t.CcvFees))
		for _, e := range t.CcvFees {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.TokenSendFee != nil {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenSendFee,
		}
	} else {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvFeeTokenAmounts"] = func() []any {
		res := make([]any, 0, len(t.CcvFeeTokenAmounts))
		for _, e := range t.CcvFeeTokenAmounts {
			res = append(res, e)
		}
		return res
	}()

	if t.TokenSendFeeTokenAmount != "" {
		args["tokenSendFeeTokenAmount"] = t.TokenSendFeeTokenAmount
	}

	if t.NetworkFeeTokenAmount != "" {
		args["networkFeeTokenAmount"] = t.NetworkFeeTokenAmount
	}

	if t.TokenSendData != nil {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenSendData,
		}
	} else {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifierData"] = func() []any {
		res := make([]any, 0, len(t.VerifierData))
		for _, e := range t.VerifierData {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	if t.Message != nil {
		args["message"] = map[string]any{
			"_type": "optional",
			"value": *t.Message,
		}
	} else {
		args["message"] = map[string]any{
			"_type": "optional",
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["encodedMessage"] = string(t.EncodedMessage)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.State != "" {
		args["state"] = func() any {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t SendingMessageV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageV1 to hex string (Canton MCMS format)
func (t SendingMessageV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageV1 from hex string (Canton MCMS format)
func (t *SendingMessageV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for SendingMessageV1

// FinalizeFee exercises the FinalizeFee choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) FinalizeFee(contractID string, args FinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFeeWithPackageID exercises the FinalizeFee choice using the provided package ID instead of package name
func (t SendingMessageV1) FinalizeFeeWithPackageID(contractID string, packageID string, args FinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSend exercises the AddTokenSend choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) AddTokenSend(contractID string, args AddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendWithPackageID exercises the AddTokenSend choice using the provided package ID instead of package name
func (t SendingMessageV1) AddTokenSendWithPackageID(contractID string, packageID string, args AddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// BuildMessage exercises the BuildMessage choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) BuildMessage(contractID string, args BuildMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "BuildMessage",
		Arguments:  argsToMap(args),
	}
}

// BuildMessageWithPackageID exercises the BuildMessage choice using the provided package ID instead of package name
func (t SendingMessageV1) BuildMessageWithPackageID(contractID string, packageID string, args BuildMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "BuildMessage",
		Arguments:  argsToMap(args),
	}
}

// AddVerifierData exercises the AddVerifierData choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) AddVerifierData(contractID string, args AddVerifierData) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddVerifierData",
		Arguments:  argsToMap(args),
	}
}

// AddVerifierDataWithPackageID exercises the AddVerifierData choice using the provided package ID instead of package name
func (t SendingMessageV1) AddVerifierDataWithPackageID(contractID string, packageID string, args AddVerifierData) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddVerifierData",
		Arguments:  argsToMap(args),
	}
}

// AddCCVFee exercises the AddCCVFee choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) AddCCVFee(contractID string, args AddCCVFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddCCVFee",
		Arguments:  argsToMap(args),
	}
}

// AddCCVFeeWithPackageID exercises the AddCCVFee choice using the provided package ID instead of package name
func (t SendingMessageV1) AddCCVFeeWithPackageID(contractID string, packageID string, args AddCCVFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddCCVFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFee exercises the AddTokenSendFee choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) AddTokenSendFee(contractID string, args AddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFeeWithPackageID exercises the AddTokenSendFee choice using the provided package ID instead of package name
func (t SendingMessageV1) AddTokenSendFeeWithPackageID(contractID string, packageID string, args AddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVs exercises the SetOutboundPoolCCVs choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) SetOutboundPoolCCVs(contractID string, args SetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVsWithPackageID exercises the SetOutboundPoolCCVs choice using the provided package ID instead of package name
func (t SendingMessageV1) SetOutboundPoolCCVsWithPackageID(contractID string, packageID string, args SetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// AddExecutorFee exercises the AddExecutorFee choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) AddExecutorFee(contractID string, args AddExecutorFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddExecutorFee",
		Arguments:  argsToMap(args),
	}
}

// AddExecutorFeeWithPackageID exercises the AddExecutorFee choice using the provided package ID instead of package name
func (t SendingMessageV1) AddExecutorFeeWithPackageID(contractID string, packageID string, args AddExecutorFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "AddExecutorFee",
		Arguments:  argsToMap(args),
	}
}

// FinalizeSend exercises the FinalizeSend choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) FinalizeSend(contractID string, args FinalizeSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeSend",
		Arguments:  argsToMap(args),
	}
}

// FinalizeSendWithPackageID exercises the FinalizeSend choice using the provided package ID instead of package name
func (t SendingMessageV1) FinalizeSendWithPackageID(contractID string, packageID string, args FinalizeSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FinalizeSend",
		Arguments:  argsToMap(args),
	}
}

// FeeTokenAmount exercises the FeeTokenAmount choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) FeeTokenAmount(contractID string, args FeeTokenAmount) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FeeTokenAmount",
		Arguments:  argsToMap(args),
	}
}

// FeeTokenAmountWithPackageID exercises the FeeTokenAmount choice using the provided package ID instead of package name
func (t SendingMessageV1) FeeTokenAmountWithPackageID(contractID string, packageID string, args FeeTokenAmount) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "FeeTokenAmount",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this SendingMessageV1 contract
// This method uses the package name in the template ID
func (t SendingMessageV1) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t SendingMessageV1) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessageV1"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetConfig is a Record type
type SetConfig struct {
	NewIsEnabled types.BOOL    `json:"newIsEnabled"`
	NewCapacity  types.NUMERIC `json:"newCapacity"`
	NewRate      types.NUMERIC `json:"newRate"`
}

// ToMap converts SetConfig to a map for DAML arguments
func (t SetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["newIsEnabled"] = bool(t.NewIsEnabled)

	m["newCapacity"] = t.NewCapacity

	m["newRate"] = t.NewRate

	return m
}

func (t SetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetConfig to hex string (Canton MCMS format)
func (t SetConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetConfig from hex string (Canton MCMS format)
func (t *SetConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetInboundPoolCCVs is a Record type
type SetInboundPoolCCVs struct {
	PoolCCVs []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// ToMap converts SetInboundPoolCCVs to a map for DAML arguments
func (t SetInboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
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

func (t SetInboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetInboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetInboundPoolCCVs to hex string (Canton MCMS format)
func (t SetInboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetInboundPoolCCVs from hex string (Canton MCMS format)
func (t *SetInboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetOutboundPoolCCVs is a Record type
type SetOutboundPoolCCVs struct {
	PoolCCVs []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// ToMap converts SetOutboundPoolCCVs to a map for DAML arguments
func (t SetOutboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
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

func (t SetOutboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetOutboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetOutboundPoolCCVs to hex string (Canton MCMS format)
func (t SetOutboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetOutboundPoolCCVs from hex string (Canton MCMS format)
func (t *SetOutboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        types.BOOL                `json:"isEnabled"`
	OnRampAddresses  []types.TEXT              `json:"onRampAddresses"`
	DefaultCCVs      []mcms.RawInstanceAddress `json:"defaultCCVs"`
	LaneMandatedCCVs []mcms.RawInstanceAddress `json:"laneMandatedCCVs"`
}

// ToMap converts SourceChainConfig to a map for DAML arguments
func (t SourceChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["onRampAddresses"] = func() []any {
		res := make([]any, 0, len(t.OnRampAddresses))
		for _, e := range t.OnRampAddresses {
			res = append(res, string(e))
		}
		return res
	}()

	m["defaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
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

func (t SourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SourceChainConfig to hex string (Canton MCMS format)
func (t SourceChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SourceChainConfig from hex string (Canton MCMS format)
func (t *SourceChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SourceChainConfigArgs is a Record type
type SourceChainConfigArgs struct {
	SourceChainSelector types.NUMERIC             `json:"sourceChainSelector"`
	IsEnabled           types.BOOL                `json:"isEnabled"`
	OnRampAddresses     []types.TEXT              `json:"onRampAddresses"`
	DefaultCCVs         []mcms.RawInstanceAddress `json:"defaultCCVs"`
	LaneMandatedCCVs    []mcms.RawInstanceAddress `json:"laneMandatedCCVs"`
}

// ToMap converts SourceChainConfigArgs to a map for DAML arguments
func (t SourceChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["isEnabled"] = bool(t.IsEnabled)

	m["onRampAddresses"] = func() []any {
		res := make([]any, 0, len(t.OnRampAddresses))
		for _, e := range t.OnRampAddresses {
			res = append(res, string(e))
		}
		return res
	}()

	m["defaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
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

func (t SourceChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SourceChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SourceChainConfigArgs to hex string (Canton MCMS format)
func (t SourceChainConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SourceChainConfigArgs from hex string (Canton MCMS format)
func (t *SourceChainConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAmount is a Record type
type TokenAmount struct {
	Token           splice_api_token_holding_v1.InstrumentId `json:"token"`
	Amount          types.NUMERIC                            `json:"amount"`
	SenderInputCids []types.CONTRACT_ID                      `json:"senderInputCids"`
}

// ToMap converts TokenAmount to a map for DAML arguments
func (t TokenAmount) ToMap() map[string]any {
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

	return m
}

func (t TokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAmount to hex string (Canton MCMS format)
func (t TokenAmount) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAmount from hex string (Canton MCMS format)
func (t *TokenAmount) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner                    types.PARTY                              `json:"ccipOwner"`
	CcvOwners                    []types.PARTY                            `json:"ccvOwners"`
	VerifiedCCVs                 []mcms.RawInstanceAddress                `json:"verifiedCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	PoolOwner                    types.PARTY                              `json:"poolOwner"`
	Receiver                     types.PARTY                              `json:"receiver"`
	TokenReceiver                types.PARTY                              `json:"tokenReceiver"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount                       types.NUMERIC                            `json:"amount"`
	SourcePoolData               types.TEXT                               `json:"sourcePoolData"`
	MessageId                    types.TEXT                               `json:"messageId"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
	Finality                     types.INT64                              `json:"finality"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenReceiveTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Tickets", "TokenReceiveTicket")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenReceiveTicket) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Tickets", "TokenReceiveTicket")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenReceiveTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = int64(t.Finality)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenReceiveTicket) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = int64(t.Finality)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicket to hex string (Canton MCMS format)
func (t TokenReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicket from hex string (Canton MCMS format)
func (t *TokenReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenReceiveTicket

// Archive exercises the Archive choice on this TokenReceiveTicket contract
// This method uses the package name in the template ID
func (t TokenReceiveTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenReceiveTicket) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TokenReceiveTicketClaimed is a Template type
type TokenReceiveTicketClaimed struct {
	CcipOwner             types.PARTY                    `json:"ccipOwner"`
	CcvOwners             []types.PARTY                  `json:"ccvOwners"`
	PoolOwner             types.PARTY                    `json:"poolOwner"`
	Receiver              types.PARTY                    `json:"receiver"`
	TokenReceiver         types.PARTY                    `json:"tokenReceiver"`
	TokenReceiveTicketCid types.CONTRACT_ID              `json:"tokenReceiveTicketCid"`
	Event                 TokenReceiveTicketClaimedEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenReceiveTicketClaimed) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "TokenReceiveTicketClaimed")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Events", "TokenReceiveTicketClaimed")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenReceiveTicketClaimed) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiveTicketCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenReceiveTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenReceiveTicketCid
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiveTicketCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenReceiveTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenReceiveTicketCid
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Event).(mapper); ok {
			return m.toMap()
		}
		return t.Event
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenReceiveTicketClaimed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimed to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimed) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimed from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimed) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenReceiveTicketClaimed

// Archive exercises the Archive choice on this TokenReceiveTicketClaimed contract
// This method uses the package name in the template ID
func (t TokenReceiveTicketClaimed) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Events", "TokenReceiveTicketClaimed"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenReceiveTicketClaimed) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Events", "TokenReceiveTicketClaimed"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TokenReceiveTicketClaimedEvent is a Record type
type TokenReceiveTicketClaimedEvent struct {
	VerifiedCCVs                 []mcms.RawInstanceAddress                `json:"verifiedCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount                       types.NUMERIC                            `json:"amount"`
	SourcePoolData               types.TEXT                               `json:"sourcePoolData"`
	MessageId                    types.TEXT                               `json:"messageId"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
	Finality                     types.INT64                              `json:"finality"`
	Output                       TokenReceiveTicketClaimedOutput          `json:"output"`
}

// ToMap converts TokenReceiveTicketClaimedEvent to a map for DAML arguments
func (t TokenReceiveTicketClaimedEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = t.Amount

	m["sourcePoolData"] = string(t.SourcePoolData)

	m["messageId"] = string(t.MessageId)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["finality"] = int64(t.Finality)

	m["output"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Output).(mapper); ok {
			return m.toMap()
		}
		return t.Output
	}()

	return m
}

func (t TokenReceiveTicketClaimedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedEvent to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedEvent from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketClaimedCompleted is a Record type
type TokenReceiveTicketClaimedCompleted struct {
	ReceiverHoldingCids []types.CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts TokenReceiveTicketClaimedCompleted to a map for DAML arguments
func (t TokenReceiveTicketClaimedCompleted) ToMap() map[string]any {
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

func (t TokenReceiveTicketClaimedCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedCompleted to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedCompleted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedCompleted from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedCompleted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketClaimedOutput is a variant/union type
type TokenReceiveTicketClaimedOutput struct {
	TokenReceiveTicketClaimedPending   *TokenReceiveTicketClaimedPending   `json:"TokenReceiveTicketClaimed_Pending,omitempty"`
	TokenReceiveTicketClaimedCompleted *TokenReceiveTicketClaimedCompleted `json:"TokenReceiveTicketClaimed_Completed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TokenReceiveTicketClaimedOutput
func (v TokenReceiveTicketClaimedOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for TokenReceiveTicketClaimedOutput
func (v *TokenReceiveTicketClaimedOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes TokenReceiveTicketClaimedOutput to hex string (Canton MCMS format)
func (v TokenReceiveTicketClaimedOutput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedOutput from hex string (Canton MCMS format)
func (v *TokenReceiveTicketClaimedOutput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TokenReceiveTicketClaimedOutput) GetVariantTag() string {

	if v.TokenReceiveTicketClaimedPending != nil {
		return "TokenReceiveTicketClaimed_Pending"
	}

	if v.TokenReceiveTicketClaimedCompleted != nil {
		return "TokenReceiveTicketClaimed_Completed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TokenReceiveTicketClaimedOutput) GetVariantValue() any {

	if v.TokenReceiveTicketClaimedPending != nil {
		return v.TokenReceiveTicketClaimedPending
	}

	if v.TokenReceiveTicketClaimedCompleted != nil {
		return v.TokenReceiveTicketClaimedCompleted
	}

	return nil
}

var _ types.VARIANT = (*TokenReceiveTicketClaimedOutput)(nil)

// TokenReceiveTicketClaimedPending is a Record type
type TokenReceiveTicketClaimedPending struct {
	TransferInstructionCid types.CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts TokenReceiveTicketClaimedPending to a map for DAML arguments
func (t TokenReceiveTicketClaimedPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferInstructionCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TransferInstructionCid).(mapper); ok {
			return m.toMap()
		}
		return t.TransferInstructionCid
	}()

	return m
}

func (t TokenReceiveTicketClaimedPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketClaimedPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketClaimedPending to hex string (Canton MCMS format)
func (t TokenReceiveTicketClaimedPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketClaimedPending from hex string (Canton MCMS format)
func (t *TokenReceiveTicketClaimedPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenSendData is a Record type
type TokenSendData struct {
	PoolInstanceId   types.TEXT                               `json:"poolInstanceId"`
	PoolOwner        types.PARTY                              `json:"poolOwner"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount           types.NUMERIC                            `json:"amount"`
	DestTokenAddress types.TEXT                               `json:"destTokenAddress"`
	ExtraData        types.TEXT                               `json:"extraData"`
}

// ToMap converts TokenSendData to a map for DAML arguments
func (t TokenSendData) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = t.Amount

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["extraData"] = string(t.ExtraData)

	return m
}

func (t TokenSendData) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenSendData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenSendData to hex string (Canton MCMS format)
func (t TokenSendData) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenSendData from hex string (Canton MCMS format)
func (t *TokenSendData) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenSendFee is a Record type
type TokenSendFee struct {
	PoolInstanceId    types.TEXT    `json:"poolInstanceId"`
	PoolOwner         types.PARTY   `json:"poolOwner"`
	FeeUSDCents       types.NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts TokenSendFee to a map for DAML arguments
func (t TokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t TokenSendFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenSendFee to hex string (Canton MCMS format)
func (t TokenSendFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenSendFee from hex string (Canton MCMS format)
func (t *TokenSendFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferV1 is a Record type
type TokenTransferV1 struct {
	Amount             types.NUMERIC `json:"amount"`
	SourcePoolAddress  types.TEXT    `json:"sourcePoolAddress"`
	SourceTokenAddress types.TEXT    `json:"sourceTokenAddress"`
	DestTokenAddress   types.TEXT    `json:"destTokenAddress"`
	TokenReceiver      types.TEXT    `json:"tokenReceiver"`
	ExtraData          types.TEXT    `json:"extraData"`
}

// ToMap converts TokenTransferV1 to a map for DAML arguments
func (t TokenTransferV1) ToMap() map[string]any {
	m := make(map[string]any)

	m["amount"] = t.Amount

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

// VerifierData is a Record type
type VerifierData struct {
	CcvInstanceId        types.TEXT    `json:"ccvInstanceId"`
	CcvOwner             types.PARTY   `json:"ccvOwner"`
	VersionTag           types.TEXT    `json:"versionTag"`
	VerifierBlob         types.TEXT    `json:"verifierBlob"`
	MessageSentObservers []types.PARTY `json:"messageSentObservers"`
}

// ToMap converts VerifierData to a map for DAML arguments
func (t VerifierData) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["versionTag"] = string(t.VersionTag)

	m["verifierBlob"] = string(t.VerifierBlob)

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t VerifierData) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *VerifierData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes VerifierData to hex string (Canton MCMS format)
func (t VerifierData) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes VerifierData from hex string (Canton MCMS format)
func (t *VerifierData) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddCCVFee(args AddCCVFee) (*bind.EncodedChoice, error)
	AddCCVFeeMCMSParams(args AddCCVFeeMCMSParams) (*bind.EncodedChoice, error)
	AddCCVVerification(args AddCCVVerification) (*bind.EncodedChoice, error)
	AddCCVVerificationMCMSParams(args AddCCVVerificationMCMSParams) (*bind.EncodedChoice, error)
	AddExecutorFee(args AddExecutorFee) (*bind.EncodedChoice, error)
	AddExecutorFeeMCMSParams(args AddExecutorFeeMCMSParams) (*bind.EncodedChoice, error)
	AddTokenSend(args AddTokenSend) (*bind.EncodedChoice, error)
	AddTokenSendFee(args AddTokenSendFee) (*bind.EncodedChoice, error)
	AddVerifierData(args AddVerifierData) (*bind.EncodedChoice, error)
	AddVerifierDataMCMSParams(args AddVerifierDataMCMSParams) (*bind.EncodedChoice, error)
	ApplyDestChainConfigUpdates(args ApplyDestChainConfigUpdates) (*bind.EncodedChoice, error)
	ApplySourceChainConfigUpdates(args ApplySourceChainConfigUpdates) (*bind.EncodedChoice, error)
	BuildMessage(args BuildMessage) (*bind.EncodedChoice, error)
	BuildMessageMCMSParams(args BuildMessageMCMSParams) (*bind.EncodedChoice, error)
	CancelExecute(args CancelExecute) (*bind.EncodedChoice, error)
	CancelExecuteMCMSParams(args CancelExecuteMCMSParams) (*bind.EncodedChoice, error)
	ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error)
	FeeTokenAmount(args FeeTokenAmount) (*bind.EncodedChoice, error)
	FeeTokenAmountMCMSParams(args FeeTokenAmountMCMSParams) (*bind.EncodedChoice, error)
	FinalizeExecute(args FinalizeExecute) (*bind.EncodedChoice, error)
	FinalizeFee(args FinalizeFee) (*bind.EncodedChoice, error)
	FinalizeSend(args FinalizeSend) (*bind.EncodedChoice, error)
	GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error)
	GetDestChainConfigMCMSParams(args GetDestChainConfigMCMSParams) (*bind.EncodedChoice, error)
	GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error)
	GetSourceChainConfigMCMSParams(args GetSourceChainConfigMCMSParams) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
	SetInboundPoolCCVs(args SetInboundPoolCCVs) (*bind.EncodedChoice, error)
	SetOutboundPoolCCVs(args SetOutboundPoolCCVs) (*bind.EncodedChoice, error)
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

// AddCCVFee encodes parameters for the AddCCVFee choice.
func (e *encoder) AddCCVFee(args AddCCVFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCCVFee", args)
}

// AddCCVFeeMCMSParams encodes MCMS parameters (without Caller) for the AddCCVFee choice.
func (e *encoder) AddCCVFeeMCMSParams(args AddCCVFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCCVFee", args)
}

// AddCCVVerification encodes parameters for the AddCCVVerification choice.
func (e *encoder) AddCCVVerification(args AddCCVVerification) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCCVVerification", args)
}

// AddCCVVerificationMCMSParams encodes MCMS parameters (without Caller) for the AddCCVVerification choice.
func (e *encoder) AddCCVVerificationMCMSParams(args AddCCVVerificationMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCCVVerification", args)
}

// AddExecutorFee encodes parameters for the AddExecutorFee choice.
func (e *encoder) AddExecutorFee(args AddExecutorFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddExecutorFee", args)
}

// AddExecutorFeeMCMSParams encodes MCMS parameters (without Caller) for the AddExecutorFee choice.
func (e *encoder) AddExecutorFeeMCMSParams(args AddExecutorFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddExecutorFee", args)
}

// AddTokenSend encodes parameters for the AddTokenSend choice.
func (e *encoder) AddTokenSend(args AddTokenSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSend", args)
}

// AddTokenSendFee encodes parameters for the AddTokenSendFee choice.
func (e *encoder) AddTokenSendFee(args AddTokenSendFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddTokenSendFee", args)
}

// AddVerifierData encodes parameters for the AddVerifierData choice.
func (e *encoder) AddVerifierData(args AddVerifierData) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddVerifierData", args)
}

// AddVerifierDataMCMSParams encodes MCMS parameters (without Caller) for the AddVerifierData choice.
func (e *encoder) AddVerifierDataMCMSParams(args AddVerifierDataMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddVerifierData", args)
}

// ApplyDestChainConfigUpdates encodes parameters for the ApplyDestChainConfigUpdates choice.
func (e *encoder) ApplyDestChainConfigUpdates(args ApplyDestChainConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyDestChainConfigUpdates", args)
}

// ApplySourceChainConfigUpdates encodes parameters for the ApplySourceChainConfigUpdates choice.
func (e *encoder) ApplySourceChainConfigUpdates(args ApplySourceChainConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplySourceChainConfigUpdates", args)
}

// BuildMessage encodes parameters for the BuildMessage choice.
func (e *encoder) BuildMessage(args BuildMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BuildMessage", args)
}

// BuildMessageMCMSParams encodes MCMS parameters (without Caller) for the BuildMessage choice.
func (e *encoder) BuildMessageMCMSParams(args BuildMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BuildMessage", args)
}

// CancelExecute encodes parameters for the CancelExecute choice.
func (e *encoder) CancelExecute(args CancelExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelExecute", args)
}

// CancelExecuteMCMSParams encodes MCMS parameters (without Caller) for the CancelExecute choice.
func (e *encoder) CancelExecuteMCMSParams(args CancelExecuteMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelExecute", args)
}

// ConsumeCapacity encodes parameters for the ConsumeCapacity choice.
func (e *encoder) ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeCapacity", args)
}

// FeeTokenAmount encodes parameters for the FeeTokenAmount choice.
func (e *encoder) FeeTokenAmount(args FeeTokenAmount) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FeeTokenAmount", args)
}

// FeeTokenAmountMCMSParams encodes MCMS parameters (without Caller) for the FeeTokenAmount choice.
func (e *encoder) FeeTokenAmountMCMSParams(args FeeTokenAmountMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FeeTokenAmount", args)
}

// FinalizeExecute encodes parameters for the FinalizeExecute choice.
func (e *encoder) FinalizeExecute(args FinalizeExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeExecute", args)
}

// FinalizeFee encodes parameters for the FinalizeFee choice.
func (e *encoder) FinalizeFee(args FinalizeFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeFee", args)
}

// FinalizeSend encodes parameters for the FinalizeSend choice.
func (e *encoder) FinalizeSend(args FinalizeSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeSend", args)
}

// GetDestChainConfig encodes parameters for the GetDestChainConfig choice.
func (e *encoder) GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig", args)
}

// GetDestChainConfigMCMSParams encodes MCMS parameters (without Caller) for the GetDestChainConfig choice.
func (e *encoder) GetDestChainConfigMCMSParams(args GetDestChainConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig", args)
}

// GetSourceChainConfig encodes parameters for the GetSourceChainConfig choice.
func (e *encoder) GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSourceChainConfig", args)
}

// GetSourceChainConfigMCMSParams encodes MCMS parameters (without Caller) for the GetSourceChainConfig choice.
func (e *encoder) GetSourceChainConfigMCMSParams(args GetSourceChainConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSourceChainConfig", args)
}

// SetConfig encodes parameters for the SetConfig choice.
func (e *encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetInboundPoolCCVs encodes parameters for the SetInboundPoolCCVs choice.
func (e *encoder) SetInboundPoolCCVs(args SetInboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetInboundPoolCCVs", args)
}

// SetOutboundPoolCCVs encodes parameters for the SetOutboundPoolCCVs choice.
func (e *encoder) SetOutboundPoolCCVs(args SetOutboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetOutboundPoolCCVs", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
