package monitoring

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"

	"github.com/smartcontractkit/chainlink-canton/eds/common"
)

type EDSBeholderMonitoring struct {
	metrics common.EDSMetricLabeler
}

func InitMonitoring(config beholder.Config) (common.EDSMonitoring, error) {
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

	return &EDSBeholderMonitoring{
		metrics: NewEDSMetricLabeler(metrics.NewLabeler(), edsMetrics),
	}, nil
}

func (m *EDSBeholderMonitoring) Metrics() common.EDSMetricLabeler {
	return m.metrics
}
