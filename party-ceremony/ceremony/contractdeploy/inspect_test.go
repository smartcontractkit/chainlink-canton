package contractdeploy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/contractdeploy"
)

func bundleWith(t *testing.T, r operations.Reporter) operations.Bundle {
	t.Helper()
	return operations.NewBundle(t.Context, logger.Nop(), r)
}

// TestState_ThresholdNotMet verifies state at the DAR upload gate.
func TestState_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	// p1 uploads DAR, threshold not met.
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, contractdeploy.PhaseDARUpload, s.Phase)
	assert.Len(t, s.DARsUploaded, 1)
}

// TestState_SigningGate verifies state at the signing gate.
func TestState_SigningGate(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	// p1 uploads DAR, threshold not met.
	_, _ = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)

	// p2 uploads DAR, all DARs done, p2 signs (1/2), signing threshold not met.
	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	s := sr.Output.State
	assert.Equal(t, contractdeploy.PhaseSigning, s.Phase)
	assert.Equal(t, 2, s.SignRequired)
	assert.Contains(t, s.Signed, "p2")
	assert.Len(t, s.Signed, 1)
	assert.NotEmpty(t, s.PreparedTxHash)
}

// TestState_Completed verifies the happy-path state snapshot.
func TestState_Completed(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle { return bundleWith(t, reporter) }
	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	// Complete the ceremony (p1 DAR, p2 DAR+sign, p1 sign+execute+verify).
	_, _ = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	_, _ = operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p2", true), input)

	sr, err := operations.ExecuteSequence(newBundle(), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.NoError(t, err)

	s := sr.Output.State
	assert.Equal(t, contractdeploy.PhaseCompleted, s.Phase)
	assert.NotEmpty(t, sr.Output.ContractID)
}

// TestState_LatestSequenceState verifies observer reads.
func TestState_LatestSequenceState(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	input := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})

	_, found := contractdeploy.LatestSequenceState(reporter, input)
	assert.False(t, found)

	_, err := operations.ExecuteSequence(bundleWith(t, reporter), contractdeploy.ContractDeploySequence, newDeps("p1", true), input)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	s, found := contractdeploy.LatestSequenceState(reporter, input)
	require.True(t, found)
	assert.Equal(t, contractdeploy.PhaseDARUpload, s.Phase)
}

// TestState_Scoping verifies two deploy ceremonies don't bleed.
func TestState_Scoping(t *testing.T) {
	t.Parallel()

	reporter := operations.NewMemoryReporter()
	inputA := baseInput([]contractdeploy.PackageRef{{Name: "mcms", Version: "current"}})
	inputB := contractdeploy.ContractDeployInput{
		DecentralizedPartyID: "test-party::1220abcdef",
		SynchronizerID:       "global",
		Packages:             []contractdeploy.PackageRef{{Name: "mcms", Version: "current"}},
		TemplateModule:       "Other.Module",
		TemplateEntity:       "Other",
		ContractArgs:         `{}`,
	}

	_, err := operations.ExecuteSequence(bundleWith(t, reporter), contractdeploy.ContractDeploySequence, newDeps("p1", true), inputA)
	require.ErrorContains(t, err, contractdeploy.ErrThresholdNotMet.Error())

	_, foundA := contractdeploy.LatestSequenceState(reporter, inputA)
	assert.True(t, foundA)

	_, foundB := contractdeploy.LatestSequenceState(reporter, inputB)
	assert.False(t, foundB)
}
