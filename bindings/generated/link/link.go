package link

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_burn_mint_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_burn_mint_v1"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	splice_api_token_transfer_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
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
	PackageName = "link"
	PackageID   = "d36cd3c181d74fb48d66e9805fd71a27afd3745a22437d266da543a28e38f99d"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	TransferPreapprovalContextKey = types.TEXT("transfer-preapproval")
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

// Cancel is a Record type
type Cancel struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts Cancel to a map for DAML arguments
func (t Cancel) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t Cancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Cancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Cancel to hex string (Canton MCMS format)
func (t Cancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Cancel from hex string (Canton MCMS format)
func (t *Cancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LinkHolding is a Template type
type LinkHolding struct {
	HoldingOwner        types.PARTY                              `json:"holdingOwner"`
	HoldingAdmin        types.PARTY                              `json:"holdingAdmin"`
	HoldingInstrumentId splice_api_token_holding_v1.InstrumentId `json:"holdingInstrumentId"`
	HoldingAmount       types.NUMERIC                            `json:"holdingAmount"`
	Meta                splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkHolding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkHolding")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkHolding) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.Token", "LinkHolding")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkHolding) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingOwner"] = t.HoldingOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingAdmin"] = t.HoldingAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingInstrumentId"] = model.NestedToDAMLValue(t.HoldingInstrumentId)

	if t.HoldingAmount != "" {
		args["holdingAmount"] = t.HoldingAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkHolding) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingOwner"] = t.HoldingOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingAdmin"] = t.HoldingAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holdingInstrumentId"] = model.NestedToDAMLValue(t.HoldingInstrumentId)

	if t.HoldingAmount != "" {
		args["holdingAmount"] = t.HoldingAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkHolding to hex string (Canton MCMS format)
func (t LinkHolding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkHolding from hex string (Canton MCMS format)
func (t *LinkHolding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkHolding

// Archive exercises the Archive choice on this LinkHolding contract via the IHolding interface
// This method uses the package name in the template ID
func (t LinkHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkHolding) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Verify interface implementations for LinkHolding

var _ splice_api_token_holding_v1.IHolding = (*LinkHolding)(nil)

// LinkRegistry is a Template type
type LinkRegistry struct {
	RegistryAdmin        types.PARTY                              `json:"registryAdmin"`
	RegistryInstrumentId splice_api_token_holding_v1.InstrumentId `json:"registryInstrumentId"`
	InstanceId           types.TEXT                               `json:"instanceId"`
	RegistryMeta         splice_api_token_metadata_v1.Metadata    `json:"registryMeta"`
	TransferPreapprovals map[types.PARTY]types.CONTRACT_ID        `json:"transferPreapprovals"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkRegistry")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkRegistry) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.Token", "LinkRegistry")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryAdmin"] = t.RegistryAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryInstrumentId"] = model.NestedToDAMLValue(t.RegistryInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryMeta"] = model.NestedToDAMLValue(t.RegistryMeta)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferPreapprovals"] = func() any {
		if t.TransferPreapprovals == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TransferPreapprovals}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryAdmin"] = t.RegistryAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryInstrumentId"] = model.NestedToDAMLValue(t.RegistryInstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registryMeta"] = model.NestedToDAMLValue(t.RegistryMeta)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferPreapprovals"] = func() any {
		if t.TransferPreapprovals == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.TransferPreapprovals}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkRegistry to hex string (Canton MCMS format)
func (t LinkRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkRegistry from hex string (Canton MCMS format)
func (t *LinkRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkRegistry

// SetTransferPreapproval exercises the SetTransferPreapproval choice on this LinkRegistry contract
// This method uses the package name in the template ID
func (t LinkRegistry) SetTransferPreapproval(contractID string, args SetTransferPreapproval) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkRegistry"),
		ContractID: contractID,
		Choice:     "SetTransferPreapproval",
		Arguments:  argsToMap(args),
	}
}

// SetTransferPreapprovalWithPackageID exercises the SetTransferPreapproval choice using the provided package ID instead of package name
func (t LinkRegistry) SetTransferPreapprovalWithPackageID(contractID string, packageID string, args SetTransferPreapproval) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "LinkRegistry"),
		ContractID: contractID,
		Choice:     "SetTransferPreapproval",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t LinkRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this LinkRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t LinkRegistry) BurnMintFactoryPublicFetch(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetchWithPackageID exercises the BurnMintFactory_PublicFetch choice using the provided package ID instead of package name
func (t LinkRegistry) BurnMintFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this LinkRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t LinkRegistry) BurnMintFactoryBurnMint(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMintWithPackageID exercises the BurnMintFactory_BurnMint choice using the provided package ID instead of package name
func (t LinkRegistry) BurnMintFactoryBurnMintWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this LinkRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t LinkRegistry) TransferFactoryTransfer(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransferWithPackageID exercises the TransferFactory_Transfer choice using the provided package ID instead of package name
func (t LinkRegistry) TransferFactoryTransferWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this LinkRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t LinkRegistry) TransferFactoryPublicFetch(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetchWithPackageID exercises the TransferFactory_PublicFetch choice using the provided package ID instead of package name
func (t LinkRegistry) TransferFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for LinkRegistry

var _ splice_api_token_burn_mint_v1.IBurnMintFactory = (*LinkRegistry)(nil)

var _ splice_api_token_transfer_instruction_v1.ITransferFactory = (*LinkRegistry)(nil)

// LinkTransferInstruction is a Template type
type LinkTransferInstruction struct {
	InstructionAdmin            types.PARTY                                       `json:"instructionAdmin"`
	InstructionTransfer         splice_api_token_transfer_instruction_v1.Transfer `json:"instructionTransfer"`
	InstructionLockedHoldingCid types.CONTRACT_ID                                 `json:"instructionLockedHoldingCid"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkTransferInstruction) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkTransferInstruction")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkTransferInstruction) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.Token", "LinkTransferInstruction")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkTransferInstruction) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionAdmin"] = t.InstructionAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionTransfer"] = model.NestedToDAMLValue(t.InstructionTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionLockedHoldingCid"] = model.NestedToDAMLValue(t.InstructionLockedHoldingCid)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkTransferInstruction) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionAdmin"] = t.InstructionAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionTransfer"] = model.NestedToDAMLValue(t.InstructionTransfer)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instructionLockedHoldingCid"] = model.NestedToDAMLValue(t.InstructionLockedHoldingCid)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkTransferInstruction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkTransferInstruction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkTransferInstruction to hex string (Canton MCMS format)
func (t LinkTransferInstruction) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkTransferInstruction from hex string (Canton MCMS format)
func (t *LinkTransferInstruction) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkTransferInstruction

// Archive exercises the Archive choice on this LinkTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t LinkTransferInstruction) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkTransferInstruction) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this LinkTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t LinkTransferInstruction) TransferInstructionAccept(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionAcceptWithPackageID exercises the TransferInstruction_Accept choice using the provided package ID instead of package name
func (t LinkTransferInstruction) TransferInstructionAcceptWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this LinkTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t LinkTransferInstruction) TransferInstructionReject(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionRejectWithPackageID exercises the TransferInstruction_Reject choice using the provided package ID instead of package name
func (t LinkTransferInstruction) TransferInstructionRejectWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this LinkTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t LinkTransferInstruction) TransferInstructionWithdraw(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdrawWithPackageID exercises the TransferInstruction_Withdraw choice using the provided package ID instead of package name
func (t LinkTransferInstruction) TransferInstructionWithdrawWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this LinkTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t LinkTransferInstruction) TransferInstructionUpdate(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdateWithPackageID exercises the TransferInstruction_Update choice using the provided package ID instead of package name
func (t LinkTransferInstruction) TransferInstructionUpdateWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for LinkTransferInstruction

var _ splice_api_token_transfer_instruction_v1.ITransferInstruction = (*LinkTransferInstruction)(nil)

// LinkTransferPreapproval is a Template type
type LinkTransferPreapproval struct {
	PreapprovalAdmin    types.PARTY `json:"preapprovalAdmin"`
	PreapprovalReceiver types.PARTY `json:"preapprovalReceiver"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkTransferPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkTransferPreapproval")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkTransferPreapproval) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.Token", "LinkTransferPreapproval")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkTransferPreapproval) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["preapprovalAdmin"] = t.PreapprovalAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["preapprovalReceiver"] = t.PreapprovalReceiver.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkTransferPreapproval) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["preapprovalAdmin"] = t.PreapprovalAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["preapprovalReceiver"] = t.PreapprovalReceiver.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkTransferPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkTransferPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkTransferPreapproval to hex string (Canton MCMS format)
func (t LinkTransferPreapproval) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkTransferPreapproval from hex string (Canton MCMS format)
func (t *LinkTransferPreapproval) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkTransferPreapproval

// Send exercises the Send choice on this LinkTransferPreapproval contract
// This method uses the package name in the template ID
func (t LinkTransferPreapproval) Send(contractID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// SendWithPackageID exercises the Send choice using the provided package ID instead of package name
func (t LinkTransferPreapproval) SendWithPackageID(contractID string, packageID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// Cancel exercises the Cancel choice on this LinkTransferPreapproval contract
// This method uses the package name in the template ID
func (t LinkTransferPreapproval) Cancel(contractID string, args Cancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Cancel",
		Arguments:  argsToMap(args),
	}
}

// CancelWithPackageID exercises the Cancel choice using the provided package ID instead of package name
func (t LinkTransferPreapproval) CancelWithPackageID(contractID string, packageID string, args Cancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkTransferPreapproval contract
// This method uses the package name in the template ID
func (t LinkTransferPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkTransferPreapproval) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "LinkTransferPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// LockedLinkHolding is a Template type
type LockedLinkHolding struct {
	LockedAdmin        types.PARTY                              `json:"lockedAdmin"`
	LockedSender       types.PARTY                              `json:"lockedSender"`
	LockedReceiver     types.PARTY                              `json:"lockedReceiver"`
	LockedInstrumentId splice_api_token_holding_v1.InstrumentId `json:"lockedInstrumentId"`
	LockedAmount       types.NUMERIC                            `json:"lockedAmount"`
	Meta               splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LockedLinkHolding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LockedLinkHolding")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LockedLinkHolding) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.Token", "LockedLinkHolding")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LockedLinkHolding) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedAdmin"] = t.LockedAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedSender"] = t.LockedSender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedReceiver"] = t.LockedReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedInstrumentId"] = model.NestedToDAMLValue(t.LockedInstrumentId)

	if t.LockedAmount != "" {
		args["lockedAmount"] = t.LockedAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LockedLinkHolding) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedAdmin"] = t.LockedAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedSender"] = t.LockedSender.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedReceiver"] = t.LockedReceiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedInstrumentId"] = model.NestedToDAMLValue(t.LockedInstrumentId)

	if t.LockedAmount != "" {
		args["lockedAmount"] = t.LockedAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LockedLinkHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LockedLinkHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LockedLinkHolding to hex string (Canton MCMS format)
func (t LockedLinkHolding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LockedLinkHolding from hex string (Canton MCMS format)
func (t *LockedLinkHolding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LockedLinkHolding

// Archive exercises the Archive choice on this LockedLinkHolding contract
// This method uses the package name in the template ID
func (t LockedLinkHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.Token", "LockedLinkHolding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LockedLinkHolding) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.Token", "LockedLinkHolding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Send is a Record type
type Send struct {
	Sender           types.PARTY                              `json:"sender"`
	Amount           types.NUMERIC                            `json:"amount"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	InputHoldingCids []types.CONTRACT_ID                      `json:"inputHoldingCids"`
	Meta             splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["sender"] = t.Sender.ToMap()

	m["amount"] = t.Amount

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t Send) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Send) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Send to hex string (Canton MCMS format)
func (t Send) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Send from hex string (Canton MCMS format)
func (t *Send) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetTransferPreapproval is a Record type
type SetTransferPreapproval struct {
	Receiver          types.PARTY        `json:"receiver"`
	NewPreapprovalCid *types.CONTRACT_ID `json:"newPreapprovalCid" hex:"optional"`
}

// ToMap converts SetTransferPreapproval to a map for DAML arguments
func (t SetTransferPreapproval) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = t.Receiver.ToMap()

	if t.NewPreapprovalCid != nil {
		m["newPreapprovalCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.NewPreapprovalCid),
		}
	} else {
		m["newPreapprovalCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetTransferPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetTransferPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetTransferPreapproval to hex string (Canton MCMS format)
func (t SetTransferPreapproval) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetTransferPreapproval from hex string (Canton MCMS format)
func (t *SetTransferPreapproval) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	Cancel(args Cancel) (*bind.EncodedChoice, error)
	Send(args Send) (*bind.EncodedChoice, error)
	SetTransferPreapproval(args SetTransferPreapproval) (*bind.EncodedChoice, error)
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

// Cancel encodes parameters for the Cancel choice.
func (e *encoder) Cancel(args Cancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Cancel", args)
}

// Send encodes parameters for the Send choice.
func (e *encoder) Send(args Send) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Send", args)
}

// SetTransferPreapproval encodes parameters for the SetTransferPreapproval choice.
func (e *encoder) SetTransferPreapproval(args SetTransferPreapproval) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetTransferPreapproval", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
