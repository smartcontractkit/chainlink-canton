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
	PackageName = "ccip-api-v2"
	PackageID   = "7fffaf108129d37413d8edfbd91ffe373051b2cb0621c26e245093c6138daf58"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IIExecutingMessage is a DAML interface
type IIExecutingMessage interface {

	// ExecutingMessageCancelExecute executes the ExecutingMessage_CancelExecute choice
	ExecutingMessageCancelExecute(contractID string, args ExecutingMessageCancelExecute) *model.ExerciseCommand

	// ExecutingMessageAddCCVVerification executes the ExecutingMessage_AddCCVVerification choice
	ExecutingMessageAddCCVVerification(contractID string, args ExecutingMessageAddCCVVerification) *model.ExerciseCommand

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

// IIFeeQuoter is a DAML interface
type IIFeeQuoter interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// FeeQuoterPublicFetch executes the FeeQuoter_PublicFetch choice
	FeeQuoterPublicFetch(contractID string, args FeeQuoterPublicFetch) *model.ExerciseCommand

	// FeeQuoterGetTokenPrice executes the FeeQuoter_GetTokenPrice choice
	FeeQuoterGetTokenPrice(contractID string, args FeeQuoterGetTokenPrice) *model.ExerciseCommand

	// FeeQuoterGetDestinationChainGasPrice executes the FeeQuoter_GetDestinationChainGasPrice choice
	FeeQuoterGetDestinationChainGasPrice(contractID string, args FeeQuoterGetDestinationChainGasPrice) *model.ExerciseCommand

	// FeeQuoterGetTokenTransferFee executes the FeeQuoter_GetTokenTransferFee choice
	FeeQuoterGetTokenTransferFee(contractID string, args FeeQuoterGetTokenTransferFee) *model.ExerciseCommand

	// FeeQuoterQuoteGasForExec executes the FeeQuoter_QuoteGasForExec choice
	FeeQuoterQuoteGasForExec(contractID string, args FeeQuoterQuoteGasForExec) *model.ExerciseCommand

	// FeeQuoterGetFeeTokens executes the FeeQuoter_GetFeeTokens choice
	FeeQuoterGetFeeTokens(contractID string, args FeeQuoterGetFeeTokens) *model.ExerciseCommand
}

// IIGlobalConfig is a DAML interface
type IIGlobalConfig interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// GlobalConfigPublicFetch executes the GlobalConfig_PublicFetch choice
	GlobalConfigPublicFetch(contractID string, args GlobalConfigPublicFetch) *model.ExerciseCommand

	// GlobalConfigGetDestChainConfig executes the GlobalConfig_GetDestChainConfig choice
	GlobalConfigGetDestChainConfig(contractID string, args GlobalConfigGetDestChainConfig) *model.ExerciseCommand

	// GlobalConfigGetSourceChainConfig executes the GlobalConfig_GetSourceChainConfig choice
	GlobalConfigGetSourceChainConfig(contractID string, args GlobalConfigGetSourceChainConfig) *model.ExerciseCommand
}

// IIRMNRemote is a DAML interface
type IIRMNRemote interface {

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

// IISendingMessage is a DAML interface
type IISendingMessage interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// SendingMessageFeeTokenAmount executes the SendingMessage_FeeTokenAmount choice
	SendingMessageFeeTokenAmount(contractID string, args SendingMessageFeeTokenAmount) *model.ExerciseCommand

	// SendingMessageAddCCVFee executes the SendingMessage_AddCCVFee choice
	SendingMessageAddCCVFee(contractID string, args SendingMessageAddCCVFee) *model.ExerciseCommand

	// SendingMessageAddVerifierData executes the SendingMessage_AddVerifierData choice
	SendingMessageAddVerifierData(contractID string, args SendingMessageAddVerifierData) *model.ExerciseCommand

	// SendingMessageAddExecutorFee executes the SendingMessage_AddExecutorFee choice
	SendingMessageAddExecutorFee(contractID string, args SendingMessageAddExecutorFee) *model.ExerciseCommand
}

// IITokenAdminRegistry is a DAML interface
type IITokenAdminRegistry interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenAdminRegistryPublicFetch executes the TokenAdminRegistry_PublicFetch choice
	TokenAdminRegistryPublicFetch(contractID string, args TokenAdminRegistryPublicFetch) *model.ExerciseCommand

	// TokenAdminRegistryFetchTokenConfig executes the TokenAdminRegistry_FetchTokenConfig choice
	TokenAdminRegistryFetchTokenConfig(contractID string, args TokenAdminRegistryFetchTokenConfig) *model.ExerciseCommand

	// TokenAdminRegistrySetPool executes the TokenAdminRegistry_SetPool choice
	TokenAdminRegistrySetPool(contractID string, args TokenAdminRegistrySetPool) *model.ExerciseCommand

	// TokenAdminRegistrySetTransferFactory executes the TokenAdminRegistry_SetTransferFactory choice
	TokenAdminRegistrySetTransferFactory(contractID string, args TokenAdminRegistrySetTransferFactory) *model.ExerciseCommand

	// TokenAdminRegistrySetBurnMintFactory executes the TokenAdminRegistry_SetBurnMintFactory choice
	TokenAdminRegistrySetBurnMintFactory(contractID string, args TokenAdminRegistrySetBurnMintFactory) *model.ExerciseCommand

	// TokenAdminRegistryProposeAdministrator executes the TokenAdminRegistry_ProposeAdministrator choice
	TokenAdminRegistryProposeAdministrator(contractID string, args TokenAdminRegistryProposeAdministrator) *model.ExerciseCommand

	// TokenAdminRegistryAcceptAdminRole executes the TokenAdminRegistry_AcceptAdminRole choice
	TokenAdminRegistryAcceptAdminRole(contractID string, args TokenAdminRegistryAcceptAdminRole) *model.ExerciseCommand

	// TokenAdminRegistryTransferAdminRole executes the TokenAdminRegistry_TransferAdminRole choice
	TokenAdminRegistryTransferAdminRole(contractID string, args TokenAdminRegistryTransferAdminRole) *model.ExerciseCommand

	// TokenAdminRegistryIsAdministrator executes the TokenAdminRegistry_IsAdministrator choice
	TokenAdminRegistryIsAdministrator(contractID string, args TokenAdminRegistryIsAdministrator) *model.ExerciseCommand

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

// IITokenConfig is a DAML interface
type IITokenConfig interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenConfigPublicFetch executes the TokenConfig_PublicFetch choice
	TokenConfigPublicFetch(contractID string, args TokenConfigPublicFetch) *model.ExerciseCommand

	// TokenConfigAssertConfiguredTransferFactory executes the TokenConfig_AssertConfiguredTransferFactory choice
	TokenConfigAssertConfiguredTransferFactory(contractID string, args TokenConfigAssertConfiguredTransferFactory) *model.ExerciseCommand

	// TokenConfigAssertConfiguredBurnMintFactory executes the TokenConfig_AssertConfiguredBurnMintFactory choice
	TokenConfigAssertConfiguredBurnMintFactory(contractID string, args TokenConfigAssertConfiguredBurnMintFactory) *model.ExerciseCommand
}

// IITokenReceiveTicket is a DAML interface
type IITokenReceiveTicket interface {

	// Consume executes the Consume choice
	Consume(contractID string, args Consume) *model.ExerciseCommand

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
}

const (
	FeeQuoterContextKey          = types.TEXT("fee-quoter")
	RmnRemoteContextKey          = types.TEXT("rmn-remote")
	GlobalConfigContextKey       = types.TEXT("global-config")
	TokenConfigContextKey        = types.TEXT("token-config")
	TokenAdminRegistryContextKey = types.TEXT("token-admin-registry")
)

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

// DestChainConfig is a Record type
type DestChainConfig struct {
	IsEnabled                 types.BOOL                                 `json:"isEnabled"`
	AddressBytesLength        types.INT64                                `json:"addressBytesLength"`
	TokenReceiverAllowed      types.BOOL                                 `json:"tokenReceiverAllowed"`
	BaseExecutionGasCost      types.INT64                                `json:"baseExecutionGasCost"`
	OffRampAddress            types.TEXT                                 `json:"offRampAddress" hex:"bytes"`
	DefaultExecutor           *chainlinkapi.RawInstanceAddress           `json:"defaultExecutor" hex:"optional"`
	LaneMandatedCCVs          []chainlinkapi.RawInstanceAddress          `json:"laneMandatedCCVs"`
	DefaultCCVs               []chainlinkapi.RawInstanceAddress          `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC                              `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC                              `json:"tokenNetworkFeeUSDCents"`
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts DestChainConfig to a map for DAML arguments
func (t DestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["addressBytesLength"] = int64(t.AddressBytesLength)

	m["tokenReceiverAllowed"] = bool(t.TokenReceiverAllowed)

	m["baseExecutionGasCost"] = int64(t.BaseExecutionGasCost)

	m["offRampAddress"] = string(t.OffRampAddress)

	if t.DefaultExecutor != nil {
		m["defaultExecutor"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.DefaultExecutor),
		}
	} else {
		m["defaultExecutor"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["defaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["messageNetworkFeeUSDCents"] = t.MessageNetworkFeeUSDCents

	m["tokenNetworkFeeUSDCents"] = t.TokenNetworkFeeUSDCents

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// ExecutingMessageCancelExecute is a Record type
type ExecutingMessageCancelExecute struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts ExecutingMessageCancelExecute to a map for DAML arguments
func (t ExecutingMessageCancelExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutingMessageCancelExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageCancelExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageCancelExecute to hex string (Canton MCMS format)
func (t ExecutingMessageCancelExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageCancelExecute from hex string (Canton MCMS format)
func (t *ExecutingMessageCancelExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterView is a Record type
type FeeQuoterView struct {
	CcipOwner  types.PARTY                                `json:"ccipOwner"`
	InstanceId types.TEXT                                 `json:"instanceId"`
	Context    splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts FeeQuoterView to a map for DAML arguments
func (t FeeQuoterView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t FeeQuoterView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterView to hex string (Canton MCMS format)
func (t FeeQuoterView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterView from hex string (Canton MCMS format)
func (t *FeeQuoterView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterGetDestinationChainGasPrice is a Record type
type FeeQuoterGetDestinationChainGasPrice struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterGetDestinationChainGasPrice to a map for DAML arguments
func (t FeeQuoterGetDestinationChainGasPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterGetDestinationChainGasPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterGetDestinationChainGasPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterGetDestinationChainGasPrice to hex string (Canton MCMS format)
func (t FeeQuoterGetDestinationChainGasPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterGetDestinationChainGasPrice from hex string (Canton MCMS format)
func (t *FeeQuoterGetDestinationChainGasPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterGetFeeTokens is a Record type
type FeeQuoterGetFeeTokens struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterGetFeeTokens to a map for DAML arguments
func (t FeeQuoterGetFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterGetFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterGetFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterGetFeeTokens to hex string (Canton MCMS format)
func (t FeeQuoterGetFeeTokens) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterGetFeeTokens from hex string (Canton MCMS format)
func (t *FeeQuoterGetFeeTokens) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterGetTokenPrice is a Record type
type FeeQuoterGetTokenPrice struct {
	Token   splice_api_token_holding_v1.InstrumentId   `json:"token"`
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterGetTokenPrice to a map for DAML arguments
func (t FeeQuoterGetTokenPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["token"] = model.NestedToDAMLValue(t.Token)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterGetTokenPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterGetTokenPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterGetTokenPrice to hex string (Canton MCMS format)
func (t FeeQuoterGetTokenPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterGetTokenPrice from hex string (Canton MCMS format)
func (t *FeeQuoterGetTokenPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterGetTokenTransferFee is a Record type
type FeeQuoterGetTokenTransferFee struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Token             splice_api_token_holding_v1.InstrumentId   `json:"token"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterGetTokenTransferFee to a map for DAML arguments
func (t FeeQuoterGetTokenTransferFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["token"] = model.NestedToDAMLValue(t.Token)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterGetTokenTransferFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterGetTokenTransferFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterGetTokenTransferFee to hex string (Canton MCMS format)
func (t FeeQuoterGetTokenTransferFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterGetTokenTransferFee from hex string (Canton MCMS format)
func (t *FeeQuoterGetTokenTransferFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterPublicFetch is a Record type
type FeeQuoterPublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress            `json:"expectedAddress"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterPublicFetch to a map for DAML arguments
func (t FeeQuoterPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterPublicFetch to hex string (Canton MCMS format)
func (t FeeQuoterPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterPublicFetch from hex string (Canton MCMS format)
func (t *FeeQuoterPublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterQuoteGasForExec is a Record type
type FeeQuoterQuoteGasForExec struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	NonCalldataGas    types.INT64                                `json:"nonCalldataGas"`
	CalldataSize      types.INT64                                `json:"calldataSize"`
	FeeToken          splice_api_token_holding_v1.InstrumentId   `json:"feeToken"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts FeeQuoterQuoteGasForExec to a map for DAML arguments
func (t FeeQuoterQuoteGasForExec) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["nonCalldataGas"] = int64(t.NonCalldataGas)

	m["calldataSize"] = int64(t.CalldataSize)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeQuoterQuoteGasForExec) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterQuoteGasForExec) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterQuoteGasForExec to hex string (Canton MCMS format)
func (t FeeQuoterQuoteGasForExec) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterQuoteGasForExec from hex string (Canton MCMS format)
func (t *FeeQuoterQuoteGasForExec) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfigView is a Record type
type GlobalConfigView struct {
	CcipOwner types.PARTY                                `json:"ccipOwner"`
	Context   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GlobalConfigView to a map for DAML arguments
func (t GlobalConfigView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t GlobalConfigView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalConfigView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalConfigView to hex string (Canton MCMS format)
func (t GlobalConfigView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalConfigView from hex string (Canton MCMS format)
func (t *GlobalConfigView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfigGetDestChainConfig is a Record type
type GlobalConfigGetDestChainConfig struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts GlobalConfigGetDestChainConfig to a map for DAML arguments
func (t GlobalConfigGetDestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GlobalConfigGetDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalConfigGetDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalConfigGetDestChainConfig to hex string (Canton MCMS format)
func (t GlobalConfigGetDestChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalConfigGetDestChainConfig from hex string (Canton MCMS format)
func (t *GlobalConfigGetDestChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfigGetSourceChainConfig is a Record type
type GlobalConfigGetSourceChainConfig struct {
	SourceChainSelector types.NUMERIC                              `json:"sourceChainSelector"`
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller              types.PARTY                                `json:"caller"`
}

// ToMap converts GlobalConfigGetSourceChainConfig to a map for DAML arguments
func (t GlobalConfigGetSourceChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GlobalConfigGetSourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalConfigGetSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalConfigGetSourceChainConfig to hex string (Canton MCMS format)
func (t GlobalConfigGetSourceChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalConfigGetSourceChainConfig from hex string (Canton MCMS format)
func (t *GlobalConfigGetSourceChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfigPublicFetch is a Record type
type GlobalConfigPublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress            `json:"expectedAddress"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts GlobalConfigPublicFetch to a map for DAML arguments
func (t GlobalConfigPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GlobalConfigPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalConfigPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalConfigPublicFetch to hex string (Canton MCMS format)
func (t GlobalConfigPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalConfigPublicFetch from hex string (Canton MCMS format)
func (t *GlobalConfigPublicFetch) UnmarshalHex(data string) error {
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
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.ExecutingMessage", "MessageExecutionState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e MessageExecutionState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.APIV2.ExecutingMessage", "MessageExecutionState")
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

// PoolRegistration is a Record type
type PoolRegistration struct {
	PoolOwner      types.PARTY `json:"poolOwner"`
	PoolInstanceId types.TEXT  `json:"poolInstanceId"`
}

// ToMap converts PoolRegistration to a map for DAML arguments
func (t PoolRegistration) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	return m
}

func (t PoolRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PoolRegistration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PoolRegistration to hex string (Canton MCMS format)
func (t PoolRegistration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PoolRegistration from hex string (Canton MCMS format)
func (t *PoolRegistration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministratorResult is a Record type
type ProposeAdministratorResult struct {
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	Created               types.BOOL                                 `json:"created"`
	Index                 types.INT64                                `json:"index"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts ProposeAdministratorResult to a map for DAML arguments
func (t ProposeAdministratorResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["created"] = bool(t.Created)

	m["index"] = int64(t.Index)

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t ProposeAdministratorResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProposeAdministratorResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProposeAdministratorResult to hex string (Canton MCMS format)
func (t ProposeAdministratorResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdministratorResult from hex string (Canton MCMS format)
func (t *ProposeAdministratorResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExecResult is a Record type
type QuoteGasForExecResult struct {
	TotalGas          types.INT64                                `json:"totalGas"`
	GasCostUSDCents   types.NUMERIC                              `json:"gasCostUSDCents"`
	FeeTokenPrice     types.NUMERIC                              `json:"feeTokenPrice"`
	PremiumMultiplier types.NUMERIC                              `json:"premiumMultiplier"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts QuoteGasForExecResult to a map for DAML arguments
func (t QuoteGasForExecResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["totalGas"] = int64(t.TotalGas)

	m["gasCostUSDCents"] = t.GasCostUSDCents

	m["feeTokenPrice"] = t.FeeTokenPrice

	m["premiumMultiplier"] = t.PremiumMultiplier

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t QuoteGasForExecResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *QuoteGasForExecResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes QuoteGasForExecResult to hex string (Canton MCMS format)
func (t QuoteGasForExecResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExecResult from hex string (Canton MCMS format)
func (t *QuoteGasForExecResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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
	ExpectedAddress chainlinkapi.RawInstanceAddress            `json:"expectedAddress"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts RMNRemotePublicFetch to a map for DAML arguments
func (t RMNRemotePublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// SendingMessageFeeTokenAmount is a Record type
type SendingMessageFeeTokenAmount struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts SendingMessageFeeTokenAmount to a map for DAML arguments
func (t SendingMessageFeeTokenAmount) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SendingMessageFeeTokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageFeeTokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageFeeTokenAmount to hex string (Canton MCMS format)
func (t SendingMessageFeeTokenAmount) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageFeeTokenAmount from hex string (Canton MCMS format)
func (t *SendingMessageFeeTokenAmount) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        types.BOOL                                 `json:"isEnabled"`
	OnRampAddresses  []types.TEXT                               `json:"onRampAddresses" hex:"[]bytes"`
	DefaultCCVs      []chainlinkapi.RawInstanceAddress          `json:"defaultCCVs"`
	LaneMandatedCCVs []chainlinkapi.RawInstanceAddress          `json:"laneMandatedCCVs"`
	Context          splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["laneMandatedCCVs"] = func() []any {
		res := make([]any, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// TimestampedPrice is a Record type
type TimestampedPrice struct {
	Price     types.NUMERIC                              `json:"price"`
	Timestamp types.TIMESTAMP                            `json:"timestamp"`
	Context   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts TimestampedPrice to a map for DAML arguments
func (t TimestampedPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["price"] = t.Price

	m["timestamp"] = t.Timestamp

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t TimestampedPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TimestampedPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TimestampedPrice to hex string (Canton MCMS format)
func (t TimestampedPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TimestampedPrice from hex string (Canton MCMS format)
func (t *TimestampedPrice) UnmarshalHex(data string) error {
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

// TokenAdminRegistryAcceptAdminRole is a Record type
type TokenAdminRegistryAcceptAdminRole struct {
	TokenConfigCid types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryAcceptAdminRole to a map for DAML arguments
func (t TokenAdminRegistryAcceptAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAcceptAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAcceptAdminRole to hex string (Canton MCMS format)
func (t TokenAdminRegistryAcceptAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAcceptAdminRole from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSend is a Record type
type TokenAdminRegistryAddTokenSend struct {
	TokenConfigCid    types.CONTRACT_ID                          `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                                 `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Amount            types.TEXT                                 `json:"amount"`
	DestTokenAddress  types.TEXT                                 `json:"destTokenAddress"`
	ExtraData         types.TEXT                                 `json:"extraData"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	TokenConfigCid    types.CONTRACT_ID                          `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                                 `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC                              `json:"feeUSDCents"`
	DestGasOverhead   types.INT64                                `json:"destGasOverhead"`
	DestBytesOverhead types.INT64                                `json:"destBytesOverhead"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	TokenConfigCid        types.CONTRACT_ID                          `json:"tokenConfigCid"`
	TokenReceiveTicketCid types.CONTRACT_ID                          `json:"tokenReceiveTicketCid"`
	InstrumentId          splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	PoolInstanceId        types.TEXT                                 `json:"poolInstanceId"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller                types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryConsumeReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryConsumeReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	TokenConfigCid types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryFetchTokenConfig to a map for DAML arguments
func (t TokenAdminRegistryFetchTokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// TokenAdminRegistryIsAdministrator is a Record type
type TokenAdminRegistryIsAdministrator struct {
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	TokenConfigCid *types.CONTRACT_ID                         `json:"tokenConfigCid" hex:"optional"`
	Administrator  types.PARTY                                `json:"administrator"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryIsAdministrator to a map for DAML arguments
func (t TokenAdminRegistryIsAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TokenConfigCid != nil {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenConfigCid),
		}
	} else {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["administrator"] = t.Administrator.ToMap()

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryIsAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryIsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryIsAdministrator to hex string (Canton MCMS format)
func (t TokenAdminRegistryIsAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryIsAdministrator from hex string (Canton MCMS format)
func (t *TokenAdminRegistryIsAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryProposeAdministrator is a Record type
type TokenAdminRegistryProposeAdministrator struct {
	TokenConfigCid *types.CONTRACT_ID                         `json:"tokenConfigCid" hex:"optional"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	NewAdmin       types.PARTY                                `json:"newAdmin"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryProposeAdministrator to a map for DAML arguments
func (t TokenAdminRegistryProposeAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	if t.TokenConfigCid != nil {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenConfigCid),
		}
	} else {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryProposeAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryProposeAdministrator to hex string (Canton MCMS format)
func (t TokenAdminRegistryProposeAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryProposeAdministrator from hex string (Canton MCMS format)
func (t *TokenAdminRegistryProposeAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryPublicFetch is a Record type
type TokenAdminRegistryPublicFetch struct {
	ExpectedAddress chainlinkapi.RawInstanceAddress            `json:"expectedAddress"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryPublicFetch to a map for DAML arguments
func (t TokenAdminRegistryPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// TokenAdminRegistrySetBurnMintFactory is a Record type
type TokenAdminRegistrySetBurnMintFactory struct {
	TokenConfigCid  types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	BurnMintFactory *types.CONTRACT_ID                         `json:"burnMintFactory" hex:"optional"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetBurnMintFactory to a map for DAML arguments
func (t TokenAdminRegistrySetBurnMintFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.BurnMintFactory != nil {
		m["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.BurnMintFactory),
		}
	} else {
		m["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetBurnMintFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetBurnMintFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetBurnMintFactory to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetBurnMintFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetBurnMintFactory from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetBurnMintFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetInboundPoolCCVs is a Record type
type TokenAdminRegistrySetInboundPoolCCVs struct {
	TokenConfigCid      types.CONTRACT_ID                          `json:"tokenConfigCid"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                                 `json:"poolInstanceId"`
	PoolCCVs            []chainlinkapi.RawInstanceAddress          `json:"poolCCVs"`
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller              types.PARTY                                `json:"caller"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	TokenConfigCid    types.CONTRACT_ID                          `json:"tokenConfigCid"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                                 `json:"poolInstanceId"`
	PoolCCVs          []chainlinkapi.RawInstanceAddress          `json:"poolCCVs"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	TokenConfigCid types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	TokenPool      *PoolRegistration                          `json:"tokenPool" hex:"optional"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetPool to a map for DAML arguments
func (t TokenAdminRegistrySetPool) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetPool to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetPool from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetTransferFactory is a Record type
type TokenAdminRegistrySetTransferFactory struct {
	TokenConfigCid  types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	TransferFactory *types.CONTRACT_ID                         `json:"transferFactory" hex:"optional"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetTransferFactory to a map for DAML arguments
func (t TokenAdminRegistrySetTransferFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TransferFactory != nil {
		m["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TransferFactory),
		}
	} else {
		m["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetTransferFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetTransferFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetTransferFactory to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetTransferFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetTransferFactory from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetTransferFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryTransferAdminRole is a Record type
type TokenAdminRegistryTransferAdminRole struct {
	TokenConfigCid types.CONTRACT_ID                          `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	NewAdmin       types.PARTY                                `json:"newAdmin"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller         types.PARTY                                `json:"caller"`
}

// ToMap converts TokenAdminRegistryTransferAdminRole to a map for DAML arguments
func (t TokenAdminRegistryTransferAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryTransferAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryTransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryTransferAdminRole to hex string (Canton MCMS format)
func (t TokenAdminRegistryTransferAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryTransferAdminRole from hex string (Canton MCMS format)
func (t *TokenAdminRegistryTransferAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfigView is a Record type
type TokenConfigView struct {
	CcipOwner    types.PARTY                                `json:"ccipOwner"`
	InstanceId   types.TEXT                                 `json:"instanceId"`
	InstrumentId splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	TokenPool    *PoolRegistration                          `json:"tokenPool" hex:"optional"`
	Context      splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts TokenConfigView to a map for DAML arguments
func (t TokenConfigView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	SuppliedFactory types.CONTRACT_ID                          `json:"suppliedFactory"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenConfigAssertConfiguredBurnMintFactory to a map for DAML arguments
func (t TokenConfigAssertConfiguredBurnMintFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["suppliedFactory"] = model.NestedToDAMLValue(t.SuppliedFactory)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	SuppliedFactory types.CONTRACT_ID                          `json:"suppliedFactory"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenConfigAssertConfiguredTransferFactory to a map for DAML arguments
func (t TokenConfigAssertConfiguredTransferFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["suppliedFactory"] = model.NestedToDAMLValue(t.SuppliedFactory)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	ExpectedAddress chainlinkapi.RawInstanceAddress            `json:"expectedAddress"`
	Context         splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller          types.PARTY                                `json:"caller"`
}

// ToMap converts TokenConfigPublicFetch to a map for DAML arguments
func (t TokenConfigPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAddress"] = model.NestedToDAMLValue(t.ExpectedAddress)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	CcipOwner types.PARTY                                `json:"ccipOwner"`
	PoolOwner types.PARTY                                `json:"poolOwner"`
	CcvOwners []types.PARTY                              `json:"ccvOwners"`
	Context   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// IIExecutingMessageInterfaceID returns the interface ID for the IIExecutingMessage interface using the package name
func IIExecutingMessageInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.ExecutingMessage", "IExecutingMessage")
}

// IIExecutingMessageInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IIExecutingMessageInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.ExecutingMessage", "IExecutingMessage")
}

// IIFeeQuoterInterfaceID returns the interface ID for the IIFeeQuoter interface using the package name
func IIFeeQuoterInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.FeeQuoter", "IFeeQuoter")
}

// IIFeeQuoterInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IIFeeQuoterInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.FeeQuoter", "IFeeQuoter")
}

// IIGlobalConfigInterfaceID returns the interface ID for the IIGlobalConfig interface using the package name
func IIGlobalConfigInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.GlobalConfig", "IGlobalConfig")
}

// IIGlobalConfigInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IIGlobalConfigInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.GlobalConfig", "IGlobalConfig")
}

// IIRMNRemoteInterfaceID returns the interface ID for the IIRMNRemote interface using the package name
func IIRMNRemoteInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.RMNRemote", "IRMNRemote")
}

// IIRMNRemoteInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IIRMNRemoteInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.RMNRemote", "IRMNRemote")
}

// IISendingMessageInterfaceID returns the interface ID for the IISendingMessage interface using the package name
func IISendingMessageInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.SendingMessage", "ISendingMessage")
}

// IISendingMessageInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IISendingMessageInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.SendingMessage", "ISendingMessage")
}

// IITokenAdminRegistryInterfaceID returns the interface ID for the IITokenAdminRegistry interface using the package name
func IITokenAdminRegistryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.TokenAdminRegistry", "ITokenAdminRegistry")
}

// IITokenAdminRegistryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IITokenAdminRegistryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.TokenAdminRegistry", "ITokenAdminRegistry")
}

// IITokenConfigInterfaceID returns the interface ID for the IITokenConfig interface using the package name
func IITokenConfigInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.TokenAdminRegistry", "ITokenConfig")
}

// IITokenConfigInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IITokenConfigInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.TokenAdminRegistry", "ITokenConfig")
}

// IITokenReceiveTicketInterfaceID returns the interface ID for the IITokenReceiveTicket interface using the package name
func IITokenReceiveTicketInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.APIV2.ExecutingMessage", "ITokenReceiveTicket")
}

// IITokenReceiveTicketInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IITokenReceiveTicketInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.APIV2.ExecutingMessage", "ITokenReceiveTicket")
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
