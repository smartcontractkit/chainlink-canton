package config

import (
	"fmt"
	"io"
	"maps"
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
		Server: ServerConfig{
			MaxRequestSizeBytes: 1024 * 1024 * 10, // 10MiB
		},
		Node:       NodeConfig{},
		Monitoring: MonitoringConfig{},
		GlobalAPIConfig: GlobalAPIConfig{
			MaxBatchSize: 1024,
		},
		CCIPAPIConfig:      CCIPAPIConfig{},
		CCVAPIConfig:       CCVAPIConfig{},
		ExecutorAPIConfig:  ExecutorAPIConfig{},
		TokenPoolAPIConfig: TokenPoolAPIConfig{},
	}
}

type Config struct {
	ChainSelector string           `toml:"chain_selector" validate:"required"`
	Server        ServerConfig     `toml:"server" validate:"required"`
	Node          NodeConfig       `toml:"node" validate:"required"`
	Monitoring    MonitoringConfig `toml:"monitoring"`

	// API Configs
	GlobalAPIConfig        GlobalAPIConfig        `toml:"global_api"`
	CCIPAPIConfig          CCIPAPIConfig          `toml:"ccip_api"`
	CCVAPIConfig           CCVAPIConfig           `toml:"ccv_api"`
	ExecutorAPIConfig      ExecutorAPIConfig      `toml:"executor_api"`
	TokenPoolAPIConfig     TokenPoolAPIConfig     `toml:"token_pool_api"`
	TokenStandardAPIConfig TokenStandardAPIConfig `toml:"token_standard_api"`
}

type ServerConfig struct {
	Host                string `toml:"host" validate:"required"`
	Port                uint16 `toml:"port" validate:"required,port"`
	MaxRequestSizeBytes int64  `toml:"max_request_size_bytes" validate:"required"`
}

type NodeConfig struct {
	URL        string                  `toml:"url" validate:"required,url"`
	AuthConfig commonconfig.AuthConfig `toml:"auth" validate:"required"`
	MaxRetries int                     `toml:"max_retries"`
}

// Global API

type GlobalAPIConfig struct {
	// The maximum number of disclosures that can be requested in a single batch request.
	MaxBatchSize int `toml:"max_batch_size" validate:"required"`
}

// CCIP API

type CCIPAPIConfig struct {
	Enabled bool `toml:"enabled"`

	PerPartyRouterFactory ContractIdentifier `toml:"per_party_router_factory" validate:"required_if=Enabled true,omitempty"`
	OnRamp                ContractIdentifier `toml:"on_ramp" validate:"required_if=Enabled true,omitempty"`
	OffRamp               ContractIdentifier `toml:"off_ramp" validate:"required_if=Enabled true,omitempty"`
	GlobalConfig          ContractIdentifier `toml:"global_config" validate:"required_if=Enabled true,omitempty"`
	TokenAdminRegistry    ContractIdentifier `toml:"token_admin_registry" validate:"required_if=Enabled true,omitempty"`
	RMNRemote             ContractIdentifier `toml:"rmn_remote" validate:"required_if=Enabled true,omitempty"`
	FeeQuoter             ContractIdentifier `toml:"fee_quoter" validate:"required_if=Enabled true,omitempty"`
}

// CCV API

type CCV struct {
	ContractIdentifier
}

type CCVAPIConfig struct {
	Enabled bool  `toml:"enabled"`
	CCVs    []CCV `toml:"ccvs" validate:"required_if=Enabled true,dive"`
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

type FactoryType string

const (
	FactoryTypeDisabled FactoryType = ""
	FactoryTypeAddress  FactoryType = "address"
	FactoryTypeURL      FactoryType = "url"
)

type TransferFactory struct {
	Type FactoryType `toml:"type" validate:"oneof='' address url"`

	TemplateId      *string                    `toml:"template_id" validate:"required_if=Type address"`
	Party           *string                    `toml:"party" validate:"required_if=Type address"`
	InstanceAddress *contracts.InstanceAddress `toml:"instance_address" validate:"required_if=Type address"`

	TokenStandardURL        *string                  `toml:"token_standard_url" validate:"excluded_unless=Type url,required_if=Type url,omitnil,url"`
	TokenStandardAuthConfig *commonconfig.AuthConfig `toml:"token_standard_auth" validate:"excluded_unless=Type url"`
}

type BurnMintFactory struct {
	Type FactoryType `toml:"type" validate:"oneof='' address url"`

	TemplateId      *string                    `toml:"template_id" validate:"required_if=Type address"`
	Party           *string                    `toml:"party" validate:"required_if=Type address"`
	InstanceAddress *contracts.InstanceAddress `toml:"instance_address" validate:"required_if=Type address"`

	TokenStandardURL        *string                  `toml:"token_standard_url" validate:"excluded_unless=Type url,required_if=Type url,omitnil,url"`
	TokenStandardAuthConfig *commonconfig.AuthConfig `toml:"token_standard_auth" validate:"excluded_unless=Type url"`
}

type TokenPool struct {
	ContractIdentifier
	Type TokenPoolType `toml:"type" validate:"required,oneof=lockRelease burnMint"`

	// The owner party of the token pool.
	PoolOwner string `toml:"pool_owner" validate:"required"`

	TransferFactory     *TransferFactory     `toml:"transfer_factory" validate:"excluded_unless=Type lockRelease"`
	BurnMintFactory     *BurnMintFactory     `toml:"burn_mint_factory" validate:"excluded_unless=Type burnMint"`
	TransferPreapproval *TransferPreapproval `toml:"transfer_preapproval" validate:"omitnil"`
}

type TransferPreapproval struct {
	ContextKey string `toml:"context_key" validate:"required"`
	TemplateId string `toml:"template_id" validate:"required"`
}

type TokenPoolAPIConfig struct {
	Enabled bool `toml:"enabled"`
	// TokenPools is keyed by instance_address (contracts.InstanceAddress.Hex()) so layered configs merge per pool.
	TokenPools map[string]TokenPool `toml:"token_pools" validate:"required_if=Enabled true,dive"`
}

// Token Standard API

type TokenType string

const (
	TokenTypeLINK TokenType = "LINK"
)

type Registry struct {
	ContractIdentifier
	TokenType TokenType `toml:"token_type" validate:"required,oneof=LINK"`
	TokenId   string    `toml:"token_id" validate:"required"`
}

type TokenStandardAPIConfig struct {
	Enabled    bool                `toml:"enabled"`
	Admin      string              `toml:"admin" validate:"required_if=Enabled true"`
	Registries map[string]Registry `toml:"registries" validate:"required_if=Enabled true,dive"`
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
	if err := cfg.validateTokenPoolMapKeys(); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func (cfg *Config) validateTokenPoolMapKeys() error {
	for key, pool := range cfg.TokenPoolAPIConfig.TokenPools {
		want := pool.InstanceAddress.Hex()
		if key != want {
			return fmt.Errorf("token_pool_api.token_pools: map key %q must equal instance_address %q", key, want)
		}
	}

	return nil
}

type configTransformer struct{}

func (t configTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	switch typ {
	// Merge doesn't handle arrays gracefully, manually check if an InstanceAddress should override the dest
	case reflect.TypeFor[contracts.InstanceAddress]():
		return func(dst, src reflect.Value) error {
			if !src.IsZero() && dst.CanSet() {
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
//
// Token pool entries are merged separately by contract-address key: mergo cannot merge struct values
// inside maps correctly, so token_pool_api.token_pools uses a manual per-key merge.
func (cfg *Config) Merge(in *Config) (*Config, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if in == nil {
		return cfg, nil
	}

	basePools := maps.Clone(cfg.TokenPoolAPIConfig.TokenPools)
	overlayPools := maps.Clone(in.TokenPoolAPIConfig.TokenPools)
	cfg.TokenPoolAPIConfig.TokenPools = nil

	inMerged := *in
	inMerged.TokenPoolAPIConfig.TokenPools = nil

	if err := mergo.Merge(cfg, &inMerged, mergo.WithOverride, mergo.WithTransformers(configTransformer{})); err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	mergedPools, err := mergeTokenPoolMaps(basePools, overlayPools)
	if err != nil {
		return nil, err
	}
	if len(mergedPools) == 0 {
		mergedPools = nil
	}
	cfg.TokenPoolAPIConfig.TokenPools = mergedPools

	return cfg, nil
}

func mergeTokenPoolMaps(base, overlay map[string]TokenPool) (map[string]TokenPool, error) {
	out := make(map[string]TokenPool)
	maps.Copy(out, base)
	for k, sv := range overlay {
		if dv, ok := out[k]; ok {
			merged := dv
			if err := mergo.Merge(&merged, &sv, mergo.WithOverride, mergo.WithTransformers(configTransformer{})); err != nil {
				return nil, fmt.Errorf("merge token pool %s: %w", k, err)
			}
			out[k] = merged
		} else {
			out[k] = sv
		}
	}

	return out, nil
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
