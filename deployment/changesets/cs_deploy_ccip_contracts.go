package changesets

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cantonclient "github.com/smartcontractkit/chainlink-canton/deployment/client"
	ccipops "github.com/smartcontractkit/chainlink-canton/deployment/ops/ccip"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

	// CCV Configuration
	CCVStorageLocation     string   `yaml:"ccvStorageLocation"`     // CCV storage location
	CCVVersionTag          string   `yaml:"ccvVersionTag"`          // CCV version tag
	CCVSigners             []string `yaml:"ccvSigners"`             // CCV signers
	CCVMessageSentObserver string   `yaml:"ccvMessageSentObserver"` // CCV message sent observer
	CCVThreshold           int64    `yaml:"ccvThreshold"`           // CCV threshold
}

var _ cldf.ChangeSetV2[DeployCCIPContractsConfig] = DeployCCIPContracts{}

// DeployCCIPContracts deploys all CCIP contracts on Canton
type DeployCCIPContracts struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployCCIPContracts) Apply(e cldf.Environment, config DeployCCIPContractsConfig) (cldf.ChangesetOutput, error) {
	ctx := context.Background()

	// Create datastore and populate it with the deployed contract information
	ds := datastore.NewMemoryDataStore()

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
	}

	// --------------------------
	// CCIP COMMON (GlobalConfig, CCVRegistry)
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPCommonOp, deps, ccipops.DeployCCIPCommonInput{
		InstanceID:         config.InstanceID + "-globalconfig",
		ChainSelectorValue: config.ChainSelectorValue,
		OnRampAddress:      config.OnRampAddress,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Common for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save GlobalConfig contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-globalconfig@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("globalconfig"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save GlobalConfig contract ID: %w", err)
	}

	// Deploy Committee Verifier
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployCommitteeVerifierOp, deps, ccipops.DeployCommitteeVerifierInput{
		InstanceID:          config.InstanceID + "-ccv",
		VersionTag:          config.CCVVersionTag,
		StorageLocation:     config.CCVStorageLocation,
		Threshold:           config.CCVThreshold,
		Signers:             config.CCVSigners,
		MessageSentObserver: config.CCVMessageSentObserver,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CommitteeVerifier for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save Committee Verifier contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-ccv@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("ccv"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Committee Verifier contract ID: %w", err)
	}

	// --------------------------
	// TOKEN ADMIN REGISTRY
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployTokenAdminRegistryOp, deps, ccipops.DeployTokenAdminRegistryInput{
		InstanceID: config.InstanceID + "-tar",
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenAdminRegistry for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save TokenAdminRegistry contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-tar@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("tar"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save TokenAdminRegistry contract ID: %w", err)
	}

	// --------------------------
	// FEE QUOTER
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployFeeQuoterOp, deps, ccipops.DeployFeeQuoterInput{
		InstanceID: config.InstanceID + "-feequoter",
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy FeeQuoter for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save FeeQuoter contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-feequoter@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("feequoter"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save FeeQuoter contract ID: %w", err)
	}

	// --------------------------
	// OFFRAMP
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployOffRampOp, deps, ccipops.DeployOffRampInput{
		InstanceID: config.InstanceID + "-offramp",
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy OffRamp for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save OffRamp contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-offramp@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("offramp"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save OffRamp contract ID: %w", err)
	}

	// --------------------------
	// PER PARTY ROUTER
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployPerPartyRouterOp, deps, ccipops.DeployPerPartyRouterInput{
		InstanceID: config.InstanceID + "-perpartyrouter",
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy PerPartyRouter for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save PerPartyRouter contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-perpartyrouter@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("perpartyrouter"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save PerPartyRouter contract ID: %w", err)
	}

	// --------------------------
	// ONRAMP
	// --------------------------
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.DeployOnRampOp, deps, ccipops.DeployOnRampInput{
		InstanceID:           config.InstanceID + "-onramp",
		DestChainSelector:    config.DestChainSelector,
		DestChainOnRampBytes: config.DestChainOnRampBytes,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy OnRamp for Canton chain %d: %w", config.ChainSelector, err)
	}

	// Save OnRamp contract ID
	err = ds.AddressRefStore.Add(
		datastore.AddressRef{
			ChainSelector: config.ChainSelector,
			Address:       fmt.Sprintf("%s-onramp@%s", config.InstanceID, deps.Party),
			Type:          datastore.ContractType("onramp"),
			Version:       semver.MustParse("1.0.0"),
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save OnRamp contract ID: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   seqReports,
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
