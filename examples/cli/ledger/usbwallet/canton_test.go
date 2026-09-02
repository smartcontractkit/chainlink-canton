package usbwallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSignatureResponse(t *testing.T) {
	t.Parallel()

	signature := make([]byte, 64)
	for i := range signature {
		signature[i] = byte(i)
	}
	challengeSignature := make([]byte, 64)
	for i := range challengeSignature {
		challengeSignature[i] = byte(255 - i)
	}

	t.Run("without challenge", func(t *testing.T) {
		t.Parallel()

		reply := append(append([]byte{0x40}, signature...), 0x00)
		got, gotChallenge, err := parseSignatureResponse(reply, false)
		require.NoError(t, err)
		require.Equal(t, signature, got)
		require.Nil(t, gotChallenge)
	})

	t.Run("with challenge", func(t *testing.T) {
		t.Parallel()

		reply := append(append([]byte{0x40}, signature...), 0x00)
		reply = append(append(reply, 0x40), challengeSignature...)
		got, gotChallenge, err := parseSignatureResponse(reply, true)
		require.NoError(t, err)
		require.Equal(t, signature, got)
		require.Equal(t, challengeSignature, gotChallenge)
	})

	t.Run("truncated", func(t *testing.T) {
		t.Parallel()

		_, _, err := parseSignatureResponse([]byte{0x40, 0x01}, false)
		require.Error(t, err)
	})
}
