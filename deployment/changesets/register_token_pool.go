package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

// RegisterTokenPoolConfig registers an already-deployed token pool with the TokenAdminRegistry.
// CantonCSDeps.Participant selects the CCIP owner participant for all TAR registration steps.
// PoolAdmin is TokenConfig admin; PoolOwner is recorded in PoolRegistration.
type RegisterTokenPoolConfig struct {
	CcipOwner      string
	PoolOwner      string
	PoolAdmin      string
	InstrumentId   splice_api_token_holding_v1.InstrumentId
	PoolInstanceID string
}

var _ cldf.ChangeSetV2[CantonCSDeps[RegisterTokenPoolConfig]] = RegisterTokenPool{}

type RegisterTokenPool struct{}

func (r RegisterTokenPool) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[RegisterTokenPoolConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}
	if config.Config.PoolInstanceID == "" {
		return fmt.Errorf("PoolInstanceID is required")
	}
	if config.Config.CcipOwner == "" {
		return fmt.Errorf("CcipOwner is required")
	}
	if config.Config.PoolOwner == "" {
		return fmt.Errorf("PoolOwner is required")
	}
	if config.Config.PoolAdmin == "" {
		return fmt.Errorf("PoolAdmin is required")
	}

	return nil
}

func (r RegisterTokenPool) Apply(e cldf.Environment, config CantonCSDeps[RegisterTokenPoolConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	cfg := config.Config

	tarRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
	}
	tarRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(tarRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve token admin registry raw address: %w", err)
	}

	_, err = cld_ops.ExecuteSequence(e.OperationsBundle, sequences.RegisterTokenPool, chain, sequences.RegisterTokenPoolInput{
		TokenAdminRegistryInstanceAddress:    contracts.HexToInstanceAddress(tarRef.Address),
		TokenAdminRegistryRawInstanceAddress: tarRaw,
		InstrumentId:                         cfg.InstrumentId,
		CcipParty:                            cfg.CcipOwner,
		PoolOwnerParty:                       cfg.PoolOwner,
		PoolAdminParty:                       cfg.PoolAdmin,
		PoolInstanceID:                       cfg.PoolInstanceID,
		CcipParticipantIndex:                 config.Participant,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("register token pool with TAR: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []cld_ops.Report[any, any]{},
	}, nil
}
