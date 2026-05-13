package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/scripts/staging/internal/stagingenv"
	"github.com/smartcontractkit/chainlink-ccip/ccv/chains/evm/gobindings/generated/latest/onramp"
	feequoter "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/fee_quoter"
	routerwrapper "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	factoryburnminterc20 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

const (
	defaultExecutionGas         = uint64(100_000)
	noExecutionTag              = "0xEBa517d200000000000000000000000000000000"
	defaultReceiptTimeout       = 2 * time.Minute
	defaultMessageData          = "Hello from EVM!"
	genericExtraArgsV3Tag       = uint32(0xa69dd4aa)
	defaultRequestedFinality    = uint64(1) // FinalityCodec: block depth 1 (FTF-style); 0 = wait-for-finality
	noGasPriceAvailableSelector = "a9674069"
)

func main() {
	if _, err := stagingenv.LoadDefault(); err != nil {
		fatalf("load scripts/staging/.env: %v", err)
	}

	srcSelectorDefault, err := stagingenv.Uint64(0, "STAGING_EVM_TO_CANTON_SRC_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}
	dstSelectorDefault, err := stagingenv.Uint64(0, "STAGING_EVM_TO_CANTON_DST_SELECTOR")
	if err != nil {
		fatalf("%v", err)
	}
	waitTimeoutDefault, err := stagingenv.Duration(defaultReceiptTimeout, "STAGING_EVM_TO_CANTON_WAIT_TIMEOUT")
	if err != nil {
		fatalf("%v", err)
	}
	requestedFinalityDefault, err := stagingenv.Uint64(defaultRequestedFinality, "STAGING_EVM_TO_CANTON_REQUESTED_FINALITY")
	if err != nil {
		fatalf("%v", err)
	}

	var (
		rpcURL             = flag.String("rpc", stagingenv.String("", "STAGING_EVM_TO_CANTON_RPC_URL"), "Source chain RPC URL")
		privateKeyHex      = flag.String("private-key", stagingenv.String("", "STAGING_EVM_TO_CANTON_PRIVATE_KEY", "PRIVATE_KEY"), "Hex private key for the EVM sender")
		routerAddr         = flag.String("router", stagingenv.String("", "STAGING_EVM_TO_CANTON_ROUTER"), "Source Router address")
		onRampAddr         = flag.String("onramp", stagingenv.String("", "STAGING_EVM_TO_CANTON_ONRAMP"), "Source OnRamp address")
		cvResolverAddr     = flag.String("committee-verifier-resolver", stagingenv.String("", "STAGING_EVM_TO_CANTON_COMMITTEE_VERIFIER_RESOLVER"), "Source CommitteeVerifierResolver address")
		srcSelector        = flag.Uint64("src", srcSelectorDefault, "Source chain selector")
		dstSelector        = flag.Uint64("dest", dstSelectorDefault, "Destination chain selector")
		receiverParty      = flag.String("receiver-party", stagingenv.String("", "STAGING_EVM_TO_CANTON_RECEIVER_PARTY"), "Destination Canton party ID")
		messageData        = flag.String("data", defaultMessageData, "Message payload")
		executionGasLimit  = flag.Uint64("execution-gas-limit", defaultExecutionGas, "Execution gas limit in extra args")
		waitTimeout        = flag.Duration("wait-timeout", waitTimeoutDefault, "How long to wait for tx receipt")
		requestedFinality  = flag.Uint64("requested-finality", requestedFinalityDefault, "GenericExtraArgsV3 requestedFinalityConfig (FinalityCodec uint32); default 1 = FTF block-depth 1; 0 = wait-for-finality; bit 16 = wait-for-safe")
		tokenAddr          = flag.String("token", stagingenv.String("", "STAGING_EVM_TO_CANTON_TOKEN"), "Sepolia transfer token address")
		tokenPoolAddr      = flag.String("token-pool", stagingenv.String("", "STAGING_EVM_TO_CANTON_TOKEN_POOL"), "Sepolia source token pool address")
		destTokenAddr      = flag.String("dest-token", stagingenv.String("", "STAGING_EVM_TO_CANTON_DEST_TOKEN"), "Canton destination token key as hex; use the hashed instrument ID TAR key, not the raw token ref")
		tokenAmountRaw     = flag.String("token-amount", stagingenv.String("", "STAGING_EVM_TO_CANTON_TOKEN_AMOUNT"), "Raw token amount in token base units; empty disables token transfer")
		tokenReceiverParty = flag.String("token-receiver-party", stagingenv.String("", "STAGING_EVM_TO_CANTON_TOKEN_RECEIVER_PARTY", "STAGING_EVM_TO_CANTON_RECEIVER_PARTY"), "Destination Canton token receiver party ID")
		tokenExtraDataHex  = flag.String("token-extra-data-hex", stagingenv.String("", "STAGING_EVM_TO_CANTON_TOKEN_EXTRA_DATA_HEX"), "Optional token transfer extraData as hex")
	)
	flag.Parse()

	requireString(*rpcURL, "-rpc", "STAGING_EVM_TO_CANTON_RPC_URL")
	requireString(*privateKeyHex, "-private-key", "STAGING_EVM_TO_CANTON_PRIVATE_KEY or PRIVATE_KEY")
	requireString(*routerAddr, "-router", "STAGING_EVM_TO_CANTON_ROUTER")
	requireString(*onRampAddr, "-onramp", "STAGING_EVM_TO_CANTON_ONRAMP")
	requireString(*cvResolverAddr, "-committee-verifier-resolver", "STAGING_EVM_TO_CANTON_COMMITTEE_VERIFIER_RESOLVER")
	requireString(*receiverParty, "-receiver-party", "STAGING_EVM_TO_CANTON_RECEIVER_PARTY")
	requireUint64(*srcSelector, "-src", "STAGING_EVM_TO_CANTON_SRC_SELECTOR")
	requireUint64(*dstSelector, "-dest", "STAGING_EVM_TO_CANTON_DST_SELECTOR")

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx, cancel := context.WithTimeout(context.Background(), *waitTimeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, *rpcURL)
	if err != nil {
		fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	privateKey, err := parsePrivateKey(*privateKeyHex)
	if err != nil {
		fatalf("parse private key: %v", err)
	}

	chainID, err := chainsel.ChainIdFromSelector(*srcSelector)
	if err != nil {
		fatalf("chain id from selector %d: %v", *srcSelector, err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
	if err != nil {
		fatalf("create transactor: %v", err)
	}

	routerAddress := common.HexToAddress(*routerAddr)
	router, err := routerwrapper.NewRouter(routerAddress, client)
	if err != nil {
		fatalf("create router wrapper: %v", err)
	}
	onRampWrapper, err := onramp.NewOnRamp(common.HexToAddress(*onRampAddr), client)
	if err != nil {
		fatalf("create onramp wrapper: %v", err)
	}

	dynCfg, err := onRampWrapper.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		fatalf("onRamp getDynamicConfig: %v", err)
	}
	fq, err := feequoter.NewFeeQuoter(dynCfg.FeeQuoter, client)
	if err != nil {
		fatalf("create fee quoter binding: %v", err)
	}
	gasPx, err := fq.GetDestinationChainGasPrice(&bind.CallOpts{Context: ctx}, *dstSelector)
	if err != nil {
		fatalf("feeQuoter getDestinationChainGasPrice: %v", err)
	}
	if gasPx.Timestamp == 0 {
		fatalf("FeeQuoter %s has no USD/gas price for dest selector %d (timestamp=0). "+
			"Router.getFee will revert with NoGasPriceAvailable. "+
			"On-chain fix: call FeeQuoter.updatePrices from a price-updater address with gasPriceUpdates for this dest "+
			"(same pattern as other staging lanes; see domains/ccv staging operations_reports configure-chains UsdPerUnitGas). "+
			"ExtraArgs / finality are not the cause when this preflight fails.",
			dynCfg.FeeQuoter.Hex(), *dstSelector)
	}

	var (
		tokenAmounts []routerwrapper.ClientEVMTokenAmount
		tokenArgs    []byte
		tokenAmount  *big.Int
	)
	if strings.TrimSpace(*tokenAmountRaw) != "" {
		requireString(*tokenAddr, "-token", "STAGING_EVM_TO_CANTON_TOKEN")

		tokenAmount, err = parseBigInt(*tokenAmountRaw)
		if err != nil {
			fatalf("parse token amount: %v", err)
		}
		_, err := decodeOptionalHex(*tokenExtraDataHex)
		if err != nil {
			fatalf("parse token extra data hex: %v", err)
		}
		tokenAmounts = []routerwrapper.ClientEVMTokenAmount{{
			Token:  common.HexToAddress(*tokenAddr),
			Amount: tokenAmount,
		}}

		if err := ensureERC20Allowance(ctx, client, auth, common.HexToAddress(*tokenAddr), routerAddress, tokenAmount, logger); err != nil {
			fatalf("prepare token allowance: %v", err)
		}
	}

	ccvs := []protocol.CCV{{
		CCVAddress: protocol.UnknownAddress(common.HexToAddress(*cvResolverAddr).Bytes()),
		Args:       nil,
		ArgsLen:    0,
	}}
	extraArgs, err := newGenericExtraArgsV3(uint32(*executionGasLimit), uint32(*requestedFinality), noExecutionTag, nil, tokenArgs, ccvs)
	if err != nil {
		fatalf("encode v3 extra args: %v", err)
	}

	receiverHash := contracts.HashedPartyFromString(*receiverParty)
	msg := routerwrapper.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(receiverHash.Bytes(), 32),
		Data:         []byte(*messageData),
		TokenAmounts: tokenAmounts,
		FeeToken:     common.Address{},
		ExtraArgs:    extraArgs,
	}

	fee, err := router.GetFee(&bind.CallOpts{Context: ctx}, *dstSelector, msg)
	if err != nil {
		hint := explainGetFeeRevert(err, dynCfg.FeeQuoter)
		fatalf("get fee: %v%s (lane: src=%d dest=%d router=%s onRamp=%s feeQuoter=%s resolver=%s receiverParty set=%t)",
			err, hint, *srcSelector, *dstSelector, *routerAddr, *onRampAddr, dynCfg.FeeQuoter.Hex(), *cvResolverAddr, strings.TrimSpace(*receiverParty) != "")
	}
	auth.Value = fee

	logEvt := logger.Info().
		Str("rpc", *rpcURL).
		Str("from", auth.From.Hex()).
		Uint64("srcSelector", *srcSelector).
		Uint64("dstSelector", *dstSelector).
		Str("router", *routerAddr).
		Str("onRamp", *onRampAddr).
		Str("committeeVerifierResolver", *cvResolverAddr).
		Str("receiverParty", *receiverParty).
		Str("receiverHash", receiverHash.Hex()).
		Str("feeWei", fee.String())
	if tokenAmount != nil {
		logEvt = logEvt.
			Str("token", *tokenAddr).
			Str("tokenPool", *tokenPoolAddr).
			Str("destToken", *destTokenAddr).
			Str("tokenAmount", tokenAmount.String()).
			Str("tokenReceiverParty", *tokenReceiverParty)
	}
	logEvt.Msg("sending EVM -> Canton message")

	tx, err := router.CcipSend(auth, *dstSelector, msg)
	if err != nil {
		fatalf("ccip send: %v", err)
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		fatalf("wait mined: %v", err)
	}
	if receipt.Status != 1 {
		fatalf("transaction reverted: tx=%s", tx.Hash().Hex())
	}

	var sentEvent *onramp.OnRampCCIPMessageSent
	eventTopic := (onramp.OnRampCCIPMessageSent{}).Topic()
	for _, lg := range receipt.Logs {
		if len(lg.Topics) == 0 || lg.Topics[0] != eventTopic {
			continue
		}
		ev, parseErr := onRampWrapper.ParseCCIPMessageSent(*lg)
		if parseErr != nil {
			fatalf("parse CCIPMessageSent event: %v", parseErr)
		}
		sentEvent = ev
		break
	}
	if sentEvent == nil {
		fatalf("no CCIPMessageSent event found in tx receipt %s", tx.Hash().Hex())
	}

	result := logger.Info().
		Str("txHash", tx.Hash().Hex()).
		Uint64("blockNumber", receipt.BlockNumber.Uint64()).
		Str("messageID", hexutil.Encode(sentEvent.MessageId[:]))
	if tokenAmount != nil {
		result = result.Str("tokenAmount", tokenAmount.String())
	}
	result.Msg("EVM -> Canton send submitted")
}

func parsePrivateKey(raw string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(raw, "0x"))
}

func parseBigInt(raw string) (*big.Int, error) {
	v, ok := new(big.Int).SetString(strings.TrimSpace(raw), 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer %q", raw)
	}
	return v, nil
}

func decodeOptionalHex(raw string) ([]byte, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	return hexutil.Decode(raw)
}

func ensureERC20Allowance(ctx context.Context, client *ethclient.Client, auth *bind.TransactOpts, tokenAddr, spender common.Address, amount *big.Int, logger zerolog.Logger) error {
	token, err := factoryburnminterc20.NewFactoryBurnMintERC20(tokenAddr, client)
	if err != nil {
		return fmt.Errorf("create token binding: %w", err)
	}
	balance, err := token.BalanceOf(&bind.CallOpts{Context: ctx}, auth.From)
	if err != nil {
		return fmt.Errorf("balanceOf: %w", err)
	}
	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient token balance: have %s need %s", balance.String(), amount.String())
	}
	allowance, err := token.Allowance(&bind.CallOpts{Context: ctx}, auth.From, spender)
	if err != nil {
		return fmt.Errorf("allowance: %w", err)
	}
	if allowance.Cmp(amount) >= 0 {
		return nil
	}

	logger.Info().
		Str("token", tokenAddr.Hex()).
		Str("spender", spender.Hex()).
		Str("currentAllowance", allowance.String()).
		Str("requiredAmount", amount.String()).
		Msg("approving token allowance for CCIP send")

	approveAuth := *auth
	approveAuth.Context = ctx
	tx, err := token.Approve(&approveAuth, spender, amount)
	if err != nil {
		return fmt.Errorf("approve: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("wait approve mined: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("approve reverted: %s", tx.Hash().Hex())
	}
	return nil
}

func newGenericExtraArgsV3(
	gasLimit uint32,
	requestedFinality uint32,
	execAddr string,
	execArgs []byte,
	tokenArgs []byte,
	ccvs []protocol.CCV,
) ([]byte, error) {
	if len(ccvs) > 255 {
		return nil, fmt.Errorf("too many CCVs: %d", len(ccvs))
	}
	if len(execArgs) > 65535 {
		return nil, fmt.Errorf("executor args too long: %d", len(execArgs))
	}
	if len(tokenArgs) > 65535 {
		return nil, fmt.Errorf("token args too long: %d", len(tokenArgs))
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, genericExtraArgsV3Tag); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, gasLimit); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, requestedFinality); err != nil {
		return nil, err
	}
	buf.WriteByte(uint8(len(ccvs)))

	for i, ccv := range ccvs {
		ccvAddr := common.BytesToAddress(ccv.CCVAddress)
		if ccvAddr == (common.Address{}) {
			buf.WriteByte(0)
		} else {
			buf.WriteByte(20)
			buf.Write(ccvAddr.Bytes())
		}
		if len(ccv.Args) > 65535 {
			return nil, fmt.Errorf("CCV[%d] args too long: %d", i, len(ccv.Args))
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(ccv.Args))); err != nil {
			return nil, err
		}
		buf.Write(ccv.Args)
	}

	execAddress := common.HexToAddress(execAddr)
	if execAddress == (common.Address{}) {
		buf.WriteByte(0)
	} else {
		buf.WriteByte(20)
		buf.Write(execAddress.Bytes())
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(len(execArgs))); err != nil {
		return nil, err
	}
	buf.Write(execArgs)

	buf.WriteByte(0)
	if err := binary.Write(buf, binary.BigEndian, uint16(len(tokenArgs))); err != nil {
		return nil, err
	}
	buf.Write(tokenArgs)

	return buf.Bytes(), nil
}

// rpcRevertData walks wrapped RPC errors for JSON-RPC revert data (hex string or []byte).
func rpcRevertData(err error) []byte {
	type dataError interface {
		ErrorData() interface{}
	}
	for err != nil {
		var de dataError
		if errors.As(err, &de) {
			switch d := de.ErrorData().(type) {
			case string:
				if b, e := hexutil.Decode(d); e == nil {
					return b
				}
			case []byte:
				return d
			}
		}
		err = errors.Unwrap(err)
	}
	return nil
}

func explainGetFeeRevert(err error, feeQuoter common.Address) string {
	data := rpcRevertData(err)
	if len(data) < 4 {
		return ""
	}
	if !strings.EqualFold(hexutil.Encode(data[:4]), "0x"+noGasPriceAvailableSelector) {
		return ""
	}
	if len(data) < 36 {
		return fmt.Sprintf(" [decoded: FeeQuoter.NoGasPriceAvailable on %s]", feeQuoter.Hex())
	}
	chainSel := binary.BigEndian.Uint64(data[28:36])
	return fmt.Sprintf(
		" [decoded: FeeQuoter.NoGasPriceAvailable(destChainSelector=%d); push usdPerUnitGas via FeeQuoter.updatePrices on %s]",
		chainSel, feeQuoter.Hex(),
	)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func requireString(value, flagName, envName string) {
	if strings.TrimSpace(value) != "" {
		return
	}
	fatalf("missing %s (or %s)", flagName, envName)
}

func requireUint64(value uint64, flagName, envName string) {
	if value != 0 {
		return
	}
	fatalf("missing %s (or %s)", flagName, envName)
}
