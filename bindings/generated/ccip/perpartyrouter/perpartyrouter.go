package perpartyrouter

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	client "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	interfaces "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/interfaces"
	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageID   = "09689fed574b4a59fd891f11969ace97167cd6290a88877952702c77a8faa66d"
	SDKVersion  = "3.4.10"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	MaxExecutedMessagesSize = types.INT64(25000)
)

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

// AddCustomObservers is a Record type
type AddCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObservers to a map for DAML arguments
func (t AddCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObservers to hex string (Canton MCMS format)
func (t AddCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObservers from hex string (Canton MCMS format)
func (t *AddCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCustomObserversParams is a Record type
type AddCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObserversParams to a map for DAML arguments
func (t AddCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t AddCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObserversParams to hex string (Canton MCMS format)
func (t AddCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObserversParams from hex string (Canton MCMS format)
func (t *AddCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ArchivedExecutedMessages is a Template type
type ArchivedExecutedMessages struct {
	CcipOwner                types.PARTY `json:"ccipOwner"`
	PartyOwner               types.PARTY `json:"partyOwner"`
	PerPartyRouterInstanceId types.TEXT  `json:"perPartyRouterInstanceId"`
	ArchiveIndex             types.INT64 `json:"archiveIndex"`
	ExecutedMessages         types.SET   `json:"executedMessages"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t ArchivedExecutedMessages) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "ArchivedExecutedMessages")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ArchivedExecutedMessages) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.PerPartyRouter", "ArchivedExecutedMessages")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t ArchivedExecutedMessages) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["perPartyRouterInstanceId"] = string(t.PerPartyRouterInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["archiveIndex"] = int64(t.ArchiveIndex)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executedMessages"] = model.NestedToDAMLValue(t.ExecutedMessages)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t ArchivedExecutedMessages) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["partyOwner"] = t.PartyOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["perPartyRouterInstanceId"] = string(t.PerPartyRouterInstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["archiveIndex"] = int64(t.ArchiveIndex)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executedMessages"] = model.NestedToDAMLValue(t.ExecutedMessages)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t ArchivedExecutedMessages) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ArchivedExecutedMessages) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ArchivedExecutedMessages to hex string (Canton MCMS format)
func (t ArchivedExecutedMessages) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ArchivedExecutedMessages from hex string (Canton MCMS format)
func (t *ArchivedExecutedMessages) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for ArchivedExecutedMessages

// IsExecuted exercises the IsExecuted choice on this ArchivedExecutedMessages contract
// This method uses the package name in the template ID
func (t ArchivedExecutedMessages) IsExecuted(contractID string, args IsExecuted) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "IsExecuted",
		Arguments:  argsToMap(args),
	}
}

// IsExecutedWithPackageID exercises the IsExecuted choice using the provided package ID instead of package name
func (t ArchivedExecutedMessages) IsExecutedWithPackageID(contractID string, packageID string, args IsExecuted) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "IsExecuted",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ArchivedExecutedMessages contract
// This method uses the package name in the template ID
func (t ArchivedExecutedMessages) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ArchivedExecutedMessages) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CCIPSend is a Record type
type CCIPSend struct {
	Context                 splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	SendingMessageCid       types.CONTRACT_ID                          `json:"sendingMessageCid"`
	FeeTokenHoldingCids     []types.CONTRACT_ID                        `json:"feeTokenHoldingCids"`
	FeeTokenTransferFactory types.CONTRACT_ID                          `json:"feeTokenTransferFactory"`
	FeeTokenExtraArgs       splice_api_token_metadata_v1.ExtraArgs     `json:"feeTokenExtraArgs"`
}

// ToMap converts CCIPSend to a map for DAML arguments
func (t CCIPSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeTokenHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["feeTokenTransferFactory"] = model.NestedToDAMLValue(t.FeeTokenTransferFactory)

	m["feeTokenExtraArgs"] = model.NestedToDAMLValue(t.FeeTokenExtraArgs)

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

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["ccipMessageSent"] = model.NestedToDAMLValue(t.CcipMessageSent)

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

	m["router"] = model.NestedToDAMLValue(t.Router)

	m["factory"] = model.NestedToDAMLValue(t.Factory)

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
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

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

	m["router"] = model.NestedToDAMLValue(t.Router)

	if t.TokenReceiveTicket != nil {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenReceiveTicket),
		}
	} else {
		m["tokenReceiveTicket"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["executionStateChanged"] = model.NestedToDAMLValue(t.ExecutionStateChanged)

	m["messageId"] = string(t.MessageId)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["state"] = model.NestedToDAMLValue(t.State)

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

// FactorySetDeps is a Record type
type FactorySetDeps struct {
	NewDeps SetDepsParams3 `json:"newDeps"`
}

// ToMap converts FactorySetDeps to a map for DAML arguments
func (t FactorySetDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = model.NestedToDAMLValue(t.NewDeps)

	return m
}

func (t FactorySetDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FactorySetDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FactorySetDeps to hex string (Canton MCMS format)
func (t FactorySetDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FactorySetDeps from hex string (Canton MCMS format)
func (t *FactorySetDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeFee2 is a Record type
type FinalizeFee2 struct {
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
}

// ToMap converts FinalizeFee2 to a map for DAML arguments
func (t FinalizeFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	return m
}

func (t FinalizeFee2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeFee2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeFee2 to hex string (Canton MCMS format)
func (t FinalizeFee2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeFee2 from hex string (Canton MCMS format)
func (t *FinalizeFee2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetExecutionState is a Record type
type GetExecutionState struct {
	MessageHash types.TEXT  `json:"messageHash"`
	Caller      types.PARTY `json:"caller"`
}

// ToMap converts GetExecutionState to a map for DAML arguments
func (t GetExecutionState) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageHash"] = string(t.MessageHash)

	m["caller"] = t.Caller.ToMap()

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

// GetExecutionStateMCMSParams is GetExecutionState without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetExecutionStateMCMSParams struct {
	MessageHash types.TEXT `json:"messageHash"`
}

// MarshalHex encodes GetExecutionStateMCMSParams to hex string for MCMS operationData.
func (t GetExecutionStateMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetExecutionStateMCMSParams from hex string.
func (t *GetExecutionStateMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFee is a Record type
type GetFee struct {
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Message           client.Canton2AnyMessage                   `json:"message"`
	CcvFeeQuotes      []common.CrossChainVerifierFeeQuote        `json:"ccvFeeQuotes"`
	TokenPoolFeeQuote *interfaces.TokenPoolFeeQuote              `json:"tokenPoolFeeQuote" hex:"optional"`
	ExecutorFeeQuote  *common.ExecutorFeeQuote                   `json:"executorFeeQuote" hex:"optional"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["ccvFeeQuotes"] = func() []any {
		res := make([]any, 0, len(t.CcvFeeQuotes))
		for _, e := range t.CcvFeeQuotes {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenPoolFeeQuote != nil {
		m["tokenPoolFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolFeeQuote),
		}
	} else {
		m["tokenPoolFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorFeeQuote != nil {
		m["executorFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorFeeQuote),
		}
	} else {
		m["executorFeeQuote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t GetFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFee to hex string (Canton MCMS format)
func (t GetFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFee from hex string (Canton MCMS format)
func (t *GetFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Message                   common.MessageV1                           `json:"message"`
	ReceiverRequiredCCVs      []mcms.RawInstanceAddress                  `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []mcms.RawInstanceAddress                  `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	TokenPoolRequiredCCVs     []mcms.RawInstanceAddress                  `json:"tokenPoolRequiredCCVs"`
}

// ToMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverOptionalCCVs))
		for _, e := range t.ReceiverOptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalThreshold"] = int64(t.ReceiverOptionalThreshold)

	m["tokenPoolRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.TokenPoolRequiredCCVs))
		for _, e := range t.TokenPoolRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t GetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForExecute to hex string (Canton MCMS format)
func (t GetRequiredCCVsForExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForExecute from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	DestChainSelector     types.NUMERIC                              `json:"destChainSelector"`
	Message               client.Canton2AnyMessage                   `json:"message"`
	TokenPoolRequiredCCVs []mcms.RawInstanceAddress                  `json:"tokenPoolRequiredCCVs"`
}

// ToMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["tokenPoolRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.TokenPoolRequiredCCVs))
		for _, e := range t.TokenPoolRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t GetRequiredCCVsForSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForSend to hex string (Canton MCMS format)
func (t GetRequiredCCVsForSend) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForSend from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForSend) UnmarshalHex(data string) error {
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

// IsExecuted is a Record type
type IsExecuted struct {
	MessageHash types.TEXT `json:"messageHash"`
}

// ToMap converts IsExecuted to a map for DAML arguments
func (t IsExecuted) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageHash"] = string(t.MessageHash)

	return m
}

func (t IsExecuted) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *IsExecuted) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes IsExecuted to hex string (Canton MCMS format)
func (t IsExecuted) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes IsExecuted from hex string (Canton MCMS format)
func (t *IsExecuted) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PerPartyRouter is a Template type
type PerPartyRouter struct {
	InstanceId                   types.TEXT                      `json:"instanceId"`
	CcipOwner                    types.PARTY                     `json:"ccipOwner"`
	PartyOwner                   types.PARTY                     `json:"partyOwner"`
	Deps                         PerPartyRouterDeps              `json:"deps"`
	OutboundSequenceNumbers      map[types.NUMERIC]types.NUMERIC `json:"outboundSequenceNumbers"`
	ExecutedMessages             types.SET                       `json:"executedMessages"`
	ArchivedExecutionContractIds []types.CONTRACT_ID             `json:"archivedExecutionContractIds"`
	CustomObservers              []types.PARTY                   `json:"customObservers"`
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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() any {
		if t.OutboundSequenceNumbers == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executedMessages"] = model.NestedToDAMLValue(t.ExecutedMessages)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["archivedExecutionContractIds"] = func() []any {
		res := make([]any, 0, len(t.ArchivedExecutionContractIds))
		for _, e := range t.ArchivedExecutionContractIds {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["outboundSequenceNumbers"] = func() any {
		if t.OutboundSequenceNumbers == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.OutboundSequenceNumbers}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["executedMessages"] = model.NestedToDAMLValue(t.ExecutedMessages)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["archivedExecutionContractIds"] = func() []any {
		res := make([]any, 0, len(t.ArchivedExecutionContractIds))
		for _, e := range t.ArchivedExecutionContractIds {
			res = append(res, e)
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
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

// GetFee exercises the GetFee choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t PerPartyRouter) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

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

// FinalizeFee exercises the FinalizeFee choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) FinalizeFee(contractID string, args FinalizeFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFeeWithPackageID exercises the FinalizeFee choice using the provided package ID instead of package name
func (t PerPartyRouter) FinalizeFeeWithPackageID(contractID string, packageID string, args FinalizeFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this PerPartyRouter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
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

// SetDeps exercises the SetDeps choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) SetDeps(contractID string, args SetDeps3) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t PerPartyRouter) SetDepsWithPackageID(contractID string, packageID string, args SetDeps3) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObservers exercises the AddCustomObservers choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) AddCustomObservers(contractID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObserversWithPackageID exercises the AddCustomObservers choice using the provided package ID instead of package name
func (t PerPartyRouter) AddCustomObserversWithPackageID(contractID string, packageID string, args AddCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObservers exercises the RemoveCustomObservers choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) RemoveCustomObservers(contractID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObserversWithPackageID exercises the RemoveCustomObservers choice using the provided package ID instead of package name
func (t PerPartyRouter) RemoveCustomObserversWithPackageID(contractID string, packageID string, args RemoveCustomObservers) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this PerPartyRouter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouter) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t PerPartyRouter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for PerPartyRouter

var _ mcms.IMCMSReceiver = (*PerPartyRouter)(nil)

// PerPartyRouterDeps is a Record type
type PerPartyRouterDeps struct {
	OnRamp             mcms.RawInstanceAddress `json:"onRamp"`
	OffRamp            mcms.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          mcms.RawInstanceAddress `json:"feeQuoter"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts PerPartyRouterDeps to a map for DAML arguments
func (t PerPartyRouterDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["onRamp"] = model.NestedToDAMLValue(t.OnRamp)

	m["offRamp"] = model.NestedToDAMLValue(t.OffRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

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
	InstanceId        types.TEXT                 `json:"instanceId"`
	CcipOwner         types.PARTY                `json:"ccipOwner"`
	Deps              PerPartyRouterDeps         `json:"deps"`
	RegisteredRouters map[types.PARTY]types.TEXT `json:"registeredRouters"`
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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

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
	args["deps"] = model.NestedToDAMLValue(t.Deps)

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

// Archive exercises the Archive choice on this PerPartyRouterFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
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

// FactorySetDeps exercises the FactorySetDeps choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) FactorySetDeps(contractID string, args FactorySetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "FactorySetDeps",
		Arguments:  argsToMap(args),
	}
}

// FactorySetDepsWithPackageID exercises the FactorySetDeps choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) FactorySetDepsWithPackageID(contractID string, packageID string, args FactorySetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "FactorySetDeps",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this PerPartyRouterFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for PerPartyRouterFactory

var _ mcms.IMCMSReceiver = (*PerPartyRouterFactory)(nil)

// PrepareExecute2 is a Record type
type PrepareExecute2 struct {
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	EncodedMessage            types.TEXT                                 `json:"encodedMessage"`
	ReceiverParty             types.PARTY                                `json:"receiverParty"`
	TokenReceiverParty        *types.PARTY                               `json:"tokenReceiverParty" hex:"optional"`
	ReceiverRequiredCCVs      []mcms.RawInstanceAddress                  `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []mcms.RawInstanceAddress                  `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	ReceiverFinalityConfig    common.FinalityConfig                      `json:"receiverFinalityConfig"`
	Caller                    types.PARTY                                `json:"caller"`
}

// ToMap converts PrepareExecute2 to a map for DAML arguments
func (t PrepareExecute2) ToMap() map[string]any {
	m := make(map[string]any)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
			"value": nil,
		}
	}

	m["receiverRequiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.ReceiverOptionalCCVs))
		for _, e := range t.ReceiverOptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["receiverOptionalThreshold"] = int64(t.ReceiverOptionalThreshold)

	m["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

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
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	EncodedMessage            types.TEXT                                 `json:"encodedMessage"`
	ReceiverParty             types.PARTY                                `json:"receiverParty"`
	TokenReceiverParty        *types.PARTY                               `json:"tokenReceiverParty" hex:"optional"`
	ReceiverRequiredCCVs      []mcms.RawInstanceAddress                  `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []mcms.RawInstanceAddress                  `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	ReceiverFinalityConfig    common.FinalityConfig                      `json:"receiverFinalityConfig"`
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
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  client.Canton2AnyMessage                   `json:"message"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts PrepareSend to a map for DAML arguments
func (t PrepareSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// RemoveCustomObservers is a Record type
type RemoveCustomObservers struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObservers to a map for DAML arguments
func (t RemoveCustomObservers) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObservers) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObservers) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObservers to hex string (Canton MCMS format)
func (t RemoveCustomObservers) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObservers from hex string (Canton MCMS format)
func (t *RemoveCustomObservers) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveCustomObserversParams is a Record type
type RemoveCustomObserversParams struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObserversParams to a map for DAML arguments
func (t RemoveCustomObserversParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["parties"] = func() []any {
		res := make([]any, 0, len(t.Parties))
		for _, e := range t.Parties {
			res = append(res, e.ToMap())
		}
		return res
	}()

	return m
}

func (t RemoveCustomObserversParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObserversParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObserversParams to hex string (Canton MCMS format)
func (t RemoveCustomObserversParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObserversParams from hex string (Canton MCMS format)
func (t *RemoveCustomObserversParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDeps3 is a Record type
type SetDeps3 struct {
	NewDeps SetDepsParams3 `json:"newDeps"`
}

// ToMap converts SetDeps3 to a map for DAML arguments
func (t SetDeps3) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = model.NestedToDAMLValue(t.NewDeps)

	return m
}

func (t SetDeps3) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDeps3) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDeps3 to hex string (Canton MCMS format)
func (t SetDeps3) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDeps3 from hex string (Canton MCMS format)
func (t *SetDeps3) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDepsParams3 is a Record type
type SetDepsParams3 struct {
	OnRamp             *mcms.RawInstanceAddress `json:"onRamp" hex:"optional"`
	OffRamp            *mcms.RawInstanceAddress `json:"offRamp" hex:"optional"`
	GlobalConfig       *mcms.RawInstanceAddress `json:"globalConfig" hex:"optional"`
	TokenAdminRegistry *mcms.RawInstanceAddress `json:"tokenAdminRegistry" hex:"optional"`
	FeeQuoter          *mcms.RawInstanceAddress `json:"feeQuoter" hex:"optional"`
	RmnRemote          *mcms.RawInstanceAddress `json:"rmnRemote" hex:"optional"`
}

// ToMap converts SetDepsParams3 to a map for DAML arguments
func (t SetDepsParams3) ToMap() map[string]any {
	m := make(map[string]any)

	if t.OnRamp != nil {
		m["onRamp"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OnRamp),
		}
	} else {
		m["onRamp"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.OffRamp != nil {
		m["offRamp"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.OffRamp),
		}
	} else {
		m["offRamp"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.GlobalConfig != nil {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.GlobalConfig),
		}
	} else {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.TokenAdminRegistry != nil {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenAdminRegistry),
		}
	} else {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.FeeQuoter != nil {
		m["feeQuoter"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.FeeQuoter),
		}
	} else {
		m["feeQuoter"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.RmnRemote != nil {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.RmnRemote),
		}
	} else {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	return m
}

func (t SetDepsParams3) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDepsParams3) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDepsParams3 to hex string (Canton MCMS format)
func (t SetDepsParams3) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDepsParams3 from hex string (Canton MCMS format)
func (t *SetDepsParams3) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error)
	AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error)
	CCIPSend(args CCIPSend) (*bind.EncodedChoice, error)
	CreateRouter(args CreateRouter) (*bind.EncodedChoice, error)
	Execute(args Execute) (*bind.EncodedChoice, error)
	FactorySetDeps(args FactorySetDeps) (*bind.EncodedChoice, error)
	FinalizeFee(args FinalizeFee2) (*bind.EncodedChoice, error)
	GetExecutionState(args GetExecutionState) (*bind.EncodedChoice, error)
	GetExecutionStateMCMSParams(args GetExecutionStateMCMSParams) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error)
	GetSequenceNumber(args GetSequenceNumber) (*bind.EncodedChoice, error)
	HasRouter(args HasRouter) (*bind.EncodedChoice, error)
	HasRouterMCMSParams(args HasRouterMCMSParams) (*bind.EncodedChoice, error)
	IsExecuted(args IsExecuted) (*bind.EncodedChoice, error)
	PrepareExecute(args PrepareExecute2) (*bind.EncodedChoice, error)
	PrepareExecuteMCMSParams(args PrepareExecute2MCMSParams) (*bind.EncodedChoice, error)
	PrepareSend(args PrepareSend) (*bind.EncodedChoice, error)
	RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error)
	RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error)
	SetDeps(args SetDeps3) (*bind.EncodedChoice, error)
	SetDepsParams(args SetDepsParams3) (*bind.EncodedChoice, error)
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

// AddCustomObservers encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObservers(args AddCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// AddCustomObserversParams encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObserversParams(args AddCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
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

// FactorySetDeps encodes parameters for the FactorySetDeps choice.
func (e *encoder) FactorySetDeps(args FactorySetDeps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FactorySetDeps", args)
}

// FinalizeFee encodes parameters for the FinalizeFee choice.
func (e *encoder) FinalizeFee(args FinalizeFee2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeFee", args)
}

// GetExecutionState encodes parameters for the GetExecutionState choice.
func (e *encoder) GetExecutionState(args GetExecutionState) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetExecutionState", args)
}

// GetExecutionStateMCMSParams encodes MCMS parameters (without Caller) for the GetExecutionState choice.
func (e *encoder) GetExecutionStateMCMSParams(args GetExecutionStateMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetExecutionState", args)
}

// GetFee encodes parameters for the GetFee choice.
func (e *encoder) GetFee(args GetFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFee", args)
}

// GetRequiredCCVsForExecute encodes parameters for the GetRequiredCCVsForExecute choice.
func (e *encoder) GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForExecute", args)
}

// GetRequiredCCVsForSend encodes parameters for the GetRequiredCCVsForSend choice.
func (e *encoder) GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForSend", args)
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

// IsExecuted encodes parameters for the IsExecuted choice.
func (e *encoder) IsExecuted(args IsExecuted) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("IsExecuted", args)
}

// PrepareExecute encodes parameters for the PrepareExecute choice.
func (e *encoder) PrepareExecute(args PrepareExecute2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// PrepareExecuteMCMSParams encodes MCMS parameters (without Caller) for the PrepareExecute choice.
func (e *encoder) PrepareExecuteMCMSParams(args PrepareExecute2MCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// PrepareSend encodes parameters for the PrepareSend choice.
func (e *encoder) PrepareSend(args PrepareSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareSend", args)
}

// RemoveCustomObservers encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObservers(args RemoveCustomObservers) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// RemoveCustomObserversParams encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObserversParams(args RemoveCustomObserversParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// SetDeps encodes parameters for the SetDeps choice.
func (e *encoder) SetDeps(args SetDeps3) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// SetDepsParams encodes parameters for the SetDeps choice.
func (e *encoder) SetDepsParams(args SetDepsParams3) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
