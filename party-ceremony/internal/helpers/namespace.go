package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// ComputeDecentralizedNamespace replicates Canton's
// DecentralizedNamespaceDefinition.computeNamespace(owners) algorithm.
//
// Algorithm (from Canton's HashBuilder / DeterministicEncoding):
//  1. Seed SHA-256 with the 4-byte big-endian purpose value 37.
//  2. Sort owner namespace fingerprints lexicographically.
//  3. For each owner: write 4-byte big-endian length, then the UTF-8 bytes.
//  4. Prepend multihash header [0x12, 0x20] (SHA-256, 32 bytes) to the digest.
//  5. Hex-encode the resulting 34 bytes → 68-character fingerprint string.
func ComputeDecentralizedNamespace(owners []string) string {
	sorted := make([]string, len(owners))
	copy(sorted, owners)
	sort.Strings(sorted)

	h := sha256.New()

	// Purpose = 37 ("DecentralizedNamespace"), big-endian int32.
	const purpose = 37
	h.Write([]byte{0, 0, 0, purpose})

	for _, ns := range sorted {
		b := []byte(ns)
		l := len(b)
		h.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)}) //nolint:gosec // TODO
		h.Write(b)
	}

	digest := h.Sum(nil)
	result := make([]byte, 0, 34)
	result = append(result, 0x12, 0x20)
	result = append(result, digest...)

	return hex.EncodeToString(result)
}
