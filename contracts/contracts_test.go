package contracts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmbeddedArtifacts(t *testing.T) {
	t.Parallel()
	// This test checks that all embedded artifacts do actually exist by iterating over them and checking that the corresponding .dar file exists
	for pkg, versions := range Versions {
		for _, v := range versions {
			content, err := GetDar(pkg, v)
			require.NoErrorf(t, err, "failed to get DAR for package %s version %s", pkg, v)
			require.NotEmptyf(t, content, "DAR content is empty for package %s version %s", pkg, v)
		}
	}
}
