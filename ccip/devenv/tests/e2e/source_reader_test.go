package canton

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	ledgerv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/committee_verifier"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/executor"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/mock_receiver"
	offrampoperations "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/offramp"
	onrampoperations "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/deployment/v1_7_0/operations/onramp"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/per_party_router_factory"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e"
	"github.com/smartcontractkit/chainlink-ccv/protocol"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
)

const (
	packageName = "ccip-perpartyrouter"
	moduleName  = "CCIP.PerPartyRouter"
	numMessages = 3
)

// Start the environment required for this test using:
// ccv up env-canton-evm.toml
// from the build/devenv directory.
//
//nolint:paralleltest // we won't run this in parallel.
func TestCantonSourceReader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CantonSourceReader test in short mode")
	}
	ccv.RegisterImplFactory(chain_selectors.FamilyCanton, cantondevenv.NewImplFactory())

	ctx := ccv.Plog.WithContext(t.Context())
	l := zerolog.Ctx(ctx)

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	var cantonChain *blockchain.Input
	for _, bc := range in.Blockchains {
		if bc.Type == blockchain.TypeCanton {
			cantonChain = bc
			break
		}
	}
	require.NotNil(t, cantonChain, "need at least one canton chain for this test")

	var evmChain *blockchain.Input
	for _, bc := range in.Blockchains {
		if bc.Type == blockchain.TypeAnvil {
			evmChain = bc
			break
		}
	}
	require.NotNil(t, evmChain, "need at least one evm chain for this test")

	cantonDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(cantonChain.ChainID, chain_selectors.FamilyCanton)
	require.NoError(t, err)

	evmDetails, err := chain_selectors.GetChainDetailsByChainIDAndFamily(evmChain.ChainID, chain_selectors.FamilyEVM)
	require.NoError(t, err)

	_, e, err := ccv.NewCLDFOperationsEnvironment(in.Blockchains, in.CLDF.DataStore)
	require.NoError(t, err)
	b := ccv.NewDefaultCLDFBundle(e)
	e.OperationsBundle = b

	lib, err := ccv.NewLib(l, configPath, chain_selectors.FamilyEVM, chain_selectors.FamilyCanton)
	require.NoError(t, err)
	chains, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	srcChain := chains[cantonDetails.ChainSelector]
	require.NotNil(t, srcChain)
	destChain := chains[evmDetails.ChainSelector]
	require.NotNil(t, destChain)

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	chain := e.BlockChains.CantonChains()[cantonDetails.ChainSelector]
	participant := chain.Participants[0]
	party := participant.PartyID

	grpcURL := cantonChain.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants[0].GRPCLedgerAPIURL
	require.NotEmpty(t, grpcURL)
	jwt := cantonChain.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants[0].JWT
	require.NotEmpty(t, jwt)

	ccipMessageSentTemplateID := &ledgerv2.Identifier{
		PackageId:  "#" + packageName,
		ModuleName: moduleName,
		EntityName: "CCIPMessageSent",
	}
	t.Logf("ccipMessageSentTemplateID being used: %s", ccipMessageSentTemplateID.String())

	sourceReader, err := sourcereader.NewSourceReader(
		logger.Test(t),
		grpcURL,
		jwt,
		sourcereader.ReaderConfig{
			NodeOperatorParty:         party,
			CCIPOwnerParty:            party,
			CCIPMessageSentTemplateID: fmt.Sprintf("%s:%s:%s", ccipMessageSentTemplateID.PackageId, ccipMessageSentTemplateID.ModuleName, ccipMessageSentTemplateID.EntityName),
		},
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	latestBefore, finalizedBefore, err := sourceReader.LatestAndFinalizedBlock(t.Context())
	require.NoError(t, err)
	require.NotNil(t, latestBefore)
	require.NotNil(t, finalizedBefore)
	t.Logf("latest block: %d, finalized block: %d before sending messages", latestBefore.Number, finalizedBefore.Number)

	// Create a few CCIPMessageSent "events" by exercising the appropriate choice on the TestRouter contract.
	messages := make([]protocol.Message, numMessages)

	addresses := getRelevantAddresses(t, in, cantonDetails, evmDetails)
	for i := range numMessages {
		sendMessageResult, err := srcChain.SendMessage(t.Context(), evmDetails.ChainSelector, cciptestinterfaces.MessageFields{
			Receiver: protocol.UnknownAddress(addresses.evmReceiver.Bytes()),
			Data:     fmt.Appendf(nil, "source-reader-message-%d", i),
		}, cciptestinterfaces.MessageOptions{
			Version:           3,
			ExecutionGasLimit: 100_000,
			CCVs: []protocol.CCV{
				{
					CCVAddress: protocol.UnknownAddress(addresses.cantonDefaultVerifier.Bytes()),
					Args:       []byte{},
					ArgsLen:    0,
				},
			},
			Executor: protocol.UnknownAddress(addresses.cantonExecutor.Bytes()),
		})
		require.NoError(t, err)
		require.NotNil(t, sendMessageResult.Message)
		messages[i] = *sendMessageResult.Message
		t.Logf("sending message seqNr %d messageID %s", sendMessageResult.Message.SequenceNumber, sendMessageResult.Message.MustMessageID().String())
	}

	latestAfter, finalizedAfter, err := sourceReader.LatestAndFinalizedBlock(t.Context())
	require.NoError(t, err)
	require.NotNil(t, latestAfter)
	require.NotNil(t, finalizedAfter)
	t.Logf("latest block: %d, finalized block: %d after sending messages", latestAfter.Number, finalizedAfter.Number)

	// query for the CCIPMessageSent events in between before and after
	events, err := sourceReader.FetchMessageSentEvents(t.Context(), new(big.Int).SetUint64(latestBefore.Number), new(big.Int).SetUint64(latestAfter.Number))
	require.NoError(t, err)
	require.Len(t, events, numMessages)

	// assert that we can find the messages in the events
	for _, event := range events {
		found := false
		for _, msg := range messages {
			if msg.MustMessageID() == event.MessageID {
				found = true
				break
			}
		}
		require.True(t, found)
	}

	var indexerMonitor *ccv.IndexerMonitor
	indexerClient, err := lib.Indexer()
	require.NoError(t, err)
	indexerMonitor, err = ccv.NewIndexerMonitor(
		zerolog.Ctx(ctx).With().Str("component", "indexer-client").Logger(),
		indexerClient)
	require.NoError(t, err)
	require.NotNil(t, indexerMonitor)

	aggregatorClients := make(map[string]*ccv.AggregatorClient)
	for qualifier := range in.AggregatorEndpoints {
		client, err := in.NewAggregatorClientForCommittee(
			zerolog.Ctx(ctx).With().Str("component", fmt.Sprintf("aggregator-client-%s", qualifier)).Logger(),
			qualifier)
		require.NoError(t, err)
		require.NotNil(t, client)
		aggregatorClients[qualifier] = client
		t.Cleanup(func() {
			client.Close()
		})
	}
	defaultAggregatorClient := aggregatorClients[devenvcommon.DefaultCommitteeVerifierQualifier]

	testCtx := e2e.NewTestingContext(t, t.Context(), chains, defaultAggregatorClient, indexerMonitor)
	for _, msg := range messages {
		result, err := testCtx.AssertMessage(msg.MustMessageID(), e2e.AssertMessageOptions{
			TickInterval:            1 * time.Second,
			ExpectedVerifierResults: 1, // just committee verifier
			Timeout:                 tests.WaitTimeout(t),
			AssertVerifierLogs:      false,
			AssertExecutorLogs:      false,
		})
		require.NoError(t, err)
		require.NotNil(t, result.AggregatedResult)
		require.Len(t, result.IndexedVerifications.Results, 1)

		ev, err := destChain.WaitOneExecEventBySeqNo(t.Context(), cantonDetails.ChainSelector, uint64(msg.SequenceNumber), tests.WaitTimeout(t))
		require.NoError(t, err)
		require.Equalf(
			t,
			cciptestinterfaces.ExecutionStateSuccess,
			ev.State,
			"message %d should have been successfully executed, return data: %x",
			msg.SequenceNumber,
			ev.ReturnData,
		)
	}
}

// relevantAddresses are the addresses required to construct a valid CCIP message from Canton -> EVM.
type relevantAddresses struct {
	cantonOnRamp                contracts.InstanceAddress
	cantonExecutor              contracts.InstanceAddress
	cantonDefaultVerifier       contracts.InstanceAddress
	cantonPerPartyRouterFactory contracts.InstanceAddress
	cantonRouter                contracts.InstanceAddress
	evmOffRamp                  common.Address
	evmReceiver                 common.Address
}

// getRelevantAddresses returns the canton and evm addresses required to construct a valid CCIP message from Canton -> EVM.
func getRelevantAddresses(t *testing.T, in *ccv.Cfg, cantonDetails, evmDetails chain_selectors.ChainDetails) relevantAddresses {
	var addresses relevantAddresses

	cantonOnRampRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(onrampoperations.ContractType),
			semver.MustParse(onrampoperations.Deploy.Version()),
			"",
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, cantonOnRampRef.Address)
	t.Logf("canton on ramp address: %s", cantonOnRampRef.Address)
	addresses.cantonOnRamp = contracts.HexToInstanceAddress(cantonOnRampRef.Address)

	cantonPerPartyRouterFactoryRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(per_party_router_factory.ContractType),
			per_party_router_factory.Version,
			"",
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, cantonPerPartyRouterFactoryRef.Address)
	t.Logf("canton per party router factory address: %s", cantonPerPartyRouterFactoryRef.Address)
	addresses.cantonPerPartyRouterFactory = contracts.HexToInstanceAddress(cantonPerPartyRouterFactoryRef.Address)

	// cantonRouterRef, err := in.CLDF.DataStore.Addresses().Get(
	// 	datastore.NewAddressRefKey(
	// 		cantonDetails.ChainSelector,
	// 		datastore.ContractType(routeroperations.ContractType),
	// 		semver.MustParse(routeroperations.Deploy.Version()),
	// 		"",
	// 	),
	// )
	// require.NoError(t, err)
	// require.NotEmpty(t, cantonRouterRef.Address)
	// t.Logf("canton router address: %s", cantonRouterRef.Address)
	addresses.cantonRouter = contracts.HexToInstanceAddress(cantonOnRampRef.Address) // TODO: fix when router is deployed

	cantonDefaultVerifierRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(committee_verifier.ContractType),
			canton_committee_verifier.Version,
			devenvcommon.DefaultCommitteeVerifierQualifier,
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, cantonDefaultVerifierRef.Address)
	t.Logf("canton default verifier address: %s", cantonDefaultVerifierRef.Address)
	addresses.cantonDefaultVerifier = contracts.HexToInstanceAddress(cantonDefaultVerifierRef.Address)

	cantonExecutorAddress, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			cantonDetails.ChainSelector,
			datastore.ContractType(executor.ProxyType),
			semver.MustParse(executor.DeployProxy.Version()),
			devenvcommon.DefaultExecutorQualifier,
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, cantonExecutorAddress.Address)
	t.Logf("canton executor address: %s", cantonExecutorAddress.Address)
	addresses.cantonExecutor = contracts.HexToInstanceAddress(cantonExecutorAddress.Address) // TODO: Mock contract

	evmOffRampRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			evmDetails.ChainSelector,
			datastore.ContractType(offrampoperations.ContractType),
			semver.MustParse(offrampoperations.Deploy.Version()),
			"",
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, evmOffRampRef.Address)
	t.Logf("evm off ramp address: %s", evmOffRampRef.Address)
	addresses.evmOffRamp = common.HexToAddress(evmOffRampRef.Address)

	evmReceiverRef, err := in.CLDF.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(
			evmDetails.ChainSelector,
			datastore.ContractType(mock_receiver.ContractType),
			semver.MustParse(mock_receiver.Deploy.Version()),
			devenvcommon.DefaultReceiverQualifier,
		),
	)
	require.NoError(t, err)
	require.NotEmpty(t, evmReceiverRef.Address)
	t.Logf("evm receiver address: %s", evmReceiverRef.Address)
	addresses.evmReceiver = common.HexToAddress(evmReceiverRef.Address)

	return addresses
}
