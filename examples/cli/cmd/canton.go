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
	"github.com/smartcontractkit/chainlink-canton/examples/cli/internal/input"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/spf13/cobra"

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
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
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
	return f.String(), nil
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
	ccvAddress := contracts.BytesToInstanceAddress(vr.VerifierDestAddress)

	ccipExecuteDisclosure, err := eds.GetCCIPExecuteDisclosure(ctx, b.CCIPEDS, encodedHex)
	if err != nil {
		return fmt.Errorf("CCIP execute disclosure: %w", err)
	}
	ccvExecuteDisclosure, err := eds.GetCCVExecuteDisclosure(ctx, b.CCVEDS, encodedHex, ccvAddress)
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
		payload      string
		executor     string
		feeTokenName string
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

			return cantonSend(ctx, b, msgReceiver, []byte(payload), "", executor, feeTokenInstrumentId, feeTokenTransferClient)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", fmt.Sprintf("destination EVM receiver address (0x-prefixed) (defaults to a CCIP Receiver contract)"))
	c.Flags().StringVar(&payload, "payload", "Hello, EVM from Canton!", "message payload (text)")
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
	c.Flags().StringVar(&feeTokenName, "fee-token", "link", "fee token (link|native)")

	return c
}

func newCantonSendTokenCmd(g *Globals) *cobra.Command {
	var (
		receiverHex string
		amountStr   string
		payload     string
		executor    string
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

			msgReceiver := common.HexToAddress(receiverHex)
			if receiverHex == "" {
				msgReceiver = b.ETHAddress
			}

			feeTokenInstrumentId := b.Profile.AmuletInstrumentID
			feeTokenTransferInstructionClient := b.AmuletTransferClient

			return cantonSend(ctx, b, msgReceiver, []byte(payload), amountStr, executor, feeTokenInstrumentId, feeTokenTransferInstructionClient)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", "destination EVM receiver address (0x-prefixed) (defaults to own address)")
	c.Flags().StringVar(&amountStr, "amount", "", "LINK token transfer amount as decimal (e.g. 0.12345, 1e-2); (required)")
	c.Flags().StringVar(&payload, "payload", "", "optional message payload (text) to attach to the token transfer")
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
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
	amountStr string,
	executorMode string,
	feeTokenInstrumentId *splice_api_token_holding_v1.InstrumentId,
	feeTokenTransferInstructionClient oapiTransferInstruction.ClientWithResponsesInterface,
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
	feeTokenHoldings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, feeTokenInstrumentId)
	if err != nil {
		return fmt.Errorf("list fee token holdings: %w", err)
	}
	feeTokenInputCids := make([]types.CONTRACT_ID, len(feeTokenHoldings))
	for i, h := range feeTokenHoldings {
		feeTokenInputCids[i] = types.CONTRACT_ID(h.ContractID)
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
		tokenHoldings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, linkInstrumentId)
		if err != nil {
			return fmt.Errorf("list LINK holdings: %w", err)
		}
		tokenTransferInputCids = make([]types.CONTRACT_ID, len(tokenHoldings))
		for i, h := range tokenHoldings {
			tokenTransferInputCids[i] = types.CONTRACT_ID(h.ContractID)
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
	gasLimit := 50_000
	if withToken && len(payload) == 0 {
		gasLimit = 0
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
	tw.AppendRow(table.Row{"Fee Token", fmt.Sprintf("%s@%s", msg.FeeToken.Id, msg.FeeToken.Admin)})
	if msg.TokenTransfer != nil {
		tw.AppendRow(table.Row{"Token Transfer - Token", fmt.Sprintf("%s@%s", msg.TokenTransfer.Token.Id, msg.TokenTransfer.Token.Admin)})
		tw.AppendRow(table.Row{"Token Transfer - Amount", msg.TokenTransfer.Amount})
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
