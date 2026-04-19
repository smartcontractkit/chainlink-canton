package config

import (
	"fmt"
	"io"
	"reflect"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func DefaultConfig() *Config {
	return &Config{
		ChainSelector: "",
		Server:        ServerConfig{},
		Node:          NodeConfig{},
		Contracts:     Contracts{},
		Monitoring:    MonitoringConfig{},
	}
}

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
	DefaultExecutor       ContractIdentifier   `toml:"default_executor" validate:"required"`
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

type configTransformer struct{}

func (t configTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	switch typ {
	// Merge doesn't handle arrays gracefully, manually check if an InstanceAddress should override the dest
	case reflect.TypeFor[contracts.InstanceAddress]():
		return func(dst, src reflect.Value) error {
			if !src.IsZero() {
				dst.Set(src)
			}

			return nil
		}
	}

	return nil
}

// Merge merges the in config into cfg and returns cfg.
// Any unexported/unset/zero-value field will be ignored.
// If in is nil, no merge will be performed and cfg will be returned as-is.
func (cfg *Config) Merge(in *Config) (*Config, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if in == nil {
		return cfg, nil
	}

	if err := mergo.Merge(cfg, in, mergo.WithOverride, mergo.WithTransformers(configTransformer{})); err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	return cfg, nil
}

func Read(configData io.Reader) (*Config, error) {
	config := DefaultConfig()

	decoder := toml.NewDecoder(configData)
	if _, err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

// ReadAndMerge merges multiple configs in low->high preceding order.
// Later passed readers will take precedence over earlier ones.
func ReadAndMerge(configDatas ...io.Reader) (*Config, error) {
	config := DefaultConfig()

	for i, data := range configDatas {
		inCfg, err := Read(data)
		if err != nil {
			return nil, fmt.Errorf("failed to read config at index %d: %w", i, err)
		}
		config, err = config.Merge(inCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config at index %d: %w", i, err)
		}
	}

	return config, nil
}
