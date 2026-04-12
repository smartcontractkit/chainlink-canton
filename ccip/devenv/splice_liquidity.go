package devenv

import (
	"context"
	"fmt"
	"net/http"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/openapi/gen/scanProxy"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/tokenMetadataV1"
	"github.com/smartcontractkit/chainlink-canton/openapi/gen/transferInstructionV1"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

const instrumentFieldAdmin = "admin"

func validatorRequestEditor(participant canton.Participant) func(context.Context, *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		token, err := participant.TokenSource.Token()
		if err != nil {
			return fmt.Errorf("retrieve participant token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

		return nil
	}
}

func newScanProxyClients(
	participant canton.Participant,
) (scanProxy.ClientWithResponsesInterface, tokenMetadataV1.ClientWithResponsesInterface, transferInstructionV1.ClientWithResponsesInterface, error) {
	requestEditor := validatorRequestEditor(participant)
	scanClient, err := scanProxy.NewClientWithResponses(participant.Endpoints.ValidatorAPIURL, scanProxy.WithRequestEditorFn(requestEditor))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create scan proxy client: %w", err)
	}
	metadataClient, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create token metadata client: %w", err)
	}
	transferClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create transfer instruction client: %w", err)
	}

	return scanClient, metadataClient, transferClient, nil
}

func newTokenMetadataClient(participant canton.Participant) (tokenMetadataV1.ClientWithResponsesInterface, error) {
	metadataClient, err := tokenMetadataV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		tokenMetadataV1.WithRequestEditorFn(validatorRequestEditor(participant)),
	)
	if err != nil {
		return nil, fmt.Errorf("create token metadata client: %w", err)
	}

	return metadataClient, nil
}

func newTransferInstructionClient(participant canton.Participant) (transferInstructionV1.ClientWithResponsesInterface, error) {
	transferClient, err := transferInstructionV1.NewClientWithResponses(
		fmt.Sprintf("%s/v0/scan-proxy", participant.Endpoints.ValidatorAPIURL),
		transferInstructionV1.WithRequestEditorFn(validatorRequestEditor(participant)),
	)
	if err != nil {
		return nil, fmt.Errorf("create transfer instruction client: %w", err)
	}

	return transferClient, nil
}

func seedAMTLiquidity(ctx context.Context, participant canton.Participant, ownerParty string, amount string) error {
	scanClient, metadataClient, transferClient, err := newScanProxyClients(participant)
	if err != nil {
		return err
	}
	if _, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount); err != nil {
		return fmt.Errorf("mint AMT for %s: %w", ownerParty, err)
	}

	return nil
}

func mintTwoAmuletHoldings(
	ctx context.Context,
	participant canton.Participant,
	ownerParty string,
	amount string,
) (types.CONTRACT_ID, types.CONTRACT_ID, []*apiv2.DisclosedContract, error) {
	scanClient, metadataClient, transferClient, err := newScanProxyClients(participant)
	if err != nil {
		return "", "", nil, err
	}
	feeHoldingCID, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount)
	if err != nil {
		return "", "", nil, fmt.Errorf("mint fee holding: %w", err)
	}
	tokenHoldingCID, err := testhelpers.MintAMT(ctx, participant, metadataClient, transferClient, scanClient, ownerParty, amount)
	if err != nil {
		return "", "", nil, fmt.Errorf("mint token-transfer holding: %w", err)
	}
	if feeHoldingCID == tokenHoldingCID {
		return "", "", nil, fmt.Errorf("mint returned same holding cid twice: %s", feeHoldingCID)
	}
	disclosedFeeHolding, err := testhelpers.GetDisclosedContractById(ctx, participant, feeHoldingCID)
	if err != nil {
		return "", "", nil, fmt.Errorf("get disclosed fee holding by id: %w", err)
	}
	disclosedTokenHolding, err := testhelpers.GetDisclosedContractById(ctx, participant, tokenHoldingCID)
	if err != nil {
		return "", "", nil, fmt.Errorf("get disclosed token holding by id: %w", err)
	}

	return types.CONTRACT_ID(feeHoldingCID), types.CONTRACT_ID(tokenHoldingCID), []*apiv2.DisclosedContract{
		disclosedFeeHolding,
		disclosedTokenHolding,
	}, nil
}

// resolveRegistryAdmin returns the token-registry admin party ID (DSO admin)
// from scan-proxy registry info for the participant's validator.
func resolveRegistryAdmin(ctx context.Context, participant canton.Participant) (string, error) {
	metadataClient, err := newTokenMetadataClient(participant)
	if err != nil {
		return "", err
	}
	registryAdmin, err := testhelpers.GetRegistryAdmin(ctx, metadataClient)
	if err != nil {
		return "", fmt.Errorf("get registry admin: %w", err)
	}

	return registryAdmin, nil
}

func selectUnlockedHoldingCIDs(holdings []*apiv2.ActiveContract, owner, admin, instrumentID string) ([]types.CONTRACT_ID, []*apiv2.DisclosedContract) {
	cids := make([]types.CONTRACT_ID, 0, len(holdings))
	disclosures := make([]*apiv2.DisclosedContract, 0, len(holdings))
	seen := make(map[string]struct{}, len(holdings))
	for _, holding := range holdings {
		created := holding.GetCreatedEvent()
		fields := created.GetInterfaceViews()[0].GetViewValue().GetFields()
		if fields[0].GetValue().GetParty() != owner {
			continue
		}

		var holdingAdmin, holdingID string
		for _, f := range fields[1].GetValue().GetRecord().GetFields() {
			switch f.GetLabel() {
			case instrumentFieldAdmin:
				holdingAdmin = f.GetValue().GetParty()
			case "id":
				holdingID = f.GetValue().GetText()
			}
		}
		if holdingAdmin != admin || holdingID != instrumentID {
			continue
		}
		if fields[3].GetValue().GetOptional().GetValue() != nil {
			continue
		}

		contractID := created.GetContractId()
		if _, exists := seen[contractID]; exists {
			continue
		}
		seen[contractID] = struct{}{}
		cid := types.CONTRACT_ID(contractID)
		cids = append(cids, cid)
		disclosures = append(disclosures, &apiv2.DisclosedContract{
			TemplateId:       created.GetTemplateId(),
			ContractId:       contractID,
			CreatedEventBlob: created.GetCreatedEventBlob(),
			SynchronizerId:   holding.GetSynchronizerId(),
		})
	}

	return cids, disclosures
}
