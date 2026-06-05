package per_party_router_factory

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonPerPartyRouterFactory")

var Version = semver.MustParse("0.1.0")

var pprfEncoder = ccipruntime.NewContract("", "CCIP.PerPartyRouter", "PerPartyRouterFactory").Encoder()

var Deploy = contract.NewDeploy(contract.DeployParams[ccipruntime.PerPartyRouterFactory]{
	Name:           "canton/ccip/per_party_router_factory/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys a CCIP PerPartyRouterFactory contract on Canton",
	Validate: func(template ccipruntime.PerPartyRouterFactory) error {
		if template.CcipOwner == "" {
			return errors.New("ccip owner cannot be empty")
		}

		return nil
	},
	PackageName: string(contracts.CCIPPerPartyRouter),
	Prefix:      "perpartyrouterfactory",
})

var CreateRouter = contract.NewExercise(contract.ExerciseParams[ccipruntime.CreateRouter]{
	Name:         "canton/ccip/per_party_router_factory/create_router",
	Version:      Version,
	Description:  "Creates a new PerPartyRouter using the PerPartyRouterFactory",
	ContractType: ContractType,
	Validate: func(input ccipruntime.CreateRouter) error {
		if input.InstanceId == "" {
			return errors.New("instance ID cannot be empty")
		}
		if input.PartyOwner == "" {
			return errors.New("router owner cannot be empty")
		}

		return nil
	},
	Template:     ccipruntime.PerPartyRouterFactory{},
	Method:       ccipruntime.PerPartyRouterFactory{}.CreateRouter,
	EncodeMethod: pprfEncoder.CreateRouter,
})
