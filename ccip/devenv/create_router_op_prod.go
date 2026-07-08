//go:build prodledger

package devenv

import (
	"errors"

	rt "github.com/smartcontractkit/chainlink-canton/bindings/generated/v1_0_0/ccip/ccipruntime"
	pprof "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var prodCreateRouterEncoder = rt.NewContract("", "CCIP.PerPartyRouter", "PerPartyRouterFactory").Encoder()

var createRouterOperation = contract.NewExercise(contract.ExerciseParams[rt.CreateRouter]{
	Name:         "canton/ccip/per_party_router_factory/create_router",
	Version:      pprof.Version,
	Description:  "Creates a new PerPartyRouter using the PerPartyRouterFactory",
	ContractType: pprof.ContractType,
	Validate: func(input rt.CreateRouter) error {
		if input.InstanceId == "" {
			return errors.New("instance ID cannot be empty")
		}
		if input.PartyOwner == "" {
			return errors.New("router owner cannot be empty")
		}

		return nil
	},
	Template:     rt.PerPartyRouterFactory{},
	Method:       rt.PerPartyRouterFactory{}.CreateRouter,
	EncodeMethod: prodCreateRouterEncoder.CreateRouter,
})
