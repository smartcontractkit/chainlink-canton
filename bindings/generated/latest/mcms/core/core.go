package core

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	api "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
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
	PackageName = "mcms-core"
	PackageID   = "3b79485f0565363dd77fbb60c379a21fa978e1cbff8ed557ccb2927c4a522d40"
	SDKVersion  = "3.4.11"
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

// CanExecuteOp is a Record type
type CanExecuteOp struct {
	Submitter  types.PARTY `json:"submitter"`
	TargetRole api.Role    `json:"targetRole"`
	Op         api.Op      `json:"op"`
}

// ToMap converts CanExecuteOp to a map for DAML arguments
func (t CanExecuteOp) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

	m["op"] = model.NestedToDAMLValue(t.Op)

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

// ExecuteOp is a Record type
type ExecuteOp struct {
	TargetRole api.Role                         `json:"targetRole"`
	Submitter  types.PARTY                      `json:"submitter"`
	Op         api.Op                           `json:"op"`
	OpProof    []types.TEXT                     `json:"opProof"`
	TargetCids map[types.TEXT]types.CONTRACT_ID `json:"targetCids"`
}

// ToMap converts ExecuteOp to a map for DAML arguments
func (t ExecuteOp) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

	m["submitter"] = t.Submitter.ToMap()

	m["op"] = model.NestedToDAMLValue(t.Op)

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
	Submitter   types.PARTY                      `json:"submitter"`
	OpId        types.TEXT                       `json:"opId"`
	Calls       []api.TimelockCall               `json:"calls"`
	Predecessor types.TEXT                       `json:"predecessor" hex:"bytes16"`
	Salt        types.TEXT                       `json:"salt" hex:"bytes16"`
	TargetCids  map[types.TEXT]types.CONTRACT_ID `json:"targetCids"`
}

// ToMap converts ExecuteScheduledBatch to a map for DAML arguments
func (t ExecuteScheduledBatch) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["opId"] = string(t.OpId)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			res = append(res, model.NestedToDAMLValue(e))
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
	TargetRole api.Role    `json:"targetRole"`
}

// ToMap converts GetState to a map for DAML arguments
func (t GetState) ToMap() map[string]any {
	m := make(map[string]any)

	m["submitter"] = t.Submitter.ToMap()

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

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
	Owner              types.PARTY                    `json:"owner"`
	InstanceId         types.TEXT                     `json:"instanceId"`
	ChainId            types.INT64                    `json:"chainId"`
	Proposer           api.RoleState                  `json:"proposer"`
	Canceller          api.RoleState                  `json:"canceller"`
	Bypasser           api.RoleState                  `json:"bypasser"`
	MinDelay           types.RELTIME                  `json:"minDelay"`
	BlockedFunctions   []api.BlockedFunction          `json:"blockedFunctions"`
	TimelockTimestamps map[types.TEXT]types.TIMESTAMP `json:"timelockTimestamps"`
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
	args["proposer"] = model.NestedToDAMLValue(t.Proposer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["canceller"] = model.NestedToDAMLValue(t.Canceller)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["bypasser"] = model.NestedToDAMLValue(t.Bypasser)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minDelay"] = model.NestedToDAMLValue(t.MinDelay)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["blockedFunctions"] = func() []any {
		res := make([]any, 0, len(t.BlockedFunctions))
		for _, e := range t.BlockedFunctions {
			res = append(res, model.NestedToDAMLValue(e))
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
	args["proposer"] = model.NestedToDAMLValue(t.Proposer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["canceller"] = model.NestedToDAMLValue(t.Canceller)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["bypasser"] = model.NestedToDAMLValue(t.Bypasser)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minDelay"] = model.NestedToDAMLValue(t.MinDelay)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["blockedFunctions"] = func() []any {
		res := make([]any, 0, len(t.BlockedFunctions))
		for _, e := range t.BlockedFunctions {
			res = append(res, model.NestedToDAMLValue(e))
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

// PruneSeenHashes exercises the PruneSeenHashes choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) PruneSeenHashes(contractID string, args PruneSeenHashes) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "PruneSeenHashes",
		Arguments:  argsToMap(args),
	}
}

// PruneSeenHashesWithPackageID exercises the PruneSeenHashes choice using the provided package ID instead of package name
func (t MCMS) PruneSeenHashesWithPackageID(contractID string, packageID string, args PruneSeenHashes) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "PruneSeenHashes",
		Arguments:  argsToMap(args),
	}
}

// PruneTimelockTimestamps exercises the PruneTimelockTimestamps choice on this MCMS contract
// This method uses the package name in the template ID
func (t MCMS) PruneTimelockTimestamps(contractID string, args PruneTimelockTimestamps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "PruneTimelockTimestamps",
		Arguments:  argsToMap(args),
	}
}

// PruneTimelockTimestampsWithPackageID exercises the PruneTimelockTimestamps choice using the provided package ID instead of package name
func (t MCMS) PruneTimelockTimestampsWithPackageID(contractID string, packageID string, args PruneTimelockTimestamps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Main", "MCMS"),
		ContractID: contractID,
		Choice:     "PruneTimelockTimestamps",
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

// MCMSState is a Record type
type MCMSState struct {
	Role          api.Role        `json:"role"`
	OpCount       types.INT64     `json:"opCount"`
	PostOpCount   types.INT64     `json:"postOpCount"`
	ValidUntil    types.TIMESTAMP `json:"validUntil"`
	HasActiveRoot types.BOOL      `json:"hasActiveRoot"`
	NumSigners    types.INT64     `json:"numSigners"`
}

// ToMap converts MCMSState to a map for DAML arguments
func (t MCMSState) ToMap() map[string]any {
	m := make(map[string]any)

	m["role"] = model.NestedToDAMLValue(t.Role)

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

// PruneSeenHashes is a Record type
type PruneSeenHashes struct {
	TargetRole api.Role `json:"targetRole"`
}

// ToMap converts PruneSeenHashes to a map for DAML arguments
func (t PruneSeenHashes) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

	return m
}

func (t PruneSeenHashes) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PruneSeenHashes) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PruneSeenHashes to hex string (Canton MCMS format)
func (t PruneSeenHashes) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PruneSeenHashes from hex string (Canton MCMS format)
func (t *PruneSeenHashes) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PruneTimelockTimestamps is a Record type
type PruneTimelockTimestamps struct {
}

// ToMap converts PruneTimelockTimestamps to a map for DAML arguments
func (t PruneTimelockTimestamps) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t PruneTimelockTimestamps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PruneTimelockTimestamps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PruneTimelockTimestamps to hex string (Canton MCMS format)
func (t PruneTimelockTimestamps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PruneTimelockTimestamps from hex string (Canton MCMS format)
func (t *PruneTimelockTimestamps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetConfig is a Record type
type SetConfig struct {
	TargetRole      api.Role         `json:"targetRole"`
	NewSigners      []api.SignerInfo `json:"newSigners"`
	NewGroupQuorums []types.INT64    `json:"newGroupQuorums"`
	NewGroupParents []types.INT64    `json:"newGroupParents"`
	ClearRoot       types.BOOL       `json:"clearRoot"`
}

// ToMap converts SetConfig to a map for DAML arguments
func (t SetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

	m["newSigners"] = func() []any {
		res := make([]any, 0, len(t.NewSigners))
		for _, e := range t.NewSigners {
			res = append(res, model.NestedToDAMLValue(e))
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

// SetRoot is a Record type
type SetRoot struct {
	TargetRole    api.Role           `json:"targetRole"`
	Submitter     types.PARTY        `json:"submitter"`
	NewRoot       types.TEXT         `json:"newRoot" hex:"bytes"`
	ValidUntil    types.TIMESTAMP    `json:"validUntil"`
	Metadata      api.RootMetadata   `json:"metadata"`
	MetadataProof []types.TEXT       `json:"metadataProof"`
	Signatures    []api.RawSignature `json:"signatures"`
}

// ToMap converts SetRoot to a map for DAML arguments
func (t SetRoot) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetRole"] = model.NestedToDAMLValue(t.TargetRole)

	m["submitter"] = t.Submitter.ToMap()

	m["newRoot"] = string(t.NewRoot)

	m["validUntil"] = t.ValidUntil

	m["metadata"] = model.NestedToDAMLValue(t.Metadata)

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
			res = append(res, model.NestedToDAMLValue(e))
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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error)
	ExecuteOp(args ExecuteOp) (*bind.EncodedChoice, error)
	ExecuteScheduledBatch(args ExecuteScheduledBatch) (*bind.EncodedChoice, error)
	GetBlockedFunctions(args GetBlockedFunctions) (*bind.EncodedChoice, error)
	GetBlockedFunctionsCount(args GetBlockedFunctionsCount) (*bind.EncodedChoice, error)
	GetMinDelay(args GetMinDelay) (*bind.EncodedChoice, error)
	GetState(args GetState) (*bind.EncodedChoice, error)
	GetTimestamp(args GetTimestamp) (*bind.EncodedChoice, error)
	IsOperation(args IsOperation) (*bind.EncodedChoice, error)
	IsOperationDone(args IsOperationDone) (*bind.EncodedChoice, error)
	IsOperationPending(args IsOperationPending) (*bind.EncodedChoice, error)
	IsOperationReady(args IsOperationReady) (*bind.EncodedChoice, error)
	PruneSeenHashes(args PruneSeenHashes) (*bind.EncodedChoice, error)
	PruneTimelockTimestamps(args PruneTimelockTimestamps) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
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

// CanExecuteOp encodes parameters for the CanExecuteOp choice.
func (e *encoder) CanExecuteOp(args CanExecuteOp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CanExecuteOp", args)
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

// PruneSeenHashes encodes parameters for the PruneSeenHashes choice.
func (e *encoder) PruneSeenHashes(args PruneSeenHashes) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PruneSeenHashes", args)
}

// PruneTimelockTimestamps encodes parameters for the PruneTimelockTimestamps choice.
func (e *encoder) PruneTimelockTimestamps(args PruneTimelockTimestamps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PruneTimelockTimestamps", args)
}

// SetConfig encodes parameters for the SetConfig choice.
func (e *encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetRoot encodes parameters for the SetRoot choice.
func (e *encoder) SetRoot(args SetRoot) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetRoot", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
