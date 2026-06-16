package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/rs/zerolog"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/scripts/prod_testnet/internal/prodtestnetenv"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

const defaultCantonGRPCURL = ""
const defaultValidatorAPIURL = ""
const defaultUserID = ""
const defaultPartyID = ""
const defaultAuthType = "clientCredentials"
const defaultAuthURL = ""
const defaultClientID = ""
const defaultClientSecret = ""
const defaultSrcSelector = uint64(9268731218649498074)
const defaultProdTestnetAmuletDSOAdmin = "DSO::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337"

type holdingRow struct {
	admin   string
	id      string
	amount  string
	locked  bool
	contractID string
}

func main() {
	if _, err := prodtestnetenv.LoadDefault(); err != nil {
		fatalf("load scripts/prod_testnet/.env: %v", err)
	}

	srcSelectorDefault, err := prodtestnetenv.Uint64(defaultSrcSelector, "PROD_TESTNET_CANTON_TO_EVM_SRC_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}

	var (
		grpcURL         = flag.String("grpc-url", prodtestnetenv.String(defaultCantonGRPCURL, "PROD_TESTNET_CANTON_GRPC_URL"), "Canton participant gRPC ledger API URL")
		validatorAPIURL = flag.String("validator-api-url", prodtestnetenv.String(defaultValidatorAPIURL, "PROD_TESTNET_CANTON_VALIDATOR_API_URL"), "Canton validator API base URL")
		userID          = flag.String("user-id", prodtestnetenv.String(defaultUserID, "PROD_TESTNET_CANTON_USER_ID"), "Canton user ID")
		partyID         = flag.String("party-id", prodtestnetenv.String(defaultPartyID, "PROD_TESTNET_CANTON_PARTY_ID"), "Canton party ID")
		authType        = flag.String("auth-type", prodtestnetenv.String(defaultAuthType, "PROD_TESTNET_CANTON_AUTH_TYPE"), "Canton auth type")
		authURL         = flag.String("auth-url", prodtestnetenv.String(defaultAuthURL, "PROD_TESTNET_CANTON_AUTH_URL"), "OIDC issuer URL")
		clientID        = flag.String("client-id", prodtestnetenv.String(defaultClientID, "PROD_TESTNET_CANTON_CLIENT_ID"), "OAuth2 client ID")
		clientSecret    = flag.String("client-secret", prodtestnetenv.String(defaultClientSecret, "PROD_TESTNET_CANTON_CLIENT_SECRET", "CLIENT_SECRET"), "OAuth2 client secret")
		jwtToken        = flag.String("jwt", prodtestnetenv.String("", "PROD_TESTNET_CANTON_JWT"), "JWT token for static auth")
		srcSelector     = flag.Uint64("src", srcSelectorDefault, "Canton chain selector")
		feeAdmin        = flag.String("fee-admin", prodtestnetenv.String(defaultProdTestnetAmuletDSOAdmin, "PROD_TESTNET_CANTON_FEE_TOKEN_ADMIN"), "Expected fee token admin for CCIP send")
		feeID           = flag.String("fee-id", prodtestnetenv.String("Amulet", "PROD_TESTNET_CANTON_FEE_TOKEN_ID"), "Expected fee token id")
	)
	flag.Parse()

	requireFlag("grpc-url", "PROD_TESTNET_CANTON_GRPC_URL", *grpcURL)
	requireFlag("validator-api-url", "PROD_TESTNET_CANTON_VALIDATOR_API_URL", *validatorAPIURL)
	requireFlag("user-id", "PROD_TESTNET_CANTON_USER_ID", *userID)
	requireFlag("party-id", "PROD_TESTNET_CANTON_PARTY_ID", *partyID)
	if *srcSelector == 0 {
		fatalf("missing -src (set PROD_TESTNET_CANTON_TO_EVM_SRC_SELECTOR)")
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	authTypeTrim := strings.TrimSpace(*authType)
	clientSecretVal := strings.TrimSpace(*clientSecret)
	if authTypeTrim != commonconfig.AuthTypeClientCredentials {
		clientSecretVal = ""
	}
	authCfg := commonconfig.AuthConfig{
		Type:         authTypeTrim,
		UserID:       *userID,
		AuthURL:      *authURL,
		ClientID:     *clientID,
		ClientSecret: clientSecretVal,
		JWT:          *jwtToken,
	}
	authProvider, err := authCfg.NewProvider(ctx)
	if err != nil {
		fatalf("build canton auth provider: %v", err)
	}

	chainProvider := provider.NewRPCChainProvider(*srcSelector, provider.RPCChainProviderConfig{
		Participants: []provider.ParticipantConfig{{
			Endpoints: provider.Endpoints{
				GRPCLedgerAPIURL: *grpcURL,
				ValidatorAPIURL:  *validatorAPIURL,
			},
			UserID:       *userID,
			PartyID:      *partyID,
			AuthProvider: authProvider,
		}},
	})
	blockChain, err := chainProvider.Initialize(ctx)
	if err != nil {
		fatalf("initialize canton chain provider: %v", err)
	}
	chain, ok := blockChain.(*cldfcanton.Chain)
	if !ok {
		fatalf("unexpected chain provider type %T", blockChain)
	}
	participant := chain.Participants[0]

	rows, err := listPartyTokenHoldings(ctx, participant, types.PARTY(*partyID))
	if err != nil {
		fatalf("list holdings: %v", err)
	}

	logger.Info().Str("party", *partyID).Int("count", len(rows)).Msg("token holdings visible to party")
	if len(rows) == 0 {
		fmt.Println("No token holdings found.")
		fmt.Println()
		fmt.Println("Canton→Sepolia send needs an unlocked Amulet holding (DSO admin) for CCIP fees.")
		fmt.Println("Prod testnet has no DevNet tap — transfer Amulet to this party from a funded wallet/validator.")
		os.Exit(0)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].admin != rows[j].admin {
			return rows[i].admin < rows[j].admin
		}
		if rows[i].id != rows[j].id {
			return rows[i].id < rows[j].id
		}
		return rows[i].contractID < rows[j].contractID
	})

	var usableFee int
	for _, row := range rows {
		usable := !row.locked && row.amount != "" && row.amount != "0" && row.amount != "0.0"
		matchesFee := strings.TrimSpace(*feeAdmin) != "" &&
			row.admin == strings.TrimSpace(*feeAdmin) &&
			row.id == strings.TrimSpace(*feeID) &&
			usable
		if matchesFee || (strings.TrimSpace(*feeAdmin) == "" && row.id == strings.TrimSpace(*feeID) && usable) {
			usableFee++
		}
		fmt.Printf("- admin=%s id=%s amount=%s locked=%t cid=%s\n", row.admin, row.id, row.amount, row.locked, row.contractID)
	}

	fmt.Println()
	if usableFee == 0 {
		fmt.Printf("No unlocked %q holdings usable for CCIP fees", strings.TrimSpace(*feeID))
		if strings.TrimSpace(*feeAdmin) != "" {
			fmt.Printf(" under admin %s", strings.TrimSpace(*feeAdmin))
		}
		fmt.Println(".")
		fmt.Println("Fund this party with a small Amulet transfer, then retry send_prod_testnet_canton_to_evm.")
	} else {
		fmt.Printf("Found %d unlocked %q holding(s) that may work for fees.\n", usableFee, strings.TrimSpace(*feeID))
	}
}

func listPartyTokenHoldings(ctx context.Context, participant cldfcanton.Participant, owner types.PARTY) ([]holdingRow, error) {
	offset, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end: %w", err)
	}

	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: offset.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{{
						IdentifierFilter: &ledgerv2.CumulativeFilter_InterfaceFilter{InterfaceFilter: &ledgerv2.InterfaceFilter{
							InterfaceId: &ledgerv2.Identifier{
								PackageId:  "#splice-api-token-holding-v1",
								ModuleName: "Splice.Api.Token.HoldingV1",
								EntityName: "Holding",
							},
							IncludeInterfaceView:    true,
							IncludeCreatedEventBlob: false,
						}},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active holdings: %w", err)
	}
	defer stream.CloseSend()

	holdingInterfaceID := &ledgerv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	}

	var rows []holdingRow
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("receive active holdings: %w", err)
		}
		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		created := entry.ActiveContract.GetCreatedEvent()
		viewRecord, err := getRelevantInterfaceViewValue(created.GetInterfaceViews(), holdingInterfaceID)
		if err != nil {
			continue
		}
		var holdingView splice_api_token_holding_v1.HoldingView
		if err := ledger.RecordToStruct(viewRecord, &holdingView); err != nil {
			return nil, fmt.Errorf("decode holding view for %s: %w", created.GetContractId(), err)
		}
		if holdingView.Owner != owner {
			continue
		}
		rows = append(rows, holdingRow{
			admin:      string(holdingView.InstrumentId.Admin),
			id:         string(holdingView.InstrumentId.Id),
			amount:     string(holdingView.Amount),
			locked:     holdingView.Lock != nil,
			contractID: created.GetContractId(),
		})
	}
	return rows, nil
}

func getRelevantInterfaceViewValue(interfaceViews []*ledgerv2.InterfaceView, expectedInterfaceID *ledgerv2.Identifier) (*ledgerv2.Record, error) {
	for _, interfaceView := range interfaceViews {
		if interfaceView.GetInterfaceId().GetModuleName() == expectedInterfaceID.GetModuleName() &&
			interfaceView.GetInterfaceId().GetEntityName() == expectedInterfaceID.GetEntityName() {
			return interfaceView.GetViewValue(), nil
		}
	}
	return nil, fmt.Errorf("no interface view found for %s:%s", expectedInterfaceID.GetModuleName(), expectedInterfaceID.GetEntityName())
}

func requireFlag(flagName, envKey, value string) {
	if strings.TrimSpace(value) == "" {
		fatalf("missing -%s (set it explicitly or via %s)", flagName, envKey)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
