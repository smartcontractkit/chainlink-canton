package mcms

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"

	"github.com/noders-team/go-daml/pkg/model"
	"github.com/noders-team/go-daml/pkg/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/mcms"
)

// ============================================================================
// MCMS Configuration Operations
// ============================================================================

// SetRootInput contains input for setting a merkle root
type SetRootInput struct {
	MCMSContractID string
	MCMSTemplateID string
	Submitter      string
	NewRoot        string
	ValidUntil     time.Time
	Metadata       mcms.RootMetadata
	MetadataProof  []string
	Signatures     []mcms.RawSignature
}

// SetRootOutput contains the transaction ID and new contract ID
type SetRootOutput struct {
	TransactionID     string
	NewMCMSContractID string
	NewMCMSTemplateID string
}

var setRootHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input SetRootInput) (output CantonOpResult[SetRootOutput], err error) {
	ctx := b.GetContext()

	// Build SetRoot args
	setRootArgs := mcms.SetRoot{
		Submitter:     types.PARTY(input.Submitter),
		NewRoot:       types.TEXT(input.NewRoot),
		ValidUntil:    types.TIMESTAMP(input.ValidUntil),
		Metadata:      input.Metadata,
		MetadataProof: make([]types.TEXT, len(input.MetadataProof)),
		Signatures:    input.Signatures,
	}

	for i, proof := range input.MetadataProof {
		setRootArgs.MetadataProof[i] = types.TEXT(proof)
	}

	// Build exercise command using generated bindings
	mcmsContract := mcms.MCMS{}
	exerciseCmd := mcmsContract.SetRoot(input.MCMSContractID, setRootArgs)

	// List known packages to find the package ID for mcms
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[SetRootOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var mcmsPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "mcms") {
			mcmsPkgID = p.PackageID
			break
		}
	}
	if mcmsPkgID == "" {
		return CantonOpResult[SetRootOutput]{}, fmt.Errorf("failed to find mcms package")
	}

	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "mcms-set-root",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", mcmsPkgID, "MCMS.Main", "MCMS"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[SetRootOutput]{}, fmt.Errorf("failed to set root: %w", err)
	}

	// Extract NEW MCMS CID from Created event
	newMCMSContractID := ""
	newMCMSTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == MCMSTemplateKey {
			newMCMSContractID = ev.Created.ContractID
			newMCMSTemplateID = ev.Created.TemplateID

			break
		}
	}

	if newMCMSContractID == "" {
		return CantonOpResult[SetRootOutput]{}, fmt.Errorf("set-root tx had no Created MCMS event; refusing to continue with old CID=%s", input.MCMSContractID)
	}

	return CantonOpResult[SetRootOutput]{
		TransactionID: commandID,
		Output: SetRootOutput{
			TransactionID:     commandID,
			NewMCMSContractID: newMCMSContractID,
			NewMCMSTemplateID: newMCMSTemplateID,
		},
	}, nil
}

// SetConfigInput contains input for setting configuration
type SetConfigInput struct {
	MCMSContractID string
	MCMSTemplateID string
	Config         mcms.SetConfig
}

// SetConfigOutput contains the transaction ID and new contract ID
type SetConfigOutput struct {
	TransactionID     string
	NewMCMSContractID string
	NewMCMSTemplateID string
}

var setConfigHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input SetConfigInput) (output CantonOpResult[SetConfigOutput], err error) {
	ctx := b.GetContext()

	// Build exercise command using generated bindings
	mcmsContract := mcms.MCMS{}
	exerciseCmd := mcmsContract.SetConfig(input.MCMSContractID, input.Config)

	// List known packages to find the package ID for mcms
	ListKnownPackagesResp, err := deps.BindingClient.PackageMng.ListKnownPackages(ctx)
	if err != nil {
		return CantonOpResult[SetConfigOutput]{}, fmt.Errorf("failed to list known packages: %w", err)
	}

	var mcmsPkgID string
	for _, p := range ListKnownPackagesResp {
		if strings.Contains(strings.ToLower(p.Name), "mcms") {
			mcmsPkgID = p.PackageID
			break
		}
	}
	if mcmsPkgID == "" {
		return CantonOpResult[SetConfigOutput]{}, fmt.Errorf("failed to find mcms package")
	}

	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "mcms-set-config",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands: []*model.Command{{
				Command: &model.ExerciseCommand{
					TemplateID: fmt.Sprintf("%s:%s:%s", mcmsPkgID, "MCMS.Main", "MCMS"),
					ContractID: exerciseCmd.ContractID,
					Choice:     exerciseCmd.Choice,
					Arguments:  exerciseCmd.Arguments,
				},
			}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[SetConfigOutput]{}, fmt.Errorf("failed to set config: %w", err)
	}

	// Extract NEW MCMS CID from Created event
	newMCMSContractID := ""
	newMCMSTemplateID := ""
	for _, ev := range submitResp.Transaction.Events {
		if ev.Created == nil {
			continue
		}
		normalized := normalizeTemplateKey(ev.Created.TemplateID)
		if normalized == MCMSTemplateKey {
			newMCMSContractID = ev.Created.ContractID
			newMCMSTemplateID = ev.Created.TemplateID

			break
		}
	}

	if newMCMSContractID == "" {
		return CantonOpResult[SetConfigOutput]{}, fmt.Errorf("set-config tx had no Created MCMS event; refusing to continue with old CID=%s", input.MCMSContractID)
	}

	return CantonOpResult[SetConfigOutput]{
		TransactionID: commandID,
		Output: SetConfigOutput{
			TransactionID:     commandID,
			NewMCMSContractID: newMCMSContractID,
			NewMCMSTemplateID: newMCMSTemplateID,
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var SetRootOp = cld_ops.NewOperation(
	"canton/mcms/set-root",
	semver.MustParse("0.1.0"),
	"Sets a merkle root for MCMS",
	setRootHandler,
)

var SetConfigOp = cld_ops.NewOperation(
	"canton/mcms/set-config",
	semver.MustParse("0.1.0"),
	"Sets configuration for MCMS",
	setConfigHandler,
)
