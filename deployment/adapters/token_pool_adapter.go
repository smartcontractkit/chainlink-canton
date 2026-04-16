package adapters

import (
	"encoding/hex"
	"fmt"
	"strings"

	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	cantonsequences "github.com/smartcontractkit/chainlink-canton/deployment/sequences"
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
	return 10, nil
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
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeployTokenVerify(e deployment.Environment, in tokenadapters.DeployTokenInput) error {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeployTokenPoolForToken() *operations.Sequence[tokenadapters.DeployTokenPoolInput, sequences.OnChainOutput, chain.BlockChains] {
	return cantonsequences.DeployTokenPoolForToken
}

func (c CantonTokenAdapter) UpdateAuthorities() *operations.Sequence[tokenadapters.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *operations.Sequence[tokenadapters.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}
