package mcms

import (
	"fmt"

	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
)

// Role links one Canton MCMS role to its timelock action, ccip contract type and SDK role.
type Role struct {
	Action       mcms_types.TimelockAction
	ContractType datastore.ContractType
	Role         mcmsApi.Role
	SDKRole      cantonsdk.TimelockRole
}

var roles = []Role{
	{
		Action:       mcms_types.TimelockActionSchedule,
		ContractType: datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig),
		Role:         mcmsApi.RoleProposer,
		SDKRole:      cantonsdk.TimelockRoleProposer,
	},
	{
		Action:       mcms_types.TimelockActionCancel,
		ContractType: datastore.ContractType(ccipdeploymentutils.CancellerManyChainMultisig),
		Role:         mcmsApi.RoleCanceller,
		SDKRole:      cantonsdk.TimelockRoleCanceller,
	},
	{
		Action:       mcms_types.TimelockActionBypass,
		ContractType: datastore.ContractType(ccipdeploymentutils.BypasserManyChainMultisig),
		Role:         mcmsApi.RoleBypasser,
		SDKRole:      cantonsdk.TimelockRoleBypasser,
	},
}

// RoleForAction returns the role for a timelock action. Unknown actions get the proposer.
func RoleForAction(action mcms_types.TimelockAction) Role {
	for _, r := range roles {
		if r.Action == action {
			return r
		}
	}

	return roles[0]
}

// RoleForContractType returns the role for a ccip MCMS contract type.
func RoleForContractType(contractType datastore.ContractType) (Role, error) {
	for _, r := range roles {
		if r.ContractType == contractType {
			return r, nil
		}
	}

	return Role{}, fmt.Errorf("unsupported MCMS contract type %q", contractType)
}
