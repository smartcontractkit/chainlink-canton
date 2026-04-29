package adapters

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonsequences "github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var (
	_ tokenadapters.TokenAdapter = &CantonTokenAdapter{}
)

type CantonTokenAdapter struct{}

func (c CantonTokenAdapter) ConfigureTokenForTransfersSequence() *operations.Sequence[tokenadapters.ConfigureTokenForTransfersInput, sequences.OnChainOutput, chain.BlockChains] {
	return cantonsequences.ConfigureTokenForTransfers
}

func (c CantonTokenAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

func (c CantonTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef) ([]byte, error) {
	tokenRefs := e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType("Token")),
		datastore.AddressRefByQualifier(poolRef.Qualifier),
	)
	if len(tokenRefs) != 1 {
		return nil, fmt.Errorf("expected exactly one precomputed token ref for qualifier %q, got %d", poolRef.Qualifier, len(tokenRefs))
	}
	addr := strings.TrimSpace(tokenRefs[0].Address)
	if addr == "" {
		return nil, fmt.Errorf("precomputed token ref for qualifier %q has empty address", poolRef.Qualifier)
	}
	if rawAddr, rawErr := contracts.RawInstanceAddressFromString(addr); rawErr == nil {
		return rawAddr.InstanceAddress().Bytes(), nil
	}
	addrBytes, err := hex.DecodeString(strings.TrimPrefix(addr, "0x"))
	if err != nil || len(addrBytes) == 0 {
		return nil, fmt.Errorf("invalid precomputed token address %q for qualifier %q", addr, poolRef.Qualifier)
	}

	return addrBytes, nil
}

func (c CantonTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef, token []byte) (uint8, error) {
	poolAddressRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		poolRef.Type,
		poolRef.Version,
		poolRef.Qualifier,
	))
	if err != nil {
		return 0, fmt.Errorf("resolve Canton token pool ref: %w", err)
	}

	chain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok || len(chain.Participants) == 0 {
		return 0, fmt.Errorf("canton chain with selector %d not found", chainSelector)
	}
	participant := chain.Participants[0]
	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}

	poolAddress := contracts.HexToInstanceAddress(poolAddressRef.Address)
	switch poolRef.Type {
	case datastore.ContractType("LockReleaseTokenPool"):
		activePool, err := opcontract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
			poolAddress,
		)
		if err != nil {
			return 0, fmt.Errorf("find active lock/release pool %s: %w", poolAddressRef.Address, err)
		}
		parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
		if err != nil {
			return 0, fmt.Errorf("parse active lock/release pool %s: %w", poolAddressRef.Address, err)
		}

		//nolint:gosec // Decimals should never exceed uint8
		return uint8(parsedPool.Decimals), nil
	case datastore.ContractType("BurnMintTokenPool"):
		activePool, err := opcontract.FindActiveContractByInstanceAddress(
			ctx,
			participant.LedgerServices.State,
			participant.PartyID,
			burnminttokenpool.BurnMintTokenPool{}.GetTemplateID(),
			poolAddress,
		)
		if err != nil {
			return 0, fmt.Errorf("find active burn/mint pool %s: %w", poolAddressRef.Address, err)
		}
		parsedPool, err := bindings.UnmarshalCreatedEvent[burnminttokenpool.BurnMintTokenPool](activePool.GetCreatedEvent())
		if err != nil {
			return 0, fmt.Errorf("parse active burn/mint pool %s: %w", poolAddressRef.Address, err)
		}

		//nolint:gosec // Decimals should never exceed uint8
		return uint8(parsedPool.Decimals), nil
	default:
		return 0, fmt.Errorf("unsupported Canton token pool type %q", poolRef.Type)
	}
}

func (c CantonTokenAdapter) DeriveTokenPoolCounterpart(e deployment.Environment, chainSelector uint64, tokenPool []byte, token []byte) ([]byte, error) {
	return tokenPool, nil
}

func (c CantonTokenAdapter) ManualRegistration() *operations.Sequence[tokenadapters.ManualRegistrationSequenceInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) SetTokenPoolRateLimits() *operations.Sequence[tokenadapters.TPRLRemotes, sequences.OnChainOutput, chain.BlockChains] {
	return cantonsequences.SetTokenPoolRateLimits
}

func (c CantonTokenAdapter) DeployToken() *operations.Sequence[tokenadapters.DeployTokenInput, sequences.OnChainOutput, chain.BlockChains] {
	return operations.NewSequence(
		"canton/token-adapter/deploy-token",
		semver.MustParse("2.0.0"),
		"Resolves an existing Canton instrument-backed token ref without deploying a token",
		func(_ operations.Bundle, _ chain.BlockChains, input tokenadapters.DeployTokenInput) (sequences.OnChainOutput, error) {
			ref, err := resolveExistingCantonTokenRef(input.ExistingDataStore, input)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}

			return sequences.OnChainOutput{Addresses: []datastore.AddressRef{ref}}, nil
		},
	)
}

func (c CantonTokenAdapter) DeployTokenVerify(e deployment.Environment, in tokenadapters.DeployTokenInput) error {
	_, err := resolveExistingCantonTokenRef(e.DataStore, in)
	return err
}

func (c CantonTokenAdapter) DeployTokenPoolForToken() *operations.Sequence[tokenadapters.DeployTokenPoolInput, sequences.OnChainOutput, chain.BlockChains] {
	return cantonsequences.DeployTokenPoolForToken
}

func (c CantonTokenAdapter) UpdateAuthorities() *operations.Sequence[tokenadapters.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	return operations.NewSequence(
		"canton/token-adapter/update-authorities",
		semver.MustParse("2.0.0"),
		"No-op for Canton instrument-backed tokens; no token or pool ownership transfer is required",
		func(_ operations.Bundle, _ *deployment.Environment, _ tokenadapters.UpdateAuthoritiesInput) (sequences.OnChainOutput, error) {
			return sequences.OnChainOutput{}, nil
		},
	)
}

func (c CantonTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *operations.Sequence[tokenadapters.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func resolveExistingCantonTokenRef(ds datastore.DataStore, in tokenadapters.DeployTokenInput) (datastore.AddressRef, error) {
	if ds == nil {
		return datastore.AddressRef{}, fmt.Errorf("existing datastore is required to resolve an existing Canton token")
	}

	refs := ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(in.ChainSelector),
		datastore.AddressRefByType(datastore.ContractType("Token")),
	)
	if len(refs) == 0 {
		return datastore.AddressRef{}, fmt.Errorf("no Canton token refs found on chain %d", in.ChainSelector)
	}

	lookupValues := uniqueNonEmpty(in.Symbol, in.Name)
	for _, value := range lookupValues {
		if matches := filterTokenRefs(refs, func(ref datastore.AddressRef) bool {
			return strings.EqualFold(ref.Qualifier, value)
		}); len(matches) == 1 {
			return matches[0], nil
		} else if len(matches) > 1 {
			return datastore.AddressRef{}, fmt.Errorf("multiple Canton token refs matched qualifier %q on chain %d", value, in.ChainSelector)
		}

		label := "instrument-id:" + value
		if matches := filterTokenRefs(refs, func(ref datastore.AddressRef) bool {
			for _, existingLabel := range ref.Labels.List() {
				if strings.EqualFold(existingLabel, label) {
					return true
				}
			}

			return false
		}); len(matches) == 1 {
			return matches[0], nil
		} else if len(matches) > 1 {
			return datastore.AddressRef{}, fmt.Errorf("multiple Canton token refs matched label %q on chain %d", label, in.ChainSelector)
		}
	}

	if len(refs) == 1 {
		return refs[0], nil
	}

	return datastore.AddressRef{}, fmt.Errorf(
		"could not resolve a unique Canton token ref on chain %d from symbol=%q name=%q; found %d token refs",
		in.ChainSelector,
		in.Symbol,
		in.Name,
		len(refs),
	)
}

func filterTokenRefs(refs []datastore.AddressRef, predicate func(datastore.AddressRef) bool) []datastore.AddressRef {
	filtered := make([]datastore.AddressRef, 0, len(refs))
	for _, ref := range refs {
		if predicate(ref) {
			filtered = append(filtered, ref)
		}
	}

	return filtered
}

func uniqueNonEmpty(values ...string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}
