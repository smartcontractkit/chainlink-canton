package tests

import (
	"fmt"
	"strings"
)

// CCIPEnv names a CCIP e2e target environment.
type CCIPEnv string

const (
	EnvDevenv      CCIPEnv = "devenv"
	EnvProdTestnet CCIPEnv = "prod-testnet"
)

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

// IsRemote reports whether the environment targets live testnet infrastructure.
func (e CCIPEnv) IsRemote() bool {
	return e == EnvProdTestnet
}
