package devenv

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
)

const validJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

func TestBuildParticipantAuthConfig(t *testing.T) {
	const localGRPC = "participant1.grpc-ledger-api.localhost:8080"
	const tlsGRPC = "canton-devnet.example.com:443"

	tests := []struct {
		name        string
		participant blockchain.CantonParticipantEndpoints
		env         map[string]string
		want        commonconfig.AuthConfig
		wantErr     bool
		validate    bool
	}{
		{
			name: "devenv_jwt_local_endpoint_insecureStatic",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: localGRPC,
				UserID:           "user-participant1",
				JWT:              validJWT,
			},
			want: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeInsecureStatic,
				JWT:  validJWT,
			},
			validate: true,
		},
		{
			name: "devenv_jwt_tls_endpoint_not_insecureStatic",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
				UserID:           "user-participant1",
				JWT:              validJWT,
			},
			env: map[string]string{
				envCantonAuthURL:           "https://auth.example.com/",
				envCantonOAuthClientID:     "client-id",
				envCantonOAuthClientSecret: "client-secret",
			},
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeClientCredentials,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			validate: true,
		},
		{
			name: "clientCredentials_from_env_without_jwt",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
			},
			env: map[string]string{
				envCantonAuthURL:           "https://auth.example.com/",
				envCantonOAuthClientID:     "client-id",
				envCantonOAuthClientSecret: "client-secret",
			},
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeClientCredentials,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			validate: true,
		},
		{
			name: "clientCredentials_from_ci_env_names",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
			},
			env: map[string]string{
				envCantonAuthURL: "https://auth.example.com/",
				envClientID:      "client-id",
				envClientSecret:  "client-secret",
			},
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeClientCredentials,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			validate: true,
		},
		{
			name: "clientCredentials_missing_client_secret_fails_validation",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
			},
			env: map[string]string{
				envCantonAuthURL:       "https://auth.example.com/",
				envCantonOAuthClientID: "client-id",
			},
			want: commonconfig.AuthConfig{
				Type:     commonconfig.AuthTypeClientCredentials,
				AuthURL:  "https://auth.example.com/",
				ClientID: "client-id",
			},
			validate: true,
			wantErr:  true,
		},
		{
			name: "explicit_static_with_env_jwt_override",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
				JWT:              validJWT,
			},
			env: map[string]string{
				envCantonAuthType:   commonconfig.AuthTypeStatic,
				envOnchainCantonJWT: validJWT,
			},
			want: commonconfig.AuthConfig{
				Type: commonconfig.AuthTypeStatic,
				JWT:  validJWT,
			},
			validate: true,
		},
		{
			name: "explicit_static_missing_jwt",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
			},
			env: map[string]string{
				envCantonAuthType: commonconfig.AuthTypeStatic,
			},
			wantErr: true,
		},
		{
			name: "local_without_jwt_defaults_clientCredentials",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: localGRPC,
			},
			env: map[string]string{
				envCantonAuthURL:           "https://auth.example.com/",
				envCantonOAuthClientID:     "client-id",
				envCantonOAuthClientSecret: "client-secret",
			},
			want: commonconfig.AuthConfig{
				Type:         commonconfig.AuthTypeClientCredentials,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			validate: true,
		},
		{
			name: "explicit_authorizationCode",
			participant: blockchain.CantonParticipantEndpoints{
				GRPCLedgerAPIURL: tlsGRPC,
			},
			env: map[string]string{
				envCantonAuthType:      commonconfig.AuthTypeAuthorizationCode,
				envCantonAuthURL:       "https://auth.example.com/",
				envCantonOAuthClientID: "client-id",
			},
			want: commonconfig.AuthConfig{
				Type:     commonconfig.AuthTypeAuthorizationCode,
				AuthURL:  "https://auth.example.com/",
				ClientID: "client-id",
			},
			validate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k := range tt.env {
				t.Setenv(k, tt.env[k])
			}
			for _, key := range []string{
				envCantonAuthType,
				envCantonAuthURL,
				envCantonOAuthClientID,
				envCantonOAuthClientSecret,
				envClientID,
				envClientSecret,
				envOnchainCantonJWT,
			} {
				if _, ok := tt.env[key]; !ok {
					t.Setenv(key, "")
				}
			}

			got, err := buildParticipantAuthConfig(tt.participant)
			if tt.wantErr {
				if err == nil && tt.validate {
					if validateErr := got.Validate(); validateErr == nil {
						t.Fatal("expected error, got nil")
					}

					return
				}
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}
			if err != nil {
				t.Fatalf("buildParticipantAuthConfig() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("buildParticipantAuthConfig() = %+v, want %+v", got, tt.want)
			}
			if tt.validate {
				if err := got.Validate(); err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
			}
		})
	}
}

func TestIsLocalDevenvEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		grpcURL string
		want    bool
	}{
		{grpcURL: "participant1.grpc-ledger-api.localhost:8080", want: true},
		{grpcURL: "canton-devnet.example.com:443", want: false},
		{grpcURL: "", want: false},
	}

	for _, tt := range tests {
		if got := isLocalDevenvEndpoint(tt.grpcURL); got != tt.want {
			t.Errorf("isLocalDevenvEndpoint(%q) = %v, want %v", tt.grpcURL, got, tt.want)
		}
	}
}

func TestBuildParticipantAuthConfig_no_jwt_priority_on_real_chain(t *testing.T) {
	t.Setenv(envOnchainCantonJWT, validJWT)
	t.Setenv(envCantonAuthURL, "https://auth.example.com/")
	t.Setenv(envCantonOAuthClientID, "client-id")
	t.Setenv(envCantonOAuthClientSecret, "client-secret")
	t.Setenv(envCantonAuthType, "")

	got, err := buildParticipantAuthConfig(blockchain.CantonParticipantEndpoints{
		GRPCLedgerAPIURL: "canton-devnet.example.com:443",
	})
	if err != nil {
		t.Fatalf("buildParticipantAuthConfig() error = %v", err)
	}
	if got.Type != commonconfig.AuthTypeClientCredentials {
		t.Fatalf("type = %q, want clientCredentials (env JWT must not imply static)", got.Type)
	}
	if got.JWT != "" {
		t.Fatalf("JWT should be empty for clientCredentials, got %q", got.JWT)
	}
}
