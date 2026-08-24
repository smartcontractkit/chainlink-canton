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

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
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
	// CcipParty is the CCIP owner party (registry owner; acts on ProposeAdministrator).
	CcipParty string
	// PoolOwnerParty is the token pool contract owner (recorded in PoolRegistration only).
	PoolOwnerParty string
	// PoolAdminParty is the TokenConfig admin (proposed/accepts admin role; acts on SetPool).
	PoolAdminParty string
	// CcipParticipantIndex selects which participant submits TAR registration choices.
	// Zero value defaults to the first participant.
	CcipParticipantIndex int `json:"ccipParticipantIndex,omitempty"`
}

var RegisterTokenPool = operations.NewSequence(
	"canton/ccip/register_token_pool",
	semver.MustParse("2.0.0"),
	"Registers a token pool with the TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole, SetPool)",
	registerTokenPool,
)

func registerTokenPool(b operations.Bundle, deps canton.Chain, input RegisterTokenPoolInput) (sequences.OnChainOutput, error) {
	ccipParticipantIndex := input.CcipParticipantIndex
	ccipParticipant, err := contract.ParticipantAt(deps, ccipParticipantIndex)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("resolve ccip participant: %w", err)
	}

	instrumentId := input.InstrumentId
	ccipParty := input.CcipParty
	poolOwnerParty := input.PoolOwnerParty
	poolAdminParty := input.PoolAdminParty
	proposeMcmsEnabled := contract.ProposalDrivenForCaller(ccipParticipant, ccipParty)
	adminMcmsEnabled := contract.ProposalDrivenForCaller(ccipParticipant, poolAdminParty)

	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).RawInstanceAddress(types.PARTY(ccipParty)).InstanceAddress()

	var proposalOutputs []contract.ExerciseOutput

	if registered, err := tokenPoolAlreadyRegisteredWithTAR(b, deps, ccipParticipantIndex, tokenConfigAddress, input); err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("check existing token pool registration: %w", err)
	} else if registered {
		b.Logger.Infof(
			"token pool %s already registered in TAR for instrument %s; skipping ProposeAdministrator/AcceptAdminRole/SetPool",
			input.PoolInstanceID,
			instrumentId.Id,
		)

		return sequences.OnChainOutput{}, nil
	}

	existingTokenConfigCid, tokenConfigFound, err := findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
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

	// Step 1: ProposeAdministrator (ccip owner proposes pool admin as TokenConfig admin).
	skipAcceptAdminRole := false
	proposeReport, err := operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[core.ProposeAdministrator]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        proposeMcmsEnabled,
		ParticipantIndex:   ccipParticipantIndex,
		Args: core.ProposeAdministrator{
			TokenConfigCid: tokenConfigCidArg,
			InstrumentId:   instrumentId,
			NewAdmin:       types.PARTY(poolAdminParty),
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
	if err == nil && proposeMcmsEnabled && !proposeReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, proposeReport.Output)
	}

	if !proposeMcmsEnabled && err == nil {
		tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after propose: %w", err)
		}
		if !tokenConfigFound {
			return sequences.OnChainOutput{}, fmt.Errorf("token config not found after propose")
		}
	}

	// Step 2: AcceptAdminRole (pool admin accepts pending admin role).
	if !skipAcceptAdminRole {
		acceptReport, err := operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[core.AcceptAdminRole]{
			InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
			RawInstanceAddress: tarRaw,
			MCMSEnabled:        adminMcmsEnabled,
			ParticipantIndex:   ccipParticipantIndex,
			Args: core.AcceptAdminRole{
				TokenConfigCid: tokenConfigCid,
				InstrumentId:   instrumentId,
				Caller:         types.PARTY(poolAdminParty),
			},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to accept admin role: %w", err)
		}
		if adminMcmsEnabled && !acceptReport.Output.Executed() {
			proposalOutputs = append(proposalOutputs, acceptReport.Output)
		}
		if !adminMcmsEnabled {
			tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config after accept admin role: %w", err)
			}
			if !tokenConfigFound {
				return sequences.OnChainOutput{}, fmt.Errorf("token config not found after accept admin role")
			}
		}
	} else if !tokenConfigFound {
		tokenConfigCid, tokenConfigFound, err = findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to lookup token config for existing admin: %w", err)
		}
		if !tokenConfigFound {
			return sequences.OnChainOutput{}, fmt.Errorf("token config not found for instrument with admin already set")
		}
	}

	// Step 3: SetPool (pool admin sets pool registration for pool owner).
	poolOwnerPartyTyped := types.PARTY(poolOwnerParty)
	setPoolReport, err := operations.ExecuteOperation(b, token_admin_registry.SetPool, deps, contract.ChoiceInput[core.SetPool]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        adminMcmsEnabled,
		ParticipantIndex:   ccipParticipantIndex,
		Args: core.SetPool{
			TokenConfigCid: tokenConfigCid,
			InstrumentId:   instrumentId,
			TokenPool: &core.PoolRegistration2{
				PoolOwner:      poolOwnerPartyTyped,
				PoolInstanceId: types.TEXT(input.PoolInstanceID),
			},
			Caller: types.PARTY(poolAdminParty),
		},
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to set pool: %w", err)
	}
	if adminMcmsEnabled && !setPoolReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, setPoolReport.Output)
	}

	if !proposeMcmsEnabled && !adminMcmsEnabled {
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

func findTokenConfigCid(b operations.Bundle, deps canton.Chain, participantIndex int, address contracts.InstanceAddress) (types.CONTRACT_ID, bool, error) {
	participant, err := contract.ParticipantAt(deps, participantIndex)
	if err != nil {
		return "", false, fmt.Errorf("resolve participant: %w", err)
	}
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

func tokenPoolAlreadyRegisteredWithTAR(
	b operations.Bundle,
	deps canton.Chain,
	participantIndex int,
	tokenConfigAddress contracts.InstanceAddress,
	input RegisterTokenPoolInput,
) (bool, error) {
	participant, err := contract.ParticipantAt(deps, participantIndex)
	if err != nil {
		return false, fmt.Errorf("resolve participant: %w", err)
	}

	activeConfig, err := contract.FindActiveContractByInstanceAddress(
		b.GetContext(),
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		core.TokenConfig{}.GetTemplateID(),
		tokenConfigAddress,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no active contract found") {
			return false, nil
		}

		return false, err
	}

	config, err := bindings.UnmarshalCreatedEvent[core.TokenConfig](activeConfig.GetCreatedEvent())
	if err != nil {
		return false, fmt.Errorf("parse token config: %w", err)
	}
	if config.Admin == nil || config.TokenPool == nil {
		return false, nil
	}
	if string(*config.Admin) != input.PoolAdminParty {
		return false, nil
	}
	if string(config.TokenPool.PoolOwner) != input.PoolOwnerParty {
		return false, nil
	}
	if string(config.TokenPool.PoolInstanceId) != input.PoolInstanceID {
		return false, nil
	}

	return true, nil
}
