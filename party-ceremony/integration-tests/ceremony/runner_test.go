//nolint:paralleltest
package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	integrationtests "github.com/smartcontractkit/chainlink-canton/party-ceremony/integration-tests"
)

func TestCeremonies(t *testing.T) {
	// Load the Canton chain once and share it across all suites.
	// Each suite's SetupSuite skips CTF startup when chain is pre-populated.
	env, err := integrationtests.LoadChainWithCTFKMS(t, 3)
	require.NoError(t, err, "failed to load shared KMS-backed Canton chain")
	chain := env.Chain
	require.Len(t, chain.Participants, 3, "expected 3 participants")

	t.Run("CantonKMSHarnessSmokeSuite", func(t *testing.T) {
		s := new(KMSHarnessSmokeSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "smoke"
		suite.Run(t, s)
	})

	t.Run("CantonOnboardingFlowTestSuite", func(t *testing.T) {
		s := new(OnboardingFlowTestSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "onboarding"
		suite.Run(t, s)
	})

	t.Run("CantonKickFlowTestSuite", func(t *testing.T) {
		s := new(KickFlowTestSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "kick"
		suite.Run(t, s)
	})

	t.Run("CantonContractDeployFlowTestSuite", func(t *testing.T) {
		s := new(ContractDeployFlowTestSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "contract-deploy"
		suite.Run(t, s)
	})

	t.Run("CantonAddParticipantFlowTestSuite", func(t *testing.T) {
		s := new(AddParticipantFlowTestSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "add-participant"
		suite.Run(t, s)
	})

	t.Run("CantonKeyRotationFlowTestSuite", func(t *testing.T) {
		s := new(KeyRotationFlowTestSuite)
		s.chain = chain
		s.KMS = env.KMS
		s.KMSRunName = "key-rotation"
		suite.Run(t, s)
	})
}
