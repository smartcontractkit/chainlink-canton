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
	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// --- DeployCoreChainContractsFromFactory ---

type DeployCoreChainContractsFromFactoryConfig struct {
	OwnerParty     string `json:"ownerParty" yaml:"ownerParty"`
	CCIPOwnerParty string `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	Params         sequences.DeployChainContractsParams
}

type DeployCoreChainContractsFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCoreChainContractsFromFactoryConfig]] = DeployCoreChainContractsFromFactory{}

func (d DeployCoreChainContractsFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCoreChainContractsFromFactoryConfig]) error {
	if config.Config.OwnerParty == "" && config.Config.CCIPOwnerParty == "" {
		return fmt.Errorf("ownerParty or ccipOwnerParty is required")
	}
	_, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCore)
	if err != nil {
		return fmt.Errorf("core CCIPFactory must be deployed first: %w", err)
	}

	return nil
}

func (d DeployCoreChainContractsFromFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployCoreChainContractsFromFactoryConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCore)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	params := config.Config.Params
	params.OwnerParty = config.Config.OwnerParty
	params.CCIPOwnerParty = config.Config.CCIPOwnerParty
	params.FactoryAddressRef = factoryRef
	params.CommitteeVerifiers = nil
	params.ProposalDriven = len(participant.ReadAsPartyIDs) > 0

	if params.CcvRegistryBinding.Unpack == "" {
		ccvRaw, err := dsutils.FirstCommitteeVerifierRawAddress(e.DataStore, config.ChainSelector, "")
		if err == nil {
			params.CcvRegistryBinding = ccvRaw.Binding()
		}
	}

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployCoreChainContractsFromFactory, chain, params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy core chain contracts from factory: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("store address ref %v: %w", addrRef, err)
		}
	}

	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}

// --- DeployCCVFromFactory ---

type DeployCCVFromFactoryConfig struct {
	OwnerParty     string `json:"ownerParty" yaml:"ownerParty"`
	CCIPOwnerParty string `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	Params         sequences.DeployChainContractsParams
}

type DeployCCVFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployCCVFromFactoryConfig]] = DeployCCVFromFactory{}

func (d DeployCCVFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployCCVFromFactoryConfig]) error {
	if config.Config.OwnerParty == "" && config.Config.CCIPOwnerParty == "" {
		return fmt.Errorf("ownerParty or ccipOwnerParty is required")
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

	return cldf.ChangesetOutput{DataStore: ds, Reports: []operations.Report[any, any]{}}, nil
}

// --- SetFactoryOwnerToMCMS ---

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
	_, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, config.Config.FactoryQualifier)
	if err != nil {
		return fmt.Errorf("factory not found: %w", err)
	}
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
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("encode SetOwnerToMCMS: %w", err)
	}

	batchOp, err := cantonmcms.BuildBatchFromOutputs([]opcontract.ExerciseOutput{exerciseOut.Output})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("build batch operation: %w", err)
	}

	mcmsRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(mcmsops.ContractType),
		mcmsops.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS contract not found: %w", err)
	}
	mcmsRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(mcmsRef)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	participant := chain.Participants[config.Participant]
	proposal, err := cantonmcms.GenerateTimelockProposal(
		e.GetContext(),
		participant.LedgerServices.State,
		participant.PartyID,
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
	mcmsParty := ownerParty
	if params.MCMSParty != "" {
		mcmsParty = types.PARTY(params.MCMSParty)
	}

	var qualifier *string
	if params.Qualifier != "" {
		qualifier = &params.Qualifier
	}

	deployReport, err := operations.ExecuteOperation(b, factoryops.Deploy, chain, opcontract.DeployInput[factorybindings.CCIPFactory]{
		Qualifier: qualifier,
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
