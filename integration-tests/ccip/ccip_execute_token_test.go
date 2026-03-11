package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/integration-tests/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

// encodeInstrumentId encodes an InstrumentId to bytes matching Daml encodeInstrumentId.
// Format: UTF-8 bytes of "id@admin" (matches Daml's toHex(id <> "@" <> partyToText admin)).
func encodeInstrumentId(admin, identifier string) []byte {
	return []byte(identifier + "@" + admin)
}

// TestLnRTokenPool_FullReceiveFlow tests the complete CCIP inbound token release flow.
// - Deploy all CCIP contracts including CommitteeVerifier, GlobalConfig, OffRamp, PerPartyRouter
// - Mint tokens to pool (simulating prior locked tokens)
// - Build inbound message with TokenTransfer
// - Generate signatures and call CommitteeVerifier_VerifyMessage to append CCV verification
// - Call PerPartyRouter.Execute to process message and get TokenReceiveTicket
// - Call TokenPool_ReleaseFromTicket to transfer tokens from pool to receiver
// - Verify receiver received the tokens
func TestLnRTokenPool_FullReceiveFlow(t *testing.T) {
	t.Parallel()

	env := testhelpers.NewTestEnvironment(t, testhelpers.WithNumberOfParticipants(3))

	// Setup participants
	ccipParticipant := env.Chain.Participants[0]
	receiverParticipant := env.Chain.Participants[1]
	tokenPoolOwnerParticipant := env.Chain.Participants[2]

	// DAR Uploading
	// Read DARs
	rmnDar, err := contracts.GetDar(contracts.CCIPRMN, contracts.CurrentVersion)
	require.NoError(t, err)
	commonDar, err := contracts.GetDar(contracts.CCIPCommon, contracts.CurrentVersion)
	require.NoError(t, err)
	offRampDar, err := contracts.GetDar(contracts.CCIPOffRamp, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenAdminRegistryDar, err := contracts.GetDar(contracts.CCIPTokenAdminRegistry, contracts.CurrentVersion)
	require.NoError(t, err)
	committeeVerifierDar, err := contracts.GetDar(contracts.CCIPCommitteeVerifier, contracts.CurrentVersion)
	require.NoError(t, err)
	tokenPoolDar, err := contracts.GetDar(contracts.CCIPLockReleaseTokenPool, contracts.CurrentVersion)
	require.NoError(t, err)
	perPartyRouterDar, err := contracts.GetDar(contracts.CCIPPerPartyRouter, contracts.CurrentVersion)
	require.NoError(t, err)
	ccipReceiverDar, err := contracts.GetDar(contracts.CCIPReceiver, contracts.CurrentVersion)
	require.NoError(t, err)

	// Upload DARs to all participants
	dars := [][]byte{rmnDar, commonDar, offRampDar, tokenAdminRegistryDar, committeeVerifierDar, tokenPoolDar, perPartyRouterDar, ccipReceiverDar}
	packageIds, err := testhelpers.UploadDARstoMultipleParticipants(t.Context(), dars, ccipParticipant, receiverParticipant, tokenPoolOwnerParticipant)
	require.NoError(t, err)
	t.Logf("Uploaded DARs to all participants: %v", packageIds)

	// Allocate parties
	partyCCIP := ccipParticipant.PartyID
	partyReceiver := receiverParticipant.PartyID
	partyTokenPoolOwner := tokenPoolOwnerParticipant.PartyID
	t.Logf("Parties: CCIP=%s, Receiver=%s, PoolOwner=%s", partyCCIP, partyReceiver, partyTokenPoolOwner)

	// Create Scan and Registry API clients
	// Using the scanProxy endpoint of the 0-th participant, all participants are able to forward requests using the BFT Scan Proxy, it doesn't matter which one we use
	tokenSource := ccipParticipant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}
	scanProxyClient, err := scanProxy.NewClientWithResponses(ccipParticipant.Endpoints.ValidatorAPIURL, scanProxy.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create scan proxy client")
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", ccipParticipant.Endpoints.ValidatorAPIURL), tokenMetadataV1.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create token metadata client")
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(fmt.Sprintf("%s/v0/scan-proxy", ccipParticipant.Endpoints.ValidatorAPIURL), transferInstructionV1.WithRequestEditorFn(interceptor))
	require.NoError(t, err, "Failed to create transfer instruction client")

	// Get DSO Admin Party (registry admin)
	registryAdmin, err := testhelpers.GetRegistryAdmin(t.Context(), tokenMetadataClient)
	require.NoError(t, err)

	// Token Setup
	// Mint tokens to Token Pool Owner (these will be "locked" in the pool)
	poolHoldingCid, err := testhelpers.MintAMT(t.Context(), tokenPoolOwnerParticipant, tokenMetadataClient, transferInstructionClient, scanProxyClient, partyTokenPoolOwner, "100.00")
	require.NoError(t, err)
	t.Logf("Minted 100 AMT to Pool Owner, Holding CID: %s", poolHoldingCid)

	// Instrument ID for AMT
	instrumentIdAmt := &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
		{Label: "admin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: registryAdmin}}},
		{Label: "id", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "Amulet"}}},
	}}}}

	// CCIP Deployment
	sourceChainSelector := "123"
	destChainSelector := "456"

	// CCV Setup
	// Generate signer keys for CommitteeVerifier
	ccvSignerKeys := make([]*ecdsa.PrivateKey, 0, 3)
	ccvSignerPubKeys := make([]string, 0, 3)
	for range 3 {
		pk, err := crypto.GenerateKey()
		require.NoError(t, err)
		ccvSignerKeys = append(ccvSignerKeys, pk)
		pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pk.PublicKey))
		ccvSignerPubKeys = append(ccvSignerPubKeys, pubKeyHex)
	}
	t.Logf("Generated %d CCV signer keys", len(ccvSignerKeys))

	// Deploy RMNRemote (required by CommitteeVerifier and PrepareExecute/Execute)
	res, err := ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-rmn-receive"}}},
						{Label: "rmnOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "customObservers", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "cursedSubjects", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	rmnRemoteCid := extractCreatedContractId(res)
	t.Logf("Deployed RMNRemote: %s", rmnRemoteCid)

	// Deploy CommitteeVerifier
	versionTag := "49ff34ed"
	ccvId := "test-ccv-receive@" + partyCCIP
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccv-receive"}}},
						{Label: "versionTag", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: versionTag}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "messageSentObserver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "storageLocations", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Text{Text: "ipfs://test-receive"}},
						}}}}},
						{Label: "storageLocationsAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "pendingStorageLocationsAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "signerConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{{
							Key: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
							Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "sourceChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}}},
								{Label: "threshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 2}}},
								{Label: "signerKeys", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
									{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[0]}},
									{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[1]}},
									{Sum: &apiv2.Value_Text{Text: ccvSignerPubKeys[2]}},
								}}}}},
							}}}},
						}}}}}},
						{Label: "deps", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "rmnRemote", Value: rawInstanceAddress("test-rmn-receive@" + partyCCIP)},
						}}}}},
						{Label: "remoteChainFeeConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	ccvCid := extractCreatedContractId(res)
	t.Logf("Deployed CommitteeVerifier: %s (ccvId: %s)", ccvCid, ccvId)

	// Deploy GlobalConfig with source chain config including the CCV
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-globalconfig-receive"}}},
						{Label: "chainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: destChainSelector}}},
						{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: ""}}},
						{Label: "destChainConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "sourceChainConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{
							{
								Key: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
								Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: true}}},
									{Label: "onRampAddress", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "0000000000000000000000000000000000000001"}}},
									{Label: "laneMandatedCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
										rawInstanceAddress(ccvId),
									}}}}},
									{Label: "defaultCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								}}}},
							},
						}}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	globalConfigCid := extractCreatedContractId(res)
	t.Logf("Deployed GlobalConfig: %s", globalConfigCid)

	// Deploy TokenAdminRegistry
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-tar-receive"}}},
						{Label: "tokenConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid := extractCreatedContractId(res)
	t.Logf("Deployed TokenAdminRegistry: %s", tokenAdminRegistryCid)

	// Deploy OffRamp with instance addresses for GlobalConfig, RMNRemote, TokenAdminRegistry
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-offramp-receive"}}},
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "deps", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "globalConfig", Value: rawInstanceAddress("test-globalconfig-receive@" + partyCCIP)},
							{Label: "rmnRemote", Value: rawInstanceAddress("test-rmn-receive@" + partyCCIP)},
							{Label: "tokenAdminRegistry", Value: rawInstanceAddress("test-tar-receive@" + partyCCIP)},
						}}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	offRampCid := extractCreatedContractId(res)
	t.Logf("Deployed OffRamp: %s", offRampCid)

	// Deploy PerPartyRouterFactory
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-factory-receive"}}},
						{Label: "deps", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "onRamp", Value: rawInstanceAddress("placeholder-onramp@" + partyCCIP)},
							{Label: "offRamp", Value: rawInstanceAddress("test-offramp-receive@" + partyCCIP)},
							{Label: "globalConfig", Value: rawInstanceAddress("test-globalconfig-receive@" + partyCCIP)},
							{Label: "tokenAdminRegistry", Value: rawInstanceAddress("test-tar-receive@" + partyCCIP)},
							{Label: "feeQuoter", Value: rawInstanceAddress("placeholder-feequoter@" + partyCCIP)},
							{Label: "rmnRemote", Value: rawInstanceAddress("test-rmn-receive@" + partyCCIP)},
						}}}}},
						{Label: "registeredRouters", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	factoryCid := extractCreatedContractId(res)
	t.Logf("Deployed PerPartyRouterFactory: %s", factoryCid)

	// Create PerPartyRouter for receiver
	disclosedFactory, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory",
	})
	require.NoError(t, err)

	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: disclosedFactory.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "partyOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-router-receiver"}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedFactory},
		},
	})
	require.NoError(t, err)
	routerCid := ""
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			if e.Created.GetTemplateId().GetEntityName() == "PerPartyRouter" {
				routerCid = e.Created.ContractId
				break
			}
		}
	}
	require.NotEmpty(t, routerCid)
	t.Logf("Created PerPartyRouter for receiver: %s", routerCid)

	// Token Pool Setup
	// Deploy inbound RateLimiter required by ReleaseFromTicket receive flow.
	res, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive-inbound-rl"}}},
						{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "remoteChainSelector", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}}},
						{Label: "direction", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitDirection_Inbound"}}}},
						{Label: "mode", Value: &apiv2.Value{Sum: &apiv2.Value_Enum{Enum: &apiv2.Enum{Constructor: "RateLimitMode_DefaultFinality"}}}},
						{Label: "isEnabled", Value: &apiv2.Value{Sum: &apiv2.Value_Bool{Bool: false}}},
						{Label: "capacity", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "rate", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "tokens", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: "0"}}},
						{Label: "lastUpdated", Value: &apiv2.Value{Sum: &apiv2.Value_Timestamp{Timestamp: time.Now().UnixMicro()}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	inboundRateLimiterCid := extractCreatedContractId(res)
	t.Logf("Deployed inbound RateLimiter: %s", inboundRateLimiterCid)

	// Deploy LockReleaseTokenPool
	res, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "ccipOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
						{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "decimals", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 6}}},
						{Label: "chainPoolConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{{
							Key: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
							Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "inboundCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								{Label: "outboundCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
								{Label: "remotePools", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
									{Sum: &apiv2.Value_Text{Text: hex.EncodeToString([]byte(partyTokenPoolOwner))}},
								}}}}},
							}}}},
						}}}}}},
						{Label: "chainFeeConfigs", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "outboundRateLimiters", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "inboundRateLimiters", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: []*apiv2.GenMap_Entry{{
							Key:   &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: sourceChainSelector}},
							Value: rawInstanceAddress("test-pool-receive-inbound-rl@" + partyTokenPoolOwner),
						}}}}}},
						{Label: "remoteTokens", Value: &apiv2.Value{Sum: &apiv2.Value_GenMap{GenMap: &apiv2.GenMap{Entries: nil}}}},
						{Label: "poolReceiveContext", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "values", Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}}},
						}}}}},
						// Use RelativeHours 24 for Amulet tokens (have expiry constraints)
						{Label: "transferTimeout", Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
							Constructor: "RelativeHours",
							Value:       &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 24}},
						}}}},
					}},
				}},
			}},
			ActAs: []string{partyTokenPoolOwner},
		},
	})
	require.NoError(t, err)
	tokenPoolCid := extractCreatedContractId(res)
	t.Logf("Deployed LockReleaseTokenPool: %s", tokenPoolCid)

	// Register pool in TokenAdminRegistry (ProposeAdministrator -> AcceptAdminRole -> SetPool)
	// Step 1: ProposeAdministrator
	res, err = ccipParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: tokenAdminRegistryCid,
					Choice:     "TokenAdminRegistry_ProposeAdministrator",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "newAdmin", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyCCIP}}},
					}}}},
				}},
			}},
			ActAs: []string{partyCCIP},
		},
	})
	require.NoError(t, err)
	tokenAdminRegistryCid = extractCreatedContractId(res)
	t.Log("Called ProposeAdministrator, tokenAdminRegistryCid:", tokenAdminRegistryCid)

	// Step 2: AcceptAdminRole
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	_, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_AcceptAdminRole",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	t.Log("Called AcceptAdminRole")

	// Step 3: SetPool
	time.Sleep(500 * time.Millisecond)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)

	_, err = tokenPoolOwnerParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry"},
					ContractId: disclosedTar.ContractId,
					Choice:     "TokenAdminRegistry_SetPool",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instrumentId", Value: instrumentIdAmt},
						{Label: "tokenPool", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "poolOwner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
							{Label: "poolInstanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-pool-receive"}}},
						}}}}}}}},
						{Label: "caller", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyTokenPoolOwner}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyTokenPoolOwner},
			DisclosedContracts: []*apiv2.DisclosedContract{disclosedTar},
		},
	})
	require.NoError(t, err)
	t.Log("Called SetPool")

	// Build Message
	// Encode instrumentId for destTokenAddress
	encodedInstrumentId := encodeInstrumentId(registryAdmin, "Amulet")
	encodedInstrumentIdHex := hex.EncodeToString(encodedInstrumentId)

	// Build token transfer (5 AMT in Splice Decimal format)
	tokenAmount := big.NewInt(5)
	encodedTokenTransfer := buildTokenTransferV1(tokenAmount, partyTokenPoolOwner, encodedInstrumentIdHex, partyReceiver)

	// Build message
	msg := &MessageV1{
		SourceChainSelector: 123,
		DestChainSelector:   456,
		SequenceNumber:      1,
		ExecutionGasLimit:   200000,
		CCIPReceiveGasLimit: 100000,
		Finality:            2000,
		CCVAndExecutorHash:  [32]byte{},
		OnRampAddress:       hexToBytes("0000000000000000000000000000000000000001"),
		OffRampAddress:      hexToBytes("0000000000000000000000000000000000000002"),
		Sender:              hexToBytes("0000000000000000000000000000000000000003"),
		Receiver:            EncodePartyID(partyReceiver),
		DestBlob:            []byte{},
		TokenTransfer:       encodedTokenTransfer,
		MessageData:         []byte{},
	}
	encodedMessage, err := EncodeMessageV1(msg)
	require.NoError(t, err)
	encodedMessageHex := hex.EncodeToString(encodedMessage)
	messageHash := crypto.Keccak256(encodedMessage)
	messageHashHex := hex.EncodeToString(messageHash)
	t.Logf("Message hash: %s", messageHashHex)

	// CCV Verification
	// Generate verifierResults with 2 of 3 signatures
	verifierResults, err := GenerateVerifierResults(encodedMessage, ccvSignerKeys[:2])
	require.NoError(t, err)
	verifierResultsHex := hex.EncodeToString(verifierResults)

	// Deploy CCIPReceiver for receiver
	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "test-ccipreceiver"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
						{Label: "requiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}},
				}},
			}},
			ActAs: []string{partyReceiver},
		},
	})
	require.NoError(t, err)
	ccipReceiverCid := extractCreatedContractId(res)
	t.Logf("Deployed CCIPReceiver: %s", ccipReceiverCid)

	// Get TransferFactory and pool holdings (needed by CCIPReceiver.Execute for token release)
	transferFactoryCid, transferFactoryDisclosures, choiceContext, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, registryAdmin, partyTokenPoolOwner, partyReceiver)
	require.NoError(t, err)

	poolHoldings, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	require.NotEmpty(t, poolHoldings, "Pool should have holdings")

	// Filter for unlocked holdings only (Amulet tokens may be locked until next mining round)
	var poolHoldingCids []*apiv2.Value
	var poolHoldingDisclosures []*apiv2.DisclosedContract
	var lockedCount, unlockedCount int
	for _, h := range poolHoldings {
		viewFields := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()
		lockField := viewFields[3].GetValue()
		isLocked := lockField.GetOptional().GetValue() != nil

		if isLocked {
			lockedCount++
			continue
		}
		unlockedCount++

		poolHoldingCids = append(poolHoldingCids, &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: h.GetCreatedEvent().GetContractId()}})
		poolHoldingDisclosures = append(poolHoldingDisclosures, &apiv2.DisclosedContract{
			ContractId:       h.GetCreatedEvent().GetContractId(),
			TemplateId:       h.GetCreatedEvent().GetTemplateId(),
			CreatedEventBlob: h.GetCreatedEvent().GetCreatedEventBlob(),
		})
	}
	t.Logf("Pool has %d holdings (%d unlocked, %d locked)", len(poolHoldings), unlockedCount, lockedCount)

	if unlockedCount == 0 {
		t.Skip("SKIPPING: All pool holdings are locked (Amulet tokens are locked until next mining round). " +
			"This is expected on fresh localnet. Either wait for mining round to complete, or use TestToken instead of Amulet.")
	}

	// Capture receiver's balance before execute
	receiverHoldingsBefore, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceBefore := getHoldingsBalance(receiverHoldingsBefore)

	// Get disclosures for CCIPReceiver.Execute
	disclosedCCIPReceiver, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver",
	})
	require.NoError(t, err)
	disclosedRouter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#ccip-perpartyrouter", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter",
	})
	require.NoError(t, err)
	disclosedOffRamp, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-offramp", ModuleName: "CCIP.OffRamp", EntityName: "OffRamp",
	})
	require.NoError(t, err)
	disclosedGlobalConfig, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.GlobalConfig", EntityName: "GlobalConfig",
	})
	require.NoError(t, err)
	disclosedTar, err = testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-tokenadminregistry", ModuleName: "CCIP.TokenAdminRegistry", EntityName: "TokenAdminRegistry",
	})
	require.NoError(t, err)
	disclosedRmnRemote, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-rmn", ModuleName: "CCIP.RMNRemote", EntityName: "RMNRemote",
	})
	require.NoError(t, err)
	disclosedCCV, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), ccipParticipant, &apiv2.Identifier{
		PackageId: "#ccip-committeeverifier", ModuleName: "CCIP.CommitteeVerifier", EntityName: "CommitteeVerifier",
	})
	require.NoError(t, err)
	disclosedPool, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#ccip-lockreleasetokenpool", ModuleName: "CCIP.LockReleaseTokenPool", EntityName: "LockReleaseTokenPool",
	})
	require.NoError(t, err)
	disclosedRateLimiter, err := testhelpers.GetDisclosedContractByTemplateId(t.Context(), tokenPoolOwnerParticipant, &apiv2.Identifier{
		PackageId: "#ccip-common", ModuleName: "CCIP.RateLimiter", EntityName: "RateLimiter",
	})
	require.NoError(t, err)

	executeDisclosures := slices.Concat(
		[]*apiv2.DisclosedContract{disclosedCCIPReceiver, disclosedRouter, disclosedOffRamp, disclosedGlobalConfig, disclosedTar, disclosedRmnRemote, disclosedCCV, disclosedPool, disclosedRateLimiter},
		transferFactoryDisclosures,
		poolHoldingDisclosures,
	)

	// Create context - replace with EDS
	executeContext := map[string]any{
		"values": map[string]any{
			"off-ramp": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosedOffRamp.ContractId,
			},
			"global-config": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosedGlobalConfig.ContractId,
			},
			"token-admin-registry": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosedTar.ContractId,
			},
			"rmn-remote": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosedRmnRemote.ContractId,
			},
		},
	}
	executeContextValue, err := testhelpers.ChoiceContextFromData(executeContext)
	require.NoError(t, err)
	poolExtraContext, err := testhelpers.ChoiceContextFromData(map[string]any{
		"values": map[string]any{
			"rate-limiter": map[string]any{
				"tag":   "AV_ContractId",
				"value": disclosedRateLimiter.ContractId,
			},
		},
	})
	require.NoError(t, err)

	// CCIPReceiver.Execute: PrepareExecute + CCV + Pool Verify + Execute + Release in one transaction
	res, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.Must(uuid.NewUUID()).String(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					ContractId: ccipReceiverCid,
					Choice:     "Execute",
					ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "context", Value: executeContextValue},
						{Label: "routerCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: routerCid}}},
						{Label: "encodedMessage", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: encodedMessageHex}}},
						{Label: "tokenTransfer", Value: &apiv2.Value{Sum: &apiv2.Value_Optional{Optional: &apiv2.Optional{Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "tokenPoolCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedPool.ContractId}}},
							{Label: "tokenReceiverParty", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: partyReceiver}}},
							{Label: "tokenInput", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "transferFactory", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: transferFactoryCid}}},
								{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "context", Value: choiceContext},
									{Label: "meta", Value: emptyMetadata},
								}}}}},
								{Label: "tokenPoolHoldings", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: poolHoldingCids}}}},
							}}}}},
							{Label: "poolExtraContext", Value: poolExtraContext},
						}}}}}}}},
						{Label: "ccvInputs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: []*apiv2.Value{
							{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "ccvCid", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: disclosedCCV.ContractId}}},
								{Label: "verifierResults", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: verifierResultsHex}}},
								{Label: "ccvExtraContext", Value: emptyCCIPContext},
							}}}},
						}}}}},
						{Label: "additionalRequiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
					}}}},
				}},
			}},
			ActAs:              []string{partyReceiver},
			DisclosedContracts: executeDisclosures,
		},
	})
	require.NoError(t, err)
	t.Log("CCIPReceiver.Execute completed")

	// Check for pending AmuletTransferInstruction (Amulet transfers without receiver preapproval)
	var pendingTransferInstructionCid string
	var releaseCompleted bool
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			templateName := e.Created.GetTemplateId().GetEntityName()
			if templateName == "AmuletTransferInstruction" {
				pendingTransferInstructionCid = e.Created.GetContractId()
			}
			if templateName == "Amulet" || templateName == "LockedAmulet" {
				releaseCompleted = true
			}
		}
	}

	// If pending, receiver must accept the TransferInstruction to complete the transfer
	if !releaseCompleted && pendingTransferInstructionCid != "" {
		time.Sleep(500 * time.Millisecond)

		acceptContextResp, err := transferInstructionClient.GetTransferInstructionAcceptContextWithResponse(t.Context(), pendingTransferInstructionCid, transferInstructionV1.GetChoiceContextRequest{})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, acceptContextResp.StatusCode(), "Failed to get accept context: %s", string(acceptContextResp.Body))

		acceptDisclosures := make([]*apiv2.DisclosedContract, 0, len(acceptContextResp.JSON200.DisclosedContracts))
		for _, contract := range acceptContextResp.JSON200.DisclosedContracts {
			id, err := testhelpers.TemplateIdFromString(contract.TemplateId)
			require.NoError(t, err)
			createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
			require.NoError(t, err)
			acceptDisclosures = append(acceptDisclosures, &apiv2.DisclosedContract{
				TemplateId:       id,
				ContractId:       contract.ContractId,
				CreatedEventBlob: createdEventBlob,
			})
		}
		acceptContext, err := testhelpers.ChoiceContextFromData(acceptContextResp.JSON200.ChoiceContextData)
		require.NoError(t, err)

		_, err = receiverParticipant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.Must(uuid.NewUUID()).String(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
						ContractId: pendingTransferInstructionCid,
						Choice:     "TransferInstruction_Accept",
						ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
							{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "context", Value: acceptContext},
								{Label: "meta", Value: emptyMetadata},
							}}}}},
						}}}},
					}},
				}},
				ActAs:              []string{partyReceiver},
				DisclosedContracts: acceptDisclosures,
			},
		})
		require.NoError(t, err)
	}

	// Verify receiver's balance increased by the expected transfer amount
	receiverHoldingsAfter, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), receiverParticipant, &apiv2.Identifier{
		PackageId: "#splice-api-token-holding-v1", ModuleName: "Splice.Api.Token.HoldingV1", EntityName: "Holding",
	})
	require.NoError(t, err)
	receiverBalanceAfter := getHoldingsBalance(receiverHoldingsAfter)

	expectedTransferAmount, _ := new(big.Float).SetInt(tokenAmount).Float64()
	actualTransferAmount := receiverBalanceAfter - receiverBalanceBefore
	require.InDelta(t, expectedTransferAmount, actualTransferAmount, 0.01, "Receiver balance should increase by transfer amount")
	t.Logf("Receiver balance: %.2f -> %.2f AMT (transferred %.2f)", receiverBalanceBefore, receiverBalanceAfter, actualTransferAmount)
}

// getHoldingsBalance returns total balance across all holdings
func getHoldingsBalance(holdings []*apiv2.ActiveContract) float64 {
	var total float64
	for _, h := range holdings {
		balanceStr := h.GetCreatedEvent().GetInterfaceViews()[0].GetViewValue().GetFields()[2].GetValue().GetSum().(*apiv2.Value_Numeric).Numeric
		balance, _ := new(big.Float).SetString(balanceStr)
		balanceFloat, _ := balance.Float64()
		total += balanceFloat
	}

	return total
}

// extractCreatedContractId returns the first created contract ID from a transaction response
func extractCreatedContractId(res *apiv2.SubmitAndWaitForTransactionResponse) string {
	for _, event := range res.GetTransaction().GetEvents() {
		if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
			return e.Created.ContractId
		}
	}

	return ""
}

// buildTokenTransferV1 builds the encoded token transfer bytes for the message.
// tokenReceiver is keccak256-hashed to match Daml encodePartyAddress.
func buildTokenTransferV1(amount *big.Int, sourcePoolOwner, destTokenAddressHex, tokenReceiverParty string) *TokenTransferV1 {
	sourcePoolAddress := []byte(sourcePoolOwner)
	sourceTokenAddress := []byte{} // Empty for Canton

	destTokenAddress, _ := hex.DecodeString(destTokenAddressHex)

	return &TokenTransferV1{
		Amount:             amount,
		SourcePoolAddress:  sourcePoolAddress,
		SourceTokenAddress: sourceTokenAddress,
		DestTokenAddress:   destTokenAddress,
		TokenReceiver:      EncodePartyID(tokenReceiverParty),
		ExtraData:          []byte{},
	}
}

var emptyMetadata = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}

var emptyCCIPContext = &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
	{
		Label: "values",
		Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
	},
}}}}
