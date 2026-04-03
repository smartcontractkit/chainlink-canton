package helpers

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
)

// GetPublicKeyFingerprint computes the Canton fingerprint of a public key.
//
// The input may be either raw key bytes or a DER-encoded SubjectPublicKeyInfo
// (PKIX) structure. When a DER-encoded key is detected, the raw key bytes are
// extracted first so that the fingerprint is always computed over the bare key.
//
// Algorithm (per Canton's HashPurpose.scala, hash purpose 12 = PublicKeyFingerprint):
//  1. Prefix the raw key bytes with the 4-byte big-endian hash purpose value 12.
//  2. SHA-256 hash the result.
//  3. Prepend the multihash header [0x12, 0x20] (SHA-256, 32 bytes).
//  4. Hex-encode the resulting 34 bytes → 68-character fingerprint string.
func GetPublicKeyFingerprint(key []byte) (string, error) {
	// If the caller passed a DER-encoded SubjectPublicKeyInfo, extract the raw
	// key bytes. Canton always fingerprints the bare key, not the DER wrapper.
	if pub, err := x509.ParsePKIXPublicKey(key); err == nil {
		if ed25519Key, ok := pub.(ed25519.PublicKey); ok {
			key = []byte(ed25519Key)
		}
	}

	h := sha256.New()

	// Purpose = 12 ("PublicKeyFingerprint"), big-endian int32.
	const purpose = 12
	h.Write([]byte{0, 0, 0, purpose})
	h.Write(key)

	digest := h.Sum(nil)
	result := make([]byte, 0, 34)
	result = append(result, 0x12, 0x20)
	result = append(result, digest...)

	return hex.EncodeToString(result), nil
}

func NewPublicKeyFromHex(hexStr string, usage cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	keyBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex public key: %w", err)
	}

	return &cryptov30.SigningPublicKey{
		PublicKey: keyBytes,
		Format:    cryptov30.CryptoKeyFormat_CRYPTO_KEY_FORMAT_RAW,
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
		Usage:     []cryptov30.SigningKeyUsage{usage},
		Scheme:    cryptov30.SigningKeyScheme_SIGNING_KEY_SCHEME_ED25519,
	}, nil
}
