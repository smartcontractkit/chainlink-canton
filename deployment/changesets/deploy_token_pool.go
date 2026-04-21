package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	lock_release_token_pool "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployTokenPoolConfig is the config for deploying a LockReleaseTokenPool.
// If TokenAdminRegistryInstanceAddress is set, the pool is also registered with that TAR in the same changeset.
type DeployTokenPoolConfig struct {
	CcipOwner    string
	PoolOwner    string
	InstrumentId splice_api_token_holding_v1.InstrumentId
	Decimals     int64
	// Optional explicit instance ID for the pool.
	// If empty, deploy operation generates one.
	InstanceID string
	// Qualifier is optional (e.g. token symbol) for AddressRef and idempotency.
	Qualifier string
	// Optional; defaults to empty. PoolReceiveContext can be set for receive context.
	PoolReceiveContext common.CCIPContext
	// Optional; defaults to 24h RelativeHours. TransferTimeout for the pool.
	TransferTimeout lockreleasetokenpool.TransferTimeout
	// Optional; defaults to empty map.
	RemoteChainConfigs types.GENMAP
	// Optional; defaults to empty map.
	TokenTransferFeeConfigs types.GENMAP
	// Optional; zero-value deps if not provided.
	Deps lockreleasetokenpool.LockReleaseTokenPoolDeps
	// If set, the pool is registered with this TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole, SetPool) in the same changeset.
	TokenAdminRegistryInstanceAddress contracts.InstanceAddress
}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployTokenPoolConfig]] = DeployTokenPool{}

type DeployTokenPool struct{}

func (d DeployTokenPool) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployTokenPoolConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d DeployTokenPool) Apply(e cldf.Environment, config CantonCSDeps[DeployTokenPoolConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	cfg := config.Config
	poolReceiveContext := cfg.PoolReceiveContext
	if poolReceiveContext.Values == nil {
		poolReceiveContext = common.CCIPContext{Values: types.TEXTMAP{}}
	}
	transferTimeout := cfg.TransferTimeout
	if transferTimeout.RelativeHours == nil && transferTimeout.Indefinite == nil {
		h := types.INT64(24)
		transferTimeout = lockreleasetokenpool.TransferTimeout{RelativeHours: &h}
	}
	remoteChainConfigs := cfg.RemoteChainConfigs
	if remoteChainConfigs == nil {
		remoteChainConfigs = types.GENMAP{}
	}
	tokenTransferFeeConfigs := cfg.TokenTransferFeeConfigs
	if tokenTransferFeeConfigs == nil {
		tokenTransferFeeConfigs = types.GENMAP{}
	}
	qualifier := new(cfg.Qualifier)
	if cfg.Qualifier == "" {
		qualifier = nil
	}
	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, lock_release_token_pool.Deploy, chain, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
		Qualifier: qualifier,
		Template: lockreleasetokenpool.LockReleaseTokenPool{
			CcipOwner:               types.PARTY(cfg.CcipOwner),
			PoolOwner:               types.PARTY(cfg.PoolOwner),
			InstanceId:              types.TEXT(cfg.InstanceID),
			InstrumentId:            cfg.InstrumentId,
			Decimals:                types.INT64(cfg.Decimals),
			RemoteChainConfigs:      remoteChainConfigs,
			TokenTransferFeeConfigs: tokenTransferFeeConfigs,
			PoolReceiveContext:      poolReceiveContext,
			TransferTimeout:         transferTimeout,
			Deps:                    cfg.Deps,
		},
		OwnerParty: types.PARTY(cfg.PoolOwner),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if len(out.Output.Labels.List()) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing raw lock/release pool label in deploy output")
	}
	rawPoolAddr, err := contracts.RawInstanceAddressFromString(out.Output.Labels.List()[0])
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("parse raw lock/release pool label: %w", err)
	}

	if err = ds.AddressRefStore.Add(out.Output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save deployed LockReleaseTokenPool contract address: %w", err)
	}

	if cfg.TokenAdminRegistryInstanceAddress != (contracts.InstanceAddress{}) {
		regInput := sequences.RegisterTokenPoolInput{
			TokenAdminRegistryInstanceAddress: cfg.TokenAdminRegistryInstanceAddress,
			InstrumentId:                      cfg.InstrumentId,
			CcipParty:                         cfg.CcipOwner,
			PoolOwnerParty:                    cfg.PoolOwner,
			PoolInstanceID:                    rawPoolAddr.InstanceID(),
		}
		_, err = cld_ops.ExecuteSequence(e.OperationsBundle, sequences.RegisterTokenPool, chain, regInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token pool with TAR: %w", err)
		}
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []cld_ops.Report[any, any]{},
	}, nil
}
