package onramp

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
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
	PackageID   = "ca0db2ad8d8fe60d0393e2b11e7b43176360da643bcf723bd99499341eac07c0"
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

// CCIPSendFromRouterResult is a Record type
type CCIPSendFromRouterResult struct {
	VerifierBlobs        []types.TEXT     `json:"verifierBlobs"`
	MessageSentObservers []types.PARTY    `json:"messageSentObservers"`
	Receipts             []common.Receipt `json:"receipts"`
}

// ToMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) ToMap() map[string]any {
	m := make(map[string]any)

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

// CancelSendFromRouter is a Record type
type CancelSendFromRouter struct {
	RouterPartyOwner      types.PARTY       `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT        `json:"routerInstanceId"`
	SendingMessageCid     types.CONTRACT_ID `json:"sendingMessageCid"`
	TokenAdminRegistryCid types.CONTRACT_ID `json:"tokenAdminRegistryCid"`
}

// ToMap converts CancelSendFromRouter to a map for DAML arguments
func (t CancelSendFromRouter) ToMap() map[string]any {
	m := make(map[string]any)

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenAdminRegistryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	return m
}

func (t CancelSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CancelSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
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

// OnRamp is a Template type
type OnRamp struct {
	InstanceId                        types.TEXT                `json:"instanceId"`
	CcipOwner                         types.PARTY               `json:"ccipOwner"`
	GlobalConfigInstanceAddress       common.RawInstanceAddress `json:"globalConfigInstanceAddress"`
	RmnRemoteInstanceAddress          common.RawInstanceAddress `json:"rmnRemoteInstanceAddress"`
	TokenAdminRegistryInstanceAddress common.RawInstanceAddress `json:"tokenAdminRegistryInstanceAddress"`
	CcvRegistryInstanceAddress        common.RawInstanceAddress `json:"ccvRegistryInstanceAddress"`
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
	args["globalConfigInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnRemoteInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.CcvRegistryInstanceAddress
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
	args["globalConfigInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfigInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnRemoteInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvRegistryInstanceAddress"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.CcvRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.CcvRegistryInstanceAddress
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

// CancelSendFromRouter exercises the CancelSendFromRouter choice on this OnRamp contract
// This method uses the package name in the template ID
func (t OnRamp) CancelSendFromRouter(contractID string, args CancelSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CancelSendFromRouter",
		Arguments:  argsToMap(args),
	}
}

// CancelSendFromRouterWithPackageID exercises the CancelSendFromRouter choice using the provided package ID instead of package name
func (t OnRamp) CancelSendFromRouterWithPackageID(contractID string, packageID string, args CancelSendFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "CancelSendFromRouter",
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

// PrepareSendFromRouter is a Record type
type PrepareSendFromRouter struct {
	RouterPartyOwner      types.PARTY                              `json:"routerPartyOwner"`
	RouterInstanceId      types.TEXT                               `json:"routerInstanceId"`
	GlobalConfigCid       types.CONTRACT_ID                        `json:"globalConfigCid"`
	TokenAdminRegistryCid types.CONTRACT_ID                        `json:"tokenAdminRegistryCid"`
	FeeQuoterCid          types.CONTRACT_ID                        `json:"feeQuoterCid"`
	RmnRemoteCid          types.CONTRACT_ID                        `json:"rmnRemoteCid"`
	DestChainSelector     types.NUMERIC                            `json:"destChainSelector"`
	Receiver              types.TEXT                               `json:"receiver"`
	Payload               types.TEXT                               `json:"payload"`
	CcipReceiveGasLimit   types.INT64                              `json:"ccipReceiveGasLimit"`
	CurrentSequenceNumber types.NUMERIC                            `json:"currentSequenceNumber"`
	SenderRequiredCCVs    []common.RawInstanceAddress              `json:"senderRequiredCCVs"`
	WithTokenTransfer     types.BOOL                               `json:"withTokenTransfer"`
	TokenReceiver         *types.TEXT                              `json:"tokenReceiver"`
	FeeToken              splice_api_token_holding_v1.InstrumentId `json:"feeToken"`
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

	m["withTokenTransfer"] = bool(t.WithTokenTransfer)

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
