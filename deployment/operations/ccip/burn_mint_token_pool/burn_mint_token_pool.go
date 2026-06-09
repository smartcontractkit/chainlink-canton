package burn_mint_token_pool

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonBurnMintTokenPool")

var Version = semver.MustParse("2.0.0")

var lrtpEncoder = burnminttokenpool.NewContract("", "CCIP.BurnMintTokenPool", "BurnMintTokenPool").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[burnminttokenpool.BurnMintTokenPool]{
	Name:           "canton/ccip/burn_mint_token_pool/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP BurnMintTokenPool contract on Canton",
	Validate: func(template burnminttokenpool.BurnMintTokenPool) error {
		if template.CcipOwner == "" {
			return errors.New("CCIP Owner is required")
		}
		if template.PoolOwner == "" {
			return errors.New("PoolOwner is required")
		}
		if template.InstrumentId == (splice_api_token_holding_v1.InstrumentId{}) {
			return errors.New("instrument ID cannot be empty")
		}
		if err := tokenpool.ValidateTokenDecimals(int64(template.Decimals)); err != nil {
			return err
		}

		return nil
	},
	PackageName: string(contracts.CCIPBurnMintTokenPool),
	Prefix:      "burnminttokenpool",
})

var SetDynamicConfig = contract.NewExercise(contract.ExerciseParams[burnminttokenpool.SetDynamicConfig]{
	Name:         "canton/ccip/burn_mint_token_pool/set_dynamic_config",
	Version:      Version,
	Description:  "Sets the Dynamic Config",
	ContractType: ContractType,
	Validate:     nil,
	Template:     burnminttokenpool.BurnMintTokenPool{},
	Method:       burnminttokenpool.BurnMintTokenPool{}.SetDynamicConfig,
	EncodeMethod: lrtpEncoder.SetDynamicConfig,
})

var ApplyChainUpdates = contract.NewExercise(contract.ExerciseParams[burnminttokenpool.ApplyChainUpdates]{
	Name:         "canton/ccip/burn_mint_token_pool/apply_chain_updates",
	Version:      Version,
	Description:  "Applies remote chain updates to a Canton BurnMintTokenPool",
	ContractType: ContractType,
	Validate: func(input burnminttokenpool.ApplyChainUpdates) error {

		return nil
	},
	Template:     burnminttokenpool.BurnMintTokenPool{},
	Method:       burnminttokenpool.BurnMintTokenPool{}.ApplyChainUpdates,
	EncodeMethod: lrtpEncoder.ApplyChainUpdates,
})

var ApplyTokenTransferFeeConfigUpdates = contract.NewExercise(contract.ExerciseParams[burnminttokenpool.ApplyTokenTransferFeeConfigUpdates]{
	Name:         "canton/ccip/burn_mint_token_pool/apply_token_transfer_fee_config_updates",
	Version:      Version,
	Description:  "Applies token transfer fee config updates to a Canton BurnMintTokenPool",
	ContractType: ContractType,
	Validate: func(input burnminttokenpool.ApplyTokenTransferFeeConfigUpdates) error {

		return nil
	},
	Template:     burnminttokenpool.BurnMintTokenPool{},
	Method:       burnminttokenpool.BurnMintTokenPool{}.ApplyTokenTransferFeeConfigUpdates,
	EncodeMethod: lrtpEncoder.ApplyTokenTransferFeeConfigUpdates,
})

var SetRateLimiterReferences = contract.NewExercise(contract.ExerciseParams[burnminttokenpool.SetRateLimiterReferences]{
	Name:         "canton/ccip/burn_mint_token_pool/set_rate_limiter_references",
	Version:      Version,
	Description:  "Updates which RateLimiter contract identities a Canton BurnMintTokenPool references per remote chain",
	ContractType: ContractType,
	Validate: func(input burnminttokenpool.SetRateLimiterReferences) error {
		if len(input.RateLimitConfigArgs) == 0 {
			return errors.New("rateLimitConfigArgs is required")
		}

		return nil
	},
	Template:     burnminttokenpool.BurnMintTokenPool{},
	Method:       burnminttokenpool.BurnMintTokenPool{}.SetRateLimiterReferences,
	EncodeMethod: lrtpEncoder.SetRateLimiterReferences,
})

var SetRateLimitConfig = contract.NewExercise(contract.ExerciseParams[burnminttokenpool.SetRateLimitConfigParams]{
	Name:         "canton/ccip/burn_mint_token_pool/set_rate_limit_config",
	Version:      Version,
	Description:  "Tunes capacity, rate, and isEnabled on a RateLimiter referenced by a Canton BurnMintTokenPool",
	ContractType: ContractType,
	Validate: func(input burnminttokenpool.SetRateLimitConfigParams) error {
		if input.RateLimiterInstanceAddress.Unpack == "" {
			return errors.New("rateLimiterInstanceAddress is required")
		}

		return nil
	},
	Template:     burnminttokenpool.BurnMintTokenPool{},
	EncodeMethod: lrtpEncoder.SetRateLimitConfigParams,
})
