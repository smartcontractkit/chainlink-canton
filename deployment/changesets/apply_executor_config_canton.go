package changesets

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/BurntSushi/toml"

	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	v2_0_0 "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/changesets"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain/sequences"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/offchain/shared"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

const executorJobTypeLabel = "executor"

// applyExecutorConfigSkippingChainSupportValidation mirrors chainlink-ccip apply-executor-config
// Apply but omits JD ListNodeChainConfigs validation so Canton destination blocks can be pushed
// to existing EVM ccvexecutor jobs before Canton node chain configs exist in JD.
func applyExecutorConfigSkippingChainSupportValidation(
	registry *adapters.ExecutorConfigRegistry,
	e deployment.Environment,
	cfg v2_0_0.ApplyExecutorConfigInput,
) (deployment.ChangesetOutput, error) {
	selectors := registry.AllDeployedChains(e.DataStore, cfg.ExecutorQualifier)
	pool := cfg.Topology.ExecutorPools[cfg.ExecutorQualifier]

	if len(selectors) == 0 {
		if !cfg.RevokeOrphanedJobs {
			e.Logger.Infow("No deployed chains found for executor pool, nothing to do",
				"qualifier", cfg.ExecutorQualifier)
			ds := datastore.NewMemoryDataStore()
			if e.DataStore != nil {
				if err := ds.Merge(e.DataStore); err != nil {
					return deployment.ChangesetOutput{}, fmt.Errorf("failed to merge datastore: %w", err)
				}
			}

			return deployment.ChangesetOutput{DataStore: ds}, nil
		}
		e.Logger.Infow("No deployed chains for executor pool, running orphan cleanup only",
			"qualifier", cfg.ExecutorQualifier)
		nopModes := executorConfigBuildNOPModes(cfg.Topology.NOPTopology.NOPs)
		scope := shared.ExecutorJobScope{ExecutorQualifier: cfg.ExecutorQualifier}
		manageReport, err := operations.ExecuteSequence(
			e.OperationsBundle,
			sequences.ManageJobProposals,
			sequences.ManageJobProposalsDeps{Env: e},
			sequences.ManageJobProposalsInput{
				JobSpecs:           nil,
				AffectedScope:      scope,
				Labels:             map[string]string{"job_type": executorJobTypeLabel, executorJobTypeLabel: cfg.ExecutorQualifier},
				NOPs:               sequences.NOPContext{Modes: nopModes, TargetNOPs: cfg.TargetNOPs, AllNOPs: executorConfigGetAllNOPAliases(cfg.Topology.NOPTopology.NOPs)},
				RevokeOrphanedJobs: true,
			},
		)
		if err != nil {
			return deployment.ChangesetOutput{Reports: manageReport.ExecutionReports}, fmt.Errorf("failed to manage job proposals (orphan cleanup): %w", err)
		}

		return deployment.ChangesetOutput{Reports: manageReport.ExecutionReports, DataStore: manageReport.Output.DataStore}, nil
	}

	chainConfigs, err := executorConfigBuildChainConfigs(registry, e.DataStore, selectors, cfg.ExecutorQualifier)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	monitoring := executorConfigConvertTopologyMonitoring(&cfg.Topology.Monitoring)
	nopModes := executorConfigBuildNOPModes(cfg.Topology.NOPTopology.NOPs)

	jobSpecs, scope, err := executorConfigBuildJobSpecs(
		chainConfigs,
		cfg.ExecutorQualifier,
		cfg.TargetNOPs,
		pool,
		cfg.Topology.IndexerAddress,
		cfg.Topology.PyroscopeURL,
		monitoring,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	manageReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.ManageJobProposals,
		sequences.ManageJobProposalsDeps{Env: e},
		sequences.ManageJobProposalsInput{
			JobSpecs:      jobSpecs,
			AffectedScope: scope,
			Labels: map[string]string{
				"job_type":           executorJobTypeLabel,
				executorJobTypeLabel: cfg.ExecutorQualifier,
			},
			NOPs: sequences.NOPContext{
				Modes:      nopModes,
				TargetNOPs: cfg.TargetNOPs,
				AllNOPs:    executorConfigGetAllNOPAliases(cfg.Topology.NOPTopology.NOPs),
			},
			RevokeOrphanedJobs: cfg.RevokeOrphanedJobs,
		},
	)
	if err != nil {
		return deployment.ChangesetOutput{
			Reports: manageReport.ExecutionReports,
		}, fmt.Errorf("failed to manage job proposals: %w", err)
	}

	e.Logger.Infow("Executor config applied",
		"jobsCount", len(manageReport.Output.Jobs),
		"revokedCount", len(manageReport.Output.RevokedJobs))

	return deployment.ChangesetOutput{
		Reports:   manageReport.ExecutionReports,
		DataStore: manageReport.Output.DataStore,
	}, nil
}

func executorConfigBuildChainConfigs(
	registry *adapters.ExecutorConfigRegistry,
	ds datastore.DataStore,
	selectors []uint64,
	qualifier string,
) (map[string]adapters.ExecutorChainConfig, error) {
	chainConfigs := make(map[string]adapters.ExecutorChainConfig, len(selectors))
	for _, sel := range selectors {
		adapter, err := registry.GetByChain(sel)
		if err != nil {
			return nil, fmt.Errorf("no adapter for chain %d: %w", sel, err)
		}
		cfg, err := adapter.BuildChainConfig(ds, sel, qualifier)
		if err != nil {
			return nil, fmt.Errorf("failed to build config for chain %d: %w", sel, err)
		}
		chainConfigs[strconv.FormatUint(sel, 10)] = cfg
	}

	return chainConfigs, nil
}

func executorConfigBuildJobSpecs(
	chainConfigs map[string]adapters.ExecutorChainConfig,
	executorQualifier string,
	targetNOPs []shared.NOPAlias,
	pool offchain.ExecutorPoolConfig,
	indexerAddress []string,
	pyroscopeURL string,
	monitoring shared.MonitoringInput,
) (shared.NOPJobSpecs, shared.ExecutorJobScope, error) {
	scope := shared.ExecutorJobScope{
		ExecutorQualifier: executorQualifier,
	}

	poolNOPs := executorConfigGetPoolNOPAliases(pool)
	nopAliases := targetNOPs
	if len(nopAliases) == 0 {
		nopAliases = poolNOPs
	}

	jobSpecs := make(shared.NOPJobSpecs)

	for _, nopAlias := range nopAliases {
		chainCfgs := make(map[string]offchain.ExecutorChainCfg)
		for chainSelectorStr, genCfg := range chainConfigs {
			chainCfg, ok := pool.ChainConfigs[chainSelectorStr]
			if !ok {
				continue
			}
			sortedPool := slices.Clone(chainCfg.NOPAliases)
			slices.Sort(sortedPool)
			chainCfgs[chainSelectorStr] = offchain.ExecutorChainCfg{
				OffRampAddress:         genCfg.OffRampAddress,
				RmnAddress:             genCfg.RmnAddress,
				DefaultExecutorAddress: genCfg.ExecutorProxyAddress,
				ExecutorPool:           sortedPool,
				ExecutionInterval:      chainCfg.ExecutionInterval,
			}
		}

		jobSpecID := shared.NewExecutorJobID(nopAlias, scope)

		executorCfg := offchain.ExecutorConfiguration{
			IndexerAddress:    indexerAddress,
			ExecutorID:        jobSpecID.GetExecutorID(),
			PyroscopeURL:      pyroscopeURL,
			NtpServer:         pool.NtpServer,
			IndexerQueryLimit: pool.IndexerQueryLimit,
			BackoffDuration:   pool.BackoffDuration,
			LookbackWindow:    pool.LookbackWindow,
			ReaderCacheExpiry: pool.ReaderCacheExpiry,
			MaxRetryDuration:  pool.MaxRetryDuration,
			WorkerCount:       pool.WorkerCount,
			Monitoring: offchain.ExecutorMonitoringConfig{
				Enabled: monitoring.Enabled,
				Type:    monitoring.Type,
				Beholder: offchain.ExecutorBeholderConfig{
					InsecureConnection:       monitoring.Beholder.InsecureConnection,
					CACertFile:               monitoring.Beholder.CACertFile,
					OtelExporterGRPCEndpoint: monitoring.Beholder.OtelExporterGRPCEndpoint,
					OtelExporterHTTPEndpoint: monitoring.Beholder.OtelExporterHTTPEndpoint,
					LogStreamingEnabled:      monitoring.Beholder.LogStreamingEnabled,
					MetricReaderInterval:     monitoring.Beholder.MetricReaderInterval,
					TraceSampleRatio:         monitoring.Beholder.TraceSampleRatio,
					TraceBatchTimeout:        monitoring.Beholder.TraceBatchTimeout,
				},
			},
			ChainConfiguration: chainCfgs,
		}

		configBytes, err := toml.Marshal(executorCfg)
		if err != nil {
			return nil, scope, fmt.Errorf("failed to marshal executor config for NOP %q: %w", nopAlias, err)
		}

		jobID := jobSpecID.ToJobID()
		jobSpec := fmt.Sprintf(`schemaVersion = 1
type = "ccvexecutor"
name = "%s"
externalJobID = "%s"
forwardingAllowed = false
executorConfig = '''
%s'''
`, string(jobID), jobID.ToExternalJobID(), string(configBytes))

		if jobSpecs[nopAlias] == nil {
			jobSpecs[nopAlias] = make(map[shared.JobID]string)
		}
		jobSpecs[nopAlias][jobID] = jobSpec
	}

	return jobSpecs, scope, nil
}

func executorConfigConvertTopologyMonitoring(m *offchain.MonitoringConfig) shared.MonitoringInput {
	if m == nil {
		return shared.MonitoringInput{}
	}

	return shared.MonitoringInput{
		Enabled: m.Enabled,
		Type:    m.Type,
		Beholder: shared.BeholderInput{
			InsecureConnection:       m.Beholder.InsecureConnection,
			CACertFile:               m.Beholder.CACertFile,
			OtelExporterGRPCEndpoint: m.Beholder.OtelExporterGRPCEndpoint,
			OtelExporterHTTPEndpoint: m.Beholder.OtelExporterHTTPEndpoint,
			LogStreamingEnabled:      m.Beholder.LogStreamingEnabled,
			MetricReaderInterval:     m.Beholder.MetricReaderInterval,
			TraceSampleRatio:         m.Beholder.TraceSampleRatio,
			TraceBatchTimeout:        m.Beholder.TraceBatchTimeout,
		},
	}
}

func executorConfigBuildNOPModes(nops []offchain.NOPConfig) map[shared.NOPAlias]shared.NOPMode {
	nopModes := make(map[shared.NOPAlias]shared.NOPMode)
	for _, nop := range nops {
		nopModes[shared.NOPAlias(nop.Alias)] = nop.GetMode()
	}

	return nopModes
}

func executorConfigGetPoolNOPAliases(pool offchain.ExecutorPoolConfig) []shared.NOPAlias {
	aliasSet := make(map[string]struct{})
	for _, chainCfg := range pool.ChainConfigs {
		for _, alias := range chainCfg.NOPAliases {
			aliasSet[alias] = struct{}{}
		}
	}
	aliases := make([]string, 0, len(aliasSet))
	for a := range aliasSet {
		aliases = append(aliases, a)
	}
	slices.Sort(aliases)

	return shared.ConvertStringToNopAliases(aliases)
}

func executorConfigGetAllNOPAliases(nops []offchain.NOPConfig) []shared.NOPAlias {
	aliases := make([]shared.NOPAlias, len(nops))
	for i, nop := range nops {
		aliases[i] = shared.NOPAlias(nop.Alias)
	}

	return aliases
}
