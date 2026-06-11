package keys

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	optest "github.com/smartcontractkit/chainlink-deployments-framework/operations/test"
)

type kmsVaultMockClient struct {
	participantID string
	keys          map[string]*cryptov30.SigningPublicKey
}

func newKmsVaultMockClient(participantID string) *kmsVaultMockClient {
	return &kmsVaultMockClient{
		participantID: participantID,
		keys:          make(map[string]*cryptov30.SigningPublicKey),
	}
}

func (m *kmsVaultMockClient) keyFor(kmsKeyID string) *cryptov30.SigningPublicKey {
	raw := sha256.Sum256([]byte(m.participantID + ":" + kmsKeyID))
	return &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}
}

func (m *kmsVaultMockClient) GetParticipantUID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *kmsVaultMockClient) GetParticipantID(context.Context) (string, error) {
	return m.participantID, nil
}

func (m *kmsVaultMockClient) GenerateSigningKey(context.Context, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) RegisterKmsSigningKey(_ context.Context, kmsKeyID, name string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	if existing, ok := m.keys[kmsKeyID]; ok {
		return nil, fmt.Errorf("RegisterKmsSigningKey: rpc error: code = Internal desc = Existing public key for %x is different than inserted key (name %q)", existing.GetPublicKey()[:4], name)
	}
	key := m.keyFor(kmsKeyID)
	m.keys[kmsKeyID] = key

	return key, nil
}

func (m *kmsVaultMockClient) LookupKmsSigningKey(_ context.Context, kmsKeyID string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	key, ok := m.keys[kmsKeyID]
	if !ok {
		return nil, fmt.Errorf("no registered KMS signing key %q found in vault", kmsKeyID)
	}

	return key, nil
}

func (m *kmsVaultMockClient) GetNamespaceFingerprint(context.Context, string, string, []string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *kmsVaultMockClient) GetNamespaceKeyName(context.Context, string, []string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *kmsVaultMockClient) GetProtocolKeyFingerprint(context.Context, []string) (string, string, error) {
	return "", "", errors.New("not implemented")
}

func (m *kmsVaultMockClient) Authorize(context.Context, uint32, *protov30.TopologyMapping, string, bool, ...string) (*protov30.SignedTopologyTransaction, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) SignTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) ([]*protov30.SignedTopologyTransaction, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) AddTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) error {
	return errors.New("not implemented")
}

func (m *kmsVaultMockClient) DNSExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *kmsVaultMockClient) P2PExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *kmsVaultMockClient) NSDExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *kmsVaultMockClient) GetDNS(context.Context, string, string) (*client.DNSState, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) GetP2P(context.Context, string, string) (*client.P2PState, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) ListDecentralizedNamespaces(context.Context, string) ([]*client.DNSState, error) {
	return nil, errors.New("not implemented")
}

func (m *kmsVaultMockClient) UploadDar(context.Context, []byte) (string, error) {
	return "", errors.New("not implemented")
}

func TestObtainSigningKey_ReusesRegisteredKmsKeyOnConflict(t *testing.T) {
	t.Parallel()

	mock := newKmsVaultMockClient("p1")
	const kmsARN = "arn:aws:kms:us-west-2:123456789:key/abc"

	_, err := mock.RegisterKmsSigningKey(context.Background(), kmsARN, "ccip-owner-mainnet", []cryptov30.SigningKeyUsage{
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
	})
	require.NoError(t, err)

	key, err := obtainSigningKey(context.Background(), mock, kmsARN, "ccv-owner-mainnet", []cryptov30.SigningKeyUsage{
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
	})
	require.NoError(t, err)
	assert.Equal(t, mock.keyFor(kmsARN).GetPublicKey(), key.GetPublicKey())
}

func TestVaultRegistrationName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "ccv-owner-mainnet", vaultRegistrationName("ccv-owner-mainnet", ""))
	assert.Equal(t, "ccip-owner-mainnet", vaultRegistrationName("ccv-owner-mainnet", "ccip-owner-mainnet"))
}

func TestCreateMemberKeyOp_ReusesExistingKmsKeysAcrossVaultNames(t *testing.T) {
	t.Parallel()

	mock := newKmsVaultMockClient("Chainlink-MainnetCV1-1::1220abc")
	const (
		nsARN = "arn:aws:kms:us-west-2:123456789:key/namespace"
		pARN  = "arn:aws:kms:us-west-2:123456789:key/protocol"
	)

	_, err := mock.RegisterKmsSigningKey(context.Background(), nsARN, "ccip-owner-mainnet", []cryptov30.SigningKeyUsage{
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
	})
	require.NoError(t, err)
	_, err = mock.RegisterKmsSigningKey(context.Background(), pARN, "ccip-owner-mainnet-protocol", []cryptov30.SigningKeyUsage{
		cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
	})
	require.NoError(t, err)

	out, err := operations.ExecuteOperation(optest.NewBundle(t), CreateMemberKeyOp, ceremony.CantonDeps{
		Client: mock,
		KMS: client.KMSConfig{
			NamespaceKeyID: nsARN,
			ProtocolKeyID:  pARN,
		},
		Logger: logger.Nop(),
	}, CreateMemberKeyInput{
		NamespaceName: "ccv-owner-mainnet",
		ParticipantID: "Chainlink-MainnetCV1-1::1220abc",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, out.Output.NamespaceFingerprint)
	assert.NotEmpty(t, out.Output.SigningKeyB64)
	assert.NotEmpty(t, out.Output.DamlKeyB64)
}

func TestCreateMemberKeyOp_KmsVaultNameOverridesRegistrationName(t *testing.T) {
	t.Parallel()

	mock := newKmsVaultMockClient("p1")
	const (
		nsARN = "arn:aws:kms:us-west-2:123456789:key/namespace"
		pARN  = "arn:aws:kms:us-west-2:123456789:key/protocol"
	)

	_, err := operations.ExecuteOperation(optest.NewBundle(t), CreateMemberKeyOp, ceremony.CantonDeps{
		Client: mock,
		KMS: client.KMSConfig{
			NamespaceKeyID: nsARN,
			ProtocolKeyID:  pARN,
		},
		Logger: logger.Nop(),
	}, CreateMemberKeyInput{
		NamespaceName: "ccv-owner-mainnet",
		KmsVaultName:  "ccip-owner-mainnet",
		ParticipantID: "p1",
	})
	require.NoError(t, err)
	assert.Contains(t, mock.keys, nsARN)
	assert.Contains(t, mock.keys, pARN)
}
