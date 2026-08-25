package sequences

import (
	"encoding/hex"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// RegisterNativeFeeTokenInTARInput registers the Canton native fee token (typically Amulet) in the
// TokenAdminRegistry by creating a TokenConfig via ProposeAdministrator + AcceptAdminRole.
// No SetPool is performed — fee tokens are not CCIP transfer tokens.
type RegisterNativeFeeTokenInTARInput struct {
	TokenAdminRegistryInstanceAddress    contracts.InstanceAddress
	TokenAdminRegistryRawInstanceAddress contracts.RawInstanceAddress
	InstrumentId                         splice_api_token_holding_v1.InstrumentId
	CcipOwnerParty                       string
	TokenQualifier                       string
	ChainSelector                        uint64
	CcipParticipantIndex                 int `json:"ccipParticipantIndex,omitempty"`
	ProposalDriven                       bool
}

var RegisterNativeFeeTokenInTAR = operations.NewSequence(
	"canton/ccip/register_native_fee_token_in_tar",
	semver.MustParse("2.0.0"),
	"Registers the native Canton fee token in TokenAdminRegistry (ProposeAdministrator, AcceptAdminRole)",
	registerNativeFeeTokenInTAR,
)

func registerNativeFeeTokenInTAR(b operations.Bundle, deps canton.Chain, input RegisterNativeFeeTokenInTARInput) (sequences.OnChainOutput, error) {
	proposalOutputs, addrRef, err := registerNativeFeeTokenInTARCore(b, deps, input)
	if err != nil {
		return sequences.OnChainOutput{}, err
	}

	var batchOps []mcms_types.BatchOperation
	if len(proposalOutputs) > 0 {
		batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for native fee token registration: %w", err)
		}
		batchOps = []mcms_types.BatchOperation{batchOp}
	}

	var addresses []datastore.AddressRef
	if addrRef != nil {
		addresses = []datastore.AddressRef{*addrRef}
	}

	return sequences.OnChainOutput{
		BatchOps:  batchOps,
		Addresses: addresses,
	}, nil
}

func registerNativeFeeTokenInTARCore(
	b operations.Bundle,
	deps canton.Chain,
	input RegisterNativeFeeTokenInTARInput,
) ([]contract.ExerciseOutput, *datastore.AddressRef, error) {
	instrumentId := input.InstrumentId
	if instrumentId.Admin == "" || instrumentId.Id == "" {
		return nil, nil, nil
	}
	if input.CcipOwnerParty == "" {
		return nil, nil, fmt.Errorf("CcipOwnerParty is required")
	}

	ccipParticipantIndex := input.CcipParticipantIndex
	ccipOwnerParty := input.CcipOwnerParty
	ccipOwner := types.PARTY(ccipOwnerParty)

	tokenConfigAddress := contracts.InstanceID(hex.EncodeToString(contracts.EncodeInstrumentID(instrumentId).Bytes())).
		RawInstanceAddress(ccipOwner).
		InstanceAddress()

	existingTokenConfigCid, found, err := findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("lookup native fee token config: %w", err)
	}
	if found {
		b.Logger.Infof(
			"native fee token %s already registered in TAR; skipping ProposeAdministrator/AcceptAdminRole",
			instrumentId.Id,
		)

		return nil, nil, nil
	}

	ccipParticipant, err := contract.ParticipantAt(deps, ccipParticipantIndex)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve ccip participant: %w", err)
	}
	adminMcmsEnabled := input.ProposalDriven && contract.ProposalDrivenForCaller(ccipParticipant, ccipOwnerParty)
	tarRaw := input.TokenAdminRegistryRawInstanceAddress.String()

	var proposalOutputs []contract.ExerciseOutput

	var tokenConfigCid types.CONTRACT_ID
	var tokenConfigCidArg *types.CONTRACT_ID
	if existingTokenConfigCid != "" {
		tokenConfigCid = existingTokenConfigCid
		tokenConfigCidArg = &existingTokenConfigCid
	}

	proposeReport, err := operations.ExecuteOperation(b, token_admin_registry.ProposeAdministrator, deps, contract.ChoiceInput[core.ProposeAdministrator]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        adminMcmsEnabled,
		ParticipantIndex:   ccipParticipantIndex,
		Args: core.ProposeAdministrator{
			TokenConfigCid: tokenConfigCidArg,
			InstrumentId:   instrumentId,
			NewAdmin:       ccipOwner,
			Caller:         ccipOwner,
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("propose native fee token admin: %w", err)
	}
	if adminMcmsEnabled && !proposeReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, proposeReport.Output)
	}

	if !adminMcmsEnabled {
		tokenConfigCid, found, err = findTokenConfigCid(b, deps, ccipParticipantIndex, tokenConfigAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("lookup native fee token config after propose: %w", err)
		}
		if !found {
			return nil, nil, fmt.Errorf("native fee token config not found after propose")
		}
	}

	acceptReport, err := operations.ExecuteOperation(b, token_admin_registry.AcceptAdminRole, deps, contract.ChoiceInput[core.AcceptAdminRole]{
		InstanceAddress:    input.TokenAdminRegistryInstanceAddress,
		RawInstanceAddress: tarRaw,
		MCMSEnabled:        adminMcmsEnabled,
		ParticipantIndex:   ccipParticipantIndex,
		Args: core.AcceptAdminRole{
			TokenConfigCid: tokenConfigCid,
			InstrumentId:   instrumentId,
			Caller:         ccipOwner,
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("accept native fee token admin role: %w", err)
	}
	if adminMcmsEnabled && !acceptReport.Output.Executed() {
		proposalOutputs = append(proposalOutputs, acceptReport.Output)
	}

	qualifier := input.TokenQualifier
	if qualifier == "" {
		qualifier = string(instrumentId.Id)
	}

	addrRef := &datastore.AddressRef{
		Address: contracts.EncodeInstrumentID(instrumentId).Hex(),
		Labels: datastore.NewLabelSet(
			fmt.Sprintf("instrument-admin:%s", instrumentId.Admin),
			fmt.Sprintf("instrument-id:%s", instrumentId.Id),
			fmt.Sprintf("ccip-owner:%s", ccipOwnerParty),
		),
		Type:          datastore.ContractType("Token"),
		Version:       token_admin_registry.Version,
		Qualifier:     qualifier,
		ChainSelector: input.ChainSelector,
	}

	return proposalOutputs, addrRef, nil
}
