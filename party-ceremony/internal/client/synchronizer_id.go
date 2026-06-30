package client

import (
	"context"
	"fmt"
	"strings"
)

// LogicalSynchronizerID converts a physical topology synchronizer ID
// (e.g. global-domain::1220...::34-0) to the logical ID expected by Ledger
// API InteractiveSubmission (global-domain::1220...). Admin topology APIs
// use the physical form via synchronizerStore; Ledger APIs reject it.
func LogicalSynchronizerID(syncID string) string {
	if strings.Count(syncID, "::") >= 2 {
		if i := strings.LastIndex(syncID, "::"); i > 0 {
			return syncID[:i]
		}
	}

	return syncID
}

// synchronizerLister is the minimal admin API surface needed to resolve aliases.
type synchronizerLister interface {
	ListConnectedSynchronizers(ctx context.Context) ([]SynchronizerInfo, error)
}

// ResolvePhysicalSynchronizerID maps a synchronizer alias or logical ID to the
// physical ID required by Canton Admin topology APIs. Values that already look
// physical (two or more "::" segments) are returned unchanged.
func ResolvePhysicalSynchronizerID(ctx context.Context, lister synchronizerLister, hint string) (string, error) {
	hint = strings.TrimSpace(hint)
	if hint == "" {
		return "", fmt.Errorf("synchronizer id is required")
	}
	if strings.Count(hint, "::") >= 2 {
		return hint, nil
	}

	syncs, err := lister.ListConnectedSynchronizers(ctx)
	if err != nil {
		return "", fmt.Errorf("listing connected synchronizers: %w", err)
	}
	if len(syncs) == 0 {
		return "", fmt.Errorf("participant has no connected synchronizers")
	}

	for _, s := range syncs {
		if s.Alias == hint || s.SynchronizerID == hint {
			return s.SynchronizerID, nil
		}
		if LogicalSynchronizerID(s.SynchronizerID) == hint {
			return s.SynchronizerID, nil
		}
	}

	// Common staging default: alias "global" when exactly one domain is connected.
	if hint == "global" {
		if len(syncs) == 1 {
			return syncs[0].SynchronizerID, nil
		}
		for _, s := range syncs {
			if s.Alias == "global" {
				return s.SynchronizerID, nil
			}
		}
	}

	return "", fmt.Errorf("synchronizer %q not found among %d connected synchronizer(s)", hint, len(syncs))
}
