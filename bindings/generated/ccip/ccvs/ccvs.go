package ccvs

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-canton/bindings/codec"
	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	model "github.com/smartcontractkit/chainlink-canton/bindings/ledger"
	"github.com/smartcontractkit/chainlink-canton/bindings/types"
)

var (
	_ = fmt.Sprintf
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = model.Command{}
)

const (
	PackageName = "ccip-committeeverifier"
	PackageID   = "734a12924b71bc59594246101aa6d63454317960829bc8e8697c48717142277a"
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

// CCVFeeConfig is a Record type
type CCVFeeConfig struct {
	FeeUSDCents        types.NUMERIC `json:"feeUSDCents"`
	GasForVerification types.INT64   `json:"gasForVerification"`
	PayloadSizeBytes   types.INT64   `json:"payloadSizeBytes"`
}

// ToMap converts CCVFeeConfig to a map for DAML arguments
func (t CCVFeeConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["feeUSDCents"] = t.FeeUSDCents

	m["gasForVerification"] = int64(t.GasForVerification)

	m["payloadSizeBytes"] = int64(t.PayloadSizeBytes)

	return m
}

func (t CCVFeeConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCVFeeConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// CommitteeVerifier is a Template type
type CommitteeVerifier struct {
	InstanceId               types.TEXT                `json:"instanceId"`
	Owner                    types.PARTY               `json:"owner"`
	CcipOwner                types.PARTY               `json:"ccipOwner"`
	VersionTag               types.TEXT                `json:"versionTag"`
	MessageSentObserver      types.PARTY               `json:"messageSentObserver"`
	StorageLocation          types.TEXT                `json:"storageLocation"`
	Threshold                types.INT64               `json:"threshold"`
	Signers                  []types.TEXT              `json:"signers"`
	RmnRemoteInstanceAddress common.RawInstanceAddress `json:"rmnRemoteInstanceAddress"`
	RemoteChainFeeConfigs    types.GENMAP              `json:"remoteChainFeeConfigs"`
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
	args := make(map[string]any)

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
	args["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, string(e))
		}
		return res
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
	args["remoteChainFeeConfigs"] = func() any {
		if t.RemoteChainFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainFeeConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CommitteeVerifier) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

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
	args["signers"] = func() []any {
		res := make([]any, 0, len(t.Signers))
		for _, e := range t.Signers {
			res = append(res, string(e))
		}
		return res
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
	args["remoteChainFeeConfigs"] = func() any {
		if t.RemoteChainFeeConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.RemoteChainFeeConfigs}
	}()

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CommitteeVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
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
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CommitteeVerifier) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
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
func (t CommitteeVerifier) CrossChainVerifierVerifyMessage(contractID string, args common.CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierVerifyMessageWithPackageID exercises the CrossChainVerifier_VerifyMessage choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierVerifyMessageWithPackageID(contractID string, packageID string, args common.CrossChainVerifierVerifyMessage) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_VerifyMessage",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFee exercises the CrossChainVerifier_CalculateFee choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierCalculateFee(contractID string, args common.CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierCalculateFeeWithPackageID exercises the CrossChainVerifier_CalculateFee choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierCalculateFeeWithPackageID(contractID string, packageID string, args common.CrossChainVerifierCalculateFee) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_CalculateFee",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifier exercises the CrossChainVerifier_ForwardToVerifier choice on this CommitteeVerifier contract via the IICrossChainVerifier interface
// This method uses the package name in the template ID
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifier(contractID string, args common.CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// CrossChainVerifierForwardToVerifierWithPackageID exercises the CrossChainVerifier_ForwardToVerifier choice using the provided package ID instead of package name
func (t CommitteeVerifier) CrossChainVerifierForwardToVerifierWithPackageID(contractID string, packageID string, args common.CrossChainVerifierForwardToVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "ICrossChainVerifier"),
		ContractID: contractID,
		Choice:     "CrossChainVerifier_ForwardToVerifier",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CommitteeVerifier

var _ common.IICrossChainVerifier = (*CommitteeVerifier)(nil)

// CommitteeVerifierCalculateFee is a Record type
type CommitteeVerifierCalculateFee struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierCalculateFee to a map for DAML arguments
func (t CommitteeVerifierCalculateFee) ToMap() map[string]any {
	m := make(map[string]any)

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierCalculateFee) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// CommitteeVerifierForwardToVerifier is a Record type
type CommitteeVerifierForwardToVerifier struct {
	RmnRemoteCid      types.CONTRACT_ID `json:"rmnRemoteCid"`
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT        `json:"verifierArgs"`
	Caller            types.PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierForwardToVerifier to a map for DAML arguments
func (t CommitteeVerifierForwardToVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["sendingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierForwardToVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// CommitteeVerifierVerifyMessage is a Record type
type CommitteeVerifierVerifyMessage struct {
	RmnRemoteCid        types.CONTRACT_ID `json:"rmnRemoteCid"`
	ExecutingMessageCid types.CONTRACT_ID `json:"executingMessageCid"`
	VerifierResults     types.TEXT        `json:"verifierResults"`
	Caller              types.PARTY       `json:"caller"`
}

// ToMap converts CommitteeVerifierVerifyMessage to a map for DAML arguments
func (t CommitteeVerifierVerifyMessage) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemoteCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemoteCid).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemoteCid
	}()

	m["executingMessageCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
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
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierVerifyMessage) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}
