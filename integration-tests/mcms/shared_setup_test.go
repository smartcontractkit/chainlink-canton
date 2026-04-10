package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// SharedCantonEnvironment holds shared test infrastructure that can be reused across tests.
// This includes the participant connection, uploaded package IDs, and pre-configured signers.
// Contract instances (MCMS, Counter, etc.) are NOT shared as each test needs its own isolated state.
type SharedCantonEnvironment struct {
	Participant   canton.Participant
	McmsPkgID     string
	McmsTestPkgID string
	McmsEncoder   mcms.MCMSEncoder
	CcipOwner     string
	Signers       []*MCMSSigner
	SortedSigners []*MCMSSigner
	Config        MCMSConfig
}

// SharedTwoParticipantEnvironment extends SharedCantonEnvironment with a second participant.
// Used for tests that require a random user on a separate participant.
type SharedTwoParticipantEnvironment struct {
	SharedCantonEnvironment
	UserParticipant canton.Participant
	RandomUser      string
}

// SharedTAREnvironment extends SharedCantonEnvironment with TokenAdminRegistry package.
// Used for tests that interact with TokenAdminRegistry contracts.
type SharedTAREnvironment struct {
	SharedCantonEnvironment
	TarPkgID string
}

// SharedCCIPMCMSEnvironment extends SharedCantonEnvironment with all CCIP + Factory packages.
// Used for tests that validate MCMS-driven CCIP deployment and configuration.
type SharedCCIPMCMSEnvironment struct {
	SharedCantonEnvironment
	CCIPCommonPkgID string
	FactoryPkgID    string
	FactoryEncoder  factory.MCMSEncoder
}

// SharedCCIPMCMSTwoParticipantEnvironment extends SharedCCIPMCMSEnvironment with a second party.
// Used for tests that validate full MCMS governance flow where a bootstrap party deploys the factory
// and then hands over ownership to MCMS party.
// Both parties are on the same participant to enable multi-party submissions.
type SharedCCIPMCMSTwoParticipantEnvironment struct {
	SharedCCIPMCMSEnvironment
	BootstrapParty string // Second party on the same participant
}

var (
	sharedEnv     *SharedCantonEnvironment
	sharedEnvOnce sync.Once
	errSharedEnv  error

	sharedTwoPartEnv     *SharedTwoParticipantEnvironment
	sharedTwoPartEnvOnce sync.Once
	errSharedTwoPartEnv  error

	sharedTAREnv     *SharedTAREnvironment
	sharedTAREnvOnce sync.Once
	errSharedTAREnv  error

	sharedCCIPMCMSEnv     *SharedCCIPMCMSEnvironment
	sharedCCIPMCMSEnvOnce sync.Once
	errSharedCCIPMCMSEnv  error

	sharedCCIPMCMSTwoPartEnv     *SharedCCIPMCMSTwoParticipantEnvironment
	sharedCCIPMCMSTwoPartEnvOnce sync.Once
	errSharedCCIPMCMSTwoPartEnv  error
)

// GetSharedEnvironment initializes the shared test environment once and returns it.
// This uses sync.Once to ensure thread-safe initialization even when tests run in parallel.
// The environment includes a participant connection, uploaded MCMS DAR, and pre-configured signers.
func GetSharedEnvironment(t *testing.T) *SharedCantonEnvironment {
	t.Helper()

	sharedEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		// Upload MCMS DAR
		mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
		if err != nil {
			errSharedEnv = err
			return
		}

		// Upload MCMS Test DAR (contains Counter)
		mcmsTestDar, err := contracts.GetDar(contracts.MCMSTest, contracts.CurrentVersion)
		if err != nil {
			errSharedEnv = err
			return
		}

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar, mcmsTestDar}, participant)
		if err != nil {
			errSharedEnv = err
			return
		}

		if len(packageIDs) < 2 {
			errSharedEnv = err
			return
		}

		mcmsPkgID := packageIDs[0]
		mcmsTestPkgID := packageIDs[1]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedEnv = &SharedCantonEnvironment{
			Participant:   participant,
			McmsPkgID:     mcmsPkgID,
			McmsTestPkgID: mcmsTestPkgID,
			McmsEncoder:   NewMCMSEncoder(mcmsPkgID),
			CcipOwner:     participant.PartyID,
			Signers:       signers,
			SortedSigners: sortedSigners,
			Config:        New2of3Config(signers),
		}
	})

	require.NoError(t, errSharedEnv, "failed to initialize shared environment")
	require.NotNil(t, sharedEnv, "shared environment is nil")

	return sharedEnv
}

// GetSharedTwoParticipantEnvironment initializes a shared environment with two participants.
// Used by tests that need a separate user participant (e.g., testing contract disclosure).
func GetSharedTwoParticipantEnvironment(t *testing.T) *SharedTwoParticipantEnvironment {
	t.Helper()

	sharedTwoPartEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))
		participant := env.Chain.Participants[0]
		userParticipant := env.Chain.Participants[1]

		// Upload MCMS DAR to both participants
		mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
		if err != nil {
			errSharedTwoPartEnv = err
			return
		}

		// Upload MCMS Test DAR (contains Counter) to both participants
		mcmsTestDar, err := contracts.GetDar(contracts.MCMSTest, contracts.CurrentVersion)
		if err != nil {
			errSharedTwoPartEnv = err
			return
		}

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar, mcmsTestDar}, participant, userParticipant)
		if err != nil {
			errSharedTwoPartEnv = err
			return
		}

		if len(packageIDs) < 2 {
			errSharedTwoPartEnv = err
			return
		}

		mcmsPkgID := packageIDs[0]
		mcmsTestPkgID := packageIDs[1]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedTwoPartEnv = &SharedTwoParticipantEnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsPkgID:     mcmsPkgID,
				McmsTestPkgID: mcmsTestPkgID,
				McmsEncoder:   NewMCMSEncoder(mcmsPkgID),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
			UserParticipant: userParticipant,
			RandomUser:      userParticipant.PartyID,
		}
	})

	require.NoError(t, errSharedTwoPartEnv, "failed to initialize two-participant environment")
	require.NotNil(t, sharedTwoPartEnv, "two-participant environment is nil")

	return sharedTwoPartEnv
}

// GetSharedTAREnvironment initializes a shared environment with TokenAdminRegistry support.
// Used by tests that interact with TokenAdminRegistry contracts.
func GetSharedTAREnvironment(t *testing.T) *SharedTAREnvironment {
	t.Helper()

	sharedTAREnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		mcmsDar, err := contracts.GetDar(contracts.MCMS, contracts.CurrentVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		mcmsTestDar, err := contracts.GetDar(contracts.MCMSTest, contracts.CurrentVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		tarDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar, mcmsTestDar, commonDar, tarDar}, participant)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		if len(packageIDs) < 4 {
			errSharedTAREnv = err

			return
		}

		mcmsPkgID := packageIDs[0]
		mcmsTestPkgID := packageIDs[1]
		tarPkgID := packageIDs[len(packageIDs)-1]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedTAREnv = &SharedTAREnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsPkgID:     mcmsPkgID,
				McmsTestPkgID: mcmsTestPkgID,
				McmsEncoder:   NewMCMSEncoder(mcmsPkgID),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
			TarPkgID: tarPkgID,
		}
	})

	require.NoError(t, errSharedTAREnv, "failed to initialize TAR environment")
	require.NotNil(t, sharedTAREnv, "TAR environment is nil")

	return sharedTAREnv
}

// GetSharedCCIPMCMSEnvironment initializes a shared environment with all CCIP and MCMS packages.
// Uploads: MCMS, MCMSTest, CCIPCommon, CCIPRMN, CCIPOffRamp, CCIPOnRamp, CCIPFeeQuoter,
// CCIPTokenAdminRegistry, CCIPCommitteeVerifier, CCIPPerPartyRouter, CCIPLockReleaseTokenPool,
// CCIPPoolInterfaces, CCIPExecutor, CCIPFactory.
func GetSharedCCIPMCMSEnvironment(t *testing.T) *SharedCCIPMCMSEnvironment {
	t.Helper()

	sharedCCIPMCMSEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		darPackages := []contracts.Package{
			contracts.MCMS,
			contracts.MCMSTest,
			contracts.CCIPCommon,
			contracts.CCIPRMN,
			contracts.CCIPOffRamp,
			contracts.CCIPOnRamp,
			contracts.CCIPFeeQuoter,
			contracts.CCIPTokenAdminRegistry,
			contracts.CCIPCommitteeVerifier,
			contracts.CCIPPerPartyRouter,
			contracts.CCIPPoolInterfaces,
			contracts.CCIPLockReleaseTokenPool,
			contracts.CCIPExecutor,
			contracts.CCIPFactory,
		}

		darBytes := make([][]byte, 0, len(darPackages))
		for _, pkg := range darPackages {
			dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
			if err != nil {
				errSharedCCIPMCMSEnv = err
				return
			}
			darBytes = append(darBytes, dar)
		}

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participant)
		if err != nil {
			errSharedCCIPMCMSEnv = err
			return
		}

		if len(packageIDs) < len(darPackages) {
			errSharedCCIPMCMSEnv = fmt.Errorf("expected %d package IDs, got %d", len(darPackages), len(packageIDs))
			return
		}

		mcmsPkgID := packageIDs[0]
		mcmsTestPkgID := packageIDs[1]
		ccipCommonPkgID := packageIDs[2]
		factoryPkgID := packageIDs[len(packageIDs)-1]

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		factoryEncoder := factory.NewContract(factoryPkgID, "CCIP.Factory", "CCIPFactory").Encoder()

		sharedCCIPMCMSEnv = &SharedCCIPMCMSEnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsPkgID:     mcmsPkgID,
				McmsTestPkgID: mcmsTestPkgID,
				McmsEncoder:   NewMCMSEncoder(mcmsPkgID),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
			CCIPCommonPkgID: ccipCommonPkgID,
			FactoryPkgID:    factoryPkgID,
			FactoryEncoder:  factoryEncoder,
		}
	})

	require.NoError(t, errSharedCCIPMCMSEnv, "failed to initialize CCIP MCMS environment")
	require.NotNil(t, sharedCCIPMCMSEnv, "CCIP MCMS environment is nil")

	return sharedCCIPMCMSEnv
}

// GetSharedCCIPMCMSTwoParticipantEnvironment initializes a shared environment with two parties
// on the same participant and all CCIP and MCMS packages. Used for tests that validate full MCMS
// governance flow where a bootstrap party deploys the factory and then hands over ownership to MCMS party.
// Both parties are on the same participant to enable multi-party submissions (ActAs with both parties).
func GetSharedCCIPMCMSTwoParticipantEnvironment(t *testing.T) *SharedCCIPMCMSTwoParticipantEnvironment {
	t.Helper()

	sharedCCIPMCMSTwoPartEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		// Ensure user has explicit CanActAs rights for the primary party
		// This is needed for multi-party submissions where both parties must be authorized
		testhelpers.GrantCanActAs(t, participant, participant.PartyID)

		// Allocate a second party on the same participant for bootstrapping
		// Use unique hint to avoid conflicts across test runs
		bootstrapParty := testhelpers.AllocateParty(t, participant, fmt.Sprintf("bootstrap-%s", uuid.New().String()[:8]))

		darPackages := []contracts.Package{
			contracts.MCMS,
			contracts.MCMSTest,
			contracts.CCIPCommon,
			contracts.CCIPRMN,
			contracts.CCIPOffRamp,
			contracts.CCIPOnRamp,
			contracts.CCIPFeeQuoter,
			contracts.CCIPTokenAdminRegistry,
			contracts.CCIPCommitteeVerifier,
			contracts.CCIPPerPartyRouter,
			contracts.CCIPPoolInterfaces,
			contracts.CCIPLockReleaseTokenPool,
			contracts.CCIPExecutor,
			contracts.CCIPFactory,
		}

		darBytes := make([][]byte, 0, len(darPackages))
		for _, pkg := range darPackages {
			dar, err := contracts.GetDar(pkg, contracts.CurrentVersion)
			if err != nil {
				errSharedCCIPMCMSTwoPartEnv = err
				return
			}
			darBytes = append(darBytes, dar)
		}

		// Upload DARs to the participant
		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participant)
		if err != nil {
			errSharedCCIPMCMSTwoPartEnv = err
			return
		}

		if len(packageIDs) < len(darPackages) {
			errSharedCCIPMCMSTwoPartEnv = fmt.Errorf("expected %d package IDs, got %d", len(darPackages), len(packageIDs))
			return
		}

		mcmsPkgID := packageIDs[0]
		mcmsTestPkgID := packageIDs[1]
		ccipCommonPkgID := packageIDs[2]
		factoryPkgID := packageIDs[len(packageIDs)-1]

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		factoryEncoder := factory.NewContract(factoryPkgID, "CCIP.Factory", "CCIPFactory").Encoder()

		sharedCCIPMCMSTwoPartEnv = &SharedCCIPMCMSTwoParticipantEnvironment{
			SharedCCIPMCMSEnvironment: SharedCCIPMCMSEnvironment{
				SharedCantonEnvironment: SharedCantonEnvironment{
					Participant:   participant,
					McmsPkgID:     mcmsPkgID,
					McmsTestPkgID: mcmsTestPkgID,
					McmsEncoder:   NewMCMSEncoder(mcmsPkgID),
					CcipOwner:     participant.PartyID,
					Signers:       signers,
					SortedSigners: sortedSigners,
					Config:        New2of3Config(signers),
				},
				CCIPCommonPkgID: ccipCommonPkgID,
				FactoryPkgID:    factoryPkgID,
				FactoryEncoder:  factoryEncoder,
			},
			BootstrapParty: bootstrapParty,
		}
	})

	require.NoError(t, errSharedCCIPMCMSTwoPartEnv, "failed to initialize CCIP MCMS two-participant environment")
	require.NotNil(t, sharedCCIPMCMSTwoPartEnv, "CCIP MCMS two-participant environment is nil")

	return sharedCCIPMCMSTwoPartEnv
}
