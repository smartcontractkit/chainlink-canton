package tests

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	apiv2 "github.com/smartcontractkit/chainlink-canton-internal/pb/gen/com/daml/ledger/api/v2"
)

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
	return a.getContract(ctx, &apiv2.Identifier{
		PackageId:  "#ccip-feequoter",
		ModuleName: "CCIP.FeeQuoter",
		EntityName: "FeeQuoter",
	})
}

func (a *CCIPApi) GetOnRamp(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return a.getContract(ctx, &apiv2.Identifier{
		PackageId:  "#ccip-onramp",
		ModuleName: "CCIP.OnRamp",
		EntityName: "OnRamp",
	})
}

func (a *CCIPApi) GetOffRamp(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return a.getContract(ctx, &apiv2.Identifier{
		PackageId:  "#ccip-offramp",
		ModuleName: "CCIP.OffRamp",
		EntityName: "OffRamp",
	})
}

func (a *CCIPApi) GetRouter(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return a.getContract(ctx, &apiv2.Identifier{
		PackageId:  "#ccip-router",
		ModuleName: "CCIP.Router",
		EntityName: "Router",
	})
}

func (a *CCIPApi) GetCCV(ctx context.Context) (*apiv2.DisclosedContract, error) {
	return a.getContract(ctx, &apiv2.Identifier{
		PackageId:  "#ccip-committeeverifier",
		ModuleName: "CCIP.CommitteeVerifier",
		EntityName: "CommitteeVerifier",
	})
}
