package registry_burn_mint_token_pool

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	registrybnm "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("CantonRegistryBurnMintTokenPool")

var Version = semver.MustParse("2.0.0")

var poolEncoder = registrybnm.NewContract("", "CCIP.Registry.BurnMintTokenPoolV2", "BurnMintTokenPool").Encoder()

var ApplyChainUpdates = contract.NewExercise(contract.ExerciseParams[registrybnm.ApplyChainUpdates]{
	Name:         "canton/ccip/registry_burn_mint_token_pool/apply_chain_updates",
	Version:      Version,
	Description:  "Applies remote chain updates to a Canton registry BurnMintTokenPool",
	ContractType: ContractType,
	Template:     registrybnm.BurnMintTokenPool{},
	Method:       registrybnm.BurnMintTokenPool{}.ApplyChainUpdates,
	EncodeMethod: poolEncoder.ApplyChainUpdates,
})
