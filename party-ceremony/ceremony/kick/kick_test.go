package kick_test

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

	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/ceremony/kick"
	"github.com/chainlink/canton-party-ceremony/internal/client"
)

// ── Constants shared across tests ────────────────────────────────────────────

const (
	testPartyID   = "test-party::1220aabbccdd"
	testSyncID    = "global"
	testKickedUID = "p3"
	testKickedFP  = "fp-owner-kicked"
)

// ── Mock shared state ────────────────────────────────────────────────────────

// mockState is shared across all mockCantonClient instances within a single
// test scenario. It tracks topology state changes (DNS submission, P2P
// proposals) that Canton would normally record.
type mockState struct {
	mu               sync.Mutex
	addCallCount     int // incremented by each AddTransactions call
	authorizePostDNS int // incremented by Authorize calls after DNS is submitted
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

// mockCantonClient is a deterministic, in-memory implementation of
// [client.CantonClient] for kick ceremony unit tests.
type mockCantonClient struct {
	participantID string
	state         *mockState // nil for single-run tests that never reach polling
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

func (m *mockCantonClient) GetNamespaceFingerprint(_ context.Context, keyName string, _ string, _ []string) (string, error) {
	raw := sha256.Sum256([]byte(m.participantID + ":" + keyName))
	return fmt.Sprintf("1220%x", raw[:8]), nil
}

func (m *mockCantonClient) Authorize(_ context.Context, serial uint32, mapping *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
	// Track P2P Authorize calls (post-DNS) so GetP2P can return updated state.
	if m.state != nil && m.state.dnsSubmitted() {
		m.state.onP2PAuthorize()
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
	return true, nil
}

func (m *mockCantonClient) P2PExists(_ context.Context, _ string, _ string) (bool, error) {
	return true, nil
}

// GetDNS returns the current DNS state. Before DNS submission it returns the
// initial 3-owner state; after SubmitKickDNSOp (AddTransactions) it returns
// the post-kick 2-owner state so the polling retry loop can complete.
func (m *mockCantonClient) GetDNS(_ context.Context, namespace string, _ string) (*client.DNSState, error) {
	if m.state != nil && m.state.dnsSubmitted() {
		return &client.DNSState{
			DecentralizedNamespace: namespace,
			Owners:                 []string{"fp-owner-p1", "fp-owner-p2"},
			Threshold:              2,
			Serial:                 2,
		}, nil
	}

	return &client.DNSState{
		DecentralizedNamespace: namespace,
		Owners:                 []string{"fp-owner-p1", "fp-owner-p2", testKickedFP},
		Threshold:              2,
		Serial:                 1,
	}, nil
}

// GetP2P returns the current P2P state. Before any post-DNS P2P Authorize calls
// it returns the initial 3-participant state; once at least one P2P proposal
// has been submitted it returns the post-kick 2-participant state so the P2P
// confirmation poll can complete.
func (m *mockCantonClient) GetP2P(_ context.Context, partyUID string, _ string) (*client.P2PState, error) {
	if m.state != nil && m.state.p2pSubmitted() {
		return &client.P2PState{
			Party: partyUID,
			Participants: []client.P2PParticipantInfo{
				{ParticipantUID: "p1", Permission: "CONFIRMATION"},
				{ParticipantUID: "p2", Permission: "CONFIRMATION"},
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
			{ParticipantUID: testKickedUID, Permission: "CONFIRMATION"},
		},
		Threshold: 2,
		Serial:    1,
	}, nil
}

func (m *mockCantonClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{
		{
			DecentralizedNamespace: "1220aabbccdd",
			Owners:                 []string{"fp-owner-p1", "fp-owner-p2", testKickedFP},
			Threshold:              2,
			Serial:                 1,
		},
	}, nil
}

func (m *mockCantonClient) UploadDar(_ context.Context, _ []byte) (string, error) {
	return "mock-package-id", nil
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
// multi-actor tests that reach SubmitKickDNSOp or ProposeKickP2POp so that
// the topology polling loops receive updated state and terminate.
func newDepsWithState(participantID string, state *mockState) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, state: state},
		Logger: logger.Nop(),
	}
}

// baseInput returns a valid 2-remaining kick input (kick "p3" from a 3-member party).
func baseInput() kick.KickInput {
	return kick.KickInput{
		DecentralizedPartyID:       testPartyID,
		KickedParticipantID:        testKickedUID,
		KickedNamespaceFingerprint: testKickedFP,
		RemainingParticipants:      []string{"p1", "p2"},
		SynchronizerID:             testSyncID,
	}
}

// ── Tests ─────────────────────────────────────────────────────────────────────

// TestKickSequence_HappyPath validates the full multi-actor happy path.
// Run 1 (p1): reads state → creates DNS proposal → signs DNS (1/2) → ErrThresholdNotMet.
// Run 2 (p2): all cached → p2 signs DNS (2/2) → DNS submitted → p2 proposes P2P (1/2) → ErrThresholdNotMet.
// Run 3 (p1): DNS cached → p1 proposes P2P (2/2) → P2P confirmed → SUCCESS.
func TestKickSequence_HappyPath(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	// Run 1 (p1): reads state, creates proposal, signs DNS (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (1/2)")
	// Run 2 (p2): signs DNS (2/2) → DNS submitted → confirms → proposes P2P (1/2) → threshold not met.
	_, err = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 2: P2P threshold not met (1/2)")
	// Run 3 (p1): P2P proposal (2/2) → P2P confirmed → success.
	sr, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err, "run 3: ceremony should complete successfully")
	assert.True(t, sr.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, sr.Output.P2PUpdated, "P2PUpdated should be true")
	assert.Equal(t, 2, sr.Output.NewThreshold, "NewThreshold should be 2 (floor(2/2)+1)")
	assert.NotContains(t, sr.Output.RemainingOwners, testKickedFP, "kicked fingerprint should be removed")
	assert.Len(t, sr.Output.RemainingOwners, 2, "should have 2 remaining owners")
}

// TestKickSequence_Idempotent verifies that re-running the ceremony after
// completion returns the same cached output.
func TestKickSequence_Idempotent(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	// Complete the ceremony (requires p1 + p2 + p1 runs).
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)
	// Re-run should return cached result.
	sr1, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err)
	sr2, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err)
	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report")
	assert.Equal(t, sr1.Output, sr2.Output)
}

// TestKickSequence_ThresholdNotMet_DNSSigning verifies that the sequence
// returns ErrThresholdNotMet when only one of two required signatures has
// been collected.
func TestKickSequence_ThresholdNotMet_DNSSigning(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	// p1 signs DNS (1/2, currentThreshold=2) → gate fires → ErrThresholdNotMet.
	_, err := operations.ExecuteSequence(b, kick.KickSequence, newDeps("p1"), baseInput())
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error())
}

// TestKickSequence_ThresholdNotMet_P2PProposal verifies that after DNS is
// confirmed, the sequence returns ErrThresholdNotMet when P2P proposals are
// still pending.
func TestKickSequence_ThresholdNotMet_P2PProposal(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	// Run 1 (p1): DNS (1/2) → threshold not met.
	_, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 1: DNS signing")
	// Run 2 (p2): DNS (2/2) + submit + P2P (1/2) → threshold not met.
	_, err = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 2: P2P proposals")
}

// TestKickSequence_ResumeAfterPartialSigning exercises the full async
// multi-actor flow with the shared reporter, verifying each resume step.
//
//	Run 1 (p1): reads state → creates DNS proposal → signs DNS (1/2) → ErrThresholdNotMet.
//	Run 2 (p2): all cached → p2 signs DNS (2/2) → DNS submitted → p2 proposes P2P (1/2) → ErrThresholdNotMet.
//	Run 3 (p1): DNS cached → p1 proposes P2P (2/2) → P2P confirmed → SUCCESS.
func TestKickSequence_ResumeAfterPartialSigning(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	_, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 1: DNS signing gate (1/2)")
	_, err = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)
	require.ErrorContains(t, err, kick.ErrThresholdNotMet.Error(), "run 2: P2P gate (1/2)")
	sr, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err, "run 3: ceremony complete")
	assert.True(t, sr.Output.DNSUpdated)
	assert.True(t, sr.Output.P2PUpdated)
	assert.Equal(t, 2, sr.Output.NewThreshold)
}

// TestKickSequence_NewThresholdOverride verifies the explicit --new-threshold
// override is respected.
func TestKickSequence_NewThresholdOverride(t *testing.T) {
	t.Parallel()
	sharedReporter := operations.NewMemoryReporter()
	state := &mockState{}
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := baseInput()
	input.NewThreshold = 1
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	_, _ = operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p2", state), input)
	sr, err := operations.ExecuteSequence(newBundle(), kick.KickSequence, newDepsWithState("p1", state), input)
	require.NoError(t, err)
	assert.Equal(t, 1, sr.Output.NewThreshold, "NewThreshold should match override")
}

// TestKickSequence_InvalidPartyID verifies that a malformed party ID
// returns an unrecoverable error before any operations execute.
func TestKickSequence_InvalidPartyID(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseInput()
	input.DecentralizedPartyID = "no-double-colon"
	_, err := operations.ExecuteSequence(b, kick.KickSequence, newDeps("p1"), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decentralized_party_id")
}

// TestKickSequence_KickedFingerprintNotInDNS verifies that passing a fingerprint
// that doesn't exist in DNS owners returns an unrecoverable error.
func TestKickSequence_KickedFingerprintNotInDNS(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseInput()
	input.KickedNamespaceFingerprint = "does-not-exist"
	_, err := operations.ExecuteSequence(b, kick.KickSequence, newDeps("p1"), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in DNS owners")
}

// TestKickSequence_KickedParticipantNotInP2P verifies that passing a
// participant UID that doesn't exist in the P2P mapping returns an error.
func TestKickSequence_KickedParticipantNotInP2P(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseInput()
	input.KickedParticipantID = "ghost-participant"
	_, err := operations.ExecuteSequence(b, kick.KickSequence, newDeps("p2"), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in P2P mapping")
}

// TestKickSequence_InsufficientRemainingParticipants verifies that if the
// current DNS threshold exceeds the remaining participant count the sequence
// returns an unrecoverable error.
func TestKickSequence_InsufficientRemainingParticipants(t *testing.T) {
	t.Parallel()
	b := optest.NewBundle(t)
	input := baseInput()
	// Only 1 remaining but currentDNSThreshold=2.
	input.RemainingParticipants = []string{"p1"}
	_, err := operations.ExecuteSequence(b, kick.KickSequence, newDeps("p1"), input)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "kick is impossible")
}

// TestKickSequence_SerializationValid checks that all IN/OUT types can be
// round-tripped through the JSON serializer used by the operations framework.
func TestKickSequence_SerializationValid(t *testing.T) {
	t.Parallel()
	lggr := logger.Nop()
	require.True(t, operations.IsSerializable(lggr, kick.KickInput{}), "KickInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.KickOutput{}), "KickOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.ReadCurrentStateInput{}), "ReadCurrentStateInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.ReadCurrentStateOutput{}), "ReadCurrentStateOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.CreateKickDNSProposalInput{}), "CreateKickDNSProposalInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.CreateKickDNSProposalOutput{}), "CreateKickDNSProposalOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.SignKickDNSProposalInput{}), "SignKickDNSProposalInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.SignKickDNSProposalOutput{}), "SignKickDNSProposalOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.SubmitKickDNSInput{}), "SubmitKickDNSInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.SubmitKickDNSOutput{}), "SubmitKickDNSOutput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.ProposeKickP2PInput{}), "ProposeKickP2PInput must be serializable")
	require.True(t, operations.IsSerializable(lggr, kick.ProposeKickP2POutput{}), "ProposeKickP2POutput must be serializable")
}

// TestGRPCCantonClientImplementsInterface confirms that *client.GRPCCantonClient
// satisfies client.CantonClient after the addition of GetDNS, GetP2P, and
// ListDecentralizedNamespaces.
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

// TestRoundTripProtoMarshal verifies that SignedTopologyTransaction proto bytes
// survive a base64 encode/decode round-trip.
func TestRoundTripProtoMarshal(t *testing.T) {
	t.Parallel()
	original := &protov30.SignedTopologyTransaction{
		Transaction: []byte("kick-test-transaction-bytes"),
		Proposal:    true,
	}
	b, err := proto.Marshal(original)
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(b)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	var restored protov30.SignedTopologyTransaction
	require.NoError(t, proto.Unmarshal(decoded, &restored))
	assert.Equal(t, original.GetTransaction(), restored.GetTransaction())
	assert.Equal(t, original.GetProposal(), restored.GetProposal())
}
