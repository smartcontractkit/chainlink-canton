package heads

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink-framework/chains"
)

const TrackableCallbackTimeout = 2 * time.Second

type callbackSet[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] map[int]Trackable[H, BLOCK_HASH]

func (set callbackSet[H, BLOCK_HASH]) values() []Trackable[H, BLOCK_HASH] {
	values := make([]Trackable[H, BLOCK_HASH], 0, len(set))
	for _, callback := range set {
		values = append(values, callback)
	}
	return values
}

// Trackable is implemented by the core txm to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
type Trackable[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] interface {
	// OnNewLongestChain sends a new head when it becomes available. Subscribers can recursively trace the parent
	// of the head to the finalized block back.
	OnNewLongestChain(ctx context.Context, head H)
}

// Broadcaster relays new Heads to all subscribers.
type Broadcaster[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] interface {
	services.Service
	BroadcastNewLongestChain(H)
	Subscribe(callback Trackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func())
}

type broadcaster[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable] struct {
	services.Service
	eng *services.Engine

	callbacks      callbackSet[H, BLOCK_HASH]
	mailbox        *mailbox.Mailbox[H]
	mutex          sync.Mutex
	latest         H
	lastCallbackID int
}

// NewBroadcaster creates a new Broadcaster
func NewBroadcaster[
	H chains.Head[BLOCK_HASH],
	BLOCK_HASH chains.Hashable,
](
	lggr logger.Logger,
) Broadcaster[H, BLOCK_HASH] {
	hb := &broadcaster[H, BLOCK_HASH]{
		callbacks: make(callbackSet[H, BLOCK_HASH]),
		mailbox:   mailbox.NewSingle[H](),
	}
	hb.Service, hb.eng = services.Config{
		Name:  "HeadBroadcaster",
		Start: hb.start,
		Close: hb.close,
	}.NewServiceEngine(lggr)
	return hb
}

func (b *broadcaster[H, BLOCK_HASH]) start(context.Context) error {
	b.eng.Go(b.run)
	return nil
}

func (b *broadcaster[H, BLOCK_HASH]) close() error {
	b.mutex.Lock()
	// clear all callbacks
	b.callbacks = make(callbackSet[H, BLOCK_HASH])
	b.mutex.Unlock()
	return nil
}

func (b *broadcaster[H, BLOCK_HASH]) BroadcastNewLongestChain(head H) {
	b.mailbox.Deliver(head)
}

// Subscribe subscribes to OnNewLongestChain and Connect until Broadcaster is closed,
// or unsubscribe callback is called explicitly
func (b *broadcaster[H, BLOCK_HASH]) Subscribe(callback Trackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func()) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	currentLongestChain = b.latest

	b.lastCallbackID++
	callbackID := b.lastCallbackID
	b.callbacks[callbackID] = callback
	unsubscribe = func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		delete(b.callbacks, callbackID)
	}

	return
}

func (b *broadcaster[H, BLOCK_HASH]) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-b.mailbox.Notify():
			b.executeCallbacks(ctx)
		}
	}
}

// DEV: the head relayer makes no promises about head delivery! Subscribing
// Jobs should expect to the relayer to skip heads if there is a large number of listeners
// and all callbacks cannot be completed in the allotted time.
func (b *broadcaster[H, BLOCK_HASH]) executeCallbacks(ctx context.Context) {
	head, exists := b.mailbox.Retrieve()
	if !exists {
		b.eng.Info("No head to retrieve. It might have been skipped")
		return
	}

	b.mutex.Lock()
	callbacks := b.callbacks.values()
	b.latest = head
	b.mutex.Unlock()

	b.eng.Debugw("Initiating callbacks",
		"headNum", head.BlockNumber(),
		"numCallbacks", len(callbacks),
	)

	wg := sync.WaitGroup{}
	wg.Add(len(callbacks))

	for _, callback := range callbacks {
		go func(trackable Trackable[H, BLOCK_HASH]) {
			defer wg.Done()
			start := time.Now()
			cctx, cancel := context.WithTimeout(ctx, TrackableCallbackTimeout)
			defer cancel()
			trackable.OnNewLongestChain(cctx, head)
			elapsed := time.Since(start)
			b.eng.Debugw(fmt.Sprintf("Finished callback in %s", elapsed),
				"callbackType", reflect.TypeOf(trackable), "blockNumber", head.BlockNumber(), "time", elapsed)
		}(callback)
	}

	wg.Wait()
}
