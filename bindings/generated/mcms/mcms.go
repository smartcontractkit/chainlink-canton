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
	PackageID   = "b107f361de3ca57a7640526c214b148d43a5d4d61baa05f0e842c7a00eaaa382"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IMCMSReceiver is a DAML interface
type IMCMSReceiver interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// MCMSReceiverEntrypoint executes the MCMSReceiver_Entrypoint choice
	MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand
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
	ApGroupQuorums []types.INT64 `json:"apGroupQuorums" hex:"[]uint32"`
	ApGroupParents []types.INT64 `json:"apGroupParents" hex:"[]uint32"`
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
	TargetInstanceAddress types.TEXT `json:"targetInstanceAddress"`
	FunctionName          types.TEXT `json:"functionName"`
}

// ToMap converts BlockedFunction to a map for DAML arguments
func (t BlockedFunction) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

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

// GetInstanceAddressChoice exercises the GetInstanceAddressChoice choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetInstanceAddressChoice(contractID string, args GetInstanceAddressChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceAddressChoice",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceAddressChoiceWithPackageID exercises the GetInstanceAddressChoice choice using the provided package ID instead of package name
func (t Counter) GetInstanceAddressChoiceWithPackageID(contractID string, packageID string, args GetInstanceAddressChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceAddressChoice",
		Arguments:  argsToMap(args),
	}
}

// Increment exercises the Increment choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) Increment(contractID string, args Increment) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Increment",
		Arguments:  argsToMap(args),
	}
}

// IncrementWithPackageID exercises the Increment choice using the provided package ID instead of package name
func (t Counter) IncrementWithPackageID(contractID string, packageID string, args Increment) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Increment",
		Arguments:  argsToMap(args),
	}
}

// SetValue exercises the SetValue choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) SetValue(contractID string, args SetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "SetValue",
		Arguments:  argsToMap(args),
	}
}

// SetValueWithPackageID exercises the SetValue choice using the provided package ID instead of package name
func (t Counter) SetValueWithPackageID(contractID string, packageID string, args SetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "SetValue",
		Arguments:  argsToMap(args),
	}
}

// Reset exercises the Reset choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) Reset(contractID string, args Reset) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Reset",
		Arguments:  argsToMap(args),
	}
}

// ResetWithPackageID exercises the Reset choice using the provided package ID instead of package name
func (t Counter) ResetWithPackageID(contractID string, packageID string, args Reset) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Reset",
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

// Verify interface implementations for Counter

var _ IMCMSReceiver = (*Counter)(nil)

// ExecuteOp is a Record type
type ExecuteOp struct {
	TargetRole Role         `json:"targetRole"`
	Submitter  types.PARTY  `json:"submitter"`
	Op         Op           `json:"op"`
	OpProof    []types.TEXT `json:"opProof"`
	TargetCids types.GENMAP `json:"targetCids"`
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

	m["targetCids"] = func() any {
		if t.TargetCids == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TargetCids}
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
	Submitter   types.PARTY    `json:"submitter"`
	OpId        types.TEXT     `json:"opId"`
	Calls       []TimelockCall `json:"calls"`
	Predecessor types.TEXT     `json:"predecessor" hex:"bytes16"`
	Salt        types.TEXT     `json:"salt" hex:"bytes16"`
	TargetCids  types.GENMAP   `json:"targetCids"`
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

	m["targetCids"] = func() any {
		if t.TargetCids == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TargetCids}
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
	Root       types.TEXT      `json:"root" hex:"bytes"`
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

// GetInstanceAddressChoice is a Record type
type GetInstanceAddressChoice struct {
	Viewer types.PARTY `json:"viewer"`
}

// ToMap converts GetInstanceAddressChoice to a map for DAML arguments
func (t GetInstanceAddressChoice) ToMap() map[string]any {
	m := make(map[string]any)

	m["viewer"] = t.Viewer.ToMap()

	return m
}

func (t GetInstanceAddressChoice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetInstanceAddressChoice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetInstanceAddressChoice to hex string (Canton MCMS format)
func (t GetInstanceAddressChoice) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetInstanceAddressChoice from hex string (Canton MCMS format)
func (t *GetInstanceAddressChoice) UnmarshalHex(data string) error {
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

// Increment is a Record type
type Increment struct {
}

// ToMap converts Increment to a map for DAML arguments
func (t Increment) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t Increment) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Increment) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Increment to hex string (Canton MCMS format)
func (t Increment) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Increment from hex string (Canton MCMS format)
func (t *Increment) UnmarshalHex(data string) error {
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

// MCMS is a Template type
type MCMS struct {
	Owner              types.PARTY       `json:"owner"`
	InstanceId         types.TEXT        `json:"instanceId"`
	ChainId            types.INT64       `json:"chainId"`
	Proposer           RoleState         `json:"proposer"`
	Canceller          RoleState         `json:"canceller"`
	Bypasser           RoleState         `json:"bypasser"`
	MinDelay           types.RELTIME     `json:"minDelay"`
	BlockedFunctions   []BlockedFunction `json:"blockedFunctions"`
	TimelockTimestamps types.GENMAP      `json:"timelockTimestamps"`
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

// MCMSReceiverView is a Record type
type MCMSReceiverView struct {
	McmsController types.PARTY `json:"mcmsController"`
	InstanceId     types.TEXT  `json:"instanceId"`
}

// ToMap converts MCMSReceiverView to a map for DAML arguments
func (t MCMSReceiverView) ToMap() map[string]any {
	m := make(map[string]any)

	m["mcmsController"] = t.McmsController.ToMap()

	m["instanceId"] = string(t.InstanceId)

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
	FunctionName  types.TEXT   `json:"functionName"`
	OperationData types.TEXT   `json:"operationData" hex:"bytes16"`
	ContractIds   types.GENMAP `json:"contractIds"`
}

// ToMap converts MCMSReceiverEntrypoint to a map for DAML arguments
func (t MCMSReceiverEntrypoint) ToMap() map[string]any {
	m := make(map[string]any)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	m["contractIds"] = func() any {
		if t.ContractIds == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ContractIds}
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
	GroupQuorums []types.INT64 `json:"groupQuorums" hex:"[]uint32"`
	GroupParents []types.INT64 `json:"groupParents" hex:"[]uint32"`
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
	ChainId               types.INT64 `json:"chainId"`
	MultisigId            types.TEXT  `json:"multisigId"`
	Nonce                 types.INT64 `json:"nonce"`
	TargetInstanceAddress types.TEXT  `json:"targetInstanceAddress"`
	FunctionName          types.TEXT  `json:"functionName"`
	OperationData         types.TEXT  `json:"operationData" hex:"bytes16"`
}

// ToMap converts Op to a map for DAML arguments
func (t Op) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["nonce"] = int64(t.Nonce)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

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

// Reset is a Record type
type Reset struct {
}

// ToMap converts Reset to a map for DAML arguments
func (t Reset) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t Reset) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Reset) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Reset to hex string (Canton MCMS format)
func (t Reset) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Reset from hex string (Canton MCMS format)
func (t *Reset) UnmarshalHex(data string) error {
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
	Predecessor types.TEXT     `json:"predecessor" hex:"bytes16"`
	Salt        types.TEXT     `json:"salt" hex:"bytes16"`
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
	GroupQuorums []types.INT64 `json:"groupQuorums" hex:"[]uint32"`
	GroupParents []types.INT64 `json:"groupParents" hex:"[]uint32"`
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

// SetRoot is a Record type
type SetRoot struct {
	TargetRole    Role            `json:"targetRole"`
	Submitter     types.PARTY     `json:"submitter"`
	NewRoot       types.TEXT      `json:"newRoot" hex:"bytes"`
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

// SetValue is a Record type
type SetValue struct {
	NewValue types.INT64 `json:"newValue"`
}

// ToMap converts SetValue to a map for DAML arguments
func (t SetValue) ToMap() map[string]any {
	m := make(map[string]any)

	m["newValue"] = int64(t.NewValue)

	return m
}

func (t SetValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetValue to hex string (Canton MCMS format)
func (t SetValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetValue from hex string (Canton MCMS format)
func (t *SetValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SignerInfo is a Record type
type SignerInfo struct {
	SignerAddress types.TEXT  `json:"signerAddress" hex:"bytes"`
	SignerIndex   types.INT64 `json:"signerIndex" hex:"uint32"`
	SignerGroup   types.INT64 `json:"signerGroup" hex:"uint32"`
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

// TimelockCall is a Record type
type TimelockCall struct {
	TargetInstanceAddress types.TEXT `json:"targetInstanceAddress"`
	FunctionName          types.TEXT `json:"functionName"`
	OperationData         types.TEXT `json:"operationData" hex:"bytes16"`
}

// ToMap converts TimelockCall to a map for DAML arguments
func (t TimelockCall) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

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
	BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error)
	CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error)
	CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error)
	ExecuteOp(args ExecuteOp) (*bind.EncodedChoice, error)
	ExecuteScheduledBatch(args ExecuteScheduledBatch) (*bind.EncodedChoice, error)
	GetBlockedFunctions(args GetBlockedFunctions) (*bind.EncodedChoice, error)
	GetBlockedFunctionsCount(args GetBlockedFunctionsCount) (*bind.EncodedChoice, error)
	GetInstanceAddressChoice(args GetInstanceAddressChoice) (*bind.EncodedChoice, error)
	GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error)
	GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error)
	GetState(args GetState) (*bind.EncodedChoice, error)
	GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error)
	GetValue(args GetValue) (*bind.EncodedChoice, error)
	Increment(args Increment) (*bind.EncodedChoice, error)
	IsOperation(args IsOperation) (*bind.EncodedChoice, error)
	IsOperationDone(args IsOperationDone) (*bind.EncodedChoice, error)
	IsOperationPending(args IsOperationPending) (*bind.EncodedChoice, error)
	IsOperationReady(args IsOperationReady) (*bind.EncodedChoice, error)
	MCMSReceiverEntrypoint(args MCMSReceiverEntrypoint) (*bind.EncodedChoice, error)
	Reset(args Reset) (*bind.EncodedChoice, error)
	ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
	SetConfigParams(args SetConfigParams) (*bind.EncodedChoice, error)
	SetRoot(args SetRoot) (*bind.EncodedChoice, error)
	SetValue(args SetValue) (*bind.EncodedChoice, error)
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

// GetInstanceAddressChoice encodes parameters for the GetInstanceAddressChoice choice.
func (e *encoder) GetInstanceAddressChoice(args GetInstanceAddressChoice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetInstanceAddressChoice", args)
}

// GetInstanceIdChoice encodes parameters for the GetInstanceIdChoice choice.
func (e *encoder) GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetInstanceIdChoice", args)
}

// GetMinDelay encodes parameters for the GetMinDelay choice.
func (e *encoder) GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetMinDelay", args)
}

// GetState encodes parameters for the GetState choice.
func (e *encoder) GetState(args GetState) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetState", args)
}

// GetTimestamp encodes parameters for the GetTimestamp choice.
func (e *encoder) GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetTimestamp", args)
}

// GetValue encodes parameters for the GetValue choice.
func (e *encoder) GetValue(args GetValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetValue", args)
}

// Increment encodes parameters for the Increment choice.
func (e *encoder) Increment(args Increment) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Increment", args)
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

// MCMSReceiverEntrypoint encodes parameters for the MCMSReceiverEntrypoint choice.
func (e *encoder) MCMSReceiverEntrypoint(args MCMSReceiverEntrypoint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MCMSReceiverEntrypoint", args)
}

// Reset encodes parameters for the Reset choice.
func (e *encoder) Reset(args Reset) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Reset", args)
}

// ScheduleBatch encodes parameters for the ScheduleBatch choice.
func (e *encoder) ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ScheduleBatch", args)
}

// SetConfig encodes parameters for the SetConfig choice.
func (e *encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetConfigParams encodes parameters for the SetConfig choice.
func (e *encoder) SetConfigParams(args SetConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetRoot encodes parameters for the SetRoot choice.
func (e *encoder) SetRoot(args SetRoot) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRoot", args)
}

// SetValue encodes parameters for the SetValue choice.
func (e *encoder) SetValue(args SetValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetValue", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
