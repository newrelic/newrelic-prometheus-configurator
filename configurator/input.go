package configurator

// Input represents the input configuration.
type Input struct {
	Name             string                  `yaml:"name,omitempty"` // It will be used as prometheus_service url parameter, TODO: check which value shall we use.
	RemoteWrite      RemoteWriteInput        `yaml:"newrelic_remote_write"`
	ExtraRemoteWrite []PrometheusExtraConfig `yaml:"extra_remote_write"`
}

// PrometheusExtraConfig represents some configuration which will be included in prometheus as it is.
type PrometheusExtraConfig interface{}
