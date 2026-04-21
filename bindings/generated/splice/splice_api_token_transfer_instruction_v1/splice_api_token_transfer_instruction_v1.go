package splice_api_token_transfer_instruction_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "splice-api-token-transfer-instruction-v1"
	PackageID   = "55ba4deb0ad4662c4168b39859738a0e91388d252286480c7331b3f71a517281"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// ITransferFactory is a DAML interface
type ITransferFactory interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TransferFactoryTransfer executes the TransferFactory_Transfer choice
	TransferFactoryTransfer(contractID string, args TransferFactoryTransfer) *model.ExerciseCommand

	// TransferFactoryPublicFetch executes the TransferFactory_PublicFetch choice
	TransferFactoryPublicFetch(contractID string, args TransferFactoryPublicFetch) *model.ExerciseCommand
}

// ITransferInstruction is a DAML interface
type ITransferInstruction interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TransferInstructionAccept executes the TransferInstruction_Accept choice
	TransferInstructionAccept(contractID string, args TransferInstructionAccept) *model.ExerciseCommand

	// TransferInstructionReject executes the TransferInstruction_Reject choice
	TransferInstructionReject(contractID string, args TransferInstructionReject) *model.ExerciseCommand

	// TransferInstructionWithdraw executes the TransferInstruction_Withdraw choice
	TransferInstructionWithdraw(contractID string, args TransferInstructionWithdraw) *model.ExerciseCommand

	// TransferInstructionUpdate executes the TransferInstruction_Update choice
	TransferInstructionUpdate(contractID string, args TransferInstructionUpdate) *model.ExerciseCommand
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

// Transfer is a Record type
type Transfer struct {
	Sender           types.PARTY                              `json:"sender"`
	Receiver         types.PARTY                              `json:"receiver"`
	Amount           types.NUMERIC                            `json:"amount"`
	InstrumentId     splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	RequestedAt      types.TIMESTAMP                          `json:"requestedAt"`
	ExecuteBefore    types.TIMESTAMP                          `json:"executeBefore"`
	InputHoldingCids []types.CONTRACT_ID                      `json:"inputHoldingCids"`
	Meta             splice_api_token_metadata_v1.Metadata    `json:"meta"`
}

// ToMap converts Transfer to a map for DAML arguments
func (t Transfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["sender"] = t.Sender.ToMap()

	m["receiver"] = t.Receiver.ToMap()

	m["amount"] = t.Amount

	m["instrumentId"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.ToMap()
		}
		return t.InstrumentId
	}()

	m["requestedAt"] = t.RequestedAt

	m["executeBefore"] = t.ExecuteBefore

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.ToMap()
		}
		return t.Meta
	}()

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

// TransferFactoryView is a Record type
type TransferFactoryView struct {
	Admin types.PARTY                           `json:"admin"`
	Meta  splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts TransferFactoryView to a map for DAML arguments
func (t TransferFactoryView) ToMap() map[string]any {
	m := make(map[string]any)

	m["admin"] = t.Admin.ToMap()

	m["meta"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.ToMap()
		}
		return t.Meta
	}()

	return m
}

func (t TransferFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferFactoryView to hex string (Canton MCMS format)
func (t TransferFactoryView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferFactoryView from hex string (Canton MCMS format)
func (t *TransferFactoryView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferFactoryPublicFetch is a Record type
type TransferFactoryPublicFetch struct {
	ExpectedAdmin types.PARTY `json:"expectedAdmin"`
	Actor         types.PARTY `json:"actor"`
}

// ToMap converts TransferFactoryPublicFetch to a map for DAML arguments
func (t TransferFactoryPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t TransferFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferFactoryPublicFetch to hex string (Canton MCMS format)
func (t TransferFactoryPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferFactoryPublicFetch from hex string (Canton MCMS format)
func (t *TransferFactoryPublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferFactoryTransfer is a Record type
type TransferFactoryTransfer struct {
	ExpectedAdmin types.PARTY                            `json:"expectedAdmin"`
	Transfer      Transfer                               `json:"transfer"`
	ExtraArgs     splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferFactoryTransfer to a map for DAML arguments
func (t TransferFactoryTransfer) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["transfer"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Transfer).(mapper); ok {
			return m.ToMap()
		}
		return t.Transfer
	}()

	m["extraArgs"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t TransferFactoryTransfer) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferFactoryTransfer) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferFactoryTransfer to hex string (Canton MCMS format)
func (t TransferFactoryTransfer) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferFactoryTransfer from hex string (Canton MCMS format)
func (t *TransferFactoryTransfer) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionResult is a Record type
type TransferInstructionResult struct {
	Output           TransferInstructionResultOutput       `json:"output"`
	SenderChangeCids []types.CONTRACT_ID                   `json:"senderChangeCids"`
	Meta             splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts TransferInstructionResult to a map for DAML arguments
func (t TransferInstructionResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["output"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Output).(mapper); ok {
			return m.ToMap()
		}
		return t.Output
	}()

	m["senderChangeCids"] = func() []any {
		res := make([]any, 0, len(t.SenderChangeCids))
		for _, e := range t.SenderChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.ToMap()
		}
		return t.Meta
	}()

	return m
}

func (t TransferInstructionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionResult to hex string (Canton MCMS format)
func (t TransferInstructionResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionResult from hex string (Canton MCMS format)
func (t *TransferInstructionResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionResultCompleted is a Record type
type TransferInstructionResultCompleted struct {
	ReceiverHoldingCids []types.CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts TransferInstructionResultCompleted to a map for DAML arguments
func (t TransferInstructionResultCompleted) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.ReceiverHoldingCids))
		for _, e := range t.ReceiverHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t TransferInstructionResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionResultCompleted to hex string (Canton MCMS format)
func (t TransferInstructionResultCompleted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionResultCompleted from hex string (Canton MCMS format)
func (t *TransferInstructionResultCompleted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionResultOutput is a variant/union type
type TransferInstructionResultOutput struct {
	TransferInstructionResultPending   *TransferInstructionResultPending   `json:"TransferInstructionResult_Pending,omitempty"`
	TransferInstructionResultCompleted *TransferInstructionResultCompleted `json:"TransferInstructionResult_Completed,omitempty"`
	TransferInstructionResultFailed    *types.UNIT                         `json:"TransferInstructionResult_Failed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionResultOutput
func (v TransferInstructionResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for TransferInstructionResultOutput
func (v *TransferInstructionResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes TransferInstructionResultOutput to hex string (Canton MCMS format)
func (v TransferInstructionResultOutput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes TransferInstructionResultOutput from hex string (Canton MCMS format)
func (v *TransferInstructionResultOutput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferInstructionResultOutput) GetVariantTag() string {

	if v.TransferInstructionResultPending != nil {
		return "TransferInstructionResult_Pending"
	}

	if v.TransferInstructionResultCompleted != nil {
		return "TransferInstructionResult_Completed"
	}

	if v.TransferInstructionResultFailed != nil {
		return "TransferInstructionResult_Failed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferInstructionResultOutput) GetVariantValue() any {

	if v.TransferInstructionResultPending != nil {
		return v.TransferInstructionResultPending
	}

	if v.TransferInstructionResultCompleted != nil {
		return v.TransferInstructionResultCompleted
	}

	if v.TransferInstructionResultFailed != nil {
		return v.TransferInstructionResultFailed
	}

	return nil
}

var _ types.VARIANT = (*TransferInstructionResultOutput)(nil)

// TransferInstructionResultPending is a Record type
type TransferInstructionResultPending struct {
	TransferInstructionCid types.CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts TransferInstructionResultPending to a map for DAML arguments
func (t TransferInstructionResultPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["transferInstructionCid"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.TransferInstructionCid).(mapper); ok {
			return m.ToMap()
		}
		return t.TransferInstructionCid
	}()

	return m
}

func (t TransferInstructionResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionResultPending to hex string (Canton MCMS format)
func (t TransferInstructionResultPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionResultPending from hex string (Canton MCMS format)
func (t *TransferInstructionResultPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionStatus is a variant/union type
type TransferInstructionStatus struct {
	TransferPendingReceiverAcceptance *types.UNIT                      `json:"TransferPendingReceiverAcceptance,omitempty"`
	TransferPendingInternalWorkflow   *TransferPendingInternalWorkflow `json:"TransferPendingInternalWorkflow,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for TransferInstructionStatus
func (v TransferInstructionStatus) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for TransferInstructionStatus
func (v *TransferInstructionStatus) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes TransferInstructionStatus to hex string (Canton MCMS format)
func (v TransferInstructionStatus) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes TransferInstructionStatus from hex string (Canton MCMS format)
func (v *TransferInstructionStatus) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v TransferInstructionStatus) GetVariantTag() string {

	if v.TransferPendingReceiverAcceptance != nil {
		return "TransferPendingReceiverAcceptance"
	}

	if v.TransferPendingInternalWorkflow != nil {
		return "TransferPendingInternalWorkflow"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v TransferInstructionStatus) GetVariantValue() any {

	if v.TransferPendingReceiverAcceptance != nil {
		return v.TransferPendingReceiverAcceptance
	}

	if v.TransferPendingInternalWorkflow != nil {
		return v.TransferPendingInternalWorkflow
	}

	return nil
}

var _ types.VARIANT = (*TransferInstructionStatus)(nil)

// TransferInstructionView is a Record type
type TransferInstructionView struct {
	OriginalInstructionCid *types.CONTRACT_ID                    `json:"originalInstructionCid" hex:"optional"`
	Transfer               Transfer                              `json:"transfer"`
	Status                 TransferInstructionStatus             `json:"status"`
	Meta                   splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts TransferInstructionView to a map for DAML arguments
func (t TransferInstructionView) ToMap() map[string]any {
	m := make(map[string]any)

	if t.OriginalInstructionCid != nil {
		m["originalInstructionCid"] = map[string]any{
			"_type": "optional",
			"value": *t.OriginalInstructionCid,
		}
	} else {
		m["originalInstructionCid"] = map[string]any{
			"_type": "optional",
		}
	}

	m["transfer"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Transfer).(mapper); ok {
			return m.ToMap()
		}
		return t.Transfer
	}()

	m["status"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Status).(mapper); ok {
			return m.ToMap()
		}
		return t.Status
	}()

	m["meta"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.Meta).(mapper); ok {
			return m.ToMap()
		}
		return t.Meta
	}()

	return m
}

func (t TransferInstructionView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionView to hex string (Canton MCMS format)
func (t TransferInstructionView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionView from hex string (Canton MCMS format)
func (t *TransferInstructionView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionAccept is a Record type
type TransferInstructionAccept struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionAccept to a map for DAML arguments
func (t TransferInstructionAccept) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t TransferInstructionAccept) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionAccept) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionAccept to hex string (Canton MCMS format)
func (t TransferInstructionAccept) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionAccept from hex string (Canton MCMS format)
func (t *TransferInstructionAccept) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionReject is a Record type
type TransferInstructionReject struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionReject to a map for DAML arguments
func (t TransferInstructionReject) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t TransferInstructionReject) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionReject) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionReject to hex string (Canton MCMS format)
func (t TransferInstructionReject) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionReject from hex string (Canton MCMS format)
func (t *TransferInstructionReject) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionUpdate is a Record type
type TransferInstructionUpdate struct {
	ExtraActors []types.PARTY                          `json:"extraActors"`
	ExtraArgs   splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionUpdate to a map for DAML arguments
func (t TransferInstructionUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraActors"] = func() []any {
		res := make([]any, 0, len(t.ExtraActors))
		for _, e := range t.ExtraActors {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["extraArgs"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t TransferInstructionUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionUpdate to hex string (Canton MCMS format)
func (t TransferInstructionUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionUpdate from hex string (Canton MCMS format)
func (t *TransferInstructionUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferInstructionWithdraw is a Record type
type TransferInstructionWithdraw struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts TransferInstructionWithdraw to a map for DAML arguments
func (t TransferInstructionWithdraw) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = func() any {
		type mapper interface{ ToMap() map[string]any }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.ToMap()
		}
		return t.ExtraArgs
	}()

	return m
}

func (t TransferInstructionWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferInstructionWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferInstructionWithdraw to hex string (Canton MCMS format)
func (t TransferInstructionWithdraw) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferInstructionWithdraw from hex string (Canton MCMS format)
func (t *TransferInstructionWithdraw) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TransferPendingInternalWorkflow is a Record type
type TransferPendingInternalWorkflow struct {
	PendingActions types.GENMAP `json:"pendingActions"`
}

// ToMap converts TransferPendingInternalWorkflow to a map for DAML arguments
func (t TransferPendingInternalWorkflow) ToMap() map[string]any {
	m := make(map[string]any)

	m["pendingActions"] = func() any {
		if t.PendingActions == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.PendingActions}
	}()

	return m
}

func (t TransferPendingInternalWorkflow) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TransferPendingInternalWorkflow) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TransferPendingInternalWorkflow to hex string (Canton MCMS format)
func (t TransferPendingInternalWorkflow) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TransferPendingInternalWorkflow from hex string (Canton MCMS format)
func (t *TransferPendingInternalWorkflow) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ITransferFactoryInterfaceID returns the interface ID for the ITransferFactory interface using the package name
func ITransferFactoryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.TransferInstructionV1", "TransferFactory")
}

// ITransferFactoryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ITransferFactoryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.TransferInstructionV1", "TransferFactory")
}

// ITransferInstructionInterfaceID returns the interface ID for the ITransferInstruction interface using the package name
func ITransferInstructionInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.TransferInstructionV1", "TransferInstruction")
}

// ITransferInstructionInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func ITransferInstructionInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.TransferInstructionV1", "TransferInstruction")
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
