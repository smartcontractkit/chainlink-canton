package lock_release_token_pool

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonLockReleaseTokenPool")

var Version = semver.MustParse("1.0.0")

var lrtpEncoder = lockreleasetokenpool.NewContract("", "CCIP.LockReleaseTokenPool", "LockReleaseTokenPool").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[lockreleasetokenpool.LockReleaseTokenPool]{
	Name:           "canton/ccip/lock_release_token_pool/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP LockReleaseTokenPool contract on Canton",
	Validate: func(template lockreleasetokenpool.LockReleaseTokenPool) error {
		if template.CcipOwner == "" {
			return errors.New("CCIP Owner is required")
		}
		if template.PoolOwner == "" {
			return errors.New("PoolOwner is required")
		}
		if template.InstrumentId == (splice_api_token_holding_v1.InstrumentId{}) {
			return errors.New("instrument ID cannot be empty")
		}
		if template.Decimals < 0 {
			return errors.New("decimals cannot be negative")
		}

		return nil
	},
	PackageName: string(contracts.CCIPLockReleaseTokenPool),
	Prefix:      "lockreleasetokenpool",
})

var SetDynamicConfig = contract.NewExercise(contract.ExerciseParams[lockreleasetokenpool.SetDynamicConfig]{
	Name:         "canton/ccip/lock_release_token_pool/set_dynamic_config",
	Version:      Version,
	Description:  "Sets the Dynamic Config",
	ContractType: ContractType,
	Validate:     nil,
	Template:     lockreleasetokenpool.LockReleaseTokenPool{},
	Method:       lockreleasetokenpool.LockReleaseTokenPool{}.SetDynamicConfig,
	EncodeMethod: lrtpEncoder.SetDynamicConfig,
})

var ApplyChainUpdates = contract.NewExercise(contract.ExerciseParams[lockreleasetokenpool.ApplyChainUpdates]{
	Name:         "canton/ccip/lock_release_token_pool/apply_chain_updates",
	Version:      Version,
	Description:  "Applies remote chain updates to a Canton LockReleaseTokenPool",
	ContractType: ContractType,
	Validate: func(input lockreleasetokenpool.ApplyChainUpdates) error {

		return nil
	},
	Template:     lockreleasetokenpool.LockReleaseTokenPool{},
	Method:       lockreleasetokenpool.LockReleaseTokenPool{}.ApplyChainUpdates,
	EncodeMethod: lrtpEncoder.ApplyChainUpdates,
})

var ApplyTokenTransferFeeConfigUpdates = contract.NewExercise(contract.ExerciseParams[lockreleasetokenpool.ApplyTokenTransferFeeConfigUpdates]{
	Name:         "canton/ccip/lock_release_token_pool/apply_token_transfer_fee_config_updates",
	Version:      Version,
	Description:  "Applies token transfer fee config updates to a Canton LockReleaseTokenPool",
	ContractType: ContractType,
	Validate: func(input lockreleasetokenpool.ApplyTokenTransferFeeConfigUpdates) error {

		return nil
	},
	Template:     lockreleasetokenpool.LockReleaseTokenPool{},
	Method:       lockreleasetokenpool.LockReleaseTokenPool{}.ApplyTokenTransferFeeConfigUpdates,
	EncodeMethod: lrtpEncoder.ApplyTokenTransferFeeConfigUpdates,
})

var SetRateLimitConfig = contract.NewExercise(contract.ExerciseParams[lockreleasetokenpool.SetRateLimitConfig]{
	Name:         "canton/ccip/lock_release_token_pool/set_rate_limit_config",
	Version:      Version,
	Description:  "Sets rate limit configs for a Canton LockReleaseTokenPool",
	ContractType: ContractType,
	Validate: func(input lockreleasetokenpool.SetRateLimitConfig) error {
		if len(input.RateLimitConfigArgs) == 0 {
			return errors.New("rateLimitConfigArgs is required")
		}

		return nil
	},
	Template:     lockreleasetokenpool.LockReleaseTokenPool{},
	Method:       lockreleasetokenpool.LockReleaseTokenPool{}.SetRateLimitConfig,
	EncodeMethod: lrtpEncoder.SetRateLimitConfig,
})
