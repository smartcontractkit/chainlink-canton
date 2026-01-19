package lockreleasetokenpool

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

// IIAny2CantonMessageReceiver is a DAML interface
type IIAny2CantonMessageReceiver interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// Any2CantonMessageReceiverGetCCVs executes the Any2CantonMessageReceiver_GetCCVs choice
	Any2CantonMessageReceiverGetCCVs(contractID string, args Any2CantonMessageReceiverGetCCVs) *model.ExerciseCommand

	// Any2CantonMessageReceiverCCIPReceive executes the Any2CantonMessageReceiver_CCIPReceive choice
	Any2CantonMessageReceiverCCIPReceive(contractID string, args Any2CantonMessageReceiverCCIPReceive) *model.ExerciseCommand
}

// IICrossChainVerifier is a DAML interface
type IICrossChainVerifier interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// CrossChainVerifierVerifyMessage executes the CrossChainVerifier_VerifyMessage choice
	CrossChainVerifierVerifyMessage(contractID string, args CrossChainVerifierVerifyMessage) *model.ExerciseCommand

	// CrossChainVerifierForwardToVerifier executes the CrossChainVerifier_ForwardToVerifier choice
	CrossChainVerifierForwardToVerifier(contractID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand
}

// Any2CantonMessage is a Record type
type Any2CantonMessage struct {
	MessageId           TEXT          `json:"messageId"`
	SourceChainSelector NUMERIC       `json:"sourceChainSelector"`
	Sender              TEXT          `json:"sender"`
	Payload             TEXT          `json:"payload"`
	TokenAmounts        []TokenAmount `json:"tokenAmounts"`
}

// toMap converts Any2CantonMessage to a map for DAML arguments
func (t Any2CantonMessage) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageId":           string(t.MessageId),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sender":              string(t.Sender),
		"payload":             string(t.Payload),
		"tokenAmounts": func() []interface{} {
			res := make([]interface{}, 0, len(t.TokenAmounts))
			for _, e := range t.TokenAmounts {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for Any2CantonMessage using JsonCodec
func (t Any2CantonMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Any2CantonMessage using JsonCodec
func (t *Any2CantonMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverView is a Record type
type Any2CantonMessageReceiverView struct {
	CcipOwner PARTY `json:"ccipOwner"`
}

// toMap converts Any2CantonMessageReceiverView to a map for DAML arguments
func (t Any2CantonMessageReceiverView) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccipOwner": t.CcipOwner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for Any2CantonMessageReceiverView using JsonCodec
func (t Any2CantonMessageReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Any2CantonMessageReceiverView using JsonCodec
func (t *Any2CantonMessageReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverCCIPReceive is a Record type
type Any2CantonMessageReceiverCCIPReceive struct {
	Message Any2CantonMessage `json:"message"`
}

// toMap converts Any2CantonMessageReceiverCCIPReceive to a map for DAML arguments
func (t Any2CantonMessageReceiverCCIPReceive) toMap() map[string]interface{} {
	return map[string]interface{}{

		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for Any2CantonMessageReceiverCCIPReceive using JsonCodec
func (t Any2CantonMessageReceiverCCIPReceive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Any2CantonMessageReceiverCCIPReceive using JsonCodec
func (t *Any2CantonMessageReceiverCCIPReceive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverGetCCVs is a Record type
type Any2CantonMessageReceiverGetCCVs struct {
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
	Caller              PARTY   `json:"caller"`
}

// toMap converts Any2CantonMessageReceiverGetCCVs to a map for DAML arguments
func (t Any2CantonMessageReceiverGetCCVs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"caller":              t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for Any2CantonMessageReceiverGetCCVs using JsonCodec
func (t Any2CantonMessageReceiverGetCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Any2CantonMessageReceiverGetCCVs using JsonCodec
func (t *Any2CantonMessageReceiverGetCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVRegistry is a Template type
type CCVRegistry struct {
	CcipOwner  PARTY `json:"ccipOwner"`
	InstanceId TEXT  `json:"instanceId"`
}

// GetTemplateID returns the template ID for this template
func (t CCVRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CCVRegistry using JsonCodec
func (t CCVRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVRegistry using JsonCodec
func (t *CCVRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCVRegistry

// CCVRegistryIssueCCVTicket exercises the CCVRegistry_IssueCCVTicket choice on this CCVRegistry contract
func (t CCVRegistry) CCVRegistryIssueCCVTicket(contractID string, args CCVRegistryIssueCCVTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry"),
		ContractID: contractID,
		Choice:     "CCVRegistry_IssueCCVTicket",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVRegistry contract
func (t CCVRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCVRegistryIssueVerifyTicket exercises the CCVRegistry_IssueVerifyTicket choice on this CCVRegistry contract
func (t CCVRegistry) CCVRegistryIssueVerifyTicket(contractID string, args CCVRegistryIssueVerifyTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry"),
		ContractID: contractID,
		Choice:     "CCVRegistry_IssueVerifyTicket",
		Arguments:  argsToMap(args),
	}
}

// CCVRegistryIssueCCVTicket is a Record type
type CCVRegistryIssueCCVTicket struct {
	CcvOwner             PARTY   `json:"ccvOwner"`
	VerifierBlob         TEXT    `json:"verifierBlob"`
	MessageSentObservers []PARTY `json:"messageSentObservers"`
	Sender               PARTY   `json:"sender"`
}

// toMap converts CCVRegistryIssueCCVTicket to a map for DAML arguments
func (t CCVRegistryIssueCCVTicket) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvOwner":     t.CcvOwner.ToMap(),
		"verifierBlob": string(t.VerifierBlob),
		"messageSentObservers": func() []interface{} {
			res := make([]interface{}, 0, len(t.MessageSentObservers))
			for _, e := range t.MessageSentObservers {
				res = append(res, e.ToMap())
			}
			return res
		}(),
		"sender": t.Sender.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CCVRegistryIssueCCVTicket using JsonCodec
func (t CCVRegistryIssueCCVTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVRegistryIssueCCVTicket using JsonCodec
func (t *CCVRegistryIssueCCVTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVRegistryIssueVerifyTicket is a Record type
type CCVRegistryIssueVerifyTicket struct {
	CcvOwner            PARTY   `json:"ccvOwner"`
	CcvData             TEXT    `json:"ccvData"`
	MessageHash         TEXT    `json:"messageHash"`
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC `json:"sequenceNumber"`
	Receiver            PARTY   `json:"receiver"`
}

// toMap converts CCVRegistryIssueVerifyTicket to a map for DAML arguments
func (t CCVRegistryIssueVerifyTicket) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvOwner":            t.CcvOwner.ToMap(),
		"ccvData":             string(t.CcvData),
		"messageHash":         string(t.MessageHash),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"receiver":            t.Receiver.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CCVRegistryIssueVerifyTicket using JsonCodec
func (t CCVRegistryIssueVerifyTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVRegistryIssueVerifyTicket using JsonCodec
func (t *CCVRegistryIssueVerifyTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVTicket is a Template type
type CCVTicket struct {
	CcvId                TEXT    `json:"ccvId"`
	CcvOwner             PARTY   `json:"ccvOwner"`
	CcipOwner            PARTY   `json:"ccipOwner"`
	Sender               PARTY   `json:"sender"`
	VerifierBlob         TEXT    `json:"verifierBlob"`
	MessageSentObservers []PARTY `json:"messageSentObservers"`
}

// GetTemplateID returns the template ID for this template
func (t CCVTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccvId"] = string(t.CcvId)

	args["ccvOwner"] = t.CcvOwner.ToMap()

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["sender"] = t.Sender.ToMap()

	args["verifierBlob"] = string(t.VerifierBlob)

	if len(t.MessageSentObservers) > 0 {
		args["messageSentObservers"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.MessageSentObservers))
			for _, e := range t.MessageSentObservers {
				res = append(res, e.ToMap())
			}
			return res
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CCVTicket using JsonCodec
func (t CCVTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVTicket using JsonCodec
func (t *CCVTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCVTicket

// CCVTicketConsume exercises the CCVTicket_Consume choice on this CCVTicket contract
func (t CCVTicket) CCVTicketConsume(contractID string, args CCVTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVTicket"),
		ContractID: contractID,
		Choice:     "CCVTicket_Consume",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVTicket contract
func (t CCVTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCVTicketConsume is a Record type
type CCVTicketConsume struct {
}

// toMap converts CCVTicketConsume to a map for DAML arguments
func (t CCVTicketConsume) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for CCVTicketConsume using JsonCodec
func (t CCVTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVTicketConsume using JsonCodec
func (t *CCVTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVVerifyTicket is a Template type
type CCVVerifyTicket struct {
	CcvId               TEXT    `json:"ccvId"`
	CcvOwner            PARTY   `json:"ccvOwner"`
	CcipOwner           PARTY   `json:"ccipOwner"`
	Caller              PARTY   `json:"caller"`
	MessageHash         TEXT    `json:"messageHash"`
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC `json:"sequenceNumber"`
}

// GetTemplateID returns the template ID for this template
func (t CCVVerifyTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVVerifyTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVVerifyTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccvId"] = string(t.CcvId)

	args["ccvOwner"] = t.CcvOwner.ToMap()

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["caller"] = t.Caller.ToMap()

	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != nil {
		args["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)
	}

	if t.SequenceNumber != nil {
		args["sequenceNumber"] = (*big.Int)(t.SequenceNumber)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CCVVerifyTicket using JsonCodec
func (t CCVVerifyTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVVerifyTicket using JsonCodec
func (t *CCVVerifyTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCVVerifyTicket

// CCVVerifyTicketConsume exercises the CCVVerifyTicket_Consume choice on this CCVVerifyTicket contract
func (t CCVVerifyTicket) CCVVerifyTicketConsume(contractID string, args CCVVerifyTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVVerifyTicket"),
		ContractID: contractID,
		Choice:     "CCVVerifyTicket_Consume",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVVerifyTicket contract
func (t CCVVerifyTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVVerifyTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCVVerifyTicketConsume is a Record type
type CCVVerifyTicketConsume struct {
}

// toMap converts CCVVerifyTicketConsume to a map for DAML arguments
func (t CCVVerifyTicketConsume) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for CCVVerifyTicketConsume using JsonCodec
func (t CCVVerifyTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCVVerifyTicketConsume using JsonCodec
func (t *CCVVerifyTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Canton2AnyMessage is a Record type
type Canton2AnyMessage struct {
	Receiver     TEXT          `json:"receiver"`
	Payload      TEXT          `json:"payload"`
	FeeToken     InstrumentId  `json:"feeToken"`
	ExtraArgs    TEXT          `json:"extraArgs"`
	TokenAmounts []TokenAmount `json:"tokenAmounts"`
}

// toMap converts Canton2AnyMessage to a map for DAML arguments
func (t Canton2AnyMessage) toMap() map[string]interface{} {
	return map[string]interface{}{

		"receiver": string(t.Receiver),
		"payload":  string(t.Payload),
		"feeToken": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeToken).(mapper); ok {
				return m.toMap()
			}
			return t.FeeToken
		}(),
		"extraArgs": string(t.ExtraArgs),
		"tokenAmounts": func() []interface{} {
			res := make([]interface{}, 0, len(t.TokenAmounts))
			for _, e := range t.TokenAmounts {
				type mapper interface{ toMap() map[string]interface{} }
				if m, ok := any(e).(mapper); ok {
					res = append(res, m.toMap())
				} else {
					res = append(res, e)
				}
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for Canton2AnyMessage using JsonCodec
func (t Canton2AnyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Canton2AnyMessage using JsonCodec
func (t *Canton2AnyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	CcipOwner       PARTY `json:"ccipOwner"`
	StorageLocation TEXT  `json:"storageLocation"`
}

// toMap converts CrossChainVerifierView to a map for DAML arguments
func (t CrossChainVerifierView) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccipOwner":       t.CcipOwner.ToMap(),
		"storageLocation": string(t.StorageLocation),
	}
}

// MarshalJSON implements custom JSON marshaling for CrossChainVerifierView using JsonCodec
func (t CrossChainVerifierView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CrossChainVerifierView using JsonCodec
func (t *CrossChainVerifierView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierForwardToVerifier is a Record type
type CrossChainVerifierForwardToVerifier struct {
	CcvRegistryCid CONTRACT_ID  `json:"ccvRegistryCid"`
	Message        MessageV1    `json:"message"`
	MessageId      TEXT         `json:"messageId"`
	FeeToken       InstrumentId `json:"feeToken"`
	FeeTokenAmount NUMERIC      `json:"feeTokenAmount"`
	VerifierArgs   TEXT         `json:"verifierArgs"`
	Caller         PARTY        `json:"caller"`
}

// toMap converts CrossChainVerifierForwardToVerifier to a map for DAML arguments
func (t CrossChainVerifierForwardToVerifier) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.CcvRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.CcvRegistryCid
		}(),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"messageId": string(t.MessageId),
		"feeToken": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeToken).(mapper); ok {
				return m.toMap()
			}
			return t.FeeToken
		}(),
		"feeTokenAmount": (*big.Int)(t.FeeTokenAmount),
		"verifierArgs":   string(t.VerifierArgs),
		"caller":         t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CrossChainVerifierForwardToVerifier using JsonCodec
func (t CrossChainVerifierForwardToVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CrossChainVerifierForwardToVerifier using JsonCodec
func (t *CrossChainVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierVerifyMessage is a Record type
type CrossChainVerifierVerifyMessage struct {
	CcvRegistryCid CONTRACT_ID `json:"ccvRegistryCid"`
	Message        MessageV1   `json:"message"`
	MessageId      TEXT        `json:"messageId"`
	CcvData        TEXT        `json:"ccvData"`
	Receiver       PARTY       `json:"receiver"`
	Caller         PARTY       `json:"caller"`
}

// toMap converts CrossChainVerifierVerifyMessage to a map for DAML arguments
func (t CrossChainVerifierVerifyMessage) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.CcvRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.CcvRegistryCid
		}(),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"messageId": string(t.MessageId),
		"ccvData":   string(t.CcvData),
		"receiver":  t.Receiver.ToMap(),
		"caller":    t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CrossChainVerifierVerifyMessage using JsonCodec
func (t CrossChainVerifierVerifyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CrossChainVerifierVerifyMessage using JsonCodec
func (t *CrossChainVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// DestChainConfig is a Record type
type DestChainConfig struct {
	IsEnabled        BOOL   `json:"isEnabled"`
	DefaultExecutor  TEXT   `json:"defaultExecutor"`
	OffRampAddress   TEXT   `json:"offRampAddress"`
	LaneMandatedCCVs []TEXT `json:"laneMandatedCCVs"`
	DefaultCCVs      []TEXT `json:"defaultCCVs"`
}

// toMap converts DestChainConfig to a map for DAML arguments
func (t DestChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"isEnabled":       bool(t.IsEnabled),
		"defaultExecutor": string(t.DefaultExecutor),
		"offRampAddress":  string(t.OffRampAddress),
		"laneMandatedCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.LaneMandatedCCVs))
			for _, e := range t.LaneMandatedCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
		"defaultCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.DefaultCCVs))
			for _, e := range t.DefaultCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for DestChainConfig using JsonCodec
func (t DestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for DestChainConfig using JsonCodec
func (t *DestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetDestChainConfig is a Record type
type GetDestChainConfig struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`
	Caller            PARTY   `json:"caller"`
}

// toMap converts GetDestChainConfig to a map for DAML arguments
func (t GetDestChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"caller":            t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetDestChainConfig using JsonCodec
func (t GetDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetDestChainConfig using JsonCodec
func (t *GetDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetSourceChainConfig is a Record type
type GetSourceChainConfig struct {
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
	Caller              PARTY   `json:"caller"`
}

// toMap converts GetSourceChainConfig to a map for DAML arguments
func (t GetSourceChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"caller":              t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for GetSourceChainConfig using JsonCodec
func (t GetSourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetSourceChainConfig using JsonCodec
func (t *GetSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GlobalConfig is a Template type
type GlobalConfig struct {
	CcipOwner          PARTY   `json:"ccipOwner"`
	InstanceId         TEXT    `json:"instanceId"`
	ChainSelector      NUMERIC `json:"chainSelector"`
	OnRampAddress      TEXT    `json:"onRampAddress"`
	DestChainConfigs   GENMAP  `json:"destChainConfigs"`
	SourceChainConfigs GENMAP  `json:"sourceChainConfigs"`
}

// GetTemplateID returns the template ID for this template
func (t GlobalConfig) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig")
}

// CreateCommand returns a CreateCommand for this template
func (t GlobalConfig) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	if t.ChainSelector != nil {
		args["chainSelector"] = (*big.Int)(t.ChainSelector)
	}

	args["onRampAddress"] = string(t.OnRampAddress)

	if t.DestChainConfigs != nil && len(t.DestChainConfigs) > 0 {
		args["destChainConfigs"] = map[string]interface{}{"_type": "genmap", "value": t.DestChainConfigs}
	}

	if t.SourceChainConfigs != nil && len(t.SourceChainConfigs) > 0 {
		args["sourceChainConfigs"] = map[string]interface{}{"_type": "genmap", "value": t.SourceChainConfigs}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for GlobalConfig using JsonCodec
func (t GlobalConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GlobalConfig using JsonCodec
func (t *GlobalConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for GlobalConfig

// GetDestChainConfig exercises the GetDestChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) GetDestChainConfig(contractID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// GetSourceChainConfig exercises the GetSourceChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) GetSourceChainConfig(contractID string, args GetSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "GetSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// UpdateDestChainConfig exercises the UpdateDestChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) UpdateDestChainConfig(contractID string, args UpdateDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateDestChainConfig",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this GlobalConfig contract
func (t GlobalConfig) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// UpdateSourceChainConfig exercises the UpdateSourceChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) UpdateSourceChainConfig(contractID string, args UpdateSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		ContractID: contractID,
		Choice:     "UpdateSourceChainConfig",
		Arguments:  argsToMap(args),
	}
}

// MessageExecutionState is an enum type
type MessageExecutionState string

const (
	MessageExecutionStateUNTOUCHED   MessageExecutionState = "UNTOUCHED"
	MessageExecutionStateIN_PROGRESS MessageExecutionState = "IN_PROGRESS"
	MessageExecutionStateSUCCESS     MessageExecutionState = "SUCCESS"
	MessageExecutionStateFAILURE     MessageExecutionState = "FAILURE"
)

// GetEnumConstructor implements types.ENUM interface
func (e MessageExecutionState) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e MessageExecutionState) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Internal", "MessageExecutionState")
}

// MarshalJSON implements custom JSON marshaling for MessageExecutionState using JsonCodec
func (e MessageExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for MessageExecutionState using JsonCodec
func (e *MessageExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
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

// toMap converts MessageV1 to a map for DAML arguments
func (t MessageV1) toMap() map[string]interface{} {
	return map[string]interface{}{

		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"destChainSelector":   (*big.Int)(t.DestChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"executionGasLimit":   int64(t.ExecutionGasLimit),
		"ccipReceiveGasLimit": int64(t.CcipReceiveGasLimit),
		"finality":            int64(t.Finality),
		"ccvAndExecutorHash":  string(t.CcvAndExecutorHash),
		"onRampAddress":       string(t.OnRampAddress),
		"offRampAddress":      string(t.OffRampAddress),
		"sender":              string(t.Sender),
		"receiver":            string(t.Receiver),
		"destBlob":            string(t.DestBlob),
		"tokenTransfer":       *t.TokenTransfer,
		"messageData":         string(t.MessageData),
	}
}

// MarshalJSON implements custom JSON marshaling for MessageV1 using JsonCodec
func (t MessageV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MessageV1 using JsonCodec
func (t *MessageV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled        BOOL   `json:"isEnabled"`
	OnRampAddress    TEXT   `json:"onRampAddress"`
	LaneMandatedCCVs []TEXT `json:"laneMandatedCCVs"`
	DefaultCCVs      []TEXT `json:"defaultCCVs"`
}

// toMap converts SourceChainConfig to a map for DAML arguments
func (t SourceChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"isEnabled":     bool(t.IsEnabled),
		"onRampAddress": string(t.OnRampAddress),
		"laneMandatedCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.LaneMandatedCCVs))
			for _, e := range t.LaneMandatedCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
		"defaultCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.DefaultCCVs))
			for _, e := range t.DefaultCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for SourceChainConfig using JsonCodec
func (t SourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SourceChainConfig using JsonCodec
func (t *SourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAmount is a Record type
type TokenAmount struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Amount       NUMERIC      `json:"amount"`
}

// toMap converts TokenAmount to a map for DAML arguments
func (t TokenAmount) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"amount": (*big.Int)(t.Amount),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAmount using JsonCodec
func (t TokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAmount using JsonCodec
func (t *TokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner           PARTY    `json:"ccipOwner"`
	Receiver            PARTY    `json:"receiver"`
	MessageHash         TEXT     `json:"messageHash"`
	SourceChainSelector NUMERIC  `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC  `json:"sequenceNumber"`
	HasTokenTransfer    BOOL     `json:"hasTokenTransfer"`
	TokenAmount         *NUMERIC `json:"tokenAmount"`
	DestTokenAddress    *TEXT    `json:"destTokenAddress"`
	TokenReceiver       *PARTY   `json:"tokenReceiver"`
}

// GetTemplateID returns the template ID for this template
func (t TokenReceiveTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenReceiveTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["receiver"] = t.Receiver.ToMap()

	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != nil {
		args["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)
	}

	if t.SequenceNumber != nil {
		args["sequenceNumber"] = (*big.Int)(t.SequenceNumber)
	}

	args["hasTokenTransfer"] = bool(t.HasTokenTransfer)

	if t.TokenAmount != nil {
		args["tokenAmount"] = map[string]interface{}{
			"_type": "optional",
			"value": (*big.Int)(*t.TokenAmount),
		}
	} else {
		args["tokenAmount"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.DestTokenAddress != nil {
		args["destTokenAddress"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.DestTokenAddress),
		}
	} else {
		args["destTokenAddress"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.TokenReceiver != nil {
		args["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.TokenReceiver).ToMap(),
		}
	} else {
		args["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TokenReceiveTicket using JsonCodec
func (t TokenReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenReceiveTicket using JsonCodec
func (t *TokenReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TokenReceiveTicket

// Archive exercises the Archive choice on this TokenReceiveTicket contract
func (t TokenReceiveTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TokenReceiveTicketConsume exercises the TokenReceiveTicket_Consume choice on this TokenReceiveTicket contract
func (t TokenReceiveTicket) TokenReceiveTicketConsume(contractID string, args TokenReceiveTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket"),
		ContractID: contractID,
		Choice:     "TokenReceiveTicket_Consume",
		Arguments:  argsToMap(args),
	}
}

// TokenReceiveTicketConsume is a Record type
type TokenReceiveTicketConsume struct {
}

// toMap converts TokenReceiveTicketConsume to a map for DAML arguments
func (t TokenReceiveTicketConsume) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for TokenReceiveTicketConsume using JsonCodec
func (t TokenReceiveTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenReceiveTicketConsume using JsonCodec
func (t *TokenReceiveTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenSendTicket is a Template type
type TokenSendTicket struct {
	PoolOwner          PARTY        `json:"poolOwner"`
	CcipOwner          PARTY        `json:"ccipOwner"`
	Sender             PARTY        `json:"sender"`
	InstrumentId       InstrumentId `json:"instrumentId"`
	SourceTokenAddress TEXT         `json:"sourceTokenAddress"`
	Amount             NUMERIC      `json:"amount"`
	DestTokenAddress   TEXT         `json:"destTokenAddress"`
	TokenReceiver      TEXT         `json:"tokenReceiver"`
	ExtraData          TEXT         `json:"extraData"`
}

// GetTemplateID returns the template ID for this template
func (t TokenSendTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenSendTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenSendTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["poolOwner"] = t.PoolOwner.ToMap()

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["sender"] = t.Sender.ToMap()

	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	args["sourceTokenAddress"] = string(t.SourceTokenAddress)

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	args["destTokenAddress"] = string(t.DestTokenAddress)

	args["tokenReceiver"] = string(t.TokenReceiver)

	args["extraData"] = string(t.ExtraData)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TokenSendTicket using JsonCodec
func (t TokenSendTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenSendTicket using JsonCodec
func (t *TokenSendTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TokenSendTicket

// TokenSendTicketConsume exercises the TokenSendTicket_Consume choice on this TokenSendTicket contract
func (t TokenSendTicket) TokenSendTicketConsume(contractID string, args TokenSendTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenSendTicket"),
		ContractID: contractID,
		Choice:     "TokenSendTicket_Consume",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenSendTicket contract
func (t TokenSendTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenSendTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TokenSendTicketConsume is a Record type
type TokenSendTicketConsume struct {
}

// toMap converts TokenSendTicketConsume to a map for DAML arguments
func (t TokenSendTicketConsume) toMap() map[string]interface{} {
	return map[string]interface{}{}
}

// MarshalJSON implements custom JSON marshaling for TokenSendTicketConsume using JsonCodec
func (t TokenSendTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenSendTicketConsume using JsonCodec
func (t *TokenSendTicketConsume) UnmarshalJSON(data []byte) error {
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

// toMap converts TokenTransferV1 to a map for DAML arguments
func (t TokenTransferV1) toMap() map[string]interface{} {
	return map[string]interface{}{

		"amount":             (*big.Int)(t.Amount),
		"sourcePoolAddress":  string(t.SourcePoolAddress),
		"sourceTokenAddress": string(t.SourceTokenAddress),
		"destTokenAddress":   string(t.DestTokenAddress),
		"tokenReceiver":      string(t.TokenReceiver),
		"extraData":          string(t.ExtraData),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenTransferV1 using JsonCodec
func (t TokenTransferV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenTransferV1 using JsonCodec
func (t *TokenTransferV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateDestChainConfig is a Record type
type UpdateDestChainConfig struct {
	DestChainSelector NUMERIC         `json:"destChainSelector"`
	Config            DestChainConfig `json:"config"`
}

// toMap converts UpdateDestChainConfig to a map for DAML arguments
func (t UpdateDestChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"config": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Config).(mapper); ok {
				return m.toMap()
			}
			return t.Config
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for UpdateDestChainConfig using JsonCodec
func (t UpdateDestChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UpdateDestChainConfig using JsonCodec
func (t *UpdateDestChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateSourceChainConfig is a Record type
type UpdateSourceChainConfig struct {
	SourceChainSelector NUMERIC           `json:"sourceChainSelector"`
	Config              SourceChainConfig `json:"config"`
}

// toMap converts UpdateSourceChainConfig to a map for DAML arguments
func (t UpdateSourceChainConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"config": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Config).(mapper); ok {
				return m.toMap()
			}
			return t.Config
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for UpdateSourceChainConfig using JsonCodec
func (t UpdateSourceChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for UpdateSourceChainConfig using JsonCodec
func (t *UpdateSourceChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IIAny2CantonMessageReceiverInterfaceID returns the interface ID for the IIAny2CantonMessageReceiver interface
func IIAny2CantonMessageReceiverInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "CCIP.Interfaces.Any2CantonMessageReceiver", "IAny2CantonMessageReceiver")
}

// IICrossChainVerifierInterfaceID returns the interface ID for the IICrossChainVerifier interface
func IICrossChainVerifierInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "CCIP.Interfaces.CrossChainVerifier", "ICrossChainVerifier")
}
