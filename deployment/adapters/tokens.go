package adapters

import (
	"github.com/Masterminds/semver/v3"

	tokenadapters "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

var (
	_ tokenadapters.TokenAdapter = &CantonTokenAdapter{}
)

type CantonTokenAdapter struct{}

// No-op sequence so devenv ConfigureTokensForTransfers succeeds on Canton without sending EVM-style
// token pool transactions to the ledger (Canton uses DAML contracts; message-only lanes do not need this).
var cantonConfigureTokenForTransfersNoOp = operations.NewSequence(
	"canton/devenv/configure_token_for_transfers_noop",
	semver.MustParse("0.0.1"),
	"Canton devenv: skip EVM token pool configure sequence; pools are not driven by the shared EVM path",
	func(_ operations.Bundle, _ chain.BlockChains, _ tokenadapters.ConfigureTokenForTransfersInput) (sequences.OnChainOutput, error) {
		return sequences.OnChainOutput{}, nil
	},
)

func (c CantonTokenAdapter) ConfigureTokenForTransfersSequence() *operations.Sequence[tokenadapters.ConfigureTokenForTransfersInput, sequences.OnChainOutput, chain.BlockChains] {
	return cantonConfigureTokenForTransfersNoOp
}

func (c CantonTokenAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	return contracts.HexToInstanceAddress(ref.Address).Bytes(), nil
}

func (c CantonTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef, token []byte) (uint8, error) {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) DeriveTokenPoolCounterpart(e deployment.Environment, chainSelector uint64, tokenPool []byte, token []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) ManualRegistration() *operations.Sequence[tokenadapters.ManualRegistrationSequenceInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) SetTokenPoolRateLimits() *operations.Sequence[tokenadapters.TPRLRemotes, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
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
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) UpdateAuthorities() *operations.Sequence[tokenadapters.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	// TODO implement me
	panic("implement me")
}

func (c CantonTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *operations.Sequence[tokenadapters.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, chain.BlockChains] {
	// TODO implement me
	panic("implement me")
}
