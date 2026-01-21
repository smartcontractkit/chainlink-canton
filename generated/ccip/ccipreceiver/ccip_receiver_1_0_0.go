package ccipreceiver

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

const PackageID = "5d7c5efa2ac789287b64fa6286367015bf32b35bbdd7d43961834909cb497321"
const SDKVersion = "3.4.8"

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

	// Check if the type has a ToMap method
	type mapper interface {
		ToMap() map[string]interface{}
	}

	if mapper, ok := args.(mapper); ok {
		return mapper.ToMap()
	}

	return map[string]interface{}{
		"args": args,
	}
}

// CCIPReceiver is a Template type
type CCIPReceiver struct {
	Owner        PARTY         `json:"owner"`
	CcipOwner    PARTY         `json:"ccipOwner"`
	RequiredCCVs []CONTRACT_ID `json:"requiredCCVs"`
	OptionalCCVs []CONTRACT_ID `json:"optionalCCVs"`
	Threshold    INT64         `json:"threshold"`
}

// GetTemplateID returns the template ID for this template
func (t CCIPReceiver) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "CCIPReceiver")
}

// CreateCommand returns a CreateCommand for this template
func (t CCIPReceiver) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["ccipOwner"] = t.CcipOwner.ToMap()

	if len(t.RequiredCCVs) > 0 {
		args["requiredCCVs"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.RequiredCCVs))
			for _, e := range t.RequiredCCVs {
				res = append(res, e)
			}
			return res
		}()
	}

	if len(t.OptionalCCVs) > 0 {
		args["optionalCCVs"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.OptionalCCVs))
			for _, e := range t.OptionalCCVs {
				res = append(res, e)
			}
			return res
		}()
	}

	args["threshold"] = int64(t.Threshold)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CCIPReceiver using JsonCodec
func (t CCIPReceiver) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPReceiver using JsonCodec
func (t *CCIPReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CCIPReceiver

// Archive exercises the Archive choice on this CCIPReceiver contract
func (t CCIPReceiver) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "CCIPReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// Any2CantonMessageReceiverGetCCVs exercises the Any2CantonMessageReceiver_GetCCVs choice on this CCIPReceiver contract via the IIAny2CantonMessageReceiver interface
func (t CCIPReceiver) Any2CantonMessageReceiverGetCCVs(contractID string, args Any2CantonMessageReceiverGetCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
		ContractID: contractID,
		Choice:     "Any2CantonMessageReceiver_GetCCVs",
		Arguments:  argsToMap(args),
	}
}

// Any2CantonMessageReceiverCCIPReceive exercises the Any2CantonMessageReceiver_CCIPReceive choice on this CCIPReceiver contract via the IIAny2CantonMessageReceiver interface
func (t CCIPReceiver) Any2CantonMessageReceiverCCIPReceive(contractID string, args Any2CantonMessageReceiverCCIPReceive) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "IAny2CantonMessageReceiver"),
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

// GetTemplateID returns the template ID for this template
func (t MessageReceived) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "MessageReceived")
}

// CreateCommand returns a CreateCommand for this template
func (t MessageReceived) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

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

// MarshalJSON implements custom JSON marshaling for MessageReceived using JsonCodec
func (t MessageReceived) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for MessageReceived using JsonCodec
func (t *MessageReceived) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for MessageReceived

// Archive exercises the Archive choice on this MessageReceived contract
func (t MessageReceived) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CCIPReceiver", "MessageReceived"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}
