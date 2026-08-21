package splice_api_token_allocation_instruction_v1

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	splice_api_token_allocation_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_allocation_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "splice-api-token-allocation-instruction-v1"
	PackageID   = "275064aacfe99cea72ee0c80563936129563776f67415ef9f13e4297eecbc520"
	SDKVersion  = "3.3.0-snapshot.20250502.13767.0.v2fc6c7e2"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IAllocationFactory is a DAML interface
type IAllocationFactory interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationFactoryAllocate executes the AllocationFactory_Allocate choice
	AllocationFactoryAllocate(contractID string, args AllocationFactoryAllocate) *model.ExerciseCommand

	// AllocationFactoryPublicFetch executes the AllocationFactory_PublicFetch choice
	AllocationFactoryPublicFetch(contractID string, args AllocationFactoryPublicFetch) *model.ExerciseCommand
}

// IAllocationInstruction is a DAML interface
type IAllocationInstruction interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// AllocationInstructionWithdraw executes the AllocationInstruction_Withdraw choice
	AllocationInstructionWithdraw(contractID string, args AllocationInstructionWithdraw) *model.ExerciseCommand

	// AllocationInstructionUpdate executes the AllocationInstruction_Update choice
	AllocationInstructionUpdate(contractID string, args AllocationInstructionUpdate) *model.ExerciseCommand
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

// AllocationFactoryView is a Record type
type AllocationFactoryView struct {
	Admin types.PARTY                           `json:"admin"`
	Meta  splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationFactoryView to a map for DAML arguments
func (t AllocationFactoryView) ToMap() map[string]any {
	m := make(map[string]any)

	m["admin"] = t.Admin.ToMap()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationFactoryView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryView to hex string (Canton MCMS format)
func (t AllocationFactoryView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryView from hex string (Canton MCMS format)
func (t *AllocationFactoryView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryAllocate is a Record type
type AllocationFactoryAllocate struct {
	ExpectedAdmin    types.PARTY                                            `json:"expectedAdmin"`
	Allocation       splice_api_token_allocation_v1.AllocationSpecification `json:"allocation"`
	RequestedAt      types.TIMESTAMP                                        `json:"requestedAt"`
	InputHoldingCids []types.CONTRACT_ID                                    `json:"inputHoldingCids"`
	ExtraArgs        splice_api_token_metadata_v1.ExtraArgs                 `json:"extraArgs"`
}

// ToMap converts AllocationFactoryAllocate to a map for DAML arguments
func (t AllocationFactoryAllocate) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["allocation"] = model.NestedToDAMLValue(t.Allocation)

	m["requestedAt"] = t.RequestedAt

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationFactoryAllocate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryAllocate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryAllocate to hex string (Canton MCMS format)
func (t AllocationFactoryAllocate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryAllocate from hex string (Canton MCMS format)
func (t *AllocationFactoryAllocate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationFactoryPublicFetch is a Record type
type AllocationFactoryPublicFetch struct {
	ExpectedAdmin types.PARTY `json:"expectedAdmin"`
	Actor         types.PARTY `json:"actor"`
}

// ToMap converts AllocationFactoryPublicFetch to a map for DAML arguments
func (t AllocationFactoryPublicFetch) ToMap() map[string]any {
	m := make(map[string]any)

	m["expectedAdmin"] = t.ExpectedAdmin.ToMap()

	m["actor"] = t.Actor.ToMap()

	return m
}

func (t AllocationFactoryPublicFetch) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationFactoryPublicFetch) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationFactoryPublicFetch to hex string (Canton MCMS format)
func (t AllocationFactoryPublicFetch) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationFactoryPublicFetch from hex string (Canton MCMS format)
func (t *AllocationFactoryPublicFetch) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionResult is a Record type
type AllocationInstructionResult struct {
	Output           AllocationInstructionResultOutput     `json:"output"`
	SenderChangeCids []types.CONTRACT_ID                   `json:"senderChangeCids"`
	Meta             splice_api_token_metadata_v1.Metadata `json:"meta"`
}

// ToMap converts AllocationInstructionResult to a map for DAML arguments
func (t AllocationInstructionResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["output"] = model.NestedToDAMLValue(t.Output)

	m["senderChangeCids"] = func() []any {
		res := make([]any, 0, len(t.SenderChangeCids))
		for _, e := range t.SenderChangeCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationInstructionResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionResult to hex string (Canton MCMS format)
func (t AllocationInstructionResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionResult from hex string (Canton MCMS format)
func (t *AllocationInstructionResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionResultCompleted is a Record type
type AllocationInstructionResultCompleted struct {
	AllocationCid types.CONTRACT_ID `json:"allocationCid"`
}

// ToMap converts AllocationInstructionResultCompleted to a map for DAML arguments
func (t AllocationInstructionResultCompleted) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocationCid"] = model.NestedToDAMLValue(t.AllocationCid)

	return m
}

func (t AllocationInstructionResultCompleted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionResultCompleted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionResultCompleted to hex string (Canton MCMS format)
func (t AllocationInstructionResultCompleted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionResultCompleted from hex string (Canton MCMS format)
func (t *AllocationInstructionResultCompleted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionResultOutput is a variant/union type
type AllocationInstructionResultOutput struct {
	AllocationInstructionResultPending   *AllocationInstructionResultPending   `json:"AllocationInstructionResult_Pending,omitempty"`
	AllocationInstructionResultCompleted *AllocationInstructionResultCompleted `json:"AllocationInstructionResult_Completed,omitempty"`
	AllocationInstructionResultFailed    *types.UNIT                           `json:"AllocationInstructionResult_Failed,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AllocationInstructionResultOutput
func (v AllocationInstructionResultOutput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for AllocationInstructionResultOutput
func (v *AllocationInstructionResultOutput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes AllocationInstructionResultOutput to hex string (Canton MCMS format)
func (v AllocationInstructionResultOutput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes AllocationInstructionResultOutput from hex string (Canton MCMS format)
func (v *AllocationInstructionResultOutput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v AllocationInstructionResultOutput) GetVariantTag() string {

	if v.AllocationInstructionResultPending != nil {
		return "AllocationInstructionResult_Pending"
	}

	if v.AllocationInstructionResultCompleted != nil {
		return "AllocationInstructionResult_Completed"
	}

	if v.AllocationInstructionResultFailed != nil {
		return "AllocationInstructionResult_Failed"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v AllocationInstructionResultOutput) GetVariantValue() any {

	if v.AllocationInstructionResultPending != nil {
		return v.AllocationInstructionResultPending
	}

	if v.AllocationInstructionResultCompleted != nil {
		return v.AllocationInstructionResultCompleted
	}

	if v.AllocationInstructionResultFailed != nil {
		return v.AllocationInstructionResultFailed
	}

	return nil
}

var _ types.VARIANT = (*AllocationInstructionResultOutput)(nil)

// AllocationInstructionResultPending is a Record type
type AllocationInstructionResultPending struct {
	AllocationInstructionCid types.CONTRACT_ID `json:"allocationInstructionCid"`
}

// ToMap converts AllocationInstructionResultPending to a map for DAML arguments
func (t AllocationInstructionResultPending) ToMap() map[string]any {
	m := make(map[string]any)

	m["allocationInstructionCid"] = model.NestedToDAMLValue(t.AllocationInstructionCid)

	return m
}

func (t AllocationInstructionResultPending) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionResultPending) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionResultPending to hex string (Canton MCMS format)
func (t AllocationInstructionResultPending) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionResultPending from hex string (Canton MCMS format)
func (t *AllocationInstructionResultPending) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionView is a Record type
type AllocationInstructionView struct {
	OriginalInstructionCid *types.CONTRACT_ID                                     `json:"originalInstructionCid" hex:"optional"`
	Allocation             splice_api_token_allocation_v1.AllocationSpecification `json:"allocation"`
	PendingActions         map[types.PARTY]types.TEXT                             `json:"pendingActions"`
	RequestedAt            types.TIMESTAMP                                        `json:"requestedAt"`
	InputHoldingCids       []types.CONTRACT_ID                                    `json:"inputHoldingCids"`
	Meta                   splice_api_token_metadata_v1.Metadata                  `json:"meta"`
}

// ToMap converts AllocationInstructionView to a map for DAML arguments
func (t AllocationInstructionView) ToMap() map[string]any {
	m := make(map[string]any)

	if t.OriginalInstructionCid != nil {
		m["originalInstructionCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OriginalInstructionCid),
		}
	} else {
		m["originalInstructionCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["allocation"] = model.NestedToDAMLValue(t.Allocation)

	m["pendingActions"] = func() any {
		if t.PendingActions == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.PendingActions}
	}()

	m["requestedAt"] = t.RequestedAt

	m["inputHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.InputHoldingCids))
		for _, e := range t.InputHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["meta"] = model.NestedToDAMLValue(t.Meta)

	return m
}

func (t AllocationInstructionView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionView to hex string (Canton MCMS format)
func (t AllocationInstructionView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionView from hex string (Canton MCMS format)
func (t *AllocationInstructionView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionUpdate is a Record type
type AllocationInstructionUpdate struct {
	ExtraActors []types.PARTY                          `json:"extraActors"`
	ExtraArgs   splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationInstructionUpdate to a map for DAML arguments
func (t AllocationInstructionUpdate) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraActors"] = func() []any {
		res := make([]any, 0, len(t.ExtraActors))
		for _, e := range t.ExtraActors {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationInstructionUpdate) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionUpdate) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionUpdate to hex string (Canton MCMS format)
func (t AllocationInstructionUpdate) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionUpdate from hex string (Canton MCMS format)
func (t *AllocationInstructionUpdate) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AllocationInstructionWithdraw is a Record type
type AllocationInstructionWithdraw struct {
	ExtraArgs splice_api_token_metadata_v1.ExtraArgs `json:"extraArgs"`
}

// ToMap converts AllocationInstructionWithdraw to a map for DAML arguments
func (t AllocationInstructionWithdraw) ToMap() map[string]any {
	m := make(map[string]any)

	m["extraArgs"] = model.NestedToDAMLValue(t.ExtraArgs)

	return m
}

func (t AllocationInstructionWithdraw) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AllocationInstructionWithdraw) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AllocationInstructionWithdraw to hex string (Canton MCMS format)
func (t AllocationInstructionWithdraw) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AllocationInstructionWithdraw from hex string (Canton MCMS format)
func (t *AllocationInstructionWithdraw) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IAllocationFactoryInterfaceID returns the interface ID for the IAllocationFactory interface using the package name
func IAllocationFactoryInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.AllocationInstructionV1", "AllocationFactory")
}

// IAllocationFactoryInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IAllocationFactoryInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.AllocationInstructionV1", "AllocationFactory")
}

// IAllocationInstructionInterfaceID returns the interface ID for the IAllocationInstruction interface using the package name
func IAllocationInstructionInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "Splice.Api.Token.AllocationInstructionV1", "AllocationInstruction")
}

// IAllocationInstructionInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IAllocationInstructionInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "Splice.Api.Token.AllocationInstructionV1", "AllocationInstruction")
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
