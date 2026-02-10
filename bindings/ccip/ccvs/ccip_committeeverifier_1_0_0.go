package ccvs

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/go-daml/pkg/codec"
	"github.com/smartcontractkit/go-daml/pkg/model"
	. "github.com/smartcontractkit/go-daml/pkg/types"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

const PackageName = "ccip-committeeverifier"
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

// CCVFeeConfig is a Record type
type CCVFeeConfig struct {
	FeeUSDCents        NUMERIC `json:"feeUSDCents"`
	GasForVerification INT64   `json:"gasForVerification"`
	PayloadSizeBytes   INT64   `json:"payloadSizeBytes"`
}

// ToMap converts CCVFeeConfig to a map for DAML arguments
func (t CCVFeeConfig) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["feeUSDCents"] = t.FeeUSDCents

	m["gasForVerification"] = int64(t.GasForVerification)

	m["payloadSizeBytes"] = int64(t.PayloadSizeBytes)

	return m
}

func (t CCVFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CCVFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CommitteeVerifier is a Template type
type CommitteeVerifier struct {
	InstanceId               TEXT               `json:"instanceId"`
	Owner                    PARTY              `json:"owner"`
	CcipOwner                PARTY              `json:"ccipOwner"`
	VersionTag               TEXT               `json:"versionTag"`
	MessageSentObserver      PARTY              `json:"messageSentObserver"`
	StorageLocation          TEXT               `json:"storageLocation"`
	Threshold                INT64              `json:"threshold"`
	Signers                  []TEXT             `json:"signers"`
	RmnRemoteInstanceAddress RawInstanceAddress `json:"rmnRemoteInstanceAddress"`
	RemoteChainFeeConfigs    GENMAP             `json:"remoteChainFeeConfigs"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CommitteeVerifier) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CommitteeVerifier) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CommitteeVerifier) CreateCommand() *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["versionTag"] = string(t.VersionTag)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageSentObserver"] = t.MessageSentObserver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocation"] = string(t.StorageLocation)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["threshold"] = int64(t.Threshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, string(e))
		}
		return res
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
	args["remoteChainFeeConfigs"] = func() interface{} {
		if t.RemoteChainFeeConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RemoteChainFeeConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CommitteeVerifier) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]interface{})

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["ccipOwner"] = t.CcipOwner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["versionTag"] = string(t.VersionTag)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["messageSentObserver"] = t.MessageSentObserver.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocation"] = string(t.StorageLocation)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["threshold"] = int64(t.Threshold)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signers"] = func() []interface{} {
		res := make([]interface{}, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, string(e))
		}
		return res
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
	args["remoteChainFeeConfigs"] = func() interface{} {
		if t.RemoteChainFeeConfigs == nil {
			return map[string]interface{}{"_type": "genmap", "value": GENMAP{}}
		}
		return map[string]interface{}{"_type": "genmap", "value": t.RemoteChainFeeConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CommitteeVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CommitteeVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// Choice methods for CommitteeVerifier

// CommitteeVerifierVerifyMessage exercises the CommitteeVerifier_VerifyMessage choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierVerifyMessage(contractID string, args CommitteeVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierVerifyMessageWithPackageID exercises the CommitteeVerifier_VerifyMessage choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierVerifyMessageWithPackageID(contractID string, packageID string, args CommitteeVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierCalculateFee exercises the CommitteeVerifier_CalculateFee choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierCalculateFee(contractID string, args CommitteeVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierCalculateFeeWithPackageID exercises the CommitteeVerifier_CalculateFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierCalculateFeeWithPackageID(contractID string, packageID string, args CommitteeVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CommitteeVerifier) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]interface{}{},
	}
}

// CommitteeVerifierForwardToVerifier exercises the CommitteeVerifier_ForwardToVerifier choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierForwardToVerifier(contractID string, args CommitteeVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierForwardToVerifierWithPackageID exercises the CommitteeVerifier_ForwardToVerifier choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierForwardToVerifierWithPackageID(contractID string, packageID string, args CommitteeVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessage exercises the CrossChainVerifier_VerifyMessage choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierVerifyMessage(contractID string, args CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessageWithPackageID exercises the CrossChainVerifier_VerifyMessage choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierVerifyMessageWithPackageID(contractID string, packageID string, args CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFee exercises the CrossChainVerifier_CalculateFee choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierCalculateFee(contractID string, args CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFeeWithPackageID exercises the CrossChainVerifier_CalculateFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierCalculateFeeWithPackageID(contractID string, packageID string, args CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifier exercises the CrossChainVerifier_ForwardToVerifier choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifier(contractID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifierWithPackageID exercises the CrossChainVerifier_ForwardToVerifier choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifierWithPackageID(contractID string, packageID string, args CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CommitteeVerifier

var _ IICrossChainVerifier = (*CommitteeVerifier)(nil)

// CommitteeVerifierCalculateFee is a Record type
type CommitteeVerifierCalculateFee struct {
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	Caller            PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierCalculateFee to a map for DAML arguments
func (t CommitteeVerifierCalculateFee) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CommitteeVerifierCalculateFee) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CommitteeVerifierCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CommitteeVerifierForwardToVerifier is a Record type
type CommitteeVerifierForwardToVerifier struct {
	RmnRemoteCid      CONTRACT_ID `json:"rmnRemoteCid"`
	SendingMessageCid CONTRACT_ID `json:"sendingMessageCid"`
	VerifierArgs      TEXT        `json:"verifierArgs"`
	Caller            PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierForwardToVerifier to a map for DAML arguments
func (t CommitteeVerifierForwardToVerifier) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["sendingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.SendingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.SendingMessageCid
	}()

	m["verifierArgs"] = string(t.VerifierArgs)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CommitteeVerifierForwardToVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CommitteeVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}

// CommitteeVerifierVerifyMessage is a Record type
type CommitteeVerifierVerifyMessage struct {
	RmnRemoteCid        CONTRACT_ID `json:"rmnRemoteCid"`
	ExecutingMessageCid CONTRACT_ID `json:"executingMessageCid"`
	VerifierResults     TEXT        `json:"verifierResults"`
	Caller              PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierVerifyMessage to a map for DAML arguments
func (t CommitteeVerifierVerifyMessage) ToMap() map[string]interface{} {
	m := make(map[string]interface{})

	m["rmnRemoteCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["executingMessageCid"] = func() interface{} {
		type mapper interface{ toMap() map[string]interface{} }
		if m, ok := any(t.ExecutingMessageCid).(mapper); ok {
			return m.toMap()
		}
		return t.ExecutingMessageCid
	}()

	m["verifierResults"] = string(t.VerifierResults)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t CommitteeVerifierVerifyMessage) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshall(t)
}

func (t *CommitteeVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshall(data, t)
}
