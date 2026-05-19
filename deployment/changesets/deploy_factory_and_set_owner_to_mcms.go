package changesets

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	ccipsequences "github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type DeployCCIPFactoryParams struct {
	OwnerParty string `json:"ownerParty" yaml:"ownerParty"`
	MCMSParty  string `json:"mcmsParty,omitempty" yaml:"mcmsParty,omitempty"`
	InstanceID string `json:"instanceID,omitempty" yaml:"instanceID,omitempty"`
	Qualifier  string `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
}

type DeployCCIPFactoryConfig struct {
	Params DeployCCIPFactoryParams `json:"params" yaml:"params"`
}

type DeployCCIPFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCCIPFactoryConfig]] = DeployCCIPFactory{}

func (d DeployCCIPFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCCIPFactoryConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", config.Participant, config.ChainSelector, len(chain.Participants))
	}
	if config.Config.Params.OwnerParty == "" {
		return fmt.Errorf("owner party is required")
	}
	if config.Config.Params.Qualifier == "" {
		return fmt.Errorf("factory qualifier is required")
	}

	return nil
}

func (d DeployCCIPFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCCIPFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	out, err := operations.ExecuteSequence(e.OperationsBundle, deployCCIPFactorySequence, chain, config.Config.Params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute DeployCCIPFactory sequence: %w", err)
	}
	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}

var deployCCIPFactorySequence = operations.NewSequence(
	"canton/ccip_factory/deploy",
	semver.MustParse("0.1.0"),
	"Deploys a CCIPFactory contract on Canton",
	func(b operations.Bundle, deps canton.Chain, input DeployCCIPFactoryParams) (ccipsequences.OnChainOutput, error) {
		if input.OwnerParty == "" {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("owner party is required")
		}

		ownerParty := types.PARTY(input.OwnerParty)
		mcmsParty := ownerParty
		if input.MCMSParty != "" {
			mcmsParty = types.PARTY(input.MCMSParty)
		}

		var qualifier *string
		if input.Qualifier != "" {
			qualifier = &input.Qualifier
		}

		deployReport, err := operations.ExecuteOperation(b, factoryops.Deploy, deps, opcontract.DeployInput[factorybindings.CCIPFactory]{
			Qualifier: qualifier,
			Template: factorybindings.CCIPFactory{
				InstanceId:                    types.TEXT(input.InstanceID),
				Owner:                         ownerParty,
				McmsParty:                     mcmsParty,
				UsedInstanceIds:               map[types.TEXT]types.BOOL{},
				DeployedContracts:             map[types.TEXT]types.CONTRACT_ID{},
				PerPartyRouterFactoryDeployed: types.BOOL(false),
			},
			OwnerParty: ownerParty,
		})
		if err != nil {
			return ccipsequences.OnChainOutput{}, fmt.Errorf("deploy CCIPFactory: %w", err)
		}

		return ccipsequences.OnChainOutput{Addresses: []datastore.AddressRef{deployReport.Output}}, nil
	},
)

// DeployFactoryAndSetOwnerToMCMSConfig holds parameters for the combined changeset
// that deploys a CCIPFactory and generates an MCMS proposal to transfer ownership.
type DeployFactoryAndSetOwnerToMCMSConfig struct {
	OwnerParty string        `json:"ownerParty" yaml:"ownerParty"`
	MCMSParty  string        `json:"mcmsParty" yaml:"mcmsParty"`
	InstanceID string        `json:"instanceID,omitempty" yaml:"instanceID,omitempty"`
	Qualifier  string        `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
	MinDelay   time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
}

type DeployFactoryAndSetOwnerToMCMS struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployFactoryAndSetOwnerToMCMSConfig]] = DeployFactoryAndSetOwnerToMCMS{}

func (d DeployFactoryAndSetOwnerToMCMS) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployFactoryAndSetOwnerToMCMSConfig]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants",
			config.Participant, config.ChainSelector, len(chain.Participants))
	}
	if config.Config.OwnerParty == "" {
		return fmt.Errorf("owner party is required")
	}
	if config.Config.MCMSParty == "" {
		return fmt.Errorf("mcms party is required")
	}

	// Verify MCMS contract exists in the datastore (needed for the proposal).
	if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(mcmsops.ContractType),
		mcmsops.Version,
		"",
	)); err != nil {
		return fmt.Errorf("MCMS contract not found in datastore (deploy MCMS first): %w", err)
	}

	return nil
}

func (d DeployFactoryAndSetOwnerToMCMS) Apply(e cldf.Environment, config CantonCSDeps[DeployFactoryAndSetOwnerToMCMSConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	// ── Step 1: Deploy CCIPFactory ──────────────────────────────────────────
	deployOut, err := operations.ExecuteSequence(e.OperationsBundle, deployCCIPFactorySequence, chain, DeployCCIPFactoryParams{
		OwnerParty: config.Config.OwnerParty,
		MCMSParty:  config.Config.MCMSParty,
		InstanceID: config.Config.InstanceID,
		Qualifier:  config.Config.Qualifier,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy CCIPFactory: %w", err)
	}

	for _, addrRef := range deployOut.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store address ref: %w", err)
		}
	}

	// Get the factory raw instance address from the deploy output.
	if len(deployOut.Output.Addresses) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("no address refs returned from factory deploy")
	}
	factoryRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(deployOut.Output.Addresses[0])
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get factory raw instance address: %w", err)
	}

	// ── Step 2: Encode SetOwnerToMCMS as MCMS proposal ──────────────────────
	exerciseOut, err := operations.ExecuteOperation(e.OperationsBundle, factoryops.SetOwnerToMCMS, chain, opcontract.ChoiceInput[factorybindings.SetOwnerToMCMS]{
		InstanceAddress:    factoryRawAddr.InstanceAddress(),
		RawInstanceAddress: string(factoryRawAddr),
		Args:               factorybindings.SetOwnerToMCMS{},
		MCMSEnabled:        true,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to encode SetOwnerToMCMS: %w", err)
	}

	batchOp, err := cantonmcms.BuildBatchFromOutputs([]opcontract.ExerciseOutput{exerciseOut.Output})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build batch operation: %w", err)
	}

	// Look up MCMS contract from the datastore for the proposal.
	mcmsRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(mcmsops.ContractType),
		mcmsops.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS contract not found in datastore: %w", err)
	}
	mcmsRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(mcmsRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS raw instance address: %w", err)
	}

	participant := chain.Participants[config.Participant]

	proposal, err := cantonmcms.GenerateTimelockProposal(
		e.GetContext(),
		participant.LedgerServices.State,
		participant.PartyID,
		cantonmcms.ProposalConfig{
			MCMSContract: cantonmcms.MCMSContractInfo{
				RawInstanceAddress: mcmsRawAddr,
				InstanceAddress:    mcmsRawAddr.InstanceAddress(),
			},
			ChainSelector: mcms_types.ChainSelector(config.ChainSelector),
			Description:   "Transfer CCIPFactory ownership to MCMS party",
			MinDelay:      config.Config.MinDelay,
			Action:        mcms_types.TimelockActionSchedule,
			Role:          cantonsdk.TimelockRoleProposer,
		},
		[]mcms_types.BatchOperation{batchOp},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate SetOwnerToMCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore:             ds,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
