package tokenadmin

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	api "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
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
	PackageName = "link-token-admin"
	PackageID   = "24c448fcf1c5bbe4657bfde6faf79991e2b7b7e44a8efa118dcbde48a93bd15a"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	IndefiniteTransferHours = types.INT64(2160)
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

// ApproveAllocate is a Record type
type ApproveAllocate struct {
	Receiver        types.PARTY   `json:"receiver"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	TransferLegId   types.TEXT    `json:"transferLegId"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveAllocate to a map for DAML arguments
func (t ApproveAllocate) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["transferLegId"] = string(t.TransferLegId)

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveAllocate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveAllocate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveAllocate to hex string (Canton MCMS format)
func (t ApproveAllocate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveAllocate from hex string (Canton MCMS format)
func (t *ApproveAllocate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveAllocateParams is a Record type
type ApproveAllocateParams struct {
	Receiver        types.PARTY   `json:"receiver"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	TransferLegId   types.TEXT    `json:"transferLegId"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveAllocateParams to a map for DAML arguments
func (t ApproveAllocateParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["transferLegId"] = string(t.TransferLegId)

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveAllocateParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveAllocateParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveAllocateParams to hex string (Canton MCMS format)
func (t ApproveAllocateParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveAllocateParams from hex string (Canton MCMS format)
func (t *ApproveAllocateParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveBurn is a Record type
type ApproveBurn struct {
	Holder          types.PARTY   `json:"holder"`
	MaxAmount       types.NUMERIC `json:"maxAmount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveBurn to a map for DAML arguments
func (t ApproveBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["holder"] = t.Holder.ToMap()

	m["maxAmount"] = t.MaxAmount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveBurn to hex string (Canton MCMS format)
func (t ApproveBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveBurn from hex string (Canton MCMS format)
func (t *ApproveBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveBurnParams is a Record type
type ApproveBurnParams struct {
	Holder          types.PARTY   `json:"holder"`
	MaxAmount       types.NUMERIC `json:"maxAmount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveBurnParams to a map for DAML arguments
func (t ApproveBurnParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["holder"] = t.Holder.ToMap()

	m["maxAmount"] = t.MaxAmount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveBurnParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveBurnParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveBurnParams to hex string (Canton MCMS format)
func (t ApproveBurnParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveBurnParams from hex string (Canton MCMS format)
func (t *ApproveBurnParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveMint is a Record type
type ApproveMint struct {
	Recipient       types.PARTY   `json:"recipient"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveMint to a map for DAML arguments
func (t ApproveMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["recipient"] = t.Recipient.ToMap()

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveMint to hex string (Canton MCMS format)
func (t ApproveMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveMint from hex string (Canton MCMS format)
func (t *ApproveMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveMintParams is a Record type
type ApproveMintParams struct {
	Recipient       types.PARTY   `json:"recipient"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveMintParams to a map for DAML arguments
func (t ApproveMintParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["recipient"] = t.Recipient.ToMap()

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveMintParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveMintParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveMintParams to hex string (Canton MCMS format)
func (t ApproveMintParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveMintParams from hex string (Canton MCMS format)
func (t *ApproveMintParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveTransfer is a Record type
type ApproveTransfer struct {
	Receiver        types.PARTY   `json:"receiver"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveTransfer to a map for DAML arguments
func (t ApproveTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveTransfer to hex string (Canton MCMS format)
func (t ApproveTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveTransfer from hex string (Canton MCMS format)
func (t *ApproveTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ApproveTransferParams is a Record type
type ApproveTransferParams struct {
	Receiver        types.PARTY   `json:"receiver"`
	Amount          types.NUMERIC `json:"amount" hex:"decimal"`
	Reference       types.TEXT    `json:"reference"`
	ValidForSeconds types.INT64   `json:"validForSeconds"`
}

// ToMap converts ApproveTransferParams to a map for DAML arguments
func (t ApproveTransferParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["reference"] = string(t.Reference)

	m["validForSeconds"] = int64(t.ValidForSeconds)

	return m
}

func (t ApproveTransferParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ApproveTransferParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ApproveTransferParams to hex string (Canton MCMS format)
func (t ApproveTransferParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ApproveTransferParams from hex string (Canton MCMS format)
func (t *ApproveTransferParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelAllocate is a Record type
type CancelAllocate struct {
}

// ToMap converts CancelAllocate to a map for DAML arguments
func (t CancelAllocate) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CancelAllocate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelAllocate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelAllocate to hex string (Canton MCMS format)
func (t CancelAllocate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelAllocate from hex string (Canton MCMS format)
func (t *CancelAllocate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelBurn is a Record type
type CancelBurn struct {
}

// ToMap converts CancelBurn to a map for DAML arguments
func (t CancelBurn) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CancelBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelBurn to hex string (Canton MCMS format)
func (t CancelBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelBurn from hex string (Canton MCMS format)
func (t *CancelBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelMint is a Record type
type CancelMint struct {
}

// ToMap converts CancelMint to a map for DAML arguments
func (t CancelMint) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CancelMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelMint to hex string (Canton MCMS format)
func (t CancelMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelMint from hex string (Canton MCMS format)
func (t *CancelMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelTransfer is a Record type
type CancelTransfer struct {
}

// ToMap converts CancelTransfer to a map for DAML arguments
func (t CancelTransfer) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CancelTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelTransfer to hex string (Canton MCMS format)
func (t CancelTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelTransfer from hex string (Canton MCMS format)
func (t *CancelTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteAllocate is a Record type
type ExecuteAllocate struct {
	TokenAdminCid        types.CONTRACT_ID                          `json:"tokenAdminCid"`
	AllocationFactoryCid types.CONTRACT_ID                          `json:"allocationFactoryCid"`
	SettlementExecutor   types.PARTY                                `json:"settlementExecutor"`
	SettlementRefId      types.TEXT                                 `json:"settlementRefId"`
	InputHoldingCids     []types.CONTRACT_ID                        `json:"inputHoldingCids"`
	AllocateBeforeHours  types.INT64                                `json:"allocateBeforeHours"`
	SettleBeforeHours    types.INT64                                `json:"settleBeforeHours"`
	RegistryContext      splice_api_token_metadata_v1.ChoiceContext `json:"registryContext"`
	Meta                 splice_api_token_metadata_v1.Metadata      `json:"meta"`
}

// ToMap converts ExecuteAllocate to a map for DAML arguments
func (t ExecuteAllocate) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminCid"] = model.NestedToDAMLValue(t.TokenAdminCid)

	m["allocationFactoryCid"] = model.NestedToDAMLValue(t.AllocationFactoryCid)

	m["settlementExecutor"] = t.SettlementExecutor.ToMap()

	m["settlementRefId"] = string(t.SettlementRefId)

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["allocateBeforeHours"] = int64(t.AllocateBeforeHours)

	m["settleBeforeHours"] = int64(t.SettleBeforeHours)

	m["registryContext"] = model.NestedToDAMLValue(t.RegistryContext)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ExecuteAllocate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteAllocate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteAllocate to hex string (Canton MCMS format)
func (t ExecuteAllocate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteAllocate from hex string (Canton MCMS format)
func (t *ExecuteAllocate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteBurn is a Record type
type ExecuteBurn struct {
	TokenAdminCid      types.CONTRACT_ID                          `json:"tokenAdminCid"`
	BurnMintFactoryCid types.CONTRACT_ID                          `json:"burnMintFactoryCid"`
	InputHoldingCids   []types.CONTRACT_ID                        `json:"inputHoldingCids"`
	RegistryContext    splice_api_token_metadata_v1.ChoiceContext `json:"registryContext"`
	Meta               splice_api_token_metadata_v1.Metadata      `json:"meta"`
}

// ToMap converts ExecuteBurn to a map for DAML arguments
func (t ExecuteBurn) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminCid"] = model.NestedToDAMLValue(t.TokenAdminCid)

	m["burnMintFactoryCid"] = model.NestedToDAMLValue(t.BurnMintFactoryCid)

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["registryContext"] = model.NestedToDAMLValue(t.RegistryContext)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ExecuteBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteBurn to hex string (Canton MCMS format)
func (t ExecuteBurn) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteBurn from hex string (Canton MCMS format)
func (t *ExecuteBurn) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteMint is a Record type
type ExecuteMint struct {
	TokenAdminCid      types.CONTRACT_ID                          `json:"tokenAdminCid"`
	BurnMintFactoryCid types.CONTRACT_ID                          `json:"burnMintFactoryCid"`
	RegistryContext    splice_api_token_metadata_v1.ChoiceContext `json:"registryContext"`
	Meta               splice_api_token_metadata_v1.Metadata      `json:"meta"`
}

// ToMap converts ExecuteMint to a map for DAML arguments
func (t ExecuteMint) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminCid"] = model.NestedToDAMLValue(t.TokenAdminCid)

	m["burnMintFactoryCid"] = model.NestedToDAMLValue(t.BurnMintFactoryCid)

	m["registryContext"] = model.NestedToDAMLValue(t.RegistryContext)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ExecuteMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteMint to hex string (Canton MCMS format)
func (t ExecuteMint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteMint from hex string (Canton MCMS format)
func (t *ExecuteMint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteTransfer is a Record type
type ExecuteTransfer struct {
	TokenAdminCid        types.CONTRACT_ID                          `json:"tokenAdminCid"`
	TransferFactoryCid   types.CONTRACT_ID                          `json:"transferFactoryCid"`
	InputHoldingCids     []types.CONTRACT_ID                        `json:"inputHoldingCids"`
	TransferTimeoutHours *types.INT64                               `json:"transferTimeoutHours" hex:"optional"`
	RegistryContext      splice_api_token_metadata_v1.ChoiceContext `json:"registryContext"`
	Meta                 splice_api_token_metadata_v1.Metadata      `json:"meta"`
}

// ToMap converts ExecuteTransfer to a map for DAML arguments
func (t ExecuteTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenAdminCid"] = model.NestedToDAMLValue(t.TokenAdminCid)

	m["transferFactoryCid"] = model.NestedToDAMLValue(t.TransferFactoryCid)

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	if t.TransferTimeoutHours != nil {
		m["transferTimeoutHours"] = map[string]any{
			"_type": "optional",
			"value": int64(*t.TransferTimeoutHours),
		}
	} else {
		m["transferTimeoutHours"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["registryContext"] = model.NestedToDAMLValue(t.RegistryContext)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t ExecuteTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteTransfer to hex string (Canton MCMS format)
func (t ExecuteTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteTransfer from hex string (Canton MCMS format)
func (t *ExecuteTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// LinkAllocationAuthorization is a Template type
type LinkAllocationAuthorization struct {
	CcipOwner       types.PARTY                              `json:"ccipOwner"`
	AdminInstanceId types.TEXT                               `json:"adminInstanceId"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Receiver        types.PARTY                              `json:"receiver"`
	Amount          types.NUMERIC                            `json:"amount" hex:"decimal"`
	TransferLegId   types.TEXT                               `json:"transferLegId"`
	Reference       types.TEXT                               `json:"reference"`
	ExpiresAt       types.TIMESTAMP                          `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkAllocationAuthorization) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkAllocationAuthorization")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkAllocationAuthorization) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.TokenAdmin", "LinkAllocationAuthorization")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkAllocationAuthorization) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferLegId"] = string(t.TransferLegId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkAllocationAuthorization) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["transferLegId"] = string(t.TransferLegId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkAllocationAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkAllocationAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkAllocationAuthorization to hex string (Canton MCMS format)
func (t LinkAllocationAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkAllocationAuthorization from hex string (Canton MCMS format)
func (t *LinkAllocationAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkAllocationAuthorization

// ExecuteAllocate exercises the ExecuteAllocate choice on this LinkAllocationAuthorization contract
// This method uses the package name in the template ID
func (t LinkAllocationAuthorization) ExecuteAllocate(contractID string, args ExecuteAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteAllocate",
		Arguments:  argsToMap(args),
	}
}

// ExecuteAllocateWithPackageID exercises the ExecuteAllocate choice using the provided package ID instead of package name
func (t LinkAllocationAuthorization) ExecuteAllocateWithPackageID(contractID string, packageID string, args ExecuteAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteAllocate",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkAllocationAuthorization contract
// This method uses the package name in the template ID
func (t LinkAllocationAuthorization) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkAllocationAuthorization) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CancelAllocate exercises the CancelAllocate choice on this LinkAllocationAuthorization contract
// This method uses the package name in the template ID
func (t LinkAllocationAuthorization) CancelAllocate(contractID string, args CancelAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "CancelAllocate",
		Arguments:  argsToMap(args),
	}
}

// CancelAllocateWithPackageID exercises the CancelAllocate choice using the provided package ID instead of package name
func (t LinkAllocationAuthorization) CancelAllocateWithPackageID(contractID string, packageID string, args CancelAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkAllocationAuthorization"),
		ContractID: contractID,
		Choice:     "CancelAllocate",
		Arguments:  argsToMap(args),
	}
}

// LinkBurnAuthorization is a Template type
type LinkBurnAuthorization struct {
	CcipOwner       types.PARTY                              `json:"ccipOwner"`
	AdminInstanceId types.TEXT                               `json:"adminInstanceId"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Holder          types.PARTY                              `json:"holder"`
	MaxAmount       types.NUMERIC                            `json:"maxAmount" hex:"decimal"`
	Reference       types.TEXT                               `json:"reference"`
	ExpiresAt       types.TIMESTAMP                          `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkBurnAuthorization) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkBurnAuthorization")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkBurnAuthorization) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.TokenAdmin", "LinkBurnAuthorization")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkBurnAuthorization) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	if t.MaxAmount != "" {
		args["maxAmount"] = t.MaxAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkBurnAuthorization) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["holder"] = t.Holder.ToMap()

	if t.MaxAmount != "" {
		args["maxAmount"] = t.MaxAmount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkBurnAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkBurnAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkBurnAuthorization to hex string (Canton MCMS format)
func (t LinkBurnAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkBurnAuthorization from hex string (Canton MCMS format)
func (t *LinkBurnAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkBurnAuthorization

// ExecuteBurn exercises the ExecuteBurn choice on this LinkBurnAuthorization contract
// This method uses the package name in the template ID
func (t LinkBurnAuthorization) ExecuteBurn(contractID string, args ExecuteBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteBurn",
		Arguments:  argsToMap(args),
	}
}

// ExecuteBurnWithPackageID exercises the ExecuteBurn choice using the provided package ID instead of package name
func (t LinkBurnAuthorization) ExecuteBurnWithPackageID(contractID string, packageID string, args ExecuteBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteBurn",
		Arguments:  argsToMap(args),
	}
}

// CancelBurn exercises the CancelBurn choice on this LinkBurnAuthorization contract
// This method uses the package name in the template ID
func (t LinkBurnAuthorization) CancelBurn(contractID string, args CancelBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "CancelBurn",
		Arguments:  argsToMap(args),
	}
}

// CancelBurnWithPackageID exercises the CancelBurn choice using the provided package ID instead of package name
func (t LinkBurnAuthorization) CancelBurnWithPackageID(contractID string, packageID string, args CancelBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "CancelBurn",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkBurnAuthorization contract
// This method uses the package name in the template ID
func (t LinkBurnAuthorization) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkBurnAuthorization) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkBurnAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// LinkMintAuthorization is a Template type
type LinkMintAuthorization struct {
	CcipOwner       types.PARTY                              `json:"ccipOwner"`
	AdminInstanceId types.TEXT                               `json:"adminInstanceId"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Recipient       types.PARTY                              `json:"recipient"`
	Amount          types.NUMERIC                            `json:"amount" hex:"decimal"`
	Reference       types.TEXT                               `json:"reference"`
	ExpiresAt       types.TIMESTAMP                          `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkMintAuthorization) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkMintAuthorization")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkMintAuthorization) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.TokenAdmin", "LinkMintAuthorization")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkMintAuthorization) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["recipient"] = t.Recipient.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkMintAuthorization) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["recipient"] = t.Recipient.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkMintAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkMintAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkMintAuthorization to hex string (Canton MCMS format)
func (t LinkMintAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkMintAuthorization from hex string (Canton MCMS format)
func (t *LinkMintAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkMintAuthorization

// ExecuteMint exercises the ExecuteMint choice on this LinkMintAuthorization contract
// This method uses the package name in the template ID
func (t LinkMintAuthorization) ExecuteMint(contractID string, args ExecuteMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteMint",
		Arguments:  argsToMap(args),
	}
}

// ExecuteMintWithPackageID exercises the ExecuteMint choice using the provided package ID instead of package name
func (t LinkMintAuthorization) ExecuteMintWithPackageID(contractID string, packageID string, args ExecuteMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteMint",
		Arguments:  argsToMap(args),
	}
}

// CancelMint exercises the CancelMint choice on this LinkMintAuthorization contract
// This method uses the package name in the template ID
func (t LinkMintAuthorization) CancelMint(contractID string, args CancelMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "CancelMint",
		Arguments:  argsToMap(args),
	}
}

// CancelMintWithPackageID exercises the CancelMint choice using the provided package ID instead of package name
func (t LinkMintAuthorization) CancelMintWithPackageID(contractID string, packageID string, args CancelMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "CancelMint",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkMintAuthorization contract
// This method uses the package name in the template ID
func (t LinkMintAuthorization) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkMintAuthorization) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkMintAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// LinkTokenAdmin is a Template type
type LinkTokenAdmin struct {
	CcipOwner    types.PARTY                              `json:"ccipOwner"`
	InstanceId   types.TEXT                               `json:"instanceId"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Paused       types.BOOL                               `json:"paused"`
	Observers    []types.PARTY                            `json:"observers"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkTokenAdmin) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkTokenAdmin) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkTokenAdmin) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["paused"] = bool(t.Paused)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkTokenAdmin) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["paused"] = bool(t.Paused)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkTokenAdmin) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkTokenAdmin) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkTokenAdmin to hex string (Canton MCMS format)
func (t LinkTokenAdmin) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkTokenAdmin from hex string (Canton MCMS format)
func (t *LinkTokenAdmin) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkTokenAdmin

// ApproveMint exercises the ApproveMint choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) ApproveMint(contractID string, args ApproveMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveMint",
		Arguments:  argsToMap(args),
	}
}

// ApproveMintWithPackageID exercises the ApproveMint choice using the provided package ID instead of package name
func (t LinkTokenAdmin) ApproveMintWithPackageID(contractID string, packageID string, args ApproveMint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveMint",
		Arguments:  argsToMap(args),
	}
}

// ApproveBurn exercises the ApproveBurn choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) ApproveBurn(contractID string, args ApproveBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveBurn",
		Arguments:  argsToMap(args),
	}
}

// ApproveBurnWithPackageID exercises the ApproveBurn choice using the provided package ID instead of package name
func (t LinkTokenAdmin) ApproveBurnWithPackageID(contractID string, packageID string, args ApproveBurn) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveBurn",
		Arguments:  argsToMap(args),
	}
}

// ApproveTransfer exercises the ApproveTransfer choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) ApproveTransfer(contractID string, args ApproveTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveTransfer",
		Arguments:  argsToMap(args),
	}
}

// ApproveTransferWithPackageID exercises the ApproveTransfer choice using the provided package ID instead of package name
func (t LinkTokenAdmin) ApproveTransferWithPackageID(contractID string, packageID string, args ApproveTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveTransfer",
		Arguments:  argsToMap(args),
	}
}

// ApproveAllocate exercises the ApproveAllocate choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) ApproveAllocate(contractID string, args ApproveAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveAllocate",
		Arguments:  argsToMap(args),
	}
}

// ApproveAllocateWithPackageID exercises the ApproveAllocate choice using the provided package ID instead of package name
func (t LinkTokenAdmin) ApproveAllocateWithPackageID(contractID string, packageID string, args ApproveAllocate) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "ApproveAllocate",
		Arguments:  argsToMap(args),
	}
}

// SetPaused exercises the SetPaused choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) SetPaused(contractID string, args SetPaused) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "SetPaused",
		Arguments:  argsToMap(args),
	}
}

// SetPausedWithPackageID exercises the SetPaused choice using the provided package ID instead of package name
func (t LinkTokenAdmin) SetPausedWithPackageID(contractID string, packageID string, args SetPaused) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "SetPaused",
		Arguments:  argsToMap(args),
	}
}

// SetObservers exercises the SetObservers choice on this LinkTokenAdmin contract
// This method uses the package name in the template ID
func (t LinkTokenAdmin) SetObservers(contractID string, args SetObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "SetObservers",
		Arguments:  argsToMap(args),
	}
}

// SetObserversWithPackageID exercises the SetObservers choice using the provided package ID instead of package name
func (t LinkTokenAdmin) SetObserversWithPackageID(contractID string, packageID string, args SetObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTokenAdmin"),
		ContractID: contractID,
		Choice:     "SetObservers",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkTokenAdmin contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t LinkTokenAdmin) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkTokenAdmin) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this LinkTokenAdmin contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t LinkTokenAdmin) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t LinkTokenAdmin) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for LinkTokenAdmin

var _ api.IMCMSReceiver = (*LinkTokenAdmin)(nil)

// LinkTransferAuthorization is a Template type
type LinkTransferAuthorization struct {
	CcipOwner       types.PARTY                              `json:"ccipOwner"`
	AdminInstanceId types.TEXT                               `json:"adminInstanceId"`
	InstrumentId    splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Receiver        types.PARTY                              `json:"receiver"`
	Amount          types.NUMERIC                            `json:"amount" hex:"decimal"`
	Reference       types.TEXT                               `json:"reference"`
	ExpiresAt       types.TIMESTAMP                          `json:"expiresAt"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t LinkTransferAuthorization) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTransferAuthorization")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t LinkTransferAuthorization) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTransferAuthorization")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t LinkTransferAuthorization) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t LinkTransferAuthorization) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["adminInstanceId"] = string(t.AdminInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiver"] = t.Receiver.ToMap()

	if t.Amount != "" {
		args["amount"] = t.Amount
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["reference"] = string(t.Reference)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["expiresAt"] = t.ExpiresAt

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t LinkTransferAuthorization) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *LinkTransferAuthorization) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes LinkTransferAuthorization to hex string (Canton MCMS format)
func (t LinkTransferAuthorization) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes LinkTransferAuthorization from hex string (Canton MCMS format)
func (t *LinkTransferAuthorization) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for LinkTransferAuthorization

// ExecuteTransfer exercises the ExecuteTransfer choice on this LinkTransferAuthorization contract
// This method uses the package name in the template ID
func (t LinkTransferAuthorization) ExecuteTransfer(contractID string, args ExecuteTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteTransfer",
		Arguments:  argsToMap(args),
	}
}

// ExecuteTransferWithPackageID exercises the ExecuteTransfer choice using the provided package ID instead of package name
func (t LinkTransferAuthorization) ExecuteTransferWithPackageID(contractID string, packageID string, args ExecuteTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "ExecuteTransfer",
		Arguments:  argsToMap(args),
	}
}

// CancelTransfer exercises the CancelTransfer choice on this LinkTransferAuthorization contract
// This method uses the package name in the template ID
func (t LinkTransferAuthorization) CancelTransfer(contractID string, args CancelTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "CancelTransfer",
		Arguments:  argsToMap(args),
	}
}

// CancelTransferWithPackageID exercises the CancelTransfer choice using the provided package ID instead of package name
func (t LinkTransferAuthorization) CancelTransferWithPackageID(contractID string, packageID string, args CancelTransfer) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "CancelTransfer",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this LinkTransferAuthorization contract
// This method uses the package name in the template ID
func (t LinkTransferAuthorization) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t LinkTransferAuthorization) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "Link.TokenAdmin", "LinkTransferAuthorization"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetObservers is a Record type
type SetObservers struct {
	Observers []types.PARTY `json:"observers"`
}

// ToMap converts SetObservers to a map for DAML arguments
func (t SetObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t SetObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetObservers to hex string (Canton MCMS format)
func (t SetObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetObservers from hex string (Canton MCMS format)
func (t *SetObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetObserversParams is a Record type
type SetObserversParams struct {
	Observers []types.PARTY `json:"observers"`
}

// ToMap converts SetObserversParams to a map for DAML arguments
func (t SetObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["observers"] = func() []any {
		res := make([]any, 0, len(t.Observers))
		for _, e := range t.Observers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t SetObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetObserversParams to hex string (Canton MCMS format)
func (t SetObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetObserversParams from hex string (Canton MCMS format)
func (t *SetObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPaused is a Record type
type SetPaused struct {
	Paused types.BOOL `json:"paused"`
}

// ToMap converts SetPaused to a map for DAML arguments
func (t SetPaused) ToMap() map[string]any {
	m := make(map[string]any)

	m["paused"] = bool(t.Paused)

	return m
}

func (t SetPaused) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPaused) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPaused to hex string (Canton MCMS format)
func (t SetPaused) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPaused from hex string (Canton MCMS format)
func (t *SetPaused) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetPausedParams is a Record type
type SetPausedParams struct {
	Paused types.BOOL `json:"paused"`
}

// ToMap converts SetPausedParams to a map for DAML arguments
func (t SetPausedParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["paused"] = bool(t.Paused)

	return m
}

func (t SetPausedParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetPausedParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetPausedParams to hex string (Canton MCMS format)
func (t SetPausedParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetPausedParams from hex string (Canton MCMS format)
func (t *SetPausedParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ApproveAllocate(args ApproveAllocate) (*bind.EncodedChoice, error)
	ApproveAllocateParams(args ApproveAllocateParams) (*bind.EncodedChoice, error)
	ApproveBurn(args ApproveBurn) (*bind.EncodedChoice, error)
	ApproveBurnParams(args ApproveBurnParams) (*bind.EncodedChoice, error)
	ApproveMint(args ApproveMint) (*bind.EncodedChoice, error)
	ApproveMintParams(args ApproveMintParams) (*bind.EncodedChoice, error)
	ApproveTransfer(args ApproveTransfer) (*bind.EncodedChoice, error)
	ApproveTransferParams(args ApproveTransferParams) (*bind.EncodedChoice, error)
	CancelAllocate(args CancelAllocate) (*bind.EncodedChoice, error)
	CancelBurn(args CancelBurn) (*bind.EncodedChoice, error)
	CancelMint(args CancelMint) (*bind.EncodedChoice, error)
	CancelTransfer(args CancelTransfer) (*bind.EncodedChoice, error)
	ExecuteAllocate(args ExecuteAllocate) (*bind.EncodedChoice, error)
	ExecuteBurn(args ExecuteBurn) (*bind.EncodedChoice, error)
	ExecuteMint(args ExecuteMint) (*bind.EncodedChoice, error)
	ExecuteTransfer(args ExecuteTransfer) (*bind.EncodedChoice, error)
	SetObservers(args SetObservers) (*bind.EncodedChoice, error)
	SetObserversParams(args SetObserversParams) (*bind.EncodedChoice, error)
	SetPaused(args SetPaused) (*bind.EncodedChoice, error)
	SetPausedParams(args SetPausedParams) (*bind.EncodedChoice, error)
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

// ApproveAllocate encodes parameters for the ApproveAllocate choice.
func (e *encoder) ApproveAllocate(args ApproveAllocate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveAllocate", args)
}

// ApproveAllocateParams encodes parameters for the ApproveAllocate choice.
func (e *encoder) ApproveAllocateParams(args ApproveAllocateParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveAllocate", args)
}

// ApproveBurn encodes parameters for the ApproveBurn choice.
func (e *encoder) ApproveBurn(args ApproveBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveBurn", args)
}

// ApproveBurnParams encodes parameters for the ApproveBurn choice.
func (e *encoder) ApproveBurnParams(args ApproveBurnParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveBurn", args)
}

// ApproveMint encodes parameters for the ApproveMint choice.
func (e *encoder) ApproveMint(args ApproveMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveMint", args)
}

// ApproveMintParams encodes parameters for the ApproveMint choice.
func (e *encoder) ApproveMintParams(args ApproveMintParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveMint", args)
}

// ApproveTransfer encodes parameters for the ApproveTransfer choice.
func (e *encoder) ApproveTransfer(args ApproveTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveTransfer", args)
}

// ApproveTransferParams encodes parameters for the ApproveTransfer choice.
func (e *encoder) ApproveTransferParams(args ApproveTransferParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ApproveTransfer", args)
}

// CancelAllocate encodes parameters for the CancelAllocate choice.
func (e *encoder) CancelAllocate(args CancelAllocate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelAllocate", args)
}

// CancelBurn encodes parameters for the CancelBurn choice.
func (e *encoder) CancelBurn(args CancelBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelBurn", args)
}

// CancelMint encodes parameters for the CancelMint choice.
func (e *encoder) CancelMint(args CancelMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelMint", args)
}

// CancelTransfer encodes parameters for the CancelTransfer choice.
func (e *encoder) CancelTransfer(args CancelTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelTransfer", args)
}

// ExecuteAllocate encodes parameters for the ExecuteAllocate choice.
func (e *encoder) ExecuteAllocate(args ExecuteAllocate) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteAllocate", args)
}

// ExecuteBurn encodes parameters for the ExecuteBurn choice.
func (e *encoder) ExecuteBurn(args ExecuteBurn) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteBurn", args)
}

// ExecuteMint encodes parameters for the ExecuteMint choice.
func (e *encoder) ExecuteMint(args ExecuteMint) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteMint", args)
}

// ExecuteTransfer encodes parameters for the ExecuteTransfer choice.
func (e *encoder) ExecuteTransfer(args ExecuteTransfer) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteTransfer", args)
}

// SetObservers encodes parameters for the SetObservers choice.
func (e *encoder) SetObservers(args SetObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetObservers", args)
}

// SetObserversParams encodes parameters for the SetObservers choice.
func (e *encoder) SetObserversParams(args SetObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetObservers", args)
}

// SetPaused encodes parameters for the SetPaused choice.
func (e *encoder) SetPaused(args SetPaused) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPaused", args)
}

// SetPausedParams encodes parameters for the SetPaused choice.
func (e *encoder) SetPausedParams(args SetPausedParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetPaused", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
