package static

import (
	"context"
	"crypto/tls"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
)

var _ authentication.Provider = &Provider{}

type Provider struct {
	accessToken          string
	transportCredentials credentials.TransportCredentials
}

type staticProviderConfig struct {
	transportCredentials credentials.TransportCredentials
}

func defaultStaticProviderConfig() staticProviderConfig {
	return staticProviderConfig{
		transportCredentials: credentials.NewTLS(
			&tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		),
	}
}

// ProviderOption configures the client credentials Provider using the functional options pattern.
// Options allow customization of the provider without breaking API compatibility.
type ProviderOption func(*staticProviderConfig)

// WithTransportCredentials configures the Provider to use the given transport credentials for gRPC connections.
// This allows customization of TLS settings, including certificate verification and minimum TLS version.
// The default transport credentials use TLS 1.2 or higher.
//
// Example:
//
//	WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))
func WithTransportCredentials(credentials credentials.TransportCredentials) ProviderOption {
	return func(config *staticProviderConfig) {
		config.transportCredentials = credentials
	}
}

func NewStaticProvider(accessToken string, options ...ProviderOption) *Provider {
	cfg := defaultStaticProviderConfig()
	for _, option := range options {
		option(&cfg)
	}

	return &Provider{
		accessToken:          accessToken,
		transportCredentials: cfg.transportCredentials,
	}
}

func (p Provider) TokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: p.accessToken,
	})
}

func (p Provider) TransportCredentials() credentials.TransportCredentials {
	return p.transportCredentials
}

func (p Provider) PerRPCCredentials() credentials.PerRPCCredentials {
	return secureTokenSource{
		TokenSource: p.TokenSource(),
	}
}

// secureTokenSource is a secure OAuth2 PerRPCCredentials implementation that
// requires transport security.
type secureTokenSource struct {
	oauth2.TokenSource
}

var _ credentials.PerRPCCredentials = secureTokenSource{}

func (ts secureTokenSource) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := ts.Token()
	if err != nil {
		return nil, err
	}
	if token == nil {
		//nolint:nilnil // nothing to do here, just returning no metadata and no error
		return nil, nil
	}

	return map[string]string{
		"authorization": "Bearer " + token.AccessToken,
	}, nil
}

func (ts secureTokenSource) RequireTransportSecurity() bool {
	return true
}
