package tests

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const envConfigFile = "CCIP_CONFIG_FILE"

// CCIPEnv names a CCIP e2e target environment.
type CCIPEnv string

const (
	EnvDevenv      CCIPEnv = "devenv"
	EnvProdTestnet CCIPEnv = "prod-testnet"
)

var ccipEnvFlag = flag.String(
	"ccip-env",
	defaultFromEnv("CCIP_ENV", string(EnvDevenv)),
	"CCIP e2e environment: devenv (default) or prod-testnet",
)

func defaultFromEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// ParseCCIPEnv validates and returns a CCIPEnv from its string form.
func ParseCCIPEnv(s string) (CCIPEnv, error) {
	switch CCIPEnv(strings.TrimSpace(s)) {
	case EnvDevenv, EnvProdTestnet:
		return CCIPEnv(strings.TrimSpace(s)), nil
	case "staging":
		return "", fmt.Errorf("ccip env %q is reserved but not yet supported", s)
	case "mainnet":
		return "", fmt.Errorf("ccip env %q is reserved but not yet supported", s)
	default:
		return "", fmt.Errorf("unknown ccip env %q: want devenv or prod-testnet", s)
	}
}

// ConfigPath returns the CCV env output TOML filename under ccip/devenv.
func (e CCIPEnv) ConfigPath() string {
	switch e {
	case EnvDevenv:
		return "env-canton-evm-out.toml"
	case EnvProdTestnet:
		return "env-prod-testnet-out.toml"
	default:
		return ""
	}
}

// ResolveConfigPath returns the CCV env output TOML filename under ccip/devenv.
// When CCIP_CONFIG_FILE is set, its basename is used instead of the default for env.
func ResolveConfigPath(env CCIPEnv) string {
	if override := strings.TrimSpace(os.Getenv(envConfigFile)); override != "" {
		return filepath.Base(override)
	}

	return env.ConfigPath()
}

// IsRemote reports whether the environment targets live testnet infrastructure.
func (e CCIPEnv) IsRemote() bool {
	return e == EnvProdTestnet
}

// ParseEnvFromFlag reads the -ccip-env flag (defaulting from CCIP_ENV).
func ParseEnvFromFlag(t *testing.T) CCIPEnv {
	t.Helper()

	t.Logf("CCIP_ENV env=%q ccip-env flag=%q",
		os.Getenv("CCIP_ENV"),
		*ccipEnvFlag,
	)

	env, err := ParseCCIPEnv(*ccipEnvFlag)
	require.NoError(t, err)

	return env
}
