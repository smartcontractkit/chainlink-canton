package monitoring

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

// InitBeholderMonitoring initialized a labeler using the given Beholder config.
func InitBeholderMonitoring(config beholder.Config) (*EDSMetricLabeler, error) {
	// All histogram buckets must be defined at time of creation.
	config.MetricViews = MetricViews()

	client, err := beholder.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create beholder client: %w", err)
	}

	beholder.SetClient(client)
	beholder.SetGlobalOtelProviders()

	edsMetrics, err := InitMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EDS metrics: %w", err)
	}

	return NewEDSMetricLabeler(metrics.NewLabeler(), edsMetrics), nil
}
