package topology

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

type p2pMockClient struct {
	uid     string
	mapping *protov30.TopologyMapping
	serial  uint32
}

func (m *p2pMockClient) GetParticipantUID(context.Context) (string, error) { return m.uid, nil }

func (m *p2pMockClient) GetParticipantID(context.Context) (string, error) { return m.uid, nil }

func (m *p2pMockClient) GenerateSigningKey(context.Context, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) RegisterKmsSigningKey(context.Context, string, string, []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) GetNamespaceFingerprint(context.Context, string, string, []string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *p2pMockClient) GetNamespaceKeyName(context.Context, string, []string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *p2pMockClient) GetProtocolKeyFingerprint(context.Context, []string) (string, string, error) {
	return "", "", errors.New("not implemented")
}

func (m *p2pMockClient) Authorize(_ context.Context, serial uint32, mapping *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
	m.mapping = mapping
	m.serial = serial

	return &protov30.SignedTopologyTransaction{Transaction: []byte("tx"), Proposal: true}, nil
}

func (m *p2pMockClient) SignTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) ([]*protov30.SignedTopologyTransaction, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) AddTransactions(context.Context, []*protov30.SignedTopologyTransaction, string) error {
	return errors.New("not implemented")
}

func (m *p2pMockClient) DNSExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *p2pMockClient) P2PExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *p2pMockClient) NSDExists(context.Context, string, string) (bool, error) {
	return false, errors.New("not implemented")
}

func (m *p2pMockClient) GetDNS(context.Context, string, string) (*client.DNSState, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) GetP2P(context.Context, string, string) (*client.P2PState, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) ListDecentralizedNamespaces(context.Context, string) ([]*client.DNSState, error) {
	return nil, errors.New("not implemented")
}

func (m *p2pMockClient) UploadDar(context.Context, []byte) (string, error) {
	return "", errors.New("not implemented")
}

func encodedP2PSigningKey(t *testing.T, key string) string {
	t.Helper()

	b, err := proto.Marshal(&cryptov30.SigningPublicKey{
		PublicKey: []byte(key),
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	})
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(b)
}

func capturedPartySigningKeys(t *testing.T, mapping *protov30.TopologyMapping) ([][]byte, uint32) {
	t.Helper()

	p2pMapping, ok := mapping.GetMapping().(*protov30.TopologyMapping_PartyToParticipant)
	require.True(t, ok)
	signingKeys := p2pMapping.PartyToParticipant.GetPartySigningKeys()
	require.NotNil(t, signingKeys)

	keys := make([][]byte, 0, len(signingKeys.GetKeys()))
	for _, key := range signingKeys.GetKeys() {
		keys = append(keys, key.GetPublicKey())
	}

	return keys, signingKeys.GetThreshold()
}

func p2pDeps(c *p2pMockClient) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: c,
		Logger: logger.Nop(),
	}
}

func TestProposeAddP2POp_PreservesAndAppendsPartySigningKeys(t *testing.T) {
	t.Parallel()

	c := &p2pMockClient{uid: "p1"}
	_, err := operations.ExecuteOperation(optest.NewBundle(t), ProposeAddP2POp, p2pDeps(c), ProposeAddP2PInput{
		ParticipantID:       "p1",
		PartyID:             "party::namespace",
		AllParticipantUIDs:  []string{"p1", "p2", "p3"},
		NewP2PThreshold:     2,
		CurrentP2PSerial:    7,
		SynchronizerID:      "global",
		PartySigningKeysB64: []string{encodedP2PSigningKey(t, "p1-key"), encodedP2PSigningKey(t, "p2-key"), encodedP2PSigningKey(t, "p3-key")},
	})
	require.NoError(t, err)

	keys, threshold := capturedPartySigningKeys(t, c.mapping)
	assert.Equal(t, uint32(8), c.serial)
	assert.Equal(t, [][]byte{[]byte("p1-key"), []byte("p2-key"), []byte("p3-key")}, keys)
	assert.Equal(t, uint32(2), threshold)
}

func TestProposeKickP2POp_RemovesKickedPartySigningKey(t *testing.T) {
	t.Parallel()

	c := &p2pMockClient{uid: "p1"}
	_, err := operations.ExecuteOperation(optest.NewBundle(t), ProposeKickP2POp, p2pDeps(c), ProposeKickP2PInput{
		ParticipantID:         "p1",
		PartyID:               "party::namespace",
		RemainingParticipants: []string{"p1", "p3"},
		NewP2PThreshold:       2,
		CurrentP2PSerial:      3,
		SynchronizerID:        "global",
		PartySigningKeysB64:   []string{encodedP2PSigningKey(t, "p1-key"), encodedP2PSigningKey(t, "p3-key")},
	})
	require.NoError(t, err)

	keys, threshold := capturedPartySigningKeys(t, c.mapping)
	assert.Equal(t, uint32(4), c.serial)
	assert.Equal(t, [][]byte{[]byte("p1-key"), []byte("p3-key")}, keys)
	assert.Equal(t, uint32(2), threshold)
}

func TestProposeRotationP2POp_ReplacesOldPartySigningKey(t *testing.T) {
	t.Parallel()

	c := &p2pMockClient{uid: "p2"}
	_, err := operations.ExecuteOperation(optest.NewBundle(t), ProposeRotationP2POp, p2pDeps(c), ProposeRotationP2PInput{
		ParticipantID:             "p2",
		PartyID:                   "party::namespace",
		AllParticipantUIDs:        []string{"p1", "p2", "p3"},
		NewP2PThreshold:           2,
		CurrentP2PSerial:          4,
		SynchronizerID:            "global",
		CurrentSigningKeysB64:     []string{encodedP2PSigningKey(t, "p1-key"), encodedP2PSigningKey(t, "old-p2-key"), encodedP2PSigningKey(t, "p3-key")},
		OldDamlKeyB64:             encodedP2PSigningKey(t, "old-p2-key"),
		NewDamlKeyB64:             encodedP2PSigningKey(t, "new-p2-key"),
		PartySigningKeysThreshold: 2,
	})
	require.NoError(t, err)

	keys, threshold := capturedPartySigningKeys(t, c.mapping)
	assert.Equal(t, uint32(5), c.serial)
	assert.Equal(t, [][]byte{[]byte("p1-key"), []byte("new-p2-key"), []byte("p3-key")}, keys)
	assert.Equal(t, uint32(2), threshold)
}
