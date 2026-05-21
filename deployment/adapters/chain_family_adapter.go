package adapters

import (
	"encoding/binary"
	"fmt"
	"math/big"

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
)

// SetRuntimeDataStore forwards to the shared runtime datastore used by lane configure sequences.
func SetRuntimeDataStore(ds datastore.DataStore) {
	dsutil.SetRuntimeDataStore(ds)
}

var _ ccipadapters.ChainFamily = (*CantonChainFamilyAdapter)(nil)

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
			ds := dsutil.RuntimeDataStore()
			if ds == nil {
				return ccipseq.OnChainOutput{}, fmt.Errorf("runtime datastore is not set")
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

			localCommitteeVerifiers := convertCommitteeVerifierConfigs(input.CommitteeVerifiers)
			var out ccipseq.OnChainOutput

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
					OnRamp:                   input.OnRamp,
					OffRamp:                  input.OffRamp,
					Router:                   input.Router,
					FeeQuoter:                input.FeeQuoter,
				}

				remoteChain, err := remoteChainDefinition(remoteSelector, remoteCfg)
				if err != nil {
					return out, err
				}

				out, err = ccipseq.RunAndMergeSequence(
					b,
					chains,
					sequences.ConfigureLaneLegAsSourceWithDataStore,
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
		},
	)
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
func (a *CantonChainFamilyAdapter) GetDefaultFeeQuoterDestChainConfig() ccipadapters.FeeQuoterDestChainConfig {
	return ccipadapters.FeeQuoterDestChainConfig{
		OverrideExistingConfig:      false,
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         CantonFamilySelector,
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		LinkFeeMultiplierPercent:    0,
		USDPerUnitGas:               big.NewInt(0),
	}
}

// GetDefaultFinalityConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultFinalityConfig() finality.Config {
	return finality.Config{
		WaitForFinality: true,
	}
}

// GetDefaultRemoteChainConfig implements [adapters.ChainFamily].
func (a *CantonChainFamilyAdapter) GetDefaultRemoteChainConfig() ccipadapters.RemoteChainDefaults {
	return ccipadapters.RemoteChainDefaults{
		AllowTrafficFrom: true, // TODO: check what this does?
		ExecutorDestChainConfig: ccipadapters.ExecutorDestChainConfig{
			USDCentsFee: 0,
			Enabled:     true,
		},
		BaseExecutionGasCost:      50_000,
		TokenReceiverAllowed:      false, // TODO: check what this does?
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

	return &lanes.ChainDefinition{
		Selector:                  remoteSelector,
		ExecutorDestChainConfig:   lanes.ExecutorDestChainConfig(remoteCfg.ExecutorDestChainConfig),
		AddressBytesLength:        remoteCfg.AddressBytesLength,
		BaseExecutionGasCost:      remoteCfg.BaseExecutionGasCost,
		TokenReceiverAllowed:      remoteCfg.TokenReceiverAllowed,
		MessageNetworkFeeUSDCents: remoteCfg.MessageNetworkFeeUSDCents,
		TokenNetworkFeeUSDCents:   remoteCfg.TokenNetworkFeeUSDCents,
		OnRamp:                    remoteCfg.OnRamps[0],
		OffRamp:                   remoteCfg.OffRamp,
		Router:                    router,
		FeeQuoterDestChainConfig: lanes.FeeQuoterDestChainConfig{
			OverrideExistingConfig:      remoteCfg.FeeQuoterDestChainConfig.OverrideExistingConfig,
			IsEnabled:                   remoteCfg.FeeQuoterDestChainConfig.IsEnabled,
			MaxDataBytes:                remoteCfg.FeeQuoterDestChainConfig.MaxDataBytes,
			MaxPerMsgGasLimit:           remoteCfg.FeeQuoterDestChainConfig.MaxPerMsgGasLimit,
			DestGasOverhead:             remoteCfg.FeeQuoterDestChainConfig.DestGasOverhead,
			DestGasPerPayloadByteBase:   remoteCfg.FeeQuoterDestChainConfig.DestGasPerPayloadByteBase,
			ChainFamilySelector:         binary.BigEndian.Uint32(remoteCfg.FeeQuoterDestChainConfig.ChainFamilySelector[:]),
			DefaultTokenFeeUSDCents:     remoteCfg.FeeQuoterDestChainConfig.DefaultTokenFeeUSDCents,
			DefaultTokenDestGasOverhead: remoteCfg.FeeQuoterDestChainConfig.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:           remoteCfg.FeeQuoterDestChainConfig.DefaultTxGasLimit,
			NetworkFeeUSDCents:          remoteCfg.FeeQuoterDestChainConfig.NetworkFeeUSDCents,
			V2Params: &lanes.FeeQuoterV2Params{
				LinkFeeMultiplierPercent: remoteCfg.FeeQuoterDestChainConfig.LinkFeeMultiplierPercent,
				USDPerUnitGas:            remoteCfg.FeeQuoterDestChainConfig.USDPerUnitGas,
			},
		},
	}, nil
}
