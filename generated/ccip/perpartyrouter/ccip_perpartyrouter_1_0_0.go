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

const PackageID = "ebb099f0b51b8ce01d2e15c1889149a6717674032cc4b1adbfb4f4ebb0961758"
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

	// Check if the type has a toMap method
	type mapper interface {
		toMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.toMap()
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
	DestChainSelector NUMERIC `json:"destChainSelector"`
	SequenceNumber    NUMERIC `json:"sequenceNumber"`
	MessageId         TEXT    `json:"messageId"`
	EncodedMessage    TEXT    `json:"encodedMessage"`
	VerifierBlobs     []TEXT  `json:"verifierBlobs"`
}

// toMap converts CCIPMessageSentEvent to a map for DAML arguments
func (t CCIPMessageSentEvent) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"sequenceNumber":    (*big.Int)(t.SequenceNumber),
		"messageId":         string(t.MessageId),
		"encodedMessage":    string(t.EncodedMessage),
		"verifierBlobs": func() []interface{} {
			res := make([]interface{}, 0, len(t.VerifierBlobs))
			for _, e := range t.VerifierBlobs {
				res = append(res, string(e))
			}
			return res
		}(),
	}
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

// toMap converts CCIPSend to a map for DAML arguments
func (t CCIPSend) toMap() map[string]interface{} {
	return map[string]interface{}{

		"onRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OnRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OnRampCid
		}(),
		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"destChainSelector":   (*big.Int)(t.DestChainSelector),
		"receiver":            string(t.Receiver),
		"payload":             string(t.Payload),
		"executionGasLimit":   int64(t.ExecutionGasLimit),
		"ccipReceiveGasLimit": int64(t.CcipReceiveGasLimit),
		"tokenSendTicket":     *t.TokenSendTicket,
		"ccvTickets": func() []interface{} {
			res := make([]interface{}, 0, len(t.CcvTickets))
			for _, e := range t.CcvTickets {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts CCIPSendResult to a map for DAML arguments
func (t CCIPSendResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"router": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Router).(mapper); ok {
				return m.toMap()
			}
			return t.Router
		}(),
		"ccipMessageSent": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.CcipMessageSent).(mapper); ok {
				return m.toMap()
			}
			return t.CcipMessageSent
		}(),
		"messageId": string(t.MessageId),
	}
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

// toMap converts CreateRouter to a map for DAML arguments
func (t CreateRouter) toMap() map[string]interface{} {
	return map[string]interface{}{

		"partyOwner": t.PartyOwner.ToMap(),
		"instanceId": string(t.InstanceId),
	}
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

// toMap converts CreateRouterResult to a map for DAML arguments
func (t CreateRouterResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"router": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Router).(mapper); ok {
				return m.toMap()
			}
			return t.Router
		}(),
		"factory": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Factory).(mapper); ok {
				return m.toMap()
			}
			return t.Factory
		}(),
	}
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
	OffRampCid            CONTRACT_ID   `json:"offRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	EncodedMessage        TEXT          `json:"encodedMessage"`
	CcvVerifyTickets      []CONTRACT_ID `json:"ccvVerifyTickets"`
}

// toMap converts Execute to a map for DAML arguments
func (t Execute) toMap() map[string]interface{} {
	return map[string]interface{}{

		"offRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OffRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OffRampCid
		}(),
		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"encodedMessage": string(t.EncodedMessage),
		"ccvVerifyTickets": func() []interface{} {
			res := make([]interface{}, 0, len(t.CcvVerifyTickets))
			for _, e := range t.CcvVerifyTickets {
				res = append(res, e)
			}
			return res
		}(),
	}
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
	TokenReceiveTicket CONTRACT_ID           `json:"tokenReceiveTicket"`
	MessageId          TEXT                  `json:"messageId"`
	Message            MessageV1             `json:"message"`
	State              MessageExecutionState `json:"state"`
}

// toMap converts ExecuteResult to a map for DAML arguments
func (t ExecuteResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"router": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Router).(mapper); ok {
				return m.toMap()
			}
			return t.Router
		}(),
		"tokenReceiveTicket": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenReceiveTicket).(mapper); ok {
				return m.toMap()
			}
			return t.TokenReceiveTicket
		}(),
		"messageId": string(t.MessageId),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"state": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}(),
	}
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

// toMap converts ExecutionStateChangedEvent to a map for DAML arguments
func (t ExecutionStateChangedEvent) toMap() map[string]interface{} {
	return map[string]interface{}{

		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"messageId":           string(t.MessageId),
		"state": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.State).(mapper); ok {
				return m.toMap()
			}
			return t.State
		}(),
		"returnData": string(t.ReturnData),
	}
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

// toMap converts GetExecutionState to a map for DAML arguments
func (t GetExecutionState) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageHash": string(t.MessageHash),
	}
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

// GetReceiverRequiredCCVs is a Record type
type GetReceiverRequiredCCVs struct {
	Caller PARTY `json:"caller"`
}

// toMap converts GetReceiverRequiredCCVs to a map for DAML arguments
func (t GetReceiverRequiredCCVs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetReceiverRequiredCCVs using JsonCodec
func (t GetReceiverRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetReceiverRequiredCCVs using JsonCodec
func (t *GetReceiverRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	OffRampCid            CONTRACT_ID   `json:"offRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	SourceChainSelector   NUMERIC       `json:"sourceChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// toMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) toMap() map[string]interface{} {
	return map[string]interface{}{

		"offRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OffRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OffRampCid
		}(),
		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"hasTokenTransfer":    bool(t.HasTokenTransfer),
		"instrumentId":        *t.InstrumentId,
	}
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForExecute using JsonCodec
func (t GetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForExecute using JsonCodec
func (t *GetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	OnRampCid             CONTRACT_ID   `json:"onRampCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	DestChainSelector     NUMERIC       `json:"destChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// toMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) toMap() map[string]interface{} {
	return map[string]interface{}{

		"onRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OnRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OnRampCid
		}(),
		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"hasTokenTransfer":  bool(t.HasTokenTransfer),
		"instrumentId":      *t.InstrumentId,
	}
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForSend using JsonCodec
func (t GetRequiredCCVsForSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForSend using JsonCodec
func (t *GetRequiredCCVsForSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetSequenceNumber is a Record type
type GetSequenceNumber struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
}

// toMap converts GetSequenceNumber to a map for DAML arguments
func (t GetSequenceNumber) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
	}
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

// toMap converts HasRouter to a map for DAML arguments
func (t HasRouter) toMap() map[string]interface{} {
	return map[string]interface{}{

		"partyOwner": t.PartyOwner.ToMap(),
		"caller":     t.Caller.ToMap(),
	}
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
	ReceiverRequiredCCVs    []TEXT `json:"receiverRequiredCCVs"`
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

	if t.OutboundSequenceNumbers != nil && len(t.OutboundSequenceNumbers) > 0 {
		args["outboundSequenceNumbers"] = map[string]interface{}{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}

	if t.ExecutionStates != nil && len(t.ExecutionStates) > 0 {
		args["executionStates"] = map[string]interface{}{"_type": "genmap", "value": t.ExecutionStates}
	}

	if len(t.ReceiverRequiredCCVs) > 0 {
		args["receiverRequiredCCVs"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
			for _, e := range t.ReceiverRequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}()
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

// GetReceiverRequiredCCVs exercises the GetReceiverRequiredCCVs choice on this PerPartyRouter contract
func (t PerPartyRouter) GetReceiverRequiredCCVs(contractID string, args GetReceiverRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetReceiverRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// UpdateReceiverRequiredCCVs exercises the UpdateReceiverRequiredCCVs choice on this PerPartyRouter contract
func (t PerPartyRouter) UpdateReceiverRequiredCCVs(contractID string, args UpdateReceiverRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "UpdateReceiverRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
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

	if t.RegisteredRouters != nil && len(t.RegisteredRouters) > 0 {
		args["registeredRouters"] = map[string]interface{}{"_type": "genmap", "value": t.RegisteredRouters}
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

// UpdateReceiverRequiredCCVs is a Record type
type UpdateReceiverRequiredCCVs struct {
	NewReceiverRequiredCCVs []TEXT `json:"newReceiverRequiredCCVs"`
}

// toMap converts UpdateReceiverRequiredCCVs to a map for DAML arguments
func (t UpdateReceiverRequiredCCVs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"newReceiverRequiredCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.NewReceiverRequiredCCVs))
			for _, e := range t.NewReceiverRequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for UpdateReceiverRequiredCCVs using JsonCodec
func (t UpdateReceiverRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UpdateReceiverRequiredCCVs using JsonCodec
func (t *UpdateReceiverRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
