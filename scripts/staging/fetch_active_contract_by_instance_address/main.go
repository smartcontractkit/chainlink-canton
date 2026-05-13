// Command fetch_active_contract_by_instance_address reads a Chainlink Canton contract from the ledger.
//
// Prefer --contract-id (single RPC, fast). By instance-address hash the tool streams ACS for the
// template + party until it finds a match — slow when many contracts exist (e.g. token pools).
//
// Create arguments are JSON-encoded via go-daml (nested maps/slices; genmaps as _type+entries).
//
// JWT defaults to env ONCHAIN_CANTON_JWT_TOKEN (same as scripts/archive_active_canton_contracts).
//
// Defaults target BCY CV1 devnet: HTTPS dashboard https://devnet.cv1.bcy-v.metalhosts.com/api/validator
// maps to ledger gRPC host devnet.cv1.bcy-v.metalhosts.com:443 . Override with --grpc-url if needed.

// go run ./scripts/staging/fetch_active_contract_by_instance_address \
// --template '#ccip-common:CCIP.GlobalConfig:GlobalConfig' \
// --instance-address '{"address":"0xa95f120fc972c72e75d74c880c26ba982c60b123c74aa9e5b18e138a59e0916a"}' \
// --instance-id globalconfig-szvgb

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

const (
	defaultJWTEnv = "ONCHAIN_CANTON_JWT_TOKEN"

	// defaultDevnetValidatorHTTPS is the validator HTTP URL; ledger clients dial gRPC host:port instead.
	defaultDevnetValidatorHTTPS = "https://devnet.cv1.bcy-v.metalhosts.com/api/validator"

	defaultPartyBootstrap = "ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51"
	// defaultPartyBootstrap = "ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d"
)

var instanceAddressHexPattern = regexp.MustCompile(`0x[0-9a-fA-F]{64}`)

func main() {
	var (
		grpcURL          string
		authority        string
		jwt              string
		jwtEnv           string
		party            string
		userID           string
		insecure         bool
		timeout          time.Duration
		templateID       string
		instanceAddress  string
		instanceIDHint   string
		contractID       string
		progressEvery    int
		maxScanContracts int
		pretty           bool
		payloadOnly      bool
	)

	flag.StringVar(&grpcURL, "grpc-url", defaultDevnetValidatorHTTPS,
		"Canton ledger gRPC dial target (host:port), or https:// URL (host and port inferred; path ignored)")
	flag.StringVar(&authority, "authority", "", "Optional gRPC authority override")
	flag.StringVar(&jwt, "jwt", "", "JWT for the ledger API; defaults to the env var specified by --jwt-env")
	flag.StringVar(&jwtEnv, "jwt-env", defaultJWTEnv, "Environment variable to read JWT from when --jwt is not set")
	flag.StringVar(&party, "party", defaultPartyBootstrap,
		"Party whose ACS is queried (default: ccipBootstrapOwner on CV1 devnet)")
	flag.StringVar(&userID, "user-id", "", "Optional user ID; primary party is used when --party is omitted")
	flag.BoolVar(&insecure, "insecure", false, "Use insecureStatic auth instead of static auth")
	flag.DurationVar(&timeout, "timeout", 10*time.Minute,
		"Per-request timeout (ACS scans over many contracts may need several minutes)")
	flag.StringVar(&templateID, "template", "",
		"Template ID packageId:module:entity (required for --instance-address ACS scan; ignored with --contract-id alone)")
	flag.StringVar(&instanceAddress, "instance-address", "",
		"Hashed InstanceAddress (0x + 64 hex), or JSON/text containing that hex (omit if using --contract-id alone)")
	flag.StringVar(&contractID, "contract-id", "",
		"Canton ledger contract id (fast path: one RPC via EventQueryService, avoids scanning the whole ACS)")
	flag.StringVar(&instanceIDHint, "instance-id", "",
		"Exact DAML instanceId text before '@' in datastore labels (e.g. lockreleasetokenpool-aswyq); skips non-matching contracts cheaply")
	flag.IntVar(&progressEvery, "progress-every", 100,
		"Print ACS scan progress to stderr every N contracts (0=disable)")
	flag.IntVar(&maxScanContracts, "max-scan", 0,
		"Abort after N ACS contracts with no match (0=unlimited); use to debug endless scans")
	flag.BoolVar(&pretty, "pretty", true, "Pretty-print JSON output")
	flag.BoolVar(&payloadOnly, "payload-only", false, "Print only the create-argument map as JSON (no metadata envelope)")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Fetch one active Canton contract by InstanceAddress hash using ACS.\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Instance-address mode scans every active contract for template+party until the hash matches — often slow.\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Use --contract-id when possible (single RPC). Match datastore labels with optional --instance-id text.\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Reads JWT from --jwt or %s.\n\n", defaultJWTEnv)
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  export ONCHAIN_CANTON_JWT_TOKEN=\"$(canton-login canton-devnet)\"\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  go run ./scripts/staging/fetch_active_contract_by_instance_address --contract-id '<ledger cid>'\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  go run ./scripts/staging/fetch_active_contract_by_instance_address \\\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "    --template '#ccip-lockreleasetokenpool:CCIP.LockReleaseTokenPool:LockReleaseTokenPool' \\\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "    --instance-address '{\"address\":\"0x9771...\"}' \\\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "    --instance-id lockreleasetokenpool-aswyq --max-scan 50000\n")
	}
	flag.Parse()

	rawInput := strings.TrimSpace(instanceAddress)
	if rawInput == "" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fatalf("read stdin: %v", err)
		}
		rawInput = strings.TrimSpace(string(b))
	}
	if strings.TrimSpace(contractID) == "" && rawInput == "" {
		fatalf("provide --contract-id and/or --instance-address (or pipe address JSON / hex on stdin)")
	}

	var addr contracts.InstanceAddress
	if rawInput != "" {
		var err error
		addr, err = parseInstanceAddressInput(rawInput)
		if err != nil {
			fatalf("parse instance address: %v", err)
		}
	}

	grpcDial := ledgerGRPCDialTarget(grpcURL)
	if strings.TrimSpace(contractID) == "" && strings.TrimSpace(templateID) == "" {
		fatalf("--template is required when using --instance-address (ACS filters by template)")
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

	conn, err := grpc.NewClient(grpcDial, dialOpts...)
	if err != nil {
		fatalf("connect to ledger API: %v", err)
	}
	defer conn.Close()

	stateClient := apiv2.NewStateServiceClient(conn)
	userClient := adminv2.NewUserManagementServiceClient(conn)

	if party == "" {
		resolved, err := resolvePrimaryParty(ctx, timeout, userClient, userID)
		if err != nil {
			fatalf("resolve primary party for user %s: %v", userID, err)
		}
		party = resolved
	}

	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var active *apiv2.ActiveContract
	switch {
	case strings.TrimSpace(contractID) != "":
		ac, err := activeContractByContractID(callCtx, party, strings.TrimSpace(contractID), conn)
		if err != nil {
			fatalf("contract-id lookup: %v", err)
		}
		active = ac
		if rawInput != "" {
			if err := verifyCreatedEventMatchesInstanceAddress(active.GetCreatedEvent(), addr); err != nil {
				fatalf("instance address mismatch for contract-id: %v", err)
			}
		}
	default:
		ac, err := findActiveContractByInstanceAddress(callCtx, stateClient, party, templateID, addr, instanceIDHint, progressEvery, maxScanContracts)
		if err != nil {
			fatalf("ACS lookup: %v", err)
		}
		active = ac
	}

	created := active.GetCreatedEvent()
	var payload map[string]interface{}
	if err := ledger.RecordToStruct(created.GetCreateArguments(), &payload); err != nil {
		fatalf("decode create arguments to JSON-friendly map: %v", err)
	}

	var out any
	if payloadOnly {
		out = payload
	} else {
		tid := created.GetTemplateId()
		envelope := map[string]interface{}{
			"templateId": map[string]string{
				"packageId":  tid.GetPackageId(),
				"moduleName": tid.GetModuleName(),
				"entityName": tid.GetEntityName(),
			},
			"contractId":  created.GetContractId(),
			"signatories": created.GetSignatories(),
			"observers":   created.GetObservers(),
			"payload":     payload,
		}
		if ts := created.GetCreatedAt(); ts != nil {
			envelope["createdAt"] = ts.AsTime().Format(time.RFC3339Nano)
		}
		if sid := active.GetSynchronizerId(); sid != "" {
			envelope["synchronizerId"] = sid
		}
		out = envelope
	}

	enc := json.NewEncoder(os.Stdout)
	if pretty {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(out); err != nil {
		fatalf("encode JSON: %v", err)
	}
}

func ledgerGRPCDialTarget(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = defaultDevnetValidatorHTTPS
	}
	if !strings.Contains(raw, "://") {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return raw
	}
	host := u.Hostname()
	if port := u.Port(); port != "" {
		return net.JoinHostPort(host, port)
	}

	return net.JoinHostPort(host, "443")
}

func parseInstanceAddressInput(raw string) (contracts.InstanceAddress, error) {
	s := strings.TrimSpace(raw)
	s = strings.TrimSuffix(s, ",")

	var wrap struct {
		Address string `json:"address"`
	}
	switch err := json.Unmarshal([]byte(s), &wrap); {
	case err == nil && wrap.Address != "":
		s = strings.TrimSpace(wrap.Address)
	default:
		if m := instanceAddressHexPattern.FindString(s); m != "" {
			s = m
		}
	}

	s = strings.Trim(s, `"'`)

	if !strings.HasPrefix(s, "0x") && len(s) == 64 {
		s = "0x" + s
	}

	hexDigits := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if len(hexDigits) != 64 {
		return contracts.InstanceAddress{}, fmt.Errorf("expected 32-byte instance address (64 hex digits), got %q", s)
	}

	return contracts.HexToInstanceAddress(s), nil
}

func wildcardEventFormatForParty(party string, includeCreatedEventBlob bool) *apiv2.EventFormat {
	return &apiv2.EventFormat{
		FiltersByParty: map[string]*apiv2.Filters{
			party: {
				Cumulative: []*apiv2.CumulativeFilter{{
					IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
						WildcardFilter: &apiv2.WildcardFilter{
							IncludeCreatedEventBlob: includeCreatedEventBlob,
						},
					},
				}},
			},
		},
		Verbose: false,
	}
}

func activeContractByContractID(ctx context.Context, party, contractID string, conn *grpc.ClientConn) (*apiv2.ActiveContract, error) {
	eq := apiv2.NewEventQueryServiceClient(conn)
	resp, err := eq.GetEventsByContractId(ctx, &apiv2.GetEventsByContractIdRequest{
		ContractId:  contractID,
		EventFormat: wildcardEventFormatForParty(party, false),
	})
	if err != nil {
		return nil, err
	}
	created := resp.GetCreated()
	if created == nil || created.GetCreatedEvent() == nil {
		return nil, fmt.Errorf("no created event for contract id %s (wrong party filter, archived contract, or unknown id)", contractID)
	}

	return &apiv2.ActiveContract{
		CreatedEvent:   created.GetCreatedEvent(),
		SynchronizerId: created.GetSynchronizerId(),
	}, nil
}

func instanceIDTextFromCreateArguments(args *apiv2.Record) (string, bool) {
	return scanRecordForInstanceID(args)
}

func scanRecordForInstanceID(r *apiv2.Record) (string, bool) {
	if r == nil {
		return "", false
	}
	for _, field := range r.GetFields() {
		if field.GetLabel() == "instanceId" {
			if t := field.GetValue().GetText(); t != "" {
				return t, true
			}
		}
		if nested := field.GetValue().GetRecord(); nested != nil {
			if t, ok := scanRecordForInstanceID(nested); ok {
				return t, ok
			}
		}
	}

	return "", false
}

func partyFromCreateArgs(args *apiv2.Record, field string) (string, bool) {
	if args == nil {
		return "", false
	}
	for _, f := range args.GetFields() {
		if f.GetLabel() != field {
			continue
		}
		p := f.GetValue().GetParty()
		if p != "" {
			return p, true
		}
	}

	return "", false
}

// partyCandidatesForInstanceHash matches deployment/utils/operations/contract/deploy.go:
// RawInstanceAddress uses OwnerParty (for pools: PoolOwner). CreatedEvent signatories should match,
// but we also try poolOwner/ccipOwner fields from create arguments for robustness.
func partyCandidatesForInstanceHash(ev *apiv2.CreatedEvent) []string {
	if ev == nil {
		return nil
	}
	seen := make(map[string]struct{})
	var out []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" {
			return
		}
		if _, ok := seen[p]; ok {
			return
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}

	for _, s := range ev.GetSignatories() {
		add(s)
	}
	args := ev.GetCreateArguments()
	if po, ok := partyFromCreateArgs(args, "poolOwner"); ok {
		add(po)
	}
	if co, ok := partyFromCreateArgs(args, "ccipOwner"); ok {
		add(co)
	}

	return out
}

// createdEventMatchesInstanceAddress returns true if create-args instanceId combined with a candidate party
// yields the same hashed InstanceAddress as deployment tooling (keccak256("instanceId@party")).
func createdEventMatchesInstanceAddress(ev *apiv2.CreatedEvent, want contracts.InstanceAddress, instanceIDHint string) bool {
	if ev == nil {
		return false
	}

	idText, ok := instanceIDTextFromCreateArguments(ev.GetCreateArguments())
	if !ok {
		return false
	}
	if instanceIDHint != "" && idText != instanceIDHint {
		return false
	}

	parties := partyCandidatesForInstanceHash(ev)
	if len(parties) == 0 {
		return false
	}

	inst := contracts.InstanceID(idText)
	for _, partyStr := range parties {
		got := inst.RawInstanceAddress(types.PARTY(partyStr)).InstanceAddress()
		if got == want {
			return true
		}
	}

	return false
}

func verifyCreatedEventMatchesInstanceAddress(ev *apiv2.CreatedEvent, want contracts.InstanceAddress) error {
	idText, ok := instanceIDTextFromCreateArguments(ev.GetCreateArguments())
	if !ok {
		return fmt.Errorf("cannot read instanceId from create arguments")
	}
	if !createdEventMatchesInstanceAddress(ev, want, "") {
		return fmt.Errorf("expected instance address %s; create-args instanceId=%q partyCandidates=%v signatories=%v",
			want.String(), idText, partyCandidatesForInstanceHash(ev), ev.GetSignatories())
	}

	return nil
}

func findActiveContractByInstanceAddress(
	ctx context.Context,
	stateService apiv2.StateServiceClient,
	party, templateID string,
	instanceAddress contracts.InstanceAddress,
	instanceIDHint string,
	progressEvery int,
	maxScanContracts int,
) (*apiv2.ActiveContract, error) {
	ledgerEndResp, err := stateService.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end: %w", err)
	}

	packageID, moduleName, entityName, err := contracts.ParseTemplateIDFromString(templateID)
	if err != nil {
		return nil, fmt.Errorf("parse template ID: %w", err)
	}

	stream, err := stateService.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: ledgerEndResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				party: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{
							TemplateFilter: &apiv2.TemplateFilter{
								TemplateId: &apiv2.Identifier{
									PackageId:  packageID,
									ModuleName: moduleName,
									EntityName: entityName,
								},
								// Large blobs are unnecessary for instance-address matching and slow ACS scans.
								IncludeCreatedEventBlob: false,
							},
						},
					}},
				},
			},
			// Verbose must be true or some Canton deployments omit create_arguments from ACS events,
			// which breaks instanceId parsing and matching (see testhelpers.ListActiveContractsByTemplateId).
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active contracts: %w", err)
	}
	defer stream.CloseSend()

	fmt.Fprintf(os.Stderr, "ACS scan: streaming template matches for party=%s (hint: use --contract-id to skip this)...\n", party)
	if instanceIDHint != "" {
		fmt.Fprintf(os.Stderr, "ACS scan: filtering to instanceId=%q before address hash compare\n", instanceIDHint)
	}

	var activeContract *apiv2.ActiveContract
	examined := 0
	var sampleInstanceIDs []string
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("recv active contracts: %w", err)
		}

		c, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok {
			continue
		}
		examined++
		if maxScanContracts > 0 && examined > maxScanContracts {
			return nil, fmt.Errorf("reached --max-scan=%d without a match (wrong party/template/address, no such pool, or raise max-scan)", maxScanContracts)
		}
		if progressEvery > 0 && examined%progressEvery == 0 {
			fmt.Fprintf(os.Stderr, "ACS scan: examined %d active contract(s) for template %s...\n", examined, templateID)
		}

		if len(sampleInstanceIDs) < 25 {
			if id, ok := instanceIDTextFromCreateArguments(c.ActiveContract.GetCreatedEvent().GetCreateArguments()); ok {
				sampleInstanceIDs = append(sampleInstanceIDs, id)
			}
		}

		if !createdEventMatchesInstanceAddress(c.ActiveContract.GetCreatedEvent(), instanceAddress, instanceIDHint) {
			continue
		}

		if activeContract != nil {
			return nil, fmt.Errorf("multiple active contracts found for InstanceAddress %s", instanceAddress.String())
		}
		activeContract = c.ActiveContract
	}

	if activeContract == nil {
		hint := ""
		switch {
		case examined > 0 && len(sampleInstanceIDs) == 0:
			hint = " Could not read instanceId from streamed create arguments (unexpected Canton payload). "
		case instanceIDHint != "" && len(sampleInstanceIDs) > 0:
			hint = fmt.Sprintf(" None of the scanned contracts used instanceId=%q (sample ids on this ledger: %v). Wrong --grpc-url / environment? ", instanceIDHint, sampleInstanceIDs)
		case len(sampleInstanceIDs) > 0:
			hint = fmt.Sprintf(" Sample instanceIds on this ledger: %v — compare to datastore labels. ", sampleInstanceIDs)
		}

		detail := strings.TrimSpace(hint)
		if detail != "" {
			detail += " "
		}

		return nil, fmt.Errorf("no active contract found for InstanceAddress %s (scanned %d contract(s)). %s"+
			"Try --contract-id, fix --grpc-url for the env that owns this address_refs row, drop --instance-id, or use --party matching the pool participant",
			instanceAddress.String(), examined, detail)
	}

	return activeContract, nil
}

func resolvePrimaryParty(ctx context.Context, timeout time.Duration, userClient adminv2.UserManagementServiceClient, userID string) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := userClient.GetUser(callCtx, &adminv2.GetUserRequest{UserId: userID})
	if err != nil {
		return "", err
	}

	p := strings.TrimSpace(resp.GetUser().GetPrimaryParty())
	if p == "" {
		return "", fmt.Errorf("user %s has no primary party", userID)
	}

	return p, nil
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
