// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"fmt"
	"time"
)

const (
	remoteWriteBaseURL       = "https://%smetric-api.%snewrelic.com/prometheus/v1/write"
	environmentStagingPrefix = "staging-"
	regionEUPrefix           = "eu."
	// prometheusServerQueryParam is added to remoteWrite url when input's name is defined.
	prometheusServerQueryParam = "prometheus_server"
)

// RemoteWriteInput defines all the NewRelic's remote write endpoint fields.
type RemoteWriteInput struct {
	LicenseKey               string                  `yaml:"license_key"`
	Staging                  bool                    `yaml:"staging"`
	ProxyURL                 string                  `yaml:"proxy_url"`
	TLSConfig                *TLSConfig              `yaml:"tls_config"`
	QueueConfig              *QueueConfig            `yaml:"queue_config"`
	RemoteTimeout            time.Duration           `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []PrometheusExtraConfig `yaml:"extra_write_relabel_configs"`
}

// QueueConfig represents the remote-write queue config.
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

// RemoteWriteOutput represents a prometheus remote_write config which can be obtained from input.
type RemoteWriteOutput struct {
	URL                 string                  `yaml:"url"`
	RemoteTimeout       time.Duration           `yaml:"remote_timeout,omitempty"`
	Authorization       Authorization           `yaml:"authorization"`
	TLSConfig           *TLSConfig              `yaml:"tls_config,omitempty"`
	ProxyURL            string                  `yaml:"proxy_url,omitempty"`
	QueueConfig         *QueueConfig            `yaml:"queue_config,omitempty"`
	WriteRelabelConfigs []PrometheusExtraConfig `yaml:"write_relabel_configs,omitempty"`
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
	if licenseIsRegionEU(licenseKey) {
		regionPrefix = regionEUPrefix
	}

	if staging {
		envPrefix = environmentStagingPrefix
	}

	url := fmt.Sprintf(remoteWriteBaseURL, envPrefix, regionPrefix)
	if dataSourceName != "" {
		url = url + "?" + prometheusServerQueryParam + "=" + dataSourceName
	}

	return url
}
