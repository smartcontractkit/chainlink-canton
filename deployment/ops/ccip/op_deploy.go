package ccip

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

// CantonOpDeps is an alias for the shared type in the client package
type CantonOpDeps = compileClient.CantonOpDeps

// CantonOpResult wraps the output for Canton operations
type CantonOpResult[T any] struct {
	TransactionID string
	Output        T
}

// normalizeTemplateKey normalizes template ID to match the pattern used in tests
func normalizeTemplateKey(tid string) string {
	tid = strings.TrimPrefix(tid, "#")
	parts := strings.Split(tid, ":")
	if len(parts) < 3 {
		return tid
	}
	return parts[len(parts)-2] + ":" + parts[len(parts)-1]
}

// ============================================================================
// CCIP Common Operations
// ============================================================================

// DeployCCIPCommonInput contains input for deploying CCIP Common contracts
type DeployCCIPCommonInput struct {
	InstanceID         string
	ChainSelectorValue string // Numeric 0 value as string
	OnRampAddress      string // OnRamp address for GlobalConfig
}

// DeployCCIPCommonOutput contains the deployed contract IDs
type DeployCCIPCommonOutput struct {
	GlobalConfigContractID string
	GlobalConfigTemplateID string
}

var deployCCIPCommonHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployCCIPCommonInput) (output CantonOpResult[DeployCCIPCommonOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP Common package
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPCommon, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, commonDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, commonDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Parse chain selector
	chainSelector, ok := new(big.Int).SetString(input.ChainSelectorValue, 10)
	if !ok {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("invalid chainSelectorValue: %s", input.ChainSelectorValue)
	}

	// Create GlobalConfig contract
	// Note: We manually construct the command to ensure empty GENMAP fields are included
	// as DAML requires all non-optional fields to be present
	commandID := uuid.Must(uuid.NewUUID()).String()
	createCmd := &model.CreateCommand{
		TemplateID: fmt.Sprintf("#%s:%s:%s", common.PackageID, "CCIP.GlobalConfig", "GlobalConfig"),
		Arguments: map[string]interface{}{
			"ccipOwner":          types.PARTY(deps.Party).ToMap(),
			"instanceId":         input.InstanceID,
			"chainSelector":      chainSelector,
			"onRampAddress":      input.OnRampAddress,
			"destChainConfigs":   map[string]interface{}{"_type": "genmap", "value": map[string]interface{}{}},
			"sourceChainConfigs": map[string]interface{}{"_type": "genmap", "value": map[string]interface{}{}},
		},
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-common-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: createCmd}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("failed to submit GlobalConfig creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	globalConfigContractID := ""
	globalConfigTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.GlobalConfig:GlobalConfig" {
			globalConfigContractID = event.Created.ContractID
			globalConfigTemplateID = event.Created.TemplateID
			break
		}
	}

	if globalConfigContractID == "" {
		return CantonOpResult[DeployCCIPCommonOutput]{}, fmt.Errorf("failed to find GlobalConfig contract in transaction events")
	}

	fmt.Printf("Deployed GlobalConfig contract   id=%s\n", globalConfigContractID)

	return CantonOpResult[DeployCCIPCommonOutput]{
		TransactionID: commandID,
		Output: DeployCCIPCommonOutput{
			GlobalConfigContractID: globalConfigContractID,
			GlobalConfigTemplateID: globalConfigTemplateID,
		},
	}, nil
}

// ============================================================================
// Token Admin Registry Operations
// ============================================================================

// DeployTokenAdminRegistryInput contains input for deploying TokenAdminRegistry
type DeployTokenAdminRegistryInput struct {
	InstanceID string
}

// DeployTokenAdminRegistryOutput contains the deployed contract ID
type DeployTokenAdminRegistryOutput struct {
	TokenAdminRegistryContractID string
	TokenAdminRegistryTemplateID string
}

var deployTokenAdminRegistryHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployTokenAdminRegistryInput) (output CantonOpResult[DeployTokenAdminRegistryOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP TokenAdminRegistry package
	tokenAdminDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployTokenAdminRegistryOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPTokenAdminRegistry, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, tokenAdminDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployTokenAdminRegistryOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, tokenAdminDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployTokenAdminRegistryOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Create TokenAdminRegistry contract
	tokenAdminRegistry := tokenadminregistry.TokenAdminRegistry{
		Owner:        types.PARTY(deps.Party),
		InstanceId:   types.TEXT(input.InstanceID),
		TokenConfigs: types.GENMAP{}, // Empty initially
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-tokenadminregistry-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: tokenAdminRegistry.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployTokenAdminRegistryOutput]{}, fmt.Errorf("failed to submit TokenAdminRegistry creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	tokenAdminRegistryContractID := ""
	tokenAdminRegistryTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.TokenAdminRegistry:TokenAdminRegistry" {
			tokenAdminRegistryContractID = event.Created.ContractID
			tokenAdminRegistryTemplateID = event.Created.TemplateID
			break
		}
	}

	if tokenAdminRegistryContractID == "" {
		return CantonOpResult[DeployTokenAdminRegistryOutput]{}, fmt.Errorf("failed to find TokenAdminRegistry contract in transaction events")
	}

	fmt.Printf("Deployed TokenAdminRegistry contract   id=%s\n", tokenAdminRegistryContractID)

	return CantonOpResult[DeployTokenAdminRegistryOutput]{
		TransactionID: commandID,
		Output: DeployTokenAdminRegistryOutput{
			TokenAdminRegistryContractID: tokenAdminRegistryContractID,
			TokenAdminRegistryTemplateID: tokenAdminRegistryTemplateID,
		},
	}, nil
}

// ============================================================================
// Fee Quoter Operations
// ============================================================================

// DeployFeeQuoterInput contains input for deploying FeeQuoter
type DeployFeeQuoterInput struct {
	InstanceID string
}

// DeployFeeQuoterOutput contains the deployed contract ID
type DeployFeeQuoterOutput struct {
	FeeQuoterContractID string
	FeeQuoterTemplateID string
}

var deployFeeQuoterHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployFeeQuoterInput) (output CantonOpResult[DeployFeeQuoterOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP FeeQuoter package
	feeQuoterDar, err := contracts.GetDar(contracts.CCIPFeeQuoter, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployFeeQuoterOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPFeeQuoter, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, feeQuoterDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployFeeQuoterOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, feeQuoterDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployFeeQuoterOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Create FeeQuoter contract
	feeQuoter := feequoter.FeeQuoter{
		Owner:                            types.PARTY(deps.Party),
		InstanceId:                       types.TEXT(input.InstanceID),
		FeeTokens:                        types.GENMAP{},                         // Empty initially
		DestChainConfigs:                 types.GENMAP{},                         // Empty initially
		TokenTransferFeeConfigs:          types.GENMAP{},                         // Empty initially
		UsdPerUnitGasByDestChainSelector: types.GENMAP{},                         // Empty initially
		UsdPerToken:                      types.GENMAP{},                         // Empty initially
		PriceUpdaters:                    []types.PARTY{types.PARTY(deps.Party)}, // Owner is price updater
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-feequoter-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: feeQuoter.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployFeeQuoterOutput]{}, fmt.Errorf("failed to submit FeeQuoter creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	feeQuoterContractID := ""
	feeQuoterTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.FeeQuoter:FeeQuoter" {
			feeQuoterContractID = event.Created.ContractID
			feeQuoterTemplateID = event.Created.TemplateID
			break
		}
	}

	if feeQuoterContractID == "" {
		return CantonOpResult[DeployFeeQuoterOutput]{}, fmt.Errorf("failed to find FeeQuoter contract in transaction events")
	}

	fmt.Printf("Deployed FeeQuoter contract   id=%s\n", feeQuoterContractID)

	return CantonOpResult[DeployFeeQuoterOutput]{
		TransactionID: commandID,
		Output: DeployFeeQuoterOutput{
			FeeQuoterContractID: feeQuoterContractID,
			FeeQuoterTemplateID: feeQuoterTemplateID,
		},
	}, nil
}

// ============================================================================
// OffRamp Operations
// ============================================================================

// DeployOffRampInput contains input for deploying OffRamp
type DeployOffRampInput struct {
	InstanceID string
}

// DeployOffRampOutput contains the deployed contract ID
type DeployOffRampOutput struct {
	OffRampContractID string
	OffRampTemplateID string
}

var deployOffRampHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployOffRampInput) (output CantonOpResult[DeployOffRampOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP OffRamp package
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployOffRampOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPOffRamp, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, offRampDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployOffRampOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, offRampDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployOffRampOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Create OffRamp contract
	offRamp := offramp.OffRamp{
		CcipOwner:  types.PARTY(deps.Party),
		InstanceId: types.TEXT(input.InstanceID),
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-offramp-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: offRamp.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployOffRampOutput]{}, fmt.Errorf("failed to submit OffRamp creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	offRampContractID := ""
	offRampTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.OffRamp:OffRamp" {
			offRampContractID = event.Created.ContractID
			offRampTemplateID = event.Created.TemplateID
			break
		}
	}

	if offRampContractID == "" {
		return CantonOpResult[DeployOffRampOutput]{}, fmt.Errorf("failed to find OffRamp contract in transaction events")
	}

	fmt.Printf("Deployed OffRamp contract   id=%s\n", offRampContractID)

	return CantonOpResult[DeployOffRampOutput]{
		TransactionID: commandID,
		Output: DeployOffRampOutput{
			OffRampContractID: offRampContractID,
			OffRampTemplateID: offRampTemplateID,
		},
	}, nil
}

// ============================================================================
// PerPartyRouter Operations
// ============================================================================

// DeployPerPartyRouterInput contains input for deploying PerPartyRouter
type DeployPerPartyRouterInput struct {
	InstanceID string
}

// DeployPerPartyRouterOutput contains the deployed contract ID
type DeployPerPartyRouterOutput struct {
	PerPartyRouterFactoryContractID string
	PerPartyRouterFactoryTemplateID string
}

var deployPerPartyRouterHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployPerPartyRouterInput) (output CantonOpResult[DeployPerPartyRouterOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP PerPartyRouter package
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployPerPartyRouterOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPPerPartyRouter, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, perPartyRouterDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployPerPartyRouterOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, perPartyRouterDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployPerPartyRouterOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Create PerPartyRouterFactory contract (factory is deployed first, then users create their routers)
	perPartyRouterFactory := perpartyrouter.PerPartyRouterFactory{
		CcipOwner:         types.PARTY(deps.Party),
		InstanceId:        types.TEXT(input.InstanceID),
		RegisteredRouters: types.GENMAP{}, // Empty initially
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-perpartyrouter-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: perPartyRouterFactory.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployPerPartyRouterOutput]{}, fmt.Errorf("failed to submit PerPartyRouterFactory creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	perPartyRouterFactoryContractID := ""
	perPartyRouterFactoryTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.PerPartyRouter:PerPartyRouterFactory" {
			perPartyRouterFactoryContractID = event.Created.ContractID
			perPartyRouterFactoryTemplateID = event.Created.TemplateID
			break
		}
	}

	if perPartyRouterFactoryContractID == "" {
		return CantonOpResult[DeployPerPartyRouterOutput]{}, fmt.Errorf("failed to find PerPartyRouterFactory contract in transaction events")
	}

	fmt.Printf("Deployed PerPartyRouterFactory contract   id=%s\n", perPartyRouterFactoryContractID)

	return CantonOpResult[DeployPerPartyRouterOutput]{
		TransactionID: commandID,
		Output: DeployPerPartyRouterOutput{
			PerPartyRouterFactoryContractID: perPartyRouterFactoryContractID,
			PerPartyRouterFactoryTemplateID: perPartyRouterFactoryTemplateID,
		},
	}, nil
}

// ============================================================================
// OnRamp Operations
// ============================================================================

// DeployOnRampInput contains input for deploying OnRamp
type DeployOnRampInput struct {
	InstanceID           string
	DestChainSelector    string // Optional: destination chain selector
	DestChainOnRampBytes []byte // Optional: destination chain onramp address bytes
}

// DeployOnRampOutput contains the deployed contract ID
type DeployOnRampOutput struct {
	OnRampContractID string
	OnRampTemplateID string
}

var deployOnRampHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployOnRampInput) (output CantonOpResult[DeployOnRampOutput], err error) {
	ctx := b.GetContext()

	// Get and upload CCIP OnRamp package
	onRampDar, err := contracts.GetDar(contracts.CCIPOnRamp, contracts.CurrentVersion)
	if err != nil {
		return CantonOpResult[DeployOnRampOutput]{}, fmt.Errorf("failed to get DAR for package %s: %w", contracts.CCIPOnRamp, err)
	}

	submissionID := "validate-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.ValidateDarFile(ctx, onRampDar, submissionID)
	if err != nil {
		return CantonOpResult[DeployOnRampOutput]{}, fmt.Errorf("failed to validate DAR file: %w", err)
	}
	uploadSubmissionID := "upload-" + time.Now().Format("20060102150405")
	err = deps.BindingClient.PackageMng.UploadDarFile(ctx, onRampDar, uploadSubmissionID)
	if err != nil {
		return CantonOpResult[DeployOnRampOutput]{}, fmt.Errorf("failed to upload DAR file: %w", err)
	}

	// Create OnRamp contract
	onRamp := onramp.OnRamp{
		CcipOwner:  types.PARTY(deps.Party),
		InstanceId: types.TEXT(input.InstanceID),
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "ccip-onramp-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: onRamp.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployOnRampOutput]{}, fmt.Errorf("failed to submit OnRamp creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	onRampContractID := ""
	onRampTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "CCIP.OnRamp:OnRamp" {
			onRampContractID = event.Created.ContractID
			onRampTemplateID = event.Created.TemplateID
			break
		}
	}

	if onRampContractID == "" {
		return CantonOpResult[DeployOnRampOutput]{}, fmt.Errorf("failed to find OnRamp contract in transaction events")
	}

	fmt.Printf("Deployed OnRamp contract   id=%s\n", onRampContractID)

	return CantonOpResult[DeployOnRampOutput]{
		TransactionID: commandID,
		Output: DeployOnRampOutput{
			OnRampContractID: onRampContractID,
			OnRampTemplateID: onRampTemplateID,
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var DeployCCIPCommonOp = cld_ops.NewOperation(
	"canton/ccip/common/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP Common contracts (GlobalConfig) on Canton",
	deployCCIPCommonHandler,
)

var DeployTokenAdminRegistryOp = cld_ops.NewOperation(
	"canton/ccip/tokenadminregistry/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP TokenAdminRegistry contract on Canton",
	deployTokenAdminRegistryHandler,
)

var DeployFeeQuoterOp = cld_ops.NewOperation(
	"canton/ccip/feequoter/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP FeeQuoter contract on Canton",
	deployFeeQuoterHandler,
)

var DeployOffRampOp = cld_ops.NewOperation(
	"canton/ccip/offramp/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP OffRamp contract on Canton",
	deployOffRampHandler,
)

var DeployPerPartyRouterOp = cld_ops.NewOperation(
	"canton/ccip/perpartyrouter/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP PerPartyRouterFactory contract on Canton",
	deployPerPartyRouterHandler,
)

var DeployOnRampOp = cld_ops.NewOperation(
	"canton/ccip/onramp/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the CCIP OnRamp contract on Canton",
	deployOnRampHandler,
)
