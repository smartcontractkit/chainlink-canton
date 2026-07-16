package devenv

import (
	"context"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_canton_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

func NewCLDF(ctx context.Context, b *blockchain.Input) (cldf_chain.BlockChain, uint64, error) {
	d, err := chainsel.GetChainDetailsByChainIDAndFamily(b.Out.ChainID, chainsel.FamilyCanton)
	if err != nil {
		return nil, 0, err
	}

	providerConfig := cldf_canton_provider.RPCChainProviderConfig{
		Participants: make([]cldf_canton_provider.ParticipantConfig, len(b.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants)),
	}

	presetPartyID := resolvePartyID()
	internalParticipants := b.Out.NetworkSpecificData.CantonData.InternalEndpoints.Participants

	for i, config := range b.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants {
		authCfg := resolveAuthConfig(config.JWT, config.UserID)
		authProvider, err := newAuthProvider(ctx, authCfg)
		if err != nil {
			return nil, 0, fmt.Errorf("create auth provider for Canton participant %d: %w", i+1, err)
		}

		userID := resolveUserID(config.UserID)
		if userID == "" && authCfg.Type == commonconfig.AuthTypeAuthorizationCode {
			userID, err = userIDFromToken(ctx, authProvider)
			if err != nil {
				return nil, 0, fmt.Errorf("resolve user id for Canton participant %d: %w", i+1, err)
			}
		}

		grpcURL := resolveGRPCLedgerURL(config.GRPCLedgerAPIURL)
		party, err := func() (string, error) {
			conn, err := grpc.NewClient(
				grpcURL,
				grpc.WithTransportCredentials(authProvider.TransportCredentials()),
				grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
			)
			if err != nil {
				return "", fmt.Errorf("create gRPC connection to Ledger API for Canton participant %d: %w", i+1, err)
			}
			defer conn.Close()

			party, err := resolveParticipantParty(ctx, conn, userID, presetPartyID)
			if err != nil {
				return "", fmt.Errorf("resolve party for Canton participant %d: %w", i+1, err)
			}

			return party, nil
		}()
		if err != nil {
			return nil, 0, err
		}

		var internalEndpoints *cldf_canton_provider.Endpoints
		if i < len(internalParticipants) {
			internal := internalParticipants[i]
			if endpointsNonEmpty(
				internal.JSONLedgerAPIURL,
				internal.GRPCLedgerAPIURL,
				internal.AdminAPIURL,
				internal.ValidatorAPIURL,
			) {
				internalEndpoints = &cldf_canton_provider.Endpoints{
					JSONLedgerAPIURL: internal.JSONLedgerAPIURL,
					GRPCLedgerAPIURL: internal.GRPCLedgerAPIURL,
					AdminAPIURL:      internal.AdminAPIURL,
					ValidatorAPIURL:  internal.ValidatorAPIURL,
				}
			}
		}

		providerConfig.Participants[i] = cldf_canton_provider.ParticipantConfig{
			Endpoints: cldf_canton_provider.Endpoints{
				JSONLedgerAPIURL: config.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: grpcURL,
				AdminAPIURL:      config.AdminAPIURL,
				ValidatorAPIURL:  resolveValidatorAPIURL(config.ValidatorAPIURL),
			},
			InternalEndpoints: internalEndpoints,
			UserID:            userID,
			PartyID:           party,
			AuthProvider:      authProvider,
		}
	}

	p, err := cldf_canton_provider.NewRPCChainProvider(d.ChainSelector, providerConfig).Initialize(ctx)
	if err != nil {
		return nil, 0, err
	}

	return p, d.ChainSelector, nil
}
