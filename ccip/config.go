package ccip

import "github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"

const DefaultCantonConfigPath = "/etc/canton/config.toml"

// Config holds chain-specific configuration for the CCIP chain integration.
type Config struct {
	// ReaderConfigs is a map of canton chain selectors to reader configurations.
	ReaderConfigs map[string]sourcereader.ReaderConfig `toml:"reader_configs"`
}
