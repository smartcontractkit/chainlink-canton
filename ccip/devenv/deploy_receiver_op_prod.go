//go:build prodledger

package devenv

import (
	"errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	v1receiver "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	receiverop "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ccipReceiverDeployOperation = contract.NewDeploy(contract.DeployParams[v1receiver.CCIPReceiver]{
	Name:           "canton/ccip/receiver/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(receiverop.ContractType, *receiverop.Version),
	Description:    "Deploys a CCIP Receiver contract on Canton",
	Validate: func(template v1receiver.CCIPReceiver) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPReceiver),
	Prefix:      "ccipreceiver",
})
