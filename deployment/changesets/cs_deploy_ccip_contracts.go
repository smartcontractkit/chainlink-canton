package changesets

import (
	"context"
	"fmt"

	cantonclient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
	ccipops "github.com/smartcontractkit/chainlink-canton-internal/deployment/ops/ccip"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployCCIPContractsConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	LedgerAPIURL  string `yaml:"ledgerApiUrl"`
	AdminAPIURL   string `yaml:"adminApiUrl"`

	JWTSecret         string `yaml:"jwtSecret"`         // Optional, defaults to "unsafe"
	DeployerParty     string `yaml:"deployerParty"`     // Optional, will allocate if not provided
	DeployerPartyHint string `yaml:"deployerPartyHint"` // Optional hint for party allocation

	// CCIP Configuration
	InstanceID           string `yaml:"instanceId"`           // Instance ID for CCIP contracts
	ChainSelectorValue   string `yaml:"chainSelectorValue"`   // Chain selector as string (Numeric 0)
	DestChainSelector    string `yaml:"destChainSelector"`    // Destination chain selector
	OnRampAddress        string `yaml:"onRampAddress"`        // OnRamp address for GlobalConfig
	DestChainOnRampBytes []byte `yaml:"destChainOnRampBytes"` // Destination chain onramp address bytes
}

var _ cldf.ChangeSetV2[DeployCCIPContractsConfig] = DeployCCIPContracts{}

// DeployCCIPContracts deploys all CCIP contracts on Canton
type DeployCCIPContracts struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployCCIPContracts) Apply(e cldf.Environment, config DeployCCIPContractsConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

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
	deps := ccipops.CantonOpDeps{
		BindingClient: setupResult.BindingClient,
		Party:         setupResult.Party,
		UserID:        setupResult.UserID,
	}

	// --------------------------
	// CCIP COMMON (GlobalConfig, CCVRegistry)
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPCommonOp, deps, ccipops.DeployCCIPCommonInput{
		InstanceID:         config.InstanceID,
		ChainSelectorValue: config.ChainSelectorValue,
		OnRampAddress:      config.OnRampAddress,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Common for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save GlobalConfig contract ID
	// TODO: Define proper type and version constants
	// typeAndVersionGlobalConfig := cldf.NewTypeAndVersion("CantonCCIPGlobalConfig", "1.0.0")
	// err = ab.Save(config.ChainSelector, commonReport.Output.GlobalConfigContractID, typeAndVersionGlobalConfig)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save GlobalConfig contract ID: %w", err)
	// }

	// --------------------------
	// TOKEN ADMIN REGISTRY
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployTokenAdminRegistryOp, deps, ccipops.DeployTokenAdminRegistryInput{
		InstanceID: config.InstanceID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenAdminRegistry for Canton chain %d: %w", config.ChainSelector, err)
	}

	// --------------------------
	// FEE QUOTER
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployFeeQuoterOp, deps, ccipops.DeployFeeQuoterInput{
		InstanceID: config.InstanceID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy FeeQuoter for Canton chain %d: %w", config.ChainSelector, err)
	}

	// --------------------------
	// OFFRAMP
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployOffRampOp, deps, ccipops.DeployOffRampInput{
		InstanceID: config.InstanceID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy OffRamp for Canton chain %d: %w", config.ChainSelector, err)
	}

	// --------------------------
	// PER PARTY ROUTER
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployPerPartyRouterOp, deps, ccipops.DeployPerPartyRouterInput{
		InstanceID: config.InstanceID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy PerPartyRouter for Canton chain %d: %w", config.ChainSelector, err)
	}

	// --------------------------
	// ONRAMP
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployOnRampOp, deps, ccipops.DeployOnRampInput{
		InstanceID:           config.InstanceID,
		DestChainSelector:    config.DestChainSelector,
		DestChainOnRampBytes: config.DestChainOnRampBytes,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy OnRamp for Canton chain %d: %w", config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployCCIPContracts) VerifyPreconditions(e cldf.Environment, config DeployCCIPContractsConfig) error {
	if config.LedgerAPIURL == "" {
		return fmt.Errorf("ledgerApiUrl is required")
	}
	if config.AdminAPIURL == "" {
		return fmt.Errorf("adminApiUrl is required")
	}
	if config.InstanceID == "" {
		return fmt.Errorf("instanceId is required")
	}
	if config.ChainSelectorValue == "" {
		return fmt.Errorf("chainSelectorValue is required")
	}
	return nil
}
