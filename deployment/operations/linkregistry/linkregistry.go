package linkregistry

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/go-daml/pkg/model"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/link"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_burn_mint_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
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
	Method:       exerciseLinkRegistryBurnMint,
})

var TransferFactory_Transfer = contract.NewExercise(contract.ExerciseParams[splice_api_token_transfer_instruction_v1.TransferFactoryTransfer]{
	Name:         "canton/link_registry/transfer",
	Version:      Version,
	Description:  "Transfers Holdings on Canton",
	ContractType: ContractType,
	Template:     link.LinkRegistry{},
	Method:       exerciseLinkRegistryTransfer,
})

var CreateTransferPreapproval = contract.NewExercise(contract.ExerciseParams[link.CreateTransferPreapproval]{
	Name:         "canton/link_registry/create_transfer_preapproval",
	Version:      Version,
	Description:  "Creates a transfer preapproval for the receiver via LinkRegistry",
	ContractType: ContractType,
	Template:     link.LinkRegistry{},
	Method:       link.LinkRegistry{}.CreateTransferPreapproval,
})

// Splice token interface choices must be exercised via the interface package on the ledger
// (see integration-tests/ccip/ccip_send_with_token_bnm_test.go), not Link.Token's local view.
func exerciseLinkRegistryBurnMint(contractID string, args splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint) *model.ExerciseCommand {
	cmd := link.LinkRegistry{}.BurnMintFactoryBurnMint(contractID, args)
	cmd.TemplateID = splice_api_token_burn_mint_v1.IBurnMintFactoryInterfaceID()

	return cmd
}

func exerciseLinkRegistryTransfer(contractID string, args splice_api_token_transfer_instruction_v1.TransferFactoryTransfer) *model.ExerciseCommand {
	cmd := link.LinkRegistry{}.TransferFactoryTransfer(contractID, args)
	cmd.TemplateID = splice_api_token_transfer_instruction_v1.ITransferFactoryInterfaceID()

	return cmd
}
