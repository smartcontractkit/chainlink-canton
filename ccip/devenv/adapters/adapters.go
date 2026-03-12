package adapters

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-ccip/deployment/v1_7_0/adapters"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var (
	_ tokenadapters.TokenAdapter = &CantonTokenAdapter{}
	_ adapters.ChainFamily       = &CantonAdapter{}
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
	// Canton "addresses" are 32-byte instance addresses, not 20-byte EVM addresses.
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
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
func (t *CantonTokenAdapter) DeployTokenVerify(e deployment.Environment, in any) error {
	return t.base.DeployTokenVerify(e, in)
}

// ManualRegistration implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) ManualRegistration() *operations.Sequence[tokenadapters.ManualRegistrationInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.ManualRegistration()
}

// RegisterToken implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) RegisterToken() *operations.Sequence[tokenadapters.RegisterTokenInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.RegisterToken()
}

// SetPool implements tokens.TokenAdapter.
func (t *CantonTokenAdapter) SetPool() *operations.Sequence[tokenadapters.SetPoolInput, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.SetPool()
}

func (t *CantonTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef) (uint8, error) {
	// TODO: Need to actually implement this behavior for Canton instead of reusing the EVM behavior.
	// return t.base.DeriveTokenDecimals(e, chainSelector, ref)
	return 0, nil
}

func (t *CantonTokenAdapter) DeriveTokenPoolCounterpart(e deployment.Environment, chainSelector uint64, tokenPoolAddress, tokenAddress []byte) ([]byte, error) {
	return t.base.DeriveTokenPoolCounterpart(e, chainSelector, tokenPoolAddress, tokenAddress)
}

func (t *CantonTokenAdapter) SetTokenPoolRateLimits() *operations.Sequence[tokenadapters.RateLimiterConfigInputs, sequences.OnChainOutput, chain.BlockChains] {
	return t.base.SetTokenPoolRateLimits()
}

func (t *CantonTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef) ([]byte, error) {
	if tokenAddress, err := deriveInstrumentTokenAddress(e, chainSelector, ref); err == nil {
		return tokenAddress, nil
	}

	addr, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return nil, fmt.Errorf("failed to get address for %v on chain %d: %w", ref, chainSelector, err)
	}
	// Raw Canton instance addresses are "<instance-id>@<party-id>" and must be keccak256 hashed.
	if rawAddr, rawErr := contracts.RawInstanceAddressFromString(addr.Address); rawErr == nil {
		return rawAddr.InstanceAddress().Bytes(), nil
	}

	// Hex-encoded addresses must be Keccak256 hashed to derive the canonical token address.
	cleanAddr := strings.TrimPrefix(addr.Address, "0x")
	addrBytes, err := hex.DecodeString(cleanAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token address %q: %w", addr.Address, err)
	}
	return crypto.Keccak256(addrBytes), nil
}

func deriveInstrumentTokenAddress(e deployment.Environment, chainSelector uint64, ref datastore.AddressRef) ([]byte, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok || len(cantonChain.Participants) == 0 {
		return nil, fmt.Errorf("canton chain %d not found", chainSelector)
	}
	addr, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return nil, fmt.Errorf("failed to get address for %v on chain %d: %w", ref, chainSelector, err)
	}
	instanceAddress := contracts.HexToInstanceAddress(addr.Address)
	participant := cantonChain.Participants[0]
	ctx := context.Background()
	if e.GetContext != nil {
		ctx = e.GetContext()
	}
	activePool, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		instanceAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("resolve lock/release token pool at %s: %w", addr.Address, err)
	}
	pool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return nil, fmt.Errorf("parse lock/release token pool at %s: %w", addr.Address, err)
	}
	instrumentCombined := string(pool.InstrumentId.Id) + "@" + string(pool.InstrumentId.Admin)
	return crypto.Keccak256([]byte(instrumentCombined)), nil
}

// CantonAdapter is an implementation of the ChainFamily interface for Canton.
type CantonAdapter struct {
	base adapters.ChainFamily
}

// NewChainFamilyAdapter creates a new Canton chain family adapter.
// A "base" adapter needs to be passed in, currently assumed to be the EVM chain family adapter,
// in order to achieve all functionality.
// TODO: this needs to be fully implemented for Canton.
func NewChainFamilyAdapter(base adapters.ChainFamily) *CantonAdapter {
	return &CantonAdapter{
		base: base,
	}
}

// AddressRefToBytes implements adapters.ChainFamily.
func (c *CantonAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	// Canton "addresses" are "InstanceAddresses", which are the 32-byte keccak256 hash of the RawInstanceAddress.
	// The RawInstanceAddress is of the form:
	// <instance-id>@<party-id>
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

// ConfigureChainForLanes implements adapters.ChainFamily.
func (c *CantonAdapter) ConfigureChainForLanes() *operations.Sequence[adapters.ConfigureChainForLanesInput, sequences.OnChainOutput, chain.BlockChains] {
	return c.base.ConfigureChainForLanes()
}
