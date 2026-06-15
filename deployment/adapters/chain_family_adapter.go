package adapters

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	committeeverifierop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	executorop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutil "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

var _ ccipadapters.ChainFamily = (*CantonChainFamilyAdapter)(nil)

// ConfigureLanesDataStoreFamilyExtraKey passes e.DataStore into Canton lane configure when the
// pinned chainlink-ccip ConfigureChainForLanesInput type does not yet expose DataStore.
const ConfigureLanesDataStoreFamilyExtraKey = "canton.dataStore"

type CantonChainFamilyAdapter struct{}

var CantonFamilySelector = [4]byte{0xdf, 0xaf, 0xaf, 0x4b}

func DefaultCantonFeeQuoterDestChainConfig() lanes.FeeQuoterDestChainConfig {
	return lanes.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      false,
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         binary.BigEndian.Uint32(CantonFamilySelector[:]),
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V1Params:                    nil,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 90,
			USDPerUnitGas:            big.NewInt(38),
		},
	}
}

func (a *CantonChainFamilyAdapter) ConfigureChainForLanes() *cldfops.Sequence[ccipadapters.ConfigureChainForLanesInput, ccipseq.OnChainOutput, cldfchain.BlockChains] {
	return cldfops.NewSequence(
		"canton/configure-chain-for-lanes",
		semver.MustParse("2.0.0"),
		"Configures CCIP lanes for a Canton chain",
		func(b cldfops.Bundle, chains cldfchain.BlockChains, input ccipadapters.ConfigureChainForLanesInput) (ccipseq.OnChainOutput, error) {
			out, err := a.configureChainForLanes(b, chains, input)
			if err != nil {
				return ccipseq.OnChainOutput{}, err
			}
			out.BatchOps = cantonmcms.ConsolidateBatchOpsPerChain(out.BatchOps)

			return out, nil
		},
	)
}

// ConfigureCommitteeVerifierForLanes configures CommitteeVerifier lane settings only.
// Used by ConfigureCantonCommitteeVerifierForLanesFromTopology (Run 2).
func (a *CantonChainFamilyAdapter) ConfigureCommitteeVerifierForLanes() *cldfops.Sequence[ccipadapters.ConfigureChainForLanesInput, ccipseq.OnChainOutput, cldfchain.BlockChains] {
	return cldfops.NewSequence(
		"canton/configure-committee-verifier-for-lanes",
		semver.MustParse("2.0.0"),
		"Configures CommitteeVerifier lane settings for a Canton chain",
		func(b cldfops.Bundle, chains cldfchain.BlockChains, input ccipadapters.ConfigureChainForLanesInput) (ccipseq.OnChainOutput, error) {
			out, err := a.configureCommitteeVerifierForLanes(b, chains, input)
			if err != nil {
				return ccipseq.OnChainOutput{}, err
			}
			out.BatchOps = cantonmcms.ConsolidateBatchOpsPerChain(out.BatchOps)

			return out, nil
		},
	)
}

// configureChainForLanes resolves deployed contract refs from the input datastore and runs
// lane configure sequences (GlobalConfig, FeeQuoter, Executor, CommitteeVerifier, …).
// MCMS transactions are returned in OnChainOutput.BatchOps for the changeset layer to split
// into ccipOwner and ccvOwner proposals.
func (a *CantonChainFamilyAdapter) configureChainForLanes(
	b cldfops.Bundle,
	chains cldfchain.BlockChains,
	input ccipadapters.ConfigureChainForLanesInput,
) (ccipseq.OnChainOutput, error) {
	ds, err := dataStoreFromConfigureChainForLanesInput(input)
	if err != nil {
		return ccipseq.OnChainOutput{}, err
	}

	localGlobalConfig, err := findContractRef(
		ds,
		input.ChainSelector,
		datastore.ContractType(global_config.ContractType),
		global_config.Version,
		"",
	)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve global config: %w", err)
	}

	router, onRamp, feeQuoter, offRamp, err := a.resolveLocalContractsForConfigureLanes(ds, input)
	if err != nil {
		return ccipseq.OnChainOutput{}, err
	}

	localCommitteeVerifiers := convertCommitteeVerifierConfigs(input.CommitteeVerifiers)
	var out ccipseq.OnChainOutput

	chain, ok := chains.CantonChains()[input.ChainSelector]
	if !ok || len(chain.Participants) == 0 {
		return ccipseq.OnChainOutput{}, fmt.Errorf("canton chain %d not found or has no participants", input.ChainSelector)
	}
	nativeInstrument, err := lookupNativeInstrumentID(b.GetContext(), chain.Participants[0])
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve Canton native fee token instrument: %w", err)
	}

	for remoteSelector, remoteCfg := range input.RemoteChains {
		localExecutor, err := resolveContractRefByAddress(
			ds,
			input.ChainSelector,
			datastore.ContractType(executorop.ContractType),
			executorop.Version,
			remoteCfg.DefaultExecutor,
		)
		if err != nil {
			return out, fmt.Errorf("resolve executor for remote chain %d: %w", remoteSelector, err)
		}
		defaultInboundCCVs, err := resolveContractRefsByAddresses(
			ds,
			input.ChainSelector,
			datastore.ContractType(committeeverifierop.ContractType),
			committeeverifierop.Version,
			remoteCfg.DefaultInboundCCVs,
		)
		if err != nil {
			return out, fmt.Errorf("resolve default inbound ccvs for remote chain %d: %w", remoteSelector, err)
		}
		laneMandatedInboundCCVs, err := resolveContractRefsByAddresses(
			ds,
			input.ChainSelector,
			datastore.ContractType(committeeverifierop.ContractType),
			committeeverifierop.Version,
			remoteCfg.LaneMandatedInboundCCVs,
		)
		if err != nil {
			return out, fmt.Errorf("resolve lane mandated inbound ccvs for remote chain %d: %w", remoteSelector, err)
		}
		defaultOutboundCCVs, err := resolveContractRefsByAddresses(
			ds,
			input.ChainSelector,
			datastore.ContractType(committeeverifierop.ContractType),
			committeeverifierop.Version,
			remoteCfg.DefaultOutboundCCVs,
		)
		if err != nil {
			return out, fmt.Errorf("resolve default outbound ccvs for remote chain %d: %w", remoteSelector, err)
		}
		laneMandatedOutboundCCVs, err := resolveContractRefsByAddresses(
			ds,
			input.ChainSelector,
			datastore.ContractType(committeeverifierop.ContractType),
			committeeverifierop.Version,
			remoteCfg.LaneMandatedOutboundCCVs,
		)
		if err != nil {
			return out, fmt.Errorf("resolve lane mandated outbound ccvs for remote chain %d: %w", remoteSelector, err)
		}

		localChain := &lanes.ChainDefinition{
			Selector:                 input.ChainSelector,
			CommitteeVerifiers:       localCommitteeVerifiers,
			DefaultInboundCCVs:       defaultInboundCCVs,
			LaneMandatedInboundCCVs:  laneMandatedInboundCCVs,
			DefaultOutboundCCVs:      defaultOutboundCCVs,
			LaneMandatedOutboundCCVs: laneMandatedOutboundCCVs,
			DefaultExecutor:          localExecutor,
			CantonLaneConfig: &lanes.CantonLaneConfig{
				GlobalConfig: localGlobalConfig,
			},
			OnRamp:    onRamp,
			OffRamp:   offRamp,
			Router:    router,
			FeeQuoter: feeQuoter,
		}

		remoteChain, err := remoteChainDefinition(remoteSelector, remoteCfg)
		if err != nil {
			return out, err
		}
		tokenPrices, err := resolveTokenPricesForRemoteDest(ds, input, remoteSelector, &nativeInstrument)
		if err != nil {
			return out, fmt.Errorf("resolve token prices for remote chain %d: %w", remoteSelector, err)
		}
		remoteChain.TokenPrices = tokenPrices

		out, err = ccipseq.RunAndMergeSequence(
			b,
			chains,
			sequences.ConfigureLaneLegAsSourceWithInput,
			sequences.ConfigureLaneLegInput{
				Lane: lanes.UpdateLanesInput{
					Source: localChain,
					Dest:   remoteChain,
				},
				DataStore: ds,
			},
			out,
		)
		if err != nil {
			return out, err
		}

		out, err = ccipseq.RunAndMergeSequence(
			b,
			chains,
			sequences.ConfigureLaneLegAsDest,
			lanes.UpdateLanesInput{
				Source: remoteChain,
				Dest:   localChain,
			},
			out,
		)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func (a *CantonChainFamilyAdapter) configureCommitteeVerifierForLanes(
	b cldfops.Bundle,
	chains cldfchain.BlockChains,
	input ccipadapters.ConfigureChainForLanesInput,
) (ccipseq.OnChainOutput, error) {
	chain, ok := chains.CantonChains()[input.ChainSelector]
	if !ok {
		return ccipseq.OnChainOutput{}, fmt.Errorf("canton chain %d not found", input.ChainSelector)
	}

	mcmsEnabled := len(chain.Participants[0].ReadAsPartyIDs) > 0
	localCommitteeVerifiers := convertCommitteeVerifierConfigs(input.CommitteeVerifiers)
	var out ccipseq.OnChainOutput

	for _, verifierConfig := range localCommitteeVerifiers {
		cvReport, err := cldfops.ExecuteSequence(b, sequences.ConfigureCommitteeVerifierAsSource, chains, sequences.ConfigureCommitteeVerifierAsSourceInput{
			ChainSelector:           input.ChainSelector,
			MCMSEnabled:             mcmsEnabled,
			CommitteeVerifierConfig: verifierConfig,
		})
		if err != nil {
			return out, fmt.Errorf("configuring committee verifier as source: %w", err)
		}
		out.BatchOps = append(out.BatchOps, cvReport.Output.BatchOps...)

		cvDestReport, err := cldfops.ExecuteSequence(b, sequences.ConfigureCommitteeVerifierAsDest, chains, sequences.ConfigureCommitteeVerifierAsDestInput{
			ChainSelector:           input.ChainSelector,
			MCMSEnabled:             mcmsEnabled,
			CommitteeVerifierConfig: verifierConfig,
		})
		if err != nil {
			return out, fmt.Errorf("configuring committee verifier as dest: %w", err)
		}
		out.BatchOps = append(out.BatchOps, cvDestReport.Output.BatchOps...)
	}

	return out, nil
}

func (a *CantonChainFamilyAdapter) resolveLocalContractsForConfigureLanes(
	ds datastore.DataStore,
	input ccipadapters.ConfigureChainForLanesInput,
) (router, onRamp, feeQuoter, offRamp []byte, err error) {
	if input.AllowOnrampOverride {
		router, err = a.GetTestRouter(ds, input.ChainSelector)
	} else {
		router, err = a.GetRouterAddress(ds, input.ChainSelector)
	}
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("resolve router: %w", err)
	}

	onRamp, err = a.GetOnRampAddress(ds, input.ChainSelector)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("resolve onRamp: %w", err)
	}
	feeQuoter, err = a.GetFQAddress(ds, input.ChainSelector)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("resolve feeQuoter: %w", err)
	}
	offRamp, err = a.GetOffRampAddress(ds, input.ChainSelector)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("resolve offRamp: %w", err)
	}

	return router, onRamp, feeQuoter, offRamp, nil
}

func (a *CantonChainFamilyAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return dsutil.ToInstanceAddressBytes(ref)
}

func (a *CantonChainFamilyAdapter) GetOnRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return findContractBytes(ds, chainSelector, datastore.ContractType(onramp.ContractType), onramp.Version)
}

func (a *CantonChainFamilyAdapter) GetOffRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return findContractBytes(ds, chainSelector, datastore.ContractType(offramp.ContractType), offramp.Version)
}

func (a *CantonChainFamilyAdapter) GetFQAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return findContractBytes(ds, chainSelector, datastore.ContractType(fee_quoter.ContractType), fee_quoter.Version)
}

func (a *CantonChainFamilyAdapter) GetRouterAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return findContractBytes(ds, chainSelector, datastore.ContractType(global_config.ContractType), global_config.Version)
}

func (a *CantonChainFamilyAdapter) GetTestRouter(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return a.GetRouterAddress(ds, chainSelector)
}

func (a *CantonChainFamilyAdapter) ResolveExecutor(ds datastore.DataStore, chainSelector uint64, qualifier string) (string, error) {
	ref, err := findContractRef(ds, chainSelector, datastore.ContractType(executorop.ContractType), executorop.Version, qualifier)
	if err != nil {
		return "", err
	}

	return ref.Address, nil
}

// GetAddressBytesLength implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetAddressBytesLength() uint8 {
	return 32
}

// GetChainFamilySelector implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetChainFamilySelector() [4]byte {
	return CantonFamilySelector
}

// GetDefaultCommitteeVerifierRemoteChainConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultCommitteeVerifierRemoteChainConfig() ccipadapters.CommitteeVerifierRemoteChainDefaults {
	return ccipadapters.CommitteeVerifierRemoteChainDefaults{
		AllowlistEnabled:   false,
		FeeUSDCents:        0,
		GasForVerification: 100_000,
		PayloadSizeBytes:   1_000,
	}
}

// GetDefaultFeeQuoterDestChainConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultFeeQuoterDestChainConfig(_, _ uint64, chainFamilySelector [4]byte) ccipadapters.FeeQuoterDestChainConfigOverrides {
	defaults := DefaultCantonFeeQuoterDestChainConfig()
	linkFeeMultiplier := uint8(0)
	if defaults.V2Params != nil {
		linkFeeMultiplier = defaults.V2Params.LinkFeeMultiplierPercent
	}
	usdPerUnitGas := big.NewInt(0)
	if defaults.V2Params != nil && defaults.V2Params.USDPerUnitGas != nil {
		usdPerUnitGas = defaults.V2Params.USDPerUnitGas
	}

	return ccipadapters.FeeQuoterDestChainConfigOverrides{
		IsEnabled:                   new(defaults.IsEnabled),
		MaxDataBytes:                new(defaults.MaxDataBytes),
		MaxPerMsgGasLimit:           new(defaults.MaxPerMsgGasLimit),
		DestGasOverhead:             new(defaults.DestGasOverhead),
		DestGasPerPayloadByteBase:   new(defaults.DestGasPerPayloadByteBase),
		ChainFamilySelector:         chainFamilySelector,
		DefaultTokenFeeUSDCents:     new(defaults.DefaultTokenFeeUSDCents),
		DefaultTokenDestGasOverhead: new(defaults.DefaultTokenDestGasOverhead),
		DefaultTxGasLimit:           new(defaults.DefaultTxGasLimit),
		NetworkFeeUSDCents:          new(defaults.NetworkFeeUSDCents),
		LinkFeeMultiplierPercent:    new(linkFeeMultiplier),
		USDPerUnitGas:               usdPerUnitGas,
	}
}

// GetDefaultFinalityConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultFinalityConfig() finality.Config {
	return finality.Config{
		WaitForFinality: true,
	}
}

// GetDefaultRemoteChainConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultRemoteChainConfig(_, _ uint64) ccipadapters.RemoteChainDefaults {
	return ccipadapters.RemoteChainDefaults{
		AllowTrafficFrom: true, // TODO: check what this does?
		ExecutorDestChainConfig: ccipadapters.ExecutorDestChainConfig{
			USDCentsFee: 0,
			Enabled:     true,
		},
		BaseExecutionGasCost:      50_000,
		TokenReceiverAllowed:      true,
		MessageNetworkFeeUSDCents: 0,
		TokenNetworkFeeUSDCents:   0,
	}
}

func findContractRef(ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType, version *semver.Version, qualifier string) (datastore.AddressRef, error) {
	return ds.Addresses().Get(datastore.NewAddressRefKey(chainSelector, contractType, version, qualifier))
}

func findContractBytes(ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType, version *semver.Version) ([]byte, error) {
	ref, err := findContractRef(ds, chainSelector, contractType, version, "")
	if err != nil {
		return nil, err
	}

	return dsutil.ToInstanceAddressBytes(ref)
}

func resolveContractRefByAddress(ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType, version *semver.Version, address string) (datastore.AddressRef, error) {
	refs, err := ds.Addresses().Fetch()
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("fetch addresses: %w", err)
	}
	for _, ref := range refs {
		if ref.ChainSelector == chainSelector && ref.Type == contractType && ref.Version.Equal(version) && ref.Address == address {
			return ref, nil
		}
	}

	return datastore.AddressRef{}, fmt.Errorf("no %s address ref found for %s on chain %d", contractType, address, chainSelector)
}

func resolveContractRefsByAddresses(ds datastore.DataStore, chainSelector uint64, contractType datastore.ContractType, version *semver.Version, addresses []string) ([]datastore.AddressRef, error) {
	refs := make([]datastore.AddressRef, 0, len(addresses))
	for _, address := range addresses {
		ref, err := resolveContractRefByAddress(ds, chainSelector, contractType, version, address)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}

	return refs, nil
}

func convertCommitteeVerifierConfigs(configs []ccipadapters.CommitteeVerifierConfig[datastore.AddressRef]) []lanes.CommitteeVerifierConfig[datastore.AddressRef] {
	out := make([]lanes.CommitteeVerifierConfig[datastore.AddressRef], 0, len(configs))
	for _, cfg := range configs {
		remoteChains := make(map[uint64]lanes.CommitteeVerifierRemoteChainConfig, len(cfg.RemoteChains))
		for selector, remoteCfg := range cfg.RemoteChains {
			remoteChains[selector] = lanes.CommitteeVerifierRemoteChainConfig{
				AllowlistEnabled:          remoteCfg.AllowlistEnabled,
				AddedAllowlistedSenders:   remoteCfg.AddedAllowlistedSenders,
				RemovedAllowlistedSenders: remoteCfg.RemovedAllowlistedSenders,
				FeeUSDCents:               remoteCfg.FeeUSDCents,
				GasForVerification:        remoteCfg.GasForVerification,
				PayloadSizeBytes:          remoteCfg.PayloadSizeBytes,
				SignatureConfig: lanes.CommitteeVerifierSignatureQuorumConfig{
					Signers:   remoteCfg.SignatureConfig.Signers,
					Threshold: remoteCfg.SignatureConfig.Threshold,
				},
			}
		}
		out = append(out, lanes.CommitteeVerifierConfig[datastore.AddressRef]{
			CommitteeVerifier: cfg.CommitteeVerifier,
			RemoteChains:      remoteChains,
		})
	}

	return out
}

func remoteChainDefinition(remoteSelector uint64, remoteCfg ccipadapters.RemoteChainConfig[[]byte, string]) (*lanes.ChainDefinition, error) {
	if len(remoteCfg.OnRamps) == 0 {
		return nil, fmt.Errorf("remote chain %d has no onramp address", remoteSelector)
	}

	router := remoteCfg.OffRamp
	if len(router) == 0 {
		router = remoteCfg.OnRamps[0]
	}

	tokenReceiverAllowed := new(false)
	if remoteCfg.TokenReceiverAllowed != nil {
		tokenReceiverAllowed = remoteCfg.TokenReceiverAllowed
	}

	return &lanes.ChainDefinition{
		Selector:                  remoteSelector,
		ExecutorDestChainConfig:   lanes.ExecutorDestChainConfig(remoteCfg.ExecutorDestChainConfig),
		AddressBytesLength:        remoteCfg.AddressBytesLength,
		BaseExecutionGasCost:      remoteCfg.BaseExecutionGasCost,
		TokenReceiverAllowed:      tokenReceiverAllowed,
		MessageNetworkFeeUSDCents: remoteCfg.MessageNetworkFeeUSDCents,
		TokenNetworkFeeUSDCents:   remoteCfg.TokenNetworkFeeUSDCents,
		OnRamp:                    remoteCfg.OnRamps[0],
		OffRamp:                   remoteCfg.OffRamp,
		Router:                    router,
		FeeQuoterDestChainConfig:  feeQuoterDestChainConfigFromOverrides(remoteCfg.FeeQuoterDestChainConfig),
	}, nil
}

func feeQuoterDestChainConfigFromOverrides(cfg ccipadapters.FeeQuoterDestChainConfigOverrides) lanes.FeeQuoterDestChainConfig {
	out := DefaultCantonFeeQuoterDestChainConfig()
	out.OverrideExistingConfig = cfg.OverrideExistingConfig
	if cfg.IsEnabled != nil {
		out.IsEnabled = *cfg.IsEnabled
	}
	if cfg.MaxDataBytes != nil {
		out.MaxDataBytes = *cfg.MaxDataBytes
	}
	if cfg.MaxPerMsgGasLimit != nil {
		out.MaxPerMsgGasLimit = *cfg.MaxPerMsgGasLimit
	}
	if cfg.DestGasOverhead != nil {
		out.DestGasOverhead = *cfg.DestGasOverhead
	}
	if cfg.DestGasPerPayloadByteBase != nil {
		out.DestGasPerPayloadByteBase = *cfg.DestGasPerPayloadByteBase
	}
	if cfg.ChainFamilySelector != [4]byte{} {
		out.ChainFamilySelector = binary.BigEndian.Uint32(cfg.ChainFamilySelector[:])
	}
	if cfg.DefaultTokenFeeUSDCents != nil {
		out.DefaultTokenFeeUSDCents = *cfg.DefaultTokenFeeUSDCents
	}
	if cfg.DefaultTokenDestGasOverhead != nil {
		out.DefaultTokenDestGasOverhead = *cfg.DefaultTokenDestGasOverhead
	}
	if cfg.DefaultTxGasLimit != nil {
		out.DefaultTxGasLimit = *cfg.DefaultTxGasLimit
	}
	if cfg.NetworkFeeUSDCents != nil {
		out.NetworkFeeUSDCents = *cfg.NetworkFeeUSDCents
	}
	if out.V2Params == nil {
		out.V2Params = &lanes.FeeQuoterV2Params{}
	}
	if cfg.LinkFeeMultiplierPercent != nil {
		out.V2Params.LinkFeeMultiplierPercent = *cfg.LinkFeeMultiplierPercent
	}
	if cfg.USDPerUnitGas != nil {
		out.V2Params.USDPerUnitGas = cfg.USDPerUnitGas
	}

	return out
}

// minProductionCantonChainNOPs is the minimum unique NOP count for Canton chains in production.
const minProductionCantonChainNOPs = 9

func (a *CantonChainFamilyAdapter) ValidateNOPsTopology(chainSelector string, nopCount int) error {
	if nopCount < minProductionCantonChainNOPs {
		return fmt.Errorf(
			"chain %q requires at least %d unique NOPs for production environments, got %d",
			chainSelector,
			minProductionCantonChainNOPs,
			nopCount,
		)
	}

	return nil
}

func dataStoreFromConfigureChainForLanesInput(input ccipadapters.ConfigureChainForLanesInput) (datastore.DataStore, error) {
	if input.FamilyExtras != nil {
		if ds, ok := input.FamilyExtras[ConfigureLanesDataStoreFamilyExtraKey].(datastore.DataStore); ok && ds != nil {
			return ds, nil
		}
	}

	if ds, ok := cachedConfigureLanesDataStore(input.ChainSelector); ok {
		return ds, nil
	}

	v := reflect.ValueOf(input)
	if f := v.FieldByName("DataStore"); f.IsValid() && !f.IsZero() {
		if ds, ok := f.Interface().(datastore.DataStore); ok && ds != nil {
			return ds, nil
		}
	}

	return nil, fmt.Errorf("datastore is required (set FamilyExtras[%q] or upgrade chainlink-ccip/deployment)", ConfigureLanesDataStoreFamilyExtraKey)
}
