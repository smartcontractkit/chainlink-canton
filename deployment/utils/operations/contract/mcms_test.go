package contract

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/go-daml/pkg/bind"
	cantonsdk "github.com/smartcontractkit/mcms/sdk/canton"
	mcms_types "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/contracts"
)

func TestNewCantonTransaction(t *testing.T) {
	t.Parallel()
	rawAddr := "globalconfig-abc12@party::somehash"
	instanceAddr := contracts.HexToInstanceAddress("0xdeadbeef")
	encodedChoice := &bind.EncodedChoice{
		TemplateID: bind.TemplateInformation{
			PackageID:    "pkg-123",
			ModuleName:   "CCIP.GlobalConfig",
			TemplateName: "GlobalConfig",
		},
		Choice:        "ApplyDestChainConfigUpdates",
		OperationData: "abcdef0123456789",
	}
	contractType := deployment.ContractType("CantonGlobalConfig")

	templateID := "#pkg-123:CCIP.GlobalConfig:GlobalConfig"
	tx, err := NewCantonTransaction(rawAddr, instanceAddr, encodedChoice, contractType, templateID)
	require.NoError(t, err)

	assert.Equal(t, instanceAddr.Hex(), tx.To)
	assert.Equal(t, string(contractType), tx.ContractType)

	expectedData, _ := hex.DecodeString("abcdef0123456789")
	assert.Equal(t, expectedData, tx.Data)

	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	assert.Equal(t, rawAddr, af.TargetInstanceAddress, "TargetInstanceAddress should be raw instanceId@partyId format")
	assert.Equal(t, "ApplyDestChainConfigUpdates", af.FunctionName)
	assert.Equal(t, "abcdef0123456789", af.OperationData)
	assert.Equal(t, templateID, af.TargetTemplateID, "TargetTemplateID should be set for dynamic CID resolution")
}

func TestNewCantonTransaction_RawAddressNotHex(t *testing.T) {
	t.Parallel()
	rawAddr := "globalconfig-abc12@party::somehash"
	instanceAddr := contracts.HexToInstanceAddress("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	encodedChoice := &bind.EncodedChoice{
		Choice:        "SetConfig",
		OperationData: "ff",
	}

	tx, err := NewCantonTransaction(rawAddr, instanceAddr, encodedChoice, "TestType", "#test:Module:Entity")
	require.NoError(t, err)

	var af cantonsdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &af))
	assert.Equal(t, rawAddr, af.TargetInstanceAddress, "should use raw format, not hex")
	assert.Equal(t, instanceAddr.Hex(), tx.To, "To field should use hex format")
}

func TestNewBatchOperationFromExercises(t *testing.T) {
	t.Parallel()
	t.Run("empty slice returns empty batch", func(t *testing.T) {
		t.Parallel()
		batch, err := NewBatchOperationFromExercises(nil)
		require.NoError(t, err)
		assert.Empty(t, batch.Transactions)
	})

	t.Run("filters out executed outputs", func(t *testing.T) {
		t.Parallel()
		outs := []ExerciseOutput{
			{
				ChainSelector: 100,
				Tx:            mcms_types.Transaction{To: "0x1"},
				ExecInfo:      &ExecInfo{UpdateID: "done"},
			},
			{
				ChainSelector: 100,
				Tx:            mcms_types.Transaction{To: "0x2"},
			},
		}
		batch, err := NewBatchOperationFromExercises(outs)
		require.NoError(t, err)
		assert.Len(t, batch.Transactions, 1)
		assert.Equal(t, "0x2", batch.Transactions[0].To)
		assert.Equal(t, mcms_types.ChainSelector(100), batch.ChainSelector)
	})

	t.Run("all executed returns empty batch", func(t *testing.T) {
		t.Parallel()
		outs := []ExerciseOutput{
			{
				ChainSelector: 100,
				Tx:            mcms_types.Transaction{To: "0x1"},
				ExecInfo:      &ExecInfo{UpdateID: "done"},
			},
		}
		batch, err := NewBatchOperationFromExercises(outs)
		require.NoError(t, err)
		assert.Empty(t, batch.Transactions)
	})

	t.Run("multiple unexecuted on same chain", func(t *testing.T) {
		t.Parallel()
		outs := []ExerciseOutput{
			{ChainSelector: 42, Tx: mcms_types.Transaction{To: "0xA"}},
			{ChainSelector: 42, Tx: mcms_types.Transaction{To: "0xB"}},
			{ChainSelector: 42, Tx: mcms_types.Transaction{To: "0xC"}},
		}
		batch, err := NewBatchOperationFromExercises(outs)
		require.NoError(t, err)
		assert.Len(t, batch.Transactions, 3)
		assert.Equal(t, mcms_types.ChainSelector(42), batch.ChainSelector)
	})

	t.Run("error on multiple chains", func(t *testing.T) {
		t.Parallel()
		outs := []ExerciseOutput{
			{ChainSelector: 1, Tx: mcms_types.Transaction{To: "0xA"}},
			{ChainSelector: 2, Tx: mcms_types.Transaction{To: "0xB"}},
		}
		_, err := NewBatchOperationFromExercises(outs)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple chains")
	})
}

func TestExerciseOutput_Executed(t *testing.T) {
	t.Parallel()
	t.Run("not executed when ExecInfo is nil", func(t *testing.T) {
		t.Parallel()
		out := ExerciseOutput{}
		assert.False(t, out.Executed())
	})

	t.Run("executed when ExecInfo is set", func(t *testing.T) {
		t.Parallel()
		out := ExerciseOutput{ExecInfo: &ExecInfo{UpdateID: "123"}}
		assert.True(t, out.Executed())
	})
}

func TestValidateCantonAdditionalFields(t *testing.T) {
	t.Parallel()
	t.Run("valid fields with raw address", func(t *testing.T) {
		t.Parallel()
		af := cantonsdk.AdditionalFields{
			TargetInstanceAddress: "globalconfig-abc12@party::somehash",
			FunctionName:          "SetConfig",
			OperationData:         "abcdef",
		}
		raw, _ := json.Marshal(af)
		require.NoError(t, ValidateCantonAdditionalFields(raw))
	})

	t.Run("missing target instance address", func(t *testing.T) {
		t.Parallel()
		af := cantonsdk.AdditionalFields{
			FunctionName:  "SetConfig",
			OperationData: "abcdef",
		}
		raw, _ := json.Marshal(af)
		err := ValidateCantonAdditionalFields(raw)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "targetInstanceAddress")
	})

	t.Run("missing function name", func(t *testing.T) {
		t.Parallel()
		af := cantonsdk.AdditionalFields{
			TargetInstanceAddress: "globalconfig-abc12@party::somehash",
			OperationData:         "abcdef",
		}
		raw, _ := json.Marshal(af)
		err := ValidateCantonAdditionalFields(raw)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "functionName")
	})

	t.Run("missing operation data", func(t *testing.T) {
		t.Parallel()
		af := cantonsdk.AdditionalFields{
			TargetInstanceAddress: "globalconfig-abc12@party::somehash",
			FunctionName:          "SetConfig",
		}
		raw, _ := json.Marshal(af)
		err := ValidateCantonAdditionalFields(raw)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "operationData")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		err := ValidateCantonAdditionalFields(json.RawMessage(`{invalid`))
		require.Error(t, err)
	})
}
