package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipantwithacs"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/archivecontracts"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/example"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// resumeCmd resumes an existing ceremony from its persisted state.
//
// Usage:
//
//	canton-party-ceremony resume <ceremony-id> \
//	  --config ./participant-config.json \
//	  --state-dir ./ceremonies
//
// On each invocation resume:
//  1. Loads workflow.json to obtain the OnboardingInput (ceremony parameters).
//  2. Loads reports.json into a MemoryReporter so all previously completed
//     operations are served from cache.
//  3. Calls OnboardingSequence — only the operations that are still pending
//     for the current participant are actually executed.
//  4. Saves the updated reports back to reports.json.
//
// If the signature threshold is not yet met, the command exits with code 2 and
// prints a human-readable message directing operators to sign and re-run.
var resumeCmd = &cobra.Command{
	Use:   "resume <ceremony-id>",
	Short: "Resume an existing ceremony",
	Long: `Resume an existing ceremony from its persisted state.

The ceremony directory (<state-dir>/<ceremony-id>/) must have been created by
a previous "init" call.  resume loads workflow.json and reports.json, rebuilds
the Operations bundle with the cached reports, and re-runs OnboardingSequence.

Operations that were already completed successfully are skipped instantly (they
are served from the cached reports).  Only pending operations for the current
participant are executed.

If the signature threshold is not yet met the command exits with code 2 and
prints a message asking additional participants to sign.`,
	Args: cobra.ExactArgs(1),
	RunE: runResume,
}

func init() {
	f := resumeCmd.Flags()

	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	rootCmd.AddCommand(resumeCmd)
}

func runResume(cmd *cobra.Command, args []string) error {
	workflowId := args[0]
	if workflowId == "" {
		return fmt.Errorf("ceremony ID is required")
	}

	configPath, _ := cmd.Flags().GetString("config")
	stateDir, _ := cmd.Flags().GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)

	workflowType, err := ceremony.PeekWorkflowType(ceremonyDir)
	if err != nil {
		return fmt.Errorf("ceremony %q: reading workflow type: %w", workflowId, err)
	}

	switch workflowType {
	case ceremony.WorkflowTypeOnboarding:
		state, err := ceremony.LoadWorkflow[onboarding.OnboardingInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeOnboardingSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeExample:
		state, err := ceremony.LoadWorkflow[example.OnboardingInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeExampleOnboardingSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeKick:
		state, err := ceremony.LoadWorkflow[kick.KickInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeKickSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeContractDeploy:
		state, err := ceremony.LoadWorkflow[contractdeploy.ContractDeployInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeContractDeploySequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeAddParticipant:
		state, err := ceremony.LoadWorkflow[addparticipant.AddParticipantInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeAddParticipantSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeKeyRotation:
		state, err := ceremony.LoadWorkflow[keyrotation.KeyRotationInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeKeyRotationSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeAddParticipantWithAcs:
		state, err := ceremony.LoadWorkflow[addparticipantwithacs.AddParticipantWithAcsInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeAddParticipantWithAcsSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	case ceremony.WorkflowTypeArchiveContracts:
		state, err := ceremony.LoadWorkflow[archivecontracts.ArchiveContractsInput](ceremonyDir)
		if err != nil {
			return fmt.Errorf("ceremony %q: loading workflow: %w", workflowId, err)
		}

		return executeArchiveContractsSequence(cmd.Context(), cfg, state.Input, stateDir, workflowId, confirmerFromFlags(cmd))

	default:
		return fmt.Errorf("ceremony %q: unknown workflow type %q; supported types: %s, %s, %s, %s, %s, %s, %s, %s",
			workflowId, workflowType,
			ceremony.WorkflowTypeOnboarding, ceremony.WorkflowTypeExample, ceremony.WorkflowTypeKick, ceremony.WorkflowTypeContractDeploy, ceremony.WorkflowTypeAddParticipant, ceremony.WorkflowTypeKeyRotation, ceremony.WorkflowTypeAddParticipantWithAcs, ceremony.WorkflowTypeArchiveContracts)
	}
}
