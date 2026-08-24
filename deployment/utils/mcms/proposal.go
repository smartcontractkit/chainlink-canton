package mcms

import (
	"context"
	"fmt"
	"strings"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"

	"github.com/smartcontractkit/mcms"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts/v2"
	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

const DefaultTimelockExpirationHours = 72

// MCMSContractInfo describes an on-chain MCMS contract needed for proposal generation.
type MCMSContractInfo struct {
	// RawInstanceAddress is the "instanceId@partyId" format of the MCMS contract.
	RawInstanceAddress contracts.RawInstanceAddress
	// InstanceAddress is the 32-byte Keccak256 hash of the raw instance address.
	InstanceAddress contracts.InstanceAddress
}

// ProposalConfig holds parameters for generating an MCMS timelock proposal.
type ProposalConfig struct {
	// MCMSContract describes the MCMS contract to build the proposal for.
	MCMSContract MCMSContractInfo
	// ChainSelector is the Canton chain selector.
	ChainSelector mcms_types.ChainSelector
	// TimelockAddress is the timelock address for this chain (same as MCMS instance address hex for Canton).
	TimelockAddress string
	// Description is a human-readable description of the proposal.
	Description string
	// MinDelay is the timelock execution delay.
	MinDelay time.Duration
	// OverridePreviousRoot overrides the existing root if true.
	OverridePreviousRoot bool
	// Action is the timelock action (schedule, cancel, bypass).
	Action mcms_types.TimelockAction
	// ValidUntil is the Unix timestamp when the proposal expires. If zero, defaults to now + 72h.
	ValidUntil uint32
	// Role determines which MCMS role state to query for OpCount (proposer, canceller, bypasser).
	Role cantonsdk.TimelockRole
}

// GenerateTimelockProposal builds a valid mcms.TimelockProposal by querying on-chain MCMS state
// and assembling all required metadata. This mirrors Sui's GenerateProposal helper.
func GenerateTimelockProposal(
	ctx context.Context,
	stateClient apiv2.StateServiceClient,
	mcmsParties []string,
	config ProposalConfig,
	batchOps []mcms_types.BatchOperation,
) (*mcms.TimelockProposal, error) {
	if len(batchOps) == 0 {
		return nil, fmt.Errorf("at least one batch operation is required")
	}
	if len(mcmsParties) == 0 {
		return nil, fmt.Errorf("at least one party is required for MCMS ledger queries")
	}

	mcmsAddrHex := config.MCMSContract.InstanceAddress.Hex()

	inspector := cantonsdk.NewInspector(stateClient, mcmsParties, config.Role)
	opCount, err := inspector.GetOpCount(ctx, mcmsAddrHex)
	if err != nil {
		return nil, fmt.Errorf("failed to get op count: %w", err)
	}

	mcmsContract, err := cantonsdk.GetMCMSContract(ctx, stateClient, mcmsParties, mcmsAddrHex)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCMS contract state: %w", err)
	}

	multisigId := makeMultisigId(string(mcmsContract.InstanceId), string(mcmsContract.Owner), config.Role)

	// postOpCount and overridePreviousRoot are applied at sign/execute time by mcms (encoder + executor), not in chainMetadata.
	metadata, err := cantonsdk.NewChainMetadata(
		opCount,
		int64(mcmsContract.ChainId),
		multisigId,
		mcmsAddrHex,
		string(mcmsContract.InstanceId),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain metadata: %w", err)
	}

	validUntil := config.ValidUntil
	if validUntil == 0 {
		const maxUint32 int64 = 1<<32 - 1
		unixTS := time.Now().Add(time.Duration(DefaultTimelockExpirationHours) * time.Hour).Unix()
		if unixTS < 0 || unixTS > maxUint32 {
			return nil, fmt.Errorf("computed validUntil timestamp %d overflows uint32", unixTS)
		}
		validUntil = uint32(unixTS)
	}

	timelockAddr := config.TimelockAddress
	if timelockAddr == "" {
		timelockAddr = mcmsAddrHex
	}

	builder := mcms.NewTimelockProposalBuilder().
		SetVersion("v1").
		SetValidUntil(validUntil).
		SetDescription(config.Description).
		SetOverridePreviousRoot(config.OverridePreviousRoot).
		AddTimelockAddress(config.ChainSelector, timelockAddr).
		AddChainMetadata(config.ChainSelector, metadata).
		SetAction(config.Action)

	if config.Action == mcms_types.TimelockActionSchedule {
		builder.SetDelay(mcms_types.NewDuration(config.MinDelay))
	}

	for _, bop := range batchOps {
		builder.AddOperation(bop)
	}

	return builder.Build()
}

// BuildBatchFromOutputs is a convenience wrapper around opcontract.NewBatchOperationFromExercises.
func BuildBatchFromOutputs(outputs []opcontract.ExerciseOutput) (mcms_types.BatchOperation, error) {
	return opcontract.NewBatchOperationFromExercises(outputs)
}

// makeMultisigId constructs the multisig ID in the format expected by the Canton MCMS SDK:
// "instanceId@partyId-role" (e.g., "mcms-abc12@party::hash-proposer").
func makeMultisigId(instanceId, party string, role cantonsdk.TimelockRole) string {
	return fmt.Sprintf("%s@%s-%s", instanceId, party, strings.ToLower(role.String()))
}

// countTransactions counts the total number of transactions across all batch operations.
func countTransactions(batchOps []mcms_types.BatchOperation) uint64 {
	var count uint64
	for _, bop := range batchOps {
		count += uint64(len(bop.Transactions))
	}

	return count
}
