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

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/ccvs"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/offramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/onramp"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/perpartyrouter"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/api"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/store"
	edsv1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds"
)

func RunEDS(ctx context.Context, logger zerolog.Logger, cfg *config.Config) error {
	chainSelector, err := strconv.ParseUint(cfg.ChainSelector, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chain selector: %w", err)
	}
	chainDetails, err := chainsel.GetChainDetails(chainSelector)
	if err != nil {
		return fmt.Errorf("failed to get chain details: %w", err)
	}

	var authProvider authentication.Provider
	switch cfg.Node.AuthConfig.Type {
	case "static":
		authProvider = authentication.NewInsecureStaticProvider(cfg.Node.AuthConfig.JWT)
	case "clientCredentials":
		authProvider, err = clientcredentials.NewDiscoveryProvider(ctx, cfg.Node.AuthConfig.AuthURL, cfg.Node.AuthConfig.ClientID, cfg.Node.AuthConfig.ClientSecret)
		if err != nil {
			return fmt.Errorf("failed to create client credentials provider: %w", err)
		}
	case "authorizationCode":
		authProvider, err = authorizationcode.NewDiscoveryProvider(ctx, cfg.Node.AuthConfig.AuthURL, cfg.Node.AuthConfig.ClientID)
		if err != nil {
			return fmt.Errorf("failed to create authorization code provider: %w", err)
		}
	default:
		return fmt.Errorf("unsupported authentication type: %s", cfg.Node.AuthConfig.Type)
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
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(perpartyrouter.PerPartyRouterFactory{}),
			PartyID:    cfg.Contracts.OnRamp.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(onramp.OnRamp{}),
			PartyID:    cfg.Contracts.OnRamp.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(offramp.OffRamp{}),
			PartyID:    cfg.Contracts.OffRamp.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(common.GlobalConfig{}),
			PartyID:    cfg.Contracts.GlobalConfig.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(tokenadminregistry.TokenAdminRegistry{}),
			PartyID:    cfg.Contracts.TokenAdminRegistry.PartyID,
		},
		store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(rmn.RMNRemote{}),
			PartyID:    cfg.Contracts.RMNRemote.PartyID,
		},
	}
	ccvCids := make([]contracts.InstanceAddress, len(cfg.Contracts.CCVs))
	for i, ccv := range cfg.Contracts.CCVs {
		templates = append(templates, store.RegisteredTemplate{
			TemplateID: store.TemplateIDFromBinding(ccvs.CommitteeVerifier{}),
			PartyID:    ccv.PartyID,
		})
		ccvCids[i] = ccv.InstanceAddress
	}

	updateStore, err := store.NewUpdateStore(ctx, store.UpdateStoreConfig{
		Logger:        logger,
		UpdateService: cantonChain.Participants[0].LedgerServices.Update,
		StateService:  cantonChain.Participants[0].LedgerServices.State,
		MaxRetries:    cfg.Node.MaxRetries,
	},
		templates...,
	)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Run store in the background
	errChan := make(chan error)
	go func(errChan chan<- error) {
		logger.Info().Msg("starting update store")
		err := updateStore.Run(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run store: %w", err)
		}
	}(errChan)

	disclosureSvc := disclosure.NewDisclosureService(ctx, disclosure.DisclosureServiceConfig{
		ContractStore:         updateStore,
		PerPartyRouterFactory: cfg.Contracts.PerPartyRouterFactory.InstanceAddress,
		OnRamp:                cfg.Contracts.OnRamp.InstanceAddress,
		OffRamp:               cfg.Contracts.OffRamp.InstanceAddress,
		GlobalConfig:          cfg.Contracts.GlobalConfig.InstanceAddress,
		TokenAdminRegistry:    cfg.Contracts.TokenAdminRegistry.InstanceAddress,
		RMNRemote:             cfg.Contracts.RMNRemote.InstanceAddress,
		CCVs:                  ccvCids,
	})

	server := api.NewServer(logger, disclosureSvc)

	r := gin.Default()

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
