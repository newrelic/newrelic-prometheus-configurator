// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statictargets_test

import (
	"net/url"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejob"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/statictargets"

	"github.com/alecthomas/units"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestBuildStaticTargetsPromConfig(t *testing.T) {
	t.Parallel()

	trueValue := true

	cases := []struct {
		Name     string
		NrConfig statictargets.Config
		Expected []promcfg.Job
	}{
		{
			Name: "All fields set with ProxyURL",
			NrConfig: statictargets.Config{
				StaticTargetJobs: []statictargets.StaticTargetJob{
					{
						ScrapeJob: scrapejob.Job{
							Job: promcfg.Job{
								JobName:               "fancy-job",
								HonorLabels:           &trueValue,
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
									InsecureSkipVerify: &trueValue,
									MinVersion:         "TLS12",
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
										InsecureSkipVerify: &trueValue,
										MinVersion:         "TLS12",
									},
									ProxyURL: "http://proxy.url.to.use:1234",
								},
							},
						},
						Targets: []string{"host:port"},
						Labels:  map[string]string{"a": "b"},
					},
				},
			},
			Expected: []promcfg.Job{
				{
					JobName:               "fancy-job",
					HonorLabels:           &trueValue,
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
						InsecureSkipVerify: &trueValue,
						MinVersion:         "TLS12",
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
							InsecureSkipVerify: &trueValue,
							MinVersion:         "TLS12",
						},
						ProxyURL: "http://proxy.url.to.use:1234",
					},
					StaticConfigs: []promcfg.StaticConfig{
						{
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
						},
					},
				},
			},
		},
		{
			Name: "All fields set with ProxyFromEnvironment",
			NrConfig: statictargets.Config{
				StaticTargetJobs: []statictargets.StaticTargetJob{
					{
						ScrapeJob: scrapejob.Job{
							Job: promcfg.Job{
								JobName:               "fancy-job",
								HonorLabels:           &trueValue,
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
									InsecureSkipVerify: &trueValue,
									MinVersion:         "TLS12",
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
										InsecureSkipVerify: &trueValue,
										MinVersion:         "TLS12",
									},
									ProxyFromEnvironment: true,
								},
							},
						},
						Targets: []string{"host:port"},
						Labels:  map[string]string{"a": "b"},
					},
				},
			},
			Expected: []promcfg.Job{
				{
					JobName:               "fancy-job",
					HonorLabels:           &trueValue,
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
						InsecureSkipVerify: &trueValue,
						MinVersion:         "TLS12",
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
							InsecureSkipVerify: &trueValue,
							MinVersion:         "TLS12",
						},
						ProxyFromEnvironment: true,
					},
					StaticConfigs: []promcfg.StaticConfig{
						{
							Targets: []string{"host:port"},
							Labels:  map[string]string{"a": "b"},
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
			prometheusConfig := c.NrConfig.Build(sharding.Config{})
			assert.EqualValues(t, c.Expected, prometheusConfig)
		})
	}
}
