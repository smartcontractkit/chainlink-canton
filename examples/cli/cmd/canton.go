package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_transfer_instruction_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/cantonops"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/clients"
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/input"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

const (
	defaultGasLimit = 50_000
)

// parseDecimalAmount parses an amount string that may include exponents
// (e.g. "1e-2") and returns it as a decimal string.
func parseDecimalAmount(s string) (string, error) {
	f := new(big.Float)
	f.SetPrec(256)
	_, ok := f.SetString(s)
	if !ok {
		return "", fmt.Errorf("invalid amount %q", s)
	}

	out := f.Text('f', -1) // convert to decimal string without exponent

	return out, nil
}

// NewCantonCmd returns the `canton` parent command.
func NewCantonCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "canton",
		Short: "Canton-side CCIP operations",
	}
	c.AddCommand(newCantonSendMessageCmd(g))
	c.AddCommand(newCantonSendTokenCmd(g))
	c.AddCommand(newCantonExecuteCmd(g))
	c.AddCommand(newCantonListEventsCmd(g))
	c.AddCommand(newCantonListHoldingsCmd(g))
	c.AddCommand(newCantonCreateTransferCmd(g))
	c.AddCommand(newCantonAcceptTransferCmd(g))

	return c
}

func resolveCantonFeeToken(b *clients.Bundle, name string) (*splice_api_token_holding_v1.InstrumentId, oapiTransferInstruction.ClientWithResponsesInterface, error) {
	switch name {
	case "link":
		return b.Profile.LinkInstrumentID, b.LinkTransferClient, nil
	case "native":
		return b.Profile.AmuletInstrumentID, b.AmuletTransferClient, nil
	default:
		return nil, nil, fmt.Errorf("invalid --fee-token %q (link|native)", name)
	}
}

// ---------------- canton list-events ----------------

func newCantonListEventsCmd(g *Globals) *cobra.Command {
	var eventName string
	c := &cobra.Command{
		Use:   "list-events",
		Short: "List active CCIPMessageSent or ExecutionStateChanged contracts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}
			var tmpl *apiv2.Identifier
			switch eventName {
			case "sent":
				tmpl = &apiv2.Identifier{PackageId: "#ccip-core", ModuleName: "CCIP.Events", EntityName: "CCIPMessageSent"}
			case "executed":
				tmpl = &apiv2.Identifier{PackageId: "#ccip-core", ModuleName: "CCIP.Events", EntityName: "ExecutionStateChanged"}
			default:
				return fmt.Errorf("invalid --event %q (sent|executed)", eventName)
			}
			active, err := testhelpers.ListActiveContractsByTemplateId(ctx, b.Participant, tmpl)
			if err != nil {
				return fmt.Errorf("list events: %w", err)
			}
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.Style().Title.Align = text.AlignCenter
			tw.SetAutoIndex(true)
			switch eventName {
			case "sent":
				tw.SetTitle("CCIPMessageSent Events")
				tw.AppendHeader(table.Row{"Message ID", "Destination Chain", "Sequence Number", "Sender"})
				for _, c := range active {
					ev, err := bindings.UnmarshalCreatedEvent[core.CCIPMessageSent](c.GetCreatedEvent())
					if err != nil {
						return fmt.Errorf("unmarshal: %w", err)
					}
					tw.AppendRow(table.Row{
						ev.Event.MessageId,
						ev.Event.DestChainSelector,
						ev.Event.SequenceNumber,
						ev.Sender,
					})
				}
			case "executed":
				tw.SetTitle("ExecutionStateChanged Events")
				tw.AppendHeader(table.Row{"Message ID", "Source Chain", "Sequence Number", "State", "Receiver"})
				for _, c := range active {
					ev, err := bindings.UnmarshalCreatedEvent[core.ExecutionStateChanged](c.GetCreatedEvent())
					if err != nil {
						return fmt.Errorf("unmarshal: %w", err)
					}
					tw.AppendRow(table.Row{
						ev.Event.MessageId,
						ev.Event.SourceChainSelector,
						ev.Event.SequenceNumber,
						ev.Event.State,
						ev.Receiver,
					})
				}
			}
			tw.AppendFooter(table.Row{"Total", len(active)})
			fmt.Println(tw.Render())

			return nil
		},
	}
	c.Flags().StringVar(&eventName, "event", "sent", "event to list (sent|executed)")

	return c
}

// ---------------- canton list-holdings ----------------

func newCantonListHoldingsCmd(g *Globals) *cobra.Command {
	var (
		withContractId bool
	)
	c := &cobra.Command{
		Use:   "list-holdings",
		Short: "List all token holdings for the configured party",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}
			holdings, err := testhelpers.ListActiveContractsByInterfaceId(ctx, b.Participant, &apiv2.Identifier{
				PackageId:  "#splice-api-token-holding-v1",
				ModuleName: "Splice.Api.Token.HoldingV1",
				EntityName: "Holding",
			})
			if err != nil {
				return fmt.Errorf("list holdings: %w", err)
			}
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.Style().Title.Align = text.AlignCenter
			tw.SetAutoIndex(true)
			if withContractId {
				tw.AppendHeader(table.Row{"Instrument ID", "Admin", "Owner", "Amount", "Contract ID"})
			} else {
				tw.AppendHeader(table.Row{"Instrument ID", "Admin", "Owner", "Amount"})
			}
			for _, h := range holdings {
				for _, view := range h.GetCreatedEvent().GetInterfaceViews() {
					var hv splice_api_token_holding_v1.HoldingView
					if err := ledger.RecordToStruct(view.GetViewValue(), &hv); err != nil {
						return fmt.Errorf("decode holding view: %w", err)
					}
					if withContractId {
						tw.AppendRow(table.Row{hv.InstrumentId.Id, hv.InstrumentId.Admin, hv.Owner, hv.Amount, h.GetCreatedEvent().GetContractId()})
					} else {
						tw.AppendRow(table.Row{hv.InstrumentId.Id, hv.InstrumentId.Admin, hv.Owner, hv.Amount})
					}
				}
			}
			fmt.Println(tw.Render())

			return nil
		},
	}
	c.Flags().BoolVar(&withContractId, "cid", false, "include contract IDs in the output")

	return c
}

// ---------------- canton create transfer ----------------
func newCantonCreateTransferCmd(g *Globals) *cobra.Command {
	var (
		tokenName        string
		receiverParty    string
		inputHoldingCids []string
		amount           string
	)
	c := &cobra.Command{
		Use:   "create-transfer",
		Short: "Create an outgoing TransferInstruction",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}

			instrumentId, transferInstructionClient, err := resolveCantonFeeToken(b, tokenName)
			if err != nil {
				return err
			}

			if receiverParty == "" {
				receiverParty = b.Participant.PartyID
			}

			// Resolve input holdings
			var inputHoldings []types.CONTRACT_ID
			if len(inputHoldingCids) > 0 {
				for _, cid := range inputHoldingCids {
					inputHoldings = append(inputHoldings, types.CONTRACT_ID(cid))
				}
			} else {
				holdings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, instrumentId)
				if err != nil {
					return fmt.Errorf("list holdings: %w", err)
				}
				for _, holding := range holdings {
					inputHoldings = append(inputHoldings, types.CONTRACT_ID(holding.ContractID))
				}
			}

			// Get disclosures
			transfer := splice_api_token_transfer_instruction_v1.Transfer{
				Sender:           types.PARTY(b.Participant.PartyID),
				Receiver:         types.PARTY(receiverParty),
				Amount:           types.NUMERIC(amount),
				InstrumentId:     *instrumentId,
				RequestedAt:      types.TIMESTAMP(time.Now()),
				ExecuteBefore:    types.TIMESTAMP(time.Now().Add(time.Hour * 24)),
				InputHoldingCids: inputHoldings,
				Meta:             splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
			}
			transferFactory, err := testhelpers.GetTransferFactoryV2(ctx, transferInstructionClient, string(instrumentId.Admin), transfer)
			if err != nil {
				return fmt.Errorf("get transfer factory: %w", err)
			}
			choiceContext, err := contracts.ChoiceContextFromData(transferFactory.ChoiceContextData)
			if err != nil {
				return fmt.Errorf("unmarshal choice context: %w", err)
			}

			// Create Transfer
			resp, err := b.Participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
				Commands: &apiv2.Commands{
					CommandId: uuid.NewString(),
					Commands: []*apiv2.Command{{
						Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferFactory"},
							ContractId: transferFactory.FactoryID,
							Choice:     "TransferFactory_Transfer",
							ChoiceArgument: ledger.MapToValue(splice_api_token_transfer_instruction_v1.TransferFactoryTransfer{
								ExpectedAdmin: instrumentId.Admin,
								Transfer:      transfer,
								ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
									Context: choiceContext,
								},
							}),
						}},
					}},
					ActAs:              []string{b.Participant.PartyID},
					DisclosedContracts: transferFactory.DisclosedContracts,
				},
				TransactionFormat: &apiv2.TransactionFormat{
					EventFormat: &apiv2.EventFormat{
						FiltersByParty: map[string]*apiv2.Filters{
							b.Participant.PartyID: {Cumulative: []*apiv2.CumulativeFilter{{IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{WildcardFilter: &apiv2.WildcardFilter{}}}}},
						},
						Verbose: true,
					},
					TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
				},
			})
			if err != nil {
				return fmt.Errorf("submit transfer: %w", err)
			}
			fmt.Printf("Submitted transfer in update: %s\n", resp.GetTransaction().GetUpdateId())
			fmt.Println(b.CantonExplorerLink(resp.GetTransaction().GetUpdateId()))

			return nil
		},
	}
	c.Flags().StringVar(&tokenName, "token", "link", "token to transfer (link|native)")
	c.Flags().StringVar(&receiverParty, "receiver", "", "party to receive the transfer (defaults to own party)")
	c.Flags().StringArrayVar(&inputHoldingCids, "input", nil, "the holding(s) to be used as an input for the transfer. If unspecified, all current holdings will be used.")
	c.Flags().StringVar(&amount, "amount", "", "the amount to transfer (required)")
	_ = c.MarkFlagRequired("amount")

	return c
}

func newCantonAcceptTransferCmd(g *Globals) *cobra.Command {
	var (
		contractID string
		tokenName  string
	)
	c := &cobra.Command{
		Use:   "accept-transfer",
		Short: "Accept an incoming TransferInstruction by contract ID",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}

			_, transferInstructionClient, err := resolveCantonFeeToken(b, tokenName)
			if err != nil {
				return err
			}

			return testhelpers.AcceptPendingTransferInstruction(
				ctx,
				b.Participant,
				transferInstructionClient,
				b.Participant.PartyID,
				contractID,
			)
		},
	}
	c.Flags().StringVar(&contractID, "contract-id", "", "TransferInstruction contract ID to accept (required)")
	c.Flags().StringVar(&tokenName, "token", "link", "token of the transfer instruction (link|native)")
	_ = c.MarkFlagRequired("contract-id")

	return c
}

// ---------------- canton execute ----------------

func newCantonExecuteCmd(g *Globals) *cobra.Command {
	var (
		messageIDHex string
		wait         time.Duration
	)
	c := &cobra.Command{
		Use:   "execute",
		Short: "Execute on Canton a message sent from EVM",
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

			return cantonExecute(ctx, b, resp.Results[0].VerifierResult)
		},
	}
	c.Flags().StringVar(&messageIDHex, "message-id", "", "CCIP message id (0x-prefixed hex) (required)")
	c.Flags().DurationVar(&wait, "wait", 15*time.Minute, "max time to wait for verifier results")
	_ = c.MarkFlagRequired("message-id")

	return c
}

func cantonExecute(ctx context.Context, b *clients.Bundle, vr protocol.VerifierResult) error {
	withToken := vr.Message.TokenTransfer != nil

	encodedMessage, err := vr.Message.Encode()
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}
	encodedHex := hex.EncodeToString(encodedMessage)
	verifierRawAddress, err := contracts.RawInstanceAddressFromString(string(vr.VerifierDestAddress))
	if err != nil {
		return fmt.Errorf("parse VerifierDestAddress: %w", err)
	}

	ccipExecuteDisclosure, err := eds.GetCCIPExecuteDisclosure(ctx, b.CCIPEDS, encodedHex)
	if err != nil {
		return fmt.Errorf("CCIP execute disclosure: %w", err)
	}
	ccvExecuteDisclosure, err := eds.GetCCVExecuteDisclosure(ctx, b.CCVEDS, encodedHex, verifierRawAddress.InstanceAddress())
	if err != nil {
		return fmt.Errorf("CCV execute disclosure: %w", err)
	}

	routerCid, err := cantonops.GetOrCreateRouter(ctx, b.Participant, b.CCIPEDS)
	if err != nil {
		return err
	}
	receiverCid, err := cantonops.GetOrCreateReceiver(ctx, b.Participant)
	if err != nil {
		return err
	}
	fmt.Printf("PerPartyRouter CID: %s\nCCIPReceiver CID: %s\n", routerCid, receiverCid)

	var tokenTransferInput *receiver.TokenTransferInput
	allDisclosures := slices.Concat(
		ccipExecuteDisclosure.DisclosedContracts,
		ccvExecuteDisclosure.DisclosedContracts,
	)

	if withToken {
		targetInstrumentId := contracts.BytesToEncodedInstrumentID(vr.Message.TokenTransfer.DestTokenAddress)
		tokenPoolAddress, err := eds.GetTokenPoolForToken(ctx, b.CCIPEDS, targetInstrumentId)
		if err != nil {
			return fmt.Errorf("get token pool: %w", err)
		}
		tokenPoolExecuteDisclosure, err := eds.GetTokenPoolExecuteDisclosure(ctx, b.TokenPoolEDS, encodedHex, tokenPoolAddress.InstanceAddress())
		if err != nil {
			return fmt.Errorf("token pool execute disclosure: %w", err)
		}
		tokenTransferInput = &receiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolExecuteDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(b.Participant.PartyID),
			PoolExtraContext:   tokenPoolExecuteDisclosure.ChoiceContext,
		}
		allDisclosures = slices.Concat(
			tokenPoolExecuteDisclosure.DisclosedContracts,
			allDisclosures,
		)
	}

	executeArgs := receiver.Execute{
		Context:        ccipExecuteDisclosure.ChoiceContext,
		RouterCid:      types.CONTRACT_ID(routerCid),
		EncodedMessage: types.TEXT(encodedHex),
		TokenTransfer:  tokenTransferInput,
		CcvInputs: []receiver.CCVInput{{
			CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
			VerifierResults: types.TEXT(hex.EncodeToString(vr.CCVData)),
			CcvExtraContext: ccvExecuteDisclosure.ChoiceContext,
		}},
	}

	fmt.Println("⏳ Executing message...")
	resp, err := b.Participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId:     receiverCid,
					Choice:         "Execute",
					ChoiceArgument: ledger.MapToValue(executeArgs),
				}},
			}},
			ActAs:              []string{b.Participant.PartyID},
			DisclosedContracts: allDisclosures,
		},
	})
	if err != nil {
		return fmt.Errorf("submit Execute: %w", err)
	}
	fmt.Printf("✅ Message executed in Update: %s\n", resp.GetTransaction().GetUpdateId())
	fmt.Println(b.CantonExplorerLink(resp.GetTransaction().GetUpdateId()))

	return nil
}

// ---------------- canton send-message / send-token ----------------

func newCantonSendMessageCmd(g *Globals) *cobra.Command {
	var (
		receiverHex  string
		gasLimit     int
		payload      string
		executor     string
		feeTokenName string
		feeInput     []string
	)
	c := &cobra.Command{
		Use:   "send-message",
		Short: "Send a message-only CCIP message from Canton to EVM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}

			if executor != "default" && executor != "none" {
				return fmt.Errorf("invalid --executor %q (default|none)", executor)
			}

			feeTokenInstrumentId, feeTokenTransferClient, err := resolveCantonFeeToken(b, feeTokenName)
			if err != nil {
				return err
			}

			msgReceiver := common.HexToAddress(receiverHex)
			if receiverHex == "" {
				msgReceiver = b.Profile.CCIPReceiverContract
			}

			return cantonSend(ctx, b, msgReceiver, []byte(payload), gasLimit, "", nil, executor, feeTokenInstrumentId, feeTokenTransferClient, feeInput)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", "destination EVM receiver address (0x-prefixed) (defaults to a CCIP Receiver contract)")
	c.Flags().StringVar(&payload, "payload", "Hello, EVM from Canton!", "message payload (text)")
	c.Flags().IntVar(&gasLimit, "gas-limit", -1, fmt.Sprintf("gas limit for EVM execution, defaults to %v for message transfers", defaultGasLimit))
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
	c.Flags().StringVar(&feeTokenName, "fee-token", "link", "fee token (link|native)")
	c.Flags().StringArrayVar(&feeInput, "fee-input", nil, "the holding(s) to be used as an input for the fee payment. If unspecified, all current holdings will be used.")

	return c
}

func newCantonSendTokenCmd(g *Globals) *cobra.Command {
	var (
		receiverHex  string
		gasLimit     int
		amountStr    string
		payload      string
		executor     string
		feeTokenName string
		feeInput     []string
		tokenInput   []string
	)
	c := &cobra.Command{
		Use:   "send-token",
		Short: "Send a LINK token transfer CCIP message from Canton to EVM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx)
			if err != nil {
				return err
			}

			if executor != "default" && executor != "none" {
				return fmt.Errorf("invalid --executor %q (default|none)", executor)
			}

			feeTokenInstrumentId, feeTokenTransferClient, err := resolveCantonFeeToken(b, feeTokenName)
			if err != nil {
				return err
			}

			msgReceiver := common.HexToAddress(receiverHex)
			if receiverHex == "" {
				msgReceiver = b.ETHAddress
			}

			return cantonSend(ctx, b, msgReceiver, []byte(payload), gasLimit, amountStr, tokenInput, executor, feeTokenInstrumentId, feeTokenTransferClient, feeInput)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", "destination EVM receiver address (0x-prefixed) (defaults to own address)")
	c.Flags().StringVar(&amountStr, "amount", "", "LINK token transfer amount as decimal (e.g. 0.12345, 1e-2); (required)")
	c.Flags().StringVar(&payload, "payload", "", "optional message payload (text) to attach to the token transfer")
	c.Flags().IntVar(&gasLimit, "gas-limit", -1, fmt.Sprintf("gas limit for EVM execution, defaults to %v for message transfers", defaultGasLimit))
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
	c.Flags().StringVar(&feeTokenName, "fee-token", "native", "fee token (link|native)")
	c.Flags().StringArrayVar(&feeInput, "fee-input", nil, "the holding(s) to be used as an input for the fee payment. If unspecified, all current holdings will be used.")
	c.Flags().StringArrayVar(&tokenInput, "token-input", nil, "the holding(s) to be used as an input for the token transfer. If unspecified, all current holdings will be used.")
	_ = c.MarkFlagRequired("amount")

	return c
}

// cantonSend implements Canton→EVM send for both message-only and
// token-transfer variants, with default or no-executor mode.
func cantonSend(
	ctx context.Context,
	b *clients.Bundle,
	receiver common.Address,
	payload []byte,
	gasLimit int,
	amountStr string,
	tokenInputHoldings []string,
	executorMode string,
	feeTokenInstrumentId *splice_api_token_holding_v1.InstrumentId,
	feeTokenTransferInstructionClient oapiTransferInstruction.ClientWithResponsesInterface,
	feeInputHoldings []string,
) error {
	withToken := amountStr != ""

	// Parse amount early to catch errors
	var normalizedAmount string
	if withToken {
		var err error
		normalizedAmount, err = parseDecimalAmount(amountStr)
		if err != nil {
			return fmt.Errorf("invalid --amount %q (supports exponents, e.g. 1e-2): %w", amountStr, err)
		}
	}

	linkInstrumentId := b.Profile.LinkInstrumentID

	// --- Resolve fee token holdings ---
	var feeTokenInputCids []types.CONTRACT_ID
	if len(feeInputHoldings) > 0 {
		for _, cid := range feeInputHoldings {
			feeTokenInputCids = append(feeTokenInputCids, types.CONTRACT_ID(cid))
		}
	} else {
		feeTokenHoldings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, feeTokenInstrumentId)
		if err != nil {
			return fmt.Errorf("list fee token holdings: %w", err)
		}
		for _, h := range feeTokenHoldings {
			feeTokenInputCids = append(feeTokenInputCids, types.CONTRACT_ID(h.ContractID))
		}
	}

	// --- Build transfer factory for fee payment ---
	transferFactory, err := testhelpers.GetTransferFactoryV2(ctx, feeTokenTransferInstructionClient, string(feeTokenInstrumentId.Admin), splice_api_token_transfer_instruction_v1.Transfer{
		Sender:           types.PARTY(b.Participant.PartyID),
		Receiver:         b.Profile.CCIPOwnerPartyID,
		Amount:           "1.0",
		InstrumentId:     *feeTokenInstrumentId,
		InputHoldingCids: feeTokenInputCids,
		Meta:             splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
	})
	if err != nil {
		return fmt.Errorf("get transfer factory: %w", err)
	}
	feeChoiceContext, err := contracts.ChoiceContextFromData(transferFactory.ChoiceContextData)
	if err != nil {
		return fmt.Errorf("fee choice context: %w", err)
	}

	// --- Resolve token transfer holdings ---
	var tokenTransferInputCids []types.CONTRACT_ID
	if withToken {
		if len(tokenInputHoldings) > 0 {
			for _, holding := range tokenInputHoldings {
				tokenTransferInputCids = append(tokenTransferInputCids, types.CONTRACT_ID(holding))
			}
		} else {
			tokenHoldings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, linkInstrumentId)
			if err != nil {
				return fmt.Errorf("list LINK holdings: %w", err)
			}
			for _, h := range tokenHoldings {
				tokenTransferInputCids = append(tokenTransferInputCids, types.CONTRACT_ID(h.ContractID))
			}
		}
	}

	// --- Resolve router + sender ---
	routerCid, err := cantonops.GetOrCreateRouter(ctx, b.Participant, b.CCIPEDS)
	if err != nil {
		return err
	}
	senderCid, err := cantonops.GetOrCreateSender(ctx, b.Participant)
	if err != nil {
		return err
	}

	// --- Build the oapi Message for EDS lookups ---
	if gasLimit < 0 {
		// gas limit is negative/unspecified, default to defaultGasLimit, except for token-only transfers
		gasLimit = defaultGasLimit
		if withToken && len(payload) == 0 {
			gasLimit = 0
		}
	}
	executorType := oapiCommon.Empty
	msg := oapiCommon.Message{
		DestinationChainSelector: strconv.FormatUint(b.Profile.EthSelector, 10),
		Executor: struct {
			Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiCommon.MessageExecutorType `json:"type"`
		}{
			Type: executorType,
		},
		FeeToken: oapiCommon.InstrumentId{
			Admin: oapiCommon.PartyId(feeTokenInstrumentId.Admin),
			Id:    string(feeTokenInstrumentId.Id),
		},
		GasLimit: gasLimit,
		Payload:  hex.EncodeToString(payload),
		Receiver: hex.EncodeToString(receiver.Bytes()),
	}
	if withToken {
		msg.TokenTransfer = &oapiCommon.TokenTransfer{
			Amount: normalizedAmount,
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(linkInstrumentId.Admin),
				Id:    string(linkInstrumentId.Id),
			},
		}
	}

	// --- Token pool disclosures (if token transfer) ---
	var (
		tokenPoolSendDisclosure *eds.TokenPoolSendDisclosure
		requiredCCVs            []string
	)
	if withToken {
		tokenPoolAddress, err := eds.GetTokenPoolForToken(ctx, b.CCIPEDS, contracts.EncodeInstrumentID(*linkInstrumentId))
		if err != nil {
			return fmt.Errorf("get LINK token pool: %w", err)
		}
		tps, err := eds.GetTokenPoolSendDisclosure(ctx, b.TokenPoolEDS, msg, tokenPoolAddress.InstanceAddress())
		if err != nil {
			return fmt.Errorf("token pool send disclosure: %w", err)
		}
		tokenPoolSendDisclosure = tps
		requiredCCVs = tps.RequiredCCVs
	}

	// --- CCIP send disclosure (resolves default CCV + default executor) ---
	ccipSendDisclosure, err := eds.GetCCIPSendDisclosure(ctx, b.CCIPEDS, msg, nil, requiredCCVs)
	if err != nil {
		return fmt.Errorf("CCIP send disclosure: %w", err)
	}
	defaultCCVAddress, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
	if err != nil {
		return fmt.Errorf("parse default CCV address: %w", err)
	}
	ccvSendDisclosure, err := eds.GetCCVSendDisclosure(ctx, b.CCVEDS, msg, defaultCCVAddress.InstanceAddress())
	if err != nil {
		return fmt.Errorf("CCV send disclosure: %w", err)
	}

	// --- Executor handling ---
	var (
		executorInput       *sender.ExecutorInput
		executorDisclosures []*apiv2.DisclosedContract
		executorExtraArg    core.ExecutorExtraArg
	)
	switch executorMode {
	case "default":
		defaultExecutorAddress, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
		if err != nil {
			return fmt.Errorf("parse default executor address: %w", err)
		}
		execDisc, err := eds.GetExecutorSendDisclosure(ctx, b.ExecutorEDS, msg, defaultExecutorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
		if err != nil {
			return fmt.Errorf("executor send disclosure: %w", err)
		}
		executorInput = &sender.ExecutorInput{
			ExecutorCid:          types.CONTRACT_ID(execDisc.ContractId),
			ExecutorExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
		}
		executorDisclosures = execDisc.DisclosedContracts
		executorExtraArg = core.ExecutorExtraArg{
			ExecutorUseDefault: &core.ExecutorUseDefault{ExecutorArgs: ""},
		}
	case "none":
		// No executor — message will not be auto-executed on the destination.
		executorInput = nil
		executorDisclosures = nil
		executorExtraArg = core.ExecutorExtraArg{
			ExecutorNoExecutor: &types.UNIT{},
		}
	}

	// --- Build sendArgs ---
	canton2Any := core.Canton2AnyMessage{
		Receiver: types.TEXT(msg.Receiver),
		Payload:  types.TEXT(msg.Payload),
		FeeToken: *feeTokenInstrumentId,
		ExtraArgs: core.ExtraArgs{
			V3: &core.GenericExtraArgsV3{
				GasLimit:      types.INT64(msg.GasLimit),
				Ccvs:          nil,
				Executor:      executorExtraArg,
				TokenReceiver: "",
				TokenArgs:     "",
			},
		},
	}
	if withToken {
		canton2Any.TokenTransfer = &core.TokenTransfer{
			Token:  *linkInstrumentId,
			Amount: types.NUMERIC(msg.TokenTransfer.Amount),
		}
	}

	sendArgs := sender.Send{
		DestinationChainSelector: types.NUMERIC(msg.DestinationChainSelector),
		Message:                  canton2Any,
		Context:                  ccipSendDisclosure.ChoiceContext,
		RouterCid:                types.CONTRACT_ID(routerCid),
		FeeTokenInput: sender.FeeTokenInput{
			SenderInputCids:         feeTokenInputCids,
			FeeTokenConfigCid:       types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid),
			FeeTokenTransferFactory: types.CONTRACT_ID(transferFactory.FactoryID),
			FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
				Context: feeChoiceContext,
				Meta:    splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
			},
		},
		CcvSendInputs: []sender.CCVSendInput{{
			CcvAddress:      ccvSendDisclosure.Address.Binding(),
			CcvCid:          types.CONTRACT_ID(ccvSendDisclosure.ContractId),
			CcvExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
		}},
		ExecutorInput: executorInput,
	}
	if withToken {
		sendArgs.TokenTransferInput = &sender.TokenTransferInput{
			SenderInputCids:  tokenTransferInputCids,
			TokenPoolCid:     types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
			PoolExtraContext: tokenPoolSendDisclosure.ChoiceContext,
		}
	}

	// --- Concatenate disclosures ---
	allDisclosures := slices.Concat(
		transferFactory.DisclosedContracts,
		ccipSendDisclosure.DisclosedContracts,
		ccvSendDisclosure.DisclosedContracts,
		executorDisclosures,
	)
	if withToken {
		allDisclosures = slices.Concat(allDisclosures, tokenPoolSendDisclosure.DisclosedContracts)
	}

	// Ask for confirmation
	fmt.Println("About to send message:")
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.AppendHeader(table.Row{"Field", "Value"})
	tw.AppendRow(table.Row{"Receiver", "0x" + msg.Receiver})
	tw.AppendRow(table.Row{"Data", "0x" + msg.Payload})
	tw.AppendRow(table.Row{"Gas Limit", msg.GasLimit})
	tw.AppendRow(table.Row{"Fee Token", fmt.Sprintf("%s@%s", msg.FeeToken.Id, msg.FeeToken.Admin)})
	tw.AppendRow(table.Row{"Fee Input Holdings", feeTokenInputCids})
	if msg.TokenTransfer != nil {
		tw.AppendRow(table.Row{"Token Transfer - Token", fmt.Sprintf("%s@%s", msg.TokenTransfer.Token.Id, msg.TokenTransfer.Token.Admin)})
		tw.AppendRow(table.Row{"Token Transfer - Amount", msg.TokenTransfer.Amount})
		tw.AppendRow(table.Row{"Token Transfer - Input Holdings", tokenTransferInputCids})
	}
	fmt.Println(tw.Render())
	fmt.Println("Confirm? (Y/N)")
	if !input.Confirm() {
		return fmt.Errorf("cancel sending message")
	}

	// --- Submit ---
	fmt.Println("⏳ Sending message...")
	resp, err := b.Participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					ContractId:     senderCid,
					Choice:         "Send",
					ChoiceArgument: ledger.MapToValue(sendArgs),
				}},
			}},
			ActAs:              []string{b.Participant.PartyID},
			DisclosedContracts: allDisclosures,
		},
	})
	if err != nil {
		return fmt.Errorf("submit Send: %w", err)
	}
	fmt.Printf("Message sent in Update: %s\n", resp.GetTransaction().GetUpdateId())
	fmt.Println(b.CantonExplorerLink(resp.GetTransaction().GetUpdateId()))

	messageId, err := cantonops.GetMessageIdFromTransaction(resp.GetTransaction())
	if err != nil {
		return err
	}
	fmt.Printf("✅ Message sent with MessageID: 0x%s\n", messageId)
	fmt.Println(b.CCIPExplorerLink(messageId))

	return nil
}
