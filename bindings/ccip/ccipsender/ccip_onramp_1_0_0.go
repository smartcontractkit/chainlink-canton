package ccipsender

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
	_ bind.BoundTemplate
)

// CCIPSendFromRouter is a Record type
type CCIPSendFromRouter struct {
	RouterPartyOwner      PARTY       `json:"routerPartyOwner"`
	RouterInstanceId      TEXT        `json:"routerInstanceId"`
	CurrentSequenceNumber NUMERIC     `json:"currentSequenceNumber"`
	RmnRemoteCid          CONTRACT_ID `json:"rmnRemoteCid"`
	GlobalConfigCid       CONTRACT_ID `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
	SendingMessageCid     CONTRACT_ID `json:"sendingMessageCid"`
}

// ToMap converts CCIPSendFromRouter to a map for DAML arguments
func (t CCIPSendFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["currentSequenceNumber"] = t.CurrentSequenceNumber

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
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

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	return m
}

func (t CCIPSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCIPSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	VerifierBlobs        []TEXT    `json:"verifierBlobs"`
	MessageSentObservers []PARTY   `json:"messageSentObservers"`
	Receipts             []Receipt `json:"receipts"`
}

// ToMap converts CCIPSendFromRouterResult to a map for DAML arguments
func (t CCIPSendFromRouterResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

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

	m["receipts"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Receipts))
		for _, e := range t.Receipts {
			type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *CCIPSendFromRouterResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// CancelSendFromRouter is a Record type
type CancelSendFromRouter struct {
	RouterPartyOwner      PARTY       `json:"routerPartyOwner"`
	RouterInstanceId      TEXT        `json:"routerInstanceId"`
	SendingMessageCid     CONTRACT_ID `json:"sendingMessageCid"`
	TokenAdminRegistryCid CONTRACT_ID `json:"tokenAdminRegistryCid"`
}

// ToMap converts CancelSendFromRouter to a map for DAML arguments
func (t CancelSendFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["tokenAdminRegistryCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryCid).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryCid
	}()

	return m
}

func (t CancelSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CancelSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// MarshalHex encodes CancelSendFromRouter to hex string (Canton MCMS format)
func (t CancelSendFromRouter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CancelSendFromRouter from hex string (Canton MCMS format)
func (t *CancelSendFromRouter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetRequiredCCVsForSend is a Record type
type GetRequiredCCVsForSend struct {
	GlobalConfigCid   CONTRACT_ID `json:"globalConfigCid"`
	DestChainSelector NUMERIC     `json:"destChainSelector"`
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

	m["destChainSelector"] = t.DestChainSelector

	return m
}

func (t GetRequiredCCVsForSend) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetRequiredCCVsForSend) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
	InstanceId                        TEXT               `json:"instanceId"`
	CcipOwner                         PARTY              `json:"ccipOwner"`
	GlobalConfigInstanceAddress       RawInstanceAddress `json:"globalConfigInstanceAddress"`
	RmnRemoteInstanceAddress          RawInstanceAddress `json:"rmnRemoteInstanceAddress"`
	TokenAdminRegistryInstanceAddress RawInstanceAddress `json:"tokenAdminRegistryInstanceAddress"`
	CcvRegistryInstanceAddress        RawInstanceAddress `json:"ccvRegistryInstanceAddress"`
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["globalConfigInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnRemoteInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvRegistryInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["globalConfigInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.GlobalConfigInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfigInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["rmnRemoteInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["tokenAdminRegistryInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.TokenAdminRegistryInstanceAddress).(mapper); ok {
			return m.toMap()
		}
		return t.TokenAdminRegistryInstanceAddress
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccvRegistryInstanceAddress"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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
	return jsonCodec.Marshall(t)
}

func (t *OnRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OnRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OnRamp", "OnRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// PrepareSendFromRouter is a Record type
type PrepareSendFromRouter struct {
	RouterPartyOwner      PARTY                `json:"routerPartyOwner"`
	RouterInstanceId      TEXT                 `json:"routerInstanceId"`
	GlobalConfigCid       CONTRACT_ID          `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID          `json:"tokenAdminRegistryCid"`
	FeeQuoterCid          CONTRACT_ID          `json:"feeQuoterCid"`
	RmnRemoteCid          CONTRACT_ID          `json:"rmnRemoteCid"`
	DestChainSelector     NUMERIC              `json:"destChainSelector"`
	Receiver              TEXT                 `json:"receiver"`
	Payload               TEXT                 `json:"payload"`
	CcipReceiveGasLimit   INT64                `json:"ccipReceiveGasLimit"`
	CurrentSequenceNumber NUMERIC              `json:"currentSequenceNumber"`
	SenderRequiredCCVs    []RawInstanceAddress `json:"senderRequiredCCVs"`
	WithTokenTransfer     BOOL                 `json:"withTokenTransfer"`
	TokenReceiver         *TEXT                `json:"tokenReceiver"`
	FeeToken              InstrumentId         `json:"feeToken"`
}

// ToMap converts PrepareSendFromRouter to a map for DAML arguments
func (t PrepareSendFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["routerInstanceId"] = string(t.RouterInstanceId)

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

	m["feeQuoterCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeQuoterCid).(mapper); ok {
			return m.toMap()
		}
		return t.FeeQuoterCid
	}()

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
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

	m["senderRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.SenderRequiredCCVs))
		for _, e := range t.SenderRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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
		m["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
			"value": string(*t.TokenReceiver),
		}
	} else {
		m["tokenReceiver"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["feeToken"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.FeeToken).(mapper); ok {
			return m.toMap()
		}
		return t.FeeToken
	}()

	return m
}

func (t PrepareSendFromRouter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PrepareSendFromRouter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
