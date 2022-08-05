// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite_test

import (
	"testing"
	"time"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/remotewrite"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildRemoteWriteOutput(t *testing.T) {
	t.Parallel()

	type args struct {
		remoteConfig   remotewrite.Config
		dataSourceName string
	}

	cases := []struct {
		Name     string
		Input    args
		Expected promcfg.RemoteWrite
	}{
		{
			Name: "Prod,  non-eu and only mandatory fields",
			Input: args{
				remoteConfig: remotewrite.Config{
					LicenseKey: "fake-prod",
				},
			},
			Expected: promcfg.RemoteWrite{
				URL: "https://metric-api.newrelic.com/prometheus/v1/write",
				Authorization: promcfg.Authorization{
					Credentials: "fake-prod",
				},
			},
		},
		{
			Name: "Staging, eu and all fields set",
			Input: args{
				dataSourceName: "source-of-metrics",
				remoteConfig: remotewrite.Config{
					LicenseKey: "eu-fake-staging",
					Staging:    true,
					ProxyURL:   "http://proxy.url",
					TLSConfig: &promcfg.TLSConfig{
						CAFile:             "ca-file",
						CertFile:           "cert-file",
						KeyFile:            "key-file",
						ServerName:         "server.name",
						InsecureSkipVerify: true,
						MinVersion:         "TLS12",
					},
					QueueConfig: &promcfg.QueueConfig{
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
					ExtraWriteRelabelConfigs: []promcfg.ExtraConfig{
						map[string]any{
							"source_labels": []any{"src.label"},
							"regex":         "to_drop.*",
							"action":        "drop",
						},
					},
				},
			},
			Expected: promcfg.RemoteWrite{
				URL:           "https://staging-metric-api.eu.newrelic.com/prometheus/v1/write?prometheus_server=source-of-metrics",
				RemoteTimeout: 10 * time.Second,
				Authorization: promcfg.Authorization{
					Credentials: "eu-fake-staging",
				},
				ProxyURL: "http://proxy.url",
				TLSConfig: &promcfg.TLSConfig{
					CAFile:             "ca-file",
					CertFile:           "cert-file",
					KeyFile:            "key-file",
					ServerName:         "server.name",
					InsecureSkipVerify: true,
					MinVersion:         "TLS12",
				},
				QueueConfig: &promcfg.QueueConfig{
					Capacity:          100,
					MaxShards:         10,
					MinShards:         2,
					MaxSamplesPerSend: 1000,
					BatchSendDeadLine: time.Second,
					MinBackoff:        100 * time.Microsecond,
					MaxBackoff:        time.Second,
					RetryOnHTTP429:    true,
				},
				WriteRelabelConfigs: []promcfg.ExtraConfig{
					map[string]any{
						"source_labels": []any{"src.label"},
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
			output := c.Input.remoteConfig.Build(c.Input.dataSourceName)
			assert.EqualValues(t, c.Expected, output)
		})
	}
}
