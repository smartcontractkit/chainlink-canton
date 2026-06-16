package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-canton/scripts/prod_testnet/internal/prodtestnetenv"
	ccvofframp "github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/gobindings/generated/latest/offramp"
	ccvexecutor "github.com/smartcontractkit/chainlink-ccv/executor"
	v1 "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/api/handlers/v1"
	indexerclient "github.com/smartcontractkit/chainlink-ccv/indexer/pkg/client"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

const (
	defaultIndexerURL   = ""
	defaultWaitTimeout  = 2 * time.Minute
	defaultPollInterval = 2 * time.Second
)

func main() {
	if _, err := prodtestnetenv.LoadDefault(); err != nil {
		fatalf("load scripts/prod_testnet/.env: %v", err)
	}

	waitTimeoutDefault, err := prodtestnetenv.Duration(defaultWaitTimeout, "PROD_TESTNET_CANTON_TO_EVM_WAIT_TIMEOUT")
	if err != nil {
		fatalf("%v", err)
	}

	var (
		indexerURL              = flag.String("indexer-url", prodtestnetenv.String(defaultIndexerURL, "PROD_TESTNET_CANTON_TO_EVM_INDEXER_URL"), "Indexer base URL")
		rpcURL                  = flag.String("rpc", prodtestnetenv.String("", "PROD_TESTNET_CANTON_TO_EVM_RPC_URL"), "Destination EVM RPC URL")
		privateKeyHex           = flag.String("private-key", prodtestnetenv.String("", "PROD_TESTNET_CANTON_TO_EVM_PRIVATE_KEY", "PRIVATE_KEY"), "Hex private key for the EVM executor")
		messageIDHex            = flag.String("message-id", "", "0x-prefixed message ID to execute")
		offRampOverride         = flag.String("offramp", prodtestnetenv.String("", "PROD_TESTNET_CANTON_TO_EVM_OFFRAMP"), "Destination OffRamp override; default uses the indexed message off-ramp")
		waitTimeout             = flag.Duration("wait-timeout", waitTimeoutDefault, "How long to wait for indexer lookup and tx mining")
		pollInterval            = flag.Duration("poll-interval", defaultPollInterval, "How often to poll the indexer")
		expectedVerifierResults = flag.Int("expected-verifier-results", 1, "Expected number of verifier results before execution")
		gasLimitOverride        = flag.Uint64("gas-limit-override", 0, "OffRamp gas limit override passed to execute")
		txGasLimit              = flag.Uint64("tx-gas-limit", 0, "Top-level transaction gas limit for the execute transaction; 0 uses node estimation")
	)
	flag.Parse()

	requireString(*rpcURL, "-rpc", "PROD_TESTNET_CANTON_TO_EVM_RPC_URL")
	requireString(*privateKeyHex, "-private-key", "PROD_TESTNET_CANTON_TO_EVM_PRIVATE_KEY or PRIVATE_KEY")
	requireString(*indexerURL, "-indexer-url", "PROD_TESTNET_CANTON_TO_EVM_INDEXER_URL")
	requireString(*messageIDHex, "-message-id", "")

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithTimeout(context.Background(), *waitTimeout)
	defer cancel()

	indexer, err := indexerclient.NewIndexerClient(*indexerURL, &http.Client{Timeout: 15 * time.Second})
	if err != nil {
		fatalf("create indexer client: %v", err)
	}

	messageID, err := protocol.NewBytes32FromString(*messageIDHex)
	if err != nil {
		fatalf("parse message id: %v", err)
	}

	resp, err := waitForVerifierResults(ctx, indexer, messageID, *pollInterval, *expectedVerifierResults, &logger)
	if err != nil {
		fatalf("wait for verifier results: %v", err)
	}
	if len(resp.Results) == 0 {
		fatalf("indexer returned no verifier results for %s", messageID.String())
	}

	message := resp.Results[0].VerifierResult.Message
	if _, err := message.MessageID(); err != nil {
		fatalf("construct message id from indexed message: %v", err)
	}

	offRampAddress := common.BytesToAddress(message.OffRampAddress.Bytes())
	if strings.TrimSpace(*offRampOverride) != "" {
		offRampAddress = common.HexToAddress(*offRampOverride)
	}
	if offRampAddress == (common.Address{}) {
		fatalf("empty offramp address; pass -offramp or ensure the indexed message includes one")
	}

	client, err := ethclient.DialContext(ctx, *rpcURL)
	if err != nil {
		fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	privateKey, err := parsePrivateKey(*privateKeyHex)
	if err != nil {
		fatalf("parse private key: %v", err)
	}

	chainID, err := chainsel.ChainIdFromSelector(uint64(message.DestChainSelector))
	if err != nil {
		fatalf("chain id from selector %d: %v", message.DestChainSelector, err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		fatalf("create transactor: %v", err)
	}
	auth.Context = ctx
	auth.GasLimit = *txGasLimit

	offRamp, err := ccvofframp.NewOffRamp(offRampAddress, client)
	if err != nil {
		fatalf("create offramp binding: %v", err)
	}

	msgIDArray := bytes32ToArray(messageID)
	execState, err := offRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, msgIDArray)
	if err != nil {
		fatalf("get execution state: %v", err)
	}
	if ccvexecutor.MessageExecutionState(execState) == ccvexecutor.SUCCESS {
		logger.Info().
			Str("messageID", messageID.String()).
			Uint8("executionState", execState).
			Msg("message already executed on destination chain")
		return
	}

	encodedMessage, err := message.Encode()
	if err != nil {
		fatalf("encode message: %v", err)
	}

	ccvInfo, err := offRamp.GetCCVsForMessage(&bind.CallOpts{Context: ctx}, encodedMessage)
	if err != nil {
		fatalf("get CCVs for message: %v", err)
	}

	verifierResults := make([]protocol.VerifierResult, 0, len(resp.Results))
	for _, result := range resp.Results {
		vr := result.VerifierResult
		if err := vr.ValidateFieldsConsistent(); err != nil {
			logger.Warn().
				Str("messageID", messageID.String()).
				Err(err).
				Str("verifierDestAddress", vr.VerifierDestAddress.String()).
				Msg("skipping inconsistent verifier result")
			continue
		}
		verifierResults = append(verifierResults, vr)
	}
	if len(verifierResults) == 0 {
		fatalf("no valid verifier results available for %s", messageID.String())
	}

	orderedCCVData, orderedCCVs, err := orderCCVData(verifierResults, ccvInfo)
	if err != nil {
		fatalf("order verifier results: %v", err)
	}

	ccvAddresses := make([]common.Address, 0, len(orderedCCVs))
	for _, ccv := range orderedCCVs {
		ccvAddresses = append(ccvAddresses, common.HexToAddress(ccv.String()))
	}

	logger.Info().
		Str("messageID", messageID.String()).
		Str("from", auth.From.Hex()).
		Str("offRamp", offRampAddress.Hex()).
		Uint64("sourceSelector", uint64(message.SourceChainSelector)).
		Uint64("destSelector", uint64(message.DestChainSelector)).
		Int("verifierResults", len(orderedCCVData)).
		Uint64("gasLimitOverride", *gasLimitOverride).
		Uint64("txGasLimit", *txGasLimit).
		Msg("executing Canton -> EVM message manually")

	tx, err := offRamp.Execute(auth, encodedMessage, ccvAddresses, orderedCCVData, uint32(*gasLimitOverride))
	if err != nil {
		fatalf("offramp execute: %v", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		fatalf("wait mined: %v", err)
	}
	if receipt.Status != 1 {
		fatalf("transaction reverted: tx=%s", tx.Hash().Hex())
	}

	finalState, err := offRamp.GetExecutionState(&bind.CallOpts{Context: ctx}, msgIDArray)
	if err != nil {
		fatalf("get final execution state: %v", err)
	}

	logger.Info().
		Str("messageID", messageID.String()).
		Str("txHash", tx.Hash().Hex()).
		Uint64("blockNumber", receipt.BlockNumber.Uint64()).
		Uint8("executionState", finalState).
		Msg("manual Canton -> EVM execution submitted")
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

func orderCCVData(
	ccvDatas []protocol.VerifierResult,
	ccvInfo ccvofframp.GetCCVsForMessage,
) ([][]byte, []protocol.UnknownAddress, error) {
	destVerifierToCCVData := make(map[string]protocol.VerifierResult, len(ccvDatas))
	for _, ccvData := range ccvDatas {
		destVerifierToCCVData[normalizeAddressKey(ccvData.VerifierDestAddress.String())] = ccvData
	}

	orderedCCVData := make([][]byte, 0, len(ccvDatas))
	orderedCCVOfframps := make([]protocol.UnknownAddress, 0, len(ccvDatas))

	for _, ccvAddress := range ccvInfo.RequiredCCVs {
		key := normalizeAddressKey(ccvAddress.Hex())
		data, ok := destVerifierToCCVData[key]
		if !ok {
			return nil, nil, fmt.Errorf(
				"required CCV %s not found in verifier results %v",
				key,
				slices.Collect(mapKeys(destVerifierToCCVData)),
			)
		}

		orderedCCVData = append(orderedCCVData, data.CCVData)
		orderedCCVOfframps = append(orderedCCVOfframps, protocol.UnknownAddress(ccvAddress.Bytes()))
	}

	optionalCount := 0
	for _, ccvAddress := range ccvInfo.OptionalCCVs {
		if optionalCount >= int(ccvInfo.Threshold) {
			break
		}

		data, ok := destVerifierToCCVData[normalizeAddressKey(ccvAddress.Hex())]
		if !ok {
			continue
		}

		orderedCCVData = append(orderedCCVData, data.CCVData)
		orderedCCVOfframps = append(orderedCCVOfframps, protocol.UnknownAddress(ccvAddress.Bytes()))
		optionalCount++
	}

	if optionalCount < int(ccvInfo.Threshold) {
		return nil, nil, fmt.Errorf(
			"not enough optional CCVs found (%d/%d) in verifier results %v",
			optionalCount,
			ccvInfo.Threshold,
			slices.Collect(mapKeys(destVerifierToCCVData)),
		)
	}

	return orderedCCVData, orderedCCVOfframps, nil
}

func mapKeys[K comparable, V any](m map[K]V) func(func(K) bool) {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

func parsePrivateKey(raw string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(raw, "0x"))
}

func normalizeAddressKey(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func bytes32ToArray(messageID protocol.Bytes32) [32]byte {
	var out [32]byte
	copy(out[:], messageID[:])
	return out
}

func requireString(value string, flagName string, envHint string) {
	if strings.TrimSpace(value) == "" {
		if envHint == "" {
			fatalf("missing required value: pass %s", flagName)
		}
		fatalf("missing required value: pass %s or set %s in scripts/prod_testnet/.env", flagName, envHint)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
