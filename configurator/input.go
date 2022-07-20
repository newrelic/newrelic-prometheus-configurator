// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

// Input represents the input configuration, it it used to load New Relic configuration so it can be parsed to
// prometheus configuration.
type Input struct {
	// DataSourceName holds the source name which will be used as `prometheus_server` parameter in New Relic remote
	// write endpoint. See:
	// <https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration/>
	// for details.
	DataSourceName string `yaml:"data_source_name"`
	// RemoteWrite holds the New Relic remote write configuration.
	RemoteWrite RemoteWriteInput `yaml:"newrelic_remote_write"`
	// ExtraRemoteWrite holds any additional remote write configuration to use as it is in prometheus configuration.
	ExtraRemoteWrite []PrometheusExtraConfig `yaml:"extra_remote_write"`
	// StaticTargets holds the static-target jobs configuration.
	StaticTargets StaticTargetsInput `yaml:"static_targets"`
}

// TLSConfig represents tls configuration, `prometheusCommonConfig.TLSConfig` cannot be used directly
// because it does not Marshal to yaml properly.
type TLSConfig struct {
	CAFile             string `yaml:"ca_file,omitempty" json:"ca_file,omitempty"`
	CertFile           string `yaml:"cert_file,omitempty" json:"cert_file,omitempty"`
	KeyFile            string `yaml:"key_file,omitempty" json:"key_file,omitempty"`
	ServerName         string `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
	MinVersion         string `yaml:"min_version,omitempty" json:"min_version,omitempty"`
}

// Authorization holds prometheus authorization information.
type Authorization struct {
	Type            string `yaml:"type,omitempty"`
	Credentials     string `yaml:"credentials,omitempty"`
	CredentialsFile string `yaml:"credentials_file,omitempty"`
}

// PrometheusExtraConfig represents some configuration which will be included in prometheus as it is.
type PrometheusExtraConfig any
