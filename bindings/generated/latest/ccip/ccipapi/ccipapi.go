package ccipapi

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipcodec "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipcodec"
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
	PackageName = "ccip-api"
	PackageID   = "78dcbd60e8aa4c6d967656ac8e23d587b2c0ea4cca4ab52391abc1733444856e"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IExecutingMessage is a DAML interface
type IExecutingMessage interface {

	// ExecutingMessageAddCCVVerification executes the ExecutingMessage_AddCCVVerification choice
	ExecutingMessageAddCCVVerification(contractID string, args ExecutingMessageAddCCVVerification) *model.ExerciseCommand

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// IRMNRemote is a DAML interface
type IRMNRemote interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// RMNRemotePublicFetch executes the RMNRemote_PublicFetch choice
	RMNRemotePublicFetch(contractID string, args RMNRemotePublicFetch) *model.ExerciseCommand

	// RMNRemoteIsCursed executes the RMNRemote_IsCursed choice
	RMNRemoteIsCursed(contractID string, args RMNRemoteIsCursed) *model.ExerciseCommand

	// RMNRemoteIsCursedForChain executes the RMNRemote_IsCursedForChain choice
	RMNRemoteIsCursedForChain(contractID string, args RMNRemoteIsCursedForChain) *model.ExerciseCommand

	// RMNRemoteGetCursedSubjects executes the RMNRemote_GetCursedSubjects choice
	RMNRemoteGetCursedSubjects(contractID string, args RMNRemoteGetCursedSubjects) *model.ExerciseCommand
}

// ISendingMessage is a DAML interface
type ISendingMessage interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// SendingMessageAddCCVFee executes the SendingMessage_AddCCVFee choice
	SendingMessageAddCCVFee(contractID string, args SendingMessageAddCCVFee) *model.ExerciseCommand

	// SendingMessageAddVerifierData executes the SendingMessage_AddVerifierData choice
	SendingMessageAddVerifierData(contractID string, args SendingMessageAddVerifierData) *model.ExerciseCommand

	// SendingMessageAddExecutorFee executes the SendingMessage_AddExecutorFee choice
	SendingMessageAddExecutorFee(contractID string, args SendingMessageAddExecutorFee) *model.ExerciseCommand
}

// ITokenAdminRegistry is a DAML interface
type ITokenAdminRegistry interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenAdminRegistryPublicFetch executes the TokenAdminRegistry_PublicFetch choice
	TokenAdminRegistryPublicFetch(contractID string, args TokenAdminRegistryPublicFetch) *model.ExerciseCommand

	// TokenAdminRegistryFetchTokenConfig executes the TokenAdminRegistry_FetchTokenConfig choice
	TokenAdminRegistryFetchTokenConfig(contractID string, args TokenAdminRegistryFetchTokenConfig) *model.ExerciseCommand

	// TokenAdminRegistrySetInboundPoolCCVs executes the TokenAdminRegistry_SetInboundPoolCCVs choice
	TokenAdminRegistrySetInboundPoolCCVs(contractID string, args TokenAdminRegistrySetInboundPoolCCVs) *model.ExerciseCommand

	// TokenAdminRegistrySetOutboundPoolCCVs executes the TokenAdminRegistry_SetOutboundPoolCCVs choice
	TokenAdminRegistrySetOutboundPoolCCVs(contractID string, args TokenAdminRegistrySetOutboundPoolCCVs) *model.ExerciseCommand

	// TokenAdminRegistryAddTokenSend executes the TokenAdminRegistry_AddTokenSend choice
	TokenAdminRegistryAddTokenSend(contractID string, args TokenAdminRegistryAddTokenSend) *model.ExerciseCommand

	// TokenAdminRegistryAddTokenSendFee executes the TokenAdminRegistry_AddTokenSendFee choice
	TokenAdminRegistryAddTokenSendFee(contractID string, args TokenAdminRegistryAddTokenSendFee) *model.ExerciseCommand

	// TokenAdminRegistryConsumeReceiveTicket executes the TokenAdminRegistry_ConsumeReceiveTicket choice
	TokenAdminRegistryConsumeReceiveTicket(contractID string, args TokenAdminRegistryConsumeReceiveTicket) *model.ExerciseCommand
}

// ITokenConfig is a DAML interface
type ITokenConfig interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenConfigPublicFetch executes the TokenConfig_PublicFetch choice
	TokenConfigPublicFetch(contractID string, args TokenConfigPublicFetch) *model.ExerciseCommand

	// TokenConfigAssertConfiguredTransferFactory executes the TokenConfig_AssertConfiguredTransferFactory choice
	TokenConfigAssertConfiguredTransferFactory(contractID string, args TokenConfigAssertConfiguredTransferFactory) *model.ExerciseCommand

	// TokenConfigAssertConfiguredBurnMintFactory executes the TokenConfig_AssertConfiguredBurnMintFactory choice
	TokenConfigAssertConfiguredBurnMintFactory(contractID string, args TokenConfigAssertConfiguredBurnMintFactory) *model.ExerciseCommand
}

// ITokenReceiveTicket is a DAML interface
type ITokenReceiveTicket interface {

	// Consume executes the Consume choice
	Consume(contractID string, args Consume) *model.ExerciseCommand

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
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

// Consume is a Record type
type Consume struct {
}

// ToMap converts Consume to a map for DAML arguments
func (t Consume) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t Consume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Consume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Consume to hex string (Canton MCMS format)
func (t Consume) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Consume from hex string (Canton MCMS format)
func (t *Consume) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutingMessageView is a Record type
type ExecutingMessageView struct {
	CcipOwner          types.PARTY                                `json:"ccipOwner"`
	Message            ccipcodec.MessageV1                        `json:"message"`
	OffRamp            chainlinkapi.RawInstanceAddress            `json:"offRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress            `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress            `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress            `json:"tokenAdminRegistry"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ExecutingMessageView to a map for DAML arguments
func (t ExecutingMessageView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["offRamp"] = model.NestedToDAMLValue(t.OffRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ExecutingMessageView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageView to hex string (Canton MCMS format)
func (t ExecutingMessageView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageView from hex string (Canton MCMS format)
func (t *ExecutingMessageView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutingMessageAddCCVVerification is a Record type
type ExecutingMessageAddCCVVerification struct {
	CcvInstanceId types.TEXT                                 `json:"ccvInstanceId"`
	VersionTag    types.TEXT                                 `json:"versionTag" hex:"bytes"`
	Context       splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller        types.PARTY                                `json:"caller"`
}

// ToMap converts ExecutingMessageAddCCVVerification to a map for DAML arguments
func (t ExecutingMessageAddCCVVerification) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["versionTag"] = string(t.VersionTag)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutingMessageAddCCVVerification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageAddCCVVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageAddCCVVerification to hex string (Canton MCMS format)
func (t ExecutingMessageAddCCVVerification) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageAddCCVVerification from hex string (Canton MCMS format)
func (t *ExecutingMessageAddCCVVerification) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.ExecutingMessageV1", "MessageExecutionState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e MessageExecutionState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.API.ExecutingMessageV1", "MessageExecutionState")
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

// RMNRemoteView is a Record type
type RMNRemoteView struct {
	CcipOwner  types.PARTY `json:"ccipOwner"`
	RmnOwner   types.PARTY `json:"rmnOwner"`
	InstanceId types.TEXT  `json:"instanceId"`
}

// ToMap converts RMNRemoteView to a map for DAML arguments
func (t RMNRemoteView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["rmnOwner"] = t.RmnOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t RMNRemoteView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemoteView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemoteView to hex string (Canton MCMS format)
func (t RMNRemoteView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemoteView from hex string (Canton MCMS format)
func (t *RMNRemoteView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RMNRemoteGetCursedSubjects is a Record type
type RMNRemoteGetCursedSubjects struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts RMNRemoteGetCursedSubjects to a map for DAML arguments
func (t RMNRemoteGetCursedSubjects) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t RMNRemoteGetCursedSubjects) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemoteGetCursedSubjects) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemoteGetCursedSubjects to hex string (Canton MCMS format)
func (t RMNRemoteGetCursedSubjects) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemoteGetCursedSubjects from hex string (Canton MCMS format)
func (t *RMNRemoteGetCursedSubjects) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RMNRemoteIsCursed is a Record type
type RMNRemoteIsCursed struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts RMNRemoteIsCursed to a map for DAML arguments
func (t RMNRemoteIsCursed) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t RMNRemoteIsCursed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemoteIsCursed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemoteIsCursed to hex string (Canton MCMS format)
func (t RMNRemoteIsCursed) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemoteIsCursed from hex string (Canton MCMS format)
func (t *RMNRemoteIsCursed) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RMNRemoteIsCursedForChain is a Record type
type RMNRemoteIsCursedForChain struct {
	ChainSelector types.NUMERIC                              `json:"chainSelector"`
	Context       splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller        types.PARTY                                `json:"caller"`
}

// ToMap converts RMNRemoteIsCursedForChain to a map for DAML arguments
func (t RMNRemoteIsCursedForChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t RMNRemoteIsCursedForChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemoteIsCursedForChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemoteIsCursedForChain to hex string (Canton MCMS format)
func (t RMNRemoteIsCursedForChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemoteIsCursedForChain from hex string (Canton MCMS format)
func (t *RMNRemoteIsCursedForChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RMNRemotePublicFetch is a Record type
type RMNRemotePublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress `json:"expectedAddress"`
	Caller          types.PARTY                     `json:"caller"`
}

// ToMap converts RMNRemotePublicFetch to a map for DAML arguments
func (t RMNRemotePublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t RMNRemotePublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemotePublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemotePublicFetch to hex string (Canton MCMS format)
func (t RMNRemotePublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemotePublicFetch from hex string (Canton MCMS format)
func (t *RMNRemotePublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageView is a Record type
type SendingMessageView struct {
	CcipOwner          types.PARTY                                `json:"ccipOwner"`
	Sender             types.PARTY                                `json:"sender"`
	DestChainSelector  types.NUMERIC                              `json:"destChainSelector"`
	RequiredCCVs       []chainlinkapi.RawInstanceAddress          `json:"requiredCCVs"`
	OutboundPoolCCVs   *[]chainlinkapi.RawInstanceAddress         `json:"outboundPoolCCVs" hex:"optional"`
	Router             chainlinkapi.RawInstanceAddress            `json:"router"`
	OnRamp             chainlinkapi.RawInstanceAddress            `json:"onRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress            `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress            `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress            `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress            `json:"feeQuoter"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts SendingMessageView to a map for DAML arguments
func (t SendingMessageView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["sender"] = t.Sender.ToMap()

	m["destChainSelector"] = t.DestChainSelector

	m["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.OutboundPoolCCVs != nil {
		m["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OutboundPoolCCVs),
		}
	} else {
		m["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["onRamp"] = model.NestedToDAMLValue(t.OnRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t SendingMessageView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageView to hex string (Canton MCMS format)
func (t SendingMessageView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageView from hex string (Canton MCMS format)
func (t *SendingMessageView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageAddCCVFee is a Record type
type SendingMessageAddCCVFee struct {
	CcvInstanceId     types.TEXT                                 `json:"ccvInstanceId"`
	FeeUSDCents       types.NUMERIC                              `json:"feeUSDCents"`
	DestGasLimit      types.INT64                                `json:"destGasLimit"`
	DestBytesOverhead types.INT64                                `json:"destBytesOverhead"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts SendingMessageAddCCVFee to a map for DAML arguments
func (t SendingMessageAddCCVFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SendingMessageAddCCVFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageAddCCVFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageAddCCVFee to hex string (Canton MCMS format)
func (t SendingMessageAddCCVFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageAddCCVFee from hex string (Canton MCMS format)
func (t *SendingMessageAddCCVFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageAddExecutorFee is a Record type
type SendingMessageAddExecutorFee struct {
	ExecutorInstanceId types.TEXT                                 `json:"executorInstanceId"`
	ExecutorArgs       types.TEXT                                 `json:"executorArgs"`
	FeeUSDCents        types.NUMERIC                              `json:"feeUSDCents"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller             types.PARTY                                `json:"caller"`
}

// ToMap converts SendingMessageAddExecutorFee to a map for DAML arguments
func (t SendingMessageAddExecutorFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorInstanceId"] = string(t.ExecutorInstanceId)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["feeUSDCents"] = t.FeeUSDCents

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SendingMessageAddExecutorFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageAddExecutorFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageAddExecutorFee to hex string (Canton MCMS format)
func (t SendingMessageAddExecutorFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageAddExecutorFee from hex string (Canton MCMS format)
func (t *SendingMessageAddExecutorFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageAddVerifierData is a Record type
type SendingMessageAddVerifierData struct {
	CcvInstanceId        types.TEXT                                 `json:"ccvInstanceId"`
	VersionTag           types.TEXT                                 `json:"versionTag" hex:"bytes"`
	VerifierBlob         types.TEXT                                 `json:"verifierBlob"`
	MessageSentObservers []types.PARTY                              `json:"messageSentObservers"`
	Context              splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller               types.PARTY                                `json:"caller"`
}

// ToMap converts SendingMessageAddVerifierData to a map for DAML arguments
func (t SendingMessageAddVerifierData) ToMap() map[string]any {
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

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SendingMessageAddVerifierData) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageAddVerifierData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageAddVerifierData to hex string (Canton MCMS format)
func (t SendingMessageAddVerifierData) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageAddVerifierData from hex string (Canton MCMS format)
func (t *SendingMessageAddVerifierData) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryView is a Record type
type TokenAdminRegistryView struct {
	CcipOwner  types.PARTY `json:"ccipOwner"`
	InstanceId types.TEXT  `json:"instanceId"`
}

// ToMap converts TokenAdminRegistryView to a map for DAML arguments
func (t TokenAdminRegistryView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t TokenAdminRegistryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryView to hex string (Canton MCMS format)
func (t TokenAdminRegistryView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryView from hex string (Canton MCMS format)
func (t *TokenAdminRegistryView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSend is a Record type
type TokenAdminRegistryAddTokenSend struct {
	TokenConfigCid    types.CONTRACT_ID                        `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                               `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount            types.TEXT                               `json:"amount"`
	DestTokenAddress  types.TEXT                               `json:"destTokenAddress"`
	ExtraData         types.TEXT                               `json:"extraData"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryAddTokenSend to a map for DAML arguments
func (t TokenAdminRegistryAddTokenSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = string(t.Amount)

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["extraData"] = string(t.ExtraData)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAddTokenSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAddTokenSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAddTokenSend to hex string (Canton MCMS format)
func (t TokenAdminRegistryAddTokenSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSend from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAddTokenSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSendFee is a Record type
type TokenAdminRegistryAddTokenSendFee struct {
	TokenConfigCid    types.CONTRACT_ID `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT        `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC     `json:"feeUSDCents"`
	DestGasOverhead   types.INT64       `json:"destGasOverhead"`
	DestBytesOverhead types.INT64       `json:"destBytesOverhead"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts TokenAdminRegistryAddTokenSendFee to a map for DAML arguments
func (t TokenAdminRegistryAddTokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAddTokenSendFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAddTokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAddTokenSendFee to hex string (Canton MCMS format)
func (t TokenAdminRegistryAddTokenSendFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSendFee from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAddTokenSendFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryConsumeReceiveTicket is a Record type
type TokenAdminRegistryConsumeReceiveTicket struct {
	TokenConfigCid        types.CONTRACT_ID                        `json:"tokenConfigCid"`
	TokenReceiveTicketCid types.CONTRACT_ID                        `json:"tokenReceiveTicketCid"`
	InstrumentId          splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	PoolInstanceId        types.TEXT                               `json:"poolInstanceId"`
	Caller                types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryConsumeReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryConsumeReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryConsumeReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryConsumeReceiveTicket to hex string (Canton MCMS format)
func (t TokenAdminRegistryConsumeReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryConsumeReceiveTicket from hex string (Canton MCMS format)
func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryFetchTokenConfig is a Record type
type TokenAdminRegistryFetchTokenConfig struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryFetchTokenConfig to a map for DAML arguments
func (t TokenAdminRegistryFetchTokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryFetchTokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryFetchTokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryFetchTokenConfig to hex string (Canton MCMS format)
func (t TokenAdminRegistryFetchTokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryFetchTokenConfig from hex string (Canton MCMS format)
func (t *TokenAdminRegistryFetchTokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryPublicFetch is a Record type
type TokenAdminRegistryPublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress `json:"expectedAddress"`
	Caller          types.PARTY                     `json:"caller"`
}

// ToMap converts TokenAdminRegistryPublicFetch to a map for DAML arguments
func (t TokenAdminRegistryPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryPublicFetch to hex string (Canton MCMS format)
func (t TokenAdminRegistryPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryPublicFetch from hex string (Canton MCMS format)
func (t *TokenAdminRegistryPublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetInboundPoolCCVs is a Record type
type TokenAdminRegistrySetInboundPoolCCVs struct {
	TokenConfigCid      types.CONTRACT_ID                 `json:"tokenConfigCid"`
	ExecutingMessageCid types.CONTRACT_ID                 `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                        `json:"poolInstanceId"`
	PoolCCVs            []chainlinkapi.RawInstanceAddress `json:"poolCCVs"`
	Caller              types.PARTY                       `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetInboundPoolCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetInboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetInboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetInboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetInboundPoolCCVs to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetInboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetInboundPoolCCVs from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetInboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetOutboundPoolCCVs is a Record type
type TokenAdminRegistrySetOutboundPoolCCVs struct {
	TokenConfigCid    types.CONTRACT_ID                 `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID                 `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                        `json:"poolInstanceId"`
	PoolCCVs          []chainlinkapi.RawInstanceAddress `json:"poolCCVs"`
	Caller            types.PARTY                       `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetOutboundPoolCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetOutboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetOutboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetOutboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetOutboundPoolCCVs to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetOutboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetOutboundPoolCCVs from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetOutboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfigView is a Record type
type TokenConfigView struct {
	CcipOwner    types.PARTY                              `json:"ccipOwner"`
	InstanceId   types.TEXT                               `json:"instanceId"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// ToMap converts TokenConfigView to a map for DAML arguments
func (t TokenConfigView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	return m
}

func (t TokenConfigView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfigView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfigView to hex string (Canton MCMS format)
func (t TokenConfigView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfigView from hex string (Canton MCMS format)
func (t *TokenConfigView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfigAssertConfiguredBurnMintFactory is a Record type
type TokenConfigAssertConfiguredBurnMintFactory struct {
	SuppliedFactory types.CONTRACT_ID `json:"suppliedFactory"`
	Caller          types.PARTY       `json:"caller"`
}

// ToMap converts TokenConfigAssertConfiguredBurnMintFactory to a map for DAML arguments
func (t TokenConfigAssertConfiguredBurnMintFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["suppliedFactory"] = model.NestedToDAMLValue(t.SuppliedFactory)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenConfigAssertConfiguredBurnMintFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfigAssertConfiguredBurnMintFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfigAssertConfiguredBurnMintFactory to hex string (Canton MCMS format)
func (t TokenConfigAssertConfiguredBurnMintFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfigAssertConfiguredBurnMintFactory from hex string (Canton MCMS format)
func (t *TokenConfigAssertConfiguredBurnMintFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfigAssertConfiguredTransferFactory is a Record type
type TokenConfigAssertConfiguredTransferFactory struct {
	SuppliedFactory types.CONTRACT_ID `json:"suppliedFactory"`
	Caller          types.PARTY       `json:"caller"`
}

// ToMap converts TokenConfigAssertConfiguredTransferFactory to a map for DAML arguments
func (t TokenConfigAssertConfiguredTransferFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["suppliedFactory"] = model.NestedToDAMLValue(t.SuppliedFactory)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenConfigAssertConfiguredTransferFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfigAssertConfiguredTransferFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfigAssertConfiguredTransferFactory to hex string (Canton MCMS format)
func (t TokenConfigAssertConfiguredTransferFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfigAssertConfiguredTransferFactory from hex string (Canton MCMS format)
func (t *TokenConfigAssertConfiguredTransferFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfigPublicFetch is a Record type
type TokenConfigPublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress `json:"expectedAddress"`
	Caller          types.PARTY                     `json:"caller"`
}

// ToMap converts TokenConfigPublicFetch to a map for DAML arguments
func (t TokenConfigPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenConfigPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfigPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfigPublicFetch to hex string (Canton MCMS format)
func (t TokenConfigPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfigPublicFetch from hex string (Canton MCMS format)
func (t *TokenConfigPublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicketView is a Record type
type TokenReceiveTicketView struct {
	CcipOwner types.PARTY   `json:"ccipOwner"`
	PoolOwner types.PARTY   `json:"poolOwner"`
	CcvOwners []types.PARTY `json:"ccvOwners"`
}

// ToMap converts TokenReceiveTicketView to a map for DAML arguments
func (t TokenReceiveTicketView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["ccvOwners"] = func() []any {
		res := make([]any, 0, len(t.CcvOwners))
		for _, e := range t.CcvOwners {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t TokenReceiveTicketView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenReceiveTicketView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenReceiveTicketView to hex string (Canton MCMS format)
func (t TokenReceiveTicketView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenReceiveTicketView from hex string (Canton MCMS format)
func (t *TokenReceiveTicketView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IExecutingMessageInterfaceID returns the interface ID for the IExecutingMessage interface using the package name
func IExecutingMessageInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.ExecutingMessageV1", "ExecutingMessage")
}

// IExecutingMessageInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IExecutingMessageInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.ExecutingMessageV1", "ExecutingMessage")
}

// IRMNRemoteInterfaceID returns the interface ID for the IRMNRemote interface using the package name
func IRMNRemoteInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.RMNRemoteV1", "RMNRemote")
}

// IRMNRemoteInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IRMNRemoteInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.RMNRemoteV1", "RMNRemote")
}

// ISendingMessageInterfaceID returns the interface ID for the ISendingMessage interface using the package name
func ISendingMessageInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.SendingMessageV1", "SendingMessage")
}

// ISendingMessageInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ISendingMessageInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.SendingMessageV1", "SendingMessage")
}

// ITokenAdminRegistryInterfaceID returns the interface ID for the ITokenAdminRegistry interface using the package name
func ITokenAdminRegistryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.TokenAdminRegistryV1", "TokenAdminRegistry")
}

// ITokenAdminRegistryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ITokenAdminRegistryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.TokenAdminRegistryV1", "TokenAdminRegistry")
}

// ITokenConfigInterfaceID returns the interface ID for the ITokenConfig interface using the package name
func ITokenConfigInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.TokenAdminRegistryV1", "TokenConfig")
}

// ITokenConfigInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ITokenConfigInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.TokenAdminRegistryV1", "TokenConfig")
}

// ITokenReceiveTicketInterfaceID returns the interface ID for the ITokenReceiveTicket interface using the package name
func ITokenReceiveTicketInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.API.ExecutingMessageV1", "TokenReceiveTicket")
}

// ITokenReceiveTicketInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ITokenReceiveTicketInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.API.ExecutingMessageV1", "TokenReceiveTicket")
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
