package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smartcontractkit/chainlink-canton/commonconfig"
	"github.com/stretchr/testify/require"
)

func TestLoadAndStateRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "registry-kit.toml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(`
network = "devnet-cv1"

[ledger]
json_api_url = "https://example.test/api/json"
grpc_ledger_api_url = "example.test:443"
admin_api_url = "example.test:443"
user_id = "user-1"
synchronizer_id = "global-domain::abc"

[ledger.auth]
type = "insecureStatic"
jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

[parties]
operator = "op::1"
provider = "prov::1"
registrar = "reg::1"

[ccip]
token_admin_registry_cid = "tar-cid"
ccip_party = "ccip::1"
burn_mint_pool_instance_id = "pool-1"
`), 0o644))

	cfg, err := Load(cfgPath)
	require.NoError(t, err)
	require.Equal(t, "devnet-cv1", cfg.Network)
	require.Equal(t, commonconfig.AuthTypeInsecureStatic, cfg.Ledger.Auth.Type)

	statePath := StatePathNextTo(cfgPath)
	st := State{InstrumentID: "TEST-USD", RegistrarServiceCID: "reg-svc"}
	require.NoError(t, st.Save(statePath))

	loaded, err := LoadState(statePath)
	require.NoError(t, err)
	require.Equal(t, "TEST-USD", loaded.InstrumentID)
}

func TestLoadAppliesOperatorBackendDefault(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "registry-kit.toml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(`
network = "devnet-cv1"

[ledger]
json_api_url = "https://example.test/api/json"
grpc_ledger_api_url = "example.test:443"
admin_api_url = "example.test:443"
user_id = "user-1"
synchronizer_id = "global-domain::abc"

[ledger.auth]
type = "insecureStatic"
jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

[parties]
operator = "op::1"
provider = "prov::1"
registrar = "reg::1"
`), 0o644))

	cfg, err := Load(cfgPath)
	require.NoError(t, err)
	require.Equal(t, "https://api.utilities.digitalasset-dev.com/api/utilities", cfg.Operator.BaseURL)
}
