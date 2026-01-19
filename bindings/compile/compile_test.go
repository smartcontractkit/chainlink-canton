package compile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-canton-internal/contracts"
)

func TestCompileAllPackages(t *testing.T) {
	t.Parallel()
	for c := range contracts.Contracts {
		t.Run(fmt.Sprintf("Compiling %s", c), func(t *testing.T) {
			t.Parallel()
			out, err := Package(c)
			require.NoError(t, err)
			require.NotEmpty(t, out)
		})
	}
}
