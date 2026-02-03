package interfaces

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

// ConsumeReceiveTicketResult is a Record type
type ConsumeReceiveTicketResult struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Amount NUMERIC `json:"amount"`

	Receiver PARTY `json:"receiver"`

	TokenReceiver PARTY `json:"tokenReceiver"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`
}

// ToMap converts ConsumeReceiveTicketResult to a map for DAML arguments
func (t ConsumeReceiveTicketResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["amount"] = t.Amount

	m["receiver"] = t.Receiver.ToMap()

	m["tokenReceiver"] = t.TokenReceiver.ToMap()

	m["messageHash"] = string(t.MessageHash)

	m["sourceChainSelector"] = t.SourceChainSelector

	return m
}

func (t ConsumeReceiveTicketResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ConsumeReceiveTicketResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	Owner PARTY `json:"owner"`

	InstanceId TEXT `json:"instanceId"`

	TokenConfigs GENMAP `json:"tokenConfigs"`
}

// GetTemplateID returns the template ID for this template
func (t TokenAdminRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenAdminRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() interface{} {
		if t.TokenConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t TokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TokenAdminRegistry

// TokenAdminRegistryGetTokenConfig exercises the TokenAdminRegistry_GetTokenConfig choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryGetTokenConfig(contractID string, args TokenAdminRegistryGetTokenConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_GetTokenConfig",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistrySetPool exercises the TokenAdminRegistry_SetPool choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistrySetPool(contractID string, args TokenAdminRegistrySetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetPool",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole exercises the TokenAdminRegistry_AcceptAdminRole choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryAcceptAdminRole(contractID string, args TokenAdminRegistryAcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AcceptAdminRole",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryTransferAdminRole exercises the TokenAdminRegistry_TransferAdminRole choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryTransferAdminRole(contractID string, args TokenAdminRegistryTransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_TransferAdminRole",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryIsAdministrator exercises the TokenAdminRegistry_IsAdministrator choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIsAdministrator(contractID string, args TokenAdminRegistryIsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IsAdministrator",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryProposeAdministrator exercises the TokenAdminRegistry_ProposeAdministrator choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryProposeAdministrator(contractID string, args TokenAdminRegistryProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ProposeAdministrator",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryIssueSendTicket exercises the TokenAdminRegistry_IssueSendTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIssueSendTicket(contractID string, args TokenAdminRegistryIssueSendTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueSendTicket",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryIssueReceiveTicket exercises the TokenAdminRegistry_IssueReceiveTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIssueReceiveTicket(contractID string, args TokenAdminRegistryIssueReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueReceiveTicket",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TokenAdminRegistryConsumeReceiveTicket exercises the TokenAdminRegistry_ConsumeReceiveTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryConsumeReceiveTicket(contractID string, args TokenAdminRegistryConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),

		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ConsumeReceiveTicket",

		Arguments: argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole is a Record type
type TokenAdminRegistryAcceptAdminRole struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryAcceptAdminRole to a map for DAML arguments
func (t TokenAdminRegistryAcceptAdminRole) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAcceptAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryConsumeReceiveTicket is a Record type
type TokenAdminRegistryConsumeReceiveTicket struct {
	TicketCid CONTRACT_ID `json:"ticketCid"`

	InstrumentId InstrumentId `json:"instrumentId"`

	PoolOwner PARTY `json:"poolOwner"`
}

// ToMap converts TokenAdminRegistryConsumeReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryConsumeReceiveTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ticketCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TicketCid
	}()

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["poolOwner"] = t.PoolOwner.ToMap()

	return m
}

func (t TokenAdminRegistryConsumeReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryGetTokenConfig is a Record type
type TokenAdminRegistryGetTokenConfig struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryGetTokenConfig to a map for DAML arguments
func (t TokenAdminRegistryGetTokenConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryGetTokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryGetTokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIsAdministrator is a Record type
type TokenAdminRegistryIsAdministrator struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Administrator PARTY `json:"administrator"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryIsAdministrator to a map for DAML arguments
func (t TokenAdminRegistryIsAdministrator) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["administrator"] = t.Administrator.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryIsAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryIsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIssueReceiveTicket is a Record type
type TokenAdminRegistryIssueReceiveTicket struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	PoolOwner PARTY `json:"poolOwner"`

	Receiver PARTY `json:"receiver"`

	TokenReceiver PARTY `json:"tokenReceiver"`

	Amount NUMERIC `json:"amount"`

	MessageHash TEXT `json:"messageHash"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryIssueReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryIssueReceiveTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["tokenReceiver"] = t.TokenReceiver.ToMap()

	m["amount"] = t.Amount

	m["messageHash"] = string(t.MessageHash)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryIssueReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryIssueReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIssueSendTicket is a Record type
type TokenAdminRegistryIssueSendTicket struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Sender PARTY `json:"sender"`

	Amount NUMERIC `json:"amount"`

	SourceTokenAddress TEXT `json:"sourceTokenAddress"`

	DestTokenAddress TEXT `json:"destTokenAddress"`

	TokenReceiver TEXT `json:"tokenReceiver"`

	ExtraData TEXT `json:"extraData"`

	Receipt Receipt `json:"receipt"`

	PoolOwner PARTY `json:"poolOwner"`
}

// ToMap converts TokenAdminRegistryIssueSendTicket to a map for DAML arguments
func (t TokenAdminRegistryIssueSendTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["sender"] = t.Sender.ToMap()

	m["amount"] = t.Amount

	m["sourceTokenAddress"] = string(t.SourceTokenAddress)

	m["destTokenAddress"] = string(t.DestTokenAddress)

	m["tokenReceiver"] = string(t.TokenReceiver)

	m["extraData"] = string(t.ExtraData)

	m["receipt"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Receipt).(mapper); ok {
			return m.toMap()
		}
		return t.Receipt
	}()

	m["poolOwner"] = t.PoolOwner.ToMap()

	return m
}

func (t TokenAdminRegistryIssueSendTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryIssueSendTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryProposeAdministrator is a Record type
type TokenAdminRegistryProposeAdministrator struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	NewAdmin PARTY `json:"newAdmin"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryProposeAdministrator to a map for DAML arguments
func (t TokenAdminRegistryProposeAdministrator) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryProposeAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	OptTokenPoolOwner *PARTY `json:"optTokenPoolOwner"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetPool to a map for DAML arguments
func (t TokenAdminRegistrySetPool) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	if t.OptTokenPoolOwner != nil {
		m["optTokenPoolOwner"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.OptTokenPoolOwner).ToMap(),
		}
	} else {
		m["optTokenPoolOwner"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistrySetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryTransferAdminRole is a Record type
type TokenAdminRegistryTransferAdminRole struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	NewAdmin PARTY `json:"newAdmin"`

	Caller PARTY `json:"caller"`
}

// ToMap converts TokenAdminRegistryTransferAdminRole to a map for DAML arguments
func (t TokenAdminRegistryTransferAdminRole) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["newAdmin"] = t.NewAdmin.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryTransferAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistryTransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenConfig is a Record type
type TokenConfig struct {
	Admin *PARTY `json:"admin"`

	PendingAdmin *PARTY `json:"pendingAdmin"`

	TokenPoolOwner *PARTY `json:"tokenPoolOwner"`
}

// ToMap converts TokenConfig to a map for DAML arguments
func (t TokenConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	if t.Admin != nil {
		m["admin"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.Admin).ToMap(),
		}
	} else {
		m["admin"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.PendingAdmin != nil {
		m["pendingAdmin"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.PendingAdmin).ToMap(),
		}
	} else {
		m["pendingAdmin"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	if t.TokenPoolOwner != nil {
		m["tokenPoolOwner"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.TokenPoolOwner).ToMap(),
		}
	} else {
		m["tokenPoolOwner"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

func (t TokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
