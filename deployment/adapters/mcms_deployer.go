package adapters

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cantonsequences "github.com/smartcontractkit/chainlink-canton/deployment/sequences"
)

var _ ccipdeploy.Deployer = (*CantonMCMSDeployer)(nil)

// CantonMCMSDeployer is the Canton MCMS deployer. The timelock is built into the MCMS
// contract, so FinalizeDeployMCMS is a no-op.
type CantonMCMSDeployer struct{}

func (d *CantonMCMSDeployer) DeployChainContracts() *cldf_ops.Sequence[ccipdeploy.ContractDeploymentConfigPerChainWithAddress, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return unsupportedSequence[ccipdeploy.ContractDeploymentConfigPerChainWithAddress](
		"canton-seq-deploy-chain-contracts-unsupported",
		"DeployChainContracts is 1.6-only and not supported on Canton",
	)
}

func (d *CantonMCMSDeployer) DeployMCMS() *cldf_ops.Sequence[ccipdeploy.MCMSDeploymentConfigPerChainWithAddress, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return unsupportedSequence[ccipdeploy.MCMSDeploymentConfigPerChainWithAddress](
		"canton-seq-deploy-mcms-unsupported",
		"DeployMCMS is not supported on Canton yet: the input has no integer Daml chain ID. Use the DeployAndConfigureMCMS changeset",
	)
}

func (d *CantonMCMSDeployer) FinalizeDeployMCMS() *cldf_ops.Sequence[ccipdeploy.MCMSDeploymentConfigPerChainWithAddress, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"canton-seq-finalize-deploy-mcms",
		semver.MustParse("1.0.0"),
		"No-op: the Canton MCMS timelock is built into the MCMS contract and needs no separate initialization.",
		func(_ cldf_ops.Bundle, chains cldf_chain.BlockChains, in ccipdeploy.MCMSDeploymentConfigPerChainWithAddress) (sequences.OnChainOutput, error) {
			if _, ok := chains.CantonChains()[in.ChainSelector]; !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("canton chain with selector %d not found in environment", in.ChainSelector)
			}

			return sequences.OnChainOutput{}, nil
		},
	)
}

func (d *CantonMCMSDeployer) SetOCR3Config() *cldf_ops.Sequence[ccipdeploy.SetOCR3ConfigInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return unsupportedSequence[ccipdeploy.SetOCR3ConfigInput](
		"canton-seq-set-ocr3-config-unsupported",
		"SetOCR3Config is 1.6-only and not supported on Canton",
	)
}

func (d *CantonMCMSDeployer) GrantAdminRoleToTimelock() *cldf_ops.Sequence[ccipdeploy.GrantAdminRoleToTimelockConfigPerChainWithSelector, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return unsupportedSequence[ccipdeploy.GrantAdminRoleToTimelockConfigPerChainWithSelector](
		"canton-seq-grant-admin-role-to-timelock-unsupported",
		"GrantAdminRoleToTimelock is not supported on Canton: the MCMS has no admin role, only its owner party",
	)
}

func (d *CantonMCMSDeployer) UpdateMCMSConfig() *cldf_ops.Sequence[ccipdeploy.UpdateMCMSConfigInputPerChainWithSelector, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cantonsequences.UpdateMCMSConfig
}

func unsupportedSequence[IN any](id, reason string) *cldf_ops.Sequence[IN, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		id,
		semver.MustParse("1.0.0"),
		reason,
		func(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ IN) (sequences.OnChainOutput, error) {
			return sequences.OnChainOutput{}, errors.New(reason)
		},
	)
}
