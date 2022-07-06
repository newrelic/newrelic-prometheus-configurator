package configurator

import (
	"crypto/tls"
	"io/ioutil"
	"testing"
	"time"

	prometheusCommonConfig "github.com/prometheus/common/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInput(t *testing.T) {
	expected := Input{
		DataSourceName: "data-source",
		RemoteWrite: RemoteWriteInput{
			LicenseKey: "nrLicenseKey",
			Staging:    true,
			ProxyURL:   "http://proxy.url.to.use:1234",
			TLSConfig: &prometheusCommonConfig.TLSConfig{
				InsecureSkipVerify: true,
				CAFile:             "/path/to/ca.crt",
				CertFile:           "/path/to/cert.crt",
				KeyFile:            "/path/to/key.crt",
				ServerName:         "server.name",
				MinVersion:         tls.VersionTLS12,
			},
			QueueConfig: &QueueConfig{
				Capacity:          2500,
				MaxShards:         200,
				MinShards:         1,
				MaxSamplesPerSend: 500,
				BatchSendDeadLine: 5 * time.Second,
				MinBackoff:        30 * time.Millisecond,
				MaxBackoff:        5 * time.Second,
				RetryOnHTTP429:    false,
			},
			RemoteTimeout: 30 * time.Second,
			ExtraWriteRelabelConfigs: []PrometheusExtraConfig{
				map[string]interface{}{
					"source_labels": []interface{}{"__name__", "instance"},
					"regex":         "node_memory_active_bytes;localhost:9100",
					"action":        "drop",
				},
			},
		},
		ExtraRemoteWrite: []PrometheusExtraConfig{
			map[string]interface{}{
				"url": "https://extra.prometheus.remote.write",
			},
		},
	}
	inputData, err := ioutil.ReadFile("testdata/input-test.yaml")
	require.NoError(t, err)
	input := Input{}
	err = yaml.Unmarshal(inputData, &input)
	require.NoError(t, err)
	require.EqualValues(t, expected, input)
}
