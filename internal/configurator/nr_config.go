// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"github.com/newrelic/newrelic-prometheus-configurator/internal/kubernetes"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/remotewrite"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/statictargets"
)

type RawPromConfig any

// NrConfig represents the nrConfig configuration, it is used to load New Relic configuration, so it can be parsed to
// prometheus configuration.
type NrConfig struct {
	// Common holds configuration for all options common to all scrape methods.
	Common promcfg.GlobalConfig `yaml:"common"`
	// Sharding holds the configuration for the sharding.
	Sharding sharding.Config `yaml:"sharding"`
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
