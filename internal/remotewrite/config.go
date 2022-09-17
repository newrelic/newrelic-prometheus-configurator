// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

import (
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

// Name identifies the newrelic remote_write entry used in logs and prometheus self metrics.
const Name = "newrelic_conf"

// Config defines all the NewRelic's remote write endpoint fields.
type Config struct {
	// LicenseKey holds the New Relic ingest license key of the account where metrics will be sent.
	LicenseKey string `yaml:"license_key"`
	// Staging configures the remote write url to point to the New Relic staging endpoint.
	Staging bool `yaml:"staging"`
	// DataSourceName holds the source name which will be used as `prometheus_server` parameter in New Relic remote
	// write endpoint. See:
	// <https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration/>
	// for details.
	DataSourceName string `yaml:"data_source_name"`
	// FedRAMP configures the remote write url to point to the New Relic FedRAMP endpoint.
	FedRAMP                  FedRAMP                 `yaml:"fedramp"`
	ProxyURL                 string                  `yaml:"proxy_url"`
	TLSConfig                *promcfg.TLSConfig      `yaml:"tls_config"`
	QueueConfig              *promcfg.QueueConfig    `yaml:"queue_config"`
	RemoteTimeout            time.Duration           `yaml:"remote_timeout"`
	ExtraWriteRelabelConfigs []promcfg.RelabelConfig `yaml:"extra_write_relabel_configs"`
}

// FedRAMP in charts are configured like `.fedramp.enabled: true` just in case we have to
// add more options to fedramp dictionary. So we add a strict for it.
type FedRAMP struct {
	Enabled bool `yaml:"enabled"`
}

// Build will create the Prometheus remote_write entry for NewRelic.
func (c Config) Build() (promcfg.RemoteWrite, error) {
	rwu := NewURL(
		WithFedRAMP(c.FedRAMP.Enabled),
		WithLicense(c.LicenseKey),
		WithStaging(c.Staging),
		WithDataSourceName(c.DataSourceName),
	)

	url, err := rwu.Build()
	if err != nil {
		return promcfg.RemoteWrite{}, err
	}

	rw := promcfg.RemoteWrite{
		Name:                Name,
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
