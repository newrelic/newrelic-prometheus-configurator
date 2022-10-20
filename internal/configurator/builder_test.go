// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/configurator"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/remotewrite"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
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
		"kubernetes-scrape-fields-test",
		"sharding-test",
		"skip-sharding-test",
	}

	for _, c := range testCases {
		name := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nrConfigFile := "testdata/" + name + ".yaml"
			expectedFile := "testdata/" + name + ".expected.yaml"

			data, err := os.ReadFile(nrConfigFile)
			require.NoError(t, err)
			expected, err := os.ReadFile(expectedFile)
			require.NoError(t, err)
			nrConfig := &configurator.NrConfig{}
			err = yaml.Unmarshal(data, nrConfig)
			require.NoError(t, err)
			prometheusConfig, err := configurator.BuildPromConfig(nrConfig)
			require.NoError(t, err)
			prometheusConfigData, err := yaml.Marshal(prometheusConfig)
			require.NoError(t, err)

			assertYamlPromConfigsAreEqual(t, expected, prometheusConfigData)
			assertIsPrometheusConfig(t, prometheusConfigData)
		})
	}
}

func TestDataSourceName(t *testing.T) { //nolint: tparallel
	t.Setenv(configurator.LicenseKeyEnvKey, "")
	t.Setenv(configurator.DataSourceNameEnvKey, "")

	configWithDataSourceName := configurator.NrConfig{
		RemoteWrite: remotewrite.Config{
			DataSourceName: "test",
			LicenseKey:     "fake",
		},
	}

	//nolint: paralleltest // need clean env variables.
	t.Run("IsSetFromConfig", func(t *testing.T) {
		configWithDataSourceName.RemoteWrite.DataSourceName = "prom-instance-name"
		promConf, err := configurator.BuildPromConfig(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(promConf)
		require.Contains(t, string(data), fmt.Sprintf("prometheus_server=%s", "prom-instance-name"))
	})

	expectedName := "prom-instance-name-from-env"
	t.Setenv(configurator.DataSourceNameEnvKey, expectedName)

	t.Run("IsSetFromEnvVar", func(t *testing.T) {
		t.Parallel()

		promConf, err := configurator.BuildPromConfig(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(promConf)
		require.Contains(t, string(data), fmt.Sprintf("prometheus_server=%s", expectedName))
	})
}

func TestLicenseKey(t *testing.T) { //nolint: tparallel
	t.Setenv(configurator.LicenseKeyEnvKey, "")
	t.Setenv(configurator.DataSourceNameEnvKey, "")

	configWithDataSourceName := configurator.NrConfig{
		RemoteWrite: remotewrite.Config{},
	}

	//nolint: paralleltest // need clean env variables.
	t.Run("FailIfNotSet", func(t *testing.T) {
		_, err := configurator.BuildPromConfig(&configWithDataSourceName)
		require.ErrorIs(t, err, configurator.ErrNoLicenseKeyFound)
	})

	//nolint: paralleltest // need clean env variables.
	t.Run("IsSetFromConfig", func(t *testing.T) {
		configWithDataSourceName.RemoteWrite.LicenseKey = "fake"
		promConf, err := configurator.BuildPromConfig(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(promConf)
		require.Contains(t, string(data), fmt.Sprintf("credentials: %s", "fake"))
	})

	expectedLicenseKey := "license-key-from-env"
	t.Setenv(configurator.LicenseKeyEnvKey, expectedLicenseKey)

	t.Run("IsSetFromEnvVar", func(t *testing.T) {
		t.Parallel()

		promConf, err := configurator.BuildPromConfig(&configWithDataSourceName)
		require.NoError(t, err)

		data, _ := yaml.Marshal(promConf)
		require.Contains(t, string(data), fmt.Sprintf("credentials: %s", expectedLicenseKey))
	})
}

func TestShardingIndex(t *testing.T) { //nolint: paralleltest
	t.Setenv(configurator.LicenseKeyEnvKey, "fake")

	testCases := []struct {
		name     string
		config   configurator.NrConfig
		expected string
		setEnv   func()
	}{
		{
			name: "IsSetFromEnvVar",
			config: configurator.NrConfig{
				Sharding: sharding.Config{
					Kind:             "hash",
					TotalShardsCount: 2,
				},
			},
			expected: "1",
			setEnv: func() {
				t.Setenv(configurator.DataSourceNameEnvKey, "newrelic-prometheus-1")
			},
		},
		{
			name: "IsSetToEmptyWhenInvalidEnvVaR",
			config: configurator.NrConfig{
				Sharding: sharding.Config{
					Kind:             "hash",
					TotalShardsCount: 2,
				},
			},
			expected: "",
			setEnv: func() {
				t.Setenv(configurator.DataSourceNameEnvKey, "invalid_name")
			},
		},
		{
			name: "HonoursConfigWhenInvalidOrEmptyEnvVar",
			config: configurator.NrConfig{
				Sharding: sharding.Config{
					Kind:             "hash",
					TotalShardsCount: 2,
					ShardIndex:       "3",
				},
			},
			expected: "3",
			setEnv: func() {
				t.Setenv(configurator.DataSourceNameEnvKey, "")
			},
		},
	}

	for _, c := range testCases { //nolint: paralleltest
		t.Run(c.name, func(t *testing.T) {
			c.setEnv()

			_, err := configurator.BuildPromConfig(&c.config)
			require.NoError(t, err)

			require.Equal(t, c.expected, c.config.Sharding.ShardIndex)
		})
	}
}

func assertYamlPromConfigsAreEqual(t *testing.T, y1, y2 []byte) {
	t.Helper()

	var o1, o2 configurator.PromConfig

	require.NoError(t, yaml.Unmarshal(y1, &o1))
	require.NoError(t, yaml.Unmarshal(y2, &o2))
	assert.EqualValues(t, o1, o2)
}

func assertIsPrometheusConfig(t *testing.T, y []byte) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "gen-prometheus-config")
	require.NoError(t, err)
	_, err = tmpFile.Write(y)
	require.NoError(t, err)
	_, err = prometheusConfig.LoadFile(tmpFile.Name(), true, false, nil)
	require.NoError(t, err, "file content was %s", string(y))
}
