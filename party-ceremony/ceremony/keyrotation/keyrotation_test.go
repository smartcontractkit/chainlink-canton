package keyrotation_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
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
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"
)

// ── Constants shared across tests ────────────────────────────────────────────

const (
	testPartyID    = "test-party::1220aabbccdd"
	testSyncID     = "global"
	testTargetUID  = "p1"
	testTargetNSFP = "fp-owner-p1"
	testNSName     = "test-rotation-ns"
)

// ── Deterministic test signing keys ─────────────────────────────────────────

// makeTestSigningKeyB64 creates a deterministic base64-encoded proto
// SigningPublicKey from a seed string.
func makeTestSigningKeyB64(seed string) string {
	raw := sha256.Sum256([]byte(seed))
	key := &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}
	keyBytes, _ := proto.Marshal(key)

	return base64.StdEncoding.EncodeToString(keyBytes)
}

// computeExpectedNewNamespaceFP computes the expected new namespace fingerprint
// that GenerateRotatedKeyOp will produce via the mock's GenerateSigningKey.
func computeExpectedNewNamespaceFP() string {
	raw := sha256.Sum256([]byte(testTargetUID + ":" + testNSName + "-rotated"))
	fp, _ := helpers.GetPublicKeyFingerprint(raw[:])

	return fp
}

func computeExpectedNewDamlKeyB64() string {
	raw := sha256.Sum256([]byte(testTargetUID + ":" + testNSName + "-protocol-rotated"))
	key := &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}
	keyBytes, _ := proto.Marshal(key)

	return base64.StdEncoding.EncodeToString(keyBytes)
}

var (
	// testSigningKeyTargetB64 is the target's (p1) current DAML signing key.
	// GetProtocolKeyFingerprint returns this as the "old" key.
	testSigningKeyTargetB64 = makeTestSigningKeyB64("target-daml")

	// testSigningKeyOtherB64 is the other participant's (p2) DAML signing key.
	testSigningKeyOtherB64 = makeTestSigningKeyB64("other-daml")

	// testNewNamespaceFP is the expected fingerprint of the rotated namespace
	// key, computed from the mock's deterministic GenerateSigningKey output.
	testNewNamespaceFP = computeExpectedNewNamespaceFP()

	// testNewDamlKeyB64 is the expected generated DAML key for the non-KMS
	// rotation path.
	testNewDamlKeyB64 = computeExpectedNewDamlKeyB64()
)

// ── Mock shared state ────────────────────────────────────────────────────────

// mockState is shared across all mockCantonClient instances within a single
// test scenario. It tracks topology state changes that Canton would record.
type mockState struct {
	mu              sync.Mutex
	nsdProposed     bool // true after ProposeNewNSDOp (NamespaceDelegation Authorize)
	dnsSubmitCount  int  // incremented by each AddTransactions call
	p2pProposeCount int  // incremented by each PartyToParticipant Authorize call
}

func (s *mockState) onNSDProposed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nsdProposed = true
}

func (s *mockState) hasNSDProposed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nsdProposed
}

func (s *mockState) onDNSSubmit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dnsSubmitCount++
}

func (s *mockState) dnsSubmitted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.dnsSubmitCount > 0
}

func (s *mockState) onP2PPropose() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.p2pProposeCount++
}

func (s *mockState) p2pReady() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.p2pProposeCount > 0
}

// ── Mock Canton client ───────────────────────────────────────────────────────

// mockCantonClient is a deterministic, in-memory implementation of
// [client.CantonClient] for key rotation ceremony unit tests.
type mockCantonClient struct {
	participantID    string
	state            *mockState // nil for single-run tests that never reach polling
	kmsRegistrations *int
}

func (m *mockCantonClient) GetParticipantUID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockCantonClient) GetParticipantID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockCantonClient) GenerateSigningKey(_ context.Context, name string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	raw := sha256.Sum256([]byte(m.participantID + ":" + name))

	return &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}, nil
}

func (m *mockCantonClient) RegisterKmsSigningKey(_ context.Context, kmsKeyID string, name string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	if m.kmsRegistrations != nil {
		(*m.kmsRegistrations)++
	}
	raw := sha256.Sum256([]byte(m.participantID + ":" + kmsKeyID + ":" + name))

	return &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}, nil
}

func (m *mockCantonClient) GetNamespaceFingerprint(_ context.Context, keyName string, _ string, _ []string) (string, error) {
	raw := sha256.Sum256([]byte(m.participantID + ":" + keyName))

	return fmt.Sprintf("1220%x", raw[:8]), nil
}

// GetNamespaceKeyName returns the test namespace key name, simulating
// auto-discovery of the existing key name from the vault.
func (m *mockCantonClient) GetNamespaceKeyName(_ context.Context, _ string, _ []string) (string, error) {
	return testNSName, nil
}

// GetProtocolKeyFingerprint returns the target's current DAML signing key.
// This is only called by GenerateRotatedKeyOp on the target participant's node.
func (m *mockCantonClient) GetProtocolKeyFingerprint(_ context.Context, _ []string) (string, string, error) {
	return "mock-old-daml-fp", testSigningKeyTargetB64, nil
}

func (m *mockCantonClient) Authorize(_ context.Context, serial uint32, mapping *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
	if m.state != nil {
		switch mapping.GetMapping().(type) {
		case *protov30.TopologyMapping_NamespaceDelegation:
			m.state.onNSDProposed()
		case *protov30.TopologyMapping_PartyToParticipant:
			m.state.onP2PPropose()
		}
	}
	payload, _ := proto.Marshal(mapping)
	raw := sha256.Sum256(append(payload, byte(serial)))

	return &protov30.SignedTopologyTransaction{
		Transaction: raw[:],
		Proposal:    true,
	}, nil
}

func (m *mockCantonClient) SignTransactions(_ context.Context, txs []*protov30.SignedTopologyTransaction, _ string) ([]*protov30.SignedTopologyTransaction, error) {
	result := make([]*protov30.SignedTopologyTransaction, len(txs))
	for i, tx := range txs {
		sigBytes := sha256.Sum256(append([]byte(m.participantID), tx.GetTransaction()...))
		result[i] = &protov30.SignedTopologyTransaction{
			Transaction: append(tx.GetTransaction(), sigBytes[:]...),
			Signatures: append(tx.GetSignatures(), &cryptov30.Signature{
				SignedBy:  m.participantID,
				Signature: sigBytes[:],
			}),
			Proposal: tx.Proposal,
		}
	}

	return result, nil
}

func (m *mockCantonClient) AddTransactions(_ context.Context, _ []*protov30.SignedTopologyTransaction, _ string) error {
	if m.state != nil {
		m.state.onDNSSubmit()
	}

	return nil
}

func (m *mockCantonClient) DNSExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

func (m *mockCantonClient) NSDExists(_ context.Context, _ string, _ string) (bool, error) {
	if m.state != nil {
		return m.state.hasNSDProposed(), nil
	}

	return true, nil
}

func (m *mockCantonClient) P2PExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// GetDNS returns the current DNS state. Before DNS submission it returns the
// initial 2-owner state with the target's fingerprint; after SubmitKickDNSOp
// (AddTransactions) it returns the post-rotation state with the new fingerprint.
func (m *mockCantonClient) GetDNS(_ context.Context, namespace string, _ string) (*client.DNSState, error) {
	if m.state != nil && m.state.dnsSubmitted() {
		return &client.DNSState{
			DecentralizedNamespace: namespace,
			Owners:                 []string{testNewNamespaceFP, "fp-owner-p2"},
			Threshold:              2,
			Serial:                 2,
		}, nil
	}

	return &client.DNSState{
		DecentralizedNamespace: namespace,
		Owners:                 []string{testTargetNSFP, "fp-owner-p2"},
		Threshold:              2,
		Serial:                 1,
	}, nil
}

// GetP2P returns the current P2P state including party signing keys.
// After at least one P2P proposal (PartyToParticipant Authorize) the signing
// keys are updated: the target's old DAML key is replaced so the poll can pass.
func (m *mockCantonClient) GetP2P(_ context.Context, partyUID string, _ string) (*client.P2PState, error) {
	signingKeys := &client.P2PSigningKeysInfo{
		Keys:      []string{testSigningKeyTargetB64, testSigningKeyOtherB64},
		Threshold: 1,
	}
	if m.state != nil && m.state.p2pReady() {
		signingKeys = &client.P2PSigningKeysInfo{
			Keys:      []string{testNewDamlKeyB64, testSigningKeyOtherB64},
			Threshold: 1,
		}
	}

	return &client.P2PState{
		Party: partyUID,
		Participants: []client.P2PParticipantInfo{
			{ParticipantUID: "p1", Permission: "CONFIRMATION"},
			{ParticipantUID: "p2", Permission: "CONFIRMATION"},
		},
		Threshold:        2,
		Serial:           1,
		PartySigningKeys: signingKeys,
	}, nil
}

func (m *mockCantonClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{
		{
			DecentralizedNamespace: "1220aabbccdd",
			Owners:                 []string{testTargetNSFP, "fp-owner-p2"},
			Threshold:              2,
			Serial:                 1,
		},
	}, nil
}

func (m *mockCantonClient) UploadDar(_ context.Context, _ []byte) (string, error) {
	return "mock-package-id", nil
}
func (m *mockCantonClient) ExportAcs(_ context.Context, _ []string, _ string, _ int64) ([]byte, error) {
	return nil, nil
}
func (m *mockCantonClient) ImportAcs(_ context.Context, _ []byte, _ string) error    { return nil }
func (m *mockCantonClient) DisconnectSynchronizer(_ context.Context, _ string) error { return nil }
func (m *mockCantonClient) ReconnectSynchronizer(_ context.Context, _ string) error  { return nil }
func (m *mockCantonClient) ListConnectedSynchronizers(_ context.Context) ([]client.SynchronizerInfo, error) {
	return nil, nil
}
func (m *mockCantonClient) ClearPartyOnboardingFlag(_ context.Context, _ string, _ string, _ int64) (bool, error) {
	return true, nil
}
func (m *mockCantonClient) LookupOffsetByTime(_ context.Context, _ int64) (int64, error) {
	return 0, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

// newDeps creates a CantonDeps without shared state — suitable for single-run
// tests that are expected to fail before any polling loop is reached.
func newDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID},
		Logger: logger.Nop(),
	}
}

// newDepsWithState creates a CantonDeps whose mock client shares the given
// mockState with all other clients in the same scenario. Required for
// multi-actor tests that reach polling loops.
func newDepsWithState(participantID string, state *mockState) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, state: state},
		Logger: logger.Nop(),
	}
}

func newKMSDeps(participantID string, kms client.KMSConfig, registrations *int) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, kmsRegistrations: registrations},
		KMS:    kms,
		Logger: logger.Nop(),
	}
}

// baseNamespaceInput returns a valid namespace-only key rotation input.
func baseNamespaceInput() keyrotation.KeyRotationInput {
	return keyrotation.KeyRotationInput{
		DecentralizedPartyID:       testPartyID,
		TargetParticipantID:        testTargetUID,
		TargetNamespaceFingerprint: testTargetNSFP,
		SynchronizerID:             testSyncID,
		RotateNamespaceKey:         true,
		RotateDamlKey:              false,
	}
}

// baseDamlInput returns a valid DAML-only key rotation input.
func baseDamlInput() keyrotation.KeyRotationInput {
	return keyrotation.KeyRotationInput{
		DecentralizedPartyID: testPartyID,
		TargetParticipantID:  testTargetUID,
		SynchronizerID:       testSyncID,
		RotateNamespaceKey:   false,
		RotateDamlKey:        true,
	}
}

// baseBothInput returns a valid key rotation input that rotates both keys.
func baseBothInput() keyrotation.KeyRotationInput {
	return keyrotation.KeyRotationInput{
		DecentralizedPartyID:       testPartyID,
		TargetParticipantID:        testTargetUID,
		TargetNamespaceFingerprint: testTargetNSFP,
		SynchronizerID:             testSyncID,
		RotateNamespaceKey:         true,
		RotateDamlKey:              true,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// TestKeyRotationSequence_NamespaceOnly_HappyPath validates the full multi-actor
// happy path for namespace key rotation only.
//
//	Run 1 (p1, target): reads state → generates key → NSD → DNS proposal → signs DNS (1/2) → ErrThresholdNotMet.
//	Run 2 (p2): all cached → p2 signs DNS (2/2) → DNS submitted → DNS confirmed → SUCCESS (no P2P step).
func TestKeyRotationSequence_NamespaceOnly_HappyPath(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseNamespaceInput()

	// Run 1 (p1, target): generates key, NSD, creates DNS proposal, signs (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (1/2)")

	// Run 2 (p2): signs DNS (2/2) → DNS submitted → DNS confirmed → SUCCESS.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.NoError(t, err, "run 2: ceremony should complete successfully")
	assert.True(t, sr.Output.NamespaceKeyRotated, "NamespaceKeyRotated should be true")
	assert.False(t, sr.Output.DamlKeyRotated, "DamlKeyRotated should be false")
	assert.True(t, sr.Output.DNSUpdated, "DNSUpdated should be true")
	assert.False(t, sr.Output.P2PUpdated, "P2PUpdated should be false")
	assert.NotEmpty(t, sr.Output.NewNamespaceFingerprint, "NewNamespaceFingerprint should be set")
	assert.Empty(t, sr.Output.NewDamlKeyFingerprint, "NewDamlKeyFingerprint should be empty")
}

// TestKeyRotationSequence_DamlOnly_HappyPath validates the full multi-actor
// happy path for DAML key rotation only.
//
//	Run 1 (p1, target): reads state → generates key (discovers old) → proposes P2P (1/2) → ErrThresholdNotMet.
//	Run 2 (p2): cached → p2 proposes P2P (2/2) → P2P confirmed → SUCCESS (no DNS step).
func TestKeyRotationSequence_DamlOnly_HappyPath(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseDamlInput()

	// Run 1 (p1, target): generates key, discovers old DAML key, proposes P2P (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: P2P threshold not met (1/2)")

	// Run 2 (p2): proposes P2P (2/2) → P2P confirmed → SUCCESS.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.NoError(t, err, "run 2: ceremony should complete successfully")
	assert.False(t, sr.Output.NamespaceKeyRotated, "NamespaceKeyRotated should be false")
	assert.True(t, sr.Output.DamlKeyRotated, "DamlKeyRotated should be true")
	assert.False(t, sr.Output.DNSUpdated, "DNSUpdated should be false")
	assert.True(t, sr.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Empty(t, sr.Output.NewNamespaceFingerprint, "NewNamespaceFingerprint should be empty")
	assert.NotEmpty(t, sr.Output.NewDamlKeyFingerprint, "NewDamlKeyFingerprint should be set")
}

// TestKeyRotationSequence_Both_HappyPath validates the full multi-actor happy
// path when both namespace and DAML keys are rotated.
//
//	Run 1 (p1, target): reads state → generates both keys → NSD → DNS proposal → signs DNS (1/2) → ErrThresholdNotMet.
//	Run 2 (p2): signs DNS (2/2) → DNS submitted → proposes P2P (1/2) → ErrThresholdNotMet.
//	Run 3 (p1): proposes P2P (2/2) → P2P confirmed → SUCCESS.
func TestKeyRotationSequence_Both_HappyPath(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseBothInput()

	// Run 1 (p1, target): generates both keys, NSD, DNS proposal, signs DNS (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (1/2)")

	// Run 2 (p2): signs DNS (2/2) → submit → DNS confirmed → proposes P2P (1/2) → threshold not met.
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 2: P2P threshold not met (1/2)")

	// Run 3 (p1): proposes P2P (2/2) → P2P confirmed → SUCCESS.
	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.NoError(t, err, "run 3: ceremony should complete successfully")
	assert.True(t, sr.Output.NamespaceKeyRotated, "NamespaceKeyRotated should be true")
	assert.True(t, sr.Output.DamlKeyRotated, "DamlKeyRotated should be true")
	assert.True(t, sr.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, sr.Output.P2PUpdated, "P2PUpdated should be true")
	assert.NotEmpty(t, sr.Output.NewNamespaceFingerprint, "NewNamespaceFingerprint should be set")
	assert.NotEmpty(t, sr.Output.NewDamlKeyFingerprint, "NewDamlKeyFingerprint should be set")
}

// TestKeyRotationSequence_Idempotent verifies that re-running the ceremony
// after completion returns the same cached output.
func TestKeyRotationSequence_Idempotent(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseNamespaceInput()

	// Complete the ceremony (namespace-only: p1 + p2).
	_, _ = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	sr1, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.NoError(t, err)

	// Re-run should return cached result.
	sr2, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.NoError(t, err)
	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report")
	assert.Equal(t, sr1.Output, sr2.Output)
}

// TestKeyRotationSequence_ThresholdNotMet_DNS verifies that the sequence
// returns ErrThresholdNotMet when only one of two required DNS signatures has
// been collected.
func TestKeyRotationSequence_ThresholdNotMet_DNS(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseNamespaceInput()

	// p1 (target) signs DNS (1/2, currentThreshold=2) → gate fires → ErrThresholdNotMet.
	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "DNS signatures collected")
}

// TestKeyRotationSequence_ThresholdNotMet_P2P verifies that after key
// generation, the sequence returns ErrThresholdNotMet when P2P proposals
// are still pending.
func TestKeyRotationSequence_ThresholdNotMet_P2P(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseDamlInput()

	// p1 (target): generates key, proposes P2P (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "P2P proposals collected")
}

// TestKeyRotationSequence_ThresholdNotMet_TargetNotRun verifies that when a
// non-target participant runs first (before the target generates keys), the
// sequence returns ErrThresholdNotMet because key generation requires the
// target participant.
func TestKeyRotationSequence_ThresholdNotMet_TargetNotRun(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	// p2 (non-target) runs first — target hasn't generated rotated keys yet.
	_, err := operations.ExecuteSequence(b, keyrotation.KeyRotationSequence, newDeps("p2"), baseDamlInput())
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "target participant has not generated rotated keys yet")
}

// TestKeyRotationSequence_ResumeAfterPartialSigning exercises the full async
// multi-actor flow with the shared reporter, verifying each resume step for
// both namespace and DAML key rotation.
//
//	Run 1 (p1, target): reads state → gen keys → NSD → DNS proposal → signs DNS (1/2) → ErrThresholdNotMet.
//	Run 2 (p2): signs DNS (2/2) → submit → DNS confirmed → proposes P2P (1/2) → ErrThresholdNotMet.
//	Run 3 (p1): proposes P2P (2/2) → P2P confirmed → SUCCESS.
func TestKeyRotationSequence_ResumeAfterPartialSigning(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseBothInput()

	_, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: DNS signing gate (1/2)")

	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 2: P2P gate (1/2)")

	sr, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, newDepsWithState(testTargetUID, state), input)
	require.NoError(t, err, "run 3: ceremony complete")
	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
}

// TestKeyRotationSequence_InvalidPartyID verifies that a malformed party ID
// returns an unrecoverable error before any operations execute.
func TestKeyRotationSequence_InvalidPartyID(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseNamespaceInput()
	input.DecentralizedPartyID = "no-double-colon"
	_, err := operations.ExecuteSequence(b, keyrotation.KeyRotationSequence, newDeps(testTargetUID), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decentralized_party_id")
}

// TestKeyRotationSequence_NeitherKeySelected verifies that setting both
// rotation flags to false returns an unrecoverable error.
func TestKeyRotationSequence_NeitherKeySelected(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := keyrotation.KeyRotationInput{
		DecentralizedPartyID: testPartyID,
		TargetParticipantID:  testTargetUID,
		SynchronizerID:       testSyncID,
		RotateNamespaceKey:   false,
		RotateDamlKey:        false,
	}
	_, err := operations.ExecuteSequence(b, keyrotation.KeyRotationSequence, newDeps(testTargetUID), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one of rotate_namespace_key or rotate_daml_key must be true")
}

// TestKeyRotationSequence_TargetNotInDNS verifies that when the target's
// namespace fingerprint is not found in the current DNS owners, the sequence
// returns an unrecoverable error.
func TestKeyRotationSequence_TargetNotInDNS(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseNamespaceInput()
	input.TargetNamespaceFingerprint = "does-not-exist"
	_, err := operations.ExecuteSequence(b, keyrotation.KeyRotationSequence, newDeps(testTargetUID), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in DNS owners")
}

// TestKeyRotationSequence_MissingTargetParticipantID verifies that an empty
// target participant ID returns an unrecoverable error.
func TestKeyRotationSequence_MissingTargetParticipantID(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseDamlInput()
	input.TargetParticipantID = ""
	_, err := operations.ExecuteSequence(b, keyrotation.KeyRotationSequence, newDeps(testTargetUID), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target_participant_id is required")
}

// TestKeyRotationSequence_SerializationValid checks that all IN/OUT types can
// be round-tripped through the JSON serializer used by the operations framework.
func TestKeyRotationSequence_SerializationValid(t *testing.T) {
	t.Parallel()
	lggr := logger.Nop()
	require.True(t, operations.IsSerializable(lggr, keyrotation.KeyRotationInput{}), "KeyRotationInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, keyrotation.KeyRotationOutput{}), "KeyRotationOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, keys.GenerateRotatedKeyInput{}), "GenerateRotatedKeyInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, keys.GenerateRotatedKeyOutput{}), "GenerateRotatedKeyOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, topology.CreateRotationDNSProposalInput{}), "CreateRotationDNSProposalInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, topology.CreateRotationDNSProposalOutput{}), "CreateRotationDNSProposalOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, topology.ProposeRotationP2PInput{}), "ProposeRotationP2PInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, topology.ProposeRotationP2POutput{}), "ProposeRotationP2POutput must be serializable")
}

// TestGRPCCantonClientImplementsInterface confirms that *client.GRPCCantonClient
// satisfies client.CantonClient.
func TestGRPCCantonClientImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ client.CantonClient = (*client.GRPCCantonClient)(nil)
}

// TestMockCantonClientImplementsInterface confirms the test mock satisfies
// the full client.CantonClient interface.
func TestMockCantonClientImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ client.CantonClient = (*mockCantonClient)(nil)
}

// TestGenerateRotatedKeyOp_ParticipantMismatch verifies that trying to
// generate rotated keys from the wrong participant returns an error.
func TestGenerateRotatedKeyOp_ParticipantMismatch(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	_, err := operations.ExecuteOperation(b, keys.GenerateRotatedKeyOp, newDeps("p1"), keys.GenerateRotatedKeyInput{
		ParticipantID:      "p2", // mismatch: client is p1 but input says p2
		SynchronizerID:     testSyncID,
		DNSOwners:          []string{testTargetNSFP, "fp-owner-p2"},
		RotateNamespaceKey: true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "participant ID mismatch")
}

// TestGenerateRotatedKeyOp_Success verifies key generation for both key types
// produces valid output.
func TestGenerateRotatedKeyOp_Success(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	r, err := operations.ExecuteOperation(b, keys.GenerateRotatedKeyOp, newDeps(testTargetUID), keys.GenerateRotatedKeyInput{
		ParticipantID:       testTargetUID,
		SynchronizerID:      testSyncID,
		DNSOwners:           []string{testTargetNSFP, "fp-owner-p2"},
		RotateNamespaceKey:  true,
		RotateDamlKey:       true,
		KnownSigningKeysB64: []string{testSigningKeyTargetB64, testSigningKeyOtherB64},
	})
	require.NoError(t, err)
	assert.Equal(t, testTargetUID, r.Output.ParticipantID)
	assert.NotEmpty(t, r.Output.NewNamespaceKeyB64)
	assert.NotEmpty(t, r.Output.NewNamespaceFingerprint)
	assert.NotEmpty(t, r.Output.NewDamlKeyB64)
	assert.NotEmpty(t, r.Output.NewDamlKeyFingerprint)
	assert.NotEmpty(t, r.Output.OldDamlKeyB64)
	assert.Equal(t, testSigningKeyTargetB64, r.Output.OldDamlKeyB64)

	// Verify the NewNamespaceKeyB64 is valid proto.
	keyBytes, err := base64.StdEncoding.DecodeString(r.Output.NewNamespaceKeyB64)
	require.NoError(t, err)
	var pk cryptov30.SigningPublicKey
	require.NoError(t, proto.Unmarshal(keyBytes, &pk))
	assert.NotEmpty(t, pk.GetPublicKey())
}

func TestGenerateRotatedKeyOp_KMSNamespaceOnly(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	var registrations int

	r, err := operations.ExecuteOperation(b, keys.GenerateRotatedKeyOp, newKMSDeps(testTargetUID, client.KMSConfig{
		NamespaceKeyID: "arn:aws:kms:us-east-1:123456789:key/rotated-namespace",
	}, &registrations), keys.GenerateRotatedKeyInput{
		ParticipantID:      testTargetUID,
		SynchronizerID:     testSyncID,
		DNSOwners:          []string{testTargetNSFP, "fp-owner-p2"},
		RotateNamespaceKey: true,
	})
	require.NoError(t, err)

	assert.NotEmpty(t, r.Output.NewNamespaceKeyB64)
	assert.NotEmpty(t, r.Output.NewNamespaceFingerprint)
	assert.Empty(t, r.Output.NewDamlKeyB64)
	assert.Equal(t, 1, registrations)
}

func TestGenerateRotatedKeyOp_KMSProtocolOnly(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	var registrations int

	r, err := operations.ExecuteOperation(b, keys.GenerateRotatedKeyOp, newKMSDeps(testTargetUID, client.KMSConfig{
		ProtocolKeyID: "arn:aws:kms:us-east-1:123456789:key/rotated-protocol",
	}, &registrations), keys.GenerateRotatedKeyInput{
		ParticipantID:       testTargetUID,
		SynchronizerID:      testSyncID,
		DNSOwners:           []string{testTargetNSFP, "fp-owner-p2"},
		RotateDamlKey:       true,
		KnownSigningKeysB64: []string{testSigningKeyTargetB64, testSigningKeyOtherB64},
	})
	require.NoError(t, err)

	assert.Empty(t, r.Output.NewNamespaceKeyB64)
	assert.NotEmpty(t, r.Output.NewDamlKeyB64)
	assert.Equal(t, testSigningKeyTargetB64, r.Output.OldDamlKeyB64)
	assert.Equal(t, 1, registrations)
}

// TestProposeRotationP2POp_ParticipantMismatch verifies that proposing a
// P2P rotation from the wrong participant is rejected.
func TestProposeRotationP2POp_ParticipantMismatch(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	_, err := operations.ExecuteOperation(b, topology.ProposeRotationP2POp, newDeps("p1"), topology.ProposeRotationP2PInput{
		ParticipantID:         "p2", // mismatch
		PartyID:               testPartyID,
		AllParticipantUIDs:    []string{"p1", "p2"},
		CurrentSigningKeysB64: []string{testSigningKeyTargetB64, testSigningKeyOtherB64},
		OldDamlKeyB64:         testSigningKeyTargetB64,
		NewDamlKeyB64:         makeTestSigningKeyB64("new-daml"),
		SynchronizerID:        testSyncID,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "participant ID mismatch")
}
