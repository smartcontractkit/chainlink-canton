package changesets

import (
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	v2_0_0 "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ApplyExecutorConfigAllowingCantonMinNOPs wraps apply-executor-config for prod Canton lanes:
//   - skips production min-NOP topology validation during VerifyPreconditions (Canton committee/pool)
//   - skips JD ListNodeChainConfigs validation during Apply so Canton destination blocks can be
//     pushed to EVM ccvexecutor jobs before Canton node chain configs exist in JD
func ApplyExecutorConfigAllowingCantonMinNOPs(registry *adapters.ExecutorConfigRegistry) deployment.ChangeSetV2[v2_0_0.ApplyExecutorConfigInput] {
	inner := v2_0_0.ApplyExecutorConfig(registry)

	return deployment.CreateChangeSet(
		func(e deployment.Environment, cfg v2_0_0.ApplyExecutorConfigInput) (deployment.ChangesetOutput, error) {
			return applyExecutorConfigSkippingChainSupportValidation(registry, e, cfg)
		},
		func(e deployment.Environment, cfg v2_0_0.ApplyExecutorConfigInput) error {
			return WithCantonProductionMinNOPCheckBypassed(cfg.Topology, func() error {
				return inner.VerifyPreconditions(e, cfg)
			})
		},
	)
}

// ApplyVerifierConfigAllowingCantonMinNOPs wraps apply-verifier-config with the same Canton production
// minimum-NOP bypass as lane configure and apply-executor-config.
func ApplyVerifierConfigAllowingCantonMinNOPs(registry *adapters.VerifierConfigRegistry) deployment.ChangeSetV2[v2_0_0.ApplyVerifierConfigInput] {
	inner := v2_0_0.ApplyVerifierConfig(registry)

	return deployment.CreateChangeSet(
		inner.Apply,
		func(e deployment.Environment, cfg v2_0_0.ApplyVerifierConfigInput) error {
			return WithCantonProductionMinNOPCheckBypassed(cfg.Topology, func() error {
				return inner.VerifyPreconditions(e, cfg)
			})
		},
	)
}

// GenerateAggregatorConfigAllowingCantonMinNOPs wraps generate-aggregator-config with the same bypass
// (resolveAggregatorChainSelectors validates topology during VerifyPreconditions).
func GenerateAggregatorConfigAllowingCantonMinNOPs(registry *adapters.AggregatorConfigRegistry) deployment.ChangeSetV2[v2_0_0.GenerateAggregatorConfigInput] {
	inner := v2_0_0.GenerateAggregatorConfig(registry)

	return deployment.CreateChangeSet(
		inner.Apply,
		func(e deployment.Environment, cfg v2_0_0.GenerateAggregatorConfigInput) error {
			return WithCantonProductionMinNOPCheckBypassed(cfg.Topology, func() error {
				return inner.VerifyPreconditions(e, cfg)
			})
		},
	)
}
