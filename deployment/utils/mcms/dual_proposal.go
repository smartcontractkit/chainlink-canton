package mcms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/mcms"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	dsutils "github.com/smartcontractkit/chainlink-canton/deployment/utils/datastore"
)

const (
	ccvFactoryInstanceSubstring = "ccip-factory-ccv@"
	rmnFactoryInstanceSubstring = "ccip-factory-rmn@"
)

// ccvOwnerFunctionNames lists MCMS entrypoint function names from CommitteeVerifier.daml.
// Keep in sync with contracts/ccip/committee-verifier/daml/CCIP/CommitteeVerifier.daml.
var ccvOwnerFunctionNames = map[string]struct{}{
	"DeployCommitteeVerifier":       {},
	"ApplySignatureConfigs":         {},
	"SetDynamicConfig":              {},
	"ApplyRemoteChainConfigUpdates": {},
	"SetDeps":                       {},
	"ApplyAllowListUpdates":         {},
	"TransferStorageLocationsAdmin": {},
	"AcceptStorageLocationsAdmin":   {},
	"UpdateStorageLocations":        {},
}

// rmnOwnerFunctionNames lists MCMS entrypoint function names from RMNRemote.daml and
// DeployRMNRemote on the rmn-qualified factory.
// Keep in sync with contracts/ccip/core/daml/CCIP/RMNRemote.daml and CCIP/Factory.daml.
var rmnOwnerFunctionNames = map[string]struct{}{
	"DeployRMNRemote":       {},
	"Curse":                 {},
	"Uncurse":               {},
	"CurseChain":            {},
	"UncurseChain":          {},
	"CurseMultiple":         {},
	"UncurseMultiple":       {},
	"AddCustomObservers":    {},
	"RemoveCustomObservers": {},
}

// SplitBatchOpsByOwner separates Canton MCMS transactions into ccipOwner, ccvOwner, and rmnOwner batches.
func SplitBatchOpsByOwner(batchOps []mcms_types.BatchOperation) (ccip, ccv, rmn []mcms_types.BatchOperation) {
	for _, op := range batchOps {
		var ccipTxs, ccvTxs, rmnTxs []mcms_types.Transaction
		for _, tx := range op.Transactions {
			switch {
			case transactionUsesRMNOwner(tx):
				rmnTxs = append(rmnTxs, tx)
			case transactionUsesCCVOwner(tx):
				ccvTxs = append(ccvTxs, tx)
			default:
				ccipTxs = append(ccipTxs, tx)
			}
		}
		if len(ccipTxs) > 0 {
			ccip = append(ccip, mcms_types.BatchOperation{
				ChainSelector: op.ChainSelector,
				Transactions:  ccipTxs,
			})
		}
		if len(ccvTxs) > 0 {
			ccv = append(ccv, mcms_types.BatchOperation{
				ChainSelector: op.ChainSelector,
				Transactions:  ccvTxs,
			})
		}
		if len(rmnTxs) > 0 {
			rmn = append(rmn, mcms_types.BatchOperation{
				ChainSelector: op.ChainSelector,
				Transactions:  rmnTxs,
			})
		}
	}

	return ccip, ccv, rmn
}

func transactionUsesCCVOwner(tx mcms_types.Transaction) bool {
	if tx.ContractType == "CommitteeVerifier" {
		return true
	}
	if len(tx.AdditionalFields) == 0 {
		return false
	}
	var fields cantonsdk.AdditionalFields
	if err := json.Unmarshal(tx.AdditionalFields, &fields); err != nil {
		return false
	}
	if strings.Contains(fields.TargetInstanceAddress, ccvFactoryInstanceSubstring) {
		return true
	}
	if _, ok := ccvOwnerFunctionNames[fields.FunctionName]; ok {
		return true
	}

	return false
}

func transactionUsesRMNOwner(tx mcms_types.Transaction) bool {
	if tx.ContractType == "RMNRemote" {
		return true
	}
	if len(tx.AdditionalFields) == 0 {
		return false
	}
	var fields cantonsdk.AdditionalFields
	if err := json.Unmarshal(tx.AdditionalFields, &fields); err != nil {
		return false
	}
	if strings.Contains(fields.TargetInstanceAddress, rmnFactoryInstanceSubstring) {
		return true
	}
	if _, ok := rmnOwnerFunctionNames[fields.FunctionName]; ok {
		return true
	}

	return false
}

// BuildDualTimelockProposalsFromBatchOps builds ccipOwner, ccvOwner, and rmnOwner timelock proposals
// for one Canton chain. Proposals target separate MCMS instances and may be signed in parallel.
func BuildDualTimelockProposalsFromBatchOps(
	ctx context.Context,
	e cldf.Environment,
	chain canton.Chain,
	chainSelector uint64,
	batchOps []mcms_types.BatchOperation,
	input ccipmcms.Input,
	description string,
) ([]mcms.TimelockProposal, error) {
	batchOps = ConsolidateBatchOpsPerChain(batchOps)
	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner(batchOps)

	var proposals []mcms.TimelockProposal

	ccipProposal, err := BuildTimelockProposalForOwner(
		ctx, e, chain, chainSelector, QualifierCCIPOwner, ccipBatch, input,
		proposalDescription(description, QualifierCCIPOwner),
	)
	if err != nil {
		return nil, fmt.Errorf("build ccipOwner proposal for chain %d: %w", chainSelector, err)
	}
	if ccipProposal != nil {
		proposals = append(proposals, *ccipProposal)
	}

	ccvProposal, err := BuildTimelockProposalForOwner(
		ctx, e, chain, chainSelector, QualifierCCVOwner, ccvBatch, input,
		proposalDescription(description, QualifierCCVOwner),
	)
	if err != nil {
		return nil, fmt.Errorf("build ccvOwner proposal for chain %d: %w", chainSelector, err)
	}
	if ccvProposal != nil {
		proposals = append(proposals, *ccvProposal)
	}

	rmnProposal, err := BuildTimelockProposalForOwner(
		ctx, e, chain, chainSelector, QualifierRMNOwner, rmnBatch, input,
		proposalDescription(description, QualifierRMNOwner),
	)
	if err != nil {
		return nil, fmt.Errorf("build rmnOwner proposal for chain %d: %w", chainSelector, err)
	}
	if rmnProposal != nil {
		proposals = append(proposals, *rmnProposal)
	}

	return proposals, nil
}

// BatchOpsFromTimelockProposals flattens MCMS timelock proposal operations into batch ops.
func BatchOpsFromTimelockProposals(proposals []mcms.TimelockProposal) []mcms_types.BatchOperation {
	var ops []mcms_types.BatchOperation
	for _, prop := range proposals {
		ops = append(ops, prop.Operations...)
	}

	return ConsolidateBatchOpsPerChain(ops)
}

// ConsolidateBatchOpsPerChain merges multiple batch operations for the same chain into one.
func ConsolidateBatchOpsPerChain(ops []mcms_types.BatchOperation) []mcms_types.BatchOperation {
	txsByChain := make(map[mcms_types.ChainSelector][]mcms_types.Transaction)
	for _, op := range ops {
		if len(op.Transactions) == 0 {
			continue
		}
		txsByChain[op.ChainSelector] = append(txsByChain[op.ChainSelector], op.Transactions...)
	}

	out := make([]mcms_types.BatchOperation, 0, len(txsByChain))
	for chainSel, txs := range txsByChain {
		out = append(out, mcms_types.BatchOperation{
			ChainSelector: chainSel,
			Transactions:  txs,
		})
	}

	return out
}

func proposalDescription(base, owner string) string {
	if base == "" {
		return owner
	}

	return base + " (" + owner + ")"
}

// BuildTimelockProposalForOwner builds one timelock proposal against the MCMS instance for ownerQualifier.
func BuildTimelockProposalForOwner(
	ctx context.Context,
	e cldf.Environment,
	chain canton.Chain,
	chainSelector uint64,
	ownerQualifier string,
	batchOps []mcms_types.BatchOperation,
	input ccipmcms.Input,
	description string,
) (*mcms.TimelockProposal, error) {
	if len(batchOps) == 0 {
		return nil, nil
	}
	if len(chain.Participants) == 0 {
		return nil, fmt.Errorf("canton chain %d has no participants", chainSelector)
	}

	raw, err := dsutils.MCMSRawInstanceAddress(e.DataStore, chainSelector, ownerQualifier)
	if err != nil {
		return nil, err
	}

	participant := chain.Participants[0]
	role := cantonsdkRoleForAction(input.TimelockAction)

	return GenerateTimelockProposal(
		ctx,
		participant.LedgerServices.State,
		participant.PartyID,
		ProposalConfig{
			MCMSContract: MCMSContractInfo{
				RawInstanceAddress: raw,
				InstanceAddress:    raw.InstanceAddress(),
			},
			ChainSelector:        mcms_types.ChainSelector(chainSelector),
			TimelockAddress:      raw.InstanceAddress().String(),
			Description:          description,
			MinDelay:             timelockDelayAsDuration(input.TimelockDelay),
			OverridePreviousRoot: input.OverridePreviousRoot,
			Action:               input.TimelockAction,
			ValidUntil:           input.ValidUntil,
			Role:                 role,
		},
		batchOps,
	)
}

func timelockDelayAsDuration(d mcms_types.Duration) time.Duration {
	parsed, err := time.ParseDuration(d.String())
	if err != nil {
		return time.Second
	}

	return parsed
}

func cantonsdkRoleForAction(action mcms_types.TimelockAction) cantonsdk.TimelockRole {
	switch action {
	case mcms_types.TimelockActionBypass:
		return cantonsdk.TimelockRoleBypasser
	case mcms_types.TimelockActionCancel:
		return cantonsdk.TimelockRoleCanceller
	default:
		return cantonsdk.TimelockRoleProposer
	}
}
