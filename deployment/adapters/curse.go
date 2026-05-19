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

	"github.com/smartcontractkit/chainlink-canton/bindings"
	common_binding "github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/common"
	"github.com/smartcontractkit/chainlink-canton/bindings/generated/ccip/rmn"
	"github.com/smartcontractkit/chainlink-canton/contracts"
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
		rmnRemoteAddressCache:    make(map[uint64]contracts.InstanceAddress),
		rmnRemoteRawAddressCache: make(map[uint64]contracts.RawInstanceAddress),
		globalConfigCache:        make(map[uint64]contracts.InstanceAddress),
	}
}

type CantonCurseAdapter struct {
	rmnRemoteAddressCache    map[uint64]contracts.InstanceAddress
	rmnRemoteRawAddressCache map[uint64]contracts.RawInstanceAddress
	globalConfigCache        map[uint64]contracts.InstanceAddress
}

// Curse implements [fastcurse.CurseAdapter].
func (c *CantonCurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, chain.BlockChains] {
	return cldf_ops.NewSequence(
		"curse_rmn_remote",
		semver.MustParse("1.0.0"),
		"Cursing subjects with RMNRemote",
		func(b cldf_ops.Bundle, chains chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			var proposalOutputs []contract.ExerciseOutput

			chain, ok := chains.CantonChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}

			// Get the RMNRemote instance address for the provided chain selector.
			instanceAddr, ok := c.rmnRemoteAddressCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote instance address cached for chain %d", chain.Selector)
			}

			// Get the RMNRemote raw instance address for MCMS
			rawInstanceAddr, ok := c.rmnRemoteRawAddressCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote raw instance address cached for chain %d", chain.Selector)
			}

			for _, subject := range in.Subjects {
				curseOut, err := cldf_ops.ExecuteOperation(b, rmn_remote.Curse, chain, contract.ChoiceInput[rmn.Curse]{
					InstanceAddress:    instanceAddr,
					RawInstanceAddress: string(rawInstanceAddr),
					Args: rmn.Curse{
						Subject: types.TEXT(hex.EncodeToString(subject[:])),
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("execute curse operation: %w", err)
				}
				proposalOutputs = append(proposalOutputs, curseOut.Output)
			}

			// Build batch operations from all proposal outputs
			if len(proposalOutputs) > 0 {
				batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to build proposal batch for curse: %w", err)
				}
				if len(batchOp.Transactions) > 0 {
					output.BatchOps = []mcms_types.BatchOperation{batchOp}
				}
			}

			return output, nil
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

	// Get raw instance address for MCMS
	rmnRemoteRawAddr, err := dsutils.GetRawInstanceAddressFromAddressRef(rmnRemoteRef)
	if err != nil {
		return fmt.Errorf("get rmn remote raw instance address: %w", err)
	}

	c.rmnRemoteAddressCache[selector] = contracts.HexToInstanceAddress(rmnRemoteRef.Address)
	c.rmnRemoteRawAddressCache[selector] = rmnRemoteRawAddr
	c.globalConfigCache[selector] = contracts.HexToInstanceAddress(globalConfigRef.Address)

	return nil
}

func (c *CantonCurseAdapter) getGlobalConfig(ctx context.Context, chain canton.Chain) (*common_binding.GlobalConfig, error) {
	globalConfig, ok := c.globalConfigCache[chain.Selector]
	if !ok {
		return nil, fmt.Errorf("global config for chain %d not found in environment", chain.Selector)
	}

	active, err := contract.FindActiveContractByInstanceAddress(
		ctx,
		chain.Participants[0].LedgerServices.State,
		chain.Participants[0].PartyID,
		common_binding.GlobalConfig{}.GetTemplateID(),
		globalConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("find active contract by instance address: %w", err)
	}

	if active == nil {
		return nil, fmt.Errorf("no active contract found for global config %s", globalConfig.String())
	}

	globalConfigCreated, err := bindings.UnmarshalCreatedEvent[common_binding.GlobalConfig](active.GetCreatedEvent())
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
	_, ok := c.rmnRemoteAddressCache[selector]
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

	rmnRemoteInstanceAddr, ok := c.rmnRemoteAddressCache[selector]
	if !ok {
		return false, fmt.Errorf("no RMNRemote instance address cached for chain %d", selector)
	}

	active, err := contract.FindActiveContractByInstanceAddress(
		e.GetContext(),
		cantonChain.Participants[0].LedgerServices.State,
		cantonChain.Participants[0].PartyID,
		rmn.RMNRemote{}.GetTemplateID(),
		rmnRemoteInstanceAddr,
	)
	if err != nil {
		return false, fmt.Errorf("find active contract by instance address: %w", err)
	}

	if active == nil {
		return false, fmt.Errorf("no active contract found for RMNRemote %s", rmnRemoteInstanceAddr.String())
	}

	rmnRemoteCreated, err := bindings.UnmarshalCreatedEvent[rmn.RMNRemote](active.GetCreatedEvent())
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
		semver.MustParse("1.0.0"),
		"Uncursing subjects with RMNRemote",
		func(b cldf_ops.Bundle, chains chain.BlockChains, in fastcurse.CurseInput) (output sequences.OnChainOutput, err error) {
			var proposalOutputs []contract.ExerciseOutput

			chain, ok := chains.CantonChains()[in.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("chain with selector %d not found in environment", in.ChainSelector)
			}

			// Get the RMNRemote instance address for the provided chain selector.
			instanceAddr, ok := c.rmnRemoteAddressCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote instance address cached for chain %d", chain.Selector)
			}

			// Get the RMNRemote raw instance address for MCMS
			rawInstanceAddr, ok := c.rmnRemoteRawAddressCache[chain.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("no RMNRemote raw instance address cached for chain %d", chain.Selector)
			}

			for _, subject := range in.Subjects {
				uncurseOut, err := cldf_ops.ExecuteOperation(b, rmn_remote.Uncurse, chain, contract.ChoiceInput[rmn.Uncurse]{
					InstanceAddress:    instanceAddr,
					RawInstanceAddress: string(rawInstanceAddr),
					Args: rmn.Uncurse{
						Subject: types.TEXT(hex.EncodeToString(subject[:])),
					},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("execute uncurse operation: %w", err)
				}
				proposalOutputs = append(proposalOutputs, uncurseOut.Output)
			}

			// Build batch operations from all proposal outputs
			if len(proposalOutputs) > 0 {
				batchOp, err := contract.NewBatchOperationFromExercises(proposalOutputs)
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to build proposal batch for uncurse: %w", err)
				}
				if len(batchOp.Transactions) > 0 {
					output.BatchOps = []mcms_types.BatchOperation{batchOp}
				}
			}

			return output, nil
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
