# Changelog
All notable changes are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#readme)

## Unreleased

## v1.4.0 - 2023-05-11

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.26.2 to 0.27.1
- Upgraded alpine from 3.17 to 3.18
- Upgraded golang.org/x/sync from 0.1.0 to 0.2.0

## v1.3.0 - 2023-03-20

### ⛓️ Dependencies
- Upgraded golang.org/x/net from 0.4.0 to 0.7.0
- Upgraded github.com/stretchr/testify from 1.8.1 to 1.8.2 - [Changelog 🔗](https://github.com/stretchr/testify/releases/tag/v1.8.2)
- Upgraded k8s.io/api from 0.26.1 to 0.26.2
- Upgraded k8s.io/client-go from 0.26.1 to 0.26.2
- Upgraded k8s.io/api from 0.26.2 to 0.26.3

## v1.2.0 - 2023-01-26

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.25.4 to 0.26.1
- Upgraded github.com/prometheus/prometheus from 0.37.3 to 0.37.5 - [Changelog 🔗](https://github.com/prometheus/prometheus/releases/tag/0.37.5)

## v1.1.0 - 2023-01-26

### 🚀 Enhancements
- Send collector_version query param to the Remote Write endpoint.

## v1.0.0 - 2022-11-29

### First stable release
- From now on the configuration is considered stable.

### ⛓️ Dependencies
- Upgraded github.com/prometheus/prometheus from 0.37.1 to 0.37.2 [Changelog](https://github.com/prometheus/prometheus/releases/tag/0.37.2)
- Upgraded k8s.io/apimachinery from 0.25.3 to 0.25.4
- Upgraded k8s.io/client-go from 0.25.3 to 0.25.4
- Upgraded alpine from 3.16 to 3.17
- Upgraded github.com/prometheus/prometheus from 0.37.2 to 0.37.3 [Changelog](https://github.com/prometheus/prometheus/releases/tag/0.37.3)

## v0.2.0 - 2022-11-03

### Integration filters
NewRelic provides a list of Dashboards, alerts and entities for several Services. The `integrations_filter` configuration allows to scrape only the targets having this experience out of the box.
- If `integrations_filter` is enabled, then the kubernetes jobs scrape merely the targets having one of the specified labels matching one of the values of `app_values`.
- Under the hood, a `relabel_configs` with `action=keep` are generated, consider it in case any custom `extra_relabel_config` is needed.

### 🚀 Enhancements
- `Integration filters` feature is now supported.

### ⛓️ Dependencies
- Upgraded github.com/stretchr/testify from 1.8.0 to 1.8.1 [Changelog](https://github.com/stretchr/testify/releases/tag/1.8.1)

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
