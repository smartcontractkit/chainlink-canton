package sequences

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/ccipcodec"
	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/globaltokensenderreceiver"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_token_sender_receiver"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployGlobalTokenSenderReceiverFromFactoryParams deploys the Canton GlobalTokenSenderReceiver
// via the ccip-qualified CCIPFactory, standalone from the full chain contract deploy.
type DeployGlobalTokenSenderReceiverFromFactoryParams struct {
	CCIPOwnerParty    string `json:"ccipOwnerParty,omitempty" yaml:"ccipOwnerParty,omitempty"`
	FactoryAddressRef datastore.AddressRef
	ProposalDriven    bool       `json:"proposalDriven,omitempty" yaml:"proposalDriven,omitempty"`
	InstanceID        types.TEXT `json:"instanceId,omitempty" yaml:"instanceId,omitempty"`
}

// DeployGlobalTokenSenderReceiverFromFactory emits an MCMS batch for DeployGlobalTokenSenderReceiver when ProposalDriven is true.
var DeployGlobalTokenSenderReceiverFromFactory = operations.NewSequence(
	"canton/ccip/deploy_global_token_sender_receiver_from_factory",
	semver.MustParse("2.0.0"),
	"Deploys GlobalTokenSenderReceiver on Canton through CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployGlobalTokenSenderReceiverFromFactoryParams) (sequences.OnChainOutput, error) {
		if input.CCIPOwnerParty == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("CCIPOwnerParty is required")
		}
		ccipOwnerParty := types.PARTY(input.CCIPOwnerParty)

		factoryRaw, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		instanceID, err := ensureInstanceID(input.InstanceID, "globaltokensenderreceiver")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("global token sender receiver instance ID: %w", err)
		}

		// Receiver-side policy is baked in; for token-only transfers the token pool decides
		// CCV and finality requirements, so the defaults below are inert.
		deployReport, err := operations.ExecuteOperation(b, factoryops.DeployGlobalTokenSenderReceiver, deps, newChoiceInput(factoryRaw, factorybindings.DeployGlobalTokenSenderReceiver{
			Contract: globaltokensenderreceiver.GlobalTokenSenderReceiver{
				InstanceId:             types.TEXT(instanceID),
				CcipOwner:              ccipOwnerParty,
				RequiredCCVs:           nil,
				OptionalCCVs:           nil,
				ReceiverFinalityConfig: ccipcodec.FinalityConfig{WaitForFinality: &types.UNIT{}},
			},
		}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("deploy GlobalTokenSenderReceiver from factory: %w", err)
		}

		globalTokenSenderReceiverRawAddr := instanceID.RawInstanceAddress(ccipOwnerParty)
		out := sequences.OnChainOutput{
			Addresses: []datastore.AddressRef{
				newAddressRef(deps.ChainSelector(), globalTokenSenderReceiverRawAddr, global_token_sender_receiver.ContractType, global_token_sender_receiver.Version, ""),
			},
		}

		if input.ProposalDriven && !deployReport.Output.Executed() {
			batchOp, err := contract.NewBatchOperationFromExercises([]contract.ExerciseOutput{deployReport.Output})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for GlobalTokenSenderReceiver deploy: %w", err)
			}
			if len(batchOp.Transactions) > 0 {
				out.BatchOps = []mcms_types.BatchOperation{batchOp}
			}
		}

		return out, nil
	},
)
