package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/ccip/devenv"
)

func TestTokenDirectionsForEnv(t *testing.T) {
	t.Parallel()

	cfg, err := decodeTokenConfigTOML(defaultTokenConfigPath())
	require.NoError(t, err)

	devenvDirs, err := tokenDirectionsForEnv(cfg, EnvDevenv)
	require.NoError(t, err)
	require.Equal(t, "BurnMintTokenPool", devenvDirs.EVMToCanton.PoolType)
	require.Contains(t, devenvDirs.EVMToCanton.PoolQualifier, "BurnMintTokenPool 2.0.0 [default]")
	require.Equal(t, "LockReleaseTokenPool", devenvDirs.CantonToEVM.PoolType)

	prodDirs, err := tokenDirectionsForEnv(cfg, EnvProdTestnet)
	require.NoError(t, err)
	require.Equal(t, "TEST", prodDirs.EVMToCanton.PoolQualifier)
	require.Equal(t, "LINK", prodDirs.EVMToCanton.RemotePoolQualifier)
	require.Equal(t, "LINK", prodDirs.CantonToEVM.PoolQualifier)
	require.Equal(t, "100", prodDirs.CantonToEVM.TransferAmount)
	require.Equal(t, "link-token", prodDirs.CantonToEVM.TransferInstrumentID)
}

func TestLoadTokenDirection_envSelection(t *testing.T) {
	t.Setenv(envTokenTestConfig, defaultTokenConfigPath())

	devenvDir := loadTokenDirection(t, EnvDevenv, directionEVMToCanton)
	require.Equal(t, "BurnMintTokenPool", string(devenvDir.PoolRef.Type))
	require.Contains(t, devenvDir.PoolRef.Qualifier, "BurnMintTokenPool 2.0.0 [default]")
	require.Nil(t, devenvDir.RemotePoolRef)

	prodDir := loadTokenDirection(t, EnvProdTestnet, directionEVMToCanton)
	require.Equal(t, "TEST", prodDir.PoolRef.Qualifier)
	require.NotNil(t, prodDir.RemotePoolRef)
	require.Equal(t, "LINK", prodDir.RemotePoolRef.Qualifier)
}

func TestLoadTokenDirection_transferAmountParsing(t *testing.T) {
	t.Setenv(envTokenTestConfig, defaultTokenConfigPath())

	evmDir := loadTokenDirection(t, EnvDevenv, directionEVMToCanton)
	require.Equal(t, "100000000001", evmDir.TransferAmount.String())

	cantonDir := loadTokenDirection(t, EnvDevenv, directionCantonToEVM)
	require.Equal(t, big.NewInt(10000000001), cantonDir.TransferAmount)

	prodCantonDir := loadTokenDirection(t, EnvProdTestnet, directionCantonToEVM)
	require.Equal(t, big.NewInt(100), prodCantonDir.TransferAmount)
	require.Equal(t, "link-token", prodCantonDir.TransferInstrumentID)

	expectedWei := new(big.Int).Mul(prodCantonDir.TransferAmount, big.NewInt(devenv.CantonFixedPointToEVMScale))
	require.Equal(t, "10000000000", expectedWei.String())
}

func TestLoadTokenDirection_missingEnvSection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token_transfer_config.toml")
	require.NoError(t, os.WriteFile(path, []byte(`
[devenv.evm_to_canton]
pool_type = "BurnMintTokenPool"
pool_version = "2.0.0"
pool_qualifier = "TEST"
transfer_amount = "1"
execution_gas_limit = 1
finality_config = 0
`), 0o600))

	t.Setenv(envTokenTestConfig, path)

	_, err := tokenDirectionsForEnv(mustDecodeTokenConfig(t, path), EnvProdTestnet)
	require.Error(t, err)
	require.Contains(t, err.Error(), "prod-testnet")
}

func mustDecodeTokenConfig(t *testing.T, path string) tokenConfigTOML {
	t.Helper()

	cfg, err := decodeTokenConfigTOML(path)
	require.NoError(t, err)

	return cfg
}

func decodeTokenConfigTOML(path string) (tokenConfigTOML, error) {
	var cfg tokenConfigTOML
	_, err := toml.DecodeFile(path, &cfg)

	return cfg, err
}

func TestParseInstrumentIDFromTokenRefLabels(t *testing.T) {
	t.Parallel()

	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"
	tokenRef := datastore.AddressRef{
		Labels: datastore.NewLabelSet(
			"instrument-admin:"+ccipOwner,
			"instrument-id:link-token",
		),
	}

	instrument, err := parseInstrumentIDFromTokenRefLabels(tokenRef)
	require.NoError(t, err)
	require.Equal(t, ccipOwner, string(instrument.Admin))
	require.Equal(t, "link-token", string(instrument.Id))

	_, err = parseInstrumentIDFromTokenRefLabels(datastore.AddressRef{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "instrument-admin")
}

func TestNormalizeTransferInstrumentID(t *testing.T) {
	t.Parallel()

	require.Equal(t, "link-token", normalizeTransferInstrumentID("LINK"))
	require.Equal(t, "link-token", normalizeTransferInstrumentID("link"))
	require.Equal(t, "link-token", normalizeTransferInstrumentID(" link "))
	require.Equal(t, "Amulet", normalizeTransferInstrumentID("Amulet"))
	require.Equal(t, "link-token", normalizeTransferInstrumentID("link-token"))
}

func TestCCIPOwnerFromRefLabels(t *testing.T) {
	t.Parallel()

	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"

	fromPrefix := datastore.AddressRef{
		Labels: datastore.NewLabelSet("ccip-owner:" + ccipOwner),
	}
	require.Equal(t, ccipOwner, ccipOwnerFromRefLabels(fromPrefix))

	fromSuffix := datastore.AddressRef{
		Labels: datastore.NewLabelSet("burnminttokenpool-LINK@" + ccipOwner),
	}
	require.Equal(t, ccipOwner, ccipOwnerFromRefLabels(fromSuffix))

	require.Empty(t, ccipOwnerFromRefLabels(datastore.AddressRef{}))
}

func TestResolveInstrumentFromTokenRefWithFallback(t *testing.T) {
	t.Parallel()

	ccipOwner := "ccipOwner::1220e382f4e57b0815e6be737006e381e6b7de448e06bd033ece6df498017879f551"

	tokenWithLabels := datastore.AddressRef{
		Labels: datastore.NewLabelSet(
			"instrument-admin:"+ccipOwner,
			"instrument-id:link-token",
		),
	}
	poolRef := datastore.AddressRef{
		Labels: datastore.NewLabelSet("burnminttokenpool-LINK@" + ccipOwner),
	}

	instrument, err := resolveInstrumentFromTokenRefWithFallback(tokenWithLabels, poolRef, "LINK")
	require.NoError(t, err)
	require.Equal(t, ccipOwner, string(instrument.Admin))
	require.Equal(t, "link-token", string(instrument.Id))

	tokenWithoutLabels := datastore.AddressRef{}
	instrument, err = resolveInstrumentFromTokenRefWithFallback(tokenWithoutLabels, poolRef, "LINK")
	require.NoError(t, err)
	require.Equal(t, ccipOwner, string(instrument.Admin))
	require.Equal(t, "link-token", string(instrument.Id))

	_, err = resolveInstrumentFromTokenRefWithFallback(tokenWithoutLabels, datastore.AddressRef{}, "LINK")
	require.Error(t, err)
	require.Contains(t, err.Error(), "instrument-admin")

	_, err = resolveInstrumentFromTokenRefWithFallback(tokenWithoutLabels, poolRef, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "transfer_instrument_id")
}
