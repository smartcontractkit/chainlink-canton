package ceremony

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	// WorkflowTypeContractDeploy identifies a contract deployment ceremony.
	WorkflowTypeContractDeploy = "contract-deploy"
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
