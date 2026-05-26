package receiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	core "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/core"
	extensionapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/chainlink/chainlinkapi"
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
	PackageName = "ccip-receiver"
	PackageID   = "2db8ed1c8b9fd7585b30b1730eeb03b130495a255f8178c585ada95f7569bef3"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
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

// CCIPMessageReceived is a Template type
type CCIPMessageReceived struct {
	Owner              types.PARTY                       `json:"owner"`
	Router             types.CONTRACT_ID                 `json:"router"`
	MessageId          types.TEXT                        `json:"messageId"`
	Message            core.MessageV1                    `json:"message"`
	TokenReleaseResult *extensionapi.ReleaseOrMintResult `json:"tokenReleaseResult" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPMessageReceived) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPMessageReceived")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPMessageReceived) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPMessageReceived")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPMessageReceived) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["router"] = model.NestedToDAMLValue(t.Router)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = model.NestedToDAMLValue(t.Message)

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenReleaseResult),
		}
	} else {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPMessageReceived) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["router"] = model.NestedToDAMLValue(t.Router)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = model.NestedToDAMLValue(t.Message)

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenReleaseResult),
		}
	} else {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPMessageReceived) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPMessageReceived) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPMessageReceived to hex string (Canton MCMS format)
func (t CCIPMessageReceived) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPMessageReceived from hex string (Canton MCMS format)
func (t *CCIPMessageReceived) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPMessageReceived

// Archive exercises the Archive choice on this CCIPMessageReceived contract
// This method uses the package name in the template ID
func (t CCIPMessageReceived) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPMessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPMessageReceived) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPMessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CCIPReceiver is a Template type
type CCIPReceiver struct {
	InstanceId             types.TEXT                        `json:"instanceId"`
	Owner                  types.PARTY                       `json:"owner"`
	RequiredCCVs           []chainlinkapi.RawInstanceAddress `json:"requiredCCVs"`
	OptionalCCVs           []chainlinkapi.RawInstanceAddress `json:"optionalCCVs"`
	OptionalThreshold      types.INT64                       `json:"optionalThreshold"`
	ReceiverFinalityConfig core.FinalityConfig               `json:"receiverFinalityConfig"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPReceiver) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPReceiver) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPReceiver) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPReceiver) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPReceiver) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPReceiver to hex string (Canton MCMS format)
func (t CCIPReceiver) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPReceiver from hex string (Canton MCMS format)
func (t *CCIPReceiver) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPReceiver

// Archive exercises the Archive choice on this CCIPReceiver contract
// This method uses the package name in the template ID
func (t CCIPReceiver) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPReceiver) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetRequiredCCVs exercises the GetRequiredCCVs choice on this CCIPReceiver contract
// This method uses the package name in the template ID
func (t CCIPReceiver) GetRequiredCCVs(contractID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsWithPackageID exercises the GetRequiredCCVs choice using the provided package ID instead of package name
func (t CCIPReceiver) GetRequiredCCVsWithPackageID(contractID string, packageID string, args GetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// Execute exercises the Execute choice on this CCIPReceiver contract
// This method uses the package name in the template ID
func (t CCIPReceiver) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t CCIPReceiver) ExecuteWithPackageID(contractID string, packageID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// UpdateRequiredCCVs exercises the UpdateRequiredCCVs choice on this CCIPReceiver contract
// This method uses the package name in the template ID
func (t CCIPReceiver) UpdateRequiredCCVs(contractID string, args UpdateRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "UpdateRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// UpdateRequiredCCVsWithPackageID exercises the UpdateRequiredCCVs choice using the provided package ID instead of package name
func (t CCIPReceiver) UpdateRequiredCCVsWithPackageID(contractID string, packageID string, args UpdateRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "UpdateRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// CCVInput is a Record type
type CCVInput struct {
	CcvCid          types.CONTRACT_ID                          `json:"ccvCid"`
	VerifierResults types.TEXT                                 `json:"verifierResults"`
	CcvExtraContext splice_api_token_metadata_v1.ChoiceContext `json:"ccvExtraContext"`
}

// ToMap converts CCVInput to a map for DAML arguments
func (t CCVInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvCid"] = model.NestedToDAMLValue(t.CcvCid)

	m["verifierResults"] = string(t.VerifierResults)

	m["ccvExtraContext"] = model.NestedToDAMLValue(t.CcvExtraContext)

	return m
}

func (t CCVInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCVInput to hex string (Canton MCMS format)
func (t CCVInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVInput from hex string (Canton MCMS format)
func (t *CCVInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Execute is a Record type
type Execute struct {
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenTransfer  *TokenTransferInput                        `json:"tokenTransfer" hex:"optional"`
	CcvInputs      []CCVInput                                 `json:"ccvInputs"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransfer),
		}
	} else {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["ccvInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvInputs))
		for _, e := range t.CcvInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t Execute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Execute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Execute to hex string (Canton MCMS format)
func (t Execute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Execute from hex string (Canton MCMS format)
func (t *Execute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVs is a Record type
type GetRequiredCCVs struct {
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenPoolCid   *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenPoolCid != nil {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolCid),
		}
	} else {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t GetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVs to hex string (Canton MCMS format)
func (t GetRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVs from hex string (Canton MCMS format)
func (t *GetRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	TokenPoolCid       types.CONTRACT_ID                          `json:"tokenPoolCid"`
	TokenReceiverParty types.PARTY                                `json:"tokenReceiverParty"`
	PoolExtraContext   splice_api_token_metadata_v1.ChoiceContext `json:"poolExtraContext"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenPoolCid"] = model.NestedToDAMLValue(t.TokenPoolCid)

	m["tokenReceiverParty"] = t.TokenReceiverParty.ToMap()

	m["poolExtraContext"] = model.NestedToDAMLValue(t.PoolExtraContext)

	return m
}

func (t TokenTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *TokenTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes TokenTransferInput to hex string (Canton MCMS format)
func (t TokenTransferInput) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes TokenTransferInput from hex string (Canton MCMS format)
func (t *TokenTransferInput) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// UpdateRequiredCCVs is a Record type
type UpdateRequiredCCVs struct {
	NewRequiredCCVs []chainlinkapi.RawInstanceAddress `json:"newRequiredCCVs"`
}

// ToMap converts UpdateRequiredCCVs to a map for DAML arguments
func (t UpdateRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["newRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.NewRequiredCCVs))
		for _, e := range t.NewRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t UpdateRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *UpdateRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes UpdateRequiredCCVs to hex string (Canton MCMS format)
func (t UpdateRequiredCCVs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes UpdateRequiredCCVs from hex string (Canton MCMS format)
func (t *UpdateRequiredCCVs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	Execute(args Execute) (*bind.EncodedChoice, error)
	GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error)
	UpdateRequiredCCVs(args UpdateRequiredCCVs) (*bind.EncodedChoice, error)
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

// Execute encodes parameters for the Execute choice.
func (e *encoder) Execute(args Execute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Execute", args)
}

// GetRequiredCCVs encodes parameters for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// UpdateRequiredCCVs encodes parameters for the UpdateRequiredCCVs choice.
func (e *encoder) UpdateRequiredCCVs(args UpdateRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateRequiredCCVs", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
