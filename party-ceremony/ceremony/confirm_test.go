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
