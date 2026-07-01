package canton

import (
	"context"
	"math/big"
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	devenvcommon "github.com/smartcontractkit/chainlink-ccv/build/devenv/common"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	canton_committee_verifier "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/committee_verifier"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/executor"
)

type canton2evmSendFixture struct {
	ctx         context.Context //nolint:containedctx // test fixture
	cantonChain cciptestinterfaces.CCIP17
	cantonImpl  *cantondevenv.Chain
	evmSelector uint64
	receiver    protocol.UnknownAddress
	ccvAddr     protocol.UnknownAddress
	executor    protocol.UnknownAddress
}

func setupCanton2EVMSendFixture(t *testing.T) canton2evmSendFixture {
	t.Helper()

	return setupCanton2EVMSendFixtureWithHoldings(t, true)
}

func setupCanton2EVMSendFixtureWithHoldings(t *testing.T, withFeeHoldings bool) canton2evmSendFixture {
	t.Helper()

	configPath := "../../env-canton-evm-out.toml"
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	lib, err := ccv.NewLibFromCCVEnv(&ccv.Plog, configPath)
	require.NoError(t, err)
	ctx := ccv.Plog.WithContext(t.Context())

	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	require.NoError(t, devenvtests.WireVerifierObservationFromLib(lib, chainMap))

	evmChain := devenvtests.GetChainFromMap(t, blockchain.TypeAnvil, in, chainMap)
	cantonChain := devenvtests.GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok)

	if withFeeHoldings {
		require.NoError(t, cantonImpl.MintTokens(ctx, new(big.Rat).SetUint64(uint64(cantondevenv.CantonToEVMFeeAmount))))
		require.NoError(t, cantonImpl.SetupSend(ctx, uint64(cantondevenv.CantonToEVMFeeAmount), nil))
	}

	ds, err := lib.DataStore()
	require.NoError(t, err)
	receiver, err := evmChain.GetEOAReceiverAddress()
	require.NoError(t, err)
	ccvAddr := devenvtests.GetContractAddress(
		t, ds, cantonChain.ChainSelector(),
		datastore.ContractType(canton_committee_verifier.ContractType),
		canton_committee_verifier.Version.String(),
		devenvcommon.DefaultCommitteeVerifierQualifier,
		"canton committee verifier",
	)
	executorAddr := devenvtests.GetContractAddress(
		t, ds, cantonChain.ChainSelector(),
		datastore.ContractType(executor.ContractType),
		executor.Version.String(),
		devenvcommon.DefaultExecutorQualifier,
		"source executor",
	)

	return canton2evmSendFixture{
		ctx:         ctx,
		cantonChain: cantonChain,
		cantonImpl:  cantonImpl,
		evmSelector: evmChain.ChainSelector(),
		receiver:    receiver,
		ccvAddr:     ccvAddr,
		executor:    executorAddr,
	}
}

func (f canton2evmSendFixture) defaultMessageOptions(gasLimit uint32) cciptestinterfaces.MessageOptions {
	return cciptestinterfaces.MessageOptions{
		ExecutionGasLimit: gasLimit,
		FinalityConfig:    1,
		Executor:          f.executor,
		CCVs: []protocol.CCV{
			{
				CCVAddress: f.ccvAddr,
				Args:       []byte{},
				ArgsLen:    0,
			},
		},
	}
}

func (f canton2evmSendFixture) send(fields cciptestinterfaces.MessageFields, opts cciptestinterfaces.MessageOptions) error {
	_, err := f.cantonChain.SendMessage(f.ctx, f.evmSelector, fields, opts, 3)

	return err
}

//nolint:paralleltest // devenv is single-flight for Canton holdings.
func TestCanton2EVM_SendValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Canton2EVM_SendValidation in short mode")
	}

	t.Run("Payload exceeds maxDataBytes", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		maxDataBytes, err := f.cantonImpl.GetMaxDataBytes(f.ctx, f.evmSelector)
		require.NoError(t, err)
		require.Positive(t, maxDataBytes)

		err = f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     make([]byte, maxDataBytes+1),
			},
			f.defaultMessageOptions(200_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "payload exceeds destination maxDataBytes")
	})

	t.Run("Gas limit exceeds maxPerMsgGasLimit", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		maxGas, err := f.cantonImpl.GetMaxPerMsgGasLimit(f.ctx, f.evmSelector)
		require.NoError(t, err)
		require.Positive(t, maxGas)

		err = f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("gas limit too high"),
			},
			f.defaultMessageOptions(maxGas+1),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "gas limit too high")
	})

	// TODO: SendMessage always builds GenericExtraArgsV3; invalid tag requires a lower-level send path
	t.Run("Invalid extraArgs tag", func(t *testing.T) {
		t.Skip("")
	})

	t.Run("Bad chain selector", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		_, err := f.cantonChain.SendMessage(
			f.ctx,
			9_999_999_999_999_999_999,
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("bad dest selector"),
			},
			f.defaultMessageOptions(200_000),
			3,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported destination chain selector")
	})

	t.Run("Invalid fee token", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		f.cantonImpl.SetFeeTokenInstrumentForTest(splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY("invalid-fee-admin"),
			Id:    types.TEXT("NotARealFeeToken"),
		})

		err := f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("invalid fee token"),
			},
			f.defaultMessageOptions(200_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no token config registered for fee token")
	})

	t.Run("No fee holding supplied", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		f.cantonImpl.ClearFeeHoldingForTest()

		err := f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("no fee holding"),
			},
			f.defaultMessageOptions(200_000),
		)
		require.Error(t, err)
		// Error: canton SendMessage: next fee holding CID is unset; call SetupSend after minting
		require.Contains(t, err.Error(), "next fee holding CID is unset")
	})

	t.Run("Insufficient fee balance", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		// Party already has seeded Amulet from devenv; force a too-small fee holding instead.
		smallCID, err := f.cantonImpl.MintTokensReturningCID(f.ctx, "0.001")
		require.NoError(t, err)
		f.cantonImpl.SetNextFeeCIDForTest(smallCID)

		err = f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("insufficient fee"),
			},
			f.defaultMessageOptions(200_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Insufficient funds")
	})

	t.Run("Unsupported token on lane", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		require.NoError(t, f.cantonImpl.MintTokens(f.ctx, new(big.Rat).SetUint64(1000)))
		require.NoError(t, f.cantonImpl.SetupSend(f.ctx, uint64(cantondevenv.CantonToEVMFeeAmount), new(big.Rat).SetUint64(1000)))
		f.cantonImpl.SetTransferTokenInstrumentForTest(splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY("unsupported-token-admin"),
			Id:    types.TEXT("UnsupportedToken"),
		})

		err := f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("unsupported token"),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount: big.NewInt(1000),
				},
			},
			f.defaultMessageOptions(500_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "wrong pool for Message.TokenTransfer")
	})

	t.Run("Zero token amount", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		require.NoError(t, f.cantonImpl.MintTokens(f.ctx, new(big.Rat).SetUint64(1000)))
		require.NoError(t, f.cantonImpl.SetupSend(f.ctx, uint64(cantondevenv.CantonToEVMFeeAmount), new(big.Rat).SetUint64(1000)))

		err := f.send(
			cciptestinterfaces.MessageFields{
				Receiver: f.receiver,
				Data:     []byte("zero token amount"),
				TokenAmount: cciptestinterfaces.TokenAmount{
					Amount: big.NewInt(0),
				},
			},
			f.defaultMessageOptions(500_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "token transfer amount must be positive")
	})

	t.Run("Invalid receiver", func(t *testing.T) {
		f := setupCanton2EVMSendFixture(t)
		// EVM dest addresses are 20 bytes; wrong length is rejected at PrepareSend.
		// Note: 20 zero bytes (0x0) is a valid encoded EVM address and is accepted on Canton send.
		err := f.send(
			cciptestinterfaces.MessageFields{
				Receiver: make([]byte, 19),
				Data:     []byte("invalid receiver length"),
			},
			f.defaultMessageOptions(200_000),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "InvalidDestChainAddress: address length (19) does not match expected length (20) for the destination chain")
	})
}
