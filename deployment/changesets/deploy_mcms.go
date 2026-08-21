package changesets

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	ccipdeploymentutils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	mcmsApi "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/mcms/api"
	mcmsCore "github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/mcms/core"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
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

type DeployAndConfigureMCMSConfig struct {
	Params DeployAndConfigureMCMSParams `json:"params" yaml:"params"`
}

type DeployAndConfigureMCMS struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployAndConfigureMCMSConfig]] = DeployAndConfigureMCMS{}

func (d DeployAndConfigureMCMS) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployAndConfigureMCMSConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}

	params := config.Config.Params
	if params.OwnerParty == "" {
		return fmt.Errorf("owner party is required")
	}
	if params.ChainID <= 0 {
		return fmt.Errorf("chain ID must be greater than zero")
	}

	return nil
}

func (d DeployAndConfigureMCMS) Apply(e cldf.Environment, config CantonCSDeps[DeployAndConfigureMCMSConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	out, err := operations.ExecuteSequence(e.OperationsBundle, deployAndConfigureMCMSSequence, chain, config.Config.Params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployAndConfigureMCMS sequence: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}

var deployAndConfigureMCMSSequence = operations.NewSequence(
	"canton/mcms/deploy_and_configure",
	semver.MustParse("0.1.0"),
	"Deploys and configures a Canton MCMS contract",
	func(b operations.Bundle, deps canton.Chain, input DeployAndConfigureMCMSParams) (ccipsequences.OnChainOutput, error) {
		if input.OwnerParty == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("owner party is required")
		}
		if input.ChainID <= 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("chain ID must be greater than zero")
		}

		initialConfig, err := buildMultisigConfig(input.InitialConfig)
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("build initial MCMS config: %w", err)
		}

		roleState := emptyRoleState(initialConfig)
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
			return ccipsequences.OnChainOutput{}, fmt.Errorf("deploy MCMS: %w", err)
		}

		if len(deployReport.Output.Labels.List()) == 0 {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("missing raw MCMS instance address label in deploy output")
		}
		rawInstanceAddress, err := contracts.RawInstanceAddressFromString(deployReport.Output.Labels.List()[0])
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("parse raw MCMS instance address label: %w", err)
		}

		for i, roleConfig := range input.RoleConfigs {
			groupConfig, err := buildNormalizedConfig(roleConfig.Config)
			if err != nil {
				return ccipsequences.OnChainOutput{}, fmt.Errorf("build MCMS config for role %s: %w", roleConfig.Role, err)
			}

			_, err = operations.ExecuteOperation(b, mcmsops.SetConfig, deps, opcontract.ChoiceInput[mcmsCore.SetConfig]{
				InstanceAddress: rawInstanceAddress.InstanceAddress(),
				Args: mcmsCore.SetConfig{
					TargetRole:      roleConfig.Role,
					NewSigners:      groupConfig.Signers,
					NewGroupQuorums: groupConfig.GroupQuorums,
					NewGroupParents: groupConfig.GroupParents,
					ClearRoot:       types.BOOL(roleConfig.Config.ClearRoot),
				},
			})
			if err != nil {
				return ccipsequences.OnChainOutput{}, fmt.Errorf("configure MCMS role %s at index %d: %w", roleConfig.Role, i, err)
			}
		}

		refs := []datastore.AddressRef{
			deployReport.Output,
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.ProposerManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.CancellerManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.BypasserManyChainMultisig), qualifierOrDefault(input.Qualifier)),
			newMCMSRoleAddressRef(deps.ChainSelector(), rawInstanceAddress, datastore.ContractType(ccipdeploymentutils.RBACTimelock), qualifierOrDefault(input.Qualifier)),
		}

		return ccipsequences.OnChainOutput{Addresses: refs}, nil
	},
)

func buildMultisigConfig(input MCMSConfigParams) (mcmsApi.MultisigConfig, error) {
	cfg, err := buildNormalizedConfig(input)
	if err != nil {
		return mcmsApi.MultisigConfig{}, err
	}

	return mcmsApi.MultisigConfig{
		Signers:      cfg.Signers,
		GroupQuorums: cfg.GroupQuorums,
		GroupParents: cfg.GroupParents,
	}, nil
}

func buildNormalizedConfig(input MCMSConfigParams) (MCMSConfigParams, error) {
	groupQuorums, err := normalizeGroups(input.GroupQuorums, "group quorums")
	if err != nil {
		return MCMSConfigParams{}, err
	}
	groupParents, err := normalizeGroups(input.GroupParents, "group parents")
	if err != nil {
		return MCMSConfigParams{}, err
	}

	return MCMSConfigParams{
		Signers:      input.Signers,
		GroupQuorums: groupQuorums,
		GroupParents: groupParents,
		ClearRoot:    input.ClearRoot,
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
