package client

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awskms "github.com/aws/aws-sdk-go-v2/service/kms"
	awskmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"google.golang.org/protobuf/proto"
)

type mockKMSAPI struct {
	keySpec           awskmstypes.KeySpec
	signingAlgos      []awskmstypes.SigningAlgorithmSpec
	publicKey         []byte
	signature         []byte
	signInput         *awskms.SignInput
	getPublicKeyInput *awskms.GetPublicKeyInput
}

type mockProtocolKeyResolver struct {
	knownSigningKeys []string
	keyB64           string
}

func (m *mockProtocolKeyResolver) GetProtocolKeyFingerprint(_ context.Context, knownSigningKeys []string) (string, string, error) {
	m.knownSigningKeys = append([]string{}, knownSigningKeys...)
	keyB64 := m.keyB64
	if keyB64 == "" {
		keyB64 = "protocol-key-b64"
	}

	return "protocol-fp", keyB64, nil
}

func (m *mockKMSAPI) GetPublicKey(_ context.Context, in *awskms.GetPublicKeyInput, _ ...func(*awskms.Options)) (*awskms.GetPublicKeyOutput, error) {
	m.getPublicKeyInput = in
	return &awskms.GetPublicKeyOutput{
		KeySpec:           m.keySpec,
		PublicKey:         m.publicKey,
		SigningAlgorithms: m.signingAlgos,
	}, nil
}

func (m *mockKMSAPI) Sign(_ context.Context, in *awskms.SignInput, _ ...func(*awskms.Options)) (*awskms.SignOutput, error) {
	m.signInput = in
	return &awskms.SignOutput{
		Signature:        m.signature,
		SigningAlgorithm: in.SigningAlgorithm,
	}, nil
}

func TestAWSKMSSigner_SignsRawPreparedHashWithP256(t *testing.T) {
	t.Parallel()

	kmsAPI := &mockKMSAPI{
		keySpec:      awskmstypes.KeySpecEccNistP256,
		signingAlgos: []awskmstypes.SigningAlgorithmSpec{awskmstypes.SigningAlgorithmSpecEcdsaSha256},
		signature:    []byte{0x30, 0x44, 0x01, 0x02},
	}

	signer, err := NewAWSKMSSigner(t.Context(), kmsAPI, "arn:aws:kms:region:acct:key/proto", "protocol-fp")
	require.NoError(t, err)

	sig, err := signer.Sign(t.Context(), []byte{0xab, 0xcd})
	require.NoError(t, err)

	require.NotNil(t, kmsAPI.getPublicKeyInput)
	assert.Equal(t, "arn:aws:kms:region:acct:key/proto", *kmsAPI.getPublicKeyInput.KeyId)
	require.NotNil(t, kmsAPI.signInput)
	assert.Equal(t, "arn:aws:kms:region:acct:key/proto", *kmsAPI.signInput.KeyId)
	assert.Equal(t, []byte{0xab, 0xcd}, kmsAPI.signInput.Message)
	assert.Equal(t, awskmstypes.MessageTypeRaw, kmsAPI.signInput.MessageType)
	assert.Equal(t, awskmstypes.SigningAlgorithmSpecEcdsaSha256, kmsAPI.signInput.SigningAlgorithm)

	assert.Equal(t, apiv2.SignatureFormat_SIGNATURE_FORMAT_DER, sig.Format)
	assert.Equal(t, []byte{0x30, 0x44, 0x01, 0x02}, sig.Signature)
	assert.Equal(t, "protocol-fp", sig.SignedBy)
	assert.Equal(t, apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_EC_DSA_SHA_256, sig.SigningAlgorithmSpec)
}

func TestAWSKMSSigner_RejectsUnsupportedKeySpec(t *testing.T) {
	t.Parallel()

	_, err := NewAWSKMSSigner(t.Context(), &mockKMSAPI{
		keySpec: awskmstypes.KeySpecRsa2048,
	}, "key-id", "protocol-fp")
	require.ErrorContains(t, err, "unsupported AWS KMS signing key spec")
}

func TestTransactionSignerFactory_UsesAWSKMSWhenProtocolKeyConfigured(t *testing.T) {
	t.Parallel()

	cantonKeyB64, publicKey := encodedSigningPublicKey(t, "protocol-key")
	resolver := &mockProtocolKeyResolver{keyB64: cantonKeyB64}
	kmsAPI := &mockKMSAPI{
		keySpec:      awskmstypes.KeySpecEccNistP256,
		signingAlgos: []awskmstypes.SigningAlgorithmSpec{awskmstypes.SigningAlgorithmSpecEcdsaSha256},
		publicKey:    publicKey,
		signature:    []byte{0x30, 0x44},
	}
	factory := NewTransactionSignerFactory(resolver, nil, KMSConfig{
		ProtocolKeyID: "arn:aws:kms:region:acct:key/proto",
	}, kmsAPI)

	signer, err := factory(t.Context(), "p1", []string{"topology-key-1", "topology-key-2"})
	require.NoError(t, err)
	require.IsType(t, &AWSKMSSigner{}, signer)

	_, err = signer.Sign(t.Context(), []byte{0xab, 0xcd})
	require.NoError(t, err)

	assert.Equal(t, []string{"topology-key-1", "topology-key-2"}, resolver.knownSigningKeys)
	require.NotNil(t, kmsAPI.signInput)
	assert.Equal(t, awskmstypes.MessageTypeRaw, kmsAPI.signInput.MessageType)
	assert.Equal(t, awskmstypes.SigningAlgorithmSpecEcdsaSha256, kmsAPI.signInput.SigningAlgorithm)
}

func TestTransactionSignerFactory_RejectsMissingKMSClient(t *testing.T) {
	t.Parallel()

	cantonKeyB64, _ := encodedSigningPublicKey(t, "protocol-key")
	factory := NewTransactionSignerFactory(&mockProtocolKeyResolver{keyB64: cantonKeyB64}, nil, KMSConfig{
		ProtocolKeyID: "arn:aws:kms:region:acct:key/proto",
	}, nil)

	_, err := factory(t.Context(), "p1", []string{cantonKeyB64})
	require.ErrorContains(t, err, "aws kms client is required")
}

func TestTransactionSignerFactory_RejectsMissingVaultClient(t *testing.T) {
	t.Parallel()

	factory := NewTransactionSignerFactory(&mockProtocolKeyResolver{}, nil, KMSConfig{}, nil)

	_, err := factory(t.Context(), "p1", []string{"topology-key"})
	require.ErrorContains(t, err, "vault client is required")
}

func TestTransactionSignerFactory_RejectsKMSPublicKeyMismatch(t *testing.T) {
	t.Parallel()

	cantonKeyB64, _ := encodedSigningPublicKey(t, "active-protocol-key")
	_, differentPublicKey := encodedSigningPublicKey(t, "different-kms-key")
	kmsAPI := &mockKMSAPI{
		keySpec:      awskmstypes.KeySpecEccNistP256,
		signingAlgos: []awskmstypes.SigningAlgorithmSpec{awskmstypes.SigningAlgorithmSpecEcdsaSha256},
		publicKey:    differentPublicKey,
	}
	factory := NewTransactionSignerFactory(&mockProtocolKeyResolver{keyB64: cantonKeyB64}, nil, KMSConfig{
		ProtocolKeyID: "arn:aws:kms:region:acct:key/proto",
	}, kmsAPI)

	_, err := factory(t.Context(), "p1", []string{cantonKeyB64})
	require.ErrorContains(t, err, "does not match active Canton protocol key")
}

func encodedSigningPublicKey(t *testing.T, seed string) (string, []byte) {
	t.Helper()

	raw := sha256.Sum256([]byte(seed))
	key := &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}
	keyBytes, err := proto.Marshal(key)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(keyBytes), raw[:]
}
