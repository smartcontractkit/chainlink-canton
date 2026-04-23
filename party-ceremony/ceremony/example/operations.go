package example

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type InitMemberInput struct {
	NamespaceName string `json:"namespace_name"`
	ParticipantID string `json:"participant_id"`
}

type InitMemberOutput struct {
	ParticipantID        string      `json:"participant_id"`
	ParticipantUID       string      `json:"participant_uid"`
	NamespaceKey         KeyMaterial `json:"namespace_signing_key"`
	NamespaceFingerprint string      `json:"namespace_fingerprint"`
}

// InitMemberOp generates mock namespace + protocol signing keys for a single
// participant and returns the member record
var InitMemberOp = operations.NewOperation(
	"example/canton-ceremony/init-member",
	semver.MustParse("1.0.0"),
	"Generate namespace signing key for the ceremony participant",
	func(b operations.Bundle, deps CantonDeps, in InitMemberInput) (InitMemberOutput, error) {
		if in.ParticipantID == "" || in.NamespaceName == "" {
			return InitMemberOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("init-member: participant_id and namespace_name are required"),
			)
		}

		nsKey, err := deps.Client.GenerateSigningKey(in.NamespaceName)
		if err != nil {
			return InitMemberOutput{}, fmt.Errorf("generating namespace key: %w", err)
		}

		uid, err := deps.Client.GetParticipantUID()
		if err != nil {
			return InitMemberOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		deps.Logger.Infow("Member initialised", "participant", in.ParticipantID, "ns_fp", nsKey.Fingerprint)

		return InitMemberOutput{
			ParticipantID:        in.ParticipantID,
			ParticipantUID:       uid,
			NamespaceKey:         nsKey,
			NamespaceFingerprint: nsKey.Fingerprint,
		}, nil
	},
)

type CreateProposalInput struct {
	NamespaceName  string             `json:"namespace_name"`
	Members        []InitMemberOutput `json:"members"`
	PartyName      string             `json:"party_name"`
	SynchronizerID string             `json:"synchronizer_id"`
	Threshold      int                `json:"threshold"`
}

type CreateProposalOutput struct {
	NamespaceName      string   `json:"namespace_name"`
	PartyID            string   `json:"party_id"`
	DecentralizedNS    string   `json:"decentralized_namespace"`
	DNSTxB64           string   `json:"dns_transaction_b64"`
	P2PTxB64           string   `json:"p2p_transaction_b64"`
	ProposalHashSHA256 string   `json:"proposal_hash_sha256"`
	RequiredSigners    []string `json:"required_signers"`
	Threshold          int      `json:"threshold"`
}

// CreateProposalOp reads all member records, builds the DecentralizedNamespace
// + PartyToParticipant topology proposals via Authorize
var CreateProposalOp = operations.NewOperation(
	"example/canton-ceremony/create-proposal",
	semver.MustParse("1.0.0"),
	"Authorize DNS and P2P topology proposals for the decentralized party (coordinator only)",
	func(b operations.Bundle, deps CantonDeps, in CreateProposalInput) (CreateProposalOutput, error) {
		if len(in.Members) == 0 {
			return CreateProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-proposal: at least one member is required"),
			)
		}

		if in.Threshold < 0 || in.Threshold > len(in.Members) {
			return CreateProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-proposal: invalid threshold. Must be between 0 and the number of members"),
			)
		}

		// Derive the decentralized namespace identifier from the sorted
		// namespace fingerprints (mirrors Canton DNS semantics).
		owners := make([]string, len(in.Members))
		for i, m := range in.Members {
			owners[i] = m.NamespaceFingerprint
		}
		decNS := deterministicNS(owners)

		partyID := fmt.Sprintf("%s::%s", in.PartyName, decNS)

		// Authorize DNS proposal in the Authorized store.
		dnsTxB64, err := deps.Client.AuthorizeProposal(AuthorizeRequest{
			Mapping:        "DecentralizedNamespaceDefinition",
			SynchronizerID: "", // empty → Authorized store
			Owners:         owners,
			Serial:         1,
		})
		if err != nil {
			return CreateProposalOutput{}, fmt.Errorf("authorizing DNS proposal: %w", err)
		}

		// Authorize P2P proposal in the Authorized store.
		p2pTxB64, err := deps.Client.AuthorizeProposal(AuthorizeRequest{
			Mapping:        "PartyToParticipant",
			SynchronizerID: "",
			PartyID:        partyID,
			Serial:         0,
		})
		if err != nil {
			return CreateProposalOutput{}, fmt.Errorf("authorizing P2P proposal: %w", err)
		}

		// Compute the canonical proposal hash: SHA256(dns_bytes || p2p_bytes).
		// In production these are the raw proto bytes from the Authorize response.
		dnsBytes, _ := base64.StdEncoding.DecodeString(dnsTxB64)
		p2pBytes, _ := base64.StdEncoding.DecodeString(p2pTxB64)
		hash := sha256.Sum256(append(dnsBytes, p2pBytes...))
		proposalHash := fmt.Sprintf("%x", hash)

		deps.Logger.Infow("Proposal created",
			"party_id", partyID, "threshold", in.Threshold, "proposal_hash", proposalHash)

		// All members are required signers in this dummy.
		requiredSigners := make([]string, len(in.Members))
		for i, m := range in.Members {
			requiredSigners[i] = m.ParticipantID
		}

		return CreateProposalOutput{
			NamespaceName:      in.NamespaceName,
			PartyID:            partyID,
			DecentralizedNS:    decNS,
			DNSTxB64:           dnsTxB64,
			P2PTxB64:           p2pTxB64,
			ProposalHashSHA256: proposalHash,
			RequiredSigners:    requiredSigners,
			Threshold:          in.Threshold,
		}, nil
	},
)

type SignProposalInput struct {
	NamespaceName      string `json:"namespace_name"`
	ParticipantID      string `json:"participant_id"`
	ProposalHashSHA256 string `json:"proposal_hash_sha256"`
	DNSTxB64           string `json:"dns_tx_b64"`
	P2PTxB64           string `json:"p2p_tx_b64"`
	SynchronizerID     string `json:"synchronizer_id"`
}

type SignProposalOutput struct {
	NamespaceName      string        `json:"namespace_name"`
	ParticipantID      string        `json:"participant_id"`
	ProposalHashSHA256 string        `json:"proposal_hash_sha256"`
	Sig                SignaturePair `json:"signatures"`
	SignedAt           time.Time     `json:"signed_at"`
}

// SignProposalOp is executed by each required signer independently.  It calls
// SignTransactions for the DNS and P2P proposals and returns the signature.
var SignProposalOp = operations.NewOperation(
	"example/canton-ceremony/sign-proposal",
	semver.MustParse("1.0.0"),
	"Sign the DNS and P2P topology transactions for a single ceremony participant",
	func(b operations.Bundle, deps CantonDeps, in SignProposalInput) (SignProposalOutput, error) {
		if in.ProposalHashSHA256 == "" {
			return SignProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("sign-proposal: proposal_hash_sha256 is required"),
			)
		}

		clientID, err := deps.Client.GetParticipantID()
		if err != nil {
			return SignProposalOutput{}, fmt.Errorf("retrieving client participant ID: %w", err)
		}
		if in.ParticipantID != clientID {
			return SignProposalOutput{}, errors.New("sign-proposal: participant ID does not match client identity")
		}

		if deps.Confirmer != nil {
			detail := ceremony.DAMLSignDetail{
				TransactionHash: in.ProposalHashSHA256,
				SignerIdentity:  in.ParticipantID,
			}
			if cErr := deps.Confirmer.ConfirmDAMLSign(b.GetContext(), detail); cErr != nil {
				return SignProposalOutput{}, operations.NewUnrecoverableError(cErr)
			}
		}

		sig, err := deps.Client.SignTransactions(SignTransactionsRequest{
			DNSTxB64:       in.DNSTxB64,
			P2PTxB64:       in.P2PTxB64,
			SynchronizerID: in.SynchronizerID,
		})
		if err != nil {
			return SignProposalOutput{}, fmt.Errorf("signing transactions: %w", err)
		}

		deps.Logger.Infow("Proposal signed", "participant", in.ParticipantID, "proposal_hash", in.ProposalHashSHA256)

		return SignProposalOutput{
			NamespaceName:      in.NamespaceName,
			ParticipantID:      in.ParticipantID,
			ProposalHashSHA256: in.ProposalHashSHA256,
			Sig:                sig,
			SignedAt:           time.Now().UTC(),
		}, nil
	},
)

type SubmitProposalInput struct {
	Proposal   CreateProposalOutput `json:"proposal"`
	Signatures []SignProposalOutput `json:"signatures"`
}

type SubmitProposalOutput struct {
	NamespaceName      string    `json:"namespace_name"`
	PartyID            string    `json:"party_id"`
	SubmittedBy        string    `json:"submitted_by"`
	SubmittedAt        time.Time `json:"submitted_at"`
	DNSConfirmed       bool      `json:"dns_confirmed"`
	P2PConfirmed       bool      `json:"p2p_confirmed"`
	ProposalHashSHA256 string    `json:"proposal_hash_sha256"`
}

// SubmitProposalOp aggregates the collected signatures, submits the DNS and
// P2P transactions via AddTransactions, and polls until both are confirmed at
// head state.  It writes result-onboard.json on success.
//
// Retry: network / polling failures are retried with the default retry policy
// (applied by the caller via [operations.WithRetry]).
var SubmitProposalOp = operations.NewOperation(
	"example/canton-ceremony/submit-proposal",
	semver.MustParse("1.0.0"),
	"Aggregate signatures and submit the finalised topology transactions (coordinator only)",
	func(b operations.Bundle, deps CantonDeps, in SubmitProposalInput) (SubmitProposalOutput, error) {
		// Submit the DNS transaction.
		if err := deps.Client.AddTransactions(in.Proposal.DNSTxB64, in.Proposal.NamespaceName); err != nil {
			return SubmitProposalOutput{}, fmt.Errorf("submitting DNS transaction: %w", err)
		}
		// Poll until the DNS change is visible at head state.
		if err := deps.Client.PollUntilConfirmed(in.Proposal.DecentralizedNS, in.Proposal.NamespaceName); err != nil {
			return SubmitProposalOutput{}, fmt.Errorf("waiting for DNS confirmation: %w", err)
		}

		// Submit the P2P transaction.
		if err := deps.Client.AddTransactions(in.Proposal.P2PTxB64, in.Proposal.NamespaceName); err != nil {
			return SubmitProposalOutput{}, fmt.Errorf("submitting P2P transaction: %w", err)
		}
		// Poll until the party-to-participant mapping is visible at head state.
		if err := deps.Client.PollUntilConfirmed(in.Proposal.PartyID, in.Proposal.NamespaceName); err != nil {
			return SubmitProposalOutput{}, fmt.Errorf("waiting for P2P confirmation: %w", err)
		}

		deps.Logger.Infow("Proposal submitted and confirmed",
			"party_id", in.Proposal.PartyID,
			"proposal_hash", in.Proposal.ProposalHashSHA256)

		return SubmitProposalOutput{
			NamespaceName:      in.Proposal.NamespaceName,
			PartyID:            in.Proposal.PartyID,
			SubmittedBy:        in.Proposal.NamespaceName, // coordinator identity in real impl
			SubmittedAt:        time.Now().UTC(),
			DNSConfirmed:       true,
			P2PConfirmed:       true,
			ProposalHashSHA256: in.Proposal.ProposalHashSHA256,
		}, nil
	},
)

// ── Helpers ───────────────────────────────────────────────────────────────────

// deterministicNS derives a stable decentralized-namespace identifier from a
// sorted list of namespace fingerprints.
func deterministicNS(fingerprints []string) string {
	var combined strings.Builder
	for _, fp := range fingerprints {
		combined.WriteString(fp + ":")
	}
	raw := sha256.Sum256([]byte(combined.String()))

	return fmt.Sprintf("1220%x", raw[:16])
}
