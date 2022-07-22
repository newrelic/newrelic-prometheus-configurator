// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"time"
)

// StaticTargetsInput defines all the static targets jobs.
type StaticTargetsInput struct {
	Jobs []Job `yaml:"jobs"`
}

// Job represents a static target job config.
type Job struct {
	Name                      string                  `yaml:"name"`
	Urls                      []string                `yaml:"urls"`
	MetricsPath               string                  `yaml:"metrics_path"`
	Labels                    map[string]string       `yaml:"labels"`
	ScrapeInterval            time.Duration           `yaml:"scrape_interval"`
	ScrapeTimeout             time.Duration           `yaml:"scrape_timeout"`
	TLSConfig                 *TLSConfig              `yaml:"tls_config"`
	BasicAuth                 *BasicAuth              `yaml:"basic_auth"`
	Authorization             Authorization           `yaml:"authorization"`
	OAuth2                    OAuth2                  `yaml:"oauth2"`
	ExtraRelabelConfigs       []PrometheusExtraConfig `yaml:"extra_relabel_config"`
	ExtraMetricRelabelConfigs []PrometheusExtraConfig `yaml:"extra_metric_relabel_config"`
}

// StaticConfigOutput defines each of the static_configs for the prometheus config.
type StaticConfigOutput struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// StaticTargetsJobOutput represents a prometheus scrape_config Job config with static_configs which can be obtained from input.
type StaticTargetsJobOutput struct {
	JobName              string                  `yaml:"job_name"`
	StaticConfigs        []StaticConfigOutput    `yaml:"static_configs,omitempty"`
	ScrapeInterval       time.Duration           `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout        time.Duration           `yaml:"scrape_timeout,omitempty"`
	MetricsPath          string                  `yaml:"metrics_path,omitempty"`
	TLSConfig            *TLSConfig              `yaml:"tls_config,omitempty"`
	BasicAuth            *BasicAuth              `yaml:"basic_auth,omitempty"`
	Authorization        Authorization           `yaml:"authorization,omitempty"`
	OAuth2               OAuth2                  `yaml:"oauth2,omitempty"`
	RelabelConfigs       []PrometheusExtraConfig `yaml:"relabel_configs,omitempty"`
	MetricRelabelConfigs []PrometheusExtraConfig `yaml:"metric_relabel_configs,omitempty"`
}

// BuildStaticTargetsOutput builds the slice of StaticTargetJobOutput given the input.
func BuildStaticTargetsOutput(i *Input) []StaticTargetsJobOutput {
	staticTargetsOutput := make([]StaticTargetsJobOutput, 0)
	for _, job := range i.StaticTargets.Jobs {
		jobOutput := StaticTargetsJobOutput{
			JobName: job.Name,
			StaticConfigs: []StaticConfigOutput{
				{
					Targets: job.Urls,
					Labels:  job.Labels,
				},
			},
			ScrapeInterval:       job.ScrapeInterval,
			ScrapeTimeout:        job.ScrapeTimeout,
			MetricsPath:          job.MetricsPath,
			TLSConfig:            job.TLSConfig,
			BasicAuth:            job.BasicAuth,
			OAuth2:               job.OAuth2,
			Authorization:        job.Authorization,
			RelabelConfigs:       job.ExtraRelabelConfigs,
			MetricRelabelConfigs: job.ExtraMetricRelabelConfigs,
		}
		staticTargetsOutput = append(staticTargetsOutput, jobOutput)
	}

	return staticTargetsOutput
}
