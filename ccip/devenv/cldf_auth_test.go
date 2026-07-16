package devenv

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

// validJWT is a well-formed JWT with sub "1234567890".
const validJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

func clearCantonEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"CANTON_AUTH_TYPE",
		"CANTON_AUTH_URL",
		"CANTON_CLIENT_ID",
		"CANTON_CLIENT_SECRET",
		"CANTON_JWT",
		"CANTON_USER_ID",
		"CANTON_PARTY_ID",
		"CANTON_GRPC_URL",
		"CANTON_VALIDATOR_API_URL",
	} {
		t.Setenv(key, "")
	}
}

func TestEnvTrim(t *testing.T) {
	clearCantonEnv(t)

	require.Empty(t, envTrim("CANTON_PARTY_ID"))

	t.Setenv("CANTON_PARTY_ID", "party-1")
	require.Equal(t, "party-1", envTrim("CANTON_PARTY_ID"))

	t.Setenv("CANTON_PARTY_ID", "  party-2  ")
	require.Equal(t, "party-2", envTrim("CANTON_PARTY_ID"))
}

func TestResolveAuthConfig(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		tomlJWT    string
		tomlUserID string
		want       commonconfig.AuthConfig
	}{
		{
			name:    "toml jwt defaults to insecureStatic",
			tomlJWT: validJWT,
			want: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeInsecureStatic,
				JWT:  validJWT,
			},
		},
		{
			name: "no jwt defaults to authorizationCode",
			want: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeAuthorizationCode,
			},
		},
		{
			name: "env overrides toml",
			env: map[string]string{
				"CANTON_AUTH_TYPE":     commonconfig.AuthTypeClientCredentials,
				"CANTON_JWT":           "env-jwt",
				"CANTON_USER_ID":       "env-user",
				"CANTON_AUTH_URL":      "https://auth.example.com/",
				"CANTON_CLIENT_ID":     "env-client",
				"CANTON_CLIENT_SECRET": "env-secret",
			},
			tomlJWT:    validJWT,
			tomlUserID: "toml-user",
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeClientCredentials,
				UserID:       "env-user",
				JWT:          "env-jwt",
				AuthURL:      "https://auth.example.com/",
				ClientID:     "env-client",
				ClientSecret: "env-secret",
			},
		},
		{
			name: "client secret cleared for authorizationCode",
			env: map[string]string{
				"CANTON_AUTH_TYPE":     commonconfig.AuthTypeAuthorizationCode,
				"CANTON_AUTH_URL":      "https://auth.example.com/",
				"CANTON_CLIENT_ID":     "env-client",
				"CANTON_CLIENT_SECRET": "should-not-apply",
			},
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeAuthorizationCode,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "env-client",
				ClientSecret: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCantonEnv(t)
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			got := resolveAuthConfig(tt.tomlJWT, tt.tomlUserID)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestResolveUserID(t *testing.T) {
	clearCantonEnv(t)

	require.Equal(t, "toml-user", resolveUserID("toml-user"))

	t.Setenv("CANTON_USER_ID", "env-user")
	require.Equal(t, "env-user", resolveUserID("toml-user"))
}

func TestResolvePartyID(t *testing.T) {
	clearCantonEnv(t)

	require.Empty(t, resolvePartyID())

	t.Setenv("CANTON_PARTY_ID", "party::1220abc")
	require.Equal(t, "party::1220abc", resolvePartyID())
}

func TestResolveGRPCLedgerURL(t *testing.T) {
	clearCantonEnv(t)

	require.Equal(t, "toml-host:443", resolveGRPCLedgerURL("toml-host:443"))

	t.Setenv("CANTON_GRPC_URL", "env-host:443")
	require.Equal(t, "env-host:443", resolveGRPCLedgerURL("toml-host:443"))
}

func TestResolveValidatorAPIURL(t *testing.T) {
	clearCantonEnv(t)

	require.Equal(t, "https://toml.example/validator/", resolveValidatorAPIURL("https://toml.example/validator/"))

	t.Setenv("CANTON_VALIDATOR_API_URL", "https://env.example/validator/")
	require.Equal(t, "https://env.example/validator/", resolveValidatorAPIURL("https://toml.example/validator/"))
}

func TestJWTSubject(t *testing.T) {
	t.Parallel()

	sub, err := jwtSubject(validJWT)
	require.NoError(t, err)
	require.Equal(t, "1234567890", sub)
}

func TestUserIDFromToken(t *testing.T) {
	t.Parallel()

	userID, err := userIDFromToken(context.Background(), testAuthProvider{token: validJWT})
	require.NoError(t, err)
	require.Equal(t, "1234567890", userID)
}

func TestEndpointsNonEmpty(t *testing.T) {
	t.Parallel()

	require.False(t, endpointsNonEmpty("", "", "", ""))
	require.True(t, endpointsNonEmpty("", "grpc:443", "", ""))
}

func TestResolveAuthConfig_envJWTOverridesToml(t *testing.T) {
	clearCantonEnv(t)

	t.Setenv("CANTON_JWT", validJWT)
	got := resolveAuthConfig("toml-jwt-should-lose", "")
	require.Equal(t, validJWT, got.JWT)
	require.Equal(t, commonconfig.AuthTypeInsecureStatic, got.Type)
}

type testAuthProvider struct {
	token string
}

func (p testAuthProvider) TokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.token})
}

func (p testAuthProvider) TransportCredentials() credentials.TransportCredentials {
	return authentication.NewInsecureStaticProvider(p.token).TransportCredentials()
}

func (p testAuthProvider) PerRPCCredentials() credentials.PerRPCCredentials {
	return oauth.TokenSource{TokenSource: p.TokenSource()}
}
