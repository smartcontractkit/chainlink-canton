package devenv

import (
	"context"
	"fmt"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"google.golang.org/grpc"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_canton_provider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
)

func NewCLDF(ctx context.Context, b *blockchain.Input) (cldf_chain.BlockChain, uint64, error) {
	d, err := chainsel.GetChainDetailsByChainIDAndFamily(b.Out.ChainID, chainsel.FamilyCanton)
	if err != nil {
		return nil, 0, err
	}

	providerConfig := cldf_canton_provider.RPCChainProviderConfig{
		Participants: make([]cldf_canton_provider.ParticipantConfig, len(b.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants)),
	}

	for i, config := range b.Out.NetworkSpecificData.CantonData.ExternalEndpoints.Participants {
		authCfg, err := buildParticipantAuthConfig(config)
		if err != nil {
			return nil, 0, fmt.Errorf("participant %d auth config: %w", i+1, err)
		}
		authProvider, err := participantAuthProvider(ctx, authCfg)
		if err != nil {
			return nil, 0, fmt.Errorf("participant %d auth provider: %w", i+1, err)
		}

		party, err := primaryPartyFromLedgerAPI(ctx, i+1, config.GRPCLedgerAPIURL, config.UserID, authProvider)
		if err != nil {
			return nil, 0, err
		}

		providerConfig.Participants[i] = cldf_canton_provider.ParticipantConfig{
			Endpoints: cldf_canton_provider.Endpoints{
				JSONLedgerAPIURL: config.JSONLedgerAPIURL,
				GRPCLedgerAPIURL: config.GRPCLedgerAPIURL,
				AdminAPIURL:      config.AdminAPIURL,
				ValidatorAPIURL:  config.ValidatorAPIURL,
			},
			InternalEndpoints: &cldf_canton_provider.Endpoints{
				JSONLedgerAPIURL: b.Out.NetworkSpecificData.CantonData.InternalEndpoints.Participants[i].JSONLedgerAPIURL,
				GRPCLedgerAPIURL: b.Out.NetworkSpecificData.CantonData.InternalEndpoints.Participants[i].GRPCLedgerAPIURL,
				AdminAPIURL:      b.Out.NetworkSpecificData.CantonData.InternalEndpoints.Participants[i].AdminAPIURL,
				ValidatorAPIURL:  b.Out.NetworkSpecificData.CantonData.InternalEndpoints.Participants[i].ValidatorAPIURL,
			},
			UserID:       config.UserID,
			PartyID:      party,
			AuthProvider: authProvider,
		}
	}
	p, err := cldf_canton_provider.NewRPCChainProvider(d.ChainSelector, providerConfig).Initialize(ctx)
	if err != nil {
		return nil, 0, err
	}

	return p, d.ChainSelector, nil
}

func primaryPartyFromLedgerAPI(
	ctx context.Context,
	participantIndex int,
	grpcLedgerAPIURL string,
	userID string,
	authProvider authentication.Provider,
) (string, error) {
	conn, err := grpc.NewClient(
		grpcLedgerAPIURL,
		grpc.WithTransportCredentials(authProvider.TransportCredentials()),
		grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create gRPC connection to Ledger API for Canton participant %d: %w", participantIndex, err)
	}
	defer conn.Close()

	userResp, err := adminv2.NewUserManagementServiceClient(conn).GetUser(ctx, &adminv2.GetUserRequest{UserId: userID})
	if err != nil {
		return "", fmt.Errorf("failed to get user info for user %s for Canton participant %d: %w", userID, participantIndex, err)
	}
	party := userResp.GetUser().GetPrimaryParty()
	if party == "" {
		return "", fmt.Errorf("no primary party found for user %s for Canton participant %d", userID, participantIndex)
	}

	return party, nil
}
