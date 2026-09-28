package authorizationcode

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials/oauth"
)

// keyringService is the service name under which tokens are stored in the OS keyring.
const keyringService = "chainlink-canton"

// ErrTokenNotCached is returned by tokenStore.Get when no token is stored for the given key.
var ErrTokenNotCached = errors.New("token not found in keyring")

// tokenStore persists OAuth2 tokens as opaque secrets. It abstracts OS keyring access
// so implementations can be swapped (e.g. in tests).
type tokenStore interface {
	// Get returns the stored secret for the service/user pair, or ErrTokenNotCached if none exists.
	Get(service, user string) (string, error)
	// Set stores (or overwrites) the secret for the service/user pair.
	Set(service, user, secret string) error
}

// keyringStore is a tokenStore backed by the operating system's native keyring
// (e.g. macOS Keychain, Windows Credential Manager, Secret Service on Linux).
type keyringStore struct{}

func (keyringStore) Get(service, user string) (string, error) {
	secret, err := keyring.Get(service, user)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrTokenNotCached
	}

	return secret, err
}

func (keyringStore) Set(service, user, secret string) error {
	return keyring.Set(service, user, secret)
}

// tokenKey identifies a token within a tokenStore.
type tokenKey struct {
	service string
	user    string
}

// keyringKeyFor derives the keyring key for a token issued by the given token endpoint
// to the given client. The user part stays readable (the client ID) while a short digest
// of the token endpoint disambiguates entries from different authorization servers.
func keyringKeyFor(tokenURL, clientID string) tokenKey {
	digest := sha256.Sum256([]byte(tokenURL + "\x00" + clientID))

	return tokenKey{
		service: keyringService,
		user:    fmt.Sprintf("%s@%s", clientID, hex.EncodeToString(digest[:8])),
	}
}

// marshalToken serializes a token (including its refresh token and expiry) into a storable secret.
func marshalToken(token *oauth2.Token) (string, error) {
	data, err := json.Marshal(token) //nolint:gosec // only used to persist token to keyring
	if err != nil {
		return "", fmt.Errorf("marshaling token: %w", err)
	}

	return string(data), nil
}

// unmarshalToken deserializes a token previously stored with marshalToken.
func unmarshalToken(secret string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(secret), &token); err != nil {
		return nil, fmt.Errorf("unmarshaling token: %w", err)
	}

	return &token, nil
}

// persistToken stores the token in the store, printing a warning (but not failing) if it cannot be saved:
// the user is already authenticated at this point, so a broken keyring should not abort the flow.
func persistToken(store tokenStore, key tokenKey, token *oauth2.Token) {
	secret, err := marshalToken(token)
	if err != nil {
		fmt.Printf("WARNING: could not serialize token for keyring: %v\n", err)

		return
	}

	if err := store.Set(key.service, key.user, secret); err != nil {
		fmt.Printf("WARNING: could not store token in keyring: %v\n", err)
	}
}

// keyringTokenSource wraps an oauth2.TokenSource and persists tokens to the store whenever they
// change. This keeps the keyring in sync when the authorization server rotates refresh tokens
// (e.g. Okta rotates by default), which would otherwise invalidate the cached token.
type keyringTokenSource struct {
	inner oauth2.TokenSource
	store tokenStore
	key   tokenKey

	mu   sync.Mutex
	last *oauth2.Token
}

func newKeyringTokenSource(inner oauth2.TokenSource, store tokenStore, key tokenKey, last *oauth2.Token) *keyringTokenSource {
	return &keyringTokenSource{inner: inner, store: store, key: key, last: last}
}

func (s *keyringTokenSource) Token() (*oauth2.Token, error) {
	token, err := s.inner.Token()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.last == nil ||
		token.AccessToken != s.last.AccessToken ||
		token.RefreshToken != s.last.RefreshToken ||
		!token.Expiry.Equal(s.last.Expiry) {
		persistToken(s.store, s.key, token)
		s.last = token
	}

	return token, nil
}

// providerFromKeyring attempts to build a Provider from a token cached in the configured
// token store, skipping the interactive login flow. It returns ok=false when no usable token
// is cached, so the caller can fall back to the login flow.
func providerFromKeyring(ctx context.Context, cfg *authorizationCodeProviderConfig, oauthCfg *oauth2.Config) (*Provider, bool) {
	key := keyringKeyFor(oauthCfg.Endpoint.TokenURL, oauthCfg.ClientID)

	secret, err := cfg.store.Get(key.service, key.user)
	switch {
	case errors.Is(err, ErrTokenNotCached):
		return nil, false
	case err != nil:
		fmt.Printf("WARNING: could not read token from keyring: %v\n", err)

		return nil, false
	}

	token, err := unmarshalToken(secret)
	if err != nil {
		fmt.Println("WARNING: ignoring corrupted token found in keyring")

		return nil, false
	}

	// A cached token is usable if it is still valid, or if it can be renewed
	// transparently via its refresh token.
	if token.AccessToken == "" || (!token.Valid() && token.RefreshToken == "") {
		return nil, false
	}

	// Validate the cached token eagerly: a valid access token is returned from the token
	// source's cache, an expired one is renewed via its refresh token. Validation is
	// bounded by the flow timeout. If the token cannot be used (e.g. the refresh token
	// was revoked), fall back to the interactive login instead of failing later on the
	// first RPC.
	validateCtx := ctx
	if cfg.timeout > 0 {
		var cancel context.CancelFunc
		validateCtx, cancel = context.WithTimeout(ctx, cfg.timeout)
		defer cancel()
	}

	validated, err := oauthCfg.TokenSource(validateCtx, token).Token()
	if err != nil {
		fmt.Printf("WARNING: token from keyring is no longer usable (%v); logging in again\n", err)

		return nil, false
	}

	if validated.AccessToken != token.AccessToken {
		// The token was renewed; keep the keyring in sync (refresh tokens may be
		// rotated on use).
		persistToken(cfg.store, key, validated)
	}

	tokenSource := newKeyringTokenSource(oauthCfg.TokenSource(ctx, validated), cfg.store, key, validated)

	fmt.Println("Using cached authentication token from keyring")

	return &Provider{
		tokenSource:          oauth.TokenSource{TokenSource: tokenSource},
		transportCredentials: cfg.transportCredentials,
	}, true
}
