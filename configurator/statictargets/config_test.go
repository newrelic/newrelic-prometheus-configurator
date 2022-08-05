// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets_test

import (
	"net/url"
	"testing"

	"github.com/alecthomas/units"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/statictargets"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildStaticTargetsOutput(t *testing.T) {
	t.Parallel()

	trueValue := true

	cases := []struct {
		Name     string
		Input    statictargets.Config
		Expected []promcfg.Job
	}{
		{
			Name: "All fields set",
			Input: statictargets.Config{
				Jobs: []statictargets.Job{
					{
						Job: promcfg.Job{
							JobName:               "fancy-job",
							HonorLabels:           true,
							HonorTimestamps:       &trueValue,
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
							TLSConfig: &promcfg.TLSConfig{
								CAFile:             "ca-file",
								CertFile:           "cert-file",
								KeyFile:            "key-file",
								ServerName:         "server.name",
								InsecureSkipVerify: true,
								MinVersion:         "TLS12",
							},
							BasicAuth: nil,
							Authorization: promcfg.Authorization{
								Type:            "Bearer",
								Credentials:     "aaa",
								CredentialsFile: "a/b",
							},
							OAuth2: promcfg.OAuth2{
								ClientID:         "client",
								ClientSecret:     "secret",
								ClientSecretFile: "a/secret",
								Scopes:           []string{"all"},
								TokenURL:         "a-url",
								EndpointParams:   map[string]string{"param": "value"},
								TLSConfig: &promcfg.TLSConfig{
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
						ExtraRelabelConfigs: []promcfg.ExtraConfig{
							map[string]any{
								"source_labels": []any{"src.label"},
								"regex":         "to_drop.*",
								"action":        "drop",
							},
						},
						ExtraMetricRelabelConfigs: []promcfg.ExtraConfig{
							map[string]any{
								"source_labels": []any{"src.label"},
								"regex":         "to_drop.*",
								"action":        "drop",
							},
						},
					},
				},
			},
			Expected: []promcfg.Job{
				{
					JobName:               "fancy-job",
					HonorLabels:           true,
					HonorTimestamps:       &trueValue,
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
					TLSConfig: &promcfg.TLSConfig{
						CAFile:             "ca-file",
						CertFile:           "cert-file",
						KeyFile:            "key-file",
						ServerName:         "server.name",
						InsecureSkipVerify: true,
						MinVersion:         "TLS12",
					},
					BasicAuth: nil,
					Authorization: promcfg.Authorization{
						Type:            "Bearer",
						Credentials:     "aaa",
						CredentialsFile: "a/b",
					},
					OAuth2: promcfg.OAuth2{
						ClientID:         "client",
						ClientSecret:     "secret",
						ClientSecretFile: "a/secret",
						Scopes:           []string{"all"},
						TokenURL:         "a-url",
						EndpointParams:   map[string]string{"param": "value"},
						TLSConfig: &promcfg.TLSConfig{
							CAFile:             "ca-file",
							CertFile:           "cert-file",
							KeyFile:            "key-file",
							ServerName:         "server.name",
							InsecureSkipVerify: true,
							MinVersion:         "TLS12",
						},
						ProxyURL: "",
					},

					StaticConfigs: []promcfg.StaticConfig{
						{
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
						},
					},

					RelabelConfigs: []any{
						map[string]any{
							"source_labels": []any{"src.label"},
							"regex":         "to_drop.*",
							"action":        "drop",
						},
					},
					MetricRelabelConfigs: []any{
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
			output := c.Input.Build()
			assert.EqualValues(t, c.Expected, output)
		})
	}
}
