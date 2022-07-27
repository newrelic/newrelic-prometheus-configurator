// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	yamlEncoderIndent    = 2
	LicenseKeyEnvKey     = "NR_PROM_LICENSE_KEY"
)

var ErrNoLicenseKeyFound = fmt.Errorf(
	"licenseKey was not set neither in yaml config or %s environment variable", LicenseKeyEnvKey,
)

// Parse loads a yaml input and returns the corresponding prometheus-agent yaml.
func Parse(newrelicConfig []byte) ([]byte, error) {
	input := &Input{}
	if err := yaml.Unmarshal(newrelicConfig, input); err != nil {
		return nil, fmt.Errorf("yaml input could not be loaded: %w", err)
	}

	if err := expandConfigsFromEnvVars(input); err != nil {
		return nil, fmt.Errorf("expanding env vars: %w", err)
	}

	output, err := BuildOutput(input)
	if err != nil {
		return nil, fmt.Errorf("output could not be built: %w", err)
	}

	prometheusConfig, err := toYaml(&output)
	if err != nil {
		return nil, fmt.Errorf("output could not be encoded to yaml: %w", err)
	}

	return prometheusConfig, nil
}

func expandConfigsFromEnvVars(i *Input) error {
	if licenseKey := os.Getenv(LicenseKeyEnvKey); licenseKey != "" {
		i.RemoteWrite.LicenseKey = licenseKey
	}

	if i.RemoteWrite.LicenseKey == "" {
		return ErrNoLicenseKeyFound
	}


	return nil
}

func toYaml(output *Output) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(yamlEncoderIndent)

	if err := encoder.Encode(output); err != nil {
		return nil, fmt.Errorf("could not encode to yaml %w", err)
	}

	return buffer.Bytes(), nil
}
