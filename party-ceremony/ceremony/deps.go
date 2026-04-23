package ceremony

import (
	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/client"
	"github.com/smartcontractkit/chainlink-deployments-framework/pkg/logger"
)

// CantonDeps is the dependency container passed to every operation handler.
type CantonDeps struct {
	Client    client.CantonClient
	Logger    logger.Logger
	Confirmer Confirmer // nil means no confirmation prompt
}
