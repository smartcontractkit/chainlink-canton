package ratelimiter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	api "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
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
	PackageName = "ccip-registry-rate-limiter-v2"
	PackageID   = "6856206c569bf6c13704eb5cd3fedecb64245fce1af80898b4ddf6580f51fa92"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	RateLimiterContextKey = types.TEXT("rate-limiter")
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

// ConsumeCapacity is a Record type
type ConsumeCapacity struct {
	Requested types.NUMERIC `json:"requested"`
}

// ToMap converts ConsumeCapacity to a map for DAML arguments
func (t ConsumeCapacity) ToMap() map[string]any {
	m := make(map[string]any)

	m["requested"] = t.Requested

	return m
}

func (t ConsumeCapacity) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeCapacity) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeCapacity to hex string (Canton MCMS format)
func (t ConsumeCapacity) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeCapacity from hex string (Canton MCMS format)
func (t *ConsumeCapacity) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ConsumeCapacityResult is a Record type
type ConsumeCapacityResult struct {
	RateLimiterCid         types.CONTRACT_ID `json:"rateLimiterCid"`
	AvailableBeforeConsume types.NUMERIC     `json:"availableBeforeConsume"`
	Consumed               types.NUMERIC     `json:"consumed"`
}

// ToMap converts ConsumeCapacityResult to a map for DAML arguments
func (t ConsumeCapacityResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["rateLimiterCid"] = model.NestedToDAMLValue(t.RateLimiterCid)

	m["availableBeforeConsume"] = t.AvailableBeforeConsume

	m["consumed"] = t.Consumed

	return m
}

func (t ConsumeCapacityResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ConsumeCapacityResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ConsumeCapacityResult to hex string (Canton MCMS format)
func (t ConsumeCapacityResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ConsumeCapacityResult from hex string (Canton MCMS format)
func (t *ConsumeCapacityResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RateLimitDirection is an enum type
type RateLimitDirection string

const (
	RateLimitDirectionRateLimitDirection_Outbound RateLimitDirection = "RateLimitDirection_Outbound"

	RateLimitDirectionRateLimitDirection_Inbound RateLimitDirection = "RateLimitDirection_Inbound"
)

func (e RateLimitDirection) GetEnumConstructor() string { return string(e) }

func (e RateLimitDirection) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimitDirection")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e RateLimitDirection) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimitDirection")
}

func (e RateLimitDirection) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *RateLimitDirection) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes RateLimitDirection to hex string (Canton MCMS format)
func (e RateLimitDirection) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes RateLimitDirection from hex string (Canton MCMS format)
func (e *RateLimitDirection) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = RateLimitDirection("")

// RateLimitMode is an enum type
type RateLimitMode string

const (
	RateLimitModeRateLimitMode_DefaultFinality RateLimitMode = "RateLimitMode_DefaultFinality"

	RateLimitModeRateLimitMode_CustomFinality RateLimitMode = "RateLimitMode_CustomFinality"
)

func (e RateLimitMode) GetEnumConstructor() string { return string(e) }

func (e RateLimitMode) GetEnumTypeID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimitMode")
}

// GetEnumTypeIDWithPackageID returns the enum type ID using the provided package ID instead of package name
func (e RateLimitMode) GetEnumTypeIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimitMode")
}

func (e RateLimitMode) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(e)
}

func (e *RateLimitMode) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, e)
}

// MarshalHex encodes RateLimitMode to hex string (Canton MCMS format)
func (e RateLimitMode) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(e)
}

// UnmarshalHex decodes RateLimitMode from hex string (Canton MCMS format)
func (e *RateLimitMode) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, e)
}

var _ types.ENUM = RateLimitMode("")

// RateLimiter is a Template type
type RateLimiter struct {
	InstanceId          types.TEXT         `json:"instanceId"`
	PoolInstanceId      types.TEXT         `json:"poolInstanceId"`
	PoolOwner           types.PARTY        `json:"poolOwner"`
	RemoteChainSelector types.NUMERIC      `json:"remoteChainSelector"`
	Direction           RateLimitDirection `json:"direction"`
	Mode                RateLimitMode      `json:"mode"`
	IsEnabled           types.BOOL         `json:"isEnabled"`
	Capacity            types.NUMERIC      `json:"capacity"`
	Rate                types.NUMERIC      `json:"rate"`
	Tokens              types.NUMERIC      `json:"tokens"`
	LastUpdated         types.TIMESTAMP    `json:"lastUpdated"`
	Observers           []types.PARTY      `json:"observers"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t RateLimiter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimiter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t RateLimiter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimiter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t RateLimiter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolInstanceId"] = string(t.PoolInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	if t.RemoteChainSelector != "" {
		args["remoteChainSelector"] = t.RemoteChainSelector
	}

	if t.Direction != "" {
		args["direction"] = model.NestedToDAMLValue(t.Direction)
	}

	if t.Mode != "" {
		args["mode"] = model.NestedToDAMLValue(t.Mode)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isEnabled"] = bool(t.IsEnabled)

	if t.Capacity != "" {
		args["capacity"] = t.Capacity
	}

	if t.Rate != "" {
		args["rate"] = t.Rate
	}

	if t.Tokens != "" {
		args["tokens"] = t.Tokens
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastUpdated"] = t.LastUpdated

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
func (t RateLimiter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolInstanceId"] = string(t.PoolInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["poolOwner"] = t.PoolOwner.ToMap()

	if t.RemoteChainSelector != "" {
		args["remoteChainSelector"] = t.RemoteChainSelector
	}

	if t.Direction != "" {
		args["direction"] = model.NestedToDAMLValue(t.Direction)
	}

	if t.Mode != "" {
		args["mode"] = model.NestedToDAMLValue(t.Mode)
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["isEnabled"] = bool(t.IsEnabled)

	if t.Capacity != "" {
		args["capacity"] = t.Capacity
	}

	if t.Rate != "" {
		args["rate"] = t.Rate
	}

	if t.Tokens != "" {
		args["tokens"] = t.Tokens
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["lastUpdated"] = t.LastUpdated

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

func (t RateLimiter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RateLimiter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RateLimiter to hex string (Canton MCMS format)
func (t RateLimiter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RateLimiter from hex string (Canton MCMS format)
func (t *RateLimiter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for RateLimiter

// ConsumeCapacity exercises the ConsumeCapacity choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) ConsumeCapacity(contractID string, args ConsumeCapacity) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "ConsumeCapacity",
		Arguments:  argsToMap(args),
	}
}

// ConsumeCapacityWithPackageID exercises the ConsumeCapacity choice using the provided package ID instead of package name
func (t RateLimiter) ConsumeCapacityWithPackageID(contractID string, packageID string, args ConsumeCapacity) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "ConsumeCapacity",
		Arguments:  argsToMap(args),
	}
}

// SetConfig exercises the SetConfig choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) SetConfig(contractID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// SetConfigWithPackageID exercises the SetConfig choice using the provided package ID instead of package name
func (t RateLimiter) SetConfigWithPackageID(contractID string, packageID string, args SetConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetConfig",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this RateLimiter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RateLimiter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t RateLimiter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetObservers exercises the SetObservers choice on this RateLimiter contract
// This method uses the package name in the template ID
func (t RateLimiter) SetObservers(contractID string, args SetObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetObservers",
		Arguments:  argsToMap(args),
	}
}

// SetObserversWithPackageID exercises the SetObservers choice using the provided package ID instead of package name
func (t RateLimiter) SetObserversWithPackageID(contractID string, packageID string, args SetObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "RateLimiter"),
		ContractID: contractID,
		Choice:     "SetObservers",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this RateLimiter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t RateLimiter) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Registry.RateLimiterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t RateLimiter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Registry.RateLimiterV2", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for RateLimiter

var _ api.IMCMSReceiver = (*RateLimiter)(nil)

// SetConfig is a Record type
type SetConfig struct {
	NewIsEnabled types.BOOL    `json:"newIsEnabled"`
	NewCapacity  types.NUMERIC `json:"newCapacity"`
	NewRate      types.NUMERIC `json:"newRate"`
}

// ToMap converts SetConfig to a map for DAML arguments
func (t SetConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["newIsEnabled"] = bool(t.NewIsEnabled)

	m["newCapacity"] = t.NewCapacity

	m["newRate"] = t.NewRate

	return m
}

func (t SetConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetConfig to hex string (Canton MCMS format)
func (t SetConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetConfig from hex string (Canton MCMS format)
func (t *SetConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error)
	SetConfig(args SetConfig) (*bind.EncodedChoice, error)
	SetObservers(args SetObservers) (*bind.EncodedChoice, error)
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

// ConsumeCapacity encodes parameters for the ConsumeCapacity choice.
func (e *encoder) ConsumeCapacity(args ConsumeCapacity) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ConsumeCapacity", args)
}

// SetConfig encodes parameters for the SetConfig choice.
func (e *encoder) SetConfig(args SetConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetConfig", args)
}

// SetObservers encodes parameters for the SetObservers choice.
func (e *encoder) SetObservers(args SetObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetObservers", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
