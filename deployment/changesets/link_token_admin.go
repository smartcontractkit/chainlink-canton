package changesets

import (
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/linkadmingov"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

// --- DeployLinkTokenAdminFromFactory ---

type DeployLinkTokenAdminFromFactoryConfig struct {
	CCIPOwnerParty string        `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	MinDelay       time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	Params         sequences.DeployLinkTokenAdminFromFactoryParams
}

type DeployLinkTokenAdminFromFactory struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[DeployLinkTokenAdminFromFactoryConfig]] = DeployLinkTokenAdminFromFactory{}

func (d DeployLinkTokenAdminFromFactory) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[DeployLinkTokenAdminFromFactoryConfig]) error {
	if err := requireCantonParticipant(e, config.ChainSelector, config.Participant); err != nil {
		return err
	}
	if config.Config.CCIPOwnerParty == "" {
		return fmt.Errorf("ccipOwnerParty is required")
	}
	if _, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCIP); err != nil {
		return fmt.Errorf("core CCIPFactory must be deployed first: %w", err)
	}
	if _, err := linkTokenAdminRef(e, config.ChainSelector); err == nil {
		return fmt.Errorf("LinkTokenAdmin is already deployed for chain %d", config.ChainSelector)
	}

	return nil
}

func (d DeployLinkTokenAdminFromFactory) Apply(e cldf.Environment, config CantonCSDeps[DeployLinkTokenAdminFromFactoryConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	factoryRef, err := dsutils.FactoryAddressRef(e.DataStore, config.ChainSelector, dsutils.QualifierCCIP)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	params := config.Config.Params
	params.CCIPOwnerParty = config.Config.CCIPOwnerParty
	params.FactoryAddressRef = factoryRef
	params.ProposalDriven = len(participant.ReadAsPartyIDs) > 0

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.DeployLinkTokenAdminFromFactory, chain, params)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("deploy link token admin from factory: %w", err)
	}

	ds := datastore.NewMemoryDataStore()
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

// --- GovernLinkTokenAdmin ---

// GovernLinkTokenAdminConfig lists LinkTokenAdmin calls for one ccipOwner MCMS proposal, in order.
type GovernLinkTokenAdminConfig struct {
	Actions     []sequences.LinkTokenAdminAction `json:"actions" yaml:"actions"`
	MinDelay    time.Duration                    `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description string                           `json:"description,omitempty" yaml:"description,omitempty"`
}

type GovernLinkTokenAdmin struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[GovernLinkTokenAdminConfig]] = GovernLinkTokenAdmin{}

func (g GovernLinkTokenAdmin) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[GovernLinkTokenAdminConfig]) error {
	if err := requireCantonParticipant(e, config.ChainSelector, config.Participant); err != nil {
		return err
	}
	input, err := governLinkTokenAdminInput(e, config)
	if err != nil {
		return err
	}
	if err := input.Validate(); err != nil {
		return err
	}
	if err := requireLiveAuthorizations(e, config, input); err != nil {
		return err
	}
	if input.ProposalDriven {
		return requireCCIPOwnerMCMSRef(e, config.ChainSelector)
	}

	return nil
}

// requireLiveAuthorizations checks the ledger: an approval must not reuse a live id, and a withdraw must have a live target.
func requireLiveAuthorizations(e cldf.Environment, config CantonCSDeps[GovernLinkTokenAdminConfig], input sequences.GovernLinkTokenAdminInput) error {
	participant := e.BlockChains.CantonChains()[config.ChainSelector].Participants[config.Participant]
	owner := types.PARTY(input.AdminAddress.Owner())
	withdrawn := make(map[sequences.LinkAuthorizationRef]bool)
	for _, action := range input.Actions {
		ref, approving := action.ApprovalRef()
		if !approving {
			if action.Withdraw == nil {
				continue
			}
			ref = *action.Withdraw
		}
		// A withdraw earlier in the batch archives the live one before this approval creates its successor.
		if approving && withdrawn[ref] {
			continue
		}
		templateID, ok := linkadmingov.AuthorizationTemplateID(ref.Kind)
		if !ok {
			return fmt.Errorf("unknown authorization kind %q", ref.Kind)
		}
		addr := contracts.InstanceID(ref.AuthorizationID).RawInstanceAddress(owner)
		_, err := opcontract.FindActiveContractByInstanceAddress(
			e.GetContext(), participant.LedgerServices.State, opcontract.LedgerQueryParties(participant), templateID, addr.InstanceAddress(),
		)
		switch {
		case approving && err == nil:
			return fmt.Errorf("authorization %q of kind %s is already live; withdraw it first", ref.AuthorizationID, ref.Kind)
		case approving && strings.Contains(err.Error(), "no active contract found"):
			continue
		case !approving && err != nil && strings.Contains(err.Error(), "no active contract found"):
			return fmt.Errorf("authorization %q of kind %s is not live", ref.AuthorizationID, ref.Kind)
		case err != nil:
			return err
		}
		withdrawn[ref] = true
	}

	return nil
}

func (g GovernLinkTokenAdmin) Apply(e cldf.Environment, config CantonCSDeps[GovernLinkTokenAdminConfig]) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	input, err := governLinkTokenAdminInput(e, config)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	out, err := operations.ExecuteSequence(e.OperationsBundle, sequences.GovernLinkTokenAdmin, chain, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("govern link token admin: %w", err)
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, input.ProposalDriven,
		cantonmcms.QualifierCCIPOwner, config.Config.MinDelay, config.Config.Description, nil, out.Output.BatchOps,
	)
}

func governLinkTokenAdminInput(e cldf.Environment, config CantonCSDeps[GovernLinkTokenAdminConfig]) (sequences.GovernLinkTokenAdminInput, error) {
	adminRef, err := linkTokenAdminRef(e, config.ChainSelector)
	if err != nil {
		return sequences.GovernLinkTokenAdminInput{}, fmt.Errorf("LinkTokenAdmin must be deployed first: %w", err)
	}
	adminRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(adminRef)
	if err != nil {
		return sequences.GovernLinkTokenAdminInput{}, fmt.Errorf("resolve LinkTokenAdmin raw address: %w", err)
	}
	participant := e.BlockChains.CantonChains()[config.ChainSelector].Participants[config.Participant]

	return sequences.GovernLinkTokenAdminInput{
		AdminAddress:   adminRaw,
		ProposalDriven: len(participant.ReadAsPartyIDs) > 0,
		Actions:        config.Config.Actions,
	}, nil
}

func linkTokenAdminRef(e cldf.Environment, chainSelector uint64) (datastore.AddressRef, error) {
	return e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(linkadmingov.ContractType),
		linkadmingov.Version,
		"",
	))
}

func requireCantonParticipant(e cldf.Environment, chainSelector uint64, participant int) error {
	chain, ok := e.BlockChains.CantonChains()[chainSelector]
	if !ok {
		return fmt.Errorf("canton chain %v not found", chainSelector)
	}
	if participant < 0 || participant >= len(chain.Participants) {
		return fmt.Errorf("participant index %d out of range for canton chain %d with %d participants", participant, chainSelector, len(chain.Participants))
	}

	return nil
}
