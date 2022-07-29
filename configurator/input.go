// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

// Input represents the input configuration, it it used to load New Relic configuration so it can be parsed to
// prometheus configuration.
type Input struct {
	// Common holds configuration for all options common to all scrape methods.
	Common GlobalConfig `yaml:"common"`
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
	// ExtraScrapeConfigs holds any additional raw scrape configuration to use as it is in prometheus configuration.
	ExtraScrapeConfigs []PrometheusExtraConfig `yaml:"extra_scrape_configs"`
	// Kubernetes holds the kubernetes-targets' configuration.
	Kubernetes KubernetesInput `yaml:"kubernetes"`
}
