package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/example"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/kick"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// executeSequence is the shared execution kernel used by both init and resume.
// It:
//  1. Loads any previously persisted reports from <ceremonyDir>/reports.json
//     (empty on first run).
//  2. Builds a MemoryReporter seeded with those reports so all previously
//     successful operations are served from cache without re-execution.
//  3. Runs OnboardingSequence.
//  4. Persists the updated report set back to reports.json — even on error so
//     that partial progress is not lost.
func executeExampleOnboardingSequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input example.OnboardingInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	if err != nil {
		return fmt.Errorf("loading reports: %w", err)
	}
	initialRun := workflowId == ""
	// ── Build bundle ──────────────────────────────────────────────────────────
	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(), // We don't want default operation logs to be noisy
		reporter,
	)

	// ── Build deps ────────────────────────────────────────────────────────────
	// NewMockCantonClientFromConfig is the injection point: swap this for a real
	// gRPC client once the Canton admin API integration is ready.
	client := example.NewMockCantonClientFromConfig(cfg)
	deps := example.CantonDeps{Client: client, Logger: lggr, Confirmer: confirmer}

	lggr.Infow("Running onboarding sequence",
		"ceremony", input.NamespaceName,
		"participant", cfg.ParticipantID,
		"participants", input.Participants,
	)

	// ── Execute sequence ──────────────────────────────────────────────────────
	sr, seqErr := operations.ExecuteSequence(bundle, example.OnboardingSequence, deps, input)

	// ── Persist reports (always, even on error) ───────────────────────────────
	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if err := os.MkdirAll(ceremonyDir, 0o755); err != nil {
		lggr.Errorw("Failed to create ceremony directory — progress may be lost on next resume",
			"dir", ceremonyDir,
			"err", err,
		)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports — progress may be lost on next resume", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports — progress may be lost on next resume",
			"err", saveErr)
	}

	// Store workflow.json on first run so resume can verify the input matches on subsequent runs.
	if initialRun {
		state := ceremony.WorkflowState[example.OnboardingInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeExample,
			Input:      input,
		}
		if err := ceremony.SaveWorkflow(ceremonyDir, state); err != nil {
			lggr.Errorw("Failed to save workflow.json — resume may fail if input cannot be reconstructed",
				"err", err)
		}
	}

	if seqErr != nil {
		// ErrThresholdNotMet is expected when not all signers have acted yet.
		// Print a friendly message and exit with a non-zero code without a stack trace.
		if strings.Contains(seqErr.Error(), example.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more participants have signed.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}

// executeOnboardingSequence is the execution kernel for the real gRPC-backed
// onboarding ceremony. It replaces the mock client with a real GRPCCantonClient
// that connects to the participant's Admin API.
func executeOnboardingSequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input onboarding.OnboardingInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	initialRun := workflowId == ""

	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(),
		reporter,
	)

	// ── Build real gRPC client ────────────────────────────────────────────
	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	grpcClient := client.NewGRPCClient(conn)
	deps := ceremony.CantonDeps{Client: grpcClient, KMS: cfg.KMS(), Logger: lggr, Confirmer: confirmer}

	lggr.Infow("Running onboarding sequence (real)",
		"ceremony", input.NamespaceName,
		"participant", cfg.ParticipantID,
		"participants", input.Participants,
		"kms_namespace_key", cfg.KmsNamespaceKeyID != "",
		"kms_protocol_key", cfg.KmsProtocolKeyID != "",
	)

	sr, seqErr := operations.ExecuteSequence(bundle, onboarding.OnboardingSequence, deps, input)

	// ── Persist reports (always, even on error) ───────────────────────────
	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if mkErr := os.MkdirAll(ceremonyDir, 0o755); mkErr != nil {
		lggr.Errorw("Failed to create ceremony directory", "dir", ceremonyDir, "err", mkErr)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports", "err", saveErr)
	}

	if initialRun {
		state := ceremony.WorkflowState[onboarding.OnboardingInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeOnboarding,
			Input:      input,
		}
		if saveErr := ceremony.SaveWorkflow(ceremonyDir, state); saveErr != nil {
			lggr.Errorw("Failed to save workflow.json", "err", saveErr)
		}
	}

	if seqErr != nil {
		if strings.Contains(seqErr.Error(), onboarding.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more participants have signed.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}

func calculateCeremonyDir(stateDir, workflowId string) string {
	return fmt.Sprintf("%s/%s", stateDir, workflowId)
}

// executeKickSequence is the execution kernel for the real gRPC-backed kick
// ceremony. It dials the Canton admin API, runs KickSequence, and persists
// the ceremony state.
func executeKickSequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input kick.KickInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	initialRun := workflowId == ""

	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(),
		reporter,
	)

	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	grpcClient := client.NewGRPCClient(conn)
	deps := ceremony.CantonDeps{Client: grpcClient, KMS: cfg.KMS(), Logger: lggr, Confirmer: confirmer}

	lggr.Infow("Running kick sequence",
		"party", input.DecentralizedPartyID,
		"kicked_participant", input.KickedParticipantID,
		"remaining_count", len(input.RemainingParticipants),
		"participant", cfg.ParticipantID,
	)

	sr, seqErr := operations.ExecuteSequence(bundle, kick.KickSequence, deps, input)

	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if mkErr := os.MkdirAll(ceremonyDir, 0o755); mkErr != nil {
		lggr.Errorw("Failed to create ceremony directory", "dir", ceremonyDir, "err", mkErr)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports", "err", saveErr)
	}

	if initialRun {
		state := ceremony.WorkflowState[kick.KickInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeKick,
			Input:      input,
		}
		if saveErr := ceremony.SaveWorkflow(ceremonyDir, state); saveErr != nil {
			lggr.Errorw("Failed to save workflow.json", "err", saveErr)
		}
	}

	if seqErr != nil {
		if strings.Contains(seqErr.Error(), kick.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "kick ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more remaining participants have acted.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}

// executeContractDeploySequence is the execution kernel for the contract
// deployment ceremony. It dials both the Admin gRPC (for DAR uploads) and
// Ledger gRPC (for interactive submission) connections, runs
// ContractDeploySequence, and persists the ceremony state.
func executeContractDeploySequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input contractdeploy.ContractDeployInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	initialRun := workflowId == ""

	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(),
		reporter,
	)

	// ── Build Admin gRPC client (for DAR uploads) ─────────────────────
	adminConn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer adminConn.Close()

	// ── Build Ledger gRPC client (for interactive submission) ─────────
	ledgerConn, err := client.DialLedger(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton ledger API: %w", err)
	}
	defer ledgerConn.Close()

	adminClient := client.NewGRPCClient(adminConn)
	var kmsAPI client.AWSKMSAPI
	if cfg.KmsProtocolKeyID != "" {
		awsCfg, cfgErr := awsconfig.LoadDefaultConfig(ctx)
		if cfgErr != nil {
			return fmt.Errorf("loading AWS config for KMS signing: %w", cfgErr)
		}
		kmsAPI = awskms.NewFromConfig(awsCfg)
	}

	deps := ledger.ContractDeployDeps{
		AdminClient:  adminClient,
		LedgerClient: client.NewGRPCLedgerClient(ledgerConn),
		// DARLoader reads DAR bytes by package name and version.
		// Callers that embed the contracts FS can swap this for contracts.GetDar.
		DARLoader: ledger.FileDARLoader("dars"),
		SignerFactory: client.NewTransactionSignerFactory(
			adminClient,
			cryptoadminv30.NewVaultServiceClient(adminConn),
			cfg.KMS(),
			kmsAPI,
		),
		Logger:    lggr,
		Confirmer: confirmer,
		UserID:    cfg.UserID,
	}

	lggr.Infow("Running contract-deploy sequence",
		"party", input.DecentralizedPartyID,
		"participant", cfg.ParticipantID,
		"package_count", len(input.Packages),
		"template", input.TemplateModule+":"+input.TemplateEntity,
		"kms_protocol_key", cfg.KmsProtocolKeyID != "",
	)

	sr, seqErr := operations.ExecuteSequence(bundle, contractdeploy.ContractDeploySequence, deps, input)

	// ── Persist reports (always, even on error) ───────────────────────
	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if mkErr := os.MkdirAll(ceremonyDir, 0o755); mkErr != nil {
		lggr.Errorw("Failed to create ceremony directory", "dir", ceremonyDir, "err", mkErr)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports", "err", saveErr)
	}

	if initialRun {
		state := ceremony.WorkflowState[contractdeploy.ContractDeployInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeContractDeploy,
			Input:      input,
		}
		if saveErr := ceremony.SaveWorkflow(ceremonyDir, state); saveErr != nil {
			lggr.Errorw("Failed to save workflow.json", "err", saveErr)
		}
	}

	if seqErr != nil {
		if strings.Contains(seqErr.Error(), contractdeploy.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "contract-deploy ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more participants have uploaded DARs.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}

// executeAddParticipantSequence is the execution kernel for the real gRPC-backed
// add-participant ceremony. It dials the Canton admin API, runs
// AddParticipantSequence, and persists the ceremony state.
func executeAddParticipantSequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input addparticipant.AddParticipantInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	initialRun := workflowId == ""

	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(),
		reporter,
	)

	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	grpcClient := client.NewGRPCClient(conn)
	deps := ceremony.CantonDeps{Client: grpcClient, KMS: cfg.KMS(), Logger: lggr, Confirmer: confirmer}

	lggr.Infow("Running add-participant sequence",
		"party", input.DecentralizedPartyID,
		"new_participant", input.NewParticipantID,
		"participant", cfg.ParticipantID,
	)

	sr, seqErr := operations.ExecuteSequence(bundle, addparticipant.AddParticipantSequence, deps, input)

	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if mkErr := os.MkdirAll(ceremonyDir, 0o755); mkErr != nil {
		lggr.Errorw("Failed to create ceremony directory", "dir", ceremonyDir, "err", mkErr)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports", "err", saveErr)
	}

	if initialRun {
		state := ceremony.WorkflowState[addparticipant.AddParticipantInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeAddParticipant,
			Input:      input,
		}
		if saveErr := ceremony.SaveWorkflow(ceremonyDir, state); saveErr != nil {
			lggr.Errorw("Failed to save workflow.json", "err", saveErr)
		}
	}

	if seqErr != nil {
		if strings.Contains(seqErr.Error(), addparticipant.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "add-participant ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more participants have acted.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}

// executeKeyRotationSequence is the execution kernel for the real gRPC-backed
// key rotation ceremony. It dials the Canton admin API, runs
// KeyRotationSequence, and persists the ceremony state.
func executeKeyRotationSequence(
	ctx context.Context,
	cfg client.ClientConfig,
	input keyrotation.KeyRotationInput,
	stateDir string,
	workflowId string,
	confirmer ceremony.Confirmer,
) error {
	lggr, err := newCLILogger()
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	var previousReports []operations.Report[any, any]
	if workflowId != "" {
		ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
		previousReports, err = ceremony.LoadReports(ceremonyDir)
		if err != nil {
			return fmt.Errorf("loading previous reports: %w", err)
		}
	}
	initialRun := workflowId == ""

	reporter := operations.NewMemoryReporter(operations.WithReports(previousReports))
	bundle := operations.NewBundle(
		func() context.Context { return ctx },
		logger.Nop(),
		reporter,
	)

	conn, err := client.Dial(cfg)
	if err != nil {
		return fmt.Errorf("connecting to Canton admin API: %w", err)
	}
	defer conn.Close()

	grpcClient := client.NewGRPCClient(conn)
	deps := ceremony.CantonDeps{Client: grpcClient, KMS: cfg.KMS(), Logger: lggr, Confirmer: confirmer}

	lggr.Infow("Running key-rotation sequence",
		"party", input.DecentralizedPartyID,
		"target_participant", input.TargetParticipantID,
		"rotate_namespace", input.RotateNamespaceKey,
		"rotate_daml", input.RotateDamlKey,
		"participant", cfg.ParticipantID,
	)

	sr, seqErr := operations.ExecuteSequence(bundle, keyrotation.KeyRotationSequence, deps, input)

	if workflowId == "" {
		workflowId = sr.ID
	}
	ceremonyDir := calculateCeremonyDir(stateDir, workflowId)
	if mkErr := os.MkdirAll(ceremonyDir, 0o755); mkErr != nil {
		lggr.Errorw("Failed to create ceremony directory", "dir", ceremonyDir, "err", mkErr)
	}
	allReports, reportErr := reporter.GetReports()
	if reportErr != nil {
		lggr.Errorw("Failed to collect reports", "err", reportErr)
	}
	if saveErr := ceremony.SaveReportUpdates(ceremonyDir, previousReports, allReports); saveErr != nil {
		lggr.Errorw("Failed to save reports", "err", saveErr)
	}

	if initialRun {
		state := ceremony.WorkflowState[keyrotation.KeyRotationInput]{
			CeremonyID: workflowId,
			Type:       ceremony.WorkflowTypeKeyRotation,
			Input:      input,
		}
		if saveErr := ceremony.SaveWorkflow(ceremonyDir, state); saveErr != nil {
			lggr.Errorw("Failed to save workflow.json", "err", saveErr)
		}
	}

	if seqErr != nil {
		if strings.Contains(seqErr.Error(), keyrotation.ErrThresholdNotMet.Error()) {
			fmt.Fprintf(os.Stderr, "key-rotation ceremony not yet complete: %v\n", seqErr)
			fmt.Fprintln(os.Stderr, "Run `resume` again after more participants have acted.")
			os.Exit(2) //nolint:gocritic // intentional early exit for UX
		}

		return seqErr
	}

	return nil
}
