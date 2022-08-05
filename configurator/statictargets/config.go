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
	Job                       promcfg.Job           `yaml:",inline"`
	Targets                   []string              `yaml:"targets"`
	Labels                    map[string]string     `yaml:"labels"`
	ExtraRelabelConfigs       []promcfg.ExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.ExtraConfig `yaml:"extra_metric_relabel_config"`
}

// Build builds the slice of StaticTargetJobOutput given the input.
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

		for _, rc := range job.ExtraRelabelConfigs {
			jobOutput.RelabelConfigs = append(jobOutput.RelabelConfigs, rc)
		}

		for _, mrc := range job.ExtraMetricRelabelConfigs {
			jobOutput.MetricRelabelConfigs = append(jobOutput.MetricRelabelConfigs, mrc)
		}

		staticTargetsOutput = append(staticTargetsOutput, jobOutput)
	}

	return staticTargetsOutput
}
