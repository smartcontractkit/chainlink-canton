// Upload dev DARs (*-current.dar under contracts/dars/current) to a local Docker Compose Canton participant.
//
// Usage (from chainlink-canton-fcr repo root):
//
//	export ONCHAIN_CANTON_JWT_TOKEN="<jwt>"  # or rely on jwt in -config TOML
//	go run ./scripts/upload-localnet-dars
//	go run ./scripts/upload-localnet-dars -all-participants
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"

	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

func main() {
	configPath := flag.String("config", "integration-tests/local-docker-compose.toml", "participant TOML (chain_selector + participants)")
	participantIdx := flag.Int("participant", 0, "participant index in TOML (0 = participant1)")
	allParticipants := flag.Bool("all-participants", false, "upload to every participant in the TOML")
	repoRoot := flag.String("repo", ".", "chainlink-canton-fcr repo root (contracts/dars and contracts/dependencies)")
	flag.Parse()

	jwt := strings.TrimSpace(os.Getenv("ONCHAIN_CANTON_JWT_TOKEN"))

	darPaths, err := collectDARs(*repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	if len(darPaths) == 0 {
		log.Fatal("no *-current.dar files found under contracts/dars/current (run make compile-contracts)")
	}

	input, err := loadParticipantInput(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	indices := []int{*participantIdx}
	if *allParticipants {
		indices = make([]int, len(input.Participants))
		for i := range input.Participants {
			indices[i] = i
		}
	}

	ctx := context.Background()
	for _, idx := range indices {
		if idx < 0 || idx >= len(input.Participants) {
			log.Fatalf("participant index %d out of range (have %d participants)", idx, len(input.Participants))
		}
		cfg := input.Participants[idx]
		token := jwt
		if token == "" {
			token = cfg.JWT
		}
		if token == "" {
			log.Fatal("JWT required: set ONCHAIN_CANTON_JWT_TOKEN or jwt in config TOML")
		}

		chain, err := buildChain(ctx, input.Selector, cfg, token)
		if err != nil {
			log.Fatalf("participant %d (%s): %v", idx, nameOrDefault(cfg, idx), err)
		}
		participant := chain.Participants[0]
		if cfg.Party != "" {
			log.Printf("participant %d (%s) party_id=%s", idx, nameOrDefault(cfg, idx), cfg.Party)
		} else {
			log.Printf("participant %d (%s) primary party=%s", idx, nameOrDefault(cfg, idx), participant.PartyID)
		}

		for _, darPath := range darPaths {
			darBytes, err := os.ReadFile(darPath) //nolint:gosec // local dev script, user-controlled repo path
			if err != nil {
				log.Fatalf("read %s: %v", darPath, err)
			}
			res, err := participant.AdminServices.Package.UploadDar(ctx, &participantv30.UploadDarRequest{
				Dars: []*participantv30.UploadDarRequest_UploadDarData{{Bytes: darBytes}},
				VetAllPackages:     true,
				SynchronizeVetting: true,
			})
			if err != nil {
				if s, ok := status.FromError(err); ok {
					log.Printf("gRPC code=%s message=%s", s.Code(), s.Message())
				}
				log.Fatalf("upload %s to %s: %v", filepath.Base(darPath), nameOrDefault(cfg, idx), err)
			}
			log.Printf("uploaded %s -> dar_ids=%v", filepath.Base(darPath), res.GetDarIds())
		}
	}

	log.Printf("done: uploaded %d DAR(s)", len(darPaths))
}

func collectDARs(repoRoot string) ([]string, error) {
	// Dev DARs aligned with bindings/generated/latest (see contracts/README.md).
	dirs := []string{
		filepath.Join(repoRoot, "contracts", "dars", "current"),
	}
	var paths []string
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(name, "-current.dar") {
				paths = append(paths, filepath.Join(dir, name))
			}
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func loadParticipantInput(path string) (testhelpers.ParticipantInput, error) {
	content, err := os.ReadFile(path) //nolint:gosec // config path from flag
	if err != nil {
		return testhelpers.ParticipantInput{}, err
	}
	var input testhelpers.ParticipantInput
	if err := toml.Unmarshal(content, &input); err != nil {
		return testhelpers.ParticipantInput{}, err
	}
	if input.Selector == 0 {
		input.Selector = chainsel.CANTON_LOCALNET.Selector
	}
	return input, nil
}

func buildChain(ctx context.Context, selector uint64, cfg testhelpers.ParticipantConfig, jwt string) (*canton.Chain, error) {
	auth := authentication.NewInsecureStaticProvider(jwt)
	party := cfg.Party
	if party == "" {
		conn, err := grpc.NewClient(
			cfg.GRPCLedgerAPIURL,
			grpc.WithTransportCredentials(auth.TransportCredentials()),
			grpc.WithPerRPCCredentials(auth.PerRPCCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("ledger grpc dial: %w", err)
		}
		defer conn.Close()
		user, err := adminv2.NewUserManagementServiceClient(conn).GetUser(ctx, &adminv2.GetUserRequest{UserId: cfg.UserName})
		if err != nil {
			return nil, fmt.Errorf("get user %q: %w", cfg.UserName, err)
		}
		if len(user.GetUser().GetPrimaryParty()) == 0 {
			return nil, fmt.Errorf("no primary party for user %q", cfg.UserName)
		}
		party = user.GetUser().GetPrimaryParty()
	}

	provider := cantonProvider.NewRPCChainProvider(selector, cantonProvider.RPCChainProviderConfig{
		Participants: []cantonProvider.ParticipantConfig{{
			Endpoints: cantonProvider.Endpoints{
				JSONLedgerAPIURL: cfg.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: cfg.GRPCLedgerAPIURL,
				AdminAPIURL:      cfg.AdminAPIURL,
				ValidatorAPIURL:  cfg.ValidatorAPIURL,
			},
			UserID:       cfg.UserName,
			PartyID:      party,
			AuthProvider: auth,
		}},
	})
	bc, err := provider.Initialize(ctx)
	if err != nil {
		return nil, err
	}
	chain, ok := bc.(*canton.Chain)
	if !ok {
		return nil, fmt.Errorf("unexpected chain type %T", bc)
	}
	return chain, nil
}

func nameOrDefault(cfg testhelpers.ParticipantConfig, idx int) string {
	if cfg.Name != "" {
		return cfg.Name
	}
	return fmt.Sprintf("participant%d", idx+1)
}
