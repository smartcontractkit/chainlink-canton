// Package authorizationcode provides OAuth2 authorization code flow authentication for gRPC connections.
// It implements the authorization code grant type as defined in RFC 6749, Section 4.1, and REQUIRES
// PKCE (RFC 7636) with the S256 challenge method to protect the authorization code exchange.
//
// This flow is designed for interactive user authentication where a browser login is required.
// It starts a local callback server to receive the authorization code and exchanges it for tokens.
//
// See: https://datatracker.ietf.org/doc/html/rfc6749#section-4.1
// See: https://datatracker.ietf.org/doc/html/rfc7636
package authorizationcode

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/icza/gox/osx"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	authentication_oauth2 "github.com/smartcontractkit/chainlink-canton/deployment/authentication"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
)

var _ authentication.Provider = Provider{}

type Provider struct {
	// tokenSource provides per-RPC OAuth2 tokens and handles refresh as needed.
	tokenSource          oauth.TokenSource
	transportCredentials credentials.TransportCredentials
}

type authorizationCodeProviderConfig struct {
	scopes               []string
	audience             string
	transportCredentials credentials.TransportCredentials
	callbackURL          string
	openBrowser          bool
	timeout              time.Duration
}

func defaultAuthorizationCodeProviderConfig() *authorizationCodeProviderConfig {
	return &authorizationCodeProviderConfig{
		scopes:   []string{"openid", "daml_ledger_api"},
		audience: "",
		transportCredentials: credentials.NewTLS(
			&tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		),
		callbackURL: "http://localhost:8400/callback",
		openBrowser: true,
	}
}

// ProviderOption configures the authorization code Provider using the functional options pattern.
// Options allow customization of scopes, audience, transport credentials, callback URL, browser
// behavior, and overall flow timeout.
type ProviderOption func(*authorizationCodeProviderConfig)

// WithScopes configures the Provider to request access tokens with the given scopes.
// Scopes define the level of access requested from the authorization server.
// The default scope is "daml_ledger_api" for Canton ledger API access.
//
// Example:
//
//	WithScopes("daml_ledger_api", "read:users")
func WithScopes(scopes ...string) ProviderOption {
	return func(config *authorizationCodeProviderConfig) {
		config.scopes = scopes
	}
}

// WithAudience configures the Provider to request access tokens with the given audience.
// The audience identifies the intended recipient of the issued access token (typically
// the resource server / API that will consume the token, e.g. the Canton ledger API).
//
// When configured, it is appended to the authorization request as an "audience" query
// parameter. This is an Auth0-specific extension: Auth0 uses the "audience" parameter to
// select which API (and therefore which value of the JWT's "aud" claim, defined in RFC 7519
// Section 4.1.3) the issued access token should target.
//
// NOTE: The "audience" parameter is NOT part of the OAuth2 specification, and most other
// authorization servers do not honor it:
//   - Okta binds the audience to the Authorization Server itself (one fixed value per AS) and
//     explicitly does not support dynamic audience switching via a request parameter, nor
//     RFC 8707 resource indicators. To use a different audience on Okta, point authURL/tokenURL
//     at a different Authorization Server.
//   - Keycloak determines the audience server-side via audience mappers on client scopes.
//
// This option therefore only has an effect with Auth0 (or an authorization server that
// explicitly emulates Auth0's behavior).
//
// As of 2025, Auth0 also offers a tenant-level "Resource Parameter Compatibility Profile"
// (Dashboard → Settings → Advanced) which, when enabled, makes Auth0 honor the standardized
// RFC 8707 "resource" request parameter as an alternative way to set the token's audience for
// the authorization code flow (and PAR/JAR/CIBA/refresh token grants). When both "resource"
// and "audience" are sent, Auth0 still gives "audience" precedence.
//
// If no audience is configured (the default), no "audience" parameter is sent and the
// authorization server determines the audience based on its own policy / client configuration.
//
// Example:
//
//	WithAudience("https://ledger.example.com")
func WithAudience(audience string) ProviderOption {
	return func(config *authorizationCodeProviderConfig) {
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
	return func(config *authorizationCodeProviderConfig) {
		config.transportCredentials = credentials
	}
}

// WithCallbackURL configures the local redirect URI used by the authorization server.
// The callback URL must be reachable by the browser and is expected to be a localhost URL.
func WithCallbackURL(callbackURL string) ProviderOption {
	return func(config *authorizationCodeProviderConfig) {
		config.callbackURL = callbackURL
	}
}

// WithOpenBrowser controls whether the default browser is opened automatically.
// When disabled, the authorization URL is printed for manual copy/paste.
func WithOpenBrowser(openBrowser bool) ProviderOption {
	return func(config *authorizationCodeProviderConfig) {
		config.openBrowser = openBrowser
	}
}

// WithTimeout configures a timeout for the overall authorization flow, including callback receipt.
func WithTimeout(timeout time.Duration) ProviderOption {
	return func(config *authorizationCodeProviderConfig) {
		config.timeout = timeout
	}
}

// NewDiscoveryProvider creates a provider using OAuth2 Authorization Server Metadata discovery
// to automatically locate the authorization and token endpoints. This implements the discovery
// mechanism defined in RFC 8414, querying the .well-known/oauth-authorization-server endpoint.
//
// PKCE with the S256 challenge method is REQUIRED; an error is returned if the server does not
// advertise support for S256.
//
// Parameters:
//   - ctx: Context for metadata discovery and the token exchange
//   - authorizationServerURL: The base URL of the authorization server (e.g., "https://auth.example.com")
//   - clientID: The OAuth2 client identifier issued by the authorization server
//   - options: Optional configuration parameters (scopes, transport credentials, callback URL, timeout, etc.)
//
// Returns an error if discovery fails or if the server does not support S256 PKCE challenges.
//
// See: https://datatracker.ietf.org/doc/html/rfc8414
func NewDiscoveryProvider(ctx context.Context, authorizationServerURL, clientID string, options ...ProviderOption) (*Provider, error) {
	metadata, err := authentication_oauth2.GetAuthorizationServerMetadata(ctx, authorizationServerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization server metadata: %w", err)
	}

	if !slices.Contains(metadata.CodeChallengeMethodsSupported, "S256") {
		// S256 PKCE is required for this flow; fail fast if unsupported.
		return nil, fmt.Errorf("authorization server does not support S256 PKCE challenges")
	}

	return NewProvider(ctx, metadata.AuthorizationEndpoint, metadata.TokenEndpoint, clientID, options...)
}

// NewProvider creates a provider that performs the OAuth2 authorization code flow with PKCE (S256).
// It starts a local callback server to receive the authorization code and exchanges it for tokens.
//
// PKCE with the S256 challenge method is REQUIRED for this flow.
//
// Parameters:
//   - ctx: Context for the overall flow; can be configured with a timeout via WithTimeout
//   - authURL: The OAuth2 authorization endpoint URL
//   - tokenURL: The OAuth2 token endpoint URL
//   - clientID: The OAuth2 client identifier issued by the authorization server
//   - options: Optional configuration parameters (scopes, transport credentials, callback URL, etc.)
//
// Returns an error if any required parameter is empty or if the callback server cannot be started.
//
// See: https://datatracker.ietf.org/doc/html/rfc6749#section-4.1
func NewProvider(ctx context.Context, authURL, tokenURL, clientID string, options ...ProviderOption) (*Provider, error) {
	cfg := defaultAuthorizationCodeProviderConfig()
	for _, option := range options {
		option(cfg)
	}

	if authURL == "" {
		return nil, fmt.Errorf("authURL cannot be empty")
	}
	if tokenURL == "" {
		return nil, fmt.Errorf("tokenURL cannot be empty")
	}
	if clientID == "" {
		return nil, fmt.Errorf("clientID cannot be empty")
	}

	if cfg.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}

	callbackURL, err := url.Parse(cfg.callbackURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse callback URL: %w", err)
	}

	oauthCfg := &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: callbackURL.String(),
		Scopes:      cfg.scopes,
		Endpoint:    oauth2.Endpoint{AuthURL: authURL, TokenURL: tokenURL},
	}

	// Generate cryptographically secure random state
	state := oauth2.GenerateVerifier()

	// Use built-in S256ChallengeOption for PKCE
	verifier := oauth2.GenerateVerifier()
	authCodeOpts := []oauth2.AuthCodeOption{oauth2.S256ChallengeOption(verifier)}
	// If an audience is configured, send it as an "audience" query parameter on the
	// authorization request. This is an Auth0-specific extension; see the docs on
	// WithAudience for details.
	if cfg.audience != "" {
		authCodeOpts = append(authCodeOpts, oauth2.SetAuthURLParam("audience", cfg.audience))
	}
	authCodeURL := oauthCfg.AuthCodeURL(state, authCodeOpts...)

	callbackChan := make(chan *oauth2.Token)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc(callbackURL.Path, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := q.Get("code")
		receivedState := q.Get("state")

		if receivedState != state {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			fmt.Println("ERROR: Invalid state parameter received")

			return
		}
		if code == "" {
			http.Error(w, "No code parameter received", http.StatusBadRequest)
			fmt.Println("ERROR: No code parameter received")

			return
		}

		// Use built-in VerifierOption for PKCE
		token, err := oauthCfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			http.Error(w, "Token exchange failed", http.StatusInternalServerError)
			fmt.Printf("ERROR: Token exchange failed: %v\n", err)

			return
		}

		callbackChan <- token

		// HTML response for the browser
		html := `
			<!DOCTYPE html>
			<html>
			<head><title>Authentication Complete</title></head>
			<body style="font-family: sans-serif; text-align: center; padding: 40px;">
				<h1>Authentication complete!</h1>
				<p>You can safely close this window.</p>
			</body>
			</html>
		`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	})

	// Start callback server
	server := http.Server{
		Addr:              callbackURL.Host,
		Handler:           serveMux,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	// Create listener to fail fast if port is unavailable
	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", server.Addr)
	if err != nil {
		return nil, fmt.Errorf("creating listener: %w", err)
	}
	fmt.Printf("Waiting for authentication on %v\n", callbackURL.String())
	serverErr := make(chan error, 1)
	go func(chan<- error) {
		serverErr <- server.Serve(listener)
	}(serverErr)

	if cfg.openBrowser {
		fmt.Println("Attempting to open your default browser.\nIf the browser does not open, open the following URL:")
		fmt.Println(authCodeURL)
		_ = osx.OpenDefault(authCodeURL)
	} else {
		fmt.Println("Visit the following URL:")
		fmt.Println(authCodeURL)
	}

	select {
	case err := <-serverErr:
		_ = server.Shutdown(ctx)
		return nil, fmt.Errorf("callback server error: %w", err)
	case token := <-callbackChan:
		fmt.Println("Authentication completed")
		tokenSource := oauthCfg.TokenSource(ctx, token)

		return &Provider{
			tokenSource:          oauth.TokenSource{TokenSource: tokenSource},
			transportCredentials: cfg.transportCredentials,
		}, server.Shutdown(ctx)
	case <-ctx.Done():
		_ = server.Shutdown(ctx)
		return nil, ctx.Err()
	}
}

// TokenSource returns the OAuth2 token source for obtaining access tokens.
// The token source automatically handles token refresh and caching.
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
