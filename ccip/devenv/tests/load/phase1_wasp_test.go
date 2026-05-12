package load

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

// TestCantonLoad_Phase1 runs a minimal WASP profile (RPS=1) with a stub Canton→EVM gun.
// No devenv stack or env file is required. If code review shows overlapping Gun.Call despite
// RPS=1, add a wasp.VirtualUser for strict sequential execution (Canton holdings 1-wide).
func TestCantonLoad_Phase1(t *testing.T) {
	stub := NewCanton2EVMStubGun()

	p := wasp.NewProfile().Add(wasp.NewGenerator(&wasp.Config{
		T:        t,
		LoadType: wasp.RPS,
		GenName:  "canton-phase1-stub",
		Schedule: wasp.Combine(
			wasp.Plain(1, 5*time.Second),
		),
		Gun: stub,
		Labels: map[string]string{
			"go_test_name": "canton-load-phase1-stub",
			"branch":       "test",
			"commit":       "test",
		},
		LokiConfig: nil,
	}))

	_, err := p.Run(true)
	require.NoError(t, err)
	p.Wait()

	require.Greater(t, stub.CallCount(), int64(0), "stub should have been exercised")
	require.LessOrEqual(t, stub.MaxConcurrentObserved(), int32(1),
		"Gun.Call must not overlap (Canton holdings 1-wide)")
}
