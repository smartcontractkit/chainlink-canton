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

// ExecuteFromRouter is a Record type
type ExecuteFromRouter struct {
	RouterPartyOwner      PARTY         `json:"routerPartyOwner"`
	ReceiverRequiredCCVs  []TEXT        `json:"receiverRequiredCCVs"`
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	EncodedMessage        TEXT          `json:"encodedMessage"`
	CcvVerifyTickets      []CONTRACT_ID `json:"ccvVerifyTickets"`
}

// toMap converts ExecuteFromRouter to a map for DAML arguments
func (t ExecuteFromRouter) toMap() map[string]interface{} {
	return map[string]interface{}{

		"routerPartyOwner": t.RouterPartyOwner.ToMap(),
		"receiverRequiredCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
			for _, e := range t.ReceiverRequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
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
		"encodedMessage": string(t.EncodedMessage),
		"ccvVerifyTickets": func() []interface{} {
			res := make([]interface{}, 0, len(t.CcvVerifyTickets))
			for _, e := range t.CcvVerifyTickets {
				res = append(res, e)
			}
			return res
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ExecuteFromRouter using JsonCodec
func (t ExecuteFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecuteFromRouter using JsonCodec
func (t *ExecuteFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecuteFromRouterResult is a Record type
type ExecuteFromRouterResult struct {
	MessageHash         TEXT        `json:"messageHash"`
	Message             MessageV1   `json:"message"`
	SourceChainSelector NUMERIC     `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC     `json:"sequenceNumber"`
	TokenReceiveTicket  CONTRACT_ID `json:"tokenReceiveTicket"`
}

// toMap converts ExecuteFromRouterResult to a map for DAML arguments
func (t ExecuteFromRouterResult) toMap() map[string]interface{} {
	return map[string]interface{}{

		"messageHash": string(t.MessageHash),
		"message": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.Message).(mapper); ok {
				return m.toMap()
			}
			return t.Message
		}(),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"sequenceNumber":      (*big.Int)(t.SequenceNumber),
		"tokenReceiveTicket": func() interface{} {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(t.TokenReceiveTicket).(mapper); ok {
				return m.toMap()
			}
			return t.TokenReceiveTicket
		}(),
	}
}

// MarshalJSON implements custom JSON marshaling for ExecuteFromRouterResult using JsonCodec
func (t ExecuteFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for ExecuteFromRouterResult using JsonCodec
func (t *ExecuteFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForExecute2 is a Record type
type GetRequiredCCVsForExecute2 struct {
	GlobalConfigCid       CONTRACT_ID   `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID   `json:"tokenAdminRegistryCid"`
	ReceiverRequiredCCVs  []TEXT        `json:"receiverRequiredCCVs"`
	SourceChainSelector   NUMERIC       `json:"sourceChainSelector"`
	HasTokenTransfer      BOOL          `json:"hasTokenTransfer"`
	InstrumentId          *InstrumentId `json:"instrumentId"`
}

// toMap converts GetRequiredCCVsForExecute2 to a map for DAML arguments
func (t GetRequiredCCVsForExecute2) toMap() map[string]interface{} {
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
		"receiverRequiredCCVs": func() []interface{} {
			res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
			for _, e := range t.ReceiverRequiredCCVs {
				res = append(res, string(e))
			}
			return res
		}(),
		"sourceChainSelector": (*big.Int)(t.SourceChainSelector),
		"hasTokenTransfer":    bool(t.HasTokenTransfer),
		"instrumentId":        *t.InstrumentId,
	}
}

// MarshalJSON implements custom JSON marshaling for GetRequiredCCVsForExecute2 using JsonCodec
func (t GetRequiredCCVsForExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for GetRequiredCCVsForExecute2 using JsonCodec
func (t *GetRequiredCCVsForExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// OffRamp is a Template type
type OffRamp struct {
	CcipOwner  PARTY `json:"ccipOwner"`
	InstanceId TEXT  `json:"instanceId"`
}

// GetTemplateID returns the template ID for this template
func (t OffRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp")
}

// CreateCommand returns a CreateCommand for this template
func (t OffRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	args["ccipOwner"] = t.CcipOwner.ToMap()

	args["instanceId"] = string(t.InstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// MarshalJSON implements custom JSON marshaling for OffRamp using JsonCodec
func (t OffRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

// UnmarshalJSON implements custom JSON unmarshaling for OffRamp using JsonCodec
func (t *OffRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for OffRamp

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this OffRamp contract
func (t OffRamp) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OffRamp contract
func (t OffRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ExecuteFromRouter exercises the ExecuteFromRouter choice on this OffRamp contract
func (t OffRamp) ExecuteFromRouter(contractID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}
