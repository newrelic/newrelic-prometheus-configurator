// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import "github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"

// StaticTargetsInput defines all the static targets jobs.
type StaticTargetsInput struct {
	Jobs []JobInput `yaml:"jobs"`
}

// BuildStaticTargetsOutput builds the slice of StaticTargetJobOutput given the input.
func BuildStaticTargetsOutput(i *Input) []any {
	staticTargetsOutput := make([]any, 0)

	for _, job := range i.StaticTargets.Jobs {
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
