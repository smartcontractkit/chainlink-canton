package cmd

import (
	"strings"
	"testing"
	"time"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPrintContracts(t *testing.T) {
	t.Parallel()

	var out strings.Builder
	require.NoError(t, printContracts(&out, []*apiv2.CreatedEvent{{
		ContractId: "00abc",
		Offset:     42,
		CreatedAt:  timestamppb.New(time.Date(2026, 11, 13, 1, 2, 3, 0, time.UTC)),
		CreateArguments: &apiv2.Record{Fields: []*apiv2.RecordField{
			{Label: "instrumentId", Value: &apiv2.Value{Sum: &apiv2.Value_Text{Text: "LINK"}}},
		}},
	}}))

	got := out.String()
	assert.Contains(t, got, "Active contracts: 1")
	assert.Contains(t, got, "Contract ID: 00abc")
	assert.Contains(t, got, "Offset:      42")
	assert.Contains(t, got, "Created at:  2026-11-13T01:02:03Z")
	assert.Contains(t, got, `"text":"LINK"`)
}
