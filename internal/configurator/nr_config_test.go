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

	t.Run("with proxy URL", func(t *testing.T) {
		t.Parallel()
		expected := testNrConfigWithProxyURL(t)
		nrConfigData, err := os.ReadFile("testdata/nr-config-test.yaml")
		require.NoError(t, err)
		checkNrConfig(t, expected, nrConfigData)
	})

	t.Run("with proxy from environment", func(t *testing.T) {
		t.Parallel()
		expected := testNrConfigWithProxyFromEnv(t)
		nrConfigData, err := os.ReadFile("testdata/nr-config-test-proxyfromenv.yaml")
		require.NoError(t, err)
		checkNrConfig(t, expected, nrConfigData)
	})

	t.Run("with full global config", func(t *testing.T) {
		t.Parallel()
		expected := testNrConfigWithFullGlobalConfig(t)
		nrConfigData, err := os.ReadFile("testdata/nr-config-test-globalconfig.yaml")
		require.NoError(t, err)
		checkNrConfig(t, expected, nrConfigData)
	})
}

func checkNrConfig(t *testing.T, expected NrConfig, nrConfigData []byte) {
	t.Helper()

	nrConfig := NrConfig{}
	err := yaml.Unmarshal(nrConfigData, &nrConfig)
	require.NoError(t, err)
	require.EqualValues(t, expected, nrConfig)
}

func baseRemoteWriteConfig() remotewrite.Config {
	trueValue := true
	falseValue := false

	return remotewrite.Config{
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
}

func baseGlobalConfig() promcfg.GlobalConfig {
	return promcfg.GlobalConfig{
		ScrapeInterval: time.Second * 60,
		ScrapeTimeout:  time.Second,
		ExternalLabels: map[string]string{
			"one":   "two",
			"three": "four",
		},
	}
}

func testNrConfigWithProxyURL(t *testing.T) NrConfig {
	t.Helper()

	remoteWriteConfig := baseRemoteWriteConfig()
	remoteWriteConfig.ProxyURL = "http://proxy.url.to.use:1234"

	return NrConfig{
		Common:      baseGlobalConfig(),
		RemoteWrite: remoteWriteConfig,
		ExtraRemoteWrite: []RawPromConfig{
			map[string]any{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}

func testNrConfigWithProxyFromEnv(t *testing.T) NrConfig {
	t.Helper()

	remoteWriteConfig := baseRemoteWriteConfig()
	remoteWriteConfig.ProxyFromEnvironment = true

	return NrConfig{
		Common:      baseGlobalConfig(),
		RemoteWrite: remoteWriteConfig,
		ExtraRemoteWrite: []RawPromConfig{
			map[string]any{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}

func testNrConfigWithFullGlobalConfig(t *testing.T) NrConfig {
	t.Helper()

	remoteWriteConfig := baseRemoteWriteConfig()
	remoteWriteConfig.ProxyURL = "http://proxy.url.to.use:1234"

	globalConfig := baseGlobalConfig()
	globalConfig.ScrapeProtocols = []string{"PrometheusProto", "OpenMetricsText1.0.0"}
	globalConfig.EvaluationInterval = 2 * time.Minute
	globalConfig.RuleQueryOffset = 5 * time.Second
	globalConfig.QueryLogFile = "/var/log/prometheus/query.log"
	globalConfig.ScrapeFailureLogFile = "/var/log/prometheus/scrape-failures.log"
	globalConfig.BodySizeLimit = 10 * 1024 * 1024 // 10MiB
	globalConfig.SampleLimit = 1000
	globalConfig.LabelLimit = 50
	globalConfig.LabelNameLengthLimit = 200
	globalConfig.LabelValueLengthLimit = 500
	globalConfig.TargetLimit = 100
	globalConfig.KeepDroppedTargets = 10
	globalConfig.MetricNameValidationScheme = "utf8"
	globalConfig.ExtraScrapeMetrics = true

	return NrConfig{
		Common:      globalConfig,
		RemoteWrite: remoteWriteConfig,
		ExtraRemoteWrite: []RawPromConfig{
			map[string]any{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}
