// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import (
	"fmt"
	"time"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
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
	LicenseKey               string                          `yaml:"license_key"`
	Staging                  bool                            `yaml:"staging"`
	ProxyURL                 string                          `yaml:"proxy_url"`
	TLSConfig                *promcfg.TLSConfig              `yaml:"tls_config"`
	QueueConfig              *promcfg.QueueConfig            `yaml:"queue_config"`
	RemoteTimeout            time.Duration                   `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []promcfg.PrometheusExtraConfig `yaml:"extra_write_relabel_configs"`
}

// BuildRemoteWriteOutput builds a RemoteWriteOutput given the input.
func BuildRemoteWriteOutput(i *Input) promcfg.RemoteWriteOutput {
	return promcfg.RemoteWriteOutput{
		URL:                 remoteWriteURL(i.RemoteWrite.Staging, i.RemoteWrite.LicenseKey, i.DataSourceName),
		RemoteTimeout:       i.RemoteWrite.RemoteTimeout,
		Authorization:       promcfg.Authorization{Credentials: i.RemoteWrite.LicenseKey},
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
