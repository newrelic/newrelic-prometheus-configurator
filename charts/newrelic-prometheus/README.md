# newrelic-prometheus

A Helm chart to deploy Prometheus.

# Helm installation

You can install this chart using [`nri-bundle`](https://github.com/newrelic/helm-charts/tree/master/charts/nri-bundle) located in the
[helm-charts repository](https://github.com/newrelic/helm-charts) or directly from this repository by adding this Helm repository:

```shell
helm repo add nri-kubernetes https://newrelic.github.io/newrelic-prometheus-configurator
helm upgrade --install newrelic-prometheus-configurator/newrelic-prometheus -f your-custom-values.yaml
```

## Values managed globally

This chart implements the [New Relic's common Helm library](https://github.com/newrelic/helm-charts/tree/master/library/common-library) which
means that it honors a wide range of defaults and globals common to most New Relic Helm charts.

Options that can be defined globally include `affinity`, `nodeSelector`, `tolerations`, `proxy` and others. The full list can be found at
[user's guide of the common library](https://github.com/newrelic/helm-charts/blob/master/library/common-library/README.md).

## Chart particularities

### Default kubernetes jobs configuration

By default some kubernetes objects are discovered and scraped by prometheus. Taking into account the snippet from `values.yaml` below:

```yaml
config:
  kubernetes:
    jobs:
    - job_name_prefix: kubernetes-job
      target_discovery:
        pod: true
        endpoints: true
        filter:
          annotations:
            prometheus.io/scrape: true
```

All pod and endpoints with the `prometheus.io/scrape: true` annotation will be scraped by default.

### Self metrics

By default it is defined a job in `static_target.jobs` to obtain self-metrics. Particularlly, a snippet like the one
below is used.

```yaml
config:
  startic_targets:
    jobs:
    - job_name: self-metrics
      targets:
        - "localhost:9090"
      extra_metric_relabel_config:
        - source_labels: [__name__]
          regex: "\
            prometheus_agent_active_series|\
            prometheus_target_interval_length_seconds|\
            prometheus_target_scrape_pool_targets|\
            prometheus_remote_storage_samples_pending|\
            prometheus_remote_storage_samples_in_total|\
            prometheus_remote_storage_samples_retried_total|\
            prometheus_agent_corruptions_total|\
            prometheus_remote_storage_shards|\
            prometheus_sd_kubernetes_events_total|\
            prometheus_agent_checkpoint_creations_failed_total|\
            prometheus_agent_checkpoint_deletions_failed_total|\
            prometheus_remote_storage_samples_dropped_total|\
            prometheus_remote_storage_samples_failed_total|\
            prometheus_sd_kubernetes_http_request_total|\
            prometheus_agent_truncate_duration_seconds_sum|\
            prometheus_build_info|\
            process_resident_memory_bytes|\
            process_virtual_memory_bytes|\
            process_cpu_seconds_total"
          action: keep
```

### Low data mode

There are two mechanisms to reduce the amount of data that this integration sends to New Relic. See this snippet from the `values.yaml` file:
```yaml
lowDataMode: false

config:
  common:
    scrape_interval: 30s
```

You might set `lowDataMode` flag to `true` (it will filter some metrics which an also be collected using New Relic Kubernetes integration), check
`values.yaml` for details.

It is also possible to adjust how frequently prometheus scrapes the targets by setting up `config.common.scrape_interval` value.

### Affinities and tolerations

The New Relic common library allows to set affinities, tolerations, and node selectors globally using e.g. `.global.affinity` to ease the configuration
when you use this chart using `nri-bundle`. This chart has an extra level of granularity to the components that it deploys:
control plane, ksm, and kubelet.

Take this snippet as an example:
```yaml
global:
  affinity: {}
affinity: {}
```

The order to set the affinity is to set `affinity` field (at root level), if that value is empty, the chart fallbacks to `global.affinity`.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Sets pod/node affinities set almost globally. (See [Affinities and tolerations](README.md#affinities-and-tolerations)) |
| cluster | string | `""` | Name of the Kubernetes cluster monitored. Can be configured also with `global.cluster`. Note it will be set as an external label in prometheus configuration, it will have precedence over `config.common.external_labels.cluster_name` and `customAttributes.cluster_name``. |
| config | object | See `values.yaml` | It holds the New Relic Prometheus configuration. Here you can easily set up prometheus to get set metrics, discover ponds and endpoints Kubernetes and send metrics to New Relic using remote-write. |
| config.common | object | `{"scrape_interval":"30s"}` | Include global configuration for prometheus agent. See `values.yaml` for details. |
| config.extra_remote_write | object | `nil` | It includes additional remote-write configuration. Note this configuration is not parsed, so valid [prometheus remote_write configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) should be provided. |
| config.extra_scrape_configs | object | `{}` | It is possible to include extra scrape configuration in [prometheus format](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config). Please note, it should be valid prometheus configuration which will not be parsed by the chart. |
| config.kubernetes | object | See `values.yaml` | It allows defining scrape jobs for kubernetes in a simple way. See `values.yaml` for details. |
| config.newrelic_remote_write | object | `nil` | Newrelic remote-write configuration settings. It should work with the default value but if you need to set it up you can customize most [prometheus remote_write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) values, as described in `values.yaml`. |
| config.static_targets.jobs | list | See `values.yaml`. | For more info about static_targets, see `values.yaml` and [scrape_config prometheus documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config). By default it defines a job to get self-metrics. Please note, if you define `static_target.jobs` and would like to keep self-metrics you need to include a job like the one defined by default. |
| containerSecurityContext | object | `{}` | Sets security context (at container level). Can be configured also with `global.containerSecurityContext` |
| customAttributes | object | `{}` | Adds extra attributes to prometheus external labels. Can be configured also with `global.customAttributes`. Please note, values defined in `common.config.externar_labels` will have precedence over `customAttributes`. |
| customSecretLicenseKey | string | `""` | In case you don't want to have the license key in you values, this allows you to point to which secret key is the license key located. Can be configured also with `global.customSecretLicenseKey` |
| customSecretName | string | `""` | In case you don't want to have the license key in you values, this allows you to point to a user created secret to get the key from there. Can be configured also with `global.customSecretName` |
| dnsConfig | object | `{}` | Sets pod's dnsConfig. Can be configured also with `global.dnsConfig` |
| fullnameOverride | string | `""` | Override the full name of the release |
| hostNetwork | bool | `false` | Sets pod's hostNetwork. Can be configured also with `global.hostNetwork` |
| images.configurator | object | See `values.yaml` | Image for New Relic configurator. |
| images.prometheus | object | See `values.yaml` | Image for prometheus which is executed in agent mode. |
| images.pullSecrets | list | `[]` | The secrets that are needed to pull images from a custom registry. |
| labels | object | `{}` | Additional labels for chart objects. Can be configured also with `global.labels` |
| licenseKey | string | `""` | This set this license key to use. Can be configured also with `global.licenseKey` |
| lowDataMode | bool | false | Reduces number of metrics sent in order to reduce costs. It can be configured also with `global.lowDataMode`. Specifically, it makes prometheus stop reporting some kubernetes cluster specific metrics, you can see details in `static/lowdatamodedefaults.yaml`. |
| nameOverride | string | `""` | Override the name of the chart |
| nodeSelector | object | `{}` | Sets pod's node selector almost globally. (See [Affinities and tolerations](README.md#affinities-and-tolerations)) |
| nrStaging | bool | `false` | Send the metrics to the staging backend. Requires a valid staging license key. Can be configured also with `global.nrStaging` |
| podAnnotations | object | `{}` | Annotations to be added to all pods created by the integration. |
| podLabels | object | `{}` | Additional labels for chart pods. Can be configured also with `global.podLabels` |
| podSecurityContext | object | `{}` | Sets security context (at pod level). Can be configured also with `global.podSecurityContext` |
| priorityClassName | string | `""` | Sets pod's priorityClassName. Can be configured also with `global.priorityClassName` |
| rbac.create | bool | `true` | Whether the chart should automatically create the RBAC objects required to run. |
| rbac.pspEnabled | bool | `false` | Whether the chart should create Pod Security Policy objects. |
| serviceAccount | object | See `values.yaml` | Settings controlling ServiceAccount creation. |
| serviceAccount.create | bool | `true` | Whether the chart should automatically create the ServiceAccount objects required to run. |
| sharding | string | `nil` | Set up prometheus replicas to allow horizontal scalability. See `values.yaml` to set it up. |
| tolerations | list | `[]` | Sets pod's tolerations to node taints almost globally. (See [Affinities and tolerations](README.md#affinities-and-tolerations)) |

## Maintainers

* [alvarocabanas](https://github.com/alvarocabanas)
* [carlossscastro](https://github.com/carlossscastro)
* [sigilioso](https://github.com/sigilioso)
* [gsanchezgavier](https://github.com/gsanchezgavier)
* [kang-makes](https://github.com/kang-makes)
* [marcsanmi](https://github.com/marcsanmi)
* [paologallinaharbur](https://github.com/paologallinaharbur)
* [roobre](https://github.com/roobre)
