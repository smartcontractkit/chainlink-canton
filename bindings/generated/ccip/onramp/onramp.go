package onramp

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	mcms "github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
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
	PackageName = "ccip-onramp"
	PackageID   = "cde3c02229c4887283ab9701223773dea65046146377c37e689d79192309bd01"
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

// CCIPSendFromRouter is a Record type
type CCIPSendFromRouter struct {
	RouterPartyOwner      types.PARTY       `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT        `json:"routerInstanceId"`
	CurrentSequenceNumber types.NUMERIC     `json:"currentSequenceNumber"`
	RmnRemoteCid          types.CONTRACT_ID `json:"rmnRemoteCid"`
	GlobalConfigCid       types.CONTRACT_ID `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID `json:"tokenAdminRegistryCid"`
	SendingMessageCid     types.CONTRACT_ID `json:"sendingMessageCid"`
}

// ToMap converts CCIPSendFromRouter to a map for DAML arguments
func (t CCIPSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["currentSequenceNumber"] = t.CurrentSequenceNumber

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

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
	Receipts             []common.Receipt  `json:"receipts"`
}

// ToMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["ccipMessageSent"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcipMessageSent).(mapper); ok {
			return m.toMap()
		}
		return t.CcipMessageSent
	}()

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

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	GlobalConfigCid   types.CONTRACT_ID `json:"globalConfigCid"`
	DestChainSelector types.NUMERIC     `json:"destChainSelector"`
}

// ToMap converts GetRequiredCCVsForSend to a map for DAML arguments
func (t GetRequiredCCVsForSend) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["destChainSelector"] = t.DestChainSelector

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

// OnRamp is a Template type
type OnRamp struct {
	InstanceId types.TEXT  `json:"instanceId"`
	CcipOwner  types.PARTY `json:"ccipOwner"`
	Deps       OnRampDeps  `json:"deps"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OnRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OnRamp) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OnRamp) CreateCommand() *model.CreateCommand {
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

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
	}()

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
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CCIPSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// CCIPSendFromRouterWithPackageID exercises the CCIPSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) CCIPSendFromRouterWithPackageID(contractID string, packageID string, args CCIPSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CCIPSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendFromRouter exercises the PrepareSendFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) PrepareSendFromRouter(contractID string, args PrepareSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "PrepareSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareSendFromRouterWithPackageID exercises the PrepareSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) PrepareSendFromRouterWithPackageID(contractID string, packageID string, args PrepareSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "PrepareSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// SetDeps exercises the SetDeps choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) SetDeps(contractID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// SetDepsWithPackageID exercises the SetDeps choice using the provided package ID instead of package name
func (t OnRamp) SetDepsWithPackageID(contractID string, packageID string, args SetDeps) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "SetDeps",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSend exercises the GetRequiredCCVsForSend choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) GetRequiredCCVsForSend(contractID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForSendWithPackageID exercises the GetRequiredCCVsForSend choice using the provided package ID instead of package name
func (t OnRamp) GetRequiredCCVsForSendWithPackageID(contractID string, packageID string, args GetRequiredCCVsForSend) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForSend",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OnRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this OnRamp contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t OnRamp) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t OnRamp) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for OnRamp

var _ mcms.IMCMSReceiver = (*OnRamp)(nil)

// OnRampDeps is a Record type
type OnRampDeps struct {
	GlobalConfig       common.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          common.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry common.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          common.RawInstanceAddress `json:"feeQuoter"`
	CcvRegistry        common.RawInstanceAddress `json:"ccvRegistry"`
}

// ToMap converts OnRampDeps to a map for DAML arguments
func (t OnRampDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
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

	m["ccvRegistry"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvRegistry).(mapper); ok {
			return m.toMap()
		}
		return t.CcvRegistry
	}()

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

// PrepareSendFromRouter is a Record type
type PrepareSendFromRouter struct {
	RouterPartyOwner      types.PARTY                               `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT                                `json:"routerInstanceId"`
	GlobalConfigCid       types.CONTRACT_ID                         `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID                         `json:"tokenAdminRegistryCid"`
	FeeQuoterCid          types.CONTRACT_ID                         `json:"feeQuoterCid"`
	RmnRemoteCid          types.CONTRACT_ID                         `json:"rmnRemoteCid"`
	DestChainSelector     types.NUMERIC                             `json:"destChainSelector"`
	Receiver              types.TEXT                                `json:"receiver"`
	Payload               types.TEXT                                `json:"payload"`
	CcipReceiveGasLimit   types.INT64                               `json:"ccipReceiveGasLimit"`
	BlockConfirmations    *types.INT64                              `json:"blockConfirmations" hex:"optional"`
	CurrentSequenceNumber types.NUMERIC                             `json:"currentSequenceNumber"`
	SenderRequiredCCVs    []common.RawInstanceAddress               `json:"senderRequiredCCVs"`
	TokenInstrumentId     *splice_api_token_holding_v1.InstrumentId `json:"tokenInstrumentId" hex:"optional"`
	TokenReceiver         *types.TEXT                               `json:"tokenReceiver" hex:"optional"`
	TokenArgs             types.TEXT                                `json:"tokenArgs"`
	FeeToken              splice_api_token_holding_v1.InstrumentId  `json:"feeToken"`
}

// ToMap converts PrepareSendFromRouter to a map for DAML arguments
func (t PrepareSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["globalConfigCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigCid).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	m["feeQuoterCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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

	m["currentSequenceNumber"] = t.CurrentSequenceNumber

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

// SetDeps is a Record type
type SetDeps struct {
	NewDeps SetDepsParams `json:"newDeps"`
}

// ToMap converts SetDeps to a map for DAML arguments
func (t SetDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["newDeps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.NewDeps).(mapper); ok {
			return m.toMap()
		}
		return t.NewDeps
	}()

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
	GlobalConfig       *common.RawInstanceAddress `json:"globalConfig" hex:"optional"`
	RmnRemote          *common.RawInstanceAddress `json:"rmnRemote" hex:"optional"`
	TokenAdminRegistry *common.RawInstanceAddress `json:"tokenAdminRegistry" hex:"optional"`
	FeeQuoter          *common.RawInstanceAddress `json:"feeQuoter" hex:"optional"`
	CcvRegistry        *common.RawInstanceAddress `json:"ccvRegistry" hex:"optional"`
}

// ToMap converts SetDepsParams to a map for DAML arguments
func (t SetDepsParams) ToMap() map[string]any {
	m := make(map[string]any)

	if t.GlobalConfig != nil {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
			"value": *t.GlobalConfig,
		}
	} else {
		m["globalConfig"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.RmnRemote != nil {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
			"value": *t.RmnRemote,
		}
	} else {
		m["rmnRemote"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.TokenAdminRegistry != nil {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
			"value": *t.TokenAdminRegistry,
		}
	} else {
		m["tokenAdminRegistry"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.FeeQuoter != nil {
		m["feeQuoter"] = map[string]any{
			"_type": "optional",
			"value": *t.FeeQuoter,
		}
	} else {
		m["feeQuoter"] = map[string]any{
			"_type": "optional",
		}
	}

	if t.CcvRegistry != nil {
		m["ccvRegistry"] = map[string]any{
			"_type": "optional",
			"value": *t.CcvRegistry,
		}
	} else {
		m["ccvRegistry"] = map[string]any{
			"_type": "optional",
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
	CCIPSendFromRouter(args CCIPSendFromRouter) (*bind.EncodedChoice, error)
	GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error)
	PrepareSendFromRouter(args PrepareSendFromRouter) (*bind.EncodedChoice, error)
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

// CCIPSendFromRouter encodes parameters for the CCIPSendFromRouter choice.
func (e *encoder) CCIPSendFromRouter(args CCIPSendFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CCIPSendFromRouter", args)
}

// GetRequiredCCVsForSend encodes parameters for the GetRequiredCCVsForSend choice.
func (e *encoder) GetRequiredCCVsForSend(args GetRequiredCCVsForSend) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetRequiredCCVsForSend", args)
}

// PrepareSendFromRouter encodes parameters for the PrepareSendFromRouter choice.
func (e *encoder) PrepareSendFromRouter(args PrepareSendFromRouter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("PrepareSendFromRouter", args)
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
