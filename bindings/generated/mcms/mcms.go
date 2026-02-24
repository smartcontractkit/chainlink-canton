package mcms

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	PackageName = "mcms"
	PackageID   = "c55a819cf2f7e1fc2dfb41eb33910e406a01af9e0a043b0720849bbf8b3be11f"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IMCMSReceiver is a DAML interface
type IMCMSReceiver interface {

	// MCMSReceiverEntrypoint executes the MCMSReceiver_Entrypoint choice
	MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand

	// MCMSReceiverGetInstanceId executes the MCMSReceiver_GetInstanceId choice
	MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand

	// MCMSReceiverGetView executes the MCMSReceiver_GetView choice
	MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand
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

// APSetConfig is a Record type
type APSetConfig struct {
	ApSigners      []SignerInfo  `json:"apSigners"`
	ApGroupQuorums []types.INT64 `json:"apGroupQuorums"`
	ApGroupParents []types.INT64 `json:"apGroupParents"`
	ApClearRoot    types.BOOL    `json:"apClearRoot"`
}

// ToMap converts APSetConfig to a map for DAML arguments
func (t APSetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["apSigners"] = func() []any {
		res := make([]any, 0, len(t.ApSigners))
		for _, e := range t.ApSigners {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["apGroupQuorums"] = func() []any {
		res := make([]any, 0, len(t.ApGroupQuorums))
		for _, e := range t.ApGroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apGroupParents"] = func() []any {
		res := make([]any, 0, len(t.ApGroupParents))
		for _, e := range t.ApGroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apClearRoot"] = bool(t.ApClearRoot)

	return m
}

func (t APSetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *APSetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes APSetConfig to hex string (Canton MCMS format)
func (t APSetConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes APSetConfig from hex string (Canton MCMS format)
func (t *APSetConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AdminParams is a variant/union type
type AdminParams struct {
	APSetConfig *APSetConfig `json:"AP_SetConfig,omitempty"`
	APClearRoot *types.UNIT  `json:"AP_ClearRoot,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AdminParams
func (v AdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for AdminParams
func (v *AdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes AdminParams to hex string (Canton MCMS format)
func (v AdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes AdminParams from hex string (Canton MCMS format)
func (v *AdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v AdminParams) GetVariantTag() string {

	if v.APSetConfig != nil {
		return "AP_SetConfig"
	}

	if v.APClearRoot != nil {
		return "AP_ClearRoot"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v AdminParams) GetVariantValue() any {

	if v.APSetConfig != nil {
		return v.APSetConfig
	}

	if v.APClearRoot != nil {
		return v.APClearRoot
	}

	return nil
}

var _ types.VARIANT = (*AdminParams)(nil)

// ArchiveCounterV1 is a Record type
type ArchiveCounterV1 struct {
}

// ToMap converts ArchiveCounterV1 to a map for DAML arguments
func (t ArchiveCounterV1) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ArchiveCounterV1) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ArchiveCounterV1) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ArchiveCounterV1 to hex string (Canton MCMS format)
func (t ArchiveCounterV1) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ArchiveCounterV1 from hex string (Canton MCMS format)
func (t *ArchiveCounterV1) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ArchiveMCMSEntrypointEvent is a Record type
type ArchiveMCMSEntrypointEvent struct {
}

// ToMap converts ArchiveMCMSEntrypointEvent to a map for DAML arguments
func (t ArchiveMCMSEntrypointEvent) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ArchiveMCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ArchiveMCMSEntrypointEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ArchiveMCMSEntrypointEvent to hex string (Canton MCMS format)
func (t ArchiveMCMSEntrypointEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ArchiveMCMSEntrypointEvent from hex string (Canton MCMS format)
func (t *ArchiveMCMSEntrypointEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ArgValue is a variant/union type
type ArgValue struct {
	AVText  *types.TEXT      `json:"AV_Text,omitempty"`
	AVInt   *types.INT64     `json:"AV_Int,omitempty"`
	AVBool  *types.BOOL      `json:"AV_Bool,omitempty"`
	AVParty *types.PARTY     `json:"AV_Party,omitempty"`
	AVTime  *types.TIMESTAMP `json:"AV_Time,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ArgValue
func (v ArgValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for ArgValue
func (v *ArgValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes ArgValue to hex string (Canton MCMS format)
func (v ArgValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes ArgValue from hex string (Canton MCMS format)
func (v *ArgValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ArgValue) GetVariantTag() string {

	if v.AVText != nil {
		return "AV_Text"
	}

	if v.AVInt != nil {
		return "AV_Int"
	}

	if v.AVBool != nil {
		return "AV_Bool"
	}

	if v.AVParty != nil {
		return "AV_Party"
	}

	if v.AVTime != nil {
		return "AV_Time"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ArgValue) GetVariantValue() any {

	if v.AVText != nil {
		return v.AVText
	}

	if v.AVInt != nil {
		return v.AVInt
	}

	if v.AVBool != nil {
		return v.AVBool
	}

	if v.AVParty != nil {
		return v.AVParty
	}

	if v.AVTime != nil {
		return v.AVTime
	}

	return nil
}

var _ types.VARIANT = (*ArgValue)(nil)

// BlockedFunction is a Record type
type BlockedFunction struct {
	TargetInstanceId types.TEXT `json:"targetInstanceId"`
	FunctionName     types.TEXT `json:"functionName"`
}

// ToMap converts BlockedFunction to a map for DAML arguments
func (t BlockedFunction) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceId"] = string(t.TargetInstanceId)

	m["functionName"] = string(t.FunctionName)

	return m
}

func (t BlockedFunction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BlockedFunction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BlockedFunction to hex string (Canton MCMS format)
func (t BlockedFunction) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BlockedFunction from hex string (Canton MCMS format)
func (t *BlockedFunction) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BypasserExecuteBatchParams is a Record type
type BypasserExecuteBatchParams struct {
	Calls []TimelockCall `json:"calls"`
}

// ToMap converts BypasserExecuteBatchParams to a map for DAML arguments
func (t BypasserExecuteBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			type mapper interface{ toMap() map[string]any }
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

func (t BypasserExecuteBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BypasserExecuteBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BypasserExecuteBatchParams to hex string (Canton MCMS format)
func (t BypasserExecuteBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BypasserExecuteBatchParams from hex string (Canton MCMS format)
func (t *BypasserExecuteBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CanExecuteOp is a Record type
type CanExecuteOp struct {
	Submitter  types.PARTY `json:"submitter"`
	TargetRole Role        `json:"targetRole"`
	Op         Op          `json:"op"`
}

// ToMap converts CanExecuteOp to a map for DAML arguments
func (t CanExecuteOp) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["targetRole"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TargetRole).(mapper); ok {
			return m.toMap()
		}
		return t.TargetRole
	}()

	m["op"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Op).(mapper); ok {
			return m.toMap()
		}
		return t.Op
	}()

	return m
}

func (t CanExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CanExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CanExecuteOp to hex string (Canton MCMS format)
func (t CanExecuteOp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CanExecuteOp from hex string (Canton MCMS format)
func (t *CanExecuteOp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelBatchParams is a Record type
type CancelBatchParams struct {
	OpId types.TEXT `json:"opId"`
}

// ToMap converts CancelBatchParams to a map for DAML arguments
func (t CancelBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["opId"] = string(t.OpId)

	return m
}

func (t CancelBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelBatchParams to hex string (Canton MCMS format)
func (t CancelBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelBatchParams from hex string (Canton MCMS format)
func (t *CancelBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelUpgradeParams is a Record type
type CancelUpgradeParams struct {
	InstanceId types.TEXT `json:"instanceId"`
}

// ToMap converts CancelUpgradeParams to a map for DAML arguments
func (t CancelUpgradeParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t CancelUpgradeParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelUpgradeParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelUpgradeParams to hex string (Canton MCMS format)
func (t CancelUpgradeParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelUpgradeParams from hex string (Canton MCMS format)
func (t *CancelUpgradeParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ClearRegistryCid is a Record type
type ClearRegistryCid struct {
}

// ToMap converts ClearRegistryCid to a map for DAML arguments
func (t ClearRegistryCid) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ClearRegistryCid) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ClearRegistryCid) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ClearRegistryCid to hex string (Canton MCMS format)
func (t ClearRegistryCid) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ClearRegistryCid from hex string (Canton MCMS format)
func (t *ClearRegistryCid) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CompleteUpgradeParams is a Record type
type CompleteUpgradeParams struct {
	InstanceId types.TEXT `json:"instanceId"`
	NewVersion types.TEXT `json:"newVersion"`
}

// ToMap converts CompleteUpgradeParams to a map for DAML arguments
func (t CompleteUpgradeParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["newVersion"] = string(t.NewVersion)

	return m
}

func (t CompleteUpgradeParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CompleteUpgradeParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CompleteUpgradeParams to hex string (Canton MCMS format)
func (t CompleteUpgradeParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CompleteUpgradeParams from hex string (Canton MCMS format)
func (t *CompleteUpgradeParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeStateMigrationTicket is a Record type
type ConsumeStateMigrationTicket struct {
}

// ToMap converts ConsumeStateMigrationTicket to a map for DAML arguments
func (t ConsumeStateMigrationTicket) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ConsumeStateMigrationTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeStateMigrationTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeStateMigrationTicket to hex string (Canton MCMS format)
func (t ConsumeStateMigrationTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeStateMigrationTicket from hex string (Canton MCMS format)
func (t *ConsumeStateMigrationTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ContractRegistration is a Record type
type ContractRegistration struct {
	InstanceId   types.TEXT      `json:"instanceId"`
	Version      types.TEXT      `json:"version"`
	Status       ContractStatus  `json:"status"`
	DeployedAt   types.TIMESTAMP `json:"deployedAt"`
	ContractType types.TEXT      `json:"contractType"`
	Metadata     types.GENMAP    `json:"metadata"`
}

// ToMap converts ContractRegistration to a map for DAML arguments
func (t ContractRegistration) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["version"] = string(t.Version)

	m["status"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Status).(mapper); ok {
			return m.toMap()
		}
		return t.Status
	}()

	m["deployedAt"] = t.DeployedAt

	m["contractType"] = string(t.ContractType)

	m["metadata"] = func() any {
		if t.Metadata == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.Metadata}
	}()

	return m
}

func (t ContractRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ContractRegistration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ContractRegistration to hex string (Canton MCMS format)
func (t ContractRegistration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ContractRegistration from hex string (Canton MCMS format)
func (t *ContractRegistration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ContractStatus is an enum type
type ContractStatus string

const (
	ContractStatusActive ContractStatus = "Active"

	ContractStatusUpgradePending ContractStatus = "UpgradePending"

	ContractStatusDeprecated ContractStatus = "Deprecated"
)

func (e ContractStatus) GetEnumConstructor() string { return string(e) }

func (e ContractStatus) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "ContractStatus")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e ContractStatus) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "ContractStatus")
}

func (e ContractStatus) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *ContractStatus) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes ContractStatus to hex string (Canton MCMS format)
func (e ContractStatus) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes ContractStatus from hex string (Canton MCMS format)
func (e *ContractStatus) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = ContractStatus("")

// Counter is a Template type
type Counter struct {
	Owner      types.PARTY `json:"owner"`
	InstanceId types.TEXT  `json:"instanceId"`
	Value      types.INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t Counter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t Counter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.Counter", "Counter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t Counter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t Counter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t Counter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Counter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Counter to hex string (Canton MCMS format)
func (t Counter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Counter from hex string (Canton MCMS format)
func (t *Counter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for Counter

// Archive exercises the Archive choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t Counter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetValue exercises the GetValue choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetValue(contractID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetValue",
		Arguments:  argsToMap(args),
	}
}

// GetValueWithPackageID exercises the GetValue choice using the provided package ID instead of package name
func (t Counter) GetValueWithPackageID(contractID string, packageID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetValue",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceIdChoice exercises the GetInstanceIdChoice choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetInstanceIdChoice(contractID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceIdChoiceWithPackageID exercises the GetInstanceIdChoice choice using the provided package ID instead of package name
func (t Counter) GetInstanceIdChoiceWithPackageID(contractID string, packageID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this Counter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Counter) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t Counter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this Counter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Counter) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceIdWithPackageID exercises the MCMSReceiver_GetInstanceId choice using the provided package ID instead of package name
func (t Counter) MCMSReceiverGetInstanceIdWithPackageID(contractID string, packageID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetView exercises the MCMSReceiver_GetView choice on this Counter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Counter) MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetViewWithPackageID exercises the MCMSReceiver_GetView choice using the provided package ID instead of package name
func (t Counter) MCMSReceiverGetViewWithPackageID(contractID string, packageID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Counter

var _ IMCMSReceiver = (*Counter)(nil)

// CounterUpgradeService is a Template type
type CounterUpgradeService struct {
	Owner             types.PARTY `json:"owner"`
	ServiceInstanceId types.TEXT  `json:"serviceInstanceId"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CounterUpgradeService) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterUpgradeService")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CounterUpgradeService) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.CounterV2", "CounterUpgradeService")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CounterUpgradeService) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["serviceInstanceId"] = string(t.ServiceInstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CounterUpgradeService) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["serviceInstanceId"] = string(t.ServiceInstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CounterUpgradeService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CounterUpgradeService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CounterUpgradeService to hex string (Canton MCMS format)
func (t CounterUpgradeService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CounterUpgradeService from hex string (Canton MCMS format)
func (t *CounterUpgradeService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CounterUpgradeService

// Archive exercises the Archive choice on this CounterUpgradeService contract
// This method uses the package name in the template ID
func (t CounterUpgradeService) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterUpgradeService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CounterUpgradeService) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterUpgradeService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this CounterUpgradeService contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterUpgradeService) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CounterUpgradeService) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this CounterUpgradeService contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterUpgradeService) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceIdWithPackageID exercises the MCMSReceiver_GetInstanceId choice using the provided package ID instead of package name
func (t CounterUpgradeService) MCMSReceiverGetInstanceIdWithPackageID(contractID string, packageID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetView exercises the MCMSReceiver_GetView choice on this CounterUpgradeService contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterUpgradeService) MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetViewWithPackageID exercises the MCMSReceiver_GetView choice using the provided package ID instead of package name
func (t CounterUpgradeService) MCMSReceiverGetViewWithPackageID(contractID string, packageID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CounterUpgradeService

var _ IMCMSReceiver = (*CounterUpgradeService)(nil)

// CounterV1Upgradeable is a Template type
type CounterV1Upgradeable struct {
	Owner      types.PARTY `json:"owner"`
	InstanceId types.TEXT  `json:"instanceId"`
	Value      types.INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CounterV1Upgradeable) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV1Upgradeable")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CounterV1Upgradeable) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV1Upgradeable")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CounterV1Upgradeable) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CounterV1Upgradeable) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CounterV1Upgradeable) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CounterV1Upgradeable) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CounterV1Upgradeable to hex string (Canton MCMS format)
func (t CounterV1Upgradeable) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CounterV1Upgradeable from hex string (Canton MCMS format)
func (t *CounterV1Upgradeable) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CounterV1Upgradeable

// Archive exercises the Archive choice on this CounterV1Upgradeable contract
// This method uses the package name in the template ID
func (t CounterV1Upgradeable) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV1Upgradeable"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CounterV1Upgradeable) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV1Upgradeable"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveCounterV1 exercises the ArchiveCounterV1 choice on this CounterV1Upgradeable contract
// This method uses the package name in the template ID
func (t CounterV1Upgradeable) ArchiveCounterV1(contractID string, args ArchiveCounterV1) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV1Upgradeable"),
		ContractID: contractID,
		Choice:     "ArchiveCounterV1",
		Arguments:  argsToMap(args),
	}
}

// ArchiveCounterV1WithPackageID exercises the ArchiveCounterV1 choice using the provided package ID instead of package name
func (t CounterV1Upgradeable) ArchiveCounterV1WithPackageID(contractID string, packageID string, args ArchiveCounterV1) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV1Upgradeable"),
		ContractID: contractID,
		Choice:     "ArchiveCounterV1",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this CounterV1Upgradeable contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV1Upgradeable) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CounterV1Upgradeable) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this CounterV1Upgradeable contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV1Upgradeable) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceIdWithPackageID exercises the MCMSReceiver_GetInstanceId choice using the provided package ID instead of package name
func (t CounterV1Upgradeable) MCMSReceiverGetInstanceIdWithPackageID(contractID string, packageID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetView exercises the MCMSReceiver_GetView choice on this CounterV1Upgradeable contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV1Upgradeable) MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetViewWithPackageID exercises the MCMSReceiver_GetView choice using the provided package ID instead of package name
func (t CounterV1Upgradeable) MCMSReceiverGetViewWithPackageID(contractID string, packageID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CounterV1Upgradeable

var _ IMCMSReceiver = (*CounterV1Upgradeable)(nil)

// CounterV2 is a Template type
type CounterV2 struct {
	Owner        types.PARTY     `json:"owner"`
	InstanceId   types.TEXT      `json:"instanceId"`
	Value        types.INT64     `json:"value"`
	Label        types.TEXT      `json:"label"`
	LastModified types.TIMESTAMP `json:"lastModified"`
	MigratedFrom *types.TEXT     `json:"migratedFrom"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CounterV2) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CounterV2) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CounterV2) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["label"] = string(t.Label)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastModified"] = t.LastModified

	if t.MigratedFrom != nil {
		args["migratedFrom"] = map[string]any{
			"_type": "optional",
			"value": string(*t.MigratedFrom),
		}
	} else {
		args["migratedFrom"] = map[string]any{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CounterV2) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["value"] = int64(t.Value)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["label"] = string(t.Label)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastModified"] = t.LastModified

	if t.MigratedFrom != nil {
		args["migratedFrom"] = map[string]any{
			"_type": "optional",
			"value": string(*t.MigratedFrom),
		}
	} else {
		args["migratedFrom"] = map[string]any{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CounterV2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CounterV2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CounterV2 to hex string (Canton MCMS format)
func (t CounterV2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CounterV2 from hex string (Canton MCMS format)
func (t *CounterV2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CounterV2

// Archive exercises the Archive choice on this CounterV2 contract
// This method uses the package name in the template ID
func (t CounterV2) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CounterV2) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetValueV2 exercises the GetValueV2 choice on this CounterV2 contract
// This method uses the package name in the template ID
func (t CounterV2) GetValueV2(contractID string, args GetValueV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetValueV2",
		Arguments:  argsToMap(args),
	}
}

// GetValueV2WithPackageID exercises the GetValueV2 choice using the provided package ID instead of package name
func (t CounterV2) GetValueV2WithPackageID(contractID string, packageID string, args GetValueV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetValueV2",
		Arguments:  argsToMap(args),
	}
}

// GetLabelV2 exercises the GetLabelV2 choice on this CounterV2 contract
// This method uses the package name in the template ID
func (t CounterV2) GetLabelV2(contractID string, args GetLabelV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetLabelV2",
		Arguments:  argsToMap(args),
	}
}

// GetLabelV2WithPackageID exercises the GetLabelV2 choice using the provided package ID instead of package name
func (t CounterV2) GetLabelV2WithPackageID(contractID string, packageID string, args GetLabelV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetLabelV2",
		Arguments:  argsToMap(args),
	}
}

// GetLastModifiedV2 exercises the GetLastModifiedV2 choice on this CounterV2 contract
// This method uses the package name in the template ID
func (t CounterV2) GetLastModifiedV2(contractID string, args GetLastModifiedV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetLastModifiedV2",
		Arguments:  argsToMap(args),
	}
}

// GetLastModifiedV2WithPackageID exercises the GetLastModifiedV2 choice using the provided package ID instead of package name
func (t CounterV2) GetLastModifiedV2WithPackageID(contractID string, packageID string, args GetLastModifiedV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetLastModifiedV2",
		Arguments:  argsToMap(args),
	}
}

// GetFullStateV2 exercises the GetFullStateV2 choice on this CounterV2 contract
// This method uses the package name in the template ID
func (t CounterV2) GetFullStateV2(contractID string, args GetFullStateV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetFullStateV2",
		Arguments:  argsToMap(args),
	}
}

// GetFullStateV2WithPackageID exercises the GetFullStateV2 choice using the provided package ID instead of package name
func (t CounterV2) GetFullStateV2WithPackageID(contractID string, packageID string, args GetFullStateV2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "CounterV2"),
		ContractID: contractID,
		Choice:     "GetFullStateV2",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this CounterV2 contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV2) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CounterV2) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this CounterV2 contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV2) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceIdWithPackageID exercises the MCMSReceiver_GetInstanceId choice using the provided package ID instead of package name
func (t CounterV2) MCMSReceiverGetInstanceIdWithPackageID(contractID string, packageID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetView exercises the MCMSReceiver_GetView choice on this CounterV2 contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CounterV2) MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetViewWithPackageID exercises the MCMSReceiver_GetView choice using the provided package ID instead of package name
func (t CounterV2) MCMSReceiverGetViewWithPackageID(contractID string, packageID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.CounterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CounterV2

var _ IMCMSReceiver = (*CounterV2)(nil)

// CounterV2State is a Record type
type CounterV2State struct {
	Value        types.INT64     `json:"value"`
	Label        types.TEXT      `json:"label"`
	LastModified types.TIMESTAMP `json:"lastModified"`
	MigratedFrom *types.TEXT     `json:"migratedFrom"`
}

// ToMap converts CounterV2State to a map for DAML arguments
func (t CounterV2State) ToMap() map[string]any {
	m := make(map[string]any)

	m["value"] = int64(t.Value)

	m["label"] = string(t.Label)

	m["lastModified"] = t.LastModified

	if t.MigratedFrom != nil {
		m["migratedFrom"] = map[string]any{
			"_type": "optional",
			"value": string(*t.MigratedFrom),
		}
	} else {
		m["migratedFrom"] = map[string]any{
			"_type": "optional",
		}
	}

	return m
}

func (t CounterV2State) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CounterV2State) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CounterV2State to hex string (Canton MCMS format)
func (t CounterV2State) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CounterV2State from hex string (Canton MCMS format)
func (t *CounterV2State) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteOp is a Record type
type ExecuteOp struct {
	TargetRole Role                `json:"targetRole"`
	Submitter  types.PARTY         `json:"submitter"`
	Op         Op                  `json:"op"`
	OpProof    []types.TEXT        `json:"opProof"`
	TargetCids []types.CONTRACT_ID `json:"targetCids"`
}

// ToMap converts ExecuteOp to a map for DAML arguments
func (t ExecuteOp) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TargetRole).(mapper); ok {
			return m.toMap()
		}
		return t.TargetRole
	}()

	m["submitter"] = t.Submitter.ToMap()

	m["op"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Op).(mapper); ok {
			return m.toMap()
		}
		return t.Op
	}()

	m["opProof"] = func() []any {
		res := make([]any, 0, len(t.OpProof))
		for _, e := range t.OpProof {
			res = append(res, string(e))
		}
		return res
	}()

	m["targetCids"] = func() []any {
		res := make([]any, 0, len(t.TargetCids))
		for _, e := range t.TargetCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t ExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteOp to hex string (Canton MCMS format)
func (t ExecuteOp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteOp from hex string (Canton MCMS format)
func (t *ExecuteOp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteScheduledBatch is a Record type
type ExecuteScheduledBatch struct {
	Submitter   types.PARTY         `json:"submitter"`
	OpId        types.TEXT          `json:"opId"`
	Calls       []TimelockCall      `json:"calls"`
	Predecessor types.TEXT          `json:"predecessor"`
	Salt        types.TEXT          `json:"salt"`
	TargetCids  []types.CONTRACT_ID `json:"targetCids"`
}

// ToMap converts ExecuteScheduledBatch to a map for DAML arguments
func (t ExecuteScheduledBatch) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["predecessor"] = string(t.Predecessor)

	m["salt"] = string(t.Salt)

	m["targetCids"] = func() []any {
		res := make([]any, 0, len(t.TargetCids))
		for _, e := range t.TargetCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t ExecuteScheduledBatch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteScheduledBatch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteScheduledBatch to hex string (Canton MCMS format)
func (t ExecuteScheduledBatch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteScheduledBatch from hex string (Canton MCMS format)
func (t *ExecuteScheduledBatch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExpiringRoot is a Record type
type ExpiringRoot struct {
	Root       types.TEXT      `json:"root"`
	ValidUntil types.TIMESTAMP `json:"validUntil"`
	OpCount    types.INT64     `json:"opCount"`
}

// ToMap converts ExpiringRoot to a map for DAML arguments
func (t ExpiringRoot) ToMap() map[string]any {
	m := make(map[string]any)

	m["root"] = string(t.Root)

	m["validUntil"] = t.ValidUntil

	m["opCount"] = int64(t.OpCount)

	return m
}

func (t ExpiringRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExpiringRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExpiringRoot to hex string (Canton MCMS format)
func (t ExpiringRoot) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExpiringRoot from hex string (Canton MCMS format)
func (t *ExpiringRoot) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetBlockedFunctions is a Record type
type GetBlockedFunctions struct {
	Submitter types.PARTY `json:"submitter"`
}

// ToMap converts GetBlockedFunctions to a map for DAML arguments
func (t GetBlockedFunctions) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	return m
}

func (t GetBlockedFunctions) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetBlockedFunctions) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetBlockedFunctions to hex string (Canton MCMS format)
func (t GetBlockedFunctions) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetBlockedFunctions from hex string (Canton MCMS format)
func (t *GetBlockedFunctions) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetBlockedFunctionsCount is a Record type
type GetBlockedFunctionsCount struct {
	Submitter types.PARTY `json:"submitter"`
}

// ToMap converts GetBlockedFunctionsCount to a map for DAML arguments
func (t GetBlockedFunctionsCount) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	return m
}

func (t GetBlockedFunctionsCount) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetBlockedFunctionsCount) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetBlockedFunctionsCount to hex string (Canton MCMS format)
func (t GetBlockedFunctionsCount) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetBlockedFunctionsCount from hex string (Canton MCMS format)
func (t *GetBlockedFunctionsCount) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetContractInfo is a Record type
type GetContractInfo struct {
	Viewer           types.PARTY `json:"viewer"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
}

// ToMap converts GetContractInfo to a map for DAML arguments
func (t GetContractInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	m["targetInstanceId"] = string(t.TargetInstanceId)

	return m
}

func (t GetContractInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetContractInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetContractInfo to hex string (Canton MCMS format)
func (t GetContractInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetContractInfo from hex string (Canton MCMS format)
func (t *GetContractInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFullStateV2 is a Record type
type GetFullStateV2 struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetFullStateV2 to a map for DAML arguments
func (t GetFullStateV2) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetFullStateV2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFullStateV2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFullStateV2 to hex string (Canton MCMS format)
func (t GetFullStateV2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFullStateV2 from hex string (Canton MCMS format)
func (t *GetFullStateV2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetInstanceIdChoice is a Record type
type GetInstanceIdChoice struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetInstanceIdChoice to a map for DAML arguments
func (t GetInstanceIdChoice) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetInstanceIdChoice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetInstanceIdChoice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetInstanceIdChoice to hex string (Canton MCMS format)
func (t GetInstanceIdChoice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetInstanceIdChoice from hex string (Canton MCMS format)
func (t *GetInstanceIdChoice) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetLabelV2 is a Record type
type GetLabelV2 struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetLabelV2 to a map for DAML arguments
func (t GetLabelV2) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetLabelV2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetLabelV2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetLabelV2 to hex string (Canton MCMS format)
func (t GetLabelV2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetLabelV2 from hex string (Canton MCMS format)
func (t *GetLabelV2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetLastModifiedV2 is a Record type
type GetLastModifiedV2 struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetLastModifiedV2 to a map for DAML arguments
func (t GetLastModifiedV2) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetLastModifiedV2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetLastModifiedV2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetLastModifiedV2 to hex string (Canton MCMS format)
func (t GetLastModifiedV2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetLastModifiedV2 from hex string (Canton MCMS format)
func (t *GetLastModifiedV2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMinDelay is a Record type
type GetMinDelay struct {
	Submitter types.PARTY `json:"submitter"`
}

// ToMap converts GetMinDelay to a map for DAML arguments
func (t GetMinDelay) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	return m
}

func (t GetMinDelay) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetMinDelay) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetMinDelay to hex string (Canton MCMS format)
func (t GetMinDelay) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMinDelay from hex string (Canton MCMS format)
func (t *GetMinDelay) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetPendingUpgrade is a Record type
type GetPendingUpgrade struct {
	Viewer           types.PARTY `json:"viewer"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
}

// ToMap converts GetPendingUpgrade to a map for DAML arguments
func (t GetPendingUpgrade) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	m["targetInstanceId"] = string(t.TargetInstanceId)

	return m
}

func (t GetPendingUpgrade) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetPendingUpgrade) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetPendingUpgrade to hex string (Canton MCMS format)
func (t GetPendingUpgrade) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetPendingUpgrade from hex string (Canton MCMS format)
func (t *GetPendingUpgrade) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRegistryCid is a Record type
type GetRegistryCid struct {
	Submitter types.PARTY `json:"submitter"`
}

// ToMap converts GetRegistryCid to a map for DAML arguments
func (t GetRegistryCid) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	return m
}

func (t GetRegistryCid) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRegistryCid) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRegistryCid to hex string (Canton MCMS format)
func (t GetRegistryCid) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRegistryCid from hex string (Canton MCMS format)
func (t *GetRegistryCid) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetState is a Record type
type GetState struct {
	Submitter  types.PARTY `json:"submitter"`
	TargetRole Role        `json:"targetRole"`
}

// ToMap converts GetState to a map for DAML arguments
func (t GetState) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["targetRole"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TargetRole).(mapper); ok {
			return m.toMap()
		}
		return t.TargetRole
	}()

	return m
}

func (t GetState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetState to hex string (Canton MCMS format)
func (t GetState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetState from hex string (Canton MCMS format)
func (t *GetState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTicketInfo is a Record type
type GetTicketInfo struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetTicketInfo to a map for DAML arguments
func (t GetTicketInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetTicketInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTicketInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTicketInfo to hex string (Canton MCMS format)
func (t GetTicketInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTicketInfo from hex string (Canton MCMS format)
func (t *GetTicketInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetTimestamp is a Record type
type GetTimestamp struct {
	Submitter types.PARTY `json:"submitter"`
	OpId      types.TEXT  `json:"opId"`
}

// ToMap converts GetTimestamp to a map for DAML arguments
func (t GetTimestamp) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	return m
}

func (t GetTimestamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetTimestamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetTimestamp to hex string (Canton MCMS format)
func (t GetTimestamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetTimestamp from hex string (Canton MCMS format)
func (t *GetTimestamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetUpgradeHistory is a Record type
type GetUpgradeHistory struct {
	Viewer           types.PARTY `json:"viewer"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
}

// ToMap converts GetUpgradeHistory to a map for DAML arguments
func (t GetUpgradeHistory) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	m["targetInstanceId"] = string(t.TargetInstanceId)

	return m
}

func (t GetUpgradeHistory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetUpgradeHistory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetUpgradeHistory to hex string (Canton MCMS format)
func (t GetUpgradeHistory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetUpgradeHistory from hex string (Canton MCMS format)
func (t *GetUpgradeHistory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetValue is a Record type
type GetValue struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetValue to a map for DAML arguments
func (t GetValue) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetValue to hex string (Canton MCMS format)
func (t GetValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetValue from hex string (Canton MCMS format)
func (t *GetValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetValueV2 is a Record type
type GetValueV2 struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetValueV2 to a map for DAML arguments
func (t GetValueV2) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetValueV2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetValueV2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetValueV2 to hex string (Canton MCMS format)
func (t GetValueV2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetValueV2 from hex string (Canton MCMS format)
func (t *GetValueV2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HasPendingUpgrade is a Record type
type HasPendingUpgrade struct {
	Viewer           types.PARTY `json:"viewer"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
}

// ToMap converts HasPendingUpgrade to a map for DAML arguments
func (t HasPendingUpgrade) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	m["targetInstanceId"] = string(t.TargetInstanceId)

	return m
}

func (t HasPendingUpgrade) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HasPendingUpgrade) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HasPendingUpgrade to hex string (Canton MCMS format)
func (t HasPendingUpgrade) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HasPendingUpgrade from hex string (Canton MCMS format)
func (t *HasPendingUpgrade) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InitiateUpgradeParams is a Record type
type InitiateUpgradeParams struct {
	InstanceId      types.TEXT `json:"instanceId"`
	FromVersion     types.TEXT `json:"fromVersion"`
	ToVersion       types.TEXT `json:"toVersion"`
	MigrationParams types.TEXT `json:"migrationParams"`
}

// ToMap converts InitiateUpgradeParams to a map for DAML arguments
func (t InitiateUpgradeParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["fromVersion"] = string(t.FromVersion)

	m["toVersion"] = string(t.ToVersion)

	m["migrationParams"] = string(t.MigrationParams)

	return m
}

func (t InitiateUpgradeParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InitiateUpgradeParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InitiateUpgradeParams to hex string (Canton MCMS format)
func (t InitiateUpgradeParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InitiateUpgradeParams from hex string (Canton MCMS format)
func (t *InitiateUpgradeParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsOperation is a Record type
type IsOperation struct {
	Submitter types.PARTY `json:"submitter"`
	OpId      types.TEXT  `json:"opId"`
}

// ToMap converts IsOperation to a map for DAML arguments
func (t IsOperation) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	return m
}

func (t IsOperation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsOperation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsOperation to hex string (Canton MCMS format)
func (t IsOperation) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsOperation from hex string (Canton MCMS format)
func (t *IsOperation) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsOperationDone is a Record type
type IsOperationDone struct {
	Submitter types.PARTY `json:"submitter"`
	OpId      types.TEXT  `json:"opId"`
}

// ToMap converts IsOperationDone to a map for DAML arguments
func (t IsOperationDone) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	return m
}

func (t IsOperationDone) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsOperationDone) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsOperationDone to hex string (Canton MCMS format)
func (t IsOperationDone) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsOperationDone from hex string (Canton MCMS format)
func (t *IsOperationDone) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsOperationPending is a Record type
type IsOperationPending struct {
	Submitter types.PARTY `json:"submitter"`
	OpId      types.TEXT  `json:"opId"`
}

// ToMap converts IsOperationPending to a map for DAML arguments
func (t IsOperationPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	return m
}

func (t IsOperationPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsOperationPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsOperationPending to hex string (Canton MCMS format)
func (t IsOperationPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsOperationPending from hex string (Canton MCMS format)
func (t *IsOperationPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsOperationReady is a Record type
type IsOperationReady struct {
	Submitter types.PARTY `json:"submitter"`
	OpId      types.TEXT  `json:"opId"`
}

// ToMap converts IsOperationReady to a map for DAML arguments
func (t IsOperationReady) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	return m
}

func (t IsOperationReady) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsOperationReady) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsOperationReady to hex string (Canton MCMS format)
func (t IsOperationReady) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsOperationReady from hex string (Canton MCMS format)
func (t *IsOperationReady) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IsRegistered is a Record type
type IsRegistered struct {
	Viewer           types.PARTY `json:"viewer"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
}

// ToMap converts IsRegistered to a map for DAML arguments
func (t IsRegistered) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	m["targetInstanceId"] = string(t.TargetInstanceId)

	return m
}

func (t IsRegistered) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsRegistered) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsRegistered to hex string (Canton MCMS format)
func (t IsRegistered) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsRegistered from hex string (Canton MCMS format)
func (t *IsRegistered) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ListRegistrations is a Record type
type ListRegistrations struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts ListRegistrations to a map for DAML arguments
func (t ListRegistrations) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t ListRegistrations) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ListRegistrations) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ListRegistrations to hex string (Canton MCMS format)
func (t ListRegistrations) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ListRegistrations from hex string (Canton MCMS format)
func (t *ListRegistrations) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMS is a Template type
type MCMS struct {
	Owner              types.PARTY        `json:"owner"`
	InstanceId         types.TEXT         `json:"instanceId"`
	ChainId            types.INT64        `json:"chainId"`
	Proposer           RoleState          `json:"proposer"`
	Canceller          RoleState          `json:"canceller"`
	Bypasser           RoleState          `json:"bypasser"`
	MinDelay           types.RELTIME      `json:"minDelay"`
	BlockedFunctions   []BlockedFunction  `json:"blockedFunctions"`
	TimelockTimestamps types.GENMAP       `json:"timelockTimestamps"`
	RegistryCid        *types.CONTRACT_ID `json:"registryCid"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MCMS) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MCMS) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.Main", "MCMS")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MCMS) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainId"] = int64(t.ChainId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["proposer"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Proposer).(mapper); ok {
			return m.toMap()
		}
		return t.Proposer
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["canceller"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Canceller).(mapper); ok {
			return m.toMap()
		}
		return t.Canceller
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["bypasser"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Bypasser).(mapper); ok {
			return m.toMap()
		}
		return t.Bypasser
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minDelay"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.MinDelay).(mapper); ok {
			return m.toMap()
		}
		return t.MinDelay
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["blockedFunctions"] = func() []any {
		res := make([]any, 0, len(t.BlockedFunctions))
		for _, e := range t.BlockedFunctions {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["timelockTimestamps"] = func() any {
		if t.TimelockTimestamps == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TimelockTimestamps}
	}()

	if t.RegistryCid != nil {
		args["registryCid"] = map[string]any{
			"_type": "optional",
			"value": *t.RegistryCid,
		}
	} else {
		args["registryCid"] = map[string]any{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MCMS) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainId"] = int64(t.ChainId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["proposer"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Proposer).(mapper); ok {
			return m.toMap()
		}
		return t.Proposer
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["canceller"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Canceller).(mapper); ok {
			return m.toMap()
		}
		return t.Canceller
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["bypasser"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Bypasser).(mapper); ok {
			return m.toMap()
		}
		return t.Bypasser
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minDelay"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.MinDelay).(mapper); ok {
			return m.toMap()
		}
		return t.MinDelay
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["blockedFunctions"] = func() []any {
		res := make([]any, 0, len(t.BlockedFunctions))
		for _, e := range t.BlockedFunctions {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["timelockTimestamps"] = func() any {
		if t.TimelockTimestamps == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TimelockTimestamps}
	}()

	if t.RegistryCid != nil {
		args["registryCid"] = map[string]any{
			"_type": "optional",
			"value": *t.RegistryCid,
		}
	} else {
		args["registryCid"] = map[string]any{
			"_type": "optional",
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MCMS) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMS) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMS to hex string (Canton MCMS format)
func (t MCMS) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMS from hex string (Canton MCMS format)
func (t *MCMS) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MCMS

// SetConfig exercises the SetConfig choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) SetConfig(contractID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// SetConfigWithPackageID exercises the SetConfig choice using the provided package ID instead of package name
func (t MCMS) SetConfigWithPackageID(contractID string, packageID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// SetRoot exercises the SetRoot choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) SetRoot(contractID string, args SetRoot) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetRoot",
		Arguments:  argsToMap(args),
	}
}

// SetRootWithPackageID exercises the SetRoot choice using the provided package ID instead of package name
func (t MCMS) SetRootWithPackageID(contractID string, packageID string, args SetRoot) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetRoot",
		Arguments:  argsToMap(args),
	}
}

// ExecuteOp exercises the ExecuteOp choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) ExecuteOp(contractID string, args ExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// ExecuteOpWithPackageID exercises the ExecuteOp choice using the provided package ID instead of package name
func (t MCMS) ExecuteOpWithPackageID(contractID string, packageID string, args ExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// ExecuteScheduledBatch exercises the ExecuteScheduledBatch choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) ExecuteScheduledBatch(contractID string, args ExecuteScheduledBatch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteScheduledBatch",
		Arguments:  argsToMap(args),
	}
}

// ExecuteScheduledBatchWithPackageID exercises the ExecuteScheduledBatch choice using the provided package ID instead of package name
func (t MCMS) ExecuteScheduledBatchWithPackageID(contractID string, packageID string, args ExecuteScheduledBatch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteScheduledBatch",
		Arguments:  argsToMap(args),
	}
}

// CanExecuteOp exercises the CanExecuteOp choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) CanExecuteOp(contractID string, args CanExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "CanExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// CanExecuteOpWithPackageID exercises the CanExecuteOp choice using the provided package ID instead of package name
func (t MCMS) CanExecuteOpWithPackageID(contractID string, packageID string, args CanExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "CanExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// IsOperationReady exercises the IsOperationReady choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) IsOperationReady(contractID string, args IsOperationReady) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationReady",
		Arguments:  argsToMap(args),
	}
}

// IsOperationReadyWithPackageID exercises the IsOperationReady choice using the provided package ID instead of package name
func (t MCMS) IsOperationReadyWithPackageID(contractID string, packageID string, args IsOperationReady) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationReady",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MCMS) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetState exercises the GetState choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetState(contractID string, args GetState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetState",
		Arguments:  argsToMap(args),
	}
}

// GetStateWithPackageID exercises the GetState choice using the provided package ID instead of package name
func (t MCMS) GetStateWithPackageID(contractID string, packageID string, args GetState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetState",
		Arguments:  argsToMap(args),
	}
}

// IsOperation exercises the IsOperation choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) IsOperation(contractID string, args IsOperation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperation",
		Arguments:  argsToMap(args),
	}
}

// IsOperationWithPackageID exercises the IsOperation choice using the provided package ID instead of package name
func (t MCMS) IsOperationWithPackageID(contractID string, packageID string, args IsOperation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperation",
		Arguments:  argsToMap(args),
	}
}

// IsOperationPending exercises the IsOperationPending choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) IsOperationPending(contractID string, args IsOperationPending) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationPending",
		Arguments:  argsToMap(args),
	}
}

// IsOperationPendingWithPackageID exercises the IsOperationPending choice using the provided package ID instead of package name
func (t MCMS) IsOperationPendingWithPackageID(contractID string, packageID string, args IsOperationPending) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationPending",
		Arguments:  argsToMap(args),
	}
}

// IsOperationDone exercises the IsOperationDone choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) IsOperationDone(contractID string, args IsOperationDone) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationDone",
		Arguments:  argsToMap(args),
	}
}

// IsOperationDoneWithPackageID exercises the IsOperationDone choice using the provided package ID instead of package name
func (t MCMS) IsOperationDoneWithPackageID(contractID string, packageID string, args IsOperationDone) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "IsOperationDone",
		Arguments:  argsToMap(args),
	}
}

// GetTimestamp exercises the GetTimestamp choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetTimestamp(contractID string, args GetTimestamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetTimestamp",
		Arguments:  argsToMap(args),
	}
}

// GetTimestampWithPackageID exercises the GetTimestamp choice using the provided package ID instead of package name
func (t MCMS) GetTimestampWithPackageID(contractID string, packageID string, args GetTimestamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetTimestamp",
		Arguments:  argsToMap(args),
	}
}

// GetMinDelay exercises the GetMinDelay choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetMinDelay(contractID string, args GetMinDelay) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetMinDelay",
		Arguments:  argsToMap(args),
	}
}

// GetMinDelayWithPackageID exercises the GetMinDelay choice using the provided package ID instead of package name
func (t MCMS) GetMinDelayWithPackageID(contractID string, packageID string, args GetMinDelay) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetMinDelay",
		Arguments:  argsToMap(args),
	}
}

// GetBlockedFunctions exercises the GetBlockedFunctions choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetBlockedFunctions(contractID string, args GetBlockedFunctions) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetBlockedFunctions",
		Arguments:  argsToMap(args),
	}
}

// GetBlockedFunctionsWithPackageID exercises the GetBlockedFunctions choice using the provided package ID instead of package name
func (t MCMS) GetBlockedFunctionsWithPackageID(contractID string, packageID string, args GetBlockedFunctions) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetBlockedFunctions",
		Arguments:  argsToMap(args),
	}
}

// GetBlockedFunctionsCount exercises the GetBlockedFunctionsCount choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetBlockedFunctionsCount(contractID string, args GetBlockedFunctionsCount) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetBlockedFunctionsCount",
		Arguments:  argsToMap(args),
	}
}

// GetBlockedFunctionsCountWithPackageID exercises the GetBlockedFunctionsCount choice using the provided package ID instead of package name
func (t MCMS) GetBlockedFunctionsCountWithPackageID(contractID string, packageID string, args GetBlockedFunctionsCount) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetBlockedFunctionsCount",
		Arguments:  argsToMap(args),
	}
}

// GetRegistryCid exercises the GetRegistryCid choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) GetRegistryCid(contractID string, args GetRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// GetRegistryCidWithPackageID exercises the GetRegistryCid choice using the provided package ID instead of package name
func (t MCMS) GetRegistryCidWithPackageID(contractID string, packageID string, args GetRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// SetRegistryCid exercises the SetRegistryCid choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) SetRegistryCid(contractID string, args SetRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// SetRegistryCidWithPackageID exercises the SetRegistryCid choice using the provided package ID instead of package name
func (t MCMS) SetRegistryCidWithPackageID(contractID string, packageID string, args SetRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// ClearRegistryCid exercises the ClearRegistryCid choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) ClearRegistryCid(contractID string, args ClearRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ClearRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// ClearRegistryCidWithPackageID exercises the ClearRegistryCid choice using the provided package ID instead of package name
func (t MCMS) ClearRegistryCidWithPackageID(contractID string, packageID string, args ClearRegistryCid) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ClearRegistryCid",
		Arguments:  argsToMap(args),
	}
}

// MCMSEntrypointEvent is a Template type
type MCMSEntrypointEvent struct {
	Owner             types.PARTY  `json:"owner"`
	InstanceId        types.TEXT   `json:"instanceId"`
	FunctionName      types.TEXT   `json:"functionName"`
	OperationData     types.TEXT   `json:"operationData"`
	ContractIdsAsText []types.TEXT `json:"contractIdsAsText"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MCMSEntrypointEvent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSEntrypointEvent")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MCMSEntrypointEvent) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.Counter", "MCMSEntrypointEvent")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MCMSEntrypointEvent) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["functionName"] = string(t.FunctionName)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operationData"] = string(t.OperationData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["contractIdsAsText"] = func() []any {
		res := make([]any, 0, len(t.ContractIdsAsText))
		for _, e := range t.ContractIdsAsText {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MCMSEntrypointEvent) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["functionName"] = string(t.FunctionName)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operationData"] = string(t.OperationData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["contractIdsAsText"] = func() []any {
		res := make([]any, 0, len(t.ContractIdsAsText))
		for _, e := range t.ContractIdsAsText {
			res = append(res, string(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSEntrypointEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSEntrypointEvent to hex string (Canton MCMS format)
func (t MCMSEntrypointEvent) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSEntrypointEvent from hex string (Canton MCMS format)
func (t *MCMSEntrypointEvent) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MCMSEntrypointEvent

// Archive exercises the Archive choice on this MCMSEntrypointEvent contract
// This method uses the package name in the template ID
func (t MCMSEntrypointEvent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MCMSEntrypointEvent) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveMCMSEntrypointEvent exercises the Archive_MCMSEntrypointEvent choice on this MCMSEntrypointEvent contract
// This method uses the package name in the template ID
func (t MCMSEntrypointEvent) ArchiveMCMSEntrypointEvent(contractID string, args ArchiveMCMSEntrypointEvent) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive_MCMSEntrypointEvent",
		Arguments:  argsToMap(args),
	}
}

// ArchiveMCMSEntrypointEventWithPackageID exercises the Archive_MCMSEntrypointEvent choice using the provided package ID instead of package name
func (t MCMSEntrypointEvent) ArchiveMCMSEntrypointEventWithPackageID(contractID string, packageID string, args ArchiveMCMSEntrypointEvent) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive_MCMSEntrypointEvent",
		Arguments:  argsToMap(args),
	}
}

// MCMSEntrypointResult is a Record type
type MCMSEntrypointResult struct {
	SelfCid    types.CONTRACT_ID   `json:"selfCid"`
	OutputCids []types.CONTRACT_ID `json:"outputCids"`
}

// ToMap converts MCMSEntrypointResult to a map for DAML arguments
func (t MCMSEntrypointResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["selfCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SelfCid).(mapper); ok {
			return m.toMap()
		}
		return t.SelfCid
	}()

	m["outputCids"] = func() []any {
		res := make([]any, 0, len(t.OutputCids))
		for _, e := range t.OutputCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t MCMSEntrypointResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSEntrypointResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSEntrypointResult to hex string (Canton MCMS format)
func (t MCMSEntrypointResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSEntrypointResult from hex string (Canton MCMS format)
func (t *MCMSEntrypointResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverView is a Record type
type MCMSReceiverView struct {
	Owner        types.PARTY `json:"owner"`
	InstanceId   types.TEXT  `json:"instanceId"`
	Version      types.TEXT  `json:"version"`
	ContractType types.TEXT  `json:"contractType"`
}

// ToMap converts MCMSReceiverView to a map for DAML arguments
func (t MCMSReceiverView) ToMap() map[string]any {
	m := make(map[string]any)

	m["owner"] = t.Owner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	m["version"] = string(t.Version)

	m["contractType"] = string(t.ContractType)

	return m
}

func (t MCMSReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverView to hex string (Canton MCMS format)
func (t MCMSReceiverView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverView from hex string (Canton MCMS format)
func (t *MCMSReceiverView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverEntrypoint is a Record type
type MCMSReceiverEntrypoint struct {
	Caller        types.PARTY         `json:"caller"`
	FunctionName  types.TEXT          `json:"functionName"`
	OperationData types.TEXT          `json:"operationData"`
	ContractIds   []types.CONTRACT_ID `json:"contractIds"`
}

// ToMap converts MCMSReceiverEntrypoint to a map for DAML arguments
func (t MCMSReceiverEntrypoint) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	m["contractIds"] = func() []any {
		res := make([]any, 0, len(t.ContractIds))
		for _, e := range t.ContractIds {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t MCMSReceiverEntrypoint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverEntrypoint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverEntrypoint to hex string (Canton MCMS format)
func (t MCMSReceiverEntrypoint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverEntrypoint from hex string (Canton MCMS format)
func (t *MCMSReceiverEntrypoint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverGetInstanceId is a Record type
type MCMSReceiverGetInstanceId struct {
	C types.PARTY `json:"c"`
}

// ToMap converts MCMSReceiverGetInstanceId to a map for DAML arguments
func (t MCMSReceiverGetInstanceId) ToMap() map[string]any {
	m := make(map[string]any)

	m["c"] = t.C.ToMap()

	return m
}

func (t MCMSReceiverGetInstanceId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverGetInstanceId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverGetInstanceId to hex string (Canton MCMS format)
func (t MCMSReceiverGetInstanceId) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverGetInstanceId from hex string (Canton MCMS format)
func (t *MCMSReceiverGetInstanceId) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverGetView is a Record type
type MCMSReceiverGetView struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts MCMSReceiverGetView to a map for DAML arguments
func (t MCMSReceiverGetView) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t MCMSReceiverGetView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverGetView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverGetView to hex string (Canton MCMS format)
func (t MCMSReceiverGetView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverGetView from hex string (Canton MCMS format)
func (t *MCMSReceiverGetView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSRegistry is a Template type
type MCMSRegistry struct {
	Owner           types.PARTY  `json:"owner"`
	InstanceId      types.TEXT   `json:"instanceId"`
	Registrations   types.GENMAP `json:"registrations"`
	PendingUpgrades types.GENMAP `json:"pendingUpgrades"`
	UpgradeHistory  types.GENMAP `json:"upgradeHistory"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MCMSRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MCMSRegistry) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MCMSRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrations"] = func() any {
		if t.Registrations == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.Registrations}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingUpgrades"] = func() any {
		if t.PendingUpgrades == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.PendingUpgrades}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["upgradeHistory"] = func() any {
		if t.UpgradeHistory == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UpgradeHistory}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MCMSRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrations"] = func() any {
		if t.Registrations == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.Registrations}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingUpgrades"] = func() any {
		if t.PendingUpgrades == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.PendingUpgrades}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["upgradeHistory"] = func() any {
		if t.UpgradeHistory == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UpgradeHistory}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MCMSRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSRegistry to hex string (Canton MCMS format)
func (t MCMSRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSRegistry from hex string (Canton MCMS format)
func (t *MCMSRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MCMSRegistry

// GetContractInfo exercises the GetContractInfo choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) GetContractInfo(contractID string, args GetContractInfo) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetContractInfo",
		Arguments:  argsToMap(args),
	}
}

// GetContractInfoWithPackageID exercises the GetContractInfo choice using the provided package ID instead of package name
func (t MCMSRegistry) GetContractInfoWithPackageID(contractID string, packageID string, args GetContractInfo) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetContractInfo",
		Arguments:  argsToMap(args),
	}
}

// GetPendingUpgrade exercises the GetPendingUpgrade choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) GetPendingUpgrade(contractID string, args GetPendingUpgrade) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetPendingUpgrade",
		Arguments:  argsToMap(args),
	}
}

// GetPendingUpgradeWithPackageID exercises the GetPendingUpgrade choice using the provided package ID instead of package name
func (t MCMSRegistry) GetPendingUpgradeWithPackageID(contractID string, packageID string, args GetPendingUpgrade) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetPendingUpgrade",
		Arguments:  argsToMap(args),
	}
}

// GetUpgradeHistory exercises the GetUpgradeHistory choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) GetUpgradeHistory(contractID string, args GetUpgradeHistory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetUpgradeHistory",
		Arguments:  argsToMap(args),
	}
}

// GetUpgradeHistoryWithPackageID exercises the GetUpgradeHistory choice using the provided package ID instead of package name
func (t MCMSRegistry) GetUpgradeHistoryWithPackageID(contractID string, packageID string, args GetUpgradeHistory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "GetUpgradeHistory",
		Arguments:  argsToMap(args),
	}
}

// IsRegistered exercises the IsRegistered choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) IsRegistered(contractID string, args IsRegistered) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "IsRegistered",
		Arguments:  argsToMap(args),
	}
}

// IsRegisteredWithPackageID exercises the IsRegistered choice using the provided package ID instead of package name
func (t MCMSRegistry) IsRegisteredWithPackageID(contractID string, packageID string, args IsRegistered) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "IsRegistered",
		Arguments:  argsToMap(args),
	}
}

// HasPendingUpgrade exercises the HasPendingUpgrade choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) HasPendingUpgrade(contractID string, args HasPendingUpgrade) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "HasPendingUpgrade",
		Arguments:  argsToMap(args),
	}
}

// HasPendingUpgradeWithPackageID exercises the HasPendingUpgrade choice using the provided package ID instead of package name
func (t MCMSRegistry) HasPendingUpgradeWithPackageID(contractID string, packageID string, args HasPendingUpgrade) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "HasPendingUpgrade",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MCMSRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ListRegistrations exercises the ListRegistrations choice on this MCMSRegistry contract
// This method uses the package name in the template ID
func (t MCMSRegistry) ListRegistrations(contractID string, args ListRegistrations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "ListRegistrations",
		Arguments:  argsToMap(args),
	}
}

// ListRegistrationsWithPackageID exercises the ListRegistrations choice using the provided package ID instead of package name
func (t MCMSRegistry) ListRegistrationsWithPackageID(contractID string, packageID string, args ListRegistrations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSRegistry"),
		ContractID: contractID,
		Choice:     "ListRegistrations",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this MCMSRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t MCMSRegistry) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t MCMSRegistry) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this MCMSRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t MCMSRegistry) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceIdWithPackageID exercises the MCMSReceiver_GetInstanceId choice using the provided package ID instead of package name
func (t MCMSRegistry) MCMSReceiverGetInstanceIdWithPackageID(contractID string, packageID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetView exercises the MCMSReceiver_GetView choice on this MCMSRegistry contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t MCMSRegistry) MCMSReceiverGetView(contractID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetViewWithPackageID exercises the MCMSReceiver_GetView choice using the provided package ID instead of package name
func (t MCMSRegistry) MCMSReceiverGetViewWithPackageID(contractID string, packageID string, args MCMSReceiverGetView) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSRegistry", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetView",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for MCMSRegistry

var _ IMCMSReceiver = (*MCMSRegistry)(nil)

// MCMSState is a Record type
type MCMSState struct {
	Role          Role            `json:"role"`
	OpCount       types.INT64     `json:"opCount"`
	PostOpCount   types.INT64     `json:"postOpCount"`
	ValidUntil    types.TIMESTAMP `json:"validUntil"`
	HasActiveRoot types.BOOL      `json:"hasActiveRoot"`
	NumSigners    types.INT64     `json:"numSigners"`
}

// ToMap converts MCMSState to a map for DAML arguments
func (t MCMSState) ToMap() map[string]any {
	m := make(map[string]any)

	m["role"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Role).(mapper); ok {
			return m.toMap()
		}
		return t.Role
	}()

	m["opCount"] = int64(t.OpCount)

	m["postOpCount"] = int64(t.PostOpCount)

	m["validUntil"] = t.ValidUntil

	m["hasActiveRoot"] = bool(t.HasActiveRoot)

	m["numSigners"] = int64(t.NumSigners)

	return m
}

func (t MCMSState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSState to hex string (Canton MCMS format)
func (t MCMSState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSState from hex string (Canton MCMS format)
func (t *MCMSState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MultisigConfig is a Record type
type MultisigConfig struct {
	Signers      []SignerInfo  `json:"signers"`
	GroupQuorums []types.INT64 `json:"groupQuorums"`
	GroupParents []types.INT64 `json:"groupParents"`
}

// ToMap converts MultisigConfig to a map for DAML arguments
func (t MultisigConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["groupQuorums"] = func() []any {
		res := make([]any, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []any {
		res := make([]any, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	return m
}

func (t MultisigConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MultisigConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MultisigConfig to hex string (Canton MCMS format)
func (t MultisigConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MultisigConfig from hex string (Canton MCMS format)
func (t *MultisigConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Op is a Record type
type Op struct {
	ChainId          types.INT64 `json:"chainId"`
	MultisigId       types.TEXT  `json:"multisigId"`
	Nonce            types.INT64 `json:"nonce"`
	TargetInstanceId types.TEXT  `json:"targetInstanceId"`
	FunctionName     types.TEXT  `json:"functionName"`
	OperationData    types.TEXT  `json:"operationData"`
}

// ToMap converts Op to a map for DAML arguments
func (t Op) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["nonce"] = int64(t.Nonce)

	m["targetInstanceId"] = string(t.TargetInstanceId)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	return m
}

func (t Op) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Op) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Op to hex string (Canton MCMS format)
func (t Op) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Op from hex string (Canton MCMS format)
func (t *Op) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PendingUpgrade is a Record type
type PendingUpgrade struct {
	InstanceId      types.TEXT      `json:"instanceId"`
	FromVersion     types.TEXT      `json:"fromVersion"`
	ToVersion       types.TEXT      `json:"toVersion"`
	InitiatedAt     types.TIMESTAMP `json:"initiatedAt"`
	MigrationParams types.TEXT      `json:"migrationParams"`
}

// ToMap converts PendingUpgrade to a map for DAML arguments
func (t PendingUpgrade) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["fromVersion"] = string(t.FromVersion)

	m["toVersion"] = string(t.ToVersion)

	m["initiatedAt"] = t.InitiatedAt

	m["migrationParams"] = string(t.MigrationParams)

	return m
}

func (t PendingUpgrade) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PendingUpgrade) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PendingUpgrade to hex string (Canton MCMS format)
func (t PendingUpgrade) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PendingUpgrade from hex string (Canton MCMS format)
func (t *PendingUpgrade) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RawSignature is a Record type
type RawSignature struct {
	PublicKey types.TEXT `json:"publicKey"`
	R         types.TEXT `json:"r"`
	S         types.TEXT `json:"s"`
}

// ToMap converts RawSignature to a map for DAML arguments
func (t RawSignature) ToMap() map[string]any {
	m := make(map[string]any)

	m["publicKey"] = string(t.PublicKey)

	m["r"] = string(t.R)

	m["s"] = string(t.S)

	return m
}

func (t RawSignature) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RawSignature) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RawSignature to hex string (Canton MCMS format)
func (t RawSignature) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RawSignature from hex string (Canton MCMS format)
func (t *RawSignature) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegisterContractParams is a Record type
type RegisterContractParams struct {
	InstanceId     types.TEXT   `json:"instanceId"`
	Version        types.TEXT   `json:"version"`
	ContractType   types.TEXT   `json:"contractType"`
	MetadataKeys   []types.TEXT `json:"metadataKeys"`
	MetadataValues []types.TEXT `json:"metadataValues"`
}

// ToMap converts RegisterContractParams to a map for DAML arguments
func (t RegisterContractParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["version"] = string(t.Version)

	m["contractType"] = string(t.ContractType)

	m["metadataKeys"] = func() []any {
		res := make([]any, 0, len(t.MetadataKeys))
		for _, e := range t.MetadataKeys {
			res = append(res, string(e))
		}
		return res
	}()

	m["metadataValues"] = func() []any {
		res := make([]any, 0, len(t.MetadataValues))
		for _, e := range t.MetadataValues {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t RegisterContractParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegisterContractParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegisterContractParams to hex string (Canton MCMS format)
func (t RegisterContractParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegisterContractParams from hex string (Canton MCMS format)
func (t *RegisterContractParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Role is an enum type
type Role string

const (
	RoleBypasser Role = "Bypasser"

	RoleCanceller Role = "Canceller"

	RoleProposer Role = "Proposer"
)

func (e Role) GetEnumConstructor() string { return string(e) }

func (e Role) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Types", "Role")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e Role) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Types", "Role")
}

func (e Role) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *Role) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes Role to hex string (Canton MCMS format)
func (e Role) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes Role from hex string (Canton MCMS format)
func (e *Role) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = Role("")

// RoleState is a Record type
type RoleState struct {
	Config       MultisigConfig `json:"config"`
	SeenHashes   types.GENMAP   `json:"seenHashes"`
	ExpiringRoot ExpiringRoot   `json:"expiringRoot"`
	RootMetadata RootMetadata   `json:"rootMetadata"`
}

// ToMap converts RoleState to a map for DAML arguments
func (t RoleState) ToMap() map[string]any {
	m := make(map[string]any)

	m["config"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	m["seenHashes"] = func() any {
		if t.SeenHashes == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SeenHashes}
	}()

	m["expiringRoot"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExpiringRoot).(mapper); ok {
			return m.toMap()
		}
		return t.ExpiringRoot
	}()

	m["rootMetadata"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RootMetadata).(mapper); ok {
			return m.toMap()
		}
		return t.RootMetadata
	}()

	return m
}

func (t RoleState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RoleState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RoleState to hex string (Canton MCMS format)
func (t RoleState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RoleState from hex string (Canton MCMS format)
func (t *RoleState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RootMetadata is a Record type
type RootMetadata struct {
	ChainId              types.INT64 `json:"chainId"`
	MultisigId           types.TEXT  `json:"multisigId"`
	PreOpCount           types.INT64 `json:"preOpCount"`
	PostOpCount          types.INT64 `json:"postOpCount"`
	OverridePreviousRoot types.BOOL  `json:"overridePreviousRoot"`
}

// ToMap converts RootMetadata to a map for DAML arguments
func (t RootMetadata) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["preOpCount"] = int64(t.PreOpCount)

	m["postOpCount"] = int64(t.PostOpCount)

	m["overridePreviousRoot"] = bool(t.OverridePreviousRoot)

	return m
}

func (t RootMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RootMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RootMetadata to hex string (Canton MCMS format)
func (t RootMetadata) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RootMetadata from hex string (Canton MCMS format)
func (t *RootMetadata) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ScheduleBatchParams is a Record type
type ScheduleBatchParams struct {
	Calls       []TimelockCall `json:"calls"`
	Predecessor types.TEXT     `json:"predecessor"`
	Salt        types.TEXT     `json:"salt"`
	DelaySecs   types.INT64    `json:"delaySecs"`
}

// ToMap converts ScheduleBatchParams to a map for DAML arguments
func (t ScheduleBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["predecessor"] = string(t.Predecessor)

	m["salt"] = string(t.Salt)

	m["delaySecs"] = int64(t.DelaySecs)

	return m
}

func (t ScheduleBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ScheduleBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ScheduleBatchParams to hex string (Canton MCMS format)
func (t ScheduleBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ScheduleBatchParams from hex string (Canton MCMS format)
func (t *ScheduleBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetConfig is a Record type
type SetConfig struct {
	TargetRole      Role          `json:"targetRole"`
	NewSigners      []SignerInfo  `json:"newSigners"`
	NewGroupQuorums []types.INT64 `json:"newGroupQuorums"`
	NewGroupParents []types.INT64 `json:"newGroupParents"`
	ClearRoot       types.BOOL    `json:"clearRoot"`
}

// ToMap converts SetConfig to a map for DAML arguments
func (t SetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TargetRole).(mapper); ok {
			return m.toMap()
		}
		return t.TargetRole
	}()

	m["newSigners"] = func() []any {
		res := make([]any, 0, len(t.NewSigners))
		for _, e := range t.NewSigners {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["newGroupQuorums"] = func() []any {
		res := make([]any, 0, len(t.NewGroupQuorums))
		for _, e := range t.NewGroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["newGroupParents"] = func() []any {
		res := make([]any, 0, len(t.NewGroupParents))
		for _, e := range t.NewGroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["clearRoot"] = bool(t.ClearRoot)

	return m
}

func (t SetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetConfig to hex string (Canton MCMS format)
func (t SetConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetConfig from hex string (Canton MCMS format)
func (t *SetConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetConfigParams is a Record type
type SetConfigParams struct {
	Signers      []SignerInfo  `json:"signers"`
	GroupQuorums []types.INT64 `json:"groupQuorums"`
	GroupParents []types.INT64 `json:"groupParents"`
	ClearRoot    types.BOOL    `json:"clearRoot"`
}

// ToMap converts SetConfigParams to a map for DAML arguments
func (t SetConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["groupQuorums"] = func() []any {
		res := make([]any, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []any {
		res := make([]any, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["clearRoot"] = bool(t.ClearRoot)

	return m
}

func (t SetConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetConfigParams to hex string (Canton MCMS format)
func (t SetConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetConfigParams from hex string (Canton MCMS format)
func (t *SetConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRegistryCid is a Record type
type SetRegistryCid struct {
	NewRegistryCid types.CONTRACT_ID `json:"newRegistryCid"`
}

// ToMap converts SetRegistryCid to a map for DAML arguments
func (t SetRegistryCid) ToMap() map[string]any {
	m := make(map[string]any)

	m["newRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.NewRegistryCid
	}()

	return m
}

func (t SetRegistryCid) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRegistryCid) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRegistryCid to hex string (Canton MCMS format)
func (t SetRegistryCid) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRegistryCid from hex string (Canton MCMS format)
func (t *SetRegistryCid) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetRoot is a Record type
type SetRoot struct {
	TargetRole    Role            `json:"targetRole"`
	Submitter     types.PARTY     `json:"submitter"`
	NewRoot       types.TEXT      `json:"newRoot"`
	ValidUntil    types.TIMESTAMP `json:"validUntil"`
	Metadata      RootMetadata    `json:"metadata"`
	MetadataProof []types.TEXT    `json:"metadataProof"`
	Signatures    []RawSignature  `json:"signatures"`
}

// ToMap converts SetRoot to a map for DAML arguments
func (t SetRoot) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TargetRole).(mapper); ok {
			return m.toMap()
		}
		return t.TargetRole
	}()

	m["submitter"] = t.Submitter.ToMap()

	m["newRoot"] = string(t.NewRoot)

	m["validUntil"] = t.ValidUntil

	m["metadata"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Metadata).(mapper); ok {
			return m.toMap()
		}
		return t.Metadata
	}()

	m["metadataProof"] = func() []any {
		res := make([]any, 0, len(t.MetadataProof))
		for _, e := range t.MetadataProof {
			res = append(res, string(e))
		}
		return res
	}()

	m["signatures"] = func() []any {
		res := make([]any, 0, len(t.Signatures))
		for _, e := range t.Signatures {
			type mapper interface{ toMap() map[string]any }
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

func (t SetRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetRoot to hex string (Canton MCMS format)
func (t SetRoot) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetRoot from hex string (Canton MCMS format)
func (t *SetRoot) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SignerInfo is a Record type
type SignerInfo struct {
	SignerAddress types.TEXT  `json:"signerAddress"`
	SignerIndex   types.INT64 `json:"signerIndex"`
	SignerGroup   types.INT64 `json:"signerGroup"`
}

// ToMap converts SignerInfo to a map for DAML arguments
func (t SignerInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["signerAddress"] = string(t.SignerAddress)

	m["signerIndex"] = int64(t.SignerIndex)

	m["signerGroup"] = int64(t.SignerGroup)

	return m
}

func (t SignerInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SignerInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SignerInfo to hex string (Canton MCMS format)
func (t SignerInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SignerInfo from hex string (Canton MCMS format)
func (t *SignerInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// StateMigrationTicket is a Template type
type StateMigrationTicket struct {
	Owner            types.PARTY     `json:"owner"`
	SourceInstanceId types.TEXT      `json:"sourceInstanceId"`
	SourceVersion    types.TEXT      `json:"sourceVersion"`
	TargetVersion    types.TEXT      `json:"targetVersion"`
	SerializedState  types.TEXT      `json:"serializedState"`
	CreatedAt        types.TIMESTAMP `json:"createdAt"`
	ExpiresAt        types.TIMESTAMP `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t StateMigrationTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "StateMigrationTicket")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t StateMigrationTicket) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.MCMSReceiver", "StateMigrationTicket")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t StateMigrationTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceInstanceId"] = string(t.SourceInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceVersion"] = string(t.SourceVersion)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["targetVersion"] = string(t.TargetVersion)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["serializedState"] = string(t.SerializedState)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["createdAt"] = t.CreatedAt

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t StateMigrationTicket) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceInstanceId"] = string(t.SourceInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sourceVersion"] = string(t.SourceVersion)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["targetVersion"] = string(t.TargetVersion)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["serializedState"] = string(t.SerializedState)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["createdAt"] = t.CreatedAt

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t StateMigrationTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *StateMigrationTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes StateMigrationTicket to hex string (Canton MCMS format)
func (t StateMigrationTicket) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes StateMigrationTicket from hex string (Canton MCMS format)
func (t *StateMigrationTicket) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for StateMigrationTicket

// ConsumeStateMigrationTicket exercises the ConsumeStateMigrationTicket choice on this StateMigrationTicket contract
// This method uses the package name in the template ID
func (t StateMigrationTicket) ConsumeStateMigrationTicket(contractID string, args ConsumeStateMigrationTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "ConsumeStateMigrationTicket",
		Arguments:  argsToMap(args),
	}
}

// ConsumeStateMigrationTicketWithPackageID exercises the ConsumeStateMigrationTicket choice using the provided package ID instead of package name
func (t StateMigrationTicket) ConsumeStateMigrationTicketWithPackageID(contractID string, packageID string, args ConsumeStateMigrationTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "ConsumeStateMigrationTicket",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this StateMigrationTicket contract
// This method uses the package name in the template ID
func (t StateMigrationTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t StateMigrationTicket) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetTicketInfo exercises the GetTicketInfo choice on this StateMigrationTicket contract
// This method uses the package name in the template ID
func (t StateMigrationTicket) GetTicketInfo(contractID string, args GetTicketInfo) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "GetTicketInfo",
		Arguments:  argsToMap(args),
	}
}

// GetTicketInfoWithPackageID exercises the GetTicketInfo choice using the provided package ID instead of package name
func (t StateMigrationTicket) GetTicketInfoWithPackageID(contractID string, packageID string, args GetTicketInfo) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.MCMSReceiver", "StateMigrationTicket"),
		ContractID: contractID,
		Choice:     "GetTicketInfo",
		Arguments:  argsToMap(args),
	}
}

// StateMigrationTicketInfo is a Record type
type StateMigrationTicketInfo struct {
	SourceInstanceId types.TEXT      `json:"sourceInstanceId"`
	SourceVersion    types.TEXT      `json:"sourceVersion"`
	TargetVersion    types.TEXT      `json:"targetVersion"`
	CreatedAt        types.TIMESTAMP `json:"createdAt"`
	ExpiresAt        types.TIMESTAMP `json:"expiresAt"`
	IsExpired        types.BOOL      `json:"isExpired"`
}

// ToMap converts StateMigrationTicketInfo to a map for DAML arguments
func (t StateMigrationTicketInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceInstanceId"] = string(t.SourceInstanceId)

	m["sourceVersion"] = string(t.SourceVersion)

	m["targetVersion"] = string(t.TargetVersion)

	m["createdAt"] = t.CreatedAt

	m["expiresAt"] = t.ExpiresAt

	m["isExpired"] = bool(t.IsExpired)

	return m
}

func (t StateMigrationTicketInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *StateMigrationTicketInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes StateMigrationTicketInfo to hex string (Canton MCMS format)
func (t StateMigrationTicketInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes StateMigrationTicketInfo from hex string (Canton MCMS format)
func (t *StateMigrationTicketInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TimelockCall is a Record type
type TimelockCall struct {
	TargetInstanceId types.TEXT `json:"targetInstanceId"`
	FunctionName     types.TEXT `json:"functionName"`
	OperationData    types.TEXT `json:"operationData"`
}

// ToMap converts TimelockCall to a map for DAML arguments
func (t TimelockCall) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceId"] = string(t.TargetInstanceId)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	return m
}

func (t TimelockCall) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TimelockCall) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TimelockCall to hex string (Canton MCMS format)
func (t TimelockCall) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TimelockCall from hex string (Canton MCMS format)
func (t *TimelockCall) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpgradeHistoryEntry is a Record type
type UpgradeHistoryEntry struct {
	FromVersion     types.TEXT      `json:"fromVersion"`
	ToVersion       types.TEXT      `json:"toVersion"`
	UpgradedAt      types.TIMESTAMP `json:"upgradedAt"`
	MigrationParams types.TEXT      `json:"migrationParams"`
}

// ToMap converts UpgradeHistoryEntry to a map for DAML arguments
func (t UpgradeHistoryEntry) ToMap() map[string]any {
	m := make(map[string]any)

	m["fromVersion"] = string(t.FromVersion)

	m["toVersion"] = string(t.ToVersion)

	m["upgradedAt"] = t.UpgradedAt

	m["migrationParams"] = string(t.MigrationParams)

	return m
}

func (t UpgradeHistoryEntry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpgradeHistoryEntry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpgradeHistoryEntry to hex string (Canton MCMS format)
func (t UpgradeHistoryEntry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpgradeHistoryEntry from hex string (Canton MCMS format)
func (t *UpgradeHistoryEntry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IMCMSReceiverInterfaceID returns the interface ID for the IMCMSReceiver interface using the package name
func IMCMSReceiverInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "MCMSReceiver")
}

// IMCMSReceiverInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IMCMSReceiverInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.MCMSReceiver", "MCMSReceiver")
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ArchiveCounterV1(args ArchiveCounterV1) (*bind.EncodedChoice, error)
	ArchiveMCMSEntrypointEvent(args ArchiveMCMSEntrypointEvent) (*bind.EncodedChoice, error)
	BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error)
	CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error)
	CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error)
	CancelUpgrade(args CancelUpgradeParams) (*bind.EncodedChoice, error)
	ClearRegistryCid(args ClearRegistryCid) (*bind.EncodedChoice, error)
	CompleteUpgrade(args CompleteUpgradeParams) (*bind.EncodedChoice, error)
	ConsumeStateMigrationTicket(args ConsumeStateMigrationTicket) (*bind.EncodedChoice, error)
	ExecuteOp(args ExecuteOp) (*bind.EncodedChoice, error)
	ExecuteScheduledBatch(args ExecuteScheduledBatch) (*bind.EncodedChoice, error)
	GetBlockedFunctions(args GetBlockedFunctions) (*bind.EncodedChoice, error)
	GetBlockedFunctionsCount(args GetBlockedFunctionsCount) (*bind.EncodedChoice, error)
	GetContractInfo(args GetContractInfo) (*bind.EncodedChoice, error)
	GetFullStateV2(args GetFullStateV2) (*bind.EncodedChoice, error)
	GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error)
	GetLabelV2(args GetLabelV2) (*bind.EncodedChoice, error)
	GetLastModifiedV2(args GetLastModifiedV2) (*bind.EncodedChoice, error)
	GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error)
	GetPendingUpgrade(args GetPendingUpgrade) (*bind.EncodedChoice, error)
	GetRegistryCid(args GetRegistryCid) (*bind.EncodedChoice, error)
	GetState(args GetState) (*bind.EncodedChoice, error)
	GetTicketInfo(args GetTicketInfo) (*bind.EncodedChoice, error)
	GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error)
	GetUpgradeHistory(args GetUpgradeHistory) (*bind.EncodedChoice, error)
	GetValue(args GetValue) (*bind.EncodedChoice, error)
	GetValueV2(args GetValueV2) (*bind.EncodedChoice, error)
	HasPendingUpgrade(args HasPendingUpgrade) (*bind.EncodedChoice, error)
	InitiateUpgrade(args InitiateUpgradeParams) (*bind.EncodedChoice, error)
	IsOperation(args IsOperation) (*bind.EncodedChoice, error)
	IsOperationDone(args IsOperationDone) (*bind.EncodedChoice, error)
	IsOperationPending(args IsOperationPending) (*bind.EncodedChoice, error)
	IsOperationReady(args IsOperationReady) (*bind.EncodedChoice, error)
	IsRegistered(args IsRegistered) (*bind.EncodedChoice, error)
	ListRegistrations(args ListRegistrations) (*bind.EncodedChoice, error)
	MCMSReceiverEntrypoint(args MCMSReceiverEntrypoint) (*bind.EncodedChoice, error)
	MCMSReceiverGetInstanceId(args MCMSReceiverGetInstanceId) (*bind.EncodedChoice, error)
	MCMSReceiverGetView(args MCMSReceiverGetView) (*bind.EncodedChoice, error)
	RegisterContract(args RegisterContractParams) (*bind.EncodedChoice, error)
	ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
	SetRegistryCid(args SetRegistryCid) (*bind.EncodedChoice, error)
	SetRoot(args SetRoot) (*bind.EncodedChoice, error)
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

// ArchiveCounterV1 encodes parameters for the ArchiveCounterV1 choice.
func (e *encoder) ArchiveCounterV1(args ArchiveCounterV1) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ArchiveCounterV1", args)
}

// ArchiveMCMSEntrypointEvent encodes parameters for the ArchiveMCMSEntrypointEvent choice.
func (e *encoder) ArchiveMCMSEntrypointEvent(args ArchiveMCMSEntrypointEvent) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ArchiveMCMSEntrypointEvent", args)
}

// BypasserExecuteBatch encodes parameters for the BypasserExecuteBatch choice.
func (e *encoder) BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BypasserExecuteBatch", args)
}

// CanExecuteOp encodes parameters for the CanExecuteOp choice.
func (e *encoder) CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CanExecuteOp", args)
}

// CancelBatch encodes parameters for the CancelBatch choice.
func (e *encoder) CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelBatch", args)
}

// CancelUpgrade encodes parameters for the CancelUpgrade choice.
func (e *encoder) CancelUpgrade(args CancelUpgradeParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelUpgrade", args)
}

// ClearRegistryCid encodes parameters for the ClearRegistryCid choice.
func (e *encoder) ClearRegistryCid(args ClearRegistryCid) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ClearRegistryCid", args)
}

// CompleteUpgrade encodes parameters for the CompleteUpgrade choice.
func (e *encoder) CompleteUpgrade(args CompleteUpgradeParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CompleteUpgrade", args)
}

// ConsumeStateMigrationTicket encodes parameters for the ConsumeStateMigrationTicket choice.
func (e *encoder) ConsumeStateMigrationTicket(args ConsumeStateMigrationTicket) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeStateMigrationTicket", args)
}

// ExecuteOp encodes parameters for the ExecuteOp choice.
func (e *encoder) ExecuteOp(args ExecuteOp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteOp", args)
}

// ExecuteScheduledBatch encodes parameters for the ExecuteScheduledBatch choice.
func (e *encoder) ExecuteScheduledBatch(args ExecuteScheduledBatch) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteScheduledBatch", args)
}

// GetBlockedFunctions encodes parameters for the GetBlockedFunctions choice.
func (e *encoder) GetBlockedFunctions(args GetBlockedFunctions) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetBlockedFunctions", args)
}

// GetBlockedFunctionsCount encodes parameters for the GetBlockedFunctionsCount choice.
func (e *encoder) GetBlockedFunctionsCount(args GetBlockedFunctionsCount) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetBlockedFunctionsCount", args)
}

// GetContractInfo encodes parameters for the GetContractInfo choice.
func (e *encoder) GetContractInfo(args GetContractInfo) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetContractInfo", args)
}

// GetFullStateV2 encodes parameters for the GetFullStateV2 choice.
func (e *encoder) GetFullStateV2(args GetFullStateV2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFullStateV2", args)
}

// GetInstanceIdChoice encodes parameters for the GetInstanceIdChoice choice.
func (e *encoder) GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetInstanceIdChoice", args)
}

// GetLabelV2 encodes parameters for the GetLabelV2 choice.
func (e *encoder) GetLabelV2(args GetLabelV2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetLabelV2", args)
}

// GetLastModifiedV2 encodes parameters for the GetLastModifiedV2 choice.
func (e *encoder) GetLastModifiedV2(args GetLastModifiedV2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetLastModifiedV2", args)
}

// GetMinDelay encodes parameters for the GetMinDelay choice.
func (e *encoder) GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetMinDelay", args)
}

// GetPendingUpgrade encodes parameters for the GetPendingUpgrade choice.
func (e *encoder) GetPendingUpgrade(args GetPendingUpgrade) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetPendingUpgrade", args)
}

// GetRegistryCid encodes parameters for the GetRegistryCid choice.
func (e *encoder) GetRegistryCid(args GetRegistryCid) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRegistryCid", args)
}

// GetState encodes parameters for the GetState choice.
func (e *encoder) GetState(args GetState) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetState", args)
}

// GetTicketInfo encodes parameters for the GetTicketInfo choice.
func (e *encoder) GetTicketInfo(args GetTicketInfo) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTicketInfo", args)
}

// GetTimestamp encodes parameters for the GetTimestamp choice.
func (e *encoder) GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTimestamp", args)
}

// GetUpgradeHistory encodes parameters for the GetUpgradeHistory choice.
func (e *encoder) GetUpgradeHistory(args GetUpgradeHistory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetUpgradeHistory", args)
}

// GetValue encodes parameters for the GetValue choice.
func (e *encoder) GetValue(args GetValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetValue", args)
}

// GetValueV2 encodes parameters for the GetValueV2 choice.
func (e *encoder) GetValueV2(args GetValueV2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetValueV2", args)
}

// HasPendingUpgrade encodes parameters for the HasPendingUpgrade choice.
func (e *encoder) HasPendingUpgrade(args HasPendingUpgrade) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HasPendingUpgrade", args)
}

// InitiateUpgrade encodes parameters for the InitiateUpgrade choice.
func (e *encoder) InitiateUpgrade(args InitiateUpgradeParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("InitiateUpgrade", args)
}

// IsOperation encodes parameters for the IsOperation choice.
func (e *encoder) IsOperation(args IsOperation) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsOperation", args)
}

// IsOperationDone encodes parameters for the IsOperationDone choice.
func (e *encoder) IsOperationDone(args IsOperationDone) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsOperationDone", args)
}

// IsOperationPending encodes parameters for the IsOperationPending choice.
func (e *encoder) IsOperationPending(args IsOperationPending) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsOperationPending", args)
}

// IsOperationReady encodes parameters for the IsOperationReady choice.
func (e *encoder) IsOperationReady(args IsOperationReady) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsOperationReady", args)
}

// IsRegistered encodes parameters for the IsRegistered choice.
func (e *encoder) IsRegistered(args IsRegistered) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsRegistered", args)
}

// ListRegistrations encodes parameters for the ListRegistrations choice.
func (e *encoder) ListRegistrations(args ListRegistrations) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ListRegistrations", args)
}

// MCMSReceiverEntrypoint encodes parameters for the MCMSReceiverEntrypoint choice.
func (e *encoder) MCMSReceiverEntrypoint(args MCMSReceiverEntrypoint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MCMSReceiverEntrypoint", args)
}

// MCMSReceiverGetInstanceId encodes parameters for the MCMSReceiverGetInstanceId choice.
func (e *encoder) MCMSReceiverGetInstanceId(args MCMSReceiverGetInstanceId) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MCMSReceiverGetInstanceId", args)
}

// MCMSReceiverGetView encodes parameters for the MCMSReceiverGetView choice.
func (e *encoder) MCMSReceiverGetView(args MCMSReceiverGetView) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MCMSReceiverGetView", args)
}

// RegisterContract encodes parameters for the RegisterContract choice.
func (e *encoder) RegisterContract(args RegisterContractParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegisterContract", args)
}

// ScheduleBatch encodes parameters for the ScheduleBatch choice.
func (e *encoder) ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ScheduleBatch", args)
}

// SetConfig encodes parameters for the SetConfig choice.
func (e *encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetRegistryCid encodes parameters for the SetRegistryCid choice.
func (e *encoder) SetRegistryCid(args SetRegistryCid) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRegistryCid", args)
}

// SetRoot encodes parameters for the SetRoot choice.
func (e *encoder) SetRoot(args SetRoot) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRoot", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
