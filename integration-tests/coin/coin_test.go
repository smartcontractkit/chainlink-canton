package tests

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func TestCoin(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(5))

	version, err := env.Chain.Participants[0].LedgerServices.Version.GetLedgerApiVersion(t.Context(), &apiv2.GetLedgerApiVersionRequest{})
	require.NoError(t, err)
	fmt.Println(version.Version)

	// Upload the DARs to all participants
	coinDar, err := contracts.GetDar(contracts.Coin, contracts.CurrentVersion)
	require.NoError(t, err)
	packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{coinDar}, env.Chain.Participants...)
	require.NoError(t, err)
	fmt.Printf("Uploaded coin DARs to all participants: %v\n", packageIDs)

	partyAlice := env.Chain.Participants[0].PartyID
	partyBob := env.Chain.Participants[1].PartyID
	partyCharlie := env.Chain.Participants[2].PartyID
	partyDave := env.Chain.Participants[3].PartyID
	partyErin := env.Chain.Participants[4].PartyID

	fmt.Println("Parties:")
	fmt.Printf(" - Alice: %s\n", partyAlice)
	fmt.Printf(" - Bob: %s\n", partyBob)
	fmt.Printf(" - Charlie: %s\n", partyCharlie)
	fmt.Printf(" - Dave: %s\n", partyDave)
	fmt.Printf(" - Erin: %s\n", partyErin)

	// Alice is the issuer, creating registry contract
	instrumentId := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "admin",
			Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
		}, {
			Label: "id",
			Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "LINK"}},
		},
	}}}}
	res, err := env.Chain.Participants[0].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "CoinRegistry",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "issuer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "instanceId",
									Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: contracts.MustNewInstanceID("coinregistry").String()}},
								}, {
									Label: "meta",
									Value: emptyMetadata,
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyAlice},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	registryCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			registryCid = e.Created.ContractId
		}
	}
	fmt.Printf("Deployed registry, CID: %v\n", registryCid)

	// Query the contract for explicit disclosure
	disclosedRegistry, err := testhelpers.GetDisclosedContractById(t.Context(), env.Chain.Participants[0], registryCid)
	require.NoError(t, err)
	fmt.Printf("Queried registry for disclosure: %v\n", disclosedRegistry.GetContractId())

	// Bob creates MintPreapproval
	res, err = env.Chain.Participants[1].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintPreapproval",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
								}, {
									Label: "sender",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyBob},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	bobMintPreapprovalCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			bobMintPreapprovalCid = e.Created.ContractId
		}
	}
	fmt.Printf("Bob created MintPreapproval, CID: %v\n", bobMintPreapprovalCid)
	time.Sleep(time.Second * 5)

	// Alice mints to Bob
	res, err = env.Chain.Participants[0].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice-api-token-burn-mint-v1",
								ModuleName: "Splice.Api.Token.BurnMintV1",
								EntityName: "BurnMintFactory",
							},
							ContractId: registryCid,
							Choice:     "BurnMintFactory_BurnMint",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "expectedAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "inputHoldingCids",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}},
								}, {
									Label: "outputs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{
											Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{
													Label: "owner",
													Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
												}, {
													Label: "amount",
													Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "42.13"}},
												}, {
													Label: "context",
													Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "values",
															Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: []*apiv2.TextMap_Entry{
																{
																	Key: "mint-preapproval",
																	Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
																		Constructor: "AV_ContractId",
																		Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: bobMintPreapprovalCid}},
																	}}},
																},
															}}}},
														},
													}}}},
												},
											}}},
										},
									}}}},
								}, {
									Label: "extraActors",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{}}}},
								}, {
									Label: "extraArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "context",
											Value: emptyChoiceContext,
										}, {
											Label: "meta",
											Value: emptyMetadata,
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyAlice},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	bobCoinHoldingCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			bobCoinHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Alice minted to Bob, CID: %v\n", bobCoinHoldingCid)
	time.Sleep(time.Second * 5)

	// Bob transfers part of their holdings to charlie
	res, err = env.Chain.Participants[1].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice-api-token-transfer-instruction-v1",
								ModuleName: "Splice.Api.Token.TransferInstructionV1",
								EntityName: "TransferFactory",
							},
							ContractId: registryCid,
							Choice:     "TransferFactory_Transfer",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "expectedAdmin",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "transfer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "sender",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyBob}},
										}, {
											Label: "receiver",
											Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCharlie}},
										}, {
											Label: "amount",
											Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "10.00"}},
										}, {
											Label: "instrumentId",
											Value: instrumentId,
										}, {
											Label: "requestedAt",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}},
										}, {
											Label: "executeBefore",
											Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().Add(time.Hour * 24).UnixMicro()}},
										}, {
											Label: "inputHoldingCids",
											Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
												{
													Sum: &apiv2.Value_ContractId{ContractId: bobCoinHoldingCid},
												},
											}}}},
										}, {
											Label: "meta",
											Value: emptyMetadata,
										},
									}}}},
								}, {
									Label: "extraArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "context",
											Value: emptyChoiceContext,
										}, {
											Label: "meta",
											Value: emptyMetadata,
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyBob},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
			},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	transferInstructionCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			// There are multiple contracts created in this transaction, since the change will be returned to the sender (bob)
			// We're interested in the TransferInstruction, since that's the one being sent to charlie
			if e.Created.GetTemplateId().GetEntityName() == "CoinTransferInstruction" {
				transferInstructionCid = e.Created.ContractId
				break
			}
		}
	}
	fmt.Printf("Bob transferred to Charlie, CID: %v\n", transferInstructionCid)
	time.Sleep(time.Second * 5)

	// Charlie accepts transfer from Bob
	res, err = env.Chain.Participants[2].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#splice-api-token-transfer-instruction-v1",
								ModuleName: "Splice.Api.Token.TransferInstructionV1",
								EntityName: "TransferInstruction",
							},
							ContractId: transferInstructionCid,
							Choice:     "TransferInstruction_Accept",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "extraArgs",
									Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
										{
											Label: "context",
											Value: emptyChoiceContext,
										}, {
											Label: "meta",
											Value: emptyMetadata,
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyCharlie},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
			},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	charlieHoldingCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			charlieHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Charlie accepted transfer, holding CID: %v\n", charlieHoldingCid)
	time.Sleep(time.Second * 5)

	// Alice grants mint rights to Dave
	res, err = env.Chain.Participants[0].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintRole",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "issuer",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								}, {
									Label: "minter",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyDave}},
								}, {
									Label: "registry",
									Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: registryCid}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyAlice},
		},
	})
	require.NoError(t, err)
	daveMintRoleCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			daveMintRoleCid = e.Created.ContractId
		}
	}
	fmt.Printf("Alice granted MintRole to Dave, CID: %v\n", daveMintRoleCid)
	time.Sleep(time.Second * 5)

	// Erin grants MintPreapproval
	res, err = env.Chain.Participants[4].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Create{
						Create: &apiv2.CreateCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintPreapproval",
							},
							CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "receiver",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyErin}},
								}, {
									Label: "sender",
									Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyAlice}},
								},
							}},
						},
					},
				},
			},
			ActAs: []string{partyErin},
		},
		TransactionFormat: nil,
	})
	require.NoError(t, err)
	erinPreApprovalCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			erinPreApprovalCid = e.Created.ContractId
		}
	}
	fmt.Printf("Erin created MintPreapproval, CID: %v\n", erinPreApprovalCid)
	disclosedErinPreApproval, err := testhelpers.GetDisclosedContractById(t.Context(), env.Chain.Participants[4], erinPreApprovalCid)
	require.NoError(t, err)
	fmt.Printf("Queried MintPreapproval for disclosure: %v\n", disclosedErinPreApproval.GetContractId())
	time.Sleep(time.Second * 5)

	// Dave uses the MintRole to mint to Erin

	// Asynchronously, listen to all updates that have Erin on Participant 5 as a stakeholder
	currentOffset, err := testhelpers.GetCurrentOffset(t.Context(), env.Chain.Participants[4].LedgerServices.State)
	require.NoError(t, err)
	updateStream, err := env.Chain.Participants[4].LedgerServices.Update.GetUpdates(t.Context(), &apiv2.GetUpdatesRequest{
		BeginExclusive: currentOffset,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						partyErin: {
							Cumulative: []*apiv2.CumulativeFilter{
								{
									IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{},
								},
							},
						},
					},
				},
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
			},
			IncludeReassignments:  nil,
			IncludeTopologyEvents: nil,
		},
	})
	require.NoError(t, err)
	defer updateStream.CloseSend()
	erinHoldingCid := ""
	var erinHoldingErr error
	go func() {
		for {
			update, err := updateStream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				erinHoldingErr = err
				return
			}
			fmt.Printf("Received update on Participant 5: %v\n", update.GetTransaction())
			for _, event := range update.GetTransaction().GetEvents() {
				if c, ok := event.GetEvent().(*apiv2.Event_Created); ok {
					erinHoldingCid = c.Created.GetContractId()
				}
			}
		}
	}()

	res, err = env.Chain.Participants[3].LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{
				{
					Command: &apiv2.Command_Exercise{
						Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{
								PackageId:  "#coin",
								ModuleName: "Coin.Registry",
								EntityName: "MintRole",
							},
							ContractId: daveMintRoleCid,
							Choice:     "Mint",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{
									Label: "instrumentId",
									Value: instrumentId,
								}, {
									Label: "outputs",
									Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										{
											Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
												{
													Label: "owner",
													Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyErin}},
												}, {
													Label: "amount",
													Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "77.77"}},
												}, {
													Label: "context",
													Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
														{
															Label: "values",
															Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: []*apiv2.TextMap_Entry{
																{
																	Key: "mint-preapproval",
																	Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
																		Constructor: "AV_ContractId",
																		Value:       &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: erinPreApprovalCid}},
																	}}},
																},
															}}}},
														},
													}}}},
												},
											}}},
										},
									}}}},
								},
							}}}},
						},
					},
				},
			},
			ActAs: []string{partyDave},
			DisclosedContracts: []*apiv2.DisclosedContract{
				disclosedRegistry,
				disclosedErinPreApproval,
			},
		},
	})
	require.NoError(t, err)
	// Since the submitting party's (dave) node (participant4) isn't a stakeholder on the created contract (CoinHolding),
	// it doesn't actually receive the CreatedEvent as part of the transaction output.
	// This check here therefore won't result in a contract ID, despite us being a witness on the created contract.
	erinCoinHoldingCid := ""
	for i, event := range res.GetTransaction().GetEvents() {
		fmt.Printf("\tEvent %v: %v\n", i, event.GetEvent())
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			erinCoinHoldingCid = e.Created.ContractId
		}
	}
	fmt.Printf("Dave minted to Erin, CID: %v\n", erinCoinHoldingCid) // Will be empty
	fmt.Printf("Transaction: %+v\n", res.GetTransaction())

	// Wait for a couple of seconds, to receive the update on participant 5
	time.Sleep(time.Second * 3)
	require.NoError(t, erinHoldingErr)
	fmt.Printf("Received CoinHolding creation update on Participant 5, CID: %v\n", erinHoldingCid)
	require.NotEmpty(t, erinHoldingCid)
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}

var emptyChoiceContext = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
