package load

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	utilstests "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	cantondevenv "github.com/smartcontractkit/chainlink-canton/ccip/devenv"
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// TestEVM2Canton_TokenLoad runs WASP RPS=1 against the EVM→Canton token transfer path.
//
// Requires a running devenv and ../../env-canton-evm-out.toml.
//
//nolint:paralleltest // single-flight exec on Canton dest; shares env with e2e.
func TestEVM2Canton_TokenLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM→Canton token load test in short mode")
	}

	configPath := "../../env-canton-evm-out.toml"
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("skipping EVM→Canton token load test: %v (start devenv to generate %s)", err, configPath)
	}

	in, err := ccv.LoadOutput[ccv.Cfg](configPath)
	require.NoError(t, err)

	ctx := ccv.Plog.WithContext(t.Context())
	lib, err := ccv.NewLibFromCCVEnv(&ccv.Plog, configPath, chainsel.FamilyEVM, chainsel.FamilyCanton)
	require.NoError(t, err)

	chainMap, err := lib.ChainsMap(ctx)
	require.NoError(t, err)
	require.NoError(t, devenvtests.WireVerifierObservationFromLib(lib, chainMap))

	evmChain := devenvtests.GetChainFromMap(t, blockchain.TypeAnvil, in, chainMap)
	cantonChain := devenvtests.GetChainFromMap(t, blockchain.TypeCanton, in, chainMap)
	cantonImpl, ok := cantonChain.(*cantondevenv.Chain)
	require.True(t, ok, "Canton dest chain must be *devenv.Chain")
	require.NoError(t, cantonImpl.SetupReceive(ctx))

	lane := devenvtests.ResolveTokenLane(t, in, lib, chainMap, evmChain.ChainSelector(), []uint64{cantonChain.ChainSelector()})
	t.Logf("Token lane: pool=%s transfer=%s srcToken=%x",
		lane.PoolRef.Qualifier,
		lane.TransferAmount.String(),
		lane.SrcToken)

	receiverParticipant, _, err := cantonImpl.ClientParticipant()
	require.NoError(t, err)
	require.NotEmpty(t, receiverParticipant.PartyID)

	receiver, err := cantonChain.GetEOAReceiverAddress()
	require.NoError(t, err)
	cantonDest := cantonTokenLoadDestination(cantonChain, receiver, lane)

	sched := loadSchedule(t)
	estimatedMessages := estimateMessages(sched)
	evmSender, err := evmChain.GetEOAReceiverAddress()
	require.NoError(t, err)
	senderBalance, err := evmChain.GetTokenBalance(ctx, evmSender, lane.SrcToken)
	require.NoError(t, err)
	requiredBalance := new(big.Int).Mul(lane.TransferAmount, big.NewInt(int64(estimatedMessages)))
	t.Logf("EVM sender token balance=%s requiredForRun=%s (estimatedMessages=%d; devenv pre-funds sender)",
		senderBalance.String(), requiredBalance.String(), estimatedMessages)
	if senderBalance.Cmp(requiredBalance) < 0 {
		t.Logf("warning: EVM sender balance may be insufficient for full run")
	}

	t.Cleanup(func() {
		_, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
	})

	ccvAddr, executorAddr := resolveEVMSourceAddrs(t, lib, evmChain.ChainSelector())

	gun, err := NewCCIPLoadGun(
		evmChain,
		[]Destination{cantonDest},
		ccvAddr,
		executorAddr,
		LoadGunOptions{
			ConfirmSend:        evmSourceConfirmSend(evmChain), // TODO: this confirmation will change on prod-testnet PR
			ConfirmExecTimeout: utilstests.WaitTimeout(t),
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-evm2canton-token", sched, "token_transfer")

	totalHoldingsRat, err := testhelpers.GetHoldingsBalance(ctx, receiverParticipant, nil)
	require.NoError(t, err)
	totalHoldingsFloat, _ := new(big.Float).SetRat(totalHoldingsRat).Float64()
	t.Logf("Canton receiver total holdings after load: %.10f (calls=%d)", totalHoldingsFloat, gun.CallCount())
}
