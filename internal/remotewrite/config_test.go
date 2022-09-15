// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite_test

import (
	"testing"
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/remotewrite"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildRemoteWritePromConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		remoteConfig remotewrite.Config
	}

	trueValue := true

	cases := []struct {
		Name     string
		NrConfig args
		Expected promcfg.RemoteWrite
	}{
		{
			Name: "Prod,  non-eu and only mandatory fields",
			NrConfig: args{
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
			NrConfig: args{
				remoteConfig: remotewrite.Config{
					DataSourceName: "source-of-metrics",
					LicenseKey:     "eu-fake-staging",
					Staging:        true,
					ProxyURL:       "http://proxy.url",
					TLSConfig: &promcfg.TLSConfig{
						CAFile:             "ca-file",
						CertFile:           "cert-file",
						KeyFile:            "key-file",
						ServerName:         "server.name",
						InsecureSkipVerify: &trueValue,
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
						RetryOnHTTP429:    &trueValue,
					},
					RemoteTimeout: 10 * time.Second,
					ExtraWriteRelabelConfigs: []promcfg.RelabelConfig{
						{
							SourceLabels: []string{"src.label"},
							Regex:        "to_drop.*",
							Action:       "drop",
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
					InsecureSkipVerify: &trueValue,
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
					RetryOnHTTP429:    &trueValue,
				},
				WriteRelabelConfigs: []promcfg.RelabelConfig{
					{
						SourceLabels: []string{"src.label"},
						Regex:        "to_drop.*",
						Action:       "drop",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			prometheusConfig, err := c.NrConfig.remoteConfig.Build()
			assert.NoError(t, err)
			assert.EqualValues(t, c.Expected, prometheusConfig)
		})
	}
}
