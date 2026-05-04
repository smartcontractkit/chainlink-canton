package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"slices"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	awskmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"
)

// AWSKMSAPI is the narrow subset of the AWS SDK v2 KMS client used by
// AWSKMSSigner. Tests can satisfy this interface with small function mocks.
type AWSKMSAPI interface {
	GetPublicKey(ctx context.Context, in *awskms.GetPublicKeyInput, optFns ...func(*awskms.Options)) (*awskms.GetPublicKeyOutput, error)
	Sign(ctx context.Context, in *awskms.SignInput, optFns ...func(*awskms.Options)) (*awskms.SignOutput, error)
}

// AWSKMSSigner signs prepared Ledger API transaction hashes with an asymmetric
// AWS KMS key registered as the participant's Canton PROTOCOL signing key.
type AWSKMSSigner struct {
	kms         AWSKMSAPI
	keyID       string
	fingerprint string
	algorithm   awskmstypes.SigningAlgorithmSpec
	ledgerAlgo  apiv2.SigningAlgorithmSpec
}

func NewAWSKMSSigner(ctx context.Context, kmsAPI AWSKMSAPI, keyID string, fingerprint string, activeSigningKeyB64 ...string) (*AWSKMSSigner, error) {
	if kmsAPI == nil {
		return nil, fmt.Errorf("aws kms client is required")
	}
	if keyID == "" {
		return nil, fmt.Errorf("aws kms key id is required")
	}
	if fingerprint == "" {
		return nil, fmt.Errorf("signing key fingerprint is required")
	}

	pub, err := kmsAPI.GetPublicKey(ctx, &awskms.GetPublicKeyInput{KeyId: &keyID})
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey %q: %w", keyID, err)
	}

	algorithm, ledgerAlgo, err := signingAlgorithmsForKey(pub.KeySpec)
	if err != nil {
		return nil, err
	}
	if len(pub.SigningAlgorithms) > 0 && !containsSigningAlgorithm(pub.SigningAlgorithms, algorithm) {
		return nil, fmt.Errorf("KMS key %q does not allow signing algorithm %s", keyID, algorithm)
	}
	if len(activeSigningKeyB64) > 0 && activeSigningKeyB64[0] != "" {
		if err := validateKMSPublicKeyMatchesCanton(pub.PublicKey, activeSigningKeyB64[0], fingerprint); err != nil {
			return nil, fmt.Errorf("KMS key %q does not match active Canton protocol key %q: %w", keyID, fingerprint, err)
		}
	}

	return &AWSKMSSigner{
		kms:         kmsAPI,
		keyID:       keyID,
		fingerprint: fingerprint,
		algorithm:   algorithm,
		ledgerAlgo:  ledgerAlgo,
	}, nil
}

func (s *AWSKMSSigner) Sign(ctx context.Context, hash []byte) (*apiv2.Signature, error) {
	out, err := s.kms.Sign(ctx, &awskms.SignInput{
		KeyId:            &s.keyID,
		Message:          hash,
		MessageType:      awskmstypes.MessageTypeRaw,
		SigningAlgorithm: s.algorithm,
	})
	if err != nil {
		return nil, fmt.Errorf("KMS Sign %q: %w", s.keyID, err)
	}

	return &apiv2.Signature{
		Format:               apiv2.SignatureFormat_SIGNATURE_FORMAT_DER,
		Signature:            out.Signature,
		SignedBy:             s.fingerprint,
		SigningAlgorithmSpec: s.ledgerAlgo,
	}, nil
}

func NewTransactionSignerFactory(
	resolver ProtocolKeyResolver,
	vault cryptoadminv30.VaultServiceClient,
	kmsCfg KMSConfig,
	kmsAPI AWSKMSAPI,
) TransactionSignerFactory {
	return func(ctx context.Context, _ string, knownSigningKeysB64 []string) (TransactionSigner, error) {
		fp, activeKeyB64, err := resolver.GetProtocolKeyFingerprint(ctx, knownSigningKeysB64)
		if err != nil {
			return nil, fmt.Errorf("discovering protocol signing key: %w", err)
		}

		if kmsCfg.ProtocolKeyID != "" {
			if kmsAPI == nil {
				return nil, fmt.Errorf("aws kms client is required when kms_protocol_key_id is configured")
			}

			return NewAWSKMSSigner(ctx, kmsAPI, kmsCfg.ProtocolKeyID, fp, activeKeyB64)
		}
		if vault == nil {
			return nil, fmt.Errorf("vault client is required when kms_protocol_key_id is not configured")
		}

		return NewVaultSigner(ctx, vault, fp)
	}
}

func validateKMSPublicKeyMatchesCanton(kmsPublicKeyDER []byte, cantonSigningKeyB64 string, fingerprint string) error {
	cantonKeyBytes, err := base64.StdEncoding.DecodeString(cantonSigningKeyB64)
	if err != nil {
		return fmt.Errorf("decoding active Canton signing key: %w", err)
	}
	var cantonKey cryptov30.SigningPublicKey
	if err := proto.Unmarshal(cantonKeyBytes, &cantonKey); err != nil {
		return fmt.Errorf("unmarshalling active Canton signing key: %w", err)
	}
	cantonPublicKey := cantonKey.GetPublicKey()
	if len(cantonPublicKey) == 0 {
		return fmt.Errorf("active Canton signing key has no public key bytes")
	}

	for _, candidate := range kmsPublicKeyCandidates(kmsPublicKeyDER) {
		if bytes.Equal(candidate, cantonPublicKey) {
			return nil
		}
		if fingerprint != "" {
			fp, fpErr := helpers.GetPublicKeyFingerprint(candidate)
			if fpErr == nil && fp == fingerprint {
				return nil
			}
		}
	}

	return fmt.Errorf("public key bytes differ")
}

func kmsPublicKeyCandidates(kmsPublicKeyDER []byte) [][]byte {
	if len(kmsPublicKeyDER) == 0 {
		return nil
	}

	candidates := [][]byte{kmsPublicKeyDER}
	pub, err := x509.ParsePKIXPublicKey(kmsPublicKeyDER)
	if err != nil {
		return candidates
	}
	if ecdsaPub, ok := pub.(*ecdsa.PublicKey); ok {
		raw, bytesErr := ecdsaPub.Bytes()
		if bytesErr == nil && len(raw) > 0 {
			candidates = append(candidates, raw)
		}
	}

	return candidates
}

func signingAlgorithmsForKey(keySpec awskmstypes.KeySpec) (awskmstypes.SigningAlgorithmSpec, apiv2.SigningAlgorithmSpec, error) {
	switch keySpec {
	case awskmstypes.KeySpecEccNistP256:
		return awskmstypes.SigningAlgorithmSpecEcdsaSha256,
			apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_EC_DSA_SHA_256,
			nil
	case awskmstypes.KeySpecEccNistP384:
		return awskmstypes.SigningAlgorithmSpecEcdsaSha384,
			apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_EC_DSA_SHA_384,
			nil
	case awskmstypes.KeySpecRsa2048,
		awskmstypes.KeySpecRsa3072,
		awskmstypes.KeySpecRsa4096,
		awskmstypes.KeySpecEccNistP521,
		awskmstypes.KeySpecEccSecgP256k1,
		awskmstypes.KeySpecSymmetricDefault,
		awskmstypes.KeySpecHmac224,
		awskmstypes.KeySpecHmac256,
		awskmstypes.KeySpecHmac384,
		awskmstypes.KeySpecHmac512,
		awskmstypes.KeySpecSm2,
		awskmstypes.KeySpecMlDsa44,
		awskmstypes.KeySpecMlDsa65,
		awskmstypes.KeySpecMlDsa87,
		awskmstypes.KeySpecEccNistEdwards25519:
		return "", apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_UNSPECIFIED,
			fmt.Errorf("unsupported AWS KMS signing key spec %s; expected ECC_NIST_P256 or ECC_NIST_P384", keySpec)
	default:
		return "", apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_UNSPECIFIED,
			fmt.Errorf("unknown AWS KMS signing key spec %s; expected ECC_NIST_P256 or ECC_NIST_P384", keySpec)
	}
}

func containsSigningAlgorithm(allowed []awskmstypes.SigningAlgorithmSpec, expected awskmstypes.SigningAlgorithmSpec) bool {
	return slices.Contains(allowed, expected)
}
