package mcms

import (
	"testing"

	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	opcontract "github.com/smartcontractkit/chainlink-canton/deployment/utils/operations/contract"
)

func TestMakeMultisigId(t *testing.T) {
	tests := []struct {
		name       string
		instanceId string
		party      string
		role       cantonsdk.TimelockRole
		expected   string
	}{
		{
			name:       "proposer role",
			instanceId: "mcms-abc12",
			party:      "party::somehash",
			role:       cantonsdk.TimelockRoleProposer,
			expected:   "mcms-abc12@party::somehash-proposer",
		},
		{
			name:       "canceller role",
			instanceId: "mcms-abc12",
			party:      "party::somehash",
			role:       cantonsdk.TimelockRoleCanceller,
			expected:   "mcms-abc12@party::somehash-canceller",
		},
		{
			name:       "bypasser role",
			instanceId: "mcms-abc12",
			party:      "party::somehash",
			role:       cantonsdk.TimelockRoleBypasser,
			expected:   "mcms-abc12@party::somehash-bypasser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := makeMultisigId(tt.instanceId, tt.party, tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountTransactions(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, uint64(0), countTransactions(nil))
	})

	t.Run("single batch", func(t *testing.T) {
		bops := []mcms_types.BatchOperation{
			{Transactions: []mcms_types.Transaction{{}, {}, {}}},
		}
		assert.Equal(t, uint64(3), countTransactions(bops))
	})

	t.Run("multiple batches", func(t *testing.T) {
		bops := []mcms_types.BatchOperation{
			{Transactions: []mcms_types.Transaction{{}, {}}},
			{Transactions: []mcms_types.Transaction{{}}},
		}
		assert.Equal(t, uint64(3), countTransactions(bops))
	})
}

func TestBuildBatchFromOutputs(t *testing.T) {
	t.Run("filters executed", func(t *testing.T) {
		outputs := []opcontract.ExerciseOutput{
			{ChainSelector: 42, Tx: mcms_types.Transaction{To: "0xA"}},
			{ChainSelector: 42, Tx: mcms_types.Transaction{To: "0xB"}, ExecInfo: &opcontract.ExecInfo{UpdateID: "done"}},
		}
		batch, err := BuildBatchFromOutputs(outputs)
		require.NoError(t, err)
		assert.Len(t, batch.Transactions, 1)
		assert.Equal(t, "0xA", batch.Transactions[0].To)
	})
}
