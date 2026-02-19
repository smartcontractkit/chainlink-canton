package changesets

import (
	"fmt"

	"github.com/aws/smithy-go/ptr"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
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
	// Qualifier is optional (e.g. token symbol) for AddressRef and idempotency.
	Qualifier string
	// Optional; defaults to empty. ChainCCVRequirements can be set for chain-specific CCV requirements.
	ChainCCVRequirements types.GENMAP
	// Optional; defaults to empty. PoolReceiveContext can be set for receive context.
	PoolReceiveContext splice_api_token_metadata_v1.ChoiceContext
	// Optional; defaults to 24h RelativeHours. TransferTimeout for the pool.
	TransferTimeout lockreleasetokenpool.TransferTimeout
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
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d DeployTokenPool) Apply(e cldf.Environment, config CantonCSDeps[DeployTokenPoolConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()

	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	deps := dependencies.CantonDeps{
		Chain: chain,
	}

	cfg := config.Config
	chainCCVReqs := cfg.ChainCCVRequirements
	if chainCCVReqs == nil {
		chainCCVReqs = types.GENMAP{}
	}
	poolReceiveContext := cfg.PoolReceiveContext
	if poolReceiveContext.Values == nil {
		poolReceiveContext = splice_api_token_metadata_v1.ChoiceContext{Values: types.TEXTMAP{}}
	}
	transferTimeout := cfg.TransferTimeout
	if transferTimeout.RelativeHours == nil && transferTimeout.Indefinite == nil {
		h := types.INT64(24)
		transferTimeout = lockreleasetokenpool.TransferTimeout{RelativeHours: &h}
	}

	template := lockreleasetokenpool.LockReleaseTokenPool{
		CcipOwner:            types.PARTY(cfg.CcipOwner),
		PoolOwner:            types.PARTY(cfg.PoolOwner),
		InstanceId:           "", // set by deploy operation
		InstrumentId:         cfg.InstrumentId,
		Decimals:             types.INT64(cfg.Decimals),
		ChainCCVRequirements: chainCCVReqs,
		PoolReceiveContext:   poolReceiveContext,
		TransferTimeout:      transferTimeout,
	}

	qualifier := ptr.String(cfg.Qualifier)
	if cfg.Qualifier == "" {
		qualifier = nil
	}

	out, err := cld_ops.ExecuteOperation(e.OperationsBundle, lock_release_token_pool.Deploy, deps, contract.DeployInput[lockreleasetokenpool.LockReleaseTokenPool]{
		ChainSelector: config.ChainSelector,
		Qualifier:     qualifier,
		ActAs:         []string{cfg.PoolOwner},
		Template:      template,
		OwnerParty:    types.PARTY(cfg.PoolOwner),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply DeployTokenPool operation: %w", err)
	}

	if err = ds.AddressRefStore.Add(out.Output); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save deployed LockReleaseTokenPool contract address: %w", err)
	}

	regInput := sequences.RegisterTokenPoolInput{
		TokenAdminRegistryInstanceAddress: cfg.TokenAdminRegistryInstanceAddress,
		InstrumentId: splice_api_token_holding_v1.InstrumentId{
			Admin: cfg.InstrumentId.Admin,
			Id:    cfg.InstrumentId.Id,
		},
		CcipParty:      cfg.CcipOwner,
		PoolOwnerParty: cfg.PoolOwner,
		PoolInstanceID: out.Output.Address,
	}
	_, err = cld_ops.ExecuteSequence(e.OperationsBundle, sequences.RegisterTokenPool, deps, regInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register token pool with TAR: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []cld_ops.Report[any, any]{},
	}, nil
}
