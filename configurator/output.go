// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package configurator holds the code to parse New Relic's configuration into a valid prometheus-agent configuration.
package configurator

// Output holds all configuration information in prometheus format which can be directly marshaled to a valid yaml
// configuration.
type Output struct {
	RemoteWrite   []any        `yaml:"remote_write"`
	ScrapeConfigs []any        `yaml:"scrape_configs,omitempty"`
	GlobalConfig  GlobalConfig `yaml:"global"`
}

// BuildOutput builds the prometheus config output from the provided input, it holds "first level" transformations
// required to obtain a valid prometheus configuration.
func BuildOutput(input *Input) (Output, error) {
	output := Output{
		RemoteWrite:  []any{BuildRemoteWriteOutput(input)},
		GlobalConfig: input.Common,
	}

	for _, extraRemoteWriteConfig := range input.ExtraRemoteWrite {
		output.RemoteWrite = append(output.RemoteWrite, extraRemoteWriteConfig)
	}

	// Include the scrape configurations corresponding to static targets
	if staticTargets := BuildStaticTargetsOutput(input); len(staticTargets) > 0 {
		output.ScrapeConfigs = append(output.ScrapeConfigs, staticTargets...)
	}

	// Include the scrape configurations corresponding to kubernetes configuration
	kubernetesTargetBuilder := newKubernetesJobBuilder()
	k8stargets, err := kubernetesTargetBuilder.Build(input)
	if err != nil {
		return output, err
	}
	for _, targetJob := range k8stargets {
		output.ScrapeConfigs = append(output.ScrapeConfigs, targetJob)
	}

	// Include "extra" scrape configuration
	for _, extraScrapeConfig := range input.ExtraScrapeConfigs {
		output.ScrapeConfigs = append(output.ScrapeConfigs, extraScrapeConfig)
	}

	return output, nil
}
