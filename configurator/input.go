package configurator

// Input represents the input configuration.
type Input struct {
	DataSourceName   string                  `yaml:"data_source_name"`
	RemoteWrite      RemoteWriteInput        `yaml:"newrelic_remote_write"`
	ExtraRemoteWrite []PrometheusExtraConfig `yaml:"extra_remote_write"`
}

// PrometheusExtraConfig represents some configuration which will be included in prometheus as it is.
type PrometheusExtraConfig interface{}
