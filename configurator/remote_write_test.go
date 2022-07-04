package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteWriteParseProd(t *testing.T) {
	prodInput := &Input{
		Name: "prod-source",
		RemoteWrite: RemoteWriteInput{
			Staging:    false,
			LicenseKey: "fake-prod",
		},
	}
	output := BuildRemoteWriteOutput(prodInput)
	assert.Equal(t, "https://metric-api.newrelic.com/prometheus/v1/write?prometheus_server=prod-source", output.URL)
	assert.Equal(t, "fake-prod", output.Authorization.Credentials)
}

func TestRemoteWriteParseStaging(t *testing.T) {
	prodInput := &Input{
		Name: "staging-source",
		RemoteWrite: RemoteWriteInput{
			Staging:    true,
			LicenseKey: "fake-staging",
		},
	}
	output := BuildRemoteWriteOutput(prodInput)
	assert.Equal(t, "https://staging-metric-api.newrelic.com/prometheus/v1/write?prometheus_server=staging-source", output.URL)
	assert.Equal(t, "fake-staging", output.Authorization.Credentials)
}
