package clientcredentials

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func newTokenServer(t *testing.T, expectedScope string, expectedAudience string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Errorf("parsing request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if values.Get("grant_type") != "client_credentials" {
			t.Errorf("expected grant_type=client_credentials, got %q", values.Get("grant_type"))
		}

		if expectedScope != "" && values.Get("scope") != expectedScope {
			t.Errorf("expected scope %q, got %q", expectedScope, values.Get("scope"))
		}

		gotAudiences := values["audience"]
		if expectedAudience == "" {
			if len(gotAudiences) != 0 {
				t.Errorf("expected no audience parameter, got %v", gotAudiences)
			}
		} else {
			if !reflect.DeepEqual(gotAudiences, []string{expectedAudience}) {
				t.Errorf("expected audience %q, got %v", expectedAudience, gotAudiences)
			}
		}

		response := tokenResponse{
			AccessToken: "test-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		payload, err := json.Marshal(response)
		if err != nil {
			t.Errorf("encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}))
}

// TestNewProvider_ValidatesInputs verifies that NewProvider returns errors when required parameters are missing.
func TestNewProvider_ValidatesInputs(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	tests := []struct {
		name         string
		tokenURL     string
		clientID     string
		clientSecret string
	}{
		{
			name:         "missing token url",
			tokenURL:     "",
			clientID:     "client-id",
			clientSecret: "client-secret",
		},
		{
			name:         "missing client id",
			tokenURL:     "https://example.test/token",
			clientID:     "",
			clientSecret: "client-secret",
		},
		{
			name:         "missing client secret",
			tokenURL:     "https://example.test/token",
			clientID:     "client-id",
			clientSecret: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewProvider(ctx, test.tokenURL, test.clientID, test.clientSecret)
			require.Error(t, err)
		})
	}
}

// TestNewProvider_UsesOptionsAndTokenSource verifies that provider options are applied and tokens can be fetched successfully.
func TestNewProvider_UsesOptionsAndTokenSource(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	server := newTokenServer(t, "scope-a scope-b", "")
	t.Cleanup(server.Close)

	customCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	provider, err := NewProvider(
		ctx,
		server.URL,
		"client-id",
		"client-secret",
		WithScopes("scope-a", "scope-b"),
		WithTransportCredentials(customCreds),
	)
	require.NoError(t, err)

	require.Same(t, customCreds, provider.TransportCredentials())

	token, err := provider.TokenSource().Token()
	require.NoError(t, err)
	require.Equal(t, "test-access-token", token.AccessToken)
}

// TestNewProvider_WithAudience verifies that an audience configured via WithAudience is sent
// to the authorization server as an "audience" parameter in the token request body.
func TestNewProvider_WithAudience(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	expectedAudience := "https://ledger.example.com"
	server := newTokenServer(t, "daml_ledger_api", expectedAudience)
	t.Cleanup(server.Close)

	provider, err := NewProvider(
		ctx,
		server.URL,
		"client-id",
		"client-secret",
		WithAudience(expectedAudience),
	)
	require.NoError(t, err)

	token, err := provider.TokenSource().Token()
	require.NoError(t, err)
	require.Equal(t, "test-access-token", token.AccessToken)
}

// TestNewProvider_WithoutAudience verifies that no "audience" parameter is sent when
// WithAudience is not used.
func TestNewProvider_WithoutAudience(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	server := newTokenServer(t, "daml_ledger_api", "")
	t.Cleanup(server.Close)

	provider, err := NewProvider(ctx, server.URL, "client-id", "client-secret")
	require.NoError(t, err)

	token, err := provider.TokenSource().Token()
	require.NoError(t, err)
	require.Equal(t, "test-access-token", token.AccessToken)
}

// TestNewDiscoveryProvider_UsesMetadataTokenEndpoint verifies that discovery metadata is fetched correctly and the token endpoint is used.
func TestNewDiscoveryProvider_UsesMetadataTokenEndpoint(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	metadataPath := "/.well-known/oauth-authorization-server"
	tokenPath := "/token"

	mux.HandleFunc(metadataPath, func(w http.ResponseWriter, r *http.Request) {
		metadata := map[string]string{
			"issuer":         server.URL,
			"token_endpoint": server.URL + tokenPath,
		}
		payload, err := json.Marshal(metadata)
		if err != nil {
			t.Errorf("encoding metadata: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})

	mux.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Errorf("parsing request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		expectedScope := "daml_ledger_api"
		if values.Get("scope") != expectedScope {
			t.Errorf("expected scope %q, got %q", expectedScope, values.Get("scope"))
		}

		response := tokenResponse{
			AccessToken: "metadata-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		payload, err := json.Marshal(response)
		if err != nil {
			t.Errorf("encoding response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})

	provider, err := NewDiscoveryProvider(ctx, server.URL, "client-id", "client-secret")
	require.NoError(t, err)

	token, err := provider.TokenSource().Token()
	require.NoError(t, err)
	require.Equal(t, "metadata-access-token", token.AccessToken)
}

// TestNewDiscoveryProvider_RequiresMetadataEndpoint verifies that discovery fails when the metadata endpoint is unavailable.
func TestNewDiscoveryProvider_RequiresMetadataEndpoint(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/.well-known/oauth-authorization-server") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	_, err := NewDiscoveryProvider(ctx, server.URL, "client-id", "client-secret")
	require.Error(t, err)
}
