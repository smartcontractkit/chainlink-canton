package ccip

import (
	"context"
	"fmt"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	ccipcore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
)

// SetBurnMintFactoryInput identifies contracts for pointing TokenConfig at a burn-mint factory.
type SetBurnMintFactoryInput struct {
	TokenAdminRegistryCID string
	TokenConfigCID        string
	InstrumentId          splice_api_token_holding_v1.InstrumentId
	BurnMintFactoryCID    string
	CcipParty             string
	PoolOwnerParty        string
	CcipClient            ledger.Client
	PoolOwnerClient       ledger.Client
}

// SetBurnMintFactory exercises TokenAdminRegistry.SetBurnMintFactory and returns the updated TokenConfig CID.
func SetBurnMintFactory(ctx context.Context, client ledger.Client, input SetBurnMintFactoryInput) (string, error) {
	ccipClient := input.CcipClient
	if ccipClient == nil {
		ccipClient = client
	}
	poolOwnerClient := input.PoolOwnerClient
	if poolOwnerClient == nil {
		poolOwnerClient = client
	}

	tarDisclosed, err := registry.DiscloseByID(ctx, ccipClient, input.CcipParty, input.TokenAdminRegistryCID)
	if err != nil {
		return "", err
	}

	factoryCID := types.CONTRACT_ID(input.BurnMintFactoryCID)
	res, err := poolOwnerClient.SubmitExerciseMulti(ctx, []string{input.PoolOwnerParty}, ccipcore.TokenAdminRegistry{}, input.TokenAdminRegistryCID, "SetBurnMintFactory",
		ccipcore.SetBurnMintFactory{
			TokenConfigCid:  types.CONTRACT_ID(input.TokenConfigCID),
			InstrumentId:    input.InstrumentId,
			BurnMintFactory: &factoryCID,
			Caller:          types.PARTY(input.PoolOwnerParty),
		}, []*apiv2.DisclosedContract{tarDisclosed})
	if err != nil {
		return "", fmt.Errorf("set burn mint factory: %w", err)
	}
	cid, ok := ledger.CreatedContractID(res.GetTransaction(), "TokenConfig")
	if !ok {
		return "", fmt.Errorf("TokenConfig not created after SetBurnMintFactory")
	}

	return cid, nil
}
