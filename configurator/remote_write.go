package configurator

import (
	"fmt"
	"time"

	"github.com/newrelic/infrastructure-agent/pkg/license"
)

const (
	remoteWriteBaseURL         = "https://%smetric-api.%snewrelic.com/prometheus/v1/write"
	environmentStagingPrefix   = "staging-"
	regionEUPrefix             = "eu."
	prometheusServerQueryParam = "prometheus_server"
)

// RemoteWriteInput defines all the NewRelic's remote write endpoint fields.
type RemoteWriteInput struct {
	LicenseKey               string                  `yaml:"license_key"`
	Staging                  bool                    `yaml:"staging"`
	ProxyURL                 string                  `yaml:"proxy_url"`
	TLSConfig                *TLSConfig              `yaml:"tls_config"` // TODO: check if we would like to use this TLSConfig or other notation common with other New Relic products.
	QueueConfig              *QueueConfig            `yaml:"queue_config"`
	RemoteTimeout            time.Duration           `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []PrometheusExtraConfig `yaml:"extra_write_relabel_configs"`
}

type QueueConfig struct {
	Capacity          int           `yaml:"capacity"`
	MaxShards         int           `yaml:"max_shards"`
	MinShards         int           `yaml:"min_shards"`
	MaxSamplesPerSend int           `yaml:"max_samples_per_send"`
	BatchSendDeadLine time.Duration `yaml:"batch_send_deadline"`
	MinBackoff        time.Duration `yaml:"min_backoff"`
	MaxBackoff        time.Duration `yaml:"max_backoff"`
	RetryOnHTTP429    bool          `yaml:"retry_on_http_429"`
}

// TLSConfig represents remote-write tls configuration, `prometheusCommonConfig.TLSConfig` cannot be used directly
// because it does not Marshal to yaml properly.
type TLSConfig struct {
	CAFile             string `yaml:"ca_file,omitempty" json:"ca_file,omitempty"`
	CertFile           string `yaml:"cert_file,omitempty" json:"cert_file,omitempty"`
	KeyFile            string `yaml:"key_file,omitempty" json:"key_file,omitempty"`
	ServerName         string `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" json:"insecure_skip_verify"`
	MinVersion         string `yaml:"min_version,omitempty" json:"min_version,omitempty"`
}

// RemoteWriteOutput represents a prometheus remote_write config which can be obtained from input.
type RemoteWriteOutput struct {
	// TODO: check if Name field is needed.
	URL                 string                  `yaml:"url"`
	RemoteTimeout       time.Duration           `yaml:"remote_timeout,omitempty"`
	Authorization       Authorization           `yaml:"authorization"`
	TLSConfig           *TLSConfig              `yaml:"tls_config,omitempty"`
	ProxyURL            string                  `yaml:"proxy_url,omitempty"`
	QueueConfig         *QueueConfig            `yaml:"queue_config,omitempty"`
	WriteRelabelConfigs []PrometheusExtraConfig `yaml:"write_relabel_configs,omitempty"`
}

// Authorization holds prometheus authorization information for remote write.
type Authorization struct {
	// TODO: check if any other authorization option may be used.
	Credentials string `yaml:"credentials"`
}

// BuildRemoteWriteOutput builds a RemoteWriteOutput given the input.
func BuildRemoteWriteOutput(i *Input) RemoteWriteOutput {
	return RemoteWriteOutput{
		URL:                 remoteWriteURL(i.RemoteWrite.Staging, i.RemoteWrite.LicenseKey, i.DataSourceName),
		RemoteTimeout:       i.RemoteWrite.RemoteTimeout,
		Authorization:       Authorization{Credentials: i.RemoteWrite.LicenseKey},
		TLSConfig:           i.RemoteWrite.TLSConfig,
		ProxyURL:            i.RemoteWrite.ProxyURL,
		QueueConfig:         i.RemoteWrite.QueueConfig,
		WriteRelabelConfigs: i.RemoteWrite.ExtraWriteRelabelConfigs,
	}
}

func remoteWriteURL(staging bool, licenseKey string, dataSourceName string) string {
	envPrefix, regionPrefix := "", ""
	if license.IsRegionEU(licenseKey) {
		regionPrefix = regionEUPrefix
	}
	if staging {
		envPrefix = environmentStagingPrefix
	}
	url := fmt.Sprintf(remoteWriteBaseURL, envPrefix, regionPrefix)
	if dataSourceName != "" {
		url = url + "?prometheus_server=" + dataSourceName
	}
	return url
}
