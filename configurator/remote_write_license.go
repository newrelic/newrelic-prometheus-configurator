// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package configurator

import "regexp"

// licenseKey code copied from infrastructure-agent to avoid the dependency. See:
// <https://github.com/newrelic/infrastructure-agent/blob/841750e718125c43a09b3270390299ba8468bff9/pkg/license/license.go#L8>

var regionLicenseRegex = regexp.MustCompile(`^([a-z]{2,3})`)

// licenseIsRegionEU returns true if license region is EU.
func licenseIsRegionEU(license string) bool {
	r := licenseGetRegion(license)
	// only EU supported
	if len(r) > 1 && r[:2] == "eu" {
		return true
	}
	return false
}

// licenseGetRegion returns license region or empty if none.
func licenseGetRegion(licenseKey string) string {
	matches := regionLicenseRegex.FindStringSubmatch(licenseKey)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
