package adapters

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
)

var _ ccipdeploy.TransferOwnershipAdapter = (*CantonTransferOwnershipAdapter)(nil)

// CantonTransferOwnershipAdapter satisfies the cross-family deploy-chain-contracts
// timelock transfer hook. Canton CCIP contracts record mcms/ccipOwner on-ledger at
// deploy time (factory choices / proposal-driven deploy); they do not use the EVM
// Ownable2Step + timelock accept path.
type CantonTransferOwnershipAdapter struct{}

func (a *CantonTransferOwnershipAdapter) InitializeTimelockAddress(_ cldf.Environment, _ mcms.Input) error {
	return nil
}

func (a *CantonTransferOwnershipAdapter) SequenceTransferOwnershipViaMCMS() *cldf_ops.Sequence[ccipdeploy.TransferOwnershipPerChainInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton-seq-transfer-ownership-via-mcms",
		semver.MustParse("1.0.0"),
		"No-op: Canton ownership for CCIP contracts is handled in Canton deploy sequences / MCMS proposals, not via the EVM timelock transfer path.",
		func(_ cldf_ops.Bundle, chains cldf_chain.BlockChains, in ccipdeploy.TransferOwnershipPerChainInput) (sequences.OnChainOutput, error) {
			if _, ok := chains.CantonChains()[in.ChainSelector]; !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found in environment", in.ChainSelector)
			}
			return sequences.OnChainOutput{}, nil
		},
	)
}

func (a *CantonTransferOwnershipAdapter) SequenceAcceptOwnership() *cldf_ops.Sequence[ccipdeploy.TransferOwnershipPerChainInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton-seq-accept-ownership",
		semver.MustParse("1.0.0"),
		"No-op: Canton contracts do not require a separate accept-ownership step.",
		func(_ cldf_ops.Bundle, chains cldf_chain.BlockChains, in ccipdeploy.TransferOwnershipPerChainInput) (sequences.OnChainOutput, error) {
			if _, ok := chains.CantonChains()[in.ChainSelector]; !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found in environment", in.ChainSelector)
			}
			return sequences.OnChainOutput{}, nil
		},
	)
}

func (a *CantonTransferOwnershipAdapter) ShouldAcceptOwnershipWithTransferOwnership(_ cldf.Environment, _ ccipdeploy.TransferOwnershipPerChainInput) (bool, error) {
	return false, nil
}
