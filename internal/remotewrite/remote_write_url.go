// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

const (
	remoteWriteScheme        = "https"
	remoteWriteHostTemplate  = "%smetric-api.%snewrelic.com"
	remoteWritePath          = "prometheus/v1/write"
	environmentStagingPrefix = "staging-"
	environmentFedRAMPPrefix = "gov-"
	// prometheusServerQueryParam is added to remoteWrite url when nrConfig's name is defined.
	prometheusServerQueryParam = "prometheus_server"
	// collectorNameQueryParam is a NR identifier of the component collecting the data. This is added as query parameter of the PRW and converted
	// to collector.name to comply with NR standards.
	collectorNameQueryParam = "collector_name"
	collectorName           = "prometheus-agent"
	// collectorVersionQueryParam is a NR version of the component collecting the data. This is added as query parameter of the PRW and converted
	// to collector.version to comply with NR standards.
	collectorVersionQueryParam = "collector_version"
)

var (
	ErrFedRAMPRegions = errors.New("FedRAMP Region Error")
)

type URLOption func(url *URL)

type URL struct {
	Staging      bool
	FedRAMP      bool
	RegionPrefix string
	Values       url.Values
}

func NewURL(opts ...URLOption) *URL {
	u := &URL{Values: url.Values{}}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// licenseGetRegion returns license region or empty if none.
func licenseGetRegion(licenseKey string) string {
	regionLicenseRegex := regexp.MustCompile(`^([a-wyz]{2,3})(?:[0-9]{2})?x{1,2}`)
	matches := regionLicenseRegex.FindStringSubmatch(licenseKey)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func WithLicense(license string) URLOption {
	return func(u *URL) {
		var region = licenseGetRegion(license)
		if region != "" && region != "gov" {
			u.RegionPrefix = region + "."
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

		u.Values.Add(prometheusServerQueryParam, dataSourceName)
	}
}

func WithCollectorName(collectorName string) URLOption {
	return func(u *URL) {
		if collectorName == "" {
			return
		}

		u.Values.Add(collectorNameQueryParam, collectorName)
	}
}

func WithCollectorVersion(collectorVersion string) URLOption {
	return func(u *URL) {
		if collectorVersion == "" {
			return
		}

		u.Values.Add(collectorVersionQueryParam, collectorVersion)
	}
}

func (u *URL) Build() (string, error) {
	if u.Staging && u.FedRAMP {
		return "", fmt.Errorf("%w: There is no FedRamp compatible endpoints for staging", ErrFedRAMPRegions)
	}
	if u.RegionPrefix != "" && u.FedRAMP {
		return "", fmt.Errorf("%w: There is no FedRamp compatible endpoints for the region %s", ErrFedRAMPRegions, u.RegionPrefix)
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
