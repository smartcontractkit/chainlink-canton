package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/ledger"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/types"
)

type mockQuerier struct {
	contracts []*ledger.ActiveContract
}

func (m *mockQuerier) GetAllContractsForParty(ctx context.Context, party string) ([]*ledger.ActiveContract, error) {
	return m.contracts, nil
}

func createContract(moduleName, entityName, environmentID, contractID string) *ledger.ActiveContract {
	return &ledger.ActiveContract{
		CreatedEvent: &apiv2.CreatedEvent{
			ContractId: contractID,
			TemplateId: &apiv2.Identifier{
				PackageId:  "test-pkg",
				ModuleName: moduleName,
				EntityName: entityName,
			},
			CreateArguments: &apiv2.Record{
				Fields: []*apiv2.RecordField{
					{
						Label: "instanceId",
						Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: environmentID}},
					},
				},
			},
			CreatedEventBlob: []byte("test-blob"),
		},
		SynchronizerID: "test-sync",
	}
}

func setupTestRouter() http.Handler {
	envConfig := &config.EnvironmentsConfig{
		Environments: map[string]config.EnvironmentConfig{
			"mainnet-v1": {
				Party:       "ccip-owner::123",
				Description: "Production mainnet",
				Contracts: config.ContractIdentifiers{
					Router:             "mainnet-v1-router",
					OnRamp:             "mainnet-v1-onramp",
					FeeQuoter:          "mainnet-v1-feequoter",
					OffRamp:            "mainnet-v1-offramp",
					CCV:                "mainnet-v1-ccv",
					TokenAdminRegistry: "mainnet-v1-tar",
				},
			},
		},
	}

	mockQ := &mockQuerier{
		contracts: []*ledger.ActiveContract{
			createContract("CCIP.Router", "Router", "mainnet-v1-router", "router-123"),
			createContract("CCIP.OnRamp", "OnRamp", "mainnet-v1-onramp", "onramp-456"),
			createContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1-feequoter", "feequoter-789"),
			createContract("CCIP.OffRamp", "OffRamp", "mainnet-v1-offramp", "offramp-111"),
			createContract("CCIP.CommitteeVerifier", "CommitteeVerifier", "mainnet-v1-ccv", "ccv-222"),
			createContract("CCIP.TokenAdminRegistry", "TokenAdminRegistry", "mainnet-v1-tar", "tar-333"),
		},
	}

	svc := disclosure.NewService(mockQ, envConfig)

	return NewRouter(svc, envConfig)
}

func TestRouter_HealthEndpoint(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
}

func TestRouter_EnvironmentsEndpoint(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ccip/environments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.EnvironmentsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Len(t, response.Environments, 1)
}

func TestRouter_CCIPSendDisclosuresEndpoint(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ccip/mainnet-v1/disclosures/send", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.CCIPSendDisclosures
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "mainnet-v1", response.EnvironmentID)
}

func TestRouter_CCIPExecuteDisclosuresEndpoint(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ccip/mainnet-v1/disclosures/execute", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.CCIPExecuteDisclosures
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "mainnet-v1", response.EnvironmentID)
}

func TestRouter_CORSHeaders(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
}

// options preflight handling is typically done by a reverse proxy in production

func TestRouter_NotFound(t *testing.T) {
	t.Parallel()

	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
