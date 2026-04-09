package contract

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

// NewCantonTransaction converts an EncodedChoice into an MCMS Transaction with Canton-specific AdditionalFields.
// rawInstanceAddress must be in "instanceId@partyId" format (e.g. "globalconfig-abc12@party::hash")
// as required by the Canton MCMS SDK for Merkle leaf hashing.
// instanceAddress is the 32-byte Keccak256 hash used for the Transaction.To field.
func NewCantonTransaction(
	rawInstanceAddress string,
	instanceAddress contracts.InstanceAddress,
	encodedChoice *bind.EncodedChoice,
	contractType deployment.ContractType,
) (mcms_types.Transaction, error) {
	af, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: rawInstanceAddress,
		FunctionName:          encodedChoice.Choice,
		OperationData:         encodedChoice.OperationData,
	})
	if err != nil {
		return mcms_types.Transaction{}, fmt.Errorf("failed to marshal canton additional fields: %w", err)
	}

	opData, _ := hex.DecodeString(encodedChoice.OperationData)

	return mcms_types.Transaction{
		OperationMetadata: mcms_types.OperationMetadata{
			ContractType: string(contractType),
		},
		To:               instanceAddress.Hex(),
		Data:             opData,
		AdditionalFields: af,
	}, nil
}

// NewBatchOperationFromExercises constructs an MCMS BatchOperation from a slice of ExerciseOutputs.
// It filters out any ExerciseOutputs that have already been executed.
// Returns an error if the ExerciseOutputs target multiple chains.
// If all ExerciseOutputs are executed, it returns an empty BatchOperation and no error.
func NewBatchOperationFromExercises(outs []ExerciseOutput) (mcms_types.BatchOperation, error) {
	if len(outs) == 0 {
		return mcms_types.BatchOperation{}, nil
	}

	var (
		chainSelector uint64
		txs           []mcms_types.Transaction
	)
	for _, out := range outs {
		if out.Executed() {
			continue
		}
		if len(txs) == 0 {
			chainSelector = out.ChainSelector
			txs = append(txs, out.Tx)

			continue
		}
		if out.ChainSelector != chainSelector {
			return mcms_types.BatchOperation{}, errors.New("failed to make batch operation: exercises target multiple chains")
		}
		txs = append(txs, out.Tx)
	}

	if len(txs) == 0 {
		return mcms_types.BatchOperation{}, nil
	}

	return mcms_types.BatchOperation{
		ChainSelector: mcms_types.ChainSelector(chainSelector),
		Transactions:  txs,
	}, nil
}

// ValidateCantonAdditionalFields validates that AdditionalFields JSON can be unmarshaled
// and has the required fields populated.
func ValidateCantonAdditionalFields(raw json.RawMessage) error {
	var af cantonsdk.AdditionalFields
	if err := json.Unmarshal(raw, &af); err != nil {
		return fmt.Errorf("failed to unmarshal Canton additional fields: %w", err)
	}
	if af.TargetInstanceAddress == "" {
		return errors.New("targetInstanceAddress is required")
	}
	if af.FunctionName == "" {
		return errors.New("functionName is required")
	}
	if af.OperationData == "" {
		return errors.New("operationData is required")
	}

	return nil
}
