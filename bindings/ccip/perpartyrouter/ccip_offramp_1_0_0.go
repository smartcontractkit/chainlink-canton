package perpartyrouter

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

// ExecuteFromRouter is a Record type
type ExecuteFromRouter struct {
	RouterPartyOwner      PARTY                `json:"routerPartyOwner"`
	ReceiverRequiredCCVs  []RawInstanceAddress `json:"receiverRequiredCCVs"`
	ExecutingMessageCid   CONTRACT_ID          `json:"executingMessageCid"`
	GlobalConfigCid       CONTRACT_ID          `json:"globalConfigCid"`
	TokenAdminRegistryCid CONTRACT_ID          `json:"tokenAdminRegistryCid"`
	RmnRemoteCid          CONTRACT_ID          `json:"rmnRemoteCid"`
}

// ToMap converts ExecuteFromRouter to a map for DAML arguments
func (t ExecuteFromRouter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["routerPartyOwner"] = t.RouterPartyOwner.ToMap()

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
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

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

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
	MessageId           TEXT         `json:"messageId"`
	Message             MessageV1    `json:"message"`
	SourceChainSelector NUMERIC      `json:"sourceChainSelector"`
	SequenceNumber      NUMERIC      `json:"sequenceNumber"`
	TokenReceiveTicket  *CONTRACT_ID `json:"tokenReceiveTicket"`
}

// ToMap converts ExecuteFromRouterResult to a map for DAML arguments
func (t ExecuteFromRouterResult) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["messageId"] = string(t.MessageId)

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

// GetRequiredCCVsForExecute is a Record type
type GetRequiredCCVsForExecute struct {
	GlobalConfigCid      CONTRACT_ID          `json:"globalConfigCid"`
	ReceiverRequiredCCVs []RawInstanceAddress `json:"receiverRequiredCCVs"`
	SourceChainSelector  NUMERIC              `json:"sourceChainSelector"`
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

	m["receiverRequiredCCVs"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.ReceiverRequiredCCVs))
		for _, e := range t.ReceiverRequiredCCVs {
			type mapper interface{ toMap() map[string]interface{} }
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

func (t GetRequiredCCVsForExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *GetRequiredCCVsForExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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

// OffRamp is a Template type
type OffRamp struct {
	InstanceId                        TEXT               `json:"instanceId"`
	CcipOwner                         PARTY              `json:"ccipOwner"`
	GlobalConfigInstanceAddress       RawInstanceAddress `json:"globalConfigInstanceAddress"`
	RmnRemoteInstanceAddress          RawInstanceAddress `json:"rmnRemoteInstanceAddress"`
	TokenAdminRegistryInstanceAddress RawInstanceAddress `json:"tokenAdminRegistryInstanceAddress"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t OffRamp) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t OffRamp) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t OffRamp) CreateCommand() *model.CreateCommand {
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

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t OffRamp) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
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

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
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

// GetRequiredCCVsForExecute exercises the GetRequiredCCVsForExecute choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) GetRequiredCCVsForExecute(contractID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// GetRequiredCCVsForExecuteWithPackageID exercises the GetRequiredCCVsForExecute choice using the provided package ID instead of package name
func (t OffRamp) GetRequiredCCVsForExecuteWithPackageID(contractID string, packageID string, args GetRequiredCCVsForExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "GetRequiredCCVsForExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute exercises the PrepareExecute choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) PrepareExecute(contractID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecuteWithPackageID exercises the PrepareExecute choice using the provided package ID instead of package name
func (t OffRamp) PrepareExecuteWithPackageID(contractID string, packageID string, args PrepareExecute) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "PrepareExecute",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t OffRamp) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ExecuteFromRouter exercises the ExecuteFromRouter choice on this OffRamp contract
// This method uses the package name in the template ID
func (t OffRamp) ExecuteFromRouter(contractID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// ExecuteFromRouterWithPackageID exercises the ExecuteFromRouter choice using the provided package ID instead of package name
func (t OffRamp) ExecuteFromRouterWithPackageID(contractID string, packageID string, args ExecuteFromRouter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.OffRamp", "OffRamp"),
		ContractID: contractID,
		Choice:     "ExecuteFromRouter",
		Arguments:  argsToMap(args),
	}
}

// PrepareExecute is a Record type
type PrepareExecute struct {
	EncodedMessage     TEXT        `json:"encodedMessage"`
	RmnRemoteCid       CONTRACT_ID `json:"rmnRemoteCid"`
	ReceiverParty      PARTY       `json:"receiverParty"`
	TokenReceiverParty *PARTY      `json:"tokenReceiverParty"`
	Caller             PARTY       `json:"caller"`
}

// ToMap converts PrepareExecute to a map for DAML arguments
func (t PrepareExecute) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["encodedMessage"] = string(t.EncodedMessage)

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["receiverParty"] = t.ReceiverParty.ToMap()

	if t.TokenReceiverParty != nil {
		m["tokenReceiverParty"] = map[string]interface{}{
			"_type": "optional",
			"value": (*t.TokenReceiverParty).ToMap(),
		}
	} else {
		m["tokenReceiverParty"] = map[string]interface{}{
			"_type": "optional",
		}
	}

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t PrepareExecute) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *PrepareExecute) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
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
