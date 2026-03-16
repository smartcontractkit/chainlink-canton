package changesets

import (
	"fmt"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rate_limiter"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type DeployRateLimiterConfig struct {
	PoolOwner           string
	PoolInstanceID      string
	RemoteChainSelector string
	Direction           common.RateLimitDirection
	Mode                common.RateLimitMode
	InstanceID          string
	Qualifier           string
	IsEnabled           bool
	Capacity            string
	Rate                string
	Tokens              string
}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployRateLimiterConfig]] = DeployRateLimiter{}

type DeployRateLimiter struct{}

func (d DeployRateLimiter) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployRateLimiterConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if len(chain.Participants) < config.Participant {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	return nil
}

func (d DeployRateLimiter) Apply(e cldf.Environment, config CantonCSDeps[DeployRateLimiterConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	deps := dependencies.CantonDeps{Chain: chain}
	cfg := config.Config

	if cfg.PoolOwner == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("pool owner is required")
	}
	if cfg.PoolInstanceID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("pool instance ID is required")
	}
	if cfg.RemoteChainSelector == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("remote chain selector is required")
	}
	if cfg.InstanceID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("instance ID is required")
	}

	mode := cfg.Mode
	if mode == "" {
		mode = common.RateLimitModeRateLimitMode_DefaultFinality
	}
	capacity := cfg.Capacity
	if capacity == "" {
		capacity = "0"
	}
	rate := cfg.Rate
	if rate == "" {
		rate = "0"
	}
	tokens := cfg.Tokens
	if tokens == "" {
		tokens = "0"
	}

	template := common.RateLimiter{
		InstanceId:          types.TEXT(cfg.InstanceID),
		PoolInstanceId:      types.TEXT(cfg.PoolInstanceID),
		PoolOwner:           types.PARTY(cfg.PoolOwner),
		RemoteChainSelector: types.NUMERIC(cfg.RemoteChainSelector),
		Direction:           cfg.Direction,
		Mode:                mode,
		IsEnabled:           types.BOOL(cfg.IsEnabled),
		Capacity:            types.NUMERIC(capacity),
		Rate:                types.NUMERIC(rate),
		Tokens:              types.NUMERIC(tokens),
		LastUpdated:         types.TIMESTAMP(time.Now()),
	}

	qualifier := ptr.String(cfg.Qualifier)
	if cfg.Qualifier == "" {
		qualifier = nil
	}

	input := contract.DeployInput[common.RateLimiter]{
		ChainSelector: config.ChainSelector,
		Qualifier:     qualifier,
		ActAs:         []string{cfg.PoolOwner},
		Template:      template,
		OwnerParty:    types.PARTY(cfg.PoolOwner),
	}

	switch cfg.Direction {
	case common.RateLimitDirectionRateLimitDirection_Inbound:
		out, err := cld_ops.ExecuteOperation(e.OperationsBundle, rate_limiter.DeployInbound, deps, input)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy rate limiter: %w", err)
		}
		if err = ds.AddressRefStore.Add(out.Output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save rate limiter address ref: %w", err)
		}
		// Add a generic type alias so devenv lookups can use one contract type regardless of direction.
		if err = ds.AddressRefStore.Upsert(datastore.AddressRef{
			Address:       out.Output.Address,
			Labels:        out.Output.Labels,
			Type:          datastore.ContractType("RateLimiter"),
			Version:       out.Output.Version,
			Qualifier:     out.Output.Qualifier,
			ChainSelector: out.Output.ChainSelector,
		}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save generic rate limiter alias address ref: %w", err)
		}
	case common.RateLimitDirectionRateLimitDirection_Outbound:
		out, err := cld_ops.ExecuteOperation(e.OperationsBundle, rate_limiter.DeployOutbound, deps, input)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy rate limiter: %w", err)
		}
		if err = ds.AddressRefStore.Add(out.Output); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save rate limiter address ref: %w", err)
		}
		// Add a generic type alias so devenv lookups can use one contract type regardless of direction.
		if err = ds.AddressRefStore.Upsert(datastore.AddressRef{
			Address:       out.Output.Address,
			Labels:        out.Output.Labels,
			Type:          datastore.ContractType("RateLimiter"),
			Version:       out.Output.Version,
			Qualifier:     out.Output.Qualifier,
			ChainSelector: out.Output.ChainSelector,
		}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save generic rate limiter alias address ref: %w", err)
		}
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("unsupported rate limiter direction %q", cfg.Direction)
	}

	return cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []cld_ops.Report[any, any]{},
	}, nil
}
