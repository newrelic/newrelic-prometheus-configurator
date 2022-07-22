// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"time"
)

// CommonConfigInput configures values that are used across other configuration
// objects.
type CommonConfigInput struct {
	GlobalConfig `yaml:",inline"`
}

// CommonConfigOutput configures values that are used across other configuration
// objects.
type CommonConfigOutput struct {
	GlobalConfig `yaml:",inline"`
}

// GlobalConfig configures values that are used across other configuration
// objects.
type GlobalConfig struct {
	// How frequently to scrape targets by default.
	ScrapeInterval time.Duration `yaml:"scrape_interval,omitempty"`
	// The default timeout when scraping targets.
	ScrapeTimeout time.Duration `yaml:"scrape_timeout,omitempty"`
	// How frequently to evaluate rules by default.
	EvaluationInterval time.Duration `yaml:"evaluation_interval,omitempty"`
	// The labels to add to any timeseries that this Prometheus instance scrapes.
	ExternalLabels map[string]string `yaml:"external_labels,omitempty"`
}

// BuildRemoteCommonConfigOutput builds a CommonConfigOutput given the input.
func BuildRemoteCommonConfigOutput(i *Input) CommonConfigOutput {
	return CommonConfigOutput{
		GlobalConfig: i.Common.GlobalConfig,
	}
}
