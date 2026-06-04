package load

import (
	"fmt"
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
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// TestEVM2Canton_Load runs WASP RPS=1 against the real EVM→Canton path (message-only).
//
// Requires a running devenv and ../../env-canton-evm-out.toml (same as the basic e2e test).
// EVM source accounts are pre-funded by devenv; no Canton MintTokens/SetupSend.
//
//nolint:paralleltest // single-flight exec on Canton dest; shares env with e2e.
func TestEVM2Canton_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM→Canton load test in short mode")
	}

	configPath := "../../env-canton-evm-out.toml"
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("skipping EVM→Canton load test: %v (start devenv to generate %s)", err, configPath)
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
	cantonDest := discoverCantonDest(t, in, chainMap)
	t.Logf("EVM→Canton load: source=%d dest=%d receiver=%x",
		evmChain.ChainSelector(), cantonDest.Chain.ChainSelector(), cantonDest.Receiver)

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
		utilstests.WaitTimeout(t),
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-evm2canton", loadSchedule(t))
}
