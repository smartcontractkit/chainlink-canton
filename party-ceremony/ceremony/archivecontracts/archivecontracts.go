package archivecontracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// ErrThresholdNotMet is returned when not all hosting participants have signed.
var ErrThresholdNotMet = errors.New("threshold not met: more participants must complete their step")

// ArchiveContractsInput is the top-level input to [ArchiveContractsSequence].
type ArchiveContractsInput struct {
	DecentralizedPartyID string                  `json:"decentralized_party_id"`
	SynchronizerID       string                  `json:"synchronizer_id"`
	Templates            []ledger.TemplateSelector `json:"templates,omitempty"`
	Targets              []ledger.ArchiveTarget    `json:"targets,omitempty"`
	// BatchSize controls how many Archive exercises are bundled per prepared
	// transaction. Zero or unset defaults to 1 because Canton Interactive
	// Submission currently supports only one command per PrepareSubmission.
	BatchSize int  `json:"batch_size,omitempty"`
	DryRun    bool `json:"dry_run,omitempty"`
}

// Phase represents the current execution phase.
type Phase string

const (
	PhaseVerifyParty  Phase = "verify-party"
	PhaseFetchMembers Phase = "fetch-members"
	PhaseDiscover     Phase = "discover"
	PhasePrepare      Phase = "prepare"
	PhaseSigning      Phase = "signing"
	PhaseExecute      Phase = "execute"
	PhaseCompleted    Phase = "completed"
)

// CeremonyState is the live snapshot embedded in every ArchiveContractsOutput.
type CeremonyState struct {
	Phase            Phase    `json:"phase"`
	Participants     []string `json:"participants,omitempty"`
	Signed           []string `json:"signed,omitempty"`
	SignRequired     int      `json:"sign_required,omitempty"`
	PreparedTxHash   string   `json:"prepared_tx_hash,omitempty"`
	CurrentBatch     int      `json:"current_batch,omitempty"`
	TotalBatches     int      `json:"total_batches,omitempty"`
	TargetsDiscovered int     `json:"targets_discovered,omitempty"`
	TargetsArchived  int      `json:"targets_archived,omitempty"`
}

// ArchiveContractsOutput is the result of [ArchiveContractsSequence].
type ArchiveContractsOutput struct {
	Targets                 []ledger.ArchiveTarget `json:"targets,omitempty"`
	PreparedTransactionHash string                 `json:"prepared_transaction_hash,omitempty"`
	ArchivedCount           int                    `json:"archived_count"`
	State                   CeremonyState          `json:"state"`
}

// ArchiveContractsSequence archives active contracts owned by a decentralized
// party using InteractiveSubmissionService (prepare → multiparty sign → execute).
var ArchiveContractsSequence = operations.NewSequence(
	"archive-contracts/canton-ceremony/archive",
	semver.MustParse("1.0.0"),
	"Archive active contracts through a decentralized party (verify → fetch → discover → prepare → sign → execute)",
	func(b operations.Bundle, deps ledger.ContractDeployDeps, in ArchiveContractsInput) (ArchiveContractsOutput, error) {
		out := ArchiveContractsOutput{
			State: CeremonyState{Phase: PhaseVerifyParty},
		}

		if _, err := operations.ExecuteOperation(b, ledger.VerifyPartyOp, deps, ledger.VerifyPartyInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
		}); err != nil {
			return out, fmt.Errorf("verify-party: %w", err)
		}

		out.State.Phase = PhaseFetchMembers
		fetchReport, err := operations.ExecuteOperation(b, ledger.FetchParticipantsOp, deps, ledger.FetchParticipantsInput{
			DecentralizedPartyID: in.DecentralizedPartyID,
			SynchronizerID:       in.SynchronizerID,
		})
		if err != nil {
			return out, fmt.Errorf("fetch-participants: %w", err)
		}

		participants := fetchReport.Output.Participants
		out.State.Participants = participants
		out.State.SignRequired = len(participants)

		uid, err := deps.AdminClient.GetParticipantUID(b.GetContext())
		if err != nil {
			return out, fmt.Errorf("getting participant UID: %w", err)
		}
		isHost := false
		for _, pid := range participants {
			if pid == uid {
				isHost = true
				break
			}
		}
		if isHost && deps.UserID != "" {
			if _, grantErr := operations.ExecuteOperation(b, ledger.GrantPartyRightsOp, deps, ledger.GrantPartyRightsInput{
				ParticipantID:        uid,
				DecentralizedPartyID: in.DecentralizedPartyID,
			}); grantErr != nil {
				deps.Logger.Infow("Party rights grant skipped", "participant", uid, "err", grantErr)
			}
		}

		targets := in.Targets
		if len(targets) == 0 {
			if len(in.Templates) == 0 {
				return out, fmt.Errorf("either templates or explicit targets are required")
			}

			out.State.Phase = PhaseDiscover
			discoverReport, err := operations.ExecuteOperation(b, ledger.DiscoverArchiveTargetsOp, deps, ledger.DiscoverArchiveTargetsInput{
				PartyID:   in.DecentralizedPartyID,
				Templates: in.Templates,
			})
			if err != nil {
				return out, fmt.Errorf("discover-archive-targets: %w", err)
			}
			targets = discoverReport.Output.Targets
		}

		out.Targets = targets
		out.State.TargetsDiscovered = len(targets)

		if len(targets) == 0 {
			out.State.Phase = PhaseCompleted
			deps.Logger.Infow("No active contracts to archive", "party", in.DecentralizedPartyID)

			return out, nil
		}

		if in.DryRun {
			out.State.Phase = PhaseCompleted
			deps.Logger.Infow("Dry run complete", "targets", len(targets))

			return out, nil
		}

		batches := splitTargets(targets, in.BatchSize)
		out.State.TotalBatches = len(batches)

		archived := 0
		for batchIndex, batch := range batches {
			out.State.CurrentBatch = batchIndex
			out.State.Phase = PhasePrepare

			prepReport, err := operations.ExecuteOperation(b, ledger.PrepareArchiveBatchOp, deps, ledger.PrepareArchiveBatchInput{
				DecentralizedPartyID: in.DecentralizedPartyID,
				SynchronizerID:       in.SynchronizerID,
				BatchIndex:           batchIndex,
				Targets:              batch,
			})
			if err != nil {
				return out, fmt.Errorf("prepare archive batch %d: %w", batchIndex, err)
			}

			prep := prepReport.Output
			out.PreparedTransactionHash = prep.PreparedTransactionHash
			out.State.PreparedTxHash = prep.PreparedTransactionHash

			out.State.Phase = PhaseSigning
			signs := make([]ledger.SignSubmissionOutput, 0, len(participants))
			for _, pid := range participants {
				r, signErr := operations.ExecuteOperation(b, ledger.SignSubmissionOp, deps, ledger.SignSubmissionInput{
					ParticipantID:           pid,
					PreparedTransactionHash: prep.PreparedTransactionHash,
					PreparedTxB64:           prep.PreparedTxB64,
					KnownSigningKeysB64:     fetchReport.Output.PartySigningKeysB64,
				})
				if signErr != nil {
					deps.Logger.Infow("Signing pending", "participant", pid, "batch", batchIndex, "err", signErr)

					continue
				}
				signs = append(signs, r.Output)
				out.State.Signed = append(out.State.Signed, pid)
			}

			if len(signs) < len(participants) {
				deps.Logger.Warnw("Not all participants have signed",
					"signed", len(signs), "required", len(participants), "batch", batchIndex)

				return out, fmt.Errorf("%w: %d/%d participants have signed batch %d",
					ErrThresholdNotMet, len(signs), len(participants), batchIndex)
			}

			sigsB64 := make([]string, len(signs))
			for i, s := range signs {
				sigsB64[i] = s.SignatureB64
			}

			out.State.Phase = PhaseExecute
			if _, err := operations.ExecuteOperation(b, ledger.ExecuteSubmissionOp, deps, ledger.ExecuteSubmissionInput{
				DecentralizedPartyID: in.DecentralizedPartyID,
				PreparedTxB64:        prep.PreparedTxB64,
				SignaturesB64:        sigsB64,
				HashingSchemeVersion: prep.HashingSchemeVersion,
			}); err != nil {
				return out, fmt.Errorf("execute archive batch %d: %w", batchIndex, err)
			}

			archived += len(batch)
			out.State.TargetsArchived = archived
			out.State.Signed = nil
		}

		out.ArchivedCount = archived
		out.State.Phase = PhaseCompleted
		deps.Logger.Infow("Archive ceremony complete",
			"party", in.DecentralizedPartyID,
			"archived", archived,
		)

		return out, nil
	},
)

func splitTargets(targets []ledger.ArchiveTarget, batchSize int) [][]ledger.ArchiveTarget {
	if batchSize <= 0 {
		batchSize = 1
	}

	batches := make([][]ledger.ArchiveTarget, 0, (len(targets)+batchSize-1)/batchSize)
	for i := 0; i < len(targets); i += batchSize {
		end := i + batchSize
		if end > len(targets) {
			end = len(targets)
		}
		batches = append(batches, targets[i:end])
	}

	return batches
}

// SplitTargetsForTest exposes [splitTargets] for unit tests.
func SplitTargetsForTest(targets []ledger.ArchiveTarget, batchSize int) [][]ledger.ArchiveTarget {
	return splitTargets(targets, batchSize)
}

// ParseTemplateSelector parses packageName:Module:Entity or packageId:Module:Entity.
// A leading "#" on the package component denotes a package name reference.
func ParseTemplateSelector(raw string) (ledger.TemplateSelector, error) {
	parts := splitTemplate(raw)
	if len(parts) != 3 {
		return ledger.TemplateSelector{}, fmt.Errorf("template %q must have format package:Module:Entity", raw)
	}

	pkg := parts[0]
	tpl := ledger.TemplateSelector{
		ModuleName: parts[1],
		EntityName: parts[2],
	}
	if len(pkg) > 0 && pkg[0] == '#' {
		tpl.PackageName = pkg[1:]
	} else if len(pkg) == 64 {
		tpl.PackageID = pkg
	} else {
		tpl.PackageName = pkg
	}

	if tpl.PackageName == "" && tpl.PackageID == "" {
		return ledger.TemplateSelector{}, fmt.Errorf("template %q missing package component", raw)
	}

	return tpl, nil
}

func splitTemplate(raw string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(raw); i++ {
		if raw[i] == ':' {
			parts = append(parts, raw[start:i])
			start = i + 1
		}
	}
	parts = append(parts, raw[start:])

	return parts
}
