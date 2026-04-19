package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/ccv"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/executor"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/tokenpool"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	executorBinding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/executor"
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
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/ccip"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	"github.com/smartcontractkit/chainlink-canton/eds/monitoring"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
)

func RunEDS(ctx context.Context, logger zerolog.Logger, cfg *config.Config) error {
	cfg, err := config.DefaultConfig().Merge(cfg)
	if err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

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

	// Stores

	activeContractStore, err := store.NewActiveContractStore(
		store.ActiveContractStoreConfig{
			Logger:        logger,
			UpdateService: cantonChain.Participants[0].LedgerServices.Update,
			StateService:  cantonChain.Participants[0].LedgerServices.State,
			StreamConfig: store.ReliableStreamConfig{
				MaxRetries: cfg.Node.MaxRetries,
			},
		},
		metrics.With("store", "ActiveContractStore"),
	)
	if err != nil {
		return fmt.Errorf("failed to create active contract store: %w", err)
	}

	var instrumentHoldingStore *store.ContractStore[splice_api_token_holding_v1.InstrumentId, *apiv2.DisclosedContract]
	/*instrumentHoldingStore = store.NewInstrumentHoldingStore(
		store.InstrumentHoldingStoreConfig{
			Logger:        logger,
			Owner:         types.PARTY(cfg.Contracts.PoolOwner),
			UpdateService: cantonChain.Participants[0].LedgerServices.Update,
			StateService:  cantonChain.Participants[0].LedgerServices.State,
			StreamConfig: store.ReliableStreamConfig{
				MaxRetries: cfg.Node.MaxRetries,
			},
		},
		metrics.With("store", "InstrumentHoldingStore"),
	)*/

	// Create HTTP Server
	router := gin.Default()
	router.Use(middleware.RequestMonitoringMiddleware(metrics))

	var templates []store.RegisteredTemplate
	if cfg.CCIPAPIConfig.Enabled {
		// Register templates
		templates = append(templates, []store.RegisteredTemplate{
			{
				TemplateID: contracts.TemplateIDFromBinding(perpartyrouter.PerPartyRouterFactory{}),
				PartyID:    cfg.CCIPAPIConfig.OnRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(onramp.OnRamp{}),
				PartyID:    cfg.CCIPAPIConfig.OnRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(offramp.OffRamp{}),
				PartyID:    cfg.CCIPAPIConfig.OffRamp.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(common.GlobalConfig{}),
				PartyID:    cfg.CCIPAPIConfig.GlobalConfig.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(tokenadminregistry.TokenAdminRegistry{}),
				PartyID:    cfg.CCIPAPIConfig.TokenAdminRegistry.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(rmn.RMNRemote{}),
				PartyID:    cfg.CCIPAPIConfig.RMNRemote.PartyID,
			}, {
				TemplateID: contracts.TemplateIDFromBinding(feequoter.FeeQuoter{}),
				PartyID:    cfg.CCIPAPIConfig.FeeQuoter.PartyID,
			},
		}...)

		// Register API server
		ccipAPIServer := ccip.NewServer(logger, activeContractStore, cfg.CCIPAPIConfig)
		oapiCCIP.RegisterHandlers(router, ccipAPIServer)
	}
	if cfg.CCVAPIConfig.Enabled {
		// Register templates
		for _, v := range cfg.CCVAPIConfig.CCVs {
			templates = append(templates, store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(ccvs.CommitteeVerifier{}),
				PartyID:    v.PartyID,
			})
		}

		// Register API server
		ccvAPIServer := ccv.NewServer(logger, activeContractStore, cfg.CCVAPIConfig)
		oapiCCV.RegisterHandlers(router, ccvAPIServer)
	}
	if cfg.ExecutorAPIConfig.Enabled {
		// Register templates
		for _, v := range cfg.ExecutorAPIConfig.Executors {
			templates = append(templates, store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(executorBinding.Executor{}),
				PartyID:    v.PartyID,
			})
		}

		// Register API server
		executorAPIServer := executor.NewServer(logger, activeContractStore, cfg.ExecutorAPIConfig)
		oapiExecutor.RegisterHandlers(router, executorAPIServer)
	}
	if cfg.TokenPoolAPIConfig.Enabled {
		// Register templates
		for _, v := range cfg.TokenPoolAPIConfig.TokenPools {
			switch v.Type {
			case config.TokenPoolTypeLockRelease:
				templates = append(templates, store.RegisteredTemplate{
					TemplateID: contracts.TemplateIDFromBinding(lockreleasetokenpool.LockReleaseTokenPool{}),
					PartyID:    v.PartyID,
				})
			case config.TokenPoolTypeBurnMint:
				fallthrough
			default:
				return fmt.Errorf("unsupported token pool type: %s", v.Type)
			}
			templates = append(templates, store.RegisteredTemplate{
				TemplateID: contracts.TemplateIDFromBinding(common.RateLimiter{}),
				PartyID:    v.PartyID,
			})
		}

		// Register API server
		tokenPoolAPIServer := tokenpool.NewServer(logger, activeContractStore, instrumentHoldingStore, cfg.TokenPoolAPIConfig)
		oapiTokenPool.RegisterHandlers(router, tokenPoolAPIServer)
	}

	// Run update store in the background
	errChan := make(chan error)
	go func(errChan chan<- error) {
		logger.Info().Msg("starting active contract store")
		err := activeContractStore.Run(ctx, store.WithFiltersByParty(store.ActiveContractStoreFilters(templates...)))
		if err != nil {
			errChan <- fmt.Errorf("failed to run active contract store: %w", err)
		}
	}(errChan)

	/*// Run instrument holding store in the background
	go func(errChan chan<- error) {
		logger.Info().Msg("starting instrument holding store")
		err := instrumentHoldingStore.Run(ctx, nil)
		if err != nil {
			errChan <- fmt.Errorf("failed to run instrument holding store: %w", err)
		}
	}(errChan)*/

	s := &http.Server{
		Handler:      router,
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
