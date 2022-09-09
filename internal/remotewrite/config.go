// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

import (
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

// Config defines all the NewRelic's remote write endpoint fields.
type Config struct {
	LicenseKey               string                  `yaml:"license_key"`
	Staging                  bool                    `yaml:"staging"`
	FedRAMP                  bool                    `yaml:"fedramp"`
	ProxyURL                 string                  `yaml:"proxy_url"`
	TLSConfig                *promcfg.TLSConfig      `yaml:"tls_config"`
	QueueConfig              *promcfg.QueueConfig    `yaml:"queue_config"`
	RemoteTimeout            time.Duration           `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_write_relabel_configs"`
}

// Build will create the Prometheus remote_write entry for NewRelic.
func (c Config) Build(dataSourceName string) (promcfg.RemoteWrite, error) {
	rwu := NewURL(
		WithFedRAMP(c.FedRAMP),
		WithLicense(c.LicenseKey),
		WithStaging(c.Staging),
		WithDataSourceName(dataSourceName),
	)

	url, err := rwu.Build()
	if err != nil {
		return promcfg.RemoteWrite{}, err
	}

	rw := promcfg.RemoteWrite{
		URL:                 url,
		RemoteTimeout:       c.RemoteTimeout,
		Authorization:       promcfg.Authorization{Credentials: c.LicenseKey},
		TLSConfig:           c.TLSConfig,
		ProxyURL:            c.ProxyURL,
		QueueConfig:         c.QueueConfig,
		WriteRelabelConfigs: c.ExtraWriteRelabelConfigs,
	}

	return rw, nil
}
