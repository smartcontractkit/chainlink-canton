package testhelpers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	apiv2 "github.com/smartcontractkit/chainlink-canton/pb/gen/com/daml/ledger/api/v2"
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

func GetAmuletRulesContract(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (string, *apiv2.Identifier, error) {
	// Get Amulet Rules contract
	amuletRulesResponse, err := scanProxyClient.GetAmuletRulesWithResponse(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("error getting amulet rules response: %w", err)
	}
	if amuletRulesResponse.StatusCode() != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d: %v", amuletRulesResponse.StatusCode(), amuletRulesResponse.Body)
	}
	amuletRulesId, err := TemplateIdFromString(amuletRulesResponse.JSON200.AmuletRules.Contract.TemplateId)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse amulet rules template id: %w", err)
	}

	return amuletRulesResponse.JSON200.AmuletRules.Contract.ContractId, amuletRulesId, nil
}

func GetFirstOpenMiningRound(ctx context.Context, scanProxyClient scanProxy.ClientWithResponsesInterface) (string, error) {
	openMiningRoundResponse, err := scanProxyClient.GetOpenAndIssuingMiningRoundsWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting open mining rounds response: %w", err)
	}
	if openMiningRoundResponse.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d: %v", openMiningRoundResponse.StatusCode(), openMiningRoundResponse.Body)
	}
	slices.SortFunc(openMiningRoundResponse.JSON200.OpenMiningRounds, func(a, b scanProxy.ContractWithState) int {
		aOpen, _ := time.Parse(time.RFC3339, a.Contract.Payload["opensAt"].(string))
		bOpen, _ := time.Parse(time.RFC3339, b.Contract.Payload["opensAt"].(string))

		return int(aOpen.UnixMilli() - bOpen.UnixMilli())
	})
	var openMiningRoundCid string
	for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
		opensAt, err := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
		if err != nil {
			return "", fmt.Errorf("failed to parse opensAt %q: %w", round.Contract.Payload["opensAt"], err)
		}
		targetClosesAt, err := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
		if err != nil {
			return "", fmt.Errorf("failed to parse targetClosesAt %q: %w", round.Contract.Payload["targetClosesAt"], err)
		}
		if opensAt.Before(time.Now()) && targetClosesAt.After(time.Now()) {
			openMiningRoundCid = round.Contract.ContractId
		}
	}

	return openMiningRoundCid, nil
}
