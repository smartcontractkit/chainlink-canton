package onboarding

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"
	"github.com/chainlink/canton-party-ceremony/internal/helpers"

	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// ── Step 1: CreateMemberKeyOp ────────────────────────────────────────────────

// CreateMemberKeyOp generates a namespace signing key for a single participant
// and fetches the participant's UID via IdentityInitializationService.
//
// Canton equivalent:
//
//	val aliceNamespaceKey = participant1.keys.secret.generate_signing_key(
//	    "decentralized-party-namespace", SigningKeyUsage.NamespaceOnly)
var CreateMemberKeyOp = operations.NewOperation(
	"onboarding/canton-ceremony/create-member-key",
	semver.MustParse("1.0.0"),
	"Generate namespace signing key for a ceremony participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateMemberKeyInput) (CreateMemberKeyOutput, error) {
		if in.ParticipantID == "" || in.NamespaceName == "" {
			return CreateMemberKeyOutput{}, operations.NewUnrecoverableError(
				errors.New("create-member-key: participant_id and namespace_name are required"),
			)
		}
		ctx := b.GetContext()

		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("fetching client participant ID: %w", err)
		}

		if in.ParticipantID != pid {
			return CreateMemberKeyOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s", pid, in.ParticipantID)
		}

		key, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName, []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE,
		})
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("generating namespace key: %w", err)
		}

		// Also generate a PROTOCOL (DAML) signing key for the participant.
		// This key is registered in PartyToParticipant.PartySigningKeys and is
		// used to authorise DAML transactions submitted via
		// InteractiveSubmissionService — distinct from the namespace key used
		// for topology operations.
		damlKey, err := deps.Client.GenerateSigningKey(ctx, in.NamespaceName+"-protocol", []cryptov30.SigningKeyUsage{
			cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL,
		})
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("generating protocol (DAML) signing key: %w", err)
		}

		uid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		fingerprint, err := helpers.GetPublicKeyFingerprint(key.GetPublicKey())
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("computing public key fingerprint: %w", err)
		}

		damlFingerprint, err := helpers.GetPublicKeyFingerprint(damlKey.GetPublicKey())
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("computing DAML key fingerprint: %w", err)
		}

		// Preserve the complete SigningPublicKey proto (Format, Scheme, KeySpec,
		// Usage) so ProposeNamespaceDelegationOp can send it back to Canton
		// verbatim, avoiding PROTO_DESERIALIZATION_FAILURE from reconstructed keys.
		keyProtoBytes, err := proto.Marshal(key)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("marshalling signing public key proto: %w", err)
		}
		keyB64 := base64.StdEncoding.EncodeToString(keyProtoBytes)

		damlKeyProtoBytes, err := proto.Marshal(damlKey)
		if err != nil {
			return CreateMemberKeyOutput{}, fmt.Errorf("marshalling DAML signing public key proto: %w", err)
		}
		damlKeyB64 := base64.StdEncoding.EncodeToString(damlKeyProtoBytes)

		return CreateMemberKeyOutput{
			ParticipantID:        in.ParticipantID,
			ParticipantUID:       uid,
			NamespaceFingerprint: fingerprint,
			SigningKeyB64:        keyB64,
			DamlKeyB64:           damlKeyB64,
			DamlKeyFingerprint:   damlFingerprint,
		}, nil
	},
)

// ── Step 2: ProposeNamespaceDelegationOp ─────────────────────────────────────

// ProposeNamespaceDelegationOp publishes a namespace delegation for the
// participant's namespace key to the synchronizer.
//
// Canton equivalent:
//
//	participant1.topology.namespace_delegations.propose_delegation(
//	    aliceNamespace, aliceNamespaceKey, DelegationRestriction.CanSignAllMappings, store = synchronizerId)
var ProposeNamespaceDelegationOp = operations.NewOperation(
	"onboarding/canton-ceremony/propose-namespace-delegation",
	semver.MustParse("1.0.0"),
	"Publish namespace delegation to the synchronizer",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeNSDInput) (ProposeNSDOutput, error) {
		if in.Namespace == "" {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-nsd: namespace is required"),
			)
		}

		if in.SigningKeyB64 == "" {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				errors.New("propose-nsd: signing_key_b64 is required"),
			)
		}

		ctx := b.GetContext()

		keyBytes, err := base64.StdEncoding.DecodeString(in.SigningKeyB64)
		if err != nil {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding signing key proto: %w", err),
			)
		}
		var pk cryptov30.SigningPublicKey
		if err := proto.Unmarshal(keyBytes, &pk); err != nil {
			return ProposeNSDOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling signing key proto: %w", err),
			)
		}

		// Build NamespaceDelegation mapping with CanSignAllMappings restriction.
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
			return ProposeNSDOutput{}, fmt.Errorf("proposing namespace delegation: %w", err)
		}

		deps.Logger.Infow("Namespace delegation proposed",
			"participant", in.ParticipantID, "namespace", in.Namespace)

		return ProposeNSDOutput{
			ParticipantID:      in.ParticipantID,
			DelegationProposed: true,
		}, nil
	},
)

// ── Step 3: CreateDNSProposalOp ──────────────────────────────────────────────

// CreateDNSProposalOp computes the decentralized namespace, builds the
// DecentralizedNamespaceDefinition mapping, and calls Authorize to create the
// initial proposal (with the proposer's signature).
//
// Canton equivalent:
//
//	val namespaceDef = DecentralizedNamespaceDefinition.tryCreate(
//	    DecentralizedNamespaceDefinition.computeNamespace(Set(aliceNS, bobNS, charlieNS)),
//	    PositiveInt.tryCreate(2),
//	    NonEmpty(Set, aliceNS, bobNS, charlieNS))
//	participant1.topology.decentralized_namespaces.propose(namespaceDef, store = synchronizerId)
var CreateDNSProposalOp = operations.NewOperation(
	"onboarding/canton-ceremony/create-dns-proposal",
	semver.MustParse("1.0.0"),
	"Create DecentralizedNamespaceDefinition proposal with first signature",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateDNSProposalInput) (CreateDNSProposalOutput, error) {
		if len(in.Members) == 0 {
			return CreateDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-dns-proposal: at least one member is required"),
			)
		}

		ctx := b.GetContext()

		owners := make([]string, len(in.Members))
		for i, m := range in.Members {
			owners[i] = m.NamespaceFingerprint
		}
		decNS := helpers.ComputeDecentralizedNamespace(owners)

		threshold := in.Threshold
		if threshold <= 0 {
			threshold = len(owners)/2 + 1
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
				DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
					DecentralizedNamespace: decNS,
					Threshold:              int32(threshold),
					Owners:                 owners,
				},
			},
		}

		// Authorize with mustFullyAuthorize=false: creates proposal with proposer's signature.
		tx, err := deps.Client.Authorize(ctx, 1, mapping, in.SynchronizerID, false)
		if err != nil {
			return CreateDNSProposalOutput{}, fmt.Errorf("authorizing DNS proposal: %w", err)
		}

		txBytes, err := proto.Marshal(tx)
		if err != nil {
			return CreateDNSProposalOutput{}, fmt.Errorf("marshalling DNS proposal: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		hash := sha256.Sum256(txBytes)
		proposalHash := fmt.Sprintf("%x", hash)

		// For initial DNS creation ALL owners must sign (not just threshold).
		requiredSigners := make([]string, len(in.Members))
		for i, m := range in.Members {
			requiredSigners[i] = m.ParticipantID
		}

		deps.Logger.Infow("DNS proposal created",
			"decentralized_namespace", decNS,
			"threshold", threshold,
			"proposal_hash", proposalHash,
			"required_signers", requiredSigners)

		return CreateDNSProposalOutput{
			DecentralizedNS:    decNS,
			ProposalHashSHA256: proposalHash,
			DNSTxB64:           txB64,
			RequiredSigners:    requiredSigners,
			Threshold:          len(in.Members), // All must sign for initial DNS.
		}, nil
	},
)

// ── Step 4: SignDNSProposalOp ────────────────────────────────────────────────

// SignDNSProposalOp signs the DNS proposal transaction for a single participant.
// Each signer runs this independently; Canton auto-selects the signer's key.
//
// Canton equivalent:
//
//	participant2.topology.decentralized_namespaces.propose(namespaceDef, store = synchronizerId)
var SignDNSProposalOp = operations.NewOperation(
	"onboarding/canton-ceremony/sign-dns-proposal",
	semver.MustParse("1.0.0"),
	"Sign the DecentralizedNamespaceDefinition proposal for a single participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SignDNSProposalInput) (SignDNSProposalOutput, error) {
		if in.ProposalHashSHA256 == "" || in.DNSTxB64 == "" {
			return SignDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("sign-dns-proposal: proposal_hash and dns_tx_b64 are required"),
			)
		}

		ctx := b.GetContext()
		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return SignDNSProposalOutput{}, fmt.Errorf("fetching client participant ID: %w", err)
		}

		if in.ParticipantID != pid {
			return SignDNSProposalOutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s", pid, in.ParticipantID)
		}

		txBytes, err := base64.StdEncoding.DecodeString(in.DNSTxB64)
		if err != nil {
			return SignDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("decoding DNS proposal: %w", err),
			)
		}

		var tx protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(txBytes, &tx); err != nil {
			return SignDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("unmarshalling DNS proposal: %w", err),
			)
		}

		if deps.Confirmer != nil {
			detail, sErr := ceremony.SummarizeTopologyTx(in.DNSTxB64, in.ProposalHashSHA256, in.ParticipantID)
			if sErr != nil {
				return SignDNSProposalOutput{}, fmt.Errorf("summarizing topology tx: %w", sErr)
			}
			if cErr := deps.Confirmer.ConfirmTopologySign(ctx, detail); cErr != nil {
				return SignDNSProposalOutput{}, operations.NewUnrecoverableError(cErr)
			}
		}

		signed, err := deps.Client.SignTransactions(ctx, []*protov30.SignedTopologyTransaction{&tx}, in.SynchronizerID)
		if err != nil {
			return SignDNSProposalOutput{}, fmt.Errorf("signing DNS proposal: %w", err)
		}
		if len(signed) == 0 {
			return SignDNSProposalOutput{}, fmt.Errorf("SignTransactions returned no transactions")
		}

		updatedBytes, err := proto.Marshal(signed[0])
		if err != nil {
			return SignDNSProposalOutput{}, fmt.Errorf("marshalling signed DNS proposal: %w", err)
		}
		updatedB64 := base64.StdEncoding.EncodeToString(updatedBytes)

		deps.Logger.Infow("DNS proposal signed", "participant", in.ParticipantID, "proposal_hash", in.ProposalHashSHA256)

		return SignDNSProposalOutput{
			ParticipantID:  in.ParticipantID,
			SignedDNSTxB64: updatedB64,
			SignedBy:       in.ParticipantID,
			SignedAt:       time.Now().UTC(),
		}, nil
	},
)

// ── Step 5: SubmitDNSOp ──────────────────────────────────────────────────────

// SubmitDNSOp merges signatures from all partially-signed DNS transactions into
// one SignedTopologyTransaction and submits it via AddTransactions, then polls
// until the DecentralizedNamespaceDefinition is confirmed at head state.
//
// Each signer independently ran SignDNSProposalOp against the same original
// proposal, so each produced a transaction with only their own signature added.
// This operation merges those per-signer Signature lists (deduped by signed_by)
// into the base transaction before submission.
var SubmitDNSOp = operations.NewOperation(
	"onboarding/canton-ceremony/submit-dns",
	semver.MustParse("1.0.0"),
	"Merge signer signatures and submit the DecentralizedNamespaceDefinition",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SubmitDNSInput) (SubmitDNSOutput, error) {
		if len(in.SignedDNSTxsB64) == 0 {
			return SubmitDNSOutput{}, operations.NewUnrecoverableError(
				errors.New("submit-dns: no signed transactions provided"),
			)
		}

		ctx := b.GetContext()

		// Decode and unmarshal the first transaction as the base.
		firstBytes, err := base64.StdEncoding.DecodeString(in.SignedDNSTxsB64[0])
		if err != nil {
			return SubmitDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-dns: decoding base signed DNS tx: %w", err),
			)
		}
		var merged protov30.SignedTopologyTransaction
		if err := proto.Unmarshal(firstBytes, &merged); err != nil {
			return SubmitDNSOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("submit-dns: unmarshalling base signed DNS tx: %w", err),
			)
		}

		// Merge signatures from all remaining partially-signed transactions.
		// Deduplicate by signed_by (key fingerprint) to avoid duplicate sigs.
		seen := make(map[string]struct{}, len(merged.Signatures))
		for _, sig := range merged.Signatures {
			seen[sig.SignedBy] = struct{}{}
		}
		for _, txB64 := range in.SignedDNSTxsB64[1:] {
			txBytes, decErr := base64.StdEncoding.DecodeString(txB64)
			if decErr != nil {
				return SubmitDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-dns: decoding signed DNS tx: %w", decErr),
				)
			}
			var partial protov30.SignedTopologyTransaction
			if unmarshalErr := proto.Unmarshal(txBytes, &partial); unmarshalErr != nil {
				return SubmitDNSOutput{}, operations.NewUnrecoverableError(
					fmt.Errorf("submit-dns: unmarshalling signed DNS tx: %w", unmarshalErr),
				)
			}
			for _, sig := range partial.Signatures {
				if _, exists := seen[sig.SignedBy]; !exists {
					seen[sig.SignedBy] = struct{}{}
					merged.Signatures = append(merged.Signatures, sig)
				}
			}
		}

		deps.Logger.Infow("Submitting merged DNS transaction",
			"signature_count", len(merged.Signatures),
			"namespace", in.FilterNamespace)

		// Mark as fully authorized (not a proposal anymore).
		merged.Proposal = false

		if err := deps.Client.AddTransactions(ctx, []*protov30.SignedTopologyTransaction{&merged}, in.SynchronizerID); err != nil {
			return SubmitDNSOutput{}, fmt.Errorf("submitting DNS transaction: %w", err)
		}

		return SubmitDNSOutput{DNSSubmitted: true}, nil
	},
)

// ── Step 6: ProposeP2POp ─────────────────────────────────────────────────────

// ProposeP2POp authorizes the PartyToParticipant mapping for a single
// participant. Each participant runs this independently.
//
// ParticipantID in the input gives each participant a unique idempotency hash,
// mirroring the SignDNSProposalOp pattern so the framework cache does not
// prevent subsequent actors from executing their own proposal.
//
// Canton accumulates proposals from each participant automatically; once the
// threshold is reached, the mapping is activated.
var ProposeP2POp = operations.NewOperation(
	"onboarding/canton-ceremony/propose-p2p",
	semver.MustParse("1.0.0"),
	"Propose PartyToParticipant mapping for a single participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in ProposeP2PInput) (ProposeP2POutput, error) {
		if in.ParticipantID == "" || in.PartyID == "" || len(in.Members) == 0 {
			return ProposeP2POutput{}, operations.NewUnrecoverableError(
				errors.New("propose-p2p: participant_id, party_id, and members are required"),
			)
		}

		ctx := b.GetContext()

		// This operation can only run on the participant that owns the signing key.
		currentUID, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return ProposeP2POutput{}, fmt.Errorf("getting participant UID: %w", err)
		}
		if currentUID != in.ParticipantID {
			return ProposeP2POutput{}, fmt.Errorf("participant ID mismatch: expected %s, got %s",
				in.ParticipantID, currentUID)
		}

		threshold := in.Threshold
		if threshold <= 0 {
			threshold = len(in.Members)/2 + 1
		}

		hostingParticipants := make([]*protov30.PartyToParticipant_HostingParticipant, len(in.Members))
		for i, m := range in.Members {
			hostingParticipants[i] = &protov30.PartyToParticipant_HostingParticipant{
				ParticipantUid: m.ParticipantUID,
				Permission:     protov30.Enums_PARTICIPANT_PERMISSION_CONFIRMATION,
			}
		}

		// Collect each member's PROTOCOL (DAML) signing key so Canton can
		// verify transaction signatures submitted via InteractiveSubmissionService.
		damlKeys := make([]*cryptov30.SigningPublicKey, 0, len(in.Members))
		for _, m := range in.Members {
			if m.DamlKeyB64 == "" {
				continue
			}
			damlKeyBytes, decErr := base64.StdEncoding.DecodeString(m.DamlKeyB64)
			if decErr != nil {
				return ProposeP2POutput{}, fmt.Errorf("decoding DAML key for %s: %w", m.ParticipantID, decErr)
			}
			var damlKey cryptov30.SigningPublicKey
			if unmErr := proto.Unmarshal(damlKeyBytes, &damlKey); unmErr != nil {
				return ProposeP2POutput{}, fmt.Errorf("unmarshalling DAML key for %s: %w", m.ParticipantID, unmErr)
			}
			damlKeys = append(damlKeys, &damlKey)
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_PartyToParticipant{
				PartyToParticipant: &protov30.PartyToParticipant{
					Party:        in.PartyID,
					Threshold:    uint32(threshold),
					Participants: hostingParticipants,
					PartySigningKeys: &cryptov30.SigningKeysWithThreshold{
						Keys:      damlKeys,
						Threshold: uint32(threshold),
					},
				},
			},
		}

		// mustFullyAuthorize=false: Canton stores the proposal and accumulates
		// signatures from other participants as they also propose.
		_, err = deps.Client.Authorize(ctx, 1, mapping, in.SynchronizerID, false)
		if err != nil {
			return ProposeP2POutput{}, fmt.Errorf("authorizing P2P mapping: %w", err)
		}

		deps.Logger.Infow("P2P proposal submitted",
			"participant", in.ParticipantID, "party", in.PartyID)

		return ProposeP2POutput{
			ParticipantID: in.ParticipantID,
			Proposed:      true,
		}, nil
	},
)
