package offramp

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

const PackageID = "9352d0543ef077e2e15006726314842d01d8d6a92ebe98a436c932e206f94994"
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

// ExecuteFromRouter is a Record type
type ExecuteFromRouter struct {
	RouterPartyOwner PARTY `json:"routerPartyOwner"`

	ReceiverRequiredCCVs []TEXT `json:"receiverRequiredCCVs"`

	GlobalConfigCid CONTRACT_ID `json:"globalConfigCid"`

	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`

	EncodedMessage TEXT `json:"encodedMessage"`

	CcvVerifyTickets []CONTRACT_ID `json:"ccvVerifyTickets"`

	TokenPoolCCVTicket *CONTRACT_ID `json:"tokenPoolCCVTicket"`
}

// ToMap converts ExecuteFromRouter to a map for DAML arguments
func (t ExecuteFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, string(e))
		}
		return res
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

	m["encodedMessage"] = string(t.EncodedMessage)

	m["ccvVerifyTickets"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.CcvVerifyTickets))
		for _, e := range t.CcvVerifyTickets {
			res = append(res, e)
		}
		return res
	}()

	if t.TokenPoolCCVTicket != nil {
		m["tokenPoolCCVTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenPoolCCVTicket,
		}
	} else {
		m["tokenPoolCCVTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

func (t ExecuteFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ExecuteFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// ExecuteFromRouterResult is a Record type
type ExecuteFromRouterResult struct {
	MessageHash TEXT `json:"messageHash"`

	Message MessageV1 `json:"message"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	SequenceNumber NUMERIC `json:"sequenceNumber"`

	TokenReceiveTicket *CONTRACT_ID `json:"tokenReceiveTicket"`
}

// ToMap converts ExecuteFromRouterResult to a map for DAML arguments
func (t ExecuteFromRouterResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["messageHash"] = string(t.MessageHash)

	m["message"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]interface{}{
			"_type": "optional",
			"value": *t.TokenReceiveTicket,
		}
	} else {
		m["tokenReceiveTicket"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	return m
}

func (t ExecuteFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *ExecuteFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	GlobalConfigCid CONTRACT_ID `json:"globalConfigCid"`

	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`

	ReceiverRequiredCCVs []TEXT `json:"receiverRequiredCCVs"`

	SourceChainSelector NUMERIC `json:"sourceChainSelector"`

	HasTokenTransfer BOOL `json:"hasTokenTransfer"`

	InstrumentId *InstrumentId `json:"instrumentId"`
}

// ToMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) ToMap() map[string]interface{} {
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

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, string(e))
		}
		return res
	}()

	m["sourceChainSelector"] = t.SourceChainSelector

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

func (t GetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// OffRamp is a Template type
type OffRamp struct {
	CcipOwner PARTY `json:"ccipOwner"`

	InstanceId TEXT `json:"instanceId"`
}

// GetTemplateID returns the template ID for this template
func (t OffRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp")
}

// CreateCommand returns a CreateCommand for this template
func (t OffRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

func (t OffRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *OffRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for OffRamp

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this OffRamp contract
func (t OffRamp) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),

		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",

		Arguments: argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OffRamp contract
func (t OffRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),

		ContractID: contractID,
		Choice:     "Archive",

		Arguments: map[string]interface{}{},
	}
}

// ExecuteFromRouter exercises the ExecuteFromRouter choice on this OffRamp contract
func (t OffRamp) ExecuteFromRouter(contractID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{

		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageID, "CCIP.OffRamp", "OffRamp"),

		ContractID: contractID,
		Choice:     "ExecuteFromRouter",

		Arguments: argsToMap(args),
	}
}
