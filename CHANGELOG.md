# Changelog
All notable changes are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#readme)

## Unreleased

## v2.1.0 - 2025-07-21

### ⛓️ Dependencies
- Updated alpine to v3.22.1
- Upgraded github.com/go-viper/mapstructure/v2 from 2.2.1 to 2.3.0 - [Changelog 🔗](https://github.com/go-viper/mapstructure/releases/tag/v2.3.0)

## v2.0.2 - 2025-07-14

### ⛓️ Dependencies
- Updated go to v1.24.5

## v2.0.1 - 2025-06-23

### ⛓️ Dependencies
- Updated go to v1.24.4

## v2.0.0 - 2025-06-02

### ⚠️️ Breaking changes ⚠️
- Upgrade to Prometheus v3. Please see our migration guide [here](https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/migrations/v2-upgrade) @dbudziwojskiNR [#458](https://github.com/newrelic/newrelic-prometheus-configurator/pull/458)

### 🚀 Enhancements
- Update namespace and serviceaccount name to match standard values from guided install @maxlemieux [#318](https://github.com/newrelic/newrelic-prometheus-configurator/pull/318)

### 🐞 Bug fixes
- Fix some bug in some file @dpacheconr [#457](https://github.com/newrelic/newrelic-prometheus-configurator/pull/457)

### ⛓️ Dependencies
- Updated go to v1.24.3
- Updated alpine to v3.22.0

## v1.20.0 - 2025-03-24

### 🚀 Enhancements
- Add v1.32 support and drop support for v1.27 @kpattaswamy [#447](https://github.com/newrelic/newrelic-prometheus-configurator/pull/447)

### ⛓️ Dependencies
- Updated golang.org/x/sync to v0.11.0
- Updated kubernetes packages to v0.32.3
- Updated go to v1.24.1

## v1.19.2 - 2025-02-17

### ⛓️ Dependencies
- Updated kubernetes packages to v0.32.1
- Updated alpine to v3.21.3

## v1.19.1 - 2025-01-27

### ⛓️ Dependencies
- Updated go to v1.23.5

## v1.19.0 - 2025-01-20

### ⛓️ Dependencies
- Upgraded golang.org/x/net from 0.30.0 to 0.33.0

## v1.18.6 - 2025-01-13

### ⛓️ Dependencies
- Updated alpine to v3.21.2

## v1.18.5 - 2024-12-23

### ⛓️ Dependencies
- Updated go to v1.23.4
- Updated golang.org/x/sync to v0.10.0
- Updated kubernetes packages to v0.32.0

## v1.18.4 - 2024-12-09

### ⛓️ Dependencies
- Updated alpine to v3.21.0

## v1.18.3 - 2024-11-25

### ⛓️ Dependencies
- Updated github.com/alecthomas/units digest
- Updated golang.org/x/sync to v0.9.0

## v1.18.2 - 2024-11-18

### ⛓️ Dependencies
- Updated go to v1.23.3

## v1.18.1 - 2024-11-04

### ⛓️ Dependencies
- Updated golang.org/x/sync to v0.8.0

## v1.18.0 - 2024-10-28

### 🚀 Enhancements
- Add 1.31 support and drop 1.26 @zeitlerc [#421](https://github.com/newrelic/newrelic-prometheus-configurator/pull/421)

## v1.17.4 - 2024-10-07

### ⛓️ Dependencies
- Updated alpine to v3.20.3

## v1.17.3 - 2024-07-29

### ⛓️ Dependencies
- Updated alpine to v3.20.2
- Updated kubernetes packages to v0.30.3

## v1.17.2 - 2024-07-22

### ⛓️ Dependencies
- Updated kubernetes packages to v0.30.2

## v1.17.1 - 2024-07-15

### ⛓️ Dependencies
- Updated github.com/alecthomas/units digest

## v1.17.0 - 2024-06-24

### 🚀 Enhancements
- Add 1.29 and 1.30 support and drop 1.25 and 1.24 @dbudziwojskiNR [#399](https://github.com/newrelic/newrelic-prometheus-configurator/pull/399)

### ⛓️ Dependencies
- Updated alpine to v3.20.1

## v1.16.4 - 2024-06-17

### ⛓️ Dependencies
- Updated go to v1.22.4

## v1.16.3 - 2024-06-10

### ⛓️ Dependencies
- Updated go to v1.22.3
- Updated kubernetes packages to v0.30.1

## v1.16.2 - 2024-05-27

### ⛓️ Dependencies
- Updated alpine to v3.20.0

## v1.16.1 - 2024-05-20

### 🐞 Bug fixes
- Update statefulset initcontainer to use container level security context from values [#385](https://github.com/newrelic/newrelic-prometheus-configurator/pull/385)

## v1.16.0 - 2024-04-29

### ⛓️ Dependencies
- Upgraded golang.org/x/net from 0.19.0 to 0.23.0

## v1.15.0 - 2024-03-25

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.3
- Upgraded google.golang.org/protobuf from 1.31.0 to 1.33.0

## v1.15.0 - 2024-03-08 

### 🚀 Enhancements

- Add `examples/metrictypes` directory for metric type conversion examples. @bpschmitt [#369] (https://github.com/newrelic/newrelic-prometheus-configurator/pull/369)

## v1.14.0 - 2024-02-26

### 🚀 Enhancements
- Add linux node selector @dbudziwojskiNR [#362](https://github.com/newrelic/newrelic-prometheus-configurator/pull/362)

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.2

## v1.13.0 - 2024-02-05

### 🚀 Enhancements
- Add Codecov @dbudziwojskiNR [#349](https://github.com/newrelic/newrelic-prometheus-configurator/pull/349)
- Add Codecov badge @dbudziwojskiNR [#352](https://github.com/newrelic/newrelic-prometheus-configurator/pull/352)

### 🐞 Bug fixes
- Fix flaky tests @dbudziwojskiNR [#344](https://github.com/newrelic/newrelic-prometheus-configurator/pull/344)

### ⛓️ Dependencies
- Updated alpine to v3.19.1

## v1.12.1 - 2024-01-22

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.1
- Updated go to v1.21.6

## v1.12.0 - 2024-01-15

### 🚀 Enhancements
- Trigger release creation by @juanjjaramillo [#341](https://github.com/newrelic/newrelic-prometheus-configurator/pull/341)

## v1.11.3 - 2024-01-09

### ⛓️ Dependencies
- Updated kubernetes packages to v0.29.0
- Updated golang.org/x/sync to v0.6.0

## v1.11.2 - 2023-12-18

### ⛓️ Dependencies
- Updated alpine to v3.19.0

## v1.11.1 - 2023-12-11

### ⛓️ Dependencies
- Updated go to v1.21.5

## v1.11.0 - 2023-12-06

### 🚀 Enhancements
- Update reusable workflow dependency by @juanjjaramillo [#325](https://github.com/newrelic/newrelic-prometheus-configurator/pull/325)

### 🐞 Bug fixes
- Fix integration test by @svetlanabrennan in [#324](https://github.com/newrelic/newrelic-prometheus-configurator/pull/324)

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.4
- Updated alpine to v3.18.5
- Updated github.com/alecthomas/units digest

## v1.10.0 - 2023-11-13

### 🚀 Enhancements
- Add full k8s 1.28 support by @svetlanabrennan in [#314](https://github.com/newrelic/newrelic-prometheus-configurator/pull/314)

## v1.9.0 - 2023-10-30

### 🚀 Enhancements
- Remove 1.23 support by @svetlanabrennan in [#303](https://github.com/newrelic/newrelic-prometheus-configurator/pull/303)
- Add k8s 1.28.0-rc.1 support by @svetlanabrennan in [#305](https://github.com/newrelic/newrelic-prometheus-configurator/pull/305)

## v1.8.1 - 2023-10-23

### ⛓️ Dependencies
- Updated kubernetes packages to v0.28.3

## v1.8.0 - 2023-10-16

### ⛓️ Dependencies
- Upgraded golang.org/x/net from 0.13.0 to 0.17.0
- Upgraded go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp from 0.32.0 to 0.44.0

## v1.7.1 - 2023-10-12

### ⛓️ Dependencies
- Updated go to 1.21

## v1.7.0 - 2023-10-11

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.28.1 to 0.28.2
- Upgraded alpine from 3.18.3 to 3.18.4
- Upgraded golang.org/x/sync from 0.3.0 to 0.4.0

## v1.6.0 - 2023-09-14

### 🚀 Enhancements
- Update K8s Versions in E2E Tests by @xqi-nr in [#265](https://github.com/newrelic/newrelic-prometheus-configurator/pull/265)

### 🐞 Bug fixes
- Add resource configuration option for initContainers. I accidentally push a commit to the repo main branch directly [https://github.com/newrelic/newrelic-prometheus-configurator/commit/cf752524b70fe4d351beb7da57a45d529b2aeece](https://github.com/newrelic/newrelic-prometheus-configurator/commit/cf752524b70fe4d351beb7da57a45d529b2aeece)

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.28.0 to 0.28.1
- Upgraded k8s.io/apimachinery from 0.28.1 to 0.28.2

## v1.5.0 - 2023-08-21

### ⛓️ Dependencies
- Upgraded alpine from 3.18.2 to 3.18.3
- Upgraded golang.org/x/sync from 0.2.0 to 0.3.0
- Upgraded k8s.io/client-go from 0.27.2 to 0.28.0

## v1.4.2 - 2023-06-15

### ⛓️ Dependencies
- Upgraded github.com/stretchr/testify from 1.8.2 to 1.8.4 - [Changelog 🔗](https://github.com/stretchr/testify/releases/tag/v1.8.4)
- Upgraded github.com/sirupsen/logrus from 1.9.0 to 1.9.3 - [Changelog 🔗](https://github.com/sirupsen/logrus/releases/tag/v1.9.3)

## v1.4.2 - 2023-06-08

### ⛓️ Dependencies
- Upgraded github.com/stretchr/testify from 1.8.2 to 1.8.4 - [Changelog 🔗](https://github.com/stretchr/testify/releases/tag/v1.8.4)
- Upgraded github.com/sirupsen/logrus from 1.9.0 to 1.9.3 - [Changelog 🔗](https://github.com/sirupsen/logrus/releases/tag/v1.9.3)

## v1.4.1 - 2023-06-03

### ⛓️ Dependencies
- Upgraded k8s.io/client-go from 0.27.1 to 0.27.2

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
