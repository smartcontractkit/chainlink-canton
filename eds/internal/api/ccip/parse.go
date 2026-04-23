package ccip

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/parse"
)

type PerPartyRouterFactory struct {
	Address contracts.RawInstanceAddress
}

func ParsePerPartyRouterFactory(createdEvent *apiv2.CreatedEvent) (*PerPartyRouterFactory, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[perpartyrouter.PerPartyRouterFactory](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PerPartyRouterFactory: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.CcipOwner)

	return &PerPartyRouterFactory{
		Address: address,
	}, nil
}

type SourceChainConfig struct {
	IsEnabled        bool
	LaneMandatedCCVs []contracts.RawInstanceAddress
	DefaultCCVs      []contracts.RawInstanceAddress
}

type DestChainConfig struct {
	IsEnabled        bool
	DefaultExecutor  *contracts.RawInstanceAddress
	LaneMandatedCCVs []contracts.RawInstanceAddress
	DefaultCCVs      []contracts.RawInstanceAddress
}

type GlobalConfig struct {
	Address            contracts.RawInstanceAddress
	SourceChainConfigs map[uint64]SourceChainConfig
	DestChainConfigs   map[uint64]DestChainConfig
}

func ParseGlobalConfig(createdEvent *apiv2.CreatedEvent) (*GlobalConfig, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[common.GlobalConfig](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal GlobalConfig: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.CcipOwner)

	sourceChainConfigs := make(map[uint64]SourceChainConfig, len(boundContract.SourceChainConfigs))
	for chainSelectorString, sourceChainConfigAny := range boundContract.SourceChainConfigs {
		chainSelector, err := parse.Uint64Checked(chainSelectorString)
		if err != nil {
			return nil, fmt.Errorf("invalid source chain config selector %q: %w", chainSelectorString, err)
		}

		sourceChainConfigs[chainSelector], err = parseSourceChainConfig(sourceChainConfigAny)
		if err != nil {
			return nil, fmt.Errorf("failed to parse source chain config for selector %d: %w", chainSelector, err)
		}
	}

	destChainConfigs := make(map[uint64]DestChainConfig, len(boundContract.DestChainConfigs))
	for chainSelectorString, destChainConfigAny := range boundContract.DestChainConfigs {
		chainSelector, err := parse.Uint64Checked(chainSelectorString)
		if err != nil {
			return nil, fmt.Errorf("invalid dest chain config selector %q: %w", chainSelectorString, err)
		}

		destChainConfigs[chainSelector], err = parseDestChainConfig(destChainConfigAny)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dest chain config for selector %d: %w", chainSelector, err)
		}
	}

	return &GlobalConfig{
		Address:            address,
		SourceChainConfigs: sourceChainConfigs,
		DestChainConfigs:   destChainConfigs,
	}, nil
}

func parseSourceChainConfig(cfg any) (SourceChainConfig, error) {
	sourceChainConfigMap, ok := cfg.(map[string]any)
	if !ok {
		return SourceChainConfig{}, fmt.Errorf("not a map")
	}
	var sourceChainConfig common.SourceChainConfig
	if err := ledger.MapToStruct(sourceChainConfigMap, &sourceChainConfig); err != nil {
		return SourceChainConfig{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	laneMandatedCCVs, err := parse.RawInstanceAddressList(sourceChainConfig.LaneMandatedCCVs)
	if err != nil {
		return SourceChainConfig{}, fmt.Errorf("failed to parse lane mandated CCVs: %w", err)
	}
	defaultCCVs, err := parse.RawInstanceAddressList(sourceChainConfig.DefaultCCVs)
	if err != nil {
		return SourceChainConfig{}, fmt.Errorf("failed to parse default CCVs: %w", err)
	}

	return SourceChainConfig{
		IsEnabled:        bool(sourceChainConfig.IsEnabled),
		LaneMandatedCCVs: laneMandatedCCVs,
		DefaultCCVs:      defaultCCVs,
	}, nil
}

func parseDestChainConfig(cfg any) (DestChainConfig, error) {
	destChainConfigMap, ok := cfg.(map[string]any)
	if !ok {
		return DestChainConfig{}, fmt.Errorf("not a map")
	}
	var destChainConfig common.DestChainConfig
	if err := ledger.MapToStruct(destChainConfigMap, &destChainConfig); err != nil {
		return DestChainConfig{}, fmt.Errorf("failed to unmarshal: %w", err)
	}

	laneMandatedCCVs, err := parse.RawInstanceAddressList(destChainConfig.LaneMandatedCCVs)
	if err != nil {
		return DestChainConfig{}, fmt.Errorf("failed to parse lane mandated CCVs: %w", err)
	}
	defaultCCVs, err := parse.RawInstanceAddressList(destChainConfig.DefaultCCVs)
	if err != nil {
		return DestChainConfig{}, fmt.Errorf("failed to parse default CCVs: %w", err)
	}

	var defaultExecutor *contracts.RawInstanceAddress
	if destChainConfig.DefaultExecutor != nil {
		executor, err := parse.RawInstanceAddress(*destChainConfig.DefaultExecutor)
		if err != nil {
			return DestChainConfig{}, fmt.Errorf("failed to parse default executor: %w", err)
		}
		defaultExecutor = &executor
	}

	return DestChainConfig{
		IsEnabled:        bool(destChainConfig.IsEnabled),
		DefaultExecutor:  defaultExecutor,
		LaneMandatedCCVs: laneMandatedCCVs,
		DefaultCCVs:      defaultCCVs,
	}, nil
}

type TokenAdminRegistry struct {
	Address      contracts.RawInstanceAddress
	TokenConfigs map[contracts.EncodedInstrumentID]tokenadminregistry.TokenConfig
}

func ParseTokenAdminRegistry(createdEvent *apiv2.CreatedEvent) (*TokenAdminRegistry, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenAdminRegistry](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal TokenAdminRegistry: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.Owner)

	tokenConfigs := make(map[contracts.EncodedInstrumentID]tokenadminregistry.TokenConfig, len(boundContract.TokenConfigs))
	for encodedInstrumentIdString, tokenConfigAny := range boundContract.TokenConfigs {
		encodedInstrumentId := contracts.HexToEncodedInstrumentID(encodedInstrumentIdString)

		tokenConfigMap, ok := tokenConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("token config for encoded instrument ID %s is not a map", encodedInstrumentIdString)
		}
		var tokenConfig tokenadminregistry.TokenConfig
		err = ledger.MapToStruct(tokenConfigMap, &tokenConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal token config for encoded instrument ID %s: %w", encodedInstrumentIdString, err)
		}

		tokenConfigs[encodedInstrumentId] = tokenConfig
	}

	return &TokenAdminRegistry{
		Address:      address,
		TokenConfigs: tokenConfigs,
	}, nil
}
