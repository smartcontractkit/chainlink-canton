package ccipreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

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
)

// CCIPMessageSent is a Template type
type CCIPMessageSent struct {
	CcipOwner PARTY                `json:"ccipOwner"`
	Sender    PARTY                `json:"sender"`
	Observers []PARTY              `json:"observers"`
	Event     CCIPMessageSentEvent `json:"event"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPMessageSent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "CCIPMessageSent")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPMessageSent) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "CCIPMessageSent")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPMessageSent) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
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

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPMessageSent) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *CCIPMessageSent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCIPMessageSent

// Archive exercises the Archive choice on this CCIPMessageSent contract
// This method uses the package name in the template ID
func (t CCIPMessageSent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "CCIPMessageSent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPMessageSent) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "CCIPMessageSent"),
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

	m["destChainSelector"] = t.DestChainSelector

	m["sequenceNumber"] = t.SequenceNumber

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

func (t CCIPMessageSentEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCIPMessageSentEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCIPSend is a Record type
type CCIPSend struct {
	OnRampCid             CONTRACT_ID   `json:"onRampCid"`
	RmnRemoteCid          CONTRACT_ID   `json:"rmnRemoteCid"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	SendingMessageCid     CONTRACT_ID   `json:"sendingMessageCid"`
	FeeTokenInput         TokenInput    `json:"feeTokenInput"`
	FeeTokenHoldingCids   []CONTRACT_ID `json:"feeTokenHoldingCids"`
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

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["feeTokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeTokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.FeeTokenInput
	}()

	m["feeTokenHoldingCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCIPSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCIPSendResult is a Record type
type CCIPSendResult struct {
	Router          CONTRACT_ID   `json:"router"`
	CcipMessageSent CONTRACT_ID   `json:"ccipMessageSent"`
	MessageId       TEXT          `json:"messageId"`
	FeeChangeCids   []CONTRACT_ID `json:"feeChangeCids"`
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

	m["feeChangeCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.FeeChangeCids))
		for _, e := range t.FeeChangeCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCIPSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CancelSend is a Record type
type CancelSend struct {
	OnRampCid             CONTRACT_ID `json:"onRampCid"`
	SendingMessageCid     CONTRACT_ID `json:"sendingMessageCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
}

// ToMap converts CancelSend to a map for DAML arguments
func (t CancelSend) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["onRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OnRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampCid
	}()

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	return m
}

func (t CancelSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CancelSend) UnmarshalJSON(data []byte) error {
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

func (t CreateRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t CreateRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CreateRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Execute is a Record type
type Execute struct {
	RmnRemoteCid          CONTRACT_ID          `json:"rmnRemoteCid"`
	OffRampCid            CONTRACT_ID          `json:"offRampCid"`
	GlobalConfigCid       CONTRACT_ID          `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID          `json:"tokenAdminRegistryCid"`
	ExecutingMessageCid   CONTRACT_ID          `json:"executingMessageCid"`
	ReceiverRequiredCCVs  []RawInstanceAddress `json:"receiverRequiredCCVs"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

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

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
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

func (t Execute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t ExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutionStateChanged) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "ExecutionStateChanged")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutionStateChanged) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "ExecutionStateChanged")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutionStateChanged) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
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

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutionStateChanged) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["event"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *ExecutionStateChanged) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for ExecutionStateChanged

// Archive exercises the Archive choice on this ExecutionStateChanged contract
// This method uses the package name in the template ID
func (t ExecutionStateChanged) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "ExecutionStateChanged"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutionStateChanged) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "ExecutionStateChanged"),
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

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

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

func (t ExecutionStateChangedEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t GetExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForExecute2 is a Record type
type GetRequiredCCVsForExecute2 struct {
	OffRampCid           CONTRACT_ID          `json:"offRampCid"`
	GlobalConfigCid      CONTRACT_ID          `json:"globalConfigCid"`
	ReceiverRequiredCCVs []RawInstanceAddress `json:"receiverRequiredCCVs"`
	SourceChainSelector  NUMERIC              `json:"sourceChainSelector"`
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

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["sourceChainSelector"] = t.SourceChainSelector

	return m
}

func (t GetRequiredCCVsForExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetRequiredCCVsForExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForSend2 is a Record type
type GetRequiredCCVsForSend2 struct {
	OnRampCid         CONTRACT_ID `json:"onRampCid"`
	GlobalConfigCid   CONTRACT_ID `json:"globalConfigCid"`
	DestChainSelector NUMERIC     `json:"destChainSelector"`
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

	m["destChainSelector"] = t.DestChainSelector

	return m
}

func (t GetRequiredCCVsForSend2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

	m["destChainSelector"] = t.DestChainSelector

	return m
}

func (t GetSequenceNumber) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t HasRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *HasRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PerPartyRouter is a Template type
type PerPartyRouter struct {
	InstanceId              TEXT   `json:"instanceId"`
	CcipOwner               PARTY  `json:"ccipOwner"`
	PartyOwner              PARTY  `json:"partyOwner"`
	OutboundSequenceNumbers GENMAP `json:"outboundSequenceNumbers"`
	ExecutionStates         GENMAP `json:"executionStates"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t PerPartyRouter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t PerPartyRouter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() interface{} {
		if t.OutboundSequenceNumbers == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionStates"] = func() interface{} {
		if t.ExecutionStates == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ExecutionStates}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t PerPartyRouter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() interface{} {
		if t.OutboundSequenceNumbers == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionStates"] = func() interface{} {
		if t.ExecutionStates == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.ExecutionStates}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t PerPartyRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PerPartyRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for PerPartyRouter

// Archive exercises the Archive choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCIPSend exercises the CCIPSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) CCIPSend(contractID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// CCIPSendWithPackageID exercises the CCIPSend choice using the provided package ID instead of package name
func (t PerPartyRouter) CCIPSendWithPackageID(contractID string, packageID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareExecute(contractID string, args PrepareExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// Execute exercises the Execute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t PerPartyRouter) ExecuteWithPackageID(contractID string, packageID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionState exercises the GetExecutionState choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetExecutionState(contractID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionStateWithPackageID exercises the GetExecutionState choice using the provided package ID instead of package name
func (t PerPartyRouter) GetExecutionStateWithPackageID(contractID string, packageID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumber exercises the GetSequenceNumber choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetSequenceNumber(contractID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumberWithPackageID exercises the GetSequenceNumber choice using the provided package ID instead of package name
func (t PerPartyRouter) GetSequenceNumberWithPackageID(contractID string, packageID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// PrepareSend exercises the PrepareSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareSend(contractID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendWithPackageID exercises the PrepareSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareSendWithPackageID(contractID string, packageID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// CancelSend exercises the CancelSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) CancelSend(contractID string, args CancelSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CancelSend",
		Arguments:  argsToMap(args),
	}
}

// CancelSendWithPackageID exercises the CancelSend choice using the provided package ID instead of package name
func (t PerPartyRouter) CancelSendWithPackageID(contractID string, packageID string, args CancelSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CancelSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterFactory is a Template type
type PerPartyRouterFactory struct {
	InstanceId        TEXT   `json:"instanceId"`
	CcipOwner         PARTY  `json:"ccipOwner"`
	RegisteredRouters GENMAP `json:"registeredRouters"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t PerPartyRouterFactory) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouterFactory) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t PerPartyRouterFactory) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registeredRouters"] = func() interface{} {
		if t.RegisteredRouters == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RegisteredRouters}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t PerPartyRouterFactory) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registeredRouters"] = func() interface{} {
		if t.RegisteredRouters == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RegisteredRouters}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t PerPartyRouterFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PerPartyRouterFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for PerPartyRouterFactory

// Archive exercises the Archive choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CreateRouter exercises the CreateRouter choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) CreateRouter(contractID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// CreateRouterWithPackageID exercises the CreateRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) CreateRouterWithPackageID(contractID string, packageID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouter exercises the HasRouter choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) HasRouter(contractID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouterWithPackageID exercises the HasRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) HasRouterWithPackageID(contractID string, packageID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute2 is a Record type
type PrepareExecute2 struct {
	OffRampCid         CONTRACT_ID `json:"offRampCid"`
	RmnRemoteCid       CONTRACT_ID `json:"rmnRemoteCid"`
	EncodedMessage     TEXT        `json:"encodedMessage"`
	ReceiverParty      PARTY       `json:"receiverParty"`
	TokenReceiverParty *PARTY      `json:"tokenReceiverParty"`
	Caller             PARTY       `json:"caller"`
}

// ToMap converts PrepareExecute2 to a map for DAML arguments
func (t PrepareExecute2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["offRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OffRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampCid
	}()

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["encodedMessage"] = string(t.EncodedMessage)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	if t.TokenReceiverParty != nil {
		m["tokenReceiverParty"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.TokenReceiverParty).ToMap(),
		}
	} else {
		m["tokenReceiverParty"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PrepareExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PrepareExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// PrepareSend is a Record type
type PrepareSend struct {
	OnRampCid             CONTRACT_ID          `json:"onRampCid"`
	GlobalConfigCid       CONTRACT_ID          `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID          `json:"tokenAdminRegistryCid"`
	FeeQuoterCid          CONTRACT_ID          `json:"feeQuoterCid"`
	RmnRemoteCid          CONTRACT_ID          `json:"rmnRemoteCid"`
	DestChainSelector     NUMERIC              `json:"destChainSelector"`
	Receiver              TEXT                 `json:"receiver"`
	Payload               TEXT                 `json:"payload"`
	CcipReceiveGasLimit   INT64                `json:"ccipReceiveGasLimit"`
	SenderRequiredCCVs    []RawInstanceAddress `json:"senderRequiredCCVs"`
	WithTokenTransfer     BOOL                 `json:"withTokenTransfer"`
	TokenReceiver         *TEXT                `json:"tokenReceiver"`
	FeeToken              InstrumentId         `json:"feeToken"`
}

// ToMap converts PrepareSend to a map for DAML arguments
func (t PrepareSend) ToMap() map[string]interface{} {
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

	m["feeQuoterCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["destChainSelector"] = t.DestChainSelector

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	m["senderRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderRequiredCCVs))
		for _, e := range t.SenderRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["withTokenTransfer"] = bool(t.WithTokenTransfer)

	if t.TokenReceiver != nil {
		m["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.TokenReceiver),
		}
	} else {
		m["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["feeToken"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	return m
}

func (t PrepareSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PrepareSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
