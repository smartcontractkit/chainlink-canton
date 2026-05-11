package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"

	edsCommon "github.com/smartcontractkit/chainlink-canton/eds/common"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/ccip"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/ccv"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/executor"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/global"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/middleware"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/token_standard"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api/tokenpool"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	"github.com/smartcontractkit/chainlink-canton/eds/monitoring"
	oapiCCIP "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccip"
	oapiCCV "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/ccv"
	oapiExecutor "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/executor"
	oapiGlobal "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/global"
	oapiTokenPool "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/tokenpool"
	oapiTokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	oapiTransferInstruction "github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
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

	activeContractStore := store.NewActiveContractStore(
		logger,
		cantonChain.Participants[0].LedgerServices.Update,
		cantonChain.Participants[0].LedgerServices.State,
		metrics.With("store", "ActiveContractStore"),
	)

	instrumentHoldingStore := store.NewInstrumentHoldingStore(
		logger,
		cantonChain.Participants[0].LedgerServices.Update,
		cantonChain.Participants[0].LedgerServices.State,
		metrics.With("store", "InstrumentHoldingStore"),
	)

	// Create HTTP Server
	router := gin.Default()
	router.Use(middleware.RequestMonitoringMiddleware(metrics))

	errChan := make(chan error)
	var globalAddressFilters []global.InstanceAddressFilter
	if cfg.CCIPAPIConfig.Enabled {
		ccipAPIServer, err := ccip.NewServer(ctx, logger, activeContractStore, cfg.CCIPAPIConfig)
		if err != nil {
			return fmt.Errorf("failed to create CCIP API: %w", err)
		}
		oapiCCIP.RegisterHandlers(router, ccipAPIServer)
		globalAddressFilters = append(globalAddressFilters, ccipAPIServer)
	}
	if cfg.CCVAPIConfig.Enabled {
		ccvAPIServer, err := ccv.NewServer(ctx, logger, activeContractStore, cfg.CCVAPIConfig)
		if err != nil {
			return fmt.Errorf("failed to create CCV API: %w", err)
		}
		oapiCCV.RegisterHandlers(router, ccvAPIServer)
		globalAddressFilters = append(globalAddressFilters, ccvAPIServer)
	}
	if cfg.ExecutorAPIConfig.Enabled {
		executorAPIServer, err := executor.NewServer(ctx, logger, activeContractStore, cfg.ExecutorAPIConfig)
		if err != nil {
			return fmt.Errorf("failed to create Executor API: %w", err)
		}
		oapiExecutor.RegisterHandlers(router, executorAPIServer)
		globalAddressFilters = append(globalAddressFilters, executorAPIServer)
	}
	if cfg.TokenPoolAPIConfig.Enabled {
		tokenPoolAPIServer, err := tokenpool.NewServer(ctx, logger, activeContractStore, instrumentHoldingStore, cfg.TokenPoolAPIConfig)
		if err != nil {
			return fmt.Errorf("failed to create TokenPool API: %w", err)
		}
		oapiTokenPool.RegisterHandlers(router, tokenPoolAPIServer)
		globalAddressFilters = append(globalAddressFilters, tokenPoolAPIServer)

		// Run instrument holding store in the background
		// This should only be run if the TokenPool API is enabled, as it will fail if no filters are specified.
		go func(errChan chan<- error) {
			logger.Info().Msg("starting instrument holding store")
			err := instrumentHoldingStore.Run(ctx, store.DefaultStreamConfig())
			if err != nil {
				errChan <- fmt.Errorf("failed to run instrument holding store: %w", err)
			}
		}(errChan)
	}
	if cfg.TokenStandardAPIConfig.Enabled {
		tokenStandardAPIServer, err := token_standard.NewServer(ctx, logger, activeContractStore, cfg.TokenStandardAPIConfig)
		if err != nil {
			return fmt.Errorf("failed to create TokenStandard API: %w", err)
		}
		oapiTokenMetadataV1.RegisterHandlers(router, tokenStandardAPIServer)
		oapiTransferInstruction.RegisterHandlers(router, tokenStandardAPIServer)
		globalAddressFilters = append(globalAddressFilters, tokenStandardAPIServer)
	}

	// Global API
	// Passing all configured API implementations as filters, as the Global API should only return disclosures for
	// addresses that are already returned as part of the other API implementations.
	globalAPIServer, err := global.NewServer(ctx, logger, activeContractStore, cfg.GlobalAPIConfig, globalAddressFilters...)
	if err != nil {
		return fmt.Errorf("failed to create Global API: %w", err)
	}
	oapiGlobal.RegisterHandlers(router, globalAPIServer)

	// Run update store in the background
	go func(errChan chan<- error) {
		logger.Info().Msg("starting active contract store")
		err := activeContractStore.Run(ctx, store.DefaultStreamConfig())
		if err != nil {
			errChan <- fmt.Errorf("failed to run active contract store: %w", err)
		}
	}(errChan)

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
