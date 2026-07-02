package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// UserConfig holds the user-supplied configuration loaded from a YAML file.
type UserConfig struct {
	Canton            CantonConfig `mapstructure:"canton"`
	EVM               EVMConfig    `mapstructure:"evm"`
	CCIPExplorerURL   string       `mapstructure:"ccip_explorer_url"`
	EVMExplorerURL    string       `mapstructure:"evm_explorer_url"`
	CantonExplorerURL string       `mapstructure:"canton_explorer_url"`
}

type CantonConfig struct {
	AuthServerURL               string `mapstructure:"authServerURL"`
	AuthClientID                string `mapstructure:"authClientID"`
	ParticipantGRPCLedgerAPIURL string `mapstructure:"participantGRPCLedgerAPIURL"`
	ValidatorAPIURL             string `mapstructure:"validatorAPIURL"`
	UserID                      string `mapstructure:"userID"`
	PartyID                     string `mapstructure:"partyID"`
}

type EVMConfig struct {
	RPCURL        string `mapstructure:"rpcURL"`
	PrivateKeyHex string `mapstructure:"privateKeyHex"`
}

// Load reads the YAML config from path and returns a validated UserConfig.
func Load(path string) (*UserConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is required (use --config)")
	}
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file %q not accessible: %w", path, err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &UserConfig{}
	if err := v.Unmarshal(cfg); err != nil {
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
		missing("canton.authServerURL", c.Canton.AuthServerURL),
		missing("canton.authClientID", c.Canton.AuthClientID),
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
