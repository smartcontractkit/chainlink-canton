package rate_limiter

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractTypeInbound = deployment.ContractType("CantonTokenPoolInboundRateLimiter")
var ContractTypeOutbound = deployment.ContractType("CantonTokenPoolOutboundRateLimiter")

var Version = semver.MustParse("1.0.0")

var DeployInbound = contract.NewDeploy(contract.DeployParams[core.RateLimiter]{
	Name:           "canton/ccip/token_pool_inbound_rate_limiter/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractTypeInbound, *Version),
	Description:    "Deploys an Inbound Token Pool Rate Limiter contract on Canton",
	Validate: func(template core.RateLimiter) error {
		if template.PoolInstanceId == "" {
			return fmt.Errorf("PoolInstanceId is required")
		}
		if template.PoolOwner == "" {
			return fmt.Errorf("PoolOwner is required")
		}
		if template.RemoteChainSelector == "" {
			return fmt.Errorf("RemoteChainSelector is required")
		}
		if template.Direction == core.RateLimitDirectionRateLimitDirection_Outbound {
			return fmt.Errorf("cannot use this operation to deploy an outbound rate limiter")
		}
		if template.Mode == "" {
			return fmt.Errorf("mode is required")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "inbound-rate-limiter",
})

var DeployOutbound = contract.NewDeploy(contract.DeployParams[core.RateLimiter]{
	Name:           "canton/ccip/token_pool_outbound_rate_limiter/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractTypeOutbound, *Version),
	Description:    "Deploys an Outbound Token Pool Rate Limiter contract on Canton",
	Validate: func(template core.RateLimiter) error {
		if template.PoolInstanceId == "" {
			return fmt.Errorf("PoolInstanceId is required")
		}
		if template.PoolOwner == "" {
			return fmt.Errorf("PoolOwner is required")
		}
		if template.RemoteChainSelector == "" {
			return fmt.Errorf("RemoteChainSelector is required")
		}
		if template.Direction == core.RateLimitDirectionRateLimitDirection_Inbound {
			return fmt.Errorf("cannot use this operation to deploy an inbound rate limiter")
		}
		if template.Mode == "" {
			return fmt.Errorf("mode is required")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommon),
	Prefix:      "outbound-rate-limiter",
})
