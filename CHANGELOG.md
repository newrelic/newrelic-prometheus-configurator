# Changelog
All notable changes are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#readme)

## Unreleased

### Integration filters
NewRelic provides a list of Dashboards, alerts and entities for several Services. The `integrations_filter` configuration
 allows to scrape only the targets having this experience out of the box.

 - If `integrations_filter` is enabled, then the kubernetes jobs scrape merely the targets having one of the specified labels matching
one of the values of `app_values`.

 - Under the hood, a `relabel_configs` with `action=keep` are generated, consider it in case any custom `extra_relabel_config` is needed.

### 🚀 Enhancements
- `Integration filters` feature is now supported.

## v0.1.0 - 2022-10-20

### 🛡️ Security notices
- Bumps `golang.org/x/sync` to `v0.1.0`.
- Replace `golang.org/x/text` to `v0.3.8` to fix `CVE-2022-32149`.

### 🚀 Enhancements
- Use Go 1.19.

### 🐞 Bug fixes
- Add missing `pod` name metadata to metrics scraped from `endpoints`.

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.25.2 to 0.25.3

## [0.0.2] - 2022-10-06
### Added
- `collector_name` metadata is added as part of New Relic remote write config, to identify the scraper component between other sources.

## [0.0.1] - 2022-09-20
### Added
- First Version of `newrelic-prometheus-configurator`.
