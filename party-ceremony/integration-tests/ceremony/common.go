package tests

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	ceremonyruntime "github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/runtime"
)

type Actor struct {
	deps ceremony.CantonDeps
	uid  string // Canton participant UID
}

type CeremonyTestSuite struct {
	suite.Suite

	chain          *canton.Chain
	SynchronizerID string
	ParticipantIDs []string
	Actors         []Actor
	Participants   []ceremonyruntime.Participant
}

func (s *CeremonyTestSuite) SetupSuite() {
	t := s.T()

	if s.chain == nil {
		// Uncomment to use local Canton environment for faster test runs (requires local setup).
		chain, err := s.NewLocalEnv()
		require.NoError(t, err, "failed to create local environment")

		// chain, err := integrationtests.LoadChainWithCTF(t, 3)
		// require.NoError(t, err, "failed to load chain with CTF")
		// require.Len(t, chain.Participants, 3, "expected 3 participants")
		s.chain = chain
	}

	runtimeParticipants := make([]ceremonyruntime.Participant, 3)
	for i := range 3 {
		p, err := ceremonyruntime.FromCantonParticipant(s.chain.Participants[i])
		require.NoError(t, err, "runtime participant %d", i+1)
		runtimeParticipants[i] = p
	}

	// Build ceremony deps for each participant through the public runtime bridge.
	actors := make([]Actor, 3)
	for i := range 3 {
		deps, cleanup, err := ceremonyruntime.NewOnboardingDeps(t.Context(), runtimeParticipants[i], logger.Test(t), nil)
		require.NoError(t, err, "onboarding deps for participant %d", i+1)
		t.Cleanup(func() { _ = cleanup() })

		uid, err := ceremonyruntime.ParticipantUID(t.Context(), runtimeParticipants[i])
		require.NoError(t, err, "participant %d UID", i+1)
		actors[i] = Actor{
			deps: deps,
			uid:  uid,
		}
		t.Logf("Participant %d: %s", i+1, actors[i].uid)
	}

	synchronizerID, err := ceremonyruntime.DiscoverSynchronizerID(t.Context(), runtimeParticipants[0])
	require.NoError(t, err, "discover synchronizer ID")
	t.Logf("Discovered synchronizer ID: %s", synchronizerID)

	participantIDs := []string{actors[0].uid, actors[1].uid, actors[2].uid}

	s.SynchronizerID = synchronizerID
	s.ParticipantIDs = participantIDs
	s.Actors = actors
	s.Participants = runtimeParticipants
}

func (s *CeremonyTestSuite) OnboardingDeps(i int) ceremony.CantonDeps {
	deps := s.Actors[i].deps
	deps.Logger = logger.Test(s.T())

	return deps
}

func (s *CeremonyTestSuite) OnboardingDepsWithConfirmer(i int, confirmer ceremony.Confirmer) ceremony.CantonDeps {
	deps := s.OnboardingDeps(i)
	deps.Confirmer = confirmer

	return deps
}

func (s *CeremonyTestSuite) NewLocalEnv() (*canton.Chain, error) {
	t := s.T()
	t.Helper()
	// Participant ports follow the localnet compose convention:
	// admin: 1201, 1202, 1203 — ledger: 1301, 1302, 1303
	type localPorts struct{ admin, ledger int }
	ports := []localPorts{{1201, 1301}, {1202, 1302}, {1203, 1303}}
	participants := make([]canton.Participant, len(ports))
	for i, p := range ports {
		adminURL := fmt.Sprintf("localhost:%d", p.admin)
		ledgerURL := fmt.Sprintf("localhost:%d", p.ledger)
		conn, err := grpc.NewClient(adminURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "failed to dial admin API for participant %d", i+1)
		t.Cleanup(func() { _ = conn.Close() })
		ledgerConn, err := grpc.NewClient(ledgerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err, "failed to dial ledger API for participant %d", i+1)
		t.Cleanup(func() { _ = ledgerConn.Close() })
		ledgerServices := canton.CreateLedgerServiceClients(ledgerConn)
		// UserID matches the additional-admin-user-id in simple-topology.conf.
		userID := fmt.Sprintf("user-participant%d", i+1)
		participants[i] = canton.Participant{
			Name: fmt.Sprintf("participant%c", 'A'+i),
			Endpoints: canton.ParticipantEndpoints{
				AdminAPIURL:      adminURL,
				GRPCLedgerAPIURL: ledgerURL,
			},
			AdminServices:  new(canton.CreateAdminServiceClients(conn)),
			LedgerServices: ledgerServices,
			UserID:         userID,
		}
	}

	return &canton.Chain{
		Participants: participants,
	}, nil
}
