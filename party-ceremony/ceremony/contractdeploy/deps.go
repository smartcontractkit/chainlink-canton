package contractdeploy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chainlink/canton-party-ceremony/internal/client"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
)

// DARLoader fetches the raw bytes of a DAR by package name and version.
// In production, wire [contracts.GetDar] from github.com/smartcontractkit/chainlink-canton/contracts.
// In tests, a simple map or fake loader can be used.
type DARLoader func(packageName, version string) ([]byte, error)

// FileDARLoader returns a [DARLoader] that reads DARs from a directory on the
// local filesystem. Files are expected to follow the naming convention
// "<name>-<version>.dar" (matching the embedded FS layout in the contracts package).
func FileDARLoader(dir string) DARLoader {
	return func(packageName, version string) ([]byte, error) {
		path := filepath.Join(dir, fmt.Sprintf("%s-%s.dar", packageName, version))
		data, err := os.ReadFile(path) //nolint:gosec // path is constructed from known package name/version
		if err != nil {
			return nil, fmt.Errorf("reading DAR %q: %w", path, err)
		}
		return data, nil
	}
}

// ContractDeployDeps is the dependency container passed to every contract-deploy
// operation handler. It holds both the Admin client (for DAR uploads) and the
// Ledger client (for interactive submission and party queries).
type ContractDeployDeps struct {
	AdminClient  client.CantonClient
	LedgerClient client.LedgerClient
	// DARLoader resolves package references to DAR bytes.
	// Wire contracts.GetDar for production usage.
	DARLoader DARLoader
	Logger    logger.Logger
}
