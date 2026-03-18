package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/feequoter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/lockreleasetokenpool"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	edsCommon "github.com/smartcontractkit/chainlink-canton/eds/common"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	"github.com/smartcontractkit/chainlink-canton/eds/monitoring"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

func RunEDS(ctx context.Context, logger zerolog.Logger, cfg *config.Config) error {
	// Validate config
	if err := cfg.Validate(); err != nil {
		return err
	}

	chainSelector, err := strconv.ParseUint(cfg.ChainSelector, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain selector: %w", err)
	}
	chainDetails, err := chainsel.GetChainDetails(chainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chain details: %w", err)
	}

	// Set up monitoring
	var metrics edsCommon.EDSMetricLabeler
	if cfg.Monitoring.Enabled {
		metrics, err = monitoring.InitBeholderMonitoring(beholder.Config{
			InsecureConnection:       cfg.Monitoring.Beholder.InsecureConnection,
			CACertFile:               cfg.Monitoring.Beholder.CACertFile,
			OtelExporterGRPCEndpoint: cfg.Monitoring.Beholder.OtelExporterGRPCEndpoint,
			OtelExporterHTTPEndpoint: cfg.Monitoring.Beholder.OtelExporterHTTPEndpoint,
			LogStreamingEnabled:      cfg.Monitoring.Beholder.LogStreamingEnabled,
			MetricReaderInterval:     time.Second * time.Duration(cfg.Monitoring.Beholder.MetricReaderInterval),
			TraceSampleRatio:         cfg.Monitoring.Beholder.TraceSampleRatio,
			TraceBatchTimeout:        time.Second * time.Duration(cfg.Monitoring.Beholder.TraceBatchTimeout),
		})
		if err != nil {
			return fmt.Errorf("failed to initialize Beholder monitoring: %w", err)
		}
	} else {
		metrics = monitoring.NoopEDSMetricLabeler{}
	}

	authProvider, err := cfg.Node.AuthConfig.NewProvider(ctx)
	if err != nil {
		return fmt.Errorf("auth config: %w", err)
	}
	chain, err := provider.NewRPCChainProvider(chainDetails.ChainSelector, provider.RPCChainProviderConfig{
		Participants: []provider.ParticipantConfig{
			{
				Endpoints: provider.Endpoints{
					JSONLedgerAPIURL: "", // not used
					GRPCLedgerAPIURL: cfg.Node.URL,
					AdminAPIURL:      "", // not used
					ValidatorAPIURL:  "", // not used
				},
				UserID:       cfg.Node.AuthConfig.UserID,
				PartyID:      "party-id", // not used, will be taken from contract config instead
				AuthProvider: authProvider,
			},
		},
	}).Initialize(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize RPC chain: %w", err)
	}
	cantonChain := chain.(*canton.Chain)

	templates := []store.RegisteredTemplate{
		{
			TemplateID: contracts.TemplateIDFromBinding(perpartyrouter.PerPartyRouterFactory{}),
			PartyID:    cfg.Contracts.OnRamp.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(onramp.OnRamp{}),
			PartyID:    cfg.Contracts.OnRamp.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(offramp.OffRamp{}),
			PartyID:    cfg.Contracts.OffRamp.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(common.GlobalConfig{}),
			PartyID:    cfg.Contracts.GlobalConfig.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(tokenadminregistry.TokenAdminRegistry{}),
			PartyID:    cfg.Contracts.TokenAdminRegistry.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(rmn.RMNRemote{}),
			PartyID:    cfg.Contracts.RMNRemote.PartyID,
		},
		{
			TemplateID: contracts.TemplateIDFromBinding(feequoter.FeeQuoter{}),
			PartyID:    cfg.Contracts.FeeQuoter.PartyID,
		},
	}
	ccvCids := make([]contracts.InstanceAddress, len(cfg.Contracts.CCVs))
	for i, ccv := range cfg.Contracts.CCVs {
		templates = append(templates, store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(ccvs.CommitteeVerifier{}),
			PartyID:    ccv.PartyID,
		})
		ccvCids[i] = ccv.InstanceAddress
	}
	tokenPoolCids := make([]contracts.InstanceAddress, len(cfg.Contracts.TokenPools))
	for i, tokenPool := range cfg.Contracts.TokenPools {
		templates = append(templates, store.RegisteredTemplate{
			TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
			PartyID:    tokenPool.PartyID,
		})
		tokenPoolCids[i] = tokenPool.InstanceAddress
	}

	updateStore, err := store.NewUpdateStore(ctx, store.UpdateStoreConfig{
		Logger:        logger,
		UpdateService: cantonChain.Participants[0].LedgerServices.Update,
		StateService:  cantonChain.Participants[0].LedgerServices.State,
		MaxRetries:    cfg.Node.MaxRetries,
	},
		metrics,
		templates...,
	)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	instrumentHoldingStore := store.NewInstrumentHoldingStore(store.InstrumentHoldingStoreConfig{
		Logger:        logger,
		UpdateService: cantonChain.Participants[0].LedgerServices.Update,
		StateService:  cantonChain.Participants[0].LedgerServices.State,
		MaxRetries:    cfg.Node.MaxRetries,
		Owner:         types.PARTY(cfg.Contracts.PoolOwner),
	})

	// Run update store in the background
	errChan := make(chan error)
	go func(errChan chan<- error) {
		logger.Info().Msg("starting update store")
		err := updateStore.Run(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run store: %w", err)
		}
	}(errChan)

	// Run instrument holding store in the background
	go func(errChan chan<- error) {
		logger.Info().Msg("starting instrument holding store")
		err := instrumentHoldingStore.Run(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run instrument holding store: %w", err)
		}
	}(errChan)

	disclosureSvc := disclosure.NewDisclosureService(ctx, disclosure.DisclosureServiceConfig{
		ContractStore:          updateStore,
		InstrumentHoldingStore: instrumentHoldingStore,
		PerPartyRouterFactory:  cfg.Contracts.PerPartyRouterFactory.InstanceAddress,
		OnRamp:                 cfg.Contracts.OnRamp.InstanceAddress,
		OffRamp:                cfg.Contracts.OffRamp.InstanceAddress,
		GlobalConfig:           cfg.Contracts.GlobalConfig.InstanceAddress,
		TokenAdminRegistry:     cfg.Contracts.TokenAdminRegistry.InstanceAddress,
		RMNRemote:              cfg.Contracts.RMNRemote.InstanceAddress,
		FeeQuoter:              cfg.Contracts.FeeQuoter.InstanceAddress,
		CCVs:                   ccvCids,
		TokenPools:             tokenPoolCids,
	})

	server := api.NewServer(logger, disclosureSvc)

	r := gin.Default()
	r.Use(middleware.RequestMonitoringMiddleware(metrics))

	edsv1.RegisterHandlers(r, server)

	s := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in the background
	go func() {
		logger.Info().Msg("starting server")
		err := s.ListenAndServe()
		if err != nil {
			errChan <- fmt.Errorf("failed to run server: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
