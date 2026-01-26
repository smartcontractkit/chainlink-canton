package deployment

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/noders-team/go-daml/pkg/client"

	"github.com/smartcontractkit/chainlink-canton-internal/deployment/view"
)

type CantonChainView struct {
	ChainSelector uint64 `json:"chainSelector,omitempty"`

	OnRamp       view.OnRampView       `json:"onRamp,omitempty"`
	GlobalConfig view.GlobalConfigView `json:"globalConfig,omitempty"`
	// TODO: Add other views as they are implemented
	// OffRamp view.OffRampView `json:"offRamp,omitempty"`
	// FeeQuoter view.FeeQuoterView `json:"feeQuoter,omitempty"`
}

type CantonChainState struct {
	// Canton connection info
	BindingClient *client.DamlBindingClient // Binding client for operations and state queries
	Party         string                    // Party that can read the contracts

	// CCIP Contract IDs
	GlobalConfigContractID       string
	TokenAdminRegistryContractID string
	FeeQuoterContractID          string
	OnRampContractID             string
	OffRampContractID            string
	PerPartyRouterContractID     string
	CommitteeVerifierContractID  string
	CCVRegistryContractID        string

	// LINK Token related
	LinkTokenRegistryContractID string
}

func (s CantonChainState) GenerateView(e *cldf.Environment, selector uint64, chainName string) (CantonChainView, error) {
	lggr := e.Logger
	chainView := CantonChainView{
		ChainSelector: selector,
	}

	lggr.Infow("generating Canton chain view", "chain", chainName, "selector", selector)

	ctx := context.Background()

	// Use StateService from binding client
	if s.BindingClient == nil {
		return CantonChainView{}, fmt.Errorf("bindingClient is required for generating views")
	}

	var mu sync.Mutex
	g, ctxG := errgroup.WithContext(ctx)

	// OnRamp
	if s.OnRampContractID != "" {
		g.Go(func() error {
			// List known packages to find the package ID for ccip-onramp
			ListKnownPackagesResp, err := s.BindingClient.PackageMng.ListKnownPackages(ctx)
			if err != nil {
				return fmt.Errorf("failed to list known packages: %w", err)
			}

			var ccipOnRampPkgID string
			for _, p := range ListKnownPackagesResp {
				if strings.Contains(strings.ToLower(p.Name), "ccip-onramp") {
					ccipOnRampPkgID = p.PackageID
					break
				}
			}
			if ccipOnRampPkgID == "" {
				return fmt.Errorf("failed to find ccip-onramp package")
			}

			onRampView, err := view.GenerateOnRampView(ctxG, s.BindingClient.StateService, s.OnRampContractID, ccipOnRampPkgID, s.Party)
			if err != nil {
				return fmt.Errorf("failed to generate onramp view for onramp %s: %w", s.OnRampContractID, err)
			}
			mu.Lock()
			chainView.OnRamp = onRampView
			mu.Unlock()
			lggr.Infow("generated onRamp view", "onRampContractID", s.OnRampContractID, "chain", chainName)
			return nil
		})
	}

	// GlobalConfig
	if s.GlobalConfigContractID != "" {
		g.Go(func() error {
			// List known packages to find the package ID for ccip-common
			ListKnownPackagesResp, err := s.BindingClient.PackageMng.ListKnownPackages(ctx)
			if err != nil {
				return fmt.Errorf("failed to list known packages: %w", err)
			}

			var ccipCommonPkgID string
			for _, p := range ListKnownPackagesResp {
				if strings.Contains(strings.ToLower(p.Name), "ccip-common") {
					ccipCommonPkgID = p.PackageID
					break
				}
			}
			if ccipCommonPkgID == "" {
				return fmt.Errorf("failed to find ccip-common package")
			}

			globalConfigView, err := view.GenerateGlobalConfigView(ctxG, s.BindingClient.StateService, s.GlobalConfigContractID, ccipCommonPkgID, s.Party)
			if err != nil {
				return fmt.Errorf("failed to generate globalconfig view for globalconfig %s: %w", s.GlobalConfigContractID, err)
			}
			mu.Lock()
			chainView.GlobalConfig = globalConfigView
			mu.Unlock()
			lggr.Infow("generated GlobalConfig view", "globalConfigContractID", s.GlobalConfigContractID, "chain", chainName)
			return nil
		})
	}

	// TODO: Add other views as they are implemented
	// OffRamp, FeeQuoter, etc.

	return chainView, g.Wait()
}

// LoadOnchainStateCanton loads chain state for Canton chains from env
func LoadOnchainStateCanton(env cldf.Environment) (map[uint64]CantonChainState, error) {
	// Canton doesn't have a BlockChains.CantonChains() method like Sui
	// We need to get chains from the address book
	// For now, we'll iterate through all chains in the address book and filter for Canton contracts

	cantonChains := make(map[uint64]CantonChainState)

	// Get all chains from the address book
	// Note: This assumes we can get chain selectors from the address book
	// In practice, you might need to maintain a separate registry of Canton chain selectors
	// or get them from the environment configuration

	// For now, we'll create a helper that loads state from addresses for a given chain selector
	// This will be called for each Canton chain selector

	return cantonChains, nil
}

// LoadCantonChainStateFromAddresses loads Canton chain state from address book entries
func LoadCantonChainStateFromAddresses(
	chainSelector uint64,
	addresses map[string]cldf.TypeAndVersion,
	bindingClient *client.DamlBindingClient,
	party string,
) (CantonChainState, error) {
	chainState := CantonChainState{
		BindingClient: bindingClient,
		Party:         party,
	}

	for addr, typeAndVersion := range addresses {
		// Parse addresses based on type
		switch typeAndVersion.Type {
		case CantonCCIPOnRampType:
			chainState.OnRampContractID = addr
		case CantonCCIPOffRampType:
			chainState.OffRampContractID = addr
		case CantonCCIPFeeQuoterType:
			chainState.FeeQuoterContractID = addr
		case CantonCCIPTokenAdminRegistryType:
			chainState.TokenAdminRegistryContractID = addr
		case CantonCCIPCommitteeVerifierType:
			chainState.CommitteeVerifierContractID = addr
		case CantonCCIPPerPartyRouterType:
			chainState.PerPartyRouterContractID = addr
		case CantonCCIPCCVRegistryType:
			chainState.CCVRegistryContractID = addr
		case CantonCCIPGlobalConfigType:
			chainState.GlobalConfigContractID = addr
		case CantonLinkTokenRegistryType:
			chainState.LinkTokenRegistryContractID = addr
		}
	}

	return chainState, nil
}
