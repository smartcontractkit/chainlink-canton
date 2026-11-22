package ledger_test

import (
	"encoding/base64"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony/ops/ledger"
)

func TestExpectedExerciseCheck_RefusesAPreparedTransactionItCannotRead(t *testing.T) {
	t.Parallel()

	for name, preparedTxB64 := range map[string]string{
		"not base64":     "not base64!",
		"not a protobuf": base64.StdEncoding.EncodeToString([]byte{0xff}),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&ledger.ExpectedExercise{Choice: "Archive"}).Check(preparedTxB64)

			require.ErrorIs(t, err, ledger.ErrPreparedTxMismatch)
		})
	}
}

func preparedExerciseB64(t *testing.T, e ledger.ExpectedExercise, chosenJSON string) string {
	t.Helper()

	return preparedViaB64(t, e, chosenJSON, &apiv2.Identifier{PackageId: e.PackageID, ModuleName: e.TemplateModule, EntityName: e.TemplateEntity}, nil, "")
}

func preparedViaB64(t *testing.T, e ledger.ExpectedExercise, chosenJSON string, templateID, interfaceID *apiv2.Identifier, packageName string) string {
	t.Helper()
	var chosen apiv2.Value
	require.NoError(t, protojson.Unmarshal([]byte(chosenJSON), &chosen))
	tx := &interactive.PreparedTransaction{Transaction: &interactive.DamlTransaction{
		Roots: []string{"0"},
		Nodes: []*interactive.DamlTransaction_Node{{
			NodeId: "0",
			VersionedNode: &interactive.DamlTransaction_Node_V1{V1: &v1.Node{NodeType: &v1.Node_Exercise{Exercise: &v1.Exercise{
				TemplateId:  templateID,
				InterfaceId: interfaceID,
				PackageName: packageName,
				ContractId:  e.ContractID,
				ChoiceId:    e.Choice,
				ChosenValue: &chosen,
			}}}},
		}},
	}}
	b, err := proto.Marshal(tx)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(b)
}

func TestExpectedExerciseCheck_ComparesTheChoiceArgument(t *testing.T) {
	t.Parallel()

	const (
		none  = `{"optional":{}}`
		some5 = `{"optional":{"value":{"int64":"5"}}}`
		cid   = `{"contractId":"00c2"}`
		someX = `{"optional":{"value":{"text":"X"}}}`
		some1 = `{"optional":{"value":{"int64":"1"}}}`
	)
	textMap := func(k1, v1, k2, v2 string) string {
		return `{"textMap":{"entries":[{"key":"` + k1 + `","value":{"text":"` + v1 + `"}},{"key":"` + k2 + `","value":{"text":"` + v2 + `"}}]}}`
	}
	for name, tc := range map[string]struct {
		want, got string
		matches   bool
	}{
		"positional None field against a labelled None field": {
			want:    `{"record":{"fields":[{"value":` + none + `},{"value":` + some5 + `}]}}`,
			got:     `{"record":{"fields":[{"label":"a","value":` + none + `},{"label":"b","value":` + some5 + `}]}}`,
			matches: true,
		},
		"omitted labelled None field": {
			want:    `{"record":{"fields":[{"label":"cid","value":` + cid + `}]}}`,
			got:     `{"record":{"fields":[{"label":"note","value":` + none + `},{"label":"cid","value":` + cid + `}]}}`,
			matches: true,
		},
		"positional field missing against a labelled None field": {
			want:    `{"record":{"fields":[{"value":` + some5 + `}]}}`,
			got:     `{"record":{"fields":[{"label":"a","value":` + none + `},{"label":"b","value":` + some5 + `}]}}`,
			matches: false,
		},
		"value moved to another label": {
			want:    `{"record":{"fields":[{"label":"newOperator","value":` + someX + `}]}}`,
			got:     `{"record":{"fields":[{"label":"newOwner","value":` + someX + `},{"label":"newOperator","value":` + none + `}]}}`,
			matches: false,
		},
		"labelled fields in another order": {
			want:    `{"record":{"fields":[{"label":"a","value":` + some1 + `},{"label":"b","value":` + some5 + `}]}}`,
			got:     `{"record":{"fields":[{"label":"b","value":` + some5 + `},{"label":"a","value":` + some1 + `}]}}`,
			matches: true,
		},
		"text map entries in another key order": {
			want:    textMap("b", "1", "a", "2"),
			got:     textMap("a", "2", "b", "1"),
			matches: true,
		},
		"text map record values in another key order": {
			want: `{"textMap":{"entries":[` +
				`{"key":"b","value":{"record":{"fields":[{"label":"x","value":` + cid + `}]}}},` +
				`{"key":"a","value":{"record":{"fields":[{"value":` + none + `},{"value":` + some5 + `}]}}}]}}`,
			got: `{"textMap":{"entries":[` +
				`{"key":"a","value":{"record":{"fields":[{"label":"p","value":` + none + `},{"label":"q","value":` + some5 + `}]}}},` +
				`{"key":"b","value":{"record":{"fields":[{"label":"note","value":` + none + `},{"label":"x","value":` + cid + `}]}}}]}}`,
			matches: true,
		},
		"text map entries with other values": {
			want:    textMap("a", "1", "b", "2"),
			got:     textMap("a", "2", "b", "1"),
			matches: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := ledger.ExpectedExercise{
				PackageID: "pkg", TemplateModule: "M", TemplateEntity: "T",
				ContractID: "00c1", Choice: "Do", ChoiceArgument: tc.want,
			}
			err := e.Check(preparedExerciseB64(t, e, tc.got))

			if tc.matches {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ledger.ErrPreparedTxMismatch)
			}
		})
	}
}

func TestExpectedExerciseCheck_AcceptsAChoiceExercisedThroughTheExpectedInterface(t *testing.T) {
	t.Parallel()

	e := ledger.ExpectedExercise{
		PackageID: "pkg", TemplateModule: "M", TemplateEntity: "Iface",
		ContractID: "00c1", Choice: "Do", ChoiceArgument: `{"unit":{}}`,
	}
	concrete := &apiv2.Identifier{PackageId: "concretePkg", ModuleName: "M", EntityName: "Concrete"}
	for name, tc := range map[string]struct {
		interfaceID *apiv2.Identifier
		matches     bool
	}{
		"interface matches":   {interfaceID: &apiv2.Identifier{PackageId: "pkg", ModuleName: "M", EntityName: "Iface"}, matches: true},
		"interface differs":   {interfaceID: &apiv2.Identifier{PackageId: "pkg", ModuleName: "M", EntityName: "Other"}, matches: false},
		"no interface is set": {interfaceID: nil, matches: false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := e.Check(preparedViaB64(t, e, `{"unit":{}}`, concrete, tc.interfaceID, ""))

			if tc.matches {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ledger.ErrPreparedTxMismatch)
			}
		})
	}
}

func TestExpectedExerciseCheck_PackageNameRef(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		module, entity string
		interfaceID    *apiv2.Identifier
		matches        bool
	}{
		"interface named by the interface's package": {
			module: "I", entity: "Iface",
			interfaceID: &apiv2.Identifier{PackageId: "pkgB", ModuleName: "I", EntityName: "Iface"},
			matches:     true,
		},
		"template named by another package": {module: "M", entity: "T", matches: false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := ledger.ExpectedExercise{
				PackageID: "#iface-pkg", TemplateModule: tc.module, TemplateEntity: tc.entity,
				ContractID: "00c1", Choice: "Do", ChoiceArgument: `{"unit":{}}`,
			}
			template := &apiv2.Identifier{PackageId: "pkgA", ModuleName: "M", EntityName: "T"}
			err := e.Check(preparedViaB64(t, e, `{"unit":{}}`, template, tc.interfaceID, "amulet"))

			if tc.matches {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ledger.ErrPreparedTxMismatch)
			}
		})
	}
}
