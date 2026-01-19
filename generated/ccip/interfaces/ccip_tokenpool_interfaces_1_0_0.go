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

const PackageID = "4c09d18cf98fd6d5c174a951ffff4fc580fa1aa96f1897e755ba603d80e4cc98"
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

	// TokenPoolReleaseOrMint executes the TokenPool_ReleaseOrMint choice
	TokenPoolReleaseOrMint(contractID string, args TokenPoolReleaseOrMint) *model.ExerciseCommand

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

// FeeInput is a Record type
type FeeInput struct {
	InstrumentId     InstrumentId  `json:"instrumentId"`
	TransferFactory  CONTRACT_ID   `json:"transferFactory"`
	ExtraArgs        ExtraArgs     `json:"extraArgs"`
	InputHoldingCids []CONTRACT_ID `json:"inputHoldingCids"`
}

// toMap converts FeeInput to a map for DAML arguments
func (t FeeInput) toMap() map[string]interface{} {
	return map[string]interface{}{

		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
		"transferFactory": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TransferFactory).(mapper); ok {
				return m.toMap()
			}
			return t.TransferFactory
		}(),
		"extraArgs": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.ExtraArgs).(mapper); ok {
				return m.toMap()
			}
			return t.ExtraArgs
		}(),
		"inputHoldingCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.InputHoldingCids))
			for _, e := range t.InputHoldingCids {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts LockOrBurnResult to a map for DAML arguments
func (t LockOrBurnResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"poolChangeCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.PoolChangeCids))
			for _, e := range t.PoolChangeCids {
				res = append(res, e)
			}
			return res
		}(),
		"senderChangeCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.SenderChangeCids))
			for _, e := range t.SenderChangeCids {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts ReleaseOrMintResult to a map for DAML arguments
func (t ReleaseOrMintResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"output": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Output).(mapper); ok {
				return m.toMap()
			}
			return t.Output
		}(),
		"poolChangeCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.PoolChangeCids))
			for _, e := range t.PoolChangeCids {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts ReleaseOrMintResultCompleted to a map for DAML arguments
func (t ReleaseOrMintResultCompleted) toMap() map[string]interface{} {
	return map[string]interface{}{

		"receiverHoldingCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.ReceiverHoldingCids))
			for _, e := range t.ReceiverHoldingCids {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts ReleaseOrMintResultPending to a map for DAML arguments
func (t ReleaseOrMintResultPending) toMap() map[string]interface{} {
	return map[string]interface{}{

		"transferInstructionCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TransferInstructionCid).(mapper); ok {
				return m.toMap()
			}
			return t.TransferInstructionCid
		}(),
	}
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

// toMap converts TokenInput to a map for DAML arguments
func (t TokenInput) toMap() map[string]interface{} {
	return map[string]interface{}{

		"transferFactory": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TransferFactory).(mapper); ok {
				return m.toMap()
			}
			return t.TransferFactory
		}(),
		"extraArgs": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.ExtraArgs).(mapper); ok {
				return m.toMap()
			}
			return t.ExtraArgs
		}(),
		"tokenPoolHoldings": func() []interface{} {
			res := make([]interface{}, 0, len(t.TokenPoolHoldings))
			for _, e := range t.TokenPoolHoldings {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts TokenPoolView to a map for DAML arguments
func (t TokenPoolView) toMap() map[string]interface{} {
	return map[string]interface{}{

		"owner":     t.Owner.ToMap(),
		"ccipOwner": t.CcipOwner.ToMap(),
		"instrumentId": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.InstrumentId).(mapper); ok {
				return m.toMap()
			}
			return t.InstrumentId
		}(),
	}
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
	Message MessageV1 `json:"message"`
	Caller  PARTY     `json:"caller"`
}

// toMap converts TokenPoolGetRequiredCCVs to a map for DAML arguments
func (t TokenPoolGetRequiredCCVs) toMap() map[string]interface{} {
	return map[string]interface{}{

		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"caller": t.Caller.ToMap(),
	}
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

// toMap converts TokenPoolLockOrBurn to a map for DAML arguments
func (t TokenPoolLockOrBurn) toMap() map[string]interface{} {
	return map[string]interface{}{

		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"tokenInput": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenInput).(mapper); ok {
				return m.toMap()
			}
			return t.TokenInput
		}(),
		"feeInput": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeInput).(mapper); ok {
				return m.toMap()
			}
			return t.FeeInput
		}(),
		"senderInputCids": func() []interface{} {
			res := make([]interface{}, 0, len(t.SenderInputCids))
			for _, e := range t.SenderInputCids {
				res = append(res, e)
			}
			return res
		}(),
		"onRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OnRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OnRampCid
		}(),
		"feeQuoterCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeQuoterCid).(mapper); ok {
				return m.toMap()
			}
			return t.FeeQuoterCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"caller": t.Caller.ToMap(),
	}
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

// TokenPoolReleaseOrMint is a Record type
type TokenPoolReleaseOrMint struct {
	EncodedMessage        TEXT          `json:"encodedMessage"`
	Ccvs                  []CONTRACT_ID `json:"ccvs"`
	CcvData               []TEXT        `json:"ccvData"`
	TokenInput            TokenInput    `json:"tokenInput"`
	OffRampCid            CONTRACT_ID   `json:"offRampCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	Caller                PARTY         `json:"caller"`
}

// toMap converts TokenPoolReleaseOrMint to a map for DAML arguments
func (t TokenPoolReleaseOrMint) toMap() map[string]interface{} {
	return map[string]interface{}{

		"encodedMessage": string(t.EncodedMessage),
		"ccvs": func() []interface{} {
			res := make([]interface{}, 0, len(t.Ccvs))
			for _, e := range t.Ccvs {
				res = append(res, e)
			}
			return res
		}(),
		"ccvData": func() []interface{} {
			res := make([]interface{}, 0, len(t.CcvData))
			for _, e := range t.CcvData {
				res = append(res, string(e))
			}
			return res
		}(),
		"tokenInput": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenInput).(mapper); ok {
				return m.toMap()
			}
			return t.TokenInput
		}(),
		"offRampCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.OffRampCid).(mapper); ok {
				return m.toMap()
			}
			return t.OffRampCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"caller": t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for TokenPoolReleaseOrMint using JsonCodec
func (t TokenPoolReleaseOrMint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for TokenPoolReleaseOrMint using JsonCodec
func (t *TokenPoolReleaseOrMint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// IITokenPoolInterfaceID returns the interface ID for the IITokenPool interface
func IITokenPoolInterfaceID(packageID *string) string {
	pkgID := PackageID
	if packageID != nil {
		pkgID = *packageID
	}
	return fmt.Sprintf("#%s:%s:%s", pkgID, "CCIP.Interfaces.TokenPool", "ITokenPool")
}
