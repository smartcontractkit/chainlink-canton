package ccip

import "github.com/smartcontractkit/chainlink-canton/ccip/sourcereader"

const DefaultCantonConfigPath = "/etc/canton/config.toml"

// Config holds chain-specific configuration for the CCIP chain integration.
type Config struct {
	// ReaderConfigs is a map of canton chain selectors to reader configurations.
	ReaderConfigs map[string]sourcereader.ReaderConfig `toml:"reader_configs"`
	// BlockchainInfos is a map of canton chain selectors to blockchain information.
	BlockchainInfos map[string]BlockchainInfo `toml:"blockchain_infos"`
}

// Supported authentication types for Canton participant APIs.
const (
	AuthTypeStatic            = "static"
	AuthTypeClientCredentials = "clientCredentials"
	AuthTypeAuthorizationCode = "authorizationCode"
)

// AuthConfig configures authentication for a Canton participant endpoint.
type AuthConfig struct {
	// Type selects the auth scheme: "static", "clientCredentials", or "authorizationCode".
	// Defaults to "static" when omitted (backward compatible).
	Type string `toml:"type"`
	// JWT is a pre-obtained token. Required when Type is "static" or empty.
	JWT string `toml:"jwt,omitempty"`
	// AuthURL is the OIDC authorization server base URL. Required for clientCredentials and authorizationCode.
	AuthURL string `toml:"auth_url,omitempty"`
	// ClientID is the OAuth2 client identifier. Required for clientCredentials and authorizationCode.
	ClientID string `toml:"client_id,omitempty"`
	// ClientSecret is the OAuth2 client secret. Required for clientCredentials only.
	ClientSecret string `toml:"client_secret,omitempty"`
}

// BlockchainInfo holds the network-specific data for a canton chain.
// It will be present in submitted job specs for canton chains, mapped to a specific chain selector.
type BlockchainInfo struct {
	GRPCLedgerAPIURL string `toml:"grpc_ledger_api_url"`
	// JWT is kept for backward compatibility. When Auth.Type is empty or "static",
	// this field is used as the bearer token.
	JWT  string     `toml:"jwt,omitempty"`
	Auth AuthConfig `toml:"auth,omitempty"`
}
