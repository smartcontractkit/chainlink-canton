package coin

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/noders-team/go-daml/pkg/codec"
	"github.com/noders-team/go-daml/pkg/model"
	. "github.com/noders-team/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageID = "831563cda769e12eaf47360bd4a44d63b108619af4c0965b0d917b7ebdab4599"
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
	View HoldingView `json:"view"`

	Issuer PARTY `json:"issuer"`
}

// GetTemplateID returns the template ID for this template
func (t CoinHolding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Holding", "CoinHolding")
}

// CreateCommand returns a CreateCommand for this template
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
func (t CoinHolding) Transfer(contractID string, args Transfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Holding", "CoinHolding"),

		ContractID: contractID,
		Choice:     "Transfer",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CoinHolding contract
func (t CoinHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Holding", "CoinHolding"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// Verify interface implementations for CoinHolding

var _ IHolding = (*CoinHolding)(nil)

// CoinRegistry is a Template type
type CoinRegistry struct {
	Issuer PARTY `json:"issuer"`

	InstrumentId InstrumentId `json:"instrumentId"`

	InstanceId TEXT `json:"instanceId"`

	Meta Metadata `json:"meta"`
}

// GetTemplateID returns the template ID for this template
func (t CoinRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "CoinRegistry")
}

// CreateCommand returns a CreateCommand for this template
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
func (t CoinRegistry) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "CoinRegistry"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this CoinRegistry contract via the ITransferFactory interface
func (t CoinRegistry) TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "TransferFactory"),

		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",

		Arguments: argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this CoinRegistry contract via the ITransferFactory interface
func (t CoinRegistry) TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "TransferFactory"),

		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",

		Arguments: argsToMap(args),
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this CoinRegistry contract via the IBurnMintFactory interface
func (t CoinRegistry) BurnMintFactoryPublicFetch(contractID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "BurnMintFactory"),

		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",

		Arguments: argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this CoinRegistry contract via the IBurnMintFactory interface
func (t CoinRegistry) BurnMintFactoryBurnMint(contractID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "BurnMintFactory"),

		ContractID: contractID,
		Choice:     "BurnMintFactory_BurnMint",

		Arguments: argsToMap(args),
	}
}

// Verify interface implementations for CoinRegistry

var _ ITransferFactory = (*CoinRegistry)(nil)

var _ IBurnMintFactory = (*CoinRegistry)(nil)

// CoinTransferInstruction is a Template type
type CoinTransferInstruction struct {
	Holding CoinHolding `json:"holding"`

	NewOwner PARTY `json:"newOwner"`

	RequestedAt TIMESTAMP `json:"requestedAt"`

	ExecuteBefore TIMESTAMP `json:"executeBefore"`
}

// GetTemplateID returns the template ID for this template
func (t CoinTransferInstruction) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "CoinTransferInstruction")
}

// CreateCommand returns a CreateCommand for this template
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
func (t CoinTransferInstruction) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "CoinTransferInstruction"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),

		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",

		Arguments: argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),

		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",

		Arguments: argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),

		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",

		Arguments: argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),

		ContractID: contractID,
		Choice:     "TransferInstruction_Update",

		Arguments: argsToMap(args),
	}
}

// Verify interface implementations for CoinTransferInstruction

var _ ITransferInstruction = (*CoinTransferInstruction)(nil)

// MintPreapproval is a Template type
type MintPreapproval struct {
	Receiver PARTY `json:"receiver"`

	Sender PARTY `json:"sender"`
}

// GetTemplateID returns the template ID for this template
func (t MintPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintPreapproval")
}

// CreateCommand returns a CreateCommand for this template
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
func (t MintPreapproval) MintPreapprovalMint(contractID string, args MintPreapprovalMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintPreapproval"),

		ContractID: contractID,
		Choice:     "MintPreapproval_Mint",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintPreapproval contract
func (t MintPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintPreapproval"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// MintPreapprovalMint is a Record type
type MintPreapprovalMint struct {
	Issuer PARTY `json:"issuer"`

	View HoldingView `json:"view"`
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
	Issuer PARTY `json:"issuer"`

	Minter PARTY `json:"minter"`

	Registry CONTRACT_ID `json:"registry"`
}

// GetTemplateID returns the template ID for this template
func (t MintRole) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintRole")
}

// CreateCommand returns a CreateCommand for this template
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
func (t MintRole) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintRole"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// MintRoleMint exercises the MintRole_Mint choice on this MintRole contract
func (t MintRole) MintRoleMint(contractID string, args MintRoleMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintRole"),

		ContractID: contractID,
		Choice:     "MintRole_Mint",

		Arguments: argsToMap(args),
	}
}

// MintRoleMint is a Record type
type MintRoleMint struct {
	InstrumentId InstrumentId `json:"instrumentId"`

	Outputs []BurnMintOutput `json:"outputs"`
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
