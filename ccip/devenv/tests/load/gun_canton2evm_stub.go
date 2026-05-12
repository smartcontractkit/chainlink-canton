package load

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

// Canton2EVMStubGun is a no-op Canton→EVM load gun for WASP wiring tests.
//
// Canton contract holdings are effectively single-threaded (1-wide): only one
// outbound message flow should be in flight at a time. This stub tracks concurrent
// Call depth so tests can prove the generator did not overlap Gun.Call invocations.
type Canton2EVMStubGun struct {
	mu            sync.Mutex
	inFlight      int32
	maxConcurrent int32
	calls         atomic.Int64
}

// NewCanton2EVMStubGun returns a Canton→EVM stub gun that succeeds immediately
// (after a short sleep to widen any scheduling race).
func NewCanton2EVMStubGun() *Canton2EVMStubGun {
	return &Canton2EVMStubGun{}
}

// Call implements wasp.Gun. At most one Call should be active; if a second
// enters while the first is still running, depth exceeds 1 and the test fails
// when *testing.T is set on wasp.Config.
func (g *Canton2EVMStubGun) Call(gen *wasp.Generator) *wasp.Response {
	var depth int32
	g.mu.Lock()
	g.inFlight++
	depth = g.inFlight
	if depth > g.maxConcurrent {
		g.maxConcurrent = depth
	}
	g.mu.Unlock()

	g.calls.Add(1)
	defer func() {
		g.mu.Lock()
		g.inFlight--
		g.mu.Unlock()
	}()

	if gen != nil && gen.Cfg != nil && gen.Cfg.T != nil && depth > 1 {
		require.FailNow(gen.Cfg.T, "overlapping Canton2EVM stub Gun.Call",
			"expected single-flight (Canton holdings 1-wide); concurrent depth=%d", depth)
	}

	// Small delay so overlapping Calls would be observable if the scheduler fired them concurrently.
	time.Sleep(2 * time.Millisecond)

	return &wasp.Response{
		Failed:     false,
		StatusCode: "200",
		Duration:   0,
	}
}

// MaxConcurrentObserved returns the peak number of overlapping Call invocations seen.
func (g *Canton2EVMStubGun) MaxConcurrentObserved() int32 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.maxConcurrent
}

// CallCount returns how many times Call ran.
func (g *Canton2EVMStubGun) CallCount() int64 {
	return g.calls.Load()
}
