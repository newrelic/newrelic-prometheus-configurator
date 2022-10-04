package remotewrite_test

import (
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/remotewrite"
	"github.com/stretchr/testify/assert"
)

func TestRemoteWriteURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name           string
		Staging        bool
		FedRAMP        bool
		LicenseKey     string
		Expected       string
		DataSourceName string
		CollectorName  string
	}{
		{
			Name:       "staging non-eu",
			Staging:    true,
			FedRAMP:    false,
			LicenseKey: "non-eu-license-key",
			Expected:   "https://staging-metric-api.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "staging eu",
			Staging:    true,
			FedRAMP:    false,
			LicenseKey: "eu-license-key",
			Expected:   "https://staging-metric-api.eu.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "prod non-eu",
			Staging:    false,
			FedRAMP:    false,
			LicenseKey: "non-eu-license-key",
			Expected:   "https://metric-api.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "prod -eu",
			Staging:    false,
			FedRAMP:    false,
			LicenseKey: "eu-license-key",
			Expected:   "https://metric-api.eu.newrelic.com/prometheus/v1/write",
		},
		{
			Name:       "fedramp",
			Staging:    false,
			FedRAMP:    true,
			LicenseKey: "non-eu-license-key",
			Expected:   "https://gov-metric-api.newrelic.com/prometheus/v1/write",
		},
		{
			Name:           "dataSourceName",
			Staging:        false,
			FedRAMP:        false,
			LicenseKey:     "non-eu-license-key",
			DataSourceName: "source",
			Expected:       "https://metric-api.newrelic.com/prometheus/v1/write?prometheus_server=source",
		},
		{
			Name:          "collectorName",
			Staging:       false,
			FedRAMP:       false,
			LicenseKey:    "non-eu-license-key",
			CollectorName: "foo",
			Expected:      "https://metric-api.newrelic.com/prometheus/v1/write?collector_name=foo",
		},
	}

	for _, testCase := range cases {
		c := testCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			rwu := remotewrite.NewURL(
				remotewrite.WithFedRAMP(c.FedRAMP),
				remotewrite.WithLicense(c.LicenseKey),
				remotewrite.WithStaging(c.Staging),
				remotewrite.WithDataSourceName(c.DataSourceName),
				remotewrite.WithCollectorName(c.CollectorName),
			)
			result, err := rwu.Build()
			if assert.NoError(t, err) {
				assert.Equal(t, c.Expected, result)
			}
		})
	}
}

func TestRemoteWriteURLErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name           string
		Staging        bool
		FedRAMP        bool
		LicenseKey     string
		Expected       error
		DataSourceName string
	}{
		{
			Name:     "staging FedRAMP",
			Staging:  true,
			FedRAMP:  true,
			Expected: remotewrite.ErrFedRAMPStaging,
		},
		{
			Name:       "European FedRAMP",
			Staging:    false,
			FedRAMP:    true,
			LicenseKey: "eu-license-key",
			Expected:   remotewrite.ErrEuFedRAMP,
		},
	}

	for _, testCase := range cases {
		c := testCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			rwu := remotewrite.NewURL(
				remotewrite.WithFedRAMP(c.FedRAMP),
				remotewrite.WithLicense(c.LicenseKey),
				remotewrite.WithStaging(c.Staging),
				remotewrite.WithDataSourceName(c.DataSourceName),
			)
			_, err := rwu.Build()
			assert.Equal(t, c.Expected, err)
		})
	}
}
