package changesets

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/ccip/factory"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// --- DeployRMNFromFactory ---

type DeployRMNFromFactoryConfig struct {
	OwnerParty     string        `json:"ownerParty" yaml:"ownerParty"`
	CCIPOwnerParty string        `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	RMNOwnerParty  string        `json:"rmnOwnerParty" yaml:"rmnOwnerParty"`
	MinDelay       time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	Params         sequences.DeployChainContractsParams
}

type DeployRMNFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployRMNFromFactoryConfig]] = DeployRMNFromFactory{}

func (d DeployRMNFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployRMNFromFactoryConfig]) error {
	if err := requireFactoryDeployOwnerParties(config.Config.OwnerParty, config.Config.CCIPOwnerParty); err != nil {
		return err
	}
	if config.Config.RMNOwnerParty == "" {
		return fmt.Errorf("rmnOwnerParty is required")
	}
	_, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierRMN)
	if err != nil {
		return fmt.Errorf("rmn CCIPFactory must be deployed first: %w", err)
	}
	_, err = dsutils.RmnRemoteRawAddress(e.DataStore, config.ChainSelector)
	if err == nil {
		return fmt.Errorf("RMNRemote is already deployed for chain %d", config.ChainSelector)
	}

	return nil
}

func (d DeployRMNFromFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployRMNFromFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierRMN)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	params := config.Config.Params
	params.OwnerParty = config.Config.OwnerParty
	params.CCIPOwnerParty = config.Config.CCIPOwnerParty
	params.RMNOwnerParty = config.Config.RMNOwnerParty
	params.FactoryAddressRef = factoryRef
	params.ProposalDriven = len(participant.ReadAsPartyIDs) > 0

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployRMNFromFactory, chain, params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy RMN from factory: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("store address ref %v: %w", addrRef, err)
		}
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, params.ProposalDriven,
		cantonmcms.QualifierRMNOwner, config.Config.MinDelay, config.Config.Description, ds, out.Output.BatchOps,
	)
}

// --- DeployCCIPChainContractsFromFactory ---

type DeployCCIPChainContractsFromFactoryConfig struct {
	OwnerParty     string        `json:"ownerParty" yaml:"ownerParty"`
	CCIPOwnerParty string        `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	MinDelay       time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	Params         sequences.DeployChainContractsParams
}

type DeployCCIPChainContractsFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCCIPChainContractsFromFactoryConfig]] = DeployCCIPChainContractsFromFactory{}

func (d DeployCCIPChainContractsFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCCIPChainContractsFromFactoryConfig]) error {
	if err := requireFactoryDeployOwnerParties(config.Config.OwnerParty, config.Config.CCIPOwnerParty); err != nil {
		return err
	}
	if config.Config.Params.CcvRegistryBinding.Unpack == "" {
		return fmt.Errorf("CcvRegistryBinding is required for core CCIP factory deploy")
	}
	_, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCIP)
	if err != nil {
		return fmt.Errorf("core CCIPFactory must be deployed first: %w", err)
	}
	_, err = dsutils.RmnRemoteRawAddress(e.DataStore, config.ChainSelector)
	if err != nil {
		return fmt.Errorf("RMNRemote must be deployed before core CCIP contracts: %w", err)
	}

	return nil
}

func (d DeployCCIPChainContractsFromFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCCIPChainContractsFromFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCIP)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	rmnRaw, err := dsutils.RmnRemoteRawAddress(e.DataStore, config.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	params := config.Config.Params
	params.OwnerParty = config.Config.OwnerParty
	params.CCIPOwnerParty = config.Config.CCIPOwnerParty
	params.FactoryAddressRef = factoryRef
	params.RmnRemoteRawInstanceAddress = rmnRaw
	params.CommitteeVerifiers = nil
	params.ProposalDriven = len(participant.ReadAsPartyIDs) > 0

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployCCIPChainContractsFromFactory, chain, params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy core chain contracts from factory: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("store address ref %v: %w", addrRef, err)
		}
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, params.ProposalDriven,
		cantonmcms.QualifierCCIPOwner, config.Config.MinDelay, config.Config.Description, ds, out.Output.BatchOps,
	)
}

// --- DeployCCVFromFactory ---

type DeployCCVFromFactoryConfig struct {
	OwnerParty     string        `json:"ownerParty" yaml:"ownerParty"`
	CCIPOwnerParty string        `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	CCVOwnerParty  string        `json:"ccvOwnerParty" yaml:"ccvOwnerParty"`
	MinDelay       time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	Params         sequences.DeployChainContractsParams
}

type DeployCCVFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCCVFromFactoryConfig]] = DeployCCVFromFactory{}

func (d DeployCCVFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCCVFromFactoryConfig]) error {
	if err := requireFactoryDeployOwnerParties(config.Config.OwnerParty, config.Config.CCIPOwnerParty); err != nil {
		return err
	}
	if config.Config.CCVOwnerParty == "" {
		return fmt.Errorf("ccvOwnerParty is required")
	}
	if len(config.Config.Params.CommitteeVerifiers) == 0 {
		return fmt.Errorf("at least one committee verifier is required in params")
	}
	_, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCV)
	if err != nil {
		return fmt.Errorf("ccv CCIPFactory must be deployed first: %w", err)
	}
	_, err = dsutils.RmnRemoteRawAddress(e.DataStore, config.ChainSelector)
	if err != nil {
		return fmt.Errorf("RMNRemote must be deployed before CCV: %w", err)
	}

	return nil
}

func (d DeployCCVFromFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCCVFromFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCV)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	rmnRaw, err := dsutils.RmnRemoteRawAddress(e.DataStore, config.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	params := config.Config.Params
	params.OwnerParty = config.Config.OwnerParty
	params.CCIPOwnerParty = config.Config.CCIPOwnerParty
	params.CCVOwnerParty = config.Config.CCVOwnerParty
	params.FactoryAddressRef = factoryRef
	params.RmnRemoteRawInstanceAddress = rmnRaw
	params.ProposalDriven = len(participant.ReadAsPartyIDs) > 0

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployCCVFromFactory, chain, params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy CCV from factory: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("store address ref %v: %w", addrRef, err)
		}
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, params.ProposalDriven,
		cantonmcms.QualifierCCVOwner, config.Config.MinDelay, config.Config.Description, ds, out.Output.BatchOps,
	)
}

// --- SetFactoryOwnerToMCMS ---
// Encodes an MCMS proposal for SetOwnerToMCMS (controller = mcmsParty only).

type SetFactoryOwnerToMCMSConfig struct {
	FactoryQualifier string        `json:"factoryQualifier" yaml:"factoryQualifier"`
	MCMSParty        string        `json:"mcmsParty" yaml:"mcmsParty"`
	MinDelay         time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
}

type SetFactoryOwnerToMCMS struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[SetFactoryOwnerToMCMSConfig]] = SetFactoryOwnerToMCMS{}

func (s SetFactoryOwnerToMCMS) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[SetFactoryOwnerToMCMSConfig]) error {
	if config.Config.FactoryQualifier == "" {
		return fmt.Errorf("factoryQualifier is required")
	}
	if config.Config.MCMSParty == "" {
		return fmt.Errorf("mcmsParty is required")
	}
	if _, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, config.Config.FactoryQualifier); err != nil {
		return fmt.Errorf("factory not found: %w", err)
	}
	mcmsOwnerQualifier, err := mcmsOwnerQualifierForFactory(config.Config.FactoryQualifier)
	if err != nil {
		return err
	}
	if _, err := dsutils.MCMSRawInstanceAddress(e.DataStore, config.ChainSelector, mcmsOwnerQualifier); err != nil {
		return fmt.Errorf("MCMS contract not found in datastore (deploy MCMS first): %w", err)
	}

	return nil
}

func (s SetFactoryOwnerToMCMS) Apply(e cldf.Environment, config CantonCSDeps[SetFactoryOwnerToMCMSConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, config.Config.FactoryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	factoryRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(factoryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	exerciseOut, err := operations.ExecuteOperation(e.OperationsBundle, factoryops.SetOwnerToMCMS, chain, opcontract.ChoiceInput[factorybindings.SetOwnerToMCMS]{
		InstanceAddress:    factoryRaw.InstanceAddress(),
		RawInstanceAddress: factoryRaw.String(),
		Args:               factorybindings.SetOwnerToMCMS{},
		MCMSEnabled:        true,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("encode SetOwnerToMCMS: %w", err)
	}

	batchOp, err := cantonmcms.BuildBatchFromOutputs([]opcontract.ExerciseOutput{exerciseOut.Output})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("build batch operation: %w", err)
	}

	mcmsOwnerQualifier, err := mcmsOwnerQualifierForFactory(config.Config.FactoryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	mcmsRaw, err := dsutils.MCMSRawInstanceAddress(e.DataStore, config.ChainSelector, mcmsOwnerQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS contract not found: %w", err)
	}

	participant := chain.Participants[config.Participant]
	proposal, err := cantonmcms.GenerateTimelockProposal(
		e.GetContext(),
		participant.LedgerServices.State,
		opcontract.LedgerQueryParties(participant),
		cantonmcms.ProposalConfig{
			MCMSContract: cantonmcms.MCMSContractInfo{
				RawInstanceAddress: mcmsRaw,
				InstanceAddress:    mcmsRaw.InstanceAddress(),
			},
			ChainSelector: mcms_types.ChainSelector(config.ChainSelector),
			Description:   fmt.Sprintf("Transfer CCIPFactory (%s) ownership to MCMS", config.Config.FactoryQualifier),
			MinDelay:      config.Config.MinDelay,
			Action:        mcms_types.TimelockActionSchedule,
			Role:          cantonsdk.TimelockRoleProposer,
		},
		[]mcms_types.BatchOperation{batchOp},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("generate SetOwnerToMCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}

// --- DeployCCIPFactory (direct deploy, no MCMS handoff; e.g. devenv pre-deploy) ---

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
	if config.Config.Params.MCMSParty == "" {
		return fmt.Errorf("mcmsParty is required")
	}

	return nil
}

func (d DeployCCIPFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCCIPFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	addrRef, err := deployCCIPFactory(e.OperationsBundle, chain, config.Config.Params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy CCIPFactory: %w", err)
	}
	if err := ds.AddressRefStore.Add(addrRef); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("store address ref: %w", err)
	}

	return cldf.ChangesetOutput{DataStore: ds}, nil
}

func deployCCIPFactory(b operations.Bundle, chain canton.Chain, params DeployCCIPFactoryParams) (datastore.AddressRef, error) {
	ownerParty := types.PARTY(params.OwnerParty)
	mcmsParty := types.PARTY(params.MCMSParty)
	qualifier := params.Qualifier

	deployReport, err := operations.ExecuteOperation(b, factoryops.Deploy, chain, opcontract.DeployInput[factorybindings.CCIPFactory]{
		Qualifier: &qualifier,
		Template: factorybindings.CCIPFactory{
			InstanceId:                    types.TEXT(params.InstanceID),
			Owner:                         ownerParty,
			McmsParty:                     mcmsParty,
			UsedInstanceIds:               map[types.TEXT]types.BOOL{},
			DeployedContracts:             map[types.TEXT]types.CONTRACT_ID{},
			PerPartyRouterFactoryDeployed: types.BOOL(false),
		},
		OwnerParty: ownerParty,
	})
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("deploy CCIPFactory: %w", err)
	}

	return deployReport.Output, nil
}

func buildFactoryDeployChangesetOutput(
	e cldf.Environment,
	chain canton.Chain,
	chainSelector uint64,
	participantIdx int,
	proposalDriven bool,
	mcmsOwnerQualifier string,
	minDelay time.Duration,
	description string,
	ds *datastore.MemoryDataStore,
	batchOps []mcms_types.BatchOperation,
) (cldf.ChangesetOutput, error) {
	output := cldf.ChangesetOutput{
		DataStore: ds,
		Reports:   []operations.Report[any, any]{},
	}
	if !proposalDriven {
		return output, nil
	}
	if description == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("description is required for proposal-driven factory deploy")
	}
	if len(batchOps) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("proposal-driven factory deploy produced no batch operations")
	}

	mcmsRaw, err := dsutils.MCMSRawInstanceAddress(e.DataStore, chainSelector, mcmsOwnerQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve MCMS for qualifier %q: %w", mcmsOwnerQualifier, err)
	}

	participant := chain.Participants[participantIdx]
	proposal, err := cantonmcms.GenerateTimelockProposal(
		e.GetContext(),
		participant.LedgerServices.State,
		opcontract.LedgerQueryParties(participant),
		cantonmcms.ProposalConfig{
			MCMSContract: cantonmcms.MCMSContractInfo{
				RawInstanceAddress: mcmsRaw,
				InstanceAddress:    mcmsRaw.InstanceAddress(),
			},
			ChainSelector: mcms_types.ChainSelector(chainSelector),
			Description:   description,
			MinDelay:      minDelay,
			Action:        mcms_types.TimelockActionSchedule,
			Role:          cantonsdk.TimelockRoleProposer,
		},
		batchOps,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("generate factory deploy proposal: %w", err)
	}

	output.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}

	return output, nil
}

func requireFactoryDeployOwnerParties(ownerParty, ccipOwnerParty string) error {
	if ownerParty == "" {
		return fmt.Errorf("ownerParty is required")
	}
	if ccipOwnerParty == "" {
		return fmt.Errorf("ccipOwnerParty is required")
	}

	return nil
}

func mcmsOwnerQualifierForFactory(factoryQualifier string) (string, error) {
	switch factoryQualifier {
	case dsutils.QualifierCCIP:
		return cantonmcms.QualifierCCIPOwner, nil
	case dsutils.QualifierCCV:
		return cantonmcms.QualifierCCVOwner, nil
	case dsutils.QualifierRMN:
		return cantonmcms.QualifierRMNOwner, nil
	default:
		return "", fmt.Errorf("unsupported factory qualifier %q for MCMS owner lookup (expected %q, %q, or %q)", factoryQualifier, dsutils.QualifierCCIP, dsutils.QualifierCCV, dsutils.QualifierRMN)
	}
}
