// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets

import "github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"

// Input defines all the static targets jobs.
type Input struct {
	Jobs []Job `yaml:"jobs"`
}

// Job represents job config for configurator.
type Job struct {
	Job                       promcfg.Job                     `yaml:",inline"`
	Targets                   []string                        `yaml:"targets"`
	Labels                    map[string]string               `yaml:"labels"`
	ExtraRelabelConfigs       []promcfg.PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

// BuildOutput builds the slice of StaticTargetJobOutput given the input.
func BuildOutput(i Input) []any {
	staticTargetsOutput := make([]any, 0)

	for _, job := range i.Jobs {
		jobOutput := job.Job

		jobOutput.StaticConfigs = []promcfg.StaticConfig{
			{
				Targets: job.Targets,
				Labels:  job.Labels,
			},
		}

		for _, c := range job.ExtraRelabelConfigs {
			jobOutput.RelabelConfigs = append(jobOutput.RelabelConfigs, c)
		}

		for _, c := range job.ExtraMetricRelabelConfigs {
			jobOutput.MetricRelabelConfigs = append(jobOutput.MetricRelabelConfigs, c)
		}

		staticTargetsOutput = append(staticTargetsOutput, jobOutput)
	}

	return staticTargetsOutput
}
