//nolint:paralleltest
package tests

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
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

	runCeremonySuites(t, chain, env.KMS, "", true)
}

func TestCeremoniesWithoutKMS(t *testing.T) {
	chain, err := integrationtests.LoadChainWithCTF(t, 3)
	require.NoError(t, err, "failed to load shared non-KMS Canton chain")
	require.Len(t, chain.Participants, 3, "expected 3 participants")

	runCeremonySuites(t, chain, nil, "non-kms-", false)
}

func runCeremonySuites(t *testing.T, chain *canton.Chain, kms *integrationtests.KMSRegistry, runNamePrefix string, includeKMSSmoke bool) {
	t.Helper()

	if includeKMSSmoke {
		t.Run("CantonKMSHarnessSmokeSuite", func(t *testing.T) {
			s := new(KMSHarnessSmokeSuite)
			s.chain = chain
			s.KMS = kms
			s.KMSRunName = runNamePrefix + "smoke"
			suite.Run(t, s)
		})
	}

	t.Run("CantonOnboardingFlowTestSuite", func(t *testing.T) {
		s := new(OnboardingFlowTestSuite)
		s.chain = chain
		s.KMS = kms
		s.KMSRunName = runNamePrefix + "onboarding"
		suite.Run(t, s)
	})

	t.Run("CantonKickFlowTestSuite", func(t *testing.T) {
		s := new(KickFlowTestSuite)
		s.chain = chain
		s.KMS = kms
		s.KMSRunName = runNamePrefix + "kick"
		suite.Run(t, s)
	})

	t.Run("CantonContractDeployFlowTestSuite", func(t *testing.T) {
		s := new(ContractDeployFlowTestSuite)
		s.chain = chain
		s.KMS = kms
		s.KMSRunName = runNamePrefix + "contract-deploy"
		suite.Run(t, s)
	})

	t.Run("CantonAddParticipantFlowTestSuite", func(t *testing.T) {
		s := new(AddParticipantFlowTestSuite)
		s.chain = chain
		s.KMS = kms
		s.KMSRunName = runNamePrefix + "add-participant"
		suite.Run(t, s)
	})

	t.Run("CantonKeyRotationFlowTestSuite", func(t *testing.T) {
		s := new(KeyRotationFlowTestSuite)
		s.chain = chain
		s.KMS = kms
		s.KMSRunName = runNamePrefix + "key-rotation"
		suite.Run(t, s)
	})
}
