package adapters

import (
	"encoding/hex"
	"fmt"
	"strings"

	evmadapters "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v2_0_0/operations/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v2_0_0/operations/offramp"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v2_0_0/operations/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v1_2_0/operations/router"
	"github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	datastore_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

var (
	_ tokenadapters.TokenAdapter = (*CantonTokenAdapter)(nil)
	_ lanes.LaneAdapter          = (*CantonLaneAdapter)(nil)
)

// TODO: move this to chainlink-canton/deployment.
type CantonTokenAdapter struct {
	base tokenadapters.TokenAdapter
}

// NewTokenAdapter creates a new Canton token adapter.
// A "base" adapter needs to be passed in, currently assumed to be the EVM token adapter,
// in order to achieve all functionality.
// TODO: this needs to be fully implemented for Canton.
func NewTokenAdapter(base tokenadapters.TokenAdapter) *CantonTokenAdapter {
	return &CantonTokenAdapter{
		base: base,
	}
}

// AddressRefToBytes implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return t.base.AddressRefToBytes(ref)
}

// ConfigureTokenForTransfersSequence implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) ConfigureTokenForTransfersSequence() *operations.Sequence[tokenadapters.ConfigureTokenForTransfersInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.ConfigureTokenForTransfersSequence()
}

// DeployToken implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) DeployToken() *operations.Sequence[tokenadapters.DeployTokenInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.DeployToken()
}

// DeployTokenPoolForToken implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) DeployTokenPoolForToken() *operations.Sequence[tokenadapters.DeployTokenPoolInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.DeployTokenPoolForToken()
}

// DeployTokenVerify implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) DeployTokenVerify(e deployment.Environment, in tokenadapters.DeployTokenInput) error {
	return t.base.DeployTokenVerify(e, in)
}

// ManualRegistration implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) ManualRegistration() *operations.Sequence[tokenadapters.ManualRegistrationSequenceInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.ManualRegistration()
}

// SetTokenPoolRateLimits implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) SetTokenPoolRateLimits() *operations.Sequence[tokenadapters.TPRLRemotes, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.SetTokenPoolRateLimits()
}

// UpdateAuthorities implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) UpdateAuthorities() *operations.Sequence[tokenadapters.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	return t.base.UpdateAuthorities()
}

// MigrateLockReleasePoolLiquiditySequence implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *operations.Sequence[tokenadapters.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.MigrateLockReleasePoolLiquiditySequence()
}

func (t *CantonTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef, token []byte) (uint8, error) {
	// TODO: Need to actually implement this behavior for Canton instead of reusing the EVM behavior.
	// return t.base.DeriveTokenDecimals(e, chainSelector, ref, token)
	return 0, nil
}

func (t *CantonTokenAdapter) DeriveTokenPoolCounterpart(e deployment.Environment, chainSelector uint64, tokenPoolAddress, tokenAddress []byte) ([]byte, error) {
	return t.base.DeriveTokenPoolCounterpart(e, chainSelector, tokenPoolAddress, tokenAddress)
}

func (t *CantonTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef) ([]byte, error) {
	addr, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return nil, fmt.Errorf("failed to get address for %v on chain %d: %w", ref, chainSelector, err)
	}
	// Address is stored as hex string
	// Remove 0x prefix if present
	cleanAddr := strings.TrimPrefix(addr.Address, "0x")

	return hex.DecodeString(cleanAddr)
}

func cantonInstanceAddressBytes(ref datastore.AddressRef) ([]byte, error) {
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

// CantonLaneAdapter implements lanes.LaneAdapter for Canton using instance-address encoding
// for on-chain contract refs in the CLDF datastore. Sequence operations delegate to the EVM
// CCV adapter; Canton-specific deployment paths (e.g. devenv impl) use separate changesets.
type CantonLaneAdapter struct {
	*evmadapters.ChainFamilyAdapter
}

// NewCantonLaneAdapter builds a lane adapter for Canton (version aligned with EVM CCV 2.0.0).
func NewCantonLaneAdapter() *CantonLaneAdapter {
	return &CantonLaneAdapter{ChainFamilyAdapter: &evmadapters.ChainFamilyAdapter{}}
}

// GetOnRampAddress implements lanes.LaneAdapter.
func (c *CantonLaneAdapter) GetOnRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore_utils.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(onramp.ContractType),
		Version:       onramp.Version,
	}, chainSelector, cantonInstanceAddressBytes)
}

// GetOffRampAddress implements lanes.LaneAdapter.
func (c *CantonLaneAdapter) GetOffRampAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore_utils.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(offramp.ContractType),
		Version:       offramp.Version,
	}, chainSelector, cantonInstanceAddressBytes)
}

// GetFQAddress implements lanes.LaneAdapter.
func (c *CantonLaneAdapter) GetFQAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore_utils.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(fee_quoter.ContractType),
		Version:       fee_quoter.Version,
	}, chainSelector, cantonInstanceAddressBytes)
}

// GetRouterAddress implements lanes.LaneAdapter.
func (c *CantonLaneAdapter) GetRouterAddress(ds datastore.DataStore, chainSelector uint64) ([]byte, error) {
	return datastore_utils.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(router.ContractType),
		Version:       router.Version,
	}, chainSelector, cantonInstanceAddressBytes)
}
