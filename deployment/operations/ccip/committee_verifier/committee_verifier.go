package committee_verifier

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CommitteeVerifier")

var Version = semver.MustParse("0.1.0")

var Deploy = contract.NewDeploy(contract.DeployParams[ccvs.CommitteeVerifier]{
	Name:           "canton/ccip/committee_verifier/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP CommitteeVerifier contract on Canton",
	Validate: func(template ccvs.CommitteeVerifier) error {
		if template.Owner == "" {
			return errors.New("owner cannot be empty")
		}
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}
		if template.VersionTag == "" {
			return errors.New("version tag cannot be empty")
		}
		if template.MessageSentObserver == "" {
			return errors.New("message sent observer cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPCommitteeVerifier),
	Prefix:      "committeeverifier",
})
