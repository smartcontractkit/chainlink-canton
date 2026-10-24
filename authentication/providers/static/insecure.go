package static

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-canton/authentication"
)

// InsecureStaticProvider is an insecure implementation of Provider that always
// returns the same static access token and does not provide/enforce transport security.
// This provider is only suitable for testing against LocalNet or other non-production environments.
type InsecureStaticProvider struct {
	AccessToken string
}

var _ authentication.Provider = InsecureStaticProvider{}

func NewInsecureStaticProvider(accessToken string) InsecureStaticProvider {
	return InsecureStaticProvider{
		AccessToken: accessToken,
	}
}

func (i InsecureStaticProvider) TokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: i.AccessToken,
	})
}

func (i InsecureStaticProvider) TransportCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}

func (i InsecureStaticProvider) PerRPCCredentials() credentials.PerRPCCredentials {
	return insecureTokenSource{
		TokenSource: i.TokenSource(),
	}
}

// insecureTokenSource is an insecure OAuth2 PerRPCCredentials implementation that retrieves tokens from an underlying oauth2.TokenSource.
// It does not enforce transport security, making it only suitable for testing against LocalNet.
type insecureTokenSource struct {
	oauth2.TokenSource
}

var _ credentials.PerRPCCredentials = insecureTokenSource{}

func (ts insecureTokenSource) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
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

func (ts insecureTokenSource) RequireTransportSecurity() bool {
	return false
}
