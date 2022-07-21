// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package configurator holds the code to parse New Relic's configuration into a valid prometheus-agent configuration.
package configurator

// Output holds all configuration information in prometheus format which can be directly marshaled to a valid yaml
// configuration.
type Output struct {
	RemoteWrite   []any                    `yaml:"remote_write"`
	ScrapeConfigs []StaticTargetsJobOutput `yaml:"scrape_configs,omitempty"`
}

// BuildOutput builds the prometheus config output from the provided input, it holds "first level" transformations
// required to obtain a valid prometheus configuration.
func BuildOutput(input *Input) (Output, error) {
	output := Output{
		RemoteWrite: []any{BuildRemoteWriteOutput(input)},
	}

	if staticTargets := BuildStaticTargetsOutput(input); len(staticTargets) > 0 {
		output.ScrapeConfigs = append(output.ScrapeConfigs, staticTargets...)
	}

	// Include extra remote-write configs
	for _, extraRemoteWriteConfig := range input.ExtraRemoteWrite {
		output.RemoteWrite = append(output.RemoteWrite, extraRemoteWriteConfig)
	}

	return output, nil
}
