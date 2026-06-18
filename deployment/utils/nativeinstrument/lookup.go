package nativeinstrument

import (
	"context"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	tokenMetadataV1 "github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
)

const DefaultNativeTokenInstrumentID = "Amulet"

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
