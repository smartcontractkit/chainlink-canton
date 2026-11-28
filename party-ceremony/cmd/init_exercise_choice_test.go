package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
)

func TestParseTemplateID(t *testing.T) {
	t.Parallel()

	pkg, module, entity, err := parseTemplateID("#utility-registry-app-v0:Utility.Registry.App.V0.Service.Provider:ProviderService")
	require.NoError(t, err)
	assert.Equal(t, "#utility-registry-app-v0", pkg)
	assert.Equal(t, "Utility.Registry.App.V0.Service.Provider", module)
	assert.Equal(t, "ProviderService", entity)

	for _, bad := range []string{"", "pkg:Module", "pkg:Module:Entity:Extra", "pkg::Entity", ":Module:Entity"} {
		_, _, _, err := parseTemplateID(bad)
		assert.Error(t, err, "template ID %q", bad)
	}
}

func TestReadChoiceArgument(t *testing.T) {
	t.Parallel()

	arg, err := readChoiceArgument("")
	require.NoError(t, err)
	assert.JSONEq(t, emptyRecordValue, arg, "no file means a choice without fields")

	dir := t.TempDir()
	good := filepath.Join(dir, "good.json")
	require.NoError(t, os.WriteFile(good, []byte(`{"record":{"fields":[{"label":"instrumentId","value":{"text":"LINK"}}]}}`), 0o600))
	arg, err = readChoiceArgument(good)
	require.NoError(t, err)
	assert.Contains(t, arg, "LINK")

	recordNotValue := filepath.Join(dir, "record.json")
	require.NoError(t, os.WriteFile(recordNotValue, []byte(`{"fields":[]}`), 0o600))
	_, err = readChoiceArgument(recordNotValue)
	require.ErrorContains(t, err, "not a Ledger API Value", "a bare Record is rejected before the ceremony starts")

	for name, body := range map[string]string{"empty.json": `{}`, "text.json": `{"text":"LINK"}`} {
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
		_, err = readChoiceArgument(path)
		require.ErrorContains(t, err, "must hold a record value", "%s parses as a Value but is not a record", body)
	}

	_, err = readChoiceArgument(filepath.Join(dir, "missing.json"))
	require.ErrorContains(t, err, "reading choice args file")
}

func TestReadDisclosedContracts(t *testing.T) {
	t.Parallel()

	disclosed, err := readDisclosedContracts("")
	require.NoError(t, err)
	assert.Nil(t, disclosed, "no file means no disclosed contracts")

	dir := t.TempDir()
	good := filepath.Join(dir, "good.json")
	require.NoError(t, os.WriteFile(good, []byte(`[{"templateId":"#pkg:Mod:Ent","contractId":"cid-1","createdEventBlob":"YmxvYg==","synchronizerId":"global"}]`), 0o600))
	disclosed, err = readDisclosedContracts(good)
	require.NoError(t, err)
	assert.Equal(t, []ledger.DisclosedContract{{TemplateID: "#pkg:Mod:Ent", ContractID: "cid-1", CreatedEventBlob: "YmxvYg==", SynchronizerID: "global"}}, disclosed)

	for name, body := range map[string]string{
		"no-contract-id.json": `[{"createdEventBlob":"YmxvYg=="}]`,
		"no-blob.json":        `[{"contractId":"cid-1"}]`,
		"bad-blob.json":       `[{"contractId":"cid-1","createdEventBlob":"not base64!"}]`,
		"not-a-list.json":     `{"contractId":"cid-1","createdEventBlob":"YmxvYg=="}`,
	} {
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
		_, err = readDisclosedContracts(path)
		assert.Error(t, err, "%s is rejected before the ceremony starts", name)
	}
}
