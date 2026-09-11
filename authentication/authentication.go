package authentication

import (
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials"
)

// Provider provides authentication credentials for connecting to a Canton participant's API endpoints.
// The Provider acts as both a raw token-source for HTTP API authentication, and a gRPC credentials provider for gRPC endpoint authentication.
//
// Implementations of this interface can implement different means of fetching and refreshing authentication tokens,
// as well as enforcing different levels of transport security. The specific implementation of the Provider
// should be chosen based on the environment being connected to (e.g. LocalNet vs. production, i.e. CI/OIDC).
type Provider interface {
	// TokenSource returns an oauth2.TokenSource that can be used to retrieve access tokens for authenticating with the participant's API endpoints.
	TokenSource() oauth2.TokenSource
	// TransportCredentials returns gRPC transport credentials to be used when connecting to the participant's RPC endpoints.
	TransportCredentials() credentials.TransportCredentials
	// PerRPCCredentials returns gRPC per-RPC credentials to be used when connecting to the participant's gRPC endpoints.
	PerRPCCredentials() credentials.PerRPCCredentials
}
