package ccip

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// NewBundle creates an operations bundle for deployment sequences.
func NewBundle(getCtx func() context.Context, log logger.Logger) cld_ops.Bundle {
	return cld_ops.NewBundle(getCtx, log, cld_ops.NewMemoryReporter())
}
