package tests

import (
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/keyrotation"
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

// TestKeyRotationFlow validates the full end-to-end namespace key rotation
// ceremony:
//
//  1. Run the onboarding ceremony to create a 3-member decentralized party.
//  2. Run the key rotation ceremony to rotate participant 1's namespace key.
//
// Key rotation ceremony progression (namespace-only, threshold=3):
//
//   - Run 1 (p1, target): reads state -> generates new key -> NSD -> DNS proposal -> signs DNS (1/3) -> ErrThresholdNotMet.
//   - Run 2 (p2): signs DNS (2/3) -> ErrThresholdNotMet.
//   - Run 3 (p3): signs DNS (3/3) -> DNS submitted -> DNS confirmed -> SUCCESS.
func (s *KeyRotationFlowTestSuite) TestKeyRotationFlow() {
	t := s.T()

	actors := s.Actors
	synchronizerID := s.SynchronizerID

	t.Log("Phase 1: onboarding a 3-member decentralized party")

	// Discover the target participant's current namespace fingerprint from the DNS state.
	decNS := strings.SplitN(s.PartyID, "::", 2)[1]
	dnsState, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after onboarding")
	require.Len(t, dnsState.Owners, 3, "should have 3 owners after onboarding")

	// Identify the target (actor[0]) namespace fingerprint.
	targetNSFP, err := actors[0].client.GetNamespaceFingerprint(t.Context(), onboardingNamespaceName, synchronizerID, dnsState.Owners)
	require.NoError(t, err, "GetNamespaceFingerprint for target participant")

	targetUID := actors[0].uid
	t.Logf("Rotating namespace key for participant: %s (ns fingerprint: %s)", targetUID, targetNSFP)

	// ── Phase 2: Rotate namespace key ────────────────────────────────────────
	t.Log("Phase 2: key rotation ceremony (namespace-only)")

	rotationInput := keyrotation.KeyRotationInput{
		DecentralizedPartyID:       s.PartyID,
		TargetParticipantID:        targetUID,
		TargetNamespaceFingerprint: targetNSFP,
		SynchronizerID:             synchronizerID,
		RotateNamespaceKey:         true,
		RotateDamlKey:              false,
	}

	sharedReporter := operations.NewMemoryReporter()
	newBundle := func() operations.Bundle {
		return operations.NewBundle(t.Context, logger.Test(t), sharedReporter)
	}

	deps := [3]ceremony.CantonDeps{
		{Client: actors[0].client, Logger: logger.Test(t)},
		{Client: actors[1].client, Logger: logger.Test(t)},
		{Client: actors[2].client, Logger: logger.Test(t)},
	}

	// Run 1 (p1, target): generates key, NSD, DNS proposal, signs DNS (1/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 1: p1 (target) generates key + signs DNS (1/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[0], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 1: DNS threshold not met (1/3)")

	// Run 2 (p2): signs DNS (2/3) -> ErrThresholdNotMet.
	t.Log("Rotation run 2: p2 signs DNS (2/3)")
	_, err = operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[1], rotationInput)
	require.ErrorContains(t, err, keyrotation.ErrThresholdNotMet.Error(), "run 2: DNS threshold not met (2/3)")

	// Run 3 (p3): signs DNS (3/3) -> DNS submitted -> DNS confirmed -> SUCCESS.
	t.Log("Rotation run 3: p3 signs DNS (3/3) — completing rotation")
	rotationResult, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[2], rotationInput)
	require.NoError(t, err, "key rotation ceremony should complete successfully")

	// ── Verify output ────────────────────────────────────────────────────────
	assert.True(t, rotationResult.Output.NamespaceKeyRotated, "NamespaceKeyRotated should be true")
	assert.False(t, rotationResult.Output.DamlKeyRotated, "DamlKeyRotated should be false")
	assert.True(t, rotationResult.Output.DNSUpdated, "DNSUpdated should be true")
	assert.False(t, rotationResult.Output.P2PUpdated, "P2PUpdated should be false")
	assert.NotEmpty(t, rotationResult.Output.NewNamespaceFingerprint, "NewNamespaceFingerprint should be set")

	newNSFP := rotationResult.Output.NewNamespaceFingerprint

	// ── Verify via topology read API ─────────────────────────────────────────
	updatedDNS, err := actors[0].client.GetDNS(t.Context(), decNS, synchronizerID)
	require.NoError(t, err, "GetDNS after rotation")
	assert.Len(t, updatedDNS.Owners, 3, "should still have 3 owners after rotation")
	assert.NotContains(t, updatedDNS.Owners, targetNSFP, "old namespace fingerprint must be removed from DNS")
	assert.Contains(t, updatedDNS.Owners, newNSFP, "new namespace fingerprint must be present in DNS")

	// ── Verify idempotency ───────────────────────────────────────────────────
	t.Log("Rotation idempotency check")
	cachedResult, err := operations.ExecuteSequence(newBundle(), keyrotation.KeyRotationSequence, deps[0], rotationInput)
	require.NoError(t, err, "idempotent re-run should succeed")
	assert.Equal(t, rotationResult.Output, cachedResult.Output, "cached result should match")
}
