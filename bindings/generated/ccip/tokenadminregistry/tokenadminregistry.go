package tokenadminregistry

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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
	PackageName = "ccip-tokenadminregistry"
	PackageID   = "d17ab23962b4c2814144b9b415f71bdee31d813dcb28ef9218fe602a32207667"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
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

// PoolRegistration is a Record type
type PoolRegistration struct {
	PoolOwner      types.PARTY `json:"poolOwner"`
	PoolInstanceId types.TEXT  `json:"poolInstanceId"`
}

// ToMap converts PoolRegistration to a map for DAML arguments
func (t PoolRegistration) ToMap() map[string]any {
	m := make(map[string]any)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	return m
}

func (t PoolRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PoolRegistration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PoolRegistration to hex string (Canton MCMS format)
func (t PoolRegistration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PoolRegistration from hex string (Canton MCMS format)
func (t *PoolRegistration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistry is a Template type
type TokenAdminRegistry struct {
	InstanceId   types.TEXT   `json:"instanceId"`
	Owner        types.PARTY  `json:"owner"`
	TokenConfigs types.GENMAP `json:"tokenConfigs"`
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
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() any {
		if t.TokenConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TokenAdminRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenConfigs"] = func() any {
		if t.TokenConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TokenConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistry to hex string (Canton MCMS format)
func (t TokenAdminRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistry from hex string (Canton MCMS format)
func (t *TokenAdminRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
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
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryAcceptAdminRole to a map for DAML arguments
func (t TokenAdminRegistryAcceptAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAcceptAdminRole to hex string (Canton MCMS format)
func (t TokenAdminRegistryAcceptAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAcceptAdminRole from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAcceptAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryConsumeReceiveTicket is a Record type
type TokenAdminRegistryConsumeReceiveTicket struct {
	TicketCid    types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryConsumeReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryConsumeReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["ticketCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TicketCid
	}()

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryConsumeReceiveTicket to hex string (Canton MCMS format)
func (t TokenAdminRegistryConsumeReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryConsumeReceiveTicket from hex string (Canton MCMS format)
func (t *TokenAdminRegistryConsumeReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryGetTokenConfig is a Record type
type TokenAdminRegistryGetTokenConfig struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryGetTokenConfig to a map for DAML arguments
func (t TokenAdminRegistryGetTokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryGetTokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryGetTokenConfig to hex string (Canton MCMS format)
func (t TokenAdminRegistryGetTokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryGetTokenConfig from hex string (Canton MCMS format)
func (t *TokenAdminRegistryGetTokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryIsAdministrator is a Record type
type TokenAdminRegistryIsAdministrator struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Administrator types.PARTY                              `json:"administrator"`
	Caller        types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryIsAdministrator to a map for DAML arguments
func (t TokenAdminRegistryIsAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryIsAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryIsAdministrator to hex string (Canton MCMS format)
func (t TokenAdminRegistryIsAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryIsAdministrator from hex string (Canton MCMS format)
func (t *TokenAdminRegistryIsAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryIssueReceiveTicket is a Record type
type TokenAdminRegistryIssueReceiveTicket struct {
	InstrumentId        splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	PoolOwner           types.PARTY                              `json:"poolOwner"`
	Receiver            types.PARTY                              `json:"receiver"`
	TokenReceiver       types.PARTY                              `json:"tokenReceiver"`
	Amount              types.NUMERIC                            `json:"amount"`
	MessageHash         types.TEXT                               `json:"messageHash"`
	SourceChainSelector types.NUMERIC                            `json:"sourceChainSelector"`
}

// ToMap converts TokenAdminRegistryIssueReceiveTicket to a map for DAML arguments
func (t TokenAdminRegistryIssueReceiveTicket) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryIssueReceiveTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryIssueReceiveTicket to hex string (Canton MCMS format)
func (t TokenAdminRegistryIssueReceiveTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryIssueReceiveTicket from hex string (Canton MCMS format)
func (t *TokenAdminRegistryIssueReceiveTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryProposeAdministrator is a Record type
type TokenAdminRegistryProposeAdministrator struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryProposeAdministrator to a map for DAML arguments
func (t TokenAdminRegistryProposeAdministrator) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryProposeAdministrator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryProposeAdministrator to hex string (Canton MCMS format)
func (t TokenAdminRegistryProposeAdministrator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryProposeAdministrator from hex string (Canton MCMS format)
func (t *TokenAdminRegistryProposeAdministrator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetInboundPoolCCVs is a Record type
type TokenAdminRegistrySetInboundPoolCCVs struct {
	ExecutingMessageCid types.CONTRACT_ID           `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                  `json:"poolInstanceId"`
	PoolCCVs            []common.RawInstanceAddress `json:"poolCCVs"`
	Caller              types.PARTY                 `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetInboundPoolCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetInboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["poolInstanceId"] = string(t.PoolInstanceId)

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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetInboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetInboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetInboundPoolCCVs to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetInboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetInboundPoolCCVs from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetInboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetPool to a map for DAML arguments
func (t TokenAdminRegistrySetPool) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenPool,
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistrySetPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetPool to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetPool from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryTransferAdminRole is a Record type
type TokenAdminRegistryTransferAdminRole struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
	Caller       types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryTransferAdminRole to a map for DAML arguments
func (t TokenAdminRegistryTransferAdminRole) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryTransferAdminRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryTransferAdminRole to hex string (Canton MCMS format)
func (t TokenAdminRegistryTransferAdminRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryTransferAdminRole from hex string (Canton MCMS format)
func (t *TokenAdminRegistryTransferAdminRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfig is a Record type
type TokenConfig struct {
	Admin        *types.PARTY      `json:"admin"`
	PendingAdmin *types.PARTY      `json:"pendingAdmin"`
	TokenPool    *PoolRegistration `json:"tokenPool"`
}

// ToMap converts TokenConfig to a map for DAML arguments
func (t TokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	if t.Admin != nil {
		m["admin"] = map[string]any{
			"_type": "optional",
			"value": (*t.Admin).ToMap(),
		}
	} else {
		m["admin"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.PendingAdmin != nil {
		m["pendingAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.PendingAdmin).ToMap(),
		}
	} else {
		m["pendingAdmin"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.TokenPool != nil {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenPool,
		}
	} else {
		m["tokenPool"] = map[string]any{
			"_type": "optional",
		}
	}

	return m
}

func (t TokenConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenConfig to hex string (Canton MCMS format)
func (t TokenConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenConfig from hex string (Canton MCMS format)
func (t *TokenConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	TokenAdminRegistryAcceptAdminRole(args TokenAdminRegistryAcceptAdminRole) (*bind.EncodedChoice, error)
	TokenAdminRegistryConsumeReceiveTicket(args TokenAdminRegistryConsumeReceiveTicket) (*bind.EncodedChoice, error)
	TokenAdminRegistryGetTokenConfig(args TokenAdminRegistryGetTokenConfig) (*bind.EncodedChoice, error)
	TokenAdminRegistryIsAdministrator(args TokenAdminRegistryIsAdministrator) (*bind.EncodedChoice, error)
	TokenAdminRegistryIssueReceiveTicket(args TokenAdminRegistryIssueReceiveTicket) (*bind.EncodedChoice, error)
	TokenAdminRegistryProposeAdministrator(args TokenAdminRegistryProposeAdministrator) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetInboundPoolCCVs(args TokenAdminRegistrySetInboundPoolCCVs) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetPool(args TokenAdminRegistrySetPool) (*bind.EncodedChoice, error)
	TokenAdminRegistryTransferAdminRole(args TokenAdminRegistryTransferAdminRole) (*bind.EncodedChoice, error)
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

// TokenAdminRegistryAcceptAdminRole encodes parameters for the TokenAdminRegistryAcceptAdminRole choice.
func (e *encoder) TokenAdminRegistryAcceptAdminRole(args TokenAdminRegistryAcceptAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAcceptAdminRole", args)
}

// TokenAdminRegistryConsumeReceiveTicket encodes parameters for the TokenAdminRegistryConsumeReceiveTicket choice.
func (e *encoder) TokenAdminRegistryConsumeReceiveTicket(args TokenAdminRegistryConsumeReceiveTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryConsumeReceiveTicket", args)
}

// TokenAdminRegistryGetTokenConfig encodes parameters for the TokenAdminRegistryGetTokenConfig choice.
func (e *encoder) TokenAdminRegistryGetTokenConfig(args TokenAdminRegistryGetTokenConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryGetTokenConfig", args)
}

// TokenAdminRegistryIsAdministrator encodes parameters for the TokenAdminRegistryIsAdministrator choice.
func (e *encoder) TokenAdminRegistryIsAdministrator(args TokenAdminRegistryIsAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryIsAdministrator", args)
}

// TokenAdminRegistryIssueReceiveTicket encodes parameters for the TokenAdminRegistryIssueReceiveTicket choice.
func (e *encoder) TokenAdminRegistryIssueReceiveTicket(args TokenAdminRegistryIssueReceiveTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryIssueReceiveTicket", args)
}

// TokenAdminRegistryProposeAdministrator encodes parameters for the TokenAdminRegistryProposeAdministrator choice.
func (e *encoder) TokenAdminRegistryProposeAdministrator(args TokenAdminRegistryProposeAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryProposeAdministrator", args)
}

// TokenAdminRegistrySetInboundPoolCCVs encodes parameters for the TokenAdminRegistrySetInboundPoolCCVs choice.
func (e *encoder) TokenAdminRegistrySetInboundPoolCCVs(args TokenAdminRegistrySetInboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetInboundPoolCCVs", args)
}

// TokenAdminRegistrySetPool encodes parameters for the TokenAdminRegistrySetPool choice.
func (e *encoder) TokenAdminRegistrySetPool(args TokenAdminRegistrySetPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetPool", args)
}

// TokenAdminRegistryTransferAdminRole encodes parameters for the TokenAdminRegistryTransferAdminRole choice.
func (e *encoder) TokenAdminRegistryTransferAdminRole(args TokenAdminRegistryTransferAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryTransferAdminRole", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
