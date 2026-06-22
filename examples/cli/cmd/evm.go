package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/offramp"
	routerwrapper "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v2_0_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/evm"
	"github.com/smartcontractkit/chainlink-ccv/protocol"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/cantonops"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/clients"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/evmops"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/input"
)

// parseAmount parses an amount string that may include exponents (e.g. "1e18")
// and returns it as a big.Int.
func parseAmount(s string) (*big.Int, error) {
	f := new(big.Float)
	f.SetPrec(256)
	_, ok := f.SetString(s)
	if !ok {
		return nil, fmt.Errorf("invalid amount %q", s)
	}
	i := new(big.Int)
	f.Int(i)

	return i, nil
}

// NewEVMCmd returns the `evm` parent command.
func NewEVMCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "evm",
		Short: "EVM-side CCIP operations (Ethereum)",
	}
	c.AddCommand(newEVMSendMessageCmd(g))
	c.AddCommand(newEVMSendTokenCmd(g))
	c.AddCommand(newEVMExecuteCmd(g))

	return c
}

// resolveEthFeeToken maps the --fee-token flag to a token address.
// "link" → profile.EthLinkTokenAddress, "native" → zero address.
func resolveEthFeeToken(b *clients.Bundle, name string) (common.Address, error) {
	switch name {
	case "link":
		return b.Profile.EthLinkTokenAddress, nil
	case "native":
		return b.Profile.EthEmptyAddress, nil
	default:
		return common.Address{}, fmt.Errorf("invalid --fee-token %q (link|native)", name)
	}
}

func newEVMSendMessageCmd(g *Globals) *cobra.Command {
	var (
		receiverParty string
		payload       string
		feeTokenName  string
	)
	c := &cobra.Command{
		Use:   "send-message",
		Short: "Send a message-only CCIP message from EVM to Canton",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}
			feeToken, err := resolveEthFeeToken(b, feeTokenName)
			if err != nil {
				return err
			}
			receiverPartyHashed := contracts.HashedPartyFromString(receiverParty)
			if receiverParty == "" {
				receiverPartyHashed = contracts.HashedPartyFromString(b.Participant.PartyID)
			}

			return evmSend(ctx, b, receiverPartyHashed, []byte(payload), feeToken, nil)
		},
	}
	c.Flags().StringVar(&receiverParty, "receiver-party", "", "destination Canton party id (defaults to own party)")
	c.Flags().StringVar(&payload, "payload", "Hello, Canton!", "message payload (text)")
	c.Flags().StringVar(&feeTokenName, "fee-token", "link", "fee token (link|native)")

	return c
}

func newEVMSendTokenCmd(g *Globals) *cobra.Command {
	var (
		receiverParty string
		amountStr     string
		feeTokenName  string
	)
	c := &cobra.Command{
		Use:   "send-token",
		Short: "Send a LINK token transfer CCIP message from EVM to Canton",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}
			amount, err := parseAmount(amountStr)
			if err != nil {
				return fmt.Errorf("invalid --amount %q (support decimal syntax and exponents, e.g. 1e18): %w", amountStr, err)
			}
			feeToken, err := resolveEthFeeToken(b, feeTokenName)
			if err != nil {
				return err
			}
			tokenAmounts := []routerwrapper.ClientEVMTokenAmount{{
				Token:  b.Profile.EthTokenAddress,
				Amount: amount,
			}}
			receiverPartyHashed := contracts.HashedPartyFromString(receiverParty)
			if receiverParty == "" {
				receiverPartyHashed = contracts.HashedPartyFromString(b.Participant.PartyID)
			}

			return evmSend(ctx, b, receiverPartyHashed, nil, feeToken, tokenAmounts)
		},
	}
	c.Flags().StringVar(&receiverParty, "receiver-party", "", "destination Canton party id (defaults to own party)")
	c.Flags().StringVar(&amountStr, "amount", "", "token amount in wei (required; supports exponents, e.g. 1e18)")
	c.Flags().StringVar(&feeTokenName, "fee-token", "link", "fee token (link|native)")
	_ = c.MarkFlagRequired("amount")

	return c
}

// evmSend is shared by send-message and send-token. It always uses the
// no-execution tag as executor because no executor currently supports Canton
// as a destination.
func evmSend(
	ctx context.Context,
	b *clients.Bundle,
	receiverParty contracts.HashedParty,
	payload []byte,
	feeToken common.Address,
	tokens []routerwrapper.ClientEVMTokenAmount,
) error {
	router, err := routerwrapper.NewRouter(b.Profile.EthRouterAddress, b.ETHClient)
	if err != nil {
		return fmt.Errorf("bind router: %w", err)
	}

	extraArgs, err := evm.NewV3ExtraArgs(
		protocol.FinalityWaitForFinality,
		uint32(0),
		b.Profile.NoExecutionTag.Hex(),
		nil, nil, nil, nil,
	)
	if err != nil {
		return fmt.Errorf("build extraArgs: %w", err)
	}

	msg := routerwrapper.ClientEVM2AnyMessage{
		Receiver:     receiverParty.Bytes(),
		Data:         payload,
		TokenAmounts: tokens,
		FeeToken:     feeToken,
		ExtraArgs:    extraArgs,
	}

	fee, err := router.GetFee(&bind.CallOpts{Context: ctx}, b.Profile.CantonSelector, msg)
	if err != nil {
		return fmt.Errorf("get fee: %w", err)
	}
	fmt.Printf("Fee for sending message: %s\n", fee.String())

	requiredAllowances := make(map[common.Address]*big.Int)
	for _, token := range tokens {
		allowance := requiredAllowances[token.Token]
		if allowance == nil {
			allowance = big.NewInt(0)
		}
		allowance.Add(allowance, token.Amount)
		requiredAllowances[token.Token] = allowance
	}

	auth := b.EthAuth
	if feeToken == b.Profile.EthEmptyAddress {
		auth.Value = fee
	} else {
		allowance := requiredAllowances[feeToken]
		if allowance == nil {
			allowance = big.NewInt(0)
		}
		allowance.Add(allowance, fee)
		requiredAllowances[feeToken] = allowance
	}

	for address, requiredAllowance := range requiredAllowances {
		if requiredAllowance.Cmp(big.NewInt(0)) > 0 {
			fmt.Printf("Ensuring router has sufficient allowance on token %s: %s\n", address.Hex(), requiredAllowance.String())
			err := evmops.EnsureERC20Allowance(ctx, b, address, b.Profile.EthRouterAddress, requiredAllowance)
			if err != nil {
				return fmt.Errorf("ensure ERC20 allowance: %w", err)
			}
		}
	}

	// Ask for confirmation
	fmt.Println("About to send message:")
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.AppendHeader(table.Row{"Field", "Value"})
	tw.AppendRow(table.Row{"Receiver", "0x" + common.Bytes2Hex(msg.Receiver)})
	tw.AppendRow(table.Row{"Data", "0x" + common.Bytes2Hex(msg.Data)})
	tw.AppendRow(table.Row{"Fee Token", msg.FeeToken.Hex()})
	tw.AppendRow(table.Row{"Tx Value", auth.Value.String()})
	if len(msg.TokenAmounts) > 0 {
		tw.AppendRow(table.Row{"Token Transfer - Token", msg.TokenAmounts[0].Token.Hex()})
		tw.AppendRow(table.Row{"Token Transfer - Amount", msg.TokenAmounts[0].Amount.String()})
	}
	fmt.Println(tw.Render())
	fmt.Println("Confirm? (Y/N)")
	if !input.Confirm() {
		return fmt.Errorf("cancel sending message")
	}

	fmt.Println("⏳ Sending message...")
	tx, err := router.CcipSend(auth, b.Profile.CantonSelector, msg)
	if err != nil {
		return fmt.Errorf("ccip send: %w", err)
	}
	fmt.Printf("Sent message with tx hash: %s\n", tx.Hash().Hex())
	fmt.Println(b.EVMExplorerLink(tx.Hash().Hex()))

	fmt.Println("⏳ Waiting for transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, b.ETHClient, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	fmt.Printf("Transaction mined in block: %d\n", receipt.BlockNumber.Uint64())

	eventTopic := (onramp.OnRampCCIPMessageSent{}).Topic()
	for _, lg := range receipt.Logs {
		if len(lg.Topics) == 0 || lg.Topics[0] != eventTopic {
			continue
		}
		onRamp, err := onramp.NewOnRamp(lg.Address, b.ETHClient)
		if err != nil {
			return fmt.Errorf("bind onramp: %w", err)
		}
		sent, err := onRamp.ParseCCIPMessageSent(*lg)
		if err != nil {
			return fmt.Errorf("parse CCIPMessageSent: %w", err)
		}
		fmt.Printf("✅ Message sent with ID: %s\n", common.BytesToHash(sent.MessageId[:]))
		fmt.Println(b.CCIPExplorerLink(common.BytesToHash(sent.MessageId[:]).Hex()))

		return nil
	}

	return fmt.Errorf("CCIPMessageSent event not found in transaction logs")
}

func newEVMExecuteCmd(g *Globals) *cobra.Command {
	var (
		messageIDHex string
		wait         time.Duration
	)
	c := &cobra.Command{
		Use:   "execute",
		Short: "Execute on EVM a message sent from Canton (used with --executor none)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}
			messageId := common.HexToHash(messageIDHex)

			fmt.Printf("Waiting for verifier results for %s (timeout %s)...\n", messageId.Hex(), wait)
			resp, err := cantonops.WaitForVerifierResult(ctx, b.IndexerClient, messageId.Hex(), wait)
			if err != nil {
				return err
			}
			fmt.Printf("Verifier results for %s successfully retrieved.\n", messageId.Hex())

			verifierResult := resp.Results[0].VerifierResult
			encodedMessage, err := verifierResult.Message.Encode()
			if err != nil {
				return fmt.Errorf("encode message: %w", err)
			}
			offRampAddress := common.BytesToAddress(verifierResult.Message.OffRampAddress)

			offRamp, err := offramp.NewOffRamp(offRampAddress, b.ETHClient)
			if err != nil {
				return fmt.Errorf("bind offramp: %w", err)
			}

			execState, err := offRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, messageId)
			if err != nil {
				return fmt.Errorf("get execution state: %w", err)
			}
			if ccip_offramp.MessageExecutionState(execState) == ccip_offramp.Success_MessageExecutionState {
				fmt.Println("Message already executed")
				return nil
			}

			fmt.Println("⏳ Executing message...")
			tx, err := offRamp.Execute(
				b.EthAuth,
				encodedMessage,
				[]common.Address{common.BytesToAddress(verifierResult.VerifierDestAddress)},
				[][]byte{verifierResult.CCVData},
				0,
			)
			if err != nil {
				return fmt.Errorf("execute: %w", err)
			}
			fmt.Printf("Message executed in transaction: %s\n", tx.Hash().Hex())
			fmt.Println(b.EVMExplorerLink(tx.Hash().Hex()))

			fmt.Println("⏳ Waiting for transaction to be mined...")
			receipt, err := bind.WaitMined(ctx, b.ETHClient, tx)
			if err != nil {
				return fmt.Errorf("wait mined: %w", err)
			}
			fmt.Printf("Transaction mined in block: %d\n", receipt.BlockNumber.Uint64())
			fmt.Printf("✅ Message %s executed on EVM\n", messageId.Hex())
			fmt.Println(b.CCIPExplorerLink(messageId.Hex()))

			return nil
		},
	}
	c.Flags().StringVar(&messageIDHex, "message-id", "", "CCIP message id (0x-prefixed hex) (required)")
	c.Flags().DurationVar(&wait, "wait", 15*time.Minute, "max time to wait for verifier results")
	_ = c.MarkFlagRequired("message-id")

	return c
}
