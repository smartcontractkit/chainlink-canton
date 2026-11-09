package linktokenadmin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

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

// countingEventQuery fails every contract lookup and counts them.
type countingEventQuery struct {
	apiv2.EventQueryServiceClient
	calls int
}

func (c *countingEventQuery) GetEventsByContractId(context.Context, *apiv2.GetEventsByContractIdRequest, ...grpc.CallOption) (*apiv2.GetEventsByContractIdResponse, error) {
	c.calls++

	return nil, errors.New("ledger lookup")
}

func TestResolveFactorySuppliedDisclosures(t *testing.T) {
	t.Parallel()

	blob := []byte("burn-mint-factory-blob")
	sync := "sync::1"
	in := ExecuteInput{Kind: KindMint, Registry: RegistryInput{
		FactoryCid: "factory-cid",
		Disclosed: []daRegistry.DisclosedContract{
			{
				ContractId:       "factory-cid",
				TemplateId:       "pkg-id:Utility.Registry.App.V0.Service.BurnMintFactory:BurnMintFactory",
				CreatedEventBlob: base64.StdEncoding.EncodeToString(blob),
				SynchronizerId:   &sync,
			},
			{
				ContractId:       "inst-config-cid",
				TemplateId:       "pkg-id:Utility.Registry.V0.Configuration.Instrument:InstrumentConfiguration",
				CreatedEventBlob: base64.StdEncoding.EncodeToString([]byte("instrument-configuration-blob")),
			},
		},
	}}
	eq := &countingEventQuery{}
	p := canton.Participant{LedgerServices: canton.LedgerServiceClients{EventQuery: eq}}

	factory, err := resolveFactory(t.Context(), p, in, linkAuth("1.0"))
	require.NoError(t, err)
	require.Zero(t, eq.calls, "supplied disclosures need no ledger lookup")
	require.Equal(t, "factory-cid", factory.factoryCid)
	require.Len(t, factory.disclosed, 2)
	require.Equal(t, "factory-cid", factory.disclosed[0].GetContractId())
	require.Equal(t, blob, factory.disclosed[0].GetCreatedEventBlob())
	require.Equal(t, "sync::1", factory.disclosed[0].GetSynchronizerId())
	require.Equal(t, "BurnMintFactory", factory.disclosed[0].GetTemplateId().GetEntityName())
	require.Equal(t, "inst-config-cid", factory.disclosed[1].GetContractId())

	// Without supplied disclosures the same input reads the factory from the participant.
	in.Registry.Disclosed = nil
	_, err = resolveFactory(t.Context(), p, in, linkAuth("1.0"))
	require.ErrorContains(t, err, "ledger lookup")
	require.Equal(t, 1, eq.calls)
}

func TestResolveFactoryBadSuppliedBlob(t *testing.T) {
	t.Parallel()

	in := ExecuteInput{Kind: KindMint, Registry: RegistryInput{
		FactoryCid: "factory-cid",
		Disclosed: []daRegistry.DisclosedContract{{
			ContractId:       "factory-cid",
			TemplateId:       "pkg-id:Utility.Registry.App.V0.Service.BurnMintFactory:BurnMintFactory",
			CreatedEventBlob: "not base64!",
		}},
	}}
	eq := &countingEventQuery{}

	_, err := resolveFactory(t.Context(), canton.Participant{LedgerServices: canton.LedgerServiceClients{EventQuery: eq}}, in, linkAuth("1.0"))
	require.ErrorContains(t, err, "decode blob")
	require.Zero(t, eq.calls)
}
