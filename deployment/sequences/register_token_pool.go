package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/tokenadminregistry"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/splice/splice_api_token_holding_v1"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/dependencies"
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

func registerTokenPool(b operations.Bundle, deps dependencies.CantonDeps, input RegisterTokenPoolInput) (sequences.OnChainOutput, error) {
	instrumentId := input.InstrumentId
	ccipParty := input.CcipParty
	poolOwnerParty := input.PoolOwnerParty

	// Step 1: ProposeAdministrator (CCIP acts)
	_, err := operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[tokenadminregistry.TokenAdminRegistryProposeAdministrator]{
		ChainSelector:   deps.Chain.Selector,
		InstanceAddress: input.TokenAdminRegistryInstanceAddress,
		ActAs:           []string{ccipParty},
		Args: tokenadminregistry.TokenAdminRegistryProposeAdministrator{
			InstrumentId: instrumentId,
			NewAdmin:     types.PARTY(poolOwnerParty),
			Caller:       types.PARTY(ccipParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to propose administrator: %w", err)
	}

	// Step 2: AcceptAdminRole (pool owner acts). Exercise resolves current TAR contract by InstanceAddress.
	_, err = operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[tokenadminregistry.TokenAdminRegistryAcceptAdminRole]{
		ChainSelector:   deps.Chain.Selector,
		InstanceAddress: input.TokenAdminRegistryInstanceAddress,
		ActAs:           []string{poolOwnerParty},
		Args: tokenadminregistry.TokenAdminRegistryAcceptAdminRole{
			InstrumentId: instrumentId,
			Caller:       types.PARTY(poolOwnerParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
	}

	// Step 3: SetPool (pool owner acts)
	poolOwnerPartyTyped := types.PARTY(poolOwnerParty)
	_, err = operations.ExecuteOperation(b, token_admin_registry.SetPool, deps, contract.ChoiceInput[tokenadminregistry.TokenAdminRegistrySetPool]{
		ChainSelector:   deps.Chain.Selector,
		InstanceAddress: input.TokenAdminRegistryInstanceAddress,
		ActAs:           []string{poolOwnerParty},
		Args: tokenadminregistry.TokenAdminRegistrySetPool{
			InstrumentId: instrumentId,
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
