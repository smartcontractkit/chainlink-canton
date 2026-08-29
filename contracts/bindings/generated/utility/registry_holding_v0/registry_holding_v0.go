package registry_holding_v0

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "utility-registry-holding-v0"
	PackageID   = "8107899ac4723ce986bf7d27416534e576e54b92161e46150a595fb78ff3d3a1"
	SDKVersion  = "3.4.9"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	GetRegistrarInternalScheme = types.TEXT("RegistrarInternalScheme")
	SpliceMetaPrefix           = types.TEXT("splice.lfdecentralizedtrust.org/")
)

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

// CollapseAction is a variant/union type
type CollapseAction struct {
	MergeSplit                    *types.UNIT   `json:"MergeSplit,omitempty"`
	MergeSplitLock                *Lock2        `json:"MergeSplitLock,omitempty"`
	MergeSplitBurn                *types.UNIT   `json:"MergeSplitBurn,omitempty"`
	UnlockMergeSplitBurn          *ExpectedLock `json:"UnlockMergeSplitBurn,omitempty"`
	MergeSplitTransfer            *types.PARTY  `json:"MergeSplitTransfer,omitempty"`
	UnlockMergeSplitTransfer      *Tuple2       `json:"UnlockMergeSplitTransfer,omitempty"`
	AutoUnlockMergeSplitTransfer  *Tuple2       `json:"AutoUnlockMergeSplitTransfer,omitempty"`
	UnlockMergeSplitLockRemaining *types.UNIT   `json:"UnlockMergeSplitLockRemaining,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for CollapseAction
func (v CollapseAction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for CollapseAction
func (v *CollapseAction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes CollapseAction to hex string (Canton MCMS format)
func (v CollapseAction) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes CollapseAction from hex string (Canton MCMS format)
func (v *CollapseAction) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v CollapseAction) GetVariantTag() string {

	if v.MergeSplit != nil {
		return "MergeSplit"
	}

	if v.MergeSplitLock != nil {
		return "MergeSplitLock"
	}

	if v.MergeSplitBurn != nil {
		return "MergeSplitBurn"
	}

	if v.UnlockMergeSplitBurn != nil {
		return "UnlockMergeSplitBurn"
	}

	if v.MergeSplitTransfer != nil {
		return "MergeSplitTransfer"
	}

	if v.UnlockMergeSplitTransfer != nil {
		return "UnlockMergeSplitTransfer"
	}

	if v.AutoUnlockMergeSplitTransfer != nil {
		return "AutoUnlockMergeSplitTransfer"
	}

	if v.UnlockMergeSplitLockRemaining != nil {
		return "UnlockMergeSplitLockRemaining"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v CollapseAction) GetVariantValue() any {

	if v.MergeSplit != nil {
		return v.MergeSplit
	}

	if v.MergeSplitLock != nil {
		return v.MergeSplitLock
	}

	if v.MergeSplitBurn != nil {
		return v.MergeSplitBurn
	}

	if v.UnlockMergeSplitBurn != nil {
		return v.UnlockMergeSplitBurn
	}

	if v.MergeSplitTransfer != nil {
		return v.MergeSplitTransfer
	}

	if v.UnlockMergeSplitTransfer != nil {
		return v.UnlockMergeSplitTransfer
	}

	if v.AutoUnlockMergeSplitTransfer != nil {
		return v.AutoUnlockMergeSplitTransfer
	}

	if v.UnlockMergeSplitLockRemaining != nil {
		return v.UnlockMergeSplitLockRemaining
	}

	return nil
}

var _ types.VARIANT = (*CollapseAction)(nil)

// CollapseActionResult is a Record type
type CollapseActionResult struct {
	Output    *Tuple2            `json:"output" hex:"optional"`
	Remaining *types.CONTRACT_ID `json:"remaining" hex:"optional"`
}

// ToMap converts CollapseActionResult to a map for DAML arguments
func (t CollapseActionResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.Output != nil {
		m["output"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Output),
		}
	} else {
		m["output"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.Remaining != nil {
		m["remaining"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Remaining),
		}
	} else {
		m["remaining"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t CollapseActionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CollapseActionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CollapseActionResult to hex string (Canton MCMS format)
func (t CollapseActionResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CollapseActionResult from hex string (Canton MCMS format)
func (t *CollapseActionResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExpectedLock is a Record type
type ExpectedLock struct {
	Lockers types.SET  `json:"lockers"`
	Context types.TEXT `json:"context"`
}

// ToMap converts ExpectedLock to a map for DAML arguments
func (t ExpectedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["lockers"] = model.NestedToDAMLValue(t.Lockers)

	m["context"] = string(t.Context)

	return m
}

func (t ExpectedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExpectedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExpectedLock to hex string (Canton MCMS format)
func (t ExpectedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExpectedLock from hex string (Canton MCMS format)
func (t *ExpectedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Holding is a Template type
type Holding struct {
	Operator   types.PARTY          `json:"operator"`
	Provider   types.PARTY          `json:"provider"`
	Registrar  types.PARTY          `json:"registrar"`
	Owner      types.PARTY          `json:"owner"`
	Instrument InstrumentIdentifier `json:"instrument"`
	Label      types.TEXT           `json:"label"`
	Amount     types.NUMERIC        `json:"amount"`
	Lock       *Lock2               `json:"lock" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t Holding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t Holding) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t Holding) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrument"] = model.NestedToDAMLValue(t.Instrument)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["label"] = string(t.Label)

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	if t.Lock != nil {
		args["lock"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Lock),
		}
	} else {
		args["lock"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t Holding) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrument"] = model.NestedToDAMLValue(t.Instrument)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["label"] = string(t.Label)

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	if t.Lock != nil {
		args["lock"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Lock),
		}
	} else {
		args["lock"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t Holding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Holding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Holding to hex string (Canton MCMS format)
func (t Holding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Holding from hex string (Canton MCMS format)
func (t *Holding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for Holding

// Archive exercises the Archive choice on this Holding contract via the IHolding interface
// This method uses the package name in the template ID
func (t Holding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t Holding) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// HoldingUnlock exercises the Holding_Unlock choice on this Holding contract
// This method uses the package name in the template ID
func (t Holding) HoldingUnlock(contractID string, args HoldingUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Unlock",
		Arguments:  argsToMap(args),
	}
}

// HoldingUnlockWithPackageID exercises the Holding_Unlock choice using the provided package ID instead of package name
func (t Holding) HoldingUnlockWithPackageID(contractID string, packageID string, args HoldingUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Unlock",
		Arguments:  argsToMap(args),
	}
}

// HoldingMerge exercises the Holding_Merge choice on this Holding contract
// This method uses the package name in the template ID
func (t Holding) HoldingMerge(contractID string, args HoldingMerge) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Merge",
		Arguments:  argsToMap(args),
	}
}

// HoldingMergeWithPackageID exercises the Holding_Merge choice using the provided package ID instead of package name
func (t Holding) HoldingMergeWithPackageID(contractID string, packageID string, args HoldingMerge) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Merge",
		Arguments:  argsToMap(args),
	}
}

// HoldingLock exercises the Holding_Lock choice on this Holding contract
// This method uses the package name in the template ID
func (t Holding) HoldingLock(contractID string, args HoldingLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Lock",
		Arguments:  argsToMap(args),
	}
}

// HoldingLockWithPackageID exercises the Holding_Lock choice using the provided package ID instead of package name
func (t Holding) HoldingLockWithPackageID(contractID string, packageID string, args HoldingLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Lock",
		Arguments:  argsToMap(args),
	}
}

// HoldingTransfer exercises the Holding_Transfer choice on this Holding contract
// This method uses the package name in the template ID
func (t Holding) HoldingTransfer(contractID string, args HoldingTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Transfer",
		Arguments:  argsToMap(args),
	}
}

// HoldingTransferWithPackageID exercises the Holding_Transfer choice using the provided package ID instead of package name
func (t Holding) HoldingTransferWithPackageID(contractID string, packageID string, args HoldingTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Transfer",
		Arguments:  argsToMap(args),
	}
}

// HoldingSplit exercises the Holding_Split choice on this Holding contract
// This method uses the package name in the template ID
func (t Holding) HoldingSplit(contractID string, args HoldingSplit) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Split",
		Arguments:  argsToMap(args),
	}
}

// HoldingSplitWithPackageID exercises the Holding_Split choice using the provided package ID instead of package name
func (t Holding) HoldingSplitWithPackageID(contractID string, packageID string, args HoldingSplit) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.Holding.V0.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Holding_Split",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Holding

var _ splice_api_token_holding_v1.IHolding = (*Holding)(nil)

// HoldingLock is a Record type
type HoldingLock struct {
	Lockers   types.SET  `json:"lockers"`
	Context   types.TEXT `json:"context"`
	Observers *types.SET `json:"observers" hex:"optional"`
}

// ToMap converts HoldingLock to a map for DAML arguments
func (t HoldingLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["lockers"] = model.NestedToDAMLValue(t.Lockers)

	m["context"] = string(t.Context)

	if t.Observers != nil {
		m["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		m["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t HoldingLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingLock to hex string (Canton MCMS format)
func (t HoldingLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingLock from hex string (Canton MCMS format)
func (t *HoldingLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingLockResult is a Record type
type HoldingLockResult struct {
	HoldingCid types.CONTRACT_ID `json:"holdingCid"`
}

// ToMap converts HoldingLockResult to a map for DAML arguments
func (t HoldingLockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	return m
}

func (t HoldingLockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingLockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingLockResult to hex string (Canton MCMS format)
func (t HoldingLockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingLockResult from hex string (Canton MCMS format)
func (t *HoldingLockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingMerge is a Record type
type HoldingMerge struct {
	HoldingCids []types.CONTRACT_ID `json:"holdingCids"`
}

// ToMap converts HoldingMerge to a map for DAML arguments
func (t HoldingMerge) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t HoldingMerge) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingMerge) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingMerge to hex string (Canton MCMS format)
func (t HoldingMerge) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingMerge from hex string (Canton MCMS format)
func (t *HoldingMerge) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingMergeResult is a Record type
type HoldingMergeResult struct {
	HoldingCid types.CONTRACT_ID                      `json:"holdingCid"`
	Meta       *splice_api_token_metadata_v1.Metadata `json:"meta" hex:"optional"`
}

// ToMap converts HoldingMergeResult to a map for DAML arguments
func (t HoldingMergeResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	if t.Meta != nil {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Meta),
		}
	} else {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t HoldingMergeResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingMergeResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingMergeResult to hex string (Canton MCMS format)
func (t HoldingMergeResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingMergeResult from hex string (Canton MCMS format)
func (t *HoldingMergeResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingSplit is a Record type
type HoldingSplit struct {
	Amounts []types.NUMERIC `json:"amounts"`
}

// ToMap converts HoldingSplit to a map for DAML arguments
func (t HoldingSplit) ToMap() map[string]any {
	m := make(map[string]any)

	m["amounts"] = func() []any {
		res := make([]any, 0, len(t.Amounts))
		for _, e := range t.Amounts {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t HoldingSplit) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingSplit) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingSplit to hex string (Canton MCMS format)
func (t HoldingSplit) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingSplit from hex string (Canton MCMS format)
func (t *HoldingSplit) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingSplitResult is a Record type
type HoldingSplitResult struct {
	SplitCids []types.CONTRACT_ID                    `json:"splitCids"`
	Remaining *types.CONTRACT_ID                     `json:"remaining" hex:"optional"`
	Meta      *splice_api_token_metadata_v1.Metadata `json:"meta" hex:"optional"`
}

// ToMap converts HoldingSplitResult to a map for DAML arguments
func (t HoldingSplitResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["splitCids"] = func() []any {
		res := make([]any, 0, len(t.SplitCids))
		for _, e := range t.SplitCids {
			res = append(res, e)
		}
		return res
	}()

	if t.Remaining != nil {
		m["remaining"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Remaining),
		}
	} else {
		m["remaining"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.Meta != nil {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Meta),
		}
	} else {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t HoldingSplitResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingSplitResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingSplitResult to hex string (Canton MCMS format)
func (t HoldingSplitResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingSplitResult from hex string (Canton MCMS format)
func (t *HoldingSplitResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingTransfer is a Record type
type HoldingTransfer struct {
	NewOwner types.PARTY `json:"newOwner"`
	NewLabel types.TEXT  `json:"newLabel"`
}

// ToMap converts HoldingTransfer to a map for DAML arguments
func (t HoldingTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["newOwner"] = t.NewOwner.ToMap()

	m["newLabel"] = string(t.NewLabel)

	return m
}

func (t HoldingTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingTransfer to hex string (Canton MCMS format)
func (t HoldingTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingTransfer from hex string (Canton MCMS format)
func (t *HoldingTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingTransferResult is a Record type
type HoldingTransferResult struct {
	HoldingCid types.CONTRACT_ID `json:"holdingCid"`
}

// ToMap converts HoldingTransferResult to a map for DAML arguments
func (t HoldingTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	return m
}

func (t HoldingTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingTransferResult to hex string (Canton MCMS format)
func (t HoldingTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingTransferResult from hex string (Canton MCMS format)
func (t *HoldingTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingUnlock is a Record type
type HoldingUnlock struct {
}

// ToMap converts HoldingUnlock to a map for DAML arguments
func (t HoldingUnlock) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t HoldingUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingUnlock to hex string (Canton MCMS format)
func (t HoldingUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingUnlock from hex string (Canton MCMS format)
func (t *HoldingUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HoldingUnlockResult is a Record type
type HoldingUnlockResult struct {
	HoldingCid types.CONTRACT_ID                      `json:"holdingCid"`
	Meta       *splice_api_token_metadata_v1.Metadata `json:"meta" hex:"optional"`
}

// ToMap converts HoldingUnlockResult to a map for DAML arguments
func (t HoldingUnlockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	if t.Meta != nil {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Meta),
		}
	} else {
		m["meta"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t HoldingUnlockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HoldingUnlockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HoldingUnlockResult to hex string (Canton MCMS format)
func (t HoldingUnlockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HoldingUnlockResult from hex string (Canton MCMS format)
func (t *HoldingUnlockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentIdentifier is a Record type
type InstrumentIdentifier struct {
	Source types.PARTY `json:"source"`
	Id     types.TEXT  `json:"id"`
	Scheme types.TEXT  `json:"scheme"`
}

// ToMap converts InstrumentIdentifier to a map for DAML arguments
func (t InstrumentIdentifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["source"] = t.Source.ToMap()

	m["id"] = string(t.Id)

	m["scheme"] = string(t.Scheme)

	return m
}

func (t InstrumentIdentifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentIdentifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentIdentifier to hex string (Canton MCMS format)
func (t InstrumentIdentifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentIdentifier from hex string (Canton MCMS format)
func (t *InstrumentIdentifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Lock2 is a Record type
type Lock2 struct {
	Lockers   types.SET  `json:"lockers"`
	Context   types.TEXT `json:"context"`
	Observers *types.SET `json:"observers" hex:"optional"`
}

// ToMap converts Lock2 to a map for DAML arguments
func (t Lock2) ToMap() map[string]any {
	m := make(map[string]any)

	m["lockers"] = model.NestedToDAMLValue(t.Lockers)

	m["context"] = string(t.Context)

	if t.Observers != nil {
		m["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		m["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t Lock2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Lock2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Lock2 to hex string (Canton MCMS format)
func (t Lock2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Lock2 from hex string (Canton MCMS format)
func (t *Lock2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperator is a Record type
type WithOperator struct {
	Operator types.PARTY `json:"operator"`
}

// ToMap converts WithOperator to a map for DAML arguments
func (t WithOperator) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	return m
}

func (t WithOperator) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperator) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperator to hex string (Canton MCMS format)
func (t WithOperator) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperator from hex string (Canton MCMS format)
func (t *WithOperator) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorProvider is a Record type
type WithOperatorProvider struct {
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
}

// ToMap converts WithOperatorProvider to a map for DAML arguments
func (t WithOperatorProvider) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	return m
}

func (t WithOperatorProvider) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorProvider) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorProvider to hex string (Canton MCMS format)
func (t WithOperatorProvider) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorProvider from hex string (Canton MCMS format)
func (t *WithOperatorProvider) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorProviderRegistrar is a Record type
type WithOperatorProviderRegistrar struct {
	Operator  types.PARTY `json:"operator"`
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
}

// ToMap converts WithOperatorProviderRegistrar to a map for DAML arguments
func (t WithOperatorProviderRegistrar) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	return m
}

func (t WithOperatorProviderRegistrar) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorProviderRegistrar) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorProviderRegistrar to hex string (Canton MCMS format)
func (t WithOperatorProviderRegistrar) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorProviderRegistrar from hex string (Canton MCMS format)
func (t *WithOperatorProviderRegistrar) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorProviderRegistrarHolder is a Record type
type WithOperatorProviderRegistrarHolder struct {
	Operator  types.PARTY `json:"operator"`
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
	Holder    types.PARTY `json:"holder"`
}

// ToMap converts WithOperatorProviderRegistrarHolder to a map for DAML arguments
func (t WithOperatorProviderRegistrarHolder) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	return m
}

func (t WithOperatorProviderRegistrarHolder) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorProviderRegistrarHolder) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorProviderRegistrarHolder to hex string (Canton MCMS format)
func (t WithOperatorProviderRegistrarHolder) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorProviderRegistrarHolder from hex string (Canton MCMS format)
func (t *WithOperatorProviderRegistrarHolder) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorProviderRegistrarHolderInstrument is a Record type
type WithOperatorProviderRegistrarHolderInstrument struct {
	Operator   types.PARTY          `json:"operator"`
	Provider   types.PARTY          `json:"provider"`
	Registrar  types.PARTY          `json:"registrar"`
	Holder     types.PARTY          `json:"holder"`
	Instrument InstrumentIdentifier `json:"instrument"`
}

// ToMap converts WithOperatorProviderRegistrarHolderInstrument to a map for DAML arguments
func (t WithOperatorProviderRegistrarHolderInstrument) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["instrument"] = model.NestedToDAMLValue(t.Instrument)

	return m
}

func (t WithOperatorProviderRegistrarHolderInstrument) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorProviderRegistrarHolderInstrument) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorProviderRegistrarHolderInstrument to hex string (Canton MCMS format)
func (t WithOperatorProviderRegistrarHolderInstrument) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorProviderRegistrarHolderInstrument from hex string (Canton MCMS format)
func (t *WithOperatorProviderRegistrarHolderInstrument) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorRegistrar is a Record type
type WithOperatorRegistrar struct {
	Operator  types.PARTY `json:"operator"`
	Registrar types.PARTY `json:"registrar"`
}

// ToMap converts WithOperatorRegistrar to a map for DAML arguments
func (t WithOperatorRegistrar) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	return m
}

func (t WithOperatorRegistrar) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorRegistrar) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorRegistrar to hex string (Canton MCMS format)
func (t WithOperatorRegistrar) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorRegistrar from hex string (Canton MCMS format)
func (t *WithOperatorRegistrar) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// WithOperatorRegistrarHolderInstrument is a Record type
type WithOperatorRegistrarHolderInstrument struct {
	Operator   types.PARTY          `json:"operator"`
	Registrar  types.PARTY          `json:"registrar"`
	Holder     types.PARTY          `json:"holder"`
	Instrument InstrumentIdentifier `json:"instrument"`
}

// ToMap converts WithOperatorRegistrarHolderInstrument to a map for DAML arguments
func (t WithOperatorRegistrarHolderInstrument) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["instrument"] = model.NestedToDAMLValue(t.Instrument)

	return m
}

func (t WithOperatorRegistrarHolderInstrument) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *WithOperatorRegistrarHolderInstrument) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes WithOperatorRegistrarHolderInstrument to hex string (Canton MCMS format)
func (t WithOperatorRegistrarHolderInstrument) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes WithOperatorRegistrarHolderInstrument from hex string (Canton MCMS format)
func (t *WithOperatorRegistrarHolderInstrument) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	HoldingLock(args HoldingLock) (*bind.EncodedChoice, error)
	HoldingMerge(args HoldingMerge) (*bind.EncodedChoice, error)
	HoldingSplit(args HoldingSplit) (*bind.EncodedChoice, error)
	HoldingTransfer(args HoldingTransfer) (*bind.EncodedChoice, error)
	HoldingUnlock(args HoldingUnlock) (*bind.EncodedChoice, error)
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

// HoldingLock encodes parameters for the Holding_Lock choice.
func (e *encoder) HoldingLock(args HoldingLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Holding_Lock", args)
}

// HoldingMerge encodes parameters for the Holding_Merge choice.
func (e *encoder) HoldingMerge(args HoldingMerge) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Holding_Merge", args)
}

// HoldingSplit encodes parameters for the Holding_Split choice.
func (e *encoder) HoldingSplit(args HoldingSplit) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Holding_Split", args)
}

// HoldingTransfer encodes parameters for the Holding_Transfer choice.
func (e *encoder) HoldingTransfer(args HoldingTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Holding_Transfer", args)
}

// HoldingUnlock encodes parameters for the Holding_Unlock choice.
func (e *encoder) HoldingUnlock(args HoldingUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Holding_Unlock", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
