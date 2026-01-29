package ccip

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/feequoter"
)

// ============================================================================
// GlobalConfig Configuration Operations
// ============================================================================

// UpdateGlobalConfigDestChainConfigInput contains input for updating dest chain config
type UpdateGlobalConfigDestChainConfigInput struct {
	GlobalConfigContractID string
	GlobalConfigTemplateID string // <-- add this (from deploy or previous update)
	DestChainSelector      string
	Config                 common.DestChainConfig
}

type UpdateGlobalConfigDestChainConfigOutput struct {
	NewGlobalConfigContractID string
	NewGlobalConfigTemplateID string
}

var updateGlobalConfigDestChainConfigHandler = func(
	b cld_ops.Bundle,
	deps CantonOpDeps,
	input UpdateGlobalConfigDestChainConfigInput,
) (output CantonOpResult[UpdateGlobalConfigDestChainConfigOutput], err error) {
	ctx := b.GetContext()

	// Parse chain selector
	chainSelector, ok := new(big.Int).SetString(input.DestChainSelector, 10)
	if !ok {
		return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{},
			fmt.Errorf("invalid destChainSelector: %s", input.DestChainSelector)
	}

	// Scale to NUMERIC(10) mantissa
	scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	mantissa := new(big.Int).Mul(chainSelector, scale10)

	// Create update args
	updateArgs := common.UpdateDestChainConfig{
		DestChainSelector: types.NUMERIC(mantissa),
		Config:            input.Config,
	}

	// Build exercise command using generated bindings
	globalConfig := common.GlobalConfig{}
	exerciseCmd := globalConfig.UpdateDestChainConfig(input.GlobalConfigContractID, updateArgs)

	// List known packages to find the package ID for ccip-common
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipCommonPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-common") {
			ccipCommonPkgID = p.PackageID
			break
		}
	}
	if ccipCommonPkgID == "" {
		return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{}, fmt.Errorf("failed to find ccip-common package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", ccipCommonPkgID, "CCIP.GlobalConfig", "GlobalConfig"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{},
			fmt.Errorf("failed to update GlobalConfig dest chain config: %w", err)
	}

	// Extract NEW GlobalConfig CID from Created event
	newGlobalConfigContractID := ""
	newGlobalConfigTemplateID := ""

	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		// NOTE: keep this consistent with your normalizeTemplateKey behavior
		// If it produces "CCIP.GlobalConfig:GlobalConfig" this is correct.
		if normalized == GlobalConfigTemplateKey {
			newGlobalConfigContractID = ev.Created.ContractID
			newGlobalConfigTemplateID = ev.Created.TemplateID

			break
		}
	}

	if newGlobalConfigContractID == "" {
		return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{},
			fmt.Errorf("dest-chain update tx had no Created GlobalConfig event; refusing to continue with old CID=%s", input.GlobalConfigContractID)
	}

	return CantonOpResult[UpdateGlobalConfigDestChainConfigOutput]{
		UpdateID: submitResp.UpdateID,
		Output: UpdateGlobalConfigDestChainConfigOutput{
			NewGlobalConfigContractID: newGlobalConfigContractID,
			NewGlobalConfigTemplateID: newGlobalConfigTemplateID,
		},
	}, nil
}

// UpdateGlobalConfigSourceChainConfigInput contains input for updating source chain config
type UpdateGlobalConfigSourceChainConfigInput struct {
	GlobalConfigContractID string
	GlobalConfigTemplateID string // <-- add this (from deploy or previous update)
	SourceChainSelector    string
	Config                 common.SourceChainConfig
}

// UpdateGlobalConfigSourceChainConfigOutput contains the transaction ID
type UpdateGlobalConfigSourceChainConfigOutput struct {
	NewGlobalConfigContractID string
	NewGlobalConfigTemplateID string
}

var updateGlobalConfigSourceChainConfigHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input UpdateGlobalConfigSourceChainConfigInput) (output CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput], err error) {
	ctx := b.GetContext()

	// Parse chain selector
	chainSelector, ok := new(big.Int).SetString(input.SourceChainSelector, 10)
	if !ok {
		return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{}, fmt.Errorf("invalid sourceChainSelector: %s", input.SourceChainSelector)
	}

	// Scale to NUMERIC(10) mantissa
	scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
	mantissa := new(big.Int).Mul(chainSelector, scale10)

	// Create update command
	updateArgs := common.UpdateSourceChainConfig{
		SourceChainSelector: types.NUMERIC(mantissa),
		Config:              input.Config,
	}

	globalConfig := common.GlobalConfig{}
	exerciseCmd := globalConfig.UpdateSourceChainConfig(input.GlobalConfigContractID, updateArgs)

	// List known packages to find the package ID for ccip-common
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipCommonPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-common") {
			ccipCommonPkgID = p.PackageID
			break
		}
	}
	if ccipCommonPkgID == "" {
		return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{}, fmt.Errorf("failed to find ccip-common package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				// TODO find a better way rather than this templateID override hack which exposes PackageID to the client
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipCommonPkgID, "CCIP.GlobalConfig", "GlobalConfig"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{}, fmt.Errorf("failed to update GlobalConfig source chain config: %w", err)
	}

	newGlobalConfigTemplateID := ""
	newGlobalConfigContractID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == GlobalConfigTemplateKey {
			newGlobalConfigContractID = ev.Created.ContractID
			newGlobalConfigTemplateID = ev.Created.TemplateID

			break
		}
	}
	if newGlobalConfigTemplateID == "" {
		return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{}, fmt.Errorf("source-chain update tx had no Created GlobalConfig event; refusing to continue with old CID=%s", input.GlobalConfigContractID)
	}

	return CantonOpResult[UpdateGlobalConfigSourceChainConfigOutput]{
		UpdateID: submitResp.UpdateID,
		Output: UpdateGlobalConfigSourceChainConfigOutput{
			NewGlobalConfigContractID: newGlobalConfigContractID,
			NewGlobalConfigTemplateID: newGlobalConfigTemplateID,
		},
	}, nil
}

// ============================================================================
// FeeQuoter Configuration Operations
// ============================================================================

// UpdateFeeQuoterPricesInput contains input for updating prices
type UpdateFeeQuoterPricesInput struct {
	FeeQuoterContractID string
	FeeQuoterTemplateID string
	PriceUpdates        feequoter.PriceUpdates
}

// UpdateFeeQuoterPricesOutput contains the transaction ID
type UpdateFeeQuoterPricesOutput struct {
	NewFeeQuoterContractID string
	NewFeeQuoterTemplateID string
}

var updateFeeQuoterPricesHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input UpdateFeeQuoterPricesInput) (output CantonOpResult[UpdateFeeQuoterPricesOutput], err error) {
	ctx := b.GetContext()

	// Create update command
	updateArgs := feequoter.UpdatePrices{
		PriceUpdates: input.PriceUpdates,
		Caller:       types.PARTY(deps.Party),
	}

	feeQuoter := feequoter.FeeQuoter{}
	exerciseCmd := feeQuoter.UpdatePrices(input.FeeQuoterContractID, updateArgs)

	// List known packages to find the package ID for ccip-feequoter
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[UpdateFeeQuoterPricesOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipFeeQuoterPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-feequoter") {
			ccipFeeQuoterPkgID = p.PackageID
			break
		}
	}
	if ccipFeeQuoterPkgID == "" {
		return CantonOpResult[UpdateFeeQuoterPricesOutput]{}, fmt.Errorf("failed to find ccip-feequoter package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				// TODO find a better way rather than this templateID override hack which exposes PackageID to the client
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipFeeQuoterPkgID, "CCIP.FeeQuoter", "FeeQuoter"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[UpdateFeeQuoterPricesOutput]{}, fmt.Errorf("failed to update FeeQuoter prices: %w", err)
	}

	newFeeQuoterContractID := ""
	newFeeQuoterTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == FeeQuoterTemplateKey {
			newFeeQuoterContractID = ev.Created.ContractID
			newFeeQuoterTemplateID = ev.Created.TemplateID

			break
		}
	}
	if newFeeQuoterTemplateID == "" {
		return CantonOpResult[UpdateFeeQuoterPricesOutput]{}, fmt.Errorf("fee-quoter-prices update tx had no Created FeeQuoter event; refusing to continue with old CID=%s", input.FeeQuoterContractID)
	}

	return CantonOpResult[UpdateFeeQuoterPricesOutput]{
		UpdateID: submitResp.UpdateID,
		Output: UpdateFeeQuoterPricesOutput{
			NewFeeQuoterContractID: newFeeQuoterContractID,
			NewFeeQuoterTemplateID: newFeeQuoterTemplateID,
		},
	}, nil
}

// ApplyFeeQuoterFeeTokenUpdatesInput contains input for applying fee token updates
type ApplyFeeQuoterFeeTokenUpdatesInput struct {
	FeeQuoterContractID string
	FeeQuoterTemplateID string
	FeeTokensToRemove   []feequoter.InstrumentId
	FeeTokensToAdd      []feequoter.FeeTokenArgs
}

// ApplyFeeQuoterFeeTokenUpdatesOutput contains the transaction ID
type ApplyFeeQuoterFeeTokenUpdatesOutput struct {
	NewFeeQuoterContractID string
	NewFeeQuoterTemplateID string
}

var applyFeeQuoterFeeTokenUpdatesHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input ApplyFeeQuoterFeeTokenUpdatesInput) (output CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput], err error) {
	ctx := b.GetContext()

	// Create update command
	updateArgs := feequoter.ApplyFeeTokenUpdates{
		FeeTokensToRemove: input.FeeTokensToRemove,
		FeeTokensToAdd:    input.FeeTokensToAdd,
		Caller:            types.PARTY(deps.Party),
	}

	feeQuoter := feequoter.FeeQuoter{}
	exerciseCmd := feeQuoter.ApplyFeeTokenUpdates(input.FeeQuoterContractID, updateArgs)

	// List known packages to find the package ID for ccip-feequoter
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipFeeQuoterPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-feequoter") {
			ccipFeeQuoterPkgID = p.PackageID
			break
		}
	}
	if ccipFeeQuoterPkgID == "" {
		return CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput]{}, fmt.Errorf("failed to find ccip-feequoter package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				// TODO find a better way rather than this templateID override hack which exposes PackageID to the client
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipFeeQuoterPkgID, "CCIP.FeeQuoter", "FeeQuoter"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput]{}, fmt.Errorf("failed to apply FeeQuoter fee token updates: %w", err)
	}

	newFeeQuoterContractID := ""
	newFeeQuoterTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == FeeQuoterTemplateKey {
			newFeeQuoterContractID = ev.Created.ContractID
			newFeeQuoterTemplateID = ev.Created.TemplateID

			break
		}
	}
	if newFeeQuoterTemplateID == "" {
		return CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput]{}, fmt.Errorf("fee-quoter-fee-token-updates update tx had no Created FeeQuoter event; refusing to continue with old CID=%s", input.FeeQuoterContractID)
	}

	return CantonOpResult[ApplyFeeQuoterFeeTokenUpdatesOutput]{
		UpdateID: submitResp.UpdateID,
		Output: ApplyFeeQuoterFeeTokenUpdatesOutput{
			NewFeeQuoterContractID: newFeeQuoterContractID,
			NewFeeQuoterTemplateID: newFeeQuoterTemplateID,
		},
	}, nil
}

// ApplyFeeQuoterDestChainConfigUpdatesInput contains input for applying dest chain config updates
type ApplyFeeQuoterDestChainConfigUpdatesInput struct {
	FeeQuoterContractID string
	FeeQuoterTemplateID string
	DestChainConfigArgs []feequoter.DestChainConfigArgs
}

// ApplyFeeQuoterDestChainConfigUpdatesOutput contains the transaction ID
type ApplyFeeQuoterDestChainConfigUpdatesOutput struct {
	NewFeeQuoterContractID string
	NewFeeQuoterTemplateID string
}

var applyFeeQuoterDestChainConfigUpdatesHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input ApplyFeeQuoterDestChainConfigUpdatesInput) (output CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput], err error) {
	ctx := b.GetContext()

	// Create update command
	updateArgs := feequoter.ApplyDestChainConfigUpdates{
		DestChainConfigArgs: input.DestChainConfigArgs,
	}

	feeQuoter := feequoter.FeeQuoter{}
	exerciseCmd := feeQuoter.ApplyDestChainConfigUpdates(input.FeeQuoterContractID, updateArgs)

	// List known packages to find the package ID for ccip-feequoter
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var ccipFeeQuoterPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "ccip-feequoter") {
			ccipFeeQuoterPkgID = p.PackageID
			break
		}
	}
	if ccipFeeQuoterPkgID == "" {
		return CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput]{}, fmt.Errorf("failed to find ccip-feequoter package")
	}

	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			CommandID: uuid.Must(uuid.NewUUID()).String(),
			ActAs:     []string{deps.Party},
			Commands: []*model.Command{{Command: &model.ExerciseCommand{
				// TODO find a better way rather than this templateID override hack which exposes PackageID to the client
				TemplateID: fmt.Sprintf("%s:%s:%s", ccipFeeQuoterPkgID, "CCIP.FeeQuoter", "FeeQuoter"),
				ContractID: exerciseCmd.ContractID,
				Choice:     exerciseCmd.Choice,
				Arguments:  exerciseCmd.Arguments,
			}}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput]{}, fmt.Errorf("failed to apply FeeQuoter dest chain config updates: %w", err)
	}

	newFeeQuoterContractID := ""
	newFeeQuoterTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == FeeQuoterTemplateKey {
			newFeeQuoterContractID = ev.Created.ContractID
			newFeeQuoterTemplateID = ev.Created.TemplateID

			break
		}
	}
	if newFeeQuoterTemplateID == "" {
		return CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput]{}, fmt.Errorf("fee-quoter-dest-chain-config-updates update tx had no Created FeeQuoter event; refusing to continue with old CID=%s", input.FeeQuoterContractID)
	}

	return CantonOpResult[ApplyFeeQuoterDestChainConfigUpdatesOutput]{
		UpdateID: submitResp.UpdateID,
		Output: ApplyFeeQuoterDestChainConfigUpdatesOutput{
			NewFeeQuoterContractID: newFeeQuoterContractID,
			NewFeeQuoterTemplateID: newFeeQuoterTemplateID,
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var UpdateGlobalConfigDestChainConfigOp = cld_ops.NewOperation(
	"canton/ccip/globalconfig/update-dest-chain-config",
	semver.MustParse("0.1.0"),
	"Updates the destination chain configuration in GlobalConfig",
	updateGlobalConfigDestChainConfigHandler,
)

var UpdateGlobalConfigSourceChainConfigOp = cld_ops.NewOperation(
	"canton/ccip/globalconfig/update-source-chain-config",
	semver.MustParse("0.1.0"),
	"Updates the source chain configuration in GlobalConfig",
	updateGlobalConfigSourceChainConfigHandler,
)

var UpdateFeeQuoterPricesOp = cld_ops.NewOperation(
	"canton/ccip/feequoter/update-prices",
	semver.MustParse("0.1.0"),
	"Updates prices in FeeQuoter",
	updateFeeQuoterPricesHandler,
)

var ApplyFeeQuoterFeeTokenUpdatesOp = cld_ops.NewOperation(
	"canton/ccip/feequoter/apply-fee-token-updates",
	semver.MustParse("0.1.0"),
	"Applies fee token updates to FeeQuoter",
	applyFeeQuoterFeeTokenUpdatesHandler,
)

var ApplyFeeQuoterDestChainConfigUpdatesOp = cld_ops.NewOperation(
	"canton/ccip/feequoter/apply-dest-chain-config-updates",
	semver.MustParse("0.1.0"),
	"Applies destination chain configuration updates to FeeQuoter",
	applyFeeQuoterDestChainConfigUpdatesHandler,
)
