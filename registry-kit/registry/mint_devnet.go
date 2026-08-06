package registry

import (
	"context"
	"fmt"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/google/uuid"
	"github.com/smartcontractkit/go-daml/pkg/types"

	splice_api_token_holding_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	registryapp "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/utility/registry_app_v0"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/operator"
)

// MintDevnetInput identifies the registrar instrument and factories for devnet mint.
type MintDevnetInput struct {
	RegistrarParty       string
	InstrumentID         string
	AllocationFactoryCID string
	Holder               string
	Amount               string
}

// RequestMintViaOperatorBackend runs AllocationFactory_RequestMint with DA operator-backend context.
func RequestMintViaOperatorBackend(
	ctx context.Context,
	client ledger.Client,
	backend *operator.Client,
	in MintDevnetInput,
) (string, error) {
	instrument := splice_api_token_holding_v1.InstrumentId{
		Admin: types.PARTY(in.RegistrarParty),
		Id:    types.TEXT(in.InstrumentID),
	}
	holder := in.Holder
	if holder == "" {
		holder = in.RegistrarParty
	}

	choiceCtx, err := backend.MintRequestContext(ctx, holder, instrument)
	if err != nil {
		return "", fmt.Errorf("fetch mint request context: %w", err)
	}

	now := time.Now()
	args := registryapp.AllocationFactoryRequestMint{
		ExpectedAdmin: types.PARTY(in.RegistrarParty),
		Mint: registryapp.Mint{
			InstrumentId:  instrument,
			Amount:        types.NUMERIC(in.Amount),
			Holder:        types.PARTY(holder),
			Reference:     types.TEXT(fmt.Sprintf("mint-%s", uuid.NewString()[:8])),
			RequestedAt:   types.TIMESTAMP(now),
			ExecuteBefore: types.TIMESTAMP(now.Add(24 * time.Hour)),
			Meta:          emptyMetadata(),
		},
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: choiceCtx.Context,
			Meta:    emptyMetadata(),
		},
	}

	res, err := client.SubmitExerciseMulti(ctx, []string{in.RegistrarParty}, registryapp.AllocationFactory{}, in.AllocationFactoryCID, "AllocationFactory_RequestMint", args, choiceCtx.DisclosedContracts)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "MintRequest")
	if !ok {
		return "", fmt.Errorf("MintRequest not created")
	}

	return cid, nil
}

// AcceptMintViaOperatorBackend runs MintRequest_Accept with DA operator-backend context.
func AcceptMintViaOperatorBackend(
	ctx context.Context,
	client ledger.Client,
	backend *operator.Client,
	registrarParty, mintRequestCID string,
) (string, error) {
	choiceCtx, err := backend.MintAcceptContext(ctx, mintRequestCID)
	if err != nil {
		return "", fmt.Errorf("fetch mint accept context: %w", err)
	}

	args := registryapp.MintRequestAccept{
		ExtraArgs: splice_api_token_metadata_v1.ExtraArgs{
			Context: choiceCtx.Context,
			Meta:    emptyMetadata(),
		},
	}

	var disclosed []*apiv2.DisclosedContract
	if len(choiceCtx.DisclosedContracts) > 0 {
		disclosed = choiceCtx.DisclosedContracts
	} else {
		instDisclosed, derr := DiscloseByID(ctx, client, registrarParty, mintRequestCID)
		if derr == nil {
			disclosed = []*apiv2.DisclosedContract{instDisclosed}
		}
	}

	res, err := client.SubmitExerciseMulti(ctx, []string{registrarParty}, registryapp.MintRequest{}, mintRequestCID, "MintRequest_Accept", args, disclosed)
	if err != nil {
		return "", err
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "Holding")
	if !ok {
		return "", fmt.Errorf("registry Holding not created")
	}

	return cid, nil
}

// FindMintRequestForInstrument returns the first pending MintRequest CID for an instrument.
func FindMintRequestForInstrument(ctx context.Context, client ledger.Client, registrarParty, instrumentID string) (string, error) {
	refs, err := FindContractsByEntity(ctx, client, registrarParty, registryapp.MintRequest{}, "MintRequest")
	if err != nil {
		return "", err
	}
	if len(refs) == 0 {
		return "", nil
	}
	// When multiple exist, return the most recently listed (ACS sort is by created time in testhelpers).
	return refs[len(refs)-1].ContractID, nil
}
