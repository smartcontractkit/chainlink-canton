package ledger_test

import (
	"encoding/base64"
	"testing"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	v1 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive/transaction/v1"
	"github.com/stretchr/testify/require"
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

func preparedExercise(t *testing.T, exercise *v1.Exercise) string {
	t.Helper()
	b, err := proto.Marshal(&interactive.PreparedTransaction{Transaction: &interactive.DamlTransaction{
		Roots: []string{"0"},
		Nodes: []*interactive.DamlTransaction_Node{{
			NodeId:        "0",
			VersionedNode: &interactive.DamlTransaction_Node_V1{V1: &v1.Node{NodeType: &v1.Node_Exercise{Exercise: exercise}}},
		}},
	}})
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(b)
}

func TestExpectedExerciseCheck_MatchesTheTemplateOrTheInterface(t *testing.T) {
	t.Parallel()

	expected := &ledger.ExpectedExercise{
		TemplateModule: "Splice.Api.Iface",
		TemplateEntity: "Iface",
		ContractID:     "c",
		Choice:         "C",
		ChoiceArgument: `{"record":{"fields":[{"label":"cid","value":{"contractId":"request-cid"}}]}}`,
	}
	concrete := &apiv2.Identifier{ModuleName: "Impl", EntityName: "Concrete"}
	for name, tc := range map[string]struct {
		interfaceID *apiv2.Identifier
		signs       bool
	}{
		"exercised through the expected interface": {interfaceID: &apiv2.Identifier{ModuleName: "Splice.Api.Iface", EntityName: "Iface"}, signs: true},
		"exercised through another interface":      {interfaceID: &apiv2.Identifier{ModuleName: "Splice.Api.Other", EntityName: "Iface"}, signs: false},
		"no interface and another template":        {interfaceID: nil, signs: false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := expected.Check(preparedExercise(t, &v1.Exercise{
				TemplateId:  concrete,
				InterfaceId: tc.interfaceID,
				ContractId:  "c",
				ChoiceId:    "C",
				ChosenValue: &apiv2.Value{Sum: &apiv2.Value_Record{Record: &apiv2.Record{Fields: []*apiv2.RecordField{
					{Value: &apiv2.Value{Sum: &apiv2.Value_ContractId{ContractId: "request-cid"}}},
				}}}},
			}))

			if tc.signs {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, ledger.ErrPreparedTxMismatch)
			}
		})
	}
}
