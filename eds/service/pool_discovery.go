package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/burnminttokenpool"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/lockreleasetokenpool"
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
	discoveredPools     map[contracts.InstanceAddress]bool
	mux                 sync.Mutex
}

// NewPoolDiscoveryService creates a new pool discovery service.
// It must be created BEFORE activeContractStore.Run() is called.
// observerPartyID is EDS's own party, named as an observer on discoverable pool contracts.
func NewPoolDiscoveryService(
	logger zerolog.Logger,
	activeContractStore store.ActiveContractStoreInterface,
	tokenPoolServer *tokenpool.Server,
	observerPartyID string,
) *PoolDiscoveryService {
	logger = logger.With().Str("component", "PoolDiscoveryService").Logger()

	// Pre-register both pool templates to watch for CreatedEvents.
	activeContractStore.RegisterTemplates(
		store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(burnminttokenpool.BurnMintTokenPool{}),
			PartyID:    observerPartyID,
		},
		store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
			PartyID:    observerPartyID,
		},
	)

	return &PoolDiscoveryService{
		logger:              logger,
		activeContractStore: activeContractStore,
		tokenPoolServer:     tokenPoolServer,
		observerParty:       types.PARTY(observerPartyID),
		discoveredPools:     make(map[contracts.InstanceAddress]bool),
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
		s.registerIfNew(ctx, pool.Address, config.TokenPoolTypeBurnMint)
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
		s.registerIfNew(ctx, pool.Address, config.TokenPoolTypeLockRelease)
	}
}

// registerIfNew registers a discovered pool if it has not already been discovered.
// Caller must hold s.mux.
func (s *PoolDiscoveryService) registerIfNew(ctx context.Context, rawAddress contracts.RawInstanceAddress, poolType config.TokenPoolType) {
	address := rawAddress.InstanceAddress()
	if s.discoveredPools[address] {
		return
	}

	if err := s.registerPool(ctx, address, rawAddress.Owner(), poolType); err != nil {
		s.logger.Err(err).Stringer("address", address).Msg("failed to register pool")
		return
	}

	s.discoveredPools[address] = true
	s.logger.Info().Stringer("address", address).Str("type", string(poolType)).Msg("discovered and registered new token pool")
}

// registerPool registers a discovered pool with the tokenpool server. PartyID stays the
// observer party (it's what EDS has rights to query the ledger with); owner is the pool's
// actual owner, read off its own instance address rather than the observer party.
func (s *PoolDiscoveryService) registerPool(ctx context.Context, address contracts.InstanceAddress, owner string, poolType config.TokenPoolType) error {
	poolConfig := config.TokenPool{
		ContractIdentifier: config.ContractIdentifier{
			PartyID:         string(s.observerParty),
			InstanceAddress: address,
		},
		Type:      poolType,
		PoolOwner: owner,
	}

	if err := s.tokenPoolServer.RegisterDiscoveredPool(ctx, poolConfig); err != nil {
		return fmt.Errorf("failed to register pool with tokenpool server: %w", err)
	}

	return nil
}
