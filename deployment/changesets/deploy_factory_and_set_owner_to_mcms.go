package changesets

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// DeployFactoryAndSetOwnerToMCMS deploys a CCIPFactory and returns an MCMS proposal to transfer ownership.
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
	if config.Config.Qualifier == "" {
		return fmt.Errorf("factory qualifier is required")
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

func (d DeployFactoryAndSetOwnerToMCMS) Apply(e cldf.Environment, config CantonCSDeps[DeployFactoryAndSetOwnerToMCMSConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	cfg := config.Config

	factoryRef, err := deployCCIPFactory(e.OperationsBundle, chain, DeployCCIPFactoryParams{
		OwnerParty: cfg.OwnerParty,
		MCMSParty:  cfg.MCMSParty,
		InstanceID: cfg.InstanceID,
		Qualifier:  cfg.Qualifier,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if err := ds.AddressRefStore.Add(factoryRef); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("store address ref: %w", err)
	}

	factoryRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(factoryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("factory raw instance address: %w", err)
	}

	// Note: default MCMS
	exerciseOut, err := operations.ExecuteOperation(e.OperationsBundle, factoryops.SetOwnerToMCMS, chain, opcontract.ChoiceInput[factorybindings.SetOwnerToMCMS]{
		InstanceAddress:    factoryRawAddr.InstanceAddress(),
		RawInstanceAddress: factoryRawAddr.String(),
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

	mcmsRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(mcmsops.ContractType),
		mcmsops.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS contract not found: %w", err)
	}
	mcmsRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(mcmsRef)
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
				RawInstanceAddress: mcmsRawAddr,
				InstanceAddress:    mcmsRawAddr.InstanceAddress(),
			},
			ChainSelector: mcms_types.ChainSelector(config.ChainSelector),
			Description:   fmt.Sprintf("Transfer CCIPFactory (%s) ownership to MCMS", cfg.Qualifier),
			MinDelay:      cfg.MinDelay,
			Action:        mcms_types.TimelockActionSchedule,
			Role:          cantonsdk.TimelockRoleProposer,
		},
		[]mcms_types.BatchOperation{batchOp},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("generate SetOwnerToMCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore:             ds,
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
