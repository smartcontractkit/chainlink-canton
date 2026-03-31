package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

const defaultJWTEnv = "ONCHAIN_CANTON_JWT_TOKEN"

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return fmt.Errorf("value cannot be empty")
	}
	*m = append(*m, v)
	return nil
}

type activeContractMatch struct {
	Template   contracts.TemplateID
	ContractID string
	CreatedAt  time.Time
}

func main() {
	var (
		grpcURL   string
		authority string
		jwt       string
		jwtEnv    string
		party     string
		userID    string
		insecure  bool
		dryRun    bool
		timeout   time.Duration

		templates   multiFlag
		contractIDs multiFlag
	)

	flag.StringVar(&grpcURL, "grpc-url", "", "Canton ledger gRPC endpoint, e.g. canton-devnet.bcy-v.metalhosts.com:443")
	flag.StringVar(&authority, "authority", "", "Optional gRPC authority override")
	flag.StringVar(&jwt, "jwt", "", "JWT for the ledger API; defaults to the env var specified by --jwt-env")
	flag.StringVar(&jwtEnv, "jwt-env", defaultJWTEnv, "Environment variable to read JWT from when --jwt is not set")
	flag.StringVar(&party, "party", "", "Party to query and act as when archiving")
	flag.StringVar(&userID, "user-id", "", "Optional user ID to resolve primary party from when --party is omitted")
	flag.BoolVar(&insecure, "insecure", false, "Use insecureStatic auth instead of static auth")
	flag.BoolVar(&dryRun, "dry-run", false, "List matching active contracts without archiving them")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Per-request timeout")
	flag.Var(&templates, "template", "Template selector in packageId:Module:Entity form; repeat the flag for multiple templates")
	flag.Var(&contractIDs, "contract-id", "Optional contract ID filter; repeat the flag for multiple IDs")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  go run ./scripts/archive_active_canton_contracts --grpc-url <host:port> --party <party> --template <pkg:Module:Entity> [--template ...] [--dry-run]\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Archives every active contract matching the provided templates for a single Canton party.\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Use --dry-run first to inspect the ACS matches before archiving.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if grpcURL == "" {
		fatalf("--grpc-url is required")
	}
	if len(templates) == 0 {
		fatalf("at least one --template is required")
	}
	if jwt == "" {
		jwt = strings.TrimSpace(os.Getenv(jwtEnv))
	}
	if jwt == "" {
		fatalf("JWT is required; set --jwt or %s", jwtEnv)
	}
	if party == "" && userID == "" {
		fatalf("either --party or --user-id is required")
	}

	ctx := context.Background()
	authType := commonconfig.AuthTypeStatic
	if insecure {
		authType = commonconfig.AuthTypeInsecureStatic
	}
	provider, err := (&commonconfig.AuthConfig{Type: authType, JWT: jwt}).NewProvider(ctx)
	if err != nil {
		fatalf("build auth provider: %v", err)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(provider.TransportCredentials()),
		grpc.WithPerRPCCredentials(provider.PerRPCCredentials()),
	}
	if authority != "" {
		dialOpts = append(dialOpts, grpc.WithAuthority(authority))
	}

	conn, err := grpc.NewClient(grpcURL, dialOpts...)
	if err != nil {
		fatalf("connect to ledger API: %v", err)
	}
	defer conn.Close()

	stateClient := apiv2.NewStateServiceClient(conn)
	commandClient := apiv2.NewCommandServiceClient(conn)
	userClient := adminv2.NewUserManagementServiceClient(conn)

	if party == "" {
		resolvedParty, err := resolvePrimaryParty(ctx, timeout, userClient, userID)
		if err != nil {
			fatalf("resolve primary party for user %s: %v", userID, err)
		}
		party = resolvedParty
	}

	parsedTemplates, err := parseTemplates(templates)
	if err != nil {
		fatalf("parse templates: %v", err)
	}

	contractIDSet := make(map[string]struct{}, len(contractIDs))
	for _, cid := range contractIDs {
		contractIDSet[cid] = struct{}{}
	}

	fmt.Printf("Using party: %s\n", party)

	var totalFound int
	var totalArchived int
	for _, template := range parsedTemplates {
		matches, err := listActiveContracts(ctx, timeout, stateClient, party, template)
		if err != nil {
			fatalf("list active contracts for %s: %v", template.String(), err)
		}
		if len(contractIDSet) > 0 {
			matches = filterByContractID(matches, contractIDSet)
		}
		totalFound += len(matches)

		fmt.Printf("\nTemplate %s\n", template.String())
		if len(matches) == 0 {
			fmt.Println("  no matching active contracts")
			continue
		}

		for _, match := range matches {
			fmt.Printf("  %s created_at=%s\n", match.ContractID, match.CreatedAt.Format(time.RFC3339Nano))
		}
		if dryRun {
			continue
		}

		for _, match := range matches {
			updateID, err := archiveContract(ctx, timeout, commandClient, party, match)
			if err != nil {
				fatalf("archive %s (%s): %v", match.ContractID, match.Template.String(), err)
			}
			totalArchived++
			fmt.Printf("  archived update_id=%s\n", updateID)
		}
	}

	fmt.Printf("\nMatched %d active contracts\n", totalFound)
	if dryRun {
		fmt.Println("Dry run complete; no contracts were archived")
		return
	}
	fmt.Printf("Archived %d contracts\n", totalArchived)
}

func resolvePrimaryParty(ctx context.Context, timeout time.Duration, userClient adminv2.UserManagementServiceClient, userID string) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := userClient.GetUser(callCtx, &adminv2.GetUserRequest{UserId: userID})
	if err != nil {
		return "", err
	}

	party := strings.TrimSpace(resp.GetUser().GetPrimaryParty())
	if party == "" {
		return "", fmt.Errorf("user %s has no primary party", userID)
	}
	return party, nil
}

func parseTemplates(inputs []string) ([]contracts.TemplateID, error) {
	seen := make(map[string]struct{}, len(inputs))
	out := make([]contracts.TemplateID, 0, len(inputs))
	for _, input := range inputs {
		parts := strings.Split(input, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("template %q must have format packageId:Module:Entity", input)
		}

		tpl := contracts.TemplateID{
			PackageID:  strings.TrimSpace(parts[0]),
			ModuleName: strings.TrimSpace(parts[1]),
			EntityName: strings.TrimSpace(parts[2]),
		}
		if tpl.PackageID == "" || tpl.ModuleName == "" || tpl.EntityName == "" {
			return nil, fmt.Errorf("template %q must not contain empty components", input)
		}

		if _, ok := seen[tpl.String()]; ok {
			continue
		}
		seen[tpl.String()] = struct{}{}
		out = append(out, tpl)
	}
	return out, nil
}

func listActiveContracts(ctx context.Context, timeout time.Duration, stateClient apiv2.StateServiceClient, party string, template contracts.TemplateID) ([]activeContractMatch, error) {
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ledgerEndResp, err := stateClient.GetLedgerEnd(callCtx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end: %w", err)
	}

	stream, err := stateClient.GetActiveContracts(callCtx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
							TemplateId:              template.ToLedgerIdentifier(),
							IncludeCreatedEventBlob: true,
						}},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active contracts: %w", err)
	}
	defer stream.CloseSend()

	var matches []activeContractMatch
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("recv active contracts: %w", err)
		}

		entry, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}

		created := entry.ActiveContract.GetCreatedEvent()
		matches = append(matches, activeContractMatch{
			Template: contracts.TemplateID{
				PackageID:  created.GetTemplateId().GetPackageId(),
				ModuleName: created.GetTemplateId().GetModuleName(),
				EntityName: created.GetTemplateId().GetEntityName(),
			},
			ContractID: created.GetContractId(),
			CreatedAt:  created.GetCreatedAt().AsTime(),
		})
	}

	slices.SortFunc(matches, func(a, b activeContractMatch) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})
	return matches, nil
}

func filterByContractID(matches []activeContractMatch, contractIDs map[string]struct{}) []activeContractMatch {
	filtered := make([]activeContractMatch, 0, len(matches))
	for _, match := range matches {
		if _, ok := contractIDs[match.ContractID]; ok {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

func archiveContract(ctx context.Context, timeout time.Duration, commandClient apiv2.CommandServiceClient, party string, match activeContractMatch) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := commandClient.SubmitAndWaitForTransaction(callCtx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.New().String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: match.Template.ToLedgerIdentifier(),
					ContractId: match.ContractID,
					Choice:     "Archive",
					ChoiceArgument: &apiv2.Value{
						Sum: &apiv2.Value_Record{Record: &apiv2.Record{}},
					},
				}},
			}},
			ActAs: []string{party},
		},
	})
	if err != nil {
		return "", err
	}
	if resp.GetTransaction() == nil {
		return "", nil
	}
	return resp.GetTransaction().GetUpdateId(), nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
