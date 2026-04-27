package topology

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

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
	"canton-ceremony/topology/create-dns-proposal",
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

// CreateKickDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the kicked owner removed and a new threshold, then calls Authorize to
// create the proposal.
var CreateKickDNSProposalOp = operations.NewOperation(
	"canton-ceremony/topology/create-kick-dns-proposal",
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
			RequiredSigners:    append(append([]string{}, in.RemainingParticipants...), in.KickedParticipantID),
		}, nil
	},
)

// CreateAddDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the new owner added and threshold adjusted, then calls Authorize to
// create the proposal.
var CreateAddDNSProposalOp = operations.NewOperation(
	"canton-ceremony/topology/create-add-dns-proposal",
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
			RequiredSigners:    append([]string{}, in.ExistingParticipantUIDs...),
		}, nil
	},
)

// CreateRotationDNSProposalOp builds the updated DecentralizedNamespaceDefinition
// with the target's old namespace fingerprint replaced by the new one.
var CreateRotationDNSProposalOp = operations.NewOperation(
	"canton-ceremony/topology/create-rotation-dns-proposal",
	semver.MustParse("1.0.0"),
	"Create updated DecentralizedNamespaceDefinition proposal with rotated owner fingerprint",
	func(b operations.Bundle, deps ceremony.CantonDeps, in CreateRotationDNSProposalInput) (CreateRotationDNSProposalOutput, error) {
		if in.DecentralizedNamespace == "" {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-rotation-dns-proposal: decentralized_namespace is required"),
			)
		}
		if in.OldNamespaceFingerprint == "" || in.NewNamespaceFingerprint == "" {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("create-rotation-dns-proposal: old and new namespace fingerprints are required"),
			)
		}

		// Build new owners list by replacing the old fingerprint with the new one.
		newOwners := make([]string, 0, len(in.CurrentOwners))
		replaced := false
		for _, owner := range in.CurrentOwners {
			if owner == in.OldNamespaceFingerprint {
				newOwners = append(newOwners, in.NewNamespaceFingerprint)
				replaced = true
			} else {
				newOwners = append(newOwners, owner)
			}
		}
		if !replaced {
			return CreateRotationDNSProposalOutput{}, operations.NewUnrecoverableError(
				fmt.Errorf("create-rotation-dns-proposal: old namespace fingerprint %q not found in current owners",
					in.OldNamespaceFingerprint),
			)
		}

		ctx := b.GetContext()

		threshold := in.CurrentThreshold
		if threshold <= 0 {
			threshold = len(newOwners)/2 + 1
		}

		mapping := &protov30.TopologyMapping{
			Mapping: &protov30.TopologyMapping_DecentralizedNamespaceDefinition{
				DecentralizedNamespaceDefinition: &protov30.DecentralizedNamespaceDefinition{
					DecentralizedNamespace: in.DecentralizedNamespace,
					Threshold:              int32(threshold),
					Owners:                 newOwners,
				},
			},
		}

		tx, err := deps.Client.Authorize(ctx, uint32(in.CurrentSerial+1), mapping, in.SynchronizerID, false)
		if err != nil {
			return CreateRotationDNSProposalOutput{}, fmt.Errorf("authorizing rotation DNS proposal: %w", err)
		}

		txBytes, err := proto.Marshal(tx)
		if err != nil {
			return CreateRotationDNSProposalOutput{}, fmt.Errorf("marshalling rotation DNS proposal: %w", err)
		}
		txB64 := base64.StdEncoding.EncodeToString(txBytes)

		hash := sha256.Sum256(txBytes)
		proposalHash := fmt.Sprintf("%x", hash)

		deps.Logger.Infow("Rotation DNS proposal created",
			"namespace", in.DecentralizedNamespace,
			"threshold", threshold,
			"old_fingerprint", in.OldNamespaceFingerprint,
			"new_fingerprint", in.NewNamespaceFingerprint,
			"proposal_hash", proposalHash,
		)

		return CreateRotationDNSProposalOutput{
			DNSTxB64:           txB64,
			ProposalHashSHA256: proposalHash,
			NewOwners:          newOwners,
			RequiredSigners:    in.AllParticipantIDs,
		}, nil
	},
)
