package ccipreceiver

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

const PackageName = "ccip-receiver"
const SDKVersion = "3.4.10"

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
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

// CCIPMessageReceived is a Template type
type CCIPMessageReceived struct {
	Owner              PARTY                `json:"owner"`
	Router             CONTRACT_ID          `json:"router"`
	MessageId          TEXT                 `json:"messageId"`
	Message            MessageV1            `json:"message"`
	TokenReleaseResult *ReleaseOrMintResult `json:"tokenReleaseResult"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["router"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenReleaseResult,
		}
	} else {
		args["tokenReleaseResult"] = map[string]interface{}{
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["router"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageId"] = string(t.MessageId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	if t.TokenReleaseResult != nil {
		args["tokenReleaseResult"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenReleaseResult,
		}
	} else {
		args["tokenReleaseResult"] = map[string]interface{}{
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
	return jsonCodec.Marshall(t)
}

func (t *CCIPMessageReceived) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCIPMessageReceived

// Archive exercises the Archive choice on this CCIPMessageReceived contract
// This method uses the package name in the template ID
func (t CCIPMessageReceived) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPMessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPMessageReceived) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPMessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCIPReceiver is a Template type
type CCIPReceiver struct {
	InstanceId   TEXT                 `json:"instanceId"`
	Owner        PARTY                `json:"owner"`
	RequiredCCVs []RawInstanceAddress `json:"requiredCCVs"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *CCIPReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCIPReceiver

// Archive exercises the Archive choice on this CCIPReceiver contract
// This method uses the package name in the template ID
func (t CCIPReceiver) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPReceiver) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
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
	CcvCid          CONTRACT_ID `json:"ccvCid"`
	VerifierResults TEXT        `json:"verifierResults"`
}

// ToMap converts CCVInput to a map for DAML arguments
func (t CCVInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["ccvCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *CCVInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Execute2 is a Record type
type Execute2 struct {
	RouterCid              CONTRACT_ID          `json:"routerCid"`
	OffRampCid             CONTRACT_ID          `json:"offRampCid"`
	GlobalConfigCid        CONTRACT_ID          `json:"globalConfigCid"`
	TokenAdminRegistryCid  CONTRACT_ID          `json:"tokenAdminRegistryCid"`
	RmnRemoteCid           CONTRACT_ID          `json:"rmnRemoteCid"`
	EncodedMessage         TEXT                 `json:"encodedMessage"`
	TokenTransfer          *TokenTransferInput  `json:"tokenTransfer"`
	CcvInputs              []CCVInput           `json:"ccvInputs"`
	AdditionalRequiredCCVs []RawInstanceAddress `json:"additionalRequiredCCVs"`
}

// ToMap converts Execute2 to a map for DAML arguments
func (t Execute2) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RouterCid).(mapper); ok {
			return m.toMap()
		}
		return t.RouterCid
	}()

	m["offRampCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.OffRampCid).(mapper); ok {
			return m.toMap()
		}
		return t.OffRampCid
	}()

	m["globalConfigCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
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

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenTransfer,
		}
	} else {
		m["tokenTransfer"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["ccvInputs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvInputs))
		for _, e := range t.CcvInputs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["additionalRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.AdditionalRequiredCCVs))
		for _, e := range t.AdditionalRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *Execute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVs is a Record type
type GetRequiredCCVs struct {
	Caller PARTY `json:"caller"`
}

// ToMap converts GetRequiredCCVs to a map for DAML arguments
func (t GetRequiredCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetRequiredCCVs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// TokenTransferInput is a Record type
type TokenTransferInput struct {
	TokenPoolCid       CONTRACT_ID `json:"tokenPoolCid"`
	TokenReceiverParty PARTY       `json:"tokenReceiverParty"`
	TokenInput         TokenInput  `json:"tokenInput"`
}

// ToMap converts TokenTransferInput to a map for DAML arguments
func (t TokenTransferInput) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["tokenPoolCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenPoolCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenPoolCid
	}()

	m["tokenReceiverParty"] = t.TokenReceiverParty.ToMap()

	m["tokenInput"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.TokenInput
	}()

	return m
}

func (t TokenTransferInput) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *TokenTransferInput) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// UpdateRequiredCCVs is a Record type
type UpdateRequiredCCVs struct {
	NewRequiredCCVs []RawInstanceAddress `json:"newRequiredCCVs"`
}

// ToMap converts UpdateRequiredCCVs to a map for DAML arguments
func (t UpdateRequiredCCVs) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["newRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.NewRequiredCCVs))
		for _, e := range t.NewRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *UpdateRequiredCCVs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
