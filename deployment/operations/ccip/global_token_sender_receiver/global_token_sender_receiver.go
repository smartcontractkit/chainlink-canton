package global_token_sender_receiver

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	globalbindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/globaltokensenderreceiver"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonGlobalTokenSenderReceiver")

var Version = semver.MustParse("2.0.0")

var Deploy = contract.NewDeploy(contract.DeployParams[globalbindings.GlobalTokenSenderReceiver]{
	Name:           "canton/ccip/global_token_sender_receiver/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP GlobalTokenSenderReceiver contract on Canton",
	Validate: func(template globalbindings.GlobalTokenSenderReceiver) error {
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPGlobalTokenSenderReceiverV1),
	Prefix:      "global-token-sender-receiver",
})
