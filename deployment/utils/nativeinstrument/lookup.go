package nativeinstrument

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
)

const DefaultNativeTokenInstrumentID = "Amulet"

const (
	instrumentAdminLabelPrefix = "instrument-admin:"
	instrumentIDLabelPrefix    = "instrument-id:"
)

// ResolveNativeInstrumentID returns the native fee token instrument. It prefers a Token
// address ref already in the datastore (Amulet registered in TAR) so lane configure and
// re-runs do not require validator scan-proxy access. Falls back to scan-proxy registry
// metadata when no ref is present (first-time bootstrap).
func ResolveNativeInstrumentID(
	ctx context.Context,
	participant canton.Participant,
	ds datastore.DataStore,
	chainSelector uint64,
) (splice_api_token_holding_v1.InstrumentId, error) {
	if ds != nil && chainSelector != 0 {
		if id, ok := nativeInstrumentFromDataStore(ds, chainSelector); ok {
			return id, nil
		}
	}

	return LookupNativeInstrumentID(ctx, participant)
}

func nativeInstrumentFromDataStore(ds datastore.DataStore, chainSelector uint64) (splice_api_token_holding_v1.InstrumentId, bool) {
	refs := ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType("Token")),
	)

	for _, ref := range refs {
		if ref.Qualifier != DefaultNativeTokenInstrumentID && !ref.Labels.Contains(instrumentIDLabelPrefix+DefaultNativeTokenInstrumentID) {
			continue
		}
		admin, id, ok := instrumentFromLabels(ref.Labels)
		if !ok {
			continue
		}

		return splice_api_token_holding_v1.InstrumentId{
			Admin: types.PARTY(admin),
			Id:    types.TEXT(id),
		}, true
	}

	return splice_api_token_holding_v1.InstrumentId{}, false
}

func instrumentFromLabels(labels datastore.LabelSet) (admin, id string, ok bool) {
	for _, label := range labels.List() {
		switch {
		case strings.HasPrefix(label, instrumentAdminLabelPrefix):
			admin = strings.TrimPrefix(label, instrumentAdminLabelPrefix)
		case strings.HasPrefix(label, instrumentIDLabelPrefix):
			id = strings.TrimPrefix(label, instrumentIDLabelPrefix)
		}
	}
	if admin == "" || id == "" {
		return "", "", false
	}

	return admin, id, true
}

// LookupNativeInstrumentID resolves the Canton network native token (Amulet) admin party
// from the participant validator scan-proxy registry metadata API.
func LookupNativeInstrumentID(ctx context.Context, participant canton.Participant) (splice_api_token_holding_v1.InstrumentId, error) {
	tokenSource := participant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}

	client, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(interceptor),
	)
	if err != nil {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("failed to create token metadata client: %w", err)
	}

	info, err := client.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("error getting registry info: %w", err)
	}
	if info.StatusCode() != http.StatusOK {
		return splice_api_token_holding_v1.InstrumentId{}, fmt.Errorf("unexpected status code from token metadata client: %d: %v", info.StatusCode(), info.Body)
	}

	return splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(info.JSON200.AdminId),
		Id:    types.TEXT(DefaultNativeTokenInstrumentID),
	}, nil
}
