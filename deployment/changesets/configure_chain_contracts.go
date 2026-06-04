package changesets

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// MCMSTimelockConfig configures MCMS timelock proposals instead of direct execution.
type MCMSTimelockConfig struct {
	MinDelay         time.Duration
	Description      string
	OverridePrevRoot bool
	Action           mcms_types.TimelockAction
	// MCMSContract describes the MCMS contract for on-chain state queries.
	MCMSContract cantonmcms.MCMSContractInfo
	// Role is the MCMS role to use (proposer, canceller, bypasser).
	Role cantonsdk.TimelockRole
}

// ConfigureGlobalConfigConfig holds the parameters for ConfigureGlobalConfig changeset.
type ConfigureGlobalConfigConfig struct {
	InstanceAddress contracts.InstanceAddress
	// RawInstanceAddress is the "instanceId@partyId" format needed for MCMS proposals.
	RawInstanceAddress string
	DestChainUpdates   []common.DestChainConfigArgs
	SourceChainUpdates []common.SourceChainConfigArgs
	TimelockConfig     *MCMSTimelockConfig
}

var _ cldf.ChangeSetV2[CantonCSDeps[ConfigureGlobalConfigConfig]] = ConfigureGlobalConfig{}

type ConfigureGlobalConfig struct{}

func (d ConfigureGlobalConfig) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[ConfigureGlobalConfigConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d ConfigureGlobalConfig) Apply(e cldf.Environment, config CantonCSDeps[ConfigureGlobalConfigConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	mcmsEnabled := config.Config.TimelockConfig != nil

	var exerciseOutputs []opcontract.ExerciseOutput

	if len(config.Config.DestChainUpdates) > 0 {
		out, err := operations.ExecuteOperation(e.OperationsBundle, global_config.ApplyDestChainConfigUpdates, chain, opcontract.ChoiceInput[common.ApplyDestChainConfigUpdates]{
			InstanceAddress:    config.Config.InstanceAddress,
			RawInstanceAddress: config.Config.RawInstanceAddress,
			Args: common.ApplyDestChainConfigUpdates{
				DestChainConfigUpdates: config.Config.DestChainUpdates,
			},
			MCMSEnabled: mcmsEnabled,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply dest chain config updates: %w", err)
		}
		exerciseOutputs = append(exerciseOutputs, out.Output)
	}

	if len(config.Config.SourceChainUpdates) > 0 {
		out, err := operations.ExecuteOperation(e.OperationsBundle, global_config.ApplySourceChainConfigUpdates, chain, opcontract.ChoiceInput[common.ApplySourceChainConfigUpdates]{
			InstanceAddress:    config.Config.InstanceAddress,
			RawInstanceAddress: config.Config.RawInstanceAddress,
			Args: common.ApplySourceChainConfigUpdates{
				SourceChainConfigUpdates: config.Config.SourceChainUpdates,
			},
			MCMSEnabled: mcmsEnabled,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply source chain config updates: %w", err)
		}
		exerciseOutputs = append(exerciseOutputs, out.Output)
	}

	if mcmsEnabled {
		batchOp, err := cantonmcms.BuildBatchFromOutputs(exerciseOutputs)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build batch operation: %w", err)
		}
		if len(batchOp.Transactions) == 0 {
			return cldf.ChangesetOutput{}, nil
		}

		tlCfg := config.Config.TimelockConfig
		participant := chain.Participants[0]

		proposal, err := cantonmcms.GenerateTimelockProposal(
			e.GetContext(),
			participant.LedgerServices.State,
			participant.PartyID,
			cantonmcms.ProposalConfig{
				MCMSContract:         tlCfg.MCMSContract,
				ChainSelector:        mcms_types.ChainSelector(config.ChainSelector),
				Description:          tlCfg.Description,
				MinDelay:             tlCfg.MinDelay,
				OverridePreviousRoot: tlCfg.OverridePrevRoot,
				Action:               tlCfg.Action,
				Role:                 tlCfg.Role,
			},
			[]mcms_types.BatchOperation{batchOp},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate timelock proposal: %w", err)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
