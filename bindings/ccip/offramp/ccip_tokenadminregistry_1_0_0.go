package offramp

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

// PoolRegistration is a Record type
type PoolRegistration struct {
	PoolOwner      PARTY `json:"poolOwner"`
	PoolInstanceId TEXT  `json:"poolInstanceId"`
}

// ToMap converts PoolRegistration to a map for DAML arguments
func (t PoolRegistration) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	return m
}

func (t PoolRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PoolRegistration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	InstanceId   TEXT   `json:"instanceId"`
	Owner        PARTY  `json:"owner"`
	TokenConfigs GENMAP `json:"tokenConfigs"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

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

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenAdminRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() interface{} {
		if t.TokenConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
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
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryGetTokenConfig(contractID string, args TokenAdminRegistryGetTokenConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_GetTokenConfig",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryGetTokenConfigWithPackageID exercises the TokenAdminRegistry_GetTokenConfig choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryGetTokenConfigWithPackageID(contractID string, packageID string, args TokenAdminRegistryGetTokenConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_GetTokenConfig",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetPool exercises the TokenAdminRegistry_SetPool choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistrySetPool(contractID string, args TokenAdminRegistrySetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetPool",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetPoolWithPackageID exercises the TokenAdminRegistry_SetPool choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistrySetPoolWithPackageID(contractID string, packageID string, args TokenAdminRegistrySetPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetPool",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole exercises the TokenAdminRegistry_AcceptAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryAcceptAdminRole(contractID string, args TokenAdminRegistryAcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRoleWithPackageID exercises the TokenAdminRegistry_AcceptAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryAcceptAdminRoleWithPackageID(contractID string, packageID string, args TokenAdminRegistryAcceptAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AcceptAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryTransferAdminRole exercises the TokenAdminRegistry_TransferAdminRole choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryTransferAdminRole(contractID string, args TokenAdminRegistryTransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryTransferAdminRoleWithPackageID exercises the TokenAdminRegistry_TransferAdminRole choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryTransferAdminRoleWithPackageID(contractID string, packageID string, args TokenAdminRegistryTransferAdminRole) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_TransferAdminRole",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIsAdministrator exercises the TokenAdminRegistry_IsAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryIsAdministrator(contractID string, args TokenAdminRegistryIsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIsAdministratorWithPackageID exercises the TokenAdminRegistry_IsAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryIsAdministratorWithPackageID(contractID string, packageID string, args TokenAdminRegistryIsAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IsAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryProposeAdministrator exercises the TokenAdminRegistry_ProposeAdministrator choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryProposeAdministrator(contractID string, args TokenAdminRegistryProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryProposeAdministratorWithPackageID exercises the TokenAdminRegistry_ProposeAdministrator choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryProposeAdministratorWithPackageID(contractID string, packageID string, args TokenAdminRegistryProposeAdministrator) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ProposeAdministrator",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIssueReceiveTicket exercises the TokenAdminRegistry_IssueReceiveTicket choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryIssueReceiveTicket(contractID string, args TokenAdminRegistryIssueReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryIssueReceiveTicketWithPackageID exercises the TokenAdminRegistry_IssueReceiveTicket choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryIssueReceiveTicketWithPackageID(contractID string, packageID string, args TokenAdminRegistryIssueReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_IssueReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryConsumeReceiveTicket exercises the TokenAdminRegistry_ConsumeReceiveTicket choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryConsumeReceiveTicket(contractID string, args TokenAdminRegistryConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryConsumeReceiveTicketWithPackageID exercises the TokenAdminRegistry_ConsumeReceiveTicket choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryConsumeReceiveTicketWithPackageID(contractID string, packageID string, args TokenAdminRegistryConsumeReceiveTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_ConsumeReceiveTicket",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TokenAdminRegistrySetInboundPoolCCVs exercises the TokenAdminRegistry_SetInboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistrySetInboundPoolCCVs(contractID string, args TokenAdminRegistrySetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetInboundPoolCCVsWithPackageID exercises the TokenAdminRegistry_SetInboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistrySetInboundPoolCCVsWithPackageID(contractID string, packageID string, args TokenAdminRegistrySetInboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetInboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAcceptAdminRole is a Record type
type TokenAdminRegistryAcceptAdminRole struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
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
	TicketCid    CONTRACT_ID  `json:"ticketCid"`
	InstrumentId InstrumentId `json:"instrumentId"`
	Caller       PARTY        `json:"caller"`
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

	m["caller"] = t.Caller.ToMap()

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
	Caller       PARTY        `json:"caller"`
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
	InstrumentId  InstrumentId `json:"instrumentId"`
	Administrator PARTY        `json:"administrator"`
	Caller        PARTY        `json:"caller"`
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
	InstrumentId        InstrumentId `json:"instrumentId"`
	PoolOwner           PARTY        `json:"poolOwner"`
	Receiver            PARTY        `json:"receiver"`
	TokenReceiver       PARTY        `json:"tokenReceiver"`
	Amount              NUMERIC      `json:"amount"`
	MessageHash         TEXT         `json:"messageHash"`
	SourceChainSelector NUMERIC      `json:"sourceChainSelector"`
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

// TokenAdminRegistryProposeAdministrator is a Record type
type TokenAdminRegistryProposeAdministrator struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	NewAdmin     PARTY        `json:"newAdmin"`
	Caller       PARTY        `json:"caller"`
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

// TokenAdminRegistrySetInboundPoolCCVs is a Record type
type TokenAdminRegistrySetInboundPoolCCVs struct {
	ExecutingMessageCid CONTRACT_ID          `json:"executingMessageCid"`
	PoolInstanceId      TEXT                 `json:"poolInstanceId"`
	PoolCCVs            []RawInstanceAddress `json:"poolCCVs"`
	Caller              PARTY                `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetInboundPoolCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetInboundPoolCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.PoolCCVs))
		for _, e := range t.PoolCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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

func (t TokenAdminRegistrySetInboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenAdminRegistrySetInboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	InstrumentId InstrumentId      `json:"instrumentId"`
	TokenPool    *PoolRegistration `json:"tokenPool"`
	Caller       PARTY             `json:"caller"`
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

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenPool,
		}
	} else {
		m["tokenPool"] = map[string]interface{}{
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
	NewAdmin     PARTY        `json:"newAdmin"`
	Caller       PARTY        `json:"caller"`
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
	Admin        *PARTY            `json:"admin"`
	PendingAdmin *PARTY            `json:"pendingAdmin"`
	TokenPool    *PoolRegistration `json:"tokenPool"`
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

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenPool,
		}
	} else {
		m["tokenPool"] = map[string]interface{}{
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
