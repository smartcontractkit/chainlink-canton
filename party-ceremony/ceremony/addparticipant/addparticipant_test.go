package addparticipant_test

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
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/addparticipant"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/topology"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// ── Constants shared across tests ────────────────────────────────────────────

const (
	testPartyID = "test-party::1220aabbccdd"
	testSyncID  = "global"
	testNewUID  = "p3"
	testNewNSFP = "fp-owner-new"
	testNSName  = "test-add-namespace"
)

// ── Mock shared state ────────────────────────────────────────────────────────

// mockState tracks topology state changes across all mockCantonClient instances
// within a single test scenario.
type mockState struct {
	mu               sync.Mutex
	keyGenerated     bool // true after new participant generates keys
	kmsRegistrations int
	nsdProposed      bool // true after NSD is proposed
	addCallCount     int  // incremented by each AddTransactions call
	authorizePostDNS int  // incremented by Authorize calls after DNS is submitted
}

func (s *mockState) onKeyGenerated() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyGenerated = true
}

func (s *mockState) onKMSRegistration() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kmsRegistrations++
}

func (s *mockState) kmsRegistrationCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.kmsRegistrations
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

func (s *mockState) onAddTransactions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addCallCount++
}

func (s *mockState) dnsSubmitted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.addCallCount > 0
}

func (s *mockState) onP2PAuthorize() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.authorizePostDNS++
}

func (s *mockState) p2pSubmitted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.authorizePostDNS > 0
}

// ── Mock Canton client ───────────────────────────────────────────────────────

type mockCantonClient struct {
	participantID string
	state         *mockState
}

func (m *mockCantonClient) GetParticipantUID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockCantonClient) GetParticipantID(_ context.Context) (string, error) {
	return m.participantID, nil
}

func (m *mockCantonClient) GenerateSigningKey(_ context.Context, name string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	if m.state != nil {
		m.state.onKeyGenerated()
	}
	raw := sha256.Sum256([]byte(m.participantID + ":" + name))

	return &cryptov30.SigningPublicKey{
		PublicKey: raw[:],
		KeySpec:   cryptov30.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
	}, nil
}

func (m *mockCantonClient) RegisterKmsSigningKey(_ context.Context, kmsKeyID string, name string, _ []cryptov30.SigningKeyUsage) (*cryptov30.SigningPublicKey, error) {
	if m.state != nil {
		m.state.onKMSRegistration()
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

func (m *mockCantonClient) GetNamespaceKeyName(_ context.Context, _ string, _ []string) (string, error) {
	return "mock-ns-key", nil
}

func (m *mockCantonClient) Authorize(_ context.Context, serial uint32, mapping *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
	// Track NSD proposals from the new participant.
	if m.state != nil {
		if _, ok := mapping.GetMapping().(*protov30.TopologyMapping_NamespaceDelegation); ok {
			m.state.onNSDProposed()
		}
		if m.state.dnsSubmitted() {
			m.state.onP2PAuthorize()
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
		m.state.onAddTransactions()
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

// GetDNS returns the current DNS state. Before DNS submission it returns a
// 2-owner state; after SubmitAddDNSOp it returns a 3-owner state.
func (m *mockCantonClient) GetDNS(_ context.Context, namespace string, _ string) (*client.DNSState, error) {
	if m.state != nil && m.state.dnsSubmitted() {
		return &client.DNSState{
			DecentralizedNamespace: namespace,
			Owners:                 []string{"fp-owner-p1", "fp-owner-p2", testNewNSFP},
			Threshold:              2,
			Serial:                 2,
		}, nil
	}

	return &client.DNSState{
		DecentralizedNamespace: namespace,
		Owners:                 []string{"fp-owner-p1", "fp-owner-p2"},
		Threshold:              2,
		Serial:                 1,
	}, nil
}

// GetP2P returns the current P2P state. Before any post-DNS P2P Authorize calls
// it returns a 2-participant state; after at least one P2P proposal it returns
// a 3-participant state.
func (m *mockCantonClient) GetP2P(_ context.Context, partyUID string, _ string) (*client.P2PState, error) {
	if m.state != nil && m.state.p2pSubmitted() {
		return &client.P2PState{
			Party: partyUID,
			Participants: []client.P2PParticipantInfo{
				{ParticipantUID: "p1", Permission: "CONFIRMATION"},
				{ParticipantUID: "p2", Permission: "CONFIRMATION"},
				{ParticipantUID: testNewUID, Permission: "CONFIRMATION"},
			},
			Threshold: 2,
			Serial:    2,
		}, nil
	}

	return &client.P2PState{
		Party: partyUID,
		Participants: []client.P2PParticipantInfo{
			{ParticipantUID: "p1", Permission: "CONFIRMATION"},
			{ParticipantUID: "p2", Permission: "CONFIRMATION"},
		},
		Threshold: 2,
		Serial:    1,
	}, nil
}

func (m *mockCantonClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{
		{
			DecentralizedNamespace: "1220aabbccdd",
			Owners:                 []string{"fp-owner-p1", "fp-owner-p2"},
			Threshold:              2,
			Serial:                 1,
		},
	}, nil
}

func (m *mockCantonClient) GetProtocolKeyFingerprint(_ context.Context, _ []string) (string, string, error) {
	return "mock-protocol-fp", "mock-protocol-key-b64", nil
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

func newDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID},
		Logger: logger.Nop(),
	}
}

func newDepsWithState(participantID string, state *mockState) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, state: state},
		Logger: logger.Nop(),
	}
}

func newKMSDepsWithState(participantID string, state *mockState) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, state: state},
		KMS: client.KMSConfig{
			NamespaceKeyID: "arn:aws:kms:us-east-1:123456789:key/" + participantID + "-namespace",
			ProtocolKeyID:  "arn:aws:kms:us-east-1:123456789:key/" + participantID + "-protocol",
		},
		Logger: logger.Nop(),
	}
}

func baseInput() addparticipant.AddParticipantInput {
	return addparticipant.AddParticipantInput{
		DecentralizedPartyID: testPartyID,
		NewParticipantID:     testNewUID,
		NamespaceName:        testNSName,
		SynchronizerID:       testSyncID,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// TestAddParticipantSequence_HappyPath validates the full multi-actor happy path.
//
//	Run 1 (p3 — new): key-gen → NSD → reads state → DNS proposal → ErrThresholdNotMet (0/2 DNS sigs).
//	Run 2 (p1): signs DNS (1/2) → ErrThresholdNotMet.
//	Run 3 (p2): signs DNS (2/2) → DNS submitted → P2P (1/2 existing) → ErrThresholdNotMet.
//	Run 4 (p1): P2P (2/2 existing) → new participant consent pending → ErrThresholdNotMet.
//	Run 5 (p3 — new): consents to P2P hosting → P2P confirmed → SUCCESS.
func TestAddParticipantSequence_HappyPath(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	// Run 1 (p3 — new): key-gen, NSD, reads state, DNS proposal → 0/2 DNS sigs.
	_, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (0/2)")

	// Run 2 (p1): signs DNS (1/2) → threshold not met.
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 2: DNS threshold not met (1/2)")

	// Run 3 (p2): signs DNS (2/2) → submits → P2P from p2 (1/2 existing) → threshold not met.
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 3: P2P existing threshold not met (1/2)")

	// Run 4 (p1): P2P from p1 (2/2 existing) → new participant consent pending.
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 4: new participant consent pending")

	// Run 5 (p3 — new): consents to P2P hosting → confirmed → SUCCESS.
	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err, "run 5: ceremony should complete successfully")

	assert.True(t, sr.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, sr.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Equal(t, 2, sr.Output.NewThreshold, "NewThreshold should default to current (2)")
	assert.Len(t, sr.Output.AllOwners, 3, "should have 3 owners after add")
}

func TestAddParticipantSequence_KMS_NewParticipantUsesLocalConfig(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	_, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newKMSDepsWithState(testNewUID, state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (0/2)")

	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 2: DNS threshold not met (1/2)")

	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 3: P2P existing threshold not met (1/2)")

	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 4: new participant consent pending")

	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newKMSDepsWithState(testNewUID, state), input)
	require.NoError(t, err, "run 5: ceremony should complete successfully")

	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
	assert.Equal(t, 2, state.kmsRegistrationCount(), "namespace and protocol keys should both be registered through KMS")
}

// TestAddParticipantSequence_Idempotent verifies that re-running after
// completion returns the same cached output.
func TestAddParticipantSequence_Idempotent(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	// Complete the ceremony (5 runs: p3, p1, p2, p1, p3).
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	sr1, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err)

	// Re-run should return cached result.
	sr2, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err)

	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report")
	assert.Equal(t, sr1.Output, sr2.Output)
}

// TestAddParticipantSequence_ThresholdNotMet_KeyGen verifies that the sequence
// returns ErrThresholdNotMet when the new participant hasn't generated keys.
func TestAddParticipantSequence_ThresholdNotMet_KeyGen(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	// p1 (existing) runs first — new participant hasn't generated keys yet.
	_, err := operations.ExecuteSequence(b, addparticipant.AddParticipantSequence, newDeps("p1"), baseInput())
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "new participant has not generated keys yet")
}

// TestAddParticipantSequence_ThresholdNotMet_DNSSigning verifies the sequence
// returns ErrThresholdNotMet when only 1 of 2 required DNS signatures collected.
func TestAddParticipantSequence_ThresholdNotMet_DNSSigning(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	// New participant generates keys and reads state.
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)

	// p1 signs DNS (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "DNS signatures collected")
}

// TestAddParticipantSequence_ThresholdNotMet_P2PProposal verifies that after
// DNS is confirmed, the sequence returns ErrThresholdNotMet when P2P proposals
// are still pending.
func TestAddParticipantSequence_ThresholdNotMet_P2PProposal(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	// New participant.
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	// p1 signs DNS (1/2).
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	// p2 signs DNS (2/2) → submit → P2P (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error())
	require.ErrorContains(t, err, "P2P proposals collected")
}

// TestAddParticipantSequence_NewThresholdOverride verifies an explicit
// --new-threshold override is propagated to the output.
func TestAddParticipantSequence_NewThresholdOverride(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	input.NewThreshold = 3 // override: propagates to ProposeAddP2POp.NewP2PThreshold

	// Same 5-run flow as HappyPath. Gate still uses currentState.DNSThreshold=2
	// for existing participants regardless of the override value.
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err)
	assert.Equal(t, 3, sr.Output.NewThreshold, "NewThreshold should match override")
}

// TestAddParticipantSequence_InvalidPartyID verifies that a malformed party ID
// returns an unrecoverable error before any operations execute.
func TestAddParticipantSequence_InvalidPartyID(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseInput()
	input.DecentralizedPartyID = "no-double-colon"
	_, err := operations.ExecuteSequence(b, addparticipant.AddParticipantSequence, newDeps(testNewUID), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decentralized_party_id")
}

// TestAddParticipantSequence_ResumeAfterPartialSigning exercises the full async
// multi-actor flow with a shared reporter, verifying each resume step.
func TestAddParticipantSequence_ResumeAfterPartialSigning(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()

	// Run 1 (new participant): generates keys, NSD, reads state → 0/2 DNS sigs.
	_, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 1: DNS signing gate")

	// Run 2 (p1): signs DNS (1/2).
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 2: DNS signing gate (1/2)")

	// Run 3 (p2): signs DNS (2/2) → submit → P2P from p2 (1/2 existing).
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 3: P2P existing gate (1/2)")

	// Run 4 (p1): P2P from p1 (2/2 existing) → new participant consent pending.
	_, err = operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, addparticipant.ErrThresholdNotMet.Error(), "run 4: new participant consent pending")

	// Run 5 (p3 — new): consents to P2P hosting → confirmed → success.
	sr, err := operations.ExecuteSequence(newBundle(), addparticipant.AddParticipantSequence, newDepsWithState(testNewUID, state), input)
	require.NoError(t, err, "run 5: ceremony complete")
	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
}

// TestGenerateNewMemberKeyOp_ParticipantMismatch verifies that trying to
// generate keys from the wrong participant returns an error.
func TestGenerateNewMemberKeyOp_ParticipantMismatch(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	_, err := operations.ExecuteOperation(b, keys.CreateMemberKeyOp, newDeps("p1"), keys.CreateMemberKeyInput{
		NamespaceName: testNSName,
		ParticipantID: "p2", // mismatch: client is p1 but input says p2
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "participant ID mismatch")
}

// TestGenerateNewMemberKeyOp_Success verifies key generation produces valid output.
func TestGenerateNewMemberKeyOp_Success(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	r, err := operations.ExecuteOperation(b, keys.CreateMemberKeyOp, newDeps(testNewUID), keys.CreateMemberKeyInput{
		NamespaceName: testNSName,
		ParticipantID: testNewUID,
	})
	require.NoError(t, err)
	assert.Equal(t, testNewUID, r.Output.ParticipantID)
	assert.NotEmpty(t, r.Output.NamespaceFingerprint)
	assert.NotEmpty(t, r.Output.SigningKeyB64)
	assert.NotEmpty(t, r.Output.DamlKeyB64)
	assert.NotEmpty(t, r.Output.DamlKeyFingerprint)

	// Verify the SigningKeyB64 is valid proto.
	keyBytes, err := base64.StdEncoding.DecodeString(r.Output.SigningKeyB64)
	require.NoError(t, err)
	var pk cryptov30.SigningPublicKey
	require.NoError(t, proto.Unmarshal(keyBytes, &pk))
	assert.NotEmpty(t, pk.GetPublicKey())
}

// TestReadCurrentStateOp_InvalidPartyID verifies that a malformed party ID
// returns an appropriate error.
func TestReadCurrentStateOp_InvalidPartyID(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	_, err := operations.ExecuteOperation(b, topology.ReadCurrentStateOp, newDeps("p1"), topology.ReadCurrentStateInput{
		DecentralizedPartyID: "no-separator",
		SynchronizerID:       testSyncID,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decentralized_party_id")
}

// TestSignAddDNSProposalOp_ParticipantMismatch verifies that signing from the
// wrong participant is rejected.
func TestSignAddDNSProposalOp_ParticipantMismatch(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)

	// Create a minimal valid proposal for the test.
	tx := &protov30.SignedTopologyTransaction{
		Transaction: []byte("test-tx"),
		Proposal:    true,
	}
	txBytes, _ := proto.Marshal(tx)
	txB64 := base64.StdEncoding.EncodeToString(txBytes)
	hash := sha256.Sum256(txBytes)

	_, err := operations.ExecuteOperation(b, topology.SignDNSProposalOp, newDeps("p1"), topology.SignDNSProposalInput{
		ParticipantID:      "p2", // mismatch
		ProposalHashSHA256: fmt.Sprintf("%x", hash),
		DNSTxB64:           txB64,
		SynchronizerID:     testSyncID,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "participant ID mismatch")
}
