package tests

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

func getDisclosedContract(ctx context.Context, participant *Participant, party string, templateId *apiv2.Identifier) (*apiv2.DisclosedContract, error) {
	jwToken, err := getJWT()
	if err != nil {
		return nil, fmt.Errorf("could not get JWToken: %w", err)
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx = metadata.NewOutgoingContext(ctx, md)

	activeContracts, err := GetActiveContractsForPartyTemplateId(ctx, participant, party, templateId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active contracts with identifier %v found for party %s", templateId, party)
	}

	contract := activeContracts[len(activeContracts)-1]
	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}

func getDisclosedInterfaceContracts(ctx context.Context, participant *Participant, party string, interfaceId *apiv2.Identifier) ([]*apiv2.DisclosedContract, error) {
	jwToken, err := getJWT()
	if err != nil {
		return nil, fmt.Errorf("could not get JWToken: %w", err)
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx = metadata.NewOutgoingContext(ctx, md)

	activeContracts, err := GetActiveContractsForPartyInterface(ctx, participant, party, interfaceId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active interface contracts with identifier %v found for party %s", interfaceId, party)
	}

	var disclosedContracts []*apiv2.DisclosedContract
	for _, contract := range activeContracts {
		disclosedContracts = append(disclosedContracts, &apiv2.DisclosedContract{
			TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
			ContractId:       contract.GetCreatedEvent().GetContractId(),
			CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
			SynchronizerId:   contract.GetSynchronizerId(),
		})
	}
	return disclosedContracts, nil
}

type CCIPApi struct {
	Participant *Participant
	CCIPParty   string
}

func (a *CCIPApi) GetCCIPParty() string {
	return a.CCIPParty
}

func (a *CCIPApi) getContract(ctx context.Context, templateId *apiv2.Identifier) (*apiv2.DisclosedContract, error) {
	jwToken, err := getJWT()
	if err != nil {
		return nil, fmt.Errorf("could not get JWToken: %w", err)
	}
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", jwToken))
	ctx = metadata.NewOutgoingContext(ctx, md)

	activeContracts, err := GetActiveContractsForPartyTemplateId(ctx, a.Participant, a.CCIPParty, templateId)
	if err != nil {
		return nil, fmt.Errorf("could not get active contracts: %w", err)
	}
	if len(activeContracts) == 0 {
		return nil, fmt.Errorf("no active OnRamp contracts found for party %s", a.CCIPParty)
	}

	contract := activeContracts[len(activeContracts)-1]
	return &apiv2.DisclosedContract{
		TemplateId:       contract.GetCreatedEvent().GetTemplateId(),
		ContractId:       contract.GetCreatedEvent().GetContractId(),
		CreatedEventBlob: contract.GetCreatedEvent().GetCreatedEventBlob(),
		SynchronizerId:   contract.GetSynchronizerId(),
	}, nil
}

func (a *CCIPApi) GetFeeQuoter(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-feequoter",
		ModuleName: "CCIP.FeeQuoter",
		EntityName: "FeeQuoter",
	})
}

func (a *CCIPApi) GetOnRamp(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-onramp",
		ModuleName: "CCIP.OnRamp",
		EntityName: "OnRamp",
	})
}

func (a *CCIPApi) GetOffRamp(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-offramp",
		ModuleName: "CCIP.OffRamp",
		EntityName: "OffRamp",
	})
}

func (a *CCIPApi) GetRouter(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-router",
		ModuleName: "CCIP.Router",
		EntityName: "Router",
	})
}

func (a *CCIPApi) GetCCV(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-committeeverifier",
		ModuleName: "CCIP.CommitteeVerifier",
		EntityName: "CommitteeVerifier",
	})
}

func (a *CCIPApi) GetTokenAdminRegistry(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.CCIPParty, &apiv2.Identifier{
		PackageId:  "#ccip-tokenadminregistry",
		ModuleName: "CCIP.TokenAdminRegistry",
		EntityName: "TokenAdminRegistry",
	})
}

type TokenPoolApi struct {
	Participant *Participant
	AdminParty  string
}

func (a *TokenPoolApi) GetTokenPool(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.AdminParty, &apiv2.Identifier{
		PackageId:  "#ccip-lockreleasetokenpool",
		ModuleName: "CCIP.LockReleaseTokenPool",
		EntityName: "LockReleaseTokenPool",
	})
}

func (a *TokenPoolApi) GetCCV(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return getDisclosedContract(ctx, a.Participant, a.AdminParty, &apiv2.Identifier{
		PackageId:  "#ccip-committeeverifier",
		ModuleName: "CCIP.CommitteeVerifier",
		EntityName: "CommitteeVerifier",
	})
}

func (a *TokenPoolApi) GetHoldings(ctx context.Context) ([]*apiv2.DisclosedContract, error) {
	return getDisclosedInterfaceContracts(ctx, a.Participant, a.AdminParty, &apiv2.Identifier{
		PackageId:  "#splice-api-token-holding-v1",
		ModuleName: "Splice.Api.Token.HoldingV1",
		EntityName: "Holding",
	})
}
