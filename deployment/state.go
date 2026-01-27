package deployment

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/noders-team/go-daml/pkg/client"
	"github.com/noders-team/go-daml/pkg/model"

	"github.com/noders-team/go-daml/pkg/service/ledger"
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
	GlobalConfigContractAddress       string
	TokenAdminRegistryContractAddress string
	FeeQuoterContractAddress          string
	OnRampContractAddress             string
	OffRampContractAddress            string
	PerPartyRouterContractAddress     string
	CommitteeVerifierContractAddress  string
	CCVRegistryContractAddress        string

	// LINK Token related
	LinkTokenRegistryContractAddress string
}

// GenerateView generates contracts views are active and visible to a specific party
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
	if s.OnRampContractAddress != "" {
		g.Go(func() error {
			// get address from datastore for onramp
			address, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType("onramp"), semver.MustParse("1.0.0"), ""))
			if err != nil {
				return fmt.Errorf("failed to get address: %w", err)
			}

			// parse party from address
			instanceID, party, err := parsePartyFromAddress(address.Address)
			if err != nil {
				return fmt.Errorf("failed to parse party from address: %w", err)
			}

			// get all contract for this party
			contracts, err := getAllActiveContractsForParty(ctx, e, selector, s.BindingClient.StateService, party)
			if err != nil {
				return fmt.Errorf("failed to get active contracts: %w", err)
			}

			// find specific contract for this party based on the createdEvent from contractCreation
			contractID, err := findContractBasedOnModuleIdAndExpectedEnvID(contracts, instanceID)
			if err != nil {
				return fmt.Errorf("failed to find contract: %w", err)
			}

			onRampView, err := view.GenerateOnRampView(ctxG, s.BindingClient.StateService, s.OnRampContractAddress, contractID, s.Party)
			if err != nil {
				return fmt.Errorf("failed to generate onramp view for onramp %s: %w", s.OnRampContractAddress, err)
			}
			mu.Lock()
			chainView.OnRamp = onRampView
			mu.Unlock()
			lggr.Infow("generated onRamp view", "onRampContractAddress", s.OnRampContractAddress, "chain", chainName)
			return nil
		})
	}

	// GlobalConfig
	if s.GlobalConfigContractAddress != "" {
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

			globalConfigView, err := view.GenerateGlobalConfigView(ctxG, s.BindingClient.StateService, s.GlobalConfigContractAddress, ccipCommonPkgID, s.Party)
			if err != nil {
				return fmt.Errorf("failed to generate globalconfig view for globalconfig %s: %w", s.GlobalConfigContractAddress, err)
			}
			mu.Lock()
			chainView.GlobalConfig = globalConfigView
			mu.Unlock()
			lggr.Infow("generated GlobalConfig view", "globalConfigContractAddress", s.GlobalConfigContractAddress, "chain", chainName)
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
			chainState.OnRampContractAddress = addr
		case CantonCCIPOffRampType:
			chainState.OffRampContractAddress = addr
		case CantonCCIPFeeQuoterType:
			chainState.FeeQuoterContractAddress = addr
		case CantonCCIPTokenAdminRegistryType:
			chainState.TokenAdminRegistryContractAddress = addr
		case CantonCCIPCommitteeVerifierType:
			chainState.CommitteeVerifierContractAddress = addr
		case CantonCCIPPerPartyRouterType:
			chainState.PerPartyRouterContractAddress = addr
		case CantonCCIPCCVRegistryType:
			chainState.CCVRegistryContractAddress = addr
		case CantonCCIPGlobalConfigType:
			chainState.GlobalConfigContractAddress = addr
		case CantonLinkTokenRegistryType:
			chainState.LinkTokenRegistryContractAddress = addr
		}
	}

	return chainState, nil
}

func getAllActiveContractsForParty(ctx context.Context, e *cldf.Environment, selector uint64, stateService ledger.StateService, party string) ([]*model.CreatedEvent, error) {
	// Get current offset
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &model.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger end: %w", err)
	}

	// Query active contracts for the party using wildcard filter
	// Following the pattern from GetActiveContractsForPartyWithFilters
	req := &model.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.Offset,
		EventFormat: &model.EventFormat{
			FiltersByParty: map[string]*model.Filters{
				party: {Inclusive: &model.InclusiveFilters{
					TemplateFilters: []*model.TemplateFilter{}, // Empty = wildcard
				}},
			},
			Verbose: true,
		},
	}

	responseChan, errorChan := stateService.GetActiveContracts(ctx, req)

	for {
		select {
		case resp, ok := <-responseChan:
			if !ok {
				return []*model.CreatedEvent{}, nil
			}
			if resp == nil || resp.ContractEntry == nil {
				continue
			}
			entry, ok := resp.ContractEntry.(*model.ActiveContractEntry)
			if !ok || entry.ActiveContract == nil || entry.ActiveContract.CreatedEvent == nil {
				continue
			}
			return []*model.CreatedEvent{entry.ActiveContract.CreatedEvent}, nil
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			if err != nil {
				return nil, fmt.Errorf("failed to receive active contract: %w", err)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func findContractBasedOnModuleIdAndExpectedEnvID(contracts []*model.CreatedEvent, expectedEnvID string) (string, error) {
	for _, contract := range contracts {
		// instanceId assertion
		for _, field := range contract.CreateArguments.(map[string]interface{}) {
			if field == "instanceId" {
				if field != expectedEnvID {
					continue
				}
			}
		}

		return contract.ContractID, nil
	}
	return "", fmt.Errorf("contract not found with expectedEnvID %s", expectedEnvID)
}

// instanceID@party
// Address looks like this local-v1-globalconfig@participant1-localparty-1::1220cc68db908cfcb2bea8383297dc05f6d2c58566866a3e47a397efbdc29c1cb0dd
func parsePartyFromAddress(address string) (string, string, error) {
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid address: %s", address)
	}
	return parts[0], parts[1], nil
}
