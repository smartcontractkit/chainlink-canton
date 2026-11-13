package globaltokensenderreceiver

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipcodec "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	clientapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	user "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/user"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
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
	PackageName = "ccip-global-token-sender-receiver-v1"
	PackageID   = "548a378d2a39a45df4ca843672201c2b1a809f7a77f20d7caa668c0e7807dee5"
	SDKVersion  = "3.4.11"
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

// Execute is a Record type
type Execute struct {
	ReceiverParty  types.PARTY                                `json:"receiverParty"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenTransfer  *user.ReceiverTokenTransferInput           `json:"tokenTransfer" hex:"optional"`
	CcvInputs      []user.ReceiverCCVInput                    `json:"ccvInputs"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts Execute to a map for DAML arguments
func (t Execute) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenTransfer != nil {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransfer),
		}
	} else {
		m["tokenTransfer"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["ccvInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvInputs))
		for _, e := range t.CcvInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

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

// GetFee is a Record type
type GetFee struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	CcvSendInputs            []user.SenderCCVSendInput                  `json:"ccvSendInputs"`
	TokenTransferInput       *user.SenderTokenTransferInput             `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *user.SenderExecutorInput                  `json:"executorInput" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetFee to a map for DAML arguments
func (t GetFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["ccvSendInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvSendInputs))
		for _, e := range t.CcvSendInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenTransferInput != nil {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransferInput),
		}
	} else {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorInput != nil {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorInput),
		}
	} else {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

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
	ReceiverParty  types.PARTY                                `json:"receiverParty"`
	RouterCid      types.CONTRACT_ID                          `json:"routerCid"`
	EncodedMessage types.TEXT                                 `json:"encodedMessage"`
	TokenPoolCid   *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
	Context        splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetRequiredCCVsForExecute to a map for DAML arguments
func (t GetRequiredCCVsForExecute) ToMap() map[string]any {
	m := make(map[string]any)

	m["receiverParty"] = t.ReceiverParty.ToMap()

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["encodedMessage"] = string(t.EncodedMessage)

	if t.TokenPoolCid != nil {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolCid),
		}
	} else {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	TokenPoolCid             *types.CONTRACT_ID                         `json:"tokenPoolCid" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	if t.TokenPoolCid != nil {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenPoolCid),
		}
	} else {
		m["tokenPoolCid"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

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

// GlobalTokenSenderReceiver is a Template type
type GlobalTokenSenderReceiver struct {
	InstanceId             types.TEXT                        `json:"instanceId"`
	CcipOwner              types.PARTY                       `json:"ccipOwner"`
	RequiredCCVs           []chainlinkapi.RawInstanceAddress `json:"requiredCCVs"`
	OptionalCCVs           []chainlinkapi.RawInstanceAddress `json:"optionalCCVs"`
	OptionalThreshold      types.INT64                       `json:"optionalThreshold"`
	ReceiverFinalityConfig ccipcodec.FinalityConfig          `json:"receiverFinalityConfig"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t GlobalTokenSenderReceiver) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t GlobalTokenSenderReceiver) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["optionalThreshold"] = int64(t.OptionalThreshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["receiverFinalityConfig"] = model.NestedToDAMLValue(t.ReceiverFinalityConfig)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t GlobalTokenSenderReceiver) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GlobalTokenSenderReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GlobalTokenSenderReceiver to hex string (Canton MCMS format)
func (t GlobalTokenSenderReceiver) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GlobalTokenSenderReceiver from hex string (Canton MCMS format)
func (t *GlobalTokenSenderReceiver) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for GlobalTokenSenderReceiver

// Send exercises the Send choice on this GlobalTokenSenderReceiver contract
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) Send(contractID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// SendWithPackageID exercises the Send choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) SendWithPackageID(contractID string, packageID string, args Send) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "Send",
		Arguments:  argsToMap(args),
	}
}

// GetFee exercises the GetFee choice on this GlobalTokenSenderReceiver contract
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) GetFee(contractID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetFeeWithPackageID exercises the GetFee choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) GetFeeWithPackageID(contractID string, packageID string, args GetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetFee",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this GlobalTokenSenderReceiver contract
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this GlobalTokenSenderReceiver contract
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this GlobalTokenSenderReceiver contract via the IICCIPSender interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "CCIPSender"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// Execute exercises the Execute choice on this GlobalTokenSenderReceiver contract
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) Execute(contractID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ExecuteWithPackageID exercises the Execute choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ExecuteWithPackageID(contractID string, packageID string, args Execute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "GlobalTokenSenderReceiver"),
		ContractID: contractID,
		Choice:     "Execute",
		Arguments:  argsToMap(args),
	}
}

// ICCIPReceiverGetRequiredCCVs exercises the ICCIPReceiver_GetRequiredCCVs choice on this GlobalTokenSenderReceiver contract via the IICCIPReceiver interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) ICCIPReceiverGetRequiredCCVs(contractID string, args user.ICCIPReceiverGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "ICCIPReceiver"),
		ContractID: contractID,
		Choice:     "ICCIPReceiver_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// ICCIPReceiverGetRequiredCCVsWithPackageID exercises the ICCIPReceiver_GetRequiredCCVs choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ICCIPReceiverGetRequiredCCVsWithPackageID(contractID string, packageID string, args user.ICCIPReceiverGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "ICCIPReceiver"),
		ContractID: contractID,
		Choice:     "ICCIPReceiver_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// ICCIPReceiverExecute exercises the ICCIPReceiver_Execute choice on this GlobalTokenSenderReceiver contract via the IICCIPReceiver interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) ICCIPReceiverExecute(contractID string, args user.ICCIPReceiverExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "ICCIPReceiver"),
		ContractID: contractID,
		Choice:     "ICCIPReceiver_Execute",
		Arguments:  argsToMap(args),
	}
}

// ICCIPReceiverExecuteWithPackageID exercises the ICCIPReceiver_Execute choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ICCIPReceiverExecuteWithPackageID(contractID string, packageID string, args user.ICCIPReceiverExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "ICCIPReceiver"),
		ContractID: contractID,
		Choice:     "ICCIPReceiver_Execute",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderGetRequiredCCVs exercises the ICCIPSender_GetRequiredCCVs choice on this GlobalTokenSenderReceiver contract via the IICCIPSender interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) ICCIPSenderGetRequiredCCVs(contractID string, args user.ICCIPSenderGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderGetRequiredCCVsWithPackageID exercises the ICCIPSender_GetRequiredCCVs choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ICCIPSenderGetRequiredCCVsWithPackageID(contractID string, packageID string, args user.ICCIPSenderGetRequiredCCVs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_GetRequiredCCVs",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderGetFee exercises the ICCIPSender_GetFee choice on this GlobalTokenSenderReceiver contract via the IICCIPSender interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) ICCIPSenderGetFee(contractID string, args user.ICCIPSenderGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_GetFee",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderGetFeeWithPackageID exercises the ICCIPSender_GetFee choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ICCIPSenderGetFeeWithPackageID(contractID string, packageID string, args user.ICCIPSenderGetFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_GetFee",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderSend exercises the ICCIPSender_Send choice on this GlobalTokenSenderReceiver contract via the IICCIPSender interface
// This method uses the package name in the template ID
func (t GlobalTokenSenderReceiver) ICCIPSenderSend(contractID string, args user.ICCIPSenderSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_Send",
		Arguments:  argsToMap(args),
	}
}

// ICCIPSenderSendWithPackageID exercises the ICCIPSender_Send choice using the provided package ID instead of package name
func (t GlobalTokenSenderReceiver) ICCIPSenderSendWithPackageID(contractID string, packageID string, args user.ICCIPSenderSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.GlobalTokenSenderReceiver", "ICCIPSender"),
		ContractID: contractID,
		Choice:     "ICCIPSender_Send",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for GlobalTokenSenderReceiver

var _ user.IICCIPReceiver = (*GlobalTokenSenderReceiver)(nil)

var _ user.IICCIPSender = (*GlobalTokenSenderReceiver)(nil)

// Send is a Record type
type Send struct {
	SenderParty              types.PARTY                                `json:"senderParty"`
	DestinationChainSelector types.NUMERIC                              `json:"destinationChainSelector"`
	Message                  clientapi.Canton2AnyMessage                `json:"message"`
	RouterCid                types.CONTRACT_ID                          `json:"routerCid"`
	FeeTokenInput            user.SenderFeeTokenInput                   `json:"feeTokenInput"`
	CcvSendInputs            []user.SenderCCVSendInput                  `json:"ccvSendInputs"`
	TokenTransferInput       *user.SenderTokenTransferInput             `json:"tokenTransferInput" hex:"optional"`
	ExecutorInput            *user.SenderExecutorInput                  `json:"executorInput" hex:"optional"`
	Context                  splice_api_token_metadata_v1.ChoiceContext `json:"context"`
}

// ToMap converts Send to a map for DAML arguments
func (t Send) ToMap() map[string]any {
	m := make(map[string]any)

	m["senderParty"] = t.SenderParty.ToMap()

	m["destinationChainSelector"] = t.DestinationChainSelector

	m["message"] = model.NestedToDAMLValue(t.Message)

	m["routerCid"] = model.NestedToDAMLValue(t.RouterCid)

	m["feeTokenInput"] = model.NestedToDAMLValue(t.FeeTokenInput)

	m["ccvSendInputs"] = func() []any {
		res := make([]any, 0, len(t.CcvSendInputs))
		for _, e := range t.CcvSendInputs {
			res = append(res, model.NestedToDAMLValue(e))
		}
		return res
	}()

	if t.TokenTransferInput != nil {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.TokenTransferInput),
		}
	} else {
		m["tokenTransferInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	if t.ExecutorInput != nil {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": model.NestedToDAMLValue(*t.ExecutorInput),
		}
	} else {
		m["executorInput"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["context"] = model.NestedToDAMLValue(t.Context)

	return m
}

func (t Send) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *Send) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes Send to hex string (Canton MCMS format)
func (t Send) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes Send from hex string (Canton MCMS format)
func (t *Send) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	Execute(args Execute) (*bind.EncodedChoice, error)
	GetFee(args GetFee) (*bind.EncodedChoice, error)
	GetRequiredCCVsForExecute(args GetRequiredCCVsForExecute) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error)
	Send(args Send) (*bind.EncodedChoice, error)
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

// Execute encodes parameters for the Execute choice.
func (e *encoder) Execute(args Execute) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Execute", args)
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

// Send encodes parameters for the Send choice.
func (e *encoder) Send(args Send) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("Send", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
