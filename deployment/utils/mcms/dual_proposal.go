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

const ccvFactoryInstanceSubstring = "ccip-factory-ccv@"

var ccvOwnerFunctionNames = map[string]struct{}{
	"DeployCommitteeVerifier":       {},
	"ApplyRemoteChainConfigUpdates": {},
	"ApplyAllowListUpdates":         {},
	"ApplySignatureConfigs":         {},
}

// SplitBatchOpsByOwner separates Canton MCMS transactions into ccip-owner vs ccv-owner batches.
func SplitBatchOpsByOwner(batchOps []mcms_types.BatchOperation) (ccip []mcms_types.BatchOperation, ccv []mcms_types.BatchOperation) {
	for _, op := range batchOps {
		var ccipTxs, ccvTxs []mcms_types.Transaction
		for _, tx := range op.Transactions {
			if transactionUsesCCVOwner(tx) {
				ccvTxs = append(ccvTxs, tx)
			} else {
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
	}

	return ccip, ccv
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

// BuildDualTimelockProposalsFromBatchOps builds ccipOwner (_0) and ccvOwner (_1) timelock proposals
// for one Canton chain. The proposals target separate MCMS instances with independent roots and
// may be signed and scheduled in parallel; neither depends on the other being executed first.
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
	ccipBatch, ccvBatch := SplitBatchOpsByOwner(batchOps)

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
