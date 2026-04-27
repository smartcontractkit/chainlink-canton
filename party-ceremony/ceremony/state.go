package ceremony

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

const (
	workflowFilename = "workflow.json"
	reportsFilename  = "reports.json"

	// WorkflowTypeOnboarding identifies a real gRPC-backed onboarding ceremony.
	WorkflowTypeOnboarding = "onboarding"

	// WorkflowTypeExample identifies a mock-backed example onboarding ceremony.
	WorkflowTypeExample = "example"

	// WorkflowTypeKick identifies a real gRPC-backed kick (participant removal) ceremony.
	WorkflowTypeKick = "kick"

	// WorkflowTypeAddParticipant identifies a real gRPC-backed add-participant ceremony.
	WorkflowTypeAddParticipant = "add-participant"

	// WorkflowTypeContractDeploy identifies a contract deployment ceremony.
	WorkflowTypeContractDeploy = "contract-deploy"

	// WorkflowTypeKeyRotation identifies a key rotation ceremony.
	WorkflowTypeKeyRotation = "key-rotation"
)

// WorkflowState is persisted to workflow.json inside the ceremony directory.
// It holds everything needed to reconstruct the ceremony input when resuming.
// T is the ceremony-specific input type (e.g. OnboardingInput).
type WorkflowState[T any] struct {
	// CeremonyID is the unique identifier for this ceremony run.
	CeremonyID string `json:"ceremony_id"`

	// Type identifies the workflow kind (e.g. "onboarding", "example").
	// resume uses this to dispatch to the correct sequence executor without
	// requiring a --real / --type flag on the command line.
	Type string `json:"type"`

	// Input is the full ceremony input passed to the sequence.
	// Persisting it means resume only needs the ceremony ID — all parameters
	// are loaded from workflow.json rather than re-supplied on the command line.
	Input T `json:"input"`
}

// PeekWorkflowType reads only the "type" field from <dir>/workflow.json
// without deserialising the full generic Input. Used by resume to dispatch to
// the correct sequence executor.
func PeekWorkflowType(dir string) (string, error) {
	path := filepath.Join(dir, workflowFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading %q: %w", path, err)
	}

	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return "", fmt.Errorf("parsing workflow type from %q: %w", path, err)
	}

	if peek.Type == "" {
		return "", fmt.Errorf("workflow %q has no type field; re-run init to recreate it", path)
	}

	return peek.Type, nil
}

// SaveWorkflow writes state to <dir>/workflow.json.
// The file is written atomically via a rename to avoid partial writes.
func SaveWorkflow[T any](dir string, state WorkflowState[T]) error {
	return writeJSONAtomic(filepath.Join(dir, workflowFilename), state)
}

// LoadWorkflow reads a WorkflowState from <dir>/workflow.json.
func LoadWorkflow[T any](dir string) (WorkflowState[T], error) {
	path := filepath.Join(dir, workflowFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return WorkflowState[T]{}, fmt.Errorf("reading %q: %w", path, err)
	}

	var state WorkflowState[T]
	if err := json.Unmarshal(data, &state); err != nil {
		return WorkflowState[T]{}, fmt.Errorf("parsing %q: %w", path, err)
	}

	return state, nil
}

// SaveReports serialises the reporter's Report slice to <dir>/reports.json.
//
// Reports are the idempotency store used by the Operations framework.
// Persisting them across CLI invocations means that any operation that already
// succeeded is not re-executed on the next resume call.
func SaveReports(dir string, reports []operations.Report[any, any]) error {
	return writeJSONAtomic(filepath.Join(dir, reportsFilename), reports)
}

// SaveReportUpdates persists reports while preserving existing reports.json
// bytes on resume. It appends only deduplicated new reports and skips
// repeated ownership-mismatch failures and stale threshold-not-met sequences.
func SaveReportUpdates(
	dir string,
	previousReports []operations.Report[any, any],
	allReports []operations.Report[any, any],
) error {
	if len(previousReports) == 0 {
		return SaveReports(dir, allReports)
	}

	reportsToAppend := deduplicateReports(previousReports, allReports)
	if len(reportsToAppend) == 0 {
		return nil
	}

	return appendReportsAtomic(filepath.Join(dir, reportsFilename), reportsToAppend)
}

// LoadReports deserialises reports from <dir>/reports.json.
// Returns an empty slice (not an error) when the file does not yet exist — this
// is the normal first-run case before any operations have been executed.
func LoadReports(dir string) ([]operations.Report[any, any], error) {
	path := filepath.Join(dir, reportsFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	var reports []operations.Report[any, any]
	if err := json.Unmarshal(data, &reports); err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	return reports, nil
}

// deduplicateReports returns only the reports from allReports that represent
// new progress compared to previousReports. It filters out:
//   - ownership-mismatch failures ("participant ID mismatch" / "does not match
//     client identity") that were already recorded with the same input;
//   - threshold-not-met sequence reports whose children were all filtered and
//     the same error was already seen.
//
// When a sequence report survives filtering, its ChildOperationReports are
// rewritten to reference only reports that will actually be persisted.
func deduplicateReports(
	previousReports []operations.Report[any, any],
	allReports []operations.Report[any, any],
) []operations.Report[any, any] {
	if len(allReports) <= len(previousReports) {
		return nil
	}

	idx := newDedupIndex(previousReports)

	newReports := allReports[len(previousReports):]
	accepted := make([]operations.Report[any, any], 0, len(newReports))
	for _, report := range newReports {
		if idx.isDuplicateOwnershipMismatch(report) {
			continue
		}
		if len(report.ChildOperationReports) > 0 {
			report.ChildOperationReports = idx.filterPersistedChildren(report.ChildOperationReports)
		}
		if idx.isDuplicateThresholdSequence(report) {
			continue
		}

		accepted = append(accepted, report)
		idx.trackSaved(report.ID)
	}

	return accepted
}

// dedupIndex tracks which reports and error patterns have already been
// persisted, so the filtering loop can decide what to skip in O(1).
type dedupIndex struct {
	savedIDs            map[string]struct{}
	ownershipMismatches map[string]struct{}
	thresholdErrors     map[string]struct{}
}

func newDedupIndex(previousReports []operations.Report[any, any]) dedupIndex {
	idx := dedupIndex{
		savedIDs:            make(map[string]struct{}, len(previousReports)),
		ownershipMismatches: make(map[string]struct{}, len(previousReports)),
		thresholdErrors:     make(map[string]struct{}, len(previousReports)),
	}
	for _, r := range previousReports {
		idx.savedIDs[r.ID] = struct{}{}
		if isOwnershipMismatchErr(r) {
			idx.ownershipMismatches[dedupKeyByInput(r)] = struct{}{}
		}
		if isThresholdNotMetErr(r) {
			idx.thresholdErrors[dedupKeyByInputAndError(r)] = struct{}{}
		}
	}

	return idx
}

func (idx *dedupIndex) trackSaved(id string) {
	idx.savedIDs[id] = struct{}{}
}

// isDuplicateOwnershipMismatch returns true if the report is an ownership-
// mismatch failure that was already seen with the same definition+input.
// On first sight it records the key and returns false.
func (idx *dedupIndex) isDuplicateOwnershipMismatch(r operations.Report[any, any]) bool {
	if !isOwnershipMismatchErr(r) {
		return false
	}
	key := dedupKeyByInput(r)
	if _, exists := idx.ownershipMismatches[key]; exists {
		return true
	}
	idx.ownershipMismatches[key] = struct{}{}

	return false
}

// isDuplicateThresholdSequence returns true if the report is a threshold-not-
// met error with no remaining children (all were duplicates) and the same error
// was already persisted.
func (idx *dedupIndex) isDuplicateThresholdSequence(r operations.Report[any, any]) bool {
	if !isThresholdNotMetErr(r) || len(r.ChildOperationReports) > 0 {
		return false
	}
	key := dedupKeyByInputAndError(r)
	if _, exists := idx.thresholdErrors[key]; exists {
		return true
	}
	idx.thresholdErrors[key] = struct{}{}

	return false
}

// filterPersistedChildren keeps only child IDs that reference already-saved
// reports, so the persisted sequence doesn't contain dangling references.
func (idx *dedupIndex) filterPersistedChildren(childIDs []string) []string {
	kept := make([]string, 0, len(childIDs))
	for _, id := range childIDs {
		if _, ok := idx.savedIDs[id]; ok {
			kept = append(kept, id)
		}
	}

	return kept
}

func isOwnershipMismatchErr(r operations.Report[any, any]) bool {
	if r.Err == nil {
		return false
	}
	msg := r.Err.Message

	return strings.Contains(msg, "participant ID mismatch") ||
		strings.Contains(msg, "does not match client identity")
}

func isThresholdNotMetErr(r operations.Report[any, any]) bool {
	return r.Err != nil &&
		strings.Contains(r.Err.Message, "threshold not met")
}

func dedupKeyByInput(r operations.Report[any, any]) string {
	return normaliseJSON(struct {
		Definition any `json:"definition"`
		Input      any `json:"input"`
	}{
		Definition: r.Def,
		Input:      r.Input,
	})
}

func dedupKeyByInputAndError(r operations.Report[any, any]) string {
	return dedupKeyByInput(r) + "\x00" + r.Err.Message
}

// normaliseJSON marshals v to JSON, then round-trips through json.Decoder with
// UseNumber to ensure Go int literals and float64 values produce identical keys
// (the tests pass map[string]any with Go ints while loaded reports use float64).
func normaliseJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	var decoded any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&decoded); err != nil {
		return string(data)
	}

	canonical, err := json.Marshal(decoded)
	if err != nil {
		return string(data)
	}

	return string(canonical)
}

func appendReportsAtomic(path string, reports []operations.Report[any, any]) error {
	if len(reports) == 0 {
		return nil
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return writeJSONAtomic(path, reports)
		}

		return fmt.Errorf("reading %q: %w", path, err)
	}

	closeIdx := bytes.LastIndexByte(existing, ']')
	if closeIdx < 0 {
		return fmt.Errorf("appending to %q: reports file is not a JSON array", path)
	}

	var addition bytes.Buffer
	if bytes.Contains(bytes.TrimSpace(existing[:closeIdx]), []byte("{")) {
		addition.WriteString(",\n")
	} else {
		addition.WriteByte('\n')
	}

	for i, report := range reports {
		if i > 0 {
			addition.WriteString(",\n")
		}
		reportJSON, err := json.MarshalIndent(report, "  ", "  ")
		if err != nil {
			return fmt.Errorf("marshalling report for %q: %w", path, err)
		}
		addition.Write(reportJSON)
	}
	if addition.Len() > 0 && addition.Bytes()[addition.Len()-1] != '\n' {
		addition.WriteByte('\n')
	}

	updated := make([]byte, 0, len(existing)+addition.Len())
	updated = append(updated, existing[:closeIdx]...)
	updated = append(updated, addition.Bytes()...)
	updated = append(updated, existing[closeIdx:]...)

	tmp := path + ".tmp"
	// #nosec G703 -- tmp path is constructed from validated ceremonyDir passed by caller
	if err := os.WriteFile(tmp, updated, 0o600); err != nil {
		return fmt.Errorf("writing %q: %w", tmp, err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("renaming %q → %q: %w", tmp, path, err)
	}

	return nil
}

// writeJSONAtomic marshals v to JSON and writes it to path, creating or
// truncating the file.
func writeJSONAtomic(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling to %q: %w", path, err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("writing %q: %w", tmp, err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("renaming %q → %q: %w", tmp, path, err)
	}

	return nil
}
