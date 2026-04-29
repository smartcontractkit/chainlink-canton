package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployBurnMintTokenPoolConfig is the config for deploying a BurnMintTokenPool.
// If TokenAdminRegistryInstanceAddress is set, the pool is also registered with that TAR in the same changeset.
type DeployBurnMintTokenPoolConfig struct {
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
	TransferTimeout burnminttokenpool.TransferTimeout
	// Optional; defaults to empty map.
	RemoteChainConfigs map[types.NUMERIC]burnminttokenpool.RemoteChainConfig
	// Optional; defaults to empty map.
	TokenTransferFeeConfigs map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig
	// Optional; zero-value deps if not provided.
	Deps burnminttokenpool.BurnMintTokenPoolDeps
	// If set, the pool is registered with this TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole, SetPool) in the same changeset.
	TokenAdminRegistryInstanceAddress contracts.InstanceAddress
}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployBurnMintTokenPoolConfig]] = DeployBurnMintTokenPool{}

type DeployBurnMintTokenPool struct{}

func (d DeployBurnMintTokenPool) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployBurnMintTokenPoolConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d DeployBurnMintTokenPool) Apply(e cldf.Environment, config CantonCSDeps[DeployBurnMintTokenPoolConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	cfg := config.Config
	poolReceiveContext := cfg.PoolReceiveContext
	if poolReceiveContext.Values == nil {
		poolReceiveContext = common.CCIPContext{Values: map[string]common.AnyValue{}}
	}
	transferTimeout := cfg.TransferTimeout
	if transferTimeout.RelativeHours == nil && transferTimeout.Indefinite == nil {
		transferTimeout = burnminttokenpool.TransferTimeout{RelativeHours: new(types.INT64(24))}
	}
	remoteChainConfigs := cfg.RemoteChainConfigs
	if remoteChainConfigs == nil {
		remoteChainConfigs = map[types.NUMERIC]burnminttokenpool.RemoteChainConfig{}
	}
	tokenTransferFeeConfigs := cfg.TokenTransferFeeConfigs
	if tokenTransferFeeConfigs == nil {
		tokenTransferFeeConfigs = map[types.NUMERIC]burnminttokenpool.TokenTransferFeeConfig{}
	}
	qualifier := new(cfg.Qualifier)
	if cfg.Qualifier == "" {
		qualifier = nil
	}
	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, burn_mint_token_pool.Deploy, chain, contract.DeployInput[burnminttokenpool.BurnMintTokenPool]{
		Qualifier: qualifier,
		Template: burnminttokenpool.BurnMintTokenPool{
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
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save deployed BurnMintTokenPool contract address: %w", err)
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
