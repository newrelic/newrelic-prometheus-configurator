// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"os"
	"testing"
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/remotewrite"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNrConfig(t *testing.T) {
	t.Parallel()

	expected := testNrConfigExpectation(t, false)
	nrConfigData, err := os.ReadFile("testdata/nr-config-test.yaml")
	require.NoError(t, err)

	checkNrConfig(t, expected, nrConfigData)

	expected = testNrConfigExpectation(t, true)
	nrConfigData, err = os.ReadFile("testdata/nr-config-test-proxyfromenv.yaml")
	require.NoError(t, err)
	checkNrConfig(t, expected, nrConfigData)
}

func checkNrConfig(t *testing.T, expected NrConfig, nrConfigData []byte) {
	t.Helper()

	nrConfig := NrConfig{}
	err := yaml.Unmarshal(nrConfigData, &nrConfig)
	require.NoError(t, err)
	require.EqualValues(t, expected, nrConfig)
}

func testNrConfigExpectation(t *testing.T, useProxyFromEnv bool) NrConfig {
	t.Helper()

	trueValue := true
	falseValue := false

	remoteWriteConfig := remotewrite.Config{
		DataSourceName: "data-source",
		LicenseKey:     "nrLicenseKey",
		Staging:        true,
		TLSConfig: &promcfg.TLSConfig{
			InsecureSkipVerify: &trueValue,
			CAFile:             "/path/to/ca.crt",
			CertFile:           "/path/to/cert.crt",
			KeyFile:            "/path/to/key.crt",
			ServerName:         "server.name",
			MinVersion:         "TLS12",
		},
		QueueConfig: &promcfg.QueueConfig{
			Capacity:          2500,
			MaxShards:         200,
			MinShards:         1,
			MaxSamplesPerSend: 500,
			BatchSendDeadLine: 5 * time.Second,
			MinBackoff:        30 * time.Millisecond,
			MaxBackoff:        5 * time.Second,
			RetryOnHTTP429:    &falseValue,
			SampleAgeLimit:    22 * time.Second,
		},
		RemoteTimeout: 30 * time.Second,
		ExtraWriteRelabelConfigs: []promcfg.RelabelConfig{
			{
				SourceLabels: []string{"__name__", "instance"},
				Regex:        "node_memory_active_bytes;localhost:9100",
				Action:       "drop",
			},
		},
	}

	// Decide whether to use ProxyURL or ProxyFromEnvironment
	if useProxyFromEnv {
		remoteWriteConfig.ProxyFromEnvironment = true
	} else {
		remoteWriteConfig.ProxyURL = "http://proxy.url.to.use:1234"
	}

	return NrConfig{
		Common: promcfg.GlobalConfig{
			ScrapeInterval: time.Second * 60,
			ScrapeTimeout:  time.Second,
			ExternalLabels: map[string]string{
				"one":   "two",
				"three": "four",
			},
		},
		RemoteWrite: remoteWriteConfig,
		ExtraRemoteWrite: []RawPromConfig{
			map[string]any{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}
