package linkregistry

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/link"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_burn_mint_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

var ContractType = deployment.ContractType("LinkRegistry")

var Version = semver.MustParse("2.0.0")

var Deploy = contract.NewDeploy(contract.DeployParams[link.LinkRegistry]{
	Name:           "canton/link_registry/deploy",
	TypeAndVersion: deployment.NewTypeAndVersion(ContractType, *Version),
	Description:    "Deploys the LinkRegistry contract on Canton",
	Validate: func(template link.LinkRegistry) error {

		return nil
	},
	PackageName: string(contracts.Link),
	Prefix:      "link-registry",
})

var BurnMintFactory_BurnMint = contract.NewExercise(contract.ExerciseParams[splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint]{
	Name:         "canton/link_registry/burn_mint",
	Version:      Version,
	Description:  "Burns/Mints Holdings on Canton",
	ContractType: ContractType,
	Template:     link.LinkRegistry{},
	Method:       link.LinkRegistry{}.BurnMintFactoryBurnMint,
})

var TransferFactory_Transfer = contract.NewExercise(contract.ExerciseParams[splice_api_token_transfer_instruction_v1.TransferFactoryTransfer]{
	Name:         "canton/link_registry/transfer",
	Version:      Version,
	Description:  "Transfers Holdings on Canton",
	ContractType: ContractType,
	Template:     link.LinkRegistry{},
	Method:       link.LinkRegistry{}.TransferFactoryTransfer,
})

var CreateTransferPreapproval = contract.NewExercise(contract.ExerciseParams[link.CreateTransferPreapproval]{
	Name:         "canton/link_registry/create_transfer_preapproval",
	Version:      Version,
	Description:  "Creates a transfer preapproval for the receiver via LinkRegistry",
	ContractType: ContractType,
	Template:     link.LinkRegistry{},
	Method:       link.LinkRegistry{}.CreateTransferPreapproval,
})
