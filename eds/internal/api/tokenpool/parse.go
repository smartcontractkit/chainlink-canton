package tokenpool

import (
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/parse"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
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
	CCIPOwner          types.PARTY
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

	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(boundContract.RemoteChainConfigs))
	for chainSelectorString, remoteChainConfig := range boundContract.RemoteChainConfigs {
		chainSelector, err := internalparse.Uint64Checked(string(chainSelectorString))
		if err != nil {
			return nil, fmt.Errorf("invalid chain selector %q: %w", chainSelectorString, err)
		}

		// Rate Limiters
		inboundRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.InboundRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound rate limiter for remote chain selector %d: %w", chainSelector, err)
		}
		inboundCustomBlockConfirmationsRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound custom block confirmations rate limiter for remote chain selector %d: %w", chainSelector, err)
		}
		outboundRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.OutboundRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid outbound rate limiter for remote chain selector %d: %w", chainSelector, err)
		}

		// CCVs
		inboundCCVs, err := parse.RawInstanceAddressList(remoteChainConfig.InboundCCVs)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound CCVs for remote chain selector %d: %w", chainSelector, err)
		}
		outboundCCVs, err := parse.RawInstanceAddressList(remoteChainConfig.OutboundCCVs)
		if err != nil {
			return nil, fmt.Errorf("invalid outbound CCVs for remote chain selector %d: %w", chainSelector, err)
		}

		remoteChainConfigs[chainSelector] = RemoteChainConfig{
			InboundRateLimiter:                         inboundRateLimiter.InstanceAddress(),
			InboundCustomBlockConfirmationsRateLimiter: inboundCustomBlockConfirmationsRateLimiter.InstanceAddress(),
			OutboundRateLimiter:                        outboundRateLimiter.InstanceAddress(),
			InboundCCVs:                                inboundCCVs,
			OutboundCCVs:                               outboundCCVs,
		}
	}

	return &LockReleaseTokenPool{
		Address:            address,
		CCIPOwner:          boundContract.CcipOwner,
		RemoteChainConfigs: remoteChainConfigs,
		InstrumentId:       boundContract.InstrumentId,
		Decimals:           boundContract.Decimals,
	}, nil
}

type BurnMintTokenPool struct {
	Address            contracts.RawInstanceAddress
	PoolOwner          types.PARTY
	InstrumentId       splice_api_token_holding_v1.InstrumentId
	Decimals           types.INT64
	RemoteChainConfigs map[uint64]RemoteChainConfig
}

func ParseBurnMintTokenPool(createdEvent *apiv2.CreatedEvent) (*BurnMintTokenPool, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[burnminttokenpool.BurnMintTokenPool](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal burn/mint token pool: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.PoolOwner)

	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(boundContract.RemoteChainConfigs))
	for chainSelectorString, remoteChainConfig := range boundContract.RemoteChainConfigs {
		chainSelector, err := internalparse.Uint64Checked(string(chainSelectorString))
		if err != nil {
			return nil, fmt.Errorf("invalid chain selector %q: %w", chainSelectorString, err)
		}

		// Rate Limiters
		inboundRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.InboundRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound rate limiter for remote chain selector %d: %w", chainSelector, err)
		}
		inboundCustomBlockConfirmationsRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound custom block confirmations rate limiter for remote chain selector %d: %w", chainSelector, err)
		}
		outboundRateLimiter, err := parse.RawInstanceAddress(remoteChainConfig.OutboundRateLimiter)
		if err != nil {
			return nil, fmt.Errorf("invalid outbound rate limiter for remote chain selector %d: %w", chainSelector, err)
		}

		// CCVs
		inboundCCVs, err := parse.RawInstanceAddressList(remoteChainConfig.InboundCCVs)
		if err != nil {
			return nil, fmt.Errorf("invalid inbound CCVs for remote chain selector %d: %w", chainSelector, err)
		}
		outboundCCVs, err := parse.RawInstanceAddressList(remoteChainConfig.OutboundCCVs)
		if err != nil {
			return nil, fmt.Errorf("invalid outbound CCVs for remote chain selector %d: %w", chainSelector, err)
		}

		remoteChainConfigs[chainSelector] = RemoteChainConfig{
			InboundRateLimiter:                         inboundRateLimiter.InstanceAddress(),
			InboundCustomBlockConfirmationsRateLimiter: inboundCustomBlockConfirmationsRateLimiter.InstanceAddress(),
			OutboundRateLimiter:                        outboundRateLimiter.InstanceAddress(),
			InboundCCVs:                                inboundCCVs,
			OutboundCCVs:                               outboundCCVs,
		}
	}

	return &BurnMintTokenPool{
		Address:            address,
		PoolOwner:          boundContract.PoolOwner,
		RemoteChainConfigs: remoteChainConfigs,
		InstrumentId:       boundContract.InstrumentId,
		Decimals:           boundContract.Decimals,
	}, nil
}
