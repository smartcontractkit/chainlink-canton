package changesets

import (
	"fmt"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-canton/bindings/generated/latest/splice/splice_api_token_holding_v1"
	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-canton/deployment/sequences"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/nativeinstrument"
	cantonmcms "github.com/smartcontractkit/chainlink-canton/deployment/utils/mcms"
)

// RegisterNativeFeeTokenInTARConfig registers the Canton native fee token (Amulet) in TAR.
// Creates TokenConfig only — no token pool. InstrumentId is resolved from the validator
// scan-proxy when admin/id are empty (requires ledger access at pipeline run time).
type RegisterNativeFeeTokenInTARConfig struct {
	CCIPOwnerParty string        `json:"ccipOwnerParty" yaml:"ccipOwnerParty"`
	MinDelay       time.Duration `json:"minDelay,omitempty" yaml:"minDelay,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	InstrumentId   splice_api_token_holding_v1.InstrumentId `json:"instrumentId,omitempty" yaml:"instrumentId,omitempty"`
	TokenQualifier string        `json:"tokenQualifier,omitempty" yaml:"tokenQualifier,omitempty"`
}

type RegisterNativeFeeTokenInTAR struct{}

var _ cldf.ChangeSetV2[CantonCSDeps[RegisterNativeFeeTokenInTARConfig]] = RegisterNativeFeeTokenInTAR{}

func (r RegisterNativeFeeTokenInTAR) VerifyPreconditions(e cldf.Environment, config CantonCSDeps[RegisterNativeFeeTokenInTARConfig]) error {
	if config.Config.CCIPOwnerParty == "" {
		return fmt.Errorf("ccipOwnerParty is required")
	}
	if err := requireCCIPOwnerMCMSRef(e, config.ChainSelector); err != nil {
		return err
	}
	if _, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	)); err != nil {
		return fmt.Errorf("token admin registry must be deployed first: %w", err)
	}

	return nil
}

func (r RegisterNativeFeeTokenInTAR) Apply(e cldf.Environment, config CantonCSDeps[RegisterNativeFeeTokenInTARConfig]) (cldf.ChangesetOutput, error) {
	ds := datastore.NewMemoryDataStore()
	chain := e.BlockChains.CantonChains()[config.ChainSelector]
	participant := chain.Participants[config.Participant]

	tarRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(
		config.ChainSelector,
		datastore.ContractType(token_admin_registry.ContractType),
		token_admin_registry.Version,
		"",
	))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve token admin registry: %w", err)
	}
	tarRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(tarRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("resolve token admin registry raw address: %w", err)
	}

	instrumentID := config.Config.InstrumentId
	if instrumentID.Admin == "" || instrumentID.Id == "" {
		instrumentID, err = nativeinstrument.LookupNativeInstrumentID(e.GetContext(), participant)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("resolve native fee token instrument: %w", err)
		}
	}

	proposalDriven := len(participant.ReadAsPartyIDs) > 0
	out, err := cld_ops.ExecuteSequence(e.OperationsBundle, sequences.RegisterNativeFeeTokenInTAR, chain, sequences.RegisterNativeFeeTokenInTARInput{
		TokenAdminRegistryInstanceAddress:    contracts.HexToInstanceAddress(tarRef.Address),
		TokenAdminRegistryRawInstanceAddress: tarRaw,
		InstrumentId:                         instrumentID,
		CcipOwnerParty:                       config.Config.CCIPOwnerParty,
		TokenQualifier:                       config.Config.TokenQualifier,
		ChainSelector:                        config.ChainSelector,
		CcipParticipantIndex:                 config.Participant,
		ProposalDriven:                       proposalDriven,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("register native fee token in TAR: %w", err)
	}

	for _, addrRef := range out.Output.Addresses {
		if err := ds.AddressRefStore.Add(addrRef); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("store address ref %v: %w", addrRef, err)
		}
	}

	if proposalDriven && len(out.Output.BatchOps) == 0 {
		e.Logger.Infof("native fee token %s already registered in TAR; no MCMS proposal generated", instrumentID.Id)
		return cldf.ChangesetOutput{DataStore: ds}, nil
	}

	return buildFactoryDeployChangesetOutput(
		e, chain, config.ChainSelector, config.Participant, proposalDriven,
		cantonmcms.QualifierCCIPOwner, config.Config.MinDelay, config.Config.Description, ds, out.Output.BatchOps,
	)
}
