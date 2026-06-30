package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipantwithacs"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/archivecontracts"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/example"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// initCmd is the parent command for all ceremony-type initialisers.
// Additional workflow types (e.g. kick) would be added as further sub-commands.
var initCmd = &cobra.Command{
	Use:   "init <workflow-type>",
	Short: "Initialise a new ceremony",
	Long:  "Create a new ceremony directory and run the first sequence step.",
}

// initOnboardingCmd initialises a decentralized-party onboarding ceremony.
//
// Usage:
//
//	canton-party-ceremony init onboarding \
//	  --namespace-name cbtc-onboard-2026 \
//	  --coordinator p1 \
//	  --participants p1,p2,p3 \
//	  --party-id-prefix cbtc-network \
//	  --synchronizer-id global \
//	  --threshold 2 \
//	  --config ./participant-config.json
var initOnboardingCmd = &cobra.Command{
	Use:   "onboarding",
	Short: "Initialise a new party-onboarding ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step.  The onboarding ceremony follows the Canton Party Management
Tooling Spec: member-init → create-proposal → sign → submit.`,
	RunE: runInitOnboarding,
}

func init() {
	f := initOnboardingCmd.Flags()

	f.String("new-namespace-name", "", "Unique ceremony identifier (required)")
	f.String("coordinator", "", "Coordinator participant ID (required)")
	f.String("participants", "", "Comma-separated list of participant IDs (required)")
	f.String("new-party-name", "", "Party ID prefix used to derive the final party identifier (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.Int("threshold", 0, "Minimum signatures required before submission. 0 = strict majority (floor(n/2)+1)")
	f.String("kms-vault-name", "", "Optional KMS vault registration name (defaults to --new-namespace-name). Use to reuse keys registered under another ceremony name.")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initOnboardingCmd.MarkFlagRequired("new-namespace-name")
	_ = initOnboardingCmd.MarkFlagRequired("participants")
	_ = initOnboardingCmd.MarkFlagRequired("new-party-name")
	_ = initOnboardingCmd.MarkFlagRequired("synchronizer-id")

	initCmd.AddCommand(initOnboardingCmd)
	rootCmd.AddCommand(initCmd)
}

func runInitOnboarding(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	namespaceName, _ := f.GetString("new-namespace-name")
	participantsRaw, _ := f.GetString("participants")
	partyName, _ := f.GetString("new-party-name")
	synchronizerID, _ := f.GetString("synchronizer-id")
	threshold, _ := f.GetInt("threshold")
	kmsVaultName, _ := f.GetString("kms-vault-name")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Normalise participant list: trim spaces, drop empty entries.
	participants := splitParticipants(participantsRaw)
	if len(participants) == 0 {
		return fmt.Errorf("--participants must contain at least one participant ID")
	}

	input := onboarding.OnboardingInput{
		NamespaceName:  namespaceName,
		KmsVaultName:   kmsVaultName,
		PartyPrefix:    partyName,
		Participants:   participants,
		SynchronizerID: synchronizerID,
		Threshold:      threshold,
	}

	return executeOnboardingSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

var initExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Initialise a mocked new party-onboarding ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step.  The onboarding ceremony follows the Canton Party Management
Tooling Spec: member-init → create-proposal → sign → submit.`,
	RunE: runInitExample,
}

func init() {
	f := initExampleCmd.Flags()

	f.String("new-namespace-name", "", "Unique ceremony identifier (required)")
	f.String("coordinator", "", "Coordinator participant ID (required)")
	f.String("participants", "", "Comma-separated list of participant IDs (required)")
	f.String("new-party-name", "", "Party ID prefix used to derive the final party identifier (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.Int("threshold", 0, "Minimum signatures required before submission. 0 = strict majority (floor(n/2)+1)")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initExampleCmd.MarkFlagRequired("new-namespace-name")
	_ = initExampleCmd.MarkFlagRequired("participants")
	_ = initExampleCmd.MarkFlagRequired("new-party-name")
	_ = initExampleCmd.MarkFlagRequired("synchronizer-id")

	initCmd.AddCommand(initExampleCmd)
	rootCmd.AddCommand(initCmd)
}

func runInitExample(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	namespaceName, _ := f.GetString("new-namespace-name")
	participantsRaw, _ := f.GetString("participants")
	partyName, _ := f.GetString("new-party-name")
	synchronizerID, _ := f.GetString("synchronizer-id")
	threshold, _ := f.GetInt("threshold")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Normalise participant list: trim spaces, drop empty entries.
	participants := splitParticipants(participantsRaw)
	if len(participants) == 0 {
		return fmt.Errorf("--participants must contain at least one participant ID")
	}

	input := example.OnboardingInput{
		NamespaceName:  namespaceName,
		PartyName:      partyName,
		Participants:   participants,
		SynchronizerID: synchronizerID,
		Threshold:      threshold,
	}

	return executeExampleOnboardingSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// initKickCmd initialises a decentralized-party kick ceremony.
//
// Usage:
//
//	canton-party-ceremony init kick \
//	  --decentralized-party-id "prefix::namespace" \
//	  --kicked-participant-id  "PAR::name::fp" \
//	  --kicked-namespace-fingerprint "1220abcd..." \
//	  --remaining-participants  "PAR::a::fp1,PAR::b::fp2" \
//	  --synchronizer-id global \
//	  --config ./participant-config.json
var initKickCmd = &cobra.Command{
	Use:   "kick",
	Short: "Initialise a new party-kick ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step. The kick ceremony removes a participant from an existing
decentralized party by updating both the DecentralizedNamespaceDefinition
and the PartyToParticipant topology mappings.`,
	RunE: runInitKick,
}

func init() {
	f := initKickCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID in the format <prefix>::<namespace> (required)")
	f.String("kicked-participant-id", "", "Canton UID of the participant to remove (required)")
	f.String("kicked-namespace-fingerprint", "", "Namespace fingerprint of the participant being kicked (required)")
	f.String("remaining-participants", "", "Comma-separated Canton UIDs of participants that will remain (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.Int("new-threshold", 0, "Signing threshold after the kick. 0 = strict majority (floor(n/2)+1)")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initKickCmd.MarkFlagRequired("decentralized-party-id")
	_ = initKickCmd.MarkFlagRequired("kicked-participant-id")
	_ = initKickCmd.MarkFlagRequired("kicked-namespace-fingerprint")
	_ = initKickCmd.MarkFlagRequired("remaining-participants")
	_ = initKickCmd.MarkFlagRequired("synchronizer-id")

	initCmd.AddCommand(initKickCmd)
}

func runInitKick(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	kickedUID, _ := f.GetString("kicked-participant-id")
	kickedFP, _ := f.GetString("kicked-namespace-fingerprint")
	remainingRaw, _ := f.GetString("remaining-participants")
	synchronizerID, _ := f.GetString("synchronizer-id")
	newThreshold, _ := f.GetInt("new-threshold")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	remaining := splitParticipants(remainingRaw)
	if len(remaining) < 2 {
		return fmt.Errorf("--remaining-participants must list at least 2 participants")
	}

	if slices.Contains(remaining, kickedUID) {
		return fmt.Errorf("kicked participant %q must not appear in --remaining-participants", kickedUID)
	}

	input := kick.KickInput{
		DecentralizedPartyID:       partyID,
		KickedParticipantID:        kickedUID,
		KickedNamespaceFingerprint: kickedFP,
		RemainingParticipants:      remaining,
		SynchronizerID:             synchronizerID,
		NewThreshold:               newThreshold,
	}

	return executeKickSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// initAddParticipantCmd initialises an add-participant ceremony.
//
// Usage:
//
//	canton-party-ceremony init add-participant \
//	  --decentralized-party-id "prefix::namespace" \
//	  --new-participant-id "PAR::newnode::fp" \
//	  --namespace-name "add-2026" \
//	  --synchronizer-id global \
//	  --config ./participant-config.json
var initAddParticipantCmd = &cobra.Command{
	Use:   "add-participant",
	Short: "Initialise a new add-participant ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step. The add-participant ceremony adds a new participant to an
existing decentralized party by updating both the DecentralizedNamespaceDefinition
and the PartyToParticipant topology mappings.`,
	RunE: runInitAddParticipant,
}

func init() {
	f := initAddParticipantCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID in the format <prefix>::<namespace> (required)")
	f.String("new-participant-id", "", "Canton UID of the participant to add (required)")
	f.String("namespace-name", "", "Label for the new participant's key generation (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.Int("new-threshold", 0, "Signing threshold after the addition. 0 = keep current")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initAddParticipantCmd.MarkFlagRequired("decentralized-party-id")
	_ = initAddParticipantCmd.MarkFlagRequired("new-participant-id")
	_ = initAddParticipantCmd.MarkFlagRequired("namespace-name")
	_ = initAddParticipantCmd.MarkFlagRequired("synchronizer-id")

	initCmd.AddCommand(initAddParticipantCmd)
}

func runInitAddParticipant(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	newUID, _ := f.GetString("new-participant-id")
	namespaceName, _ := f.GetString("namespace-name")
	synchronizerID, _ := f.GetString("synchronizer-id")
	newThreshold, _ := f.GetInt("new-threshold")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	input := addparticipant.AddParticipantInput{
		DecentralizedPartyID: partyID,
		NewParticipantID:     newUID,
		NamespaceName:        namespaceName,
		SynchronizerID:       synchronizerID,
		NewThreshold:         newThreshold,
	}

	return executeAddParticipantSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// initKeyRotationCmd initialises a key rotation ceremony.
//
// Usage:
//
//	canton-party-ceremony init key-rotation \
//	  --decentralized-party-id "prefix::namespace" \
//	  --target-participant-id "PAR::name::fp" \
//	  --target-namespace-fingerprint "1220abcd..." \
//	  --synchronizer-id global \
//	  --rotate-namespace-key \
//	  --rotate-daml-key \
//	  --config ./participant-config.json
var initKeyRotationCmd = &cobra.Command{
	Use:   "key-rotation",
	Short: "Initialise a new key rotation ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step. The key-rotation ceremony replaces a participant's namespace
key and/or DAML signing key in the topology. Participants are dynamically
fetched from the current P2P topology state.`,
	RunE: runInitKeyRotation,
}

func init() {
	f := initKeyRotationCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID in the format <prefix>::<namespace> (required)")
	f.String("target-participant-id", "", "Canton UID of the participant whose key is being rotated (required)")
	f.String("target-namespace-fingerprint", "", "Current namespace fingerprint of the target (required when --rotate-namespace-key)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.Bool("rotate-namespace-key", false, "Rotate the namespace signing key")
	f.Bool("rotate-daml-key", false, "Rotate the DAML (protocol) signing key")
	f.Int("new-threshold", 0, "Signing threshold after rotation. 0 = keep current")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initKeyRotationCmd.MarkFlagRequired("decentralized-party-id")
	_ = initKeyRotationCmd.MarkFlagRequired("target-participant-id")
	_ = initKeyRotationCmd.MarkFlagRequired("synchronizer-id")

	initCmd.AddCommand(initKeyRotationCmd)
}

func runInitKeyRotation(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	targetUID, _ := f.GetString("target-participant-id")
	targetNSFP, _ := f.GetString("target-namespace-fingerprint")
	synchronizerID, _ := f.GetString("synchronizer-id")
	rotateNS, _ := f.GetBool("rotate-namespace-key")
	rotateDAML, _ := f.GetBool("rotate-daml-key")
	newThreshold, _ := f.GetInt("new-threshold")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	if !rotateNS && !rotateDAML {
		return fmt.Errorf("at least one of --rotate-namespace-key or --rotate-daml-key must be set")
	}

	if rotateNS && targetNSFP == "" {
		return fmt.Errorf("--target-namespace-fingerprint is required when --rotate-namespace-key is set")
	}

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	input := keyrotation.KeyRotationInput{
		DecentralizedPartyID:       partyID,
		TargetParticipantID:        targetUID,
		TargetNamespaceFingerprint: targetNSFP,
		SynchronizerID:             synchronizerID,
		RotateNamespaceKey:         rotateNS,
		RotateDamlKey:              rotateDAML,
		NewThreshold:               newThreshold,
	}

	return executeKeyRotationSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// splitParticipants splits a comma-separated participant string, trimming
// whitespace and dropping empty entries.
func splitParticipants(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}

	return out
}

// contractProfile holds the known defaults for a supported contract package.
type contractProfile struct {
	module   string
	entity   string
	argsFile string
}

// knownContracts maps package names (as passed to --packages) to their
// built-in template defaults. Add new entries here when new contract types are
// supported.
var knownContracts = map[string]contractProfile{
	"mcms-core": {
		module:   "MCMS.Main",
		entity:   "MCMS",
		argsFile: "mcms-args.json",
	},
}

// initContractDeployCmd initialises a contract deployment ceremony.
//
// Usage:
//
//	canton-party-ceremony init contract-deploy \
//	  --decentralized-party-id "prefix::namespace" \
//	  --synchronizer-id global \
//	  --packages "mcms-core:1.0.0" \
//	  --contract-args-file ./mcms-args.json \
//	  --config ./participant-config.json
var initContractDeployCmd = &cobra.Command{
	Use:   "contract-deploy",
	Short: "Initialise a new contract deployment ceremony",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step. The contract-deploy ceremony uploads DARs to all participants,
verifies the decentralized party, and prepares a contract creation transaction
via InteractiveSubmissionService. Signing and execution are not yet implemented.

For known contract packages (e.g. "mcms-core"), --template-module, --template-entity,
and --contract-args-file are auto-populated with built-in defaults. For any other
package, all three flags must be provided explicitly.`,
	RunE: runInitContractDeploy,
}

func init() {
	f := initContractDeployCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID in the format <prefix>::<namespace> (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.String("packages", "", "Comma-separated list of packages in name:version format, e.g. mcms-core:1.0.0,globalconfig:1.0.0 (required)")
	f.String("template-module", "", "Fully-qualified DAML module name, e.g. MCMS.Main")
	f.String("template-entity", "", "DAML template entity name, e.g. MCMS")
	f.String("contract-args-file", "", "Path to JSON file containing contract creation arguments")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initContractDeployCmd.MarkFlagRequired("decentralized-party-id")
	_ = initContractDeployCmd.MarkFlagRequired("synchronizer-id")
	_ = initContractDeployCmd.MarkFlagRequired("packages")

	initCmd.AddCommand(initContractDeployCmd)
}

func runInitContractDeploy(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	synchronizerID, _ := f.GetString("synchronizer-id")
	packagesRaw, _ := f.GetString("packages")
	templateModule, _ := f.GetString("template-module")
	templateEntity, _ := f.GetString("template-entity")
	argsFile, _ := f.GetString("contract-args-file")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	packages, err := parsePackageRefs(packagesRaw)
	if err != nil {
		return err
	}
	if len(packages) == 0 {
		return fmt.Errorf("--packages must list at least one package")
	}

	templateModule, templateEntity, argsFile, err = applyContractDefaults(packages, templateModule, templateEntity, argsFile)
	if err != nil {
		return err
	}

	var contractArgs string
	if argsFile != "" {
		data, readErr := readFileBytes(argsFile)
		if readErr != nil {
			return fmt.Errorf("reading contract args file %q: %w", argsFile, readErr)
		}
		contractArgs = string(data)
	}

	input := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: partyID,
		SynchronizerID:       synchronizerID,
		Packages:             packages,
		TemplateModule:       templateModule,
		TemplateEntity:       templateEntity,
		ContractArgs:         contractArgs,
	}

	return executeContractDeploySequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// readFileBytes reads a file and returns its contents as a byte slice.
func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// applyContractDefaults fills in templateModule, entity, and argsFile from the
// built-in profile for the first package when its name matches a known contract.
// For unknown package names, it returns an error when module or entity are empty.
func applyContractDefaults(
	pkgs []contractdeploy.PackageRef,
	module, entity, argsFile string,
) (string, string, string, error) {
	if len(pkgs) == 0 {
		return module, entity, argsFile, nil
	}

	primary := pkgs[0].Name
	profile, known := knownContracts[primary]
	if known {
		if module == "" {
			module = profile.module
		}
		if entity == "" {
			entity = profile.entity
		}
		if argsFile == "" {
			argsFile = profile.argsFile
		}

		return module, entity, argsFile, nil
	}

	// Unknown contract: require both module and entity from the caller.
	if module == "" || entity == "" {
		return "", "", "", fmt.Errorf(
			"package %q is not a known contract type; --template-module and --template-entity are required",
			primary,
		)
	}

	return module, entity, argsFile, nil
}

// parsePackageRefs parses a comma-separated list of "name:version" entries
// into a slice of [contractdeploy.PackageRef]s.
// Example input: "mcms-core:1.0.0,globalconfig:1.0.0"
func parsePackageRefs(raw string) ([]contractdeploy.PackageRef, error) {
	parts := splitParticipants(raw) // reuses comma splitter
	refs := make([]contractdeploy.PackageRef, 0, len(parts))
	for _, p := range parts {
		idx := strings.LastIndex(p, ":")
		if idx <= 0 || idx == len(p)-1 {
			return nil, fmt.Errorf("invalid package reference %q: expected \"name:version\" format", p)
		}
		refs = append(refs, contractdeploy.PackageRef{
			Name:    p[:idx],
			Version: p[idx+1:],
		})
	}

	return refs, nil
}

// initAddParticipantWithAcsCmd initialises a combined add-participant + ACS
// replication ceremony.
//
// Usage:
//
//	canton-party-ceremony init add-participant-with-acs \
//	  --decentralized-party-id "prefix::namespace" \
//	  --new-participant-id "PAR::newnode::fp" \
//	  --namespace-name "add-2026" \
//	  --source-participant-id "PAR::source::fp" \
//	  --synchronizer-id global \
//	  --synchronizer-alias global \
//	  --config ./participant-config.json
var initAddParticipantWithAcsCmd = &cobra.Command{
	Use:   "add-participant-with-acs",
	Short: "Initialise a new add-participant ceremony with ACS replication",
	Long: `Create the ceremony directory, write workflow.json, and run the first
sequence step. This ceremony adds a new participant to an existing decentralized
party with the onboarding flag set, exports the Active Contract Set from a
source participant, imports it into the new participant, and clears the flag.`,
	RunE: runInitAddParticipantWithAcs,
}

func init() {
	f := initAddParticipantWithAcsCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID in the format <prefix>::<namespace> (required)")
	f.String("new-participant-id", "", "Canton UID of the participant to add (required)")
	f.String("namespace-name", "", "Label for the new participant's key generation (required)")
	f.String("source-participant-id", "", "Canton UID of the existing participant that exports ACS (required)")
	f.String("synchronizer-id", "", "Canton synchronizer ID (required)")
	f.String("synchronizer-alias", "", "Human-readable synchronizer alias for disconnect/reconnect (required)")
	f.Int("new-threshold", 0, "Signing threshold after the addition. 0 = keep current")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("decentralized-party-id")
	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("new-participant-id")
	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("namespace-name")
	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("source-participant-id")
	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("synchronizer-id")
	_ = initAddParticipantWithAcsCmd.MarkFlagRequired("synchronizer-alias")

	initCmd.AddCommand(initAddParticipantWithAcsCmd)
}

func runInitAddParticipantWithAcs(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	newUID, _ := f.GetString("new-participant-id")
	namespaceName, _ := f.GetString("namespace-name")
	sourceUID, _ := f.GetString("source-participant-id")
	synchronizerID, _ := f.GetString("synchronizer-id")
	synchronizerAlias, _ := f.GetString("synchronizer-alias")
	newThreshold, _ := f.GetInt("new-threshold")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	input := addparticipantwithacs.AddParticipantWithAcsInput{
		DecentralizedPartyID: partyID,
		NewParticipantID:     newUID,
		NamespaceName:        namespaceName,
		SynchronizerID:       synchronizerID,
		SynchronizerAlias:    synchronizerAlias,
		SourceParticipantID:  sourceUID,
		NewThreshold:         newThreshold,
	}

	return executeAddParticipantWithAcsSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}

// initArchiveContractsCmd initialises a multiparty contract archive ceremony.
//
// Usage:
//
//	canton-party-ceremony init archive-contracts \
//	  --decentralized-party-id "ccipOwner::1220..." \
//	  --synchronizer-id global \
//	  --template "#ccip-common:CCIP.GlobalConfig:GlobalConfig" \
//	  --config ./cv1.participant-config.json \
//	  --dry-run
var initArchiveContractsCmd = &cobra.Command{
	Use:   "archive-contracts",
	Short: "Initialise a multiparty contract archive ceremony",
	Long: `Archive active contracts owned by a decentralized party via InteractiveSubmission.

Run init on the coordinator participant (cv1) with a JWT that has CanActAs the
decentralized party. Each hosting participant must resume with the same
ceremony ID and a shared --state-dir (copy reports.json between nodes if needed).

Flow: discover ACS targets → prepare Archive batch → sign on every host → execute.

Use --dry-run to list matching contracts without submitting. Archives run one
contract per prepared transaction by default (Canton Interactive Submission
limit); increase --batch-size only if your participant supports multi-command prepare.`,
	RunE: runInitArchiveContracts,
}

func init() {
	f := initArchiveContractsCmd.Flags()

	f.String("decentralized-party-id", "", "Full party ID (required)")
	f.String("synchronizer-id", "global", "Canton synchronizer ID")
	f.StringArray("template", nil, "Template selector package:Module:Entity (#name for package name); repeat")
	f.String("inventory-file", "", "JSON file with explicit archive targets [{package_id,module_name,entity_name,contract_id}, ...]")
	f.Int("batch-size", 1, "Archives per prepared transaction (default 1; Canton interactive prepare supports one command)")
	f.Bool("dry-run", false, "Discover and list targets without archiving")
	f.String("config", "participant-config.json", "Path to participant config JSON file")
	f.String("state-dir", "ceremonies", "Root directory under which ceremony state is stored")

	_ = initArchiveContractsCmd.MarkFlagRequired("decentralized-party-id")

	initCmd.AddCommand(initArchiveContractsCmd)
}

func runInitArchiveContracts(cmd *cobra.Command, _ []string) error {
	f := cmd.Flags()

	partyID, _ := f.GetString("decentralized-party-id")
	synchronizerID, _ := f.GetString("synchronizer-id")
	templateFlags, _ := f.GetStringArray("template")
	inventoryFile, _ := f.GetString("inventory-file")
	batchSize, _ := f.GetInt("batch-size")
	dryRun, _ := f.GetBool("dry-run")
	configPath, _ := f.GetString("config")
	stateDir, _ := f.GetString("state-dir")

	cfg, err := client.LoadConfig(configPath)
	if err != nil {
		return err
	}

	var templates []ledger.TemplateSelector
	for _, raw := range templateFlags {
		tpl, err := archivecontracts.ParseTemplateSelector(strings.TrimSpace(raw))
		if err != nil {
			return err
		}
		templates = append(templates, tpl)
	}

	var targets []ledger.ArchiveTarget
	if inventoryFile != "" {
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			return fmt.Errorf("reading inventory file %q: %w", inventoryFile, err)
		}
		if err := json.Unmarshal(data, &targets); err != nil {
			return fmt.Errorf("parsing inventory file %q: %w", inventoryFile, err)
		}
	}

	if len(templates) == 0 && len(targets) == 0 {
		return fmt.Errorf("at least one --template or --inventory-file is required")
	}

	input := archivecontracts.ArchiveContractsInput{
		DecentralizedPartyID: partyID,
		SynchronizerID:       synchronizerID,
		Templates:            templates,
		Targets:              targets,
		BatchSize:            batchSize,
		DryRun:               dryRun,
	}

	return executeArchiveContractsSequence(cmd.Context(), cfg, input, stateDir, "", confirmerFromFlags(cmd))
}
