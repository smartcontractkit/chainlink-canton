package registry_v0

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_featured_app_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_featured_app_v1"
	splice_api_token_allocation_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_allocation_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	splice_api_token_transfer_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	credential_v0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/credential_v0"
	registry_holding_v0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_holding_v0"
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
	PackageName = "utility-registry-v0"
	PackageID   = "a236e8e22a3b5f199e37d5554e82bafd2df688f901de02b00be3964bdfa8c1ab"
	SDKVersion  = "3.4.9"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	UtilityPrefix = types.TEXT("utility.digitalasset.com/")
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

// AcceptedBurn is a Template type
type AcceptedBurn struct {
	Burn         Burn       `json:"burn"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AcceptedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurn to hex string (Canton MCMS format)
func (t AcceptedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurn from hex string (Canton MCMS format)
func (t *AcceptedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedBurn

// AcceptedBurnClean exercises the AcceptedBurn_Clean choice on this AcceptedBurn contract
// This method uses the package name in the template ID
func (t AcceptedBurn) AcceptedBurnClean(contractID string, args AcceptedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedBurnCleanWithPackageID exercises the AcceptedBurn_Clean choice using the provided package ID instead of package name
func (t AcceptedBurn) AcceptedBurnCleanWithPackageID(contractID string, packageID string, args AcceptedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedBurnExecute exercises the AcceptedBurn_Execute choice on this AcceptedBurn contract
// This method uses the package name in the template ID
func (t AcceptedBurn) AcceptedBurnExecute(contractID string, args AcceptedBurnExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedBurnExecuteWithPackageID exercises the AcceptedBurn_Execute choice using the provided package ID instead of package name
func (t AcceptedBurn) AcceptedBurnExecuteWithPackageID(contractID string, packageID string, args AcceptedBurnExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedBurnFail exercises the AcceptedBurn_Fail choice on this AcceptedBurn contract
// This method uses the package name in the template ID
func (t AcceptedBurn) AcceptedBurnFail(contractID string, args AcceptedBurnFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedBurnFailWithPackageID exercises the AcceptedBurn_Fail choice using the provided package ID instead of package name
func (t AcceptedBurn) AcceptedBurnFailWithPackageID(contractID string, packageID string, args AcceptedBurnFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "AcceptedBurn_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedBurn contract
// This method uses the package name in the template ID
func (t AcceptedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "AcceptedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedBurnClean is a Record type
type AcceptedBurnClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts AcceptedBurnClean to a map for DAML arguments
func (t AcceptedBurnClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AcceptedBurnClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurnClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurnClean to hex string (Canton MCMS format)
func (t AcceptedBurnClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurnClean from hex string (Canton MCMS format)
func (t *AcceptedBurnClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedBurnExecute is a Record type
type AcceptedBurnExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	CredentialCids             []types.CONTRACT_ID `json:"credentialCids"`
	HoldingCid                 types.CONTRACT_ID   `json:"holdingCid"`
}

// ToMap converts AcceptedBurnExecute to a map for DAML arguments
func (t AcceptedBurnExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	return m
}

func (t AcceptedBurnExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurnExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurnExecute to hex string (Canton MCMS format)
func (t AcceptedBurnExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurnExecute from hex string (Canton MCMS format)
func (t *AcceptedBurnExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedBurnExecuteResult is a Record type
type AcceptedBurnExecuteResult struct {
	ExecutedBurnCid types.CONTRACT_ID                      `json:"executedBurnCid"`
	Meta            *splice_api_token_metadata_v1.Metadata `json:"meta" hex:"optional"`
}

// ToMap converts AcceptedBurnExecuteResult to a map for DAML arguments
func (t AcceptedBurnExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["executedBurnCid"] = model.NestedToDAMLValue(t.ExecutedBurnCid)

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

func (t AcceptedBurnExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurnExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurnExecuteResult to hex string (Canton MCMS format)
func (t AcceptedBurnExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurnExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedBurnExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedBurnFail is a Record type
type AcceptedBurnFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedBurnFail to a map for DAML arguments
func (t AcceptedBurnFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedBurnFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurnFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurnFail to hex string (Canton MCMS format)
func (t AcceptedBurnFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurnFail from hex string (Canton MCMS format)
func (t *AcceptedBurnFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedBurnFailResult is a Record type
type AcceptedBurnFailResult struct {
	FailedBurnCid types.CONTRACT_ID `json:"failedBurnCid"`
}

// ToMap converts AcceptedBurnFailResult to a map for DAML arguments
func (t AcceptedBurnFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedBurnCid"] = model.NestedToDAMLValue(t.FailedBurnCid)

	return m
}

func (t AcceptedBurnFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedBurnFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedBurnFailResult to hex string (Canton MCMS format)
func (t AcceptedBurnFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedBurnFailResult from hex string (Canton MCMS format)
func (t *AcceptedBurnFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedForceTransfer is a Template type
type AcceptedForceTransfer struct {
	ForceTransfer      ForceTransfer `json:"forceTransfer"`
	RegistrarRationale types.TEXT    `json:"registrarRationale"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedForceTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedForceTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedForceTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRationale"] = string(t.RegistrarRationale)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AcceptedForceTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRationale"] = string(t.RegistrarRationale)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedForceTransfer to hex string (Canton MCMS format)
func (t AcceptedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedForceTransfer from hex string (Canton MCMS format)
func (t *AcceptedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedForceTransfer

// AcceptedForceTransferExecute exercises the AcceptedForceTransfer_Execute choice on this AcceptedForceTransfer contract
// This method uses the package name in the template ID
func (t AcceptedForceTransfer) AcceptedForceTransferExecute(contractID string, args AcceptedForceTransferExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedForceTransfer_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedForceTransferExecuteWithPackageID exercises the AcceptedForceTransfer_Execute choice using the provided package ID instead of package name
func (t AcceptedForceTransfer) AcceptedForceTransferExecuteWithPackageID(contractID string, packageID string, args AcceptedForceTransferExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedForceTransfer_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedForceTransferFail exercises the AcceptedForceTransfer_Fail choice on this AcceptedForceTransfer contract
// This method uses the package name in the template ID
func (t AcceptedForceTransfer) AcceptedForceTransferFail(contractID string, args AcceptedForceTransferFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedForceTransfer_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedForceTransferFailWithPackageID exercises the AcceptedForceTransfer_Fail choice using the provided package ID instead of package name
func (t AcceptedForceTransfer) AcceptedForceTransferFailWithPackageID(contractID string, packageID string, args AcceptedForceTransferFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedForceTransfer_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedForceTransfer contract
// This method uses the package name in the template ID
func (t AcceptedForceTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedForceTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "AcceptedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedForceTransferExecute is a Record type
type AcceptedForceTransferExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	HoldingCids                []types.CONTRACT_ID `json:"holdingCids"`
	RequestorCredentialCids    []types.CONTRACT_ID `json:"requestorCredentialCids"`
}

// ToMap converts AcceptedForceTransferExecute to a map for DAML arguments
func (t AcceptedForceTransferExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["requestorCredentialCids"] = func() []any {
		res := make([]any, 0, len(t.RequestorCredentialCids))
		for _, e := range t.RequestorCredentialCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t AcceptedForceTransferExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedForceTransferExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedForceTransferExecute to hex string (Canton MCMS format)
func (t AcceptedForceTransferExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedForceTransferExecute from hex string (Canton MCMS format)
func (t *AcceptedForceTransferExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedForceTransferExecuteResult is a Record type
type AcceptedForceTransferExecuteResult struct {
	HoldingTransferResult    registry_holding_v0.HoldingTransferResult `json:"holdingTransferResult"`
	ExecutedForceTransferCid types.CONTRACT_ID                         `json:"executedForceTransferCid"`
	RemainingHoldingCids     []types.CONTRACT_ID                       `json:"remainingHoldingCids"`
}

// ToMap converts AcceptedForceTransferExecuteResult to a map for DAML arguments
func (t AcceptedForceTransferExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingTransferResult"] = model.NestedToDAMLValue(t.HoldingTransferResult)

	m["executedForceTransferCid"] = model.NestedToDAMLValue(t.ExecutedForceTransferCid)

	m["remainingHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.RemainingHoldingCids))
		for _, e := range t.RemainingHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t AcceptedForceTransferExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedForceTransferExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedForceTransferExecuteResult to hex string (Canton MCMS format)
func (t AcceptedForceTransferExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedForceTransferExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedForceTransferExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedForceTransferFail is a Record type
type AcceptedForceTransferFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedForceTransferFail to a map for DAML arguments
func (t AcceptedForceTransferFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedForceTransferFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedForceTransferFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedForceTransferFail to hex string (Canton MCMS format)
func (t AcceptedForceTransferFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedForceTransferFail from hex string (Canton MCMS format)
func (t *AcceptedForceTransferFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedForceTransferFailResult is a Record type
type AcceptedForceTransferFailResult struct {
	FailedForceTransferCid types.CONTRACT_ID `json:"failedForceTransferCid"`
}

// ToMap converts AcceptedForceTransferFailResult to a map for DAML arguments
func (t AcceptedForceTransferFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedForceTransferCid"] = model.NestedToDAMLValue(t.FailedForceTransferCid)

	return m
}

func (t AcceptedForceTransferFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedForceTransferFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedForceTransferFailResult to hex string (Canton MCMS format)
func (t AcceptedForceTransferFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedForceTransferFailResult from hex string (Canton MCMS format)
func (t *AcceptedForceTransferFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedLock is a Template type
type AcceptedLock struct {
	Lock         Lock3      `json:"lock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedLock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "AcceptedLock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedLock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "AcceptedLock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedLock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AcceptedLock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLock to hex string (Canton MCMS format)
func (t AcceptedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLock from hex string (Canton MCMS format)
func (t *AcceptedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedLock

// AcceptedLockClean exercises the AcceptedLock_Clean choice on this AcceptedLock contract
// This method uses the package name in the template ID
func (t AcceptedLock) AcceptedLockClean(contractID string, args AcceptedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedLockCleanWithPackageID exercises the AcceptedLock_Clean choice using the provided package ID instead of package name
func (t AcceptedLock) AcceptedLockCleanWithPackageID(contractID string, packageID string, args AcceptedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedLockExecute exercises the AcceptedLock_Execute choice on this AcceptedLock contract
// This method uses the package name in the template ID
func (t AcceptedLock) AcceptedLockExecute(contractID string, args AcceptedLockExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedLockExecuteWithPackageID exercises the AcceptedLock_Execute choice using the provided package ID instead of package name
func (t AcceptedLock) AcceptedLockExecuteWithPackageID(contractID string, packageID string, args AcceptedLockExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedLockFail exercises the AcceptedLock_Fail choice on this AcceptedLock contract
// This method uses the package name in the template ID
func (t AcceptedLock) AcceptedLockFail(contractID string, args AcceptedLockFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedLockFailWithPackageID exercises the AcceptedLock_Fail choice using the provided package ID instead of package name
func (t AcceptedLock) AcceptedLockFailWithPackageID(contractID string, packageID string, args AcceptedLockFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "AcceptedLock_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedLock contract
// This method uses the package name in the template ID
func (t AcceptedLock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedLock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "AcceptedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedLockClean is a Record type
type AcceptedLockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts AcceptedLockClean to a map for DAML arguments
func (t AcceptedLockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AcceptedLockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLockClean to hex string (Canton MCMS format)
func (t AcceptedLockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLockClean from hex string (Canton MCMS format)
func (t *AcceptedLockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedLockExecute is a Record type
type AcceptedLockExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	CredentialCids             []types.CONTRACT_ID `json:"credentialCids"`
	HoldingCid                 types.CONTRACT_ID   `json:"holdingCid"`
}

// ToMap converts AcceptedLockExecute to a map for DAML arguments
func (t AcceptedLockExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	return m
}

func (t AcceptedLockExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLockExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLockExecute to hex string (Canton MCMS format)
func (t AcceptedLockExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLockExecute from hex string (Canton MCMS format)
func (t *AcceptedLockExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedLockExecuteResult is a Record type
type AcceptedLockExecuteResult struct {
	HoldingLockResult registry_holding_v0.HoldingLockResult `json:"holdingLockResult"`
	ExecutedLockCid   types.CONTRACT_ID                     `json:"executedLockCid"`
}

// ToMap converts AcceptedLockExecuteResult to a map for DAML arguments
func (t AcceptedLockExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingLockResult"] = model.NestedToDAMLValue(t.HoldingLockResult)

	m["executedLockCid"] = model.NestedToDAMLValue(t.ExecutedLockCid)

	return m
}

func (t AcceptedLockExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLockExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLockExecuteResult to hex string (Canton MCMS format)
func (t AcceptedLockExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLockExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedLockExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedLockFail is a Record type
type AcceptedLockFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedLockFail to a map for DAML arguments
func (t AcceptedLockFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedLockFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLockFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLockFail to hex string (Canton MCMS format)
func (t AcceptedLockFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLockFail from hex string (Canton MCMS format)
func (t *AcceptedLockFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedLockFailResult is a Record type
type AcceptedLockFailResult struct {
	FailedLockCid types.CONTRACT_ID `json:"failedLockCid"`
}

// ToMap converts AcceptedLockFailResult to a map for DAML arguments
func (t AcceptedLockFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedLockCid"] = model.NestedToDAMLValue(t.FailedLockCid)

	return m
}

func (t AcceptedLockFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedLockFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedLockFailResult to hex string (Canton MCMS format)
func (t AcceptedLockFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedLockFailResult from hex string (Canton MCMS format)
func (t *AcceptedLockFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedMint is a Template type
type AcceptedMint struct {
	Mint         Mint       `json:"mint"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "AcceptedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "AcceptedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AcceptedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMint to hex string (Canton MCMS format)
func (t AcceptedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMint from hex string (Canton MCMS format)
func (t *AcceptedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedMint

// AcceptedMintClean exercises the AcceptedMint_Clean choice on this AcceptedMint contract
// This method uses the package name in the template ID
func (t AcceptedMint) AcceptedMintClean(contractID string, args AcceptedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedMintCleanWithPackageID exercises the AcceptedMint_Clean choice using the provided package ID instead of package name
func (t AcceptedMint) AcceptedMintCleanWithPackageID(contractID string, packageID string, args AcceptedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedMintExecute exercises the AcceptedMint_Execute choice on this AcceptedMint contract
// This method uses the package name in the template ID
func (t AcceptedMint) AcceptedMintExecute(contractID string, args AcceptedMintExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedMintExecuteWithPackageID exercises the AcceptedMint_Execute choice using the provided package ID instead of package name
func (t AcceptedMint) AcceptedMintExecuteWithPackageID(contractID string, packageID string, args AcceptedMintExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedMintFail exercises the AcceptedMint_Fail choice on this AcceptedMint contract
// This method uses the package name in the template ID
func (t AcceptedMint) AcceptedMintFail(contractID string, args AcceptedMintFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedMintFailWithPackageID exercises the AcceptedMint_Fail choice using the provided package ID instead of package name
func (t AcceptedMint) AcceptedMintFailWithPackageID(contractID string, packageID string, args AcceptedMintFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "AcceptedMint_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedMint contract
// This method uses the package name in the template ID
func (t AcceptedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "AcceptedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedMintClean is a Record type
type AcceptedMintClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts AcceptedMintClean to a map for DAML arguments
func (t AcceptedMintClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AcceptedMintClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMintClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMintClean to hex string (Canton MCMS format)
func (t AcceptedMintClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMintClean from hex string (Canton MCMS format)
func (t *AcceptedMintClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedMintExecute is a Record type
type AcceptedMintExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	CredentialCids             []types.CONTRACT_ID `json:"credentialCids"`
}

// ToMap converts AcceptedMintExecute to a map for DAML arguments
func (t AcceptedMintExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t AcceptedMintExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMintExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMintExecute to hex string (Canton MCMS format)
func (t AcceptedMintExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMintExecute from hex string (Canton MCMS format)
func (t *AcceptedMintExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedMintExecuteResult is a Record type
type AcceptedMintExecuteResult struct {
	HoldingCid      types.CONTRACT_ID                      `json:"holdingCid"`
	ExecutedMintCid types.CONTRACT_ID                      `json:"executedMintCid"`
	Meta            *splice_api_token_metadata_v1.Metadata `json:"meta" hex:"optional"`
}

// ToMap converts AcceptedMintExecuteResult to a map for DAML arguments
func (t AcceptedMintExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	m["executedMintCid"] = model.NestedToDAMLValue(t.ExecutedMintCid)

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

func (t AcceptedMintExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMintExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMintExecuteResult to hex string (Canton MCMS format)
func (t AcceptedMintExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMintExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedMintExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedMintFail is a Record type
type AcceptedMintFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedMintFail to a map for DAML arguments
func (t AcceptedMintFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedMintFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMintFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMintFail to hex string (Canton MCMS format)
func (t AcceptedMintFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMintFail from hex string (Canton MCMS format)
func (t *AcceptedMintFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedMintFailResult is a Record type
type AcceptedMintFailResult struct {
	FailedMintCid types.CONTRACT_ID `json:"failedMintCid"`
}

// ToMap converts AcceptedMintFailResult to a map for DAML arguments
func (t AcceptedMintFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedMintCid"] = model.NestedToDAMLValue(t.FailedMintCid)

	return m
}

func (t AcceptedMintFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedMintFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedMintFailResult to hex string (Canton MCMS format)
func (t AcceptedMintFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedMintFailResult from hex string (Canton MCMS format)
func (t *AcceptedMintFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedTransfer is a Template type
type AcceptedTransfer struct {
	Transfer      Transfer2  `json:"transfer"`
	SenderLabel   types.TEXT `json:"senderLabel"`
	ReceiverLabel types.TEXT `json:"receiverLabel"`
	Observers     *types.SET `json:"observers" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
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
func (t AcceptedTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransfer to hex string (Canton MCMS format)
func (t AcceptedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransfer from hex string (Canton MCMS format)
func (t *AcceptedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedTransfer

// AcceptedTransferClean exercises the AcceptedTransfer_Clean choice on this AcceptedTransfer contract
// This method uses the package name in the template ID
func (t AcceptedTransfer) AcceptedTransferClean(contractID string, args AcceptedTransferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedTransferCleanWithPackageID exercises the AcceptedTransfer_Clean choice using the provided package ID instead of package name
func (t AcceptedTransfer) AcceptedTransferCleanWithPackageID(contractID string, packageID string, args AcceptedTransferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedTransferExecute exercises the AcceptedTransfer_Execute choice on this AcceptedTransfer contract
// This method uses the package name in the template ID
func (t AcceptedTransfer) AcceptedTransferExecute(contractID string, args AcceptedTransferExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedTransferExecuteWithPackageID exercises the AcceptedTransfer_Execute choice using the provided package ID instead of package name
func (t AcceptedTransfer) AcceptedTransferExecuteWithPackageID(contractID string, packageID string, args AcceptedTransferExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedTransferFail exercises the AcceptedTransfer_Fail choice on this AcceptedTransfer contract
// This method uses the package name in the template ID
func (t AcceptedTransfer) AcceptedTransferFail(contractID string, args AcceptedTransferFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedTransferFailWithPackageID exercises the AcceptedTransfer_Fail choice using the provided package ID instead of package name
func (t AcceptedTransfer) AcceptedTransferFailWithPackageID(contractID string, packageID string, args AcceptedTransferFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "AcceptedTransfer_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedTransfer contract
// This method uses the package name in the template ID
func (t AcceptedTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "AcceptedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedTransferClean is a Record type
type AcceptedTransferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts AcceptedTransferClean to a map for DAML arguments
func (t AcceptedTransferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AcceptedTransferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransferClean to hex string (Canton MCMS format)
func (t AcceptedTransferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransferClean from hex string (Canton MCMS format)
func (t *AcceptedTransferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedTransferExecute is a Record type
type AcceptedTransferExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	SenderCredentialCids       []types.CONTRACT_ID `json:"senderCredentialCids"`
	ReceiverCredentialCids     []types.CONTRACT_ID `json:"receiverCredentialCids"`
	HoldingCid                 types.CONTRACT_ID   `json:"holdingCid"`
}

// ToMap converts AcceptedTransferExecute to a map for DAML arguments
func (t AcceptedTransferExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["senderCredentialCids"] = func() []any {
		res := make([]any, 0, len(t.SenderCredentialCids))
		for _, e := range t.SenderCredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["receiverCredentialCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverCredentialCids))
		for _, e := range t.ReceiverCredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	return m
}

func (t AcceptedTransferExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransferExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransferExecute to hex string (Canton MCMS format)
func (t AcceptedTransferExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransferExecute from hex string (Canton MCMS format)
func (t *AcceptedTransferExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedTransferExecuteResult is a Record type
type AcceptedTransferExecuteResult struct {
	HoldingTransferResult registry_holding_v0.HoldingTransferResult `json:"holdingTransferResult"`
	ExecutedTransferCid   types.CONTRACT_ID                         `json:"executedTransferCid"`
	Meta                  *splice_api_token_metadata_v1.Metadata    `json:"meta" hex:"optional"`
}

// ToMap converts AcceptedTransferExecuteResult to a map for DAML arguments
func (t AcceptedTransferExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingTransferResult"] = model.NestedToDAMLValue(t.HoldingTransferResult)

	m["executedTransferCid"] = model.NestedToDAMLValue(t.ExecutedTransferCid)

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

func (t AcceptedTransferExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransferExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransferExecuteResult to hex string (Canton MCMS format)
func (t AcceptedTransferExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransferExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedTransferExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedTransferFail is a Record type
type AcceptedTransferFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedTransferFail to a map for DAML arguments
func (t AcceptedTransferFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedTransferFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransferFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransferFail to hex string (Canton MCMS format)
func (t AcceptedTransferFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransferFail from hex string (Canton MCMS format)
func (t *AcceptedTransferFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedTransferFailResult is a Record type
type AcceptedTransferFailResult struct {
	FailedTransferCid types.CONTRACT_ID `json:"failedTransferCid"`
}

// ToMap converts AcceptedTransferFailResult to a map for DAML arguments
func (t AcceptedTransferFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedTransferCid"] = model.NestedToDAMLValue(t.FailedTransferCid)

	return m
}

func (t AcceptedTransferFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedTransferFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedTransferFailResult to hex string (Canton MCMS format)
func (t AcceptedTransferFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedTransferFailResult from hex string (Canton MCMS format)
func (t *AcceptedTransferFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedUnlock is a Template type
type AcceptedUnlock struct {
	Unlock       Unlock     `json:"unlock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AcceptedUnlock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AcceptedUnlock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AcceptedUnlock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AcceptedUnlock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AcceptedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlock to hex string (Canton MCMS format)
func (t AcceptedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlock from hex string (Canton MCMS format)
func (t *AcceptedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AcceptedUnlock

// AcceptedUnlockClean exercises the AcceptedUnlock_Clean choice on this AcceptedUnlock contract
// This method uses the package name in the template ID
func (t AcceptedUnlock) AcceptedUnlockClean(contractID string, args AcceptedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedUnlockCleanWithPackageID exercises the AcceptedUnlock_Clean choice using the provided package ID instead of package name
func (t AcceptedUnlock) AcceptedUnlockCleanWithPackageID(contractID string, packageID string, args AcceptedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// AcceptedUnlockExecute exercises the AcceptedUnlock_Execute choice on this AcceptedUnlock contract
// This method uses the package name in the template ID
func (t AcceptedUnlock) AcceptedUnlockExecute(contractID string, args AcceptedUnlockExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedUnlockExecuteWithPackageID exercises the AcceptedUnlock_Execute choice using the provided package ID instead of package name
func (t AcceptedUnlock) AcceptedUnlockExecuteWithPackageID(contractID string, packageID string, args AcceptedUnlockExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Execute",
		Arguments:  argsToMap(args),
	}
}

// AcceptedUnlockFail exercises the AcceptedUnlock_Fail choice on this AcceptedUnlock contract
// This method uses the package name in the template ID
func (t AcceptedUnlock) AcceptedUnlockFail(contractID string, args AcceptedUnlockFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Fail",
		Arguments:  argsToMap(args),
	}
}

// AcceptedUnlockFailWithPackageID exercises the AcceptedUnlock_Fail choice using the provided package ID instead of package name
func (t AcceptedUnlock) AcceptedUnlockFailWithPackageID(contractID string, packageID string, args AcceptedUnlockFail) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "AcceptedUnlock_Fail",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AcceptedUnlock contract
// This method uses the package name in the template ID
func (t AcceptedUnlock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AcceptedUnlock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "AcceptedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AcceptedUnlockClean is a Record type
type AcceptedUnlockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts AcceptedUnlockClean to a map for DAML arguments
func (t AcceptedUnlockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AcceptedUnlockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlockClean to hex string (Canton MCMS format)
func (t AcceptedUnlockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlockClean from hex string (Canton MCMS format)
func (t *AcceptedUnlockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedUnlockExecute is a Record type
type AcceptedUnlockExecute struct {
	InstrumentConfigurationCid types.CONTRACT_ID   `json:"instrumentConfigurationCid"`
	CredentialCids             []types.CONTRACT_ID `json:"credentialCids"`
	HoldingCids                []types.CONTRACT_ID `json:"holdingCids"`
}

// ToMap converts AcceptedUnlockExecute to a map for DAML arguments
func (t AcceptedUnlockExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t AcceptedUnlockExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlockExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlockExecute to hex string (Canton MCMS format)
func (t AcceptedUnlockExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlockExecute from hex string (Canton MCMS format)
func (t *AcceptedUnlockExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedUnlockExecuteResult is a Record type
type AcceptedUnlockExecuteResult struct {
	HoldingCid        types.CONTRACT_ID   `json:"holdingCid"`
	RemainingCids     []types.CONTRACT_ID `json:"remainingCids"`
	ExecutedUnlockCid types.CONTRACT_ID   `json:"executedUnlockCid"`
}

// ToMap converts AcceptedUnlockExecuteResult to a map for DAML arguments
func (t AcceptedUnlockExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	m["remainingCids"] = func() []any {
		res := make([]any, 0, len(t.RemainingCids))
		for _, e := range t.RemainingCids {
			res = append(res, e)
		}
		return res
	}()

	m["executedUnlockCid"] = model.NestedToDAMLValue(t.ExecutedUnlockCid)

	return m
}

func (t AcceptedUnlockExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlockExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlockExecuteResult to hex string (Canton MCMS format)
func (t AcceptedUnlockExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlockExecuteResult from hex string (Canton MCMS format)
func (t *AcceptedUnlockExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedUnlockFail is a Record type
type AcceptedUnlockFail struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts AcceptedUnlockFail to a map for DAML arguments
func (t AcceptedUnlockFail) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t AcceptedUnlockFail) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlockFail) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlockFail to hex string (Canton MCMS format)
func (t AcceptedUnlockFail) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlockFail from hex string (Canton MCMS format)
func (t *AcceptedUnlockFail) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AcceptedUnlockFailResult is a Record type
type AcceptedUnlockFailResult struct {
	FailedUnlockCid types.CONTRACT_ID `json:"failedUnlockCid"`
}

// ToMap converts AcceptedUnlockFailResult to a map for DAML arguments
func (t AcceptedUnlockFailResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["failedUnlockCid"] = model.NestedToDAMLValue(t.FailedUnlockCid)

	return m
}

func (t AcceptedUnlockFailResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AcceptedUnlockFailResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AcceptedUnlockFailResult to hex string (Canton MCMS format)
func (t AcceptedUnlockFailResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AcceptedUnlockFailResult from hex string (Canton MCMS format)
func (t *AcceptedUnlockFailResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AppRewardConfiguration is a Template type
type AppRewardConfiguration struct {
	Operator types.PARTY                   `json:"operator"`
	Provider types.PARTY                   `json:"provider"`
	Details  AppRewardConfigurationDetails `json:"details"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AppRewardConfiguration) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AppRewardConfiguration) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AppRewardConfiguration) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["details"] = model.NestedToDAMLValue(t.Details)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AppRewardConfiguration) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["details"] = model.NestedToDAMLValue(t.Details)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AppRewardConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AppRewardConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AppRewardConfiguration to hex string (Canton MCMS format)
func (t AppRewardConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AppRewardConfiguration from hex string (Canton MCMS format)
func (t *AppRewardConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AppRewardConfiguration

// Archive exercises the Archive choice on this AppRewardConfiguration contract
// This method uses the package name in the template ID
func (t AppRewardConfiguration) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AppRewardConfiguration) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AppRewardConfigurationModify exercises the AppRewardConfiguration_Modify choice on this AppRewardConfiguration contract
// This method uses the package name in the template ID
func (t AppRewardConfiguration) AppRewardConfigurationModify(contractID string, args AppRewardConfigurationModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration"),
		ContractID: contractID,
		Choice:     "AppRewardConfiguration_Modify",
		Arguments:  argsToMap(args),
	}
}

// AppRewardConfigurationModifyWithPackageID exercises the AppRewardConfiguration_Modify choice using the provided package ID instead of package name
func (t AppRewardConfiguration) AppRewardConfigurationModifyWithPackageID(contractID string, packageID string, args AppRewardConfigurationModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.AppReward", "AppRewardConfiguration"),
		ContractID: contractID,
		Choice:     "AppRewardConfiguration_Modify",
		Arguments:  argsToMap(args),
	}
}

// AppRewardConfigurationDetails is a Record type
type AppRewardConfigurationDetails struct {
	Dso                          types.PARTY                                     `json:"dso"`
	OperatorAppRewardBeneficiary splice_api_featured_app_v1.AppRewardBeneficiary `json:"operatorAppRewardBeneficiary"`
}

// ToMap converts AppRewardConfigurationDetails to a map for DAML arguments
func (t AppRewardConfigurationDetails) ToMap() map[string]any {
	m := make(map[string]any)

	m["dso"] = t.Dso.ToMap()

	m["operatorAppRewardBeneficiary"] = model.NestedToDAMLValue(t.OperatorAppRewardBeneficiary)

	return m
}

func (t AppRewardConfigurationDetails) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AppRewardConfigurationDetails) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AppRewardConfigurationDetails to hex string (Canton MCMS format)
func (t AppRewardConfigurationDetails) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AppRewardConfigurationDetails from hex string (Canton MCMS format)
func (t *AppRewardConfigurationDetails) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AppRewardConfigurationModify is a Record type
type AppRewardConfigurationModify struct {
	Details AppRewardConfigurationDetails `json:"details"`
}

// ToMap converts AppRewardConfigurationModify to a map for DAML arguments
func (t AppRewardConfigurationModify) ToMap() map[string]any {
	m := make(map[string]any)

	m["details"] = model.NestedToDAMLValue(t.Details)

	return m
}

func (t AppRewardConfigurationModify) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AppRewardConfigurationModify) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AppRewardConfigurationModify to hex string (Canton MCMS format)
func (t AppRewardConfigurationModify) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AppRewardConfigurationModify from hex string (Canton MCMS format)
func (t *AppRewardConfigurationModify) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Batch is a Record type
type Batch struct {
	Id              types.TEXT       `json:"id"`
	Size            types.INT64      `json:"size"`
	SettlementFrom  *types.TIMESTAMP `json:"settlementFrom" hex:"optional"`
	SettlementUntil *types.TIMESTAMP `json:"settlementUntil" hex:"optional"`
}

// ToMap converts Batch to a map for DAML arguments
func (t Batch) ToMap() map[string]any {
	m := make(map[string]any)

	m["id"] = string(t.Id)

	m["size"] = int64(t.Size)

	if t.SettlementFrom != nil {
		m["settlementFrom"] = map[string]any{
			"_type": "optional",
			"value": *t.SettlementFrom,
		}
	} else {
		m["settlementFrom"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.SettlementUntil != nil {
		m["settlementUntil"] = map[string]any{
			"_type": "optional",
			"value": *t.SettlementUntil,
		}
	} else {
		m["settlementUntil"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t Batch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Batch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Batch to hex string (Canton MCMS format)
func (t Batch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Batch from hex string (Canton MCMS format)
func (t *Batch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Burn is a Record type
type Burn struct {
	Operator             types.PARTY                              `json:"operator"`
	Provider             types.PARTY                              `json:"provider"`
	Registrar            types.PARTY                              `json:"registrar"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Holder               types.PARTY                              `json:"holder"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                Batch                                    `json:"batch"`
}

// ToMap converts Burn to a map for DAML arguments
func (t Burn) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t Burn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Burn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Burn to hex string (Canton MCMS format)
func (t Burn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Burn from hex string (Canton MCMS format)
func (t *Burn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOffer is a Template type
type BurnOffer struct {
	Burn Burn `json:"burn"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t BurnOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t BurnOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t BurnOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t BurnOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t BurnOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOffer to hex string (Canton MCMS format)
func (t BurnOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOffer from hex string (Canton MCMS format)
func (t *BurnOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for BurnOffer

// BurnOfferClean exercises the BurnOffer_Clean choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferClean(contractID string, args BurnOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferCleanWithPackageID exercises the BurnOffer_Clean choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferCleanWithPackageID(contractID string, packageID string, args BurnOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferAccept exercises the BurnOffer_Accept choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferAccept(contractID string, args BurnOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferAcceptWithPackageID exercises the BurnOffer_Accept choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferAcceptWithPackageID(contractID string, packageID string, args BurnOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferReject exercises the BurnOffer_Reject choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferReject(contractID string, args BurnOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferRejectWithPackageID exercises the BurnOffer_Reject choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferRejectWithPackageID(contractID string, packageID string, args BurnOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferCancel exercises the BurnOffer_Cancel choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferCancel(contractID string, args BurnOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferCancelWithPackageID exercises the BurnOffer_Cancel choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferCancelWithPackageID(contractID string, packageID string, args BurnOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t BurnOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// BurnOfferAccept is a Record type
type BurnOfferAccept struct {
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// ToMap converts BurnOfferAccept to a map for DAML arguments
func (t BurnOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingLabel"] = string(t.HoldingLabel)

	return m
}

func (t BurnOfferAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferAccept to hex string (Canton MCMS format)
func (t BurnOfferAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferAccept from hex string (Canton MCMS format)
func (t *BurnOfferAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferAcceptResult is a Record type
type BurnOfferAcceptResult struct {
	AcceptedBurnCid types.CONTRACT_ID `json:"acceptedBurnCid"`
}

// ToMap converts BurnOfferAcceptResult to a map for DAML arguments
func (t BurnOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedBurnCid"] = model.NestedToDAMLValue(t.AcceptedBurnCid)

	return m
}

func (t BurnOfferAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferAcceptResult to hex string (Canton MCMS format)
func (t BurnOfferAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferAcceptResult from hex string (Canton MCMS format)
func (t *BurnOfferAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferCancel is a Record type
type BurnOfferCancel struct {
}

// ToMap converts BurnOfferCancel to a map for DAML arguments
func (t BurnOfferCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t BurnOfferCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferCancel to hex string (Canton MCMS format)
func (t BurnOfferCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferCancel from hex string (Canton MCMS format)
func (t *BurnOfferCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferCancelResult is a Record type
type BurnOfferCancelResult struct {
}

// ToMap converts BurnOfferCancelResult to a map for DAML arguments
func (t BurnOfferCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t BurnOfferCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferCancelResult to hex string (Canton MCMS format)
func (t BurnOfferCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferCancelResult from hex string (Canton MCMS format)
func (t *BurnOfferCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferClean is a Record type
type BurnOfferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts BurnOfferClean to a map for DAML arguments
func (t BurnOfferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t BurnOfferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferClean to hex string (Canton MCMS format)
func (t BurnOfferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferClean from hex string (Canton MCMS format)
func (t *BurnOfferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferReject is a Record type
type BurnOfferReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts BurnOfferReject to a map for DAML arguments
func (t BurnOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t BurnOfferReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferReject to hex string (Canton MCMS format)
func (t BurnOfferReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferReject from hex string (Canton MCMS format)
func (t *BurnOfferReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnOfferRejectResult is a Record type
type BurnOfferRejectResult struct {
	RejectedBurnCid types.CONTRACT_ID `json:"rejectedBurnCid"`
}

// ToMap converts BurnOfferRejectResult to a map for DAML arguments
func (t BurnOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedBurnCid"] = model.NestedToDAMLValue(t.RejectedBurnCid)

	return m
}

func (t BurnOfferRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnOfferRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnOfferRejectResult to hex string (Canton MCMS format)
func (t BurnOfferRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnOfferRejectResult from hex string (Canton MCMS format)
func (t *BurnOfferRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequest is a Template type
type BurnRequest struct {
	Burn         Burn       `json:"burn"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t BurnRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t BurnRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t BurnRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t BurnRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t BurnRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequest to hex string (Canton MCMS format)
func (t BurnRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequest from hex string (Canton MCMS format)
func (t *BurnRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for BurnRequest

// BurnRequestClean exercises the BurnRequest_Clean choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestClean(contractID string, args BurnRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestCleanWithPackageID exercises the BurnRequest_Clean choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestCleanWithPackageID(contractID string, packageID string, args BurnRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestAccept exercises the BurnRequest_Accept choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestAccept(contractID string, args BurnRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestAcceptWithPackageID exercises the BurnRequest_Accept choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestAcceptWithPackageID(contractID string, packageID string, args BurnRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestReject exercises the BurnRequest_Reject choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestReject(contractID string, args BurnRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestRejectWithPackageID exercises the BurnRequest_Reject choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestRejectWithPackageID(contractID string, packageID string, args BurnRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestCancel exercises the BurnRequest_Cancel choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestCancel(contractID string, args BurnRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestCancelWithPackageID exercises the BurnRequest_Cancel choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestCancelWithPackageID(contractID string, packageID string, args BurnRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t BurnRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// BurnRequestAccept is a Record type
type BurnRequestAccept struct {
}

// ToMap converts BurnRequestAccept to a map for DAML arguments
func (t BurnRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t BurnRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestAccept to hex string (Canton MCMS format)
func (t BurnRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestAccept from hex string (Canton MCMS format)
func (t *BurnRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestAcceptResult is a Record type
type BurnRequestAcceptResult struct {
	AcceptedBurnCid types.CONTRACT_ID `json:"acceptedBurnCid"`
}

// ToMap converts BurnRequestAcceptResult to a map for DAML arguments
func (t BurnRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedBurnCid"] = model.NestedToDAMLValue(t.AcceptedBurnCid)

	return m
}

func (t BurnRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestAcceptResult to hex string (Canton MCMS format)
func (t BurnRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestAcceptResult from hex string (Canton MCMS format)
func (t *BurnRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestCancel is a Record type
type BurnRequestCancel struct {
}

// ToMap converts BurnRequestCancel to a map for DAML arguments
func (t BurnRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t BurnRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestCancel to hex string (Canton MCMS format)
func (t BurnRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestCancel from hex string (Canton MCMS format)
func (t *BurnRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestCancelResult is a Record type
type BurnRequestCancelResult struct {
}

// ToMap converts BurnRequestCancelResult to a map for DAML arguments
func (t BurnRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t BurnRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestCancelResult to hex string (Canton MCMS format)
func (t BurnRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestCancelResult from hex string (Canton MCMS format)
func (t *BurnRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestClean is a Record type
type BurnRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts BurnRequestClean to a map for DAML arguments
func (t BurnRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t BurnRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestClean to hex string (Canton MCMS format)
func (t BurnRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestClean from hex string (Canton MCMS format)
func (t *BurnRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestReject is a Record type
type BurnRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts BurnRequestReject to a map for DAML arguments
func (t BurnRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t BurnRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestReject to hex string (Canton MCMS format)
func (t BurnRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestReject from hex string (Canton MCMS format)
func (t *BurnRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BurnRequestRejectResult is a Record type
type BurnRequestRejectResult struct {
	RejectedBurnCid types.CONTRACT_ID `json:"rejectedBurnCid"`
}

// ToMap converts BurnRequestRejectResult to a map for DAML arguments
func (t BurnRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedBurnCid"] = model.NestedToDAMLValue(t.RejectedBurnCid)

	return m
}

func (t BurnRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BurnRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BurnRequestRejectResult to hex string (Canton MCMS format)
func (t BurnRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BurnRequestRejectResult from hex string (Canton MCMS format)
func (t *BurnRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DvpLegAllocation is a Template type
type DvpLegAllocation struct {
	Allocation       splice_api_token_allocation_v1.AllocationSpecification `json:"allocation"`
	LockedHoldingCid types.CONTRACT_ID                                      `json:"lockedHoldingCid"`
	Operator         types.PARTY                                            `json:"operator"`
	Provider         *types.PARTY                                           `json:"provider" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t DvpLegAllocation) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Allocation", "DvpLegAllocation")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t DvpLegAllocation) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Allocation", "DvpLegAllocation")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t DvpLegAllocation) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allocation"] = model.NestedToDAMLValue(t.Allocation)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedHoldingCid"] = model.NestedToDAMLValue(t.LockedHoldingCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	if t.Provider != nil {
		args["provider"] = map[string]any{
			"_type": "optional",
			"value": (*t.Provider).ToMap(),
		}
	} else {
		args["provider"] = map[string]any{
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
func (t DvpLegAllocation) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["allocation"] = model.NestedToDAMLValue(t.Allocation)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedHoldingCid"] = model.NestedToDAMLValue(t.LockedHoldingCid)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	if t.Provider != nil {
		args["provider"] = map[string]any{
			"_type": "optional",
			"value": (*t.Provider).ToMap(),
		}
	} else {
		args["provider"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t DvpLegAllocation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DvpLegAllocation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DvpLegAllocation to hex string (Canton MCMS format)
func (t DvpLegAllocation) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DvpLegAllocation from hex string (Canton MCMS format)
func (t *DvpLegAllocation) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for DvpLegAllocation

// Archive exercises the Archive choice on this DvpLegAllocation contract via the IAllocation interface
// This method uses the package name in the template ID
func (t DvpLegAllocation) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t DvpLegAllocation) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// AllocationWithdraw exercises the Allocation_Withdraw choice on this DvpLegAllocation contract via the IAllocation interface
// This method uses the package name in the template ID
func (t DvpLegAllocation) AllocationWithdraw(contractID string, args splice_api_token_allocation_v1.AllocationWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// AllocationWithdrawWithPackageID exercises the Allocation_Withdraw choice using the provided package ID instead of package name
func (t DvpLegAllocation) AllocationWithdrawWithPackageID(contractID string, packageID string, args splice_api_token_allocation_v1.AllocationWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// AllocationCancel exercises the Allocation_Cancel choice on this DvpLegAllocation contract via the IAllocation interface
// This method uses the package name in the template ID
func (t DvpLegAllocation) AllocationCancel(contractID string, args splice_api_token_allocation_v1.AllocationCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Cancel",
		Arguments:  argsToMap(args),
	}
}

// AllocationCancelWithPackageID exercises the Allocation_Cancel choice using the provided package ID instead of package name
func (t DvpLegAllocation) AllocationCancelWithPackageID(contractID string, packageID string, args splice_api_token_allocation_v1.AllocationCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_Cancel",
		Arguments:  argsToMap(args),
	}
}

// AllocationExecuteTransfer exercises the Allocation_ExecuteTransfer choice on this DvpLegAllocation contract via the IAllocation interface
// This method uses the package name in the template ID
func (t DvpLegAllocation) AllocationExecuteTransfer(contractID string, args splice_api_token_allocation_v1.AllocationExecuteTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_ExecuteTransfer",
		Arguments:  argsToMap(args),
	}
}

// AllocationExecuteTransferWithPackageID exercises the Allocation_ExecuteTransfer choice using the provided package ID instead of package name
func (t DvpLegAllocation) AllocationExecuteTransferWithPackageID(contractID string, packageID string, args splice_api_token_allocation_v1.AllocationExecuteTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Allocation", "Allocation"),
		ContractID: contractID,
		Choice:     "Allocation_ExecuteTransfer",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for DvpLegAllocation

var _ splice_api_token_allocation_v1.IAllocation = (*DvpLegAllocation)(nil)

// ExecutedBurn is a Template type
type ExecutedBurn struct {
	Burn         Burn       `json:"burn"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedBurn to hex string (Canton MCMS format)
func (t ExecutedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedBurn from hex string (Canton MCMS format)
func (t *ExecutedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedBurn

// ExecutedBurnClean exercises the ExecutedBurn_Clean choice on this ExecutedBurn contract
// This method uses the package name in the template ID
func (t ExecutedBurn) ExecutedBurnClean(contractID string, args ExecutedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// ExecutedBurnCleanWithPackageID exercises the ExecutedBurn_Clean choice using the provided package ID instead of package name
func (t ExecutedBurn) ExecutedBurnCleanWithPackageID(contractID string, packageID string, args ExecutedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedBurn contract
// This method uses the package name in the template ID
func (t ExecutedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedBurnDelete exercises the ExecutedBurn_Delete choice on this ExecutedBurn contract
// This method uses the package name in the template ID
func (t ExecutedBurn) ExecutedBurnDelete(contractID string, args ExecutedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedBurnDeleteWithPackageID exercises the ExecutedBurn_Delete choice using the provided package ID instead of package name
func (t ExecutedBurn) ExecutedBurnDeleteWithPackageID(contractID string, packageID string, args ExecutedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedBurnClean is a Record type
type ExecutedBurnClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedBurnClean to a map for DAML arguments
func (t ExecutedBurnClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t ExecutedBurnClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedBurnClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedBurnClean to hex string (Canton MCMS format)
func (t ExecutedBurnClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedBurnClean from hex string (Canton MCMS format)
func (t *ExecutedBurnClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedBurnDelete is a Record type
type ExecutedBurnDelete struct {
}

// ToMap converts ExecutedBurnDelete to a map for DAML arguments
func (t ExecutedBurnDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedBurnDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedBurnDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedBurnDelete to hex string (Canton MCMS format)
func (t ExecutedBurnDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedBurnDelete from hex string (Canton MCMS format)
func (t *ExecutedBurnDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedBurnDeleteResult is a Record type
type ExecutedBurnDeleteResult struct {
}

// ToMap converts ExecutedBurnDeleteResult to a map for DAML arguments
func (t ExecutedBurnDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedBurnDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedBurnDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedBurnDeleteResult to hex string (Canton MCMS format)
func (t ExecutedBurnDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedBurnDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedBurnDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedForceTransfer is a Template type
type ExecutedForceTransfer struct {
	ForceTransfer      ForceTransfer `json:"forceTransfer"`
	RegistrarRationale types.TEXT    `json:"registrarRationale"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedForceTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedForceTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedForceTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRationale"] = string(t.RegistrarRationale)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutedForceTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRationale"] = string(t.RegistrarRationale)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedForceTransfer to hex string (Canton MCMS format)
func (t ExecutedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedForceTransfer from hex string (Canton MCMS format)
func (t *ExecutedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedForceTransfer

// Archive exercises the Archive choice on this ExecutedForceTransfer contract
// This method uses the package name in the template ID
func (t ExecutedForceTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedForceTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedForceTransferDelete exercises the ExecutedForceTransfer_Delete choice on this ExecutedForceTransfer contract
// This method uses the package name in the template ID
func (t ExecutedForceTransfer) ExecutedForceTransferDelete(contractID string, args ExecutedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer"),
		ContractID: contractID,
		Choice:     "ExecutedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedForceTransferDeleteWithPackageID exercises the ExecutedForceTransfer_Delete choice using the provided package ID instead of package name
func (t ExecutedForceTransfer) ExecutedForceTransferDeleteWithPackageID(contractID string, packageID string, args ExecutedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ExecutedForceTransfer"),
		ContractID: contractID,
		Choice:     "ExecutedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedForceTransferDelete is a Record type
type ExecutedForceTransferDelete struct {
}

// ToMap converts ExecutedForceTransferDelete to a map for DAML arguments
func (t ExecutedForceTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedForceTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedForceTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedForceTransferDelete to hex string (Canton MCMS format)
func (t ExecutedForceTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedForceTransferDelete from hex string (Canton MCMS format)
func (t *ExecutedForceTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedForceTransferDeleteResult is a Record type
type ExecutedForceTransferDeleteResult struct {
}

// ToMap converts ExecutedForceTransferDeleteResult to a map for DAML arguments
func (t ExecutedForceTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedForceTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedForceTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedForceTransferDeleteResult to hex string (Canton MCMS format)
func (t ExecutedForceTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedForceTransferDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedForceTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedLock is a Template type
type ExecutedLock struct {
	Lock         Lock3      `json:"lock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedLock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "ExecutedLock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedLock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "ExecutedLock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedLock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutedLock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedLock to hex string (Canton MCMS format)
func (t ExecutedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedLock from hex string (Canton MCMS format)
func (t *ExecutedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedLock

// ExecutedLockClean exercises the ExecutedLock_Clean choice on this ExecutedLock contract
// This method uses the package name in the template ID
func (t ExecutedLock) ExecutedLockClean(contractID string, args ExecutedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "ExecutedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// ExecutedLockCleanWithPackageID exercises the ExecutedLock_Clean choice using the provided package ID instead of package name
func (t ExecutedLock) ExecutedLockCleanWithPackageID(contractID string, packageID string, args ExecutedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "ExecutedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedLock contract
// This method uses the package name in the template ID
func (t ExecutedLock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedLock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedLockDelete exercises the ExecutedLock_Delete choice on this ExecutedLock contract
// This method uses the package name in the template ID
func (t ExecutedLock) ExecutedLockDelete(contractID string, args ExecutedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "ExecutedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedLockDeleteWithPackageID exercises the ExecutedLock_Delete choice using the provided package ID instead of package name
func (t ExecutedLock) ExecutedLockDeleteWithPackageID(contractID string, packageID string, args ExecutedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "ExecutedLock"),
		ContractID: contractID,
		Choice:     "ExecutedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedLockClean is a Record type
type ExecutedLockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedLockClean to a map for DAML arguments
func (t ExecutedLockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t ExecutedLockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedLockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedLockClean to hex string (Canton MCMS format)
func (t ExecutedLockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedLockClean from hex string (Canton MCMS format)
func (t *ExecutedLockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedLockDelete is a Record type
type ExecutedLockDelete struct {
}

// ToMap converts ExecutedLockDelete to a map for DAML arguments
func (t ExecutedLockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedLockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedLockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedLockDelete to hex string (Canton MCMS format)
func (t ExecutedLockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedLockDelete from hex string (Canton MCMS format)
func (t *ExecutedLockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedLockDeleteResult is a Record type
type ExecutedLockDeleteResult struct {
}

// ToMap converts ExecutedLockDeleteResult to a map for DAML arguments
func (t ExecutedLockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedLockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedLockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedLockDeleteResult to hex string (Canton MCMS format)
func (t ExecutedLockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedLockDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedLockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedMint is a Template type
type ExecutedMint struct {
	Mint         Mint       `json:"mint"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "ExecutedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "ExecutedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedMint to hex string (Canton MCMS format)
func (t ExecutedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedMint from hex string (Canton MCMS format)
func (t *ExecutedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedMint

// ExecutedMintClean exercises the ExecutedMint_Clean choice on this ExecutedMint contract
// This method uses the package name in the template ID
func (t ExecutedMint) ExecutedMintClean(contractID string, args ExecutedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// ExecutedMintCleanWithPackageID exercises the ExecutedMint_Clean choice using the provided package ID instead of package name
func (t ExecutedMint) ExecutedMintCleanWithPackageID(contractID string, packageID string, args ExecutedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedMint contract
// This method uses the package name in the template ID
func (t ExecutedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedMintDelete exercises the ExecutedMint_Delete choice on this ExecutedMint contract
// This method uses the package name in the template ID
func (t ExecutedMint) ExecutedMintDelete(contractID string, args ExecutedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedMintDeleteWithPackageID exercises the ExecutedMint_Delete choice using the provided package ID instead of package name
func (t ExecutedMint) ExecutedMintDeleteWithPackageID(contractID string, packageID string, args ExecutedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedMintClean is a Record type
type ExecutedMintClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedMintClean to a map for DAML arguments
func (t ExecutedMintClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t ExecutedMintClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedMintClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedMintClean to hex string (Canton MCMS format)
func (t ExecutedMintClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedMintClean from hex string (Canton MCMS format)
func (t *ExecutedMintClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedMintDelete is a Record type
type ExecutedMintDelete struct {
}

// ToMap converts ExecutedMintDelete to a map for DAML arguments
func (t ExecutedMintDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedMintDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedMintDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedMintDelete to hex string (Canton MCMS format)
func (t ExecutedMintDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedMintDelete from hex string (Canton MCMS format)
func (t *ExecutedMintDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedMintDeleteResult is a Record type
type ExecutedMintDeleteResult struct {
}

// ToMap converts ExecutedMintDeleteResult to a map for DAML arguments
func (t ExecutedMintDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedMintDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedMintDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedMintDeleteResult to hex string (Canton MCMS format)
func (t ExecutedMintDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedMintDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedMintDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedTransfer is a Template type
type ExecutedTransfer struct {
	Transfer           Transfer2   `json:"transfer"`
	SenderLabel        types.TEXT  `json:"senderLabel"`
	ReceiverLabel      types.TEXT  `json:"receiverLabel"`
	Observers          *types.SET  `json:"observers" hex:"optional"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OperatorIsObserver != nil {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.OperatorIsObserver),
		}
	} else {
		args["operatorIsObserver"] = map[string]any{
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
func (t ExecutedTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OperatorIsObserver != nil {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.OperatorIsObserver),
		}
	} else {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedTransfer to hex string (Canton MCMS format)
func (t ExecutedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedTransfer from hex string (Canton MCMS format)
func (t *ExecutedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedTransfer

// Archive exercises the Archive choice on this ExecutedTransfer contract
// This method uses the package name in the template ID
func (t ExecutedTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedTransferDelete exercises the ExecutedTransfer_Delete choice on this ExecutedTransfer contract
// This method uses the package name in the template ID
func (t ExecutedTransfer) ExecutedTransferDelete(contractID string, args ExecutedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer"),
		ContractID: contractID,
		Choice:     "ExecutedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedTransferDeleteWithPackageID exercises the ExecutedTransfer_Delete choice using the provided package ID instead of package name
func (t ExecutedTransfer) ExecutedTransferDeleteWithPackageID(contractID string, packageID string, args ExecutedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "ExecutedTransfer"),
		ContractID: contractID,
		Choice:     "ExecutedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedTransferDelete is a Record type
type ExecutedTransferDelete struct {
}

// ToMap converts ExecutedTransferDelete to a map for DAML arguments
func (t ExecutedTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedTransferDelete to hex string (Canton MCMS format)
func (t ExecutedTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedTransferDelete from hex string (Canton MCMS format)
func (t *ExecutedTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedTransferDeleteResult is a Record type
type ExecutedTransferDeleteResult struct {
}

// ToMap converts ExecutedTransferDeleteResult to a map for DAML arguments
func (t ExecutedTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedTransferDeleteResult to hex string (Canton MCMS format)
func (t ExecutedTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedTransferDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedUnlock is a Template type
type ExecutedUnlock struct {
	Unlock       Unlock     `json:"unlock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedUnlock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedUnlock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedUnlock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ExecutedUnlock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ExecutedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedUnlock to hex string (Canton MCMS format)
func (t ExecutedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedUnlock from hex string (Canton MCMS format)
func (t *ExecutedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ExecutedUnlock

// ExecutedUnlockClean exercises the ExecutedUnlock_Clean choice on this ExecutedUnlock contract
// This method uses the package name in the template ID
func (t ExecutedUnlock) ExecutedUnlockClean(contractID string, args ExecutedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "ExecutedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// ExecutedUnlockCleanWithPackageID exercises the ExecutedUnlock_Clean choice using the provided package ID instead of package name
func (t ExecutedUnlock) ExecutedUnlockCleanWithPackageID(contractID string, packageID string, args ExecutedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "ExecutedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedUnlock contract
// This method uses the package name in the template ID
func (t ExecutedUnlock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedUnlock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedUnlockDelete exercises the ExecutedUnlock_Delete choice on this ExecutedUnlock contract
// This method uses the package name in the template ID
func (t ExecutedUnlock) ExecutedUnlockDelete(contractID string, args ExecutedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "ExecutedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedUnlockDeleteWithPackageID exercises the ExecutedUnlock_Delete choice using the provided package ID instead of package name
func (t ExecutedUnlock) ExecutedUnlockDeleteWithPackageID(contractID string, packageID string, args ExecutedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "ExecutedUnlock"),
		ContractID: contractID,
		Choice:     "ExecutedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedUnlockClean is a Record type
type ExecutedUnlockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedUnlockClean to a map for DAML arguments
func (t ExecutedUnlockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t ExecutedUnlockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedUnlockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedUnlockClean to hex string (Canton MCMS format)
func (t ExecutedUnlockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedUnlockClean from hex string (Canton MCMS format)
func (t *ExecutedUnlockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedUnlockDelete is a Record type
type ExecutedUnlockDelete struct {
}

// ToMap converts ExecutedUnlockDelete to a map for DAML arguments
func (t ExecutedUnlockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedUnlockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedUnlockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedUnlockDelete to hex string (Canton MCMS format)
func (t ExecutedUnlockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedUnlockDelete from hex string (Canton MCMS format)
func (t *ExecutedUnlockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedUnlockDeleteResult is a Record type
type ExecutedUnlockDeleteResult struct {
}

// ToMap converts ExecutedUnlockDeleteResult to a map for DAML arguments
func (t ExecutedUnlockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ExecutedUnlockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecutedUnlockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecutedUnlockDeleteResult to hex string (Canton MCMS format)
func (t ExecutedUnlockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecutedUnlockDeleteResult from hex string (Canton MCMS format)
func (t *ExecutedUnlockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExpectedInputHoldingLockState is an enum type
type ExpectedInputHoldingLockState string

const (
	ExpectedInputHoldingLockStateExpectedUnlocked ExpectedInputHoldingLockState = "ExpectedUnlocked"

	ExpectedInputHoldingLockStateExpectedLocked ExpectedInputHoldingLockState = "ExpectedLocked"
)

func (e ExpectedInputHoldingLockState) GetEnumConstructor() string { return string(e) }

func (e ExpectedInputHoldingLockState) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "ExpectedInputHoldingLockState")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e ExpectedInputHoldingLockState) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "ExpectedInputHoldingLockState")
}

func (e ExpectedInputHoldingLockState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *ExpectedInputHoldingLockState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes ExpectedInputHoldingLockState to hex string (Canton MCMS format)
func (e ExpectedInputHoldingLockState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes ExpectedInputHoldingLockState from hex string (Canton MCMS format)
func (e *ExpectedInputHoldingLockState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = ExpectedInputHoldingLockState("")

// FailedBurn is a Template type
type FailedBurn struct {
	Burn   Burn       `json:"burn"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "FailedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "FailedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FailedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedBurn to hex string (Canton MCMS format)
func (t FailedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedBurn from hex string (Canton MCMS format)
func (t *FailedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedBurn

// FailedBurnClean exercises the FailedBurn_Clean choice on this FailedBurn contract
// This method uses the package name in the template ID
func (t FailedBurn) FailedBurnClean(contractID string, args FailedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "FailedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedBurnCleanWithPackageID exercises the FailedBurn_Clean choice using the provided package ID instead of package name
func (t FailedBurn) FailedBurnCleanWithPackageID(contractID string, packageID string, args FailedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "FailedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedBurnDelete exercises the FailedBurn_Delete choice on this FailedBurn contract
// This method uses the package name in the template ID
func (t FailedBurn) FailedBurnDelete(contractID string, args FailedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "FailedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedBurnDeleteWithPackageID exercises the FailedBurn_Delete choice using the provided package ID instead of package name
func (t FailedBurn) FailedBurnDeleteWithPackageID(contractID string, packageID string, args FailedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "FailedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedBurn contract
// This method uses the package name in the template ID
func (t FailedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "FailedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedBurnClean is a Record type
type FailedBurnClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts FailedBurnClean to a map for DAML arguments
func (t FailedBurnClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t FailedBurnClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedBurnClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedBurnClean to hex string (Canton MCMS format)
func (t FailedBurnClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedBurnClean from hex string (Canton MCMS format)
func (t *FailedBurnClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedBurnDelete is a Record type
type FailedBurnDelete struct {
}

// ToMap converts FailedBurnDelete to a map for DAML arguments
func (t FailedBurnDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedBurnDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedBurnDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedBurnDelete to hex string (Canton MCMS format)
func (t FailedBurnDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedBurnDelete from hex string (Canton MCMS format)
func (t *FailedBurnDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedBurnDeleteResult is a Record type
type FailedBurnDeleteResult struct {
}

// ToMap converts FailedBurnDeleteResult to a map for DAML arguments
func (t FailedBurnDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedBurnDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedBurnDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedBurnDeleteResult to hex string (Canton MCMS format)
func (t FailedBurnDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedBurnDeleteResult from hex string (Canton MCMS format)
func (t *FailedBurnDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedForceTransfer is a Template type
type FailedForceTransfer struct {
	ForceTransfer ForceTransfer `json:"forceTransfer"`
	Reason        types.TEXT    `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedForceTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedForceTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedForceTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FailedForceTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedForceTransfer to hex string (Canton MCMS format)
func (t FailedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedForceTransfer from hex string (Canton MCMS format)
func (t *FailedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedForceTransfer

// FailedForceTransferDelete exercises the FailedForceTransfer_Delete choice on this FailedForceTransfer contract
// This method uses the package name in the template ID
func (t FailedForceTransfer) FailedForceTransferDelete(contractID string, args FailedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer"),
		ContractID: contractID,
		Choice:     "FailedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedForceTransferDeleteWithPackageID exercises the FailedForceTransfer_Delete choice using the provided package ID instead of package name
func (t FailedForceTransfer) FailedForceTransferDeleteWithPackageID(contractID string, packageID string, args FailedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer"),
		ContractID: contractID,
		Choice:     "FailedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedForceTransfer contract
// This method uses the package name in the template ID
func (t FailedForceTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedForceTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "FailedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedForceTransferDelete is a Record type
type FailedForceTransferDelete struct {
}

// ToMap converts FailedForceTransferDelete to a map for DAML arguments
func (t FailedForceTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedForceTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedForceTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedForceTransferDelete to hex string (Canton MCMS format)
func (t FailedForceTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedForceTransferDelete from hex string (Canton MCMS format)
func (t *FailedForceTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedForceTransferDeleteResult is a Record type
type FailedForceTransferDeleteResult struct {
}

// ToMap converts FailedForceTransferDeleteResult to a map for DAML arguments
func (t FailedForceTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedForceTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedForceTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedForceTransferDeleteResult to hex string (Canton MCMS format)
func (t FailedForceTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedForceTransferDeleteResult from hex string (Canton MCMS format)
func (t *FailedForceTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedLock is a Template type
type FailedLock struct {
	Lock   Lock3      `json:"lock"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedLock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "FailedLock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedLock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "FailedLock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedLock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FailedLock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedLock to hex string (Canton MCMS format)
func (t FailedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedLock from hex string (Canton MCMS format)
func (t *FailedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedLock

// FailedLockClean exercises the FailedLock_Clean choice on this FailedLock contract
// This method uses the package name in the template ID
func (t FailedLock) FailedLockClean(contractID string, args FailedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "FailedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedLockCleanWithPackageID exercises the FailedLock_Clean choice using the provided package ID instead of package name
func (t FailedLock) FailedLockCleanWithPackageID(contractID string, packageID string, args FailedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "FailedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedLockDelete exercises the FailedLock_Delete choice on this FailedLock contract
// This method uses the package name in the template ID
func (t FailedLock) FailedLockDelete(contractID string, args FailedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "FailedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedLockDeleteWithPackageID exercises the FailedLock_Delete choice using the provided package ID instead of package name
func (t FailedLock) FailedLockDeleteWithPackageID(contractID string, packageID string, args FailedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "FailedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedLock contract
// This method uses the package name in the template ID
func (t FailedLock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedLock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "FailedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedLockClean is a Record type
type FailedLockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts FailedLockClean to a map for DAML arguments
func (t FailedLockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t FailedLockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedLockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedLockClean to hex string (Canton MCMS format)
func (t FailedLockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedLockClean from hex string (Canton MCMS format)
func (t *FailedLockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedLockDelete is a Record type
type FailedLockDelete struct {
}

// ToMap converts FailedLockDelete to a map for DAML arguments
func (t FailedLockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedLockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedLockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedLockDelete to hex string (Canton MCMS format)
func (t FailedLockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedLockDelete from hex string (Canton MCMS format)
func (t *FailedLockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedLockDeleteResult is a Record type
type FailedLockDeleteResult struct {
}

// ToMap converts FailedLockDeleteResult to a map for DAML arguments
func (t FailedLockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedLockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedLockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedLockDeleteResult to hex string (Canton MCMS format)
func (t FailedLockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedLockDeleteResult from hex string (Canton MCMS format)
func (t *FailedLockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedMint is a Template type
type FailedMint struct {
	Mint   Mint       `json:"mint"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "FailedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "FailedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FailedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedMint to hex string (Canton MCMS format)
func (t FailedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedMint from hex string (Canton MCMS format)
func (t *FailedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedMint

// FailedMintClean exercises the FailedMint_Clean choice on this FailedMint contract
// This method uses the package name in the template ID
func (t FailedMint) FailedMintClean(contractID string, args FailedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "FailedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedMintCleanWithPackageID exercises the FailedMint_Clean choice using the provided package ID instead of package name
func (t FailedMint) FailedMintCleanWithPackageID(contractID string, packageID string, args FailedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "FailedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedMintDelete exercises the FailedMint_Delete choice on this FailedMint contract
// This method uses the package name in the template ID
func (t FailedMint) FailedMintDelete(contractID string, args FailedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "FailedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedMintDeleteWithPackageID exercises the FailedMint_Delete choice using the provided package ID instead of package name
func (t FailedMint) FailedMintDeleteWithPackageID(contractID string, packageID string, args FailedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "FailedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedMint contract
// This method uses the package name in the template ID
func (t FailedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "FailedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedMintClean is a Record type
type FailedMintClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts FailedMintClean to a map for DAML arguments
func (t FailedMintClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t FailedMintClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedMintClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedMintClean to hex string (Canton MCMS format)
func (t FailedMintClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedMintClean from hex string (Canton MCMS format)
func (t *FailedMintClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedMintDelete is a Record type
type FailedMintDelete struct {
}

// ToMap converts FailedMintDelete to a map for DAML arguments
func (t FailedMintDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedMintDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedMintDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedMintDelete to hex string (Canton MCMS format)
func (t FailedMintDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedMintDelete from hex string (Canton MCMS format)
func (t *FailedMintDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedMintDeleteResult is a Record type
type FailedMintDeleteResult struct {
}

// ToMap converts FailedMintDeleteResult to a map for DAML arguments
func (t FailedMintDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedMintDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedMintDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedMintDeleteResult to hex string (Canton MCMS format)
func (t FailedMintDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedMintDeleteResult from hex string (Canton MCMS format)
func (t *FailedMintDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedTransfer is a Template type
type FailedTransfer struct {
	Transfer  Transfer2  `json:"transfer"`
	Reason    types.TEXT `json:"reason"`
	Observers *types.SET `json:"observers" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
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
func (t FailedTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedTransfer to hex string (Canton MCMS format)
func (t FailedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedTransfer from hex string (Canton MCMS format)
func (t *FailedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedTransfer

// FailedTransferClean exercises the FailedTransfer_Clean choice on this FailedTransfer contract
// This method uses the package name in the template ID
func (t FailedTransfer) FailedTransferClean(contractID string, args FailedTransferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "FailedTransfer_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedTransferCleanWithPackageID exercises the FailedTransfer_Clean choice using the provided package ID instead of package name
func (t FailedTransfer) FailedTransferCleanWithPackageID(contractID string, packageID string, args FailedTransferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "FailedTransfer_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedTransferDelete exercises the FailedTransfer_Delete choice on this FailedTransfer contract
// This method uses the package name in the template ID
func (t FailedTransfer) FailedTransferDelete(contractID string, args FailedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "FailedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedTransferDeleteWithPackageID exercises the FailedTransfer_Delete choice using the provided package ID instead of package name
func (t FailedTransfer) FailedTransferDeleteWithPackageID(contractID string, packageID string, args FailedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "FailedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedTransfer contract
// This method uses the package name in the template ID
func (t FailedTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "FailedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedTransferClean is a Record type
type FailedTransferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts FailedTransferClean to a map for DAML arguments
func (t FailedTransferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t FailedTransferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedTransferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedTransferClean to hex string (Canton MCMS format)
func (t FailedTransferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedTransferClean from hex string (Canton MCMS format)
func (t *FailedTransferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedTransferDelete is a Record type
type FailedTransferDelete struct {
}

// ToMap converts FailedTransferDelete to a map for DAML arguments
func (t FailedTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedTransferDelete to hex string (Canton MCMS format)
func (t FailedTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedTransferDelete from hex string (Canton MCMS format)
func (t *FailedTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedTransferDeleteResult is a Record type
type FailedTransferDeleteResult struct {
}

// ToMap converts FailedTransferDeleteResult to a map for DAML arguments
func (t FailedTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedTransferDeleteResult to hex string (Canton MCMS format)
func (t FailedTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedTransferDeleteResult from hex string (Canton MCMS format)
func (t *FailedTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedUnlock is a Template type
type FailedUnlock struct {
	Unlock Unlock     `json:"unlock"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t FailedUnlock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t FailedUnlock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t FailedUnlock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t FailedUnlock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t FailedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedUnlock to hex string (Canton MCMS format)
func (t FailedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedUnlock from hex string (Canton MCMS format)
func (t *FailedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for FailedUnlock

// FailedUnlockClean exercises the FailedUnlock_Clean choice on this FailedUnlock contract
// This method uses the package name in the template ID
func (t FailedUnlock) FailedUnlockClean(contractID string, args FailedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "FailedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedUnlockCleanWithPackageID exercises the FailedUnlock_Clean choice using the provided package ID instead of package name
func (t FailedUnlock) FailedUnlockCleanWithPackageID(contractID string, packageID string, args FailedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "FailedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// FailedUnlockDelete exercises the FailedUnlock_Delete choice on this FailedUnlock contract
// This method uses the package name in the template ID
func (t FailedUnlock) FailedUnlockDelete(contractID string, args FailedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "FailedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// FailedUnlockDeleteWithPackageID exercises the FailedUnlock_Delete choice using the provided package ID instead of package name
func (t FailedUnlock) FailedUnlockDeleteWithPackageID(contractID string, packageID string, args FailedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "FailedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this FailedUnlock contract
// This method uses the package name in the template ID
func (t FailedUnlock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t FailedUnlock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "FailedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// FailedUnlockClean is a Record type
type FailedUnlockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts FailedUnlockClean to a map for DAML arguments
func (t FailedUnlockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t FailedUnlockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedUnlockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedUnlockClean to hex string (Canton MCMS format)
func (t FailedUnlockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedUnlockClean from hex string (Canton MCMS format)
func (t *FailedUnlockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedUnlockDelete is a Record type
type FailedUnlockDelete struct {
}

// ToMap converts FailedUnlockDelete to a map for DAML arguments
func (t FailedUnlockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedUnlockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedUnlockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedUnlockDelete to hex string (Canton MCMS format)
func (t FailedUnlockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedUnlockDelete from hex string (Canton MCMS format)
func (t *FailedUnlockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FailedUnlockDeleteResult is a Record type
type FailedUnlockDeleteResult struct {
}

// ToMap converts FailedUnlockDeleteResult to a map for DAML arguments
func (t FailedUnlockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t FailedUnlockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FailedUnlockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FailedUnlockDeleteResult to hex string (Canton MCMS format)
func (t FailedUnlockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FailedUnlockDeleteResult from hex string (Canton MCMS format)
func (t *FailedUnlockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransfer is a Record type
type ForceTransfer struct {
	Requestor          types.PARTY `json:"requestor"`
	RequestorRationale types.TEXT  `json:"requestorRationale"`
	Transfer           Transfer2   `json:"transfer"`
	SenderLabel        types.TEXT  `json:"senderLabel"`
	ReceiverLabel      types.TEXT  `json:"receiverLabel"`
}

// ToMap converts ForceTransfer to a map for DAML arguments
func (t ForceTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["requestor"] = t.Requestor.ToMap()

	m["requestorRationale"] = string(t.RequestorRationale)

	m["transfer"] = model.NestedToDAMLValue(t.Transfer)

	m["senderLabel"] = string(t.SenderLabel)

	m["receiverLabel"] = string(t.ReceiverLabel)

	return m
}

func (t ForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransfer to hex string (Canton MCMS format)
func (t ForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransfer from hex string (Canton MCMS format)
func (t *ForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequest is a Template type
type ForceTransferRequest struct {
	ForceTransfer ForceTransfer `json:"forceTransfer"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ForceTransferRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ForceTransferRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ForceTransferRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ForceTransferRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ForceTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequest to hex string (Canton MCMS format)
func (t ForceTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequest from hex string (Canton MCMS format)
func (t *ForceTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ForceTransferRequest

// ForceTransferRequestAccept exercises the ForceTransferRequest_Accept choice on this ForceTransferRequest contract
// This method uses the package name in the template ID
func (t ForceTransferRequest) ForceTransferRequestAccept(contractID string, args ForceTransferRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// ForceTransferRequestAcceptWithPackageID exercises the ForceTransferRequest_Accept choice using the provided package ID instead of package name
func (t ForceTransferRequest) ForceTransferRequestAcceptWithPackageID(contractID string, packageID string, args ForceTransferRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// ForceTransferRequestReject exercises the ForceTransferRequest_Reject choice on this ForceTransferRequest contract
// This method uses the package name in the template ID
func (t ForceTransferRequest) ForceTransferRequestReject(contractID string, args ForceTransferRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// ForceTransferRequestRejectWithPackageID exercises the ForceTransferRequest_Reject choice using the provided package ID instead of package name
func (t ForceTransferRequest) ForceTransferRequestRejectWithPackageID(contractID string, packageID string, args ForceTransferRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// ForceTransferRequestCancel exercises the ForceTransferRequest_Cancel choice on this ForceTransferRequest contract
// This method uses the package name in the template ID
func (t ForceTransferRequest) ForceTransferRequestCancel(contractID string, args ForceTransferRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// ForceTransferRequestCancelWithPackageID exercises the ForceTransferRequest_Cancel choice using the provided package ID instead of package name
func (t ForceTransferRequest) ForceTransferRequestCancelWithPackageID(contractID string, packageID string, args ForceTransferRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "ForceTransferRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ForceTransferRequest contract
// This method uses the package name in the template ID
func (t ForceTransferRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ForceTransferRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "ForceTransferRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ForceTransferRequestAccept is a Record type
type ForceTransferRequestAccept struct {
	RegistrarRationale types.TEXT `json:"registrarRationale"`
}

// ToMap converts ForceTransferRequestAccept to a map for DAML arguments
func (t ForceTransferRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarRationale"] = string(t.RegistrarRationale)

	return m
}

func (t ForceTransferRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestAccept to hex string (Canton MCMS format)
func (t ForceTransferRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestAccept from hex string (Canton MCMS format)
func (t *ForceTransferRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequestAcceptResult is a Record type
type ForceTransferRequestAcceptResult struct {
	AcceptedForceTransferCid types.CONTRACT_ID `json:"acceptedForceTransferCid"`
}

// ToMap converts ForceTransferRequestAcceptResult to a map for DAML arguments
func (t ForceTransferRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedForceTransferCid"] = model.NestedToDAMLValue(t.AcceptedForceTransferCid)

	return m
}

func (t ForceTransferRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestAcceptResult to hex string (Canton MCMS format)
func (t ForceTransferRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestAcceptResult from hex string (Canton MCMS format)
func (t *ForceTransferRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequestCancel is a Record type
type ForceTransferRequestCancel struct {
}

// ToMap converts ForceTransferRequestCancel to a map for DAML arguments
func (t ForceTransferRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ForceTransferRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestCancel to hex string (Canton MCMS format)
func (t ForceTransferRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestCancel from hex string (Canton MCMS format)
func (t *ForceTransferRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequestCancelResult is a Record type
type ForceTransferRequestCancelResult struct {
}

// ToMap converts ForceTransferRequestCancelResult to a map for DAML arguments
func (t ForceTransferRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ForceTransferRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestCancelResult to hex string (Canton MCMS format)
func (t ForceTransferRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestCancelResult from hex string (Canton MCMS format)
func (t *ForceTransferRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequestReject is a Record type
type ForceTransferRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts ForceTransferRequestReject to a map for DAML arguments
func (t ForceTransferRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t ForceTransferRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestReject to hex string (Canton MCMS format)
func (t ForceTransferRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestReject from hex string (Canton MCMS format)
func (t *ForceTransferRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ForceTransferRequestRejectResult is a Record type
type ForceTransferRequestRejectResult struct {
	RejectedForceTransferCid types.CONTRACT_ID `json:"rejectedForceTransferCid"`
}

// ToMap converts ForceTransferRequestRejectResult to a map for DAML arguments
func (t ForceTransferRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedForceTransferCid"] = model.NestedToDAMLValue(t.RejectedForceTransferCid)

	return m
}

func (t ForceTransferRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForceTransferRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForceTransferRequestRejectResult to hex string (Canton MCMS format)
func (t ForceTransferRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForceTransferRequestRejectResult from hex string (Canton MCMS format)
func (t *ForceTransferRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentConfiguration is a Template type
type InstrumentConfiguration struct {
	Operator                       types.PARTY                                        `json:"operator"`
	Provider                       types.PARTY                                        `json:"provider"`
	Registrar                      types.PARTY                                        `json:"registrar"`
	DefaultIdentifier              registry_holding_v0.InstrumentIdentifier           `json:"defaultIdentifier"`
	AdditionalIdentifiers          []registry_holding_v0.InstrumentIdentifier         `json:"additionalIdentifiers"`
	IssuerRequirements             []credential_v0.PartyCredentialRequirement         `json:"issuerRequirements"`
	HolderRequirements             []credential_v0.PartyCredentialRequirement         `json:"holderRequirements"`
	ProviderAppRewardBeneficiaries *[]splice_api_featured_app_v1.AppRewardBeneficiary `json:"providerAppRewardBeneficiaries" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t InstrumentConfiguration) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t InstrumentConfiguration) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t InstrumentConfiguration) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["defaultIdentifier"] = model.NestedToDAMLValue(t.DefaultIdentifier)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["additionalIdentifiers"] = func() []any {
		res := make([]any, 0, len(t.AdditionalIdentifiers))
		for _, e := range t.AdditionalIdentifiers {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuerRequirements"] = func() []any {
		res := make([]any, 0, len(t.IssuerRequirements))
		for _, e := range t.IssuerRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holderRequirements"] = func() []any {
		res := make([]any, 0, len(t.HolderRequirements))
		for _, e := range t.HolderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.ProviderAppRewardBeneficiaries != nil {
		args["providerAppRewardBeneficiaries"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ProviderAppRewardBeneficiaries),
		}
	} else {
		args["providerAppRewardBeneficiaries"] = map[string]any{
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
func (t InstrumentConfiguration) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["defaultIdentifier"] = model.NestedToDAMLValue(t.DefaultIdentifier)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["additionalIdentifiers"] = func() []any {
		res := make([]any, 0, len(t.AdditionalIdentifiers))
		for _, e := range t.AdditionalIdentifiers {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuerRequirements"] = func() []any {
		res := make([]any, 0, len(t.IssuerRequirements))
		for _, e := range t.IssuerRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holderRequirements"] = func() []any {
		res := make([]any, 0, len(t.HolderRequirements))
		for _, e := range t.HolderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.ProviderAppRewardBeneficiaries != nil {
		args["providerAppRewardBeneficiaries"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ProviderAppRewardBeneficiaries),
		}
	} else {
		args["providerAppRewardBeneficiaries"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t InstrumentConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentConfiguration to hex string (Canton MCMS format)
func (t InstrumentConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentConfiguration from hex string (Canton MCMS format)
func (t *InstrumentConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for InstrumentConfiguration

// InstrumentConfigurationSetProviderAppRewardBeneficiaries exercises the InstrumentConfiguration_SetProviderAppRewardBeneficiaries choice on this InstrumentConfiguration contract
// This method uses the package name in the template ID
func (t InstrumentConfiguration) InstrumentConfigurationSetProviderAppRewardBeneficiaries(contractID string, args InstrumentConfigurationSetProviderAppRewardBeneficiaries) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "InstrumentConfiguration_SetProviderAppRewardBeneficiaries",
		Arguments:  argsToMap(args),
	}
}

// InstrumentConfigurationSetProviderAppRewardBeneficiariesWithPackageID exercises the InstrumentConfiguration_SetProviderAppRewardBeneficiaries choice using the provided package ID instead of package name
func (t InstrumentConfiguration) InstrumentConfigurationSetProviderAppRewardBeneficiariesWithPackageID(contractID string, packageID string, args InstrumentConfigurationSetProviderAppRewardBeneficiaries) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "InstrumentConfiguration_SetProviderAppRewardBeneficiaries",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this InstrumentConfiguration contract
// This method uses the package name in the template ID
func (t InstrumentConfiguration) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t InstrumentConfiguration) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// InstrumentConfigurationGet exercises the InstrumentConfiguration_Get choice on this InstrumentConfiguration contract
// This method uses the package name in the template ID
func (t InstrumentConfiguration) InstrumentConfigurationGet(contractID string, args InstrumentConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "InstrumentConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// InstrumentConfigurationGetWithPackageID exercises the InstrumentConfiguration_Get choice using the provided package ID instead of package name
func (t InstrumentConfiguration) InstrumentConfigurationGetWithPackageID(contractID string, packageID string, args InstrumentConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Configuration.Instrument", "InstrumentConfiguration"),
		ContractID: contractID,
		Choice:     "InstrumentConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// InstrumentConfigurationGet is a Record type
type InstrumentConfigurationGet struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts InstrumentConfigurationGet to a map for DAML arguments
func (t InstrumentConfigurationGet) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t InstrumentConfigurationGet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentConfigurationGet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentConfigurationGet to hex string (Canton MCMS format)
func (t InstrumentConfigurationGet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentConfigurationGet from hex string (Canton MCMS format)
func (t *InstrumentConfigurationGet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentConfigurationGetResult is a Record type
type InstrumentConfigurationGetResult struct {
	InstrumentConfiguration InstrumentConfiguration `json:"instrumentConfiguration"`
}

// ToMap converts InstrumentConfigurationGetResult to a map for DAML arguments
func (t InstrumentConfigurationGetResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfiguration"] = model.NestedToDAMLValue(t.InstrumentConfiguration)

	return m
}

func (t InstrumentConfigurationGetResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentConfigurationGetResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentConfigurationGetResult to hex string (Canton MCMS format)
func (t InstrumentConfigurationGetResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentConfigurationGetResult from hex string (Canton MCMS format)
func (t *InstrumentConfigurationGetResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentConfigurationSetProviderAppRewardBeneficiaries is a Record type
type InstrumentConfigurationSetProviderAppRewardBeneficiaries struct {
	ProviderAppRewardBeneficiaries *[]splice_api_featured_app_v1.AppRewardBeneficiary `json:"providerAppRewardBeneficiaries" hex:"optional"`
}

// ToMap converts InstrumentConfigurationSetProviderAppRewardBeneficiaries to a map for DAML arguments
func (t InstrumentConfigurationSetProviderAppRewardBeneficiaries) ToMap() map[string]any {
	m := make(map[string]any)

	if t.ProviderAppRewardBeneficiaries != nil {
		m["providerAppRewardBeneficiaries"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ProviderAppRewardBeneficiaries),
		}
	} else {
		m["providerAppRewardBeneficiaries"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t InstrumentConfigurationSetProviderAppRewardBeneficiaries) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentConfigurationSetProviderAppRewardBeneficiaries) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentConfigurationSetProviderAppRewardBeneficiaries to hex string (Canton MCMS format)
func (t InstrumentConfigurationSetProviderAppRewardBeneficiaries) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentConfigurationSetProviderAppRewardBeneficiaries from hex string (Canton MCMS format)
func (t *InstrumentConfigurationSetProviderAppRewardBeneficiaries) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Lock3 is a Record type
type Lock3 struct {
	Operator             types.PARTY                              `json:"operator"`
	Provider             types.PARTY                              `json:"provider"`
	Registrar            types.PARTY                              `json:"registrar"`
	Holder               types.PARTY                              `json:"holder"`
	Locker               types.PARTY                              `json:"locker"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Context              types.TEXT                               `json:"context"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                Batch                                    `json:"batch"`
}

// ToMap converts Lock3 to a map for DAML arguments
func (t Lock3) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["locker"] = t.Locker.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["context"] = string(t.Context)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t Lock3) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Lock3) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Lock3 to hex string (Canton MCMS format)
func (t Lock3) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Lock3 from hex string (Canton MCMS format)
func (t *Lock3) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOffer is a Template type
type LockOffer struct {
	Lock         Lock3      `json:"lock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LockOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LockOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LockOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LockOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOffer to hex string (Canton MCMS format)
func (t LockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOffer from hex string (Canton MCMS format)
func (t *LockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LockOffer

// LockOfferClean exercises the LockOffer_Clean choice on this LockOffer contract
// This method uses the package name in the template ID
func (t LockOffer) LockOfferClean(contractID string, args LockOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// LockOfferCleanWithPackageID exercises the LockOffer_Clean choice using the provided package ID instead of package name
func (t LockOffer) LockOfferCleanWithPackageID(contractID string, packageID string, args LockOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// LockOfferAccept exercises the LockOffer_Accept choice on this LockOffer contract
// This method uses the package name in the template ID
func (t LockOffer) LockOfferAccept(contractID string, args LockOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// LockOfferAcceptWithPackageID exercises the LockOffer_Accept choice using the provided package ID instead of package name
func (t LockOffer) LockOfferAcceptWithPackageID(contractID string, packageID string, args LockOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// LockOfferReject exercises the LockOffer_Reject choice on this LockOffer contract
// This method uses the package name in the template ID
func (t LockOffer) LockOfferReject(contractID string, args LockOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// LockOfferRejectWithPackageID exercises the LockOffer_Reject choice using the provided package ID instead of package name
func (t LockOffer) LockOfferRejectWithPackageID(contractID string, packageID string, args LockOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// LockOfferCancel exercises the LockOffer_Cancel choice on this LockOffer contract
// This method uses the package name in the template ID
func (t LockOffer) LockOfferCancel(contractID string, args LockOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// LockOfferCancelWithPackageID exercises the LockOffer_Cancel choice using the provided package ID instead of package name
func (t LockOffer) LockOfferCancelWithPackageID(contractID string, packageID string, args LockOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "LockOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LockOffer contract
// This method uses the package name in the template ID
func (t LockOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// LockOfferAccept is a Record type
type LockOfferAccept struct {
}

// ToMap converts LockOfferAccept to a map for DAML arguments
func (t LockOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t LockOfferAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferAccept to hex string (Canton MCMS format)
func (t LockOfferAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferAccept from hex string (Canton MCMS format)
func (t *LockOfferAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferAcceptResult is a Record type
type LockOfferAcceptResult struct {
	AcceptedLockCid types.CONTRACT_ID `json:"acceptedLockCid"`
}

// ToMap converts LockOfferAcceptResult to a map for DAML arguments
func (t LockOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedLockCid"] = model.NestedToDAMLValue(t.AcceptedLockCid)

	return m
}

func (t LockOfferAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferAcceptResult to hex string (Canton MCMS format)
func (t LockOfferAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferAcceptResult from hex string (Canton MCMS format)
func (t *LockOfferAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferCancel is a Record type
type LockOfferCancel struct {
}

// ToMap converts LockOfferCancel to a map for DAML arguments
func (t LockOfferCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t LockOfferCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferCancel to hex string (Canton MCMS format)
func (t LockOfferCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferCancel from hex string (Canton MCMS format)
func (t *LockOfferCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferCancelResult is a Record type
type LockOfferCancelResult struct {
}

// ToMap converts LockOfferCancelResult to a map for DAML arguments
func (t LockOfferCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t LockOfferCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferCancelResult to hex string (Canton MCMS format)
func (t LockOfferCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferCancelResult from hex string (Canton MCMS format)
func (t *LockOfferCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferClean is a Record type
type LockOfferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts LockOfferClean to a map for DAML arguments
func (t LockOfferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t LockOfferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferClean to hex string (Canton MCMS format)
func (t LockOfferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferClean from hex string (Canton MCMS format)
func (t *LockOfferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferReject is a Record type
type LockOfferReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts LockOfferReject to a map for DAML arguments
func (t LockOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t LockOfferReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferReject to hex string (Canton MCMS format)
func (t LockOfferReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferReject from hex string (Canton MCMS format)
func (t *LockOfferReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockOfferRejectResult is a Record type
type LockOfferRejectResult struct {
	RejectedLockCid types.CONTRACT_ID `json:"rejectedLockCid"`
}

// ToMap converts LockOfferRejectResult to a map for DAML arguments
func (t LockOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedLockCid"] = model.NestedToDAMLValue(t.RejectedLockCid)

	return m
}

func (t LockOfferRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockOfferRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockOfferRejectResult to hex string (Canton MCMS format)
func (t LockOfferRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockOfferRejectResult from hex string (Canton MCMS format)
func (t *LockOfferRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequest is a Template type
type LockRequest struct {
	Lock Lock3 `json:"lock"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LockRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LockRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LockRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LockRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequest to hex string (Canton MCMS format)
func (t LockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequest from hex string (Canton MCMS format)
func (t *LockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LockRequest

// LockRequestClean exercises the LockRequest_Clean choice on this LockRequest contract
// This method uses the package name in the template ID
func (t LockRequest) LockRequestClean(contractID string, args LockRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// LockRequestCleanWithPackageID exercises the LockRequest_Clean choice using the provided package ID instead of package name
func (t LockRequest) LockRequestCleanWithPackageID(contractID string, packageID string, args LockRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// LockRequestAccept exercises the LockRequest_Accept choice on this LockRequest contract
// This method uses the package name in the template ID
func (t LockRequest) LockRequestAccept(contractID string, args LockRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// LockRequestAcceptWithPackageID exercises the LockRequest_Accept choice using the provided package ID instead of package name
func (t LockRequest) LockRequestAcceptWithPackageID(contractID string, packageID string, args LockRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// LockRequestReject exercises the LockRequest_Reject choice on this LockRequest contract
// This method uses the package name in the template ID
func (t LockRequest) LockRequestReject(contractID string, args LockRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// LockRequestRejectWithPackageID exercises the LockRequest_Reject choice using the provided package ID instead of package name
func (t LockRequest) LockRequestRejectWithPackageID(contractID string, packageID string, args LockRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// LockRequestCancel exercises the LockRequest_Cancel choice on this LockRequest contract
// This method uses the package name in the template ID
func (t LockRequest) LockRequestCancel(contractID string, args LockRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// LockRequestCancelWithPackageID exercises the LockRequest_Cancel choice using the provided package ID instead of package name
func (t LockRequest) LockRequestCancelWithPackageID(contractID string, packageID string, args LockRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "LockRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LockRequest contract
// This method uses the package name in the template ID
func (t LockRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "LockRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// LockRequestAccept is a Record type
type LockRequestAccept struct {
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// ToMap converts LockRequestAccept to a map for DAML arguments
func (t LockRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingLabel"] = string(t.HoldingLabel)

	return m
}

func (t LockRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestAccept to hex string (Canton MCMS format)
func (t LockRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestAccept from hex string (Canton MCMS format)
func (t *LockRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestAcceptResult is a Record type
type LockRequestAcceptResult struct {
	AcceptedLockCid types.CONTRACT_ID `json:"acceptedLockCid"`
}

// ToMap converts LockRequestAcceptResult to a map for DAML arguments
func (t LockRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedLockCid"] = model.NestedToDAMLValue(t.AcceptedLockCid)

	return m
}

func (t LockRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestAcceptResult to hex string (Canton MCMS format)
func (t LockRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestAcceptResult from hex string (Canton MCMS format)
func (t *LockRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestCancel is a Record type
type LockRequestCancel struct {
}

// ToMap converts LockRequestCancel to a map for DAML arguments
func (t LockRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t LockRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestCancel to hex string (Canton MCMS format)
func (t LockRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestCancel from hex string (Canton MCMS format)
func (t *LockRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestCancelResult is a Record type
type LockRequestCancelResult struct {
}

// ToMap converts LockRequestCancelResult to a map for DAML arguments
func (t LockRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t LockRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestCancelResult to hex string (Canton MCMS format)
func (t LockRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestCancelResult from hex string (Canton MCMS format)
func (t *LockRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestClean is a Record type
type LockRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts LockRequestClean to a map for DAML arguments
func (t LockRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t LockRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestClean to hex string (Canton MCMS format)
func (t LockRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestClean from hex string (Canton MCMS format)
func (t *LockRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestReject is a Record type
type LockRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts LockRequestReject to a map for DAML arguments
func (t LockRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t LockRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestReject to hex string (Canton MCMS format)
func (t LockRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestReject from hex string (Canton MCMS format)
func (t *LockRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LockRequestRejectResult is a Record type
type LockRequestRejectResult struct {
	RejectedLockCid types.CONTRACT_ID `json:"rejectedLockCid"`
}

// ToMap converts LockRequestRejectResult to a map for DAML arguments
func (t LockRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedLockCid"] = model.NestedToDAMLValue(t.RejectedLockCid)

	return m
}

func (t LockRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockRequestRejectResult to hex string (Canton MCMS format)
func (t LockRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockRequestRejectResult from hex string (Canton MCMS format)
func (t *LockRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Mint is a Record type
type Mint struct {
	Operator             types.PARTY                              `json:"operator"`
	Provider             types.PARTY                              `json:"provider"`
	Registrar            types.PARTY                              `json:"registrar"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Holder               types.PARTY                              `json:"holder"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                Batch                                    `json:"batch"`
}

// ToMap converts Mint to a map for DAML arguments
func (t Mint) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t Mint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Mint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Mint to hex string (Canton MCMS format)
func (t Mint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Mint from hex string (Canton MCMS format)
func (t *Mint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOffer is a Template type
type MintOffer struct {
	Mint Mint `json:"mint"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MintOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOffer to hex string (Canton MCMS format)
func (t MintOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOffer from hex string (Canton MCMS format)
func (t *MintOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MintOffer

// MintOfferClean exercises the MintOffer_Clean choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferClean(contractID string, args MintOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// MintOfferCleanWithPackageID exercises the MintOffer_Clean choice using the provided package ID instead of package name
func (t MintOffer) MintOfferCleanWithPackageID(contractID string, packageID string, args MintOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// MintOfferAccept exercises the MintOffer_Accept choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferAccept(contractID string, args MintOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintOfferAcceptWithPackageID exercises the MintOffer_Accept choice using the provided package ID instead of package name
func (t MintOffer) MintOfferAcceptWithPackageID(contractID string, packageID string, args MintOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintOfferReject exercises the MintOffer_Reject choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferReject(contractID string, args MintOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintOfferRejectWithPackageID exercises the MintOffer_Reject choice using the provided package ID instead of package name
func (t MintOffer) MintOfferRejectWithPackageID(contractID string, packageID string, args MintOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintOfferCancel exercises the MintOffer_Cancel choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferCancel(contractID string, args MintOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// MintOfferCancelWithPackageID exercises the MintOffer_Cancel choice using the provided package ID instead of package name
func (t MintOffer) MintOfferCancelWithPackageID(contractID string, packageID string, args MintOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintOfferAccept is a Record type
type MintOfferAccept struct {
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// ToMap converts MintOfferAccept to a map for DAML arguments
func (t MintOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingLabel"] = string(t.HoldingLabel)

	return m
}

func (t MintOfferAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferAccept to hex string (Canton MCMS format)
func (t MintOfferAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferAccept from hex string (Canton MCMS format)
func (t *MintOfferAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferAcceptResult is a Record type
type MintOfferAcceptResult struct {
	AcceptedMintCid types.CONTRACT_ID `json:"acceptedMintCid"`
}

// ToMap converts MintOfferAcceptResult to a map for DAML arguments
func (t MintOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedMintCid"] = model.NestedToDAMLValue(t.AcceptedMintCid)

	return m
}

func (t MintOfferAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferAcceptResult to hex string (Canton MCMS format)
func (t MintOfferAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferAcceptResult from hex string (Canton MCMS format)
func (t *MintOfferAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferCancel is a Record type
type MintOfferCancel struct {
}

// ToMap converts MintOfferCancel to a map for DAML arguments
func (t MintOfferCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t MintOfferCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferCancel to hex string (Canton MCMS format)
func (t MintOfferCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferCancel from hex string (Canton MCMS format)
func (t *MintOfferCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferCancelResult is a Record type
type MintOfferCancelResult struct {
}

// ToMap converts MintOfferCancelResult to a map for DAML arguments
func (t MintOfferCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t MintOfferCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferCancelResult to hex string (Canton MCMS format)
func (t MintOfferCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferCancelResult from hex string (Canton MCMS format)
func (t *MintOfferCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferClean is a Record type
type MintOfferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts MintOfferClean to a map for DAML arguments
func (t MintOfferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t MintOfferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferClean to hex string (Canton MCMS format)
func (t MintOfferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferClean from hex string (Canton MCMS format)
func (t *MintOfferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferReject is a Record type
type MintOfferReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts MintOfferReject to a map for DAML arguments
func (t MintOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t MintOfferReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferReject to hex string (Canton MCMS format)
func (t MintOfferReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferReject from hex string (Canton MCMS format)
func (t *MintOfferReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintOfferRejectResult is a Record type
type MintOfferRejectResult struct {
	RejectedMintCid types.CONTRACT_ID `json:"rejectedMintCid"`
}

// ToMap converts MintOfferRejectResult to a map for DAML arguments
func (t MintOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedMintCid"] = model.NestedToDAMLValue(t.RejectedMintCid)

	return m
}

func (t MintOfferRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintOfferRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintOfferRejectResult to hex string (Canton MCMS format)
func (t MintOfferRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintOfferRejectResult from hex string (Canton MCMS format)
func (t *MintOfferRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequest is a Template type
type MintRequest struct {
	Mint         Mint       `json:"mint"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MintRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequest to hex string (Canton MCMS format)
func (t MintRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequest from hex string (Canton MCMS format)
func (t *MintRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MintRequest

// MintRequestClean exercises the MintRequest_Clean choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestClean(contractID string, args MintRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// MintRequestCleanWithPackageID exercises the MintRequest_Clean choice using the provided package ID instead of package name
func (t MintRequest) MintRequestCleanWithPackageID(contractID string, packageID string, args MintRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// MintRequestAccept exercises the MintRequest_Accept choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestAccept(contractID string, args MintRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintRequestAcceptWithPackageID exercises the MintRequest_Accept choice using the provided package ID instead of package name
func (t MintRequest) MintRequestAcceptWithPackageID(contractID string, packageID string, args MintRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintRequestReject exercises the MintRequest_Reject choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestReject(contractID string, args MintRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintRequestRejectWithPackageID exercises the MintRequest_Reject choice using the provided package ID instead of package name
func (t MintRequest) MintRequestRejectWithPackageID(contractID string, packageID string, args MintRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintRequestCancel exercises the MintRequest_Cancel choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestCancel(contractID string, args MintRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// MintRequestCancelWithPackageID exercises the MintRequest_Cancel choice using the provided package ID instead of package name
func (t MintRequest) MintRequestCancelWithPackageID(contractID string, packageID string, args MintRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintRequestAccept is a Record type
type MintRequestAccept struct {
}

// ToMap converts MintRequestAccept to a map for DAML arguments
func (t MintRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t MintRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestAccept to hex string (Canton MCMS format)
func (t MintRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestAccept from hex string (Canton MCMS format)
func (t *MintRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestAcceptResult is a Record type
type MintRequestAcceptResult struct {
	AcceptedMintCid types.CONTRACT_ID `json:"acceptedMintCid"`
}

// ToMap converts MintRequestAcceptResult to a map for DAML arguments
func (t MintRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedMintCid"] = model.NestedToDAMLValue(t.AcceptedMintCid)

	return m
}

func (t MintRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestAcceptResult to hex string (Canton MCMS format)
func (t MintRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestAcceptResult from hex string (Canton MCMS format)
func (t *MintRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestCancel is a Record type
type MintRequestCancel struct {
}

// ToMap converts MintRequestCancel to a map for DAML arguments
func (t MintRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t MintRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestCancel to hex string (Canton MCMS format)
func (t MintRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestCancel from hex string (Canton MCMS format)
func (t *MintRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestCancelResult is a Record type
type MintRequestCancelResult struct {
}

// ToMap converts MintRequestCancelResult to a map for DAML arguments
func (t MintRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t MintRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestCancelResult to hex string (Canton MCMS format)
func (t MintRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestCancelResult from hex string (Canton MCMS format)
func (t *MintRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestClean is a Record type
type MintRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts MintRequestClean to a map for DAML arguments
func (t MintRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t MintRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestClean to hex string (Canton MCMS format)
func (t MintRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestClean from hex string (Canton MCMS format)
func (t *MintRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestReject is a Record type
type MintRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts MintRequestReject to a map for DAML arguments
func (t MintRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t MintRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestReject to hex string (Canton MCMS format)
func (t MintRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestReject from hex string (Canton MCMS format)
func (t *MintRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRequestRejectResult is a Record type
type MintRequestRejectResult struct {
	RejectedMintCid types.CONTRACT_ID `json:"rejectedMintCid"`
}

// ToMap converts MintRequestRejectResult to a map for DAML arguments
func (t MintRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedMintCid"] = model.NestedToDAMLValue(t.RejectedMintCid)

	return m
}

func (t MintRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRequestRejectResult to hex string (Canton MCMS format)
func (t MintRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRequestRejectResult from hex string (Canton MCMS format)
func (t *MintRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedBurn is a Template type
type RejectedBurn struct {
	Burn   Burn       `json:"burn"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "RejectedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "RejectedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedBurn to hex string (Canton MCMS format)
func (t RejectedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedBurn from hex string (Canton MCMS format)
func (t *RejectedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedBurn

// RejectedBurnClean exercises the RejectedBurn_Clean choice on this RejectedBurn contract
// This method uses the package name in the template ID
func (t RejectedBurn) RejectedBurnClean(contractID string, args RejectedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedBurnCleanWithPackageID exercises the RejectedBurn_Clean choice using the provided package ID instead of package name
func (t RejectedBurn) RejectedBurnCleanWithPackageID(contractID string, packageID string, args RejectedBurnClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedBurnDelete exercises the RejectedBurn_Delete choice on this RejectedBurn contract
// This method uses the package name in the template ID
func (t RejectedBurn) RejectedBurnDelete(contractID string, args RejectedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedBurnDeleteWithPackageID exercises the RejectedBurn_Delete choice using the provided package ID instead of package name
func (t RejectedBurn) RejectedBurnDeleteWithPackageID(contractID string, packageID string, args RejectedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedBurn contract
// This method uses the package name in the template ID
func (t RejectedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedBurnClean is a Record type
type RejectedBurnClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedBurnClean to a map for DAML arguments
func (t RejectedBurnClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RejectedBurnClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedBurnClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedBurnClean to hex string (Canton MCMS format)
func (t RejectedBurnClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedBurnClean from hex string (Canton MCMS format)
func (t *RejectedBurnClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedBurnDelete is a Record type
type RejectedBurnDelete struct {
}

// ToMap converts RejectedBurnDelete to a map for DAML arguments
func (t RejectedBurnDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedBurnDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedBurnDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedBurnDelete to hex string (Canton MCMS format)
func (t RejectedBurnDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedBurnDelete from hex string (Canton MCMS format)
func (t *RejectedBurnDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedBurnDeleteResult is a Record type
type RejectedBurnDeleteResult struct {
}

// ToMap converts RejectedBurnDeleteResult to a map for DAML arguments
func (t RejectedBurnDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedBurnDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedBurnDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedBurnDeleteResult to hex string (Canton MCMS format)
func (t RejectedBurnDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedBurnDeleteResult from hex string (Canton MCMS format)
func (t *RejectedBurnDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedForceTransfer is a Template type
type RejectedForceTransfer struct {
	ForceTransfer ForceTransfer `json:"forceTransfer"`
	Reason        types.TEXT    `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedForceTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedForceTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedForceTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedForceTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forceTransfer"] = model.NestedToDAMLValue(t.ForceTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedForceTransfer to hex string (Canton MCMS format)
func (t RejectedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedForceTransfer from hex string (Canton MCMS format)
func (t *RejectedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedForceTransfer

// RejectedForceTransferDelete exercises the RejectedForceTransfer_Delete choice on this RejectedForceTransfer contract
// This method uses the package name in the template ID
func (t RejectedForceTransfer) RejectedForceTransferDelete(contractID string, args RejectedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer"),
		ContractID: contractID,
		Choice:     "RejectedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedForceTransferDeleteWithPackageID exercises the RejectedForceTransfer_Delete choice using the provided package ID instead of package name
func (t RejectedForceTransfer) RejectedForceTransferDeleteWithPackageID(contractID string, packageID string, args RejectedForceTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer"),
		ContractID: contractID,
		Choice:     "RejectedForceTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedForceTransfer contract
// This method uses the package name in the template ID
func (t RejectedForceTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedForceTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.ForceTransfer", "RejectedForceTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedForceTransferDelete is a Record type
type RejectedForceTransferDelete struct {
}

// ToMap converts RejectedForceTransferDelete to a map for DAML arguments
func (t RejectedForceTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedForceTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedForceTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedForceTransferDelete to hex string (Canton MCMS format)
func (t RejectedForceTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedForceTransferDelete from hex string (Canton MCMS format)
func (t *RejectedForceTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedForceTransferDeleteResult is a Record type
type RejectedForceTransferDeleteResult struct {
}

// ToMap converts RejectedForceTransferDeleteResult to a map for DAML arguments
func (t RejectedForceTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedForceTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedForceTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedForceTransferDeleteResult to hex string (Canton MCMS format)
func (t RejectedForceTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedForceTransferDeleteResult from hex string (Canton MCMS format)
func (t *RejectedForceTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedLock is a Template type
type RejectedLock struct {
	Lock   Lock3      `json:"lock"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedLock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "RejectedLock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedLock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "RejectedLock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedLock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedLock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lock"] = model.NestedToDAMLValue(t.Lock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedLock to hex string (Canton MCMS format)
func (t RejectedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedLock from hex string (Canton MCMS format)
func (t *RejectedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedLock

// RejectedLockClean exercises the RejectedLock_Clean choice on this RejectedLock contract
// This method uses the package name in the template ID
func (t RejectedLock) RejectedLockClean(contractID string, args RejectedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "RejectedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedLockCleanWithPackageID exercises the RejectedLock_Clean choice using the provided package ID instead of package name
func (t RejectedLock) RejectedLockCleanWithPackageID(contractID string, packageID string, args RejectedLockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "RejectedLock_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedLockDelete exercises the RejectedLock_Delete choice on this RejectedLock contract
// This method uses the package name in the template ID
func (t RejectedLock) RejectedLockDelete(contractID string, args RejectedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "RejectedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedLockDeleteWithPackageID exercises the RejectedLock_Delete choice using the provided package ID instead of package name
func (t RejectedLock) RejectedLockDeleteWithPackageID(contractID string, packageID string, args RejectedLockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "RejectedLock_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedLock contract
// This method uses the package name in the template ID
func (t RejectedLock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedLock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Lock", "RejectedLock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedLockClean is a Record type
type RejectedLockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedLockClean to a map for DAML arguments
func (t RejectedLockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RejectedLockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedLockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedLockClean to hex string (Canton MCMS format)
func (t RejectedLockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedLockClean from hex string (Canton MCMS format)
func (t *RejectedLockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedLockDelete is a Record type
type RejectedLockDelete struct {
}

// ToMap converts RejectedLockDelete to a map for DAML arguments
func (t RejectedLockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedLockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedLockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedLockDelete to hex string (Canton MCMS format)
func (t RejectedLockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedLockDelete from hex string (Canton MCMS format)
func (t *RejectedLockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedLockDeleteResult is a Record type
type RejectedLockDeleteResult struct {
}

// ToMap converts RejectedLockDeleteResult to a map for DAML arguments
func (t RejectedLockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedLockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedLockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedLockDeleteResult to hex string (Canton MCMS format)
func (t RejectedLockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedLockDeleteResult from hex string (Canton MCMS format)
func (t *RejectedLockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedMint is a Template type
type RejectedMint struct {
	Mint   Mint       `json:"mint"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "RejectedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "RejectedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedMint to hex string (Canton MCMS format)
func (t RejectedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedMint from hex string (Canton MCMS format)
func (t *RejectedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedMint

// RejectedMintClean exercises the RejectedMint_Clean choice on this RejectedMint contract
// This method uses the package name in the template ID
func (t RejectedMint) RejectedMintClean(contractID string, args RejectedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedMintCleanWithPackageID exercises the RejectedMint_Clean choice using the provided package ID instead of package name
func (t RejectedMint) RejectedMintCleanWithPackageID(contractID string, packageID string, args RejectedMintClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedMintDelete exercises the RejectedMint_Delete choice on this RejectedMint contract
// This method uses the package name in the template ID
func (t RejectedMint) RejectedMintDelete(contractID string, args RejectedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedMintDeleteWithPackageID exercises the RejectedMint_Delete choice using the provided package ID instead of package name
func (t RejectedMint) RejectedMintDeleteWithPackageID(contractID string, packageID string, args RejectedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedMint contract
// This method uses the package name in the template ID
func (t RejectedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedMintClean is a Record type
type RejectedMintClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedMintClean to a map for DAML arguments
func (t RejectedMintClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RejectedMintClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedMintClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedMintClean to hex string (Canton MCMS format)
func (t RejectedMintClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedMintClean from hex string (Canton MCMS format)
func (t *RejectedMintClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedMintDelete is a Record type
type RejectedMintDelete struct {
}

// ToMap converts RejectedMintDelete to a map for DAML arguments
func (t RejectedMintDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedMintDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedMintDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedMintDelete to hex string (Canton MCMS format)
func (t RejectedMintDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedMintDelete from hex string (Canton MCMS format)
func (t *RejectedMintDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedMintDeleteResult is a Record type
type RejectedMintDeleteResult struct {
}

// ToMap converts RejectedMintDeleteResult to a map for DAML arguments
func (t RejectedMintDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedMintDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedMintDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedMintDeleteResult to hex string (Canton MCMS format)
func (t RejectedMintDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedMintDeleteResult from hex string (Canton MCMS format)
func (t *RejectedMintDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedTransfer is a Template type
type RejectedTransfer struct {
	Transfer           Transfer2   `json:"transfer"`
	Reason             types.TEXT  `json:"reason"`
	Observers          *types.SET  `json:"observers" hex:"optional"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedTransfer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedTransfer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedTransfer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OperatorIsObserver != nil {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.OperatorIsObserver),
		}
	} else {
		args["operatorIsObserver"] = map[string]any{
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
func (t RejectedTransfer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	if t.Observers != nil {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Observers),
		}
	} else {
		args["observers"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OperatorIsObserver != nil {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.OperatorIsObserver),
		}
	} else {
		args["operatorIsObserver"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedTransfer to hex string (Canton MCMS format)
func (t RejectedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedTransfer from hex string (Canton MCMS format)
func (t *RejectedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedTransfer

// RejectedTransferDelete exercises the RejectedTransfer_Delete choice on this RejectedTransfer contract
// This method uses the package name in the template ID
func (t RejectedTransfer) RejectedTransferDelete(contractID string, args RejectedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer"),
		ContractID: contractID,
		Choice:     "RejectedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedTransferDeleteWithPackageID exercises the RejectedTransfer_Delete choice using the provided package ID instead of package name
func (t RejectedTransfer) RejectedTransferDeleteWithPackageID(contractID string, packageID string, args RejectedTransferDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer"),
		ContractID: contractID,
		Choice:     "RejectedTransfer_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedTransfer contract
// This method uses the package name in the template ID
func (t RejectedTransfer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedTransfer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "RejectedTransfer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedTransferDelete is a Record type
type RejectedTransferDelete struct {
}

// ToMap converts RejectedTransferDelete to a map for DAML arguments
func (t RejectedTransferDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedTransferDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedTransferDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedTransferDelete to hex string (Canton MCMS format)
func (t RejectedTransferDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedTransferDelete from hex string (Canton MCMS format)
func (t *RejectedTransferDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedTransferDeleteResult is a Record type
type RejectedTransferDeleteResult struct {
}

// ToMap converts RejectedTransferDeleteResult to a map for DAML arguments
func (t RejectedTransferDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedTransferDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedTransferDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedTransferDeleteResult to hex string (Canton MCMS format)
func (t RejectedTransferDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedTransferDeleteResult from hex string (Canton MCMS format)
func (t *RejectedTransferDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedUnlock is a Template type
type RejectedUnlock struct {
	Unlock Unlock     `json:"unlock"`
	Reason types.TEXT `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedUnlock) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedUnlock) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedUnlock) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedUnlock) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedUnlock to hex string (Canton MCMS format)
func (t RejectedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedUnlock from hex string (Canton MCMS format)
func (t *RejectedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedUnlock

// RejectedUnlockClean exercises the RejectedUnlock_Clean choice on this RejectedUnlock contract
// This method uses the package name in the template ID
func (t RejectedUnlock) RejectedUnlockClean(contractID string, args RejectedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "RejectedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedUnlockCleanWithPackageID exercises the RejectedUnlock_Clean choice using the provided package ID instead of package name
func (t RejectedUnlock) RejectedUnlockCleanWithPackageID(contractID string, packageID string, args RejectedUnlockClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "RejectedUnlock_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedUnlockDelete exercises the RejectedUnlock_Delete choice on this RejectedUnlock contract
// This method uses the package name in the template ID
func (t RejectedUnlock) RejectedUnlockDelete(contractID string, args RejectedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "RejectedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedUnlockDeleteWithPackageID exercises the RejectedUnlock_Delete choice using the provided package ID instead of package name
func (t RejectedUnlock) RejectedUnlockDeleteWithPackageID(contractID string, packageID string, args RejectedUnlockDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "RejectedUnlock_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedUnlock contract
// This method uses the package name in the template ID
func (t RejectedUnlock) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedUnlock) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "RejectedUnlock"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedUnlockClean is a Record type
type RejectedUnlockClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedUnlockClean to a map for DAML arguments
func (t RejectedUnlockClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RejectedUnlockClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedUnlockClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedUnlockClean to hex string (Canton MCMS format)
func (t RejectedUnlockClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedUnlockClean from hex string (Canton MCMS format)
func (t *RejectedUnlockClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedUnlockDelete is a Record type
type RejectedUnlockDelete struct {
}

// ToMap converts RejectedUnlockDelete to a map for DAML arguments
func (t RejectedUnlockDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedUnlockDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedUnlockDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedUnlockDelete to hex string (Canton MCMS format)
func (t RejectedUnlockDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedUnlockDelete from hex string (Canton MCMS format)
func (t *RejectedUnlockDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedUnlockDeleteResult is a Record type
type RejectedUnlockDeleteResult struct {
}

// ToMap converts RejectedUnlockDeleteResult to a map for DAML arguments
func (t RejectedUnlockDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedUnlockDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedUnlockDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedUnlockDeleteResult to hex string (Canton MCMS format)
func (t RejectedUnlockDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedUnlockDeleteResult from hex string (Canton MCMS format)
func (t *RejectedUnlockDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Transfer2 is a Record type
type Transfer2 struct {
	Operator             types.PARTY                              `json:"operator"`
	Provider             types.PARTY                              `json:"provider"`
	Registrar            types.PARTY                              `json:"registrar"`
	Sender               types.PARTY                              `json:"sender"`
	Receiver             types.PARTY                              `json:"receiver"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                Batch                                    `json:"batch"`
}

// ToMap converts Transfer2 to a map for DAML arguments
func (t Transfer2) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t Transfer2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Transfer2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Transfer2 to hex string (Canton MCMS format)
func (t Transfer2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Transfer2 from hex string (Canton MCMS format)
func (t *Transfer2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOffer is a Template type
type TransferOffer struct {
	Transfer    Transfer2  `json:"transfer"`
	SenderLabel types.TEXT `json:"senderLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TransferOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TransferOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TransferOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TransferOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["senderLabel"] = string(t.SenderLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TransferOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOffer to hex string (Canton MCMS format)
func (t TransferOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOffer from hex string (Canton MCMS format)
func (t *TransferOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TransferOffer

// TransferOfferClean exercises the TransferOffer_Clean choice on this TransferOffer contract
// This method uses the package name in the template ID
func (t TransferOffer) TransferOfferClean(contractID string, args TransferOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferCleanWithPackageID exercises the TransferOffer_Clean choice using the provided package ID instead of package name
func (t TransferOffer) TransferOfferCleanWithPackageID(contractID string, packageID string, args TransferOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferAccept exercises the TransferOffer_Accept choice on this TransferOffer contract
// This method uses the package name in the template ID
func (t TransferOffer) TransferOfferAccept(contractID string, args TransferOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferAcceptWithPackageID exercises the TransferOffer_Accept choice using the provided package ID instead of package name
func (t TransferOffer) TransferOfferAcceptWithPackageID(contractID string, packageID string, args TransferOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferReject exercises the TransferOffer_Reject choice on this TransferOffer contract
// This method uses the package name in the template ID
func (t TransferOffer) TransferOfferReject(contractID string, args TransferOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferRejectWithPackageID exercises the TransferOffer_Reject choice using the provided package ID instead of package name
func (t TransferOffer) TransferOfferRejectWithPackageID(contractID string, packageID string, args TransferOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferCancel exercises the TransferOffer_Cancel choice on this TransferOffer contract
// This method uses the package name in the template ID
func (t TransferOffer) TransferOfferCancel(contractID string, args TransferOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// TransferOfferCancelWithPackageID exercises the TransferOffer_Cancel choice using the provided package ID instead of package name
func (t TransferOffer) TransferOfferCancelWithPackageID(contractID string, packageID string, args TransferOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "TransferOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferOffer contract
// This method uses the package name in the template ID
func (t TransferOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TransferOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferOfferAccept is a Record type
type TransferOfferAccept struct {
	ReceiverLabel types.TEXT `json:"receiverLabel"`
}

// ToMap converts TransferOfferAccept to a map for DAML arguments
func (t TransferOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverLabel"] = string(t.ReceiverLabel)

	return m
}

func (t TransferOfferAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferAccept to hex string (Canton MCMS format)
func (t TransferOfferAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferAccept from hex string (Canton MCMS format)
func (t *TransferOfferAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferAcceptResult is a Record type
type TransferOfferAcceptResult struct {
	AcceptedTransferCid types.CONTRACT_ID `json:"acceptedTransferCid"`
}

// ToMap converts TransferOfferAcceptResult to a map for DAML arguments
func (t TransferOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedTransferCid"] = model.NestedToDAMLValue(t.AcceptedTransferCid)

	return m
}

func (t TransferOfferAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferAcceptResult to hex string (Canton MCMS format)
func (t TransferOfferAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferAcceptResult from hex string (Canton MCMS format)
func (t *TransferOfferAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferCancel is a Record type
type TransferOfferCancel struct {
}

// ToMap converts TransferOfferCancel to a map for DAML arguments
func (t TransferOfferCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t TransferOfferCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferCancel to hex string (Canton MCMS format)
func (t TransferOfferCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferCancel from hex string (Canton MCMS format)
func (t *TransferOfferCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferCancelResult is a Record type
type TransferOfferCancelResult struct {
}

// ToMap converts TransferOfferCancelResult to a map for DAML arguments
func (t TransferOfferCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t TransferOfferCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferCancelResult to hex string (Canton MCMS format)
func (t TransferOfferCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferCancelResult from hex string (Canton MCMS format)
func (t *TransferOfferCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferClean is a Record type
type TransferOfferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts TransferOfferClean to a map for DAML arguments
func (t TransferOfferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t TransferOfferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferClean to hex string (Canton MCMS format)
func (t TransferOfferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferClean from hex string (Canton MCMS format)
func (t *TransferOfferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferReject is a Record type
type TransferOfferReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts TransferOfferReject to a map for DAML arguments
func (t TransferOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t TransferOfferReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferReject to hex string (Canton MCMS format)
func (t TransferOfferReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferReject from hex string (Canton MCMS format)
func (t *TransferOfferReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOfferRejectResult is a Record type
type TransferOfferRejectResult struct {
	RejectedTransferCid types.CONTRACT_ID `json:"rejectedTransferCid"`
}

// ToMap converts TransferOfferRejectResult to a map for DAML arguments
func (t TransferOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedTransferCid"] = model.NestedToDAMLValue(t.RejectedTransferCid)

	return m
}

func (t TransferOfferRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferOfferRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferOfferRejectResult to hex string (Canton MCMS format)
func (t TransferOfferRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferOfferRejectResult from hex string (Canton MCMS format)
func (t *TransferOfferRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequest is a Template type
type TransferRequest struct {
	Transfer      Transfer2  `json:"transfer"`
	ReceiverLabel types.TEXT `json:"receiverLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TransferRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TransferRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TransferRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TransferRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverLabel"] = string(t.ReceiverLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequest to hex string (Canton MCMS format)
func (t TransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequest from hex string (Canton MCMS format)
func (t *TransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TransferRequest

// TransferRequestClean exercises the TransferRequest_Clean choice on this TransferRequest contract
// This method uses the package name in the template ID
func (t TransferRequest) TransferRequestClean(contractID string, args TransferRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestCleanWithPackageID exercises the TransferRequest_Clean choice using the provided package ID instead of package name
func (t TransferRequest) TransferRequestCleanWithPackageID(contractID string, packageID string, args TransferRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestAccept exercises the TransferRequest_Accept choice on this TransferRequest contract
// This method uses the package name in the template ID
func (t TransferRequest) TransferRequestAccept(contractID string, args TransferRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestAcceptWithPackageID exercises the TransferRequest_Accept choice using the provided package ID instead of package name
func (t TransferRequest) TransferRequestAcceptWithPackageID(contractID string, packageID string, args TransferRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestReject exercises the TransferRequest_Reject choice on this TransferRequest contract
// This method uses the package name in the template ID
func (t TransferRequest) TransferRequestReject(contractID string, args TransferRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestRejectWithPackageID exercises the TransferRequest_Reject choice using the provided package ID instead of package name
func (t TransferRequest) TransferRequestRejectWithPackageID(contractID string, packageID string, args TransferRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestCancel exercises the TransferRequest_Cancel choice on this TransferRequest contract
// This method uses the package name in the template ID
func (t TransferRequest) TransferRequestCancel(contractID string, args TransferRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// TransferRequestCancelWithPackageID exercises the TransferRequest_Cancel choice using the provided package ID instead of package name
func (t TransferRequest) TransferRequestCancelWithPackageID(contractID string, packageID string, args TransferRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "TransferRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferRequest contract
// This method uses the package name in the template ID
func (t TransferRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TransferRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Transfer", "TransferRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferRequestAccept is a Record type
type TransferRequestAccept struct {
	SenderLabel types.TEXT `json:"senderLabel"`
}

// ToMap converts TransferRequestAccept to a map for DAML arguments
func (t TransferRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderLabel"] = string(t.SenderLabel)

	return m
}

func (t TransferRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestAccept to hex string (Canton MCMS format)
func (t TransferRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestAccept from hex string (Canton MCMS format)
func (t *TransferRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestAcceptResult is a Record type
type TransferRequestAcceptResult struct {
	AcceptedTransferCid types.CONTRACT_ID `json:"acceptedTransferCid"`
}

// ToMap converts TransferRequestAcceptResult to a map for DAML arguments
func (t TransferRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedTransferCid"] = model.NestedToDAMLValue(t.AcceptedTransferCid)

	return m
}

func (t TransferRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestAcceptResult to hex string (Canton MCMS format)
func (t TransferRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestAcceptResult from hex string (Canton MCMS format)
func (t *TransferRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestCancel is a Record type
type TransferRequestCancel struct {
}

// ToMap converts TransferRequestCancel to a map for DAML arguments
func (t TransferRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t TransferRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestCancel to hex string (Canton MCMS format)
func (t TransferRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestCancel from hex string (Canton MCMS format)
func (t *TransferRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestCancelResult is a Record type
type TransferRequestCancelResult struct {
}

// ToMap converts TransferRequestCancelResult to a map for DAML arguments
func (t TransferRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t TransferRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestCancelResult to hex string (Canton MCMS format)
func (t TransferRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestCancelResult from hex string (Canton MCMS format)
func (t *TransferRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestClean is a Record type
type TransferRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts TransferRequestClean to a map for DAML arguments
func (t TransferRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t TransferRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestClean to hex string (Canton MCMS format)
func (t TransferRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestClean from hex string (Canton MCMS format)
func (t *TransferRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestReject is a Record type
type TransferRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts TransferRequestReject to a map for DAML arguments
func (t TransferRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t TransferRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestReject to hex string (Canton MCMS format)
func (t TransferRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestReject from hex string (Canton MCMS format)
func (t *TransferRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRequestRejectResult is a Record type
type TransferRequestRejectResult struct {
	RejectedTransferCid types.CONTRACT_ID `json:"rejectedTransferCid"`
}

// ToMap converts TransferRequestRejectResult to a map for DAML arguments
func (t TransferRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedTransferCid"] = model.NestedToDAMLValue(t.RejectedTransferCid)

	return m
}

func (t TransferRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRequestRejectResult to hex string (Canton MCMS format)
func (t TransferRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRequestRejectResult from hex string (Canton MCMS format)
func (t *TransferRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRule is a Template type
type TransferRule struct {
	Operator  types.PARTY `json:"operator"`
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TransferRule) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TransferRule) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TransferRule) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TransferRule) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TransferRule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRule to hex string (Canton MCMS format)
func (t TransferRule) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRule from hex string (Canton MCMS format)
func (t *TransferRule) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TransferRule

// TransferRuleDirectTransfer exercises the TransferRule_DirectTransfer choice on this TransferRule contract
// This method uses the package name in the template ID
func (t TransferRule) TransferRuleDirectTransfer(contractID string, args TransferRuleDirectTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_DirectTransfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleDirectTransferWithPackageID exercises the TransferRule_DirectTransfer choice using the provided package ID instead of package name
func (t TransferRule) TransferRuleDirectTransferWithPackageID(contractID string, packageID string, args TransferRuleDirectTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_DirectTransfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleTwoStepTransfer exercises the TransferRule_TwoStepTransfer choice on this TransferRule contract
// This method uses the package name in the template ID
func (t TransferRule) TransferRuleTwoStepTransfer(contractID string, args TransferRuleTwoStepTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_TwoStepTransfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleTwoStepTransferWithPackageID exercises the TransferRule_TwoStepTransfer choice using the provided package ID instead of package name
func (t TransferRule) TransferRuleTwoStepTransferWithPackageID(contractID string, packageID string, args TransferRuleTwoStepTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_TwoStepTransfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleExecuteAllocation exercises the TransferRule_ExecuteAllocation choice on this TransferRule contract
// This method uses the package name in the template ID
func (t TransferRule) TransferRuleExecuteAllocation(contractID string, args TransferRuleExecuteAllocation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_ExecuteAllocation",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleExecuteAllocationWithPackageID exercises the TransferRule_ExecuteAllocation choice using the provided package ID instead of package name
func (t TransferRule) TransferRuleExecuteAllocationWithPackageID(contractID string, packageID string, args TransferRuleExecuteAllocation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_ExecuteAllocation",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferRule contract
// This method uses the package name in the template ID
func (t TransferRule) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TransferRule) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferRuleTransfer exercises the TransferRule_Transfer choice on this TransferRule contract
// This method uses the package name in the template ID
func (t TransferRule) TransferRuleTransfer(contractID string, args TransferRuleTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleTransferWithPackageID exercises the TransferRule_Transfer choice using the provided package ID instead of package name
func (t TransferRule) TransferRuleTransferWithPackageID(contractID string, packageID string, args TransferRuleTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Rule.Transfer", "TransferRule"),
		ContractID: contractID,
		Choice:     "TransferRule_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferRuleAcceptTransferOfferResult is a Record type
type TransferRuleAcceptTransferOfferResult struct {
	ReceiverHoldingCid types.CONTRACT_ID  `json:"receiverHoldingCid"`
	SenderHoldingCid   *types.CONTRACT_ID `json:"senderHoldingCid" hex:"optional"`
}

// ToMap converts TransferRuleAcceptTransferOfferResult to a map for DAML arguments
func (t TransferRuleAcceptTransferOfferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCid"] = model.NestedToDAMLValue(t.ReceiverHoldingCid)

	if t.SenderHoldingCid != nil {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.SenderHoldingCid),
		}
	} else {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleAcceptTransferOfferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleAcceptTransferOfferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleAcceptTransferOfferResult to hex string (Canton MCMS format)
func (t TransferRuleAcceptTransferOfferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleAcceptTransferOfferResult from hex string (Canton MCMS format)
func (t *TransferRuleAcceptTransferOfferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleDirectTransfer is a Record type
type TransferRuleDirectTransfer struct {
	Transfer         splice_api_token_transfer_instruction_v1.Transfer `json:"transfer"`
	ExtraArgs        splice_api_token_metadata_v1.ExtraArgs             `json:"extraArgs"`
	ExpectedOperator types.PARTY                                        `json:"expectedOperator"`
	ExpectedProvider *types.PARTY                                       `json:"expectedProvider" hex:"optional"`
}

// ToMap converts TransferRuleDirectTransfer to a map for DAML arguments
func (t TransferRuleDirectTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["transfer"] = model.NestedToDAMLValue(t.Transfer)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	m["expectedOperator"] = t.ExpectedOperator.ToMap()

	if t.ExpectedProvider != nil {
		m["expectedProvider"] = map[string]any{
			"_type": "optional",
			"value": (*t.ExpectedProvider).ToMap(),
		}
	} else {
		m["expectedProvider"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleDirectTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleDirectTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleDirectTransfer to hex string (Canton MCMS format)
func (t TransferRuleDirectTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleDirectTransfer from hex string (Canton MCMS format)
func (t *TransferRuleDirectTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleDirectTransferResult is a Record type
type TransferRuleDirectTransferResult struct {
	ReceiverHoldingCid types.CONTRACT_ID  `json:"receiverHoldingCid"`
	SenderHoldingCid   *types.CONTRACT_ID `json:"senderHoldingCid" hex:"optional"`
}

// ToMap converts TransferRuleDirectTransferResult to a map for DAML arguments
func (t TransferRuleDirectTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCid"] = model.NestedToDAMLValue(t.ReceiverHoldingCid)

	if t.SenderHoldingCid != nil {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.SenderHoldingCid),
		}
	} else {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleDirectTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleDirectTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleDirectTransferResult to hex string (Canton MCMS format)
func (t TransferRuleDirectTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleDirectTransferResult from hex string (Canton MCMS format)
func (t *TransferRuleDirectTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleExecuteAllocation is a Record type
type TransferRuleExecuteAllocation struct {
	Allocation       splice_api_token_allocation_v1.AllocationView `json:"allocation"`
	ExtraArgs        splice_api_token_metadata_v1.ExtraArgs        `json:"extraArgs"`
	ExpectedOperator types.PARTY                                   `json:"expectedOperator"`
	ExpectedProvider types.PARTY                                   `json:"expectedProvider"`
}

// ToMap converts TransferRuleExecuteAllocation to a map for DAML arguments
func (t TransferRuleExecuteAllocation) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocation"] = model.NestedToDAMLValue(t.Allocation)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	m["expectedOperator"] = t.ExpectedOperator.ToMap()

	m["expectedProvider"] = t.ExpectedProvider.ToMap()

	return m
}

func (t TransferRuleExecuteAllocation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleExecuteAllocation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleExecuteAllocation to hex string (Canton MCMS format)
func (t TransferRuleExecuteAllocation) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleExecuteAllocation from hex string (Canton MCMS format)
func (t *TransferRuleExecuteAllocation) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleExecuteAllocationResult is a Record type
type TransferRuleExecuteAllocationResult struct {
	ReceiverHoldingCid types.CONTRACT_ID  `json:"receiverHoldingCid"`
	SenderHoldingCid   *types.CONTRACT_ID `json:"senderHoldingCid" hex:"optional"`
}

// ToMap converts TransferRuleExecuteAllocationResult to a map for DAML arguments
func (t TransferRuleExecuteAllocationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCid"] = model.NestedToDAMLValue(t.ReceiverHoldingCid)

	if t.SenderHoldingCid != nil {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.SenderHoldingCid),
		}
	} else {
		m["senderHoldingCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleExecuteAllocationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleExecuteAllocationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleExecuteAllocationResult to hex string (Canton MCMS format)
func (t TransferRuleExecuteAllocationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleExecuteAllocationResult from hex string (Canton MCMS format)
func (t *TransferRuleExecuteAllocationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleTransfer is a Record type
type TransferRuleTransfer struct {
	Transfer                   splice_api_token_transfer_instruction_v1.Transfer `json:"transfer"`
	InstrumentConfigurationCid types.CONTRACT_ID                                  `json:"instrumentConfigurationCid"`
	SenderCredentialCids       []types.CONTRACT_ID                                `json:"senderCredentialCids"`
	ReceiverCredentialCids     []types.CONTRACT_ID                                `json:"receiverCredentialCids"`
	AppRewardConfigurationCid  *types.CONTRACT_ID                                 `json:"appRewardConfigurationCid" hex:"optional"`
	FeaturedAppRightCid        *types.CONTRACT_ID                                 `json:"featuredAppRightCid" hex:"optional"`
}

// ToMap converts TransferRuleTransfer to a map for DAML arguments
func (t TransferRuleTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["transfer"] = model.NestedToDAMLValue(t.Transfer)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	m["senderCredentialCids"] = func() []any {
		res := make([]any, 0, len(t.SenderCredentialCids))
		for _, e := range t.SenderCredentialCids {
			res = append(res, e)
		}
		return res
	}()

	m["receiverCredentialCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverCredentialCids))
		for _, e := range t.ReceiverCredentialCids {
			res = append(res, e)
		}
		return res
	}()

	if t.AppRewardConfigurationCid != nil {
		m["appRewardConfigurationCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.AppRewardConfigurationCid),
		}
	} else {
		m["appRewardConfigurationCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.FeaturedAppRightCid != nil {
		m["featuredAppRightCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.FeaturedAppRightCid),
		}
	} else {
		m["featuredAppRightCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleTransfer to hex string (Canton MCMS format)
func (t TransferRuleTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleTransfer from hex string (Canton MCMS format)
func (t *TransferRuleTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleTransferResult is a Record type
type TransferRuleTransferResult struct {
	ReceiverHoldingCid types.CONTRACT_ID  `json:"receiverHoldingCid"`
	SenderChangeCid    *types.CONTRACT_ID `json:"senderChangeCid" hex:"optional"`
}

// ToMap converts TransferRuleTransferResult to a map for DAML arguments
func (t TransferRuleTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCid"] = model.NestedToDAMLValue(t.ReceiverHoldingCid)

	if t.SenderChangeCid != nil {
		m["senderChangeCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.SenderChangeCid),
		}
	} else {
		m["senderChangeCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t TransferRuleTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleTransferResult to hex string (Canton MCMS format)
func (t TransferRuleTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleTransferResult from hex string (Canton MCMS format)
func (t *TransferRuleTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferRuleTwoStepTransfer is a Record type
type TransferRuleTwoStepTransfer struct {
	Transfer         splice_api_token_transfer_instruction_v1.Transfer `json:"transfer"`
	ExtraArgs        splice_api_token_metadata_v1.ExtraArgs             `json:"extraArgs"`
	ExpectedOperator types.PARTY                                        `json:"expectedOperator"`
	ExpectedProvider types.PARTY                                        `json:"expectedProvider"`
}

// ToMap converts TransferRuleTwoStepTransfer to a map for DAML arguments
func (t TransferRuleTwoStepTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["transfer"] = model.NestedToDAMLValue(t.Transfer)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	m["expectedOperator"] = t.ExpectedOperator.ToMap()

	m["expectedProvider"] = t.ExpectedProvider.ToMap()

	return m
}

func (t TransferRuleTwoStepTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferRuleTwoStepTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferRuleTwoStepTransfer to hex string (Canton MCMS format)
func (t TransferRuleTwoStepTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferRuleTwoStepTransfer from hex string (Canton MCMS format)
func (t *TransferRuleTwoStepTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TxKind is an enum type
type TxKind string

const (
	TxKindTxKind_Transfer TxKind = "TxKind_Transfer"

	TxKindTxKind_Unlock TxKind = "TxKind_Unlock"

	TxKindTxKind_MergeSplit TxKind = "TxKind_MergeSplit"

	TxKindTxKind_Burn TxKind = "TxKind_Burn"

	TxKindTxKind_Mint TxKind = "TxKind_Mint"

	TxKindTxKind_ExpireDust TxKind = "TxKind_ExpireDust"
)

func (e TxKind) GetEnumConstructor() string { return string(e) }

func (e TxKind) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.TokenApiUtils", "TxKind")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e TxKind) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.TokenApiUtils", "TxKind")
}

func (e TxKind) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *TxKind) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes TxKind to hex string (Canton MCMS format)
func (e TxKind) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes TxKind from hex string (Canton MCMS format)
func (e *TxKind) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = TxKind("")

// Unlock is a Record type
type Unlock struct {
	Operator             types.PARTY                              `json:"operator"`
	Provider             types.PARTY                              `json:"provider"`
	Registrar            types.PARTY                              `json:"registrar"`
	Holder               types.PARTY                              `json:"holder"`
	Locker               types.PARTY                              `json:"locker"`
	LockContext          types.TEXT                               `json:"lockContext"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                Batch                                    `json:"batch"`
}

// ToMap converts Unlock to a map for DAML arguments
func (t Unlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["operator"] = t.Operator.ToMap()

	m["provider"] = t.Provider.ToMap()

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["locker"] = t.Locker.ToMap()

	m["lockContext"] = string(t.LockContext)

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t Unlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Unlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Unlock to hex string (Canton MCMS format)
func (t Unlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Unlock from hex string (Canton MCMS format)
func (t *Unlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOffer is a Template type
type UnlockOffer struct {
	Unlock       Unlock     `json:"unlock"`
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t UnlockOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t UnlockOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t UnlockOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t UnlockOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingLabel"] = string(t.HoldingLabel)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t UnlockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOffer to hex string (Canton MCMS format)
func (t UnlockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOffer from hex string (Canton MCMS format)
func (t *UnlockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for UnlockOffer

// UnlockOfferClean exercises the UnlockOffer_Clean choice on this UnlockOffer contract
// This method uses the package name in the template ID
func (t UnlockOffer) UnlockOfferClean(contractID string, args UnlockOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferCleanWithPackageID exercises the UnlockOffer_Clean choice using the provided package ID instead of package name
func (t UnlockOffer) UnlockOfferCleanWithPackageID(contractID string, packageID string, args UnlockOfferClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Clean",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferAccept exercises the UnlockOffer_Accept choice on this UnlockOffer contract
// This method uses the package name in the template ID
func (t UnlockOffer) UnlockOfferAccept(contractID string, args UnlockOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferAcceptWithPackageID exercises the UnlockOffer_Accept choice using the provided package ID instead of package name
func (t UnlockOffer) UnlockOfferAcceptWithPackageID(contractID string, packageID string, args UnlockOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferReject exercises the UnlockOffer_Reject choice on this UnlockOffer contract
// This method uses the package name in the template ID
func (t UnlockOffer) UnlockOfferReject(contractID string, args UnlockOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferRejectWithPackageID exercises the UnlockOffer_Reject choice using the provided package ID instead of package name
func (t UnlockOffer) UnlockOfferRejectWithPackageID(contractID string, packageID string, args UnlockOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferCancel exercises the UnlockOffer_Cancel choice on this UnlockOffer contract
// This method uses the package name in the template ID
func (t UnlockOffer) UnlockOfferCancel(contractID string, args UnlockOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// UnlockOfferCancelWithPackageID exercises the UnlockOffer_Cancel choice using the provided package ID instead of package name
func (t UnlockOffer) UnlockOfferCancelWithPackageID(contractID string, packageID string, args UnlockOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "UnlockOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this UnlockOffer contract
// This method uses the package name in the template ID
func (t UnlockOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t UnlockOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// UnlockOfferAccept is a Record type
type UnlockOfferAccept struct {
}

// ToMap converts UnlockOfferAccept to a map for DAML arguments
func (t UnlockOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UnlockOfferAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferAccept to hex string (Canton MCMS format)
func (t UnlockOfferAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferAccept from hex string (Canton MCMS format)
func (t *UnlockOfferAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferAcceptResult is a Record type
type UnlockOfferAcceptResult struct {
	AcceptedUnlockCid types.CONTRACT_ID `json:"acceptedUnlockCid"`
}

// ToMap converts UnlockOfferAcceptResult to a map for DAML arguments
func (t UnlockOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedUnlockCid"] = model.NestedToDAMLValue(t.AcceptedUnlockCid)

	return m
}

func (t UnlockOfferAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferAcceptResult to hex string (Canton MCMS format)
func (t UnlockOfferAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferAcceptResult from hex string (Canton MCMS format)
func (t *UnlockOfferAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferCancel is a Record type
type UnlockOfferCancel struct {
}

// ToMap converts UnlockOfferCancel to a map for DAML arguments
func (t UnlockOfferCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UnlockOfferCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferCancel to hex string (Canton MCMS format)
func (t UnlockOfferCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferCancel from hex string (Canton MCMS format)
func (t *UnlockOfferCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferCancelResult is a Record type
type UnlockOfferCancelResult struct {
}

// ToMap converts UnlockOfferCancelResult to a map for DAML arguments
func (t UnlockOfferCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UnlockOfferCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferCancelResult to hex string (Canton MCMS format)
func (t UnlockOfferCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferCancelResult from hex string (Canton MCMS format)
func (t *UnlockOfferCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferClean is a Record type
type UnlockOfferClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts UnlockOfferClean to a map for DAML arguments
func (t UnlockOfferClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t UnlockOfferClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferClean to hex string (Canton MCMS format)
func (t UnlockOfferClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferClean from hex string (Canton MCMS format)
func (t *UnlockOfferClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferReject is a Record type
type UnlockOfferReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts UnlockOfferReject to a map for DAML arguments
func (t UnlockOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t UnlockOfferReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferReject to hex string (Canton MCMS format)
func (t UnlockOfferReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferReject from hex string (Canton MCMS format)
func (t *UnlockOfferReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockOfferRejectResult is a Record type
type UnlockOfferRejectResult struct {
	RejectedUnlockCid types.CONTRACT_ID `json:"rejectedUnlockCid"`
}

// ToMap converts UnlockOfferRejectResult to a map for DAML arguments
func (t UnlockOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedUnlockCid"] = model.NestedToDAMLValue(t.RejectedUnlockCid)

	return m
}

func (t UnlockOfferRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockOfferRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockOfferRejectResult to hex string (Canton MCMS format)
func (t UnlockOfferRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockOfferRejectResult from hex string (Canton MCMS format)
func (t *UnlockOfferRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequest is a Template type
type UnlockRequest struct {
	Unlock Unlock `json:"unlock"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t UnlockRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t UnlockRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t UnlockRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t UnlockRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["unlock"] = model.NestedToDAMLValue(t.Unlock)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t UnlockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequest to hex string (Canton MCMS format)
func (t UnlockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequest from hex string (Canton MCMS format)
func (t *UnlockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for UnlockRequest

// UnlockRequestClean exercises the UnlockRequest_Clean choice on this UnlockRequest contract
// This method uses the package name in the template ID
func (t UnlockRequest) UnlockRequestClean(contractID string, args UnlockRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestCleanWithPackageID exercises the UnlockRequest_Clean choice using the provided package ID instead of package name
func (t UnlockRequest) UnlockRequestCleanWithPackageID(contractID string, packageID string, args UnlockRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestAccept exercises the UnlockRequest_Accept choice on this UnlockRequest contract
// This method uses the package name in the template ID
func (t UnlockRequest) UnlockRequestAccept(contractID string, args UnlockRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestAcceptWithPackageID exercises the UnlockRequest_Accept choice using the provided package ID instead of package name
func (t UnlockRequest) UnlockRequestAcceptWithPackageID(contractID string, packageID string, args UnlockRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestReject exercises the UnlockRequest_Reject choice on this UnlockRequest contract
// This method uses the package name in the template ID
func (t UnlockRequest) UnlockRequestReject(contractID string, args UnlockRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestRejectWithPackageID exercises the UnlockRequest_Reject choice using the provided package ID instead of package name
func (t UnlockRequest) UnlockRequestRejectWithPackageID(contractID string, packageID string, args UnlockRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestCancel exercises the UnlockRequest_Cancel choice on this UnlockRequest contract
// This method uses the package name in the template ID
func (t UnlockRequest) UnlockRequestCancel(contractID string, args UnlockRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// UnlockRequestCancelWithPackageID exercises the UnlockRequest_Cancel choice using the provided package ID instead of package name
func (t UnlockRequest) UnlockRequestCancelWithPackageID(contractID string, packageID string, args UnlockRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "UnlockRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this UnlockRequest contract
// This method uses the package name in the template ID
func (t UnlockRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t UnlockRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.V0.Holding.Unlock", "UnlockRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// UnlockRequestAccept is a Record type
type UnlockRequestAccept struct {
	HoldingLabel types.TEXT `json:"holdingLabel"`
}

// ToMap converts UnlockRequestAccept to a map for DAML arguments
func (t UnlockRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingLabel"] = string(t.HoldingLabel)

	return m
}

func (t UnlockRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestAccept to hex string (Canton MCMS format)
func (t UnlockRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestAccept from hex string (Canton MCMS format)
func (t *UnlockRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestAcceptResult is a Record type
type UnlockRequestAcceptResult struct {
	AcceptedUnlockCid types.CONTRACT_ID `json:"acceptedUnlockCid"`
}

// ToMap converts UnlockRequestAcceptResult to a map for DAML arguments
func (t UnlockRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedUnlockCid"] = model.NestedToDAMLValue(t.AcceptedUnlockCid)

	return m
}

func (t UnlockRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestAcceptResult to hex string (Canton MCMS format)
func (t UnlockRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestAcceptResult from hex string (Canton MCMS format)
func (t *UnlockRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestCancel is a Record type
type UnlockRequestCancel struct {
}

// ToMap converts UnlockRequestCancel to a map for DAML arguments
func (t UnlockRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UnlockRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestCancel to hex string (Canton MCMS format)
func (t UnlockRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestCancel from hex string (Canton MCMS format)
func (t *UnlockRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestCancelResult is a Record type
type UnlockRequestCancelResult struct {
}

// ToMap converts UnlockRequestCancelResult to a map for DAML arguments
func (t UnlockRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t UnlockRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestCancelResult to hex string (Canton MCMS format)
func (t UnlockRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestCancelResult from hex string (Canton MCMS format)
func (t *UnlockRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestClean is a Record type
type UnlockRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts UnlockRequestClean to a map for DAML arguments
func (t UnlockRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t UnlockRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestClean to hex string (Canton MCMS format)
func (t UnlockRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestClean from hex string (Canton MCMS format)
func (t *UnlockRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestReject is a Record type
type UnlockRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts UnlockRequestReject to a map for DAML arguments
func (t UnlockRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t UnlockRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestReject to hex string (Canton MCMS format)
func (t UnlockRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestReject from hex string (Canton MCMS format)
func (t *UnlockRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UnlockRequestRejectResult is a Record type
type UnlockRequestRejectResult struct {
	RejectedUnlockCid types.CONTRACT_ID `json:"rejectedUnlockCid"`
}

// ToMap converts UnlockRequestRejectResult to a map for DAML arguments
func (t UnlockRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedUnlockCid"] = model.NestedToDAMLValue(t.RejectedUnlockCid)

	return m
}

func (t UnlockRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UnlockRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UnlockRequestRejectResult to hex string (Canton MCMS format)
func (t UnlockRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UnlockRequestRejectResult from hex string (Canton MCMS format)
func (t *UnlockRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AcceptedBurnClean(args AcceptedBurnClean) (*bind.EncodedChoice, error)
	AcceptedBurnExecute(args AcceptedBurnExecute) (*bind.EncodedChoice, error)
	AcceptedBurnFail(args AcceptedBurnFail) (*bind.EncodedChoice, error)
	AcceptedForceTransferExecute(args AcceptedForceTransferExecute) (*bind.EncodedChoice, error)
	AcceptedForceTransferFail(args AcceptedForceTransferFail) (*bind.EncodedChoice, error)
	AcceptedLockClean(args AcceptedLockClean) (*bind.EncodedChoice, error)
	AcceptedLockExecute(args AcceptedLockExecute) (*bind.EncodedChoice, error)
	AcceptedLockFail(args AcceptedLockFail) (*bind.EncodedChoice, error)
	AcceptedMintClean(args AcceptedMintClean) (*bind.EncodedChoice, error)
	AcceptedMintExecute(args AcceptedMintExecute) (*bind.EncodedChoice, error)
	AcceptedMintFail(args AcceptedMintFail) (*bind.EncodedChoice, error)
	AcceptedTransferClean(args AcceptedTransferClean) (*bind.EncodedChoice, error)
	AcceptedTransferExecute(args AcceptedTransferExecute) (*bind.EncodedChoice, error)
	AcceptedTransferFail(args AcceptedTransferFail) (*bind.EncodedChoice, error)
	AcceptedUnlockClean(args AcceptedUnlockClean) (*bind.EncodedChoice, error)
	AcceptedUnlockExecute(args AcceptedUnlockExecute) (*bind.EncodedChoice, error)
	AcceptedUnlockFail(args AcceptedUnlockFail) (*bind.EncodedChoice, error)
	AppRewardConfigurationModify(args AppRewardConfigurationModify) (*bind.EncodedChoice, error)
	BurnOfferAccept(args BurnOfferAccept) (*bind.EncodedChoice, error)
	BurnOfferCancel(args BurnOfferCancel) (*bind.EncodedChoice, error)
	BurnOfferClean(args BurnOfferClean) (*bind.EncodedChoice, error)
	BurnOfferReject(args BurnOfferReject) (*bind.EncodedChoice, error)
	BurnRequestAccept(args BurnRequestAccept) (*bind.EncodedChoice, error)
	BurnRequestCancel(args BurnRequestCancel) (*bind.EncodedChoice, error)
	BurnRequestClean(args BurnRequestClean) (*bind.EncodedChoice, error)
	BurnRequestReject(args BurnRequestReject) (*bind.EncodedChoice, error)
	ExecutedBurnClean(args ExecutedBurnClean) (*bind.EncodedChoice, error)
	ExecutedBurnDelete(args ExecutedBurnDelete) (*bind.EncodedChoice, error)
	ExecutedForceTransferDelete(args ExecutedForceTransferDelete) (*bind.EncodedChoice, error)
	ExecutedLockClean(args ExecutedLockClean) (*bind.EncodedChoice, error)
	ExecutedLockDelete(args ExecutedLockDelete) (*bind.EncodedChoice, error)
	ExecutedMintClean(args ExecutedMintClean) (*bind.EncodedChoice, error)
	ExecutedMintDelete(args ExecutedMintDelete) (*bind.EncodedChoice, error)
	ExecutedTransferDelete(args ExecutedTransferDelete) (*bind.EncodedChoice, error)
	ExecutedUnlockClean(args ExecutedUnlockClean) (*bind.EncodedChoice, error)
	ExecutedUnlockDelete(args ExecutedUnlockDelete) (*bind.EncodedChoice, error)
	FailedBurnClean(args FailedBurnClean) (*bind.EncodedChoice, error)
	FailedBurnDelete(args FailedBurnDelete) (*bind.EncodedChoice, error)
	FailedForceTransferDelete(args FailedForceTransferDelete) (*bind.EncodedChoice, error)
	FailedLockClean(args FailedLockClean) (*bind.EncodedChoice, error)
	FailedLockDelete(args FailedLockDelete) (*bind.EncodedChoice, error)
	FailedMintClean(args FailedMintClean) (*bind.EncodedChoice, error)
	FailedMintDelete(args FailedMintDelete) (*bind.EncodedChoice, error)
	FailedTransferClean(args FailedTransferClean) (*bind.EncodedChoice, error)
	FailedTransferDelete(args FailedTransferDelete) (*bind.EncodedChoice, error)
	FailedUnlockClean(args FailedUnlockClean) (*bind.EncodedChoice, error)
	FailedUnlockDelete(args FailedUnlockDelete) (*bind.EncodedChoice, error)
	ForceTransferRequestAccept(args ForceTransferRequestAccept) (*bind.EncodedChoice, error)
	ForceTransferRequestCancel(args ForceTransferRequestCancel) (*bind.EncodedChoice, error)
	ForceTransferRequestReject(args ForceTransferRequestReject) (*bind.EncodedChoice, error)
	InstrumentConfigurationGet(args InstrumentConfigurationGet) (*bind.EncodedChoice, error)
	InstrumentConfigurationSetProviderAppRewardBeneficiaries(args InstrumentConfigurationSetProviderAppRewardBeneficiaries) (*bind.EncodedChoice, error)
	LockOfferAccept(args LockOfferAccept) (*bind.EncodedChoice, error)
	LockOfferCancel(args LockOfferCancel) (*bind.EncodedChoice, error)
	LockOfferClean(args LockOfferClean) (*bind.EncodedChoice, error)
	LockOfferReject(args LockOfferReject) (*bind.EncodedChoice, error)
	LockRequestAccept(args LockRequestAccept) (*bind.EncodedChoice, error)
	LockRequestCancel(args LockRequestCancel) (*bind.EncodedChoice, error)
	LockRequestClean(args LockRequestClean) (*bind.EncodedChoice, error)
	LockRequestReject(args LockRequestReject) (*bind.EncodedChoice, error)
	MintOfferAccept(args MintOfferAccept) (*bind.EncodedChoice, error)
	MintOfferCancel(args MintOfferCancel) (*bind.EncodedChoice, error)
	MintOfferClean(args MintOfferClean) (*bind.EncodedChoice, error)
	MintOfferReject(args MintOfferReject) (*bind.EncodedChoice, error)
	MintRequestAccept(args MintRequestAccept) (*bind.EncodedChoice, error)
	MintRequestCancel(args MintRequestCancel) (*bind.EncodedChoice, error)
	MintRequestClean(args MintRequestClean) (*bind.EncodedChoice, error)
	MintRequestReject(args MintRequestReject) (*bind.EncodedChoice, error)
	RejectedBurnClean(args RejectedBurnClean) (*bind.EncodedChoice, error)
	RejectedBurnDelete(args RejectedBurnDelete) (*bind.EncodedChoice, error)
	RejectedForceTransferDelete(args RejectedForceTransferDelete) (*bind.EncodedChoice, error)
	RejectedLockClean(args RejectedLockClean) (*bind.EncodedChoice, error)
	RejectedLockDelete(args RejectedLockDelete) (*bind.EncodedChoice, error)
	RejectedMintClean(args RejectedMintClean) (*bind.EncodedChoice, error)
	RejectedMintDelete(args RejectedMintDelete) (*bind.EncodedChoice, error)
	RejectedTransferDelete(args RejectedTransferDelete) (*bind.EncodedChoice, error)
	RejectedUnlockClean(args RejectedUnlockClean) (*bind.EncodedChoice, error)
	RejectedUnlockDelete(args RejectedUnlockDelete) (*bind.EncodedChoice, error)
	TransferOfferAccept(args TransferOfferAccept) (*bind.EncodedChoice, error)
	TransferOfferCancel(args TransferOfferCancel) (*bind.EncodedChoice, error)
	TransferOfferClean(args TransferOfferClean) (*bind.EncodedChoice, error)
	TransferOfferReject(args TransferOfferReject) (*bind.EncodedChoice, error)
	TransferRequestAccept(args TransferRequestAccept) (*bind.EncodedChoice, error)
	TransferRequestCancel(args TransferRequestCancel) (*bind.EncodedChoice, error)
	TransferRequestClean(args TransferRequestClean) (*bind.EncodedChoice, error)
	TransferRequestReject(args TransferRequestReject) (*bind.EncodedChoice, error)
	TransferRuleDirectTransfer(args TransferRuleDirectTransfer) (*bind.EncodedChoice, error)
	TransferRuleExecuteAllocation(args TransferRuleExecuteAllocation) (*bind.EncodedChoice, error)
	TransferRuleTransfer(args TransferRuleTransfer) (*bind.EncodedChoice, error)
	TransferRuleTwoStepTransfer(args TransferRuleTwoStepTransfer) (*bind.EncodedChoice, error)
	UnlockOfferAccept(args UnlockOfferAccept) (*bind.EncodedChoice, error)
	UnlockOfferCancel(args UnlockOfferCancel) (*bind.EncodedChoice, error)
	UnlockOfferClean(args UnlockOfferClean) (*bind.EncodedChoice, error)
	UnlockOfferReject(args UnlockOfferReject) (*bind.EncodedChoice, error)
	UnlockRequestAccept(args UnlockRequestAccept) (*bind.EncodedChoice, error)
	UnlockRequestCancel(args UnlockRequestCancel) (*bind.EncodedChoice, error)
	UnlockRequestClean(args UnlockRequestClean) (*bind.EncodedChoice, error)
	UnlockRequestReject(args UnlockRequestReject) (*bind.EncodedChoice, error)
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

// AcceptedBurnClean encodes parameters for the AcceptedBurn_Clean choice.
func (e *encoder) AcceptedBurnClean(args AcceptedBurnClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedBurn_Clean", args)
}

// AcceptedBurnExecute encodes parameters for the AcceptedBurn_Execute choice.
func (e *encoder) AcceptedBurnExecute(args AcceptedBurnExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedBurn_Execute", args)
}

// AcceptedBurnFail encodes parameters for the AcceptedBurn_Fail choice.
func (e *encoder) AcceptedBurnFail(args AcceptedBurnFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedBurn_Fail", args)
}

// AcceptedForceTransferExecute encodes parameters for the AcceptedForceTransfer_Execute choice.
func (e *encoder) AcceptedForceTransferExecute(args AcceptedForceTransferExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedForceTransfer_Execute", args)
}

// AcceptedForceTransferFail encodes parameters for the AcceptedForceTransfer_Fail choice.
func (e *encoder) AcceptedForceTransferFail(args AcceptedForceTransferFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedForceTransfer_Fail", args)
}

// AcceptedLockClean encodes parameters for the AcceptedLock_Clean choice.
func (e *encoder) AcceptedLockClean(args AcceptedLockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedLock_Clean", args)
}

// AcceptedLockExecute encodes parameters for the AcceptedLock_Execute choice.
func (e *encoder) AcceptedLockExecute(args AcceptedLockExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedLock_Execute", args)
}

// AcceptedLockFail encodes parameters for the AcceptedLock_Fail choice.
func (e *encoder) AcceptedLockFail(args AcceptedLockFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedLock_Fail", args)
}

// AcceptedMintClean encodes parameters for the AcceptedMint_Clean choice.
func (e *encoder) AcceptedMintClean(args AcceptedMintClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedMint_Clean", args)
}

// AcceptedMintExecute encodes parameters for the AcceptedMint_Execute choice.
func (e *encoder) AcceptedMintExecute(args AcceptedMintExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedMint_Execute", args)
}

// AcceptedMintFail encodes parameters for the AcceptedMint_Fail choice.
func (e *encoder) AcceptedMintFail(args AcceptedMintFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedMint_Fail", args)
}

// AcceptedTransferClean encodes parameters for the AcceptedTransfer_Clean choice.
func (e *encoder) AcceptedTransferClean(args AcceptedTransferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedTransfer_Clean", args)
}

// AcceptedTransferExecute encodes parameters for the AcceptedTransfer_Execute choice.
func (e *encoder) AcceptedTransferExecute(args AcceptedTransferExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedTransfer_Execute", args)
}

// AcceptedTransferFail encodes parameters for the AcceptedTransfer_Fail choice.
func (e *encoder) AcceptedTransferFail(args AcceptedTransferFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedTransfer_Fail", args)
}

// AcceptedUnlockClean encodes parameters for the AcceptedUnlock_Clean choice.
func (e *encoder) AcceptedUnlockClean(args AcceptedUnlockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedUnlock_Clean", args)
}

// AcceptedUnlockExecute encodes parameters for the AcceptedUnlock_Execute choice.
func (e *encoder) AcceptedUnlockExecute(args AcceptedUnlockExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedUnlock_Execute", args)
}

// AcceptedUnlockFail encodes parameters for the AcceptedUnlock_Fail choice.
func (e *encoder) AcceptedUnlockFail(args AcceptedUnlockFail) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AcceptedUnlock_Fail", args)
}

// AppRewardConfigurationModify encodes parameters for the AppRewardConfiguration_Modify choice.
func (e *encoder) AppRewardConfigurationModify(args AppRewardConfigurationModify) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AppRewardConfiguration_Modify", args)
}

// BurnOfferAccept encodes parameters for the BurnOffer_Accept choice.
func (e *encoder) BurnOfferAccept(args BurnOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Accept", args)
}

// BurnOfferCancel encodes parameters for the BurnOffer_Cancel choice.
func (e *encoder) BurnOfferCancel(args BurnOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Cancel", args)
}

// BurnOfferClean encodes parameters for the BurnOffer_Clean choice.
func (e *encoder) BurnOfferClean(args BurnOfferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Clean", args)
}

// BurnOfferReject encodes parameters for the BurnOffer_Reject choice.
func (e *encoder) BurnOfferReject(args BurnOfferReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Reject", args)
}

// BurnRequestAccept encodes parameters for the BurnRequest_Accept choice.
func (e *encoder) BurnRequestAccept(args BurnRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnRequest_Accept", args)
}

// BurnRequestCancel encodes parameters for the BurnRequest_Cancel choice.
func (e *encoder) BurnRequestCancel(args BurnRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnRequest_Cancel", args)
}

// BurnRequestClean encodes parameters for the BurnRequest_Clean choice.
func (e *encoder) BurnRequestClean(args BurnRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnRequest_Clean", args)
}

// BurnRequestReject encodes parameters for the BurnRequest_Reject choice.
func (e *encoder) BurnRequestReject(args BurnRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnRequest_Reject", args)
}

// ExecutedBurnClean encodes parameters for the ExecutedBurn_Clean choice.
func (e *encoder) ExecutedBurnClean(args ExecutedBurnClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedBurn_Clean", args)
}

// ExecutedBurnDelete encodes parameters for the ExecutedBurn_Delete choice.
func (e *encoder) ExecutedBurnDelete(args ExecutedBurnDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedBurn_Delete", args)
}

// ExecutedForceTransferDelete encodes parameters for the ExecutedForceTransfer_Delete choice.
func (e *encoder) ExecutedForceTransferDelete(args ExecutedForceTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedForceTransfer_Delete", args)
}

// ExecutedLockClean encodes parameters for the ExecutedLock_Clean choice.
func (e *encoder) ExecutedLockClean(args ExecutedLockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedLock_Clean", args)
}

// ExecutedLockDelete encodes parameters for the ExecutedLock_Delete choice.
func (e *encoder) ExecutedLockDelete(args ExecutedLockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedLock_Delete", args)
}

// ExecutedMintClean encodes parameters for the ExecutedMint_Clean choice.
func (e *encoder) ExecutedMintClean(args ExecutedMintClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedMint_Clean", args)
}

// ExecutedMintDelete encodes parameters for the ExecutedMint_Delete choice.
func (e *encoder) ExecutedMintDelete(args ExecutedMintDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedMint_Delete", args)
}

// ExecutedTransferDelete encodes parameters for the ExecutedTransfer_Delete choice.
func (e *encoder) ExecutedTransferDelete(args ExecutedTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedTransfer_Delete", args)
}

// ExecutedUnlockClean encodes parameters for the ExecutedUnlock_Clean choice.
func (e *encoder) ExecutedUnlockClean(args ExecutedUnlockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedUnlock_Clean", args)
}

// ExecutedUnlockDelete encodes parameters for the ExecutedUnlock_Delete choice.
func (e *encoder) ExecutedUnlockDelete(args ExecutedUnlockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedUnlock_Delete", args)
}

// FailedBurnClean encodes parameters for the FailedBurn_Clean choice.
func (e *encoder) FailedBurnClean(args FailedBurnClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedBurn_Clean", args)
}

// FailedBurnDelete encodes parameters for the FailedBurn_Delete choice.
func (e *encoder) FailedBurnDelete(args FailedBurnDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedBurn_Delete", args)
}

// FailedForceTransferDelete encodes parameters for the FailedForceTransfer_Delete choice.
func (e *encoder) FailedForceTransferDelete(args FailedForceTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedForceTransfer_Delete", args)
}

// FailedLockClean encodes parameters for the FailedLock_Clean choice.
func (e *encoder) FailedLockClean(args FailedLockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedLock_Clean", args)
}

// FailedLockDelete encodes parameters for the FailedLock_Delete choice.
func (e *encoder) FailedLockDelete(args FailedLockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedLock_Delete", args)
}

// FailedMintClean encodes parameters for the FailedMint_Clean choice.
func (e *encoder) FailedMintClean(args FailedMintClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedMint_Clean", args)
}

// FailedMintDelete encodes parameters for the FailedMint_Delete choice.
func (e *encoder) FailedMintDelete(args FailedMintDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedMint_Delete", args)
}

// FailedTransferClean encodes parameters for the FailedTransfer_Clean choice.
func (e *encoder) FailedTransferClean(args FailedTransferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedTransfer_Clean", args)
}

// FailedTransferDelete encodes parameters for the FailedTransfer_Delete choice.
func (e *encoder) FailedTransferDelete(args FailedTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedTransfer_Delete", args)
}

// FailedUnlockClean encodes parameters for the FailedUnlock_Clean choice.
func (e *encoder) FailedUnlockClean(args FailedUnlockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedUnlock_Clean", args)
}

// FailedUnlockDelete encodes parameters for the FailedUnlock_Delete choice.
func (e *encoder) FailedUnlockDelete(args FailedUnlockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FailedUnlock_Delete", args)
}

// ForceTransferRequestAccept encodes parameters for the ForceTransferRequest_Accept choice.
func (e *encoder) ForceTransferRequestAccept(args ForceTransferRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ForceTransferRequest_Accept", args)
}

// ForceTransferRequestCancel encodes parameters for the ForceTransferRequest_Cancel choice.
func (e *encoder) ForceTransferRequestCancel(args ForceTransferRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ForceTransferRequest_Cancel", args)
}

// ForceTransferRequestReject encodes parameters for the ForceTransferRequest_Reject choice.
func (e *encoder) ForceTransferRequestReject(args ForceTransferRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ForceTransferRequest_Reject", args)
}

// InstrumentConfigurationGet encodes parameters for the InstrumentConfiguration_Get choice.
func (e *encoder) InstrumentConfigurationGet(args InstrumentConfigurationGet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("InstrumentConfiguration_Get", args)
}

// InstrumentConfigurationSetProviderAppRewardBeneficiaries encodes parameters for the InstrumentConfiguration_SetProviderAppRewardBeneficiaries choice.
func (e *encoder) InstrumentConfigurationSetProviderAppRewardBeneficiaries(args InstrumentConfigurationSetProviderAppRewardBeneficiaries) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("InstrumentConfiguration_SetProviderAppRewardBeneficiaries", args)
}

// LockOfferAccept encodes parameters for the LockOffer_Accept choice.
func (e *encoder) LockOfferAccept(args LockOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOffer_Accept", args)
}

// LockOfferCancel encodes parameters for the LockOffer_Cancel choice.
func (e *encoder) LockOfferCancel(args LockOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOffer_Cancel", args)
}

// LockOfferClean encodes parameters for the LockOffer_Clean choice.
func (e *encoder) LockOfferClean(args LockOfferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOffer_Clean", args)
}

// LockOfferReject encodes parameters for the LockOffer_Reject choice.
func (e *encoder) LockOfferReject(args LockOfferReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockOffer_Reject", args)
}

// LockRequestAccept encodes parameters for the LockRequest_Accept choice.
func (e *encoder) LockRequestAccept(args LockRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockRequest_Accept", args)
}

// LockRequestCancel encodes parameters for the LockRequest_Cancel choice.
func (e *encoder) LockRequestCancel(args LockRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockRequest_Cancel", args)
}

// LockRequestClean encodes parameters for the LockRequest_Clean choice.
func (e *encoder) LockRequestClean(args LockRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockRequest_Clean", args)
}

// LockRequestReject encodes parameters for the LockRequest_Reject choice.
func (e *encoder) LockRequestReject(args LockRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("LockRequest_Reject", args)
}

// MintOfferAccept encodes parameters for the MintOffer_Accept choice.
func (e *encoder) MintOfferAccept(args MintOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Accept", args)
}

// MintOfferCancel encodes parameters for the MintOffer_Cancel choice.
func (e *encoder) MintOfferCancel(args MintOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Cancel", args)
}

// MintOfferClean encodes parameters for the MintOffer_Clean choice.
func (e *encoder) MintOfferClean(args MintOfferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Clean", args)
}

// MintOfferReject encodes parameters for the MintOffer_Reject choice.
func (e *encoder) MintOfferReject(args MintOfferReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Reject", args)
}

// MintRequestAccept encodes parameters for the MintRequest_Accept choice.
func (e *encoder) MintRequestAccept(args MintRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRequest_Accept", args)
}

// MintRequestCancel encodes parameters for the MintRequest_Cancel choice.
func (e *encoder) MintRequestCancel(args MintRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRequest_Cancel", args)
}

// MintRequestClean encodes parameters for the MintRequest_Clean choice.
func (e *encoder) MintRequestClean(args MintRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRequest_Clean", args)
}

// MintRequestReject encodes parameters for the MintRequest_Reject choice.
func (e *encoder) MintRequestReject(args MintRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRequest_Reject", args)
}

// RejectedBurnClean encodes parameters for the RejectedBurn_Clean choice.
func (e *encoder) RejectedBurnClean(args RejectedBurnClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedBurn_Clean", args)
}

// RejectedBurnDelete encodes parameters for the RejectedBurn_Delete choice.
func (e *encoder) RejectedBurnDelete(args RejectedBurnDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedBurn_Delete", args)
}

// RejectedForceTransferDelete encodes parameters for the RejectedForceTransfer_Delete choice.
func (e *encoder) RejectedForceTransferDelete(args RejectedForceTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedForceTransfer_Delete", args)
}

// RejectedLockClean encodes parameters for the RejectedLock_Clean choice.
func (e *encoder) RejectedLockClean(args RejectedLockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedLock_Clean", args)
}

// RejectedLockDelete encodes parameters for the RejectedLock_Delete choice.
func (e *encoder) RejectedLockDelete(args RejectedLockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedLock_Delete", args)
}

// RejectedMintClean encodes parameters for the RejectedMint_Clean choice.
func (e *encoder) RejectedMintClean(args RejectedMintClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedMint_Clean", args)
}

// RejectedMintDelete encodes parameters for the RejectedMint_Delete choice.
func (e *encoder) RejectedMintDelete(args RejectedMintDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedMint_Delete", args)
}

// RejectedTransferDelete encodes parameters for the RejectedTransfer_Delete choice.
func (e *encoder) RejectedTransferDelete(args RejectedTransferDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedTransfer_Delete", args)
}

// RejectedUnlockClean encodes parameters for the RejectedUnlock_Clean choice.
func (e *encoder) RejectedUnlockClean(args RejectedUnlockClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedUnlock_Clean", args)
}

// RejectedUnlockDelete encodes parameters for the RejectedUnlock_Delete choice.
func (e *encoder) RejectedUnlockDelete(args RejectedUnlockDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedUnlock_Delete", args)
}

// TransferOfferAccept encodes parameters for the TransferOffer_Accept choice.
func (e *encoder) TransferOfferAccept(args TransferOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferOffer_Accept", args)
}

// TransferOfferCancel encodes parameters for the TransferOffer_Cancel choice.
func (e *encoder) TransferOfferCancel(args TransferOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferOffer_Cancel", args)
}

// TransferOfferClean encodes parameters for the TransferOffer_Clean choice.
func (e *encoder) TransferOfferClean(args TransferOfferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferOffer_Clean", args)
}

// TransferOfferReject encodes parameters for the TransferOffer_Reject choice.
func (e *encoder) TransferOfferReject(args TransferOfferReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferOffer_Reject", args)
}

// TransferRequestAccept encodes parameters for the TransferRequest_Accept choice.
func (e *encoder) TransferRequestAccept(args TransferRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRequest_Accept", args)
}

// TransferRequestCancel encodes parameters for the TransferRequest_Cancel choice.
func (e *encoder) TransferRequestCancel(args TransferRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRequest_Cancel", args)
}

// TransferRequestClean encodes parameters for the TransferRequest_Clean choice.
func (e *encoder) TransferRequestClean(args TransferRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRequest_Clean", args)
}

// TransferRequestReject encodes parameters for the TransferRequest_Reject choice.
func (e *encoder) TransferRequestReject(args TransferRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRequest_Reject", args)
}

// TransferRuleDirectTransfer encodes parameters for the TransferRule_DirectTransfer choice.
func (e *encoder) TransferRuleDirectTransfer(args TransferRuleDirectTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRule_DirectTransfer", args)
}

// TransferRuleExecuteAllocation encodes parameters for the TransferRule_ExecuteAllocation choice.
func (e *encoder) TransferRuleExecuteAllocation(args TransferRuleExecuteAllocation) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRule_ExecuteAllocation", args)
}

// TransferRuleTransfer encodes parameters for the TransferRule_Transfer choice.
func (e *encoder) TransferRuleTransfer(args TransferRuleTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRule_Transfer", args)
}

// TransferRuleTwoStepTransfer encodes parameters for the TransferRule_TwoStepTransfer choice.
func (e *encoder) TransferRuleTwoStepTransfer(args TransferRuleTwoStepTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferRule_TwoStepTransfer", args)
}

// UnlockOfferAccept encodes parameters for the UnlockOffer_Accept choice.
func (e *encoder) UnlockOfferAccept(args UnlockOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockOffer_Accept", args)
}

// UnlockOfferCancel encodes parameters for the UnlockOffer_Cancel choice.
func (e *encoder) UnlockOfferCancel(args UnlockOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockOffer_Cancel", args)
}

// UnlockOfferClean encodes parameters for the UnlockOffer_Clean choice.
func (e *encoder) UnlockOfferClean(args UnlockOfferClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockOffer_Clean", args)
}

// UnlockOfferReject encodes parameters for the UnlockOffer_Reject choice.
func (e *encoder) UnlockOfferReject(args UnlockOfferReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockOffer_Reject", args)
}

// UnlockRequestAccept encodes parameters for the UnlockRequest_Accept choice.
func (e *encoder) UnlockRequestAccept(args UnlockRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockRequest_Accept", args)
}

// UnlockRequestCancel encodes parameters for the UnlockRequest_Cancel choice.
func (e *encoder) UnlockRequestCancel(args UnlockRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockRequest_Cancel", args)
}

// UnlockRequestClean encodes parameters for the UnlockRequest_Clean choice.
func (e *encoder) UnlockRequestClean(args UnlockRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockRequest_Clean", args)
}

// UnlockRequestReject encodes parameters for the UnlockRequest_Reject choice.
func (e *encoder) UnlockRequestReject(args UnlockRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UnlockRequest_Reject", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
