package sequences

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// RegisterTokenPoolInput is the input for registering a token pool with the TokenAdminRegistry.
type RegisterTokenPoolInput struct {
	// TokenAdminRegistryInstanceAddress is the instance address of the TokenAdminRegistry contract.
	TokenAdminRegistryInstanceAddress contracts.InstanceAddress
	// InstrumentId identifies the token (admin party + token id).
	InstrumentId splice_api_token_holding_v1.InstrumentId
	// PoolInstanceID is the instance ID of the token pool.
	PoolInstanceID string
	// CcipParty is the CCIP owner party (acts on ProposeAdministrator).
	CcipParty string
	// PoolOwnerParty is the token pool owner party (acts on AcceptAdminRole and SetPool).
	PoolOwnerParty string
}

var RegisterTokenPool = operations.NewSequence(
	"canton/ccip/register_token_pool",
	semver.MustParse("1.0.0"),
	"Registers a token pool with the TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole, SetPool)",
	registerTokenPool,
)

func registerTokenPool(b operations.Bundle, deps canton.Chain, input RegisterTokenPoolInput) (sequences.OnChainOutput, error) {
	instrumentId := input.InstrumentId
	ccipParty := input.CcipParty
	poolOwnerParty := input.PoolOwnerParty
	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).RawInstanceAddress(types.PARTY(ccipParty)).InstanceAddress()

	existingTokenConfigCid, tokenConfigFound, err := findTokenConfigCid(b, deps, tokenConfigAddress)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config: %w", err)
	}
	var tokenConfigCidArg *types.CONTRACT_ID
	if tokenConfigFound {
		tokenConfigCidArg = new(existingTokenConfigCid)
	}

	// Step 1: ProposeAdministrator (CCIP acts)
	skipAcceptAdminRole := false
	_, err = operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[tokenadminregistry.ProposeAdministrator]{
		InstanceAddress: input.TokenAdminRegistryInstanceAddress,
		Args: tokenadminregistry.ProposeAdministrator{
			TokenConfigCid: tokenConfigCidArg,
			InstrumentId:   instrumentId,
			NewAdmin:       types.PARTY(poolOwnerParty),
			Caller:         types.PARTY(ccipParty),
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "admin already set") {
			skipAcceptAdminRole = true
		} else {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to propose administrator: %w", err)
		}
	}

	tokenConfigCid, tokenConfigFound, err := findTokenConfigCid(b, deps, tokenConfigAddress)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after propose: %w", err)
	}
	if !tokenConfigFound {
		return sequences.OnChainOutput{}, fmt.Errorf("token config not found after propose")
	}

	// Step 2: AcceptAdminRole (pool owner acts). Exercise resolves current TAR contract by InstanceAddress.
	if !skipAcceptAdminRole {
		_, err = operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[tokenadminregistry.AcceptAdminRole]{
			InstanceAddress: input.TokenAdminRegistryInstanceAddress,
			Args: tokenadminregistry.AcceptAdminRole{
				TokenConfigCid: tokenConfigCid,
				InstrumentId:   instrumentId,
				Caller:         types.PARTY(poolOwnerParty),
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
		}
		tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, tokenConfigAddress)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after accept admin role: %w", err)
		}
		if !tokenConfigFound {
			return sequences.OnChainOutput{}, fmt.Errorf("token config not found after accept admin role")
		}
	}

	// Step 3: SetPool (pool owner acts)
	poolOwnerPartyTyped := types.PARTY(poolOwnerParty)
	_, err = operations.ExecuteOperation(b, token_admin_registry.SetPool, deps, contract.ChoiceInput[tokenadminregistry.SetPool]{
		InstanceAddress: input.TokenAdminRegistryInstanceAddress,
		Args: tokenadminregistry.SetPool{
			TokenConfigCid: tokenConfigCid,
			InstrumentId:   instrumentId,
			TokenPool: &tokenadminregistry.PoolRegistration{
				PoolOwner:      poolOwnerPartyTyped,
				PoolInstanceId: types.TEXT(input.PoolInstanceID),
			},
			Caller: types.PARTY(poolOwnerParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to set pool: %w", err)
	}

	return sequences.OnChainOutput{}, nil
}

func findTokenConfigCid(b operations.Bundle, deps canton.Chain, address contracts.InstanceAddress) (types.CONTRACT_ID, bool, error) {
	participant := deps.Participants[0]
	contractID, err := contract.FindActiveContractIDByInstanceAddress(
		b.GetContext(),
		participant.LedgerServices.State,
		participant.PartyID,
		tokenadminregistry.TokenConfig{}.GetTemplateID(),
		address,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no active contract found") {
			return "", false, nil
		}

		return "", false, err
	}

	return types.CONTRACT_ID(contractID), true, nil
}
