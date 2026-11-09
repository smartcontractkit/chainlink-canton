package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linktokenadmin"
)

// ExecuteLinkTokenAdminAuthorization executes an MCMS-approved LinkTokenAdmin authorization as the
// selected participant's party. Config.ParticipantIndex is ignored in favor of Participant.
type ExecuteLinkTokenAdminAuthorization struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[linktokenadmin.ExecuteInput]] = ExecuteLinkTokenAdminAuthorization{}

func (ExecuteLinkTokenAdminAuthorization) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[linktokenadmin.ExecuteInput]) error {
	if err := verifyParticipant(e, config.ChainSelector, config.Participant); err != nil {
		return err
	}

	return config.Config.Validate()
}

func (ExecuteLinkTokenAdminAuthorization) Apply(e cldf.Environment, config CantonCSDeps[linktokenadmin.ExecuteInput]) (cldf.ChangesetOutput, error) {
	input := config.Config
	input.ParticipantIndex = config.Participant
	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, linktokenadmin.Execute, e.BlockChains.CantonChains()[config.ChainSelector], input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("execute %s authorization %s: %w", input.Kind, input.Authorization, err)
	}

	return cldf.ChangesetOutput{Reports: []cld_ops.Report[any, any]{report.ToGenericReport()}}, nil
}

// CleanupExpiredLinkTokenAdminAuthorization archives an expired LinkTokenAdmin authorization as the
// selected participant's party.
type CleanupExpiredLinkTokenAdminAuthorization struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[linktokenadmin.CleanupInput]] = CleanupExpiredLinkTokenAdminAuthorization{}

func (CleanupExpiredLinkTokenAdminAuthorization) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[linktokenadmin.CleanupInput]) error {
	if err := verifyParticipant(e, config.ChainSelector, config.Participant); err != nil {
		return err
	}

	return config.Config.Validate()
}

func (CleanupExpiredLinkTokenAdminAuthorization) Apply(e cldf.Environment, config CantonCSDeps[linktokenadmin.CleanupInput]) (cldf.ChangesetOutput, error) {
	input := config.Config
	input.ParticipantIndex = config.Participant
	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, linktokenadmin.CleanupExpired, e.BlockChains.CantonChains()[config.ChainSelector], input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("clean up %s authorization %s: %w", input.Kind, input.Authorization, err)
	}

	return cldf.ChangesetOutput{Reports: []cld_ops.Report[any, any]{report.ToGenericReport()}}, nil
}

func verifyParticipant(e cldf.Environment, chainSelector uint64, participant int) error {
	chain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", chainSelector)
	}
	if participant < 0 || participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", participant, chainSelector, len(chain.Participants))
	}

	return nil
}
