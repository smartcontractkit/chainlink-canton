package perpartyrouter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
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
	PackageName = "ccip-perpartyrouter"
	PackageID   = "8c46d52011243e9f1165c7439e06cec6cba58a8e10e3224837be80d5a660c1c4"
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

// CCIPSend is a Record type
type CCIPSend struct {
	Context             common.CCIPContext    `json:"context"`
	SendingMessageCid   types.CONTRACT_ID     `json:"sendingMessageCid"`
	FeeTokenInput       interfaces.TokenInput `json:"feeTokenInput"`
	FeeTokenHoldingCids []types.CONTRACT_ID   `json:"feeTokenHoldingCids"`
}

// ToMap converts CCIPSend to a map for DAML arguments
func (t CCIPSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["feeTokenInput"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeTokenInput).(mapper); ok {
			return m.toMap()
		}
		return t.FeeTokenInput
	}()

	m["feeTokenHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSend to hex string (Canton MCMS format)
func (t CCIPSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSend from hex string (Canton MCMS format)
func (t *CCIPSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCIPSendResult is a Record type
type CCIPSendResult struct {
	Router          types.CONTRACT_ID   `json:"router"`
	CcipMessageSent types.CONTRACT_ID   `json:"ccipMessageSent"`
	MessageId       types.TEXT          `json:"messageId"`
	FeeChangeCids   []types.CONTRACT_ID `json:"feeChangeCids"`
}

// ToMap converts CCIPSendResult to a map for DAML arguments
func (t CCIPSendResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	m["ccipMessageSent"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcipMessageSent).(mapper); ok {
			return m.toMap()
		}
		return t.CcipMessageSent
	}()

	m["messageId"] = string(t.MessageId)

	m["feeChangeCids"] = func() []any {
		res := make([]any, 0, len(t.FeeChangeCids))
		for _, e := range t.FeeChangeCids {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSendResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSendResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSendResult to hex string (Canton MCMS format)
func (t CCIPSendResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSendResult from hex string (Canton MCMS format)
func (t *CCIPSendResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CreateRouter is a Record type
type CreateRouter struct {
	PartyOwner types.PARTY `json:"partyOwner"`
	InstanceId types.TEXT  `json:"instanceId"`
}

// ToMap converts CreateRouter to a map for DAML arguments
func (t CreateRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["partyOwner"] = t.PartyOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t CreateRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CreateRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CreateRouter to hex string (Canton MCMS format)
func (t CreateRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CreateRouter from hex string (Canton MCMS format)
func (t *CreateRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CreateRouterResult is a Record type
type CreateRouterResult struct {
	Router  types.CONTRACT_ID `json:"router"`
	Factory types.CONTRACT_ID `json:"factory"`
}

// ToMap converts CreateRouterResult to a map for DAML arguments
func (t CreateRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	m["factory"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Factory).(mapper); ok {
			return m.toMap()
		}
		return t.Factory
	}()

	return m
}

func (t CreateRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CreateRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CreateRouterResult to hex string (Canton MCMS format)
func (t CreateRouterResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CreateRouterResult from hex string (Canton MCMS format)
func (t *CreateRouterResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Execute is a Record type
type Execute struct {
	Context              common.CCIPContext          `json:"context"`
	ExecutingMessageCid  types.CONTRACT_ID           `json:"executingMessageCid"`
	ReceiverRequiredCCVs []common.RawInstanceAddress `json:"receiverRequiredCCVs"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
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

func (t Execute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Execute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Execute to hex string (Canton MCMS format)
func (t Execute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Execute from hex string (Canton MCMS format)
func (t *Execute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteResult is a Record type
type ExecuteResult struct {
	Router                types.CONTRACT_ID            `json:"router"`
	TokenReceiveTicket    *types.CONTRACT_ID           `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID            `json:"executionStateChanged"`
	MessageId             types.TEXT                   `json:"messageId"`
	Message               common.MessageV1             `json:"message"`
	State                 common.MessageExecutionState `json:"state"`
}

// ToMap converts ExecuteResult to a map for DAML arguments
func (t ExecuteResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["router"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Router).(mapper); ok {
			return m.toMap()
		}
		return t.Router
	}()

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenReceiveTicket,
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
		}
	}

	m["executionStateChanged"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExecutionStateChanged).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutionStateChanged
	}()

	m["messageId"] = string(t.MessageId)

	m["message"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Message).(mapper); ok {
			return m.toMap()
		}
		return t.Message
	}()

	m["state"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.State).(mapper); ok {
			return m.toMap()
		}
		return t.State
	}()

	return m
}

func (t ExecuteResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteResult to hex string (Canton MCMS format)
func (t ExecuteResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteResult from hex string (Canton MCMS format)
func (t *ExecuteResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetExecutionState is a Record type
type GetExecutionState struct {
	MessageHash types.TEXT `json:"messageHash"`
}

// ToMap converts GetExecutionState to a map for DAML arguments
func (t GetExecutionState) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageHash"] = string(t.MessageHash)

	return m
}

func (t GetExecutionState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetExecutionState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetExecutionState to hex string (Canton MCMS format)
func (t GetExecutionState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetExecutionState from hex string (Canton MCMS format)
func (t *GetExecutionState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForExecute2 is a Record type
type GetRequiredCCVsForExecute2 struct {
	Context              common.CCIPContext          `json:"context"`
	ReceiverRequiredCCVs []common.RawInstanceAddress `json:"receiverRequiredCCVs"`
	SourceChainSelector  types.NUMERIC               `json:"sourceChainSelector"`
}

// ToMap converts GetRequiredCCVsForExecute2 to a map for DAML arguments
func (t GetRequiredCCVsForExecute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["sourceChainSelector"] = t.SourceChainSelector

	return m
}

func (t GetRequiredCCVsForExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForExecute2 to hex string (Canton MCMS format)
func (t GetRequiredCCVsForExecute2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForExecute2 from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForExecute2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForSend2 is a Record type
type GetRequiredCCVsForSend2 struct {
	Context           common.CCIPContext `json:"context"`
	DestChainSelector types.NUMERIC      `json:"destChainSelector"`
}

// ToMap converts GetRequiredCCVsForSend2 to a map for DAML arguments
func (t GetRequiredCCVsForSend2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["destChainSelector"] = t.DestChainSelector

	return m
}

func (t GetRequiredCCVsForSend2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForSend2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForSend2 to hex string (Canton MCMS format)
func (t GetRequiredCCVsForSend2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForSend2 from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForSend2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetSequenceNumber is a Record type
type GetSequenceNumber struct {
	DestChainSelector types.NUMERIC `json:"destChainSelector"`
}

// ToMap converts GetSequenceNumber to a map for DAML arguments
func (t GetSequenceNumber) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	return m
}

func (t GetSequenceNumber) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetSequenceNumber) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetSequenceNumber to hex string (Canton MCMS format)
func (t GetSequenceNumber) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetSequenceNumber from hex string (Canton MCMS format)
func (t *GetSequenceNumber) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HasRouter is a Record type
type HasRouter struct {
	PartyOwner types.PARTY `json:"partyOwner"`
	Caller     types.PARTY `json:"caller"`
}

// ToMap converts HasRouter to a map for DAML arguments
func (t HasRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["partyOwner"] = t.PartyOwner.ToMap()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t HasRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *HasRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes HasRouter to hex string (Canton MCMS format)
func (t HasRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HasRouter from hex string (Canton MCMS format)
func (t *HasRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// HasRouterMCMSParams is HasRouter without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type HasRouterMCMSParams struct {
	PartyOwner types.PARTY `json:"partyOwner"`
}

// MarshalHex encodes HasRouterMCMSParams to hex string for MCMS operationData.
func (t HasRouterMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes HasRouterMCMSParams from hex string.
func (t *HasRouterMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouter is a Template type
type PerPartyRouter struct {
	InstanceId              types.TEXT         `json:"instanceId"`
	CcipOwner               types.PARTY        `json:"ccipOwner"`
	PartyOwner              types.PARTY        `json:"partyOwner"`
	Deps                    PerPartyRouterDeps `json:"deps"`
	OutboundSequenceNumbers types.GENMAP       `json:"outboundSequenceNumbers"`
	ExecutionStates         types.GENMAP       `json:"executionStates"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t PerPartyRouter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t PerPartyRouter) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() any {
		if t.OutboundSequenceNumbers == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionStates"] = func() any {
		if t.ExecutionStates == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ExecutionStates}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t PerPartyRouter) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() any {
		if t.OutboundSequenceNumbers == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executionStates"] = func() any {
		if t.ExecutionStates == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.ExecutionStates}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t PerPartyRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouter to hex string (Canton MCMS format)
func (t PerPartyRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouter from hex string (Canton MCMS format)
func (t *PerPartyRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for PerPartyRouter

// PrepareSend exercises the PrepareSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareSend(contractID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendWithPackageID exercises the PrepareSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareSendWithPackageID(contractID string, packageID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// CCIPSend exercises the CCIPSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) CCIPSend(contractID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// CCIPSendWithPackageID exercises the CCIPSend choice using the provided package ID instead of package name
func (t PerPartyRouter) CCIPSendWithPackageID(contractID string, packageID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareExecute(contractID string, args PrepareExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// Execute exercises the Execute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t PerPartyRouter) ExecuteWithPackageID(contractID string, packageID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionState exercises the GetExecutionState choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetExecutionState(contractID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionStateWithPackageID exercises the GetExecutionState choice using the provided package ID instead of package name
func (t PerPartyRouter) GetExecutionStateWithPackageID(contractID string, packageID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumber exercises the GetSequenceNumber choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetSequenceNumber(contractID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumberWithPackageID exercises the GetSequenceNumber choice using the provided package ID instead of package name
func (t PerPartyRouter) GetSequenceNumberWithPackageID(contractID string, packageID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterDeps is a Record type
type PerPartyRouterDeps struct {
	OnRamp             common.RawInstanceAddress `json:"onRamp"`
	OffRamp            common.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       common.RawInstanceAddress `json:"globalConfig"`
	TokenAdminRegistry common.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          common.RawInstanceAddress `json:"feeQuoter"`
	RmnRemote          common.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts PerPartyRouterDeps to a map for DAML arguments
func (t PerPartyRouterDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["onRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OnRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OnRamp
	}()

	m["offRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OffRamp
	}()

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
	}()

	m["tokenAdminRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistry
	}()

	m["feeQuoter"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoter).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoter
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	return m
}

func (t PerPartyRouterDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterDeps to hex string (Canton MCMS format)
func (t PerPartyRouterDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterDeps from hex string (Canton MCMS format)
func (t *PerPartyRouterDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouterFactory is a Template type
type PerPartyRouterFactory struct {
	InstanceId        types.TEXT         `json:"instanceId"`
	CcipOwner         types.PARTY        `json:"ccipOwner"`
	Deps              PerPartyRouterDeps `json:"deps"`
	RegisteredRouters types.GENMAP       `json:"registeredRouters"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t PerPartyRouterFactory) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouterFactory) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t PerPartyRouterFactory) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registeredRouters"] = func() any {
		if t.RegisteredRouters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RegisteredRouters}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t PerPartyRouterFactory) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["registeredRouters"] = func() any {
		if t.RegisteredRouters == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RegisteredRouters}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t PerPartyRouterFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PerPartyRouterFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PerPartyRouterFactory to hex string (Canton MCMS format)
func (t PerPartyRouterFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PerPartyRouterFactory from hex string (Canton MCMS format)
func (t *PerPartyRouterFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for PerPartyRouterFactory

// CreateRouter exercises the CreateRouter choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) CreateRouter(contractID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// CreateRouterWithPackageID exercises the CreateRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) CreateRouterWithPackageID(contractID string, packageID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouter exercises the HasRouter choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) HasRouter(contractID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouterWithPackageID exercises the HasRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) HasRouterWithPackageID(contractID string, packageID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// PrepareExecute2 is a Record type
type PrepareExecute2 struct {
	Context            common.CCIPContext `json:"context"`
	EncodedMessage     types.TEXT         `json:"encodedMessage"`
	ReceiverParty      types.PARTY        `json:"receiverParty"`
	TokenReceiverParty *types.PARTY       `json:"tokenReceiverParty" hex:"optional"`
	Caller             types.PARTY        `json:"caller"`
}

// ToMap converts PrepareExecute2 to a map for DAML arguments
func (t PrepareExecute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["encodedMessage"] = string(t.EncodedMessage)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	if t.TokenReceiverParty != nil {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
			"value": (*t.TokenReceiverParty).ToMap(),
		}
	} else {
		m["tokenReceiverParty"] = map[string]any{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PrepareExecute2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PrepareExecute2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PrepareExecute2 to hex string (Canton MCMS format)
func (t PrepareExecute2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecute2 from hex string (Canton MCMS format)
func (t *PrepareExecute2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareExecute2MCMSParams is PrepareExecute2 without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type PrepareExecute2MCMSParams struct {
	Context            common.CCIPContext `json:"context"`
	EncodedMessage     types.TEXT         `json:"encodedMessage"`
	ReceiverParty      types.PARTY        `json:"receiverParty"`
	TokenReceiverParty *types.PARTY       `json:"tokenReceiverParty" hex:"optional"`
}

// MarshalHex encodes PrepareExecute2MCMSParams to hex string for MCMS operationData.
func (t PrepareExecute2MCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecute2MCMSParams from hex string.
func (t *PrepareExecute2MCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareSend is a Record type
type PrepareSend struct {
	Context             common.CCIPContext                        `json:"context"`
	DestChainSelector   types.NUMERIC                             `json:"destChainSelector"`
	Receiver            types.TEXT                                `json:"receiver"`
	Payload             types.TEXT                                `json:"payload"`
	CcipReceiveGasLimit types.INT64                               `json:"ccipReceiveGasLimit"`
	BlockConfirmations  *types.INT64                              `json:"blockConfirmations" hex:"optional"`
	SenderRequiredCCVs  []common.RawInstanceAddress               `json:"senderRequiredCCVs"`
	TokenInstrumentId   *splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId" hex:"optional"`
	TokenReceiver       *types.TEXT                               `json:"tokenReceiver" hex:"optional"`
	TokenArgs           types.TEXT                                `json:"tokenArgs"`
	FeeToken            splice_api_token_holding_v1.InstrumentId  `json:"feeToken"`
}

// ToMap converts PrepareSend to a map for DAML arguments
func (t PrepareSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Context).(mapper); ok {
			return m.toMap()
		}
		return t.Context
	}()

	m["destChainSelector"] = t.DestChainSelector

	m["receiver"] = string(t.Receiver)

	m["payload"] = string(t.Payload)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	if t.BlockConfirmations != nil {
		m["blockConfirmations"] = map[string]any{
			"_type": "optional",
			"value": int64(*t.BlockConfirmations),
		}
	} else {
		m["blockConfirmations"] = map[string]any{
			"_type": "optional",
		}
	}

	m["senderRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.SenderRequiredCCVs))
		for _, e := range t.SenderRequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	if t.TokenInstrumentId != nil {
		m["tokenInstrumentId"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenInstrumentId,
		}
	} else {
		m["tokenInstrumentId"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.TokenReceiver != nil {
		m["tokenReceiver"] = map[string]any{
			"_type": "optional",
			"value": string(*t.TokenReceiver),
		}
	} else {
		m["tokenReceiver"] = map[string]any{
			"_type": "optional",
		}
	}

	m["tokenArgs"] = string(t.TokenArgs)

	m["feeToken"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	return m
}

func (t PrepareSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PrepareSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PrepareSend to hex string (Canton MCMS format)
func (t PrepareSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareSend from hex string (Canton MCMS format)
func (t *PrepareSend) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	CCIPSend(args CCIPSend) (*bind.EncodedChoice, error)
	CreateRouter(args CreateRouter) (*bind.EncodedChoice, error)
	Execute(args Execute) (*bind.EncodedChoice, error)
	GetExecutionState(args GetExecutionState) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecute2(args GetRequiredCCVsForExecute2) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSend2(args GetRequiredCCVsForSend2) (*bind.EncodedChoice, error)
	GetSequenceNumber(args GetSequenceNumber) (*bind.EncodedChoice, error)
	HasRouter(args HasRouter) (*bind.EncodedChoice, error)
	HasRouterMCMSParams(args HasRouterMCMSParams) (*bind.EncodedChoice, error)
	PrepareExecute2(args PrepareExecute2) (*bind.EncodedChoice, error)
	PrepareExecute2MCMSParams(args PrepareExecute2MCMSParams) (*bind.EncodedChoice, error)
	PrepareSend(args PrepareSend) (*bind.EncodedChoice, error)
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

// CCIPSend encodes parameters for the CCIPSend choice.
func (e *encoder) CCIPSend(args CCIPSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CCIPSend", args)
}

// CreateRouter encodes parameters for the CreateRouter choice.
func (e *encoder) CreateRouter(args CreateRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CreateRouter", args)
}

// Execute encodes parameters for the Execute choice.
func (e *encoder) Execute(args Execute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Execute", args)
}

// GetExecutionState encodes parameters for the GetExecutionState choice.
func (e *encoder) GetExecutionState(args GetExecutionState) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetExecutionState", args)
}

// GetRequiredCCVsForExecute2 encodes parameters for the GetRequiredCCVsForExecute2 choice.
func (e *encoder) GetRequiredCCVsForExecute2(args GetRequiredCCVsForExecute2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForExecute2", args)
}

// GetRequiredCCVsForSend2 encodes parameters for the GetRequiredCCVsForSend2 choice.
func (e *encoder) GetRequiredCCVsForSend2(args GetRequiredCCVsForSend2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForSend2", args)
}

// GetSequenceNumber encodes parameters for the GetSequenceNumber choice.
func (e *encoder) GetSequenceNumber(args GetSequenceNumber) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetSequenceNumber", args)
}

// HasRouter encodes parameters for the HasRouter choice.
func (e *encoder) HasRouter(args HasRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HasRouter", args)
}

// HasRouterMCMSParams encodes MCMS parameters (without Caller) for the HasRouter choice.
func (e *encoder) HasRouterMCMSParams(args HasRouterMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("HasRouter", args)
}

// PrepareExecute2 encodes parameters for the PrepareExecute2 choice.
func (e *encoder) PrepareExecute2(args PrepareExecute2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute2", args)
}

// PrepareExecute2MCMSParams encodes MCMS parameters (without Caller) for the PrepareExecute2 choice.
func (e *encoder) PrepareExecute2MCMSParams(args PrepareExecute2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute2", args)
}

// PrepareSend encodes parameters for the PrepareSend choice.
func (e *encoder) PrepareSend(args PrepareSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareSend", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
