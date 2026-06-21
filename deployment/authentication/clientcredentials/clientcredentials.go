// Package clientcredentials provides OAuth2 client credentials flow authentication for gRPC connections.
// It implements the client credentials grant type as defined in RFC 6749, Section 4.4.
// This flow is designed for machine-to-machine authentication where the client can securely maintain
// a client secret, making it ideal for server-to-server communication and CI/CD environments.
//
// See: https://datatracker.ietf.org/doc/html/rfc6749#section-4.4
package clientcredentials

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	authentication_oauth2 "github.com/smartcontractkit/chainlink-canton/deployment/authentication"
)

var _ authentication.Provider = Provider{}

// Provider implements authentication.Provider using the OAuth2 client credentials flow for gRPC.
// Tokens are automatically refreshed as needed
// using the configured token source.
type Provider struct {
	tokenSource          oauth.TokenSource
	transportCredentials credentials.TransportCredentials
}

type clientCredentialsProviderConfig struct {
	scopes               []string
	audience             string
	transportCredentials credentials.TransportCredentials
}

func defaultClientCredentialsProviderConfig() *clientCredentialsProviderConfig {
	return &clientCredentialsProviderConfig{
		scopes:   []string{"daml_ledger_api"},
		audience: "",
		transportCredentials: credentials.NewTLS(
			&tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		),
	}
}

// ProviderOption configures the client credentials Provider using the functional options pattern.
// Options allow customization of scopes, audience, and transport credentials without breaking API compatibility.
type ProviderOption func(*clientCredentialsProviderConfig)

// WithScopes configures the Provider to request access tokens with the given scopes.
// Scopes define the level of access requested from the authorization server.
// The default scope is "daml_ledger_api" for Canton ledger API access.
//
// Example:
//
//	WithScopes("daml_ledger_api", "read:users")
func WithScopes(scopes ...string) ProviderOption {
	return func(config *clientCredentialsProviderConfig) {
		config.scopes = scopes
	}
}

// WithAudience configures the Provider to request access tokens with the given audience.
// The audience identifies the intended recipient of the issued access token (typically
// the resource server / API that will consume the token, e.g. the Canton ledger API).
//
// When configured, it is sent to the authorization server as an "audience" parameter in
// the token request body. This is an Auth0-specific extension: Auth0 uses the "audience"
// parameter to select which API (and therefore which value of the JWT's "aud" claim,
// defined in RFC 7519 Section 4.1.3) the issued access token should target.
//
// NOTE: The "audience" parameter is NOT part of the OAuth2 specification, and most other
// authorization servers do not honor it:
//   - Okta binds the audience to the Authorization Server itself (one fixed value per AS) and
//     explicitly does not support dynamic audience switching via a request parameter, nor
//     RFC 8707 resource indicators. To use a different audience on Okta, point tokenURL at a
//     different Authorization Server.
//   - Keycloak determines the audience server-side via audience mappers on client scopes.
//
// This option therefore only has an effect with Auth0 (or an authorization server that
// explicitly emulates Auth0's behavior).
//
// Note: Auth0's tenant-level "Resource Parameter Compatibility Profile" (which makes Auth0
// honor the standardized RFC 8707 "resource" parameter as an alternative to "audience") does
// NOT cover the client credentials grant — its supported flows are limited to the
// authorization code flow, PAR, JAR, CIBA, and refresh token grants.
//
// If no audience is configured (the default), no "audience" parameter is sent and the
// authorization server determines the audience based on its own policy / client configuration.
//
// Example:
//
//	WithAudience("https://ledger.example.com")
func WithAudience(audience string) ProviderOption {
	return func(config *clientCredentialsProviderConfig) {
		config.audience = audience
	}
}

// WithTransportCredentials configures the Provider to use the given transport credentials for gRPC connections.
// This allows customization of TLS settings, including certificate verification and minimum TLS version.
// The default transport credentials use TLS 1.2 or higher.
//
// Example:
//
//	WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))
func WithTransportCredentials(credentials credentials.TransportCredentials) ProviderOption {
	return func(config *clientCredentialsProviderConfig) {
		config.transportCredentials = credentials
	}
}

// NewDiscoveryProvider creates a provider using OAuth2 Authorization Server Metadata discovery
// to automatically locate the token endpoint. This implements the discovery mechanism defined
// in RFC 8414, querying the .well-known/oauth-authorization-server endpoint.
//
// This is the recommended approach when the authorization server supports metadata discovery,
// as it eliminates the need to manually specify endpoint URLs and reduces configuration errors.
//
// Parameters:
//   - ctx: Context for the initial metadata + token request
//   - authorizationServerURL: The base URL of the authorization server (e.g., "https://auth.example.com")
//   - clientID: The OAuth2 client identifier issued by the authorization server
//   - clientSecret: The OAuth2 client secret (should be kept secure)
//   - options: Optional configuration parameters (scopes, transport credentials)
//
// See: https://datatracker.ietf.org/doc/html/rfc8414
func NewDiscoveryProvider(ctx context.Context, authorizationServerURL, clientID, clientSecret string, options ...ProviderOption) (*Provider, error) {
	metadata, err := authentication_oauth2.GetAuthorizationServerMetadata(ctx, authorizationServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization server metadata: %w", err)
	}

	return NewProvider(ctx, metadata.TokenEndpoint, clientID, clientSecret, options...)
}

// NewProvider creates a provider that fetches tokens using the OAuth2 client credentials flow
// as defined in RFC 6749, Section 4.4. This constructor requires the token endpoint URL to be
// explicitly provided, making it suitable for environments where the endpoint is known in advance.
//
// The provider automatically handles token acquisition and refresh. Tokens are requested from the
// authorization server using the client credentials (client ID and secret) and are cached until
// they expire. The token source will automatically refresh expired tokens before making requests.
//
// This is ideal for CI/CD pipelines, automated testing, and server-to-server communication where
// the client credentials and token endpoint URL are available as configuration or environment variables.
// For authorization servers that support OAuth2 discovery (RFC 8414), consider using NewDiscoveryProvider
// instead to automatically locate the token endpoint.
//
// Parameters:
//   - ctx: Context for the initial token request
//   - tokenURL: The OAuth2 token endpoint URL (e.g., "https://auth.example.com/v1/token")
//   - clientID: The OAuth2 client identifier issued by the authorization server
//   - clientSecret: The OAuth2 client secret (should be kept secure and never logged)
//   - options: Optional configuration parameters (scopes, transport credentials)
//
// Returns an error if any required parameter is empty or if the provider cannot be initialized.
//
// See: https://datatracker.ietf.org/doc/html/rfc6749#section-4.4
func NewProvider(ctx context.Context, tokenURL, clientID, clientSecret string, options ...ProviderOption) (*Provider, error) {
	cfg := defaultClientCredentialsProviderConfig()
	for _, option := range options {
		option(cfg)
	}

	if tokenURL == "" {
		return nil, fmt.Errorf("tokenURL cannot be empty")
	}
	if clientID == "" {
		return nil, fmt.Errorf("clientID cannot be empty")
	}
	if clientSecret == "" {
		return nil, fmt.Errorf("clientSecret cannot be empty")
	}

	oauthCfg := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       cfg.scopes,
	}

	// If an audience is set, add it to the token request parameters.
	if cfg.audience != "" {
		oauthCfg.EndpointParams = map[string][]string{
			"audience": {cfg.audience},
		}
	}

	refreshCtx := context.WithoutCancel(ctx)
	tokenSource := oauthCfg.TokenSource(refreshCtx)

	return &Provider{
		tokenSource:          oauth.TokenSource{TokenSource: tokenSource},
		transportCredentials: cfg.transportCredentials,
	}, nil
}

// TokenSource returns the OAuth2 token source for obtaining access tokens.
// The token source automatically handles token refresh and caching, ensuring that
// valid tokens are provided for each request without manual intervention.
//
// Returns the underlying oauth2.TokenSource that can be used independently of gRPC
// for other OAuth2-authenticated requests.
func (p Provider) TokenSource() oauth2.TokenSource {
	return p.tokenSource.TokenSource
}

// TransportCredentials returns the transport credentials for establishing gRPC connections.
func (p Provider) TransportCredentials() credentials.TransportCredentials {
	return p.transportCredentials
}

// PerRPCCredentials returns the per-RPC credentials that attach OAuth2 tokens to each gRPC call.
// The credentials implement RFC 6750 Bearer Token usage, adding the access token to the
// Authorization header of each RPC using the format: "Authorization: Bearer <token>".
//
// See: https://datatracker.ietf.org/doc/html/rfc6750
func (p Provider) PerRPCCredentials() credentials.PerRPCCredentials {
	return p.tokenSource
}
