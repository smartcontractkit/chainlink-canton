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

const PackageID = "3fd403d741c2ea5312ac0f94ef5461ea209f0a933a96322618dd477a3510e442"
const SDKVersion = "3.4.8"

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

	// Check if the type has a toMap method
	type mapper interface {
		toMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.toMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// CoinHolding is a Template type
type CoinHolding struct {
	View   HoldingView `json:"view"`
	Issuer PARTY       `json:"issuer"`
}

// GetTemplateID returns the template ID for this template
func (t CoinHolding) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Holding", "CoinHolding")
}

// CreateCommand returns a CreateCommand for this template
func (t CoinHolding) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["view"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.View).(mapper); ok {
			return m.toMap()
		}
		return t.View
	}()

	args["issuer"] = t.Issuer.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CoinHolding using JsonCodec
func (t CoinHolding) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CoinHolding using JsonCodec
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
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CoinHolding contract
func (t CoinHolding) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Holding", "CoinHolding"),
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
	Meta         Metadata     `json:"meta"`
}

// GetTemplateID returns the template ID for this template
func (t CoinRegistry) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "CoinRegistry")
}

// CreateCommand returns a CreateCommand for this template
func (t CoinRegistry) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["issuer"] = t.Issuer.ToMap()

	args["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

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

// MarshalJSON implements custom JSON marshaling for CoinRegistry using JsonCodec
func (t CoinRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CoinRegistry using JsonCodec
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
		Arguments:  map[string]interface{}{},
	}
}

// TransferFactoryTransfer exercises the TransferFactory_Transfer choice on this CoinRegistry contract via the ITransferFactory interface
func (t CoinRegistry) TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", "55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281", "Splice.Api.Token.TransferInstructionV1", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_Transfer",
		Arguments:  argsToMap(args),
	}
}

// TransferFactoryPublicFetch exercises the TransferFactory_PublicFetch choice on this CoinRegistry contract via the ITransferFactory interface
func (t CoinRegistry) TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "TransferFactory"),
		ContractID: contractID,
		Choice:     "TransferFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryPublicFetch exercises the BurnMintFactory_PublicFetch choice on this CoinRegistry contract via the IBurnMintFactory interface
func (t CoinRegistry) BurnMintFactoryPublicFetch(contractID string, args BurnMintFactoryPublicFetch) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "BurnMintFactory"),
		ContractID: contractID,
		Choice:     "BurnMintFactory_PublicFetch",
		Arguments:  argsToMap(args),
	}
}

// BurnMintFactoryBurnMint exercises the BurnMintFactory_BurnMint choice on this CoinRegistry contract via the IBurnMintFactory interface
func (t CoinRegistry) BurnMintFactoryBurnMint(contractID string, args BurnMintFactoryBurnMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("%s:%s:%s", "9cc2cbc838ef38dc2c7f34014c9c452bcf71b8e2a4f939235fc0b5d0924b185e", "Splice.Api.Token.BurnMintV1", "BurnMintFactory"),
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

// GetTemplateID returns the template ID for this template
func (t CoinTransferInstruction) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "CoinTransferInstruction")
}

// CreateCommand returns a CreateCommand for this template
func (t CoinTransferInstruction) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["holding"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Holding).(mapper); ok {
			return m.toMap()
		}
		return t.Holding
	}()

	args["newOwner"] = t.NewOwner.ToMap()

	args["requestedAt"] = t.RequestedAt

	args["executeBefore"] = t.ExecuteBefore

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CoinTransferInstruction using JsonCodec
func (t CoinTransferInstruction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CoinTransferInstruction using JsonCodec
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
		Arguments:  map[string]interface{}{},
	}
}

// TransferInstructionAccept exercises the TransferInstruction_Accept choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Accept",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionReject exercises the TransferInstruction_Reject choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Reject",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionWithdraw exercises the TransferInstruction_Withdraw choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),
		ContractID: contractID,
		Choice:     "TransferInstruction_Withdraw",
		Arguments:  argsToMap(args),
	}
}

// TransferInstructionUpdate exercises the TransferInstruction_Update choice on this CoinTransferInstruction contract via the ITransferInstruction interface
func (t CoinTransferInstruction) TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Transfer", "TransferInstruction"),
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

// GetTemplateID returns the template ID for this template
func (t MintPreapproval) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintPreapproval")
}

// CreateCommand returns a CreateCommand for this template
func (t MintPreapproval) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["receiver"] = t.Receiver.ToMap()

	args["sender"] = t.Sender.ToMap()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MintPreapproval using JsonCodec
func (t MintPreapproval) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MintPreapproval using JsonCodec
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
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this MintPreapproval contract
func (t MintPreapproval) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintPreapproval"),
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

// toMap converts MintPreapprovalMint to a map for DAML arguments
func (t MintPreapprovalMint) toMap() map[string]interface{} {
	return map[string]interface{}{

		"issuer": t.Issuer.ToMap(),
		"view": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.View).(mapper); ok {
				return m.toMap()
			}
			return t.View
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for MintPreapprovalMint using JsonCodec
func (t MintPreapprovalMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MintPreapprovalMint using JsonCodec
func (t *MintPreapprovalMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MintRole is a Template type
type MintRole struct {
	Issuer   PARTY  `json:"issuer"`
	Minter   PARTY  `json:"minter"`
	Registry string `json:"registry"`
}

// GetTemplateID returns the template ID for this template
func (t MintRole) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintRole")
}

// CreateCommand returns a CreateCommand for this template
func (t MintRole) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["issuer"] = t.Issuer.ToMap()

	args["minter"] = t.Minter.ToMap()

	args["registry"] = string(t.Registry)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for MintRole using JsonCodec
func (t MintRole) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MintRole using JsonCodec
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
		Arguments:  map[string]interface{}{},
	}
}

// MintRoleMint exercises the MintRole_Mint choice on this MintRole contract
func (t MintRole) MintRoleMint(contractID string, args MintRoleMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "Coin.Registry", "MintRole"),
		ContractID: contractID,
		Choice:     "MintRole_Mint",
		Arguments:  argsToMap(args),
	}
}

// MintRoleMint is a Record type
type MintRoleMint struct {
	InstrumentId InstrumentId `json:"instrumentId"`
	Outputs      LIST         `json:"outputs"`
}

// toMap converts MintRoleMint to a map for DAML arguments
func (t MintRoleMint) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"outputs": t.Outputs,
	}
}

// MarshalJSON implements custom JSON marshaling for MintRoleMint using JsonCodec
func (t MintRoleMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MintRoleMint using JsonCodec
func (t *MintRoleMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Transfer is a Record type
type Transfer struct {
	To PARTY `json:"to"`
}

// toMap converts Transfer to a map for DAML arguments
func (t Transfer) toMap() map[string]interface{} {
	return map[string]interface{}{

		"to": t.To.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for Transfer using JsonCodec
func (t Transfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for Transfer using JsonCodec
func (t *Transfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
