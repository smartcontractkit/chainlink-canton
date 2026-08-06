package ccip

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/go-daml/pkg/types"

	burnminttokenpool "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/burnminttokenpool"
	splice_api_token_metadata_v1 "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_metadata_v1"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
)

// ReleaseFromTicketInput identifies contracts for a pool release exercise.
type ReleaseFromTicketInput struct {
	PoolCID               string
	TokenAdminRegistryCID string
	TokenConfigCID        string
	RMNRemoteCID          string
	TokenReceiveTicketCID string
	Receiver              string
	PoolOwner             string
	ExtraContext          splice_api_token_metadata_v1.ChoiceContext
	Disclosures           PoolReleaseDisclosures
}

// ReleaseFromTicket exercises BurnMintTokenPool.ReleaseFromTicket and returns the created Holding CID.
func ReleaseFromTicket(ctx context.Context, client ledger.Client, input ReleaseFromTicketInput) (string, error) {
	args := burnminttokenpool.ReleaseFromTicket{
		TokenAdminRegistryCid: types.CONTRACT_ID(input.TokenAdminRegistryCID),
		TokenConfigCid:        types.CONTRACT_ID(input.TokenConfigCID),
		RmnRemoteCid:          types.CONTRACT_ID(input.RMNRemoteCID),
		Context:               input.ExtraContext,
		TokenReceiveTicketCid: types.CONTRACT_ID(input.TokenReceiveTicketCID),
		Caller:                types.PARTY(input.Receiver),
	}

	actAs := []string{input.Receiver}
	if input.PoolOwner != "" && input.PoolOwner != input.Receiver {
		actAs = append(actAs, input.PoolOwner)
	}

	res, err := client.SubmitExerciseMulti(ctx, actAs, burnminttokenpool.BurnMintTokenPool{}, input.PoolCID, "ReleaseFromTicket", args, input.Disclosures.All())
	if err != nil {
		return "", fmt.Errorf("release from ticket: %w", err)
	}

	holdingCID, ok := ledger.CreatedHoldingForOwner(res.GetTransaction(), input.Receiver)
	if !ok {
		return "", fmt.Errorf("registry Holding for receiver not created")
	}

	return holdingCID, nil
}
