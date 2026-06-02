package factory

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/types"

	common "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/common"
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CCIPFactory")

var Version = semver.MustParse("0.1.0")

var factoryEncoder = factorybindings.NewContract("", "CCIP.Factory", "CCIPFactory").Encoder()

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
	PackageName: string(contracts.CCIPFactory),
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
		CcvRegistry:        args.Contract.Deps.CcvRegistry,
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
	var payload bytes.Buffer

	writeLenPrefixedText(&payload, args.Contract.InstanceId)
	writeLenPrefixedText(&payload, types.TEXT(args.Contract.Owner))
	writeInt64(&payload, args.Contract.MaxCCVsPerMsg)
	if err := writeRequestedFinality(&payload, args.Contract.DynamicConfig.AllowedFinalityConfig); err != nil {
		return nil, err
	}
	writeBool(&payload, args.Contract.DynamicConfig.CcvAllowlistEnabled)

	return &bind.EncodedChoice{
		Choice:        "DeployExecutor",
		OperationData: hex.EncodeToString(payload.Bytes()),
	}, nil
}

func writeLenPrefixedText(buf *bytes.Buffer, value types.TEXT) {
	writeUvarint(buf, uint64(len(value)))
	buf.WriteString(string(value))
}

func writeUvarint(buf *bytes.Buffer, value uint64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], value)
	buf.Write(scratch[:n])
}

func writeInt64(buf *bytes.Buffer, value types.INT64) {
	_ = binary.Write(buf, binary.BigEndian, int64(value))
}

func writeBool(buf *bytes.Buffer, value types.BOOL) {
	if value {
		buf.WriteByte(0x01)
		return
	}
	buf.WriteByte(0x00)
}

func writeRequestedFinality(buf *bytes.Buffer, finality common.FinalityConfig) error {
	switch {
	case finality.WaitForFinality != nil:
		buf.WriteByte(0x00)
	case finality.WaitForSafe != nil:
		buf.WriteByte(0x01)
	case finality.BlockDepth != nil:
		buf.WriteByte(0x02)
		writeInt64(buf, *finality.BlockDepth)
	default:
		return fmt.Errorf("unsupported executor finality config: %+v", finality)
	}

	return nil
}
