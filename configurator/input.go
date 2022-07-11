// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

// Input represents the input configuration, it it used to load New Relic configuration so it can be parsed to
// prometheus configuration.
type Input struct {
	// DataSourceName holds the source name which will be used as `prometheus_server` parameter in New Relic remote
	// write endpoint.
	DataSourceName string `yaml:"data_source_name"`
	// RemoteWrite holds the New Relic remote write configuration.
	RemoteWrite RemoteWriteInput `yaml:"newrelic_remote_write"`
	// ExtraRemoteWrite holds any additional remote write configuration to use as it is in prometheus configuration.
	ExtraRemoteWrite []PrometheusExtraConfig `yaml:"extra_remote_write"`
}

// PrometheusExtraConfig represents some configuration which will be included in prometheus as it is.
type PrometheusExtraConfig interface{}
