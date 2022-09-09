// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	remoteWriteScheme        = "https"
	remoteWriteHostTemplate  = "%smetric-api.%snewrelic.com"
	remoteWritePath          = "prometheus/v1/write"
	environmentStagingPrefix = "staging-"
	environmentFedRAMPPrefix = "gov-"
	regionEUPrefix           = "eu."
	// prometheusServerQueryParam is added to remoteWrite url when nrConfig's name is defined.
	prometheusServerQueryParam = "prometheus_server"
)

var (
	ErrFedRAMPStaging = errors.New("there is no staging environment for FedRAMP")
	ErrEuFedRAMP      = errors.New("there is no European FedRAMP region")
)

type URLOption func(url *URL)

type URL struct {
	Staging      bool
	FedRAMP      bool
	RegionPrefix string
	Values       url.Values
}

func NewURL(opts ...URLOption) *URL {
	rwu := &URL{}

	for _, opt := range opts {
		opt(rwu)
	}

	return rwu
}

func WithLicense(license string) URLOption {
	return func(rwu *URL) {
		// only EU supported
		if licenseIsRegionEU(license) {
			rwu.RegionPrefix = regionEUPrefix
		} else {
			rwu.RegionPrefix = ""
		}
	}
}

func WithStaging(staging bool) URLOption {
	return func(rwu *URL) {
		rwu.Staging = staging
	}
}

func WithFedRAMP(fedramp bool) URLOption {
	return func(rwu *URL) {
		rwu.FedRAMP = fedramp
	}
}

func WithDataSourceName(dataSourceName string) URLOption {
	return func(rwu *URL) {
		if dataSourceName == "" {
			return
		}

		if rwu.Values != nil {
			rwu.Values.Add(prometheusServerQueryParam, dataSourceName)
		} else {
			rwu.Values = url.Values{
				prometheusServerQueryParam: []string{dataSourceName},
			}
		}
	}
}

func (rwu *URL) Build() (string, error) {
	if rwu.Staging && rwu.FedRAMP {
		return "", ErrFedRAMPStaging
	}
	if rwu.RegionPrefix == regionEUPrefix && rwu.FedRAMP {
		return "", ErrEuFedRAMP
	}

	var prefix string
	if rwu.Staging {
		prefix = environmentStagingPrefix
	}
	if rwu.FedRAMP {
		prefix = environmentFedRAMPPrefix
	}

	url := url.URL{
		Scheme:   remoteWriteScheme,
		Host:     fmt.Sprintf(remoteWriteHostTemplate, prefix, rwu.RegionPrefix),
		Path:     remoteWritePath,
		RawQuery: rwu.Values.Encode(),
	}

	return url.String(), nil
}
