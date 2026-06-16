package prodtestnetenv

import (
	"time"

	"github.com/smartcontractkit/chainlink-canton/scripts/internal/scriptenv"
)

const (
	defaultEnvRelPath = "scripts/prod_testnet/.env"
	// DefaultEDSURL is the prod_testnet Canton EDS deployed at eds.testnet.ccip.chain.link.
	DefaultEDSURL = "https://eds.testnet.ccip.chain.link"
)

// LoadDefault loads scripts/prod_testnet/.env if it exists. Existing process env vars win.
func LoadDefault() (string, error) {
	return scriptenv.Load(defaultEnvRelPath)
}

func Load(path string) error {
	_, err := scriptenv.Load(path)
	return err
}

func String(defaultValue string, keys ...string) string {
	return scriptenv.String(defaultValue, keys...)
}

func Uint64(defaultValue uint64, keys ...string) (uint64, error) {
	return scriptenv.Uint64(defaultValue, keys...)
}

func Int64(defaultValue int64, keys ...string) (int64, error) {
	return scriptenv.Int64(defaultValue, keys...)
}

func Duration(defaultValue time.Duration, keys ...string) (time.Duration, error) {
	return scriptenv.Duration(defaultValue, keys...)
}

func First(keys ...string) string {
	return scriptenv.First(keys...)
}
