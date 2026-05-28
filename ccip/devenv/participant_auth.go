package devenv

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

const (
	envCantonAuthType          = "CANTON_AUTH_TYPE"
	envCantonAuthURL           = "CANTON_AUTH_URL"
	envCantonOAuthClientID     = "CANTON_OAUTH_CLIENT_ID"
	envCantonOAuthClientSecret = "CANTON_OAUTH_CLIENT_SECRET"
	envOnchainCantonJWT        = "ONCHAIN_CANTON_JWT_TOKEN"
	envClientID                = "CLIENT_ID"
	envClientSecret            = "CLIENT_SECRET"
)

// buildParticipantAuthConfig returns an AuthConfig with an explicit type.
// Devenv participants with a JWT and a local (non-TLS) gRPC endpoint use insecureStatic.
// Real-chain connections default to clientCredentials from env (same as canton-login --ci).
// Static auth is used only when explicitly requested via CANTON_AUTH_TYPE=static.
func buildParticipantAuthConfig(participant blockchain.CantonParticipantEndpoints) (commonconfig.AuthConfig, error) {
	if isDevenvParticipant(participant) {
		return commonconfig.AuthConfig{
			Type: commonconfig.AuthTypeInsecureStatic,
			JWT:  participant.JWT,
		}, nil
	}

	authType := strings.TrimSpace(os.Getenv(envCantonAuthType))
	if authType == "" {
		authType = commonconfig.AuthTypeClientCredentials
	}

	switch authType {
	case commonconfig.AuthTypeInsecureStatic:
		jwt := participant.JWT
		if jwt == "" {
			return commonconfig.AuthConfig{}, fmt.Errorf("insecureStatic auth requires a JWT on the participant config")
		}

		return commonconfig.AuthConfig{
			Type: commonconfig.AuthTypeInsecureStatic,
			JWT:  jwt,
		}, nil

	case commonconfig.AuthTypeStatic:
		jwt := participant.JWT
		if override := strings.TrimSpace(os.Getenv(envOnchainCantonJWT)); override != "" {
			jwt = override
		}
		if jwt == "" {
			return commonconfig.AuthConfig{}, fmt.Errorf("static auth requires jwt on participant config or %s", envOnchainCantonJWT)
		}

		return commonconfig.AuthConfig{
			Type: commonconfig.AuthTypeStatic,
			JWT:  jwt,
		}, nil

	case commonconfig.AuthTypeAuthorizationCode:
		return commonconfig.AuthConfig{
			Type:     commonconfig.AuthTypeAuthorizationCode,
			AuthURL:  strings.TrimSpace(os.Getenv(envCantonAuthURL)),
			ClientID: oauthClientID(),
		}, nil

	case commonconfig.AuthTypeClientCredentials:
		return commonconfig.AuthConfig{
			Type:         commonconfig.AuthTypeClientCredentials,
			AuthURL:      strings.TrimSpace(os.Getenv(envCantonAuthURL)),
			ClientID:     oauthClientID(),
			ClientSecret: oauthClientSecret(),
		}, nil

	default:
		return commonconfig.AuthConfig{}, fmt.Errorf("unsupported %s: %q", envCantonAuthType, authType)
	}
}

func oauthClientID() string {
	return firstNonEmpty(
		strings.TrimSpace(os.Getenv(envCantonOAuthClientID)),
		strings.TrimSpace(os.Getenv(envClientID)),
	)
}

func oauthClientSecret() string {
	return firstNonEmpty(
		strings.TrimSpace(os.Getenv(envCantonOAuthClientSecret)),
		strings.TrimSpace(os.Getenv(envClientSecret)),
	)
}

func isDevenvParticipant(participant blockchain.CantonParticipantEndpoints) bool {
	if strings.TrimSpace(participant.JWT) == "" {
		return false
	}

	return isLocalDevenvEndpoint(participant.GRPCLedgerAPIURL)
}

// isLocalDevenvEndpoint reports whether the gRPC ledger URL is a local devenv endpoint
// (plain gRPC, not TLS on :443).
func isLocalDevenvEndpoint(grpcURL string) bool {
	grpcURL = strings.TrimSpace(grpcURL)
	if grpcURL == "" {
		return false
	}

	hostPort := grpcURL
	if idx := strings.LastIndex(hostPort, "/"); idx != -1 {
		hostPort = hostPort[idx+1:]
	}

	return !strings.HasSuffix(hostPort, ":443")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}

// participantAuthProvider validates auth config and builds an authentication.Provider.
func participantAuthProvider(ctx context.Context, auth commonconfig.AuthConfig) (authentication.Provider, error) {
	if err := auth.Validate(); err != nil {
		return nil, err
	}

	return auth.NewProvider(ctx)
}
