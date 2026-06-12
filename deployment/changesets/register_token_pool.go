package changesets

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

// RegisterTokenPoolConfig registers an already-deployed token pool with the TokenAdminRegistry.
// CantonCSDeps.Participant selects the CCIP owner participant (ProposeAdministrator).
// PoolParticipant selects the pool owner participant (AcceptAdminRole, SetPool); defaults to the CCIP participant when zero.
type RegisterTokenPoolConfig struct {
	CcipOwner      string
	PoolOwner      string
	InstrumentId   splice_api_token_holding_v1.InstrumentId
	PoolInstanceID string
	// PoolParticipant is the participant index for pool-owner steps. When zero, CantonCSDeps.Participant is used.
	PoolParticipant int
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
	poolParticipant := poolParticipantIndex(config)
	if poolParticipant < 0 || poolParticipant >= len(chain.Participants) {
		return fmt.Errorf("pool participant index %d out of range for canton chain %d with %d participants", poolParticipant, config.ChainSelector, len(chain.Participants))
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
		PoolInstanceID:                       cfg.PoolInstanceID,
		CcipParticipantIndex:                 config.Participant,
		PoolParticipantIndex:                 poolParticipantIndex(config),
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("register token pool with TAR: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []cld_ops.Report[any, any]{},
	}, nil
}

func poolParticipantIndex(config CantonCSDeps[RegisterTokenPoolConfig]) int {
	if config.Config.PoolParticipant != 0 {
		return config.Config.PoolParticipant
	}

	return config.Participant
}
