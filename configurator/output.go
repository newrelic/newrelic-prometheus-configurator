// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package configurator holds the code to parse New Relic's configuration into a valid prometheus-agent configuration.
package configurator

import (
	"fmt"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
)

// Output holds all configuration information in prometheus format which can be directly marshaled to a valid yaml
// configuration.
type Output struct {
	RemoteWrite   []RawPromConfig      `yaml:"remote_write"`
	ScrapeConfigs []RawPromConfig      `yaml:"scrape_configs,omitempty"`
	GlobalConfig  promcfg.GlobalConfig `yaml:"global"`
}

// BuildOutput builds the prometheus config output from the provided input, it holds "first level" transformations
// required to obtain a valid prometheus configuration.
func BuildOutput(input *Input) (Output, error) {
	output := Output{
		RemoteWrite:  []RawPromConfig{input.RemoteWrite.Build(input.DataSourceName)},
		GlobalConfig: input.Common,
	}

	output.RemoteWrite = append(output.RemoteWrite, input.ExtraRemoteWrite...)

	for _, staticTargets := range input.StaticTargets.Build() {
		output.ScrapeConfigs = append(output.ScrapeConfigs, staticTargets)
	}

	k8sJobs, err := input.Kubernetes.Build()
	if err != nil {
		return output, fmt.Errorf("building k8s config: %w", err)
	}

	for _, job := range k8sJobs {
		output.ScrapeConfigs = append(output.ScrapeConfigs, job)
	}

	output.ScrapeConfigs = append(output.ScrapeConfigs, input.ExtraScrapeConfigs...)

	return output, nil
}
