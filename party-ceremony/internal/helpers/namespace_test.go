package helpers_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"
)

func TestComputeDecentralizedNamespace(t *testing.T) {
	t.Parallel()

	// Calculated using Canton's Console
	expected := "1220f6d3baf6aad6ee26a519cc1be0aef278ef3b7debbfcd0e2f781ec0798b953522"
	n1 := "1220315314ddb36cb2ee73a284c01ce2e70c4d8ef1fd07cd7e7e68f9ac69f9d4d8b0"
	n2 := "1220f7a1b4e384170598f547dfad2e961436cd8a325efe3b1b91d94ce38b2b9c727a"
	n3 := "12203669b0e6cb2d88dfaaeb7409cd996d67202d84f505b01bd188c4d62e590f027d"

	result := helpers.ComputeDecentralizedNamespace([]string{n1, n2, n3})
	assert.Equal(t, expected, result)

	result2 := helpers.ComputeDecentralizedNamespace([]string{n1, n2, n3})
	assert.Equal(t, expected, result2)
}

// reference: https://docs.digitalasset.com/operate/3.4/howtos/operate/parties/decentralized_parties.html#decentralized-namespace-computation
func TestComputeDecentralizedNamespace_MatchesDocs(t *testing.T) {
	t.Parallel()

	h := sha256.New()
	h.Write([]byte{0, 0, 0, 37})
	h.Write([]byte{0, 0, 0, 5})
	h.Write([]byte("alice"))
	h.Write([]byte{0, 0, 0, 3})
	h.Write([]byte("bob"))
	digest := h.Sum(nil)
	expected := fmt.Sprintf("1220%x", digest)

	result := helpers.ComputeDecentralizedNamespace([]string{"bob", "alice"})
	assert.Equal(t, expected, result)
}
