package changesets

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// UpdateMCMSConfig updates signer, quorum and parent configuration for one or more
// roles on an existing Canton MCMS instance, resolved from the datastore by qualifier.
//
// The SetConfig choice is controlled by the MCMS owner party, so this changeset only
// supports direct execution by a participant that can ActAs that party (devnet, test
// environments, or before owner-party authorization is locked down). Proposal-driven
// updates via MCMS self-dispatch (ScheduleBatch → ExecuteScheduledBatch) are not
// supported yet.
type UpdateMCMSConfigParams struct {
	// Qualifier selects the MCMS instance in the datastore. Defaults to the CLL qualifier.
	Qualifier   string                 `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
	RoleConfigs []MCMSRoleConfigParams `json:"roleConfigs" yaml:"roleConfigs"`
}

type UpdateMCMSConfigConfig struct {
	Params UpdateMCMSConfigParams `json:"params" yaml:"params"`
}

type UpdateMCMSConfig struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[UpdateMCMSConfigConfig]] = UpdateMCMSConfig{}

func (u UpdateMCMSConfig) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[UpdateMCMSConfigConfig]) error {
	params := config.Config.Params
	if len(params.RoleConfigs) == 0 {
		return fmt.Errorf("roleConfigs is required")
	}
	for _, roleConfig := range params.RoleConfigs {
		if _, err := buildNormalizedConfig(roleConfig.Config); err != nil {
			return fmt.Errorf("build MCMS config for role %s: %w", roleConfig.Role, err)
		}
	}

	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	qualifier := qualifierOrDefault(params.Qualifier)
	if _, err := dsutils.ProposerMCMSAddressRef(e.DataStore, config.ChainSelector, qualifier); err != nil {
		return fmt.Errorf("MCMS instance for qualifier %q must be deployed first: %w", qualifier, err)
	}
	rawInstanceAddress, err := dsutils.MCMSRawInstanceAddress(e.DataStore, config.ChainSelector, qualifier)
	if err != nil {
		return fmt.Errorf("resolve MCMS raw instance address for qualifier %q: %w", qualifier, err)
	}

	participant := chain.Participants[config.Participant]
	if opcontract.ProposalDrivenForCaller(participant, rawInstanceAddress.Owner()) {
		return fmt.Errorf(
			"participant %s cannot ActAs MCMS owner party %s required by the SetConfig choice; proposal-driven config updates are not supported yet",
			participant.PartyID, rawInstanceAddress.Owner(),
		)
	}

	return nil
}

func (u UpdateMCMSConfig) Apply(e cldf.Environment, config CantonCSDeps[UpdateMCMSConfigConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	rawInstanceAddress, err := dsutils.MCMSRawInstanceAddress(e.DataStore, config.ChainSelector, qualifierOrDefault(config.Config.Params.Qualifier))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve MCMS raw instance address: %w", err)
	}

	_, err = operations.ExecuteSequence(e.OperationsBundle, updateMCMSConfigSequence, chain, updateMCMSConfigInput{
		RawInstanceAddress: rawInstanceAddress,
		ParticipantIndex:   config.Participant,
		RoleConfigs:        config.Config.Params.RoleConfigs,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute UpdateMCMSConfig sequence: %w", err)
	}

	// SetConfig archives and recreates the MCMS contract at the same instance address,
	// so no new address refs are produced.
	return cldf.ChangesetOutput{
		DataStore: datastore.NewMemoryDataStore(),
		Reports:   []operations.Report[any, any]{},
	}, nil
}

type updateMCMSConfigInput struct {
	RawInstanceAddress contracts.RawInstanceAddress `json:"rawInstanceAddress"`
	ParticipantIndex   int                          `json:"participantIndex"`
	RoleConfigs        []MCMSRoleConfigParams       `json:"roleConfigs"`
}

var updateMCMSConfigSequence = operations.NewSequence(
	"canton/mcms/update_config",
	semver.MustParse("0.1.0"),
	"Updates signer configuration on an existing Canton MCMS contract",
	func(b operations.Bundle, deps canton.Chain, input updateMCMSConfigInput) (ccipsequences.OnChainOutput, error) {
		for i, roleConfig := range input.RoleConfigs {
			groupConfig, err := buildNormalizedConfig(roleConfig.Config)
			if err != nil {
				return ccipsequences.OnChainOutput{}, fmt.Errorf("build MCMS config for role %s: %w", roleConfig.Role, err)
			}

			_, err = operations.ExecuteOperation(b, mcmsops.SetConfig, deps, opcontract.ChoiceInput[mcmsCore.SetConfig]{
				InstanceAddress:  input.RawInstanceAddress.InstanceAddress(),
				ParticipantIndex: input.ParticipantIndex,
				Args: mcmsCore.SetConfig{
					TargetRole:      roleConfig.Role,
					NewSigners:      groupConfig.Signers,
					NewGroupQuorums: groupConfig.GroupQuorums,
					NewGroupParents: groupConfig.GroupParents,
					ClearRoot:       types.BOOL(roleConfig.Config.ClearRoot),
				},
			})
			if err != nil {
				return ccipsequences.OnChainOutput{}, fmt.Errorf("update MCMS role %s at index %d: %w", roleConfig.Role, i, err)
			}
		}

		return ccipsequences.OnChainOutput{}, nil
	},
)
