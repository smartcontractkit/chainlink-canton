package ledger_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

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
