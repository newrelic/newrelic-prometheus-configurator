package remotewrite

import (
	"testing"

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
			Expected:       "https://metric-api.newrelic.com/prometheus/v1/write?prometheus_server=source",
			DataSourceName: "source",
		},
	}

	for _, testCase := range cases {
		c := testCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			rwu := NewURL(
				WithFedRAMP(c.FedRAMP),
				WithLicense(c.LicenseKey),
				WithStaging(c.Staging),
				WithDataSourceName(c.DataSourceName),
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
			Expected: ErrFedRAMPStaging,
		},
		{
			Name:       "European FedRAMP",
			Staging:    false,
			FedRAMP:    true,
			LicenseKey: "eu-license-key",
			Expected:   ErrEuFedRAMP,
		},
	}

	for _, testCase := range cases {
		c := testCase
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			rwu := NewURL(
				WithFedRAMP(c.FedRAMP),
				WithLicense(c.LicenseKey),
				WithStaging(c.Staging),
				WithDataSourceName(c.DataSourceName),
			)
			_, err := rwu.Build()
			assert.Equal(t, c.Expected, err)
		})
	}
}
