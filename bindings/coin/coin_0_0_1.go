package coin

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

const PackageName = "coin"
const SDKVersion = "3.4.10"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
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

// CoinHolding is a Template type
type CoinHolding struct {
	View   HoldingView `json:"view"`
	Issuer PARTY       `json:"issuer"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["view"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.View).(mapper); ok {
			return m.toMap()
		}
		return t.View
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CoinHolding) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["view"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.View).(mapper); ok {
			return m.toMap()
		}
		return t.View
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CoinHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CoinHolding) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// Archive exercises the Archive choice on this CoinHolding contract
// This method uses the package name in the template ID
func (t CoinHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Holding", "CoinHolding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinHolding) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Holding", "CoinHolding"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Verify interface implementations for CoinHolding

var _ IHolding = (*CoinHolding)(nil)

// CoinRegistry is a Template type
type CoinRegistry struct {
	Issuer       PARTY        `json:"issuer"`
	InstrumentId InstrumentId `json:"instrumentId"`
	InstanceId   TEXT         `json:"instanceId"`
	Meta         Metadata     `json:"meta"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.toMap()
		}
		return t.Meta
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CoinRegistry) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["meta"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.toMap()
		}
		return t.Meta
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CoinRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CoinRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CoinRegistry

// Archive exercises the Archive choice on this CoinRegistry contract
// This method uses the package name in the template ID
func (t CoinRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "CoinRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinRegistry) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "CoinRegistry"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this CoinRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryTransferWithPackageID exercises the TransferFactory_Transfer choice using the provided package ID instead of package name
func (t CoinRegistry) TransferFactoryTransferWithPackageID(contractID string, packageID string, args TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this CoinRegistry contract via the ITransferFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetchWithPackageID exercises the TransferFactory_PublicFetch choice using the provided package ID instead of package name
func (t CoinRegistry) TransferFactoryPublicFetchWithPackageID(contractID string, packageID string, args TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this CoinRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) BurnMintFactoryPublicFetch(contractID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetchWithPackageID exercises the BurnMintFactory_PublicFetch choice using the provided package ID instead of package name
func (t CoinRegistry) BurnMintFactoryPublicFetchWithPackageID(contractID string, packageID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this CoinRegistry contract via the IBurnMintFactory interface
// This method uses the package name in the template ID
func (t CoinRegistry) BurnMintFactoryBurnMint(contractID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMintWithPackageID exercises the BurnMintFactory_BurnMint choice using the provided package ID instead of package name
func (t CoinRegistry) BurnMintFactoryBurnMintWithPackageID(contractID string, packageID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CoinRegistry

var _ ITransferFactory = (*CoinRegistry)(nil)

var _ IBurnMintFactory = (*CoinRegistry)(nil)

// CoinTransferInstruction is a Template type
type CoinTransferInstruction struct {
	Holding       CoinHolding `json:"holding"`
	NewOwner      PARTY       `json:"newOwner"`
	RequestedAt   TIMESTAMP   `json:"requestedAt"`
	ExecuteBefore TIMESTAMP   `json:"executeBefore"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holding"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Holding).(mapper); ok {
			return m.toMap()
		}
		return t.Holding
	}()

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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holding"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Holding).(mapper); ok {
			return m.toMap()
		}
		return t.Holding
	}()

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
	return jsonCodec.Marshall(t)
}

func (t *CoinTransferInstruction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CoinTransferInstruction

// Archive exercises the Archive choice on this CoinTransferInstruction contract
// This method uses the package name in the template ID
func (t CoinTransferInstruction) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "CoinTransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CoinTransferInstruction) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "CoinTransferInstruction"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionAcceptWithPackageID exercises the TransferInstruction_Accept choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionAcceptWithPackageID(contractID string, packageID string, args TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionRejectWithPackageID exercises the TransferInstruction_Reject choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionRejectWithPackageID(contractID string, packageID string, args TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdrawWithPackageID exercises the TransferInstruction_Withdraw choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionWithdrawWithPackageID(contractID string, packageID string, args TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this CoinTransferInstruction contract via the ITransferInstruction interface
// This method uses the package name in the template ID
func (t CoinTransferInstruction) TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdateWithPackageID exercises the TransferInstruction_Update choice using the provided package ID instead of package name
func (t CoinTransferInstruction) TransferInstructionUpdateWithPackageID(contractID string, packageID string, args TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Update",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CoinTransferInstruction

var _ ITransferInstruction = (*CoinTransferInstruction)(nil)

// MintPreapproval is a Template type
type MintPreapproval struct {
	Receiver PARTY `json:"receiver"`
	Sender   PARTY `json:"sender"`
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
	args := make(map[string]interface{})

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
	args := make(map[string]interface{})

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
	return jsonCodec.Marshall(t)
}

func (t *MintPreapproval) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintPreapproval) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintPreapproval"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// MintPreapprovalMint is a Record type
type MintPreapprovalMint struct {
	Issuer PARTY       `json:"issuer"`
	View   HoldingView `json:"view"`
}

// ToMap converts MintPreapprovalMint to a map for DAML arguments
func (t MintPreapprovalMint) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["issuer"] = t.Issuer.ToMap()

	m["view"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.View).(mapper); ok {
			return m.toMap()
		}
		return t.View
	}()

	return m
}

func (t MintPreapprovalMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MintPreapprovalMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MintRole is a Template type
type MintRole struct {
	Issuer   PARTY       `json:"issuer"`
	Minter   PARTY       `json:"minter"`
	Registry CONTRACT_ID `json:"registry"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minter"] = t.Minter.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registry"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Registry).(mapper); ok {
			return m.toMap()
		}
		return t.Registry
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MintRole) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["issuer"] = t.Issuer.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minter"] = t.Minter.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registry"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Registry).(mapper); ok {
			return m.toMap()
		}
		return t.Registry
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MintRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MintRole) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MintRole

// Archive exercises the Archive choice on this MintRole contract
// This method uses the package name in the template ID
func (t MintRole) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MintRole) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
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
	InstrumentId InstrumentId     `json:"instrumentId"`
	Outputs      []BurnMintOutput `json:"outputs"`
}

// ToMap converts MintRoleMint to a map for DAML arguments
func (t MintRoleMint) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["outputs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Outputs))
		for _, e := range t.Outputs {
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

func (t MintRoleMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MintRoleMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Transfer is a Record type
type Transfer struct {
	To PARTY `json:"to"`
}

// ToMap converts Transfer to a map for DAML arguments
func (t Transfer) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["to"] = t.To.ToMap()

	return m
}

func (t Transfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *Transfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
