package client

import "strings"

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
