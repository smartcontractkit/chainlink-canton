package client

import (
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"fmt"
	"strings"

	v2crypto "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	versionv1 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/version/v1"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"
)

// VaultSigner implements [TransactionSigner] by exporting a signing key from
// Canton's VaultService and using it to sign externally via the Go standard
// library Ed25519 implementation.
//
// In Canton, DAML transactions that go through InteractiveSubmissionService are
// authorised by the signatories' namespace keys (the same keys used during the
// decentralized-party onboarding ceremony). Each participant signs with their
// own individual namespace signing key.
//
// Usage:
//
//	fingerprint, err := FindNamespaceKeyFingerprint(ctx, vaultClient)
//	signer, err := NewVaultSigner(ctx, vaultClient, fingerprint)
type VaultSigner struct {
	privateKey  ed25519.PrivateKey // 64-byte Ed25519 private key
	fingerprint string             // Canton key fingerprint ("1220" + hex(sha256(pubKey)))
}

// NewVaultSigner creates a VaultSigner by calling ExportKeyPair on the given
// vault client and extracting the Ed25519 seed from the returned proto.
//
// fingerprint identifies which key to export (see [FindNamespaceKeyFingerprint]).
// protocolVersion 34 matches Canton 3.4.x node serialisation.
func NewVaultSigner(ctx context.Context, vault cryptoadminv30.VaultServiceClient, fingerprint string) (*VaultSigner, error) {
	resp, err := vault.ExportKeyPair(ctx, &cryptoadminv30.ExportKeyPairRequest{
		Fingerprint:     fingerprint,
		ProtocolVersion: 34,
	})
	if err != nil {
		if strings.Contains(err.Error(), "does not support exporting") {
			return nil, fmt.Errorf("ExportKeyPair %q: %w (participant keys are KMS-backed: set kms_protocol_key_id in participant-config.json and use AWS credentials with kms:Sign)", fingerprint, err)
		}
		return nil, fmt.Errorf("ExportKeyPair %q: %w", fingerprint, err)
	}

	// ExportKeyPairResponse.key_pair is a serialised UntypedVersionedMessage
	// (Canton's generic versioned wrapper). Unwrap it to get the inner
	// CryptoKeyPair bytes.
	var versioned versionv1.UntypedVersionedMessage
	if err := proto.Unmarshal(resp.GetKeyPair(), &versioned); err != nil {
		return nil, fmt.Errorf("deserialising versioned wrapper for %q: %w", fingerprint, err)
	}

	// The response bytes are a serialised v30.CryptoKeyPair proto.
	var pair cryptov30.CryptoKeyPair
	if err := proto.Unmarshal(versioned.GetData(), &pair); err != nil {
		return nil, fmt.Errorf("deserialising exported key pair for %q: %w", fingerprint, err)
	}

	skp := pair.GetSigningKeyPair()
	if skp == nil {
		return nil, fmt.Errorf("exported key for fingerprint %q is not a signing key pair", fingerprint)
	}

	seed := skp.GetPrivateKey().GetPrivateKey()

	// Canton may export the private key as a raw 32-byte seed or as a
	// PKCS8 DER-encoded key (48 bytes for Ed25519). Handle both.
	var privKey ed25519.PrivateKey
	switch len(seed) {
	case ed25519.SeedSize:
		// Raw 32-byte seed — expand directly.
		privKey = ed25519.NewKeyFromSeed(seed)
	case 48:
		// PKCS8 DER-encoded Ed25519 private key.
		parsed, err := x509.ParsePKCS8PrivateKey(seed)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS8 private key for %q: %w", fingerprint, err)
		}
		var ok bool
		privKey, ok = parsed.(ed25519.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PKCS8 key for %q is not Ed25519", fingerprint)
		}
	default:
		return nil, fmt.Errorf("unexpected Ed25519 seed length for %q: got %d bytes, want %d",
			fingerprint, len(seed), ed25519.SeedSize)
	}

	return &VaultSigner{
		privateKey:  privKey,
		fingerprint: fingerprint,
	}, nil
}

// Sign signs hash using the exported Ed25519 key and returns a Ledger API
// Signature proto with format RAW and algorithm ED25519.
func (s *VaultSigner) Sign(_ context.Context, hash []byte) (*v2crypto.Signature, error) {
	sig := ed25519.Sign(s.privateKey, hash)
	return &v2crypto.Signature{
		Format:               v2crypto.SignatureFormat_SIGNATURE_FORMAT_RAW,
		Signature:            sig,
		SignedBy:             s.fingerprint,
		SigningAlgorithmSpec: v2crypto.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_ED25519,
	}, nil
}

// FindNamespaceKeyFingerprint lists the NAMESPACE signing keys in the vault
// and returns the fingerprint of the first one found.
//
// Each participant in a decentralized-party ceremony has exactly one NAMESPACE
// signing key (generated during onboarding). In Canton, the signing key fingerprint
// is computed as the multihash "1220" + hex(sha256(rawPublicKeyBytes)).
//
// This fingerprint is also the participant's individual namespace identifier
// within the DecentralizedNamespaceDefinition.
func FindNamespaceKeyFingerprint(ctx context.Context, vault cryptoadminv30.VaultServiceClient) (string, error) {
	return findSigningKeyFingerprint(ctx, vault,
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
		"NAMESPACE",
	)
}

// FindProtocolKeyFingerprint lists the PROTOCOL (DAML) signing keys in the
// vault and returns the fingerprint of the first one found.
//
// Each participant generates a PROTOCOL signing key during the onboarding
// ceremony (named "<namespace>-protocol"). This key is registered in
// PartyToParticipant.PartySigningKeys and is used to authorise DAML
// transactions submitted via InteractiveSubmissionService.
func FindProtocolKeyFingerprint(ctx context.Context, vault cryptoadminv30.VaultServiceClient) (string, error) {
	return findSigningKeyFingerprint(ctx, vault,
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
		"PROTOCOL",
	)
}

// findSigningKeyFingerprint is the shared implementation for FindNamespaceKeyFingerprint
// and FindProtocolKeyFingerprint.
func findSigningKeyFingerprint(ctx context.Context, vault cryptoadminv30.VaultServiceClient, usage cryptov30.SigningKeyUsage, usageName string) (string, error) {
	resp, err := vault.ListMyKeys(ctx, &cryptoadminv30.ListMyKeysRequest{
		Filters: &cryptoadminv30.ListKeysFilters{
			Usage: []cryptov30.SigningKeyUsage{usage},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ListMyKeys (%s usage): %w", usageName, err)
	}

	keys := resp.GetPrivateKeysMetadata()
	if len(keys) == 0 {
		return "", fmt.Errorf("no %s signing key found in vault", usageName)
	}

	// Iterate over all matching keys and return the fingerprint of the first
	// one that has a non-empty raw public key. Pre-existing internal Canton
	// keys may have a nil SigningPublicKey in the PublicKeyWithName wrapper;
	// we skip those so that ceremony-generated keys are found reliably.
	for _, k := range keys {
		spk := k.GetPublicKeyWithName().GetPublicKey().GetSigningPublicKey()
		if spk == nil {
			continue
		}
		rawBytes := spk.GetPublicKey()
		if len(rawBytes) == 0 {
			continue
		}
		// Use the canonical Canton fingerprint algorithm:
		// sha256([0,0,0,12] + rawKey), prefixed with multihash header [0x12, 0x20].
		// Purpose 12 = PublicKeyFingerprint (HashPurpose.scala).
		return helpers.GetPublicKeyFingerprint(rawBytes)
	}

	return "", fmt.Errorf("no %s signing key with accessible public key bytes found in vault", usageName)
}
