// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scrapejob

import (
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
)

// Job is a wrapper of `promcfg.Job` to include some extra fields used in all scrape jobs.
type Job struct {
	promcfg.Job `yaml:",inline"`

	SkipSharding              bool                    `yaml:"skip_sharding"`
	ExtraRelabelConfigs       []promcfg.RelabelConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_metric_relabel_config"`
}

// WithRelabelConfigs returns Config with the provided relabel configs added in the underlying `promcfg.Job`.
func (j Job) WithRelabelConfigs(relabelConfigs []promcfg.RelabelConfig) Job {
	j.RelabelConfigs = append(j.RelabelConfigs, relabelConfigs...)
	return j
}

// WithName returns Config with the underlying `promcfg.Job` name updated.
func (j Job) WithName(name string) Job {
	j.JobName = name
	return j
}

// Includes the sharding rules corresponding to the job configuration and the sharding config provided.
func (j Job) includeShardingRules(shardingConfig sharding.Config, job promcfg.Job) promcfg.Job {
	if shardingConfig.ShouldIncludeShardingRules() && !j.SkipSharding {
		job.RelabelConfigs = append(shardingConfig.RelabelConfigs(), job.RelabelConfigs...)
	}
	return job
}

// BuildPrometheusJob returns the underlying `promcfg.Job` setting up the extra fields and additional rules.
func (j Job) BuildPrometheusJob(shardingConfig sharding.Config) promcfg.Job {
	job := j.includeShardingRules(shardingConfig, j.Job)

	job.RelabelConfigs = append(job.RelabelConfigs, j.ExtraRelabelConfigs...)
	job.MetricRelabelConfigs = append(job.MetricRelabelConfigs, j.ExtraMetricRelabelConfigs...)

	return job
}
