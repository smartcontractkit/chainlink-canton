package helpers_test

import (
	"encoding/base64"
	"testing"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chainlink/canton-party-ceremony/internal/helpers"
)

// TestGetPublicKeyFingerprint_KnownVector uses the example from the Canton docs:
// https://docs.digitalasset.com/operate/3.4/tutorials/external-signing/topology-transaction.html
//
//	> compute_canton_fingerprint_from_base64 "2RwUiIHVUVdulxzD8NKtPmIaaBqMer1A90rDjoklJPY="
//	1220205057e331cc8929dd217e2f8e63f503b7081773de60d01fb46839700bc5caaa
func TestGetPublicKeyFingerprint_KnownVector(t *testing.T) {
	t.Parallel()

	keyBase64 := "2RwUiIHVUVdulxzD8NKtPmIaaBqMer1A90rDjoklJPY="
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	require.NoError(t, err)

	fingerprint, err := helpers.GetPublicKeyFingerprint(keyBytes)
	require.NoError(t, err)
	assert.Equal(t, "1220205057e331cc8929dd217e2f8e63f503b7081773de60d01fb46839700bc5caaa", fingerprint)
}

func TestGetPublicKeyFingerprint_KnownVector2(t *testing.T) {
	t.Parallel()

	keyBase64 := "MCowBQYDK2VwAyEATrcRcTwTV44cpws+3IAmiSlzWK3DqnMnAsMztzt84Tc="
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	require.NoError(t, err)

	fingerprint, err := helpers.GetPublicKeyFingerprint(keyBytes)
	require.NoError(t, err)
	assert.Equal(t, "1220e4593688f14d6a917034ce54c4a9fc5410221de4dadb93ac4326fb4221d70b6d", fingerprint)
}

func TestGetPublicKeyFingerprint_Format(t *testing.T) {
	t.Parallel()

	key := []byte("any 32-byte-ish key material here")
	fp, err := helpers.GetPublicKeyFingerprint(key)
	require.NoError(t, err)
	assert.Len(t, fp, 68, "fingerprint must be 68 hex chars (34 bytes)")
	assert.Equal(t, "1220", fp[:4], "must start with multihash prefix 0x12 0x20")
}

func TestGetPublicKeyFingerprint_Deterministic(t *testing.T) {
	t.Parallel()

	key := []byte{0x01, 0x02, 0x03, 0x04}
	fp1, err := helpers.GetPublicKeyFingerprint(key)
	require.NoError(t, err)
	fp2, err := helpers.GetPublicKeyFingerprint(key)
	require.NoError(t, err)
	assert.Equal(t, fp1, fp2)
}

func TestGetPublicKeyFingerprint_DifferentInputsDifferentOutputs(t *testing.T) {
	t.Parallel()

	fp1, err := helpers.GetPublicKeyFingerprint([]byte{0x01})
	require.NoError(t, err)
	fp2, err := helpers.GetPublicKeyFingerprint([]byte{0x02})
	require.NoError(t, err)
	assert.NotEqual(t, fp1, fp2)
}

func TestNewPublicKeyFromHex_Valid(t *testing.T) {
	t.Parallel()

	keyBytes := []byte{0xde, 0xad, 0xbe, 0xef, 0x01, 0x02, 0x03, 0x04}
	hexStr := "deadbeef01020304"

	key, err := helpers.NewPublicKeyFromHex(hexStr, cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE)
	require.NoError(t, err)
	assert.Equal(t, keyBytes, key.GetPublicKey())
	assert.Equal(t, cryptov30.CryptoKeyFormat_CRYPTO_KEY_FORMAT_RAW, key.GetFormat())
	assert.Equal(t, cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519, key.GetKeySpec())
	assert.Equal(t, []cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE}, key.GetUsage())
}

func TestNewPublicKeyFromHex_InvalidHex(t *testing.T) {
	t.Parallel()

	_, err := helpers.NewPublicKeyFromHex("not-valid-hex", cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid hex public key")
}

func TestNewPublicKeyFromHex_FingerprintRoundtrip(t *testing.T) {
	t.Parallel()

	// Verify that bytes survive the hex round-trip and produce the expected fingerprint.
	keyBase64 := "2RwUiIHVUVdulxzD8NKtPmIaaBqMer1A90rDjoklJPY="
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	require.NoError(t, err)

	hexStr := make([]byte, len(keyBytes)*2)
	const hextable = "0123456789abcdef"
	for i, b := range keyBytes {
		hexStr[i*2] = hextable[b>>4]
		hexStr[i*2+1] = hextable[b&0x0f]
	}

	key, err := helpers.NewPublicKeyFromHex(string(hexStr), cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE)
	require.NoError(t, err)

	fp, err := helpers.GetPublicKeyFingerprint(key.GetPublicKey())
	require.NoError(t, err)
	assert.Equal(t, "1220205057e331cc8929dd217e2f8e63f503b7081773de60d01fb46839700bc5caaa", fp)
}
