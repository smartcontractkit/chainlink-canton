package mcms

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

const PackageID = "c72fb09e1c4bfcf76004f3c8079de2ce0670d9e78c06aea1d48abfae80a5c199"
const SDKVersion = "3.4.10"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IMCMSReceiver is a DAML interface
type IMCMSReceiver interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// MCMSReceiverGetInstanceId executes the MCMSReceiver_GetInstanceId choice
	MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand

	// MCMSReceiverEntrypoint executes the MCMSReceiver_Entrypoint choice
	MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a ToMap method
	type mapper interface {
		ToMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// APSetConfig is a Record type
type APSetConfig struct {
	ApSigners      []SignerInfo `json:"apSigners"`
	ApGroupQuorums []INT64      `json:"apGroupQuorums"`
	ApGroupParents []INT64      `json:"apGroupParents"`
	ApClearRoot    BOOL         `json:"apClearRoot"`
}

// ToMap converts APSetConfig to a map for DAML arguments
func (t APSetConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["apSigners"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ApSigners))
		for _, e := range t.ApSigners {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["apGroupQuorums"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ApGroupQuorums))
		for _, e := range t.ApGroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apGroupParents"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ApGroupParents))
		for _, e := range t.ApGroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apClearRoot"] = bool(t.ApClearRoot)

	return m
}

// MarshalJSON implements custom JSON marshaling for APSetConfig using JsonCodec
func (t APSetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for APSetConfig using JsonCodec
func (t *APSetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// AdminParams is a variant/union type
type AdminParams struct {
	APSetConfig *SET  `json:"AP_SetConfig,omitempty"`
	APClearRoot *UNIT `json:"AP_ClearRoot,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AdminParams
func (v AdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for AdminParams
func (v *AdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
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
func (v AdminParams) GetVariantValue() interface{} {

	if v.APSetConfig != nil {
		return v.APSetConfig
	}

	if v.APClearRoot != nil {
		return v.APClearRoot
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*AdminParams)(nil)

// ArchiveMCMSEntrypointEvent is a Record type
type ArchiveMCMSEntrypointEvent struct {
}

// ToMap converts ArchiveMCMSEntrypointEvent to a map for DAML arguments
func (t ArchiveMCMSEntrypointEvent) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

// MarshalJSON implements custom JSON marshaling for ArchiveMCMSEntrypointEvent using JsonCodec
func (t ArchiveMCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ArchiveMCMSEntrypointEvent using JsonCodec
func (t *ArchiveMCMSEntrypointEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ArgValue is a variant/union type
type ArgValue struct {
	AVText  *TEXT      `json:"AV_Text,omitempty"`
	AVInt   *INT64     `json:"AV_Int,omitempty"`
	AVBool  *BOOL      `json:"AV_Bool,omitempty"`
	AVParty *PARTY     `json:"AV_Party,omitempty"`
	AVTime  *TIMESTAMP `json:"AV_Time,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ArgValue
func (v ArgValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for ArgValue
func (v *ArgValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
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
func (v ArgValue) GetVariantValue() interface{} {

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

// Verify interface implementation
var _ VARIANT = (*ArgValue)(nil)

// CanExecuteOp is a Record type
type CanExecuteOp struct {
	Submitter PARTY `json:"submitter"`
	Op        Op    `json:"op"`
}

// ToMap converts CanExecuteOp to a map for DAML arguments
func (t CanExecuteOp) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["submitter"] = t.Submitter.ToMap()

	m["op"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Op).(mapper); ok {
			return m.toMap()
		}
		return t.Op
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for CanExecuteOp using JsonCodec
func (t CanExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CanExecuteOp using JsonCodec
func (t *CanExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Counter is a Template type
type Counter struct {
	Owner      PARTY `json:"owner"`
	InstanceId TEXT  `json:"instanceId"`
	Value      INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template
func (t Counter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter")
}

// CreateCommand returns a CreateCommand for this template
func (t Counter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	args["value"] = int64(t.Value)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for Counter using JsonCodec
func (t Counter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Counter using JsonCodec
func (t *Counter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for Counter

// Archive exercises the Archive choice on this Counter contract
func (t Counter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// GetValue exercises the GetValue choice on this Counter contract
func (t Counter) GetValue(contractID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetValue",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceIdChoice exercises the GetInstanceIdChoice choice on this Counter contract
func (t Counter) GetInstanceIdChoice(contractID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this Counter contract via the IMCMSReceiver interface
func (t Counter) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this Counter contract via the IMCMSReceiver interface
func (t Counter) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Counter

var _ IMCMSReceiver = (*Counter)(nil)

// ExecuteMcmsOp is a Record type
type ExecuteMcmsOp struct {
	Submitter PARTY  `json:"submitter"`
	Op        Op     `json:"op"`
	OpProof   []TEXT `json:"opProof"`
}

// ToMap converts ExecuteMcmsOp to a map for DAML arguments
func (t ExecuteMcmsOp) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["submitter"] = t.Submitter.ToMap()

	m["op"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Op).(mapper); ok {
			return m.toMap()
		}
		return t.Op
	}()

	m["opProof"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OpProof))
		for _, e := range t.OpProof {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ExecuteMcmsOp using JsonCodec
func (t ExecuteMcmsOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecuteMcmsOp using JsonCodec
func (t *ExecuteMcmsOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecuteOp is a Record type
type ExecuteOp struct {
	Submitter   PARTY         `json:"submitter"`
	TargetCid   CONTRACT_ID   `json:"targetCid"`
	Op          Op            `json:"op"`
	OpProof     []TEXT        `json:"opProof"`
	ContractIds []CONTRACT_ID `json:"contractIds"`
}

// ToMap converts ExecuteOp to a map for DAML arguments
func (t ExecuteOp) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["submitter"] = t.Submitter.ToMap()

	m["targetCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TargetCid).(mapper); ok {
			return m.toMap()
		}
		return t.TargetCid
	}()

	m["op"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Op).(mapper); ok {
			return m.toMap()
		}
		return t.Op
	}()

	m["opProof"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OpProof))
		for _, e := range t.OpProof {
			res = append(res, string(e))
		}
		return res
	}()

	m["contractIds"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ContractIds))
		for _, e := range t.ContractIds {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ExecuteOp using JsonCodec
func (t ExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecuteOp using JsonCodec
func (t *ExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExpiringRoot is a Record type
type ExpiringRoot struct {
	Root       TEXT      `json:"root"`
	ValidUntil TIMESTAMP `json:"validUntil"`
	OpCount    INT64     `json:"opCount"`
}

// ToMap converts ExpiringRoot to a map for DAML arguments
func (t ExpiringRoot) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["root"] = string(t.Root)

	m["validUntil"] = t.ValidUntil

	m["opCount"] = int64(t.OpCount)

	return m
}

// MarshalJSON implements custom JSON marshaling for ExpiringRoot using JsonCodec
func (t ExpiringRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExpiringRoot using JsonCodec
func (t *ExpiringRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetInstanceIdChoice is a Record type
type GetInstanceIdChoice struct {
	Viewer PARTY `json:"viewer"`
}

// ToMap converts GetInstanceIdChoice to a map for DAML arguments
func (t GetInstanceIdChoice) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["viewer"] = t.Viewer.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for GetInstanceIdChoice using JsonCodec
func (t GetInstanceIdChoice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetInstanceIdChoice using JsonCodec
func (t *GetInstanceIdChoice) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetState is a Record type
type GetState struct {
	Submitter PARTY `json:"submitter"`
}

// ToMap converts GetState to a map for DAML arguments
func (t GetState) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["submitter"] = t.Submitter.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for GetState using JsonCodec
func (t GetState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetState using JsonCodec
func (t *GetState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetValue is a Record type
type GetValue struct {
	Viewer PARTY `json:"viewer"`
}

// ToMap converts GetValue to a map for DAML arguments
func (t GetValue) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["viewer"] = t.Viewer.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for GetValue using JsonCodec
func (t GetValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetValue using JsonCodec
func (t *GetValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMS is a Template type
type MCMS struct {
	Owner        PARTY          `json:"owner"`
	Role         Role           `json:"role"`
	ChainId      INT64          `json:"chainId"`
	McmsId       TEXT           `json:"mcmsId"`
	Config       MultisigConfig `json:"config"`
	SeenHashes   GENMAP         `json:"seenHashes"`
	ExpiringRoot ExpiringRoot   `json:"expiringRoot"`
	RootMetadata RootMetadata   `json:"rootMetadata"`
}

// GetTemplateID returns the template ID for this template
func (t MCMS) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS")
}

// CreateCommand returns a CreateCommand for this template
func (t MCMS) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	if t.Role != "" {
		args["role"] = func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Role).(mapper); ok {
				return m.toMap()
			}
			return t.Role
		}()
	}

	args["chainId"] = int64(t.ChainId)

	args["mcmsId"] = string(t.McmsId)

	args["config"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	if t.SeenHashes != nil && len(t.SeenHashes) > 0 {
		args["seenHashes"] = map[string]interface{}{"_type": "genmap", "value": t.SeenHashes}
	}

	args["expiringRoot"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExpiringRoot).(mapper); ok {
			return m.toMap()
		}
		return t.ExpiringRoot
	}()

	args["rootMetadata"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RootMetadata).(mapper); ok {
			return m.toMap()
		}
		return t.RootMetadata
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MCMS using JsonCodec
func (t MCMS) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMS using JsonCodec
func (t *MCMS) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MCMS

// SetRoot exercises the SetRoot choice on this MCMS contract
func (t MCMS) SetRoot(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetRoot",
		Arguments:  argsToMap(args),
	}
}

// ExecuteOp exercises the ExecuteOp choice on this MCMS contract
func (t MCMS) ExecuteOp(contractID string, args ExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// ExecuteMcmsOp exercises the ExecuteMcmsOp choice on this MCMS contract
func (t MCMS) ExecuteMcmsOp(contractID string, args ExecuteMcmsOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "ExecuteMcmsOp",
		Arguments:  argsToMap(args),
	}
}

// SetConfig exercises the SetConfig choice on this MCMS contract
func (t MCMS) SetConfig(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// CanExecuteOp exercises the CanExecuteOp choice on this MCMS contract
func (t MCMS) CanExecuteOp(contractID string, args CanExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "CanExecuteOp",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMS contract
func (t MCMS) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// GetState exercises the GetState choice on this MCMS contract
func (t MCMS) GetState(contractID string, args GetState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "GetState",
		Arguments:  argsToMap(args),
	}
}

// MCMSEntrypointEvent is a Template type
type MCMSEntrypointEvent struct {
	Owner             PARTY  `json:"owner"`
	InstanceId        TEXT   `json:"instanceId"`
	FunctionName      TEXT   `json:"functionName"`
	OperationData     TEXT   `json:"operationData"`
	ContractIdsAsText []TEXT `json:"contractIdsAsText"`
}

// GetTemplateID returns the template ID for this template
func (t MCMSEntrypointEvent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSEntrypointEvent")
}

// CreateCommand returns a CreateCommand for this template
func (t MCMSEntrypointEvent) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	args["functionName"] = string(t.FunctionName)

	args["operationData"] = string(t.OperationData)

	if len(t.ContractIdsAsText) > 0 {
		args["contractIdsAsText"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.ContractIdsAsText))
			for _, e := range t.ContractIdsAsText {
				res = append(res, string(e))
			}
			return res
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MCMSEntrypointEvent using JsonCodec
func (t MCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMSEntrypointEvent using JsonCodec
func (t *MCMSEntrypointEvent) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MCMSEntrypointEvent

// Archive exercises the Archive choice on this MCMSEntrypointEvent contract
func (t MCMSEntrypointEvent) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveMCMSEntrypointEvent exercises the Archive_MCMSEntrypointEvent choice on this MCMSEntrypointEvent contract
func (t MCMSEntrypointEvent) ArchiveMCMSEntrypointEvent(contractID string, args ArchiveMCMSEntrypointEvent) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSEntrypointEvent"),
		ContractID: contractID,
		Choice:     "Archive_MCMSEntrypointEvent",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverView is a Record type
type MCMSReceiverView struct {
	Owner      PARTY `json:"owner"`
	InstanceId TEXT  `json:"instanceId"`
}

// ToMap converts MCMSReceiverView to a map for DAML arguments
func (t MCMSReceiverView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

// MarshalJSON implements custom JSON marshaling for MCMSReceiverView using JsonCodec
func (t MCMSReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMSReceiverView using JsonCodec
func (t *MCMSReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSReceiverEntrypoint is a Record type
type MCMSReceiverEntrypoint struct {
	Caller        PARTY         `json:"caller"`
	FunctionName  TEXT          `json:"functionName"`
	OperationData TEXT          `json:"operationData"`
	ContractIds   []CONTRACT_ID `json:"contractIds"`
}

// ToMap converts MCMSReceiverEntrypoint to a map for DAML arguments
func (t MCMSReceiverEntrypoint) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	m["contractIds"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ContractIds))
		for _, e := range t.ContractIds {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for MCMSReceiverEntrypoint using JsonCodec
func (t MCMSReceiverEntrypoint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMSReceiverEntrypoint using JsonCodec
func (t *MCMSReceiverEntrypoint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSReceiverGetInstanceId is a Record type
type MCMSReceiverGetInstanceId struct {
	C PARTY `json:"c"`
}

// ToMap converts MCMSReceiverGetInstanceId to a map for DAML arguments
func (t MCMSReceiverGetInstanceId) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["c"] = t.C.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for MCMSReceiverGetInstanceId using JsonCodec
func (t MCMSReceiverGetInstanceId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMSReceiverGetInstanceId using JsonCodec
func (t *MCMSReceiverGetInstanceId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSState is a Record type
type MCMSState struct {
	Role          Role      `json:"role"`
	OpCount       INT64     `json:"opCount"`
	PostOpCount   INT64     `json:"postOpCount"`
	ValidUntil    TIMESTAMP `json:"validUntil"`
	HasActiveRoot BOOL      `json:"hasActiveRoot"`
	NumSigners    INT64     `json:"numSigners"`
}

// ToMap converts MCMSState to a map for DAML arguments
func (t MCMSState) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["role"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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

// MarshalJSON implements custom JSON marshaling for MCMSState using JsonCodec
func (t MCMSState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MCMSState using JsonCodec
func (t *MCMSState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MultisigConfig is a Record type
type MultisigConfig struct {
	Signers      []SignerInfo `json:"signers"`
	GroupQuorums []INT64      `json:"groupQuorums"`
	GroupParents []INT64      `json:"groupParents"`
}

// ToMap converts MultisigConfig to a map for DAML arguments
func (t MultisigConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["signers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Signers))
		for _, e := range t.Signers {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["groupQuorums"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for MultisigConfig using JsonCodec
func (t MultisigConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MultisigConfig using JsonCodec
func (t *MultisigConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Op is a Record type
type Op struct {
	ChainId          INT64 `json:"chainId"`
	MultisigId       TEXT  `json:"multisigId"`
	Nonce            INT64 `json:"nonce"`
	TargetInstanceId TEXT  `json:"targetInstanceId"`
	FunctionName     TEXT  `json:"functionName"`
	OperationData    TEXT  `json:"operationData"`
}

// ToMap converts Op to a map for DAML arguments
func (t Op) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["nonce"] = int64(t.Nonce)

	m["targetInstanceId"] = string(t.TargetInstanceId)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	return m
}

// MarshalJSON implements custom JSON marshaling for Op using JsonCodec
func (t Op) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Op using JsonCodec
func (t *Op) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RawSignature is a Record type
type RawSignature struct {
	PublicKey TEXT `json:"publicKey"`
	R         TEXT `json:"r"`
	S         TEXT `json:"s"`
}

// ToMap converts RawSignature to a map for DAML arguments
func (t RawSignature) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["publicKey"] = string(t.PublicKey)

	m["r"] = string(t.R)

	m["s"] = string(t.S)

	return m
}

// MarshalJSON implements custom JSON marshaling for RawSignature using JsonCodec
func (t RawSignature) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RawSignature using JsonCodec
func (t *RawSignature) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Role is an enum type
type Role string

const (
	RoleProposer Role = "Proposer"
)

// GetEnumConstructor implements types.ENUM interface
func (e Role) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e Role) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Types", "Role")
}

// MarshalJSON implements custom JSON marshaling for Role using JsonCodec
func (e Role) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for Role using JsonCodec
func (e *Role) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = Role("")

// RootMetadata is a Record type
type RootMetadata struct {
	ChainId              INT64 `json:"chainId"`
	MultisigId           TEXT  `json:"multisigId"`
	PreOpCount           INT64 `json:"preOpCount"`
	PostOpCount          INT64 `json:"postOpCount"`
	OverridePreviousRoot BOOL  `json:"overridePreviousRoot"`
}

// ToMap converts RootMetadata to a map for DAML arguments
func (t RootMetadata) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["preOpCount"] = int64(t.PreOpCount)

	m["postOpCount"] = int64(t.PostOpCount)

	m["overridePreviousRoot"] = bool(t.OverridePreviousRoot)

	return m
}

// MarshalJSON implements custom JSON marshaling for RootMetadata using JsonCodec
func (t RootMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for RootMetadata using JsonCodec
func (t *RootMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetConfig is a Record type
type SetConfig struct {
	NewSigners      []SignerInfo `json:"newSigners"`
	NewGroupQuorums []INT64      `json:"newGroupQuorums"`
	NewGroupParents []INT64      `json:"newGroupParents"`
	ClearRoot       BOOL         `json:"clearRoot"`
}

// ToMap converts SetConfig to a map for DAML arguments
func (t SetConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newSigners"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.NewSigners))
		for _, e := range t.NewSigners {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["newGroupQuorums"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.NewGroupQuorums))
		for _, e := range t.NewGroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["newGroupParents"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.NewGroupParents))
		for _, e := range t.NewGroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["clearRoot"] = bool(t.ClearRoot)

	return m
}

// MarshalJSON implements custom JSON marshaling for SetConfig using JsonCodec
func (t SetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SetConfig using JsonCodec
func (t *SetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetConfigParams is a Record type
type SetConfigParams struct {
	Signers      []SignerInfo `json:"signers"`
	GroupQuorums []INT64      `json:"groupQuorums"`
	GroupParents []INT64      `json:"groupParents"`
	ClearRoot    BOOL         `json:"clearRoot"`
}

// ToMap converts SetConfigParams to a map for DAML arguments
func (t SetConfigParams) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["signers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Signers))
		for _, e := range t.Signers {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["groupQuorums"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["clearRoot"] = bool(t.ClearRoot)

	return m
}

// MarshalJSON implements custom JSON marshaling for SetConfigParams using JsonCodec
func (t SetConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SetConfigParams using JsonCodec
func (t *SetConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetRoot is a Record type
type SetRoot struct {
	Submitter     PARTY          `json:"submitter"`
	NewRoot       TEXT           `json:"newRoot"`
	ValidUntil    TIMESTAMP      `json:"validUntil"`
	Metadata      RootMetadata   `json:"metadata"`
	MetadataProof []TEXT         `json:"metadataProof"`
	Signatures    []RawSignature `json:"signatures"`
}

// ToMap converts SetRoot to a map for DAML arguments
func (t SetRoot) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["submitter"] = t.Submitter.ToMap()

	m["newRoot"] = string(t.NewRoot)

	m["validUntil"] = t.ValidUntil

	m["metadata"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Metadata).(mapper); ok {
			return m.toMap()
		}
		return t.Metadata
	}()

	m["metadataProof"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.MetadataProof))
		for _, e := range t.MetadataProof {
			res = append(res, string(e))
		}
		return res
	}()

	m["signatures"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Signatures))
		for _, e := range t.Signatures {
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

// MarshalJSON implements custom JSON marshaling for SetRoot using JsonCodec
func (t SetRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SetRoot using JsonCodec
func (t *SetRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SignerInfo is a Record type
type SignerInfo struct {
	SignerAddress TEXT  `json:"signerAddress"`
	SignerIndex   INT64 `json:"signerIndex"`
	SignerGroup   INT64 `json:"signerGroup"`
}

// ToMap converts SignerInfo to a map for DAML arguments
func (t SignerInfo) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["signerAddress"] = string(t.SignerAddress)

	m["signerIndex"] = int64(t.SignerIndex)

	m["signerGroup"] = int64(t.SignerGroup)

	return m
}

// MarshalJSON implements custom JSON marshaling for SignerInfo using JsonCodec
func (t SignerInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for SignerInfo using JsonCodec
func (t *SignerInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IMCMSReceiverInterfaceID returns the interface ID for the IMCMSReceiver interface
func IMCMSReceiverInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "MCMS.MCMSReceiver", "MCMSReceiver")
}
