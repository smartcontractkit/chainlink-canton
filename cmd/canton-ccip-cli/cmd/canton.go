package cmd

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	indexerCommon "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/common"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/cantonops"
	"github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/clients"
	"github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/finality"
	"github.com/smartcontractkit/chainlink-canton/cmd/canton-ccip-cli/internal/input"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/clientapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/events"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/chainlink/chainlinkapi"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_transfer_instruction_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/eds/api/common"
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
	c.AddCommand(newCantonSyncReceiverCCVCmd(g))
	c.AddCommand(newCantonSetReceiverRequiredCCVsCmd(g))
	c.AddCommand(newCantonArchiveReceiverCmd(g))
	c.AddCommand(newCantonListEventsCmd(g))
	c.AddCommand(newCantonListReceiversCmd(g))
	c.AddCommand(newCantonListHoldingsCmd(g))
	c.AddCommand(newCantonListTransferInstructionsCmd(g))
	c.AddCommand(newCantonCreateTransferCmd(g))
	c.AddCommand(newCantonAcceptTransferCmd(g))

	c.AddCommand(newCantonExternalPartyCommand(g))

	return c
}

func resolveCantonToken(b *clients.Bundle, token string) (*splice_api_token_holding_v1.InstrumentId, oapiTransferInstruction.ClientWithResponsesInterface, error) {
	switch token {
	case "link":
		linkEdsClients, err := b.GetTokenStandardClients(b.Profile.LinkInstrumentID.Admin)
		if err != nil {
			return nil, nil, fmt.Errorf("get link EDS clients: %w", err)
		}

		return b.Profile.LinkInstrumentID, linkEdsClients.TransferInstructionClient, nil
	case "native":
		nativeEdsClients, err := b.GetTokenStandardClients(b.Profile.AmuletInstrumentID.Admin)
		if err != nil {
			return nil, nil, fmt.Errorf("get native EDS clients: %w", err)
		}

		return b.Profile.AmuletInstrumentID, nativeEdsClients.TransferInstructionClient, nil
	default:
		// Try to parse InstrumentId, should be in format <ID>@<ADMIN>
		if split := strings.Split(token, "@"); len(split) == 2 {
			instrumentId := &splice_api_token_holding_v1.InstrumentId{
				Id:    types.TEXT(split[0]),
				Admin: types.PARTY(split[1]),
			}
			tokenEdsClients, err := b.GetTokenStandardClients(instrumentId.Admin)
			if err != nil {
				return nil, nil, fmt.Errorf("no Token Standard URL found for token %q: %w", token, err)
			}

			return instrumentId, tokenEdsClients.TransferInstructionClient, nil
		}

		return nil, nil, fmt.Errorf("invalid token %q, allowed values: (link|native|<InstrumentId>)", token)
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
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			var tmpl *apiv2.Identifier
			switch eventName {
			case "sent":
				tmpl = contracts.IdentifierFromBinding(events.CCIPMessageSent{})
			case "executed":
				tmpl = contracts.IdentifierFromBinding(events.ExecutionStateChanged{})
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
					ev, err := bindings.UnmarshalCreatedEvent[events.CCIPMessageSent](c.GetCreatedEvent())
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
					ev, err := bindings.UnmarshalCreatedEvent[events.ExecutionStateChanged](c.GetCreatedEvent())
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
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			holdings, err := testhelpers.ListActiveContractsByInterfaceId(ctx, b.Participant, contracts.MustTemplateIDFromString(splice_api_token_holding_v1.IHoldingInterfaceID()).ToLedgerIdentifier())
			if err != nil {
				return fmt.Errorf("list holdings: %w", err)
			}
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.Style().Title.Align = text.AlignCenter
			tw.SetAutoIndex(true)
			if withContractId {
				tw.AppendHeader(table.Row{"Instrument ID", "Owner", "Amount", "Locked", "Symbol", "Name", "Contract ID"})
			} else {
				tw.AppendHeader(table.Row{"Instrument ID", "Owner", "Amount", "Locked", "Symbol", "Name"})
			}
			for _, h := range holdings {
				for _, view := range h.GetCreatedEvent().GetInterfaceViews() {
					var hv splice_api_token_holding_v1.HoldingView
					if err := ledger.RecordToStruct(view.GetViewValue(), &hv); err != nil {
						return fmt.Errorf("decode holding view: %w", err)
					}

					// Query TokenMetadata API
					var (
						tokenName, tokenSymbol string
					)
					if edsClients, err := b.GetTokenStandardClients(hv.InstrumentId.Admin); err == nil {
						instrumentInfo, err := edsClients.MetadataClient.GetInstrumentWithResponse(ctx, string(hv.InstrumentId.Id))
						if err == nil {
							if instrumentInfo.StatusCode() == http.StatusOK && instrumentInfo.JSON200 != nil {
								tokenName, tokenSymbol = instrumentInfo.JSON200.Name, instrumentInfo.JSON200.Symbol
							}
						}
					}

					if withContractId {
						tw.AppendRow(table.Row{fmt.Sprintf("%s@%s", hv.InstrumentId.Id, hv.InstrumentId.Admin), hv.InstrumentId.Admin, hv.Owner, hv.Amount, hv.Lock != nil, tokenSymbol, tokenName, h.GetCreatedEvent().GetContractId()})
					} else {
						tw.AppendRow(table.Row{fmt.Sprintf("%s@%s", hv.InstrumentId.Id, hv.InstrumentId.Admin), hv.Owner, hv.Amount, hv.Lock != nil, tokenSymbol, tokenName})
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

// ---------------- canton list-transfer-instructions ----------------

func newCantonListTransferInstructionsCmd(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "list-transfer-instructions",
		Short: "List all transfer instructions for the configured party",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			holdings, err := testhelpers.ListActiveContractsByInterfaceId(ctx, b.Participant, contracts.MustTemplateIDFromString(splice_api_token_transfer_instruction_v1.ITransferInstructionInterfaceID()).ToLedgerIdentifier())
			if err != nil {
				return fmt.Errorf("list holdings: %w", err)
			}
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.Style().Title.Align = text.AlignCenter
			tw.SetAutoIndex(true)
			tw.AppendHeader(table.Row{"Instrument ID", "Admin", "Sender", "Receiver", "Amount", "Requested At", "Execute Before", "ContractId"})
			for _, h := range holdings {
				for _, view := range h.GetCreatedEvent().GetInterfaceViews() {
					var tiv splice_api_token_transfer_instruction_v1.TransferInstructionView
					if err := ledger.RecordToStruct(view.GetViewValue(), &tiv); err != nil {
						return fmt.Errorf("decode transfer instruction view: %w", err)
					}
					tw.AppendRow(table.Row{fmt.Sprintf("%s@%s", tiv.Transfer.InstrumentId.Id, tiv.Transfer.InstrumentId.Admin), tiv.Transfer.Sender, tiv.Transfer.Receiver, tiv.Transfer.Amount, time.Time(tiv.Transfer.RequestedAt).String(), time.Time(tiv.Transfer.ExecuteBefore).String(), h.GetCreatedEvent().GetContractId()})
				}
			}
			fmt.Println(tw.Render())

			return nil
		},
	}

	return c
}

// ---------------- canton create transfer ----------------
func newCantonCreateTransferCmd(g *Globals) *cobra.Command {
	var (
		token                       string
		receiverParty               string
		inputHoldingCids            []string
		amount                      string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "create-transfer",
		Short: "Create an outgoing TransferInstruction",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			instrumentId, transferInstructionClient, err := resolveCantonToken(b, token)
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
				holdings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, instrumentId, testhelpers.WithUnlockedHoldingsOnly())
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
				return fmt.Errorf("get transfer factory for token: %w", err)
			}
			choiceContext, err := contracts.ChoiceContextFromData(transferFactory.ChoiceContextData)
			if err != nil {
				return fmt.Errorf("unmarshal choice context: %w", err)
			}

			// Ask for confirmation
			inputHoldingStrings := make([]string, len(inputHoldings))
			for i, cid := range inputHoldings {
				inputHoldingStrings[i] = string(cid)
			}
			fmt.Println("About to send transfer:")
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.AppendHeader(table.Row{"Field", "Value"})
			tw.AppendRow(table.Row{"Sender", transfer.Sender})
			tw.AppendRow(table.Row{"Receiver", transfer.Receiver})
			tw.AppendRow(table.Row{"Amount", transfer.Amount})
			tw.AppendRow(table.Row{"Instrument ID", fmt.Sprintf("%s@%s", transfer.InstrumentId.Id, transfer.InstrumentId.Admin)})
			tw.AppendRow(table.Row{"Requested at", time.Time(transfer.RequestedAt).String()})
			tw.AppendRow(table.Row{"Execute before", time.Time(transfer.ExecuteBefore).String()})
			tw.AppendRow(table.Row{"Input Holdings", strings.Join(inputHoldingStrings, ",")})
			fmt.Println(tw.Render())
			fmt.Println("Confirm? (Y/N)")
			if !input.Confirm() {
				return fmt.Errorf("cancel creating transfer")
			}

			// Create Transfer
			tx, err := cantonops.CantonSubmit(
				ctx,
				b.Participant,
				useLedger,
				[]*apiv2.Command{{
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
				transferFactory.DisclosedContracts,
				&apiv2.TransactionFormat{
					EventFormat: &apiv2.EventFormat{
						FiltersByParty: map[string]*apiv2.Filters{
							b.Participant.PartyID: &apiv2.Filters{Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &apiv2.InterfaceFilter{
										InterfaceId:          contracts.MustTemplateIDFromString(splice_api_token_transfer_instruction_v1.ITransferInstructionInterfaceID()).ToLedgerIdentifier(),
										IncludeInterfaceView: true,
									}},
								},
							}},
						},
						Verbose: true,
					},
					TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
				},
				packageSelectionPreferences,
			)
			if err != nil {
				return fmt.Errorf("submit transfer: %w", err)
			}
			fmt.Printf("✅ Submitted transfer in update: %s\n", tx.GetUpdateId())
			fmt.Println(b.CantonExplorerLink(tx.GetUpdateId()))

			_, ce, err := testhelpers.GetCreatedInterfaceViewFromTransaction[splice_api_token_transfer_instruction_v1.TransferInstructionView](tx, contracts.MustTemplateIDFromString(splice_api_token_transfer_instruction_v1.ITransferInstructionInterfaceIDWithPackageID(splice_api_token_transfer_instruction_v1.PackageID)).ToLedgerIdentifier())
			if err == nil {
				fmt.Printf("📑 Pending TransferInstruction created: %s\n", ce.GetContractId())
			}

			return nil
		},
	}
	c.Flags().StringVar(&token, "token", "link", "token to transfer (link|native|<InstrumentId>)")
	c.Flags().StringVar(&receiverParty, "receiver", "", "party to receive the transfer (defaults to own party)")
	c.Flags().StringArrayVar(&inputHoldingCids, "input", nil, "the holding(s) to be used as an input for the transfer. If unspecified, all current holdings will be used.")
	c.Flags().StringVar(&amount, "amount", "", "the amount to transfer (required)")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")
	_ = c.MarkFlagRequired("amount")

	return c
}

func newCantonAcceptTransferCmd(g *Globals) *cobra.Command {
	var (
		contractID                  string
		token                       string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "accept-transfer",
		Short: "Accept an incoming TransferInstruction by contract ID",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			_, transferInstructionClient, err := resolveCantonToken(b, token)
			if err != nil {
				return err
			}

			acceptContextResp, err := transferInstructionClient.GetTransferInstructionAcceptContextWithResponse(ctx, contractID, oapiTransferInstruction.GetChoiceContextRequest{})
			if err != nil {
				return fmt.Errorf("get transfer instruction accept context: %w", err)
			}
			if acceptContextResp.StatusCode() != http.StatusOK || acceptContextResp.JSON200 == nil {
				return fmt.Errorf("unexpected transfer instruction accept context status=%d", acceptContextResp.StatusCode())
			}

			disclosedContracts := make([]*apiv2.DisclosedContract, 0, len(acceptContextResp.JSON200.DisclosedContracts))
			for _, contract := range acceptContextResp.JSON200.DisclosedContracts {
				id, err := testhelpers.TemplateIdFromString(contract.TemplateId)
				if err != nil {
					return fmt.Errorf("parse accept-context template id: %w", err)
				}
				createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
				if err != nil {
					return fmt.Errorf("decode accept-context created event blob: %w", err)
				}
				disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
					TemplateId:       id,
					ContractId:       contract.ContractId,
					CreatedEventBlob: createdEventBlob,
					SynchronizerId:   contract.SynchronizerId,
				})
			}

			acceptContext, err := contracts.ChoiceContextFromData(acceptContextResp.JSON200.ChoiceContextData)
			if err != nil {
				return fmt.Errorf("convert transfer instruction accept context: %w", err)
			}

			// Accept Transfer
			tx, err := cantonops.CantonSubmit(
				ctx,
				b.Participant,
				useLedger,
				[]*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
						ContractId: contractID,
						Choice:     "TransferInstruction_Accept",
						ChoiceArgument: ledger.MapToValue(splice_api_token_transfer_instruction_v1.TransferInstructionAccept{
							ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
								Context: acceptContext,
							},
						}),
					}},
				}},
				disclosedContracts,
				nil,
				packageSelectionPreferences,
			)
			if err != nil {
				return fmt.Errorf("submit accept transfer: %w", err)
			}
			fmt.Println("✅ Transfer accepted in update:", tx.GetUpdateId())
			fmt.Println(b.CantonExplorerLink(tx.GetUpdateId()))

			return nil
		},
	}
	c.Flags().StringVar(&contractID, "contract-id", "", "TransferInstruction contract ID to accept (required)")
	c.Flags().StringVar(&token, "token", "link", "token of the transfer instruction (link|native|<InstrumentId>)")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")
	_ = c.MarkFlagRequired("contract-id")

	return c
}

// ---------------- canton execute ----------------

func newCantonExecuteCmd(g *Globals) *cobra.Command {
	var (
		messageIDHex                string
		wait                        time.Duration
		finalityName                string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "execute",
		Short: "Execute on Canton a message sent from EVM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			fin, err := finality.Parse(finalityName)
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

			return cantonExecute(ctx, b, resp.Results, fin, useLedger, packageSelectionPreferences)
		},
	}
	c.Flags().StringVar(&messageIDHex, "message-id", "", "CCIP message id (0x-prefixed hex) (required)")
	c.Flags().DurationVar(&wait, "wait", 15*time.Minute, "max time to wait for verifier results")
	c.Flags().StringVar(&finalityName, "finality", "finality", "must match the send: finality (full), safe, or block depth 1-65535")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")
	_ = c.MarkFlagRequired("message-id")

	return c
}

func newCantonListReceiversCmd(g *Globals) *cobra.Command {
	var withContractId bool
	c := &cobra.Command{
		Use:   "list-receivers",
		Short: "List all CCIPReceiver contracts visible to the configured party",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			active, err := testhelpers.ListActiveContractsByTemplateId(ctx, b.Participant, contracts.IdentifierFromBinding(receiver.CCIPReceiver{}))
			if err != nil {
				return fmt.Errorf("list receivers: %w", err)
			}
			tw := table.NewWriter()
			tw.SetStyle(table.StyleLight)
			tw.Style().Title.Align = text.AlignCenter
			tw.SetAutoIndex(true)
			tw.SetTitle("CCIPReceivers")
			if withContractId {
				tw.AppendHeader(table.Row{"Instance ID", "Owner", "Finality", "Required CCVs", "Optional CCVs", "Optional Threshold", "Contract ID"})
			} else {
				tw.AppendHeader(table.Row{"Instance ID", "Owner", "Finality", "Required CCVs", "Optional CCVs", "Optional Threshold"})
			}
			for _, ac := range active {
				recv, err := bindings.UnmarshalCreatedEvent[receiver.CCIPReceiver](ac.GetCreatedEvent())
				if err != nil {
					return fmt.Errorf("unmarshal CCIPReceiver: %w", err)
				}
				ccvStrings := func(ccvs []chainlinkapi.RawInstanceAddress) string {
					out := make([]string, len(ccvs))
					for i, ccv := range ccvs {
						out[i] = string(ccv.Unpack)
					}

					return strings.Join(out, ",")
				}
				if withContractId {
					tw.AppendRow(table.Row{
						recv.InstanceId,
						recv.Owner,
						cantonops.ReceiverFinalityLabel(recv.ReceiverFinalityConfig),
						ccvStrings(recv.RequiredCCVs),
						ccvStrings(recv.OptionalCCVs),
						recv.OptionalThreshold,
						ac.GetCreatedEvent().GetContractId(),
					})
				} else {
					tw.AppendRow(table.Row{
						recv.InstanceId,
						recv.Owner,
						cantonops.ReceiverFinalityLabel(recv.ReceiverFinalityConfig),
						ccvStrings(recv.RequiredCCVs),
						ccvStrings(recv.OptionalCCVs),
						recv.OptionalThreshold,
					})
				}
			}
			tw.AppendFooter(table.Row{"Total", len(active)})
			fmt.Println(tw.Render())

			return nil
		},
	}
	c.Flags().BoolVar(&withContractId, "cid", false, "include contract IDs in the output")

	return c
}

func newCantonSetReceiverRequiredCCVsCmd(g *Globals) *cobra.Command {
	var useLedger string
	c := &cobra.Command{
		Use:   "set-receiver-required-ccvs <contract-id> <ccv> [ccv...]",
		Short: "Update a CCIPReceiver's required CCVs by contract ID",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			receiverCid := args[0]
			requiredCCVs := make([]contracts.RawInstanceAddress, 0, len(args)-1)
			for _, ccv := range args[1:] {
				// intentionally not doing validation here to be able to set invalid CCVs for testing purposes
				requiredCCVs = append(requiredCCVs, contracts.RawInstanceAddress(ccv))
			}

			newCid, err := cantonops.UpdateReceiverRequiredCCVs(ctx, b.Participant, receiverCid, requiredCCVs, useLedger)
			if err != nil {
				return err
			}
			fmt.Printf("New CCIPReceiver CID: %s\n", newCid)

			return nil
		},
	}
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")

	return c
}

func newCantonArchiveReceiverCmd(g *Globals) *cobra.Command {
	var useLedger string
	c := &cobra.Command{
		Use:   "archive-receiver <contract-id>",
		Short: "Archive a CCIPReceiver contract by contract ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			receiverCid := args[0]

			tx, err := cantonops.CantonSubmit(
				ctx,
				b.Participant,
				useLedger,
				[]*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId:     contracts.IdentifierFromBinding(receiver.CCIPReceiver{}),
						ContractId:     receiverCid,
						Choice:         "Archive",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{}}},
					}},
				}},
				nil,
				nil,
			)
			if err != nil {
				return fmt.Errorf("submit Archive: %w", err)
			}
			fmt.Println("✅ CCIPReceiver archived in update:", tx.GetUpdateId())
			fmt.Println(b.CantonExplorerLink(tx.GetUpdateId()))

			return nil
		},
	}
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")

	return c
}

func newCantonSyncReceiverCCVCmd(g *Globals) *cobra.Command {
	var (
		requiredCCV                 string
		finalityName                string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "sync-receiver-ccv",
		Short: "Deploy or update CCIPReceiver required CCVs (run once before inbound load/e2e on prod)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}
			fin, err := finality.Parse(finalityName)
			if err != nil {
				return err
			}
			raw, err := contracts.RawInstanceAddressFromString(requiredCCV)
			if err != nil {
				return fmt.Errorf("parse --required-ccv: %w", err)
			}

			_, err = cantonops.GetOrCreateReceiver(ctx, b.Participant, fin.Receiver, []contracts.RawInstanceAddress{raw}, useLedger, packageSelectionPreferences)

			return err
		},
	}
	c.Flags().StringVar(&requiredCCV, "required-ccv", "", "Canton committee verifier raw address (instanceId@owner)")
	c.Flags().StringVar(&finalityName, "finality", "1", "receiver finality profile: finality|safe|1-65535")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")
	_ = c.MarkFlagRequired("required-ccv")

	return c
}

// execute a message on Canton
// verifierResults must be non-empty
func cantonExecute(ctx context.Context, b *clients.Bundle, verifierResults []indexerCommon.VerifierResultWithMetadata, fin finality.Parsed, useLedger string, packageSelectionPreferences []string) error {
	message := verifierResults[0].VerifierResult.Message
	withToken := message.TokenTransfer != nil
	receiverParty := types.PARTY(b.Participant.PartyID)

	encodedMessage, err := message.Encode()
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}
	encodedHex := hex.EncodeToString(encodedMessage)

	// Get CCIP disclosures
	ccipEdsClients, err := b.GetEDSClients(b.Profile.CCIPOwnerPartyID)
	if err != nil {
		return fmt.Errorf("get CCIP EDS clients: %w", err)
	}
	ccipExecuteDisclosure, err := eds.GetCCIPExecuteDisclosure(ctx, ccipEdsClients.CCIPEDS, encodedHex, receiverParty)
	if err != nil {
		return fmt.Errorf("CCIP execute disclosure: %w", err)
	}
	allDisclosures := ccipExecuteDisclosure.DisclosedContracts

	// Get CCV disclosures
	verifierAddresses := make([]contracts.RawInstanceAddress, len(verifierResults))
	ccvInputs := make([]receiver.CCVInput, len(verifierResults))
	for i, result := range verifierResults {
		verifierRawAddress, err := contracts.RawInstanceAddressFromString(string(result.VerifierResult.VerifierDestAddress))
		if err != nil {
			return fmt.Errorf("parse VerifierDestAddress: %w", err)
		}
		verifierAddresses[i] = verifierRawAddress
		ccvEdsClients, err := b.GetEDSClients(types.PARTY(verifierRawAddress.Owner()))
		if err != nil {
			return fmt.Errorf("get CCV EDS clients: %w", err)
		}
		ccvExecuteDisclosure, err := eds.GetCCVExecuteDisclosure(ctx, ccvEdsClients.CCVEDS, encodedHex, verifierRawAddress.InstanceAddress(), receiverParty)
		if err != nil {
			return fmt.Errorf("CCV execute disclosure: %w", err)
		}

		fmt.Printf("Got CCV disclosure for %s: CID: %v, %d disclosed contract(s)\n", verifierRawAddress, ccvExecuteDisclosure.ContractId, len(ccvExecuteDisclosure.DisclosedContracts))
		allDisclosures = append(allDisclosures, ccvExecuteDisclosure.DisclosedContracts...)
		ccvInputs[i] = receiver.CCVInput{
			CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
			VerifierResults: types.TEXT(hex.EncodeToString(result.VerifierResult.CCVData)),
			Context:         ccvExecuteDisclosure.ChoiceContext,
		}
	}

	routerCid, err := cantonops.GetOrCreateRouter(ctx, b.Participant, ccipEdsClients.CCIPEDS, useLedger, packageSelectionPreferences)
	if err != nil {
		return err
	}
	receiverCid, err := cantonops.GetOrCreateReceiver(ctx, b.Participant, fin.Receiver, verifierAddresses, useLedger, packageSelectionPreferences)
	if err != nil {
		return err
	}
	fmt.Printf("PerPartyRouter CID: %s\nCCIPReceiver CID: %s\n", routerCid, receiverCid)

	var tokenTransferInput *receiver.TokenTransferInput

	if withToken {
		targetInstrumentId := contracts.BytesToEncodedInstrumentID(message.TokenTransfer.DestTokenAddress)
		tokenPoolAddress, err := eds.GetTokenPoolForToken(ctx, ccipEdsClients.CCIPEDS, targetInstrumentId)
		if err != nil {
			return fmt.Errorf("get token pool: %w", err)
		}
		tokenPoolEdsClients, err := b.GetEDSClients(types.PARTY(tokenPoolAddress.Owner()))
		if err != nil {
			return fmt.Errorf("get token pool EDS clients: %w", err)
		}
		tokenPoolExecuteDisclosure, err := eds.GetTokenPoolExecuteDisclosure(ctx, tokenPoolEdsClients.TokenPoolEDS, encodedHex, tokenPoolAddress.InstanceAddress(), receiverParty)
		if err != nil {
			return fmt.Errorf("token pool execute disclosure: %w", err)
		}
		tokenTransferInput = &receiver.TokenTransferInput{
			TokenPoolCid:       types.CONTRACT_ID(tokenPoolExecuteDisclosure.ContractId),
			TokenReceiverParty: types.PARTY(b.Participant.PartyID),
			Context:            tokenPoolExecuteDisclosure.ChoiceContext,
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
		CcvInputs:      ccvInputs,
	}

	fmt.Println("⏳ Executing message...")
	tx, err := cantonops.CantonSubmit(
		ctx,
		b.Participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
				TemplateId:     contracts.IdentifierFromBinding(receiver.CCIPReceiver{}),
				ContractId:     receiverCid,
				Choice:         "Execute",
				ChoiceArgument: ledger.MapToValue(executeArgs),
			}},
		}},
		allDisclosures,
		nil,
		packageSelectionPreferences,
	)
	if err != nil {
		return fmt.Errorf("submit Execute: %w", err)
	}
	fmt.Printf("✅ Message executed in Update: %s\n", tx.GetUpdateId())
	fmt.Println(b.CantonExplorerLink(tx.GetUpdateId()))

	return nil
}

// ---------------- canton send-message / send-token ----------------

func newCantonSendMessageCmd(g *Globals) *cobra.Command {
	var (
		receiverHex                 string
		gasLimit                    int
		payload                     string
		executor                    string
		feeToken                    string
		feeInput                    []string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "send-message",
		Short: "Send a message-only CCIP message from Canton to EVM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			if executor != "default" && executor != "none" {
				return fmt.Errorf("invalid --executor %q (default|none)", executor)
			}

			feeTokenInstrumentId, feeTokenTransferClient, err := resolveCantonToken(b, feeToken)
			if err != nil {
				return err
			}

			msgReceiver := common.HexToAddress(receiverHex)
			if receiverHex == "" {
				msgReceiver = b.Profile.CCIPReceiverContract
			}

			return cantonSend(ctx, b, msgReceiver, []byte(payload), gasLimit, nil, "", nil, executor, feeTokenInstrumentId, feeTokenTransferClient, feeInput, useLedger, packageSelectionPreferences)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", "destination EVM receiver address (0x-prefixed) (defaults to a CCIP Receiver contract)")
	c.Flags().StringVar(&payload, "payload", "Hello, EVM from Canton!", "message payload (text)")
	c.Flags().IntVar(&gasLimit, "gas-limit", -1, fmt.Sprintf("gas limit for EVM execution, defaults to %v for message transfers", defaultGasLimit))
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
	c.Flags().StringVar(&feeToken, "fee-token", "link", "fee token (link|native|<InstrumentId>)")
	c.Flags().StringArrayVar(&feeInput, "fee-input", nil, "the holding(s) to be used as an input for the fee payment. If unspecified, all current holdings will be used.")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")

	return c
}

func newCantonSendTokenCmd(g *Globals) *cobra.Command {
	var (
		receiverHex                 string
		gasLimit                    int
		token                       string
		amountStr                   string
		payload                     string
		executor                    string
		feeToken                    string
		feeInput                    []string
		tokenInput                  []string
		useLedger                   string
		packageSelectionPreferences []string
	)
	c := &cobra.Command{
		Use:   "send-token",
		Short: "Send a LINK token transfer CCIP message from Canton to EVM",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			tokenTransferInstrumentId := b.Profile.LinkInstrumentID
			if token != "" {
				split := strings.Split(token, "@")
				if len(split) != 2 || split[0] == "" || split[1] == "" {
					return fmt.Errorf("invalid --token %q (must be in the format of <id>@<admin>)", token)
				}
				tokenTransferInstrumentId = &splice_api_token_holding_v1.InstrumentId{
					Id:    types.TEXT(split[0]),
					Admin: types.PARTY(split[1]),
				}
			}

			if executor != "default" && executor != "none" {
				return fmt.Errorf("invalid --executor %q (default|none)", executor)
			}

			feeTokenInstrumentId, feeTokenTransferClient, err := resolveCantonToken(b, feeToken)
			if err != nil {
				return err
			}

			msgReceiver := common.HexToAddress(receiverHex)
			if receiverHex == "" {
				msgReceiver = b.ETHAddress
			}

			return cantonSend(ctx, b, msgReceiver, []byte(payload), gasLimit, tokenTransferInstrumentId, amountStr, tokenInput, executor, feeTokenInstrumentId, feeTokenTransferClient, feeInput, useLedger, packageSelectionPreferences)
		},
	}
	c.Flags().StringVar(&receiverHex, "receiver", "", "destination EVM receiver address (0x-prefixed) (defaults to own address)")
	c.Flags().StringVar(&token, "token", "", "instrumentId of the token to transfer, in the format of <id>@<admin> (defaults to LINK)")
	c.Flags().StringVar(&amountStr, "amount", "", "LINK token transfer amount as decimal (e.g. 0.12345, 1e-2); (required)")
	c.Flags().StringVar(&payload, "payload", "", "optional message payload (text) to attach to the token transfer")
	c.Flags().IntVar(&gasLimit, "gas-limit", -1, fmt.Sprintf("gas limit for EVM execution, defaults to %v for message transfers", defaultGasLimit))
	c.Flags().StringVar(&executor, "executor", "default", "executor mode (default|none)")
	c.Flags().StringVar(&feeToken, "fee-token", "native", "fee token (link|native)")
	c.Flags().StringArrayVar(&feeInput, "fee-input", nil, "the holding(s) to be used as an input for the fee payment. If unspecified, all current holdings will be used.")
	c.Flags().StringArrayVar(&tokenInput, "token-input", nil, "the holding(s) to be used as an input for the token transfer. If unspecified, all current holdings will be used.")
	c.Flags().StringVar(&useLedger, "ledger", "", "enable interactive Ledger signing if set. Accepts a derivation path value, either a full path like m/44'/6767'/0'/0'/0' or a depth like 42 in which case it will increment the last component, e.g. m/44'/6767'/0'/0'/42'")
	c.Flags().StringSliceVar(&packageSelectionPreferences, "package-selection-preferences", nil, "comma-separated list of package IDs for package selection preference")
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
	// token transfers
	tokenTransferInstrumentId *splice_api_token_holding_v1.InstrumentId,
	amountStr string,
	tokenInputHoldings []string,
	// executor
	executorMode string,
	// fee token
	feeTokenInstrumentId *splice_api_token_holding_v1.InstrumentId,
	feeTokenTransferInstructionClient oapiTransferInstruction.ClientWithResponsesInterface,
	feeInputHoldings []string,
	useLedger string,
	packageSelectionPreferences []string,
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
		return fmt.Errorf("get transfer factory for fee token: %w", err)
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
			tokenHoldings, err := testhelpers.ListHoldingsForInstrument(ctx, b.Participant, tokenTransferInstrumentId, testhelpers.WithUnlockedHoldingsOnly())
			if err != nil {
				return fmt.Errorf("list LINK holdings: %w", err)
			}
			for _, h := range tokenHoldings {
				tokenTransferInputCids = append(tokenTransferInputCids, types.CONTRACT_ID(h.ContractID))
			}
		}
	}

	// Get the ccipOwner party's EDS clients
	ccipEdsClients, err := b.GetEDSClients(b.Profile.CCIPOwnerPartyID)
	if err != nil {
		return fmt.Errorf("get CCIP EDS clients: %w", err)
	}

	// --- Resolve router + sender ---
	routerCid, err := cantonops.GetOrCreateRouter(ctx, b.Participant, ccipEdsClients.CCIPEDS, useLedger, packageSelectionPreferences)
	if err != nil {
		return err
	}
	senderCid, err := cantonops.GetOrCreateSender(ctx, b.Participant, useLedger, packageSelectionPreferences)
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
		Sender:   b.Participant.PartyID,
		Receiver: hex.EncodeToString(receiver.Bytes()),
	}
	tokenTransferHoldings := make([]string, len(tokenTransferInputCids))
	for i, cid := range tokenTransferInputCids {
		tokenTransferHoldings[i] = string(cid)
	}
	if withToken {
		msg.TokenTransfer = &oapiCommon.TokenTransfer{
			Amount: normalizedAmount,
			Token: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(tokenTransferInstrumentId.Admin),
				Id:    string(tokenTransferInstrumentId.Id),
			},
			HoldingContractIds: new(tokenTransferHoldings),
		}
	}

	// --- Token pool disclosures (if token transfer) ---
	var (
		tokenPoolSendDisclosure *eds.TokenPoolSendDisclosure
		requiredCCVs            []string
	)
	if withToken {
		tokenPoolAddress, err := eds.GetTokenPoolForToken(ctx, ccipEdsClients.CCIPEDS, contracts.EncodeInstrumentID(*tokenTransferInstrumentId))
		if err != nil {
			return fmt.Errorf("get LINK token pool: %w", err)
		}
		tokenPoolEdsClients, err := b.GetEDSClients(types.PARTY(tokenPoolAddress.Owner()))
		if err != nil {
			return fmt.Errorf("get token pool EDS clients: %w", err)
		}
		tps, err := eds.GetTokenPoolSendDisclosure(ctx, tokenPoolEdsClients.TokenPoolEDS, msg, tokenPoolAddress.InstanceAddress())
		if err != nil {
			return fmt.Errorf("token pool send disclosure: %w", err)
		}
		tokenPoolSendDisclosure = tps
		requiredCCVs = tps.RequiredCCVs
	}

	// --- CCIP send disclosure (resolves default CCV + default executor) ---
	ccipSendDisclosure, err := eds.GetCCIPSendDisclosure(ctx, ccipEdsClients.CCIPEDS, msg, nil, requiredCCVs)
	if err != nil {
		return fmt.Errorf("CCIP send disclosure: %w", err)
	}
	defaultCCVAddress, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
	if err != nil {
		return fmt.Errorf("parse default CCV address: %w", err)
	}
	ccvEdsClients, err := b.GetEDSClients(types.PARTY(defaultCCVAddress.Owner()))
	if err != nil {
		return fmt.Errorf("get CCV EDS clients: %w", err)
	}
	ccvSendDisclosure, err := eds.GetCCVSendDisclosure(ctx, ccvEdsClients.CCVEDS, msg, defaultCCVAddress.InstanceAddress())
	if err != nil {
		return fmt.Errorf("CCV send disclosure: %w", err)
	}

	// --- Executor handling ---
	var (
		executorInput       *sender.ExecutorInput
		executorDisclosures []*apiv2.DisclosedContract
		executorExtraArg    clientapi.ExecutorExtraArg
	)
	switch executorMode {
	case "default":
		defaultExecutorAddress, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
		if err != nil {
			return fmt.Errorf("parse default executor address: %w", err)
		}
		executorEdsClients, err := b.GetEDSClients(types.PARTY(defaultExecutorAddress.Owner()))
		if err != nil {
			return fmt.Errorf("get executor EDS clients: %w", err)
		}
		execDisc, err := eds.GetExecutorSendDisclosure(ctx, executorEdsClients.ExecutorEDS, msg, defaultExecutorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
		if err != nil {
			return fmt.Errorf("executor send disclosure: %w", err)
		}
		executorInput = &sender.ExecutorInput{
			ExecutorCid: types.CONTRACT_ID(execDisc.ContractId),
			Context:     splice_api_token_metadata_v1.ChoiceContext{},
		}
		executorDisclosures = execDisc.DisclosedContracts
		executorExtraArg = clientapi.ExecutorExtraArg{
			ExecutorUseDefault: &clientapi.ExecutorUseDefault{ExecutorArgs: ""},
		}
	case "none":
		// No executor — message will not be auto-executed on the destination.
		executorInput = nil
		executorDisclosures = nil
		executorExtraArg = clientapi.ExecutorExtraArg{
			ExecutorNoExecutor: &types.UNIT{},
		}
	}

	// --- Build sendArgs ---
	canton2Any := clientapi.Canton2AnyMessage{
		Receiver: types.TEXT(msg.Receiver),
		Payload:  types.TEXT(msg.Payload),
		FeeToken: *feeTokenInstrumentId,
		ExtraArgs: clientapi.ExtraArgs{
			V3: &clientapi.GenericExtraArgsV3{
				GasLimit:      types.INT64(msg.GasLimit),
				Ccvs:          nil,
				Executor:      executorExtraArg,
				TokenReceiver: "",
				TokenArgs:     "",
			},
		},
	}
	if withToken {
		canton2Any.TokenTransfer = &clientapi.TokenTransfer{
			Token:  *tokenTransferInstrumentId,
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
			CcvAddress: ccvSendDisclosure.Address.Binding(),
			CcvCid:     types.CONTRACT_ID(ccvSendDisclosure.ContractId),
			Context:    splice_api_token_metadata_v1.ChoiceContext{},
		}},
		ExecutorInput: executorInput,
	}
	if withToken {
		sendArgs.TokenTransferInput = &sender.TokenTransferInput{
			SenderInputCids: tokenTransferInputCids,
			TokenPoolCid:    types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
			Context:         tokenPoolSendDisclosure.ChoiceContext,
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
	tx, err := cantonops.CantonSubmit(
		ctx,
		b.Participant,
		useLedger,
		[]*apiv2.Command{{
			Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
				TemplateId:     contracts.IdentifierFromBinding(sender.CCIPSender{}),
				ContractId:     senderCid,
				Choice:         "Send",
				ChoiceArgument: ledger.MapToValue(sendArgs),
			}},
		}},
		allDisclosures,
		nil,
		packageSelectionPreferences,
	)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	fmt.Printf("Message sent in Update: %s\n", tx.GetUpdateId())
	fmt.Println(b.CantonExplorerLink(tx.GetUpdateId()))

	messageId, err := cantonops.GetMessageIdFromTransaction(tx)
	if err != nil {
		return err
	}
	fmt.Printf("✅ Message sent with MessageID: 0x%s\n", messageId)
	fmt.Println(b.CCIPExplorerLink(messageId))

	return nil
}
