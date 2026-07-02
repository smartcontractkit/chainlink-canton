package evmops

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/factory_burn_mint_erc20"

	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/clients"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/input"
)

func EnsureERC20Allowance(
	ctx context.Context,
	b *clients.Bundle,
	token common.Address,
	spender common.Address,
	amount *big.Int,
) error {
	erc20, err := factory_burn_mint_erc20.NewFactoryBurnMintERC20(token, b.ETHClient)
	if err != nil {
		return err
	}

	balance, err := erc20.BalanceOf(&bind.CallOpts{Context: ctx}, b.EthAuth.From)
	if err != nil {
		return err
	}
	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance: got: %s, want at least: %s", balance.String(), amount.String())
	}

	allowance, err := erc20.Allowance(&bind.CallOpts{Context: ctx}, b.EthAuth.From, spender)
	if err != nil {
		return err
	}
	if allowance.Cmp(amount) >= 0 {
		fmt.Printf("Sufficient allowance: %s\n", allowance.String())
		return nil
	}

	// Ask for confirmation
	fmt.Printf("Insufficient ERC20 allowance (%s) for spender %s on token %s\n", allowance.String(), spender.Hex(), token.Hex())
	fmt.Printf("Approving %s to spend %s...\n", spender.Hex(), amount.String())
	fmt.Println("Confirm? (Y/N)")
	if !input.Confirm() {
		return fmt.Errorf("cancel approval")
	}

	fmt.Println("⏳ Approving...")
	opts := *b.EthAuth
	opts.Context = ctx
	tx, err := erc20.Approve(&opts, spender, amount)
	if err != nil {
		return fmt.Errorf("approve ERC20: %w", err)
	}
	fmt.Printf("Sent message with tx hash: %s\n", tx.Hash().Hex())
	fmt.Println(b.EVMExplorerLink(tx.Hash().Hex()))

	fmt.Println("⏳ Waiting for transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, b.ETHClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	fmt.Printf("Transaction mined in block: %d\n", receipt.BlockNumber.Uint64())

	return nil
}
