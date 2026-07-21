package ledger

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/config"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cantonProvider "github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider"
)

// DevnetClient is a ledger.Client backed by a devnet/RPC Canton participant (JWT auth).
type DevnetClient = CTFClient

// NewDevnetClient wraps a CLDF RPC participant.
func NewDevnetClient(participant canton.Participant) *DevnetClient {
	return NewCTFClient(participant)
}

// ConnectDevnet builds a ledger client from registry-kit config.
func ConnectDevnet(ctx context.Context, cfg config.Config, actAsParty string) (*DevnetClient, canton.Participant, error) {
	authProvider, err := cfg.Ledger.Auth.NewProvider(ctx)
	if err != nil {
		return nil, canton.Participant{}, fmt.Errorf("auth provider: %w", err)
	}

	providerCfg := cantonProvider.RPCChainProviderConfig{
		Participants: []cantonProvider.ParticipantConfig{{
			Endpoints: cantonProvider.Endpoints{
				JSONLedgerAPIURL: cfg.Ledger.JSONAPIURL,
				GRPCLedgerAPIURL: cfg.Ledger.GRPCLedgerAPIURL,
				AdminAPIURL:      cfg.Ledger.AdminAPIURL,
				ValidatorAPIURL:  cfg.Ledger.ValidatorAPIURL,
			},
			UserID:       cfg.Ledger.UserID,
			PartyID:      actAsParty,
			AuthProvider: authProvider,
		}},
	}

	chainProvider := cantonProvider.NewRPCChainProvider(cfg.ChainSelector, providerCfg)
	bc, err := chainProvider.Initialize(ctx)
	if err != nil {
		return nil, canton.Participant{}, fmt.Errorf("initialize RPC chain: %w", err)
	}

	chain, ok := bc.(*canton.Chain)
	if !ok || len(chain.Participants) == 0 {
		return nil, canton.Participant{}, fmt.Errorf("RPC chain has no participants")
	}

	participant := chain.Participants[0]
	participant.PartyID = actAsParty

	return NewDevnetClient(participant), participant, nil
}

// ConnectDevnetWithStaticJWT is a test/helper entry when auth is a pre-set JWT env var.
func ConnectDevnetWithStaticJWT(ctx context.Context, cfg config.Config, actAsParty, jwt string) (*DevnetClient, canton.Participant, error) {
	cfgCopy := cfg
	cfgCopy.Ledger.Auth = commonconfig.AuthConfig{
		Type: commonconfig.AuthTypeInsecureStatic,
		JWT:  jwt,
	}

	return ConnectDevnet(ctx, cfgCopy, actAsParty)
}
