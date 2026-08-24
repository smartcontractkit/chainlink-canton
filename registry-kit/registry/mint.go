package registry

import (
	"context"
	"fmt"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/google/uuid"

	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_metadata_v1"
	registryapp "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/utility/registry_app_v0"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/testhelpers"
)

// MintViaAllocationFactory runs AllocationFactory_RequestMint + MintRequest_Accept
// with self-assembled choice context (no DA operator backend).
func MintViaAllocationFactory(ctx context.Context, client ledger.Client, bootstrap BootstrapResult, owner, amount string) (string, error) {
	instDisclosed, err := DiscloseByID(ctx, client, bootstrap.Party, bootstrap.InstrumentConfiguration)
	if err != nil {
		return "", fmt.Errorf("disclose instrument configuration: %w", err)
	}

	mintCtx := MintChoiceContext(bootstrap.InstrumentConfiguration, false)
	mintReqCID, err := requestMint(ctx, client, bootstrap, owner, amount, mintCtx, instDisclosed)
	if err != nil {
		return "", err
	}

	acceptCtx := MintChoiceContext(bootstrap.InstrumentConfiguration, true)

	return acceptMintRequest(ctx, client, bootstrap.Party, owner, mintReqCID, acceptCtx, instDisclosed)
}

func requestMint(
	ctx context.Context,
	client ledger.Client,
	bootstrap BootstrapResult,
	owner, amount string,
	ctxData splice_api_token_metadata_v1.ChoiceContext,
	disclosed *apiv2.DisclosedContract,
) (string, error) {
	now := time.Now()
	args := registryapp.AllocationFactoryRequestMint{
		ExpectedAdmin: types.PARTY(bootstrap.Party),
		Mint: registryapp.Mint{
			InstrumentId: splice_api_token_holding_v1.InstrumentId{
				Admin: types.PARTY(bootstrap.Party),
				Id:    types.TEXT(bootstrap.InstrumentID),
			},
			Amount:        types.NUMERIC(amount),
			Holder:        types.PARTY(owner),
			Reference:     types.TEXT(fmt.Sprintf("mint-%s", uuid.NewString()[:8])),
			RequestedAt:   types.TIMESTAMP(now),
			ExecuteBefore: types.TIMESTAMP(now.Add(24 * time.Hour)),
			Meta:          emptyMetadata(),
		},
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: ctxData,
			Meta:    emptyMetadata(),
		},
	}
	res, err := client.SubmitExerciseMulti(ctx, mintActAs(bootstrap.Party, owner), registryapp.AllocationFactory{}, bootstrap.AllocationFactory, "AllocationFactory_RequestMint", args, []*apiv2.DisclosedContract{disclosed})
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "MintRequest")
	if !ok {
		return "", fmt.Errorf("MintRequest not created")
	}

	return cid, nil
}

func acceptMintRequest(
	ctx context.Context,
	client ledger.Client,
	registrar, party, mintReqCID string,
	ctxData splice_api_token_metadata_v1.ChoiceContext,
	disclosed *apiv2.DisclosedContract,
) (string, error) {
	args := registryapp.MintRequestAccept{
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: ctxData,
			Meta:    emptyMetadata(),
		},
	}
	res, err := client.SubmitExerciseMulti(ctx, mintActAs(registrar, party), registryapp.MintRequest{}, mintReqCID, "MintRequest_Accept", args, []*apiv2.DisclosedContract{disclosed})
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "Holding")
	if !ok {
		return "", fmt.Errorf("registry Holding not created")
	}

	return cid, nil
}

func mintActAs(registrar, owner string) []string {
	actAs := []string{registrar}
	if owner != "" && owner != registrar {
		actAs = append(actAs, owner)
	}

	return actAs
}

// DiscloseByID fetches a disclosed contract blob for JSON API submission.
func DiscloseByID(ctx context.Context, client ledger.Client, party, contractID string) (*apiv2.DisclosedContract, error) {
	return testhelpers.GetDisclosedContractById(ctx, client.ForParty(party), contractID)
}
