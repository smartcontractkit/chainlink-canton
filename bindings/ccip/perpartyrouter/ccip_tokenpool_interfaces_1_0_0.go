package perpartyrouter

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
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
)

// IITokenPool is a DAML interface
type IITokenPool interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// TokenPoolGetRequiredCCVs executes the TokenPool_GetRequiredCCVs choice
	TokenPoolGetRequiredCCVs(contractID string, args TokenPoolGetRequiredCCVs) *model.ExerciseCommand

	// TokenPoolVerifyInboundCCVs executes the TokenPool_VerifyInboundCCVs choice
	TokenPoolVerifyInboundCCVs(contractID string, args TokenPoolVerifyInboundCCVs) *model.ExerciseCommand

	// TokenPoolReleaseFromTicket executes the TokenPool_ReleaseFromTicket choice
	TokenPoolReleaseFromTicket(contractID string, args TokenPoolReleaseFromTicket) *model.ExerciseCommand

	// TokenPoolLockOrBurn executes the TokenPool_LockOrBurn choice
	TokenPoolLockOrBurn(contractID string, args TokenPoolLockOrBurn) *model.ExerciseCommand

	// TokenPoolCalculateFee executes the TokenPool_CalculateFee choice
	TokenPoolCalculateFee(contractID string, args TokenPoolCalculateFee) *model.ExerciseCommand
}

// LockOrBurnResult is a Record type
type LockOrBurnResult struct {
	PoolChangeCids    []CONTRACT_ID `json:"poolChangeCids"`
	SenderChangeCids  []CONTRACT_ID `json:"senderChangeCids"`
	SendingMessageCid CONTRACT_ID   `json:"sendingMessageCid"`
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

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
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

// TokenPoolCalculateFee is a Record type
type TokenPoolCalculateFee struct {
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	Caller            PARTY       `json:"caller"`
}

// ToMap converts TokenPoolCalculateFee to a map for DAML arguments
func (t TokenPoolCalculateFee) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t TokenPoolCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenPoolCalculateFee) UnmarshalJSON(data []byte) error {
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
	SendingMessageCid CONTRACT_ID   `json:"sendingMessageCid"`
	TokenInput        TokenInput    `json:"tokenInput"`
	SenderInputCids   []CONTRACT_ID `json:"senderInputCids"`
	Amount            NUMERIC       `json:"amount"`
	RmnRemoteCid      CONTRACT_ID   `json:"rmnRemoteCid"`
	Caller            PARTY         `json:"caller"`
}

// ToMap converts TokenPoolLockOrBurn to a map for DAML arguments
func (t TokenPoolLockOrBurn) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["senderInputCids"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderInputCids))
		for _, e := range t.SenderInputCids {
			res = append(res, e)
		}
		return res
	}()

	m["amount"] = t.Amount

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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
	RmnRemoteCid          CONTRACT_ID `json:"rmnRemoteCid"`
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

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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

// TokenPoolVerifyInboundCCVs is a Record type
type TokenPoolVerifyInboundCCVs struct {
	ExecutingMessageCid   CONTRACT_ID `json:"executingMessageCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
	Caller                PARTY       `json:"caller"`
}

// ToMap converts TokenPoolVerifyInboundCCVs to a map for DAML arguments
func (t TokenPoolVerifyInboundCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
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

func (t TokenPoolVerifyInboundCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenPoolVerifyInboundCCVs) UnmarshalJSON(data []byte) error {
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
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Interfaces.TokenPool", "ITokenPool")
}
