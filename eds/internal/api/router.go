package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton-internal/eds/internal/disclosure"
)

func NewRouter(disclosureSvc *disclosure.Service, envConfig *config.EnvironmentsConfig) http.Handler {
	r := mux.NewRouter()
	handlers := NewHandlers(disclosureSvc, envConfig)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/health", handlers.Health).Methods("GET")
	api.HandleFunc("/ccip/environments", handlers.ListEnvironments).Methods("GET")
	api.HandleFunc("/ccip/{environmentId}/disclosures/send", handlers.GetCCIPSendDisclosures).Methods("GET")
	api.HandleFunc("/ccip/{environmentId}/disclosures/execute", handlers.GetCCIPExecuteDisclosures).Methods("GET")

	r.Use(LoggingMiddleware)
	r.Use(CORSMiddleware)

	return r
}
