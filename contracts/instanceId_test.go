package contracts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstanceId(t *testing.T) {
	t.Parallel()

	t.Run("valid InstanceID", func(t *testing.T) {
		t.Parallel()

		prefix := "onramp"
		party := "ledger-api-user"

		instanceID, err := NewInstanceID(prefix, party)
		require.NoError(t, err)
		require.Len(t, instanceID, len(prefix)+1+InstanceIDRandomPartLength+1+len(party))

		require.True(t, instanceID.Valid())

		gotParty, gotErr := instanceID.Party()
		require.NoError(t, gotErr)
		require.Equal(t, party, gotParty)

		gotPrefix, gotErr := instanceID.Prefix()
		require.NoError(t, gotErr)
		require.Truef(t, strings.HasPrefix(gotPrefix, prefix), "expected prefix to start with %s, got %s", prefix, gotPrefix)

		instanceAddress := instanceID.InstanceAddress()
		require.Len(t, instanceAddress, 32)
		require.NotZero(t, instanceAddress)
	})
	t.Run("invalid prefix (length)", func(t *testing.T) {
		t.Parallel()

		_, err := NewInstanceID("", "some-party")
		require.Error(t, err)

		_, err = NewInstanceID("thisprefixistoolongwaywaywaytoolongasitismorethanfiftycharacters", "some-party")
		require.Error(t, err)
	})
	t.Run("invalid prefix (invalid chars)", func(t *testing.T) {
		t.Parallel()
		instanceID, err := NewInstanceID("On-Ramp_32", "some-party")
		require.NoError(t, err)
		require.True(t, instanceID.Valid())

		_, err = NewInstanceID("invalid hint", "some-party")
		require.Error(t, err)
		_, err = NewInstanceID("ccv-öäü", "some-party")
		require.Error(t, err)
	})
	t.Run("invalid party", func(t *testing.T) {
		t.Parallel()

		_, err := NewInstanceID("valid", "")
		require.Error(t, err)
	})
	t.Run("invalid InstanceID format", func(t *testing.T) {
		t.Parallel()

		invalidID := InstanceID("invalid-format-id")

		require.False(t, invalidID.Valid())

		_, err := invalidID.Party()
		require.Error(t, err)

		_, err = invalidID.Prefix()
		require.Error(t, err)
	})
}
