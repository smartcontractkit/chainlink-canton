package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/mcms/api"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// SharedCantonEnvironment holds shared test infrastructure that can be reused across tests.
// This includes the participant connection, uploaded package IDs, and pre-configured signers.
// Contract instances (MCMS, Counter, etc.) are NOT shared as each test needs its own isolated state.
type SharedCantonEnvironment struct {
	Participant   canton.Participant
	McmsEncoder   mcmsApi.MCMSEncoder
	CcipOwner     string
	Signers       []*MCMSSigner
	SortedSigners []*MCMSSigner
	Config        MCMSConfig
}

// SharedTAREnvironment extends SharedCantonEnvironment with TokenAdminRegistry package.
// Used for tests that interact with TokenAdminRegistry contracts.
type SharedTAREnvironment struct {
	SharedCantonEnvironment
}

// SharedCCIPMCMSEnvironment extends SharedCantonEnvironment with all CCIP + Factory packages.
// Used for tests that validate MCMS-driven CCIP deployment and configuration.
type SharedCCIPMCMSEnvironment struct {
	SharedCantonEnvironment
	FactoryEncoder factory.MCMSEncoder
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

	sharedTAREnv     *SharedTAREnvironment
	sharedTAREnvOnce sync.Once
	errSharedTAREnv  error

	sharedCCIPMCMSEnv     *SharedCCIPMCMSEnvironment
	sharedCCIPMCMSEnvOnce sync.Once
	errSharedCCIPMCMSEnv  error

	sharedCCIPMCMSTwoPartEnv     *SharedCCIPMCMSTwoParticipantEnvironment
	sharedCCIPMCMSTwoPartEnvOnce sync.Once
	errSharedCCIPMCMSTwoPartEnv  error

	sharedFeeTreasuryEnv     *SharedFeeTreasuryEnvironment
	sharedFeeTreasuryEnvOnce sync.Once
	errSharedFeeTreasuryEnv  error
)

// SharedFeeTreasuryEnvironment adds a second participant hosting the withdrawal recipient,
// so a real Amulet transfer can be received and its balance read as the recipient party.
type SharedFeeTreasuryEnvironment struct {
	SharedCantonEnvironment
	RecipientParticipant canton.Participant
	Recipient            string
}

// GetSharedFeeTreasuryEnvironment initializes a 2-participant environment with the MCMS core
// and CCIPFeeTreasury packages on participant[0] (feeOwner / MCMS owner). participant[1] hosts
// the fee-withdrawal recipient; the Splice token-standard packages it needs are pre-deployed on localnet.
func GetSharedFeeTreasuryEnvironment(t *testing.T) *SharedFeeTreasuryEnvironment {
	t.Helper()

	sharedFeeTreasuryEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(2))
		participant := env.Chain.Participants[0]
		recipientParticipant := env.Chain.Participants[1]

		darPackages := []contracts.Package{
			contracts.MCMSCore,
			contracts.CCIPFeeTreasury,
		}

		darBytes := make([][]byte, 0, len(darPackages))
		for _, pkg := range darPackages {
			dar, err := contracts.GetDar(pkg, contracts.DevVersion)
			if err != nil {
				errSharedFeeTreasuryEnv = err
				return
			}
			darBytes = append(darBytes, dar)
		}

		// Both participants need the fee-treasury package: the recipient (participant[1]) is an
		// observer on the FeeWithdrawalAuthorization, so its participant must resolve/host it.
		if _, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participant, recipientParticipant); err != nil {
			errSharedFeeTreasuryEnv = err
			return
		}

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedFeeTreasuryEnv = &SharedFeeTreasuryEnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsEncoder:   NewMCMSEncoder(),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
			RecipientParticipant: recipientParticipant,
			Recipient:            recipientParticipant.PartyID,
		}
	})

	require.NoError(t, errSharedFeeTreasuryEnv, "failed to initialize fee treasury environment")
	require.NotNil(t, sharedFeeTreasuryEnv, "fee treasury environment is nil")

	return sharedFeeTreasuryEnv
}

// GetSharedEnvironment initializes the shared test environment once and returns it.
// This uses sync.Once to ensure thread-safe initialization even when tests run in parallel.
// The environment includes a participant connection, uploaded MCMS DAR, and pre-configured signers.
func GetSharedEnvironment(t *testing.T) *SharedCantonEnvironment {
	t.Helper()

	sharedEnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		// Upload MCMS DAR
		mcmsDar, err := contracts.GetDar(contracts.MCMSCore, contracts.DevVersion)
		if err != nil {
			errSharedEnv = err
			return
		}

		// Upload MCMS Test DAR (contains Counter)
		mcmsTestDar, err := contracts.GetDar(contracts.MCMSTest, contracts.DevVersion)
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

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedEnv = &SharedCantonEnvironment{
			Participant:   participant,
			McmsEncoder:   NewMCMSEncoder(),
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

// GetSharedTAREnvironment initializes a shared environment with TokenAdminRegistry support.
// Used by tests that interact with TokenAdminRegistry contracts.
func GetSharedTAREnvironment(t *testing.T) *SharedTAREnvironment {
	t.Helper()

	sharedTAREnvOnce.Do(func() {
		env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(1))
		participant := env.Chain.Participants[0]

		mcmsDar, err := contracts.GetDar(contracts.MCMSCore, contracts.DevVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		mcmsTestDar, err := contracts.GetDar(contracts.MCMSTest, contracts.DevVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		commonDar, err := contracts.GetDar(contracts.CCIPRuntimeV2, contracts.DevVersion)
		if err != nil {
			errSharedTAREnv = err

			return
		}

		tarDar, err := contracts.GetDar(contracts.CCIPCoreV2, contracts.DevVersion)
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

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		sharedTAREnv = &SharedTAREnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsEncoder:   NewMCMSEncoder(),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
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
			contracts.MCMSCore,
			contracts.MCMSTest,
			contracts.CCIPRuntimeV2,
			contracts.CCIPCoreV2,
			contracts.CCIPCommitteeVerifierV2,
			contracts.CCIPExtensionAPIV2,
			contracts.CCIPLockReleaseTokenPoolV2,
			contracts.CCIPExecutorV2,
			contracts.CCIPFactoryV2,
		}

		darBytes := make([][]byte, 0, len(darPackages))
		for _, pkg := range darPackages {
			dar, err := contracts.GetDar(pkg, contracts.DevVersion)
			if err != nil {
				errSharedCCIPMCMSEnv = err
				return
			}
			darBytes = append(darBytes, dar)
		}

		_, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participant)
		if err != nil {
			errSharedCCIPMCMSEnv = err
			return
		}

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		factoryEncoder := factory.NewContract(fmt.Sprintf("#%s", factory.PackageName), "CCIP.FactoryV2", "CCIPFactory").Encoder()

		sharedCCIPMCMSEnv = &SharedCCIPMCMSEnvironment{
			SharedCantonEnvironment: SharedCantonEnvironment{
				Participant:   participant,
				McmsEncoder:   NewMCMSEncoder(),
				CcipOwner:     participant.PartyID,
				Signers:       signers,
				SortedSigners: sortedSigners,
				Config:        New2of3Config(signers),
			},
			FactoryEncoder: factoryEncoder,
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

		// Allocate a second party on the same participant for bootstrapping
		// Use unique hint to avoid conflicts across test runs
		bootstrapParty := testhelpers.AllocateParty(t, participant, fmt.Sprintf("bootstrap-%s", uuid.NewString()[:8]))
		testhelpers.GrantCanActAs(t, participant, bootstrapParty)

		darPackages := []contracts.Package{
			contracts.MCMSCore,
			contracts.MCMSTest,
			contracts.CCIPRuntimeV2,
			contracts.CCIPCoreV2,
			contracts.CCIPCommitteeVerifierV2,
			contracts.CCIPExtensionAPIV2,
			contracts.CCIPLockReleaseTokenPoolV2,
			contracts.CCIPExecutorV2,
			contracts.CCIPFactoryV2,
		}

		darBytes := make([][]byte, 0, len(darPackages))
		for _, pkg := range darPackages {
			dar, err := contracts.GetDar(pkg, contracts.DevVersion)
			if err != nil {
				errSharedCCIPMCMSTwoPartEnv = err
				return
			}
			darBytes = append(darBytes, dar)
		}

		// Upload DARs to the participant
		_, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), darBytes, participant)
		if err != nil {
			errSharedCCIPMCMSTwoPartEnv = err
			return
		}

		signers := createSigners(t)
		sortedSigners := SortSignersByAddress(signers)

		factoryEncoder := factory.NewContract(fmt.Sprintf("#%s", factory.PackageName), "CCIP.FactoryV2", "CCIPFactory").Encoder()

		sharedCCIPMCMSTwoPartEnv = &SharedCCIPMCMSTwoParticipantEnvironment{
			SharedCCIPMCMSEnvironment: SharedCCIPMCMSEnvironment{
				SharedCantonEnvironment: SharedCantonEnvironment{
					Participant:   participant,
					McmsEncoder:   NewMCMSEncoder(),
					CcipOwner:     participant.PartyID,
					Signers:       signers,
					SortedSigners: sortedSigners,
					Config:        New2of3Config(signers),
				},
				FactoryEncoder: factoryEncoder,
			},
			BootstrapParty: bootstrapParty,
		}
	})

	require.NoError(t, errSharedCCIPMCMSTwoPartEnv, "failed to initialize CCIP MCMS two-participant environment")
	require.NotNil(t, sharedCCIPMCMSTwoPartEnv, "CCIP MCMS two-participant environment is nil")

	return sharedCCIPMCMSTwoPartEnv
}
