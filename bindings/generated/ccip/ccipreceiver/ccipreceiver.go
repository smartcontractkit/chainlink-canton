package ccipreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-canton/bindings/codec"
	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	model "github.com/smartcontractkit/chainlink-canton/bindings/ledger"
	"github.com/smartcontractkit/chainlink-canton/bindings/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
)

const (
	PackageName = "ccip-receiver"
	PackageID   = "897f6bab414efa1d2e4659f36131af9565882876040b8356ee47599c2ad92760"
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
	TokenReleaseResult *interfaces.ReleaseOrMintResult `json:"tokenReleaseResult"`
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
	InstanceId   types.TEXT                  `json:"instanceId"`
	Owner        types.PARTY                 `json:"owner"`
	RequiredCCVs []common.RawInstanceAddress `json:"requiredCCVs"`
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
	CcvCid          types.CONTRACT_ID `json:"ccvCid"`
	VerifierResults types.TEXT        `json:"verifierResults"`
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

// Execute2 is a Record type
type Execute2 struct {
	RouterCid              types.CONTRACT_ID           `json:"routerCid"`
	OffRampCid             types.CONTRACT_ID           `json:"offRampCid"`
	GlobalConfigCid        types.CONTRACT_ID           `json:"globalConfigCid"`
	TokenAdminRegistryCid  types.CONTRACT_ID           `json:"tokenAdminRegistryCid"`
	RmnRemoteCid           types.CONTRACT_ID           `json:"rmnRemoteCid"`
	EncodedMessage         types.TEXT                  `json:"encodedMessage"`
	TokenTransfer          *TokenTransferInput         `json:"tokenTransfer"`
	CcvInputs              []CCVInput                  `json:"ccvInputs"`
	AdditionalRequiredCCVs []common.RawInstanceAddress `json:"additionalRequiredCCVs"`
}

// ToMap converts Execute2 to a map for DAML arguments
func (t Execute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	m["offRampCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampCid
	}()

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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

	m["additionalRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.AdditionalRequiredCCVs))
		for _, e := range t.AdditionalRequiredCCVs {
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

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	TokenPoolCid       types.CONTRACT_ID     `json:"tokenPoolCid"`
	TokenReceiverParty types.PARTY           `json:"tokenReceiverParty"`
	TokenInput         interfaces.TokenInput `json:"tokenInput"`
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
