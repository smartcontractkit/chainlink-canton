package ccipreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
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
	PackageID   = "74414d9387e515b14a3dbb964f1c951fa81d468d8b8c8489e9a58d71e46071ce"
	SDKVersion  = "3.4.10"
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
	Owner              types.PARTY                     `json:"owner"`
	Router             types.CONTRACT_ID               `json:"router"`
	MessageId          types.TEXT                      `json:"messageId"`
	Message            common.MessageV1                `json:"message"`
	TokenReleaseResult *interfaces.ReleaseOrMintResult `json:"tokenReleaseResult" hex:"optional"`
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
	args["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenReleaseResult,
		}
	} else {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
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
	args["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenReleaseResult,
		}
	} else {
		args["tokenReleaseResult"] = map[string]any{
			"_type": "optional",
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
	InstanceId            types.TEXT                  `json:"instanceId"`
	Owner                 types.PARTY                 `json:"owner"`
	RequiredCCVs          []common.RawInstanceAddress `json:"requiredCCVs"`
	OptionalCCVs          []common.RawInstanceAddress `json:"optionalCCVs"`
	OptionalThreshold     types.INT64                 `json:"optionalThreshold"`
	MinBlockConfirmations types.INT64                 `json:"minBlockConfirmations"`
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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minBlockConfirmations"] = int64(t.MinBlockConfirmations)

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
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["minBlockConfirmations"] = int64(t.MinBlockConfirmations)

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
func (t CCIPReceiver) Execute(contractID string, args Execute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t CCIPReceiver) ExecuteWithPackageID(contractID string, packageID string, args Execute2) *model.ExerciseCommand {
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
	CcvCid          types.CONTRACT_ID  `json:"ccvCid"`
	VerifierResults types.TEXT         `json:"verifierResults"`
	CcvExtraContext common.CCIPContext `json:"ccvExtraContext"`
}

// ToMap converts CCVInput to a map for DAML arguments
func (t CCVInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccvCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvCid).(mapper); ok {
			return m.toMap()
		}
		return t.CcvCid
	}()

	m["verifierResults"] = string(t.VerifierResults)

	m["ccvExtraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.CcvExtraContext
	}()

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

// Execute2 is a Record type
type Execute2 struct {
	Context        common.CCIPContext  `json:"context"`
	RouterCid      types.CONTRACT_ID   `json:"routerCid"`
	EncodedMessage types.TEXT          `json:"encodedMessage"`
	TokenTransfer  *TokenTransferInput `json:"tokenTransfer" hex:"optional"`
	CcvInputs      []CCVInput          `json:"ccvInputs"`
}

// ToMap converts Execute2 to a map for DAML arguments
func (t Execute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["routerCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenTransfer,
		}
	} else {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
		}
	}

	m["ccvInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvInputs))
		for _, e := range t.CcvInputs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	return m
}

func (t Execute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Execute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Execute2 to hex string (Canton MCMS format)
func (t Execute2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Execute2 from hex string (Canton MCMS format)
func (t *Execute2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVs is a Record type
type GetRequiredCCVs struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

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

// GetRequiredCCVsMCMSParams is GetRequiredCCVs without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetRequiredCCVsMCMSParams struct {
}

// MarshalHex encodes GetRequiredCCVsMCMSParams to hex string for MCMS operationData.
func (t GetRequiredCCVsMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsMCMSParams from hex string.
func (t *GetRequiredCCVsMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	TokenPoolCid       types.CONTRACT_ID     `json:"tokenPoolCid"`
	TokenReceiverParty types.PARTY           `json:"tokenReceiverParty"`
	TokenInput         interfaces.TokenInput `json:"tokenInput"`
	PoolExtraContext   common.CCIPContext    `json:"poolExtraContext"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]any {
	m := make(map[string]any)

	m["tokenPoolCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenPoolCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenPoolCid
	}()

	m["tokenReceiverParty"] = t.TokenReceiverParty.ToMap()

	m["tokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	m["poolExtraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PoolExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolExtraContext
	}()

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
	NewRequiredCCVs []common.RawInstanceAddress `json:"newRequiredCCVs"`
}

// ToMap converts UpdateRequiredCCVs to a map for DAML arguments
func (t UpdateRequiredCCVs) ToMap() map[string]any {
	m := make(map[string]any)

	m["newRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.NewRequiredCCVs))
		for _, e := range t.NewRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
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
	Execute2(args Execute2) (*bind.EncodedChoice, error)
	GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error)
	GetRequiredCCVsMCMSParams(args GetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error)
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

// Execute2 encodes parameters for the Execute2 choice.
func (e *encoder) Execute2(args Execute2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Execute2", args)
}

// GetRequiredCCVs encodes parameters for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVs(args GetRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// GetRequiredCCVsMCMSParams encodes MCMS parameters (without Caller) for the GetRequiredCCVs choice.
func (e *encoder) GetRequiredCCVsMCMSParams(args GetRequiredCCVsMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVs", args)
}

// UpdateRequiredCCVs encodes parameters for the UpdateRequiredCCVs choice.
func (e *encoder) UpdateRequiredCCVs(args UpdateRequiredCCVs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("UpdateRequiredCCVs", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
