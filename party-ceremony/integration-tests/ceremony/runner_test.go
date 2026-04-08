//nolint:paralleltest
package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCeremonies(t *testing.T) {
	t.Run("CantonOnboardingFlowTestSuite", func(t *testing.T) {
		suite.Run(t, new(OnboardingFlowTestSuite))
	})

	t.Run("CantonKickFlowTestSuite", func(t *testing.T) {
		suite.Run(t, new(KickFlowTestSuite))
	})

	t.Run("CantonContractDeployFlowTestSuite", func(t *testing.T) {
		suite.Run(t, new(ContractDeployFlowTestSuite))
	})
}
