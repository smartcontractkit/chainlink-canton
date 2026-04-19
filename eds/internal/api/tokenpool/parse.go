package tokenpool

import (
	"fmt"
	"math/big"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/service/ledger"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
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
	RemoteChainConfigs map[uint64]RemoteChainConfig
}

func ParseLockReleaseTokenPool(createdEvent *apiv2.CreatedEvent) (*LockReleaseTokenPool, error) {
	boundContract, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](createdEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal lock release token pool: %w", err)
	}

	address := contracts.NewRawInstanceAddress(contracts.InstanceID(boundContract.InstanceId), boundContract.PoolOwner)

	remoteChainConfigs := make(map[uint64]RemoteChainConfig, len(boundContract.RemoteChainConfigs))
	for chainSelectorString, remoteChainConfigAny := range boundContract.RemoteChainConfigs {
		parsedSelector, ok := new(big.Rat).SetString(chainSelectorString)
		if !ok {
			return nil, fmt.Errorf("failed to parse chain selector %q", chainSelectorString)
		}
		if !parsedSelector.IsInt() {
			return nil, fmt.Errorf("chain selector %q is not an integer", chainSelectorString)
		}
		chainSelector := parsedSelector.Num().Uint64()

		remoteChainConfigMap, ok := remoteChainConfigAny.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("remote chain config for chain selector %d is not a map", chainSelector)
		}
		// unmarshal the remote chain config using ledger.MapToStruct
		var remoteChainConfig lockreleasetokenpool.RemoteChainConfig
		err = ledger.MapToStruct(remoteChainConfigMap, &remoteChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal remote chain config: %w", err)
		}

		// Rate Limiters
		inboundRateLimiter, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.InboundRateLimiter.Unpack))
		if err != nil {
			return nil, fmt.Errorf("invalid inbound rate limiter for remote chain selector %q: %w", chainSelector, err)
		}
		inboundCustomBlockConfirmationsRateLimiter, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.InboundCustomBlockConfirmationsRateLimiter.Unpack))
		if err != nil {
			return nil, fmt.Errorf("invalid inbound custom block confirmations rate limiter for remote chain selector %q: %w", chainSelector, err)
		}
		outboundRateLimiter, err := contracts.RawInstanceAddressFromString(string(remoteChainConfig.OutboundRateLimiter.Unpack))
		if err != nil {
			return nil, fmt.Errorf("invalid outbound rate limiter for remote chain selector %q: %w", chainSelector, err)
		}

		// CCVs
		inboundCCVs := make([]contracts.RawInstanceAddress, len(remoteChainConfig.InboundCCVs))
		for i, inboundCCV := range remoteChainConfig.InboundCCVs {
			inboundCCVAddress, err := contracts.RawInstanceAddressFromString(string(inboundCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid inbound CCV at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			inboundCCVs[i] = inboundCCVAddress
		}
		outboundCCVs := make([]contracts.RawInstanceAddress, len(remoteChainConfig.OutboundCCVs))
		for i, outboundCCV := range remoteChainConfig.OutboundCCVs {
			outboundCCVAddress, err := contracts.RawInstanceAddressFromString(string(outboundCCV.Unpack))
			if err != nil {
				return nil, fmt.Errorf("invalid outbound CCV at index %d for remote chain selector %q: %w", i, chainSelector, err)
			}
			outboundCCVs[i] = outboundCCVAddress
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
		RemoteChainConfigs: remoteChainConfigs,
		InstrumentId:       boundContract.InstrumentId,
	}, nil
}
