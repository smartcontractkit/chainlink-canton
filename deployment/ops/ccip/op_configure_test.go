package ccip

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/ccip/feequoter"
	compileClient "github.com/smartcontractkit/chainlink-canton-internal/deployment/client"
)

func TestConfigureCCIPContracts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	setupResult, err := compileClient.Setup(ctx, compileClient.Config{
		LedgerAPIURL:      "participant1.grpc-ledger-api.localhost:8080",
		AdminAPIURL:       "participant1.admin-api.localhost:8080",
		JWTSecret:         "unsafe",
		DeployerParty:     "", // Empty to use primary party or allocate new one
		DeployerPartyHint: "ledger-api-user",
	})
	require.NoError(t, err, "Failed to setup Canton client")

	t.Cleanup(setupResult.BindingClient.Close)

	deps := CantonOpDeps{
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

	instanceID := "test-ccip-instance"
	chainSelectorValue := "1111111111"
	onRampAddress := "0000000000000000000000000000000000000001"
	evmChainSelector := "2222222222"

	// --------------------------
	// Deploy contracts first (required for configuration)
	// --------------------------
	var globalConfigContractID string
	var feeQuoterContractID string

	t.Run("DeployContracts", func(t *testing.T) {
		t.Parallel()

		// Deploy GlobalConfig
		commonResult, err := cld_ops.ExecuteOperation(bundle, DeployCCIPCommonOp, deps, DeployCCIPCommonInput{
			InstanceID:         instanceID,
			ChainSelectorValue: chainSelectorValue,
			OnRampAddress:      onRampAddress,
		})
		require.NoError(t, err, "failed to deploy CCIP Common")
		globalConfigContractID = commonResult.Output.Output.GlobalConfigContractID
		require.NotEmpty(t, globalConfigContractID, "GlobalConfig contract ID should not be empty")
		t.Logf("Deployed GlobalConfig contract ID: %s", globalConfigContractID)

		// Deploy FeeQuoter
		feeQuoterResult, err := cld_ops.ExecuteOperation(bundle, DeployFeeQuoterOp, deps, DeployFeeQuoterInput{
			InstanceID: instanceID,
		})
		require.NoError(t, err, "failed to deploy FeeQuoter")
		feeQuoterContractID = feeQuoterResult.Output.Output.FeeQuoterContractID
		require.NotEmpty(t, feeQuoterContractID, "FeeQuoter contract ID should not be empty")
		t.Logf("Deployed FeeQuoter contract ID: %s", feeQuoterContractID)
	})

	// --------------------------
	// Test GlobalConfig Configuration
	// --------------------------

	t.Run("UpdateGlobalConfigDestChainConfig", func(t *testing.T) {
		t.Parallel()

		destChainConfig := common.DestChainConfig{
			IsEnabled:        types.BOOL(true),
			DefaultExecutor:  types.TEXT("executor-party"),
			OffRampAddress:   types.TEXT("0000000000000000000000000000000000000002"),
			LaneMandatedCCVs: []types.TEXT{types.TEXT("ccv1"), types.TEXT("ccv2")},
			DefaultCCVs:      []types.TEXT{types.TEXT("default-ccv")},
		}

		result, err := cld_ops.ExecuteOperation(bundle, UpdateGlobalConfigDestChainConfigOp, deps, UpdateGlobalConfigDestChainConfigInput{
			GlobalConfigContractID: globalConfigContractID,
			DestChainSelector:      evmChainSelector,
			Config:                 destChainConfig,
		})
		require.NoError(t, err, "failed to update GlobalConfig dest chain config")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, result.Output.Output.NewGlobalConfigContractID, "new GlobalConfig contract ID should not be empty")
		globalConfigContractID = result.Output.Output.NewGlobalConfigContractID

		t.Logf("Updated GlobalConfig dest chain config, tx=%s newCID=%s",
			result.Output.Output.TransactionID,
			result.Output.Output.NewGlobalConfigContractID,
		)
	})

	t.Run("UpdateGlobalConfigSourceChainConfig", func(t *testing.T) {
		t.Parallel()

		sourceChainConfig := common.SourceChainConfig{
			IsEnabled:        types.BOOL(true),
			OnRampAddress:    types.TEXT(onRampAddress),
			LaneMandatedCCVs: []types.TEXT{types.TEXT("ccv1")},
			DefaultCCVs:      []types.TEXT{types.TEXT("default-ccv")},
		}

		result, err := cld_ops.ExecuteOperation(bundle, UpdateGlobalConfigSourceChainConfigOp, deps, UpdateGlobalConfigSourceChainConfigInput{
			GlobalConfigContractID: globalConfigContractID, // fresh created contractID
			SourceChainSelector:    evmChainSelector,
			Config:                 sourceChainConfig,
		})
		require.NoError(t, err, "failed to update GlobalConfig source chain config")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")
		// carry forward the new CID returned by the update op
		require.NotEmpty(t, result.Output.Output.NewGlobalConfigContractID, "new GlobalConfig contract ID should not be empty")
		globalConfigContractID = result.Output.Output.NewGlobalConfigContractID
		t.Logf("Updated GlobalConfig source chain config, new contract ID: %s", globalConfigContractID)
	})

	// --------------------------
	// Test FeeQuoter Configuration
	// --------------------------

	t.Run("UpdateFeeQuoterPrices", func(t *testing.T) {
		t.Parallel()

		// Create test instrument ID
		instrumentId := feequoter.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    types.TEXT("test-token"),
		}

		// Scale USD per token to NUMERIC(10) mantissa
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
		usdPerToken := new(big.Int).SetInt64(100) // $1.00 (scaled)
		usdPerTokenMantissa := new(big.Int).Mul(usdPerToken, scale10)

		// Scale dest chain selector to NUMERIC(10) mantissa
		destChainSelector, _ := new(big.Int).SetString(evmChainSelector, 10)
		destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

		// Scale USD per unit gas to NUMERIC(10) mantissa
		usdPerUnitGas := new(big.Int).SetInt64(50) // $0.50 (scaled)
		usdPerUnitGasMantissa := new(big.Int).Mul(usdPerUnitGas, scale10)

		priceUpdates := feequoter.PriceUpdates{
			TokenPriceUpdates: []feequoter.TokenPriceUpdate{
				{
					InstrumentId: instrumentId,
					UsdPerToken:  types.NUMERIC(usdPerTokenMantissa),
				},
			},
			GasPriceUpdates: []feequoter.GasPriceUpdate{
				{
					DestChainSelector: types.NUMERIC(destChainSelectorMantissa),
					UsdPerUnitGas:     types.NUMERIC(usdPerUnitGasMantissa),
				},
			},
		}

		result, err := cld_ops.ExecuteOperation(bundle, UpdateFeeQuoterPricesOp, deps, UpdateFeeQuoterPricesInput{
			FeeQuoterContractID: feeQuoterContractID,
			PriceUpdates:        priceUpdates,
		})
		require.NoError(t, err, "failed to update FeeQuoter prices")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, result.Output.Output.NewFeeQuoterContractID, "new FeeQuoter contract ID should not be empty")
		feeQuoterContractID = result.Output.Output.NewFeeQuoterContractID
		t.Logf("Updated FeeQuoter prices, new contract ID: %s", feeQuoterContractID)
	})

	t.Run("ApplyFeeQuoterFeeTokenUpdates", func(t *testing.T) {
		t.Parallel()

		// Create test instrument IDs
		instrumentId1 := feequoter.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    types.TEXT("token-1"),
		}

		instrumentId2 := feequoter.InstrumentId{
			Admin: types.PARTY(deps.Party),
			Id:    types.TEXT("token-2"),
		}

		// Scale premium multiplier to NUMERIC(10) mantissa
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
		premiumMultiplier := new(big.Int).SetInt64(120) // 1.2x (scaled)
		premiumMultiplierMantissa := new(big.Int).Mul(premiumMultiplier, scale10)

		feeTokensToAdd := []feequoter.FeeTokenArgs{
			{
				InstrumentId:      instrumentId1,
				PremiumMultiplier: types.NUMERIC(premiumMultiplierMantissa),
			},
			{
				InstrumentId:      instrumentId2,
				PremiumMultiplier: types.NUMERIC(premiumMultiplierMantissa),
			},
		}

		result, err := cld_ops.ExecuteOperation(bundle, ApplyFeeQuoterFeeTokenUpdatesOp, deps, ApplyFeeQuoterFeeTokenUpdatesInput{
			FeeQuoterContractID: feeQuoterContractID,
			FeeTokensToRemove:   []feequoter.InstrumentId{}, // Empty for this test
			FeeTokensToAdd:      feeTokensToAdd,
		})
		require.NoError(t, err, "failed to apply FeeQuoter fee token updates")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, result.Output.Output.NewFeeQuoterContractID, "new FeeQuoter contract ID should not be empty")
		feeQuoterContractID = result.Output.Output.NewFeeQuoterContractID
		t.Logf("Applied FeeQuoter fee token updates, new contract ID: %s", feeQuoterContractID)
	})

	t.Run("ApplyFeeQuoterDestChainConfigUpdates", func(t *testing.T) {
		t.Parallel()

		// Scale dest chain selector to NUMERIC(10) mantissa
		scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
		destChainSelector, _ := new(big.Int).SetString(evmChainSelector, 10)
		destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

		// Scale network fee USD to NUMERIC(10) mantissa
		networkFeeUSD := new(big.Int).SetInt64(100) // $1.00 (scaled)
		networkFeeUSDMantissa := new(big.Int).Mul(networkFeeUSD, scale10)

		// Scale default token fee USD to NUMERIC(10) mantissa
		defaultTokenFeeUSD := new(big.Int).SetInt64(50) // $0.50 (scaled)
		defaultTokenFeeUSDMantissa := new(big.Int).Mul(defaultTokenFeeUSD, scale10)

		destChainConfigArgs := []feequoter.DestChainConfigArgs{
			{
				DestChainSelector: types.NUMERIC(destChainSelectorMantissa),
				DestChainConfig: feequoter.DestChainConfig2{
					IsEnabled:                   types.BOOL(true),
					MaxDataBytes:                types.INT64(5000),
					MaxPerMsgGasLimit:           types.INT64(2000000),
					DestGasOverhead:             types.INT64(30000),
					DestGasPerPayloadByteBase:   types.INT64(12),
					ChainFamilySelector:         types.TEXT("2812d52c"),
					DefaultTxGasLimit:           types.INT64(500000),
					NetworkFeeUSD:               types.NUMERIC(networkFeeUSDMantissa),
					DefaultTokenFeeUSD:          types.NUMERIC(defaultTokenFeeUSDMantissa),
					DefaultTokenDestGasOverhead: types.INT64(0),
				},
			},
		}

		result, err := cld_ops.ExecuteOperation(bundle, ApplyFeeQuoterDestChainConfigUpdatesOp, deps, ApplyFeeQuoterDestChainConfigUpdatesInput{
			FeeQuoterContractID: feeQuoterContractID,
			DestChainConfigArgs: destChainConfigArgs,
		})
		require.NoError(t, err, "failed to apply FeeQuoter dest chain config updates")
		require.NotEmpty(t, result.Output.Output.TransactionID, "Transaction ID should not be empty")

		// carry forward the new CID returned by the update op
		require.NotEmpty(t, result.Output.Output.NewFeeQuoterContractID, "new FeeQuoter contract ID should not be empty")
		feeQuoterContractID = result.Output.Output.NewFeeQuoterContractID
		t.Logf("Applied FeeQuoter dest chain config updates, new contract ID: %s", feeQuoterContractID)
	})
}
