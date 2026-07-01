package core

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	api "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
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
	PackageName = "ccip-core"
	PackageID   = "78669a275f92a3daad67ef946e890e7e1cbc32e302e659284277b7d606865156"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	TokenConfigKey           = types.TEXT("token-config")
	TokenAdminRegistryKey    = types.TEXT("token-admin-registry")
	GlobalCurseSubject       = types.TEXT("01000000000000000000000000000001")
	UsdPerUsdCent            = types.NUMERIC("100000000.")
	PremiumIdentity          = types.NUMERIC("10000000000.")
	MinTokenDecimals         = types.INT64(0)
	MaxTokenDecimals         = types.INT64(37)
	E10PerPercent            = types.NUMERIC("100000000.")
	BaseNumeric              = types.NUMERIC("100000000.")
	BaseInt                  = types.INT64(100000000)
	MaxUint256DecimalText    = types.TEXT("115792089237316195423570985008687907853269984665640564039457584007913129639935")
	MaxNumeric0DecimalText   = types.TEXT("99999999999999999999999999999999999999")
	RmnRemoteKey             = types.TEXT("rmn-remote")
	FeeQuoterKey             = types.TEXT("fee-quoter")
	WaitForFinalityFlag      = types.TEXT("00000000")
	MinBlockDepth            = types.INT64(1)
	MaxBlockDepth            = types.INT64(65535)
	FinalityConfigByteLength = types.INT64(4)
	MaxNumeric0IntegerText   = types.TEXT("99999999999999999999999999999999999999")
	GlobalConfigKey          = types.TEXT("global-config")
	MaxCCVsPerMessage        = types.INT64(255)
	RateLimiterKey           = types.TEXT("rate-limiter")
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

// AcceptAdminParams is a Record type
type AcceptAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// ToMap converts AcceptAdminParams to a map for DAML arguments
func (t AcceptAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	return m
}

func (t AcceptAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptAdminParams to hex string (Canton MCMS format)
func (t AcceptAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminParams from hex string (Canton MCMS format)
func (t *AcceptAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptAdminRole is a Record type
type AcceptAdminRole struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts AcceptAdminRole to a map for DAML arguments
func (t AcceptAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t AcceptAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptAdminRole to hex string (Canton MCMS format)
func (t AcceptAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminRole from hex string (Canton MCMS format)
func (t *AcceptAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptAdminRoleMCMSParams is AcceptAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type AcceptAdminRoleMCMSParams struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes AcceptAdminRoleMCMSParams to hex string for MCMS operationData.
func (t AcceptAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptAdminRoleMCMSParams from hex string.
func (t *AcceptAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCCVFee is a Record type
type AddCCVFee struct {
	CcvInstanceId     types.TEXT                                 `json:"ccvInstanceId"`
	FeeUSDCents       types.NUMERIC                              `json:"feeUSDCents"`
	DestGasLimit      types.INT64                                `json:"destGasLimit"`
	DestBytesOverhead types.INT64                                `json:"destBytesOverhead"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts AddCCVFee to a map for DAML arguments
func (t AddCCVFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvInstanceId"] = string(t.CcvInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	CcvInstanceId     types.TEXT                                 `json:"ccvInstanceId"`
	FeeUSDCents       types.NUMERIC                              `json:"feeUSDCents"`
	DestGasLimit      types.INT64                                `json:"destGasLimit"`
	DestBytesOverhead types.INT64                                `json:"destBytesOverhead"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	VersionTag    types.TEXT  `json:"versionTag" hex:"bytes"`
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
	VersionTag    types.TEXT `json:"versionTag" hex:"bytes"`
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

// AddCustomObservers is a Record type
type AddCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObservers to a map for DAML arguments
func (t AddCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObservers to hex string (Canton MCMS format)
func (t AddCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObservers from hex string (Canton MCMS format)
func (t *AddCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCustomObserversParams is a Record type
type AddCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObserversParams to a map for DAML arguments
func (t AddCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObserversParams to hex string (Canton MCMS format)
func (t AddCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObserversParams from hex string (Canton MCMS format)
func (t *AddCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddExecutorFee is a Record type
type AddExecutorFee struct {
	ExecutorInstanceId types.TEXT                                 `json:"executorInstanceId"`
	ExecutorArgs       types.TEXT                                 `json:"executorArgs"`
	FeeUSDCents        types.NUMERIC                              `json:"feeUSDCents"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller             types.PARTY                                `json:"caller"`
}

// ToMap converts AddExecutorFee to a map for DAML arguments
func (t AddExecutorFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["executorInstanceId"] = string(t.ExecutorInstanceId)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["feeUSDCents"] = t.FeeUSDCents

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	ExecutorInstanceId types.TEXT                                 `json:"executorInstanceId"`
	ExecutorArgs       types.TEXT                                 `json:"executorArgs"`
	FeeUSDCents        types.NUMERIC                              `json:"feeUSDCents"`
	Context            splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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

// AddPriceUpdaters is a Record type
type AddPriceUpdaters struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddPriceUpdaters to a map for DAML arguments
func (t AddPriceUpdaters) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddPriceUpdaters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddPriceUpdaters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddPriceUpdaters to hex string (Canton MCMS format)
func (t AddPriceUpdaters) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddPriceUpdaters from hex string (Canton MCMS format)
func (t *AddPriceUpdaters) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddTokenSend is a Record type
type AddTokenSend struct {
	PoolInstanceId   types.TEXT                                 `json:"poolInstanceId"`
	PoolOwner        types.PARTY                                `json:"poolOwner"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Amount           types.TEXT                                 `json:"amount"`
	DestTokenAddress types.TEXT                                 `json:"destTokenAddress"`
	ExtraData        types.TEXT                                 `json:"extraData"`
	Context          splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts AddTokenSend to a map for DAML arguments
func (t AddTokenSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = string(t.Amount)

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["extraData"] = string(t.ExtraData)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	PoolInstanceId    types.TEXT                                 `json:"poolInstanceId"`
	PoolOwner         types.PARTY                                `json:"poolOwner"`
	FeeUSDCents       types.NUMERIC                              `json:"feeUSDCents"`
	DestGasOverhead   types.INT64                                `json:"destGasOverhead"`
	DestBytesOverhead types.INT64                                `json:"destBytesOverhead"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts AddTokenSendFee to a map for DAML arguments
func (t AddTokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	CcvInstanceId        types.TEXT                                 `json:"ccvInstanceId"`
	VersionTag           types.TEXT                                 `json:"versionTag" hex:"bytes"`
	VerifierBlob         types.TEXT                                 `json:"verifierBlob"`
	MessageSentObservers []types.PARTY                              `json:"messageSentObservers"`
	Context              splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller               types.PARTY                                `json:"caller"`
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

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	CcvInstanceId        types.TEXT                                 `json:"ccvInstanceId"`
	VersionTag           types.TEXT                                 `json:"versionTag" hex:"bytes"`
	VerifierBlob         types.TEXT                                 `json:"verifierBlob"`
	MessageSentObservers []types.PARTY                              `json:"messageSentObservers"`
	Context              splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
			res = append(res, model.NestedToDAMLValue(e))
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

// ApplyDestChainConfigUpdatesParams is a Record type
type ApplyDestChainConfigUpdatesParams struct {
	DestChainConfigArgs []DestChainConfigArgs `json:"destChainConfigArgs"`
}

// ToMap converts ApplyDestChainConfigUpdatesParams to a map for DAML arguments
func (t ApplyDestChainConfigUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyDestChainConfigUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainConfigUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainConfigUpdatesParams to hex string (Canton MCMS format)
func (t ApplyDestChainConfigUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainConfigUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyDestChainConfigUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyFeeQuoterDestChainConfigUpdates is a Record type
type ApplyFeeQuoterDestChainConfigUpdates struct {
	DestChainConfigArgs []FeeQuoterDestChainConfigArgs `json:"destChainConfigArgs"`
}

// ToMap converts ApplyFeeQuoterDestChainConfigUpdates to a map for DAML arguments
func (t ApplyFeeQuoterDestChainConfigUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyFeeQuoterDestChainConfigUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyFeeQuoterDestChainConfigUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyFeeQuoterDestChainConfigUpdates to hex string (Canton MCMS format)
func (t ApplyFeeQuoterDestChainConfigUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyFeeQuoterDestChainConfigUpdates from hex string (Canton MCMS format)
func (t *ApplyFeeQuoterDestChainConfigUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyFeeQuoterDestChainConfigUpdatesParams is a Record type
type ApplyFeeQuoterDestChainConfigUpdatesParams struct {
	DestChainConfigArgs []FeeQuoterDestChainConfigArgs `json:"destChainConfigArgs"`
}

// ToMap converts ApplyFeeQuoterDestChainConfigUpdatesParams to a map for DAML arguments
func (t ApplyFeeQuoterDestChainConfigUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.DestChainConfigArgs))
		for _, e := range t.DestChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplyFeeQuoterDestChainConfigUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyFeeQuoterDestChainConfigUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyFeeQuoterDestChainConfigUpdatesParams to hex string (Canton MCMS format)
func (t ApplyFeeQuoterDestChainConfigUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyFeeQuoterDestChainConfigUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyFeeQuoterDestChainConfigUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyPriceUpdatersUpdate is a Record type
type ApplyPriceUpdatersUpdate struct {
	AddedPriceUpdaters   []types.PARTY `json:"addedPriceUpdaters"`
	RemovedPriceUpdaters []types.PARTY `json:"removedPriceUpdaters"`
}

// ToMap converts ApplyPriceUpdatersUpdate to a map for DAML arguments
func (t ApplyPriceUpdatersUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["addedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.AddedPriceUpdaters))
		for _, e := range t.AddedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["removedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.RemovedPriceUpdaters))
		for _, e := range t.RemovedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t ApplyPriceUpdatersUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyPriceUpdatersUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyPriceUpdatersUpdate to hex string (Canton MCMS format)
func (t ApplyPriceUpdatersUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyPriceUpdatersUpdate from hex string (Canton MCMS format)
func (t *ApplyPriceUpdatersUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyPriceUpdatersUpdateParams is a Record type
type ApplyPriceUpdatersUpdateParams struct {
	AddedPriceUpdaters   []types.PARTY `json:"addedPriceUpdaters"`
	RemovedPriceUpdaters []types.PARTY `json:"removedPriceUpdaters"`
}

// ToMap converts ApplyPriceUpdatersUpdateParams to a map for DAML arguments
func (t ApplyPriceUpdatersUpdateParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["addedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.AddedPriceUpdaters))
		for _, e := range t.AddedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["removedPriceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.RemovedPriceUpdaters))
		for _, e := range t.RemovedPriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t ApplyPriceUpdatersUpdateParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyPriceUpdatersUpdateParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyPriceUpdatersUpdateParams to hex string (Canton MCMS format)
func (t ApplyPriceUpdatersUpdateParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyPriceUpdatersUpdateParams from hex string (Canton MCMS format)
func (t *ApplyPriceUpdatersUpdateParams) UnmarshalHex(data string) error {
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
			res = append(res, model.NestedToDAMLValue(e))
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

// ApplySourceChainConfigUpdatesParams is a Record type
type ApplySourceChainConfigUpdatesParams struct {
	SourceChainConfigArgs []SourceChainConfigArgs `json:"sourceChainConfigArgs"`
}

// ToMap converts ApplySourceChainConfigUpdatesParams to a map for DAML arguments
func (t ApplySourceChainConfigUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainConfigArgs"] = func() []any {
		res := make([]any, 0, len(t.SourceChainConfigArgs))
		for _, e := range t.SourceChainConfigArgs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ApplySourceChainConfigUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplySourceChainConfigUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplySourceChainConfigUpdatesParams to hex string (Canton MCMS format)
func (t ApplySourceChainConfigUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplySourceChainConfigUpdatesParams from hex string (Canton MCMS format)
func (t *ApplySourceChainConfigUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BuildMessage is a Record type
type BuildMessage struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts BuildMessage to a map for DAML arguments
func (t BuildMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	args["event"] = model.NestedToDAMLValue(t.Event)

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
	args["event"] = model.NestedToDAMLValue(t.Event)

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
			res = append(res, model.NestedToDAMLValue(e))
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
	VersionTag    types.TEXT  `json:"versionTag" hex:"bytes"`
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

	m["rateLimiterCid"] = model.NestedToDAMLValue(t.RateLimiterCid)

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

// ConsumeReceiveTicket is a Record type
type ConsumeReceiveTicket struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	TicketCid      types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	PoolInstanceId types.TEXT                               `json:"poolInstanceId"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts ConsumeReceiveTicket to a map for DAML arguments
func (t ConsumeReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["ticketCid"] = model.NestedToDAMLValue(t.TicketCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ConsumeReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeReceiveTicket to hex string (Canton MCMS format)
func (t ConsumeReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeReceiveTicket from hex string (Canton MCMS format)
func (t *ConsumeReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeReceiveTicketMCMSParams is ConsumeReceiveTicket without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ConsumeReceiveTicketMCMSParams struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	TicketCid      types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	PoolInstanceId types.TEXT                               `json:"poolInstanceId"`
}

// MarshalHex encodes ConsumeReceiveTicketMCMSParams to hex string for MCMS operationData.
func (t ConsumeReceiveTicketMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeReceiveTicketMCMSParams from hex string.
func (t *ConsumeReceiveTicketMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Curse is a Record type
type Curse struct {
	Subject types.TEXT `json:"subject" hex:"bytes"`
}

// ToMap converts Curse to a map for DAML arguments
func (t Curse) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t Curse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Curse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Curse to hex string (Canton MCMS format)
func (t Curse) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Curse from hex string (Canton MCMS format)
func (t *Curse) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseChain is a Record type
type CurseChain struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts CurseChain to a map for DAML arguments
func (t CurseChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t CurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseChain to hex string (Canton MCMS format)
func (t CurseChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseChain from hex string (Canton MCMS format)
func (t *CurseChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseChainParams is a Record type
type CurseChainParams struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts CurseChainParams to a map for DAML arguments
func (t CurseChainParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t CurseChainParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseChainParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseChainParams to hex string (Canton MCMS format)
func (t CurseChainParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseChainParams from hex string (Canton MCMS format)
func (t *CurseChainParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseGlobal is a Record type
type CurseGlobal struct {
}

// ToMap converts CurseGlobal to a map for DAML arguments
func (t CurseGlobal) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseGlobal to hex string (Canton MCMS format)
func (t CurseGlobal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseGlobal from hex string (Canton MCMS format)
func (t *CurseGlobal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseMultiple is a Record type
type CurseMultiple struct {
	Subjects []types.TEXT `json:"subjects" hex:"[]bytes"`
}

// ToMap converts CurseMultiple to a map for DAML arguments
func (t CurseMultiple) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseMultiple to hex string (Canton MCMS format)
func (t CurseMultiple) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseMultiple from hex string (Canton MCMS format)
func (t *CurseMultiple) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseMultipleParams is a Record type
type CurseMultipleParams struct {
	Subjects []types.TEXT `json:"subjects" hex:"[]bytes"`
}

// ToMap converts CurseMultipleParams to a map for DAML arguments
func (t CurseMultipleParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CurseMultipleParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseMultipleParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseMultipleParams to hex string (Canton MCMS format)
func (t CurseMultipleParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseMultipleParams from hex string (Canton MCMS format)
func (t *CurseMultipleParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CurseParams is a Record type
type CurseParams struct {
	Subject types.TEXT `json:"subject" hex:"bytes"`
}

// ToMap converts CurseParams to a map for DAML arguments
func (t CurseParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t CurseParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CurseParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CurseParams to hex string (Canton MCMS format)
func (t CurseParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CurseParams from hex string (Canton MCMS format)
func (t *CurseParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DecodedFinality is a Record type
type DecodedFinality struct {
	Raw       types.TEXT     `json:"raw"`
	Requested FinalityConfig `json:"requested"`
}

// ToMap converts DecodedFinality to a map for DAML arguments
func (t DecodedFinality) ToMap() map[string]any {
	m := make(map[string]any)

	m["raw"] = string(t.Raw)

	m["requested"] = model.NestedToDAMLValue(t.Requested)

	return m
}

func (t DecodedFinality) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DecodedFinality) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DecodedFinality to hex string (Canton MCMS format)
func (t DecodedFinality) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DecodedFinality from hex string (Canton MCMS format)
func (t *DecodedFinality) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DestChainConfig is a Record type
type DestChainConfig struct {
	IsEnabled                 types.BOOL                        `json:"isEnabled"`
	AddressBytesLength        types.INT64                       `json:"addressBytesLength"`
	TokenReceiverAllowed      types.BOOL                        `json:"tokenReceiverAllowed"`
	BaseExecutionGasCost      types.INT64                       `json:"baseExecutionGasCost"`
	OffRampAddress            types.TEXT                        `json:"offRampAddress" hex:"bytes"`
	DefaultExecutor           *chainlinkapi.RawInstanceAddress  `json:"defaultExecutor" hex:"optional"`
	LaneMandatedCCVs          []chainlinkapi.RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []chainlinkapi.RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC                     `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC                     `json:"tokenNetworkFeeUSDCents"`
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
	DestChainSelector         types.NUMERIC                     `json:"destChainSelector"`
	IsEnabled                 types.BOOL                        `json:"isEnabled"`
	AddressBytesLength        types.INT64                       `json:"addressBytesLength"`
	TokenReceiverAllowed      types.BOOL                        `json:"tokenReceiverAllowed"`
	BaseExecutionGasCost      types.INT64                       `json:"baseExecutionGasCost"`
	OffRampAddress            types.TEXT                        `json:"offRampAddress" hex:"bytes"`
	DefaultExecutor           *chainlinkapi.RawInstanceAddress  `json:"defaultExecutor" hex:"optional"`
	LaneMandatedCCVs          []chainlinkapi.RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []chainlinkapi.RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC                     `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC                     `json:"tokenNetworkFeeUSDCents"`
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
	OffRamp            chainlinkapi.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts ExecutingMessageDeps to a map for DAML arguments
func (t ExecutingMessageDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["offRamp"] = model.NestedToDAMLValue(t.OffRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

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
	CcipOwner               types.PARTY                       `json:"ccipOwner"`
	Message                 MessageV1                         `json:"message"`
	MessageId               types.TEXT                        `json:"messageId"`
	Receiver                types.PARTY                       `json:"receiver"`
	TokenReceiver           *types.PARTY                      `json:"tokenReceiver" hex:"optional"`
	Executor                types.PARTY                       `json:"executor"`
	ObservingParties        []types.PARTY                     `json:"observingParties"`
	CcvVerifications        []CCVVerification                 `json:"ccvVerifications"`
	CcvOwners               []types.PARTY                     `json:"ccvOwners"`
	RequiredCCVs            []chainlinkapi.RawInstanceAddress `json:"requiredCCVs"`
	OptionalCCVs            []chainlinkapi.RawInstanceAddress `json:"optionalCCVs"`
	OptionalCCVThreshold    types.INT64                       `json:"optionalCCVThreshold"`
	ReceiverFinalityConfig  FinalityConfig                    `json:"receiverFinalityConfig"`
	SourceDefaultCCVs       []chainlinkapi.RawInstanceAddress `json:"sourceDefaultCCVs"`
	InboundPoolVerification *InboundPoolVerification          `json:"inboundPoolVerification" hex:"optional"`
	Deps                    ExecutingMessageDeps              `json:"deps"`
	State                   ExecutingMessageState             `json:"state"`
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
	args["message"] = model.NestedToDAMLValue(t.Message)

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
			"value": nil,
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
			res = append(res, model.NestedToDAMLValue(e))
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVThreshold"] = int64(t.OptionalCCVThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.SourceDefaultCCVs))
		for _, e := range t.SourceDefaultCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.InboundPoolVerification != nil {
		args["inboundPoolVerification"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.InboundPoolVerification),
		}
	} else {
		args["inboundPoolVerification"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	if t.State != "" {
		args["state"] = model.NestedToDAMLValue(t.State)
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
	args["message"] = model.NestedToDAMLValue(t.Message)

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
			"value": nil,
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
			res = append(res, model.NestedToDAMLValue(e))
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVThreshold"] = int64(t.OptionalCCVThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.SourceDefaultCCVs))
		for _, e := range t.SourceDefaultCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.InboundPoolVerification != nil {
		args["inboundPoolVerification"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.InboundPoolVerification),
		}
	} else {
		args["inboundPoolVerification"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	if t.State != "" {
		args["state"] = model.NestedToDAMLValue(t.State)
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
	args["event"] = model.NestedToDAMLValue(t.Event)

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
	args["event"] = model.NestedToDAMLValue(t.Event)

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

	m["state"] = model.NestedToDAMLValue(t.State)

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

// FeeQuoter is a Template type
type FeeQuoter struct {
	InstanceId                       types.TEXT                                 `json:"instanceId"`
	Owner                            types.PARTY                                `json:"owner"`
	FeeTokens                        types.SET                                  `json:"feeTokens"`
	DestChainConfigs                 map[types.NUMERIC]FeeQuoterDestChainConfig `json:"destChainConfigs"`
	TokenTransferFeeConfigs          map[types.NUMERIC]types.GENMAP             `json:"tokenTransferFeeConfigs"`
	UsdPerUnitGasByDestChainSelector map[types.NUMERIC]TimestampedPrice         `json:"usdPerUnitGasByDestChainSelector"`
	UsdPerToken                      types.GENMAP                               `json:"usdPerToken"`
	LinkTokenInstrumentId            splice_api_token_holding_v1.InstrumentId   `json:"linkTokenInstrumentId"`
	PriceUpdaters                    []types.PARTY                              `json:"priceUpdaters"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FeeQuoter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FeeQuoter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FeeQuoter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeTokens"] = model.NestedToDAMLValue(t.FeeTokens)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["linkTokenInstrumentId"] = model.NestedToDAMLValue(t.LinkTokenInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["priceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.PriceUpdaters))
		for _, e := range t.PriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FeeQuoter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["feeTokens"] = model.NestedToDAMLValue(t.FeeTokens)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() any {
		if t.DestChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenTransferFeeConfigs"] = func() any {
		if t.TokenTransferFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenTransferFeeConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerUnitGasByDestChainSelector"] = func() any {
		if t.UsdPerUnitGasByDestChainSelector == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerUnitGasByDestChainSelector}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usdPerToken"] = func() any {
		if t.UsdPerToken == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsdPerToken}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["linkTokenInstrumentId"] = model.NestedToDAMLValue(t.LinkTokenInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["priceUpdaters"] = func() []any {
		res := make([]any, 0, len(t.PriceUpdaters))
		for _, e := range t.PriceUpdaters {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FeeQuoter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoter to hex string (Canton MCMS format)
func (t FeeQuoter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoter from hex string (Canton MCMS format)
func (t *FeeQuoter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FeeQuoter

// QuoteGasForExec exercises the QuoteGasForExec choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) QuoteGasForExec(contractID string, args QuoteGasForExec) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "QuoteGasForExec",
		Arguments:  argsToMap(args),
	}
}

// QuoteGasForExecWithPackageID exercises the QuoteGasForExec choice using the provided package ID instead of package name
func (t FeeQuoter) QuoteGasForExecWithPackageID(contractID string, packageID string, args QuoteGasForExec) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "QuoteGasForExec",
		Arguments:  argsToMap(args),
	}
}

// GetTokenPrice exercises the GetTokenPrice choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetTokenPrice(contractID string, args GetTokenPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenPrice",
		Arguments:  argsToMap(args),
	}
}

// GetTokenPriceWithPackageID exercises the GetTokenPrice choice using the provided package ID instead of package name
func (t FeeQuoter) GetTokenPriceWithPackageID(contractID string, packageID string, args GetTokenPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenPrice",
		Arguments:  argsToMap(args),
	}
}

// GetDestinationChainGasPrice exercises the GetDestinationChainGasPrice choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetDestinationChainGasPrice(contractID string, args GetDestinationChainGasPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestinationChainGasPrice",
		Arguments:  argsToMap(args),
	}
}

// GetDestinationChainGasPriceWithPackageID exercises the GetDestinationChainGasPrice choice using the provided package ID instead of package name
func (t FeeQuoter) GetDestinationChainGasPriceWithPackageID(contractID string, packageID string, args GetDestinationChainGasPrice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestinationChainGasPrice",
		Arguments:  argsToMap(args),
	}
}

// UpdatePrices exercises the UpdatePrices choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) UpdatePrices(contractID string, args UpdatePrices) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "UpdatePrices",
		Arguments:  argsToMap(args),
	}
}

// UpdatePricesWithPackageID exercises the UpdatePrices choice using the provided package ID instead of package name
func (t FeeQuoter) UpdatePricesWithPackageID(contractID string, packageID string, args UpdatePrices) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "UpdatePrices",
		Arguments:  argsToMap(args),
	}
}

// GetTokenTransferFee exercises the GetTokenTransferFee choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetTokenTransferFee(contractID string, args GetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// GetTokenTransferFeeWithPackageID exercises the GetTokenTransferFee choice using the provided package ID instead of package name
func (t FeeQuoter) GetTokenTransferFeeWithPackageID(contractID string, packageID string, args GetTokenTransferFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetTokenTransferFee",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfig exercises the GetDestChainConfig choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetDestChainConfig(contractID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainConfigWithPackageID exercises the GetDestChainConfig choice using the provided package ID instead of package name
func (t FeeQuoter) GetDestChainConfigWithPackageID(contractID string, packageID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyFeeQuoterDestChainConfigUpdates exercises the ApplyFeeQuoterDestChainConfigUpdates choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyFeeQuoterDestChainConfigUpdates(contractID string, args ApplyFeeQuoterDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyFeeQuoterDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyFeeQuoterDestChainConfigUpdatesWithPackageID exercises the ApplyFeeQuoterDestChainConfigUpdates choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyFeeQuoterDestChainConfigUpdatesWithPackageID(contractID string, packageID string, args ApplyFeeQuoterDestChainConfigUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyFeeQuoterDestChainConfigUpdates",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FeeQuoter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeQuoter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FeeQuoter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) Get(contractID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t FeeQuoter) GetWithPackageID(contractID string, packageID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetFeeTokens exercises the GetFeeTokens choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) GetFeeTokens(contractID string, args GetFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// GetFeeTokensWithPackageID exercises the GetFeeTokens choice using the provided package ID instead of package name
func (t FeeQuoter) GetFeeTokensWithPackageID(contractID string, packageID string, args GetFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "GetFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// ApplyPriceUpdatersUpdate exercises the ApplyPriceUpdatersUpdate choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) ApplyPriceUpdatersUpdate(contractID string, args ApplyPriceUpdatersUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyPriceUpdatersUpdate",
		Arguments:  argsToMap(args),
	}
}

// ApplyPriceUpdatersUpdateWithPackageID exercises the ApplyPriceUpdatersUpdate choice using the provided package ID instead of package name
func (t FeeQuoter) ApplyPriceUpdatersUpdateWithPackageID(contractID string, packageID string, args ApplyPriceUpdatersUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "ApplyPriceUpdatersUpdate",
		Arguments:  argsToMap(args),
	}
}

// AddPriceUpdaters exercises the AddPriceUpdaters choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) AddPriceUpdaters(contractID string, args AddPriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "AddPriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// AddPriceUpdatersWithPackageID exercises the AddPriceUpdaters choice using the provided package ID instead of package name
func (t FeeQuoter) AddPriceUpdatersWithPackageID(contractID string, packageID string, args AddPriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "AddPriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemovePriceUpdaters exercises the RemovePriceUpdaters choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) RemovePriceUpdaters(contractID string, args RemovePriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemovePriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemovePriceUpdatersWithPackageID exercises the RemovePriceUpdaters choice using the provided package ID instead of package name
func (t FeeQuoter) RemovePriceUpdatersWithPackageID(contractID string, packageID string, args RemovePriceUpdaters) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemovePriceUpdaters",
		Arguments:  argsToMap(args),
	}
}

// RemoveFeeTokens exercises the RemoveFeeTokens choice on this FeeQuoter contract
// This method uses the package name in the template ID
func (t FeeQuoter) RemoveFeeTokens(contractID string, args RemoveFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemoveFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// RemoveFeeTokensWithPackageID exercises the RemoveFeeTokens choice using the provided package ID instead of package name
func (t FeeQuoter) RemoveFeeTokensWithPackageID(contractID string, packageID string, args RemoveFeeTokens) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "FeeQuoter"),
		ContractID: contractID,
		Choice:     "RemoveFeeTokens",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this FeeQuoter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t FeeQuoter) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t FeeQuoter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.FeeQuoter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for FeeQuoter

var _ api.IMCMSReceiver = (*FeeQuoter)(nil)

// FeeQuoterDestChainConfig is a Record type
type FeeQuoterDestChainConfig struct {
	IsEnabled                   types.BOOL    `json:"isEnabled"`
	MaxDataBytes                types.INT64   `json:"maxDataBytes"`
	MaxPerMsgGasLimit           types.INT64   `json:"maxPerMsgGasLimit"`
	DestGasOverhead             types.INT64   `json:"destGasOverhead"`
	DestGasPerPayloadByteBase   types.INT64   `json:"destGasPerPayloadByteBase"`
	DefaultTxGasLimit           types.INT64   `json:"defaultTxGasLimit"`
	LinkFeeMultiplierPercent    types.NUMERIC `json:"linkFeeMultiplierPercent"`
	DefaultTokenFeeUSD          types.NUMERIC `json:"defaultTokenFeeUSD"`
	DefaultTokenDestGasOverhead types.INT64   `json:"defaultTokenDestGasOverhead"`
}

// ToMap converts FeeQuoterDestChainConfig to a map for DAML arguments
func (t FeeQuoterDestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["maxDataBytes"] = int64(t.MaxDataBytes)

	m["maxPerMsgGasLimit"] = int64(t.MaxPerMsgGasLimit)

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destGasPerPayloadByteBase"] = int64(t.DestGasPerPayloadByteBase)

	m["defaultTxGasLimit"] = int64(t.DefaultTxGasLimit)

	m["linkFeeMultiplierPercent"] = t.LinkFeeMultiplierPercent

	m["defaultTokenFeeUSD"] = t.DefaultTokenFeeUSD

	m["defaultTokenDestGasOverhead"] = int64(t.DefaultTokenDestGasOverhead)

	return m
}

func (t FeeQuoterDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterDestChainConfig to hex string (Canton MCMS format)
func (t FeeQuoterDestChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterDestChainConfig from hex string (Canton MCMS format)
func (t *FeeQuoterDestChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeQuoterDestChainConfigArgs is a Record type
type FeeQuoterDestChainConfigArgs struct {
	DestChainSelector types.NUMERIC            `json:"destChainSelector"`
	DestChainConfig   FeeQuoterDestChainConfig `json:"destChainConfig"`
}

// ToMap converts FeeQuoterDestChainConfigArgs to a map for DAML arguments
func (t FeeQuoterDestChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["destChainConfig"] = model.NestedToDAMLValue(t.DestChainConfig)

	return m
}

func (t FeeQuoterDestChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeQuoterDestChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeQuoterDestChainConfigArgs to hex string (Canton MCMS format)
func (t FeeQuoterDestChainConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeQuoterDestChainConfigArgs from hex string (Canton MCMS format)
func (t *FeeQuoterDestChainConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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

// FinalityConfig is a variant/union type
type FinalityConfig struct {
	WaitForFinality *types.UNIT  `json:"WaitForFinality,omitempty"`
	WaitForSafe     *types.UNIT  `json:"WaitForSafe,omitempty"`
	BlockDepth      *types.INT64 `json:"BlockDepth,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for FinalityConfig
func (v FinalityConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for FinalityConfig
func (v *FinalityConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes FinalityConfig to hex string (Canton MCMS format)
func (v FinalityConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes FinalityConfig from hex string (Canton MCMS format)
func (v *FinalityConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v FinalityConfig) GetVariantTag() string {

	if v.WaitForFinality != nil {
		return "WaitForFinality"
	}

	if v.WaitForSafe != nil {
		return "WaitForSafe"
	}

	if v.BlockDepth != nil {
		return "BlockDepth"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v FinalityConfig) GetVariantValue() any {

	if v.WaitForFinality != nil {
		return v.WaitForFinality
	}

	if v.WaitForSafe != nil {
		return v.WaitForSafe
	}

	if v.BlockDepth != nil {
		return v.BlockDepth
	}

	return nil
}

var _ types.VARIANT = (*FinalityConfig)(nil)

// GetVariantTagByte implements types.VariantWithTagByte interface for MCMS numeric tag encoding
func (v FinalityConfig) GetVariantTagByte() byte {

	if v.WaitForFinality != nil {
		return 0
	}

	if v.WaitForSafe != nil {
		return 1
	}

	if v.BlockDepth != nil {
		return 2
	}

	return 0xFF // Invalid/unknown variant
}

var _ types.VariantWithTagByte = (*FinalityConfig)(nil)

// FinalizeExecute is a Record type
type FinalizeExecute struct {
	TokenAdminRegistryInstanceId types.TEXT                                `json:"tokenAdminRegistryInstanceId"`
	MaybePoolAddress             *chainlinkapi.RawInstanceAddress          `json:"maybePoolAddress" hex:"optional"`
	MaybeTicketReceiver          *types.PARTY                              `json:"maybeTicketReceiver" hex:"optional"`
	MaybeTokenReceiver           *types.PARTY                              `json:"maybeTokenReceiver" hex:"optional"`
	MaybeInstrumentId            *splice_api_token_holding_v1.InstrumentId `json:"maybeInstrumentId" hex:"optional"`
	MaybeAmount                  *types.TEXT                               `json:"maybeAmount" hex:"optional"`
	ReturnData                   types.TEXT                                `json:"returnData"`
}

// ToMap converts FinalizeExecute to a map for DAML arguments
func (t FinalizeExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	if t.MaybePoolAddress != nil {
		m["maybePoolAddress"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.MaybePoolAddress),
		}
	} else {
		m["maybePoolAddress"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			"value": nil,
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
			"value": nil,
		}
	}

	if t.MaybeInstrumentId != nil {
		m["maybeInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.MaybeInstrumentId),
		}
	} else {
		m["maybeInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.MaybeAmount != nil {
		m["maybeAmount"] = map[string]any{
			"_type": "optional",
			"value": string(*t.MaybeAmount),
		}
	} else {
		m["maybeAmount"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			"value": model.NestedToDAMLValue(*t.TokenReceiveTicket),
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["executionStateChanged"] = model.NestedToDAMLValue(t.ExecutionStateChanged)

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
			res = append(res, model.NestedToDAMLValue(e))
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

	m["ccipMessageSent"] = model.NestedToDAMLValue(t.CcipMessageSent)

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

// GasPriceUpdate is a Record type
type GasPriceUpdate struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	UsdPerUnitGas     types.NUMERIC `json:"usdPerUnitGas"`
}

// ToMap converts GasPriceUpdate to a map for DAML arguments
func (t GasPriceUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["usdPerUnitGas"] = t.UsdPerUnitGas

	return m
}

func (t GasPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GasPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GasPriceUpdate to hex string (Canton MCMS format)
func (t GasPriceUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GasPriceUpdate from hex string (Canton MCMS format)
func (t *GasPriceUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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

// Get is a Record type
type Get struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts Get to a map for DAML arguments
func (t Get) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Get) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Get) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Get to hex string (Canton MCMS format)
func (t Get) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get from hex string (Canton MCMS format)
func (t *Get) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMCMSParams is Get without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetMCMSParams struct {
}

// MarshalHex encodes GetMCMSParams to hex string for MCMS operationData.
func (t GetMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMCMSParams from hex string.
func (t *GetMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetCursedSubjects is a Record type
type GetCursedSubjects struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts GetCursedSubjects to a map for DAML arguments
func (t GetCursedSubjects) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetCursedSubjects) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetCursedSubjects) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetCursedSubjects to hex string (Canton MCMS format)
func (t GetCursedSubjects) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetCursedSubjects from hex string (Canton MCMS format)
func (t *GetCursedSubjects) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetCursedSubjectsMCMSParams is GetCursedSubjects without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetCursedSubjectsMCMSParams struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// MarshalHex encodes GetCursedSubjectsMCMSParams to hex string for MCMS operationData.
func (t GetCursedSubjectsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetCursedSubjectsMCMSParams from hex string.
func (t *GetCursedSubjectsMCMSParams) UnmarshalHex(data string) error {
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

// GetDestinationChainGasPrice is a Record type
type GetDestinationChainGasPrice struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts GetDestinationChainGasPrice to a map for DAML arguments
func (t GetDestinationChainGasPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDestinationChainGasPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDestinationChainGasPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDestinationChainGasPrice to hex string (Canton MCMS format)
func (t GetDestinationChainGasPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestinationChainGasPrice from hex string (Canton MCMS format)
func (t *GetDestinationChainGasPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestinationChainGasPriceMCMSParams is GetDestinationChainGasPrice without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDestinationChainGasPriceMCMSParams struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// MarshalHex encodes GetDestinationChainGasPriceMCMSParams to hex string for MCMS operationData.
func (t GetDestinationChainGasPriceMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestinationChainGasPriceMCMSParams from hex string.
func (t *GetDestinationChainGasPriceMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeTokens is a Record type
type GetFeeTokens struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetFeeTokens to a map for DAML arguments
func (t GetFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeTokens to hex string (Canton MCMS format)
func (t GetFeeTokens) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeTokens from hex string (Canton MCMS format)
func (t *GetFeeTokens) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeTokensMCMSParams is GetFeeTokens without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFeeTokensMCMSParams struct {
}

// MarshalHex encodes GetFeeTokensMCMSParams to hex string for MCMS operationData.
func (t GetFeeTokensMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeTokensMCMSParams from hex string.
func (t *GetFeeTokensMCMSParams) UnmarshalHex(data string) error {
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

// GetTokenConfigByCid is a Record type
type GetTokenConfigByCid struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenConfigByCid to a map for DAML arguments
func (t GetTokenConfigByCid) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenConfigByCid) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenConfigByCid) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenConfigByCid to hex string (Canton MCMS format)
func (t GetTokenConfigByCid) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenConfigByCid from hex string (Canton MCMS format)
func (t *GetTokenConfigByCid) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenConfigByCidMCMSParams is GetTokenConfigByCid without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenConfigByCidMCMSParams struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes GetTokenConfigByCidMCMSParams to hex string for MCMS operationData.
func (t GetTokenConfigByCidMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenConfigByCidMCMSParams from hex string.
func (t *GetTokenConfigByCidMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenPrice is a Record type
type GetTokenPrice struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenPrice to a map for DAML arguments
func (t GetTokenPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenPrice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenPrice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenPrice to hex string (Canton MCMS format)
func (t GetTokenPrice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenPrice from hex string (Canton MCMS format)
func (t *GetTokenPrice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenPriceMCMSParams is GetTokenPrice without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenPriceMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes GetTokenPriceMCMSParams to hex string for MCMS operationData.
func (t GetTokenPriceMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenPriceMCMSParams from hex string.
func (t *GetTokenPriceMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenTransferFee is a Record type
type GetTokenTransferFee struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	Token             splice_api_token_holding_v1.InstrumentId `json:"token"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts GetTokenTransferFee to a map for DAML arguments
func (t GetTokenTransferFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["token"] = model.NestedToDAMLValue(t.Token)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetTokenTransferFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTokenTransferFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTokenTransferFee to hex string (Canton MCMS format)
func (t GetTokenTransferFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenTransferFee from hex string (Canton MCMS format)
func (t *GetTokenTransferFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTokenTransferFeeMCMSParams is GetTokenTransferFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetTokenTransferFeeMCMSParams struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	Token             splice_api_token_holding_v1.InstrumentId `json:"token"`
}

// MarshalHex encodes GetTokenTransferFeeMCMSParams to hex string for MCMS operationData.
func (t GetTokenTransferFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTokenTransferFeeMCMSParams from hex string.
func (t *GetTokenTransferFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GlobalConfig is a Template type
type GlobalConfig struct {
	InstanceId         types.TEXT                          `json:"instanceId"`
	CcipOwner          types.PARTY                         `json:"ccipOwner"`
	ChainSelector      types.NUMERIC                       `json:"chainSelector"`
	DestChainConfigs   map[types.NUMERIC]DestChainConfig   `json:"destChainConfigs"`
	SourceChainConfigs map[types.NUMERIC]SourceChainConfig `json:"sourceChainConfigs"`
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
func (t GlobalConfig) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t GlobalConfig) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for GlobalConfig

var _ api.IMCMSReceiver = (*GlobalConfig)(nil)

// InboundPoolVerification is a Record type
type InboundPoolVerification struct {
	PoolInstanceId types.TEXT                        `json:"poolInstanceId"`
	PoolOwner      types.PARTY                       `json:"poolOwner"`
	PoolCCVs       []chainlinkapi.RawInstanceAddress `json:"poolCCVs"`
}

// ToMap converts InboundPoolVerification to a map for DAML arguments
func (t InboundPoolVerification) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t InboundPoolVerification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InboundPoolVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InboundPoolVerification to hex string (Canton MCMS format)
func (t InboundPoolVerification) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InboundPoolVerification from hex string (Canton MCMS format)
func (t *InboundPoolVerification) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsAdministrator is a Record type
type IsAdministrator struct {
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenConfigCid *types.CONTRACT_ID                       `json:"tokenConfigCid" hex:"optional"`
	Administrator  types.PARTY                              `json:"administrator"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts IsAdministrator to a map for DAML arguments
func (t IsAdministrator) ToMap() map[string]any {
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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsAdministrator to hex string (Canton MCMS format)
func (t IsAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsAdministrator from hex string (Canton MCMS format)
func (t *IsAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsAdministratorMCMSParams is IsAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsAdministratorMCMSParams struct {
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenConfigCid *types.CONTRACT_ID                       `json:"tokenConfigCid" hex:"optional"`
	Administrator  types.PARTY                              `json:"administrator"`
}

// MarshalHex encodes IsAdministratorMCMSParams to hex string for MCMS operationData.
func (t IsAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsAdministratorMCMSParams from hex string.
func (t *IsAdministratorMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursed is a Record type
type IsCursed struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller  types.PARTY                                `json:"caller"`
}

// ToMap converts IsCursed to a map for DAML arguments
func (t IsCursed) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsCursed) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsCursed) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsCursed to hex string (Canton MCMS format)
func (t IsCursed) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursed from hex string (Canton MCMS format)
func (t *IsCursed) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedMCMSParams is IsCursed without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsCursedMCMSParams struct {
	Context splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// MarshalHex encodes IsCursedMCMSParams to hex string for MCMS operationData.
func (t IsCursedMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedMCMSParams from hex string.
func (t *IsCursedMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedForChain is a Record type
type IsCursedForChain struct {
	ChainSelector types.NUMERIC                              `json:"chainSelector"`
	Context       splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller        types.PARTY                                `json:"caller"`
}

// ToMap converts IsCursedForChain to a map for DAML arguments
func (t IsCursedForChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t IsCursedForChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsCursedForChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsCursedForChain to hex string (Canton MCMS format)
func (t IsCursedForChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedForChain from hex string (Canton MCMS format)
func (t *IsCursedForChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsCursedForChainMCMSParams is IsCursedForChain without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type IsCursedForChainMCMSParams struct {
	ChainSelector types.NUMERIC                              `json:"chainSelector"`
	Context       splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// MarshalHex encodes IsCursedForChainMCMSParams to hex string for MCMS operationData.
func (t IsCursedForChainMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsCursedForChainMCMSParams from hex string.
func (t *IsCursedForChainMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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

// LocalAmountConversionResult is a Record type
type LocalAmountConversionResult struct {
	LocalAmount        types.NUMERIC `json:"localAmount"`
	TruncatedRemainder types.TEXT    `json:"truncatedRemainder"`
	WasTruncated       types.BOOL    `json:"wasTruncated"`
}

// ToMap converts LocalAmountConversionResult to a map for DAML arguments
func (t LocalAmountConversionResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["localAmount"] = t.LocalAmount

	m["truncatedRemainder"] = string(t.TruncatedRemainder)

	m["wasTruncated"] = bool(t.WasTruncated)

	return m
}

func (t LocalAmountConversionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LocalAmountConversionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LocalAmountConversionResult to hex string (Canton MCMS format)
func (t LocalAmountConversionResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LocalAmountConversionResult from hex string (Canton MCMS format)
func (t *LocalAmountConversionResult) UnmarshalHex(data string) error {
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
	Finality            DecodedFinality  `json:"finality"`
	CcvAndExecutorHash  types.TEXT       `json:"ccvAndExecutorHash"`
	OnRampAddress       types.TEXT       `json:"onRampAddress"`
	OffRampAddress      types.TEXT       `json:"offRampAddress" hex:"bytes"`
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

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	m["onRampAddress"] = string(t.OnRampAddress)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["sender"] = string(t.Sender)

	m["receiver"] = string(t.Receiver)

	m["destBlob"] = string(t.DestBlob)

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

// PriceUpdates is a Record type
type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate `json:"tokenPriceUpdates"`
	GasPriceUpdates   []GasPriceUpdate   `json:"gasPriceUpdates"`
}

// ToMap converts PriceUpdates to a map for DAML arguments
func (t PriceUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenPriceUpdates"] = func() []any {
		res := make([]any, 0, len(t.TokenPriceUpdates))
		for _, e := range t.TokenPriceUpdates {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["gasPriceUpdates"] = func() []any {
		res := make([]any, 0, len(t.GasPriceUpdates))
		for _, e := range t.GasPriceUpdates {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t PriceUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PriceUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PriceUpdates to hex string (Canton MCMS format)
func (t PriceUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PriceUpdates from hex string (Canton MCMS format)
func (t *PriceUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdminParams is a Record type
type ProposeAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// ToMap converts ProposeAdminParams to a map for DAML arguments
func (t ProposeAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t ProposeAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProposeAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProposeAdminParams to hex string (Canton MCMS format)
func (t ProposeAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdminParams from hex string (Canton MCMS format)
func (t *ProposeAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministrator is a Record type
type ProposeAdministrator struct {
	TokenConfigCid *types.CONTRACT_ID                       `json:"tokenConfigCid" hex:"optional"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin       types.PARTY                              `json:"newAdmin"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts ProposeAdministrator to a map for DAML arguments
func (t ProposeAdministrator) ToMap() map[string]any {
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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ProposeAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProposeAdministrator to hex string (Canton MCMS format)
func (t ProposeAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdministrator from hex string (Canton MCMS format)
func (t *ProposeAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministratorMCMSParams is ProposeAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ProposeAdministratorMCMSParams struct {
	TokenConfigCid *types.CONTRACT_ID                       `json:"tokenConfigCid" hex:"optional"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin       types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes ProposeAdministratorMCMSParams to hex string for MCMS operationData.
func (t ProposeAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProposeAdministratorMCMSParams from hex string.
func (t *ProposeAdministratorMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProposeAdministratorResult is a Record type
type ProposeAdministratorResult struct {
	TokenAdminRegistryCid types.CONTRACT_ID `json:"tokenAdminRegistryCid"`
	TokenConfigCid        types.CONTRACT_ID `json:"tokenConfigCid"`
	Created               types.BOOL        `json:"created"`
	Index                 types.INT64       `json:"index"`
}

// ToMap converts ProposeAdministratorResult to a map for DAML arguments
func (t ProposeAdministratorResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["created"] = bool(t.Created)

	m["index"] = int64(t.Index)

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

// QuoteGasForExec is a Record type
type QuoteGasForExec struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	NonCalldataGas    types.INT64                              `json:"nonCalldataGas"`
	CalldataSize      types.INT64                              `json:"calldataSize"`
	FeeToken          splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts QuoteGasForExec to a map for DAML arguments
func (t QuoteGasForExec) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["nonCalldataGas"] = int64(t.NonCalldataGas)

	m["calldataSize"] = int64(t.CalldataSize)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t QuoteGasForExec) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *QuoteGasForExec) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes QuoteGasForExec to hex string (Canton MCMS format)
func (t QuoteGasForExec) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExec from hex string (Canton MCMS format)
func (t *QuoteGasForExec) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExecMCMSParams is QuoteGasForExec without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type QuoteGasForExecMCMSParams struct {
	DestChainSelector types.NUMERIC                            `json:"destChainSelector"`
	NonCalldataGas    types.INT64                              `json:"nonCalldataGas"`
	CalldataSize      types.INT64                              `json:"calldataSize"`
	FeeToken          splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
}

// MarshalHex encodes QuoteGasForExecMCMSParams to hex string for MCMS operationData.
func (t QuoteGasForExecMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes QuoteGasForExecMCMSParams from hex string.
func (t *QuoteGasForExecMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// QuoteGasForExecResult is a Record type
type QuoteGasForExecResult struct {
	TotalGas          types.INT64   `json:"totalGas"`
	GasCostUSDCents   types.NUMERIC `json:"gasCostUSDCents"`
	FeeTokenPrice     types.NUMERIC `json:"feeTokenPrice"`
	PremiumMultiplier types.NUMERIC `json:"premiumMultiplier"`
}

// ToMap converts QuoteGasForExecResult to a map for DAML arguments
func (t QuoteGasForExecResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["totalGas"] = int64(t.TotalGas)

	m["gasCostUSDCents"] = t.GasCostUSDCents

	m["feeTokenPrice"] = t.FeeTokenPrice

	m["premiumMultiplier"] = t.PremiumMultiplier

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

// RMNRemote is a Template type
type RMNRemote struct {
	InstanceId      types.TEXT    `json:"instanceId"`
	RmnOwner        types.PARTY   `json:"rmnOwner"`
	CcipOwner       types.PARTY   `json:"ccipOwner"`
	CustomObservers []types.PARTY `json:"customObservers"`
	CursedSubjects  []types.TEXT  `json:"cursedSubjects" hex:"[]bytes"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RMNRemote) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RMNRemote) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RMNRemote) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []any {
		res := make([]any, 0, len(t.CursedSubjects))
		for _, e := range t.CursedSubjects {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RMNRemote) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnOwner"] = t.RmnOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["cursedSubjects"] = func() []any {
		res := make([]any, 0, len(t.CursedSubjects))
		for _, e := range t.CursedSubjects {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RMNRemote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RMNRemote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RMNRemote to hex string (Canton MCMS format)
func (t RMNRemote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RMNRemote from hex string (Canton MCMS format)
func (t *RMNRemote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RMNRemote

// UncurseChain exercises the UncurseChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseChain(contractID string, args UncurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseChain",
		Arguments:  argsToMap(args),
	}
}

// UncurseChainWithPackageID exercises the UncurseChain choice using the provided package ID instead of package name
func (t RMNRemote) UncurseChainWithPackageID(contractID string, packageID string, args UncurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseChain",
		Arguments:  argsToMap(args),
	}
}

// CurseChain exercises the CurseChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseChain(contractID string, args CurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseChain",
		Arguments:  argsToMap(args),
	}
}

// CurseChainWithPackageID exercises the CurseChain choice using the provided package ID instead of package name
func (t RMNRemote) CurseChainWithPackageID(contractID string, packageID string, args CurseChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseChain",
		Arguments:  argsToMap(args),
	}
}

// UncurseGlobal exercises the UncurseGlobal choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseGlobal(contractID string, args UncurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// UncurseGlobalWithPackageID exercises the UncurseGlobal choice using the provided package ID instead of package name
func (t RMNRemote) UncurseGlobalWithPackageID(contractID string, packageID string, args UncurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// CurseGlobal exercises the CurseGlobal choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseGlobal(contractID string, args CurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// CurseGlobalWithPackageID exercises the CurseGlobal choice using the provided package ID instead of package name
func (t RMNRemote) CurseGlobalWithPackageID(contractID string, packageID string, args CurseGlobal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseGlobal",
		Arguments:  argsToMap(args),
	}
}

// IsCursedForChain exercises the IsCursedForChain choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) IsCursedForChain(contractID string, args IsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// IsCursedForChainWithPackageID exercises the IsCursedForChain choice using the provided package ID instead of package name
func (t RMNRemote) IsCursedForChainWithPackageID(contractID string, packageID string, args IsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// Curse exercises the Curse choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Curse(contractID string, args Curse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Curse",
		Arguments:  argsToMap(args),
	}
}

// CurseWithPackageID exercises the Curse choice using the provided package ID instead of package name
func (t RMNRemote) CurseWithPackageID(contractID string, packageID string, args Curse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Curse",
		Arguments:  argsToMap(args),
	}
}

// Uncurse exercises the Uncurse choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) Uncurse(contractID string, args Uncurse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Uncurse",
		Arguments:  argsToMap(args),
	}
}

// UncurseWithPackageID exercises the Uncurse choice using the provided package ID instead of package name
func (t RMNRemote) UncurseWithPackageID(contractID string, packageID string, args Uncurse) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Uncurse",
		Arguments:  argsToMap(args),
	}
}

// CurseMultiple exercises the CurseMultiple choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) CurseMultiple(contractID string, args CurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// CurseMultipleWithPackageID exercises the CurseMultiple choice using the provided package ID instead of package name
func (t RMNRemote) CurseMultipleWithPackageID(contractID string, packageID string, args CurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "CurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// UncurseMultiple exercises the UncurseMultiple choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UncurseMultiple(contractID string, args UncurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// UncurseMultipleWithPackageID exercises the UncurseMultiple choice using the provided package ID instead of package name
func (t RMNRemote) UncurseMultipleWithPackageID(contractID string, packageID string, args UncurseMultiple) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UncurseMultiple",
		Arguments:  argsToMap(args),
	}
}

// IsCursed exercises the IsCursed choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) IsCursed(contractID string, args IsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursed",
		Arguments:  argsToMap(args),
	}
}

// IsCursedWithPackageID exercises the IsCursed choice using the provided package ID instead of package name
func (t RMNRemote) IsCursedWithPackageID(contractID string, packageID string, args IsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "IsCursed",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RMNRemote contract via the IRMNRemote interface
// This method uses the package name in the template ID
func (t RMNRemote) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RMNRemote) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetCursedSubjects exercises the GetCursedSubjects choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) GetCursedSubjects(contractID string, args GetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// GetCursedSubjectsWithPackageID exercises the GetCursedSubjects choice using the provided package ID instead of package name
func (t RMNRemote) GetCursedSubjectsWithPackageID(contractID string, packageID string, args GetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObservers exercises the AddCustomObservers choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) AddCustomObservers(contractID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObserversWithPackageID exercises the AddCustomObservers choice using the provided package ID instead of package name
func (t RMNRemote) AddCustomObserversWithPackageID(contractID string, packageID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObservers exercises the RemoveCustomObservers choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) RemoveCustomObservers(contractID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObserversWithPackageID exercises the RemoveCustomObservers choice using the provided package ID instead of package name
func (t RMNRemote) RemoveCustomObserversWithPackageID(contractID string, packageID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// UpdateCCIPOwner exercises the UpdateCCIPOwner choice on this RMNRemote contract
// This method uses the package name in the template ID
func (t RMNRemote) UpdateCCIPOwner(contractID string, args UpdateCCIPOwner) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UpdateCCIPOwner",
		Arguments:  argsToMap(args),
	}
}

// UpdateCCIPOwnerWithPackageID exercises the UpdateCCIPOwner choice using the provided package ID instead of package name
func (t RMNRemote) UpdateCCIPOwnerWithPackageID(contractID string, packageID string, args UpdateCCIPOwner) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "UpdateCCIPOwner",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this RMNRemote contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RMNRemote) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t RMNRemote) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// RMNRemotePublicFetch exercises the RMNRemote_PublicFetch choice on this RMNRemote contract via the IRMNRemote interface
// This method uses the package name in the template ID
func (t RMNRemote) RMNRemotePublicFetch(contractID string, args ccipapi.RMNRemotePublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// RMNRemotePublicFetchWithPackageID exercises the RMNRemote_PublicFetch choice using the provided package ID instead of package name
func (t RMNRemote) RMNRemotePublicFetchWithPackageID(contractID string, packageID string, args ccipapi.RMNRemotePublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteIsCursed exercises the RMNRemote_IsCursed choice on this RMNRemote contract via the IRMNRemote interface
// This method uses the package name in the template ID
func (t RMNRemote) RMNRemoteIsCursed(contractID string, args ccipapi.RMNRemoteIsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_IsCursed",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteIsCursedWithPackageID exercises the RMNRemote_IsCursed choice using the provided package ID instead of package name
func (t RMNRemote) RMNRemoteIsCursedWithPackageID(contractID string, packageID string, args ccipapi.RMNRemoteIsCursed) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_IsCursed",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteIsCursedForChain exercises the RMNRemote_IsCursedForChain choice on this RMNRemote contract via the IRMNRemote interface
// This method uses the package name in the template ID
func (t RMNRemote) RMNRemoteIsCursedForChain(contractID string, args ccipapi.RMNRemoteIsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteIsCursedForChainWithPackageID exercises the RMNRemote_IsCursedForChain choice using the provided package ID instead of package name
func (t RMNRemote) RMNRemoteIsCursedForChainWithPackageID(contractID string, packageID string, args ccipapi.RMNRemoteIsCursedForChain) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_IsCursedForChain",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteGetCursedSubjects exercises the RMNRemote_GetCursedSubjects choice on this RMNRemote contract via the IRMNRemote interface
// This method uses the package name in the template ID
func (t RMNRemote) RMNRemoteGetCursedSubjects(contractID string, args ccipapi.RMNRemoteGetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// RMNRemoteGetCursedSubjectsWithPackageID exercises the RMNRemote_GetCursedSubjects choice using the provided package ID instead of package name
func (t RMNRemote) RMNRemoteGetCursedSubjectsWithPackageID(contractID string, packageID string, args ccipapi.RMNRemoteGetCursedSubjects) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RMNRemote", "RMNRemote"),
		ContractID: contractID,
		Choice:     "RMNRemote_GetCursedSubjects",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for RMNRemote

var _ api.IMCMSReceiver = (*RMNRemote)(nil)

var _ ccipapi.IRMNRemote = (*RMNRemote)(nil)

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
		args["direction"] = model.NestedToDAMLValue(t.Direction)
	}

	if t.Mode != "" {
		args["mode"] = model.NestedToDAMLValue(t.Mode)
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
		args["direction"] = model.NestedToDAMLValue(t.Direction)
	}

	if t.Mode != "" {
		args["mode"] = model.NestedToDAMLValue(t.Mode)
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

// Archive exercises the Archive choice on this RateLimiter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RateLimiter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RateLimiter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this RateLimiter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RateLimiter) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RateLimiter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t RateLimiter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RateLimiter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for RateLimiter

var _ api.IMCMSReceiver = (*RateLimiter)(nil)

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

	m["issuerType"] = model.NestedToDAMLValue(t.IssuerType)

	m["issuerAddress"] = string(t.IssuerAddress)

	if t.VersionTag != nil {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": string(*t.VersionTag),
		}
	} else {
		m["versionTag"] = map[string]any{
			"_type": "optional",
			"value": nil,
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

// RemoveCustomObservers is a Record type
type RemoveCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObservers to a map for DAML arguments
func (t RemoveCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObservers to hex string (Canton MCMS format)
func (t RemoveCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObservers from hex string (Canton MCMS format)
func (t *RemoveCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveCustomObserversParams is a Record type
type RemoveCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObserversParams to a map for DAML arguments
func (t RemoveCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObserversParams to hex string (Canton MCMS format)
func (t RemoveCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObserversParams from hex string (Canton MCMS format)
func (t *RemoveCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveFeeTokens is a Record type
type RemoveFeeTokens struct {
	FeeTokensToRemove []splice_api_token_holding_v1.InstrumentId `json:"feeTokensToRemove"`
}

// ToMap converts RemoveFeeTokens to a map for DAML arguments
func (t RemoveFeeTokens) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokensToRemove"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToRemove))
		for _, e := range t.FeeTokensToRemove {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RemoveFeeTokens) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveFeeTokens) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveFeeTokens to hex string (Canton MCMS format)
func (t RemoveFeeTokens) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveFeeTokens from hex string (Canton MCMS format)
func (t *RemoveFeeTokens) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveFeeTokensParams is a Record type
type RemoveFeeTokensParams struct {
	FeeTokensToRemove []splice_api_token_holding_v1.InstrumentId `json:"feeTokensToRemove"`
}

// ToMap converts RemoveFeeTokensParams to a map for DAML arguments
func (t RemoveFeeTokensParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokensToRemove"] = func() []any {
		res := make([]any, 0, len(t.FeeTokensToRemove))
		for _, e := range t.FeeTokensToRemove {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RemoveFeeTokensParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveFeeTokensParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveFeeTokensParams to hex string (Canton MCMS format)
func (t RemoveFeeTokensParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveFeeTokensParams from hex string (Canton MCMS format)
func (t *RemoveFeeTokensParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemovePriceUpdaters is a Record type
type RemovePriceUpdaters struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemovePriceUpdaters to a map for DAML arguments
func (t RemovePriceUpdaters) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemovePriceUpdaters) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemovePriceUpdaters) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemovePriceUpdaters to hex string (Canton MCMS format)
func (t RemovePriceUpdaters) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemovePriceUpdaters from hex string (Canton MCMS format)
func (t *RemovePriceUpdaters) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SendingMessageDeps is a Record type
type SendingMessageDeps struct {
	Router             chainlinkapi.RawInstanceAddress `json:"router"`
	OnRamp             chainlinkapi.RawInstanceAddress `json:"onRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
}

// ToMap converts SendingMessageDeps to a map for DAML arguments
func (t SendingMessageDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["onRamp"] = model.NestedToDAMLValue(t.OnRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

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
	DestAddressBytesLength    types.INT64                               `json:"destAddressBytesLength"`
	SequenceNumber            types.NUMERIC                             `json:"sequenceNumber"`
	DestDefaultCCVs           []chainlinkapi.RawInstanceAddress         `json:"destDefaultCCVs"`
	RequiredCCVs              []chainlinkapi.RawInstanceAddress         `json:"requiredCCVs"`
	RequiredExecutor          *chainlinkapi.RawInstanceAddress          `json:"requiredExecutor" hex:"optional"`
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
	OffRampAddress            types.TEXT                                `json:"offRampAddress" hex:"bytes"`
	TokenReceiver             types.TEXT                                `json:"tokenReceiver"`
	TokenArgs                 types.TEXT                                `json:"tokenArgs"`
	FeeToken                  splice_api_token_holding_v1.InstrumentId  `json:"feeToken"`
	NetworkFeeUSDCents        types.NUMERIC                             `json:"networkFeeUSDCents"`
	ExpectedTokenInstrumentId *splice_api_token_holding_v1.InstrumentId `json:"expectedTokenInstrumentId" hex:"optional"`
	OutboundPoolCCVs          *[]chainlinkapi.RawInstanceAddress        `json:"outboundPoolCCVs" hex:"optional"`
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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	if t.DestChainSelector != "" {
		args["destChainSelector"] = t.DestChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destAddressBytesLength"] = int64(t.DestAddressBytesLength)

	if t.SequenceNumber != "" {
		args["sequenceNumber"] = t.SequenceNumber
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DestDefaultCCVs))
		for _, e := range t.DestDefaultCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.RequiredExecutor != nil {
		args["requiredExecutor"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RequiredExecutor),
		}
	} else {
		args["requiredExecutor"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorAddress"] = string(t.ExecutorAddress)

	if t.ExecutionMode != nil {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutionMode),
		}
	} else {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
	args["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	if t.NetworkFeeUSDCents != "" {
		args["networkFeeUSDCents"] = t.NetworkFeeUSDCents
	}

	if t.ExpectedTokenInstrumentId != nil {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExpectedTokenInstrumentId),
		}
	} else {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OutboundPoolCCVs != nil {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OutboundPoolCCVs),
		}
	} else {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorArgs"] = string(t.ExecutorArgs)

	if t.ExecutorFee != nil {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorFee),
		}
	} else {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenSendFee != nil {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenSendFee),
		}
	} else {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			"value": model.NestedToDAMLValue(*t.TokenSendData),
		}
	} else {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifierData"] = func() []any {
		res := make([]any, 0, len(t.VerifierData))
		for _, e := range t.VerifierData {
			res = append(res, model.NestedToDAMLValue(e))
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
			"value": model.NestedToDAMLValue(*t.Message),
		}
	} else {
		args["message"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["encodedMessage"] = string(t.EncodedMessage)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.State != "" {
		args["state"] = model.NestedToDAMLValue(t.State)
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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	if t.DestChainSelector != "" {
		args["destChainSelector"] = t.DestChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destAddressBytesLength"] = int64(t.DestAddressBytesLength)

	if t.SequenceNumber != "" {
		args["sequenceNumber"] = t.SequenceNumber
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destDefaultCCVs"] = func() []any {
		res := make([]any, 0, len(t.DestDefaultCCVs))
		for _, e := range t.DestDefaultCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.RequiredExecutor != nil {
		args["requiredExecutor"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RequiredExecutor),
		}
	} else {
		args["requiredExecutor"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorAddress"] = string(t.ExecutorAddress)

	if t.ExecutionMode != nil {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutionMode),
		}
	} else {
		args["executionMode"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
	args["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	if t.NetworkFeeUSDCents != "" {
		args["networkFeeUSDCents"] = t.NetworkFeeUSDCents
	}

	if t.ExpectedTokenInstrumentId != nil {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExpectedTokenInstrumentId),
		}
	} else {
		args["expectedTokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OutboundPoolCCVs != nil {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OutboundPoolCCVs),
		}
	} else {
		args["outboundPoolCCVs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executorArgs"] = string(t.ExecutorArgs)

	if t.ExecutorFee != nil {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorFee),
		}
	} else {
		args["executorFee"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenSendFee != nil {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenSendFee),
		}
	} else {
		args["tokenSendFee"] = map[string]any{
			"_type": "optional",
			"value": nil,
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
			"value": model.NestedToDAMLValue(*t.TokenSendData),
		}
	} else {
		args["tokenSendData"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifierData"] = func() []any {
		res := make([]any, 0, len(t.VerifierData))
		for _, e := range t.VerifierData {
			res = append(res, model.NestedToDAMLValue(e))
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
			"value": model.NestedToDAMLValue(*t.Message),
		}
	} else {
		args["message"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["encodedMessage"] = string(t.EncodedMessage)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.State != "" {
		args["state"] = model.NestedToDAMLValue(t.State)
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

// Archive exercises the Archive choice on this SendingMessageV1 contract via the ISendingMessage interface
// This method uses the package name in the template ID
func (t SendingMessageV1) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t SendingMessageV1) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SendingMessageAddCCVFee exercises the SendingMessage_AddCCVFee choice on this SendingMessageV1 contract via the ISendingMessage interface
// This method uses the package name in the template ID
func (t SendingMessageV1) SendingMessageAddCCVFee(contractID string, args ccipapi.SendingMessageAddCCVFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddCCVFee",
		Arguments:  argsToMap(args),
	}
}

// SendingMessageAddCCVFeeWithPackageID exercises the SendingMessage_AddCCVFee choice using the provided package ID instead of package name
func (t SendingMessageV1) SendingMessageAddCCVFeeWithPackageID(contractID string, packageID string, args ccipapi.SendingMessageAddCCVFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddCCVFee",
		Arguments:  argsToMap(args),
	}
}

// SendingMessageAddVerifierData exercises the SendingMessage_AddVerifierData choice on this SendingMessageV1 contract via the ISendingMessage interface
// This method uses the package name in the template ID
func (t SendingMessageV1) SendingMessageAddVerifierData(contractID string, args ccipapi.SendingMessageAddVerifierData) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddVerifierData",
		Arguments:  argsToMap(args),
	}
}

// SendingMessageAddVerifierDataWithPackageID exercises the SendingMessage_AddVerifierData choice using the provided package ID instead of package name
func (t SendingMessageV1) SendingMessageAddVerifierDataWithPackageID(contractID string, packageID string, args ccipapi.SendingMessageAddVerifierData) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddVerifierData",
		Arguments:  argsToMap(args),
	}
}

// SendingMessageAddExecutorFee exercises the SendingMessage_AddExecutorFee choice on this SendingMessageV1 contract via the ISendingMessage interface
// This method uses the package name in the template ID
func (t SendingMessageV1) SendingMessageAddExecutorFee(contractID string, args ccipapi.SendingMessageAddExecutorFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddExecutorFee",
		Arguments:  argsToMap(args),
	}
}

// SendingMessageAddExecutorFeeWithPackageID exercises the SendingMessage_AddExecutorFee choice using the provided package ID instead of package name
func (t SendingMessageV1) SendingMessageAddExecutorFeeWithPackageID(contractID string, packageID string, args ccipapi.SendingMessageAddExecutorFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.SendingMessageV1", "SendingMessage"),
		ContractID: contractID,
		Choice:     "SendingMessage_AddExecutorFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for SendingMessageV1

var _ ccipapi.ISendingMessage = (*SendingMessageV1)(nil)

// SetBurnMintFactory is a Record type
type SetBurnMintFactory struct {
	TokenConfigCid  types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	BurnMintFactory *types.CONTRACT_ID                       `json:"burnMintFactory" hex:"optional"`
	Caller          types.PARTY                              `json:"caller"`
}

// ToMap converts SetBurnMintFactory to a map for DAML arguments
func (t SetBurnMintFactory) ToMap() map[string]any {
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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetBurnMintFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetBurnMintFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetBurnMintFactory to hex string (Canton MCMS format)
func (t SetBurnMintFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetBurnMintFactory from hex string (Canton MCMS format)
func (t *SetBurnMintFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetBurnMintFactoryMCMSParams is SetBurnMintFactory without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetBurnMintFactoryMCMSParams struct {
	TokenConfigCid  types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	BurnMintFactory *types.CONTRACT_ID                       `json:"burnMintFactory" hex:"optional"`
}

// MarshalHex encodes SetBurnMintFactoryMCMSParams to hex string for MCMS operationData.
func (t SetBurnMintFactoryMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetBurnMintFactoryMCMSParams from hex string.
func (t *SetBurnMintFactoryMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetBurnMintFactoryParams is a Record type
type SetBurnMintFactoryParams struct {
	InstrumentId           splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	BurnMintFactoryAddress *chainlinkapi.RawInstanceAddress         `json:"burnMintFactoryAddress" hex:"optional"`
}

// ToMap converts SetBurnMintFactoryParams to a map for DAML arguments
func (t SetBurnMintFactoryParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.BurnMintFactoryAddress != nil {
		m["burnMintFactoryAddress"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.BurnMintFactoryAddress),
		}
	} else {
		m["burnMintFactoryAddress"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetBurnMintFactoryParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetBurnMintFactoryParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetBurnMintFactoryParams to hex string (Canton MCMS format)
func (t SetBurnMintFactoryParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetBurnMintFactoryParams from hex string (Canton MCMS format)
func (t *SetBurnMintFactoryParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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
	PoolInstanceId types.TEXT                        `json:"poolInstanceId"`
	PoolOwner      types.PARTY                       `json:"poolOwner"`
	PoolCCVs       []chainlinkapi.RawInstanceAddress `json:"poolCCVs"`
}

// ToMap converts SetInboundPoolCCVs to a map for DAML arguments
func (t SetInboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
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
	PoolCCVs []chainlinkapi.RawInstanceAddress `json:"poolCCVs"`
}

// ToMap converts SetOutboundPoolCCVs to a map for DAML arguments
func (t SetOutboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
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

// SetPool is a Record type
type SetPool struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool      *PoolRegistration                        `json:"tokenPool" hex:"optional"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts SetPool to a map for DAML arguments
func (t SetPool) ToMap() map[string]any {
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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPool to hex string (Canton MCMS format)
func (t SetPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPool from hex string (Canton MCMS format)
func (t *SetPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPoolMCMSParams is SetPool without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetPoolMCMSParams struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool      *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// MarshalHex encodes SetPoolMCMSParams to hex string for MCMS operationData.
func (t SetPoolMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPoolMCMSParams from hex string.
func (t *SetPoolMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPoolParams is a Record type
type SetPoolParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// ToMap converts SetPoolParams to a map for DAML arguments
func (t SetPoolParams) ToMap() map[string]any {
	m := make(map[string]any)

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

	return m
}

func (t SetPoolParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPoolParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPoolParams to hex string (Canton MCMS format)
func (t SetPoolParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPoolParams from hex string (Canton MCMS format)
func (t *SetPoolParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferFactory is a Record type
type SetTransferFactory struct {
	TokenConfigCid  types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TransferFactory *types.CONTRACT_ID                       `json:"transferFactory" hex:"optional"`
	Caller          types.PARTY                              `json:"caller"`
}

// ToMap converts SetTransferFactory to a map for DAML arguments
func (t SetTransferFactory) ToMap() map[string]any {
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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t SetTransferFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetTransferFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetTransferFactory to hex string (Canton MCMS format)
func (t SetTransferFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferFactory from hex string (Canton MCMS format)
func (t *SetTransferFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferFactoryMCMSParams is SetTransferFactory without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type SetTransferFactoryMCMSParams struct {
	TokenConfigCid  types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TransferFactory *types.CONTRACT_ID                       `json:"transferFactory" hex:"optional"`
}

// MarshalHex encodes SetTransferFactoryMCMSParams to hex string for MCMS operationData.
func (t SetTransferFactoryMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferFactoryMCMSParams from hex string.
func (t *SetTransferFactoryMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferFactoryParams is a Record type
type SetTransferFactoryParams struct {
	InstrumentId           splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TransferFactoryAddress *chainlinkapi.RawInstanceAddress         `json:"transferFactoryAddress" hex:"optional"`
}

// ToMap converts SetTransferFactoryParams to a map for DAML arguments
func (t SetTransferFactoryParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.TransferFactoryAddress != nil {
		m["transferFactoryAddress"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TransferFactoryAddress),
		}
	} else {
		m["transferFactoryAddress"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetTransferFactoryParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetTransferFactoryParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetTransferFactoryParams to hex string (Canton MCMS format)
func (t SetTransferFactoryParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferFactoryParams from hex string (Canton MCMS format)
func (t *SetTransferFactoryParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        types.BOOL                        `json:"isEnabled"`
	OnRampAddresses  []types.TEXT                      `json:"onRampAddresses" hex:"[]bytes"`
	DefaultCCVs      []chainlinkapi.RawInstanceAddress `json:"defaultCCVs"`
	LaneMandatedCCVs []chainlinkapi.RawInstanceAddress `json:"laneMandatedCCVs"`
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
	SourceChainSelector types.NUMERIC                     `json:"sourceChainSelector"`
	IsEnabled           types.BOOL                        `json:"isEnabled"`
	OnRampAddresses     []types.TEXT                      `json:"onRampAddresses" hex:"[]bytes"`
	DefaultCCVs         []chainlinkapi.RawInstanceAddress `json:"defaultCCVs"`
	LaneMandatedCCVs    []chainlinkapi.RawInstanceAddress `json:"laneMandatedCCVs"`
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

// TimestampedPrice is a Record type
type TimestampedPrice struct {
	Price     types.NUMERIC   `json:"price"`
	Timestamp types.TIMESTAMP `json:"timestamp"`
}

// ToMap converts TimestampedPrice to a map for DAML arguments
func (t TimestampedPrice) ToMap() map[string]any {
	m := make(map[string]any)

	m["price"] = t.Price

	m["timestamp"] = t.Timestamp

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

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	InstanceId types.TEXT  `json:"instanceId"`
	Owner      types.PARTY `json:"owner"`
	EntryCount types.INT64 `json:"entryCount"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenAdminRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenAdminRegistry) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenAdminRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["entryCount"] = int64(t.EntryCount)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenAdminRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["entryCount"] = int64(t.EntryCount)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistry to hex string (Canton MCMS format)
func (t TokenAdminRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistry from hex string (Canton MCMS format)
func (t *TokenAdminRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenAdminRegistry

// ProposeAdministrator exercises the ProposeAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) ProposeAdministrator(contractID string, args ProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// ProposeAdministratorWithPackageID exercises the ProposeAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ProposeAdministratorWithPackageID(contractID string, packageID string, args ProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecute exercises the FinalizeExecute choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) FinalizeExecute(contractID string, args FinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// FinalizeExecuteWithPackageID exercises the FinalizeExecute choice using the provided package ID instead of package name
func (t TokenAdminRegistry) FinalizeExecuteWithPackageID(contractID string, packageID string, args FinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVs exercises the SetInboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetInboundPoolCCVs(contractID string, args SetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetInboundPoolCCVsWithPackageID exercises the SetInboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetInboundPoolCCVsWithPackageID(contractID string, packageID string, args SetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSend exercises the AddTokenSend choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AddTokenSend(contractID string, args AddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendWithPackageID exercises the AddTokenSend choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AddTokenSendWithPackageID(contractID string, packageID string, args AddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFee exercises the AddTokenSendFee choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AddTokenSendFee(contractID string, args AddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// AddTokenSendFeeWithPackageID exercises the AddTokenSendFee choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AddTokenSendFeeWithPackageID(contractID string, packageID string, args AddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVs exercises the SetOutboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetOutboundPoolCCVs(contractID string, args SetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// SetOutboundPoolCCVsWithPackageID exercises the SetOutboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetOutboundPoolCCVsWithPackageID(contractID string, packageID string, args SetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// ConsumeReceiveTicket exercises the ConsumeReceiveTicket choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) ConsumeReceiveTicket(contractID string, args ConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// ConsumeReceiveTicketWithPackageID exercises the ConsumeReceiveTicket choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ConsumeReceiveTicketWithPackageID(contractID string, packageID string, args ConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// IsAdministrator exercises the IsAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) IsAdministrator(contractID string, args IsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// IsAdministratorWithPackageID exercises the IsAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) IsAdministratorWithPackageID(contractID string, packageID string, args IsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TransferAdminRole exercises the TransferAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TransferAdminRole(contractID string, args TransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TransferAdminRoleWithPackageID exercises the TransferAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TransferAdminRoleWithPackageID(contractID string, packageID string, args TransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// AcceptAdminRole exercises the AcceptAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) AcceptAdminRole(contractID string, args AcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// AcceptAdminRoleWithPackageID exercises the AcceptAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) AcceptAdminRoleWithPackageID(contractID string, packageID string, args AcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// SetBurnMintFactory exercises the SetBurnMintFactory choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetBurnMintFactory(contractID string, args SetBurnMintFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetBurnMintFactory",
		Arguments:  argsToMap(args),
	}
}

// SetBurnMintFactoryWithPackageID exercises the SetBurnMintFactory choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetBurnMintFactoryWithPackageID(contractID string, packageID string, args SetBurnMintFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetBurnMintFactory",
		Arguments:  argsToMap(args),
	}
}

// SetTransferFactory exercises the SetTransferFactory choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetTransferFactory(contractID string, args SetTransferFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetTransferFactory",
		Arguments:  argsToMap(args),
	}
}

// SetTransferFactoryWithPackageID exercises the SetTransferFactory choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetTransferFactoryWithPackageID(contractID string, packageID string, args SetTransferFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetTransferFactory",
		Arguments:  argsToMap(args),
	}
}

// SetPool exercises the SetPool choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) SetPool(contractID string, args SetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetPool",
		Arguments:  argsToMap(args),
	}
}

// SetPoolWithPackageID exercises the SetPool choice using the provided package ID instead of package name
func (t TokenAdminRegistry) SetPoolWithPackageID(contractID string, packageID string, args SetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "SetPool",
		Arguments:  argsToMap(args),
	}
}

// GetTokenConfigByCid exercises the GetTokenConfigByCid choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) GetTokenConfigByCid(contractID string, args GetTokenConfigByCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "GetTokenConfigByCid",
		Arguments:  argsToMap(args),
	}
}

// GetTokenConfigByCidWithPackageID exercises the GetTokenConfigByCid choice using the provided package ID instead of package name
func (t TokenAdminRegistry) GetTokenConfigByCidWithPackageID(contractID string, packageID string, args GetTokenConfigByCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "GetTokenConfigByCid",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Get(contractID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t TokenAdminRegistry) GetWithPackageID(contractID string, packageID string, args Get) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t TokenAdminRegistry) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for TokenAdminRegistry

var _ api.IMCMSReceiver = (*TokenAdminRegistry)(nil)

// TokenConfig is a Template type
type TokenConfig struct {
	InstanceId         types.TEXT                               `json:"instanceId"`
	RegistryInstanceId types.TEXT                               `json:"registryInstanceId"`
	RegistryOwner      types.PARTY                              `json:"registryOwner"`
	Index              types.INT64                              `json:"index"`
	IsCCIPManaged      types.BOOL                               `json:"isCCIPManaged"`
	InstrumentId       splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Admin              *types.PARTY                             `json:"admin" hex:"optional"`
	PendingAdmin       *types.PARTY                             `json:"pendingAdmin" hex:"optional"`
	TokenPool          *PoolRegistration                        `json:"tokenPool" hex:"optional"`
	TransferFactory    *types.CONTRACT_ID                       `json:"transferFactory" hex:"optional"`
	BurnMintFactory    *types.CONTRACT_ID                       `json:"burnMintFactory" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TokenConfig) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenConfig")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TokenConfig) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenConfig")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TokenConfig) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryInstanceId"] = string(t.RegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryOwner"] = t.RegistryOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["index"] = int64(t.Index)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isCCIPManaged"] = bool(t.IsCCIPManaged)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.Admin != nil {
		args["admin"] = map[string]any{
			"_type": "optional",
			"value": (*t.Admin).ToMap(),
		}
	} else {
		args["admin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.PendingAdmin != nil {
		args["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.PendingAdmin).ToMap(),
		}
	} else {
		args["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TokenPool != nil {
		args["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		args["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TransferFactory != nil {
		args["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TransferFactory),
		}
	} else {
		args["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.BurnMintFactory != nil {
		args["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.BurnMintFactory),
		}
	} else {
		args["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenConfig) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryInstanceId"] = string(t.RegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryOwner"] = t.RegistryOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["index"] = int64(t.Index)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isCCIPManaged"] = bool(t.IsCCIPManaged)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	if t.Admin != nil {
		args["admin"] = map[string]any{
			"_type": "optional",
			"value": (*t.Admin).ToMap(),
		}
	} else {
		args["admin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.PendingAdmin != nil {
		args["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.PendingAdmin).ToMap(),
		}
	} else {
		args["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TokenPool != nil {
		args["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPool),
		}
	} else {
		args["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TransferFactory != nil {
		args["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TransferFactory),
		}
	} else {
		args["transferFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.BurnMintFactory != nil {
		args["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.BurnMintFactory),
		}
	} else {
		args["burnMintFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfig to hex string (Canton MCMS format)
func (t TokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfig from hex string (Canton MCMS format)
func (t *TokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TokenConfig

// Archive exercises the Archive choice on this TokenConfig contract
// This method uses the package name in the template ID
func (t TokenConfig) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenConfig"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenConfig) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenConfig"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TokenPriceUpdate is a Record type
type TokenPriceUpdate struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	UsdPerToken  types.NUMERIC                            `json:"usdPerToken"`
}

// ToMap converts TokenPriceUpdate to a map for DAML arguments
func (t TokenPriceUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["usdPerToken"] = t.UsdPerToken

	return m
}

func (t TokenPriceUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenPriceUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenPriceUpdate to hex string (Canton MCMS format)
func (t TokenPriceUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenPriceUpdate from hex string (Canton MCMS format)
func (t *TokenPriceUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner                    types.PARTY                              `json:"ccipOwner"`
	CcvOwners                    []types.PARTY                            `json:"ccvOwners"`
	VerifiedCCVs                 []chainlinkapi.RawInstanceAddress        `json:"verifiedCCVs"`
	RequiredInboundPoolCCVs      []chainlinkapi.RawInstanceAddress        `json:"requiredInboundPoolCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	PoolAddress                  chainlinkapi.RawInstanceAddress          `json:"poolAddress"`
	PoolOwner                    types.PARTY                              `json:"poolOwner"`
	Receiver                     types.PARTY                              `json:"receiver"`
	TokenReceiver                types.PARTY                              `json:"tokenReceiver"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount                       types.TEXT                               `json:"amount"`
	SourcePoolData               types.TEXT                               `json:"sourcePoolData"`
	MessageId                    types.TEXT                               `json:"messageId"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
	Finality                     FinalityConfig                           `json:"finality"`
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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredInboundPoolCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredInboundPoolCCVs))
		for _, e := range t.RequiredInboundPoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["amount"] = string(t.Amount)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = model.NestedToDAMLValue(t.Finality)

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
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredInboundPoolCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredInboundPoolCCVs))
		for _, e := range t.RequiredInboundPoolCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["amount"] = string(t.Amount)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourcePoolData"] = string(t.SourcePoolData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["finality"] = model.NestedToDAMLValue(t.Finality)

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

// Consume exercises the Consume choice on this TokenReceiveTicket contract
// This method uses the package name in the template ID
func (t TokenReceiveTicket) Consume(contractID string, args Consume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Consume",
		Arguments:  argsToMap(args),
	}
}

// ConsumeWithPackageID exercises the Consume choice using the provided package ID instead of package name
func (t TokenReceiveTicket) ConsumeWithPackageID(contractID string, packageID string, args Consume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Consume",
		Arguments:  argsToMap(args),
	}
}

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
	args["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

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
	args["tokenReceiveTicketCid"] = model.NestedToDAMLValue(t.TokenReceiveTicketCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = model.NestedToDAMLValue(t.Event)

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
	VerifiedCCVs                 []chainlinkapi.RawInstanceAddress        `json:"verifiedCCVs"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	PoolAddress                  chainlinkapi.RawInstanceAddress          `json:"poolAddress"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	SourceAmount                 types.TEXT                               `json:"sourceAmount"`
	LocalAmount                  types.NUMERIC                            `json:"localAmount"`
	SourcePoolData               types.TEXT                               `json:"sourcePoolData"`
	MessageId                    types.TEXT                               `json:"messageId"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
	Finality                     FinalityConfig                           `json:"finality"`
	Output                       TokenReceiveTicketClaimedOutput          `json:"output"`
}

// ToMap converts TokenReceiveTicketClaimedEvent to a map for DAML arguments
func (t TokenReceiveTicketClaimedEvent) ToMap() map[string]any {
	m := make(map[string]any)

	m["verifiedCCVs"] = func() []any {
		res := make([]any, 0, len(t.VerifiedCCVs))
		for _, e := range t.VerifiedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["tokenAdminRegistryInstanceId"] = string(t.TokenAdminRegistryInstanceId)

	m["poolAddress"] = model.NestedToDAMLValue(t.PoolAddress)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["sourceAmount"] = string(t.SourceAmount)

	m["localAmount"] = t.LocalAmount

	m["sourcePoolData"] = string(t.SourcePoolData)

	m["messageId"] = string(t.MessageId)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["finality"] = model.NestedToDAMLValue(t.Finality)

	m["output"] = model.NestedToDAMLValue(t.Output)

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

	m["transferInstructionCid"] = model.NestedToDAMLValue(t.TransferInstructionCid)

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
	Amount           types.TEXT                               `json:"amount"`
	DestTokenAddress types.TEXT                               `json:"destTokenAddress"`
	ExtraData        types.TEXT                               `json:"extraData"`
}

// ToMap converts TokenSendData to a map for DAML arguments
func (t TokenSendData) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = string(t.Amount)

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

// TokenTransferFeeConfig is a Record type
type TokenTransferFeeConfig struct {
	IsEnabled         types.BOOL    `json:"isEnabled"`
	FeeUSD            types.NUMERIC `json:"feeUSD"`
	DestGasOverhead   types.INT64   `json:"destGasOverhead"`
	DestBytesOverhead types.INT64   `json:"destBytesOverhead"`
}

// ToMap converts TokenTransferFeeConfig to a map for DAML arguments
func (t TokenTransferFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["feeUSD"] = t.FeeUSD

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	return m
}

func (t TokenTransferFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferFeeConfig to hex string (Canton MCMS format)
func (t TokenTransferFeeConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferFeeConfig from hex string (Canton MCMS format)
func (t *TokenTransferFeeConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferV1 is a Record type
type TokenTransferV1 struct {
	Amount             types.TEXT `json:"amount"`
	SourcePoolAddress  types.TEXT `json:"sourcePoolAddress"`
	SourceTokenAddress types.TEXT `json:"sourceTokenAddress"`
	DestTokenAddress   types.TEXT `json:"destTokenAddress"`
	TokenReceiver      types.TEXT `json:"tokenReceiver"`
	ExtraData          types.TEXT `json:"extraData"`
}

// ToMap converts TokenTransferV1 to a map for DAML arguments
func (t TokenTransferV1) ToMap() map[string]any {
	m := make(map[string]any)

	m["amount"] = string(t.Amount)

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

// TransferAdminParams is a Record type
type TransferAdminParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// ToMap converts TransferAdminParams to a map for DAML arguments
func (t TransferAdminParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t TransferAdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferAdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferAdminParams to hex string (Canton MCMS format)
func (t TransferAdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminParams from hex string (Canton MCMS format)
func (t *TransferAdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferAdminRole is a Record type
type TransferAdminRole struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin       types.PARTY                              `json:"newAdmin"`
	Caller         types.PARTY                              `json:"caller"`
}

// ToMap converts TransferAdminRole to a map for DAML arguments
func (t TransferAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenConfigCid"] = model.NestedToDAMLValue(t.TokenConfigCid)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TransferAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferAdminRole to hex string (Canton MCMS format)
func (t TransferAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminRole from hex string (Canton MCMS format)
func (t *TransferAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferAdminRoleMCMSParams is TransferAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TransferAdminRoleMCMSParams struct {
	TokenConfigCid types.CONTRACT_ID                        `json:"tokenConfigCid"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin       types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes TransferAdminRoleMCMSParams to hex string for MCMS operationData.
func (t TransferAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferAdminRoleMCMSParams from hex string.
func (t *TransferAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Uncurse is a Record type
type Uncurse struct {
	Subject types.TEXT `json:"subject" hex:"bytes"`
}

// ToMap converts Uncurse to a map for DAML arguments
func (t Uncurse) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t Uncurse) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Uncurse) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Uncurse to hex string (Canton MCMS format)
func (t Uncurse) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Uncurse from hex string (Canton MCMS format)
func (t *Uncurse) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseChain is a Record type
type UncurseChain struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts UncurseChain to a map for DAML arguments
func (t UncurseChain) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t UncurseChain) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseChain) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseChain to hex string (Canton MCMS format)
func (t UncurseChain) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseChain from hex string (Canton MCMS format)
func (t *UncurseChain) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseChainParams is a Record type
type UncurseChainParams struct {
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts UncurseChainParams to a map for DAML arguments
func (t UncurseChainParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t UncurseChainParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseChainParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseChainParams to hex string (Canton MCMS format)
func (t UncurseChainParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseChainParams from hex string (Canton MCMS format)
func (t *UncurseChainParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseGlobal is a Record type
type UncurseGlobal struct {
}

// ToMap converts UncurseGlobal to a map for DAML arguments
func (t UncurseGlobal) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UncurseGlobal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseGlobal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseGlobal to hex string (Canton MCMS format)
func (t UncurseGlobal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseGlobal from hex string (Canton MCMS format)
func (t *UncurseGlobal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseMultiple is a Record type
type UncurseMultiple struct {
	Subjects []types.TEXT `json:"subjects" hex:"[]bytes"`
}

// ToMap converts UncurseMultiple to a map for DAML arguments
func (t UncurseMultiple) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UncurseMultiple) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseMultiple) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseMultiple to hex string (Canton MCMS format)
func (t UncurseMultiple) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseMultiple from hex string (Canton MCMS format)
func (t *UncurseMultiple) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseMultipleParams is a Record type
type UncurseMultipleParams struct {
	Subjects []types.TEXT `json:"subjects" hex:"[]bytes"`
}

// ToMap converts UncurseMultipleParams to a map for DAML arguments
func (t UncurseMultipleParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subjects"] = func() []any {
		res := make([]any, 0, len(t.Subjects))
		for _, e := range t.Subjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t UncurseMultipleParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseMultipleParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseMultipleParams to hex string (Canton MCMS format)
func (t UncurseMultipleParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseMultipleParams from hex string (Canton MCMS format)
func (t *UncurseMultipleParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UncurseParams is a Record type
type UncurseParams struct {
	Subject types.TEXT `json:"subject" hex:"bytes"`
}

// ToMap converts UncurseParams to a map for DAML arguments
func (t UncurseParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["subject"] = string(t.Subject)

	return m
}

func (t UncurseParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UncurseParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UncurseParams to hex string (Canton MCMS format)
func (t UncurseParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UncurseParams from hex string (Canton MCMS format)
func (t *UncurseParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateCCIPOwner is a Record type
type UpdateCCIPOwner struct {
	NewCCIPOwner types.PARTY `json:"newCCIPOwner"`
}

// ToMap converts UpdateCCIPOwner to a map for DAML arguments
func (t UpdateCCIPOwner) ToMap() map[string]any {
	m := make(map[string]any)

	m["newCCIPOwner"] = t.NewCCIPOwner.ToMap()

	return m
}

func (t UpdateCCIPOwner) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateCCIPOwner) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateCCIPOwner to hex string (Canton MCMS format)
func (t UpdateCCIPOwner) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateCCIPOwner from hex string (Canton MCMS format)
func (t *UpdateCCIPOwner) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePrices is a Record type
type UpdatePrices struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       types.PARTY  `json:"caller"`
}

// ToMap converts UpdatePrices to a map for DAML arguments
func (t UpdatePrices) ToMap() map[string]any {
	m := make(map[string]any)

	m["priceUpdates"] = model.NestedToDAMLValue(t.PriceUpdates)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t UpdatePrices) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdatePrices) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdatePrices to hex string (Canton MCMS format)
func (t UpdatePrices) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePrices from hex string (Canton MCMS format)
func (t *UpdatePrices) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePricesMCMSParams is UpdatePrices without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type UpdatePricesMCMSParams struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
}

// MarshalHex encodes UpdatePricesMCMSParams to hex string for MCMS operationData.
func (t UpdatePricesMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePricesMCMSParams from hex string.
func (t *UpdatePricesMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdatePricesParams is a Record type
type UpdatePricesParams struct {
	PriceUpdates PriceUpdates `json:"priceUpdates"`
	Caller       types.PARTY  `json:"caller"`
}

// ToMap converts UpdatePricesParams to a map for DAML arguments
func (t UpdatePricesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["priceUpdates"] = model.NestedToDAMLValue(t.PriceUpdates)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t UpdatePricesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdatePricesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdatePricesParams to hex string (Canton MCMS format)
func (t UpdatePricesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdatePricesParams from hex string (Canton MCMS format)
func (t *UpdatePricesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// VerifierData is a Record type
type VerifierData struct {
	CcvInstanceId        types.TEXT    `json:"ccvInstanceId"`
	CcvOwner             types.PARTY   `json:"ccvOwner"`
	VersionTag           types.TEXT    `json:"versionTag" hex:"bytes"`
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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AcceptAdminRole(args AcceptAdminRole) (*bind.EncodedChoice, error)
	AcceptAdminRoleMCMSParams(args AcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error)
	AddCCVFee(args AddCCVFee) (*bind.EncodedChoice, error)
	AddCCVFeeMCMSParams(args AddCCVFeeMCMSParams) (*bind.EncodedChoice, error)
	AddCCVVerification(args AddCCVVerification) (*bind.EncodedChoice, error)
	AddCCVVerificationMCMSParams(args AddCCVVerificationMCMSParams) (*bind.EncodedChoice, error)
	AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error)
	AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error)
	AddExecutorFee(args AddExecutorFee) (*bind.EncodedChoice, error)
	AddExecutorFeeMCMSParams(args AddExecutorFeeMCMSParams) (*bind.EncodedChoice, error)
	AddPriceUpdaters(args AddPriceUpdaters) (*bind.EncodedChoice, error)
	AddTokenSend(args AddTokenSend) (*bind.EncodedChoice, error)
	AddTokenSendFee(args AddTokenSendFee) (*bind.EncodedChoice, error)
	AddVerifierData(args AddVerifierData) (*bind.EncodedChoice, error)
	AddVerifierDataMCMSParams(args AddVerifierDataMCMSParams) (*bind.EncodedChoice, error)
	ApplyDestChainConfigUpdates(args ApplyDestChainConfigUpdates) (*bind.EncodedChoice, error)
	ApplyDestChainConfigUpdatesParams(args ApplyDestChainConfigUpdatesParams) (*bind.EncodedChoice, error)
	ApplyFeeQuoterDestChainConfigUpdates(args ApplyFeeQuoterDestChainConfigUpdates) (*bind.EncodedChoice, error)
	ApplyFeeQuoterDestChainConfigUpdatesParams(args ApplyFeeQuoterDestChainConfigUpdatesParams) (*bind.EncodedChoice, error)
	ApplyPriceUpdatersUpdate(args ApplyPriceUpdatersUpdate) (*bind.EncodedChoice, error)
	ApplyPriceUpdatersUpdateParams(args ApplyPriceUpdatersUpdateParams) (*bind.EncodedChoice, error)
	ApplySourceChainConfigUpdates(args ApplySourceChainConfigUpdates) (*bind.EncodedChoice, error)
	ApplySourceChainConfigUpdatesParams(args ApplySourceChainConfigUpdatesParams) (*bind.EncodedChoice, error)
	BuildMessage(args BuildMessage) (*bind.EncodedChoice, error)
	CancelExecute(args CancelExecute) (*bind.EncodedChoice, error)
	CancelExecuteMCMSParams(args CancelExecuteMCMSParams) (*bind.EncodedChoice, error)
	Consume(args Consume) (*bind.EncodedChoice, error)
	ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error)
	ConsumeReceiveTicket(args ConsumeReceiveTicket) (*bind.EncodedChoice, error)
	ConsumeReceiveTicketMCMSParams(args ConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error)
	Curse(args Curse) (*bind.EncodedChoice, error)
	CurseChain(args CurseChain) (*bind.EncodedChoice, error)
	CurseChainParams(args CurseChainParams) (*bind.EncodedChoice, error)
	CurseGlobal(args CurseGlobal) (*bind.EncodedChoice, error)
	CurseMultiple(args CurseMultiple) (*bind.EncodedChoice, error)
	CurseMultipleParams(args CurseMultipleParams) (*bind.EncodedChoice, error)
	CurseParams(args CurseParams) (*bind.EncodedChoice, error)
	FeeTokenAmount(args FeeTokenAmount) (*bind.EncodedChoice, error)
	FeeTokenAmountMCMSParams(args FeeTokenAmountMCMSParams) (*bind.EncodedChoice, error)
	FinalizeExecute(args FinalizeExecute) (*bind.EncodedChoice, error)
	FinalizeFee(args FinalizeFee) (*bind.EncodedChoice, error)
	FinalizeSend(args FinalizeSend) (*bind.EncodedChoice, error)
	Get(args Get) (*bind.EncodedChoice, error)
	GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error)
	GetCursedSubjects(args GetCursedSubjects) (*bind.EncodedChoice, error)
	GetCursedSubjectsMCMSParams(args GetCursedSubjectsMCMSParams) (*bind.EncodedChoice, error)
	GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error)
	GetDestChainConfigMCMSParams(args GetDestChainConfigMCMSParams) (*bind.EncodedChoice, error)
	GetDestinationChainGasPrice(args GetDestinationChainGasPrice) (*bind.EncodedChoice, error)
	GetDestinationChainGasPriceMCMSParams(args GetDestinationChainGasPriceMCMSParams) (*bind.EncodedChoice, error)
	GetFeeTokens(args GetFeeTokens) (*bind.EncodedChoice, error)
	GetFeeTokensMCMSParams(args GetFeeTokensMCMSParams) (*bind.EncodedChoice, error)
	GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error)
	GetSourceChainConfigMCMSParams(args GetSourceChainConfigMCMSParams) (*bind.EncodedChoice, error)
	GetTokenConfigByCid(args GetTokenConfigByCid) (*bind.EncodedChoice, error)
	GetTokenConfigByCidMCMSParams(args GetTokenConfigByCidMCMSParams) (*bind.EncodedChoice, error)
	GetTokenPrice(args GetTokenPrice) (*bind.EncodedChoice, error)
	GetTokenPriceMCMSParams(args GetTokenPriceMCMSParams) (*bind.EncodedChoice, error)
	GetTokenTransferFee(args GetTokenTransferFee) (*bind.EncodedChoice, error)
	GetTokenTransferFeeMCMSParams(args GetTokenTransferFeeMCMSParams) (*bind.EncodedChoice, error)
	IsAdministrator(args IsAdministrator) (*bind.EncodedChoice, error)
	IsAdministratorMCMSParams(args IsAdministratorMCMSParams) (*bind.EncodedChoice, error)
	IsCursed(args IsCursed) (*bind.EncodedChoice, error)
	IsCursedMCMSParams(args IsCursedMCMSParams) (*bind.EncodedChoice, error)
	IsCursedForChain(args IsCursedForChain) (*bind.EncodedChoice, error)
	IsCursedForChainMCMSParams(args IsCursedForChainMCMSParams) (*bind.EncodedChoice, error)
	ProposeAdministrator(args ProposeAdministrator) (*bind.EncodedChoice, error)
	ProposeAdministratorMCMSParams(args ProposeAdministratorMCMSParams) (*bind.EncodedChoice, error)
	QuoteGasForExec(args QuoteGasForExec) (*bind.EncodedChoice, error)
	QuoteGasForExecMCMSParams(args QuoteGasForExecMCMSParams) (*bind.EncodedChoice, error)
	RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error)
	RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error)
	RemoveFeeTokens(args RemoveFeeTokens) (*bind.EncodedChoice, error)
	RemoveFeeTokensParams(args RemoveFeeTokensParams) (*bind.EncodedChoice, error)
	RemovePriceUpdaters(args RemovePriceUpdaters) (*bind.EncodedChoice, error)
	SetBurnMintFactory(args SetBurnMintFactory) (*bind.EncodedChoice, error)
	SetBurnMintFactoryMCMSParams(args SetBurnMintFactoryMCMSParams) (*bind.EncodedChoice, error)
	SetBurnMintFactoryParams(args SetBurnMintFactoryParams) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
	SetInboundPoolCCVs(args SetInboundPoolCCVs) (*bind.EncodedChoice, error)
	SetOutboundPoolCCVs(args SetOutboundPoolCCVs) (*bind.EncodedChoice, error)
	SetPool(args SetPool) (*bind.EncodedChoice, error)
	SetPoolMCMSParams(args SetPoolMCMSParams) (*bind.EncodedChoice, error)
	SetPoolParams(args SetPoolParams) (*bind.EncodedChoice, error)
	SetTransferFactory(args SetTransferFactory) (*bind.EncodedChoice, error)
	SetTransferFactoryMCMSParams(args SetTransferFactoryMCMSParams) (*bind.EncodedChoice, error)
	SetTransferFactoryParams(args SetTransferFactoryParams) (*bind.EncodedChoice, error)
	TransferAdminRole(args TransferAdminRole) (*bind.EncodedChoice, error)
	TransferAdminRoleMCMSParams(args TransferAdminRoleMCMSParams) (*bind.EncodedChoice, error)
	Uncurse(args Uncurse) (*bind.EncodedChoice, error)
	UncurseChain(args UncurseChain) (*bind.EncodedChoice, error)
	UncurseChainParams(args UncurseChainParams) (*bind.EncodedChoice, error)
	UncurseGlobal(args UncurseGlobal) (*bind.EncodedChoice, error)
	UncurseMultiple(args UncurseMultiple) (*bind.EncodedChoice, error)
	UncurseMultipleParams(args UncurseMultipleParams) (*bind.EncodedChoice, error)
	UncurseParams(args UncurseParams) (*bind.EncodedChoice, error)
	UpdateCCIPOwner(args UpdateCCIPOwner) (*bind.EncodedChoice, error)
	UpdatePrices(args UpdatePrices) (*bind.EncodedChoice, error)
	UpdatePricesMCMSParams(args UpdatePricesMCMSParams) (*bind.EncodedChoice, error)
	UpdatePricesParams(args UpdatePricesParams) (*bind.EncodedChoice, error)
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

// AcceptAdminRole encodes parameters for the AcceptAdminRole choice.
func (e *encoder) AcceptAdminRole(args AcceptAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptAdminRole", args)
}

// AcceptAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the AcceptAdminRole choice.
func (e *encoder) AcceptAdminRoleMCMSParams(args AcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptAdminRole", args)
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

// AddCustomObservers encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// AddCustomObserversParams encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// AddExecutorFee encodes parameters for the AddExecutorFee choice.
func (e *encoder) AddExecutorFee(args AddExecutorFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddExecutorFee", args)
}

// AddExecutorFeeMCMSParams encodes MCMS parameters (without Caller) for the AddExecutorFee choice.
func (e *encoder) AddExecutorFeeMCMSParams(args AddExecutorFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddExecutorFee", args)
}

// AddPriceUpdaters encodes parameters for the AddPriceUpdaters choice.
func (e *encoder) AddPriceUpdaters(args AddPriceUpdaters) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddPriceUpdaters", args)
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

// ApplyDestChainConfigUpdatesParams encodes parameters for the ApplyDestChainConfigUpdates choice.
func (e *encoder) ApplyDestChainConfigUpdatesParams(args ApplyDestChainConfigUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyDestChainConfigUpdates", args)
}

// ApplyFeeQuoterDestChainConfigUpdates encodes parameters for the ApplyFeeQuoterDestChainConfigUpdates choice.
func (e *encoder) ApplyFeeQuoterDestChainConfigUpdates(args ApplyFeeQuoterDestChainConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyFeeQuoterDestChainConfigUpdates", args)
}

// ApplyFeeQuoterDestChainConfigUpdatesParams encodes parameters for the ApplyFeeQuoterDestChainConfigUpdates choice.
func (e *encoder) ApplyFeeQuoterDestChainConfigUpdatesParams(args ApplyFeeQuoterDestChainConfigUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyFeeQuoterDestChainConfigUpdates", args)
}

// ApplyPriceUpdatersUpdate encodes parameters for the ApplyPriceUpdatersUpdate choice.
func (e *encoder) ApplyPriceUpdatersUpdate(args ApplyPriceUpdatersUpdate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyPriceUpdatersUpdate", args)
}

// ApplyPriceUpdatersUpdateParams encodes parameters for the ApplyPriceUpdatersUpdate choice.
func (e *encoder) ApplyPriceUpdatersUpdateParams(args ApplyPriceUpdatersUpdateParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyPriceUpdatersUpdate", args)
}

// ApplySourceChainConfigUpdates encodes parameters for the ApplySourceChainConfigUpdates choice.
func (e *encoder) ApplySourceChainConfigUpdates(args ApplySourceChainConfigUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplySourceChainConfigUpdates", args)
}

// ApplySourceChainConfigUpdatesParams encodes parameters for the ApplySourceChainConfigUpdates choice.
func (e *encoder) ApplySourceChainConfigUpdatesParams(args ApplySourceChainConfigUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplySourceChainConfigUpdates", args)
}

// BuildMessage encodes parameters for the BuildMessage choice.
func (e *encoder) BuildMessage(args BuildMessage) (*bind.EncodedChoice, error) {
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

// Consume encodes parameters for the Consume choice.
func (e *encoder) Consume(args Consume) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Consume", args)
}

// ConsumeCapacity encodes parameters for the ConsumeCapacity choice.
func (e *encoder) ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeCapacity", args)
}

// ConsumeReceiveTicket encodes parameters for the ConsumeReceiveTicket choice.
func (e *encoder) ConsumeReceiveTicket(args ConsumeReceiveTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeReceiveTicket", args)
}

// ConsumeReceiveTicketMCMSParams encodes MCMS parameters (without Caller) for the ConsumeReceiveTicket choice.
func (e *encoder) ConsumeReceiveTicketMCMSParams(args ConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeReceiveTicket", args)
}

// Curse encodes parameters for the Curse choice.
func (e *encoder) Curse(args Curse) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Curse", args)
}

// CurseChain encodes parameters for the CurseChain choice.
func (e *encoder) CurseChain(args CurseChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseChain", args)
}

// CurseChainParams encodes parameters for the CurseChain choice.
func (e *encoder) CurseChainParams(args CurseChainParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseChain", args)
}

// CurseGlobal encodes parameters for the CurseGlobal choice.
func (e *encoder) CurseGlobal(args CurseGlobal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseGlobal", args)
}

// CurseMultiple encodes parameters for the CurseMultiple choice.
func (e *encoder) CurseMultiple(args CurseMultiple) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseMultiple", args)
}

// CurseMultipleParams encodes parameters for the CurseMultiple choice.
func (e *encoder) CurseMultipleParams(args CurseMultipleParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CurseMultiple", args)
}

// CurseParams encodes parameters for the Curse choice.
func (e *encoder) CurseParams(args CurseParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Curse", args)
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

// Get encodes parameters for the Get choice.
func (e *encoder) Get(args Get) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetMCMSParams encodes MCMS parameters (without Caller) for the Get choice.
func (e *encoder) GetMCMSParams(args GetMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get", args)
}

// GetCursedSubjects encodes parameters for the GetCursedSubjects choice.
func (e *encoder) GetCursedSubjects(args GetCursedSubjects) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetCursedSubjects", args)
}

// GetCursedSubjectsMCMSParams encodes MCMS parameters (without Caller) for the GetCursedSubjects choice.
func (e *encoder) GetCursedSubjectsMCMSParams(args GetCursedSubjectsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetCursedSubjects", args)
}

// GetDestChainConfig encodes parameters for the GetDestChainConfig choice.
func (e *encoder) GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig", args)
}

// GetDestChainConfigMCMSParams encodes MCMS parameters (without Caller) for the GetDestChainConfig choice.
func (e *encoder) GetDestChainConfigMCMSParams(args GetDestChainConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig", args)
}

// GetDestinationChainGasPrice encodes parameters for the GetDestinationChainGasPrice choice.
func (e *encoder) GetDestinationChainGasPrice(args GetDestinationChainGasPrice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestinationChainGasPrice", args)
}

// GetDestinationChainGasPriceMCMSParams encodes MCMS parameters (without Caller) for the GetDestinationChainGasPrice choice.
func (e *encoder) GetDestinationChainGasPriceMCMSParams(args GetDestinationChainGasPriceMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestinationChainGasPrice", args)
}

// GetFeeTokens encodes parameters for the GetFeeTokens choice.
func (e *encoder) GetFeeTokens(args GetFeeTokens) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFeeTokens", args)
}

// GetFeeTokensMCMSParams encodes MCMS parameters (without Caller) for the GetFeeTokens choice.
func (e *encoder) GetFeeTokensMCMSParams(args GetFeeTokensMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFeeTokens", args)
}

// GetSourceChainConfig encodes parameters for the GetSourceChainConfig choice.
func (e *encoder) GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSourceChainConfig", args)
}

// GetSourceChainConfigMCMSParams encodes MCMS parameters (without Caller) for the GetSourceChainConfig choice.
func (e *encoder) GetSourceChainConfigMCMSParams(args GetSourceChainConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSourceChainConfig", args)
}

// GetTokenConfigByCid encodes parameters for the GetTokenConfigByCid choice.
func (e *encoder) GetTokenConfigByCid(args GetTokenConfigByCid) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenConfigByCid", args)
}

// GetTokenConfigByCidMCMSParams encodes MCMS parameters (without Caller) for the GetTokenConfigByCid choice.
func (e *encoder) GetTokenConfigByCidMCMSParams(args GetTokenConfigByCidMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenConfigByCid", args)
}

// GetTokenPrice encodes parameters for the GetTokenPrice choice.
func (e *encoder) GetTokenPrice(args GetTokenPrice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenPrice", args)
}

// GetTokenPriceMCMSParams encodes MCMS parameters (without Caller) for the GetTokenPrice choice.
func (e *encoder) GetTokenPriceMCMSParams(args GetTokenPriceMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenPrice", args)
}

// GetTokenTransferFee encodes parameters for the GetTokenTransferFee choice.
func (e *encoder) GetTokenTransferFee(args GetTokenTransferFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenTransferFee", args)
}

// GetTokenTransferFeeMCMSParams encodes MCMS parameters (without Caller) for the GetTokenTransferFee choice.
func (e *encoder) GetTokenTransferFeeMCMSParams(args GetTokenTransferFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTokenTransferFee", args)
}

// IsAdministrator encodes parameters for the IsAdministrator choice.
func (e *encoder) IsAdministrator(args IsAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsAdministrator", args)
}

// IsAdministratorMCMSParams encodes MCMS parameters (without Caller) for the IsAdministrator choice.
func (e *encoder) IsAdministratorMCMSParams(args IsAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsAdministrator", args)
}

// IsCursed encodes parameters for the IsCursed choice.
func (e *encoder) IsCursed(args IsCursed) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursed", args)
}

// IsCursedMCMSParams encodes MCMS parameters (without Caller) for the IsCursed choice.
func (e *encoder) IsCursedMCMSParams(args IsCursedMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursed", args)
}

// IsCursedForChain encodes parameters for the IsCursedForChain choice.
func (e *encoder) IsCursedForChain(args IsCursedForChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursedForChain", args)
}

// IsCursedForChainMCMSParams encodes MCMS parameters (without Caller) for the IsCursedForChain choice.
func (e *encoder) IsCursedForChainMCMSParams(args IsCursedForChainMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsCursedForChain", args)
}

// ProposeAdministrator encodes parameters for the ProposeAdministrator choice.
func (e *encoder) ProposeAdministrator(args ProposeAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProposeAdministrator", args)
}

// ProposeAdministratorMCMSParams encodes MCMS parameters (without Caller) for the ProposeAdministrator choice.
func (e *encoder) ProposeAdministratorMCMSParams(args ProposeAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProposeAdministrator", args)
}

// QuoteGasForExec encodes parameters for the QuoteGasForExec choice.
func (e *encoder) QuoteGasForExec(args QuoteGasForExec) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("QuoteGasForExec", args)
}

// QuoteGasForExecMCMSParams encodes MCMS parameters (without Caller) for the QuoteGasForExec choice.
func (e *encoder) QuoteGasForExecMCMSParams(args QuoteGasForExecMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("QuoteGasForExec", args)
}

// RemoveCustomObservers encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// RemoveCustomObserversParams encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// RemoveFeeTokens encodes parameters for the RemoveFeeTokens choice.
func (e *encoder) RemoveFeeTokens(args RemoveFeeTokens) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveFeeTokens", args)
}

// RemoveFeeTokensParams encodes parameters for the RemoveFeeTokens choice.
func (e *encoder) RemoveFeeTokensParams(args RemoveFeeTokensParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveFeeTokens", args)
}

// RemovePriceUpdaters encodes parameters for the RemovePriceUpdaters choice.
func (e *encoder) RemovePriceUpdaters(args RemovePriceUpdaters) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemovePriceUpdaters", args)
}

// SetBurnMintFactory encodes parameters for the SetBurnMintFactory choice.
func (e *encoder) SetBurnMintFactory(args SetBurnMintFactory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetBurnMintFactory", args)
}

// SetBurnMintFactoryMCMSParams encodes MCMS parameters (without Caller) for the SetBurnMintFactory choice.
func (e *encoder) SetBurnMintFactoryMCMSParams(args SetBurnMintFactoryMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetBurnMintFactory", args)
}

// SetBurnMintFactoryParams encodes parameters for the SetBurnMintFactory choice.
func (e *encoder) SetBurnMintFactoryParams(args SetBurnMintFactoryParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetBurnMintFactory", args)
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

// SetPool encodes parameters for the SetPool choice.
func (e *encoder) SetPool(args SetPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// SetPoolMCMSParams encodes MCMS parameters (without Caller) for the SetPool choice.
func (e *encoder) SetPoolMCMSParams(args SetPoolMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// SetPoolParams encodes parameters for the SetPool choice.
func (e *encoder) SetPoolParams(args SetPoolParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPool", args)
}

// SetTransferFactory encodes parameters for the SetTransferFactory choice.
func (e *encoder) SetTransferFactory(args SetTransferFactory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferFactory", args)
}

// SetTransferFactoryMCMSParams encodes MCMS parameters (without Caller) for the SetTransferFactory choice.
func (e *encoder) SetTransferFactoryMCMSParams(args SetTransferFactoryMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferFactory", args)
}

// SetTransferFactoryParams encodes parameters for the SetTransferFactory choice.
func (e *encoder) SetTransferFactoryParams(args SetTransferFactoryParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferFactory", args)
}

// TransferAdminRole encodes parameters for the TransferAdminRole choice.
func (e *encoder) TransferAdminRole(args TransferAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferAdminRole", args)
}

// TransferAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the TransferAdminRole choice.
func (e *encoder) TransferAdminRoleMCMSParams(args TransferAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferAdminRole", args)
}

// Uncurse encodes parameters for the Uncurse choice.
func (e *encoder) Uncurse(args Uncurse) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Uncurse", args)
}

// UncurseChain encodes parameters for the UncurseChain choice.
func (e *encoder) UncurseChain(args UncurseChain) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseChain", args)
}

// UncurseChainParams encodes parameters for the UncurseChain choice.
func (e *encoder) UncurseChainParams(args UncurseChainParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseChain", args)
}

// UncurseGlobal encodes parameters for the UncurseGlobal choice.
func (e *encoder) UncurseGlobal(args UncurseGlobal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseGlobal", args)
}

// UncurseMultiple encodes parameters for the UncurseMultiple choice.
func (e *encoder) UncurseMultiple(args UncurseMultiple) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseMultiple", args)
}

// UncurseMultipleParams encodes parameters for the UncurseMultiple choice.
func (e *encoder) UncurseMultipleParams(args UncurseMultipleParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UncurseMultiple", args)
}

// UncurseParams encodes parameters for the Uncurse choice.
func (e *encoder) UncurseParams(args UncurseParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Uncurse", args)
}

// UpdateCCIPOwner encodes parameters for the UpdateCCIPOwner choice.
func (e *encoder) UpdateCCIPOwner(args UpdateCCIPOwner) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateCCIPOwner", args)
}

// UpdatePrices encodes parameters for the UpdatePrices choice.
func (e *encoder) UpdatePrices(args UpdatePrices) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// UpdatePricesMCMSParams encodes MCMS parameters (without Caller) for the UpdatePrices choice.
func (e *encoder) UpdatePricesMCMSParams(args UpdatePricesMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// UpdatePricesParams encodes parameters for the UpdatePrices choice.
func (e *encoder) UpdatePricesParams(args UpdatePricesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdatePrices", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
