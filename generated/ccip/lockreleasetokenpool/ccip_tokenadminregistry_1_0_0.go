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

// ConsumeReceiveTicketResult is a Record type
type ConsumeReceiveTicketResult struct {
	MessageHash         TEXT     `json:"messageHash"`
	SourceChainSelector NUMERIC  `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC  `json:"sequenceNumber"`
	HasTokenTransfer    BOOL     `json:"hasTokenTransfer"`
	TokenAmount         *NUMERIC `json:"tokenAmount"`
	DestTokenAddress    *TEXT    `json:"destTokenAddress"`
	TokenReceiver       *PARTY   `json:"tokenReceiver"`
}

// toMap converts ConsumeReceiveTicketResult to a map for DAML arguments
func (t ConsumeReceiveTicketResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageHash":         string(t.MessageHash),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"hasTokenTransfer":    bool(t.HasTokenTransfer),
		"tokenAmount":         (*big.Int)(*t.TokenAmount),
		"destTokenAddress":    string(*t.DestTokenAddress),
		"tokenReceiver":       (*t.TokenReceiver).ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for ConsumeReceiveTicketResult using JsonCodec
func (t ConsumeReceiveTicketResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ConsumeReceiveTicketResult using JsonCodec
func (t *ConsumeReceiveTicketResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	Owner        PARTY  `json:"owner"`
	InstanceId   TEXT   `json:"instanceId"`
	TokenConfigs GENMAP `json:"tokenConfigs"`
}

// GetTemplateID returns the template ID for this template
func (t TokenAdminRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t TokenAdminRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	if t.TokenConfigs != nil && len(t.TokenConfigs) > 0 {
		args["tokenConfigs"] = map[string]interface{}{"_type": "genmap", "value": t.TokenConfigs}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistry using JsonCodec
func (t TokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistry using JsonCodec
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
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetPool exercises the TokenAdminRegistry_SetPool choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistrySetPool(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetPool",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole exercises the TokenAdminRegistry_AcceptAdminRole choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryAcceptAdminRole(contractID string, args TokenAdminRegistryAcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryTransferAdminRole exercises the TokenAdminRegistry_TransferAdminRole choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryTransferAdminRole(contractID string, args TokenAdminRegistryTransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIsAdministrator exercises the TokenAdminRegistry_IsAdministrator choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIsAdministrator(contractID string, args TokenAdminRegistryIsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryProposeAdministrator exercises the TokenAdminRegistry_ProposeAdministrator choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryProposeAdministrator(contractID string, args TokenAdminRegistryProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetRequiredCCVs exercises the TokenAdminRegistry_SetRequiredCCVs choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistrySetRequiredCCVs(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIssueSendTicket exercises the TokenAdminRegistry_IssueSendTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIssueSendTicket(contractID string, args TokenAdminRegistryIssueSendTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueSendTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIssueReceiveTicket exercises the TokenAdminRegistry_IssueReceiveTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryIssueReceiveTicket(contractID string, args TokenAdminRegistryIssueReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TokenAdminRegistryConsumeReceiveTicket exercises the TokenAdminRegistry_ConsumeReceiveTicket choice on this TokenAdminRegistry contract
func (t TokenAdminRegistry) TokenAdminRegistryConsumeReceiveTicket(contractID string, args TokenAdminRegistryConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole is a Record type
type TokenAdminRegistryAcceptAdminRole struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistryAcceptAdminRole to a map for DAML arguments
func (t TokenAdminRegistryAcceptAdminRole) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryAcceptAdminRole using JsonCodec
func (t TokenAdminRegistryAcceptAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryAcceptAdminRole using JsonCodec
func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryConsumeReceiveTicket is a Record type
type TokenAdminRegistryConsumeReceiveTicket struct {
	TicketCid    CONTRACT_ID  `json:"ticketCid"`
	InstrumentId InstrumentId `json:"instrumentId"`
	PoolOwner    PARTY        `json:"poolOwner"`
}

// toMap converts TokenAdminRegistryConsumeReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryConsumeReceiveTicket) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ticketCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TicketCid).(mapper); ok {
				return m.toMap()
			}
			return t.TicketCid
		}(),
		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"poolOwner": t.PoolOwner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryConsumeReceiveTicket using JsonCodec
func (t TokenAdminRegistryConsumeReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryConsumeReceiveTicket using JsonCodec
func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryGetTokenConfig is a Record type
type TokenAdminRegistryGetTokenConfig struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistryGetTokenConfig to a map for DAML arguments
func (t TokenAdminRegistryGetTokenConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryGetTokenConfig using JsonCodec
func (t TokenAdminRegistryGetTokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryGetTokenConfig using JsonCodec
func (t *TokenAdminRegistryGetTokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIsAdministrator is a Record type
type TokenAdminRegistryIsAdministrator struct {
	InstrumentId  InstrumentId `json:"instrumentId"`
	Administrator PARTY        `json:"administrator"`
	Caller        PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistryIsAdministrator to a map for DAML arguments
func (t TokenAdminRegistryIsAdministrator) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"administrator": t.Administrator.ToMap(),
		"caller":        t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryIsAdministrator using JsonCodec
func (t TokenAdminRegistryIsAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryIsAdministrator using JsonCodec
func (t *TokenAdminRegistryIsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIssueReceiveTicket is a Record type
type TokenAdminRegistryIssueReceiveTicket struct {
	MessageHash         TEXT     `json:"messageHash"`
	SourceChainSelector NUMERIC  `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC  `json:"sequenceNumber"`
	HasTokenTransfer    BOOL     `json:"hasTokenTransfer"`
	TokenAmount         *NUMERIC `json:"tokenAmount"`
	DestTokenAddress    *TEXT    `json:"destTokenAddress"`
	TokenReceiver       *PARTY   `json:"tokenReceiver"`
	Receiver            PARTY    `json:"receiver"`
	Caller              PARTY    `json:"caller"`
}

// toMap converts TokenAdminRegistryIssueReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryIssueReceiveTicket) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageHash":         string(t.MessageHash),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"hasTokenTransfer":    bool(t.HasTokenTransfer),
		"tokenAmount":         (*big.Int)(*t.TokenAmount),
		"destTokenAddress":    string(*t.DestTokenAddress),
		"tokenReceiver":       (*t.TokenReceiver).ToMap(),
		"receiver":            t.Receiver.ToMap(),
		"caller":              t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryIssueReceiveTicket using JsonCodec
func (t TokenAdminRegistryIssueReceiveTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryIssueReceiveTicket using JsonCodec
func (t *TokenAdminRegistryIssueReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryIssueSendTicket is a Record type
type TokenAdminRegistryIssueSendTicket struct {
	InstrumentId       InstrumentId `json:"instrumentId"`
	Sender             PARTY        `json:"sender"`
	Amount             NUMERIC      `json:"amount"`
	SourceTokenAddress TEXT         `json:"sourceTokenAddress"`
	DestTokenAddress   TEXT         `json:"destTokenAddress"`
	TokenReceiver      TEXT         `json:"tokenReceiver"`
	ExtraData          TEXT         `json:"extraData"`
	PoolOwner          PARTY        `json:"poolOwner"`
}

// toMap converts TokenAdminRegistryIssueSendTicket to a map for DAML arguments
func (t TokenAdminRegistryIssueSendTicket) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"sender":             t.Sender.ToMap(),
		"amount":             (*big.Int)(t.Amount),
		"sourceTokenAddress": string(t.SourceTokenAddress),
		"destTokenAddress":   string(t.DestTokenAddress),
		"tokenReceiver":      string(t.TokenReceiver),
		"extraData":          string(t.ExtraData),
		"poolOwner":          t.PoolOwner.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryIssueSendTicket using JsonCodec
func (t TokenAdminRegistryIssueSendTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryIssueSendTicket using JsonCodec
func (t *TokenAdminRegistryIssueSendTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryProposeAdministrator is a Record type
type TokenAdminRegistryProposeAdministrator struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	NewAdmin     PARTY        `json:"newAdmin"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistryProposeAdministrator to a map for DAML arguments
func (t TokenAdminRegistryProposeAdministrator) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"newAdmin": t.NewAdmin.ToMap(),
		"caller":   t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryProposeAdministrator using JsonCodec
func (t TokenAdminRegistryProposeAdministrator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryProposeAdministrator using JsonCodec
func (t *TokenAdminRegistryProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	InstrumentId      InstrumentId `json:"instrumentId"`
	OptTokenPoolOwner *PARTY       `json:"optTokenPoolOwner"`
	Caller            PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistrySetPool to a map for DAML arguments
func (t TokenAdminRegistrySetPool) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"optTokenPoolOwner": (*t.OptTokenPoolOwner).ToMap(),
		"caller":            t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistrySetPool using JsonCodec
func (t TokenAdminRegistrySetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistrySetPool using JsonCodec
func (t *TokenAdminRegistrySetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistrySetRequiredCCVs is a Record type
type TokenAdminRegistrySetRequiredCCVs struct {
	InstrumentId    InstrumentId `json:"instrumentId"`
	NewRequiredCCVs []TEXT       `json:"newRequiredCCVs"`
	Caller          PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistrySetRequiredCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetRequiredCCVs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"newRequiredCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.NewRequiredCCVs))
			for _, e := range t.NewRequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistrySetRequiredCCVs using JsonCodec
func (t TokenAdminRegistrySetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistrySetRequiredCCVs using JsonCodec
func (t *TokenAdminRegistrySetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistryTransferAdminRole is a Record type
type TokenAdminRegistryTransferAdminRole struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	NewAdmin     PARTY        `json:"newAdmin"`
	Caller       PARTY        `json:"caller"`
}

// toMap converts TokenAdminRegistryTransferAdminRole to a map for DAML arguments
func (t TokenAdminRegistryTransferAdminRole) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"newAdmin": t.NewAdmin.ToMap(),
		"caller":   t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenAdminRegistryTransferAdminRole using JsonCodec
func (t TokenAdminRegistryTransferAdminRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenAdminRegistryTransferAdminRole using JsonCodec
func (t *TokenAdminRegistryTransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenConfig is a Record type
type TokenConfig struct {
	Admin          *PARTY `json:"admin"`
	PendingAdmin   *PARTY `json:"pendingAdmin"`
	TokenPoolOwner *PARTY `json:"tokenPoolOwner"`
	RequiredCCVs   []TEXT `json:"requiredCCVs"`
}

// toMap converts TokenConfig to a map for DAML arguments
func (t TokenConfig) toMap() map[string]interface{} {
	return map[string]interface{}{

		"admin":          (*t.Admin).ToMap(),
		"pendingAdmin":   (*t.PendingAdmin).ToMap(),
		"tokenPoolOwner": (*t.TokenPoolOwner).ToMap(),
		"requiredCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.RequiredCCVs))
			for _, e := range t.RequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenConfig using JsonCodec
func (t TokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenConfig using JsonCodec
func (t *TokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
