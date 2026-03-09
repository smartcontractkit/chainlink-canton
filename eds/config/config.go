package config

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

type Config struct {
	ChainSelector string       `toml:"chain_selector" validate:"required"`
	Server        ServerConfig `toml:"server" validate:"required"`
	Node          NodeConfig   `toml:"node" validate:"required"`
	Contracts     Contracts    `toml:"contracts" validate:"required"`
}

type ServerConfig struct {
	Host string `toml:"host" validate:"required"`
	Port uint16 `toml:"port" validate:"required,port"`
}

type NodeConfig struct {
	URL        string                  `toml:"url" validate:"required,url"`
	AuthConfig commonconfig.AuthConfig `toml:"auth" validate:"required"`
	MaxRetries int                     `toml:"max_retries"`
}

type Contracts struct {
	PerPartyRouterFactory ContractIdentifier   `toml:"per_party_router_factory" validate:"required"`
	OnRamp                ContractIdentifier   `toml:"on_ramp" validate:"required"`
	OffRamp               ContractIdentifier   `toml:"off_ramp" validate:"required"`
	GlobalConfig          ContractIdentifier   `toml:"global_config" validate:"required"`
	TokenAdminRegistry    ContractIdentifier   `toml:"token_admin_registry" validate:"required"`
	RMNRemote             ContractIdentifier   `toml:"rmn_remote" validate:"required"`
	FeeQuoter             ContractIdentifier   `toml:"fee_quoter" validate:"required"`
	CCVs                  []ContractIdentifier `toml:"ccvs"`
}

type ContractIdentifier struct {
	PartyID         string                    `toml:"party_id" validate:"required"`
	InstanceAddress contracts.InstanceAddress `toml:"instance_address" validate:"required"`
}

func (cfg *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := commonconfig.RegisterAuthValidators(validate); err != nil {
		return fmt.Errorf("register auth validators: %w", err)
	}
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func Read(configData io.Reader) (*Config, error) {
	var config Config
	decoder := toml.NewDecoder(configData)
	if _, err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}
