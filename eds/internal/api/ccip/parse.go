package ccip

import (
	"fmt"
	"math/big"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
)

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
		return nil, fmt.Errorf("failed to unmarshal global config: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.CcipOwner)

	sourceChainConfigs := make(map[uint64]SourceChainConfig, len(boundContract.SourceChainConfigs))
	for chainSelectorString, sourceChainConfigAny := range boundContract.SourceChainConfigs {
		parsedSelector, ok := new(big.Rat).SetString(chainSelectorString)
		if !ok {
			return nil, fmt.Errorf("failed to parse chain selector %q", chainSelectorString)
		}
		if !parsedSelector.IsInt() {
			return nil, fmt.Errorf("chain selector %q is not an integer", chainSelectorString)
		}
		chainSelector := parsedSelector.Num().Uint64()

		sourceChainConfigMap, ok := sourceChainConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("source chain config for chain selector %d is not a map", chainSelector)
		}
		var sourceChainConfig common.SourceChainConfig
		err = ledger.MapToStruct(sourceChainConfigMap, &sourceChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal source chain config for selector %q: %w", chainSelector, err)
		}

		laneMandatedCCVs := make([]contracts.RawInstanceAddress, len(sourceChainConfig.LaneMandatedCCVs))
		for i, laneMandatedCCV := range sourceChainConfig.LaneMandatedCCVs {
			ccvAddress, err := contracts.RawInstanceAddressFromString(string(laneMandatedCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid lane mandated CCV at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			laneMandatedCCVs[i] = ccvAddress
		}
		defaultCCVs := make([]contracts.RawInstanceAddress, len(sourceChainConfig.DefaultCCVs))
		for i, defaultCCV := range sourceChainConfig.DefaultCCVs {
			ccvAddress, err := contracts.RawInstanceAddressFromString(string(defaultCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid default ccv at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			defaultCCVs[i] = ccvAddress
		}

		sourceChainConfigs[chainSelector] = SourceChainConfig{
			IsEnabled:        bool(sourceChainConfig.IsEnabled),
			LaneMandatedCCVs: laneMandatedCCVs,
			DefaultCCVs:      defaultCCVs,
		}
	}

	destChainConfigs := make(map[uint64]DestChainConfig, len(boundContract.DestChainConfigs))
	for chainSelectorString, destChainConfigAny := range boundContract.DestChainConfigs {
		parsedSelector, ok := new(big.Rat).SetString(chainSelectorString)
		if !ok {
			return nil, fmt.Errorf("failed to parse chain selector %q", chainSelectorString)
		}
		if !parsedSelector.IsInt() {
			return nil, fmt.Errorf("chain selector %q is not an integer", chainSelectorString)
		}
		chainSelector := parsedSelector.Num().Uint64()

		destChainConfigMap, ok := destChainConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("dest chain config for chain selector %d is not a map", chainSelector)
		}
		var destChainConfig common.DestChainConfig
		err = ledger.MapToStruct(destChainConfigMap, &destChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dest chain config for selector %q: %w", chainSelector, err)
		}

		laneMandatedCCVs := make([]contracts.RawInstanceAddress, len(destChainConfig.LaneMandatedCCVs))
		for i, laneMandatedCCV := range destChainConfig.LaneMandatedCCVs {
			ccvAddress, err := contracts.RawInstanceAddressFromString(string(laneMandatedCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid lane mandated CCV at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			laneMandatedCCVs[i] = ccvAddress
		}
		defaultCCVs := make([]contracts.RawInstanceAddress, len(destChainConfig.DefaultCCVs))
		for i, defaultCCV := range destChainConfig.DefaultCCVs {
			ccvAddress, err := contracts.RawInstanceAddressFromString(string(defaultCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid default ccv at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			defaultCCVs[i] = ccvAddress
		}

		var defaultExecutor *contracts.RawInstanceAddress
		if destChainConfig.DefaultExecutor != nil {
			executor, err := contracts.RawInstanceAddressFromString(string((*destChainConfig.DefaultExecutor).Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid default executor for chain selector %q: %w", chainSelector, err)
			}
			defaultExecutor = &executor
		}

		destChainConfigs[chainSelector] = DestChainConfig{
			IsEnabled:        bool(destChainConfig.IsEnabled),
			DefaultExecutor:  defaultExecutor,
			LaneMandatedCCVs: laneMandatedCCVs,
			DefaultCCVs:      defaultCCVs,
		}
	}

	return &GlobalConfig{
		Address:            address,
		SourceChainConfigs: sourceChainConfigs,
		DestChainConfigs:   destChainConfigs,
	}, nil
}

type TokenAdminRegistry struct {
	Address      contracts.RawInstanceAddress
	TokenConfigs map[contracts.EncodedInstrumentID]tokenadminregistry.TokenConfig
}

func ParseTokenAdminRegistry(createdEvent *apiv2.CreatedEvent) (*TokenAdminRegistry, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[tokenadminregistry.TokenAdminRegistry](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token admin registry: %w", err)
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
