package sourcereader

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"testing"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/internal/mocks"
)

// Ledger field labels used when building CreatedEvent / receipt records in tests.
const (
	ccipMessageSentCCIPOwnerLabel = "ccipOwner"
	ccipMessageSentEventLabel     = "event"

	ccipMessageSentEventDestChainSelectorLabel = "destChainSelector"
	ccipMessageSentEventSequenceNumberLabel    = "sequenceNumber"
	ccipMessageSentEventMessageIDLabel         = "messageId"
	ccipMessageSentEventEncodedMessageLabel    = "encodedMessage"
	ccipMessageSentEventVerifierBlobsLabel     = "verifierBlobs"
	ccipMessageSentEventReceiptsLabel          = "receipts"

	ccipMessageSentEventReceiptIssuerTypeLabel        = "issuerType"
	ccipMessageSentEventReceiptIssuerAddressLabel     = "issuerAddress"
	ccipMessageSentEventReceiptDestGasLimitLabel      = "destGasLimit"
	ccipMessageSentEventReceiptDestBytesOverheadLabel = "destBytesOverhead"
	ccipMessageSentEventReceiptFeeTokenAmountLabel    = "feeTokenAmount"
	ccipMessageSentEventReceiptExtraArgsLabel         = "extraArgs"
)

func TestSourceReader_LatestAndFinalizedBlock(t *testing.T) {
	t.Parallel()
	t.Run("returns latest and finalized headers", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		offset := int64(42)

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(
			mock.Anything,
			mock.Anything,
		).Return(&ledgerv2.GetLedgerEndResponse{Offset: offset}, nil)

		reader := &sourceReader{
			stateServiceClient: stateClient,
		}

		latest, finalized, err := reader.LatestAndFinalizedBlock(ctx)
		offsetUint64 := uint64(offset)
		require.NoError(t, err)
		require.NotNil(t, latest)
		require.NotNil(t, finalized)
		require.Equal(t, offsetUint64, latest.Number)
		require.Equal(t, intToBytes32(offsetUint64), latest.Hash)
		require.Equal(t, parentHash(offsetUint64), latest.ParentHash)
		require.Equal(t, *latest, *finalized)
	})

	t.Run("surfaces ledger end error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		stateClient := mocks.NewMockStateServiceClient(t)
		expectedErr := errors.New("boom")

		stateClient.EXPECT().GetLedgerEnd(
			mock.Anything,
			mock.Anything,
		).Return((*ledgerv2.GetLedgerEndResponse)(nil), expectedErr)

		reader := &sourceReader{
			stateServiceClient: stateClient,
		}

		latest, finalized, err := reader.LatestAndFinalizedBlock(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get ledger end")
		require.Nil(t, latest)
		require.Nil(t, finalized)
	})
}

func TestSourceReader_GetBlocksHeaders(t *testing.T) {
	t.Run("builds headers for requested blocks", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		stateClient := mocks.NewMockStateServiceClient(t)

		stateClient.EXPECT().GetLedgerEnd(
			mock.Anything,
			mock.Anything,
		).Return(&ledgerv2.GetLedgerEndResponse{Offset: 10}, nil)

		reader := &sourceReader{
			stateServiceClient: stateClient,
		}

		blockZero := big.NewInt(0)
		blockZeroUint64 := blockZero.Uint64()
		blockFive := big.NewInt(5)
		blockFiveUint64 := blockFive.Uint64()
		headers, err := reader.GetBlocksHeaders(ctx, []*big.Int{blockZero, blockFive})
		require.NoError(t, err)
		require.Len(t, headers, 2)
		require.Equal(t, uint64(0), headers[blockZeroUint64].Number)
		require.Equal(t, protocol.Bytes32{}, headers[blockZeroUint64].ParentHash)
		require.Equal(t, intToBytes32(0), headers[blockZeroUint64].Hash)
		require.Equal(t, uint64(5), headers[blockFiveUint64].Number)
		require.Equal(t, intToBytes32(4), headers[blockFiveUint64].ParentHash)
		require.Equal(t, intToBytes32(5), headers[blockFiveUint64].Hash)
	})

	t.Run("errors when block exceeds latest offset", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		stateClient := mocks.NewMockStateServiceClient(t)

		stateClient.EXPECT().GetLedgerEnd(
			mock.Anything,
			mock.Anything,
		).Return(&ledgerv2.GetLedgerEndResponse{Offset: 3}, nil)

		reader := &sourceReader{
			stateServiceClient: stateClient,
		}

		_, err := reader.GetBlocksHeaders(ctx, []*big.Int{big.NewInt(4)})
		require.Error(t, err)
		require.ErrorContains(t, err, "block number is greater than latest offset")
	})
}

func TestSourceReader_GetRMNCursedSubjects(t *testing.T) {
	t.Parallel()
	const (
		nopParty      = "node-operator-party"
		rmnOwner      = "owner"
		ccipOwner     = "ccip-owner"
		rmnInstanceID = "rmn-1"
	)
	var (
		rmnRemoteTemplateID = contracts.TemplateID{
			PackageID:  fmt.Sprintf("#%s", rmn.PackageName),
			ModuleName: "CCIP.RMNRemote",
			EntityName: "RMNRemote",
		}
		rmnInstanceAddress = contracts.InstanceID(rmnInstanceID).RawInstanceAddress(rmnOwner).InstanceAddress()
	)

	t.Run("returns cursed subjects from first active RMNRemote", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		// Two subjects: with and without 0x prefix (code adds 0x if missing).
		subject1Hex := "0102030405060708090a0b0c0d0e0f10"
		subject2Hex := "0x1112131415161718191a1b1c1d1e1f20"

		createdEvent := &ledgerv2.CreatedEvent{
			TemplateId: rmnRemoteTemplateID.ToLedgerIdentifier(),
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: rmnInstanceID}}},
					{Label: "rmnOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: rmnOwner}}},
					{Label: "ccipOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}}},
					{Label: "customObservers", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
					{Label: "cursedSubjects", Value: &ledgerv2.Value{
						Sum: &ledgerv2.Value_List{
							List: &ledgerv2.List{Elements: []*ledgerv2.Value{
								{Sum: &ledgerv2.Value_Text{Text: subject1Hex}},
								{Sum: &ledgerv2.Value_Text{Text: subject2Hex}},
							}},
						},
					}},
				},
			},
			Signatories: []string{rmnOwner},
		}

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.MatchedBy(func(req *ledgerv2.GetActiveContractsRequest) bool {
			return req.ActiveAtOffset == 5 &&
				req.EventFormat != nil &&
				req.EventFormat.FiltersByParty[nopParty] != nil
		}), mock.Anything).
			Return(&fakeActiveContractsStream{
				ctx: ctx,
				responses: []*ledgerv2.GetActiveContractsResponse{
					{
						ContractEntry: &ledgerv2.GetActiveContractsResponse_ActiveContract{
							ActiveContract: &ledgerv2.ActiveContract{
								CreatedEvent: createdEvent,
							},
						},
					},
				},
			}, nil)

		reader := &sourceReader{
			stateServiceClient: stateClient,
			config: ReaderConfig{
				NodeOperatorParty:   nopParty,
				RMNRemoteTemplateID: rmnRemoteTemplateID,
			},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		subjects, err := reader.GetRMNCursedSubjects(ctx)
		require.NoError(t, err)
		require.Len(t, subjects, 2)

		expected1, err := protocol.NewBytes16FromString("0x" + subject1Hex)
		require.NoError(t, err)
		expected2, err := protocol.NewBytes16FromString(subject2Hex)
		require.NoError(t, err)
		require.Equal(t, expected1, subjects[0])
		require.Equal(t, expected2, subjects[1])
	})

	t.Run("returns empty list when RMNRemote has no cursed subjects", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		createdEvent := &ledgerv2.CreatedEvent{
			TemplateId: rmnRemoteTemplateID.ToLedgerIdentifier(),
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: rmnInstanceID}}},
					{Label: "rmnOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: rmnOwner}}},
					{Label: "ccipOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}}},
					{Label: "customObservers", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
					{Label: "cursedSubjects", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
				},
			},
			Signatories: []string{rmnOwner},
		}

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx, responses: []*ledgerv2.GetActiveContractsResponse{
				{ContractEntry: &ledgerv2.GetActiveContractsResponse_ActiveContract{
					ActiveContract: &ledgerv2.ActiveContract{CreatedEvent: createdEvent},
				}},
			}}, nil)

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{NodeOperatorParty: nopParty, RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		subjects, err := reader.GetRMNCursedSubjects(ctx)
		require.NoError(t, err)
		require.Empty(t, subjects)
	})

	t.Run("returns error when no active RMNRemote found", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx}, nil) // no responses, EOF immediately

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{NodeOperatorParty: nopParty, RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		_, err := reader.GetRMNCursedSubjects(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "no active contract found for InstanceAddress")
	})

	t.Run("surfaces LatestAndFinalizedBlock error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return((*ledgerv2.GetLedgerEndResponse)(nil), errors.New("ledger end failed"))

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		_, err := reader.GetRMNCursedSubjects(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get latest block")
	})

	t.Run("surfaces GetActiveContracts error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return((grpc.ServerStreamingClient[ledgerv2.GetActiveContractsResponse])(nil), errors.New("get active contracts failed"))

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{NodeOperatorParty: nopParty, RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		_, err := reader.GetRMNCursedSubjects(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get active contracts")
	})

	t.Run("surfaces Recv error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx, err: errors.New("recv failed")}, nil)

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{NodeOperatorParty: nopParty, RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		_, err := reader.GetRMNCursedSubjects(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to receive active contract")
	})

	t.Run("returns error when subject is not valid hex", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		createdEvent := &ledgerv2.CreatedEvent{
			TemplateId: rmnRemoteTemplateID.ToLedgerIdentifier(),
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: rmnInstanceID}}},
					{Label: "rmnOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: rmnOwner}}},
					{Label: "ccipOwner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}}},
					{Label: "customObservers", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_List{List: &ledgerv2.List{Elements: []*ledgerv2.Value{}}}}},
					{Label: "cursedSubjects", Value: &ledgerv2.Value{
						Sum: &ledgerv2.Value_List{
							List: &ledgerv2.List{Elements: []*ledgerv2.Value{
								{Sum: &ledgerv2.Value_Text{Text: "not-valid-hex!!"}},
							}},
						},
					}},
				},
			},
			Signatories: []string{rmnOwner},
		}

		stateClient := mocks.NewMockStateServiceClient(t)
		stateClient.EXPECT().GetLedgerEnd(mock.Anything, mock.Anything).
			Return(&ledgerv2.GetLedgerEndResponse{Offset: 5}, nil)
		stateClient.EXPECT().GetActiveContracts(mock.Anything, mock.Anything, mock.Anything).
			Return(&fakeActiveContractsStream{ctx: ctx, responses: []*ledgerv2.GetActiveContractsResponse{
				{ContractEntry: &ledgerv2.GetActiveContractsResponse_ActiveContract{
					ActiveContract: &ledgerv2.ActiveContract{CreatedEvent: createdEvent},
				}},
			}}, nil)

		reader := &sourceReader{
			stateServiceClient:       stateClient,
			config:                   ReaderConfig{NodeOperatorParty: nopParty, RMNRemoteTemplateID: rmnRemoteTemplateID},
			rmnRemoteInstanceAddress: rmnInstanceAddress,
		}

		_, err := reader.GetRMNCursedSubjects(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to decode subject")
	})
}

func TestSourceReader_FetchMessageSentEvents(t *testing.T) {
	t.Parallel()
	const ccipOwner = "owner-party"
	const nopParty = "node-operator-party"
	var (
		templateID = contracts.TemplateID{
			PackageID:  common.PackageName,
			ModuleName: "CCIP.Events",
			EntityName: "CCIPMessageSent",
		}
	)

	t.Run("returns error when fromBlock is nil", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		events, err := reader.FetchMessageSentEvents(ctx, nil, big.NewInt(1))
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "fromBlock is nil")
	})

	t.Run("returns error when fromBlock is negative", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(-1), big.NewInt(5))
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "fromBlock is negative")
	})

	t.Run("returns error when toBlock is negative", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(-3))
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "toBlock is negative")
	})

	t.Run("returns error when toBlock exceeds max Int64", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		toBlock := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(1))
		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), toBlock)
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "toBlock is larger than the max Int64 value")
		require.ErrorContains(t, err, toBlock.String())
	})

	t.Run("returns error when fromBlock is greater than toBlock", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(10), big.NewInt(5))
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "fromBlock is greater than toBlock")
	})

	t.Run("returns error when exclusive begin offset exceeds max Int64", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		reader := &sourceReader{}

		fromBlock := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(2))
		events, err := reader.FetchMessageSentEvents(ctx, fromBlock, nil)
		require.Nil(t, events)
		require.Error(t, err)
		require.ErrorContains(t, err, "exclusive begin offset derived from fromBlock exceeds max Int64")
	})

	t.Run("ignores event when ccipOwner does not match", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		msg, err := protocol.NewMessage(
			protocol.ChainSelector(1),
			protocol.ChainSelector(2),
			protocol.SequenceNumber(7),
			protocol.UnknownAddress{0x01},
			protocol.UnknownAddress{0x02},
			1,
			100,
			200,
			protocol.Bytes32{},
			protocol.UnknownAddress{0x03},
			protocol.UnknownAddress{0x04},
			[]byte{0xAA},
			[]byte{0xBB},
			nil,
		)
		require.NoError(t, err)

		encodedMsg, err := msg.Encode()
		require.NoError(t, err)
		msgID := msg.MustMessageID()
		msgIDHex := hex.EncodeToString(msgID[:])
		encodedMsgHex := hex.EncodeToString(encodedMsg)

		created := &ledgerv2.CreatedEvent{
			TemplateId: templateID.ToLedgerIdentifier(),
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{
						Label: ccipMessageSentCCIPOwnerLabel,
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: "wrong-owner"}}, // wrong owner, should not get processed
					},
					{
						Label: ccipMessageSentEventLabel,
						Value: &ledgerv2.Value{
							Sum: &ledgerv2.Value_Record{
								Record: &ledgerv2.Record{
									Fields: []*ledgerv2.RecordField{
										{
											Label: ccipMessageSentEventDestChainSelectorLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 2}},
										},
										{
											Label: ccipMessageSentEventSequenceNumberLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 7}},
										},
										{
											Label: ccipMessageSentEventMessageIDLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: msgIDHex}},
										},
										{
											Label: ccipMessageSentEventEncodedMessageLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: encodedMsgHex}},
										},
										{
											Label: ccipMessageSentEventVerifierBlobsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
										{
											Label: ccipMessageSentEventReceiptsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		tx := &ledgerv2.Transaction{
			UpdateId: "0xdeadbeef",
			Offset:   10,
			Events: []*ledgerv2.Event{
				{Event: &ledgerv2.Event_Created{Created: created}},
			},
		}

		stream := &fakeUpdateStream{
			ctx: ctx,
			responses: []*ledgerv2.GetUpdatesResponse{
				{Update: &ledgerv2.GetUpdatesResponse_Transaction{Transaction: tx}},
			},
		}

		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(
			mock.Anything,
			mock.Anything,
		).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(5))
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("ignores event when ccipOwnerParty is not in signatories", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		msg, err := protocol.NewMessage(
			protocol.ChainSelector(1),
			protocol.ChainSelector(2),
			protocol.SequenceNumber(7),
			protocol.UnknownAddress{0x01},
			protocol.UnknownAddress{0x02},
			1,
			100,
			200,
			protocol.Bytes32{},
			protocol.UnknownAddress{0x03},
			protocol.UnknownAddress{0x04},
			[]byte{0xAA},
			[]byte{0xBB},
			nil,
		)
		require.NoError(t, err)

		encodedMsg, err := msg.Encode()
		require.NoError(t, err)
		msgID := msg.MustMessageID()
		msgIDHex := hex.EncodeToString(msgID[:])
		encodedMsgHex := hex.EncodeToString(encodedMsg)

		// Event has correct ccipOwner in CreateArguments but signatories do not include ccipOwnerParty.
		created := &ledgerv2.CreatedEvent{
			TemplateId:  templateID.ToLedgerIdentifier(),
			Signatories: []string{"other-party"}, // ccipOwner not in signatories - should be skipped
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{
						Label: ccipMessageSentCCIPOwnerLabel,
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}},
					},
					{
						Label: ccipMessageSentEventLabel,
						Value: &ledgerv2.Value{
							Sum: &ledgerv2.Value_Record{
								Record: &ledgerv2.Record{
									Fields: []*ledgerv2.RecordField{
										{
											Label: ccipMessageSentEventDestChainSelectorLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 2}},
										},
										{
											Label: ccipMessageSentEventSequenceNumberLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 7}},
										},
										{
											Label: ccipMessageSentEventMessageIDLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: msgIDHex}},
										},
										{
											Label: ccipMessageSentEventEncodedMessageLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: encodedMsgHex}},
										},
										{
											Label: ccipMessageSentEventVerifierBlobsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
										{
											Label: ccipMessageSentEventReceiptsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		tx := &ledgerv2.Transaction{
			UpdateId: "0xdeadbeef",
			Offset:   10,
			Events: []*ledgerv2.Event{
				{Event: &ledgerv2.Event_Created{Created: created}},
			},
		}

		stream := &fakeUpdateStream{
			ctx: ctx,
			responses: []*ledgerv2.GetUpdatesResponse{
				{Update: &ledgerv2.GetUpdatesResponse_Transaction{Transaction: tx}},
			},
		}

		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(5))
		require.NoError(t, err)
		require.Empty(t, events, "event with ccipOwner not in signatories should be ignored")
	})

	t.Run("returns error when stream recv fails", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		updateClient := mocks.NewMockUpdateServiceClient(t)
		stream := &fakeUpdateStream{
			ctx: ctx,
			err: errors.New("recv failed"),
		}

		updateClient.EXPECT().GetUpdates(
			mock.Anything,
			mock.Anything,
		).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		_, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(2))
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to get updates")
	})

	t.Run("uses zero begin exclusive when fromBlock is zero", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		updateClient := mocks.NewMockUpdateServiceClient(t)
		stream := &fakeUpdateStream{ctx: ctx}

		updateClient.EXPECT().GetUpdates(
			mock.Anything,
			mock.MatchedBy(func(req *ledgerv2.GetUpdatesRequest) bool {
				if req.GetBeginExclusive() != 0 {
					return false
				}
				if req.EndInclusive == nil || *req.EndInclusive != 2 {
					return false
				}

				return true
			}),
		).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(0), big.NewInt(2))
		require.NoError(t, err)
		require.Empty(t, events)
	})

	t.Run("parses events with verifier blobs and receipts", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		verifierBlobHex := "deadbeef"
		extraArgsHex := "cafebabe"

		ccvIssuer := protocol.Keccak256([]byte("issuer1"))
		execIssuer := protocol.Keccak256([]byte("executor"))
		networkIssuer := protocol.Keccak256([]byte("network"))

		bindingReceipts := []common.Receipt{
			{IssuerAddress: types.TEXT(hex.EncodeToString(ccvIssuer[:])), DestGasLimit: 100000, DestBytesOverhead: 500, FeeTokenAmount: types.NUMERIC("1000000."), ExtraArgs: types.TEXT(extraArgsHex)},
			{IssuerAddress: types.TEXT(hex.EncodeToString(execIssuer[:])), DestGasLimit: 0, DestBytesOverhead: 0, FeeTokenAmount: types.NUMERIC("500000."), ExtraArgs: types.TEXT("")},
			{IssuerAddress: types.TEXT(hex.EncodeToString(networkIssuer[:])), DestGasLimit: 0, DestBytesOverhead: 0, FeeTokenAmount: types.NUMERIC("500000."), ExtraArgs: types.TEXT("")},
		}
		receiptsWithBlobs, err := receiptsBindingToProtocol(bindingReceipts)
		require.NoError(t, err)
		require.Len(t, receiptsWithBlobs, 3) // 1 verifier blob, 1 executor receipt, 1 network fee receipt

		// Ledger-style receipts list for the CreatedEvent (so UnmarshalCreatedEvent can populate Event.Receipts).
		receiptsList := &ledgerv2.List{Elements: []*ledgerv2.Value{
			{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: ccipMessageSentEventReceiptIssuerTypeLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: "ccv"}}},
					{Label: ccipMessageSentEventReceiptIssuerAddressLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: hex.EncodeToString(ccvIssuer[:])}}},
					{Label: ccipMessageSentEventReceiptDestGasLimitLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 100000}}},
					{Label: ccipMessageSentEventReceiptDestBytesOverheadLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 500}}},
					{Label: ccipMessageSentEventReceiptFeeTokenAmountLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "1000000."}}},
					{Label: ccipMessageSentEventReceiptExtraArgsLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: extraArgsHex}}},
				},
			}}},
			{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: ccipMessageSentEventReceiptIssuerTypeLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: "executor"}}},
					{Label: ccipMessageSentEventReceiptIssuerAddressLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: hex.EncodeToString(execIssuer[:])}}},
					{Label: ccipMessageSentEventReceiptDestGasLimitLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
					{Label: ccipMessageSentEventReceiptDestBytesOverheadLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
					{Label: ccipMessageSentEventReceiptFeeTokenAmountLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "500000."}}},
					{Label: ccipMessageSentEventReceiptExtraArgsLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: ""}}},
				},
			}}},
			{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{Label: ccipMessageSentEventReceiptIssuerTypeLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: "network"}}},
					{Label: ccipMessageSentEventReceiptIssuerAddressLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: hex.EncodeToString(networkIssuer[:])}}},
					{Label: ccipMessageSentEventReceiptDestGasLimitLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
					{Label: ccipMessageSentEventReceiptDestBytesOverheadLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 0}}},
					{Label: ccipMessageSentEventReceiptFeeTokenAmountLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: "500000."}}},
					{Label: ccipMessageSentEventReceiptExtraArgsLabel, Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: ""}}},
				},
			}}},
		}}

		structure, err := protocol.ParseReceiptStructure(receiptsWithBlobs, 1, 0)
		require.NoError(t, err)

		ccvAndExecutorHash, err := protocol.ComputeCCVAndExecutorHash(structure.CCVAddresses, structure.ExecutorAddress)
		require.NoError(t, err)

		msg, err := protocol.NewMessage(
			protocol.ChainSelector(1),
			protocol.ChainSelector(2),
			protocol.SequenceNumber(7),
			protocol.UnknownAddress{0x01},
			protocol.UnknownAddress{0x02},
			1,
			100,
			200,
			ccvAndExecutorHash,
			protocol.UnknownAddress{0x03},
			protocol.UnknownAddress{0x04},
			[]byte{0xAA},
			[]byte{0xBB},
			nil,
		)
		require.NoError(t, err)

		encodedMsg, err := msg.Encode()
		require.NoError(t, err)
		msgID := msg.MustMessageID()
		msgIDHex := hex.EncodeToString(msgID[:])
		encodedMsgHex := hex.EncodeToString(encodedMsg)

		created := &ledgerv2.CreatedEvent{
			TemplateId:  templateID.ToLedgerIdentifier(),
			Signatories: []string{ccipOwner},
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{
						Label: ccipMessageSentCCIPOwnerLabel,
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}},
					},
					{
						Label: ccipMessageSentEventLabel,
						Value: &ledgerv2.Value{
							Sum: &ledgerv2.Value_Record{
								Record: &ledgerv2.Record{
									Fields: []*ledgerv2.RecordField{
										{
											Label: ccipMessageSentEventDestChainSelectorLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 2}},
										},
										{
											Label: ccipMessageSentEventSequenceNumberLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 7}},
										},
										{
											Label: ccipMessageSentEventMessageIDLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: msgIDHex}},
										},
										{
											Label: ccipMessageSentEventEncodedMessageLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: encodedMsgHex}},
										},
										{
											Label: ccipMessageSentEventVerifierBlobsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{
														{Sum: &ledgerv2.Value_Text{Text: verifierBlobHex}},
													}},
												},
											},
										},
										{
											Label: ccipMessageSentEventReceiptsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: receiptsList,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		tx := &ledgerv2.Transaction{
			UpdateId: "0xdeadbeef",
			Offset:   10,
			Events: []*ledgerv2.Event{
				{Event: &ledgerv2.Event_Created{Created: created}},
			},
		}

		stream := &fakeUpdateStream{
			ctx: ctx,
			responses: []*ledgerv2.GetUpdatesResponse{
				{Update: &ledgerv2.GetUpdatesResponse_Transaction{Transaction: tx}},
			},
		}

		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		events, err := reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(5))
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, msg.MustMessageID(), events[0].MessageID)

		// Assert receipts were properly parsed
		require.Len(t, events[0].Receipts, 3)

		// First receipt - should have verifier blob populated
		receipt1 := events[0].Receipts[0]
		require.Equal(t, protocol.UnknownAddress(ccvIssuer[:]), receipt1.Issuer)
		require.Equal(t, uint64(100000), receipt1.DestGasLimit)
		require.Equal(t, uint32(500), receipt1.DestBytesOverhead)
		require.Equal(t, big.NewInt(1000000), receipt1.FeeTokenAmount)
		expectedExtraArgs, _ := hex.DecodeString(extraArgsHex)
		require.Equal(t, protocol.ByteSlice(expectedExtraArgs), receipt1.ExtraArgs)
		expectedBlob, _ := hex.DecodeString(verifierBlobHex)
		require.Equal(t, protocol.ByteSlice(expectedBlob), receipt1.Blob)

		// Second receipt - no verifier blob (executor fee receipt)
		receipt2 := events[0].Receipts[1]
		require.Equal(t, protocol.UnknownAddress(execIssuer[:]), receipt2.Issuer)
		require.Equal(t, uint64(0), receipt2.DestGasLimit)
		require.Equal(t, uint32(0), receipt2.DestBytesOverhead)
		require.Equal(t, big.NewInt(500000), receipt2.FeeTokenAmount)
		require.Empty(t, receipt2.ExtraArgs)
		require.Nil(t, receipt2.Blob) // No corresponding verifier blob

		// Third receipt - no verifier blob (network fee receipt)
		receipt3 := events[0].Receipts[2]
		require.Equal(t, protocol.UnknownAddress(networkIssuer[:]), receipt3.Issuer)
		require.Equal(t, uint64(0), receipt3.DestGasLimit)
		require.Equal(t, uint32(0), receipt3.DestBytesOverhead)
		require.Equal(t, big.NewInt(500000), receipt3.FeeTokenAmount)
		require.Empty(t, receipt3.ExtraArgs)
		require.Nil(t, receipt3.Blob) // No corresponding verifier blob
	})

	t.Run("returns error when receipts fewer than verifier blobs", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		msg, err := protocol.NewMessage(
			protocol.ChainSelector(1),
			protocol.ChainSelector(2),
			protocol.SequenceNumber(7),
			protocol.UnknownAddress{0x01},
			protocol.UnknownAddress{0x02},
			1,
			100,
			200,
			protocol.Bytes32{},
			protocol.UnknownAddress{0x03},
			protocol.UnknownAddress{0x04},
			[]byte{0xAA},
			[]byte{0xBB},
			nil,
		)
		require.NoError(t, err)

		encodedMsg, err := msg.Encode()
		require.NoError(t, err)
		msgID := msg.MustMessageID()
		msgIDHex := hex.EncodeToString(msgID[:])
		encodedMsgHex := hex.EncodeToString(encodedMsg)

		// Two verifier blobs but zero receipts - should fail
		created := &ledgerv2.CreatedEvent{
			TemplateId:  templateID.ToLedgerIdentifier(),
			Signatories: []string{ccipOwner},
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{
						Label: ccipMessageSentCCIPOwnerLabel,
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}},
					},
					{
						Label: ccipMessageSentEventLabel,
						Value: &ledgerv2.Value{
							Sum: &ledgerv2.Value_Record{
								Record: &ledgerv2.Record{
									Fields: []*ledgerv2.RecordField{
										{
											Label: ccipMessageSentEventDestChainSelectorLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 2}},
										},
										{
											Label: ccipMessageSentEventSequenceNumberLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 7}},
										},
										{
											Label: ccipMessageSentEventMessageIDLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: msgIDHex}},
										},
										{
											Label: ccipMessageSentEventEncodedMessageLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: encodedMsgHex}},
										},
										{
											Label: ccipMessageSentEventVerifierBlobsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{
														{Sum: &ledgerv2.Value_Text{Text: "deadbeef"}},
														{Sum: &ledgerv2.Value_Text{Text: "cafebabe"}},
													}},
												},
											},
										},
										{
											Label: ccipMessageSentEventReceiptsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		tx := &ledgerv2.Transaction{
			UpdateId: "0xdeadbeef",
			Offset:   10,
			Events: []*ledgerv2.Event{
				{Event: &ledgerv2.Event_Created{Created: created}},
			},
		}

		stream := &fakeUpdateStream{
			ctx: ctx,
			responses: []*ledgerv2.GetUpdatesResponse{
				{Update: &ledgerv2.GetUpdatesResponse_Transaction{Transaction: tx}},
			},
		}

		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		_, err = reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(5))
		require.Error(t, err)
		require.ErrorContains(t, err, "expected more receipts than verifier blobs")
	})

	t.Run("returns error on unknown receipt field", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		msg, err := protocol.NewMessage(
			protocol.ChainSelector(1),
			protocol.ChainSelector(2),
			protocol.SequenceNumber(7),
			protocol.UnknownAddress{0x01},
			protocol.UnknownAddress{0x02},
			1,
			100,
			200,
			protocol.Bytes32{},
			protocol.UnknownAddress{0x03},
			protocol.UnknownAddress{0x04},
			[]byte{0xAA},
			[]byte{0xBB},
			nil,
		)
		require.NoError(t, err)

		encodedMsg, err := msg.Encode()
		require.NoError(t, err)
		msgID := msg.MustMessageID()
		msgIDHex := hex.EncodeToString(msgID[:])
		encodedMsgHex := hex.EncodeToString(encodedMsg)

		created := &ledgerv2.CreatedEvent{
			TemplateId:  templateID.ToLedgerIdentifier(),
			Signatories: []string{ccipOwner},
			CreateArguments: &ledgerv2.Record{
				Fields: []*ledgerv2.RecordField{
					{
						Label: ccipMessageSentCCIPOwnerLabel,
						Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: ccipOwner}},
					},
					{
						Label: ccipMessageSentEventLabel,
						Value: &ledgerv2.Value{
							Sum: &ledgerv2.Value_Record{
								Record: &ledgerv2.Record{
									Fields: []*ledgerv2.RecordField{
										{
											Label: ccipMessageSentEventDestChainSelectorLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 2}},
										},
										{
											Label: ccipMessageSentEventSequenceNumberLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: 7}},
										},
										{
											Label: ccipMessageSentEventMessageIDLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: msgIDHex}},
										},
										{
											Label: ccipMessageSentEventEncodedMessageLabel,
											Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: encodedMsgHex}},
										},
										{
											Label: ccipMessageSentEventVerifierBlobsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{}},
												},
											},
										},
										{
											Label: ccipMessageSentEventReceiptsLabel,
											Value: &ledgerv2.Value{
												Sum: &ledgerv2.Value_List{
													List: &ledgerv2.List{Elements: []*ledgerv2.Value{
														{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{
															Fields: []*ledgerv2.RecordField{
																{Label: "unknownField", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: "value"}}},
															},
														}}},
													}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		tx := &ledgerv2.Transaction{
			UpdateId: "0xdeadbeef",
			Offset:   10,
			Events: []*ledgerv2.Event{
				{Event: &ledgerv2.Event_Created{Created: created}},
			},
		}

		stream := &fakeUpdateStream{
			ctx: ctx,
			responses: []*ledgerv2.GetUpdatesResponse{
				{Update: &ledgerv2.GetUpdatesResponse_Transaction{Transaction: tx}},
			},
		}

		updateClient := mocks.NewMockUpdateServiceClient(t)
		updateClient.EXPECT().GetUpdates(mock.Anything, mock.Anything).Return(stream, nil)

		reader := &sourceReader{
			updateServiceClient: updateClient,
			config: ReaderConfig{
				NodeOperatorParty:         nopParty,
				CCIPOwnerParty:            ccipOwner,
				CCIPMessageSentTemplateID: templateID,
			},
		}

		_, err = reader.FetchMessageSentEvents(ctx, big.NewInt(1), big.NewInt(5))
		require.Error(t, err)
		// With UnmarshalCreatedEvent, receipts are unmarshaled into binding structs; a receipt
		// with only unknown fields has zero-valued required fields, so we fail on parse (e.g. fee token amount).
		require.ErrorContains(t, err, "failed to process receipts")
	})
}

type fakeUpdateStream struct {
	ctx       context.Context //nolint:containedctx
	responses []*ledgerv2.GetUpdatesResponse
	err       error
	idx       int
}

func (s *fakeUpdateStream) Recv() (*ledgerv2.GetUpdatesResponse, error) {
	if s.idx < len(s.responses) {
		resp := s.responses[s.idx]
		s.idx++

		return resp, nil
	}
	if s.err != nil {
		return nil, s.err
	}

	return nil, io.EOF
}

func (s *fakeUpdateStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s *fakeUpdateStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s *fakeUpdateStream) CloseSend() error {
	return nil
}

func (s *fakeUpdateStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}

	return context.Background()
}

func (s *fakeUpdateStream) SendMsg(any) error {
	return nil
}

func (s *fakeUpdateStream) RecvMsg(any) error {
	return nil
}

var _ grpc.ServerStreamingClient[ledgerv2.GetUpdatesResponse] = (*fakeUpdateStream)(nil)

// fakeActiveContractsStream implements grpc.ServerStreamingClient[ledgerv2.GetActiveContractsResponse] for tests.
type fakeActiveContractsStream struct {
	ctx       context.Context //nolint:containedctx
	responses []*ledgerv2.GetActiveContractsResponse
	err       error
	idx       int
}

func (s *fakeActiveContractsStream) Recv() (*ledgerv2.GetActiveContractsResponse, error) {
	if s.idx < len(s.responses) {
		resp := s.responses[s.idx]
		s.idx++

		return resp, nil
	}
	if s.err != nil {
		return nil, s.err
	}

	return nil, io.EOF
}

func (s *fakeActiveContractsStream) Header() (metadata.MD, error) {
	return metadata.MD{}, nil
}

func (s *fakeActiveContractsStream) Trailer() metadata.MD {
	return metadata.MD{}
}

func (s *fakeActiveContractsStream) CloseSend() error {
	return nil
}

func (s *fakeActiveContractsStream) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}

	return context.Background()
}

func (s *fakeActiveContractsStream) SendMsg(any) error {
	return nil
}

func (s *fakeActiveContractsStream) RecvMsg(any) error {
	return nil
}

var _ grpc.ServerStreamingClient[ledgerv2.GetActiveContractsResponse] = (*fakeActiveContractsStream)(nil)
