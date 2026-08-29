package factory

import (
	"errors"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CCIPFactory")

var Version = semver.MustParse("2.0.0")

var factoryEncoder = factorybindings.NewContract("", "CCIP.FactoryV2", "CCIPFactory").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[factorybindings.CCIPFactory]{
	Name:           "canton/ccip/factory/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIPFactory contract on Canton",
	Validate: func(template factorybindings.CCIPFactory) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.McmsParty == "" {
			return errors.New("mcmsParty cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPFactoryV2),
	Prefix:      "factory",
})

var DeployRMNRemote = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployRMNRemote]{
	Name:         "canton/ccip/factory/deploy_rmn_remote",
	Version:      Version,
	Description:  "Deploys an RMNRemote through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployRMNRemote,
	EncodeMethod: encodeDeployRMNRemote,
})

var DeployGlobalConfig = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployGlobalConfig]{
	Name:         "canton/ccip/factory/deploy_global_config",
	Version:      Version,
	Description:  "Deploys a GlobalConfig through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployGlobalConfig,
	EncodeMethod: encodeDeployGlobalConfig,
})

var DeployTokenAdminRegistry = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployTokenAdminRegistry]{
	Name:         "canton/ccip/factory/deploy_token_admin_registry",
	Version:      Version,
	Description:  "Deploys a TokenAdminRegistry through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployTokenAdminRegistry,
	EncodeMethod: encodeDeployTokenAdminRegistry,
})

var DeployFeeQuoter = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployFeeQuoter]{
	Name:         "canton/ccip/factory/deploy_fee_quoter",
	Version:      Version,
	Description:  "Deploys a FeeQuoter through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployFeeQuoter,
	EncodeMethod: encodeDeployFeeQuoter,
})

var DeployLinkToken = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployLinkToken]{
	Name:         "canton/ccip/factory/deploy_link_token",
	Version:      Version,
	Description:  "Deploys a LinkRegistry token through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployLinkToken,
	EncodeMethod: encodeDeployLinkToken,
})

var DeployCommitteeVerifier = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployCommitteeVerifier]{
	Name:         "canton/ccip/factory/deploy_committee_verifier",
	Version:      Version,
	Description:  "Deploys a CommitteeVerifier through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployCommitteeVerifier,
	EncodeMethod: encodeDeployCommitteeVerifier,
})

var DeployOffRamp = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployOffRamp]{
	Name:         "canton/ccip/factory/deploy_offramp",
	Version:      Version,
	Description:  "Deploys an OffRamp through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployOffRamp,
	EncodeMethod: encodeDeployOffRamp,
})

var DeployOnRamp = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployOnRamp]{
	Name:         "canton/ccip/factory/deploy_onramp",
	Version:      Version,
	Description:  "Deploys an OnRamp through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployOnRamp,
	EncodeMethod: encodeDeployOnRamp,
})

var DeployPerPartyRouterFactory = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployPerPartyRouterFactory]{
	Name:         "canton/ccip/factory/deploy_per_party_router_factory",
	Version:      Version,
	Description:  "Deploys a PerPartyRouterFactory through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployPerPartyRouterFactory,
	EncodeMethod: encodeDeployPerPartyRouterFactory,
})

var SetOwnerToMCMS = contract.NewExercise(contract.ExerciseParams[factorybindings.SetOwnerToMCMS]{
	Name:         "canton/ccip/factory/set_owner_to_mcms",
	Version:      Version,
	Description:  "Transfers CCIPFactory ownership from bootstrap owner to MCMS party",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.SetOwnerToMCMS,
	EncodeMethod: factoryEncoder.SetOwnerToMCMS,
})

var DeployExecutor = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployExecutor]{
	Name:         "canton/ccip/factory/deploy_executor",
	Version:      Version,
	Description:  "Deploys an Executor through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployExecutor,
	EncodeMethod: encodeDeployExecutor,
})

var DeployLockReleaseTokenPool = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployLockReleaseTokenPool]{
	Name:         "canton/ccip/factory/deploy_lock_release_token_pool",
	Version:      Version,
	Description:  "Deploys a LockReleaseTokenPool through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployLockReleaseTokenPool,
	EncodeMethod: encodeDeployLockReleaseTokenPool,
})

var DeployBurnMintTokenPool = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployBurnMintTokenPool]{
	Name:         "canton/ccip/factory/deploy_burn_mint_token_pool",
	Version:      Version,
	Description:  "Deploys a BurnMintTokenPool through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployBurnMintTokenPool,
	EncodeMethod: encodeDeployBurnMintTokenPool,
})

var DeployRateLimiter = contract.NewExercise(contract.ExerciseParams[factorybindings.DeployRateLimiter]{
	Name:         "canton/ccip/factory/deploy_rate_limiter",
	Version:      Version,
	Description:  "Deploys a token pool rate limiter through the CCIPFactory",
	ContractType: ContractType,
	Template:     factorybindings.CCIPFactory{},
	Method:       factorybindings.CCIPFactory{}.DeployRateLimiter,
	EncodeMethod: encodeDeployRateLimiter,
})

func encodeDeployRMNRemote(args factorybindings.DeployRMNRemote) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployRMNRemoteParams(factorybindings.DeployRMNRemoteParams{
		InstanceId:      args.Contract.InstanceId,
		RmnOwner:        args.Contract.RmnOwner,
		CcipOwner:       args.Contract.CcipOwner,
		CustomObservers: args.Contract.CustomObservers,
		CursedSubjects:  args.Contract.CursedSubjects,
	})
}

func encodeDeployGlobalConfig(args factorybindings.DeployGlobalConfig) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployGlobalConfigParams(factorybindings.DeployGlobalConfigParams{
		InstanceId:    args.Contract.InstanceId,
		ChainSelector: args.Contract.ChainSelector,
	})
}

func encodeDeployTokenAdminRegistry(args factorybindings.DeployTokenAdminRegistry) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployTokenAdminRegistryParams(factorybindings.DeployTokenAdminRegistryParams{
		InstanceId: args.Contract.InstanceId,
	})
}

func encodeDeployFeeQuoter(args factorybindings.DeployFeeQuoter) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployFeeQuoterParams(factorybindings.DeployFeeQuoterParams{
		InstanceId:            args.Contract.InstanceId,
		LinkTokenInstrumentId: args.Contract.LinkTokenInstrumentId,
	})
}

func encodeDeployLinkToken(args factorybindings.DeployLinkToken) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployLinkTokenParams(factorybindings.DeployLinkTokenParams{
		InstanceId:   args.Contract.InstanceId,
		InstrumentId: args.Contract.RegistryInstrumentId,
	})
}

func encodeDeployCommitteeVerifier(args factorybindings.DeployCommitteeVerifier) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployCommitteeVerifierParams(factorybindings.DeployCommitteeVerifierParams{
		InstanceId:                   args.Contract.InstanceId,
		Owner:                        args.Contract.Owner,
		CcipOwner:                    args.Contract.CcipOwner,
		VersionTag:                   args.Contract.VersionTag,
		AllowListAdmin:               args.Contract.AllowListAdmin,
		MessageSentObservers:         args.Contract.MessageSentObservers,
		RmnRemote:                    args.Contract.Deps.RmnRemote,
		StorageLocations:             args.Contract.StorageLocations,
		StorageLocationsAdmin:        args.Contract.StorageLocationsAdmin,
		PendingStorageLocationsAdmin: args.Contract.PendingStorageLocationsAdmin,
	})
}

func encodeDeployOffRamp(args factorybindings.DeployOffRamp) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployOffRampParams(factorybindings.DeployOffRampParams{
		InstanceId:         args.Contract.InstanceId,
		GlobalConfig:       args.Contract.Deps.GlobalConfig,
		RmnRemote:          args.Contract.Deps.RmnRemote,
		TokenAdminRegistry: args.Contract.Deps.TokenAdminRegistry,
	})
}

func encodeDeployOnRamp(args factorybindings.DeployOnRamp) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployOnRampParams(factorybindings.DeployOnRampParams{
		InstanceId:         args.Contract.InstanceId,
		GlobalConfig:       args.Contract.Deps.GlobalConfig,
		RmnRemote:          args.Contract.Deps.RmnRemote,
		TokenAdminRegistry: args.Contract.Deps.TokenAdminRegistry,
		FeeQuoter:          args.Contract.Deps.FeeQuoter,
		MaxUSDCentsPerMsg:  args.Contract.MaxUSDCentsPerMsg,
	})
}

func encodeDeployPerPartyRouterFactory(args factorybindings.DeployPerPartyRouterFactory) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployPerPartyRouterFactoryParams(factorybindings.DeployPerPartyRouterFactoryParams{
		InstanceId:         args.Contract.InstanceId,
		OnRamp:             args.Contract.Deps.OnRamp,
		OffRamp:            args.Contract.Deps.OffRamp,
		GlobalConfig:       args.Contract.Deps.GlobalConfig,
		TokenAdminRegistry: args.Contract.Deps.TokenAdminRegistry,
		FeeQuoter:          args.Contract.Deps.FeeQuoter,
		RmnRemote:          args.Contract.Deps.RmnRemote,
	})
}

func encodeDeployLockReleaseTokenPool(args factorybindings.DeployLockReleaseTokenPool) (*bind.EncodedChoice, error) {
	c := args.Contract
	return factoryEncoder.DeployLockReleaseTokenPoolParams(factorybindings.DeployLockReleaseTokenPoolParams{
		InstanceId:         c.InstanceId,
		PoolOwner:          c.PoolOwner,
		CcipOwner:          c.CcipOwner,
		InstrumentId:       c.InstrumentId,
		Decimals:           c.Decimals,
		RateLimitAdmin:     c.RateLimitAdmin,
		TokenAdminRegistry: c.Deps.TokenAdminRegistry,
		FeeQuoter:          c.Deps.FeeQuoter,
		RmnRemote:          c.Deps.RmnRemote,
		PoolReceiveContext: c.PoolReceiveContext,
		TransferTimeout:    c.TransferTimeout,
	})
}

func encodeDeployBurnMintTokenPool(args factorybindings.DeployBurnMintTokenPool) (*bind.EncodedChoice, error) {
	c := args.Contract
	return factoryEncoder.DeployBurnMintTokenPoolParams(factorybindings.DeployBurnMintTokenPoolParams{
		InstanceId:         c.InstanceId,
		PoolOwner:          c.PoolOwner,
		CcipOwner:          c.CcipOwner,
		InstrumentId:       c.InstrumentId,
		Decimals:           c.Decimals,
		RateLimitAdmin:     c.RateLimitAdmin,
		TokenAdminRegistry: c.Deps.TokenAdminRegistry,
		FeeQuoter:          c.Deps.FeeQuoter,
		RmnRemote:          c.Deps.RmnRemote,
		PoolReceiveContext: c.PoolReceiveContext,
		TransferTimeout:    c.TransferTimeout,
	})
}

func encodeDeployRateLimiter(args factorybindings.DeployRateLimiter) (*bind.EncodedChoice, error) {
	c := args.Contract
	return factoryEncoder.DeployRateLimiterParams(factorybindings.DeployRateLimiterParams{
		InstanceId:          c.InstanceId,
		PoolInstanceId:      c.PoolInstanceId,
		PoolOwner:           c.PoolOwner,
		RemoteChainSelector: c.RemoteChainSelector,
		Direction:           c.Direction,
		Mode:                c.Mode,
		IsEnabled:           c.IsEnabled,
		Capacity:            c.Capacity,
		Rate:                c.Rate,
	})
}

func encodeDeployExecutor(args factorybindings.DeployExecutor) (*bind.EncodedChoice, error) {
	return factoryEncoder.DeployExecutorParams(factorybindings.DeployExecutorParams{
		InstanceId:            args.Contract.InstanceId,
		Owner:                 args.Contract.Owner,
		MaxCCVsPerMsg:         args.Contract.MaxCCVsPerMsg,
		AllowedFinalityConfig: args.Contract.DynamicConfig.AllowedFinalityConfig,
		CcvAllowlistEnabled:   args.Contract.DynamicConfig.CcvAllowlistEnabled,
	})
}
