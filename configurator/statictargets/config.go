// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets

import "github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"

// Config defines all the static targets jobs.
type Config struct {
	Jobs []Job `yaml:"jobs"`
}

// Job represents job config for configurator.
type Job struct {
	Job                       promcfg.Job             `yaml:",inline"`
	Targets                   []string                `yaml:"targets"`
	Labels                    map[string]string       `yaml:"labels"`
	ExtraRelabelConfigs       []promcfg.RelabelConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_metric_relabel_config"`
}

// Build will create a Prometheus Job list based on the static targets configuration.
func (c Config) Build() []promcfg.Job {
	staticTargetsOutput := []promcfg.Job{}

	for _, job := range c.Jobs {
		jobOutput := job.Job

		jobOutput.StaticConfigs = []promcfg.StaticConfig{
			{
				Targets: job.Targets,
				Labels:  job.Labels,
			},
		}

		jobOutput.RelabelConfigs = append(jobOutput.RelabelConfigs, job.ExtraRelabelConfigs...)

		jobOutput.MetricRelabelConfigs = append(jobOutput.MetricRelabelConfigs, job.ExtraMetricRelabelConfigs...)

		staticTargetsOutput = append(staticTargetsOutput, jobOutput)
	}

	return staticTargetsOutput
}
