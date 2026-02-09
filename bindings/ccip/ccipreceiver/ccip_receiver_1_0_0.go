package ccipreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"https://github.com/smartcontractkit/go-daml/pkg/codec"
	"https://github.com/smartcontractkit/go-daml/pkg/model"
	. "https://github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
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

// CCIPReceiver is a Template type
type CCIPReceiver struct {
	Owner        PARTY         `json:"owner"`
	CcipOwner    PARTY         `json:"ccipOwner"`
	RequiredCCVs []CONTRACT_ID `json:"requiredCCVs"`
	OptionalCCVs []CONTRACT_ID `json:"optionalCCVs"`
	Threshold    INT64         `json:"threshold"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPReceiver) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "CCIPReceiver")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPReceiver) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "CCIPReceiver")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPReceiver) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["threshold"] = int64(t.Threshold)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPReceiver) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["threshold"] = int64(t.Threshold)

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

// Any2CantonMessageReceiverGetCCVs exercises the Any2CantonMessageReceiver_GetCCVs choice on this CCIPReceiver contract via the IIAny2CantonMessageReceiver interface
// This method uses the package name in the template ID
func (t CCIPReceiver) Any2CantonMessageReceiverGetCCVs(contractID string, args Any2CantonMessageReceiverGetCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
		ContractID: contractID,
		Choice:     "Any2CantonMessageReceiver_GetCCVs",
		Arguments:  argsToMap(args),
	}
}

// Any2CantonMessageReceiverGetCCVsWithPackageID exercises the Any2CantonMessageReceiver_GetCCVs choice using the provided package ID instead of package name
func (t CCIPReceiver) Any2CantonMessageReceiverGetCCVsWithPackageID(contractID string, packageID string, args Any2CantonMessageReceiverGetCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
		ContractID: contractID,
		Choice:     "Any2CantonMessageReceiver_GetCCVs",
		Arguments:  argsToMap(args),
	}
}

// Any2CantonMessageReceiverCCIPReceive exercises the Any2CantonMessageReceiver_CCIPReceive choice on this CCIPReceiver contract via the IIAny2CantonMessageReceiver interface
// This method uses the package name in the template ID
func (t CCIPReceiver) Any2CantonMessageReceiverCCIPReceive(contractID string, args Any2CantonMessageReceiverCCIPReceive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
		ContractID: contractID,
		Choice:     "Any2CantonMessageReceiver_CCIPReceive",
		Arguments:  argsToMap(args),
	}
}

// Any2CantonMessageReceiverCCIPReceiveWithPackageID exercises the Any2CantonMessageReceiver_CCIPReceive choice using the provided package ID instead of package name
func (t CCIPReceiver) Any2CantonMessageReceiverCCIPReceiveWithPackageID(contractID string, packageID string, args Any2CantonMessageReceiverCCIPReceive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
		ContractID: contractID,
		Choice:     "Any2CantonMessageReceiver_CCIPReceive",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CCIPReceiver

var _ IIAny2CantonMessageReceiver = (*CCIPReceiver)(nil)

// MessageReceived is a Template type
type MessageReceived struct {
	Owner   PARTY             `json:"owner"`
	Message Any2CantonMessage `json:"message"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t MessageReceived) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "MessageReceived")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t MessageReceived) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "MessageReceived")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t MessageReceived) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t MessageReceived) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t MessageReceived) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *MessageReceived) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MessageReceived

// Archive exercises the Archive choice on this MessageReceived contract
// This method uses the package name in the template ID
func (t MessageReceived) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CCIPReceiver", "MessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t MessageReceived) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CCIPReceiver", "MessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}
