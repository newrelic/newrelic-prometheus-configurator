// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets

import "github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"

// Config defines all the static targets jobs.
type Config struct {
	StaticTargetJobs []StaticTargetJob `yaml:"jobs"`
}

// StaticTargetJob represents job config for configurator.
type StaticTargetJob struct {
	PromScrapeJob             promcfg.Job             `yaml:",inline"`
	Targets                   []string                `yaml:"targets"`
	Labels                    map[string]string       `yaml:"labels"`
	ExtraRelabelConfigs       []promcfg.RelabelConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_metric_relabel_config"`
}

// Build will create a Prometheus Job list based on the static targets configuration.
func (c Config) Build() []promcfg.Job {
	promScrapeJobs := []promcfg.Job{}

	for _, staticTargetJob := range c.StaticTargetJobs {
		promScrapeJob := staticTargetJob.PromScrapeJob

		promScrapeJob.StaticConfigs = []promcfg.StaticConfig{
			{
				Targets: staticTargetJob.Targets,
				Labels:  staticTargetJob.Labels,
			},
		}

		promScrapeJob.RelabelConfigs = append(promScrapeJob.RelabelConfigs, staticTargetJob.ExtraRelabelConfigs...)

		promScrapeJob.MetricRelabelConfigs = append(promScrapeJob.MetricRelabelConfigs, staticTargetJob.ExtraMetricRelabelConfigs...)

		promScrapeJobs = append(promScrapeJobs, promScrapeJob)
	}

	return promScrapeJobs
}
