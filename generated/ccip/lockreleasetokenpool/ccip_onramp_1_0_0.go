package lockreleasetokenpool

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

// toMap converts CCIPSendFromRouter to a map for DAML arguments
func (t CCIPSendFromRouter) toMap() map[string]interface{} {
	return map[string]interface{}{

		"routerPartyOwner": t.RouterPartyOwner.ToMap(),
		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"destChainSelector":     (*big.Int)(t.DestChainSelector),
		"receiver":              string(t.Receiver),
		"payload":               string(t.Payload),
		"executionGasLimit":     int64(t.ExecutionGasLimit),
		"ccipReceiveGasLimit":   int64(t.CcipReceiveGasLimit),
		"currentSequenceNumber": (*big.Int)(t.CurrentSequenceNumber),
		"tokenSendTicket":       *t.TokenSendTicket,
		"ccvTickets": func() []interface{} {
			res := make([]interface{}, 0, len(t.CcvTickets))
			for _, e := range t.CcvTickets {
				res = append(res, e)
			}
			return res
		}(),
	}
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

// toMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageId":         string(t.MessageId),
		"encodedMessage":    string(t.EncodedMessage),
		"newSequenceNumber": (*big.Int)(t.NewSequenceNumber),
		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"verifierBlobs": func() []interface{} {
			res := make([]interface{}, 0, len(t.VerifierBlobs))
			for _, e := range t.VerifierBlobs {
				res = append(res, string(e))
			}
			return res
		}(),
		"messageSentObservers": func() []interface{} {
			res := make([]interface{}, 0, len(t.MessageSentObservers))
			for _, e := range t.MessageSentObservers {
				res = append(res, e.ToMap())
			}
			return res
		}(),
	}
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

// toMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) toMap() map[string]interface{} {
	return map[string]interface{}{

		"globalConfigCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.GlobalConfigCid).(mapper); ok {
				return m.toMap()
			}
			return t.GlobalConfigCid
		}(),
		"tokenAdminRegistryCid": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
				return m.toMap()
			}
			return t.TokenAdminRegistryCid
		}(),
		"destChainSelector": (*big.Int)(t.DestChainSelector),
		"hasTokenTransfer":  bool(t.HasTokenTransfer),
		"instrumentId":      *t.InstrumentId,
	}
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
