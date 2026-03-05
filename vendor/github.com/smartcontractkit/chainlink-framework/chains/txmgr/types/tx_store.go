package types

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-framework/chains"
	"github.com/smartcontractkit/chainlink-framework/chains/fees"
)

// TxStore is a superset of all the needed persistence layer methods
type TxStore[ADDR chains.Hashable, CID chains.ID, THASH chains.Hashable, BHASH chains.Hashable, R ChainReceipt[THASH, BHASH], SEQ chains.Sequence, FEE fees.Fee] interface {
	UnstartedTxQueuePruner
	TxHistoryReaper[CID]
	TransactionStore[ADDR, CID, THASH, BHASH, SEQ, FEE]

	// Find confirmed txes beyond the minConfirmations param that require callback but have not yet been signaled
	FindTxesPendingCallback(ctx context.Context, latest, finalized int64, chainID CID) (receiptsPlus []ReceiptPlus[R], err error)
	// Update tx to mark that its callback has been signaled
	UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainID CID) error
	SaveFetchedReceipts(ctx context.Context, r []R) error

	// additional methods for tx store management
	CheckTxQueueCapacity(ctx context.Context, fromAddress ADDR, maxQueuedTransactions uint64, chainID CID) (err error)
	Close()
	Abandon(ctx context.Context, id CID, addr ADDR) error
	// Find transactions by a field in the TxMeta blob and transaction states
	FindTxesByMetaFieldAndStates(ctx context.Context, metaField string, metaValue string, states []TxState, chainID *big.Int) (tx []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided by transaction states
	FindTxesWithMetaFieldByStates(ctx context.Context, metaField string, states []TxState, chainID *big.Int) (tx []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	// Find transactions with a non-null TxMeta field that was provided and a receipt block number greater than or equal to the one provided
	FindTxesWithMetaFieldByReceiptBlockNum(ctx context.Context, metaField string, blockNum int64, chainID *big.Int) (tx []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	// Find transactions loaded with transaction attempts and receipts by transaction IDs and states
	FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx context.Context, ids []int64, states []TxState, chainID *big.Int) (tx []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CID) (tx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindReceiptWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CID) (ChainReceipt[THASH, BHASH], error)
}

// TransactionStore contains the persistence layer methods needed to manage Txs and TxAttempts
type TransactionStore[ADDR chains.Hashable, CID chains.ID, THASH chains.Hashable, BHASH chains.Hashable, SEQ chains.Sequence, FEE fees.Fee] interface {
	CountUnconfirmedTransactions(ctx context.Context, fromAddress ADDR, chainID CID) (count uint32, err error)
	CountTransactionsByState(ctx context.Context, state TxState, chainID CID) (count uint32, err error)
	CountUnstartedTransactions(ctx context.Context, fromAddress ADDR, chainID CID) (count uint32, err error)
	CreateTransaction(ctx context.Context, txRequest TxRequest[ADDR, THASH], chainID CID) (tx Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	DeleteInProgressAttempt(ctx context.Context, attempt TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	FindLatestSequence(ctx context.Context, fromAddress ADDR, chainID CID) (SEQ, error)
	// FindReorgOrIncludedTxs returns either a list of re-org'd transactions or included transactions based on the provided sequence
	FindReorgOrIncludedTxs(ctx context.Context, fromAddress ADDR, nonce SEQ, chainID CID) (reorgTx []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], includedTxs []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindTxsRequiringGasBump(ctx context.Context, address ADDR, blockNum, gasBumpThreshold, depth int64, chainID CID) (etxs []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindTxsRequiringResubmissionDueToInsufficientFunds(ctx context.Context, address ADDR, chainID CID) (etxs []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindTxAttemptsConfirmedMissingReceipt(ctx context.Context, chainID CID) (attempts []TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindTxAttemptsRequiringResend(ctx context.Context, olderThan time.Time, maxInFlightTransactions uint32, chainID CID, address ADDR) (attempts []TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	// Search for Tx using the idempotencyKey and chainID
	FindTxWithIdempotencyKey(ctx context.Context, idempotencyKey string, chainID CID) (tx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	// Search for Tx using the fromAddress and sequence
	FindTxWithSequence(ctx context.Context, fromAddress ADDR, seq SEQ) (etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	FindNextUnstartedTransactionFromAddress(ctx context.Context, fromAddress ADDR, chainID CID) (*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], error)

	FindEarliestUnconfirmedBroadcastTime(ctx context.Context, chainID CID) (null.Time, error)
	FindEarliestUnconfirmedTxAttemptBlock(ctx context.Context, chainID CID) (null.Int, error)
	GetTxInProgress(ctx context.Context, fromAddress ADDR) (etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	GetInProgressTxAttempts(ctx context.Context, address ADDR, chainID CID) (attempts []TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	GetAbandonedTransactionsByBatch(ctx context.Context, chainID CID, enabledAddrs []ADDR, offset, limit uint) (txs []*Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	GetTxByID(ctx context.Context, id int64) (tx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], err error)
	HasInProgressTransaction(ctx context.Context, account ADDR, chainID CID) (exists bool, err error)
	LoadTxAttempts(ctx context.Context, etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	PreloadTxes(ctx context.Context, attempts []TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	SaveConfirmedAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], broadcastAt time.Time) error
	SaveInProgressAttempt(ctx context.Context, attempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	SaveInsufficientFundsAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], broadcastAt time.Time) error
	SaveReplacementInProgressAttempt(ctx context.Context, oldAttempt TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], replacementAttempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	SaveSentAttempt(ctx context.Context, timeout time.Duration, attempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], broadcastAt time.Time) error
	SetBroadcastBeforeBlockNum(ctx context.Context, blockNum int64, chainID CID) error
	UpdateBroadcastAts(ctx context.Context, now time.Time, etxIDs []int64) error
	UpdateTxAttemptInProgressToBroadcast(ctx context.Context, etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], attempt TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE], NewAttemptState TxAttemptState) error
	// UpdateTxCallbackCompleted updates tx to mark that its callback has been signaled
	UpdateTxCallbackCompleted(ctx context.Context, pipelineTaskRunRid uuid.UUID, chainID CID) error
	// UpdateTxConfirmed updates transaction states to confirmed
	UpdateTxConfirmed(ctx context.Context, etxIDs []int64) error
	// UpdateTxFatalErrorAndDeleteAttempts updates transaction states to fatal error, deletes attempts, and clears broadcast info and sequence
	UpdateTxFatalErrorAndDeleteAttempts(ctx context.Context, etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
	// UpdateTxFatalError updates transaction states to fatal error with error message
	UpdateTxFatalError(ctx context.Context, etxIDs []int64, errMsg string) error
	UpdateTxsForRebroadcast(ctx context.Context, etxIDs []int64, attemptIDs []int64) error
	UpdateTxsUnconfirmed(ctx context.Context, etxIDs []int64) error
	UpdateTxUnstartedToInProgress(ctx context.Context, etx *Tx[CID, ADDR, THASH, BHASH, SEQ, FEE], attempt *TxAttempt[CID, ADDR, THASH, BHASH, SEQ, FEE]) error
}

type TxHistoryReaper[CID chains.ID] interface {
	ReapTxHistory(ctx context.Context, timeThreshold time.Time, chainID CID) error
}

type UnstartedTxQueuePruner interface {
	PruneUnstartedTxQueue(ctx context.Context, queueSize uint32, subject uuid.UUID) (ids []int64, err error)
}

// R is the raw unparsed transaction receipt
type ReceiptPlus[R any] struct {
	ID           uuid.UUID `db:"pipeline_run_id"`
	Receipt      R         `db:"receipt"`
	FailOnRevert bool      `db:"fail_on_revert"`
}

type ChainReceipt[THASH, BHASH chains.Hashable] interface {
	GetStatus() uint64
	GetTxHash() THASH
	GetBlockNumber() *big.Int
	IsZero() bool
	IsUnmined() bool
	GetFeeUsed() uint64
	GetTransactionIndex() uint
	GetBlockHash() BHASH
	GetRevertReason() *string
	GetEffectiveGasPrice() *big.Int
	GetL1Fee() *big.Int
}
