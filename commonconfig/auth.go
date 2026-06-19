package commonconfig

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"

	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"
	"github.com/smartcontractkit/chainlink-canton/deployment/authentication/static"
)

// Supported authentication types for Canton participant APIs.
const (
	// AuthTypeStatic should be used when we have a JWT token.
	AuthTypeStatic = "static"

	// AuthTypeInsecureStatic should be used when we have a JWT token but we're not using TLS for transport security.
	AuthTypeInsecureStatic = "insecureStatic"

	// AuthTypeClientCredentials should be used when we're using client credentials authentication.
	AuthTypeClientCredentials = "clientCredentials"

	// AuthTypeAuthorizationCode should be used when we're using authorization code authentication.
	AuthTypeAuthorizationCode = "authorizationCode"
)

// AuthConfig configures authentication for a Canton participant endpoint.
type AuthConfig struct {
	// Type selects the auth scheme: "static", "insecureStatic", "clientCredentials", or "authorizationCode".
	// Defaults to "static" when omitted (backward compatible).
	Type string `toml:"type" validate:"required,oneof=static insecureStatic clientCredentials authorizationCode"`

	// InsecureSkipVerify configures the authentication provider to disable TLS certification validation.
	InsecureSkipVerify bool `toml:"insecure_skip_verify,omitempty"`

	// UserID is the user ID for the authentication. Required for clientCredentials and authorizationCode only.
	UserID string `toml:"user_id" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode"`

	// JWT is a pre-obtained token. Required when Type is "static" or "insecureStatic".
	JWT string `toml:"jwt,omitempty" validate:"required_if=Type static,required_if=Type insecureStatic,excluded_unless=Type static|excluded_unless=Type insecureStatic,omitempty,jwt"`

	// AuthURL is the OIDC authorization server base URL. Required for clientCredentials and authorizationCode.
	AuthURL string `toml:"auth_url,omitempty" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode,omitempty,url"`

	// ClientID is the OAuth2 client identifier. Required for clientCredentials and authorizationCode.
	ClientID string `toml:"client_id,omitempty" validate:"required_if=Type clientCredentials,required_if=Type authorizationCode"`

	// ClientSecret is the OAuth2 client secret. Required for clientCredentials only.
	ClientSecret string `toml:"client_secret,omitempty" validate:"required_if=Type clientCredentials,excluded_unless=Type clientCredentials"`
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
	case AuthTypeInsecureStatic:
		if a.JWT == "" {
			return nil, fmt.Errorf("insecureStatic auth requires a JWT token (set auth.jwt)")
		}

		return authentication.NewInsecureStaticProvider(a.JWT), nil

	case AuthTypeStatic:
		if a.JWT == "" {
			return nil, fmt.Errorf("static auth requires a JWT token (set auth.jwt)")
		}

		var options []static.ProviderOption
		if a.InsecureSkipVerify {
			options = append(options, static.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})))
		}

		return static.NewStaticProvider(a.JWT, options...), nil

	case AuthTypeClientCredentials:
		if a.AuthURL == "" || a.ClientID == "" || a.ClientSecret == "" {
			return nil, fmt.Errorf("clientCredentials auth requires auth_url, client_id, and client_secret")
		}

		var options []clientcredentials.ProviderOption
		if a.InsecureSkipVerify {
			options = append(options, clientcredentials.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})))
		}

		return clientcredentials.NewDiscoveryProvider(ctx, a.AuthURL, a.ClientID, a.ClientSecret, options...)

	case AuthTypeAuthorizationCode:
		if a.AuthURL == "" || a.ClientID == "" {
			return nil, fmt.Errorf("authorizationCode auth requires auth_url and client_id")
		}

		var options []authorizationcode.ProviderOption
		if a.InsecureSkipVerify {
			options = append(options, authorizationcode.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})))
		}

		return authorizationcode.NewDiscoveryProvider(ctx, a.AuthURL, a.ClientID, options...)

	default:
		return nil, fmt.Errorf("unsupported auth type: %q (expected static, insecureStatic, clientCredentials, or authorizationCode)", authType)
	}
}
