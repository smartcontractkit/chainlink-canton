package changesets

import (
	"context"
	"fmt"
	"math/big"

	cantonclient "github.com/smartcontractkit/chainlink-canton/deployment/client"
	ccipops "github.com/smartcontractkit/chainlink-canton/deployment/ops/ccip"

	"github.com/noders-team/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/ccip/feequoter"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type ConfigureEVMChainConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	LedgerAPIURL  string `yaml:"ledgerApiUrl"`
	AdminAPIURL   string `yaml:"adminApiUrl"`

	JWTSecret         string `yaml:"jwtSecret"`         // Optional, defaults to "unsafe"
	DeployerParty     string `yaml:"deployerParty"`     // Optional, will allocate if not provided
	DeployerPartyHint string `yaml:"deployerPartyHint"` // Optional hint for party allocation

	// Contract IDs (required - should be retrieved from address book or passed from previous changeset)
	GlobalConfigContractID string `yaml:"globalConfigContractId"`
	FeeQuoterContractID    string `yaml:"feeQuoterContractId"`

	// EVM Chain Configuration
	EVMChainSelector string `yaml:"evmChainSelector"` // EVM chain selector as string

	// GlobalConfig DestChainConfig
	DestChainConfig *DestChainConfigInput `yaml:"destChainConfig"`

	// GlobalConfig SourceChainConfig
	SourceChainConfig *SourceChainConfigInput `yaml:"sourceChainConfig"`

	// FeeQuoter Price Updates
	PriceUpdates *PriceUpdatesInput `yaml:"priceUpdates"`

	// FeeQuoter Fee Token Updates
	FeeTokenUpdates *FeeTokenUpdatesInput `yaml:"feeTokenUpdates"`

	// FeeQuoter Dest Chain Config Updates
	FeeQuoterDestChainConfigUpdates []FeeQuoterDestChainConfigInput `yaml:"feeQuoterDestChainConfigUpdates"`
}

// DestChainConfigInput represents the input structure for DestChainConfig
type DestChainConfigInput struct {
	IsEnabled        bool     `yaml:"isEnabled"`
	DefaultExecutor  string   `yaml:"defaultExecutor"`
	OffRampAddress   string   `yaml:"offRampAddress"`
	LaneMandatedCCVs []string `yaml:"laneMandatedCCVs"`
	DefaultCCVs      []string `yaml:"defaultCCVs"`
}

// SourceChainConfigInput represents the input structure for SourceChainConfig
type SourceChainConfigInput struct {
	IsEnabled        bool     `yaml:"isEnabled"`
	OnRampAddress    string   `yaml:"onRampAddress"`
	LaneMandatedCCVs []string `yaml:"laneMandatedCCVs"`
	DefaultCCVs      []string `yaml:"defaultCCVs"`
}

// PriceUpdatesInput represents the input structure for PriceUpdates
type PriceUpdatesInput struct {
	TokenPriceUpdates []TokenPriceUpdateInput `yaml:"tokenPriceUpdates"`
	GasPriceUpdates   []GasPriceUpdateInput   `yaml:"gasPriceUpdates"`
}

// TokenPriceUpdateInput represents the input structure for TokenPriceUpdate
type TokenPriceUpdateInput struct {
	InstrumentId InstrumentIdInput `yaml:"instrumentId"`
	UsdPerToken  string            `yaml:"usdPerToken"` // Numeric value as string
}

// GasPriceUpdateInput represents the input structure for GasPriceUpdate
type GasPriceUpdateInput struct {
	DestChainSelector string `yaml:"destChainSelector"` // Chain selector as string
	UsdPerUnitGas     string `yaml:"usdPerUnitGas"`     // Numeric value as string
}

// FeeTokenUpdatesInput represents the input structure for FeeTokenUpdates
type FeeTokenUpdatesInput struct {
	FeeTokensToRemove []InstrumentIdInput `yaml:"feeTokensToRemove"`
	FeeTokensToAdd    []FeeTokenArgsInput `yaml:"feeTokensToAdd"`
}

// InstrumentIdInput represents the input structure for InstrumentId
type InstrumentIdInput struct {
	Admin string `yaml:"admin"` // Party ID
	Id    string `yaml:"id"`    // Token identifier
}

// FeeTokenArgsInput represents the input structure for FeeTokenArgs
type FeeTokenArgsInput struct {
	InstrumentId      InstrumentIdInput `yaml:"instrumentId"`
	PremiumMultiplier string            `yaml:"premiumMultiplier"` // Numeric value as string
}

// FeeQuoterDestChainConfigInput represents the input structure for FeeQuoter DestChainConfig
type FeeQuoterDestChainConfigInput struct {
	DestChainSelector string                         `yaml:"destChainSelector"` // Chain selector as string
	DestChainConfig   FeeQuoterDestChainConfig2Input `yaml:"destChainConfig"`
}

// FeeQuoterDestChainConfig2Input represents the input structure for DestChainConfig2
type FeeQuoterDestChainConfig2Input struct {
	IsEnabled                   bool   `yaml:"isEnabled"`
	MaxDataBytes                int64  `yaml:"maxDataBytes"`
	MaxPerMsgGasLimit           int64  `yaml:"maxPerMsgGasLimit"`
	DestGasOverhead             int64  `yaml:"destGasOverhead"`
	DestGasPerPayloadByteBase   int64  `yaml:"destGasPerPayloadByteBase"`
	ChainFamilySelector         string `yaml:"chainFamilySelector"`
	DefaultTxGasLimit           int64  `yaml:"defaultTxGasLimit"`
	NetworkFeeUSD               string `yaml:"networkFeeUSD"`      // Numeric value as string
	DefaultTokenFeeUSD          string `yaml:"defaultTokenFeeUSD"` // Numeric value as string
	DefaultTokenDestGasOverhead int64  `yaml:"defaultTokenDestGasOverhead"`
}

var _ cldf.ChangeSetV2[ConfigureEVMChainConfig] = ConfigureEVMChain{}

// ConfigureEVMChain configures an EVM chain on Canton CCIP contracts
type ConfigureEVMChain struct{}

// Apply implements deployment.ChangeSetV2.
func (c ConfigureEVMChain) Apply(e cldf.Environment, config ConfigureEVMChainConfig) (cldf.ChangesetOutput, error) {
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
	// GlobalConfig Updates
	// --------------------------

	// Update DestChainConfig if provided
	if config.DestChainConfig != nil {
		destChainConfig := common.DestChainConfig{
			IsEnabled:        types.BOOL(config.DestChainConfig.IsEnabled),
			DefaultExecutor:  types.TEXT(config.DestChainConfig.DefaultExecutor),
			OffRampAddress:   types.TEXT(config.DestChainConfig.OffRampAddress),
			LaneMandatedCCVs: make([]types.TEXT, len(config.DestChainConfig.LaneMandatedCCVs)),
			DefaultCCVs:      make([]types.TEXT, len(config.DestChainConfig.DefaultCCVs)),
		}
		for i, ccv := range config.DestChainConfig.LaneMandatedCCVs {
			destChainConfig.LaneMandatedCCVs[i] = types.TEXT(ccv)
		}
		for i, ccv := range config.DestChainConfig.DefaultCCVs {
			destChainConfig.DefaultCCVs[i] = types.TEXT(ccv)
		}

		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.UpdateGlobalConfigDestChainConfigOp, deps, ccipops.UpdateGlobalConfigDestChainConfigInput{
			GlobalConfigContractID: config.GlobalConfigContractID,
			DestChainSelector:      config.EVMChainSelector,
			Config:                 destChainConfig,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to update GlobalConfig dest chain config for EVM chain %s: %w", config.EVMChainSelector, err)
		}
	}

	// Update SourceChainConfig if provided
	if config.SourceChainConfig != nil {
		sourceChainConfig := common.SourceChainConfig{
			IsEnabled:        types.BOOL(config.SourceChainConfig.IsEnabled),
			OnRampAddress:    types.TEXT(config.SourceChainConfig.OnRampAddress),
			LaneMandatedCCVs: make([]types.TEXT, len(config.SourceChainConfig.LaneMandatedCCVs)),
			DefaultCCVs:      make([]types.TEXT, len(config.SourceChainConfig.DefaultCCVs)),
		}
		for i, ccv := range config.SourceChainConfig.LaneMandatedCCVs {
			sourceChainConfig.LaneMandatedCCVs[i] = types.TEXT(ccv)
		}
		for i, ccv := range config.SourceChainConfig.DefaultCCVs {
			sourceChainConfig.DefaultCCVs[i] = types.TEXT(ccv)
		}

		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.UpdateGlobalConfigSourceChainConfigOp, deps, ccipops.UpdateGlobalConfigSourceChainConfigInput{
			GlobalConfigContractID: config.GlobalConfigContractID,
			SourceChainSelector:    config.EVMChainSelector,
			Config:                 sourceChainConfig,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to update GlobalConfig source chain config for EVM chain %s: %w", config.EVMChainSelector, err)
		}
	}

	// --------------------------
	// FeeQuoter Updates
	// --------------------------

	// Update Prices if provided
	if config.PriceUpdates != nil {
		priceUpdates := feequoter.PriceUpdates{
			TokenPriceUpdates: make([]feequoter.TokenPriceUpdate, len(config.PriceUpdates.TokenPriceUpdates)),
			GasPriceUpdates:   make([]feequoter.GasPriceUpdate, len(config.PriceUpdates.GasPriceUpdates)),
		}

		// Convert token price updates
		for i, tpu := range config.PriceUpdates.TokenPriceUpdates {
			// Parse instrument ID
			instrumentId := feequoter.InstrumentId{
				Admin: types.PARTY(tpu.InstrumentId.Admin),
				Id:    types.TEXT(tpu.InstrumentId.Id),
			}

			// Parse USD per token
			usdPerToken, ok := new(big.Int).SetString(tpu.UsdPerToken, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid usdPerToken: %s", tpu.UsdPerToken)
			}
			scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
			usdPerTokenMantissa := new(big.Int).Mul(usdPerToken, scale10)

			priceUpdates.TokenPriceUpdates[i] = feequoter.TokenPriceUpdate{
				InstrumentId: instrumentId,
				UsdPerToken:  types.NUMERIC(usdPerTokenMantissa),
			}
		}

		// Convert gas price updates
		for i, gpu := range config.PriceUpdates.GasPriceUpdates {
			// Parse dest chain selector
			destChainSelector, ok := new(big.Int).SetString(gpu.DestChainSelector, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid destChainSelector: %s", gpu.DestChainSelector)
			}
			scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
			destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

			// Parse USD per unit gas
			usdPerUnitGas, ok := new(big.Int).SetString(gpu.UsdPerUnitGas, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid usdPerUnitGas: %s", gpu.UsdPerUnitGas)
			}
			usdPerUnitGasMantissa := new(big.Int).Mul(usdPerUnitGas, scale10)

			priceUpdates.GasPriceUpdates[i] = feequoter.GasPriceUpdate{
				DestChainSelector: types.NUMERIC(destChainSelectorMantissa),
				UsdPerUnitGas:     types.NUMERIC(usdPerUnitGasMantissa),
			}
		}

		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.UpdateFeeQuoterPricesOp, deps, ccipops.UpdateFeeQuoterPricesInput{
			FeeQuoterContractID: config.FeeQuoterContractID,
			PriceUpdates:        priceUpdates,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to update FeeQuoter prices: %w", err)
		}
	}

	// Apply Fee Token Updates if provided
	if config.FeeTokenUpdates != nil {
		feeTokensToRemove := make([]feequoter.InstrumentId, len(config.FeeTokenUpdates.FeeTokensToRemove))
		for i, ft := range config.FeeTokenUpdates.FeeTokensToRemove {
			feeTokensToRemove[i] = feequoter.InstrumentId{
				Admin: types.PARTY(ft.Admin),
				Id:    types.TEXT(ft.Id),
			}
		}

		feeTokensToAdd := make([]feequoter.FeeTokenArgs, len(config.FeeTokenUpdates.FeeTokensToAdd))
		for i, fta := range config.FeeTokenUpdates.FeeTokensToAdd {
			instrumentId := feequoter.InstrumentId{
				Admin: types.PARTY(fta.InstrumentId.Admin),
				Id:    types.TEXT(fta.InstrumentId.Id),
			}

			premiumMultiplier, ok := new(big.Int).SetString(fta.PremiumMultiplier, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid premiumMultiplier: %s", fta.PremiumMultiplier)
			}
			scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
			premiumMultiplierMantissa := new(big.Int).Mul(premiumMultiplier, scale10)

			feeTokensToAdd[i] = feequoter.FeeTokenArgs{
				InstrumentId:      instrumentId,
				PremiumMultiplier: types.NUMERIC(premiumMultiplierMantissa),
			}
		}

		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.ApplyFeeQuoterFeeTokenUpdatesOp, deps, ccipops.ApplyFeeQuoterFeeTokenUpdatesInput{
			FeeQuoterContractID: config.FeeQuoterContractID,
			FeeTokensToRemove:   feeTokensToRemove,
			FeeTokensToAdd:      feeTokensToAdd,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply FeeQuoter fee token updates: %w", err)
		}
	}

	// Apply Dest Chain Config Updates if provided
	if len(config.FeeQuoterDestChainConfigUpdates) > 0 {
		destChainConfigArgs := make([]feequoter.DestChainConfigArgs, len(config.FeeQuoterDestChainConfigUpdates))
		for i, dccu := range config.FeeQuoterDestChainConfigUpdates {
			// Parse dest chain selector
			destChainSelector, ok := new(big.Int).SetString(dccu.DestChainSelector, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid destChainSelector: %s", dccu.DestChainSelector)
			}
			scale10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil)
			destChainSelectorMantissa := new(big.Int).Mul(destChainSelector, scale10)

			// Parse network fee USD
			networkFeeUSD, ok := new(big.Int).SetString(dccu.DestChainConfig.NetworkFeeUSD, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid networkFeeUSD: %s", dccu.DestChainConfig.NetworkFeeUSD)
			}
			networkFeeUSDMantissa := new(big.Int).Mul(networkFeeUSD, scale10)

			// Parse default token fee USD
			defaultTokenFeeUSD, ok := new(big.Int).SetString(dccu.DestChainConfig.DefaultTokenFeeUSD, 10)
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("invalid defaultTokenFeeUSD: %s", dccu.DestChainConfig.DefaultTokenFeeUSD)
			}
			defaultTokenFeeUSDMantissa := new(big.Int).Mul(defaultTokenFeeUSD, scale10)

			destChainConfigArgs[i] = feequoter.DestChainConfigArgs{
				DestChainSelector: types.NUMERIC(destChainSelectorMantissa),
				DestChainConfig: feequoter.DestChainConfig2{
					IsEnabled:                   types.BOOL(dccu.DestChainConfig.IsEnabled),
					MaxDataBytes:                types.INT64(dccu.DestChainConfig.MaxDataBytes),
					MaxPerMsgGasLimit:           types.INT64(dccu.DestChainConfig.MaxPerMsgGasLimit),
					DestGasOverhead:             types.INT64(dccu.DestChainConfig.DestGasOverhead),
					DestGasPerPayloadByteBase:   types.INT64(dccu.DestChainConfig.DestGasPerPayloadByteBase),
					ChainFamilySelector:         types.TEXT(dccu.DestChainConfig.ChainFamilySelector),
					DefaultTxGasLimit:           types.INT64(dccu.DestChainConfig.DefaultTxGasLimit),
					NetworkFeeUSD:               types.NUMERIC(networkFeeUSDMantissa),
					DefaultTokenFeeUSD:          types.NUMERIC(defaultTokenFeeUSDMantissa),
					DefaultTokenDestGasOverhead: types.INT64(dccu.DestChainConfig.DefaultTokenDestGasOverhead),
				},
			}
		}

		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.ApplyFeeQuoterDestChainConfigUpdatesOp, deps, ccipops.ApplyFeeQuoterDestChainConfigUpdatesInput{
			FeeQuoterContractID: config.FeeQuoterContractID,
			DestChainConfigArgs: destChainConfigArgs,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply FeeQuoter dest chain config updates: %w", err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (c ConfigureEVMChain) VerifyPreconditions(e cldf.Environment, config ConfigureEVMChainConfig) error {
	if config.LedgerAPIURL == "" {
		return fmt.Errorf("ledgerApiUrl is required")
	}
	if config.AdminAPIURL == "" {
		return fmt.Errorf("adminApiUrl is required")
	}
	if config.GlobalConfigContractID == "" {
		return fmt.Errorf("globalConfigContractId is required")
	}
	if config.FeeQuoterContractID == "" {
		return fmt.Errorf("feeQuoterContractId is required")
	}
	if config.EVMChainSelector == "" {
		return fmt.Errorf("evmChainSelector is required")
	}

	return nil
}
