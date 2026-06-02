package tests

import (
	"context"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// signerActor pairs a participant index with its KMS config so the helper
// knows whose signing key to mint a signature with.
type signerActor struct {
	actorIndex int
	kmsCfg     client.KMSConfig
}

// exerciseChoiceAsParty drives a full multi-party InteractiveSubmission flow
// outside of the ceremony framework: PrepareSubmission on `submitter`'s ledger,
// then each `signers` actor signs the prepared hash with its own KMS-backed
// (or vault) signer, then ExecuteSubmission. Returns the first created
// contract ID (empty for archive-only submissions).
//
// Use this when a test needs to perform a single ad-hoc DAML action on behalf
// of the decentralized party — e.g. archive a contract between ceremony runs,
// or submit a transaction once the ceremony has completed.
func (s *CeremonyTestSuite) exerciseChoiceAsParty(
	t *testing.T,
	submitterActorIndex int,
	signers []signerActor,
	partyID string,
	synchronizerID string,
	templateID *apiv2.Identifier,
	contractID string,
	choice string,
	choiceArg *apiv2.Value,
) string {
	t.Helper()

	if choiceArg == nil {
		choiceArg = &apiv2.Value{
			Sum: &apiv2.Value_Record{Record: &apiv2.Record{}},
		}
	}
	commands := []*apiv2.Command{{
		Command: &apiv2.Command_Exercise{
			Exercise: &apiv2.ExerciseCommand{
				TemplateId:     templateID,
				ContractId:     contractID,
				Choice:         choice,
				ChoiceArgument: choiceArg,
			},
		},
	}}
	return s.submitCommandsAsParty(t, submitterActorIndex, signers, partyID, synchronizerID, commands)
}

// submitCommandsAsParty is the underlying primitive used by
// [exerciseChoiceAsParty]: PrepareSubmission with the given commands, sign
// with all `signers`, then ExecuteSubmission.
func (s *CeremonyTestSuite) submitCommandsAsParty(
	t *testing.T,
	submitterActorIndex int,
	signers []signerActor,
	partyID string,
	synchronizerID string,
	commands []*apiv2.Command,
) string {
	t.Helper()
	require.NotEmpty(t, signers, "at least one signer required")
	require.NotEmpty(t, commands, "at least one command required")

	ctx := t.Context()

	// 1. Submitter prepares the transaction.
	submitter := s.chain.Participants[submitterActorIndex]
	submitterLedger, submitterConn := s.NewLedgerClient(submitter)
	t.Cleanup(func() { _ = submitterConn.Close() })

	if submitter.UserID != "" {
		// User must have CanActAs rights on the party to prepare.
		err := submitterLedger.GrantPartyRights(ctx, submitter.UserID, partyID)
		require.NoError(t, err, "GrantPartyRights to submitter user before submit")
	}

	prep, err := submitterLedger.PrepareSubmission(ctx, commands, []string{partyID}, []string{partyID}, synchronizerID)
	require.NoError(t, err, "PrepareSubmission")
	hashBytes := prep.GetPreparedTransactionHash()
	require.NotEmpty(t, hashBytes, "prepared transaction hash should be non-empty")

	// 2. Discover the party's current signing keys via topology.
	adminClient := s.Actors[submitterActorIndex].deps.Client
	p2pState, err := adminClient.GetP2P(ctx, partyID, synchronizerID)
	require.NoError(t, err, "GetP2P to discover party signing keys")
	require.NotNil(t, p2pState.PartySigningKeys, "P2P signing keys must be present")
	knownKeys := p2pState.PartySigningKeys.Keys
	require.NotEmpty(t, knownKeys, "party signing key set must be non-empty")

	// 3. Each signer signs the prepared hash with its own participant's key.
	sigs := make([]*apiv2.Signature, 0, len(signers))
	for _, sa := range signers {
		sigs = append(sigs, s.signHashAs(t, ctx, sa, knownKeys, hashBytes))
	}

	// 4. Submitter executes with the collected signatures.
	partySigs := &interactive.PartySignatures{
		Signatures: []*interactive.SinglePartySignatures{{
			Party:      partyID,
			Signatures: sigs,
		}},
	}

	createdCID, err := submitterLedger.ExecuteSubmission(
		ctx,
		prep.GetPreparedTransaction(),
		partySigs,
		prep.GetHashingSchemeVersion(),
	)
	require.NoError(t, err, "ExecuteSubmission")

	return createdCID
}

// signHashAs builds a signer for one participant using its KMS config (or
// vault fallback) and signs `hash`, returning a Ledger API Signature proto.
func (s *CeremonyTestSuite) signHashAs(
	t *testing.T,
	ctx context.Context,
	sa signerActor,
	knownKeys []string,
	hash []byte,
) *apiv2.Signature {
	t.Helper()

	participant := s.chain.Participants[sa.actorIndex]
	adminConn := s.NewAdminConn(participant)
	t.Cleanup(func() { _ = adminConn.Close() })
	vault := cryptoadminv30.NewVaultServiceClient(adminConn)

	var kmsAPI client.AWSKMSAPI
	if sa.kmsCfg.ProtocolKeyID != "" {
		require.NotNil(t, s.KMS, "KMS registry required when ProtocolKeyID is set")
		kmsAPI = s.KMS.AWSClient()
	}

	factory := client.NewTransactionSignerFactory(s.Actors[sa.actorIndex].deps.Client, vault, sa.kmsCfg, kmsAPI)
	signer, err := factory(ctx, s.Actors[sa.actorIndex].uid, knownKeys)
	require.NoError(t, err, "build signer for actor %d", sa.actorIndex+1)

	sig, err := signer.Sign(ctx, hash)
	require.NoError(t, err, "sign hash with actor %d's protocol key", sa.actorIndex+1)

	return sig
}
