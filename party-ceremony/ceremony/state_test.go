package ceremony

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeduplicateReports_DropsRepeatedMismatchAndStaleSequence(t *testing.T) {
	t.Parallel()

	prevP2 := ownershipMismatchReport("old-p2", "p2")
	prevP3 := ownershipMismatchReport("old-p3", "p3")
	prevSeq := thresholdSequenceReport("old-seq", 2, "1/2")
	prevSeq.ChildOperationReports = []string{prevP2.ID, prevP3.ID}
	previous := []operations.Report[any, any]{prevP2, prevP3, prevSeq}

	newP2 := ownershipMismatchReport("new-p2", "p2")
	newP3 := ownershipMismatchReport("new-p3", "p3")
	newSeq := thresholdSequenceReport("new-seq", 2, "1/2")
	newSeq.ChildOperationReports = []string{newP2.ID, newP3.ID}

	got := deduplicateReports(previous, append(previous, newP2, newP3, newSeq))
	assert.Empty(t, got)
}

func TestDeduplicateReports_RewritesChildrenAfterFiltering(t *testing.T) {
	t.Parallel()

	prevP2 := ownershipMismatchReport("old-p2", "p2")
	prevP3 := ownershipMismatchReport("old-p3", "p3")
	prevSeq := thresholdSequenceReport("old-seq", 3, "1/3")
	previous := []operations.Report[any, any]{prevP2, prevP3, prevSeq}

	newP2Success := successReport("new-p2-ok", "example/sign", map[string]any{"participant_id": "p2"}, map[string]any{"signed": true})
	newP3Duplicate := ownershipMismatchReport("new-p3", "p3")
	newSubmit := successReport("new-submit", "example/submit", map[string]any{"proposal": "abc"}, map[string]any{"submitted": true})
	newSeq := testReport("new-seq", "example/sequence", map[string]any{"threshold": 3}, map[string]any{"phase": "signing"}, "signature threshold not met: more signers must resume the sequence: 2/3")
	newSeq.ChildOperationReports = []string{newP2Success.ID, newP3Duplicate.ID, newSubmit.ID}

	got := deduplicateReports(previous, append(previous, newP2Success, newP3Duplicate, newSubmit, newSeq))
	require.Len(t, got, 3)
	assert.Equal(t, []string{newP2Success.ID, newSubmit.ID}, got[2].ChildOperationReports)
}

func TestSaveReportUpdates_NoOpPreservesBytes(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	prevP2 := ownershipMismatchReport("old-p2", "p2")
	prevSeq := thresholdSequenceReport("old-seq", 2, "1/2")
	prevSeq.ChildOperationReports = []string{prevP2.ID}
	previous := []operations.Report[any, any]{prevP2, prevSeq}
	require.NoError(t, SaveReports(dir, previous))

	path := dir + "/" + reportsFilename
	before, err := os.ReadFile(path)
	require.NoError(t, err)

	newP2 := ownershipMismatchReport("new-p2", "p2")
	newSeq := thresholdSequenceReport("new-seq", 2, "1/2")
	newSeq.ChildOperationReports = []string{newP2.ID}
	require.NoError(t, SaveReportUpdates(dir, previous, append(previous, newP2, newSeq)))

	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestSaveReportUpdates_AppendsPreservingPrefix(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	prev := []operations.Report[any, any]{
		ownershipMismatchReport("old-p2", "p2"),
	}
	require.NoError(t, SaveReports(dir, prev))

	path := dir + "/" + reportsFilename
	before, err := os.ReadFile(path)
	require.NoError(t, err)
	beforeClose := bytes.LastIndexByte(before, ']')
	require.NotEqual(t, -1, beforeClose)

	accepted := successReport("new-p2-ok", "example/sign", map[string]any{"participant_id": "p2"}, map[string]any{"signed": true})
	require.NoError(t, SaveReportUpdates(dir, prev, append(prev, accepted)))

	after, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.True(t, bytes.HasPrefix(after, before[:beforeClose]))

	loaded, err := LoadReports(dir)
	require.NoError(t, err)
	require.Len(t, loaded, 2)
	assert.Equal(t, accepted.ID, loaded[1].ID)
}

// --- test helpers ---

func ownershipMismatchReport(id, participantID string) operations.Report[any, any] {
	return testReport(id, "example/sign",
		map[string]any{"participant_id": participantID}, nil,
		"sign-proposal: participant ID does not match client identity")
}

func thresholdSequenceReport(id string, threshold int, progress string) operations.Report[any, any] {
	return testReport(id, "example/sequence",
		map[string]any{"threshold": threshold}, nil,
		fmt.Sprintf("signature threshold not met: more signers must resume the sequence: %s", progress))
}

func successReport(id, defID string, input, output any) operations.Report[any, any] {
	return testReport(id, defID, input, output, "")
}

func testReport(id, defID string, input any, output any, errMsg string) operations.Report[any, any] {
	report := operations.Report[any, any]{
		ID: id,
		Def: operations.Definition{
			ID:          defID,
			Version:     semver.MustParse("1.0.0"),
			Description: "test",
		},
		Input:  input,
		Output: output,
	}
	if errMsg != "" {
		report.Err = &operations.ReportError{Message: errMsg}
	}

	return report
}
