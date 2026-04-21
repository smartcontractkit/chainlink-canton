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
		Monitoring:    MonitoringConfig{},
		GlobalAPIConfig: GlobalAPIConfig{
			MaxBatchSize: 1024,
		},
		CCIPAPIConfig: CCIPAPIConfig{},
	}
}

type Config struct {
	ChainSelector string           `toml:"chain_selector" validate:"required"`
	Server        ServerConfig     `toml:"server" validate:"required"`
	Node          NodeConfig       `toml:"node" validate:"required"`
	Monitoring    MonitoringConfig `toml:"monitoring"`

	// API Configs
	GlobalAPIConfig    GlobalAPIConfig    `toml:"global_api"`
	CCIPAPIConfig      CCIPAPIConfig      `toml:"ccip_api"`
	CCVAPIConfig       CCVAPIConfig       `toml:"ccv_api"`
	ExecutorAPIConfig  ExecutorAPIConfig  `toml:"executor_api"`
	TokenPoolAPIConfig TokenPoolAPIConfig `toml:"token_pool_api"`
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

// Global API

type GlobalAPIConfig struct {
	MaxBatchSize int `toml:"max_batch_size" validate:"required"`
}

// CCIP API

type CCIPAPIConfig struct {
	Enabled bool `toml:"enabled"`

	PerPartyRouterFactory ContractIdentifier `toml:"per_party_router_factory" validate:"required_if=Enabled true"`
	OnRamp                ContractIdentifier `toml:"on_ramp" validate:"required_if=Enabled true"`
	OffRamp               ContractIdentifier `toml:"off_ramp" validate:"required_if=Enabled true"`
	GlobalConfig          ContractIdentifier `toml:"global_config" validate:"required_if=Enabled true"`
	TokenAdminRegistry    ContractIdentifier `toml:"token_admin_registry" validate:"required_if=Enabled true"`
	RMNRemote             ContractIdentifier `toml:"rmn_remote" validate:"required_if=Enabled true"`
	FeeQuoter             ContractIdentifier `toml:"fee_quoter" validate:"required_if=Enabled true"`
}

// CCV API

type CCV struct {
	ContractIdentifier
}

type CCVAPIConfig struct {
	Enabled bool  `toml:"enabled"`
	CCVs    []CCV `toml:"ccvs" validate:"required,dive"`
}

// Executor API

type Executor struct {
	ContractIdentifier
}

type ExecutorAPIConfig struct {
	Enabled   bool       `toml:"enabled"`
	Executors []Executor `toml:"executors" validate:"required_if=Enabled true,dive"`
}

// Token Pool API

type TokenPoolType string

const (
	TokenPoolTypeLockRelease TokenPoolType = "lockRelease"
	TokenPoolTypeBurnMint    TokenPoolType = "burnMint"
)

type TokenPool struct {
	ContractIdentifier
	Type TokenPoolType `toml:"type" validate:"required,oneof=lockRelease burnMint"`

	// The owner party of the token pool that is the owner of locked holdings
	PoolOwner string `toml:"pool_owner" validate:"required"`
	// The URL of the Token Standard API to use for this token pool.
	// If not set, fetching the transfer factory will be disabled for this pool.
	TokenStandardURL        *string                  `toml:"token_standard_url" validate:"omitnil,url"`
	TokenStandardAuthConfig *commonconfig.AuthConfig `toml:"token_standard_auth" validate:"omitnil,required"`
}

type TokenPoolAPIConfig struct {
	Enabled    bool        `toml:"enabled"`
	TokenPools []TokenPool `toml:"token_pools" validate:"required_if=Enabled true,dive"`
}

// ContractIdentifier uniquely identifies a contract using an InstanceAddress.
// It also contains a PartyID of a party that must be a stakeholder on the contract.
// This party will be used to retrieve contracts from the Active Contract Set, it can be the same party that owns
// the contract, but does not have to be.
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
