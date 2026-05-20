// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package remotewrite

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	remoteWriteScheme        = "https"
	remoteWriteHostTemplate  = "%smetric-api.%s%s%s"
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
	legacyCollectionDomain     = "newrelic.com"
	collectionDomain           = "nr-data.net"
)

var (
	ErrFedRAMPRegions = errors.New("FedRAMP Region Error")
	newDomainRegions  = []string{"jp"}
)

type URLOption func(url *URL)

type URL struct {
	Staging bool
	FedRAMP bool
	Region  string
	Values  url.Values
}

func NewURL(opts ...URLOption) *URL {
	u := &URL{Values: url.Values{}}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func getPrefix(staging, fedRAMP bool) string {
	if staging {
		return environmentStagingPrefix
	}
	if fedRAMP {
		return environmentFedRAMPPrefix
	}
	return ""
}

func getDomain(region string) string {
	for _, r := range newDomainRegions {
		if strings.EqualFold(region, r) {
			return collectionDomain
		}
	}
	return legacyCollectionDomain
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
		region := licenseGetRegion(license)
		if region != "" && region != "gov" {
			u.Region = region
		} else {
			u.Region = ""
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

func (u *URL) preconditions() error {
	if u.Staging && u.FedRAMP {
		return fmt.Errorf("%w: There is no FedRamp compatible endpoints for staging", ErrFedRAMPRegions)
	}
	if u.Region != "" && u.FedRAMP {
		return fmt.Errorf("%w: There is no FedRamp compatible endpoints for the region %s", ErrFedRAMPRegions, u.Region)
	}
	return nil
}

func (u *URL) Build() (string, error) {
	if err := u.preconditions(); err != nil {
		return "", err
	}

	prefix := getPrefix(u.Staging, u.FedRAMP)
	domain := getDomain(u.Region)
	regionPostfix := ""

	if u.Region != "" {
		regionPostfix = "."
	}

	url := url.URL{
		Scheme:   remoteWriteScheme,
		Host:     fmt.Sprintf(remoteWriteHostTemplate, prefix, u.Region, regionPostfix, domain),
		Path:     remoteWritePath,
		RawQuery: u.Values.Encode(),
	}

	return url.String(), nil
}
