package runtime

import (
	"context"
	"fmt"

	apiv2interactive "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// Participant contains the endpoint and auth data needed to run ceremony flows
// against one Canton participant.
type Participant struct {
	Name             string
	AdminAPIURL      string
	GRPCLedgerAPIURL string
	UserID           string
	AdminJWT         string
	LedgerJWT        string
}

// FromCantonParticipant converts a CLDF Canton participant into the endpoint
// and auth shape expected by the party-ceremony runtime helpers.
func FromCantonParticipant(participant canton.Participant) (Participant, error) {
	var jwt string
	if participant.TokenSource != nil {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return Participant{}, fmt.Errorf("getting token for participant %q: %w", participant.Name, err)
		}
		jwt = token.AccessToken
	}

	return Participant{
		Name:             participant.Name,
		AdminAPIURL:      participant.Endpoints.AdminAPIURL,
		GRPCLedgerAPIURL: participant.Endpoints.GRPCLedgerAPIURL,
		UserID:           participant.UserID,
		AdminJWT:         jwt,
		LedgerJWT:        jwt,
	}, nil
}

// ParticipantUID returns the Canton unique participant identifier used by the
// party-ceremony topology sequences.
func ParticipantUID(ctx context.Context, participant Participant) (string, error) {
	conn, err := dialAdmin(ctx, participant)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return client.NewGRPCClient(conn).GetParticipantUID(ctx)
}

// DiscoverSynchronizerID returns the first connected synchronizer for a
// participant.
func DiscoverSynchronizerID(ctx context.Context, participant Participant) (string, error) {
	conn, err := dialAdmin(ctx, participant)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	synchronizers := participantv30.NewSynchronizerConnectivityServiceClient(conn)
	resp, err := synchronizers.ListConnectedSynchronizers(
		ctx,
		&participantv30.ListConnectedSynchronizersRequest{},
	)
	if err != nil {
		return "", err
	}
	if len(resp.GetConnectedSynchronizers()) == 0 {
		return "", fmt.Errorf("participant %q has no connected synchronizers", participant.Name)
	}

	return resp.GetConnectedSynchronizers()[0].GetSynchronizerId(), nil
}

// NewOnboardingDeps creates real gRPC-backed dependencies for onboarding a
// decentralized party from an existing CLDF Canton participant.
func NewOnboardingDeps(
	ctx context.Context,
	participant Participant,
	lggr logger.Logger,
	confirmer ceremony.Confirmer,
) (ceremony.CantonDeps, func() error, error) {
	conn, err := dialAdmin(ctx, participant)
	if err != nil {
		return ceremony.CantonDeps{}, nil, err
	}

	return ceremony.CantonDeps{
		Client:    client.NewGRPCClient(conn),
		Logger:    lggr,
		Confirmer: confirmer,
	}, conn.Close, nil
}

// NewContractDeployDeps creates real gRPC-backed dependencies for deploying a
// contract through an already-onboarded decentralized party.
func NewContractDeployDeps(
	ctx context.Context,
	participant Participant,
	darLoader ledger.DARLoader,
	damlKeyFingerprint string,
	lggr logger.Logger,
	confirmer ceremony.Confirmer,
) (ledger.ContractDeployDeps, func() error, error) {
	adminConn, err := dialAdmin(ctx, participant)
	if err != nil {
		return ledger.ContractDeployDeps{}, nil, err
	}

	ledgerConn, err := dialLedger(ctx, participant)
	if err != nil {
		_ = adminConn.Close()

		return ledger.ContractDeployDeps{}, nil, err
	}

	vault := cryptoadminv30.NewVaultServiceClient(adminConn)
	signer, err := client.NewVaultSigner(ctx, vault, damlKeyFingerprint)
	if err != nil {
		_ = ledgerConn.Close()
		_ = adminConn.Close()

		return ledger.ContractDeployDeps{}, nil, err
	}

	cleanup := func() error {
		ledgerErr := ledgerConn.Close()
		adminErr := adminConn.Close()
		if ledgerErr != nil {
			return ledgerErr
		}

		return adminErr
	}

	return ledger.ContractDeployDeps{
		AdminClient:  client.NewGRPCClient(adminConn),
		LedgerClient: client.NewGRPCLedgerClient(ledgerConn),
		DARLoader:    darLoader,
		Signer:       signer,
		Logger:       lggr,
		Confirmer:    confirmer,
		UserID:       participant.UserID,
	}, cleanup, nil
}

func dialAdmin(_ context.Context, participant Participant) (*grpc.ClientConn, error) {
	opts, err := dialOptions(participant, false)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(participant.AdminAPIURL, opts...)
}

func dialLedger(_ context.Context, participant Participant) (*grpc.ClientConn, error) {
	opts, err := dialOptions(participant, true)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(participant.GRPCLedgerAPIURL, opts...)
}

func dialOptions(participant Participant, ledgerAPI bool) ([]grpc.DialOption, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	jwt := participant.AdminJWT
	if ledgerAPI && participant.LedgerJWT != "" {
		jwt = participant.LedgerJWT
	}

	if jwt != "" {
		opts = append(opts,
			grpc.WithUnaryInterceptor(jwtUnary(jwt)),
			grpc.WithStreamInterceptor(jwtStream(jwt)),
		)
	} else if ledgerAPI && participant.UserID != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(userIDUnary(participant.UserID)))
	}

	return opts, nil
}

func jwtUnary(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), method, req, reply, cc, opts...)
	}
}

func jwtStream(token string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), desc, cc, method, opts...)
	}
}

func userIDUnary(userID string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if r, ok := req.(*apiv2interactive.PrepareSubmissionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}
		if r, ok := req.(*apiv2interactive.ExecuteSubmissionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}
		if r, ok := req.(*apiv2interactive.ExecuteSubmissionAndWaitForTransactionRequest); ok && r.UserId == "" {
			r.UserId = userID
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
