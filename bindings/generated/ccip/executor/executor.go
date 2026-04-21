package executor

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
	PackageName = "ccip-executor"
	PackageID   = "a348c31180506a0f0c4b8060132a6bd86b85919aae22f973dfb98866633e30eb"
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

// ApplyAllowedCCVUpdates is a Record type
type ApplyAllowedCCVUpdates struct {
	CcvsToRemove        []mcms.RawInstanceAddress `json:"ccvsToRemove"`
	CcvsToAdd           []mcms.RawInstanceAddress `json:"ccvsToAdd"`
	CcvAllowlistEnabled types.BOOL                `json:"ccvAllowlistEnabled"`
}

// ToMap converts ApplyAllowedCCVUpdates to a map for DAML arguments
func (t ApplyAllowedCCVUpdates) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvsToRemove"] = func() []any {
		res := make([]any, 0, len(t.CcvsToRemove))
		for _, e := range t.CcvsToRemove {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvsToAdd"] = func() []any {
		res := make([]any, 0, len(t.CcvsToAdd))
		for _, e := range t.CcvsToAdd {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvAllowlistEnabled"] = bool(t.CcvAllowlistEnabled)

	return m
}

func (t ApplyAllowedCCVUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyAllowedCCVUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyAllowedCCVUpdates to hex string (Canton MCMS format)
func (t ApplyAllowedCCVUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyAllowedCCVUpdates from hex string (Canton MCMS format)
func (t *ApplyAllowedCCVUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyAllowedCCVUpdatesParams is a Record type
type ApplyAllowedCCVUpdatesParams struct {
	CcvsToRemove        []mcms.RawInstanceAddress `json:"ccvsToRemove"`
	CcvsToAdd           []mcms.RawInstanceAddress `json:"ccvsToAdd"`
	CcvAllowlistEnabled types.BOOL                `json:"ccvAllowlistEnabled"`
}

// ToMap converts ApplyAllowedCCVUpdatesParams to a map for DAML arguments
func (t ApplyAllowedCCVUpdatesParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvsToRemove"] = func() []any {
		res := make([]any, 0, len(t.CcvsToRemove))
		for _, e := range t.CcvsToRemove {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvsToAdd"] = func() []any {
		res := make([]any, 0, len(t.CcvsToAdd))
		for _, e := range t.CcvsToAdd {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["ccvAllowlistEnabled"] = bool(t.CcvAllowlistEnabled)

	return m
}

func (t ApplyAllowedCCVUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyAllowedCCVUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyAllowedCCVUpdatesParams to hex string (Canton MCMS format)
func (t ApplyAllowedCCVUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyAllowedCCVUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyAllowedCCVUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyDestChainUpdates is a Record type
type ApplyDestChainUpdates struct {
	DestChainSelectorsToRemove []types.NUMERIC         `json:"destChainSelectorsToRemove"`
	DestChainSelectorsToAdd    []RemoteChainConfigArgs `json:"destChainSelectorsToAdd"`
}

// ToMap converts ApplyDestChainUpdates to a map for DAML arguments
func (t ApplyDestChainUpdates) ToMap() map[string]any {
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
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t ApplyDestChainUpdates) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainUpdates) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainUpdates to hex string (Canton MCMS format)
func (t ApplyDestChainUpdates) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainUpdates from hex string (Canton MCMS format)
func (t *ApplyDestChainUpdates) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApplyDestChainUpdatesParams is a Record type
type ApplyDestChainUpdatesParams struct {
	DestChainSelectorsToRemove []types.NUMERIC         `json:"destChainSelectorsToRemove"`
	DestChainSelectorsToAdd    []RemoteChainConfigArgs `json:"destChainSelectorsToAdd"`
}

// ToMap converts ApplyDestChainUpdatesParams to a map for DAML arguments
func (t ApplyDestChainUpdatesParams) ToMap() map[string]any {
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
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t ApplyDestChainUpdatesParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApplyDestChainUpdatesParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApplyDestChainUpdatesParams to hex string (Canton MCMS format)
func (t ApplyDestChainUpdatesParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApplyDestChainUpdatesParams from hex string (Canton MCMS format)
func (t *ApplyDestChainUpdatesParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFee is a Record type
type CalculateFee struct {
	SendingMessageCid types.CONTRACT_ID  `json:"sendingMessageCid"`
	ExecutorArgs      types.TEXT         `json:"executorArgs"`
	ExtraContext      common.CCIPContext `json:"extraContext"`
	Caller            types.PARTY        `json:"caller"`
}

// ToMap converts CalculateFee to a map for DAML arguments
func (t CalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.ToMap()
		}
		return t.SendingMessageCid
	}()

	m["executorArgs"] = string(t.ExecutorArgs)

	m["extraContext"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraContext
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CalculateFee to hex string (Canton MCMS format)
func (t CalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFee from hex string (Canton MCMS format)
func (t *CalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CalculateFeeMCMSParams is CalculateFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CalculateFeeMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID  `json:"sendingMessageCid"`
	ExecutorArgs      types.TEXT         `json:"executorArgs"`
	ExtraContext      common.CCIPContext `json:"extraContext"`
}

// MarshalHex encodes CalculateFeeMCMSParams to hex string for MCMS operationData.
func (t CalculateFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CalculateFeeMCMSParams from hex string.
func (t *CalculateFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DynamicConfig is a Record type
type DynamicConfig struct {
	FeeAggregator         *types.PARTY          `json:"feeAggregator" hex:"optional"`
	AllowedFinalityConfig common.FinalityConfig `json:"allowedFinalityConfig"`
	CcvAllowlistEnabled   types.BOOL            `json:"ccvAllowlistEnabled"`
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

	m["allowedFinalityConfig"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.AllowedFinalityConfig).(mapper); ok {
			return m.ToMap()
		}
		return t.AllowedFinalityConfig
	}()

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
	InstanceId         types.TEXT                `json:"instanceId"`
	Owner              types.PARTY               `json:"owner"`
	MaxCCVsPerMsg      types.INT64               `json:"maxCCVsPerMsg"`
	DynamicConfig      DynamicConfig             `json:"dynamicConfig"`
	AllowedCCVs        []mcms.RawInstanceAddress `json:"allowedCCVs"`
	RemoteChainConfigs types.GENMAP              `json:"remoteChainConfigs"`
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
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.DynamicConfig).(mapper); ok {
			return m.ToMap()
		}
		return t.DynamicConfig
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allowedCCVs"] = func() []any {
		res := make([]any, 0, len(t.AllowedCCVs))
		for _, e := range t.AllowedCCVs {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
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
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.DynamicConfig).(mapper); ok {
			return m.ToMap()
		}
		return t.DynamicConfig
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allowedCCVs"] = func() []any {
		res := make([]any, 0, len(t.AllowedCCVs))
		for _, e := range t.AllowedCCVs {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
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

// CalculateFee exercises the CalculateFee choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) CalculateFee(contractID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CalculateFeeWithPackageID exercises the CalculateFee choice using the provided package ID instead of package name
func (t Executor) CalculateFeeWithPackageID(contractID string, packageID string, args CalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// GetFee exercises the GetFee choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t Executor) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this Executor contract via the IIExecutor interface
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

// GetDestChains exercises the GetDestChains choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetDestChains(contractID string, args GetDestChains) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetDestChains",
		Arguments:  argsToMap(args),
	}
}

// GetDestChainsWithPackageID exercises the GetDestChains choice using the provided package ID instead of package name
func (t Executor) GetDestChainsWithPackageID(contractID string, packageID string, args GetDestChains) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetDestChains",
		Arguments:  argsToMap(args),
	}
}

// GetDynamicConfig exercises the GetDynamicConfig choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetDynamicConfig(contractID string, args GetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// GetDynamicConfigWithPackageID exercises the GetDynamicConfig choice using the provided package ID instead of package name
func (t Executor) GetDynamicConfigWithPackageID(contractID string, packageID string, args GetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// GetAllowedCCVs exercises the GetAllowedCCVs choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetAllowedCCVs(contractID string, args GetAllowedCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetAllowedCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetAllowedCCVsWithPackageID exercises the GetAllowedCCVs choice using the provided package ID instead of package name
func (t Executor) GetAllowedCCVsWithPackageID(contractID string, packageID string, args GetAllowedCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetAllowedCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetMaxCCVsPerMessage exercises the GetMaxCCVsPerMessage choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetMaxCCVsPerMessage(contractID string, args GetMaxCCVsPerMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetMaxCCVsPerMessage",
		Arguments:  argsToMap(args),
	}
}

// GetMaxCCVsPerMessageWithPackageID exercises the GetMaxCCVsPerMessage choice using the provided package ID instead of package name
func (t Executor) GetMaxCCVsPerMessageWithPackageID(contractID string, packageID string, args GetMaxCCVsPerMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetMaxCCVsPerMessage",
		Arguments:  argsToMap(args),
	}
}

// GetAllowedFinalityConfig exercises the GetAllowedFinalityConfig choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) GetAllowedFinalityConfig(contractID string, args GetAllowedFinalityConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetAllowedFinalityConfig",
		Arguments:  argsToMap(args),
	}
}

// GetAllowedFinalityConfigWithPackageID exercises the GetAllowedFinalityConfig choice using the provided package ID instead of package name
func (t Executor) GetAllowedFinalityConfigWithPackageID(contractID string, packageID string, args GetAllowedFinalityConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "GetAllowedFinalityConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainUpdates exercises the ApplyDestChainUpdates choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ApplyDestChainUpdates(contractID string, args ApplyDestChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "ApplyDestChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyDestChainUpdatesWithPackageID exercises the ApplyDestChainUpdates choice using the provided package ID instead of package name
func (t Executor) ApplyDestChainUpdatesWithPackageID(contractID string, packageID string, args ApplyDestChainUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "ApplyDestChainUpdates",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfig exercises the SetDynamicConfig choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) SetDynamicConfig(contractID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// SetDynamicConfigWithPackageID exercises the SetDynamicConfig choice using the provided package ID instead of package name
func (t Executor) SetDynamicConfigWithPackageID(contractID string, packageID string, args SetDynamicConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "SetDynamicConfig",
		Arguments:  argsToMap(args),
	}
}

// ApplyAllowedCCVUpdates exercises the ApplyAllowedCCVUpdates choice on this Executor contract
// This method uses the package name in the template ID
func (t Executor) ApplyAllowedCCVUpdates(contractID string, args ApplyAllowedCCVUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "ApplyAllowedCCVUpdates",
		Arguments:  argsToMap(args),
	}
}

// ApplyAllowedCCVUpdatesWithPackageID exercises the ApplyAllowedCCVUpdates choice using the provided package ID instead of package name
func (t Executor) ApplyAllowedCCVUpdatesWithPackageID(contractID string, packageID string, args ApplyAllowedCCVUpdates) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "Executor"),
		ContractID: contractID,
		Choice:     "ApplyAllowedCCVUpdates",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this Executor contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Executor) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t Executor) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// ExecutorCalculateFee exercises the Executor_CalculateFee choice on this Executor contract via the IIExecutor interface
// This method uses the package name in the template ID
func (t Executor) ExecutorCalculateFee(contractID string, args common.ExecutorCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "IExecutor"),
		ContractID: contractID,
		Choice:     "Executor_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// ExecutorCalculateFeeWithPackageID exercises the Executor_CalculateFee choice using the provided package ID instead of package name
func (t Executor) ExecutorCalculateFeeWithPackageID(contractID string, packageID string, args common.ExecutorCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "IExecutor"),
		ContractID: contractID,
		Choice:     "Executor_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetFee exercises the Executor_GetFee choice on this Executor contract via the IIExecutor interface
// This method uses the package name in the template ID
func (t Executor) ExecutorGetFee(contractID string, args common.ExecutorGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Executor", "IExecutor"),
		ContractID: contractID,
		Choice:     "Executor_GetFee",
		Arguments:  argsToMap(args),
	}
}

// ExecutorGetFeeWithPackageID exercises the Executor_GetFee choice using the provided package ID instead of package name
func (t Executor) ExecutorGetFeeWithPackageID(contractID string, packageID string, args common.ExecutorGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Executor", "IExecutor"),
		ContractID: contractID,
		Choice:     "Executor_GetFee",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Executor

var _ mcms.IMCMSReceiver = (*Executor)(nil)

var _ common.IIExecutor = (*Executor)(nil)

// GetAllowedCCVs is a Record type
type GetAllowedCCVs struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetAllowedCCVs to a map for DAML arguments
func (t GetAllowedCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetAllowedCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetAllowedCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetAllowedCCVs to hex string (Canton MCMS format)
func (t GetAllowedCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetAllowedCCVs from hex string (Canton MCMS format)
func (t *GetAllowedCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetAllowedCCVsMCMSParams is GetAllowedCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetAllowedCCVsMCMSParams struct {
}

// MarshalHex encodes GetAllowedCCVsMCMSParams to hex string for MCMS operationData.
func (t GetAllowedCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetAllowedCCVsMCMSParams from hex string.
func (t *GetAllowedCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetAllowedFinalityConfig is a Record type
type GetAllowedFinalityConfig struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetAllowedFinalityConfig to a map for DAML arguments
func (t GetAllowedFinalityConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetAllowedFinalityConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetAllowedFinalityConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetAllowedFinalityConfig to hex string (Canton MCMS format)
func (t GetAllowedFinalityConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetAllowedFinalityConfig from hex string (Canton MCMS format)
func (t *GetAllowedFinalityConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetAllowedFinalityConfigMCMSParams is GetAllowedFinalityConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetAllowedFinalityConfigMCMSParams struct {
}

// MarshalHex encodes GetAllowedFinalityConfigMCMSParams to hex string for MCMS operationData.
func (t GetAllowedFinalityConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetAllowedFinalityConfigMCMSParams from hex string.
func (t *GetAllowedFinalityConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestChains is a Record type
type GetDestChains struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetDestChains to a map for DAML arguments
func (t GetDestChains) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDestChains) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDestChains) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDestChains to hex string (Canton MCMS format)
func (t GetDestChains) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChains from hex string (Canton MCMS format)
func (t *GetDestChains) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDestChainsMCMSParams is GetDestChains without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDestChainsMCMSParams struct {
}

// MarshalHex encodes GetDestChainsMCMSParams to hex string for MCMS operationData.
func (t GetDestChainsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDestChainsMCMSParams from hex string.
func (t *GetDestChainsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDynamicConfig is a Record type
type GetDynamicConfig struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetDynamicConfig to a map for DAML arguments
func (t GetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetDynamicConfig to hex string (Canton MCMS format)
func (t GetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDynamicConfig from hex string (Canton MCMS format)
func (t *GetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetDynamicConfigMCMSParams is GetDynamicConfig without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetDynamicConfigMCMSParams struct {
}

// MarshalHex encodes GetDynamicConfigMCMSParams to hex string for MCMS operationData.
func (t GetDynamicConfigMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetDynamicConfigMCMSParams from hex string.
func (t *GetDynamicConfigMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee is a Record type
type GetFee struct {
	DestChainSelector types.NUMERIC             `json:"destChainSelector"`
	RequiredCCVs      []mcms.RawInstanceAddress `json:"requiredCCVs"`
	Caller            types.PARTY               `json:"caller"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ ToMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.ToMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee to hex string (Canton MCMS format)
func (t GetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee from hex string (Canton MCMS format)
func (t *GetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeMCMSParams is GetFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFeeMCMSParams struct {
	DestChainSelector types.NUMERIC             `json:"destChainSelector"`
	RequiredCCVs      []mcms.RawInstanceAddress `json:"requiredCCVs"`
}

// MarshalHex encodes GetFeeMCMSParams to hex string for MCMS operationData.
func (t GetFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeMCMSParams from hex string.
func (t *GetFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMaxCCVsPerMessage is a Record type
type GetMaxCCVsPerMessage struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetMaxCCVsPerMessage to a map for DAML arguments
func (t GetMaxCCVsPerMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetMaxCCVsPerMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetMaxCCVsPerMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetMaxCCVsPerMessage to hex string (Canton MCMS format)
func (t GetMaxCCVsPerMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMaxCCVsPerMessage from hex string (Canton MCMS format)
func (t *GetMaxCCVsPerMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetMaxCCVsPerMessageMCMSParams is GetMaxCCVsPerMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetMaxCCVsPerMessageMCMSParams struct {
}

// MarshalHex encodes GetMaxCCVsPerMessageMCMSParams to hex string for MCMS operationData.
func (t GetMaxCCVsPerMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetMaxCCVsPerMessageMCMSParams from hex string.
func (t *GetMaxCCVsPerMessageMCMSParams) UnmarshalHex(data string) error {
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
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Config).(mapper); ok {
			return m.ToMap()
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

// SetDynamicConfig is a Record type
type SetDynamicConfig struct {
	NewDynamicConfig DynamicConfig `json:"newDynamicConfig"`
}

// ToMap converts SetDynamicConfig to a map for DAML arguments
func (t SetDynamicConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDynamicConfig"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.NewDynamicConfig).(mapper); ok {
			return m.ToMap()
		}
		return t.NewDynamicConfig
	}()

	return m
}

func (t SetDynamicConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfig to hex string (Canton MCMS format)
func (t SetDynamicConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfig from hex string (Canton MCMS format)
func (t *SetDynamicConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDynamicConfigParams is a Record type
type SetDynamicConfigParams struct {
	NewDynamicConfig DynamicConfig `json:"newDynamicConfig"`
}

// ToMap converts SetDynamicConfigParams to a map for DAML arguments
func (t SetDynamicConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDynamicConfig"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.NewDynamicConfig).(mapper); ok {
			return m.ToMap()
		}
		return t.NewDynamicConfig
	}()

	return m
}

func (t SetDynamicConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDynamicConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDynamicConfigParams to hex string (Canton MCMS format)
func (t SetDynamicConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDynamicConfigParams from hex string (Canton MCMS format)
func (t *SetDynamicConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ApplyAllowedCCVUpdates(args ApplyAllowedCCVUpdates) (*bind.EncodedChoice, error)
	ApplyAllowedCCVUpdatesParams(args ApplyAllowedCCVUpdatesParams) (*bind.EncodedChoice, error)
	ApplyDestChainUpdates(args ApplyDestChainUpdates) (*bind.EncodedChoice, error)
	ApplyDestChainUpdatesParams(args ApplyDestChainUpdatesParams) (*bind.EncodedChoice, error)
	CalculateFee(args CalculateFee) (*bind.EncodedChoice, error)
	CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	GetAllowedCCVs(args GetAllowedCCVs) (*bind.EncodedChoice, error)
	GetAllowedCCVsMCMSParams(args GetAllowedCCVsMCMSParams) (*bind.EncodedChoice, error)
	GetAllowedFinalityConfig(args GetAllowedFinalityConfig) (*bind.EncodedChoice, error)
	GetAllowedFinalityConfigMCMSParams(args GetAllowedFinalityConfigMCMSParams) (*bind.EncodedChoice, error)
	GetDestChains(args GetDestChains) (*bind.EncodedChoice, error)
	GetDestChainsMCMSParams(args GetDestChainsMCMSParams) (*bind.EncodedChoice, error)
	GetDynamicConfig(args GetDynamicConfig) (*bind.EncodedChoice, error)
	GetDynamicConfigMCMSParams(args GetDynamicConfigMCMSParams) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error)
	GetMaxCCVsPerMessage(args GetMaxCCVsPerMessage) (*bind.EncodedChoice, error)
	GetMaxCCVsPerMessageMCMSParams(args GetMaxCCVsPerMessageMCMSParams) (*bind.EncodedChoice, error)
	SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error)
	SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error)
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

// ApplyAllowedCCVUpdates encodes parameters for the ApplyAllowedCCVUpdates choice.
func (e *encoder) ApplyAllowedCCVUpdates(args ApplyAllowedCCVUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyAllowedCCVUpdates", args)
}

// ApplyAllowedCCVUpdatesParams encodes parameters for the ApplyAllowedCCVUpdates choice.
func (e *encoder) ApplyAllowedCCVUpdatesParams(args ApplyAllowedCCVUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyAllowedCCVUpdates", args)
}

// ApplyDestChainUpdates encodes parameters for the ApplyDestChainUpdates choice.
func (e *encoder) ApplyDestChainUpdates(args ApplyDestChainUpdates) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyDestChainUpdates", args)
}

// ApplyDestChainUpdatesParams encodes parameters for the ApplyDestChainUpdates choice.
func (e *encoder) ApplyDestChainUpdatesParams(args ApplyDestChainUpdatesParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApplyDestChainUpdates", args)
}

// CalculateFee encodes parameters for the CalculateFee choice.
func (e *encoder) CalculateFee(args CalculateFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// CalculateFeeMCMSParams encodes MCMS parameters (without Caller) for the CalculateFee choice.
func (e *encoder) CalculateFeeMCMSParams(args CalculateFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CalculateFee", args)
}

// GetAllowedCCVs encodes parameters for the GetAllowedCCVs choice.
func (e *encoder) GetAllowedCCVs(args GetAllowedCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetAllowedCCVs", args)
}

// GetAllowedCCVsMCMSParams encodes MCMS parameters (without Caller) for the GetAllowedCCVs choice.
func (e *encoder) GetAllowedCCVsMCMSParams(args GetAllowedCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetAllowedCCVs", args)
}

// GetAllowedFinalityConfig encodes parameters for the GetAllowedFinalityConfig choice.
func (e *encoder) GetAllowedFinalityConfig(args GetAllowedFinalityConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetAllowedFinalityConfig", args)
}

// GetAllowedFinalityConfigMCMSParams encodes MCMS parameters (without Caller) for the GetAllowedFinalityConfig choice.
func (e *encoder) GetAllowedFinalityConfigMCMSParams(args GetAllowedFinalityConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetAllowedFinalityConfig", args)
}

// GetDestChains encodes parameters for the GetDestChains choice.
func (e *encoder) GetDestChains(args GetDestChains) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChains", args)
}

// GetDestChainsMCMSParams encodes MCMS parameters (without Caller) for the GetDestChains choice.
func (e *encoder) GetDestChainsMCMSParams(args GetDestChainsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDestChains", args)
}

// GetDynamicConfig encodes parameters for the GetDynamicConfig choice.
func (e *encoder) GetDynamicConfig(args GetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDynamicConfig", args)
}

// GetDynamicConfigMCMSParams encodes MCMS parameters (without Caller) for the GetDynamicConfig choice.
func (e *encoder) GetDynamicConfigMCMSParams(args GetDynamicConfigMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetDynamicConfig", args)
}

// GetFee encodes parameters for the GetFee choice.
func (e *encoder) GetFee(args GetFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetFeeMCMSParams encodes MCMS parameters (without Caller) for the GetFee choice.
func (e *encoder) GetFeeMCMSParams(args GetFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetMaxCCVsPerMessage encodes parameters for the GetMaxCCVsPerMessage choice.
func (e *encoder) GetMaxCCVsPerMessage(args GetMaxCCVsPerMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetMaxCCVsPerMessage", args)
}

// GetMaxCCVsPerMessageMCMSParams encodes MCMS parameters (without Caller) for the GetMaxCCVsPerMessage choice.
func (e *encoder) GetMaxCCVsPerMessageMCMSParams(args GetMaxCCVsPerMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetMaxCCVsPerMessage", args)
}

// SetDynamicConfig encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfig(args SetDynamicConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// SetDynamicConfigParams encodes parameters for the SetDynamicConfig choice.
func (e *encoder) SetDynamicConfigParams(args SetDynamicConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDynamicConfig", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
