package commonconfig

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"

	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"
)

// Supported authentication types for Canton participant APIs.
const (
	AuthTypeStatic            = "static"
	AuthTypeClientCredentials = "clientCredentials"
	AuthTypeAuthorizationCode = "authorizationCode"
)

// AuthConfig configures authentication for a Canton participant endpoint.
type AuthConfig struct {
	// Type selects the auth scheme: "static", "clientCredentials", or "authorizationCode".
	// Defaults to "static" when omitted (backward compatible).
	Type string `toml:"type" validate:"required,oneof=static clientCredentials authorizationCode"`

	// UserID is the user ID for the authentication. Required for all auth types except for static.
	UserID string `toml:"user_id" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode"`

	// JWT is a pre-obtained token. Required when Type is "static" or empty.
	// optional_jwt is a custom validator registered in Validate(); revive's struct-tag only knows built-in validator options.
	JWT string `toml:"jwt,omitempty" validate:"required_if=Type static,excluded_unless=Type static,omitempty,jwt"`

	// AuthURL is the OIDC authorization server base URL. Required for clientCredentials and authorizationCode.
	// optional_url is a custom validator registered in Validate(); revive's struct-tag only knows built-in validator options.
	AuthURL string `toml:"auth_url,omitempty" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode,omitempty,url"`

	// ClientID is the OAuth2 client identifier. Required for clientCredentials and authorizationCode.
	ClientID string `toml:"client_id,omitempty" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode"`

	// ClientSecret is the OAuth2 client secret. Required for clientCredentials only.
	ClientSecret string `toml:"client_secret,omitempty" validate:"required_if=Type clientCredentials,excluded_unless=Type clientCredentials"`

	// InsecureTransport disables TLS transport security when using static auth.
	// Default is false (secure transport). Set to true for local Canton setups where TLS is not available.
	InsecureTransport bool `toml:"insecure_transport,omitempty"`
}

func (a *AuthConfig) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(a); err != nil {
		return fmt.Errorf("auth config: %w", err)
	}

	return nil
}

// NewProvider builds an authentication.Provider from this config.
// It validates the config (including JWT/URL format) before building the provider.
func (a *AuthConfig) NewProvider(ctx context.Context) (authentication.Provider, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	authType := a.Type
	if authType == "" {
		authType = AuthTypeStatic
	}

	switch authType {
	case AuthTypeStatic:
		if a.JWT == "" {
			return nil, fmt.Errorf("static auth requires a JWT token (set auth.jwt)")
		}
		if a.InsecureTransport {
			return authentication.NewInsecureStaticProvider(a.JWT), nil
		}

		return authentication.NewStaticProvider(a.JWT), nil

	case AuthTypeClientCredentials:
		if a.AuthURL == "" || a.ClientID == "" || a.ClientSecret == "" {
			return nil, fmt.Errorf("clientCredentials auth requires auth_url, client_id, and client_secret")
		}

		return clientcredentials.NewDiscoveryProvider(ctx, a.AuthURL, a.ClientID, a.ClientSecret)

	case AuthTypeAuthorizationCode:
		if a.AuthURL == "" || a.ClientID == "" {
			return nil, fmt.Errorf("authorizationCode auth requires auth_url and client_id")
		}

		return authorizationcode.NewDiscoveryProvider(ctx, a.AuthURL, a.ClientID)

	default:
		return nil, fmt.Errorf("unsupported auth type: %q (expected static, clientCredentials, or authorizationCode)", authType)
	}
}
