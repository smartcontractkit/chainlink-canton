package authorizationcode

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// safeBuffer is a thread-safe wrapper around bytes.Buffer that protects concurrent writes.
// It is used to safely capture stdout output from goroutines that may run concurrently with
// the test. The mutex ensures that Write and String operations do not race when tests are
// run in parallel via t.Parallel().
type safeBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

var stdoutMu sync.Mutex

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.String()
}

// captureStdout temporarily redirects os.Stdout to capture output from the provider initialization.
// This is necessary because NewProvider prints the authorization URL and status messages to stdout,
// which the test needs to extract to complete the OAuth2 flow simulation.
//
// The function uses a mutex-protected section to ensure only one test captures stdout at a time,
// preventing race conditions when tests run in parallel. It returns a safeBuffer to collect the
// output and a restore function that closes the pipe and restores the original stdout.
//
// Usage:
//
//	output, restore := captureStdout(t)
//	defer restore()
//	// code that writes to stdout
//	url := extractFirstURL(output.String())
func captureStdout(t *testing.T) (*safeBuffer, func()) {
	t.Helper()

	// Mutex to only capture stdout in one test at a time
	stdoutMu.Lock()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err, "creating stdout pipe")

	buffer := &safeBuffer{}
	done := make(chan struct{})

	os.Stdout = writer
	go func() {
		_, _ = io.Copy(buffer, reader)
		close(done)
	}()

	// Restore os.StdOut and close the writer when done
	return buffer, func() {
		_ = writer.Close()
		os.Stdout = original
		<-done
		stdoutMu.Unlock()
	}
}

func freePort(t *testing.T) string {
	t.Helper()

	//nolint:noctx
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "allocating port")
	addr := listener.Addr().String()
	_ = listener.Close()

	return addr
}

func extractFirstURL(output string) string {
	for line := range strings.SplitSeq(output, "\n") {
		candidate := strings.TrimSpace(line)
		if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") {
			if parsed, err := url.Parse(candidate); err == nil && parsed.Scheme != "" && parsed.Host != "" {
				return candidate
			}
		}
	}

	return ""
}

func newTokenServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parsing form: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		switch grantType := r.Form.Get("grant_type"); grantType {
		case "authorization_code":
			if r.Form.Get("code") == "" {
				t.Errorf("expected code to be set")
			}
			if r.Form.Get("code_verifier") == "" {
				t.Errorf("expected code_verifier to be set")
			}

			payload, err := json.Marshal(map[string]any{
				"access_token": "auth-code-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
			if err != nil {
				t.Errorf("encoding response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(payload)
		default:
			t.Errorf("unexpected grant_type %q", grantType)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

func newRefreshingTokenServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parsing form: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		switch grantType := r.Form.Get("grant_type"); grantType {
		case "authorization_code":
			if r.Form.Get("code") == "" {
				t.Errorf("expected code to be set")
			}
			if r.Form.Get("code_verifier") == "" {
				t.Errorf("expected code_verifier to be set")
			}

			payload, err := json.Marshal(map[string]any{
				"access_token":  "auth-code-token",
				"refresh_token": "refresh-token",
				"token_type":    "Bearer",
				"expires_in":    -1,
			})
			if err != nil {
				t.Errorf("encoding response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(payload)
		case "refresh_token":
			if r.Form.Get("refresh_token") != "refresh-token" {
				t.Errorf("expected refresh_token to be set")
			}

			payload, err := json.Marshal(map[string]any{
				"access_token": "refreshed-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
			if err != nil {
				t.Errorf("encoding response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(payload)
		default:
			t.Errorf("unexpected grant_type %q", grantType)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}

// TestNewProvider_ValidatesInputs verifies that NewProvider returns errors when required parameters are missing.
func TestNewProvider_ValidatesInputs(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tests := []struct {
		name     string
		authURL  string
		tokenURL string
		clientID string
	}{
		{
			name:     "missing auth url",
			authURL:  "",
			tokenURL: "https://example.test/token",
			clientID: "client-id",
		},
		{
			name:     "missing token url",
			authURL:  "https://example.test/auth",
			tokenURL: "",
			clientID: "client-id",
		},
		{
			name:     "missing client id",
			authURL:  "https://example.test/auth",
			tokenURL: "https://example.test/token",
			clientID: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewProvider(ctx, test.authURL, test.tokenURL, test.clientID)
			require.Error(t, err)
		})
	}
}

// TestNewDiscoveryProvider_RequiresS256 verifies that the provider rejects authorization servers that don't support S256 PKCE.
func TestNewDiscoveryProvider_RequiresS256(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuer := "http://" + r.Host
		payload, err := json.Marshal(map[string]any{
			"issuer":                 issuer,
			"authorization_endpoint": "https://example.test/auth",
			"token_endpoint":         "https://example.test/token",
			"code_challenge_methods": []string{"plain"},
		})
		if err != nil {
			t.Errorf("encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
	t.Cleanup(server.Close)

	_, err := NewDiscoveryProvider(ctx, server.URL, "client-id")
	require.Error(t, err)
}

// TestNewProvider_FlowCompletes simulates a complete authorization code flow with PKCE, from printing the auth URL to exchanging the code for a token.
func TestNewProvider_FlowCompletes(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	callbackHost := freePort(t)

	tokenServer := newTokenServer(t)
	t.Cleanup(tokenServer.Close)

	output, restore := captureStdout(t)
	defer restore()

	resultCh := make(chan struct {
		provider *Provider
		err      error
	}, 1)

	// Start Provider in the background to capture callback URL
	go func() {
		provider, err := NewProvider(
			ctx,
			tokenServer.URL+"/auth",
			tokenServer.URL+"/token",
			"client-id",
			WithCallbackURL("http://"+callbackHost+"/callback"),
			WithOpenBrowser(false),
			WithTimeout(5*time.Second),
		)
		resultCh <- struct {
			provider *Provider
			err      error
		}{provider: provider, err: err}
	}()

	// Capture callback URL from stdout
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
	response, err := http.Get(callbackURL) //nolint:gosec,noctx
	require.NoError(t, err, "requesting callback")
	require.NoError(t, response.Body.Close())

	result := <-resultCh
	require.NoError(t, result.err)
	require.NotNil(t, result.provider)

	token, err := result.provider.TokenSource().Token()
	require.NoError(t, err, "requesting token")
	require.Equal(t, "auth-code-token", token.AccessToken)
}

func TestNewProvider_RefreshUsesContextWithoutCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	callbackHost := freePort(t)

	tokenServer := newRefreshingTokenServer(t)
	t.Cleanup(tokenServer.Close)

	output, restore := captureStdout(t)
	defer restore()

	resultCh := make(chan struct {
		provider *Provider
		err      error
	}, 1)

	go func() {
		provider, err := NewProvider(
			ctx,
			tokenServer.URL+"/auth",
			tokenServer.URL+"/token",
			"client-id",
			WithCallbackURL("http://"+callbackHost+"/callback"),
			WithOpenBrowser(false),
			WithTimeout(5*time.Second),
		)
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
	response, err := http.Get(callbackURL) //nolint:gosec,noctx
	require.NoError(t, err, "requesting callback")
	require.NoError(t, response.Body.Close())

	result := <-resultCh
	require.NoError(t, result.err)
	require.NotNil(t, result.provider)

	cancel()

	refreshedToken, err := result.provider.TokenSource().Token()
	require.NoError(t, err, "refreshing token after parent context cancellation")
	require.Equal(t, "refreshed-token", refreshedToken.AccessToken)
}
