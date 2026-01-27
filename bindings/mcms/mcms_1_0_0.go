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

const PackageID = "fbb6f00c5aa22cb2efa84283b86a80b04fd6c7458bc75588a495771f8f1bc526"
const SDKVersion = "3.4.9"

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

	// MCMSReceiverIssueAuthTicket executes the MCMSReceiver_IssueAuthTicket choice
	MCMSReceiverIssueAuthTicket(contractID string, args MCMSReceiverIssueAuthTicket) *model.ExerciseCommand
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	type mapper interface {
		ToMap() map[string]interface{}
	}
	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{"args": args}
}

// APSetConfig is a Record type
type APSetConfig struct {
	ApSigners []SignerInfo `json:"apSigners"`

	ApGroupQuorums []INT64 `json:"apGroupQuorums"`

	ApGroupParents []INT64 `json:"apGroupParents"`

	ApClearRoot BOOL `json:"apClearRoot"`
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

func (t APSetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

var _ VARIANT = (*AdminParams)(nil)

// ArchiveMCMSEntrypointEvent is a Record type
type ArchiveMCMSEntrypointEvent struct {
}

// ToMap converts ArchiveMCMSEntrypointEvent to a map for DAML arguments
func (t ArchiveMCMSEntrypointEvent) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t ArchiveMCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

var _ VARIANT = (*ArgValue)(nil)

// CanExecuteOp is a Record type
type CanExecuteOp struct {
	Submitter PARTY `json:"submitter"`

	Op Op `json:"op"`
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

func (t CanExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CanExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Consume is a Record type
type Consume struct {
}

// ToMap converts Consume to a map for DAML arguments
func (t Consume) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	return m
}

func (t Consume) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Consume) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Counter is a Template type
type Counter struct {
	Owner PARTY `json:"owner"`

	InstanceId TEXT `json:"instanceId"`

	Value INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template
func (t Counter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter")
}

// CreateCommand returns a CreateCommand for this template
func (t Counter) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

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

func (t Counter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: map[string]interface{}{},
	}
}

// GetValue exercises the GetValue choice on this Counter contract
func (t Counter) GetValue(contractID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter"),

		ContractID: contractID,
		Choice:     "GetValue",

		Arguments: argsToMap(args),
	}
}

// GetInstanceIdChoice exercises the GetInstanceIdChoice choice on this Counter contract
func (t Counter) GetInstanceIdChoice(contractID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "Counter"),

		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",

		Arguments: argsToMap(args),
	}
}

// MCMSReceiverGetInstanceId exercises the MCMSReceiver_GetInstanceId choice on this Counter contract via the IMCMSReceiver interface
func (t Counter) MCMSReceiverGetInstanceId(contractID string, args MCMSReceiverGetInstanceId) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSReceiver"),

		ContractID: contractID,
		Choice:     "MCMSReceiver_GetInstanceId",

		Arguments: argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this Counter contract via the IMCMSReceiver interface
func (t Counter) MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSReceiver"),

		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",

		Arguments: argsToMap(args),
	}
}

// MCMSReceiverIssueAuthTicket exercises the MCMSReceiver_IssueAuthTicket choice on this Counter contract via the IMCMSReceiver interface
func (t Counter) MCMSReceiverIssueAuthTicket(contractID string, args MCMSReceiverIssueAuthTicket) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSReceiver"),

		ContractID: contractID,
		Choice:     "MCMSReceiver_IssueAuthTicket",

		Arguments: argsToMap(args),
	}
}

// Verify interface implementations for Counter

var _ IMCMSReceiver = (*Counter)(nil)

// ExecuteMcmsOp is a Record type
type ExecuteMcmsOp struct {
	Submitter PARTY `json:"submitter"`

	Op Op `json:"op"`

	OpProof []TEXT `json:"opProof"`
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

func (t ExecuteMcmsOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ExecuteMcmsOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecuteOp is a Record type
type ExecuteOp struct {
	Submitter PARTY `json:"submitter"`

	TargetCid CONTRACT_ID `json:"targetCid"`

	AuthTicketCid CONTRACT_ID `json:"authTicketCid"`

	Op Op `json:"op"`

	OpProof []TEXT `json:"opProof"`

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

	m["authTicketCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.AuthTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.AuthTicketCid
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

func (t ExecuteOp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ExecuteOp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExpiringRoot is a Record type
type ExpiringRoot struct {
	Root TEXT `json:"root"`

	ValidUntil TIMESTAMP `json:"validUntil"`

	OpCount INT64 `json:"opCount"`
}

// ToMap converts ExpiringRoot to a map for DAML arguments
func (t ExpiringRoot) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["root"] = string(t.Root)

	m["validUntil"] = t.ValidUntil

	m["opCount"] = int64(t.OpCount)

	return m
}

func (t ExpiringRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t GetInstanceIdChoice) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t GetState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t GetValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Lookup is a Record type
type Lookup struct {
	Caller PARTY `json:"caller"`

	InstanceId TEXT `json:"instanceId"`
}

// ToMap converts Lookup to a map for DAML arguments
func (t Lookup) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t Lookup) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Lookup) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMS is a Template type
type MCMS struct {
	Owner PARTY `json:"owner"`

	Role Role `json:"role"`

	ChainId INT64 `json:"chainId"`

	McmsId TEXT `json:"mcmsId"`

	Config MultisigConfig `json:"config"`

	SeenHashes GENMAP `json:"seenHashes"`

	ExpiringRoot ExpiringRoot `json:"expiringRoot"`

	RootMetadata RootMetadata `json:"rootMetadata"`
}

// GetTemplateID returns the template ID for this template
func (t MCMS) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS")
}

// CreateCommand returns a CreateCommand for this template
func (t MCMS) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
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

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["chainId"] = int64(t.ChainId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsId"] = string(t.McmsId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["config"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["seenHashes"] = func() interface{} {
		if t.SeenHashes == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.SeenHashes}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiringRoot"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExpiringRoot).(mapper); ok {
			return m.toMap()
		}
		return t.ExpiringRoot
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
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

func (t MCMS) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: argsToMap(args),
	}
}

// ExecuteOp exercises the ExecuteOp choice on this MCMS contract
func (t MCMS) ExecuteOp(contractID string, args ExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "ExecuteOp",

		Arguments: argsToMap(args),
	}
}

// ExecuteMcmsOp exercises the ExecuteMcmsOp choice on this MCMS contract
func (t MCMS) ExecuteMcmsOp(contractID string, args ExecuteMcmsOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "ExecuteMcmsOp",

		Arguments: argsToMap(args),
	}
}

// SetConfig exercises the SetConfig choice on this MCMS contract
func (t MCMS) SetConfig(contractID string, args SET) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "SetConfig",

		Arguments: argsToMap(args),
	}
}

// CanExecuteOp exercises the CanExecuteOp choice on this MCMS contract
func (t MCMS) CanExecuteOp(contractID string, args CanExecuteOp) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "CanExecuteOp",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMS contract
func (t MCMS) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// GetState exercises the GetState choice on this MCMS contract
func (t MCMS) GetState(contractID string, args GetState) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Main", "MCMS"),

		ContractID: contractID,
		Choice:     "GetState",

		Arguments: argsToMap(args),
	}
}

// MCMSEntrypointEvent is a Template type
type MCMSEntrypointEvent struct {
	Owner PARTY `json:"owner"`

	InstanceId TEXT `json:"instanceId"`

	FunctionName TEXT `json:"functionName"`

	OperationData TEXT `json:"operationData"`

	ContractIdsAsText []TEXT `json:"contractIdsAsText"`
}

// GetTemplateID returns the template ID for this template
func (t MCMSEntrypointEvent) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSEntrypointEvent")
}

// CreateCommand returns a CreateCommand for this template
func (t MCMSEntrypointEvent) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["functionName"] = string(t.FunctionName)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operationData"] = string(t.OperationData)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["contractIdsAsText"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ContractIdsAsText))
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

func (t MCMSEntrypointEvent) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

		Arguments: map[string]interface{}{},
	}
}

// ArchiveMCMSEntrypointEvent exercises the Archive_MCMSEntrypointEvent choice on this MCMSEntrypointEvent contract
func (t MCMSEntrypointEvent) ArchiveMCMSEntrypointEvent(contractID string, args ArchiveMCMSEntrypointEvent) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Counter", "MCMSEntrypointEvent"),

		ContractID: contractID,
		Choice:     "Archive_MCMSEntrypointEvent",

		Arguments: argsToMap(args),
	}
}

// MCMSReceiverView is a Record type
type MCMSReceiverView struct {
	Owner PARTY `json:"owner"`

	InstanceId TEXT `json:"instanceId"`
}

// ToMap converts MCMSReceiverView to a map for DAML arguments
func (t MCMSReceiverView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t MCMSReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MCMSReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSReceiverEntrypoint is a Record type
type MCMSReceiverEntrypoint struct {
	Caller PARTY `json:"caller"`

	FunctionName TEXT `json:"functionName"`

	OperationData TEXT `json:"operationData"`

	ContractIds []CONTRACT_ID `json:"contractIds"`
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

func (t MCMSReceiverEntrypoint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t MCMSReceiverGetInstanceId) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MCMSReceiverGetInstanceId) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSReceiverIssueAuthTicket is a Record type
type MCMSReceiverIssueAuthTicket struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts MCMSReceiverIssueAuthTicket to a map for DAML arguments
func (t MCMSReceiverIssueAuthTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t MCMSReceiverIssueAuthTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MCMSReceiverIssueAuthTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MCMSRegistry is a Template type
type MCMSRegistry struct {
	McmsOwner PARTY `json:"mcmsOwner"`

	Registrations GENMAP `json:"registrations"`
}

// GetTemplateID returns the template ID for this template
func (t MCMSRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t MCMSRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsOwner"] = t.McmsOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrations"] = func() interface{} {
		if t.Registrations == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.Registrations}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t MCMSRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MCMSRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MCMSRegistry

// Register exercises the Register choice on this MCMSRegistry contract
func (t MCMSRegistry) Register(contractID string, args Register) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry"),

		ContractID: contractID,
		Choice:     "Register",

		Arguments: argsToMap(args),
	}
}

// UpdateRegistration exercises the UpdateRegistration choice on this MCMSRegistry contract
func (t MCMSRegistry) UpdateRegistration(contractID string, args UpdateRegistration) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry"),

		ContractID: contractID,
		Choice:     "UpdateRegistration",

		Arguments: argsToMap(args),
	}
}

// Unregister exercises the Unregister choice on this MCMSRegistry contract
func (t MCMSRegistry) Unregister(contractID string, args Unregister) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry"),

		ContractID: contractID,
		Choice:     "Unregister",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MCMSRegistry contract
func (t MCMSRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// Lookup exercises the Lookup choice on this MCMSRegistry contract
func (t MCMSRegistry) Lookup(contractID string, args Lookup) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSRegistry", "MCMSRegistry"),

		ContractID: contractID,
		Choice:     "Lookup",

		Arguments: argsToMap(args),
	}
}

// MCMSState is a Record type
type MCMSState struct {
	Role Role `json:"role"`

	OpCount INT64 `json:"opCount"`

	PostOpCount INT64 `json:"postOpCount"`

	ValidUntil TIMESTAMP `json:"validUntil"`

	HasActiveRoot BOOL `json:"hasActiveRoot"`

	NumSigners INT64 `json:"numSigners"`
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

func (t MCMSState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MCMSState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MultisigConfig is a Record type
type MultisigConfig struct {
	Signers []SignerInfo `json:"signers"`

	GroupQuorums []INT64 `json:"groupQuorums"`

	GroupParents []INT64 `json:"groupParents"`
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

func (t MultisigConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MultisigConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Op is a Record type
type Op struct {
	ChainId INT64 `json:"chainId"`

	MultisigId TEXT `json:"multisigId"`

	Nonce INT64 `json:"nonce"`

	TargetInstanceId TEXT `json:"targetInstanceId"`

	FunctionName TEXT `json:"functionName"`

	OperationData TEXT `json:"operationData"`
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

func (t Op) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Op) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// RawSignature is a Record type
type RawSignature struct {
	PublicKey TEXT `json:"publicKey"`

	R TEXT `json:"r"`

	S TEXT `json:"s"`
}

// ToMap converts RawSignature to a map for DAML arguments
func (t RawSignature) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["publicKey"] = string(t.PublicKey)

	m["r"] = string(t.R)

	m["s"] = string(t.S)

	return m
}

func (t RawSignature) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *RawSignature) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Register is a Record type
type Register struct {
	InstanceId TEXT `json:"instanceId"`

	TargetCid CONTRACT_ID `json:"targetCid"`
}

// ToMap converts Register to a map for DAML arguments
func (t Register) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instanceId"] = string(t.InstanceId)

	m["targetCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TargetCid).(mapper); ok {
			return m.toMap()
		}
		return t.TargetCid
	}()

	return m
}

func (t Register) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Register) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Role is an enum type
type Role string

const (
	RoleProposer Role = "Proposer"
)

func (e Role) GetEnumConstructor() string { return string(e) }

func (e Role) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.Types", "Role")
}

func (e Role) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

func (e *Role) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

var _ ENUM = Role("")

// RootMetadata is a Record type
type RootMetadata struct {
	ChainId INT64 `json:"chainId"`

	MultisigId TEXT `json:"multisigId"`

	PreOpCount INT64 `json:"preOpCount"`

	PostOpCount INT64 `json:"postOpCount"`

	OverridePreviousRoot BOOL `json:"overridePreviousRoot"`
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

func (t RootMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *RootMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetConfig is a Record type
type SetConfig struct {
	NewSigners []SignerInfo `json:"newSigners"`

	NewGroupQuorums []INT64 `json:"newGroupQuorums"`

	NewGroupParents []INT64 `json:"newGroupParents"`

	ClearRoot BOOL `json:"clearRoot"`
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

func (t SetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *SetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetConfigParams is a Record type
type SetConfigParams struct {
	Signers []SignerInfo `json:"signers"`

	GroupQuorums []INT64 `json:"groupQuorums"`

	GroupParents []INT64 `json:"groupParents"`

	ClearRoot BOOL `json:"clearRoot"`
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

func (t SetConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *SetConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SetRoot is a Record type
type SetRoot struct {
	Submitter PARTY `json:"submitter"`

	NewRoot TEXT `json:"newRoot"`

	ValidUntil TIMESTAMP `json:"validUntil"`

	Metadata RootMetadata `json:"metadata"`

	MetadataProof []TEXT `json:"metadataProof"`

	Signatures []RawSignature `json:"signatures"`
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

func (t SetRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *SetRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// SignerInfo is a Record type
type SignerInfo struct {
	SignerAddress TEXT `json:"signerAddress"`

	SignerIndex INT64 `json:"signerIndex"`

	SignerGroup INT64 `json:"signerGroup"`
}

// ToMap converts SignerInfo to a map for DAML arguments
func (t SignerInfo) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["signerAddress"] = string(t.SignerAddress)

	m["signerIndex"] = int64(t.SignerIndex)

	m["signerGroup"] = int64(t.SignerGroup)

	return m
}

func (t SignerInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *SignerInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TargetAuthTicket is a Template type
type TargetAuthTicket struct {
	Owner PARTY `json:"owner"`

	InstanceId TEXT `json:"instanceId"`

	TargetCid CONTRACT_ID `json:"targetCid"`
}

// GetTemplateID returns the template ID for this template
func (t TargetAuthTicket) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSReceiver", "TargetAuthTicket")
}

// CreateCommand returns a CreateCommand for this template
func (t TargetAuthTicket) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["targetCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TargetCid).(mapper); ok {
			return m.toMap()
		}
		return t.TargetCid
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t TargetAuthTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TargetAuthTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for TargetAuthTicket

// Archive exercises the Archive choice on this TargetAuthTicket contract
func (t TargetAuthTicket) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSReceiver", "TargetAuthTicket"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// Consume exercises the Consume choice on this TargetAuthTicket contract
func (t TargetAuthTicket) Consume(contractID string, args Consume) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "MCMS.MCMSReceiver", "TargetAuthTicket"),

		ContractID: contractID,
		Choice:     "Consume",

		Arguments: argsToMap(args),
	}
}

// Unregister is a Record type
type Unregister struct {
	InstanceId TEXT `json:"instanceId"`
}

// ToMap converts Unregister to a map for DAML arguments
func (t Unregister) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t Unregister) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Unregister) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateRegistration is a Record type
type UpdateRegistration struct {
	InstanceId TEXT `json:"instanceId"`

	NewTargetCid CONTRACT_ID `json:"newTargetCid"`
}

// ToMap converts UpdateRegistration to a map for DAML arguments
func (t UpdateRegistration) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instanceId"] = string(t.InstanceId)

	m["newTargetCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.NewTargetCid).(mapper); ok {
			return m.toMap()
		}
		return t.NewTargetCid
	}()

	return m
}

func (t UpdateRegistration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *UpdateRegistration) UnmarshalJSON(data []byte) error {
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
