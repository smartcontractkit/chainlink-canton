//go:build prodledger

package devenv

import (
	"errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	v1sender "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	senderop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ccipSenderDeployOperation = contract.NewDeploy(contract.DeployParams[v1sender.CCIPSender]{
	Name:           "canton/ccip/sender/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(senderop.ContractType, *senderop.Version),
	Description:    "Deploys a CCIP Sender contract on Canton",
	Validate: func(template v1sender.CCIPSender) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPSender),
	Prefix:      "ccipsender",
})
