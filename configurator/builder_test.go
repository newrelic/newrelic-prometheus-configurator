// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/remotewrite"
	prometheusConfig "github.com/prometheus/prometheus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	// kubernetes module needed to be imported in order to support 'kubernetes_sd_configs' field in prometheus scrape
	// configs. see: <https://github.com/prometheus/prometheus/tree/main/discovery> for details.
	_ "github.com/prometheus/prometheus/discovery/kubernetes"
)

func TestBuilder(t *testing.T) { //nolint: paralleltest,tparallel
	t.Setenv(configurator.LicenseKeyEnvKey, "")
	t.Setenv(configurator.DataSourceNameEnvKey, "")

	// it relies on testdata/<placeholder>.yaml and testdata/<placeholder>.expected.yaml
	testCases := []string{
		"remote-write-test",
		"static-targets-test",
		"external-labels-test",
		"endpoints-test",
		"filter-test",
		"pods-test",
	}

	for _, c := range testCases {
		name := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			inputFile := "testdata/" + name + ".yaml"
			expectedFile := "testdata/" + name + ".expected.yaml"

			data, err := ioutil.ReadFile(inputFile)
			require.NoError(t, err)
			expected, err := ioutil.ReadFile(expectedFile)
			require.NoError(t, err)
			input := &configurator.Input{}
			err = yaml.Unmarshal(data, input)
			require.NoError(t, err)
			output, err := configurator.BuildOutput(input)
			require.NoError(t, err)
			outputData, err := yaml.Marshal(output)
			require.NoError(t, err)

			assertYamlOutputsAreEqual(t, expected, outputData)
			assertIsPrometheusConfig(t, outputData)
		})
	}
}

func TestDataSourceName(t *testing.T) { //nolint: tparallel
	t.Setenv(configurator.LicenseKeyEnvKey, "")
	t.Setenv(configurator.DataSourceNameEnvKey, "")

	configWithDataSourceName := configurator.Input{
		DataSourceName: "test",
		RemoteWrite: remotewrite.Config{
			LicenseKey: "fake",
		},
	}

	//nolint: paralleltest // need clean env variables.
	t.Run("IsSetFromConfig", func(t *testing.T) {
		configWithDataSourceName.DataSourceName = "prom-instance-name"
		prometheusConfig, err := configurator.BuildOutput(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(prometheusConfig)
		require.Contains(t, string(data), fmt.Sprintf("prometheus_server=%s", "prom-instance-name"))
	})

	expectedName := "prom-instance-name-from-env"
	t.Setenv(configurator.DataSourceNameEnvKey, expectedName)

	t.Run("IsSetFromEnvVar", func(t *testing.T) {
		t.Parallel()

		prometheusConfig, err := configurator.BuildOutput(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(prometheusConfig)
		require.Contains(t, string(data), fmt.Sprintf("prometheus_server=%s", expectedName))
	})
}

func TestLicenseKey(t *testing.T) { //nolint: tparallel
	t.Setenv(configurator.LicenseKeyEnvKey, "")
	t.Setenv(configurator.DataSourceNameEnvKey, "")

	configWithDataSourceName := configurator.Input{
		RemoteWrite: remotewrite.Config{},
	}

	//nolint: paralleltest // need clean env variables.
	t.Run("FailIfNotSet", func(t *testing.T) {
		_, err := configurator.BuildOutput(&configWithDataSourceName)
		require.ErrorIs(t, err, configurator.ErrNoLicenseKeyFound)
	})

	//nolint: paralleltest // need clean env variables.
	t.Run("IsSetFromConfig", func(t *testing.T) {
		configWithDataSourceName.RemoteWrite.LicenseKey = "fake"
		prometheusConfig, err := configurator.BuildOutput(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(prometheusConfig)
		require.Contains(t, string(data), fmt.Sprintf("credentials: %s", "fake"))
	})

	expectedLicenseKey := "license-key-from-env"
	t.Setenv(configurator.LicenseKeyEnvKey, expectedLicenseKey)

	t.Run("IsSetFromEnvVar", func(t *testing.T) {
		t.Parallel()

		prometheusConfig, err := configurator.BuildOutput(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(prometheusConfig)
		require.Contains(t, string(data), fmt.Sprintf("credentials: %s", expectedLicenseKey))
	})
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
	require.NoError(t, err, "file content was %s", string(y))
}
