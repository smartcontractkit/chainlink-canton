package interfaces

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

const PackageID = "e3e5bc3d80f265d41e4e9bd293f695460d94ed8d2364bd80a625ecf3e1f725fa"
const SDKVersion = "3.4.8"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IITokenPool is a DAML interface
type IITokenPool interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenPoolGetRequiredCCVs executes the TokenPool_GetRequiredCCVs choice
	TokenPoolGetRequiredCCVs(contractID string, args TokenPoolGetRequiredCCVs) *model.ExerciseCommand

	// TokenPoolVerifyCCVs executes the TokenPool_VerifyCCVs choice
	TokenPoolVerifyCCVs(contractID string, args TokenPoolVerifyCCVs) *model.ExerciseCommand

	// TokenPoolReleaseFromTicket executes the TokenPool_ReleaseFromTicket choice
	TokenPoolReleaseFromTicket(contractID string, args TokenPoolReleaseFromTicket) *model.ExerciseCommand

	// TokenPoolLockOrBurn executes the TokenPool_LockOrBurn choice
	TokenPoolLockOrBurn(contractID string, args TokenPoolLockOrBurn) *model.ExerciseCommand
}

func argsToMap(args interface{}) map[string]interface{} {
	if args == nil {
		return map[string]interface{}{}
	}

	if m, ok := args.(map[string]interface{}); ok {
		return m
	}

	// Check if the type has a ToMap method
	type mapper interface {
		ToMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// FeeInput is a Record type
type FeeInput struct {
	InstrumentId     InstrumentId  `json:"instrumentId"`
	TransferFactory  CONTRACT_ID   `json:"transferFactory"`
	ExtraArgs        ExtraArgs     `json:"extraArgs"`
	InputHoldingCids []CONTRACT_ID `json:"inputHoldingCids"`
}

// ToMap converts FeeInput to a map for DAML arguments
func (t FeeInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["transferFactory"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TransferFactory).(mapper); ok {
			return m.toMap()
		}
		return t.TransferFactory
	}()

	m["extraArgs"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

	m["inputHoldingCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for FeeInput using JsonCodec
func (t FeeInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for FeeInput using JsonCodec
func (t *FeeInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// LockOrBurnResult is a Record type
type LockOrBurnResult struct {
	PoolChangeCids   []CONTRACT_ID `json:"poolChangeCids"`
	SenderChangeCids []CONTRACT_ID `json:"senderChangeCids"`
}

// ToMap converts LockOrBurnResult to a map for DAML arguments
func (t LockOrBurnResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["poolChangeCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.PoolChangeCids))
		for _, e := range t.PoolChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["senderChangeCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderChangeCids))
		for _, e := range t.SenderChangeCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for LockOrBurnResult using JsonCodec
func (t LockOrBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for LockOrBurnResult using JsonCodec
func (t *LockOrBurnResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ReleaseOrMintResult is a Record type
type ReleaseOrMintResult struct {
	Output         ReleaseOrMintResultOutput `json:"output"`
	PoolChangeCids []CONTRACT_ID             `json:"poolChangeCids"`
}

// ToMap converts ReleaseOrMintResult to a map for DAML arguments
func (t ReleaseOrMintResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["output"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Output).(mapper); ok {
			return m.toMap()
		}
		return t.Output
	}()

	m["poolChangeCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.PoolChangeCids))
		for _, e := range t.PoolChangeCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ReleaseOrMintResult using JsonCodec
func (t ReleaseOrMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ReleaseOrMintResult using JsonCodec
func (t *ReleaseOrMintResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ReleaseOrMintResultCompleted is a Record type
type ReleaseOrMintResultCompleted struct {
	ReceiverHoldingCids []CONTRACT_ID `json:"receiverHoldingCids"`
}

// ToMap converts ReleaseOrMintResultCompleted to a map for DAML arguments
func (t ReleaseOrMintResultCompleted) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["receiverHoldingCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverHoldingCids))
		for _, e := range t.ReceiverHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ReleaseOrMintResultCompleted using JsonCodec
func (t ReleaseOrMintResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ReleaseOrMintResultCompleted using JsonCodec
func (t *ReleaseOrMintResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ReleaseOrMintResultOutput is a variant/union type
type ReleaseOrMintResultOutput struct {
	ReleaseOrMintResultPending   *ReleaseOrMintResultPending   `json:"ReleaseOrMintResult_Pending,omitempty"`
	ReleaseOrMintResultCompleted *ReleaseOrMintResultCompleted `json:"ReleaseOrMintResult_Completed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ReleaseOrMintResultOutput
func (v ReleaseOrMintResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(v)
}

// UnmarshalJSON implements custom JSON unmarshaling for ReleaseOrMintResultOutput
func (v *ReleaseOrMintResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ReleaseOrMintResultOutput) GetVariantTag() string {

	if v.ReleaseOrMintResultPending != nil {
		return "ReleaseOrMintResult_Pending"
	}

	if v.ReleaseOrMintResultCompleted != nil {
		return "ReleaseOrMintResult_Completed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ReleaseOrMintResultOutput) GetVariantValue() interface{} {

	if v.ReleaseOrMintResultPending != nil {
		return v.ReleaseOrMintResultPending
	}

	if v.ReleaseOrMintResultCompleted != nil {
		return v.ReleaseOrMintResultCompleted
	}

	return nil
}

// Verify interface implementation
var _ VARIANT = (*ReleaseOrMintResultOutput)(nil)

// ReleaseOrMintResultPending is a Record type
type ReleaseOrMintResultPending struct {
	TransferInstructionCid CONTRACT_ID `json:"transferInstructionCid"`
}

// ToMap converts ReleaseOrMintResultPending to a map for DAML arguments
func (t ReleaseOrMintResultPending) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferInstructionCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TransferInstructionCid).(mapper); ok {
			return m.toMap()
		}
		return t.TransferInstructionCid
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for ReleaseOrMintResultPending using JsonCodec
func (t ReleaseOrMintResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ReleaseOrMintResultPending using JsonCodec
func (t *ReleaseOrMintResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenInput is a Record type
type TokenInput struct {
	TransferFactory   CONTRACT_ID   `json:"transferFactory"`
	ExtraArgs         ExtraArgs     `json:"extraArgs"`
	TokenPoolHoldings []CONTRACT_ID `json:"tokenPoolHoldings"`
}

// ToMap converts TokenInput to a map for DAML arguments
func (t TokenInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["transferFactory"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TransferFactory).(mapper); ok {
			return m.toMap()
		}
		return t.TransferFactory
	}()

	m["extraArgs"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExtraArgs).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraArgs
	}()

	m["tokenPoolHoldings"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.TokenPoolHoldings))
		for _, e := range t.TokenPoolHoldings {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenInput using JsonCodec
func (t TokenInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenInput using JsonCodec
func (t *TokenInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolView is a Record type
type TokenPoolView struct {
	Owner        PARTY        `json:"owner"`
	CcipOwner    PARTY        `json:"ccipOwner"`
	InstrumentId InstrumentId `json:"instrumentId"`
}

// ToMap converts TokenPoolView to a map for DAML arguments
func (t TokenPoolView) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["owner"] = t.Owner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instrumentId"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenPoolView using JsonCodec
func (t TokenPoolView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolView using JsonCodec
func (t *TokenPoolView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolGetRequiredCCVs is a Record type
type TokenPoolGetRequiredCCVs struct {
	RemoteChainSelector NUMERIC           `json:"remoteChainSelector"`
	Amount              NUMERIC           `json:"amount"`
	Finality            INT64             `json:"finality"`
	ExtraData           TEXT              `json:"extraData"`
	Direction           TransferDirection `json:"direction"`
	Caller              PARTY             `json:"caller"`
}

// ToMap converts TokenPoolGetRequiredCCVs to a map for DAML arguments
func (t TokenPoolGetRequiredCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["remoteChainSelector"] = (*big.Int)(t.RemoteChainSelector)

	m["amount"] = (*big.Int)(t.Amount)

	m["finality"] = int64(t.Finality)

	m["extraData"] = string(t.ExtraData)

	m["direction"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Direction).(mapper); ok {
			return m.toMap()
		}
		return t.Direction
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenPoolGetRequiredCCVs using JsonCodec
func (t TokenPoolGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolGetRequiredCCVs using JsonCodec
func (t *TokenPoolGetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolLockOrBurn is a Record type
type TokenPoolLockOrBurn struct {
	DestChainSelector     NUMERIC           `json:"destChainSelector"`
	Message               Canton2AnyMessage `json:"message"`
	TokenInput            TokenInput        `json:"tokenInput"`
	FeeInput              FeeInput          `json:"feeInput"`
	SenderInputCids       []CONTRACT_ID     `json:"senderInputCids"`
	OnRampCid             CONTRACT_ID       `json:"onRampCid"`
	FeeQuoterCid          CONTRACT_ID       `json:"feeQuoterCid"`
	TokenAdminRegistryCid CONTRACT_ID       `json:"tokenAdminRegistryCid"`
	Caller                PARTY             `json:"caller"`
}

// ToMap converts TokenPoolLockOrBurn to a map for DAML arguments
func (t TokenPoolLockOrBurn) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["feeInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeInput).(mapper); ok {
			return m.toMap()
		}
		return t.FeeInput
	}()

	m["senderInputCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["onRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OnRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OnRampCid
	}()

	m["feeQuoterCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenPoolLockOrBurn using JsonCodec
func (t TokenPoolLockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolLockOrBurn using JsonCodec
func (t *TokenPoolLockOrBurn) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolReleaseFromTicket is a Record type
type TokenPoolReleaseFromTicket struct {
	TokenReceiveTicketCid CONTRACT_ID `json:"tokenReceiveTicketCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
	TokenInput            TokenInput  `json:"tokenInput"`
	Caller                PARTY       `json:"caller"`
}

// ToMap converts TokenPoolReleaseFromTicket to a map for DAML arguments
func (t TokenPoolReleaseFromTicket) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["tokenReceiveTicketCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenReceiveTicketCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenReceiveTicketCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenPoolReleaseFromTicket using JsonCodec
func (t TokenPoolReleaseFromTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolReleaseFromTicket using JsonCodec
func (t *TokenPoolReleaseFromTicket) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenPoolVerifyCCVs is a Record type
type TokenPoolVerifyCCVs struct {
	CcvVerifyTickets    []CONTRACT_ID `json:"ccvVerifyTickets"`
	MessageHash         TEXT          `json:"messageHash"`
	SourceChainSelector NUMERIC       `json:"sourceChainSelector"`
	Amount              NUMERIC       `json:"amount"`
	Receiver            PARTY         `json:"receiver"`
	Caller              PARTY         `json:"caller"`
}

// ToMap converts TokenPoolVerifyCCVs to a map for DAML arguments
func (t TokenPoolVerifyCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvVerifyTickets"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvVerifyTickets))
		for _, e := range t.CcvVerifyTickets {
			res = append(res, e)
		}
		return res
	}()

	m["messageHash"] = string(t.MessageHash)

	m["sourceChainSelector"] = (*big.Int)(t.SourceChainSelector)

	m["amount"] = (*big.Int)(t.Amount)

	m["receiver"] = t.Receiver.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

// MarshalJSON implements custom JSON marshaling for TokenPoolVerifyCCVs using JsonCodec
func (t TokenPoolVerifyCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolVerifyCCVs using JsonCodec
func (t *TokenPoolVerifyCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferDirection is an enum type
type TransferDirection string

const (
	TransferDirectionOutbound TransferDirection = "Outbound"
	TransferDirectionInbound  TransferDirection = "Inbound"
)

// GetEnumConstructor implements types.ENUM interface
func (e TransferDirection) GetEnumConstructor() string {
	return string(e)
}

// GetEnumTypeID implements types.ENUM interface
func (e TransferDirection) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.Interfaces.TokenPool", "TransferDirection")
}

// MarshalJSON implements custom JSON marshaling for TransferDirection using JsonCodec
func (e TransferDirection) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

// UnmarshalJSON implements custom JSON unmarshaling for TransferDirection using JsonCodec
func (e *TransferDirection) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

// Verify interface implementation
var _ ENUM = TransferDirection("")

// IITokenPoolInterfaceID returns the interface ID for the IITokenPool interface
func IITokenPoolInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "CCIP.Interfaces.TokenPool", "ITokenPool")
}
