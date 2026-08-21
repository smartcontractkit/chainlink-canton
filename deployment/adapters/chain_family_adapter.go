package adapters

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/finality"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	ccipseq "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	committeeverifierop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	executorop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
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

func isEVMSelector(chainSelector uint64) bool {
	family, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil {
		return false
	}

	return family == chainsel.FamilyEVM
}

func DefaultCantonFeeQuoterDestChainConfig() lanes.FeeQuoterDestChainConfig {
	return lanes.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      false,
		IsEnabled:                   true,
		MaxDataBytes:                32_000,
		MaxPerMsgGasLimit:           15_000_000,
		DestGasOverhead:             0,
		DestGasPerPayloadByteBase:   20,
		ChainFamilySelector:         binary.BigEndian.Uint32(CantonFamilySelector[:]),
		DefaultTokenFeeUSDCents:     0,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          50,
		V1Params:                    nil,
		V2Params: &lanes.FeeQuoterV2Params{
			LinkFeeMultiplierPercent: 90,
			USDPerUnitGas:            big.NewInt(38),
		},
	}
}

func defaultCantonFeeQuoterDestChainConfig() lanes.FeeQuoterDestChainConfig {
	return lanes.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      false,
		IsEnabled:                   true,
		MaxDataBytes:                32_000,
		MaxPerMsgGasLimit:           15_000_000,
		DestGasOverhead:             0,
		DestGasPerPayloadByteBase:   20,
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
	nativeInstrument, err := lookupNativeInstrumentID(b.GetContext(), chain.Participants[0], ds, input.ChainSelector)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve Canton native fee token instrument: %w", err)
	}

	// Register native fee token (Amulet) TokenConfig in TAR once per lane configure run.
	// Skipped when already on-ledger. Not inlined in proposal-driven core deploy because
	// timelock-execute pre-resolves TAR before the deploy batch creates it.
	tarRef, err := findContractRef(
		ds,
		input.ChainSelector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
	}
	tarRaw, err := dsutil.GetRawInstanceAddressFromAddressRef(tarRef)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve token admin registry raw address: %w", err)
	}
	ccipOwnerParty, err := resolveCcipOwnerParty(ds, input.ChainSelector)
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("resolve ccipOwner party: %w", err)
	}
	mcmsEnabled := len(chain.Participants[0].ReadAsPartyIDs) > 0
	tarReport, err := cldfops.ExecuteSequence(b, sequences.RegisterNativeFeeTokenInTAR, chain, sequences.RegisterNativeFeeTokenInTARInput{
		TokenAdminRegistryInstanceAddress:    contracts.HexToInstanceAddress(tarRef.Address),
		TokenAdminRegistryRawInstanceAddress: tarRaw,
		InstrumentId:                         nativeInstrument,
		CcipOwnerParty:                       ccipOwnerParty,
		TokenQualifier:                       string(nativeInstrument.Id),
		ChainSelector:                        input.ChainSelector,
		ProposalDriven:                       mcmsEnabled,
	})
	if err != nil {
		return ccipseq.OnChainOutput{}, fmt.Errorf("register native fee token in TAR: %w", err)
	}
	out.BatchOps = append(out.BatchOps, tarReport.Output.BatchOps...)
	out.Addresses = append(out.Addresses, tarReport.Output.Addresses...)

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
		defaultInboundCCVs, err := resolveDefaultInboundCCVs(input.ChainSelector, remoteCfg)
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

		remoteChain, err := BuildRemoteChainDefinition(remoteSelector, remoteCfg)
		if err != nil {
			return out, err
		}
		tokenPrices, err := ResolveTokenPricesForRemoteDest(ds, input, remoteSelector, &nativeInstrument)
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
//
// The lane changeset asks the *remote* chain's adapter how the local chain should configure
// its FeeQuoter for that remote as a destination, so this is reached with Canton as the
// destination and sourceChainSelector naming the sending chain.
func (a *CantonChainFamilyAdapter) GetDefaultFeeQuoterDestChainConfig(sourceChainSelector, _ uint64, chainFamilySelector [4]byte) ccipadapters.FeeQuoterDestChainConfigOverrides {
	if isEVMSelector(sourceChainSelector) {
		// Canton has no gas market and execution is manual, so every gas-derived field is
		// zeroed and the lane is priced by a flat network fee instead. USDPerUnitGas is still
		// required: the EVM FeeQuoter 2.0 unconditionally reverts NoGasPriceAvailable in
		// quoteGasForExec for any enabled destination that has no gas-price entry, regardless of
		// family, so an EVM source must seed a Canton gas price via FeeQuoter::UpdatePrices even
		// though Canton is priced flat. The value mirrors the non-EVM branch
		// (defaultCantonFeeQuoterDestChainConfig().V2Params.USDPerUnitGas = 38, i.e. 0.0000000038
		// USD/gas) and the gas price the Canton-side configure sequence pushes for its EVM remote.
		return ccipadapters.FeeQuoterDestChainConfigOverrides{
			IsEnabled:                   new(true),
			MaxDataBytes:                new(uint32(32_000)),
			MaxPerMsgGasLimit:           new(uint32(15_000_000)),
			DestGasOverhead:             new(uint32(0)),
			DestGasPerPayloadByteBase:   new(uint8(20)),
			ChainFamilySelector:         chainFamilySelector,
			DefaultTokenFeeUSDCents:     new(uint16(0)),
			DefaultTokenDestGasOverhead: new(uint32(0)),
			DefaultTxGasLimit:           new(uint32(1)),
			NetworkFeeUSDCents:          new(uint16(50)),
			LinkFeeMultiplierPercent:    new(uint8(90)),
			USDPerUnitGas:               defaultCantonFeeQuoterDestChainConfig().V2Params.USDPerUnitGas,
		}
	}

	defaults := defaultCantonFeeQuoterDestChainConfig()
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
//
// The lane changeset asks the *remote* chain's adapter for the defaults the local chain should
// apply to that remote, so this is reached with Canton as the remote and sourceChainSelector
// naming the sending chain.
func (a *CantonChainFamilyAdapter) GetDefaultRemoteChainConfig(sourceChainSelector, _ uint64) ccipadapters.RemoteChainDefaults {
	if isEVMSelector(sourceChainSelector) {
		// Canton messages are executed manually (ManuallyExecuteMessage), so no EVM Executor
		// supports Canton as a destination. SkipExecutorConfig keeps the EVM sequence from
		// reading or configuring an Executor contract for this lane; the EVM sequence writes
		// its own no-execution sentinel into the OnRamp's default-executor field, so no
		// family-specific address crosses this boundary.
		return ccipadapters.RemoteChainDefaults{
			AllowTrafficFrom: true,
			ExecutorDestChainConfig: ccipadapters.ExecutorDestChainConfig{
				USDCentsFee: 0,
				Enabled:     false,
			},
			BaseExecutionGasCost:      200_000,
			TokenReceiverAllowed:      false,
			MessageNetworkFeeUSDCents: 50,
			TokenNetworkFeeUSDCents:   50,
			SkipExecutorConfig:        true,
		}
	}

	return ccipadapters.RemoteChainDefaults{
		AllowTrafficFrom:          true,
		BaseExecutionGasCost:      50_000,
		TokenReceiverAllowed:      true,
		MessageNetworkFeeUSDCents: 10,
		TokenNetworkFeeUSDCents:   25,
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

func resolveDefaultInboundCCVs(
	chainSelector uint64,
	remoteCfg ccipadapters.RemoteChainConfig[[]byte, string],
) ([]datastore.AddressRef, error) {
	// A Canton CommitteeVerifier verifies Canton-origin messages on the destination chain (it is
	// the *outbound* CCV consumed by the remote EVM side). It cannot verify EVM-origin messages
	// on Canton, so it is never a valid *inbound* CCV for a Canton chain. chainlink-ccip's
	// resolveDefaultCCVs auto-resolves the "default"-qualifier CommitteeVerifier into
	// DefaultInboundCCVs when the lane YAML leaves it unset; honoring that here would silently
	// produce an inbound lane backed by a verifier that can never verify — the opposite of
	// fail-closed. Force the invalid-ccv sentinel as the default inbound CCV regardless of what
	// was auto-resolved. Explicit per-lane inbound CCVs set via laneMandatedInboundCCVs are
	// resolved separately (see configureChainForLanes) and remain honored.
	_ = remoteCfg.DefaultInboundCCVs // explicitly ignored; see comment above

	// Hardening: default inbound CCV is invalid-ccv (fail-closed). The sentinel is hardcoded
	// rather than resolved from the datastore so lanes stay fail-closed even when no invalid-ccv
	// qualifier was ever written to the address book.
	ref, err := invalidCCVAddressRef(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("build default invalid-ccv inbound ref: %w", err)
	}

	return []datastore.AddressRef{ref}, nil
}

// InvalidCCVRawInstanceAddress is the fail-closed sentinel CCV used as the default inbound CCV
// when a remote chain config does not specify one. It intentionally points at a CommitteeVerifier
// instance that can never verify, so unconfigured inbound lanes reject messages.
const InvalidCCVRawInstanceAddress = "invalid-ccv@ccvOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"

// invalidCCVAddressRef synthesizes the datastore ref for InvalidCCVRawInstanceAddress. The label
// carries the raw instance address, which is what dsutil.GetRawInstanceAddressFromAddressRef reads
// when the ref is converted to a binding.
func invalidCCVAddressRef(chainSelector uint64) (datastore.AddressRef, error) {
	raw, err := contracts.RawInstanceAddressFromString(InvalidCCVRawInstanceAddress)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("parse %q: %w", InvalidCCVRawInstanceAddress, err)
	}

	return datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(committeeverifierop.ContractType),
		Version:       committeeverifierop.Version,
		Qualifier:     "invalid-ccv",
		Address:       raw.InstanceAddress().Hex(),
		Labels:        datastore.NewLabelSet(raw.String()),
	}, nil
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

// BuildRemoteChainDefinition maps a resolved remote chain config to a lane ChainDefinition.
func BuildRemoteChainDefinition(remoteSelector uint64, remoteCfg ccipadapters.RemoteChainConfig[[]byte, string]) (*lanes.ChainDefinition, error) {
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
		FeeQuoterDestChainConfig:  BuildFeeQuoterDestChainConfig(remoteCfg.FeeQuoterDestChainConfig),
	}, nil
}

// BuildFeeQuoterDestChainConfig merges adapter overrides onto Canton FQ dest defaults.
func BuildFeeQuoterDestChainConfig(cfg ccipadapters.FeeQuoterDestChainConfigOverrides) lanes.FeeQuoterDestChainConfig {
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
const (
	minProductionCantonChainNOPs = 9
	minTestnetCantonChainNOPS    = 4
)

func (a *CantonChainFamilyAdapter) ValidateNOPsTopology(chainSelector string, nopCount int) error {
	chainSelectorUint, err := strconv.ParseUint(chainSelector, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid chain selector %q: %w", chainSelector, err)
	}
	isTestnet, err := chainsel.IsTestnetChain(chainSelectorUint)
	if err != nil {
		return fmt.Errorf("failed to determine if chain selector %q is testnet: %w", chainSelector, err)
	}
	if isTestnet {
		if nopCount < minTestnetCantonChainNOPS {
			return fmt.Errorf(
				"chain %q requires at least %d unique NOPs for testnet environments, got %d",
				chainSelector,
				minTestnetCantonChainNOPS,
				nopCount,
			)
		}

		return nil
	}
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
