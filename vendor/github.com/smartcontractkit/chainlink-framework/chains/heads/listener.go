package heads

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jpillora/backoff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

var (
	promNumHeadsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	}, []string{"ChainID"})
	promEthConnectionErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_connection_errors",
		Help: "The total number of node connection errors",
	}, []string{"ChainID"})
)

// Handler is a callback that handles incoming heads
type Handler[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] func(ctx context.Context, header H) error

// Listener is a chain agnostic interface that manages connection of Client that receives heads from the blockchain node
type Listener[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] interface {
	services.Service

	// ListenForNewHeads runs the listen loop (not thread safe)
	ListenForNewHeads(ctx context.Context)

	// ReceivingHeads returns true if the listener is receiving heads (thread safe)
	ReceivingHeads() bool

	// Connected returns true if the listener is connected (thread safe)
	Connected() bool

	// HealthReport returns report of errors within Listener
	HealthReport() map[string]error
}

type ListenerConfig interface {
	BlockEmissionIdleWarningThreshold() time.Duration
}

type listener[
	HTH Head[BLOCK_HASH, ID],
	S chains.Subscription,
	ID chains.ID,
	BLOCK_HASH chains.Hashable,
] struct {
	services.Service
	eng *services.Engine

	config           ListenerConfig
	client           Client[HTH, S, ID, BLOCK_HASH]
	onSubscription   func(context.Context)
	handleNewHead    Handler[HTH, BLOCK_HASH]
	chHeaders        <-chan HTH
	headSubscription chains.Subscription
	connected        atomic.Bool
	receivingHeads   atomic.Bool
}

func NewListener[
	HTH Head[BLOCK_HASH, ID],
	S chains.Subscription,
	ID chains.ID,
	BLOCK_HASH chains.Hashable,
	CLIENT Client[HTH, S, ID, BLOCK_HASH],
](
	lggr logger.Logger,
	client CLIENT,
	config ListenerConfig,
	onSubscription func(context.Context),
	handleNewHead Handler[HTH, BLOCK_HASH],
) Listener[HTH, BLOCK_HASH] {
	hl := &listener[HTH, S, ID, BLOCK_HASH]{
		config:         config,
		client:         client,
		onSubscription: onSubscription,
		handleNewHead:  handleNewHead,
	}
	hl.Service, hl.eng = services.Config{
		Name:  "HeadListener",
		Start: hl.start,
	}.NewServiceEngine(lggr)
	return hl
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) start(context.Context) error {
	l.eng.Go(l.ListenForNewHeads)
	return nil
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) ListenForNewHeads(ctx context.Context) {
	defer l.unsubscribe()

	for {
		if !l.subscribe(ctx) {
			break
		}

		if l.onSubscription != nil {
			l.onSubscription(ctx)
		}
		err := l.receiveHeaders(ctx, l.handleNewHead)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			l.eng.Errorw("Error in new head subscription, unsubscribed", "err", err)
			continue
		}
		break
	}
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) ReceivingHeads() bool {
	return l.receivingHeads.Load()
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) Connected() bool {
	return l.connected.Load()
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	receivingHeads := l.ReceivingHeads()
	connected := l.Connected()
	var err error
	if !receivingHeads || !connected {
		err = fmt.Errorf("Listener connected = %t, receiving heads = %t", connected, receivingHeads)
	}

	return map[string]error{l.Name(): err}
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) receiveHeaders(ctx context.Context, handleNewHead Handler[HTH, BLOCK_HASH]) error {
	var noHeadsAlarmC <-chan time.Time
	var noHeadsAlarmT *time.Ticker
	noHeadsAlarmDuration := l.config.BlockEmissionIdleWarningThreshold()
	if noHeadsAlarmDuration > 0 {
		noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
		noHeadsAlarmC = noHeadsAlarmT.C
		defer noHeadsAlarmT.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case blockHeader, open := <-l.chHeaders:
			chainID := l.client.ConfiguredChainID()
			if noHeadsAlarmT != nil {
				// We've received a head, reset the no heads alarm
				noHeadsAlarmT.Stop()
				noHeadsAlarmT = time.NewTicker(noHeadsAlarmDuration)
				noHeadsAlarmC = noHeadsAlarmT.C
			}
			l.receivingHeads.Store(true)
			if !open {
				return errors.New("head listener: chHeaders prematurely closed")
			}
			if !blockHeader.IsValid() {
				l.eng.Error("got nil block header")
				continue
			}

			// Compare the chain ID of the block header to the chain ID of the client
			if !blockHeader.HasChainID() || blockHeader.ChainID().String() != chainID.String() {
				l.eng.Panicf("head listener for %s received block header for %s", chainID, blockHeader.ChainID())
			}
			promNumHeadsReceived.WithLabelValues(chainID.String()).Inc()

			err := handleNewHead(ctx, blockHeader)
			if ctx.Err() != nil {
				return nil
			} else if err != nil {
				return err
			}

		case err, open := <-l.headSubscription.Err():
			// err can be nil, because of using chainIDSubForwarder
			if !open || err == nil {
				return errors.New("head listener: subscription Err channel prematurely closed")
			}
			return err

		case <-noHeadsAlarmC:
			// We haven't received a head on the channel for a long time, log a warning
			l.eng.Warnf("have not received a head for %v", noHeadsAlarmDuration)
			l.receivingHeads.Store(false)
		}
	}
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) subscribe(ctx context.Context) bool {
	subscribeRetryBackoff := backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}

	chainID := l.client.ConfiguredChainID()

	for {
		l.unsubscribe()

		l.eng.Debugf("Subscribing to new heads on chain %s", chainID.String())

		select {
		case <-ctx.Done():
			return false

		case <-time.After(subscribeRetryBackoff.Duration()):
			err := l.subscribeToHead(ctx)
			if err != nil {
				promEthConnectionErrors.WithLabelValues(chainID.String()).Inc()
				l.eng.Warnw("Failed to subscribe to heads on chain", "chainID", chainID.String(), "err", err)
			} else {
				l.eng.Debugf("Subscribed to heads on chain %s", chainID.String())
				return true
			}
		}
	}
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) subscribeToHead(ctx context.Context) error {
	var err error
	l.chHeaders, l.headSubscription, err = l.client.SubscribeToHeads(ctx)
	if err != nil {
		return fmt.Errorf("Client#SubscribeToHeads: %w", err)
	}

	l.connected.Store(true)

	return nil
}

func (l *listener[HTH, S, ID, BLOCK_HASH]) unsubscribe() {
	if l.headSubscription != nil {
		l.connected.Store(false)
		l.headSubscription.Unsubscribe()
		l.headSubscription = nil
	}
}
