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

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/mcms"
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
// MCMS Operations
// ============================================================================

// DeployMCMSInput contains input for deploying MCMS
type DeployMCMSInput struct {
	InstanceID string
	ChainID    int64
	MCMSID     string
	Role       string // "Proposer" or "Verifier"
	Config     mcms.MultisigConfig
}

// DeployMCMSOutput contains the deployed contract ID
type DeployMCMSOutput struct {
	MCMSContractID string
	MCMSTemplateID string
}

var deployMCMSHandler = func(b cld_ops.Bundle, deps CantonOpDeps, input DeployMCMSInput) (output CantonOpResult[DeployMCMSOutput], err error) {
	ctx := b.GetContext()

	// Parse role
	var role mcms.Role
	switch strings.ToLower(input.Role) {
	case "proposer":
		role = mcms.RoleProposer
	default:
		return CantonOpResult[DeployMCMSOutput]{}, fmt.Errorf("invalid role: %s (must be 'Proposer' or 'Verifier')", input.Role)
	}

	// Create empty expiring root
	emptyExpiringRoot := mcms.ExpiringRoot{
		Root:       types.TEXT(""),
		ValidUntil: types.TIMESTAMP(time.Unix(0, 0).UTC()),
		OpCount:    types.INT64(0),
	}

	// Create empty root metadata
	emptyRootMetadata := mcms.RootMetadata{
		ChainId:              types.INT64(0),
		MultisigId:           types.TEXT(""),
		PreOpCount:           types.INT64(0),
		PostOpCount:          types.INT64(0),
		OverridePreviousRoot: types.BOOL(false),
	}

	// Create MCMS contract
	mcmsContract := mcms.MCMS{
		Owner:        types.PARTY(deps.Party),
		Role:         role,
		ChainId:      types.INT64(input.ChainID),
		McmsId:       types.TEXT(input.MCMSID),
		Config:       input.Config,
		SeenHashes:   types.GENMAP{}, // Empty map
		ExpiringRoot: emptyExpiringRoot,
		RootMetadata: emptyRootMetadata,
	}

	// Submit via binding client's CommandService
	commandID := uuid.Must(uuid.NewUUID()).String()
	cmds := &model.SubmitAndWaitRequest{
		Commands: &model.Commands{
			WorkflowID: "mcms-deploy",
			UserID:     deps.UserID,
			CommandID:  commandID,
			ActAs:      []string{deps.Party},
			Commands:   []*model.Command{{Command: mcmsContract.CreateCommand()}},
		},
	}

	submitResp, err := deps.BindingClient.CommandService.SubmitAndWaitForTransaction(ctx, cmds)
	if err != nil {
		return CantonOpResult[DeployMCMSOutput]{}, fmt.Errorf("failed to submit MCMS creation: %w", err)
	}

	// Retrieve the contract ID and template ID from the create event
	mcmsContractID := ""
	mcmsTemplateID := ""
	for _, event := range submitResp.Transaction.Events {
		if event.Created == nil {
			continue
		}

		normalizedTemplateID := normalizeTemplateKey(event.Created.TemplateID)
		if normalizedTemplateID == "MCMS.Main:MCMS" {
			mcmsContractID = event.Created.ContractID
			mcmsTemplateID = event.Created.TemplateID
			break
		}
	}

	if mcmsContractID == "" {
		return CantonOpResult[DeployMCMSOutput]{}, fmt.Errorf("failed to find MCMS contract in transaction events")
	}

	fmt.Printf("Deployed MCMS contract   id=%s\n", mcmsContractID)

	return CantonOpResult[DeployMCMSOutput]{
		TransactionID: commandID,
		Output: DeployMCMSOutput{
			MCMSContractID: mcmsContractID,
			MCMSTemplateID: mcmsTemplateID,
		},
	}, nil
}

// ============================================================================
// Operation Declarations
// ============================================================================

var DeployMCMSOp = cld_ops.NewOperation(
	"canton/mcms/deploy",
	semver.MustParse("0.1.0"),
	"Deploys the MCMS contract on Canton",
	deployMCMSHandler,
)
