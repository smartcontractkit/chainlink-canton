package linktokenadmin

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"

	spliceHolding "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/daRegistry"
)

const burnMintFactoryPath = "/api/utilities/v0/registry/burn-mint-instruction/v0/burn-mint-factory"

func linkAuth(amount string) authorization {
	return authorization{
		ccipOwner:    ccipOwner,
		instrumentID: spliceHolding.InstrumentId{Admin: ccipOwner, Id: "LINK"},
		counterparty: recipient,
		amount:       types.NUMERIC(amount),
	}
}

func TestBurnMintRequestBody(t *testing.T) {
	t.Parallel()

	mint, err := json.Marshal(burnMintRequest(KindMint, linkAuth("10.0"), nil))
	require.NoError(t, err)
	require.JSONEq(t, `{
		"instrumentId": {"admin": "`+ccipOwner+`", "id": "LINK"},
		"inputHoldingCids": [],
		"outputs": [{"owner": "`+recipient+`", "amount": "10.0"}]
	}`, string(mint))

	burn, err := json.Marshal(burnMintRequest(KindBurn, linkAuth("10.0"), []string{"h1", "h2"}))
	require.NoError(t, err)
	require.JSONEq(t, `{
		"instrumentId": {"admin": "`+ccipOwner+`", "id": "LINK"},
		"inputHoldingCids": ["h1", "h2"],
		"outputs": []
	}`, string(burn))
}

func TestFetchBurnMintFactory(t *testing.T) {
	t.Parallel()

	blob := []byte("instrument-configuration-blob")
	var gotBody daRegistry.GetBurnMintFactoryRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != burnMintFactoryPath {
			http.NotFound(w, r)

			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || json.Unmarshal(body, &gotBody) != nil {
			http.Error(w, "bad body", http.StatusBadRequest)

			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"factoryId": "factory-cid",
			"choiceContext": {
				"choiceContextData": {"values": {
					"utility.digitalasset.com/instrument-configuration": {"tag": "AV_ContractId", "value": "inst-config-cid"},
					"utility.digitalasset.com/issuer-credentials": {"tag": "AV_List", "value": []}
				}},
				"disclosedContracts": [{
					"contractId": "inst-config-cid",
					"templateId": "pkg-id:Utility.Registry.V0.Configuration.Instrument:InstrumentConfiguration",
					"createdEventBlob": "` + base64.StdEncoding.EncodeToString(blob) + `",
					"synchronizerId": "sync::1"
				}]
			}
		}`))
	}))
	t.Cleanup(server.Close)

	client, err := daRegistry.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	req := burnMintRequest(KindMint, linkAuth("10.0"), nil)
	factory, err := fetchBurnMintFactory(t.Context(), client, req)
	require.NoError(t, err)

	require.Equal(t, req, gotBody, "registry receives the mint output")
	require.Equal(t, "factory-cid", factory.factoryCid)
	require.Contains(t, factory.context.Values, "utility.digitalasset.com/instrument-configuration")
	require.Contains(t, factory.context.Values, "utility.digitalasset.com/issuer-credentials")
	instConfig := factory.context.Values["utility.digitalasset.com/instrument-configuration"]
	require.NotNil(t, instConfig.AVContractId)
	require.Equal(t, "inst-config-cid", string(*instConfig.AVContractId))

	require.Len(t, factory.disclosed, 1)
	d := factory.disclosed[0]
	require.Equal(t, "inst-config-cid", d.GetContractId())
	require.Equal(t, blob, d.GetCreatedEventBlob())
	require.Equal(t, "sync::1", d.GetSynchronizerId())
	require.Equal(t, "pkg-id", d.GetTemplateId().GetPackageId())
	require.Equal(t, "InstrumentConfiguration", d.GetTemplateId().GetEntityName())
}

func TestFetchBurnMintFactoryErrorStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"unknown instrument"}`, http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client, err := daRegistry.NewClientWithResponses(server.URL)
	require.NoError(t, err)

	_, err = fetchBurnMintFactory(t.Context(), client, burnMintRequest(KindBurn, linkAuth("4.0"), []string{"h1"}))
	require.ErrorContains(t, err, "status 404")
	require.ErrorContains(t, err, "unknown instrument")
}
