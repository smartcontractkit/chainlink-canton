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
	var chosen apiv2.Value
	require.NoError(t, protojson.Unmarshal([]byte(chosenJSON), &chosen))
	tx := &interactive.PreparedTransaction{Transaction: &interactive.DamlTransaction{
		Roots: []string{"0"},
		Nodes: []*interactive.DamlTransaction_Node{{
			NodeId: "0",
			VersionedNode: &interactive.DamlTransaction_Node_V1{V1: &v1.Node{NodeType: &v1.Node_Exercise{Exercise: &v1.Exercise{
				TemplateId:  &apiv2.Identifier{PackageId: e.PackageID, ModuleName: e.TemplateModule, EntityName: e.TemplateEntity},
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
		"text map entries in another key order": {
			want:    textMap("b", "1", "a", "2"),
			got:     textMap("a", "2", "b", "1"),
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
