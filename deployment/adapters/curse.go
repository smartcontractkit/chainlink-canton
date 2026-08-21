package adapters

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/go-daml/pkg/types"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings"
	"github.com/smartcontractkit/chainlink-canton/contracts/bindings/generated/ccip/core"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/global_config"
	"github.com/smartcontractkit/chainlink-canton/deployment/operations/ccip/rmn_remote"
	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
	internalparse "github.com/smartcontractkit/chainlink-canton/internal/parse"
)

var (
	_ fastcurse.CurseAdapter        = (*CantonCurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CantonCurseSubjectAdapter)(nil)
)

func NewCantonCurseAdapter() *CantonCurseAdapter {
	return &CantonCurseAdapter{
		rmnRemoteRawCache: make(map[uint64]contracts.RawInstanceAddress),
		globalConfigCache: make(map[uint64]contracts.InstanceAddress),
	}
}

type CantonCurseAdapter struct {
	rmnRemoteRawCache map[uint64]contracts.RawInstanceAddress
	globalConfigCache map[uint64]contracts.InstanceAddress
}

// Curse implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, chain.BlockChains] {
	return cldf_ops.NewSequence(
		"curse_rmn_remote",
		semver.MustParse("2.0.0"),
		"Cursing subjects with RMNRemote",
		func(b cldf_ops.Bundle, chains chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			chain, ok := chains.CantonChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}

			rmnRemoteRaw, ok := c.rmnRemoteRawCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote instance address cached for chain %d", chain.Selector)
			}

			participant := chain.Participants[0]
			// ReadAsPartyIDs are CanReadAs rights for parties the operator cannot ActAs (e.g. ccip owner).
			// When present, exercises must be encoded as MCMS proposals instead of submitted directly.
			mcmsEnabled := len(participant.ReadAsPartyIDs) > 0
			var proposalOutputs []contract.ExerciseOutput

			for _, subject := range in.Subjects {
				report, err := cldf_ops.ExecuteOperation(b, rmn_remote.Curse, chain, contract.ChoiceInput[core.Curse]{
					InstanceAddress:    rmnRemoteRaw.InstanceAddress(),
					RawInstanceAddress: rmnRemoteRaw.String(),
					MCMSEnabled:        mcmsEnabled,
					Args: core.Curse{
						Subject: types.TEXT(hex.EncodeToString(subject[:])),
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("execute curse operation: %w", err)
				}
				if mcmsEnabled && !report.Output.Executed() {
					proposalOutputs = append(proposalOutputs, report.Output)
				}
			}

			if !mcmsEnabled {
				return output, nil
			}
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for curse: %w", err)
			}
			if len(batchOp.Transactions) == 0 {
				return sequences.OnChainOutput{}, nil
			}

			return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
		},
	)
}

// Initialize implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) Initialize(e deployment.Environment, selector uint64) error {
	rmnRemoteRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType(rmn_remote.ContractType), rmn_remote.Version, ""))
	if err != nil {
		return fmt.Errorf("get rmn remote address: %w", err)
	}

	globalConfigRef, err := e.DataStore.Addresses().Get(datastore.NewAddressRefKey(selector, datastore.ContractType(global_config.ContractType), global_config.Version, ""))
	if err != nil {
		return fmt.Errorf("get global config address: %w", err)
	}

	rmnRemoteRaw, err := dsutils.GetRawInstanceAddressFromAddressRef(rmnRemoteRef)
	if err != nil {
		return fmt.Errorf("rmn remote raw instance address: %w", err)
	}
	c.rmnRemoteRawCache[selector] = rmnRemoteRaw
	c.globalConfigCache[selector] = contracts.HexToInstanceAddress(globalConfigRef.Address)

	return nil
}

func (c *CantonCurseAdapter) getGlobalConfig(ctx context.Context, chain canton.Chain) (*core.GlobalConfig, error) {
	globalConfig, ok := c.globalConfigCache[chain.Selector]
	if !ok {
		return nil, fmt.Errorf("global config for chain %d not found in environment", chain.Selector)
	}

	participant := chain.Participants[0]
	active, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		core.GlobalConfig{}.GetTemplateID(),
		globalConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("find active contract by instance address: %w", err)
	}

	if active == nil {
		return nil, fmt.Errorf("no active contract found for global config %s", globalConfig.String())
	}

	globalConfigCreated, err := bindings.UnmarshalCreatedEvent[core.GlobalConfig](active.GetCreatedEvent())
	if err != nil {
		return nil, fmt.Errorf("unmarshal global config: %w", err)
	}

	return globalConfigCreated, nil
}

// IsChainConnectedToTargetChain implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) IsChainConnectedToTargetChain(e deployment.Environment, selector uint64, targetSel uint64) (bool, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[selector]
	if !ok {
		return false, fmt.Errorf("chain with selector %d not found in environment", selector)
	}

	globalConfigCreated, err := c.getGlobalConfig(e.GetContext(), cantonChain)
	if err != nil {
		return false, fmt.Errorf("get global config: %w", err)
	}

	// TODO: we need a helper for this numeric conversion.
	destNumeric := types.NUMERIC(fmt.Sprintf("%d.", targetSel))
	destEnabled := globalConfigCreated.DestChainConfigs[destNumeric].IsEnabled

	return bool(destEnabled), nil
}

// IsCurseEnabledForChain implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) IsCurseEnabledForChain(e deployment.Environment, selector uint64) (bool, error) {
	_, ok := c.rmnRemoteRawCache[selector]
	if !ok {
		return false, fmt.Errorf("no RMNRemote instance address cached for chain %d", selector)
	}

	return true, nil
}

// IsSubjectCursedOnChain implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) IsSubjectCursedOnChain(e deployment.Environment, selector uint64, subject fastcurse.Subject) (bool, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[selector]
	if !ok {
		return false, fmt.Errorf("chain with selector %d not found in environment", selector)
	}

	rmnRemoteRaw, ok := c.rmnRemoteRawCache[selector]
	if !ok {
		return false, fmt.Errorf("no RMNRemote instance address cached for chain %d", selector)
	}

	participant := cantonChain.Participants[0]
	active, err := contract.FindActiveContractByInstanceAddress(
		e.GetContext(),
		participant.LedgerServices.State,
		contract.LedgerQueryParties(participant),
		core.RMNRemote{}.GetTemplateID(),
		rmnRemoteRaw.InstanceAddress(),
	)
	if err != nil {
		return false, fmt.Errorf("find active contract by instance address: %w", err)
	}

	if active == nil {
		return false, fmt.Errorf("no active contract found for RMNRemote %s", rmnRemoteRaw.InstanceAddress().String())
	}

	rmnRemoteCreated, err := bindings.UnmarshalCreatedEvent[core.RMNRemote](active.GetCreatedEvent())
	if err != nil {
		return false, fmt.Errorf("unmarshal RMNRemote: %w", err)
	}

	inputSubjectHex := types.TEXT(hex.EncodeToString(subject[:]))
	containsSubject := slices.ContainsFunc(rmnRemoteCreated.CursedSubjects, func(subject types.TEXT) bool {
		return subject == inputSubjectHex
	})

	return containsSubject, nil
}

// ListConnectedChains implements [fastcurse.CurseAdapter].
// TODO: go doc comment on the interface for this method is pretty confusing.
// It seems to want all configured source chains for the given chain selector.
func (c *CantonCurseAdapter) ListConnectedChains(e deployment.Environment, selector uint64) ([]uint64, error) {
	cantonChain, ok := e.BlockChains.CantonChains()[selector]
	if !ok {
		return nil, fmt.Errorf("chain with selector %d not found in environment", selector)
	}

	globalConfigCreated, err := c.getGlobalConfig(e.GetContext(), cantonChain)
	if err != nil {
		return nil, fmt.Errorf("get global config: %w", err)
	}

	// Get all configured source chains from global config
	var sourceChainSelectors []uint64
	for selectorNumeric, config := range globalConfigCreated.SourceChainConfigs {
		if bool(config.IsEnabled) {
			parsed, err := internalparse.Uint64Checked(string(selectorNumeric))
			if err != nil {
				return nil, fmt.Errorf("parse source chain selector from GlobalConfig SourceChainConfigs map key %s: %w", selectorNumeric, err)
			}
			sourceChainSelectors = append(sourceChainSelectors, parsed)
		}
	}

	return sourceChainSelectors, nil
}

// SubjectToSelector implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) SubjectToSelector(subject fastcurse.Subject) (uint64, error) {
	return fastcurse.GenericSubjectToSelector(subject)
}

// Uncurse implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, chain.BlockChains] {
	return cldf_ops.NewSequence(
		"uncurse_rmn_remote",
		semver.MustParse("2.0.0"),
		"Uncursing subjects with RMNRemote",
		func(b cldf_ops.Bundle, chains chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			chain, ok := chains.CantonChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}

			rmnRemoteRaw, ok := c.rmnRemoteRawCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote instance address cached for chain %d", chain.Selector)
			}

			participant := chain.Participants[0]
			// ReadAsPartyIDs are CanReadAs rights for parties the operator cannot ActAs (e.g. ccip owner).
			// When present, exercises must be encoded as MCMS proposals instead of submitted directly.
			mcmsEnabled := len(participant.ReadAsPartyIDs) > 0
			var proposalOutputs []contract.ExerciseOutput

			for _, subject := range in.Subjects {
				report, err := cldf_ops.ExecuteOperation(b, rmn_remote.Uncurse, chain, contract.ChoiceInput[core.Uncurse]{
					InstanceAddress:    rmnRemoteRaw.InstanceAddress(),
					RawInstanceAddress: rmnRemoteRaw.String(),
					MCMSEnabled:        mcmsEnabled,
					Args: core.Uncurse{
						Subject: types.TEXT(hex.EncodeToString(subject[:])),
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("execute uncurse operation: %w", err)
				}
				if mcmsEnabled && !report.Output.Executed() {
					proposalOutputs = append(proposalOutputs, report.Output)
				}
			}

			if !mcmsEnabled {
				return output, nil
			}
			batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("build MCMS batch for uncurse: %w", err)
			}
			if len(batchOp.Transactions) == 0 {
				return sequences.OnChainOutput{}, nil
			}

			return sequences.OnChainOutput{BatchOps: []mcms_types.BatchOperation{batchOp}}, nil
		},
	)
}

type CantonCurseSubjectAdapter struct {
}

func NewCantonCurseSubjectAdapter() *CantonCurseSubjectAdapter {
	return &CantonCurseSubjectAdapter{}
}

// DeriveCurseAdapterVersion implements [fastcurse.CurseSubjectAdapter].
func (c *CantonCurseSubjectAdapter) DeriveCurseAdapterVersion(e deployment.Environment, selector uint64) (*semver.Version, error) {
	return rmn_remote.Version, nil
}

// SelectorToSubject implements [fastcurse.CurseSubjectAdapter].
func (c *CantonCurseSubjectAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return fastcurse.GenericSelectorToSubject(selector)
}
