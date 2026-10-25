package authorizationcode

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// fakeTokenStore is an in-memory tokenStore for testing keyring behavior without
// touching the operating system's keyring.
type fakeTokenStore struct {
	mu      sync.Mutex
	secrets map[string]string
	getErr  error
}

func newFakeTokenStore(secrets map[string]string) *fakeTokenStore {
	if secrets == nil {
		secrets = make(map[string]string)
	}

	return &fakeTokenStore{secrets: secrets}
}

func (f *fakeTokenStore) Get(service, user string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.getErr != nil {
		return "", f.getErr
	}

	secret, ok := f.secrets[service+"/"+user]
	if !ok {
		return "", ErrTokenNotCached
	}

	return secret, nil
}

func (f *fakeTokenStore) Set(service, user, secret string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.secrets[service+"/"+user] = secret

	return nil
}

// newOIDCTokenServer returns a token endpoint that grants tokens for both the
// authorization_code and the refresh_token grant types, simulating refresh token rotation.
func newOIDCTokenServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parsing form: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		accessToken, refreshToken := "auth-code-token", "refresh-token-1"
		switch r.Form.Get("grant_type") {
		case "authorization_code":
			if r.Form.Get("code") == "" || r.Form.Get("code_verifier") == "" {
				t.Errorf("expected code and code_verifier to be set")
				w.WriteHeader(http.StatusBadRequest)

				return
			}
		case "refresh_token":
			if r.Form.Get("refresh_token") == "" {
				t.Errorf("expected refresh_token to be set")
				w.WriteHeader(http.StatusBadRequest)

				return
			}
			accessToken, refreshToken = "renewed-token", "refresh-token-2"
		default:
			t.Errorf("unexpected grant_type %q", r.Form.Get("grant_type"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		payload, err := json.Marshal(map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
		if err != nil {
			t.Errorf("encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
}

// completeLoginFlow runs NewProvider in the background and simulates the browser callback,
// returning the resulting provider.
func completeLoginFlow(t *testing.T, authURL, tokenURL string, options ...ProviderOption) *Provider {
	t.Helper()

	callbackHost := freePort(t)
	allOptions := append([]ProviderOption{
		WithCallbackURL("http://" + callbackHost + "/callback"),
		WithOpenBrowser(false),
		WithTimeout(5 * time.Second),
	}, options...)

	output, restore := captureStdout(t)
	defer restore()

	resultCh := make(chan struct {
		provider *Provider
		err      error
	}, 1)

	go func() {
		provider, err := NewProvider(t.Context(), authURL, tokenURL, "client-id", allOptions...)
		resultCh <- struct {
			provider *Provider
			err      error
		}{provider: provider, err: err}
	}()

	var authCodeURL string
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		authCodeURL = extractFirstURL(output.String())
		if authCodeURL != "" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.NotEmpty(t, authCodeURL, "auth code URL not found in output")

	parsed, err := url.Parse(authCodeURL)
	require.NoError(t, err, "parsing auth URL")
	state := parsed.Query().Get("state")
	require.NotEmpty(t, state, "state not found in auth URL")

	callbackURL := "http://" + callbackHost + "/callback?code=code123&state=" + url.QueryEscape(state)
	response, err := http.Get(callbackURL) //nolint:noctx
	require.NoError(t, err, "requesting callback")
	require.NoError(t, response.Body.Close())

	result := <-resultCh
	require.NoError(t, result.err)
	require.NotNil(t, result.provider)

	return result.provider
}

func storeKey(tokenURL, clientID string) string {
	key := keyringKeyFor(tokenURL, clientID)

	return key.service + "/" + key.user
}

// TestNewProvider_KeyringValidToken_SkipsLogin verifies that a still-valid token from the
// keyring is used directly, without starting the interactive login flow (no token server
// is reachable here, so any login attempt would fail).
func TestNewProvider_KeyringValidToken_SkipsLogin(t *testing.T) {
	t.Parallel()

	store := newFakeTokenStore(nil)
	secret, err := marshalToken(&oauth2.Token{
		AccessToken: "cached-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	})
	require.NoError(t, err, "marshaling token")
	store.secrets[storeKey("https://example.test/token", "client-id")] = secret

	provider, err := NewProvider(
		t.Context(),
		"https://example.test/auth",
		"https://example.test/token",
		"client-id",
		withTokenStore(store),
	)
	require.NoError(t, err)

	token, err := provider.TokenSource().Token()
	require.NoError(t, err, "requesting token")
	require.Equal(t, "cached-token", token.AccessToken)

	// The unchanged cached token must not be rewritten to the store.
	require.Equal(t, secret, store.secrets[storeKey("https://example.test/token", "client-id")])
}

// TestNewProvider_KeyringExpiredToken_RenewsAndPersists verifies that an expired cached token
// with a refresh token is renewed via the token endpoint and the renewed token (including the
// rotated refresh token) is written back to the keyring.
func TestNewProvider_KeyringExpiredToken_RenewsAndPersists(t *testing.T) {
	t.Parallel()

	tokenServer := newOIDCTokenServer(t)
	t.Cleanup(tokenServer.Close)

	tokenURL := tokenServer.URL + "/token"

	store := newFakeTokenStore(nil)
	secret, err := marshalToken(&oauth2.Token{
		AccessToken:  "expired-token",
		RefreshToken: "refresh-token-1",
		Expiry:       time.Now().Add(-time.Hour),
	})
	require.NoError(t, err, "marshaling token")
	store.secrets[storeKey(tokenURL, "client-id-123")] = secret

	provider, err := NewProvider(
		t.Context(),
		tokenServer.URL+"/auth",
		tokenURL,
		"client-id-123",
		withTokenStore(store),
	)
	require.NoError(t, err, "cached refresh token should skip the login flow")

	token, err := provider.TokenSource().Token()
	require.NoError(t, err, "requesting token")
	require.Equal(t, "renewed-token", token.AccessToken)

	stored, err := unmarshalToken(store.secrets[storeKey(tokenURL, "client-id-123")])
	require.NoError(t, err, "unmarshaling persisted token")
	require.Equal(t, "renewed-token", stored.AccessToken)
	require.Equal(t, "refresh-token-2", stored.RefreshToken)
}

// TestNewProvider_KeyringFallbackToLogin verifies that when no usable token is cached
// (nothing cached, expired without refresh token, or keyring read failure), the provider
// falls back to the interactive login flow and persists the fetched token.
func TestNewProvider_KeyringFallbackToLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T) *fakeTokenStore
	}{
		{
			name: "no token cached",
			setup: func(t *testing.T) *fakeTokenStore {
				return newFakeTokenStore(nil)
			},
		},
		{
			name: "expired token without refresh token",
			setup: func(t *testing.T) *fakeTokenStore {
				secret, err := marshalToken(&oauth2.Token{
					AccessToken: "expired-token",
					Expiry:      time.Now().Add(-time.Hour),
				})
				require.NoError(t, err, "marshaling token")

				return newFakeTokenStore(map[string]string{"placeholder": secret})
			},
		},
		{
			name: "corrupted token",
			setup: func(t *testing.T) *fakeTokenStore {
				return newFakeTokenStore(map[string]string{"placeholder": "not-json"})
			},
		},
		{
			name: "keyring read error",
			setup: func(t *testing.T) *fakeTokenStore {
				store := newFakeTokenStore(nil)
				store.getErr = errors.New("keyring unavailable")

				return store
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tokenServer := newOIDCTokenServer(t)
			t.Cleanup(tokenServer.Close)

			tokenURL := tokenServer.URL + "/token"

			store := test.setup(t)
			// Place the token under the key the provider will look up.
			store.secrets[storeKey(tokenURL, "client-id")] = store.secrets["placeholder"]

			provider := completeLoginFlow(t, tokenServer.URL+"/auth", tokenURL, withTokenStore(store))

			token, err := provider.TokenSource().Token()
			require.NoError(t, err, "requesting token")
			require.Equal(t, "auth-code-token", token.AccessToken)

			// The fetched token must be persisted, including its refresh token.
			stored, err := unmarshalToken(store.secrets[storeKey(tokenURL, "client-id")])
			require.NoError(t, err, "unmarshaling persisted token")
			require.Equal(t, "auth-code-token", stored.AccessToken)
			require.Equal(t, "refresh-token-1", stored.RefreshToken)
		})
	}
}

// TestKeyringKeyFor_DistinctAndStable verifies that keyring keys differ per token endpoint
// and client ID, and are stable for repeated calls.
func TestKeyringKeyFor_DistinctAndStable(t *testing.T) {
	t.Parallel()

	a := keyringKeyFor("https://a.test/token", "client-id")
	b := keyringKeyFor("https://b.test/token", "client-id")
	c := keyringKeyFor("https://a.test/token", "other-client")

	require.Equal(t, keyringService, a.service)
	require.NotEqual(t, a.user, b.user)
	require.NotEqual(t, a.user, c.user)
	require.Equal(t, a, keyringKeyFor("https://a.test/token", "client-id"))
}

// TestTokenMarshalRoundTrip verifies that tokens survive serialization, including expiry
// and refresh token.
func TestTokenMarshalRoundTrip(t *testing.T) {
	t.Parallel()

	token := &oauth2.Token{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour).Round(time.Second),
	}

	secret, err := marshalToken(token)
	require.NoError(t, err, "marshaling token")

	got, err := unmarshalToken(secret)
	require.NoError(t, err, "unmarshaling token")
	require.Equal(t, token, got)
}
