// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator_test

import (
	"net/url"
	"testing"

	"github.com/alecthomas/units"
	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildStaticTargetsOutput(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Input    *configurator.Input
		Expected []any
	}{
		{
			Name: "All fields set",
			Input: &configurator.Input{
				DataSourceName: "source-of-metrics",
				StaticTargets: configurator.StaticTargetsInput{
					Jobs: []configurator.ScrapeJobInput{
						{
							// nolint: dupl // TargetJob should be the same
							TargetJob: configurator.TargetJob{
								JobName:               "fancy-job",
								HonorLabels:           true,
								HonorTimestamps:       true,
								Params:                url.Values{"q": {"puppies"}, "oe": {"utf8"}},
								Scheme:                "https",
								BodySizeLimit:         units.Base2Bytes(1025),
								SampleLimit:           uint(2000),
								TargetLimit:           uint(2000),
								LabelLimit:            uint(2000),
								LabelNameLengthLimit:  uint(2000),
								LabelValueLengthLimit: uint(2000),
								MetricsPath:           "/metrics",
								ScrapeInterval:        10000,
								ScrapeTimeout:         10000,
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
							},
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
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
			Expected: []any{
				configurator.TargetJobOutput{
					// nolint: dupl // TargetJob should be the same.
					TargetJob: configurator.TargetJob{
						JobName:               "fancy-job",
						HonorLabels:           true,
						HonorTimestamps:       true,
						Params:                url.Values{"q": {"puppies"}, "oe": {"utf8"}},
						Scheme:                "https",
						BodySizeLimit:         units.Base2Bytes(1025),
						SampleLimit:           uint(2000),
						TargetLimit:           uint(2000),
						LabelLimit:            uint(2000),
						LabelNameLengthLimit:  uint(2000),
						LabelValueLengthLimit: uint(2000),
						MetricsPath:           "/metrics",
						ScrapeInterval:        10000,
						ScrapeTimeout:         10000,
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
					},
					StaticConfigs: []configurator.StaticConfig{
						{
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
						},
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
