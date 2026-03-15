package sequences

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	clccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type ConfigureTokenPoolRateLimitersInput struct {
	PoolInstanceAddress contracts.InstanceAddress
	PoolOwnerParty      string
	Config              ConfigureTokenPoolRateLimiterConfig
}

type ConfigureTokenPoolRateLimiterConfig struct {
	IsEnabled *bool
	Capacity  types.NUMERIC
	Rate      types.NUMERIC
	Tokens    types.NUMERIC
}

var ConfigureTokenPoolRateLimiters = operations.NewSequence(
	"canton/ccip/configure_token_pool_rate_limiters",
	semver.MustParse("1.0.0"),
	"Deploys rate limiters for a token pool and wires them into the pool config",
	configureTokenPoolRateLimiters,
)

func configureTokenPoolRateLimiters(b operations.Bundle, deps dependencies.CantonDeps, input ConfigureTokenPoolRateLimitersInput) (clccipsequences.OnChainOutput, error) {
	participant := deps.Chain.Participants[deps.Participant]
	activePool, err := contract.FindActiveContractByInstanceAddress(
		b.GetContext(),
		participant.LedgerServices.State,
		participant.PartyID,
		lockreleasetokenpool.LockReleaseTokenPool{}.GetTemplateID(),
		input.PoolInstanceAddress,
	)
	if err != nil {
		return clccipsequences.OnChainOutput{}, fmt.Errorf("find active token pool: %w", err)
	}
	parsedPool, err := bindings.UnmarshalCreatedEvent[lockreleasetokenpool.LockReleaseTokenPool](activePool.GetCreatedEvent())
	if err != nil {
		return clccipsequences.OnChainOutput{}, fmt.Errorf("parse active token pool: %w", err)
	}
	if len(parsedPool.ChainPoolConfigs) == 0 || len(parsedPool.OutboundRateLimiters) > 0 {
		return clccipsequences.OnChainOutput{}, nil
	}

	newOutbound := make(types.GENMAP)
	newInbound := make(types.GENMAP)
	newInboundCustom := make(types.GENMAP)
	now := types.TIMESTAMP(time.Now().UTC())
	cfg := input.Config
	if cfg.Capacity == "" {
		cfg.Capacity = types.NUMERIC("999999999999999999")
	}
	if cfg.Rate == "" {
		cfg.Rate = types.NUMERIC("999999999999999999")
	}
	if cfg.Tokens == "" {
		cfg.Tokens = types.NUMERIC("0")
	}
	isEnabled := types.BOOL(false)
	if cfg.IsEnabled != nil {
		isEnabled = types.BOOL(*cfg.IsEnabled)
	}

	for selectorKey := range parsedPool.ChainPoolConfigs {
		remoteChainSelector := types.NUMERIC(strconv.FormatInt(normalizeSelectorKey(selectorKey), 10))
		outbound, err := operations.ExecuteOperation(b, rate_limiter.Deploy, deps, contract.DeployInput[common.RateLimiter]{
			ChainSelector: deps.Chain.Selector,
			Qualifier:     nil,
			ActAs:         []string{input.PoolOwnerParty},
			Template: common.RateLimiter{
				PoolInstanceId:      parsedPool.InstanceId,
				PoolOwner:           parsedPool.PoolOwner,
				RemoteChainSelector: remoteChainSelector,
				Direction:           common.RateLimitDirectionRateLimitDirection_Outbound,
				Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
				IsEnabled:           isEnabled,
				Capacity:            cfg.Capacity,
				Rate:                cfg.Rate,
				Tokens:              cfg.Tokens,
				LastUpdated:         now,
			},
			OwnerParty: parsedPool.PoolOwner,
		})
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("deploy outbound rate limiter for selector %s: %w", selectorKey, err)
		}
		inbound, err := operations.ExecuteOperation(b, rate_limiter.Deploy, deps, contract.DeployInput[common.RateLimiter]{
			ChainSelector: deps.Chain.Selector,
			Qualifier:     nil,
			ActAs:         []string{input.PoolOwnerParty},
			Template: common.RateLimiter{
				PoolInstanceId:      parsedPool.InstanceId,
				PoolOwner:           parsedPool.PoolOwner,
				RemoteChainSelector: remoteChainSelector,
				Direction:           common.RateLimitDirectionRateLimitDirection_Inbound,
				Mode:                common.RateLimitModeRateLimitMode_DefaultFinality,
				IsEnabled:           isEnabled,
				Capacity:            cfg.Capacity,
				Rate:                cfg.Rate,
				Tokens:              cfg.Tokens,
				LastUpdated:         now,
			},
			OwnerParty: parsedPool.PoolOwner,
		})
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("deploy inbound rate limiter for selector %s: %w", selectorKey, err)
		}
		inboundCustom, err := operations.ExecuteOperation(b, rate_limiter.Deploy, deps, contract.DeployInput[common.RateLimiter]{
			ChainSelector: deps.Chain.Selector,
			Qualifier:     nil,
			ActAs:         []string{input.PoolOwnerParty},
			Template: common.RateLimiter{
				PoolInstanceId:      parsedPool.InstanceId,
				PoolOwner:           parsedPool.PoolOwner,
				RemoteChainSelector: remoteChainSelector,
				Direction:           common.RateLimitDirectionRateLimitDirection_Inbound,
				Mode:                common.RateLimitModeRateLimitMode_CustomFinality,
				IsEnabled:           isEnabled,
				Capacity:            cfg.Capacity,
				Rate:                cfg.Rate,
				Tokens:              cfg.Tokens,
				LastUpdated:         now,
			},
			OwnerParty: parsedPool.PoolOwner,
		})
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("deploy inbound custom-finality rate limiter for selector %s: %w", selectorKey, err)
		}

		outboundLabels := outbound.Output.Labels.List()
		inboundLabels := inbound.Output.Labels.List()
		inboundCustomLabels := inboundCustom.Output.Labels.List()
		if len(outboundLabels) == 0 || len(inboundLabels) == 0 || len(inboundCustomLabels) == 0 {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("missing rate limiter raw address labels for selector %s", selectorKey)
		}

		outboundRaw, err := contracts.RawInstanceAddressFromString(outboundLabels[0])
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("parse outbound rate limiter raw address for selector %s: %w", selectorKey, err)
		}
		inboundRaw, err := contracts.RawInstanceAddressFromString(inboundLabels[0])
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("parse inbound rate limiter raw address for selector %s: %w", selectorKey, err)
		}
		inboundCustomRaw, err := contracts.RawInstanceAddressFromString(inboundCustomLabels[0])
		if err != nil {
			return clccipsequences.OnChainOutput{}, fmt.Errorf("parse inbound custom-finality rate limiter raw address for selector %s: %w", selectorKey, err)
		}

		newOutbound[selectorKey] = outboundRaw.Binding()
		newInbound[selectorKey] = inboundRaw.Binding()
		newInboundCustom[selectorKey] = inboundCustomRaw.Binding()
	}

	_, err = operations.ExecuteOperation(b, lock_release_token_pool.UpdateRateLimiters, deps, contract.ChoiceInput[lockreleasetokenpool.LockReleaseTokenPoolUpdateRateLimiters]{
		ChainSelector:   deps.Chain.Selector,
		InstanceAddress: input.PoolInstanceAddress,
		ActAs:           []string{input.PoolOwnerParty},
		Args: lockreleasetokenpool.LockReleaseTokenPoolUpdateRateLimiters{
			NewOutboundRateLimiters:      newOutbound,
			NewInboundRateLimiters:       newInbound,
			NewInboundCustomRateLimiters: newInboundCustom,
		},
	})
	if err != nil {
		return clccipsequences.OnChainOutput{}, fmt.Errorf("update token pool rate limiters: %w", err)
	}

	return clccipsequences.OnChainOutput{}, nil
}

func normalizeSelectorKey(selectorKey string) int64 {
	selector, err := strconv.ParseInt(strings.TrimSuffix(selectorKey, "."), 10, 64)
	if err == nil {
		return selector
	}
	return 0
}
