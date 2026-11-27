package changesets

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// UpdateCantonTokenPoolFinalityConfigConfig configures UpdateCantonTokenPoolFinalityConfig.
type UpdateCantonTokenPoolFinalityConfigConfig struct {
	// TokenPoolRef points at the deployed Canton token pool in the datastore. The pool type is
	// detected from the ref, supporting BurnMintTokenPool, LockReleaseTokenPool, and
	// RegistryBurnMintTokenPool (and their Canton-prefixed canonical names).
	TokenPoolRef datastore.AddressRef `json:"tokenPoolRef" yaml:"tokenPoolRef"`
	// RemoteChainSelectors lists the lanes whose finality config is updated. Each lane must
	// already be configured on the pool via ApplyChainUpdates.
	RemoteChainSelectors []uint64 `json:"remoteChainSelectors" yaml:"remoteChainSelectors"`
	// AllowedFinalityConfig is the finality config set on every listed lane, e.g. blockDepth: 1
	// to allow fast transfer finality requests. Exactly one of waitForFinality, waitForSafe, or
	// blockDepth must be set.
	AllowedFinalityConfig finality.Config `json:"allowedFinalityConfig" yaml:"allowedFinalityConfig"`
	// MinDelay is the MCMS timelock min delay (proposal-driven runs only).
	MinDelay time.Duration `json:"minDelay" yaml:"minDelay"`
	// Description is the MCMS proposal description (proposal-driven runs only).
	Description string `json:"description" yaml:"description"`
}

var _ cldf.ChangeSetV2[CantonCSDeps[UpdateCantonTokenPoolFinalityConfigConfig]] = UpdateCantonTokenPoolFinalityConfig{}

type UpdateCantonTokenPoolFinalityConfig struct{}

func (u UpdateCantonTokenPoolFinalityConfig) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[UpdateCantonTokenPoolFinalityConfigConfig]) error {
	cfg := config.Config
	if len(cfg.RemoteChainSelectors) == 0 {
		return fmt.Errorf("at least one remote chain selector is required")
	}
	if cfg.TokenPoolRef.Address == "" {
		return fmt.Errorf("tokenPoolRef with an address is required")
	}
	if _, err := sequences.ResolveCantonTokenPoolTypeFromRef(cfg.TokenPoolRef.Type); err != nil {
		return err
	}
	if err := sequences.ValidateCantonFinalityConfig(cfg.AllowedFinalityConfig); err != nil {
		return err
	}

	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants",
			config.Participant, config.ChainSelector, len(chain.Participants))
	}
	// MCMS proposals are only used when the operator cannot act for the pool owner directly.
	if len(chain.Participants[config.Participant].ReadAsPartyIDs) > 0 {
		if err := requireCCIPOwnerMCMSRef(e, config.ChainSelector); err != nil {
			return err
		}
	}

	return nil
}

func (u UpdateCantonTokenPoolFinalityConfig) Apply(e cldf.Environment, config CantonCSDeps[UpdateCantonTokenPoolFinalityConfigConfig]) (cldf.ChangesetOutput, error) {
	cfg := config.Config
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participantIdx := config.Participant
	proposalDriven := len(chain.Participants[participantIdx].ReadAsPartyIDs) > 0

	out, err := cld_ops.ExecuteSequence(e.OperationsBundle, sequences.UpdateTokenPoolFinalityConfig, e.BlockChains, sequences.UpdateTokenPoolFinalityInput{
		CantonChainSelector:   config.ChainSelector,
		TokenPoolRef:          cfg.TokenPoolRef,
		RemoteChainSelectors:  cfg.RemoteChainSelectors,
		AllowedFinalityConfig: cfg.AllowedFinalityConfig,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("update canton token pool finality config: %w", err)
	}

	if proposalDriven && len(out.Output.BatchOps) == 0 {
		e.Logger.Infof(
			"lanes %v on token pool %s already allow the requested finality; no MCMS proposal generated",
			cfg.RemoteChainSelectors, cfg.TokenPoolRef.Address,
		)

		return cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}, nil
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, participantIdx, proposalDriven,
		cantonmcms.QualifierCCIPOwner, cfg.MinDelay, cfg.Description,
		datastore.NewMemoryDataStore(), out.Output.BatchOps,
	)
}
