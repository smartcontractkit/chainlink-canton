package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipantwithacs"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/example"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <ceremony-id>",
	Short: "Inspect the current state of a ceremony",
	Long: `Read the persisted ceremony state and display progress without executing
any operations.

Prints the current phase, which members have initialised, signature collection
progress, and the final result when the ceremony is complete.

Use --json to emit machine-readable JSON suitable for scripting or dashboards.`,
	Args: cobra.ExactArgs(1),
	RunE: runInspect,
}

func init() {
	f := inspectCmd.Flags()
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")
	f.Bool("json", false, "Emit state as JSON instead of human-readable text")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	workflowID := args[0]
	stateDir, _ := cmd.Flags().GetString("state-dir")
	asJSON, _ := cmd.Flags().GetBool("json")

	ceremonyDir := calculateCeremonyDir(stateDir, workflowID)

	workflowType, err := ceremony.PeekWorkflowType(ceremonyDir)
	if err != nil {
		return fmt.Errorf("ceremony %q: %w", workflowID, err)
	}

	switch workflowType {
	case ceremony.WorkflowTypeExample:
		return runInspectExample(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeOnboarding:
		return runInspectOnboarding(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeKick:
		return runInspectKick(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeAddParticipant:
		return runInspectAddParticipant(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeContractDeploy:
		return runInspectContractDeploy(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeKeyRotation:
		return runInspectKeyRotation(ceremonyDir, workflowID, workflowType, asJSON)
	case ceremony.WorkflowTypeAddParticipantWithAcs:
		return runInspectAddParticipantWithAcs(ceremonyDir, workflowID, workflowType, asJSON)
	default:
		return fmt.Errorf("ceremony %q: state inspection is not yet supported for workflow type %q", workflowID, workflowType)
	}
}

func runInspectExample(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[example.OnboardingInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := example.QueryState(reporter, wf.Input)

	if asJSON {
		payload := struct {
			CeremonyID   string                 `json:"ceremony_id"`
			WorkflowType string                 `json:"workflow_type"`
			Found        bool                   `json:"found"`
			Phase        example.Phase          `json:"phase,omitempty"`
			State        *example.CeremonyState `json:"state,omitempty"`
			PartyID      string                 `json:"party_id,omitempty"`
			DNSConfirmed *bool                  `json:"dns_confirmed,omitempty"`
			P2PConfirmed *bool                  `json:"p2p_confirmed,omitempty"`
		}{
			CeremonyID:   workflowID,
			WorkflowType: workflowType,
			Found:        found,
		}
		if found {
			payload.Phase = out.State.Phase
			payload.State = &out.State
			if out.PartyID != "" {
				payload.PartyID = out.PartyID
				payload.DNSConfirmed = &out.DNSConfirmed
				payload.P2PConfirmed = &out.P2PConfirmed
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(payload)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", inspectPhaseLabel(s.Phase))
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Members initialized:\t%d / %d  (%s)\n",
		len(s.InitializedMembers), len(wf.Input.Participants),
		strings.Join(s.InitializedMembers, ", "))

	if s.ProposalHash != "" {
		short := s.ProposalHash
		if len(short) > 16 {
			short = short[:16] + "…"
		}
		fmt.Fprintf(w, "Proposal hash:\t%s\n", short)
	}

	if len(s.RequiredSigners) > 0 {
		fmt.Fprintf(w, "Signatures:\t%d / %d  (threshold %d)\n",
			len(s.CollectedSigners), len(s.RequiredSigners), s.Threshold)
		fmt.Fprintln(w)

		collected := make(map[string]bool, len(s.CollectedSigners))
		for _, id := range s.CollectedSigners {
			collected[id] = true
		}
		for _, signer := range s.RequiredSigners {
			if collected[signer] {
				fmt.Fprintf(w, "  [signed]\t%s\n", signer)
			} else {
				fmt.Fprintf(w, "  [pending]\t%s\n", signer)
			}
		}
	}

	if out.PartyID != "" {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Party ID:\t%s\n", out.PartyID)
		fmt.Fprintf(w, "DNS confirmed:\t%v\n", out.DNSConfirmed)
		fmt.Fprintf(w, "P2P confirmed:\t%v\n", out.P2PConfirmed)
	}

	return nil
}

// Common phases
const completed = "completed"
const dnsProposal = "dns-proposal — creating DNS proposal"
const dnsSigning = "dns-signing — collecting DNS signatures"
const dnsSubmit = "dns-submit — submitting DNS update"
const p2pProposal = "p2p — collecting P2P proposals"
const readState = "read-state — reading current topology"

func inspectPhaseLabel(p example.Phase) string {
	switch p {
	case example.PhaseInit:
		return "init — waiting for all members to initialise"
	case example.PhaseProposal:
		return "proposal — waiting for coordinator to create proposal"
	case example.PhaseSigning:
		return "signing — collecting signatures"
	case example.PhaseSubmit:
		return "submit — threshold met, waiting for coordinator to submit"
	case example.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Onboarding ───────────────────────────────────────────────────────────────

func runInspectOnboarding(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[onboarding.OnboardingInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := onboarding.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", onboardingPhaseLabel(s.Phase))
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Keys generated:\t%d / %d  (%s)\n",
		len(s.KeysGenerated), s.Threshold, strings.Join(s.KeysGenerated, ", "))

	printDNSProgress(w, s.ProposalHash, s.RequiredSigners, s.CollectedSigners, s.Threshold)

	if out.PartyID != "" {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Party ID:\t%s\n", out.PartyID)
		fmt.Fprintf(w, "DNS confirmed:\t%v\n", out.DNSConfirmed)
		fmt.Fprintf(w, "P2P confirmed:\t%v\n", out.P2PConfirmed)
	}

	return nil
}

func onboardingPhaseLabel(p onboarding.Phase) string {
	switch p {
	case onboarding.PhaseKeyGen:
		return "key-gen — waiting for all members to generate keys"
	case onboarding.PhaseNSD:
		return "nsd — waiting for NSD proposals"
	case onboarding.PhaseDNSProposal:
		return dnsProposal
	case onboarding.PhaseDNSSigning:
		return dnsSigning
	case onboarding.PhaseDNSSubmit:
		return dnsSubmit
	case onboarding.PhaseP2P:
		return p2pProposal
	case onboarding.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Kick ─────────────────────────────────────────────────────────────────────

func runInspectKick(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[kick.KickInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := kick.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", kickPhaseLabel(s.Phase))
	fmt.Fprintf(w, "Kicked participant:\t%s\n", s.KickedParticipant)
	fmt.Fprintln(w)

	printDNSProgress(w, s.ProposalHash, s.RequiredSigners, s.CollectedSigners, s.DNSThreshold)

	if out.DNSUpdated {
		fmt.Fprintf(w, "DNS updated:\t%v\n", out.DNSUpdated)
	}
	if out.P2PUpdated {
		fmt.Fprintf(w, "P2P updated:\t%v\n", out.P2PUpdated)
	}

	return nil
}

func kickPhaseLabel(p kick.Phase) string {
	switch p {
	case kick.PhaseReadState:
		return readState
	case kick.PhaseDNSProposal:
		return dnsProposal
	case kick.PhaseDNSSigning:
		return dnsSigning
	case kick.PhaseDNSSubmit:
		return dnsSubmit
	case kick.PhaseP2P:
		return p2pProposal
	case kick.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Add Participant ──────────────────────────────────────────────────────────

func runInspectAddParticipant(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[addparticipant.AddParticipantInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := addparticipant.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", addParticipantPhaseLabel(s.Phase))
	fmt.Fprintf(w, "New member key ready:\t%v\n", s.NewMemberKeyReady)
	fmt.Fprintf(w, "NSD proposed:\t%v\n", s.NSDProposed)
	fmt.Fprintln(w)

	printDNSProgress(w, s.ProposalHash, s.RequiredSigners, s.CollectedSigners, s.DNSThreshold)

	if out.DNSUpdated {
		fmt.Fprintf(w, "DNS updated:\t%v\n", out.DNSUpdated)
	}
	if out.P2PUpdated {
		fmt.Fprintf(w, "P2P updated:\t%v\n", out.P2PUpdated)
	}

	return nil
}

func addParticipantPhaseLabel(p addparticipant.Phase) string {
	switch p {
	case addparticipant.PhaseKeyGen:
		return "key-gen — waiting for new member to generate keys"
	case addparticipant.PhaseNSD:
		return "nsd — waiting for NSD proposal"
	case addparticipant.PhaseReadState:
		return readState
	case addparticipant.PhaseDNSProposal:
		return dnsProposal
	case addparticipant.PhaseDNSSigning:
		return dnsSigning
	case addparticipant.PhaseDNSSubmit:
		return dnsSubmit
	case addparticipant.PhaseP2P:
		return p2pProposal
	case addparticipant.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Contract Deploy ──────────────────────────────────────────────────────────

func runInspectContractDeploy(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[contractdeploy.ContractDeployInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := contractdeploy.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", contractDeployPhaseLabel(s.Phase))
	fmt.Fprintf(w, "DARs uploaded:\t%d / %d  (%s)\n", len(s.DARsUploaded), s.DARsRequired,
		strings.Join(s.DARsUploaded, ", "))
	fmt.Fprintf(w, "Signed:\t%d / %d  (%s)\n",
		len(s.Signed), s.SignRequired, strings.Join(s.Signed, ", "))
	if s.PreparedTxHash != "" {
		short := s.PreparedTxHash
		if len(short) > 16 {
			short = short[:16] + "…"
		}
		fmt.Fprintf(w, "Prepared TX hash:\t%s\n", short)
	}

	if out.ContractID != "" {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Contract ID:\t%s\n", out.ContractID)
	}

	return nil
}

func contractDeployPhaseLabel(p contractdeploy.Phase) string {
	switch p {
	case contractdeploy.PhaseVerifyParty:
		return "verify-party — verifying party exists"
	case contractdeploy.PhaseFetchMembers:
		return "fetch-members — fetching member list"
	case contractdeploy.PhaseDARUpload:
		return "dar-upload — uploading DARs"
	case contractdeploy.PhasePrepare:
		return "prepare — preparing transaction"
	case contractdeploy.PhaseSigning:
		return "signing — collecting transaction signatures"
	case contractdeploy.PhaseExecute:
		return "execute — executing transaction"
	case contractdeploy.PhaseVerifyContract:
		return "verify-contract — verifying deployed contract"
	case contractdeploy.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Key Rotation ─────────────────────────────────────────────────────────────

func runInspectKeyRotation(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[keyrotation.KeyRotationInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := keyrotation.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", keyRotationPhaseLabel(s.Phase))
	fmt.Fprintf(w, "Rotate namespace:\t%v\n", s.RotateNamespace)
	fmt.Fprintf(w, "Rotate DAML:\t%v\n", s.RotateDaml)
	fmt.Fprintf(w, "Target key gen ready:\t%v\n", s.TargetKeyGenReady)
	fmt.Fprintf(w, "NSD proposed:\t%v\n", s.NSDProposed)
	fmt.Fprintln(w)

	printDNSProgress(w, s.ProposalHash, s.RequiredSigners, s.CollectedSigners, s.DNSThreshold)

	if s.P2PRequired > 0 {
		fmt.Fprintf(w, "P2P proposals:\t%d / %d\n", s.P2PProposedCount, s.P2PRequired)
	}

	if out.DNSUpdated {
		fmt.Fprintf(w, "DNS updated:\t%v\n", out.DNSUpdated)
	}
	if out.P2PUpdated {
		fmt.Fprintf(w, "P2P updated:\t%v\n", out.P2PUpdated)
	}

	return nil
}

func keyRotationPhaseLabel(p keyrotation.Phase) string {
	switch p {
	case keyrotation.PhaseReadState:
		return readState
	case keyrotation.PhaseKeyGen:
		return "key-gen — waiting for target to generate rotated keys"
	case keyrotation.PhaseNSD:
		return "nsd — waiting for rotated NSD proposal"
	case keyrotation.PhaseDNSProposal:
		return "dns-proposal — creating rotation DNS proposal"
	case keyrotation.PhaseDNSSigning:
		return dnsSigning
	case keyrotation.PhaseDNSSubmit:
		return "dns-submit — submitting rotation DNS update"
	case keyrotation.PhaseP2P:
		return "p2p — collecting P2P rotation proposals"
	case keyrotation.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}

// ── Shared helpers ───────────────────────────────────────────────────────────

// emitJSON is a generic JSON emitter for inspect commands.
func emitJSON(workflowID, workflowType string, found bool, phase string, state any) error {
	payload := struct {
		CeremonyID   string `json:"ceremony_id"`
		WorkflowType string `json:"workflow_type"`
		Found        bool   `json:"found"`
		Phase        string `json:"phase,omitempty"`
		State        any    `json:"state,omitempty"`
	}{
		CeremonyID:   workflowID,
		WorkflowType: workflowType,
		Found:        found,
	}
	if found {
		payload.Phase = phase
		payload.State = state
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(payload)
}

// printDNSProgress prints common DNS signature progress to a tabwriter.
func printDNSProgress(w *tabwriter.Writer, proposalHash string, requiredSigners, collectedSigners []string, threshold int) {
	if proposalHash != "" {
		short := proposalHash
		if len(short) > 16 {
			short = short[:16] + "…"
		}
		fmt.Fprintf(w, "Proposal hash:\t%s\n", short)
	}

	if len(requiredSigners) > 0 {
		fmt.Fprintf(w, "Signatures:\t%d / %d  (threshold %d)\n",
			len(collectedSigners), len(requiredSigners), threshold)
		fmt.Fprintln(w)

		collected := make(map[string]bool, len(collectedSigners))
		for _, id := range collectedSigners {
			collected[id] = true
		}
		for _, signer := range requiredSigners {
			if collected[signer] {
				fmt.Fprintf(w, "  [signed]\t%s\n", signer)
			} else {
				fmt.Fprintf(w, "  [pending]\t%s\n", signer)
			}
		}
	}
}

// ── Add Participant With ACS ─────────────────────────────────────────────────

func runInspectAddParticipantWithAcs(ceremonyDir, workflowID, workflowType string, asJSON bool) error {
	wf, err := ceremony.LoadWorkflow[addparticipantwithacs.AddParticipantWithAcsInput](ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading workflow: %w", err)
	}

	reports, err := ceremony.LoadReports(ceremonyDir)
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}

	reporter := operations.NewMemoryReporter(operations.WithReports(reports))
	out, found := addparticipantwithacs.QueryState(reporter, wf.Input)

	if asJSON {
		return emitJSON(workflowID, workflowType, found, string(out.State.Phase), &out.State)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Ceremony ID:\t%s\n", workflowID)
	fmt.Fprintf(w, "Type:\t%s\n", workflowType)

	if !found {
		fmt.Fprintf(w, "Phase:\tinit  (no runs recorded yet)\n")
		return nil
	}

	s := out.State
	fmt.Fprintf(w, "Phase:\t%s\n", addParticipantWithAcsPhaseLabel(s.Phase))
	fmt.Fprintf(w, "New member key ready:\t%v\n", s.NewMemberKeyReady)
	fmt.Fprintf(w, "NSD proposed:\t%v\n", s.NSDProposed)
	fmt.Fprintf(w, "Ledger offset recorded:\t%v\n", s.LedgerOffsetRecorded)
	fmt.Fprintf(w, "ACS exported:\t%v\n", s.AcsExported)
	fmt.Fprintf(w, "ACS imported:\t%v\n", s.AcsImported)
	fmt.Fprintf(w, "Onboarding flag cleared:\t%v\n", s.OnboardingFlagCleared)
	fmt.Fprintln(w)

	printDNSProgress(w, s.ProposalHash, s.RequiredSigners, s.CollectedSigners, s.DNSThreshold)

	if out.DNSUpdated {
		fmt.Fprintf(w, "DNS updated:\t%v\n", out.DNSUpdated)
	}
	if out.P2PUpdated {
		fmt.Fprintf(w, "P2P updated:\t%v\n", out.P2PUpdated)
	}

	return nil
}

func addParticipantWithAcsPhaseLabel(p addparticipantwithacs.Phase) string {
	switch p {
	case addparticipantwithacs.PhaseKeyGen:
		return "key-gen — waiting for new member to generate keys"
	case addparticipantwithacs.PhaseRecordTargetOffset:
		return "record-target-offset — recording ledger offset on target"
	case addparticipantwithacs.PhaseNSD:
		return "nsd — waiting for NSD proposal"
	case addparticipantwithacs.PhaseReadState:
		return readState
	case addparticipantwithacs.PhaseDNSProposal:
		return dnsProposal
	case addparticipantwithacs.PhaseDNSSigning:
		return dnsSigning
	case addparticipantwithacs.PhaseDNSSubmit:
		return dnsSubmit
	case addparticipantwithacs.PhaseRecordOffset:
		return "record-offset — recording ledger offset on source"
	case addparticipantwithacs.PhaseP2POnboarding:
		return "p2p-onboarding — collecting P2P proposals with onboarding flag"
	case addparticipantwithacs.PhaseTargetDisconnect:
		return "target-disconnect — target disconnects from synchronizer"
	case addparticipantwithacs.PhaseAcsExport:
		return "acs-export — exporting ACS from source participant"
	case addparticipantwithacs.PhaseAcsImport:
		return "acs-import — importing ACS into target participant"
	case addparticipantwithacs.PhaseTargetReconnect:
		return "target-reconnect — target reconnects to synchronizer"
	case addparticipantwithacs.PhaseClearOnboarding:
		return "clear-onboarding — clearing onboarding flag on target"
	case addparticipantwithacs.PhaseCompleted:
		return completed
	default:
		return string(p)
	}
}
