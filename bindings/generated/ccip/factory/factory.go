package factory

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ccipreceiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	ccipsender "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	ccvs "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	feequoter "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	lockreleasetokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	offramp "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	onramp "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	perpartyrouter "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	rmn "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	tokenadminregistry "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
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
	PackageName = "ccip-factory"
	PackageID   = "88fd99add287be5ada6b41690180f7339a8a090a3e289c06e8f574a250472484"
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

// CCIPFactory is a Template type
type CCIPFactory struct {
	InstanceId                    types.TEXT   `json:"instanceId"`
	Owner                         types.PARTY  `json:"owner"`
	McmsParty                     types.PARTY  `json:"mcmsParty"`
	UsedInstanceIds               types.GENMAP `json:"usedInstanceIds"`
	DeployedContracts             types.GENMAP `json:"deployedContracts"`
	PerPartyRouterFactoryDeployed types.BOOL   `json:"perPartyRouterFactoryDeployed"`
}

// GetTemplateID returns the template ID for this template using the package name
func (t CCIPFactory) GetTemplateID() string {
	return fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory")
}

// GetTemplateIDWithPackageID returns the template ID using the provided package ID instead of package name
func (t CCIPFactory) GetTemplateIDWithPackageID(packageID string) string {
	return fmt.Sprintf("%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory")
}

// CreateCommand returns a CreateCommand for this template using the package name
func (t CCIPFactory) CreateCommand() *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsParty"] = t.McmsParty.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usedInstanceIds"] = func() any {
		if t.UsedInstanceIds == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsedInstanceIds}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deployedContracts"] = func() any {
		if t.DeployedContracts == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DeployedContracts}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["perPartyRouterFactoryDeployed"] = bool(t.PerPartyRouterFactoryDeployed)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateID(),
		Arguments:  args,
	}
}

// CreateCommandWithPackageID returns a CreateCommand using the provided package ID instead of package name
func (t CCIPFactory) CreateCommandWithPackageID(packageID string) *model.CreateCommand {
	args := make(map[string]any)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["instanceId"] = string(t.InstanceId)

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["owner"] = t.Owner.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["mcmsParty"] = t.McmsParty.ToMap()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["usedInstanceIds"] = func() any {
		if t.UsedInstanceIds == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsedInstanceIds}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["deployedContracts"] = func() any {
		if t.DeployedContracts == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DeployedContracts}
	}()

	// IMPORTANT: always include non-optional fields (GENMAP/MAP/LIST/[] etc), even if empty
	args["perPartyRouterFactoryDeployed"] = bool(t.PerPartyRouterFactoryDeployed)

	return &model.CreateCommand{
		TemplateID: t.GetTemplateIDWithPackageID(packageID),
		Arguments:  args,
	}
}

func (t CCIPFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *CCIPFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes CCIPFactory to hex string (Canton MCMS format)
func (t CCIPFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes CCIPFactory from hex string (Canton MCMS format)
func (t *CCIPFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// Choice methods for CCIPFactory

// DeployGlobalConfig exercises the DeployGlobalConfig choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployGlobalConfig(contractID string, args DeployGlobalConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployGlobalConfig",
		Arguments:  argsToMap(args),
	}
}

// DeployGlobalConfigWithPackageID exercises the DeployGlobalConfig choice using the provided package ID instead of package name
func (t CCIPFactory) DeployGlobalConfigWithPackageID(contractID string, packageID string, args DeployGlobalConfig) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployGlobalConfig",
		Arguments:  argsToMap(args),
	}
}

// DeployFeeQuoter exercises the DeployFeeQuoter choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployFeeQuoter(contractID string, args DeployFeeQuoter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployFeeQuoter",
		Arguments:  argsToMap(args),
	}
}

// DeployFeeQuoterWithPackageID exercises the DeployFeeQuoter choice using the provided package ID instead of package name
func (t CCIPFactory) DeployFeeQuoterWithPackageID(contractID string, packageID string, args DeployFeeQuoter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployFeeQuoter",
		Arguments:  argsToMap(args),
	}
}

// DeployTokenAdminRegistry exercises the DeployTokenAdminRegistry choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployTokenAdminRegistry(contractID string, args DeployTokenAdminRegistry) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployTokenAdminRegistry",
		Arguments:  argsToMap(args),
	}
}

// DeployTokenAdminRegistryWithPackageID exercises the DeployTokenAdminRegistry choice using the provided package ID instead of package name
func (t CCIPFactory) DeployTokenAdminRegistryWithPackageID(contractID string, packageID string, args DeployTokenAdminRegistry) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployTokenAdminRegistry",
		Arguments:  argsToMap(args),
	}
}

// DeployOnRamp exercises the DeployOnRamp choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployOnRamp(contractID string, args DeployOnRamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployOnRamp",
		Arguments:  argsToMap(args),
	}
}

// DeployOnRampWithPackageID exercises the DeployOnRamp choice using the provided package ID instead of package name
func (t CCIPFactory) DeployOnRampWithPackageID(contractID string, packageID string, args DeployOnRamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployOnRamp",
		Arguments:  argsToMap(args),
	}
}

// DeployOffRamp exercises the DeployOffRamp choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployOffRamp(contractID string, args DeployOffRamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployOffRamp",
		Arguments:  argsToMap(args),
	}
}

// DeployOffRampWithPackageID exercises the DeployOffRamp choice using the provided package ID instead of package name
func (t CCIPFactory) DeployOffRampWithPackageID(contractID string, packageID string, args DeployOffRamp) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployOffRamp",
		Arguments:  argsToMap(args),
	}
}

// DeployPerPartyRouterFactory exercises the DeployPerPartyRouterFactory choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployPerPartyRouterFactory(contractID string, args DeployPerPartyRouterFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployPerPartyRouterFactory",
		Arguments:  argsToMap(args),
	}
}

// DeployPerPartyRouterFactoryWithPackageID exercises the DeployPerPartyRouterFactory choice using the provided package ID instead of package name
func (t CCIPFactory) DeployPerPartyRouterFactoryWithPackageID(contractID string, packageID string, args DeployPerPartyRouterFactory) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployPerPartyRouterFactory",
		Arguments:  argsToMap(args),
	}
}

// DeployRMNRemote exercises the DeployRMNRemote choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployRMNRemote(contractID string, args DeployRMNRemote) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployRMNRemote",
		Arguments:  argsToMap(args),
	}
}

// DeployRMNRemoteWithPackageID exercises the DeployRMNRemote choice using the provided package ID instead of package name
func (t CCIPFactory) DeployRMNRemoteWithPackageID(contractID string, packageID string, args DeployRMNRemote) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployRMNRemote",
		Arguments:  argsToMap(args),
	}
}

// DeployCommitteeVerifier exercises the DeployCommitteeVerifier choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployCommitteeVerifier(contractID string, args DeployCommitteeVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCommitteeVerifier",
		Arguments:  argsToMap(args),
	}
}

// DeployCommitteeVerifierWithPackageID exercises the DeployCommitteeVerifier choice using the provided package ID instead of package name
func (t CCIPFactory) DeployCommitteeVerifierWithPackageID(contractID string, packageID string, args DeployCommitteeVerifier) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCommitteeVerifier",
		Arguments:  argsToMap(args),
	}
}

// DeployLockReleaseTokenPool exercises the DeployLockReleaseTokenPool choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployLockReleaseTokenPool(contractID string, args DeployLockReleaseTokenPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployLockReleaseTokenPool",
		Arguments:  argsToMap(args),
	}
}

// DeployLockReleaseTokenPoolWithPackageID exercises the DeployLockReleaseTokenPool choice using the provided package ID instead of package name
func (t CCIPFactory) DeployLockReleaseTokenPoolWithPackageID(contractID string, packageID string, args DeployLockReleaseTokenPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployLockReleaseTokenPool",
		Arguments:  argsToMap(args),
	}
}

// DeployRateLimiter exercises the DeployRateLimiter choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployRateLimiter(contractID string, args DeployRateLimiter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployRateLimiter",
		Arguments:  argsToMap(args),
	}
}

// DeployRateLimiterWithPackageID exercises the DeployRateLimiter choice using the provided package ID instead of package name
func (t CCIPFactory) DeployRateLimiterWithPackageID(contractID string, packageID string, args DeployRateLimiter) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployRateLimiter",
		Arguments:  argsToMap(args),
	}
}

// DeployCCIPSender exercises the DeployCCIPSender choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployCCIPSender(contractID string, args DeployCCIPSender) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCCIPSender",
		Arguments:  argsToMap(args),
	}
}

// DeployCCIPSenderWithPackageID exercises the DeployCCIPSender choice using the provided package ID instead of package name
func (t CCIPFactory) DeployCCIPSenderWithPackageID(contractID string, packageID string, args DeployCCIPSender) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCCIPSender",
		Arguments:  argsToMap(args),
	}
}

// DeployCCIPReceiver exercises the DeployCCIPReceiver choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployCCIPReceiver(contractID string, args DeployCCIPReceiver) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCCIPReceiver",
		Arguments:  argsToMap(args),
	}
}

// DeployCCIPReceiverWithPackageID exercises the DeployCCIPReceiver choice using the provided package ID instead of package name
func (t CCIPFactory) DeployCCIPReceiverWithPackageID(contractID string, packageID string, args DeployCCIPReceiver) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployCCIPReceiver",
		Arguments:  argsToMap(args),
	}
}

// SetOwnerToMCMS exercises the SetOwnerToMCMS choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) SetOwnerToMCMS(contractID string, args SetOwnerToMCMS) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "SetOwnerToMCMS",
		Arguments:  argsToMap(args),
	}
}

// SetOwnerToMCMSWithPackageID exercises the SetOwnerToMCMS choice using the provided package ID instead of package name
func (t CCIPFactory) SetOwnerToMCMSWithPackageID(contractID string, packageID string, args SetOwnerToMCMS) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "SetOwnerToMCMS",
		Arguments:  argsToMap(args),
	}
}

// Archive exercises the Archive choice on this CCIPFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CCIPFactory) Archive(contractID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// ArchiveWithPackageID exercises the Archive choice using the provided package ID instead of package name
func (t CCIPFactory) ArchiveWithPackageID(contractID string, packageID string) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "Archive",
		Arguments:  map[string]any{},
	}
}

// GetFactoryState exercises the GetFactoryState choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) GetFactoryState(contractID string, args GetFactoryState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "GetFactoryState",
		Arguments:  argsToMap(args),
	}
}

// GetFactoryStateWithPackageID exercises the GetFactoryState choice using the provided package ID instead of package name
func (t CCIPFactory) GetFactoryStateWithPackageID(contractID string, packageID string, args GetFactoryState) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "GetFactoryState",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypoint exercises the MCMSReceiver_Entrypoint choice on this CCIPFactory contract via the IMCMSReceiver interface
// This method uses the package name in the template ID
func (t CCIPFactory) MCMSReceiverEntrypoint(contractID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CCIPFactory) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args mcms.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CCIPFactory

var _ mcms.IMCMSReceiver = (*CCIPFactory)(nil)

// DeployCCIPReceiver is a Record type
type DeployCCIPReceiver struct {
	Contract ccipreceiver.CCIPReceiver `json:"contract"`
}

// ToMap converts DeployCCIPReceiver to a map for DAML arguments
func (t DeployCCIPReceiver) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployCCIPReceiver) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCCIPReceiver) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCCIPReceiver to hex string (Canton MCMS format)
func (t DeployCCIPReceiver) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCCIPReceiver from hex string (Canton MCMS format)
func (t *DeployCCIPReceiver) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCCIPReceiverParams is a Record type
type DeployCCIPReceiverParams struct {
	InstanceId             types.TEXT                `json:"instanceId"`
	Owner                  types.PARTY               `json:"owner"`
	RequiredCCVs           []mcms.RawInstanceAddress `json:"requiredCCVs"`
	OptionalCCVs           []mcms.RawInstanceAddress `json:"optionalCCVs"`
	OptionalThreshold      types.INT64               `json:"optionalThreshold"`
	ReceiverFinalityConfig common.FinalityConfig     `json:"receiverFinalityConfig"`
}

// ToMap converts DeployCCIPReceiverParams to a map for DAML arguments
func (t DeployCCIPReceiverParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	m["requiredCCVs"] = func() []any {
		res := make([]any, 0, len(t.RequiredCCVs))
		for _, e := range t.RequiredCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["optionalCCVs"] = func() []any {
		res := make([]any, 0, len(t.OptionalCCVs))
		for _, e := range t.OptionalCCVs {
			type mapper interface{ toMap() map[string]any }
			if m, ok := any(e).(mapper); ok {
				res = append(res, m.toMap())
			} else {
				res = append(res, e)
			}
		}
		return res
	}()

	m["optionalThreshold"] = int64(t.OptionalThreshold)

	m["receiverFinalityConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.ReceiverFinalityConfig).(mapper); ok {
			return m.toMap()
		}
		return t.ReceiverFinalityConfig
	}()

	return m
}

func (t DeployCCIPReceiverParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCCIPReceiverParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCCIPReceiverParams to hex string (Canton MCMS format)
func (t DeployCCIPReceiverParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCCIPReceiverParams from hex string (Canton MCMS format)
func (t *DeployCCIPReceiverParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCCIPSender is a Record type
type DeployCCIPSender struct {
	Contract ccipsender.CCIPSender `json:"contract"`
}

// ToMap converts DeployCCIPSender to a map for DAML arguments
func (t DeployCCIPSender) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployCCIPSender) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCCIPSender) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCCIPSender to hex string (Canton MCMS format)
func (t DeployCCIPSender) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCCIPSender from hex string (Canton MCMS format)
func (t *DeployCCIPSender) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCCIPSenderParams is a Record type
type DeployCCIPSenderParams struct {
	InstanceId types.TEXT  `json:"instanceId"`
	Owner      types.PARTY `json:"owner"`
}

// ToMap converts DeployCCIPSenderParams to a map for DAML arguments
func (t DeployCCIPSenderParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	return m
}

func (t DeployCCIPSenderParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCCIPSenderParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCCIPSenderParams to hex string (Canton MCMS format)
func (t DeployCCIPSenderParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCCIPSenderParams from hex string (Canton MCMS format)
func (t *DeployCCIPSenderParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCommitteeVerifier is a Record type
type DeployCommitteeVerifier struct {
	Contract ccvs.CommitteeVerifier `json:"contract"`
}

// ToMap converts DeployCommitteeVerifier to a map for DAML arguments
func (t DeployCommitteeVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployCommitteeVerifier) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCommitteeVerifier) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCommitteeVerifier to hex string (Canton MCMS format)
func (t DeployCommitteeVerifier) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCommitteeVerifier from hex string (Canton MCMS format)
func (t *DeployCommitteeVerifier) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCommitteeVerifierParams is a Record type
type DeployCommitteeVerifierParams struct {
	InstanceId                   types.TEXT              `json:"instanceId"`
	Owner                        types.PARTY             `json:"owner"`
	CcipOwner                    types.PARTY             `json:"ccipOwner"`
	VersionTag                   types.TEXT              `json:"versionTag"`
	AllowListAdmin               *types.PARTY            `json:"allowListAdmin" hex:"optional"`
	MessageSentObservers         []types.PARTY           `json:"messageSentObservers"`
	RmnRemote                    mcms.RawInstanceAddress `json:"rmnRemote"`
	StorageLocations             []types.TEXT            `json:"storageLocations"`
	StorageLocationsAdmin        types.PARTY             `json:"storageLocationsAdmin"`
	PendingStorageLocationsAdmin types.PARTY             `json:"pendingStorageLocationsAdmin"`
}

// ToMap converts DeployCommitteeVerifierParams to a map for DAML arguments
func (t DeployCommitteeVerifierParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["versionTag"] = string(t.VersionTag)

	if t.AllowListAdmin != nil {
		m["allowListAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.AllowListAdmin).ToMap(),
		}
	} else {
		m["allowListAdmin"] = map[string]any{
			"_type": "optional",
		}
	}

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["storageLocations"] = func() []any {
		res := make([]any, 0, len(t.StorageLocations))
		for _, e := range t.StorageLocations {
			res = append(res, string(e))
		}
		return res
	}()

	m["storageLocationsAdmin"] = t.StorageLocationsAdmin.ToMap()

	m["pendingStorageLocationsAdmin"] = t.PendingStorageLocationsAdmin.ToMap()

	return m
}

func (t DeployCommitteeVerifierParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployCommitteeVerifierParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployCommitteeVerifierParams to hex string (Canton MCMS format)
func (t DeployCommitteeVerifierParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployCommitteeVerifierParams from hex string (Canton MCMS format)
func (t *DeployCommitteeVerifierParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployFeeQuoter is a Record type
type DeployFeeQuoter struct {
	Contract feequoter.FeeQuoter `json:"contract"`
}

// ToMap converts DeployFeeQuoter to a map for DAML arguments
func (t DeployFeeQuoter) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployFeeQuoter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployFeeQuoter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployFeeQuoter to hex string (Canton MCMS format)
func (t DeployFeeQuoter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployFeeQuoter from hex string (Canton MCMS format)
func (t *DeployFeeQuoter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployFeeQuoterParams is a Record type
type DeployFeeQuoterParams struct {
	InstanceId            types.TEXT                               `json:"instanceId"`
	LinkTokenInstrumentId splice_api_token_holding_v1.InstrumentId `json:"linkTokenInstrumentId"`
}

// ToMap converts DeployFeeQuoterParams to a map for DAML arguments
func (t DeployFeeQuoterParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["linkTokenInstrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.LinkTokenInstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.LinkTokenInstrumentId
	}()

	return m
}

func (t DeployFeeQuoterParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployFeeQuoterParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployFeeQuoterParams to hex string (Canton MCMS format)
func (t DeployFeeQuoterParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployFeeQuoterParams from hex string (Canton MCMS format)
func (t *DeployFeeQuoterParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployGlobalConfig is a Record type
type DeployGlobalConfig struct {
	Contract common.GlobalConfig `json:"contract"`
}

// ToMap converts DeployGlobalConfig to a map for DAML arguments
func (t DeployGlobalConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployGlobalConfig) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployGlobalConfig) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployGlobalConfig to hex string (Canton MCMS format)
func (t DeployGlobalConfig) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployGlobalConfig from hex string (Canton MCMS format)
func (t *DeployGlobalConfig) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployGlobalConfigParams is a Record type
type DeployGlobalConfigParams struct {
	InstanceId    types.TEXT    `json:"instanceId"`
	ChainSelector types.NUMERIC `json:"chainSelector"`
}

// ToMap converts DeployGlobalConfigParams to a map for DAML arguments
func (t DeployGlobalConfigParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["chainSelector"] = t.ChainSelector

	return m
}

func (t DeployGlobalConfigParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployGlobalConfigParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployGlobalConfigParams to hex string (Canton MCMS format)
func (t DeployGlobalConfigParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployGlobalConfigParams from hex string (Canton MCMS format)
func (t *DeployGlobalConfigParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployLockReleaseTokenPool is a Record type
type DeployLockReleaseTokenPool struct {
	Contract lockreleasetokenpool.LockReleaseTokenPool `json:"contract"`
}

// ToMap converts DeployLockReleaseTokenPool to a map for DAML arguments
func (t DeployLockReleaseTokenPool) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployLockReleaseTokenPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployLockReleaseTokenPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployLockReleaseTokenPool to hex string (Canton MCMS format)
func (t DeployLockReleaseTokenPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployLockReleaseTokenPool from hex string (Canton MCMS format)
func (t *DeployLockReleaseTokenPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployLockReleaseTokenPoolParams is a Record type
type DeployLockReleaseTokenPoolParams struct {
	InstanceId         types.TEXT                               `json:"instanceId"`
	PoolOwner          types.PARTY                              `json:"poolOwner"`
	CcipOwner          types.PARTY                              `json:"ccipOwner"`
	InstrumentId       splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
	Decimals           types.INT64                              `json:"decimals"`
	RateLimitAdmin     *types.PARTY                             `json:"rateLimitAdmin" hex:"optional"`
	TokenAdminRegistry mcms.RawInstanceAddress                  `json:"tokenAdminRegistry"`
	FeeQuoter          mcms.RawInstanceAddress                  `json:"feeQuoter"`
	RmnRemote          mcms.RawInstanceAddress                  `json:"rmnRemote"`
	PoolReceiveContext common.CCIPContext                       `json:"poolReceiveContext"`
	TransferTimeout    lockreleasetokenpool.TransferTimeout     `json:"transferTimeout"`
}

// ToMap converts DeployLockReleaseTokenPoolParams to a map for DAML arguments
func (t DeployLockReleaseTokenPoolParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instrumentId"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.InstrumentId).(mapper); ok {
			return m.toMap()
		}
		return t.InstrumentId
	}()

	m["decimals"] = int64(t.Decimals)

	if t.RateLimitAdmin != nil {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
		}
	}

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

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	m["poolReceiveContext"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.PoolReceiveContext).(mapper); ok {
			return m.toMap()
		}
		return t.PoolReceiveContext
	}()

	m["transferTimeout"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.TransferTimeout).(mapper); ok {
			return m.toMap()
		}
		return t.TransferTimeout
	}()

	return m
}

func (t DeployLockReleaseTokenPoolParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployLockReleaseTokenPoolParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployLockReleaseTokenPoolParams to hex string (Canton MCMS format)
func (t DeployLockReleaseTokenPoolParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployLockReleaseTokenPoolParams from hex string (Canton MCMS format)
func (t *DeployLockReleaseTokenPoolParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployOffRamp is a Record type
type DeployOffRamp struct {
	Contract offramp.OffRamp `json:"contract"`
}

// ToMap converts DeployOffRamp to a map for DAML arguments
func (t DeployOffRamp) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployOffRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployOffRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployOffRamp to hex string (Canton MCMS format)
func (t DeployOffRamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployOffRamp from hex string (Canton MCMS format)
func (t *DeployOffRamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployOffRampParams is a Record type
type DeployOffRampParams struct {
	InstanceId         types.TEXT              `json:"instanceId"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts DeployOffRampParams to a map for DAML arguments
func (t DeployOffRampParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

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

	return m
}

func (t DeployOffRampParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployOffRampParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployOffRampParams to hex string (Canton MCMS format)
func (t DeployOffRampParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployOffRampParams from hex string (Canton MCMS format)
func (t *DeployOffRampParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployOnRamp is a Record type
type DeployOnRamp struct {
	Contract onramp.OnRamp `json:"contract"`
}

// ToMap converts DeployOnRamp to a map for DAML arguments
func (t DeployOnRamp) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployOnRamp) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployOnRamp) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployOnRamp to hex string (Canton MCMS format)
func (t DeployOnRamp) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployOnRamp from hex string (Canton MCMS format)
func (t *DeployOnRamp) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployOnRampParams is a Record type
type DeployOnRampParams struct {
	InstanceId         types.TEXT              `json:"instanceId"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          mcms.RawInstanceAddress `json:"feeQuoter"`
	CcvRegistry        mcms.RawInstanceAddress `json:"ccvRegistry"`
	MaxUSDCentsPerMsg  types.NUMERIC           `json:"maxUSDCentsPerMsg"`
}

// ToMap converts DeployOnRampParams to a map for DAML arguments
func (t DeployOnRampParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

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

	m["maxUSDCentsPerMsg"] = t.MaxUSDCentsPerMsg

	return m
}

func (t DeployOnRampParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployOnRampParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployOnRampParams to hex string (Canton MCMS format)
func (t DeployOnRampParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployOnRampParams from hex string (Canton MCMS format)
func (t *DeployOnRampParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployPerPartyRouterFactory is a Record type
type DeployPerPartyRouterFactory struct {
	Contract perpartyrouter.PerPartyRouterFactory `json:"contract"`
}

// ToMap converts DeployPerPartyRouterFactory to a map for DAML arguments
func (t DeployPerPartyRouterFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployPerPartyRouterFactory) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployPerPartyRouterFactory) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployPerPartyRouterFactory to hex string (Canton MCMS format)
func (t DeployPerPartyRouterFactory) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployPerPartyRouterFactory from hex string (Canton MCMS format)
func (t *DeployPerPartyRouterFactory) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployPerPartyRouterFactoryParams is a Record type
type DeployPerPartyRouterFactoryParams struct {
	InstanceId         types.TEXT              `json:"instanceId"`
	OnRamp             mcms.RawInstanceAddress `json:"onRamp"`
	OffRamp            mcms.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       mcms.RawInstanceAddress `json:"globalConfig"`
	TokenAdminRegistry mcms.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          mcms.RawInstanceAddress `json:"feeQuoter"`
	RmnRemote          mcms.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts DeployPerPartyRouterFactoryParams to a map for DAML arguments
func (t DeployPerPartyRouterFactoryParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["onRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OnRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OnRamp
	}()

	m["offRamp"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.OffRamp).(mapper); ok {
			return m.toMap()
		}
		return t.OffRamp
	}()

	m["globalConfig"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.GlobalConfig).(mapper); ok {
			return m.toMap()
		}
		return t.GlobalConfig
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

	m["rmnRemote"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.RmnRemote).(mapper); ok {
			return m.toMap()
		}
		return t.RmnRemote
	}()

	return m
}

func (t DeployPerPartyRouterFactoryParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployPerPartyRouterFactoryParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployPerPartyRouterFactoryParams to hex string (Canton MCMS format)
func (t DeployPerPartyRouterFactoryParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployPerPartyRouterFactoryParams from hex string (Canton MCMS format)
func (t *DeployPerPartyRouterFactoryParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployRMNRemote is a Record type
type DeployRMNRemote struct {
	Contract rmn.RMNRemote `json:"contract"`
}

// ToMap converts DeployRMNRemote to a map for DAML arguments
func (t DeployRMNRemote) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployRMNRemote) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployRMNRemote) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployRMNRemote to hex string (Canton MCMS format)
func (t DeployRMNRemote) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployRMNRemote from hex string (Canton MCMS format)
func (t *DeployRMNRemote) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployRMNRemoteParams is a Record type
type DeployRMNRemoteParams struct {
	InstanceId      types.TEXT    `json:"instanceId"`
	RmnOwner        types.PARTY   `json:"rmnOwner"`
	CcipOwner       types.PARTY   `json:"ccipOwner"`
	CustomObservers []types.PARTY `json:"customObservers"`
	CursedSubjects  []types.TEXT  `json:"cursedSubjects"`
}

// ToMap converts DeployRMNRemoteParams to a map for DAML arguments
func (t DeployRMNRemoteParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["rmnOwner"] = t.RmnOwner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["customObservers"] = func() []any {
		res := make([]any, 0, len(t.CustomObservers))
		for _, e := range t.CustomObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["cursedSubjects"] = func() []any {
		res := make([]any, 0, len(t.CursedSubjects))
		for _, e := range t.CursedSubjects {
			res = append(res, string(e))
		}
		return res
	}()

	return m
}

func (t DeployRMNRemoteParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployRMNRemoteParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployRMNRemoteParams to hex string (Canton MCMS format)
func (t DeployRMNRemoteParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployRMNRemoteParams from hex string (Canton MCMS format)
func (t *DeployRMNRemoteParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployRateLimiter is a Record type
type DeployRateLimiter struct {
	Contract common.RateLimiter `json:"contract"`
}

// ToMap converts DeployRateLimiter to a map for DAML arguments
func (t DeployRateLimiter) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployRateLimiter) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployRateLimiter) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployRateLimiter to hex string (Canton MCMS format)
func (t DeployRateLimiter) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployRateLimiter from hex string (Canton MCMS format)
func (t *DeployRateLimiter) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployRateLimiterParams is a Record type
type DeployRateLimiterParams struct {
	InstanceId          types.TEXT                `json:"instanceId"`
	PoolInstanceId      types.TEXT                `json:"poolInstanceId"`
	PoolOwner           types.PARTY               `json:"poolOwner"`
	RemoteChainSelector types.NUMERIC             `json:"remoteChainSelector"`
	Direction           common.RateLimitDirection `json:"direction"`
	Mode                common.RateLimitMode      `json:"mode"`
	IsEnabled           types.BOOL                `json:"isEnabled"`
	Capacity            types.NUMERIC             `json:"capacity"`
	Rate                types.NUMERIC             `json:"rate"`
}

// ToMap converts DeployRateLimiterParams to a map for DAML arguments
func (t DeployRateLimiterParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["direction"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Direction).(mapper); ok {
			return m.toMap()
		}
		return t.Direction
	}()

	m["mode"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Mode).(mapper); ok {
			return m.toMap()
		}
		return t.Mode
	}()

	m["isEnabled"] = bool(t.IsEnabled)

	m["capacity"] = t.Capacity

	m["rate"] = t.Rate

	return m
}

func (t DeployRateLimiterParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployRateLimiterParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployRateLimiterParams to hex string (Canton MCMS format)
func (t DeployRateLimiterParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployRateLimiterParams from hex string (Canton MCMS format)
func (t *DeployRateLimiterParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployResult is a Record type
type DeployResult struct {
	FactoryCid  types.CONTRACT_ID `json:"factoryCid"`
	DeployedCid types.CONTRACT_ID `json:"deployedCid"`
}

// ToMap converts DeployResult to a map for DAML arguments
func (t DeployResult) ToMap() map[string]any {
	m := make(map[string]any)

	m["factoryCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.FactoryCid).(mapper); ok {
			return m.toMap()
		}
		return t.FactoryCid
	}()

	m["deployedCid"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.DeployedCid).(mapper); ok {
			return m.toMap()
		}
		return t.DeployedCid
	}()

	return m
}

func (t DeployResult) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployResult) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployResult to hex string (Canton MCMS format)
func (t DeployResult) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployResult from hex string (Canton MCMS format)
func (t *DeployResult) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployTokenAdminRegistry is a Record type
type DeployTokenAdminRegistry struct {
	Contract tokenadminregistry.TokenAdminRegistry `json:"contract"`
}

// ToMap converts DeployTokenAdminRegistry to a map for DAML arguments
func (t DeployTokenAdminRegistry) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = func() any {
		type mapper interface{ toMap() map[string]any }
		if m, ok := any(t.Contract).(mapper); ok {
			return m.toMap()
		}
		return t.Contract
	}()

	return m
}

func (t DeployTokenAdminRegistry) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployTokenAdminRegistry) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployTokenAdminRegistry to hex string (Canton MCMS format)
func (t DeployTokenAdminRegistry) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployTokenAdminRegistry from hex string (Canton MCMS format)
func (t *DeployTokenAdminRegistry) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployTokenAdminRegistryParams is a Record type
type DeployTokenAdminRegistryParams struct {
	InstanceId types.TEXT `json:"instanceId"`
}

// ToMap converts DeployTokenAdminRegistryParams to a map for DAML arguments
func (t DeployTokenAdminRegistryParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	return m
}

func (t DeployTokenAdminRegistryParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployTokenAdminRegistryParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployTokenAdminRegistryParams to hex string (Canton MCMS format)
func (t DeployTokenAdminRegistryParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployTokenAdminRegistryParams from hex string (Canton MCMS format)
func (t *DeployTokenAdminRegistryParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// FactoryState is a Record type
type FactoryState struct {
	InstanceId                    types.TEXT   `json:"instanceId"`
	Owner                         types.PARTY  `json:"owner"`
	McmsParty                     types.PARTY  `json:"mcmsParty"`
	UsedInstanceIds               types.GENMAP `json:"usedInstanceIds"`
	DeployedContracts             types.GENMAP `json:"deployedContracts"`
	PerPartyRouterFactoryDeployed types.BOOL   `json:"perPartyRouterFactoryDeployed"`
}

// ToMap converts FactoryState to a map for DAML arguments
func (t FactoryState) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	m["mcmsParty"] = t.McmsParty.ToMap()

	m["usedInstanceIds"] = func() any {
		if t.UsedInstanceIds == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.UsedInstanceIds}
	}()

	m["deployedContracts"] = func() any {
		if t.DeployedContracts == nil {
			return map[string]any{"_type": "genmap", "value": types.GENMAP{}}
		}
		return map[string]any{"_type": "genmap", "value": t.DeployedContracts}
	}()

	m["perPartyRouterFactoryDeployed"] = bool(t.PerPartyRouterFactoryDeployed)

	return m
}

func (t FactoryState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *FactoryState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes FactoryState to hex string (Canton MCMS format)
func (t FactoryState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes FactoryState from hex string (Canton MCMS format)
func (t *FactoryState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFactoryState is a Record type
type GetFactoryState struct {
	Caller types.PARTY `json:"caller"`
}

// ToMap converts GetFactoryState to a map for DAML arguments
func (t GetFactoryState) ToMap() map[string]any {
	m := make(map[string]any)

	m["caller"] = t.Caller.ToMap()

	return m
}

func (t GetFactoryState) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *GetFactoryState) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes GetFactoryState to hex string (Canton MCMS format)
func (t GetFactoryState) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFactoryState from hex string (Canton MCMS format)
func (t *GetFactoryState) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// GetFactoryStateMCMSParams is GetFactoryState without the Caller field for MCMS operationData encoding.
// Use this when encoding choice arguments for MCMS timelock operations.
type GetFactoryStateMCMSParams struct {
}

// MarshalHex encodes GetFactoryStateMCMSParams to hex string for MCMS operationData.
func (t GetFactoryStateMCMSParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes GetFactoryStateMCMSParams from hex string.
func (t *GetFactoryStateMCMSParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// SetOwnerToMCMS is a Record type
type SetOwnerToMCMS struct {
}

// ToMap converts SetOwnerToMCMS to a map for DAML arguments
func (t SetOwnerToMCMS) ToMap() map[string]any {
	m := make(map[string]any)
	return m
}

func (t SetOwnerToMCMS) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *SetOwnerToMCMS) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes SetOwnerToMCMS to hex string (Canton MCMS format)
func (t SetOwnerToMCMS) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes SetOwnerToMCMS from hex string (Canton MCMS format)
func (t *SetOwnerToMCMS) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// MCMSEncoder interface for typed encoding methods.
// Implemented by Encoder for method-based encoding.
type MCMSEncoder interface {
	DeployCCIPReceiver(args DeployCCIPReceiver) (*bind.EncodedChoice, error)
	DeployCCIPReceiverParams(args DeployCCIPReceiverParams) (*bind.EncodedChoice, error)
	DeployCCIPSender(args DeployCCIPSender) (*bind.EncodedChoice, error)
	DeployCCIPSenderParams(args DeployCCIPSenderParams) (*bind.EncodedChoice, error)
	DeployCommitteeVerifier(args DeployCommitteeVerifier) (*bind.EncodedChoice, error)
	DeployCommitteeVerifierParams(args DeployCommitteeVerifierParams) (*bind.EncodedChoice, error)
	DeployFeeQuoter(args DeployFeeQuoter) (*bind.EncodedChoice, error)
	DeployFeeQuoterParams(args DeployFeeQuoterParams) (*bind.EncodedChoice, error)
	DeployGlobalConfig(args DeployGlobalConfig) (*bind.EncodedChoice, error)
	DeployGlobalConfigParams(args DeployGlobalConfigParams) (*bind.EncodedChoice, error)
	DeployLockReleaseTokenPool(args DeployLockReleaseTokenPool) (*bind.EncodedChoice, error)
	DeployLockReleaseTokenPoolParams(args DeployLockReleaseTokenPoolParams) (*bind.EncodedChoice, error)
	DeployOffRamp(args DeployOffRamp) (*bind.EncodedChoice, error)
	DeployOffRampParams(args DeployOffRampParams) (*bind.EncodedChoice, error)
	DeployOnRamp(args DeployOnRamp) (*bind.EncodedChoice, error)
	DeployOnRampParams(args DeployOnRampParams) (*bind.EncodedChoice, error)
	DeployPerPartyRouterFactory(args DeployPerPartyRouterFactory) (*bind.EncodedChoice, error)
	DeployPerPartyRouterFactoryParams(args DeployPerPartyRouterFactoryParams) (*bind.EncodedChoice, error)
	DeployRMNRemote(args DeployRMNRemote) (*bind.EncodedChoice, error)
	DeployRMNRemoteParams(args DeployRMNRemoteParams) (*bind.EncodedChoice, error)
	DeployRateLimiter(args DeployRateLimiter) (*bind.EncodedChoice, error)
	DeployRateLimiterParams(args DeployRateLimiterParams) (*bind.EncodedChoice, error)
	DeployTokenAdminRegistry(args DeployTokenAdminRegistry) (*bind.EncodedChoice, error)
	DeployTokenAdminRegistryParams(args DeployTokenAdminRegistryParams) (*bind.EncodedChoice, error)
	GetFactoryState(args GetFactoryState) (*bind.EncodedChoice, error)
	GetFactoryStateMCMSParams(args GetFactoryStateMCMSParams) (*bind.EncodedChoice, error)
	SetOwnerToMCMS(args SetOwnerToMCMS) (*bind.EncodedChoice, error)
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

// DeployCCIPReceiver encodes parameters for the DeployCCIPReceiver choice.
func (e *encoder) DeployCCIPReceiver(args DeployCCIPReceiver) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCCIPReceiver", args)
}

// DeployCCIPReceiverParams encodes parameters for the DeployCCIPReceiver choice.
func (e *encoder) DeployCCIPReceiverParams(args DeployCCIPReceiverParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCCIPReceiver", args)
}

// DeployCCIPSender encodes parameters for the DeployCCIPSender choice.
func (e *encoder) DeployCCIPSender(args DeployCCIPSender) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCCIPSender", args)
}

// DeployCCIPSenderParams encodes parameters for the DeployCCIPSender choice.
func (e *encoder) DeployCCIPSenderParams(args DeployCCIPSenderParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCCIPSender", args)
}

// DeployCommitteeVerifier encodes parameters for the DeployCommitteeVerifier choice.
func (e *encoder) DeployCommitteeVerifier(args DeployCommitteeVerifier) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCommitteeVerifier", args)
}

// DeployCommitteeVerifierParams encodes parameters for the DeployCommitteeVerifier choice.
func (e *encoder) DeployCommitteeVerifierParams(args DeployCommitteeVerifierParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCommitteeVerifier", args)
}

// DeployFeeQuoter encodes parameters for the DeployFeeQuoter choice.
func (e *encoder) DeployFeeQuoter(args DeployFeeQuoter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployFeeQuoter", args)
}

// DeployFeeQuoterParams encodes parameters for the DeployFeeQuoter choice.
func (e *encoder) DeployFeeQuoterParams(args DeployFeeQuoterParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployFeeQuoter", args)
}

// DeployGlobalConfig encodes parameters for the DeployGlobalConfig choice.
func (e *encoder) DeployGlobalConfig(args DeployGlobalConfig) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployGlobalConfig", args)
}

// DeployGlobalConfigParams encodes parameters for the DeployGlobalConfig choice.
func (e *encoder) DeployGlobalConfigParams(args DeployGlobalConfigParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployGlobalConfig", args)
}

// DeployLockReleaseTokenPool encodes parameters for the DeployLockReleaseTokenPool choice.
func (e *encoder) DeployLockReleaseTokenPool(args DeployLockReleaseTokenPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployLockReleaseTokenPool", args)
}

// DeployLockReleaseTokenPoolParams encodes parameters for the DeployLockReleaseTokenPool choice.
func (e *encoder) DeployLockReleaseTokenPoolParams(args DeployLockReleaseTokenPoolParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployLockReleaseTokenPool", args)
}

// DeployOffRamp encodes parameters for the DeployOffRamp choice.
func (e *encoder) DeployOffRamp(args DeployOffRamp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployOffRamp", args)
}

// DeployOffRampParams encodes parameters for the DeployOffRamp choice.
func (e *encoder) DeployOffRampParams(args DeployOffRampParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployOffRamp", args)
}

// DeployOnRamp encodes parameters for the DeployOnRamp choice.
func (e *encoder) DeployOnRamp(args DeployOnRamp) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployOnRamp", args)
}

// DeployOnRampParams encodes parameters for the DeployOnRamp choice.
func (e *encoder) DeployOnRampParams(args DeployOnRampParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployOnRamp", args)
}

// DeployPerPartyRouterFactory encodes parameters for the DeployPerPartyRouterFactory choice.
func (e *encoder) DeployPerPartyRouterFactory(args DeployPerPartyRouterFactory) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployPerPartyRouterFactory", args)
}

// DeployPerPartyRouterFactoryParams encodes parameters for the DeployPerPartyRouterFactory choice.
func (e *encoder) DeployPerPartyRouterFactoryParams(args DeployPerPartyRouterFactoryParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployPerPartyRouterFactory", args)
}

// DeployRMNRemote encodes parameters for the DeployRMNRemote choice.
func (e *encoder) DeployRMNRemote(args DeployRMNRemote) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployRMNRemote", args)
}

// DeployRMNRemoteParams encodes parameters for the DeployRMNRemote choice.
func (e *encoder) DeployRMNRemoteParams(args DeployRMNRemoteParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployRMNRemote", args)
}

// DeployRateLimiter encodes parameters for the DeployRateLimiter choice.
func (e *encoder) DeployRateLimiter(args DeployRateLimiter) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployRateLimiter", args)
}

// DeployRateLimiterParams encodes parameters for the DeployRateLimiter choice.
func (e *encoder) DeployRateLimiterParams(args DeployRateLimiterParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployRateLimiter", args)
}

// DeployTokenAdminRegistry encodes parameters for the DeployTokenAdminRegistry choice.
func (e *encoder) DeployTokenAdminRegistry(args DeployTokenAdminRegistry) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployTokenAdminRegistry", args)
}

// DeployTokenAdminRegistryParams encodes parameters for the DeployTokenAdminRegistry choice.
func (e *encoder) DeployTokenAdminRegistryParams(args DeployTokenAdminRegistryParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployTokenAdminRegistry", args)
}

// GetFactoryState encodes parameters for the GetFactoryState choice.
func (e *encoder) GetFactoryState(args GetFactoryState) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFactoryState", args)
}

// GetFactoryStateMCMSParams encodes MCMS parameters (without Caller) for the GetFactoryState choice.
func (e *encoder) GetFactoryStateMCMSParams(args GetFactoryStateMCMSParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("GetFactoryState", args)
}

// SetOwnerToMCMS encodes parameters for the SetOwnerToMCMS choice.
func (e *encoder) SetOwnerToMCMS(args SetOwnerToMCMS) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("SetOwnerToMCMS", args)
}

// Verify MCMSEncoder interface implementation
var _ MCMSEncoder = (*encoder)(nil)
