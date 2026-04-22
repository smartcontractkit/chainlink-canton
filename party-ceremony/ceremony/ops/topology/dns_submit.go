package topology

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/chainlink/canton-party-ceremony/ceremony"

	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"google.golang.org/protobuf/proto"
)

// SubmitDNSOp merges signatures from all partially-signed DNS transactions into
// one SignedTopologyTransaction and submits it via AddTransactions.
//
// Each signer independently ran SignDNSProposalOp against the same original
// proposal, so each produced a transaction with only their own signature added.
// This operation merges those per-signer Signature lists (deduped by signed_by)
// into the base transaction before submission.
var SubmitDNSOp = operations.NewOperation(
	"canton-ceremony/topology/submit-dns",
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
