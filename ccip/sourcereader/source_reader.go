package sourcereader

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"slices"
	"strings"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-ccv/pkg/chainaccess"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// ReaderConfig is the configuration required to create a canton source reader.
type ReaderConfig struct {
	// NodeOperatorParty is the observer party that node operators will use to observe CCIPMessageSent events.
	// This is used to filter out events that are not observed by the node operator.
	NodeOperatorParty string `toml:"node_operator_party"`

	// CCIPOwnerParty is the party that we expect to be present in the CCIPMessageSent.ccipOwner field.
	// This proves that the ccipOwner is a signatory on the CCIPMessageSent contract(event).
	CCIPOwnerParty string `toml:"ccip_owner_party"`

	// CCIPMessageSentTemplateID is the template ID of the CCIPMessageSent contract.
	CCIPMessageSentTemplateID contracts.TemplateID `toml:"ccip_message_sent_template_id"`

	// RMNRemoteTemplateID is the template ID of the RMNRemote contract.
	RMNRemoteTemplateID contracts.TemplateID `toml:"rmn_remote_template_id"`

	// Authority is the authority to use for the gRPC connection.
	// Connecting to the gRPC API via nginx usually requires this to be set.
	Authority string `toml:"authority"`
}

type sourceReader struct {
	lggr                     logger.Logger
	stateServiceClient       ledgerv2.StateServiceClient
	updateServiceClient      ledgerv2.UpdateServiceClient
	rmnRemoteInstanceAddress contracts.InstanceAddress

	config ReaderConfig
}

// NewSourceReader creates a new canton source reader.
// Auth is handled via gRPC dial options (WithTransportCredentials, WithPerRPCCredentials).
func NewSourceReader(
	lggr logger.Logger,
	grpcEndpoint string,
	config ReaderConfig,
	rmnRemoteInstanceAddress contracts.InstanceAddress,
	opts ...grpc.DialOption,
) (chainaccess.SourceReader, error) {
	lggr.Infow("creating gRPC connection to canton node", "grpcEndpoint", grpcEndpoint, "config", config)

	if config.Authority != "" {
		opts = append(opts, grpc.WithAuthority(config.Authority))
	}

	conn, err := grpc.NewClient(grpcEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to canton node: %w", err)
	}

	return &sourceReader{
		lggr:                     lggr,
		stateServiceClient:       ledgerv2.NewStateServiceClient(conn),
		updateServiceClient:      ledgerv2.NewUpdateServiceClient(conn),
		config:                   config,
		rmnRemoteInstanceAddress: rmnRemoteInstanceAddress,
	}, nil
}

// FetchMessageSentEvents implements chainaccess.SourceReader.
//
// It streams ledger updates in the requested offset range and collects CCIPMessageSent created
// events that pass the configured node-operator filter, template id match, and signatory
// checks (including ccipOwner alignment with the configured CCIP owner party).
//
// Resilience: after those checks, if decoding or validation fails for a single event (message
// payload, receipts, verifier blobs, transaction update id, cross-field checks, etc.), the
// failure is logged at error level with ledger context and that event is skipped. Other events
// in the same range are still returned. Metadata mismatches (wrong owner / template) are
// skipped without an error log. Transport and RPC errors (GetUpdates, Recv, GetLedgerEnd
// when toBlock is nil) still fail the entire call.
func (c *sourceReader) FetchMessageSentEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]protocol.MessageSentEvent, error) {
	if fromBlock == nil {
		return nil, fmt.Errorf("fromBlock is nil, it must be provided")
	}

	// Check that neither fromBlock nor toBlock are negative.
	if fromBlock.Sign() < 0 {
		return nil, fmt.Errorf("fromBlock is negative, cannot proceed: %s", fromBlock.String())
	}
	if toBlock != nil && toBlock.Sign() < 0 {
		return nil, fmt.Errorf("toBlock is negative, cannot proceed: %s", toBlock.String())
	}

	maxInt64 := big.NewInt(math.MaxInt64)
	if toBlock != nil {
		if fromBlock.Cmp(toBlock) > 0 {
			return nil, fmt.Errorf(
				"fromBlock is greater than toBlock, cannot proceed: %s > %s",
				fromBlock.String(),
				toBlock.String(),
			)
		}
		if toBlock.Cmp(maxInt64) > 0 {
			return nil, fmt.Errorf(
				"toBlock is larger than the max Int64 value, cannot proceed: %s > %d",
				toBlock.String(),
				math.MaxInt64,
			)
		}
	}

	// since begin is exclusive we need to subtract 1 from fromBlock
	begin := new(big.Int).Sub(fromBlock, big.NewInt(1))
	// check that begin is not negative
	if begin.Sign() < 0 {
		begin = big.NewInt(0)
	}
	if begin.Cmp(maxInt64) > 0 {
		return nil, fmt.Errorf(
			"exclusive begin offset derived from fromBlock exceeds max Int64, cannot proceed: %s > %d",
			begin.String(),
			math.MaxInt64,
		)
	}

	var end *int64
	if toBlock != nil {
		// Safe to convert: toBlock is non-negative and <= max Int64.
		end = new(toBlock.Int64())
	} else {
		// If toBlock is nil, we need to get the latest ledger end to avoid streaming indefinitely
		// and to ensure we return a slice as expected by the interface.
		ledgerEnd, err := c.stateServiceClient.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to get ledger end for open-ended query: %w", err)
		}
		end = new(ledgerEnd.GetOffset())
	}

	ccipMessageSentIdentifier := c.config.CCIPMessageSentTemplateID.ToLedgerIdentifier()
	updates, err := c.updateServiceClient.GetUpdates(ctx, &ledgerv2.GetUpdatesRequest{
		BeginExclusive: begin.Int64(), // safe to convert: begin is non-negative and <= max Int64.
		EndInclusive:   end,
		UpdateFormat: &ledgerv2.UpdateFormat{
			IncludeTransactions: &ledgerv2.TransactionFormat{
				TransactionShape: ledgerv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				EventFormat: &ledgerv2.EventFormat{
					FiltersByParty: map[string]*ledgerv2.Filters{
						c.config.NodeOperatorParty: {
							Cumulative: []*ledgerv2.CumulativeFilter{
								{
									IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
										TemplateFilter: &ledgerv2.TemplateFilter{
											TemplateId:              ccipMessageSentIdentifier,
											IncludeCreatedEventBlob: true,
										},
									},
								},
							},
						},
					},
					Verbose: true,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}

	var transactions []*ledgerv2.Transaction
	for {
		update, err := updates.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to get updates: %w", err)
		}
		transactions = append(transactions, update.GetTransaction())
	}

	return extractEvents(c.lggr, transactions, c.config.CCIPOwnerParty, ccipMessageSentIdentifier), nil
}

func extractEvents(lggr logger.Logger, transactions []*ledgerv2.Transaction, ccipOwnerParty string, ccipMessageSentTemplateID *ledgerv2.Identifier) []protocol.MessageSentEvent {
	if lggr == nil {
		lggr = logger.Nop()
	}

	var events []protocol.MessageSentEvent
	for _, tx := range transactions {
		if tx == nil {
			continue
		}

		for _, event := range tx.GetEvents() {
			created := event.GetCreated()
			if created == nil {
				continue
			}
			// check that ccipOwnerParty is in the signatories of the created event.
			// This is a defense in depth measure to ensure that we only process events that are signed by the ccipOwnerParty.
			if !slices.Contains(created.GetSignatories(), ccipOwnerParty) {
				continue
			}
			messageSentEvent, err := processCreatedEvent(tx, created, ccipOwnerParty, ccipMessageSentTemplateID)
			if err != nil {
				if errors.Is(err, errMetadataMismatch) {
					continue
				}

				lggr.Errorw(
					"skipping CCIPMessageSent created event after post-filter processing failed",
					"err", err,
					"ledgerOffset", tx.GetOffset(),
					"updateId", tx.GetUpdateId(),
				)

				continue
			}
			if messageSentEvent != nil {
				events = append(events, *messageSentEvent)
			}
		}
	}

	return events
}

// This is a sentinel error that is returned in the event of a metadata mismatch
// in the created event.
var errMetadataMismatch = errors.New("metadata mismatch")

func processCreatedEvent(
	tx *ledgerv2.Transaction,
	created *ledgerv2.CreatedEvent,
	expectedCCIPOwnerParty string,
	ccipMessageSentTemplateID *ledgerv2.Identifier,
) (*protocol.MessageSentEvent, error) {
	if !identifiersClose(created.GetTemplateId(), ccipMessageSentTemplateID) {
		return nil, errMetadataMismatch
	}

	parsed, err := bindings.UnmarshalCreatedEvent[common.CCIPMessageSent](created)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CCIPMessageSent created event: %w", err)
	}

	if string(parsed.CcipOwner) != expectedCCIPOwnerParty {
		return nil, errMetadataMismatch
	}

	messageSentEvent, err := ccipMessageSentEventToProtocol(&parsed.Event)
	if err != nil {
		return nil, fmt.Errorf("failed to process CCIPMessageSent event: %w", err)
	}

	messageSentEvent.BlockNumber = uint64(tx.GetOffset()) //nolint:gosec // offset is always non-negative
	txHash, err := protocol.NewByteSliceFromHex(tx.GetUpdateId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse tx hash from update ID %s: %w", tx.GetUpdateId(), err)
	}
	messageSentEvent.TxHash = txHash

	return messageSentEvent, nil
}

// ccipMessageSentEventToProtocol converts the binding type common.CCIPMessageSentEvent
// to protocol.MessageSentEvent (hex decoding, message decode, receipt mapping, validations).
func ccipMessageSentEventToProtocol(evt *common.CCIPMessageSentEvent) (*protocol.MessageSentEvent, error) {
	messageSentEvent := &protocol.MessageSentEvent{}

	messageID, err := hex.DecodeString(string(evt.MessageId))
	if err != nil {
		return nil, fmt.Errorf("failed to decode message ID: %w, input: %s", err, evt.MessageId)
	}
	copy(messageSentEvent.MessageID[:], messageID)

	encodedMessage, err := hex.DecodeString(string(evt.EncodedMessage))
	if err != nil {
		return nil, fmt.Errorf("failed to decode encoded message: %w, input: %s", err, evt.EncodedMessage)
	}
	msg, err := protocol.DecodeMessage(encodedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w, input: %s", err, evt.EncodedMessage)
	}
	messageSentEvent.Message = *msg

	var verifierBlobs [][]byte
	for _, blobHex := range evt.VerifierBlobs {
		verifierBlobBytes, err := hex.DecodeString(string(blobHex))
		if err != nil {
			return nil, fmt.Errorf("failed to decode verifier blob: %w, input: %s", err, blobHex)
		}
		verifierBlobs = append(verifierBlobs, verifierBlobBytes)
	}

	messageSentEvent.Receipts, err = receiptsBindingToProtocol(evt.Receipts)
	if err != nil {
		return nil, fmt.Errorf("failed to process receipts: %w", err)
	}

	// There are more receipts than verifierBlobs.
	// https://github.com/smartcontractkit/chainlink-ccip/blob/f47d23c550cefae31f13ee7368b747018c5035f4/chains/evm/contracts/onRamp/OnRamp.sol#L129-L132
	if len(messageSentEvent.Receipts) < len(verifierBlobs) {
		return nil, fmt.Errorf(
			"expected more receipts than verifier blobs, got %d receipts and %d verifier blobs",
			len(messageSentEvent.Receipts), len(verifierBlobs),
		)
	}

	for i, blob := range verifierBlobs {
		messageSentEvent.Receipts[i].Blob = blob
	}

	if messageSentEvent.Message.MustMessageID() != messageSentEvent.MessageID {
		return nil, fmt.Errorf("message ID mismatch, from event: %s, from message: %s", messageSentEvent.MessageID.String(), messageSentEvent.Message.MustMessageID().String())
	}

	if err := protocol.ValidateCCVAndExecutorHash(messageSentEvent.Message, messageSentEvent.Receipts); err != nil {
		return nil, fmt.Errorf("ccvAndExecutorHash validation failed: %w", err)
	}

	return messageSentEvent, nil
}

// receiptsBindingToProtocol converts binding []common.Receipt to []protocol.ReceiptWithBlob.
func receiptsBindingToProtocol(receipts []common.Receipt) ([]protocol.ReceiptWithBlob, error) {
	protoReceipts := make([]protocol.ReceiptWithBlob, 0, len(receipts))
	for _, r := range receipts {
		decoded, err := protocol.NewUnknownAddressFromHex(string(r.IssuerAddress))
		if err != nil {
			return nil, fmt.Errorf("failed to decode issuerAddress: %w, input: %s", err, r.IssuerAddress)
		}
		feeTokenAmountFloat, ok := new(big.Float).SetString(string(r.FeeTokenAmount))
		if !ok {
			return nil, fmt.Errorf("failed to parse fee token amount numeric, input: %s", r.FeeTokenAmount)
		}
		feeTokenAmount, _ := feeTokenAmountFloat.Int(nil)
		extraArgs, err := hex.DecodeString(string(r.ExtraArgs))
		if err != nil {
			return nil, fmt.Errorf("failed to decode extra args: %w, input: %s", err, r.ExtraArgs)
		}
		protoReceipts = append(protoReceipts, protocol.ReceiptWithBlob{
			Issuer:            decoded,
			DestGasLimit:      uint64(r.DestGasLimit),      //nolint:gosec // int64 is always non-negative
			DestBytesOverhead: uint32(r.DestBytesOverhead), //nolint:gosec // int64 is always non-negative
			FeeTokenAmount:    feeTokenAmount,
			ExtraArgs:         protocol.ByteSlice(extraArgs),
		})
	}

	return protoReceipts, nil
}

func identifiersClose(a, b *ledgerv2.Identifier) bool {
	// handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Note: we don't check package ID because what is returned from the server is not the package name
	// but the package hex ID. Since this hex ID may change frequently, we don't want to rely on it.
	return a.GetModuleName() == b.GetModuleName() && a.GetEntityName() == b.GetEntityName()
}

// GetBlocksHeaders implements chainaccess.SourceReader.
// The blockNumbers passed in are offset numbers, since that's all we ever return from LatestAndFinalizedBlock.
func (c *sourceReader) GetBlocksHeaders(ctx context.Context, blockNumbers []*big.Int) (map[uint64]protocol.BlockHeader, error) {
	latest, _, err := c.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	if latest == nil {
		return nil, fmt.Errorf("latest block is nil")
	}

	headers := make(map[uint64]protocol.BlockHeader)
	for _, blockNum := range blockNumbers {
		if blockNum.Uint64() > latest.Number {
			return nil, fmt.Errorf("block number is greater than latest offset: %d > %d", blockNum.Uint64(), latest.Number)
		}

		h := intToBytes32(blockNum.Uint64())
		headers[blockNum.Uint64()] = protocol.BlockHeader{
			Number:     blockNum.Uint64(),
			Hash:       h,
			ParentHash: parentHash(blockNum.Uint64()),
			// TODO: determine if we can get an offset's timestamp.
			// Timestamp: time.Time{},
		}
	}

	return headers, nil
}

// GetRMNCursedSubjects implements chainaccess.SourceReader.
func (c *sourceReader) GetRMNCursedSubjects(ctx context.Context) ([]protocol.Bytes16, error) {
	latest, _, err := c.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	if latest == nil {
		return nil, fmt.Errorf("latest block is nil")
	}

	activeContract, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		c.stateServiceClient,
		c.config.NodeOperatorParty,
		c.config.RMNRemoteTemplateID.String(),
		c.rmnRemoteInstanceAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find active contract: %w", err)
	}

	createdEvent := activeContract.GetCreatedEvent()
	if createdEvent == nil {
		return nil, fmt.Errorf("created event is nil")
	}

	rmnRemote, err := bindings.UnmarshalCreatedEvent[rmn.RMNRemote](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RMNRemote created event: %w", err)
	}

	var cursedSubjects []protocol.Bytes16
	for _, subject := range rmnRemote.CursedSubjects {
		subjectStr := string(subject)
		if !strings.HasPrefix(subjectStr, "0x") {
			subjectStr = "0x" + subjectStr
		}
		subjectBytes16, err := protocol.NewBytes16FromString(subjectStr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode subject from BytesHex: %w, input: %s", err, subject)
		}
		cursedSubjects = append(cursedSubjects, subjectBytes16)
	}

	return cursedSubjects, nil
}

// LatestAndFinalizedBlock returns the latest offset of the canton validator we are connected to.
// The latest "block" on Canton is always finalized.
func (c *sourceReader) LatestAndFinalizedBlock(ctx context.Context) (latest, finalized *protocol.BlockHeader, err error) {
	end, err := c.stateServiceClient.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ledger end: %w", err)
	}
	offsetUint64 := uint64(end.GetOffset()) //nolint:gosec // offset is always non-negative
	h := intToBytes32(offsetUint64)
	parentHash := parentHash(offsetUint64)

	return &protocol.BlockHeader{
			Number:     offsetUint64,
			Hash:       h,
			ParentHash: parentHash,
			// TODO: determine if we can get an offset's timestamp.
			// Timestamp: time.Time{},
		}, &protocol.BlockHeader{
			Number:     offsetUint64,
			Hash:       h,
			ParentHash: parentHash,
			// TODO: determine if we can get an offset's timestamp.
			// Timestamp: time.Time{},
		}, nil
}

func (c *sourceReader) LatestSafeBlock(ctx context.Context) (safe *protocol.BlockHeader, err error) {
	//nolint:nilnil // Returns nil without an error when the chain does not support the safe tag.
	return nil, nil
}

func intToBytes32(i uint64) protocol.Bytes32 {
	var b protocol.Bytes32
	binary.BigEndian.PutUint64(b[:], i)

	return b
}

func parentHash(i uint64) protocol.Bytes32 {
	// subtract 1 from i in a checked manner
	// i.e. making sure we don't underflow
	if i == 0 {
		return protocol.Bytes32{}
	}

	return intToBytes32(i - 1)
}

var _ chainaccess.SourceReader = (*sourceReader)(nil)
