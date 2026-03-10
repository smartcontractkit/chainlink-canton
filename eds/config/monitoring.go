package config

type MonitoringConfig struct {
	Enabled  bool           `toml:"enabled"`
	Beholder BeholderConfig `toml:"beholder" validate:"required_if=Enabled true"`
}

type BeholderConfig struct {
	InsecureConnection       bool    `toml:"insecure_connection"`
	CACertFile               string  `toml:"ca_cert_file"`
	OtelExporterGRPCEndpoint string  `toml:"otel_exporter_grpc_endpoint"`
	OtelExporterHTTPEndpoint string  `toml:"otel_exporter_http_endpoint"`
	LogStreamingEnabled      bool    `toml:"log_streaming_enabled"`
	MetricReaderInterval     int64   `toml:"metric_reader_interval"`
	TraceSampleRatio         float64 `toml:"trace_sample_ratio"`
	TraceBatchTimeout        int64   `toml:"trace_batch_timeout"`
}
