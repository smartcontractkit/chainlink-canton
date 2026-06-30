package ledger

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// DiscoverArchiveTargetsOp queries the ACS for active contracts matching the
// provided template selectors visible to the decentralized party.
var DiscoverArchiveTargetsOp = operations.NewOperation(
	"canton-ceremony/ledger/discover-archive-targets",
	semver.MustParse("1.0.0"),
	"Discover active contracts to archive from ACS template filters",
	func(b operations.Bundle, deps ContractDeployDeps, in DiscoverArchiveTargetsInput) (DiscoverArchiveTargetsOutput, error) {
		ctx := b.GetContext()

		seen := make(map[string]struct{})
		var targets []ArchiveTarget

		for _, tpl := range in.Templates {
			if tpl.ModuleName == "" || tpl.EntityName == "" {
				return DiscoverArchiveTargetsOutput{}, fmt.Errorf("template %q:%q missing module or entity", tpl.PackageName, tpl.PackageID)
			}

			var events []*apiv2.CreatedEvent
			var err error
			switch {
			case tpl.PackageName != "":
				events, err = deps.LedgerClient.GetActiveContractsByTemplateForParty(
					ctx, in.PartyID, tpl.PackageName, tpl.ModuleName, tpl.EntityName,
				)
			case tpl.PackageID != "":
				events, err = deps.LedgerClient.GetActiveContractsByTemplate(
					ctx, in.PartyID, tpl.PackageID, tpl.ModuleName, tpl.EntityName,
				)
			default:
				return DiscoverArchiveTargetsOutput{}, fmt.Errorf("template %s.%s requires package_name or package_id", tpl.ModuleName, tpl.EntityName)
			}
			if err != nil {
				return DiscoverArchiveTargetsOutput{}, fmt.Errorf("listing %s.%s: %w", tpl.ModuleName, tpl.EntityName, err)
			}

			for _, ev := range events {
				cid := ev.GetContractId()
				if cid == "" {
					continue
				}
				if _, ok := seen[cid]; ok {
					continue
				}
				seen[cid] = struct{}{}

				tid := ev.GetTemplateId()
				targets = append(targets, ArchiveTarget{
					PackageID:  tid.GetPackageId(),
					ModuleName: tid.GetModuleName(),
					EntityName: tid.GetEntityName(),
					ContractID: cid,
				})
			}
		}

		slices.SortFunc(targets, func(a, b ArchiveTarget) int {
			if c := stringsCompare(a.ModuleName, b.ModuleName); c != 0 {
				return c
			}
			if c := stringsCompare(a.EntityName, b.EntityName); c != 0 {
				return c
			}

			return stringsCompare(a.ContractID, b.ContractID)
		})

		deps.Logger.Infow("Discovered archive targets",
			"party", in.PartyID,
			"count", len(targets),
		)

		return DiscoverArchiveTargetsOutput{Targets: targets}, nil
	},
)

// PrepareArchiveBatchOp prepares one or more Archive exercises via
// InteractiveSubmissionService for a decentralized party.
var PrepareArchiveBatchOp = operations.NewOperation(
	"canton-ceremony/ledger/prepare-archive-batch",
	semver.MustParse("1.0.0"),
	"Prepare Archive exercises via InteractiveSubmissionService",
	func(b operations.Bundle, deps ContractDeployDeps, in PrepareArchiveBatchInput) (PrepareArchiveBatchOutput, error) {
		ctx := b.GetContext()

		if len(in.Targets) == 0 {
			return PrepareArchiveBatchOutput{}, fmt.Errorf("no archive targets in batch %d", in.BatchIndex)
		}

		commands := make([]*apiv2.Command, 0, len(in.Targets))
		for _, target := range in.Targets {
			commands = append(commands, &apiv2.Command{
				Command: &apiv2.Command_Exercise{
					Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{
							PackageId:  target.PackageID,
							ModuleName: target.ModuleName,
							EntityName: target.EntityName,
						},
						ContractId: target.ContractID,
						Choice:     "Archive",
						ChoiceArgument: &apiv2.Value{
							Sum: &apiv2.Value_Record{Record: &apiv2.Record{}},
						},
					},
				},
			})
		}

		resp, err := deps.LedgerClient.PrepareSubmission(
			ctx,
			commands,
			[]string{in.DecentralizedPartyID},
			[]string{in.DecentralizedPartyID},
			in.SynchronizerID,
		)
		if err != nil {
			return PrepareArchiveBatchOutput{}, fmt.Errorf("preparing archive batch %d: %w", in.BatchIndex, err)
		}

		hashHex := hex.EncodeToString(resp.GetPreparedTransactionHash())

		txBytes, err := proto.Marshal(resp.GetPreparedTransaction())
		if err != nil {
			return PrepareArchiveBatchOutput{}, fmt.Errorf("marshalling prepared transaction: %w", err)
		}

		deps.Logger.Infow("Archive batch prepared",
			"batch", in.BatchIndex,
			"hash", hashHex,
			"archives", len(in.Targets),
		)

		return PrepareArchiveBatchOutput{
			PreparedTransactionHash: hashHex,
			PreparedTxB64:           base64.StdEncoding.EncodeToString(txBytes),
			HashingSchemeVersion:    int32(resp.GetHashingSchemeVersion()),
			ArchiveCount:            len(in.Targets),
		}, nil
	},
)

func stringsCompare(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}

	return 0
}
