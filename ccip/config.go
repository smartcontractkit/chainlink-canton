package ccip

import (
	"github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"
	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

const DefaultCantonConfigPath = "/etc/canton/config.toml"

// Config holds chain-specific configuration for the CCIP chain integration.
type Config struct {
	// ReaderConfigs is a map of canton chain selectors to reader configurations.
	ReaderConfigs map[string]sourcereader.ReaderConfig `toml:"reader_configs"`
	// BlockchainInfos is a map of canton chain selectors to blockchain information.
	BlockchainInfos map[string]BlockchainInfo `toml:"blockchain_infos"`
}

// BlockchainInfo holds the network-specific data for a canton chain.
// It will be present in submitted job specs for canton chains, mapped to a specific chain selector.
type BlockchainInfo struct {
	GRPCLedgerAPIURL string                  `toml:"grpc_ledger_api_url"`
	Auth             commonconfig.AuthConfig `toml:"auth,omitempty"`
}
