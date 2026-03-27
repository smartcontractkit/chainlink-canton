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
	ChainSelector string           `toml:"chain_selector" validate:"required"`
	Server        ServerConfig     `toml:"server" validate:"required"`
	Node          NodeConfig       `toml:"node" validate:"required"`
	Contracts     Contracts        `toml:"contracts" validate:"required"`
	Monitoring    MonitoringConfig `toml:"monitoring"`
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

type TokenPoolContracts struct {
	TokenPool                                  ContractIdentifier `toml:"token_pool" validate:"required"`
	InboundRateLimiter                         ContractIdentifier `toml:"inbound_rate_limiter" validate:"required"`
	InboundCustomBlockConfirmationsRateLimiter ContractIdentifier `toml:"inbound_custom_block_confirmations_rate_limiter" validate:"required"`
	OutboundRateLimiter                        ContractIdentifier `toml:"outbound_rate_limiter" validate:"required"`
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
	TokenPoolContracts    []TokenPoolContracts `toml:"token_pool_contracts"`

	// PoolOwner is the party that owns the token pools.
	// The instrument holdings of this owner is what will end up getting tracked by EDS,
	// and we will serve disclosures for them if needed (i.e. for token transfer executions).
	PoolOwner string `toml:"pool_owner"`
}

type ContractIdentifier struct {
	PartyID         string                    `toml:"party_id" validate:"required"`
	InstanceAddress contracts.InstanceAddress `toml:"instance_address" validate:"required"`
}

func (cfg *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
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
