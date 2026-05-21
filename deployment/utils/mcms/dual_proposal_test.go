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

	ccipBatch, ccvBatch := SplitBatchOpsByOwner(batchOps)
	require.Len(t, ccipBatch, 1)
	require.Len(t, ccvBatch, 1)
	assert.Len(t, ccipBatch[0].Transactions, 2)
	assert.Len(t, ccvBatch[0].Transactions, 1)
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
