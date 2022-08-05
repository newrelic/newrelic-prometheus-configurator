// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

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

// Config defines all the NewRelic's remote write endpoint fields.
type Config struct {
	LicenseKey               string                  `yaml:"license_key"`
	Staging                  bool                    `yaml:"staging"`
	ProxyURL                 string                  `yaml:"proxy_url"`
	TLSConfig                *promcfg.TLSConfig      `yaml:"tls_config"`
	QueueConfig              *promcfg.QueueConfig    `yaml:"queue_config"`
	RemoteTimeout            time.Duration           `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_write_relabel_configs"`
}

// Build will create the Prometheus remote_write entry for NewRelic.
func (c Config) Build(dataSourceName string) promcfg.RemoteWrite {
	return promcfg.RemoteWrite{
		URL:                 remoteWriteURL(c.Staging, c.LicenseKey, dataSourceName),
		RemoteTimeout:       c.RemoteTimeout,
		Authorization:       promcfg.Authorization{Credentials: c.LicenseKey},
		TLSConfig:           c.TLSConfig,
		ProxyURL:            c.ProxyURL,
		QueueConfig:         c.QueueConfig,
		WriteRelabelConfigs: c.ExtraWriteRelabelConfigs,
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
