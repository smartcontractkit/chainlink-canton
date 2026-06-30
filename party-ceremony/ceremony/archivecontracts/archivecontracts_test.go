package archivecontracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/archivecontracts"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
)

func TestParseTemplateSelector(t *testing.T) {
	t.Parallel()

	tpl, err := archivecontracts.ParseTemplateSelector("#ccip-common:CCIP.GlobalConfig:GlobalConfig")
	require.NoError(t, err)
	assert.Equal(t, "ccip-common", tpl.PackageName)
	assert.Equal(t, "CCIP.GlobalConfig", tpl.ModuleName)
	assert.Equal(t, "GlobalConfig", tpl.EntityName)

	tpl, err = archivecontracts.ParseTemplateSelector("c24ecaf0a38cbbd7943b191208728703ae1fb060cc9971fd49161da17efee29c:CCIP.GlobalConfig:GlobalConfig")
	require.NoError(t, err)
	assert.Equal(t, "c24ecaf0a38cbbd7943b191208728703ae1fb060cc9971fd49161da17efee29c", tpl.PackageID)
}

func TestSplitTargets(t *testing.T) {
	t.Parallel()

	targets := []ledger.ArchiveTarget{
		{ContractID: "a"},
		{ContractID: "b"},
		{ContractID: "c"},
	}
	batches := archivecontracts.SplitTargetsForTest(targets, 2)
	require.Len(t, batches, 2)
	assert.Len(t, batches[0], 2)
	assert.Len(t, batches[1], 1)

	defaultBatches := archivecontracts.SplitTargetsForTest(targets, 0)
	require.Len(t, defaultBatches, 3)
	for _, b := range defaultBatches {
		assert.Len(t, b, 1)
	}
}
