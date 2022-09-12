// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scrapejobs

import "github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"

// Config is a wrapper of `promcfg.Job` to include some extra fields used in all scrape jobs.
type Config struct {
	promcfg.Job `yaml:",inline"`

	ExtraRelabelConfigs       []promcfg.RelabelConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_metric_relabel_config"`
}

// WithRelabelConfigs returns Config with the provided relabel configs added in the underlying `promcfg.Job`.
func (c Config) WithRelabelConfigs(relabelConfigs []promcfg.RelabelConfig) Config {
	c.Job.RelabelConfigs = append(c.Job.RelabelConfigs, relabelConfigs...)
	return c
}

// WithName returns Config with the underlying `promcfg.Job` name updated.
func (c Config) WithName(name string) Config {
	c.Job.JobName = name
	return c
}

// BuildPrometheusJob returns the underlying `promcfg.Job` setting up the extra fields.
func (c Config) BuildPrometheusJob() promcfg.Job {
	c.Job.RelabelConfigs = append(c.Job.RelabelConfigs, c.ExtraRelabelConfigs...)
	c.Job.MetricRelabelConfigs = append(c.Job.MetricRelabelConfigs, c.ExtraMetricRelabelConfigs...)

	return c.Job
}
