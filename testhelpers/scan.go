package testhelpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

func TemplateIdFromString(s string) (*apiv2.Identifier, error) {
	split := strings.Split(s, ":")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid template id format: %s", s)
	}

	return &apiv2.Identifier{
		PackageId:  split[0],
		ModuleName: split[1],
		EntityName: split[2],
	}, nil
}

func ContractToDisclosedContract(contract scanProxy.ContractWithState) (*apiv2.DisclosedContract, error) {
	templateId, err := TemplateIdFromString(contract.Contract.TemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template id: %w", err)
	}
	createdEventBlob, err := base64.StdEncoding.DecodeString(contract.Contract.CreatedEventBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to decode created event blob: %w", err)
	}

	return &apiv2.DisclosedContract{
		TemplateId:       templateId,
		ContractId:       contract.Contract.ContractId,
		CreatedEventBlob: createdEventBlob,
		SynchronizerId:   *contract.DomainId,
	}, nil
}

func GetAmuletRulesContract(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (string, *apiv2.DisclosedContract, error) {
	// Get Amulet Rules contract
	amuletRulesResponse, err := scanProxyClient.GetDsoInfoWithResponse(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("error getting amulet rules response: %w", err)
	}
	if amuletRulesResponse.StatusCode() != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d: %v", amuletRulesResponse.StatusCode(), amuletRulesResponse.Body)
	}

	amuletRules, err := ContractToDisclosedContract(amuletRulesResponse.JSON200.AmuletRules)

	return amuletRulesResponse.JSON200.DsoPartyId, amuletRules, err
}

func GetFirstOpenMiningRound(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (*apiv2.DisclosedContract, error) {
	openMiningRoundResponse, err := scanProxyClient.GetOpenAndIssuingMiningRoundsWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting open mining rounds response: %w", err)
	}
	if openMiningRoundResponse.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %v", openMiningRoundResponse.StatusCode(), openMiningRoundResponse.Body)
	}
	slices.SortFunc(openMiningRoundResponse.JSON200.OpenMiningRounds, func(a, b scanProxy.ContractWithState) int {
		aOpen, _ := time.Parse(time.RFC3339, a.Contract.Payload["opensAt"].(string))
		bOpen, _ := time.Parse(time.RFC3339, b.Contract.Payload["opensAt"].(string))

		return int(aOpen.UnixMilli() - bOpen.UnixMilli())
	})

	for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
		opensAt, err := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse opensAt %q: %w", round.Contract.Payload["opensAt"], err)
		}
		targetClosesAt, err := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse targetClosesAt %q: %w", round.Contract.Payload["targetClosesAt"], err)
		}
		if opensAt.Before(time.Now()) && targetClosesAt.After(time.Now()) {
			return ContractToDisclosedContract(round)
		}
	}

	return nil, fmt.Errorf("failed to find open mining round contract")
}

func NewValidatorAPIClients(
	participant canton.Participant,
) (
	scanProxy.ClientWithResponsesInterface,
	tokenMetadataV1.ClientWithResponsesInterface,
	transferInstructionV1.ClientWithResponsesInterface,
	error,
) {
	tokenSource := participant.TokenSource
	interceptor := func(ctx context.Context, req *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}

	scanProxyClient, err := scanProxy.NewClientWithResponses(
		participant.Endpoints.ValidatorAPIURL,
		scanProxy.WithRequestEditorFn(interceptor),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create scan proxy client: %w", err)
	}
	tokenMetadataClient, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(interceptor),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create token metadata client: %w", err)
	}
	transferInstructionClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(interceptor),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create transfer instruction client: %w", err)
	}

	return scanProxyClient, tokenMetadataClient, transferInstructionClient, nil
}

func ResolveRegistryAdmin(ctx context.Context, participant canton.Participant) (string, error) {
	_, tokenMetadataClient, _, err := NewValidatorAPIClients(participant)
	if err != nil {
		return "", err
	}

	registryAdmin, err := GetRegistryAdmin(ctx, tokenMetadataClient)
	if err != nil {
		return "", fmt.Errorf("get registry admin: %w", err)
	}

	return registryAdmin, nil
}
