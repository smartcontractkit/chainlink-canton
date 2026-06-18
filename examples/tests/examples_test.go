package tests

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	routerwrapper "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v2_0_0/onramp"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/evm"
	v1 "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/api/handlers/v1"
	indexerclient "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/client"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/ccipruntime"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/receiver"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/sender"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
	"github.com/smartcontractkit/chainlink-canton/testhelpers/eds"
)

var (
	// EDS
	edsURL     = ""
	indexerURL = ""

	// Canton - Local Chain
	cantonSelector              = chainsel.CANTON_TESTNET.Selector
	authorizationServerURL      = ""
	authClientID                = ""
	participantGrpcLedgerApiURL = ""
	validatorApiURL             = ""
	userID                      = ""
	partyID                     = ""

	// Eth - Remote Chain
	ethSelector      = chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	ethRpcURL        = ""
	ethPrivateKeyHex = ""
	ethRouterAddress = common.HexToAddress("0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59")
	ethTokenAddress  = common.HexToAddress("0xeEe6675b20fE5950eb51361b93021D076289F612")
	noExecutionTag   = common.HexToAddress("0xEBa517d200000000000000000000000000000000")

	// Example Receiver contract that emits an event if a message is received
	ccipReceiverContract = common.HexToAddress("0x1E340B34Cc732a71f5e59804da7E645a52e10E1B")
	receiverAddress      = common.HexToAddress("0x0")
)

const (
	dsoPartyID       = "DSO::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337"
	ccipOwnerPartyID = "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"
)

var amuletInstrumentID = &splice_api_token_holding_v1.InstrumentId{
	Admin: dsoPartyID,
	Id:    "Amulet",
}

var linkInstrumentId = &splice_api_token_holding_v1.InstrumentId{
	Admin: ccipOwnerPartyID,
	Id:    "link-token",
}

func TestMulti(t *testing.T) {
	authProvider, err := authorizationcode.NewDiscoveryProvider(t.Context(), authorizationServerURL, authClientID)
	require.NoError(t, err)
	require.NotNil(t, authProvider)

	token, err := authProvider.TokenSource().Token()
	require.NoError(t, err)
	require.NotNil(t, token)

	fmt.Println(token.AccessToken)

	rpcProviderConfig := provider.RPCChainProviderConfig{
		Participants: []provider.ParticipantConfig{
			{
				Endpoints: provider.Endpoints{
					JSONLedgerAPIURL: "json",
					GRPCLedgerAPIURL: participantGrpcLedgerApiURL,
					ValidatorAPIURL:  validatorApiURL,
				},
				UserID:       userID,
				PartyID:      partyID,
				AuthProvider: authProvider,
			},
		},
	}

	chain, err := provider.NewRPCChainProvider(chainsel.CANTON_LOCALNET.Selector, rpcProviderConfig).Initialize(t.Context())
	require.NoError(t, err)
	require.NotNil(t, chain)

	cantonChain := chain.(*canton.Chain)
	participant := cantonChain.Participants[0]

	_, _, transferInstructionClient, err := testhelpers.NewValidatorAPIClients(participant)
	require.NoError(t, err)

	t.Run("GetVersion", func(t *testing.T) {
		versionResp, err := participant.LedgerServices.Version.GetLedgerApiVersion(t.Context(), &apiv2.GetLedgerApiVersionRequest{})
		require.NoError(t, err)

		fmt.Println(versionResp.GetVersion())
	})

	t.Run("AcceptIncomingTransferInstruction", func(t *testing.T) {
		transferInstructionCid := ""

		contextResp, err := transferInstructionClient.GetTransferInstructionAcceptContextWithResponse(t.Context(), transferInstructionCid, oapiTransferInstruction.GetTransferInstructionAcceptContextJSONRequestBody{
			Meta: nil,
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, contextResp.StatusCode(), "Unexpected status code, response: %v", string(contextResp.Body))

		var disclosedContracts []*apiv2.DisclosedContract
		for _, contract := range contextResp.JSON200.DisclosedContracts {
			id, err := testhelpers.TemplateIdFromString(contract.TemplateId)
			require.NoError(t, err)
			createdEventBlob, err := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
			require.NoError(t, err)
			disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
				TemplateId:       id,
				ContractId:       contract.ContractId,
				CreatedEventBlob: createdEventBlob,
				SynchronizerId:   contract.SynchronizerId,
			})
		}

		acceptContext, err := testhelpers.ChoiceContextFromData(contextResp.JSON200.ChoiceContextData)
		require.NoError(t, err)
		fmt.Println(acceptContext)

		resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.NewString(),
				Commands: []*apiv2.Command{
					{
						Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
							TemplateId: &apiv2.Identifier{PackageId: "#splice-api-token-transfer-instruction-v1", ModuleName: "Splice.Api.Token.TransferInstructionV1", EntityName: "TransferInstruction"},
							ContractId: transferInstructionCid,
							Choice:     "TransferInstruction_Accept",
							ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
								{Label: "extraArgs", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "context", Value: acceptContext},
									{Label: "meta", Value: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{{
										Label: "values",
										Value: &apiv2.Value{Sum: &apiv2.Value_TextMap{TextMap: &apiv2.TextMap{Entries: nil}}},
									}}}}}},
								}}}}},
							}}}},
						}},
					},
				},
				ActAs:              []string{participant.PartyID},
				DisclosedContracts: disclosedContracts,
			},
		})
		require.NoError(t, err)

		fmt.Println("Accepted in update: ", resp.GetTransaction().GetUpdateId())
	})

	t.Run("FilterContracts", func(t *testing.T) {
		activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-core", ModuleName: "CCIP.Events", EntityName: "CCIPMessageSent"})
		require.NoError(t, err)
		for i, contract := range activeContracts {
			fmt.Printf("Active contract %d: ID=%s, CreateArguments=%s\n", i, contract.GetCreatedEvent().GetContractId(), contract.GetCreatedEvent().CreateArguments)
		}
	})

	t.Run("Enumerate Holdings", func(t *testing.T) {
		holdings, err := testhelpers.ListActiveContractsByInterfaceId(t.Context(), participant, &apiv2.Identifier{
			PackageId:  "#splice-api-token-holding-v1",
			ModuleName: "Splice.Api.Token.HoldingV1",
			EntityName: "Holding",
		})
		require.NoError(t, err)
		for _, holding := range holdings {
			for _, view := range holding.GetCreatedEvent().GetInterfaceViews() {
				var holdingView splice_api_token_holding_v1.HoldingView
				err = ledger.RecordToStruct(view.GetViewValue(), &holdingView)
				require.NoError(t, err)
				fmt.Printf("InstrumentId: %v@%v Owner: %v Amount: %v \n", holdingView.InstrumentId.Id, holdingView.InstrumentId.Admin, holdingView.Owner, holdingView.Amount)
			}
		}
	})

	t.Run("Create PerPartyRouter", func(t *testing.T) {
		ccipEdsClient, err := oapiCCIP.NewClientWithResponses(edsURL)
		require.NoError(t, err)

		routerCid := getRouter(t, participant, ccipEdsClient)
		fmt.Println("PerPartyRouter Contract ID: ", routerCid)
	})

	t.Run("Create CCIPReceiver", func(t *testing.T) {
		receiverCid := getReceiver(t, participant)
		fmt.Println("CCIPReceiver Contract ID: ", receiverCid)
	})

	t.Run("Create CCIPSender", func(t *testing.T) {
		senderCid := getSender(t, participant)
		fmt.Println("CCIPSender Contract ID: ", senderCid)
	})

	t.Run("Send: EVM -> Canton (Token Transfer, Native)", func(t *testing.T) {
		// Source
		finality := protocol.FinalityWaitForFinality
		execGasLimit := uint32(0)
		executorAddress := noExecutionTag

		// Dest
		receiverParty := participant.PartyID

		// Message
		tokenAmount := big.NewInt(1e18)

		chainIdString, err := chainsel.GetChainIDFromSelector(ethSelector)
		require.NoError(t, err)
		sourceChainId, err := strconv.ParseUint(chainIdString, 10, 64)
		require.NoError(t, err)
		client, err := ethclient.DialContext(t.Context(), ethRpcURL)
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(ethPrivateKeyHex, "0x"))
		require.NoError(t, err)
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(sourceChainId))
		require.NoError(t, err)

		router, err := routerwrapper.NewRouter(ethRouterAddress, client)
		require.NoError(t, err)

		receiverPartyHashed := contracts.HashedPartyFromString(receiverParty)

		extraArgs, err := evm.NewV3ExtraArgs(finality, execGasLimit, executorAddress.Hex(), nil, nil, nil, nil)
		require.NoError(t, err)

		msg := routerwrapper.ClientEVM2AnyMessage{
			Receiver: receiverPartyHashed.Bytes(),
			Data:     nil,
			TokenAmounts: []routerwrapper.ClientEVMTokenAmount{{
				Token:  ethTokenAddress,
				Amount: tokenAmount,
			}},
			FeeToken:  common.Address{}, // Native
			ExtraArgs: extraArgs,
		}

		fee, err := router.GetFee(&bind.CallOpts{Context: t.Context()}, cantonSelector, msg)
		require.NoError(t, err)
		fmt.Printf("Fee for sending message: %s\n", fee.String())
		auth.Value = fee

		tx, err := router.CcipSend(auth, cantonSelector, msg)
		require.NoError(t, err)
		fmt.Printf("Sent message with tx hash: %s\n", tx.Hash().Hex())

		receipt, err := bind.WaitMined(t.Context(), client, tx)
		require.NoError(t, err)
		fmt.Printf("Transaction mined in block: %d\n", receipt.BlockNumber.Uint64())

		var sentEvent *onramp.OnRampCCIPMessageSent
		eventTopic := (onramp.OnRampCCIPMessageSent{}).Topic()
		for _, lg := range receipt.Logs {
			if len(lg.Topics) == 0 || lg.Topics[0] != eventTopic {
				continue
			}

			onRamp, err := onramp.NewOnRamp(lg.Address, client)
			require.NoError(t, err)

			messageSentEvent, err := onRamp.ParseCCIPMessageSent(*lg)
			require.NoError(t, err)
			sentEvent = messageSentEvent
			break
		}
		require.NotNil(t, sentEvent, "Sent event not found in transaction logs")
		fmt.Printf("Message sent with ID: %s\n", common.Bytes2Hex(sentEvent.MessageId[:]))
	})

	t.Run("Execute: Canton (Token Transfer)", func(t *testing.T) {
		messageId := common.HexToHash("3bb03d186c6adb90baab2004d683dca394804e721c9bda207d4a72ba3d86a351")
		timeout := time.Minute

		indexer, err := indexerclient.NewIndexerClient(indexerURL, &http.Client{Timeout: 15 * time.Second})
		require.NoError(t, err)
		ccipEdsClient, err := oapiCCIP.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		ccvEdsClient, err := oapiCCV.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		tokenPoolEdsClient, err := oapiTokenPool.NewClientWithResponses(edsURL)
		require.NoError(t, err)

		start := time.Now()
		var response v1.VerifierResultsByMessageIDResponse
		for {
			if time.Now().Sub(start) > timeout {
				t.Fatal("Timed out waiting for message to be processed by indexer")
			}
			status, resp, err := indexer.VerifierResultsByMessageID(t.Context(), v1.VerifierResultsByMessageIDInput{MessageID: messageId.Hex()})
			if err != nil {
				t.Logf("Error querying indexer: %v", err)
				time.Sleep(5 * time.Second)
				continue
			} else if status != http.StatusOK {
				t.Logf("Unexpected status code from indexer: %d", status)
				time.Sleep(5 * time.Second)
				continue
			}

			response = resp
			break
		}
		require.NotEmptyf(t, response.Results, "Expected at least one verifier result for message ID %s", messageId.Hex())
		t.Logf("Received response from indexer: %+v", response)

		verifierResult := response.Results[0].VerifierResult
		message := verifierResult.Message
		encodedMessage, err := message.Encode()
		require.NoError(t, err)

		targetInstrumentId := contracts.BytesToEncodedInstrumentID(message.TokenTransfer.DestTokenAddress)
		ccvAddress := contracts.BytesToInstanceAddress(verifierResult.VerifierDestAddress)

		tokenPoolAddress, err := eds.GetTokenPoolForToken(t.Context(), ccipEdsClient, targetInstrumentId)
		require.NoError(t, err)
		ccipExecuteDisclosure, err := eds.GetCCIPExecuteDisclosure(t.Context(), ccipEdsClient, hex.EncodeToString(encodedMessage))
		require.NoError(t, err)
		ccvExecuteDisclosure, err := eds.GetCCVExecuteDisclosure(t.Context(), ccvEdsClient, hex.EncodeToString(encodedMessage), ccvAddress)
		require.NoError(t, err)
		tokenPoolExecuteDisclosure, err := eds.GetTokenPoolExecuteDisclosure(t.Context(), tokenPoolEdsClient, hex.EncodeToString(encodedMessage), tokenPoolAddress.InstanceAddress())
		require.NoError(t, err)

		fmt.Println("CCIP Execute Disclosure:", ccipExecuteDisclosure)
		fmt.Println("CCV Execute Disclosure:", ccvExecuteDisclosure)
		fmt.Println("Token Pool Execute Disclosure:", tokenPoolExecuteDisclosure)

		// Get PerPartyRouter
		routerCid := getRouter(t, participant, ccipEdsClient)

		// Get CCIPReceiver
		receiverCid := getReceiver(t, participant)

		fmt.Println("PerPartyRouter CID:", routerCid)
		fmt.Println("CCIPReceiver CID:", receiverCid)

		executeArgs := receiver.Execute{
			Context:        ccipExecuteDisclosure.ChoiceContext,
			RouterCid:      types.CONTRACT_ID(routerCid),
			EncodedMessage: types.TEXT(hex.EncodeToString(encodedMessage)),
			TokenTransfer: &receiver.TokenTransferInput{
				TokenPoolCid:       types.CONTRACT_ID(tokenPoolExecuteDisclosure.ContractId),
				TokenReceiverParty: types.PARTY(participant.PartyID),
				PoolExtraContext:   tokenPoolExecuteDisclosure.ChoiceContext,
			},
			CcvInputs: []receiver.CCVInput{
				{
					CcvCid:          types.CONTRACT_ID(ccvExecuteDisclosure.ContractId),
					VerifierResults: types.TEXT(hex.EncodeToString(verifierResult.CCVData)),
					CcvExtraContext: ccvExecuteDisclosure.ChoiceContext,
				},
			},
		}

		resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.NewString(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId:     &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
						ContractId:     receiverCid,
						Choice:         "Execute",
						ChoiceArgument: ledger.MapToValue(executeArgs),
					}},
				}},
				ActAs: []string{participant.PartyID},
				DisclosedContracts: slices.Concat(
					tokenPoolExecuteDisclosure.DisclosedContracts,
					ccipExecuteDisclosure.DisclosedContracts,
					ccvExecuteDisclosure.DisclosedContracts,
				),
			},
		})
		require.NoError(t, err)
		fmt.Println("Message executed in Update: ", resp.GetTransaction().GetUpdateId())
	})

	t.Run("Send: Canton->EVM (Message Only, Native)", func(t *testing.T) {
		messageReceiver := ccipReceiverContract

		ccipEdsClient, err := oapiCCIP.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		ccvEdsClient, err := oapiCCV.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		executorEdsClient, err := oapiExecutor.NewClientWithResponses(edsURL)
		require.NoError(t, err)

		// Get Fee Input holding
		feeTokenHoldings, err := testhelpers.ListHoldingsForInstrument(t.Context(), participant, amuletInstrumentID)
		require.NoError(t, err)

		transferFactoryCid, transferFactoryDisclosures, choiceContextRaw, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, string(amuletInstrumentID.Admin), participant.PartyID, ccipOwnerPartyID)
		require.NoError(t, err)
		choiceContext, err := contracts.ChoiceContextFromData(choiceContextRaw)
		require.NoError(t, err)

		// Get PerPartyRouter
		routerCid := getRouter(t, participant, ccipEdsClient)

		// Get CCIPSender
		senderCid := getSender(t, participant)

		msg := oapiCommon.Message{
			DestinationChainSelector: strconv.FormatUint(ethSelector, 10),
			Executor: struct {
				Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
				Type    oapiCommon.MessageExecutorType `json:"type"`
			}{
				Type: oapiCommon.Empty,
			},
			FeeToken: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(amuletInstrumentID.Admin),
				Id:    string(amuletInstrumentID.Id),
			},
			GasLimit:      50_000,
			Payload:       hex.EncodeToString([]byte("Hello, EVM from Canton!")),
			Receiver:      hex.EncodeToString(messageReceiver.Bytes()),
			TokenTransfer: nil,
		}

		ccipSendDisclosure, err := eds.GetCCIPSendDisclosure(t.Context(), ccipEdsClient, msg, nil, nil)
		require.NoError(t, err)
		// EDS returns the default CCV(s) if no input is specified
		defaultCCVAddress, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
		require.NoError(t, err)
		// EDS returns the default Executor (if set) if no specific executor input is specified
		defaultExecutorAddress, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
		require.NoError(t, err)
		ccvSendDisclosure, err := eds.GetCCVSendDisclosure(t.Context(), ccvEdsClient, msg, defaultCCVAddress.InstanceAddress())
		require.NoError(t, err)
		executorSendDisclosure, err := eds.GetExecutorSendDisclosure(t.Context(), executorEdsClient, msg, defaultExecutorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
		require.NoError(t, err)

		feeTokenInputCids := make([]types.CONTRACT_ID, len(feeTokenHoldings))
		for i, holding := range feeTokenHoldings {
			feeTokenInputCids[i] = types.CONTRACT_ID(holding.ContractID)
		}

		sendArgs := sender.Send{
			DestinationChainSelector: types.NUMERIC(msg.DestinationChainSelector),
			Message: core.Canton2AnyMessage{
				Receiver:      types.TEXT(msg.Receiver),
				Payload:       types.TEXT(msg.Payload),
				TokenTransfer: nil,
				FeeToken:      *amuletInstrumentID,
				ExtraArgs: core.ExtraArgs{
					V3: &core.GenericExtraArgsV3{
						GasLimit: types.INT64(msg.GasLimit),
						Ccvs:     nil,
						Executor: core.ExecutorExtraArg{
							ExecutorUseDefault: &core.ExecutorUseDefault{ExecutorArgs: ""},
						},
						TokenReceiver: "",
						TokenArgs:     "",
					},
				},
			},
			Context:   ccipSendDisclosure.ChoiceContext,
			RouterCid: types.CONTRACT_ID(routerCid),
			FeeTokenInput: sender.FeeTokenInput{
				SenderInputCids:         feeTokenInputCids,
				FeeTokenConfigCid:       types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid),
				FeeTokenTransferFactory: types.CONTRACT_ID(transferFactoryCid),
				FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
					Context: choiceContext,
					Meta: splice_api_token_metadata_v1.Metadata{
						Values: map[string]types.TEXT{},
					},
				},
			},
			CcvSendInputs: []sender.CCVSendInput{
				{
					CcvAddress:      ccvSendDisclosure.Address.Binding(),
					CcvCid:          types.CONTRACT_ID(ccvSendDisclosure.ContractId),
					CcvExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
				},
			},
			TokenTransferInput: nil,
			ExecutorInput: &sender.ExecutorInput{
				ExecutorCid:          types.CONTRACT_ID(executorSendDisclosure.ContractId),
				ExecutorExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
			},
		}
		allDisclosures := slices.Concat(
			transferFactoryDisclosures,
			ccipSendDisclosure.DisclosedContracts,
			ccvSendDisclosure.DisclosedContracts,
			executorSendDisclosure.DisclosedContracts,
		)
		resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.NewString(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
						ContractId:     senderCid,
						Choice:         "Send",
						ChoiceArgument: ledger.MapToValue(sendArgs),
					}},
				}},
				ActAs:              []string{participant.PartyID},
				DisclosedContracts: allDisclosures,
			},
		})
		require.NoError(t, err)
		fmt.Println("Message sent in Update: ", resp.GetTransaction().GetUpdateId())

		var returnedMessageId string
		for _, event := range resp.GetTransaction().GetEvents() {
			if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
				if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
					ccipMessageSent, err := bindings.UnmarshalCreatedEvent[core.CCIPMessageSent](e.Created)
					require.NoError(t, err)
					returnedMessageId = string(ccipMessageSent.Event.MessageId)
					break
				}
			}
		}
		require.NotEmpty(t, returnedMessageId)
		fmt.Println("Message sent with MessageID: ", returnedMessageId)
	})

	t.Run("Send: Canton->EVM (Token Only, Native)", func(t *testing.T) {
		// The address that will receive the tokens
		messageReceiver := common.HexToAddress("0x90392A1E8A941098a3C75E0BDB172cFdE7E4f1f4")

		ccipEdsClient, err := oapiCCIP.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		ccvEdsClient, err := oapiCCV.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		executorEdsClient, err := oapiExecutor.NewClientWithResponses(edsURL)
		require.NoError(t, err)
		tokenPoolEdsClient, err := oapiTokenPool.NewClientWithResponses(edsURL)
		require.NoError(t, err)

		// Get Fee Input holding
		feeTokenHoldings, err := testhelpers.ListHoldingsForInstrument(t.Context(), participant, amuletInstrumentID)
		require.NoError(t, err)

		transferFactoryCid, transferFactoryDisclosures, choiceContextRaw, err := testhelpers.GetTransferFactory(t.Context(), transferInstructionClient, string(amuletInstrumentID.Admin), participant.PartyID, ccipOwnerPartyID)
		require.NoError(t, err)
		choiceContext, err := contracts.ChoiceContextFromData(choiceContextRaw)
		require.NoError(t, err)

		// Get Token input holding
		tokenHoldings, err := testhelpers.ListHoldingsForInstrument(t.Context(), participant, linkInstrumentId)
		require.NoError(t, err)

		// Get PerPartyRouter
		routerCid := getRouter(t, participant, ccipEdsClient)

		// Get CCIPSender
		senderCid := getSender(t, participant)

		msg := oapiCommon.Message{
			DestinationChainSelector: strconv.FormatUint(ethSelector, 10),
			Executor: struct {
				Address *oapiCommon.RawOrHashedAddress `json:"address,omitempty"`
				Type    oapiCommon.MessageExecutorType `json:"type"`
			}{
				Type: oapiCommon.Empty,
			},
			FeeToken: oapiCommon.InstrumentId{
				Admin: oapiCommon.PartyId(amuletInstrumentID.Admin),
				Id:    string(amuletInstrumentID.Id),
			},
			GasLimit: 0,
			Payload:  "",
			Receiver: hex.EncodeToString(messageReceiver.Bytes()),
			TokenTransfer: &oapiCommon.TokenTransfer{
				Amount: "0.5",
				Token: oapiCommon.InstrumentId{
					Admin: oapiCommon.PartyId(linkInstrumentId.Admin),
					Id:    string(linkInstrumentId.Id),
				},
			},
		}

		tokenPoolAddress, err := eds.GetTokenPoolForToken(t.Context(), ccipEdsClient, contracts.EncodeInstrumentID(*linkInstrumentId))
		require.NoError(t, err)
		tokenPoolSendDisclosure, err := eds.GetTokenPoolSendDisclosure(t.Context(), tokenPoolEdsClient, msg, tokenPoolAddress.InstanceAddress())
		require.NoError(t, err)
		ccipSendDisclosure, err := eds.GetCCIPSendDisclosure(t.Context(), ccipEdsClient, msg, nil, tokenPoolSendDisclosure.RequiredCCVs)
		require.NoError(t, err)
		// EDS returns the default CCV(s) if no input is specified
		defaultCCVAddress, err := contracts.RawInstanceAddressFromString(ccipSendDisclosure.CCVs[0])
		require.NoError(t, err)
		// EDS returns the default Executor (if set) if no specific executor input is specified
		defaultExecutorAddress, err := contracts.RawInstanceAddressFromString(*ccipSendDisclosure.Executor)
		require.NoError(t, err)
		ccvSendDisclosure, err := eds.GetCCVSendDisclosure(t.Context(), ccvEdsClient, msg, defaultCCVAddress.InstanceAddress())
		require.NoError(t, err)
		executorSendDisclosure, err := eds.GetExecutorSendDisclosure(t.Context(), executorEdsClient, msg, defaultExecutorAddress.InstanceAddress(), ccipSendDisclosure.CCVs)
		require.NoError(t, err)

		feeTokenInputCids := make([]types.CONTRACT_ID, len(feeTokenHoldings))
		for i, holding := range feeTokenHoldings {
			feeTokenInputCids[i] = types.CONTRACT_ID(holding.ContractID)
		}
		tokenTransferInputCids := make([]types.CONTRACT_ID, len(tokenHoldings))
		for i, holding := range tokenHoldings {
			tokenTransferInputCids[i] = types.CONTRACT_ID(holding.ContractID)
		}

		sendArgs := sender.Send{
			DestinationChainSelector: types.NUMERIC(msg.DestinationChainSelector),
			Message: core.Canton2AnyMessage{
				Receiver: types.TEXT(msg.Receiver),
				Payload:  types.TEXT(msg.Payload),
				TokenTransfer: &core.TokenTransfer{
					Token:  *linkInstrumentId,
					Amount: types.NUMERIC(msg.TokenTransfer.Amount),
				},
				FeeToken: *amuletInstrumentID,
				ExtraArgs: core.ExtraArgs{
					V3: &core.GenericExtraArgsV3{
						GasLimit: types.INT64(msg.GasLimit),
						Ccvs:     nil,
						Executor: core.ExecutorExtraArg{
							ExecutorUseDefault: &core.ExecutorUseDefault{ExecutorArgs: ""},
						},
						TokenReceiver: "",
						TokenArgs:     "",
					},
				},
			},
			Context:   ccipSendDisclosure.ChoiceContext,
			RouterCid: types.CONTRACT_ID(routerCid),
			FeeTokenInput: sender.FeeTokenInput{
				SenderInputCids:         feeTokenInputCids,
				FeeTokenConfigCid:       types.CONTRACT_ID(ccipSendDisclosure.FeeTokenConfigCid),
				FeeTokenTransferFactory: types.CONTRACT_ID(transferFactoryCid),
				FeeTokenExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
					Context: choiceContext,
					Meta: splice_api_token_metadata_v1.Metadata{
						Values: map[string]types.TEXT{},
					},
				},
			},
			CcvSendInputs: []sender.CCVSendInput{
				{
					CcvAddress:      ccvSendDisclosure.Address.Binding(),
					CcvCid:          types.CONTRACT_ID(ccvSendDisclosure.ContractId),
					CcvExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
				},
			},
			TokenTransferInput: &sender.TokenTransferInput{
				SenderInputCids:  tokenTransferInputCids,
				TokenPoolCid:     types.CONTRACT_ID(tokenPoolSendDisclosure.ContractId),
				PoolExtraContext: tokenPoolSendDisclosure.ChoiceContext,
			},
			ExecutorInput: &sender.ExecutorInput{
				ExecutorCid:          types.CONTRACT_ID(executorSendDisclosure.ContractId),
				ExecutorExtraContext: splice_api_token_metadata_v1.ChoiceContext{},
			},
		}
		allDisclosures := slices.Concat(
			transferFactoryDisclosures,
			tokenPoolSendDisclosure.DisclosedContracts,
			ccipSendDisclosure.DisclosedContracts,
			ccvSendDisclosure.DisclosedContracts,
			executorSendDisclosure.DisclosedContracts,
		)
		resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.NewString(),
				Commands: []*apiv2.Command{{
					Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
						TemplateId:     &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
						ContractId:     senderCid,
						Choice:         "Send",
						ChoiceArgument: ledger.MapToValue(sendArgs),
					}},
				}},
				ActAs:              []string{participant.PartyID},
				DisclosedContracts: allDisclosures,
			},
		})
		require.NoError(t, err)
		fmt.Println("Message sent in Update: ", resp.GetTransaction().GetUpdateId())

		var returnedMessageId string
		for _, event := range resp.GetTransaction().GetEvents() {
			if e, ok := event.GetEvent().(*apiv2.Event_Created); ok {
				if e.Created.GetTemplateId().GetEntityName() == "CCIPMessageSent" {
					ccipMessageSent, err := bindings.UnmarshalCreatedEvent[core.CCIPMessageSent](e.Created)
					require.NoError(t, err)
					returnedMessageId = string(ccipMessageSent.Event.MessageId)
					break
				}
			}
		}
		require.NotEmpty(t, returnedMessageId)
		fmt.Println("Message sent with MessageID: ", returnedMessageId)
	})
}

// getRouter returns the PerPartyRouter for the participant's party, or creates one if it doesn't exist yet.
func getRouter(t *testing.T, participant canton.Participant, ccipEdsClient oapiCCIP.ClientWithResponsesInterface) string {
	// Get active contracts. If one exists, return it
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-runtime", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouter"})
	require.NoError(t, err)
	if len(activeContracts) > 0 {
		return activeContracts[0].GetCreatedEvent().GetContractId()
	}

	t.Logf("No active PerPartyRouter found for party %s, deploying one...", participant.PartyID)
	disclosures, err := eds.GetPerPartyRouterFactoryDisclosure(t.Context(), ccipEdsClient, participant.PartyID)
	require.NoError(t, err)

	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Exercise{Exercise: &apiv2.ExerciseCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-runtime", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"},
					ContractId: disclosures.ContractId,
					Choice:     "CreateRouter",
					ChoiceArgument: ledger.MapToValue(ccipruntime.CreateRouter{
						PartyOwner: types.PARTY(participant.PartyID),
						InstanceId: "default-router",
					}),
				}},
			}},
			ActAs:              []string{participant.PartyID},
			DisclosedContracts: disclosures.DisclosedContracts,
		},
	})
	require.NoError(t, err)
	t.Logf("Created PerPartyRouter in update: %v", resp.GetTransaction().GetUpdateId())

	activeContracts, err = testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-runtime", ModuleName: "CCIP.PerPartyRouter", EntityName: "PerPartyRouterFactory"})
	require.NoError(t, err)
	require.NotEmptyf(t, activeContracts, "Expected to find active PerPartyRouter after creation for party %s", participant.PartyID)

	return activeContracts[0].GetCreatedEvent().GetContractId()
}

func getSender(t *testing.T, participant canton.Participant) string {
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"})
	require.NoError(t, err)
	if len(activeContracts) > 0 {
		return activeContracts[0].GetCreatedEvent().GetContractId()
	}

	t.Logf("No active CCIPSender found for party %s, deploying one...", participant.PartyID)
	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ccipsender"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: participant.PartyID}}},
					}},
				}},
			}},
			ActAs: []string{participant.PartyID},
		},
	})
	require.NoError(t, err)
	t.Logf("Created CCIPSender in update: %v", resp.GetTransaction().GetUpdateId())

	activeContracts, err = testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-sender", ModuleName: "CCIP.CCIPSender", EntityName: "CCIPSender"})
	require.NoError(t, err)
	require.NotEmptyf(t, activeContracts, "Expected to find active CCIPSender after creation for party %s", participant.PartyID)

	return activeContracts[0].GetCreatedEvent().GetContractId()
}

func getReceiver(t *testing.T, participant canton.Participant) string {
	activeContracts, err := testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"})
	require.NoError(t, err)
	if len(activeContracts) > 0 {
		return activeContracts[0].GetCreatedEvent().GetContractId()
	}

	t.Logf("No active CCIPReceiver found for party %s, deploying one...", participant.PartyID)
	resp, err := participant.LedgerServices.Command.SubmitAndWaitForTransaction(t.Context(), &apiv2.SubmitAndWaitForTransactionRequest{
		Commands: &apiv2.Commands{
			CommandId: uuid.NewString(),
			Commands: []*apiv2.Command{{
				Command: &apiv2.Command_Create{Create: &apiv2.CreateCommand{
					TemplateId: &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"},
					CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
						{Label: "instanceId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "ccipreceiver-waitForFinality"}}},
						{Label: "owner", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: participant.PartyID}}},
						{Label: "receiverFinalityConfig", Value: &apiv2.Value{Sum: &apiv2.Value_Variant{Variant: &apiv2.Variant{
							Constructor: "WaitForFinality",
							Value:       &apiv2.Value{Sum: &apiv2.Value_Unit{}},
						}}}},
						{Label: "requiredCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalCCVs", Value: &apiv2.Value{Sum: &apiv2.Value_List{List: &apiv2.List{Elements: nil}}}},
						{Label: "optionalThreshold", Value: &apiv2.Value{Sum: &apiv2.Value_Int64{Int64: 0}}},
					}},
				}},
			}},
			ActAs: []string{participant.PartyID},
		},
	})
	require.NoError(t, err)
	t.Logf("Created CCIPReceiver in update: %v", resp.GetTransaction().GetUpdateId())

	activeContracts, err = testhelpers.ListActiveContractsByTemplateId(t.Context(), participant, &apiv2.Identifier{PackageId: "#ccip-receiver", ModuleName: "CCIP.CCIPReceiver", EntityName: "CCIPReceiver"})
	require.NoError(t, err)
	require.NotEmptyf(t, activeContracts, "Expected to find active CCIPReceiver after creation for party %s", participant.PartyID)

	return activeContracts[0].GetCreatedEvent().GetContractId()
}
