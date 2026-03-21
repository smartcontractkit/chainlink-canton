package offramp

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
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
	PackageName = "ccip-offramp"
	PackageID   = "b6cb8741bd6730bcb1a641e8a7f148d1b7f3a6979378d7e3df7b12419fe4caed"
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

// ExecuteFromRouter is a Record type
type ExecuteFromRouter struct {
	RouterPartyOwner      types.PARTY       `json:"routerPartyOwner"`
	ExecutingMessageCid   types.CONTRACT_ID `json:"executingMessageCid"`
	GlobalConfigCid       types.CONTRACT_ID `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          types.CONTRACT_ID `json:"rmnRemoteCid"`
}

// ToMap converts ExecuteFromRouter to a map for DAML arguments
func (t ExecuteFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	return m
}

func (t ExecuteFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteFromRouter to hex string (Canton MCMS format)
func (t ExecuteFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteFromRouter from hex string (Canton MCMS format)
func (t *ExecuteFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteFromRouterResult is a Record type
type ExecuteFromRouterResult struct {
	MessageId             types.TEXT         `json:"messageId"`
	Message               common.MessageV1   `json:"message"`
	SourceChainSelector   types.NUMERIC      `json:"sourceChainSelector"`
	SequenceNumber        types.NUMERIC      `json:"sequenceNumber"`
	TokenReceiveTicket    *types.CONTRACT_ID `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID  `json:"executionStateChanged"`
}

// ToMap converts ExecuteFromRouterResult to a map for DAML arguments
func (t ExecuteFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageId"] = string(t.MessageId)

	m["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenReceiveTicket,
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
		}
	}

	m["executionStateChanged"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutionStateChanged).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutionStateChanged
	}()

	return m
}

func (t ExecuteFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteFromRouterResult to hex string (Canton MCMS format)
func (t ExecuteFromRouterResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteFromRouterResult from hex string (Canton MCMS format)
func (t *ExecuteFromRouterResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	GlobalConfigCid      types.CONTRACT_ID           `json:"globalConfigCid"`
	ReceiverRequiredCCVs []common.RawInstanceAddress `json:"receiverRequiredCCVs"`
	SourceChainSelector  types.NUMERIC               `json:"sourceChainSelector"`
}

// ToMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
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

func (t GetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForExecute to hex string (Canton MCMS format)
func (t GetRequiredCCVsForExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForExecute from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// OffRamp is a Template type
type OffRamp struct {
	InstanceId types.TEXT  `json:"instanceId"`
	CcipOwner  types.PARTY `json:"ccipOwner"`
	Deps       OffRampDeps `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OffRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OffRamp) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OffRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t OffRamp) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t OffRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OffRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OffRamp to hex string (Canton MCMS format)
func (t OffRamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OffRamp from hex string (Canton MCMS format)
func (t *OffRamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for OffRamp

// ExecuteFromRouter exercises the ExecuteFromRouter choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) ExecuteFromRouter(contractID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// ExecuteFromRouterWithPackageID exercises the ExecuteFromRouter choice using the provided package ID instead of package name
func (t OffRamp) ExecuteFromRouterWithPackageID(contractID string, packageID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) PrepareExecute(contractID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t OffRamp) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// SetDeps exercises the SetDeps choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t OffRamp) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t OffRamp) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OffRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OffRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OffRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this OffRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OffRamp) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t OffRamp) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for OffRamp

var _ mcms.IMCMSReceiver = (*OffRamp)(nil)

// OffRampDeps is a Record type
type OffRampDeps struct {
	GlobalConfig       common.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          common.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry common.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts OffRampDeps to a map for DAML arguments
func (t OffRampDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["tokenAdminRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistry
	}()

	return m
}

func (t OffRampDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OffRampDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OffRampDeps to hex string (Canton MCMS format)
func (t OffRampDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OffRampDeps from hex string (Canton MCMS format)
func (t *OffRampDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareExecute is a Record type
type PrepareExecute struct {
	EncodedMessage                types.TEXT                  `json:"encodedMessage"`
	GlobalConfigCid               types.CONTRACT_ID           `json:"globalConfigCid"`
	TokenAdminRegistryCid         types.CONTRACT_ID           `json:"tokenAdminRegistryCid"`
	ReceiverRequiredCCVs          []common.RawInstanceAddress `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs          []common.RawInstanceAddress `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold     types.INT64                 `json:"receiverOptionalThreshold"`
	ReceiverMinBlockConfirmations types.INT64                 `json:"receiverMinBlockConfirmations"`
	RmnRemoteCid                  types.CONTRACT_ID           `json:"rmnRemoteCid"`
	ReceiverParty                 types.PARTY                 `json:"receiverParty"`
	TokenReceiverParty            *types.PARTY                `json:"tokenReceiverParty" hex:"optional"`
	Caller                        types.PARTY                 `json:"caller"`
}

// ToMap converts PrepareExecute to a map for DAML arguments
func (t PrepareExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["receiverOptionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverOptionalCCVs))
		for _, e := range t.ReceiverOptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["receiverOptionalThreshold"] = int64(t.ReceiverOptionalThreshold)

	m["receiverMinBlockConfirmations"] = int64(t.ReceiverMinBlockConfirmations)

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["receiverParty"] = t.ReceiverParty.ToMap()

	if t.TokenReceiverParty != nil {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
			"value": (*t.TokenReceiverParty).ToMap(),
		}
	} else {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PrepareExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PrepareExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PrepareExecute to hex string (Canton MCMS format)
func (t PrepareExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecute from hex string (Canton MCMS format)
func (t *PrepareExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareExecuteMCMSParams is PrepareExecute without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type PrepareExecuteMCMSParams struct {
	EncodedMessage                types.TEXT                  `json:"encodedMessage"`
	GlobalConfigCid               types.CONTRACT_ID           `json:"globalConfigCid"`
	TokenAdminRegistryCid         types.CONTRACT_ID           `json:"tokenAdminRegistryCid"`
	ReceiverRequiredCCVs          []common.RawInstanceAddress `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs          []common.RawInstanceAddress `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold     types.INT64                 `json:"receiverOptionalThreshold"`
	ReceiverMinBlockConfirmations types.INT64                 `json:"receiverMinBlockConfirmations"`
	RmnRemoteCid                  types.CONTRACT_ID           `json:"rmnRemoteCid"`
	ReceiverParty                 types.PARTY                 `json:"receiverParty"`
	TokenReceiverParty            *types.PARTY                `json:"tokenReceiverParty" hex:"optional"`
}

// MarshalHex encodes PrepareExecuteMCMSParams to hex string for MCMS operationData.
func (t PrepareExecuteMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecuteMCMSParams from hex string.
func (t *PrepareExecuteMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDeps is a Record type
type SetDeps struct {
	NewDeps SetDepsParams `json:"newDeps"`
}

// ToMap converts SetDeps to a map for DAML arguments
func (t SetDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewDeps).(mapper); ok {
			return m.toMap()
		}
		return t.NewDeps
	}()

	return m
}

func (t SetDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDeps to hex string (Canton MCMS format)
func (t SetDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDeps from hex string (Canton MCMS format)
func (t *SetDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDepsParams is a Record type
type SetDepsParams struct {
	GlobalConfig       *common.RawInstanceAddress `json:"globalConfig" hex:"optional"`
	RmnRemote          *common.RawInstanceAddress `json:"rmnRemote" hex:"optional"`
	TokenAdminRegistry *common.RawInstanceAddress `json:"tokenAdminRegistry" hex:"optional"`
}

// ToMap converts SetDepsParams to a map for DAML arguments
func (t SetDepsParams) ToMap() map[string]any {
	m := make(map[string]any)

	if t.GlobalConfig != nil {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
			"value": *t.GlobalConfig,
		}
	} else {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.RmnRemote != nil {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": *t.RmnRemote,
		}
	} else {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.TokenAdminRegistry != nil {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenAdminRegistry,
		}
	} else {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
		}
	}

	return m
}

func (t SetDepsParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDepsParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDepsParams to hex string (Canton MCMS format)
func (t SetDepsParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDepsParams from hex string (Canton MCMS format)
func (t *SetDepsParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ExecuteFromRouter(args ExecuteFromRouter) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error)
	PrepareExecute(args PrepareExecute) (*bind.EncodedChoice, error)
	PrepareExecuteMCMSParams(args PrepareExecuteMCMSParams) (*bind.EncodedChoice, error)
	SetDeps(args SetDeps) (*bind.EncodedChoice, error)
	SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error)
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

// ExecuteFromRouter encodes parameters for the ExecuteFromRouter choice.
func (e *encoder) ExecuteFromRouter(args ExecuteFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteFromRouter", args)
}

// GetRequiredCCVsForExecute encodes parameters for the GetRequiredCCVsForExecute choice.
func (e *encoder) GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForExecute", args)
}

// PrepareExecute encodes parameters for the PrepareExecute choice.
func (e *encoder) PrepareExecute(args PrepareExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// PrepareExecuteMCMSParams encodes MCMS parameters (without Caller) for the PrepareExecute choice.
func (e *encoder) PrepareExecuteMCMSParams(args PrepareExecuteMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// SetDeps encodes parameters for the SetDeps choice.
func (e *encoder) SetDeps(args SetDeps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// SetDepsParams encodes parameters for the SetDeps choice.
func (e *encoder) SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
