package devenv

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	adminv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
)

func resolveAuthConfig(tomlJWT, tomlUserID string) commonconfig.AuthConfig {
	jwt := envTrim("CANTON_JWT")
	if jwt == "" {
		jwt = tomlJWT
	}

	authType := envTrim("CANTON_AUTH_TYPE")
	switch {
	case authType == "" && jwt != "":
		authType = commonconfig.AuthTypeInsecureStatic
	case authType == "":
		authType = commonconfig.AuthTypeAuthorizationCode
	}

	authURL := envTrim("CANTON_AUTH_URL")
	clientID := envTrim("CANTON_CLIENT_ID")
	clientSecret := envTrim("CANTON_CLIENT_SECRET")
	if authType != commonconfig.AuthTypeClientCredentials {
		clientSecret = ""
	}

	return commonconfig.AuthConfig{
		Type:         authType,
		UserID:       resolveUserID(tomlUserID),
		JWT:          jwt,
		AuthURL:      authURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func resolveUserID(tomlUserID string) string {
	if userID := envTrim("CANTON_USER_ID"); userID != "" {
		return userID
	}

	return tomlUserID
}

func resolvePartyID() string {
	return envTrim("CANTON_PARTY_ID")
}

func resolveGRPCLedgerURL(tomlURL string) string {
	if url := envTrim("CANTON_GRPC_URL"); url != "" {
		return url
	}

	return tomlURL
}

func resolveValidatorAPIURL(tomlURL string) string {
	if url := envTrim("CANTON_VALIDATOR_API_URL"); url != "" {
		return url
	}

	return tomlURL
}

func newAuthProvider(ctx context.Context, authCfg commonconfig.AuthConfig) (authentication.Provider, error) {
	if authCfg.Type == commonconfig.AuthTypeAuthorizationCode && authCfg.UserID == "" {
		return authorizationcode.NewDiscoveryProvider(ctx, authCfg.AuthURL, authCfg.ClientID)
	}

	return authCfg.NewProvider(ctx)
}

func userIDFromToken(ctx context.Context, provider authentication.Provider) (string, error) {
	_ = ctx

	token, err := provider.TokenSource().Token()
	if err != nil {
		return "", fmt.Errorf("get oauth token: %w", err)
	}

	sub, err := jwtSubject(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("extract user id from token sub: %w", err)
	}
	if sub == "" {
		return "", fmt.Errorf("token sub claim is empty")
	}

	return sub, nil
}

func jwtSubject(accessToken string) (string, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode JWT payload: %w", err)
	}

	var claims struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("parse JWT payload: %w", err)
	}

	return claims.Sub, nil
}

func resolveParticipantParty(ctx context.Context, conn *grpc.ClientConn, userID, partyID string) (string, error) {
	if party := strings.TrimSpace(partyID); party != "" {
		return party, nil
	}

	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("party ID not preset and user ID unset")
	}

	userResp, err := adminv2.NewUserManagementServiceClient(conn).GetUser(ctx, &adminv2.GetUserRequest{UserId: userID})
	if err != nil {
		return "", fmt.Errorf("get user %s: %w", userID, err)
	}

	party := userResp.GetUser().GetPrimaryParty()
	if party == "" {
		return "", fmt.Errorf("no primary party found for user %s", userID)
	}

	return party, nil
}

func endpointsNonEmpty(jsonURL, grpcURL, adminURL, validatorURL string) bool {
	return strings.TrimSpace(jsonURL) != "" ||
		strings.TrimSpace(grpcURL) != "" ||
		strings.TrimSpace(adminURL) != "" ||
		strings.TrimSpace(validatorURL) != ""
}

func envTrim(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
