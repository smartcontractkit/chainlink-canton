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
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// RegisterTokenPoolInput is the input for registering a token pool with the TokenAdminRegistry.
type RegisterTokenPoolInput struct {
	// TokenAdminRegistryInstanceAddress is the instance address of the TokenAdminRegistry contract.
	TokenAdminRegistryInstanceAddress contracts.InstanceAddress
	// TokenAdminRegistryRawInstanceAddress is the raw instance address label for MCMS proposals.
	TokenAdminRegistryRawInstanceAddress contracts.RawInstanceAddress
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
	semver.MustParse("2.0.0"),
	"Registers a token pool with the TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole, SetPool)",
	registerTokenPool,
)

func registerTokenPool(b operations.Bundle, deps canton.Chain, input RegisterTokenPoolInput) (sequences.OnChainOutput, error) {
	participant := deps.Participants[0]
	// ReadAsPartyIDs are CanReadAs rights for parties the operator cannot ActAs (e.g. ccip owner).
	// When present, exercises must be encoded as MCMS proposals instead of submitted directly.
	mcmsEnabled := len(participant.ReadAsPartyIDs) > 0

	instrumentId := input.InstrumentId
	ccipParty := input.CcipParty
	poolOwnerParty := input.PoolOwnerParty
	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).RawInstanceAddress(types.PARTY(ccipParty)).InstanceAddress()

	var proposalOutputs []contract.ExerciseOutput

	existingTokenConfigCid, tokenConfigFound, err := findTokenConfigCid(b, deps, tokenConfigAddress)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config: %w", err)
	}
	var tokenConfigCid types.CONTRACT_ID
	var tokenConfigCidArg *types.CONTRACT_ID
	if tokenConfigFound {
		tokenConfigCid = existingTokenConfigCid
		tokenConfigCidArg = &existingTokenConfigCid
	}

	tarRaw := input.TokenAdminRegistryRawInstanceAddress.String()

	// Step 1: ProposeAdministrator (CCIP acts)
	skipAcceptAdminRole := false
	proposeReport, err := operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[core.ProposeAdministrator]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        mcmsEnabled,
		Args: core.ProposeAdministrator{
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
	if err == nil && mcmsEnabled && !proposeReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, proposeReport.Output)
	}

	if !mcmsEnabled && err == nil {
		tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, tokenConfigAddress)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after propose: %w", err)
		}
		if !tokenConfigFound {
			return sequences.OnChainOutput{}, fmt.Errorf("token config not found after propose")
		}
	}

	// Step 2: AcceptAdminRole (pool owner acts). Exercise resolves current TAR contract by InstanceAddress.
	if !skipAcceptAdminRole {
		acceptReport, err := operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[core.AcceptAdminRole]{
			InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
			RawInstanceAddress: tarRaw,
			MCMSEnabled:        mcmsEnabled,
			Args: core.AcceptAdminRole{
				TokenConfigCid: tokenConfigCid,
				InstrumentId:   instrumentId,
				Caller:         types.PARTY(poolOwnerParty),
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
		}
		if mcmsEnabled && !acceptReport.Output.Executed() {
			proposalOutputs = append(proposalOutputs, acceptReport.Output)
		}
		if !mcmsEnabled {
			tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, tokenConfigAddress)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after accept admin role: %w", err)
			}
			if !tokenConfigFound {
				return sequences.OnChainOutput{}, fmt.Errorf("token config not found after accept admin role")
			}
		}
	}

	// Step 3: SetPool (pool owner acts)
	poolOwnerPartyTyped := types.PARTY(poolOwnerParty)
	setPoolReport, err := operations.ExecuteOperation(b, token_admin_registry.SetPool, deps, contract.ChoiceInput[core.SetPool]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        mcmsEnabled,
		Args: core.SetPool{
			TokenConfigCid: tokenConfigCid,
			InstrumentId:   instrumentId,
			TokenPool: &core.PoolRegistration{
				PoolOwner:      poolOwnerPartyTyped,
				PoolInstanceId: types.TEXT(input.PoolInstanceID),
			},
			Caller: types.PARTY(poolOwnerParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to set pool: %w", err)
	}
	if mcmsEnabled && !setPoolReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, setPoolReport.Output)
	}

	if !mcmsEnabled {
		return sequences.OnChainOutput{}, nil
	}
	batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for token pool registration: %w", err)
	}
	if len(batchOp.Transactions) == 0 {
		return sequences.OnChainOutput{}, nil
	}

	return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
}

func findTokenConfigCid(b operations.Bundle, deps canton.Chain, address contracts.InstanceAddress) (types.CONTRACT_ID, bool, error) {
	participant := deps.Participants[0]
	contractID, err := contract.FindActiveContractIDByInstanceAddress(
		b.GetContext(),
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		core.TokenConfig{}.GetTemplateID(),
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
