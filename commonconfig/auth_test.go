package commonconfig

import (
	"errors"
	"testing"
)

// validJWT is a well-formed JWT (header.payload.signature) accepted by the validator's jwt tag.
const validJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

func TestAuthConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  AuthConfig
		wantErr bool
	}{
		// --- static auth: only Type and JWT required ---
		{
			name: "static_valid",
			config: AuthConfig{
				Type: AuthTypeStatic,
				JWT:  validJWT,
			},
			wantErr: false,
		},
		{
			name: "static_missing_jwt",
			config: AuthConfig{
				Type: AuthTypeStatic,
			},
			wantErr: true,
		},
		{
			name: "static_invalid_jwt_format",
			config: AuthConfig{
				Type: AuthTypeStatic,
				JWT:  "not-a-valid-jwt",
			},
			wantErr: true,
		},
		{
			name: "static_jwt_only_two_parts",
			config: AuthConfig{
				Type: AuthTypeStatic,
				JWT:  "part1.part2",
			},
			wantErr: true,
		},

		// --- insecureStatic auth: same as static (Type + JWT required), but uses insecure transport ---
		{
			name: "insecureStatic_valid",
			config: AuthConfig{
				Type: AuthTypeInsecureStatic,
				JWT:  validJWT,
			},
			wantErr: false,
		},
		{
			name: "insecureStatic_missing_jwt",
			config: AuthConfig{
				Type: AuthTypeInsecureStatic,
			},
			wantErr: true,
		},
		{
			name: "insecureStatic_invalid_jwt_format",
			config: AuthConfig{
				Type: AuthTypeInsecureStatic,
				JWT:  "not-a-valid-jwt",
			},
			wantErr: true,
		},
		{
			name: "insecureStatic_jwt_only_two_parts",
			config: AuthConfig{
				Type: AuthTypeInsecureStatic,
				JWT:  "part1.part2",
			},
			wantErr: true,
		},

		// --- clientCredentials: Type, UserID, AuthURL, ClientID, ClientSecret required; AuthURL must be valid URL ---
		{
			name: "clientCredentials_valid",
			config: AuthConfig{
				Type:         AuthTypeClientCredentials,
				UserID:       "user-1",
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			wantErr: false,
		},
		{
			name: "clientCredentials_missing_user_id",
			config: AuthConfig{
				Type:         AuthTypeClientCredentials,
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			wantErr: true,
		},
		{
			name: "clientCredentials_missing_auth_url",
			config: AuthConfig{
				Type:         AuthTypeClientCredentials,
				UserID:       "user-1",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			wantErr: true,
		},
		{
			name: "clientCredentials_invalid_auth_url",
			config: AuthConfig{
				Type:         AuthTypeClientCredentials,
				UserID:       "user-1",
				AuthURL:      "not-a-url",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			wantErr: true,
		},
		{
			name: "clientCredentials_missing_client_id",
			config: AuthConfig{
				Type:         AuthTypeClientCredentials,
				UserID:       "user-1",
				AuthURL:      "https://auth.example.com/",
				ClientSecret: "client-secret",
			},
			wantErr: true,
		},
		{
			name: "clientCredentials_missing_client_secret",
			config: AuthConfig{
				Type:     AuthTypeClientCredentials,
				UserID:   "user-1",
				AuthURL:  "https://auth.example.com/",
				ClientID: "client-id",
			},
			wantErr: true,
		},

		// --- authorizationCode: Type, UserID, AuthURL, ClientID required; ClientSecret must be unset (excluded_unless clientCredentials) ---
		{
			name: "authorizationCode_valid",
			config: AuthConfig{
				Type:     AuthTypeAuthorizationCode,
				UserID:   "user-1",
				AuthURL:  "https://auth.example.com/",
				ClientID: "client-id",
			},
			wantErr: false,
		},
		{
			name: "authorizationCode_missing_user_id",
			config: AuthConfig{
				Type:     AuthTypeAuthorizationCode,
				AuthURL:  "https://auth.example.com/",
				ClientID: "client-id",
			},
			wantErr: true,
		},
		{
			name: "authorizationCode_missing_auth_url",
			config: AuthConfig{
				Type:     AuthTypeAuthorizationCode,
				UserID:   "user-1",
				ClientID: "client-id",
			},
			wantErr: true,
		},
		{
			name: "authorizationCode_missing_client_id",
			config: AuthConfig{
				Type:    AuthTypeAuthorizationCode,
				UserID:  "user-1",
				AuthURL: "https://auth.example.com/",
			},
			wantErr: true,
		},
		{
			name: "authorizationCode_client_secret_must_be_unset",
			config: AuthConfig{
				Type:         AuthTypeAuthorizationCode,
				UserID:       "user-1",
				AuthURL:      "https://auth.example.com/",
				ClientID:     "client-id",
				ClientSecret: "must-not-set",
			},
			wantErr: true,
		},

		// --- Type validation ---
		{
			name: "missing_type",
			config: AuthConfig{
				JWT: validJWT,
			},
			wantErr: true,
		},
		{
			name: "invalid_type",
			config: AuthConfig{
				Type: "invalidAuthType",
				JWT:  validJWT,
			},
			wantErr: true,
		},
		{
			name: "empty_type",
			config: AuthConfig{
				Type: "",
				JWT:  validJWT,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthConfig_Validate_error_wrapping(t *testing.T) {
	t.Parallel()

	cfg := AuthConfig{Type: AuthTypeStatic} // missing JWT
	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() expected error, got nil")
	}
	// Wrapped error should be unwrappable and message prefixed
	if len(err.Error()) < 13 || err.Error()[:13] != "auth config: " {
		t.Errorf("Validate() error should be wrapped with 'auth config:' prefix, got: %q", err.Error())
	}
	_ = errors.Unwrap(err) // should not panic
}
