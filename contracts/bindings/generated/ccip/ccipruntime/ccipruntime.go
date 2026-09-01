package ccipruntime

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipapi"
	ccipcodec "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	clientapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	events "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/events"
	extensionapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/extensionapi"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	api "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
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
	PackageName = "ccip-runtime-v2"
	PackageID   = "e0bea36920a86d9372140068ee47b0e80cd10d6318adcebe0bec62f7d8a96b85"
	SDKVersion  = "3.4.11"
)

type Template interface {
	CreateCommand() *model.CreateCommand
	GetTemplateID() string
}

const (
	PerPartyRouterInstanceIdPrefix = types.TEXT("per-party-router")
	MaxExecutedMessagesSize        = types.INT64(25000)
	OnRampContextKey               = types.TEXT("on-ramp")
	NoExecutionAddressBytes        = types.TEXT("eba517d200000000000000000000000000000000000000000000000000000000")
	MessageStaticSize              = types.INT64(69)
	OffRampContextKey              = types.TEXT("off-ramp")
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

// AddCustomObservers2 is a Record type
type AddCustomObservers2 struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObservers2 to a map for DAML arguments
func (t AddCustomObservers2) ToMap() map[string]any {
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

func (t AddCustomObservers2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObservers2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObservers2 to hex string (Canton MCMS format)
func (t AddCustomObservers2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObservers2 from hex string (Canton MCMS format)
func (t *AddCustomObservers2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// AddCustomObserversParams2 is a Record type
type AddCustomObserversParams2 struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts AddCustomObserversParams2 to a map for DAML arguments
func (t AddCustomObserversParams2) ToMap() map[string]any {
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

func (t AddCustomObserversParams2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *AddCustomObserversParams2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes AddCustomObserversParams2 to hex string (Canton MCMS format)
func (t AddCustomObserversParams2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes AddCustomObserversParams2 from hex string (Canton MCMS format)
func (t *AddCustomObserversParams2) UnmarshalHex(data string) error {
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
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t ArchivedExecutedMessages) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages")
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
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "IsExecuted",
		Arguments:  argsToMap(args),
	}
}

// IsExecutedWithPackageID exercises the IsExecuted choice using the provided package ID instead of package name
func (t ArchivedExecutedMessages) IsExecutedWithPackageID(contractID string, packageID string, args IsExecuted) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "IsExecuted",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this ArchivedExecutedMessages contract
// This method uses the package name in the template ID
func (t ArchivedExecutedMessages) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t ArchivedExecutedMessages) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "ArchivedExecutedMessages"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// CCIPSend is a Record type
type CCIPSend struct {
	SendingMessageCid       types.CONTRACT_ID                          `json:"sendingMessageCid"`
	FeeTokenHoldingCids     []types.CONTRACT_ID                        `json:"feeTokenHoldingCids"`
	FeeTokenConfigCid       types.CONTRACT_ID                          `json:"feeTokenConfigCid"`
	FeeTokenTransferFactory types.CONTRACT_ID                          `json:"feeTokenTransferFactory"`
	FeeTokenExtraArgs       splice_api_token_metadata_v1.ExtraArgs     `json:"feeTokenExtraArgs"`
	Context                 splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts CCIPSend to a map for DAML arguments
func (t CCIPSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["feeTokenHoldingCids"] = func() []any {
		res := make([]any, 0, len(t.FeeTokenHoldingCids))
		for _, e := range t.FeeTokenHoldingCids {
			res = append(res, e)
		}
		return res
	}()

	m["feeTokenConfigCid"] = model.NestedToDAMLValue(t.FeeTokenConfigCid)

	m["feeTokenTransferFactory"] = model.NestedToDAMLValue(t.FeeTokenTransferFactory)

	m["feeTokenExtraArgs"] = model.NestedToDAMLValue(t.FeeTokenExtraArgs)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// CCIPSendFromRouter is a Record type
type CCIPSendFromRouter struct {
	RouterPartyOwner      types.PARTY        `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT         `json:"routerInstanceId"`
	CurrentSequenceNumber types.NUMERIC      `json:"currentSequenceNumber"`
	RmnRemoteCid          types.CONTRACT_ID  `json:"rmnRemoteCid"`
	GlobalConfigCid       types.CONTRACT_ID  `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        *types.CONTRACT_ID `json:"tokenConfigCid" hex:"optional"`
	SendingMessageCid     types.CONTRACT_ID  `json:"sendingMessageCid"`
}

// ToMap converts CCIPSendFromRouter to a map for DAML arguments
func (t CCIPSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["currentSequenceNumber"] = t.CurrentSequenceNumber

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	if t.TokenConfigCid != nil {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenConfigCid),
		}
	} else {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	return m
}

func (t CCIPSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSendFromRouter to hex string (Canton MCMS format)
func (t CCIPSendFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSendFromRouter from hex string (Canton MCMS format)
func (t *CCIPSendFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCIPSendFromRouterResult is a Record type
type CCIPSendFromRouterResult struct {
	CcipMessageSent      types.CONTRACT_ID `json:"ccipMessageSent"`
	VerifierBlobs        []types.TEXT      `json:"verifierBlobs"`
	MessageSentObservers []types.PARTY     `json:"messageSentObservers"`
	Receipts             []events.Receipt  `json:"receipts"`
}

// ToMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipMessageSent"] = model.NestedToDAMLValue(t.CcipMessageSent)

	m["verifierBlobs"] = func() []any {
		res := make([]any, 0, len(t.VerifierBlobs))
		for _, e := range t.VerifierBlobs {
			res = append(res, string(e))
		}
		return res
	}()

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["receipts"] = func() []any {
		res := make([]any, 0, len(t.Receipts))
		for _, e := range t.Receipts {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t CCIPSendFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSendFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSendFromRouterResult to hex string (Canton MCMS format)
func (t CCIPSendFromRouterResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSendFromRouterResult from hex string (Canton MCMS format)
func (t *CCIPSendFromRouterResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CCIPSendResult2 is a Record type
type CCIPSendResult2 struct {
	Router                 types.CONTRACT_ID   `json:"router"`
	CcipMessageSent        types.CONTRACT_ID   `json:"ccipMessageSent"`
	MessageId              types.TEXT          `json:"messageId"`
	FeeChangeCids          []types.CONTRACT_ID `json:"feeChangeCids"`
	PendingFeeInstructions []types.CONTRACT_ID `json:"pendingFeeInstructions"`
}

// ToMap converts CCIPSendResult2 to a map for DAML arguments
func (t CCIPSendResult2) ToMap() map[string]any {
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

	m["pendingFeeInstructions"] = func() []any {
		res := make([]any, 0, len(t.PendingFeeInstructions))
		for _, e := range t.PendingFeeInstructions {
			res = append(res, e)
		}
		return res
	}()

	return m
}

func (t CCIPSendResult2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPSendResult2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPSendResult2 to hex string (Canton MCMS format)
func (t CCIPSendResult2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPSendResult2 from hex string (Canton MCMS format)
func (t *CCIPSendResult2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CreateRouter is a Record type
type CreateRouter struct {
	PartyOwner          types.PARTY    `json:"partyOwner"`
	InstanceId          types.TEXT     `json:"instanceId"`
	FeeTransferLifetime *types.RELTIME `json:"feeTransferLifetime" hex:"optional"`
}

// ToMap converts CreateRouter to a map for DAML arguments
func (t CreateRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["partyOwner"] = t.PartyOwner.ToMap()

	m["instanceId"] = string(t.InstanceId)

	if t.FeeTransferLifetime != nil {
		m["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.FeeTransferLifetime),
		}
	} else {
		m["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
	ExecutingMessageCid types.CONTRACT_ID                          `json:"executingMessageCid"`
	Context             splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]any {
	m := make(map[string]any)

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// ExecuteFromRouter is a Record type
type ExecuteFromRouter struct {
	RouterPartyOwner      types.PARTY        `json:"routerPartyOwner"`
	ExecutingMessageCid   types.CONTRACT_ID  `json:"executingMessageCid"`
	GlobalConfigCid       types.CONTRACT_ID  `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID  `json:"tokenAdminRegistryCid"`
	TokenConfigCid        *types.CONTRACT_ID `json:"tokenConfigCid" hex:"optional"`
	RmnRemoteCid          types.CONTRACT_ID  `json:"rmnRemoteCid"`
}

// ToMap converts ExecuteFromRouter to a map for DAML arguments
func (t ExecuteFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["executingMessageCid"] = model.NestedToDAMLValue(t.ExecutingMessageCid)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	if t.TokenConfigCid != nil {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenConfigCid),
		}
	} else {
		m["tokenConfigCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	return m
}

func (t ExecuteFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteFromRouter to hex string (Canton MCMS format)
func (t ExecuteFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteFromRouter from hex string (Canton MCMS format)
func (t *ExecuteFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteFromRouterResult is a Record type
type ExecuteFromRouterResult struct {
	MessageId             types.TEXT          `json:"messageId"`
	Message               ccipcodec.MessageV1 `json:"message"`
	SourceChainSelector   types.NUMERIC       `json:"sourceChainSelector"`
	SequenceNumber        types.NUMERIC       `json:"sequenceNumber"`
	TokenReceiveTicket    *types.CONTRACT_ID  `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID   `json:"executionStateChanged"`
}

// ToMap converts ExecuteFromRouterResult to a map for DAML arguments
func (t ExecuteFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageId"] = string(t.MessageId)

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["sequenceNumber"] = t.SequenceNumber

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

	return m
}

func (t ExecuteFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteFromRouterResult to hex string (Canton MCMS format)
func (t ExecuteFromRouterResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteFromRouterResult from hex string (Canton MCMS format)
func (t *ExecuteFromRouterResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// ExecuteResult2 is a Record type
type ExecuteResult2 struct {
	Router                types.CONTRACT_ID             `json:"router"`
	TokenReceiveTicket    *types.CONTRACT_ID            `json:"tokenReceiveTicket" hex:"optional"`
	ExecutionStateChanged types.CONTRACT_ID             `json:"executionStateChanged"`
	MessageId             types.TEXT                    `json:"messageId"`
	Message               ccipcodec.MessageV1           `json:"message"`
	State                 ccipapi.MessageExecutionState `json:"state"`
}

// ToMap converts ExecuteResult2 to a map for DAML arguments
func (t ExecuteResult2) ToMap() map[string]any {
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

func (t ExecuteResult2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *ExecuteResult2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes ExecuteResult2 to hex string (Canton MCMS format)
func (t ExecuteResult2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes ExecuteResult2 from hex string (Canton MCMS format)
func (t *ExecuteResult2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FactorySetDeps is a Record type
type FactorySetDeps struct {
	NewDeps SetDepsParams `json:"newDeps"`
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

// FeeCalculationInputs is a Record type
type FeeCalculationInputs struct {
	FeeToken            splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
	Payload             types.TEXT                               `json:"payload"`
	ExecutorArgs        types.TEXT                               `json:"executorArgs"`
	HasTokenTransfer    types.BOOL                               `json:"hasTokenTransfer"`
	CcipReceiveGasLimit types.INT64                              `json:"ccipReceiveGasLimit"`
	CcvFeeUSDCents      types.NUMERIC                            `json:"ccvFeeUSDCents"`
	CcvGasSum           types.INT64                              `json:"ccvGasSum"`
	CcvBytesSum         types.INT64                              `json:"ccvBytesSum"`
	PoolFeeUSDCents     types.NUMERIC                            `json:"poolFeeUSDCents"`
	PoolGas             types.INT64                              `json:"poolGas"`
	PoolBytes           types.INT64                              `json:"poolBytes"`
	ExecutorFeeUSDCents types.NUMERIC                            `json:"executorFeeUSDCents"`
	HasExecutorFee      types.BOOL                               `json:"hasExecutorFee"`
}

// ToMap converts FeeCalculationInputs to a map for DAML arguments
func (t FeeCalculationInputs) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeToken"] = model.NestedToDAMLValue(t.FeeToken)

	m["payload"] = string(t.Payload)

	m["executorArgs"] = string(t.ExecutorArgs)

	m["hasTokenTransfer"] = bool(t.HasTokenTransfer)

	m["ccipReceiveGasLimit"] = int64(t.CcipReceiveGasLimit)

	m["ccvFeeUSDCents"] = t.CcvFeeUSDCents

	m["ccvGasSum"] = int64(t.CcvGasSum)

	m["ccvBytesSum"] = int64(t.CcvBytesSum)

	m["poolFeeUSDCents"] = t.PoolFeeUSDCents

	m["poolGas"] = int64(t.PoolGas)

	m["poolBytes"] = int64(t.PoolBytes)

	m["executorFeeUSDCents"] = t.ExecutorFeeUSDCents

	m["hasExecutorFee"] = bool(t.HasExecutorFee)

	return m
}

func (t FeeCalculationInputs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeCalculationInputs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeCalculationInputs to hex string (Canton MCMS format)
func (t FeeCalculationInputs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeCalculationInputs from hex string (Canton MCMS format)
func (t *FeeCalculationInputs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FeeCalculationResult is a Record type
type FeeCalculationResult struct {
	TotalFeeTokenAmount       types.NUMERIC `json:"totalFeeTokenAmount"`
	TotalUSDCents             types.NUMERIC `json:"totalUSDCents"`
	ExecutionCostUSDCents     types.NUMERIC `json:"executionCostUSDCents"`
	ExecutorDestGasLimit      types.INT64   `json:"executorDestGasLimit"`
	ExecutorDestBytesOverhead types.INT64   `json:"executorDestBytesOverhead"`
	TotalExecutionGasLimit    types.INT64   `json:"totalExecutionGasLimit"`
	PremiumMultiplier         types.NUMERIC `json:"premiumMultiplier"`
	FeeTokenPrice             types.NUMERIC `json:"feeTokenPrice"`
}

// ToMap converts FeeCalculationResult to a map for DAML arguments
func (t FeeCalculationResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["totalFeeTokenAmount"] = t.TotalFeeTokenAmount

	m["totalUSDCents"] = t.TotalUSDCents

	m["executionCostUSDCents"] = t.ExecutionCostUSDCents

	m["executorDestGasLimit"] = int64(t.ExecutorDestGasLimit)

	m["executorDestBytesOverhead"] = int64(t.ExecutorDestBytesOverhead)

	m["totalExecutionGasLimit"] = int64(t.TotalExecutionGasLimit)

	m["premiumMultiplier"] = t.PremiumMultiplier

	m["feeTokenPrice"] = t.FeeTokenPrice

	return m
}

func (t FeeCalculationResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FeeCalculationResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FeeCalculationResult to hex string (Canton MCMS format)
func (t FeeCalculationResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FeeCalculationResult from hex string (Canton MCMS format)
func (t *FeeCalculationResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FinalizeFee2 is a Record type
type FinalizeFee2 struct {
	SendingMessageCid types.CONTRACT_ID                          `json:"sendingMessageCid"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts FinalizeFee2 to a map for DAML arguments
func (t FinalizeFee2) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// FinalizeFeeFromRouter is a Record type
type FinalizeFeeFromRouter struct {
	RouterPartyOwner  types.PARTY       `json:"routerPartyOwner"`
	RouterInstanceId  types.TEXT        `json:"routerInstanceId"`
	GlobalConfigCid   types.CONTRACT_ID `json:"globalConfigCid"`
	FeeQuoterCid      types.CONTRACT_ID `json:"feeQuoterCid"`
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
}

// ToMap converts FinalizeFeeFromRouter to a map for DAML arguments
func (t FinalizeFeeFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["sendingMessageCid"] = model.NestedToDAMLValue(t.SendingMessageCid)

	return m
}

func (t FinalizeFeeFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FinalizeFeeFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FinalizeFeeFromRouter to hex string (Canton MCMS format)
func (t FinalizeFeeFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FinalizeFeeFromRouter from hex string (Canton MCMS format)
func (t *FinalizeFeeFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetExecutionState is a Record type
type GetExecutionState struct {
	MessageHash types.TEXT                                 `json:"messageHash"`
	Context     splice_api_token_metadata_v1.ChoiceContext `json:"context"`
	Caller      types.PARTY                                `json:"caller"`
}

// ToMap converts GetExecutionState to a map for DAML arguments
func (t GetExecutionState) ToMap() map[string]any {
	m := make(map[string]any)

	m["messageHash"] = string(t.MessageHash)

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type GetExecutionStateMCMSParams struct {
	MessageHash types.TEXT                                 `json:"messageHash"`
	Context     splice_api_token_metadata_v1.ChoiceContext `json:"context"`
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
	Message           clientapi.Canton2AnyMessage                `json:"message"`
	CcvFeeQuotes      []extensionapi.CrossChainVerifierFeeQuote  `json:"ccvFeeQuotes"`
	TokenPoolFeeQuote *extensionapi.TokenPoolFeeQuote            `json:"tokenPoolFeeQuote" hex:"optional"`
	ExecutorFeeQuote  *extensionapi.ExecutorFeeQuote             `json:"executorFeeQuote" hex:"optional"`
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

// GetFeeFromRouter is a Record type
type GetFeeFromRouter struct {
	RouterPartyOwner  types.PARTY                               `json:"routerPartyOwner"`
	RouterInstanceId  types.TEXT                                `json:"routerInstanceId"`
	GlobalConfigCid   types.CONTRACT_ID                         `json:"globalConfigCid"`
	FeeQuoterCid      types.CONTRACT_ID                         `json:"feeQuoterCid"`
	DestChainSelector types.NUMERIC                             `json:"destChainSelector"`
	Message           clientapi.Canton2AnyMessage               `json:"message"`
	CcvFeeQuotes      []extensionapi.CrossChainVerifierFeeQuote `json:"ccvFeeQuotes"`
	TokenPoolFeeQuote *extensionapi.TokenPoolFeeQuote           `json:"tokenPoolFeeQuote" hex:"optional"`
	ExecutorFeeQuote  *extensionapi.ExecutorFeeQuote            `json:"executorFeeQuote" hex:"optional"`
}

// ToMap converts GetFeeFromRouter to a map for DAML arguments
func (t GetFeeFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

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

func (t GetFeeFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeFromRouter to hex string (Canton MCMS format)
func (t GetFeeFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeFromRouter from hex string (Canton MCMS format)
func (t *GetFeeFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFeeFromRouterResult is a Record type
type GetFeeFromRouterResult struct {
	FeeTokenAmount types.NUMERIC `json:"feeTokenAmount"`
}

// ToMap converts GetFeeFromRouterResult to a map for DAML arguments
func (t GetFeeFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeTokenAmount"] = t.FeeTokenAmount

	return m
}

func (t GetFeeFromRouterResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFeeFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFeeFromRouterResult to hex string (Canton MCMS format)
func (t GetFeeFromRouterResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFeeFromRouterResult from hex string (Canton MCMS format)
func (t *GetFeeFromRouterResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	Message                   ccipcodec.MessageV1                        `json:"message"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress          `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                                `json:"receiverOptionalThreshold"`
	TokenPoolRequiredCCVs     []chainlinkapi.RawInstanceAddress          `json:"tokenPoolRequiredCCVs"`
	Context                   splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) ToMap() map[string]any {
	m := make(map[string]any)

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

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// GetRequiredCCVsForExecuteFromRouter is a Record type
type GetRequiredCCVsForExecuteFromRouter struct {
	GlobalConfigCid           types.CONTRACT_ID                 `json:"globalConfigCid"`
	Message                   ccipcodec.MessageV1               `json:"message"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                       `json:"receiverOptionalThreshold"`
	TokenPoolRequiredCCVs     []chainlinkapi.RawInstanceAddress `json:"tokenPoolRequiredCCVs"`
}

// ToMap converts GetRequiredCCVsForExecuteFromRouter to a map for DAML arguments
func (t GetRequiredCCVsForExecuteFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

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

func (t GetRequiredCCVsForExecuteFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForExecuteFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForExecuteFromRouter to hex string (Canton MCMS format)
func (t GetRequiredCCVsForExecuteFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForExecuteFromRouter from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForExecuteFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Message           clientapi.Canton2AnyMessage                `json:"message"`
	PoolReportedCCVs  []chainlinkapi.RawInstanceAddress          `json:"poolReportedCCVs"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["poolReportedCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolReportedCCVs))
		for _, e := range t.PoolReportedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	m["context"] = model.NestedToDAMLValue(t.Context)

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

// GetRequiredCCVsForSendFromRouter is a Record type
type GetRequiredCCVsForSendFromRouter struct {
	GlobalConfigCid   types.CONTRACT_ID                 `json:"globalConfigCid"`
	DestChainSelector types.NUMERIC                     `json:"destChainSelector"`
	Message           clientapi.Canton2AnyMessage       `json:"message"`
	PoolReportedCCVs  []chainlinkapi.RawInstanceAddress `json:"poolReportedCCVs"`
}

// ToMap converts GetRequiredCCVsForSendFromRouter to a map for DAML arguments
func (t GetRequiredCCVsForSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["poolReportedCCVs"] = func() []any {
		res := make([]any, 0, len(t.PoolReportedCCVs))
		for _, e := range t.PoolReportedCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	return m
}

func (t GetRequiredCCVsForSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetRequiredCCVsForSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetRequiredCCVsForSendFromRouter to hex string (Canton MCMS format)
func (t GetRequiredCCVsForSendFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetRequiredCCVsForSendFromRouter from hex string (Canton MCMS format)
func (t *GetRequiredCCVsForSendFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetSequenceNumber is a Record type
type GetSequenceNumber struct {
	DestChainSelector types.NUMERIC                              `json:"destChainSelector"`
	Context           splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetSequenceNumber to a map for DAML arguments
func (t GetSequenceNumber) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["context"] = model.NestedToDAMLValue(t.Context)

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
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
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

// OffRamp is a Template type
type OffRamp struct {
	InstanceId types.TEXT  `json:"instanceId"`
	CcipOwner  types.PARTY `json:"ccipOwner"`
	Deps       OffRampDeps `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OffRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "OffRamp")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OffRamp) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "OffRamp")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OffRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t OffRamp) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t OffRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OffRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OffRamp to hex string (Canton MCMS format)
func (t OffRamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OffRamp from hex string (Canton MCMS format)
func (t *OffRamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for OffRamp

// ExecuteFromRouter exercises the ExecuteFromRouter choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) ExecuteFromRouter(contractID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// ExecuteFromRouterWithPackageID exercises the ExecuteFromRouter choice using the provided package ID instead of package name
func (t OffRamp) ExecuteFromRouterWithPackageID(contractID string, packageID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) PrepareExecute(contractID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t OffRamp) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteFromRouter exercises the GetRequiredCCVsForExecuteFromRouter choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) GetRequiredCCVsForExecuteFromRouter(contractID string, args GetRequiredCCVsForExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteFromRouterWithPackageID exercises the GetRequiredCCVsForExecuteFromRouter choice using the provided package ID instead of package name
func (t OffRamp) GetRequiredCCVsForExecuteFromRouterWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OffRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OffRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OffRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetDeps exercises the SetDeps choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t OffRamp) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this OffRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OffRamp) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t OffRamp) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OffRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for OffRamp

var _ api.IMCMSReceiver = (*OffRamp)(nil)

// OffRampDeps is a Record type
type OffRampDeps struct {
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts OffRampDeps to a map for DAML arguments
func (t OffRampDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	return m
}

func (t OffRampDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OffRampDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OffRampDeps to hex string (Canton MCMS format)
func (t OffRampDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OffRampDeps from hex string (Canton MCMS format)
func (t *OffRampDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// OnRamp is a Template type
type OnRamp struct {
	InstanceId        types.TEXT    `json:"instanceId"`
	CcipOwner         types.PARTY   `json:"ccipOwner"`
	MaxUSDCentsPerMsg types.NUMERIC `json:"maxUSDCentsPerMsg"`
	Deps              OnRampDeps    `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OnRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OnRamp) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OnRamp) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	if t.MaxUSDCentsPerMsg != "" {
		args["maxUSDCentsPerMsg"] = t.MaxUSDCentsPerMsg
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t OnRamp) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	if t.MaxUSDCentsPerMsg != "" {
		args["maxUSDCentsPerMsg"] = t.MaxUSDCentsPerMsg
	}

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = model.NestedToDAMLValue(t.Deps)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t OnRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OnRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OnRamp to hex string (Canton MCMS format)
func (t OnRamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OnRamp from hex string (Canton MCMS format)
func (t *OnRamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for OnRamp

// CCIPSendFromRouter exercises the CCIPSendFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) CCIPSendFromRouter(contractID string, args CCIPSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CCIPSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// CCIPSendFromRouterWithPackageID exercises the CCIPSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) CCIPSendFromRouterWithPackageID(contractID string, packageID string, args CCIPSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CCIPSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// GetFeeFromRouter exercises the GetFeeFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) GetFeeFromRouter(contractID string, args GetFeeFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetFeeFromRouter",
		Arguments:  argsToMap(args),
	}
}

// GetFeeFromRouterWithPackageID exercises the GetFeeFromRouter choice using the provided package ID instead of package name
func (t OnRamp) GetFeeFromRouterWithPackageID(contractID string, packageID string, args GetFeeFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetFeeFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendFromRouter exercises the PrepareSendFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) PrepareSendFromRouter(contractID string, args PrepareSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "PrepareSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendFromRouterWithPackageID exercises the PrepareSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) PrepareSendFromRouterWithPackageID(contractID string, packageID string, args PrepareSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "PrepareSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendFromRouter exercises the GetRequiredCCVsForSendFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) GetRequiredCCVsForSendFromRouter(contractID string, args GetRequiredCCVsForSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendFromRouterWithPackageID exercises the GetRequiredCCVsForSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) GetRequiredCCVsForSendFromRouterWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFeeFromRouter exercises the FinalizeFeeFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) FinalizeFeeFromRouter(contractID string, args FinalizeFeeFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "FinalizeFeeFromRouter",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFeeFromRouterWithPackageID exercises the FinalizeFeeFromRouter choice using the provided package ID instead of package name
func (t OnRamp) FinalizeFeeFromRouterWithPackageID(contractID string, packageID string, args FinalizeFeeFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "FinalizeFeeFromRouter",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OnRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OnRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OnRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// SetDeps exercises the SetDeps choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t OnRamp) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this OnRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OnRamp) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t OnRamp) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for OnRamp

var _ api.IMCMSReceiver = (*OnRamp)(nil)

// OnRampDeps is a Record type
type OnRampDeps struct {
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
}

// ToMap converts OnRampDeps to a map for DAML arguments
func (t OnRampDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	return m
}

func (t OnRampDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *OnRampDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes OnRampDeps to hex string (Canton MCMS format)
func (t OnRampDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes OnRampDeps from hex string (Canton MCMS format)
func (t *OnRampDeps) UnmarshalHex(data string) error {
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
	FeeTransferLifetime          *types.RELTIME                  `json:"feeTransferLifetime" hex:"optional"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t PerPartyRouter) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouter) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter")
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

	if t.FeeTransferLifetime != nil {
		args["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.FeeTransferLifetime),
		}
	} else {
		args["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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

	if t.FeeTransferLifetime != nil {
		args["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.FeeTransferLifetime),
		}
	} else {
		args["feeTransferLifetime"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// CCIPSendWithPackageID exercises the CCIPSend choice using the provided package ID instead of package name
func (t PerPartyRouter) CCIPSendWithPackageID(contractID string, packageID string, args CCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// Execute exercises the Execute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t PerPartyRouter) ExecuteWithPackageID(contractID string, packageID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionState exercises the GetExecutionState choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetExecutionState(contractID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// GetExecutionStateWithPackageID exercises the GetExecutionState choice using the provided package ID instead of package name
func (t PerPartyRouter) GetExecutionStateWithPackageID(contractID string, packageID string, args GetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareExecute(contractID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// GetFee exercises the GetFee choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t PerPartyRouter) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// PrepareSend exercises the PrepareSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) PrepareSend(contractID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendWithPackageID exercises the PrepareSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PrepareSendWithPackageID(contractID string, packageID string, args PrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFee exercises the FinalizeFee choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) FinalizeFee(contractID string, args FinalizeFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// FinalizeFeeWithPackageID exercises the FinalizeFee choice using the provided package ID instead of package name
func (t PerPartyRouter) FinalizeFeeWithPackageID(contractID string, packageID string, args FinalizeFee2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouter) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetSequenceNumber exercises the GetSequenceNumber choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) GetSequenceNumber(contractID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// GetSequenceNumberWithPackageID exercises the GetSequenceNumber choice using the provided package ID instead of package name
func (t PerPartyRouter) GetSequenceNumberWithPackageID(contractID string, packageID string, args GetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// SetDeps exercises the SetDeps choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t PerPartyRouter) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObservers exercises the AddCustomObservers choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) AddCustomObservers(contractID string, args AddCustomObservers2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// AddCustomObserversWithPackageID exercises the AddCustomObservers choice using the provided package ID instead of package name
func (t PerPartyRouter) AddCustomObserversWithPackageID(contractID string, packageID string, args AddCustomObservers2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "AddCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObservers exercises the RemoveCustomObservers choice on this PerPartyRouter contract
// This method uses the package name in the template ID
func (t PerPartyRouter) RemoveCustomObservers(contractID string, args RemoveCustomObservers2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// RemoveCustomObserversWithPackageID exercises the RemoveCustomObservers choice using the provided package ID instead of package name
func (t PerPartyRouter) RemoveCustomObserversWithPackageID(contractID string, packageID string, args RemoveCustomObservers2) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "RemoveCustomObservers",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this PerPartyRouter contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouter) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t PerPartyRouter) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetSequenceNumber exercises the PerPartyRouter_GetSequenceNumber choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterGetSequenceNumber(contractID string, args clientapi.PerPartyRouterGetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetSequenceNumberWithPackageID exercises the PerPartyRouter_GetSequenceNumber choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterGetSequenceNumberWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterGetSequenceNumber) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetSequenceNumber",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetRequiredCCVsForSend exercises the PerPartyRouter_GetRequiredCCVsForSend choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterGetRequiredCCVsForSend(contractID string, args clientapi.PerPartyRouterGetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetRequiredCCVsForSendWithPackageID exercises the PerPartyRouter_GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterGetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterGetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetFee exercises the PerPartyRouter_GetFee choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterGetFee(contractID string, args clientapi.PerPartyRouterGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetFee",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetFeeWithPackageID exercises the PerPartyRouter_GetFee choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterGetFeeWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetFee",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterPrepareSend exercises the PerPartyRouter_PrepareSend choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterPrepareSend(contractID string, args clientapi.PerPartyRouterPrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterPrepareSendWithPackageID exercises the PerPartyRouter_PrepareSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterPrepareSendWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterPrepareSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_PrepareSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterFinalizeFee exercises the PerPartyRouter_FinalizeFee choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterFinalizeFee(contractID string, args clientapi.PerPartyRouterFinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterFinalizeFeeWithPackageID exercises the PerPartyRouter_FinalizeFee choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterFinalizeFeeWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterFinalizeFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_FinalizeFee",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterCCIPSend exercises the PerPartyRouter_CCIPSend choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterCCIPSend(contractID string, args clientapi.PerPartyRouterCCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterCCIPSendWithPackageID exercises the PerPartyRouter_CCIPSend choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterCCIPSendWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterCCIPSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_CCIPSend",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetExecutionState exercises the PerPartyRouter_GetExecutionState choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterGetExecutionState(contractID string, args clientapi.PerPartyRouterGetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetExecutionStateWithPackageID exercises the PerPartyRouter_GetExecutionState choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterGetExecutionStateWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterGetExecutionState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetExecutionState",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetRequiredCCVsForExecute exercises the PerPartyRouter_GetRequiredCCVsForExecute choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterGetRequiredCCVsForExecute(contractID string, args clientapi.PerPartyRouterGetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterGetRequiredCCVsForExecuteWithPackageID exercises the PerPartyRouter_GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterGetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterGetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterPrepareExecute exercises the PerPartyRouter_PrepareExecute choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterPrepareExecute(contractID string, args clientapi.PerPartyRouterPrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterPrepareExecuteWithPackageID exercises the PerPartyRouter_PrepareExecute choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterPrepareExecuteWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterPrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterExecute exercises the PerPartyRouter_Execute choice on this PerPartyRouter contract via the IPerPartyRouter interface
// This method uses the package name in the template ID
func (t PerPartyRouter) PerPartyRouterExecute(contractID string, args clientapi.PerPartyRouterExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_Execute",
		Arguments:  argsToMap(args),
	}
}

// PerPartyRouterExecuteWithPackageID exercises the PerPartyRouter_Execute choice using the provided package ID instead of package name
func (t PerPartyRouter) PerPartyRouterExecuteWithPackageID(contractID string, packageID string, args clientapi.PerPartyRouterExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouter"),
		ContractID: contractID,
		Choice:     "PerPartyRouter_Execute",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for PerPartyRouter

var _ api.IMCMSReceiver = (*PerPartyRouter)(nil)

var _ clientapi.IPerPartyRouter = (*PerPartyRouter)(nil)

// PerPartyRouterDeps is a Record type
type PerPartyRouterDeps struct {
	OnRamp             chainlinkapi.RawInstanceAddress `json:"onRamp"`
	OffRamp            chainlinkapi.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
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
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t PerPartyRouterFactory) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory")
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
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// CreateRouterWithPackageID exercises the CreateRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) CreateRouterWithPackageID(contractID string, packageID string, args CreateRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "CreateRouter",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this PerPartyRouterFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// HasRouter exercises the HasRouter choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) HasRouter(contractID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// HasRouterWithPackageID exercises the HasRouter choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) HasRouterWithPackageID(contractID string, packageID string, args HasRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "HasRouter",
		Arguments:  argsToMap(args),
	}
}

// FactorySetDeps exercises the FactorySetDeps choice on this PerPartyRouterFactory contract
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) FactorySetDeps(contractID string, args FactorySetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "FactorySetDeps",
		Arguments:  argsToMap(args),
	}
}

// FactorySetDepsWithPackageID exercises the FactorySetDeps choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) FactorySetDepsWithPackageID(contractID string, packageID string, args FactorySetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "PerPartyRouterFactory"),
		ContractID: contractID,
		Choice:     "FactorySetDeps",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this PerPartyRouterFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t PerPartyRouterFactory) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t PerPartyRouterFactory) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.RuntimeV2.PerPartyRouter", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for PerPartyRouterFactory

var _ api.IMCMSReceiver = (*PerPartyRouterFactory)(nil)

// PrepareExecute is a Record type
type PrepareExecute struct {
	EncodedMessage            types.TEXT                        `json:"encodedMessage"`
	GlobalConfigCid           types.CONTRACT_ID                 `json:"globalConfigCid"`
	TokenAdminRegistryCid     types.CONTRACT_ID                 `json:"tokenAdminRegistryCid"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                       `json:"receiverOptionalThreshold"`
	ReceiverFinalityConfig    ccipcodec.FinalityConfig          `json:"receiverFinalityConfig"`
	RmnRemoteCid              types.CONTRACT_ID                 `json:"rmnRemoteCid"`
	ReceiverParty             types.PARTY                       `json:"receiverParty"`
	TokenReceiverParty        *types.PARTY                      `json:"tokenReceiverParty" hex:"optional"`
	Caller                    types.PARTY                       `json:"caller"`
}

// ToMap converts PrepareExecute to a map for DAML arguments
func (t PrepareExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["encodedMessage"] = string(t.EncodedMessage)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

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

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

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

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PrepareExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PrepareExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PrepareExecute to hex string (Canton MCMS format)
func (t PrepareExecute) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecute from hex string (Canton MCMS format)
func (t *PrepareExecute) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareExecuteMCMSParams is PrepareExecute without the Caller field for MCMS operationData encoding.
// ContractId fields are omitted; pass them via the MCMS targetCids map at execution time.
type PrepareExecuteMCMSParams struct {
	EncodedMessage            types.TEXT                        `json:"encodedMessage"`
	ReceiverRequiredCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverRequiredCCVs"`
	ReceiverOptionalCCVs      []chainlinkapi.RawInstanceAddress `json:"receiverOptionalCCVs"`
	ReceiverOptionalThreshold types.INT64                       `json:"receiverOptionalThreshold"`
	ReceiverFinalityConfig    ccipcodec.FinalityConfig          `json:"receiverFinalityConfig"`
	ReceiverParty             types.PARTY                       `json:"receiverParty"`
	TokenReceiverParty        *types.PARTY                      `json:"tokenReceiverParty" hex:"optional"`
}

// MarshalHex encodes PrepareExecuteMCMSParams to hex string for MCMS operationData.
func (t PrepareExecuteMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareExecuteMCMSParams from hex string.
func (t *PrepareExecuteMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// PrepareSend is a Record type
type PrepareSend struct {
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
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

// PrepareSendFromRouter is a Record type
type PrepareSendFromRouter struct {
	DestChainSelector     types.NUMERIC                              `json:"destChainSelector"`
	Message               clientapi.Canton2AnyMessage                `json:"message"`
	RouterPartyOwner      types.PARTY                                `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT                                 `json:"routerInstanceId"`
	GlobalConfigCid       types.CONTRACT_ID                          `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID                          `json:"tokenAdminRegistryCid"`
	FeeQuoterCid          types.CONTRACT_ID                          `json:"feeQuoterCid"`
	RmnRemoteCid          types.CONTRACT_ID                          `json:"rmnRemoteCid"`
	CurrentSequenceNumber types.NUMERIC                              `json:"currentSequenceNumber"`
	Context               splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts PrepareSendFromRouter to a map for DAML arguments
func (t PrepareSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["destChainSelector"] = t.DestChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["globalConfigCid"] = model.NestedToDAMLValue(t.GlobalConfigCid)

	m["tokenAdminRegistryCid"] = model.NestedToDAMLValue(t.TokenAdminRegistryCid)

	m["feeQuoterCid"] = model.NestedToDAMLValue(t.FeeQuoterCid)

	m["rmnRemoteCid"] = model.NestedToDAMLValue(t.RmnRemoteCid)

	m["currentSequenceNumber"] = t.CurrentSequenceNumber

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t PrepareSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *PrepareSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes PrepareSendFromRouter to hex string (Canton MCMS format)
func (t PrepareSendFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes PrepareSendFromRouter from hex string (Canton MCMS format)
func (t *PrepareSendFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveCustomObservers2 is a Record type
type RemoveCustomObservers2 struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObservers2 to a map for DAML arguments
func (t RemoveCustomObservers2) ToMap() map[string]any {
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

func (t RemoveCustomObservers2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObservers2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObservers2 to hex string (Canton MCMS format)
func (t RemoveCustomObservers2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObservers2 from hex string (Canton MCMS format)
func (t *RemoveCustomObservers2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// RemoveCustomObserversParams2 is a Record type
type RemoveCustomObserversParams2 struct {
	Parties []types.PARTY `json:"parties"`
}

// ToMap converts RemoveCustomObserversParams2 to a map for DAML arguments
func (t RemoveCustomObserversParams2) ToMap() map[string]any {
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

func (t RemoveCustomObserversParams2) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *RemoveCustomObserversParams2) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes RemoveCustomObserversParams2 to hex string (Canton MCMS format)
func (t RemoveCustomObserversParams2) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes RemoveCustomObserversParams2 from hex string (Canton MCMS format)
func (t *RemoveCustomObserversParams2) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDeps is a Record type
type SetDeps struct {
	NewDeps SetDepsParams `json:"newDeps"`
}

// ToMap converts SetDeps to a map for DAML arguments
func (t SetDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = model.NestedToDAMLValue(t.NewDeps)

	return m
}

func (t SetDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDeps to hex string (Canton MCMS format)
func (t SetDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDeps from hex string (Canton MCMS format)
func (t *SetDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetDepsParams is a Record type
type SetDepsParams struct {
	OnRamp             *chainlinkapi.RawInstanceAddress `json:"onRamp" hex:"optional"`
	OffRamp            *chainlinkapi.RawInstanceAddress `json:"offRamp" hex:"optional"`
	GlobalConfig       *chainlinkapi.RawInstanceAddress `json:"globalConfig" hex:"optional"`
	TokenAdminRegistry *chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry" hex:"optional"`
	FeeQuoter          *chainlinkapi.RawInstanceAddress `json:"feeQuoter" hex:"optional"`
	RmnRemote          *chainlinkapi.RawInstanceAddress `json:"rmnRemote" hex:"optional"`
}

// ToMap converts SetDepsParams to a map for DAML arguments
func (t SetDepsParams) ToMap() map[string]any {
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

func (t SetDepsParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetDepsParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetDepsParams to hex string (Canton MCMS format)
func (t SetDepsParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetDepsParams from hex string (Canton MCMS format)
func (t *SetDepsParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	AddCustomObservers(args AddCustomObservers2) (*bind.EncodedChoice, error)
	AddCustomObserversParams(args AddCustomObserversParams2) (*bind.EncodedChoice, error)
	CCIPSend(args CCIPSend) (*bind.EncodedChoice, error)
	CCIPSendFromRouter(args CCIPSendFromRouter) (*bind.EncodedChoice, error)
	CreateRouter(args CreateRouter) (*bind.EncodedChoice, error)
	Execute(args Execute) (*bind.EncodedChoice, error)
	ExecuteFromRouter(args ExecuteFromRouter) (*bind.EncodedChoice, error)
	FactorySetDeps(args FactorySetDeps) (*bind.EncodedChoice, error)
	FinalizeFee(args FinalizeFee2) (*bind.EncodedChoice, error)
	FinalizeFeeFromRouter(args FinalizeFeeFromRouter) (*bind.EncodedChoice, error)
	GetExecutionState(args GetExecutionState) (*bind.EncodedChoice, error)
	GetExecutionStateMCMSParams(args GetExecutionStateMCMSParams) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetFeeFromRouter(args GetFeeFromRouter) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecuteFromRouter(args GetRequiredCCVsForExecuteFromRouter) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSendFromRouter(args GetRequiredCCVsForSendFromRouter) (*bind.EncodedChoice, error)
	GetSequenceNumber(args GetSequenceNumber) (*bind.EncodedChoice, error)
	HasRouter(args HasRouter) (*bind.EncodedChoice, error)
	HasRouterMCMSParams(args HasRouterMCMSParams) (*bind.EncodedChoice, error)
	IsExecuted(args IsExecuted) (*bind.EncodedChoice, error)
	PrepareExecute(args PrepareExecute) (*bind.EncodedChoice, error)
	PrepareExecuteMCMSParams(args PrepareExecuteMCMSParams) (*bind.EncodedChoice, error)
	PrepareSend(args PrepareSend) (*bind.EncodedChoice, error)
	PrepareSendFromRouter(args PrepareSendFromRouter) (*bind.EncodedChoice, error)
	RemoveCustomObservers(args RemoveCustomObservers2) (*bind.EncodedChoice, error)
	RemoveCustomObserversParams(args RemoveCustomObserversParams2) (*bind.EncodedChoice, error)
	SetDeps(args SetDeps) (*bind.EncodedChoice, error)
	SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error)
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
func (e *encoder) AddCustomObservers(args AddCustomObservers2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// AddCustomObserversParams encodes parameters for the AddCustomObservers choice.
func (e *encoder) AddCustomObserversParams(args AddCustomObserversParams2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("AddCustomObservers", args)
}

// CCIPSend encodes parameters for the CCIPSend choice.
func (e *encoder) CCIPSend(args CCIPSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CCIPSend", args)
}

// CCIPSendFromRouter encodes parameters for the CCIPSendFromRouter choice.
func (e *encoder) CCIPSendFromRouter(args CCIPSendFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CCIPSendFromRouter", args)
}

// CreateRouter encodes parameters for the CreateRouter choice.
func (e *encoder) CreateRouter(args CreateRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CreateRouter", args)
}

// Execute encodes parameters for the Execute choice.
func (e *encoder) Execute(args Execute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Execute", args)
}

// ExecuteFromRouter encodes parameters for the ExecuteFromRouter choice.
func (e *encoder) ExecuteFromRouter(args ExecuteFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("ExecuteFromRouter", args)
}

// FactorySetDeps encodes parameters for the FactorySetDeps choice.
func (e *encoder) FactorySetDeps(args FactorySetDeps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FactorySetDeps", args)
}

// FinalizeFee encodes parameters for the FinalizeFee choice.
func (e *encoder) FinalizeFee(args FinalizeFee2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeFee", args)
}

// FinalizeFeeFromRouter encodes parameters for the FinalizeFeeFromRouter choice.
func (e *encoder) FinalizeFeeFromRouter(args FinalizeFeeFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("FinalizeFeeFromRouter", args)
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

// GetFeeFromRouter encodes parameters for the GetFeeFromRouter choice.
func (e *encoder) GetFeeFromRouter(args GetFeeFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFeeFromRouter", args)
}

// GetRequiredCCVsForExecute encodes parameters for the GetRequiredCCVsForExecute choice.
func (e *encoder) GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForExecute", args)
}

// GetRequiredCCVsForExecuteFromRouter encodes parameters for the GetRequiredCCVsForExecuteFromRouter choice.
func (e *encoder) GetRequiredCCVsForExecuteFromRouter(args GetRequiredCCVsForExecuteFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForExecuteFromRouter", args)
}

// GetRequiredCCVsForSend encodes parameters for the GetRequiredCCVsForSend choice.
func (e *encoder) GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForSend", args)
}

// GetRequiredCCVsForSendFromRouter encodes parameters for the GetRequiredCCVsForSendFromRouter choice.
func (e *encoder) GetRequiredCCVsForSendFromRouter(args GetRequiredCCVsForSendFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForSendFromRouter", args)
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
func (e *encoder) PrepareExecute(args PrepareExecute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// PrepareExecuteMCMSParams encodes MCMS parameters (without Caller) for the PrepareExecute choice.
func (e *encoder) PrepareExecuteMCMSParams(args PrepareExecuteMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareExecute", args)
}

// PrepareSend encodes parameters for the PrepareSend choice.
func (e *encoder) PrepareSend(args PrepareSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareSend", args)
}

// PrepareSendFromRouter encodes parameters for the PrepareSendFromRouter choice.
func (e *encoder) PrepareSendFromRouter(args PrepareSendFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareSendFromRouter", args)
}

// RemoveCustomObservers encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObservers(args RemoveCustomObservers2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// RemoveCustomObserversParams encodes parameters for the RemoveCustomObservers choice.
func (e *encoder) RemoveCustomObserversParams(args RemoveCustomObserversParams2) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("RemoveCustomObservers", args)
}

// SetDeps encodes parameters for the SetDeps choice.
func (e *encoder) SetDeps(args SetDeps) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// SetDepsParams encodes parameters for the SetDeps choice.
func (e *encoder) SetDepsParams(args SetDepsParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetDeps", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
