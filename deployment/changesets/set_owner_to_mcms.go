package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	factorybindings "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/factory"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	factoryops "github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/factory"
	mcmsops "github.com/smartcontractkit/chainlink-canton/deployment/operations/mcms"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

type SetOwnerToMCMS struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[struct{}]] = SetOwnerToMCMS{}

func (s SetOwnerToMCMS) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[struct{}]) error {
	chain, ok := e.BlockChains.CantonChains()[config.ChainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", config.ChainSelector)
	}
	if config.Participant < 0 || config.Participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants",
			config.Participant, config.ChainSelector, len(chain.Participants))
	}

	for _, key := range []datastore.AddressRefKey{
		datastore.NewAddressRefKey(config.ChainSelector, datastore.ContractType(factoryops.ContractType), factoryops.Version, ""),
		datastore.NewAddressRefKey(config.ChainSelector, datastore.ContractType(mcmsops.ContractType), mcmsops.Version, ""),
	} {
		if _, err := e.DataStore.Addresses().Get(key); err != nil {
			return fmt.Errorf("contract not found in datastore (%s): %w", key.Type, err)
		}
	}

	return nil
}

func (s SetOwnerToMCMS) Apply(e cldf.Environment, config CantonCSDeps[struct{}]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]

	// Look up factory address from the datastore.
	factoryRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(factoryops.ContractType),
		factoryops.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("CCIPFactory not found in datastore: %w", err)
	}

	factoryRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(factoryRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get factory raw instance address: %w", err)
	}

	// Look up MCMS address from the datastore.
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

	// Encode the SetOwnerToMCMS exercise for the MCMS proposal.
	out, err := operations.ExecuteOperation(e.OperationsBundle, factoryops.SetOwnerToMCMS, chain, opcontract.ChoiceInput[factorybindings.SetOwnerToMCMS]{
		InstanceAddress:    factoryRawAddr.InstanceAddress(),
		RawInstanceAddress: string(factoryRawAddr),
		Args:               factorybindings.SetOwnerToMCMS{},
		MCMSEnabled:        true,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to encode SetOwnerToMCMS: %w", err)
	}

	batchOp, err := cantonmcms.BuildBatchFromOutputs([]opcontract.ExerciseOutput{out.Output})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build batch operation: %w", err)
	}
	if len(batchOp.Transactions) == 0 {
		return cldf.ChangesetOutput{}, nil
	}

	participant := chain.Participants[config.Participant]

	proposal, err := cantonmcms.GenerateTimelockProposal(
		e.GetContext(),
		participant.LedgerServices.State,
		participant.PartyID,
		cantonmcms.ProposalConfig{
			MCMSContract: cantonmcms.MCMSContractInfo{
				RawInstanceAddress: contracts.RawInstanceAddress(mcmsRawAddr),
				InstanceAddress:    mcmsRawAddr.InstanceAddress(),
			},
			ChainSelector: mcms_types.ChainSelector(config.ChainSelector),
			Description:   "Transfer CCIPFactory ownership to MCMS party",
			Action:        mcms_types.TimelockActionBypass,
			Role:          cantonsdk.TimelockRoleBypasser,
		},
		[]mcms_types.BatchOperation{batchOp},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate timelock proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
