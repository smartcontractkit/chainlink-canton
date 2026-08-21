package registry_app_v0

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_allocation_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_allocation_instruction_v1"
	splice_api_token_allocation_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_allocation_v1"
	splice_api_token_burn_mint_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_burn_mint_v1"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	splice_api_token_transfer_instruction_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	credential_v0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/credential_v0"
	registry_holding_v0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_holding_v0"
	registry_v0 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_v0"
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
	PackageName = "utility-registry-app-v0"
	PackageID   = "7a75ef6e69f69395a4e60919e228528bb8f3881150ccfde3f31bcc73864b18ab"
	SDKVersion  = "3.4.9"
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

// AllocationFactory is a Template type
type AllocationFactory struct {
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
	Operator  types.PARTY `json:"operator"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t AllocationFactory) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t AllocationFactory) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t AllocationFactory) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t AllocationFactory) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t AllocationFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactory to hex string (Canton MCMS format)
func (t AllocationFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactory from hex string (Canton MCMS format)
func (t *AllocationFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for AllocationFactory

// AllocationFactoryAllocateInternal exercises the AllocationFactory_AllocateInternal choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryAllocateInternal(contractID string, args AllocationFactoryAllocateInternal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_AllocateInternal",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryAllocateInternalWithPackageID exercises the AllocationFactory_AllocateInternal choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryAllocateInternalWithPackageID(contractID string, packageID string, args AllocationFactoryAllocateInternal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_AllocateInternal",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryInternalBurnMint exercises the AllocationFactory_InternalBurnMint choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryInternalBurnMint(contractID string, args AllocationFactoryInternalBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_InternalBurnMint",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryInternalBurnMintWithPackageID exercises the AllocationFactory_InternalBurnMint choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryInternalBurnMintWithPackageID(contractID string, packageID string, args AllocationFactoryInternalBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_InternalBurnMint",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryTransferInternal exercises the AllocationFactory_TransferInternal choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryTransferInternal(contractID string, args AllocationFactoryTransferInternal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_TransferInternal",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryTransferInternalWithPackageID exercises the AllocationFactory_TransferInternal choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryTransferInternalWithPackageID(contractID string, packageID string, args AllocationFactoryTransferInternal) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_TransferInternal",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryRequestBurn exercises the AllocationFactory_RequestBurn choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryRequestBurn(contractID string, args AllocationFactoryRequestBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_RequestBurn",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryRequestBurnWithPackageID exercises the AllocationFactory_RequestBurn choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryRequestBurnWithPackageID(contractID string, packageID string, args AllocationFactoryRequestBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_RequestBurn",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryOfferBurn exercises the AllocationFactory_OfferBurn choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryOfferBurn(contractID string, args AllocationFactoryOfferBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_OfferBurn",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryOfferBurnWithPackageID exercises the AllocationFactory_OfferBurn choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryOfferBurnWithPackageID(contractID string, packageID string, args AllocationFactoryOfferBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_OfferBurn",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryRequestMint exercises the AllocationFactory_RequestMint choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryRequestMint(contractID string, args AllocationFactoryRequestMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_RequestMint",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryRequestMintWithPackageID exercises the AllocationFactory_RequestMint choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryRequestMintWithPackageID(contractID string, packageID string, args AllocationFactoryRequestMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_RequestMint",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryOfferMint exercises the AllocationFactory_OfferMint choice on this AllocationFactory contract
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryOfferMint(contractID string, args AllocationFactoryOfferMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_OfferMint",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryOfferMintWithPackageID exercises the AllocationFactory_OfferMint choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryOfferMintWithPackageID(contractID string, packageID string, args AllocationFactoryOfferMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_OfferMint",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this AllocationFactory contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t AllocationFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this AllocationFactory contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) TransferFactoryTransfer(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransferWithPackageID exercises the TransferFactory_Transfer choice using the provided package ID instead of package name
func (t AllocationFactory) TransferFactoryTransferWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this AllocationFactory contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) TransferFactoryPublicFetch(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetchWithPackageID exercises the TransferFactory_PublicFetch choice using the provided package ID instead of package name
func (t AllocationFactory) TransferFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryAllocate exercises the AllocationFactory_Allocate choice on this AllocationFactory contract via the IAllocationFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryAllocate(contractID string, args splice_api_token_allocation_instruction_v1.AllocationFactoryAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_Allocate",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryAllocateWithPackageID exercises the AllocationFactory_Allocate choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryAllocateWithPackageID(contractID string, packageID string, args splice_api_token_allocation_instruction_v1.AllocationFactoryAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_Allocate",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryPublicFetch exercises the AllocationFactory_PublicFetch choice on this AllocationFactory contract via the IAllocationFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) AllocationFactoryPublicFetch(contractID string, args splice_api_token_allocation_instruction_v1.AllocationFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// AllocationFactoryPublicFetchWithPackageID exercises the AllocationFactory_PublicFetch choice using the provided package ID instead of package name
func (t AllocationFactory) AllocationFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_allocation_instruction_v1.AllocationFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "AllocationFactory"),
		ContractID: contractID,
		Choice:     "AllocationFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this AllocationFactory contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) BurnMintFactoryPublicFetch(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetchWithPackageID exercises the BurnMintFactory_PublicFetch choice using the provided package ID instead of package name
func (t AllocationFactory) BurnMintFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this AllocationFactory contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t AllocationFactory) BurnMintFactoryBurnMint(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMintWithPackageID exercises the BurnMintFactory_BurnMint choice using the provided package ID instead of package name
func (t AllocationFactory) BurnMintFactoryBurnMintWithPackageID(contractID string, packageID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.AllocationFactory", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for AllocationFactory

var _ splice_api_token_transfer_instruction_v1.ITransferFactory = (*AllocationFactory)(nil)

var _ splice_api_token_allocation_instruction_v1.IAllocationFactory = (*AllocationFactory)(nil)

var _ splice_api_token_burn_mint_v1.IBurnMintFactory = (*AllocationFactory)(nil)

// AllocationFactoryAllocateInternal is a Record type
type AllocationFactoryAllocateInternal struct {
	Payload splice_api_token_allocation_instruction_v1.AllocationFactoryAllocate `json:"payload"`
}

// ToMap converts AllocationFactoryAllocateInternal to a map for DAML arguments
func (t AllocationFactoryAllocateInternal) ToMap() map[string]any {
	m := make(map[string]any)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t AllocationFactoryAllocateInternal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryAllocateInternal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryAllocateInternal to hex string (Canton MCMS format)
func (t AllocationFactoryAllocateInternal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryAllocateInternal from hex string (Canton MCMS format)
func (t *AllocationFactoryAllocateInternal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryInternalBurnMint is a Record type
type AllocationFactoryInternalBurnMint struct {
	Payload splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint `json:"payload"`
}

// ToMap converts AllocationFactoryInternalBurnMint to a map for DAML arguments
func (t AllocationFactoryInternalBurnMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t AllocationFactoryInternalBurnMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryInternalBurnMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryInternalBurnMint to hex string (Canton MCMS format)
func (t AllocationFactoryInternalBurnMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryInternalBurnMint from hex string (Canton MCMS format)
func (t *AllocationFactoryInternalBurnMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryOfferBurn is a Record type
type AllocationFactoryOfferBurn struct {
	ExpectedAdmin types.PARTY                            `json:"expectedAdmin"`
	Burn          Burn                                   `json:"burn"`
	ExtraArgs     splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationFactoryOfferBurn to a map for DAML arguments
func (t AllocationFactoryOfferBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["burn"] = model.NestedToDAMLValue(t.Burn)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationFactoryOfferBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryOfferBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryOfferBurn to hex string (Canton MCMS format)
func (t AllocationFactoryOfferBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryOfferBurn from hex string (Canton MCMS format)
func (t *AllocationFactoryOfferBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryOfferBurnResult is a Record type
type AllocationFactoryOfferBurnResult struct {
	BurnOfferCid types.CONTRACT_ID                     `json:"burnOfferCid"`
	Meta         splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryOfferBurnResult to a map for DAML arguments
func (t AllocationFactoryOfferBurnResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["burnOfferCid"] = model.NestedToDAMLValue(t.BurnOfferCid)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationFactoryOfferBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryOfferBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryOfferBurnResult to hex string (Canton MCMS format)
func (t AllocationFactoryOfferBurnResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryOfferBurnResult from hex string (Canton MCMS format)
func (t *AllocationFactoryOfferBurnResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryOfferMint is a Record type
type AllocationFactoryOfferMint struct {
	ExpectedAdmin types.PARTY                            `json:"expectedAdmin"`
	Mint          Mint                                   `json:"mint"`
	ExtraArgs     splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationFactoryOfferMint to a map for DAML arguments
func (t AllocationFactoryOfferMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["mint"] = model.NestedToDAMLValue(t.Mint)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationFactoryOfferMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryOfferMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryOfferMint to hex string (Canton MCMS format)
func (t AllocationFactoryOfferMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryOfferMint from hex string (Canton MCMS format)
func (t *AllocationFactoryOfferMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryOfferMintResult is a Record type
type AllocationFactoryOfferMintResult struct {
	MintOfferCid types.CONTRACT_ID                     `json:"mintOfferCid"`
	Meta         splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryOfferMintResult to a map for DAML arguments
func (t AllocationFactoryOfferMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["mintOfferCid"] = model.NestedToDAMLValue(t.MintOfferCid)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationFactoryOfferMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryOfferMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryOfferMintResult to hex string (Canton MCMS format)
func (t AllocationFactoryOfferMintResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryOfferMintResult from hex string (Canton MCMS format)
func (t *AllocationFactoryOfferMintResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryRequestBurn is a Record type
type AllocationFactoryRequestBurn struct {
	ExpectedAdmin types.PARTY                            `json:"expectedAdmin"`
	Burn          Burn                                   `json:"burn"`
	HoldingCids   []types.CONTRACT_ID                    `json:"holdingCids"`
	ExtraArgs     splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationFactoryRequestBurn to a map for DAML arguments
func (t AllocationFactoryRequestBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["burn"] = model.NestedToDAMLValue(t.Burn)

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationFactoryRequestBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryRequestBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryRequestBurn to hex string (Canton MCMS format)
func (t AllocationFactoryRequestBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryRequestBurn from hex string (Canton MCMS format)
func (t *AllocationFactoryRequestBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryRequestBurnResult is a Record type
type AllocationFactoryRequestBurnResult struct {
	BurnRequestCid types.CONTRACT_ID                     `json:"burnRequestCid"`
	Remaining      *types.CONTRACT_ID                    `json:"remaining" hex:"optional"`
	Meta           splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryRequestBurnResult to a map for DAML arguments
func (t AllocationFactoryRequestBurnResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["burnRequestCid"] = model.NestedToDAMLValue(t.BurnRequestCid)

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

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationFactoryRequestBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryRequestBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryRequestBurnResult to hex string (Canton MCMS format)
func (t AllocationFactoryRequestBurnResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryRequestBurnResult from hex string (Canton MCMS format)
func (t *AllocationFactoryRequestBurnResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryRequestMint is a Record type
type AllocationFactoryRequestMint struct {
	ExpectedAdmin types.PARTY                            `json:"expectedAdmin"`
	Mint          Mint                                   `json:"mint"`
	ExtraArgs     splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationFactoryRequestMint to a map for DAML arguments
func (t AllocationFactoryRequestMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["mint"] = model.NestedToDAMLValue(t.Mint)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationFactoryRequestMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryRequestMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryRequestMint to hex string (Canton MCMS format)
func (t AllocationFactoryRequestMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryRequestMint from hex string (Canton MCMS format)
func (t *AllocationFactoryRequestMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryRequestMintResult is a Record type
type AllocationFactoryRequestMintResult struct {
	MintRequestCid types.CONTRACT_ID                     `json:"mintRequestCid"`
	Meta           splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryRequestMintResult to a map for DAML arguments
func (t AllocationFactoryRequestMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["mintRequestCid"] = model.NestedToDAMLValue(t.MintRequestCid)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationFactoryRequestMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryRequestMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryRequestMintResult to hex string (Canton MCMS format)
func (t AllocationFactoryRequestMintResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryRequestMintResult from hex string (Canton MCMS format)
func (t *AllocationFactoryRequestMintResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryTransferInternal is a Record type
type AllocationFactoryTransferInternal struct {
	Payload splice_api_token_transfer_instruction_v1.TransferFactoryTransfer `json:"payload"`
}

// ToMap converts AllocationFactoryTransferInternal to a map for DAML arguments
func (t AllocationFactoryTransferInternal) ToMap() map[string]any {
	m := make(map[string]any)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t AllocationFactoryTransferInternal) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryTransferInternal) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryTransferInternal to hex string (Canton MCMS format)
func (t AllocationFactoryTransferInternal) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryTransferInternal from hex string (Canton MCMS format)
func (t *AllocationFactoryTransferInternal) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Burn is a Record type
type Burn struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount        types.NUMERIC                            `json:"amount"`
	Holder        types.PARTY                              `json:"holder"`
	Reference     types.TEXT                               `json:"reference"`
	RequestedAt   types.TIMESTAMP                          `json:"requestedAt"`
	ExecuteBefore types.TIMESTAMP                          `json:"executeBefore"`
	Meta          splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// ToMap converts Burn to a map for DAML arguments
func (t Burn) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["requestedAt"] = t.RequestedAt

	m["executeBefore"] = t.ExecuteBefore

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
	Burn     Burn        `json:"burn"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t BurnOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t BurnOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t BurnOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

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
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

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

// BurnOfferAccept exercises the BurnOffer_Accept choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferAccept(contractID string, args BurnOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferAcceptWithPackageID exercises the BurnOffer_Accept choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferAcceptWithPackageID(contractID string, packageID string, args BurnOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferReject exercises the BurnOffer_Reject choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferReject(contractID string, args BurnOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferRejectWithPackageID exercises the BurnOffer_Reject choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferRejectWithPackageID(contractID string, packageID string, args BurnOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferCancel exercises the BurnOffer_Cancel choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) BurnOfferCancel(contractID string, args BurnOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// BurnOfferCancelWithPackageID exercises the BurnOffer_Cancel choice using the provided package ID instead of package name
func (t BurnOffer) BurnOfferCancelWithPackageID(contractID string, packageID string, args BurnOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "BurnOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BurnOffer contract
// This method uses the package name in the template ID
func (t BurnOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t BurnOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// BurnOfferAccept is a Record type
type BurnOfferAccept struct {
	HoldingCids []types.CONTRACT_ID                    `json:"holdingCids"`
	ExtraArgs   splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts BurnOfferAccept to a map for DAML arguments
func (t BurnOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

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
	ExecutedBurnCid *types.CONTRACT_ID                    `json:"executedBurnCid" hex:"optional"`
	Remaining       *types.CONTRACT_ID                    `json:"remaining" hex:"optional"`
	Meta            splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts BurnOfferAcceptResult to a map for DAML arguments
func (t BurnOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.ExecutedBurnCid != nil {
		m["executedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutedBurnCid),
		}
	} else {
		m["executedBurnCid"] = map[string]any{
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

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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

// BurnOfferReject is a Record type
type BurnOfferReject struct {
	Reason    types.TEXT                              `json:"reason"`
	ExtraArgs *splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs" hex:"optional"`
}

// ToMap converts BurnOfferReject to a map for DAML arguments
func (t BurnOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	if t.ExtraArgs != nil {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExtraArgs),
		}
	} else {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	RejectedBurnCid *types.CONTRACT_ID `json:"rejectedBurnCid" hex:"optional"`
}

// ToMap converts BurnOfferRejectResult to a map for DAML arguments
func (t BurnOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RejectedBurnCid != nil {
		m["rejectedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RejectedBurnCid),
		}
	} else {
		m["rejectedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	Operator         types.PARTY       `json:"operator"`
	Provider         types.PARTY       `json:"provider"`
	Burn             Burn              `json:"burn"`
	LockedHoldingCid types.CONTRACT_ID `json:"lockedHoldingCid"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t BurnRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t BurnRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t BurnRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedHoldingCid"] = model.NestedToDAMLValue(t.LockedHoldingCid)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t BurnRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lockedHoldingCid"] = model.NestedToDAMLValue(t.LockedHoldingCid)

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

// BurnRequestAccept exercises the BurnRequest_Accept choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestAccept(contractID string, args BurnRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestAcceptWithPackageID exercises the BurnRequest_Accept choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestAcceptWithPackageID(contractID string, packageID string, args BurnRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestReject exercises the BurnRequest_Reject choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestReject(contractID string, args BurnRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestRejectWithPackageID exercises the BurnRequest_Reject choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestRejectWithPackageID(contractID string, packageID string, args BurnRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestCancel exercises the BurnRequest_Cancel choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) BurnRequestCancel(contractID string, args BurnRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// BurnRequestCancelWithPackageID exercises the BurnRequest_Cancel choice using the provided package ID instead of package name
func (t BurnRequest) BurnRequestCancelWithPackageID(contractID string, packageID string, args BurnRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "BurnRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this BurnRequest contract
// This method uses the package name in the template ID
func (t BurnRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t BurnRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "BurnRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// BurnRequestAccept is a Record type
type BurnRequestAccept struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts BurnRequestAccept to a map for DAML arguments
func (t BurnRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

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
	ExecutedBurnCid *types.CONTRACT_ID                    `json:"executedBurnCid" hex:"optional"`
	Meta            splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts BurnRequestAcceptResult to a map for DAML arguments
func (t BurnRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.ExecutedBurnCid != nil {
		m["executedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutedBurnCid),
		}
	} else {
		m["executedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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
	HoldingCid types.CONTRACT_ID `json:"holdingCid"`
}

// ToMap converts BurnRequestCancelResult to a map for DAML arguments
func (t BurnRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

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

// BurnRequestReject is a Record type
type BurnRequestReject struct {
	Reason    types.TEXT                              `json:"reason"`
	ExtraArgs *splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs" hex:"optional"`
}

// ToMap converts BurnRequestReject to a map for DAML arguments
func (t BurnRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	if t.ExtraArgs != nil {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExtraArgs),
		}
	} else {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	RejectedBurnCid *types.CONTRACT_ID `json:"rejectedBurnCid" hex:"optional"`
	HoldingCid      types.CONTRACT_ID  `json:"holdingCid"`
}

// ToMap converts BurnRequestRejectResult to a map for DAML arguments
func (t BurnRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RejectedBurnCid != nil {
		m["rejectedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RejectedBurnCid),
		}
	} else {
		m["rejectedBurnCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

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

// EnforcementService is a Template type
type EnforcementService struct {
	Operator  types.PARTY `json:"operator"`
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
	Holder    types.PARTY `json:"holder"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t EnforcementService) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t EnforcementService) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t EnforcementService) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t EnforcementService) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t EnforcementService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementService to hex string (Canton MCMS format)
func (t EnforcementService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementService from hex string (Canton MCMS format)
func (t *EnforcementService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for EnforcementService

// EnforcementServiceTerminate exercises the EnforcementService_Terminate choice on this EnforcementService contract
// This method uses the package name in the template ID
func (t EnforcementService) EnforcementServiceTerminate(contractID string, args EnforcementServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceTerminateWithPackageID exercises the EnforcementService_Terminate choice using the provided package ID instead of package name
func (t EnforcementService) EnforcementServiceTerminateWithPackageID(contractID string, packageID string, args EnforcementServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceAcceptForceTransferRequest exercises the EnforcementService_AcceptForceTransferRequest choice on this EnforcementService contract
// This method uses the package name in the template ID
func (t EnforcementService) EnforcementServiceAcceptForceTransferRequest(contractID string, args EnforcementServiceAcceptForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_AcceptForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceAcceptForceTransferRequestWithPackageID exercises the EnforcementService_AcceptForceTransferRequest choice using the provided package ID instead of package name
func (t EnforcementService) EnforcementServiceAcceptForceTransferRequestWithPackageID(contractID string, packageID string, args EnforcementServiceAcceptForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_AcceptForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization exercises the EnforcementService_AcceptForceTransferRequestWithSenderAuthorization choice on this EnforcementService contract
// This method uses the package name in the template ID
func (t EnforcementService) EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization(contractID string, args EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_AcceptForceTransferRequestWithSenderAuthorization",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceAcceptForceTransferRequestWithSenderAuthorizationWithPackageID exercises the EnforcementService_AcceptForceTransferRequestWithSenderAuthorization choice using the provided package ID instead of package name
func (t EnforcementService) EnforcementServiceAcceptForceTransferRequestWithSenderAuthorizationWithPackageID(contractID string, packageID string, args EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "EnforcementService_AcceptForceTransferRequestWithSenderAuthorization",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this EnforcementService contract
// This method uses the package name in the template ID
func (t EnforcementService) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t EnforcementService) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// EnforcementServiceRequest is a Template type
type EnforcementServiceRequest struct {
	Operator  types.PARTY `json:"operator"`
	Provider  types.PARTY `json:"provider"`
	Registrar types.PARTY `json:"registrar"`
	Holder    types.PARTY `json:"holder"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t EnforcementServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t EnforcementServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t EnforcementServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t EnforcementServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t EnforcementServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequest to hex string (Canton MCMS format)
func (t EnforcementServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequest from hex string (Canton MCMS format)
func (t *EnforcementServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for EnforcementServiceRequest

// EnforcementServiceRequestAccept exercises the EnforcementServiceRequest_Accept choice on this EnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t EnforcementServiceRequest) EnforcementServiceRequestAccept(contractID string, args EnforcementServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceRequestAcceptWithPackageID exercises the EnforcementServiceRequest_Accept choice using the provided package ID instead of package name
func (t EnforcementServiceRequest) EnforcementServiceRequestAcceptWithPackageID(contractID string, packageID string, args EnforcementServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceRequestReject exercises the EnforcementServiceRequest_Reject choice on this EnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t EnforcementServiceRequest) EnforcementServiceRequestReject(contractID string, args EnforcementServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceRequestRejectWithPackageID exercises the EnforcementServiceRequest_Reject choice using the provided package ID instead of package name
func (t EnforcementServiceRequest) EnforcementServiceRequestRejectWithPackageID(contractID string, packageID string, args EnforcementServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceRequestCancel exercises the EnforcementServiceRequest_Cancel choice on this EnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t EnforcementServiceRequest) EnforcementServiceRequestCancel(contractID string, args EnforcementServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// EnforcementServiceRequestCancelWithPackageID exercises the EnforcementServiceRequest_Cancel choice using the provided package ID instead of package name
func (t EnforcementServiceRequest) EnforcementServiceRequestCancelWithPackageID(contractID string, packageID string, args EnforcementServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "EnforcementServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this EnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t EnforcementServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t EnforcementServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "EnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// EnforcementServiceRequestAccept is a Record type
type EnforcementServiceRequestAccept struct {
	RegistrarConfigurationCid types.CONTRACT_ID   `json:"registrarConfigurationCid"`
	CredentialCids            []types.CONTRACT_ID `json:"credentialCids"`
}

// ToMap converts EnforcementServiceRequestAccept to a map for DAML arguments
func (t EnforcementServiceRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarConfigurationCid"] = model.NestedToDAMLValue(t.RegistrarConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t EnforcementServiceRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestAccept to hex string (Canton MCMS format)
func (t EnforcementServiceRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestAccept from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceRequestAcceptResult is a Record type
type EnforcementServiceRequestAcceptResult struct {
	EnforcementServiceCid types.CONTRACT_ID `json:"enforcementServiceCid"`
}

// ToMap converts EnforcementServiceRequestAcceptResult to a map for DAML arguments
func (t EnforcementServiceRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["enforcementServiceCid"] = model.NestedToDAMLValue(t.EnforcementServiceCid)

	return m
}

func (t EnforcementServiceRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestAcceptResult to hex string (Canton MCMS format)
func (t EnforcementServiceRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestAcceptResult from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceRequestCancel is a Record type
type EnforcementServiceRequestCancel struct {
}

// ToMap converts EnforcementServiceRequestCancel to a map for DAML arguments
func (t EnforcementServiceRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t EnforcementServiceRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestCancel to hex string (Canton MCMS format)
func (t EnforcementServiceRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestCancel from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceRequestCancelResult is a Record type
type EnforcementServiceRequestCancelResult struct {
}

// ToMap converts EnforcementServiceRequestCancelResult to a map for DAML arguments
func (t EnforcementServiceRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t EnforcementServiceRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestCancelResult to hex string (Canton MCMS format)
func (t EnforcementServiceRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestCancelResult from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceRequestReject is a Record type
type EnforcementServiceRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts EnforcementServiceRequestReject to a map for DAML arguments
func (t EnforcementServiceRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t EnforcementServiceRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestReject to hex string (Canton MCMS format)
func (t EnforcementServiceRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestReject from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceRequestRejectResult is a Record type
type EnforcementServiceRequestRejectResult struct {
	RejectedEnforcementServiceRequestCid types.CONTRACT_ID `json:"rejectedEnforcementServiceRequestCid"`
}

// ToMap converts EnforcementServiceRequestRejectResult to a map for DAML arguments
func (t EnforcementServiceRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedEnforcementServiceRequestCid"] = model.NestedToDAMLValue(t.RejectedEnforcementServiceRequestCid)

	return m
}

func (t EnforcementServiceRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceRequestRejectResult to hex string (Canton MCMS format)
func (t EnforcementServiceRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceRequestRejectResult from hex string (Canton MCMS format)
func (t *EnforcementServiceRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceAcceptForceTransferRequest is a Record type
type EnforcementServiceAcceptForceTransferRequest struct {
	Cid                           types.CONTRACT_ID                      `json:"cid"`
	Payload                       registry_v0.ForceTransferRequestAccept `json:"payload"`
	ReceiverEnforcementServiceCid types.CONTRACT_ID                      `json:"receiverEnforcementServiceCid"`
}

// ToMap converts EnforcementServiceAcceptForceTransferRequest to a map for DAML arguments
func (t EnforcementServiceAcceptForceTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	m["receiverEnforcementServiceCid"] = model.NestedToDAMLValue(t.ReceiverEnforcementServiceCid)

	return m
}

func (t EnforcementServiceAcceptForceTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceAcceptForceTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceAcceptForceTransferRequest to hex string (Canton MCMS format)
func (t EnforcementServiceAcceptForceTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceAcceptForceTransferRequest from hex string (Canton MCMS format)
func (t *EnforcementServiceAcceptForceTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization is a Record type
type EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization struct {
	Cid     types.CONTRACT_ID                      `json:"cid"`
	Payload registry_v0.ForceTransferRequestAccept `json:"payload"`
	Sender  types.PARTY                            `json:"sender"`
}

// ToMap converts EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization to a map for DAML arguments
func (t EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	m["sender"] = t.Sender.ToMap()

	return m
}

func (t EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization to hex string (Canton MCMS format)
func (t EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization from hex string (Canton MCMS format)
func (t *EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceTerminate is a Record type
type EnforcementServiceTerminate struct {
}

// ToMap converts EnforcementServiceTerminate to a map for DAML arguments
func (t EnforcementServiceTerminate) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t EnforcementServiceTerminate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceTerminate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceTerminate to hex string (Canton MCMS format)
func (t EnforcementServiceTerminate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceTerminate from hex string (Canton MCMS format)
func (t *EnforcementServiceTerminate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// EnforcementServiceTerminateResult is a Record type
type EnforcementServiceTerminateResult struct {
}

// ToMap converts EnforcementServiceTerminateResult to a map for DAML arguments
func (t EnforcementServiceTerminateResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t EnforcementServiceTerminateResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *EnforcementServiceTerminateResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes EnforcementServiceTerminateResult to hex string (Canton MCMS format)
func (t EnforcementServiceTerminateResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes EnforcementServiceTerminateResult from hex string (Canton MCMS format)
func (t *EnforcementServiceTerminateResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecutedBurn is a Template type
type ExecutedBurn struct {
	Operator           types.PARTY `json:"operator"`
	Provider           types.PARTY `json:"provider"`
	Burn               Burn        `json:"burn"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

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
func (t ExecutedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

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

// ExecutedBurnDelete exercises the ExecutedBurn_Delete choice on this ExecutedBurn contract
// This method uses the package name in the template ID
func (t ExecutedBurn) ExecutedBurnDelete(contractID string, args ExecutedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedBurnDeleteWithPackageID exercises the ExecutedBurn_Delete choice using the provided package ID instead of package name
func (t ExecutedBurn) ExecutedBurnDeleteWithPackageID(contractID string, packageID string, args ExecutedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "ExecutedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedBurn contract
// This method uses the package name in the template ID
func (t ExecutedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "ExecutedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedBurnDelete is a Record type
type ExecutedBurnDelete struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedBurnDelete to a map for DAML arguments
func (t ExecutedBurnDelete) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

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

// ExecutedMint is a Template type
type ExecutedMint struct {
	Operator           types.PARTY `json:"operator"`
	Provider           types.PARTY `json:"provider"`
	Mint               Mint        `json:"mint"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ExecutedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ExecutedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ExecutedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

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
func (t ExecutedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

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

// ExecutedMintDelete exercises the ExecutedMint_Delete choice on this ExecutedMint contract
// This method uses the package name in the template ID
func (t ExecutedMint) ExecutedMintDelete(contractID string, args ExecutedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// ExecutedMintDeleteWithPackageID exercises the ExecutedMint_Delete choice using the provided package ID instead of package name
func (t ExecutedMint) ExecutedMintDeleteWithPackageID(contractID string, packageID string, args ExecutedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "ExecutedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ExecutedMint contract
// This method uses the package name in the template ID
func (t ExecutedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ExecutedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "ExecutedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ExecutedMintDelete is a Record type
type ExecutedMintDelete struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ExecutedMintDelete to a map for DAML arguments
func (t ExecutedMintDelete) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

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

// HolderService is a Template type
type HolderService struct {
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
	Holder   types.PARTY `json:"holder"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t HolderService) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t HolderService) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t HolderService) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t HolderService) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t HolderService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderService to hex string (Canton MCMS format)
func (t HolderService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderService from hex string (Canton MCMS format)
func (t *HolderService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for HolderService

// HolderServiceClean exercises the HolderService_Clean choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceClean(contractID string, args HolderServiceClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_Clean",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCleanWithPackageID exercises the HolderService_Clean choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCleanWithPackageID(contractID string, packageID string, args HolderServiceClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_Clean",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceTerminate exercises the HolderService_Terminate choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceTerminate(contractID string, args HolderServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceTerminateWithPackageID exercises the HolderService_Terminate choice using the provided package ID instead of package name
func (t HolderService) HolderServiceTerminateWithPackageID(contractID string, packageID string, args HolderServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestEnforcementService exercises the HolderService_RequestEnforcementService choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestEnforcementService(contractID string, args HolderServiceRequestEnforcementService) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestEnforcementService",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestEnforcementServiceWithPackageID exercises the HolderService_RequestEnforcementService choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestEnforcementServiceWithPackageID(contractID string, packageID string, args HolderServiceRequestEnforcementService) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestEnforcementService",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelEnforcementServiceRequest exercises the HolderService_CancelEnforcementServiceRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelEnforcementServiceRequest(contractID string, args HolderServiceCancelEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelEnforcementServiceRequestWithPackageID exercises the HolderService_CancelEnforcementServiceRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelEnforcementServiceRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestForceTransfer exercises the HolderService_RequestForceTransfer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestForceTransfer(contractID string, args HolderServiceRequestForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestForceTransferWithPackageID exercises the HolderService_RequestForceTransfer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestForceTransferWithPackageID(contractID string, packageID string, args HolderServiceRequestForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelForceTransferRequest exercises the HolderService_CancelForceTransferRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelForceTransferRequest(contractID string, args HolderServiceCancelForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelForceTransferRequestWithPackageID exercises the HolderService_CancelForceTransferRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelForceTransferRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestMint exercises the HolderService_RequestMint choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestMint(contractID string, args HolderServiceRequestMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestMint",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestMintWithPackageID exercises the HolderService_RequestMint choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestMintWithPackageID(contractID string, packageID string, args HolderServiceRequestMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestMint",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelMintRequest exercises the HolderService_CancelMintRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelMintRequest(contractID string, args HolderServiceCancelMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelMintRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelMintRequestWithPackageID exercises the HolderService_CancelMintRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelMintRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelMintRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptMintOffer exercises the HolderService_AcceptMintOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptMintOffer(contractID string, args HolderServiceAcceptMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptMintOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptMintOfferWithPackageID exercises the HolderService_AcceptMintOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptMintOfferWithPackageID(contractID string, packageID string, args HolderServiceAcceptMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptMintOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectMintOffer exercises the HolderService_RejectMintOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectMintOffer(contractID string, args HolderServiceRejectMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectMintOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectMintOfferWithPackageID exercises the HolderService_RejectMintOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectMintOfferWithPackageID(contractID string, packageID string, args HolderServiceRejectMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectMintOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestBurn exercises the HolderService_RequestBurn choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestBurn(contractID string, args HolderServiceRequestBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestBurn",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestBurnWithPackageID exercises the HolderService_RequestBurn choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestBurnWithPackageID(contractID string, packageID string, args HolderServiceRequestBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestBurn",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelBurnRequest exercises the HolderService_CancelBurnRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelBurnRequest(contractID string, args HolderServiceCancelBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelBurnRequestWithPackageID exercises the HolderService_CancelBurnRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelBurnRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptBurnOffer exercises the HolderService_AcceptBurnOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptBurnOffer(contractID string, args HolderServiceAcceptBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptBurnOfferWithPackageID exercises the HolderService_AcceptBurnOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptBurnOfferWithPackageID(contractID string, packageID string, args HolderServiceAcceptBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectBurnOffer exercises the HolderService_RejectBurnOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectBurnOffer(contractID string, args HolderServiceRejectBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectBurnOfferWithPackageID exercises the HolderService_RejectBurnOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectBurnOfferWithPackageID(contractID string, packageID string, args HolderServiceRejectBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferTransfer exercises the HolderService_OfferTransfer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceOfferTransfer(contractID string, args HolderServiceOfferTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferTransferWithPackageID exercises the HolderService_OfferTransfer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceOfferTransferWithPackageID(contractID string, packageID string, args HolderServiceOfferTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelTransferOffer exercises the HolderService_CancelTransferOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelTransferOffer(contractID string, args HolderServiceCancelTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelTransferOfferWithPackageID exercises the HolderService_CancelTransferOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelTransferOfferWithPackageID(contractID string, packageID string, args HolderServiceCancelTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptTransferOffer exercises the HolderService_AcceptTransferOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptTransferOffer(contractID string, args HolderServiceAcceptTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptTransferOfferWithPackageID exercises the HolderService_AcceptTransferOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptTransferOfferWithPackageID(contractID string, packageID string, args HolderServiceAcceptTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectTransferOffer exercises the HolderService_RejectTransferOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectTransferOffer(contractID string, args HolderServiceRejectTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectTransferOfferWithPackageID exercises the HolderService_RejectTransferOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectTransferOfferWithPackageID(contractID string, packageID string, args HolderServiceRejectTransferOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectTransferOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestTransfer exercises the HolderService_RequestTransfer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestTransfer(contractID string, args HolderServiceRequestTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestTransferWithPackageID exercises the HolderService_RequestTransfer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestTransferWithPackageID(contractID string, packageID string, args HolderServiceRequestTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestTransfer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelTransferRequest exercises the HolderService_CancelTransferRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelTransferRequest(contractID string, args HolderServiceCancelTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelTransferRequestWithPackageID exercises the HolderService_CancelTransferRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelTransferRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptTransferRequest exercises the HolderService_AcceptTransferRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptTransferRequest(contractID string, args HolderServiceAcceptTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptTransferRequestWithPackageID exercises the HolderService_AcceptTransferRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptTransferRequestWithPackageID(contractID string, packageID string, args HolderServiceAcceptTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectTransferRequest exercises the HolderService_RejectTransferRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectTransferRequest(contractID string, args HolderServiceRejectTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectTransferRequestWithPackageID exercises the HolderService_RejectTransferRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectTransferRequestWithPackageID(contractID string, packageID string, args HolderServiceRejectTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferLock exercises the HolderService_OfferLock choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceOfferLock(contractID string, args HolderServiceOfferLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferLock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferLockWithPackageID exercises the HolderService_OfferLock choice using the provided package ID instead of package name
func (t HolderService) HolderServiceOfferLockWithPackageID(contractID string, packageID string, args HolderServiceOfferLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferLock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelLockOffer exercises the HolderService_CancelLockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelLockOffer(contractID string, args HolderServiceCancelLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelLockOfferWithPackageID exercises the HolderService_CancelLockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelLockOfferWithPackageID(contractID string, packageID string, args HolderServiceCancelLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptLockOffer exercises the HolderService_AcceptLockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptLockOffer(contractID string, args HolderServiceAcceptLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptLockOfferWithPackageID exercises the HolderService_AcceptLockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptLockOfferWithPackageID(contractID string, packageID string, args HolderServiceAcceptLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectLockOffer exercises the HolderService_RejectLockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectLockOffer(contractID string, args HolderServiceRejectLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectLockOfferWithPackageID exercises the HolderService_RejectLockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectLockOfferWithPackageID(contractID string, packageID string, args HolderServiceRejectLockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectLockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestLock exercises the HolderService_RequestLock choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestLock(contractID string, args HolderServiceRequestLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestLock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestLockWithPackageID exercises the HolderService_RequestLock choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestLockWithPackageID(contractID string, packageID string, args HolderServiceRequestLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestLock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelLockRequest exercises the HolderService_CancelLockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelLockRequest(contractID string, args HolderServiceCancelLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelLockRequestWithPackageID exercises the HolderService_CancelLockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelLockRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptLockRequest exercises the HolderService_AcceptLockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptLockRequest(contractID string, args HolderServiceAcceptLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptLockRequestWithPackageID exercises the HolderService_AcceptLockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptLockRequestWithPackageID(contractID string, packageID string, args HolderServiceAcceptLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectLockRequest exercises the HolderService_RejectLockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectLockRequest(contractID string, args HolderServiceRejectLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectLockRequestWithPackageID exercises the HolderService_RejectLockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectLockRequestWithPackageID(contractID string, packageID string, args HolderServiceRejectLockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectLockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferUnlock exercises the HolderService_OfferUnlock choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceOfferUnlock(contractID string, args HolderServiceOfferUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferUnlock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceOfferUnlockWithPackageID exercises the HolderService_OfferUnlock choice using the provided package ID instead of package name
func (t HolderService) HolderServiceOfferUnlockWithPackageID(contractID string, packageID string, args HolderServiceOfferUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_OfferUnlock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelUnlockOffer exercises the HolderService_CancelUnlockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelUnlockOffer(contractID string, args HolderServiceCancelUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelUnlockOfferWithPackageID exercises the HolderService_CancelUnlockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelUnlockOfferWithPackageID(contractID string, packageID string, args HolderServiceCancelUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptUnlockOffer exercises the HolderService_AcceptUnlockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptUnlockOffer(contractID string, args HolderServiceAcceptUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptUnlockOfferWithPackageID exercises the HolderService_AcceptUnlockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptUnlockOfferWithPackageID(contractID string, packageID string, args HolderServiceAcceptUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectUnlockOffer exercises the HolderService_RejectUnlockOffer choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectUnlockOffer(contractID string, args HolderServiceRejectUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectUnlockOfferWithPackageID exercises the HolderService_RejectUnlockOffer choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectUnlockOfferWithPackageID(contractID string, packageID string, args HolderServiceRejectUnlockOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectUnlockOffer",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestUnlock exercises the HolderService_RequestUnlock choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRequestUnlock(contractID string, args HolderServiceRequestUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestUnlock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestUnlockWithPackageID exercises the HolderService_RequestUnlock choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRequestUnlockWithPackageID(contractID string, packageID string, args HolderServiceRequestUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RequestUnlock",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelUnlockRequest exercises the HolderService_CancelUnlockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCancelUnlockRequest(contractID string, args HolderServiceCancelUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCancelUnlockRequestWithPackageID exercises the HolderService_CancelUnlockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCancelUnlockRequestWithPackageID(contractID string, packageID string, args HolderServiceCancelUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CancelUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptUnlockRequest exercises the HolderService_AcceptUnlockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceAcceptUnlockRequest(contractID string, args HolderServiceAcceptUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceAcceptUnlockRequestWithPackageID exercises the HolderService_AcceptUnlockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceAcceptUnlockRequestWithPackageID(contractID string, packageID string, args HolderServiceAcceptUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_AcceptUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectUnlockRequest exercises the HolderService_RejectUnlockRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectUnlockRequest(contractID string, args HolderServiceRejectUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectUnlockRequestWithPackageID exercises the HolderService_RejectUnlockRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectUnlockRequestWithPackageID(contractID string, packageID string, args HolderServiceRejectUnlockRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectUnlockRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCreateAllocation exercises the HolderService_CreateAllocation choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceCreateAllocation(contractID string, args HolderServiceCreateAllocation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CreateAllocation",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceCreateAllocationWithPackageID exercises the HolderService_CreateAllocation choice using the provided package ID instead of package name
func (t HolderService) HolderServiceCreateAllocationWithPackageID(contractID string, packageID string, args HolderServiceCreateAllocation) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_CreateAllocation",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectAllocationRequest exercises the HolderService_RejectAllocationRequest choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) HolderServiceRejectAllocationRequest(contractID string, args HolderServiceRejectAllocationRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectAllocationRequest",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRejectAllocationRequestWithPackageID exercises the HolderService_RejectAllocationRequest choice using the provided package ID instead of package name
func (t HolderService) HolderServiceRejectAllocationRequestWithPackageID(contractID string, packageID string, args HolderServiceRejectAllocationRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "HolderService_RejectAllocationRequest",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this HolderService contract
// This method uses the package name in the template ID
func (t HolderService) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t HolderService) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// HolderServiceRequest is a Template type
type HolderServiceRequest struct {
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
	Holder   types.PARTY `json:"holder"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t HolderServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t HolderServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t HolderServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t HolderServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t HolderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequest to hex string (Canton MCMS format)
func (t HolderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequest from hex string (Canton MCMS format)
func (t *HolderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for HolderServiceRequest

// HolderServiceRequestClean exercises the HolderServiceRequest_Clean choice on this HolderServiceRequest contract
// This method uses the package name in the template ID
func (t HolderServiceRequest) HolderServiceRequestClean(contractID string, args HolderServiceRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestCleanWithPackageID exercises the HolderServiceRequest_Clean choice using the provided package ID instead of package name
func (t HolderServiceRequest) HolderServiceRequestCleanWithPackageID(contractID string, packageID string, args HolderServiceRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestAccept exercises the HolderServiceRequest_Accept choice on this HolderServiceRequest contract
// This method uses the package name in the template ID
func (t HolderServiceRequest) HolderServiceRequestAccept(contractID string, args HolderServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestAcceptWithPackageID exercises the HolderServiceRequest_Accept choice using the provided package ID instead of package name
func (t HolderServiceRequest) HolderServiceRequestAcceptWithPackageID(contractID string, packageID string, args HolderServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestReject exercises the HolderServiceRequest_Reject choice on this HolderServiceRequest contract
// This method uses the package name in the template ID
func (t HolderServiceRequest) HolderServiceRequestReject(contractID string, args HolderServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestRejectWithPackageID exercises the HolderServiceRequest_Reject choice using the provided package ID instead of package name
func (t HolderServiceRequest) HolderServiceRequestRejectWithPackageID(contractID string, packageID string, args HolderServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestCancel exercises the HolderServiceRequest_Cancel choice on this HolderServiceRequest contract
// This method uses the package name in the template ID
func (t HolderServiceRequest) HolderServiceRequestCancel(contractID string, args HolderServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// HolderServiceRequestCancelWithPackageID exercises the HolderServiceRequest_Cancel choice using the provided package ID instead of package name
func (t HolderServiceRequest) HolderServiceRequestCancelWithPackageID(contractID string, packageID string, args HolderServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "HolderServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this HolderServiceRequest contract
// This method uses the package name in the template ID
func (t HolderServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t HolderServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "HolderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// HolderServiceRequestAccept is a Record type
type HolderServiceRequestAccept struct {
	ProviderConfigurationCid types.CONTRACT_ID   `json:"providerConfigurationCid"`
	CredentialCids           []types.CONTRACT_ID `json:"credentialCids"`
}

// ToMap converts HolderServiceRequestAccept to a map for DAML arguments
func (t HolderServiceRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerConfigurationCid"] = model.NestedToDAMLValue(t.ProviderConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t HolderServiceRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestAccept to hex string (Canton MCMS format)
func (t HolderServiceRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestAccept from hex string (Canton MCMS format)
func (t *HolderServiceRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestAcceptResult is a Record type
type HolderServiceRequestAcceptResult struct {
	HolderServiceCid types.CONTRACT_ID `json:"holderServiceCid"`
}

// ToMap converts HolderServiceRequestAcceptResult to a map for DAML arguments
func (t HolderServiceRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holderServiceCid"] = model.NestedToDAMLValue(t.HolderServiceCid)

	return m
}

func (t HolderServiceRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestAcceptResult to hex string (Canton MCMS format)
func (t HolderServiceRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestAcceptResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestCancel is a Record type
type HolderServiceRequestCancel struct {
}

// ToMap converts HolderServiceRequestCancel to a map for DAML arguments
func (t HolderServiceRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t HolderServiceRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestCancel to hex string (Canton MCMS format)
func (t HolderServiceRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestCancel from hex string (Canton MCMS format)
func (t *HolderServiceRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestCancelResult is a Record type
type HolderServiceRequestCancelResult struct {
}

// ToMap converts HolderServiceRequestCancelResult to a map for DAML arguments
func (t HolderServiceRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t HolderServiceRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestCancelResult to hex string (Canton MCMS format)
func (t HolderServiceRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestCancelResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestClean is a Record type
type HolderServiceRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts HolderServiceRequestClean to a map for DAML arguments
func (t HolderServiceRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t HolderServiceRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestClean to hex string (Canton MCMS format)
func (t HolderServiceRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestClean from hex string (Canton MCMS format)
func (t *HolderServiceRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestReject is a Record type
type HolderServiceRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts HolderServiceRequestReject to a map for DAML arguments
func (t HolderServiceRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t HolderServiceRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestReject to hex string (Canton MCMS format)
func (t HolderServiceRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestReject from hex string (Canton MCMS format)
func (t *HolderServiceRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestRejectResult is a Record type
type HolderServiceRequestRejectResult struct {
	RejectedHolderServiceRequestCid types.CONTRACT_ID `json:"rejectedHolderServiceRequestCid"`
}

// ToMap converts HolderServiceRequestRejectResult to a map for DAML arguments
func (t HolderServiceRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedHolderServiceRequestCid"] = model.NestedToDAMLValue(t.RejectedHolderServiceRequestCid)

	return m
}

func (t HolderServiceRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestRejectResult to hex string (Canton MCMS format)
func (t HolderServiceRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestRejectResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptBurnOffer is a Record type
type HolderServiceAcceptBurnOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.BurnOfferAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptBurnOffer to a map for DAML arguments
func (t HolderServiceAcceptBurnOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptBurnOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptBurnOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptBurnOffer to hex string (Canton MCMS format)
func (t HolderServiceAcceptBurnOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptBurnOffer from hex string (Canton MCMS format)
func (t *HolderServiceAcceptBurnOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptLockOffer is a Record type
type HolderServiceAcceptLockOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.LockOfferAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptLockOffer to a map for DAML arguments
func (t HolderServiceAcceptLockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptLockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptLockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptLockOffer to hex string (Canton MCMS format)
func (t HolderServiceAcceptLockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptLockOffer from hex string (Canton MCMS format)
func (t *HolderServiceAcceptLockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptLockRequest is a Record type
type HolderServiceAcceptLockRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.LockRequestAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptLockRequest to a map for DAML arguments
func (t HolderServiceAcceptLockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptLockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptLockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptLockRequest to hex string (Canton MCMS format)
func (t HolderServiceAcceptLockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptLockRequest from hex string (Canton MCMS format)
func (t *HolderServiceAcceptLockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptMintOffer is a Record type
type HolderServiceAcceptMintOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.MintOfferAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptMintOffer to a map for DAML arguments
func (t HolderServiceAcceptMintOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptMintOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptMintOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptMintOffer to hex string (Canton MCMS format)
func (t HolderServiceAcceptMintOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptMintOffer from hex string (Canton MCMS format)
func (t *HolderServiceAcceptMintOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptTransferOffer is a Record type
type HolderServiceAcceptTransferOffer struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.TransferOfferAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptTransferOffer to a map for DAML arguments
func (t HolderServiceAcceptTransferOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptTransferOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptTransferOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptTransferOffer to hex string (Canton MCMS format)
func (t HolderServiceAcceptTransferOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptTransferOffer from hex string (Canton MCMS format)
func (t *HolderServiceAcceptTransferOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptTransferOfferResult is a Record type
type HolderServiceAcceptTransferOfferResult struct {
	AcceptedTransferCid types.CONTRACT_ID `json:"acceptedTransferCid"`
}

// ToMap converts HolderServiceAcceptTransferOfferResult to a map for DAML arguments
func (t HolderServiceAcceptTransferOfferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["acceptedTransferCid"] = model.NestedToDAMLValue(t.AcceptedTransferCid)

	return m
}

func (t HolderServiceAcceptTransferOfferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptTransferOfferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptTransferOfferResult to hex string (Canton MCMS format)
func (t HolderServiceAcceptTransferOfferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptTransferOfferResult from hex string (Canton MCMS format)
func (t *HolderServiceAcceptTransferOfferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptTransferRequest is a Record type
type HolderServiceAcceptTransferRequest struct {
	Cid     types.CONTRACT_ID                 `json:"cid"`
	Payload registry_v0.TransferRequestAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptTransferRequest to a map for DAML arguments
func (t HolderServiceAcceptTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptTransferRequest to hex string (Canton MCMS format)
func (t HolderServiceAcceptTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptTransferRequest from hex string (Canton MCMS format)
func (t *HolderServiceAcceptTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptUnlockOffer is a Record type
type HolderServiceAcceptUnlockOffer struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.UnlockOfferAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptUnlockOffer to a map for DAML arguments
func (t HolderServiceAcceptUnlockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptUnlockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptUnlockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptUnlockOffer to hex string (Canton MCMS format)
func (t HolderServiceAcceptUnlockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptUnlockOffer from hex string (Canton MCMS format)
func (t *HolderServiceAcceptUnlockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceAcceptUnlockRequest is a Record type
type HolderServiceAcceptUnlockRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.UnlockRequestAccept `json:"payload"`
}

// ToMap converts HolderServiceAcceptUnlockRequest to a map for DAML arguments
func (t HolderServiceAcceptUnlockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceAcceptUnlockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceAcceptUnlockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceAcceptUnlockRequest to hex string (Canton MCMS format)
func (t HolderServiceAcceptUnlockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceAcceptUnlockRequest from hex string (Canton MCMS format)
func (t *HolderServiceAcceptUnlockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelBurnRequest is a Record type
type HolderServiceCancelBurnRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.BurnRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelBurnRequest to a map for DAML arguments
func (t HolderServiceCancelBurnRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelBurnRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelBurnRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelBurnRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelBurnRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelBurnRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelBurnRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelEnforcementServiceRequest is a Record type
type HolderServiceCancelEnforcementServiceRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload EnforcementServiceRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelEnforcementServiceRequest to a map for DAML arguments
func (t HolderServiceCancelEnforcementServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelEnforcementServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelEnforcementServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelEnforcementServiceRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelEnforcementServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelEnforcementServiceRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelEnforcementServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelForceTransferRequest is a Record type
type HolderServiceCancelForceTransferRequest struct {
	Cid     types.CONTRACT_ID                      `json:"cid"`
	Payload registry_v0.ForceTransferRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelForceTransferRequest to a map for DAML arguments
func (t HolderServiceCancelForceTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelForceTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelForceTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelForceTransferRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelForceTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelForceTransferRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelForceTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelLockOffer is a Record type
type HolderServiceCancelLockOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.LockOfferCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelLockOffer to a map for DAML arguments
func (t HolderServiceCancelLockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelLockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelLockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelLockOffer to hex string (Canton MCMS format)
func (t HolderServiceCancelLockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelLockOffer from hex string (Canton MCMS format)
func (t *HolderServiceCancelLockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelLockRequest is a Record type
type HolderServiceCancelLockRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.LockRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelLockRequest to a map for DAML arguments
func (t HolderServiceCancelLockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelLockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelLockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelLockRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelLockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelLockRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelLockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelMintRequest is a Record type
type HolderServiceCancelMintRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.MintRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelMintRequest to a map for DAML arguments
func (t HolderServiceCancelMintRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelMintRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelMintRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelMintRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelMintRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelMintRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelMintRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelTransferOffer is a Record type
type HolderServiceCancelTransferOffer struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.TransferOfferCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelTransferOffer to a map for DAML arguments
func (t HolderServiceCancelTransferOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelTransferOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelTransferOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelTransferOffer to hex string (Canton MCMS format)
func (t HolderServiceCancelTransferOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelTransferOffer from hex string (Canton MCMS format)
func (t *HolderServiceCancelTransferOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelTransferRequest is a Record type
type HolderServiceCancelTransferRequest struct {
	Cid     types.CONTRACT_ID                 `json:"cid"`
	Payload registry_v0.TransferRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelTransferRequest to a map for DAML arguments
func (t HolderServiceCancelTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelTransferRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelTransferRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelUnlockOffer is a Record type
type HolderServiceCancelUnlockOffer struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.UnlockOfferCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelUnlockOffer to a map for DAML arguments
func (t HolderServiceCancelUnlockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelUnlockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelUnlockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelUnlockOffer to hex string (Canton MCMS format)
func (t HolderServiceCancelUnlockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelUnlockOffer from hex string (Canton MCMS format)
func (t *HolderServiceCancelUnlockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCancelUnlockRequest is a Record type
type HolderServiceCancelUnlockRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.UnlockRequestCancel `json:"payload"`
}

// ToMap converts HolderServiceCancelUnlockRequest to a map for DAML arguments
func (t HolderServiceCancelUnlockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceCancelUnlockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCancelUnlockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCancelUnlockRequest to hex string (Canton MCMS format)
func (t HolderServiceCancelUnlockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCancelUnlockRequest from hex string (Canton MCMS format)
func (t *HolderServiceCancelUnlockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceClean is a Record type
type HolderServiceClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts HolderServiceClean to a map for DAML arguments
func (t HolderServiceClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t HolderServiceClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceClean to hex string (Canton MCMS format)
func (t HolderServiceClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceClean from hex string (Canton MCMS format)
func (t *HolderServiceClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceCreateAllocation is a Record type
type HolderServiceCreateAllocation struct {
	Registrar            types.PARTY                                            `json:"registrar"`
	AllocationFactoryCid types.CONTRACT_ID                                      `json:"allocationFactoryCid"`
	Allocation           splice_api_token_allocation_v1.AllocationSpecification `json:"allocation"`
	InputHoldings        []types.CONTRACT_ID                                    `json:"inputHoldings"`
	ExtraArgs            splice_api_token_metadata_v1.ExtraArgs                 `json:"extraArgs"`
}

// ToMap converts HolderServiceCreateAllocation to a map for DAML arguments
func (t HolderServiceCreateAllocation) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["allocationFactoryCid"] = model.NestedToDAMLValue(t.AllocationFactoryCid)

	m["allocation"] = model.NestedToDAMLValue(t.Allocation)

	m["inputHoldings"] = func() []any {
		res := make([]any, 0, len(t.InputHoldings))
		for _, e := range t.InputHoldings {
			res = append(res, e)
		}
		return res
	}()

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t HolderServiceCreateAllocation) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceCreateAllocation) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceCreateAllocation to hex string (Canton MCMS format)
func (t HolderServiceCreateAllocation) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceCreateAllocation from hex string (Canton MCMS format)
func (t *HolderServiceCreateAllocation) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferLock is a Record type
type HolderServiceOfferLock struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Locker               types.PARTY                              `json:"locker"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Context              types.TEXT                               `json:"context"`
	HoldingLabel         types.TEXT                               `json:"holdingLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceOfferLock to a map for DAML arguments
func (t HolderServiceOfferLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["locker"] = t.Locker.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["context"] = string(t.Context)

	m["holdingLabel"] = string(t.HoldingLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceOfferLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferLock to hex string (Canton MCMS format)
func (t HolderServiceOfferLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferLock from hex string (Canton MCMS format)
func (t *HolderServiceOfferLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferLockResult is a Record type
type HolderServiceOfferLockResult struct {
	LockOfferCid types.CONTRACT_ID `json:"lockOfferCid"`
}

// ToMap converts HolderServiceOfferLockResult to a map for DAML arguments
func (t HolderServiceOfferLockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["lockOfferCid"] = model.NestedToDAMLValue(t.LockOfferCid)

	return m
}

func (t HolderServiceOfferLockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferLockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferLockResult to hex string (Canton MCMS format)
func (t HolderServiceOfferLockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferLockResult from hex string (Canton MCMS format)
func (t *HolderServiceOfferLockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferTransfer is a Record type
type HolderServiceOfferTransfer struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Receiver             types.PARTY                              `json:"receiver"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	SenderLabel          types.TEXT                               `json:"senderLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceOfferTransfer to a map for DAML arguments
func (t HolderServiceOfferTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["senderLabel"] = string(t.SenderLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceOfferTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferTransfer to hex string (Canton MCMS format)
func (t HolderServiceOfferTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferTransfer from hex string (Canton MCMS format)
func (t *HolderServiceOfferTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferTransferResult is a Record type
type HolderServiceOfferTransferResult struct {
	TransferOfferCid types.CONTRACT_ID `json:"transferOfferCid"`
}

// ToMap converts HolderServiceOfferTransferResult to a map for DAML arguments
func (t HolderServiceOfferTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferOfferCid"] = model.NestedToDAMLValue(t.TransferOfferCid)

	return m
}

func (t HolderServiceOfferTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferTransferResult to hex string (Canton MCMS format)
func (t HolderServiceOfferTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferTransferResult from hex string (Canton MCMS format)
func (t *HolderServiceOfferTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferUnlock is a Record type
type HolderServiceOfferUnlock struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Locker               types.PARTY                              `json:"locker"`
	LockContext          types.TEXT                               `json:"lockContext"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	HoldingLabel         types.TEXT                               `json:"holdingLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceOfferUnlock to a map for DAML arguments
func (t HolderServiceOfferUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["locker"] = t.Locker.ToMap()

	m["lockContext"] = string(t.LockContext)

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holdingLabel"] = string(t.HoldingLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceOfferUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferUnlock to hex string (Canton MCMS format)
func (t HolderServiceOfferUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferUnlock from hex string (Canton MCMS format)
func (t *HolderServiceOfferUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceOfferUnlockResult is a Record type
type HolderServiceOfferUnlockResult struct {
	UnlockOfferCid types.CONTRACT_ID `json:"unlockOfferCid"`
}

// ToMap converts HolderServiceOfferUnlockResult to a map for DAML arguments
func (t HolderServiceOfferUnlockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["unlockOfferCid"] = model.NestedToDAMLValue(t.UnlockOfferCid)

	return m
}

func (t HolderServiceOfferUnlockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceOfferUnlockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceOfferUnlockResult to hex string (Canton MCMS format)
func (t HolderServiceOfferUnlockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceOfferUnlockResult from hex string (Canton MCMS format)
func (t *HolderServiceOfferUnlockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectAllocationRequest is a Record type
type HolderServiceRejectAllocationRequest struct {
	AllocationRequestCid types.CONTRACT_ID `json:"allocationRequestCid"`
}

// ToMap converts HolderServiceRejectAllocationRequest to a map for DAML arguments
func (t HolderServiceRejectAllocationRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocationRequestCid"] = model.NestedToDAMLValue(t.AllocationRequestCid)

	return m
}

func (t HolderServiceRejectAllocationRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectAllocationRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectAllocationRequest to hex string (Canton MCMS format)
func (t HolderServiceRejectAllocationRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectAllocationRequest from hex string (Canton MCMS format)
func (t *HolderServiceRejectAllocationRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectBurnOffer is a Record type
type HolderServiceRejectBurnOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.BurnOfferReject `json:"payload"`
}

// ToMap converts HolderServiceRejectBurnOffer to a map for DAML arguments
func (t HolderServiceRejectBurnOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectBurnOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectBurnOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectBurnOffer to hex string (Canton MCMS format)
func (t HolderServiceRejectBurnOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectBurnOffer from hex string (Canton MCMS format)
func (t *HolderServiceRejectBurnOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectLockOffer is a Record type
type HolderServiceRejectLockOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.LockOfferReject `json:"payload"`
}

// ToMap converts HolderServiceRejectLockOffer to a map for DAML arguments
func (t HolderServiceRejectLockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectLockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectLockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectLockOffer to hex string (Canton MCMS format)
func (t HolderServiceRejectLockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectLockOffer from hex string (Canton MCMS format)
func (t *HolderServiceRejectLockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectLockRequest is a Record type
type HolderServiceRejectLockRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.LockRequestReject `json:"payload"`
}

// ToMap converts HolderServiceRejectLockRequest to a map for DAML arguments
func (t HolderServiceRejectLockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectLockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectLockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectLockRequest to hex string (Canton MCMS format)
func (t HolderServiceRejectLockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectLockRequest from hex string (Canton MCMS format)
func (t *HolderServiceRejectLockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectMintOffer is a Record type
type HolderServiceRejectMintOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.MintOfferReject `json:"payload"`
}

// ToMap converts HolderServiceRejectMintOffer to a map for DAML arguments
func (t HolderServiceRejectMintOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectMintOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectMintOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectMintOffer to hex string (Canton MCMS format)
func (t HolderServiceRejectMintOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectMintOffer from hex string (Canton MCMS format)
func (t *HolderServiceRejectMintOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectTransferOffer is a Record type
type HolderServiceRejectTransferOffer struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.TransferOfferReject `json:"payload"`
}

// ToMap converts HolderServiceRejectTransferOffer to a map for DAML arguments
func (t HolderServiceRejectTransferOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectTransferOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectTransferOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectTransferOffer to hex string (Canton MCMS format)
func (t HolderServiceRejectTransferOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectTransferOffer from hex string (Canton MCMS format)
func (t *HolderServiceRejectTransferOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectTransferRequest is a Record type
type HolderServiceRejectTransferRequest struct {
	Cid     types.CONTRACT_ID                 `json:"cid"`
	Payload registry_v0.TransferRequestReject `json:"payload"`
}

// ToMap converts HolderServiceRejectTransferRequest to a map for DAML arguments
func (t HolderServiceRejectTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectTransferRequest to hex string (Canton MCMS format)
func (t HolderServiceRejectTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectTransferRequest from hex string (Canton MCMS format)
func (t *HolderServiceRejectTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectUnlockOffer is a Record type
type HolderServiceRejectUnlockOffer struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.UnlockOfferReject `json:"payload"`
}

// ToMap converts HolderServiceRejectUnlockOffer to a map for DAML arguments
func (t HolderServiceRejectUnlockOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectUnlockOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectUnlockOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectUnlockOffer to hex string (Canton MCMS format)
func (t HolderServiceRejectUnlockOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectUnlockOffer from hex string (Canton MCMS format)
func (t *HolderServiceRejectUnlockOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRejectUnlockRequest is a Record type
type HolderServiceRejectUnlockRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.UnlockRequestReject `json:"payload"`
}

// ToMap converts HolderServiceRejectUnlockRequest to a map for DAML arguments
func (t HolderServiceRejectUnlockRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t HolderServiceRejectUnlockRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRejectUnlockRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRejectUnlockRequest to hex string (Canton MCMS format)
func (t HolderServiceRejectUnlockRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRejectUnlockRequest from hex string (Canton MCMS format)
func (t *HolderServiceRejectUnlockRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestBurn is a Record type
type HolderServiceRequestBurn struct {
	Registrar            types.PARTY                              `json:"registrar"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	HoldingLabel         types.TEXT                               `json:"holdingLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceRequestBurn to a map for DAML arguments
func (t HolderServiceRequestBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holdingLabel"] = string(t.HoldingLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceRequestBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestBurn to hex string (Canton MCMS format)
func (t HolderServiceRequestBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestBurn from hex string (Canton MCMS format)
func (t *HolderServiceRequestBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestBurnResult is a Record type
type HolderServiceRequestBurnResult struct {
	BurnRequestCid types.CONTRACT_ID `json:"burnRequestCid"`
}

// ToMap converts HolderServiceRequestBurnResult to a map for DAML arguments
func (t HolderServiceRequestBurnResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["burnRequestCid"] = model.NestedToDAMLValue(t.BurnRequestCid)

	return m
}

func (t HolderServiceRequestBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestBurnResult to hex string (Canton MCMS format)
func (t HolderServiceRequestBurnResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestBurnResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestBurnResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestEnforcementService is a Record type
type HolderServiceRequestEnforcementService struct {
	Registrar types.PARTY `json:"registrar"`
}

// ToMap converts HolderServiceRequestEnforcementService to a map for DAML arguments
func (t HolderServiceRequestEnforcementService) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	return m
}

func (t HolderServiceRequestEnforcementService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestEnforcementService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestEnforcementService to hex string (Canton MCMS format)
func (t HolderServiceRequestEnforcementService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestEnforcementService from hex string (Canton MCMS format)
func (t *HolderServiceRequestEnforcementService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestEnforcementServiceResult is a Record type
type HolderServiceRequestEnforcementServiceResult struct {
	EnforcementServiceRequestCid types.CONTRACT_ID `json:"enforcementServiceRequestCid"`
}

// ToMap converts HolderServiceRequestEnforcementServiceResult to a map for DAML arguments
func (t HolderServiceRequestEnforcementServiceResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["enforcementServiceRequestCid"] = model.NestedToDAMLValue(t.EnforcementServiceRequestCid)

	return m
}

func (t HolderServiceRequestEnforcementServiceResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestEnforcementServiceResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestEnforcementServiceResult to hex string (Canton MCMS format)
func (t HolderServiceRequestEnforcementServiceResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestEnforcementServiceResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestEnforcementServiceResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestForceTransfer is a Record type
type HolderServiceRequestForceTransfer struct {
	RequestorRationale   types.TEXT                               `json:"requestorRationale"`
	Registrar            types.PARTY                              `json:"registrar"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
	Sender               types.PARTY                              `json:"sender"`
	SenderLabel          types.TEXT                               `json:"senderLabel"`
	Receiver             types.PARTY                              `json:"receiver"`
	ReceiverLabel        types.TEXT                               `json:"receiverLabel"`
}

// ToMap converts HolderServiceRequestForceTransfer to a map for DAML arguments
func (t HolderServiceRequestForceTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["requestorRationale"] = string(t.RequestorRationale)

	m["registrar"] = t.Registrar.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	m["sender"] = t.Sender.ToMap()

	m["senderLabel"] = string(t.SenderLabel)

	m["receiver"] = t.Receiver.ToMap()

	m["receiverLabel"] = string(t.ReceiverLabel)

	return m
}

func (t HolderServiceRequestForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestForceTransfer to hex string (Canton MCMS format)
func (t HolderServiceRequestForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestForceTransfer from hex string (Canton MCMS format)
func (t *HolderServiceRequestForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestForceTransferResult is a Record type
type HolderServiceRequestForceTransferResult struct {
	ForceTransferRequestCid types.CONTRACT_ID `json:"forceTransferRequestCid"`
}

// ToMap converts HolderServiceRequestForceTransferResult to a map for DAML arguments
func (t HolderServiceRequestForceTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["forceTransferRequestCid"] = model.NestedToDAMLValue(t.ForceTransferRequestCid)

	return m
}

func (t HolderServiceRequestForceTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestForceTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestForceTransferResult to hex string (Canton MCMS format)
func (t HolderServiceRequestForceTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestForceTransferResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestForceTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestLock is a Record type
type HolderServiceRequestLock struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Holder               types.PARTY                              `json:"holder"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Context              types.TEXT                               `json:"context"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceRequestLock to a map for DAML arguments
func (t HolderServiceRequestLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["context"] = string(t.Context)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceRequestLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestLock to hex string (Canton MCMS format)
func (t HolderServiceRequestLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestLock from hex string (Canton MCMS format)
func (t *HolderServiceRequestLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestLockResult is a Record type
type HolderServiceRequestLockResult struct {
	LockRequestCid types.CONTRACT_ID `json:"lockRequestCid"`
}

// ToMap converts HolderServiceRequestLockResult to a map for DAML arguments
func (t HolderServiceRequestLockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["lockRequestCid"] = model.NestedToDAMLValue(t.LockRequestCid)

	return m
}

func (t HolderServiceRequestLockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestLockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestLockResult to hex string (Canton MCMS format)
func (t HolderServiceRequestLockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestLockResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestLockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestMint is a Record type
type HolderServiceRequestMint struct {
	Registrar            types.PARTY                              `json:"registrar"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
	HoldingLabel         types.TEXT                               `json:"holdingLabel"`
}

// ToMap converts HolderServiceRequestMint to a map for DAML arguments
func (t HolderServiceRequestMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	m["holdingLabel"] = string(t.HoldingLabel)

	return m
}

func (t HolderServiceRequestMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestMint to hex string (Canton MCMS format)
func (t HolderServiceRequestMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestMint from hex string (Canton MCMS format)
func (t *HolderServiceRequestMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestMintResult is a Record type
type HolderServiceRequestMintResult struct {
	MintRequestCid types.CONTRACT_ID `json:"mintRequestCid"`
}

// ToMap converts HolderServiceRequestMintResult to a map for DAML arguments
func (t HolderServiceRequestMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["mintRequestCid"] = model.NestedToDAMLValue(t.MintRequestCid)

	return m
}

func (t HolderServiceRequestMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestMintResult to hex string (Canton MCMS format)
func (t HolderServiceRequestMintResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestMintResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestMintResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestTransfer is a Record type
type HolderServiceRequestTransfer struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Sender               types.PARTY                              `json:"sender"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	ReceiverLabel        types.TEXT                               `json:"receiverLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceRequestTransfer to a map for DAML arguments
func (t HolderServiceRequestTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["sender"] = t.Sender.ToMap()

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["receiverLabel"] = string(t.ReceiverLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceRequestTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestTransfer to hex string (Canton MCMS format)
func (t HolderServiceRequestTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestTransfer from hex string (Canton MCMS format)
func (t *HolderServiceRequestTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestTransferResult is a Record type
type HolderServiceRequestTransferResult struct {
	TransferRequestCid types.CONTRACT_ID `json:"transferRequestCid"`
}

// ToMap converts HolderServiceRequestTransferResult to a map for DAML arguments
func (t HolderServiceRequestTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferRequestCid"] = model.NestedToDAMLValue(t.TransferRequestCid)

	return m
}

func (t HolderServiceRequestTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestTransferResult to hex string (Canton MCMS format)
func (t HolderServiceRequestTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestTransferResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestUnlock is a Record type
type HolderServiceRequestUnlock struct {
	Registrar            types.PARTY                              `json:"registrar"`
	Holder               types.PARTY                              `json:"holder"`
	LockContext          types.TEXT                               `json:"lockContext"`
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	HoldingLabel         types.TEXT                               `json:"holdingLabel"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts HolderServiceRequestUnlock to a map for DAML arguments
func (t HolderServiceRequestUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrar"] = t.Registrar.ToMap()

	m["holder"] = t.Holder.ToMap()

	m["lockContext"] = string(t.LockContext)

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holdingLabel"] = string(t.HoldingLabel)

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t HolderServiceRequestUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestUnlock to hex string (Canton MCMS format)
func (t HolderServiceRequestUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestUnlock from hex string (Canton MCMS format)
func (t *HolderServiceRequestUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceRequestUnlockResult is a Record type
type HolderServiceRequestUnlockResult struct {
	UnlockRequestCid types.CONTRACT_ID `json:"unlockRequestCid"`
}

// ToMap converts HolderServiceRequestUnlockResult to a map for DAML arguments
func (t HolderServiceRequestUnlockResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["unlockRequestCid"] = model.NestedToDAMLValue(t.UnlockRequestCid)

	return m
}

func (t HolderServiceRequestUnlockResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceRequestUnlockResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceRequestUnlockResult to hex string (Canton MCMS format)
func (t HolderServiceRequestUnlockResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceRequestUnlockResult from hex string (Canton MCMS format)
func (t *HolderServiceRequestUnlockResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceTerminate is a Record type
type HolderServiceTerminate struct {
}

// ToMap converts HolderServiceTerminate to a map for DAML arguments
func (t HolderServiceTerminate) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t HolderServiceTerminate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceTerminate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceTerminate to hex string (Canton MCMS format)
func (t HolderServiceTerminate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceTerminate from hex string (Canton MCMS format)
func (t *HolderServiceTerminate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HolderServiceTerminateResult is a Record type
type HolderServiceTerminateResult struct {
}

// ToMap converts HolderServiceTerminateResult to a map for DAML arguments
func (t HolderServiceTerminateResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t HolderServiceTerminateResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HolderServiceTerminateResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HolderServiceTerminateResult to hex string (Canton MCMS format)
func (t HolderServiceTerminateResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HolderServiceTerminateResult from hex string (Canton MCMS format)
func (t *HolderServiceTerminateResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// InstrumentAllowance is a Record type
type InstrumentAllowance struct {
	Id types.TEXT `json:"id"`
}

// ToMap converts InstrumentAllowance to a map for DAML arguments
func (t InstrumentAllowance) ToMap() map[string]any {
	m := make(map[string]any)

	m["id"] = string(t.Id)

	return m
}

func (t InstrumentAllowance) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *InstrumentAllowance) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes InstrumentAllowance to hex string (Canton MCMS format)
func (t InstrumentAllowance) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes InstrumentAllowance from hex string (Canton MCMS format)
func (t *InstrumentAllowance) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Mint is a Record type
type Mint struct {
	InstrumentId  splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Amount        types.NUMERIC                            `json:"amount"`
	Holder        types.PARTY                              `json:"holder"`
	Reference     types.TEXT                               `json:"reference"`
	RequestedAt   types.TIMESTAMP                          `json:"requestedAt"`
	ExecuteBefore types.TIMESTAMP                          `json:"executeBefore"`
	Meta          splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// ToMap converts Mint to a map for DAML arguments
func (t Mint) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["requestedAt"] = t.RequestedAt

	m["executeBefore"] = t.ExecuteBefore

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
	Mint     Mint        `json:"mint"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

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
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

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

// MintOfferAccept exercises the MintOffer_Accept choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferAccept(contractID string, args MintOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintOfferAcceptWithPackageID exercises the MintOffer_Accept choice using the provided package ID instead of package name
func (t MintOffer) MintOfferAcceptWithPackageID(contractID string, packageID string, args MintOfferAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintOfferReject exercises the MintOffer_Reject choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferReject(contractID string, args MintOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintOfferRejectWithPackageID exercises the MintOffer_Reject choice using the provided package ID instead of package name
func (t MintOffer) MintOfferRejectWithPackageID(contractID string, packageID string, args MintOfferReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintOfferCancel exercises the MintOffer_Cancel choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) MintOfferCancel(contractID string, args MintOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// MintOfferCancelWithPackageID exercises the MintOffer_Cancel choice using the provided package ID instead of package name
func (t MintOffer) MintOfferCancelWithPackageID(contractID string, packageID string, args MintOfferCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "MintOffer_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintOffer contract
// This method uses the package name in the template ID
func (t MintOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintOffer"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintOfferAccept is a Record type
type MintOfferAccept struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts MintOfferAccept to a map for DAML arguments
func (t MintOfferAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

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
	HoldingCid      types.CONTRACT_ID                     `json:"holdingCid"`
	ExecutedMintCid *types.CONTRACT_ID                    `json:"executedMintCid" hex:"optional"`
	Meta            splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts MintOfferAcceptResult to a map for DAML arguments
func (t MintOfferAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	if t.ExecutedMintCid != nil {
		m["executedMintCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutedMintCid),
		}
	} else {
		m["executedMintCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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

// MintOfferReject is a Record type
type MintOfferReject struct {
	Reason    types.TEXT                              `json:"reason"`
	ExtraArgs *splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs" hex:"optional"`
}

// ToMap converts MintOfferReject to a map for DAML arguments
func (t MintOfferReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	if t.ExtraArgs != nil {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExtraArgs),
		}
	} else {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	RejectedMintCid *types.CONTRACT_ID `json:"rejectedMintCid" hex:"optional"`
}

// ToMap converts MintOfferRejectResult to a map for DAML arguments
func (t MintOfferRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RejectedMintCid != nil {
		m["rejectedMintCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RejectedMintCid),
		}
	} else {
		m["rejectedMintCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
	Mint     Mint        `json:"mint"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MintRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MintRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MintRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

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

// MintRequestAccept exercises the MintRequest_Accept choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestAccept(contractID string, args MintRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintRequestAcceptWithPackageID exercises the MintRequest_Accept choice using the provided package ID instead of package name
func (t MintRequest) MintRequestAcceptWithPackageID(contractID string, packageID string, args MintRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// MintRequestReject exercises the MintRequest_Reject choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestReject(contractID string, args MintRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// MintRequestRejectWithPackageID exercises the MintRequest_Reject choice using the provided package ID instead of package name
func (t MintRequest) MintRequestRejectWithPackageID(contractID string, packageID string, args MintRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MintRequestCancel exercises the MintRequest_Cancel choice on this MintRequest contract
// This method uses the package name in the template ID
func (t MintRequest) MintRequestCancel(contractID string, args MintRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// MintRequestCancelWithPackageID exercises the MintRequest_Cancel choice using the provided package ID instead of package name
func (t MintRequest) MintRequestCancelWithPackageID(contractID string, packageID string, args MintRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "MintRequest"),
		ContractID: contractID,
		Choice:     "MintRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// MintRequestAccept is a Record type
type MintRequestAccept struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts MintRequestAccept to a map for DAML arguments
func (t MintRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

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
	HoldingCid      types.CONTRACT_ID                     `json:"holdingCid"`
	ExecutedMintCid *types.CONTRACT_ID                    `json:"executedMintCid" hex:"optional"`
	Meta            splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts MintRequestAcceptResult to a map for DAML arguments
func (t MintRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["holdingCid"] = model.NestedToDAMLValue(t.HoldingCid)

	if t.ExecutedMintCid != nil {
		m["executedMintCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutedMintCid),
		}
	} else {
		m["executedMintCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["meta"] = model.NestedToDAMLValue(t.Meta)

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

// MintRequestReject is a Record type
type MintRequestReject struct {
	Reason    types.TEXT                              `json:"reason"`
	ExtraArgs *splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs" hex:"optional"`
}

// ToMap converts MintRequestReject to a map for DAML arguments
func (t MintRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	if t.ExtraArgs != nil {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExtraArgs),
		}
	} else {
		m["extraArgs"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	RejectedMintCid *types.CONTRACT_ID `json:"rejectedMintCid" hex:"optional"`
}

// ToMap converts MintRequestRejectResult to a map for DAML arguments
func (t MintRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	if t.RejectedMintCid != nil {
		m["rejectedMintCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RejectedMintCid),
		}
	} else {
		m["rejectedMintCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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

// OperatorConfiguration is a Template type
type OperatorConfiguration struct {
	Operator             types.PARTY                                `json:"operator"`
	ProviderRequirements []credential_v0.PartyCredentialRequirement `json:"providerRequirements"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OperatorConfiguration) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OperatorConfiguration) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OperatorConfiguration) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["providerRequirements"] = func() []any {
		res := make([]any, 0, len(t.ProviderRequirements))
		for _, e := range t.ProviderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t OperatorConfiguration) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["providerRequirements"] = func() []any {
		res := make([]any, 0, len(t.ProviderRequirements))
		for _, e := range t.ProviderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t OperatorConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OperatorConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OperatorConfiguration to hex string (Canton MCMS format)
func (t OperatorConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OperatorConfiguration from hex string (Canton MCMS format)
func (t *OperatorConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for OperatorConfiguration

// OperatorConfigurationGet exercises the OperatorConfiguration_Get choice on this OperatorConfiguration contract
// This method uses the package name in the template ID
func (t OperatorConfiguration) OperatorConfigurationGet(contractID string, args OperatorConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "OperatorConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// OperatorConfigurationGetWithPackageID exercises the OperatorConfiguration_Get choice using the provided package ID instead of package name
func (t OperatorConfiguration) OperatorConfigurationGetWithPackageID(contractID string, packageID string, args OperatorConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "OperatorConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OperatorConfiguration contract
// This method uses the package name in the template ID
func (t OperatorConfiguration) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OperatorConfiguration) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// OperatorConfigurationModify exercises the OperatorConfiguration_Modify choice on this OperatorConfiguration contract
// This method uses the package name in the template ID
func (t OperatorConfiguration) OperatorConfigurationModify(contractID string, args OperatorConfigurationModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "OperatorConfiguration_Modify",
		Arguments:  argsToMap(args),
	}
}

// OperatorConfigurationModifyWithPackageID exercises the OperatorConfiguration_Modify choice using the provided package ID instead of package name
func (t OperatorConfiguration) OperatorConfigurationModifyWithPackageID(contractID string, packageID string, args OperatorConfigurationModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Operator", "OperatorConfiguration"),
		ContractID: contractID,
		Choice:     "OperatorConfiguration_Modify",
		Arguments:  argsToMap(args),
	}
}

// OperatorConfigurationGet is a Record type
type OperatorConfigurationGet struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts OperatorConfigurationGet to a map for DAML arguments
func (t OperatorConfigurationGet) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t OperatorConfigurationGet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OperatorConfigurationGet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OperatorConfigurationGet to hex string (Canton MCMS format)
func (t OperatorConfigurationGet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OperatorConfigurationGet from hex string (Canton MCMS format)
func (t *OperatorConfigurationGet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// OperatorConfigurationGetResult is a Record type
type OperatorConfigurationGetResult struct {
	OperatorConfiguration OperatorConfiguration `json:"operatorConfiguration"`
}

// ToMap converts OperatorConfigurationGetResult to a map for DAML arguments
func (t OperatorConfigurationGetResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["operatorConfiguration"] = model.NestedToDAMLValue(t.OperatorConfiguration)

	return m
}

func (t OperatorConfigurationGetResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OperatorConfigurationGetResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OperatorConfigurationGetResult to hex string (Canton MCMS format)
func (t OperatorConfigurationGetResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OperatorConfigurationGetResult from hex string (Canton MCMS format)
func (t *OperatorConfigurationGetResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// OperatorConfigurationModify is a Record type
type OperatorConfigurationModify struct {
	NewProviderRequirements []credential_v0.PartyCredentialRequirement `json:"newProviderRequirements"`
}

// ToMap converts OperatorConfigurationModify to a map for DAML arguments
func (t OperatorConfigurationModify) ToMap() map[string]any {
	m := make(map[string]any)

	m["newProviderRequirements"] = func() []any {
		res := make([]any, 0, len(t.NewProviderRequirements))
		for _, e := range t.NewProviderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t OperatorConfigurationModify) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OperatorConfigurationModify) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OperatorConfigurationModify to hex string (Canton MCMS format)
func (t OperatorConfigurationModify) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OperatorConfigurationModify from hex string (Canton MCMS format)
func (t *OperatorConfigurationModify) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// OperatorConfigurationModifyResult is a Record type
type OperatorConfigurationModifyResult struct {
	OperatorConfigurationCid types.CONTRACT_ID `json:"operatorConfigurationCid"`
}

// ToMap converts OperatorConfigurationModifyResult to a map for DAML arguments
func (t OperatorConfigurationModifyResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["operatorConfigurationCid"] = model.NestedToDAMLValue(t.OperatorConfigurationCid)

	return m
}

func (t OperatorConfigurationModifyResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OperatorConfigurationModifyResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OperatorConfigurationModifyResult to hex string (Canton MCMS format)
func (t OperatorConfigurationModifyResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OperatorConfigurationModifyResult from hex string (Canton MCMS format)
func (t *OperatorConfigurationModifyResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderConfiguration is a Template type
type ProviderConfiguration struct {
	Operator              types.PARTY                                `json:"operator"`
	Provider              types.PARTY                                `json:"provider"`
	RegistrarRequirements []credential_v0.PartyCredentialRequirement `json:"registrarRequirements"`
	HolderRequirements    []credential_v0.PartyCredentialRequirement `json:"holderRequirements"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ProviderConfiguration) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ProviderConfiguration) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ProviderConfiguration) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRequirements"] = func() []any {
		res := make([]any, 0, len(t.RegistrarRequirements))
		for _, e := range t.RegistrarRequirements {
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

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ProviderConfiguration) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrarRequirements"] = func() []any {
		res := make([]any, 0, len(t.RegistrarRequirements))
		for _, e := range t.RegistrarRequirements {
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

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ProviderConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderConfiguration to hex string (Canton MCMS format)
func (t ProviderConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderConfiguration from hex string (Canton MCMS format)
func (t *ProviderConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ProviderConfiguration

// Archive exercises the Archive choice on this ProviderConfiguration contract
// This method uses the package name in the template ID
func (t ProviderConfiguration) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ProviderConfiguration) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ProviderConfigurationGet exercises the ProviderConfiguration_Get choice on this ProviderConfiguration contract
// This method uses the package name in the template ID
func (t ProviderConfiguration) ProviderConfigurationGet(contractID string, args ProviderConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration"),
		ContractID: contractID,
		Choice:     "ProviderConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// ProviderConfigurationGetWithPackageID exercises the ProviderConfiguration_Get choice using the provided package ID instead of package name
func (t ProviderConfiguration) ProviderConfigurationGetWithPackageID(contractID string, packageID string, args ProviderConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Provider", "ProviderConfiguration"),
		ContractID: contractID,
		Choice:     "ProviderConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// ProviderConfigurationGet is a Record type
type ProviderConfigurationGet struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts ProviderConfigurationGet to a map for DAML arguments
func (t ProviderConfigurationGet) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t ProviderConfigurationGet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderConfigurationGet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderConfigurationGet to hex string (Canton MCMS format)
func (t ProviderConfigurationGet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderConfigurationGet from hex string (Canton MCMS format)
func (t *ProviderConfigurationGet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderConfigurationGetResult is a Record type
type ProviderConfigurationGetResult struct {
	ProviderConfiguration ProviderConfiguration `json:"providerConfiguration"`
}

// ToMap converts ProviderConfigurationGetResult to a map for DAML arguments
func (t ProviderConfigurationGetResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerConfiguration"] = model.NestedToDAMLValue(t.ProviderConfiguration)

	return m
}

func (t ProviderConfigurationGetResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderConfigurationGetResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderConfigurationGetResult to hex string (Canton MCMS format)
func (t ProviderConfigurationGetResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderConfigurationGetResult from hex string (Canton MCMS format)
func (t *ProviderConfigurationGetResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderService is a Template type
type ProviderService struct {
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ProviderService) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ProviderService) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ProviderService) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ProviderService) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ProviderService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderService to hex string (Canton MCMS format)
func (t ProviderService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderService from hex string (Canton MCMS format)
func (t *ProviderService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ProviderService

// ProviderServiceTerminate exercises the ProviderService_Terminate choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceTerminate(contractID string, args ProviderServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceTerminateWithPackageID exercises the ProviderService_Terminate choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceTerminateWithPackageID(contractID string, packageID string, args ProviderServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceCreateProviderConfiguration exercises the ProviderService_CreateProviderConfiguration choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceCreateProviderConfiguration(contractID string, args ProviderServiceCreateProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_CreateProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceCreateProviderConfigurationWithPackageID exercises the ProviderService_CreateProviderConfiguration choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceCreateProviderConfigurationWithPackageID(contractID string, packageID string, args ProviderServiceCreateProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_CreateProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceArchiveProviderConfiguration exercises the ProviderService_ArchiveProviderConfiguration choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceArchiveProviderConfiguration(contractID string, args ProviderServiceArchiveProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_ArchiveProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceArchiveProviderConfigurationWithPackageID exercises the ProviderService_ArchiveProviderConfiguration choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceArchiveProviderConfigurationWithPackageID(contractID string, packageID string, args ProviderServiceArchiveProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_ArchiveProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceArchiveAndCreateProviderConfiguration exercises the ProviderService_ArchiveAndCreateProviderConfiguration choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceArchiveAndCreateProviderConfiguration(contractID string, args ProviderServiceArchiveAndCreateProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_ArchiveAndCreateProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceArchiveAndCreateProviderConfigurationWithPackageID exercises the ProviderService_ArchiveAndCreateProviderConfiguration choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceArchiveAndCreateProviderConfigurationWithPackageID(contractID string, packageID string, args ProviderServiceArchiveAndCreateProviderConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_ArchiveAndCreateProviderConfiguration",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceAcceptRegistrarServiceRequest exercises the ProviderService_AcceptRegistrarServiceRequest choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceAcceptRegistrarServiceRequest(contractID string, args ProviderServiceAcceptRegistrarServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_AcceptRegistrarServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceAcceptRegistrarServiceRequestWithPackageID exercises the ProviderService_AcceptRegistrarServiceRequest choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceAcceptRegistrarServiceRequestWithPackageID(contractID string, packageID string, args ProviderServiceAcceptRegistrarServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_AcceptRegistrarServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRejectRegistrarServiceRequest exercises the ProviderService_RejectRegistrarServiceRequest choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceRejectRegistrarServiceRequest(contractID string, args ProviderServiceRejectRegistrarServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_RejectRegistrarServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRejectRegistrarServiceRequestWithPackageID exercises the ProviderService_RejectRegistrarServiceRequest choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceRejectRegistrarServiceRequestWithPackageID(contractID string, packageID string, args ProviderServiceRejectRegistrarServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_RejectRegistrarServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceAcceptHolderServiceRequest exercises the ProviderService_AcceptHolderServiceRequest choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceAcceptHolderServiceRequest(contractID string, args ProviderServiceAcceptHolderServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_AcceptHolderServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceAcceptHolderServiceRequestWithPackageID exercises the ProviderService_AcceptHolderServiceRequest choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceAcceptHolderServiceRequestWithPackageID(contractID string, packageID string, args ProviderServiceAcceptHolderServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_AcceptHolderServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRejectHolderServiceRequest exercises the ProviderService_RejectHolderServiceRequest choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) ProviderServiceRejectHolderServiceRequest(contractID string, args ProviderServiceRejectHolderServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_RejectHolderServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRejectHolderServiceRequestWithPackageID exercises the ProviderService_RejectHolderServiceRequest choice using the provided package ID instead of package name
func (t ProviderService) ProviderServiceRejectHolderServiceRequestWithPackageID(contractID string, packageID string, args ProviderServiceRejectHolderServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "ProviderService_RejectHolderServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ProviderService contract
// This method uses the package name in the template ID
func (t ProviderService) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ProviderService) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ProviderServiceRequest is a Template type
type ProviderServiceRequest struct {
	Operator types.PARTY `json:"operator"`
	Provider types.PARTY `json:"provider"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ProviderServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ProviderServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ProviderServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ProviderServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ProviderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequest to hex string (Canton MCMS format)
func (t ProviderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequest from hex string (Canton MCMS format)
func (t *ProviderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ProviderServiceRequest

// ProviderServiceRequestAccept exercises the ProviderServiceRequest_Accept choice on this ProviderServiceRequest contract
// This method uses the package name in the template ID
func (t ProviderServiceRequest) ProviderServiceRequestAccept(contractID string, args ProviderServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRequestAcceptWithPackageID exercises the ProviderServiceRequest_Accept choice using the provided package ID instead of package name
func (t ProviderServiceRequest) ProviderServiceRequestAcceptWithPackageID(contractID string, packageID string, args ProviderServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRequestReject exercises the ProviderServiceRequest_Reject choice on this ProviderServiceRequest contract
// This method uses the package name in the template ID
func (t ProviderServiceRequest) ProviderServiceRequestReject(contractID string, args ProviderServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRequestRejectWithPackageID exercises the ProviderServiceRequest_Reject choice using the provided package ID instead of package name
func (t ProviderServiceRequest) ProviderServiceRequestRejectWithPackageID(contractID string, packageID string, args ProviderServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRequestCancel exercises the ProviderServiceRequest_Cancel choice on this ProviderServiceRequest contract
// This method uses the package name in the template ID
func (t ProviderServiceRequest) ProviderServiceRequestCancel(contractID string, args ProviderServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// ProviderServiceRequestCancelWithPackageID exercises the ProviderServiceRequest_Cancel choice using the provided package ID instead of package name
func (t ProviderServiceRequest) ProviderServiceRequestCancelWithPackageID(contractID string, packageID string, args ProviderServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "ProviderServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ProviderServiceRequest contract
// This method uses the package name in the template ID
func (t ProviderServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ProviderServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "ProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ProviderServiceRequestAccept is a Record type
type ProviderServiceRequestAccept struct {
	OperatorConfigurationCid      types.CONTRACT_ID                          `json:"operatorConfigurationCid"`
	CredentialCids                []types.CONTRACT_ID                        `json:"credentialCids"`
	AppRewardConfigurationDetails *registry_v0.AppRewardConfigurationDetails `json:"appRewardConfigurationDetails" hex:"optional"`
}

// ToMap converts ProviderServiceRequestAccept to a map for DAML arguments
func (t ProviderServiceRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["operatorConfigurationCid"] = model.NestedToDAMLValue(t.OperatorConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	if t.AppRewardConfigurationDetails != nil {
		m["appRewardConfigurationDetails"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.AppRewardConfigurationDetails),
		}
	} else {
		m["appRewardConfigurationDetails"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t ProviderServiceRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestAccept to hex string (Canton MCMS format)
func (t ProviderServiceRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestAccept from hex string (Canton MCMS format)
func (t *ProviderServiceRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRequestAcceptResult is a Record type
type ProviderServiceRequestAcceptResult struct {
	ProviderServiceCid        types.CONTRACT_ID  `json:"providerServiceCid"`
	AppRewardConfigurationCid *types.CONTRACT_ID `json:"appRewardConfigurationCid" hex:"optional"`
}

// ToMap converts ProviderServiceRequestAcceptResult to a map for DAML arguments
func (t ProviderServiceRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerServiceCid"] = model.NestedToDAMLValue(t.ProviderServiceCid)

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

	return m
}

func (t ProviderServiceRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestAcceptResult to hex string (Canton MCMS format)
func (t ProviderServiceRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestAcceptResult from hex string (Canton MCMS format)
func (t *ProviderServiceRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRequestCancel is a Record type
type ProviderServiceRequestCancel struct {
}

// ToMap converts ProviderServiceRequestCancel to a map for DAML arguments
func (t ProviderServiceRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ProviderServiceRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestCancel to hex string (Canton MCMS format)
func (t ProviderServiceRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestCancel from hex string (Canton MCMS format)
func (t *ProviderServiceRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRequestCancelResult is a Record type
type ProviderServiceRequestCancelResult struct {
}

// ToMap converts ProviderServiceRequestCancelResult to a map for DAML arguments
func (t ProviderServiceRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ProviderServiceRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestCancelResult to hex string (Canton MCMS format)
func (t ProviderServiceRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestCancelResult from hex string (Canton MCMS format)
func (t *ProviderServiceRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRequestReject is a Record type
type ProviderServiceRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts ProviderServiceRequestReject to a map for DAML arguments
func (t ProviderServiceRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t ProviderServiceRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestReject to hex string (Canton MCMS format)
func (t ProviderServiceRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestReject from hex string (Canton MCMS format)
func (t *ProviderServiceRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRequestRejectResult is a Record type
type ProviderServiceRequestRejectResult struct {
	RejectedProviderServiceRequestCid types.CONTRACT_ID `json:"rejectedProviderServiceRequestCid"`
}

// ToMap converts ProviderServiceRequestRejectResult to a map for DAML arguments
func (t ProviderServiceRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedProviderServiceRequestCid"] = model.NestedToDAMLValue(t.RejectedProviderServiceRequestCid)

	return m
}

func (t ProviderServiceRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRequestRejectResult to hex string (Canton MCMS format)
func (t ProviderServiceRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRequestRejectResult from hex string (Canton MCMS format)
func (t *ProviderServiceRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceAcceptHolderServiceRequest is a Record type
type ProviderServiceAcceptHolderServiceRequest struct {
	Cid     types.CONTRACT_ID          `json:"cid"`
	Payload HolderServiceRequestAccept `json:"payload"`
}

// ToMap converts ProviderServiceAcceptHolderServiceRequest to a map for DAML arguments
func (t ProviderServiceAcceptHolderServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t ProviderServiceAcceptHolderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceAcceptHolderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceAcceptHolderServiceRequest to hex string (Canton MCMS format)
func (t ProviderServiceAcceptHolderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceAcceptHolderServiceRequest from hex string (Canton MCMS format)
func (t *ProviderServiceAcceptHolderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceAcceptRegistrarServiceRequest is a Record type
type ProviderServiceAcceptRegistrarServiceRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload RegistrarServiceRequestAccept `json:"payload"`
}

// ToMap converts ProviderServiceAcceptRegistrarServiceRequest to a map for DAML arguments
func (t ProviderServiceAcceptRegistrarServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t ProviderServiceAcceptRegistrarServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceAcceptRegistrarServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceAcceptRegistrarServiceRequest to hex string (Canton MCMS format)
func (t ProviderServiceAcceptRegistrarServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceAcceptRegistrarServiceRequest from hex string (Canton MCMS format)
func (t *ProviderServiceAcceptRegistrarServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceArchiveAndCreateProviderConfiguration is a Record type
type ProviderServiceArchiveAndCreateProviderConfiguration struct {
	Cid     types.CONTRACT_ID                          `json:"cid"`
	Payload ProviderServiceCreateProviderConfiguration `json:"payload"`
}

// ToMap converts ProviderServiceArchiveAndCreateProviderConfiguration to a map for DAML arguments
func (t ProviderServiceArchiveAndCreateProviderConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t ProviderServiceArchiveAndCreateProviderConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceArchiveAndCreateProviderConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceArchiveAndCreateProviderConfiguration to hex string (Canton MCMS format)
func (t ProviderServiceArchiveAndCreateProviderConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceArchiveAndCreateProviderConfiguration from hex string (Canton MCMS format)
func (t *ProviderServiceArchiveAndCreateProviderConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceArchiveAndCreateProviderConfigurationResult is a Record type
type ProviderServiceArchiveAndCreateProviderConfigurationResult struct {
	ArchiveResult ProviderServiceArchiveProviderConfigurationResult `json:"archiveResult"`
	CreateResult  ProviderServiceCreateProviderConfigurationResult  `json:"createResult"`
}

// ToMap converts ProviderServiceArchiveAndCreateProviderConfigurationResult to a map for DAML arguments
func (t ProviderServiceArchiveAndCreateProviderConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["archiveResult"] = model.NestedToDAMLValue(t.ArchiveResult)

	m["createResult"] = model.NestedToDAMLValue(t.CreateResult)

	return m
}

func (t ProviderServiceArchiveAndCreateProviderConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceArchiveAndCreateProviderConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceArchiveAndCreateProviderConfigurationResult to hex string (Canton MCMS format)
func (t ProviderServiceArchiveAndCreateProviderConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceArchiveAndCreateProviderConfigurationResult from hex string (Canton MCMS format)
func (t *ProviderServiceArchiveAndCreateProviderConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceArchiveProviderConfiguration is a Record type
type ProviderServiceArchiveProviderConfiguration struct {
	Cid types.CONTRACT_ID `json:"cid"`
}

// ToMap converts ProviderServiceArchiveProviderConfiguration to a map for DAML arguments
func (t ProviderServiceArchiveProviderConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	return m
}

func (t ProviderServiceArchiveProviderConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceArchiveProviderConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceArchiveProviderConfiguration to hex string (Canton MCMS format)
func (t ProviderServiceArchiveProviderConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceArchiveProviderConfiguration from hex string (Canton MCMS format)
func (t *ProviderServiceArchiveProviderConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceArchiveProviderConfigurationResult is a Record type
type ProviderServiceArchiveProviderConfigurationResult struct {
}

// ToMap converts ProviderServiceArchiveProviderConfigurationResult to a map for DAML arguments
func (t ProviderServiceArchiveProviderConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ProviderServiceArchiveProviderConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceArchiveProviderConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceArchiveProviderConfigurationResult to hex string (Canton MCMS format)
func (t ProviderServiceArchiveProviderConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceArchiveProviderConfigurationResult from hex string (Canton MCMS format)
func (t *ProviderServiceArchiveProviderConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceCreateProviderConfiguration is a Record type
type ProviderServiceCreateProviderConfiguration struct {
	RegistrarRequirements []credential_v0.PartyCredentialRequirement `json:"registrarRequirements"`
	HolderRequirements    []credential_v0.PartyCredentialRequirement `json:"holderRequirements"`
}

// ToMap converts ProviderServiceCreateProviderConfiguration to a map for DAML arguments
func (t ProviderServiceCreateProviderConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarRequirements"] = func() []any {
		res := make([]any, 0, len(t.RegistrarRequirements))
		for _, e := range t.RegistrarRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["holderRequirements"] = func() []any {
		res := make([]any, 0, len(t.HolderRequirements))
		for _, e := range t.HolderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t ProviderServiceCreateProviderConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceCreateProviderConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceCreateProviderConfiguration to hex string (Canton MCMS format)
func (t ProviderServiceCreateProviderConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceCreateProviderConfiguration from hex string (Canton MCMS format)
func (t *ProviderServiceCreateProviderConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceCreateProviderConfigurationResult is a Record type
type ProviderServiceCreateProviderConfigurationResult struct {
	ProviderConfigurationCid types.CONTRACT_ID `json:"providerConfigurationCid"`
}

// ToMap converts ProviderServiceCreateProviderConfigurationResult to a map for DAML arguments
func (t ProviderServiceCreateProviderConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerConfigurationCid"] = model.NestedToDAMLValue(t.ProviderConfigurationCid)

	return m
}

func (t ProviderServiceCreateProviderConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceCreateProviderConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceCreateProviderConfigurationResult to hex string (Canton MCMS format)
func (t ProviderServiceCreateProviderConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceCreateProviderConfigurationResult from hex string (Canton MCMS format)
func (t *ProviderServiceCreateProviderConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceCreateResult is a Record type
type ProviderServiceCreateResult struct {
	ProviderServiceCid types.CONTRACT_ID `json:"providerServiceCid"`
}

// ToMap converts ProviderServiceCreateResult to a map for DAML arguments
func (t ProviderServiceCreateResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerServiceCid"] = model.NestedToDAMLValue(t.ProviderServiceCid)

	return m
}

func (t ProviderServiceCreateResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceCreateResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceCreateResult to hex string (Canton MCMS format)
func (t ProviderServiceCreateResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceCreateResult from hex string (Canton MCMS format)
func (t *ProviderServiceCreateResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRejectHolderServiceRequest is a Record type
type ProviderServiceRejectHolderServiceRequest struct {
	Cid     types.CONTRACT_ID          `json:"cid"`
	Payload HolderServiceRequestReject `json:"payload"`
}

// ToMap converts ProviderServiceRejectHolderServiceRequest to a map for DAML arguments
func (t ProviderServiceRejectHolderServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t ProviderServiceRejectHolderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRejectHolderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRejectHolderServiceRequest to hex string (Canton MCMS format)
func (t ProviderServiceRejectHolderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRejectHolderServiceRequest from hex string (Canton MCMS format)
func (t *ProviderServiceRejectHolderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceRejectRegistrarServiceRequest is a Record type
type ProviderServiceRejectRegistrarServiceRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload RegistrarServiceRequestReject `json:"payload"`
}

// ToMap converts ProviderServiceRejectRegistrarServiceRequest to a map for DAML arguments
func (t ProviderServiceRejectRegistrarServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t ProviderServiceRejectRegistrarServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceRejectRegistrarServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceRejectRegistrarServiceRequest to hex string (Canton MCMS format)
func (t ProviderServiceRejectRegistrarServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceRejectRegistrarServiceRequest from hex string (Canton MCMS format)
func (t *ProviderServiceRejectRegistrarServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceTerminate is a Record type
type ProviderServiceTerminate struct {
}

// ToMap converts ProviderServiceTerminate to a map for DAML arguments
func (t ProviderServiceTerminate) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ProviderServiceTerminate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceTerminate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceTerminate to hex string (Canton MCMS format)
func (t ProviderServiceTerminate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceTerminate from hex string (Canton MCMS format)
func (t *ProviderServiceTerminate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ProviderServiceTerminateResult is a Record type
type ProviderServiceTerminateResult struct {
}

// ToMap converts ProviderServiceTerminateResult to a map for DAML arguments
func (t ProviderServiceTerminateResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t ProviderServiceTerminateResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ProviderServiceTerminateResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ProviderServiceTerminateResult to hex string (Canton MCMS format)
func (t ProviderServiceTerminateResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ProviderServiceTerminateResult from hex string (Canton MCMS format)
func (t *ProviderServiceTerminateResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarConfiguration is a Template type
type RegistrarConfiguration struct {
	Operator                types.PARTY                                `json:"operator"`
	Provider                types.PARTY                                `json:"provider"`
	Registrar               types.PARTY                                `json:"registrar"`
	EnforcementRequirements []credential_v0.PartyCredentialRequirement `json:"enforcementRequirements"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RegistrarConfiguration) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RegistrarConfiguration) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RegistrarConfiguration) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["enforcementRequirements"] = func() []any {
		res := make([]any, 0, len(t.EnforcementRequirements))
		for _, e := range t.EnforcementRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RegistrarConfiguration) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["enforcementRequirements"] = func() []any {
		res := make([]any, 0, len(t.EnforcementRequirements))
		for _, e := range t.EnforcementRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RegistrarConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarConfiguration to hex string (Canton MCMS format)
func (t RegistrarConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarConfiguration from hex string (Canton MCMS format)
func (t *RegistrarConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RegistrarConfiguration

// Archive exercises the Archive choice on this RegistrarConfiguration contract
// This method uses the package name in the template ID
func (t RegistrarConfiguration) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RegistrarConfiguration) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RegistrarConfigurationGet exercises the RegistrarConfiguration_Get choice on this RegistrarConfiguration contract
// This method uses the package name in the template ID
func (t RegistrarConfiguration) RegistrarConfigurationGet(contractID string, args RegistrarConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration"),
		ContractID: contractID,
		Choice:     "RegistrarConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// RegistrarConfigurationGetWithPackageID exercises the RegistrarConfiguration_Get choice using the provided package ID instead of package name
func (t RegistrarConfiguration) RegistrarConfigurationGetWithPackageID(contractID string, packageID string, args RegistrarConfigurationGet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Configuration.Registrar", "RegistrarConfiguration"),
		ContractID: contractID,
		Choice:     "RegistrarConfiguration_Get",
		Arguments:  argsToMap(args),
	}
}

// RegistrarConfigurationGet is a Record type
type RegistrarConfigurationGet struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RegistrarConfigurationGet to a map for DAML arguments
func (t RegistrarConfigurationGet) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarConfigurationGet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarConfigurationGet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarConfigurationGet to hex string (Canton MCMS format)
func (t RegistrarConfigurationGet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarConfigurationGet from hex string (Canton MCMS format)
func (t *RegistrarConfigurationGet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarConfigurationGetResult is a Record type
type RegistrarConfigurationGetResult struct {
	RegistrarConfiguration RegistrarConfiguration `json:"registrarConfiguration"`
}

// ToMap converts RegistrarConfigurationGetResult to a map for DAML arguments
func (t RegistrarConfigurationGetResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarConfiguration"] = model.NestedToDAMLValue(t.RegistrarConfiguration)

	return m
}

func (t RegistrarConfigurationGetResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarConfigurationGetResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarConfigurationGetResult to hex string (Canton MCMS format)
func (t RegistrarConfigurationGetResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarConfigurationGetResult from hex string (Canton MCMS format)
func (t *RegistrarConfigurationGetResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarService is a Template type
type RegistrarService struct {
	Operator              types.PARTY `json:"operator"`
	Provider              types.PARTY `json:"provider"`
	Registrar             types.PARTY `json:"registrar"`
	EnableResultContracts *types.BOOL `json:"enableResultContracts" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RegistrarService) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RegistrarService) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RegistrarService) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	if t.EnableResultContracts != nil {
		args["enableResultContracts"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.EnableResultContracts),
		}
	} else {
		args["enableResultContracts"] = map[string]any{
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
func (t RegistrarService) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	if t.EnableResultContracts != nil {
		args["enableResultContracts"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.EnableResultContracts),
		}
	} else {
		args["enableResultContracts"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RegistrarService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarService to hex string (Canton MCMS format)
func (t RegistrarService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarService from hex string (Canton MCMS format)
func (t *RegistrarService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RegistrarService

// RegistrarServiceSet exercises the RegistrarService_Set choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceSet(contractID string, args RegistrarServiceSet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_Set",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceSetWithPackageID exercises the RegistrarService_Set choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceSetWithPackageID(contractID string, packageID string, args RegistrarServiceSet) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_Set",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceTerminate exercises the RegistrarService_Terminate choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceTerminate(contractID string, args RegistrarServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceTerminateWithPackageID exercises the RegistrarService_Terminate choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceTerminateWithPackageID(contractID string, packageID string, args RegistrarServiceTerminate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_Terminate",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptEnforcementServiceRequest exercises the RegistrarService_AcceptEnforcementServiceRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceAcceptEnforcementServiceRequest(contractID string, args RegistrarServiceAcceptEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptEnforcementServiceRequestWithPackageID exercises the RegistrarService_AcceptEnforcementServiceRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceAcceptEnforcementServiceRequestWithPackageID(contractID string, packageID string, args RegistrarServiceAcceptEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectEnforcementServiceRequest exercises the RegistrarService_RejectEnforcementServiceRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceRejectEnforcementServiceRequest(contractID string, args RegistrarServiceRejectEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectEnforcementServiceRequestWithPackageID exercises the RegistrarService_RejectEnforcementServiceRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceRejectEnforcementServiceRequestWithPackageID(contractID string, packageID string, args RegistrarServiceRejectEnforcementServiceRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectEnforcementServiceRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceTerminateEnforcementService exercises the RegistrarService_TerminateEnforcementService choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceTerminateEnforcementService(contractID string, args RegistrarServiceTerminateEnforcementService) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_TerminateEnforcementService",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceTerminateEnforcementServiceWithPackageID exercises the RegistrarService_TerminateEnforcementService choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceTerminateEnforcementServiceWithPackageID(contractID string, packageID string, args RegistrarServiceTerminateEnforcementService) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_TerminateEnforcementService",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateRegistrarConfiguration exercises the RegistrarService_CreateRegistrarConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCreateRegistrarConfiguration(contractID string, args RegistrarServiceCreateRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateRegistrarConfigurationWithPackageID exercises the RegistrarService_CreateRegistrarConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCreateRegistrarConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceCreateRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptForceTransferRequest exercises the RegistrarService_AcceptForceTransferRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceAcceptForceTransferRequest(contractID string, args RegistrarServiceAcceptForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptForceTransferRequestWithPackageID exercises the RegistrarService_AcceptForceTransferRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceAcceptForceTransferRequestWithPackageID(contractID string, packageID string, args RegistrarServiceAcceptForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedForceTransfer exercises the RegistrarService_ExecuteAcceptedForceTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedForceTransfer(contractID string, args RegistrarServiceExecuteAcceptedForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedForceTransferWithPackageID exercises the RegistrarService_ExecuteAcceptedForceTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedForceTransferWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedForceTransfer exercises the RegistrarService_FailAcceptedForceTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedForceTransfer(contractID string, args RegistrarServiceFailAcceptedForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedForceTransferWithPackageID exercises the RegistrarService_FailAcceptedForceTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedForceTransferWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedForceTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedForceTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectForceTransferRequest exercises the RegistrarService_RejectForceTransferRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceRejectForceTransferRequest(contractID string, args RegistrarServiceRejectForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectForceTransferRequestWithPackageID exercises the RegistrarService_RejectForceTransferRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceRejectForceTransferRequestWithPackageID(contractID string, packageID string, args RegistrarServiceRejectForceTransferRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectForceTransferRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveRegistrarConfiguration exercises the RegistrarService_ArchiveRegistrarConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveRegistrarConfiguration(contractID string, args RegistrarServiceArchiveRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveRegistrarConfigurationWithPackageID exercises the RegistrarService_ArchiveRegistrarConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveRegistrarConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAndCreateRegistrarConfiguration exercises the RegistrarService_ArchiveAndCreateRegistrarConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveAndCreateRegistrarConfiguration(contractID string, args RegistrarServiceArchiveAndCreateRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAndCreateRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAndCreateRegistrarConfigurationWithPackageID exercises the RegistrarService_ArchiveAndCreateRegistrarConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveAndCreateRegistrarConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveAndCreateRegistrarConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAndCreateRegistrarConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateInstrumentConfiguration exercises the RegistrarService_CreateInstrumentConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCreateInstrumentConfiguration(contractID string, args RegistrarServiceCreateInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateInstrumentConfigurationWithPackageID exercises the RegistrarService_CreateInstrumentConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCreateInstrumentConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceCreateInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveInstrumentConfiguration exercises the RegistrarService_ArchiveInstrumentConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveInstrumentConfiguration(contractID string, args RegistrarServiceArchiveInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveInstrumentConfigurationWithPackageID exercises the RegistrarService_ArchiveInstrumentConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveInstrumentConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAndCreateInstrumentConfiguration exercises the RegistrarService_ArchiveAndCreateInstrumentConfiguration choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveAndCreateInstrumentConfiguration(contractID string, args RegistrarServiceArchiveAndCreateInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAndCreateInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAndCreateInstrumentConfigurationWithPackageID exercises the RegistrarService_ArchiveAndCreateInstrumentConfiguration choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveAndCreateInstrumentConfigurationWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveAndCreateInstrumentConfiguration) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAndCreateInstrumentConfiguration",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceOfferMint exercises the RegistrarService_OfferMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceOfferMint(contractID string, args RegistrarServiceOfferMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_OfferMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceOfferMintWithPackageID exercises the RegistrarService_OfferMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceOfferMintWithPackageID(contractID string, packageID string, args RegistrarServiceOfferMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_OfferMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCancelMintOffer exercises the RegistrarService_CancelMintOffer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCancelMintOffer(contractID string, args RegistrarServiceCancelMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CancelMintOffer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCancelMintOfferWithPackageID exercises the RegistrarService_CancelMintOffer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCancelMintOfferWithPackageID(contractID string, packageID string, args RegistrarServiceCancelMintOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CancelMintOffer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptMintRequest exercises the RegistrarService_AcceptMintRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceAcceptMintRequest(contractID string, args RegistrarServiceAcceptMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptMintRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptMintRequestWithPackageID exercises the RegistrarService_AcceptMintRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceAcceptMintRequestWithPackageID(contractID string, packageID string, args RegistrarServiceAcceptMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptMintRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectMintRequest exercises the RegistrarService_RejectMintRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceRejectMintRequest(contractID string, args RegistrarServiceRejectMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectMintRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectMintRequestWithPackageID exercises the RegistrarService_RejectMintRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceRejectMintRequestWithPackageID(contractID string, packageID string, args RegistrarServiceRejectMintRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectMintRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceOfferBurn exercises the RegistrarService_OfferBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceOfferBurn(contractID string, args RegistrarServiceOfferBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_OfferBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceOfferBurnWithPackageID exercises the RegistrarService_OfferBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceOfferBurnWithPackageID(contractID string, packageID string, args RegistrarServiceOfferBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_OfferBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCancelBurnOffer exercises the RegistrarService_CancelBurnOffer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCancelBurnOffer(contractID string, args RegistrarServiceCancelBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CancelBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCancelBurnOfferWithPackageID exercises the RegistrarService_CancelBurnOffer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCancelBurnOfferWithPackageID(contractID string, packageID string, args RegistrarServiceCancelBurnOffer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CancelBurnOffer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptBurnRequest exercises the RegistrarService_AcceptBurnRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceAcceptBurnRequest(contractID string, args RegistrarServiceAcceptBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceAcceptBurnRequestWithPackageID exercises the RegistrarService_AcceptBurnRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceAcceptBurnRequestWithPackageID(contractID string, packageID string, args RegistrarServiceAcceptBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_AcceptBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectBurnRequest exercises the RegistrarService_RejectBurnRequest choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceRejectBurnRequest(contractID string, args RegistrarServiceRejectBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRejectBurnRequestWithPackageID exercises the RegistrarService_RejectBurnRequest choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceRejectBurnRequestWithPackageID(contractID string, packageID string, args RegistrarServiceRejectBurnRequest) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_RejectBurnRequest",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceSplitHolding exercises the RegistrarService_SplitHolding choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceSplitHolding(contractID string, args RegistrarServiceSplitHolding) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_SplitHolding",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceSplitHoldingWithPackageID exercises the RegistrarService_SplitHolding choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceSplitHoldingWithPackageID(contractID string, packageID string, args RegistrarServiceSplitHolding) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_SplitHolding",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceMergeHolding exercises the RegistrarService_MergeHolding choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceMergeHolding(contractID string, args RegistrarServiceMergeHolding) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_MergeHolding",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceMergeHoldingWithPackageID exercises the RegistrarService_MergeHolding choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceMergeHoldingWithPackageID(contractID string, packageID string, args RegistrarServiceMergeHolding) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_MergeHolding",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedTransfer exercises the RegistrarService_ExecuteAcceptedTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedTransfer(contractID string, args RegistrarServiceExecuteAcceptedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedTransferWithPackageID exercises the RegistrarService_ExecuteAcceptedTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedTransferWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedTransfer exercises the RegistrarService_FailAcceptedTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedTransfer(contractID string, args RegistrarServiceFailAcceptedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedTransferWithPackageID exercises the RegistrarService_FailAcceptedTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedTransferWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedLock exercises the RegistrarService_ExecuteAcceptedLock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedLock(contractID string, args RegistrarServiceExecuteAcceptedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedLockWithPackageID exercises the RegistrarService_ExecuteAcceptedLock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedLockWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedLock exercises the RegistrarService_FailAcceptedLock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedLock(contractID string, args RegistrarServiceFailAcceptedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedLockWithPackageID exercises the RegistrarService_FailAcceptedLock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedLockWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedUnlock exercises the RegistrarService_ExecuteAcceptedUnlock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedUnlock(contractID string, args RegistrarServiceExecuteAcceptedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedUnlockWithPackageID exercises the RegistrarService_ExecuteAcceptedUnlock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedUnlockWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedUnlock exercises the RegistrarService_FailAcceptedUnlock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedUnlock(contractID string, args RegistrarServiceFailAcceptedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedUnlockWithPackageID exercises the RegistrarService_FailAcceptedUnlock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedUnlockWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedMint exercises the RegistrarService_ExecuteAcceptedMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedMint(contractID string, args RegistrarServiceExecuteAcceptedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedMintWithPackageID exercises the RegistrarService_ExecuteAcceptedMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedMintWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedMint exercises the RegistrarService_FailAcceptedMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedMint(contractID string, args RegistrarServiceFailAcceptedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedMintWithPackageID exercises the RegistrarService_FailAcceptedMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedMintWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedBurn exercises the RegistrarService_ExecuteAcceptedBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceExecuteAcceptedBurn(contractID string, args RegistrarServiceExecuteAcceptedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceExecuteAcceptedBurnWithPackageID exercises the RegistrarService_ExecuteAcceptedBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceExecuteAcceptedBurnWithPackageID(contractID string, packageID string, args RegistrarServiceExecuteAcceptedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ExecuteAcceptedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedBurn exercises the RegistrarService_FailAcceptedBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceFailAcceptedBurn(contractID string, args RegistrarServiceFailAcceptedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceFailAcceptedBurnWithPackageID exercises the RegistrarService_FailAcceptedBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceFailAcceptedBurnWithPackageID(contractID string, packageID string, args RegistrarServiceFailAcceptedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_FailAcceptedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedMint exercises the RegistrarService_DeleteRejectedMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedMint(contractID string, args RegistrarServiceDeleteRejectedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedMintWithPackageID exercises the RegistrarService_DeleteRejectedMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedMintWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedMint exercises the RegistrarService_DeleteFailedMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedMint(contractID string, args RegistrarServiceDeleteFailedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedMintWithPackageID exercises the RegistrarService_DeleteFailedMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedMintWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedMint exercises the RegistrarService_DeleteExecutedMint choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedMint(contractID string, args RegistrarServiceDeleteExecutedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedMintWithPackageID exercises the RegistrarService_DeleteExecutedMint choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedMintWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedMint",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedBurn exercises the RegistrarService_DeleteRejectedBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedBurn(contractID string, args RegistrarServiceDeleteRejectedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedBurnWithPackageID exercises the RegistrarService_DeleteRejectedBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedBurnWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedBurn exercises the RegistrarService_DeleteFailedBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedBurn(contractID string, args RegistrarServiceDeleteFailedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedBurnWithPackageID exercises the RegistrarService_DeleteFailedBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedBurnWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedBurn exercises the RegistrarService_DeleteExecutedBurn choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedBurn(contractID string, args RegistrarServiceDeleteExecutedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedBurnWithPackageID exercises the RegistrarService_DeleteExecutedBurn choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedBurnWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedBurn",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedLock exercises the RegistrarService_DeleteRejectedLock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedLock(contractID string, args RegistrarServiceDeleteRejectedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedLockWithPackageID exercises the RegistrarService_DeleteRejectedLock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedLockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedLock exercises the RegistrarService_DeleteFailedLock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedLock(contractID string, args RegistrarServiceDeleteFailedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedLockWithPackageID exercises the RegistrarService_DeleteFailedLock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedLockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedLock exercises the RegistrarService_DeleteExecutedLock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedLock(contractID string, args RegistrarServiceDeleteExecutedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedLockWithPackageID exercises the RegistrarService_DeleteExecutedLock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedLockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedLock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedLock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedUnlock exercises the RegistrarService_DeleteRejectedUnlock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedUnlock(contractID string, args RegistrarServiceDeleteRejectedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedUnlockWithPackageID exercises the RegistrarService_DeleteRejectedUnlock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedUnlockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedUnlock exercises the RegistrarService_DeleteFailedUnlock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedUnlock(contractID string, args RegistrarServiceDeleteFailedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedUnlockWithPackageID exercises the RegistrarService_DeleteFailedUnlock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedUnlockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedUnlock exercises the RegistrarService_DeleteExecutedUnlock choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedUnlock(contractID string, args RegistrarServiceDeleteExecutedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedUnlockWithPackageID exercises the RegistrarService_DeleteExecutedUnlock choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedUnlockWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedUnlock) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedUnlock",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedTransfer exercises the RegistrarService_DeleteRejectedTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedTransfer(contractID string, args RegistrarServiceDeleteRejectedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedTransferWithPackageID exercises the RegistrarService_DeleteRejectedTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedTransferWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedTransfers exercises the RegistrarService_DeleteRejectedTransfers choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteRejectedTransfers(contractID string, args RegistrarServiceDeleteRejectedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteRejectedTransfersWithPackageID exercises the RegistrarService_DeleteRejectedTransfers choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteRejectedTransfersWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteRejectedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteRejectedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedTransfer exercises the RegistrarService_DeleteFailedTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedTransfer(contractID string, args RegistrarServiceDeleteFailedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedTransferWithPackageID exercises the RegistrarService_DeleteFailedTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedTransferWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedTransfers exercises the RegistrarService_DeleteFailedTransfers choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteFailedTransfers(contractID string, args RegistrarServiceDeleteFailedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteFailedTransfersWithPackageID exercises the RegistrarService_DeleteFailedTransfers choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteFailedTransfersWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteFailedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteFailedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedTransfer exercises the RegistrarService_DeleteExecutedTransfer choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedTransfer(contractID string, args RegistrarServiceDeleteExecutedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedTransferWithPackageID exercises the RegistrarService_DeleteExecutedTransfer choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedTransferWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedTransfer",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedTransfers exercises the RegistrarService_DeleteExecutedTransfers choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceDeleteExecutedTransfers(contractID string, args RegistrarServiceDeleteExecutedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceDeleteExecutedTransfersWithPackageID exercises the RegistrarService_DeleteExecutedTransfers choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceDeleteExecutedTransfersWithPackageID(contractID string, packageID string, args RegistrarServiceDeleteExecutedTransfers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_DeleteExecutedTransfers",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateAllocationFactory exercises the RegistrarService_CreateAllocationFactory choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCreateAllocationFactory(contractID string, args RegistrarServiceCreateAllocationFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateAllocationFactory",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateAllocationFactoryWithPackageID exercises the RegistrarService_CreateAllocationFactory choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCreateAllocationFactoryWithPackageID(contractID string, packageID string, args RegistrarServiceCreateAllocationFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateAllocationFactory",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAllocationFactory exercises the RegistrarService_ArchiveAllocationFactory choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveAllocationFactory(contractID string, args RegistrarServiceArchiveAllocationFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAllocationFactory",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveAllocationFactoryWithPackageID exercises the RegistrarService_ArchiveAllocationFactory choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveAllocationFactoryWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveAllocationFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveAllocationFactory",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateTransferRule exercises the RegistrarService_CreateTransferRule choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceCreateTransferRule(contractID string, args RegistrarServiceCreateTransferRule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateTransferRule",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceCreateTransferRuleWithPackageID exercises the RegistrarService_CreateTransferRule choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceCreateTransferRuleWithPackageID(contractID string, packageID string, args RegistrarServiceCreateTransferRule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_CreateTransferRule",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveTransferRule exercises the RegistrarService_ArchiveTransferRule choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) RegistrarServiceArchiveTransferRule(contractID string, args RegistrarServiceArchiveTransferRule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveTransferRule",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceArchiveTransferRuleWithPackageID exercises the RegistrarService_ArchiveTransferRule choice using the provided package ID instead of package name
func (t RegistrarService) RegistrarServiceArchiveTransferRuleWithPackageID(contractID string, packageID string, args RegistrarServiceArchiveTransferRule) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "RegistrarService_ArchiveTransferRule",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RegistrarService contract
// This method uses the package name in the template ID
func (t RegistrarService) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RegistrarService) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarService"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RegistrarServiceRequest is a Template type
type RegistrarServiceRequest struct {
	Operator                types.PARTY `json:"operator"`
	Provider                types.PARTY `json:"provider"`
	Registrar               types.PARTY `json:"registrar"`
	CreateTransferRule      *types.BOOL `json:"createTransferRule" hex:"optional"`
	CreateAllocationFactory *types.BOOL `json:"createAllocationFactory" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RegistrarServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RegistrarServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RegistrarServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	if t.CreateTransferRule != nil {
		args["createTransferRule"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.CreateTransferRule),
		}
	} else {
		args["createTransferRule"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.CreateAllocationFactory != nil {
		args["createAllocationFactory"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.CreateAllocationFactory),
		}
	} else {
		args["createAllocationFactory"] = map[string]any{
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
func (t RegistrarServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registrar"] = t.Registrar.ToMap()

	if t.CreateTransferRule != nil {
		args["createTransferRule"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.CreateTransferRule),
		}
	} else {
		args["createTransferRule"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.CreateAllocationFactory != nil {
		args["createAllocationFactory"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.CreateAllocationFactory),
		}
	} else {
		args["createAllocationFactory"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RegistrarServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequest to hex string (Canton MCMS format)
func (t RegistrarServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RegistrarServiceRequest

// RegistrarServiceRequestAccept exercises the RegistrarServiceRequest_Accept choice on this RegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RegistrarServiceRequest) RegistrarServiceRequestAccept(contractID string, args RegistrarServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRequestAcceptWithPackageID exercises the RegistrarServiceRequest_Accept choice using the provided package ID instead of package name
func (t RegistrarServiceRequest) RegistrarServiceRequestAcceptWithPackageID(contractID string, packageID string, args RegistrarServiceRequestAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Accept",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRequestReject exercises the RegistrarServiceRequest_Reject choice on this RegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RegistrarServiceRequest) RegistrarServiceRequestReject(contractID string, args RegistrarServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRequestRejectWithPackageID exercises the RegistrarServiceRequest_Reject choice using the provided package ID instead of package name
func (t RegistrarServiceRequest) RegistrarServiceRequestRejectWithPackageID(contractID string, packageID string, args RegistrarServiceRequestReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Reject",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRequestCancel exercises the RegistrarServiceRequest_Cancel choice on this RegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RegistrarServiceRequest) RegistrarServiceRequestCancel(contractID string, args RegistrarServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// RegistrarServiceRequestCancelWithPackageID exercises the RegistrarServiceRequest_Cancel choice using the provided package ID instead of package name
func (t RegistrarServiceRequest) RegistrarServiceRequestCancelWithPackageID(contractID string, packageID string, args RegistrarServiceRequestCancel) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RegistrarServiceRequest_Cancel",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RegistrarServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RegistrarServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RegistrarServiceRequestAccept is a Record type
type RegistrarServiceRequestAccept struct {
	ProviderConfigurationCid types.CONTRACT_ID   `json:"providerConfigurationCid"`
	CredentialCids           []types.CONTRACT_ID `json:"credentialCids"`
}

// ToMap converts RegistrarServiceRequestAccept to a map for DAML arguments
func (t RegistrarServiceRequestAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["providerConfigurationCid"] = model.NestedToDAMLValue(t.ProviderConfigurationCid)

	m["credentialCids"] = func() []any {
		res := make([]any, 0, len(t.CredentialCids))
		for _, e := range t.CredentialCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t RegistrarServiceRequestAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestAccept to hex string (Canton MCMS format)
func (t RegistrarServiceRequestAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestAccept from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRequestAcceptResult is a Record type
type RegistrarServiceRequestAcceptResult struct {
	RegistrarServiceCid  types.CONTRACT_ID  `json:"registrarServiceCid"`
	TransferRuleCid      *types.CONTRACT_ID `json:"transferRuleCid" hex:"optional"`
	AllocationFactoryCid *types.CONTRACT_ID `json:"allocationFactoryCid" hex:"optional"`
}

// ToMap converts RegistrarServiceRequestAcceptResult to a map for DAML arguments
func (t RegistrarServiceRequestAcceptResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarServiceCid"] = model.NestedToDAMLValue(t.RegistrarServiceCid)

	if t.TransferRuleCid != nil {
		m["transferRuleCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TransferRuleCid),
		}
	} else {
		m["transferRuleCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.AllocationFactoryCid != nil {
		m["allocationFactoryCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.AllocationFactoryCid),
		}
	} else {
		m["allocationFactoryCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t RegistrarServiceRequestAcceptResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestAcceptResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestAcceptResult to hex string (Canton MCMS format)
func (t RegistrarServiceRequestAcceptResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestAcceptResult from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestAcceptResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRequestCancel is a Record type
type RegistrarServiceRequestCancel struct {
}

// ToMap converts RegistrarServiceRequestCancel to a map for DAML arguments
func (t RegistrarServiceRequestCancel) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceRequestCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestCancel to hex string (Canton MCMS format)
func (t RegistrarServiceRequestCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestCancel from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRequestCancelResult is a Record type
type RegistrarServiceRequestCancelResult struct {
}

// ToMap converts RegistrarServiceRequestCancelResult to a map for DAML arguments
func (t RegistrarServiceRequestCancelResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceRequestCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestCancelResult to hex string (Canton MCMS format)
func (t RegistrarServiceRequestCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestCancelResult from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRequestReject is a Record type
type RegistrarServiceRequestReject struct {
	Reason types.TEXT `json:"reason"`
}

// ToMap converts RegistrarServiceRequestReject to a map for DAML arguments
func (t RegistrarServiceRequestReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["reason"] = string(t.Reason)

	return m
}

func (t RegistrarServiceRequestReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestReject to hex string (Canton MCMS format)
func (t RegistrarServiceRequestReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestReject from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRequestRejectResult is a Record type
type RegistrarServiceRequestRejectResult struct {
	RejectedRegistrarServiceRequestCid types.CONTRACT_ID `json:"rejectedRegistrarServiceRequestCid"`
}

// ToMap converts RegistrarServiceRequestRejectResult to a map for DAML arguments
func (t RegistrarServiceRequestRejectResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rejectedRegistrarServiceRequestCid"] = model.NestedToDAMLValue(t.RejectedRegistrarServiceRequestCid)

	return m
}

func (t RegistrarServiceRequestRejectResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRequestRejectResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRequestRejectResult to hex string (Canton MCMS format)
func (t RegistrarServiceRequestRejectResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRequestRejectResult from hex string (Canton MCMS format)
func (t *RegistrarServiceRequestRejectResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceAcceptBurnRequest is a Record type
type RegistrarServiceAcceptBurnRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.BurnRequestAccept `json:"payload"`
}

// ToMap converts RegistrarServiceAcceptBurnRequest to a map for DAML arguments
func (t RegistrarServiceAcceptBurnRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceAcceptBurnRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceAcceptBurnRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceAcceptBurnRequest to hex string (Canton MCMS format)
func (t RegistrarServiceAcceptBurnRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceAcceptBurnRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceAcceptBurnRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceAcceptEnforcementServiceRequest is a Record type
type RegistrarServiceAcceptEnforcementServiceRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload EnforcementServiceRequestAccept `json:"payload"`
}

// ToMap converts RegistrarServiceAcceptEnforcementServiceRequest to a map for DAML arguments
func (t RegistrarServiceAcceptEnforcementServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceAcceptEnforcementServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceAcceptEnforcementServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceAcceptEnforcementServiceRequest to hex string (Canton MCMS format)
func (t RegistrarServiceAcceptEnforcementServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceAcceptEnforcementServiceRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceAcceptEnforcementServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceAcceptForceTransferRequest is a Record type
type RegistrarServiceAcceptForceTransferRequest struct {
	SenderEnforcementServiceCid   types.CONTRACT_ID                      `json:"senderEnforcementServiceCid"`
	ReceiverEnforcementServiceCid types.CONTRACT_ID                      `json:"receiverEnforcementServiceCid"`
	Cid                           types.CONTRACT_ID                      `json:"cid"`
	Payload                       registry_v0.ForceTransferRequestAccept `json:"payload"`
}

// ToMap converts RegistrarServiceAcceptForceTransferRequest to a map for DAML arguments
func (t RegistrarServiceAcceptForceTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderEnforcementServiceCid"] = model.NestedToDAMLValue(t.SenderEnforcementServiceCid)

	m["receiverEnforcementServiceCid"] = model.NestedToDAMLValue(t.ReceiverEnforcementServiceCid)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceAcceptForceTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceAcceptForceTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceAcceptForceTransferRequest to hex string (Canton MCMS format)
func (t RegistrarServiceAcceptForceTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceAcceptForceTransferRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceAcceptForceTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceAcceptMintRequest is a Record type
type RegistrarServiceAcceptMintRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.MintRequestAccept `json:"payload"`
}

// ToMap converts RegistrarServiceAcceptMintRequest to a map for DAML arguments
func (t RegistrarServiceAcceptMintRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceAcceptMintRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceAcceptMintRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceAcceptMintRequest to hex string (Canton MCMS format)
func (t RegistrarServiceAcceptMintRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceAcceptMintRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceAcceptMintRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAllocationFactory is a Record type
type RegistrarServiceArchiveAllocationFactory struct {
	AllocationFactoryCid types.CONTRACT_ID `json:"allocationFactoryCid"`
}

// ToMap converts RegistrarServiceArchiveAllocationFactory to a map for DAML arguments
func (t RegistrarServiceArchiveAllocationFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocationFactoryCid"] = model.NestedToDAMLValue(t.AllocationFactoryCid)

	return m
}

func (t RegistrarServiceArchiveAllocationFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAllocationFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAllocationFactory to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAllocationFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAllocationFactory from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAllocationFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAllocationFactoryResult is a Record type
type RegistrarServiceArchiveAllocationFactoryResult struct {
}

// ToMap converts RegistrarServiceArchiveAllocationFactoryResult to a map for DAML arguments
func (t RegistrarServiceArchiveAllocationFactoryResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceArchiveAllocationFactoryResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAllocationFactoryResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAllocationFactoryResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAllocationFactoryResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAllocationFactoryResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAllocationFactoryResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAndCreateInstrumentConfiguration is a Record type
type RegistrarServiceArchiveAndCreateInstrumentConfiguration struct {
	Cid     types.CONTRACT_ID                             `json:"cid"`
	Payload RegistrarServiceCreateInstrumentConfiguration `json:"payload"`
}

// ToMap converts RegistrarServiceArchiveAndCreateInstrumentConfiguration to a map for DAML arguments
func (t RegistrarServiceArchiveAndCreateInstrumentConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceArchiveAndCreateInstrumentConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAndCreateInstrumentConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAndCreateInstrumentConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAndCreateInstrumentConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAndCreateInstrumentConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAndCreateInstrumentConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAndCreateInstrumentConfigurationResult is a Record type
type RegistrarServiceArchiveAndCreateInstrumentConfigurationResult struct {
	ArchiveResult RegistrarServiceArchiveInstrumentConfigurationResult `json:"archiveResult"`
	CreateResult  RegistrarServiceCreateInstrumentConfigurationResult  `json:"createResult"`
}

// ToMap converts RegistrarServiceArchiveAndCreateInstrumentConfigurationResult to a map for DAML arguments
func (t RegistrarServiceArchiveAndCreateInstrumentConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["archiveResult"] = model.NestedToDAMLValue(t.ArchiveResult)

	m["createResult"] = model.NestedToDAMLValue(t.CreateResult)

	return m
}

func (t RegistrarServiceArchiveAndCreateInstrumentConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAndCreateInstrumentConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAndCreateInstrumentConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAndCreateInstrumentConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAndCreateInstrumentConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAndCreateInstrumentConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAndCreateRegistrarConfiguration is a Record type
type RegistrarServiceArchiveAndCreateRegistrarConfiguration struct {
	Cid     types.CONTRACT_ID                            `json:"cid"`
	Payload RegistrarServiceCreateRegistrarConfiguration `json:"payload"`
}

// ToMap converts RegistrarServiceArchiveAndCreateRegistrarConfiguration to a map for DAML arguments
func (t RegistrarServiceArchiveAndCreateRegistrarConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceArchiveAndCreateRegistrarConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAndCreateRegistrarConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAndCreateRegistrarConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAndCreateRegistrarConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAndCreateRegistrarConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAndCreateRegistrarConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveAndCreateRegistrarConfigurationResult is a Record type
type RegistrarServiceArchiveAndCreateRegistrarConfigurationResult struct {
	ArchiveResult RegistrarServiceArchiveRegistrarConfigurationResult `json:"archiveResult"`
	CreateResult  RegistrarServiceCreateRegistrarConfigurationResult  `json:"createResult"`
}

// ToMap converts RegistrarServiceArchiveAndCreateRegistrarConfigurationResult to a map for DAML arguments
func (t RegistrarServiceArchiveAndCreateRegistrarConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["archiveResult"] = model.NestedToDAMLValue(t.ArchiveResult)

	m["createResult"] = model.NestedToDAMLValue(t.CreateResult)

	return m
}

func (t RegistrarServiceArchiveAndCreateRegistrarConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveAndCreateRegistrarConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveAndCreateRegistrarConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveAndCreateRegistrarConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveAndCreateRegistrarConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveAndCreateRegistrarConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveInstrumentConfiguration is a Record type
type RegistrarServiceArchiveInstrumentConfiguration struct {
	Cid types.CONTRACT_ID `json:"cid"`
}

// ToMap converts RegistrarServiceArchiveInstrumentConfiguration to a map for DAML arguments
func (t RegistrarServiceArchiveInstrumentConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	return m
}

func (t RegistrarServiceArchiveInstrumentConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveInstrumentConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveInstrumentConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveInstrumentConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveInstrumentConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveInstrumentConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveInstrumentConfigurationResult is a Record type
type RegistrarServiceArchiveInstrumentConfigurationResult struct {
}

// ToMap converts RegistrarServiceArchiveInstrumentConfigurationResult to a map for DAML arguments
func (t RegistrarServiceArchiveInstrumentConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceArchiveInstrumentConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveInstrumentConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveInstrumentConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveInstrumentConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveInstrumentConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveInstrumentConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveRegistrarConfiguration is a Record type
type RegistrarServiceArchiveRegistrarConfiguration struct {
	Cid types.CONTRACT_ID `json:"cid"`
}

// ToMap converts RegistrarServiceArchiveRegistrarConfiguration to a map for DAML arguments
func (t RegistrarServiceArchiveRegistrarConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	return m
}

func (t RegistrarServiceArchiveRegistrarConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveRegistrarConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveRegistrarConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveRegistrarConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveRegistrarConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveRegistrarConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveRegistrarConfigurationResult is a Record type
type RegistrarServiceArchiveRegistrarConfigurationResult struct {
}

// ToMap converts RegistrarServiceArchiveRegistrarConfigurationResult to a map for DAML arguments
func (t RegistrarServiceArchiveRegistrarConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceArchiveRegistrarConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveRegistrarConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveRegistrarConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveRegistrarConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveRegistrarConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveRegistrarConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveTransferRule is a Record type
type RegistrarServiceArchiveTransferRule struct {
	Cid types.CONTRACT_ID `json:"cid"`
}

// ToMap converts RegistrarServiceArchiveTransferRule to a map for DAML arguments
func (t RegistrarServiceArchiveTransferRule) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	return m
}

func (t RegistrarServiceArchiveTransferRule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveTransferRule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveTransferRule to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveTransferRule) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveTransferRule from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveTransferRule) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceArchiveTransferRuleResult is a Record type
type RegistrarServiceArchiveTransferRuleResult struct {
}

// ToMap converts RegistrarServiceArchiveTransferRuleResult to a map for DAML arguments
func (t RegistrarServiceArchiveTransferRuleResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceArchiveTransferRuleResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceArchiveTransferRuleResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceArchiveTransferRuleResult to hex string (Canton MCMS format)
func (t RegistrarServiceArchiveTransferRuleResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceArchiveTransferRuleResult from hex string (Canton MCMS format)
func (t *RegistrarServiceArchiveTransferRuleResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCancelBurnOffer is a Record type
type RegistrarServiceCancelBurnOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.BurnOfferCancel `json:"payload"`
}

// ToMap converts RegistrarServiceCancelBurnOffer to a map for DAML arguments
func (t RegistrarServiceCancelBurnOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceCancelBurnOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCancelBurnOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCancelBurnOffer to hex string (Canton MCMS format)
func (t RegistrarServiceCancelBurnOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCancelBurnOffer from hex string (Canton MCMS format)
func (t *RegistrarServiceCancelBurnOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCancelMintOffer is a Record type
type RegistrarServiceCancelMintOffer struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload registry_v0.MintOfferCancel `json:"payload"`
}

// ToMap converts RegistrarServiceCancelMintOffer to a map for DAML arguments
func (t RegistrarServiceCancelMintOffer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceCancelMintOffer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCancelMintOffer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCancelMintOffer to hex string (Canton MCMS format)
func (t RegistrarServiceCancelMintOffer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCancelMintOffer from hex string (Canton MCMS format)
func (t *RegistrarServiceCancelMintOffer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateAllocationFactory is a Record type
type RegistrarServiceCreateAllocationFactory struct {
}

// ToMap converts RegistrarServiceCreateAllocationFactory to a map for DAML arguments
func (t RegistrarServiceCreateAllocationFactory) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceCreateAllocationFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateAllocationFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateAllocationFactory to hex string (Canton MCMS format)
func (t RegistrarServiceCreateAllocationFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateAllocationFactory from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateAllocationFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateAllocationFactoryResult is a Record type
type RegistrarServiceCreateAllocationFactoryResult struct {
	AllocationFactoryCid types.CONTRACT_ID `json:"allocationFactoryCid"`
}

// ToMap converts RegistrarServiceCreateAllocationFactoryResult to a map for DAML arguments
func (t RegistrarServiceCreateAllocationFactoryResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocationFactoryCid"] = model.NestedToDAMLValue(t.AllocationFactoryCid)

	return m
}

func (t RegistrarServiceCreateAllocationFactoryResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateAllocationFactoryResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateAllocationFactoryResult to hex string (Canton MCMS format)
func (t RegistrarServiceCreateAllocationFactoryResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateAllocationFactoryResult from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateAllocationFactoryResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateInstrumentConfiguration is a Record type
type RegistrarServiceCreateInstrumentConfiguration struct {
	InstrumentId          types.TEXT                                 `json:"instrumentId"`
	AdditionalIdentifiers []registry_holding_v0.InstrumentIdentifier `json:"additionalIdentifiers"`
	IssuerRequirements    []credential_v0.PartyCredentialRequirement `json:"issuerRequirements"`
	HolderRequirements    []credential_v0.PartyCredentialRequirement `json:"holderRequirements"`
}

// ToMap converts RegistrarServiceCreateInstrumentConfiguration to a map for DAML arguments
func (t RegistrarServiceCreateInstrumentConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentId"] = string(t.InstrumentId)

	m["additionalIdentifiers"] = func() []any {
		res := make([]any, 0, len(t.AdditionalIdentifiers))
		for _, e := range t.AdditionalIdentifiers {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["issuerRequirements"] = func() []any {
		res := make([]any, 0, len(t.IssuerRequirements))
		for _, e := range t.IssuerRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["holderRequirements"] = func() []any {
		res := make([]any, 0, len(t.HolderRequirements))
		for _, e := range t.HolderRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RegistrarServiceCreateInstrumentConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateInstrumentConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateInstrumentConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceCreateInstrumentConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateInstrumentConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateInstrumentConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateInstrumentConfigurationResult is a Record type
type RegistrarServiceCreateInstrumentConfigurationResult struct {
	InstrumentConfigurationCid types.CONTRACT_ID `json:"instrumentConfigurationCid"`
}

// ToMap converts RegistrarServiceCreateInstrumentConfigurationResult to a map for DAML arguments
func (t RegistrarServiceCreateInstrumentConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentConfigurationCid"] = model.NestedToDAMLValue(t.InstrumentConfigurationCid)

	return m
}

func (t RegistrarServiceCreateInstrumentConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateInstrumentConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateInstrumentConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceCreateInstrumentConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateInstrumentConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateInstrumentConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateRegistrarConfiguration is a Record type
type RegistrarServiceCreateRegistrarConfiguration struct {
	EnforcementRequirements []credential_v0.PartyCredentialRequirement `json:"enforcementRequirements"`
}

// ToMap converts RegistrarServiceCreateRegistrarConfiguration to a map for DAML arguments
func (t RegistrarServiceCreateRegistrarConfiguration) ToMap() map[string]any {
	m := make(map[string]any)

	m["enforcementRequirements"] = func() []any {
		res := make([]any, 0, len(t.EnforcementRequirements))
		for _, e := range t.EnforcementRequirements {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t RegistrarServiceCreateRegistrarConfiguration) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateRegistrarConfiguration) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateRegistrarConfiguration to hex string (Canton MCMS format)
func (t RegistrarServiceCreateRegistrarConfiguration) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateRegistrarConfiguration from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateRegistrarConfiguration) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateRegistrarConfigurationResult is a Record type
type RegistrarServiceCreateRegistrarConfigurationResult struct {
	RegistrarConfigurationCid types.CONTRACT_ID `json:"registrarConfigurationCid"`
}

// ToMap converts RegistrarServiceCreateRegistrarConfigurationResult to a map for DAML arguments
func (t RegistrarServiceCreateRegistrarConfigurationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["registrarConfigurationCid"] = model.NestedToDAMLValue(t.RegistrarConfigurationCid)

	return m
}

func (t RegistrarServiceCreateRegistrarConfigurationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateRegistrarConfigurationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateRegistrarConfigurationResult to hex string (Canton MCMS format)
func (t RegistrarServiceCreateRegistrarConfigurationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateRegistrarConfigurationResult from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateRegistrarConfigurationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateTransferRule is a Record type
type RegistrarServiceCreateTransferRule struct {
}

// ToMap converts RegistrarServiceCreateTransferRule to a map for DAML arguments
func (t RegistrarServiceCreateTransferRule) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceCreateTransferRule) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateTransferRule) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateTransferRule to hex string (Canton MCMS format)
func (t RegistrarServiceCreateTransferRule) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateTransferRule from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateTransferRule) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceCreateTransferRuleResult is a Record type
type RegistrarServiceCreateTransferRuleResult struct {
	TransferRuleCid types.CONTRACT_ID `json:"transferRuleCid"`
}

// ToMap converts RegistrarServiceCreateTransferRuleResult to a map for DAML arguments
func (t RegistrarServiceCreateTransferRuleResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferRuleCid"] = model.NestedToDAMLValue(t.TransferRuleCid)

	return m
}

func (t RegistrarServiceCreateTransferRuleResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceCreateTransferRuleResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceCreateTransferRuleResult to hex string (Canton MCMS format)
func (t RegistrarServiceCreateTransferRuleResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceCreateTransferRuleResult from hex string (Canton MCMS format)
func (t *RegistrarServiceCreateTransferRuleResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedBurn is a Record type
type RegistrarServiceDeleteExecutedBurn struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.ExecutedBurnDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteExecutedBurn to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteExecutedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedBurn to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedLock is a Record type
type RegistrarServiceDeleteExecutedLock struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.ExecutedLockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteExecutedLock to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteExecutedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedLock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedLock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedMint is a Record type
type RegistrarServiceDeleteExecutedMint struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.ExecutedMintDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteExecutedMint to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteExecutedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedMint to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedMint from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedTransfer is a Record type
type RegistrarServiceDeleteExecutedTransfer struct {
	Cid     types.CONTRACT_ID                  `json:"cid"`
	Payload registry_v0.ExecutedTransferDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteExecutedTransfer to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteExecutedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedTransfers is a Record type
type RegistrarServiceDeleteExecutedTransfers struct {
	Cids            []types.CONTRACT_ID `json:"cids"`
	ChoiceObservers []types.PARTY       `json:"choiceObservers"`
	Actor           types.PARTY         `json:"actor"`
}

// ToMap converts RegistrarServiceDeleteExecutedTransfers to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedTransfers) ToMap() map[string]any {
	m := make(map[string]any)

	m["cids"] = func() []any {
		res := make([]any, 0, len(t.Cids))
		for _, e := range t.Cids {
			res = append(res, e)
		}
		return res
	}()

	m["choiceObservers"] = func() []any {
		res := make([]any, 0, len(t.ChoiceObservers))
		for _, e := range t.ChoiceObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarServiceDeleteExecutedTransfers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedTransfers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedTransfers to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedTransfers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedTransfers from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedTransfers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteExecutedUnlock is a Record type
type RegistrarServiceDeleteExecutedUnlock struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_v0.ExecutedUnlockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteExecutedUnlock to a map for DAML arguments
func (t RegistrarServiceDeleteExecutedUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteExecutedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteExecutedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteExecutedUnlock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteExecutedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteExecutedUnlock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteExecutedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedBurn is a Record type
type RegistrarServiceDeleteFailedBurn struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.FailedBurnDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteFailedBurn to a map for DAML arguments
func (t RegistrarServiceDeleteFailedBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteFailedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedBurn to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedLock is a Record type
type RegistrarServiceDeleteFailedLock struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.FailedLockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteFailedLock to a map for DAML arguments
func (t RegistrarServiceDeleteFailedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteFailedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedLock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedLock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedMint is a Record type
type RegistrarServiceDeleteFailedMint struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.FailedMintDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteFailedMint to a map for DAML arguments
func (t RegistrarServiceDeleteFailedMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteFailedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedMint to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedMint from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedTransfer is a Record type
type RegistrarServiceDeleteFailedTransfer struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_v0.FailedTransferDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteFailedTransfer to a map for DAML arguments
func (t RegistrarServiceDeleteFailedTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteFailedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedTransfers is a Record type
type RegistrarServiceDeleteFailedTransfers struct {
	Cids            []types.CONTRACT_ID `json:"cids"`
	ChoiceObservers []types.PARTY       `json:"choiceObservers"`
	Actor           types.PARTY         `json:"actor"`
}

// ToMap converts RegistrarServiceDeleteFailedTransfers to a map for DAML arguments
func (t RegistrarServiceDeleteFailedTransfers) ToMap() map[string]any {
	m := make(map[string]any)

	m["cids"] = func() []any {
		res := make([]any, 0, len(t.Cids))
		for _, e := range t.Cids {
			res = append(res, e)
		}
		return res
	}()

	m["choiceObservers"] = func() []any {
		res := make([]any, 0, len(t.ChoiceObservers))
		for _, e := range t.ChoiceObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarServiceDeleteFailedTransfers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedTransfers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedTransfers to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedTransfers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedTransfers from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedTransfers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteFailedUnlock is a Record type
type RegistrarServiceDeleteFailedUnlock struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.FailedUnlockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteFailedUnlock to a map for DAML arguments
func (t RegistrarServiceDeleteFailedUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteFailedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteFailedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteFailedUnlock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteFailedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteFailedUnlock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteFailedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedBurn is a Record type
type RegistrarServiceDeleteRejectedBurn struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.RejectedBurnDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteRejectedBurn to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteRejectedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedBurn to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedLock is a Record type
type RegistrarServiceDeleteRejectedLock struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.RejectedLockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteRejectedLock to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteRejectedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedLock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedLock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedMint is a Record type
type RegistrarServiceDeleteRejectedMint struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.RejectedMintDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteRejectedMint to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteRejectedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedMint to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedMint from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedTransfer is a Record type
type RegistrarServiceDeleteRejectedTransfer struct {
	Cid     types.CONTRACT_ID                  `json:"cid"`
	Payload registry_v0.RejectedTransferDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteRejectedTransfer to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteRejectedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedTransfers is a Record type
type RegistrarServiceDeleteRejectedTransfers struct {
	Cids            []types.CONTRACT_ID `json:"cids"`
	ChoiceObservers []types.PARTY       `json:"choiceObservers"`
	Actor           types.PARTY         `json:"actor"`
}

// ToMap converts RegistrarServiceDeleteRejectedTransfers to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedTransfers) ToMap() map[string]any {
	m := make(map[string]any)

	m["cids"] = func() []any {
		res := make([]any, 0, len(t.Cids))
		for _, e := range t.Cids {
			res = append(res, e)
		}
		return res
	}()

	m["choiceObservers"] = func() []any {
		res := make([]any, 0, len(t.ChoiceObservers))
		for _, e := range t.ChoiceObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarServiceDeleteRejectedTransfers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedTransfers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedTransfers to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedTransfers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedTransfers from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedTransfers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceDeleteRejectedUnlock is a Record type
type RegistrarServiceDeleteRejectedUnlock struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_v0.RejectedUnlockDelete `json:"payload"`
}

// ToMap converts RegistrarServiceDeleteRejectedUnlock to a map for DAML arguments
func (t RegistrarServiceDeleteRejectedUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceDeleteRejectedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceDeleteRejectedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceDeleteRejectedUnlock to hex string (Canton MCMS format)
func (t RegistrarServiceDeleteRejectedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceDeleteRejectedUnlock from hex string (Canton MCMS format)
func (t *RegistrarServiceDeleteRejectedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedBurn is a Record type
type RegistrarServiceExecuteAcceptedBurn struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.AcceptedBurnExecute `json:"payload"`
}

// ToMap converts RegistrarServiceExecuteAcceptedBurn to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceExecuteAcceptedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedBurn to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedForceTransfer is a Record type
type RegistrarServiceExecuteAcceptedForceTransfer struct {
	Cid     types.CONTRACT_ID                        `json:"cid"`
	Payload registry_v0.AcceptedForceTransferExecute `json:"payload"`
	Actor   types.PARTY                              `json:"actor"`
}

// ToMap converts RegistrarServiceExecuteAcceptedForceTransfer to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedForceTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarServiceExecuteAcceptedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedForceTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedForceTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedLock is a Record type
type RegistrarServiceExecuteAcceptedLock struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.AcceptedLockExecute `json:"payload"`
}

// ToMap converts RegistrarServiceExecuteAcceptedLock to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceExecuteAcceptedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedLock to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedLock from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedMint is a Record type
type RegistrarServiceExecuteAcceptedMint struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload registry_v0.AcceptedMintExecute `json:"payload"`
}

// ToMap converts RegistrarServiceExecuteAcceptedMint to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceExecuteAcceptedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedMint to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedMint from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedTransfer is a Record type
type RegistrarServiceExecuteAcceptedTransfer struct {
	Cid     types.CONTRACT_ID                   `json:"cid"`
	Payload registry_v0.AcceptedTransferExecute `json:"payload"`
}

// ToMap converts RegistrarServiceExecuteAcceptedTransfer to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceExecuteAcceptedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceExecuteAcceptedUnlock is a Record type
type RegistrarServiceExecuteAcceptedUnlock struct {
	Cid     types.CONTRACT_ID                 `json:"cid"`
	Payload registry_v0.AcceptedUnlockExecute `json:"payload"`
}

// ToMap converts RegistrarServiceExecuteAcceptedUnlock to a map for DAML arguments
func (t RegistrarServiceExecuteAcceptedUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceExecuteAcceptedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceExecuteAcceptedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceExecuteAcceptedUnlock to hex string (Canton MCMS format)
func (t RegistrarServiceExecuteAcceptedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceExecuteAcceptedUnlock from hex string (Canton MCMS format)
func (t *RegistrarServiceExecuteAcceptedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedBurn is a Record type
type RegistrarServiceFailAcceptedBurn struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.AcceptedBurnFail `json:"payload"`
}

// ToMap converts RegistrarServiceFailAcceptedBurn to a map for DAML arguments
func (t RegistrarServiceFailAcceptedBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceFailAcceptedBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedBurn to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedForceTransfer is a Record type
type RegistrarServiceFailAcceptedForceTransfer struct {
	Cid     types.CONTRACT_ID                     `json:"cid"`
	Payload registry_v0.AcceptedForceTransferFail `json:"payload"`
	Actor   types.PARTY                           `json:"actor"`
}

// ToMap converts RegistrarServiceFailAcceptedForceTransfer to a map for DAML arguments
func (t RegistrarServiceFailAcceptedForceTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RegistrarServiceFailAcceptedForceTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedForceTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedForceTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedForceTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedForceTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedForceTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedLock is a Record type
type RegistrarServiceFailAcceptedLock struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.AcceptedLockFail `json:"payload"`
}

// ToMap converts RegistrarServiceFailAcceptedLock to a map for DAML arguments
func (t RegistrarServiceFailAcceptedLock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceFailAcceptedLock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedLock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedLock to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedLock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedLock from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedLock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedMint is a Record type
type RegistrarServiceFailAcceptedMint struct {
	Cid     types.CONTRACT_ID            `json:"cid"`
	Payload registry_v0.AcceptedMintFail `json:"payload"`
}

// ToMap converts RegistrarServiceFailAcceptedMint to a map for DAML arguments
func (t RegistrarServiceFailAcceptedMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceFailAcceptedMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedMint to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedMint from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedTransfer is a Record type
type RegistrarServiceFailAcceptedTransfer struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_v0.AcceptedTransferFail `json:"payload"`
}

// ToMap converts RegistrarServiceFailAcceptedTransfer to a map for DAML arguments
func (t RegistrarServiceFailAcceptedTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceFailAcceptedTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedTransfer to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedTransfer from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceFailAcceptedUnlock is a Record type
type RegistrarServiceFailAcceptedUnlock struct {
	Cid     types.CONTRACT_ID              `json:"cid"`
	Payload registry_v0.AcceptedUnlockFail `json:"payload"`
}

// ToMap converts RegistrarServiceFailAcceptedUnlock to a map for DAML arguments
func (t RegistrarServiceFailAcceptedUnlock) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceFailAcceptedUnlock) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceFailAcceptedUnlock) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceFailAcceptedUnlock to hex string (Canton MCMS format)
func (t RegistrarServiceFailAcceptedUnlock) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceFailAcceptedUnlock from hex string (Canton MCMS format)
func (t *RegistrarServiceFailAcceptedUnlock) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceMergeHolding is a Record type
type RegistrarServiceMergeHolding struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_holding_v0.HoldingMerge `json:"payload"`
}

// ToMap converts RegistrarServiceMergeHolding to a map for DAML arguments
func (t RegistrarServiceMergeHolding) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceMergeHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceMergeHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceMergeHolding to hex string (Canton MCMS format)
func (t RegistrarServiceMergeHolding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceMergeHolding from hex string (Canton MCMS format)
func (t *RegistrarServiceMergeHolding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceOfferBurn is a Record type
type RegistrarServiceOfferBurn struct {
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Holder               types.PARTY                              `json:"holder"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts RegistrarServiceOfferBurn to a map for DAML arguments
func (t RegistrarServiceOfferBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t RegistrarServiceOfferBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceOfferBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceOfferBurn to hex string (Canton MCMS format)
func (t RegistrarServiceOfferBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceOfferBurn from hex string (Canton MCMS format)
func (t *RegistrarServiceOfferBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceOfferBurnResult is a Record type
type RegistrarServiceOfferBurnResult struct {
	BurnOfferCid types.CONTRACT_ID `json:"burnOfferCid"`
}

// ToMap converts RegistrarServiceOfferBurnResult to a map for DAML arguments
func (t RegistrarServiceOfferBurnResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["burnOfferCid"] = model.NestedToDAMLValue(t.BurnOfferCid)

	return m
}

func (t RegistrarServiceOfferBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceOfferBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceOfferBurnResult to hex string (Canton MCMS format)
func (t RegistrarServiceOfferBurnResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceOfferBurnResult from hex string (Canton MCMS format)
func (t *RegistrarServiceOfferBurnResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceOfferMint is a Record type
type RegistrarServiceOfferMint struct {
	InstrumentIdentifier registry_holding_v0.InstrumentIdentifier `json:"instrumentIdentifier"`
	Amount               types.NUMERIC                            `json:"amount"`
	Holder               types.PARTY                              `json:"holder"`
	Reference            types.TEXT                               `json:"reference"`
	Batch                registry_v0.Batch                        `json:"batch"`
}

// ToMap converts RegistrarServiceOfferMint to a map for DAML arguments
func (t RegistrarServiceOfferMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["instrumentIdentifier"] = model.NestedToDAMLValue(t.InstrumentIdentifier)

	m["amount"] = t.Amount

	m["holder"] = t.Holder.ToMap()

	m["reference"] = string(t.Reference)

	m["batch"] = model.NestedToDAMLValue(t.Batch)

	return m
}

func (t RegistrarServiceOfferMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceOfferMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceOfferMint to hex string (Canton MCMS format)
func (t RegistrarServiceOfferMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceOfferMint from hex string (Canton MCMS format)
func (t *RegistrarServiceOfferMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceOfferMintResult is a Record type
type RegistrarServiceOfferMintResult struct {
	MintOfferCid types.CONTRACT_ID `json:"mintOfferCid"`
}

// ToMap converts RegistrarServiceOfferMintResult to a map for DAML arguments
func (t RegistrarServiceOfferMintResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["mintOfferCid"] = model.NestedToDAMLValue(t.MintOfferCid)

	return m
}

func (t RegistrarServiceOfferMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceOfferMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceOfferMintResult to hex string (Canton MCMS format)
func (t RegistrarServiceOfferMintResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceOfferMintResult from hex string (Canton MCMS format)
func (t *RegistrarServiceOfferMintResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRejectBurnRequest is a Record type
type RegistrarServiceRejectBurnRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.BurnRequestReject `json:"payload"`
}

// ToMap converts RegistrarServiceRejectBurnRequest to a map for DAML arguments
func (t RegistrarServiceRejectBurnRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceRejectBurnRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRejectBurnRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRejectBurnRequest to hex string (Canton MCMS format)
func (t RegistrarServiceRejectBurnRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRejectBurnRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceRejectBurnRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRejectEnforcementServiceRequest is a Record type
type RegistrarServiceRejectEnforcementServiceRequest struct {
	Cid     types.CONTRACT_ID               `json:"cid"`
	Payload EnforcementServiceRequestReject `json:"payload"`
}

// ToMap converts RegistrarServiceRejectEnforcementServiceRequest to a map for DAML arguments
func (t RegistrarServiceRejectEnforcementServiceRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceRejectEnforcementServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRejectEnforcementServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRejectEnforcementServiceRequest to hex string (Canton MCMS format)
func (t RegistrarServiceRejectEnforcementServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRejectEnforcementServiceRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceRejectEnforcementServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRejectForceTransferRequest is a Record type
type RegistrarServiceRejectForceTransferRequest struct {
	Cid     types.CONTRACT_ID                      `json:"cid"`
	Payload registry_v0.ForceTransferRequestReject `json:"payload"`
}

// ToMap converts RegistrarServiceRejectForceTransferRequest to a map for DAML arguments
func (t RegistrarServiceRejectForceTransferRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceRejectForceTransferRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRejectForceTransferRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRejectForceTransferRequest to hex string (Canton MCMS format)
func (t RegistrarServiceRejectForceTransferRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRejectForceTransferRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceRejectForceTransferRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceRejectMintRequest is a Record type
type RegistrarServiceRejectMintRequest struct {
	Cid     types.CONTRACT_ID             `json:"cid"`
	Payload registry_v0.MintRequestReject `json:"payload"`
}

// ToMap converts RegistrarServiceRejectMintRequest to a map for DAML arguments
func (t RegistrarServiceRejectMintRequest) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceRejectMintRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceRejectMintRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceRejectMintRequest to hex string (Canton MCMS format)
func (t RegistrarServiceRejectMintRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceRejectMintRequest from hex string (Canton MCMS format)
func (t *RegistrarServiceRejectMintRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceSet is a Record type
type RegistrarServiceSet struct {
	EnableResultContracts *types.BOOL `json:"enableResultContracts" hex:"optional"`
}

// ToMap converts RegistrarServiceSet to a map for DAML arguments
func (t RegistrarServiceSet) ToMap() map[string]any {
	m := make(map[string]any)

	if t.EnableResultContracts != nil {
		m["enableResultContracts"] = map[string]any{
			"_type": "optional",
			"value": bool(*t.EnableResultContracts),
		}
	} else {
		m["enableResultContracts"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t RegistrarServiceSet) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceSet) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceSet to hex string (Canton MCMS format)
func (t RegistrarServiceSet) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceSet from hex string (Canton MCMS format)
func (t *RegistrarServiceSet) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceSplitHolding is a Record type
type RegistrarServiceSplitHolding struct {
	Cid     types.CONTRACT_ID                `json:"cid"`
	Payload registry_holding_v0.HoldingSplit `json:"payload"`
}

// ToMap converts RegistrarServiceSplitHolding to a map for DAML arguments
func (t RegistrarServiceSplitHolding) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceSplitHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceSplitHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceSplitHolding to hex string (Canton MCMS format)
func (t RegistrarServiceSplitHolding) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceSplitHolding from hex string (Canton MCMS format)
func (t *RegistrarServiceSplitHolding) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceTerminate is a Record type
type RegistrarServiceTerminate struct {
}

// ToMap converts RegistrarServiceTerminate to a map for DAML arguments
func (t RegistrarServiceTerminate) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceTerminate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceTerminate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceTerminate to hex string (Canton MCMS format)
func (t RegistrarServiceTerminate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceTerminate from hex string (Canton MCMS format)
func (t *RegistrarServiceTerminate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceTerminateEnforcementService is a Record type
type RegistrarServiceTerminateEnforcementService struct {
	Cid     types.CONTRACT_ID           `json:"cid"`
	Payload EnforcementServiceTerminate `json:"payload"`
}

// ToMap converts RegistrarServiceTerminateEnforcementService to a map for DAML arguments
func (t RegistrarServiceTerminateEnforcementService) ToMap() map[string]any {
	m := make(map[string]any)

	m["cid"] = model.NestedToDAMLValue(t.Cid)

	m["payload"] = model.NestedToDAMLValue(t.Payload)

	return m
}

func (t RegistrarServiceTerminateEnforcementService) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceTerminateEnforcementService) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceTerminateEnforcementService to hex string (Canton MCMS format)
func (t RegistrarServiceTerminateEnforcementService) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceTerminateEnforcementService from hex string (Canton MCMS format)
func (t *RegistrarServiceTerminateEnforcementService) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RegistrarServiceTerminateResult is a Record type
type RegistrarServiceTerminateResult struct {
}

// ToMap converts RegistrarServiceTerminateResult to a map for DAML arguments
func (t RegistrarServiceTerminateResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RegistrarServiceTerminateResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RegistrarServiceTerminateResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RegistrarServiceTerminateResult to hex string (Canton MCMS format)
func (t RegistrarServiceTerminateResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RegistrarServiceTerminateResult from hex string (Canton MCMS format)
func (t *RegistrarServiceTerminateResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedBurn is a Template type
type RejectedBurn struct {
	Operator           types.PARTY `json:"operator"`
	Provider           types.PARTY `json:"provider"`
	Burn               Burn        `json:"burn"`
	Reason             types.TEXT  `json:"reason"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedBurn) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedBurn) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedBurn) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

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
func (t RejectedBurn) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["burn"] = model.NestedToDAMLValue(t.Burn)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

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

// RejectedBurnDelete exercises the RejectedBurn_Delete choice on this RejectedBurn contract
// This method uses the package name in the template ID
func (t RejectedBurn) RejectedBurnDelete(contractID string, args RejectedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedBurnDeleteWithPackageID exercises the RejectedBurn_Delete choice using the provided package ID instead of package name
func (t RejectedBurn) RejectedBurnDeleteWithPackageID(contractID string, packageID string, args RejectedBurnDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "RejectedBurn_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedBurn contract
// This method uses the package name in the template ID
func (t RejectedBurn) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedBurn) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Burn", "RejectedBurn"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedBurnDelete is a Record type
type RejectedBurnDelete struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedBurnDelete to a map for DAML arguments
func (t RejectedBurnDelete) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

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

// RejectedEnforcementServiceRequest is a Template type
type RejectedEnforcementServiceRequest struct {
	Request EnforcementServiceRequest `json:"request"`
	Reason  types.TEXT                `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedEnforcementServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedEnforcementServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedEnforcementServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedEnforcementServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedEnforcementServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedEnforcementServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedEnforcementServiceRequest to hex string (Canton MCMS format)
func (t RejectedEnforcementServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedEnforcementServiceRequest from hex string (Canton MCMS format)
func (t *RejectedEnforcementServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedEnforcementServiceRequest

// Archive exercises the Archive choice on this RejectedEnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedEnforcementServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedEnforcementServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedEnforcementServiceRequestDelete exercises the RejectedEnforcementServiceRequest_Delete choice on this RejectedEnforcementServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedEnforcementServiceRequest) RejectedEnforcementServiceRequestDelete(contractID string, args RejectedEnforcementServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedEnforcementServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedEnforcementServiceRequestDeleteWithPackageID exercises the RejectedEnforcementServiceRequest_Delete choice using the provided package ID instead of package name
func (t RejectedEnforcementServiceRequest) RejectedEnforcementServiceRequestDeleteWithPackageID(contractID string, packageID string, args RejectedEnforcementServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Enforcement", "RejectedEnforcementServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedEnforcementServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedEnforcementServiceRequestDelete is a Record type
type RejectedEnforcementServiceRequestDelete struct {
}

// ToMap converts RejectedEnforcementServiceRequestDelete to a map for DAML arguments
func (t RejectedEnforcementServiceRequestDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedEnforcementServiceRequestDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedEnforcementServiceRequestDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedEnforcementServiceRequestDelete to hex string (Canton MCMS format)
func (t RejectedEnforcementServiceRequestDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedEnforcementServiceRequestDelete from hex string (Canton MCMS format)
func (t *RejectedEnforcementServiceRequestDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedEnforcementServiceRequestDeleteResult is a Record type
type RejectedEnforcementServiceRequestDeleteResult struct {
}

// ToMap converts RejectedEnforcementServiceRequestDeleteResult to a map for DAML arguments
func (t RejectedEnforcementServiceRequestDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedEnforcementServiceRequestDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedEnforcementServiceRequestDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedEnforcementServiceRequestDeleteResult to hex string (Canton MCMS format)
func (t RejectedEnforcementServiceRequestDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedEnforcementServiceRequestDeleteResult from hex string (Canton MCMS format)
func (t *RejectedEnforcementServiceRequestDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedHolderServiceRequest is a Template type
type RejectedHolderServiceRequest struct {
	Request HolderServiceRequest `json:"request"`
	Reason  types.TEXT           `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedHolderServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedHolderServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedHolderServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedHolderServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedHolderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedHolderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedHolderServiceRequest to hex string (Canton MCMS format)
func (t RejectedHolderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedHolderServiceRequest from hex string (Canton MCMS format)
func (t *RejectedHolderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedHolderServiceRequest

// RejectedHolderServiceRequestClean exercises the RejectedHolderServiceRequest_Clean choice on this RejectedHolderServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedHolderServiceRequest) RejectedHolderServiceRequestClean(contractID string, args RejectedHolderServiceRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedHolderServiceRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// RejectedHolderServiceRequestCleanWithPackageID exercises the RejectedHolderServiceRequest_Clean choice using the provided package ID instead of package name
func (t RejectedHolderServiceRequest) RejectedHolderServiceRequestCleanWithPackageID(contractID string, packageID string, args RejectedHolderServiceRequestClean) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedHolderServiceRequest_Clean",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedHolderServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedHolderServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedHolderServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedHolderServiceRequestDelete exercises the RejectedHolderServiceRequest_Delete choice on this RejectedHolderServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedHolderServiceRequest) RejectedHolderServiceRequestDelete(contractID string, args RejectedHolderServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedHolderServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedHolderServiceRequestDeleteWithPackageID exercises the RejectedHolderServiceRequest_Delete choice using the provided package ID instead of package name
func (t RejectedHolderServiceRequest) RejectedHolderServiceRequestDeleteWithPackageID(contractID string, packageID string, args RejectedHolderServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Holder", "RejectedHolderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedHolderServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedHolderServiceRequestClean is a Record type
type RejectedHolderServiceRequestClean struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedHolderServiceRequestClean to a map for DAML arguments
func (t RejectedHolderServiceRequestClean) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t RejectedHolderServiceRequestClean) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedHolderServiceRequestClean) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedHolderServiceRequestClean to hex string (Canton MCMS format)
func (t RejectedHolderServiceRequestClean) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedHolderServiceRequestClean from hex string (Canton MCMS format)
func (t *RejectedHolderServiceRequestClean) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedHolderServiceRequestDelete is a Record type
type RejectedHolderServiceRequestDelete struct {
}

// ToMap converts RejectedHolderServiceRequestDelete to a map for DAML arguments
func (t RejectedHolderServiceRequestDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedHolderServiceRequestDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedHolderServiceRequestDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedHolderServiceRequestDelete to hex string (Canton MCMS format)
func (t RejectedHolderServiceRequestDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedHolderServiceRequestDelete from hex string (Canton MCMS format)
func (t *RejectedHolderServiceRequestDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedHolderServiceRequestDeleteResult is a Record type
type RejectedHolderServiceRequestDeleteResult struct {
}

// ToMap converts RejectedHolderServiceRequestDeleteResult to a map for DAML arguments
func (t RejectedHolderServiceRequestDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedHolderServiceRequestDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedHolderServiceRequestDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedHolderServiceRequestDeleteResult to hex string (Canton MCMS format)
func (t RejectedHolderServiceRequestDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedHolderServiceRequestDeleteResult from hex string (Canton MCMS format)
func (t *RejectedHolderServiceRequestDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedMint is a Template type
type RejectedMint struct {
	Operator           types.PARTY `json:"operator"`
	Provider           types.PARTY `json:"provider"`
	Mint               Mint        `json:"mint"`
	Reason             types.TEXT  `json:"reason"`
	OperatorIsObserver *types.BOOL `json:"operatorIsObserver" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedMint) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "RejectedMint")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedMint) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "RejectedMint")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedMint) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

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
func (t RejectedMint) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mint"] = model.NestedToDAMLValue(t.Mint)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

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

// RejectedMintDelete exercises the RejectedMint_Delete choice on this RejectedMint contract
// This method uses the package name in the template ID
func (t RejectedMint) RejectedMintDelete(contractID string, args RejectedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedMintDeleteWithPackageID exercises the RejectedMint_Delete choice using the provided package ID instead of package name
func (t RejectedMint) RejectedMintDeleteWithPackageID(contractID string, packageID string, args RejectedMintDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "RejectedMint_Delete",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RejectedMint contract
// This method uses the package name in the template ID
func (t RejectedMint) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedMint) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Mint", "RejectedMint"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedMintDelete is a Record type
type RejectedMintDelete struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts RejectedMintDelete to a map for DAML arguments
func (t RejectedMintDelete) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

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

// RejectedProviderServiceRequest is a Template type
type RejectedProviderServiceRequest struct {
	Request ProviderServiceRequest `json:"request"`
	Reason  types.TEXT             `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedProviderServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedProviderServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedProviderServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedProviderServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedProviderServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedProviderServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedProviderServiceRequest to hex string (Canton MCMS format)
func (t RejectedProviderServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedProviderServiceRequest from hex string (Canton MCMS format)
func (t *RejectedProviderServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedProviderServiceRequest

// Archive exercises the Archive choice on this RejectedProviderServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedProviderServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedProviderServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedProviderServiceRequestDelete exercises the RejectedProviderServiceRequest_Delete choice on this RejectedProviderServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedProviderServiceRequest) RejectedProviderServiceRequestDelete(contractID string, args RejectedProviderServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedProviderServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedProviderServiceRequestDeleteWithPackageID exercises the RejectedProviderServiceRequest_Delete choice using the provided package ID instead of package name
func (t RejectedProviderServiceRequest) RejectedProviderServiceRequestDeleteWithPackageID(contractID string, packageID string, args RejectedProviderServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Provider", "RejectedProviderServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedProviderServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedProviderServiceRequestDelete is a Record type
type RejectedProviderServiceRequestDelete struct {
}

// ToMap converts RejectedProviderServiceRequestDelete to a map for DAML arguments
func (t RejectedProviderServiceRequestDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedProviderServiceRequestDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedProviderServiceRequestDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedProviderServiceRequestDelete to hex string (Canton MCMS format)
func (t RejectedProviderServiceRequestDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedProviderServiceRequestDelete from hex string (Canton MCMS format)
func (t *RejectedProviderServiceRequestDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedProviderServiceRequestDeleteResult is a Record type
type RejectedProviderServiceRequestDeleteResult struct {
}

// ToMap converts RejectedProviderServiceRequestDeleteResult to a map for DAML arguments
func (t RejectedProviderServiceRequestDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedProviderServiceRequestDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedProviderServiceRequestDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedProviderServiceRequestDeleteResult to hex string (Canton MCMS format)
func (t RejectedProviderServiceRequestDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedProviderServiceRequestDeleteResult from hex string (Canton MCMS format)
func (t *RejectedProviderServiceRequestDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedRegistrarServiceRequest is a Template type
type RejectedRegistrarServiceRequest struct {
	Request RegistrarServiceRequest `json:"request"`
	Reason  types.TEXT              `json:"reason"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RejectedRegistrarServiceRequest) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RejectedRegistrarServiceRequest) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RejectedRegistrarServiceRequest) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t RejectedRegistrarServiceRequest) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["request"] = model.NestedToDAMLValue(t.Request)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reason"] = string(t.Reason)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t RejectedRegistrarServiceRequest) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedRegistrarServiceRequest) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedRegistrarServiceRequest to hex string (Canton MCMS format)
func (t RejectedRegistrarServiceRequest) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedRegistrarServiceRequest from hex string (Canton MCMS format)
func (t *RejectedRegistrarServiceRequest) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RejectedRegistrarServiceRequest

// Archive exercises the Archive choice on this RejectedRegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedRegistrarServiceRequest) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RejectedRegistrarServiceRequest) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// RejectedRegistrarServiceRequestDelete exercises the RejectedRegistrarServiceRequest_Delete choice on this RejectedRegistrarServiceRequest contract
// This method uses the package name in the template ID
func (t RejectedRegistrarServiceRequest) RejectedRegistrarServiceRequestDelete(contractID string, args RejectedRegistrarServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedRegistrarServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedRegistrarServiceRequestDeleteWithPackageID exercises the RejectedRegistrarServiceRequest_Delete choice using the provided package ID instead of package name
func (t RejectedRegistrarServiceRequest) RejectedRegistrarServiceRequestDeleteWithPackageID(contractID string, packageID string, args RejectedRegistrarServiceRequestDelete) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Service.Registrar", "RejectedRegistrarServiceRequest"),
		ContractID: contractID,
		Choice:     "RejectedRegistrarServiceRequest_Delete",
		Arguments:  argsToMap(args),
	}
}

// RejectedRegistrarServiceRequestDelete is a Record type
type RejectedRegistrarServiceRequestDelete struct {
}

// ToMap converts RejectedRegistrarServiceRequestDelete to a map for DAML arguments
func (t RejectedRegistrarServiceRequestDelete) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedRegistrarServiceRequestDelete) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedRegistrarServiceRequestDelete) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedRegistrarServiceRequestDelete to hex string (Canton MCMS format)
func (t RejectedRegistrarServiceRequestDelete) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedRegistrarServiceRequestDelete from hex string (Canton MCMS format)
func (t *RejectedRegistrarServiceRequestDelete) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RejectedRegistrarServiceRequestDeleteResult is a Record type
type RejectedRegistrarServiceRequestDeleteResult struct {
}

// ToMap converts RejectedRegistrarServiceRequestDeleteResult to a map for DAML arguments
func (t RejectedRegistrarServiceRequestDeleteResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t RejectedRegistrarServiceRequestDeleteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RejectedRegistrarServiceRequestDeleteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RejectedRegistrarServiceRequestDeleteResult to hex string (Canton MCMS format)
func (t RejectedRegistrarServiceRequestDeleteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RejectedRegistrarServiceRequestDeleteResult from hex string (Canton MCMS format)
func (t *RejectedRegistrarServiceRequestDeleteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferOffer is a Template type
type TransferOffer struct {
	Operator types.PARTY                                       `json:"operator"`
	Provider types.PARTY                                       `json:"provider"`
	Transfer splice_api_token_transfer_instruction_v1.Transfer `json:"transfer"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TransferOffer) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferOffer")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TransferOffer) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferOffer")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TransferOffer) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TransferOffer) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["provider"] = t.Provider.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transfer"] = model.NestedToDAMLValue(t.Transfer)

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

// Archive exercises the Archive choice on this TransferOffer contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t TransferOffer) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TransferOffer) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this TransferOffer contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t TransferOffer) TransferInstructionAccept(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionAcceptWithPackageID exercises the TransferInstruction_Accept choice using the provided package ID instead of package name
func (t TransferOffer) TransferInstructionAcceptWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this TransferOffer contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t TransferOffer) TransferInstructionReject(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionRejectWithPackageID exercises the TransferInstruction_Reject choice using the provided package ID instead of package name
func (t TransferOffer) TransferInstructionRejectWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this TransferOffer contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t TransferOffer) TransferInstructionWithdraw(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdrawWithPackageID exercises the TransferInstruction_Withdraw choice using the provided package ID instead of package name
func (t TransferOffer) TransferInstructionWithdrawWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this TransferOffer contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t TransferOffer) TransferInstructionUpdate(contractID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdateWithPackageID exercises the TransferInstruction_Update choice using the provided package ID instead of package name
func (t TransferOffer) TransferInstructionUpdateWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for TransferOffer

var _ splice_api_token_transfer_instruction_v1.ITransferInstruction = (*TransferOffer)(nil)

// TransferPreapproval is a Template type
type TransferPreapproval struct {
	Operator             types.PARTY           `json:"operator"`
	Receiver             types.PARTY           `json:"receiver"`
	InstrumentAdmin      types.PARTY           `json:"instrumentAdmin"`
	InstrumentAllowances []InstrumentAllowance `json:"instrumentAllowances"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t TransferPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t TransferPreapproval) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t TransferPreapproval) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentAdmin"] = t.InstrumentAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentAllowances"] = func() []any {
		res := make([]any, 0, len(t.InstrumentAllowances))
		for _, e := range t.InstrumentAllowances {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t TransferPreapproval) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["operator"] = t.Operator.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentAdmin"] = t.InstrumentAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentAllowances"] = func() []any {
		res := make([]any, 0, len(t.InstrumentAllowances))
		for _, e := range t.InstrumentAllowances {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t TransferPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPreapproval to hex string (Canton MCMS format)
func (t TransferPreapproval) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPreapproval from hex string (Canton MCMS format)
func (t *TransferPreapproval) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for TransferPreapproval

// TransferPreapprovalWithdraw exercises the TransferPreapproval_Withdraw choice on this TransferPreapproval contract
// This method uses the package name in the template ID
func (t TransferPreapproval) TransferPreapprovalWithdraw(contractID string, args TransferPreapprovalWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalWithdrawWithPackageID exercises the TransferPreapproval_Withdraw choice using the provided package ID instead of package name
func (t TransferPreapproval) TransferPreapprovalWithdrawWithPackageID(contractID string, packageID string, args TransferPreapprovalWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this TransferPreapproval contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t TransferPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t TransferPreapproval) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// TransferPreapprovalModify exercises the TransferPreapproval_Modify choice on this TransferPreapproval contract
// This method uses the package name in the template ID
func (t TransferPreapproval) TransferPreapprovalModify(contractID string, args TransferPreapprovalModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Modify",
		Arguments:  argsToMap(args),
	}
}

// TransferPreapprovalModifyWithPackageID exercises the TransferPreapproval_Modify choice using the provided package ID instead of package name
func (t TransferPreapproval) TransferPreapprovalModifyWithPackageID(contractID string, packageID string, args TransferPreapprovalModify) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferPreapproval"),
		ContractID: contractID,
		Choice:     "TransferPreapproval_Modify",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this TransferPreapproval contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t TransferPreapproval) TransferFactoryTransfer(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransferWithPackageID exercises the TransferFactory_Transfer choice using the provided package ID instead of package name
func (t TransferPreapproval) TransferFactoryTransferWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this TransferPreapproval contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t TransferPreapproval) TransferFactoryPublicFetch(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetchWithPackageID exercises the TransferFactory_PublicFetch choice using the provided package ID instead of package name
func (t TransferPreapproval) TransferFactoryPublicFetchWithPackageID(contractID string, packageID string, args splice_api_token_transfer_instruction_v1.TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Utility.Registry.App.V0.Model.TransferPreapproval", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for TransferPreapproval

var _ splice_api_token_transfer_instruction_v1.ITransferFactory = (*TransferPreapproval)(nil)

// TransferPreapprovalModify is a Record type
type TransferPreapprovalModify struct {
	NewInstrumentAllowances []InstrumentAllowance `json:"newInstrumentAllowances"`
}

// ToMap converts TransferPreapprovalModify to a map for DAML arguments
func (t TransferPreapprovalModify) ToMap() map[string]any {
	m := make(map[string]any)

	m["newInstrumentAllowances"] = func() []any {
		res := make([]any, 0, len(t.NewInstrumentAllowances))
		for _, e := range t.NewInstrumentAllowances {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t TransferPreapprovalModify) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPreapprovalModify) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPreapprovalModify to hex string (Canton MCMS format)
func (t TransferPreapprovalModify) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPreapprovalModify from hex string (Canton MCMS format)
func (t *TransferPreapprovalModify) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferPreapprovalModifyResult is a Record type
type TransferPreapprovalModifyResult struct {
	TransferPreapprovalCid types.CONTRACT_ID `json:"transferPreapprovalCid"`
}

// ToMap converts TransferPreapprovalModifyResult to a map for DAML arguments
func (t TransferPreapprovalModifyResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferPreapprovalCid"] = model.NestedToDAMLValue(t.TransferPreapprovalCid)

	return m
}

func (t TransferPreapprovalModifyResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPreapprovalModifyResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPreapprovalModifyResult to hex string (Canton MCMS format)
func (t TransferPreapprovalModifyResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPreapprovalModifyResult from hex string (Canton MCMS format)
func (t *TransferPreapprovalModifyResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferPreapprovalWithdraw is a Record type
type TransferPreapprovalWithdraw struct {
	Actor types.PARTY `json:"actor"`
}

// ToMap converts TransferPreapprovalWithdraw to a map for DAML arguments
func (t TransferPreapprovalWithdraw) ToMap() map[string]any {
	m := make(map[string]any)

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t TransferPreapprovalWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPreapprovalWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPreapprovalWithdraw to hex string (Canton MCMS format)
func (t TransferPreapprovalWithdraw) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPreapprovalWithdraw from hex string (Canton MCMS format)
func (t *TransferPreapprovalWithdraw) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferPreapprovalWithdrawResult is a Record type
type TransferPreapprovalWithdrawResult struct {
}

// ToMap converts TransferPreapprovalWithdrawResult to a map for DAML arguments
func (t TransferPreapprovalWithdrawResult) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t TransferPreapprovalWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPreapprovalWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPreapprovalWithdrawResult to hex string (Canton MCMS format)
func (t TransferPreapprovalWithdrawResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPreapprovalWithdrawResult from hex string (Canton MCMS format)
func (t *TransferPreapprovalWithdrawResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AllocationFactoryAllocateInternal(args AllocationFactoryAllocateInternal) (*bind.EncodedChoice, error)
	AllocationFactoryInternalBurnMint(args AllocationFactoryInternalBurnMint) (*bind.EncodedChoice, error)
	AllocationFactoryOfferBurn(args AllocationFactoryOfferBurn) (*bind.EncodedChoice, error)
	AllocationFactoryOfferMint(args AllocationFactoryOfferMint) (*bind.EncodedChoice, error)
	AllocationFactoryRequestBurn(args AllocationFactoryRequestBurn) (*bind.EncodedChoice, error)
	AllocationFactoryRequestMint(args AllocationFactoryRequestMint) (*bind.EncodedChoice, error)
	AllocationFactoryTransferInternal(args AllocationFactoryTransferInternal) (*bind.EncodedChoice, error)
	BurnOfferAccept(args BurnOfferAccept) (*bind.EncodedChoice, error)
	BurnOfferCancel(args BurnOfferCancel) (*bind.EncodedChoice, error)
	BurnOfferReject(args BurnOfferReject) (*bind.EncodedChoice, error)
	BurnRequestAccept(args BurnRequestAccept) (*bind.EncodedChoice, error)
	BurnRequestCancel(args BurnRequestCancel) (*bind.EncodedChoice, error)
	BurnRequestReject(args BurnRequestReject) (*bind.EncodedChoice, error)
	EnforcementServiceRequestAccept(args EnforcementServiceRequestAccept) (*bind.EncodedChoice, error)
	EnforcementServiceRequestCancel(args EnforcementServiceRequestCancel) (*bind.EncodedChoice, error)
	EnforcementServiceRequestReject(args EnforcementServiceRequestReject) (*bind.EncodedChoice, error)
	EnforcementServiceAcceptForceTransferRequest(args EnforcementServiceAcceptForceTransferRequest) (*bind.EncodedChoice, error)
	EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization(args EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) (*bind.EncodedChoice, error)
	EnforcementServiceTerminate(args EnforcementServiceTerminate) (*bind.EncodedChoice, error)
	ExecutedBurnDelete(args ExecutedBurnDelete) (*bind.EncodedChoice, error)
	ExecutedMintDelete(args ExecutedMintDelete) (*bind.EncodedChoice, error)
	HolderServiceRequestAccept(args HolderServiceRequestAccept) (*bind.EncodedChoice, error)
	HolderServiceRequestCancel(args HolderServiceRequestCancel) (*bind.EncodedChoice, error)
	HolderServiceRequestClean(args HolderServiceRequestClean) (*bind.EncodedChoice, error)
	HolderServiceRequestReject(args HolderServiceRequestReject) (*bind.EncodedChoice, error)
	HolderServiceAcceptBurnOffer(args HolderServiceAcceptBurnOffer) (*bind.EncodedChoice, error)
	HolderServiceAcceptLockOffer(args HolderServiceAcceptLockOffer) (*bind.EncodedChoice, error)
	HolderServiceAcceptLockRequest(args HolderServiceAcceptLockRequest) (*bind.EncodedChoice, error)
	HolderServiceAcceptMintOffer(args HolderServiceAcceptMintOffer) (*bind.EncodedChoice, error)
	HolderServiceAcceptTransferOffer(args HolderServiceAcceptTransferOffer) (*bind.EncodedChoice, error)
	HolderServiceAcceptTransferRequest(args HolderServiceAcceptTransferRequest) (*bind.EncodedChoice, error)
	HolderServiceAcceptUnlockOffer(args HolderServiceAcceptUnlockOffer) (*bind.EncodedChoice, error)
	HolderServiceAcceptUnlockRequest(args HolderServiceAcceptUnlockRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelBurnRequest(args HolderServiceCancelBurnRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelEnforcementServiceRequest(args HolderServiceCancelEnforcementServiceRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelForceTransferRequest(args HolderServiceCancelForceTransferRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelLockOffer(args HolderServiceCancelLockOffer) (*bind.EncodedChoice, error)
	HolderServiceCancelLockRequest(args HolderServiceCancelLockRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelMintRequest(args HolderServiceCancelMintRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelTransferOffer(args HolderServiceCancelTransferOffer) (*bind.EncodedChoice, error)
	HolderServiceCancelTransferRequest(args HolderServiceCancelTransferRequest) (*bind.EncodedChoice, error)
	HolderServiceCancelUnlockOffer(args HolderServiceCancelUnlockOffer) (*bind.EncodedChoice, error)
	HolderServiceCancelUnlockRequest(args HolderServiceCancelUnlockRequest) (*bind.EncodedChoice, error)
	HolderServiceClean(args HolderServiceClean) (*bind.EncodedChoice, error)
	HolderServiceCreateAllocation(args HolderServiceCreateAllocation) (*bind.EncodedChoice, error)
	HolderServiceOfferLock(args HolderServiceOfferLock) (*bind.EncodedChoice, error)
	HolderServiceOfferTransfer(args HolderServiceOfferTransfer) (*bind.EncodedChoice, error)
	HolderServiceOfferUnlock(args HolderServiceOfferUnlock) (*bind.EncodedChoice, error)
	HolderServiceRejectAllocationRequest(args HolderServiceRejectAllocationRequest) (*bind.EncodedChoice, error)
	HolderServiceRejectBurnOffer(args HolderServiceRejectBurnOffer) (*bind.EncodedChoice, error)
	HolderServiceRejectLockOffer(args HolderServiceRejectLockOffer) (*bind.EncodedChoice, error)
	HolderServiceRejectLockRequest(args HolderServiceRejectLockRequest) (*bind.EncodedChoice, error)
	HolderServiceRejectMintOffer(args HolderServiceRejectMintOffer) (*bind.EncodedChoice, error)
	HolderServiceRejectTransferOffer(args HolderServiceRejectTransferOffer) (*bind.EncodedChoice, error)
	HolderServiceRejectTransferRequest(args HolderServiceRejectTransferRequest) (*bind.EncodedChoice, error)
	HolderServiceRejectUnlockOffer(args HolderServiceRejectUnlockOffer) (*bind.EncodedChoice, error)
	HolderServiceRejectUnlockRequest(args HolderServiceRejectUnlockRequest) (*bind.EncodedChoice, error)
	HolderServiceRequestBurn(args HolderServiceRequestBurn) (*bind.EncodedChoice, error)
	HolderServiceRequestEnforcementService(args HolderServiceRequestEnforcementService) (*bind.EncodedChoice, error)
	HolderServiceRequestForceTransfer(args HolderServiceRequestForceTransfer) (*bind.EncodedChoice, error)
	HolderServiceRequestLock(args HolderServiceRequestLock) (*bind.EncodedChoice, error)
	HolderServiceRequestMint(args HolderServiceRequestMint) (*bind.EncodedChoice, error)
	HolderServiceRequestTransfer(args HolderServiceRequestTransfer) (*bind.EncodedChoice, error)
	HolderServiceRequestUnlock(args HolderServiceRequestUnlock) (*bind.EncodedChoice, error)
	HolderServiceTerminate(args HolderServiceTerminate) (*bind.EncodedChoice, error)
	MintOfferAccept(args MintOfferAccept) (*bind.EncodedChoice, error)
	MintOfferCancel(args MintOfferCancel) (*bind.EncodedChoice, error)
	MintOfferReject(args MintOfferReject) (*bind.EncodedChoice, error)
	MintRequestAccept(args MintRequestAccept) (*bind.EncodedChoice, error)
	MintRequestCancel(args MintRequestCancel) (*bind.EncodedChoice, error)
	MintRequestReject(args MintRequestReject) (*bind.EncodedChoice, error)
	OperatorConfigurationGet(args OperatorConfigurationGet) (*bind.EncodedChoice, error)
	OperatorConfigurationModify(args OperatorConfigurationModify) (*bind.EncodedChoice, error)
	ProviderConfigurationGet(args ProviderConfigurationGet) (*bind.EncodedChoice, error)
	ProviderServiceRequestAccept(args ProviderServiceRequestAccept) (*bind.EncodedChoice, error)
	ProviderServiceRequestCancel(args ProviderServiceRequestCancel) (*bind.EncodedChoice, error)
	ProviderServiceRequestReject(args ProviderServiceRequestReject) (*bind.EncodedChoice, error)
	ProviderServiceAcceptHolderServiceRequest(args ProviderServiceAcceptHolderServiceRequest) (*bind.EncodedChoice, error)
	ProviderServiceAcceptRegistrarServiceRequest(args ProviderServiceAcceptRegistrarServiceRequest) (*bind.EncodedChoice, error)
	ProviderServiceArchiveAndCreateProviderConfiguration(args ProviderServiceArchiveAndCreateProviderConfiguration) (*bind.EncodedChoice, error)
	ProviderServiceArchiveProviderConfiguration(args ProviderServiceArchiveProviderConfiguration) (*bind.EncodedChoice, error)
	ProviderServiceCreateProviderConfiguration(args ProviderServiceCreateProviderConfiguration) (*bind.EncodedChoice, error)
	ProviderServiceRejectHolderServiceRequest(args ProviderServiceRejectHolderServiceRequest) (*bind.EncodedChoice, error)
	ProviderServiceRejectRegistrarServiceRequest(args ProviderServiceRejectRegistrarServiceRequest) (*bind.EncodedChoice, error)
	ProviderServiceTerminate(args ProviderServiceTerminate) (*bind.EncodedChoice, error)
	RegistrarConfigurationGet(args RegistrarConfigurationGet) (*bind.EncodedChoice, error)
	RegistrarServiceRequestAccept(args RegistrarServiceRequestAccept) (*bind.EncodedChoice, error)
	RegistrarServiceRequestCancel(args RegistrarServiceRequestCancel) (*bind.EncodedChoice, error)
	RegistrarServiceRequestReject(args RegistrarServiceRequestReject) (*bind.EncodedChoice, error)
	RegistrarServiceAcceptBurnRequest(args RegistrarServiceAcceptBurnRequest) (*bind.EncodedChoice, error)
	RegistrarServiceAcceptEnforcementServiceRequest(args RegistrarServiceAcceptEnforcementServiceRequest) (*bind.EncodedChoice, error)
	RegistrarServiceAcceptForceTransferRequest(args RegistrarServiceAcceptForceTransferRequest) (*bind.EncodedChoice, error)
	RegistrarServiceAcceptMintRequest(args RegistrarServiceAcceptMintRequest) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveAllocationFactory(args RegistrarServiceArchiveAllocationFactory) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveAndCreateInstrumentConfiguration(args RegistrarServiceArchiveAndCreateInstrumentConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveAndCreateRegistrarConfiguration(args RegistrarServiceArchiveAndCreateRegistrarConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveInstrumentConfiguration(args RegistrarServiceArchiveInstrumentConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveRegistrarConfiguration(args RegistrarServiceArchiveRegistrarConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceArchiveTransferRule(args RegistrarServiceArchiveTransferRule) (*bind.EncodedChoice, error)
	RegistrarServiceCancelBurnOffer(args RegistrarServiceCancelBurnOffer) (*bind.EncodedChoice, error)
	RegistrarServiceCancelMintOffer(args RegistrarServiceCancelMintOffer) (*bind.EncodedChoice, error)
	RegistrarServiceCreateAllocationFactory(args RegistrarServiceCreateAllocationFactory) (*bind.EncodedChoice, error)
	RegistrarServiceCreateInstrumentConfiguration(args RegistrarServiceCreateInstrumentConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceCreateRegistrarConfiguration(args RegistrarServiceCreateRegistrarConfiguration) (*bind.EncodedChoice, error)
	RegistrarServiceCreateTransferRule(args RegistrarServiceCreateTransferRule) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedBurn(args RegistrarServiceDeleteExecutedBurn) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedLock(args RegistrarServiceDeleteExecutedLock) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedMint(args RegistrarServiceDeleteExecutedMint) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedTransfer(args RegistrarServiceDeleteExecutedTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedTransfers(args RegistrarServiceDeleteExecutedTransfers) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteExecutedUnlock(args RegistrarServiceDeleteExecutedUnlock) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedBurn(args RegistrarServiceDeleteFailedBurn) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedLock(args RegistrarServiceDeleteFailedLock) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedMint(args RegistrarServiceDeleteFailedMint) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedTransfer(args RegistrarServiceDeleteFailedTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedTransfers(args RegistrarServiceDeleteFailedTransfers) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteFailedUnlock(args RegistrarServiceDeleteFailedUnlock) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedBurn(args RegistrarServiceDeleteRejectedBurn) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedLock(args RegistrarServiceDeleteRejectedLock) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedMint(args RegistrarServiceDeleteRejectedMint) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedTransfer(args RegistrarServiceDeleteRejectedTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedTransfers(args RegistrarServiceDeleteRejectedTransfers) (*bind.EncodedChoice, error)
	RegistrarServiceDeleteRejectedUnlock(args RegistrarServiceDeleteRejectedUnlock) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedBurn(args RegistrarServiceExecuteAcceptedBurn) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedForceTransfer(args RegistrarServiceExecuteAcceptedForceTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedLock(args RegistrarServiceExecuteAcceptedLock) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedMint(args RegistrarServiceExecuteAcceptedMint) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedTransfer(args RegistrarServiceExecuteAcceptedTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceExecuteAcceptedUnlock(args RegistrarServiceExecuteAcceptedUnlock) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedBurn(args RegistrarServiceFailAcceptedBurn) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedForceTransfer(args RegistrarServiceFailAcceptedForceTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedLock(args RegistrarServiceFailAcceptedLock) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedMint(args RegistrarServiceFailAcceptedMint) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedTransfer(args RegistrarServiceFailAcceptedTransfer) (*bind.EncodedChoice, error)
	RegistrarServiceFailAcceptedUnlock(args RegistrarServiceFailAcceptedUnlock) (*bind.EncodedChoice, error)
	RegistrarServiceMergeHolding(args RegistrarServiceMergeHolding) (*bind.EncodedChoice, error)
	RegistrarServiceOfferBurn(args RegistrarServiceOfferBurn) (*bind.EncodedChoice, error)
	RegistrarServiceOfferMint(args RegistrarServiceOfferMint) (*bind.EncodedChoice, error)
	RegistrarServiceRejectBurnRequest(args RegistrarServiceRejectBurnRequest) (*bind.EncodedChoice, error)
	RegistrarServiceRejectEnforcementServiceRequest(args RegistrarServiceRejectEnforcementServiceRequest) (*bind.EncodedChoice, error)
	RegistrarServiceRejectForceTransferRequest(args RegistrarServiceRejectForceTransferRequest) (*bind.EncodedChoice, error)
	RegistrarServiceRejectMintRequest(args RegistrarServiceRejectMintRequest) (*bind.EncodedChoice, error)
	RegistrarServiceSet(args RegistrarServiceSet) (*bind.EncodedChoice, error)
	RegistrarServiceSplitHolding(args RegistrarServiceSplitHolding) (*bind.EncodedChoice, error)
	RegistrarServiceTerminate(args RegistrarServiceTerminate) (*bind.EncodedChoice, error)
	RegistrarServiceTerminateEnforcementService(args RegistrarServiceTerminateEnforcementService) (*bind.EncodedChoice, error)
	RejectedBurnDelete(args RejectedBurnDelete) (*bind.EncodedChoice, error)
	RejectedEnforcementServiceRequestDelete(args RejectedEnforcementServiceRequestDelete) (*bind.EncodedChoice, error)
	RejectedHolderServiceRequestClean(args RejectedHolderServiceRequestClean) (*bind.EncodedChoice, error)
	RejectedHolderServiceRequestDelete(args RejectedHolderServiceRequestDelete) (*bind.EncodedChoice, error)
	RejectedMintDelete(args RejectedMintDelete) (*bind.EncodedChoice, error)
	RejectedProviderServiceRequestDelete(args RejectedProviderServiceRequestDelete) (*bind.EncodedChoice, error)
	RejectedRegistrarServiceRequestDelete(args RejectedRegistrarServiceRequestDelete) (*bind.EncodedChoice, error)
	TransferPreapprovalModify(args TransferPreapprovalModify) (*bind.EncodedChoice, error)
	TransferPreapprovalWithdraw(args TransferPreapprovalWithdraw) (*bind.EncodedChoice, error)
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

// AllocationFactoryAllocateInternal encodes parameters for the AllocationFactory_AllocateInternal choice.
func (e *encoder) AllocationFactoryAllocateInternal(args AllocationFactoryAllocateInternal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_AllocateInternal", args)
}

// AllocationFactoryInternalBurnMint encodes parameters for the AllocationFactory_InternalBurnMint choice.
func (e *encoder) AllocationFactoryInternalBurnMint(args AllocationFactoryInternalBurnMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_InternalBurnMint", args)
}

// AllocationFactoryOfferBurn encodes parameters for the AllocationFactory_OfferBurn choice.
func (e *encoder) AllocationFactoryOfferBurn(args AllocationFactoryOfferBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_OfferBurn", args)
}

// AllocationFactoryOfferMint encodes parameters for the AllocationFactory_OfferMint choice.
func (e *encoder) AllocationFactoryOfferMint(args AllocationFactoryOfferMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_OfferMint", args)
}

// AllocationFactoryRequestBurn encodes parameters for the AllocationFactory_RequestBurn choice.
func (e *encoder) AllocationFactoryRequestBurn(args AllocationFactoryRequestBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_RequestBurn", args)
}

// AllocationFactoryRequestMint encodes parameters for the AllocationFactory_RequestMint choice.
func (e *encoder) AllocationFactoryRequestMint(args AllocationFactoryRequestMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_RequestMint", args)
}

// AllocationFactoryTransferInternal encodes parameters for the AllocationFactory_TransferInternal choice.
func (e *encoder) AllocationFactoryTransferInternal(args AllocationFactoryTransferInternal) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AllocationFactory_TransferInternal", args)
}

// BurnOfferAccept encodes parameters for the BurnOffer_Accept choice.
func (e *encoder) BurnOfferAccept(args BurnOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Accept", args)
}

// BurnOfferCancel encodes parameters for the BurnOffer_Cancel choice.
func (e *encoder) BurnOfferCancel(args BurnOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnOffer_Cancel", args)
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

// BurnRequestReject encodes parameters for the BurnRequest_Reject choice.
func (e *encoder) BurnRequestReject(args BurnRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BurnRequest_Reject", args)
}

// EnforcementServiceRequestAccept encodes parameters for the EnforcementServiceRequest_Accept choice.
func (e *encoder) EnforcementServiceRequestAccept(args EnforcementServiceRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementServiceRequest_Accept", args)
}

// EnforcementServiceRequestCancel encodes parameters for the EnforcementServiceRequest_Cancel choice.
func (e *encoder) EnforcementServiceRequestCancel(args EnforcementServiceRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementServiceRequest_Cancel", args)
}

// EnforcementServiceRequestReject encodes parameters for the EnforcementServiceRequest_Reject choice.
func (e *encoder) EnforcementServiceRequestReject(args EnforcementServiceRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementServiceRequest_Reject", args)
}

// EnforcementServiceAcceptForceTransferRequest encodes parameters for the EnforcementService_AcceptForceTransferRequest choice.
func (e *encoder) EnforcementServiceAcceptForceTransferRequest(args EnforcementServiceAcceptForceTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementService_AcceptForceTransferRequest", args)
}

// EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization encodes parameters for the EnforcementService_AcceptForceTransferRequestWithSenderAuthorization choice.
func (e *encoder) EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization(args EnforcementServiceAcceptForceTransferRequestWithSenderAuthorization) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementService_AcceptForceTransferRequestWithSenderAuthorization", args)
}

// EnforcementServiceTerminate encodes parameters for the EnforcementService_Terminate choice.
func (e *encoder) EnforcementServiceTerminate(args EnforcementServiceTerminate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("EnforcementService_Terminate", args)
}

// ExecutedBurnDelete encodes parameters for the ExecutedBurn_Delete choice.
func (e *encoder) ExecutedBurnDelete(args ExecutedBurnDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedBurn_Delete", args)
}

// ExecutedMintDelete encodes parameters for the ExecutedMint_Delete choice.
func (e *encoder) ExecutedMintDelete(args ExecutedMintDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecutedMint_Delete", args)
}

// HolderServiceRequestAccept encodes parameters for the HolderServiceRequest_Accept choice.
func (e *encoder) HolderServiceRequestAccept(args HolderServiceRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderServiceRequest_Accept", args)
}

// HolderServiceRequestCancel encodes parameters for the HolderServiceRequest_Cancel choice.
func (e *encoder) HolderServiceRequestCancel(args HolderServiceRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderServiceRequest_Cancel", args)
}

// HolderServiceRequestClean encodes parameters for the HolderServiceRequest_Clean choice.
func (e *encoder) HolderServiceRequestClean(args HolderServiceRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderServiceRequest_Clean", args)
}

// HolderServiceRequestReject encodes parameters for the HolderServiceRequest_Reject choice.
func (e *encoder) HolderServiceRequestReject(args HolderServiceRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderServiceRequest_Reject", args)
}

// HolderServiceAcceptBurnOffer encodes parameters for the HolderService_AcceptBurnOffer choice.
func (e *encoder) HolderServiceAcceptBurnOffer(args HolderServiceAcceptBurnOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptBurnOffer", args)
}

// HolderServiceAcceptLockOffer encodes parameters for the HolderService_AcceptLockOffer choice.
func (e *encoder) HolderServiceAcceptLockOffer(args HolderServiceAcceptLockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptLockOffer", args)
}

// HolderServiceAcceptLockRequest encodes parameters for the HolderService_AcceptLockRequest choice.
func (e *encoder) HolderServiceAcceptLockRequest(args HolderServiceAcceptLockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptLockRequest", args)
}

// HolderServiceAcceptMintOffer encodes parameters for the HolderService_AcceptMintOffer choice.
func (e *encoder) HolderServiceAcceptMintOffer(args HolderServiceAcceptMintOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptMintOffer", args)
}

// HolderServiceAcceptTransferOffer encodes parameters for the HolderService_AcceptTransferOffer choice.
func (e *encoder) HolderServiceAcceptTransferOffer(args HolderServiceAcceptTransferOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptTransferOffer", args)
}

// HolderServiceAcceptTransferRequest encodes parameters for the HolderService_AcceptTransferRequest choice.
func (e *encoder) HolderServiceAcceptTransferRequest(args HolderServiceAcceptTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptTransferRequest", args)
}

// HolderServiceAcceptUnlockOffer encodes parameters for the HolderService_AcceptUnlockOffer choice.
func (e *encoder) HolderServiceAcceptUnlockOffer(args HolderServiceAcceptUnlockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptUnlockOffer", args)
}

// HolderServiceAcceptUnlockRequest encodes parameters for the HolderService_AcceptUnlockRequest choice.
func (e *encoder) HolderServiceAcceptUnlockRequest(args HolderServiceAcceptUnlockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_AcceptUnlockRequest", args)
}

// HolderServiceCancelBurnRequest encodes parameters for the HolderService_CancelBurnRequest choice.
func (e *encoder) HolderServiceCancelBurnRequest(args HolderServiceCancelBurnRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelBurnRequest", args)
}

// HolderServiceCancelEnforcementServiceRequest encodes parameters for the HolderService_CancelEnforcementServiceRequest choice.
func (e *encoder) HolderServiceCancelEnforcementServiceRequest(args HolderServiceCancelEnforcementServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelEnforcementServiceRequest", args)
}

// HolderServiceCancelForceTransferRequest encodes parameters for the HolderService_CancelForceTransferRequest choice.
func (e *encoder) HolderServiceCancelForceTransferRequest(args HolderServiceCancelForceTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelForceTransferRequest", args)
}

// HolderServiceCancelLockOffer encodes parameters for the HolderService_CancelLockOffer choice.
func (e *encoder) HolderServiceCancelLockOffer(args HolderServiceCancelLockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelLockOffer", args)
}

// HolderServiceCancelLockRequest encodes parameters for the HolderService_CancelLockRequest choice.
func (e *encoder) HolderServiceCancelLockRequest(args HolderServiceCancelLockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelLockRequest", args)
}

// HolderServiceCancelMintRequest encodes parameters for the HolderService_CancelMintRequest choice.
func (e *encoder) HolderServiceCancelMintRequest(args HolderServiceCancelMintRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelMintRequest", args)
}

// HolderServiceCancelTransferOffer encodes parameters for the HolderService_CancelTransferOffer choice.
func (e *encoder) HolderServiceCancelTransferOffer(args HolderServiceCancelTransferOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelTransferOffer", args)
}

// HolderServiceCancelTransferRequest encodes parameters for the HolderService_CancelTransferRequest choice.
func (e *encoder) HolderServiceCancelTransferRequest(args HolderServiceCancelTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelTransferRequest", args)
}

// HolderServiceCancelUnlockOffer encodes parameters for the HolderService_CancelUnlockOffer choice.
func (e *encoder) HolderServiceCancelUnlockOffer(args HolderServiceCancelUnlockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelUnlockOffer", args)
}

// HolderServiceCancelUnlockRequest encodes parameters for the HolderService_CancelUnlockRequest choice.
func (e *encoder) HolderServiceCancelUnlockRequest(args HolderServiceCancelUnlockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CancelUnlockRequest", args)
}

// HolderServiceClean encodes parameters for the HolderService_Clean choice.
func (e *encoder) HolderServiceClean(args HolderServiceClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_Clean", args)
}

// HolderServiceCreateAllocation encodes parameters for the HolderService_CreateAllocation choice.
func (e *encoder) HolderServiceCreateAllocation(args HolderServiceCreateAllocation) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_CreateAllocation", args)
}

// HolderServiceOfferLock encodes parameters for the HolderService_OfferLock choice.
func (e *encoder) HolderServiceOfferLock(args HolderServiceOfferLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_OfferLock", args)
}

// HolderServiceOfferTransfer encodes parameters for the HolderService_OfferTransfer choice.
func (e *encoder) HolderServiceOfferTransfer(args HolderServiceOfferTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_OfferTransfer", args)
}

// HolderServiceOfferUnlock encodes parameters for the HolderService_OfferUnlock choice.
func (e *encoder) HolderServiceOfferUnlock(args HolderServiceOfferUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_OfferUnlock", args)
}

// HolderServiceRejectAllocationRequest encodes parameters for the HolderService_RejectAllocationRequest choice.
func (e *encoder) HolderServiceRejectAllocationRequest(args HolderServiceRejectAllocationRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectAllocationRequest", args)
}

// HolderServiceRejectBurnOffer encodes parameters for the HolderService_RejectBurnOffer choice.
func (e *encoder) HolderServiceRejectBurnOffer(args HolderServiceRejectBurnOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectBurnOffer", args)
}

// HolderServiceRejectLockOffer encodes parameters for the HolderService_RejectLockOffer choice.
func (e *encoder) HolderServiceRejectLockOffer(args HolderServiceRejectLockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectLockOffer", args)
}

// HolderServiceRejectLockRequest encodes parameters for the HolderService_RejectLockRequest choice.
func (e *encoder) HolderServiceRejectLockRequest(args HolderServiceRejectLockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectLockRequest", args)
}

// HolderServiceRejectMintOffer encodes parameters for the HolderService_RejectMintOffer choice.
func (e *encoder) HolderServiceRejectMintOffer(args HolderServiceRejectMintOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectMintOffer", args)
}

// HolderServiceRejectTransferOffer encodes parameters for the HolderService_RejectTransferOffer choice.
func (e *encoder) HolderServiceRejectTransferOffer(args HolderServiceRejectTransferOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectTransferOffer", args)
}

// HolderServiceRejectTransferRequest encodes parameters for the HolderService_RejectTransferRequest choice.
func (e *encoder) HolderServiceRejectTransferRequest(args HolderServiceRejectTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectTransferRequest", args)
}

// HolderServiceRejectUnlockOffer encodes parameters for the HolderService_RejectUnlockOffer choice.
func (e *encoder) HolderServiceRejectUnlockOffer(args HolderServiceRejectUnlockOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectUnlockOffer", args)
}

// HolderServiceRejectUnlockRequest encodes parameters for the HolderService_RejectUnlockRequest choice.
func (e *encoder) HolderServiceRejectUnlockRequest(args HolderServiceRejectUnlockRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RejectUnlockRequest", args)
}

// HolderServiceRequestBurn encodes parameters for the HolderService_RequestBurn choice.
func (e *encoder) HolderServiceRequestBurn(args HolderServiceRequestBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestBurn", args)
}

// HolderServiceRequestEnforcementService encodes parameters for the HolderService_RequestEnforcementService choice.
func (e *encoder) HolderServiceRequestEnforcementService(args HolderServiceRequestEnforcementService) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestEnforcementService", args)
}

// HolderServiceRequestForceTransfer encodes parameters for the HolderService_RequestForceTransfer choice.
func (e *encoder) HolderServiceRequestForceTransfer(args HolderServiceRequestForceTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestForceTransfer", args)
}

// HolderServiceRequestLock encodes parameters for the HolderService_RequestLock choice.
func (e *encoder) HolderServiceRequestLock(args HolderServiceRequestLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestLock", args)
}

// HolderServiceRequestMint encodes parameters for the HolderService_RequestMint choice.
func (e *encoder) HolderServiceRequestMint(args HolderServiceRequestMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestMint", args)
}

// HolderServiceRequestTransfer encodes parameters for the HolderService_RequestTransfer choice.
func (e *encoder) HolderServiceRequestTransfer(args HolderServiceRequestTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestTransfer", args)
}

// HolderServiceRequestUnlock encodes parameters for the HolderService_RequestUnlock choice.
func (e *encoder) HolderServiceRequestUnlock(args HolderServiceRequestUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_RequestUnlock", args)
}

// HolderServiceTerminate encodes parameters for the HolderService_Terminate choice.
func (e *encoder) HolderServiceTerminate(args HolderServiceTerminate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HolderService_Terminate", args)
}

// MintOfferAccept encodes parameters for the MintOffer_Accept choice.
func (e *encoder) MintOfferAccept(args MintOfferAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Accept", args)
}

// MintOfferCancel encodes parameters for the MintOffer_Cancel choice.
func (e *encoder) MintOfferCancel(args MintOfferCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintOffer_Cancel", args)
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

// MintRequestReject encodes parameters for the MintRequest_Reject choice.
func (e *encoder) MintRequestReject(args MintRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("MintRequest_Reject", args)
}

// OperatorConfigurationGet encodes parameters for the OperatorConfiguration_Get choice.
func (e *encoder) OperatorConfigurationGet(args OperatorConfigurationGet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("OperatorConfiguration_Get", args)
}

// OperatorConfigurationModify encodes parameters for the OperatorConfiguration_Modify choice.
func (e *encoder) OperatorConfigurationModify(args OperatorConfigurationModify) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("OperatorConfiguration_Modify", args)
}

// ProviderConfigurationGet encodes parameters for the ProviderConfiguration_Get choice.
func (e *encoder) ProviderConfigurationGet(args ProviderConfigurationGet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderConfiguration_Get", args)
}

// ProviderServiceRequestAccept encodes parameters for the ProviderServiceRequest_Accept choice.
func (e *encoder) ProviderServiceRequestAccept(args ProviderServiceRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderServiceRequest_Accept", args)
}

// ProviderServiceRequestCancel encodes parameters for the ProviderServiceRequest_Cancel choice.
func (e *encoder) ProviderServiceRequestCancel(args ProviderServiceRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderServiceRequest_Cancel", args)
}

// ProviderServiceRequestReject encodes parameters for the ProviderServiceRequest_Reject choice.
func (e *encoder) ProviderServiceRequestReject(args ProviderServiceRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderServiceRequest_Reject", args)
}

// ProviderServiceAcceptHolderServiceRequest encodes parameters for the ProviderService_AcceptHolderServiceRequest choice.
func (e *encoder) ProviderServiceAcceptHolderServiceRequest(args ProviderServiceAcceptHolderServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_AcceptHolderServiceRequest", args)
}

// ProviderServiceAcceptRegistrarServiceRequest encodes parameters for the ProviderService_AcceptRegistrarServiceRequest choice.
func (e *encoder) ProviderServiceAcceptRegistrarServiceRequest(args ProviderServiceAcceptRegistrarServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_AcceptRegistrarServiceRequest", args)
}

// ProviderServiceArchiveAndCreateProviderConfiguration encodes parameters for the ProviderService_ArchiveAndCreateProviderConfiguration choice.
func (e *encoder) ProviderServiceArchiveAndCreateProviderConfiguration(args ProviderServiceArchiveAndCreateProviderConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_ArchiveAndCreateProviderConfiguration", args)
}

// ProviderServiceArchiveProviderConfiguration encodes parameters for the ProviderService_ArchiveProviderConfiguration choice.
func (e *encoder) ProviderServiceArchiveProviderConfiguration(args ProviderServiceArchiveProviderConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_ArchiveProviderConfiguration", args)
}

// ProviderServiceCreateProviderConfiguration encodes parameters for the ProviderService_CreateProviderConfiguration choice.
func (e *encoder) ProviderServiceCreateProviderConfiguration(args ProviderServiceCreateProviderConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_CreateProviderConfiguration", args)
}

// ProviderServiceRejectHolderServiceRequest encodes parameters for the ProviderService_RejectHolderServiceRequest choice.
func (e *encoder) ProviderServiceRejectHolderServiceRequest(args ProviderServiceRejectHolderServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_RejectHolderServiceRequest", args)
}

// ProviderServiceRejectRegistrarServiceRequest encodes parameters for the ProviderService_RejectRegistrarServiceRequest choice.
func (e *encoder) ProviderServiceRejectRegistrarServiceRequest(args ProviderServiceRejectRegistrarServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_RejectRegistrarServiceRequest", args)
}

// ProviderServiceTerminate encodes parameters for the ProviderService_Terminate choice.
func (e *encoder) ProviderServiceTerminate(args ProviderServiceTerminate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ProviderService_Terminate", args)
}

// RegistrarConfigurationGet encodes parameters for the RegistrarConfiguration_Get choice.
func (e *encoder) RegistrarConfigurationGet(args RegistrarConfigurationGet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarConfiguration_Get", args)
}

// RegistrarServiceRequestAccept encodes parameters for the RegistrarServiceRequest_Accept choice.
func (e *encoder) RegistrarServiceRequestAccept(args RegistrarServiceRequestAccept) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarServiceRequest_Accept", args)
}

// RegistrarServiceRequestCancel encodes parameters for the RegistrarServiceRequest_Cancel choice.
func (e *encoder) RegistrarServiceRequestCancel(args RegistrarServiceRequestCancel) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarServiceRequest_Cancel", args)
}

// RegistrarServiceRequestReject encodes parameters for the RegistrarServiceRequest_Reject choice.
func (e *encoder) RegistrarServiceRequestReject(args RegistrarServiceRequestReject) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarServiceRequest_Reject", args)
}

// RegistrarServiceAcceptBurnRequest encodes parameters for the RegistrarService_AcceptBurnRequest choice.
func (e *encoder) RegistrarServiceAcceptBurnRequest(args RegistrarServiceAcceptBurnRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_AcceptBurnRequest", args)
}

// RegistrarServiceAcceptEnforcementServiceRequest encodes parameters for the RegistrarService_AcceptEnforcementServiceRequest choice.
func (e *encoder) RegistrarServiceAcceptEnforcementServiceRequest(args RegistrarServiceAcceptEnforcementServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_AcceptEnforcementServiceRequest", args)
}

// RegistrarServiceAcceptForceTransferRequest encodes parameters for the RegistrarService_AcceptForceTransferRequest choice.
func (e *encoder) RegistrarServiceAcceptForceTransferRequest(args RegistrarServiceAcceptForceTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_AcceptForceTransferRequest", args)
}

// RegistrarServiceAcceptMintRequest encodes parameters for the RegistrarService_AcceptMintRequest choice.
func (e *encoder) RegistrarServiceAcceptMintRequest(args RegistrarServiceAcceptMintRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_AcceptMintRequest", args)
}

// RegistrarServiceArchiveAllocationFactory encodes parameters for the RegistrarService_ArchiveAllocationFactory choice.
func (e *encoder) RegistrarServiceArchiveAllocationFactory(args RegistrarServiceArchiveAllocationFactory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveAllocationFactory", args)
}

// RegistrarServiceArchiveAndCreateInstrumentConfiguration encodes parameters for the RegistrarService_ArchiveAndCreateInstrumentConfiguration choice.
func (e *encoder) RegistrarServiceArchiveAndCreateInstrumentConfiguration(args RegistrarServiceArchiveAndCreateInstrumentConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveAndCreateInstrumentConfiguration", args)
}

// RegistrarServiceArchiveAndCreateRegistrarConfiguration encodes parameters for the RegistrarService_ArchiveAndCreateRegistrarConfiguration choice.
func (e *encoder) RegistrarServiceArchiveAndCreateRegistrarConfiguration(args RegistrarServiceArchiveAndCreateRegistrarConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveAndCreateRegistrarConfiguration", args)
}

// RegistrarServiceArchiveInstrumentConfiguration encodes parameters for the RegistrarService_ArchiveInstrumentConfiguration choice.
func (e *encoder) RegistrarServiceArchiveInstrumentConfiguration(args RegistrarServiceArchiveInstrumentConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveInstrumentConfiguration", args)
}

// RegistrarServiceArchiveRegistrarConfiguration encodes parameters for the RegistrarService_ArchiveRegistrarConfiguration choice.
func (e *encoder) RegistrarServiceArchiveRegistrarConfiguration(args RegistrarServiceArchiveRegistrarConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveRegistrarConfiguration", args)
}

// RegistrarServiceArchiveTransferRule encodes parameters for the RegistrarService_ArchiveTransferRule choice.
func (e *encoder) RegistrarServiceArchiveTransferRule(args RegistrarServiceArchiveTransferRule) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ArchiveTransferRule", args)
}

// RegistrarServiceCancelBurnOffer encodes parameters for the RegistrarService_CancelBurnOffer choice.
func (e *encoder) RegistrarServiceCancelBurnOffer(args RegistrarServiceCancelBurnOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CancelBurnOffer", args)
}

// RegistrarServiceCancelMintOffer encodes parameters for the RegistrarService_CancelMintOffer choice.
func (e *encoder) RegistrarServiceCancelMintOffer(args RegistrarServiceCancelMintOffer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CancelMintOffer", args)
}

// RegistrarServiceCreateAllocationFactory encodes parameters for the RegistrarService_CreateAllocationFactory choice.
func (e *encoder) RegistrarServiceCreateAllocationFactory(args RegistrarServiceCreateAllocationFactory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CreateAllocationFactory", args)
}

// RegistrarServiceCreateInstrumentConfiguration encodes parameters for the RegistrarService_CreateInstrumentConfiguration choice.
func (e *encoder) RegistrarServiceCreateInstrumentConfiguration(args RegistrarServiceCreateInstrumentConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CreateInstrumentConfiguration", args)
}

// RegistrarServiceCreateRegistrarConfiguration encodes parameters for the RegistrarService_CreateRegistrarConfiguration choice.
func (e *encoder) RegistrarServiceCreateRegistrarConfiguration(args RegistrarServiceCreateRegistrarConfiguration) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CreateRegistrarConfiguration", args)
}

// RegistrarServiceCreateTransferRule encodes parameters for the RegistrarService_CreateTransferRule choice.
func (e *encoder) RegistrarServiceCreateTransferRule(args RegistrarServiceCreateTransferRule) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_CreateTransferRule", args)
}

// RegistrarServiceDeleteExecutedBurn encodes parameters for the RegistrarService_DeleteExecutedBurn choice.
func (e *encoder) RegistrarServiceDeleteExecutedBurn(args RegistrarServiceDeleteExecutedBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedBurn", args)
}

// RegistrarServiceDeleteExecutedLock encodes parameters for the RegistrarService_DeleteExecutedLock choice.
func (e *encoder) RegistrarServiceDeleteExecutedLock(args RegistrarServiceDeleteExecutedLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedLock", args)
}

// RegistrarServiceDeleteExecutedMint encodes parameters for the RegistrarService_DeleteExecutedMint choice.
func (e *encoder) RegistrarServiceDeleteExecutedMint(args RegistrarServiceDeleteExecutedMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedMint", args)
}

// RegistrarServiceDeleteExecutedTransfer encodes parameters for the RegistrarService_DeleteExecutedTransfer choice.
func (e *encoder) RegistrarServiceDeleteExecutedTransfer(args RegistrarServiceDeleteExecutedTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedTransfer", args)
}

// RegistrarServiceDeleteExecutedTransfers encodes parameters for the RegistrarService_DeleteExecutedTransfers choice.
func (e *encoder) RegistrarServiceDeleteExecutedTransfers(args RegistrarServiceDeleteExecutedTransfers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedTransfers", args)
}

// RegistrarServiceDeleteExecutedUnlock encodes parameters for the RegistrarService_DeleteExecutedUnlock choice.
func (e *encoder) RegistrarServiceDeleteExecutedUnlock(args RegistrarServiceDeleteExecutedUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteExecutedUnlock", args)
}

// RegistrarServiceDeleteFailedBurn encodes parameters for the RegistrarService_DeleteFailedBurn choice.
func (e *encoder) RegistrarServiceDeleteFailedBurn(args RegistrarServiceDeleteFailedBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedBurn", args)
}

// RegistrarServiceDeleteFailedLock encodes parameters for the RegistrarService_DeleteFailedLock choice.
func (e *encoder) RegistrarServiceDeleteFailedLock(args RegistrarServiceDeleteFailedLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedLock", args)
}

// RegistrarServiceDeleteFailedMint encodes parameters for the RegistrarService_DeleteFailedMint choice.
func (e *encoder) RegistrarServiceDeleteFailedMint(args RegistrarServiceDeleteFailedMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedMint", args)
}

// RegistrarServiceDeleteFailedTransfer encodes parameters for the RegistrarService_DeleteFailedTransfer choice.
func (e *encoder) RegistrarServiceDeleteFailedTransfer(args RegistrarServiceDeleteFailedTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedTransfer", args)
}

// RegistrarServiceDeleteFailedTransfers encodes parameters for the RegistrarService_DeleteFailedTransfers choice.
func (e *encoder) RegistrarServiceDeleteFailedTransfers(args RegistrarServiceDeleteFailedTransfers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedTransfers", args)
}

// RegistrarServiceDeleteFailedUnlock encodes parameters for the RegistrarService_DeleteFailedUnlock choice.
func (e *encoder) RegistrarServiceDeleteFailedUnlock(args RegistrarServiceDeleteFailedUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteFailedUnlock", args)
}

// RegistrarServiceDeleteRejectedBurn encodes parameters for the RegistrarService_DeleteRejectedBurn choice.
func (e *encoder) RegistrarServiceDeleteRejectedBurn(args RegistrarServiceDeleteRejectedBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedBurn", args)
}

// RegistrarServiceDeleteRejectedLock encodes parameters for the RegistrarService_DeleteRejectedLock choice.
func (e *encoder) RegistrarServiceDeleteRejectedLock(args RegistrarServiceDeleteRejectedLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedLock", args)
}

// RegistrarServiceDeleteRejectedMint encodes parameters for the RegistrarService_DeleteRejectedMint choice.
func (e *encoder) RegistrarServiceDeleteRejectedMint(args RegistrarServiceDeleteRejectedMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedMint", args)
}

// RegistrarServiceDeleteRejectedTransfer encodes parameters for the RegistrarService_DeleteRejectedTransfer choice.
func (e *encoder) RegistrarServiceDeleteRejectedTransfer(args RegistrarServiceDeleteRejectedTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedTransfer", args)
}

// RegistrarServiceDeleteRejectedTransfers encodes parameters for the RegistrarService_DeleteRejectedTransfers choice.
func (e *encoder) RegistrarServiceDeleteRejectedTransfers(args RegistrarServiceDeleteRejectedTransfers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedTransfers", args)
}

// RegistrarServiceDeleteRejectedUnlock encodes parameters for the RegistrarService_DeleteRejectedUnlock choice.
func (e *encoder) RegistrarServiceDeleteRejectedUnlock(args RegistrarServiceDeleteRejectedUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_DeleteRejectedUnlock", args)
}

// RegistrarServiceExecuteAcceptedBurn encodes parameters for the RegistrarService_ExecuteAcceptedBurn choice.
func (e *encoder) RegistrarServiceExecuteAcceptedBurn(args RegistrarServiceExecuteAcceptedBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedBurn", args)
}

// RegistrarServiceExecuteAcceptedForceTransfer encodes parameters for the RegistrarService_ExecuteAcceptedForceTransfer choice.
func (e *encoder) RegistrarServiceExecuteAcceptedForceTransfer(args RegistrarServiceExecuteAcceptedForceTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedForceTransfer", args)
}

// RegistrarServiceExecuteAcceptedLock encodes parameters for the RegistrarService_ExecuteAcceptedLock choice.
func (e *encoder) RegistrarServiceExecuteAcceptedLock(args RegistrarServiceExecuteAcceptedLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedLock", args)
}

// RegistrarServiceExecuteAcceptedMint encodes parameters for the RegistrarService_ExecuteAcceptedMint choice.
func (e *encoder) RegistrarServiceExecuteAcceptedMint(args RegistrarServiceExecuteAcceptedMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedMint", args)
}

// RegistrarServiceExecuteAcceptedTransfer encodes parameters for the RegistrarService_ExecuteAcceptedTransfer choice.
func (e *encoder) RegistrarServiceExecuteAcceptedTransfer(args RegistrarServiceExecuteAcceptedTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedTransfer", args)
}

// RegistrarServiceExecuteAcceptedUnlock encodes parameters for the RegistrarService_ExecuteAcceptedUnlock choice.
func (e *encoder) RegistrarServiceExecuteAcceptedUnlock(args RegistrarServiceExecuteAcceptedUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_ExecuteAcceptedUnlock", args)
}

// RegistrarServiceFailAcceptedBurn encodes parameters for the RegistrarService_FailAcceptedBurn choice.
func (e *encoder) RegistrarServiceFailAcceptedBurn(args RegistrarServiceFailAcceptedBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedBurn", args)
}

// RegistrarServiceFailAcceptedForceTransfer encodes parameters for the RegistrarService_FailAcceptedForceTransfer choice.
func (e *encoder) RegistrarServiceFailAcceptedForceTransfer(args RegistrarServiceFailAcceptedForceTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedForceTransfer", args)
}

// RegistrarServiceFailAcceptedLock encodes parameters for the RegistrarService_FailAcceptedLock choice.
func (e *encoder) RegistrarServiceFailAcceptedLock(args RegistrarServiceFailAcceptedLock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedLock", args)
}

// RegistrarServiceFailAcceptedMint encodes parameters for the RegistrarService_FailAcceptedMint choice.
func (e *encoder) RegistrarServiceFailAcceptedMint(args RegistrarServiceFailAcceptedMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedMint", args)
}

// RegistrarServiceFailAcceptedTransfer encodes parameters for the RegistrarService_FailAcceptedTransfer choice.
func (e *encoder) RegistrarServiceFailAcceptedTransfer(args RegistrarServiceFailAcceptedTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedTransfer", args)
}

// RegistrarServiceFailAcceptedUnlock encodes parameters for the RegistrarService_FailAcceptedUnlock choice.
func (e *encoder) RegistrarServiceFailAcceptedUnlock(args RegistrarServiceFailAcceptedUnlock) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_FailAcceptedUnlock", args)
}

// RegistrarServiceMergeHolding encodes parameters for the RegistrarService_MergeHolding choice.
func (e *encoder) RegistrarServiceMergeHolding(args RegistrarServiceMergeHolding) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_MergeHolding", args)
}

// RegistrarServiceOfferBurn encodes parameters for the RegistrarService_OfferBurn choice.
func (e *encoder) RegistrarServiceOfferBurn(args RegistrarServiceOfferBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_OfferBurn", args)
}

// RegistrarServiceOfferMint encodes parameters for the RegistrarService_OfferMint choice.
func (e *encoder) RegistrarServiceOfferMint(args RegistrarServiceOfferMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_OfferMint", args)
}

// RegistrarServiceRejectBurnRequest encodes parameters for the RegistrarService_RejectBurnRequest choice.
func (e *encoder) RegistrarServiceRejectBurnRequest(args RegistrarServiceRejectBurnRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_RejectBurnRequest", args)
}

// RegistrarServiceRejectEnforcementServiceRequest encodes parameters for the RegistrarService_RejectEnforcementServiceRequest choice.
func (e *encoder) RegistrarServiceRejectEnforcementServiceRequest(args RegistrarServiceRejectEnforcementServiceRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_RejectEnforcementServiceRequest", args)
}

// RegistrarServiceRejectForceTransferRequest encodes parameters for the RegistrarService_RejectForceTransferRequest choice.
func (e *encoder) RegistrarServiceRejectForceTransferRequest(args RegistrarServiceRejectForceTransferRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_RejectForceTransferRequest", args)
}

// RegistrarServiceRejectMintRequest encodes parameters for the RegistrarService_RejectMintRequest choice.
func (e *encoder) RegistrarServiceRejectMintRequest(args RegistrarServiceRejectMintRequest) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_RejectMintRequest", args)
}

// RegistrarServiceSet encodes parameters for the RegistrarService_Set choice.
func (e *encoder) RegistrarServiceSet(args RegistrarServiceSet) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_Set", args)
}

// RegistrarServiceSplitHolding encodes parameters for the RegistrarService_SplitHolding choice.
func (e *encoder) RegistrarServiceSplitHolding(args RegistrarServiceSplitHolding) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_SplitHolding", args)
}

// RegistrarServiceTerminate encodes parameters for the RegistrarService_Terminate choice.
func (e *encoder) RegistrarServiceTerminate(args RegistrarServiceTerminate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_Terminate", args)
}

// RegistrarServiceTerminateEnforcementService encodes parameters for the RegistrarService_TerminateEnforcementService choice.
func (e *encoder) RegistrarServiceTerminateEnforcementService(args RegistrarServiceTerminateEnforcementService) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RegistrarService_TerminateEnforcementService", args)
}

// RejectedBurnDelete encodes parameters for the RejectedBurn_Delete choice.
func (e *encoder) RejectedBurnDelete(args RejectedBurnDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedBurn_Delete", args)
}

// RejectedEnforcementServiceRequestDelete encodes parameters for the RejectedEnforcementServiceRequest_Delete choice.
func (e *encoder) RejectedEnforcementServiceRequestDelete(args RejectedEnforcementServiceRequestDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedEnforcementServiceRequest_Delete", args)
}

// RejectedHolderServiceRequestClean encodes parameters for the RejectedHolderServiceRequest_Clean choice.
func (e *encoder) RejectedHolderServiceRequestClean(args RejectedHolderServiceRequestClean) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedHolderServiceRequest_Clean", args)
}

// RejectedHolderServiceRequestDelete encodes parameters for the RejectedHolderServiceRequest_Delete choice.
func (e *encoder) RejectedHolderServiceRequestDelete(args RejectedHolderServiceRequestDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedHolderServiceRequest_Delete", args)
}

// RejectedMintDelete encodes parameters for the RejectedMint_Delete choice.
func (e *encoder) RejectedMintDelete(args RejectedMintDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedMint_Delete", args)
}

// RejectedProviderServiceRequestDelete encodes parameters for the RejectedProviderServiceRequest_Delete choice.
func (e *encoder) RejectedProviderServiceRequestDelete(args RejectedProviderServiceRequestDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedProviderServiceRequest_Delete", args)
}

// RejectedRegistrarServiceRequestDelete encodes parameters for the RejectedRegistrarServiceRequest_Delete choice.
func (e *encoder) RejectedRegistrarServiceRequestDelete(args RejectedRegistrarServiceRequestDelete) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RejectedRegistrarServiceRequest_Delete", args)
}

// TransferPreapprovalModify encodes parameters for the TransferPreapproval_Modify choice.
func (e *encoder) TransferPreapprovalModify(args TransferPreapprovalModify) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferPreapproval_Modify", args)
}

// TransferPreapprovalWithdraw encodes parameters for the TransferPreapproval_Withdraw choice.
func (e *encoder) TransferPreapprovalWithdraw(args TransferPreapprovalWithdraw) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("TransferPreapproval_Withdraw", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
