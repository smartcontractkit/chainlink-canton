package ledger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/ceremony"
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
)

// DARLoader fetches the raw bytes of a DAR by package name and version.
// In production, wire [contracts.GetDar] from github.com/smartcontractkit/chainlink-canton/contracts.
// In tests, a simple map or fake loader can be used.
type DARLoader func(packageName, version string) ([]byte, error)

// releaseDir must stay in sync with contracts.ReleaseDir in the parent module.
const releaseDir = "v1_0_0"

func darVersionDir(version string) string {
	if version == "current" {
		return "current"
	}

	return releaseDir
}

// FileDARLoader returns a [DARLoader] that reads DARs from a directory on the
// local filesystem. Files live under versioned subdirectories (e.g. current/,
// v1_0_0/) matching the embedded FS layout in the contracts package.
func FileDARLoader(dir string) DARLoader {
	return func(packageName, version string) ([]byte, error) {
		path := filepath.Join(dir, darVersionDir(version), fmt.Sprintf("%s-%s.dar", packageName, version))
		data, err := os.ReadFile(path)
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
	// Signer signs prepared transaction hashes for the interactive submission.
	Signer client.TransactionSigner
	// SignerFactory creates a signer from current topology signing keys. It is
	// used in production; Signer remains available for focused tests.
	SignerFactory client.TransactionSignerFactory
	Logger        logger.Logger
	Confirmer     ceremony.Confirmer // nil means no confirmation prompt
	// UserID is the Ledger API user to grant actAs/readAs rights for the
	// decentralized party. When empty the grant step is a no-op (no-auth environments).
	UserID string
}
