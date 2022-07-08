package configurator

// Input represents the input configuration.
type Input struct {
	// DataSourceName holds the source name which will be used as `prometheus_server` parameter in New Relic remote write endpoint.
	DataSourceName string `yaml:"data_source_name"`
	// RemoteWrite holds the New Relic remote write configuration.
	RemoteWrite RemoteWriteInput `yaml:"newrelic_remote_write"`
	// ExtraRemoteWrite holds any additional remote write configuration to use as it is in prometheus configuration.
	ExtraRemoteWrite []PrometheusExtraConfig `yaml:"extra_remote_write"`
}

// PrometheusExtraConfig represents some configuration which will be included in prometheus as it is.
type PrometheusExtraConfig interface{}
