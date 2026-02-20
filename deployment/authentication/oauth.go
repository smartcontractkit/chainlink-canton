package authentication

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/icza/gox/osx"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton/provider/authentication"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var _ authentication.Provider = OIDCProvider{}

type OIDCProvider struct {
	s oauth.TokenSource
}

func NewClientCredentialsProvider(ctx context.Context, authURL, clientID, clientSecret string) (OIDCProvider, error) {
	tokenURL := fmt.Sprintf("%s/v1/token", authURL)

	oauthCfg := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{"daml_ledger_api"},
	}

	tokenSource := oauthCfg.TokenSource(ctx)

	return OIDCProvider{
		s: oauth.TokenSource{TokenSource: tokenSource},
	}, nil
}

func NewAuthorizationCodeProvider(ctx context.Context, authURL, clientID string) (OIDCProvider, error) {
	// Generate PKCE verifier using built-in oauth2 support
	verifier := oauth2.GenerateVerifier()

	port := 8400
	authEndpoint := fmt.Sprintf("%s/v1/authorize", authURL)
	tokenEndpoint := fmt.Sprintf("%s/v1/token", authURL)
	redirectURL := fmt.Sprintf("http://localhost:%d", port)

	oauthCfg := &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: redirectURL + "/callback",
		Scopes:      []string{"openid", "daml_ledger_api"},
		Endpoint:    oauth2.Endpoint{AuthURL: authEndpoint, TokenURL: tokenEndpoint},
	}

	// Generate cryptographically secure random state
	state := generateState()

	// Use built-in S256ChallengeOption for PKCE
	authCodeURL := oauthCfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))

	callbackChan := make(chan *oauth2.Token)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := q.Get("code")
		receivedState := q.Get("state")

		if receivedState != state {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		// Use built-in VerifierOption for PKCE
		token, err := oauthCfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
		if err != nil {
			http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		callbackChan <- token

		// HTML response for the browser
		html := `
			<!DOCTYPE html>
			<html>
			<head><title>Authentication Complete</title></head>
			<body style="font-family: sans-serif; text-align: center; padding: 40px;">
				<h1>Authentication complete!</h1>
				<p>You can safely close this window.</p>
			</body>
			</html>
		`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	})

	// Start callback server
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           serveMux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Println("Waiting for authentication...")
	go func() {
		_ = server.ListenAndServe()
	}()
	fmt.Println("Attempting to open your default browser.\nIf the browser does not open, open the following URL:")
	fmt.Println(authCodeURL)
	_ = osx.OpenDefault(authCodeURL)
	select {
	case token := <-callbackChan:
		fmt.Println("Authentication complete")
		tokenSource := oauthCfg.TokenSource(ctx, token)

		return OIDCProvider{
			s: oauth.TokenSource{TokenSource: tokenSource},
		}, nil
	case <-ctx.Done():
		return OIDCProvider{}, ctx.Err()
	}
}

func generateState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

func (p OIDCProvider) TokenSource() oauth2.TokenSource {
	return p.s.TokenSource
}

func (p OIDCProvider) TransportCredentials() credentials.TransportCredentials {
	return credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
	})
}

func (p OIDCProvider) PerRPCCredentials() credentials.PerRPCCredentials {
	return p.s
}
