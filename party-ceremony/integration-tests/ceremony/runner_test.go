//nolint:paralleltest
package tests

import (
	"testing"

	integrationtests "github.com/smartcontractkit/chainlink-canton/party-ceremony/integration-tests"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCeremonies(t *testing.T) {
	// Load the Canton chain once and share it across all suites.
	// Each suite's SetupSuite skips the LoadChainWithCTF call when chain is pre-populated.
	chain, err := integrationtests.LoadChainWithCTF(t, 3)
	require.NoError(t, err, "failed to load shared Canton chain")
	require.Len(t, chain.Participants, 3, "expected 3 participants")

	// t.Run("CantonOnboardingFlowTestSuite", func(t *testing.T) {
	// 	s := new(OnboardingFlowTestSuite)
	// 	s.chain = chain
	// 	suite.Run(t, s)
	// })

	// t.Run("CantonKickFlowTestSuite", func(t *testing.T) {
	// 	s := new(KickFlowTestSuite)
	// 	s.chain = chain
	// 	suite.Run(t, s)
	// })

	t.Run("CantonContractDeployFlowTestSuite", func(t *testing.T) {
		s := new(ContractDeployFlowTestSuite)
		s.chain = chain
		suite.Run(t, s)
	})

	// t.Run("CantonAddParticipantFlowTestSuite", func(t *testing.T) {
	// 	s := new(AddParticipantFlowTestSuite)
	// 	s.chain = chain
	// 	suite.Run(t, s)
	// })

	// t.Run("CantonKeyRotationFlowTestSuite", func(t *testing.T) {
	// 	s := new(KeyRotationFlowTestSuite)
	// 	s.chain = chain
	// 	suite.Run(t, s)
	// })
}
