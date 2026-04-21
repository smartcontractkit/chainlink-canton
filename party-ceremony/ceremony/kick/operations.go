package kick

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ── Step 1: ReadCurrentStateOp ───────────────────────────────────────────────

// ReadCurrentStateOp queries the Canton topology store for the current
// DecentralizedNamespaceDefinition and PartyToParticipant state. It is the
// first step in the kick ceremony and provides the serial numbers and existing
// owner/participant lists required by all subsequent operations.
//
// This operation is read-only and naturally idempotent.
var ReadCurrentStateOp = operations.NewOperation(
	"kick/canton-ceremony/read-current-state",
	semver.MustParse("1.0.0"),
	"Query current DNS and P2P topology state from the synchronizer",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ReadCurrentStateInput) (ReadCurrentStateOutput, error) {
		if in.DecentralizedPartyID == "" {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				errors.New("read-current-state: decentralized_party_id is required"),
			)
		}

		parts := strings.SplitN(in.DecentralizedPartyID, "::", 2)
		if len(parts) != 2 || parts[1] == "" {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("read-current-state: invalid decentralized_party_id format %q: expected <prefix>::<namespace>",
					in.DecentralizedPartyID),
			)
		}
		decNS := parts[1]

		ctx := b.GetContext()

		dnsState, err := deps.Client.GetDNS(ctx, decNS, in.SynchronizerID)
		if err != nil {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("reading DNS state for namespace %q: %w", decNS, err),
			)
		}

		p2pState, err := deps.Client.GetP2P(ctx, in.DecentralizedPartyID, in.SynchronizerID)
		if err != nil {
			return ReadCurrentStateOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("reading P2P state for party %q: %w", in.DecentralizedPartyID, err),
			)
		}

		participantUIDs := make([]string, len(p2pState.Participants))
		for i, p := range p2pState.Participants {
			participantUIDs[i] = p.ParticipantUID
		}

		deps.Logger.Infow("Current state read",
			"namespace", decNS,
			"dns_owners_count", len(dnsState.Owners),
			"dns_threshold", dnsState.Threshold,
			"dns_serial", dnsState.Serial,
			"p2p_participants_count", len(participantUIDs),
			"p2p_threshold", p2pState.Threshold,
			"p2p_serial", p2pState.Serial,
		)

		out := ReadCurrentStateOutput{
			DecentralizedNamespace: dnsState.DecentralizedNamespace,
			DNSOwners:              dnsState.Owners,
			DNSThreshold:           dnsState.Threshold,
			DNSSerial:              dnsState.Serial,
			P2PParticipantUIDs:     participantUIDs,
			P2PThreshold:           p2pState.Threshold,
			P2PSerial:              p2pState.Serial,
		}
		if p2pState.PartySigningKeys != nil {
			out.PartySigningKeysB64 = p2pState.PartySigningKeys.Keys
			out.PartySigningThreshold = p2pState.PartySigningKeys.Threshold
		}

		return out, nil
	},
)

// ── Step 2: CreateKickDNSProposalOp ─────────────────────────────────────────

// CreateKickDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the kicked owner removed and a new threshold, then calls Authorize to
// create the proposal. The serial is explicitly set to currentSerial+1 to
// safely update the existing mapping.
//
// Canton equivalent:
//
//	val updatedDns = DecentralizedNamespaceDefinition.tryCreate(
//	    currentNamespace, newThreshold, currentOwners - kickedNamespace)
//	participant.topology.decentralized_namespaces.propose(updatedDns, store = syncId)
var CreateKickDNSProposalOp = operations.NewOperation(
	"kick/canton-ceremony/create-kick-dns-proposal",
	semver.MustParse("1.0.0"),
	"Create updated DecentralizedNamespaceDefinition proposal with kicked owner removed",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateKickDNSProposalInput) (CreateKickDNSProposalOutput, error) {
		if in.DecentralizedNamespace == "" {
			return CreateKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-kick-dns-proposal: decentralized_namespace is required"),
			)
		}
		if in.KickedNamespaceFingerprint == "" {
			return CreateKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-kick-dns-proposal: kicked_namespace_fingerprint is required"),
			)
		}

		// Build new owners list with the kicked participant removed.
		newOwners := make([]string, 0, len(in.CurrentOwners))
		for _, owner := range in.CurrentOwners {
			if owner != in.KickedNamespaceFingerprint {
				newOwners = append(newOwners, owner)
			}
		}
		if len(newOwners) == len(in.CurrentOwners) {
			return CreateKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("create-kick-dns-proposal: kicked namespace fingerprint %q not found in current owners",
					in.KickedNamespaceFingerprint),
			)
		}

		ctx := b.GetContext()

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
				DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
					DecentralizedNamespace: in.DecentralizedNamespace,
					Threshold:              int32(in.NewThreshold),
					Owners:                 newOwners,
				},
			},
		}

		// Serial must be explicitly incremented for updates to existing mappings.
		tx, err := deps.Client.Authorize(ctx, uint32(in.CurrentSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return CreateKickDNSProposalOutput{}, fmt.Errorf("authorizing kick DNS proposal: %w", err)
		}

		txBytes, err := proto.Marshal(tx)
		if err != nil {
			return CreateKickDNSProposalOutput{}, fmt.Errorf("marshalling kick DNS proposal: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		hash := sha256.Sum256(txBytes)
		proposalHash := fmt.Sprintf("%x", hash)

		deps.Logger.Infow("Kick DNS proposal created",
			"namespace", in.DecentralizedNamespace,
			"new_threshold", in.NewThreshold,
			"new_owners_count", len(newOwners),
			"proposal_hash", proposalHash,
		)

		return CreateKickDNSProposalOutput{
			DNSTxB64:           txB64,
			ProposalHashSHA256: proposalHash,
			NewOwners:          newOwners,
			NewThreshold:       in.NewThreshold,
			// All current owners (remaining + kicked) must sign the DNS update
			// because Canton requires threshold-of-current-owners for serial > 1.
			// The kicked participant is still a current owner until the update lands.
			RequiredSigners: append(append([]string{}, in.RemainingParticipants...), in.KickedParticipantID),
		}, nil
	},
)

// ── Step 3: SignKickDNSProposalOp ────────────────────────────────────────────

// SignKickDNSProposalOp signs the kick DNS proposal for a single remaining
// participant. Each remaining participant runs this independently; Canton
// auto-selects the appropriate namespace key.
//
// ParticipantID is included in the input so each actor gets a unique
// idempotency hash, preventing cross-actor cache collisions.
//
// Canton equivalent:
//
//	participant.topology.decentralized_namespaces.propose(updatedDns, store = syncId)
var SignKickDNSProposalOp = operations.NewOperation(
	"kick/canton-ceremony/sign-kick-dns-proposal",
	semver.MustParse("1.0.0"),
	"Sign the kick DecentralizedNamespaceDefinition proposal for a single remaining participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SignKickDNSProposalInput) (SignKickDNSProposalOutput, error) {
		if in.ProposalHashSHA256 == "" || in.DNSTxB64 == "" {
			return SignKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("sign-kick-dns-proposal: proposal_hash_sha256 and dns_tx_b64 are required"),
			)
		}

		ctx := b.GetContext()
		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return SignKickDNSProposalOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		// Only the participant whose UID matches the expected signer can execute
		// this step. Mismatches are expected when iterating over all signers.
		if in.ParticipantID != pid {
			return SignKickDNSProposalOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, pid,
			)
		}

		txBytes, err := base64.StdEncoding.DecodeString(in.DNSTxB64)
		if err != nil {
			return SignKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding kick DNS proposal: %w", err),
			)
		}

		var tx protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(txBytes, &tx); err != nil {
			return SignKickDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling kick DNS proposal: %w", err),
			)
		}

		if deps.Confirmer != nil {
			detail, sErr := ceremony.SummarizeTopologyTx(in.DNSTxB64, in.ProposalHashSHA256, in.ParticipantID)
			if sErr != nil {
				return SignKickDNSProposalOutput{}, fmt.Errorf("summarizing topology tx: %w", sErr)
			}
			if cErr := deps.Confirmer.ConfirmTopologySign(ctx, detail); cErr != nil {
				return SignKickDNSProposalOutput{}, operations.NewUnrecoverableError(cErr)
			}
		}

		signed, err := deps.Client.SignTransactions(ctx, []*protov30.SignedTopologyTransaction{&tx}, in.SynchronizerID)
		if err != nil {
			return SignKickDNSProposalOutput{}, fmt.Errorf("signing kick DNS proposal: %w", err)
		}
		if len(signed) == 0 {
			return SignKickDNSProposalOutput{}, fmt.Errorf("SignTransactions returned no transactions")
		}

		updatedBytes, err := proto.Marshal(signed[0])
		if err != nil {
			return SignKickDNSProposalOutput{}, fmt.Errorf("marshalling signed kick DNS proposal: %w", err)
		}
		updatedB64 := base64.StdEncoding.EncodeToString(updatedBytes)

		deps.Logger.Infow("Kick DNS proposal signed",
			"participant", in.ParticipantID,
			"proposal_hash", in.ProposalHashSHA256,
		)

		return SignKickDNSProposalOutput{
			ParticipantID:  in.ParticipantID,
			SignedDNSTxB64: updatedB64,
			SignedBy:       in.ParticipantID,
			SignedAt:       time.Now().UTC(),
		}, nil
	},
)

// ── Step 4: SubmitKickDNSOp ──────────────────────────────────────────────────

// SubmitKickDNSOp merges per-signer signatures from all partially-signed kick
// DNS transactions into one SignedTopologyTransaction and submits it via
// AddTransactions.
//
// Each remaining signer independently ran SignKickDNSProposalOp against the
// same original proposal, so each produced a transaction with only their own
// signature. This operation merges the Signature lists (deduped by signed_by)
// before submission.
var SubmitKickDNSOp = operations.NewOperation(
	"kick/canton-ceremony/submit-kick-dns",
	semver.MustParse("1.0.0"),
	"Merge signer signatures and submit the updated DecentralizedNamespaceDefinition",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SubmitKickDNSInput) (SubmitKickDNSOutput, error) {
		if len(in.SignedDNSTxsB64) == 0 {
			return SubmitKickDNSOutput{}, operations.NewUnrecoverableError(
				errors.New("submit-kick-dns: no signed transactions provided"),
			)
		}

		ctx := b.GetContext()

		firstBytes, err := base64.StdEncoding.DecodeString(in.SignedDNSTxsB64[0])
		if err != nil {
			return SubmitKickDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-kick-dns: decoding base signed DNS tx: %w", err),
			)
		}

		var merged protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(firstBytes, &merged); err != nil {
			return SubmitKickDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-kick-dns: unmarshalling base signed DNS tx: %w", err),
			)
		}

		// Merge signatures from all remaining partially-signed transactions.
		// Deduplicate by signed_by to avoid submitting duplicate signatures.
		seen := make(map[string]struct{}, len(merged.Signatures))
		for _, sig := range merged.Signatures {
			seen[sig.SignedBy] = struct{}{}
		}
		for _, txB64 := range in.SignedDNSTxsB64[1:] {
			txBytes, decErr := base64.StdEncoding.DecodeString(txB64)
			if decErr != nil {
				return SubmitKickDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-kick-dns: decoding signed DNS tx: %w", decErr),
				)
			}
			var partial protov30.SignedTopologyTransaction
			if unmarshalErr := proto.Unmarshal(txBytes, &partial); unmarshalErr != nil {
				return SubmitKickDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-kick-dns: unmarshalling signed DNS tx: %w", unmarshalErr),
				)
			}
			for _, sig := range partial.Signatures {
				if _, exists := seen[sig.SignedBy]; !exists {
					seen[sig.SignedBy] = struct{}{}
					merged.Signatures = append(merged.Signatures, sig)
				}
			}
		}

		deps.Logger.Infow("Submitting merged kick DNS transaction",
			"signature_count", len(merged.Signatures),
			"namespace", in.FilterNamespace,
		)

		merged.Proposal = false

		if err := deps.Client.AddTransactions(ctx, []*protov30.SignedTopologyTransaction{&merged}, in.SynchronizerID); err != nil {
			return SubmitKickDNSOutput{}, fmt.Errorf("submitting kick DNS transaction: %w", err)
		}

		return SubmitKickDNSOutput{DNSSubmitted: true}, nil
	},
)

// ── Step 5: ProposeKickP2POp ─────────────────────────────────────────────────

// ProposeKickP2POp authorizes the updated PartyToParticipant mapping with the
// kicked participant removed. Each remaining participant runs this independently.
//
// Canton accumulates proposals from each remaining participant. Once the new
// DNS threshold is reached, the mapping is activated, completing the kick.
//
// Canton equivalent:
//
//	participant.topology.party_to_participant_mappings.propose(
//	    partyId, currentParticipants.filterNot(kicked),
//	    confirmationThreshold, store = syncId)
var ProposeKickP2POp = operations.NewOperation(
	"kick/canton-ceremony/propose-kick-p2p",
	semver.MustParse("1.0.0"),
	"Propose updated PartyToParticipant mapping with kicked participant removed",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeKickP2PInput) (ProposeKickP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.RemainingParticipants) == 0 {
			return ProposeKickP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-kick-p2p: participant_id, party_id, and remaining_participants are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeKickP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}

		// Only the participant whose UID matches can propose for themselves.
		if currentUID != in.ParticipantID {
			return ProposeKickP2POutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.RemainingParticipants))
		for i, uid := range in.RemainingParticipants {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: uid,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(in.NewP2PThreshold),
					Participants: hostingParticipants,
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, uint32(in.CurrentP2PSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeKickP2POutput{}, fmt.Errorf("authorizing kick P2P proposal: %w", err)
		}

		deps.Logger.Infow("Kick P2P proposal submitted",
			"participant", in.ParticipantID,
			"party", in.PartyID,
			"remaining_count", len(in.RemainingParticipants),
		)

		return ProposeKickP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)
