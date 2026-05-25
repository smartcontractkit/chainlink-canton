package mcms

import (
	"encoding/json"
	"testing"

	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitBatchOpsByOwner_laneConfigure(t *testing.T) {
	t.Parallel()

	ccipFields, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: "globalconfig-abc@party::1220",
		FunctionName:          "ApplyDestChainConfigUpdates",
	})
	require.NoError(t, err)

	ccvFields, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: "committeeverifier-xyz@party::1220",
		FunctionName:          "ApplySignatureConfigs",
	})
	require.NoError(t, err)

	batchOps := []mcms_types.BatchOperation{{
		ChainSelector: 8706591216959472610,
		Transactions: []mcms_types.Transaction{
			{AdditionalFields: ccipFields},
			{AdditionalFields: ccvFields},
			{AdditionalFields: ccipFields},
		},
	}}

	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner(batchOps)
	require.Len(t, ccipBatch, 1)
	require.Len(t, ccvBatch, 1)
	require.Empty(t, rmnBatch)
	assert.Len(t, ccipBatch[0].Transactions, 2)
	assert.Len(t, ccvBatch[0].Transactions, 1)
}

func TestSplitBatchOpsByOwner_rmnCurse(t *testing.T) {
	t.Parallel()

	rmnFields, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: "rmn-abc@party::1220",
		FunctionName:          "CurseChain",
	})
	require.NoError(t, err)

	ccipFields, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: "globalconfig-abc@party::1220",
		FunctionName:          "ApplyDestChainConfigUpdates",
	})
	require.NoError(t, err)

	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner([]mcms_types.BatchOperation{{
		ChainSelector: 8706591216959472610,
		Transactions: []mcms_types.Transaction{
			{AdditionalFields: ccipFields},
			{AdditionalFields: rmnFields},
		},
	}})
	require.Len(t, ccipBatch, 1)
	require.Empty(t, ccvBatch)
	require.Len(t, rmnBatch, 1)
	assert.Len(t, ccipBatch[0].Transactions, 1)
	assert.Len(t, rmnBatch[0].Transactions, 1)
}

func TestSplitBatchOpsByOwner_rmnFactoryDeploy(t *testing.T) {
	t.Parallel()

	fields, err := json.Marshal(cantonsdk.AdditionalFields{
		TargetInstanceAddress: "ccip-factory-rmn@party::1220",
		FunctionName:          "DeployRMNRemote",
	})
	require.NoError(t, err)

	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner([]mcms_types.BatchOperation{{
		ChainSelector: 8706591216959472610,
		Transactions:  []mcms_types.Transaction{{AdditionalFields: fields}},
	}})
	require.Empty(t, ccipBatch)
	require.Empty(t, ccvBatch)
	require.Len(t, rmnBatch, 1)
}

func TestSplitBatchOpsByOwner_committeeVerifierMCMSEntrypoints(t *testing.T) {
	t.Parallel()

	sel := mcms_types.ChainSelector(8706591216959472610)
	ccvFunctions := []string{
		"SetDynamicConfig",
		"SetDeps",
		"TransferStorageLocationsAdmin",
		"AcceptStorageLocationsAdmin",
		"UpdateStorageLocations",
	}

	for _, fn := range ccvFunctions {
		t.Run(fn, func(t *testing.T) {
			t.Parallel()

			fields, err := json.Marshal(cantonsdk.AdditionalFields{
				TargetInstanceAddress: "committeeverifier-xyz@party::1220",
				FunctionName:          fn,
			})
			require.NoError(t, err)

			ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner([]mcms_types.BatchOperation{{
				ChainSelector: sel,
				Transactions:  []mcms_types.Transaction{{AdditionalFields: fields}},
			}})
			require.Empty(t, ccipBatch)
			require.Len(t, ccvBatch, 1)
			require.Empty(t, rmnBatch)
			assert.Len(t, ccvBatch[0].Transactions, 1)
		})
	}
}

func TestSplitBatchOpsByOwner_contractTypeCommitteeVerifier(t *testing.T) {
	t.Parallel()

	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner([]mcms_types.BatchOperation{{
		ChainSelector: 8706591216959472610,
		Transactions: []mcms_types.Transaction{{
			OperationMetadata: mcms_types.OperationMetadata{ContractType: "CommitteeVerifier"},
		}},
	}})
	require.Empty(t, ccipBatch)
	require.Len(t, ccvBatch, 1)
	require.Empty(t, rmnBatch)
}

func TestSplitBatchOpsByOwner_contractTypeRMNRemote(t *testing.T) {
	t.Parallel()

	ccipBatch, ccvBatch, rmnBatch := SplitBatchOpsByOwner([]mcms_types.BatchOperation{{
		ChainSelector: 8706591216959472610,
		Transactions: []mcms_types.Transaction{{
			OperationMetadata: mcms_types.OperationMetadata{ContractType: "RMNRemote"},
		}},
	}})
	require.Empty(t, ccipBatch)
	require.Empty(t, ccvBatch)
	require.Len(t, rmnBatch, 1)
}

func TestConsolidateBatchOpsPerChain(t *testing.T) {
	t.Parallel()

	sel := mcms_types.ChainSelector(8706591216959472610)
	ops := []mcms_types.BatchOperation{
		{ChainSelector: sel, Transactions: []mcms_types.Transaction{{To: "0x1"}}},
		{ChainSelector: sel, Transactions: []mcms_types.Transaction{{To: "0x2"}}},
	}

	consolidated := ConsolidateBatchOpsPerChain(ops)
	require.Len(t, consolidated, 1)
	assert.Len(t, consolidated[0].Transactions, 2)
}
