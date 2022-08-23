// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package configurator holds the code to parse New Relic's configuration into a valid prometheus-agent configuration.
package configurator

import (
	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
)

// Output holds all configuration information in prometheus format which can be directly marshaled to a valid yaml
// configuration.
type Output struct {
	RemoteWrite   []RawPromConfig      `yaml:"remote_write"`
	ScrapeConfigs []RawPromConfig      `yaml:"scrape_configs,omitempty"`
	GlobalConfig  promcfg.GlobalConfig `yaml:"global"`
}
