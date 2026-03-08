package tests

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
)

// SharedCantonEnvironment holds shared test infrastructure that can be reused across tests.
// This includes the participant connection, uploaded package IDs, and pre-configured signers.
// Contract instances (MCMS, Counter, etc.) are NOT shared as each test needs its own isolated state.
type SharedCantonEnvironment struct {
	Participant   canton.Participant
	McmsPkgID     string
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

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant)
		if err != nil {
			errSharedEnv = err
			return
		}

		if len(packageIDs) == 0 {
			errSharedEnv = err
			return
		}

		mcmsPkgID := packageIDs[0]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedEnv = &SharedCantonEnvironment{
			Participant:   participant,
			McmsPkgID:     mcmsPkgID,
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

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar}, participant, userParticipant)
		if err != nil {
			errSharedTwoPartEnv = err
			return
		}

		if len(packageIDs) == 0 {
			errSharedTwoPartEnv = err
			return
		}

		mcmsPkgID := packageIDs[0]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedTwoPartEnv = &SharedTwoParticipantEnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsPkgID:     mcmsPkgID,
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

		packageIDs, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), [][]byte{mcmsDar, commonDar, tarDar}, participant)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		if len(packageIDs) < 3 {
			errSharedTAREnv = err

			return
		}

		mcmsPkgID := packageIDs[0]
		tarPkgID := packageIDs[len(packageIDs)-1]
		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedTAREnv = &SharedTAREnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsPkgID:     mcmsPkgID,
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
