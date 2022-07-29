// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

// StaticTargetsInput defines all the static targets jobs.
type StaticTargetsInput struct {
	Jobs []JobInput `yaml:"jobs"`
}

// BuildStaticTargetsOutput builds the slice of StaticTargetJobOutput given the input.
func BuildStaticTargetsOutput(i *Input) []any {
	staticTargetsOutput := make([]any, 0)
	for _, job := range i.StaticTargets.Jobs {
		jobOutput := BuildJobOutput(job).WithExtraConfigs(job)
		staticTargetsOutput = append(staticTargetsOutput, jobOutput)
	}

	return staticTargetsOutput
}
