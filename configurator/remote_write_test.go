// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildRemoteWriteOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Input    *Input
		Expected RemoteWriteOutput
	}{
		{
			Name: "Prod,  non-eu and only mandatory fields",
			Input: &Input{
				RemoteWrite: RemoteWriteInput{
					LicenseKey: "fake-prod",
				},
			},
			Expected: RemoteWriteOutput{
				URL: "https://metric-api.newrelic.com/prometheus/v1/write",
				Authorization: Authorization{
					Credentials: "fake-prod",
				},
			},
		},
		{
			Name: "Staging, eu and all fields set",
			Input: &Input{
				DataSourceName: "source-of-metrics",
				RemoteWrite: RemoteWriteInput{
					LicenseKey: "eu-fake-staging",
					Staging:    true,
					ProxyURL:   "http://proxy.url",
					TLSConfig: &TLSConfig{
						CAFile:             "ca-file",
						CertFile:           "cert-file",
						KeyFile:            "key-file",
						ServerName:         "server.name",
						InsecureSkipVerify: true,
						MinVersion:         "TLS12",
					},
					QueueConfig: &QueueConfig{
						Capacity:          100,
						MaxShards:         10,
						MinShards:         2,
						MaxSamplesPerSend: 1000,
						BatchSendDeadLine: time.Second,
						MinBackoff:        100 * time.Microsecond,
						MaxBackoff:        time.Second,
						RetryOnHTTP429:    true,
					},
					RemoteTimeout: 10 * time.Second,
					ExtraWriteRelabelConfigs: []PrometheusExtraConfig{
						map[string]interface{}{
							"source_labels": []interface{}{"src.label"},
							"regex":         "to_drop.*",
							"action":        "drop",
						},
					},
				},
			},
			Expected: RemoteWriteOutput{
				URL:           "https://staging-metric-api.eu.newrelic.com/prometheus/v1/write?prometheus_server=source-of-metrics",
				RemoteTimeout: 10 * time.Second,
				Authorization: Authorization{
					Credentials: "eu-fake-staging",
				},
				ProxyURL: "http://proxy.url",
				TLSConfig: &TLSConfig{
					CAFile:             "ca-file",
					CertFile:           "cert-file",
					KeyFile:            "key-file",
					ServerName:         "server.name",
					InsecureSkipVerify: true,
					MinVersion:         "TLS12",
				},
				QueueConfig: &QueueConfig{
					Capacity:          100,
					MaxShards:         10,
					MinShards:         2,
					MaxSamplesPerSend: 1000,
					BatchSendDeadLine: time.Second,
					MinBackoff:        100 * time.Microsecond,
					MaxBackoff:        time.Second,
					RetryOnHTTP429:    true,
				},
				WriteRelabelConfigs: []PrometheusExtraConfig{
					map[string]interface{}{
						"source_labels": []interface{}{"src.label"},
						"regex":         "to_drop.*",
						"action":        "drop",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			output := BuildRemoteWriteOutput(c.Input)
			assert.EqualValues(t, c.Expected, output)
		})
	}
}

func TestRemoteWriteURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name           string
		Staging        bool
		LicenseKey     string
		Expected       string
		DataSourceName string
	}{
		{
			Name:       "staging non-eu",
			Staging:    true,
			LicenseKey: "non-eu-license-key",
			Expected:   "https://staging-metric-api.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "staging eu",
			Staging:    true,
			LicenseKey: "eu-license-key",
			Expected:   "https://staging-metric-api.eu.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "prod non-eu",
			Staging:    false,
			LicenseKey: "non-eu-license-key",
			Expected:   "https://metric-api.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "prod -eu",
			Staging:    false,
			LicenseKey: "eu-license-key",
			Expected:   "https://metric-api.eu.newrelic.com/prometheus/v1/write",
		},
		{
			Name:           "dataSourceName",
			Staging:        false,
			LicenseKey:     "non-eu-license-key",
			Expected:       "https://metric-api.newrelic.com/prometheus/v1/write?prometheus_server=source",
			DataSourceName: "source",
		},
	}

	for _, testCase := range cases {
		c := testCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := remoteWriteURL(c.Staging, c.LicenseKey, c.DataSourceName)
			assert.Equal(t, c.Expected, result)
		})
	}
}
