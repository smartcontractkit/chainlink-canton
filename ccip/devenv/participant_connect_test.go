package devenv

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

func TestConnectCantonParticipant(t *testing.T) {
	grpcURL := strings.TrimSpace(os.Getenv("CANTON_GRPC_URL"))
	authURL := strings.TrimSpace(os.Getenv(envCantonAuthURL))
	clientID := oauthClientID()
	clientSecret := oauthClientSecret()

	if grpcURL == "" || authURL == "" || clientID == "" || clientSecret == "" {
		t.Skip("skipping: set CANTON_GRPC_URL, CANTON_AUTH_URL, and CLIENT_ID/CLIENT_SECRET (or CANTON_OAUTH_*) to run")
	}

	authCfg, err := buildParticipantAuthConfig(blockchain.CantonParticipantEndpoints{
		GRPCLedgerAPIURL: grpcURL,
	})
	if err != nil {
		t.Fatalf("buildParticipantAuthConfig: %v", err)
	}
	if authCfg.Type != commonconfig.AuthTypeClientCredentials {
		t.Fatalf("expected clientCredentials auth, got %q", authCfg.Type)
	}

	ctx := context.Background()
	provider, err := participantAuthProvider(ctx, authCfg)
	if err != nil {
		t.Fatalf("participantAuthProvider: %v", err)
	}

	requireBearerToken(t, ctx, provider)

	conn, err := grpc.NewClient(
		grpcURL,
		grpc.WithTransportCredentials(provider.TransportCredentials()),
		grpc.WithPerRPCCredentials(provider.PerRPCCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	defer conn.Close()

	waitForGRPCReady(t, ctx, conn)
	t.Log("connected to Canton participant (client credentials OAuth + authenticated gRPC channel)")
}

func requireBearerToken(t *testing.T, ctx context.Context, provider interface {
	PerRPCCredentials() credentials.PerRPCCredentials
}) {
	t.Helper()

	md, err := provider.PerRPCCredentials().GetRequestMetadata(ctx, "https://canton/")
	if err != nil {
		t.Fatalf("fetch OAuth access token: %v", err)
	}

	for key, value := range md {
		if strings.EqualFold(key, "authorization") && strings.HasPrefix(strings.ToLower(value), "bearer ") {
			return
		}
	}

	t.Fatalf("expected bearer token in RPC metadata, got %#v", md)
}

func waitForGRPCReady(t *testing.T, ctx context.Context, conn *grpc.ClientConn) {
	t.Helper()

	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	conn.Connect()
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return
		}
		if !conn.WaitForStateChange(waitCtx, state) {
			t.Fatalf("timed out waiting for gRPC connection (last state: %s)", state)
		}
	}
}
