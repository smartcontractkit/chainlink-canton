// cleanup_staging_cv_dars archives active non-MCMS CCIP contracts and removes
// non-MCMS DARs from staging cv0–cv3 participants (DevNet endpoints).
//
// MCMS DARs (mcms-core, etc.) are left installed. Only ccip-* and standalone
// chainlink-api DARs are targeted for removal.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	stagingCCIPOwnerParty = "ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51"
	stagingBootstrapParty = "ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d"

	defaultCV0Ledger = "canton-devnet.bcy-v.metalhosts.com:443"
	defaultCV0Admin  = "admin.canton-devnet.bcy-v.metalhosts.com:443"
	defaultCV1Ledger = "devnet.cv1.bcy-v.metalhosts.com:443"
	defaultCV1Admin  = "admin.devnet.cv1.bcy-v.metalhosts.com:443"
	defaultCV2Admin  = "admin.devnet.cv2.bcy-v.metalhosts.com:443"
	defaultCV3Admin  = "admin.devnet.cv3.bcy-v.metalhosts.com:443"

	// ccipOwner is hosted on cv1/cv3; cv0 aggregate ledger JWT cannot query its ACS.
	defaultArchiveNode    = "cv1"
	defaultRMNArchiveNode = "cv3" // optional separate terminal; signatory is still ccipOwner on staging
)

var rmnRemoteTemplates = []string{
	"#ccip-rmn:CCIP.RMNRemote:RMNRemote",
}

type node struct {
	name         string
	adminTarget  string
	ledgerTarget string
	tokenEnv     string
}

func stagingNodes() []node {
	return []node{
		{name: "cv0", adminTarget: envOr("CV0_GRPC_PARTICIPANT_ADMIN", defaultCV0Admin), ledgerTarget: envOr("CV0_GRPC_LEDGER_URL", defaultCV0Ledger), tokenEnv: "CV0_TOKEN"},
		{name: "cv1", adminTarget: envOr("CV1_GRPC_PARTICIPANT_ADMIN", defaultCV1Admin), ledgerTarget: envOr("CV1_GRPC_LEDGER_URL", defaultCV1Ledger), tokenEnv: "CV1_TOKEN"},
		{name: "cv2", adminTarget: envOr("CV2_GRPC_PARTICIPANT_ADMIN", defaultCV2Admin), ledgerTarget: envOr("CV2_GRPC_LEDGER_URL", "devnet.cv2.bcy-v.metalhosts.com:443"), tokenEnv: "CV2_TOKEN"},
		{name: "cv3", adminTarget: envOr("CV3_GRPC_PARTICIPANT_ADMIN", defaultCV3Admin), ledgerTarget: envOr("CV3_GRPC_LEDGER_URL", "devnet.cv3.bcy-v.metalhosts.com:443"), tokenEnv: "CV3_TOKEN"},
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// desiredDARVersions is the target CCIP/MCMS package set (source of truth for staging).
// Any ccip-*/mcms*/chainlink-api/link DAR not matching these name+version pairs is stale.
var desiredDARVersions = map[string][]string{
	"chainlink-api":              {"2.0.0"},
	"mcms-api":                   {"1.0.0"},
	"link":                       {"2.0.0"},
	"ccip-core":                  {"2.0.0"},
	"ccip-extension-api":         {"2.0.0"},
	"ccip-runtime":               {"2.0.0"},
	"ccip-sender":                {"2.0.0"},
	"ccip-receiver":              {"2.0.0"},
	"ccip-executor":              {"2.0.0"},
	"ccip-committee-verifier":    {"2.0.0"},
	"ccip-lock-release-token-pool": {"2.0.0"},
	"ccip-burn-mint-token-pool":  {"2.0.0"},
	"ccip-factory":               {"2.0.0"},
	"mcms-core":                  {"1.0.0", "2.0.0"},
}

var nonMCMSTemplates = []string{
	"#ccip-common:CCIP.GlobalConfig:GlobalConfig",
	"#ccip-common:CCIP.RateLimiter:RateLimiter",
	"#ccip-common:CCIP.SendingMessageV1:SendingMessageV1",
	"#ccip-common:CCIP.ExecutingMessageV1:ExecutingMessageV1",
	"#ccip-common:CCIP.Tickets:TokenReceiveTicket",
	"#ccip-common:CCIP.Events:CCIPMessageSent",
	"#ccip-common:CCIP.Events:ExecutionStateChanged",
	"#ccip-common:CCIP.Events:TokenReceiveTicketClaimed",
	"#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenAdminRegistry",
	"#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenConfig",
	"#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter",
	"#ccip-offramp:CCIP.OffRamp:OffRamp",
	"#ccip-onramp:CCIP.OnRamp:OnRamp",
	"#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouterFactory",
	"#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouter",
	"#ccip-perpartyrouter:CCIP.PerPartyRouter:ArchivedExecutedMessages",
	"#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier",
	"#ccip-executor:CCIP.Executor:Executor",
	"#ccip-sender:CCIP.CCIPSender:CCIPSender",
	"#ccip-receiver:CCIP.CCIPReceiver:CCIPReceiver",
	"#ccip-receiver:CCIP.CCIPReceiver:CCIPMessageReceived",
	"#ccip-lockreleasetokenpool:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
	"#ccip-rmn:CCIP.RMNRemote:RMNRemote",
	"#ccip-factory:CCIP.Factory:CCIPFactory",
	"#mcms:MCMS.Main:MCMS",
}

func main() {
	var (
		listDars       bool
		dryRun         bool
		skipArchive    bool
		withArchive    bool
		archiveRMNOnly bool
		skipRemove     bool
		nodeFilter     string
		archiveNode    string
		rmnArchiveNode string
	)

	flag.BoolVar(&listDars, "list-dars", false, "list all installed DARs on each node (read-only)")
	flag.BoolVar(&dryRun, "dry-run", false, "list targets without archiving or removing")
	flag.BoolVar(&skipArchive, "skip-archive", false, "skip contract archive step")
	flag.BoolVar(&withArchive, "with-archive", false, "force archive even for cv0/cv2/cv3 (needs archive-node token in this shell)")
	flag.BoolVar(&archiveRMNOnly, "archive-rmn-only", false, "archive RMNRemote only (same ccipOwner signatory on staging; optional cv3 terminal)")
	flag.BoolVar(&skipRemove, "skip-remove", false, "skip DAR removal step (archive only)")
	flag.StringVar(&nodeFilter, "node", "all", "participant: cv0, cv1, cv2, cv3, or all")
	flag.StringVar(&archiveNode, "archive-node", defaultArchiveNode, "ledger for ccipOwner contract archive (default cv1)")
	flag.StringVar(&rmnArchiveNode, "rmn-archive-node", defaultRMNArchiveNode, "ledger for RMNRemote dry-run/archive when using --archive-rmn-only (default cv3)")
	flag.Parse()

	nodes, err := selectNodes(nodeFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	if listDars {
		for _, n := range nodes {
			if err := listNodeDars(n); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] list DARs: %v\n", n.name, err)
				os.Exit(1)
			}
		}
		return
	}

	runArchive := shouldRunArchive(nodeFilter, skipArchive, withArchive, archiveRMNOnly)
	if runArchive && skipArchive && !archiveRMNOnly {
		fmt.Fprintf(os.Stderr, "conflicting flags: --skip-archive and --with-archive\n")
		os.Exit(2)
	}
	if archiveRMNOnly && skipArchive {
		fmt.Fprintf(os.Stderr, "conflicting flags: --skip-archive and --archive-rmn-only\n")
		os.Exit(2)
	}

	if runArchive {
		if archiveRMNOnly {
			rmnN, err := selectNode(rmnArchiveNode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "rmn archive node: %v\n", err)
				os.Exit(2)
			}
			if err := archiveRMNRemote(rmnN, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "archive RMNRemote: %v\n", err)
				os.Exit(1)
			}
		} else {
			archiveN, err := selectNode(archiveNode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "archive node: %v\n", err)
				os.Exit(2)
			}
			if err := archiveCCIPOwnerContracts(archiveN, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "archive contracts: %v\n", err)
				os.Exit(1)
			}
			if err := tryArchiveRMNRemote(rmnArchiveNode, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "archive RMNRemote: %v\n", err)
				os.Exit(1)
			}
		}
	} else if !skipArchive && !listDars && nodeFilter != "all" {
		fmt.Printf("Note: skipping archive (--node %s). Archive once from the cv1 terminal, then remove DARs here.\n\n", nodeFilter)
	}

	if !skipRemove {
		for _, n := range nodes {
			if err := cleanupNodeDars(n, dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "[%s] cleanup DARs: %v\n", n.name, err)
				os.Exit(1)
			}
		}
	}

	if dryRun {
		fmt.Println("\nDry run complete. Re-run without --dry-run to archive and remove.")
	}
}

func shouldRunArchive(nodeFilter string, skipArchive, withArchive, archiveRMNOnly bool) bool {
	if archiveRMNOnly {
		return true
	}
	if skipArchive {
		return false
	}
	if withArchive {
		return true
	}
	// Per-terminal workflow: ccipOwner archive from cv1; RMNRemote from cv3 (--archive-rmn-only).
	switch nodeFilter {
	case "cv1", "all":
		return true
	default:
		return false
	}
}

func selectNodes(filter string) ([]node, error) {
	switch filter {
	case "all":
		return stagingNodes(), nil
	default:
		n, err := selectNode(filter)
		if err != nil {
			return nil, err
		}
		return []node{n}, nil
	}
}

func selectNode(name string) (node, error) {
	for _, n := range stagingNodes() {
		if n.name == name {
			return n, nil
		}
	}
	return node{}, fmt.Errorf("invalid node %q (want cv0, cv1, cv2, or cv3)", name)
}

func archiveCCIPOwnerContracts(n node, dryRun bool) error {
	jwt := os.Getenv(n.tokenEnv)
	if jwt == "" {
		return fmt.Errorf("%s is required for contract archive (%s ledger hosts ccipOwner).\n"+
			"  Run this from the cv1 terminal (cv1 Okta client), not cv0/cv2/cv3:\n"+
			"    go run ./scripts/cleanup_staging_cv_dars --skip-remove --node cv1", n.tokenEnv, n.name)
	}

	repoRoot, err := repoRoot()
	if err != nil {
		return err
	}

	archiveScript := filepath.Join(repoRoot, "scripts", "archive_active_canton_contracts")

	fmt.Printf("=== Archive non-MCMS CCIP contracts (%s ledger: %s) ===\n", n.name, n.ledgerTarget)
	fmt.Printf("Party: %s\n", stagingCCIPOwnerParty)
	if err := runArchive(buildArchiveArgs(archiveScript, n.ledgerTarget, jwt, stagingCCIPOwnerParty, nonMCMSTemplates, dryRun)); err != nil {
		return err
	}

	// Factory may still be owned by bootstrap party before SetOwnerToMCMS.
	fmt.Printf("\nParty: %s (factory)\n", stagingBootstrapParty)
	return runArchive(buildArchiveArgs(archiveScript, n.ledgerTarget, jwt, stagingBootstrapParty, []string{
		"#ccip-factory:CCIP.Factory:CCIPFactory",
	}, dryRun))
}

func tryArchiveRMNRemote(rmnArchiveNode string, dryRun bool) error {
	rmnN, err := selectNode(rmnArchiveNode)
	if err != nil {
		return err
	}
	if os.Getenv(rmnN.tokenEnv) == "" {
		fmt.Printf("\nNote: RMNRemote not archived (%s not set in this terminal).\n", rmnN.tokenEnv)
		fmt.Printf("  Run from the cv3 terminal:\n")
		fmt.Printf("    go run ./scripts/cleanup_staging_cv_dars --skip-remove --archive-rmn-only --node cv3\n")
		return nil
	}
	return archiveRMNRemote(rmnN, dryRun)
}

func archiveRMNRemote(n node, dryRun bool) error {
	jwt := os.Getenv(n.tokenEnv)
	if jwt == "" {
		return fmt.Errorf("%s is required for RMNRemote archive (%s hosts rmnOwner signatory).\n"+
			"  Run from the cv3 terminal:\n"+
			"    go run ./scripts/cleanup_staging_cv_dars --skip-remove --archive-rmn-only --node cv3", n.tokenEnv, n.name)
	}

	repoRoot, err := repoRoot()
	if err != nil {
		return err
	}

	archiveScript := filepath.Join(repoRoot, "scripts", "archive_active_canton_contracts")

	fmt.Printf("\n=== Archive RMNRemote (%s ledger: %s) ===\n", n.name, n.ledgerTarget)
	fmt.Printf("Query party: %s\n", stagingCCIPOwnerParty)
	fmt.Printf("Note: staging DevNet deployed RMNRemote with rmnOwner=ccipOwner (signatory is ccipOwner; multiparty archive required).\n")
	return runArchive(buildArchiveArgs(archiveScript, n.ledgerTarget, jwt, stagingCCIPOwnerParty, rmnRemoteTemplates, dryRun))
}

func buildArchiveArgs(script, ledgerURL, jwt, party string, templates []string, dryRun bool) []string {
	args := []string{"run", script,
		"--grpc-url", ledgerURL,
		"--jwt", jwt,
		"--party", party,
	}
	for _, tmpl := range templates {
		args = append(args, "--template", tmpl)
	}
	if dryRun {
		args = append(args, "--dry-run")
	}
	return args
}

func runArchive(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repo root from %s", wd)
		}
		dir = parent
	}
}

func listNodeDars(n node) error {
	jwt := os.Getenv(n.tokenEnv)
	if jwt == "" {
		return fmt.Errorf("missing %s", n.tokenEnv)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

	conn, err := grpc.NewClient(n.adminTarget, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return fmt.Errorf("dial admin %s: %w", n.adminTarget, err)
	}
	defer conn.Close()

	pkgService := participantv30.NewPackageServiceClient(conn)
	resp, err := pkgService.ListDars(ctx, &participantv30.ListDarsRequest{Limit: 500})
	if err != nil {
		return fmt.Errorf("list DARs: %w", err)
	}

	dars := resp.GetDars()
	var shown, excludedSplice, desired, keep, remove int
	fmt.Printf("\n=== %s (%s) — %d DARs (excluding splice) ===\n", n.name, n.adminTarget, len(dars))
	fmt.Printf("%-8s %-36s %-12s %s\n", "ACTION", "NAME", "VERSION", "MAIN_PACKAGE_ID")
	fmt.Println(strings.Repeat("-", 104))

	for _, dar := range dars {
		if isSpliceDar(dar.GetName()) {
			excludedSplice++
			continue
		}
		shown++
		name, version := dar.GetName(), dar.GetVersion()
		var action string
		switch {
		case isDesiredDar(name, version):
			action = "DESIRED"
			desired++
		case shouldRemoveDar(name, version):
			action = "REMOVE"
			remove++
		default:
			action = "KEEP"
			keep++
		}
		fmt.Printf("%-8s %-36s %-12s %s\n", action, name, version, dar.GetMain())
	}

	fmt.Printf("\nSummary: %d shown (%d splice hidden), %d desired, %d keep (platform/other), %d remove (stale)\n",
		shown, excludedSplice, desired, keep, remove)

	installed := map[string]map[string]struct{}{}
	for _, dar := range dars {
		if isSpliceDar(dar.GetName()) {
			continue
		}
		name, version := dar.GetName(), dar.GetVersion()
		if installed[name] == nil {
			installed[name] = map[string]struct{}{}
		}
		installed[name][version] = struct{}{}
	}
	var missing []string
	for name, versions := range desiredDARVersions {
		for _, version := range versions {
			if _, ok := installed[name][version]; !ok {
				missing = append(missing, fmt.Sprintf("%s@%s", name, version))
			}
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Printf("Missing desired: %s\n", strings.Join(missing, ", "))
	} else {
		fmt.Println("Missing desired: none")
	}
	return nil
}

func cleanupNodeDars(n node, dryRun bool) error {
	jwt := os.Getenv(n.tokenEnv)
	if jwt == "" {
		return fmt.Errorf("missing %s", n.tokenEnv)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

	conn, err := grpc.NewClient(n.adminTarget, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return fmt.Errorf("dial admin: %w", err)
	}
	defer conn.Close()

	pkgService := participantv30.NewPackageServiceClient(conn)

	resp, err := pkgService.ListDars(ctx, &participantv30.ListDarsRequest{Limit: 500})
	if err != nil {
		return fmt.Errorf("list DARs: %w", err)
	}

	var targets []*participantv30.DarDescription
	for _, dar := range resp.GetDars() {
		if shouldRemoveDar(dar.GetName(), dar.GetVersion()) {
			targets = append(targets, dar)
		}
	}

	fmt.Printf("\n=== %s: stale DARs to remove (%d) ===\n", n.name, len(targets))
	for _, dar := range targets {
		fmt.Printf("  %-32s v%-10s main=%s\n", dar.GetName(), dar.GetVersion(), dar.GetMain())
	}

	if dryRun || len(targets) == 0 {
		return nil
	}

	for _, dar := range targets {
		main := dar.GetMain()
		name := dar.GetName()

		fmt.Printf("[%s] unvet %s (%s)...\n", n.name, name, shortID(main))
		if _, err := pkgService.UnvetDar(ctx, &participantv30.UnvetDarRequest{MainPackageId: main}); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] unvet %s: %v\n", n.name, name, err)
		}

		fmt.Printf("[%s] remove DAR %s...\n", n.name, name)
		if _, err := pkgService.RemoveDar(ctx, &participantv30.RemoveDarRequest{MainPackageId: main}); err != nil {
			if darRemovalBlockedByLiveContracts(err) {
				return fmt.Errorf("remove DAR %s blocked by active contracts — run archive from the cv1 terminal first:\n"+
					"  go run ./scripts/cleanup_staging_cv_dars --skip-remove --node cv1   # cv1 terminal\n"+
					"  go run ./scripts/cleanup_staging_cv_dars --node %s                      # then this terminal", name, n.name)
			}
			return fmt.Errorf("remove DAR %s: %w", name, err)
		}

		fmt.Printf("[%s] remove package %s...\n", n.name, shortID(main))
		if _, err := pkgService.RemovePackage(ctx, &participantv30.RemovePackageRequest{
			PackageId: main,
			Force:     true,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] remove package %s: %v\n", n.name, shortID(main), err)
		}
	}

	remaining, err := pkgService.ListDars(ctx, &participantv30.ListDarsRequest{Limit: 500})
	if err != nil {
		return fmt.Errorf("verify list DARs: %w", err)
	}

	var left int
	for _, dar := range remaining.GetDars() {
		if shouldRemoveDar(dar.GetName(), dar.GetVersion()) {
			left++
		}
	}
	if left > 0 {
		return fmt.Errorf("%d stale DARs still remain on %s", left, n.name)
	}

	fmt.Printf("[%s] stale DAR cleanup complete\n", n.name)
	return nil
}

func isDesiredDar(name, version string) bool {
	for _, v := range desiredDARVersions[name] {
		if v == version {
			return true
		}
	}
	return false
}

// shouldRemoveDar reports stale chainlink packages: ccip-*/mcms*/chainlink-api/link
// that are not in desiredDARVersions.
func shouldRemoveDar(name, version string) bool {
	if isDesiredDar(name, version) {
		return false
	}
	n := strings.ToLower(name)
	if strings.HasPrefix(n, "ccip-") {
		return true
	}
	if strings.HasPrefix(n, "mcms") {
		return true
	}
	return n == "chainlink-api" || n == "link"
}

func isSpliceDar(name string) bool {
	return strings.HasPrefix(strings.ToLower(name), "splice-")
}

func darRemovalBlockedByLiveContracts(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if strings.Contains(msg, "PACKAGE_OR_DAR_REMOVAL_ERROR") || strings.Contains(msg, "is in-use") {
		return true
	}
	if s, ok := status.FromError(err); ok && s.Code() == codes.FailedPrecondition {
		return strings.Contains(strings.ToLower(msg), "cannot be removed")
	}
	return false
}

func shortID(id string) string {
	if len(id) <= 16 {
		return id
	}
	return id[:16] + "..."
}
