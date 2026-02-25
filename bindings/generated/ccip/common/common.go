package common

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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
	PackageName = "ccip-common"
	PackageID   = "f73c3eced5f3dc106c9162398a76e31d4614e5c425c32d50b10ee11d89f80ba7"
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

// AddTokenSend is a Record type
type AddTokenSend struct {
	PoolInstanceId   types.TEXT                               `json:"poolInstanceId"`
	PoolOwner        types.PARTY                              `json:"poolOwner"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount           types.NUMERIC                            `json:"amount"`
	DestTokenAddress types.TEXT                               `json:"destTokenAddress"`
	ExtraData        types.TEXT                               `json:"extraData"`
	PoolRequiredCCVs []RawInstanceAddress                     `json:"poolRequiredCCVs"`
	Caller           types.PARTY                              `json:"caller"`
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

	m["poolRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolRequiredCCVs))
		for _, e := range t.PoolRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

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
	Caller            types.PARTY   `json:"caller"`
}

// ToMap converts AddTokenSendFee to a map for DAML arguments
func (t AddTokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["caller"] = t.Caller.ToMap()

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

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	CcipOwner       types.PARTY `json:"ccipOwner"`
	StorageLocation types.TEXT  `json:"storageLocation"`
}

// ToMap converts CrossChainVerifierView to a map for DAML arguments
func (t CrossChainVerifierView) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["storageLocation"] = string(t.StorageLocation)

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
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT                                 `json:"verifierArgs"`
	Caller            types.PARTY                                `json:"caller"`
}

// ToMap converts CrossChainVerifierForwardToVerifier to a map for DAML arguments
func (t CrossChainVerifierForwardToVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
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
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	VerifierResults     types.TEXT                                 `json:"verifierResults"`
	Caller              types.PARTY                                `json:"caller"`
}

// ToMap converts CrossChainVerifierVerifyMessage to a map for DAML arguments
func (t CrossChainVerifierVerifyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
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
	IsEnabled                 types.BOOL           `json:"isEnabled"`
	DefaultExecutor           RawInstanceAddress   `json:"defaultExecutor"`
	OffRampAddress            types.TEXT           `json:"offRampAddress"`
	LaneMandatedCCVs          []RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents types.NUMERIC        `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   types.NUMERIC        `json:"tokenNetworkFeeUSDCents"`
}

// ToMap converts DestChainConfig to a map for DAML arguments
func (t DestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["defaultExecutor"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DefaultExecutor).(mapper); ok {
			return m.toMap()
		}
		return t.DefaultExecutor
	}()

	m["offRampAddress"] = string(t.OffRampAddress)

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

// ExecutingMessageV1 is a Template type
type ExecutingMessageV1 struct {
	CcipOwner                         types.PARTY           `json:"ccipOwner"`
	Message                           MessageV1             `json:"message"`
	MessageId                         types.TEXT            `json:"messageId"`
	Receiver                          types.PARTY           `json:"receiver"`
	TokenReceiver                     *types.PARTY          `json:"tokenReceiver"`
	Executor                          types.PARTY           `json:"executor"`
	ObservingParties                  []types.PARTY         `json:"observingParties"`
	CcvVerifications                  []CCVVerification     `json:"ccvVerifications"`
	InboundPoolCCVs                   *[]RawInstanceAddress `json:"inboundPoolCCVs"`
	OffRampInstanceAddress            RawInstanceAddress    `json:"offRampInstanceAddress"`
	TokenAdminRegistryInstanceAddress RawInstanceAddress    `json:"tokenAdminRegistryInstanceAddress"`
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
	args["offRampInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRampInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

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
	args["offRampInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRampInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

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

// ExecutingMessageV1Archive exercises the ExecutingMessageV1_Archive choice on this ExecutingMessageV1 contract
// This method uses the package name in the template ID
func (t ExecutingMessageV1) ExecutingMessageV1Archive(contractID string, args ExecutingMessageV1Archive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "ExecutingMessageV1_Archive",
		Arguments:  argsToMap(args),
	}
}

// ExecutingMessageV1ArchiveWithPackageID exercises the ExecutingMessageV1_Archive choice using the provided package ID instead of package name
func (t ExecutingMessageV1) ExecutingMessageV1ArchiveWithPackageID(contractID string, packageID string, args ExecutingMessageV1Archive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.ExecutingMessageV1", "ExecutingMessageV1"),
		ContractID: contractID,
		Choice:     "ExecutingMessageV1_Archive",
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

// ExecutingMessageV1Archive is a Record type
type ExecutingMessageV1Archive struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutingMessageV1Archive to a map for DAML arguments
func (t ExecutingMessageV1Archive) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutingMessageV1Archive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutingMessageV1Archive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutingMessageV1Archive to hex string (Canton MCMS format)
func (t ExecutingMessageV1Archive) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutingMessageV1Archive from hex string (Canton MCMS format)
func (t *ExecutingMessageV1Archive) UnmarshalHex(data string) error {
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

// FinalizeFee is a Record type
type FinalizeFee struct {
	FeeTokenPrice     types.NUMERIC `json:"feeTokenPrice"`
	PremiumMultiplier types.NUMERIC `json:"premiumMultiplier"`
}

// ToMap converts FinalizeFee to a map for DAML arguments
func (t FinalizeFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenPrice"] = t.FeeTokenPrice

	m["premiumMultiplier"] = t.PremiumMultiplier

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

// GlobalConfig is a Template type
type GlobalConfig struct {
	InstanceId         types.TEXT    `json:"instanceId"`
	CcipOwner          types.PARTY   `json:"ccipOwner"`
	ChainSelector      types.NUMERIC `json:"chainSelector"`
	OnRampAddress      types.TEXT    `json:"onRampAddress"`
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
	args["onRampAddress"] = string(t.OnRampAddress)

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
	args["onRampAddress"] = string(t.OnRampAddress)

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

// UpdateDestChainConfig exercises the UpdateDestChainConfig choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) UpdateDestChainConfig(contractID string, args UpdateDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// UpdateDestChainConfigWithPackageID exercises the UpdateDestChainConfig choice using the provided package ID instead of package name
func (t GlobalConfig) UpdateDestChainConfigWithPackageID(contractID string, packageID string, args UpdateDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t GlobalConfig) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// UpdateSourceChainConfig exercises the UpdateSourceChainConfig choice on this GlobalConfig contract
// This method uses the package name in the template ID
func (t GlobalConfig) UpdateSourceChainConfig(contractID string, args UpdateSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// UpdateSourceChainConfigWithPackageID exercises the UpdateSourceChainConfig choice using the provided package ID instead of package name
func (t GlobalConfig) UpdateSourceChainConfigWithPackageID(contractID string, packageID string, args UpdateSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// IssuerType is an enum type
type IssuerType string

const (
	IssuerTypeIssuerType_CCV IssuerType = "IssuerType_CCV"

	IssuerTypeIssuerType_Pool IssuerType = "IssuerType_Pool"

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
	TokenTransfer       *TokenTransferV1 `json:"tokenTransfer"`
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

// RawInstanceAddress is a Record type
type RawInstanceAddress struct {
	Unpack types.TEXT `json:"unpack"`
}

// ToMap converts RawInstanceAddress to a map for DAML arguments
func (t RawInstanceAddress) ToMap() map[string]any {
	m := make(map[string]any)

	m["unpack"] = string(t.Unpack)

	return m
}

func (t RawInstanceAddress) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RawInstanceAddress) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RawInstanceAddress to hex string (Canton MCMS format)
func (t RawInstanceAddress) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RawInstanceAddress from hex string (Canton MCMS format)
func (t *RawInstanceAddress) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Receipt is a Record type
type Receipt struct {
	IssuerType        IssuerType    `json:"issuerType"`
	IssuerAddress     types.TEXT    `json:"issuerAddress"`
	VersionTag        *types.TEXT   `json:"versionTag"`
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

// SendingMessageV1 is a Template type
type SendingMessageV1 struct {
	RouterInstanceAddress             RawInstanceAddress                       `json:"routerInstanceAddress"`
	OnRampInstanceAddress             RawInstanceAddress                       `json:"onRampInstanceAddress"`
	TokenAdminRegistryInstanceAddress RawInstanceAddress                       `json:"tokenAdminRegistryInstanceAddress"`
	CcipOwner                         types.PARTY                              `json:"ccipOwner"`
	Sender                            types.PARTY                              `json:"sender"`
	DestChainSelector                 types.NUMERIC                            `json:"destChainSelector"`
	SequenceNumber                    types.NUMERIC                            `json:"sequenceNumber"`
	RequiredCCVs                      []RawInstanceAddress                     `json:"requiredCCVs"`
	ExecutorAddress                   RawInstanceAddress                       `json:"executorAddress"`
	SourceChainSelector               types.NUMERIC                            `json:"sourceChainSelector"`
	SenderAddress                     types.TEXT                               `json:"senderAddress"`
	Receiver                          types.TEXT                               `json:"receiver"`
	Payload                           types.TEXT                               `json:"payload"`
	ExecutionGasLimit                 types.INT64                              `json:"executionGasLimit"`
	CcipReceiveGasLimit               types.INT64                              `json:"ccipReceiveGasLimit"`
	Finality                          types.INT64                              `json:"finality"`
	CcvAndExecutorHash                types.TEXT                               `json:"ccvAndExecutorHash"`
	OnRampAddress                     types.TEXT                               `json:"onRampAddress"`
	OffRampAddress                    types.TEXT                               `json:"offRampAddress"`
	TokenReceiver                     types.TEXT                               `json:"tokenReceiver"`
	FeeToken                          splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	NetworkFeeUSDCents                types.NUMERIC                            `json:"networkFeeUSDCents"`
	ObservingParties                  []types.PARTY                            `json:"observingParties"`
	CcvFees                           []CCVFee                                 `json:"ccvFees"`
	TokenSendFee                      *TokenSendFee                            `json:"tokenSendFee"`
	FeesFinalized                     types.BOOL                               `json:"feesFinalized"`
	CcvFeeTokenAmounts                []types.NUMERIC                          `json:"ccvFeeTokenAmounts"`
	TokenSendFeeTokenAmount           types.NUMERIC                            `json:"tokenSendFeeTokenAmount"`
	NetworkFeeTokenAmount             types.NUMERIC                            `json:"networkFeeTokenAmount"`
	TokenSendData                     *TokenSendData                           `json:"tokenSendData"`
	VerifierData                      []VerifierData                           `json:"verifierData"`
	Message                           *MessageV1                               `json:"message"`
	EncodedMessage                    types.TEXT                               `json:"encodedMessage"`
	MessageId                         types.TEXT                               `json:"messageId"`
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
	args["routerInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RouterInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OnRampInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
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
	args["executorAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutorAddress).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutorAddress
	}()

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
	args["finality"] = int64(t.Finality)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampAddress"] = string(t.OnRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["offRampAddress"] = string(t.OffRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = string(t.TokenReceiver)

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
	args["feesFinalized"] = bool(t.FeesFinalized)

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

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t SendingMessageV1) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["routerInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RouterInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OnRampInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
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
	args["executorAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutorAddress).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutorAddress
	}()

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
	args["finality"] = int64(t.Finality)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvAndExecutorHash"] = string(t.CcvAndExecutorHash)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampAddress"] = string(t.OnRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["offRampAddress"] = string(t.OffRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = string(t.TokenReceiver)

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
	args["feesFinalized"] = bool(t.FeesFinalized)

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

// SetInboundPoolCCVs is a Record type
type SetInboundPoolCCVs struct {
	PoolCCVs []RawInstanceAddress `json:"poolCCVs"`
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

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        types.BOOL           `json:"isEnabled"`
	OnRampAddress    types.TEXT           `json:"onRampAddress"`
	LaneMandatedCCVs []RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs      []RawInstanceAddress `json:"defaultCCVs"`
}

// ToMap converts SourceChainConfig to a map for DAML arguments
func (t SourceChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["isEnabled"] = bool(t.IsEnabled)

	m["onRampAddress"] = string(t.OnRampAddress)

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

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner                    types.PARTY                              `json:"ccipOwner"`
	TokenAdminRegistryInstanceId types.TEXT                               `json:"tokenAdminRegistryInstanceId"`
	PoolOwner                    types.PARTY                              `json:"poolOwner"`
	Receiver                     types.PARTY                              `json:"receiver"`
	TokenReceiver                types.PARTY                              `json:"tokenReceiver"`
	InstrumentId                 splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount                       types.NUMERIC                            `json:"amount"`
	MessageHash                  types.TEXT                               `json:"messageHash"`
	SourceChainSelector          types.NUMERIC                            `json:"sourceChainSelector"`
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
	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

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
	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != "" {
		args["sourceChainSelector"] = t.SourceChainSelector
	}

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

// UpdateDestChainConfig is a Record type
type UpdateDestChainConfig struct {
	DestChainSelector types.NUMERIC   `json:"destChainSelector"`
	Config            DestChainConfig `json:"config"`
}

// ToMap converts UpdateDestChainConfig to a map for DAML arguments
func (t UpdateDestChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["config"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	return m
}

func (t UpdateDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateDestChainConfig to hex string (Canton MCMS format)
func (t UpdateDestChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateDestChainConfig from hex string (Canton MCMS format)
func (t *UpdateDestChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateSourceChainConfig is a Record type
type UpdateSourceChainConfig struct {
	SourceChainSelector types.NUMERIC     `json:"sourceChainSelector"`
	Config              SourceChainConfig `json:"config"`
}

// ToMap converts UpdateSourceChainConfig to a map for DAML arguments
func (t UpdateSourceChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["config"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	return m
}

func (t UpdateSourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateSourceChainConfig to hex string (Canton MCMS format)
func (t UpdateSourceChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateSourceChainConfig from hex string (Canton MCMS format)
func (t *UpdateSourceChainConfig) UnmarshalHex(data string) error {
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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddCCVFee(args AddCCVFee) (*bind.EncodedChoice, error)
	AddCCVVerification(args AddCCVVerification) (*bind.EncodedChoice, error)
	AddTokenSend(args AddTokenSend) (*bind.EncodedChoice, error)
	AddTokenSendFee(args AddTokenSendFee) (*bind.EncodedChoice, error)
	AddVerifierData(args AddVerifierData) (*bind.EncodedChoice, error)
	ExecutingMessageV1Archive(args ExecutingMessageV1Archive) (*bind.EncodedChoice, error)
	FeeTokenAmount(args FeeTokenAmount) (*bind.EncodedChoice, error)
	FinalizeFee(args FinalizeFee) (*bind.EncodedChoice, error)
	GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error)
	GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error)
	SetInboundPoolCCVs(args SetInboundPoolCCVs) (*bind.EncodedChoice, error)
	UpdateDestChainConfig(args UpdateDestChainConfig) (*bind.EncodedChoice, error)
	UpdateSourceChainConfig(args UpdateSourceChainConfig) (*bind.EncodedChoice, error)
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

// AddCCVVerification encodes parameters for the AddCCVVerification choice.
func (e *encoder) AddCCVVerification(args AddCCVVerification) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCCVVerification", args)
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

// ExecutingMessageV1Archive encodes parameters for the ExecutingMessageV1Archive choice.
func (e *encoder) ExecutingMessageV1Archive(args ExecutingMessageV1Archive) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutingMessageV1Archive", args)
}

// FeeTokenAmount encodes parameters for the FeeTokenAmount choice.
func (e *encoder) FeeTokenAmount(args FeeTokenAmount) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FeeTokenAmount", args)
}

// FinalizeFee encodes parameters for the FinalizeFee choice.
func (e *encoder) FinalizeFee(args FinalizeFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeFee", args)
}

// GetDestChainConfig encodes parameters for the GetDestChainConfig choice.
func (e *encoder) GetDestChainConfig(args GetDestChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChainConfig", args)
}

// GetSourceChainConfig encodes parameters for the GetSourceChainConfig choice.
func (e *encoder) GetSourceChainConfig(args GetSourceChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSourceChainConfig", args)
}

// SetInboundPoolCCVs encodes parameters for the SetInboundPoolCCVs choice.
func (e *encoder) SetInboundPoolCCVs(args SetInboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetInboundPoolCCVs", args)
}

// UpdateDestChainConfig encodes parameters for the UpdateDestChainConfig choice.
func (e *encoder) UpdateDestChainConfig(args UpdateDestChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateDestChainConfig", args)
}

// UpdateSourceChainConfig encodes parameters for the UpdateSourceChainConfig choice.
func (e *encoder) UpdateSourceChainConfig(args UpdateSourceChainConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateSourceChainConfig", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
