package splice_api_token_allocation_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
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
	PackageName = "splice-api-token-allocation-v1"
	PackageID   = "93c942ae2b4c2ba674fb152fe38473c507bda4e82b4e4c5da55a552a9d8cce1d"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IAllocation is a DAML interface
type IAllocation interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationWithdraw executes the Allocation_Withdraw choice
	AllocationWithdraw(contractID string, args AllocationWithdraw) *model.ExerciseCommand

	// AllocationCancel executes the Allocation_Cancel choice
	AllocationCancel(contractID string, args AllocationCancel) *model.ExerciseCommand

	// AllocationExecuteTransfer executes the Allocation_ExecuteTransfer choice
	AllocationExecuteTransfer(contractID string, args AllocationExecuteTransfer) *model.ExerciseCommand
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

// AllocationSpecification is a Record type
type AllocationSpecification struct {
	Settlement    SettlementInfo `json:"settlement"`
	TransferLegId types.TEXT     `json:"transferLegId"`
	TransferLeg   TransferLeg    `json:"transferLeg"`
}

// ToMap converts AllocationSpecification to a map for DAML arguments
func (t AllocationSpecification) ToMap() map[string]any {
	m := make(map[string]any)

	m["settlement"] = model.NestedToDAMLValue(t.Settlement)

	m["transferLegId"] = string(t.TransferLegId)

	m["transferLeg"] = model.NestedToDAMLValue(t.TransferLeg)

	return m
}

func (t AllocationSpecification) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationSpecification) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationSpecification to hex string (Canton MCMS format)
func (t AllocationSpecification) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationSpecification from hex string (Canton MCMS format)
func (t *AllocationSpecification) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationView is a Record type
type AllocationView struct {
	Allocation  AllocationSpecification               `json:"allocation"`
	HoldingCids []types.CONTRACT_ID                   `json:"holdingCids"`
	Meta        splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationView to a map for DAML arguments
func (t AllocationView) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocation"] = model.NestedToDAMLValue(t.Allocation)

	m["holdingCids"] = func() []any {
		res := make([]any, 0, len(t.HoldingCids))
		for _, e := range t.HoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationView to hex string (Canton MCMS format)
func (t AllocationView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationView from hex string (Canton MCMS format)
func (t *AllocationView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationCancel is a Record type
type AllocationCancel struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationCancel to a map for DAML arguments
func (t AllocationCancel) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationCancel) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationCancel) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationCancel to hex string (Canton MCMS format)
func (t AllocationCancel) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationCancel from hex string (Canton MCMS format)
func (t *AllocationCancel) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationCancelResult is a Record type
type AllocationCancelResult struct {
	SenderHoldingCids []types.CONTRACT_ID                   `json:"senderHoldingCids"`
	Meta              splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationCancelResult to a map for DAML arguments
func (t AllocationCancelResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.SenderHoldingCids))
		for _, e := range t.SenderHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationCancelResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationCancelResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationCancelResult to hex string (Canton MCMS format)
func (t AllocationCancelResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationCancelResult from hex string (Canton MCMS format)
func (t *AllocationCancelResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationExecuteTransfer is a Record type
type AllocationExecuteTransfer struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationExecuteTransfer to a map for DAML arguments
func (t AllocationExecuteTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationExecuteTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationExecuteTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationExecuteTransfer to hex string (Canton MCMS format)
func (t AllocationExecuteTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationExecuteTransfer from hex string (Canton MCMS format)
func (t *AllocationExecuteTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationExecuteTransferResult is a Record type
type AllocationExecuteTransferResult struct {
	SenderHoldingCids   []types.CONTRACT_ID                   `json:"senderHoldingCids"`
	ReceiverHoldingCids []types.CONTRACT_ID                   `json:"receiverHoldingCids"`
	Meta                splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationExecuteTransferResult to a map for DAML arguments
func (t AllocationExecuteTransferResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.SenderHoldingCids))
		for _, e := range t.SenderHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["receiverHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverHoldingCids))
		for _, e := range t.ReceiverHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationExecuteTransferResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationExecuteTransferResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationExecuteTransferResult to hex string (Canton MCMS format)
func (t AllocationExecuteTransferResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationExecuteTransferResult from hex string (Canton MCMS format)
func (t *AllocationExecuteTransferResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationWithdraw is a Record type
type AllocationWithdraw struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationWithdraw to a map for DAML arguments
func (t AllocationWithdraw) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationWithdraw to hex string (Canton MCMS format)
func (t AllocationWithdraw) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationWithdraw from hex string (Canton MCMS format)
func (t *AllocationWithdraw) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationWithdrawResult is a Record type
type AllocationWithdrawResult struct {
	SenderHoldingCids []types.CONTRACT_ID                   `json:"senderHoldingCids"`
	Meta              splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationWithdrawResult to a map for DAML arguments
func (t AllocationWithdrawResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.SenderHoldingCids))
		for _, e := range t.SenderHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationWithdrawResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationWithdrawResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationWithdrawResult to hex string (Canton MCMS format)
func (t AllocationWithdrawResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationWithdrawResult from hex string (Canton MCMS format)
func (t *AllocationWithdrawResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Reference is a Record type
type Reference struct {
	Id  types.TEXT         `json:"id"`
	Cid *types.CONTRACT_ID `json:"cid" hex:"optional"`
}

// ToMap converts Reference to a map for DAML arguments
func (t Reference) ToMap() map[string]any {
	m := make(map[string]any)

	m["id"] = string(t.Id)

	if t.Cid != nil {
		m["cid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.Cid),
		}
	} else {
		m["cid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t Reference) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Reference) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Reference to hex string (Canton MCMS format)
func (t Reference) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Reference from hex string (Canton MCMS format)
func (t *Reference) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SettlementInfo is a Record type
type SettlementInfo struct {
	Executor       types.PARTY                           `json:"executor"`
	SettlementRef  Reference                             `json:"settlementRef"`
	RequestedAt    types.TIMESTAMP                       `json:"requestedAt"`
	AllocateBefore types.TIMESTAMP                       `json:"allocateBefore"`
	SettleBefore   types.TIMESTAMP                       `json:"settleBefore"`
	Meta           splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts SettlementInfo to a map for DAML arguments
func (t SettlementInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["executor"] = t.Executor.ToMap()

	m["settlementRef"] = model.NestedToDAMLValue(t.SettlementRef)

	m["requestedAt"] = t.RequestedAt

	m["allocateBefore"] = t.AllocateBefore

	m["settleBefore"] = t.SettleBefore

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t SettlementInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SettlementInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SettlementInfo to hex string (Canton MCMS format)
func (t SettlementInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SettlementInfo from hex string (Canton MCMS format)
func (t *SettlementInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferLeg is a Record type
type TransferLeg struct {
	Sender       types.PARTY                              `json:"sender"`
	Receiver     types.PARTY                              `json:"receiver"`
	Amount       types.NUMERIC                            `json:"amount"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Meta         splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// ToMap converts TransferLeg to a map for DAML arguments
func (t TransferLeg) ToMap() map[string]any {
	m := make(map[string]any)

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t TransferLeg) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferLeg) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferLeg to hex string (Canton MCMS format)
func (t TransferLeg) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferLeg from hex string (Canton MCMS format)
func (t *TransferLeg) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IAllocationInterfaceID returns the interface ID for the IAllocation interface using the package name
func IAllocationInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.AllocationV1", "Allocation")
}

// IAllocationInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IAllocationInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.AllocationV1", "Allocation")
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
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

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
