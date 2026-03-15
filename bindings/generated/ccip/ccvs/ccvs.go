package ccvs

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
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
	PackageName = "ccip-committeeverifier"
	PackageID   = "5d37bb112e02f511f5659d92c1f7e3cb3548d2e3a3d0e8ac9905e0b21162acbe"
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

// MarshalHex encodes CCVFeeConfig to hex string (Canton MCMS format)
func (t CCVFeeConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCVFeeConfig from hex string (Canton MCMS format)
func (t *CCVFeeConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifier is a Template type
type CommitteeVerifier struct {
	InstanceId                   types.TEXT            `json:"instanceId"`
	Owner                        types.PARTY           `json:"owner"`
	CcipOwner                    types.PARTY           `json:"ccipOwner"`
	VersionTag                   types.TEXT            `json:"versionTag"`
	MessageSentObserver          types.PARTY           `json:"messageSentObserver"`
	StorageLocations             []types.TEXT          `json:"storageLocations"`
	StorageLocationsAdmin        types.PARTY           `json:"storageLocationsAdmin"`
	PendingStorageLocationsAdmin types.PARTY           `json:"pendingStorageLocationsAdmin"`
	SignerConfigs                types.GENMAP          `json:"signerConfigs"`
	Deps                         CommitteeVerifierDeps `json:"deps"`
	RemoteChainFeeConfigs        types.GENMAP          `json:"remoteChainFeeConfigs"`
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
	args["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocationsAdmin"] = t.StorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingStorageLocationsAdmin"] = t.PendingStorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signerConfigs"] = func() any {
		if t.SignerConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SignerConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
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
	args["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["storageLocationsAdmin"] = t.StorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["pendingStorageLocationsAdmin"] = t.PendingStorageLocationsAdmin.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["signerConfigs"] = func() any {
		if t.SignerConfigs == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.SignerConfigs}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deps"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Deps).(mapper); ok {
			return m.toMap()
		}
		return t.Deps
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

// MarshalHex encodes CommitteeVerifier to hex string (Canton MCMS format)
func (t CommitteeVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifier from hex string (Canton MCMS format)
func (t *CommitteeVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
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

// CommitteeVerifierApplySignatureConfigs exercises the CommitteeVerifier_ApplySignatureConfigs choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierApplySignatureConfigs(contractID string, args CommitteeVerifierApplySignatureConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_ApplySignatureConfigs",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierApplySignatureConfigsWithPackageID exercises the CommitteeVerifier_ApplySignatureConfigs choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierApplySignatureConfigsWithPackageID(contractID string, packageID string, args CommitteeVerifierApplySignatureConfigs) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_ApplySignatureConfigs",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierTransferStorageLocationsAdmin exercises the CommitteeVerifier_TransferStorageLocationsAdmin choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierTransferStorageLocationsAdmin(contractID string, args CommitteeVerifierTransferStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_TransferStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierTransferStorageLocationsAdminWithPackageID exercises the CommitteeVerifier_TransferStorageLocationsAdmin choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierTransferStorageLocationsAdminWithPackageID(contractID string, packageID string, args CommitteeVerifierTransferStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_TransferStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierAcceptStorageLocationsAdmin exercises the CommitteeVerifier_AcceptStorageLocationsAdmin choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierAcceptStorageLocationsAdmin(contractID string, args CommitteeVerifierAcceptStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_AcceptStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierAcceptStorageLocationsAdminWithPackageID exercises the CommitteeVerifier_AcceptStorageLocationsAdmin choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierAcceptStorageLocationsAdminWithPackageID(contractID string, packageID string, args CommitteeVerifierAcceptStorageLocationsAdmin) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_AcceptStorageLocationsAdmin",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierUpdateStorageLocations exercises the CommitteeVerifier_UpdateStorageLocations choice on this CommitteeVerifier contract
// This method uses the package name in the template ID
func (t CommitteeVerifier) CommitteeVerifierUpdateStorageLocations(contractID string, args CommitteeVerifierUpdateStorageLocations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_UpdateStorageLocations",
		Arguments:  argsToMap(args),
	}
}

// CommitteeVerifierUpdateStorageLocationsWithPackageID exercises the CommitteeVerifier_UpdateStorageLocations choice using the provided package ID instead of package name
func (t CommitteeVerifier) CommitteeVerifierUpdateStorageLocationsWithPackageID(contractID string, packageID string, args CommitteeVerifierUpdateStorageLocations) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.CommitteeVerifier", "CommitteeVerifier"),
		ContractID: contractID,
		Choice:     "CommitteeVerifier_UpdateStorageLocations",
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

// CommitteeVerifierDeps is a Record type
type CommitteeVerifierDeps struct {
	RmnRemote common.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts CommitteeVerifierDeps to a map for DAML arguments
func (t CommitteeVerifierDeps) ToMap() map[string]any {
	m := make(map[string]any)

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	return m
}

func (t CommitteeVerifierDeps) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierDeps) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierDeps to hex string (Canton MCMS format)
func (t CommitteeVerifierDeps) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierDeps from hex string (Canton MCMS format)
func (t *CommitteeVerifierDeps) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierAcceptStorageLocationsAdmin is a Record type
type CommitteeVerifierAcceptStorageLocationsAdmin struct {
}

// ToMap converts CommitteeVerifierAcceptStorageLocationsAdmin to a map for DAML arguments
func (t CommitteeVerifierAcceptStorageLocationsAdmin) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t CommitteeVerifierAcceptStorageLocationsAdmin) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierAcceptStorageLocationsAdmin) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierAcceptStorageLocationsAdmin to hex string (Canton MCMS format)
func (t CommitteeVerifierAcceptStorageLocationsAdmin) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierAcceptStorageLocationsAdmin from hex string (Canton MCMS format)
func (t *CommitteeVerifierAcceptStorageLocationsAdmin) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierApplySignatureConfigs is a Record type
type CommitteeVerifierApplySignatureConfigs struct {
	SourceChainSelectorsToRemove []types.NUMERIC   `json:"sourceChainSelectorsToRemove"`
	SignatureConfigs             []SignatureConfig `json:"signatureConfigs"`
}

// ToMap converts CommitteeVerifierApplySignatureConfigs to a map for DAML arguments
func (t CommitteeVerifierApplySignatureConfigs) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelectorsToRemove"] = func() []any {
		res := make([]any, 0, len(t.SourceChainSelectorsToRemove))
		for _, e := range t.SourceChainSelectorsToRemove {
			res = append(res, e)
		}
		return res
	}()

	m["signatureConfigs"] = func() []any {
		res := make([]any, 0, len(t.SignatureConfigs))
		for _, e := range t.SignatureConfigs {
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

func (t CommitteeVerifierApplySignatureConfigs) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierApplySignatureConfigs) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierApplySignatureConfigs to hex string (Canton MCMS format)
func (t CommitteeVerifierApplySignatureConfigs) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierApplySignatureConfigs from hex string (Canton MCMS format)
func (t *CommitteeVerifierApplySignatureConfigs) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

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

// MarshalHex encodes CommitteeVerifierCalculateFee to hex string (Canton MCMS format)
func (t CommitteeVerifierCalculateFee) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierCalculateFee from hex string (Canton MCMS format)
func (t *CommitteeVerifierCalculateFee) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierCalculateFeeMCMSParams is CommitteeVerifierCalculateFee without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CommitteeVerifierCalculateFeeMCMSParams struct {
	SendingMessageCid types.CONTRACT_ID `json:"sendingMessageCid"`
}

// MarshalHex encodes CommitteeVerifierCalculateFeeMCMSParams to hex string for MCMS operationData.
func (t CommitteeVerifierCalculateFeeMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierCalculateFeeMCMSParams from hex string.
func (t *CommitteeVerifierCalculateFeeMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierForwardToVerifier is a Record type
type CommitteeVerifierForwardToVerifier struct {
	RmnRemoteCid      types.CONTRACT_ID  `json:"rmnRemoteCid"`
	ExtraContext      common.CCIPContext `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID  `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT         `json:"verifierArgs"`
	Caller            types.PARTY        `json:"caller"`
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

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
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

// MarshalHex encodes CommitteeVerifierForwardToVerifier to hex string (Canton MCMS format)
func (t CommitteeVerifierForwardToVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierForwardToVerifier from hex string (Canton MCMS format)
func (t *CommitteeVerifierForwardToVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierForwardToVerifierMCMSParams is CommitteeVerifierForwardToVerifier without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CommitteeVerifierForwardToVerifierMCMSParams struct {
	RmnRemoteCid      types.CONTRACT_ID  `json:"rmnRemoteCid"`
	ExtraContext      common.CCIPContext `json:"extraContext"`
	SendingMessageCid types.CONTRACT_ID  `json:"sendingMessageCid"`
	VerifierArgs      types.TEXT         `json:"verifierArgs"`
}

// MarshalHex encodes CommitteeVerifierForwardToVerifierMCMSParams to hex string for MCMS operationData.
func (t CommitteeVerifierForwardToVerifierMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierForwardToVerifierMCMSParams from hex string.
func (t *CommitteeVerifierForwardToVerifierMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierTransferStorageLocationsAdmin is a Record type
type CommitteeVerifierTransferStorageLocationsAdmin struct {
	NewAdmin types.PARTY `json:"newAdmin"`
}

// ToMap converts CommitteeVerifierTransferStorageLocationsAdmin to a map for DAML arguments
func (t CommitteeVerifierTransferStorageLocationsAdmin) ToMap() map[string]any {
	m := make(map[string]any)

	m["newAdmin"] = t.NewAdmin.ToMap()

	return m
}

func (t CommitteeVerifierTransferStorageLocationsAdmin) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierTransferStorageLocationsAdmin) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierTransferStorageLocationsAdmin to hex string (Canton MCMS format)
func (t CommitteeVerifierTransferStorageLocationsAdmin) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierTransferStorageLocationsAdmin from hex string (Canton MCMS format)
func (t *CommitteeVerifierTransferStorageLocationsAdmin) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierUpdateStorageLocations is a Record type
type CommitteeVerifierUpdateStorageLocations struct {
	NewLocations []types.TEXT `json:"newLocations"`
}

// ToMap converts CommitteeVerifierUpdateStorageLocations to a map for DAML arguments
func (t CommitteeVerifierUpdateStorageLocations) ToMap() map[string]any {
	m := make(map[string]any)

	m["newLocations"] = func() []any {
		res := make([]any, 0, len(t.NewLocations))
		for _, e := range t.NewLocations {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t CommitteeVerifierUpdateStorageLocations) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CommitteeVerifierUpdateStorageLocations) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CommitteeVerifierUpdateStorageLocations to hex string (Canton MCMS format)
func (t CommitteeVerifierUpdateStorageLocations) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierUpdateStorageLocations from hex string (Canton MCMS format)
func (t *CommitteeVerifierUpdateStorageLocations) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierVerifyMessage is a Record type
type CommitteeVerifierVerifyMessage struct {
	RmnRemoteCid        types.CONTRACT_ID  `json:"rmnRemoteCid"`
	ExtraContext        common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID  `json:"executingMessageCid"`
	VerifierResults     types.TEXT         `json:"verifierResults"`
	Caller              types.PARTY        `json:"caller"`
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

	m["extraContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ExtraContext).(mapper); ok {
			return m.toMap()
		}
		return t.ExtraContext
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

// MarshalHex encodes CommitteeVerifierVerifyMessage to hex string (Canton MCMS format)
func (t CommitteeVerifierVerifyMessage) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierVerifyMessage from hex string (Canton MCMS format)
func (t *CommitteeVerifierVerifyMessage) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// CommitteeVerifierVerifyMessageMCMSParams is CommitteeVerifierVerifyMessage without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type CommitteeVerifierVerifyMessageMCMSParams struct {
	RmnRemoteCid        types.CONTRACT_ID  `json:"rmnRemoteCid"`
	ExtraContext        common.CCIPContext `json:"extraContext"`
	ExecutingMessageCid types.CONTRACT_ID  `json:"executingMessageCid"`
	VerifierResults     types.TEXT         `json:"verifierResults"`
}

// MarshalHex encodes CommitteeVerifierVerifyMessageMCMSParams to hex string for MCMS operationData.
func (t CommitteeVerifierVerifyMessageMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CommitteeVerifierVerifyMessageMCMSParams from hex string.
func (t *CommitteeVerifierVerifyMessageMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SignatureConfig is a Record type
type SignatureConfig struct {
	SourceChainSelector types.NUMERIC `json:"sourceChainSelector"`
	Threshold           types.INT64   `json:"threshold"`
	SignerKeys          []types.TEXT  `json:"signerKeys"`
}

// ToMap converts SignatureConfig to a map for DAML arguments
func (t SignatureConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["sourceChainSelector"] = t.SourceChainSelector

	m["threshold"] = int64(t.Threshold)

	m["signerKeys"] = func() []any {
		res := make([]any, 0, len(t.SignerKeys))
		for _, e := range t.SignerKeys {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t SignatureConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SignatureConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SignatureConfig to hex string (Canton MCMS format)
func (t SignatureConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SignatureConfig from hex string (Canton MCMS format)
func (t *SignatureConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	CommitteeVerifierAcceptStorageLocationsAdmin(args CommitteeVerifierAcceptStorageLocationsAdmin) (*bind.EncodedChoice, error)
	CommitteeVerifierApplySignatureConfigs(args CommitteeVerifierApplySignatureConfigs) (*bind.EncodedChoice, error)
	CommitteeVerifierCalculateFee(args CommitteeVerifierCalculateFee) (*bind.EncodedChoice, error)
	CommitteeVerifierCalculateFeeMCMSParams(args CommitteeVerifierCalculateFeeMCMSParams) (*bind.EncodedChoice, error)
	CommitteeVerifierForwardToVerifier(args CommitteeVerifierForwardToVerifier) (*bind.EncodedChoice, error)
	CommitteeVerifierForwardToVerifierMCMSParams(args CommitteeVerifierForwardToVerifierMCMSParams) (*bind.EncodedChoice, error)
	CommitteeVerifierTransferStorageLocationsAdmin(args CommitteeVerifierTransferStorageLocationsAdmin) (*bind.EncodedChoice, error)
	CommitteeVerifierUpdateStorageLocations(args CommitteeVerifierUpdateStorageLocations) (*bind.EncodedChoice, error)
	CommitteeVerifierVerifyMessage(args CommitteeVerifierVerifyMessage) (*bind.EncodedChoice, error)
	CommitteeVerifierVerifyMessageMCMSParams(args CommitteeVerifierVerifyMessageMCMSParams) (*bind.EncodedChoice, error)
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

// CommitteeVerifierAcceptStorageLocationsAdmin encodes parameters for the CommitteeVerifierAcceptStorageLocationsAdmin choice.
func (e *encoder) CommitteeVerifierAcceptStorageLocationsAdmin(args CommitteeVerifierAcceptStorageLocationsAdmin) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierAcceptStorageLocationsAdmin", args)
}

// CommitteeVerifierApplySignatureConfigs encodes parameters for the CommitteeVerifierApplySignatureConfigs choice.
func (e *encoder) CommitteeVerifierApplySignatureConfigs(args CommitteeVerifierApplySignatureConfigs) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierApplySignatureConfigs", args)
}

// CommitteeVerifierCalculateFee encodes parameters for the CommitteeVerifierCalculateFee choice.
func (e *encoder) CommitteeVerifierCalculateFee(args CommitteeVerifierCalculateFee) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierCalculateFee", args)
}

// CommitteeVerifierCalculateFeeMCMSParams encodes MCMS parameters (without Caller) for the CommitteeVerifierCalculateFee choice.
func (e *encoder) CommitteeVerifierCalculateFeeMCMSParams(args CommitteeVerifierCalculateFeeMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierCalculateFee", args)
}

// CommitteeVerifierForwardToVerifier encodes parameters for the CommitteeVerifierForwardToVerifier choice.
func (e *encoder) CommitteeVerifierForwardToVerifier(args CommitteeVerifierForwardToVerifier) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierForwardToVerifier", args)
}

// CommitteeVerifierForwardToVerifierMCMSParams encodes MCMS parameters (without Caller) for the CommitteeVerifierForwardToVerifier choice.
func (e *encoder) CommitteeVerifierForwardToVerifierMCMSParams(args CommitteeVerifierForwardToVerifierMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierForwardToVerifier", args)
}

// CommitteeVerifierTransferStorageLocationsAdmin encodes parameters for the CommitteeVerifierTransferStorageLocationsAdmin choice.
func (e *encoder) CommitteeVerifierTransferStorageLocationsAdmin(args CommitteeVerifierTransferStorageLocationsAdmin) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierTransferStorageLocationsAdmin", args)
}

// CommitteeVerifierUpdateStorageLocations encodes parameters for the CommitteeVerifierUpdateStorageLocations choice.
func (e *encoder) CommitteeVerifierUpdateStorageLocations(args CommitteeVerifierUpdateStorageLocations) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierUpdateStorageLocations", args)
}

// CommitteeVerifierVerifyMessage encodes parameters for the CommitteeVerifierVerifyMessage choice.
func (e *encoder) CommitteeVerifierVerifyMessage(args CommitteeVerifierVerifyMessage) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierVerifyMessage", args)
}

// CommitteeVerifierVerifyMessageMCMSParams encodes MCMS parameters (without Caller) for the CommitteeVerifierVerifyMessage choice.
func (e *encoder) CommitteeVerifierVerifyMessageMCMSParams(args CommitteeVerifierVerifyMessageMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("CommitteeVerifierVerifyMessage", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
