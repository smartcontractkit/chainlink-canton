package adapters

import (
	"context"
	"os"
	"time"
)

type vpnPauseLogger interface {
	Warnw(msg string, keysAndValues ...any)
}

// pauseBeforeCantonLedgerReadIfConfigured sleeps before Canton ledger RPCs so operators
// can switch VPN (chainlink-legacy for JD/EVM → SmartContract for Canton MCMS reads).
// Set CANTON_CONFIGURE_LANES_VPN_PAUSE=30s (any Go duration).
func pauseBeforeCantonLedgerReadIfConfigured(ctx context.Context, logger vpnPauseLogger, reason string) {
	raw := os.Getenv("CANTON_CONFIGURE_LANES_VPN_PAUSE")
	if raw == "" {
		return
	}

	pause, err := time.ParseDuration(raw)
	if err != nil || pause <= 0 {
		pause = 30 * time.Second
	}

	if logger != nil {
		logger.Warnw(
			"switch vpn now",
			"targetVPN", "SmartContract",
			"reason", reason,
			"pause", pause.String(),
		)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-time.After(pause):
	case <-ctx.Done():
	}
}
