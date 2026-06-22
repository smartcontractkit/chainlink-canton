package changesets

import (
	"fmt"
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	ccipadapters "github.com/smartcontractkit/chainlink-ccip/deployment/v2_0_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/deployment/adapters"
	committeeverifierop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	feequoterop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/nativeinstrument"
)

// HardenCantonInboundLaneConfig applies inbound-only hardening on Canton for one remote EVM source.
// Use instead of canton-build-lanes-cross-family when only UpdatePrices + invalid default inbound CCV are needed.
type HardenCantonInboundLaneConfig struct {
	RemoteSourceChainSelector uint64 `json:"remoteSourceChainSelector" yaml:"remoteSourceChainSelector"`
	// DefaultInboundCCVQualifiers resolves CommitteeVerifier refs from datastore (preferred).
	DefaultInboundCCVQualifiers []string `json:"defaultInboundCCVQualifiers,omitempty" yaml:"defaultInboundCCVQualifiers,omitempty"`
	// DefaultInboundCCVs optional inline refs; looked up in datastore when labels are missing.
	DefaultInboundCCVs []datastore.AddressRef `json:"defaultInboundCCVs,omitempty" yaml:"defaultInboundCCVs,omitempty"`
	USDPerUnitGas      int64                  `json:"usdPerUnitGas,omitempty" yaml:"usdPerUnitGas,omitempty"`
	MinDelay                  time.Duration          `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description               string                 `json:"description,omitempty" yaml:"description,omitempty"`
	// RemoteTokenPrices mirrors adapters.CantonRemoteTokenPricesFamilyExtraKey when non-default prices are needed.
	RemoteTokenPrices map[string]any `json:"remoteTokenPrices,omitempty" yaml:"remoteTokenPrices,omitempty"`
}

type HardenCantonInboundLane struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[HardenCantonInboundLaneConfig]] = HardenCantonInboundLane{}

func (h HardenCantonInboundLane) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[HardenCantonInboundLaneConfig]) error {
	cfg := config.Config
	if cfg.RemoteSourceChainSelector == 0 {
		return fmt.Errorf("remoteSourceChainSelector is required")
	}
	if len(cfg.DefaultInboundCCVQualifiers) == 0 && len(cfg.DefaultInboundCCVs) == 0 {
		return fmt.Errorf("defaultInboundCCVQualifiers or defaultInboundCCVs is required")
	}
	if cfg.MinDelay <= 0 {
		return fmt.Errorf("minDelay must be greater than zero")
	}
	if cfg.Description == "" {
		return fmt.Errorf("description is required")
	}

	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants",
			config.Participant, config.ChainSelector, len(chain.Participants))
	}

	if err := requireCCIPOwnerMCMSRef(e, config.ChainSelector); err != nil {
		return err
	}

	if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(global_config.ContractType),
		global_config.Version,
		"",
	)); err != nil {
		return fmt.Errorf("CantonGlobalConfig must be deployed on chain %d: %w", config.ChainSelector, err)
	}
	if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	)); err != nil {
		return fmt.Errorf("FeeQuoter must be deployed on chain %d: %w", config.ChainSelector, err)
	}

	if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		cfg.RemoteSourceChainSelector,
		datastore.ContractType(onramp.ContractType),
		onramp.Version,
		"",
	)); err != nil {
		return fmt.Errorf("remote OnRamp for chain %d must exist in datastore: %w", cfg.RemoteSourceChainSelector, err)
	}

	defaultInboundCCVs, err := resolveDefaultInboundCCVs(e.DataStore, config.ChainSelector, cfg)
	if err != nil {
		return err
	}
	for _, ref := range defaultInboundCCVs {
		if len(ref.Labels.List()) == 0 {
			return fmt.Errorf("default inbound CCV %q has no Canton deploy labels in datastore", ref.Qualifier)
		}
	}

	return nil
}

func resolveDefaultInboundCCVs(
	ds datastore.DataStore,
	cantonChainSelector uint64,
	cfg HardenCantonInboundLaneConfig,
) ([]datastore.AddressRef, error) {
	if len(cfg.DefaultInboundCCVQualifiers) > 0 {
		refs := make([]datastore.AddressRef, 0, len(cfg.DefaultInboundCCVQualifiers))
		for _, qualifier := range cfg.DefaultInboundCCVQualifiers {
			ref, err := ds.Addresses().Get(datastore.NewAddressRefKey(
				cantonChainSelector,
				datastore.ContractType(committeeverifierop.ContractType),
				committeeverifierop.Version,
				qualifier,
			))
			if err != nil {
				return nil, fmt.Errorf("resolve default inbound CCV qualifier %q: %w", qualifier, err)
			}
			refs = append(refs, ref)
		}

		return refs, nil
	}

	refs := make([]datastore.AddressRef, 0, len(cfg.DefaultInboundCCVs))
	for _, ref := range cfg.DefaultInboundCCVs {
		if len(ref.Labels.List()) > 0 {
			refs = append(refs, ref)
			continue
		}
		chainSelector := ref.ChainSelector
		if chainSelector == 0 {
			chainSelector = cantonChainSelector
		}
		contractType := ref.Type
		if contractType == "" {
			contractType = datastore.ContractType(committeeverifierop.ContractType)
		}
		version := committeeverifierop.Version
		if ref.Version != nil {
			version = ref.Version
		}
		resolved, err := ds.Addresses().Get(datastore.NewAddressRefKey(
			chainSelector, contractType, version, ref.Qualifier,
		))
		if err != nil {
			return nil, fmt.Errorf("resolve default inbound CCV from datastore (qualifier %q): %w", ref.Qualifier, err)
		}
		refs = append(refs, resolved)
	}

	return refs, nil
}

func (h HardenCantonInboundLane) Apply(e cldf.Environment, config CantonCSDeps[HardenCantonInboundLaneConfig]) (cldf.ChangesetOutput, error) {
	cfg := config.Config
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]
	proposalDriven := len(participant.ReadAsPartyIDs) > 0

	globalConfigRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(global_config.ContractType),
		global_config.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve global config: %w", err)
	}

	feeQuoterRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(feequoterop.ContractType),
		feequoterop.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve fee quoter: %w", err)
	}

	remoteOnRampRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		cfg.RemoteSourceChainSelector,
		datastore.ContractType(onramp.ContractType),
		onramp.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve remote OnRamp: %w", err)
	}
	if !gethcommon.IsHexAddress(remoteOnRampRef.Address) {
		return cldf.ChangesetOutput{}, fmt.Errorf("remote OnRamp address %q is not valid hex", remoteOnRampRef.Address)
	}

	nativeInstrument, err := nativeinstrument.ResolveNativeInstrumentID(
		e.GetContext(), participant, e.DataStore, config.ChainSelector,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve native fee token instrument: %w", err)
	}

	familyExtras := map[string]any(nil)
	if len(cfg.RemoteTokenPrices) > 0 {
		familyExtras = map[string]any{
			adapters.CantonRemoteTokenPricesFamilyExtraKey: map[string]any{
				fmt.Sprintf("%d", cfg.RemoteSourceChainSelector): cfg.RemoteTokenPrices,
			},
		}
	}

	tokenPrices, err := adapters.ResolveTokenPricesForRemoteDest(
		e.DataStore,
		ccipadapters.ConfigureChainForLanesInput{
			ChainSelector:  config.ChainSelector,
			FamilyExtras:   familyExtras,
		},
		cfg.RemoteSourceChainSelector,
		&nativeInstrument,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve token prices: %w", err)
	}

	usdPerUnitGas := big.NewInt(cfg.USDPerUnitGas)
	if cfg.USDPerUnitGas == 0 {
		usdPerUnitGas = big.NewInt(38)
	}

	defaultInboundCCVs, err := resolveDefaultInboundCCVs(e.DataStore, config.ChainSelector, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	out, err := cld_ops.ExecuteSequence(e.OperationsBundle, sequences.HardenCantonInboundLane, e.BlockChains, sequences.HardenCantonInboundLaneInput{
		CantonChainSelector:       config.ChainSelector,
		RemoteSourceChainSelector: cfg.RemoteSourceChainSelector,
		GlobalConfigRef:           globalConfigRef,
		FeeQuoterRef:              feeQuoterRef,
		DefaultInboundCCVs:        defaultInboundCCVs,
		RemoteOnRampAddress:       gethcommon.HexToAddress(remoteOnRampRef.Address),
		TokenPrices:               tokenPrices,
		USDPerUnitGas:             usdPerUnitGas,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("harden canton inbound lane: %w", err)
	}

	if proposalDriven && len(out.Output.BatchOps) == 0 {
		e.Logger.Infof("Canton inbound hardening already applied for remote source %d; no MCMS proposal generated", cfg.RemoteSourceChainSelector)
		return cldf.ChangesetOutput{}, nil
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, proposalDriven,
		cantonmcms.QualifierCCIPOwner, cfg.MinDelay, cfg.Description, nil, out.Output.BatchOps,
	)
}
