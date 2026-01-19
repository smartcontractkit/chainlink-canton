package changesets

import (
	"fmt"

	linkops "github.com/smartcontractkit/chainlink-canton-internal/deployment/ops/link"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type DeployLinkTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

// DeployLinkToken deploys Sui chain packages and modules
type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	suiChains := e.BlockChains.CantonChains()

	deps := sui_ops.OpTxDeps{}

	// Run DeployLinkToken Operation
	_, err := cld_ops.ExecuteOperation(e.OperationsBundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// TODO: figure out how to handle storage of contracts
	// save LinkToken address to the addressbook
	// typeAndVersionLinkToken := cldf.NewTypeAndVersion(deployment.SuiLinkTokenType, deployment.Version1_0_0)
	// err = ab.Save(config.ChainSelector, linkTokenReport.Output.PackageId, typeAndVersionLinkToken)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken address %s for Sui chain %d: %w", linkTokenReport.Output.PackageId, config.ChainSelector, err)
	// }

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployLinkToken) VerifyPreconditions(e cldf.Environment, config DeployLinkTokenConfig) error {
	return nil
}
