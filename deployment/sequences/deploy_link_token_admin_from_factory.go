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

	factorybindings "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/link/tokenadmin"
	"github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/splice/splice_api_token_holding_v1"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// DeployLinkTokenAdminFromFactoryParams deploys LinkTokenAdmin via the ccip-qualified CCIPFactory.
// InstrumentID.Id is required: it must name the registry-issued LINK instrument this admin governs.
type DeployLinkTokenAdminFromFactoryParams struct {
	CCIPOwnerParty    string `json:"ccipOwnerParty,omitempty" yaml:"ccipOwnerParty,omitempty"`
	FactoryAddressRef datastore.AddressRef
	ProposalDriven    bool                                     `json:"proposalDriven,omitempty" yaml:"proposalDriven,omitempty"`
	InstanceID        types.TEXT                               `json:"instanceId,omitempty" yaml:"instanceId,omitempty"`
	InstrumentID      splice_api_token_holding_v1.InstrumentId `json:"instrumentId" yaml:"instrumentId"`
	Observers         []types.PARTY                            `json:"observers,omitempty" yaml:"observers,omitempty"`
}

// DeployLinkTokenAdminFromFactory emits an MCMS batch for DeployLinkTokenAdmin when ProposalDriven is true.
var DeployLinkTokenAdminFromFactory = operations.NewSequence(
	"canton/ccip/deploy_link_token_admin_from_factory",
	semver.MustParse("1.0.0"),
	"Deploys LinkTokenAdmin on Canton through CCIPFactory",
	func(b operations.Bundle, deps canton.Chain, input DeployLinkTokenAdminFromFactoryParams) (sequences.OnChainOutput, error) {
		if input.CCIPOwnerParty == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("CCIPOwnerParty is required")
		}
		ccipOwnerParty := types.PARTY(input.CCIPOwnerParty)

		instrumentID := input.InstrumentID
		if instrumentID.Id == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("InstrumentID.Id is required")
		}
		if instrumentID.Admin == "" {
			instrumentID.Admin = ccipOwnerParty
		}
		if instrumentID.Admin != ccipOwnerParty {
			return sequences.OnChainOutput{}, fmt.Errorf("instrument admin must match ccipOwnerParty")
		}

		factoryRaw, err := rawInstanceAddressFromAddressRef(input.FactoryAddressRef)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		adminInstanceID, err := ensureInstanceID(input.InstanceID, "link-token-admin")
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("link token admin instance ID: %w", err)
		}

		deployReport, err := operations.ExecuteOperation(b, factoryops.DeployLinkTokenAdmin, deps, newChoiceInput(factoryRaw, factorybindings.DeployLinkTokenAdmin{
			Contract: tokenadmin.LinkTokenAdmin{
				CcipOwner:    ccipOwnerParty,
				InstanceId:   types.TEXT(adminInstanceID),
				InstrumentId: instrumentID,
				Observers:    input.Observers,
			},
		}, input.ProposalDriven))
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("deploy LinkTokenAdmin from factory: %w", err)
		}

		out := sequences.OnChainOutput{
			Addresses: []datastore.AddressRef{
				newAddressRef(deps.ChainSelector(), adminInstanceID.RawInstanceAddress(ccipOwnerParty), linkadmingov.ContractType, linkadmingov.Version, ""),
			},
		}

		if input.ProposalDriven && !deployReport.Output.Executed() {
			batchOp, err := contract.NewBatchOperationFromExercises([]contract.ExerciseOutput{deployReport.Output})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for link token admin deploy: %w", err)
			}
			out.BatchOps = []mcms_types.BatchOperation{batchOp}
		}

		return out, nil
	},
)
