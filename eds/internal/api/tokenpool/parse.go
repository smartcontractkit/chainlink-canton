package tokenpool

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/service/ledger"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/parse"
)

type RemoteChainConfig struct {
	InboundRateLimiter                         contracts.InstanceAddress
	InboundCustomBlockConfirmationsRateLimiter contracts.InstanceAddress
	OutboundRateLimiter                        contracts.InstanceAddress

	InboundCCVs  []contracts.RawInstanceAddress
	OutboundCCVs []contracts.RawInstanceAddress
}

type LockReleaseTokenPool struct {
	Address            contracts.RawInstanceAddress
	InstrumentId       splice_api_token_holding_v1.InstrumentId
	Decimals           types.INT64
	RemoteChainConfigs map[uint64]RemoteChainConfig
}

type BurnMintTokenPool struct {
	Address            contracts.RawInstanceAddress
	InstrumentId       splice_api_token_holding_v1.InstrumentId
	Decimals           types.INT64
	RemoteChainConfigs map[uint64]RemoteChainConfig
}

func ParseLockReleaseTokenPool(createdEvent *apiv2.CreatedEvent) (*LockReleaseTokenPool, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal lock release token pool: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.PoolOwner)
	remoteChainConfigs, err := parseLockReleaseRemoteChainConfigs(boundContract.RemoteChainConfigs)
	if err != nil {
		return nil, err
	}

	return &LockReleaseTokenPool{
		Address:            address,
		RemoteChainConfigs: remoteChainConfigs,
		InstrumentId:       boundContract.InstrumentId,
		Decimals:           boundContract.Decimals,
	}, nil
}

func ParseBurnMintTokenPool(createdEvent *apiv2.CreatedEvent) (*BurnMintTokenPool, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[burnminttokenpool.BurnMintTokenPool](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal burn/mint token pool: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.PoolOwner)
	remoteChainConfigs, err := parseBurnMintRemoteChainConfigs(boundContract.RemoteChainConfigs)
	if err != nil {
		return nil, err
	}

	return &BurnMintTokenPool{
		Address:            address,
		RemoteChainConfigs: remoteChainConfigs,
		InstrumentId:       boundContract.InstrumentId,
		Decimals:           boundContract.Decimals,
	}, nil
}

func parseLockReleaseRemoteChainConfigs(remoteChainConfigsMap types.GENMAP) (map[uint64]RemoteChainConfig, error) {
	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(remoteChainConfigsMap))
	for chainSelectorString, remoteChainConfigAny := range remoteChainConfigsMap {
		chainSelector, err := parse.Uint64Checked(chainSelectorString)
		if err != nil {
			return nil, fmt.Errorf("invalid chain selector %q: %w", chainSelectorString, err)
		}

		remoteChainConfigMap, ok := remoteChainConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("remote chain config for chain selector %d is not a map", chainSelector)
		}
		var remoteChainConfig lockreleasetokenpool.RemoteChainConfig
		err = ledger.MapToStruct(remoteChainConfigMap, &remoteChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal lock release remote chain config: %w", err)
		}

		parsedConfig, err := parseRemoteChainConfig(
			chainSelector,
			remoteChainConfig.InboundRateLimiter,
			remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter,
			remoteChainConfig.OutboundRateLimiter,
			remoteChainConfig.InboundCCVs,
			remoteChainConfig.OutboundCCVs,
		)
		if err != nil {
			return nil, err
		}
		remoteChainConfigs[chainSelector] = parsedConfig
	}

	return remoteChainConfigs, nil
}

func parseBurnMintRemoteChainConfigs(remoteChainConfigsMap types.GENMAP) (map[uint64]RemoteChainConfig, error) {
	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(remoteChainConfigsMap))
	for chainSelectorString, remoteChainConfigAny := range remoteChainConfigsMap {
		chainSelector, err := parse.Uint64Checked(chainSelectorString)
		if err != nil {
			return nil, fmt.Errorf("invalid chain selector %q: %w", chainSelectorString, err)
		}

		remoteChainConfigMap, ok := remoteChainConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("remote chain config for chain selector %d is not a map", chainSelector)
		}
		var remoteChainConfig burnminttokenpool.RemoteChainConfig
		err = ledger.MapToStruct(remoteChainConfigMap, &remoteChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal burn/mint remote chain config: %w", err)
		}

		parsedConfig, err := parseRemoteChainConfig(
			chainSelector,
			remoteChainConfig.InboundRateLimiter,
			remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter,
			remoteChainConfig.OutboundRateLimiter,
			remoteChainConfig.InboundCCVs,
			remoteChainConfig.OutboundCCVs,
		)
		if err != nil {
			return nil, err
		}
		remoteChainConfigs[chainSelector] = parsedConfig
	}

	return remoteChainConfigs, nil
}

func parseRemoteChainConfig(
	chainSelector uint64,
	inboundRateLimiterRaw mcms.RawInstanceAddress,
	inboundCustomRateLimiterRaw mcms.RawInstanceAddress,
	outboundRateLimiterRaw mcms.RawInstanceAddress,
	inboundCCVsRaw []mcms.RawInstanceAddress,
	outboundCCVsRaw []mcms.RawInstanceAddress,
) (RemoteChainConfig, error) {
	inboundRateLimiter, err := parse.RawInstanceAddress(inboundRateLimiterRaw)
	if err != nil {
		return RemoteChainConfig{}, fmt.Errorf("invalid inbound rate limiter for remote chain selector %d: %w", chainSelector, err)
	}
	inboundCustomBlockConfirmationsRateLimiter, err := parse.RawInstanceAddress(inboundCustomRateLimiterRaw)
	if err != nil {
		return RemoteChainConfig{}, fmt.Errorf("invalid inbound custom block confirmations rate limiter for remote chain selector %d: %w", chainSelector, err)
	}
	outboundRateLimiter, err := parse.RawInstanceAddress(outboundRateLimiterRaw)
	if err != nil {
		return RemoteChainConfig{}, fmt.Errorf("invalid outbound rate limiter for remote chain selector %d: %w", chainSelector, err)
	}

	inboundCCVs, err := parse.RawInstanceAddressList(inboundCCVsRaw)
	if err != nil {
		return RemoteChainConfig{}, fmt.Errorf("invalid inbound CCVs for remote chain selector %d: %w", chainSelector, err)
	}
	outboundCCVs, err := parse.RawInstanceAddressList(outboundCCVsRaw)
	if err != nil {
		return RemoteChainConfig{}, fmt.Errorf("invalid outbound CCVs for remote chain selector %d: %w", chainSelector, err)
	}

	return RemoteChainConfig{
		InboundRateLimiter:                         inboundRateLimiter.InstanceAddress(),
		InboundCustomBlockConfirmationsRateLimiter: inboundCustomBlockConfirmationsRateLimiter.InstanceAddress(),
		OutboundRateLimiter:                        outboundRateLimiter.InstanceAddress(),
		InboundCCVs:                                inboundCCVs,
		OutboundCCVs:                               outboundCCVs,
	}, nil
}
