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
