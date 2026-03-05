package heads

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

var (
	promCurrentHead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	}, []string{"evmChainID"})

	promOldHead = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_very_old_head",
		Help: "Counter is incremented every time we get a head that is much lower than the highest seen head ('much lower' is defined as a block that is EVM.FinalityDepth or greater below the highest seen head)",
	}, []string{"evmChainID"})
)

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

// Tracker holds and stores the block experienced by a particular node in a thread safe manner.
type Tracker[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] interface {
	services.Service
	// Backfill given a head will fill in any missing heads up to latestFinalized
	Backfill(ctx context.Context, headWithChain H, prevHeadWithChain H) (err error)
	LatestChain() H
	// LatestSafeBlock returns the latest block that is considered safe to use.
	LatestSafeBlock(ctx context.Context) (safe H, err error)
	// LatestAndFinalizedBlock - returns latest and latest finalized blocks.
	// NOTE: Returns latest finalized block as is, ignoring the FinalityTagBypass feature flag.
	LatestAndFinalizedBlock(ctx context.Context) (latest, finalized H, err error)
}

type ChainConfig interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	FinalityDepth() uint32
	SafeDepth() uint32
	FinalityTagEnabled() bool
	SafeTagSupported() bool
	FinalizedBlockOffset() uint32
}

type TrackerConfig interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
	FinalityTagBypass() bool
	MaxAllowedFinalityDepth() uint32
	PersistenceEnabled() bool
	PersistenceBatchSize() int64
}

type headPair[HTH any] struct {
	head     HTH
	prevHead HTH
}

type tracker[HTH Head[BHASH, ID], S chains.Subscription, ID chains.ID, BHASH chains.Hashable] struct {
	services.Service
	eng *services.Engine

	log             logger.SugaredLogger
	headBroadcaster Broadcaster[HTH, BHASH]
	headSaver       Saver[HTH, BHASH]
	mailMon         *mailbox.Monitor
	client          Client[HTH, S, ID, BHASH]
	chainID         chains.ID
	config          ChainConfig
	htConfig        TrackerConfig

	backfillMB   *mailbox.Mailbox[headPair[HTH]]
	broadcastMB  *mailbox.Mailbox[HTH]
	headListener Listener[HTH, BHASH]
	getNilHead   func() HTH
}

// NewTracker instantiates a new Tracker using Saver to persist new block numbers.
func NewTracker[HTH Head[BHASH, ID], S chains.Subscription, ID chains.ID, BHASH chains.Hashable](
	lggr logger.Logger,
	client Client[HTH, S, ID, BHASH],
	config ChainConfig,
	htConfig TrackerConfig,
	headBroadcaster Broadcaster[HTH, BHASH],
	headSaver Saver[HTH, BHASH],
	mailMon *mailbox.Monitor,
	getNilHead func() HTH,
) Tracker[HTH, BHASH] {
	ht := &tracker[HTH, S, ID, BHASH]{
		headBroadcaster: headBroadcaster,
		client:          client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		backfillMB:      mailbox.NewSingle[headPair[HTH]](),
		broadcastMB:     mailbox.New[HTH](HeadsBufferSize),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
	ht.Service, ht.eng = services.Config{
		Name: "HeadTracker",
		NewSubServices: func(lggr logger.Logger) []services.Service {
			ht.headListener = NewListener[HTH, S, ID, BHASH](lggr, client, config,
				// NOTE: Always try to start the head tracker off with whatever the
				// latest head is, without waiting for the subscription to send us one.
				//
				// In some cases the subscription will send us the most recent head
				// anyway when we connect (but we should not rely on this because it is
				// not specced). If it happens this is fine, and the head will be
				// ignored as a duplicate.
				func(ctx context.Context) {
					err := ht.handleInitialHead(ctx)
					if err != nil {
						ht.log.Errorw("Error handling initial head", "err", err.Error())
					}
				}, ht.handleNewHead)
			return []services.Service{ht.headListener}
		},
		Start: ht.start,
		Close: ht.close,
	}.NewServiceEngine(lggr)
	ht.log = logger.Sugared(ht.eng)
	return ht
}

// Start starts Tracker service.
func (t *tracker[HTH, S, ID, BHASH]) start(context.Context) error {
	t.eng.Go(t.backfillLoop)
	t.eng.Go(t.broadcastLoop)

	t.mailMon.Monitor(t.broadcastMB, "HeadTracker", "Broadcast", t.chainID.String())

	return nil
}

func (t *tracker[HTH, S, ID, BHASH]) handleInitialHead(ctx context.Context) error {
	initialHead, err := t.client.HeadByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch initial head: %w", err)
	}

	if !initialHead.IsValid() {
		t.log.Warnw("Got nil initial head", "head", initialHead)
		return nil
	}
	t.log.Debugw("Got initial head", "head", initialHead, "blockNumber", initialHead.BlockNumber(), "blockHash", initialHead.BlockHash())

	latestFinalized, err := t.calculateLatestFinalized(ctx, initialHead, t.htConfig.FinalityTagBypass())
	if err != nil {
		return fmt.Errorf("failed to calculate latest finalized head: %w", err)
	}

	if !latestFinalized.IsValid() {
		return fmt.Errorf("latest finalized block is not valid")
	}

	latestChain, err := t.headSaver.Load(ctx, latestFinalized.BlockNumber())
	if err != nil {
		return fmt.Errorf("failed to initialized headSaver: %w", err)
	}

	if latestChain.IsValid() {
		earliest := latestChain.EarliestHeadInChain()
		t.log.Debugw(
			"Loaded chain from DB",
			"latest_blockNumber", latestChain.BlockNumber(),
			"latest_blockHash", latestChain.BlockHash(),
			"earliest_blockNumber", earliest.BlockNumber(),
			"earliest_blockHash", earliest.BlockHash(),
		)
	}
	if err := t.handleNewHead(ctx, initialHead); err != nil {
		return fmt.Errorf("error handling initial head: %w", err)
	}

	return nil
}

func (t *tracker[HTH, S, ID, BHASH]) close() error {
	return t.broadcastMB.Close()
}

// verifyFinalizedBlockHashes returns finality violated error if a block hash mismatch is found in provided chains
func (t *tracker[HTH, S, ID, BHASH]) verifyFinalizedBlockHashes(finalizedHeadWithChain chains.Head[BHASH], prevHeadWithChain chains.Head[BHASH]) error {
	if finalizedHeadWithChain == nil || prevHeadWithChain == nil {
		return nil
	}

	prevLatestFinalized := prevHeadWithChain.LatestFinalizedHead()
	if prevLatestFinalized == nil {
		return nil
	}

	prevLatestFinalizedBlockNum := prevLatestFinalized.BlockNumber()
	prevLatestFinalizedHash := prevLatestFinalized.BlockHash()
	finalizedHead, err := finalizedHeadWithChain.HeadAtHeight(prevLatestFinalizedBlockNum)
	if err != nil {
		return nil
	}

	finalizedBlockNum := finalizedHead.BlockNumber()
	if finalizedBlockNum < prevLatestFinalizedBlockNum {
		return fmt.Errorf("latest finalized head at height %d is behind previously seen finalized head at height %d: %w",
			finalizedBlockNum, prevLatestFinalizedBlockNum, types.ErrFinalityViolated)
	}

	finalizedHash := finalizedHead.BlockHash()
	if finalizedHash != prevLatestFinalizedHash {
		return fmt.Errorf("block hash mismatch at height %d: expected %s, got %s: %w",
			prevLatestFinalizedBlockNum, prevLatestFinalizedHash, finalizedHash, types.ErrFinalityViolated)
	}
	return nil
}

func (t *tracker[HTH, S, ID, BHASH]) instantFinality() bool {
	return !t.config.FinalityTagEnabled() && t.config.FinalityDepth() == 0 && t.config.FinalizedBlockOffset() == 0
}

func (t *tracker[HTH, S, ID, BHASH]) Backfill(ctx context.Context, headWithChain HTH, prevHeadWithChain HTH) (err error) {
	latestFinalized, err := t.calculateLatestFinalized(ctx, headWithChain, t.htConfig.FinalityTagBypass())
	if err != nil {
		return fmt.Errorf("failed to calculate finalized block: %w", err)
	}

	if !latestFinalized.IsValid() {
		return errors.New("can not perform backfill without a valid latestFinalized head")
	}

	if headWithChain.BlockNumber() < latestFinalized.BlockNumber() {
		const warnMsg = "expected head of canonical chain to be ahead of the latestFinalized, but this may be normal on chains with fast finality due to fetch timing"
		t.log.With("head_block_num", headWithChain.BlockNumber(),
			"latest_finalized_block_number", latestFinalized.BlockNumber()).
			Warnf(warnMsg)
		return errors.New(warnMsg)
	}

	if headWithChain.BlockNumber()-latestFinalized.BlockNumber() > int64(t.htConfig.MaxAllowedFinalityDepth()) {
		return fmt.Errorf("gap between latest finalized block (%d) and current head (%d) is too large (> %d)",
			latestFinalized.BlockNumber(), headWithChain.BlockNumber(), t.htConfig.MaxAllowedFinalityDepth())
	}

	if !t.instantFinality() {
		finalityViolationCondition := fmt.Sprintf("backfill %v", types.ErrFinalityViolated)
		// verify block hashes since calculateLatestFinalized made an additional RPC call
		err = t.verifyFinalizedBlockHashes(latestFinalized, prevHeadWithChain.LatestFinalizedHead())
		if err != nil {
			t.eng.SetHealthCond(finalityViolationCondition, err)
			return err
		}
		t.eng.ClearHealthCond(finalityViolationCondition)
	}

	return t.backfill(ctx, headWithChain, latestFinalized)
}

func (t *tracker[HTH, S, ID, BHASH]) LatestChain() HTH {
	return t.headSaver.LatestChain()
}

func (t *tracker[HTH, S, ID, BHASH]) handleNewHead(ctx context.Context, head HTH) error {
	prevHead := t.headSaver.LatestChain()

	t.log.Debugw(fmt.Sprintf("Received new head %v", head.BlockNumber()),
		"blockHash", head.BlockHash(),
		"parentHeadHash", head.GetParentHash(),
		"blockTs", head.GetTimestamp(),
		"blockTsUnix", head.GetTimestamp().Unix(),
		"blockDifficulty", head.BlockDifficulty(),
	)

	var prevLatestFinalized chains.Head[BHASH]
	if prevHead.IsValid() {
		prevLatestFinalized = prevHead.LatestFinalizedHead()
	}

	finalityViolationCondition := fmt.Sprintf("handleNewHead %v", types.ErrFinalityViolated)
	if prevLatestFinalized != nil && head.BlockNumber() < prevLatestFinalized.BlockNumber() {
		promOldHead.WithLabelValues(t.chainID.String()).Inc()
		t.log.Critical("Got very old block. Either a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers. This node may not function correctly without manual intervention.", "err", types.ErrFinalityViolated)
		oldBlockErr := fmt.Errorf("got very old block with number %d (highest seen was %d)", head.BlockNumber(), prevHead.BlockNumber())
		err := fmt.Errorf("%w: %w", oldBlockErr, types.ErrFinalityViolated)
		t.eng.SetHealthCond(finalityViolationCondition, err)
		return err
	}

	if err := t.verifyFinalizedBlockHashes(head.LatestFinalizedHead(), prevHead); err != nil {
		t.log.Critical(err)
		t.eng.SetHealthCond(finalityViolationCondition, err)
		return err
	}
	t.eng.ClearHealthCond(finalityViolationCondition)

	if err := t.headSaver.Save(ctx, head); ctx.Err() != nil {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to save head: %#v: %w", head, err)
	}

	if !prevHead.IsValid() || head.BlockNumber() > prevHead.BlockNumber() {
		promCurrentHead.WithLabelValues(t.chainID.String()).Set(float64(head.BlockNumber()))

		headWithChain := t.headSaver.Chain(head.BlockHash())
		if !headWithChain.IsValid() {
			return fmt.Errorf("heads.tracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}
		t.backfillMB.Deliver(headPair[HTH]{headWithChain, prevHead})
		t.broadcastMB.Deliver(headWithChain)
	} else if head.BlockNumber() == prevHead.BlockNumber() {
		if head.BlockHash() != prevHead.BlockHash() {
			t.log.Debugw("Got duplicate head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockHash())
		} else {
			t.log.Debugw("Head already in the database", "head", head.BlockHash())
		}
	} else {
		t.log.Debugw("Got out of order head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockNumber())
		promOldHead.WithLabelValues(t.chainID.String()).Inc()
		if prevLatestFinalized == nil {
			// sanity check
			finalityDepth := int64(t.config.FinalityDepth())
			if head.BlockNumber() < prevHead.BlockNumber()-finalityDepth {
				t.log.Warnf("Received old block at height %d past finality depth of %d. Either a re-org occurred, one of the RPC nodes has gotten out of sync, or the chain went backwards in block numbers.", head.BlockNumber(), finalityDepth)
			}
		}
	}
	return nil
}

func (t *tracker[HTH, S, ID, BHASH]) broadcastLoop(ctx context.Context) {
	samplingInterval := t.htConfig.SamplingInterval()
	if samplingInterval > 0 {
		t.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-debounceHead.C:
				item := t.broadcastMB.RetrieveLatestAndClear()
				if !item.IsValid() {
					continue
				}
				t.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		t.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.broadcastMB.Notify():
				for {
					item, exists := t.broadcastMB.Retrieve()
					if !exists {
						break
					}
					t.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}

func (t *tracker[HTH, S, ID, BHASH]) backfillLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.backfillMB.Notify():
			for {
				backfillHeadPair, exists := t.backfillMB.Retrieve()
				if !exists {
					break
				}
				{
					err := t.Backfill(ctx, backfillHeadPair.head, backfillHeadPair.prevHead)
					if err != nil {
						t.log.Warnw("Unexpected error while backfilling heads", "err", err)
					} else if ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

func (t *tracker[HTH, S, ID, BHASH]) LatestSafeBlock(ctx context.Context) (safe HTH, err error) {
	if t.config.FinalityTagEnabled() {
		var latestSafe HTH
		if t.config.SafeTagSupported() {
			latestSafe, err = t.client.LatestSafeBlock(ctx)
			if err != nil {
				return latestSafe, fmt.Errorf("failed to get latest finalized block: %w", err)
			}
		} else {
			// return latest finalized block if safe tag is not enabled
			_, finalized, err2 := t.LatestAndFinalizedBlock(ctx)
			if err2 != nil {
				return finalized, fmt.Errorf("failed to get latest finalized block: %w", err2)
			}
			latestSafe = finalized
		}

		if !latestSafe.IsValid() {
			return latestSafe, fmt.Errorf("failed to get valid latest finalized block")
		}
		return latestSafe, nil
	}
	latest, err := t.client.HeadByNumber(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to get latest block: %w", err)
		return
	}

	if !latest.IsValid() {
		err = fmt.Errorf("expected latest block to be valid")
		return
	}
	if t.instantFinality() {
		return latest, nil
	}
	safeDepth := int64(t.config.SafeDepth())
	if safeDepth <= 0 {
		safeDepth = int64(t.config.FinalityDepth())
	}
	safeBlockNumber := latest.BlockNumber() - safeDepth
	if safeBlockNumber <= 0 {
		safeBlockNumber = 0
	}
	return t.getHeadAtHeight(ctx, latest.BlockHash(), safeBlockNumber)
}

// LatestAndFinalizedBlock - returns latest and latest finalized blocks.
// NOTE: Returns latest finalized block as is, ignoring the FinalityTagBypass feature flag.
// TODO: BCI-3321 use cached values instead of making RPC requests
func (t *tracker[HTH, S, ID, BHASH]) LatestAndFinalizedBlock(ctx context.Context) (latest, finalized HTH, err error) {
	latest, err = t.client.HeadByNumber(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to get latest block: %w", err)
		return
	}

	if !latest.IsValid() {
		err = fmt.Errorf("expected latest block to be valid")
		return
	}

	finalized, err = t.calculateLatestFinalized(ctx, latest, false)
	if err != nil {
		err = fmt.Errorf("failed to calculate latest finalized block: %w", err)
		return
	}
	if !finalized.IsValid() {
		err = fmt.Errorf("expected finalized block to be valid")
		return
	}

	return
}

func (t *tracker[HTH, S, ID, BHASH]) getHeadAtHeight(ctx context.Context, chainHeadHash BHASH, blockHeight int64) (HTH, error) {
	chainHead := t.headSaver.Chain(chainHeadHash)
	if chainHead.IsValid() {
		// check if provided chain contains a block of specified height
		headAtHeight, err := chainHead.HeadAtHeight(blockHeight)
		if err == nil {
			// we are forced to reload the block due to type mismatched caused by generics
			hthAtHeight := t.headSaver.Chain(headAtHeight.BlockHash())
			// ensure that the block was not removed from the chain by another goroutine
			if hthAtHeight.IsValid() {
				return hthAtHeight, nil
			}
		}
	}

	return t.client.HeadByNumber(ctx, big.NewInt(blockHeight))
}

// calculateLatestFinalized - returns latest finalized block. It's expected that currentHeadNumber - is the head of
// canonical chain. There is no guaranties that returned block belongs to the canonical chain. Additional verification
// must be performed before usage.
func (t *tracker[HTH, S, ID, BHASH]) calculateLatestFinalized(ctx context.Context, currentHead HTH, finalityTagBypass bool) (HTH, error) {
	if t.config.FinalityTagEnabled() && !finalityTagBypass {
		latestFinalized, err := t.client.LatestFinalizedBlock(ctx)
		if err != nil {
			return latestFinalized, fmt.Errorf("failed to get latest finalized block: %w", err)
		}

		if !latestFinalized.IsValid() {
			return latestFinalized, fmt.Errorf("failed to get valid latest finalized block")
		}

		if t.config.FinalizedBlockOffset() == 0 {
			return latestFinalized, nil
		}

		finalizedBlockNumber := max(latestFinalized.BlockNumber()-int64(t.config.FinalizedBlockOffset()), 0)
		return t.getHeadAtHeight(ctx, latestFinalized.BlockHash(), finalizedBlockNumber)
	}
	// no need to make an additional RPC call on chains with instant finality
	if t.instantFinality() {
		return currentHead, nil
	}
	finalizedBlockNumber := currentHead.BlockNumber() - int64(t.config.FinalityDepth()) - int64(t.config.FinalizedBlockOffset())
	if finalizedBlockNumber <= 0 {
		finalizedBlockNumber = 0
	}
	return t.getHeadAtHeight(ctx, currentHead.BlockHash(), finalizedBlockNumber)
}

// backfill fetches all missing heads up until the latestFinalizedHead
func (t *tracker[HTH, S, ID, BHASH]) backfill(ctx context.Context, head, latestFinalizedHead HTH) (err error) {
	headBlockNumber := head.BlockNumber()
	mark := time.Now()
	fetched := 0
	baseHeight := latestFinalizedHead.BlockNumber()
	l := t.log.With("blockNumber", headBlockNumber,
		"n", headBlockNumber-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", headBlockNumber-1)
	l.Debug("Starting backfill")
	defer func() {
		if ctx.Err() != nil {
			l.Warnw("Backfill context error", "err", ctx.Err())
			return
		}
		l.Debugw("Finished backfill",
			"fetched", fetched,
			"time", time.Since(mark),
			"err", err)
	}()

	for i := head.BlockNumber() - 1; i >= baseHeight; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		existingHead := t.headSaver.Chain(head.GetParentHash())
		if existingHead.IsValid() {
			head = existingHead
			continue
		}
		head, err = t.fetchAndSaveHead(ctx, i, head.GetParentHash())
		fetched++
		if ctx.Err() != nil {
			t.log.Debugw("context canceled, aborting backfill", "err", err, "ctx.Err", ctx.Err())
			return fmt.Errorf("fetchAndSaveHead failed: %w", ctx.Err())
		} else if err != nil {
			return fmt.Errorf("fetchAndSaveHead failed: %w", err)
		}
	}

	if head.BlockHash() != latestFinalizedHead.BlockHash() {
		t.log.Criticalw("Finalized block missing from canonical chain",
			"finalized_block_number", latestFinalizedHead.BlockNumber(), "finalized_hash", latestFinalizedHead.BlockHash(),
			"canonical_chain_block_number", head.BlockNumber(), "canonical_chain_hash", head.BlockHash())
		return FinalizedMissingError[BHASH]{latestFinalizedHead.BlockHash(), head.BlockHash()}
	}

	l = l.With("latest_finalized_block_hash", latestFinalizedHead.BlockHash(),
		"latest_finalized_block_number", latestFinalizedHead.BlockNumber())

	err = t.headSaver.MarkFinalized(ctx, latestFinalizedHead)
	if err != nil {
		l.Debugw("failed to mark block as finalized", "err", err)
		return nil
	}

	l.Debugw("marked block as finalized")

	return
}

type FinalizedMissingError[BLOCK_HASH chains.Hashable] struct {
	Finalized, Canonical BLOCK_HASH
}

func (e FinalizedMissingError[BLOCK_HASH]) Error() string {
	return fmt.Sprintf("finalized block %s missing from canonical chain %s", e.Finalized, e.Canonical)
}

func (t *tracker[HTH, S, ID, BHASH]) fetchAndSaveHead(ctx context.Context, n int64, hash BHASH) (HTH, error) {
	t.log.Debugw("Fetching head", "blockHeight", n, "blockHash", hash)
	head, err := t.client.HeadByHash(ctx, hash)
	if err != nil {
		return t.getNilHead(), err
	} else if !head.IsValid() {
		return t.getNilHead(), errors.New("got nil head")
	}
	err = t.headSaver.Save(ctx, head)
	if err != nil {
		return t.getNilHead(), err
	}
	return head, nil
}
