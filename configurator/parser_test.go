// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator_test

import (
	"io/ioutil"
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	prometheusConfig "github.com/prometheus/prometheus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestParser(t *testing.T) {
	t.Parallel()

	// it relies on testdata/<placeholder>.yaml and testdata/<placeholder>.expected.yaml
	testCases := []string{
		"remote-write-test",
		"static-targets-test",
		"external-labels-test",
	}

	for _, c := range testCases {
		name := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			inputFile := "testdata/" + name + ".yaml"
			expectedFile := "testdata/" + name + ".expected.yaml"
			input, err := ioutil.ReadFile(inputFile)
			require.NoError(t, err)
			expected, err := ioutil.ReadFile(expectedFile)
			require.NoError(t, err)
			output, err := configurator.Parse(input)
			require.NoError(t, err)
			assertYamlOutputsAreEqual(t, expected, output)
			assertIsPrometheusConfig(t, output)
		})
	}
}

func TestParserInvalidInputYamlError(t *testing.T) {
	t.Parallel()

	input := []byte(`}invalid-yml`)
	_, err := configurator.Parse(input)
	assert.Error(t, err)
}

func assertYamlOutputsAreEqual(t *testing.T, y1, y2 []byte) {
	t.Helper()

	var o1, o2 configurator.Output

	require.NoError(t, yaml.Unmarshal(y1, &o1))
	require.NoError(t, yaml.Unmarshal(y2, &o2))
	assert.EqualValues(t, o1, o2)
}

func assertIsPrometheusConfig(t *testing.T, y []byte) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "gen-prometheus-config")
	require.NoError(t, err)
	_, err = tmpFile.Write(y)
	require.NoError(t, err)
	_, err = prometheusConfig.LoadFile(tmpFile.Name(), true, false, nil)
	require.NoError(t, err)
}
