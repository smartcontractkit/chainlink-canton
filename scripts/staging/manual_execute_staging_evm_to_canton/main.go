package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipreceiver"
	ccipcommon "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiEDSCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/scripts/staging/internal/stagingenv"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
	v1 "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/api/handlers/v1"
	indexerclient "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/client"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

// This script calls the CCV indexer, Canton EDS (HTTPS under *.ccip.stage.internal.griddle.sh), then
// the Canton participant ledger. Those internal hostnames usually require corporate VPN; if a single
// VPN profile cannot reach both indexer and Canton, run through indexer/auth/EDS first, then use
// -vpn-switch-wait before any ledger Submit (the program sleeps and logs when to switch).

const (
	defaultIndexerURL             = ""
	defaultEDSURL                 = "https://chainlink-ccv-canton-eds-canton-1.ccip.stage.internal.griddle.sh"
	defaultCantonGRPCURL          = ""
	defaultValidatorAPIURL        = ""
	defaultUserID                 = ""
	defaultPartyID                = ""
	defaultAuthType               = commonconfig.AuthTypeAuthorizationCode
	defaultAuthURL                = ""
	defaultClientID               = ""
	defaultDestSelector    uint64 = 0
	defaultTimeout                = 2 * time.Minute
	defaultPollInterval           = 2 * time.Second
	defaultExpectedResults        = 1
	defaultVPNSwitchWait          = 15 * time.Second
	defaultLedgerTimeout          = 2 * time.Minute
	// Mainnet-family staging EVM selector used when configuring GlobalConfig source chains for Canton staging lanes.
	defaultSourceChainSelector uint64 = 16015286601757825753
)

func main() {
	if _, err := stagingenv.LoadDefault(); err != nil {
		fatalf("load scripts/staging/.env: %v", err)
	}

	destSelectorDefault, err := stagingenv.Uint64(defaultDestSelector, "STAGING_EVM_TO_CANTON_DST_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}
	timeoutDefault, err := stagingenv.Duration(defaultTimeout, "STAGING_EVM_TO_CANTON_WAIT_TIMEOUT")
	if err != nil {
		fatalf("%v", err)
	}
	ledgerTimeoutDefault, err := stagingenv.Duration(defaultLedgerTimeout, "STAGING_EVM_TO_CANTON_LEDGER_TIMEOUT")
	if err != nil {
		fatalf("%v", err)
	}
	pollIntervalDefault, err := stagingenv.Duration(defaultPollInterval, "STAGING_EVM_TO_CANTON_POLL_INTERVAL")
	if err != nil {
		fatalf("%v", err)
	}
	vpnSwitchWaitDefault, err := stagingenv.Duration(defaultVPNSwitchWait, "STAGING_EVM_TO_CANTON_VPN_SWITCH_WAIT")
	if err != nil {
		fatalf("%v", err)
	}
	sourceChainSelectorDefault, err := stagingenv.Uint64(defaultSourceChainSelector, "STAGING_EVM_TO_CANTON_SOURCE_CHAIN_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}

	var (
		indexerURL              = flag.String("indexer-url", stagingenv.String(defaultIndexerURL, "STAGING_EVM_TO_CANTON_INDEXER_URL"), "Indexer base URL")
		edsURL                  = flag.String("eds-url", stagingenv.String(defaultEDSURL, "STAGING_CANTON_EDS_URL"), "EDS base URL")
		grpcURL                 = flag.String("grpc-url", stagingenv.String(defaultCantonGRPCURL, "STAGING_CANTON_GRPC_URL"), "Canton participant gRPC ledger API URL")
		validatorAPIURL         = flag.String("validator-api-url", stagingenv.String(defaultValidatorAPIURL, "STAGING_CANTON_VALIDATOR_API_URL"), "Canton validator API base URL")
		messageIDHex            = flag.String("message-id", "", "0x-prefixed message ID to execute")
		queryUpdateID           = flag.String("query-update-id", "", "Query an already committed Canton update by update ID and print created events")
		destSelector            = flag.Uint64("dest", destSelectorDefault, "Destination Canton chain selector")
		timeout                 = flag.Duration("timeout", timeoutDefault, "How long to wait for indexed verifications")
		ledgerTimeout           = flag.Duration("ledger-timeout", ledgerTimeoutDefault, "How long to allow post-handoff Canton ledger operations")
		pollInterval            = flag.Duration("poll-interval", pollIntervalDefault, "How often to poll the indexer")
		expectedVerifierResults = flag.Int("expected-verifier-results", defaultExpectedResults, "Expected number of verifier results before execution")
		vpnSwitchWait           = flag.Duration("vpn-switch-wait", vpnSwitchWaitDefault, "Pause after auth/indexer fetch before the first Canton ledger transaction")
		printJSONOnly           = flag.Bool("print-json-only", false, "Print the fetched indexer result JSON and exit")
		receiverMinBlockConfs   = flag.Int64("receiver-min-block-confirmations", -1, "CCIPReceiver receiverFinalityConfig: block depth (BlockDepth); 0 = WaitForFinality; default uses message finality")
		userID                  = flag.String("user-id", stagingenv.String(defaultUserID, "STAGING_CANTON_USER_ID"), "Canton user ID")
		partyID                 = flag.String("party-id", stagingenv.String(defaultPartyID, "STAGING_CANTON_PARTY_ID"), "Canton party ID used to execute")
		authType                = flag.String("auth-type", stagingenv.String(defaultAuthType, "STAGING_CANTON_AUTH_TYPE"), "Canton auth type: authorizationCode, clientCredentials, static, insecureStatic")
		authURL                 = flag.String("auth-url", stagingenv.String(defaultAuthURL, "STAGING_CANTON_AUTH_URL"), "OIDC auth URL for authorizationCode")
		clientID                = flag.String("client-id", stagingenv.String(defaultClientID, "STAGING_CANTON_CLIENT_ID"), "OIDC client ID for authorizationCode")
		clientSecret            = flag.String("client-secret", stagingenv.String("", "STAGING_CANTON_CLIENT_SECRET"), "OIDC client secret for clientCredentials")
		jwtToken                = flag.String("jwt", stagingenv.String("", "STAGING_CANTON_JWT"), "JWT token for static/insecureStatic auth")
		ccvOverride             = flag.String("ccv", stagingenv.String("", "STAGING_EVM_TO_CANTON_CCV", "STAGING_EVM_TO_CANTON_COMMITTEE_VERIFIER"), "Canton CommitteeVerifier instance address override for execute disclosures")
		routerCIDOverride       = flag.String("router-cid", "", "Existing PerPartyRouter contract ID override")
		receiverCIDOverride     = flag.String("receiver-cid", "", "Existing CCIPReceiver contract ID override")
		routerInstanceID        = flag.String("router-instance-id", "", "Router instance ID override")
		receiverInstanceID      = flag.String("receiver-instance-id", "", "Receiver instance ID override")
		sourceChainSelector     = flag.Uint64("source-chain-selector", sourceChainSelectorDefault, "Source chain selector for GlobalConfig GetSourceChainConfig (EVM lane)")
		globalConfigCID         = flag.String("global-config-cid", stagingenv.String("", "STAGING_EVM_TO_CANTON_GLOBAL_CONFIG_CID"), "GlobalConfig contract ID for optional preflight cross-check only — must match execute choice context global-config (defaults to that context)")
		skipGlobalConfigCheck   = flag.Bool("skip-global-config-check", false, "Do not exercise GlobalConfig GetSourceChainConfig before Execute")
	)
	flag.Parse()

	if *messageIDHex == "" && *queryUpdateID == "" {
		fatalf("missing -message-id")
	}
	requireString(*indexerURL, "-indexer-url", "STAGING_EVM_TO_CANTON_INDEXER_URL")
	requireString(*edsURL, "-eds-url", "STAGING_CANTON_EDS_URL")
	requireString(*grpcURL, "-grpc-url", "STAGING_CANTON_GRPC_URL")
	requireString(*validatorAPIURL, "-validator-api-url", "STAGING_CANTON_VALIDATOR_API_URL")
	requireString(*userID, "-user-id", "STAGING_CANTON_USER_ID")
	requireString(*partyID, "-party-id", "STAGING_CANTON_PARTY_ID")
	requireUint64(*destSelector, "-dest", "STAGING_EVM_TO_CANTON_DST_SELECTOR")

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	indexer, err := indexerclient.NewIndexerClient(*indexerURL, &http.Client{Timeout: 15 * time.Second})
	if err != nil {
		fatalf("create indexer client: %v", err)
	}

	authCfg := commonconfig.AuthConfig{
		Type:         *authType,
		UserID:       *userID,
		AuthURL:      *authURL,
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		JWT:          *jwtToken,
	}
	authProvider, err := authCfg.NewProvider(ctx)
	if err != nil {
		fatalf("build canton auth provider: %v", err)
	}

	chainProvider := provider.NewRPCChainProvider(*destSelector, provider.RPCChainProviderConfig{
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

	if *queryUpdateID != "" {
		if err := queryCommittedUpdate(ctx, participant, *queryUpdateID, &logger); err != nil {
			fatalf("query committed update: %v", err)
		}
		return
	}

	messageIDBytes := hexutil.MustDecode(*messageIDHex)
	var messageID protocol.Bytes32
	copy(messageID[:], messageIDBytes)

	resp, err := waitForVerifierResults(ctx, indexer, messageID, *pollInterval, *expectedVerifierResults, &logger)
	if err != nil {
		fatalf("wait for verifier results: %v", err)
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		fatalf("marshal indexer response: %v", err)
	}
	if *printJSONOnly {
		fmt.Println(string(respJSON))
		return
	}

	ccipEDS, err := oapiCCIP.NewClientWithResponses(strings.TrimSuffix(strings.TrimSpace(*edsURL), "/"), oapiCCIP.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}))
	if err != nil {
		fatalf("create ccip eds client: %v", err)
	}
	validatorAuth := func(ctx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve validator token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(validatorAuth),
	)
	if err != nil {
		fatalf("create token metadata client: %v", err)
	}
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(validatorAuth),
	)
	if err != nil {
		fatalf("create transfer instruction client: %v", err)
	}

	message := resp.Results[0].VerifierResult.Message
	encodedMessage, err := message.Encode()
	if err != nil {
		fatalf("encode message: %v", err)
	}
	decSanity, err := protocol.DecodeMessage(encodedMessage)
	if err != nil {
		fatalf("decode encoded message (sanity): %v", err)
	}
	if !bytes.Equal(decSanity.OnRampAddress, message.OnRampAddress) {
		fatalf("re-decoded message on-ramp %x does not match indexer message on-ramp %x", decSanity.OnRampAddress, message.OnRampAddress)
	}

	execInputs, err := getExecuteInputs(ctx, strings.TrimSuffix(strings.TrimSpace(*edsURL), "/"), message, resp, encodedMessage, *ccvOverride, &logger)
	if err != nil {
		fatalf("get execute disclosures: %v", err)
	}

	if *vpnSwitchWait > 0 {
		logger.Info().
			Dur("wait", *vpnSwitchWait).
			Str("hint", "Canton ledger gRPC and/or EDS may need a different VPN path than the indexer").
			Msg("indexer/auth/EDS disclosures done — switch VPN now before Canton ledger transactions")
		time.Sleep(*vpnSwitchWait)
	}

	ledgerCtx, ledgerCancel := context.WithTimeout(context.Background(), *ledgerTimeout)
	defer ledgerCancel()

	routerID := *routerInstanceID
	if routerID == "" {
		routerID = "router-" + shortMessageID(messageIDHex)
	}
	receiverID := *receiverInstanceID
	if receiverID == "" {
		receiverID = "receiver-" + shortMessageID(messageIDHex)
	}

	routerCID := *routerCIDOverride
	if routerCID == "" {
		var err error
		routerCID, err = ensurePerPartyRouter(ledgerCtx, participant, ccipEDS, *partyID, routerID)
		if err != nil {
			fatalf("create per-party router: %v", err)
		}
	}
	logger.Info().Str("routerCID", routerCID).Msg("resolved per-party router")

	minBlockConfs := *receiverMinBlockConfs
	if minBlockConfs < 0 {
		minBlockConfs = int64(message.Finality)
	}

	receiverCID := *receiverCIDOverride
	if receiverCID == "" {
		var err error
		receiverCID, err = findExistingReceiverCID(ledgerCtx, participant, *partyID, receiverID)
		if err != nil {
			fatalf("find existing receiver: %v", err)
		}
		if receiverCID == "" {
			receiverCID, err = deployReceiver(ledgerCtx, participant, *partyID, receiverID, minBlockConfs)
			if err != nil {
				fatalf("deploy receiver: %v", err)
			}
			logger.Info().
				Str("receiverCID", receiverCID).
				Int64("receiverFinalityBlockDepth", minBlockConfs).
				Msg("deployed ccip receiver")
		} else {
			logger.Info().
				Str("receiverCID", receiverCID).
				Int64("receiverFinalityBlockDepth", minBlockConfs).
				Msg("resolved existing ccip receiver")
		}
	} else {
		logger.Info().
			Str("receiverCID", receiverCID).
			Int64("receiverFinalityBlockDepth", minBlockConfs).
			Msg("using existing ccip receiver")
	}

	if !*skipGlobalConfigCheck {
		gcFromCtx, ctxErr := globalConfigContractIDFromChoiceContext(execInputs.ChoiceContext)
		gcCID := gcFromCtx
		if gcCID == "" {
			gcCID = globalConfigContractIDFromDisclosures(execInputs.DisclosedContracts)
			if gcCID != "" && ctxErr != nil {
				logger.Warn().Err(ctxErr).Msg("GlobalConfig: could not read global-config from execute choice context; using first GlobalConfig in disclosures (may not match PrepareExecute — fix EDS/context if preflight succeeds but Execute fails)")
			}
		} else if ctxErr != nil {
			fatalf("GlobalConfig check: %v", ctxErr)
		}
		gcOverride := strings.TrimSpace(*globalConfigCID)
		if gcCID == "" {
			fatalf("GlobalConfig check: missing contract id (CCIP execute choice context should include global-config as AV_ContractId, or disclosures must list GlobalConfig)")
		}
		if gcOverride != "" && !strings.EqualFold(gcOverride, gcCID) {
			fatalf("GlobalConfig mismatch: -global-config-cid / env is %q but execute path uses %q (from choice context when available, else disclosures); PrepareExecute always uses the context — unset STAGING_EVM_TO_CANTON_GLOBAL_CONFIG_CID or align the override", gcOverride, gcCID)
		}
		if gcFromCtx != "" && gcOverride == "" {
			logger.Info().Str("globalConfigCID", gcCID).Msg("GlobalConfig contract id from execute choice context (matches PrepareExecute)")
		}
		expectedOnRamp := normalizeHexLower(hex.EncodeToString(message.OnRampAddress))
		if err := verifyGlobalConfigSourceOnRamp(ledgerCtx, participant, *partyID, gcCID, *sourceChainSelector, message.OnRampAddress, expectedOnRamp, execInputs.DisclosedContracts, &logger); err != nil {
			fatalf("GlobalConfig GetSourceChainConfig check: %v", err)
		}
	}

	disclosedContracts := execInputs.DisclosedContracts
	var receiverHoldingsBefore []*apiv2.ActiveContract
	if execInputs.RequiresTokenRelease() {
		receiverHoldingsBefore, err = listActiveContractsByInterfaceID(ledgerCtx, participant, &apiv2.Identifier{
			PackageId:  "#splice-api-token-holding-v1",
			ModuleName: "Splice.Api.Token.HoldingV1",
			EntityName: "Holding",
		})
		if err != nil {
			fatalf("list receiver holdings before execute: %v", err)
		}
		logger.Info().
			Float64("receiverHoldingsBefore", holdingsBalance(receiverHoldingsBefore)).
			Msg("canton receiver total holdings before execute")

		registryAdmin, err := getRegistryAdmin(ledgerCtx, tokenMetadataClient)
		if err != nil {
			fatalf("get registry admin: %v", err)
		}
		transferFactoryCID, transferFactoryDisclosures, transferFactoryContext, err := getTransferFactory(ledgerCtx, transferInstructionClient, registryAdmin, *partyID, *partyID)
		if err != nil {
			fatalf("get transfer factory: %v", err)
		}
		disclosedContracts = append(disclosedContracts, transferFactoryDisclosures...)
		execInputs.TransferFactoryCID = transferFactoryCID
		execInputs.TransferFactoryContext = transferFactoryContext
	}

	ccvInputs := make([]*apiv2.Value, len(execInputs.CCVContractIDs))
	for i := range execInputs.CCVContractIDs {
		ccvExtra := emptyCCIPContext()
		if i < len(execInputs.CCVExtraContexts) && execInputs.CCVExtraContexts[i] != nil {
			ccvExtra = execInputs.CCVExtraContexts[i]
		}
		ccvInputs[i] = &apiv2.Value{
			Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "ccvCid", Value: &apiv2.Value{Sum: execInputs.CCVContractIDs[i]}},
				{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: execInputs.VerifierResultsHex[i]}}},
				{Label: "ccvExtraContext", Value: ccvExtra},
			}}},
		}
	}

	execRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ledgerCtx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: receiverCID,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "context", Value: execInputs.ChoiceContext},
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCID}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: hex.EncodeToString(encodedMessage)}}},
						{Label: "tokenTransfer", Value: execInputs.TokenTransferValue(*partyID)},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: ccvInputs}}}},
					}}}},
				}},
			}},
			ActAs:              []string{*partyID},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		fatalf("submit execute transaction: %v", err)
	}

	receivedMessageID := ""
	var receivedEventFields map[string]string
	for _, event := range execRes.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok && e.Created.GetTemplateId().GetEntityName() == "CCIPMessageReceived" {
			fields := e.Created.GetCreateArguments().GetFields()
			if len(fields) > 2 {
				receivedMessageID = fields[2].GetValue().GetText()
			}
			receivedEventFields = flattenRecordFields(e.Created.GetCreateArguments())
			break
		}
	}

	if receivedEventFields != nil {
		logger.Info().
			Str("updateID", execRes.GetTransaction().GetUpdateId()).
			Str("contractID", receiverCID).
			Str("eventTemplate", "CCIPMessageReceived").
			Fields(zerolog.Dict().Fields(receivedEventFields)).
			Msg("canton execution event")
	}

	logger.Info().
		Str("messageID", *messageIDHex).
		Str("updateID", execRes.GetTransaction().GetUpdateId()).
		Str("receivedMessageID", receivedMessageID).
		Msg("manual Canton execution submitted")

	if execInputs.RequiresTokenRelease() {
		receiverHoldingsAfter, err := listActiveContractsByInterfaceID(ledgerCtx, participant, &apiv2.Identifier{
			PackageId:  "#splice-api-token-holding-v1",
			ModuleName: "Splice.Api.Token.HoldingV1",
			EntityName: "Holding",
		})
		if err != nil {
			fatalf("list receiver holdings after execute: %v", err)
		}
		beforeBalance := holdingsBalance(receiverHoldingsBefore)
		afterBalance := holdingsBalance(receiverHoldingsAfter)
		logger.Info().
			Float64("receiverHoldingsBefore", beforeBalance).
			Float64("receiverHoldingsAfter", afterBalance).
			Float64("receiverHoldingsDelta", afterBalance-beforeBalance).
			Msg("canton receiver holdings validated after token execute")
	}
}

// verifyGlobalConfigSourceOnRamp exercises GlobalConfig GetSourceChainConfig (non-consuming) and
// checks that onRampAddresses contains the same BytesHex text DAML uses for (message.onRampAddress `elem` ...).
// DAML equality is exact Text match (case-sensitive hex; no trimming or alternate encodings).
func verifyGlobalConfigSourceOnRamp(
	ctx context.Context,
	participant cldfcanton.Participant,
	partyID string,
	globalConfigCID string,
	sourceChainSelector uint64,
	onRampMsgBytes []byte,
	expectedOnRampHexLower string,
	disclosed []*apiv2.DisclosedContract,
	logger *zerolog.Logger,
) error {
	tpl := contracts.TemplateIDFromBinding(ccipcommon.GlobalConfig{}).ToLedgerIdentifier()
	choiceArg := ledger.MapToValue(ccipcommon.GetSourceChainConfig{
		SourceChainSelector: types.NUMERIC(strconv.FormatUint(sourceChainSelector, 10)),
		Caller:              types.PARTY(partyID),
	}.ToMap())

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId:     tpl,
					ContractId:     globalConfigCID,
					Choice:         "GetSourceChainConfig",
					ChoiceArgument: choiceArg,
				}},
			}},
			ActAs:              []string{partyID},
			DisclosedContracts: disclosed,
		},
		TransactionFormat: &apiv2.TransactionFormat{
			EventFormat: &apiv2.EventFormat{
				FiltersByParty: map[string]*apiv2.Filters{
					partyID: {},
				},
			},
			TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_LEDGER_EFFECTS,
		},
	})
	if err != nil {
		return fmt.Errorf("submit GetSourceChainConfig: %w", err)
	}

	var exercised *apiv2.ExercisedEvent
	for _, event := range res.GetTransaction().GetEvents() {
		if ex := event.GetExercised(); ex != nil && ex.GetChoice() == "GetSourceChainConfig" {
			exercised = ex
			break
		}
	}
	if exercised == nil {
		return errors.New("ledger transaction had no GetSourceChainConfig exercised event")
	}

	cfgRec, ok := unwrapOptionalSourceChainConfigRecord(exercised.GetExerciseResult())
	if !ok {
		return fmt.Errorf("no SourceChainConfig for sourceChainSelector=%d (absent in GlobalConfig)", sourceChainSelector)
	}

	onRampAddrs, err := sourceChainConfigOnRampHexList(cfgRec)
	if err != nil {
		return err
	}
	// Exact BytesHex string the ledger + Daml decode use for MessageV1.onRampAddress when the wire
	// encoding is produced by Go's hex.EncodeToString(Encode()) (lowercase a-f).
	msgOnRampHex := hex.EncodeToString(onRampMsgBytes)
	// Do not TrimSpace ledger entries: Daml `elem` compares the stored Text exactly.

	var strictMatch bool
	var casingOnly string
	for _, entry := range onRampAddrs {
		entryNorm := normalizeHexLower(entry)
		if entryNorm == "" {
			continue
		}
		if entry == msgOnRampHex {
			strictMatch = true
			break
		}
		if entryNorm == msgOnRampHex {
			casingOnly = entry
		}
	}

	// Optional: same address bytes with different width (20 vs 32 byte EVM encoding) — elem still fails; surface clearly.
	if !strictMatch && casingOnly == "" {
		mb := onRampMsgBytes
		for _, entry := range onRampAddrs {
			eb, err := hex.DecodeString(normalizeHexLower(entry))
			if err != nil || len(eb) == 0 {
				continue
			}
			// Wrong pipeline encoded onRamp as hex(ASCII(hex_digits)): ledger Text is ~128 hex chars
			// whose byte payload spells the same string as msgOnRampHex. Daml `elem` needs the
			// single-encoded BytesHex Text (same 64-char string as message.onRampAddress), not this.
			if bytesLookLikeASCIIHexDigits(eb) && strings.EqualFold(string(eb), msgOnRampHex) {
				return fmt.Errorf("GlobalConfig onRampAddresses entry stores hex-of-ASCII-hex (double encoding): ledger BytesHex text %q decodes to bytes that read %q — same logical on-ramp as message %q, but Daml compares Text equality; fix deployment / ApplySourceChainConfigUpdates so the list contains literal BytesHex %q", entry, string(eb), msgOnRampHex, msgOnRampHex)
			}
			if bytes.Equal(mb, eb) {
				return fmt.Errorf("GlobalConfig on-ramp %q decodes to the same %d bytes as the message on-ramp, but DAML requires exact BytesHex text match (elem): ledger entry %q != message %q — store onRampAddresses using the same hex string as in CCIP MessageV1 (from wire)", entry, len(mb), entry, msgOnRampHex)
			}
			const evmAddrLen = 20
			if len(mb) == evmAddrLen && len(eb) == 32 && bytes.Equal(bytes.Repeat([]byte{0}, 12), eb[:12]) && bytes.Equal(mb, eb[12:]) {
				return fmt.Errorf("message on-ramp is %d bytes (%s) but GlobalConfig lists 32-byte padded %q; DAML elem will fail — add the on-ramp as the same %d-byte or 32-byte hex encoding used in outbound messages", evmAddrLen, msgOnRampHex, entry, evmAddrLen)
			}
			if len(mb) == 32 && len(eb) == evmAddrLen && bytes.Equal(bytes.Repeat([]byte{0}, 12), mb[:12]) && bytes.Equal(mb[12:], eb) {
				return fmt.Errorf("message on-ramp is 32-byte padded (%s) but GlobalConfig lists 20-byte %q; DAML elem will fail — align onRampAddresses encoding with MessageV1", msgOnRampHex, entry)
			}
		}
	}

	if logger != nil {
		logger.Info().
			Str("globalConfigCID", globalConfigCID).
			Uint64("sourceChainSelector", sourceChainSelector).
			Str("messageOnRampHexForElem", msgOnRampHex).
			Str("expectedOnRampHexLower", expectedOnRampHexLower).
			Strs("onRampAddressesFromLedgerRaw", onRampAddrs).
			Bool("damlElemMatch", strictMatch).
			Msg("GlobalConfig GetSourceChainConfig preflight (strict BytesHex elem)")
	}

	if strictMatch {
		return nil
	}
	if casingOnly != "" {
		return fmt.Errorf("on-ramp differs only by hex letter case: GlobalConfig has %q but message/codec uses %q; DAML (message.onRampAddress `elem` onRampAddresses) is case-sensitive — re-apply onRampAddresses with lowercase a-f", casingOnly, msgOnRampHex)
	}
	return fmt.Errorf("onRampAddresses must contain exact BytesHex %q for offramp validate (ledger has %v); if a tool shows a match using normalized ASCII-hex, it is not what Daml elem checks", msgOnRampHex, onRampAddrs)
}

// globalConfigContractIDFromChoiceContext returns the GlobalConfig contract id that PerPartyRouter.PrepareExecute
// uses (CCIP.GlobalConfig.lookupGlobalConfigCid). Preflight must use this same cid — not an arbitrary disclosure
// or env override pointing at a different GlobalConfig instance.
func globalConfigContractIDFromChoiceContext(ctx *apiv2.Value) (string, error) {
	const key = "global-config"
	if ctx == nil {
		return "", errors.New("execute choice context is nil")
	}
	rec := ctx.GetRecord()
	if rec == nil {
		return "", errors.New("execute choice context is not a record")
	}
	var tm *apiv2.TextMap
	for _, f := range rec.GetFields() {
		if f.GetLabel() == "values" && f.GetValue().GetTextMap() != nil {
			tm = f.GetValue().GetTextMap()
			break
		}
	}
	if tm == nil && len(rec.GetFields()) > 0 && rec.GetFields()[0].GetValue().GetTextMap() != nil {
		tm = rec.GetFields()[0].GetValue().GetTextMap()
	}
	if tm == nil {
		return "", errors.New("execute choice context has no values map")
	}
	for _, e := range tm.GetEntries() {
		if e.GetKey() != key {
			continue
		}
		if cid, ok := choiceContextVariantContractID(e.GetValue()); ok {
			return cid, nil
		}
		return "", fmt.Errorf("execute choice context %q is not AV_ContractId", key)
	}
	return "", fmt.Errorf("execute choice context missing %q (PerPartyRouter cannot resolve GlobalConfig)", key)
}

func choiceContextVariantContractID(v *apiv2.Value) (cid string, ok bool) {
	if v == nil {
		return "", false
	}
	if vt := v.GetVariant(); vt != nil {
		c := vt.GetConstructor()
		if strings.HasSuffix(c, "AV_ContractId") || strings.Contains(c, ":AV_ContractId") {
			if inner := vt.GetValue(); inner != nil {
				if cid := inner.GetContractId(); cid != "" {
					return cid, true
				}
			}
		}
	}
	if cid := v.GetContractId(); cid != "" {
		return cid, true
	}
	return "", false
}

func globalConfigContractIDFromDisclosures(disclosed []*apiv2.DisclosedContract) string {
	for _, d := range disclosed {
		tid := d.GetTemplateId()
		if tid != nil && tid.GetEntityName() == "GlobalConfig" && tid.GetModuleName() == "CCIP.GlobalConfig" {
			return d.GetContractId()
		}
	}
	return ""
}

func normalizeHexLower(s string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "0x"))
}

func bytesLookLikeASCIIHexDigits(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	for _, c := range b {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func unwrapOptionalSourceChainConfigRecord(v *apiv2.Value) (*apiv2.Record, bool) {
	if v == nil {
		return nil, false
	}
	if opt := v.GetOptional(); opt != nil {
		if opt.GetValue() == nil {
			return nil, false
		}
		if r := unwrapRecordMaybeUnlabeledProduct(opt.GetValue()); r != nil {
			return r, true
		}
		return nil, false
	}
	if varnt := v.GetVariant(); varnt != nil {
		if isDamlOptionalSomeConstructor(varnt.GetConstructor()) && varnt.GetValue() != nil {
			if r := unwrapRecordMaybeUnlabeledProduct(varnt.GetValue()); r != nil {
				return r, true
			}
		}
		return nil, false
	}
	if r := unwrapRecordMaybeUnlabeledProduct(v); r != nil {
		return r, true
	}
	return nil, false
}

// isDamlOptionalSomeConstructor accepts "Some", package-qualified constructors, etc.
func isDamlOptionalSomeConstructor(constructor string) bool {
	c := strings.TrimSpace(constructor)
	if c == "" {
		return false
	}
	if c == "Some" {
		return true
	}
	if strings.HasSuffix(c, ":Some") {
		return true
	}
	if strings.HasSuffix(c, ".Some") {
		return true
	}
	return false
}

// unwrapRecordMaybeUnlabeledProduct unwraps a single-field record wrapping the payload (e.g. Optional Some, data constructors).
func unwrapRecordMaybeUnlabeledProduct(v *apiv2.Value) *apiv2.Record {
	if v == nil {
		return nil
	}
	r := v.GetRecord()
	if r == nil {
		return nil
	}
	fields := r.GetFields()
	if len(fields) == 1 && fields[0].GetValue() != nil {
		if inner := fields[0].GetValue().GetRecord(); inner != nil && len(inner.GetFields()) > 1 {
			return inner
		}
	}
	return r
}

// sourceChainConfigOnRampHexList reads onRampAddresses from CCIP.GlobalConfigTypes:SourceChainConfig.
// Ledger API records may omit field labels; fall back to declaration order: isEnabled, onRampAddresses, ...
func sourceChainConfigOnRampHexList(rec *apiv2.Record) ([]string, error) {
	if rec == nil {
		return nil, fmt.Errorf("record is nil")
	}
	fields := rec.GetFields()
	for _, field := range fields {
		if strings.EqualFold(field.GetLabel(), "onRampAddresses") {
			return hexTextListFromValue(field.GetValue())
		}
	}
	// Declaration order: isEnabled, onRampAddresses, defaultCCVs, laneMandatedCCVs (labels may be empty on ledger).
	if len(fields) >= 2 && fields[1].GetValue().GetList() != nil {
		return hexTextListFromValue(fields[1].GetValue())
	}
	labels := make([]string, 0, len(fields))
	for _, f := range fields {
		labels = append(labels, f.GetLabel())
	}
	return nil, fmt.Errorf("could not locate onRampAddresses (field labels: %v)", labels)
}

func hexTextListFromValue(v *apiv2.Value) ([]string, error) {
	if v == nil {
		return nil, fmt.Errorf("value is nil")
	}
	list := v.GetList()
	if list == nil {
		return nil, fmt.Errorf("value is not a list")
	}
	out := make([]string, 0, len(list.GetElements()))
	for _, el := range list.GetElements() {
		if el == nil {
			return nil, fmt.Errorf("onRamp list has nil element")
		}
		if _, ok := el.GetSum().(*apiv2.Value_Text); !ok {
			return nil, fmt.Errorf("onRamp list element is not Text (BytesHex): %s", valueToString(el))
		}
		out = append(out, el.GetText())
	}
	return out, nil
}

func ensurePerPartyRouter(
	ctx context.Context,
	participant cldfcanton.Participant,
	edsClient *oapiCCIP.ClientWithResponses,
	partyID string,
	instanceID string,
) (string, error) {
	existingRouterCID, err := findExistingRouterCID(ctx, participant, partyID)
	if err != nil {
		return "", fmt.Errorf("find existing router: %w", err)
	}
	if existingRouterCID != "" {
		return existingRouterCID, nil
	}

	factoryCID, disclosedContracts, err := getPerPartyRouterFactoryDisclosures(ctx, edsClient, partyID)
	if err != nil {
		return "", err
	}

	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: factoryCID,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{
						PartyOwner: types.PARTY(partyID),
						InstanceId: types.TEXT(instanceID),
					}),
				}},
			}},
			ActAs:              []string{partyID},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		return "", err
	}
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok && e.Created.GetTemplateId().GetEntityName() == "PerPartyRouter" {
			return e.Created.ContractId, nil
		}
	}

	return "", errors.New("per-party router create transaction produced no PerPartyRouter contract")
}

func queryCommittedUpdate(ctx context.Context, participant cldfcanton.Participant, updateID string, logger *zerolog.Logger) error {
	updateRes, err := participant.LedgerServices.Update.GetUpdateById(ctx, &apiv2.GetUpdateByIdRequest{
		UpdateId: updateID,
		UpdateFormat: &apiv2.UpdateFormat{
			IncludeTransactions: &apiv2.TransactionFormat{
				TransactionShape: apiv2.TransactionShape_TRANSACTION_SHAPE_ACS_DELTA,
				EventFormat: &apiv2.EventFormat{
					FiltersByParty: map[string]*apiv2.Filters{
						participant.PartyID: {
							Cumulative: []*apiv2.CumulativeFilter{{
								IdentifierFilter: &apiv2.CumulativeFilter_WildcardFilter{
									WildcardFilter: &apiv2.WildcardFilter{IncludeCreatedEventBlob: false},
								},
							}},
						},
					},
					Verbose: true,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	found := false
	for _, event := range updateRes.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil {
			fields := flattenRecordFields(created.GetCreateArguments())
			ev := logger.Info().
				Str("updateID", updateID).
				Str("contractID", created.GetContractId()).
				Str("eventTemplate", created.GetTemplateId().GetEntityName())
			if len(fields) > 0 {
				ev = ev.Fields(fields)
			} else if raw, err := json.Marshal(created); err == nil {
				ev = ev.Str("rawCreatedEvent", string(raw))
			}
			if created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
				for _, field := range created.GetCreateArguments().GetFields() {
					if field.GetLabel() != "event" || field.GetValue().GetRecord() == nil {
						continue
					}
					eventFields := flattenRecordFields(field.GetValue().GetRecord())
					if len(eventFields) == 0 {
						break
					}
					if raw, err := json.Marshal(eventFields); err == nil {
						ev = ev.Str("ccipMessageSentEvent", string(raw))
					}
					if receipts, ok := eventFields["receipts"]; ok {
						ev = ev.Str("receipts", receipts)
					}
					break
				}
			}
			ev.Msg("canton committed event")
			found = true
		}
	}
	if !found {
		logger.Info().Str("updateID", updateID).Msg("no created events found in committed update")
	}
	return nil
}

func listActiveContractsByInterfaceID(
	ctx context.Context,
	participant cldfcanton.Participant,
	interfaceID *apiv2.Identifier,
) ([]*apiv2.ActiveContract, error) {
	offsetResp, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end: %w", err)
	}

	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offsetResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{{
						IdentifierFilter: &apiv2.CumulativeFilter_InterfaceFilter{
							InterfaceFilter: &apiv2.InterfaceFilter{
								InterfaceId:             interfaceID,
								IncludeInterfaceView:    true,
								IncludeCreatedEventBlob: false,
							},
						},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active contracts by interface: %w", err)
	}
	defer stream.CloseSend()

	var activeContracts []*apiv2.ActiveContract
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("receive active contracts by interface: %w", err)
		}
		active, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || active.ActiveContract == nil || active.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		activeContracts = append(activeContracts, active.ActiveContract)
	}
	return activeContracts, nil
}

func holdingsBalance(holdings []*apiv2.ActiveContract) float64 {
	var total float64
	for _, h := range holdings {
		views := h.GetCreatedEvent().GetInterfaceViews()
		if len(views) == 0 || views[0].GetViewValue() == nil {
			continue
		}
		fields := views[0].GetViewValue().GetFields()
		if len(fields) < 3 {
			continue
		}
		balanceStr := fields[2].GetValue().GetNumeric()
		balance, ok := new(big.Float).SetString(balanceStr)
		if !ok {
			continue
		}
		v, _ := balance.Float64()
		total += v
	}
	return total
}

func flattenRecordFields(record *apiv2.Record) map[string]string {
	if record == nil {
		return nil
	}
	out := make(map[string]string, len(record.GetFields()))
	for i, field := range record.GetFields() {
		key := field.GetLabel()
		if key == "" {
			key = fmt.Sprintf("field_%d", i)
		}
		out[key] = valueToString(field.GetValue())
	}
	return out
}

func valueToString(v *apiv2.Value) string {
	if v == nil || v.GetSum() == nil {
		return ""
	}
	switch x := v.GetSum().(type) {
	case *apiv2.Value_Text:
		return x.Text
	case *apiv2.Value_Int64:
		return fmt.Sprintf("%d", x.Int64)
	case *apiv2.Value_Bool:
		return fmt.Sprintf("%t", x.Bool)
	case *apiv2.Value_ContractId:
		return x.ContractId
	case *apiv2.Value_Party:
		return x.Party
	case *apiv2.Value_Unit:
		return "unit"
	case *apiv2.Value_Optional:
		if x.Optional == nil || x.Optional.Value == nil {
			return "null"
		}
		return valueToString(x.Optional.Value)
	case *apiv2.Value_Record:
		b, _ := json.Marshal(flattenRecordFields(x.Record))
		return string(b)
	case *apiv2.Value_List:
		parts := make([]string, 0, len(x.List.GetElements()))
		for _, el := range x.List.GetElements() {
			parts = append(parts, valueToString(el))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case *apiv2.Value_Numeric:
		return x.Numeric
	case *apiv2.Value_Timestamp:
		return fmt.Sprintf("%d", x.Timestamp)
	case *apiv2.Value_Date:
		return fmt.Sprintf("%d", x.Date)
	case *apiv2.Value_Enum:
		return x.Enum.GetConstructor()
	case *apiv2.Value_Variant:
		return x.Variant.GetConstructor()
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func findExistingRouterCID(ctx context.Context, participant cldfcanton.Participant, partyID string) (string, error) {
	offsetResp, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", err
	}

	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offsetResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
								TemplateId: &apiv2.Identifier{
									PackageId:  "#ccip-perpartyrouter",
									ModuleName: "CCIP.PerPartyRouter",
									EntityName: "PerPartyRouter",
								},
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return "", err
	}
	defer stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		active, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || active.ActiveContract == nil || active.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		fields := active.ActiveContract.GetCreatedEvent().GetCreateArguments().GetFields()
		for _, field := range fields {
			if field.GetLabel() == "partyOwner" && field.GetValue().GetParty() == partyID {
				return active.ActiveContract.GetCreatedEvent().GetContractId(), nil
			}
		}
	}
}

func findExistingReceiverCID(ctx context.Context, participant cldfcanton.Participant, partyID, instanceID string) (string, error) {
	offsetResp, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &apiv2.GetLedgerEndRequest{})
	if err != nil {
		return "", err
	}

	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &apiv2.GetActiveContractsRequest{
		ActiveAtOffset: offsetResp.GetOffset(),
		EventFormat: &apiv2.EventFormat{
			FiltersByParty: map[string]*apiv2.Filters{
				participant.PartyID: {
					Cumulative: []*apiv2.CumulativeFilter{
						{
							IdentifierFilter: &apiv2.CumulativeFilter_TemplateFilter{TemplateFilter: &apiv2.TemplateFilter{
								TemplateId: &apiv2.Identifier{
									PackageId:  "#ccip-receiver",
									ModuleName: "CCIP.CCIPReceiver",
									EntityName: "CCIPReceiver",
								},
								IncludeCreatedEventBlob: true,
							}},
						},
					},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return "", err
	}
	defer stream.CloseSend()

	var ownerMatch string
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return ownerMatch, nil
		}
		if err != nil {
			return "", err
		}
		active, ok := resp.GetContractEntry().(*apiv2.GetActiveContractsResponse_ActiveContract)
		if !ok || active.ActiveContract == nil || active.ActiveContract.GetCreatedEvent() == nil {
			continue
		}

		created := active.ActiveContract.GetCreatedEvent()
		matchedOwner := false
		matchedInstance := instanceID == ""
		for _, field := range created.GetCreateArguments().GetFields() {
			switch field.GetLabel() {
			case "owner":
				matchedOwner = field.GetValue().GetParty() == partyID
			case "instanceId":
				matchedInstance = instanceID == "" || field.GetValue().GetText() == instanceID
			}
		}

		if matchedOwner && matchedInstance {
			return created.GetContractId(), nil
		}
		if matchedOwner && ownerMatch == "" {
			ownerMatch = created.GetContractId()
		}
	}
}

// receiverFinalityConfigValue encodes CCIP.FinalityConfig for CCIPReceiver create
// (same shape as integration-tests/ccip/ccip_execute_test.go finalityConfigValueFromBlockConfirmations).
func receiverFinalityConfigValue(blockConfirmations int64) *apiv2.Value {
	if blockConfirmations == 0 {
		return &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
			Constructor: "WaitForFinality",
			Value:       &apiv2.Value{Sum: &apiv2.Value_Unit{}},
		}}}
	}
	return &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
		Constructor: "BlockDepth",
		Value:       &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: blockConfirmations}},
	}}}
}

func deployReceiver(ctx context.Context, participant cldfcanton.Participant, partyID string, instanceID string, minBlockConfirmations int64) (string, error) {
	res, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: instanceID}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyID}}},
						{Label: "receiverFinalityConfig", Value: receiverFinalityConfigValue(minBlockConfirmations)},
						{Label: "requiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalThreshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
					}},
				}},
			}},
			ActAs: []string{partyID},
		},
	})
	if err != nil {
		return "", err
	}

	contractID := extractFirstCreatedContractID(res)
	if contractID == "" {
		return "", errors.New("receiver create transaction produced no created contract")
	}
	return contractID, nil
}

func getExecuteInputs(
	ctx context.Context,
	edsBaseURL string,
	message protocol.Message,
	resp v1.VerifierResultsByMessageIDResponse,
	encodedMessage []byte,
	ccvOverride string,
	logger *zerolog.Logger,
) (*executeInputs, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	ccipCl, err := oapiCCIP.NewClientWithResponses(edsBaseURL, oapiCCIP.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("ccip eds client: %w", err)
	}
	ccvCl, err := oapiCCV.NewClientWithResponses(edsBaseURL, oapiCCV.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("ccv eds client: %w", err)
	}
	tpCl, err := oapiTokenPool.NewClientWithResponses(edsBaseURL, oapiTokenPool.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("token pool eds client: %w", err)
	}

	ccvs := make([]contracts.InstanceAddress, len(resp.Results))
	verifierResultsHex := make([]string, len(resp.Results))
	for i, result := range resp.Results {
		if ccvOverride != "" {
			ccvs[i] = contracts.HexToInstanceAddress(ccvOverride)
			if logger != nil {
				logger.Info().
					Str("overrideCCV", ccvOverride).
					Str("verifierDestAddress", result.VerifierResult.VerifierDestAddress.String()).
					Msg("using explicit Canton CCV override for execute disclosures")
			}
		} else {
			ccvs[i] = contracts.HexToInstanceAddress(result.VerifierResult.VerifierDestAddress.String())
		}
		verifierResultsHex[i] = hex.EncodeToString(result.VerifierResult.CCVData)
	}

	encHex := hex.EncodeToString(encodedMessage)
	ccipDisc, err := eds.GetCCIPExecuteDisclosure(ctx, ccipCl, encHex)
	if err != nil {
		return nil, fmt.Errorf("ccip execute disclosure: %w", err)
	}
	choiceContext, err := testhelpers.ChoiceContextFromStruct(ccipDisc.ChoiceContext)
	if err != nil {
		return nil, fmt.Errorf("ccip execute choice context: %w", err)
	}

	disclosedContracts := append([]*apiv2.DisclosedContract(nil), ccipDisc.DisclosedContracts...)

	ccvContractIDs := make([]*apiv2.Value_ContractId, len(ccvs))
	ccvExtraContexts := make([]*apiv2.Value, len(ccvs))
	for i, ccv := range ccvs {
		ccvDisc, err := eds.GetCCVExecuteDisclosure(ctx, ccvCl, encHex, ccv)
		if err != nil {
			return nil, fmt.Errorf("ccv execute disclosure for %s: %w", ccv.String(), err)
		}
		ccvContractIDs[i] = &apiv2.Value_ContractId{ContractId: ccvDisc.ContractId}
		ccvCtx, err := testhelpers.ChoiceContextFromStruct(ccvDisc.ChoiceContext)
		if err != nil {
			return nil, fmt.Errorf("ccv execute choice context for %s: %w", ccv.String(), err)
		}
		ccvExtraContexts[i] = ccvCtx
		disclosedContracts = append(disclosedContracts, ccvDisc.DisclosedContracts...)
	}

	var poolExtraContext *apiv2.Value
	var tokenPoolContractID *apiv2.Value_ContractId
	var tokenPoolHoldingContractIDs []*apiv2.Value

	if message.TokenTransfer != nil && ccipDisc.TokenPool != nil {
		tpInst := (*ccipDisc.TokenPool).InstanceAddress()
		tpd, err := eds.GetTokenPoolExecuteDisclosure(ctx, tpCl, encHex, tpInst)
		if err != nil {
			return nil, fmt.Errorf("token pool execute disclosure: %w", err)
		}
		tokenPoolContractID = &apiv2.Value_ContractId{ContractId: tpd.ContractId}
		poolExtraContext, err = testhelpers.ChoiceContextFromStruct(tpd.ChoiceContext)
		if err != nil {
			return nil, fmt.Errorf("token pool execute choice context: %w", err)
		}
		disclosedContracts = append(disclosedContracts, tpd.DisclosedContracts...)
	}

	return &executeInputs{
		DisclosedContracts:   disclosedContracts,
		ChoiceContext:        choiceContext,
		PoolExtraContext:     poolExtraContext,
		CCVContractIDs:       ccvContractIDs,
		CCVExtraContexts:     ccvExtraContexts,
		VerifierResultsHex:   verifierResultsHex,
		TokenPoolContractID:  tokenPoolContractID,
		TokenPoolHoldingCIDs: tokenPoolHoldingContractIDs,
	}, nil
}

type executeInputs struct {
	DisclosedContracts     []*apiv2.DisclosedContract
	ChoiceContext          *apiv2.Value
	PoolExtraContext       *apiv2.Value
	CCVContractIDs         []*apiv2.Value_ContractId
	CCVExtraContexts       []*apiv2.Value
	VerifierResultsHex     []string
	TokenPoolContractID    *apiv2.Value_ContractId
	TokenPoolHoldingCIDs   []*apiv2.Value
	TransferFactoryCID     string
	TransferFactoryContext *apiv2.Value
}

func (e *executeInputs) RequiresTokenRelease() bool {
	return e != nil && e.TokenPoolContractID != nil
}

func (e *executeInputs) TokenTransferValue(tokenReceiverParty string) *apiv2.Value {
	if e == nil || e.TokenPoolContractID == nil {
		return &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: nil}}}
	}

	poolExtraContext := e.PoolExtraContext
	if poolExtraContext == nil {
		poolExtraContext = emptyCCIPContext()
	}
	extraArgsContext := e.TransferFactoryContext
	if extraArgsContext == nil {
		extraArgsContext = emptyCCIPContext()
	}
	tokenPoolHoldings := e.TokenPoolHoldingCIDs
	if tokenPoolHoldings == nil {
		tokenPoolHoldings = []*apiv2.Value{}
	}

	return &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: e.TokenPoolContractID}},
		{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: tokenReceiverParty}}},
		{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: e.TransferFactoryCID}}},
			{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
				{Label: "context", Value: extraArgsContext},
				{Label: "meta", Value: emptyMetadata()},
			}}}}},
			{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: tokenPoolHoldings}}}},
		}}}}},
		{Label: "poolExtraContext", Value: poolExtraContext},
	}}}}}}}
}

func getPerPartyRouterFactoryDisclosures(
	ctx context.Context,
	edsClient *oapiCCIP.ClientWithResponses,
	partyID string,
) (string, []*apiv2.DisclosedContract, error) {
	resp, err := edsClient.PostPerPartyRouterFactoryWithResponse(ctx, oapiCCIP.CCIPPerPartyRouterFactoryRequest{PartyID: oapiEDSCommon.PartyId(partyID)})
	if err != nil {
		return "", nil, fmt.Errorf("failed to get per-party router factory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return "", nil, fmt.Errorf("unexpected per-party-router-factory status: %d", resp.StatusCode())
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.DisclosedContracts {
		disclosedContract, err := disclosedContractToProto(contract)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert disclosed factory contract: %w", err)
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	return resp.JSON200.ContractId, disclosedContracts, nil
}

func waitForVerifierResults(
	ctx context.Context,
	indexer *indexerclient.IndexerClient,
	messageID protocol.Bytes32,
	pollInterval time.Duration,
	expectedVerifierResults int,
	logger *zerolog.Logger,
) (v1.VerifierResultsByMessageIDResponse, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	msgIDHex := messageID.String()
	for {
		status, resp, err := indexer.VerifierResultsByMessageID(ctx, v1.VerifierResultsByMessageIDInput{MessageID: msgIDHex})
		if err == nil && status == http.StatusOK && resp.Success && len(resp.Results) >= expectedVerifierResults {
			logger.Info().
				Str("messageID", msgIDHex).
				Int("verifierResults", len(resp.Results)).
				Msg("found verifier results in indexer")
			return resp, nil
		}
		if err == nil {
			logger.Info().
				Str("messageID", msgIDHex).
				Int("status", status).
				Bool("success", resp.Success).
				Int("verifierResults", len(resp.Results)).
				Msg("waiting for verifier results in indexer")
		} else {
			logger.Info().Err(err).Str("messageID", msgIDHex).Msg("waiting for verifier results in indexer")
		}

		select {
		case <-ctx.Done():
			return v1.VerifierResultsByMessageIDResponse{}, fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func disclosedContractToProto(contract oapiEDSCommon.DisclosedContract) (*apiv2.DisclosedContract, error) {
	id, err := templateIDFromString(contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       id,
		ContractId:       contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   contract.SynchronizerId,
	}, nil
}

func templateIDFromString(s string) (*apiv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}

	return &apiv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func choiceContextFromData(choiceContextData map[string]any) (*apiv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	fields := make([]*apiv2.TextMap_Entry, 0, len(values))
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *apiv2.Value
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Text{Text: valueString}}
		case "AV_Int":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: valueBool}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: valueString}}
		default:
			return nil, fmt.Errorf("unimplemented tag: %v", tag)
		}

		fields = append(fields, &apiv2.TextMap_Entry{
			Key: k,
			Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
				Constructor: tag,
				Value:       value,
			}}},
		})
	}

	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: fields}}},
		},
	}}}}, nil
}

func emptyCCIPContext() *apiv2.Value {
	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
		},
	}}}}
}

func emptyMetadata() *apiv2.Value {
	return &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{
			Label: "values",
			Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
		},
	}}}}
}

func getRegistryAdmin(ctx context.Context, metadataClient tokenMetadataV1.ClientWithResponsesInterface) (string, error) {
	resp, err := metadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting registry info: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %v", resp.StatusCode(), resp.Body)
	}
	return resp.JSON200.AdminId, nil
}

func getTransferFactory(
	ctx context.Context,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	registryAdmin string,
	sender string,
	receiver string,
) (string, []*apiv2.DisclosedContract, *apiv2.Value, error) {
	resp, err := transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   sender,
				"receiver": receiver,
				"amount":   "100.00",
				"instrumentId": map[string]any{
					"admin": registryAdmin,
					"id":    "Amulet",
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(-time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"inputHoldingCids": []string{},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{
					"values": map[string]any{},
				},
				"meta": map[string]any{
					"values": map[string]any{},
				},
			},
		},
	})
	if err != nil {
		return "", nil, nil, fmt.Errorf("error getting transfer factory response: %w", err)
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return "", nil, nil, fmt.Errorf("unexpected status code: %d: %v", resp.StatusCode(), resp.Body)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		disclosedContract, err := disclosedContractToProto(oapiEDSCommon.DisclosedContract{
			TemplateId:       contract.TemplateId,
			ContractId:       contract.ContractId,
			CreatedEventBlob: contract.CreatedEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
		if err != nil {
			return "", nil, nil, fmt.Errorf("failed to convert transfer factory disclosed contract: %w", err)
		}
		disclosedContracts = append(disclosedContracts, disclosedContract)
	}

	choiceContext, err := choiceContextFromData(resp.JSON200.ChoiceContext.ChoiceContextData)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to convert transfer factory choice context: %w", err)
	}

	return resp.JSON200.FactoryId, disclosedContracts, choiceContext, nil
}

func extractFirstCreatedContractID(res *apiv2.SubmitAndWaitForTransactionResponse) string {
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			return e.Created.ContractId
		}
	}
	return ""
}

func shortMessageID(messageIDHex *string) string {
	msg := strings.TrimPrefix(*messageIDHex, "0x")
	if len(msg) <= 12 {
		return msg
	}
	return msg[:12]
}

func requireString(value string, flagName string, envHint string) {
	if strings.TrimSpace(value) == "" {
		if envHint == "" {
			fatalf("missing required value: pass %s", flagName)
		}
		fatalf("missing required value: pass %s or set %s in scripts/staging/.env", flagName, envHint)
	}
}

func requireUint64(value uint64, flagName string, envHint string) {
	if value != 0 {
		return
	}
	if envHint == "" {
		fatalf("missing required value: pass %s", flagName)
	}
	fatalf("missing required value: pass %s or set %s in scripts/staging/.env", flagName, envHint)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

var _ = ccipreceiver.CCIPReceiver{}
var _ = ccipcommon.GlobalConfig{}
var _ io.Reader
