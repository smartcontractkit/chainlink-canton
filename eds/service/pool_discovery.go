package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/registry/ratelimiter"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
)

// PoolDiscoveryService watches for BurnMintTokenPool and LockReleaseTokenPool CreatedEvents
// naming EDS's observer party as an observer, and dynamically registers them for serving.
type PoolDiscoveryService struct {
	logger              zerolog.Logger
	activeContractStore store.ActiveContractStoreInterface
	tokenPoolServer     *tokenpool.Server
	observerParty       types.PARTY
	tokenStandardURL    string
	// discoveredPools tracks the ledger offset of the contract currently registered for each
	// address. InstanceIds aren't guaranteed unique, so two distinct contracts can derive the
	// same address; the offset lets a later-created contract always win, regardless of the
	// (effectively random) order poll results are processed in.
	discoveredPools map[contracts.InstanceAddress]int64
	mux             sync.Mutex
}

// NewPoolDiscoveryService creates a new pool discovery service.
// It must be created BEFORE activeContractStore.Run() is called.
// cfg.party_id is EDS's own party, named as an observer on discoverable pool contracts.
// cfg.token_standard_url is used to wire a URL-mode factory for each discovered pool,
// so send/execute disclosures include the token's burn-mint/transfer factory without
// any per-pool static config.
func NewPoolDiscoveryService(
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStoreInterface,
	tokenPoolServer *tokenpool.Server,
	cfg config.RegistryAPIConfig,
) *PoolDiscoveryService {
	logger = logger.With().Str("component", "PoolDiscoveryService").Logger()

	// Pre-register both pool templates to watch for CreatedEvents, plus RateLimiter so a
	// discovered pool's rate limiters resolve when serving send/execute. RegisterDiscoveredPool
	// cannot register these later, as RegisterTemplates must run before activeContractStore.Run.
	//
	// RateLimiter has `observer observers`, so this resolves as long as the observer party is
	// named in the RateLimiter's own observers field (not just the pool's).
	activeContractStore.RegisterTemplates(
		store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{}),
			PartyID:    cfg.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
			PartyID:    cfg.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(ratelimiter.RateLimiter{}),
			PartyID:    cfg.PartyID,
		},
	)

	return &PoolDiscoveryService{
		logger:              logger,
		activeContractStore: activeContractStore,
		tokenPoolServer:     tokenPoolServer,
		observerParty:       types.PARTY(cfg.PartyID),
		tokenStandardURL:    cfg.TokenStandardURL,
		discoveredPools:     make(map[contracts.InstanceAddress]int64),
	}
}

// Watch starts watching for new token pool contracts.
func (s *PoolDiscoveryService) Watch(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.checkForNewPools(ctx)
		}
	}
}

// checkForNewPools polls the activeContractStore for BurnMintTokenPool and LockReleaseTokenPool contracts.
func (s *PoolDiscoveryService) checkForNewPools(ctx context.Context) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.checkBurnMintPools(ctx)
	s.checkLockReleasePools(ctx)
}

func (s *PoolDiscoveryService) checkBurnMintPools(ctx context.Context) {
	templates, ok := s.activeContractStore.GetByTemplateId(s.observerParty, contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{}))
	if !ok {
		return // No pools yet
	}

	for _, activeContract := range templates {
		pool, err := tokenpool.ParseBurnMintTokenPool(activeContract.GetCreatedEvent())
		if err != nil {
			s.logger.Err(err).Msg("failed to parse burn mint token pool")
			continue
		}
		s.registerIfNewer(ctx, pool.Address, activeContract.GetCreatedEvent().GetOffset(), config.TokenPoolTypeBurnMint)
	}
}

func (s *PoolDiscoveryService) checkLockReleasePools(ctx context.Context) {
	templates, ok := s.activeContractStore.GetByTemplateId(s.observerParty, contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}))
	if !ok {
		return // No pools yet
	}

	for _, activeContract := range templates {
		pool, err := tokenpool.ParseLockReleaseTokenPool(activeContract.GetCreatedEvent())
		if err != nil {
			s.logger.Err(err).Msg("failed to parse lock release token pool")
			continue
		}
		s.registerIfNewer(ctx, pool.Address, activeContract.GetCreatedEvent().GetOffset(), config.TokenPoolTypeLockRelease)
	}
}

// registerIfNewer registers a discovered pool if its ledger offset is newer than whatever is
// currently registered for the same address, so a later-created contract always wins over an
// earlier one that happens to derive the same address (see discoveredPools).
// Caller must hold s.mux.
func (s *PoolDiscoveryService) registerIfNewer(ctx context.Context, rawAddress contracts.RawInstanceAddress, offset int64, poolType config.TokenPoolType) {
	address := rawAddress.InstanceAddress()
	if existingOffset, ok := s.discoveredPools[address]; ok && existingOffset >= offset {
		return
	}

	if err := s.registerPool(ctx, address, rawAddress.Owner(), poolType); err != nil {
		s.logger.Err(err).Stringer("address", address).Msg("failed to register pool")
		return
	}

	s.discoveredPools[address] = offset
	s.logger.Info().Stringer("address", address).Str("type", string(poolType)).Msg("discovered and registered new token pool")
}

// registerPool registers a discovered pool with the tokenpool server. PartyID stays the
// observer party (it's what EDS has rights to query the ledger with); owner is the pool's
// actual owner, read off its own instance address rather than the observer party.
// The factory is always wired in URL mode against the configured token-standard backend:
// the pool's LockOrBurn/ReleaseFromTicket require the factory CID in the choice context,
// and URL mode resolves it over HTTP, so EDS needs no ledger visibility into the factory.
func (s *PoolDiscoveryService) registerPool(ctx context.Context, address contracts.InstanceAddress, owner string, poolType config.TokenPoolType) error {
	poolConfig := config.TokenPool{
		ContractIdentifier: config.ContractIdentifier{
			PartyID:         string(s.observerParty),
			InstanceAddress: address,
		},
		Type:      poolType,
		PoolOwner: owner,
		Factory: &config.Factory{
			Type:             config.FactoryTypeURL,
			TokenStandardURL: &s.tokenStandardURL,
		},
	}

	if err := s.tokenPoolServer.RegisterDiscoveredPool(ctx, poolConfig); err != nil {
		return fmt.Errorf("failed to register pool with tokenpool server: %w", err)
	}

	return nil
}
