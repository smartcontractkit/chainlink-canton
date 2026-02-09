package interfaces

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageName = "ccip-tokenpool-interfaces"
const SDKVersion = "3.4.10"

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

	type mapper interface {
		ToMap() map[string]interface{}
	}
	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{"args": args}
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

func (t FeeInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t LockOrBurnResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t ReleaseOrMintResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t ReleaseOrMintResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t ReleaseOrMintResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t TokenInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t TokenPoolView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["amount"] = t.Amount

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

func (t TokenPoolGetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

	m["destChainSelector"] = t.DestChainSelector

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

func (t TokenPoolLockOrBurn) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

func (t TokenPoolReleaseFromTicket) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

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

	m["sourceChainSelector"] = t.SourceChainSelector

	m["amount"] = t.Amount

	m["receiver"] = t.Receiver.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolVerifyCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenPoolVerifyCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TransferDirection is an enum type
type TransferDirection string

const (
	TransferDirectionOutbound TransferDirection = "Outbound"

	TransferDirectionInbound TransferDirection = "Inbound"
)

func (e TransferDirection) GetEnumConstructor() string { return string(e) }

func (e TransferDirection) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.TokenPool", "TransferDirection")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e TransferDirection) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Interfaces.TokenPool", "TransferDirection")
}

func (e TransferDirection) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(e)
}

func (e *TransferDirection) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, e)
}

var _ ENUM = TransferDirection("")

// IITokenPoolInterfaceID returns the interface ID for the IITokenPool interface using the package name
func IITokenPoolInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Interfaces.TokenPool", "ITokenPool")
}

// IITokenPoolInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IITokenPoolInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Interfaces.TokenPool", "ITokenPool")
}
