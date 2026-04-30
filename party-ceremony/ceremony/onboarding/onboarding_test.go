package onboarding_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/onboarding"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/keys"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// ── Mock Canton client ───────────────────────────────────────────────────────

type mockCantonClient struct {
	participantID    string
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

func (m *mockCantonClient) GetNamespaceKeyName(_ context.Context, _ string, _ []string) (string, error) {
	return "mock-ns-key", nil
}

func (m *mockCantonClient) Authorize(_ context.Context, serial uint32, mapping *protov30.TopologyMapping, _ string, _ bool, _ ...string) (*protov30.SignedTopologyTransaction, error) {
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
			Proposal:    tx.Proposal,
		}
	}

	return result, nil
}

func (m *mockCantonClient) AddTransactions(_ context.Context, _ []*protov30.SignedTopologyTransaction, _ string) error {
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

func (m *mockCantonClient) GetDNS(_ context.Context, _ string, _ string) (*client.DNSState, error) {
	return &client.DNSState{}, nil
}

func (m *mockCantonClient) GetP2P(_ context.Context, _ string, _ string) (*client.P2PState, error) {
	return &client.P2PState{}, nil
}

func (m *mockCantonClient) ListDecentralizedNamespaces(_ context.Context, _ string) ([]*client.DNSState, error) {
	return []*client.DNSState{}, nil
}

func (m *mockCantonClient) GetProtocolKeyFingerprint(_ context.Context, _ []string) (string, string, error) {
	return "mock-protocol-fp", "mock-protocol-key-b64", nil
}

func (m *mockCantonClient) UploadDar(_ context.Context, _ []byte) (string, error) {
	return "mock-package-id", nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func newDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID},
		Logger: logger.Nop(),
	}
}

func newKMSDeps(participantID string) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID},
		KMS: client.KMSConfig{
			NamespaceKeyID: "arn:aws:kms:us-east-1:123456789:key/ns-key-123",
			ProtocolKeyID:  "arn:aws:kms:us-east-1:123456789:key/proto-key-456",
		},
		Logger: logger.Nop(),
	}
}

func newTrackedKMSDeps(participantID string, registrations *int) ceremony.CantonDeps {
	return ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: participantID, kmsRegistrations: registrations},
		KMS: client.KMSConfig{
			NamespaceKeyID: "arn:aws:kms:us-east-1:123456789:key/" + participantID + "-namespace",
			ProtocolKeyID:  "arn:aws:kms:us-east-1:123456789:key/" + participantID + "-protocol",
		},
		Logger: logger.Nop(),
	}
}

// baseInput returns a single-participant input so tests that do not need
// multi-actor coordination stay simple and self-contained.
func baseInput() onboarding.OnboardingInput {
	return onboarding.OnboardingInput{
		NamespaceName:  "test-namespace",
		PartyPrefix:    "test-party",
		Participants:   []string{"p1"},
		SynchronizerID: "global",
		Threshold:      1,
	}
}

// multiActorInput returns a two-participant input for tests that exercise
// the distributed signing / resume logic.
func multiActorInput() onboarding.OnboardingInput {
	return onboarding.OnboardingInput{
		NamespaceName:  "test-namespace",
		PartyPrefix:    "test-party",
		Participants:   []string{"p1", "p2"},
		SynchronizerID: "global",
		Threshold:      2,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

func TestOnboardingSequence_HappyPath(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	sr, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, newDeps("p1"), baseInput())
	require.NoError(t, err)

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
	assert.Contains(t, sr.Output.PartyID, "test-party::")
}

func TestOnboardingSequence_Idempotent(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	deps := newDeps("p1")

	sr1, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, deps, baseInput())
	require.NoError(t, err)

	sr2, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, deps, baseInput())
	require.NoError(t, err)

	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report")
	assert.Equal(t, sr1.Output, sr2.Output)
}

// TestOnboardingSequence_ThresholdNotMet verifies that when a participant runs
// before all required actors have generated their keys, the sequence returns
// ErrThresholdNotMet so the caller knows to retry later.
func TestOnboardingSequence_ThresholdNotMet(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	// p1's node generates p1's key but not p2's → gate fires (1/2) → ErrThresholdNotMet.
	_, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, newDeps("p1"), multiActorInput())
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error())
}

// TestOnboardingSequence_AsyncMultiActorFlow verifies the shared-reporter
// caching contract: once the ceremony is complete, any subsequent invocation
// with the same input returns the same cached sequence report.
func TestOnboardingSequence_AsyncMultiActorFlow(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}

	// Single-participant ceremony — actor "p1" completes it end-to-end.
	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), baseInput())
	require.NoError(t, err)
	assert.NotEmpty(t, sr.Output.PartyID)

	// Any subsequent invocation with the same input must return the cached report.
	srCached, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p2"), baseInput())
	require.NoError(t, err)
	assert.Equal(t, sr.ID, srCached.ID, "repeated invocation must return the cached report")
}

// TestOnboardingSequence_ResumeAfterPartialSigning exercises the full async
// ceremony with 2 participants and threshold=2. Each participant can only act
// for themselves, so four independent runs are required:
//
//	Run 1 (p1): p1 generates key + NSD → gate fires (1/2 keys) → ErrThresholdNotMet.
//	Run 2 (p2): p1 cached + p2 generates key + NSD → gate passes →
//	            DNS proposal → p2 signs DNS (1/2) → ErrThresholdNotMet.
//	Run 3 (p1): all keys cached → p1 signs DNS (2/2) → DNS submitted →
//	            p1 proposes P2P (1/2) → ErrThresholdNotMet.
//	Run 4 (p2): DNS cached → p2 proposes P2P (2/2) → P2P confirmed → SUCCESS.
func TestOnboardingSequence_ResumeAfterPartialSigning(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := multiActorInput() // 2 participants: p1, p2

	// Run 1: p1 generates its key and proposes its NSD. p2's key is missing → gate.
	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 1: key-gen gate (1/2)")

	// Run 2: p1 key cached, p2 generates its key + NSD; gate passes; DNS proposal
	// created; p2 signs (1/2 DNS sigs) → ErrThresholdNotMet.
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p2"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 2: DNS signing threshold (1/2)")

	// Run 3: all keys cached; p1 signs DNS (2/2) → DNS submitted; p1 proposes P2P
	// (1/2) → ErrThresholdNotMet (p2's P2P proposal is still pending).
	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p1"), input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 3: P2P threshold (1/2)")

	// Run 4: DNS cached; p2 proposes P2P (2/2) → mapping confirmed → success.
	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, newDeps("p2"), input)
	require.NoError(t, err, "run 4: all P2P proposals collected — ceremony complete")

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
}

func TestGRPCCantonClientImplementsInterface(t *testing.T) {
	t.Parallel()
	var _ client.CantonClient = (*client.GRPCCantonClient)(nil)
	var _ client.CantonClient = (*mockCantonClient)(nil)
}

func TestRoundTripProtoMarshal(t *testing.T) {
	t.Parallel()

	original := &protov30.SignedTopologyTransaction{
		Transaction: []byte("test-transaction-bytes"),
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

// ── KMS tests ────────────────────────────────────────────────────────────────

// kmsInput returns a single-participant input. KMS key IDs are intentionally
// carried by deps, not shared workflow input.
func kmsInput() onboarding.OnboardingInput {
	return baseInput()
}

func TestOnboardingSequence_KMS_HappyPath(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	sr, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, newKMSDeps("p1"), kmsInput())
	require.NoError(t, err)

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
	assert.Contains(t, sr.Output.PartyID, "test-party::")
}

func TestOnboardingSequence_KMS_DifferentOutputFromGenerate(t *testing.T) {
	t.Parallel()

	// Run with generated keys.
	b1 := optest.NewBundle(t)
	srGen, err := operations.ExecuteSequence(b1, onboarding.OnboardingSequence, newDeps("p1"), baseInput())
	require.NoError(t, err)

	// Run with KMS keys.
	b2 := optest.NewBundle(t)
	srKms, err := operations.ExecuteSequence(b2, onboarding.OnboardingSequence, newKMSDeps("p1"), kmsInput())
	require.NoError(t, err)

	// The KMS path uses different key material (kmsKeyID is hashed into mock
	// output), so the party ID should differ.
	assert.NotEqual(t, srGen.Output.PartyID, srKms.Output.PartyID,
		"KMS and generated paths should produce different keys/party IDs")
}

func TestOnboardingSequence_KMS_Idempotent(t *testing.T) {
	t.Parallel()

	b := optest.NewBundle(t)
	deps := newKMSDeps("p1")
	in := kmsInput()

	sr1, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, deps, in)
	require.NoError(t, err)

	sr2, err := operations.ExecuteSequence(b, onboarding.OnboardingSequence, deps, in)
	require.NoError(t, err)

	assert.Equal(t, sr1.ID, sr2.ID, "second call must return the cached report")
	assert.Equal(t, sr1.Output, sr2.Output)
}

func TestOnboardingSequence_KMS_ResumeUsesLocalConfigPerActor(t *testing.T) {
	t.Parallel()

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Nop(), sharedReporter)
	}
	input := multiActorInput()

	var p1Registrations int
	var p2Registrations int
	deps1 := newTrackedKMSDeps("p1", &p1Registrations)
	deps2 := newTrackedKMSDeps("p2", &p2Registrations)

	_, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 1: key-gen gate (1/2)")

	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps2, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 2: DNS signing threshold (1/2)")

	_, err = operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps1, input)
	require.ErrorContains(t, err, onboarding.ErrThresholdNotMet.Error(), "run 3: P2P threshold (1/2)")

	sr, err := operations.ExecuteSequence(newBundle(), onboarding.OnboardingSequence, deps2, input)
	require.NoError(t, err, "run 4: all P2P proposals collected")

	assert.True(t, sr.Output.DNSConfirmed)
	assert.True(t, sr.Output.P2PConfirmed)
	assert.Equal(t, 2, p1Registrations, "p1 should register namespace and protocol keys from local config")
	assert.Equal(t, 2, p2Registrations, "p2 should register namespace and protocol keys from local config")
}

func TestCreateMemberKeyOp_KMSRequiresBothKeys(t *testing.T) {
	t.Parallel()

	_, err := operations.ExecuteOperation(optest.NewBundle(t), keys.CreateMemberKeyOp, ceremony.CantonDeps{
		Client: &mockCantonClient{participantID: "p1"},
		KMS: client.KMSConfig{
			NamespaceKeyID: "arn:aws:kms:us-east-1:123456789:key/namespace-only",
		},
		Logger: logger.Nop(),
	}, keys.CreateMemberKeyInput{
		NamespaceName: "test-namespace",
		ParticipantID: "p1",
	})
	require.ErrorContains(t, err, "kms_namespace_key_id and kms_protocol_key_id must both be set")
}

// ── Mock Canton client ───────────────────────────────────────────────────────
