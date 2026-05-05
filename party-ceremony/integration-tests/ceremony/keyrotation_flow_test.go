package tests

import (
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

type KeyRotationFlowTestSuite struct {
	OnboardingFlowTestSuite
}

func (s *KeyRotationFlowTestSuite) SetupSuite() {
	s.OnboardingFlowTestSuite.SetupSuite()
	s.performOnboarding(s.T(), operations.NewMemoryReporter())
}

// TestOnboardingFlow overrides the inherited method to be a lightweight
// assertion. The full onboarding was already run inside SetupSuite.
func (s *KeyRotationFlowTestSuite) TestOnboardingFlow() {
	t := s.T()
	require.NotEmpty(t, s.PartyID, "PartyID should have been set during SetupSuite")
	t.Logf("Onboarding already performed in SetupSuite; PartyID=%s", s.PartyID)
}

// ── Test ─────────────────────────────────────────────────────────────────────

// TestKeyRotationFlow validates the full end-to-end namespace and protocol key rotation
// ceremony:
//
//  1. Run the onboarding ceremony to create a 3-member decentralized party.
//  2. Run the key rotation ceremony to rotate participant 1's namespace and DAML keys.
//
// Key rotation ceremony progression (namespace + DAML, threshold=3):
//
//   - Run 1 (p1, target): reads state -> generates new key -> NSD -> DNS proposal -> signs DNS (1/3) -> ErrThresholdNotMet.
//   - Run 2 (p2): signs DNS (2/3) -> ErrThresholdNotMet.
//   - Run 3 (p3): signs DNS (3/3) -> DNS submitted -> P2P proposal (1/3) -> ErrThresholdNotMet.
//   - Run 4 (p1): P2P proposal (2/3) -> ErrThresholdNotMet.
//   - Run 5 (p2): P2P proposal (3/3) -> SUCCESS.
func (s *KeyRotationFlowTestSuite) TestKeyRotationFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	t.Log("Phase 1: onboarding a 3-member decentralized party")

	// Discover the target participant's current namespace fingerprint from the DNS state.
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]
	dnsState, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after onboarding")
	require.Len(t, dnsState.Owners, 3, "should have 3 owners after onboarding")

	// Identify the target (actor[0]) namespace fingerprint.
	targetNSFP, err := actors[0].deps.Client.GetNamespaceFingerprint(t.Context(), s.onboardingNamespaceName(), synchronizerID, dnsState.Owners)
	require.NoError(t, err, "GetNamespaceFingerprint for target participant")
	p2pState, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P before rotation")
	require.NotNil(t, p2pState.PartySigningKeys, "P2P signing keys should be present before rotation")
	oldTargetProtocolFP, oldTargetProtocolKeyB64, err := actors[0].deps.Client.GetProtocolKeyFingerprint(t.Context(), p2pState.PartySigningKeys.Keys)
	require.NoError(t, err, "GetProtocolKeyFingerprint for target participant before rotation")

	targetUID := actors[0].uid
	t.Logf("Rotating namespace key for participant: %s (ns fingerprint: %s)", targetUID, targetNSFP)

	// ── Phase 2: Rotate namespace and protocol keys ──────────────────────────
	t.Log("Phase 2: key rotation ceremony (namespace + protocol)")

	rotationInput := keyrotation.KeyRotationInput{
		DecentralizedPartyID:       s.PartyID,
		TargetParticipantID:        targetUID,
		TargetNamespaceFingerprint: targetNSFP,
		SynchronizerID:             synchronizerID,
		RotateNamespaceKey:         true,
		RotateDamlKey:              true,
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	rotationKMS := s.kmsConfigFor(0, "rotation")
	deps := [3]ceremony.CantonDeps{
		s.depsFor(0, rotationKMS),
		s.OnboardingDeps(1),
		s.OnboardingDeps(2),
	}

	// Run 1 (p1, target): generates key, NSD, DNS proposal, signs DNS (1/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 1: p1 (target) generates key + signs DNS (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[0], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (1/3)")

	// Run 2 (p2): signs DNS (2/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 2: p2 signs DNS (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[1], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 2: DNS threshold not met (2/3)")

	// Run 3 (p3): signs DNS (3/3) -> DNS submitted -> P2P proposal (1/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 3: p3 signs DNS (3/3) + P2P proposal (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[2], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 3: P2P threshold not met (1/3)")

	// Run 4 (p1): P2P proposal (2/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 4: p1 proposes P2P (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[0], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 4: P2P threshold not met (2/3)")

	// Run 5 (p2): P2P proposal (3/3) -> SUCCESS.
	t.Log("Rotation run 5: p2 proposes P2P (3/3) — completing rotation")
	rotationResult, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[1], rotationInput)
	require.NoError(t, err, "key rotation ceremony should complete successfully")

	// ── Verify output ────────────────────────────────────────────────────────
	assert.True(t, rotationResult.Output.NamespaceKeyRotated, "NamespaceKeyRotated should be true")
	assert.True(t, rotationResult.Output.DamlKeyRotated, "DamlKeyRotated should be true")
	assert.True(t, rotationResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.True(t, rotationResult.Output.P2PUpdated, "P2PUpdated should be true")
	assert.NotEmpty(t, rotationResult.Output.NewNamespaceFingerprint, "NewNamespaceFingerprint should be set")
	assert.NotEmpty(t, rotationResult.Output.NewDamlKeyFingerprint, "NewDamlKeyFingerprint should be set")

	newNSFP := rotationResult.Output.NewNamespaceFingerprint

	// ── Verify via topology read API ─────────────────────────────────────────
	updatedDNS, err := actors[0].deps.Client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after rotation")
	assert.Len(t, updatedDNS.Owners, 3, "should still have 3 owners after rotation")
	assert.NotContains(t, updatedDNS.Owners, targetNSFP, "old namespace fingerprint must be removed from DNS")
	assert.Contains(t, updatedDNS.Owners, newNSFP, "new namespace fingerprint must be present in DNS")
	updatedP2P, err := actors[0].deps.Client.GetP2P(t.Context(), s.PartyID, synchronizerID)
	require.NoError(t, err, "GetP2P after rotation")
	require.NotNil(t, updatedP2P.PartySigningKeys, "P2P signing keys should be present after rotation")
	assert.Len(t, updatedP2P.PartySigningKeys.Keys, len(p2pState.PartySigningKeys.Keys), "rotation should not change signing-key count")
	assert.NotContains(t, updatedP2P.PartySigningKeys.Keys, oldTargetProtocolKeyB64, "old protocol signing key must be removed from P2P")
	newTargetProtocolFP, newTargetProtocolKeyB64, err := actors[0].deps.Client.GetProtocolKeyFingerprint(t.Context(), updatedP2P.PartySigningKeys.Keys)
	require.NoError(t, err, "GetProtocolKeyFingerprint after rotation")
	assert.NotEqual(t, oldTargetProtocolFP, newTargetProtocolFP, "target protocol fingerprint should change")
	assert.Equal(t, rotationResult.Output.NewDamlKeyFingerprint, newTargetProtocolFP, "new protocol fingerprint should match rotation output")
	assert.Contains(t, updatedP2P.PartySigningKeys.Keys, newTargetProtocolKeyB64, "new protocol signing key must be present in P2P")
	s.assertKMSKeysRegistered(0, rotationKMS)

	// ── Verify idempotency ───────────────────────────────────────────────────
	t.Log("Rotation idempotency check")
	cachedResult, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[0], rotationInput)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, rotationResult.Output, cachedResult.Output, "cached result should match")
	s.assertReportsDoNotContainKMS(sharedReporter)

	contractKMS := []client.KMSConfig{
		rotationKMS,
		s.kmsConfigFor(1, "onboarding"),
		s.kmsConfigFor(2, "onboarding"),
	}
	_, recorders := s.runContractDeployFlow(t, s.PartyID, "post-rotation-contract-deploy", contractKMS)
	assertRecordersMatchKMSConfig(t, contractKMS, recorders)
}
