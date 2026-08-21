package ccip

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	ccipcore "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/ledger"
	"github.com/smartcontractkit/chainlink-canton/registry-kit/registry"
)

// RegisterTokenPoolClientInput is the input for TAR registration via ledger client with per-step actAs parties.
type RegisterTokenPoolClientInput struct {
	TokenAdminRegistryCID string
	InstrumentId          splice_api_token_holding_v1.InstrumentId
	PoolInstanceID        string
	CcipParty             string
	PoolOwnerParty        string
	// PoolOwnerClient optionally supplies the ledger client for pool-owner steps when the
	// pool owner is hosted on a different CTF participant than the CCIP party.
	PoolOwnerClient ledger.Client
}

// RegisterTokenPoolViaClient runs ProposeAdministrator → AcceptAdminRole → SetPool using explicit actAs parties.
// Use when PoolOwnerParty differs from the CTF participant party (CLDF deploy ops always act as participant.PartyID).
func RegisterTokenPoolViaClient(ctx context.Context, client ledger.Client, input RegisterTokenPoolClientInput) (tokenConfigCID string, tarCID string, err error) {
	poolOwnerClient := input.PoolOwnerClient
	if poolOwnerClient == nil {
		poolOwnerClient = client
	}

	tokenConfigAddr := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(input.InstrumentId).Bytes())).
		RawInstanceAddress(types.PARTY(input.CcipParty)).
		InstanceAddress()

	tarCID = input.TokenAdminRegistryCID
	tokenConfigCID, found, err := findTokenConfigCID(ctx, client, input.CcipParty, tokenConfigAddr)
	if err != nil {
		return "", tarCID, fmt.Errorf("lookup token config: %w", err)
	}
	var tokenConfigCIDArg *types.CONTRACT_ID
	if found {
		cid := types.CONTRACT_ID(tokenConfigCID)
		tokenConfigCIDArg = &cid
	}

	skipAccept := false
	proposeRes, err := client.SubmitExercise(ctx, input.CcipParty, ccipcore.TokenAdminRegistry{}, tarCID, "ProposeAdministrator",
		ccipcore.ProposeAdministrator{
			TokenConfigCid: tokenConfigCIDArg,
			InstrumentId:   input.InstrumentId,
			NewAdmin:       types.PARTY(input.PoolOwnerParty),
			Caller:         types.PARTY(input.CcipParty),
		})
	if err != nil {
		if strings.Contains(err.Error(), "admin already set") {
			skipAccept = true
		} else {
			return "", tarCID, fmt.Errorf("propose administrator: %w", err)
		}
	} else {
		if newTarCID, ok := ledger.CreatedContractID(proposeRes.GetTransaction(), "TokenAdminRegistry"); ok {
			tarCID = newTarCID
		}
		cid, ok := ledger.CreatedContractID(proposeRes.GetTransaction(), "TokenConfig")
		if !ok {
			return "", tarCID, fmt.Errorf("TokenConfig not created after ProposeAdministrator")
		}
		tokenConfigCID = cid
	}

	if !skipAccept {
		tarDisclosed, err := registry.DiscloseByID(ctx, client, input.CcipParty, tarCID)
		if err != nil {
			return "", tarCID, fmt.Errorf("disclose TAR for AcceptAdminRole: %w", err)
		}

		acceptRes, err := poolOwnerClient.SubmitExerciseMulti(ctx, []string{input.PoolOwnerParty}, ccipcore.TokenAdminRegistry{}, tarCID, "AcceptAdminRole",
			ccipcore.AcceptAdminRole{
				TokenConfigCid: types.CONTRACT_ID(tokenConfigCID),
				InstrumentId:   input.InstrumentId,
				Caller:         types.PARTY(input.PoolOwnerParty),
			}, []*apiv2.DisclosedContract{tarDisclosed})
		if err != nil {
			return "", tarCID, fmt.Errorf("accept admin role: %w", err)
		}
		cid, ok := ledger.CreatedContractID(acceptRes.GetTransaction(), "TokenConfig")
		if !ok {
			return "", tarCID, fmt.Errorf("TokenConfig not created after AcceptAdminRole")
		}
		tokenConfigCID = cid
	}

	tarDisclosed, err := registry.DiscloseByID(ctx, client, input.CcipParty, tarCID)
	if err != nil {
		return "", tarCID, fmt.Errorf("disclose TAR for SetPool: %w", err)
	}

	setPoolRes, err := poolOwnerClient.SubmitExerciseMulti(ctx, []string{input.PoolOwnerParty}, ccipcore.TokenAdminRegistry{}, tarCID, "SetPool",
		ccipcore.SetPool{
			TokenConfigCid: types.CONTRACT_ID(tokenConfigCID),
			InstrumentId:   input.InstrumentId,
			TokenPool: &ccipcore.PoolRegistration2{
				PoolOwner:      types.PARTY(input.PoolOwnerParty),
				PoolInstanceId: types.TEXT(input.PoolInstanceID),
			},
			Caller: types.PARTY(input.PoolOwnerParty),
		}, []*apiv2.DisclosedContract{tarDisclosed})
	if err != nil {
		return "", tarCID, fmt.Errorf("set pool: %w", err)
	}
	cid, ok := ledger.CreatedContractID(setPoolRes.GetTransaction(), "TokenConfig")
	if !ok {
		return "", tarCID, fmt.Errorf("TokenConfig not created after SetPool")
	}

	return cid, tarCID, nil
}

func findTokenConfigCID(ctx context.Context, client ledger.Client, ccipParty string, addr contracts.InstanceAddress) (string, bool, error) {
	participant := client.Participant()
	cid, err := contract.FindActiveContractIDByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		[]string{ccipParty},
		ccipcore.TokenConfig{}.GetTemplateID(),
		addr,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no active contract found") {
			return "", false, nil
		}

		return "", false, err
	}

	return cid, true, nil
}
