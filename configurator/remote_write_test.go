package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteWriteParseProd(t *testing.T) {
	prodInput := &Input{
		RemoteWrite: RemoteWriteInput{
			Staging:    false,
			LicenseKey: "fake-prod",
		},
	}
	output := BuildRemoteWriteOutput(prodInput)
	assert.Equal(t, "https://metric-api.newrelic.com/prometheus/v1/write", output.URL)
	assert.Equal(t, "fake-prod", output.Authorization.Credentials)
}

func TestRemoteWriteParseStaging(t *testing.T) {
	prodInput := &Input{
		RemoteWrite: RemoteWriteInput{
			Staging:    true,
			LicenseKey: "fake-staging",
		},
	}
	output := BuildRemoteWriteOutput(prodInput)
	assert.Equal(t, "https://staging-metric-api.newrelic.com/prometheus/v1/write", output.URL)
	assert.Equal(t, "fake-staging", output.Authorization.Credentials)
}

func TestRemoteWriteURL(t *testing.T) {

	cases := []struct {
		Name       string
		Staging    bool
		LicenseKey string
		Expected   string
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
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			result := remoteWriteURL(c.Staging, c.LicenseKey)
			assert.Equal(t, c.Expected, result)
		})
	}
}
