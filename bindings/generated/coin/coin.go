package coin

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
	PackageName = "coin"
	PackageID   = "96f6ec9667c5ba42b8207d955a2d23328fb7f567c716fdc5891b04bd71bf64ac"
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

// CoinHolding is a Template type
type CoinHolding struct {
	View   splice_api_token_holding_v1.HoldingView `json:"view"`
	Issuer types.PARTY                             `json:"issuer"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CoinHolding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Holding", "CoinHolding")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CoinHolding) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Coin.Holding", "CoinHolding")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CoinHolding) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["view"] = model.NestedToDAMLValue(t.View)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CoinHolding) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["view"] = model.NestedToDAMLValue(t.View)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CoinHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CoinHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CoinHolding to hex string (Canton MCMS format)
func (t CoinHolding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CoinHolding from hex string (Canton MCMS format)
func (t *CoinHolding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CoinHolding

// Transfer exercises the Transfer choice on this CoinHolding contract
// This method uses the package name in the template ID
func (t CoinHolding) Transfer(contractID string, args Transfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Holding", "CoinHolding"),
		ContractID: contractID,
		Choice:     "Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferWithPackageID exercises the Transfer choice using the provided package ID instead of package name
func (t CoinHolding) TransferWithPackageID(contractID string, packageID string, args Transfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Holding", "CoinHolding"),
		ContractID: contractID,
		Choice:     "Transfer",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CoinHolding contract via the IHolding interface
// This method uses the package name in the template ID
func (t CoinHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinHolding) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Holding", "Holding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Verify interface implementations for CoinHolding

var _ splice_api_token_holding_v1.IHolding = (*CoinHolding)(nil)

// CoinRegistry is a Template type
type CoinRegistry struct {
	Issuer       types.PARTY                              `json:"issuer"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	InstanceId   types.TEXT                               `json:"instanceId"`
	Meta         splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CoinRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "CoinRegistry")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CoinRegistry) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Coin.Registry", "CoinRegistry")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CoinRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CoinRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = model.NestedToDAMLValue(t.Meta)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CoinRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CoinRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CoinRegistry to hex string (Canton MCMS format)
func (t CoinRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CoinRegistry from hex string (Canton MCMS format)
func (t *CoinRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CoinRegistry

// Archive exercises the Archive choice on this CoinRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this CoinRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) TransferFactoryTransfer(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransferWithPackageID exercises the TransferFactory_Transfer choice using the provided package ID instead of package name
func (t CoinRegistry) TransferFactoryTransferWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this CoinRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) TransferFactoryPublicFetch(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetchWithPackageID exercises the TransferFactory_PublicFetch choice using the provided package ID instead of package name
func (t CoinRegistry) TransferFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this CoinRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) BurnMintFactoryPublicFetch(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetchWithPackageID exercises the BurnMintFactory_PublicFetch choice using the provided package ID instead of package name
func (t CoinRegistry) BurnMintFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this CoinRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) BurnMintFactoryBurnMint(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMintWithPackageID exercises the BurnMintFactory_BurnMint choice using the provided package ID instead of package name
func (t CoinRegistry) BurnMintFactoryBurnMintWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CoinRegistry

var _ splice_api_token_transfer_instruction_v1.ITransferFactory = (*CoinRegistry)(nil)

var _ splice_api_token_burn_mint_v1.IBurnMintFactory = (*CoinRegistry)(nil)

// CoinTransferInstruction is a Template type
type CoinTransferInstruction struct {
	Holding       CoinHolding     `json:"holding"`
	NewOwner      types.PARTY     `json:"newOwner"`
	RequestedAt   types.TIMESTAMP `json:"requestedAt"`
	ExecuteBefore types.TIMESTAMP `json:"executeBefore"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CoinTransferInstruction) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "CoinTransferInstruction")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CoinTransferInstruction) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Coin.Transfer", "CoinTransferInstruction")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CoinTransferInstruction) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holding"] = model.NestedToDAMLValue(t.Holding)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["newOwner"] = t.NewOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requestedAt"] = t.RequestedAt

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executeBefore"] = t.ExecuteBefore

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CoinTransferInstruction) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holding"] = model.NestedToDAMLValue(t.Holding)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["newOwner"] = t.NewOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requestedAt"] = t.RequestedAt

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executeBefore"] = t.ExecuteBefore

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CoinTransferInstruction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CoinTransferInstruction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CoinTransferInstruction to hex string (Canton MCMS format)
func (t CoinTransferInstruction) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CoinTransferInstruction from hex string (Canton MCMS format)
func (t *CoinTransferInstruction) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CoinTransferInstruction

// Archive exercises the Archive choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinTransferInstruction) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionAccept(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionAcceptWithPackageID exercises the TransferInstruction_Accept choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionAcceptWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionReject(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionRejectWithPackageID exercises the TransferInstruction_Reject choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionRejectWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionWithdraw(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdrawWithPackageID exercises the TransferInstruction_Withdraw choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionWithdrawWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionUpdate(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdateWithPackageID exercises the TransferInstruction_Update choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionUpdateWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CoinTransferInstruction

var _ splice_api_token_transfer_instruction_v1.ITransferInstruction = (*CoinTransferInstruction)(nil)

// MintPreapproval is a Template type
type MintPreapproval struct {
	Receiver types.PARTY `json:"receiver"`
	Sender   types.PARTY `json:"sender"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintPreapproval")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintPreapproval) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Coin.Registry", "MintPreapproval")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintPreapproval) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintPreapproval) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["sender"] = t.Sender.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MintPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintPreapproval to hex string (Canton MCMS format)
func (t MintPreapproval) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintPreapproval from hex string (Canton MCMS format)
func (t *MintPreapproval) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MintPreapproval

// MintPreapprovalMint exercises the MintPreapproval_Mint choice on this MintPreapproval contract
// This method uses the package name in the template ID
func (t MintPreapproval) MintPreapprovalMint(contractID string, args MintPreapprovalMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintPreapproval"),
		ContractID: contractID,
		Choice:     "MintPreapproval_Mint",
		Arguments:  argsToMap(args),
	}
}

// MintPreapprovalMintWithPackageID exercises the MintPreapproval_Mint choice using the provided package ID instead of package name
func (t MintPreapproval) MintPreapprovalMintWithPackageID(contractID string, packageID string, args MintPreapprovalMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintPreapproval"),
		ContractID: contractID,
		Choice:     "MintPreapproval_Mint",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintPreapproval contract
// This method uses the package name in the template ID
func (t MintPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintPreapproval) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintPreapprovalMint is a Record type
type MintPreapprovalMint struct {
	Issuer types.PARTY                             `json:"issuer"`
	View   splice_api_token_holding_v1.HoldingView `json:"view"`
}

// ToMap converts MintPreapprovalMint to a map for DAML arguments
func (t MintPreapprovalMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["issuer"] = t.Issuer.ToMap()

	m["view"] = model.NestedToDAMLValue(t.View)

	return m
}

func (t MintPreapprovalMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintPreapprovalMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintPreapprovalMint to hex string (Canton MCMS format)
func (t MintPreapprovalMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintPreapprovalMint from hex string (Canton MCMS format)
func (t *MintPreapprovalMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MintRole is a Template type
type MintRole struct {
	Issuer   types.PARTY       `json:"issuer"`
	Minter   types.PARTY       `json:"minter"`
	Registry types.CONTRACT_ID `json:"registry"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintRole) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintRole")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintRole) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Coin.Registry", "MintRole")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintRole) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minter"] = t.Minter.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registry"] = model.NestedToDAMLValue(t.Registry)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintRole) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minter"] = t.Minter.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registry"] = model.NestedToDAMLValue(t.Registry)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MintRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRole to hex string (Canton MCMS format)
func (t MintRole) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRole from hex string (Canton MCMS format)
func (t *MintRole) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for MintRole

// Archive exercises the Archive choice on this MintRole contract
// This method uses the package name in the template ID
func (t MintRole) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintRole) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintRoleMint exercises the MintRole_Mint choice on this MintRole contract
// This method uses the package name in the template ID
func (t MintRole) MintRoleMint(contractID string, args MintRoleMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "MintRole_Mint",
		Arguments:  argsToMap(args),
	}
}

// MintRoleMintWithPackageID exercises the MintRole_Mint choice using the provided package ID instead of package name
func (t MintRole) MintRoleMintWithPackageID(contractID string, packageID string, args MintRoleMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "MintRole_Mint",
		Arguments:  argsToMap(args),
	}
}

// MintRoleMint is a Record type
type MintRoleMint struct {
	InstrumentId splice_api_token_holding_v1.InstrumentId       `json:"instrumentId"`
	Outputs      []splice_api_token_burn_mint_v1.BurnMintOutput `json:"outputs"`
}

// ToMap converts MintRoleMint to a map for DAML arguments
func (t MintRoleMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["outputs"] = func() []any {
		res := make([]any, 0, len(t.Outputs))
		for _, e := range t.Outputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t MintRoleMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MintRoleMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MintRoleMint to hex string (Canton MCMS format)
func (t MintRoleMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MintRoleMint from hex string (Canton MCMS format)
func (t *MintRoleMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Transfer is a Record type
type Transfer struct {
	To types.PARTY `json:"to"`
}

// ToMap converts Transfer to a map for DAML arguments
func (t Transfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["to"] = t.To.ToMap()

	return m
}

func (t Transfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Transfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Transfer to hex string (Canton MCMS format)
func (t Transfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Transfer from hex string (Canton MCMS format)
func (t *Transfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	MintPreapprovalMint(args MintPreapprovalMint) (*bind.EncodedChoice, error)
	MintRoleMint(args MintRoleMint) (*bind.EncodedChoice, error)
	Transfer(args Transfer) (*bind.EncodedChoice, error)
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

// MintPreapprovalMint encodes parameters for the MintPreapproval_Mint choice.
func (e *encoder) MintPreapprovalMint(args MintPreapprovalMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintPreapproval_Mint", args)
}

// MintRoleMint encodes parameters for the MintRole_Mint choice.
func (e *encoder) MintRoleMint(args MintRoleMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRole_Mint", args)
}

// Transfer encodes parameters for the Transfer choice.
func (e *encoder) Transfer(args Transfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Transfer", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
