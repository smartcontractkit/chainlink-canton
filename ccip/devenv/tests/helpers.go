package tests

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/tests/e2e/tcapi"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

const (
	envConfirmExecTimeout     = "CANTON_CONFIRM_EXEC_TIMEOUT"
	defaultConfirmExecTimeout = 5 * time.Minute

	envConfirmSendTimeout = "CANTON_CONFIRM_SEND_TIMEOUT"
)

// ConfirmExecTimeout returns the timeout for Canton ConfirmExecOnDest polling.
// Default is 5 minutes; override with CANTON_CONFIRM_EXEC_TIMEOUT (e.g. "10m").
func ConfirmExecTimeout(t *testing.T) time.Duration {
	t.Helper()

	timeout := defaultConfirmExecTimeout
	if d := os.Getenv(envConfirmExecTimeout); d != "" {
		parsed, err := time.ParseDuration(d)
		require.NoError(t, err, "%s=%q invalid", envConfirmExecTimeout, d)
		timeout = parsed
	}

	return timeout
}

// ConfirmSendTimeout returns the timeout for ConfirmSendOnSource polling.
// Defaults: devenv 15s, prod-testnet 10m; override with CANTON_CONFIRM_SEND_TIMEOUT (e.g. "30s").
func ConfirmSendTimeout(t *testing.T, env CCIPEnv) time.Duration {
	t.Helper()

	var timeout time.Duration
	switch env {
	case EnvDevenv:
		timeout = 15 * time.Second
	case EnvProdTestnet:
		timeout = 10 * time.Minute
	default:
		timeout = 15 * time.Second
	}

	if d := os.Getenv(envConfirmSendTimeout); d != "" {
		parsed, err := time.ParseDuration(d)
		require.NoError(t, err, "%s=%q invalid", envConfirmSendTimeout, d)
		timeout = parsed
	}

	return timeout
}

// EVMToCantonMessageOptions returns standard message options for EVM→Canton sends with FTF.
func EVMToCantonMessageOptions(gasLimit uint32, executor, ccvAddr protocol.UnknownAddress) cciptestinterfaces.MessageOptions {
	return cciptestinterfaces.MessageOptions{
		ExecutionGasLimit: gasLimit,
		FinalityConfig:    cantondevenv.EVMToCantonFinalityConfig,
		Executor:          executor,
		CCVs: []protocol.CCV{
			{CCVAddress: ccvAddr, Args: []byte{}, ArgsLen: 0},
		},
	}
}

// E2EBootstrap holds shared CCIP e2e setup for a selected environment.
type E2EBootstrap struct {
	Env      CCIPEnv
	Cfg      *ccv.Cfg
	Lib      ccv.Lib
	ChainMap map[uint64]cciptestinterfaces.CCIP17
	Canton   *cantondevenv.Chain
	EVM      cciptestinterfaces.CCIP17
}

// BootstrapE2E loads config, wires verifier observation, and resolves Canton + EVM chains.
func BootstrapE2E(t *testing.T, env CCIPEnv) E2EBootstrap {
	t.Helper()

	if env.IsRemote() && os.Getenv("CANTON_GRPC_URL") == "" {
		t.Skip("CANTON_GRPC_URL unset: not configured for remote Canton")
	}

	configPath := filepath.Join("..", "..", ResolveConfigPath(env))
	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	lib, err := ccv.NewLibFromCCVEnv(&ccv.Plog, configPath)
	require.NoError(t, err)

	ctx := t.Context()
	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	require.NoError(t, WireVerifierObservationFromLib(lib, chainMap))

	evmChain := GetChainFromMap(t, blockchain.TypeAnvil, in, chainMap)
	cantonChain := GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok, "Canton chain must be *devenv.Chain")

	if !env.IsRemote() {
		t.Cleanup(func() {
			_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
			require.NoError(t, err)
		})
	}

	return E2EBootstrap{
		Env:      env,
		Cfg:      in,
		Lib:      lib,
		ChainMap: chainMap,
		Canton:   cantonImpl,
		EVM:      evmChain,
	}
}

// SetupCantonSend prepares Canton for send in both envs.
// devenv mints fee Amulet before SetupSend; prod-testnet assumes the party is already funded.
func (b E2EBootstrap) SetupCantonSend(t *testing.T, ctx context.Context, transferAmount uint64) {
	t.Helper()

	fee := uint64(cantondevenv.CantonToEVMFeeAmount)
	if !b.Env.IsRemote() {
		require.NoError(t, b.Canton.MintTokens(ctx, new(big.Rat).SetUint64(fee)))
	}
	require.NoError(t, b.Canton.SetupSend(ctx, fee, new(big.Rat).SetUint64(transferAmount)))
}

// SetupCantonTokenSend prepares Canton for token sends (fee + transfer holdings).
// Devenv mints Amulet; prod-testnet logs required Amulet and transfer instrument balances.
func (b E2EBootstrap) SetupCantonTokenSend(t *testing.T, ctx context.Context, lane TokenLane, sends int) {
	t.Helper()

	fee := uint64(cantondevenv.CantonToEVMTokenTransferFeeAmount)
	feeTotal := new(big.Rat).SetUint64(uint64(sends) * fee)
	transferTotalFP := new(big.Int).Mul(lane.TransferAmount, big.NewInt(int64(sends)))
	transferTotal := new(big.Rat).SetFrac(transferTotalFP, big.NewInt(cantondevenv.CantonFixedPointScale))
	transferPerSend := new(big.Rat).SetFrac(lane.TransferAmount, big.NewInt(cantondevenv.CantonFixedPointScale))

	if !b.Env.IsRemote() {
		require.NoError(t, b.Canton.MintTokens(ctx, feeTotal))
		require.NoError(t, b.Canton.MintTokens(ctx, transferTotal))
	} else {
		instrumentLabel := lane.TransferInstrumentID
		if instrumentLabel == "" {
			instrumentLabel = string(cantondevenv.AMTInstrument)
		}
		t.Logf("prod-testnet: ensure Canton party holds at least %s Amulet for fees (%d sends × %d)",
			feeTotal.FloatString(10), sends, cantondevenv.CantonToEVMTokenTransferFeeAmount)
		t.Logf("prod-testnet: ensure Canton party holds at least %s %s for transfers (%d sends × %s)",
			transferTotal.FloatString(10), instrumentLabel, sends, transferPerSend.FloatString(10))
		if lane.TransferInstrument.Admin != "" {
			participant, _, err := b.Canton.ClientParticipant()
			require.NoError(t, err)
			inst := &lane.TransferInstrument
			filters := []testhelpers.Filter{
				testhelpers.WithHoldingOwner(participant.PartyID),
				testhelpers.WithUnlockedHoldingsOnly(),
			}
			rows, err := testhelpers.ListHoldingsForInstrument(ctx, participant, inst, filters...)
			require.NoError(t, err)
			balance, err := testhelpers.GetHoldingsBalance(ctx, participant, inst, filters...)
			require.NoError(t, err)
			t.Logf("prod-testnet: unlocked holdings for instrument admin=%s id=%s: count=%d total=%s",
				lane.TransferInstrument.Admin, lane.TransferInstrument.Id, len(rows), balance.FloatString(10))
		}
	}

	if lane.TransferInstrument.Admin != "" {
		require.NoError(t, b.Canton.SetupSend(ctx, fee, transferPerSend, lane.TransferInstrument))
	} else {
		require.NoError(t, b.Canton.SetupSend(ctx, fee, transferPerSend))
	}
}

// SetupCantonReceive deploys the client party's PerPartyRouter before inbound messages
// are executed on Canton (e.g. EVM→Canton).
func (b E2EBootstrap) SetupCantonReceive(t *testing.T, ctx context.Context) {
	t.Helper()
	require.NoError(t, b.Canton.SetupReceive(ctx))
}

// ResolveEVMReceiver returns the EVM wallet address derived from PRIVATE_KEY on prod-testnet,
// or the devenv EOA otherwise. Used as the Canton→EVM message receiver, the EVM→Canton sender,
// and in load tests for EVM-side balance checks.
func (b E2EBootstrap) ResolveEVMReceiver(t *testing.T) protocol.UnknownAddress {
	t.Helper()

	if b.Env.IsRemote() {
		pkHex := strings.TrimSpace(os.Getenv("PRIVATE_KEY"))
		require.NotEmpty(t, pkHex, "PRIVATE_KEY required for prod-testnet EVM receiver")
		pkHex = strings.TrimPrefix(pkHex, "0x")
		pk, err := gethcrypto.HexToECDSA(pkHex)
		require.NoError(t, err)
		addr := gethcrypto.PubkeyToAddress(pk.PublicKey)
		return protocol.UnknownAddress(addr.Bytes())
	}

	receiver, err := b.EVM.GetEOAReceiverAddress()
	require.NoError(t, err)
	return receiver
}

// ConfirmEVMSendOnSource confirms an EVM-side CCIP send after SendMessage.
//
// On devenv we poll CCIPMessageSent via ConfirmSendOnSource (local Anvil, small block range).
// On prod-testnet we use sendResult from SendMessage (tx receipt) only: ConfirmSendOnSource
// should be used here but is broken in chainlink-ccv devenv — the event poller scans from
// block 1 to latest on public RPCs and times out with 504 Gateway Timeout.
func (b E2EBootstrap) ConfirmEVMSendOnSource(
	t *testing.T,
	ctx context.Context,
	destSelector uint64,
	seqNo uint64,
	sendResult cciptestinterfaces.MessageSentEvent,
) (cciptestinterfaces.MessageSentEvent, error) {
	t.Helper()

	if b.Env.IsRemote() {
		return sendResult, nil
	}

	return b.EVM.ConfirmSendOnSource(
		ctx,
		destSelector,
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
		ConfirmSendTimeout(t, b.Env),
	)
}

// ConfirmCantonSendOnSource confirms a Canton-side CCIP send after SendMessage.
// Canton tracks the last sent event in memory and polls by sequence number when needed.
func (b E2EBootstrap) ConfirmCantonSendOnSource(
	t *testing.T,
	ctx context.Context,
	destSelector uint64,
	seqNo uint64,
) (cciptestinterfaces.MessageSentEvent, error) {
	t.Helper()

	return b.Canton.ConfirmSendOnSource(
		ctx,
		destSelector,
		cciptestinterfaces.MessageEventKey{SeqNum: seqNo},
		ConfirmSendTimeout(t, b.Env),
	)
}

// SkipIfRemote skips token subtests that are not supported on prod-testnet.
func (b E2EBootstrap) SkipIfRemote(t *testing.T, reason string) {
	t.Helper()
	if b.Env.IsRemote() {
		t.Skip(reason)
	}
}

func GetContractAddress(
	t *testing.T,
	ds datastore.DataStore,
	chainSelector uint64,
	contractType datastore.ContractType,
	version, qualifier, contractName string,
) protocol.UnknownAddress {
	t.Helper()

	ref, err := ds.Addresses().Get(
		datastore.NewAddressRefKey(chainSelector, contractType, semver.MustParse(version), qualifier),
	)
	require.NoErrorf(t, err, "failed to get %s address for chain selector %d, ContractType: %s, ContractVersion: %s",
		contractName, chainSelector, contractType, version)

	addr, err := protocol.NewUnknownAddressFromHex(ref.Address)
	require.NoError(t, err)

	return addr
}

func GetChain(t *testing.T, chainType string, cfg *ccv.Cfg, lib ccv.Lib) cciptestinterfaces.CCIP17 {
	chainMap, err := lib.ChainsMap(t.Context())
	require.NoError(t, err)
	return GetChainFromMap(t, chainType, cfg, chainMap)
}

// GetChainFromMap returns a chain from an existing ChainsMap result. lib.ChainsMap
// constructs new impls on every call, so tests must reuse the same map they wired
// (via WireVerifierObservationFromLib) rather than calling ChainsMap again through GetChain.
func GetChainFromMap(
	t *testing.T,
	chainType string,
	cfg *ccv.Cfg,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
) cciptestinterfaces.CCIP17 {
	t.Helper()

	selector := chainSelectorForType(t, chainType, cfg)

	c, ok := chainMap[selector]
	require.True(t, ok, "chain selector %d not in ChainsMap", selector)

	return c
}

func chainSelectorForType(t *testing.T, chainType string, cfg *ccv.Cfg) uint64 {
	t.Helper()

	var chain *blockchain.Input
	for _, bc := range cfg.Blockchains {
		if bc.Type == chainType {
			chain = bc
			break
		}
	}
	require.NotNil(t, chain, "need at least one chain for this test")

	var family string
	switch chainType {
	case blockchain.TypeCanton:
		family = chainsel.FamilyCanton
	case blockchain.TypeAnvil:
		family = chainsel.FamilyEVM
	default:
		t.Fatalf("unsupported chain type %q", chainType)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(chain.ChainID, family)
	require.NoError(t, err)

	return chainDetails.ChainSelector
}

// verifierObservationAware is implemented by chain impls that need off-chain verifier
// clients inside ConfirmExecOnDest. EVM does not implement this; Canton does.
type verifierObservationAware interface {
	SetVerifierObservation(cantondevenv.VerifierObservation)
}

// WireVerifierObservationFromLib builds aggregator/indexer clients from lib and
// injects them into every chain that implements verifierObservationAware. Call once
// after lib.ChainsMap. Requires NewLibFromCCVEnv (not CLDF-only Lib).
func WireVerifierObservationFromLib(lib ccv.Lib, chains map[uint64]cciptestinterfaces.CCIP17) error {
	obs, err := cantondevenv.VerifierObservationFromLib(lib)
	if err != nil {
		return err
	}
	for _, c := range chains {
		if vo, ok := c.(verifierObservationAware); ok {
			vo.SetVerifierObservation(obs)
		}
	}

	return nil
}

func AssertSingleVerifierResult(
	t *testing.T,
	ctx context.Context,
	lib ccv.Lib,
	messageID [32]byte,
) tcapi.AssertionResult {
	t.Helper()

	obs, err := cantondevenv.VerifierObservationFromLib(lib)
	require.NoError(t, err)

	result, err := cantondevenv.AssertMessageWithVerifierObservation(ctx, obs, messageID, tcapi.AssertMessageOptions{
		TickInterval:            time.Second,
		Timeout:                 ConfirmExecTimeout(t),
		ExpectedVerifierResults: 1,
		AssertVerifierLogs:      false,
		AssertExecutorLogs:      false,
	})
	require.NoError(t, err)
	require.Len(t, result.IndexedVerifications.Results, 1)

	return result
}
