// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/remotewrite"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInput(t *testing.T) {
	t.Parallel()

	expected := testInputExpectation(t)
	inputData, err := ioutil.ReadFile("testdata/input-test.yaml")
	require.NoError(t, err)

	checkInput(t, expected, inputData)
}

func checkInput(t *testing.T, expected Input, inputData []byte) {
	t.Helper()

	input := Input{}
	err := yaml.Unmarshal(inputData, &input)
	require.NoError(t, err)
	require.EqualValues(t, expected, input)
}

func testInputExpectation(t *testing.T) Input {
	t.Helper()

	return Input{
		Common: promcfg.GlobalConfig{
			ScrapeInterval: time.Second * 60,
			ScrapeTimeout:  time.Second,
			ExternalLabels: map[string]string{
				"one":   "two",
				"three": "four",
			},
		},
		DataSourceName: "data-source",
		RemoteWrite: remotewrite.Input{
			LicenseKey: "nrLicenseKey",
			Staging:    true,
			ProxyURL:   "http://proxy.url.to.use:1234",
			TLSConfig: &promcfg.TLSConfig{
				InsecureSkipVerify: true,
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
				RetryOnHTTP429:    false,
			},
			RemoteTimeout: 30 * time.Second,
			ExtraWriteRelabelConfigs: []promcfg.PrometheusExtraConfig{
				map[string]any{
					"source_labels": []any{"__name__", "instance"},
					"regex":         "node_memory_active_bytes;localhost:9100",
					"action":        "drop",
				},
			},
		},
		ExtraRemoteWrite: []promcfg.PrometheusExtraConfig{
			map[string]any{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}
