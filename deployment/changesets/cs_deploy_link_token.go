package changesets

import (
	"context"
	"fmt"

	cantonclient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
	linkops "github.com/smartcontractkit/chainlink-canton-internal/deployment/ops/link"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployLinkTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	LedgerAPIURL  string `yaml:"ledgerApiUrl"`
	AdminAPIURL   string `yaml:"adminApiUrl"`

	JWTSecret         string `yaml:"jwtSecret"`         // Optional, defaults to "unsafe"
	DeployerParty     string `yaml:"deployerParty"`     // Optional, will allocate if not provided
	DeployerPartyHint string `yaml:"deployerPartyHint"` // Optional hint for party allocation
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

// DeployLinkToken deploys LINK token CoinRegistry contract on Canton
type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	// TODO; This will be abstracted to CLD / CLDF
	// Setup Canton client
	setupResult, err := cantonclient.Setup(ctx, cantonclient.Config{
		LedgerAPIURL:      config.LedgerAPIURL,
		AdminAPIURL:       config.AdminAPIURL,
		JWTSecret:         config.JWTSecret,
		DeployerParty:     config.DeployerParty,
		DeployerPartyHint: config.DeployerPartyHint,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to setup Canton client: %w", err)
	}
	defer setupResult.BindingClient.Close()

	// Create Canton operation dependencies
	deps := linkops.CantonOpDeps{
		BindingClient: setupResult.BindingClient,
		Party:         setupResult.Party,
		UserID:        setupResult.UserID,
	}

	// Run DeployLinkToken Operation
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Extract output
	// deployResult, ok := result.Output.(linkops.CantonOpResult[linkops.DeployLinkTokenOutput])
	// if !ok {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("unexpected output type from DeployLINKOp")
	// }

	// Save LinkToken registry contract ID to the addressbook
	// TODO: Define proper type and version constants
	// typeAndVersionLinkToken := cldf.NewTypeAndVersion("CantonLinkTokenRegistry", "1.0.0")
	// err = ab.Save(config.ChainSelector, deployResult.Output.RegistryContractID, typeAndVersionLinkToken)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken registry contract ID %s for Canton chain %d: %w", deployResult.Output.RegistryContractID, config.ChainSelector, err)
	// }

	// Add report

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployLinkToken) VerifyPreconditions(e cldf.Environment, config DeployLinkTokenConfig) error {
	if config.LedgerAPIURL == "" {
		return fmt.Errorf("ledgerApiUrl is required")
	}
	if config.AdminAPIURL == "" {
		return fmt.Errorf("adminApiUrl is required")
	}

	return nil
}
