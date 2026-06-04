package factory

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	burnminttokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	ccipruntime "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	committeeverifier "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/committeeverifier"
	core "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	executor "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/executor"
	lockreleasetokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/lockreleasetokenpool"
	chainlinkapi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/chainlink/chainlinkapi"
	link "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	api "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
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
	PackageID   = "4171d8083aeb81fabbff6720b35f22f5e74f1f377141f16003d5cc7509a954c3"
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

// CCIPFactory is a Template type
type CCIPFactory struct {
	InstanceId                    types.TEXT                       `json:"instanceId"`
	Owner                         types.PARTY                      `json:"owner"`
	McmsParty                     types.PARTY                      `json:"mcmsParty"`
	UsedInstanceIds               map[types.TEXT]types.BOOL        `json:"usedInstanceIds"`
	DeployedContracts             map[types.TEXT]types.CONTRACT_ID `json:"deployedContracts"`
	PerPartyRouterFactoryDeployed types.BOOL                       `json:"perPartyRouterFactoryDeployed"`
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

// DeployLinkToken exercises the DeployLinkToken choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployLinkToken(contractID string, args DeployLinkToken) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployLinkToken",
		Arguments:  argsToMap(args),
	}
}

// DeployLinkTokenWithPackageID exercises the DeployLinkToken choice using the provided package ID instead of package name
func (t CCIPFactory) DeployLinkTokenWithPackageID(contractID string, packageID string, args DeployLinkToken) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployLinkToken",
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

// DeployExecutor exercises the DeployExecutor choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployExecutor(contractID string, args DeployExecutor) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployExecutor",
		Arguments:  argsToMap(args),
	}
}

// DeployExecutorWithPackageID exercises the DeployExecutor choice using the provided package ID instead of package name
func (t CCIPFactory) DeployExecutorWithPackageID(contractID string, packageID string, args DeployExecutor) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployExecutor",
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

// DeployBurnMintTokenPool exercises the DeployBurnMintTokenPool choice on this CCIPFactory contract
// This method uses the package name in the template ID
func (t CCIPFactory) DeployBurnMintTokenPool(contractID string, args DeployBurnMintTokenPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployBurnMintTokenPool",
		Arguments:  argsToMap(args),
	}
}

// DeployBurnMintTokenPoolWithPackageID exercises the DeployBurnMintTokenPool choice using the provided package ID instead of package name
func (t CCIPFactory) DeployBurnMintTokenPoolWithPackageID(contractID string, packageID string, args DeployBurnMintTokenPool) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "CCIPFactory"),
		ContractID: contractID,
		Choice:     "DeployBurnMintTokenPool",
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
func (t CCIPFactory) MCMSReceiverEntrypoint(contractID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", PackageName, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// MCMSReceiverEntrypointWithPackageID exercises the MCMSReceiver_Entrypoint choice using the provided package ID instead of package name
func (t CCIPFactory) MCMSReceiverEntrypointWithPackageID(contractID string, packageID string, args api.MCMSReceiverEntrypoint) *model.ExerciseCommand {
	return &model.ExerciseCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", packageID, "CCIP.Factory", "MCMSReceiver"),
		ContractID: contractID,
		Choice:     "MCMSReceiver_Entrypoint",
		Arguments:  argsToMap(args),
	}
}

// Verify interface implementations for CCIPFactory

var _ api.IMCMSReceiver = (*CCIPFactory)(nil)

// DeployBurnMintTokenPool is a Record type
type DeployBurnMintTokenPool struct {
	Contract burnminttokenpool.BurnMintTokenPool `json:"contract"`
}

// ToMap converts DeployBurnMintTokenPool to a map for DAML arguments
func (t DeployBurnMintTokenPool) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

	return m
}

func (t DeployBurnMintTokenPool) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployBurnMintTokenPool) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployBurnMintTokenPool to hex string (Canton MCMS format)
func (t DeployBurnMintTokenPool) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployBurnMintTokenPool from hex string (Canton MCMS format)
func (t *DeployBurnMintTokenPool) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployBurnMintTokenPoolParams is a Record type
type DeployBurnMintTokenPoolParams struct {
	InstanceId         types.TEXT                                 `json:"instanceId"`
	PoolOwner          types.PARTY                                `json:"poolOwner"`
	CcipOwner          types.PARTY                                `json:"ccipOwner"`
	InstrumentId       splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Decimals           types.INT64                                `json:"decimals"`
	RateLimitAdmin     *types.PARTY                               `json:"rateLimitAdmin" hex:"optional"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress            `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress            `json:"feeQuoter"`
	RmnRemote          chainlinkapi.RawInstanceAddress            `json:"rmnRemote"`
	PoolReceiveContext splice_api_token_metadata_v1.ChoiceContext `json:"poolReceiveContext"`
	TransferTimeout    burnminttokenpool.TransferTimeout          `json:"transferTimeout"`
}

// ToMap converts DeployBurnMintTokenPoolParams to a map for DAML arguments
func (t DeployBurnMintTokenPoolParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["decimals"] = int64(t.Decimals)

	if t.RateLimitAdmin != nil {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["poolReceiveContext"] = model.NestedToDAMLValue(t.PoolReceiveContext)

	m["transferTimeout"] = model.NestedToDAMLValue(t.TransferTimeout)

	return m
}

func (t DeployBurnMintTokenPoolParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployBurnMintTokenPoolParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployBurnMintTokenPoolParams to hex string (Canton MCMS format)
func (t DeployBurnMintTokenPoolParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployBurnMintTokenPoolParams from hex string (Canton MCMS format)
func (t *DeployBurnMintTokenPoolParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployCommitteeVerifier is a Record type
type DeployCommitteeVerifier struct {
	Contract committeeverifier.CommitteeVerifier `json:"contract"`
}

// ToMap converts DeployCommitteeVerifier to a map for DAML arguments
func (t DeployCommitteeVerifier) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId                   types.TEXT                      `json:"instanceId"`
	Owner                        types.PARTY                     `json:"owner"`
	CcipOwner                    types.PARTY                     `json:"ccipOwner"`
	VersionTag                   types.TEXT                      `json:"versionTag" hex:"bytes"`
	AllowListAdmin               *types.PARTY                    `json:"allowListAdmin" hex:"optional"`
	MessageSentObservers         []types.PARTY                   `json:"messageSentObservers"`
	RmnRemote                    chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	StorageLocations             []types.TEXT                    `json:"storageLocations"`
	StorageLocationsAdmin        types.PARTY                     `json:"storageLocationsAdmin"`
	PendingStorageLocationsAdmin types.PARTY                     `json:"pendingStorageLocationsAdmin"`
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
			"value": nil,
		}
	}

	m["messageSentObservers"] = func() []any {
		res := make([]any, 0, len(t.MessageSentObservers))
		for _, e := range t.MessageSentObservers {
			res = append(res, e.ToMap())
		}
		return res
	}()

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

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

// DeployExecutor is a Record type
type DeployExecutor struct {
	Contract executor.Executor `json:"contract"`
}

// ToMap converts DeployExecutor to a map for DAML arguments
func (t DeployExecutor) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

	return m
}

func (t DeployExecutor) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployExecutor) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployExecutor to hex string (Canton MCMS format)
func (t DeployExecutor) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployExecutor from hex string (Canton MCMS format)
func (t *DeployExecutor) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployExecutorParams is a Record type
type DeployExecutorParams struct {
	InstanceId            types.TEXT          `json:"instanceId"`
	Owner                 types.PARTY         `json:"owner"`
	MaxCCVsPerMsg         types.INT64         `json:"maxCCVsPerMsg"`
	AllowedFinalityConfig core.FinalityConfig `json:"allowedFinalityConfig"`
	CcvAllowlistEnabled   types.BOOL          `json:"ccvAllowlistEnabled"`
}

// ToMap converts DeployExecutorParams to a map for DAML arguments
func (t DeployExecutorParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["owner"] = t.Owner.ToMap()

	m["maxCCVsPerMsg"] = int64(t.MaxCCVsPerMsg)

	m["allowedFinalityConfig"] = model.NestedToDAMLValue(t.AllowedFinalityConfig)

	m["ccvAllowlistEnabled"] = bool(t.CcvAllowlistEnabled)

	return m
}

func (t DeployExecutorParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployExecutorParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployExecutorParams to hex string (Canton MCMS format)
func (t DeployExecutorParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployExecutorParams from hex string (Canton MCMS format)
func (t *DeployExecutorParams) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployFeeQuoter is a Record type
type DeployFeeQuoter struct {
	Contract core.FeeQuoter `json:"contract"`
}

// ToMap converts DeployFeeQuoter to a map for DAML arguments
func (t DeployFeeQuoter) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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

	m["linkTokenInstrumentId"] = model.NestedToDAMLValue(t.LinkTokenInstrumentId)

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
	Contract core.GlobalConfig `json:"contract"`
}

// ToMap converts DeployGlobalConfig to a map for DAML arguments
func (t DeployGlobalConfig) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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

// DeployLinkToken is a Record type
type DeployLinkToken struct {
	Contract link.LinkRegistry `json:"contract"`
}

// ToMap converts DeployLinkToken to a map for DAML arguments
func (t DeployLinkToken) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

	return m
}

func (t DeployLinkToken) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployLinkToken) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployLinkToken to hex string (Canton MCMS format)
func (t DeployLinkToken) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployLinkToken from hex string (Canton MCMS format)
func (t *DeployLinkToken) UnmarshalHex(data string) error {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Unmarshal(data, t)
}

// DeployLinkTokenParams is a Record type
type DeployLinkTokenParams struct {
	InstanceId   types.TEXT                               `json:"instanceId"`
	InstrumentId splice_api_token_holding_v1.InstrumentId `json:"instrumentId"`
}

// ToMap converts DeployLinkTokenParams to a map for DAML arguments
func (t DeployLinkTokenParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	return m
}

func (t DeployLinkTokenParams) MarshalJSON() ([]byte, error) {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Marshal(t)
}

func (t *DeployLinkTokenParams) UnmarshalJSON(data []byte) error {
	jsonCodec := codec.NewJsonCodec()
	return jsonCodec.Unmarshal(data, t)
}

// MarshalHex encodes DeployLinkTokenParams to hex string (Canton MCMS format)
func (t DeployLinkTokenParams) MarshalHex() (string, error) {
	hexCodec := codec.NewHexCodec()
	return hexCodec.Marshal(t)
}

// UnmarshalHex decodes DeployLinkTokenParams from hex string (Canton MCMS format)
func (t *DeployLinkTokenParams) UnmarshalHex(data string) error {
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

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId         types.TEXT                                 `json:"instanceId"`
	PoolOwner          types.PARTY                                `json:"poolOwner"`
	CcipOwner          types.PARTY                                `json:"ccipOwner"`
	InstrumentId       splice_api_token_holding_v1.InstrumentId   `json:"instrumentId"`
	Decimals           types.INT64                                `json:"decimals"`
	RateLimitAdmin     *types.PARTY                               `json:"rateLimitAdmin" hex:"optional"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress            `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress            `json:"feeQuoter"`
	RmnRemote          chainlinkapi.RawInstanceAddress            `json:"rmnRemote"`
	PoolReceiveContext splice_api_token_metadata_v1.ChoiceContext `json:"poolReceiveContext"`
	TransferTimeout    lockreleasetokenpool.TransferTimeout       `json:"transferTimeout"`
}

// ToMap converts DeployLockReleaseTokenPoolParams to a map for DAML arguments
func (t DeployLockReleaseTokenPoolParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["ccipOwner"] = t.CcipOwner.ToMap()

	m["instrumentId"] = model.NestedToDAMLValue(t.InstrumentId)

	m["decimals"] = int64(t.Decimals)

	if t.RateLimitAdmin != nil {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": (*t.RateLimitAdmin).ToMap(),
		}
	} else {
		m["rateLimitAdmin"] = map[string]any{
			"_type": "optional",
			"value": nil,
		}
	}

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["poolReceiveContext"] = model.NestedToDAMLValue(t.PoolReceiveContext)

	m["transferTimeout"] = model.NestedToDAMLValue(t.TransferTimeout)

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
	Contract ccipruntime.OffRamp `json:"contract"`
}

// ToMap converts DeployOffRamp to a map for DAML arguments
func (t DeployOffRamp) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId         types.TEXT                      `json:"instanceId"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
}

// ToMap converts DeployOffRampParams to a map for DAML arguments
func (t DeployOffRampParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

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
	Contract ccipruntime.OnRamp `json:"contract"`
}

// ToMap converts DeployOnRamp to a map for DAML arguments
func (t DeployOnRamp) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId         types.TEXT                      `json:"instanceId"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
	CcvRegistry        chainlinkapi.RawInstanceAddress `json:"ccvRegistry"`
	MaxUSDCentsPerMsg  types.NUMERIC                   `json:"maxUSDCentsPerMsg"`
}

// ToMap converts DeployOnRampParams to a map for DAML arguments
func (t DeployOnRampParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["ccvRegistry"] = model.NestedToDAMLValue(t.CcvRegistry)

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
	Contract ccipruntime.PerPartyRouterFactory `json:"contract"`
}

// ToMap converts DeployPerPartyRouterFactory to a map for DAML arguments
func (t DeployPerPartyRouterFactory) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId         types.TEXT                      `json:"instanceId"`
	OnRamp             chainlinkapi.RawInstanceAddress `json:"onRamp"`
	OffRamp            chainlinkapi.RawInstanceAddress `json:"offRamp"`
	GlobalConfig       chainlinkapi.RawInstanceAddress `json:"globalConfig"`
	TokenAdminRegistry chainlinkapi.RawInstanceAddress `json:"tokenAdminRegistry"`
	FeeQuoter          chainlinkapi.RawInstanceAddress `json:"feeQuoter"`
	RmnRemote          chainlinkapi.RawInstanceAddress `json:"rmnRemote"`
}

// ToMap converts DeployPerPartyRouterFactoryParams to a map for DAML arguments
func (t DeployPerPartyRouterFactoryParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["onRamp"] = model.NestedToDAMLValue(t.OnRamp)

	m["offRamp"] = model.NestedToDAMLValue(t.OffRamp)

	m["globalConfig"] = model.NestedToDAMLValue(t.GlobalConfig)

	m["tokenAdminRegistry"] = model.NestedToDAMLValue(t.TokenAdminRegistry)

	m["feeQuoter"] = model.NestedToDAMLValue(t.FeeQuoter)

	m["rmnRemote"] = model.NestedToDAMLValue(t.RmnRemote)

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
	Contract core.RMNRemote `json:"contract"`
}

// ToMap converts DeployRMNRemote to a map for DAML arguments
func (t DeployRMNRemote) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	Contract core.RateLimiter `json:"contract"`
}

// ToMap converts DeployRateLimiter to a map for DAML arguments
func (t DeployRateLimiter) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId          types.TEXT              `json:"instanceId"`
	PoolInstanceId      types.TEXT              `json:"poolInstanceId"`
	PoolOwner           types.PARTY             `json:"poolOwner"`
	RemoteChainSelector types.NUMERIC           `json:"remoteChainSelector"`
	Direction           core.RateLimitDirection `json:"direction"`
	Mode                core.RateLimitMode      `json:"mode"`
	IsEnabled           types.BOOL              `json:"isEnabled"`
	Capacity            types.NUMERIC           `json:"capacity"`
	Rate                types.NUMERIC           `json:"rate"`
}

// ToMap converts DeployRateLimiterParams to a map for DAML arguments
func (t DeployRateLimiterParams) ToMap() map[string]any {
	m := make(map[string]any)

	m["instanceId"] = string(t.InstanceId)

	m["poolInstanceId"] = string(t.PoolInstanceId)

	m["poolOwner"] = t.PoolOwner.ToMap()

	m["remoteChainSelector"] = t.RemoteChainSelector

	m["direction"] = model.NestedToDAMLValue(t.Direction)

	m["mode"] = model.NestedToDAMLValue(t.Mode)

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

	m["factoryCid"] = model.NestedToDAMLValue(t.FactoryCid)

	m["deployedCid"] = model.NestedToDAMLValue(t.DeployedCid)

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
	Contract core.TokenAdminRegistry `json:"contract"`
}

// ToMap converts DeployTokenAdminRegistry to a map for DAML arguments
func (t DeployTokenAdminRegistry) ToMap() map[string]any {
	m := make(map[string]any)

	m["contract"] = model.NestedToDAMLValue(t.Contract)

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
	InstanceId                    types.TEXT                       `json:"instanceId"`
	Owner                         types.PARTY                      `json:"owner"`
	McmsParty                     types.PARTY                      `json:"mcmsParty"`
	UsedInstanceIds               map[types.TEXT]types.BOOL        `json:"usedInstanceIds"`
	DeployedContracts             map[types.TEXT]types.CONTRACT_ID `json:"deployedContracts"`
	PerPartyRouterFactoryDeployed types.BOOL                       `json:"perPartyRouterFactoryDeployed"`
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
	DeployBurnMintTokenPool(args DeployBurnMintTokenPool) (*bind.EncodedChoice, error)
	DeployBurnMintTokenPoolParams(args DeployBurnMintTokenPoolParams) (*bind.EncodedChoice, error)
	DeployCommitteeVerifier(args DeployCommitteeVerifier) (*bind.EncodedChoice, error)
	DeployCommitteeVerifierParams(args DeployCommitteeVerifierParams) (*bind.EncodedChoice, error)
	DeployExecutor(args DeployExecutor) (*bind.EncodedChoice, error)
	DeployExecutorParams(args DeployExecutorParams) (*bind.EncodedChoice, error)
	DeployFeeQuoter(args DeployFeeQuoter) (*bind.EncodedChoice, error)
	DeployFeeQuoterParams(args DeployFeeQuoterParams) (*bind.EncodedChoice, error)
	DeployGlobalConfig(args DeployGlobalConfig) (*bind.EncodedChoice, error)
	DeployGlobalConfigParams(args DeployGlobalConfigParams) (*bind.EncodedChoice, error)
	DeployLinkToken(args DeployLinkToken) (*bind.EncodedChoice, error)
	DeployLinkTokenParams(args DeployLinkTokenParams) (*bind.EncodedChoice, error)
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

// DeployBurnMintTokenPool encodes parameters for the DeployBurnMintTokenPool choice.
func (e *encoder) DeployBurnMintTokenPool(args DeployBurnMintTokenPool) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployBurnMintTokenPool", args)
}

// DeployBurnMintTokenPoolParams encodes parameters for the DeployBurnMintTokenPool choice.
func (e *encoder) DeployBurnMintTokenPoolParams(args DeployBurnMintTokenPoolParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployBurnMintTokenPool", args)
}

// DeployCommitteeVerifier encodes parameters for the DeployCommitteeVerifier choice.
func (e *encoder) DeployCommitteeVerifier(args DeployCommitteeVerifier) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCommitteeVerifier", args)
}

// DeployCommitteeVerifierParams encodes parameters for the DeployCommitteeVerifier choice.
func (e *encoder) DeployCommitteeVerifierParams(args DeployCommitteeVerifierParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployCommitteeVerifier", args)
}

// DeployExecutor encodes parameters for the DeployExecutor choice.
func (e *encoder) DeployExecutor(args DeployExecutor) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployExecutor", args)
}

// DeployExecutorParams encodes parameters for the DeployExecutor choice.
func (e *encoder) DeployExecutorParams(args DeployExecutorParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployExecutor", args)
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

// DeployLinkToken encodes parameters for the DeployLinkToken choice.
func (e *encoder) DeployLinkToken(args DeployLinkToken) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployLinkToken", args)
}

// DeployLinkTokenParams encodes parameters for the DeployLinkToken choice.
func (e *encoder) DeployLinkTokenParams(args DeployLinkTokenParams) (*bind.EncodedChoice, error) {
	return e.EncodeChoiceArgs("DeployLinkToken", args)
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
