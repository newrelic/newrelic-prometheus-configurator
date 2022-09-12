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
	u := &URL{}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func WithLicense(license string) URLOption {
	return func(u *URL) {
		if licenseIsRegionEU(license) {
			u.RegionPrefix = regionEUPrefix
		} else {
			u.RegionPrefix = ""
		}
	}
}

func WithStaging(staging bool) URLOption {
	return func(u *URL) {
		u.Staging = staging
	}
}

func WithFedRAMP(fedramp bool) URLOption {
	return func(u *URL) {
		u.FedRAMP = fedramp
	}
}

func WithDataSourceName(dataSourceName string) URLOption {
	return func(u *URL) {
		if dataSourceName == "" {
			return
		}

		if u.Values != nil {
			u.Values.Add(prometheusServerQueryParam, dataSourceName)
		} else {
			u.Values = url.Values{
				prometheusServerQueryParam: []string{dataSourceName},
			}
		}
	}
}

func (u *URL) Build() (string, error) {
	if u.Staging && u.FedRAMP {
		return "", ErrFedRAMPStaging
	}
	if u.RegionPrefix == regionEUPrefix && u.FedRAMP {
		return "", ErrEuFedRAMP
	}

	var prefix string
	if u.Staging {
		prefix = environmentStagingPrefix
	}
	if u.FedRAMP {
		prefix = environmentFedRAMPPrefix
	}

	url := url.URL{
		Scheme:   remoteWriteScheme,
		Host:     fmt.Sprintf(remoteWriteHostTemplate, prefix, u.RegionPrefix),
		Path:     remoteWritePath,
		RawQuery: u.Values.Encode(),
	}

	return url.String(), nil
}
