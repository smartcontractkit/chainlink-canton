package topology

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// SignDNSProposalOp signs a DNS proposal transaction for a single participant.
// Each signer runs this independently; Canton auto-selects the signer's key.
//
// ParticipantID is included in the input so each actor gets a unique
// idempotency hash, preventing cross-actor cache collisions.
var SignDNSProposalOp = operations.NewOperation(
	"canton-ceremony/topology/sign-dns-proposal",
	semver.MustParse("1.0.0"),
	"Sign a DecentralizedNamespaceDefinition proposal for a single participant",
	func(b operations.Bundle, deps ceremony.CantonDeps, in SignDNSProposalInput) (SignDNSProposalOutput, error) {
		if in.ProposalHashSHA256 == "" || in.DNSTxB64 == "" {
			return SignDNSProposalOutput{}, operations.NewUnrecoverableError(
				errors.New("sign-dns-proposal: proposal_hash_sha256 and dns_tx_b64 are required"),
			)
		}

		ctx := b.GetContext()
		pid, err := deps.Client.GetParticipantUID(ctx)
		if err != nil {
			return SignDNSProposalOutput{}, fmt.Errorf("fetching participant UID: %w", err)
		}

		// Only the participant whose UID matches the expected signer can execute
		// this step. Mismatches are expected when iterating over all signers.
		if in.ParticipantID != pid {
			return SignDNSProposalOutput{}, fmt.Errorf(
				"participant ID mismatch: expected %s, got %s", in.ParticipantID, pid,
			)
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

		deps.Logger.Infow("DNS proposal signed",
			"participant", in.ParticipantID,
			"proposal_hash", in.ProposalHashSHA256,
		)

		return SignDNSProposalOutput{
			ParticipantID:  in.ParticipantID,
			SignedDNSTxB64: updatedB64,
			SignedBy:       in.ParticipantID,
			SignedAt:       time.Now().UTC(),
		}, nil
	},
)
