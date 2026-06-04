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
		authProvider := authentication.NewInsecureStaticProvider(config.JWT)
		// Get Primary Party for user
		ledgerApiConn, err := grpc.NewClient(
			config.GRPCLedgerAPIURL,
			grpc.WithTransportCredentials(authProvider.TransportCredentials()),
			grpc.WithPerRPCCredentials(authProvider.PerRPCCredentials()),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create gRPC connection to Ledger API for Canton participant %d: %w", i+1, err)
		}
		userResp, err := adminv2.NewUserManagementServiceClient(ledgerApiConn).GetUser(context.Background(), &adminv2.GetUserRequest{UserId: config.UserID})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get user info for user %s for Canton participant %d: %w", config.UserID, i+1, err)
		}
		party := userResp.GetUser().GetPrimaryParty()
		if party == "" {
			return nil, 0, fmt.Errorf("no primary party found for user %s for Canton participant %d", config.UserID, i+1)
		}
		_ = ledgerApiConn.Close()

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
	p, err := cldf_canton_provider.NewRPCChainProvider(d.ChainSelector, providerConfig).Initialize(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	return p, d.ChainSelector, nil
}
