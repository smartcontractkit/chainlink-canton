package deployment

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
	"github.com/smartcontractkit/chainlink-canton-internal/deployment/ops/ccip"
)

func TestCantonChainStateGenerateView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Setup Canton client
	setupResult, err := client.Setup(ctx, client.Config{
		LedgerAPIURL:      "participant1.grpc-ledger-api.localhost:8080",
		AdminAPIURL:       "participant1.admin-api.localhost:8080",
		JWTSecret:         "unsafe",
		DeployerParty:     "",
		DeployerPartyHint: "ledger-api-user",
	})
	require.NoError(t, err, "Failed to setup Canton client")
	t.Cleanup(func() { setupResult.BindingClient.Close() })

	deps := ccip.CantonOpDeps{
		BindingClient: setupResult.BindingClient,
		Party:         setupResult.Party,
		UserID:        setupResult.UserID,
	}

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	instanceID := "test-state-view-instance"
	chainSelectorValue := "1111111111"
	onRampAddress := "0000000000000000000000000000000000000001"
	destChainSelector := "2222222222"
	sourceChainSelector := "3333333333"
	chainSelector := uint64(1111111111)

	var globalConfigContractID string
	var globalConfigTemplateID string
	var onRampContractID string

	// --------------------------
	// Step 1: Deploy Contracts
	// --------------------------
	t.Run("DeployContracts", func(t *testing.T) {
		// Deploy GlobalConfig
		commonResult, err := cld_ops.ExecuteOperation(bundle, ccip.DeployCCIPCommonOp, deps, ccip.DeployCCIPCommonInput{
			InstanceID:         instanceID,
			ChainSelectorValue: chainSelectorValue,
			OnRampAddress:      onRampAddress,
		})
		require.NoError(t, err, "failed to deploy GlobalConfig")
		require.NotEmpty(t, commonResult.Output.Output.GlobalConfigContractID, "GlobalConfig contract ID should not be empty")
		globalConfigContractID = commonResult.Output.Output.GlobalConfigContractID
		globalConfigTemplateID = commonResult.Output.Output.GlobalConfigTemplateID
		t.Logf("Deployed GlobalConfig contract ID: %s", globalConfigContractID)

		// Deploy OnRamp
		onRampResult, err := cld_ops.ExecuteOperation(bundle, ccip.DeployOnRampOp, deps, ccip.DeployOnRampInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy OnRamp")
		require.NotEmpty(t, onRampResult.Output.Output.OnRampContractID, "OnRamp contract ID should not be empty")
		onRampContractID = onRampResult.Output.Output.OnRampContractID
		t.Logf("Deployed OnRamp contract ID: %s", onRampContractID)
	})

	// --------------------------
	// Step 2: Configure GlobalConfig
	// --------------------------
	t.Run("ConfigureGlobalConfig", func(t *testing.T) {
		// Add destination chain config
		destChainConfig := common.DestChainConfig{
			IsEnabled:        types.BOOL(true),
			DefaultExecutor:  types.TEXT("executor-party"),
			OffRampAddress:   types.TEXT("0000000000000000000000000000000000000002"),
			LaneMandatedCCVs: []types.TEXT{types.TEXT("ccv1"), types.TEXT("ccv2")},
			DefaultCCVs:      []types.TEXT{types.TEXT("default-ccv")},
		}

		destResult, err := cld_ops.ExecuteOperation(bundle, ccip.UpdateGlobalConfigDestChainConfigOp, deps, ccip.UpdateGlobalConfigDestChainConfigInput{
			GlobalConfigContractID: globalConfigContractID,
			GlobalConfigTemplateID: globalConfigTemplateID,
			DestChainSelector:      destChainSelector,
			Config:                 destChainConfig,
		})
		require.NoError(t, err, "failed to update dest chain config")
		require.NotEmpty(t, destResult.Output.Output.NewGlobalConfigContractID, "New GlobalConfig contract ID should not be empty")
		globalConfigContractID = destResult.Output.Output.NewGlobalConfigContractID // Update to new contract ID
		t.Logf("Updated GlobalConfig dest chain config, new contract ID: %s", globalConfigContractID)

		// Add source chain config
		sourceChainConfig := common.SourceChainConfig{
			IsEnabled:        types.BOOL(true),
			OnRampAddress:    types.TEXT("0000000000000000000000000000000000000003"),
			LaneMandatedCCVs: []types.TEXT{types.TEXT("ccv3")},
			DefaultCCVs:      []types.TEXT{types.TEXT("default-source-ccv")},
		}

		sourceResult, err := cld_ops.ExecuteOperation(bundle, ccip.UpdateGlobalConfigSourceChainConfigOp, deps, ccip.UpdateGlobalConfigSourceChainConfigInput{
			GlobalConfigContractID: globalConfigContractID,
			GlobalConfigTemplateID: globalConfigTemplateID,
			SourceChainSelector:    sourceChainSelector,
			Config:                 sourceChainConfig,
		})
		require.NoError(t, err, "failed to update source chain config")
		require.NotEmpty(t, sourceResult.Output.Output.NewGlobalConfigContractID, "New GlobalConfig contract ID should not be empty")
		globalConfigContractID = sourceResult.Output.Output.NewGlobalConfigContractID // Update to new contract ID
		t.Logf("Updated GlobalConfig source chain config, new contract ID: %s", globalConfigContractID)
	})

	// --------------------------
	// Step 3: Generate Chain View
	// --------------------------
	t.Run("GenerateChainView", func(t *testing.T) {
		// Create chain state
		chainState := CantonChainState{
			BindingClient:          setupResult.BindingClient,
			Party:                  deps.Party,
			GlobalConfigContractID: globalConfigContractID,
			OnRampContractID:       onRampContractID,
		}

		// Create environment
		env := &cldf.Environment{
			Logger: logger.Test(t),
		}

		// Generate view
		chainView, err := chainState.GenerateView(env, chainSelector, "test-canton-chain")
		require.NoError(t, err, "failed to generate chain view")

		// Verify chain selector
		require.Equal(t, chainSelector, chainView.ChainSelector, "Chain selector should match")

		// --------------------------
		// Verify GlobalConfig View
		// --------------------------
		require.NotEmpty(t, chainView.GlobalConfig.Address, "GlobalConfig address should not be empty")
		require.Equal(t, globalConfigContractID, chainView.GlobalConfig.Address, "GlobalConfig contract ID should match")
		require.Equal(t, deps.Party, chainView.GlobalConfig.Owner, "GlobalConfig owner should match deployer party")
		require.Equal(t, deps.Party, chainView.GlobalConfig.CcipOwner, "GlobalConfig ccipOwner should match deployer party")
		require.Equal(t, instanceID, chainView.GlobalConfig.InstanceId, "GlobalConfig instanceId should match")
		require.Equal(t, chainSelectorValue, chainView.GlobalConfig.ChainSelector, "GlobalConfig chainSelector should match")
		require.Equal(t, onRampAddress, chainView.GlobalConfig.OnRampAddress, "GlobalConfig onRampAddress should match")

		// Verify destination chain config
		require.NotNil(t, chainView.GlobalConfig.DestChainConfigs, "DestChainConfigs should not be nil")
		require.Contains(t, chainView.GlobalConfig.DestChainConfigs, destChainSelector, "DestChainConfigs should contain dest chain selector")

		destConfig, ok := chainView.GlobalConfig.DestChainConfigs[destChainSelector]
		require.True(t, ok, "DestChainConfig should exist for dest chain selector")
		require.True(t, destConfig.IsEnabled, "DestChainConfig should be enabled")
		require.Equal(t, "executor-party", destConfig.DefaultExecutor, "DefaultExecutor should match")
		require.Equal(t, "0000000000000000000000000000000000000002", destConfig.OffRampAddress, "OffRampAddress should match")
		require.Len(t, destConfig.LaneMandatedCCVs, 2, "LaneMandatedCCVs should have 2 items")
		require.Contains(t, destConfig.LaneMandatedCCVs, "ccv1", "LaneMandatedCCVs should contain ccv1")
		require.Contains(t, destConfig.LaneMandatedCCVs, "ccv2", "LaneMandatedCCVs should contain ccv2")
		require.Len(t, destConfig.DefaultCCVs, 1, "DefaultCCVs should have 1 item")
		require.Contains(t, destConfig.DefaultCCVs, "default-ccv", "DefaultCCVs should contain default-ccv")

		// Verify source chain config
		require.NotNil(t, chainView.GlobalConfig.SourceChainConfigs, "SourceChainConfigs should not be nil")
		require.Contains(t, chainView.GlobalConfig.SourceChainConfigs, sourceChainSelector, "SourceChainConfigs should contain source chain selector")

		sourceConfig, ok := chainView.GlobalConfig.SourceChainConfigs[sourceChainSelector]
		require.True(t, ok, "SourceChainConfig should exist for source chain selector")
		require.True(t, sourceConfig.IsEnabled, "SourceChainConfig should be enabled")
		require.Equal(t, "0000000000000000000000000000000000000003", sourceConfig.OnRampAddress, "OnRampAddress should match")
		require.Len(t, sourceConfig.LaneMandatedCCVs, 1, "LaneMandatedCCVs should have 1 item")
		require.Contains(t, sourceConfig.LaneMandatedCCVs, "ccv3", "LaneMandatedCCVs should contain ccv3")
		require.Len(t, sourceConfig.DefaultCCVs, 1, "DefaultCCVs should have 1 item")
		require.Contains(t, sourceConfig.DefaultCCVs, "default-source-ccv", "DefaultCCVs should contain default-source-ccv")

		// --------------------------
		// Verify OnRamp View
		// --------------------------
		require.NotEmpty(t, chainView.OnRamp.Address, "OnRamp address should not be empty")
		require.Equal(t, onRampContractID, chainView.OnRamp.Address, "OnRamp contract ID should match")
		require.Equal(t, deps.Party, chainView.OnRamp.Owner, "OnRamp owner should match deployer party")
		require.Equal(t, deps.Party, chainView.OnRamp.CcipOwner, "OnRamp ccipOwner should match deployer party")
		require.Equal(t, instanceID, chainView.OnRamp.InstanceId, "OnRamp instanceId should match")

		t.Logf("Successfully generated chain view")
		t.Logf("  Chain Selector: %d", chainView.ChainSelector)
		t.Logf("  GlobalConfig Contract ID: %s", chainView.GlobalConfig.Address)
		t.Logf("  GlobalConfig Dest Chain Configs: %d", len(chainView.GlobalConfig.DestChainConfigs))
		t.Logf("  GlobalConfig Source Chain Configs: %d", len(chainView.GlobalConfig.SourceChainConfigs))
		t.Logf("  OnRamp Contract ID: %s", chainView.OnRamp.Address)

		// Write view to JSON file
		jsonData, err := json.MarshalIndent(chainView, "", "  ")
		require.NoError(t, err, "failed to marshal chain view to JSON")

		// Create output directory if it doesn't exist
		outputDir := "test-output"
		err = os.MkdirAll(outputDir, 0755)
		require.NoError(t, err, "failed to create output directory")

		// Write to file
		outputFile := filepath.Join(outputDir, "canton-chain-view.json")
		err = os.WriteFile(outputFile, jsonData, 0644)
		require.NoError(t, err, "failed to write JSON file")
		t.Logf("Chain view written to: %s", outputFile)
	})
}
