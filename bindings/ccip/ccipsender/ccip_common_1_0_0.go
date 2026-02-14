package ccipsender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

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

// AddCCVFee is a Record type
type AddCCVFee struct {
	CcvInstanceId     TEXT    `json:"ccvInstanceId"`
	FeeUSDCents       NUMERIC `json:"feeUSDCents"`
	DestGasLimit      INT64   `json:"destGasLimit"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
	Caller            PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *AddCCVFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AddCCVVerification is a Record type
type AddCCVVerification struct {
	CcvInstanceId TEXT  `json:"ccvInstanceId"`
	VersionTag    TEXT  `json:"versionTag"`
	Caller        PARTY `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *AddCCVVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AddTokenSend is a Record type
type AddTokenSend struct {
	PoolInstanceId   TEXT                 `json:"poolInstanceId"`
	PoolOwner        PARTY                `json:"poolOwner"`
	InstrumentId     InstrumentId         `json:"instrumentId"`
	Amount           NUMERIC              `json:"amount"`
	DestTokenAddress TEXT                 `json:"destTokenAddress"`
	ExtraData        TEXT                 `json:"extraData"`
	PoolRequiredCCVs []RawInstanceAddress `json:"poolRequiredCCVs"`
	Caller           PARTY                `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *AddTokenSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AddTokenSendFee is a Record type
type AddTokenSendFee struct {
	PoolInstanceId    TEXT    `json:"poolInstanceId"`
	PoolOwner         PARTY   `json:"poolOwner"`
	FeeUSDCents       NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   INT64   `json:"destGasOverhead"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
	Caller            PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *AddTokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AddVerifierData is a Record type
type AddVerifierData struct {
	CcvInstanceId        TEXT    `json:"ccvInstanceId"`
	VersionTag           TEXT    `json:"versionTag"`
	VerifierBlob         TEXT    `json:"verifierBlob"`
	MessageSentObservers []PARTY `json:"messageSentObservers"`
	Caller               PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *AddVerifierData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVFee is a Record type
type CCVFee struct {
	CcvInstanceId     TEXT    `json:"ccvInstanceId"`
	CcvOwner          PARTY   `json:"ccvOwner"`
	FeeUSDCents       NUMERIC `json:"feeUSDCents"`
	DestGasLimit      INT64   `json:"destGasLimit"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
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
	return jsonCodec.Marshall(t)
}

func (t *CCVFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVVerification is a Record type
type CCVVerification struct {
	CcvInstanceId TEXT  `json:"ccvInstanceId"`
	CcvOwner      PARTY `json:"ccvOwner"`
	VersionTag    TEXT  `json:"versionTag"`
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
	return jsonCodec.Marshall(t)
}

func (t *CCVVerification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	CcipOwner       PARTY `json:"ccipOwner"`
	StorageLocation TEXT  `json:"storageLocation"`
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
	return jsonCodec.Marshall(t)
}

func (t *CrossChainVerifierView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierCalculateFee is a Record type
type CrossChainVerifierCalculateFee struct {
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	Caller            PARTY       `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *CrossChainVerifierCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierForwardToVerifier is a Record type
type CrossChainVerifierForwardToVerifier struct {
	RmnRemoteCid      CONTRACT_ID `json:"rmnRemoteCid"`
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	VerifierArgs      TEXT        `json:"verifierArgs"`
	Caller            PARTY       `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *CrossChainVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierVerifyMessage is a Record type
type CrossChainVerifierVerifyMessage struct {
	RmnRemoteCid        CONTRACT_ID `json:"rmnRemoteCid"`
	ExecutingMessageCid CONTRACT_ID `json:"executingMessageCid"`
	VerifierResults     TEXT        `json:"verifierResults"`
	Caller              PARTY       `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *CrossChainVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfig is a Record type
type DestChainConfig struct {
	IsEnabled                 BOOL                 `json:"isEnabled"`
	DefaultExecutor           RawInstanceAddress   `json:"defaultExecutor"`
	OffRampAddress            TEXT                 `json:"offRampAddress"`
	LaneMandatedCCVs          []RawInstanceAddress `json:"laneMandatedCCVs"`
	DefaultCCVs               []RawInstanceAddress `json:"defaultCCVs"`
	MessageNetworkFeeUSDCents NUMERIC              `json:"messageNetworkFeeUSDCents"`
	TokenNetworkFeeUSDCents   NUMERIC              `json:"tokenNetworkFeeUSDCents"`
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
	return jsonCodec.Marshall(t)
}

func (t *DestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecutingMessageV1 is a Template type
type ExecutingMessageV1 struct {
	CcipOwner                         PARTY                 `json:"ccipOwner"`
	Message                           MessageV1             `json:"message"`
	MessageId                         TEXT                  `json:"messageId"`
	Receiver                          PARTY                 `json:"receiver"`
	TokenReceiver                     *PARTY                `json:"tokenReceiver"`
	Executor                          PARTY                 `json:"executor"`
	ObservingParties                  []PARTY               `json:"observingParties"`
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
	return jsonCodec.Marshall(t)
}

func (t *ExecutingMessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	Caller PARTY `json:"caller"`
}

// ToMap converts ExecutingMessageV1Archive to a map for DAML arguments
func (t ExecutingMessageV1Archive) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutingMessageV1Archive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ExecutingMessageV1Archive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FeeTokenAmount is a Record type
type FeeTokenAmount struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts FeeTokenAmount to a map for DAML arguments
func (t FeeTokenAmount) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t FeeTokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *FeeTokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// FinalizeFee is a Record type
type FinalizeFee struct {
	FeeTokenPrice     NUMERIC `json:"feeTokenPrice"`
	PremiumMultiplier NUMERIC `json:"premiumMultiplier"`
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
	return jsonCodec.Marshall(t)
}

func (t *FinalizeFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestChainConfig is a Record type
type GetDestChainConfig struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *GetDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetSourceChainConfig is a Record type
type GetSourceChainConfig struct {
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
	Caller              PARTY   `json:"caller"`
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
	return jsonCodec.Marshall(t)
}

func (t *GetSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GlobalConfig is a Template type
type GlobalConfig struct {
	InstanceId         TEXT    `json:"instanceId"`
	CcipOwner          PARTY   `json:"ccipOwner"`
	ChainSelector      NUMERIC `json:"chainSelector"`
	OnRampAddress      TEXT    `json:"onRampAddress"`
	DestChainConfigs   GENMAP  `json:"destChainConfigs"`
	SourceChainConfigs GENMAP  `json:"sourceChainConfigs"`
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
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceChainConfigs"] = func() any {
		if t.SourceChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
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
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceChainConfigs"] = func() any {
		if t.SourceChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": GENMAP{}}
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
	return jsonCodec.Marshall(t)
}

func (t *GlobalConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	return jsonCodec.Marshall(e)
}

func (e *IssuerType) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

var _ ENUM = IssuerType("")

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
	return jsonCodec.Marshall(e)
}

func (e *MessageExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

var _ ENUM = MessageExecutionState("")

// MessageV1 is a Record type
type MessageV1 struct {
	SourceChainSelector NUMERIC          `json:"sourceChainSelector"`
	DestChainSelector   NUMERIC          `json:"destChainSelector"`
	SequenceNumber      NUMERIC          `json:"sequenceNumber"`
	ExecutionGasLimit   INT64            `json:"executionGasLimit"`
	CcipReceiveGasLimit INT64            `json:"ccipReceiveGasLimit"`
	Finality            INT64            `json:"finality"`
	CcvAndExecutorHash  TEXT             `json:"ccvAndExecutorHash"`
	OnRampAddress       TEXT             `json:"onRampAddress"`
	OffRampAddress      TEXT             `json:"offRampAddress"`
	Sender              TEXT             `json:"sender"`
	Receiver            TEXT             `json:"receiver"`
	DestBlob            TEXT             `json:"destBlob"`
	TokenTransfer       *TokenTransferV1 `json:"tokenTransfer"`
	MessageData         TEXT             `json:"messageData"`
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
	return jsonCodec.Marshall(t)
}

func (t *MessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RawInstanceAddress is a Record type
type RawInstanceAddress struct {
	Unpack TEXT `json:"unpack"`
}

// ToMap converts RawInstanceAddress to a map for DAML arguments
func (t RawInstanceAddress) ToMap() map[string]any {
	m := make(map[string]any)

	m["unpack"] = string(t.Unpack)

	return m
}

func (t RawInstanceAddress) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *RawInstanceAddress) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Receipt is a Record type
type Receipt struct {
	IssuerType        IssuerType `json:"issuerType"`
	IssuerAddress     TEXT       `json:"issuerAddress"`
	VersionTag        *TEXT      `json:"versionTag"`
	DestGasLimit      INT64      `json:"destGasLimit"`
	DestBytesOverhead INT64      `json:"destBytesOverhead"`
	FeeTokenAmount    NUMERIC    `json:"feeTokenAmount"`
	ExtraArgs         TEXT       `json:"extraArgs"`
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
	return jsonCodec.Marshall(t)
}

func (t *Receipt) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SendingMessageV1 is a Template type
type SendingMessageV1 struct {
	RouterInstanceAddress             RawInstanceAddress   `json:"routerInstanceAddress"`
	OnRampInstanceAddress             RawInstanceAddress   `json:"onRampInstanceAddress"`
	TokenAdminRegistryInstanceAddress RawInstanceAddress   `json:"tokenAdminRegistryInstanceAddress"`
	CcipOwner                         PARTY                `json:"ccipOwner"`
	Sender                            PARTY                `json:"sender"`
	DestChainSelector                 NUMERIC              `json:"destChainSelector"`
	SequenceNumber                    NUMERIC              `json:"sequenceNumber"`
	RequiredCCVs                      []RawInstanceAddress `json:"requiredCCVs"`
	ExecutorAddress                   RawInstanceAddress   `json:"executorAddress"`
	SourceChainSelector               NUMERIC              `json:"sourceChainSelector"`
	SenderAddress                     TEXT                 `json:"senderAddress"`
	Receiver                          TEXT                 `json:"receiver"`
	Payload                           TEXT                 `json:"payload"`
	ExecutionGasLimit                 INT64                `json:"executionGasLimit"`
	CcipReceiveGasLimit               INT64                `json:"ccipReceiveGasLimit"`
	Finality                          INT64                `json:"finality"`
	CcvAndExecutorHash                TEXT                 `json:"ccvAndExecutorHash"`
	OnRampAddress                     TEXT                 `json:"onRampAddress"`
	OffRampAddress                    TEXT                 `json:"offRampAddress"`
	TokenReceiver                     TEXT                 `json:"tokenReceiver"`
	FeeToken                          InstrumentId         `json:"feeToken"`
	NetworkFeeUSDCents                NUMERIC              `json:"networkFeeUSDCents"`
	ObservingParties                  []PARTY              `json:"observingParties"`
	CcvFees                           []CCVFee             `json:"ccvFees"`
	TokenSendFee                      *TokenSendFee        `json:"tokenSendFee"`
	FeesFinalized                     BOOL                 `json:"feesFinalized"`
	CcvFeeTokenAmounts                []NUMERIC            `json:"ccvFeeTokenAmounts"`
	TokenSendFeeTokenAmount           NUMERIC              `json:"tokenSendFeeTokenAmount"`
	NetworkFeeTokenAmount             NUMERIC              `json:"networkFeeTokenAmount"`
	TokenSendData                     *TokenSendData       `json:"tokenSendData"`
	VerifierData                      []VerifierData       `json:"verifierData"`
	Message                           *MessageV1           `json:"message"`
	EncodedMessage                    TEXT                 `json:"encodedMessage"`
	MessageId                         TEXT                 `json:"messageId"`
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
	return jsonCodec.Marshall(t)
}

func (t *SendingMessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	return jsonCodec.Marshall(t)
}

func (t *SetInboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        BOOL                 `json:"isEnabled"`
	OnRampAddress    TEXT                 `json:"onRampAddress"`
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
	return jsonCodec.Marshall(t)
}

func (t *SourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner                    PARTY        `json:"ccipOwner"`
	TokenAdminRegistryInstanceId TEXT         `json:"tokenAdminRegistryInstanceId"`
	PoolOwner                    PARTY        `json:"poolOwner"`
	Receiver                     PARTY        `json:"receiver"`
	TokenReceiver                PARTY        `json:"tokenReceiver"`
	InstrumentId                 InstrumentId `json:"instrumentId"`
	Amount                       NUMERIC      `json:"amount"`
	MessageHash                  TEXT         `json:"messageHash"`
	SourceChainSelector          NUMERIC      `json:"sourceChainSelector"`
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
	return jsonCodec.Marshall(t)
}

func (t *TokenReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	PoolInstanceId   TEXT         `json:"poolInstanceId"`
	PoolOwner        PARTY        `json:"poolOwner"`
	InstrumentId     InstrumentId `json:"instrumentId"`
	Amount           NUMERIC      `json:"amount"`
	DestTokenAddress TEXT         `json:"destTokenAddress"`
	ExtraData        TEXT         `json:"extraData"`
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
	return jsonCodec.Marshall(t)
}

func (t *TokenSendData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenSendFee is a Record type
type TokenSendFee struct {
	PoolInstanceId    TEXT    `json:"poolInstanceId"`
	PoolOwner         PARTY   `json:"poolOwner"`
	FeeUSDCents       NUMERIC `json:"feeUSDCents"`
	DestGasOverhead   INT64   `json:"destGasOverhead"`
	DestBytesOverhead INT64   `json:"destBytesOverhead"`
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
	return jsonCodec.Marshall(t)
}

func (t *TokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenTransferV1 is a Record type
type TokenTransferV1 struct {
	Amount             NUMERIC `json:"amount"`
	SourcePoolAddress  TEXT    `json:"sourcePoolAddress"`
	SourceTokenAddress TEXT    `json:"sourceTokenAddress"`
	DestTokenAddress   TEXT    `json:"destTokenAddress"`
	TokenReceiver      TEXT    `json:"tokenReceiver"`
	ExtraData          TEXT    `json:"extraData"`
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
	return jsonCodec.Marshall(t)
}

func (t *TokenTransferV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateDestChainConfig is a Record type
type UpdateDestChainConfig struct {
	DestChainSelector NUMERIC         `json:"destChainSelector"`
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
	return jsonCodec.Marshall(t)
}

func (t *UpdateDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateSourceChainConfig is a Record type
type UpdateSourceChainConfig struct {
	SourceChainSelector NUMERIC           `json:"sourceChainSelector"`
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
	return jsonCodec.Marshall(t)
}

func (t *UpdateSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// VerifierData is a Record type
type VerifierData struct {
	CcvInstanceId        TEXT    `json:"ccvInstanceId"`
	CcvOwner             PARTY   `json:"ccvOwner"`
	VersionTag           TEXT    `json:"versionTag"`
	VerifierBlob         TEXT    `json:"verifierBlob"`
	MessageSentObservers []PARTY `json:"messageSentObservers"`
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
	return jsonCodec.Marshall(t)
}

func (t *VerifierData) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IICrossChainVerifierInterfaceID returns the interface ID for the IICrossChainVerifier interface using the package name
func IICrossChainVerifierInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.CrossChainVerifier", "ICrossChainVerifier")
}

// IICrossChainVerifierInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IICrossChainVerifierInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Interfaces.CrossChainVerifier", "ICrossChainVerifier")
}
