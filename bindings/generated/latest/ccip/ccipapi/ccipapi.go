package ccipapi

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
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
	PackageID   = "33db14c5f121ea80a05f47663c03a14bad793a97051afc7d80d1ce16028ba779"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IRMNRemote is a DAML interface
type IRMNRemote interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

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

	// SendingMessageBuildMessage executes the SendingMessage_BuildMessage choice
	SendingMessageBuildMessage(contractID string, args SendingMessageBuildMessage) *model.ExerciseCommand
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

// SendingMessageBuildMessage is a Record type
type SendingMessageBuildMessage struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts SendingMessageBuildMessage to a map for DAML arguments
func (t SendingMessageBuildMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SendingMessageBuildMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SendingMessageBuildMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SendingMessageBuildMessage to hex string (Canton MCMS format)
func (t SendingMessageBuildMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SendingMessageBuildMessage from hex string (Canton MCMS format)
func (t *SendingMessageBuildMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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
