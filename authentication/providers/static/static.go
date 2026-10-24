package static

import (
	"context"
	"crypto/tls"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-canton/authentication"
)

// StaticProvider is a secure implementation of Provider that returns
// a static access token and enforces TLS transport security.
// This provider is suitable for remote environments.
type StaticProvider struct {
	AccessToken string
}

var _ authentication.Provider = StaticProvider{}

func NewStaticProvider(accessToken string) StaticProvider {
	return StaticProvider{
		AccessToken: accessToken,
	}
}

func (s StaticProvider) TokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: s.AccessToken,
	})
}

func (s StaticProvider) TransportCredentials() credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
}

func (s StaticProvider) PerRPCCredentials() credentials.PerRPCCredentials {
	return secureTokenSource{
		TokenSource: s.TokenSource(),
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
