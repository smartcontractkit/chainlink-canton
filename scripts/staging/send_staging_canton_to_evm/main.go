package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccipsender"
	ccipclient "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/client"
	executorbinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiEDSCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/scripts/staging/internal/stagingeds"
	"github.com/smartcontractkit/chainlink-canton/scripts/staging/internal/stagingenv"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	cldfcanton "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"
)

const (
	defaultCantonGRPCURL    = ""
	defaultValidatorAPIURL  = ""
	defaultEDSURL           = "https://chainlink-ccv-canton-eds-canton-1.ccip.stage.internal.griddle.sh"
	defaultUserID           = ""
	defaultPartyID          = ""
	defaultAuthType         = commonconfig.AuthTypeAuthorizationCode
	defaultAuthURL          = ""
	defaultClientID         = ""
	defaultClientSecret     = ""
	defaultSrcSelector      = uint64(0)
	defaultDstSelector      = uint64(0)
	defaultReceiver         = ""
	defaultData             = "hello from canton to evm"
	defaultExecutionGas     = int64(200_000)
	defaultFinalityConfig   = int64(0)
	defaultWaitTimeout      = 2 * time.Minute
	defaultVPNSwitchWait    = 15 * time.Second
	defaultScanRetryTimeout = 45 * time.Second
	defaultSenderInstanceID = ""

	defaultCommitteeVerifier = ""
	defaultExecutor          = ""
)

// Staging ledgers sometimes register the CCIP runtime DAR under a short package alias; bindings use "ccip-runtime".
const perPartyRouterPackageAlias = "#ccip-perpartyrouter"

// CCIPSender DAR is often registered as #ccip-sender (see integration-tests/ccip/ccip_send_test.go).
const ccipSenderPackageAlias = "#ccip-sender"

func ccipsenderTemplateIDs() []string {
	primary := ccipsender.CCIPSender{}.GetTemplateID()
	alt := ccipSenderPackageAlias + ":CCIP.CCIPSender:CCIPSender"
	if primary == alt {
		return []string{primary}
	}
	return []string{primary, alt}
}

func perPartyRouterFactoryTemplateIDs() []string {
	primary := perpartyrouter.PerPartyRouterFactory{}.GetTemplateID()
	alt := perPartyRouterPackageAlias + ":CCIP.PerPartyRouter:PerPartyRouterFactory"
	if primary == alt {
		return []string{primary}
	}
	return []string{primary, alt}
}

func perPartyRouterTemplateIDs() []string {
	primary := perpartyrouter.PerPartyRouter{}.GetTemplateID()
	alt := perPartyRouterPackageAlias + ":CCIP.PerPartyRouter:PerPartyRouter"
	if primary == alt {
		return []string{primary}
	}
	return []string{primary, alt}
}

// edsBaseRequiresSplitVPNPhase is true for hosted staging EDS on *.griddle.sh, which is on SmartContracts VPN
// while Canton gRPC, validator API, and scan-proxy require Chainlink Legacy VPN.
func edsBaseRequiresSplitVPNPhase(edsBase string) bool {
	return strings.Contains(strings.ToLower(edsBase), ".griddle.")
}

func waitForStagingVPN(logger *zerolog.Logger, wait time.Duration, networkLabel, instruction string) {
	if wait <= 0 {
		return
	}
	logger.Info().Dur("wait", wait).Str("expectedVPN", networkLabel).Msg(instruction)
	time.Sleep(wait)
}

// perPartyRouterFactoryPreload is filled during the SmartContracts VPN phase (hosted EDS only).
type perPartyRouterFactoryPreload struct {
	ContractID string
	Disclosed  []*ledgerv2.DisclosedContract
}

func fetchPerPartyRouterFactoryPreloadFromHostedEDS(ctx context.Context, edsBaseURL, partyID string, logger *zerolog.Logger) (*perPartyRouterFactoryPreload, error) {
	edsBase := strings.TrimSpace(edsBaseURL)
	if edsBase == "" {
		return nil, fmt.Errorf("empty EDS base URL")
	}
	ccipEDS, err := oapiCCIP.NewClientWithResponses(edsBase, oapiCCIP.WithHTTPClient(&http.Client{Timeout: 20 * time.Second}))
	if err != nil {
		return nil, fmt.Errorf("create CCIP EDS client: %w", err)
	}
	fd, derr := eds.GetPerPartyRouterFactoryDisclosure(ctx, ccipEDS, partyID)
	if derr != nil {
		return nil, derr
	}
	if fd == nil || fd.ContractId == "" || len(fd.DisclosedContracts) == 0 {
		return nil, fmt.Errorf("empty PerPartyRouterFactory disclosure from EDS")
	}
	if logger != nil {
		logger.Info().Str("factoryContractId", fd.ContractId).Msg("fetched PerPartyRouterFactory from hosted EDS (SmartContracts VPN phase)")
	}
	return &perPartyRouterFactoryPreload{ContractID: fd.ContractId, Disclosed: fd.DisclosedContracts}, nil
}

func main() {
	if _, err := stagingenv.LoadDefault(); err != nil {
		fatalf("load scripts/staging/.env: %v", err)
	}

	srcSelectorDefault, err := stagingenv.Uint64(defaultSrcSelector, "STAGING_CANTON_TO_EVM_SRC_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}
	dstSelectorDefault, err := stagingenv.Uint64(defaultDstSelector, "STAGING_CANTON_TO_EVM_DST_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}
	waitTimeoutDefault, err := stagingenv.Duration(defaultWaitTimeout, "STAGING_CANTON_TO_EVM_WAIT_TIMEOUT")
	if err != nil {
		fatalf("%v", err)
	}
	vpnSwitchWaitDefault, err := stagingenv.Duration(defaultVPNSwitchWait, "STAGING_CANTON_TO_EVM_VPN_SWITCH_WAIT")
	if err != nil {
		fatalf("%v", err)
	}

	var (
		grpcURL           = flag.String("grpc-url", stagingenv.String(defaultCantonGRPCURL, "STAGING_CANTON_GRPC_URL"), "Canton participant gRPC ledger API URL")
		validatorAPIURL   = flag.String("validator-api-url", stagingenv.String(defaultValidatorAPIURL, "STAGING_CANTON_VALIDATOR_API_URL"), "Canton validator API base URL")
		edsURL            = flag.String("eds-url", stagingenv.String(defaultEDSURL, "STAGING_CANTON_EDS_URL"), "EDS base URL")
		userID            = flag.String("user-id", stagingenv.String(defaultUserID, "STAGING_CANTON_USER_ID"), "Canton user ID")
		partyID           = flag.String("party-id", stagingenv.String(defaultPartyID, "STAGING_CANTON_PARTY_ID"), "Canton party ID used to send")
		authType          = flag.String("auth-type", stagingenv.String(defaultAuthType, "STAGING_CANTON_AUTH_TYPE"), "Canton auth type: authorizationCode, clientCredentials, static, insecureStatic")
		authURL           = flag.String("auth-url", stagingenv.String(defaultAuthURL, "STAGING_CANTON_AUTH_URL"), "OIDC issuer URL (authorizationCode / clientCredentials)")
		clientID          = flag.String("client-id", stagingenv.String(defaultClientID, "STAGING_CANTON_CLIENT_ID"), "OAuth2 client ID (authorizationCode / clientCredentials)")
		clientSecret      = flag.String("client-secret", stagingenv.String(defaultClientSecret, "STAGING_CANTON_CLIENT_SECRET", "CLIENT_SECRET"), "OAuth2 client secret (clientCredentials only; ignored for other auth types)")
		jwtToken          = flag.String("jwt", stagingenv.String("", "STAGING_CANTON_JWT"), "JWT token for static/insecureStatic auth")
		srcSelector       = flag.Uint64("src", srcSelectorDefault, "Source Canton chain selector")
		dstSelector       = flag.Uint64("dest", dstSelectorDefault, "Destination EVM chain selector")
		receiverHex       = flag.String("receiver", stagingenv.String(defaultReceiver, "STAGING_CANTON_TO_EVM_RECEIVER"), "Destination EVM receiver address")
		messageData       = flag.String("data", defaultData, "Message payload")
		executionGasLimit = flag.Int64("execution-gas-limit", defaultExecutionGas, "Execution gas limit")
		finalityConfig    = flag.Int64("finality-config", defaultFinalityConfig, "Finality config / block confirmations")
		waitTimeout       = flag.Duration("wait-timeout", waitTimeoutDefault, "How long to wait for the send transaction")
		vpnSwitchWait     = flag.Duration("vpn-switch-wait", vpnSwitchWaitDefault, "Pause before each VPN phase when EDS URL is *.griddle.* (Legacy→SmartContracts→Legacy). Ignored sub-waits can be skipped with 0 only if both networks are reachable")
		scanRetryTimeout  = flag.Duration("scan-retry-timeout", defaultScanRetryTimeout, "How long to retry scan-proxy backed calls")

		senderInstanceID  = flag.String("sender-instance-id", stagingenv.String(defaultSenderInstanceID, "STAGING_CANTON_TO_EVM_SENDER_INSTANCE_ID"), "CCIPSender instance ID")
		routerInstanceIDFlag = flag.String("router-instance-id", stagingenv.String("", "STAGING_CANTON_TO_EVM_ROUTER_INSTANCE_ID"), "PerPartyRouter instance id (defaults to sender-instance-id when empty)")
		perPartyRouterFactoryCIDFlag = flag.String("per-party-router-factory-cid", stagingenv.String("", "STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_CID"), "Optional PerPartyRouterFactory contract ID if your party cannot see the factory via ACS")
		perPartyRouterFactoryAddrFlag = flag.String("per-party-router-factory-address", stagingenv.String("", "STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_ADDRESS"), "Optional PerPartyRouterFactory instance address (0x… hex from address_refs) if ACS template scan misses")
		committeeVerifier = flag.String("committee-verifier", stagingenv.String(defaultCommitteeVerifier, "STAGING_CANTON_TO_EVM_COMMITTEE_VERIFIER"), "CommitteeVerifier CCV for EDS: hex InstanceAddress (0x…) or raw instanceId@party (same as ccip/devenv opts.CCVs; no ledger ACS lookup)")
		executorAddr      = flag.String("executor", stagingenv.String(defaultExecutor, "STAGING_CANTON_TO_EVM_EXECUTOR"), "Source Canton Executor instance address")
	)
	flag.Parse()

	requireFlag := func(flagName, envKey, value string) {
		if strings.TrimSpace(value) == "" {
			fatalf("missing -%s (set it explicitly or via %s)", flagName, envKey)
		}
	}
	requireUint64Flag := func(flagName, envKey string, value uint64) {
		if value == 0 {
			fatalf("missing -%s (set it explicitly or via %s)", flagName, envKey)
		}
	}

	requireFlag("grpc-url", "STAGING_CANTON_GRPC_URL", *grpcURL)
	requireFlag("validator-api-url", "STAGING_CANTON_VALIDATOR_API_URL", *validatorAPIURL)
	requireFlag("eds-url", "STAGING_CANTON_EDS_URL", *edsURL)
	requireFlag("user-id", "STAGING_CANTON_USER_ID", *userID)
	requireFlag("party-id", "STAGING_CANTON_PARTY_ID", *partyID)
	requireFlag("receiver", "STAGING_CANTON_TO_EVM_RECEIVER", *receiverHex)
	requireFlag("sender-instance-id", "STAGING_CANTON_TO_EVM_SENDER_INSTANCE_ID", *senderInstanceID)
	routerInstanceID := strings.TrimSpace(*routerInstanceIDFlag)
	if routerInstanceID == "" {
		routerInstanceID = strings.TrimSpace(*senderInstanceID)
	}
	requireFlag("committee-verifier", "STAGING_CANTON_TO_EVM_COMMITTEE_VERIFIER", *committeeVerifier)
	requireFlag("executor", "STAGING_CANTON_TO_EVM_EXECUTOR", *executorAddr)
	requireUint64Flag("src", "STAGING_CANTON_TO_EVM_SRC_SELECTOR", *srcSelector)
	requireUint64Flag("dest", "STAGING_CANTON_TO_EVM_DST_SELECTOR", *dstSelector)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithTimeout(context.Background(), *waitTimeout)
	defer cancel()

	authTypeTrim := strings.TrimSpace(*authType)
	clientSecretVal := strings.TrimSpace(*clientSecret)
	if authTypeTrim != commonconfig.AuthTypeClientCredentials {
		// Validator requires client_secret only for clientCredentials; it must be unset for authorizationCode.
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

	edsBaseTrim := strings.TrimSuffix(strings.TrimSpace(*edsURL), "/")
	splitGriddleEDS := edsBaseRequiresSplitVPNPhase(edsBaseTrim)
	if splitGriddleEDS {
		logger.Info().Msg(`split VPN mode: hosted EDS on *.griddle.sh uses SmartContracts.com VPN; Canton gRPC, validator API, and scan-proxy use Chainlink Legacy VPN. Script order: Legacy (init + prep) → SmartContracts (EDS) → Legacy (router + submit).`)
		waitForStagingVPN(&logger, *vpnSwitchWait, "chainlink_legacy", "Connect Chainlink Legacy VPN before Canton participant gRPC Initialize. Waiting…")
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
	validatorAuth := func(ctx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve validator token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}
	scanProxyClient, err := scanProxy.NewClientWithResponses(participant.Endpoints.ValidatorAPIURL, scanProxy.WithRequestEditorFn(validatorAuth))
	if err != nil {
		fatalf("create scan-proxy client: %v", err)
	}
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL), tokenMetadataV1.WithRequestEditorFn(validatorAuth))
	if err != nil {
		fatalf("create token metadata client: %v", err)
	}
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL), transferInstructionV1.WithRequestEditorFn(validatorAuth))
	if err != nil {
		fatalf("create transfer instruction client: %v", err)
	}

	resolveDisclosedByAddress := func(templateID string, address contracts.InstanceAddress) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error) {
		active, err := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, templateID, address)
		if err != nil {
			if !strings.Contains(err.Error(), "multiple active contracts found") {
				return "", nil, err
			}
			parts := strings.Split(templateID, ":")
			if len(parts) != 3 {
				return "", nil, fmt.Errorf("invalid template ID for fallback lookup %q: %w", templateID, err)
			}
			lookupCtx, cancelLookup := context.WithTimeout(ctx, 20*time.Second)
			defer cancelLookup()

			ledgerEnd, endErr := participant.LedgerServices.State.GetLedgerEnd(lookupCtx, &ledgerv2.GetLedgerEndRequest{})
			if endErr != nil {
				return "", nil, fmt.Errorf("get ledger end for fallback lookup: %w", endErr)
			}
			stream, streamErr := participant.LedgerServices.State.GetActiveContracts(lookupCtx, &ledgerv2.GetActiveContractsRequest{
				ActiveAtOffset: ledgerEnd.GetOffset(),
				EventFormat: &ledgerv2.EventFormat{
					FiltersByParty: map[string]*ledgerv2.Filters{
						participant.PartyID: {
							Cumulative: []*ledgerv2.CumulativeFilter{{
								IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{
									TemplateFilter: &ledgerv2.TemplateFilter{
										TemplateId: &ledgerv2.Identifier{
											PackageId:  parts[0],
											ModuleName: parts[1],
											EntityName: parts[2],
										},
										IncludeCreatedEventBlob: true,
									},
								},
							}},
						},
					},
					Verbose: true,
				},
			})
			if streamErr != nil {
				return "", nil, fmt.Errorf("get active contracts for fallback lookup: %w", streamErr)
			}
			defer stream.CloseSend()

			var latestMatch *ledgerv2.ActiveContract
			for {
				resp, recvErr := stream.Recv()
				if recvErr != nil {
					if lookupCtx.Err() != nil {
						return "", nil, fmt.Errorf("fallback lookup timed out while reading active contracts for %s: %w", address.String(), lookupCtx.Err())
					}
					if errors.Is(recvErr, io.EOF) {
						break
					}
					return "", nil, fmt.Errorf("receive active contracts for fallback lookup: %w", recvErr)
				}
				entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
				if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
					continue
				}
				created := entry.ActiveContract.GetCreatedEvent()
				createArgs := created.GetCreateArguments()
				if createArgs == nil {
					continue
				}
				var instanceIDText string
				for _, f := range createArgs.GetFields() {
					if f.GetLabel() == "instanceId" {
						instanceIDText = f.GetValue().GetText()
						break
					}
				}
				if instanceIDText == "" || len(created.GetSignatories()) != 1 {
					continue
				}
				gotAddr := contracts.InstanceID(instanceIDText).RawInstanceAddress(types.PARTY(created.GetSignatories()[0])).InstanceAddress()
				if gotAddr != address {
					continue
				}
				if latestMatch == nil || created.GetOffset() > latestMatch.GetCreatedEvent().GetOffset() {
					latestMatch = entry.ActiveContract
				}
			}
			if latestMatch == nil {
				return "", nil, err
			}
			active = latestMatch
		}
		return types.CONTRACT_ID(active.GetCreatedEvent().GetContractId()), convertToDisclosedContract(active), nil
	}

	registryAdmin, err := getRegistryAdmin(withRetry(ctx, *scanRetryTimeout), tokenMetadataClient)
	if err != nil {
		fatalf("get registry admin: %v", err)
	}
	feeTokenInstrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(registryAdmin),
		Id:    types.TEXT("Amulet"),
	}

	feeTokenHoldingCID, disclosedFeeTokenHolding, feeTokenHoldingAmount, err := ensureAmuletFeeTokenHolding(
		ctx,
		participant,
		tokenMetadataClient,
		transferInstructionClient,
		scanProxyClient,
		feeTokenInstrument,
		*partyID,
		"10000000000",
	)
	if err != nil {
		fatalf("ensure amulet fee token holding: %v", err)
	}
	logger.Info().
		Str("holdingCID", feeTokenHoldingCID).
		Str("amount", feeTokenHoldingAmount).
		Msg("minted dedicated amulet holding for fee payment")
	transferFactoryCID, transferFactoryDisclosures, choiceContext, err := getTransferFactory(withRetry(ctx, *scanRetryTimeout), transferInstructionClient, registryAdmin, *partyID, *partyID)
	if err != nil {
		fatalf("get transfer factory: %v", err)
	}
	transferFactoryContextValues, err := transferFactoryContextFromChoiceContext(choiceContext)
	if err != nil {
		fatalf("decode transfer factory choice context: %v", err)
	}

	committeeVerifierCCV := strings.TrimSpace(*committeeVerifier)
	if committeeVerifierCCV == "" {
		fatalf("missing committee verifier CCV (set -committee-verifier or STAGING_CANTON_TO_EVM_COMMITTEE_VERIFIER)")
	}

	receiverBytes, err := decodeEVMAddress(*receiverHex)
	if err != nil {
		fatalf("decode receiver: %v", err)
	}

	outgoing := oapiEDSCommon.Message{
		DestinationChainSelector: fmt.Sprintf("%d", *dstSelector),
		FeeToken: oapiEDSCommon.InstrumentId{
			Admin: oapiEDSCommon.PartyId(feeTokenInstrument.Admin),
			Id:    string(feeTokenInstrument.Id),
		},
		Executor: struct {
			Address *oapiEDSCommon.RawOrHashedAddress `json:"address,omitempty"`
			Type    oapiEDSCommon.MessageExecutorType `json:"type"`
		}{Type: oapiEDSCommon.Empty},
		Payload:  hex.EncodeToString([]byte(*messageData)),
		Receiver: hex.EncodeToString(receiverBytes),
	}

	edsHTTP := &http.Client{Timeout: 15 * time.Second}
	var sendEDS *stagingeds.SendEDSOutcome
	var factoryPreload *perPartyRouterFactoryPreload
	if splitGriddleEDS {
		waitForStagingVPN(&logger, *vpnSwitchWait, "smartcontracts_com", "Switch to SmartContracts.com VPN for hosted EDS (*.griddle.sh). Waiting before factory prefetch and send disclosure collection…")
		if strings.TrimSpace(*perPartyRouterFactoryCIDFlag) == "" {
			p, preloadErr := fetchPerPartyRouterFactoryPreloadFromHostedEDS(ctx, edsBaseTrim, *partyID, &logger)
			if preloadErr != nil {
				logger.Warn().Err(preloadErr).Msg("PerPartyRouterFactory EDS prefetch failed; will try ledger on Chainlink Legacy VPN or use -per-party-router-factory-cid")
			} else {
				factoryPreload = p
			}
		}
		sendEDS, err = stagingeds.CollectSendDisclosures(ctx, edsBaseTrim, edsHTTP, outgoing, []string{committeeVerifierCCV}, nil, nil)
		if err != nil {
			fatalf("collect send disclosures from EDS: %v", err)
		}
		waitForStagingVPN(&logger, *vpnSwitchWait, "chainlink_legacy", "Switch back to Chainlink Legacy VPN for per-party router (CreateRouter / ACS) and send submission…")
	} else {
		if *vpnSwitchWait > 0 {
			logger.Info().Dur("wait", *vpnSwitchWait).Msg("pause: ensure VPN reaches Canton (gRPC, validator, scan-proxy) and EDS")
			time.Sleep(*vpnSwitchWait)
		}
		sendEDS, err = stagingeds.CollectSendDisclosures(ctx, edsBaseTrim, edsHTTP, outgoing, []string{committeeVerifierCCV}, nil, nil)
		if err != nil {
			fatalf("collect send disclosures from EDS: %v", err)
		}
	}

	routerAddress, err := ensurePerPartyRouter(ctx, participant, *partyID, routerInstanceID,
		strings.TrimSpace(*perPartyRouterFactoryCIDFlag),
		strings.TrimSpace(*perPartyRouterFactoryAddrFlag),
		edsBaseTrim,
		factoryPreload,
		splitGriddleEDS,
		&logger,
	)
	if err != nil {
		fatalf("ensure per-party router: %v", err)
	}
	routerCID, disclosedRouter, err := resolveDisclosedByAddress(perpartyrouter.PerPartyRouter{}.GetTemplateID(), routerAddress)
	if err != nil {
		fatalf("resolve per-party router disclosed contract: %v", err)
	}
	logger.Info().
		Str("routerCID", string(routerCID)).
		Str("routerAddress", routerAddress.String()).
		Msg("resolved per-party router")

	var executorCID types.CONTRACT_ID
	var disclosedExecutor *ledgerv2.DisclosedContract
	if sendEDS.ExecutorInput != nil {
		executorCID = sendEDS.ExecutorInput.ExecutorCid
		disclosedExecutor = stagingeds.FindDisclosedContractByContractID(sendEDS.DisclosedContracts, string(executorCID))
		if disclosedExecutor == nil {
			fatalf("EDS send disclosures did not include a DisclosedContract for executor cid %s (ccip/devenv loads executor via GetExecutorSendDisclosure)", executorCID)
		}
		logger.Info().Str("executorCID", string(executorCID)).Msg("resolved executor from EDS send disclosures")
	} else {
		var resErr error
		executorCID, disclosedExecutor, resErr = resolveDisclosedByAddress(executorbinding.Executor{}.GetTemplateID(), contracts.HexToInstanceAddress(*executorAddr))
		if resErr != nil {
			fatalf("resolve executor disclosed contract: %v", resErr)
		}
	}

	senderCID, disclosedSender, senderAddress, err := ensureCCIPSender(ctx, participant, resolveDisclosedByAddress, *partyID, *senderInstanceID, &logger)
	if err != nil {
		fatalf("ensure ccip sender: %v", err)
	}
	logger.Info().
		Str("senderInstanceID", *senderInstanceID).
		Str("senderCID", string(senderCID)).
		Str("senderAddress", senderAddress.String()).
		Msg("resolved ccip sender")

	execInput := sendEDS.ExecutorInput
	if execInput == nil {
		execInput = &ccipsender.ExecutorInput{
			ExecutorCid: executorCID,
			ExecutorExtraContext: splice_api_token_metadata_v1.ChoiceContext{
				Values: map[string]splice_api_token_metadata_v1.AnyValue{},
			},
		}
	}

	sendArgs := ccipsender.Send{
		Context:                  sendEDS.SendContext,
		RouterCid:                routerCID,
		DestinationChainSelector: types.NUMERIC(fmt.Sprintf("%d", *dstSelector)),
		Message: ccipclient.Canton2AnyMessage{
			Receiver: types.TEXT(hex.EncodeToString(receiverBytes)),
			Payload:  types.TEXT(hex.EncodeToString([]byte(*messageData))),
			FeeToken: feeTokenInstrument,
			ExtraArgs: ccipclient.ExtraArgs{V3: &ccipclient.GenericExtraArgsV3{
				GasLimit: types.INT64(*executionGasLimit),
				Ccvs:     sendEDS.CcvExtraArgs,
				Executor: ccipclient.ExecutorExtraArg{ExecutorUseDefault: &ccipclient.ExecutorUseDefault{
					ExecutorArgs: types.TEXT(""),
				}},
				TokenReceiver: types.TEXT(""),
				TokenArgs:     types.TEXT(""),
			}},
		},
		FeeTokenInput: ccipsender.FeeTokenInput{
			SenderInputCids:         []types.CONTRACT_ID{types.CONTRACT_ID(feeTokenHoldingCID)},
			FeeTokenTransferFactory: types.CONTRACT_ID(transferFactoryCID),
			FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
				Context: splice_api_token_metadata_v1.ChoiceContext{Values: transferFactoryContextValues},
				Meta:    splice_api_token_metadata_v1.Metadata{Values: map[string]types.TEXT{}},
			},
		},
		CcvSendInputs: sendEDS.CcvSendInputs,
		ExecutorInput: execInput,
	}

	disclosedContracts := []*ledgerv2.DisclosedContract{
		disclosedSender,
		disclosedExecutor,
		disclosedRouter,
		disclosedFeeTokenHolding,
	}
	disclosedContracts = append(disclosedContracts, sendEDS.DisclosedContracts...)
	disclosedContracts = append(disclosedContracts, transferFactoryDisclosures...)
	disclosedContracts = dedupDisclosedContracts(disclosedContracts)
	for _, dc := range disclosedContracts {
		if dc == nil || dc.GetContractId() == "" {
			fatalf("empty disclosed contract ID before send")
		}
	}
	activeContractIDs, err := activeContractIDSet(ctx, participant)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to build active contract ID set before send submit")
	} else {
		for _, dc := range disclosedContracts {
			if _, ok := activeContractIDs[dc.GetContractId()]; ok {
				continue
			}
			logger.Warn().
				Str("contractID", dc.GetContractId()).
				Str("templateID", identifierString(dc.GetTemplateId())).
				Msg("disclosed contract already inactive before send submit")
		}
		logInactiveKeyContract(logger, activeContractIDs, "sender", string(senderCID))
		logInactiveKeyContract(logger, activeContractIDs, "router", string(routerCID))
		logInactiveKeyContract(logger, activeContractIDs, "executor", string(executorCID))
		logInactiveKeyContract(logger, activeContractIDs, "feeHolding", feeTokenHoldingCID)
		logInactiveKeyContract(logger, activeContractIDs, "transferFactory", transferFactoryCID)
		for i, input := range sendEDS.CcvSendInputs {
			logInactiveKeyContract(logger, activeContractIDs, fmt.Sprintf("ccv[%d]", i), string(input.CcvCid))
		}
	}

	sendRes, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
					TemplateId:     disclosedSender.GetTemplateId(),
					ContractId:     string(senderCID),
					Choice:         "Send",
					ChoiceArgument: ledger.MapToValue(sendArgs.ToMap()),
				}},
			}},
			ActAs:              []string{*partyID},
			DisclosedContracts: disclosedContracts,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("submit ccip send via ccip sender: %v", err)
		if strings.Contains(err.Error(), "No price found for") {
			msg += "; CCIP OnRamp requires FeeQuoter.usdPerToken to include the fee instrument (this script uses registry-admin + id \"Amulet\"). Add a TimestampedPrice for that instrument on the staging FeeQuoter (same as chain ops / MCMS), or align the send fee token with a priced instrument."
		}
		fatalf("%s", msg)
	}

	var messageID string
	var seqNo uint64
	for _, event := range sendRes.GetTransaction().GetEvents() {
		created := event.GetCreated()
		if created == nil || created.GetTemplateId() == nil || created.GetTemplateId().GetEntityName() != "CCIPMessageSent" {
			continue
		}
		for _, f := range created.GetCreateArguments().GetFields() {
			if f.GetLabel() != "event" || f.GetValue().GetRecord() == nil {
				continue
			}
			for _, eventField := range f.GetValue().GetRecord().GetFields() {
				switch eventField.GetLabel() {
				case "encodedMessage":
					encodedMessage, decErr := hex.DecodeString(eventField.GetValue().GetText())
					if decErr != nil {
						fatalf("decode encodedMessage from CCIPMessageSent event: %v", decErr)
					}
					decodedMessagePtr, decErr := protocol.DecodeMessage(encodedMessage)
					if decErr != nil {
						fatalf("decode protocol message from encodedMessage: %v", decErr)
					}
					seqNo = uint64(decodedMessagePtr.SequenceNumber)
					messageID = hex.EncodeToString(gethcrypto.Keccak256(encodedMessage))
				case "messageId":
					if messageID == "" {
						messageID = eventField.GetValue().GetText()
					}
				}
			}
		}
		break
	}
	if messageID == "" {
		fatalf("no CCIPMessageSent event found in sender transaction")
	}

	logger.Info().
		Uint64("srcSelector", *srcSelector).
		Uint64("dstSelector", *dstSelector).
		Str("receiver", *receiverHex).
		Str("data", *messageData).
		Int64("executionGasLimit", *executionGasLimit).
		Int64("finalityConfig", *finalityConfig).
		Str("updateID", sendRes.GetTransaction().GetUpdateId()).
		Uint64("sequenceNumber", seqNo).
		Str("messageID", "0x"+strings.TrimPrefix(messageID, "0x")).
		Msg("Canton -> EVM send submitted")
}

// resolvePerPartyRouterFactoryForCreateRouter matches ccip/devenv: primary source is hosted EDS
// GetPerPartyRouterFactoryDisclosure (see Chain.GetPerPartyRouterFactoryDisclosure + DeployPerPartyRouter
// in manual_execution.go). When preloaded is set (SmartContracts VPN phase), uses it before live EDS.
// Falls back to explicit CID, factory instance address, then ACS template scan.
func resolvePerPartyRouterFactoryForCreateRouter(
	ctx context.Context,
	participant cldfcanton.Participant,
	partyID string,
	factoryCIDOverride string,
	factoryAddrOverride string,
	edsBaseURL string,
	preloaded *perPartyRouterFactoryPreload,
	suppressLiveEDS bool,
	logger *zerolog.Logger,
) (factoryCID string, disclosed []*ledgerv2.DisclosedContract, err error) {
	if strings.TrimSpace(factoryCIDOverride) != "" {
		factoryCID = strings.TrimSpace(factoryCIDOverride)
		dc, err := getDisclosedContractByID(ctx, participant, factoryCID)
		if err != nil {
			return "", nil, fmt.Errorf("load factory by -per-party-router-factory-cid / STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_CID: %w", err)
		}
		return factoryCID, []*ledgerv2.DisclosedContract{dc}, nil
	}

	if preloaded != nil && preloaded.ContractID != "" && len(preloaded.Disclosed) > 0 {
		if logger != nil {
			logger.Info().Str("factoryContractId", preloaded.ContractID).Msg("using pre-fetched PerPartyRouterFactory disclosure (EDS phase)")
		}
		return preloaded.ContractID, preloaded.Disclosed, nil
	}

	edsBase := strings.TrimSpace(edsBaseURL)
	if !suppressLiveEDS && edsBase != "" {
		ccipEDS, cerr := oapiCCIP.NewClientWithResponses(edsBase, oapiCCIP.WithHTTPClient(&http.Client{Timeout: 20 * time.Second}))
		if cerr != nil {
			return "", nil, fmt.Errorf("create CCIP EDS client for PerPartyRouterFactory: %w", cerr)
		}
		fd, derr := eds.GetPerPartyRouterFactoryDisclosure(ctx, ccipEDS, partyID)
		if derr == nil && fd != nil && fd.ContractId != "" && len(fd.DisclosedContracts) > 0 {
			if logger != nil {
				logger.Info().Str("factoryContractId", fd.ContractId).Msg("resolved PerPartyRouterFactory via hosted EDS (ccip/devenv/manual_execution.go pattern)")
			}
			return fd.ContractId, fd.DisclosedContracts, nil
		}
		if logger != nil && derr != nil {
			logger.Warn().Err(derr).Msg("hosted EDS PerPartyRouterFactory disclosure failed; trying ledger fallbacks")
		}
	}

	if addr := strings.TrimSpace(factoryAddrOverride); addr != "" {
		ia := contracts.HexToInstanceAddress(addr)
		for _, tid := range perPartyRouterFactoryTemplateIDs() {
			active, ferr := contract.FindActiveContractByInstanceAddress(ctx, participant.LedgerServices.State, participant.PartyID, tid, ia)
			if ferr == nil && active != nil && active.GetCreatedEvent() != nil {
				return active.GetCreatedEvent().GetContractId(), []*ledgerv2.DisclosedContract{convertToDisclosedContract(active)}, nil
			}
		}
	}

	active, err := findActivePerPartyRouterFactory(ctx, participant)
	if err != nil {
		return "", nil, err
	}
	if active == nil || active.GetCreatedEvent() == nil {
		if suppressLiveEDS {
			return "", nil, fmt.Errorf("no PerPartyRouterFactory on ledger after split VPN EDS phase (set STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_CID / -per-party-router-factory-address, or ensure SmartContracts VPN EDS prefetch succeeded)")
		}
		return "", nil, fmt.Errorf("no PerPartyRouterFactory (EDS at %q failed or unreachable; set STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_CID, -per-party-router-factory-address from address_refs, or fix VPN/DNS for EDS)", edsBase)
	}
	return active.GetCreatedEvent().GetContractId(), []*ledgerv2.DisclosedContract{convertToDisclosedContract(active)}, nil
}

func findActivePerPartyRouterFactory(ctx context.Context, participant cldfcanton.Participant) (*ledgerv2.ActiveContract, error) {
	for _, tid := range perPartyRouterFactoryTemplateIDs() {
		pkg, mod, ent, err := contracts.ParseTemplateIDFromString(tid)
		if err != nil {
			continue
		}
		active, err := scanFirstActiveContractByTemplate(ctx, participant, pkg, mod, ent)
		if err != nil {
			return nil, err
		}
		if active != nil {
			return active, nil
		}
	}
	return nil, nil
}

func scanFirstActiveContractByTemplate(ctx context.Context, participant cldfcanton.Participant, packageID, moduleName, entityName string) (*ledgerv2.ActiveContract, error) {
	offsetResp, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, err
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: offsetResp.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{{
						IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{TemplateFilter: &ledgerv2.TemplateFilter{
							TemplateId: &ledgerv2.Identifier{
								PackageId:  packageID,
								ModuleName: moduleName,
								EntityName: entityName,
							},
							IncludeCreatedEventBlob: true,
						}},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return nil, err
	}
	defer stream.CloseSend()
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		active, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || active.ActiveContract == nil || active.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		return active.ActiveContract, nil
	}
}

// ensurePerPartyRouter mirrors ccip/devenv.DeployPerPartyRouter (manual_execution.go): EDS-backed
// factory disclosure + CreateRouter, idempotent if the router already exists.
func ensurePerPartyRouter(
	ctx context.Context,
	participant cldfcanton.Participant,
	partyID string,
	routerInstanceID string,
	perPartyRouterFactoryCIDOverride string,
	perPartyRouterFactoryAddrOverride string,
	edsBaseURL string,
	factoryPreload *perPartyRouterFactoryPreload,
	suppressLiveEDS bool,
	logger *zerolog.Logger,
) (contracts.InstanceAddress, error) {
	existingRouterAddr, err := findExistingRouterAddress(ctx, participant, partyID, routerInstanceID)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("find existing router: %w", err)
	}
	if existingRouterAddr != (contracts.InstanceAddress{}) {
		return existingRouterAddr, nil
	}

	factoryCID, factoryDisclosed, err := resolvePerPartyRouterFactoryForCreateRouter(ctx, participant, partyID, perPartyRouterFactoryCIDOverride, perPartyRouterFactoryAddrOverride, edsBaseURL, factoryPreload, suppressLiveEDS, logger)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("per-party router not found for %s (instance %q); resolve PerPartyRouterFactory: %w", partyID, routerInstanceID, err)
	}
	if len(factoryDisclosed) == 0 || factoryDisclosed[0].GetTemplateId() == nil {
		return contracts.InstanceAddress{}, fmt.Errorf("internal: empty factory disclosure for CreateRouter")
	}
	factoryTemplateID := factoryDisclosed[0].GetTemplateId()
	res, submitErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{Exercise: &ledgerv2.ExerciseCommand{
					TemplateId:     factoryTemplateID,
					ContractId:     factoryCID,
					Choice:         "CreateRouter",
					ChoiceArgument: ledger.MapToValue(perpartyrouter.CreateRouter{PartyOwner: types.PARTY(partyID), InstanceId: types.TEXT(routerInstanceID)}),
				}},
			}},
			ActAs:              []string{partyID},
			DisclosedContracts: factoryDisclosed,
		},
	})
	if submitErr != nil && logger != nil {
		logger.Warn().Err(submitErr).Msg("CreateRouter submit failed or duplicate (devenv ignores this); re-checking ACS for PerPartyRouter")
	}
	if submitErr == nil {
		for _, event := range res.GetTransaction().GetEvents() {
			if e := event.GetCreated(); e != nil && e.GetTemplateId() != nil && e.GetTemplateId().GetEntityName() == "PerPartyRouter" {
				if logger != nil {
					logger.Info().
						Str("routerContractId", e.GetContractId()).
						Str("routerInstanceId", routerInstanceID).
						Msg("created per-party router on ledger")
				}
				break
			}
		}
	}

	existingRouterAddr, err = findExistingRouterAddress(ctx, participant, partyID, routerInstanceID)
	if err != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("reload router after CreateRouter: %w", err)
	}
	if existingRouterAddr != (contracts.InstanceAddress{}) {
		return existingRouterAddr, nil
	}
	if submitErr != nil {
		return contracts.InstanceAddress{}, fmt.Errorf("create per-party router: %w", submitErr)
	}
	return contracts.InstanceAddress{}, fmt.Errorf("CreateRouter reported success but PerPartyRouter not found for party %s instance %q", partyID, routerInstanceID)
}

func findExistingRouterAddress(ctx context.Context, participant cldfcanton.Participant, partyID, routerInstanceID string) (contracts.InstanceAddress, error) {
	for _, tid := range perPartyRouterTemplateIDs() {
		pkg, mod, ent, err := contracts.ParseTemplateIDFromString(tid)
		if err != nil {
			continue
		}
		addr, err := scanPerPartyRouterForParty(ctx, participant, partyID, routerInstanceID, pkg, mod, ent)
		if err != nil {
			return contracts.InstanceAddress{}, err
		}
		if addr != (contracts.InstanceAddress{}) {
			return addr, nil
		}
	}
	return contracts.InstanceAddress{}, nil
}

func scanPerPartyRouterForParty(ctx context.Context, participant cldfcanton.Participant, partyID, routerInstanceID, packageID, moduleName, entityName string) (contracts.InstanceAddress, error) {
	offsetResp, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return contracts.InstanceAddress{}, err
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: offsetResp.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{{
						IdentifierFilter: &ledgerv2.CumulativeFilter_TemplateFilter{TemplateFilter: &ledgerv2.TemplateFilter{
							TemplateId: &ledgerv2.Identifier{
								PackageId:  packageID,
								ModuleName: moduleName,
								EntityName: entityName,
							},
							IncludeCreatedEventBlob: true,
						}},
					}},
				},
			},
			Verbose: true,
		},
	})
	if err != nil {
		return contracts.InstanceAddress{}, err
	}
	defer stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return contracts.InstanceAddress{}, nil
		}
		if err != nil {
			return contracts.InstanceAddress{}, err
		}
		active, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || active.ActiveContract == nil || active.ActiveContract.GetCreatedEvent() == nil {
			continue
		}

		created := active.ActiveContract.GetCreatedEvent()
		fields := created.GetCreateArguments().GetFields()
		var instanceIDText string
		var partyOwner string
		for _, field := range fields {
			switch field.GetLabel() {
			case "instanceId":
				instanceIDText = field.GetValue().GetText()
			case "partyOwner":
				partyOwner = field.GetValue().GetParty()
			}
		}
		if partyOwner != partyID || instanceIDText == "" || len(created.GetSignatories()) != 1 {
			continue
		}
		if routerInstanceID != "" && instanceIDText != routerInstanceID {
			continue
		}

		return contracts.InstanceID(instanceIDText).RawInstanceAddress(types.PARTY(created.GetSignatories()[0])).InstanceAddress(), nil
	}
}

func createdContractIDFromTransaction(tx *ledgerv2.Transaction, entityName string) string {
	if tx == nil {
		return ""
	}
	for _, event := range tx.GetEvents() {
		if e := event.GetCreated(); e != nil && e.GetTemplateId() != nil && e.GetTemplateId().GetEntityName() == entityName {
			return e.GetContractId()
		}
	}
	return ""
}

func ensureCCIPSender(
	ctx context.Context,
	participant cldfcanton.Participant,
	resolve func(string, contracts.InstanceAddress) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, error),
	partyID string,
	instanceID string,
	logger *zerolog.Logger,
) (types.CONTRACT_ID, *ledgerv2.DisclosedContract, contracts.InstanceAddress, error) {
	senderAddr := contracts.InstanceID(instanceID).RawInstanceAddress(types.PARTY(partyID)).InstanceAddress()
	senderCID, disclosedSender, err := resolve(ccipsender.CCIPSender{}.GetTemplateID(), senderAddr)
	if err == nil && senderCID != "" {
		return senderCID, disclosedSender, senderAddr, nil
	}
	if err != nil && !strings.Contains(err.Error(), "no active contract found") {
		return "", nil, contracts.InstanceAddress{}, fmt.Errorf("resolve existing ccip sender: %w", err)
	}
	// Missing sender on ACS (or empty cid): ledger Create, same pattern as integration-tests/ccip/ccip_send_test.go.
	// v2 CreateCommand expects *Record, not *Value (unlike ExerciseCommand.ChoiceArgument).
	ccipSenderCreateArgs := &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
		{Label: "instanceId", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: instanceID}}},
		{Label: "owner", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: partyID}}},
	}}

	var lastCreateErr error
	for _, tidStr := range ccipsenderTemplateIDs() {
		tplID, terr := templateIDFromString(tidStr)
		if terr != nil {
			continue
		}
		res, createErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
			Commands: &ledgerv2.Commands{
				CommandId: uuid.NewString(),
				Commands: []*ledgerv2.Command{{
					Command: &ledgerv2.Command_Create{Create: &ledgerv2.CreateCommand{
						TemplateId:      tplID,
						CreateArguments: ccipSenderCreateArgs,
					}},
				}},
				ActAs: []string{partyID},
			},
		})
		if createErr != nil {
			lastCreateErr = createErr
			if logger != nil {
				logger.Warn().Err(createErr).Str("templateID", tidStr).Msg("CCIPSender Create failed (trying alternate package id if any)")
			}
			senderCID2, disclosed2, err2 := resolve(ccipsender.CCIPSender{}.GetTemplateID(), senderAddr)
			if err2 == nil && senderCID2 != "" {
				if logger != nil {
					logger.Warn().Err(createErr).Msg("CCIPSender Create failed but sender already exists; using ACS resolution")
				}
				return senderCID2, disclosed2, senderAddr, nil
			}
			continue
		}
		newCID := createdContractIDFromTransaction(res.GetTransaction(), "CCIPSender")
		if newCID == "" {
			lastCreateErr = fmt.Errorf("create transaction had no CCIPSender Created event for template %s", tidStr)
			continue
		}
		if logger != nil {
			logger.Info().Str("senderContractId", newCID).Str("senderInstanceId", instanceID).Msg("created CCIPSender on ledger")
		}
		dc, gerr := getDisclosedContractByID(ctx, participant, newCID)
		if gerr != nil {
			return "", nil, senderAddr, fmt.Errorf("load disclosed ccip sender after create: %w", gerr)
		}
		return types.CONTRACT_ID(newCID), dc, senderAddr, nil
	}
	if lastCreateErr != nil {
		return "", nil, senderAddr, fmt.Errorf("create ccip sender: %w", lastCreateErr)
	}
	return "", nil, senderAddr, fmt.Errorf("create ccip sender: no template id worked")
}

func disclosedContractToProto(contract oapiEDSCommon.DisclosedContract) (*ledgerv2.DisclosedContract, error) {
	id, err := templateIDFromString(contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}
	return &ledgerv2.DisclosedContract{
		TemplateId:       id,
		ContractId:       contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   contract.SynchronizerId,
	}, nil
}

func templateIDFromString(s string) (*ledgerv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}
	return &ledgerv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func transferFactoryContextFromChoiceContext(choiceContext *ledgerv2.Value) (map[string]splice_api_token_metadata_v1.AnyValue, error) {
	if choiceContext == nil || choiceContext.GetRecord() == nil {
		return nil, fmt.Errorf("choice context is nil or not a record")
	}
	values := make(map[string]splice_api_token_metadata_v1.AnyValue)
	for _, field := range choiceContext.GetRecord().GetFields() {
		if field.GetLabel() != "values" || field.GetValue().GetTextMap() == nil {
			continue
		}
		for _, entry := range field.GetValue().GetTextMap().GetEntries() {
			if v := entry.GetValue().GetVariant(); v != nil {
				cid := types.CONTRACT_ID(v.GetValue().GetContractId())
				values[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVContractId: &cid}
				continue
			}
			if txt := entry.GetValue().GetText(); txt != "" {
				t := types.TEXT(txt)
				values[entry.GetKey()] = splice_api_token_metadata_v1.AnyValue{AVText: &t}
			}
		}
	}
	return values, nil
}

func withRetry(parent context.Context, timeout time.Duration) context.Context {
	if timeout <= 0 {
		return parent
	}
	ctx, _ := context.WithTimeout(parent, timeout)
	return ctx
}

func retryScanCall[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		res, err := fn(ctx)
		if err == nil {
			return res, nil
		}
		if !isRetryableScanErr(err) {
			return zero, err
		}
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("scan-proxy retry timeout: %w", err)
		case <-ticker.C:
		}
	}
}

func isRetryableScanErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unexpected status code: 502") ||
		strings.Contains(msg, "unexpected status code: 503") ||
		strings.Contains(msg, "Failed to reach consensus") ||
		strings.Contains(msg, "bad gateway") ||
		strings.Contains(msg, "Bad Gateway")
}

func getRegistryAdmin(ctx context.Context, metadataClient tokenMetadataV1.ClientWithResponsesInterface) (string, error) {
	registryAdminFallback := stagingenv.String("", "STAGING_CANTON_TO_EVM_REGISTRY_ADMIN", "STAGING_CANTON_REGISTRY_ADMIN")
	return retryScanCall(ctx, func(ctx context.Context) (string, error) {
		registryInfoResponse, err := metadataClient.GetRegistryInfoWithResponse(ctx)
		if err != nil {
			if registryAdminFallback != "" {
				return registryAdminFallback, nil
			}
			return "", fmt.Errorf("error getting registry info: %w", err)
		}
		if registryInfoResponse.StatusCode() != http.StatusOK {
			if registryAdminFallback != "" {
				return registryAdminFallback, nil
			}
			return "", fmt.Errorf("unexpected status code: %d: %v", registryInfoResponse.StatusCode(), registryInfoResponse.Body)
		}
		return registryInfoResponse.JSON200.AdminId, nil
	})
}

func getTransferFactory(ctx context.Context, transferInstructionClient transferInstructionV1.ClientWithResponsesInterface, registryAdmin, sender, receiver string) (string, []*ledgerv2.DisclosedContract, *ledgerv2.Value, error) {
	res, err := retryScanCall(ctx, func(ctx context.Context) (struct {
		factoryID   string
		disclosures []*ledgerv2.DisclosedContract
		choiceCtx   *ledgerv2.Value
	}, error) {
		transferFactoryResponse, err := transferInstructionClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
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
			return struct {
				factoryID   string
				disclosures []*ledgerv2.DisclosedContract
				choiceCtx   *ledgerv2.Value
			}{}, fmt.Errorf("error getting transfer factory response: %w", err)
		}
		if transferFactoryResponse.StatusCode() != http.StatusOK {
			return struct {
				factoryID   string
				disclosures []*ledgerv2.DisclosedContract
				choiceCtx   *ledgerv2.Value
			}{}, fmt.Errorf("unexpected status code: %d: %v", transferFactoryResponse.StatusCode(), transferFactoryResponse.Body)
		}

		var disclosedContracts []*ledgerv2.DisclosedContract
		for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
			disclosedContract, err := disclosedContractToProto(oapiEDSCommon.DisclosedContract{
				TemplateId:       contract.TemplateId,
				ContractId:       contract.ContractId,
				CreatedEventBlob: contract.CreatedEventBlob,
				SynchronizerId:   contract.SynchronizerId,
			})
			if err != nil {
				return struct {
					factoryID   string
					disclosures []*ledgerv2.DisclosedContract
					choiceCtx   *ledgerv2.Value
				}{}, fmt.Errorf("failed to convert transfer factory disclosed contract: %w", err)
			}
			disclosedContracts = append(disclosedContracts, disclosedContract)
		}

		choiceContext, err := choiceContextFromData(transferFactoryResponse.JSON200.ChoiceContext.ChoiceContextData)
		if err != nil {
			return struct {
				factoryID   string
				disclosures []*ledgerv2.DisclosedContract
				choiceCtx   *ledgerv2.Value
			}{}, fmt.Errorf("failed to convert choice context: %w", err)
		}

		return struct {
			factoryID   string
			disclosures []*ledgerv2.DisclosedContract
			choiceCtx   *ledgerv2.Value
		}{
			factoryID:   transferFactoryResponse.JSON200.FactoryId,
			disclosures: disclosedContracts,
			choiceCtx:   choiceContext,
		}, nil
	})
	if err != nil {
		return "", nil, nil, err
	}
	return res.factoryID, res.disclosures, res.choiceCtx, nil
}

func getAmuletRulesContract(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (string, *ledgerv2.DisclosedContract, error) {
	res, err := retryScanCall(ctx, func(ctx context.Context) (struct {
		dsoPartyID string
		contract   *ledgerv2.DisclosedContract
	}, error) {
		// Avoid the broader /dso aggregate when this flow only needs dso_party_id and amulet_rules.
		dsoPartyIDResponse, err := scanProxyClient.GetDsoPartyIdWithResponse(ctx)
		if err != nil {
			return struct {
				dsoPartyID string
				contract   *ledgerv2.DisclosedContract
			}{}, fmt.Errorf("error getting dso party id response: %w", err)
		}
		if dsoPartyIDResponse.StatusCode() != http.StatusOK {
			return struct {
				dsoPartyID string
				contract   *ledgerv2.DisclosedContract
			}{}, fmt.Errorf("unexpected dso party id status code: %d: %v", dsoPartyIDResponse.StatusCode(), dsoPartyIDResponse.Body)
		}

		amuletRulesResponse, err := scanProxyClient.GetAmuletRulesWithResponse(ctx)
		if err != nil {
			return struct {
				dsoPartyID string
				contract   *ledgerv2.DisclosedContract
			}{}, fmt.Errorf("error getting amulet rules response: %w", err)
		}
		if amuletRulesResponse.StatusCode() != http.StatusOK {
			return struct {
				dsoPartyID string
				contract   *ledgerv2.DisclosedContract
			}{}, fmt.Errorf("unexpected amulet rules status code: %d: %v", amuletRulesResponse.StatusCode(), amuletRulesResponse.Body)
		}
		amuletRules, err := contractWithStateToDisclosedContract(amuletRulesResponse.JSON200.AmuletRules)
		if err != nil {
			return struct {
				dsoPartyID string
				contract   *ledgerv2.DisclosedContract
			}{}, err
		}
		return struct {
			dsoPartyID string
			contract   *ledgerv2.DisclosedContract
		}{
			dsoPartyID: dsoPartyIDResponse.JSON200.DsoPartyId,
			contract:   amuletRules,
		}, nil
	})
	if err != nil {
		return "", nil, err
	}
	return res.dsoPartyID, res.contract, nil
}

func choiceContextFromData(choiceContextData map[string]any) (*ledgerv2.Value, error) {
	values, ok := choiceContextData["values"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("no values found in choice context")
	}

	var fields []*ledgerv2.TextMap_Entry
	for k, v := range values {
		f := v.(map[string]any)
		tag := f["tag"].(string)
		rawValue := f["value"]

		var value *ledgerv2.Value
		switch tag {
		case "AV_Text":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Text value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Text{Text: valueString}}
		case "AV_Int":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_Int value is not a number: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: int64(valueFloat)}}
		case "AV_Decimal":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Decimal value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: valueString}}
		case "AV_Bool":
			valueBool, ok := rawValue.(bool)
			if !ok {
				return nil, fmt.Errorf("AV_Bool value is not a bool: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Bool{Bool: valueBool}}
		case "AV_Date":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Date value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Date value is not a RFC3339 time: %s", valueString)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Date{Date: int32(t.Unix() / 86400)}}
		case "AV_Time":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_Time value is not a string: %T", rawValue)
			}
			t, err := time.Parse(time.RFC3339, valueString)
			if err != nil {
				return nil, fmt.Errorf("AV_Time value is not a RFC3339 time: %s", valueString)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Timestamp{Timestamp: t.UnixMicro()}}
		case "AV_RelTime":
			valueFloat, ok := rawValue.(float64)
			if !ok {
				return nil, fmt.Errorf("AV_RelTime value is not a number: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
				{Label: "microseconds", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Int64{Int64: int64(valueFloat)}}},
			}}}}
		case "AV_ContractId":
			valueString, ok := rawValue.(string)
			if !ok {
				return nil, fmt.Errorf("AV_ContractId value is not a string: %T", rawValue)
			}
			value = &ledgerv2.Value{Sum: &ledgerv2.Value_ContractId{ContractId: valueString}}
		default:
			return nil, fmt.Errorf("unimplemented tag: %v", tag)
		}

		fields = append(fields, &ledgerv2.TextMap_Entry{
			Key: k,
			Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Variant{Variant: &ledgerv2.Variant{
				Constructor: tag,
				Value:       value,
			}}},
		})
	}

	return &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{{
		Label: "values",
		Value: &ledgerv2.Value{Sum: &ledgerv2.Value_TextMap{TextMap: &ledgerv2.TextMap{Entries: fields}}},
	}}}}}, nil
}

func getFirstOpenMiningRound(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (*ledgerv2.DisclosedContract, error) {
	openMiningRoundResponse, err := scanProxyClient.GetOpenAndIssuingMiningRoundsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting open mining rounds response: %w", err)
	}
	if openMiningRoundResponse.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %v", openMiningRoundResponse.StatusCode(), openMiningRoundResponse.Body)
	}
	for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
		opensAt, err := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse opensAt %q: %w", round.Contract.Payload["opensAt"], err)
		}
		targetClosesAt, err := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse targetClosesAt %q: %w", round.Contract.Payload["targetClosesAt"], err)
		}
		if opensAt.Before(time.Now()) && targetClosesAt.After(time.Now()) {
			return contractWithStateToDisclosedContract(round)
		}
	}
	return nil, fmt.Errorf("failed to find open mining round contract")
}

// Kept local because the similar MintAMT helper lives under the separate
// integration-tests module; importing that from this root-module script would
// create the wrong dependency direction.
func mintAMT(
	ctx context.Context,
	participant cldfcanton.Participant,
	metadataClient tokenMetadataV1.ClientWithResponsesInterface,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	scanProxyClient scanProxy.ClientWithResponsesInterface,
	toParty string,
	amount string,
) (string, error) {
	registryAdmin, err := getRegistryAdmin(ctx, metadataClient)
	if err != nil {
		return "", fmt.Errorf("failed to get registry admin: %w", err)
	}
	_, amuletRulesContract, err := getAmuletRulesContract(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get amulet rules contract: %w", err)
	}
	_, disclosedContracts, _, err := getTransferFactory(ctx, transferInstructionClient, registryAdmin, registryAdmin, toParty)
	if err != nil {
		return "", fmt.Errorf("failed to get transfer factory: %w", err)
	}
	openMiningRoundContract, err := getFirstOpenMiningRound(ctx, scanProxyClient)
	if err != nil {
		return "", fmt.Errorf("failed to get open mining round: %w", err)
	}

	response, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &ledgerv2.SubmitAndWaitForTransactionRequest{
		Commands: &ledgerv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*ledgerv2.Command{{
				Command: &ledgerv2.Command_Exercise{
					Exercise: &ledgerv2.ExerciseCommand{
						TemplateId: amuletRulesContract.TemplateId,
						ContractId: amuletRulesContract.ContractId,
						Choice:     "AmuletRules_DevNet_Tap",
						ChoiceArgument: &ledgerv2.Value{Sum: &ledgerv2.Value_Record{Record: &ledgerv2.Record{Fields: []*ledgerv2.RecordField{
							{Label: "receiver", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Party{Party: toParty}}},
							{Label: "amount", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_Numeric{Numeric: amount}}},
							{Label: "openRound", Value: &ledgerv2.Value{Sum: &ledgerv2.Value_ContractId{ContractId: openMiningRoundContract.ContractId}}},
						}}}},
					},
				},
			}},
			ActAs:              []string{toParty},
			DisclosedContracts: dedupDisclosedContracts(append(disclosedContracts, amuletRulesContract, openMiningRoundContract)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to mint AMT: %w", err)
	}

	tokenHoldingCID := ""
	for _, event := range response.GetTransaction().GetEvents() {
		if created := event.GetCreated(); created != nil {
			tokenHoldingCID = created.GetContractId()
		}
	}
	return tokenHoldingCID, nil
}

func ensureAmuletFeeTokenHolding(
	ctx context.Context,
	participant cldfcanton.Participant,
	metadataClient tokenMetadataV1.ClientWithResponsesInterface,
	transferInstructionClient transferInstructionV1.ClientWithResponsesInterface,
	scanProxyClient scanProxy.ClientWithResponsesInterface,
	instrument splice_api_token_holding_v1.InstrumentId,
	partyID string,
	mintAmount string,
) (string, *ledgerv2.DisclosedContract, string, error) {
	mintedCID, err := mintAMT(ctx, participant, metadataClient, transferInstructionClient, scanProxyClient, partyID, mintAmount)
	if err != nil {
		return "", nil, "", fmt.Errorf("mint amulet fee tokens: %w", err)
	}
	disclosedHolding, err := getDisclosedContractByID(ctx, participant, mintedCID)
	if err != nil {
		return "", nil, "", fmt.Errorf("get disclosed minted amulet holding: %w", err)
	}
	return mintedCID, disclosedHolding, mintAmount, nil
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

func getDisclosedContractByID(ctx context.Context, participant cldfcanton.Participant, contractID string) (*ledgerv2.DisclosedContract, error) {
	offset, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger end: %w", err)
	}
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: offset.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersByParty: map[string]*ledgerv2.Filters{
				participant.PartyID: {
					Cumulative: []*ledgerv2.CumulativeFilter{{
						IdentifierFilter: &ledgerv2.CumulativeFilter_WildcardFilter{
							WildcardFilter: &ledgerv2.WildcardFilter{IncludeCreatedEventBlob: true},
						},
					}},
				},
			},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active contracts using wildcard filter: %w", err)
	}
	defer stream.CloseSend()
	for {
		activeContract, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to receive active contracts: %w", err)
		}
		if c, ok := activeContract.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract); ok {
			if c.ActiveContract.GetCreatedEvent().GetContractId() == contractID {
				return &ledgerv2.DisclosedContract{
					TemplateId:       c.ActiveContract.GetCreatedEvent().GetTemplateId(),
					ContractId:       c.ActiveContract.GetCreatedEvent().GetContractId(),
					CreatedEventBlob: c.ActiveContract.GetCreatedEvent().GetCreatedEventBlob(),
					SynchronizerId:   c.ActiveContract.GetSynchronizerId(),
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to find active contract with id %s", contractID)
}

func contractWithStateToDisclosedContract(contract scanProxy.ContractWithState) (*ledgerv2.DisclosedContract, error) {
	id, err := templateIDFromString(contract.Contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.Contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}
	return &ledgerv2.DisclosedContract{
		TemplateId:       id,
		ContractId:       contract.Contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   *contract.DomainId,
	}, nil
}

func activeContractIDSet(ctx context.Context, participant cldfcanton.Participant) (map[string]struct{}, error) {
	offset, err := participant.LedgerServices.State.GetLedgerEnd(ctx, &ledgerv2.GetLedgerEndRequest{})
	if err != nil {
		return nil, fmt.Errorf("get ledger end for active contract scan: %w", err)
	}
	// FiltersForAnyParty matches contracts visible to any party hosted by this participant
	// (see party-ceremony/internal/client/grpc_ledger_client.go). FiltersByParty(sender) alone
	// omits many CCIP / Splice admin contracts the sender still uses via disclosure, producing false
	// "already inactive" warnings.
	stream, err := participant.LedgerServices.State.GetActiveContracts(ctx, &ledgerv2.GetActiveContractsRequest{
		ActiveAtOffset: offset.GetOffset(),
		EventFormat: &ledgerv2.EventFormat{
			FiltersForAnyParty: &ledgerv2.Filters{
				Cumulative: []*ledgerv2.CumulativeFilter{{
					IdentifierFilter: &ledgerv2.CumulativeFilter_WildcardFilter{
						WildcardFilter: &ledgerv2.WildcardFilter{},
					},
				}},
			},
			Verbose: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get active contracts for active contract scan: %w", err)
	}
	defer stream.CloseSend()

	active := make(map[string]struct{})
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("receive active contracts for active contract scan: %w", err)
		}
		entry, ok := resp.GetContractEntry().(*ledgerv2.GetActiveContractsResponse_ActiveContract)
		if !ok || entry.ActiveContract == nil || entry.ActiveContract.GetCreatedEvent() == nil {
			continue
		}
		active[entry.ActiveContract.GetCreatedEvent().GetContractId()] = struct{}{}
	}
	return active, nil
}

func logInactiveKeyContract(logger zerolog.Logger, active map[string]struct{}, label string, contractID string) {
	if contractID == "" {
		return
	}
	if _, ok := active[contractID]; ok {
		return
	}
	logger.Warn().
		Str("label", label).
		Str("contractID", contractID).
		Msg("direct contract input already inactive before send submit")
}

func identifierString(id *ledgerv2.Identifier) string {
	if id == nil {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s", id.GetPackageId(), id.GetModuleName(), id.GetEntityName())
}

func dedupDisclosedContracts(in []*ledgerv2.DisclosedContract) []*ledgerv2.DisclosedContract {
	seen := make(map[string]struct{}, len(in))
	out := make([]*ledgerv2.DisclosedContract, 0, len(in))
	for _, dc := range in {
		if dc == nil || dc.GetContractId() == "" {
			continue
		}
		if _, ok := seen[dc.GetContractId()]; ok {
			continue
		}
		seen[dc.GetContractId()] = struct{}{}
		out = append(out, dc)
	}
	return out
}

func findDisclosedContract(disclosedContracts []*ledgerv2.DisclosedContract, contractID string) *ledgerv2.DisclosedContract {
	for _, dc := range disclosedContracts {
		if dc != nil && dc.GetContractId() == contractID {
			return dc
		}
	}
	return nil
}

func decodeEVMAddress(addr string) ([]byte, error) {
	trimmed := strings.TrimPrefix(addr, "0x")
	if len(trimmed) != 40 {
		return nil, fmt.Errorf("expected 20-byte EVM address, got %q", addr)
	}
	return hex.DecodeString(trimmed)
}

func convertToDisclosedContract(active *ledgerv2.ActiveContract) *ledgerv2.DisclosedContract {
	if active == nil || active.GetCreatedEvent() == nil {
		return nil
	}
	created := active.GetCreatedEvent()
	return &ledgerv2.DisclosedContract{
		TemplateId:       created.GetTemplateId(),
		ContractId:       created.GetContractId(),
		CreatedEventBlob: created.GetCreatedEventBlob(),
		SynchronizerId:   active.GetSynchronizerId(),
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
