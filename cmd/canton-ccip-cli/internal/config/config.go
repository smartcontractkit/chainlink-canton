package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// UserConfig holds the user-supplied configuration loaded from a YAML file.
type UserConfig struct {
	Network           string       `yaml:"network"`
	Canton            CantonConfig `yaml:"canton"`
	EVM               EVMConfig    `yaml:"evm"`
	CCIPExplorerURL   string       `yaml:"ccip_explorer_url"`
	EVMExplorerURL    string       `yaml:"evm_explorer_url"`
	CantonExplorerURL string       `yaml:"canton_explorer_url"`
}

type CantonConfig struct {
	Disabled                    bool              `yaml:"disabled"`
	AuthType                    string            `yaml:"authType"`
	AuthServerURL               string            `yaml:"authServerURL"`
	AuthClientID                string            `yaml:"authClientID"`
	AuthClientSecret            string            `yaml:"authClientSecret"`
	AuthJWT                     string            `yaml:"authJWT"`
	ParticipantGRPCLedgerAPIURL string            `yaml:"participantGRPCLedgerAPIURL"`
	ValidatorAPIURL             string            `yaml:"validatorAPIURL"`
	UserID                      string            `yaml:"userID"`
	PartyID                     string            `yaml:"partyID"`
	EDSURLs                     map[string]string `yaml:"edsURLs"`
	TokenStandardURLs           map[string]string `yaml:"tokenStandardURLs"`
}

type EVMConfig struct {
	RPCURL        string `yaml:"rpcURL"`
	PrivateKeyHex string `yaml:"privateKeyHex"`
}

// Load reads the YAML config from path and returns a validated UserConfig.
func Load(path string) (*UserConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is required (use --config)")
	}
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file %q not accessible: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &UserConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *UserConfig) validate() error {
	missing := func(field, val string) error {
		if val == "" {
			return fmt.Errorf("config field %s is required", field)
		}

		return nil
	}
	for _, e := range []error{
		missing("canton.participantGRPCLedgerAPIURL", c.Canton.ParticipantGRPCLedgerAPIURL),
		missing("canton.validatorAPIURL", c.Canton.ValidatorAPIURL),
		missing("canton.userID", c.Canton.UserID),
		missing("canton.partyID", c.Canton.PartyID),
		missing("evm.rpcURL", c.EVM.RPCURL),
		missing("evm.privateKeyHex", c.EVM.PrivateKeyHex),
	} {
		if e != nil {
			return e
		}
	}

	return nil
}
