package load

import (
	"testing"

	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	_ "github.com/smartcontractkit/chainlink-ccv/build/devenv/evm" // register EVM ImplFactory
	"github.com/stretchr/testify/require"

	_ "github.com/smartcontractkit/chainlink-canton/ccip/devenv" // registers Canton via init
	devenvtests "github.com/smartcontractkit/chainlink-canton/ccip/devenv/tests"
)

// TestEVM2Canton_Load runs WASP RPS=1 against the real EVM→Canton path (message-only).
//
// Devenv: requires a running devenv and env-canton-evm-out.toml; EVM accounts are pre-funded.
// Prod-testnet: set CANTON_GRPC_URL, CANTON_PARTY_ID, CANTON_AUTH_*, and PRIVATE_KEY; Canton
// party must already be funded (no MintTokens on this path).
//
//nolint:paralleltest // single-flight exec on Canton dest; shares env with e2e.
func TestEVM2Canton_Load(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping EVM→Canton load test in short mode")
	}

	env := devenvtests.ParseEnvFromFlag(t)
	boot := devenvtests.BootstrapE2E(t, env)
	ctx := ccv.Plog.WithContext(t.Context())
	boot.SetupCantonReceive(t, ctx)

	cantonDest := discoverCantonDestFromBoot(t, boot)
	t.Logf("EVM→Canton load: source=%d dest=%d receiver=%x",
		boot.EVM.ChainSelector(), cantonDest.Chain.ChainSelector(), cantonDest.Receiver)

	ccvAddr, _ := resolveEVMSourceAddrs(t, boot.Lib, boot.EVM.ChainSelector())

	gun, err := NewCCIPLoadGun(
		boot.EVM,
		[]Destination{cantonDest},
		ccvAddr,
		nil,
		LoadGunOptions{
			ConfirmSend:        EVMSourceConfirmSend(boot),
			ConfirmExecTimeout: devenvtests.ConfirmExecTimeout(t),
			SkipExecConfirm:    false,
		},
	)
	require.NoError(t, err)

	runWASP(t, gun, "canton-load-evm2canton", loadSchedule(t), "message_only", false, boot.Cfg.IndexerEndpoints)
}
