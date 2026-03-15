package rate_limiter

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("RateLimiter")

var Version = semver.MustParse("1.0.0")

var Deploy = contract.NewDeploy(contract.DeployParams[common.RateLimiter]{
	Name:           "canton/ccip/rate_limiter/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP RateLimiter contract on Canton",
	Validate: func(template common.RateLimiter) error {
		if template.PoolInstanceId == "" {
			return errors.New("poolInstanceId cannot be empty")
		}
		if template.PoolOwner == "" {
			return errors.New("poolOwner cannot be empty")
		}
		if template.RemoteChainSelector == "" {
			return errors.New("remoteChainSelector cannot be empty")
		}
		if template.Direction == "" {
			return errors.New("direction cannot be empty")
		}
		if template.Mode == "" {
			return errors.New("mode cannot be empty")
		}
		if template.Capacity == "" {
			return errors.New("capacity cannot be empty")
		}
		if template.Rate == "" {
			return errors.New("rate cannot be empty")
		}
		if template.Tokens == "" {
			return errors.New("tokens cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "ratelimiter",
})
