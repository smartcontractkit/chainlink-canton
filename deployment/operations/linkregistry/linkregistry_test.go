package linkregistry

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_burn_mint_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
)

func TestExerciseLinkRegistryUsesSpliceInterfaceTemplateIDs(t *testing.T) {
	t.Parallel()

	burnMintCmd := exerciseLinkRegistryBurnMint("cid", splice_api_token_burn_mint_v1.BurnMintFactoryBurnMint{})
	require.Equal(t, splice_api_token_burn_mint_v1.IBurnMintFactoryInterfaceID(), burnMintCmd.TemplateID)
	require.Equal(t, "BurnMintFactory_BurnMint", burnMintCmd.Choice)

	transferCmd := exerciseLinkRegistryTransfer("cid", splice_api_token_transfer_instruction_v1.TransferFactoryTransfer{})
	require.Equal(t, splice_api_token_transfer_instruction_v1.ITransferFactoryInterfaceID(), transferCmd.TemplateID)
	require.Equal(t, "TransferFactory_Transfer", transferCmd.Choice)
}
