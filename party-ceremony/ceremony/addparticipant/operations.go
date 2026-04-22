package addparticipant

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/internal/helpers"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ── Step 1: GenerateNewMemberKeyOp ───────────────────────────────────────────

// GenerateNewMemberKeyOp generates namespace and DAML signing keys for the new
// participant. Only runs on the new participant's node (UID mismatch → error,
// cached by framework for other actors).
//
// Canton equivalent:
//
//	val nsKey = participant.keys.secret.generate_signing_key("ns", SigningKeyUsage.NamespaceOnly)
//	val damlKey = participant.keys.secret.generate_signing_key("daml", SigningKeyUsage.ProtocolOnly)
var GenerateNewMemberKeyOp = operations.NewOperation(
	"add-participant/canton-ceremony/generate-new-member-key",
	semver.MustParse("1.0.0"),
	"Generate namespace and DAML signing keys for the new participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in GenerateNewMemberKeyInput) (GenerateNewMemberKeyOutput, error) {
		if in.ParticipantID == "" || in.NamespaceName == "" {
			return GenerateNewMemberKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("generate-new-member-key: participant_id and namespace_name are required"),
			)
		}
		ctx := b.GetContext()

		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		if in.ParticipantID != pid {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, pid,
			)
		}

		key, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName, []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
		})
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("generating namespace key: %w", err)
		}

		damlKey, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName+"-protocol", []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
		})
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("generating protocol (DAML) signing key: %w", err)
		}

		uid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		fingerprint, err := helpers.GetPublicKeyFingerprint(key.GetPublicKey())
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("computing public key fingerprint: %w", err)
		}

		damlFingerprint, err := helpers.GetPublicKeyFingerprint(damlKey.GetPublicKey())
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("computing DAML key fingerprint: %w", err)
		}

		keyProtoBytes, err := proto.Marshal(key)
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("marshalling signing public key proto: %w", err)
		}
		keyB64 := base64.StdEncoding.EncodeToString(keyProtoBytes)

		damlKeyProtoBytes, err := proto.Marshal(damlKey)
		if err != nil {
			return GenerateNewMemberKeyOutput{}, fmt.Errorf("marshalling DAML signing key proto: %w", err)
		}
		damlKeyB64 := base64.StdEncoding.EncodeToString(damlKeyProtoBytes)

		deps.Logger.Infow("New member keys generated",
			"participant", in.ParticipantID,
			"uid", uid,
			"namespace_fingerprint", fingerprint,
			"daml_key_fingerprint", damlFingerprint,
		)

		return GenerateNewMemberKeyOutput{
			ParticipantID:        in.ParticipantID,
			ParticipantUID:       uid,
			NamespaceFingerprint: fingerprint,
			SigningKeyB64:        keyB64,
			DamlKeyB64:           damlKeyB64,
			DamlKeyFingerprint:   damlFingerprint,
		}, nil
	},
)

// ── Step 2: ProposeNewNSDOp ──────────────────────────────────────────────────

// ProposeNewNSDOp publishes a namespace delegation for the new participant's
// namespace key to the synchronizer.
//
// Canton equivalent:
//
//	participant.topology.namespace_delegations.propose_delegation(ns, nsKey, ...)
var ProposeNewNSDOp = operations.NewOperation(
	"add-participant/canton-ceremony/propose-new-nsd",
	semver.MustParse("1.0.0"),
	"Publish namespace delegation for the new participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeNewNSDInput) (ProposeNewNSDOutput, error) {
		if in.Namespace == "" {
			return ProposeNewNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-new-nsd: namespace is required"),
			)
		}
		if in.SigningKeyB64 == "" {
			return ProposeNewNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-new-nsd: signing_key_b64 is required"),
			)
		}

		ctx := b.GetContext()

		keyBytes, err := base64.StdEncoding.DecodeString(in.SigningKeyB64)
		if err != nil {
			return ProposeNewNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding signing key proto: %w", err),
			)
		}
		var pk cryptov30.SigningPublicKey
		if err := proto.Unmarshal(keyBytes, &pk); err != nil {
			return ProposeNewNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling signing key proto: %w", err),
			)
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_NamespaceDelegation{
				NamespaceDelegation: &protov30.NamespaceDelegation{
					Namespace: in.Namespace,
					TargetKey: &pk,
					Restriction: &protov30.NamespaceDelegation_CanSignAllMappings_{
						CanSignAllMappings: &protov30.NamespaceDelegation_CanSignAllMappings{},
					},
				},
			},
		}

		_, err = deps.Client.Authorize(ctx, 1, mapping, in.SynchronizerID, true)
		if err != nil {
			return ProposeNewNSDOutput{}, fmt.Errorf("proposing namespace delegation: %w", err)
		}

		deps.Logger.Infow("New participant namespace delegation proposed",
			"participant", in.ParticipantID, "namespace", in.Namespace)

		return ProposeNewNSDOutput{
			ParticipantID:      in.ParticipantID,
			DelegationProposed: true,
		}, nil
	},
)

// ── Step 3: ReadCurrentStateOp ───────────────────────────────────────────────

// ReadCurrentStateOp queries the Canton topology store for the current
// DecentralizedNamespaceDefinition and PartyToParticipant state.
// Read-only and naturally idempotent.
var ReadCurrentStateOp = operations.NewOperation(
	"add-participant/canton-ceremony/read-current-state",
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

		return ReadCurrentStateOutput{
			DecentralizedNamespace: dnsState.DecentralizedNamespace,
			DNSOwners:              dnsState.Owners,
			DNSThreshold:           dnsState.Threshold,
			DNSSerial:              dnsState.Serial,
			P2PParticipantUIDs:     participantUIDs,
			P2PThreshold:           p2pState.Threshold,
			P2PSerial:              p2pState.Serial,
		}, nil
	},
)

// ── Step 4: CreateAddDNSProposalOp ──────────────────────────────────────────

// CreateAddDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the new owner added and threshold adjusted, then calls Authorize to
// create the proposal. Serial is explicitly set to currentSerial+1.
//
// Canton equivalent:
//
//	val updatedDns = DecentralizedNamespaceDefinition.tryCreate(
//	    currentNamespace, newThreshold, currentOwners + newNamespace)
//	participant.topology.decentralized_namespaces.propose(updatedDns, store = syncId)
var CreateAddDNSProposalOp = operations.NewOperation(
	"add-participant/canton-ceremony/create-add-dns-proposal",
	semver.MustParse("1.0.0"),
	"Create updated DecentralizedNamespaceDefinition proposal with new owner added",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateAddDNSProposalInput) (CreateAddDNSProposalOutput, error) {
		if in.DecentralizedNamespace == "" {
			return CreateAddDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-add-dns-proposal: decentralized_namespace is required"),
			)
		}
		if in.NewOwnerFingerprint == "" {
			return CreateAddDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-add-dns-proposal: new_owner_fingerprint is required"),
			)
		}

		// Build new owners list with the new participant added.
		newOwners := append(append([]string{}, in.CurrentOwners...), in.NewOwnerFingerprint)

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
			return CreateAddDNSProposalOutput{}, fmt.Errorf("authorizing add DNS proposal: %w", err)
		}

		txBytes, err := proto.Marshal(tx)
		if err != nil {
			return CreateAddDNSProposalOutput{}, fmt.Errorf("marshalling add DNS proposal: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		hash := sha256.Sum256(txBytes)
		proposalHash := fmt.Sprintf("%x", hash)

		deps.Logger.Infow("Add DNS proposal created",
			"namespace", in.DecentralizedNamespace,
			"new_threshold", in.NewThreshold,
			"new_owners_count", len(newOwners),
			"proposal_hash", proposalHash,
		)

		return CreateAddDNSProposalOutput{
			DNSTxB64:           txB64,
			ProposalHashSHA256: proposalHash,
			NewOwners:          newOwners,
			NewThreshold:       in.NewThreshold,
			// Only existing members sign — the new participant is not yet an owner.
			RequiredSigners: append([]string{}, in.ExistingParticipantUIDs...),
		}, nil
	},
)

// ── Step 5: SignAddDNSProposalOp ─────────────────────────────────────────────

// SignAddDNSProposalOp signs the add DNS proposal for a single existing
// participant. Each existing participant runs this independently; Canton
// auto-selects the appropriate namespace key.
var SignAddDNSProposalOp = operations.NewOperation(
	"add-participant/canton-ceremony/sign-add-dns-proposal",
	semver.MustParse("1.0.0"),
	"Sign the add DecentralizedNamespaceDefinition proposal for a single existing participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SignAddDNSProposalInput) (SignAddDNSProposalOutput, error) {
		if in.ProposalHashSHA256 == "" || in.DNSTxB64 == "" {
			return SignAddDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("sign-add-dns-proposal: proposal_hash_sha256 and dns_tx_b64 are required"),
			)
		}

		ctx := b.GetContext()
		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return SignAddDNSProposalOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		if in.ParticipantID != pid {
			return SignAddDNSProposalOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, pid,
			)
		}

		txBytes, err := base64.StdEncoding.DecodeString(in.DNSTxB64)
		if err != nil {
			return SignAddDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding add DNS proposal: %w", err),
			)
		}

		var tx protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(txBytes, &tx); err != nil {
			return SignAddDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling add DNS proposal: %w", err),
			)
		}

		if deps.Confirmer != nil {
			detail, sErr := ceremony.SummarizeTopologyTx(in.DNSTxB64, in.ProposalHashSHA256, in.ParticipantID)
			if sErr != nil {
				return SignAddDNSProposalOutput{}, fmt.Errorf("summarizing topology tx: %w", sErr)
			}
			if cErr := deps.Confirmer.ConfirmTopologySign(ctx, detail); cErr != nil {
				return SignAddDNSProposalOutput{}, operations.NewUnrecoverableError(cErr)
			}
		}

		signed, err := deps.Client.SignTransactions(ctx, []*protov30.SignedTopologyTransaction{&tx}, in.SynchronizerID)
		if err != nil {
			return SignAddDNSProposalOutput{}, fmt.Errorf("signing add DNS proposal: %w", err)
		}
		if len(signed) == 0 {
			return SignAddDNSProposalOutput{}, fmt.Errorf("SignTransactions returned no transactions")
		}

		updatedBytes, err := proto.Marshal(signed[0])
		if err != nil {
			return SignAddDNSProposalOutput{}, fmt.Errorf("marshalling signed add DNS proposal: %w", err)
		}
		updatedB64 := base64.StdEncoding.EncodeToString(updatedBytes)

		deps.Logger.Infow("Add DNS proposal signed",
			"participant", in.ParticipantID,
			"proposal_hash", in.ProposalHashSHA256,
		)

		return SignAddDNSProposalOutput{
			ParticipantID:  in.ParticipantID,
			SignedDNSTxB64: updatedB64,
			SignedBy:       in.ParticipantID,
			SignedAt:       time.Now().UTC(),
		}, nil
	},
)

// ── Step 6: SubmitAddDNSOp ───────────────────────────────────────────────────

// SubmitAddDNSOp merges per-signer signatures from all partially-signed add
// DNS transactions into one SignedTopologyTransaction and submits it via
// AddTransactions.
var SubmitAddDNSOp = operations.NewOperation(
	"add-participant/canton-ceremony/submit-add-dns",
	semver.MustParse("1.0.0"),
	"Merge signer signatures and submit the updated DecentralizedNamespaceDefinition",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SubmitAddDNSInput) (SubmitAddDNSOutput, error) {
		if len(in.SignedDNSTxsB64) == 0 {
			return SubmitAddDNSOutput{}, operations.NewUnrecoverableError(
				errors.New("submit-add-dns: no signed transactions provided"),
			)
		}

		ctx := b.GetContext()

		firstBytes, err := base64.StdEncoding.DecodeString(in.SignedDNSTxsB64[0])
		if err != nil {
			return SubmitAddDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-add-dns: decoding base signed DNS tx: %w", err),
			)
		}

		var merged protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(firstBytes, &merged); err != nil {
			return SubmitAddDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-add-dns: unmarshalling base signed DNS tx: %w", err),
			)
		}

		// Merge signatures from all remaining partially-signed transactions.
		seen := make(map[string]struct{}, len(merged.Signatures))
		for _, sig := range merged.Signatures {
			seen[sig.SignedBy] = struct{}{}
		}
		for _, txB64 := range in.SignedDNSTxsB64[1:] {
			txBytes, decErr := base64.StdEncoding.DecodeString(txB64)
			if decErr != nil {
				return SubmitAddDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-add-dns: decoding signed DNS tx: %w", decErr),
				)
			}
			var partial protov30.SignedTopologyTransaction
			if unmarshalErr := proto.Unmarshal(txBytes, &partial); unmarshalErr != nil {
				return SubmitAddDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-add-dns: unmarshalling signed DNS tx: %w", unmarshalErr),
				)
			}
			for _, sig := range partial.Signatures {
				if _, exists := seen[sig.SignedBy]; !exists {
					seen[sig.SignedBy] = struct{}{}
					merged.Signatures = append(merged.Signatures, sig)
				}
			}
		}

		deps.Logger.Infow("Submitting merged add DNS transaction",
			"signature_count", len(merged.Signatures),
			"namespace", in.FilterNamespace,
		)

		merged.Proposal = false

		if err := deps.Client.AddTransactions(ctx, []*protov30.SignedTopologyTransaction{&merged}, in.SynchronizerID); err != nil {
			return SubmitAddDNSOutput{}, fmt.Errorf("submitting add DNS transaction: %w", err)
		}

		return SubmitAddDNSOutput{DNSSubmitted: true}, nil
	},
)

// ── Step 7: ProposeAddP2POp ──────────────────────────────────────────────────

// ProposeAddP2POp authorizes the updated PartyToParticipant mapping with the
// new participant added. Each existing participant runs this independently.
//
// Canton accumulates proposals from each existing participant. Once the
// threshold is reached, the mapping is activated.
//
// Canton equivalent:
//
//	participant.topology.party_to_participant_mappings.propose(
//	    partyId, currentParticipants :+ (newParticipant, Confirmation),
//	    confirmationThreshold, store = syncId)
var ProposeAddP2POp = operations.NewOperation(
	"add-participant/canton-ceremony/propose-add-p2p",
	semver.MustParse("1.0.0"),
	"Propose updated PartyToParticipant mapping with new participant added",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeAddP2PInput) (ProposeAddP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.AllParticipantUIDs) == 0 {
			return ProposeAddP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-add-p2p: participant_id, party_id, and all_participant_uids are required"),
			)
		}

		ctx := b.GetContext()

		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeAddP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}

		if currentUID != in.ParticipantID {
			return ProposeAddP2POutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, currentUID,
			)
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.AllParticipantUIDs))
		for i, uid := range in.AllParticipantUIDs {
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
			return ProposeAddP2POutput{}, fmt.Errorf("authorizing add P2P proposal: %w", err)
		}

		deps.Logger.Infow("Add P2P proposal submitted",
			"participant", in.ParticipantID,
			"party", in.PartyID,
			"total_participants", len(in.AllParticipantUIDs),
		)

		return ProposeAddP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)
