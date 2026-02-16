package lock_release_token_pool

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonLockReleaseTokenPool")

var Version = semver.MustParse("1.0.0")

var Deploy = contract.NewDeploy(contract.DeployParams[lockreleasetokenpool.LockReleaseTokenPool]{
	Name:           "canton/ccip/lock_release_token_pool/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the CCIP LockReleaseTokenPool contract on Canton",
	Validate: func(template lockreleasetokenpool.LockReleaseTokenPool) error {
		if template.InstrumentId == (splice_api_token_holding_v1.InstrumentId{}) {
			return errors.New("instrument ID cannot be empty")
		}
		if template.Decimals < 0 {
			return errors.New("decimals cannot be negative")
		}

		return nil
	},
	PackageName: string(contracts.CCIPLockReleaseTokenPool),
	Prefix:      "lockreleasetokenpool",
})
