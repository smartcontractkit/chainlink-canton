package perpartyrouter

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

// CCIPSendFromRouter is a Record type
type CCIPSendFromRouter struct {
	RouterPartyOwner      PARTY         `json:"routerPartyOwner"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	DestChainSelector     NUMERIC       `json:"destChainSelector"`
	Receiver              TEXT          `json:"receiver"`
	Payload               TEXT          `json:"payload"`
	ExecutionGasLimit     INT64         `json:"executionGasLimit"`
	CcipReceiveGasLimit   INT64         `json:"ccipReceiveGasLimit"`
	CurrentSequenceNumber NUMERIC       `json:"currentSequenceNumber"`
	TokenSendTicket       *CONTRACT_ID  `json:"tokenSendTicket"`
	CcvTickets            []CONTRACT_ID `json:"ccvTickets"`
}

// ToMap converts CCIPSendFromRouter to a map for DAML arguments
func (t CCIPSendFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

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

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["executionGasLimit"] = int64(t.ExecutionGasLimit)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	m["currentSequenceNumber"] = (*big.Int)(t.CurrentSequenceNumber)

	if t.TokenSendTicket != nil {
		m["tokenSendTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenSendTicket,
		}
	} else {
		m["tokenSendTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["ccvTickets"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvTickets))
		for _, e := range t.CcvTickets {
			res = append(res, e)
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for CCIPSendFromRouter using JsonCodec
func (t CCIPSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPSendFromRouter using JsonCodec
func (t *CCIPSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CCIPSendFromRouterResult is a Record type
type CCIPSendFromRouterResult struct {
	MessageId            TEXT    `json:"messageId"`
	EncodedMessage       TEXT    `json:"encodedMessage"`
	NewSequenceNumber    NUMERIC `json:"newSequenceNumber"`
	DestChainSelector    NUMERIC `json:"destChainSelector"`
	VerifierBlobs        []TEXT  `json:"verifierBlobs"`
	MessageSentObservers []PARTY `json:"messageSentObservers"`
}

// ToMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["messageId"] = string(t.MessageId)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["newSequenceNumber"] = (*big.Int)(t.NewSequenceNumber)

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["verifierBlobs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["messageSentObservers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

// MarshalJSON implements custom JSON marshaling for CCIPSendFromRouterResult using JsonCodec
func (t CCIPSendFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for CCIPSendFromRouterResult using JsonCodec
func (t *CCIPSendFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	DestChainSelector     NUMERIC       `json:"destChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// ToMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

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

	m["destChainSelector"] = (*big.Int)(t.DestChainSelector)

	m["hasTokenTransfer"] = bool(t.HasTokenTransfer)

	if t.InstrumentId != nil {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.InstrumentId,
		}
	} else {
		m["instrumentId"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForSend using JsonCodec
func (t GetRequiredCCVsForSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForSend using JsonCodec
func (t *GetRequiredCCVsForSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// OnRamp is a Template type
type OnRamp struct {
	CcipOwner  PARTY `json:"ccipOwner"`
	InstanceId TEXT  `json:"instanceId"`
}

// GetTemplateID returns the template ID for this template
func (t OnRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OnRamp", "OnRamp")
}

// CreateCommand returns a CreateCommand for this template
func (t OnRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for OnRamp using JsonCodec
func (t OnRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for OnRamp using JsonCodec
func (t *OnRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for OnRamp

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this OnRamp contract
func (t OnRamp) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OnRamp contract
func (t OnRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CCIPSendFromRouter exercises the CCIPSendFromRouter choice on this OnRamp contract
func (t OnRamp) CCIPSendFromRouter(contractID string, args CCIPSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CCIPSendFromRouter",
		Arguments:  argsToMap(args),
	}
}
