// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"github.com/newrelic-forks/newrelic-prometheus/configurator/kubernetes"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/remotewrite"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/statictargets"
)

type RawPromConfig any

// Input represents the input configuration, it is used to load New Relic configuration, so it can be parsed to
// prometheus configuration.
type Input struct {
	// Common holds configuration for all options common to all scrape methods.
	Common promcfg.GlobalConfig `yaml:"common"`
	// DataSourceName holds the source name which will be used as `prometheus_server` parameter in New Relic remote
	// write endpoint. See:
	// <https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration/>
	// for details.
	DataSourceName string `yaml:"data_source_name"`
	// Sharding holds the configuration for the sharding.
	Sharding promcfg.Sharding `yaml:"sharding"`
	// RemoteWrite holds the New Relic remote write configuration.
	RemoteWrite remotewrite.Config `yaml:"newrelic_remote_write"`
	// ExtraRemoteWrite holds any additional remote write configuration to use as it is in prometheus configuration.
	ExtraRemoteWrite []RawPromConfig `yaml:"extra_remote_write"`
	// StaticTargets holds the static-target jobs configuration.
	StaticTargets statictargets.Config `yaml:"static_targets"`
	// ExtraScrapeConfigs holds any additional raw scrape configuration to use as it is in prometheus configuration.
	ExtraScrapeConfigs []RawPromConfig `yaml:"extra_scrape_configs"`
	// Kubernetes holds the kubernetes-targets' configuration.
	Kubernetes kubernetes.Config `yaml:"kubernetes"`
}
