// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"io/ioutil"
	"testing"
	"time"

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

//nolint:paralleltest
// t.SetEnv cannot be used in parallel tests.
func TestInputLicenseKeyFromEnv(t *testing.T) {
	expectedLicenseKey := "license-key-from-env"
	expected := testInputExpectation(t)
	expected.RemoteWrite.LicenseKey = expectedLicenseKey
	t.Setenv("NR_PROM_LICENSE_KEY", expectedLicenseKey)
	inputData, err := ioutil.ReadFile("testdata/input-test.no-license-key.yaml")
	require.NoError(t, err)

	checkInput(t, expected, inputData)
}

//nolint:paralleltest
// t.SetEnv cannot be used in parallel tests.
func TestInputLicenseKeyFromEnvPrecedence(t *testing.T) {
	expectedLicenseKey := "license-key-from-env"
	expected := testInputExpectation(t)
	expected.RemoteWrite.LicenseKey = expectedLicenseKey
	t.Setenv("NR_PROM_LICENSE_KEY", expectedLicenseKey)
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
		DataSourceName: "data-source",
		RemoteWrite: RemoteWriteInput{
			LicenseKey: "nrLicenseKey",
			Staging:    true,
			ProxyURL:   "http://proxy.url.to.use:1234",
			TLSConfig: &TLSConfig{
				InsecureSkipVerify: true,
				CAFile:             "/path/to/ca.crt",
				CertFile:           "/path/to/cert.crt",
				KeyFile:            "/path/to/key.crt",
				ServerName:         "server.name",
				MinVersion:         "TLS12",
			},
			QueueConfig: &QueueConfig{
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
			ExtraWriteRelabelConfigs: []PrometheusExtraConfig{
				map[string]interface{}{
					"source_labels": []interface{}{"__name__", "instance"},
					"regex":         "node_memory_active_bytes;localhost:9100",
					"action":        "drop",
				},
			},
		},
		ExtraRemoteWrite: []PrometheusExtraConfig{
			map[string]interface{}{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
}
