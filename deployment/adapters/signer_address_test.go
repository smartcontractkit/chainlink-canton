package adapters

import (
	"encoding/hex"
	"strings"
	"testing"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A canton NOP's signer identity is a raw 65-byte secp256k1 public key, but the
// on-chain committee verifier stores the derived 20-byte ECDSA address. The canton
// address normalizer must collapse both representations to the same canonical form
// so ccv state inference can map on-chain committee signers back to canton NOPs.
func TestNormalizeCantonSignerAddress_PubkeyAndAddressCollapse(t *testing.T) {
	key, err := gethcrypto.GenerateKey()
	require.NoError(t, err)

	pubkey := gethcrypto.FromECDSAPub(&key.PublicKey) // 65-byte uncompressed
	require.Len(t, pubkey, 65)
	pubkeyHex := hex.EncodeToString(pubkey)
	addr := gethcrypto.PubkeyToAddress(key.PublicKey).Hex() // checksummed 0x...

	want := strings.ToLower(addr)

	// Raw pubkey (with and without 0x), the checksummed address, and the already
	// canonical form all normalize to the same lowercase derived address.
	assert.Equal(t, want, normalizeCantonSignerAddress(pubkeyHex), "raw pubkey")
	assert.Equal(t, want, normalizeCantonSignerAddress("0x"+pubkeyHex), "0x-prefixed pubkey")
	assert.Equal(t, want, normalizeCantonSignerAddress(addr), "checksummed address")
	assert.Equal(t, want, normalizeCantonSignerAddress(want), "lowercase address (idempotent)")
}
