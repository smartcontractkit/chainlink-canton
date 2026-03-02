# Authentication Package

OAuth 2.0 authentication providers for gRPC connections in Chainlink Canton. Supports Authorization Code Flow (
interactive user login) and Client Credentials Flow (machine-to-machine).

## Quick Start

### Client Credentials

For server-to-server authentication:

```go
import "github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"

provider, err := clientcredentials.NewDiscoveryProvider(
    ctx,
    "https://auth.example.com",
    "client-id",
    "client-secret",
)
```

### Authorization Code (Interactive Login)

For user login flows:

```go
import "github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"

provider, err := authorizationcode.NewDiscoveryProvider(
    ctx,
    "https://auth.example.com",
    "client-id",
)
```

## Authorization Code Flow

For interactive user authentication via browser login. **Requires OAuth server to support PKCE with S256 challenge
method.** Automatically handles state validation to prevent CSRF attacks. By default, automatically opens the
authorization URL in your browser.

### Metadata Discovery

Automatically discovers authorization and token endpoints via RFC 8414:

```go
import "github.com/smartcontractkit/chainlink-canton/deployment/authentication/authorizationcode"

provider, err := authorizationcode.NewDiscoveryProvider(
    ctx,
    "https://auth.example.com",
    "client-id",
    authorizationcode.WithScopes("daml_ledger_api", "offline_access"),
    authorizationcode.WithOpenBrowser(false), // Disables automatically opening the login URL in the default browser
)
```

### Direct Configuration

If metadata discovery is unavailable:

```go
provider, err := authorizationcode.NewProvider(
    ctx,
    "https://auth.example.com/oauth/authorize",
    "https://auth.example.com/oauth/token",
    "client-id",
)
```

Additional options are available via the functional options pattern (see package documentation for `WithScopes`,
`WithCallbackURL`, `WithOpenBrowser`, etc).

## Client Credentials Flow

For server-to-server authentication. Ideal for automated systems, CI/CD pipelines, and service-to-service communication.

### Metadata Discovery

Automatically discovers token endpoint via RFC 8414:

```go
import "github.com/smartcontractkit/chainlink-canton/deployment/authentication/clientcredentials"

provider, err := clientcredentials.NewDiscoveryProvider(
    ctx,
    "https://auth.example.com",
    "client-id",
    "client-secret",
    clientcredentials.WithScopes("daml_ledger_api", "admin"),
)
```

### Direct Configuration

If metadata discovery is unavailable:

```go
provider, err := clientcredentials.NewProvider(
    ctx,
    "https://auth.example.com/oauth/token",
    "client-id",
    "client-secret",
)

```

Additional options are available via the functional options pattern (see package documentation for `WithScopes`,
`WithTransportCredentials`, etc).

## Token Management

Both providers automatically handle token lifecycle management: caching, automatic refresh, and thread-safe operations.
