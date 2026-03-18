package mcmstest

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	PackageName = "mcms-test"
	PackageID   = "31e900835409a728eb05a9a0bd3e883f4db558d6e8fc34944cd8dc120b273650"
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

// Counter is a Template type
type Counter struct {
	Owner      types.PARTY `json:"owner"`
	InstanceId types.TEXT  `json:"instanceId"`
	Value      types.INT64 `json:"value"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t Counter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t Counter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter")
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

// Archive exercises the Archive choice on this Counter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Counter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t Counter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetValue exercises the GetValue choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetValue(contractID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetValue",
		Arguments:  argsToMap(args),
	}
}

// GetValueWithPackageID exercises the GetValue choice using the provided package ID instead of package name
func (t Counter) GetValueWithPackageID(contractID string, packageID string, args GetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetValue",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceIdChoice exercises the GetInstanceIdChoice choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetInstanceIdChoice(contractID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceIdChoiceWithPackageID exercises the GetInstanceIdChoice choice using the provided package ID instead of package name
func (t Counter) GetInstanceIdChoiceWithPackageID(contractID string, packageID string, args GetInstanceIdChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceIdChoice",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceAddressChoice exercises the GetInstanceAddressChoice choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) GetInstanceAddressChoice(contractID string, args GetInstanceAddressChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceAddressChoice",
		Arguments:  argsToMap(args),
	}
}

// GetInstanceAddressChoiceWithPackageID exercises the GetInstanceAddressChoice choice using the provided package ID instead of package name
func (t Counter) GetInstanceAddressChoiceWithPackageID(contractID string, packageID string, args GetInstanceAddressChoice) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "GetInstanceAddressChoice",
		Arguments:  argsToMap(args),
	}
}

// Increment exercises the Increment choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) Increment(contractID string, args Increment) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Increment",
		Arguments:  argsToMap(args),
	}
}

// IncrementWithPackageID exercises the Increment choice using the provided package ID instead of package name
func (t Counter) IncrementWithPackageID(contractID string, packageID string, args Increment) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Increment",
		Arguments:  argsToMap(args),
	}
}

// SetValue exercises the SetValue choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) SetValue(contractID string, args SetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "SetValue",
		Arguments:  argsToMap(args),
	}
}

// SetValueWithPackageID exercises the SetValue choice using the provided package ID instead of package name
func (t Counter) SetValueWithPackageID(contractID string, packageID string, args SetValue) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "SetValue",
		Arguments:  argsToMap(args),
	}
}

// Reset exercises the Reset choice on this Counter contract
// This method uses the package name in the template ID
func (t Counter) Reset(contractID string, args Reset) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Reset",
		Arguments:  argsToMap(args),
	}
}

// ResetWithPackageID exercises the Reset choice using the provided package ID instead of package name
func (t Counter) ResetWithPackageID(contractID string, packageID string, args Reset) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "Counter"),
		ContractID: contractID,
		Choice:     "Reset",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this Counter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t Counter) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t Counter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.Counter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for Counter

var _ mcms.IMCMSReceiver = (*Counter)(nil)

// ForgedReceiver is a Template type
type ForgedReceiver struct {
	Owner             types.PARTY `json:"owner"`
	InstanceId        types.TEXT  `json:"instanceId"`
	ForgedPartyIdText types.TEXT  `json:"forgedPartyIdText"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ForgedReceiver) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.ForgedReceiver", "ForgedReceiver")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ForgedReceiver) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.Mock.ForgedReceiver", "ForgedReceiver")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ForgedReceiver) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forgedPartyIdText"] = string(t.ForgedPartyIdText)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ForgedReceiver) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["forgedPartyIdText"] = string(t.ForgedPartyIdText)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ForgedReceiver) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ForgedReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ForgedReceiver to hex string (Canton MCMS format)
func (t ForgedReceiver) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ForgedReceiver from hex string (Canton MCMS format)
func (t *ForgedReceiver) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ForgedReceiver

// Archive exercises the Archive choice on this ForgedReceiver contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t ForgedReceiver) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.ForgedReceiver", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ForgedReceiver) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.ForgedReceiver", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this ForgedReceiver contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t ForgedReceiver) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Mock.ForgedReceiver", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t ForgedReceiver) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Mock.ForgedReceiver", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for ForgedReceiver

var _ mcms.IMCMSReceiver = (*ForgedReceiver)(nil)

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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	GetInstanceAddressChoice(args GetInstanceAddressChoice) (*bind.EncodedChoice, error)
	GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error)
	GetValue(args GetValue) (*bind.EncodedChoice, error)
	Increment(args Increment) (*bind.EncodedChoice, error)
	Reset(args Reset) (*bind.EncodedChoice, error)
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

// GetInstanceAddressChoice encodes parameters for the GetInstanceAddressChoice choice.
func (e *encoder) GetInstanceAddressChoice(args GetInstanceAddressChoice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetInstanceAddressChoice", args)
}

// GetInstanceIdChoice encodes parameters for the GetInstanceIdChoice choice.
func (e *encoder) GetInstanceIdChoice(args GetInstanceIdChoice) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetInstanceIdChoice", args)
}

// GetValue encodes parameters for the GetValue choice.
func (e *encoder) GetValue(args GetValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetValue", args)
}

// Increment encodes parameters for the Increment choice.
func (e *encoder) Increment(args Increment) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Increment", args)
}

// Reset encodes parameters for the Reset choice.
func (e *encoder) Reset(args Reset) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Reset", args)
}

// SetValue encodes parameters for the SetValue choice.
func (e *encoder) SetValue(args SetValue) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetValue", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
