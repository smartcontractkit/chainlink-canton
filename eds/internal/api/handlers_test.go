package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/ledger"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/types"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

// MockContractQuerier for testing
type MockContractQuerier struct {
	Contracts []*ledger.ActiveContract
	Error     error
}

func (m *MockContractQuerier) GetAllContractsForParty(ctx context.Context, party string) ([]*ledger.ActiveContract, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Contracts, nil
}

func createMockContract(moduleName, entityName, environmentID, contractID string) *ledger.ActiveContract {
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
						Label: "environmentId",
						Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: environmentID}},
					},
				},
			},
			CreatedEventBlob: []byte("test-blob"),
		},
		SynchronizerID: "test-sync",
	}
}

func setupTestHandlers() (*Handlers, *config.EnvironmentsConfig, *MockContractQuerier) {
	envConfig := &config.EnvironmentsConfig{
		Environments: map[string]config.EnvironmentConfig{
			"mainnet-v1": {
				Party:       "ccip-owner::123",
				Description: "Production mainnet",
				Contracts: config.ContractIdentifiers{
					Router:             "mainnet-v1",
					OnRamp:             "mainnet-v1",
					FeeQuoter:          "mainnet-v1",
					OffRamp:            "mainnet-v1",
					CCV:                "mainnet-v1",
					TokenAdminRegistry: "mainnet-v1",
				},
			},
			"testnet": {
				Party:       "ccip-test::456",
				Description: "Test network",
				Contracts: config.ContractIdentifiers{
					Router:    "testnet",
					OnRamp:    "testnet",
					FeeQuoter: "testnet",
				},
			},
		},
	}

	mockQuerier := &MockContractQuerier{
		Contracts: []*ledger.ActiveContract{
			createMockContract("CCIP.Router", "Router", "mainnet-v1", "router-123"),
			createMockContract("CCIP.OnRamp", "OnRamp", "mainnet-v1", "onramp-456"),
			createMockContract("CCIP.FeeQuoter", "FeeQuoter", "mainnet-v1", "feequoter-789"),
			createMockContract("CCIP.OffRamp", "OffRamp", "mainnet-v1", "offramp-111"),
			createMockContract("CCIP.CommitteeVerifier", "CommitteeVerifier", "mainnet-v1", "ccv-222"),
			createMockContract("CCIP.TokenAdminRegistry", "TokenAdminRegistry", "mainnet-v1", "tar-333"),
		},
	}

	svc := disclosure.NewService(mockQuerier, envConfig)
	handlers := NewHandlers(svc, envConfig)

	return handlers, envConfig, mockQuerier
}

func TestHandlers_GetCCIPSendDisclosures(t *testing.T) {
	handlers, _, _ := setupTestHandlers()

	t.Run("returns disclosures for valid environment", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ccip/mainnet-v1/disclosures/send", nil)
		req = mux.SetURLVars(req, map[string]string{"environmentId": "mainnet-v1"})
		w := httptest.NewRecorder()

		handlers.GetCCIPSendDisclosures(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response types.CCIPSendDisclosures
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "mainnet-v1", response.EnvironmentID)
		assert.NotNil(t, response.Contracts.Router)
		assert.NotNil(t, response.Contracts.OnRamp)
		assert.NotNil(t, response.Contracts.FeeQuoter)
		assert.Equal(t, "router-123", response.Contracts.Router.ContractID)
	})

	t.Run("returns 404 for unknown environment", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ccip/unknown/disclosures/send", nil)
		req = mux.SetURLVars(req, map[string]string{"environmentId": "unknown"})
		w := httptest.NewRecorder()

		handlers.GetCCIPSendDisclosures(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response types.ErrorResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "DISCLOSURES_NOT_FOUND", response.Code)
	})
}

func TestHandlers_GetCCIPExecuteDisclosures(t *testing.T) {
	handlers, _, _ := setupTestHandlers()

	t.Run("returns disclosures for valid environment", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ccip/mainnet-v1/disclosures/execute", nil)
		req = mux.SetURLVars(req, map[string]string{"environmentId": "mainnet-v1"})
		w := httptest.NewRecorder()

		handlers.GetCCIPExecuteDisclosures(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response types.CCIPExecuteDisclosures
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "mainnet-v1", response.EnvironmentID)
		assert.NotNil(t, response.Contracts.OffRamp)
		assert.NotNil(t, response.Contracts.CCV)
		assert.NotNil(t, response.Contracts.TokenAdminRegistry)
	})
}

func TestHandlers_ListEnvironments(t *testing.T) {
	handlers, _, _ := setupTestHandlers()

	req := httptest.NewRequest("GET", "/api/v1/ccip/environments", nil)
	w := httptest.NewRecorder()

	handlers.ListEnvironments(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.EnvironmentsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response.Environments, 2)

	// Find mainnet in response
	var foundMainnet, foundTestnet bool
	for _, env := range response.Environments {
		if env.ID == "mainnet-v1" {
			foundMainnet = true
			assert.Equal(t, "ccip-owner::123", env.Party)
			assert.Equal(t, "Production mainnet", env.Description)
		}
		if env.ID == "testnet" {
			foundTestnet = true
			assert.Equal(t, "ccip-test::456", env.Party)
		}
	}
	assert.True(t, foundMainnet)
	assert.True(t, foundTestnet)
}

func TestHandlers_Health(t *testing.T) {
	handlers, _, _ := setupTestHandlers()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	handlers.Health(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.True(t, response.LedgerAPIConnected)
	assert.Len(t, response.Environments, 2)
	assert.Contains(t, response.Environments, "mainnet-v1")
	assert.Contains(t, response.Environments, "testnet")
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"key": "value"}
	WriteJSON(w, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "value", response["key"])
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Something went wrong")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response types.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Something went wrong", response.Error)
	assert.Equal(t, "BAD_REQUEST", response.Code)
}
