package tokenadminregistry

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
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
	PackageID   = "53df6801d3850e6c0ca52ebf82013fb8974175ec41aa56e4f9f7fc44ee5d72bd"
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

// Get2 is a Record type
type Get2 struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts Get2 to a map for DAML arguments
func (t Get2) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t Get2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Get2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Get2 to hex string (Canton MCMS format)
func (t Get2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get2 from hex string (Canton MCMS format)
func (t *Get2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Get2MCMSParams is Get2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type Get2MCMSParams struct {
}

// MarshalHex encodes Get2MCMSParams to hex string for MCMS operationData.
func (t Get2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Get2MCMSParams from hex string.
func (t *Get2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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

// TokenAdminRegistrySetOutboundPoolCCVs exercises the TokenAdminRegistry_SetOutboundPoolCCVs choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistrySetOutboundPoolCCVs(contractID string, args TokenAdminRegistrySetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistrySetOutboundPoolCCVsWithPackageID exercises the TokenAdminRegistry_SetOutboundPoolCCVs choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistrySetOutboundPoolCCVsWithPackageID(contractID string, packageID string, args TokenAdminRegistrySetOutboundPoolCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_SetOutboundPoolCCVs",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAddTokenSendFee exercises the TokenAdminRegistry_AddTokenSendFee choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryAddTokenSendFee(contractID string, args TokenAdminRegistryAddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAddTokenSendFeeWithPackageID exercises the TokenAdminRegistry_AddTokenSendFee choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryAddTokenSendFeeWithPackageID(contractID string, packageID string, args TokenAdminRegistryAddTokenSendFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AddTokenSendFee",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAddTokenSend exercises the TokenAdminRegistry_AddTokenSend choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryAddTokenSend(contractID string, args TokenAdminRegistryAddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AddTokenSend",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryAddTokenSendWithPackageID exercises the TokenAdminRegistry_AddTokenSend choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryAddTokenSendWithPackageID(contractID string, packageID string, args TokenAdminRegistryAddTokenSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_AddTokenSend",
		Arguments:  argsToMap(args),
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

// TokenAdminRegistryFinalizeExecute exercises the TokenAdminRegistry_FinalizeExecute choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) TokenAdminRegistryFinalizeExecute(contractID string, args TokenAdminRegistryFinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

// TokenAdminRegistryFinalizeExecuteWithPackageID exercises the TokenAdminRegistry_FinalizeExecute choice using the provided package ID instead of package name
func (t TokenAdminRegistry) TokenAdminRegistryFinalizeExecuteWithPackageID(contractID string, packageID string, args TokenAdminRegistryFinalizeExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "TokenAdminRegistry_FinalizeExecute",
		Arguments:  argsToMap(args),
	}
}

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

// Archive exercises the Archive choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TokenAdminRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Get exercises the Get choice on this TokenAdminRegistry contract
// This method uses the package name in the template ID
func (t TokenAdminRegistry) Get(contractID string, args Get2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// GetWithPackageID exercises the Get choice using the provided package ID instead of package name
func (t TokenAdminRegistry) GetWithPackageID(contractID string, packageID string, args Get2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "TokenAdminRegistry"),
		ContractID: contractID,
		Choice:     "Get",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this TokenAdminRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t TokenAdminRegistry) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t TokenAdminRegistry) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.TokenAdminRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for TokenAdminRegistry

var _ mcms.IMCMSReceiver = (*TokenAdminRegistry)(nil)

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

// TokenAdminRegistryAcceptAdminRoleMCMSParams is TokenAdminRegistryAcceptAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryAcceptAdminRoleMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes TokenAdminRegistryAcceptAdminRoleMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryAcceptAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAcceptAdminRoleMCMSParams from hex string.
func (t *TokenAdminRegistryAcceptAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSend is a Record type
type TokenAdminRegistryAddTokenSend struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                               `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount            types.NUMERIC                            `json:"amount"`
	DestTokenAddress  types.TEXT                               `json:"destTokenAddress"`
	ExtraData         types.TEXT                               `json:"extraData"`
	Caller            types.PARTY                              `json:"caller"`
}

// ToMap converts TokenAdminRegistryAddTokenSend to a map for DAML arguments
func (t TokenAdminRegistryAddTokenSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["poolInstanceId"] = string(t.PoolInstanceId)

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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAddTokenSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAddTokenSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAddTokenSend to hex string (Canton MCMS format)
func (t TokenAdminRegistryAddTokenSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSend from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAddTokenSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSendMCMSParams is TokenAdminRegistryAddTokenSend without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryAddTokenSendMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID                        `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                               `json:"poolInstanceId"`
	InstrumentId      splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount            types.NUMERIC                            `json:"amount"`
	DestTokenAddress  types.TEXT                               `json:"destTokenAddress"`
	ExtraData         types.TEXT                               `json:"extraData"`
}

// MarshalHex encodes TokenAdminRegistryAddTokenSendMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryAddTokenSendMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSendMCMSParams from hex string.
func (t *TokenAdminRegistryAddTokenSendMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSendFee is a Record type
type TokenAdminRegistryAddTokenSendFee struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT        `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC     `json:"feeUSDCents"`
	DestGasOverhead   types.INT64       `json:"destGasOverhead"`
	DestBytesOverhead types.INT64       `json:"destBytesOverhead"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts TokenAdminRegistryAddTokenSendFee to a map for DAML arguments
func (t TokenAdminRegistryAddTokenSendFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["feeUSDCents"] = t.FeeUSDCents

	m["destGasOverhead"] = int64(t.DestGasOverhead)

	m["destBytesOverhead"] = int64(t.DestBytesOverhead)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryAddTokenSendFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryAddTokenSendFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryAddTokenSendFee to hex string (Canton MCMS format)
func (t TokenAdminRegistryAddTokenSendFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSendFee from hex string (Canton MCMS format)
func (t *TokenAdminRegistryAddTokenSendFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryAddTokenSendFeeMCMSParams is TokenAdminRegistryAddTokenSendFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryAddTokenSendFeeMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT        `json:"poolInstanceId"`
	FeeUSDCents       types.NUMERIC     `json:"feeUSDCents"`
	DestGasOverhead   types.INT64       `json:"destGasOverhead"`
	DestBytesOverhead types.INT64       `json:"destBytesOverhead"`
}

// MarshalHex encodes TokenAdminRegistryAddTokenSendFeeMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryAddTokenSendFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryAddTokenSendFeeMCMSParams from hex string.
func (t *TokenAdminRegistryAddTokenSendFeeMCMSParams) UnmarshalHex(data string) error {
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

// TokenAdminRegistryConsumeReceiveTicketMCMSParams is TokenAdminRegistryConsumeReceiveTicket without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryConsumeReceiveTicketMCMSParams struct {
	TicketCid    types.CONTRACT_ID                        `json:"ticketCid"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes TokenAdminRegistryConsumeReceiveTicketMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryConsumeReceiveTicketMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryConsumeReceiveTicketMCMSParams from hex string.
func (t *TokenAdminRegistryConsumeReceiveTicketMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryFinalizeExecute is a Record type
type TokenAdminRegistryFinalizeExecute struct {
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	TicketReceiver      types.PARTY       `json:"ticketReceiver"`
	ReturnData          types.TEXT        `json:"returnData"`
	Caller              types.PARTY       `json:"caller"`
}

// ToMap converts TokenAdminRegistryFinalizeExecute to a map for DAML arguments
func (t TokenAdminRegistryFinalizeExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["ticketReceiver"] = t.TicketReceiver.ToMap()

	m["returnData"] = string(t.ReturnData)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenAdminRegistryFinalizeExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistryFinalizeExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistryFinalizeExecute to hex string (Canton MCMS format)
func (t TokenAdminRegistryFinalizeExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryFinalizeExecute from hex string (Canton MCMS format)
func (t *TokenAdminRegistryFinalizeExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistryFinalizeExecuteMCMSParams is TokenAdminRegistryFinalizeExecute without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryFinalizeExecuteMCMSParams struct {
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	TicketReceiver      types.PARTY       `json:"ticketReceiver"`
	ReturnData          types.TEXT        `json:"returnData"`
}

// MarshalHex encodes TokenAdminRegistryFinalizeExecuteMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryFinalizeExecuteMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryFinalizeExecuteMCMSParams from hex string.
func (t *TokenAdminRegistryFinalizeExecuteMCMSParams) UnmarshalHex(data string) error {
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

// TokenAdminRegistryGetTokenConfigMCMSParams is TokenAdminRegistryGetTokenConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryGetTokenConfigMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// MarshalHex encodes TokenAdminRegistryGetTokenConfigMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryGetTokenConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryGetTokenConfigMCMSParams from hex string.
func (t *TokenAdminRegistryGetTokenConfigMCMSParams) UnmarshalHex(data string) error {
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

// TokenAdminRegistryIsAdministratorMCMSParams is TokenAdminRegistryIsAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryIsAdministratorMCMSParams struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Administrator types.PARTY                              `json:"administrator"`
}

// MarshalHex encodes TokenAdminRegistryIsAdministratorMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryIsAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryIsAdministratorMCMSParams from hex string.
func (t *TokenAdminRegistryIsAdministratorMCMSParams) UnmarshalHex(data string) error {
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

// TokenAdminRegistryProposeAdministratorMCMSParams is TokenAdminRegistryProposeAdministrator without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryProposeAdministratorMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes TokenAdminRegistryProposeAdministratorMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryProposeAdministratorMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryProposeAdministratorMCMSParams from hex string.
func (t *TokenAdminRegistryProposeAdministratorMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetInboundPoolCCVs is a Record type
type TokenAdminRegistrySetInboundPoolCCVs struct {
	ExecutingMessageCid types.CONTRACT_ID         `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                `json:"poolInstanceId"`
	PoolCCVs            []mcms.RawInstanceAddress `json:"poolCCVs"`
	Caller              types.PARTY               `json:"caller"`
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

// TokenAdminRegistrySetInboundPoolCCVsMCMSParams is TokenAdminRegistrySetInboundPoolCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistrySetInboundPoolCCVsMCMSParams struct {
	ExecutingMessageCid types.CONTRACT_ID         `json:"executingMessageCid"`
	PoolInstanceId      types.TEXT                `json:"poolInstanceId"`
	PoolCCVs            []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// MarshalHex encodes TokenAdminRegistrySetInboundPoolCCVsMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistrySetInboundPoolCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetInboundPoolCCVsMCMSParams from hex string.
func (t *TokenAdminRegistrySetInboundPoolCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetOutboundPoolCCVs is a Record type
type TokenAdminRegistrySetOutboundPoolCCVs struct {
	SendingMessageCid types.CONTRACT_ID         `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                `json:"poolInstanceId"`
	PoolCCVs          []mcms.RawInstanceAddress `json:"poolCCVs"`
	Caller            types.PARTY               `json:"caller"`
}

// ToMap converts TokenAdminRegistrySetOutboundPoolCCVs to a map for DAML arguments
func (t TokenAdminRegistrySetOutboundPoolCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
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

func (t TokenAdminRegistrySetOutboundPoolCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenAdminRegistrySetOutboundPoolCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenAdminRegistrySetOutboundPoolCCVs to hex string (Canton MCMS format)
func (t TokenAdminRegistrySetOutboundPoolCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetOutboundPoolCCVs from hex string (Canton MCMS format)
func (t *TokenAdminRegistrySetOutboundPoolCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetOutboundPoolCCVsMCMSParams is TokenAdminRegistrySetOutboundPoolCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistrySetOutboundPoolCCVsMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID         `json:"sendingMessageCid"`
	PoolInstanceId    types.TEXT                `json:"poolInstanceId"`
	PoolCCVs          []mcms.RawInstanceAddress `json:"poolCCVs"`
}

// MarshalHex encodes TokenAdminRegistrySetOutboundPoolCCVsMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistrySetOutboundPoolCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetOutboundPoolCCVsMCMSParams from hex string.
func (t *TokenAdminRegistrySetOutboundPoolCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenAdminRegistrySetPool is a Record type
type TokenAdminRegistrySetPool struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
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

// TokenAdminRegistrySetPoolMCMSParams is TokenAdminRegistrySetPool without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistrySetPoolMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// MarshalHex encodes TokenAdminRegistrySetPoolMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistrySetPoolMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistrySetPoolMCMSParams from hex string.
func (t *TokenAdminRegistrySetPoolMCMSParams) UnmarshalHex(data string) error {
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

// TokenAdminRegistryTransferAdminRoleMCMSParams is TokenAdminRegistryTransferAdminRole without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type TokenAdminRegistryTransferAdminRoleMCMSParams struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	NewAdmin     types.PARTY                              `json:"newAdmin"`
}

// MarshalHex encodes TokenAdminRegistryTransferAdminRoleMCMSParams to hex string for MCMS operationData.
func (t TokenAdminRegistryTransferAdminRoleMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenAdminRegistryTransferAdminRoleMCMSParams from hex string.
func (t *TokenAdminRegistryTransferAdminRoleMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenConfig is a Record type
type TokenConfig struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Admin        *types.PARTY                             `json:"admin" hex:"optional"`
	PendingAdmin *types.PARTY                             `json:"pendingAdmin" hex:"optional"`
	TokenPool    *PoolRegistration                        `json:"tokenPool" hex:"optional"`
}

// ToMap converts TokenConfig to a map for DAML arguments
func (t TokenConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

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
	Get2(args Get2) (*bind.EncodedChoice, error)
	Get2MCMSParams(args Get2MCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryAcceptAdminRole(args TokenAdminRegistryAcceptAdminRole) (*bind.EncodedChoice, error)
	TokenAdminRegistryAcceptAdminRoleMCMSParams(args TokenAdminRegistryAcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryAddTokenSend(args TokenAdminRegistryAddTokenSend) (*bind.EncodedChoice, error)
	TokenAdminRegistryAddTokenSendMCMSParams(args TokenAdminRegistryAddTokenSendMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryAddTokenSendFee(args TokenAdminRegistryAddTokenSendFee) (*bind.EncodedChoice, error)
	TokenAdminRegistryAddTokenSendFeeMCMSParams(args TokenAdminRegistryAddTokenSendFeeMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryConsumeReceiveTicket(args TokenAdminRegistryConsumeReceiveTicket) (*bind.EncodedChoice, error)
	TokenAdminRegistryConsumeReceiveTicketMCMSParams(args TokenAdminRegistryConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryFinalizeExecute(args TokenAdminRegistryFinalizeExecute) (*bind.EncodedChoice, error)
	TokenAdminRegistryFinalizeExecuteMCMSParams(args TokenAdminRegistryFinalizeExecuteMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryGetTokenConfig(args TokenAdminRegistryGetTokenConfig) (*bind.EncodedChoice, error)
	TokenAdminRegistryGetTokenConfigMCMSParams(args TokenAdminRegistryGetTokenConfigMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryIsAdministrator(args TokenAdminRegistryIsAdministrator) (*bind.EncodedChoice, error)
	TokenAdminRegistryIsAdministratorMCMSParams(args TokenAdminRegistryIsAdministratorMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryProposeAdministrator(args TokenAdminRegistryProposeAdministrator) (*bind.EncodedChoice, error)
	TokenAdminRegistryProposeAdministratorMCMSParams(args TokenAdminRegistryProposeAdministratorMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetInboundPoolCCVs(args TokenAdminRegistrySetInboundPoolCCVs) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetInboundPoolCCVsMCMSParams(args TokenAdminRegistrySetInboundPoolCCVsMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetOutboundPoolCCVs(args TokenAdminRegistrySetOutboundPoolCCVs) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetOutboundPoolCCVsMCMSParams(args TokenAdminRegistrySetOutboundPoolCCVsMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetPool(args TokenAdminRegistrySetPool) (*bind.EncodedChoice, error)
	TokenAdminRegistrySetPoolMCMSParams(args TokenAdminRegistrySetPoolMCMSParams) (*bind.EncodedChoice, error)
	TokenAdminRegistryTransferAdminRole(args TokenAdminRegistryTransferAdminRole) (*bind.EncodedChoice, error)
	TokenAdminRegistryTransferAdminRoleMCMSParams(args TokenAdminRegistryTransferAdminRoleMCMSParams) (*bind.EncodedChoice, error)
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

// Get2 encodes parameters for the Get2 choice.
func (e *encoder) Get2(args Get2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get2", args)
}

// Get2MCMSParams encodes MCMS parameters (without Caller) for the Get2 choice.
func (e *encoder) Get2MCMSParams(args Get2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Get2", args)
}

// TokenAdminRegistryAcceptAdminRole encodes parameters for the TokenAdminRegistryAcceptAdminRole choice.
func (e *encoder) TokenAdminRegistryAcceptAdminRole(args TokenAdminRegistryAcceptAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAcceptAdminRole", args)
}

// TokenAdminRegistryAcceptAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryAcceptAdminRole choice.
func (e *encoder) TokenAdminRegistryAcceptAdminRoleMCMSParams(args TokenAdminRegistryAcceptAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAcceptAdminRole", args)
}

// TokenAdminRegistryAddTokenSend encodes parameters for the TokenAdminRegistryAddTokenSend choice.
func (e *encoder) TokenAdminRegistryAddTokenSend(args TokenAdminRegistryAddTokenSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAddTokenSend", args)
}

// TokenAdminRegistryAddTokenSendMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryAddTokenSend choice.
func (e *encoder) TokenAdminRegistryAddTokenSendMCMSParams(args TokenAdminRegistryAddTokenSendMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAddTokenSend", args)
}

// TokenAdminRegistryAddTokenSendFee encodes parameters for the TokenAdminRegistryAddTokenSendFee choice.
func (e *encoder) TokenAdminRegistryAddTokenSendFee(args TokenAdminRegistryAddTokenSendFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAddTokenSendFee", args)
}

// TokenAdminRegistryAddTokenSendFeeMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryAddTokenSendFee choice.
func (e *encoder) TokenAdminRegistryAddTokenSendFeeMCMSParams(args TokenAdminRegistryAddTokenSendFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryAddTokenSendFee", args)
}

// TokenAdminRegistryConsumeReceiveTicket encodes parameters for the TokenAdminRegistryConsumeReceiveTicket choice.
func (e *encoder) TokenAdminRegistryConsumeReceiveTicket(args TokenAdminRegistryConsumeReceiveTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryConsumeReceiveTicket", args)
}

// TokenAdminRegistryConsumeReceiveTicketMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryConsumeReceiveTicket choice.
func (e *encoder) TokenAdminRegistryConsumeReceiveTicketMCMSParams(args TokenAdminRegistryConsumeReceiveTicketMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryConsumeReceiveTicket", args)
}

// TokenAdminRegistryFinalizeExecute encodes parameters for the TokenAdminRegistryFinalizeExecute choice.
func (e *encoder) TokenAdminRegistryFinalizeExecute(args TokenAdminRegistryFinalizeExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryFinalizeExecute", args)
}

// TokenAdminRegistryFinalizeExecuteMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryFinalizeExecute choice.
func (e *encoder) TokenAdminRegistryFinalizeExecuteMCMSParams(args TokenAdminRegistryFinalizeExecuteMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryFinalizeExecute", args)
}

// TokenAdminRegistryGetTokenConfig encodes parameters for the TokenAdminRegistryGetTokenConfig choice.
func (e *encoder) TokenAdminRegistryGetTokenConfig(args TokenAdminRegistryGetTokenConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryGetTokenConfig", args)
}

// TokenAdminRegistryGetTokenConfigMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryGetTokenConfig choice.
func (e *encoder) TokenAdminRegistryGetTokenConfigMCMSParams(args TokenAdminRegistryGetTokenConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryGetTokenConfig", args)
}

// TokenAdminRegistryIsAdministrator encodes parameters for the TokenAdminRegistryIsAdministrator choice.
func (e *encoder) TokenAdminRegistryIsAdministrator(args TokenAdminRegistryIsAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryIsAdministrator", args)
}

// TokenAdminRegistryIsAdministratorMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryIsAdministrator choice.
func (e *encoder) TokenAdminRegistryIsAdministratorMCMSParams(args TokenAdminRegistryIsAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryIsAdministrator", args)
}

// TokenAdminRegistryProposeAdministrator encodes parameters for the TokenAdminRegistryProposeAdministrator choice.
func (e *encoder) TokenAdminRegistryProposeAdministrator(args TokenAdminRegistryProposeAdministrator) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryProposeAdministrator", args)
}

// TokenAdminRegistryProposeAdministratorMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryProposeAdministrator choice.
func (e *encoder) TokenAdminRegistryProposeAdministratorMCMSParams(args TokenAdminRegistryProposeAdministratorMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryProposeAdministrator", args)
}

// TokenAdminRegistrySetInboundPoolCCVs encodes parameters for the TokenAdminRegistrySetInboundPoolCCVs choice.
func (e *encoder) TokenAdminRegistrySetInboundPoolCCVs(args TokenAdminRegistrySetInboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetInboundPoolCCVs", args)
}

// TokenAdminRegistrySetInboundPoolCCVsMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistrySetInboundPoolCCVs choice.
func (e *encoder) TokenAdminRegistrySetInboundPoolCCVsMCMSParams(args TokenAdminRegistrySetInboundPoolCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetInboundPoolCCVs", args)
}

// TokenAdminRegistrySetOutboundPoolCCVs encodes parameters for the TokenAdminRegistrySetOutboundPoolCCVs choice.
func (e *encoder) TokenAdminRegistrySetOutboundPoolCCVs(args TokenAdminRegistrySetOutboundPoolCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetOutboundPoolCCVs", args)
}

// TokenAdminRegistrySetOutboundPoolCCVsMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistrySetOutboundPoolCCVs choice.
func (e *encoder) TokenAdminRegistrySetOutboundPoolCCVsMCMSParams(args TokenAdminRegistrySetOutboundPoolCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetOutboundPoolCCVs", args)
}

// TokenAdminRegistrySetPool encodes parameters for the TokenAdminRegistrySetPool choice.
func (e *encoder) TokenAdminRegistrySetPool(args TokenAdminRegistrySetPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetPool", args)
}

// TokenAdminRegistrySetPoolMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistrySetPool choice.
func (e *encoder) TokenAdminRegistrySetPoolMCMSParams(args TokenAdminRegistrySetPoolMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistrySetPool", args)
}

// TokenAdminRegistryTransferAdminRole encodes parameters for the TokenAdminRegistryTransferAdminRole choice.
func (e *encoder) TokenAdminRegistryTransferAdminRole(args TokenAdminRegistryTransferAdminRole) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryTransferAdminRole", args)
}

// TokenAdminRegistryTransferAdminRoleMCMSParams encodes MCMS parameters (without Caller) for the TokenAdminRegistryTransferAdminRole choice.
func (e *encoder) TokenAdminRegistryTransferAdminRoleMCMSParams(args TokenAdminRegistryTransferAdminRoleMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TokenAdminRegistryTransferAdminRole", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
