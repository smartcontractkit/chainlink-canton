package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	// Save original env vars
	origJWTSecret := os.Getenv("JWT_SECRET")
	origPort := os.Getenv("PORT")
	origHost := os.Getenv("HOST")
	defer func() {
		os.Setenv("JWT_SECRET", origJWTSecret)
		os.Setenv("PORT", origPort)
		os.Setenv("HOST", origHost)
	}()

	t.Run("returns error when JWT_SECRET is missing", func(t *testing.T) {
		t.Parallel()

		os.Setenv("JWT_SECRET", "")
		_, err := Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "JWT_SECRET")
	})

	t.Run("loads config with defaults", func(t *testing.T) {
		t.Parallel()

		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("PORT", "")
		os.Setenv("HOST", "")

		cfg, err := Load()
		require.NoError(t, err)

		assert.Equal(t, 8090, cfg.Port)
		assert.Equal(t, "0.0.0.0", cfg.Host)
		assert.Equal(t, "localhost", cfg.LedgerAPIHost)
		assert.Equal(t, 10001, cfg.LedgerAPIPort)
		assert.Equal(t, "test-secret", cfg.JWTSecret)
		assert.Equal(t, "https://canton.network.global", cfg.JWTAudience)
		assert.Equal(t, "./environments.yaml", cfg.EnvironmentsConfigPath)
	})

	t.Run("loads config from env vars", func(t *testing.T) {
		t.Parallel()

		os.Setenv("JWT_SECRET", "my-secret")
		os.Setenv("PORT", "9090")
		os.Setenv("HOST", "127.0.0.1")
		os.Setenv("LEDGER_API_HOST", "canton-node")
		os.Setenv("LEDGER_API_PORT", "5001")

		cfg, err := Load()
		require.NoError(t, err)

		assert.Equal(t, 9090, cfg.Port)
		assert.Equal(t, "127.0.0.1", cfg.Host)
		assert.Equal(t, "canton-node", cfg.LedgerAPIHost)
		assert.Equal(t, 5001, cfg.LedgerAPIPort)
		assert.Equal(t, "my-secret", cfg.JWTSecret)

		// Cleanup
		os.Unsetenv("LEDGER_API_HOST")
		os.Unsetenv("LEDGER_API_PORT")
	})

	t.Run("handles invalid port gracefully", func(t *testing.T) {
		t.Parallel()

		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("PORT", "invalid")

		cfg, err := Load()
		require.NoError(t, err)
		assert.Equal(t, 8090, cfg.Port) // Falls back to default
	})
}

func TestGetEnv(t *testing.T) {
	t.Parallel()

	t.Run("returns env value when set", func(t *testing.T) {
		t.Parallel()

		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default")
		assert.Equal(t, "test-value", result)
	})

	t.Run("returns default when env not set", func(t *testing.T) {
		t.Parallel()

		os.Unsetenv("TEST_VAR_UNSET")

		result := getEnv("TEST_VAR_UNSET", "default-value")
		assert.Equal(t, "default-value", result)
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Parallel()

	t.Run("returns env value as int when valid", func(t *testing.T) {
		t.Parallel()

		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		result := getEnvInt("TEST_INT", 0)
		assert.Equal(t, 42, result)
	})

	t.Run("returns default when env not set", func(t *testing.T) {
		t.Parallel()

		os.Unsetenv("TEST_INT_UNSET")

		result := getEnvInt("TEST_INT_UNSET", 100)
		assert.Equal(t, 100, result)
	})

	t.Run("returns default when env value is invalid", func(t *testing.T) {
		t.Parallel()

		os.Setenv("TEST_INT_INVALID", "not-a-number")
		defer os.Unsetenv("TEST_INT_INVALID")

		result := getEnvInt("TEST_INT_INVALID", 50)
		assert.Equal(t, 50, result)
	})
}
