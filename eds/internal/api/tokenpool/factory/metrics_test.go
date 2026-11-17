package factory

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	oapiCommon "github.com/smartcontractkit/chainlink-canton/eds/api/common"
	"github.com/smartcontractkit/chainlink-canton/eds/config"
)

type recordedDARequest struct {
	operation string
	success   bool
}

// recordingMetrics captures every DA request recorded by a factory.
type recordingMetrics struct {
	mu   sync.Mutex
	reqs []recordedDARequest
}

func (r *recordingMetrics) RecordDARequest(_ context.Context, operation string, success bool, _ time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqs = append(r.reqs, recordedDARequest{operation: operation, success: success})
}

func (r *recordingMetrics) recorded() []recordedDARequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]recordedDARequest(nil), r.reqs...)
}

const (
	validFactoryResponse  = `{"factoryId":"factory-cid","transferKind":"direct","choiceContext":{"choiceContextData":{},"disclosedContracts":[]}}`
	validHoldingsResponse = `{"disclosedHoldings":[]}`
)

// daCase describes one response of the fake DA server and the result that must be recorded for it.
type daCase struct {
	name        string
	status      int
	body        string
	wantSuccess bool
}

func daCases(validBody string) []daCase {
	return []daCase{
		{name: "200 with valid body is a success", status: http.StatusOK, body: validBody, wantSuccess: true},
		{name: "404 is a failure", status: http.StatusNotFound, body: `{"error":"not found"}`, wantSuccess: false},
		{name: "500 is a failure", status: http.StatusInternalServerError, body: `oops`, wantSuccess: false},
		{name: "200 with malformed body is a failure", status: http.StatusOK, body: `not json`, wantSuccess: false},
	}
}

func fakeDA(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return server
}

func urlFactoryConfig(url string) config.Factory {
	return config.Factory{Type: config.FactoryTypeURL, TokenStandardURL: &url}
}

func TestTransferFactoryRecordsDARequest(t *testing.T) {
	t.Parallel()
	for _, tc := range daCases(validFactoryResponse) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := fakeDA(t, tc.status, tc.body)
			metrics := &recordingMetrics{}
			f, err := NewTransferFactory(t.Context(), "owner", nil, urlFactoryConfig(server.URL), WithMetrics(metrics))
			require.NoError(t, err)

			_, _, _, err = f.GetFeeTokenDisclosures(t.Context(), splice_api_token_holding_v1.InstrumentId{Admin: "admin", Id: "id"}, "sender", "receiver")
			if tc.wantSuccess {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			assert.Equal(t, []recordedDARequest{{operation: OperationTransferFactory, success: tc.wantSuccess}}, metrics.recorded())
		})
	}
}

func TestBurnMintFactoryRecordsDARequest(t *testing.T) {
	t.Parallel()
	message := oapiCommon.Message{TokenTransfer: &oapiCommon.TokenTransfer{
		Amount: "1.0",
		Token:  oapiCommon.InstrumentId{Admin: "admin", Id: "id"},
	}}

	for _, tc := range daCases(validFactoryResponse) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := fakeDA(t, tc.status, tc.body)
			metrics := &recordingMetrics{}
			f, err := NewBurnMintFactory(t.Context(), "owner", nil, urlFactoryConfig(server.URL), WithMetrics(metrics))
			require.NoError(t, err)

			_, _, _, err = f.GetSendDisclosures(t.Context(), message, nil)
			if tc.wantSuccess {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			assert.Equal(t, []recordedDARequest{{operation: OperationBurnMintFactory, success: tc.wantSuccess}}, metrics.recorded())
		})
	}
}

func TestBridgeHoldingsRecordsDARequest(t *testing.T) {
	t.Parallel()
	for _, tc := range daCases(validHoldingsResponse) {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := fakeDA(t, tc.status, tc.body)
			cfg := urlFactoryConfig(server.URL)
			cfg.BridgeHoldingsURL = &server.URL
			metrics := &recordingMetrics{}
			f, err := NewTransferFactory(t.Context(), "owner", nil, cfg, WithMetrics(metrics))
			require.NoError(t, err)

			_, err = f.(BridgeHoldingsDisclosureFactory).GetBridgeHoldings(t.Context(), "owner", splice_api_token_holding_v1.InstrumentId{Id: "id"})
			if tc.wantSuccess {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			assert.Equal(t, []recordedDARequest{{operation: OperationBridgeHoldings, success: tc.wantSuccess}}, metrics.recorded())
		})
	}
}

// Without a configured bridge URL no request is made, so nothing must be recorded.
func TestBridgeHoldingsNotConfiguredRecordsNothing(t *testing.T) {
	t.Parallel()
	server := fakeDA(t, http.StatusOK, validHoldingsResponse)
	metrics := &recordingMetrics{}
	f, err := NewTransferFactory(t.Context(), "owner", nil, urlFactoryConfig(server.URL), WithMetrics(metrics))
	require.NoError(t, err)

	_, err = f.(BridgeHoldingsDisclosureFactory).GetBridgeHoldings(t.Context(), "owner", splice_api_token_holding_v1.InstrumentId{Id: "id"})
	require.Error(t, err)

	assert.Empty(t, metrics.recorded())
}

// Factories created without WithMetrics must keep working.
func TestFactoriesWithoutMetrics(t *testing.T) {
	t.Parallel()
	server := fakeDA(t, http.StatusOK, validFactoryResponse)
	f, err := NewTransferFactory(t.Context(), "owner", nil, urlFactoryConfig(server.URL))
	require.NoError(t, err)

	_, _, _, err = f.GetFeeTokenDisclosures(t.Context(), splice_api_token_holding_v1.InstrumentId{Admin: "admin", Id: "id"}, "sender", "receiver")
	require.NoError(t, err)
}
