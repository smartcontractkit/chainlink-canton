package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

const DefaultConfigPath = "registry-kit.toml"

// Config is stable operator input for canton-registry-kit (devnet.cv1).
type Config struct {
	Network       string         `toml:"network"`
	ChainSelector uint64         `toml:"chain_selector"`
	Ledger        LedgerConfig   `toml:"ledger"`
	Parties       PartiesConfig  `toml:"parties"`
	CCIP          CCIPConfig     `toml:"ccip"`
	Operator      OperatorConfig `toml:"operator_backend"`
}

// LedgerConfig holds Canton participant API endpoints and auth.
type LedgerConfig struct {
	JSONAPIURL       string                  `toml:"json_api_url"`
	GRPCLedgerAPIURL string                  `toml:"grpc_ledger_api_url"`
	AdminAPIURL      string                  `toml:"admin_api_url"`
	ValidatorAPIURL  string                  `toml:"validator_api_url"`
	UserID           string                  `toml:"user_id"`
	SynchronizerID   string                  `toml:"synchronizer_id"`
	Auth             commonconfig.AuthConfig `toml:"auth"`
}

// PartiesConfig lists Registry role parties on devnet.
type PartiesConfig struct {
	Operator  string `toml:"operator"`
	Provider  string `toml:"provider"`
	Registrar string `toml:"registrar"`
	Holder    string `toml:"holder"`
}

// CCIPConfig references pre-deployed CCIP contracts on the participant.
type CCIPConfig struct {
	TokenAdminRegistryCID  string `toml:"token_admin_registry_cid"`
	CCIPParty              string `toml:"ccip_party"`
	BurnMintPoolInstanceID string `toml:"burn_mint_pool_instance_id"`
}

// OperatorConfig is DA's hosted Utilities operator backend (mint/burn choice context).
type OperatorConfig struct {
	BaseURL string `toml:"base_url"`
}

// Load reads registry-kit.toml from path.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	cfg.applyDefaults()

	return cfg, nil
}

func (c *Config) applyDefaults() {
	if c.Operator.BaseURL == "" {
		c.Operator.BaseURL = "https://api.utilities.digitalasset-dev.com/api/utilities"
	}
}

// Validate checks required fields for devnet CLI usage.
func (c Config) Validate() error {
	if c.Ledger.JSONAPIURL == "" {
		return fmt.Errorf("ledger.json_api_url is required")
	}
	if c.Ledger.GRPCLedgerAPIURL == "" {
		return fmt.Errorf("ledger.grpc_ledger_api_url is required")
	}
	if c.Ledger.UserID == "" {
		return fmt.Errorf("ledger.user_id is required")
	}
	if c.Ledger.SynchronizerID == "" {
		return fmt.Errorf("ledger.synchronizer_id is required")
	}
	if err := c.Ledger.Auth.Validate(); err != nil {
		return fmt.Errorf("ledger.auth: %w", err)
	}
	if c.Parties.Operator == "" {
		return fmt.Errorf("parties.operator is required")
	}
	if c.Parties.Provider == "" {
		return fmt.Errorf("parties.provider is required")
	}
	if c.Parties.Registrar == "" {
		return fmt.Errorf("parties.registrar is required")
	}

	return nil
}

// ActingParty returns the party used for ledger commands when a role flag is omitted.
func (c Config) ActingParty(role string) (string, error) {
	switch role {
	case "", "registrar":
		return c.Parties.Registrar, nil
	case "provider":
		return c.Parties.Provider, nil
	case "operator":
		return c.Parties.Operator, nil
	case "holder":
		if c.Parties.Holder != "" {
			return c.Parties.Holder, nil
		}

		return c.Parties.Registrar, nil
	default:
		return "", fmt.Errorf("unknown party role %q", role)
	}
}
