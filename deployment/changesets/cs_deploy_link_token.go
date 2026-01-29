package changesets

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cantonclient "github.com/smartcontractkit/chainlink-canton/deployment/client"
	linkops "github.com/smartcontractkit/chainlink-canton/deployment/ops/link"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

	InstanceID string `yaml:"instanceId"` // Instance ID for LINK token registry
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

// DeployLinkToken deploys LINK token CoinRegistry contract on Canton
type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()

	// Create datastore and populate it with the deployed contract information
	ds := datastore.NewMemoryDataStore()

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

	// Add report to datastore
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-linktokenregistry@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("linktokenregistry"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken registry contract ID: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   seqReports,
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
