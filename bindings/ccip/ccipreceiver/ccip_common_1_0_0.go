package ccipreceiver

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
	MessageId TEXT `json:"messageId"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	Sender TEXT `json:"sender"`

	Payload TEXT `json:"payload"`

	TokenAmounts []TokenAmount `json:"tokenAmounts"`
}

// ToMap converts Any2CantonMessage to a map for DAML arguments
func (t Any2CantonMessage) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["messageId"] = string(t.MessageId)

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["sender"] = string(t.Sender)

	m["payload"] = string(t.Payload)

	m["tokenAmounts"] = func() []interface{} {
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
	}()

	return m
}

func (t Any2CantonMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Any2CantonMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverView is a Record type
type Any2CantonMessageReceiverView struct {
	CcipOwner PARTY `json:"ccipOwner"`
}

// ToMap converts Any2CantonMessageReceiverView to a map for DAML arguments
func (t Any2CantonMessageReceiverView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccipOwner"] = t.CcipOwner.ToMap()

	return m
}

func (t Any2CantonMessageReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Any2CantonMessageReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverCCIPReceive is a Record type
type Any2CantonMessageReceiverCCIPReceive struct {
	Message Any2CantonMessage `json:"message"`
}

// ToMap converts Any2CantonMessageReceiverCCIPReceive to a map for DAML arguments
func (t Any2CantonMessageReceiverCCIPReceive) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	return m
}

func (t Any2CantonMessageReceiverCCIPReceive) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Any2CantonMessageReceiverCCIPReceive) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Any2CantonMessageReceiverGetCCVs is a Record type
type Any2CantonMessageReceiverGetCCVs struct {
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	Caller PARTY `json:"caller"`
}

// ToMap converts Any2CantonMessageReceiverGetCCVs to a map for DAML arguments
func (t Any2CantonMessageReceiverGetCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Any2CantonMessageReceiverGetCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Any2CantonMessageReceiverGetCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVRegistry is a Template type
type CCVRegistry struct {
	CcipOwner PARTY `json:"ccipOwner"`

	InstanceId TEXT `json:"instanceId"`
}

// GetTemplateID returns the template ID for this template
func (t CCVRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t CCVRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVRegistry contract
func (t CCVRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// CCVRegistryIssueVerifyTicket exercises the CCVRegistry_IssueVerifyTicket choice on this CCVRegistry contract
func (t CCVRegistry) CCVRegistryIssueVerifyTicket(contractID string, args CCVRegistryIssueVerifyTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCVRegistry", "CCVRegistry"),

		ContractID: contractID,
		Choice:     "CCVRegistry_IssueVerifyTicket",

		Arguments: argsToMap(args),
	}
}

// CCVRegistryIssueCCVTicket is a Record type
type CCVRegistryIssueCCVTicket struct {
	CcvOwner PARTY `json:"ccvOwner"`

	VerifierBlob TEXT `json:"verifierBlob"`

	MessageSentObservers []PARTY `json:"messageSentObservers"`

	Sender PARTY `json:"sender"`

	Receipt Receipt `json:"receipt"`
}

// ToMap converts CCVRegistryIssueCCVTicket to a map for DAML arguments
func (t CCVRegistryIssueCCVTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["verifierBlob"] = string(t.VerifierBlob)

	m["messageSentObservers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["sender"] = t.Sender.ToMap()

	m["receipt"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Receipt).(mapper); ok {
			return m.toMap()
		}
		return t.Receipt
	}()

	return m
}

func (t CCVRegistryIssueCCVTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCVRegistryIssueCCVTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVRegistryIssueVerifyTicket is a Record type
type CCVRegistryIssueVerifyTicket struct {
	CcvOwner PARTY `json:"ccvOwner"`

	VerifierResults TEXT `json:"verifierResults"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	SequenceNumber NUMERIC `json:"sequenceNumber"`

	Receiver PARTY `json:"receiver"`
}

// ToMap converts CCVRegistryIssueVerifyTicket to a map for DAML arguments
func (t CCVRegistryIssueVerifyTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvOwner"] = t.CcvOwner.ToMap()

	m["verifierResults"] = string(t.VerifierResults)

	m["messageHash"] = string(t.MessageHash)

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["sequenceNumber"] = (*big.Int)(t.SequenceNumber)

	m["receiver"] = t.Receiver.ToMap()

	return m
}

func (t CCVRegistryIssueVerifyTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCVRegistryIssueVerifyTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVTicket is a Template type
type CCVTicket struct {
	CcvId TEXT `json:"ccvId"`

	CcvOwner PARTY `json:"ccvOwner"`

	CcipOwner PARTY `json:"ccipOwner"`

	Sender PARTY `json:"sender"`

	VerifierBlob TEXT `json:"verifierBlob"`

	MessageSentObservers []PARTY `json:"messageSentObservers"`

	Receipt Receipt `json:"receipt"`
}

// GetTemplateID returns the template ID for this template
func (t CCVTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvId"] = string(t.CcvId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwner"] = t.CcvOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifierBlob"] = string(t.VerifierBlob)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageSentObservers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receipt"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Receipt).(mapper); ok {
			return m.toMap()
		}
		return t.Receipt
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t CCVTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVTicket contract
func (t CCVTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// CCVTicketConsume is a Record type
type CCVTicketConsume struct {
}

// ToMap converts CCVTicketConsume to a map for DAML arguments
func (t CCVTicketConsume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t CCVTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCVTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCVVerifyTicket is a Template type
type CCVVerifyTicket struct {
	CcvId TEXT `json:"ccvId"`

	CcvOwner PARTY `json:"ccvOwner"`

	CcipOwner PARTY `json:"ccipOwner"`

	Caller PARTY `json:"caller"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	SequenceNumber NUMERIC `json:"sequenceNumber"`
}

// GetTemplateID returns the template ID for this template
func (t CCVVerifyTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVVerifyTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t CCVVerifyTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvId"] = string(t.CcvId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvOwner"] = t.CcvOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["caller"] = t.Caller.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
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

func (t CCVVerifyTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCVVerifyTicket contract
func (t CCVVerifyTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "CCVVerifyTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// CCVVerifyTicketConsume is a Record type
type CCVVerifyTicketConsume struct {
}

// ToMap converts CCVVerifyTicketConsume to a map for DAML arguments
func (t CCVVerifyTicketConsume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t CCVVerifyTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCVVerifyTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Canton2AnyMessage is a Record type
type Canton2AnyMessage struct {
	Receiver TEXT `json:"receiver"`

	Payload TEXT `json:"payload"`

	FeeToken InstrumentId `json:"feeToken"`

	ExtraArgs TEXT `json:"extraArgs"`

	TokenAmounts []TokenAmount `json:"tokenAmounts"`
}

// ToMap converts Canton2AnyMessage to a map for DAML arguments
func (t Canton2AnyMessage) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["feeToken"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	m["extraArgs"] = string(t.ExtraArgs)

	m["tokenAmounts"] = func() []interface{} {
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
	}()

	return m
}

func (t Canton2AnyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Canton2AnyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CrossChainVerifierView is a Record type
type CrossChainVerifierView struct {
	CcipOwner PARTY `json:"ccipOwner"`

	StorageLocation TEXT `json:"storageLocation"`
}

// ToMap converts CrossChainVerifierView to a map for DAML arguments
func (t CrossChainVerifierView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

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

// CrossChainVerifierForwardToVerifier is a Record type
type CrossChainVerifierForwardToVerifier struct {
	CcvRegistryCid CONTRACT_ID `json:"ccvRegistryCid"`

	Message MessageV1 `json:"message"`

	MessageId TEXT `json:"messageId"`

	FeeToken InstrumentId `json:"feeToken"`

	FeeTokenAmount NUMERIC `json:"feeTokenAmount"`

	VerifierArgs TEXT `json:"verifierArgs"`

	Caller PARTY `json:"caller"`
}

// ToMap converts CrossChainVerifierForwardToVerifier to a map for DAML arguments
func (t CrossChainVerifierForwardToVerifier) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.CcvRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.CcvRegistryCid
	}()

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["messageId"] = string(t.MessageId)

	m["feeToken"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	m["feeTokenAmount"] = (*big.Int)(t.FeeTokenAmount)

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
	CcvRegistryCid CONTRACT_ID `json:"ccvRegistryCid"`

	Message MessageV1 `json:"message"`

	MessageId TEXT `json:"messageId"`

	VerifierResults TEXT `json:"verifierResults"`

	Receiver PARTY `json:"receiver"`

	Caller PARTY `json:"caller"`
}

// ToMap converts CrossChainVerifierVerifyMessage to a map for DAML arguments
func (t CrossChainVerifierVerifyMessage) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.CcvRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.CcvRegistryCid
	}()

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["messageId"] = string(t.MessageId)

	m["verifierResults"] = string(t.VerifierResults)

	m["receiver"] = t.Receiver.ToMap()

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
	IsEnabled BOOL `json:"isEnabled"`

	DefaultExecutor TEXT `json:"defaultExecutor"`

	OffRampAddress TEXT `json:"offRampAddress"`

	LaneMandatedCCVs []TEXT `json:"laneMandatedCCVs"`

	DefaultCCVs []TEXT `json:"defaultCCVs"`
}

// ToMap converts DestChainConfig to a map for DAML arguments
func (t DestChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["isEnabled"] = bool(t.IsEnabled)

	m["defaultExecutor"] = string(t.DefaultExecutor)

	m["offRampAddress"] = string(t.OffRampAddress)

	m["laneMandatedCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	m["defaultCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			res = append(res, string(e))
		}
		return res
	}()

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

// GetDestChainConfig is a Record type
type GetDestChainConfig struct {
	DestChainSelector NUMERIC `json:"destChainSelector"`

	Caller PARTY `json:"caller"`
}

// ToMap converts GetDestChainConfig to a map for DAML arguments
func (t GetDestChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

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

	Caller PARTY `json:"caller"`
}

// ToMap converts GetSourceChainConfig to a map for DAML arguments
func (t GetSourceChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

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
	CcipOwner PARTY `json:"ccipOwner"`

	InstanceId TEXT `json:"instanceId"`

	ChainSelector NUMERIC `json:"chainSelector"`

	OnRampAddress TEXT `json:"onRampAddress"`

	DestChainConfigs GENMAP `json:"destChainConfigs"`

	SourceChainConfigs GENMAP `json:"sourceChainConfigs"`
}

// GetTemplateID returns the template ID for this template
func (t GlobalConfig) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig")
}

// CreateCommand returns a CreateCommand for this template
func (t GlobalConfig) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	if t.ChainSelector != nil {
		args["chainSelector"] = (*big.Int)(t.ChainSelector)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["onRampAddress"] = string(t.OnRampAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destChainConfigs"] = func() interface{} {
		if t.DestChainConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.DestChainConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceChainConfigs"] = func() interface{} {
		if t.SourceChainConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.SourceChainConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
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
func (t GlobalConfig) GetDestChainConfig(contractID string, args GetDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),

		ContractID: contractID,
		Choice:     "GetDestChainConfig",

		Arguments: argsToMap(args),
	}
}

// GetSourceChainConfig exercises the GetSourceChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) GetSourceChainConfig(contractID string, args GetSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),

		ContractID: contractID,
		Choice:     "GetSourceChainConfig",

		Arguments: argsToMap(args),
	}
}

// UpdateDestChainConfig exercises the UpdateDestChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) UpdateDestChainConfig(contractID string, args UpdateDestChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),

		ContractID: contractID,
		Choice:     "UpdateDestChainConfig",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this GlobalConfig contract
func (t GlobalConfig) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// UpdateSourceChainConfig exercises the UpdateSourceChainConfig choice on this GlobalConfig contract
func (t GlobalConfig) UpdateSourceChainConfig(contractID string, args UpdateSourceChainConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.GlobalConfig", "GlobalConfig"),

		ContractID: contractID,
		Choice:     "UpdateSourceChainConfig",

		Arguments: argsToMap(args),
	}
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
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Internal", "MessageExecutionState")
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
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	DestChainSelector NUMERIC `json:"destChainSelector"`

	SequenceNumber NUMERIC `json:"sequenceNumber"`

	ExecutionGasLimit INT64 `json:"executionGasLimit"`

	CcipReceiveGasLimit INT64 `json:"ccipReceiveGasLimit"`

	Finality INT64 `json:"finality"`

	CcvAndExecutorHash TEXT `json:"ccvAndExecutorHash"`

	OnRampAddress TEXT `json:"onRampAddress"`

	OffRampAddress TEXT `json:"offRampAddress"`

	Sender TEXT `json:"sender"`

	Receiver TEXT `json:"receiver"`

	DestBlob TEXT `json:"destBlob"`

	TokenTransfer *TokenTransferV1 `json:"tokenTransfer"`

	MessageData TEXT `json:"messageData"`
}

// ToMap converts MessageV1 to a map for DAML arguments
func (t MessageV1) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["sequenceNumber"] = (*big.Int)(t.SequenceNumber)

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
		m["tokenTransfer"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenTransfer,
		}
	} else {
		m["tokenTransfer"] = map[string]interface{}{
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

// Receipt is a Record type
type Receipt struct {
	Issuer TEXT `json:"issuer"`

	DestGasLimit INT64 `json:"destGasLimit"`

	DestBytesOverhead INT64 `json:"destBytesOverhead"`

	FeeTokenAmount NUMERIC `json:"feeTokenAmount"`

	ExtraArgs TEXT `json:"extraArgs"`
}

// ToMap converts Receipt to a map for DAML arguments
func (t Receipt) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issuer"] = string(t.Issuer)

	m["destGasLimit"] = int64(t.DestGasLimit)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["feeTokenAmount"] = (*big.Int)(t.FeeTokenAmount)

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

// SourceChainConfig is a Record type
type SourceChainConfig struct {
	IsEnabled BOOL `json:"isEnabled"`

	OnRampAddress TEXT `json:"onRampAddress"`

	LaneMandatedCCVs []TEXT `json:"laneMandatedCCVs"`

	DefaultCCVs []TEXT `json:"defaultCCVs"`
}

// ToMap converts SourceChainConfig to a map for DAML arguments
func (t SourceChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["isEnabled"] = bool(t.IsEnabled)

	m["onRampAddress"] = string(t.OnRampAddress)

	m["laneMandatedCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.LaneMandatedCCVs))
		for _, e := range t.LaneMandatedCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	m["defaultCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.DefaultCCVs))
		for _, e := range t.DefaultCCVs {
			res = append(res, string(e))
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

// TokenAmount is a Record type
type TokenAmount struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Amount NUMERIC `json:"amount"`
}

// ToMap converts TokenAmount to a map for DAML arguments
func (t TokenAmount) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = (*big.Int)(t.Amount)

	return m
}

func (t TokenAmount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAmount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolCCVVerifiedTicket is a Template type
type TokenPoolCCVVerifiedTicket struct {
	PoolOwner PARTY `json:"poolOwner"`

	CcipOwner PARTY `json:"ccipOwner"`

	Receiver PARTY `json:"receiver"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	VerifiedCCVIds []TEXT `json:"verifiedCCVIds"`
}

// GetTemplateID returns the template ID for this template
func (t TokenPoolCCVVerifiedTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenPoolCCVVerifiedTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenPoolCCVVerifiedTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != nil {
		args["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["verifiedCCVIds"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.VerifiedCCVIds))
		for _, e := range t.VerifiedCCVIds {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t TokenPoolCCVVerifiedTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenPoolCCVVerifiedTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TokenPoolCCVVerifiedTicket

// TokenPoolCCVVerifiedTicketConsume exercises the TokenPoolCCVVerifiedTicket_Consume choice on this TokenPoolCCVVerifiedTicket contract
func (t TokenPoolCCVVerifiedTicket) TokenPoolCCVVerifiedTicketConsume(contractID string, args TokenPoolCCVVerifiedTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenPoolCCVVerifiedTicket"),

		ContractID: contractID,
		Choice:     "TokenPoolCCVVerifiedTicket_Consume",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenPoolCCVVerifiedTicket contract
func (t TokenPoolCCVVerifiedTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenPoolCCVVerifiedTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TokenPoolCCVVerifiedTicketConsume is a Record type
type TokenPoolCCVVerifiedTicketConsume struct {
}

// ToMap converts TokenPoolCCVVerifiedTicketConsume to a map for DAML arguments
func (t TokenPoolCCVVerifiedTicketConsume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t TokenPoolCCVVerifiedTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenPoolCCVVerifiedTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenReceiveTicket is a Template type
type TokenReceiveTicket struct {
	CcipOwner PARTY `json:"ccipOwner"`

	PoolOwner PARTY `json:"poolOwner"`

	Receiver PARTY `json:"receiver"`

	TokenReceiver PARTY `json:"tokenReceiver"`

	InstrumentId InstrumentId `json:"instrumentId"`

	Amount NUMERIC `json:"amount"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
}

// GetTemplateID returns the template ID for this template
func (t TokenReceiveTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenReceiveTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = t.TokenReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageHash"] = string(t.MessageHash)

	if t.SourceChainSelector != nil {
		args["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
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
func (t TokenReceiveTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TokenReceiveTicketConsume exercises the TokenReceiveTicket_Consume choice on this TokenReceiveTicket contract
func (t TokenReceiveTicket) TokenReceiveTicketConsume(contractID string, args TokenReceiveTicketConsume) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenReceiveTicket"),

		ContractID: contractID,
		Choice:     "TokenReceiveTicket_Consume",

		Arguments: argsToMap(args),
	}
}

// TokenReceiveTicketConsume is a Record type
type TokenReceiveTicketConsume struct {
}

// ToMap converts TokenReceiveTicketConsume to a map for DAML arguments
func (t TokenReceiveTicketConsume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t TokenReceiveTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenReceiveTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenSendTicket is a Template type
type TokenSendTicket struct {
	PoolOwner PARTY `json:"poolOwner"`

	CcipOwner PARTY `json:"ccipOwner"`

	Sender PARTY `json:"sender"`

	InstrumentId InstrumentId `json:"instrumentId"`

	SourceTokenAddress TEXT `json:"sourceTokenAddress"`

	Amount NUMERIC `json:"amount"`

	DestTokenAddress TEXT `json:"destTokenAddress"`

	TokenReceiver TEXT `json:"tokenReceiver"`

	ExtraData TEXT `json:"extraData"`

	Receipt Receipt `json:"receipt"`
}

// GetTemplateID returns the template ID for this template
func (t TokenSendTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenSendTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenSendTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceTokenAddress"] = string(t.SourceTokenAddress)

	if t.Amount != nil {
		args["amount"] = (*big.Int)(t.Amount)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["destTokenAddress"] = string(t.DestTokenAddress)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenReceiver"] = string(t.TokenReceiver)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["extraData"] = string(t.ExtraData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receipt"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Receipt).(mapper); ok {
			return m.toMap()
		}
		return t.Receipt
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t TokenSendTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenSendTicket contract
func (t TokenSendTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Tickets", "TokenSendTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TokenSendTicketConsume is a Record type
type TokenSendTicketConsume struct {
}

// ToMap converts TokenSendTicketConsume to a map for DAML arguments
func (t TokenSendTicketConsume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t TokenSendTicketConsume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenSendTicketConsume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenTransferV1 is a Record type
type TokenTransferV1 struct {
	Amount NUMERIC `json:"amount"`

	SourcePoolAddress TEXT `json:"sourcePoolAddress"`

	SourceTokenAddress TEXT `json:"sourceTokenAddress"`

	DestTokenAddress TEXT `json:"destTokenAddress"`

	TokenReceiver TEXT `json:"tokenReceiver"`

	ExtraData TEXT `json:"extraData"`
}

// ToMap converts TokenTransferV1 to a map for DAML arguments
func (t TokenTransferV1) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["amount"] = (*big.Int)(t.Amount)

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
	DestChainSelector NUMERIC `json:"destChainSelector"`

	Config DestChainConfig `json:"config"`
}

// ToMap converts UpdateDestChainConfig to a map for DAML arguments
func (t UpdateDestChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["config"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	Config SourceChainConfig `json:"config"`
}

// ToMap converts UpdateSourceChainConfig to a map for DAML arguments
func (t UpdateSourceChainConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["config"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
