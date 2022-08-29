# newrelic-prometheus

![Version: 0.0.11-devel](https://img.shields.io/badge/Version-0.0.11--devel-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v2.37.0](https://img.shields.io/badge/AppVersion-v2.37.0-informational?style=flat-square)

A Helm chart to deploy Prometheus.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| alvarocabanas |  | <https://github.com/alvarocabanas> |
| carlossscastro |  | <https://github.com/carlossscastro> |
| sigilioso |  | <https://github.com/sigilioso> |
| gsanchezgavier |  | <https://github.com/gsanchezgavier> |
| kang-makes |  | <https://github.com/kang-makes> |
| marcsanmi |  | <https://github.com/marcsanmi> |
| paologallinaharbur |  | <https://github.com/paologallinaharbur> |
| roobre |  | <https://github.com/roobre> |

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://helm-charts.newrelic.com | common-library | 1.0.6 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Sets pod/node affinities set almost globally. (See [Affinities and tolerations](README.md#affinities-and-tolerations)) |
| cluster | string | `""` | Name of the Kubernetes cluster monitored. Can be configured also with `global.cluster`. Note it will be set as an external label in prometheus configuration, it will have precedence over `config.common.external_labels.cluster_name` and `customAttributes.cluster_name``. |
| config | object | It sets self metrics and target scraping for pod and endpoints kubernetes objects which include a specific annotation. Check `values.yaml` for details. | It holds the New Relic prometheus configuration. Check `values.yaml` for details. |
| config.common | object | `{"scrape_interval":"30s"}` | Include global configuration for prometheus agent. See `values.yaml` for details. |
| config.extra_remote_write | object | `nil` | It includes additional remote-write configuration. Note this configuration is not parsed, so valid [prometheus remote_write configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) should be provided. |
| config.extra_scrape_configs | object | `{}` | It is possible to include extra scrape configuration in [prometheus format](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config). Please note, it should be valid prometheus configuration which will not be parsed by the chart. |
| config.kubernetes | object | By default it defines jobs to scrape all pod and endpoints objects with `"prometheus.io/scrape"` annotation. | It allows defining scrape jobs for kubernetes in a simple way. |
| config.kubernetes.jobs[0].target_discovery.endpoints | bool | `true` | Whether endpoints should be discovered. |
| config.kubernetes.jobs[0].target_discovery.filter | object | `{"annotations":{"prometheus.io/scrape":true}}` | Define filtering criteria, it is possible to set labels and/or annotations. All filters will apply (defined filters are taken into account as an "and operation"). |
| config.kubernetes.jobs[0].target_discovery.pod | bool | `true` | Whether pods should be discovered. |
| config.newrelic_remote_write | object | `nil` | Newrelic remote-write configuration settings. It should work with the default value but if you need to set it up you can customize most [prometheus remote_write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) values, as described in `values.yaml`. |
| config.static_targets | object | Includes a self-metrics job, for more info see `values.yaml`. | For more info about static_targets, see `values.yaml` and [scrape_config prometheus documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config). Please note, if you define `static_target.jobs` and still need self metrics you should also include the job to do so (as defined by default). |
| containerSecurityContext | object | `{}` | Sets security context (at container level). Can be configured also with `global.containerSecurityContext` |
| customAttributes | object | `{}` | Adds extra attributes to prometheus external labels. Can be configured also with `global.customAttributes`. Please note, values defined in `common.config.externar_labels` will have precedence over `customAttributes`. |
| customSecretLicenseKey | string | `""` | In case you don't want to have the license key in you values, this allows you to point to which secret key is the license key located. Can be configured also with `global.customSecretLicenseKey` |
| customSecretName | string | `""` | In case you don't want to have the license key in you values, this allows you to point to a user created secret to get the key from there. Can be configured also with `global.customSecretName` |
| dnsConfig | object | `{}` | Sets pod's dnsConfig. Can be configured also with `global.dnsConfig` |
| fullnameOverride | string | `""` | Override the full name of the release |
| hostNetwork | bool | `false` | Sets pod's hostNetwork. Can be configured also with `global.hostNetwork` |
| images.configurator | object | `{"pullPolicy":"IfNotPresent","registry":"","repository":"newrelic/newrelic-prometheus-configurator","tag":""}` | Image for New Relic configurator. @default See `values.yaml` |
| images.prometheus | object | `{"pullPolicy":"IfNotPresent","registry":"","repository":"quay.io/prometheus/prometheus","tag":""}` | Image for prometheus which is executed in agent mode. @default See `values.yaml` |
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
| rbac | object | `{"create":true,"pspEnabled":false}` | Settings controlling RBAC objects creation. |
| rbac.create | bool | `true` | Whether the chart should automatically create the RBAC objects required to run. |
| rbac.pspEnabled | bool | `false` | Whether the chart should create Pod Security Policy objects. |
| serviceAccount | object | See `values.yaml` | Settings controlling ServiceAccount creation. |
| serviceAccount.create | bool | `true` | Whether the chart should automatically create the ServiceAccount objects required to run. |
| tolerations | list | `[]` | Sets pod's tolerations to node taints almost globally. (See [Affinities and tolerations](README.md#affinities-and-tolerations)) |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
