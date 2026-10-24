package sequences

import (
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	ccipdeploy "github.com/smartcontractkit/chainlink-ccip/deployment/deploy"
	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/v2/bindings/generated/mcms/core"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

const mcmsGroupCount = 32

type MCMSConfigParams struct {
	Signers      []mcmsApi.SignerInfo `json:"signers" yaml:"signers"`
	GroupQuorums []types.INT64        `json:"groupQuorums" yaml:"groupQuorums"`
	GroupParents []types.INT64        `json:"groupParents" yaml:"groupParents"`
	ClearRoot    bool                 `json:"clearRoot" yaml:"clearRoot"`
}

type MCMSRoleConfigParams struct {
	Role   mcmsApi.Role     `json:"role" yaml:"role"`
	Config MCMSConfigParams `json:"config" yaml:"config"`
}

type DeployAndConfigureMCMSParams struct {
	OwnerParty       string                    `json:"ownerParty" yaml:"ownerParty"`
	InstanceID       string                    `json:"instanceID,omitempty" yaml:"instanceID,omitempty"`
	ChainID          int64                     `json:"chainID" yaml:"chainID"`
	Qualifier        string                    `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
	MinDelay         time.Duration             `json:"minDelay" yaml:"minDelay"`
	BlockedFunctions []mcmsApi.BlockedFunction `json:"blockedFunctions" yaml:"blockedFunctions"`
	InitialConfig    MCMSConfigParams          `json:"initialConfig" yaml:"initialConfig"`
	RoleConfigs      []MCMSRoleConfigParams    `json:"roleConfigs" yaml:"roleConfigs"`
}

// MCMSRoleSetConfig is the new signer config for one MCMS role.
type MCMSRoleSetConfig struct {
	Role   mcmsApi.Role
	Params mcmsApi.SetConfigParams
}

type SetMCMSConfigInput struct {
	RawInstanceAddress contracts.RawInstanceAddress
	Roles              []MCMSRoleSetConfig
}

var DeployAndConfigureMCMS = operations.NewSequence(
	"canton/mcms/deploy_and_configure",
	semver.MustParse("0.1.0"),
	"Deploys and configures a Canton MCMS contract",
	func(b operations.Bundle, deps canton.Chain, input DeployAndConfigureMCMSParams) (sequences.OnChainOutput, error) {
		if input.OwnerParty == "" {
			return sequences.OnChainOutput{}, fmt.Errorf("owner party is required")
		}
		if input.ChainID <= 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("chain ID must be greater than zero")
		}

		initialConfig, err := buildSetConfigParams(input.InitialConfig)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("build initial MCMS config: %w", err)
		}

		roleState := emptyRoleState(mcmsApi.MultisigConfig{
			Signers:      initialConfig.Signers,
			GroupQuorums: initialConfig.GroupQuorums,
			GroupParents: initialConfig.GroupParents,
		})
		ownerParty := types.PARTY(input.OwnerParty)

		deployReport, err := operations.ExecuteOperation(b, mcmsops.Deploy, deps, opcontract.DeployInput[mcmsCore.MCMS]{
			Template: mcmsCore.MCMS{
				Owner:              ownerParty,
				InstanceId:         types.TEXT(input.InstanceID),
				ChainId:            types.INT64(input.ChainID),
				Proposer:           roleState,
				Canceller:          roleState,
				Bypasser:           roleState,
				MinDelay:           types.RELTIME(input.MinDelay),
				BlockedFunctions:   input.BlockedFunctions,
				TimelockTimestamps: map[types.TEXT]types.TIMESTAMP{},
			},
			OwnerParty: ownerParty,
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("deploy MCMS: %w", err)
		}

		rawInstanceAddress, err := dsutils.MCMSRawInstanceAddressFromRef(deployReport.Output)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		roles := make([]MCMSRoleSetConfig, len(input.RoleConfigs))
		for i, roleConfig := range input.RoleConfigs {
			params, err := buildSetConfigParams(roleConfig.Config)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS config for role %s: %w", roleConfig.Role, err)
			}
			roles[i] = MCMSRoleSetConfig{Role: roleConfig.Role, Params: params}
		}
		if len(roles) > 0 {
			// The deployer is the owner here, so SetConfig runs directly.
			_, err = operations.ExecuteSequence(b, SetMCMSConfig, deps, SetMCMSConfigInput{RawInstanceAddress: rawInstanceAddress, Roles: roles})
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
		}

		refs := []datastore.AddressRef{
			deployReport.Output,
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.CancellerManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.BypasserManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.RBACTimelock), qualifierOrDefault(input.Qualifier)),
		}

		return sequences.OnChainOutput{Addresses: refs}, nil
	},
)

// SetMCMSConfig runs SetConfig directly if the deployer is the MCMS owner.
// Otherwise it returns a self-dispatch batch op for a proposal.
var SetMCMSConfig = operations.NewSequence(
	"canton/mcms/set_config",
	semver.MustParse("1.0.0"),
	"Sets the signer config of Canton MCMS roles, directly or through an MCMS proposal",
	func(b operations.Bundle, deps canton.Chain, input SetMCMSConfigInput) (sequences.OnChainOutput, error) {
		if len(deps.Participants) == 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("canton chain %d has no participants", deps.ChainSelector())
		}

		if !opcontract.ProposalDrivenForCaller(deps.Participants[0], input.RawInstanceAddress.Owner()) {
			for _, r := range input.Roles {
				_, err := operations.ExecuteOperation(b, mcmsops.SetConfig, deps, opcontract.ChoiceInput[mcmsCore.SetConfig]{
					InstanceAddress: input.RawInstanceAddress.InstanceAddress(),
					Args: mcmsCore.SetConfig{
						TargetRole:      r.Role,
						NewSigners:      r.Params.Signers,
						NewGroupQuorums: r.Params.GroupQuorums,
						NewGroupParents: r.Params.GroupParents,
						ClearRoot:       r.Params.ClearRoot,
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("set config for MCMS role %s on %s: %w", r.Role, input.RawInstanceAddress, err)
				}
			}

			return sequences.OnChainOutput{}, nil
		}

		batch := mcms_types.BatchOperation{ChainSelector: mcms_types.ChainSelector(deps.ChainSelector())}
		for _, r := range input.Roles {
			tx, err := selfDispatchSetConfigTransaction(input.RawInstanceAddress, r.Role, r.Params)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build SetConfig proposal transaction for MCMS role %s on %s: %w", r.Role, input.RawInstanceAddress, err)
			}
			batch.Transactions = append(batch.Transactions, tx)
		}

		return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batch}}, nil
	},
)

// UpdateMCMSConfig is the ccip Deployer UpdateMCMSConfig step for Canton. It calls SetMCMSConfig per MCMS ref
// and merges any proposal transactions into one batch op.
var UpdateMCMSConfig = operations.NewSequence(
	"canton/mcms/update_mcms_config",
	semver.MustParse("1.0.0"),
	"Updates the signer config of Canton MCMS roles, directly or through an MCMS proposal",
	func(b operations.Bundle, chains cldf_chain.BlockChains, in ccipdeploy.UpdateMCMSConfigInputPerChainWithSelector) (sequences.OnChainOutput, error) {
		chain, ok := chains.CantonChains()[in.ChainSelector]
		if !ok || len(chain.Participants) == 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("canton chain %d not found or has no participants", in.ChainSelector)
		}

		params, err := cantonSetConfigParams(in.MCMConfig)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		batch := mcms_types.BatchOperation{ChainSelector: mcms_types.ChainSelector(in.ChainSelector)}
		for _, ref := range in.MCMContracts {
			role, err := cantonmcms.RoleForContractType(ref.Type)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
			rawInstanceAddress, err := dsutils.MCMSRawInstanceAddressFromRef(ref)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}

			report, err := operations.ExecuteSequence(b, SetMCMSConfig, chain, SetMCMSConfigInput{
				RawInstanceAddress: rawInstanceAddress,
				Roles:              []MCMSRoleSetConfig{{Role: role.Role, Params: params}},
			})
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
			for _, op := range report.Output.BatchOps {
				batch.Transactions = append(batch.Transactions, op.Transactions...)
			}
		}

		if len(batch.Transactions) == 0 {
			return sequences.OnChainOutput{}, nil
		}

		return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batch}}, nil
	},
)

// cantonSetConfigParams converts cfg like the mcms Canton SDK Configurer.
func cantonSetConfigParams(cfg mcms_types.Config) (mcmsApi.SetConfigParams, error) {
	groupQuorums, groupParents, signerAddresses, signerGroups, err := mcmssdk.ExtractSetConfigInputs(&cfg)
	if err != nil {
		return mcmsApi.SetConfigParams{}, fmt.Errorf("extract MCMS config: %w", err)
	}

	signers := make([]mcmsApi.SignerInfo, len(signerAddresses))
	for i, addr := range signerAddresses {
		signers[i] = mcmsApi.SignerInfo{
			SignerAddress: types.TEXT(strings.TrimPrefix(strings.ToLower(addr.Hex()), "0x")),
			SignerIndex:   types.INT64(i),
			SignerGroup:   types.INT64(signerGroups[i]),
		}
	}

	quorums := make([]types.INT64, len(groupQuorums))
	for i, q := range groupQuorums {
		quorums[i] = types.INT64(q)
	}
	parents := make([]types.INT64, len(groupParents))
	for i, p := range groupParents {
		parents[i] = types.INT64(p)
	}

	return mcmsApi.SetConfigParams{
		Signers:      signers,
		GroupQuorums: quorums,
		GroupParents: parents,
		ClearRoot:    false,
	}, nil
}

func selfDispatchSetConfigTransaction(
	rawInstanceAddress contracts.RawInstanceAddress,
	role mcmsApi.Role,
	params mcmsApi.SetConfigParams,
) (mcms_types.Transaction, error) {
	operationData, err := mcmsops.EncodeSelfDispatchSetConfig(role, params)
	if err != nil {
		return mcms_types.Transaction{}, err
	}

	// No template ID: MCMS rejects a target contract ID for calls to itself.
	return opcontract.NewCantonTransaction(
		rawInstanceAddress.String(),
		rawInstanceAddress.InstanceAddress(),
		&bind.EncodedChoice{Choice: mcmsops.SelfDispatchSetConfigFunction, OperationData: operationData},
		mcmsops.ContractType,
		"",
	)
}

func buildSetConfigParams(input MCMSConfigParams) (mcmsApi.SetConfigParams, error) {
	groupQuorums, err := normalizeGroups(input.GroupQuorums, "group quorums")
	if err != nil {
		return mcmsApi.SetConfigParams{}, err
	}
	groupParents, err := normalizeGroups(input.GroupParents, "group parents")
	if err != nil {
		return mcmsApi.SetConfigParams{}, err
	}

	return mcmsApi.SetConfigParams{
		Signers:      input.Signers,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
		ClearRoot:    types.BOOL(input.ClearRoot),
	}, nil
}

func normalizeGroups(input []types.INT64, name string) ([]types.INT64, error) {
	if len(input) > mcmsGroupCount {
		return nil, fmt.Errorf("%s length %d exceeds max %d", name, len(input), mcmsGroupCount)
	}

	out := make([]types.INT64, mcmsGroupCount)
	copy(out, input)

	return out, nil
}

func emptyRoleState(config mcmsApi.MultisigConfig) mcmsApi.RoleState {
	return mcmsApi.RoleState{
		Config:     config,
		SeenHashes: map[types.TEXT]types.TIMESTAMP{},
		ExpiringRoot: mcmsApi.ExpiringRoot{
			Root:       types.TEXT(""),
			ValidUntil: types.TIMESTAMP(time.Unix(0, 0)),
			OpCount:    types.INT64(0),
		},
		RootMetadata: mcmsApi.RootMetadata{
			ChainId:              types.INT64(0),
			MultisigId:           types.TEXT(""),
			PreOpCount:           types.INT64(0),
			PostOpCount:          types.INT64(0),
			OverridePreviousRoot: types.BOOL(false),
		},
	}
}

func qualifierOrDefault(qualifier string) string {
	if qualifier != "" {
		return qualifier
	}

	return ccipdeploymentutils.CLLQualifier
}

func newMCMSRoleAddressRef(chainSelector uint64, rawInstanceAddress contracts.RawInstanceAddress, contractType datastore.ContractType, qualifier string) datastore.AddressRef {
	return datastore.AddressRef{
		Address:       rawInstanceAddress.InstanceAddress().String(),
		Labels:        datastore.NewLabelSet(rawInstanceAddress.String()),
		ChainSelector: chainSelector,
		Type:          contractType,
		Version:       mcmsops.Version,
		Qualifier:     qualifier,
	}
}
