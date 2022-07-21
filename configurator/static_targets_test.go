// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator_test

import (
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildStaticTargetsOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Input    *configurator.Input
		Expected []configurator.StaticTargetsJobOutput
	}{
		{
			Name: "Staging, eu and all fields set",
			Input: &configurator.Input{
				DataSourceName: "source-of-metrics",
				StaticTargets: configurator.StaticTargetsInput{
					Jobs: []configurator.Job{
						{
							Name:           "fancy-job",
							Urls:           []string{"host:port"},
							MetricsPath:    "/metrics",
							Labels:         map[string]string{"a": "b"},
							ScrapeInterval: 10000,
							ScrapeTimeout:  10000,
							TLSConfig: &configurator.TLSConfig{
								CAFile:             "ca-file",
								CertFile:           "cert-file",
								KeyFile:            "key-file",
								ServerName:         "server.name",
								InsecureSkipVerify: true,
								MinVersion:         "TLS12",
							},
							BasicAuth: nil,
							Authorization: configurator.Authorization{
								Type:            "Bearer",
								Credentials:     "aaa",
								CredentialsFile: "a/b",
							},
							OAuth2: configurator.OAuth2{
								ClientID:         "client",
								ClientSecret:     "secret",
								ClientSecretFile: "a/secret",
								Scopes:           []string{"all"},
								TokenURL:         "a-url",
								EndpointParams:   map[string]string{"param": "value"},
								TLSConfig: &configurator.TLSConfig{
									CAFile:             "ca-file",
									CertFile:           "cert-file",
									KeyFile:            "key-file",
									ServerName:         "server.name",
									InsecureSkipVerify: true,
									MinVersion:         "TLS12",
								},
								ProxyURL: "",
							},
							ExtraRelabelConfigs: []configurator.PrometheusExtraConfig{
								map[string]any{
									"source_labels": []any{"src.label"},
									"regex":         "to_drop.*",
									"action":        "drop",
								},
							},
							ExtraMetricRelabelConfigs: []configurator.PrometheusExtraConfig{
								map[string]any{
									"source_labels": []any{"src.label"},
									"regex":         "to_drop.*",
									"action":        "drop",
								},
							},
						},
					},
				},
			},
			Expected: []configurator.StaticTargetsJobOutput{
				{
					JobName: "fancy-job",
					StaticConfigs: []configurator.StaticConfigOutput{
						{
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
						},
					},
					MetricsPath:    "/metrics",
					ScrapeInterval: 10000,
					ScrapeTimeout:  10000,
					TLSConfig: &configurator.TLSConfig{
						CAFile:             "ca-file",
						CertFile:           "cert-file",
						KeyFile:            "key-file",
						ServerName:         "server.name",
						InsecureSkipVerify: true,
						MinVersion:         "TLS12",
					},
					BasicAuth: nil,
					Authorization: configurator.Authorization{
						Type:            "Bearer",
						Credentials:     "aaa",
						CredentialsFile: "a/b",
					},
					OAuth2: configurator.OAuth2{
						ClientID:         "client",
						ClientSecret:     "secret",
						ClientSecretFile: "a/secret",
						Scopes:           []string{"all"},
						TokenURL:         "a-url",
						EndpointParams:   map[string]string{"param": "value"},
						TLSConfig: &configurator.TLSConfig{
							CAFile:             "ca-file",
							CertFile:           "cert-file",
							KeyFile:            "key-file",
							ServerName:         "server.name",
							InsecureSkipVerify: true,
							MinVersion:         "TLS12",
						},
						ProxyURL: "",
					},
					RelabelConfigs: []configurator.PrometheusExtraConfig{
						map[string]any{
							"source_labels": []any{"src.label"},
							"regex":         "to_drop.*",
							"action":        "drop",
						},
					},
					MetricRelabelConfigs: []configurator.PrometheusExtraConfig{
						map[string]any{
							"source_labels": []any{"src.label"},
							"regex":         "to_drop.*",
							"action":        "drop",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			output := configurator.BuildStaticTargetsOutput(c.Input)
			assert.EqualValues(t, c.Expected, output)
		})
	}
}
