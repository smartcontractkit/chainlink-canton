package ccvs

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

const PackageID = "60ab65f876d9fcb12cd327fff8aafd91c8c42d2bc4cfe0722025267484666724"
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

// CommitteeVerifier is a Template type
type CommitteeVerifier struct {
	Owner               PARTY  `json:"owner"`
	InstanceId          TEXT   `json:"instanceId"`
	CcipOwner           PARTY  `json:"ccipOwner"`
	VersionTag          TEXT   `json:"versionTag"`
	MessageSentObserver PARTY  `json:"messageSentObserver"`
	StorageLocation     TEXT   `json:"storageLocation"`
	Threshold           INT64  `json:"threshold"`
	Signers             []TEXT `json:"signers"`
}

// GetTemplateID returns the template ID for this template
func (t CommitteeVerifier) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "CommitteeVerifier")
}

// CreateCommand returns a CreateCommand for this template
func (t CommitteeVerifier) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["owner"] = t.Owner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["versionTag"] = string(t.VersionTag)

	args["messageSentObserver"] = t.MessageSentObserver.ToMap()

	args["storageLocation"] = string(t.StorageLocation)

	args["threshold"] = int64(t.Threshold)

	if len(t.Signers) > 0 {
		args["signers"] = func() []interface{} {
			res := make([]interface{}, 0, len(t.Signers))
			for _, e := range t.Signers {
				res = append(res, string(e))
			}
			return res
		}()
	}

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for CommitteeVerifier using JsonCodec
func (t CommitteeVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CommitteeVerifier using JsonCodec
func (t *CommitteeVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CommitteeVerifier

// CommitteeVerifierVerifyMessage exercises the CommitteeVerifier_VerifyMessage choice on this CommitteeVerifier contract
func (t CommitteeVerifier) CommitteeVerifierVerifyMessage(contractID string, args CommitteeVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CommitteeVerifier contract
func (t CommitteeVerifier) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CommitteeVerifierForwardToVerifier exercises the CommitteeVerifier_ForwardToVerifier choice on this CommitteeVerifier contract
func (t CommitteeVerifier) CommitteeVerifierForwardToVerifier(contractID string, args CommitteeVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessage exercises the CrossChainVerifier_VerifyMessage choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
func (t CommitteeVerifier) CrossChainVerifierVerifyMessage(contractID string, args CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifier exercises the CrossChainVerifier_ForwardToVerifier choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifier(contractID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CommitteeVerifier

var _ IICrossChainVerifier = (*CommitteeVerifier)(nil)

// CommitteeVerifierForwardToVerifier is a Record type
type CommitteeVerifierForwardToVerifier struct {
	CcvRegistryCid CONTRACT_ID  `json:"ccvRegistryCid"`
	Message        MessageV1    `json:"message"`
	MessageId      TEXT         `json:"messageId"`
	FeeToken       InstrumentId `json:"feeToken"`
	FeeTokenAmount NUMERIC      `json:"feeTokenAmount"`
	VerifierArgs   TEXT         `json:"verifierArgs"`
	Caller         PARTY        `json:"caller"`
}

// toMap converts CommitteeVerifierForwardToVerifier to a map for DAML arguments
func (t CommitteeVerifierForwardToVerifier) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.CcvRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.CcvRegistryCid
		}(),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"messageId": string(t.MessageId),
		"feeToken": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.FeeToken).(mapper); ok {
				return m.toMap()
			}
			return t.FeeToken
		}(),
		"feeTokenAmount": (*big.Int)(t.FeeTokenAmount),
		"verifierArgs":   string(t.VerifierArgs),
		"caller":         t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CommitteeVerifierForwardToVerifier using JsonCodec
func (t CommitteeVerifierForwardToVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CommitteeVerifierForwardToVerifier using JsonCodec
func (t *CommitteeVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CommitteeVerifierVerifyMessage is a Record type
type CommitteeVerifierVerifyMessage struct {
	CcvRegistryCid CONTRACT_ID `json:"ccvRegistryCid"`
	Message        MessageV1   `json:"message"`
	MessageId      TEXT        `json:"messageId"`
	CcvData        TEXT        `json:"ccvData"`
	Receiver       PARTY       `json:"receiver"`
	Caller         PARTY       `json:"caller"`
}

// toMap converts CommitteeVerifierVerifyMessage to a map for DAML arguments
func (t CommitteeVerifierVerifyMessage) toMap() map[string]interface{} {
	return map[string]interface{}{

		"ccvRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.CcvRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.CcvRegistryCid
		}(),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"messageId": string(t.MessageId),
		"ccvData":   string(t.CcvData),
		"receiver":  t.Receiver.ToMap(),
		"caller":    t.Caller.ToMap(),
	}
}

// MarshalJSON implements custom JSON marshaling for CommitteeVerifierVerifyMessage using JsonCodec
func (t CommitteeVerifierVerifyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CommitteeVerifierVerifyMessage using JsonCodec
func (t *CommitteeVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
