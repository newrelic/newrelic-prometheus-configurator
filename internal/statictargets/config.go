// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets

import (
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejobs"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
)

type Config struct {
	StaticTargetJobs []StaticTargetJob `yaml:"jobs"`
}

// StaticTargetJob represents job config for configurator.
type StaticTargetJob struct {
	ScrapeJob scrapejobs.Job    `yaml:",inline"`
	Targets   []string          `yaml:"targets"`
	Labels    map[string]string `yaml:"labels"`
}

// Build will create a Prometheus Job list based on the static targets configuration.
func (c Config) Build(shardingConfig sharding.Config) []promcfg.Job {
	promScrapeJobs := []promcfg.Job{}

	for _, staticTargetJob := range c.StaticTargetJobs {
		promScrapeJob := staticTargetJob.ScrapeJob.BuildPrometheusJob(shardingConfig)

		promScrapeJob.StaticConfigs = []promcfg.StaticConfig{
			{
				Targets: staticTargetJob.Targets,
				Labels:  staticTargetJob.Labels,
			},
		}

		promScrapeJobs = append(promScrapeJobs, promScrapeJob)
	}

	return promScrapeJobs
}
