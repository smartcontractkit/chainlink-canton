package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	tokenscore "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	ccv "github.com/smartcontractkit/chainlink-ccv/build/devenv"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/cciptestinterfaces"
	"github.com/smartcontractkit/chainlink-ccv/build/devenv/evm"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

const (
	envTokenTestConfig     = "CANTON_TOKEN_TEST_CONFIG" //nolint:gosec // env var name, not a credential
	defaultTokenConfigFile = "token_transfer_config.toml"
	directionEVMToCanton   = "evm_to_canton"
	directionCantonToEVM   = "canton_to_evm"
)

// TokenLane describes a resolved token pool pairing for load or e2e tests.
type TokenLane struct {
	PoolRef datastore.AddressRef
	// TransferAmount is the per-message token transfer amount.
	TransferAmount *big.Int
	// ExecutionGasLimit is the per-message execution gas limit.
	ExecutionGasLimit uint32
	// FinalityConfig is the per-message finality configuration.
	FinalityConfig protocol.Finality
	// SrcToken is the source-chain token address; set for EVM source, empty for Canton source.
	// TODO: support multiple source tokens when production tests cover more than one instrument.
	SrcToken protocol.UnknownAddress
	// DestTokenBySelector maps each EVM destination selector to its token address
	// for balance assertions. Canton destinations are omitted.
	DestTokenBySelector map[uint64]protocol.UnknownAddress
}

type tokenConfigTOML struct {
	EVMToCanton tokenDirectionTOML `toml:"evm_to_canton"`
	CantonToEVM tokenDirectionTOML `toml:"canton_to_evm"`
}

type tokenDirectionTOML struct {
	PoolType          string `toml:"pool_type"`
	PoolVersion       string `toml:"pool_version"`
	PoolQualifier     string `toml:"pool_qualifier"`
	TransferAmount    string `toml:"transfer_amount"`
	ExecutionGasLimit uint32 `toml:"execution_gas_limit"`
	FinalityConfig    uint32 `toml:"finality_config"`
}

// ResolveTokenLane loads send params from TOML, matches the pool on the source chain,
// and resolves source/destination token addresses for the requested destinations.
func ResolveTokenLane(
	t *testing.T,
	in *ccv.Cfg,
	lib ccv.Lib,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
	srcSelector uint64,
	destSelectors []uint64,
) TokenLane {
	t.Helper()

	// Load the CLDF environment and datastore used to resolve on-chain addresses.
	env, err := lib.CLDFEnvironment()
	require.NoError(t, err)
	require.NotNil(t, env)

	srcChain, ok := chainMap[srcSelector]
	require.True(t, ok, "source chain %d not in harness chain map", srcSelector)

	// Pick the TOML block from the source chain family (Canton<->EVM only).
	direction := directionEVMToCanton
	if isCantonSelector(srcSelector) {
		direction = directionCantonToEVM
	}
	// Read pool identity and per-message send params from token_transfer_config.toml.
	poolRef, transferAmount, gasLimit, finalityConfig := loadTokenDirection(t, direction)

	// List token transfer configs deployed on the source chain for the requested destinations.
	srcProvider := tokenConfigProvider(srcChain)
	cfgs, err := srcProvider.GetTokenTransferConfigs(env, srcSelector, destSelectors, in.EnvironmentTopology)
	require.NoError(t, err, "get token transfer configs for source chain %d", srcSelector)

	// Require exactly one config whose pool ref matches the TOML declaration.
	cfg := selectTokenConfig(t, cfgs, poolRef, srcSelector)

	// Fail fast if any requested destination is missing from that lane's RemoteChains.
	for _, sel := range destSelectors {
		if _, present := cfg.RemoteChains[sel]; !present {
			t.Fatalf("destination %d not configured for pool %s on chain %d (have %v)",
				sel, poolRefString(poolRef), srcSelector, sortedRemoteSelectors(cfg))
		}
	}

	// Assemble the lane with TOML send params; token addresses are filled in below.
	lane := TokenLane{
		PoolRef:             poolRef,
		TransferAmount:      transferAmount,
		ExecutionGasLimit:   gasLimit,
		FinalityConfig:      finalityConfig,
		DestTokenBySelector: make(map[uint64]protocol.UnknownAddress, len(destSelectors)),
	}

	// EVM source: resolve the ERC-20 from the matched config's TokenRef in the datastore.
	// Canton source: instrument is chosen in SetupSend, so SrcToken stays empty.
	if !isCantonSelector(srcSelector) {
		srcToken, err := resolveTokenRef(env.DataStore, srcSelector, cfg.TokenRef)
		require.NoError(t, err, "resolve source token on chain %d", srcSelector)
		lane.SrcToken = srcToken
	}

	// EVM destinations: look up each dest chain's config for the remote pool and resolve TokenRef.
	// Canton destinations are skipped (no EVM-style token address for balance checks).
	for _, sel := range destSelectors {
		if isCantonSelector(sel) {
			continue
		}
		lane.DestTokenBySelector[sel] = resolveDestToken(t, env, in, chainMap, srcSelector, sel, cfg.RemoteChains[sel], poolRef)
	}

	return lane
}

func loadTokenDirection(t *testing.T, direction string) (poolRef datastore.AddressRef, transferAmount *big.Int, gasLimit uint32, finalityConfig protocol.Finality) {
	t.Helper()

	path := os.Getenv(envTokenTestConfig)
	if path == "" {
		path = defaultTokenConfigPath()
	}

	var cfg tokenConfigTOML
	_, err := toml.DecodeFile(path, &cfg)
	require.NoError(t, err, "decode token transfer config %q (set %s to override)", path, envTokenTestConfig)

	var dir tokenDirectionTOML
	switch direction {
	case directionEVMToCanton:
		dir = cfg.EVMToCanton
	case directionCantonToEVM:
		dir = cfg.CantonToEVM
	default:
		t.Fatalf("unknown token transfer direction %q (expected %q or %q)", direction, directionEVMToCanton, directionCantonToEVM)
	}

	require.NotEmpty(t, dir.PoolType, "%s: pool_type is required for direction %q", path, direction)
	require.NotEmpty(t, dir.PoolVersion, "%s: pool_version is required for direction %q", path, direction)
	require.NotEmpty(t, dir.PoolQualifier, "%s: pool_qualifier is required for direction %q", path, direction)

	version, err := semver.NewVersion(dir.PoolVersion)
	require.NoError(t, err, "%s: invalid pool_version %q for direction %q", path, dir.PoolVersion, direction)

	amount, ok := new(big.Int).SetString(strings.TrimSpace(dir.TransferAmount), 10)
	require.True(t, ok && amount.Sign() > 0, "%s: transfer_amount %q must be a positive integer for direction %q", path, dir.TransferAmount, direction)

	require.NotZero(t, dir.ExecutionGasLimit, "%s: execution_gas_limit is required for direction %q", path, direction)

	poolRef = datastore.AddressRef{
		Type:      datastore.ContractType(dir.PoolType),
		Version:   version,
		Qualifier: dir.PoolQualifier,
	}

	return poolRef, amount, dir.ExecutionGasLimit, protocol.Finality(dir.FinalityConfig)
}

func selectTokenConfig(t *testing.T, cfgs []tokenscore.TokenTransferConfig, poolRef datastore.AddressRef, srcSelector uint64) tokenscore.TokenTransferConfig {
	t.Helper()

	matches := make([]tokenscore.TokenTransferConfig, 0, 1)
	for _, cfg := range cfgs {
		if poolRefEqual(cfg.TokenPoolRef, poolRef) {
			matches = append(matches, cfg)
		}
	}
	switch len(matches) {
	case 1:
		return matches[0]
	case 0:
		t.Fatalf("no token transfer config on chain %d matches pool %s (have %s)",
			srcSelector, poolRefString(poolRef), poolRefsString(cfgs))
	default:
		t.Fatalf("pool %s matched %d configs on chain %d (expected one): %s",
			poolRefString(poolRef), len(matches), srcSelector, poolRefsString(matches))
	}

	return tokenscore.TokenTransferConfig{}
}

func resolveDestToken(
	t *testing.T,
	env *deployment.Environment,
	in *ccv.Cfg,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
	srcSelector uint64,
	destSelector uint64,
	remoteChain tokenscore.RemoteChainConfig[*datastore.AddressRef, datastore.AddressRef],
	srcPoolRef datastore.AddressRef,
) protocol.UnknownAddress {
	t.Helper()

	require.NotNil(t, remoteChain.RemotePool, "remote pool unset for destination %d (source pool %s)",
		destSelector, poolRefString(srcPoolRef))

	destChain, ok := chainMap[destSelector]
	require.True(t, ok, "destination chain %d not in harness chain map", destSelector)
	destProvider := tokenConfigProvider(destChain)
	destCfgs, err := destProvider.GetTokenTransferConfigs(env, destSelector, []uint64{srcSelector}, in.EnvironmentTopology)
	require.NoError(t, err, "get token transfer configs for destination chain %d", destSelector)

	for _, destCfg := range destCfgs {
		if poolRefEqual(destCfg.TokenPoolRef, *remoteChain.RemotePool) {
			destToken, err := resolveTokenRef(env.DataStore, destSelector, destCfg.TokenRef)
			require.NoError(t, err, "resolve destination token on chain %d", destSelector)
			return destToken
		}
	}
	t.Fatalf("no config on destination chain %d matches remote pool %s (have %s)",
		destSelector, poolRefString(*remoteChain.RemotePool), poolRefsString(destCfgs))

	return protocol.UnknownAddress{}
}

func tokenConfigProvider(chain cciptestinterfaces.CCIP17) cciptestinterfaces.TokenConfigProvider {
	if p, ok := chain.(cciptestinterfaces.TokenConfigProvider); ok {
		return p
	}

	return evm.NewEmptyCCIP17EVM()
}

func resolveTokenRef(ds datastore.DataStore, chainSelector uint64, ref datastore.AddressRef) (protocol.UnknownAddress, error) {
	addrRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return protocol.UnknownAddress{}, fmt.Errorf("resolve token %s on chain %d: %w", poolRefString(ref), chainSelector, err)
	}

	return protocol.NewUnknownAddressFromHex(addrRef.Address)
}

func isCantonSelector(selector uint64) bool {
	family, err := chainsel.GetSelectorFamily(selector)
	return err == nil && family == chainsel.FamilyCanton
}

func poolRefEqual(a, b datastore.AddressRef) bool {
	if a.Type != b.Type || a.Qualifier != b.Qualifier {
		return false
	}
	switch {
	case a.Version == nil && b.Version == nil:
		return true
	case a.Version == nil || b.Version == nil:
		return false
	default:
		return a.Version.Equal(b.Version)
	}
}

func defaultTokenConfigPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return defaultTokenConfigFile
	}

	return filepath.Join(filepath.Dir(file), defaultTokenConfigFile)
}

// Helpers for logging //
func poolRefsString(cfgs []tokenscore.TokenTransferConfig) string {
	if len(cfgs) == 0 {
		return "none"
	}
	parts := make([]string, len(cfgs))
	for i, cfg := range cfgs {
		parts[i] = poolRefString(cfg.TokenPoolRef)
	}
	sort.Strings(parts)

	return strings.Join(parts, ", ")
}

func poolRefString(ref datastore.AddressRef) string {
	ver := "<nil>"
	if ref.Version != nil {
		ver = ref.Version.String()
	}

	return fmt.Sprintf("%s/%s/%s", ref.Type, ver, ref.Qualifier)
}

func sortedRemoteSelectors(cfg tokenscore.TokenTransferConfig) []uint64 {
	out := make([]uint64, 0, len(cfg.RemoteChains))
	for sel := range cfg.RemoteChains {
		out = append(out, sel)
	}
	slices.Sort(out)

	return out
}
