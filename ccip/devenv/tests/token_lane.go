package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
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
	"github.com/stretchr/testify/require"
)

const (
	// envTokenTestConfig overrides the path to the token transfer TOML config.
	envTokenTestConfig = "CANTON_TOKEN_TEST_CONFIG"

	// defaultTokenConfigFile is the committed TOML config, resolved relative to
	// this package's source directory when the env override is not set.
	defaultTokenConfigFile = "token_transfer_config.toml"

	// DirectionEVMToCanton / DirectionCantonToEVM select a TOML config block.
	DirectionEVMToCanton = "evm_to_canton"
	DirectionCantonToEVM = "canton_to_evm"

	// Numeric send-param fallback defaults (token identity is never defaulted).
	CantonToEVMFeeAmount           int64 = 2_000
	CantonToEVMTokenTransferAmount int64 = 1_000
	EVMToCantonTransferAmount      int64 = 100_000_000_000
	EVMDecimalsScale               int64 = 1_000_000_000_000_000_000
)

// TokenLane describes a resolved token pool pairing for load or e2e tests.
type TokenLane struct {
	// PoolRef is the source-chain token pool identity (for logging/debugging).
	PoolRef datastore.AddressRef
	// TransferAmount is the per-message token transfer amount.
	TransferAmount *big.Int
	// ExecutionGasLimit is the per-message execution gas limit.
	ExecutionGasLimit uint32
	// FinalityConfig is the per-message finality configuration.
	FinalityConfig protocol.Finality
	// SrcToken is the source-chain token address; set for EVM source, empty for Canton source.
	SrcToken protocol.UnknownAddress
	// DestTokenBySelector maps each EVM destination selector to its token address
	// for balance assertions. Canton destinations are omitted.
	DestTokenBySelector map[uint64]protocol.UnknownAddress
}

// TokenTransferInput is the TOML-declared token lane to resolve.
type TokenTransferInput struct {
	PoolRef           datastore.AddressRef
	TransferAmount    *big.Int
	ExecutionGasLimit uint32
	FinalityConfig    protocol.Finality
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

// LoadTokenTransferInput loads the token lane input for the given direction from
// the committed TOML config (or the CANTON_TOKEN_TEST_CONFIG override). Token
// identity (type/version/qualifier) is required; numeric send params fall back to
// per-direction defaults when omitted.
func LoadTokenTransferInput(t *testing.T, direction string) TokenTransferInput {
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
	case DirectionEVMToCanton:
		dir = cfg.EVMToCanton
	case DirectionCantonToEVM:
		dir = cfg.CantonToEVM
	default:
		t.Fatalf("unknown token transfer direction %q (expected %q or %q)", direction, DirectionEVMToCanton, DirectionCantonToEVM)
	}

	require.NotEmpty(t, dir.PoolType, "%s: pool_type is required for direction %q", path, direction)
	require.NotEmpty(t, dir.PoolVersion, "%s: pool_version is required for direction %q", path, direction)
	require.NotEmpty(t, dir.PoolQualifier, "%s: pool_qualifier is required for direction %q", path, direction)

	version, err := semver.NewVersion(dir.PoolVersion)
	require.NoError(t, err, "%s: invalid pool_version %q for direction %q", path, dir.PoolVersion, direction)

	defAmount, defGas, defFinality := directionDefaults(direction)

	amount := defAmount
	if strings.TrimSpace(dir.TransferAmount) != "" {
		parsed, ok := new(big.Int).SetString(strings.TrimSpace(dir.TransferAmount), 10)
		require.True(t, ok, "%s: invalid transfer_amount %q for direction %q", path, dir.TransferAmount, direction)
		amount = parsed
	}

	gasLimit := defGas
	if dir.ExecutionGasLimit != 0 {
		gasLimit = dir.ExecutionGasLimit
	}

	finalityConfig := defFinality
	if dir.FinalityConfig != 0 {
		finalityConfig = protocol.Finality(dir.FinalityConfig)
	}

	return TokenTransferInput{
		PoolRef: datastore.AddressRef{
			Type:      datastore.ContractType(dir.PoolType),
			Version:   version,
			Qualifier: dir.PoolQualifier,
		},
		TransferAmount:    amount,
		ExecutionGasLimit: gasLimit,
		FinalityConfig:    finalityConfig,
	}
}

// ResolveTokenLane discovers the token lane declared by input against the source
// chain's GetTokenTransferConfigs, validates that every destination selector is
// configured for that lane, and resolves the source/destination token addresses.
func ResolveTokenLane(
	t *testing.T,
	in *ccv.Cfg,
	lib ccv.Lib,
	chainMap map[uint64]cciptestinterfaces.CCIP17,
	srcSelector uint64,
	destSelectors []uint64,
	input TokenTransferInput,
) TokenLane {
	t.Helper()

	env, err := lib.CLDFEnvironment()
	require.NoError(t, err)
	require.NotNil(t, env)

	srcChain, ok := chainMap[srcSelector]
	require.True(t, ok, "source chain %d not in harness chain map", srcSelector)

	srcProvider := tokenConfigProvider(srcChain)
	cfgs, err := srcProvider.GetTokenTransferConfigs(env, srcSelector, destSelectors, in.EnvironmentTopology)
	require.NoError(t, err, "get token transfer configs for source chain %d", srcSelector)

	matches := make([]tokenscore.TokenTransferConfig, 0, 1)
	for _, cfg := range cfgs {
		if poolRefMatches(cfg.TokenPoolRef, input.PoolRef) {
			matches = append(matches, cfg)
		}
	}
	switch len(matches) {
	case 1: // ok
	case 0:
		t.Fatalf("no token transfer config on chain %d matches pool ref %s; available pool refs: %s",
			srcSelector, formatPoolRef(input.PoolRef), formatAvailablePoolRefs(cfgs))
	default:
		t.Fatalf("pool ref %s matched %d token transfer configs on chain %d (expected exactly one); matched: %s",
			formatPoolRef(input.PoolRef), len(matches), srcSelector, formatAvailablePoolRefs(matches))
	}
	cfg := matches[0]

	// Preflight: every requested destination must have this lane configured.
	for _, sel := range destSelectors {
		if _, present := cfg.RemoteChains[sel]; !present {
			t.Fatalf("destination selector %d not configured for pool ref %s on chain %d; available destinations: %s",
				sel, formatPoolRef(input.PoolRef), srcSelector, formatSelectors(remoteSelectors(cfg)))
		}
	}

	lane := TokenLane{
		PoolRef:             input.PoolRef,
		TransferAmount:      input.TransferAmount,
		ExecutionGasLimit:   input.ExecutionGasLimit,
		FinalityConfig:      input.FinalityConfig,
		DestTokenBySelector: make(map[uint64]protocol.UnknownAddress, len(destSelectors)),
	}

	ds := env.DataStore

	// Source token: EVM source resolves from cfg.TokenRef; Canton source uses the
	// Amulet instrument selected in SetupSend, so SrcToken stays empty.
	if !isCantonSelector(srcSelector) {
		srcToken, err := resolveTokenRef(ds, srcSelector, cfg.TokenRef)
		require.NoError(t, err, "resolve source token on chain %d", srcSelector)
		lane.SrcToken = srcToken
	}

	// Destination tokens: query each EVM destination's own configs and match the
	// pool advertised by the source lane, then resolve its TokenRef. RemoteChainConfig.RemoteToken
	// is left unset by the impls, so the dest token must come from the dest chain's own config.
	for _, sel := range destSelectors {
		if isCantonSelector(sel) {
			continue
		}
		remoteChain := cfg.RemoteChains[sel]
		require.NotNil(t, remoteChain.RemotePool, "remote pool ref unset for destination %d on pool ref %s",
			sel, formatPoolRef(input.PoolRef))

		destChain, ok := chainMap[sel]
		require.True(t, ok, "destination chain %d not in harness chain map", sel)
		destProvider := tokenConfigProvider(destChain)
		destCfgs, err := destProvider.GetTokenTransferConfigs(env, sel, []uint64{srcSelector}, in.EnvironmentTopology)
		require.NoError(t, err, "get token transfer configs for destination chain %d", sel)

		var (
			destToken protocol.UnknownAddress
			found     bool
		)
		for _, destCfg := range destCfgs {
			if poolRefEqual(destCfg.TokenPoolRef, *remoteChain.RemotePool) {
				destToken, err = resolveTokenRef(ds, sel, destCfg.TokenRef)
				require.NoError(t, err, "resolve destination token on chain %d", sel)
				found = true
				break
			}
		}
		require.True(t, found, "no token transfer config on destination chain %d matches remote pool %s; available pool refs: %s",
			sel, formatPoolRef(*remoteChain.RemotePool), formatAvailablePoolRefs(destCfgs))
		lane.DestTokenBySelector[sel] = destToken
	}

	return lane
}

// tokenConfigProvider returns the TokenConfigProvider for a chain. Canton's *Chain
// implements it directly; EVM chains do not (the runtime impl differs from the
// config impl), so we fall back to a fresh EVM config provider.
func tokenConfigProvider(chain cciptestinterfaces.CCIP17) cciptestinterfaces.TokenConfigProvider {
	if p, ok := chain.(cciptestinterfaces.TokenConfigProvider); ok {
		return p
	}
	return evm.NewEmptyCCIP17EVM()
}

func resolveTokenRef(ds datastore.DataStore, chainSelector uint64, ref datastore.AddressRef) (protocol.UnknownAddress, error) {
	addrRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(chainSelector, ref.Type, ref.Version, ref.Qualifier))
	if err != nil {
		return protocol.UnknownAddress{}, fmt.Errorf("resolve token ref %s on chain %d: %w", formatPoolRef(ref), chainSelector, err)
	}
	return protocol.NewUnknownAddressFromHex(addrRef.Address)
}

func isCantonSelector(selector uint64) bool {
	family, err := chainsel.GetSelectorFamily(selector)
	return err == nil && family == chainsel.FamilyCanton
}

func poolRefMatches(have, want datastore.AddressRef) bool {
	return poolRefEqual(have, want)
}

func poolRefEqual(a, b datastore.AddressRef) bool {
	return a.Type == b.Type &&
		versionString(a.Version) == versionString(b.Version) &&
		a.Qualifier == b.Qualifier
}

func versionString(v *semver.Version) string {
	if v == nil {
		return ""
	}
	return v.String()
}

func remoteSelectors(cfg tokenscore.TokenTransferConfig) []uint64 {
	out := make([]uint64, 0, len(cfg.RemoteChains))
	for sel := range cfg.RemoteChains {
		out = append(out, sel)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func formatPoolRef(ref datastore.AddressRef) string {
	return fmt.Sprintf("%s %s (%q)", ref.Type, versionString(ref.Version), ref.Qualifier)
}

func formatAvailablePoolRefs(cfgs []tokenscore.TokenTransferConfig) string {
	if len(cfgs) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(cfgs))
	for _, cfg := range cfgs {
		parts = append(parts, formatPoolRef(cfg.TokenPoolRef))
	}
	sort.Strings(parts)
	return strings.Join(parts, "; ")
}

func formatSelectors(selectors []uint64) string {
	if len(selectors) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(selectors))
	for _, sel := range selectors {
		parts = append(parts, fmt.Sprintf("%d", sel))
	}
	return strings.Join(parts, ", ")
}

func directionDefaults(direction string) (amount *big.Int, gasLimit uint32, finalityConfig protocol.Finality) {
	if direction == DirectionCantonToEVM {
		return big.NewInt(CantonToEVMTokenTransferAmount), 500_000, protocol.Finality(1)
	}
	return big.NewInt(EVMToCantonTransferAmount), 200_000, protocol.Finality(0)
}

func defaultTokenConfigPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return defaultTokenConfigFile
	}
	return filepath.Join(filepath.Dir(file), defaultTokenConfigFile)
}
