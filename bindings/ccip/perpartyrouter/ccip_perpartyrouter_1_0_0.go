package perpartyrouter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "7fe78b8b1675aeb1bc043235d3123a0d4303f4fbab9813360de880b3b1c08c96"
const SDKVersion = "3.4.8"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a ToMap method
	type mapper interface {
		ToMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// CCIPMessageSent is a Template type
type CCIPMessageSent struct {
	CcipOwner PARTY                `json:"ccipOwner"`
	Sender    PARTY                `json:"sender"`
	Observers []PARTY              `json:"observers"`
	Event     CCIPMessageSentEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template
func (t CCIPMessageSent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "CCIPMessageSent")
}

// CreateCommand returns a CreateCommand for this template
func (t CCIPMessageSent) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["sender"] = t.Sender.ToMap()

	if len(t.Observers) > 0 {
		args["observers"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.Observers))
			for _, e := range t.Observers {
				res = append(res, e.ToMap())
			}
			return res
		}()
	}

	args["event"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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

// MarshalJSON implements custom JSON marshaling for CCIPMessageSent using JsonCodec
func (t CCIPMessageSent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPMessageSent using JsonCodec
func (t *CCIPMessageSent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCIPMessageSent

// Archive exercises the Archive choice on this CCIPMessageSent contract
func (t CCIPMessageSent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCIPMessageSentEvent is a Record type
type CCIPMessageSentEvent struct {
	DestChainSelector NUMERIC   `json:"destChainSelector"`
	SequenceNumber    NUMERIC   `json:"sequenceNumber"`
	MessageId         TEXT      `json:"messageId"`
	EncodedMessage    TEXT      `json:"encodedMessage"`
	VerifierBlobs     []TEXT    `json:"verifierBlobs"`
	Receipts          []Receipt `json:"receipts"`
}

// ToMap converts CCIPMessageSentEvent to a map for DAML arguments
func (t CCIPMessageSentEvent) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["sequenceNumber"] = (*big.Int)(t.SequenceNumber)

	m["messageId"] = string(t.MessageId)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["verifierBlobs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["receipts"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Receipts))
		for _, e := range t.Receipts {
			type mapper interface{ toMap() map[string]interface{} }
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

// MarshalJSON implements custom JSON marshaling for CCIPMessageSentEvent using JsonCodec
func (t CCIPMessageSentEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPMessageSentEvent using JsonCodec
func (t *CCIPMessageSentEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCIPSend is a Record type
type CCIPSend struct {
	OnRampCid             CONTRACT_ID   `json:"onRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	DestChainSelector     NUMERIC       `json:"destChainSelector"`
	Receiver              TEXT          `json:"receiver"`
	Payload               TEXT          `json:"payload"`
	ExecutionGasLimit     INT64         `json:"executionGasLimit"`
	CcipReceiveGasLimit   INT64         `json:"ccipReceiveGasLimit"`
	TokenSendTicket       *CONTRACT_ID  `json:"tokenSendTicket"`
	CcvTickets            []CONTRACT_ID `json:"ccvTickets"`
}

// ToMap converts CCIPSend to a map for DAML arguments
func (t CCIPSend) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["onRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OnRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampCid
	}()

	m["globalConfigCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["executionGasLimit"] = int64(t.ExecutionGasLimit)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	if t.TokenSendTicket != nil {
		m["tokenSendTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenSendTicket,
		}
	} else {
		m["tokenSendTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["ccvTickets"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvTickets))
		for _, e := range t.CcvTickets {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for CCIPSend using JsonCodec
func (t CCIPSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPSend using JsonCodec
func (t *CCIPSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCIPSendResult is a Record type
type CCIPSendResult struct {
	Router          CONTRACT_ID `json:"router"`
	CcipMessageSent CONTRACT_ID `json:"ccipMessageSent"`
	MessageId       TEXT        `json:"messageId"`
}

// ToMap converts CCIPSendResult to a map for DAML arguments
func (t CCIPSendResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["router"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	m["ccipMessageSent"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.CcipMessageSent).(mapper); ok {
			return m.toMap()
		}
		return t.CcipMessageSent
	}()

	m["messageId"] = string(t.MessageId)

	return m
}

// MarshalJSON implements custom JSON marshaling for CCIPSendResult using JsonCodec
func (t CCIPSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPSendResult using JsonCodec
func (t *CCIPSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CreateRouter is a Record type
type CreateRouter struct {
	PartyOwner PARTY `json:"partyOwner"`
	InstanceId TEXT  `json:"instanceId"`
}

// ToMap converts CreateRouter to a map for DAML arguments
func (t CreateRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["partyOwner"] = t.PartyOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

// MarshalJSON implements custom JSON marshaling for CreateRouter using JsonCodec
func (t CreateRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CreateRouter using JsonCodec
func (t *CreateRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CreateRouterResult is a Record type
type CreateRouterResult struct {
	Router  CONTRACT_ID `json:"router"`
	Factory CONTRACT_ID `json:"factory"`
}

// ToMap converts CreateRouterResult to a map for DAML arguments
func (t CreateRouterResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["router"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	m["factory"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Factory).(mapper); ok {
			return m.toMap()
		}
		return t.Factory
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for CreateRouterResult using JsonCodec
func (t CreateRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CreateRouterResult using JsonCodec
func (t *CreateRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Execute is a Record type
type Execute struct {
	OffRampCid             CONTRACT_ID   `json:"offRampCid"`
	GlobalConfigCid        CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid  CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	EncodedMessage         TEXT          `json:"encodedMessage"`
	CcvVerifyTickets       []CONTRACT_ID `json:"ccvVerifyTickets"`
	TokenPoolCCVTicket     *CONTRACT_ID  `json:"tokenPoolCCVTicket"`
	ReceiverRequiredCCVIds []TEXT        `json:"receiverRequiredCCVIds"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["offRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OffRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampCid
	}()

	m["globalConfigCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["encodedMessage"] = string(t.EncodedMessage)

	m["ccvVerifyTickets"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvVerifyTickets))
		for _, e := range t.CcvVerifyTickets {
			res = append(res, e)
		}
		return res
	}()

	if t.TokenPoolCCVTicket != nil {
		m["tokenPoolCCVTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenPoolCCVTicket,
		}
	} else {
		m["tokenPoolCCVTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["receiverRequiredCCVIds"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVIds))
		for _, e := range t.ReceiverRequiredCCVIds {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for Execute using JsonCodec
func (t Execute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Execute using JsonCodec
func (t *Execute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecuteResult is a Record type
type ExecuteResult struct {
	Router             CONTRACT_ID           `json:"router"`
	TokenReceiveTicket *CONTRACT_ID          `json:"tokenReceiveTicket"`
	MessageId          TEXT                  `json:"messageId"`
	Message            MessageV1             `json:"message"`
	State              MessageExecutionState `json:"state"`
}

// ToMap converts ExecuteResult to a map for DAML arguments
func (t ExecuteResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["router"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenReceiveTicket,
		}
	} else {
		m["tokenReceiveTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["messageId"] = string(t.MessageId)

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["state"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.State).(mapper); ok {
			return m.toMap()
		}
		return t.State
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ExecuteResult using JsonCodec
func (t ExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecuteResult using JsonCodec
func (t *ExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecutionStateChanged is a Template type
type ExecutionStateChanged struct {
	CcipOwner PARTY                      `json:"ccipOwner"`
	Receiver  PARTY                      `json:"receiver"`
	Event     ExecutionStateChangedEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template
func (t ExecutionStateChanged) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "ExecutionStateChanged")
}

// CreateCommand returns a CreateCommand for this template
func (t ExecutionStateChanged) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["receiver"] = t.Receiver.ToMap()

	args["event"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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

// MarshalJSON implements custom JSON marshaling for ExecutionStateChanged using JsonCodec
func (t ExecutionStateChanged) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecutionStateChanged using JsonCodec
func (t *ExecutionStateChanged) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ExecutionStateChanged

// Archive exercises the Archive choice on this ExecutionStateChanged contract
func (t ExecutionStateChanged) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ExecutionStateChangedEvent is a Record type
type ExecutionStateChangedEvent struct {
	SourceChainSelector NUMERIC               `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC               `json:"sequenceNumber"`
	MessageId           TEXT                  `json:"messageId"`
	State               MessageExecutionState `json:"state"`
	ReturnData          TEXT                  `json:"returnData"`
}

// ToMap converts ExecutionStateChangedEvent to a map for DAML arguments
func (t ExecutionStateChangedEvent) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["sequenceNumber"] = (*big.Int)(t.SequenceNumber)

	m["messageId"] = string(t.MessageId)

	m["state"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.State).(mapper); ok {
			return m.toMap()
		}
		return t.State
	}()

	m["returnData"] = string(t.ReturnData)

	return m
}

// MarshalJSON implements custom JSON marshaling for ExecutionStateChangedEvent using JsonCodec
func (t ExecutionStateChangedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecutionStateChangedEvent using JsonCodec
func (t *ExecutionStateChangedEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetExecutionState is a Record type
type GetExecutionState struct {
	MessageHash TEXT `json:"messageHash"`
}

// ToMap converts GetExecutionState to a map for DAML arguments
func (t GetExecutionState) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["messageHash"] = string(t.MessageHash)

	return m
}

// MarshalJSON implements custom JSON marshaling for GetExecutionState using JsonCodec
func (t GetExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetExecutionState using JsonCodec
func (t *GetExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForExecute2 is a Record type
type GetRequiredCCVsForExecute2 struct {
	OffRampCid            CONTRACT_ID   `json:"offRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	ReceiverRequiredCCVs  []TEXT        `json:"receiverRequiredCCVs"`
	SourceChainSelector   NUMERIC       `json:"sourceChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// ToMap converts GetRequiredCCVsForExecute2 to a map for DAML arguments
func (t GetRequiredCCVsForExecute2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["offRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OffRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampCid
	}()

	m["globalConfigCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["hasTokenTransfer"] = bool(t.HasTokenTransfer)

	if t.InstrumentId != nil {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.InstrumentId,
		}
	} else {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForExecute2 using JsonCodec
func (t GetRequiredCCVsForExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForExecute2 using JsonCodec
func (t *GetRequiredCCVsForExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForSend2 is a Record type
type GetRequiredCCVsForSend2 struct {
	OnRampCid             CONTRACT_ID   `json:"onRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	DestChainSelector     NUMERIC       `json:"destChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// ToMap converts GetRequiredCCVsForSend2 to a map for DAML arguments
func (t GetRequiredCCVsForSend2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["onRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OnRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampCid
	}()

	m["globalConfigCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["hasTokenTransfer"] = bool(t.HasTokenTransfer)

	if t.InstrumentId != nil {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.InstrumentId,
		}
	} else {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForSend2 using JsonCodec
func (t GetRequiredCCVsForSend2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForSend2 using JsonCodec
func (t *GetRequiredCCVsForSend2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetSequenceNumber is a Record type
type GetSequenceNumber struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
}

// ToMap converts GetSequenceNumber to a map for DAML arguments
func (t GetSequenceNumber) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	return m
}

// MarshalJSON implements custom JSON marshaling for GetSequenceNumber using JsonCodec
func (t GetSequenceNumber) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetSequenceNumber using JsonCodec
func (t *GetSequenceNumber) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// HasRouter is a Record type
type HasRouter struct {
	PartyOwner PARTY `json:"partyOwner"`
	Caller     PARTY `json:"caller"`
}

// ToMap converts HasRouter to a map for DAML arguments
func (t HasRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["partyOwner"] = t.PartyOwner.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for HasRouter using JsonCodec
func (t HasRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for HasRouter using JsonCodec
func (t *HasRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PerPartyRouter is a Template type
type PerPartyRouter struct {
	CcipOwner               PARTY  `json:"ccipOwner"`
	PartyOwner              PARTY  `json:"partyOwner"`
	InstanceId              TEXT   `json:"instanceId"`
	OutboundSequenceNumbers GENMAP `json:"outboundSequenceNumbers"`
	ExecutionStates         GENMAP `json:"executionStates"`
}

// GetTemplateID returns the template ID for this template
func (t PerPartyRouter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter")
}

// CreateCommand returns a CreateCommand for this template
func (t PerPartyRouter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()
	args["partyOwner"] = t.PartyOwner.ToMap()
	args["instanceId"] = string(t.InstanceId)

	// Non-optional GENMAP fields must be present even if empty
	if t.OutboundSequenceNumbers == nil {
		t.OutboundSequenceNumbers = GENMAP{}
	}
	args["outboundSequenceNumbers"] = map[string]interface{}{
		"_type": "genmap",
		"value": t.OutboundSequenceNumbers,
	}

	if t.ExecutionStates == nil {
		t.ExecutionStates = GENMAP{}
	}
	args["executionStates"] = map[string]interface{}{
		"_type": "genmap",
		"value": t.ExecutionStates,
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for PerPartyRouter using JsonCodec
func (t PerPartyRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PerPartyRouter using JsonCodec
func (t *PerPartyRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for PerPartyRouter

// Archive exercises the Archive choice on this PerPartyRouter contract
func (t PerPartyRouter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCIPSend exercises the CCIPSend choice on this PerPartyRouter contract
func (t PerPartyRouter) CCIPSend(contractID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// Execute exercises the Execute choice on this PerPartyRouter contract
func (t PerPartyRouter) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionState exercises the GetExecutionState choice on this PerPartyRouter contract
func (t PerPartyRouter) GetExecutionState(contractID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumber exercises the GetSequenceNumber choice on this PerPartyRouter contract
func (t PerPartyRouter) GetSequenceNumber(contractID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterFactory is a Template type
type PerPartyRouterFactory struct {
	CcipOwner         PARTY  `json:"ccipOwner"`
	InstanceId        TEXT   `json:"instanceId"`
	RegisteredRouters GENMAP `json:"registeredRouters"`
}

// GetTemplateID returns the template ID for this template
func (t PerPartyRouterFactory) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory")
}

// CreateCommand returns a CreateCommand for this template
func (t PerPartyRouterFactory) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()
	args["instanceId"] = string(t.InstanceId)

	// Non-optional GENMAP must always be present even if empty
	if t.RegisteredRouters == nil {
		t.RegisteredRouters = GENMAP{} // or types.GENMAP{} depending on package
	}
	args["registeredRouters"] = map[string]interface{}{
		"_type": "genmap",
		"value": t.RegisteredRouters,
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for PerPartyRouterFactory using JsonCodec
func (t PerPartyRouterFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for PerPartyRouterFactory using JsonCodec
func (t *PerPartyRouterFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for PerPartyRouterFactory

// Archive exercises the Archive choice on this PerPartyRouterFactory contract
func (t PerPartyRouterFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CreateRouter exercises the CreateRouter choice on this PerPartyRouterFactory contract
func (t PerPartyRouterFactory) CreateRouter(contractID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouter exercises the HasRouter choice on this PerPartyRouterFactory contract
func (t PerPartyRouterFactory) HasRouter(contractID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}
