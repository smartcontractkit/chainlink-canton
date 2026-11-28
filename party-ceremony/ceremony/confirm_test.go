package ceremony_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
)

func TestNoOpConfirmer(t *testing.T) {
	t.Parallel()
	c := ceremony.NoOpConfirmer{}
	require.NoError(t, c.ConfirmTopologySign(context.Background(), ceremony.TopologySignDetail{}))
	require.NoError(t, c.ConfirmDAMLSign(context.Background(), ceremony.DAMLSignDetail{}))
}

func TestAlwaysRejectConfirmer(t *testing.T) {
	t.Parallel()
	c := ceremony.AlwaysRejectConfirmer{}
	err := c.ConfirmTopologySign(context.Background(), ceremony.TopologySignDetail{})
	require.ErrorIs(t, err, ceremony.ErrUserRejected)
	err = c.ConfirmDAMLSign(context.Background(), ceremony.DAMLSignDetail{})
	require.ErrorIs(t, err, ceremony.ErrUserRejected)
}

func TestInteractiveConfirmer_Y(t *testing.T) {
	t.Parallel()
	in := bytes.NewBufferString("y\n")
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: in, Out: out}
	err := c.ConfirmTopologySign(context.Background(), ceremony.TopologySignDetail{
		MappingType: "DecentralizedNamespaceDefinition",
		DNSOwners:   []string{"fp-a"},
	})
	require.NoError(t, err)
	assert.Contains(t, out.String(), "TOPOLOGY TRANSACTION")
}

func TestInteractiveConfirmer_Yes(t *testing.T) {
	t.Parallel()
	in := bytes.NewBufferString("yes\n")
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: in, Out: out}
	err := c.ConfirmDAMLSign(context.Background(), ceremony.DAMLSignDetail{
		TransactionHash: "deadbeef",
	})
	require.NoError(t, err)
	assert.Contains(t, out.String(), "DAML TRANSACTION")
}

func TestInteractiveConfirmer_N(t *testing.T) {
	t.Parallel()
	in := bytes.NewBufferString("n\n")
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: in, Out: out}
	err := c.ConfirmTopologySign(context.Background(), ceremony.TopologySignDetail{})
	require.ErrorIs(t, err, ceremony.ErrUserRejected)
}

func TestInteractiveConfirmer_EOF(t *testing.T) {
	t.Parallel()
	in := bytes.NewReader(nil)
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: in, Out: out}
	err := c.ConfirmTopologySign(context.Background(), ceremony.TopologySignDetail{})
	require.ErrorIs(t, err, ceremony.ErrUserRejected)
}

func TestInteractiveConfirmer_DAMLSignShowsTheExercise(t *testing.T) {
	t.Parallel()
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: bytes.NewBufferString("y\n"), Out: out}
	err := c.ConfirmDAMLSign(context.Background(), ceremony.DAMLSignDetail{
		TransactionHash: "deadbeef",
		Template:        "Registry.Service:RegistrarService",
		ContractID:      "registrar-cid",
		Choice:          "RegistrarService_CreateInstrumentConfiguration",
		ChoiceArgument:  `{"record":{}}`,
	})
	require.NoError(t, err)
	for _, want := range []string{"Registry.Service:RegistrarService", "registrar-cid", "RegistrarService_CreateInstrumentConfiguration", `{"record":{}}`} {
		assert.Contains(t, out.String(), want)
	}
}

func TestInteractiveConfirmer_DAMLSignWithoutAChoiceOmitsTheExercise(t *testing.T) {
	t.Parallel()
	out := &bytes.Buffer{}
	c := &ceremony.InteractiveConfirmer{In: bytes.NewBufferString("y\n"), Out: out}
	require.NoError(t, c.ConfirmDAMLSign(context.Background(), ceremony.DAMLSignDetail{TransactionHash: "deadbeef"}))
	for _, label := range []string{"Template:", "Contract ID:", "Choice:", "Choice Argument:"} {
		assert.NotContains(t, out.String(), label)
	}
}
