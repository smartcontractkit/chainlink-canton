package api

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	PackageName = "mcms-api"
	PackageID   = "72e12cf8b6aa3db9a2e7c325dfe36a56d4268a4b66b6fc73123ffe520614faba"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

// IMCMSReceiver is a DAML interface
type IMCMSReceiver interface {

	// Archive executes the Archive choice
	Archive(contractID string) *model.ExerciseCommand

	// MCMSReceiverEntrypoint executes the MCMSReceiver_Entrypoint choice
	MCMSReceiverEntrypoint(contractID string, args MCMSReceiverEntrypoint) *model.ExerciseCommand
}

const (
	ZeroHash      = types.TEXT("0000000000000000000000000000000000000000000000000000000000000000")
	NumGroups     = types.INT64(32)
	MaxNumSigners = types.INT64(200)
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

// APSetConfig is a Record type
type APSetConfig struct {
	ApSigners      []SignerInfo  `json:"apSigners"`
	ApGroupQuorums []types.INT64 `json:"apGroupQuorums" hex:"[]uint32"`
	ApGroupParents []types.INT64 `json:"apGroupParents" hex:"[]uint32"`
	ApClearRoot    types.BOOL    `json:"apClearRoot"`
}

// ToMap converts APSetConfig to a map for DAML arguments
func (t APSetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["apSigners"] = func() []any {
		res := make([]any, 0, len(t.ApSigners))
		for _, e := range t.ApSigners {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["apGroupQuorums"] = func() []any {
		res := make([]any, 0, len(t.ApGroupQuorums))
		for _, e := range t.ApGroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apGroupParents"] = func() []any {
		res := make([]any, 0, len(t.ApGroupParents))
		for _, e := range t.ApGroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["apClearRoot"] = bool(t.ApClearRoot)

	return m
}

func (t APSetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *APSetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes APSetConfig to hex string (Canton MCMS format)
func (t APSetConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes APSetConfig from hex string (Canton MCMS format)
func (t *APSetConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AdminParams is a variant/union type
type AdminParams struct {
	APSetConfig *APSetConfig `json:"AP_SetConfig,omitempty"`
	APClearRoot *types.UNIT  `json:"AP_ClearRoot,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for AdminParams
func (v AdminParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for AdminParams
func (v *AdminParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes AdminParams to hex string (Canton MCMS format)
func (v AdminParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes AdminParams from hex string (Canton MCMS format)
func (v *AdminParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v AdminParams) GetVariantTag() string {

	if v.APSetConfig != nil {
		return "AP_SetConfig"
	}

	if v.APClearRoot != nil {
		return "AP_ClearRoot"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v AdminParams) GetVariantValue() any {

	if v.APSetConfig != nil {
		return v.APSetConfig
	}

	if v.APClearRoot != nil {
		return v.APClearRoot
	}

	return nil
}

var _ types.VARIANT = (*AdminParams)(nil)

// ArgValue is a variant/union type
type ArgValue struct {
	AVText  *types.TEXT      `json:"AV_Text,omitempty"`
	AVInt   *types.INT64     `json:"AV_Int,omitempty"`
	AVBool  *types.BOOL      `json:"AV_Bool,omitempty"`
	AVParty *types.PARTY     `json:"AV_Party,omitempty"`
	AVTime  *types.TIMESTAMP `json:"AV_Time,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ArgValue
func (v ArgValue) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(v)
}

// UnmarshalJSON implements custom JSON unmarshalling for ArgValue
func (v *ArgValue) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, v)
}

// MarshalHex encodes ArgValue to hex string (Canton MCMS format)
func (v ArgValue) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(v)
}

// UnmarshalHex decodes ArgValue from hex string (Canton MCMS format)
func (v *ArgValue) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, v)
}

// GetVariantTag implements types.VARIANT interface
func (v ArgValue) GetVariantTag() string {

	if v.AVText != nil {
		return "AV_Text"
	}

	if v.AVInt != nil {
		return "AV_Int"
	}

	if v.AVBool != nil {
		return "AV_Bool"
	}

	if v.AVParty != nil {
		return "AV_Party"
	}

	if v.AVTime != nil {
		return "AV_Time"
	}

	return ""
}

// GetVariantValue implements types.VARIANT interface
func (v ArgValue) GetVariantValue() any {

	if v.AVText != nil {
		return v.AVText
	}

	if v.AVInt != nil {
		return v.AVInt
	}

	if v.AVBool != nil {
		return v.AVBool
	}

	if v.AVParty != nil {
		return v.AVParty
	}

	if v.AVTime != nil {
		return v.AVTime
	}

	return nil
}

var _ types.VARIANT = (*ArgValue)(nil)

// BlockedFunction is a Record type
type BlockedFunction struct {
	TargetInstanceAddress types.TEXT `json:"targetInstanceAddress"`
	FunctionName          types.TEXT `json:"functionName"`
}

// ToMap converts BlockedFunction to a map for DAML arguments
func (t BlockedFunction) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

	m["functionName"] = string(t.FunctionName)

	return m
}

func (t BlockedFunction) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BlockedFunction) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BlockedFunction to hex string (Canton MCMS format)
func (t BlockedFunction) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BlockedFunction from hex string (Canton MCMS format)
func (t *BlockedFunction) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// BypasserExecuteBatchParams is a Record type
type BypasserExecuteBatchParams struct {
	Calls []TimelockCall `json:"calls"`
}

// ToMap converts BypasserExecuteBatchParams to a map for DAML arguments
func (t BypasserExecuteBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t BypasserExecuteBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *BypasserExecuteBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes BypasserExecuteBatchParams to hex string (Canton MCMS format)
func (t BypasserExecuteBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes BypasserExecuteBatchParams from hex string (Canton MCMS format)
func (t *BypasserExecuteBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CancelBatchParams is a Record type
type CancelBatchParams struct {
	OpId types.TEXT `json:"opId"`
}

// ToMap converts CancelBatchParams to a map for DAML arguments
func (t CancelBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["opId"] = string(t.OpId)

	return m
}

func (t CancelBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CancelBatchParams to hex string (Canton MCMS format)
func (t CancelBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelBatchParams from hex string (Canton MCMS format)
func (t *CancelBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExpiringRoot is a Record type
type ExpiringRoot struct {
	Root       types.TEXT      `json:"root" hex:"bytes"`
	ValidUntil types.TIMESTAMP `json:"validUntil"`
	OpCount    types.INT64     `json:"opCount"`
}

// ToMap converts ExpiringRoot to a map for DAML arguments
func (t ExpiringRoot) ToMap() map[string]any {
	m := make(map[string]any)

	m["root"] = string(t.Root)

	m["validUntil"] = t.ValidUntil

	m["opCount"] = int64(t.OpCount)

	return m
}

func (t ExpiringRoot) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExpiringRoot) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExpiringRoot to hex string (Canton MCMS format)
func (t ExpiringRoot) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExpiringRoot from hex string (Canton MCMS format)
func (t *ExpiringRoot) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverView is a Record type
type MCMSReceiverView struct {
	McmsController types.PARTY `json:"mcmsController"`
	InstanceId     types.TEXT  `json:"instanceId"`
}

// ToMap converts MCMSReceiverView to a map for DAML arguments
func (t MCMSReceiverView) ToMap() map[string]any {
	m := make(map[string]any)

	m["mcmsController"] = t.McmsController.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t MCMSReceiverView) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverView) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverView to hex string (Canton MCMS format)
func (t MCMSReceiverView) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverView from hex string (Canton MCMS format)
func (t *MCMSReceiverView) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSReceiverEntrypoint is a Record type
type MCMSReceiverEntrypoint struct {
	FunctionName  types.TEXT                       `json:"functionName"`
	OperationData types.TEXT                       `json:"operationData" hex:"bytes16"`
	ContractIds   map[types.TEXT]types.CONTRACT_ID `json:"contractIds"`
}

// ToMap converts MCMSReceiverEntrypoint to a map for DAML arguments
func (t MCMSReceiverEntrypoint) ToMap() map[string]any {
	m := make(map[string]any)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	m["contractIds"] = func() any {
		if t.ContractIds == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ContractIds}
	}()

	return m
}

func (t MCMSReceiverEntrypoint) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MCMSReceiverEntrypoint) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MCMSReceiverEntrypoint to hex string (Canton MCMS format)
func (t MCMSReceiverEntrypoint) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MCMSReceiverEntrypoint from hex string (Canton MCMS format)
func (t *MCMSReceiverEntrypoint) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MultisigConfig is a Record type
type MultisigConfig struct {
	Signers      []SignerInfo  `json:"signers"`
	GroupQuorums []types.INT64 `json:"groupQuorums" hex:"[]uint32"`
	GroupParents []types.INT64 `json:"groupParents" hex:"[]uint32"`
}

// ToMap converts MultisigConfig to a map for DAML arguments
func (t MultisigConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["groupQuorums"] = func() []any {
		res := make([]any, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []any {
		res := make([]any, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	return m
}

func (t MultisigConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *MultisigConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes MultisigConfig to hex string (Canton MCMS format)
func (t MultisigConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes MultisigConfig from hex string (Canton MCMS format)
func (t *MultisigConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Op is a Record type
type Op struct {
	ChainId               types.INT64 `json:"chainId"`
	MultisigId            types.TEXT  `json:"multisigId"`
	Nonce                 types.INT64 `json:"nonce"`
	TargetInstanceAddress types.TEXT  `json:"targetInstanceAddress"`
	FunctionName          types.TEXT  `json:"functionName"`
	OperationData         types.TEXT  `json:"operationData" hex:"bytes16"`
}

// ToMap converts Op to a map for DAML arguments
func (t Op) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["nonce"] = int64(t.Nonce)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	return m
}

func (t Op) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Op) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Op to hex string (Canton MCMS format)
func (t Op) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Op from hex string (Canton MCMS format)
func (t *Op) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RawSignature is a Record type
type RawSignature struct {
	PublicKey types.TEXT `json:"publicKey"`
	R         types.TEXT `json:"r"`
	S         types.TEXT `json:"s"`
}

// ToMap converts RawSignature to a map for DAML arguments
func (t RawSignature) ToMap() map[string]any {
	m := make(map[string]any)

	m["publicKey"] = string(t.PublicKey)

	m["r"] = string(t.R)

	m["s"] = string(t.S)

	return m
}

func (t RawSignature) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RawSignature) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RawSignature to hex string (Canton MCMS format)
func (t RawSignature) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RawSignature from hex string (Canton MCMS format)
func (t *RawSignature) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Role is an enum type
type Role string

const (
	RoleBypasser Role = "Bypasser"

	RoleCanceller Role = "Canceller"

	RoleProposer Role = "Proposer"
)

func (e Role) GetEnumConstructor() string { return string(e) }

func (e Role) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.Types", "Role")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e Role) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "MCMS.Types", "Role")
}

func (e Role) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *Role) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes Role to hex string (Canton MCMS format)
func (e Role) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes Role from hex string (Canton MCMS format)
func (e *Role) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = Role("")

// RoleState is a Record type
type RoleState struct {
	Config       MultisigConfig                 `json:"config"`
	SeenHashes   map[types.TEXT]types.TIMESTAMP `json:"seenHashes"`
	ExpiringRoot ExpiringRoot                   `json:"expiringRoot"`
	RootMetadata RootMetadata                   `json:"rootMetadata"`
}

// ToMap converts RoleState to a map for DAML arguments
func (t RoleState) ToMap() map[string]any {
	m := make(map[string]any)

	m["config"] = model.NestedToDAMLValue(t.Config)

	m["seenHashes"] = func() any {
		if t.SeenHashes == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SeenHashes}
	}()

	m["expiringRoot"] = model.NestedToDAMLValue(t.ExpiringRoot)

	m["rootMetadata"] = model.NestedToDAMLValue(t.RootMetadata)

	return m
}

func (t RoleState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RoleState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RoleState to hex string (Canton MCMS format)
func (t RoleState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RoleState from hex string (Canton MCMS format)
func (t *RoleState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RootMetadata is a Record type
type RootMetadata struct {
	ChainId              types.INT64 `json:"chainId"`
	MultisigId           types.TEXT  `json:"multisigId"`
	PreOpCount           types.INT64 `json:"preOpCount"`
	PostOpCount          types.INT64 `json:"postOpCount"`
	OverridePreviousRoot types.BOOL  `json:"overridePreviousRoot"`
}

// ToMap converts RootMetadata to a map for DAML arguments
func (t RootMetadata) ToMap() map[string]any {
	m := make(map[string]any)

	m["chainId"] = int64(t.ChainId)

	m["multisigId"] = string(t.MultisigId)

	m["preOpCount"] = int64(t.PreOpCount)

	m["postOpCount"] = int64(t.PostOpCount)

	m["overridePreviousRoot"] = bool(t.OverridePreviousRoot)

	return m
}

func (t RootMetadata) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RootMetadata) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RootMetadata to hex string (Canton MCMS format)
func (t RootMetadata) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RootMetadata from hex string (Canton MCMS format)
func (t *RootMetadata) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ScheduleBatchParams is a Record type
type ScheduleBatchParams struct {
	Calls       []TimelockCall `json:"calls"`
	Predecessor types.TEXT     `json:"predecessor" hex:"bytes16"`
	Salt        types.TEXT     `json:"salt" hex:"bytes16"`
	DelaySecs   types.INT64    `json:"delaySecs"`
}

// ToMap converts ScheduleBatchParams to a map for DAML arguments
func (t ScheduleBatchParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["calls"] = func() []any {
		res := make([]any, 0, len(t.Calls))
		for _, e := range t.Calls {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["predecessor"] = string(t.Predecessor)

	m["salt"] = string(t.Salt)

	m["delaySecs"] = int64(t.DelaySecs)

	return m
}

func (t ScheduleBatchParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ScheduleBatchParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ScheduleBatchParams to hex string (Canton MCMS format)
func (t ScheduleBatchParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ScheduleBatchParams from hex string (Canton MCMS format)
func (t *ScheduleBatchParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetConfigParams is a Record type
type SetConfigParams struct {
	Signers      []SignerInfo  `json:"signers"`
	GroupQuorums []types.INT64 `json:"groupQuorums" hex:"[]uint32"`
	GroupParents []types.INT64 `json:"groupParents" hex:"[]uint32"`
	ClearRoot    types.BOOL    `json:"clearRoot"`
}

// ToMap converts SetConfigParams to a map for DAML arguments
func (t SetConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["groupQuorums"] = func() []any {
		res := make([]any, 0, len(t.GroupQuorums))
		for _, e := range t.GroupQuorums {
			res = append(res, int64(e))
		}
		return res
	}()

	m["groupParents"] = func() []any {
		res := make([]any, 0, len(t.GroupParents))
		for _, e := range t.GroupParents {
			res = append(res, int64(e))
		}
		return res
	}()

	m["clearRoot"] = bool(t.ClearRoot)

	return m
}

func (t SetConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetConfigParams to hex string (Canton MCMS format)
func (t SetConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetConfigParams from hex string (Canton MCMS format)
func (t *SetConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SignerInfo is a Record type
type SignerInfo struct {
	SignerAddress types.TEXT  `json:"signerAddress" hex:"bytes"`
	SignerIndex   types.INT64 `json:"signerIndex" hex:"uint32"`
	SignerGroup   types.INT64 `json:"signerGroup" hex:"uint32"`
}

// ToMap converts SignerInfo to a map for DAML arguments
func (t SignerInfo) ToMap() map[string]any {
	m := make(map[string]any)

	m["signerAddress"] = string(t.SignerAddress)

	m["signerIndex"] = int64(t.SignerIndex)

	m["signerGroup"] = int64(t.SignerGroup)

	return m
}

func (t SignerInfo) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SignerInfo) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SignerInfo to hex string (Canton MCMS format)
func (t SignerInfo) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SignerInfo from hex string (Canton MCMS format)
func (t *SignerInfo) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TimelockCall is a Record type
type TimelockCall struct {
	TargetInstanceAddress types.TEXT `json:"targetInstanceAddress"`
	FunctionName          types.TEXT `json:"functionName"`
	OperationData         types.TEXT `json:"operationData" hex:"bytes16"`
}

// ToMap converts TimelockCall to a map for DAML arguments
func (t TimelockCall) ToMap() map[string]any {
	m := make(map[string]any)

	m["targetInstanceAddress"] = string(t.TargetInstanceAddress)

	m["functionName"] = string(t.FunctionName)

	m["operationData"] = string(t.OperationData)

	return m
}

func (t TimelockCall) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TimelockCall) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TimelockCall to hex string (Canton MCMS format)
func (t TimelockCall) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TimelockCall from hex string (Canton MCMS format)
func (t *TimelockCall) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// IMCMSReceiverInterfaceID returns the interface ID for the IMCMSReceiver interface using the package name
func IMCMSReceiverInterfaceID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "MCMS.MCMSReceiver", "MCMSReceiver")
}

// IMCMSReceiverInterfaceIDWithPackageID returns the interface ID using the provided package ID instead of package name
func IMCMSReceiverInterfaceIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "MCMS.MCMSReceiver", "MCMSReceiver")
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error)
	CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error)
	ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error)
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

// BypasserExecuteBatch encodes parameters for the BypasserExecuteBatch choice.
func (e *encoder) BypasserExecuteBatch(args BypasserExecuteBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("BypasserExecuteBatch", args)
}

// CancelBatch encodes parameters for the CancelBatch choice.
func (e *encoder) CancelBatch(args CancelBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CancelBatch", args)
}

// ScheduleBatch encodes parameters for the ScheduleBatch choice.
func (e *encoder) ScheduleBatch(args ScheduleBatchParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ScheduleBatch", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
