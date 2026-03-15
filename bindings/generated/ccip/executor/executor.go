package executor

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
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
	PackageName = "ccip-executor"
	PackageID   = "a5119d605c8f4fc5ebb7e4920a8be9e0780a1d85b8d632502c599c02309c9c01"
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

// DynamicConfig is a Record type
type DynamicConfig struct {
	FeeAggregator         *types.PARTY `json:"feeAggregator" hex:"optional"`
	MinBlockConfirmations types.INT64  `json:"minBlockConfirmations"`
	CcvAllowlistEnabled   types.BOOL   `json:"ccvAllowlistEnabled"`
}

// ToMap converts DynamicConfig to a map for DAML arguments
func (t DynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	if t.FeeAggregator != nil {
		m["feeAggregator"] = map[string]any{
			"_type": "optional",
			"value": (*t.FeeAggregator).ToMap(),
		}
	} else {
		m["feeAggregator"] = map[string]any{
			"_type": "optional",
		}
	}

	m["minBlockConfirmations"] = int64(t.MinBlockConfirmations)

	m["ccvAllowlistEnabled"] = bool(t.CcvAllowlistEnabled)

	return m
}

func (t DynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DynamicConfig to hex string (Canton MCMS format)
func (t DynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DynamicConfig from hex string (Canton MCMS format)
func (t *DynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Executor is a Template type
type Executor struct {
	InstanceId         types.TEXT                  `json:"instanceId"`
	Owner              types.PARTY                 `json:"owner"`
	MaxCCVsPerMsg      types.INT64                 `json:"maxCCVsPerMsg"`
	DynamicConfig      DynamicConfig               `json:"dynamicConfig"`
	AllowedCCVs        []common.RawInstanceAddress `json:"allowedCCVs"`
	RemoteChainConfigs types.GENMAP                `json:"remoteChainConfigs"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t Executor) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t Executor) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Executor", "Executor")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t Executor) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["maxCCVsPerMsg"] = int64(t.MaxCCVsPerMsg)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["dynamicConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DynamicConfig).(mapper); ok {
			return m.toMap()
		}
		return t.DynamicConfig
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allowedCCVs"] = func() []any {
		res := make([]any, 0, len(t.AllowedCCVs))
		for _, e := range t.AllowedCCVs {
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
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t Executor) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["maxCCVsPerMsg"] = int64(t.MaxCCVsPerMsg)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["dynamicConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DynamicConfig).(mapper); ok {
			return m.toMap()
		}
		return t.DynamicConfig
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allowedCCVs"] = func() []any {
		res := make([]any, 0, len(t.AllowedCCVs))
		for _, e := range t.AllowedCCVs {
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
	args["remoteChainConfigs"] = func() any {
		if t.RemoteChainConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t Executor) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Executor) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Executor to hex string (Canton MCMS format)
func (t Executor) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Executor from hex string (Canton MCMS format)
func (t *Executor) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for Executor

// ExecutorApplyDestChainUpdates exercises the Executor_ApplyDestChainUpdates choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorApplyDestChainUpdates(contractID string, args ExecutorApplyDestChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_ApplyDestChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// ExecutorApplyDestChainUpdatesWithPackageID exercises the Executor_ApplyDestChainUpdates choice using the provided package ID instead of package name
func (t Executor) ExecutorApplyDestChainUpdatesWithPackageID(contractID string, packageID string, args ExecutorApplyDestChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_ApplyDestChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// ExecutorSetDynamicConfig exercises the Executor_SetDynamicConfig choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorSetDynamicConfig(contractID string, args ExecutorSetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// ExecutorSetDynamicConfigWithPackageID exercises the Executor_SetDynamicConfig choice using the provided package ID instead of package name
func (t Executor) ExecutorSetDynamicConfigWithPackageID(contractID string, packageID string, args ExecutorSetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// ExecutorApplyAllowedCCVUpdates exercises the Executor_ApplyAllowedCCVUpdates choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorApplyAllowedCCVUpdates(contractID string, args ExecutorApplyAllowedCCVUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_ApplyAllowedCCVUpdates",
		Arguments:  argsToMap(args),
	}
}

// ExecutorApplyAllowedCCVUpdatesWithPackageID exercises the Executor_ApplyAllowedCCVUpdates choice using the provided package ID instead of package name
func (t Executor) ExecutorApplyAllowedCCVUpdatesWithPackageID(contractID string, packageID string, args ExecutorApplyAllowedCCVUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_ApplyAllowedCCVUpdates",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetDestChains exercises the Executor_GetDestChains choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorGetDestChains(contractID string, args ExecutorGetDestChains) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetDestChains",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetDestChainsWithPackageID exercises the Executor_GetDestChains choice using the provided package ID instead of package name
func (t Executor) ExecutorGetDestChainsWithPackageID(contractID string, packageID string, args ExecutorGetDestChains) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetDestChains",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetDynamicConfig exercises the Executor_GetDynamicConfig choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorGetDynamicConfig(contractID string, args ExecutorGetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetDynamicConfigWithPackageID exercises the Executor_GetDynamicConfig choice using the provided package ID instead of package name
func (t Executor) ExecutorGetDynamicConfigWithPackageID(contractID string, packageID string, args ExecutorGetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetAllowedCCVs exercises the Executor_GetAllowedCCVs choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorGetAllowedCCVs(contractID string, args ExecutorGetAllowedCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetAllowedCCVs",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetAllowedCCVsWithPackageID exercises the Executor_GetAllowedCCVs choice using the provided package ID instead of package name
func (t Executor) ExecutorGetAllowedCCVsWithPackageID(contractID string, packageID string, args ExecutorGetAllowedCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetAllowedCCVs",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetMaxCCVsPerMessage exercises the Executor_GetMaxCCVsPerMessage choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorGetMaxCCVsPerMessage(contractID string, args ExecutorGetMaxCCVsPerMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetMaxCCVsPerMessage",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetMaxCCVsPerMessageWithPackageID exercises the Executor_GetMaxCCVsPerMessage choice using the provided package ID instead of package name
func (t Executor) ExecutorGetMaxCCVsPerMessageWithPackageID(contractID string, packageID string, args ExecutorGetMaxCCVsPerMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetMaxCCVsPerMessage",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetMinBlockConfirmations exercises the Executor_GetMinBlockConfirmations choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorGetMinBlockConfirmations(contractID string, args ExecutorGetMinBlockConfirmations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetMinBlockConfirmations",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetMinBlockConfirmationsWithPackageID exercises the Executor_GetMinBlockConfirmations choice using the provided package ID instead of package name
func (t Executor) ExecutorGetMinBlockConfirmationsWithPackageID(contractID string, packageID string, args ExecutorGetMinBlockConfirmations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_GetMinBlockConfirmations",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t Executor) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutorCalculateFee exercises the Executor_CalculateFee choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ExecutorCalculateFee(contractID string, args common.ExecutorCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// ExecutorCalculateFeeWithPackageID exercises the Executor_CalculateFee choice using the provided package ID instead of package name
func (t Executor) ExecutorCalculateFeeWithPackageID(contractID string, packageID string, args common.ExecutorCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "Executor_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Executor

var _ common.IIExecutor = (*Executor)(nil)

// ExecutorApplyAllowedCCVUpdates is a Record type
type ExecutorApplyAllowedCCVUpdates struct {
	CcvsToRemove        []common.RawInstanceAddress `json:"ccvsToRemove"`
	CcvsToAdd           []common.RawInstanceAddress `json:"ccvsToAdd"`
	CcvAllowlistEnabled types.BOOL                  `json:"ccvAllowlistEnabled"`
}

// ToMap converts ExecutorApplyAllowedCCVUpdates to a map for DAML arguments
func (t ExecutorApplyAllowedCCVUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvsToRemove"] = func() []any {
		res := make([]any, 0, len(t.CcvsToRemove))
		for _, e := range t.CcvsToRemove {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvsToAdd"] = func() []any {
		res := make([]any, 0, len(t.CcvsToAdd))
		for _, e := range t.CcvsToAdd {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvAllowlistEnabled"] = bool(t.CcvAllowlistEnabled)

	return m
}

func (t ExecutorApplyAllowedCCVUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorApplyAllowedCCVUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorApplyAllowedCCVUpdates to hex string (Canton MCMS format)
func (t ExecutorApplyAllowedCCVUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorApplyAllowedCCVUpdates from hex string (Canton MCMS format)
func (t *ExecutorApplyAllowedCCVUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorApplyDestChainUpdates is a Record type
type ExecutorApplyDestChainUpdates struct {
	DestChainSelectorsToRemove []types.NUMERIC         `json:"destChainSelectorsToRemove"`
	DestChainSelectorsToAdd    []RemoteChainConfigArgs `json:"destChainSelectorsToAdd"`
}

// ToMap converts ExecutorApplyDestChainUpdates to a map for DAML arguments
func (t ExecutorApplyDestChainUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.DestChainSelectorsToRemove))
		for _, e := range t.DestChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["destChainSelectorsToAdd"] = func() []any {
		res := make([]any, 0, len(t.DestChainSelectorsToAdd))
		for _, e := range t.DestChainSelectorsToAdd {
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

func (t ExecutorApplyDestChainUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorApplyDestChainUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorApplyDestChainUpdates to hex string (Canton MCMS format)
func (t ExecutorApplyDestChainUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorApplyDestChainUpdates from hex string (Canton MCMS format)
func (t *ExecutorApplyDestChainUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorCalculateFee2 is a Record type
type ExecutorCalculateFee2 struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	ExecutorArgs      types.TEXT        `json:"executorArgs"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts ExecutorCalculateFee2 to a map for DAML arguments
func (t ExecutorCalculateFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["executorArgs"] = string(t.ExecutorArgs)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorCalculateFee2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorCalculateFee2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorCalculateFee2 to hex string (Canton MCMS format)
func (t ExecutorCalculateFee2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorCalculateFee2 from hex string (Canton MCMS format)
func (t *ExecutorCalculateFee2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetAllowedCCVs is a Record type
type ExecutorGetAllowedCCVs struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutorGetAllowedCCVs to a map for DAML arguments
func (t ExecutorGetAllowedCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetAllowedCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetAllowedCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetAllowedCCVs to hex string (Canton MCMS format)
func (t ExecutorGetAllowedCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetAllowedCCVs from hex string (Canton MCMS format)
func (t *ExecutorGetAllowedCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetAllowedCCVsMCMSParams is ExecutorGetAllowedCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ExecutorGetAllowedCCVsMCMSParams struct {
}

// MarshalHex encodes ExecutorGetAllowedCCVsMCMSParams to hex string for MCMS operationData.
func (t ExecutorGetAllowedCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetAllowedCCVsMCMSParams from hex string.
func (t *ExecutorGetAllowedCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetDestChains is a Record type
type ExecutorGetDestChains struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutorGetDestChains to a map for DAML arguments
func (t ExecutorGetDestChains) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetDestChains) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetDestChains) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetDestChains to hex string (Canton MCMS format)
func (t ExecutorGetDestChains) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetDestChains from hex string (Canton MCMS format)
func (t *ExecutorGetDestChains) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetDestChainsMCMSParams is ExecutorGetDestChains without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ExecutorGetDestChainsMCMSParams struct {
}

// MarshalHex encodes ExecutorGetDestChainsMCMSParams to hex string for MCMS operationData.
func (t ExecutorGetDestChainsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetDestChainsMCMSParams from hex string.
func (t *ExecutorGetDestChainsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetDynamicConfig is a Record type
type ExecutorGetDynamicConfig struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutorGetDynamicConfig to a map for DAML arguments
func (t ExecutorGetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetDynamicConfig to hex string (Canton MCMS format)
func (t ExecutorGetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetDynamicConfig from hex string (Canton MCMS format)
func (t *ExecutorGetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetDynamicConfigMCMSParams is ExecutorGetDynamicConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ExecutorGetDynamicConfigMCMSParams struct {
}

// MarshalHex encodes ExecutorGetDynamicConfigMCMSParams to hex string for MCMS operationData.
func (t ExecutorGetDynamicConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetDynamicConfigMCMSParams from hex string.
func (t *ExecutorGetDynamicConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetMaxCCVsPerMessage is a Record type
type ExecutorGetMaxCCVsPerMessage struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutorGetMaxCCVsPerMessage to a map for DAML arguments
func (t ExecutorGetMaxCCVsPerMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetMaxCCVsPerMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetMaxCCVsPerMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetMaxCCVsPerMessage to hex string (Canton MCMS format)
func (t ExecutorGetMaxCCVsPerMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetMaxCCVsPerMessage from hex string (Canton MCMS format)
func (t *ExecutorGetMaxCCVsPerMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetMaxCCVsPerMessageMCMSParams is ExecutorGetMaxCCVsPerMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ExecutorGetMaxCCVsPerMessageMCMSParams struct {
}

// MarshalHex encodes ExecutorGetMaxCCVsPerMessageMCMSParams to hex string for MCMS operationData.
func (t ExecutorGetMaxCCVsPerMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetMaxCCVsPerMessageMCMSParams from hex string.
func (t *ExecutorGetMaxCCVsPerMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetMinBlockConfirmations is a Record type
type ExecutorGetMinBlockConfirmations struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts ExecutorGetMinBlockConfirmations to a map for DAML arguments
func (t ExecutorGetMinBlockConfirmations) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t ExecutorGetMinBlockConfirmations) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorGetMinBlockConfirmations) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorGetMinBlockConfirmations to hex string (Canton MCMS format)
func (t ExecutorGetMinBlockConfirmations) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetMinBlockConfirmations from hex string (Canton MCMS format)
func (t *ExecutorGetMinBlockConfirmations) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorGetMinBlockConfirmationsMCMSParams is ExecutorGetMinBlockConfirmations without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type ExecutorGetMinBlockConfirmationsMCMSParams struct {
}

// MarshalHex encodes ExecutorGetMinBlockConfirmationsMCMSParams to hex string for MCMS operationData.
func (t ExecutorGetMinBlockConfirmationsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorGetMinBlockConfirmationsMCMSParams from hex string.
func (t *ExecutorGetMinBlockConfirmationsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutorSetDynamicConfig is a Record type
type ExecutorSetDynamicConfig struct {
	NewDynamicConfig DynamicConfig `json:"newDynamicConfig"`
}

// ToMap converts ExecutorSetDynamicConfig to a map for DAML arguments
func (t ExecutorSetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDynamicConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewDynamicConfig).(mapper); ok {
			return m.toMap()
		}
		return t.NewDynamicConfig
	}()

	return m
}

func (t ExecutorSetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutorSetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutorSetDynamicConfig to hex string (Canton MCMS format)
func (t ExecutorSetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutorSetDynamicConfig from hex string (Canton MCMS format)
func (t *ExecutorSetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoteChainConfig is a Record type
type RemoteChainConfig struct {
	FeeUSDCents types.NUMERIC `json:"feeUSDCents"`
	Enabled     types.BOOL    `json:"enabled"`
}

// ToMap converts RemoteChainConfig to a map for DAML arguments
func (t RemoteChainConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeUSDCents"] = t.FeeUSDCents

	m["enabled"] = bool(t.Enabled)

	return m
}

func (t RemoteChainConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoteChainConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoteChainConfig to hex string (Canton MCMS format)
func (t RemoteChainConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoteChainConfig from hex string (Canton MCMS format)
func (t *RemoteChainConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoteChainConfigArgs is a Record type
type RemoteChainConfigArgs struct {
	DestChainSelector types.NUMERIC     `json:"destChainSelector"`
	Config            RemoteChainConfig `json:"config"`
}

// ToMap converts RemoteChainConfigArgs to a map for DAML arguments
func (t RemoteChainConfigArgs) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["config"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Config).(mapper); ok {
			return m.toMap()
		}
		return t.Config
	}()

	return m
}

func (t RemoteChainConfigArgs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoteChainConfigArgs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoteChainConfigArgs to hex string (Canton MCMS format)
func (t RemoteChainConfigArgs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoteChainConfigArgs from hex string (Canton MCMS format)
func (t *RemoteChainConfigArgs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ExecutorApplyAllowedCCVUpdates(args ExecutorApplyAllowedCCVUpdates) (*bind.EncodedChoice, error)
	ExecutorApplyDestChainUpdates(args ExecutorApplyDestChainUpdates) (*bind.EncodedChoice, error)
	ExecutorGetAllowedCCVs(args ExecutorGetAllowedCCVs) (*bind.EncodedChoice, error)
	ExecutorGetAllowedCCVsMCMSParams(args ExecutorGetAllowedCCVsMCMSParams) (*bind.EncodedChoice, error)
	ExecutorGetDestChains(args ExecutorGetDestChains) (*bind.EncodedChoice, error)
	ExecutorGetDestChainsMCMSParams(args ExecutorGetDestChainsMCMSParams) (*bind.EncodedChoice, error)
	ExecutorGetDynamicConfig(args ExecutorGetDynamicConfig) (*bind.EncodedChoice, error)
	ExecutorGetDynamicConfigMCMSParams(args ExecutorGetDynamicConfigMCMSParams) (*bind.EncodedChoice, error)
	ExecutorGetMaxCCVsPerMessage(args ExecutorGetMaxCCVsPerMessage) (*bind.EncodedChoice, error)
	ExecutorGetMaxCCVsPerMessageMCMSParams(args ExecutorGetMaxCCVsPerMessageMCMSParams) (*bind.EncodedChoice, error)
	ExecutorGetMinBlockConfirmations(args ExecutorGetMinBlockConfirmations) (*bind.EncodedChoice, error)
	ExecutorGetMinBlockConfirmationsMCMSParams(args ExecutorGetMinBlockConfirmationsMCMSParams) (*bind.EncodedChoice, error)
	ExecutorSetDynamicConfig(args ExecutorSetDynamicConfig) (*bind.EncodedChoice, error)
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

// ExecutorApplyAllowedCCVUpdates encodes parameters for the ExecutorApplyAllowedCCVUpdates choice.
func (e *encoder) ExecutorApplyAllowedCCVUpdates(args ExecutorApplyAllowedCCVUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorApplyAllowedCCVUpdates", args)
}

// ExecutorApplyDestChainUpdates encodes parameters for the ExecutorApplyDestChainUpdates choice.
func (e *encoder) ExecutorApplyDestChainUpdates(args ExecutorApplyDestChainUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorApplyDestChainUpdates", args)
}

// ExecutorGetAllowedCCVs encodes parameters for the ExecutorGetAllowedCCVs choice.
func (e *encoder) ExecutorGetAllowedCCVs(args ExecutorGetAllowedCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetAllowedCCVs", args)
}

// ExecutorGetAllowedCCVsMCMSParams encodes MCMS parameters (without Caller) for the ExecutorGetAllowedCCVs choice.
func (e *encoder) ExecutorGetAllowedCCVsMCMSParams(args ExecutorGetAllowedCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetAllowedCCVs", args)
}

// ExecutorGetDestChains encodes parameters for the ExecutorGetDestChains choice.
func (e *encoder) ExecutorGetDestChains(args ExecutorGetDestChains) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetDestChains", args)
}

// ExecutorGetDestChainsMCMSParams encodes MCMS parameters (without Caller) for the ExecutorGetDestChains choice.
func (e *encoder) ExecutorGetDestChainsMCMSParams(args ExecutorGetDestChainsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetDestChains", args)
}

// ExecutorGetDynamicConfig encodes parameters for the ExecutorGetDynamicConfig choice.
func (e *encoder) ExecutorGetDynamicConfig(args ExecutorGetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetDynamicConfig", args)
}

// ExecutorGetDynamicConfigMCMSParams encodes MCMS parameters (without Caller) for the ExecutorGetDynamicConfig choice.
func (e *encoder) ExecutorGetDynamicConfigMCMSParams(args ExecutorGetDynamicConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetDynamicConfig", args)
}

// ExecutorGetMaxCCVsPerMessage encodes parameters for the ExecutorGetMaxCCVsPerMessage choice.
func (e *encoder) ExecutorGetMaxCCVsPerMessage(args ExecutorGetMaxCCVsPerMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetMaxCCVsPerMessage", args)
}

// ExecutorGetMaxCCVsPerMessageMCMSParams encodes MCMS parameters (without Caller) for the ExecutorGetMaxCCVsPerMessage choice.
func (e *encoder) ExecutorGetMaxCCVsPerMessageMCMSParams(args ExecutorGetMaxCCVsPerMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetMaxCCVsPerMessage", args)
}

// ExecutorGetMinBlockConfirmations encodes parameters for the ExecutorGetMinBlockConfirmations choice.
func (e *encoder) ExecutorGetMinBlockConfirmations(args ExecutorGetMinBlockConfirmations) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetMinBlockConfirmations", args)
}

// ExecutorGetMinBlockConfirmationsMCMSParams encodes MCMS parameters (without Caller) for the ExecutorGetMinBlockConfirmations choice.
func (e *encoder) ExecutorGetMinBlockConfirmationsMCMSParams(args ExecutorGetMinBlockConfirmationsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorGetMinBlockConfirmations", args)
}

// ExecutorSetDynamicConfig encodes parameters for the ExecutorSetDynamicConfig choice.
func (e *encoder) ExecutorSetDynamicConfig(args ExecutorSetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutorSetDynamicConfig", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
