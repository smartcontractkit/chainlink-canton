package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smartcontractkit/chainlink-canton/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/types"
)

type Handlers struct {
	disclosureSvc *disclosure.Service
	envConfig     *config.EnvironmentsConfig
}

func NewHandlers(disclosureSvc *disclosure.Service, envConfig *config.EnvironmentsConfig) *Handlers {
	return &Handlers{
		disclosureSvc: disclosureSvc,
		envConfig:     envConfig,
	}
}

func (h *Handlers) GetCCIPSendDisclosures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	environmentID := vars["environmentId"]

	disclosures, err := h.disclosureSvc.GetCCIPSendDisclosures(r.Context(), environmentID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "DISCLOSURES_NOT_FOUND", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, disclosures)
}

func (h *Handlers) GetCCIPExecuteDisclosures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	environmentID := vars["environmentId"]

	disclosures, err := h.disclosureSvc.GetCCIPExecuteDisclosures(r.Context(), environmentID)
	if err != nil {
		WriteError(w, http.StatusNotFound, "DISCLOSURES_NOT_FOUND", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, disclosures)
}

func (h *Handlers) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	envs := make([]types.EnvironmentInfo, 0, len(h.envConfig.Environments))
	for id, env := range h.envConfig.Environments {
		envs = append(envs, types.EnvironmentInfo{
			ID:          id,
			Party:       env.Party,
			Description: env.Description,
		})
	}

	WriteJSON(w, http.StatusOK, types.EnvironmentsResponse{
		Environments: envs,
	})
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	envNames := make([]string, 0, len(h.envConfig.Environments))
	for name := range h.envConfig.Environments {
		envNames = append(envNames, name)
	}

	WriteJSON(w, http.StatusOK, types.HealthResponse{
		Status:             "healthy",
		LedgerAPIConnected: true,
		Environments:       envNames,
	})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	//nolint:errchkjson // temp
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, types.ErrorResponse{
		Error: message,
		Code:  code,
	})
}
