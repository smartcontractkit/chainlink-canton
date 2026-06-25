package changesets

import (
	"strconv"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"
)

const (
	minProductionChainNOPs       = 15
	minProductionCantonChainNOPs = 9
)

type chainCommitteeConfigBackup struct {
	qualifier     string
	chainSelector string
	cfg           offchain.ChainCommitteeConfig
}

type executorPoolChainConfigBackup struct {
	poolName      string
	chainSelector string
	cfg           offchain.ChainExecutorPoolConfig
}

// WithCantonProductionMinNOPCheckBypassed temporarily inflates Canton-family committee and
// executor-pool chain_configs during VerifyPreconditions only. The on-disk topology and Apply
// path keep the real membership.
//
// Deprecated: generic validation now delegates to ChainFamily.ValidateNOPsTopology; Canton uses
// 9 (mainnet) / 4 (testnet) via CantonChainFamilyAdapter — not the EVM 15-NOP rule. This bypass
// remains for canary topologies with fewer NOPs than the adapter minimum. See
// docs/issues/deprecate-min-nop-bypass-wrappers.md.
func WithCantonProductionMinNOPCheckBypassed(topology *offchain.EnvironmentTopology, fn func() error) error {
	if topology == nil || topology.NOPTopology == nil {
		return fn()
	}

	committeeBackups := inflateCantonCommitteeNOPsForProductionValidation(topology.NOPTopology)
	defer restoreCommitteeChainConfigBackups(topology.NOPTopology, committeeBackups)

	executorPoolBackups := inflateCantonExecutorPoolNOPsForProductionValidation(topology.ExecutorPools)
	defer restoreExecutorPoolChainConfigBackups(topology.ExecutorPools, executorPoolBackups)

	return fn()
}

func inflateCantonCommitteeNOPsForProductionValidation(n *offchain.NOPTopology) []chainCommitteeConfigBackup {
	var backups []chainCommitteeConfigBackup

	for qualifier, committee := range n.Committees {
		filler := largestCommitteeChainConfig(committee)
		if filler == nil {
			continue
		}

		for chainSelector, chainCfg := range committee.ChainConfigs {
			if !isCantonCommitteeChainSelector(chainSelector) {
				continue
			}
			if uniqueTopologyNOPCount(chainCfg.NOPAliases) >= minProductionCantonChainNOPs {
				continue
			}

			backups = append(backups, chainCommitteeConfigBackup{
				qualifier:     qualifier,
				chainSelector: chainSelector,
				cfg:           chainCfg,
			})

			patched := chainCfg
			patched.NOPAliases = append([]string(nil), filler.NOPAliases...)
			committee.ChainConfigs[chainSelector] = patched
		}

		n.Committees[qualifier] = committee
	}

	return backups
}

func restoreCommitteeChainConfigBackups(n *offchain.NOPTopology, backups []chainCommitteeConfigBackup) {
	for _, backup := range backups {
		committee, ok := n.Committees[backup.qualifier]
		if !ok {
			continue
		}
		committee.ChainConfigs[backup.chainSelector] = backup.cfg
		n.Committees[backup.qualifier] = committee
	}
}

func largestCommitteeChainConfig(committee offchain.CommitteeConfig) *offchain.ChainCommitteeConfig {
	var best *offchain.ChainCommitteeConfig
	bestCount := 0

	for chainSelector, chainCfg := range committee.ChainConfigs {
		if isCantonCommitteeChainSelector(chainSelector) {
			continue
		}
		count := uniqueTopologyNOPCount(chainCfg.NOPAliases)
		if count > bestCount {
			cfg := chainCfg
			best = &cfg
			bestCount = count
		}
	}

	return best
}

func isCantonCommitteeChainSelector(chainSelector string) bool {
	return isCantonChainSelector(chainSelector)
}

func isCantonChainSelector(chainSelector string) bool {
	sel, err := strconv.ParseUint(chainSelector, 10, 64)
	if err != nil {
		return false
	}
	family, err := chainsel.GetSelectorFamily(sel)
	if err != nil {
		return false
	}

	return family == chainsel.FamilyCanton
}

func inflateCantonExecutorPoolNOPsForProductionValidation(pools map[string]offchain.ExecutorPoolConfig) []executorPoolChainConfigBackup {
	var backups []executorPoolChainConfigBackup

	for poolName, pool := range pools {
		filler := largestExecutorPoolChainConfig(pool)
		if filler == nil {
			continue
		}

		for chainSelector, chainCfg := range pool.ChainConfigs {
			if !isCantonChainSelector(chainSelector) {
				continue
			}
			if uniqueTopologyNOPCount(chainCfg.NOPAliases) >= minProductionChainNOPs {
				continue
			}

			backups = append(backups, executorPoolChainConfigBackup{
				poolName:      poolName,
				chainSelector: chainSelector,
				cfg:           chainCfg,
			})

			patched := chainCfg
			patched.NOPAliases = append([]string(nil), filler.NOPAliases...)
			pool.ChainConfigs[chainSelector] = patched
		}

		pools[poolName] = pool
	}

	return backups
}

func restoreExecutorPoolChainConfigBackups(pools map[string]offchain.ExecutorPoolConfig, backups []executorPoolChainConfigBackup) {
	for _, backup := range backups {
		pool, ok := pools[backup.poolName]
		if !ok {
			continue
		}
		pool.ChainConfigs[backup.chainSelector] = backup.cfg
		pools[backup.poolName] = pool
	}
}

func largestExecutorPoolChainConfig(pool offchain.ExecutorPoolConfig) *offchain.ChainExecutorPoolConfig {
	var best *offchain.ChainExecutorPoolConfig
	bestCount := 0

	for chainSelector, chainCfg := range pool.ChainConfigs {
		if isCantonChainSelector(chainSelector) {
			continue
		}
		count := uniqueTopologyNOPCount(chainCfg.NOPAliases)
		if count > bestCount {
			cfg := chainCfg
			best = &cfg
			bestCount = count
		}
	}

	return best
}

func uniqueTopologyNOPCount(aliases []string) int {
	seen := make(map[string]struct{}, len(aliases))
	for _, alias := range aliases {
		seen[alias] = struct{}{}
	}

	return len(seen)
}
