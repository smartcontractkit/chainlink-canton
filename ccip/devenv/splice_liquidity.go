package devenv

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
)

func seedAMTLiquidity(ctx context.Context, participant canton.Participant, ownerParty string, amount string) error {
	requestEditor := func(reqCtx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("retrieve participant token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}

	scanClient, err := scanProxy.NewClientWithResponses(participant.Endpoints.ValidatorAPIURL, scanProxy.WithRequestEditorFn(requestEditor))
	if err != nil {
		return fmt.Errorf("create scan proxy client: %w", err)
	}
	metadataClient, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return fmt.Errorf("create token metadata client: %w", err)
	}
	transferClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return fmt.Errorf("create transfer instruction client: %w", err)
	}

	registryInfoResponse, err := metadataClient.GetRegistryInfoWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("get registry info: %w", err)
	}
	if registryInfoResponse.StatusCode() != http.StatusOK || registryInfoResponse.JSON200 == nil || registryInfoResponse.JSON200.AdminId == "" {
		return fmt.Errorf("unexpected registry info response status=%d", registryInfoResponse.StatusCode())
	}
	registryAdmin := registryInfoResponse.JSON200.AdminId

	transferFactoryResponse, err := transferClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": registryAdmin,
			"transfer": map[string]any{
				"sender":   registryAdmin,
				"receiver": ownerParty,
				"amount":   amount,
				"instrumentId": map[string]any{
					"admin": registryAdmin,
					"id":    "Amulet",
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"inputHoldingCids": []string{},
				"meta":             map[string]any{"values": map[string]any{}},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{"values": map[string]any{}},
				"meta":    map[string]any{"values": map[string]any{}},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("get transfer factory for mint: %w", err)
	}
	if transferFactoryResponse.StatusCode() != http.StatusOK || transferFactoryResponse.JSON200 == nil {
		return fmt.Errorf("unexpected transfer factory response status=%d", transferFactoryResponse.StatusCode())
	}
	disclosedContracts := make([]*apiv2.DisclosedContract, 0, len(transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts))
	for _, contract := range transferFactoryResponse.JSON200.ChoiceContext.DisclosedContracts {
		id, parseErr := TemplateIdFromString(contract.TemplateId)
		if parseErr != nil {
			return fmt.Errorf("parse transfer factory template id: %w", parseErr)
		}
		createdEventBlob, decodeErr := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if decodeErr != nil {
			return fmt.Errorf("decode transfer factory created event blob: %w", decodeErr)
		}
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}

	amuletRulesResponse, err := scanClient.GetAmuletRulesWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("get amulet rules: %w", err)
	}
	if amuletRulesResponse.StatusCode() != http.StatusOK || amuletRulesResponse.JSON200 == nil {
		return fmt.Errorf("unexpected amulet rules response status=%d", amuletRulesResponse.StatusCode())
	}
	amuletRulesID, err := TemplateIdFromString(amuletRulesResponse.JSON200.AmuletRules.Contract.TemplateId)
	if err != nil {
		return fmt.Errorf("parse amulet rules template id: %w", err)
	}
	amuletRulesCID := amuletRulesResponse.JSON200.AmuletRules.Contract.ContractId

	getActiveOpenMiningRoundCID := func() (string, error) {
		openMiningRoundResponse, roundErr := scanClient.GetOpenAndIssuingMiningRoundsWithResponse(ctx)
		if roundErr != nil {
			return "", fmt.Errorf("get open mining rounds: %w", roundErr)
		}
		if openMiningRoundResponse.StatusCode() != http.StatusOK || openMiningRoundResponse.JSON200 == nil {
			return "", fmt.Errorf("unexpected open mining rounds response status=%d", openMiningRoundResponse.StatusCode())
		}
		slices.SortFunc(openMiningRoundResponse.JSON200.OpenMiningRounds, func(a, b scanProxy.ContractWithState) int {
			aOpen, _ := time.Parse(time.RFC3339, a.Contract.Payload["opensAt"].(string))
			bOpen, _ := time.Parse(time.RFC3339, b.Contract.Payload["opensAt"].(string))
			return int(aOpen.UnixMilli() - bOpen.UnixMilli())
		})
		var openMiningRoundCID string
		for _, round := range openMiningRoundResponse.JSON200.OpenMiningRounds {
			opensAt, parseOpenErr := time.Parse(time.RFC3339, round.Contract.Payload["opensAt"].(string))
			if parseOpenErr != nil {
				continue
			}
			targetClosesAt, parseCloseErr := time.Parse(time.RFC3339, round.Contract.Payload["targetClosesAt"].(string))
			if parseCloseErr != nil {
				continue
			}
			now := time.Now()
			if opensAt.Before(now) && targetClosesAt.After(now) {
				openMiningRoundCID = round.Contract.ContractId
			}
		}
		if openMiningRoundCID == "" {
			return "", fmt.Errorf("no active open mining round found")
		}
		return openMiningRoundCID, nil
	}
	openMiningRoundCID, err := getActiveOpenMiningRoundCID()
	if err != nil {
		return err
	}

	submitTap := func(disclosures []*apiv2.DisclosedContract) error {
		_, submitErr := participant.LedgerServices.Command.SubmitAndWaitForTransaction(ctx, &apiv2.SubmitAndWaitForTransactionRequest{
			Commands: &apiv2.Commands{
				CommandId: uuid.Must(uuid.NewUUID()).String(),
				Commands: []*apiv2.Command{
					{
						Command: &apiv2.Command_Exercise{
							Exercise: &apiv2.ExerciseCommand{
								TemplateId: amuletRulesID,
								ContractId: amuletRulesCID,
								Choice:     "AmuletRules_DevNet_Tap",
								ChoiceArgument: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
									{Label: "receiver", Value: &apiv2.Value{Sum: &apiv2.Value_Party{Party: ownerParty}}},
									{Label: "amount", Value: &apiv2.Value{Sum: &apiv2.Value_Numeric{Numeric: amount}}},
									{Label: "openRound", Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: openMiningRoundCID}}},
								}}}},
							},
						},
					},
				},
				ActAs:              []string{ownerParty},
				DisclosedContracts: disclosures,
			},
		})
		return submitErr
	}

	err = submitTap(disclosedContracts)
	if err != nil && (strings.Contains(err.Error(), "CONTRACT_NOT_FOUND") || strings.Contains(err.Error(), "LOCAL_VERDICT_INACTIVE_CONTRACTS")) {
		err = submitTap(nil)
	}
	if err != nil && strings.Contains(err.Error(), "CONTRACT_NOT_FOUND") {
		if refreshedCID, refreshErr := getActiveOpenMiningRoundCID(); refreshErr == nil {
			openMiningRoundCID = refreshedCID
			err = submitTap(nil)
		}
	}
	if err != nil {
		return fmt.Errorf("mint AMT for %s: %w", ownerParty, err)
	}
	return nil
}

func getTransferFactoryFromScanProxy(
	ctx context.Context,
	participant canton.Participant,
	expectedAdmin string,
	sender string,
	receiver string,
	amount string,
	instrumentAdmin string,
	instrumentID string,
	inputHoldingCids []string,
) (types.CONTRACT_ID, []*apiv2.DisclosedContract, splice_api_token_metadata_v1.ChoiceContext, error) {
	requestEditor := func(reqCtx context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("retrieve participant token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
		return nil
	}
	transferClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("create transfer instruction client: %w", err)
	}

	resp, err := transferClient.GetTransferFactoryWithResponse(ctx, transferInstructionV1.GetFactoryRequest{
		ChoiceArguments: map[string]any{
			"expectedAdmin": expectedAdmin,
			"transfer": map[string]any{
				"sender":   sender,
				"receiver": receiver,
				"amount":   amount,
				"instrumentId": map[string]any{
					"admin": instrumentAdmin,
					"id":    instrumentID,
				},
				"lock":             nil,
				"requestedAt":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"executeBefore":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				"inputHoldingCids": inputHoldingCids,
				"meta":             map[string]any{"values": map[string]any{}},
			},
			"extraArgs": map[string]any{
				"context": map[string]any{"values": map[string]any{}},
				"meta":    map[string]any{"values": map[string]any{}},
			},
		},
	})
	if err != nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("get transfer factory: %w", err)
	}
	if resp.StatusCode() != http.StatusOK || resp.JSON200 == nil {
		return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("unexpected transfer factory response status=%d", resp.StatusCode())
	}

	disclosures := make([]*apiv2.DisclosedContract, 0, len(resp.JSON200.ChoiceContext.DisclosedContracts))
	for _, contract := range resp.JSON200.ChoiceContext.DisclosedContracts {
		id, parseErr := TemplateIdFromString(contract.TemplateId)
		if parseErr != nil {
			return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("parse transfer factory disclosure template id: %w", parseErr)
		}
		createdEventBlob, decodeErr := base64.StdEncoding.DecodeString(contract.CreatedEventBlob)
		if decodeErr != nil {
			return "", nil, splice_api_token_metadata_v1.ChoiceContext{}, fmt.Errorf("decode transfer factory disclosure created event blob: %w", decodeErr)
		}
		disclosures = append(disclosures, &apiv2.DisclosedContract{
			TemplateId:       id,
			ContractId:       contract.ContractId,
			CreatedEventBlob: createdEventBlob,
			SynchronizerId:   contract.SynchronizerId,
		})
	}

	choiceCtx := splice_api_token_metadata_v1.ChoiceContext{Values: types.TEXTMAP{}}
	if rawValues, ok := resp.JSON200.ChoiceContext.ChoiceContextData["values"].(map[string]any); ok {
		for k, rawValue := range rawValues {
			vMap, ok := rawValue.(map[string]any)
			if !ok {
				continue
			}
			tag, _ := vMap["tag"].(string)
			value := vMap["value"]
			av := splice_api_token_metadata_v1.AnyValue{}
			switch tag {
			case "AV_ContractId":
				if s, ok := value.(string); ok {
					cid := types.CONTRACT_ID(s)
					av.AVContractId = &cid
				}
			case "AV_Text":
				if s, ok := value.(string); ok {
					t := types.TEXT(s)
					av.AVText = &t
				}
			case "AV_Bool":
				if b, ok := value.(bool); ok {
					t := types.BOOL(b)
					av.AVBool = &t
				}
			}
			if av.GetVariantTag() != "" {
				choiceCtx.Values[k] = av
			}
		}
	}

	return types.CONTRACT_ID(resp.JSON200.FactoryId), disclosures, choiceCtx, nil
}

func selectUnlockedHoldingCIDs(holdings []*apiv2.ActiveContract, owner, admin, instrumentID string) ([]types.CONTRACT_ID, []*apiv2.DisclosedContract) {
	cids := make([]types.CONTRACT_ID, 0, len(holdings))
	disclosures := make([]*apiv2.DisclosedContract, 0, len(holdings))
	for _, holding := range holdings {
		views := holding.GetCreatedEvent().GetInterfaceViews()
		if len(views) == 0 || views[0].GetViewValue() == nil {
			continue
		}
		fields := views[0].GetViewValue().GetFields()
		if len(fields) < 4 {
			continue
		}
		ownerParty := fields[0].GetValue().GetParty()
		if ownerParty != owner {
			continue
		}
		instrumentRecord := fields[1].GetValue().GetRecord()
		if instrumentRecord == nil || len(instrumentRecord.GetFields()) < 2 {
			continue
		}
		var holdingAdmin, holdingID string
		for _, f := range instrumentRecord.GetFields() {
			switch f.GetLabel() {
			case "admin":
				holdingAdmin = f.GetValue().GetParty()
			case "id":
				holdingID = f.GetValue().GetText()
			}
		}
		if holdingAdmin != admin || holdingID != instrumentID {
			continue
		}
		isLocked := fields[3].GetValue().GetOptional().GetValue() != nil
		if isLocked {
			continue
		}
		cid := types.CONTRACT_ID(holding.GetCreatedEvent().GetContractId())
		cids = append(cids, cid)
		disclosures = append(disclosures, &apiv2.DisclosedContract{
			TemplateId:       holding.GetCreatedEvent().GetTemplateId(),
			ContractId:       holding.GetCreatedEvent().GetContractId(),
			CreatedEventBlob: holding.GetCreatedEvent().GetCreatedEventBlob(),
			SynchronizerId:   holding.GetSynchronizerId(),
		})
	}
	return cids, disclosures
}
